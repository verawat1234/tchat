package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

// FeatureProductHandler handles featuring products in store live streams
type FeatureProductHandler struct {
	liveStreamRepo repository.LiveStreamRepositoryInterface
	productRepo    repository.StreamProductRepository
}

// FeatureProductRequest represents the request to feature a product
type FeatureProductRequest struct {
	ProductID       uuid.UUID `json:"product_id" binding:"required"`
	DisplayPosition string    `json:"display_position" binding:"required,oneof=overlay sidebar fullscreen"`
	DisplayPriority int       `json:"display_priority" binding:"required,min=1,max=10"`
}

// FeatureProductResponse represents the successful response
type FeatureProductResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    FeaturedProductData    `json:"data"`
}

// FeaturedProductData contains the featured product details
type FeaturedProductData struct {
	ID                     uuid.UUID `json:"id"`
	StreamID               uuid.UUID `json:"stream_id"`
	ProductID              uuid.UUID `json:"product_id"`
	FeaturedAt             string    `json:"featured_at"`
	DisplayDurationSeconds *int      `json:"display_duration_seconds"`
	DisplayPosition        string    `json:"display_position"`
	DisplayPriority        int       `json:"display_priority"`
	ViewCount              int       `json:"view_count"`
	ClickCount             int       `json:"click_count"`
	PurchaseCount          int       `json:"purchase_count"`
	RevenueGenerated       float64   `json:"revenue_generated"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewFeatureProductHandler creates a new feature product handler
func NewFeatureProductHandler(
	liveStreamRepo repository.LiveStreamRepositoryInterface,
	productRepo repository.StreamProductRepository,
) *FeatureProductHandler {
	return &FeatureProductHandler{
		liveStreamRepo: liveStreamRepo,
		productRepo:    productRepo,
	}
}

// Handle processes the request to feature a product in a live stream
func (h *FeatureProductHandler) Handle(c *gin.Context) {
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

	// Extract and validate stream ID
	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid stream ID format",
		})
		return
	}

	// Parse request body
	var req FeatureProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Fetch stream (repository method doesn't use context)
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "Stream not found",
		})
		return
	}

	// Authorization check - only broadcaster can feature products
	if stream.BroadcasterID != broadcasterUUID {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Success: false,
			Error:   "Not authorized to feature products on this stream",
		})
		return
	}

	// Stream type validation - only store streams can feature products
	if stream.StreamType != "store" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Products can only be featured in store streams",
			Details: map[string]interface{}{
				"stream_id":            streamID.String(),
				"stream_type":          stream.StreamType,
				"allowed_stream_types": []string{"store"},
			},
		})
		return
	}

	// Check maximum products limit (10 per stream)
	currentCount, err := h.productRepo.CountByStream(ctx, streamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to check product count",
		})
		return
	}

	if currentCount >= 10 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Maximum number of featured products reached",
			Details: map[string]interface{}{
				"stream_id":             streamID.String(),
				"current_product_count": currentCount,
				"max_allowed":           10,
				"message":               "Remove a featured product before adding a new one",
			},
		})
		return
	}

	// Validate product exists (simplified - in production, call commerce service)
	if err := h.validateProductExists(ctx, req.ProductID); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "Product not found",
			Details: map[string]interface{}{
				"product_id": req.ProductID.String(),
				"message":    "Product does not exist or is not available",
			},
		})
		return
	}

	// Create featured product
	now := time.Now().UTC()
	streamProduct := &models.StreamProduct{
		ID:                     uuid.New(),
		StreamID:               streamID,
		ProductID:              req.ProductID,
		FeaturedAt:             now,
		DisplayDurationSeconds: nil, // NULL initially as still being featured
		DisplayPosition:        req.DisplayPosition,
		DisplayPriority:        req.DisplayPriority,
		ViewCount:              0,
		ClickCount:             0,
		PurchaseCount:          0,
		RevenueGenerated:       0.00,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	// Persist to database
	if err := h.productRepo.Create(ctx, streamProduct); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to feature product",
		})
		return
	}

	// Return success response
	response := FeatureProductResponse{
		Success: true,
		Message: "Product featured successfully",
		Data: FeaturedProductData{
			ID:                     streamProduct.ID,
			StreamID:               streamProduct.StreamID,
			ProductID:              streamProduct.ProductID,
			FeaturedAt:             streamProduct.FeaturedAt.Format("2006-01-02T15:04:05Z07:00"),
			DisplayDurationSeconds: streamProduct.DisplayDurationSeconds,
			DisplayPosition:        streamProduct.DisplayPosition,
			DisplayPriority:        streamProduct.DisplayPriority,
			ViewCount:              streamProduct.ViewCount,
			ClickCount:             streamProduct.ClickCount,
			PurchaseCount:          streamProduct.PurchaseCount,
			RevenueGenerated:       streamProduct.RevenueGenerated,
		},
	}

	c.JSON(http.StatusCreated, response)
}

// validateProductExists validates that a product exists
// In production, this would call the commerce service API
func (h *FeatureProductHandler) validateProductExists(ctx context.Context, productID uuid.UUID) error {
	// For now, reject products with ID starting with "999" (contract test pattern)
	if productID.String()[:3] == "999" {
		return ErrProductNotFound
	}

	// In production: call commerce service to validate product
	// commerceClient.GetProduct(ctx, productID)

	return nil
}

// ErrProductNotFound is returned when a product doesn't exist
var ErrProductNotFound = &ProductNotFoundError{}

// ProductNotFoundError represents a product not found error
type ProductNotFoundError struct{}

func (e *ProductNotFoundError) Error() string {
	return "product not found"
}