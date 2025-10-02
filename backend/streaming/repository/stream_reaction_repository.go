package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"tchat.dev/streaming/models"
)

// StreamReactionRepository defines the interface for stream reaction data access
type StreamReactionRepository interface {
	// Create inserts a new reaction into ScyllaDB with 30-day TTL
	Create(ctx context.Context, reaction *models.StreamReaction) error

	// GetByID retrieves a specific reaction by stream and reaction IDs
	GetByID(ctx context.Context, streamID, reactionID uuid.UUID) (*models.StreamReaction, error)

	// ListByStream retrieves reactions for a stream with pagination
	// beforeTimestamp is optional for pagination (reactions before this time)
	ListByStream(ctx context.Context, streamID uuid.UUID, limit int, beforeTimestamp *time.Time) ([]*models.StreamReaction, error)

	// GetReactionCounts retrieves aggregated reaction counts from Redis
	// Returns a map of reaction_type to count
	GetReactionCounts(ctx context.Context, streamID uuid.UUID) (map[string]int, error)

	// Delete removes a reaction from ScyllaDB
	Delete(ctx context.Context, streamID, reactionID uuid.UUID) error

	// IncrementReactionCount increments the Redis counter for a reaction type
	IncrementReactionCount(ctx context.Context, streamID uuid.UUID, reactionType string) error

	// DecrementReactionCount decrements the Redis counter for a reaction type
	DecrementReactionCount(ctx context.Context, streamID uuid.UUID, reactionType string) error

	// SetReactionTTL sets TTL on Redis counters matching stream duration
	SetReactionTTL(ctx context.Context, streamID uuid.UUID, ttl time.Duration) error
}

// streamReactionRepository implements StreamReactionRepository using ScyllaDB and Redis
type streamReactionRepository struct {
	scyllaSession *gocql.Session
	redisClient   *redis.Client
}

// NewStreamReactionRepository creates a new stream reaction repository
func NewStreamReactionRepository(scyllaSession *gocql.Session, redisClient *redis.Client) StreamReactionRepository {
	return &streamReactionRepository{
		scyllaSession: scyllaSession,
		redisClient:   redisClient,
	}
}

// Create inserts a new reaction into ScyllaDB with 30-day TTL
func (r *streamReactionRepository) Create(ctx context.Context, reaction *models.StreamReaction) error {
	// Validate rate limiting: 10 reactions/second per user
	rateLimitKey := fmt.Sprintf("reaction_rate:%s:%s", reaction.StreamID, reaction.ViewerID)
	count, err := r.redisClient.Incr(ctx, rateLimitKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Set expiry on first increment
	if count == 1 {
		r.redisClient.Expire(ctx, rateLimitKey, time.Second)
	}

	// Check rate limit: 10 reactions/second
	if count > 10 {
		return fmt.Errorf("rate limit exceeded: maximum 10 reactions per second")
	}

	// Prepare CQL statement with TTL
	query := `INSERT INTO stream_reactions (stream_id, timestamp, reaction_id, viewer_id, reaction_type)
	          VALUES (?, ?, ?, ?, ?)
	          USING TTL ?`

	// Execute with 30-day TTL
	if err := r.scyllaSession.Query(query,
		reaction.StreamID,
		reaction.Timestamp,
		reaction.ReactionID,
		reaction.ViewerID,
		reaction.ReactionType,
		models.StreamReactionDefaultTTL,
	).WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to create reaction: %w", err)
	}

	// Increment Redis counter for real-time aggregation
	if err := r.IncrementReactionCount(ctx, reaction.StreamID, reaction.ReactionType); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to increment reaction count: %v\n", err)
	}

	return nil
}

// GetByID retrieves a specific reaction by stream and reaction IDs
func (r *streamReactionRepository) GetByID(ctx context.Context, streamID, reactionID uuid.UUID) (*models.StreamReaction, error) {
	query := `SELECT stream_id, timestamp, reaction_id, viewer_id, reaction_type
	          FROM stream_reactions
	          WHERE stream_id = ? AND reaction_id = ?`

	var reaction models.StreamReaction
	if err := r.scyllaSession.Query(query, streamID, reactionID).
		WithContext(ctx).
		Scan(
			&reaction.StreamID,
			&reaction.Timestamp,
			&reaction.ReactionID,
			&reaction.ViewerID,
			&reaction.ReactionType,
		); err != nil {
		if err == gocql.ErrNotFound {
			return nil, fmt.Errorf("reaction not found")
		}
		return nil, fmt.Errorf("failed to get reaction: %w", err)
	}

	return &reaction, nil
}

// ListByStream retrieves reactions for a stream with pagination
func (r *streamReactionRepository) ListByStream(ctx context.Context, streamID uuid.UUID, limit int, beforeTimestamp *time.Time) ([]*models.StreamReaction, error) {
	var query string
	var args []interface{}

	if beforeTimestamp != nil {
		// Pagination: get reactions before the specified timestamp
		query = `SELECT stream_id, timestamp, reaction_id, viewer_id, reaction_type
		         FROM stream_reactions
		         WHERE stream_id = ? AND timestamp < ?
		         ORDER BY timestamp DESC
		         LIMIT ?`
		args = []interface{}{streamID, *beforeTimestamp, limit}
	} else {
		// First page: get most recent reactions
		query = `SELECT stream_id, timestamp, reaction_id, viewer_id, reaction_type
		         FROM stream_reactions
		         WHERE stream_id = ?
		         ORDER BY timestamp DESC
		         LIMIT ?`
		args = []interface{}{streamID, limit}
	}

	iter := r.scyllaSession.Query(query, args...).WithContext(ctx).Iter()
	defer iter.Close()

	reactions := make([]*models.StreamReaction, 0, limit)
	var reaction models.StreamReaction

	for iter.Scan(
		&reaction.StreamID,
		&reaction.Timestamp,
		&reaction.ReactionID,
		&reaction.ViewerID,
		&reaction.ReactionType,
	) {
		// Create a new instance for each iteration
		r := models.StreamReaction{
			StreamID:     reaction.StreamID,
			Timestamp:    reaction.Timestamp,
			ReactionID:   reaction.ReactionID,
			ViewerID:     reaction.ViewerID,
			ReactionType: reaction.ReactionType,
		}
		reactions = append(reactions, &r)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to list reactions: %w", err)
	}

	return reactions, nil
}

// GetReactionCounts retrieves aggregated reaction counts from Redis
func (r *streamReactionRepository) GetReactionCounts(ctx context.Context, streamID uuid.UUID) (map[string]int, error) {
	// Use Redis hash to store reaction counts
	hashKey := fmt.Sprintf("reactions:%s", streamID)

	// Get all fields from hash
	result, err := r.redisClient.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get reaction counts: %w", err)
	}

	// Convert string values to int
	counts := make(map[string]int, len(result))
	for reactionType, countStr := range result {
		var count int
		fmt.Sscanf(countStr, "%d", &count)
		counts[reactionType] = count
	}

	return counts, nil
}

// Delete removes a reaction from ScyllaDB
func (r *streamReactionRepository) Delete(ctx context.Context, streamID, reactionID uuid.UUID) error {
	// First, get the reaction to know its type for Redis decrement
	reaction, err := r.GetByID(ctx, streamID, reactionID)
	if err != nil {
		return err
	}

	// Delete from ScyllaDB
	query := `DELETE FROM stream_reactions
	          WHERE stream_id = ? AND timestamp = ? AND reaction_id = ?`

	if err := r.scyllaSession.Query(query, streamID, reaction.Timestamp, reactionID).
		WithContext(ctx).
		Exec(); err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}

	// Decrement Redis counter
	if err := r.DecrementReactionCount(ctx, streamID, reaction.ReactionType); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to decrement reaction count: %v\n", err)
	}

	return nil
}

// IncrementReactionCount increments the Redis counter for a reaction type
func (r *streamReactionRepository) IncrementReactionCount(ctx context.Context, streamID uuid.UUID, reactionType string) error {
	hashKey := fmt.Sprintf("reactions:%s", streamID)

	// Increment the counter for this reaction type
	if err := r.redisClient.HIncrBy(ctx, hashKey, reactionType, 1).Err(); err != nil {
		return fmt.Errorf("failed to increment reaction count: %w", err)
	}

	return nil
}

// DecrementReactionCount decrements the Redis counter for a reaction type
func (r *streamReactionRepository) DecrementReactionCount(ctx context.Context, streamID uuid.UUID, reactionType string) error {
	hashKey := fmt.Sprintf("reactions:%s", streamID)

	// Decrement the counter for this reaction type
	count, err := r.redisClient.HIncrBy(ctx, hashKey, reactionType, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to decrement reaction count: %w", err)
	}

	// Remove the field if count reaches 0
	if count <= 0 {
		r.redisClient.HDel(ctx, hashKey, reactionType)
	}

	return nil
}

// SetReactionTTL sets TTL on Redis counters matching stream duration
func (r *streamReactionRepository) SetReactionTTL(ctx context.Context, streamID uuid.UUID, ttl time.Duration) error {
	hashKey := fmt.Sprintf("reactions:%s", streamID)

	if err := r.redisClient.Expire(ctx, hashKey, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set reaction TTL: %w", err)
	}

	return nil
}
