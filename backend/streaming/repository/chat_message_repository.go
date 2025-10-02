package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
)

// Helper function to convert google/uuid.UUID to gocql.UUID
func toGocqlUUID(googleUUID uuid.UUID) gocql.UUID {
	var gocqlUUID gocql.UUID
	copy(gocqlUUID[:], googleUUID[:])
	return gocqlUUID
}

// Helper function to convert gocql.UUID to google/uuid.UUID
func fromGocqlUUID(gocqlUUID gocql.UUID) uuid.UUID {
	var googleUUID uuid.UUID
	copy(googleUUID[:], gocqlUUID[:])
	return googleUUID
}

// ChatMessageRepository defines the interface for chat message data access
type ChatMessageRepository interface {
	// Create inserts a new chat message with 30-day TTL
	Create(ctx context.Context, message *models.ChatMessage) error

	// GetByID retrieves a specific chat message by stream and message ID
	GetByID(ctx context.Context, streamID, messageID uuid.UUID) (*models.ChatMessage, error)

	// ListByStream retrieves chat messages for a stream with pagination
	// Uses timestamp for cursor-based pagination (descending order)
	ListByStream(ctx context.Context, streamID uuid.UUID, limit int, beforeTimestamp *time.Time) ([]*models.ChatMessage, error)

	// Update modifies a chat message (primarily for moderation status)
	Update(ctx context.Context, streamID, messageID uuid.UUID, updates map[string]interface{}) error

	// Delete removes a chat message (soft delete via moderation_status)
	Delete(ctx context.Context, streamID, messageID uuid.UUID) error

	// GetByModeration retrieves messages by moderation status for moderator review
	GetByModeration(ctx context.Context, streamID uuid.UUID, status string, limit int) ([]*models.ChatMessage, error)

	// CountByStream returns the total message count for a stream
	CountByStream(ctx context.Context, streamID uuid.UUID) (int64, error)

	// BatchCreate inserts multiple messages efficiently
	BatchCreate(ctx context.Context, messages []*models.ChatMessage) error
}

// chatMessageRepository implements ChatMessageRepository using gocql
type chatMessageRepository struct {
	session *gocql.Session

	// Prepared statements for performance
	insertStmt          *gocql.Query
	selectByIDStmt      *gocql.Query
	selectByStreamStmt  *gocql.Query
	updateStmt          *gocql.Query
	selectByModerationStmt *gocql.Query
}

// NewChatMessageRepository creates a new chat message repository
func NewChatMessageRepository(session *gocql.Session) (ChatMessageRepository, error) {
	if session == nil {
		return nil, fmt.Errorf("gocql session cannot be nil")
	}

	repo := &chatMessageRepository{
		session: session,
	}

	// Initialize prepared statements for performance
	if err := repo.prepareCQLStatements(); err != nil {
		return nil, fmt.Errorf("failed to prepare CQL statements: %w", err)
	}

	return repo, nil
}

// prepareCQLStatements prepares all CQL statements for reuse
func (r *chatMessageRepository) prepareCQLStatements() error {
	// Insert statement with 30-day TTL (2592000 seconds)
	insertCQL := `
		INSERT INTO chat_messages (
			stream_id, timestamp, message_id, sender_id,
			sender_display_name, message_text, moderation_status, message_type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		USING TTL ?
	`
	r.insertStmt = r.session.Query(insertCQL)

	// Select by ID statement
	// Note: Uses ALLOW FILTERING since message_id is not the first clustering column
	// This is acceptable for GetByID as it's used sparingly in production
	selectByIDCQL := `
		SELECT stream_id, timestamp, message_id, sender_id,
			   sender_display_name, message_text, moderation_status, message_type
		FROM chat_messages
		WHERE stream_id = ? AND message_id = ?
		LIMIT 1
		ALLOW FILTERING
	`
	r.selectByIDStmt = r.session.Query(selectByIDCQL)

	// Select by stream with pagination
	selectByStreamCQL := `
		SELECT stream_id, timestamp, message_id, sender_id,
			   sender_display_name, message_text, moderation_status, message_type
		FROM chat_messages
		WHERE stream_id = ? AND timestamp < ?
		ORDER BY timestamp DESC
		LIMIT ?
	`
	r.selectByStreamStmt = r.session.Query(selectByStreamCQL)

	// Update statement (for moderation)
	updateCQL := `
		UPDATE chat_messages
		SET moderation_status = ?
		WHERE stream_id = ? AND timestamp = ? AND message_id = ?
	`
	r.updateStmt = r.session.Query(updateCQL)

	// Select by moderation status
	selectByModerationCQL := `
		SELECT stream_id, timestamp, message_id, sender_id,
			   sender_display_name, message_text, moderation_status, message_type
		FROM chat_messages
		WHERE stream_id = ? AND moderation_status = ?
		LIMIT ?
		ALLOW FILTERING
	`
	r.selectByModerationStmt = r.session.Query(selectByModerationCQL)

	return nil
}

// Create inserts a new chat message with 30-day TTL
func (r *chatMessageRepository) Create(ctx context.Context, message *models.ChatMessage) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	// Set default timestamp if not provided
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	// Generate message ID if not provided
	if message.MessageID == uuid.Nil {
		message.MessageID = uuid.New()
	}

	// Set default moderation status if not provided
	if message.ModerationStatus == "" {
		message.ModerationStatus = models.ModerationStatusVisible
	}

	// Set default message type if not provided
	if message.MessageType == "" {
		message.MessageType = models.MessageTypeText
	}

	// Execute insert with TTL (convert google UUID to gocql UUID)
	query := r.session.Query(r.insertStmt.Statement(),
		toGocqlUUID(message.StreamID),
		message.Timestamp,
		toGocqlUUID(message.MessageID),
		toGocqlUUID(message.SenderID),
		message.SenderDisplayName,
		message.MessageText,
		message.ModerationStatus,
		message.MessageType,
		models.DefaultTTL, // 30 days in seconds
	).WithContext(ctx)

	if err := query.Exec(); err != nil {
		return fmt.Errorf("failed to create chat message: %w", err)
	}

	return nil
}

// GetByID retrieves a specific chat message
func (r *chatMessageRepository) GetByID(ctx context.Context, streamID, messageID uuid.UUID) (*models.ChatMessage, error) {
	if streamID == uuid.Nil || messageID == uuid.Nil {
		return nil, fmt.Errorf("streamID and messageID cannot be nil")
	}

	var message models.ChatMessage

	query := r.session.Query(r.selectByIDStmt.Statement(),
		toGocqlUUID(streamID),
		toGocqlUUID(messageID),
	).WithContext(ctx)

	var gocqlStreamID, gocqlMessageID, gocqlSenderID gocql.UUID
	err := query.Scan(
		&gocqlStreamID,
		&message.Timestamp,
		&gocqlMessageID,
		&gocqlSenderID,
		&message.SenderDisplayName,
		&message.MessageText,
		&message.ModerationStatus,
		&message.MessageType,
	)

	// Convert gocql UUIDs back to google UUIDs
	message.StreamID = fromGocqlUUID(gocqlStreamID)
	message.MessageID = fromGocqlUUID(gocqlMessageID)
	message.SenderID = fromGocqlUUID(gocqlSenderID)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, fmt.Errorf("chat message not found")
		}
		return nil, fmt.Errorf("failed to retrieve chat message: %w", err)
	}

	return &message, nil
}

// ListByStream retrieves chat messages for a stream with cursor-based pagination
func (r *chatMessageRepository) ListByStream(ctx context.Context, streamID uuid.UUID, limit int, beforeTimestamp *time.Time) ([]*models.ChatMessage, error) {
	if streamID == uuid.Nil {
		return nil, fmt.Errorf("streamID cannot be nil")
	}

	if limit <= 0 || limit > 500 {
		limit = 100 // Default limit with reasonable upper bound
	}

	// Use current time as cursor if not provided
	cursor := time.Now()
	if beforeTimestamp != nil {
		cursor = *beforeTimestamp
	}

	query := r.session.Query(r.selectByStreamStmt.Statement(),
		toGocqlUUID(streamID),
		cursor,
		limit,
	).WithContext(ctx)

	iter := query.Iter()
	defer iter.Close()

	messages := make([]*models.ChatMessage, 0, limit)
	var gocqlStreamID, gocqlMessageID, gocqlSenderID gocql.UUID
	var message models.ChatMessage

	for iter.Scan(
		&gocqlStreamID,
		&message.Timestamp,
		&gocqlMessageID,
		&gocqlSenderID,
		&message.SenderDisplayName,
		&message.MessageText,
		&message.ModerationStatus,
		&message.MessageType,
	) {
		// Convert gocql UUIDs back to google UUIDs
		message.StreamID = fromGocqlUUID(gocqlStreamID)
		message.MessageID = fromGocqlUUID(gocqlMessageID)
		message.SenderID = fromGocqlUUID(gocqlSenderID)

		// Create a copy to avoid pointer reuse issues
		msg := message
		messages = append(messages, &msg)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to iterate chat messages: %w", err)
	}

	return messages, nil
}

// Update modifies a chat message (primarily for moderation)
func (r *chatMessageRepository) Update(ctx context.Context, streamID, messageID uuid.UUID, updates map[string]interface{}) error {
	if streamID == uuid.Nil || messageID == uuid.Nil {
		return fmt.Errorf("streamID and messageID cannot be nil")
	}

	if len(updates) == 0 {
		return fmt.Errorf("updates cannot be empty")
	}

	// First, retrieve the message to get its timestamp (required for update)
	message, err := r.GetByID(ctx, streamID, messageID)
	if err != nil {
		return fmt.Errorf("failed to retrieve message for update: %w", err)
	}

	// Only moderation_status updates are supported for ScyllaDB efficiency
	moderationStatus, ok := updates["moderation_status"].(string)
	if !ok {
		return fmt.Errorf("only moderation_status updates are supported")
	}

	// Validate moderation status
	if moderationStatus != models.ModerationStatusVisible &&
		moderationStatus != models.ModerationStatusRemoved &&
		moderationStatus != models.ModerationStatusFlagged {
		return fmt.Errorf("invalid moderation_status: %s", moderationStatus)
	}

	query := r.session.Query(r.updateStmt.Statement(),
		moderationStatus,
		toGocqlUUID(streamID),
		message.Timestamp,
		toGocqlUUID(messageID),
	).WithContext(ctx)

	if err := query.Exec(); err != nil {
		return fmt.Errorf("failed to update chat message: %w", err)
	}

	return nil
}

// Delete soft deletes a chat message by setting moderation status to 'removed'
func (r *chatMessageRepository) Delete(ctx context.Context, streamID, messageID uuid.UUID) error {
	return r.Update(ctx, streamID, messageID, map[string]interface{}{
		"moderation_status": models.ModerationStatusRemoved,
	})
}

// GetByModeration retrieves messages by moderation status
// Note: Uses ALLOW FILTERING - should be used sparingly in production
func (r *chatMessageRepository) GetByModeration(ctx context.Context, streamID uuid.UUID, status string, limit int) ([]*models.ChatMessage, error) {
	if streamID == uuid.Nil {
		return nil, fmt.Errorf("streamID cannot be nil")
	}

	// Validate moderation status
	if status != models.ModerationStatusVisible &&
		status != models.ModerationStatusRemoved &&
		status != models.ModerationStatusFlagged {
		return nil, fmt.Errorf("invalid moderation_status: %s", status)
	}

	if limit <= 0 || limit > 500 {
		limit = 100
	}

	query := r.session.Query(r.selectByModerationStmt.Statement(),
		toGocqlUUID(streamID),
		status,
		limit,
	).WithContext(ctx)

	iter := query.Iter()
	defer iter.Close()

	messages := make([]*models.ChatMessage, 0, limit)
	var gocqlStreamID, gocqlMessageID, gocqlSenderID gocql.UUID
	var message models.ChatMessage

	for iter.Scan(
		&gocqlStreamID,
		&message.Timestamp,
		&gocqlMessageID,
		&gocqlSenderID,
		&message.SenderDisplayName,
		&message.MessageText,
		&message.ModerationStatus,
		&message.MessageType,
	) {
		// Convert gocql UUIDs back to google UUIDs
		message.StreamID = fromGocqlUUID(gocqlStreamID)
		message.MessageID = fromGocqlUUID(gocqlMessageID)
		message.SenderID = fromGocqlUUID(gocqlSenderID)

		msg := message
		messages = append(messages, &msg)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to iterate moderated messages: %w", err)
	}

	return messages, nil
}

// CountByStream returns the total message count for a stream
// Note: COUNT queries can be expensive in ScyllaDB - use with caution
func (r *chatMessageRepository) CountByStream(ctx context.Context, streamID uuid.UUID) (int64, error) {
	if streamID == uuid.Nil {
		return 0, fmt.Errorf("streamID cannot be nil")
	}

	countCQL := `SELECT COUNT(*) FROM chat_messages WHERE stream_id = ?`
	query := r.session.Query(countCQL, toGocqlUUID(streamID)).WithContext(ctx)

	var count int64
	if err := query.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

// BatchCreate inserts multiple messages efficiently using batch operations
func (r *chatMessageRepository) BatchCreate(ctx context.Context, messages []*models.ChatMessage) error {
	if len(messages) == 0 {
		return fmt.Errorf("messages cannot be empty")
	}

	// ScyllaDB recommends batch sizes of 10-100 statements
	const maxBatchSize = 50
	for i := 0; i < len(messages); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(messages) {
			end = len(messages)
		}

		batch := r.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)

		for _, message := range messages[i:end] {
			// Set defaults
			if message.Timestamp.IsZero() {
				message.Timestamp = time.Now()
			}
			if message.MessageID == uuid.Nil {
				message.MessageID = uuid.New()
			}
			if message.ModerationStatus == "" {
				message.ModerationStatus = models.ModerationStatusVisible
			}
			if message.MessageType == "" {
				message.MessageType = models.MessageTypeText
			}

			batch.Query(r.insertStmt.Statement(),
				toGocqlUUID(message.StreamID),
				message.Timestamp,
				toGocqlUUID(message.MessageID),
				toGocqlUUID(message.SenderID),
				message.SenderDisplayName,
				message.MessageText,
				message.ModerationStatus,
				message.MessageType,
				models.DefaultTTL,
			)
		}

		if err := r.session.ExecuteBatch(batch); err != nil {
			return fmt.Errorf("failed to batch create messages (batch %d-%d): %w", i, end, err)
		}
	}

	return nil
}

// RateLimitCheck validates message rate (5 messages/second per user)
// This should be called by the service layer before Create
func RateLimitCheck(ctx context.Context, session *gocql.Session, streamID, senderID uuid.UUID, windowSeconds int) (int64, error) {
	if streamID == uuid.Nil || senderID == uuid.Nil {
		return 0, fmt.Errorf("streamID and senderID cannot be nil")
	}

	windowStart := time.Now().Add(-time.Duration(windowSeconds) * time.Second)

	countCQL := `
		SELECT COUNT(*)
		FROM chat_messages
		WHERE stream_id = ? AND sender_id = ? AND timestamp >= ?
		ALLOW FILTERING
	`
	query := session.Query(countCQL, streamID, senderID, windowStart).WithContext(ctx)

	var count int64
	if err := query.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to check rate limit: %w", err)
	}

	return count, nil
}