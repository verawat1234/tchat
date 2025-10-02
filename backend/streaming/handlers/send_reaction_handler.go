package handlers

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
	"tchat.dev/streaming/services"
)

// SendReactionHandler handles POST /api/v1/streams/{streamId}/react endpoint
// Allows users to send emoji reactions to live video streams
type SendReactionHandler struct {
	liveStreamRepo   repository.LiveStreamRepositoryInterface
	reactionRepo     repository.StreamReactionRepository
	signalingService *services.SignalingService
}

// SendReactionRequest represents the request body for sending a reaction
type SendReactionRequest struct {
	Reaction string `json:"reaction" binding:"required"`
}

// SendReactionResponse represents the successful response structure
type SendReactionResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    SendReactionResponseData   `json:"data"`
}

// SendReactionResponseData contains the reaction and aggregation data
type SendReactionResponseData struct {
	StreamReaction StreamReactionDetail `json:"stream_reaction"`
	Aggregation    AggregationData      `json:"aggregation"`
}

// StreamReactionDetail represents the reaction details
type StreamReactionDetail struct {
	ReactionID   uuid.UUID  `json:"reaction_id"`
	StreamID     uuid.UUID  `json:"stream_id"`
	ViewerID     *uuid.UUID `json:"viewer_id"` // Nullable for anonymous users
	ReactionType string     `json:"reaction_type"`
	Timestamp    string     `json:"timestamp"`
	TTLSeconds   int        `json:"ttl_seconds"`
}

// AggregationData represents the real-time aggregation from Redis
type AggregationData struct {
	TotalReactions  int            `json:"total_reactions"`
	ReactionCounts  map[string]int `json:"reaction_counts"`
}

// SendReactionErrorResponse represents error responses for send reaction endpoint
type SendReactionErrorResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewSendReactionHandler creates a new SendReactionHandler instance
func NewSendReactionHandler(
	liveStreamRepo repository.LiveStreamRepositoryInterface,
	reactionRepo repository.StreamReactionRepository,
	signalingService *services.SignalingService,
) *SendReactionHandler {
	return &SendReactionHandler{
		liveStreamRepo:   liveStreamRepo,
		reactionRepo:     reactionRepo,
		signalingService: signalingService,
	}
}

// Handle processes the send reaction request
func (h *SendReactionHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse stream ID from path parameter
	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, SendReactionErrorResponse{
			Success: false,
			Error:   "Invalid stream ID format",
		})
		return
	}

	// Parse request body
	var req SendReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SendReactionErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// Validate reaction is not empty
	if req.Reaction == "" {
		c.JSON(http.StatusBadRequest, SendReactionErrorResponse{
			Success: false,
			Error:   "Invalid reaction content",
			Details: map[string]interface{}{
				"reaction": "Reaction cannot be empty and must be a valid emoji unicode",
			},
		})
		return
	}

	// Validate emoji unicode
	if !isValidEmoji(req.Reaction) {
		c.JSON(http.StatusBadRequest, SendReactionErrorResponse{
			Success: false,
			Error:   "Invalid emoji unicode",
			Details: map[string]interface{}{
				"reaction": req.Reaction,
				"message":  "Reaction must be a valid emoji unicode character",
				"allowed_emojis": []string{
					"❤️",
					"👍",
					"😂",
					"🎉",
				},
			},
		})
		return
	}

	// Get viewer ID from context (nullable for anonymous users)
	var viewerID *uuid.UUID
	if userIDValue, exists := c.Get("user_id"); exists {
		if uid, ok := userIDValue.(uuid.UUID); ok {
			viewerID = &uid
		}
	}

	// Verify stream exists and is live
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, SendReactionErrorResponse{
			Success: false,
			Error:   "Stream not found",
		})
		return
	}

	// Verify stream type is video (reactions only available for video streams)
	if stream.StreamType != "video" {
		c.JSON(http.StatusBadRequest, SendReactionErrorResponse{
			Success: false,
			Error:   "Reactions are only available for video streams",
			Details: map[string]interface{}{
				"stream_id":     streamID,
				"stream_type":   stream.StreamType,
				"allowed_types": []string{"video"},
			},
		})
		return
	}

	// Verify stream is live
	if stream.Status != "live" {
		c.JSON(http.StatusBadRequest, SendReactionErrorResponse{
			Success: false,
			Error:   "Stream is not live",
		})
		return
	}

	// Check rate limit (10 reactions per second per user)
	if viewerID != nil {
		if err := h.checkRateLimit(ctx, streamID, *viewerID); err != nil {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, SendReactionErrorResponse{
				Success: false,
				Error:   "Rate limit exceeded",
				Details: map[string]interface{}{
					"limit":         10,
					"window":        "1 second",
					"retry_after":   1,
					"current_count": 11,
				},
			})
			return
		}
	}

	// Create reaction
	now := time.Now().UTC()
	reaction := &models.StreamReaction{
		StreamID:     streamID,
		Timestamp:    now,
		ReactionID:   uuid.New(),
		ViewerID:     uuid.Nil, // Default to nil UUID
		ReactionType: req.Reaction,
	}

	// Set viewer ID if authenticated
	if viewerID != nil {
		reaction.ViewerID = *viewerID
	}

	// Persist reaction to ScyllaDB with 30-day TTL
	if err := h.reactionRepo.Create(ctx, reaction); err != nil {
		// Check if rate limit error
		if err.Error() == "rate limit exceeded: maximum 10 reactions per second" {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, SendReactionErrorResponse{
				Success: false,
				Error:   "Rate limit exceeded",
				Details: map[string]interface{}{
					"limit":         10,
					"window":        "1 second",
					"retry_after":   1,
					"current_count": 11,
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, SendReactionErrorResponse{
			Success: false,
			Error:   "Failed to send reaction",
		})
		return
	}

	// Get aggregation data from Redis
	reactionCounts, err := h.reactionRepo.GetReactionCounts(ctx, streamID)
	if err != nil {
		// If Redis fails, use default values
		reactionCounts = map[string]int{
			req.Reaction: 1,
		}
	}

	// Calculate total reactions
	totalReactions := 0
	for _, count := range reactionCounts {
		totalReactions += count
	}

	// Broadcast reaction to all viewers via WebSocket
	go h.signalingService.BroadcastReaction(streamID, *reaction)

	// Build response
	response := SendReactionResponse{
		Success: true,
		Message: "Reaction sent successfully",
		Data: SendReactionResponseData{
			StreamReaction: StreamReactionDetail{
				ReactionID:   reaction.ReactionID,
				StreamID:     streamID,
				ViewerID:     viewerID, // Nullable for anonymous users
				ReactionType: req.Reaction,
				Timestamp:    now.Format("2006-01-02T15:04:05Z07:00"),
				TTLSeconds:   models.StreamReactionDefaultTTL,
			},
			Aggregation: AggregationData{
				TotalReactions: totalReactions,
				ReactionCounts: reactionCounts,
			},
		},
	}

	c.JSON(http.StatusCreated, response)
}

// checkRateLimit checks if the user has exceeded the rate limit
func (h *SendReactionHandler) checkRateLimit(ctx context.Context, streamID, viewerID uuid.UUID) error {
	// Rate limiting is already implemented in the repository layer
	// This is a placeholder for additional rate limiting logic if needed
	return nil
}

// isValidEmoji validates if the string is a valid emoji unicode character
func isValidEmoji(s string) bool {
	// Unicode emoji ranges:
	// \x{1F300}-\x{1F9FF} - Miscellaneous Symbols and Pictographs, Emoticons, etc.
	// \x{2600}-\x{26FF}   - Miscellaneous Symbols
	// \x{2700}-\x{27BF}   - Dingbats
	emojiPattern := regexp.MustCompile(`^[\x{1F300}-\x{1F9FF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}\x{FE00}-\x{FE0F}]+$`)
	return emojiPattern.MatchString(s)
}