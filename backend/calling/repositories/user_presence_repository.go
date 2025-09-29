package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"tchat.dev/calling/models"
)

// UserPresenceRepository defines the interface for user presence data access
type UserPresenceRepository interface {
	GetByUserID(userID uuid.UUID) (*models.UserPresence, error)
	Update(presence *models.UserPresence) error
	Delete(userID uuid.UUID) error
	SetOnline(userID uuid.UUID) error
	SetOffline(userID uuid.UUID) error
	SetInCall(userID uuid.UUID) error
	SetBusy(userID uuid.UUID) error
	GetOnlineUsers() ([]models.UserPresence, error)
	GetUsersByStatus(status models.PresenceStatus) ([]models.UserPresence, error)
	IsUserOnline(userID uuid.UUID) (bool, error)
	GetBulkPresence(userIDs []uuid.UUID) (map[uuid.UUID]*models.UserPresence, error)
}

// RedisUserPresenceRepository implements UserPresenceRepository using Redis
type RedisUserPresenceRepository struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisUserPresenceRepository creates a new Redis-based user presence repository
func NewRedisUserPresenceRepository(client *redis.Client, ttl time.Duration) UserPresenceRepository {
	if ttl == 0 {
		ttl = 30 * time.Minute // Default TTL for presence data
	}
	return &RedisUserPresenceRepository{
		client: client,
		ttl:    ttl,
	}
}

// getUserPresenceKey returns the Redis key for a user's presence
func (r *RedisUserPresenceRepository) getUserPresenceKey(userID uuid.UUID) string {
	return fmt.Sprintf("presence:user:%s", userID.String())
}

// getStatusSetKey returns the Redis key for a status set
func (r *RedisUserPresenceRepository) getStatusSetKey(status models.PresenceStatus) string {
	return fmt.Sprintf("presence:status:%s", string(status))
}

// GetByUserID retrieves user presence by user ID
func (r *RedisUserPresenceRepository) GetByUserID(userID uuid.UUID) (*models.UserPresence, error) {
	ctx := context.Background()
	key := r.getUserPresenceKey(userID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// User not found, return default offline presence
			return &models.UserPresence{
				UserID:    userID,
				Status:    models.PresenceStatusOffline,
				LastSeen:  time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, err
	}

	var presence models.UserPresence
	if err := json.Unmarshal([]byte(data), &presence); err != nil {
		return nil, err
	}

	return &presence, nil
}

// Update updates user presence data
func (r *RedisUserPresenceRepository) Update(presence *models.UserPresence) error {
	ctx := context.Background()
	key := r.getUserPresenceKey(presence.UserID)

	// Remove from old status set if it exists
	oldPresence, _ := r.GetByUserID(presence.UserID)
	if oldPresence != nil && oldPresence.Status != presence.Status {
		oldStatusKey := r.getStatusSetKey(oldPresence.Status)
		r.client.SRem(ctx, oldStatusKey, presence.UserID.String())
	}

	// Update timestamp
	presence.UpdatedAt = time.Now()
	if presence.Status == models.PresenceStatusOnline {
		presence.LastSeen = time.Now()
	}

	// Serialize presence data
	data, err := json.Marshal(presence)
	if err != nil {
		return err
	}

	// Store in Redis with TTL
	if err := r.client.Set(ctx, key, data, r.ttl).Err(); err != nil {
		return err
	}

	// Add to status set
	statusKey := r.getStatusSetKey(presence.Status)
	if err := r.client.SAdd(ctx, statusKey, presence.UserID.String()).Err(); err != nil {
		return err
	}

	// Set TTL on status set as well
	r.client.Expire(ctx, statusKey, r.ttl)

	return nil
}

// Delete removes user presence data
func (r *RedisUserPresenceRepository) Delete(userID uuid.UUID) error {
	ctx := context.Background()

	// Get current presence to remove from status sets
	presence, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	// Remove from status set
	statusKey := r.getStatusSetKey(presence.Status)
	r.client.SRem(ctx, statusKey, userID.String())

	// Delete presence key
	key := r.getUserPresenceKey(userID)
	return r.client.Del(ctx, key).Err()
}

// SetOnline sets a user as online
func (r *RedisUserPresenceRepository) SetOnline(userID uuid.UUID) error {
	presence := &models.UserPresence{
		UserID:    userID,
		Status:    models.PresenceStatusOnline,
		LastSeen:  time.Now(),
		UpdatedAt: time.Now(),
	}
	return r.Update(presence)
}

// SetOffline sets a user as offline
func (r *RedisUserPresenceRepository) SetOffline(userID uuid.UUID) error {
	presence, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}
	presence.SetOffline()
	return r.Update(presence)
}

// SetInCall sets a user as in a call
func (r *RedisUserPresenceRepository) SetInCall(userID uuid.UUID) error {
	presence, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}
	presence.SetInCall()
	return r.Update(presence)
}

// SetBusy sets a user as busy
func (r *RedisUserPresenceRepository) SetBusy(userID uuid.UUID) error {
	presence, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}
	presence.SetBusy()
	return r.Update(presence)
}

// GetOnlineUsers retrieves all online users
func (r *RedisUserPresenceRepository) GetOnlineUsers() ([]models.UserPresence, error) {
	return r.GetUsersByStatus(models.PresenceStatusOnline)
}

// GetUsersByStatus retrieves all users with a specific status
func (r *RedisUserPresenceRepository) GetUsersByStatus(status models.PresenceStatus) ([]models.UserPresence, error) {
	ctx := context.Background()
	statusKey := r.getStatusSetKey(status)

	userIDs, err := r.client.SMembers(ctx, statusKey).Result()
	if err != nil {
		return nil, err
	}

	var presences []models.UserPresence
	for _, userIDStr := range userIDs {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			continue // Skip invalid UUIDs
		}

		presence, err := r.GetByUserID(userID)
		if err != nil {
			continue // Skip errored users
		}

		// Double-check status matches (in case of stale data)
		if presence.Status == status {
			presences = append(presences, *presence)
		} else {
			// Clean up stale data
			r.client.SRem(ctx, statusKey, userIDStr)
		}
	}

	return presences, nil
}

// IsUserOnline checks if a user is currently online
func (r *RedisUserPresenceRepository) IsUserOnline(userID uuid.UUID) (bool, error) {
	presence, err := r.GetByUserID(userID)
	if err != nil {
		return false, err
	}
	return presence.IsOnline(), nil
}

// GetBulkPresence retrieves presence for multiple users efficiently
func (r *RedisUserPresenceRepository) GetBulkPresence(userIDs []uuid.UUID) (map[uuid.UUID]*models.UserPresence, error) {
	ctx := context.Background()
	result := make(map[uuid.UUID]*models.UserPresence)

	if len(userIDs) == 0 {
		return result, nil
	}

	// Prepare keys for bulk get
	keys := make([]string, len(userIDs))
	for i, userID := range userIDs {
		keys[i] = r.getUserPresenceKey(userID)
	}

	// Bulk get from Redis
	values, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	// Parse results
	for i, value := range values {
		userID := userIDs[i]

		if value == nil {
			// User not found, set default offline presence
			result[userID] = &models.UserPresence{
				UserID:    userID,
				Status:    models.PresenceStatusOffline,
				LastSeen:  time.Now(),
				UpdatedAt: time.Now(),
			}
			continue
		}

		var presence models.UserPresence
		if err := json.Unmarshal([]byte(value.(string)), &presence); err == nil {
			result[userID] = &presence
		} else {
			// Failed to parse, set default offline presence
			result[userID] = &models.UserPresence{
				UserID:    userID,
				Status:    models.PresenceStatusOffline,
				LastSeen:  time.Now(),
				UpdatedAt: time.Now(),
			}
		}
	}

	return result, nil
}

// CleanupExpiredPresence removes expired presence entries
func (r *RedisUserPresenceRepository) CleanupExpiredPresence() error {
	ctx := context.Background()

	// Get all status sets
	statuses := []models.PresenceStatus{
		models.PresenceStatusOnline,
		models.PresenceStatusBusy,
		models.PresenceStatusInCall,
		models.PresenceStatusOffline,
	}

	for _, status := range statuses {
		statusKey := r.getStatusSetKey(status)
		userIDs, _ := r.client.SMembers(ctx, statusKey).Result()

		for _, userIDStr := range userIDs {
			userKey := fmt.Sprintf("presence:user:%s", userIDStr)

			// Check if user key exists
			exists, _ := r.client.Exists(ctx, userKey).Result()
			if exists == 0 {
				// Remove from status set if user key doesn't exist
				r.client.SRem(ctx, statusKey, userIDStr)
			}
		}
	}

	return nil
}
