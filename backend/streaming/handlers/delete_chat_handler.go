package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

type DeleteChatHandler struct {
	liveStreamRepo repository.LiveStreamRepositoryInterface
	chatRepo       repository.ChatMessageRepository
}

func NewDeleteChatHandler(
	liveStreamRepo repository.LiveStreamRepositoryInterface,
	chatRepo repository.ChatMessageRepository,
) *DeleteChatHandler {
	return &DeleteChatHandler{
		liveStreamRepo: liveStreamRepo,
		chatRepo:       chatRepo,
	}
}

func (h *DeleteChatHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract user ID from JWT middleware (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Missing or invalid JWT token",
		})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Missing or invalid JWT token",
		})
		return
	}

	// Parse stream ID
	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid stream ID format",
		})
		return
	}

	// Parse message ID
	messageIDStr := c.Param("messageId")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid message ID format",
		})
		return
	}

	// Verify stream exists
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Stream not found",
			"details": gin.H{
				"stream_id": streamID,
			},
		})
		return
	}

	// Fetch the message to verify it exists
	message, err := h.chatRepo.GetByID(ctx, streamID, messageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Chat message not found",
			"details": gin.H{
				"stream_id":  streamID,
				"message_id": messageID,
			},
		})
		return
	}

	// Check if message is already deleted/removed
	if message.ModerationStatus == models.ModerationStatusRemoved {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Chat message not found or already removed",
			"details": gin.H{
				"stream_id":         streamID,
				"message_id":        messageID,
				"moderation_status": message.ModerationStatus,
			},
		})
		return
	}

	// Authorization: Check if user is broadcaster or message author
	isBroadcaster := stream.BroadcasterID == userUUID
	isAuthor := message.SenderID == userUUID

	if !isBroadcaster && !isAuthor {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Only broadcaster or message sender can delete this message",
			"details": gin.H{
				"stream_id":      streamID,
				"message_id":     messageID,
				"required_roles": []string{"broadcaster", "message_sender"},
			},
		})
		return
	}

	// Soft delete by updating moderation status to 'removed'
	// ScyllaDB doesn't support true deletes with TTL, so we mark as removed
	updates := map[string]interface{}{
		"moderation_status": models.ModerationStatusRemoved,
	}

	if err := h.chatRepo.Update(ctx, streamID, messageID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete message",
		})
		return
	}

	// Return 204 No Content on successful deletion
	c.Status(http.StatusNoContent)
}