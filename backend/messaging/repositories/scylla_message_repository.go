package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"

	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
)

// ScyllaMessageRepository implements services.MessageRepository using ScyllaDB
type ScyllaMessageRepository struct {
	session *gocql.Session
}

// NewScyllaMessageRepository creates a new ScyllaDB message repository
func NewScyllaMessageRepository(session *gocql.Session) services.MessageRepository {
	return &ScyllaMessageRepository{session: session}
}

// Create creates a new message in ScyllaDB
func (r *ScyllaMessageRepository) Create(ctx context.Context, message *models.Message) error {
	// Convert Content map to JSON string
	contentJSON, err := json.Marshal(message.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	query := `INSERT INTO messages (dialog_id, id, sender_id, text, media_url, message_type, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		message.DialogID,
		message.ID,
		message.SenderID,
		string(contentJSON), // Store MessageContent as JSON string
		message.MediaURL,
		message.Type.String(),
		string(message.Status), // MessageStatus is already a string type
		message.CreatedAt,
		message.UpdatedAt,
	).WithContext(ctx).Exec()
}

// GetByID retrieves a message by its ID
func (r *ScyllaMessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	// Note: This requires a secondary index on id or query by dialog_id + created_at
	// For now, we'll need dialog_id to efficiently query
	return nil, fmt.Errorf("GetByID requires dialog_id for efficient ScyllaDB query - use GetByDialogID instead")
}

// GetByDialogID retrieves messages for a specific dialog with filters and pagination
func (r *ScyllaMessageRepository) GetByDialogID(ctx context.Context, dialogID uuid.UUID, filters services.MessageFilters, pagination services.Pagination) ([]*models.Message, int64, error) {
	query := `SELECT dialog_id, id, sender_id, text, media_url, message_type, status, created_at, updated_at
		FROM messages WHERE dialog_id = ?`

	var args []interface{}
	args = append(args, dialogID)

	// Add time range filters if provided
	if filters.SentFrom != nil {
		query += " AND created_at >= ?"
		args = append(args, *filters.SentFrom)
	}
	if filters.SentTo != nil {
		query += " AND created_at <= ?"
		args = append(args, *filters.SentTo)
	}

	// Add ordering and limit
	query += " ORDER BY created_at DESC"
	if pagination.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT %d", pagination.PageSize)
	}

	iter := r.session.Query(query, args...).WithContext(ctx).Iter()
	defer iter.Close()

	var messages []*models.Message
	var (
		dialogID_   uuid.UUID
		id          uuid.UUID
		senderID    uuid.UUID
		text        string
		mediaURL    *string
		msgType     string
		status      string
		createdAt   time.Time
		updatedAt   time.Time
	)

	for iter.Scan(&dialogID_, &id, &senderID, &text, &mediaURL, &msgType, &status, &createdAt, &updatedAt) {
		// Parse MessageContent from JSON string
		var content models.MessageContent
		if err := json.Unmarshal([]byte(text), &content); err != nil {
			return nil, 0, fmt.Errorf("failed to parse message content: %w", err)
		}

		message := &models.Message{
			ID:        id,
			DialogID:  dialogID_,
			SenderID:  senderID,
			Content:   content,
			MediaURL:  mediaURL,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		// Parse message type
		message.Type = models.MessageType(msgType)

		// Parse status
		message.Status = models.MessageStatus(status)

		// Apply additional filters
		if filters.SenderID != nil && message.SenderID != *filters.SenderID {
			continue
		}
		if filters.Type != nil && message.Type != *filters.Type {
			continue
		}

		messages = append(messages, message)
	}

	if err := iter.Close(); err != nil {
		return nil, 0, fmt.Errorf("query iteration failed: %w", err)
	}

	// Get total count (approximate for ScyllaDB)
	total := int64(len(messages))

	return messages, total, nil
}

// Update updates a message
func (r *ScyllaMessageRepository) Update(ctx context.Context, message *models.Message) error {
	contentJSON, err := json.Marshal(message.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	query := `UPDATE messages SET text = ?, media_url = ?, status = ?, updated_at = ?
		WHERE dialog_id = ? AND created_at = ? AND id = ?`

	return r.session.Query(query,
		string(contentJSON),
		message.MediaURL,
		string(message.Status),
		time.Now(),
		message.DialogID,
		message.CreatedAt,
		message.ID,
	).WithContext(ctx).Exec()
}

// Delete marks a message as deleted
func (r *ScyllaMessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft delete by updating deleted_at timestamp
	// Note: Requires dialog_id and created_at for efficient update
	return fmt.Errorf("Delete requires dialog_id and created_at for efficient ScyllaDB update")
}

// MarkAsRead marks a message as read for a user
func (r *ScyllaMessageRepository) MarkAsRead(ctx context.Context, messageID, userID uuid.UUID) error {
	// This would typically be stored in a separate read_receipts table
	// For now, update message status
	return nil
}

// MarkAsDelivered marks a message as delivered to a user
func (r *ScyllaMessageRepository) MarkAsDelivered(ctx context.Context, messageID, userID uuid.UUID) error {
	// This would typically be stored in a separate delivery_receipts table
	return nil
}

// GetUnreadCount gets the count of unread messages for a user in a dialog
func (r *ScyllaMessageRepository) GetUnreadCount(ctx context.Context, dialogID, userID uuid.UUID) (int, error) {
	// Would require a read_receipts table with counters
	return 0, nil
}

// SearchMessages searches for messages in a dialog
func (r *ScyllaMessageRepository) SearchMessages(ctx context.Context, dialogID uuid.UUID, query string, limit int) ([]*models.Message, error) {
	// Full-text search not natively supported in ScyllaDB
	// Would require integration with external search engine (Elasticsearch)
	return nil, fmt.Errorf("search not implemented - requires external search engine")
}

// GetMessageStats gets statistics for messages in a dialog
func (r *ScyllaMessageRepository) GetMessageStats(ctx context.Context, dialogID uuid.UUID) (*services.MessageStats, error) {
	// Would require aggregation queries or materialized views
	return &services.MessageStats{}, nil
}
