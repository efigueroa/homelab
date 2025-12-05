package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite" // Pure Go SQLite driver (no CGO required)
)

// DB wraps the sql.DB connection and provides methods for our application
type DB struct {
	conn *sql.DB
}

// DetectedMusic represents a music link detected in Telegram
type DetectedMusic struct {
	ID              int       // Auto-incrementing primary key
	MessageID       int64     // Telegram message ID
	GroupID         int64     // Telegram group/chat ID
	GroupName       string    // Human-readable group name
	SenderID        int64     // Telegram user ID who sent the message
	SenderUsername  string    // Username of the sender
	OriginalLink    string    // The original URL shared
	Platform        string    // Platform (spotify, youtube, etc.)
	DetectedAt      time.Time // When the link was detected
	SongName        string    // Detected song name
	AlbumName       string    // Detected album name
	ArtistName      string    // Detected artist name
	MusicBrainzID   string    // MusicBrainz ID for the item
	LidarrStatus    string    // Status: pending, added, failed, skipped
	LidarrAddedAt   *time.Time // When it was added to Lidarr (nil if not added)
	LidarrErrorMsg  string    // Error message if Lidarr add failed
}

// GroupConfig stores per-group configuration
type GroupConfig struct {
	GroupID   int64  // Telegram group/chat ID
	GroupName string // Human-readable group name
	Enabled   bool   // Whether monitoring is enabled for this group
	AutoAdd   bool   // Whether to auto-add to Lidarr (always true in this version)
}

// ConfigOverride stores runtime configuration overrides
// These values override environment variables and can be updated via the web UI
type ConfigOverride struct {
	Key       string    // Configuration key name
	Value     string    // Configuration value (encrypted/sensitive)
	UpdatedAt time.Time // When this was last updated
}

// New creates a new database connection and initializes the schema
func New(dbPath string) (*DB, error) {
	// Open SQLite database
	// The ?_pragma parameters optimize SQLite for our use case
	conn, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Initialize the database schema
	if err := db.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// initSchema creates all necessary tables if they don't exist
func (db *DB) initSchema() error {
	schema := `
	-- Table to store detected music links and metadata
	CREATE TABLE IF NOT EXISTS detected_music (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message_id INTEGER NOT NULL,
		group_id INTEGER NOT NULL,
		group_name TEXT NOT NULL,
		sender_id INTEGER NOT NULL,
		sender_username TEXT,
		original_link TEXT NOT NULL,
		platform TEXT NOT NULL,
		detected_at TIMESTAMP NOT NULL,
		song_name TEXT,
		album_name TEXT,
		artist_name TEXT,
		musicbrainz_id TEXT,
		lidarr_status TEXT NOT NULL DEFAULT 'pending',
		lidarr_added_at TIMESTAMP,
		lidarr_error_msg TEXT
	);

	-- Table to store group-specific configuration
	CREATE TABLE IF NOT EXISTS group_configs (
		group_id INTEGER PRIMARY KEY,
		group_name TEXT NOT NULL,
		enabled BOOLEAN NOT NULL DEFAULT 1,
		auto_add BOOLEAN NOT NULL DEFAULT 1
	);

	-- Table to store runtime configuration overrides
	-- This allows updating API keys and other settings via the web UI
	CREATE TABLE IF NOT EXISTS config_overrides (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	-- Indexes for faster queries
	CREATE INDEX IF NOT EXISTS idx_detected_music_group_id ON detected_music(group_id);
	CREATE INDEX IF NOT EXISTS idx_detected_music_artist ON detected_music(artist_name);
	CREATE INDEX IF NOT EXISTS idx_detected_music_platform ON detected_music(platform);
	CREATE INDEX IF NOT EXISTS idx_detected_music_detected_at ON detected_music(detected_at);
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// InsertDetectedMusic adds a new detected music entry to the database
func (db *DB) InsertDetectedMusic(music *DetectedMusic) (int64, error) {
	query := `
		INSERT INTO detected_music (
			message_id, group_id, group_name, sender_id, sender_username,
			original_link, platform, detected_at, song_name, album_name,
			artist_name, musicbrainz_id, lidarr_status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query,
		music.MessageID, music.GroupID, music.GroupName, music.SenderID,
		music.SenderUsername, music.OriginalLink, music.Platform,
		music.DetectedAt, music.SongName, music.AlbumName, music.ArtistName,
		music.MusicBrainzID, music.LidarrStatus,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert detected music: %w", err)
	}

	// Return the ID of the newly inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// UpdateLidarrStatus updates the Lidarr status for a detected music entry
func (db *DB) UpdateLidarrStatus(id int64, status string, errorMsg string) error {
	query := `
		UPDATE detected_music
		SET lidarr_status = ?, lidarr_added_at = ?, lidarr_error_msg = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := db.conn.Exec(query, status, now, errorMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update lidarr status: %w", err)
	}

	return nil
}

// GetGroupConfig retrieves the configuration for a specific group
func (db *DB) GetGroupConfig(groupID int64) (*GroupConfig, error) {
	query := `SELECT group_id, group_name, enabled, auto_add FROM group_configs WHERE group_id = ?`

	var config GroupConfig
	err := db.conn.QueryRow(query, groupID).Scan(
		&config.GroupID, &config.GroupName, &config.Enabled, &config.AutoAdd,
	)

	if err == sql.ErrNoRows {
		// If no config exists, return default config (enabled, auto-add)
		return &GroupConfig{
			GroupID: groupID,
			Enabled: true,
			AutoAdd: true,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get group config: %w", err)
	}

	return &config, nil
}

// UpsertGroupConfig inserts or updates a group configuration
func (db *DB) UpsertGroupConfig(config *GroupConfig) error {
	query := `
		INSERT INTO group_configs (group_id, group_name, enabled, auto_add)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(group_id) DO UPDATE SET
			group_name = excluded.group_name,
			enabled = excluded.enabled,
			auto_add = excluded.auto_add
	`

	_, err := db.conn.Exec(query, config.GroupID, config.GroupName, config.Enabled, config.AutoAdd)
	if err != nil {
		return fmt.Errorf("failed to upsert group config: %w", err)
	}

	return nil
}

// GetRecentMusic retrieves the most recently detected music
func (db *DB) GetRecentMusic(limit int) ([]DetectedMusic, error) {
	query := `
		SELECT id, message_id, group_id, group_name, sender_id, sender_username,
		       original_link, platform, detected_at, song_name, album_name,
		       artist_name, musicbrainz_id, lidarr_status, lidarr_added_at, lidarr_error_msg
		FROM detected_music
		ORDER BY detected_at DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent music: %w", err)
	}
	defer rows.Close()

	var results []DetectedMusic
	for rows.Next() {
		var music DetectedMusic
		err := rows.Scan(
			&music.ID, &music.MessageID, &music.GroupID, &music.GroupName,
			&music.SenderID, &music.SenderUsername, &music.OriginalLink,
			&music.Platform, &music.DetectedAt, &music.SongName, &music.AlbumName,
			&music.ArtistName, &music.MusicBrainzID, &music.LidarrStatus,
			&music.LidarrAddedAt, &music.LidarrErrorMsg,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, music)
	}

	return results, nil
}

// GetTopArtists retrieves the most frequently detected artists
func (db *DB) GetTopArtists(limit int) ([]struct{ Artist string; Count int }, error) {
	query := `
		SELECT artist_name, COUNT(*) as count
		FROM detected_music
		WHERE artist_name IS NOT NULL AND artist_name != ''
		GROUP BY artist_name
		ORDER BY count DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top artists: %w", err)
	}
	defer rows.Close()

	var results []struct{ Artist string; Count int }
	for rows.Next() {
		var item struct{ Artist string; Count int }
		if err := rows.Scan(&item.Artist, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

// GetTopAlbums retrieves the most frequently detected albums
func (db *DB) GetTopAlbums(limit int) ([]struct{ Album string; Artist string; Count int }, error) {
	query := `
		SELECT album_name, artist_name, COUNT(*) as count
		FROM detected_music
		WHERE album_name IS NOT NULL AND album_name != ''
		GROUP BY album_name, artist_name
		ORDER BY count DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top albums: %w", err)
	}
	defer rows.Close()

	var results []struct{ Album string; Artist string; Count int }
	for rows.Next() {
		var item struct{ Album string; Artist string; Count int }
		if err := rows.Scan(&item.Album, &item.Artist, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

// GetGroupActivity retrieves activity statistics by group
func (db *DB) GetGroupActivity() ([]struct{ GroupName string; Count int }, error) {
	query := `
		SELECT group_name, COUNT(*) as count
		FROM detected_music
		GROUP BY group_name
		ORDER BY count DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query group activity: %w", err)
	}
	defer rows.Close()

	var results []struct{ GroupName string; Count int }
	for rows.Next() {
		var item struct{ GroupName string; Count int }
		if err := rows.Scan(&item.GroupName, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

// GetPlatformStats retrieves statistics by platform
func (db *DB) GetPlatformStats() ([]struct{ Platform string; Count int }, error) {
	query := `
		SELECT platform, COUNT(*) as count
		FROM detected_music
		GROUP BY platform
		ORDER BY count DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform stats: %w", err)
	}
	defer rows.Close()

	var results []struct{ Platform string; Count int }
	for rows.Next() {
		var item struct{ Platform string; Count int }
		if err := rows.Scan(&item.Platform, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

// SetConfigOverride stores or updates a configuration override
// This allows runtime updates to configuration values via the web UI
func (db *DB) SetConfigOverride(key, value string) error {
	query := `
		INSERT INTO config_overrides (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`

	_, err := db.conn.Exec(query, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to set config override: %w", err)
	}

	return nil
}

// GetConfigOverride retrieves a configuration override value
// Returns empty string if the key doesn't exist
func (db *DB) GetConfigOverride(key string) (string, error) {
	query := `SELECT value FROM config_overrides WHERE key = ?`

	var value string
	err := db.conn.QueryRow(query, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil // Key doesn't exist, return empty string
	}
	if err != nil {
		return "", fmt.Errorf("failed to get config override: %w", err)
	}

	return value, nil
}

// GetAllConfigOverrides retrieves all configuration overrides
// Returns a map of key-value pairs
func (db *DB) GetAllConfigOverrides() (map[string]string, error) {
	query := `SELECT key, value FROM config_overrides`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query config overrides: %w", err)
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result[key] = value
	}

	return result, nil
}

// DeleteConfigOverride removes a configuration override
func (db *DB) DeleteConfigOverride(key string) error {
	query := `DELETE FROM config_overrides WHERE key = ?`

	_, err := db.conn.Exec(query, key)
	if err != nil {
		return fmt.Errorf("failed to delete config override: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
