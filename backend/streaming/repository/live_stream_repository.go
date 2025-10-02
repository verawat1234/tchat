// backend/streaming/repository/live_stream_repository.go
// LiveStream Repository - Data access layer for live streaming operations
// Implements T025: LiveStreamRepository with GORM PostgreSQL integration

package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/streaming/models"
)

// LiveStreamRepositoryInterface defines the contract for live stream data operations
type LiveStreamRepositoryInterface interface {
	// Core CRUD operations
	Create(stream *models.LiveStream) error
	GetByID(streamID uuid.UUID) (*models.LiveStream, error)
	List(filters map[string]interface{}, limit, offset int) ([]*models.LiveStream, int64, error)
	Update(streamID uuid.UUID, updates map[string]interface{}) error
	Delete(streamID uuid.UUID) error

	// Specialized queries
	GetByStreamKey(streamKey string) (*models.LiveStream, error)
	ListByBroadcaster(broadcasterID uuid.UUID, limit, offset int) ([]*models.LiveStream, int64, error)
	GetLiveStreams(limit, offset int) ([]*models.LiveStream, int64, error)

	// Status management
	UpdateStatus(streamID uuid.UUID, status string) error
	GetByStatus(status string, limit, offset int) ([]*models.LiveStream, int64, error)

	// Viewer management
	IncrementViewerCount(streamID uuid.UUID) error
	DecrementViewerCount(streamID uuid.UUID) error
	UpdatePeakViewerCount(streamID uuid.UUID, currentCount int) error

	// Recording management
	UpdateRecordingURL(streamID uuid.UUID, recordingURL string, expiryDate time.Time) error
	GetExpiredRecordings(cutoffDate time.Time) ([]*models.LiveStream, error)
	CleanupExpiredRecordings() (int64, error)

	// Analytics and discovery
	GetStreamsByTags(tags []string, limit, offset int) ([]*models.LiveStream, int64, error)
	GetScheduledStreams(afterTime time.Time, limit, offset int) ([]*models.LiveStream, int64, error)
	GetStreamAnalytics(streamID uuid.UUID) (*StreamAnalytics, error)
	GetBroadcasterAnalytics(broadcasterID uuid.UUID) (*BroadcasterAnalytics, error)

	// Health and maintenance
	GetRepositoryHealth() (*RepositoryHealth, error)
}

// LiveStreamRepository implements LiveStreamRepositoryInterface
type LiveStreamRepository struct {
	db *gorm.DB
}

// NewLiveStreamRepository creates a new live stream repository instance
func NewLiveStreamRepository(db *gorm.DB) LiveStreamRepositoryInterface {
	return &LiveStreamRepository{db: db}
}

// Core CRUD operations

func (r *LiveStreamRepository) Create(stream *models.LiveStream) error {
	// Validate required fields
	if stream.BroadcasterID == uuid.Nil {
		return fmt.Errorf("broadcaster_id is required")
	}

	if stream.StreamType == "" {
		return fmt.Errorf("stream_type is required")
	}

	if stream.Title == "" {
		return fmt.Errorf("title is required")
	}

	if stream.StreamKey == "" {
		return fmt.Errorf("stream_key is required")
	}

	// Set defaults
	if stream.Status == "" {
		stream.Status = "scheduled"
	}

	if stream.MaxCapacity == 0 {
		stream.MaxCapacity = 50000
	}

	if stream.Language == "" {
		stream.Language = "en"
	}

	return r.db.Create(stream).Error
}

func (r *LiveStreamRepository) GetByID(streamID uuid.UUID) (*models.LiveStream, error) {
	var stream models.LiveStream
	err := r.db.First(&stream, "id = ?", streamID).Error
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

func (r *LiveStreamRepository) List(filters map[string]interface{}, limit, offset int) ([]*models.LiveStream, int64, error) {
	var streams []*models.LiveStream
	var total int64

	db := r.db.Model(&models.LiveStream{})

	// Apply filters
	if status, ok := filters["status"]; ok {
		db = db.Where("status = ?", status)
	}

	if streamType, ok := filters["stream_type"]; ok {
		db = db.Where("stream_type = ?", streamType)
	}

	if broadcasterID, ok := filters["broadcaster_id"]; ok {
		db = db.Where("broadcaster_id = ?", broadcasterID)
	}

	if privacySetting, ok := filters["privacy_setting"]; ok {
		db = db.Where("privacy_setting = ?", privacySetting)
	}

	if kycTier, ok := filters["broadcaster_kyc_tier"]; ok {
		db = db.Where("broadcaster_kyc_tier >= ?", kycTier)
	}

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := db.Order("scheduled_start_time DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&streams).Error

	return streams, total, err
}

func (r *LiveStreamRepository) Update(streamID uuid.UUID, updates map[string]interface{}) error {
	// Validate streamID
	if streamID == uuid.Nil {
		return fmt.Errorf("stream_id is required")
	}

	// Ensure updated_at is set
	updates["updated_at"] = time.Now()

	result := r.db.Model(&models.LiveStream{}).
		Where("id = ?", streamID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream not found: %s", streamID)
	}

	return nil
}

func (r *LiveStreamRepository) Delete(streamID uuid.UUID) error {
	result := r.db.Delete(&models.LiveStream{}, "id = ?", streamID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream not found: %s", streamID)
	}

	return nil
}

// Specialized queries

func (r *LiveStreamRepository) GetByStreamKey(streamKey string) (*models.LiveStream, error) {
	var stream models.LiveStream
	err := r.db.Where("stream_key = ?", streamKey).First(&stream).Error
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

func (r *LiveStreamRepository) ListByBroadcaster(broadcasterID uuid.UUID, limit, offset int) ([]*models.LiveStream, int64, error) {
	var streams []*models.LiveStream
	var total int64

	db := r.db.Model(&models.LiveStream{}).Where("broadcaster_id = ?", broadcasterID)

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := db.Order("scheduled_start_time DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&streams).Error

	return streams, total, err
}

func (r *LiveStreamRepository) GetLiveStreams(limit, offset int) ([]*models.LiveStream, int64, error) {
	var streams []*models.LiveStream
	var total int64

	db := r.db.Model(&models.LiveStream{}).Where("status = ?", "live")

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results ordered by viewer count
	err := db.Order("viewer_count DESC, actual_start_time DESC").
		Limit(limit).
		Offset(offset).
		Find(&streams).Error

	return streams, total, err
}

// Status management

func (r *LiveStreamRepository) UpdateStatus(streamID uuid.UUID, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// Set actual_start_time when going live
	if status == "live" {
		updates["actual_start_time"] = time.Now()
	}

	// Set end_time when ending or terminating
	if status == "ended" || status == "terminated" {
		updates["end_time"] = time.Now()
	}

	return r.Update(streamID, updates)
}

func (r *LiveStreamRepository) GetByStatus(status string, limit, offset int) ([]*models.LiveStream, int64, error) {
	filters := map[string]interface{}{
		"status": status,
	}
	return r.List(filters, limit, offset)
}

// Viewer management

func (r *LiveStreamRepository) IncrementViewerCount(streamID uuid.UUID) error {
	return r.db.Model(&models.LiveStream{}).
		Where("id = ?", streamID).
		UpdateColumn("viewer_count", gorm.Expr("viewer_count + ?", 1)).Error
}

func (r *LiveStreamRepository) DecrementViewerCount(streamID uuid.UUID) error {
	return r.db.Model(&models.LiveStream{}).
		Where("id = ? AND viewer_count > 0", streamID).
		UpdateColumn("viewer_count", gorm.Expr("viewer_count - ?", 1)).Error
}

func (r *LiveStreamRepository) UpdatePeakViewerCount(streamID uuid.UUID, currentCount int) error {
	return r.db.Model(&models.LiveStream{}).
		Where("id = ? AND peak_viewer_count < ?", streamID, currentCount).
		UpdateColumn("peak_viewer_count", currentCount).Error
}

// Recording management

func (r *LiveStreamRepository) UpdateRecordingURL(streamID uuid.UUID, recordingURL string, expiryDate time.Time) error {
	updates := map[string]interface{}{
		"recording_url":         recordingURL,
		"recording_expiry_date": expiryDate,
		"updated_at":            time.Now(),
	}
	return r.Update(streamID, updates)
}

func (r *LiveStreamRepository) GetExpiredRecordings(cutoffDate time.Time) ([]*models.LiveStream, error) {
	var streams []*models.LiveStream
	err := r.db.Where("recording_expiry_date < ? AND recording_url IS NOT NULL", cutoffDate).
		Find(&streams).Error
	return streams, err
}

func (r *LiveStreamRepository) CleanupExpiredRecordings() (int64, error) {
	cutoffDate := time.Now()
	result := r.db.Model(&models.LiveStream{}).
		Where("recording_expiry_date < ?", cutoffDate).
		Updates(map[string]interface{}{
			"recording_url":         nil,
			"recording_expiry_date": nil,
			"updated_at":            time.Now(),
		})
	return result.RowsAffected, result.Error
}

// Analytics and discovery

func (r *LiveStreamRepository) GetStreamsByTags(tags []string, limit, offset int) ([]*models.LiveStream, int64, error) {
	var streams []*models.LiveStream
	var total int64

	db := r.db.Model(&models.LiveStream{}).
		Where("tags @> ?", tags)

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&streams).Error

	return streams, total, err
}

func (r *LiveStreamRepository) GetScheduledStreams(afterTime time.Time, limit, offset int) ([]*models.LiveStream, int64, error) {
	var streams []*models.LiveStream
	var total int64

	db := r.db.Model(&models.LiveStream{}).
		Where("status = ? AND scheduled_start_time > ?", "scheduled", afterTime)

	// Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results ordered by scheduled time
	err := db.Order("scheduled_start_time ASC").
		Limit(limit).
		Offset(offset).
		Find(&streams).Error

	return streams, total, err
}

func (r *LiveStreamRepository) GetStreamAnalytics(streamID uuid.UUID) (*StreamAnalytics, error) {
	var stream models.LiveStream
	if err := r.db.First(&stream, "id = ?", streamID).Error; err != nil {
		return nil, err
	}

	analytics := &StreamAnalytics{
		StreamID:          stream.ID,
		Title:             stream.Title,
		StreamType:        stream.StreamType,
		Status:            stream.Status,
		ViewerCount:       stream.ViewerCount,
		PeakViewerCount:   stream.PeakViewerCount,
		ScheduledStart:    stream.ScheduledStartTime,
		ActualStart:       stream.ActualStartTime,
		EndTime:           stream.EndTime,
		Language:          stream.Language,
		Tags:              stream.Tags,
	}

	// Calculate duration if stream has ended
	if stream.ActualStartTime != nil && stream.EndTime != nil {
		duration := stream.EndTime.Sub(*stream.ActualStartTime)
		analytics.DurationSeconds = int(duration.Seconds())
	}

	return analytics, nil
}

func (r *LiveStreamRepository) GetBroadcasterAnalytics(broadcasterID uuid.UUID) (*BroadcasterAnalytics, error) {
	analytics := &BroadcasterAnalytics{
		BroadcasterID: broadcasterID,
	}

	// Get total streams count
	var totalStreams int64
	if err := r.db.Model(&models.LiveStream{}).
		Where("broadcaster_id = ?", broadcasterID).
		Count(&totalStreams).Error; err != nil {
		return nil, err
	}
	analytics.TotalStreams = totalStreams

	// Get live streams count
	var liveStreams int64
	if err := r.db.Model(&models.LiveStream{}).
		Where("broadcaster_id = ? AND status = ?", broadcasterID, "live").
		Count(&liveStreams).Error; err != nil {
		return nil, err
	}
	analytics.LiveStreams = liveStreams

	// Get aggregate metrics
	var aggregates struct {
		TotalViewers    int64   `json:"total_viewers"`
		MaxPeakViewers  int     `json:"max_peak_viewers"`
		AvgViewerCount  float64 `json:"avg_viewer_count"`
	}

	err := r.db.Model(&models.LiveStream{}).
		Select(`
			SUM(peak_viewer_count) as total_viewers,
			MAX(peak_viewer_count) as max_peak_viewers,
			AVG(NULLIF(peak_viewer_count, 0)) as avg_viewer_count
		`).
		Where("broadcaster_id = ? AND status != ?", broadcasterID, "scheduled").
		Scan(&aggregates).Error

	if err != nil {
		return nil, err
	}

	analytics.TotalViewers = aggregates.TotalViewers
	analytics.MaxPeakViewers = aggregates.MaxPeakViewers
	analytics.AvgViewerCount = aggregates.AvgViewerCount

	// Get most recent stream
	var mostRecentStream models.LiveStream
	if err := r.db.Where("broadcaster_id = ?", broadcasterID).
		Order("created_at DESC").
		First(&mostRecentStream).Error; err == nil {
		analytics.MostRecentStream = &mostRecentStream.CreatedAt
	}

	return analytics, nil
}

// Health and maintenance

func (r *LiveStreamRepository) GetRepositoryHealth() (*RepositoryHealth, error) {
	var health RepositoryHealth

	// Check database connection
	sqlDB, err := r.db.DB()
	if err != nil {
		health.DatabaseConnected = false
		health.Errors = append(health.Errors, "Failed to get database connection")
	} else {
		err = sqlDB.Ping()
		health.DatabaseConnected = err == nil
		if err != nil {
			health.Errors = append(health.Errors, fmt.Sprintf("Database ping failed: %v", err))
		}
	}

	// Get table counts by status
	var counts struct {
		Total       int64 `json:"total"`
		Scheduled   int64 `json:"scheduled"`
		Live        int64 `json:"live"`
		Ended       int64 `json:"ended"`
		Terminated  int64 `json:"terminated"`
	}

	r.db.Model(&models.LiveStream{}).Count(&counts.Total)
	r.db.Model(&models.LiveStream{}).Where("status = ?", "scheduled").Count(&counts.Scheduled)
	r.db.Model(&models.LiveStream{}).Where("status = ?", "live").Count(&counts.Live)
	r.db.Model(&models.LiveStream{}).Where("status = ?", "ended").Count(&counts.Ended)
	r.db.Model(&models.LiveStream{}).Where("status = ?", "terminated").Count(&counts.Terminated)

	health.TableCounts = map[string]int64{
		"total":      counts.Total,
		"scheduled":  counts.Scheduled,
		"live":       counts.Live,
		"ended":      counts.Ended,
		"terminated": counts.Terminated,
	}

	// Get active viewers count
	var activeViewers int64
	r.db.Model(&models.LiveStream{}).
		Where("status = ?", "live").
		Select("COALESCE(SUM(viewer_count), 0)").
		Scan(&activeViewers)
	health.ActiveViewers = activeViewers

	health.Healthy = health.DatabaseConnected && len(health.Errors) == 0
	health.CheckedAt = time.Now()

	return &health, nil
}

// Supporting types for analytics

type StreamAnalytics struct {
	StreamID        uuid.UUID      `json:"stream_id"`
	Title           string         `json:"title"`
	StreamType      string         `json:"stream_type"`
	Status          string         `json:"status"`
	ViewerCount     int            `json:"viewer_count"`
	PeakViewerCount int            `json:"peak_viewer_count"`
	ScheduledStart  *time.Time     `json:"scheduled_start"`
	ActualStart     *time.Time     `json:"actual_start"`
	EndTime         *time.Time     `json:"end_time"`
	DurationSeconds int            `json:"duration_seconds"`
	Language        string         `json:"language"`
	Tags            []string       `json:"tags"`
}

type BroadcasterAnalytics struct {
	BroadcasterID   uuid.UUID  `json:"broadcaster_id"`
	TotalStreams    int64      `json:"total_streams"`
	LiveStreams     int64      `json:"live_streams"`
	TotalViewers    int64      `json:"total_viewers"`
	MaxPeakViewers  int        `json:"max_peak_viewers"`
	AvgViewerCount  float64    `json:"avg_viewer_count"`
	MostRecentStream *time.Time `json:"most_recent_stream,omitempty"`
}

type RepositoryHealth struct {
	Healthy           bool              `json:"healthy"`
	DatabaseConnected bool              `json:"database_connected"`
	TableCounts       map[string]int64  `json:"table_counts"`
	ActiveViewers     int64             `json:"active_viewers"`
	Errors            []string          `json:"errors"`
	CheckedAt         time.Time         `json:"checked_at"`
}