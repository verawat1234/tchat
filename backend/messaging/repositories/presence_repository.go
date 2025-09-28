package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/messaging/models"
	"tchat.dev/messaging/services"
)

// PresenceRepository implements services.PresenceRepository using GORM
type PresenceRepository struct {
	db *gorm.DB
}

// NewPresenceRepository creates a new presence repository
func NewPresenceRepository(db *gorm.DB) services.PresenceRepository {
	return &PresenceRepository{db: db}
}

// Create creates a new presence record in the database
func (r *PresenceRepository) Create(ctx context.Context, presence *models.Presence) error {
	return r.db.WithContext(ctx).Create(presence).Error
}

// GetByUserID retrieves presence information for a specific user
func (r *PresenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Presence, error) {
	var presence models.Presence
	err := r.db.WithContext(ctx).First(&presence, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &presence, nil
}

// Update updates an existing presence record
func (r *PresenceRepository) Update(ctx context.Context, presence *models.Presence) error {
	return r.db.WithContext(ctx).Save(presence).Error
}

// GetByUserIDs retrieves presence information for multiple users
func (r *PresenceRepository) GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.Presence, error) {
	var presences []*models.Presence
	err := r.db.WithContext(ctx).
		Where("user_id IN ?", userIDs).
		Find(&presences).Error
	return presences, err
}

// GetOnlineUsers retrieves a list of currently online users
func (r *PresenceRepository) GetOnlineUsers(ctx context.Context, limit int) ([]*models.Presence, error) {
	var presences []*models.Presence
	err := r.db.WithContext(ctx).
		Where("is_online = true").
		Order("last_updated DESC").
		Limit(limit).
		Find(&presences).Error
	return presences, err
}

// CleanupStalePresence marks users as offline if they haven't been seen for a while
func (r *PresenceRepository) CleanupStalePresence(ctx context.Context, staleThreshold time.Duration) error {
	threshold := time.Now().Add(-staleThreshold)
	return r.db.WithContext(ctx).Model(&models.Presence{}).
		Where("last_updated < ? AND is_online = true", threshold).
		Updates(map[string]interface{}{
			"is_online":    false,
			"last_updated": time.Now(),
		}).Error
}

// GetPresenceStats returns statistics about user presence
func (r *PresenceRepository) GetPresenceStats(ctx context.Context) (*services.PresenceStats, error) {
	var stats services.PresenceStats

	// Total users
	err := r.db.WithContext(ctx).Model(&models.Presence{}).Count(&stats.TotalUsers).Error
	if err != nil {
		return nil, err
	}

	// Online users
	err = r.db.WithContext(ctx).Model(&models.Presence{}).
		Where("is_online = true").Count(&stats.OnlineUsers).Error
	if err != nil {
		return nil, err
	}

	// Away users
	err = r.db.WithContext(ctx).Model(&models.Presence{}).
		Where("status = ?", models.PresenceStatusAway).Count(&stats.AwayUsers).Error
	if err != nil {
		return nil, err
	}

	// Busy users
	err = r.db.WithContext(ctx).Model(&models.Presence{}).
		Where("status = ?", models.PresenceStatusBusy).Count(&stats.BusyUsers).Error
	if err != nil {
		return nil, err
	}

	// Offline users
	err = r.db.WithContext(ctx).Model(&models.Presence{}).
		Where("is_online = false OR is_online IS NULL").Count(&stats.OfflineUsers).Error
	if err != nil {
		return nil, err
	}

	// Average uptime (calculate from session start times)
	var avgUptimeMinutes float64
	r.db.WithContext(ctx).Model(&models.Presence{}).
		Where("is_online = true AND last_updated > NOW() - INTERVAL '24 hours'").
		Select("AVG(EXTRACT(EPOCH FROM (NOW() - created_at))/60)").
		Row().Scan(&avgUptimeMinutes)
	stats.AverageUptime = time.Duration(avgUptimeMinutes) * time.Minute

	// Peak online time and count (simplified calculation)
	var peakTime time.Time
	var peakCount int64

	// Get the hour with most online activity in the last 24 hours
	row := r.db.WithContext(ctx).Raw(`
		SELECT
			date_trunc('hour', last_updated) AS hour,
			COUNT(*) AS count
		FROM presences
		WHERE last_updated > NOW() - INTERVAL '24 hours'
		GROUP BY hour
		ORDER BY count DESC
		LIMIT 1
	`).Row()

	row.Scan(&peakTime, &peakCount)
	stats.PeakOnlineTime = peakTime
	stats.PeakOnlineCount = peakCount

	return &stats, nil
}