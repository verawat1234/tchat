package models

import (
	"time"

	"github.com/google/uuid"
)

// StreamReaction represents a real-time reaction to a live stream
// Stored in ScyllaDB for high-velocity time-series data with 30-day TTL
type StreamReaction struct {
	StreamID     uuid.UUID `json:"stream_id"`
	Timestamp    time.Time `json:"timestamp"`
	ReactionID   uuid.UUID `json:"reaction_id"`
	ViewerID     uuid.UUID `json:"viewer_id"`
	ReactionType string    `json:"reaction_type"` // Emoji unicode (e.g., '👍', '❤️', '😂')
}

const (
	// StreamReactionTableName is the ScyllaDB table name for stream reactions
	StreamReactionTableName = "stream_reactions"

	// StreamReactionDefaultTTL is the default time-to-live in seconds (30 days)
	StreamReactionDefaultTTL = 2592000
)