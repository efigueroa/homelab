package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration loaded from environment variables
// This struct contains all the credentials and settings needed to run the application
type Config struct {
	// Telegram configuration
	TelegramAPIID   int    // Telegram API ID from my.telegram.org
	TelegramAPIHash string // Telegram API Hash from my.telegram.org
	TelegramPhone   string // Phone number in international format (e.g., +1234567890)

	// Spotify configuration
	SpotifyClientID     string // Spotify Client ID from developer.spotify.com
	SpotifyClientSecret string // Spotify Client Secret from developer.spotify.com

	// YouTube configuration
	YouTubeAPIKey string // YouTube Data API key from console.cloud.google.com

	// Lidarr configuration
	LidarrURL    string // Lidarr instance URL (e.g., http://localhost:8686)
	LidarrAPIKey string // Lidarr API key from Settings > General

	// Database configuration
	DatabasePath string // Path to SQLite database file

	// API Server configuration
	ServerPort string // Port for the API server to listen on (default: 8080)
}

// Load reads configuration from environment variables
// It returns an error if any required environment variable is missing
func Load() (*Config, error) {
	cfg := &Config{}

	// Load Telegram configuration
	// API ID must be converted from string to int
	apiIDStr := os.Getenv("TELEGRAM_API_ID")
	if apiIDStr == "" {
		return nil, fmt.Errorf("TELEGRAM_API_ID environment variable is required")
	}
	apiID, err := strconv.Atoi(apiIDStr)
	if err != nil {
		return nil, fmt.Errorf("TELEGRAM_API_ID must be a valid integer: %w", err)
	}
	cfg.TelegramAPIID = apiID

	cfg.TelegramAPIHash = os.Getenv("TELEGRAM_API_HASH")
	if cfg.TelegramAPIHash == "" {
		return nil, fmt.Errorf("TELEGRAM_API_HASH environment variable is required")
	}

	cfg.TelegramPhone = os.Getenv("TELEGRAM_PHONE")
	if cfg.TelegramPhone == "" {
		return nil, fmt.Errorf("TELEGRAM_PHONE environment variable is required")
	}

	// Load Spotify configuration
	cfg.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	if cfg.SpotifyClientID == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_ID environment variable is required")
	}

	cfg.SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	if cfg.SpotifyClientSecret == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_SECRET environment variable is required")
	}

	// Load YouTube configuration
	cfg.YouTubeAPIKey = os.Getenv("YOUTUBE_API_KEY")
	if cfg.YouTubeAPIKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY environment variable is required")
	}

	// Load Lidarr configuration
	cfg.LidarrURL = os.Getenv("LIDARR_URL")
	if cfg.LidarrURL == "" {
		return nil, fmt.Errorf("LIDARR_URL environment variable is required")
	}

	cfg.LidarrAPIKey = os.Getenv("LIDARR_API_KEY")
	if cfg.LidarrAPIKey == "" {
		return nil, fmt.Errorf("LIDARR_API_KEY environment variable is required")
	}

	// Load database configuration with default
	cfg.DatabasePath = os.Getenv("DATABASE_PATH")
	if cfg.DatabasePath == "" {
		cfg.DatabasePath = "/data/music_monitor.db" // Default path for Docker volume
	}

	// Load server configuration with default
	cfg.ServerPort = os.Getenv("SERVER_PORT")
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080" // Default port
	}

	return cfg, nil
}

// ApplyOverrides applies database configuration overrides to the config
// This allows runtime updates to configuration values
func (c *Config) ApplyOverrides(overrides map[string]string) {
	// Apply each override if it exists
	if val, ok := overrides["SPOTIFY_CLIENT_ID"]; ok && val != "" {
		c.SpotifyClientID = val
	}
	if val, ok := overrides["SPOTIFY_CLIENT_SECRET"]; ok && val != "" {
		c.SpotifyClientSecret = val
	}
	if val, ok := overrides["YOUTUBE_API_KEY"]; ok && val != "" {
		c.YouTubeAPIKey = val
	}
	if val, ok := overrides["LIDARR_URL"]; ok && val != "" {
		c.LidarrURL = val
	}
	if val, ok := overrides["LIDARR_API_KEY"]; ok && val != "" {
		c.LidarrAPIKey = val
	}
}

// GetMaskedConfig returns a map of configuration keys to their masked values
// SECURITY: This should be used when displaying config to users
func (c *Config) GetMaskedConfig() map[string]string {
	return map[string]string{
		"TELEGRAM_API_ID":       MaskAPIKey(strconv.Itoa(c.TelegramAPIID)),
		"TELEGRAM_API_HASH":     MaskAPIKey(c.TelegramAPIHash),
		"TELEGRAM_PHONE":        MaskPhoneNumber(c.TelegramPhone),
		"SPOTIFY_CLIENT_ID":     MaskAPIKey(c.SpotifyClientID),
		"SPOTIFY_CLIENT_SECRET": MaskAPIKey(c.SpotifyClientSecret),
		"YOUTUBE_API_KEY":       MaskAPIKey(c.YouTubeAPIKey),
		"LIDARR_URL":            MaskURL(c.LidarrURL),
		"LIDARR_API_KEY":        MaskAPIKey(c.LidarrAPIKey),
		"DATABASE_PATH":         c.DatabasePath,
		"SERVER_PORT":           c.ServerPort,
	}
}
