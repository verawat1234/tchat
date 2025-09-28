package main

import (
	"log"

	"tchat.dev/shared/app"
	"tchat.dev/shared/config"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Override port for video service if not set
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8091 // Video service port
	}

	// Create service implementation
	videoService := NewVideoServiceImpl()

	// Create base application with video service implementation
	baseApp := app.NewBaseApp(cfg, videoService)

	// Initialize application
	if err := baseApp.Initialize(); err != nil {
		log.Fatalf("Failed to initialize video service: %v", err)
	}

	// Run with graceful shutdown
	baseApp.RunWithGracefulShutdown()
}