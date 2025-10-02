package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tchat.dev/video/handlers"
	"tchat.dev/video/models"
	"tchat.dev/video/services"
)

// VideoServiceImpl implements ServiceInitializer for video service
type VideoServiceImpl struct {
	videoRepo     services.VideoRepository
	videoService  *services.VideoService
	videoHandlers *handlers.VideoHandlers
}

// NewVideoServiceImpl creates a new video service implementation
func NewVideoServiceImpl() *VideoServiceImpl {
	return &VideoServiceImpl{}
}

// GetModels returns all models for this service
func (v *VideoServiceImpl) GetModels() []interface{} {
	return []interface{}{
		&models.Channel{},
		&models.Video{},
		&models.VideoInteraction{},
		&models.VideoComment{},
		&models.VideoShare{},
	}
}

// InitializeRepositories initializes data access repositories
func (v *VideoServiceImpl) InitializeRepositories(db *gorm.DB) error {
	v.videoRepo = services.NewPostgreSQLVideoRepository(db)
	return nil
}

// InitializeServices initializes business logic services
func (v *VideoServiceImpl) InitializeServices(db *gorm.DB) error {
	v.videoService = services.NewVideoService(v.videoRepo, db)
	return nil
}

// InitializeHandlers initializes HTTP handlers
func (v *VideoServiceImpl) InitializeHandlers() error {
	v.videoHandlers = handlers.NewVideoHandlers(v.videoService)
	return nil
}

// RegisterRoutes registers all video-related routes
func (v *VideoServiceImpl) RegisterRoutes(router *gin.Engine) error {
	// API routes
	v1 := router.Group("/api/v1")
	{
		// Video routes
		videos := v1.Group("/videos")
		{
			// Video CRUD operations
			videos.GET("", v.videoHandlers.GetVideos)
			videos.POST("", v.videoHandlers.CreateVideo)
			videos.GET("/:id", v.videoHandlers.GetVideo)
			videos.PUT("/:id", v.videoHandlers.UpdateVideo)
			videos.DELETE("/:id", v.videoHandlers.DeleteVideo)

			// Video interactions
			videos.POST("/:id/like", v.videoHandlers.LikeVideo)
			videos.POST("/:id/share", v.videoHandlers.ShareVideo)

			// Video comments
			videos.GET("/:id/comments", v.videoHandlers.GetVideoComments)
			videos.POST("/:id/comments", v.videoHandlers.AddVideoComment)

			// Video upload
			videos.POST("/upload", v.videoHandlers.UploadVideo)

			// Video search and filtering
			videos.GET("/search", v.videoHandlers.SearchVideos)
			videos.GET("/category/:category", v.videoHandlers.GetVideoByCategory)
			videos.GET("/trending", v.videoHandlers.GetTrendingVideos)

			// Short videos (TikTok style)
			videos.GET("/shorts", v.videoHandlers.GetShortVideos)

			// Health check
			videos.GET("/health", v.videoHandlers.VideoHealth)
		}

		// Channel routes
		channels := v1.Group("/channels")
		{
			channels.POST("", v.videoHandlers.CreateChannel)
			channels.GET("/:id", v.videoHandlers.GetChannel)
		}
	}

	return nil
}

// GetServiceInfo returns service name and version
func (v *VideoServiceImpl) GetServiceInfo() (string, string) {
	return "video-service", "1.0.0"
}