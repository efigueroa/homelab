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

// SpotifyClient handles Spotify API interactions
type SpotifyClient struct {
	clientID     string
	clientSecret string
	accessToken  string
	tokenExpiry  time.Time
	httpClient   *http.Client
}

// SpotifyTrack represents a single track from Spotify
type SpotifyTrack struct {
	Name    string   // Track name
	Artists []string // List of artist names
	Album   string   // Album name
}

// SpotifyAlbum represents an album from Spotify
type SpotifyAlbum struct {
	Name    string   // Album name
	Artists []string // List of artist names
	Tracks  []string // List of track names in the album
}

// SpotifyPlaylist represents a playlist from Spotify
type SpotifyPlaylist struct {
	Name   string          // Playlist name
	Tracks []SpotifyTrack // List of tracks in the playlist
}

// SpotifyArtist represents an artist from Spotify
type SpotifyArtist struct {
	Name string // Artist name
}

// Regular expressions to match Spotify URLs
var (
	spotifyTrackRegex    = regexp.MustCompile(`spotify\.com/track/([a-zA-Z0-9]+)`)
	spotifyAlbumRegex    = regexp.MustCompile(`spotify\.com/album/([a-zA-Z0-9]+)`)
	spotifyPlaylistRegex = regexp.MustCompile(`spotify\.com/playlist/([a-zA-Z0-9]+)`)
	spotifyArtistRegex   = regexp.MustCompile(`spotify\.com/artist/([a-zA-Z0-9]+)`)
)

// NewSpotifyClient creates a new Spotify API client
func NewSpotifyClient(clientID, clientSecret string) *SpotifyClient {
	return &SpotifyClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

// authenticate gets a new access token from Spotify using client credentials flow
// This is called automatically when needed
func (s *SpotifyClient) authenticate(ctx context.Context) error {
	// Spotify uses OAuth 2.0 Client Credentials flow
	// We need to exchange our client ID and secret for an access token

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	// Set basic auth header with client credentials
	req.SetBasicAuth(s.clientID, s.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("spotify auth failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response to get the access token
	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"` // Token lifetime in seconds
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	// Store the token and calculate expiry time
	s.accessToken = result.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return nil
}

// ensureAuthenticated checks if we have a valid token, and gets a new one if needed
func (s *SpotifyClient) ensureAuthenticated(ctx context.Context) error {
	// If token is missing or expired, authenticate
	if s.accessToken == "" || time.Now().After(s.tokenExpiry) {
		return s.authenticate(ctx)
	}
	return nil
}

// IsSpotifyURL checks if a URL is a Spotify link
func IsSpotifyURL(urlStr string) bool {
	return strings.Contains(urlStr, "spotify.com/")
}

// GetSpotifyLinkType determines what type of Spotify link it is
func GetSpotifyLinkType(urlStr string) string {
	if spotifyTrackRegex.MatchString(urlStr) {
		return "track"
	}
	if spotifyAlbumRegex.MatchString(urlStr) {
		return "album"
	}
	if spotifyPlaylistRegex.MatchString(urlStr) {
		return "playlist"
	}
	if spotifyArtistRegex.MatchString(urlStr) {
		return "artist"
	}
	return "unknown"
}

// GetTrack fetches track information from Spotify
func (s *SpotifyClient) GetTrack(ctx context.Context, urlStr string) (*SpotifyTrack, error) {
	// Ensure we have a valid access token
	if err := s.ensureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Extract track ID from URL using regex
	matches := spotifyTrackRegex.FindStringSubmatch(urlStr)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid Spotify track URL")
	}
	trackID := matches[1]

	// Call Spotify API to get track details
	apiURL := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header with our access token
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract artist names from the nested structure
	artists := make([]string, len(result.Artists))
	for i, artist := range result.Artists {
		artists[i] = artist.Name
	}

	return &SpotifyTrack{
		Name:    result.Name,
		Artists: artists,
		Album:   result.Album.Name,
	}, nil
}

// GetAlbum fetches album information from Spotify
func (s *SpotifyClient) GetAlbum(ctx context.Context, urlStr string) (*SpotifyAlbum, error) {
	if err := s.ensureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	matches := spotifyAlbumRegex.FindStringSubmatch(urlStr)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid Spotify album URL")
	}
	albumID := matches[1]

	apiURL := fmt.Sprintf("https://api.spotify.com/v1/albums/%s", albumID)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Tracks struct {
			Items []struct {
				Name string `json:"name"`
			} `json:"items"`
		} `json:"tracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	artists := make([]string, len(result.Artists))
	for i, artist := range result.Artists {
		artists[i] = artist.Name
	}

	tracks := make([]string, len(result.Tracks.Items))
	for i, track := range result.Tracks.Items {
		tracks[i] = track.Name
	}

	return &SpotifyAlbum{
		Name:    result.Name,
		Artists: artists,
		Tracks:  tracks,
	}, nil
}

// GetPlaylist fetches playlist information from Spotify
func (s *SpotifyClient) GetPlaylist(ctx context.Context, urlStr string) (*SpotifyPlaylist, error) {
	if err := s.ensureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	matches := spotifyPlaylistRegex.FindStringSubmatch(urlStr)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid Spotify playlist URL")
	}
	playlistID := matches[1]

	apiURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", playlistID)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Name   string `json:"name"`
		Tracks struct {
			Items []struct {
				Track struct {
					Name    string `json:"name"`
					Artists []struct {
						Name string `json:"name"`
					} `json:"artists"`
					Album struct {
						Name string `json:"name"`
					} `json:"album"`
				} `json:"track"`
			} `json:"items"`
		} `json:"tracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	tracks := make([]SpotifyTrack, 0, len(result.Tracks.Items))
	for _, item := range result.Tracks.Items {
		// Skip null tracks (can happen in playlists)
		if item.Track.Name == "" {
			continue
		}

		artists := make([]string, len(item.Track.Artists))
		for i, artist := range item.Track.Artists {
			artists[i] = artist.Name
		}

		tracks = append(tracks, SpotifyTrack{
			Name:    item.Track.Name,
			Artists: artists,
			Album:   item.Track.Album.Name,
		})
	}

	return &SpotifyPlaylist{
		Name:   result.Name,
		Tracks: tracks,
	}, nil
}

// GetArtist fetches artist information from Spotify
func (s *SpotifyClient) GetArtist(ctx context.Context, urlStr string) (*SpotifyArtist, error) {
	if err := s.ensureAuthenticated(ctx); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	matches := spotifyArtistRegex.FindStringSubmatch(urlStr)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid Spotify artist URL")
	}
	artistID := matches[1]

	apiURL := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", artistID)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &SpotifyArtist{
		Name: result.Name,
	}, nil
}
