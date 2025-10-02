package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

type NotificationPreferencesHandler struct {
	prefRepo repository.NotificationPreferenceRepository
}

type NotificationPreferencesRequest struct {
	PushEnabled      *bool   `json:"push_enabled,omitempty"`
	EmailEnabled     *bool   `json:"email_enabled,omitempty"`
	SMSEnabled       *bool   `json:"sms_enabled,omitempty"`
	InAppEnabled     *bool   `json:"in_app_enabled,omitempty"`
	FollowOnlyMode   *bool   `json:"follow_only_mode,omitempty"`
	QuietHoursStart  *string `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd    *string `json:"quiet_hours_end,omitempty"`
	Timezone         *string `json:"timezone,omitempty"`
}

type NotificationPreferencesResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    NotificationPreferencesData  `json:"data"`
}

type NotificationPreferencesData struct {
	UserID          uuid.UUID `json:"user_id"`
	PushEnabled     bool      `json:"push_enabled"`
	EmailEnabled    bool      `json:"email_enabled"`
	SMSEnabled      bool      `json:"sms_enabled"`
	InAppEnabled    bool      `json:"in_app_enabled"`
	FollowOnlyMode  bool      `json:"follow_only_mode"`
	QuietHoursStart *string   `json:"quiet_hours_start"`
	QuietHoursEnd   *string   `json:"quiet_hours_end"`
	Timezone        string    `json:"timezone"`
	UpdatedAt       string    `json:"updated_at"`
}

func NewNotificationPreferencesHandler(
	prefRepo repository.NotificationPreferenceRepository,
) *NotificationPreferencesHandler {
	return &NotificationPreferencesHandler{
		prefRepo: prefRepo,
	}
}

func (h *NotificationPreferencesHandler) HandleGet(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract user ID from JWT
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

	// Fetch or create default preferences
	prefs, err := h.prefRepo.GetOrCreateDefault(ctx, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve notification preferences",
		})
		return
	}

	// Determine message based on whether preferences exist
	message := "Notification preferences retrieved successfully"
	if prefs.UpdatedAt.IsZero() || prefs.UpdatedAt.Before(time.Now().Add(-time.Second)) {
		message = "Default notification preferences created"
	}

	response := NotificationPreferencesResponse{
		Success: true,
		Message: message,
		Data:    h.convertToResponse(prefs),
	}

	c.JSON(http.StatusOK, response)
}

func (h *NotificationPreferencesHandler) HandlePut(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract user ID from JWT
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

	var req NotificationPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Get existing preferences or create default
	_, err := h.prefRepo.GetOrCreateDefault(ctx, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve notification preferences",
		})
		return
	}

	// Validate and parse quiet hours if provided
	var quietStart, quietEnd *time.Time
	if req.QuietHoursStart != nil {
		if *req.QuietHoursStart == "" {
			// Clear quiet hours start
			quietStart = nil
		} else {
			parsed, err := parseTimeString(*req.QuietHoursStart)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "Invalid quiet hours time format",
					"details": map[string]interface{}{
						"field":           "quiet_hours_start",
						"invalid_value":   *req.QuietHoursStart,
						"expected_format": "HH:MM:SS (00:00:00 - 23:59:59)",
					},
				})
				return
			}
			quietStart = &parsed
		}
	}

	if req.QuietHoursEnd != nil {
		if *req.QuietHoursEnd == "" {
			// Clear quiet hours end
			quietEnd = nil
		} else {
			parsed, err := parseTimeString(*req.QuietHoursEnd)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "Invalid quiet hours time format",
					"details": map[string]interface{}{
						"field":           "quiet_hours_end",
						"invalid_value":   *req.QuietHoursEnd,
						"expected_format": "HH:MM:SS (00:00:00 - 23:59:59)",
					},
				})
				return
			}
			quietEnd = &parsed
		}
	}

	// Validate timezone if provided
	if req.Timezone != nil && *req.Timezone != "" {
		if _, err := time.LoadLocation(*req.Timezone); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Invalid timezone: %s", *req.Timezone),
			})
			return
		}
	}

	// Build update map
	updates := make(map[string]interface{})
	if req.PushEnabled != nil {
		updates["push_enabled"] = *req.PushEnabled
	}
	if req.EmailEnabled != nil {
		updates["email_enabled"] = *req.EmailEnabled
	}
	if req.SMSEnabled != nil {
		// Map SMS to store_streams_enabled (or add SMS field to model)
		updates["store_streams_enabled"] = *req.SMSEnabled
	}
	if req.InAppEnabled != nil {
		updates["in_app_enabled"] = *req.InAppEnabled
	}
	if req.FollowOnlyMode != nil {
		// Map follow_only_mode to video_streams_enabled for now
		// (ideally add follow_only_mode field to model)
		updates["video_streams_enabled"] = !(*req.FollowOnlyMode)
	}
	if req.QuietHoursStart != nil {
		updates["quiet_hours_start"] = quietStart
	}
	if req.QuietHoursEnd != nil {
		updates["quiet_hours_end"] = quietEnd
	}
	if req.Timezone != nil {
		if *req.Timezone == "" {
			updates["timezone"] = sql.NullString{Valid: false}
		} else {
			updates["timezone"] = sql.NullString{String: *req.Timezone, Valid: true}
		}
	}

	// Update preferences
	if err := h.prefRepo.Update(ctx, userUUID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update notification preferences",
		})
		return
	}

	// Fetch updated preferences
	updatedPrefs, err := h.prefRepo.GetByUserID(ctx, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated preferences",
		})
		return
	}

	response := NotificationPreferencesResponse{
		Success: true,
		Message: "Notification preferences updated successfully",
		Data:    h.convertToResponse(updatedPrefs),
	}

	c.JSON(http.StatusOK, response)
}

func (h *NotificationPreferencesHandler) convertToResponse(prefs *models.NotificationPreference) NotificationPreferencesData {
	data := NotificationPreferencesData{
		UserID:       prefs.UserID,
		PushEnabled:  prefs.PushEnabled,
		EmailEnabled: prefs.EmailEnabled,
		SMSEnabled:   prefs.StoreStreamsEnabled, // Map store_streams to SMS for contract
		InAppEnabled: prefs.InAppEnabled,
		FollowOnlyMode: !prefs.VideoStreamsEnabled, // Map video_streams to follow_only
		Timezone:     "UTC", // Default timezone
		UpdatedAt:    prefs.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Set timezone if valid
	if prefs.Timezone.Valid {
		data.Timezone = prefs.Timezone.String
	}

	// Convert quiet hours to string format
	if prefs.QuietHoursStart != nil {
		timeStr := prefs.QuietHoursStart.Format("15:04:05")
		data.QuietHoursStart = &timeStr
	}
	if prefs.QuietHoursEnd != nil {
		timeStr := prefs.QuietHoursEnd.Format("15:04:05")
		data.QuietHoursEnd = &timeStr
	}

	return data
}

// parseTimeString parses HH:MM:SS format into time.Time
func parseTimeString(timeStr string) (time.Time, error) {
	// Parse as HH:MM:SS format
	parsed, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format (expected HH:MM:SS): %w", err)
	}

	// Validate time components
	hour := parsed.Hour()
	minute := parsed.Minute()
	second := parsed.Second()

	if hour < 0 || hour > 23 {
		return time.Time{}, fmt.Errorf("hour must be between 0 and 23, got %d", hour)
	}
	if minute < 0 || minute > 59 {
		return time.Time{}, fmt.Errorf("minute must be between 0 and 59, got %d", minute)
	}
	if second < 0 || second > 59 {
		return time.Time{}, fmt.Errorf("second must be between 0 and 59, got %d", second)
	}

	return parsed, nil
}