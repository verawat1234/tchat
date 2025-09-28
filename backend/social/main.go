package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"tchat/social/database"
	"tchat/social/handlers"
	"tchat/social/repository"
	"tchat/social/services"
	"tchat.dev/shared/config"
)

func main() {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Initialize configuration
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	cfg := &config.Config{
		Debug: os.Getenv("GIN_MODE") == "debug" || os.Getenv("DEBUG") == "true",
		Database: config.DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "tchat_social"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}

	// Initialize database
	dbConfig := database.NewSocialDatabaseConfig()
	if err := dbConfig.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if sqlDB, err := dbConfig.GetDB().DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Initialize repository manager
	repoManager := repository.NewManager(dbConfig.GetDB())
	defer repoManager.Close()

	// Initialize services with real database
	userService := services.NewUserService(repoManager)
	postService := services.NewPostService(repoManager)
	syncService := services.NewMobileSyncService(repoManager)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	mobileHandler := handlers.NewMobileHandler(syncService, userService)

	// Setup routes
	setupRoutes(router, userHandler, postHandler, mobileHandler)

	// Get port from environment or default to 8092
	port := os.Getenv("PORT")
	if port == "" {
		port = "8092"
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("🚀 Social service starting on port %s\n", port)
		fmt.Printf("📱 Health check: http://localhost:%s/health\n", port)
		fmt.Printf("📋 API docs: http://localhost:%s/api/v1/social\n", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down social service...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Social service forced to shutdown: %v", err)
	}

	fmt.Println("✅ Social service exited gracefully")
}

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupRoutes configures all API routes
func setupRoutes(router *gin.Engine, userHandler *handlers.UserHandler, postHandler *handlers.PostHandler, mobileHandler *handlers.MobileHandler) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "social",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Social routes group
		social := v1.Group("/social")
		{
			// User profile routes
			social.GET("/profiles/:userId", userHandler.GetSocialProfile)
			social.PUT("/profiles/:userId", userHandler.UpdateSocialProfile)

			// User discovery and relationships
			social.GET("/discover/users", userHandler.DiscoverUsers)
			social.POST("/follow", userHandler.FollowUser)
			social.DELETE("/follow/:followerId/:followingId", userHandler.UnfollowUser)
			social.GET("/followers/:userId", userHandler.GetFollowers)
			social.GET("/following/:userId", userHandler.GetFollowing)

			// User analytics
			social.GET("/analytics/users/:userId", userHandler.GetUserAnalytics)

			// Post routes
			social.POST("/posts", postHandler.CreatePost)
			social.GET("/posts/:postId", postHandler.GetPost)
			social.PUT("/posts/:postId", postHandler.UpdatePost)
			social.DELETE("/posts/:postId", postHandler.DeletePost)

			// Interaction routes
			social.POST("/reactions", postHandler.AddReaction)
			social.DELETE("/reactions/:targetType/:targetId", postHandler.RemoveReaction)
			social.POST("/comments", postHandler.CreateComment)

			// Content routes
			social.GET("/feed", postHandler.GetSocialFeed)
			social.GET("/trending", postHandler.GetTrendingContent)
			social.POST("/share", postHandler.ShareContent)
		}

		// Register mobile-optimized routes
		mobileHandler.RegisterMobileRoutes(v1)
	}

	// Add a route to list all available endpoints
	router.GET("/api/v1/social", func(c *gin.Context) {
		endpoints := map[string]interface{}{
			"service": "Social Service API",
			"version": "1.0.0",
			"endpoints": map[string]interface{}{
				"profiles": map[string]string{
					"GET /api/v1/social/profiles/:userId":    "Get user social profile",
					"PUT /api/v1/social/profiles/:userId":    "Update user social profile",
				},
				"discovery": map[string]string{
					"GET /api/v1/social/discover/users": "Discover users",
				},
				"relationships": map[string]string{
					"POST /api/v1/social/follow":                         "Follow a user",
					"DELETE /api/v1/social/follow/:followerId/:followingId": "Unfollow a user",
					"GET /api/v1/social/followers/:userId":               "Get user followers",
					"GET /api/v1/social/following/:userId":               "Get users being followed",
				},
				"posts": map[string]string{
					"POST /api/v1/social/posts":          "Create a post",
					"GET /api/v1/social/posts/:postId":   "Get a post",
					"PUT /api/v1/social/posts/:postId":   "Update a post",
					"DELETE /api/v1/social/posts/:postId": "Delete a post",
				},
				"interactions": map[string]string{
					"POST /api/v1/social/reactions":                      "Add reaction to post/comment",
					"DELETE /api/v1/social/reactions/:targetType/:targetId": "Remove reaction",
					"POST /api/v1/social/comments":                       "Create comment on post",
				},
				"content": map[string]string{
					"GET /api/v1/social/feed":     "Get personalized social feed",
					"GET /api/v1/social/trending": "Get trending content",
					"POST /api/v1/social/share":   "Share content externally",
				},
				"analytics": map[string]string{
					"GET /api/v1/social/analytics/users/:userId": "Get user analytics",
				},
			},
			"regions_supported": []string{"TH", "SG", "MY", "ID", "PH", "VN"},
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		c.JSON(http.StatusOK, endpoints)
	})
}

// corsMiddleware adds CORS headers for cross-origin requests
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}