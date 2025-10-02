package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

type GetStreamHandler struct {
	liveStreamRepo    repository.LiveStreamRepositoryInterface
	streamProductRepo repository.StreamProductRepository
}

type FeaturedProductResponse struct {
	ID               uuid.UUID `json:"id"`
	ProductID        uuid.UUID `json:"product_id"`
	ProductName      string    `json:"product_name"`      // Fetched from commerce service
	ProductPrice     float64   `json:"product_price"`     // Fetched from commerce service
	ProductImageURL  string    `json:"product_image_url"` // Fetched from commerce service
	DisplayPosition  string    `json:"display_position"`
	DisplayPriority  int       `json:"display_priority"`
	FeaturedAt       string    `json:"featured_at"`
	ViewCount        int       `json:"view_count"`
	ClickCount       int       `json:"click_count"`
	PurchaseCount    int       `json:"purchase_count"`
	RevenueGenerated float64   `json:"revenue_generated"`
}

type GetStreamResponse struct {
	ID                  uuid.UUID                 `json:"id"`
	BroadcasterID       uuid.UUID                 `json:"broadcaster_id"`
	BroadcasterKYCTier  int                       `json:"broadcaster_kyc_tier"`
	StreamType          string                    `json:"stream_type"`
	Title               string                    `json:"title"`
	Description         *string                   `json:"description"`
	PrivacySetting      string                    `json:"privacy_setting"`
	Status              string                    `json:"status"`
	ScheduledStartTime  *string                   `json:"scheduled_start_time"`
	ActualStartTime     *string                   `json:"actual_start_time"`
	EndTime             *string                   `json:"end_time"`
	RecordingURL        *string                   `json:"recording_url"`
	RecordingExpiryDate *string                   `json:"recording_expiry_date,omitempty"`
	ViewerCount         int                       `json:"viewer_count"`
	PeakViewerCount     int                       `json:"peak_viewer_count"`
	MaxCapacity         int                       `json:"max_capacity"`
	ThumbnailURL        *string                   `json:"thumbnail_url"`
	Language            string                    `json:"language"`
	Tags                []string                  `json:"tags"`
	FeaturedProducts    []FeaturedProductResponse `json:"featured_products"`
	CreatedAt           string                    `json:"created_at"`
	UpdatedAt           string                    `json:"updated_at"`
}

func NewGetStreamHandler(liveStreamRepo repository.LiveStreamRepositoryInterface, streamProductRepo repository.StreamProductRepository) *GetStreamHandler {
	return &GetStreamHandler{
		liveStreamRepo:    liveStreamRepo,
		streamProductRepo: streamProductRepo,
	}
}

func (h *GetStreamHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse stream ID from path parameter
	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stream ID format"})
		return
	}

	// Get stream from repository
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
		return
	}

	// Get featured products for store streams
	var featuredProducts []FeaturedProductResponse
	if stream.StreamType == "store" {
		featuredProducts = h.getFeaturedProducts(ctx, streamID)
	}

	// Build response
	response := h.buildStreamResponse(stream, featuredProducts)

	c.JSON(http.StatusOK, response)
}

func (h *GetStreamHandler) buildStreamResponse(stream *models.LiveStream, featuredProducts []FeaturedProductResponse) GetStreamResponse {
	// Convert sql.NullString to *string
	var description *string
	if stream.Description.Valid {
		description = &stream.Description.String
	}

	var recordingURL *string
	if stream.RecordingURL.Valid {
		recordingURL = &stream.RecordingURL.String
	}

	var thumbnailURL *string
	if stream.ThumbnailURL.Valid {
		thumbnailURL = &stream.ThumbnailURL.String
	}

	// Format timestamps
	var scheduledStartTime, actualStartTime, endTime *string
	if stream.ScheduledStartTime != nil {
		t := stream.ScheduledStartTime.Format("2006-01-02T15:04:05Z")
		scheduledStartTime = &t
	}
	if stream.ActualStartTime != nil {
		t := stream.ActualStartTime.Format("2006-01-02T15:04:05Z")
		actualStartTime = &t
	}
	if stream.EndTime != nil {
		t := stream.EndTime.Format("2006-01-02T15:04:05Z")
		endTime = &t
	}

	var recordingExpiryDate *string
	if stream.RecordingExpiryDate != nil {
		t := stream.RecordingExpiryDate.Format("2006-01-02T15:04:05Z")
		recordingExpiryDate = &t
	}

	// Ensure tags is not nil
	tags := stream.Tags
	if tags == nil {
		tags = []string{}
	}

	// Ensure featured products is not nil
	if featuredProducts == nil {
		featuredProducts = []FeaturedProductResponse{}
	}

	return GetStreamResponse{
		ID:                  stream.ID,
		BroadcasterID:       stream.BroadcasterID,
		BroadcasterKYCTier:  stream.BroadcasterKYCTier,
		StreamType:          stream.StreamType,
		Title:               stream.Title,
		Description:         description,
		PrivacySetting:      stream.PrivacySetting,
		Status:              stream.Status,
		ScheduledStartTime:  scheduledStartTime,
		ActualStartTime:     actualStartTime,
		EndTime:             endTime,
		RecordingURL:        recordingURL,
		RecordingExpiryDate: recordingExpiryDate,
		ViewerCount:         stream.ViewerCount,
		PeakViewerCount:     stream.PeakViewerCount,
		MaxCapacity:         stream.MaxCapacity,
		ThumbnailURL:        thumbnailURL,
		Language:            stream.Language,
		Tags:                tags,
		FeaturedProducts:    featuredProducts,
		CreatedAt:           stream.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:           stream.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (h *GetStreamHandler) getFeaturedProducts(ctx context.Context, streamID uuid.UUID) []FeaturedProductResponse {
	// Get stream products from repository
	products, err := h.streamProductRepo.ListByStream(ctx, streamID)
	if err != nil || len(products) == 0 {
		return []FeaturedProductResponse{}
	}

	// Build featured products response
	// In a real microservice architecture, this would call the commerce service
	// to get product details (name, price, image_url)
	featuredProducts := make([]FeaturedProductResponse, 0, len(products))
	for _, p := range products {
		// TODO: Fetch product details from commerce service
		// For now, return the stream product data with placeholder values
		// The contract test expects these fields to come from the commerce service
		featuredProduct := FeaturedProductResponse{
			ID:               p.ID,
			ProductID:        p.ProductID,
			ProductName:      "Product Name", // Placeholder - should come from commerce service
			ProductPrice:     0.0,             // Placeholder - should come from commerce service
			ProductImageURL:  "",              // Placeholder - should come from commerce service
			DisplayPosition:  p.DisplayPosition,
			DisplayPriority:  p.DisplayPriority,
			FeaturedAt:       p.FeaturedAt.Format("2006-01-02T15:04:05Z"),
			ViewCount:        p.ViewCount,
			ClickCount:       p.ClickCount,
			PurchaseCount:    p.PurchaseCount,
			RevenueGenerated: p.RevenueGenerated,
		}
		featuredProducts = append(featuredProducts, featuredProduct)
	}

	return featuredProducts
}

// Helper function to convert sql.NullString to *string
func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}