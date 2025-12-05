package telegram

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/efigueroa/telegram-music-monitor/database"
	"github.com/efigueroa/telegram-music-monitor/lidarr"
	"github.com/efigueroa/telegram-music-monitor/musicbrainz"
	"github.com/efigueroa/telegram-music-monitor/platforms"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
)

// Client wraps the Telegram client and handles message processing
type Client struct {
	apiID         int
	apiHash       string
	phone         string
	sessionPath   string
	db            *database.DB
	spotifyClient *platforms.SpotifyClient
	youtubeClient *platforms.YouTubeClient
	mbClient      *musicbrainz.Client
	lidarrClient  *lidarr.Client
}

// Config contains configuration for the Telegram client
type Config struct {
	APIID         int
	APIHash       string
	Phone         string
	SessionPath   string
	DB            *database.DB
	SpotifyClient *platforms.SpotifyClient
	YouTubeClient *platforms.YouTubeClient
	MBClient      *musicbrainz.Client
	LidarrClient  *lidarr.Client
}

// Regular expression to extract URLs from messages
var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

// NewClient creates a new Telegram client
func NewClient(cfg Config) *Client {
	return &Client{
		apiID:         cfg.APIID,
		apiHash:       cfg.APIHash,
		phone:         cfg.Phone,
		sessionPath:   cfg.SessionPath,
		db:            cfg.DB,
		spotifyClient: cfg.SpotifyClient,
		youtubeClient: cfg.YouTubeClient,
		mbClient:      cfg.MBClient,
		lidarrClient:  cfg.LidarrClient,
	}
}

// Run starts the Telegram client and begins monitoring messages
func (c *Client) Run(ctx context.Context) error {
	// Create a new Telegram client
	client := telegram.NewClient(c.apiID, c.apiHash, telegram.Options{
		// Store session data in a file so we don't need to re-authenticate every time
		SessionStorage: &telegram.FileSessionStorage{
			Path: c.sessionPath,
		},
	})

	// Run the client
	return client.Run(ctx, func(ctx context.Context) error {
		// Get the API client
		api := client.API()

		// Authenticate if not already authenticated
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return fmt.Errorf("failed to get auth status: %w", err)
		}

		if !status.Authorized {
			// Need to authenticate
			log.Println("Not authenticated, starting authentication flow...")

			// Use phone authentication flow
			if _, err := client.Auth().Phone(ctx, c.phone, auth.CodeAuthenticatorFunc(
				func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
					// In a real application, you'd get this from user input
					// For Docker, you'll need to enter this on first run
					var code string
					fmt.Print("Enter code: ")
					fmt.Scanln(&code)
					return code, nil
				},
			)); err != nil {
				return fmt.Errorf("failed to authenticate: %w", err)
			}

			log.Println("Successfully authenticated!")
		} else {
			log.Println("Already authenticated")
		}

		// Set up update handler to receive new messages
		gaps := updates.New(updates.Config{
			Handler: c, // We implement the Handler interface below
		})

		// Start the update handler
		log.Println("Starting message monitoring...")
		return gaps.Run(ctx, api, tg.NewUpdateDispatcher(api).OnNewMessage(c.handleMessage))
	})
}

// handleMessage is called when a new message is received
func (c *Client) handleMessage(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	msg, ok := update.Message.(*tg.Message)
	if !ok {
		// Not a regular message, ignore
		return nil
	}

	// Get message text
	text := msg.Message
	if text == "" {
		// No text content, ignore
		return nil
	}

	// Extract all URLs from the message
	urls := urlRegex.FindAllString(text, -1)
	if len(urls) == 0 {
		// No URLs found, ignore
		return nil
	}

	// Get chat/group information
	peer := msg.PeerID
	var groupID int64
	var groupName string

	// Determine the type of chat
	switch p := peer.(type) {
	case *tg.PeerChannel:
		groupID = p.ChannelID
		// Get channel info to get the name
		if ch, ok := e.Channels[p.ChannelID]; ok {
			groupName = ch.GetTitle()
		}
	case *tg.PeerChat:
		groupID = p.ChatID
		if ch, ok := e.Chats[p.ChatID]; ok {
			if chat, ok := ch.(*tg.Chat); ok {
				groupName = chat.Title
			}
		}
	case *tg.PeerUser:
		// Private message, we might not want to monitor these
		// For now, we'll skip private messages
		return nil
	default:
		return nil
	}

	// Check if monitoring is enabled for this group
	groupConfig, err := c.db.GetGroupConfig(groupID)
	if err != nil {
		log.Printf("Error getting group config for %d: %v", groupID, err)
		return nil
	}

	if !groupConfig.Enabled {
		// Monitoring disabled for this group
		return nil
	}

	// Get sender information
	var senderID int64
	var senderUsername string

	if msg.FromID != nil {
		if fromUser, ok := msg.FromID.(*tg.PeerUser); ok {
			senderID = fromUser.UserID
			if user, ok := e.Users[fromUser.UserID]; ok {
				senderUsername = user.GetUsername()
			}
		}
	}

	// Process each URL
	for _, url := range urls {
		// Process the URL in a goroutine to not block message handling
		go c.processURL(context.Background(), url, msg.ID, groupID, groupName, senderID, senderUsername)
	}

	return nil
}

// processURL handles a detected music URL
func (c *Client) processURL(ctx context.Context, url string, messageID int, groupID int64, groupName string, senderID int64, senderUsername string) {
	log.Printf("Processing URL: %s from group: %s", url, groupName)

	var platform string
	var artistName string
	var albumName string
	var songName string
	var musicbrainzID string

	// Determine the platform and extract metadata
	if platforms.IsSpotifyURL(url) {
		platform = "spotify"
		if err := c.processSpotifyURL(ctx, url, &artistName, &albumName, &songName); err != nil {
			log.Printf("Error processing Spotify URL %s: %v", url, err)
			return
		}
	} else if platforms.IsYouTubeURL(url) {
		platform = "youtube"
		if err := c.processYouTubeURL(ctx, url, &artistName, &albumName, &songName); err != nil {
			log.Printf("Error processing YouTube URL %s: %v", url, err)
			return
		}
	} else {
		// Not a supported platform
		return
	}

	// Enrich with MusicBrainz data
	if artistName != "" {
		artists, err := c.mbClient.SearchArtist(ctx, artistName)
		if err != nil {
			log.Printf("Error searching MusicBrainz for artist %s: %v", artistName, err)
		} else if best := musicbrainz.GetBestArtistMatch(artists); best != nil {
			musicbrainzID = best.ID
			// Use MusicBrainz canonical name
			artistName = best.Name
		}
	}

	// Insert into database
	detectedMusic := &database.DetectedMusic{
		MessageID:      int64(messageID),
		GroupID:        groupID,
		GroupName:      groupName,
		SenderID:       senderID,
		SenderUsername: senderUsername,
		OriginalLink:   url,
		Platform:       platform,
		DetectedAt:     time.Now(),
		SongName:       songName,
		AlbumName:      albumName,
		ArtistName:     artistName,
		MusicBrainzID:  musicbrainzID,
		LidarrStatus:   "pending",
	}

	id, err := c.db.InsertDetectedMusic(detectedMusic)
	if err != nil {
		log.Printf("Error inserting detected music: %v", err)
		return
	}

	// Add to Lidarr if we have a MusicBrainz ID and artist name
	if musicbrainzID != "" && artistName != "" {
		if err := c.addToLidarr(ctx, id, artistName, musicbrainzID); err != nil {
			log.Printf("Error adding to Lidarr: %v", err)
		}
	} else {
		log.Printf("Skipping Lidarr add: missing MusicBrainz ID or artist name")
		c.db.UpdateLidarrStatus(id, "skipped", "Missing MusicBrainz ID or artist name")
	}
}

// processSpotifyURL extracts metadata from a Spotify URL
func (c *Client) processSpotifyURL(ctx context.Context, url string, artistName, albumName, songName *string) error {
	linkType := platforms.GetSpotifyLinkType(url)

	switch linkType {
	case "track":
		track, err := c.spotifyClient.GetTrack(ctx, url)
		if err != nil {
			return err
		}
		*songName = track.Name
		*albumName = track.Album
		if len(track.Artists) > 0 {
			*artistName = track.Artists[0]
		}

	case "album":
		album, err := c.spotifyClient.GetAlbum(ctx, url)
		if err != nil {
			return err
		}
		*albumName = album.Name
		if len(album.Artists) > 0 {
			*artistName = album.Artists[0]
		}

	case "playlist":
		playlist, err := c.spotifyClient.GetPlaylist(ctx, url)
		if err != nil {
			return err
		}
		// For playlists, we'll process each track separately
		// For now, just take the first track
		if len(playlist.Tracks) > 0 {
			*songName = playlist.Tracks[0].Name
			*albumName = playlist.Tracks[0].Album
			if len(playlist.Tracks[0].Artists) > 0 {
				*artistName = playlist.Tracks[0].Artists[0]
			}
		}

	case "artist":
		artist, err := c.spotifyClient.GetArtist(ctx, url)
		if err != nil {
			return err
		}
		*artistName = artist.Name

	default:
		return fmt.Errorf("unsupported Spotify link type: %s", linkType)
	}

	return nil
}

// processYouTubeURL extracts metadata from a YouTube URL
func (c *Client) processYouTubeURL(ctx context.Context, url string, artistName, albumName, songName *string) error {
	linkType := platforms.GetYouTubeLinkType(url)

	switch linkType {
	case "video":
		video, err := c.youtubeClient.GetVideo(ctx, url)
		if err != nil {
			return err
		}
		*songName = video.SongName
		if len(video.Artists) > 0 {
			*artistName = video.Artists[0]
		}

	case "playlist":
		playlist, err := c.youtubeClient.GetPlaylist(ctx, url)
		if err != nil {
			return err
		}
		// For playlists, take the first video
		if len(playlist.Videos) > 0 {
			*songName = playlist.Videos[0].SongName
			if len(playlist.Videos[0].Artists) > 0 {
				*artistName = playlist.Videos[0].Artists[0]
			}
		}

	default:
		return fmt.Errorf("unsupported YouTube link type: %s", linkType)
	}

	return nil
}

// addToLidarr adds an artist to Lidarr
func (c *Client) addToLidarr(ctx context.Context, dbID int64, artistName, musicbrainzID string) error {
	// Check if artist already exists in Lidarr
	exists, err := c.lidarrClient.ArtistExists(ctx, musicbrainzID)
	if err != nil {
		c.db.UpdateLidarrStatus(dbID, "failed", fmt.Sprintf("Error checking Lidarr: %v", err))
		return err
	}

	if exists {
		log.Printf("Artist %s already exists in Lidarr", artistName)
		c.db.UpdateLidarrStatus(dbID, "skipped", "Artist already exists in Lidarr")
		return nil
	}

	// Add to Lidarr
	if err := c.lidarrClient.AddArtist(ctx, artistName, musicbrainzID); err != nil {
		c.db.UpdateLidarrStatus(dbID, "failed", fmt.Sprintf("Error adding to Lidarr: %v", err))
		return err
	}

	log.Printf("Successfully added artist %s to Lidarr", artistName)
	c.db.UpdateLidarrStatus(dbID, "added", "")
	return nil
}

// Handle is required to implement the updates.Handler interface
// This is called when updates need to be processed
func (c *Client) Handle(ctx context.Context, updates tg.UpdatesClass) error {
	// The actual message handling is done in handleMessage
	// This is just required by the interface
	return nil
}
