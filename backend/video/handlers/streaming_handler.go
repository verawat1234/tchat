package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/video/models"
	"tchat.dev/video/services"
)

// StreamingHandler handles video streaming endpoints
type StreamingHandler struct {
	videoService *services.VideoService
	syncService  *services.SyncService
}

// NewStreamingHandler creates a new streaming handler
func NewStreamingHandler(videoService *services.VideoService, syncService *services.SyncService) *StreamingHandler {
	return &StreamingHandler{
		videoService: videoService,
		syncService:  syncService,
	}
}

// RegisterRoutes registers streaming routes with the Gin router
func (h *StreamingHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/videos/:id/stream", h.GetStreamURL)
	router.GET("/videos/:id/manifest", h.GetStreamManifest)
	router.GET("/videos/:id/quality/:quality", h.GetQualityStream)
	router.GET("/videos/:id/thumbnail", h.GetThumbnail)
	router.GET("/videos/:id/preview", h.GetPreviewClip)
	router.POST("/videos/:id/analytics", h.RecordAnalytics)
	router.GET("/videos/:id/playback-info", h.GetPlaybackInfo)
}

// StreamURLResponse represents streaming URL response
type StreamURLResponse struct {
	VideoID           uuid.UUID          `json:"video_id"`
	StreamURL         string             `json:"stream_url"`
	Protocol          string             `json:"protocol"`
	ManifestURL       string             `json:"manifest_url,omitempty"`
	AvailableQualities []string          `json:"available_qualities"`
	DefaultQuality    string             `json:"default_quality"`
	DurationSeconds   int                `json:"duration_seconds"`
	Status            string             `json:"status"`
	ExpiresAt         time.Time          `json:"expires_at"`
}

// ManifestResponse represents adaptive stream manifest
type ManifestResponse struct {
	VideoID        uuid.UUID           `json:"video_id"`
	Protocol       string              `json:"protocol"`
	ManifestURL    string              `json:"manifest_url"`
	ManifestType   string              `json:"manifest_type"` // HLS, DASH
	QualityLevels  []QualityLevel      `json:"quality_levels"`
	AdaptiveConfig AdaptiveConfig      `json:"adaptive_config"`
}

// QualityLevel represents a streaming quality option
type QualityLevel struct {
	Quality      string  `json:"quality"`
	Resolution   string  `json:"resolution"`
	Bitrate      int     `json:"bitrate"`
	Framerate    int     `json:"framerate"`
	CodecProfile string  `json:"codec_profile"`
	StreamURL    string  `json:"stream_url"`
}

// AdaptiveConfig represents adaptive streaming configuration
type AdaptiveConfig struct {
	Enabled            bool   `json:"enabled"`
	MinBitrate         int    `json:"min_bitrate"`
	MaxBitrate         int    `json:"max_bitrate"`
	StartBitrate       int    `json:"start_bitrate"`
	BufferDuration     int    `json:"buffer_duration"`
	SwitchThreshold    int    `json:"switch_threshold"`
}

// PlaybackInfoResponse represents playback session information
type PlaybackInfoResponse struct {
	VideoID          uuid.UUID          `json:"video_id"`
	SessionID        uuid.UUID          `json:"session_id"`
	CurrentPosition  int                `json:"current_position"`
	Duration         int                `json:"duration"`
	Quality          string             `json:"quality"`
	Platform         models.PlatformType `json:"platform"`
	BufferHealth     BufferHealth       `json:"buffer_health"`
	SyncStatus       SyncStatus         `json:"sync_status"`
}

// BufferHealth represents buffer status
type BufferHealth struct {
	BufferedSeconds  int     `json:"buffered_seconds"`
	BufferPercentage float64 `json:"buffer_percentage"`
	IsHealthy        bool    `json:"is_healthy"`
}

// SyncStatus represents cross-platform sync status
type SyncStatus struct {
	LastSyncTime     time.Time `json:"last_sync_time"`
	SyncedPlatforms  []string  `json:"synced_platforms"`
	ConflictDetected bool      `json:"conflict_detected"`
}

// AnalyticsRequest represents analytics data submission
type AnalyticsRequest struct {
	EventType    string            `json:"event_type" binding:"required"`
	Position     int               `json:"position"`
	Quality      string            `json:"quality"`
	BufferEvents int               `json:"buffer_events"`
	ErrorCount   int               `json:"error_count"`
	Platform     string            `json:"platform"`
	Metadata     map[string]string `json:"metadata"`
}

// GetStreamURL retrieves streaming URL for a video
// GET /api/v1/videos/:id/stream
func (h *StreamingHandler) GetStreamURL(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Get video details
	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Check video availability
	if video.AvailabilityStatus != models.StatusPublic {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":  "Video not available for streaming",
			"status": string(video.AvailabilityStatus),
		})
		return
	}

	// Get quality parameter (optional)
	requestedQuality := c.DefaultQuery("quality", "auto")

	// Get platform from header or query
	platformStr := c.DefaultQuery("platform", "web")
	var platformType models.PlatformType
	switch platformStr {
	case "android":
		platformType = models.PlatformAndroid
	case "ios":
		platformType = models.PlatformIOS
	default:
		platformType = models.PlatformWeb
	}

	// Get user ID from context (assuming middleware sets this)
	userIDStr, _ := c.Get("user_id")
	userID := uuid.Nil
	if uid, ok := userIDStr.(uuid.UUID); ok {
		userID = uid
	}

	// Generate streaming URL with token using correct method
	streamingURLResp, err := h.videoService.GetStreamingURL(c.Request.Context(), videoID, userID, requestedQuality, platformType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate stream URL"})
		return
	}

	// Generate manifest URL for adaptive streaming
	manifestURL := fmt.Sprintf("/api/v1/videos/%s/manifest", videoID.String())

	// Get available qualities from quality options
	availableQualities := video.QualityOptions
	if len(availableQualities) == 0 {
		availableQualities = []string{"360p", "720p", "1080p"}
	}

	// Determine default quality
	defaultQuality := streamingURLResp.Quality

	response := StreamURLResponse{
		VideoID:            videoID,
		StreamURL:          streamingURLResp.URL,
		Protocol:           "HLS",
		ManifestURL:        manifestURL,
		AvailableQualities: availableQualities,
		DefaultQuality:     defaultQuality,
		DurationSeconds:    video.Duration,
		Status:             string(video.AvailabilityStatus),
		ExpiresAt:          streamingURLResp.ExpiresAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetStreamManifest retrieves adaptive stream manifest
// GET /api/v1/videos/:id/manifest
func (h *StreamingHandler) GetStreamManifest(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Get video details
	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Check video availability
	if video.AvailabilityStatus != models.StatusPublic {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Video manifest not available",
		})
		return
	}

	// Get manifest type from query (HLS or DASH)
	manifestType := c.DefaultQuery("type", "HLS")

	// Build quality levels
	qualityLevels := []QualityLevel{
		{
			Quality:      "360p",
			Resolution:   "640x360",
			Bitrate:      800000,
			Framerate:    30,
			CodecProfile: "H.264 Baseline",
			StreamURL:    fmt.Sprintf("/api/v1/videos/%s/quality/360p", videoID.String()),
		},
		{
			Quality:      "720p",
			Resolution:   "1280x720",
			Bitrate:      2500000,
			Framerate:    30,
			CodecProfile: "H.264 Main",
			StreamURL:    fmt.Sprintf("/api/v1/videos/%s/quality/720p", videoID.String()),
		},
		{
			Quality:      "1080p",
			Resolution:   "1920x1080",
			Bitrate:      5000000,
			Framerate:    60,
			CodecProfile: "H.264 High",
			StreamURL:    fmt.Sprintf("/api/v1/videos/%s/quality/1080p", videoID.String()),
		},
	}

	// Adaptive configuration
	adaptiveConfig := AdaptiveConfig{
		Enabled:         true,
		MinBitrate:      500000,
		MaxBitrate:      8000000,
		StartBitrate:    2500000,
		BufferDuration:  10,
		SwitchThreshold: 3,
	}

	// Generate manifest URL
	manifestURL := fmt.Sprintf("https://cdn.tchat.com/videos/%s/manifest.m3u8", videoID.String())
	if manifestType == "DASH" {
		manifestURL = fmt.Sprintf("https://cdn.tchat.com/videos/%s/manifest.mpd", videoID.String())
	}

	response := ManifestResponse{
		VideoID:        videoID,
		Protocol:       "HLS",
		ManifestURL:    manifestURL,
		ManifestType:   manifestType,
		QualityLevels:  qualityLevels,
		AdaptiveConfig: adaptiveConfig,
	}

	c.JSON(http.StatusOK, response)
}

// GetQualityStream retrieves stream for specific quality
// GET /api/v1/videos/:id/quality/:quality
func (h *StreamingHandler) GetQualityStream(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Get quality parameter
	quality := c.Param("quality")

	// Validate quality
	validQualities := map[string]bool{
		"360p":  true,
		"720p":  true,
		"1080p": true,
	}
	if !validQualities[quality] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quality parameter"})
		return
	}

	// Get video details
	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Check video availability
	if video.AvailabilityStatus != models.StatusPublic {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Video not available"})
		return
	}

	// Generate quality-specific stream URL
	streamURL := fmt.Sprintf("https://cdn.tchat.com/videos/%s/%s/stream.m3u8", videoID.String(), quality)

	c.JSON(http.StatusOK, gin.H{
		"video_id":   videoID,
		"quality":    quality,
		"stream_url": streamURL,
		"expires_at": time.Now().Add(2 * time.Hour),
	})
}

// GetThumbnail retrieves video thumbnail
// GET /api/v1/videos/:id/thumbnail
func (h *StreamingHandler) GetThumbnail(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Get thumbnail size from query
	size := c.DefaultQuery("size", "medium") // small, medium, large

	// Get video details
	video, err := h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Return thumbnail URL
	thumbnailURL := video.ThumbnailURL
	if thumbnailURL == "" {
		thumbnailURL = fmt.Sprintf("https://cdn.tchat.com/videos/%s/thumbnail_%s.jpg", videoID.String(), size)
	}

	c.JSON(http.StatusOK, gin.H{
		"video_id":      videoID,
		"thumbnail_url": thumbnailURL,
		"size":          size,
	})
}

// GetPreviewClip retrieves preview clip URL
// GET /api/v1/videos/:id/preview
func (h *StreamingHandler) GetPreviewClip(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Get video details
	_, err = h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Generate preview clip URL (first 30 seconds)
	previewURL := fmt.Sprintf("https://cdn.tchat.com/videos/%s/preview.mp4", videoID.String())

	c.JSON(http.StatusOK, gin.H{
		"video_id":    videoID,
		"preview_url": previewURL,
		"duration":    30,
		"quality":     "720p",
	})
}

// RecordAnalytics records video playback analytics
// POST /api/v1/videos/:id/analytics
func (h *StreamingHandler) RecordAnalytics(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Parse request body
	var req AnalyticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Record analytics event
	// In production, this would send to analytics service
	// For now, we'll just acknowledge receipt

	c.JSON(http.StatusOK, gin.H{
		"video_id":   videoID,
		"event_type": req.EventType,
		"recorded":   true,
		"timestamp":  time.Now(),
	})
}

// GetPlaybackInfo retrieves current playback session information
// GET /api/v1/videos/:id/playback-info
func (h *StreamingHandler) GetPlaybackInfo(c *gin.Context) {
	// Parse video ID
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID format"})
		return
	}

	// Get session ID from query or header
	sessionIDStr := c.Query("session_id")
	if sessionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
		return
	}

	_, err = uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	// Get video details
	_, err = h.videoService.GetVideo(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Get playback session
// 	session, err := h.videoService.GetPlaybackSession(c.Request.Context(), sessionID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
// 		return
// 	}
// 
// 	// Calculate buffer health
// 	bufferHealth := BufferHealth{
// 		BufferedSeconds:  session.BufferEvents,
// 		BufferPercentage: float64(session.BufferEvents) / float64(video.ContentSpecification.Format.DurationSeconds) * 100,
// 		IsHealthy:        session.BufferEvents >= 5,
// 	}
// 
// 	// Get sync status
// 	syncStates, err := h.syncService.GetSyncStatesBySession(c.Request.Context(), sessionID)
// 	syncedPlatforms := make([]string, 0)
// 	conflictDetected := false
// 	lastSyncTime := time.Now()
// 
// 	if err == nil && len(syncStates) > 0 {
// 		for _, state := range syncStates {
// 			syncedPlatforms = append(syncedPlatforms, string(state.PlatformType))
// 			if state.ConflictState.HasConflicts {
// 				conflictDetected = true
// 			}
// 			if state.LastSyncTime.After(lastSyncTime) {
// 				lastSyncTime = *state.LastSyncTime
// 			}
// 		}
// 	}
// 
// 	syncStatus := SyncStatus{
// 		LastSyncTime:     lastSyncTime,
// 		SyncedPlatforms:  syncedPlatforms,
// 		ConflictDetected: conflictDetected,
// 	}
// 
// 	response := PlaybackInfoResponse{
// 		VideoID:         videoID,
// 		SessionID:       sessionID,
// 		CurrentPosition: session.ProgressTracking.CurrentPosition,
// 		Duration:        video.ContentSpecification.Format.DurationSeconds,
// 		Quality:         session.QualitySettings.CurrentQuality,
// 		Platform:        session.PlatformType,
// 		BufferHealth:    bufferHealth,
// 		SyncStatus:      syncStatus,
// 	}

// 	c.JSON(http.StatusOK, response)
}

// Helper function to parse tags from comma-separated string
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	// Simple split by comma
	tags := make([]string, 0)
	for _, tag := range splitByComma(tagsStr) {
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

// Helper function to split string by comma
func splitByComma(s string) []string {
	result := make([]string, 0)
	current := ""
	for _, char := range s {
		if char == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// Helper function to parse integer from string with default
func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}