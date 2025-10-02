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

	"github.com/gin-gonic/gin"
	"tchat.dev/shared/config"
)

// SimpleVideoApp provides basic health endpoint without full service initialization
type SimpleVideoApp struct {
	config *config.Config
	server *http.Server
}

// NewSimpleVideoApp creates a new minimal video service
func NewSimpleVideoApp(cfg *config.Config) *SimpleVideoApp {
	return &SimpleVideoApp{
		config: cfg,
	}
}

// Initialize sets up the minimal HTTP server
func (a *SimpleVideoApp) Initialize() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "video-service",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
		})
	})

	// API health endpoint
	router.GET("/api/v1/videos/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "video service healthy",
		})
	})

	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
	a.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("Video service initialized on %s", addr)
	return nil
}

// Run starts the server
func (a *SimpleVideoApp) Run() error {
	log.Printf("Starting video service on %s", a.server.Addr)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Video service is running on %s", a.server.Addr)
	return nil
}

// Shutdown gracefully shuts down the server
func (a *SimpleVideoApp) Shutdown(ctx context.Context) error {
	log.Println("Shutting down video service...")
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	log.Println("Video service shutdown completed")
	return nil
}

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Force video service port
	cfg.Server.Port = 8091

	// Create and initialize simple app
	app := NewSimpleVideoApp(cfg)
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the application
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown gracefully: %v", err)
		os.Exit(1)
	}

	log.Println("Video service stopped")
}
