package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"tchat.dev/calling/config"
	"tchat.dev/calling/handlers"
	"tchat.dev/calling/repositories"
	"tchat.dev/calling/services"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize database (already connects in constructor)
	dbManager, err := config.NewDatabaseManager()
	if err != nil {
		log.Fatalf("Failed to create database manager: %v", err)
	}
	defer dbManager.Close()

	// Initialize Redis (already connects in constructor)
	redisClient, err := config.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}
	defer redisClient.Close()

	// CORS middleware for WebRTC
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Initialize repositories
	callSessionRepo := repositories.NewGormCallSessionRepository(dbManager.GetDB())
	userPresenceRepo := repositories.NewRedisUserPresenceRepository(redisClient.Client, 24*time.Hour)
	callHistoryRepo := repositories.NewGormCallHistoryRepository(dbManager.GetDB())

	// Initialize services
	callService := services.NewCallService(callSessionRepo, userPresenceRepo, callHistoryRepo)
	presenceService := services.NewPresenceService(userPresenceRepo)
	signalingService := services.NewSignalingService(callService, presenceService)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(dbManager.GetDB(), redisClient)
	callHandler := handlers.NewCallHandlers(callService)
	presenceHandler := handlers.NewPresenceHandlers(presenceService)
	historyHandler := handlers.NewHistoryHandlers(callHistoryRepo)
	websocketHandler := handlers.NewWebSocketHandler(signalingService)

	// Health check endpoints
	r.GET("/health", healthHandler.BasicHealthCheck)
	r.GET("/health/detailed", healthHandler.DetailedHealthCheck)
	r.GET("/ready", healthHandler.ReadinessCheck)
	r.GET("/live", healthHandler.LivenessCheck)
	r.GET("/metrics", healthHandler.MetricsEndpoint)

	// API v1 group
	api := r.Group("/api/v1")

	// Register handler routes
	callHandler.RegisterRoutes(api)
	presenceHandler.RegisterRoutes(api)
	historyHandler.RegisterRoutes(api)

	// WebSocket signaling endpoint
	api.GET("/signaling", websocketHandler.HandleWebSocket)

	// Start server
	srv := &http.Server{
		Addr:    ":8093",
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("Calling service started on :8093")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
