package musicbrainz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client handles MusicBrainz API interactions
// MusicBrainz is a free, open music encyclopedia that collects music metadata
type Client struct {
	httpClient *http.Client
	userAgent  string // MusicBrainz requires a User-Agent header
	rateLimit  time.Duration // Time to wait between requests (MusicBrainz requires rate limiting)
	lastRequest time.Time
}

// Artist represents a MusicBrainz artist
type Artist struct {
	ID   string // MusicBrainz ID (MBID)
	Name string // Artist name
	Score int   // Search relevance score (0-100)
}

// Release represents a MusicBrainz release (album)
type Release struct {
	ID     string // MusicBrainz ID (MBID)
	Title  string // Album/Release title
	Artist string // Primary artist name
	Score  int    // Search relevance score (0-100)
}

// Recording represents a MusicBrainz recording (track/song)
type Recording struct {
	ID     string // MusicBrainz ID (MBID)
	Title  string // Track/Song title
	Artist string // Primary artist name
	Score  int    // Search relevance score (0-100)
}

// NewClient creates a new MusicBrainz API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		// MusicBrainz requires a descriptive User-Agent
		userAgent: "TelegramMusicMonitor/1.0 (https://github.com/efigueroa/telegram-music-monitor)",
		// MusicBrainz has a rate limit: no more than 1 request per second
		rateLimit: 1 * time.Second,
	}
}

// waitForRateLimit ensures we don't exceed MusicBrainz's rate limit
func (c *Client) waitForRateLimit() {
	if !c.lastRequest.IsZero() {
		elapsed := time.Since(c.lastRequest)
		if elapsed < c.rateLimit {
			time.Sleep(c.rateLimit - elapsed)
		}
	}
	c.lastRequest = time.Now()
}

// SearchArtist searches for an artist by name
// Returns a list of potential matches sorted by relevance
func (c *Client) SearchArtist(ctx context.Context, artistName string) ([]Artist, error) {
	c.waitForRateLimit()

	// Build the search query
	// MusicBrainz uses Lucene query syntax
	query := url.QueryEscape(fmt.Sprintf("artist:%s", artistName))
	apiURL := fmt.Sprintf("https://musicbrainz.org/ws/2/artist?query=%s&fmt=json&limit=5", query)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// MusicBrainz requires a User-Agent header
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("musicbrainz API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result struct {
		Artists []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Score int    `json:"score"` // MusicBrainz returns a score for search relevance
		} `json:"artists"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our Artist struct
	artists := make([]Artist, len(result.Artists))
	for i, a := range result.Artists {
		artists[i] = Artist{
			ID:   a.ID,
			Name: a.Name,
			Score: a.Score,
		}
	}

	return artists, nil
}

// SearchRelease searches for a release (album) by artist and title
func (c *Client) SearchRelease(ctx context.Context, artistName, albumTitle string) ([]Release, error) {
	c.waitForRateLimit()

	// Build search query with both artist and release title for better matches
	query := url.QueryEscape(fmt.Sprintf("artist:%s AND release:%s", artistName, albumTitle))
	apiURL := fmt.Sprintf("https://musicbrainz.org/ws/2/release?query=%s&fmt=json&limit=5", query)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("musicbrainz API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Releases []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Score int    `json:"score"`
			// Artist credits is an array, we'll take the first one
			ArtistCredit []struct {
				Name string `json:"name"`
			} `json:"artist-credit"`
		} `json:"releases"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	releases := make([]Release, 0, len(result.Releases))
	for _, r := range result.Releases {
		artistName := ""
		if len(r.ArtistCredit) > 0 {
			artistName = r.ArtistCredit[0].Name
		}

		releases = append(releases, Release{
			ID:     r.ID,
			Title:  r.Title,
			Artist: artistName,
			Score:  r.Score,
		})
	}

	return releases, nil
}

// SearchRecording searches for a recording (track/song) by artist and title
func (c *Client) SearchRecording(ctx context.Context, artistName, songTitle string) ([]Recording, error) {
	c.waitForRateLimit()

	// Build search query
	var queryParts []string
	if artistName != "" {
		queryParts = append(queryParts, fmt.Sprintf("artist:%s", artistName))
	}
	if songTitle != "" {
		queryParts = append(queryParts, fmt.Sprintf("recording:%s", songTitle))
	}

	if len(queryParts) == 0 {
		return nil, fmt.Errorf("must provide at least artist or song title")
	}

	query := url.QueryEscape(strings.Join(queryParts, " AND "))
	apiURL := fmt.Sprintf("https://musicbrainz.org/ws/2/recording?query=%s&fmt=json&limit=5", query)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("musicbrainz API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Recordings []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Score int    `json:"score"`
			ArtistCredit []struct {
				Name string `json:"name"`
			} `json:"artist-credit"`
		} `json:"recordings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	recordings := make([]Recording, 0, len(result.Recordings))
	for _, r := range result.Recordings {
		artistName := ""
		if len(r.ArtistCredit) > 0 {
			artistName = r.ArtistCredit[0].Name
		}

		recordings = append(recordings, Recording{
			ID:     r.ID,
			Title:  r.Title,
			Artist: artistName,
			Score:  r.Score,
		})
	}

	return recordings, nil
}

// GetBestArtistMatch returns the best matching artist (highest score) or nil if no matches
func GetBestArtistMatch(artists []Artist) *Artist {
	if len(artists) == 0 {
		return nil
	}
	// The results are already sorted by score, so just return the first one
	return &artists[0]
}

// GetBestReleaseMatch returns the best matching release (highest score) or nil if no matches
func GetBestReleaseMatch(releases []Release) *Release {
	if len(releases) == 0 {
		return nil
	}
	return &releases[0]
}

// GetBestRecordingMatch returns the best matching recording (highest score) or nil if no matches
func GetBestRecordingMatch(recordings []Recording) *Recording {
	if len(recordings) == 0 {
		return nil
	}
	return &recordings[0]
}
