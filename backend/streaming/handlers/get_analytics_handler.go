// backend/streaming/handlers/get_analytics_handler.go
// Get Analytics Handler - Retrieves comprehensive analytics for completed or live streams
// Implements T023: Get Stream Analytics contract test

package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

// GetAnalyticsHandler handles GET /api/v1/streams/{streamId}/analytics
type GetAnalyticsHandler struct {
	liveStreamRepo repository.LiveStreamRepositoryInterface
	analyticsRepo  repository.StreamAnalyticsRepository
}

// AnalyticsResponse represents the complete analytics response
type AnalyticsResponse struct {
	StreamID string `json:"stream_id"`

	// Viewer metrics
	TotalUniqueViewers          int `json:"total_unique_viewers"`
	PeakConcurrentViewers       int `json:"peak_concurrent_viewers"`
	AverageWatchDurationSeconds int `json:"average_watch_duration_seconds"`

	// Engagement metrics
	TotalChatMessages int `json:"total_chat_messages"`
	TotalReactions    int `json:"total_reactions"`
	UniqueChatter     int `json:"unique_chatters"`

	// Store context metrics (nullable for video streams)
	ProductsFeatured   *int     `json:"products_featured,omitempty"`
	TotalProductViews  *int     `json:"total_product_views,omitempty"`
	TotalProductClicks *int     `json:"total_product_clicks,omitempty"`
	TotalPurchases     *int     `json:"total_purchases,omitempty"`
	TotalRevenue       *float64 `json:"total_revenue,omitempty"`

	// Quality metrics
	AverageViewerQuality string `json:"average_viewer_quality,omitempty"`
	TotalRebufferEvents  int    `json:"total_rebuffer_events"`

	// Geographic distribution
	ViewerCountries map[string]int `json:"viewer_countries,omitempty"`

	// Calculated timestamp
	CalculatedAt string `json:"calculated_at,omitempty"`
}

// NewGetAnalyticsHandler creates a new get analytics handler
func NewGetAnalyticsHandler(
	liveStreamRepo repository.LiveStreamRepositoryInterface,
	analyticsRepo repository.StreamAnalyticsRepository,
) *GetAnalyticsHandler {
	return &GetAnalyticsHandler{
		liveStreamRepo: liveStreamRepo,
		analyticsRepo:  analyticsRepo,
	}
}

// Handle processes the get analytics request
func (h *GetAnalyticsHandler) Handle(c *gin.Context) {
	ctx := context.Background()

	// Extract broadcaster ID from JWT (stored in context by auth middleware)
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

	// Parse stream ID from URL parameter
	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stream ID format"})
		return
	}

	// Fetch the stream to verify it exists and check broadcaster authorization
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
		return
	}

	// Authorization check: only broadcaster can view analytics
	if stream.BroadcasterID != broadcasterUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only broadcaster can view analytics"})
		return
	}

	// Fetch analytics from repository
	analytics, err := h.analyticsRepo.GetByStreamID(ctx, streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
		return
	}

	// Convert analytics to response format
	response := h.buildAnalyticsResponse(analytics)

	c.JSON(http.StatusOK, response)
}

// buildAnalyticsResponse converts models.StreamAnalytics to AnalyticsResponse
func (h *GetAnalyticsHandler) buildAnalyticsResponse(analytics *models.StreamAnalytics) AnalyticsResponse {
	response := AnalyticsResponse{
		StreamID:                    analytics.StreamID.String(),
		TotalUniqueViewers:          analytics.TotalUniqueViewers,
		PeakConcurrentViewers:       analytics.PeakConcurrentViewers,
		AverageWatchDurationSeconds: analytics.AverageWatchDurationSeconds,
		TotalChatMessages:           analytics.TotalChatMessages,
		TotalReactions:              analytics.TotalReactions,
		UniqueChatter:               analytics.UniqueChatter,
		ProductsFeatured:            analytics.ProductsFeatured,
		TotalProductViews:           analytics.TotalProductViews,
		TotalProductClicks:          analytics.TotalProductClicks,
		TotalPurchases:              analytics.TotalPurchases,
		TotalRevenue:                analytics.TotalRevenue,
		TotalRebufferEvents:         analytics.TotalRebufferEvents,
	}

	// Handle average viewer quality
	if analytics.AverageViewerQuality.Valid {
		response.AverageViewerQuality = analytics.AverageViewerQuality.String
	}

	// Parse viewer countries from JSONB
	if len(analytics.ViewerCountries) > 0 {
		var countries map[string]int
		if err := json.Unmarshal(analytics.ViewerCountries, &countries); err == nil {
			response.ViewerCountries = countries
		}
	}

	// Format calculated timestamp
	if analytics.CalculatedAt != nil {
		response.CalculatedAt = analytics.CalculatedAt.UTC().Format("2006-01-02T15:04:05Z")
	}

	return response
}
