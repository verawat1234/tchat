package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
	"tchat.dev/streaming/services"
)

type CreateStreamHandler struct {
	liveStreamRepo repository.LiveStreamRepository
	kycService     services.KYCService
	webrtcService  services.WebRTCService
}

type CreateStreamRequest struct {
	StreamType         string                 `json:"stream_type" binding:"required,oneof=store video"`
	Title              string                 `json:"title" binding:"required,min=1,max=200"`
	Description        *string                `json:"description,omitempty"`
	ScheduledStartTime *time.Time             `json:"scheduled_start_time,omitempty"`
	MaxCapacity        *int                   `json:"max_capacity,omitempty"`
	StreamSettings     map[string]interface{} `json:"stream_settings,omitempty"`
}

type CreateStreamResponse struct {
	ID                  uuid.UUID              `json:"id"`
	StreamType          string                 `json:"stream_type"`
	Title               string                 `json:"title"`
	Description         *string                `json:"description"`
	BroadcasterID       uuid.UUID              `json:"broadcaster_id"`
	BroadcasterKYCTier  int                    `json:"broadcaster_kyc_tier"`
	Status              string                 `json:"status"`
	StreamKey           string                 `json:"stream_key"`
	MaxCapacity         int                    `json:"max_capacity"`
	ScheduledStartTime  *time.Time             `json:"scheduled_start_time"`
	CreatedAt           time.Time              `json:"created_at"`
}

func NewCreateStreamHandler(
	liveStreamRepo repository.LiveStreamRepository,
	kycService services.KYCService,
	webrtcService services.WebRTCService,
) *CreateStreamHandler {
	return &CreateStreamHandler{
		liveStreamRepo: liveStreamRepo,
		kycService:     kycService,
		webrtcService:  webrtcService,
	}
}

func (h *CreateStreamHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract broadcaster ID from JWT (assumes auth middleware)
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

	var req CreateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// KYC validation based on stream type
	var kycValid bool
	var err error
	if req.StreamType == "store" {
		kycValid, err = h.kycService.ValidateStoreSellerKYC(broadcasterUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "KYC validation failed"})
			return
		}
		if !kycValid {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "KYC Tier 1+ verification required for store streaming",
				"required_tier": 1,
			})
			return
		}
	} else if req.StreamType == "video" {
		kycValid, err = h.kycService.ValidateVideoCreatorAuth(broadcasterUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication validation failed"})
			return
		}
		if !kycValid {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Email or phone verification required for video streaming",
			})
			return
		}
	}

	// Get broadcaster KYC tier
	kycTier, err := h.kycService.GetKYCTier(broadcasterUUID)
	if err != nil {
		kycTier = 0 // Default to Tier 0 if lookup fails
	}

	// Generate stream key
	streamKey := generateStreamKey()

	// Set default max capacity
	maxCapacity := 50000
	if req.MaxCapacity != nil && *req.MaxCapacity > 0 && *req.MaxCapacity <= 100000 {
		maxCapacity = *req.MaxCapacity
	}

	// Create LiveStream entity
	stream := &models.LiveStream{
		BroadcasterID:      broadcasterUUID,
		StreamType:         req.StreamType,
		Title:              req.Title,
		Description:        req.Description,
		BroadcasterKYCTier: kycTier,
		Status:             "scheduled",
		StreamKey:          streamKey,
		MaxCapacity:        maxCapacity,
		ViewerCount:        0,
		PeakViewerCount:    0,
		TotalViewTime:      0,
		ScheduledStartTime: req.ScheduledStartTime,
		StreamSettings:     req.StreamSettings,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	// Persist to database
	if err := h.liveStreamRepo.Create(ctx, stream); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stream"})
		return
	}

	// Return response
	response := CreateStreamResponse{
		ID:                 stream.ID,
		StreamType:         stream.StreamType,
		Title:              stream.Title,
		Description:        stream.Description,
		BroadcasterID:      stream.BroadcasterID,
		BroadcasterKYCTier: stream.BroadcasterKYCTier,
		Status:             stream.Status,
		StreamKey:          stream.StreamKey,
		MaxCapacity:        stream.MaxCapacity,
		ScheduledStartTime: stream.ScheduledStartTime,
		CreatedAt:          stream.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func generateStreamKey() string {
	return "rtmp://stream.tchat.com/" + uuid.New().String()[:16]
}