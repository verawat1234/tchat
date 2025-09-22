package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/messaging/models"
	sharedModels "backend/shared/models"
)

type DialogRepository interface {
	Create(ctx context.Context, dialog *models.Dialog) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Dialog, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, filters DialogFilters, pagination Pagination) ([]*models.Dialog, int64, error)
	Update(ctx context.Context, dialog *models.Dialog) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error)
	AddParticipant(ctx context.Context, participant *models.DialogParticipant) error
	RemoveParticipant(ctx context.Context, dialogID, userID uuid.UUID) error
	UpdateParticipant(ctx context.Context, participant *models.DialogParticipant) error
	GetAdmins(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error)
	SearchDialogs(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Dialog, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event *sharedModels.Event) error
}

type NotificationService interface {
	SendNotification(ctx context.Context, userID uuid.UUID, notificationType string, data map[string]interface{}) error
}

type DialogFilters struct {
	Type        *models.DialogType `json:"type,omitempty"`
	IsArchived  *bool              `json:"is_archived,omitempty"`
	IsMuted     *bool              `json:"is_muted,omitempty"`
	HasUnread   *bool              `json:"has_unread,omitempty"`
	UpdatedFrom *time.Time         `json:"updated_from,omitempty"`
	UpdatedTo   *time.Time         `json:"updated_to,omitempty"`
}

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	OrderBy  string `json:"order_by"`
	Order    string `json:"order"` // asc, desc
}

type DialogService struct {
	dialogRepo          DialogRepository
	eventPublisher      EventPublisher
	notificationService NotificationService
	db                  *gorm.DB
}

func NewDialogService(
	dialogRepo DialogRepository,
	eventPublisher EventPublisher,
	notificationService NotificationService,
	db *gorm.DB,
) *DialogService {
	return &DialogService{
		dialogRepo:          dialogRepo,
		eventPublisher:      eventPublisher,
		notificationService: notificationService,
		db:                  db,
	}
}

func (ds *DialogService) CreateDialog(ctx context.Context, req *CreateDialogRequest) (*models.Dialog, error) {
	// Validate request
	if err := ds.validateCreateDialogRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create dialog
	dialog := &models.Dialog{
		ID:          uuid.New(),
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Settings:    req.Settings,
		Metadata:    req.Metadata,
		CreatedBy:   req.CreatorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Set default settings based on dialog type
	dialog.SetDefaultSettings()

	// Validate dialog
	if err := dialog.Validate(); err != nil {
		return nil, fmt.Errorf("dialog validation failed: %w", err)
	}

	// Save dialog
	if err := ds.dialogRepo.Create(ctx, dialog); err != nil {
		return nil, fmt.Errorf("failed to create dialog: %w", err)
	}

	// Add creator as admin participant
	creatorParticipant := &models.DialogParticipant{
		ID:        uuid.New(),
		DialogID:  dialog.ID,
		UserID:    req.CreatorID,
		Role:      models.ParticipantRoleAdmin,
		Status:    models.ParticipantStatusActive,
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ds.dialogRepo.AddParticipant(ctx, creatorParticipant); err != nil {
		return nil, fmt.Errorf("failed to add creator as participant: %w", err)
	}

	// Add other participants if specified
	for _, userID := range req.ParticipantIDs {
		if userID != req.CreatorID { // Skip creator as already added
			participant := &models.DialogParticipant{
				ID:        uuid.New(),
				DialogID:  dialog.ID,
				UserID:    userID,
				Role:      models.ParticipantRoleMember,
				Status:    models.ParticipantStatusActive,
				JoinedAt:  time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := ds.dialogRepo.AddParticipant(ctx, participant); err != nil {
				fmt.Printf("Failed to add participant %s to dialog %s: %v\n", userID, dialog.ID, err)
			}
		}
	}

	// Update participant count
	dialog.ParticipantCount = len(req.ParticipantIDs)
	if err := ds.dialogRepo.Update(ctx, dialog); err != nil {
		fmt.Printf("Failed to update participant count: %v\n", err)
	}

	// Publish dialog created event
	if err := ds.publishDialogEvent(ctx, sharedModels.EventTypeDialogCreated, dialog.ID, req.CreatorID, map[string]interface{}{
		"dialog_id":         dialog.ID,
		"dialog_type":       dialog.Type,
		"participant_count": dialog.ParticipantCount,
		"title":            dialog.Title,
	}); err != nil {
		fmt.Printf("Failed to publish dialog created event: %v\n", err)
	}

	// Send notifications to participants (except creator)
	for _, userID := range req.ParticipantIDs {
		if userID != req.CreatorID {
			go ds.notificationService.SendNotification(context.Background(), userID, "dialog_invitation", map[string]interface{}{
				"dialog_id":   dialog.ID,
				"dialog_title": dialog.Title,
				"creator_id":  req.CreatorID,
			})
		}
	}

	return dialog, nil
}

func (ds *DialogService) GetDialogByID(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (*models.Dialog, error) {
	if dialogID == uuid.Nil {
		return nil, fmt.Errorf("dialog ID is required")
	}

	dialog, err := ds.dialogRepo.GetByID(ctx, dialogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("dialog not found")
		}
		return nil, fmt.Errorf("failed to get dialog: %w", err)
	}

	// Check if user has access to this dialog
	if !ds.userHasAccess(ctx, dialogID, userID) {
		return nil, fmt.Errorf("access denied")
	}

	return dialog, nil
}

func (ds *DialogService) GetUserDialogs(ctx context.Context, userID uuid.UUID, filters DialogFilters, pagination Pagination) ([]*models.Dialog, int64, error) {
	if userID == uuid.Nil {
		return nil, 0, fmt.Errorf("user ID is required")
	}

	// Validate pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}
	if pagination.OrderBy == "" {
		pagination.OrderBy = "updated_at"
	}
	if pagination.Order != "asc" && pagination.Order != "desc" {
		pagination.Order = "desc"
	}

	dialogs, total, err := ds.dialogRepo.GetByUserID(ctx, userID, filters, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user dialogs: %w", err)
	}

	return dialogs, total, nil
}

func (ds *DialogService) UpdateDialog(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID, req *UpdateDialogRequest) (*models.Dialog, error) {
	// Get dialog
	dialog, err := ds.GetDialogByID(ctx, dialogID, userID)
	if err != nil {
		return nil, err
	}

	// Check if user is admin
	if !ds.userIsAdmin(ctx, dialogID, userID) {
		return nil, fmt.Errorf("only admins can update dialog")
	}

	// Track changes for event
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.Title != nil && *req.Title != dialog.Title {
		changes["title"] = map[string]string{"from": dialog.Title, "to": *req.Title}
		dialog.Title = *req.Title
	}

	if req.Description != nil && *req.Description != dialog.Description {
		changes["description"] = map[string]string{"from": dialog.Description, "to": *req.Description}
		dialog.Description = *req.Description
	}

	if req.Settings != nil {
		changes["settings"] = map[string]interface{}{
			"from": dialog.Settings,
			"to":   *req.Settings,
		}
		dialog.Settings = *req.Settings
	}

	if req.Metadata != nil {
		// Merge metadata
		if dialog.Metadata == nil {
			dialog.Metadata = make(map[string]interface{})
		}
		for key, value := range req.Metadata {
			dialog.Metadata[key] = value
		}
		changes["metadata"] = req.Metadata
	}

	// Update timestamp
	dialog.UpdatedAt = time.Now()

	// Validate updated dialog
	if err := dialog.Validate(); err != nil {
		return nil, fmt.Errorf("updated dialog validation failed: %w", err)
	}

	// Save to database
	if err := ds.dialogRepo.Update(ctx, dialog); err != nil {
		return nil, fmt.Errorf("failed to update dialog: %w", err)
	}

	// Publish dialog updated event if there were changes
	if len(changes) > 0 {
		if err := ds.publishDialogEvent(ctx, "dialog.updated", dialogID, userID, map[string]interface{}{
			"changes": changes,
		}); err != nil {
			fmt.Printf("Failed to publish dialog updated event: %v\n", err)
		}
	}

	return dialog, nil
}

func (ds *DialogService) AddParticipant(ctx context.Context, dialogID uuid.UUID, adminUserID uuid.UUID, req *AddParticipantRequest) error {
	// Check if user is admin
	if !ds.userIsAdmin(ctx, dialogID, adminUserID) {
		return fmt.Errorf("only admins can add participants")
	}

	// Get dialog to check limits
	dialog, err := ds.dialogRepo.GetByID(ctx, dialogID)
	if err != nil {
		return fmt.Errorf("dialog not found")
	}

	// Check participant limit
	maxParticipants := dialog.GetMaxParticipants()
	if dialog.ParticipantCount >= maxParticipants {
		return fmt.Errorf("dialog has reached maximum participant limit (%d)", maxParticipants)
	}

	// Check if user is already a participant
	participants, err := ds.dialogRepo.GetParticipants(ctx, dialogID)
	if err != nil {
		return fmt.Errorf("failed to check existing participants: %w", err)
	}

	for _, participant := range participants {
		if participant.UserID == req.UserID {
			if participant.Status == models.ParticipantStatusActive {
				return fmt.Errorf("user is already a participant")
			} else if participant.Status == models.ParticipantStatusLeft {
				// Reactivate participant
				participant.Status = models.ParticipantStatusActive
				participant.JoinedAt = time.Now()
				participant.UpdatedAt = time.Now()
				return ds.dialogRepo.UpdateParticipant(ctx, participant)
			}
		}
	}

	// Add new participant
	participant := &models.DialogParticipant{
		ID:        uuid.New(),
		DialogID:  dialogID,
		UserID:    req.UserID,
		Role:      req.Role,
		Status:    models.ParticipantStatusActive,
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ds.dialogRepo.AddParticipant(ctx, participant); err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}

	// Update participant count
	dialog.ParticipantCount++
	dialog.UpdatedAt = time.Now()
	if err := ds.dialogRepo.Update(ctx, dialog); err != nil {
		fmt.Printf("Failed to update participant count: %v\n", err)
	}

	// Publish participant added event
	if err := ds.publishDialogEvent(ctx, sharedModels.EventTypeDialogParticipantAdded, dialogID, adminUserID, map[string]interface{}{
		"participant_id": req.UserID,
		"role":          req.Role,
		"added_by":      adminUserID,
	}); err != nil {
		fmt.Printf("Failed to publish participant added event: %v\n", err)
	}

	// Send notification to new participant
	go ds.notificationService.SendNotification(context.Background(), req.UserID, "dialog_invitation", map[string]interface{}{
		"dialog_id":    dialogID,
		"dialog_title": dialog.Title,
		"invited_by":   adminUserID,
	})

	return nil
}

func (ds *DialogService) RemoveParticipant(ctx context.Context, dialogID uuid.UUID, adminUserID uuid.UUID, userID uuid.UUID) error {
	// Check if user is admin (or removing themselves)
	if adminUserID != userID && !ds.userIsAdmin(ctx, dialogID, adminUserID) {
		return fmt.Errorf("only admins can remove other participants")
	}

	// Cannot remove the last admin
	if ds.userIsAdmin(ctx, dialogID, userID) {
		admins, err := ds.dialogRepo.GetAdmins(ctx, dialogID)
		if err != nil {
			return fmt.Errorf("failed to check admin count: %w", err)
		}
		if len(admins) <= 1 {
			return fmt.Errorf("cannot remove the last admin")
		}
	}

	// Remove participant
	if err := ds.dialogRepo.RemoveParticipant(ctx, dialogID, userID); err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	// Update participant count
	dialog, err := ds.dialogRepo.GetByID(ctx, dialogID)
	if err == nil {
		dialog.ParticipantCount--
		dialog.UpdatedAt = time.Now()
		ds.dialogRepo.Update(ctx, dialog)
	}

	// Publish participant removed event
	if err := ds.publishDialogEvent(ctx, "dialog.participant_removed", dialogID, adminUserID, map[string]interface{}{
		"participant_id": userID,
		"removed_by":     adminUserID,
		"is_self_leave":  adminUserID == userID,
	}); err != nil {
		fmt.Printf("Failed to publish participant removed event: %v\n", err)
	}

	return nil
}

func (ds *DialogService) PromoteParticipant(ctx context.Context, dialogID uuid.UUID, adminUserID uuid.UUID, userID uuid.UUID) error {
	// Check if user is admin
	if !ds.userIsAdmin(ctx, dialogID, adminUserID) {
		return fmt.Errorf("only admins can promote participants")
	}

	// Get participant
	participants, err := ds.dialogRepo.GetParticipants(ctx, dialogID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %w", err)
	}

	var targetParticipant *models.DialogParticipant
	for _, participant := range participants {
		if participant.UserID == userID {
			targetParticipant = participant
			break
		}
	}

	if targetParticipant == nil {
		return fmt.Errorf("user is not a participant")
	}

	if targetParticipant.Role == models.ParticipantRoleAdmin {
		return fmt.Errorf("user is already an admin")
	}

	// Promote to admin
	targetParticipant.Role = models.ParticipantRoleAdmin
	targetParticipant.UpdatedAt = time.Now()

	if err := ds.dialogRepo.UpdateParticipant(ctx, targetParticipant); err != nil {
		return fmt.Errorf("failed to promote participant: %w", err)
	}

	// Publish participant promoted event
	if err := ds.publishDialogEvent(ctx, "dialog.participant_promoted", dialogID, adminUserID, map[string]interface{}{
		"participant_id": userID,
		"promoted_by":    adminUserID,
		"new_role":       models.ParticipantRoleAdmin,
	}); err != nil {
		fmt.Printf("Failed to publish participant promoted event: %v\n", err)
	}

	return nil
}

func (ds *DialogService) ArchiveDialog(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) error {
	// Check if user is admin
	if !ds.userIsAdmin(ctx, dialogID, userID) {
		return fmt.Errorf("only admins can archive dialog")
	}

	dialog, err := ds.dialogRepo.GetByID(ctx, dialogID)
	if err != nil {
		return fmt.Errorf("dialog not found")
	}

	if dialog.IsArchived {
		return fmt.Errorf("dialog is already archived")
	}

	dialog.IsArchived = true
	dialog.UpdatedAt = time.Now()

	if err := ds.dialogRepo.Update(ctx, dialog); err != nil {
		return fmt.Errorf("failed to archive dialog: %w", err)
	}

	// Publish dialog archived event
	if err := ds.publishDialogEvent(ctx, "dialog.archived", dialogID, userID, map[string]interface{}{
		"archived_by": userID,
	}); err != nil {
		fmt.Printf("Failed to publish dialog archived event: %v\n", err)
	}

	return nil
}

func (ds *DialogService) SearchDialogs(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Dialog, error) {
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	dialogs, err := ds.dialogRepo.SearchDialogs(ctx, userID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search dialogs: %w", err)
	}

	return dialogs, nil
}

func (ds *DialogService) GetDialogParticipants(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) ([]*models.DialogParticipant, error) {
	// Check if user has access to this dialog
	if !ds.userHasAccess(ctx, dialogID, userID) {
		return nil, fmt.Errorf("access denied")
	}

	participants, err := ds.dialogRepo.GetParticipants(ctx, dialogID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	return participants, nil
}

// Private helper methods

func (ds *DialogService) userHasAccess(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) bool {
	participants, err := ds.dialogRepo.GetParticipants(ctx, dialogID)
	if err != nil {
		return false
	}

	for _, participant := range participants {
		if participant.UserID == userID && participant.Status == models.ParticipantStatusActive {
			return true
		}
	}

	return false
}

func (ds *DialogService) userIsAdmin(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) bool {
	participants, err := ds.dialogRepo.GetParticipants(ctx, dialogID)
	if err != nil {
		return false
	}

	for _, participant := range participants {
		if participant.UserID == userID &&
		   participant.Status == models.ParticipantStatusActive &&
		   (participant.Role == models.ParticipantRoleAdmin || participant.Role == models.ParticipantRoleOwner) {
			return true
		}
	}

	return false
}

func (ds *DialogService) validateCreateDialogRequest(req *CreateDialogRequest) error {
	if req.CreatorID == uuid.Nil {
		return fmt.Errorf("creator ID is required")
	}

	if req.Type == "" {
		return fmt.Errorf("dialog type is required")
	}

	if req.Title == "" {
		return fmt.Errorf("dialog title is required")
	}

	if len(req.ParticipantIDs) == 0 {
		req.ParticipantIDs = []uuid.UUID{req.CreatorID}
	}

	// Check for duplicate participant IDs
	seen := make(map[uuid.UUID]bool)
	for _, id := range req.ParticipantIDs {
		if seen[id] {
			return fmt.Errorf("duplicate participant ID: %s", id)
		}
		seen[id] = true
	}

	return nil
}

func (ds *DialogService) publishDialogEvent(ctx context.Context, eventType sharedModels.EventType, dialogID uuid.UUID, userID uuid.UUID, data map[string]interface{}) error {
	event := &sharedModels.Event{
		ID:            uuid.New(),
		Type:          eventType,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       fmt.Sprintf("Dialog event: %s", eventType),
		AggregateID:   dialogID.String(),
		AggregateType: "dialog",
		EventVersion:  1,
		OccurredAt:    time.Now(),
		Status:        sharedModels.EventStatusPending,
		Metadata: sharedModels.EventMetadata{
			Source:      "messaging-service",
			Environment: "production",
			Region:      "sea",
		},
	}

	// Add user context to data
	data["user_id"] = userID

	if err := event.MarshalData(data); err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	return ds.eventPublisher.Publish(ctx, event)
}

// Request/Response structures

type CreateDialogRequest struct {
	Type           models.DialogType      `json:"type" binding:"required"`
	Title          string                 `json:"title" binding:"required"`
	Description    string                 `json:"description"`
	CreatorID      uuid.UUID              `json:"creator_id" binding:"required"`
	ParticipantIDs []uuid.UUID            `json:"participant_ids"`
	Settings       models.DialogSettings  `json:"settings"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type UpdateDialogRequest struct {
	Title       *string                 `json:"title"`
	Description *string                 `json:"description"`
	Settings    *models.DialogSettings  `json:"settings"`
	Metadata    map[string]interface{}  `json:"metadata"`
}

type AddParticipantRequest struct {
	UserID uuid.UUID                   `json:"user_id" binding:"required"`
	Role   models.ParticipantRole      `json:"role"`
}

type DialogResponse struct {
	ID               uuid.UUID              `json:"id"`
	Type             models.DialogType      `json:"type"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	ParticipantCount int                    `json:"participant_count"`
	UnreadCount      int                    `json:"unread_count"`
	LastMessage      *LastMessageInfo       `json:"last_message,omitempty"`
	Settings         models.DialogSettings  `json:"settings"`
	IsArchived       bool                   `json:"is_archived"`
	IsMuted          bool                   `json:"is_muted"`
	CreatedBy        uuid.UUID              `json:"created_by"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type LastMessageInfo struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	SenderID  uuid.UUID `json:"sender_id"`
	SentAt    time.Time `json:"sent_at"`
}

type DialogParticipantResponse struct {
	ID         uuid.UUID                 `json:"id"`
	UserID     uuid.UUID                 `json:"user_id"`
	Role       models.ParticipantRole    `json:"role"`
	Status     models.ParticipantStatus  `json:"status"`
	JoinedAt   time.Time                 `json:"joined_at"`
	LastSeen   *time.Time                `json:"last_seen,omitempty"`
}

type DialogListResponse struct {
	Dialogs    []*DialogResponse `json:"dialogs"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

func (dialog *models.Dialog) ToResponse() *DialogResponse {
	return &DialogResponse{
		ID:               dialog.ID,
		Type:             dialog.Type,
		Title:            dialog.Title,
		Description:      dialog.Description,
		ParticipantCount: dialog.ParticipantCount,
		UnreadCount:      dialog.UnreadCount,
		Settings:         dialog.Settings,
		IsArchived:       dialog.IsArchived,
		IsMuted:          dialog.IsMuted,
		CreatedBy:        dialog.CreatedBy,
		CreatedAt:        dialog.CreatedAt,
		UpdatedAt:        dialog.UpdatedAt,
	}
}

func (participant *models.DialogParticipant) ToResponse() *DialogParticipantResponse {
	return &DialogParticipantResponse{
		ID:       participant.ID,
		UserID:   participant.UserID,
		Role:     participant.Role,
		Status:   participant.Status,
		JoinedAt: participant.JoinedAt,
		LastSeen: participant.LastSeenAt,
	}
}