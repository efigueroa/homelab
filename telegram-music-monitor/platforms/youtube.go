package platforms

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// YouTubeClient handles YouTube API interactions
type YouTubeClient struct {
	apiKey     string
	httpClient *http.Client
}

// YouTubeVideo represents a single video from YouTube
type YouTubeVideo struct {
	Title       string   // Video title
	Description string   // Video description
	Artists     []string // Extracted artist names (from title parsing)
	SongName    string   // Extracted song name (from title parsing)
}

// YouTubePlaylist represents a playlist from YouTube
type YouTubePlaylist struct {
	Title  string          // Playlist title
	Videos []YouTubeVideo // Videos in the playlist
}

// Regular expressions to match YouTube URLs
var (
	// Matches youtube.com/watch?v=VIDEO_ID or youtu.be/VIDEO_ID
	youtubeVideoRegex = regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/)([a-zA-Z0-9_-]{11})`)
	// Matches youtube.com/playlist?list=PLAYLIST_ID
	youtubePlaylistRegex = regexp.MustCompile(`youtube\.com/playlist\?list=([a-zA-Z0-9_-]+)`)
)

// Regular expression to parse common music video title formats
// Common formats: "Artist - Song", "Artist: Song", "Song by Artist", etc.
var musicTitleRegex = regexp.MustCompile(`^([^-:|]+?)(?:\s*[-:|]\s*|\s+by\s+)(.+?)(?:\s*\(.*\)|\s*\[.*\]|$)`)

// NewYouTubeClient creates a new YouTube API client
func NewYouTubeClient(apiKey string) *YouTubeClient {
	return &YouTubeClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// IsYouTubeURL checks if a URL is a YouTube link
func IsYouTubeURL(urlStr string) bool {
	return strings.Contains(urlStr, "youtube.com/") || strings.Contains(urlStr, "youtu.be/")
}

// GetYouTubeLinkType determines what type of YouTube link it is
func GetYouTubeLinkType(urlStr string) string {
	if youtubePlaylistRegex.MatchString(urlStr) {
		return "playlist"
	}
	if youtubeVideoRegex.MatchString(urlStr) {
		return "video"
	}
	return "unknown"
}

// parseMusicTitle attempts to extract artist and song name from a YouTube video title
// Common formats:
// - "Artist - Song Title"
// - "Artist: Song Title"
// - "Song Title by Artist"
// - "Artist | Song Title"
func parseMusicTitle(title string) (artist string, song string) {
	// Try the regex pattern first
	matches := musicTitleRegex.FindStringSubmatch(title)
	if len(matches) >= 3 {
		// Depending on the format, matches[1] and matches[2] could be artist/song or song/artist
		part1 := strings.TrimSpace(matches[1])
		part2 := strings.TrimSpace(matches[2])

		// If "by" is in the title, part1 is song and part2 is artist
		if strings.Contains(strings.ToLower(title), " by ") {
			return part2, part1
		}

		// Otherwise, assume part1 is artist and part2 is song
		return part1, part2
	}

	// If parsing fails, return the whole title as the song name
	return "", title
}

// GetVideo fetches video information from YouTube
func (y *YouTubeClient) GetVideo(ctx context.Context, urlStr string) (*YouTubeVideo, error) {
	// Extract video ID from URL using regex
	matches := youtubeVideoRegex.FindStringSubmatch(urlStr)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid YouTube video URL")
	}
	videoID := matches[1]

	// Build YouTube Data API request
	// We use the videos endpoint to get video details
	apiURL := "https://www.googleapis.com/youtube/v3/videos"
	params := url.Values{}
	params.Set("part", "snippet") // Request the snippet part which contains title and description
	params.Set("id", videoID)
	params.Set("key", y.apiKey)

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := y.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("youtube API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result struct {
		Items []struct {
			Snippet struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	}

	title := result.Items[0].Snippet.Title
	description := result.Items[0].Snippet.Description

	// Try to parse artist and song from the title
	artist, song := parseMusicTitle(title)

	video := &YouTubeVideo{
		Title:       title,
		Description: description,
		SongName:    song,
	}

	// If we successfully extracted an artist, add it to the list
	if artist != "" {
		video.Artists = []string{artist}
	}

	return video, nil
}

// GetPlaylist fetches playlist information from YouTube
func (y *YouTubeClient) GetPlaylist(ctx context.Context, urlStr string) (*YouTubePlaylist, error) {
	// Extract playlist ID from URL
	matches := youtubePlaylistRegex.FindStringSubmatch(urlStr)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid YouTube playlist URL")
	}
	playlistID := matches[1]

	// First, get playlist details
	apiURL := "https://www.googleapis.com/youtube/v3/playlists"
	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("id", playlistID)
	params.Set("key", y.apiKey)

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := y.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("youtube API returned status %d: %s", resp.StatusCode, string(body))
	}

	var playlistResult struct {
		Items []struct {
			Snippet struct {
				Title string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&playlistResult); err != nil {
		return nil, fmt.Errorf("failed to decode playlist response: %w", err)
	}

	if len(playlistResult.Items) == 0 {
		return nil, fmt.Errorf("playlist not found")
	}

	playlistTitle := playlistResult.Items[0].Snippet.Title

	// Now get playlist items (videos in the playlist)
	// Note: This returns up to 50 items. For larger playlists, we'd need to handle pagination
	itemsURL := "https://www.googleapis.com/youtube/v3/playlistItems"
	itemsParams := url.Values{}
	itemsParams.Set("part", "snippet")
	itemsParams.Set("playlistId", playlistID)
	itemsParams.Set("maxResults", "50") // Get up to 50 videos
	itemsParams.Set("key", y.apiKey)

	fullItemsURL := fmt.Sprintf("%s?%s", itemsURL, itemsParams.Encode())
	itemsReq, err := http.NewRequestWithContext(ctx, "GET", fullItemsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create items request: %w", err)
	}

	itemsResp, err := y.httpClient.Do(itemsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute items request: %w", err)
	}
	defer itemsResp.Body.Close()

	if itemsResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(itemsResp.Body)
		return nil, fmt.Errorf("youtube API returned status %d: %s", itemsResp.StatusCode, string(body))
	}

	var itemsResult struct {
		Items []struct {
			Snippet struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(itemsResp.Body).Decode(&itemsResult); err != nil {
		return nil, fmt.Errorf("failed to decode items response: %w", err)
	}

	// Parse each video in the playlist
	videos := make([]YouTubeVideo, 0, len(itemsResult.Items))
	for _, item := range itemsResult.Items {
		title := item.Snippet.Title
		description := item.Snippet.Description

		// Try to parse artist and song from the title
		artist, song := parseMusicTitle(title)

		video := YouTubeVideo{
			Title:       title,
			Description: description,
			SongName:    song,
		}

		if artist != "" {
			video.Artists = []string{artist}
		}

		videos = append(videos, video)
	}

	return &YouTubePlaylist{
		Title:  playlistTitle,
		Videos: videos,
	}, nil
}
