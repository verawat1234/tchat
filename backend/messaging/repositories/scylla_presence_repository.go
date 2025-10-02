package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"

	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
)

// ScyllaPresenceRepository implements services.PresenceRepository using ScyllaDB
type ScyllaPresenceRepository struct {
	session *gocql.Session
}

// NewScyllaPresenceRepository creates a new ScyllaDB presence repository
func NewScyllaPresenceRepository(session *gocql.Session) services.PresenceRepository {
	return &ScyllaPresenceRepository{session: session}
}

// Create creates a new presence record in ScyllaDB
func (r *ScyllaPresenceRepository) Create(ctx context.Context, presence *models.Presence) error {
	query := `INSERT INTO presence (user_id, status, last_seen, location, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	// Convert location to string if present
	var locationStr *string
	if presence.Location != nil {
		locStr := fmt.Sprintf("%s,%s", presence.Location.City, presence.Location.Country)
		locationStr = &locStr
	}

	return r.session.Query(query,
		presence.UserID,
		string(presence.Status),
		presence.LastSeen,
		locationStr,
		presence.CreatedAt,
		presence.UpdatedAt,
	).WithContext(ctx).Exec()
}

// GetByUserID retrieves presence for a specific user
func (r *ScyllaPresenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Presence, error) {
	query := `SELECT user_id, status, last_seen, location, created_at, updated_at
		FROM presence WHERE user_id = ?`

	var (
		uid          uuid.UUID
		status       string
		lastSeen     *time.Time
		locationStr  *string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.session.Query(query, userID).WithContext(ctx).Scan(
		&uid, &status, &lastSeen, &locationStr, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, fmt.Errorf("presence not found")
		}
		return nil, fmt.Errorf("failed to get presence: %w", err)
	}

	presence := &models.Presence{
		UserID:    uid,
		LastSeen:  lastSeen,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// Parse location if present
	if locationStr != nil {
		presence.Location = &models.UserLocation{
			City: *locationStr,
		}
	}

	// Parse status
	presence.Status = models.PresenceStatus(status)

	return presence, nil
}

// Update updates a presence record
func (r *ScyllaPresenceRepository) Update(ctx context.Context, presence *models.Presence) error {
	query := `UPDATE presence SET status = ?, last_seen = ?, location = ?, updated_at = ?
		WHERE user_id = ?`

	// Convert location to string if present
	var locationStr *string
	if presence.Location != nil {
		locStr := fmt.Sprintf("%s,%s", presence.Location.City, presence.Location.Country)
		locationStr = &locStr
	}

	return r.session.Query(query,
		string(presence.Status),
		presence.LastSeen,
		locationStr,
		time.Now(),
		presence.UserID,
	).WithContext(ctx).Exec()
}

// GetByUserIDs retrieves presence for multiple users
func (r *ScyllaPresenceRepository) GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Presence, error) {
	var presences []*models.Presence

	for _, userID := range userIDs {
		presence, err := r.GetByUserID(ctx, userID)
		if err != nil {
			// Skip not found errors
			if err.Error() == "presence not found" {
				continue
			}
			return nil, err
		}
		presences = append(presences, presence)
	}

	return presences, nil
}

// GetOnlineUsers retrieves online users
func (r *ScyllaPresenceRepository) GetOnlineUsers(ctx context.Context, limit int) ([]*models.Presence, error) {
	// Would need a secondary index or materialized view on status
	query := `SELECT user_id, status, last_seen, location, created_at, updated_at
		FROM presence LIMIT ?`

	iter := r.session.Query(query, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var presences []*models.Presence
	var (
		userID      uuid.UUID
		status      string
		lastSeen    *time.Time
		locationStr *string
		createdAt   time.Time
		updatedAt   time.Time
	)

	for iter.Scan(&userID, &status, &lastSeen, &locationStr, &createdAt, &updatedAt) {
		var presStatus models.PresenceStatus
		presStatus = models.PresenceStatus(status)

		// Filter for online users
		if presStatus == models.PresenceStatusOnline {
			presence := &models.Presence{
				UserID:    userID,
				Status:    presStatus,
				LastSeen:  lastSeen,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}

			// Parse location if present
			if locationStr != nil {
				presence.Location = &models.UserLocation{
					City: *locationStr,
				}
			}

			presences = append(presences, presence)
		}
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("query iteration failed: %w", err)
	}

	return presences, nil
}

// CleanupStalePresence removes stale presence records
func (r *ScyllaPresenceRepository) CleanupStalePresence(ctx context.Context, staleThreshold time.Duration) error {
	// ScyllaDB TTL would be better for this
	cutoff := time.Now().Add(-staleThreshold)

	query := `SELECT user_id FROM presence WHERE last_seen < ? ALLOW FILTERING`
	iter := r.session.Query(query, cutoff).WithContext(ctx).Iter()

	var userIDs []uuid.UUID
	var userID uuid.UUID
	for iter.Scan(&userID) {
		userIDs = append(userIDs, userID)
	}

	if err := iter.Close(); err != nil {
		return fmt.Errorf("query iteration failed: %w", err)
	}

	// Delete stale records
	deleteQuery := `DELETE FROM presence WHERE user_id = ?`
	for _, uid := range userIDs {
		if err := r.session.Query(deleteQuery, uid).WithContext(ctx).Exec(); err != nil {
			return fmt.Errorf("failed to delete stale presence: %w", err)
		}
	}

	return nil
}

// GetPresenceStats retrieves presence statistics
func (r *ScyllaPresenceRepository) GetPresenceStats(ctx context.Context) (*services.PresenceStats, error) {
	// Would require aggregation queries or materialized views
	return &services.PresenceStats{}, nil
}
