package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/calling/models"
	"tchat.dev/calling/repositories"
)

// HistoryHandlers handles HTTP requests for call history management
type HistoryHandlers struct {
	historyRepo repositories.CallHistoryRepository
}

// NewHistoryHandlers creates a new HistoryHandlers instance
func NewHistoryHandlers(historyRepo repositories.CallHistoryRepository) *HistoryHandlers {
	return &HistoryHandlers{
		historyRepo: historyRepo,
	}
}

// CallHistoryResponse represents a call history item in the response
type CallHistoryResponse struct {
	ID                   string  `json:"id"`
	CallSessionID        string  `json:"call_session_id"`
	UserID               string  `json:"user_id"`
	OtherParticipantID   string  `json:"other_participant_id"`
	OtherParticipantName *string `json:"other_participant_name,omitempty"`
	CallType             string  `json:"call_type"`
	CallStatus           string  `json:"call_status"`
	InitiatedByMe        bool    `json:"initiated_by_me"`
	Duration             int     `json:"duration"` // Duration in seconds
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	Direction            string  `json:"direction"` // "incoming" or "outgoing"
	Outcome              string  `json:"outcome"`   // "completed", "missed", "failed", etc.
}

// CallHistoryStatsResponse represents call statistics in the response
type CallHistoryStatsResponse struct {
	TotalCalls      int     `json:"total_calls"`
	OutgoingCalls   int     `json:"outgoing_calls"`
	IncomingCalls   int     `json:"incoming_calls"`
	SuccessfulCalls int     `json:"successful_calls"`
	MissedCalls     int     `json:"missed_calls"`
	VoiceCalls      int     `json:"voice_calls"`
	VideoCalls      int     `json:"video_calls"`
	TotalDuration   int     `json:"total_duration"`   // In seconds
	AverageDuration float64 `json:"average_duration"` // In seconds
}

// GetCallHistory handles GET /history requests
func (h *HistoryHandlers) GetCallHistory(c *gin.Context) {
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

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid limit parameter",
			Details: "limit must be a number between 0 and 100",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid offset parameter",
			Details: "offset must be a non-negative number",
		})
		return
	}

	// Parse optional call type filter
	var callTypeFilter *models.CallType
	if callTypeStr := c.Query("call_type"); callTypeStr != "" {
		switch callTypeStr {
		case "voice":
			ct := models.CallTypeVoice
			callTypeFilter = &ct
		case "video":
			ct := models.CallTypeVideo
			callTypeFilter = &ct
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid call_type parameter",
				Details: "call_type must be 'voice' or 'video'",
			})
			return
		}
	}

	// Get call history
	var histories []models.CallHistory
	var getErr error

	if callTypeFilter != nil {
		histories, getErr = h.historyRepo.GetCallHistoryByType(userID, *callTypeFilter, limit, offset)
	} else {
		histories, getErr = h.historyRepo.GetCallHistory(userID, limit, offset)
	}

	if getErr != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
			Code:  "internal_error",
		})
		return
	}

	// Convert to response format
	responses := make([]CallHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = CallHistoryResponse{
			ID:                   history.ID.String(),
			CallSessionID:        history.CallSessionID.String(),
			UserID:               history.UserID.String(),
			OtherParticipantID:   history.OtherParticipantID.String(),
			OtherParticipantName: history.OtherParticipantName,
			CallType:             string(history.CallType),
			CallStatus:           string(history.CallStatus),
			InitiatedByMe:        history.InitiatedByMe,
			Duration:             history.Duration,
			Direction:            history.GetCallDirection(),
			Outcome:              history.GetCallOutcome(),
			CreatedAt:            history.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:            history.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Call history retrieved successfully",
		Data: map[string]interface{}{
			"items":  responses,
			"limit":  limit,
			"offset": offset,
			"count":  len(responses),
		},
	})
}

// GetCallStats handles GET /history/stats requests
func (h *HistoryHandlers) GetCallStats(c *gin.Context) {
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

	// Get call statistics
	stats, err := h.historyRepo.GetCallStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
			Code:  "internal_error",
		})
		return
	}

	// Convert to response format with type assertions
	response := CallHistoryStatsResponse{
		TotalCalls:      int(stats["total_calls"].(int64)),
		OutgoingCalls:   int(stats["outgoing_calls"].(int64)),
		IncomingCalls:   int(stats["incoming_calls"].(int64)),
		SuccessfulCalls: int(stats["successful_calls"].(int64)),
		MissedCalls:     int(stats["missed_calls"].(int64)),
		VoiceCalls:      int(stats["voice_calls"].(int64)),
		VideoCalls:      int(stats["video_calls"].(int64)),
		TotalDuration:   int(stats["total_duration"].(int64)),
		AverageDuration: stats["average_duration"].(float64),
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Call statistics retrieved successfully",
		Data:    response,
	})
}

// GetRecentHistory handles GET /history/recent requests
func (h *HistoryHandlers) GetRecentHistory(c *gin.Context) {
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

	// Parse optional parameters
	daysStr := c.DefaultQuery("days", "7")
	limitStr := c.DefaultQuery("limit", "10")

	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 0 || days > 365 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid days parameter",
			Details: "days must be a number between 0 and 365",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid limit parameter",
			Details: "limit must be a number between 0 and 100",
		})
		return
	}

	// Get recent call history
	histories, err := h.historyRepo.GetRecentCallHistory(userID, days, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
			Code:  "internal_error",
		})
		return
	}

	// Convert to response format
	responses := make([]CallHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = CallHistoryResponse{
			ID:                   history.ID.String(),
			CallSessionID:        history.CallSessionID.String(),
			UserID:               history.UserID.String(),
			OtherParticipantID:   history.OtherParticipantID.String(),
			OtherParticipantName: history.OtherParticipantName,
			CallType:             string(history.CallType),
			CallStatus:           string(history.CallStatus),
			InitiatedByMe:        history.InitiatedByMe,
			Duration:             history.Duration,
			Direction:            history.GetCallDirection(),
			Outcome:              history.GetCallOutcome(),
			CreatedAt:            history.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:            history.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Recent call history retrieved successfully",
		Data: map[string]interface{}{
			"items": responses,
			"days":  days,
			"limit": limit,
			"count": len(responses),
		},
	})
}

// RegisterRoutes registers all history-related routes
func (h *HistoryHandlers) RegisterRoutes(router *gin.RouterGroup) {
	history := router.Group("/history")
	{
		history.GET("", h.GetCallHistory)
		history.GET("/stats", h.GetCallStats)
		history.GET("/recent", h.GetRecentHistory)
	}
}
