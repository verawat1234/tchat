package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tchat.dev/streaming/models"
)

// StreamAnalyticsRepository defines the interface for stream analytics data access
type StreamAnalyticsRepository interface {
	// Create creates a new stream analytics record
	Create(ctx context.Context, analytics *models.StreamAnalytics) error

	// GetByStreamID retrieves analytics by stream ID
	GetByStreamID(ctx context.Context, streamID uuid.UUID) (*models.StreamAnalytics, error)

	// Update updates stream analytics with a map of updates
	Update(ctx context.Context, streamID uuid.UUID, updates map[string]interface{}) error

	// IncrementMetric atomically increments a metric by a specified value
	IncrementMetric(ctx context.Context, streamID uuid.UUID, metric string, value int) error

	// UpdatePeakViewers updates the peak concurrent viewers if new count is higher
	UpdatePeakViewers(ctx context.Context, streamID uuid.UUID, count int) error

	// UpdateCommerceMetrics updates commerce-related metrics (store context only)
	UpdateCommerceMetrics(ctx context.Context, streamID uuid.UUID, revenue float64, purchases, products int) error

	// GetTopStreamsByMetric retrieves top streams ordered by a specific metric
	GetTopStreamsByMetric(ctx context.Context, metric string, limit int) ([]*models.StreamAnalytics, error)

	// UpdateViewerCountries updates the geographic distribution JSONB field
	UpdateViewerCountries(ctx context.Context, streamID uuid.UUID, countries map[string]int) error

	// GetAverageMetrics calculates average metrics across all streams
	GetAverageMetrics(ctx context.Context) (map[string]float64, error)

	// DeleteByStreamID deletes analytics for a specific stream (for cleanup)
	DeleteByStreamID(ctx context.Context, streamID uuid.UUID) error
}

// streamAnalyticsRepository implements StreamAnalyticsRepository interface
type streamAnalyticsRepository struct {
	db *gorm.DB
}

// NewStreamAnalyticsRepository creates a new stream analytics repository
func NewStreamAnalyticsRepository(db *gorm.DB) StreamAnalyticsRepository {
	return &streamAnalyticsRepository{
		db: db,
	}
}

// Create creates a new stream analytics record
func (r *streamAnalyticsRepository) Create(ctx context.Context, analytics *models.StreamAnalytics) error {
	if err := r.db.WithContext(ctx).Create(analytics).Error; err != nil {
		return fmt.Errorf("failed to create stream analytics: %w", err)
	}
	return nil
}

// GetByStreamID retrieves analytics by stream ID
func (r *streamAnalyticsRepository) GetByStreamID(ctx context.Context, streamID uuid.UUID) (*models.StreamAnalytics, error) {
	var analytics models.StreamAnalytics

	if err := r.db.WithContext(ctx).Where("stream_id = ?", streamID).First(&analytics).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("stream analytics not found")
		}
		return nil, fmt.Errorf("failed to find stream analytics: %w", err)
	}

	return &analytics, nil
}

// Update updates stream analytics with a map of updates
func (r *streamAnalyticsRepository) Update(ctx context.Context, streamID uuid.UUID, updates map[string]interface{}) error {
	// Add calculated_at timestamp to updates
	updates["calculated_at"] = time.Now()

	result := r.db.WithContext(ctx).
		Model(&models.StreamAnalytics{}).
		Where("stream_id = ?", streamID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update stream analytics: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream analytics not found")
	}

	return nil
}

// IncrementMetric atomically increments a metric by a specified value
func (r *streamAnalyticsRepository) IncrementMetric(ctx context.Context, streamID uuid.UUID, metric string, value int) error {
	// Build update expression with atomic increment
	result := r.db.WithContext(ctx).
		Model(&models.StreamAnalytics{}).
		Where("stream_id = ?", streamID).
		Updates(map[string]interface{}{
			metric:          gorm.Expr(fmt.Sprintf("%s + ?", metric), value),
			"calculated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to increment metric %s: %w", metric, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream analytics not found")
	}

	return nil
}

// UpdatePeakViewers updates the peak concurrent viewers if new count is higher
func (r *streamAnalyticsRepository) UpdatePeakViewers(ctx context.Context, streamID uuid.UUID, count int) error {
	// Use GREATEST function to only update if new count is higher
	result := r.db.WithContext(ctx).
		Model(&models.StreamAnalytics{}).
		Where("stream_id = ?", streamID).
		Updates(map[string]interface{}{
			"peak_concurrent_viewers": gorm.Expr("GREATEST(peak_concurrent_viewers, ?)", count),
			"calculated_at":           time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update peak viewers: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream analytics not found")
	}

	return nil
}

// UpdateCommerceMetrics updates commerce-related metrics (store context only)
func (r *streamAnalyticsRepository) UpdateCommerceMetrics(ctx context.Context, streamID uuid.UUID, revenue float64, purchases, products int) error {
	updates := map[string]interface{}{
		"calculated_at": time.Now(),
	}

	// Increment revenue atomically
	if revenue > 0 {
		updates["total_revenue"] = gorm.Expr("COALESCE(total_revenue, 0) + ?", revenue)
	}

	// Increment purchases atomically
	if purchases > 0 {
		updates["total_purchases"] = gorm.Expr("COALESCE(total_purchases, 0) + ?", purchases)
	}

	// Update products featured count
	if products > 0 {
		updates["products_featured"] = products
	}

	result := r.db.WithContext(ctx).
		Model(&models.StreamAnalytics{}).
		Where("stream_id = ?", streamID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update commerce metrics: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream analytics not found")
	}

	return nil
}

// GetTopStreamsByMetric retrieves top streams ordered by a specific metric
func (r *streamAnalyticsRepository) GetTopStreamsByMetric(ctx context.Context, metric string, limit int) ([]*models.StreamAnalytics, error) {
	var analytics []*models.StreamAnalytics

	// Validate metric name to prevent SQL injection
	validMetrics := map[string]bool{
		"total_unique_viewers":         true,
		"peak_concurrent_viewers":      true,
		"average_watch_duration_seconds": true,
		"total_chat_messages":          true,
		"total_reactions":              true,
		"unique_chatter":               true,
		"total_revenue":                true,
		"total_purchases":              true,
	}

	if !validMetrics[metric] {
		return nil, fmt.Errorf("invalid metric name: %s", metric)
	}

	query := r.db.WithContext(ctx).
		Model(&models.StreamAnalytics{}).
		Order(fmt.Sprintf("%s DESC NULLS LAST", metric))

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&analytics).Error; err != nil {
		return nil, fmt.Errorf("failed to get top streams by %s: %w", metric, err)
	}

	return analytics, nil
}

// UpdateViewerCountries updates the geographic distribution JSONB field
func (r *streamAnalyticsRepository) UpdateViewerCountries(ctx context.Context, streamID uuid.UUID, countries map[string]int) error {
	// Convert map to JSON bytes
	countriesJSON, err := json.Marshal(countries)
	if err != nil {
		return fmt.Errorf("failed to marshal countries: %w", err)
	}

	result := r.db.WithContext(ctx).
		Model(&models.StreamAnalytics{}).
		Where("stream_id = ?", streamID).
		Updates(map[string]interface{}{
			"viewer_countries": countriesJSON,
			"calculated_at":    time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update viewer countries: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream analytics not found")
	}

	return nil
}

// GetAverageMetrics calculates average metrics across all streams
func (r *streamAnalyticsRepository) GetAverageMetrics(ctx context.Context) (map[string]float64, error) {
	type AverageResult struct {
		AvgUniqueViewers    float64 `gorm:"column:avg_unique_viewers"`
		AvgPeakViewers      float64 `gorm:"column:avg_peak_viewers"`
		AvgWatchDuration    float64 `gorm:"column:avg_watch_duration"`
		AvgChatMessages     float64 `gorm:"column:avg_chat_messages"`
		AvgReactions        float64 `gorm:"column:avg_reactions"`
		AvgRevenue          float64 `gorm:"column:avg_revenue"`
	}

	var result AverageResult

	err := r.db.WithContext(ctx).
		Model(&models.StreamAnalytics{}).
		Select(`
			AVG(total_unique_viewers) as avg_unique_viewers,
			AVG(peak_concurrent_viewers) as avg_peak_viewers,
			AVG(average_watch_duration_seconds) as avg_watch_duration,
			AVG(total_chat_messages) as avg_chat_messages,
			AVG(total_reactions) as avg_reactions,
			AVG(COALESCE(total_revenue, 0)) as avg_revenue
		`).
		Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("failed to calculate average metrics: %w", err)
	}

	averages := map[string]float64{
		"avg_unique_viewers":    result.AvgUniqueViewers,
		"avg_peak_viewers":      result.AvgPeakViewers,
		"avg_watch_duration":    result.AvgWatchDuration,
		"avg_chat_messages":     result.AvgChatMessages,
		"avg_reactions":         result.AvgReactions,
		"avg_revenue":           result.AvgRevenue,
	}

	return averages, nil
}

// DeleteByStreamID deletes analytics for a specific stream (for cleanup)
func (r *streamAnalyticsRepository) DeleteByStreamID(ctx context.Context, streamID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("stream_id = ?", streamID).
		Delete(&models.StreamAnalytics{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete stream analytics: %w", result.Error)
	}

	return nil
}

// BatchCreate creates multiple analytics records using upsert (ON CONFLICT)
func (r *streamAnalyticsRepository) BatchCreate(ctx context.Context, analyticsList []*models.StreamAnalytics) error {
	if len(analyticsList) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "stream_id"}},
			UpdateAll: true,
		}).
		Create(&analyticsList).Error

	if err != nil {
		return fmt.Errorf("failed to batch create stream analytics: %w", err)
	}

	return nil
}