package models

import (
	"time"

	"github.com/google/uuid"
)

// ChatMessage represents a chat message in a live stream
// This model is designed for ScyllaDB/Cassandra with a 30-day TTL
// Messages are partitioned by stream_id and ordered by timestamp (descending)
type ChatMessage struct {
	StreamID          uuid.UUID `cql:"stream_id"`
	Timestamp         time.Time `cql:"timestamp"`
	MessageID         uuid.UUID `cql:"message_id"`
	SenderID          uuid.UUID `cql:"sender_id"`
	SenderDisplayName string    `cql:"sender_display_name"`
	MessageText       string    `cql:"message_text"`
	ModerationStatus  string    `cql:"moderation_status"` // 'visible', 'removed', 'flagged'
	MessageType       string    `cql:"message_type"`      // 'text', 'emoji', 'system'
}

// Table and TTL constants for ScyllaDB
const (
	// TableName is the ScyllaDB table name for chat messages
	TableName = "chat_messages"

	// DefaultTTL is the time-to-live in seconds (30 days)
	DefaultTTL = 2592000
)

// ModerationStatus constants
const (
	ModerationStatusVisible = "visible"
	ModerationStatusRemoved = "removed"
	ModerationStatusFlagged = "flagged"
)

// MessageType constants
const (
	MessageTypeText   = "text"
	MessageTypeEmoji  = "emoji"
	MessageTypeSystem = "system"
)