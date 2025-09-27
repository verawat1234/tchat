package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/messaging/models"
	sharedModels "tchat.dev/shared/models"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	GetByDialogID(ctx context.Context, dialogID uuid.UUID, filters MessageFilters, pagination Pagination) ([]*models.Message, int64, error)
	Update(ctx context.Context, message *models.Message) error
	Delete(ctx context.Context, id uuid.UUID) error
	MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error
	MarkAsDelivered(ctx context.Context, messageID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, dialogID, userID uuid.UUID) (int, error)
	SearchMessages(ctx context.Context, dialogID uuid.UUID, query string, limit int) ([]*models.Message, error)
	GetMessageStats(ctx context.Context, dialogID uuid.UUID) (*MessageStats, error)
}

type MessageDeliveryService interface {
	DeliverMessage(ctx context.Context, message *models.Message, recipientIDs []uuid.UUID) error
	SendPushNotification(ctx context.Context, userID uuid.UUID, message *models.Message) error
}

type ContentModerator interface {
	ModerateContent(ctx context.Context, content string, contentType models.MessageType) (*ModerationResult, error)
	DetectSpam(ctx context.Context, senderID uuid.UUID, content string) (*SpamDetectionResult, error)
}

type MediaProcessor interface {
	ProcessImageUpload(ctx context.Context, imageData []byte, metadata map[string]interface{}) (*ProcessedMedia, error)
	ProcessVideoUpload(ctx context.Context, videoData []byte, metadata map[string]interface{}) (*ProcessedMedia, error)
	ProcessAudioUpload(ctx context.Context, audioData []byte, metadata map[string]interface{}) (*ProcessedMedia, error)
	GenerateThumbnail(ctx context.Context, mediaURL string, mediaType string) (string, error)
}

type MessageFilters struct {
	Type        *models.MessageType `json:"type,omitempty"`
	SenderID    *uuid.UUID          `json:"sender_id,omitempty"`
	IsEdited    *bool               `json:"is_edited,omitempty"`
	IsDeleted   *bool               `json:"is_deleted,omitempty"`
	SentFrom    *time.Time          `json:"sent_from,omitempty"`
	SentTo      *time.Time          `json:"sent_to,omitempty"`
	HasMedia    *bool               `json:"has_media,omitempty"`
}

type MessageStats struct {
	TotalMessages    int64   `json:"total_messages"`
	TextMessages     int64   `json:"text_messages"`
	MediaMessages    int64   `json:"media_messages"`
	SystemMessages   int64   `json:"system_messages"`
	AverageLength    float64 `json:"average_length"`
	MessagesPerDay   int64   `json:"messages_per_day"`
	ActiveSenders    int64   `json:"active_senders"`
}

type ModerationResult struct {
	IsApproved    bool     `json:"is_approved"`
	Violations    []string `json:"violations"`
	Confidence    float64  `json:"confidence"`
	FilteredContent string `json:"filtered_content"`
}

type SpamDetectionResult struct {
	IsSpam      bool    `json:"is_spam"`
	Confidence  float64 `json:"confidence"`
	Reasons     []string `json:"reasons"`
}

type ProcessedMedia struct {
	URL         string            `json:"url"`
	ThumbnailURL string           `json:"thumbnail_url"`
	Size        int64             `json:"size"`
	Width       int               `json:"width"`
	Height      int               `json:"height"`
	Duration    *float64          `json:"duration,omitempty"`
	Format      string            `json:"format"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type MessageService struct {
	messageRepo       MessageRepository
	dialogService     *DialogService
	deliveryService   MessageDeliveryService
	contentModerator  ContentModerator
	mediaProcessor    MediaProcessor
	eventPublisher    EventPublisher
	db                *gorm.DB
}

func NewMessageService(
	messageRepo MessageRepository,
	dialogService *DialogService,
	deliveryService MessageDeliveryService,
	contentModerator ContentModerator,
	mediaProcessor MediaProcessor,
	eventPublisher EventPublisher,
	db *gorm.DB,
) *MessageService {
	return &MessageService{
		messageRepo:      messageRepo,
		dialogService:    dialogService,
		deliveryService:  deliveryService,
		contentModerator: contentModerator,
		mediaProcessor:   mediaProcessor,
		eventPublisher:   eventPublisher,
		db:               db,
	}
}

func (ms *MessageService) SendMessage(ctx context.Context, req *SendMessageRequest) (*models.Message, error) {
	// Validate request
	if err := ms.validateSendMessageRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user has access to dialog
	_, err := ms.dialogService.GetDialogByID(ctx, req.DialogID, req.SenderID)
	if err != nil {
		return nil, fmt.Errorf("access denied or dialog not found: %w", err)
	}

	// Moderate content if it's text message
	if req.Type == models.MessageTypeText && req.Content != "" {
		moderation, err := ms.contentModerator.ModerateContent(ctx, req.Content, req.Type)
		if err != nil {
			return nil, fmt.Errorf("content moderation failed: %w", err)
		}

		if !moderation.IsApproved {
			return nil, fmt.Errorf("message content violates community guidelines: %v", moderation.Violations)
		}

		// Use filtered content if available
		if moderation.FilteredContent != "" {
			req.Content = moderation.FilteredContent
		}
	}

	// Check for spam
	spam, err := ms.contentModerator.DetectSpam(ctx, req.SenderID, req.Content)
	if err != nil {
		fmt.Printf("Spam detection failed: %v\n", err)
	} else if spam.IsSpam {
		return nil, fmt.Errorf("message detected as spam")
	}

	// Process media if present
	if req.MediaData != nil {
		processed, err := ms.processMediaContent(ctx, req.Type, req.MediaData, req.MediaMetadata)
		if err != nil {
			return nil, fmt.Errorf("media processing failed: %w", err)
		}
		req.MediaURL = processed.URL
		req.ThumbnailURL = processed.ThumbnailURL
		req.MediaMetadata = processed.Metadata
	}

	// Create message
	message := &models.Message{
		ID:           uuid.New(),
		DialogID:     req.DialogID,
		SenderID:     req.SenderID,
		Type:         req.Type,
		Content:      models.MessageContent{"text": req.Content},
		MediaURL:     &req.MediaURL,
		ThumbnailURL: &req.ThumbnailURL,
		Status:       models.MessageStatusSent,
		Metadata:     req.Metadata,
		SentAt:       time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Validate message
	if err := message.Validate(); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// Handle reply context
	if req.ReplyToID != nil {
		replyToMessage, err := ms.messageRepo.GetByID(ctx, *req.ReplyToID)
		if err == nil && replyToMessage.DialogID == req.DialogID {
			message.ReplyToID = req.ReplyToID
			message.ReplyTo = &models.MessageReply{
				MessageID:   replyToMessage.ID,
				SenderID:    replyToMessage.SenderID,
				Content:     ms.truncateContent(replyToMessage.Content, 100),
				MessageType: string(replyToMessage.Type),
			}
		}
	}

	// Handle forward context - temporarily disabled until model supports it
	// TODO: Add ForwardFromID and ForwardFrom fields to Message model
	_ = req.ForwardFromID // Suppress unused variable warning

	// Save message
	if err := ms.messageRepo.Create(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Get dialog participants for delivery
	participants, err := ms.dialogService.GetDialogParticipants(ctx, req.DialogID, req.SenderID)
	if err != nil {
		fmt.Printf("Failed to get dialog participants for message delivery: %v\n", err)
	} else {
		// Extract recipient IDs (exclude sender)
		recipientIDs := make([]uuid.UUID, 0)
		for _, participant := range participants {
			if participant.UserID != req.SenderID && participant.Status == models.ParticipantStatusActive {
				recipientIDs = append(recipientIDs, participant.UserID)
			}
		}

		// Deliver message asynchronously
		go ms.deliverMessageAsync(context.Background(), message, recipientIDs)
	}

	// Publish message sent event
	if err := ms.publishMessageEvent(ctx, sharedModels.EventTypeMessageSent, message.ID, message.SenderID, map[string]interface{}{
		"message_id":  message.ID,
		"dialog_id":   message.DialogID,
		"message_type": message.Type,
		"has_media":   message.MediaURL != nil && *message.MediaURL != "",
		"is_reply":    message.ReplyToID != nil,
		"is_forward":  false, // TODO: Add forward support
	}); err != nil {
		fmt.Printf("Failed to publish message sent event: %v\n", err)
	}

	return message, nil
}

func (ms *MessageService) GetMessageByID(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (*models.Message, error) {
	if messageID == uuid.Nil {
		return nil, fmt.Errorf("message ID is required")
	}

	message, err := ms.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Check if user has access to this dialog
	_, err = ms.dialogService.GetDialogByID(ctx, message.DialogID, userID)
	if err != nil {
		return nil, fmt.Errorf("access denied")
	}

	return message, nil
}

func (ms *MessageService) GetDialogMessages(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID, filters MessageFilters, pagination Pagination) ([]*models.Message, int64, error) {
	// Check if user has access to dialog
	_, err := ms.dialogService.GetDialogByID(ctx, dialogID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("access denied or dialog not found: %w", err)
	}

	// Validate pagination
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > 100 {
		pagination.PageSize = 50
	}
	if pagination.OrderBy == "" {
		pagination.OrderBy = "sent_at"
	}
	if pagination.Order != "asc" && pagination.Order != "desc" {
		pagination.Order = "desc"
	}

	messages, total, err := ms.messageRepo.GetByDialogID(ctx, dialogID, filters, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get dialog messages: %w", err)
	}

	return messages, total, nil
}

func (ms *MessageService) EditMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, newContent string) (*models.Message, error) {
	// Get message
	message, err := ms.GetMessageByID(ctx, messageID, userID)
	if err != nil {
		return nil, err
	}

	// Check if user is the sender
	if message.SenderID != userID {
		return nil, fmt.Errorf("only the sender can edit their message")
	}

	// Check if message is editable
	if !message.CanEdit() {
		return nil, fmt.Errorf("message cannot be edited")
	}

	// Check edit time limit (e.g., 24 hours)
	if time.Since(message.SentAt) > 24*time.Hour {
		return nil, fmt.Errorf("edit time limit exceeded")
	}

	// Moderate new content
	if message.Type == models.MessageTypeText {
		moderation, err := ms.contentModerator.ModerateContent(ctx, newContent, message.Type)
		if err != nil {
			return nil, fmt.Errorf("content moderation failed: %w", err)
		}

		if !moderation.IsApproved {
			return nil, fmt.Errorf("edited content violates community guidelines: %v", moderation.Violations)
		}

		if moderation.FilteredContent != "" {
			newContent = moderation.FilteredContent
		}
	}

	// TODO: Implement edit history functionality
	// For now, just update the content directly

	// Update message content
	message.Content = models.MessageContent{"text": newContent}
	message.IsEdited = true
	message.EditedAt = &time.Time{}
	*message.EditedAt = time.Now()
	message.UpdatedAt = time.Now()

	// Save updated message
	if err := ms.messageRepo.Update(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	// Publish message edited event
	if err := ms.publishMessageEvent(ctx, "message.edited", messageID, userID, map[string]interface{}{
		"new_content": newContent,
		"is_edited":   true,
	}); err != nil {
		fmt.Printf("Failed to publish message edited event: %v\n", err)
	}

	return message, nil
}

func (ms *MessageService) DeleteMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, deleteForEveryone bool) error {
	// Get message
	message, err := ms.GetMessageByID(ctx, messageID, userID)
	if err != nil {
		return err
	}

	// Check permissions
	if deleteForEveryone {
		// Only sender or admin can delete for everyone
		if message.SenderID != userID {
			// Check if user is admin in the dialog
			// This would require dialog service method to check admin status
			return fmt.Errorf("only sender or admin can delete message for everyone")
		}
	}

	if deleteForEveryone {
		// Soft delete for everyone
		message.IsDeleted = true
		message.DeletedAt = &time.Time{}
		*message.DeletedAt = time.Now()
		// TODO: Add DeletedBy field to Message model if needed
		message.UpdatedAt = time.Now()

		if err := ms.messageRepo.Update(ctx, message); err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
	} else {
		// Delete only for this user (add to hidden messages list)
		// This would require a separate table/field to track per-user deletions
		return fmt.Errorf("delete for self not implemented yet")
	}

	// Publish message deleted event
	if err := ms.publishMessageEvent(ctx, "message.deleted", messageID, userID, map[string]interface{}{
		"delete_for_everyone": deleteForEveryone,
	}); err != nil {
		fmt.Printf("Failed to publish message deleted event: %v\n", err)
	}

	return nil
}

func (ms *MessageService) MarkAsRead(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error {
	// Get message to verify access
	message, err := ms.GetMessageByID(ctx, messageID, userID)
	if err != nil {
		return err
	}

	// Don't mark own messages as read
	if message.SenderID == userID {
		return nil
	}

	// Mark as read
	if err := ms.messageRepo.MarkAsRead(ctx, messageID, userID); err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	// Publish message read event
	if err := ms.publishMessageEvent(ctx, sharedModels.EventTypeMessageRead, messageID, userID, map[string]interface{}{
		"reader_id": userID,
		"sender_id": message.SenderID,
		"dialog_id": message.DialogID,
	}); err != nil {
		fmt.Printf("Failed to publish message read event: %v\n", err)
	}

	return nil
}

func (ms *MessageService) MarkAsDelivered(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) error {
	if err := ms.messageRepo.MarkAsDelivered(ctx, messageID, userID); err != nil {
		return fmt.Errorf("failed to mark message as delivered: %w", err)
	}

	// Publish message delivered event
	if err := ms.publishMessageEvent(ctx, sharedModels.EventTypeMessageDelivered, messageID, userID, map[string]interface{}{
		"recipient_id": userID,
	}); err != nil {
		fmt.Printf("Failed to publish message delivered event: %v\n", err)
	}

	return nil
}

func (ms *MessageService) GetUnreadCount(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (int, error) {
	// Check if user has access to dialog
	_, err := ms.dialogService.GetDialogByID(ctx, dialogID, userID)
	if err != nil {
		return 0, fmt.Errorf("access denied or dialog not found: %w", err)
	}

	count, err := ms.messageRepo.GetUnreadCount(ctx, dialogID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

func (ms *MessageService) SearchMessages(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID, query string, limit int) ([]*models.Message, error) {
	// Check if user has access to dialog
	_, err := ms.dialogService.GetDialogByID(ctx, dialogID, userID)
	if err != nil {
		return nil, fmt.Errorf("access denied or dialog not found: %w", err)
	}

	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	messages, err := ms.messageRepo.SearchMessages(ctx, dialogID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}

	return messages, nil
}

func (ms *MessageService) GetMessageStats(ctx context.Context, dialogID uuid.UUID, userID uuid.UUID) (*MessageStats, error) {
	// Check if user has access to dialog
	_, err := ms.dialogService.GetDialogByID(ctx, dialogID, userID)
	if err != nil {
		return nil, fmt.Errorf("access denied or dialog not found: %w", err)
	}

	stats, err := ms.messageRepo.GetMessageStats(ctx, dialogID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message stats: %w", err)
	}

	return stats, nil
}

// Private helper methods

func (ms *MessageService) processMediaContent(ctx context.Context, messageType models.MessageType, mediaData []byte, metadata map[string]interface{}) (*ProcessedMedia, error) {
	switch messageType {
	case models.MessageTypeImage:
		return ms.mediaProcessor.ProcessImageUpload(ctx, mediaData, metadata)
	case models.MessageTypeVideo:
		return ms.mediaProcessor.ProcessVideoUpload(ctx, mediaData, metadata)
	case models.MessageTypeVoice:
		return ms.mediaProcessor.ProcessAudioUpload(ctx, mediaData, metadata)
	default:
		return nil, fmt.Errorf("unsupported media type: %s", messageType)
	}
}

func (ms *MessageService) deliverMessageAsync(ctx context.Context, message *models.Message, recipientIDs []uuid.UUID) {
	// Deliver message to all recipients
	if err := ms.deliveryService.DeliverMessage(ctx, message, recipientIDs); err != nil {
		fmt.Printf("Failed to deliver message %s: %v\n", message.ID, err)
	}

	// Send push notifications
	for _, recipientID := range recipientIDs {
		if err := ms.deliveryService.SendPushNotification(ctx, recipientID, message); err != nil {
			fmt.Printf("Failed to send push notification to %s: %v\n", recipientID, err)
		}
	}
}

func (ms *MessageService) truncateContent(content models.MessageContent, maxLength int) string {
	// Extract text content from MessageContent map
	if textContent, exists := content["text"]; exists {
		if textStr, ok := textContent.(string); ok {
			if len(textStr) <= maxLength {
				return textStr
			}
			return textStr[:maxLength] + "..."
		}
	}
	return "[Media content]"
}

func (ms *MessageService) validateSendMessageRequest(req *SendMessageRequest) error {
	if req.DialogID == uuid.Nil {
		return fmt.Errorf("dialog ID is required")
	}

	if req.SenderID == uuid.Nil {
		return fmt.Errorf("sender ID is required")
	}

	if req.Type == "" {
		return fmt.Errorf("message type is required")
	}

	if req.Type == models.MessageTypeText && req.Content == "" {
		return fmt.Errorf("content is required for text messages")
	}

	if req.Type != models.MessageTypeText && req.Type != models.MessageTypeSystem && req.MediaData == nil && req.MediaURL == "" {
		return fmt.Errorf("media data or URL is required for media messages")
	}

	return nil
}

func (ms *MessageService) publishMessageEvent(ctx context.Context, eventType sharedModels.EventType, messageID uuid.UUID, userID uuid.UUID, data map[string]interface{}) error {
	event := &sharedModels.Event{
		ID:            uuid.New(),
		Type:          eventType,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       fmt.Sprintf("Message event: %s", eventType),
		AggregateID:   messageID.String(),
		AggregateType: "message",
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

	return ms.eventPublisher.Publish(ctx, event)
}

// Request/Response structures

type SendMessageRequest struct {
	DialogID      uuid.UUID              `json:"dialog_id" binding:"required"`
	SenderID      uuid.UUID              `json:"sender_id" binding:"required"`
	Type          models.MessageType     `json:"type" binding:"required"`
	Content       string                 `json:"content"`
	MediaData     []byte                 `json:"media_data,omitempty"`
	MediaURL      string                 `json:"media_url,omitempty"`
	ThumbnailURL  string                 `json:"thumbnail_url,omitempty"`
	MediaMetadata map[string]interface{} `json:"media_metadata,omitempty"`
	ReplyToID     *uuid.UUID             `json:"reply_to_id,omitempty"`
	ForwardFromID *uuid.UUID             `json:"forward_from_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type MessageResponse struct {
	ID           uuid.UUID              `json:"id"`
	DialogID     uuid.UUID              `json:"dialog_id"`
	SenderID     uuid.UUID              `json:"sender_id"`
	Type         models.MessageType     `json:"type"`
	Content      string                 `json:"content"`
	MediaURL     string                 `json:"media_url,omitempty"`
	ThumbnailURL string                 `json:"thumbnail_url,omitempty"`
	Status       models.MessageStatus   `json:"status"`
	IsEdited     bool                   `json:"is_edited"`
	IsDeleted    bool                   `json:"is_deleted"`
	ReplyTo      *models.MessageReply   `json:"reply_to,omitempty"`
	ForwardFrom  *models.MessageForward `json:"forward_from,omitempty"`
	ReadBy       []MessageReadInfo      `json:"read_by,omitempty"`
	SentAt       time.Time              `json:"sent_at"`
	EditedAt     *time.Time             `json:"edited_at,omitempty"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`
}

type MessageReadInfo struct {
	UserID uuid.UUID `json:"user_id"`
	ReadAt time.Time `json:"read_at"`
}

type MessageListResponse struct {
	Messages   []*MessageResponse `json:"messages"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

func MessageToResponse(message *models.Message) *MessageResponse {
	// Extract text content from MessageContent
	content := ""
	if textContent, exists := message.Content["text"]; exists {
		if textStr, ok := textContent.(string); ok {
			content = textStr
		}
	}

	// Handle MediaURL pointer
	mediaURL := ""
	if message.MediaURL != nil {
		mediaURL = *message.MediaURL
	}

	// Handle ThumbnailURL pointer
	thumbnailURL := ""
	if message.ThumbnailURL != nil {
		thumbnailURL = *message.ThumbnailURL
	}

	response := &MessageResponse{
		ID:           message.ID,
		DialogID:     message.DialogID,
		SenderID:     message.SenderID,
		Type:         message.Type,
		Content:      content,
		MediaURL:     mediaURL,
		ThumbnailURL: thumbnailURL,
		Status:       message.Status,
		IsEdited:     message.IsEdited,
		IsDeleted:    message.IsDeleted,
		ReplyTo:      message.ReplyTo,
		ForwardFrom:  nil, // TODO: Add ForwardFrom support
		SentAt:       message.SentAt,
		EditedAt:     message.EditedAt,
		DeletedAt:    message.DeletedAt,
	}

	// TODO: Implement read receipts functionality
	// ReadBy field not implemented in Message model yet

	return response
}