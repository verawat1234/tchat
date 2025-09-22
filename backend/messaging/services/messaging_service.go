package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"tchat.dev/messaging/models"
)

// MessagingService provides messaging functionality
type MessagingService struct {
	dialogRepo  DialogRepository
	messageRepo MessageRepository
	cache       CacheService
	events      EventService
	config      *MessagingConfig
}

// MessagingConfig holds messaging service configuration
type MessagingConfig struct {
	MaxMessageLength       int
	MaxParticipants       int
	DefaultRetentionPeriod time.Duration
	MaxFileSize           int64
	AllowedFileTypes      []string
	EnableEncryption      bool
	MaxDialogsPerUser     int
}

// DefaultMessagingConfig returns default messaging configuration
func DefaultMessagingConfig() *MessagingConfig {
	return &MessagingConfig{
		MaxMessageLength:       4000,
		MaxParticipants:       200000, // Broadcast channels
		DefaultRetentionPeriod: 365 * 24 * time.Hour, // 1 year
		MaxFileSize:           50 * 1024 * 1024, // 50MB
		AllowedFileTypes:      []string{"jpg", "jpeg", "png", "gif", "mp4", "mp3", "pdf", "doc", "docx"},
		EnableEncryption:      true,
		MaxDialogsPerUser:     1000,
	}
}

// DialogRepository interface for dialog data access
type DialogRepository interface {
	Create(ctx context.Context, dialog *models.Dialog) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Dialog, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Dialog, error)
	Update(ctx context.Context, dialog *models.Dialog) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddParticipant(ctx context.Context, dialogID, userID uuid.UUID, role models.ParticipantRole) error
	RemoveParticipant(ctx context.Context, dialogID, userID uuid.UUID) error
	GetParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error)
	UpdateParticipant(ctx context.Context, participant *models.DialogParticipant) error
}

// MessageRepository interface for message data access
type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	GetByDialogID(ctx context.Context, dialogID uuid.UUID, limit, offset int) ([]*models.Message, error)
	GetThreadMessages(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*models.Message, error)
	Update(ctx context.Context, message *models.Message) error
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	GetUnreadCount(ctx context.Context, dialogID, userID uuid.UUID) (int64, error)
	MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error
	AddReaction(ctx context.Context, messageID, userID uuid.UUID, reaction string) error
	RemoveReaction(ctx context.Context, messageID, userID uuid.UUID, reaction string) error
}

// CacheService interface for caching operations
type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
	SetWithExpiry(ctx context.Context, key string, value interface{}, expiry time.Duration) error
}

// EventService interface for event publishing
type EventService interface {
	PublishMessage(ctx context.Context, event *MessageEvent) error
	PublishDialogUpdate(ctx context.Context, event *DialogEvent) error
}

// MessageEvent represents a message-related event
type MessageEvent struct {
	Type      string         `json:"type"`
	DialogID  uuid.UUID      `json:"dialog_id"`
	MessageID uuid.UUID      `json:"message_id"`
	SenderID  uuid.UUID      `json:"sender_id"`
	Message   *models.Message `json:"message,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// DialogEvent represents a dialog-related event
type DialogEvent struct {
	Type      string         `json:"type"`
	DialogID  uuid.UUID      `json:"dialog_id"`
	UserID    uuid.UUID      `json:"user_id"`
	Dialog    *models.Dialog `json:"dialog,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// CreateDialogRequest represents a request to create a dialog
type CreateDialogRequest struct {
	Type          models.DialogType `json:"type"`
	Name          *string           `json:"name,omitempty"`
	Description   *string           `json:"description,omitempty"`
	CreatorID     uuid.UUID         `json:"creator_id"`
	ParticipantIDs []uuid.UUID       `json:"participant_ids,omitempty"`
	IsPublic      bool              `json:"is_public"`
	Settings      *models.DialogSettings `json:"settings,omitempty"`
}

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	DialogID    uuid.UUID           `json:"dialog_id"`
	SenderID    uuid.UUID           `json:"sender_id"`
	Type        models.MessageType  `json:"type"`
	Content     models.MessageContent `json:"content"`
	ParentID    *uuid.UUID          `json:"parent_id,omitempty"`
	Mentions    []uuid.UUID         `json:"mentions,omitempty"`
	ReplyToID   *uuid.UUID          `json:"reply_to_id,omitempty"`
}

// NewMessagingService creates a new messaging service
func NewMessagingService(
	dialogRepo DialogRepository,
	messageRepo MessageRepository,
	cache CacheService,
	events EventService,
	config *MessagingConfig,
) *MessagingService {
	if config == nil {
		config = DefaultMessagingConfig()
	}

	return &MessagingService{
		dialogRepo:  dialogRepo,
		messageRepo: messageRepo,
		cache:       cache,
		events:      events,
		config:      config,
	}
}

// CreateDialog creates a new dialog/conversation
func (m *MessagingService) CreateDialog(ctx context.Context, req *CreateDialogRequest) (*models.Dialog, error) {
	// Validate request
	if err := m.validateCreateDialogRequest(req); err != nil {
		return nil, fmt.Errorf("invalid create dialog request: %v", err)
	}

	// Check user's dialog limit
	userDialogs, err := m.dialogRepo.GetByUserID(ctx, req.CreatorID, m.config.MaxDialogsPerUser, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check user dialogs: %v", err)
	}
	if len(userDialogs) >= m.config.MaxDialogsPerUser {
		return nil, fmt.Errorf("user has reached maximum number of dialogs (%d)", m.config.MaxDialogsPerUser)
	}

	// Create dialog
	dialog := &models.Dialog{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   req.CreatorID,
		IsPublic:    req.IsPublic,
	}

	// Set default settings if not provided
	if req.Settings != nil {
		dialog.Settings = *req.Settings
	} else {
		dialog.Settings = models.DialogSettings{
			AllowInvites:       true,
			MessageHistory:     true,
			AllowFileSharing:   true,
			AllowVoiceMessages: true,
			AutoDeleteAfter:    nil,
			OnlyAdminsCanPost:  false,
			RequireApproval:    false,
		}
	}

	// Set up dialog before creation
	if err := dialog.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("failed to prepare dialog: %v", err)
	}

	// Validate dialog
	if err := dialog.Validate(); err != nil {
		return nil, fmt.Errorf("dialog validation failed: %v", err)
	}

	// Save dialog
	if err := m.dialogRepo.Create(ctx, dialog); err != nil {
		return nil, fmt.Errorf("failed to create dialog: %v", err)
	}

	// Add creator as admin participant
	if err := m.dialogRepo.AddParticipant(ctx, dialog.ID, req.CreatorID, models.ParticipantRoleAdmin); err != nil {
		return nil, fmt.Errorf("failed to add creator as participant: %v", err)
	}

	// Add other participants
	for _, participantID := range req.ParticipantIDs {
		if participantID == req.CreatorID {
			continue // Skip creator as already added
		}

		role := models.ParticipantRoleMember
		if dialog.Type == models.DialogTypeUser {
			role = models.ParticipantRoleMember // Both users are equal in private chat
		}

		if err := m.dialogRepo.AddParticipant(ctx, dialog.ID, participantID, role); err != nil {
			// Log error but continue with other participants
			fmt.Printf("Warning: failed to add participant %s to dialog %s: %v\n", participantID, dialog.ID, err)
		}
	}

	// Update participant count
	participants, err := m.dialogRepo.GetParticipants(ctx, dialog.ID)
	if err == nil {
		dialog.ParticipantCount = len(participants)
		m.dialogRepo.Update(ctx, dialog)
	}

	// Publish dialog creation event
	event := &DialogEvent{
		Type:      "dialog_created",
		DialogID:  dialog.ID,
		UserID:    req.CreatorID,
		Dialog:    dialog,
		Timestamp: time.Now().UTC(),
	}
	m.events.PublishDialogUpdate(ctx, event)

	return dialog, nil
}

// SendMessage sends a message to a dialog
func (m *MessagingService) SendMessage(ctx context.Context, req *SendMessageRequest) (*models.Message, error) {
	// Validate request
	if err := m.validateSendMessageRequest(req); err != nil {
		return nil, fmt.Errorf("invalid send message request: %v", err)
	}

	// Get dialog
	dialog, err := m.dialogRepo.GetByID(ctx, req.DialogID)
	if err != nil {
		return nil, fmt.Errorf("dialog not found: %v", err)
	}

	// Check if user is participant
	participants, err := m.dialogRepo.GetParticipants(ctx, req.DialogID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %v", err)
	}

	userParticipant := m.findParticipant(participants, req.SenderID)
	if userParticipant == nil {
		return nil, fmt.Errorf("user is not a participant in this dialog")
	}

	// Check posting permissions
	if dialog.Settings.OnlyAdminsCanPost && userParticipant.Role != models.ParticipantRoleAdmin {
		return nil, fmt.Errorf("only admins can post in this dialog")
	}

	// Create message
	message := &models.Message{
		DialogID:  req.DialogID,
		SenderID:  req.SenderID,
		Type:      req.Type,
		Content:   req.Content,
		ParentID:  req.ParentID,
		Mentions:  req.Mentions,
		ReplyToID: req.ReplyToID,
	}

	// Set up message before creation
	if err := message.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("failed to prepare message: %v", err)
	}

	// Validate message
	if err := message.Validate(); err != nil {
		return nil, fmt.Errorf("message validation failed: %v", err)
	}

	// Additional content validation based on type
	if err := m.validateMessageContent(message); err != nil {
		return nil, fmt.Errorf("message content validation failed: %v", err)
	}

	// Save message
	if err := m.messageRepo.Create(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %v", err)
	}

	// Update dialog's last message info
	dialog.LastMessageID = &message.ID
	dialog.LastMessageAt = &message.CreatedAt
	dialog.UpdatedAt = message.CreatedAt
	if err := m.dialogRepo.Update(ctx, dialog); err != nil {
		// Log error but don't fail the message sending
		fmt.Printf("Warning: failed to update dialog last message: %v\n", err)
	}

	// Publish message event
	event := &MessageEvent{
		Type:      "message_sent",
		DialogID:  req.DialogID,
		MessageID: message.ID,
		SenderID:  req.SenderID,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
	m.events.PublishMessage(ctx, event)

	return message, nil
}

// GetDialog retrieves a dialog by ID
func (m *MessagingService) GetDialog(ctx context.Context, dialogID, userID uuid.UUID) (*models.Dialog, error) {
	// Get dialog
	dialog, err := m.dialogRepo.GetByID(ctx, dialogID)
	if err != nil {
		return nil, fmt.Errorf("dialog not found: %v", err)
	}

	// Check if user is participant (for private dialogs)
	if !dialog.IsPublic {
		participants, err := m.dialogRepo.GetParticipants(ctx, dialogID)
		if err != nil {
			return nil, fmt.Errorf("failed to get participants: %v", err)
		}

		userParticipant := m.findParticipant(participants, userID)
		if userParticipant == nil {
			return nil, fmt.Errorf("access denied: user is not a participant")
		}
	}

	return dialog, nil
}

// GetMessages retrieves messages from a dialog
func (m *MessagingService) GetMessages(ctx context.Context, dialogID, userID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	// Verify user access to dialog
	_, err := m.GetDialog(ctx, dialogID, userID)
	if err != nil {
		return nil, err
	}

	// Get messages
	messages, err := m.messageRepo.GetByDialogID(ctx, dialogID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %v", err)
	}

	return messages, nil
}

// AddParticipant adds a user to a dialog
func (m *MessagingService) AddParticipant(ctx context.Context, dialogID, userID, inviterID uuid.UUID, role models.ParticipantRole) error {
	// Get dialog
	dialog, err := m.dialogRepo.GetByID(ctx, dialogID)
	if err != nil {
		return fmt.Errorf("dialog not found: %v", err)
	}

	// Check inviter permissions
	participants, err := m.dialogRepo.GetParticipants(ctx, dialogID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %v", err)
	}

	inviterParticipant := m.findParticipant(participants, inviterID)
	if inviterParticipant == nil {
		return fmt.Errorf("inviter is not a participant in this dialog")
	}

	// Check if invites are allowed
	if !dialog.Settings.AllowInvites && inviterParticipant.Role != models.ParticipantRoleAdmin {
		return fmt.Errorf("only admins can invite participants when invites are disabled")
	}

	// Check participant limit
	if len(participants) >= m.config.MaxParticipants {
		return fmt.Errorf("dialog has reached maximum number of participants")
	}

	// Check if user is already a participant
	if m.findParticipant(participants, userID) != nil {
		return fmt.Errorf("user is already a participant in this dialog")
	}

	// Add participant
	if err := m.dialogRepo.AddParticipant(ctx, dialogID, userID, role); err != nil {
		return fmt.Errorf("failed to add participant: %v", err)
	}

	// Update participant count
	dialog.ParticipantCount = len(participants) + 1
	dialog.UpdatedAt = time.Now().UTC()
	if err := m.dialogRepo.Update(ctx, dialog); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update participant count: %v\n", err)
	}

	// Publish participant added event
	event := &DialogEvent{
		Type:      "participant_added",
		DialogID:  dialogID,
		UserID:    userID,
		Timestamp: time.Now().UTC(),
	}
	m.events.PublishDialogUpdate(ctx, event)

	return nil
}

// RemoveParticipant removes a user from a dialog
func (m *MessagingService) RemoveParticipant(ctx context.Context, dialogID, userID, removerID uuid.UUID) error {
	// Get dialog
	dialog, err := m.dialogRepo.GetByID(ctx, dialogID)
	if err != nil {
		return fmt.Errorf("dialog not found: %v", err)
	}

	// Get participants
	participants, err := m.dialogRepo.GetParticipants(ctx, dialogID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %v", err)
	}

	removerParticipant := m.findParticipant(participants, removerID)
	userParticipant := m.findParticipant(participants, userID)

	if removerParticipant == nil {
		return fmt.Errorf("remover is not a participant in this dialog")
	}

	if userParticipant == nil {
		return fmt.Errorf("user is not a participant in this dialog")
	}

	// Check permissions (admin can remove anyone, users can remove themselves)
	if removerID != userID && removerParticipant.Role != models.ParticipantRoleAdmin {
		return fmt.Errorf("insufficient permissions to remove participant")
	}

	// Prevent removing the creator unless they're removing themselves
	if userID == dialog.CreatorID && removerID != userID {
		return fmt.Errorf("cannot remove dialog creator")
	}

	// Remove participant
	if err := m.dialogRepo.RemoveParticipant(ctx, dialogID, userID); err != nil {
		return fmt.Errorf("failed to remove participant: %v", err)
	}

	// Update participant count
	dialog.ParticipantCount = len(participants) - 1
	dialog.UpdatedAt = time.Now().UTC()
	if err := m.dialogRepo.Update(ctx, dialog); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update participant count: %v\n", err)
	}

	// Publish participant removed event
	event := &DialogEvent{
		Type:      "participant_removed",
		DialogID:  dialogID,
		UserID:    userID,
		Timestamp: time.Now().UTC(),
	}
	m.events.PublishDialogUpdate(ctx, event)

	return nil
}

// AddReaction adds a reaction to a message
func (m *MessagingService) AddReaction(ctx context.Context, messageID, userID uuid.UUID, reaction string) error {
	// Get message
	message, err := m.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message not found: %v", err)
	}

	// Verify user access to dialog
	_, err = m.GetDialog(ctx, message.DialogID, userID)
	if err != nil {
		return err
	}

	// Validate reaction
	if strings.TrimSpace(reaction) == "" {
		return fmt.Errorf("reaction cannot be empty")
	}

	// Add reaction
	if err := m.messageRepo.AddReaction(ctx, messageID, userID, reaction); err != nil {
		return fmt.Errorf("failed to add reaction: %v", err)
	}

	return nil
}

// MarkAsRead marks a message as read
func (m *MessagingService) MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error {
	// Get message
	message, err := m.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message not found: %v", err)
	}

	// Verify user access to dialog
	_, err = m.GetDialog(ctx, message.DialogID, userID)
	if err != nil {
		return err
	}

	// Mark as read
	if err := m.messageRepo.MarkAsRead(ctx, messageID, userID); err != nil {
		return fmt.Errorf("failed to mark message as read: %v", err)
	}

	return nil
}

// GetUnreadCount gets the number of unread messages in a dialog for a user
func (m *MessagingService) GetUnreadCount(ctx context.Context, dialogID, userID uuid.UUID) (int64, error) {
	// Verify user access to dialog
	_, err := m.GetDialog(ctx, dialogID, userID)
	if err != nil {
		return 0, err
	}

	// Get unread count
	count, err := m.messageRepo.GetUnreadCount(ctx, dialogID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %v", err)
	}

	return count, nil
}

// Helper methods

func (m *MessagingService) validateCreateDialogRequest(req *CreateDialogRequest) error {
	if !req.Type.IsValid() {
		return fmt.Errorf("invalid dialog type: %s", req.Type)
	}

	if req.CreatorID == uuid.Nil {
		return errors.New("creator_id is required")
	}

	// Validate participant limits based on dialog type
	maxParticipants := m.getMaxParticipantsForType(req.Type)
	if len(req.ParticipantIDs) > maxParticipants {
		return fmt.Errorf("too many participants: max %d allowed for %s", maxParticipants, req.Type)
	}

	// Validate dialog name for group/channel types
	if (req.Type == models.DialogTypeGroup || req.Type == models.DialogTypeChannel || req.Type == models.DialogTypeBusiness) {
		if req.Name == nil || strings.TrimSpace(*req.Name) == "" {
			return fmt.Errorf("name is required for %s dialogs", req.Type)
		}
	}

	return nil
}

func (m *MessagingService) validateSendMessageRequest(req *SendMessageRequest) error {
	if req.DialogID == uuid.Nil {
		return errors.New("dialog_id is required")
	}

	if req.SenderID == uuid.Nil {
		return errors.New("sender_id is required")
	}

	if !req.Type.IsValid() {
		return fmt.Errorf("invalid message type: %s", req.Type)
	}

	return nil
}

func (m *MessagingService) validateMessageContent(message *models.Message) error {
	switch message.Type {
	case models.MessageTypeText:
		if message.Content.Text == nil || strings.TrimSpace(*message.Content.Text) == "" {
			return errors.New("text content is required for text messages")
		}
		if len(*message.Content.Text) > m.config.MaxMessageLength {
			return fmt.Errorf("message text exceeds maximum length of %d characters", m.config.MaxMessageLength)
		}

	case models.MessageTypeFile, models.MessageTypeImage, models.MessageTypeVideo:
		if message.Content.File == nil {
			return fmt.Errorf("file content is required for %s messages", message.Type)
		}
		if message.Content.File.Size > m.config.MaxFileSize {
			return fmt.Errorf("file size exceeds maximum limit of %d bytes", m.config.MaxFileSize)
		}

	case models.MessageTypeVoice:
		if message.Content.Voice == nil {
			return errors.New("voice content is required for voice messages")
		}
		if message.Content.Voice.Duration <= 0 {
			return errors.New("voice message duration must be positive")
		}
	}

	return nil
}

func (m *MessagingService) findParticipant(participants []*models.DialogParticipant, userID uuid.UUID) *models.DialogParticipant {
	for _, participant := range participants {
		if participant.UserID == userID {
			return participant
		}
	}
	return nil
}

func (m *MessagingService) getMaxParticipantsForType(dialogType models.DialogType) int {
	switch dialogType {
	case models.DialogTypeUser:
		return 2
	case models.DialogTypeGroup:
		return 1000
	case models.DialogTypeChannel:
		return 200000
	case models.DialogTypeBusiness:
		return 10000
	default:
		return 1000
	}
}