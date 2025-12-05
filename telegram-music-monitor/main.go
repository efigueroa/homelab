package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/efigueroa/telegram-music-monitor/api"
	"github.com/efigueroa/telegram-music-monitor/config"
	"github.com/efigueroa/telegram-music-monitor/database"
	"github.com/efigueroa/telegram-music-monitor/lidarr"
	"github.com/efigueroa/telegram-music-monitor/musicbrainz"
	"github.com/efigueroa/telegram-music-monitor/platforms"
	"github.com/efigueroa/telegram-music-monitor/telegram"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Telegram Music Monitor...")

	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database initialized successfully")

	// Load configuration overrides from database
	// This allows runtime updates via the web UI to persist
	overrides, err := db.GetAllConfigOverrides()
	if err != nil {
		log.Printf("Warning: Failed to load config overrides: %v", err)
	} else if len(overrides) > 0 {
		cfg.ApplyOverrides(overrides)
		log.Printf("Applied %d configuration overrides from database", len(overrides))
	}

	// Initialize API clients
	spotifyClient := platforms.NewSpotifyClient(cfg.SpotifyClientID, cfg.SpotifyClientSecret)
	log.Println("Spotify client initialized")

	youtubeClient := platforms.NewYouTubeClient(cfg.YouTubeAPIKey)
	log.Println("YouTube client initialized")

	mbClient := musicbrainz.NewClient()
	log.Println("MusicBrainz client initialized")

	lidarrClient := lidarr.NewClient(cfg.LidarrURL, cfg.LidarrAPIKey)
	log.Println("Lidarr client initialized")

	// Initialize API server with database and config
	apiServer := api.NewServer(db, cfg)
	log.Printf("API server initialized on port %s", cfg.ServerPort)

	// Create context that listens for shutdown signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	httpServer := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: apiServer.Router(),
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.ServerPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Initialize and start Telegram client in a goroutine
	telegramClient := telegram.NewClient(telegram.Config{
		APIID:         cfg.TelegramAPIID,
		APIHash:       cfg.TelegramAPIHash,
		Phone:         cfg.TelegramPhone,
		SessionPath:   "/data/telegram.session", // Store session in data volume
		DB:            db,
		SpotifyClient: spotifyClient,
		YouTubeClient: youtubeClient,
		MBClient:      mbClient,
		LidarrClient:  lidarrClient,
	})

	// Run Telegram client in a goroutine
	telegramDone := make(chan error, 1)
	go func() {
		log.Println("Starting Telegram client...")
		telegramDone <- telegramClient.Run(ctx)
	}()

	// Wait for shutdown signal or Telegram client error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		log.Println("Shutting down gracefully...")

		// Cancel context to stop Telegram client
		cancel()

		// Shutdown HTTP server with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		} else {
			log.Println("HTTP server stopped")
		}

		// Wait for Telegram client to stop
		select {
		case <-telegramDone:
			log.Println("Telegram client stopped")
		case <-time.After(10 * time.Second):
			log.Println("Telegram client shutdown timeout")
		}

	case err := <-telegramDone:
		// Telegram client exited
		if err != nil {
			log.Printf("Telegram client error: %v", err)
		}
		log.Println("Telegram client stopped")

		// Shutdown HTTP server
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}

	log.Println("Application stopped")
}
