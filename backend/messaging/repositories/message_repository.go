package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
)

// MessageRepository implements services.MessageRepository using GORM
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) services.MessageRepository {
	return &MessageRepository{db: db}
}

// Create creates a new message in the database
func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetByID retrieves a message by its ID
func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := r.db.WithContext(ctx).First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetByDialogID retrieves messages for a specific dialog with filters and pagination
func (r *MessageRepository) GetByDialogID(ctx context.Context, dialogID uuid.UUID, filters services.MessageFilters, pagination services.Pagination) ([]*models.Message, int64, error) {
	var messages []*models.Message
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Message{}).Where("dialog_id = ?", dialogID)

	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.SenderID != nil {
		query = query.Where("sender_id = ?", *filters.SenderID)
	}
	if filters.IsEdited != nil {
		query = query.Where("is_edited = ?", *filters.IsEdited)
	}
	if filters.IsDeleted != nil {
		query = query.Where("is_deleted = ?", *filters.IsDeleted)
	}
	if filters.SentFrom != nil {
		query = query.Where("sent_at >= ?", *filters.SentFrom)
	}
	if filters.SentTo != nil {
		query = query.Where("sent_at <= ?", *filters.SentTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	orderBy := "sent_at"
	if pagination.OrderBy != "" {
		orderBy = pagination.OrderBy
	}
	order := "DESC"
	if pagination.Order != "" {
		order = pagination.Order
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.Order(fmt.Sprintf("%s %s", orderBy, order)).
		Limit(pagination.PageSize).
		Offset(offset).
		Find(&messages).Error

	return messages, total, err
}

// Update updates an existing message
func (r *MessageRepository) Update(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Save(message).Error
}

// Delete soft deletes a message
func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Message{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// MarkAsRead marks a message as read for a specific user
func (r *MessageRepository) MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error {
	// This would typically update a separate read_receipts table
	// For now, we'll implement a simple approach
	return r.db.WithContext(ctx).Exec(
		"UPDATE messages SET read_by = array_append(COALESCE(read_by, '{}'), ?) WHERE id = ? AND NOT (? = ANY(COALESCE(read_by, '{}')))",
		userID.String(), messageID, userID.String(),
	).Error
}

// MarkAsDelivered marks a message as delivered for a specific user
func (r *MessageRepository) MarkAsDelivered(ctx context.Context, messageID, userID uuid.UUID) error {
	// This would typically update a separate delivery_receipts table
	// For now, we'll implement a simple approach
	return r.db.WithContext(ctx).Exec(
		"UPDATE messages SET delivered_to = array_append(COALESCE(delivered_to, '{}'), ?) WHERE id = ? AND NOT (? = ANY(COALESCE(delivered_to, '{}')))",
		userID.String(), messageID, userID.String(),
	).Error
}

// GetUnreadCount returns the number of unread messages in a dialog for a user
func (r *MessageRepository) GetUnreadCount(ctx context.Context, dialogID, userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Message{}).
		Where("dialog_id = ? AND sender_id != ? AND NOT (? = ANY(COALESCE(read_by, '{}')))",
			dialogID, userID, userID.String()).
		Count(&count).Error
	return int(count), err
}

// SearchMessages searches for messages containing a query string within a dialog
func (r *MessageRepository) SearchMessages(ctx context.Context, dialogID uuid.UUID, query string, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	err := r.db.WithContext(ctx).
		Where("dialog_id = ? AND content ILIKE ?", dialogID, "%"+query+"%").
		Order("sent_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// GetMessageStats returns statistics about messages in a dialog
func (r *MessageRepository) GetMessageStats(ctx context.Context, dialogID uuid.UUID) (*services.MessageStats, error) {
	var stats services.MessageStats

	// Total messages
	err := r.db.WithContext(ctx).Model(&models.Message{}).
		Where("dialog_id = ?", dialogID).
		Count(&stats.TotalMessages).Error
	if err != nil {
		return nil, err
	}

	// Messages by type
	var textCount, mediaCount, systemCount int64
	r.db.WithContext(ctx).Model(&models.Message{}).
		Where("dialog_id = ? AND type = ?", dialogID, models.MessageTypeText).
		Count(&textCount)
	r.db.WithContext(ctx).Model(&models.Message{}).
		Where("dialog_id = ? AND type IN ?", dialogID, []models.MessageType{
			models.MessageTypeImage, models.MessageTypeVideo, models.MessageTypeVoice, models.MessageTypeFile,
		}).Count(&mediaCount)
	r.db.WithContext(ctx).Model(&models.Message{}).
		Where("dialog_id = ? AND type = ?", dialogID, models.MessageTypeSystem).
		Count(&systemCount)

	stats.TextMessages = textCount
	stats.MediaMessages = mediaCount
	stats.SystemMessages = systemCount

	// Average message length for text messages
	var avgLength float64
	r.db.WithContext(ctx).Model(&models.Message{}).
		Where("dialog_id = ? AND type = ?", dialogID, models.MessageTypeText).
		Select("AVG(LENGTH(content))").
		Row().Scan(&avgLength)
	stats.AverageLength = avgLength

	// Active senders (unique senders in the last 30 days)
	var activeSenders int64
	r.db.WithContext(ctx).Model(&models.Message{}).
		Where("dialog_id = ? AND sent_at > NOW() - INTERVAL '30 days'", dialogID).
		Select("COUNT(DISTINCT sender_id)").
		Row().Scan(&activeSenders)
	stats.ActiveSenders = activeSenders

	// Messages per day (average over last 30 days)
	var messagesPerDay float64
	r.db.WithContext(ctx).Model(&models.Message{}).
		Where("dialog_id = ? AND sent_at > NOW() - INTERVAL '30 days'", dialogID).
		Select("COUNT(*) / 30.0").
		Row().Scan(&messagesPerDay)
	stats.MessagesPerDay = int64(messagesPerDay)

	return &stats, nil
}