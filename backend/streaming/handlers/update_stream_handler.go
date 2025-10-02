// backend/streaming/handlers/update_stream_handler.go
// UpdateStreamHandler - Handles PATCH /api/v1/streams/{streamId} endpoint
// Implements T041: Update stream metadata with authorization and validation

package handlers

import (
	"database/sql"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

// UpdateStreamHandler handles stream update operations with authorization
type UpdateStreamHandler struct {
	liveStreamRepo repository.LiveStreamRepositoryInterface
}

// UpdateStreamRequest defines the updatable fields for a live stream
// All fields are optional (pointer types) for partial updates
type UpdateStreamRequest struct {
	Title            *string    `json:"title,omitempty"`
	Description      *string    `json:"description,omitempty"`
	PrivacySetting   *string    `json:"privacy_setting,omitempty"`
	ScheduledStart   *time.Time `json:"scheduled_start_time,omitempty"`
	MaxCapacity      *int       `json:"max_capacity,omitempty"`
	ThumbnailURL     *string    `json:"thumbnail_url,omitempty"`
	Language         *string    `json:"language,omitempty"`
	Tags             *[]string  `json:"tags,omitempty"`
}

// UpdateStreamResponse wraps the updated stream data
type UpdateStreamResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *models.LiveStream `json:"data,omitempty"`
	Error   string             `json:"error,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewUpdateStreamHandler creates a new update stream handler instance
func NewUpdateStreamHandler(liveStreamRepo repository.LiveStreamRepositoryInterface) *UpdateStreamHandler {
	return &UpdateStreamHandler{
		liveStreamRepo: liveStreamRepo,
	}
}

// Handle processes PATCH /api/v1/streams/{streamId} requests
func (h *UpdateStreamHandler) Handle(c *gin.Context) {
	// Extract broadcaster ID from JWT token
	broadcasterID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, UpdateStreamResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	broadcasterUUID, ok := broadcasterID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, UpdateStreamResponse{
			Success: false,
			Error:   "Invalid user ID format",
		})
		return
	}

	// Parse and validate stream ID from URL parameter
	streamIDStr := c.Param("streamId")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, UpdateStreamResponse{
			Success: false,
			Error:   "Invalid stream ID format",
			Details: map[string]interface{}{
				"stream_id": streamIDStr,
			},
		})
		return
	}

	// Fetch existing stream from database
	stream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, UpdateStreamResponse{
			Success: false,
			Error:   "Stream not found",
			Details: map[string]interface{}{
				"stream_id": streamID.String(),
			},
		})
		return
	}

	// Authorization check: Only stream broadcaster can update
	if stream.BroadcasterID != broadcasterUUID {
		c.JSON(http.StatusForbidden, UpdateStreamResponse{
			Success: false,
			Error:   "Only stream broadcaster can perform this action",
			Details: map[string]interface{}{
				"stream_id":     streamID.String(),
				"required_role": "owner",
			},
		})
		return
	}

	// Parse request body
	var req UpdateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UpdateStreamResponse{
			Success: false,
			Error:   "Invalid request body",
			Details: map[string]interface{}{
				"parse_error": err.Error(),
			},
		})
		return
	}

	// Validate stream status restrictions
	if stream.Status == "live" {
		// Check if any restricted fields are being updated
		restrictedFields := h.getRestrictedFieldsForLiveStream(&req)
		if len(restrictedFields) > 0 {
			c.JSON(http.StatusBadRequest, UpdateStreamResponse{
				Success: false,
				Error:   "Cannot update certain fields while stream is live",
				Details: map[string]interface{}{
					"stream_id":         streamID.String(),
					"current_status":    "live",
					"restricted_fields": restrictedFields,
					"allowed_fields":    []string{"title", "description", "tags"},
				},
			})
			return
		}
	}

	// Build update map with validation
	updates := make(map[string]interface{})

	// Validate and add title
	if req.Title != nil {
		if len(*req.Title) < 1 || len(*req.Title) > 200 {
			c.JSON(http.StatusBadRequest, UpdateStreamResponse{
				Success: false,
				Error:   "Title must be between 1 and 200 characters",
			})
			return
		}
		updates["title"] = *req.Title
	}

	// Validate and add description
	if req.Description != nil {
		if len(*req.Description) > 5000 {
			c.JSON(http.StatusBadRequest, UpdateStreamResponse{
				Success: false,
				Error:   "Description cannot exceed 5000 characters",
			})
			return
		}
		updates["description"] = sql.NullString{
			String: *req.Description,
			Valid:  *req.Description != "",
		}
	}

	// Validate and add privacy setting
	if req.PrivacySetting != nil {
		validPrivacy := map[string]bool{
			"public":         true,
			"followers_only": true,
			"private":        true,
		}
		if !validPrivacy[*req.PrivacySetting] {
			c.JSON(http.StatusBadRequest, UpdateStreamResponse{
				Success: false,
				Error:   "Invalid privacy setting. Must be: public, followers_only, or private",
			})
			return
		}
		updates["privacy_setting"] = *req.PrivacySetting
	}

	// Validate and add scheduled start time
	if req.ScheduledStart != nil {
		if req.ScheduledStart.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, UpdateStreamResponse{
				Success: false,
				Error:   "Scheduled start time must be in the future",
			})
			return
		}
		updates["scheduled_start_time"] = req.ScheduledStart
	}

	// Validate and add max capacity
	if req.MaxCapacity != nil {
		if *req.MaxCapacity < 1 || *req.MaxCapacity > 100000 {
			c.JSON(http.StatusBadRequest, UpdateStreamResponse{
				Success: false,
				Error:   "Max capacity must be between 1 and 100000",
			})
			return
		}
		updates["max_capacity"] = *req.MaxCapacity
	}

	// Validate and add thumbnail URL
	if req.ThumbnailURL != nil {
		if *req.ThumbnailURL != "" {
			if !h.isValidURL(*req.ThumbnailURL) {
				c.JSON(http.StatusBadRequest, UpdateStreamResponse{
					Success: false,
					Error:   "Invalid thumbnail URL format",
				})
				return
			}
			updates["thumbnail_url"] = sql.NullString{
				String: *req.ThumbnailURL,
				Valid:  true,
			}
		} else {
			updates["thumbnail_url"] = sql.NullString{Valid: false}
		}
	}

	// Validate and add language
	if req.Language != nil {
		if len(*req.Language) < 2 || len(*req.Language) > 10 {
			c.JSON(http.StatusBadRequest, UpdateStreamResponse{
				Success: false,
				Error:   "Language code must be between 2 and 10 characters",
			})
			return
		}
		updates["language"] = *req.Language
	}

	// Add tags
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}

	// Check if there are any updates to apply
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, UpdateStreamResponse{
			Success: false,
			Error:   "No valid fields provided for update",
		})
		return
	}

	// Set updated_at timestamp
	updates["updated_at"] = time.Now()

	// Perform update in database
	if err := h.liveStreamRepo.Update(streamID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, UpdateStreamResponse{
			Success: false,
			Error:   "Failed to update stream",
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	// Fetch updated stream
	updatedStream, err := h.liveStreamRepo.GetByID(streamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateStreamResponse{
			Success: false,
			Error:   "Failed to fetch updated stream",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, UpdateStreamResponse{
		Success: true,
		Message: "Stream updated successfully",
		Data:    updatedStream,
	})
}

// getRestrictedFieldsForLiveStream returns fields that cannot be updated during live stream
func (h *UpdateStreamHandler) getRestrictedFieldsForLiveStream(req *UpdateStreamRequest) []string {
	restricted := []string{}

	if req.PrivacySetting != nil {
		restricted = append(restricted, "privacy_setting")
	}

	if req.ScheduledStart != nil {
		restricted = append(restricted, "scheduled_start_time")
	}

	if req.MaxCapacity != nil {
		restricted = append(restricted, "max_capacity")
	}

	if req.Language != nil {
		restricted = append(restricted, "language")
	}

	return restricted
}

// isValidURL validates URL format with basic regex
func (h *UpdateStreamHandler) isValidURL(urlStr string) bool {
	// Basic URL validation for http/https URLs
	pattern := `^https?://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(/[a-zA-Z0-9\-\_/\.]*)?(\.(jpg|jpeg|png|webp|gif))?$`
	matched, _ := regexp.MatchString(pattern, urlStr)
	return matched
}