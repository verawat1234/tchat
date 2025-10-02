package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/services"
)

// SignalingWebSocketHandler handles WebSocket connections for signaling
func SignalingWebSocket(signalingService *services.SignalingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user_id from context (set by WebSocketAuth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		userUUID, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := signalingService.UpgradeConnection(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
			return
		}

		// Handle WebSocket connection
		signalingService.HandleConnection(conn, userUUID)
	}
}