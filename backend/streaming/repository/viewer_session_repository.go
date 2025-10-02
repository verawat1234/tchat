package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"gorm.io/gorm"
)

// ViewerSessionRepository defines the interface for viewer session data access
type ViewerSessionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, session *models.ViewerSession) error
	GetByID(ctx context.Context, sessionID uuid.UUID) (*models.ViewerSession, error)
	GetActiveByUser(ctx context.Context, userID uuid.UUID) ([]*models.ViewerSession, error)
	GetByStream(ctx context.Context, streamID uuid.UUID, limit, offset int) ([]*models.ViewerSession, int64, error)
	Update(ctx context.Context, sessionID uuid.UUID, updates map[string]interface{}) error
	EndSession(ctx context.Context, sessionID uuid.UUID) error

	// Analytics operations
	GetViewerCount(ctx context.Context, streamID uuid.UUID) (int, error)
	GetWatchTimeByUser(ctx context.Context, userID uuid.UUID, since time.Time) (int, error)
	GetViewerDemographics(ctx context.Context, streamID uuid.UUID) (map[string]interface{}, error)
}

// viewerSessionRepository implements the ViewerSessionRepository interface
type viewerSessionRepository struct {
	db *gorm.DB
}

// NewViewerSessionRepository creates a new viewer session repository
func NewViewerSessionRepository(db *gorm.DB) ViewerSessionRepository {
	return &viewerSessionRepository{
		db: db,
	}
}

// Create creates a new viewer session
func (r *viewerSessionRepository) Create(ctx context.Context, session *models.ViewerSession) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("failed to create viewer session: %w", err)
	}
	return nil
}

// GetByID retrieves a viewer session by its ID
func (r *viewerSessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.ViewerSession, error) {
	var session models.ViewerSession

	if err := r.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("viewer session not found")
		}
		return nil, fmt.Errorf("failed to find viewer session: %w", err)
	}

	return &session, nil
}

// GetActiveByUser retrieves all active sessions for a user
func (r *viewerSessionRepository) GetActiveByUser(ctx context.Context, userID uuid.UUID) ([]*models.ViewerSession, error) {
	var sessions []*models.ViewerSession

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND left_at IS NULL", userID).
		Order("joined_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to find active sessions: %w", err)
	}

	return sessions, nil
}

// GetByStream retrieves viewer sessions for a specific stream with pagination
func (r *viewerSessionRepository) GetByStream(ctx context.Context, streamID uuid.UUID, limit, offset int) ([]*models.ViewerSession, int64, error) {
	var sessions []*models.ViewerSession
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ViewerSession{}).Where("stream_id = ?", streamID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count viewer sessions: %w", err)
	}

	// Apply pagination and ordering
	if err := query.
		Order("joined_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find viewer sessions: %w", err)
	}

	return sessions, total, nil
}

// Update updates a viewer session
func (r *viewerSessionRepository) Update(ctx context.Context, sessionID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&models.ViewerSession{}).Where("id = ?", sessionID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update viewer session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("viewer session not found")
	}
	return nil
}

// EndSession marks a viewer session as ended
func (r *viewerSessionRepository) EndSession(ctx context.Context, sessionID uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.ViewerSession{}).
		Where("id = ? AND left_at IS NULL", sessionID).
		Update("left_at", now)

	if result.Error != nil {
		return fmt.Errorf("failed to end viewer session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("viewer session not found or already ended")
	}
	return nil
}

// GetViewerCount returns the current viewer count for a stream
func (r *viewerSessionRepository) GetViewerCount(ctx context.Context, streamID uuid.UUID) (int, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.ViewerSession{}).
		Where("stream_id = ? AND left_at IS NULL", streamID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count viewers: %w", err)
	}

	return int(count), nil
}

// GetWatchTimeByUser returns total watch time in seconds for a user since a given date
func (r *viewerSessionRepository) GetWatchTimeByUser(ctx context.Context, userID uuid.UUID, since time.Time) (int, error) {
	var result struct {
		TotalSeconds int
	}

	query := `
		SELECT COALESCE(SUM(
			CASE
				WHEN left_at IS NOT NULL THEN EXTRACT(EPOCH FROM (left_at - joined_at))
				ELSE EXTRACT(EPOCH FROM (NOW() - joined_at))
			END
		), 0)::INTEGER AS total_seconds
		FROM viewer_sessions
		WHERE user_id = $1 AND joined_at >= $2
	`

	if err := r.db.WithContext(ctx).Raw(query, userID, since).Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate watch time: %w", err)
	}

	return result.TotalSeconds, nil
}

// GetViewerDemographics returns demographic information for stream viewers
func (r *viewerSessionRepository) GetViewerDemographics(ctx context.Context, streamID uuid.UUID) (map[string]interface{}, error) {
	// Country distribution
	var countryStats []struct {
		CountryCode string
		Count       int64
	}

	if err := r.db.WithContext(ctx).
		Model(&models.ViewerSession{}).
		Select("country_code, COUNT(*) as count").
		Where("stream_id = ? AND country_code IS NOT NULL", streamID).
		Group("country_code").
		Order("count DESC").
		Limit(10).
		Scan(&countryStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get country distribution: %w", err)
	}

	// Quality layer distribution
	var qualityStats []struct {
		QualityLayer string
		Count        int64
	}

	if err := r.db.WithContext(ctx).
		Model(&models.ViewerSession{}).
		Select("average_quality_layer, COUNT(*) as count").
		Where("stream_id = ? AND average_quality_layer IS NOT NULL", streamID).
		Group("average_quality_layer").
		Order("count DESC").
		Scan(&qualityStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get quality distribution: %w", err)
	}

	// Average watch duration
	var avgDuration struct {
		AvgSeconds float64
	}

	query := `
		SELECT COALESCE(AVG(
			CASE
				WHEN left_at IS NOT NULL THEN EXTRACT(EPOCH FROM (left_at - joined_at))
				ELSE EXTRACT(EPOCH FROM (NOW() - joined_at))
			END
		), 0) AS avg_seconds
		FROM viewer_sessions
		WHERE stream_id = $1
	`

	if err := r.db.WithContext(ctx).Raw(query, streamID).Scan(&avgDuration).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average watch duration: %w", err)
	}

	// Authenticated vs anonymous viewers
	var authStats struct {
		AuthenticatedCount int64
		AnonymousCount     int64
	}

	if err := r.db.WithContext(ctx).
		Model(&models.ViewerSession{}).
		Select(`
			COUNT(CASE WHEN user_id IS NOT NULL THEN 1 END) as authenticated_count,
			COUNT(CASE WHEN user_id IS NULL THEN 1 END) as anonymous_count
		`).
		Where("stream_id = ?", streamID).
		Scan(&authStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get authentication stats: %w", err)
	}

	// Total rebuffer events
	var rebufferStats struct {
		TotalRebuffers int64
		AvgRebuffers   float64
	}

	if err := r.db.WithContext(ctx).
		Model(&models.ViewerSession{}).
		Select(`
			COALESCE(SUM(rebuffer_count), 0) as total_rebuffers,
			COALESCE(AVG(rebuffer_count), 0) as avg_rebuffers
		`).
		Where("stream_id = ?", streamID).
		Scan(&rebufferStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get rebuffer stats: %w", err)
	}

	// Build response map
	demographics := map[string]interface{}{
		"country_distribution": countryStats,
		"quality_distribution": qualityStats,
		"average_watch_duration_seconds": avgDuration.AvgSeconds,
		"authenticated_viewers": authStats.AuthenticatedCount,
		"anonymous_viewers":     authStats.AnonymousCount,
		"total_rebuffers":       rebufferStats.TotalRebuffers,
		"average_rebuffers":     rebufferStats.AvgRebuffers,
	}

	return demographics, nil
}