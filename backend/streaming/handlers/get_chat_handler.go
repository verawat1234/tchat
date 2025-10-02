package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/repository"
)

type GetChatHandler struct {
	liveStreamRepo repository.LiveStreamRepository
	chatRepo       repository.ChatMessageRepository
}

type ChatMessageResponse struct {
	MessageID        uuid.UUID `json:"message_id"`
	StreamID         uuid.UUID `json:"stream_id"`
	UserID           uuid.UUID `json:"user_id"`
	Message          string    `json:"message"`
	Timestamp        string    `json:"timestamp"`
	ModerationStatus string    `json:"moderation_status"`
}

type GetChatResponse struct {
	Messages   []ChatMessageResponse `json:"messages"`
	HasMore    bool                  `json:"has_more"`
	NextCursor *string               `json:"next_cursor"`
}

func NewGetChatHandler(
	liveStreamRepo repository.LiveStreamRepository,
	chatRepo repository.ChatMessageRepository,
) *GetChatHandler {
	return &GetChatHandler{
		liveStreamRepo: liveStreamRepo,
		chatRepo:       chatRepo,
	}
}

func (h *GetChatHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stream ID format"})
		return
	}

	// Verify stream exists
	_, err = h.liveStreamRepo.GetByID(ctx, streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
		return
	}

	// Parse pagination parameters
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	var beforeTimestamp *time.Time
	if beforeStr := c.Query("before_timestamp"); beforeStr != "" {
		if t, err := time.Parse(time.RFC3339, beforeStr); err == nil {
			beforeTimestamp = &t
		}
	}

	// Fetch messages from ScyllaDB (reverse chronological order)
	messages, err := h.chatRepo.ListByStream(ctx, streamID, limit+1, beforeTimestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat messages"})
		return
	}

	// Determine if there are more messages
	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	// Convert to response format
	responseMessages := make([]ChatMessageResponse, len(messages))
	for i, msg := range messages {
		responseMessages[i] = ChatMessageResponse{
			MessageID:        msg.MessageID,
			StreamID:         msg.StreamID,
			UserID:           msg.UserID,
			Message:          msg.Message,
			Timestamp:        msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			ModerationStatus: msg.ModerationStatus,
		}
	}

	// Generate next cursor if there are more messages
	var nextCursor *string
	if hasMore && len(messages) > 0 {
		cursor := messages[len(messages)-1].Timestamp.Format(time.RFC3339)
		nextCursor = &cursor
	}

	response := GetChatResponse{
		Messages:   responseMessages,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}

	c.JSON(http.StatusOK, response)
}