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
)

func main() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("ENVIRONMENT") == "development" {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database and Redis connections
	dbManager, redisClient, err := config.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize connections: %v", err)
	}
	defer dbManager.Close()
	defer redisClient.Close()

	// Create Gin router
	r := gin.Default()

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

	// Initialize health handler
	healthHandler := handlers.NewHealthHandler(dbManager.GetDB(), redisClient)

	// Health check endpoints
	r.GET("/health", healthHandler.BasicHealthCheck)
	r.GET("/health/detailed", healthHandler.DetailedHealthCheck)
	r.GET("/ready", healthHandler.ReadinessCheck)
	r.GET("/live", healthHandler.LivenessCheck)
	r.GET("/metrics", healthHandler.MetricsEndpoint)

	// Initialize handlers with dependencies
	callHandler := handlers.NewCallHandler(dbManager.GetDB(), redisClient)
	presenceHandler := handlers.NewPresenceHandler(redisClient)
	historyHandler := handlers.NewHistoryHandler(dbManager.GetDB())
	websocketHandler := handlers.NewWebSocketHandler(redisClient)

	// API v1 group
	api := r.Group("/api/v1")
	{
		// Call management endpoints
		calls := api.Group("/calls")
		{
			calls.POST("/initiate", callHandler.InitiateCall)
			calls.POST("/:id/answer", callHandler.AnswerCall)
			calls.POST("/:id/end", callHandler.EndCall)
			calls.GET("/:id/status", callHandler.GetCallStatus)
		}

		// Presence management endpoints
		presence := api.Group("/presence")
		{
			presence.GET("/status", presenceHandler.GetPresenceStatus)
			presence.PUT("/status", presenceHandler.UpdatePresenceStatus)
			presence.GET("/check/:user_id", presenceHandler.CheckUserPresence)
		}

		// Call history endpoints
		api.GET("/history", historyHandler.GetCallHistory)

		// Service status endpoint
		api.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"service": "calling",
				"status":  "operational",
				"message": "Calling service is ready",
			})
		})
	}

	// WebSocket endpoint for signaling
	r.GET("/ws/calling", websocketHandler.HandleWebSocket)

	// Get port from environment
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8093" // Updated default port for calling service
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting calling service on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down calling service...")

	// Give the server 30 seconds to finish handling requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Calling service exited")
}
