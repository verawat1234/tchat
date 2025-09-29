package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/calling/services"
)

// WebSocketHandler handles WebSocket connections for real-time signaling
type WebSocketHandler struct {
	signalingService *services.SignalingService
}

// NewWebSocketHandler creates a new WebSocketHandler instance
func NewWebSocketHandler(signalingService *services.SignalingService) *WebSocketHandler {
	return &WebSocketHandler{
		signalingService: signalingService,
	}
}

// HandleWebSocket handles WebSocket connection upgrades
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from query parameter (in production, get from JWT token)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing user_id parameter",
			Details: "user_id query parameter is required",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user_id format",
			Details: "user_id must be a valid UUID",
		})
		return
	}

	// Handle WebSocket connection through SignalingService
	if err := h.signalingService.HandleWebSocketConnection(c.Writer, c.Request, userID); err != nil {
		log.Printf("Failed to handle WebSocket connection: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to handle WebSocket connection",
			Details: err.Error(),
		})
		return
	}
}

// RegisterRoutes registers WebSocket routes
func (h *WebSocketHandler) RegisterRoutes(router *gin.RouterGroup) {
	// WebSocket endpoint for signaling
	router.GET("/ws", h.HandleWebSocket)

	// Optional: Health check endpoint for WebSocket service
	router.GET("/ws/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, SuccessResponse{
			Message: "WebSocket service is healthy",
			Data: map[string]interface{}{
				"status": "healthy",
			},
		})
	})
}
