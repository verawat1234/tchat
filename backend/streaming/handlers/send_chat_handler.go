package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
	"tchat.dev/streaming/services"
)

type SendChatHandler struct {
	liveStreamRepo   repository.LiveStreamRepositoryInterface
	chatRepo         repository.ChatMessageRepository
	signalingService *services.SignalingService
}

type SendChatRequest struct {
	Message     string  `json:"message" binding:"required,min=1,max=500"`
	MessageType *string `json:"message_type,omitempty"`
}

type ChatMessageResponse struct {
	MessageID        uuid.UUID  `json:"message_id"`
	StreamID         uuid.UUID  `json:"stream_id"`
	SenderID         *uuid.UUID `json:"sender_id"`
	SenderDisplayName string    `json:"sender_display_name"`
	MessageText      string     `json:"message_text"`
	MessageType      string     `json:"message_type"`
	ModerationStatus string     `json:"moderation_status"`
	Timestamp        string     `json:"timestamp"`
	TTLSeconds       int        `json:"ttl_seconds"`
}

type SendChatResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    *SendChatResponseData `json:"data,omitempty"`
}

type SendChatResponseData struct {
	ChatMessage ChatMessageResponse `json:"chat_message"`
}

type SendChatErrorResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func NewSendChatHandler(
	liveStreamRepo repository.LiveStreamRepositoryInterface,
	chatRepo repository.ChatMessageRepository,
	signalingService *services.SignalingService,
) *SendChatHandler {
	return &SendChatHandler{
		liveStreamRepo:   liveStreamRepo,
		chatRepo:         chatRepo,
		signalingService: signalingService,
	}
}

func (h *SendChatHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract stream ID from path parameter
	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, SendChatErrorResponse{
			Success: false,
			Error:   "Invalid stream ID format",
			Details: map[string]interface{}{
				"stream_id": streamIDStr,
			},
		})
		return
	}

	// Verify stream exists and is live
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, SendChatErrorResponse{
			Success: false,
			Error:   "Stream not found",
			Details: map[string]interface{}{
				"stream_id": streamID.String(),
			},
		})
		return
	}

	if stream.Status != "live" {
		endTime := ""
		if stream.EndTime != nil {
			endTime = stream.EndTime.UTC().Format("2006-01-02T15:04:05Z")
		}
		c.JSON(http.StatusBadRequest, SendChatErrorResponse{
			Success: false,
			Error:   "Stream is not live",
			Details: map[string]interface{}{
				"stream_id":      streamID.String(),
				"current_status": stream.Status,
				"end_time":       endTime,
			},
		})
		return
	}

	// Parse request body
	var req SendChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Check if message is empty
		if len(req.Message) == 0 {
			c.JSON(http.StatusBadRequest, SendChatErrorResponse{
				Success: false,
				Error:   "Invalid message content",
				Details: map[string]interface{}{
					"message": "Message cannot be empty and must be between 1 and 500 characters",
				},
			})
			return
		}

		// Check if message is too long
		if len(req.Message) > 500 {
			c.JSON(http.StatusBadRequest, SendChatErrorResponse{
				Success: false,
				Error:   "Invalid message content",
				Details: map[string]interface{}{
					"message":        "Message cannot be empty and must be between 1 and 500 characters",
					"message_length": len(req.Message),
					"max_length":     500,
				},
			})
			return
		}

		c.JSON(http.StatusBadRequest, SendChatErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate message length explicitly
	if len(req.Message) == 0 {
		c.JSON(http.StatusBadRequest, SendChatErrorResponse{
			Success: false,
			Error:   "Invalid message content",
			Details: map[string]interface{}{
				"message": "Message cannot be empty and must be between 1 and 500 characters",
			},
		})
		return
	}

	if len(req.Message) > 500 {
		c.JSON(http.StatusBadRequest, SendChatErrorResponse{
			Success: false,
			Error:   "Invalid message content",
			Details: map[string]interface{}{
				"message":        "Message cannot be empty and must be between 1 and 500 characters",
				"message_length": len(req.Message),
				"max_length":     500,
			},
		})
		return
	}

	// Extract user ID from JWT (optional - can be null for anonymous users)
	var userUUID *uuid.UUID
	var senderDisplayName string
	userID, exists := c.Get("user_id")
	if exists && userID != nil {
		if parsedUserID, ok := userID.(uuid.UUID); ok {
			userUUID = &parsedUserID
			// Get sender display name from context or default
			if displayName, exists := c.Get("display_name"); exists {
				senderDisplayName = displayName.(string)
			} else {
				senderDisplayName = "User"
			}
		}
	} else {
		// Anonymous user
		userUUID = nil
		senderDisplayName = "Anonymous"
	}

	// Rate limiting check (5 messages per second per user)
	if userUUID != nil {
		// TODO: Implement rate limiting using Redis or ScyllaDB
		// For now, this is a placeholder
		// count, err := h.rateLimitCheck(ctx, streamID, *userUUID)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, SendChatErrorResponse{
		// 		Success: false,
		// 		Error:   "Rate limit check failed",
		// 	})
		// 	return
		// }
		// if count >= 5 {
		// 	c.Header("Retry-After", "1")
		// 	c.JSON(http.StatusTooManyRequests, SendChatErrorResponse{
		// 		Success: false,
		// 		Error:   "Rate limit exceeded",
		// 		Details: map[string]interface{}{
		// 			"limit":         5,
		// 			"window":        "1 second",
		// 			"retry_after":   1,
		// 			"current_count": count + 1,
		// 		},
		// 	})
		// 	return
		// }
	}

	// Set default message type
	messageType := "text"
	if req.MessageType != nil && *req.MessageType != "" {
		messageType = *req.MessageType
	}

	// Create chat message
	now := time.Now().UTC()
	chatMessage := &models.ChatMessage{
		StreamID:          streamID,
		Timestamp:         now,
		MessageID:         uuid.New(),
		SenderID:          uuid.Nil,
		SenderDisplayName: senderDisplayName,
		MessageText:       req.Message,
		ModerationStatus:  models.ModerationStatusVisible,
		MessageType:       messageType,
	}

	// Set sender ID if not anonymous
	if userUUID != nil {
		chatMessage.SenderID = *userUUID
	}

	// Persist to ScyllaDB with 30-day TTL
	if err := h.chatRepo.Create(ctx, chatMessage); err != nil {
		c.JSON(http.StatusInternalServerError, SendChatErrorResponse{
			Success: false,
			Error:   "Failed to send chat message",
		})
		return
	}

	// Broadcast to all viewers via WebSocket
	go h.signalingService.BroadcastChatMessage(streamID, *chatMessage)

	// Prepare response
	var responseSenderID *uuid.UUID
	if userUUID != nil {
		responseSenderID = userUUID
	}

	response := SendChatResponse{
		Success: true,
		Message: "Message sent successfully",
		Data: &SendChatResponseData{
			ChatMessage: ChatMessageResponse{
				MessageID:         chatMessage.MessageID,
				StreamID:          streamID,
				SenderID:          responseSenderID,
				SenderDisplayName: senderDisplayName,
				MessageText:       req.Message,
				MessageType:       messageType,
				ModerationStatus:  models.ModerationStatusVisible,
				Timestamp:         now.Format("2006-01-02T15:04:05Z"),
				TTLSeconds:        models.DefaultTTL,
			},
		},
	}

	c.JSON(http.StatusCreated, response)
}