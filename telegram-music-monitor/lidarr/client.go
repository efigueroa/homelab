package lidarr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client handles Lidarr API interactions
// Lidarr is a music collection manager that can automatically download music
type Client struct {
	baseURL    string // Lidarr instance URL (e.g., http://localhost:8686)
	apiKey     string // Lidarr API key
	httpClient *http.Client
}

// Artist represents a Lidarr artist
type Artist struct {
	ID              int    `json:"id"`
	ArtistName      string `json:"artistName"`
	ForeignArtistID string `json:"foreignArtistId"` // MusicBrainz ID
	Monitored       bool   `json:"monitored"`
}

// Album represents a Lidarr album
type Album struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	ForeignAlbumID string `json:"foreignAlbumId"` // MusicBrainz ID
	Monitored     bool   `json:"monitored"`
}

// SearchResult represents a search result from Lidarr's lookup endpoint
type SearchResult struct {
	ForeignArtistID string `json:"foreignArtistId"` // MusicBrainz ID
	ArtistName      string `json:"artistName"`
}

// AddArtistRequest represents the request body for adding an artist to Lidarr
type AddArtistRequest struct {
	ArtistName             string         `json:"artistName"`
	ForeignArtistID        string         `json:"foreignArtistId"` // MusicBrainz ID
	Monitored              bool           `json:"monitored"`
	QualityProfileID       int            `json:"qualityProfileId"`
	MetadataProfileID      int            `json:"metadataProfileId"`
	RootFolderPath         string         `json:"rootFolderPath"`
	AddOptions             AddOptions     `json:"addOptions"`
	MonitorNewItems        string         `json:"monitorNewItems"`
}

// AddOptions contains options for when adding an artist
type AddOptions struct {
	SearchForMissingAlbums bool `json:"searchForMissingAlbums"`
}

// RootFolder represents a root folder in Lidarr where music is stored
type RootFolder struct {
	Path string `json:"path"`
	ID   int    `json:"id"`
}

// QualityProfile represents a quality profile in Lidarr
type QualityProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MetadataProfile represents a metadata profile in Lidarr
type MetadataProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// NewClient creates a new Lidarr API client
func NewClient(baseURL, apiKey string) *Client {
	// Remove trailing slash from base URL if present
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// doRequest performs an HTTP request to the Lidarr API
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s/api/v1/%s", c.baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Lidarr requires the API key in a custom header
	req.Header.Set("X-Api-Key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// LookupArtist searches for an artist by name or MusicBrainz ID
func (c *Client) LookupArtist(ctx context.Context, term string) ([]SearchResult, error) {
	// The lookup endpoint can search by name or MusicBrainz ID
	endpoint := fmt.Sprintf("search?term=%s", term)

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lidarr API returned status %d: %s", resp.StatusCode, string(body))
	}

	var results []SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return results, nil
}

// GetRootFolders retrieves all root folders configured in Lidarr
func (c *Client) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
	resp, err := c.doRequest(ctx, "GET", "rootfolder", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lidarr API returned status %d: %s", resp.StatusCode, string(body))
	}

	var folders []RootFolder
	if err := json.NewDecoder(resp.Body).Decode(&folders); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return folders, nil
}

// GetQualityProfiles retrieves all quality profiles configured in Lidarr
func (c *Client) GetQualityProfiles(ctx context.Context) ([]QualityProfile, error) {
	resp, err := c.doRequest(ctx, "GET", "qualityprofile", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lidarr API returned status %d: %s", resp.StatusCode, string(body))
	}

	var profiles []QualityProfile
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return profiles, nil
}

// GetMetadataProfiles retrieves all metadata profiles configured in Lidarr
func (c *Client) GetMetadataProfiles(ctx context.Context) ([]MetadataProfile, error) {
	resp, err := c.doRequest(ctx, "GET", "metadataprofile", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lidarr API returned status %d: %s", resp.StatusCode, string(body))
	}

	var profiles []MetadataProfile
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return profiles, nil
}

// AddArtist adds an artist to Lidarr
// It requires a MusicBrainz ID and will use the first available root folder and profiles
func (c *Client) AddArtist(ctx context.Context, artistName, musicbrainzID string) error {
	// Get necessary configuration
	rootFolders, err := c.GetRootFolders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get root folders: %w", err)
	}
	if len(rootFolders) == 0 {
		return fmt.Errorf("no root folders configured in Lidarr")
	}

	qualityProfiles, err := c.GetQualityProfiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get quality profiles: %w", err)
	}
	if len(qualityProfiles) == 0 {
		return fmt.Errorf("no quality profiles configured in Lidarr")
	}

	metadataProfiles, err := c.GetMetadataProfiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get metadata profiles: %w", err)
	}
	if len(metadataProfiles) == 0 {
		return fmt.Errorf("no metadata profiles configured in Lidarr")
	}

	// Build the add artist request
	// Use the first root folder and profiles (you could make this configurable)
	addReq := AddArtistRequest{
		ArtistName:        artistName,
		ForeignArtistID:   musicbrainzID,
		Monitored:         true, // Monitor the artist
		QualityProfileID:  qualityProfiles[0].ID,
		MetadataProfileID: metadataProfiles[0].ID,
		RootFolderPath:    rootFolders[0].Path,
		MonitorNewItems:   "all", // Monitor all new items
		AddOptions: AddOptions{
			SearchForMissingAlbums: true, // Automatically search for albums
		},
	}

	resp, err := c.doRequest(ctx, "POST", "artist", addReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Lidarr returns 201 Created on success
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lidarr API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetArtist retrieves an artist from Lidarr by MusicBrainz ID
// Returns nil if the artist is not found
func (c *Client) GetArtist(ctx context.Context, musicbrainzID string) (*Artist, error) {
	endpoint := fmt.Sprintf("artist?mbId=%s", musicbrainzID)

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Artist not found
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("lidarr API returned status %d: %s", resp.StatusCode, string(body))
	}

	// The endpoint returns an array, even for a single artist
	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(artists) == 0 {
		return nil, nil
	}

	return &artists[0], nil
}

// ArtistExists checks if an artist already exists in Lidarr
func (c *Client) ArtistExists(ctx context.Context, musicbrainzID string) (bool, error) {
	artist, err := c.GetArtist(ctx, musicbrainzID)
	if err != nil {
		return false, err
	}
	return artist != nil, nil
}
