package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/calling/models"
	"tchat.dev/calling/services"
)

// CallHandlers handles HTTP requests for call management
type CallHandlers struct {
	callService *services.CallService
}

// NewCallHandlers creates a new CallHandlers instance
func NewCallHandlers(callService *services.CallService) *CallHandlers {
	return &CallHandlers{
		callService: callService,
	}
}

// InitiateCallRequest represents the request body for initiating a call
type InitiateCallRequest struct {
	CallerID string `json:"caller_id" binding:"required,uuid"`
	CalleeID string `json:"callee_id" binding:"required,uuid"`
	CallType string `json:"call_type" binding:"required,oneof=voice video"`
}

// AnswerCallRequest represents the request body for answering a call
type AnswerCallRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	Accept bool   `json:"accept"`
}

// EndCallRequest represents the request body for ending a call
type EndCallRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a successful response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// InitiateCall handles POST /calls/initiate requests
func (h *CallHandlers) InitiateCall(c *gin.Context) {
	var req InitiateCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Parse UUIDs
	callerID, err := uuid.Parse(req.CallerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid caller_id format",
			Details: "caller_id must be a valid UUID",
		})
		return
	}

	calleeID, err := uuid.Parse(req.CalleeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid callee_id format",
			Details: "callee_id must be a valid UUID",
		})
		return
	}

	// Convert call type
	var callType models.CallType
	switch req.CallType {
	case "voice":
		callType = models.CallTypeVoice
	case "video":
		callType = models.CallTypeVideo
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid call_type",
			Details: "call_type must be 'voice' or 'video'",
		})
		return
	}

	// Create service request
	serviceReq := services.InitiateCallRequest{
		CallerID: callerID,
		CalleeID: calleeID,
		CallType: callType,
	}

	// Initiate the call
	callSession, err := h.callService.InitiateCall(serviceReq)
	if err != nil {
		// Map service errors to HTTP status codes
		switch err {
		case models.ErrUserAlreadyInCall:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "User already in call",
				Code:  "user_already_in_call",
			})
		case models.ErrUserNotAvailable:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "User not available",
				Code:  "user_not_available",
			})
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
				Code:  "user_not_found",
			})
		case models.ErrInvalidCallState:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid call state",
				Code:  "invalid_call_state",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
				Code:  "internal_error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Call initiated successfully",
		Data:    callSession,
	})
}

// AnswerCall handles POST /calls/{id}/answer requests
func (h *CallHandlers) AnswerCall(c *gin.Context) {
	// Get call ID from path parameter
	callIDStr := c.Param("id")
	callID, err := uuid.Parse(callIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid call ID format",
			Details: "Call ID must be a valid UUID",
		})
		return
	}

	var req AnswerCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user_id format",
			Details: "user_id must be a valid UUID",
		})
		return
	}

	// Create service request
	serviceReq := services.AnswerCallRequest{
		CallID: callID,
		UserID: userID,
		Accept: req.Accept,
	}

	// Answer the call
	callSession, err := h.callService.AnswerCall(serviceReq)
	if err != nil {
		// Map service errors to HTTP status codes
		switch err {
		case models.ErrCallNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Call not found",
				Code:  "call_not_found",
			})
		case models.ErrCallAlreadyAnswered:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Call already answered",
				Code:  "call_already_answered",
			})
		case models.ErrParticipantNotFound:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "User not participant in call",
				Code:  "participant_not_found",
			})
		case models.ErrPermissionDenied:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Permission denied",
				Code:  "permission_denied",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
				Code:  "internal_error",
			})
		}
		return
	}

	action := "declined"
	if req.Accept {
		action = "accepted"
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Call " + action + " successfully",
		Data:    callSession,
	})
}

// EndCall handles POST /calls/{id}/end requests
func (h *CallHandlers) EndCall(c *gin.Context) {
	// Get call ID from path parameter
	callIDStr := c.Param("id")
	callID, err := uuid.Parse(callIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid call ID format",
			Details: "Call ID must be a valid UUID",
		})
		return
	}

	var req EndCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user_id format",
			Details: "user_id must be a valid UUID",
		})
		return
	}

	// Create service request
	serviceReq := services.EndCallRequest{
		CallID: callID,
		UserID: userID,
	}

	// End the call
	callSession, err := h.callService.EndCall(serviceReq)
	if err != nil {
		// Map service errors to HTTP status codes
		switch err {
		case models.ErrCallNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Call not found",
				Code:  "call_not_found",
			})
		case models.ErrCallAlreadyEnded:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Call already ended",
				Code:  "call_already_ended",
			})
		case models.ErrParticipantNotFound:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "User not participant in call",
				Code:  "participant_not_found",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
				Code:  "internal_error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Call ended successfully",
		Data:    callSession,
	})
}

// GetCallStatus handles GET /calls/{id}/status requests
func (h *CallHandlers) GetCallStatus(c *gin.Context) {
	// Get call ID from path parameter
	callIDStr := c.Param("id")
	callID, err := uuid.Parse(callIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid call ID format",
			Details: "Call ID must be a valid UUID",
		})
		return
	}

	// Get user ID from query parameter for authorization
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

	// Get call status
	callSession, err := h.callService.GetCallStatus(callID, userID)
	if err != nil {
		// Map service errors to HTTP status codes
		switch err {
		case models.ErrCallNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Call not found",
				Code:  "call_not_found",
			})
		case models.ErrPermissionDenied:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "Permission denied",
				Code:  "permission_denied",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
				Code:  "internal_error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Call status retrieved successfully",
		Data:    callSession,
	})
}

// RegisterRoutes registers all call-related routes
func (h *CallHandlers) RegisterRoutes(router *gin.RouterGroup) {
	calls := router.Group("/calls")
	{
		calls.POST("/initiate", h.InitiateCall)
		calls.POST("/:id/answer", h.AnswerCall)
		calls.POST("/:id/end", h.EndCall)
		calls.GET("/:id/status", h.GetCallStatus)
	}
}
