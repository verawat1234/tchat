package database

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
	"tchat-backend/shared/database"
)

// MessageTimelineManager handles ScyllaDB operations for message timelines
type MessageTimelineManager struct {
	db *database.ScyllaDB
}

// NewMessageTimelineManager creates a new message timeline manager
func NewMessageTimelineManager(db *database.ScyllaDB) *MessageTimelineManager {
	return &MessageTimelineManager{
		db: db,
	}
}

// RunMigrations creates the necessary tables for message timelines
func (m *MessageTimelineManager) RunMigrations() error {
	migrations := []string{
		// Message timeline table - optimized for time-series queries
		`CREATE TABLE IF NOT EXISTS message_timeline (
			dialog_id UUID,
			bucket_time timestamp,
			message_id UUID,
			sender_id UUID,
			message_type text,
			content text,
			media_url text,
			reply_to_id UUID,
			forward_from_id UUID,
			is_edited boolean,
			is_deleted boolean,
			created_at timestamp,
			PRIMARY KEY (dialog_id, bucket_time, message_id)
		) WITH CLUSTERING ORDER BY (bucket_time DESC, message_id DESC)
		AND compaction = {'class': 'TimeWindowCompactionStrategy', 'compaction_window_unit': 'HOURS', 'compaction_window_size': 24}
		AND gc_grace_seconds = 7776000`,

		// User message timeline - for user's message history across dialogs
		`CREATE TABLE IF NOT EXISTS user_message_timeline (
			user_id UUID,
			bucket_time timestamp,
			dialog_id UUID,
			message_id UUID,
			message_type text,
			content text,
			created_at timestamp,
			PRIMARY KEY (user_id, bucket_time, message_id)
		) WITH CLUSTERING ORDER BY (bucket_time DESC, message_id DESC)
		AND compaction = {'class': 'TimeWindowCompactionStrategy', 'compaction_window_unit': 'HOURS', 'compaction_window_size': 24}
		AND gc_grace_seconds = 7776000`,

		// Dialog statistics - for dialog message counts and activity
		`CREATE TABLE IF NOT EXISTS dialog_stats (
			dialog_id UUID,
			stat_date date,
			message_count counter,
			participant_count counter,
			media_count counter,
			PRIMARY KEY (dialog_id, stat_date)
		) WITH compaction = {'class': 'TimeWindowCompactionStrategy', 'compaction_window_unit': 'DAYS', 'compaction_window_size': 30}`,

		// User activity statistics
		`CREATE TABLE IF NOT EXISTS user_activity_stats (
			user_id UUID,
			stat_date date,
			messages_sent counter,
			dialogs_active counter,
			media_shared counter,
			PRIMARY KEY (user_id, stat_date)
		) WITH compaction = {'class': 'TimeWindowCompactionStrategy', 'compaction_window_unit': 'DAYS', 'compaction_window_size': 30}`,

		// Message search index - for full-text search capabilities
		`CREATE TABLE IF NOT EXISTS message_search_index (
			search_term text,
			dialog_id UUID,
			message_id UUID,
			sender_id UUID,
			content text,
			created_at timestamp,
			PRIMARY KEY (search_term, dialog_id, message_id)
		) WITH CLUSTERING ORDER BY (dialog_id ASC, message_id DESC)
		AND compaction = {'class': 'SizeTieredCompactionStrategy'}`,

		// Media timeline - for media messages
		`CREATE TABLE IF NOT EXISTS media_timeline (
			dialog_id UUID,
			media_type text,
			bucket_time timestamp,
			message_id UUID,
			sender_id UUID,
			media_url text,
			media_metadata text,
			created_at timestamp,
			PRIMARY KEY ((dialog_id, media_type), bucket_time, message_id)
		) WITH CLUSTERING ORDER BY (bucket_time DESC, message_id DESC)
		AND compaction = {'class': 'TimeWindowCompactionStrategy', 'compaction_window_unit': 'HOURS', 'compaction_window_size': 24}
		AND gc_grace_seconds = 7776000`,

		// Reaction timeline - for message reactions
		`CREATE TABLE IF NOT EXISTS reaction_timeline (
			message_id UUID,
			user_id UUID,
			reaction_type text,
			created_at timestamp,
			PRIMARY KEY (message_id, user_id, reaction_type)
		) WITH compaction = {'class': 'SizeTieredCompactionStrategy'}`,

		// Read receipts - for message read status
		`CREATE TABLE IF NOT EXISTS message_read_receipts (
			message_id UUID,
			user_id UUID,
			read_at timestamp,
			PRIMARY KEY (message_id, user_id)
		) WITH compaction = {'class': 'SizeTieredCompactionStrategy'}
		AND default_time_to_live = 2592000`,

		// Typing indicators - for real-time typing status
		`CREATE TABLE IF NOT EXISTS typing_indicators (
			dialog_id UUID,
			user_id UUID,
			is_typing boolean,
			updated_at timestamp,
			PRIMARY KEY (dialog_id, user_id)
		) WITH compaction = {'class': 'SizeTieredCompactionStrategy'}
		AND default_time_to_live = 300`,

		// Message delivery status
		`CREATE TABLE IF NOT EXISTS message_delivery_status (
			message_id UUID,
			user_id UUID,
			status text,
			delivered_at timestamp,
			PRIMARY KEY (message_id, user_id)
		) WITH compaction = {'class': 'SizeTieredCompactionStrategy'}
		AND default_time_to_live = 2592000`,
	}

	return m.db.RunMigrations(migrations)
}

// InsertMessage inserts a message into the timeline
func (m *MessageTimelineManager) InsertMessage(
	dialogID, messageID, senderID gocql.UUID,
	messageType, content, mediaURL string,
	replyToID, forwardFromID *gocql.UUID,
	isEdited, isDeleted bool,
	createdAt, bucketTime int64,
) error {
	batch := m.db.CreateBatch(gocql.LoggedBatch)

	// Insert into message timeline
	batch.Query(`INSERT INTO message_timeline
		(dialog_id, bucket_time, message_id, sender_id, message_type, content, media_url,
		 reply_to_id, forward_from_id, is_edited, is_deleted, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dialogID, bucketTime, messageID, senderID, messageType, content, mediaURL,
		replyToID, forwardFromID, isEdited, isDeleted, createdAt)

	// Insert into user timeline
	batch.Query(`INSERT INTO user_message_timeline
		(user_id, bucket_time, dialog_id, message_id, message_type, content, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		senderID, bucketTime, dialogID, messageID, messageType, content, createdAt)

	// Update dialog stats
	batch.Query(`UPDATE dialog_stats SET message_count = message_count + 1
		WHERE dialog_id = ? AND stat_date = ?`,
		dialogID, bucketTime)

	// Update user activity stats
	batch.Query(`UPDATE user_activity_stats SET messages_sent = messages_sent + 1
		WHERE user_id = ? AND stat_date = ?`,
		senderID, bucketTime)

	// If media message, insert into media timeline
	if mediaURL != "" {
		batch.Query(`INSERT INTO media_timeline
			(dialog_id, media_type, bucket_time, message_id, sender_id, media_url, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			dialogID, messageType, bucketTime, messageID, senderID, mediaURL, createdAt)

		// Update media count
		batch.Query(`UPDATE dialog_stats SET media_count = media_count + 1
			WHERE dialog_id = ? AND stat_date = ?`,
			dialogID, bucketTime)

		batch.Query(`UPDATE user_activity_stats SET media_shared = media_shared + 1
			WHERE user_id = ? AND stat_date = ?`,
			senderID, bucketTime)
	}

	return m.db.ExecuteBatch(batch)
}

// GetDialogMessages retrieves messages for a dialog with pagination
func (m *MessageTimelineManager) GetDialogMessages(
	dialogID gocql.UUID,
	limit int,
	lastBucketTime *int64,
	lastMessageID *gocql.UUID,
) (*gocql.Iter, error) {
	var query string
	var values []interface{}

	if lastBucketTime != nil && lastMessageID != nil {
		query = `SELECT bucket_time, message_id, sender_id, message_type, content, media_url,
			reply_to_id, forward_from_id, is_edited, is_deleted, created_at
			FROM message_timeline
			WHERE dialog_id = ? AND (bucket_time, message_id) < (?, ?)
			ORDER BY bucket_time DESC, message_id DESC
			LIMIT ?`
		values = []interface{}{dialogID, *lastBucketTime, *lastMessageID, limit}
	} else {
		query = `SELECT bucket_time, message_id, sender_id, message_type, content, media_url,
			reply_to_id, forward_from_id, is_edited, is_deleted, created_at
			FROM message_timeline
			WHERE dialog_id = ?
			ORDER BY bucket_time DESC, message_id DESC
			LIMIT ?`
		values = []interface{}{dialogID, limit}
	}

	return m.db.QueryRows(query, values...), nil
}

// GetUserMessages retrieves messages sent by a user
func (m *MessageTimelineManager) GetUserMessages(
	userID gocql.UUID,
	limit int,
	lastBucketTime *int64,
	lastMessageID *gocql.UUID,
) (*gocql.Iter, error) {
	var query string
	var values []interface{}

	if lastBucketTime != nil && lastMessageID != nil {
		query = `SELECT bucket_time, dialog_id, message_id, message_type, content, created_at
			FROM user_message_timeline
			WHERE user_id = ? AND (bucket_time, message_id) < (?, ?)
			ORDER BY bucket_time DESC, message_id DESC
			LIMIT ?`
		values = []interface{}{userID, *lastBucketTime, *lastMessageID, limit}
	} else {
		query = `SELECT bucket_time, dialog_id, message_id, message_type, content, created_at
			FROM user_message_timeline
			WHERE user_id = ?
			ORDER BY bucket_time DESC, message_id DESC
			LIMIT ?`
		values = []interface{}{userID, limit}
	}

	return m.db.QueryRows(query, values...), nil
}

// AddReaction adds a reaction to a message
func (m *MessageTimelineManager) AddReaction(
	messageID, userID gocql.UUID,
	reactionType string,
	createdAt int64,
) error {
	return m.db.ExecuteQuery(`INSERT INTO reaction_timeline
		(message_id, user_id, reaction_type, created_at)
		VALUES (?, ?, ?, ?)`,
		messageID, userID, reactionType, createdAt)
}

// RemoveReaction removes a reaction from a message
func (m *MessageTimelineManager) RemoveReaction(
	messageID, userID gocql.UUID,
	reactionType string,
) error {
	return m.db.ExecuteQuery(`DELETE FROM reaction_timeline
		WHERE message_id = ? AND user_id = ? AND reaction_type = ?`,
		messageID, userID, reactionType)
}

// MarkMessageAsRead marks a message as read by a user
func (m *MessageTimelineManager) MarkMessageAsRead(
	messageID, userID gocql.UUID,
	readAt int64,
) error {
	return m.db.ExecuteQuery(`INSERT INTO message_read_receipts
		(message_id, user_id, read_at)
		VALUES (?, ?, ?)`,
		messageID, userID, readAt)
}

// UpdateTypingIndicator updates typing status for a user in a dialog
func (m *MessageTimelineManager) UpdateTypingIndicator(
	dialogID, userID gocql.UUID,
	isTyping bool,
	updatedAt int64,
) error {
	return m.db.ExecuteQuery(`INSERT INTO typing_indicators
		(dialog_id, user_id, is_typing, updated_at)
		VALUES (?, ?, ?, ?)`,
		dialogID, userID, isTyping, updatedAt)
}

// GetDialogStats retrieves statistics for a dialog
func (m *MessageTimelineManager) GetDialogStats(
	dialogID gocql.UUID,
	statDate int64,
) (messageCount, participantCount, mediaCount int64, err error) {
	err = m.db.QueryRow(
		[]interface{}{&messageCount, &participantCount, &mediaCount},
		`SELECT message_count, participant_count, media_count
		FROM dialog_stats
		WHERE dialog_id = ? AND stat_date = ?`,
		dialogID, statDate,
	)
	return
}

// GetMediaTimeline retrieves media messages for a dialog
func (m *MessageTimelineManager) GetMediaTimeline(
	dialogID gocql.UUID,
	mediaType string,
	limit int,
	lastBucketTime *int64,
	lastMessageID *gocql.UUID,
) (*gocql.Iter, error) {
	var query string
	var values []interface{}

	if lastBucketTime != nil && lastMessageID != nil {
		query = `SELECT bucket_time, message_id, sender_id, media_url, media_metadata, created_at
			FROM media_timeline
			WHERE dialog_id = ? AND media_type = ? AND (bucket_time, message_id) < (?, ?)
			ORDER BY bucket_time DESC, message_id DESC
			LIMIT ?`
		values = []interface{}{dialogID, mediaType, *lastBucketTime, *lastMessageID, limit}
	} else {
		query = `SELECT bucket_time, message_id, sender_id, media_url, media_metadata, created_at
			FROM media_timeline
			WHERE dialog_id = ? AND media_type = ?
			ORDER BY bucket_time DESC, message_id DESC
			LIMIT ?`
		values = []interface{}{dialogID, mediaType, limit}
	}

	return m.db.QueryRows(query, values...), nil
}

// GetMessageReactions retrieves reactions for a message
func (m *MessageTimelineManager) GetMessageReactions(messageID gocql.UUID) (*gocql.Iter, error) {
	return m.db.QueryRows(`SELECT user_id, reaction_type, created_at
		FROM reaction_timeline
		WHERE message_id = ?`,
		messageID), nil
}

// Helper function to calculate bucket time (hour-based bucketing)
func CalculateBucketTime(timestamp int64) int64 {
	// Round down to the nearest hour
	return (timestamp / 3600) * 3600
}

// Helper function to calculate date bucket for stats
func CalculateDateBucket(timestamp int64) int64 {
	// Round down to the nearest day
	return (timestamp / 86400) * 86400
}

// Cleanup old data based on retention policy
func (m *MessageTimelineManager) CleanupOldData(retentionDays int) error {
	cutoffTime := fmt.Sprintf("%d", (time.Now().Unix()-int64(retentionDays*86400))*1000)

	queries := []string{
		fmt.Sprintf("DELETE FROM message_timeline WHERE bucket_time < %s", cutoffTime),
		fmt.Sprintf("DELETE FROM user_message_timeline WHERE bucket_time < %s", cutoffTime),
		fmt.Sprintf("DELETE FROM media_timeline WHERE bucket_time < %s", cutoffTime),
	}

	for _, query := range queries {
		if err := m.db.ExecuteQuery(query); err != nil {
			log.Printf("Warning: Failed to cleanup old data with query '%s': %v", query, err)
		}
	}

	return nil
}