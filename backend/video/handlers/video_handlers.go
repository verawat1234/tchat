package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/shared/responses"
	"tchat.dev/video/models"
	"tchat.dev/video/services"
)

// VideoHandlers handles HTTP requests for video operations
type VideoHandlers struct {
	videoService *services.VideoService
}

// NewVideoHandlers creates a new video handlers instance
func NewVideoHandlers(videoService *services.VideoService) *VideoHandlers {
	return &VideoHandlers{
		videoService: videoService,
	}
}

// GetVideos retrieves videos with pagination
func (h *VideoHandlers) GetVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	category := c.Query("category")

	videos, err := h.videoService.GetVideos(page, perPage, category)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve videos")
		return
	}

	// Build response
	response := gin.H{
		"videos":     videos,
		"page":       page,
		"per_page":   perPage,
		"total":      len(videos),
		"has_more":   len(videos) == perPage,
	}

	responses.SendSuccessResponse(c, response)
}

// GetShortVideos retrieves short-form videos (TikTok style)
func (h *VideoHandlers) GetShortVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	videos, err := h.videoService.GetShortVideos(page, perPage)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve short videos")
		return
	}

	// Build response compatible with existing video structure
	response := map[string]interface{}{
		"videos": videos,
	}

	responses.SendSuccessResponse(c, response)
}

// GetVideo retrieves a specific video by ID
func (h *VideoHandlers) GetVideo(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	// Get user ID from context if available (for tracking)
	var userID *uuid.UUID
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if parsedUserID, err := uuid.Parse(userIDStr); err == nil {
			userID = &parsedUserID
		}
	}

	video, err := h.videoService.GetVideoByID(videoID, userID)
	if err != nil {
		if err.Error() == "video not found" {
			responses.NotFoundResponse(c, "Video not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to retrieve video")
		return
	}

	responses.SendSuccessResponse(c, video)
}

// CreateVideo creates a new video
func (h *VideoHandlers) CreateVideo(c *gin.Context) {
	var video models.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.videoService.CreateVideo(&video); err != nil {
		responses.InternalErrorResponse(c, "Failed to create video")
		return
	}

	responses.SuccessMessageResponse(c, "Video created successfully")
}

// UpdateVideo updates an existing video
func (h *VideoHandlers) UpdateVideo(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	var video models.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	video.ID = videoID
	if err := h.videoService.UpdateVideo(&video); err != nil {
		if err.Error() == "video not found" {
			responses.NotFoundResponse(c, "Video not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to update video")
		return
	}

	responses.SuccessMessageResponse(c, "Video updated successfully")
}

// DeleteVideo deletes a video
func (h *VideoHandlers) DeleteVideo(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.videoService.DeleteVideo(videoID); err != nil {
		if err.Error() == "video not found" {
			responses.NotFoundResponse(c, "Video not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to delete video")
		return
	}

	responses.SuccessMessageResponse(c, "Video deleted successfully")
}

// LikeVideo likes a video
func (h *VideoHandlers) LikeVideo(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	// Get user ID from context
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.videoService.LikeVideo(videoID, userID); err != nil {
		if err.Error() == "video not found" {
			responses.NotFoundResponse(c, "Video not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to like video")
		return
	}

	responses.SuccessMessageResponse(c, "Video liked successfully")
}

// CreateChannel creates a new channel
func (h *VideoHandlers) CreateChannel(c *gin.Context) {
	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.videoService.CreateChannel(&channel); err != nil {
		responses.InternalErrorResponse(c, "Failed to create channel")
		return
	}

	responses.SuccessMessageResponse(c, "Channel created successfully")
}

// GetChannel retrieves a channel by ID
func (h *VideoHandlers) GetChannel(c *gin.Context) {
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	channel, err := h.videoService.GetChannelByID(channelID)
	if err != nil {
		if err.Error() == "channel not found" {
			responses.NotFoundResponse(c, "Channel not found")
			return
		}
		responses.InternalErrorResponse(c, "Failed to retrieve channel")
		return
	}

	responses.SendSuccessResponse(c, channel)
}

// UploadVideo handles video file upload
func (h *VideoHandlers) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	// Validate file type
	if !isValidVideoType(file.Header.Get("Content-Type")) {
		responses.BadRequestResponse(c, "Invalid format")
		return
	}

	// TODO: Upload to cloud storage (S3, GCS, etc.)
	// For now, return a mock URL
	videoURL := fmt.Sprintf("https://cdn.tchat.dev/videos/%s", file.Filename)

	responses.SendSuccessResponse(c, gin.H{
		"videoUrl": videoURL,
		"filename": file.Filename,
		"size":     file.Size,
	})
}

// SearchVideos searches videos by title, description, or tags
func (h *VideoHandlers) SearchVideos(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		responses.BadRequestResponse(c, "Invalid format")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	category := c.Query("category")

	videos, err := h.videoService.SearchVideos(query, page, perPage, category)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to search videos")
		return
	}

	response := gin.H{
		"videos":     videos,
		"query":      query,
		"page":       page,
		"per_page":   perPage,
		"total":      len(videos),
		"has_more":   len(videos) == perPage,
	}

	responses.SendSuccessResponse(c, response)
}

// GetVideoComments retrieves comments for a video
func (h *VideoHandlers) GetVideoComments(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	comments, err := h.videoService.GetVideoComments(videoID, page, perPage)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve comments")
		return
	}

	response := gin.H{
		"comments":   comments,
		"page":       page,
		"per_page":   perPage,
		"total":      len(comments),
		"has_more":   len(comments) == perPage,
	}

	responses.SendSuccessResponse(c, response)
}

// AddVideoComment adds a comment to a video
func (h *VideoHandlers) AddVideoComment(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	var commentData struct {
		Content string `json:"content" binding:"required"`
		ParentID *uuid.UUID `json:"parentId,omitempty"`
	}

	if err := c.ShouldBindJSON(&commentData); err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	// Get user ID from context
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	_, err = h.videoService.AddVideoComment(videoID, userID, commentData.Content, commentData.ParentID)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to add comment")
		return
	}

	responses.SuccessMessageResponse(c, "Comment added successfully")
}

// ShareVideo handles video sharing
func (h *VideoHandlers) ShareVideo(c *gin.Context) {
	videoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	var shareData struct {
		Platform string `json:"platform" binding:"required"`
		Message  string `json:"message,omitempty"`
	}

	if err := c.ShouldBindJSON(&shareData); err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	// Get user ID from context
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		responses.UnauthorizedResponse(c, "User authentication required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responses.BadRequestResponse(c, err.Error())
		return
	}

	_, err = h.videoService.ShareVideo(videoID, userID, shareData.Platform, shareData.Message)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to share video")
		return
	}

	responses.SuccessMessageResponse(c, "Video shared successfully")
}

// GetVideoByCategory retrieves videos by category
func (h *VideoHandlers) GetVideoByCategory(c *gin.Context) {
	category := c.Param("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	videos, err := h.videoService.GetVideosByCategory(category, page, perPage)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve videos")
		return
	}

	response := gin.H{
		"videos":     videos,
		"category":   category,
		"page":       page,
		"per_page":   perPage,
		"total":      len(videos),
		"has_more":   len(videos) == perPage,
	}

	responses.SendSuccessResponse(c, response)
}

// GetTrendingVideos retrieves trending videos
func (h *VideoHandlers) GetTrendingVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	timeframe := c.DefaultQuery("timeframe", "day") // day, week, month

	videos, err := h.videoService.GetTrendingVideos(timeframe, page, perPage)
	if err != nil {
		responses.InternalErrorResponse(c, "Failed to retrieve trending videos")
		return
	}

	response := gin.H{
		"videos":     videos,
		"timeframe":  timeframe,
		"page":       page,
		"per_page":   perPage,
		"total":      len(videos),
		"has_more":   len(videos) == perPage,
	}

	responses.SendSuccessResponse(c, response)
}

// VideoHealth provides health check for video service
func (h *VideoHandlers) VideoHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "video-service",
		"api":       "available",
		"timestamp": gin.H{},
	})
}

// Helper function to validate video file types
func isValidVideoType(contentType string) bool {
	validTypes := []string{
		"video/mp4",
		"video/avi",
		"video/quicktime", // .mov files
		"video/x-msvideo", // .avi files
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}