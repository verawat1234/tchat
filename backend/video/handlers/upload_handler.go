// backend/video/handlers/upload_handler.go
// Video Upload Handler - HTTP endpoints for video upload operations
// Implements T033: Video upload handler POST /api/v1/videos

package handlers

import (
	"tchat.dev/video/repository"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/video/models"
	"tchat.dev/video/services"
)

// UploadHandler handles video upload operations
type UploadHandler struct {
	videoService services.VideoServiceInterface
}

// NewUploadHandler creates a new upload handler instance
func NewUploadHandler(videoService services.VideoServiceInterface) *UploadHandler {
	return &UploadHandler{
		videoService: videoService,
	}
}

// RegisterRoutes registers upload routes with the router
func (h *UploadHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/videos", h.UploadVideo)
	router.GET("/videos/:id", h.GetVideo)
	router.PUT("/videos/:id", h.UpdateVideo)
	router.DELETE("/videos/:id", h.DeleteVideo)
	router.GET("/videos", h.ListVideos)
	router.GET("/videos/creator/:creatorId", h.GetVideosByCreator)
}

// UploadVideo handles POST /api/v1/videos
// @Summary Upload a new video
// @Description Upload video file with metadata for processing and streaming
// @Tags videos
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Video file"
// @Param title formData string true "Video title"
// @Param description formData string false "Video description"
// @Param tags formData string false "Comma-separated tags"
// @Param content_rating formData string false "Content rating (G, PG, PG13, R, NC17)"
// @Success 201 {object} VideoUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse "File too large"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/videos [post]
func (h *UploadHandler) UploadVideo(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	creatorID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(5 << 30); err != nil { // 5GB max
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "form_parse_error",
			Message: fmt.Sprintf("Failed to parse form: %v", err),
		})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "file_required",
			Message: "Video file is required",
		})
		return
	}
	defer file.Close()

	// Parse metadata
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "title_required",
			Message: "Video title is required",
		})
		return
	}

	description := c.PostForm("description")
	tagsStr := c.PostForm("tags")
	contentRatingStr := c.PostForm("content_rating")

	// Parse tags
	var tags []string
	if tagsStr != "" {
		tags = parseTags(tagsStr)
	}

	// Parse content rating
	var contentRating models.ContentRating
	if contentRatingStr != "" {
		contentRating = models.ContentRating(contentRatingStr)
		// Validate content rating
		validRatings := []models.ContentRating{
			models.RatingG, models.RatingPG, models.RatingPG13,
			models.RatingR, models.RatingNC17,
		}
		isValid := false
		for _, rating := range validRatings {
			if contentRating == rating {
				isValid = true
				break
			}
		}
		if !isValid {
			contentRating = models.RatingG // Default to G
		}
	} else {
		contentRating = models.RatingG
	}

	// Create video metadata
	metadata := services.VideoMetadata{
		Title:         title,
		Description:   description,
		CreatorID:     creatorID,
		ContentRating: contentRating,
		Tags:          tags,
	}

	// Upload video
	video, err := h.videoService.UploadVideo(c.Request.Context(), file, header, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "upload_failed",
			Message: fmt.Sprintf("Failed to upload video: %v", err),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, VideoUploadResponse{
		VideoID:       video.ID,
		Title:         video.Title,
		Status:        string(video.AvailabilityStatus),
		UploadURL:     video.ThumbnailURL,
		ProcessingETA: "2-5 minutes", // Estimate based on file size
		Message:       "Video upload started. Processing will begin shortly.",
	})
}

// GetVideo handles GET /api/v1/videos/:id
// @Summary Get video details
// @Description Retrieve video metadata and streaming information
// @Tags videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} VideoDetailsResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/videos/{id} [get]
func (h *UploadHandler) GetVideo(c *gin.Context) {
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_video_id",
			Message: "Invalid video ID format",
		})
		return
	}

	// Get video
	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "video_not_found",
			Message: "Video not found",
		})
		return
	}

	// Return video details
	c.JSON(http.StatusOK, VideoDetailsResponse{
		ID:                 video.ID,
		Title:              video.Title,
		Description:        video.Description,
		CreatorID:          video.CreatorID,
		Duration:           video.Duration,
		ThumbnailURL:       video.ThumbnailURL,
		AvailabilityStatus: string(video.AvailabilityStatus),
		ContentRating:      string(video.ContentRating),
		Tags:               video.Tags,
		QualityOptions:     video.QualityOptions,
		UploadTimestamp:    video.UploadTimestamp,
		PublishTimestamp:   &video.UploadTimestamp,
		ViewCount:          video.SocialMetrics.ViewCount,
		LikeCount:          video.SocialMetrics.LikeCount,
		CommentCount:       video.SocialMetrics.CommentCount,
		ShareCount:         video.SocialMetrics.ShareCount,
	})
}

// UpdateVideo handles PUT /api/v1/videos/:id
// @Summary Update video metadata
// @Description Update video title, description, tags, and other metadata
// @Tags videos
// @Accept json
// @Produce json
// @Param id path string true "Video ID"
// @Param video body UpdateVideoRequest true "Video metadata"
// @Success 200 {object} VideoDetailsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/videos/{id} [put]
func (h *UploadHandler) UpdateVideo(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_video_id",
			Message: "Invalid video ID format",
		})
		return
	}

	// Parse request body
	var req UpdateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	// Get video to check ownership
	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "video_not_found",
			Message: "Video not found",
		})
		return
	}

	// Check if user is the creator
	if video.CreatorID.String() != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "You can only update your own videos",
		})
		return
	}

	// Update video metadata
	metadata := services.VideoMetadata{
		Title:         req.Title,
		Description:   req.Description,
		Tags:          req.Tags,
		ContentRating: models.ContentRating(req.ContentRating),
	}

	if err := h.videoService.UpdateVideoMetadata(c.Request.Context(), videoID, metadata); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: fmt.Sprintf("Failed to update video: %v", err),
		})
		return
	}

	// Get updated video
	updatedVideo, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "fetch_failed",
			Message: "Video updated but failed to fetch updated details",
		})
		return
	}

	// Return updated video
	c.JSON(http.StatusOK, VideoDetailsResponse{
		ID:                 updatedVideo.ID,
		Title:              updatedVideo.Title,
		Description:        updatedVideo.Description,
		CreatorID:          updatedVideo.CreatorID,
		Duration:           updatedVideo.Duration,
		ThumbnailURL:       updatedVideo.ThumbnailURL,
		AvailabilityStatus: string(updatedVideo.AvailabilityStatus),
		ContentRating:      string(updatedVideo.ContentRating),
		Tags:               updatedVideo.Tags,
		QualityOptions:     updatedVideo.QualityOptions,
		UploadTimestamp:    updatedVideo.UploadTimestamp,
		PublishTimestamp:   &updatedVideo.UploadTimestamp,
		ViewCount:          updatedVideo.SocialMetrics.ViewCount,
		LikeCount:          updatedVideo.SocialMetrics.LikeCount,
		CommentCount:       updatedVideo.SocialMetrics.CommentCount,
		ShareCount:         updatedVideo.SocialMetrics.ShareCount,
	})
}

// DeleteVideo handles DELETE /api/v1/videos/:id
// @Summary Delete a video
// @Description Delete video and all associated data
// @Tags videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/videos/{id} [delete]
func (h *UploadHandler) DeleteVideo(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_video_id",
			Message: "Invalid video ID format",
		})
		return
	}

	// Get video to check ownership
	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "video_not_found",
			Message: "Video not found",
		})
		return
	}

	// Check if user is the creator
	if video.CreatorID.String() != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "You can only delete your own videos",
		})
		return
	}

	// Delete video
	if err := h.videoService.DeleteVideo(c.Request.Context(), videoID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "delete_failed",
			Message: fmt.Sprintf("Failed to delete video: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Video deleted successfully",
	})
}

// ListVideos handles GET /api/v1/videos
// @Summary List videos
// @Description List videos with pagination and filtering
// @Tags videos
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Param status query string false "Filter by status (public, private, processing)"
// @Param tags query string false "Filter by tags (comma-separated)"
// @Success 200 {object} VideoListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/videos [get]
func (h *UploadHandler) ListVideos(c *gin.Context) {
	// Parse pagination parameters
	page := parseIntQueryParam(c, "page", 1)
	perPage := parseIntQueryParam(c, "per_page", 20)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	// Parse filters
	statusStr := c.Query("status")
	tagsStr := c.Query("tags")

	// Build search filters
	filters := repository.VideoSearchFilters{
		Status: models.AvailabilityStatus(statusStr),
	}

	if tagsStr != "" {
		filters.Tags = parseTags(tagsStr)
	}

	// Search videos
	pagination := services.Pagination{
		Limit:  perPage,
		Offset: offset,
	}

	results, err := h.videoService.SearchVideos(c.Request.Context(), "", filters, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "search_failed",
			Message: fmt.Sprintf("Failed to search videos: %v", err),
		})
		return
	}

	// Convert to response format
	videos := make([]VideoSummary, len(results.Videos))
	for i, video := range results.Videos {
		videos[i] = VideoSummary{
			ID:           video.ID,
			Title:        video.Title,
			CreatorID:    video.CreatorID,
			ThumbnailURL: video.ThumbnailURL,
			Duration:     video.Duration,
			ViewCount:    video.SocialMetrics.ViewCount,
			UploadedAt:   video.UploadTimestamp,
		}
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Videos:  videos,
		Total:   results.Total,
		Page:    results.Page,
		PerPage: results.PageSize,
		HasMore: results.HasMore,
	})
}

// GetVideosByCreator handles GET /api/v1/videos/creator/:creatorId
// @Summary Get videos by creator
// @Description List all videos uploaded by a specific creator
// @Tags videos
// @Produce json
// @Param creatorId path string true "Creator ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} VideoListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/videos/creator/{creatorId} [get]
func (h *UploadHandler) GetVideosByCreator(c *gin.Context) {
	creatorIDStr := c.Param("creatorId")
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_creator_id",
			Message: "Invalid creator ID format",
		})
		return
	}

	// Parse pagination
	page := parseIntQueryParam(c, "page", 1)
	perPage := parseIntQueryParam(c, "per_page", 20)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	pagination := services.Pagination{
		Limit:  perPage,
		Offset: offset,
	}

	// Get videos by creator
	videos, err := h.videoService.GetVideosByCreator(c.Request.Context(), creatorID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "fetch_failed",
			Message: fmt.Sprintf("Failed to fetch videos: %v", err),
		})
		return
	}

	// Convert to response format
	videoSummaries := make([]VideoSummary, len(videos))
	for i, video := range videos {
		videoSummaries[i] = VideoSummary{
			ID:           video.ID,
			Title:        video.Title,
			CreatorID:    video.CreatorID,
			ThumbnailURL: video.ThumbnailURL,
			Duration:     video.Duration,
			ViewCount:    video.SocialMetrics.ViewCount,
			UploadedAt:   video.UploadTimestamp,
		}
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Videos:  videoSummaries,
		Total:   int64(len(videoSummaries)),
		Page:    page,
		PerPage: perPage,
		HasMore: len(videoSummaries) == perPage,
	})
}

// Helper functions
// Note: parseTags is defined in streaming_handler.go to avoid duplication

func splitCSV(s string) []string {
	var result []string
	current := ""

	for _, char := range s {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

func parseIntQueryParam(c *gin.Context, param string, defaultValue int) int {
	valueStr := c.Query(param)
	if valueStr == "" {
		return defaultValue
	}

	var value int
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return defaultValue
	}

	return value
}

// Response types

type VideoUploadResponse struct {
	VideoID       uuid.UUID `json:"video_id"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	UploadURL     string    `json:"upload_url"`
	ProcessingETA string    `json:"processing_eta"`
	Message       string    `json:"message"`
}

type VideoDetailsResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	CreatorID          uuid.UUID  `json:"creator_id"`
	Duration           int        `json:"duration"`
	ThumbnailURL       string     `json:"thumbnail_url"`
	AvailabilityStatus string     `json:"availability_status"`
	ContentRating      string     `json:"content_rating"`
	Tags               []string   `json:"tags"`
	QualityOptions     []string   `json:"quality_options"`
	UploadTimestamp    time.Time  `json:"uploaded_at"`
	PublishTimestamp   *time.Time `json:"published_at,omitempty"`
	ViewCount          int64      `json:"view_count"`
	LikeCount          int64      `json:"like_count"`
	CommentCount       int64      `json:"comment_count"`
	ShareCount         int64      `json:"share_count"`
}

type VideoSummary struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	CreatorID    uuid.UUID `json:"creator_id"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Duration     int       `json:"duration"`
	ViewCount    int64     `json:"view_count"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

type VideoListResponse struct {
	Videos  []VideoSummary `json:"videos"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
	HasMore bool           `json:"has_more"`
}

type UpdateVideoRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	ContentRating string   `json:"content_rating"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}