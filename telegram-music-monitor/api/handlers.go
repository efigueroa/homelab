package api

import (
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/efigueroa/telegram-music-monitor/config"
	"github.com/efigueroa/telegram-music-monitor/database"
	"github.com/gorilla/mux"
)

// Server wraps the HTTP server and database
type Server struct {
	db      *database.DB
	config  *config.Config
	router  *mux.Router
	configMu sync.RWMutex // Protects config updates
}

//go:embed dashboard.html settings.html
var dashboardHTML embed.FS

// NewServer creates a new API server
func NewServer(db *database.DB, cfg *config.Config) *Server {
	s := &Server{
		db:     db,
		config: cfg,
		router: mux.NewRouter(),
	}

	// Set up routes
	s.setupRoutes()

	return s
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Dashboard (serves the embedded HTML file)
	s.router.HandleFunc("/", s.handleDashboard).Methods("GET")
	s.router.HandleFunc("/settings", s.handleSettingsPage).Methods("GET")

	// API endpoints
	api := s.router.PathPrefix("/api").Subrouter()

	// Metrics endpoints
	api.HandleFunc("/metrics/recent", s.handleRecentMusic).Methods("GET")
	api.HandleFunc("/metrics/top-artists", s.handleTopArtists).Methods("GET")
	api.HandleFunc("/metrics/top-albums", s.handleTopAlbums).Methods("GET")
	api.HandleFunc("/metrics/group-activity", s.handleGroupActivity).Methods("GET")
	api.HandleFunc("/metrics/platform-stats", s.handlePlatformStats).Methods("GET")

	// Settings/Config endpoints
	api.HandleFunc("/settings", s.handleGetSettings).Methods("GET")
	api.HandleFunc("/settings", s.handleUpdateSettings).Methods("POST")

	// Health check
	api.HandleFunc("/health", s.handleHealth).Methods("GET")
}

// Router returns the mux router
func (s *Server) Router() *mux.Router {
	return s.router
}

// handleDashboard serves the embedded HTML dashboard
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	// Read the embedded dashboard HTML file
	data, err := dashboardHTML.ReadFile("dashboard.html")
	if err != nil {
		http.Error(w, "Failed to load dashboard", http.StatusInternalServerError)
		log.Printf("Error reading dashboard.html: %v", err)
		return
	}

	// Set content type and write the HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// handleRecentMusic returns recently detected music
// Query parameter: limit (default: 50)
func (s *Server) handleRecentMusic(w http.ResponseWriter, r *http.Request) {
	// Get limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Query database
	music, err := s.db.GetRecentMusic(limit)
	if err != nil {
		http.Error(w, "Failed to fetch recent music", http.StatusInternalServerError)
		log.Printf("Error fetching recent music: %v", err)
		return
	}

	// Return JSON response
	respondJSON(w, music)
}

// handleTopArtists returns the most frequently detected artists
// Query parameter: limit (default: 10)
func (s *Server) handleTopArtists(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	artists, err := s.db.GetTopArtists(limit)
	if err != nil {
		http.Error(w, "Failed to fetch top artists", http.StatusInternalServerError)
		log.Printf("Error fetching top artists: %v", err)
		return
	}

	respondJSON(w, artists)
}

// handleTopAlbums returns the most frequently detected albums
// Query parameter: limit (default: 10)
func (s *Server) handleTopAlbums(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	albums, err := s.db.GetTopAlbums(limit)
	if err != nil {
		http.Error(w, "Failed to fetch top albums", http.StatusInternalServerError)
		log.Printf("Error fetching top albums: %v", err)
		return
	}

	respondJSON(w, albums)
}

// handleGroupActivity returns activity statistics by group
func (s *Server) handleGroupActivity(w http.ResponseWriter, r *http.Request) {
	activity, err := s.db.GetGroupActivity()
	if err != nil {
		http.Error(w, "Failed to fetch group activity", http.StatusInternalServerError)
		log.Printf("Error fetching group activity: %v", err)
		return
	}

	respondJSON(w, activity)
}

// handlePlatformStats returns statistics by platform
func (s *Server) handlePlatformStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db.GetPlatformStats()
	if err != nil {
		http.Error(w, "Failed to fetch platform stats", http.StatusInternalServerError)
		log.Printf("Error fetching platform stats: %v", err)
		return
	}

	respondJSON(w, stats)
}

// handleHealth returns a simple health check response
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{
		"status": "ok",
	})
}

// handleSettingsPage serves the settings HTML page
func (s *Server) handleSettingsPage(w http.ResponseWriter, r *http.Request) {
	data, err := dashboardHTML.ReadFile("settings.html")
	if err != nil {
		http.Error(w, "Failed to load settings page", http.StatusInternalServerError)
		log.Printf("Error reading settings.html: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// handleGetSettings returns the current configuration (with masked sensitive values)
// SECURITY: This endpoint never returns full API keys, only masked versions
func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	s.configMu.RLock()
	defer s.configMu.RUnlock()

	// Get masked configuration values
	maskedConfig := s.config.GetMaskedConfig()

	// Add which fields can be updated
	response := map[string]interface{}{
		"config": maskedConfig,
		"updatableFields": []string{
			"SPOTIFY_CLIENT_ID",
			"SPOTIFY_CLIENT_SECRET",
			"YOUTUBE_API_KEY",
			"LIDARR_URL",
			"LIDARR_API_KEY",
		},
	}

	respondJSON(w, response)
}

// handleUpdateSettings updates configuration values
// Request body should be JSON: {"key": "SPOTIFY_CLIENT_ID", "value": "new_value"}
func (s *Server) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the key is updatable
	if !config.ValidConfigKey(req.Key) {
		http.Error(w, "Invalid configuration key", http.StatusBadRequest)
		return
	}

	// Validate the value is not empty
	if req.Value == "" {
		http.Error(w, "Configuration value cannot be empty", http.StatusBadRequest)
		return
	}

	// Save to database
	if err := s.db.SetConfigOverride(req.Key, req.Value); err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		log.Printf("Error saving config override: %v", err)
		return
	}

	// Apply the override to the current config
	s.configMu.Lock()
	s.config.ApplyOverrides(map[string]string{req.Key: req.Value})
	s.configMu.Unlock()

	log.Printf("Configuration updated: %s", req.Key)

	// Return success with masked value
	respondJSON(w, map[string]string{
		"status":  "success",
		"key":     req.Key,
		"message": "Configuration updated successfully",
		"masked":  config.GetMaskedConfigValue(req.Key, req.Value),
	})
}

// respondJSON is a helper function to send JSON responses
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
