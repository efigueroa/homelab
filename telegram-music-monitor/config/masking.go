package config

import (
	"strings"
)

// MaskAPIKey masks an API key or secret for display purposes
// Shows only the first 4 and last 4 characters, with asterisks in between
// SECURITY: Never display full API keys to users
func MaskAPIKey(key string) string {
	if key == "" {
		return ""
	}

	// If key is too short, just mask it entirely
	if len(key) <= 8 {
		return "****"
	}

	// Show first 4 and last 4 characters
	prefix := key[:4]
	suffix := key[len(key)-4:]
	middle := strings.Repeat("*", 12) // Fixed length for consistency

	return prefix + middle + suffix
}

// MaskPhoneNumber masks a phone number for display
// Shows only the last 4 digits
func MaskPhoneNumber(phone string) string {
	if phone == "" {
		return ""
	}

	if len(phone) <= 4 {
		return "****"
	}

	// Show only last 4 digits
	suffix := phone[len(phone)-4:]
	return "****" + suffix
}

// MaskURL masks a URL while keeping the protocol and domain visible
// Example: http://example.com/path -> http://example.com/****
func MaskURL(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	// Find the last slash
	lastSlash := strings.LastIndex(urlStr, "/")
	if lastSlash == -1 {
		return urlStr // No path to mask
	}

	// Check if this is just the protocol slashes
	if lastSlash < len(urlStr)-1 {
		return urlStr[:lastSlash+1] + "****"
	}

	return urlStr
}

// ValidConfigKey checks if a configuration key is valid and can be updated
func ValidConfigKey(key string) bool {
	validKeys := []string{
		"SPOTIFY_CLIENT_ID",
		"SPOTIFY_CLIENT_SECRET",
		"YOUTUBE_API_KEY",
		"LIDARR_URL",
		"LIDARR_API_KEY",
	}

	for _, valid := range validKeys {
		if key == valid {
			return true
		}
	}

	return false
}

// GetMaskedConfigValue returns a masked version of a config value based on its key
func GetMaskedConfigValue(key, value string) string {
	if value == "" {
		return ""
	}

	switch key {
	case "TELEGRAM_PHONE":
		return MaskPhoneNumber(value)
	case "LIDARR_URL":
		return MaskURL(value)
	case "SPOTIFY_CLIENT_ID", "SPOTIFY_CLIENT_SECRET", "YOUTUBE_API_KEY", "LIDARR_API_KEY", "TELEGRAM_API_HASH":
		return MaskAPIKey(value)
	case "TELEGRAM_API_ID":
		// API ID is numeric, just show it's set
		return "****"
	default:
		return MaskAPIKey(value)
	}
}
