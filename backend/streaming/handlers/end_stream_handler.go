package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
	"tchat.dev/streaming/services"
)

type EndStreamHandler struct {
	liveStreamRepo      repository.LiveStreamRepositoryInterface
	streamAnalyticsRepo repository.StreamAnalyticsRepository
	recordingService    services.RecordingService
}

type EndStreamResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type EndStreamErrorResponse struct {
	Success       bool                   `json:"success"`
	Error         string                 `json:"error"`
	StreamID      string                 `json:"stream_id,omitempty"`
	BroadcasterID string                 `json:"broadcaster_id,omitempty"`
	CurrentStatus string                 `json:"current_status,omitempty"`
	Details       map[string]interface{} `json:"details,omitempty"`
}

func NewEndStreamHandler(
	liveStreamRepo repository.LiveStreamRepositoryInterface,
	streamAnalyticsRepo repository.StreamAnalyticsRepository,
	recordingService services.RecordingService,
) *EndStreamHandler {
	return &EndStreamHandler{
		liveStreamRepo:      liveStreamRepo,
		streamAnalyticsRepo: streamAnalyticsRepo,
		recordingService:    recordingService,
	}
}

func (h *EndStreamHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract broadcaster ID from JWT
	broadcasterID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	broadcasterUUID, ok := broadcasterID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stream ID format"})
		return
	}

	// Fetch stream
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, EndStreamErrorResponse{
			Success:  false,
			Error:    "Stream not found",
			StreamID: streamID.String(),
			Details: map[string]interface{}{
				"message": "The requested stream does not exist or has been deleted",
			},
		})
		return
	}

	// Authorization check - only broadcaster can end their stream
	if stream.BroadcasterID != broadcasterUUID {
		c.JSON(http.StatusForbidden, EndStreamErrorResponse{
			Success:       false,
			Error:         "Only the broadcaster can end their stream",
			StreamID:      streamID.String(),
			BroadcasterID: stream.BroadcasterID.String(),
			Details: map[string]interface{}{
				"message":    "This stream belongs to another user",
				"permission": "broadcaster_only",
			},
		})
		return
	}

	// Status validation - stream must be live
	if stream.Status != "live" {
		c.JSON(http.StatusConflict, EndStreamErrorResponse{
			Success:       false,
			Error:         "Cannot end stream that is not live",
			CurrentStatus: stream.Status,
			StreamID:      streamID.String(),
			Details: map[string]interface{}{
				"message":          "Stream must be in 'live' status to be ended",
				"allowed_statuses": []string{"live"},
			},
		})
		return
	}

	// Stop recording and get URL
	recordingURL, err := h.recordingService.StopRecording(streamID)
	if err != nil {
		// Log error but continue - recording is not critical
		fmt.Printf("Warning: Failed to stop recording for stream %s: %v\n", streamID, err)
		recordingURL = ""
	}

	// Calculate stream metrics
	now := time.Now().UTC()
	var durationSeconds int
	if stream.ActualStartTime != nil {
		durationSeconds = int(now.Sub(*stream.ActualStartTime).Seconds())
	}

	// Set recording expiry (30 days from now)
	recordingExpiryDate := now.Add(30 * 24 * time.Hour)

	// Update stream status to ended
	updates := map[string]interface{}{
		"status":     "ended",
		"end_time":   now,
		"updated_at": now,
	}

	if recordingURL != "" {
		updates["recording_url"] = recordingURL
		updates["recording_expiry_date"] = recordingExpiryDate

		// Set CDN lifecycle policy
		if err := h.recordingService.SetLifecyclePolicy(recordingURL, recordingExpiryDate); err != nil {
			// Log error but continue
			fmt.Printf("Warning: Failed to set lifecycle policy for %s: %v\n", recordingURL, err)
		}
	}

	if err := h.liveStreamRepo.Update(streamID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stream status"})
		return
	}

	// Create analytics summary in background
	go h.calculateStreamAnalytics(context.Background(), streamID, stream, durationSeconds)

	// Build response
	response := h.buildSuccessResponse(stream, recordingURL, recordingExpiryDate, durationSeconds, now)
	c.JSON(http.StatusOK, response)
}

func (h *EndStreamHandler) buildSuccessResponse(
	stream *models.LiveStream,
	recordingURL string,
	recordingExpiryDate time.Time,
	durationSeconds int,
	endTime time.Time,
) EndStreamResponse {
	// Format times in RFC3339 format
	scheduledStartTime := ""
	if stream.ScheduledStartTime != nil {
		scheduledStartTime = stream.ScheduledStartTime.Format(time.RFC3339)
	}

	actualStartTime := ""
	if stream.ActualStartTime != nil {
		actualStartTime = stream.ActualStartTime.Format(time.RFC3339)
	}

	endTimeStr := endTime.Format(time.RFC3339)
	recordingExpiryDateStr := recordingExpiryDate.Format(time.RFC3339)

	// Build stream object
	streamObj := map[string]interface{}{
		"id":                     stream.ID.String(),
		"status":                 "ended",
		"broadcaster_id":         stream.BroadcasterID.String(),
		"stream_type":            stream.StreamType,
		"title":                  stream.Title,
		"privacy_setting":        stream.PrivacySetting,
		"scheduled_start_time":   scheduledStartTime,
		"actual_start_time":      actualStartTime,
		"end_time":               endTimeStr,
		"recording_url":          recordingURL,
		"recording_expiry_date":  recordingExpiryDateStr,
		"viewer_count":           0, // Stream has ended, so viewer count is 0
		"peak_viewer_count":      stream.PeakViewerCount,
	}

	// Build recording metadata
	recordingMetadata := map[string]interface{}{
		"processing_status": "processing",
		"duration_seconds":  durationSeconds,
		"file_size_bytes":   0, // Not yet available during processing
		"retention_days":    30,
	}

	return EndStreamResponse{
		Success: true,
		Message: "Stream ended successfully",
		Data: map[string]interface{}{
			"stream":             streamObj,
			"recording_metadata": recordingMetadata,
		},
	}
}

func (h *EndStreamHandler) calculateStreamAnalytics(ctx context.Context, streamID uuid.UUID, stream *models.LiveStream, durationSeconds int) {
	// This would calculate comprehensive analytics in a real implementation
	// For now, we create a basic analytics record

	analytics := &models.StreamAnalytics{
		StreamID:              streamID,
		PeakConcurrentViewers: stream.PeakViewerCount,
		// Additional analytics would be calculated here based on viewer sessions, chat, reactions, etc.
	}

	if err := h.streamAnalyticsRepo.Create(ctx, analytics); err != nil {
		fmt.Printf("Warning: Failed to create analytics for stream %s: %v\n", streamID, err)
	}
}