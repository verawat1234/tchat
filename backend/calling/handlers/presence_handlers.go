package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/calling/models"
	"tchat.dev/calling/services"
)

// PresenceHandlers handles HTTP requests for user presence management
type PresenceHandlers struct {
	presenceService *services.PresenceService
}

// NewPresenceHandlers creates a new PresenceHandlers instance
func NewPresenceHandlers(presenceService *services.PresenceService) *PresenceHandlers {
	return &PresenceHandlers{
		presenceService: presenceService,
	}
}

// UpdatePresenceRequest represents the request body for updating presence
type UpdatePresenceRequest struct {
	Status string `json:"status" binding:"required,oneof=online busy offline"`
}

// GetPresenceResponse represents the response for presence information
type GetPresenceResponse struct {
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	LastSeen  string `json:"last_seen"`
	UpdatedAt string `json:"updated_at"`
}

// CheckAvailabilityResponse represents the response for availability check
type CheckAvailabilityResponse struct {
	UserID    string `json:"user_id"`
	Available bool   `json:"available"`
	Status    string `json:"status"`
	LastSeen  string `json:"last_seen"`
}

// GetPresenceStatus handles GET /presence/status requests
func (h *PresenceHandlers) GetPresenceStatus(c *gin.Context) {
	// Get user ID from query parameter or JWT token (simplified for now)
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

	// Get presence status
	presence, err := h.presenceService.GetPresence(userID)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
				Code:  "user_not_found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
				Code:  "internal_error",
			})
		}
		return
	}

	response := GetPresenceResponse{
		UserID:    presence.UserID.String(),
		Status:    string(presence.Status),
		LastSeen:  presence.LastSeen.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: presence.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Presence status retrieved successfully",
		Data:    response,
	})
}

// UpdatePresenceStatus handles PUT /presence/status requests
func (h *PresenceHandlers) UpdatePresenceStatus(c *gin.Context) {
	// Get user ID from query parameter or JWT token (simplified for now)
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

	var req UpdatePresenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Convert status string to PresenceStatus
	var status models.PresenceStatus
	switch req.Status {
	case "online":
		status = models.PresenceStatusOnline
	case "busy":
		status = models.PresenceStatusBusy
	case "offline":
		status = models.PresenceStatusOffline
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid status",
			Details: "status must be 'online', 'busy', or 'offline'",
		})
		return
	}

	// Update presence
	serviceReq := services.UpdatePresenceRequest{
		UserID: userID,
		Status: status,
	}

	presence, err := h.presenceService.UpdatePresence(serviceReq)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
				Code:  "user_not_found",
			})
		case models.ErrInvalidPresenceStatus:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid presence status",
				Code:  "invalid_presence_status",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
				Code:  "internal_error",
			})
		}
		return
	}

	response := GetPresenceResponse{
		UserID:    presence.UserID.String(),
		Status:    string(presence.Status),
		LastSeen:  presence.LastSeen.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: presence.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Presence status updated successfully",
		Data:    response,
	})
}

// CheckUserAvailability handles GET /presence/check/{user_id} requests
func (h *PresenceHandlers) CheckUserAvailability(c *gin.Context) {
	// Get user ID from path parameter
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user_id format",
			Details: "user_id must be a valid UUID",
		})
		return
	}

	// Check availability
	availability, err := h.presenceService.CheckAvailability(userID)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
				Code:  "user_not_found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
				Code:  "internal_error",
			})
		}
		return
	}

	response := CheckAvailabilityResponse{
		UserID:    availability.UserID.String(),
		Available: availability.Available,
		Status:    string(availability.Status),
		LastSeen:  availability.LastSeen.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User availability checked successfully",
		Data:    response,
	})
}

// GetOnlineUsers handles GET /presence/online requests
func (h *PresenceHandlers) GetOnlineUsers(c *gin.Context) {
	// Get online users
	presences, err := h.presenceService.GetOnlineUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
			Code:  "internal_error",
		})
		return
	}

	// Convert to response format
	responses := make([]GetPresenceResponse, len(presences))
	for i, presence := range presences {
		responses[i] = GetPresenceResponse{
			UserID:    presence.UserID.String(),
			Status:    string(presence.Status),
			LastSeen:  presence.LastSeen.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: presence.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Online users retrieved successfully",
		Data:    responses,
	})
}

// GetPresenceStats handles GET /presence/stats requests
func (h *PresenceHandlers) GetPresenceStats(c *gin.Context) {
	// Get presence statistics
	stats, err := h.presenceService.GetPresenceStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
			Code:  "internal_error",
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Presence statistics retrieved successfully",
		Data:    stats,
	})
}

// RegisterRoutes registers all presence-related routes
func (h *PresenceHandlers) RegisterRoutes(router *gin.RouterGroup) {
	presence := router.Group("/presence")
	{
		presence.GET("/status", h.GetPresenceStatus)
		presence.PUT("/status", h.UpdatePresenceStatus)
		presence.GET("/check/:user_id", h.CheckUserAvailability)
		presence.GET("/online", h.GetOnlineUsers)
		presence.GET("/stats", h.GetPresenceStats)
	}
}
