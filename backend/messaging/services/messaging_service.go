package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tchat.dev/messaging/models"
)

// MessagingService interface defines messaging service operations
type MessagingService interface {
	CreateDialog(ctx context.Context, req *CreateDialogRequest) (*models.Dialog, error)
	SendMessage(ctx context.Context, req *SendMessageRequest) (*models.Message, error)
	GetDialogByID(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (*models.Dialog, error)
	GetMessageByID(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*models.Message, error)
	GetDialogParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error)
	GetUserDialogs(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Dialog, error)
	GetMessages(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID, limit, offset int) ([]*models.Message, error)
	AddReaction(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, reaction string) error
	MarkAsRead(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (int, error)
}

// MessagingServiceImpl provides messaging functionality
type MessagingServiceImpl struct {
	dialogService  *DialogService
	messageService *MessageService
	cache          CacheService
	events         EventService
	config         *MessagingConfig
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


// NewMessagingService creates a new messaging service
func NewMessagingService(
	dialogService *DialogService,
	messageService *MessageService,
	cache CacheService,
	events EventService,
	config *MessagingConfig,
) MessagingService {
	if config == nil {
		config = DefaultMessagingConfig()
	}

	return &MessagingServiceImpl{
		dialogService:  dialogService,
		messageService: messageService,
		cache:          cache,
		events:         events,
		config:         config,
	}
}

// CreateDialog delegates dialog creation to the dialog service
func (m *MessagingServiceImpl) CreateDialog(ctx context.Context, req *CreateDialogRequest) (*models.Dialog, error) {
	return m.dialogService.CreateDialog(ctx, req)
}

// SendMessage delegates message sending to the message service
func (m *MessagingServiceImpl) SendMessage(ctx context.Context, req *SendMessageRequest) (*models.Message, error) {
	return m.messageService.SendMessage(ctx, req)
}

// GetDialogByID delegates to dialog service
func (m *MessagingServiceImpl) GetDialogByID(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (*models.Dialog, error) {
	return m.dialogService.GetDialogByID(ctx, dialogID, userID)
}

// GetMessageByID delegates to message service
func (m *MessagingServiceImpl) GetMessageByID(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*models.Message, error) {
	return m.messageService.GetMessageByID(ctx, messageID, userID)
}

// GetDialogParticipants delegates to dialog service
func (m *MessagingServiceImpl) GetDialogParticipants(ctx context.Context, dialogID uuid.UUID) ([]*models.DialogParticipant, error) {
	// Create a dummy user ID for the internal call since this method doesn't require user validation
	dummyUserID := uuid.New()
	return m.dialogService.GetDialogParticipants(ctx, dialogID, dummyUserID)
}

// GetUserDialogs delegates to dialog service
func (m *MessagingServiceImpl) GetUserDialogs(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Dialog, error) {
	filters := DialogFilters{}
	pagination := Pagination{
		Page:     (offset / limit) + 1,
		PageSize: limit,
		OrderBy:  "updated_at",
		Order:    "desc",
	}
	dialogs, _, err := m.dialogService.GetUserDialogs(ctx, userID, filters, pagination)
	return dialogs, err
}

// GetMessages delegates to message service
func (m *MessagingServiceImpl) GetMessages(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	filters := MessageFilters{}
	pagination := Pagination{
		Page:     (offset / limit) + 1,
		PageSize: limit,
		OrderBy:  "sent_at",
		Order:    "desc",
	}
	messages, _, err := m.messageService.GetDialogMessages(ctx, dialogID, userID, filters, pagination)
	return messages, err
}

// AddReaction delegates to message service
func (m *MessagingServiceImpl) AddReaction(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, reaction string) error {
	// This would need to be implemented in the message service
	// For now, return not implemented
	return fmt.Errorf("add reaction not implemented yet")
}

// MarkAsRead delegates to message service
func (m *MessagingServiceImpl) MarkAsRead(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error {
	return m.messageService.MarkAsRead(ctx, messageID, userID)
}

// GetUnreadCount delegates to message service
func (m *MessagingServiceImpl) GetUnreadCount(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (int, error) {
	return m.messageService.GetUnreadCount(ctx, dialogID, userID)
}