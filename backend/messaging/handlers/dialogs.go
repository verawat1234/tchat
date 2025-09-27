package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat.dev/messaging/services"
	"tchat.dev/shared/middleware"
	"tchat.dev/shared/responses"
)

type DialogHandler struct {
	dialogService *services.DialogService
	validator     *validator.Validate
}

func NewDialogHandler(dialogService *services.DialogService) *DialogHandler {
	return &DialogHandler{
		dialogService: dialogService,
		validator:     validator.New(),
	}
}

// CreateDialogRequest represents the request to create a new dialog
type CreateDialogRequest struct {
	Type         string      `json:"type" validate:"required,oneof=user group channel business" example:"user"`
	Title        string      `json:"title,omitempty" validate:"omitempty,min=1,max=100" example:"Project Discussion"`
	Description  string      `json:"description,omitempty" validate:"omitempty,max=500" example:"Discussion about the new project features"`
	ParticipantIDs []uuid.UUID `json:"participant_ids" validate:"required,min=1,max=1000" example:"[\"123e4567-e89b-12d3-a456-426614174000\"]"`
	IsPrivate    bool        `json:"is_private,omitempty" example:"false"`
	Settings     *DialogSettings `json:"settings,omitempty"`
}

type DialogSettings struct {
	MuteNotifications bool     `json:"mute_notifications,omitempty" example:"false"`
	DisableInvites    bool     `json:"disable_invites,omitempty" example:"false"`
	ReadReceiptsOff   bool     `json:"read_receipts_off,omitempty" example:"false"`
	MessageRetention  int      `json:"message_retention,omitempty" example:"365"`
	AllowedFileTypes  []string `json:"allowed_file_types,omitempty" example:"[\"image\", \"document\"]"`
	MaxFileSize       int64    `json:"max_file_size,omitempty" example:"52428800"`
}

// UpdateDialogRequest represents the request to update dialog settings
type UpdateDialogRequest struct {
	Title       string          `json:"title,omitempty" validate:"omitempty,min=1,max=100" example:"Updated Project Discussion"`
	Description string          `json:"description,omitempty" validate:"omitempty,max=500" example:"Updated description"`
	Settings    *DialogSettings `json:"settings,omitempty"`
}

// AddParticipantsRequest represents the request to add participants to a dialog
type AddParticipantsRequest struct {
	ParticipantIDs []uuid.UUID `json:"participant_ids" validate:"required,min=1,max=100" example:"[\"123e4567-e89b-12d3-a456-426614174001\"]"`
	Role           string      `json:"role,omitempty" validate:"omitempty,oneof=member admin moderator" example:"member"`
}

// RemoveParticipantRequest represents the request to remove a participant
type RemoveParticipantRequest struct {
	ParticipantID uuid.UUID `json:"participant_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174001"`
	Reason        string    `json:"reason,omitempty" validate:"omitempty,max=200" example:"User requested to leave"`
}

// DialogResponse represents a dialog in API responses
type DialogResponse struct {
	ID                uuid.UUID             `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Type              string                `json:"type" example:"group"`
	Title             string                `json:"title" example:"Project Discussion"`
	Description       string                `json:"description,omitempty" example:"Discussion about the new project features"`
	CreatorID         uuid.UUID             `json:"creator_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	ParticipantCount  int                   `json:"participant_count" example:"5"`
	UnreadCount       int                   `json:"unread_count" example:"3"`
	LastMessage       *LastMessageInfo      `json:"last_message,omitempty"`
	LastActivity      string                `json:"last_activity" example:"2024-01-20T15:45:00Z"`
	IsPrivate         bool                  `json:"is_private" example:"false"`
	IsMuted           bool                  `json:"is_muted" example:"false"`
	IsArchived        bool                  `json:"is_archived" example:"false"`
	CreatedAt         string                `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt         string                `json:"updated_at" example:"2024-01-20T15:45:00Z"`
	Settings          *DialogSettings       `json:"settings,omitempty"`
	Participants      []ParticipantInfo     `json:"participants,omitempty"`
	UserRole          string                `json:"user_role,omitempty" example:"member"`
	CanInvite         bool                  `json:"can_invite" example:"true"`
	CanLeave          bool                  `json:"can_leave" example:"true"`
}

type LastMessageInfo struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174002"`
	Type      string    `json:"type" example:"text"`
	Content   string    `json:"content" example:"Hey everyone! How's the project going?"`
	SenderID  uuid.UUID `json:"sender_id" example:"123e4567-e89b-12d3-a456-426614174003"`
	SenderName string   `json:"sender_name" example:"John Doe"`
	Timestamp string    `json:"timestamp" example:"2024-01-20T15:45:00Z"`
}

type ParticipantInfo struct {
	UserID      uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	DisplayName string    `json:"display_name" example:"John Doe"`
	Avatar      string    `json:"avatar,omitempty" example:"https://cdn.tchat.com/avatars/user.jpg"`
	Role        string    `json:"role" example:"member"`
	JoinedAt    string    `json:"joined_at" example:"2024-01-15T10:30:00Z"`
	IsOnline    bool      `json:"is_online" example:"true"`
	LastSeen    string    `json:"last_seen,omitempty" example:"2024-01-20T14:30:00Z"`
}

// @Summary Get user dialogs
// @Description Get list of dialogs for the authenticated user
// @Tags dialogs
// @Produce json
// @Security BearerAuth
// @Param type query string false "Dialog type filter" Enums(user,group,channel,business)
// @Param limit query int false "Number of dialogs to return" default(20) minimum(1) maximum(100)
// @Param offset query int false "Number of dialogs to skip" default(0) minimum(0)
// @Param search query string false "Search term for dialog titles"
// @Param include_archived query bool false "Include archived dialogs" default(false)
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs [get]
func (h *DialogHandler) GetDialogs(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Parse query parameters
	dialogType := c.Query("type")
	search := c.Query("search")
	includeArchived, _ := strconv.ParseBool(c.Query("include_archived"))

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Build service request
	req := &services.GetDialogsRequest{
		UserID:          userUUID,
		Type:            dialogType,
		Search:          search,
		IncludeArchived: includeArchived,
		Limit:           limit,
		Offset:          offset,
	}

	// Get dialogs from service
	result, err := h.dialogService.GetUserDialogs(c.Request.Context(), req)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get dialogs", "Unable to retrieve dialogs.")
		return
	}

	// Convert to response format
	dialogResponses := make([]DialogResponse, len(result.Dialogs))
	for i, dialog := range result.Dialogs {
		dialogResponses[i] = h.convertToDialogResponse(dialog, userUUID)
	}

	// Build response data
	data := gin.H{
		"dialogs":     dialogResponses,
		"total":       result.Total,
		"limit":       limit,
		"offset":      offset,
		"has_more":    result.HasMore,
		"unread_total": result.UnreadTotal,
	}

	responses.DataResponse(c, data)
}

// @Summary Create a new dialog
// @Description Create a new dialog (conversation)
// @Tags dialogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateDialogRequest true "Dialog creation request"
// @Success 201 {object} responses.DataResponse{data=DialogResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs [post]
func (h *DialogHandler) CreateDialog(c *gin.Context) {
	var req CreateDialogRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.CreateDialogRequest{
		Type:           req.Type,
		Title:          req.Title,
		Description:    req.Description,
		CreatorID:      userUUID,
		ParticipantIDs: req.ParticipantIDs,
		IsPrivate:      req.IsPrivate,
		Settings:       h.convertDialogSettings(req.Settings),
	}

	// Create dialog
	dialog, err := h.dialogService.CreateDialog(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "validation failed"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
			return
		case strings.Contains(err.Error(), "participant not found"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Invalid participant", "One or more participants not found.")
			return
		case strings.Contains(err.Error(), "permission denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Permission denied", "You don't have permission to create this type of dialog.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Dialog creation failed", "Failed to create dialog.")
			return
		}
	}

	// Convert to response format
	dialogResponse := h.convertToDialogResponse(dialog, userUUID)

	// Log dialog creation
	middleware.LogInfo(c, "Dialog created", gin.H{
		"dialog_id":        dialog.ID,
		"type":             dialog.Type,
		"creator_id":       userUUID,
		"participant_count": len(req.ParticipantIDs),
	})

	c.JSON(http.StatusCreated, responses.DataResponse{
		Success: true,
		Data:    dialogResponse,
	})
}

// @Summary Get dialog by ID
// @Description Get detailed information about a specific dialog
// @Tags dialogs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dialog ID"
// @Param include_participants query bool false "Include participant details" default(false)
// @Success 200 {object} responses.DataResponse{data=DialogResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Dialog not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs/{id} [get]
func (h *DialogHandler) GetDialog(c *gin.Context) {
	// Parse dialog ID
	dialogIDStr := c.Param("id")
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid dialog ID", "Dialog ID must be a valid UUID.")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Parse query parameters
	includeParticipants, _ := strconv.ParseBool(c.Query("include_participants"))

	// Get dialog from service
	dialog, err := h.dialogService.GetDialogByID(c.Request.Context(), dialogID, userUUID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Dialog not found", "The specified dialog does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have access to this dialog.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get dialog", "Unable to retrieve dialog.")
			return
		}
	}

	// Convert to response format
	dialogResponse := h.convertToDialogResponse(dialog, userUUID)

	// Include participants if requested
	if includeParticipants {
		participants, err := h.dialogService.GetDialogParticipants(c.Request.Context(), dialogID, userUUID)
		if err == nil {
			dialogResponse.Participants = h.convertToParticipantInfo(participants)
		}
	}

	responses.DataResponse(c, dialogResponse)
}

// @Summary Update dialog
// @Description Update dialog settings and information
// @Tags dialogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dialog ID"
// @Param request body UpdateDialogRequest true "Dialog update request"
// @Success 200 {object} responses.DataResponse{data=DialogResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Permission denied"
// @Failure 404 {object} responses.ErrorResponse "Dialog not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs/{id} [put]
func (h *DialogHandler) UpdateDialog(c *gin.Context) {
	var req UpdateDialogRequest

	// Parse dialog ID
	dialogIDStr := c.Param("id")
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid dialog ID", "Dialog ID must be a valid UUID.")
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.UpdateDialogRequest{
		DialogID:    dialogID,
		UserID:      userUUID,
		Title:       req.Title,
		Description: req.Description,
		Settings:    h.convertDialogSettings(req.Settings),
	}

	// Update dialog
	dialog, err := h.dialogService.UpdateDialog(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Dialog not found", "The specified dialog does not exist.")
			return
		case strings.Contains(err.Error(), "permission denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Permission denied", "You don't have permission to update this dialog.")
			return
		case strings.Contains(err.Error(), "validation failed"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Dialog update failed", "Failed to update dialog.")
			return
		}
	}

	// Convert to response format
	dialogResponse := h.convertToDialogResponse(dialog, userUUID)

	// Log dialog update
	middleware.LogInfo(c, "Dialog updated", gin.H{
		"dialog_id": dialogID,
		"user_id":   userUUID,
	})

	responses.DataResponse(c, dialogResponse)
}

// @Summary Add participants to dialog
// @Description Add new participants to an existing dialog
// @Tags dialogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dialog ID"
// @Param request body AddParticipantsRequest true "Add participants request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Permission denied"
// @Failure 404 {object} responses.ErrorResponse "Dialog not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs/{id}/participants [post]
func (h *DialogHandler) AddParticipants(c *gin.Context) {
	var req AddParticipantsRequest

	// Parse dialog ID
	dialogIDStr := c.Param("id")
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid dialog ID", "Dialog ID must be a valid UUID.")
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "member"
	}

	// Build service request
	serviceReq := &services.AddParticipantsRequest{
		DialogID:       dialogID,
		UserID:         userUUID,
		ParticipantIDs: req.ParticipantIDs,
		Role:           req.Role,
	}

	// Add participants
	err = h.dialogService.AddParticipants(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Dialog not found", "The specified dialog does not exist.")
			return
		case strings.Contains(err.Error(), "permission denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Permission denied", "You don't have permission to add participants to this dialog.")
			return
		case strings.Contains(err.Error(), "participant already exists"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Participant already exists", "One or more participants are already in this dialog.")
			return
		case strings.Contains(err.Error(), "participant limit"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Participant limit exceeded", "Cannot add more participants to this dialog.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to add participants", "Unable to add participants to dialog.")
			return
		}
	}

	// Log participants addition
	middleware.LogInfo(c, "Participants added", gin.H{
		"dialog_id":        dialogID,
		"user_id":          userUUID,
		"participant_count": len(req.ParticipantIDs),
		"role":             req.Role,
	})

	responses.SuccessMessageResponse(c, "Participants added successfully")
}

// @Summary Remove participant from dialog
// @Description Remove a participant from an existing dialog
// @Tags dialogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dialog ID"
// @Param request body RemoveParticipantRequest true "Remove participant request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Permission denied"
// @Failure 404 {object} responses.ErrorResponse "Dialog or participant not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs/{id}/participants/remove [post]
func (h *DialogHandler) RemoveParticipant(c *gin.Context) {
	var req RemoveParticipantRequest

	// Parse dialog ID
	dialogIDStr := c.Param("id")
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid dialog ID", "Dialog ID must be a valid UUID.")
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.RemoveParticipantRequest{
		DialogID:      dialogID,
		UserID:        userUUID,
		ParticipantID: req.ParticipantID,
		Reason:        req.Reason,
	}

	// Remove participant
	err = h.dialogService.RemoveParticipant(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "dialog not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Dialog not found", "The specified dialog does not exist.")
			return
		case strings.Contains(err.Error(), "participant not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Participant not found", "The specified participant is not in this dialog.")
			return
		case strings.Contains(err.Error(), "permission denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Permission denied", "You don't have permission to remove participants from this dialog.")
			return
		case strings.Contains(err.Error(), "cannot remove creator"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Cannot remove creator", "The dialog creator cannot be removed.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove participant", "Unable to remove participant from dialog.")
			return
		}
	}

	// Log participant removal
	middleware.LogInfo(c, "Participant removed", gin.H{
		"dialog_id":      dialogID,
		"user_id":        userUUID,
		"participant_id": req.ParticipantID,
		"reason":         req.Reason,
	})

	responses.SuccessMessageResponse(c, "Participant removed successfully")
}

// @Summary Leave dialog
// @Description Leave a dialog (remove yourself as participant)
// @Tags dialogs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dialog ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Cannot leave dialog"
// @Failure 404 {object} responses.ErrorResponse "Dialog not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /dialogs/{id}/leave [post]
func (h *DialogHandler) LeaveDialog(c *gin.Context) {
	// Parse dialog ID
	dialogIDStr := c.Param("id")
	dialogID, err := uuid.Parse(dialogIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid dialog ID", "Dialog ID must be a valid UUID.")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Leave dialog
	err = h.dialogService.LeaveDialog(c.Request.Context(), dialogID, userUUID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Dialog not found", "The specified dialog does not exist.")
			return
		case strings.Contains(err.Error(), "not a participant"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Not a participant", "You are not a participant in this dialog.")
			return
		case strings.Contains(err.Error(), "cannot leave"):
			responses.ErrorResponse(c, http.StatusForbidden, "Cannot leave dialog", "You cannot leave this dialog.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to leave dialog", "Unable to leave dialog.")
			return
		}
	}

	// Log dialog leave
	middleware.LogInfo(c, "Left dialog", gin.H{
		"dialog_id": dialogID,
		"user_id":   userUUID,
	})

	responses.SuccessMessageResponse(c, "Left dialog successfully")
}

// Helper functions

func (h *DialogHandler) convertToDialogResponse(dialog interface{}, userID uuid.UUID) DialogResponse {
	// This is a simplified conversion - in a real implementation,
	// you would properly convert from your dialog model to the response
	return DialogResponse{
		// Populate fields from dialog model
		// This is just a placeholder structure
	}
}

func (h *DialogHandler) convertToParticipantInfo(participants interface{}) []ParticipantInfo {
	// This is a simplified conversion - in a real implementation,
	// you would properly convert from your participant models to the response
	return []ParticipantInfo{}
}

func (h *DialogHandler) convertDialogSettings(settings *DialogSettings) interface{} {
	if settings == nil {
		return nil
	}
	// Convert to service layer settings structure
	return settings
}

// RegisterDialogRoutes registers all dialog-related routes
func RegisterDialogRoutes(router *gin.RouterGroup, dialogService *services.DialogService) {
	handler := NewDialogHandler(dialogService)

	// Protected dialog routes
	dialogs := router.Group("/dialogs")
	dialogs.Use(middleware.AuthRequired())
	{
		dialogs.GET("", handler.GetDialogs)
		dialogs.POST("", handler.CreateDialog)
		dialogs.GET("/:id", handler.GetDialog)
		dialogs.PUT("/:id", handler.UpdateDialog)
		dialogs.POST("/:id/participants", handler.AddParticipants)
		dialogs.POST("/:id/participants/remove", handler.RemoveParticipant)
		dialogs.POST("/:id/leave", handler.LeaveDialog)
	}
}

// RegisterDialogRoutesWithMiddleware registers dialog routes with custom middleware
func RegisterDialogRoutesWithMiddleware(
	router *gin.RouterGroup,
	dialogService *services.DialogService,
	middlewares ...gin.HandlerFunc,
) {
	handler := NewDialogHandler(dialogService)

	// Protected dialog routes with middleware
	dialogs := router.Group("/dialogs")
	allMiddlewares := append(middlewares, middleware.AuthRequired())
	dialogs.Use(allMiddlewares...)
	{
		dialogs.GET("", handler.GetDialogs)
		dialogs.POST("", handler.CreateDialog)
		dialogs.GET("/:id", handler.GetDialog)
		dialogs.PUT("/:id", handler.UpdateDialog)
		dialogs.POST("/:id/participants", handler.AddParticipants)
		dialogs.POST("/:id/participants/remove", handler.RemoveParticipant)
		dialogs.POST("/:id/leave", handler.LeaveDialog)
	}
}