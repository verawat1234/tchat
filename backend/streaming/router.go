package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"tchat.dev/streaming/handlers"
	"tchat.dev/streaming/middleware"
)

// SetupRouter configures all routes and middleware for the streaming service
func SetupRouter(h *handlers.Handlers, authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	// CORS configuration for web and mobile clients
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://tchat.com", "https://www.tchat.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Error recovery middleware
	router.Use(gin.Recovery())

	// Request logging middleware
	router.Use(middleware.RequestLogger())

	// Health check endpoint (public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "streaming",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Stream lifecycle endpoints (authenticated)
		streams := v1.Group("/streams")
		streams.Use(authMiddleware)
		{
			// Stream CRUD operations
			streams.POST("", h.CreateStream.Handle)           // Create stream
			streams.GET("", h.ListStreams.Handle)             // List streams
			streams.GET("/:streamId", h.GetStream.Handle)     // Get stream details
			streams.PATCH("/:streamId", h.UpdateStream.Handle) // Update stream

			// Stream control operations
			streams.POST("/:streamId/start", h.StartStream.Handle) // Start stream (WebRTC)
			streams.POST("/:streamId/end", h.EndStream.Handle)     // End stream

			// Chat endpoints with rate limiting
			chatGroup := streams.Group("/:streamId/chat")
			chatGroup.Use(middleware.RateLimiter(10, time.Second)) // 10 messages/sec per user
			{
				chatGroup.POST("", h.SendChat.Handle)              // Send chat message
				chatGroup.GET("", h.GetChat.Handle)                // Get chat history
				chatGroup.DELETE("/:messageId", h.DeleteChat.Handle) // Delete chat message
			}

			// Reaction endpoint with rate limiting
			streams.POST("/:streamId/react",
				middleware.RateLimiter(10, time.Second), // 10 reactions/sec per user
				h.SendReaction.Handle)

			// Store product endpoints
			streams.POST("/:streamId/products", h.FeatureProduct.Handle) // Feature product in stream
			streams.GET("/:streamId/products", h.ListProducts.Handle)    // List featured products

			// Analytics endpoint
			streams.GET("/:streamId/analytics", h.GetAnalytics.Handle) // Get stream analytics
		}

		// Notification preferences (authenticated)
		prefs := v1.Group("/notification-preferences")
		prefs.Use(authMiddleware)
		{
			prefs.GET("", h.NotificationPreferences.HandleGet)  // Get preferences
			prefs.PUT("", h.NotificationPreferences.HandlePut)  // Update preferences
		}

		// WebSocket signaling endpoint (authenticated via query token)
		v1.GET("/ws/signaling", middleware.WebSocketAuth(), func(c *gin.Context) {
			handlers.SignalingWebSocket(h.SignalingService)(c)
		})
	}

	return router
}