// backend/streaming/handlers/list_streams_handler.go
// List Streams Handler - Retrieves live streams with filtering and pagination
// Implements contract test: TestListStreamsPactConsumer

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

// ListStreamsHandler handles GET /api/v1/streams endpoint
type ListStreamsHandler struct {
	liveStreamRepo repository.LiveStreamRepositoryInterface
}

// ListStreamsResponse defines the response structure for stream listing
type ListStreamsResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    ListStreamsDataPayload `json:"data"`
}

// ListStreamsDataPayload contains the actual stream data and pagination metadata
type ListStreamsDataPayload struct {
	Streams []StreamSummary `json:"streams"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
}

// StreamSummary provides a summary view of a live stream
// Matches the contract test expectations with all required fields
type StreamSummary struct {
	// Core identifiers
	ID            uuid.UUID `json:"id"`
	BroadcasterID uuid.UUID `json:"broadcaster_id"`

	// KYC and stream metadata
	BroadcasterKYCTier int    `json:"broadcaster_kyc_tier"`
	StreamType         string `json:"stream_type"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	PrivacySetting     string `json:"privacy_setting"`
	Status             string `json:"status"`

	// Stream metrics
	ViewerCount     int `json:"viewer_count"`
	PeakViewerCount int `json:"peak_viewer_count"`
	MaxCapacity     int `json:"max_capacity"`

	// Timestamps (RFC3339 format)
	ScheduledStartTime string  `json:"scheduled_start_time"`
	ActualStartTime    *string `json:"actual_start_time,omitempty"`
	EndTime            *string `json:"end_time,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`

	// Optional fields
	ThumbnailURL         *string  `json:"thumbnail_url,omitempty"`
	Language             string   `json:"language"`
	Tags                 []string `json:"tags"`
	RecordingURL         *string  `json:"recording_url,omitempty"`
	RecordingExpiryDate  *string  `json:"recording_expiry_date,omitempty"`
}

// NewListStreamsHandler creates a new list streams handler instance
func NewListStreamsHandler(liveStreamRepo repository.LiveStreamRepositoryInterface) *ListStreamsHandler {
	return &ListStreamsHandler{liveStreamRepo: liveStreamRepo}
}

// Handle processes the GET /api/v1/streams request with filtering and pagination
func (h *ListStreamsHandler) Handle(c *gin.Context) {
	// Parse pagination parameters with validation
	limit := 20 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0 // default offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Build filters map for repository query
	filters := make(map[string]interface{})

	// Filter by stream_type (store/video)
	if streamType := c.Query("stream_type"); streamType != "" {
		if streamType == "store" || streamType == "video" {
			filters["stream_type"] = streamType
		}
	}

	// Filter by status (scheduled/live/ended)
	if status := c.Query("status"); status != "" {
		if status == "scheduled" || status == "live" || status == "ended" || status == "terminated" {
			filters["status"] = status
		}
	}

	// Filter by broadcaster_id (UUID validation)
	if broadcasterIDStr := c.Query("broadcaster_id"); broadcasterIDStr != "" {
		if broadcasterID, err := uuid.Parse(broadcasterIDStr); err == nil {
			filters["broadcaster_id"] = broadcasterID
		}
	}

	// Fetch streams from repository with filters and pagination
	streams, totalCount, err := h.liveStreamRepo.List(filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve streams",
			"error":   err.Error(),
		})
		return
	}

	// Convert database models to API response format
	summaries := make([]StreamSummary, len(streams))
	for i, stream := range streams {
		summaries[i] = convertToStreamSummary(stream)
	}

	// Build response with pagination metadata
	response := ListStreamsResponse{
		Success: true,
		Message: "Streams retrieved successfully",
		Data: ListStreamsDataPayload{
			Streams: summaries,
			Total:   totalCount,
			Limit:   limit,
			Offset:  offset,
		},
	}

	c.JSON(http.StatusOK, response)
}

// convertToStreamSummary transforms a database LiveStream model to API StreamSummary
func convertToStreamSummary(stream *models.LiveStream) StreamSummary {
	summary := StreamSummary{
		// Core identifiers
		ID:            stream.ID,
		BroadcasterID: stream.BroadcasterID,

		// KYC and stream metadata
		BroadcasterKYCTier: stream.BroadcasterKYCTier,
		StreamType:         stream.StreamType,
		Title:              stream.Title,
		PrivacySetting:     stream.PrivacySetting,
		Status:             stream.Status,

		// Stream metrics
		ViewerCount:     stream.ViewerCount,
		PeakViewerCount: stream.PeakViewerCount,
		MaxCapacity:     stream.MaxCapacity,

		// Language and tags
		Language: stream.Language,
		Tags:     stream.Tags,

		// Required timestamps (always present)
		CreatedAt: stream.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: stream.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle nullable description
	if stream.Description.Valid {
		summary.Description = stream.Description.String
	} else {
		summary.Description = ""
	}

	// Handle optional scheduled start time
	if stream.ScheduledStartTime != nil {
		summary.ScheduledStartTime = stream.ScheduledStartTime.Format("2006-01-02T15:04:05Z07:00")
	} else {
		summary.ScheduledStartTime = ""
	}

	// Handle optional actual start time
	if stream.ActualStartTime != nil {
		actualStart := stream.ActualStartTime.Format("2006-01-02T15:04:05Z07:00")
		summary.ActualStartTime = &actualStart
	}

	// Handle optional end time
	if stream.EndTime != nil {
		endTime := stream.EndTime.Format("2006-01-02T15:04:05Z07:00")
		summary.EndTime = &endTime
	}

	// Handle optional thumbnail URL
	if stream.ThumbnailURL.Valid {
		summary.ThumbnailURL = &stream.ThumbnailURL.String
	}

	// Handle optional recording URL
	if stream.RecordingURL.Valid {
		summary.RecordingURL = &stream.RecordingURL.String
	}

	// Handle optional recording expiry date
	if stream.RecordingExpiryDate != nil {
		expiryDate := stream.RecordingExpiryDate.Format("2006-01-02T15:04:05Z07:00")
		summary.RecordingExpiryDate = &expiryDate
	}

	return summary
}