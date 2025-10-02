package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tchat.dev/shared/logger"
)

// VideoSocialHandler handles social interactions for video content
type VideoSocialHandler struct {
	logger *logger.TchatLogger
}

// NewVideoSocialHandler creates a new video social handler instance
func NewVideoSocialHandler(logger *logger.TchatLogger) *VideoSocialHandler {
	return &VideoSocialHandler{
		logger: logger,
	}
}

// VideoInteraction represents a social interaction on a video
type VideoInteraction struct {
	ID          string `json:"id"`
	VideoID     string `json:"video_id"`
	UserID      string `json:"user_id"`
	Type        string `json:"type"`         // like, dislike, share, comment, favorite
	CommentText string `json:"comment_text"` // For comment interactions
	CreatedAt   string `json:"created_at"`
}

// VideoInteractionStats represents aggregated video interaction statistics
type VideoInteractionStats struct {
	VideoID        string `json:"video_id"`
	LikeCount      int64  `json:"like_count"`
	DislikeCount   int64  `json:"dislike_count"`
	CommentCount   int64  `json:"comment_count"`
	ShareCount     int64  `json:"share_count"`
	FavoriteCount  int64  `json:"favorite_count"`
	ViewCount      int64  `json:"view_count"`
	EngagementRate float64 `json:"engagement_rate"`
}

// LikeVideo handles video like action
// POST /api/v1/videos/:id/like
func (h *VideoSocialHandler) LikeVideo(c *gin.Context) {
	videoID := c.Param("id")
	if _, err := uuid.Parse(videoID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_video_id",
			"message": "Invalid video ID format",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	interaction := VideoInteraction{
		ID:        uuid.New().String(),
		VideoID:   videoID,
		UserID:    userID.(string),
		Type:      "like",
		CreatedAt: "2024-01-01T00:00:00Z", // Would use time.Now()
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id": videoID,
		"user_id":  userID,
		"action":   "like",
	}).Info("Video liked")

	c.JSON(http.StatusCreated, gin.H{
		"interaction": interaction,
		"message":     "Video liked successfully",
	})
}

// UnlikeVideo handles video unlike action
// DELETE /api/v1/videos/:id/like
func (h *VideoSocialHandler) UnlikeVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	h.logger.WithFields(map[string]interface{}{
		"video_id": videoID,
		"user_id":  userID,
		"action":   "unlike",
	}).Info("Video unliked")

	c.JSON(http.StatusOK, gin.H{
		"message": "Video unliked successfully",
	})
}

// CommentOnVideo handles video comment creation
// POST /api/v1/videos/:id/comments
func (h *VideoSocialHandler) CommentOnVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req struct {
		CommentText string `json:"comment_text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Comment text is required",
		})
		return
	}

	interaction := VideoInteraction{
		ID:          uuid.New().String(),
		VideoID:     videoID,
		UserID:      userID.(string),
		Type:        "comment",
		CommentText: req.CommentText,
		CreatedAt:   "2024-01-01T00:00:00Z",
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id":     videoID,
		"user_id":      userID,
		"comment_text": req.CommentText,
		"action":       "comment",
	}).Info("Video commented")

	c.JSON(http.StatusCreated, gin.H{
		"interaction": interaction,
		"message":     "Comment added successfully",
	})
}

// ShareVideo handles video share action
// POST /api/v1/videos/:id/share
func (h *VideoSocialHandler) ShareVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req struct {
		Platform string `json:"platform"` // facebook, twitter, whatsapp, etc.
	}

	c.ShouldBindJSON(&req)

	interaction := VideoInteraction{
		ID:        uuid.New().String(),
		VideoID:   videoID,
		UserID:    userID.(string),
		Type:      "share",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id": videoID,
		"user_id":  userID,
		"platform": req.Platform,
		"action":   "share",
	}).Info("Video shared")

	c.JSON(http.StatusCreated, gin.H{
		"interaction": interaction,
		"share_url":   "https://tchat.dev/videos/" + videoID,
		"message":     "Video shared successfully",
	})
}

// FavoriteVideo handles video favorite action
// POST /api/v1/videos/:id/favorite
func (h *VideoSocialHandler) FavoriteVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	interaction := VideoInteraction{
		ID:        uuid.New().String(),
		VideoID:   videoID,
		UserID:    userID.(string),
		Type:      "favorite",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id": videoID,
		"user_id":  userID,
		"action":   "favorite",
	}).Info("Video favorited")

	c.JSON(http.StatusCreated, gin.H{
		"interaction": interaction,
		"message":     "Video added to favorites",
	})
}

// UnfavoriteVideo handles video unfavorite action
// DELETE /api/v1/videos/:id/favorite
func (h *VideoSocialHandler) UnfavoriteVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	h.logger.WithFields(map[string]interface{}{
		"video_id": videoID,
		"user_id":  userID,
		"action":   "unfavorite",
	}).Info("Video unfavorited")

	c.JSON(http.StatusOK, gin.H{
		"message": "Video removed from favorites",
	})
}

// GetVideoComments retrieves comments for a video
// GET /api/v1/videos/:id/comments
func (h *VideoSocialHandler) GetVideoComments(c *gin.Context) {
	videoID := c.Param("id")

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Mock data - would fetch from database
	comments := []VideoInteraction{
		{
			ID:          uuid.New().String(),
			VideoID:     videoID,
			UserID:      uuid.New().String(),
			Type:        "comment",
			CommentText: "Great video!",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id": videoID,
		"page":     page,
		"limit":    limit,
	}).Info("Video comments retrieved")

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"page":     page,
		"limit":    limit,
		"total":    len(comments),
	})
}

// GetVideoStats retrieves aggregated interaction statistics for a video
// GET /api/v1/videos/:id/stats
func (h *VideoSocialHandler) GetVideoStats(c *gin.Context) {
	videoID := c.Param("id")

	// Mock data - would aggregate from database
	stats := VideoInteractionStats{
		VideoID:        videoID,
		LikeCount:      1250,
		DislikeCount:   23,
		CommentCount:   387,
		ShareCount:     156,
		FavoriteCount:  542,
		ViewCount:      15678,
		EngagementRate: 14.8, // (likes + comments + shares) / views * 100
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id":        videoID,
		"like_count":      stats.LikeCount,
		"comment_count":   stats.CommentCount,
		"engagement_rate": stats.EngagementRate,
	}).Info("Video stats retrieved")

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GetUserVideoInteractions retrieves a user's interactions with a specific video
// GET /api/v1/videos/:id/user-interactions
func (h *VideoSocialHandler) GetUserVideoInteractions(c *gin.Context) {
	videoID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	// Mock data - would fetch from database
	interactions := map[string]bool{
		"liked":     true,
		"favorited": false,
		"shared":    true,
		"commented": true,
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id": videoID,
		"user_id":  userID,
	}).Info("User video interactions retrieved")

	c.JSON(http.StatusOK, gin.H{
		"video_id":     videoID,
		"user_id":      userID,
		"interactions": interactions,
	})
}

// GetTrendingVideos retrieves trending videos based on social engagement
// GET /api/v1/videos/trending
func (h *VideoSocialHandler) GetTrendingVideos(c *gin.Context) {
	timeframe := c.DefaultQuery("timeframe", "day") // day, week, month
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Mock data - would calculate from database
	trendingVideos := []gin.H{
		{
			"video_id":        uuid.New().String(),
			"engagement_rate": 18.5,
			"viral_score":     92.3,
			"trend_velocity":  "rising",
		},
	}

	h.logger.WithFields(map[string]interface{}{
		"timeframe": timeframe,
		"limit":     limit,
	}).Info("Trending videos retrieved")

	c.JSON(http.StatusOK, gin.H{
		"videos":    trendingVideos,
		"timeframe": timeframe,
		"limit":     limit,
		"total":     len(trendingVideos),
	})
}