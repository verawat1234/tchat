// backend/video/repository/video_repository.go
// Video Repository - Data access layer for video operations
// Implements T030: VideoRepository with CRUD operations

package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
	"tchat.dev/video/models"
)

// VideoRepositoryInterface defines the contract for video data operations
type VideoRepositoryInterface interface {
	// Video Content operations
	CreateVideo(video *models.VideoContent) error
	GetVideoByID(id uuid.UUID) (*models.VideoContent, error)
	GetVideosByCreator(creatorID uuid.UUID, limit, offset int) ([]*models.VideoContent, error)
	GetVideosByStatus(status models.AvailabilityStatus, limit, offset int) ([]*models.VideoContent, error)
	GetVideosByTags(tags []string, limit, offset int) ([]*models.VideoContent, error)
	UpdateVideo(video *models.VideoContent) error
	DeleteVideo(id uuid.UUID) error

	// Playback Session operations
	CreatePlaybackSession(session *models.PlaybackSession) error
	GetPlaybackSession(id uuid.UUID) (*models.PlaybackSession, error)
	GetActiveSessionsByUser(userID uuid.UUID) ([]*models.PlaybackSession, error)
	GetSessionsByVideo(videoID uuid.UUID, limit, offset int) ([]*models.PlaybackSession, error)
	UpdatePlaybackSession(session *models.PlaybackSession) error
	CompletePlaybackSession(id uuid.UUID) error
	CleanupInactiveSessions(inactiveThreshold time.Duration) (int64, error)

	// Viewing History operations
	CreateViewingHistory(history *models.ViewingHistory) error
	GetViewingHistoryByUser(userID uuid.UUID, limit, offset int) ([]*models.ViewingHistory, error)
	GetViewingHistoryByVideo(videoID uuid.UUID, limit, offset int) ([]*models.ViewingHistory, error)
	GetRecentViewingHistory(userID uuid.UUID, since time.Time) ([]*models.ViewingHistory, error)
	UpdateViewingHistory(history *models.ViewingHistory) error

	// Platform Configuration operations
	CreatePlatformConfig(config *models.PlatformConfiguration) error
	GetPlatformConfig(platformType models.PlatformType, deviceCategory string) (*models.PlatformConfiguration, error)
	GetActivePlatformConfigs() ([]*models.PlatformConfiguration, error)
	UpdatePlatformConfig(config *models.PlatformConfiguration) error

	// Synchronization State operations
	CreateSyncState(syncState *models.SynchronizationState) error
	GetSyncState(id uuid.UUID) (*models.SynchronizationState, error)
	GetSyncStateBySession(sessionID uuid.UUID) ([]*models.SynchronizationState, error)
	GetSyncStateByDevice(deviceID string) ([]*models.SynchronizationState, error)
	UpdateSyncState(syncState *models.SynchronizationState) error
	GetPendingSyncStates() ([]*models.SynchronizationState, error)
	CleanupOldSyncStates(retentionDays int) (int64, error)

	// Analytics and reporting
	GetVideoAnalytics(videoID uuid.UUID) (*VideoAnalytics, error)
	GetUserAnalytics(userID uuid.UUID) (*UserAnalytics, error)
	GetPlatformAnalytics(platformType models.PlatformType, since time.Time) (*PlatformAnalytics, error)
	GetPopularVideos(limit int, timeframe time.Duration) ([]*models.VideoContent, error)
	GetTrendingVideos(limit int) ([]*models.VideoContent, error)

	// Search and discovery
	SearchVideos(query string, filters VideoSearchFilters, limit, offset int) ([]*models.VideoContent, error)
	GetRecommendations(userID uuid.UUID, limit int) ([]*models.VideoContent, error)
	GetSimilarVideos(videoID uuid.UUID, limit int) ([]*models.VideoContent, error)

	// Health and maintenance
	GetRepositoryHealth() (*RepositoryHealth, error)
	OptimizeDatabase() error
	GetDatabaseStats() (*DatabaseStats, error)
}

// VideoRepository implements VideoRepositoryInterface
type VideoRepository struct {
	db *gorm.DB
}

// NewVideoRepository creates a new video repository instance
func NewVideoRepository(db *gorm.DB) VideoRepositoryInterface {
	return &VideoRepository{db: db}
}

// Video Content operations

func (r *VideoRepository) CreateVideo(video *models.VideoContent) error {
	if err := video.ValidateAvailabilityTransition(video.AvailabilityStatus); err != nil {
		return fmt.Errorf("invalid video status: %w", err)
	}

	return r.db.Create(video).Error
}

func (r *VideoRepository) GetVideoByID(id uuid.UUID) (*models.VideoContent, error) {
	var video models.VideoContent
	err := r.db.Preload("PlaybackSessions").Preload("ViewingHistory").First(&video, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *VideoRepository) GetVideosByCreator(creatorID uuid.UUID, limit, offset int) ([]*models.VideoContent, error) {
	var videos []*models.VideoContent
	err := r.db.Where("creator_id = ?", creatorID).
		Order("upload_timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) GetVideosByStatus(status models.AvailabilityStatus, limit, offset int) ([]*models.VideoContent, error) {
	var videos []*models.VideoContent
	err := r.db.Where("availability_status = ?", status).
		Order("upload_timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) GetVideosByTags(tags []string, limit, offset int) ([]*models.VideoContent, error) {
	var videos []*models.VideoContent
	err := r.db.Where("tags @> ?", tags).
		Order("upload_timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error
	return videos, err
}

func (r *VideoRepository) UpdateVideo(video *models.VideoContent) error {
	return r.db.Save(video).Error
}

func (r *VideoRepository) DeleteVideo(id uuid.UUID) error {
	return r.db.Delete(&models.VideoContent{}, "id = ?", id).Error
}

// Playback Session operations

func (r *VideoRepository) CreatePlaybackSession(session *models.PlaybackSession) error {
	return r.db.Create(session).Error
}

func (r *VideoRepository) GetPlaybackSession(id uuid.UUID) (*models.PlaybackSession, error) {
	var session models.PlaybackSession
	err := r.db.Preload("VideoContent").Preload("SynchronizationStates").First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *VideoRepository) GetActiveSessionsByUser(userID uuid.UUID) ([]*models.PlaybackSession, error) {
	var sessions []*models.PlaybackSession
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
		Preload("VideoContent").
		Order("last_updated DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *VideoRepository) GetSessionsByVideo(videoID uuid.UUID, limit, offset int) ([]*models.PlaybackSession, error) {
	var sessions []*models.PlaybackSession
	err := r.db.Where("video_id = ?", videoID).
		Order("session_start DESC").
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error
	return sessions, err
}

func (r *VideoRepository) UpdatePlaybackSession(session *models.PlaybackSession) error {
	return r.db.Save(session).Error
}

func (r *VideoRepository) CompletePlaybackSession(id uuid.UUID) error {
	return r.db.Model(&models.PlaybackSession{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":     false,
			"session_state": models.SessionCompleted,
			"session_end":   time.Now(),
		}).Error
}

func (r *VideoRepository) CleanupInactiveSessions(inactiveThreshold time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-inactiveThreshold)
	result := r.db.Where("last_updated < ? AND is_active = ?", cutoffTime, false).
		Delete(&models.PlaybackSession{})
	return result.RowsAffected, result.Error
}

// Viewing History operations

func (r *VideoRepository) CreateViewingHistory(history *models.ViewingHistory) error {
	// Calculate engagement score before saving
	history.CalculateEngagementScore()
	return r.db.Create(history).Error
}

func (r *VideoRepository) GetViewingHistoryByUser(userID uuid.UUID, limit, offset int) ([]*models.ViewingHistory, error) {
	var history []*models.ViewingHistory
	err := r.db.Where("user_id = ?", userID).
		Preload("VideoContent").
		Order("watched_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&history).Error
	return history, err
}

func (r *VideoRepository) GetViewingHistoryByVideo(videoID uuid.UUID, limit, offset int) ([]*models.ViewingHistory, error) {
	var history []*models.ViewingHistory
	err := r.db.Where("video_id = ?", videoID).
		Order("watched_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&history).Error
	return history, err
}

func (r *VideoRepository) GetRecentViewingHistory(userID uuid.UUID, since time.Time) ([]*models.ViewingHistory, error) {
	var history []*models.ViewingHistory
	err := r.db.Where("user_id = ? AND watched_at > ?", userID, since).
		Preload("VideoContent").
		Order("watched_at DESC").
		Find(&history).Error
	return history, err
}

func (r *VideoRepository) UpdateViewingHistory(history *models.ViewingHistory) error {
	// Recalculate engagement score
	history.CalculateEngagementScore()
	return r.db.Save(history).Error
}

// Platform Configuration operations

func (r *VideoRepository) CreatePlatformConfig(config *models.PlatformConfiguration) error {
	return r.db.Create(config).Error
}

func (r *VideoRepository) GetPlatformConfig(platformType models.PlatformType, deviceCategory string) (*models.PlatformConfiguration, error) {
	var config models.PlatformConfiguration
	err := r.db.Where("platform_type = ? AND device_category = ? AND is_active = ?",
		platformType, deviceCategory, true).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *VideoRepository) GetActivePlatformConfigs() ([]*models.PlatformConfiguration, error) {
	var configs []*models.PlatformConfiguration
	err := r.db.Where("is_active = ?", true).Find(&configs).Error
	return configs, err
}

func (r *VideoRepository) UpdatePlatformConfig(config *models.PlatformConfiguration) error {
	return r.db.Save(config).Error
}

// Synchronization State operations

func (r *VideoRepository) CreateSyncState(syncState *models.SynchronizationState) error {
	return r.db.Create(syncState).Error
}

func (r *VideoRepository) GetSyncState(id uuid.UUID) (*models.SynchronizationState, error) {
	var syncState models.SynchronizationState
	err := r.db.Preload("PlaybackSession").First(&syncState, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &syncState, nil
}

func (r *VideoRepository) GetSyncStateBySession(sessionID uuid.UUID) ([]*models.SynchronizationState, error) {
	var syncStates []*models.SynchronizationState
	err := r.db.Where("session_id = ?", sessionID).Find(&syncStates).Error
	return syncStates, err
}

func (r *VideoRepository) GetSyncStateByDevice(deviceID string) ([]*models.SynchronizationState, error) {
	var syncStates []*models.SynchronizationState
	err := r.db.Where("device_id = ?", deviceID).
		Order("created_at DESC").
		Find(&syncStates).Error
	return syncStates, err
}

func (r *VideoRepository) UpdateSyncState(syncState *models.SynchronizationState) error {
	return r.db.Save(syncState).Error
}

func (r *VideoRepository) GetPendingSyncStates() ([]*models.SynchronizationState, error) {
	var syncStates []*models.SynchronizationState
	err := r.db.Where("sync_status = ? OR sync_status = ?",
		models.SyncPending, models.SyncInProgress).
		Order("next_sync_time ASC").
		Find(&syncStates).Error
	return syncStates, err
}

func (r *VideoRepository) CleanupOldSyncStates(retentionDays int) (int64, error) {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	result := r.db.Where("created_at < ?", cutoffTime).
		Delete(&models.SynchronizationState{})
	return result.RowsAffected, result.Error
}

// Analytics and reporting

func (r *VideoRepository) GetVideoAnalytics(videoID uuid.UUID) (*VideoAnalytics, error) {
	var analytics VideoAnalytics

	// Get basic video info
	if err := r.db.Select("title, duration, upload_timestamp").
		First(&analytics.Video, "id = ?", videoID).Error; err != nil {
		return nil, err
	}

	// Get view count and engagement metrics
	var viewStats struct {
		TotalViews      int64   `json:"total_views"`
		UniqueViewers   int64   `json:"unique_viewers"`
		AvgCompletionRate float64 `json:"avg_completion_rate"`
		AvgEngagementScore float64 `json:"avg_engagement_score"`
		TotalWatchTime  int64   `json:"total_watch_time"`
	}

	err := r.db.Model(&models.ViewingHistory{}).
		Select(`
			COUNT(*) as total_views,
			COUNT(DISTINCT user_id) as unique_viewers,
			AVG(completion_rate) as avg_completion_rate,
			AVG(engagement_score) as avg_engagement_score,
			SUM(watched_duration) as total_watch_time
		`).
		Where("video_id = ?", videoID).
		Scan(&viewStats).Error

	if err != nil {
		return nil, err
	}

	analytics.TotalViews = viewStats.TotalViews
	analytics.UniqueViewers = viewStats.UniqueViewers
	analytics.AvgCompletionRate = viewStats.AvgCompletionRate
	analytics.AvgEngagementScore = viewStats.AvgEngagementScore
	analytics.TotalWatchTime = viewStats.TotalWatchTime

	return &analytics, nil
}

func (r *VideoRepository) GetUserAnalytics(userID uuid.UUID) (*UserAnalytics, error) {
	var analytics UserAnalytics
	analytics.UserID = userID

	// Get viewing statistics
	var viewStats struct {
		TotalVideosWatched int64   `json:"total_videos_watched"`
		TotalWatchTime     int64   `json:"total_watch_time"`
		AvgSessionDuration float64 `json:"avg_session_duration"`
		PreferredQuality   string  `json:"preferred_quality"`
		PreferredPlatform  string  `json:"preferred_platform"`
	}

	err := r.db.Model(&models.ViewingHistory{}).
		Select(`
			COUNT(DISTINCT video_id) as total_videos_watched,
			SUM(watched_duration) as total_watch_time,
			AVG(watched_duration) as avg_session_duration
		`).
		Where("user_id = ?", userID).
		Scan(&viewStats).Error

	if err != nil {
		return nil, err
	}

	analytics.TotalVideosWatched = viewStats.TotalVideosWatched
	analytics.TotalWatchTime = viewStats.TotalWatchTime
	analytics.AvgSessionDuration = viewStats.AvgSessionDuration

	// Get most used quality and platform
	r.db.Model(&models.ViewingHistory{}).
		Select("quality_watched").
		Where("user_id = ?", userID).
		Group("quality_watched").
		Order("COUNT(*) DESC").
		Limit(1).
		Pluck("quality_watched", &analytics.PreferredQuality)

	r.db.Model(&models.ViewingHistory{}).
		Select("platform_used").
		Where("user_id = ?", userID).
		Group("platform_used").
		Order("COUNT(*) DESC").
		Limit(1).
		Pluck("platform_used", &analytics.PreferredPlatform)

	return &analytics, nil
}

func (r *VideoRepository) GetPlatformAnalytics(platformType models.PlatformType, since time.Time) (*PlatformAnalytics, error) {
	var analytics PlatformAnalytics
	analytics.Platform = string(platformType)
	analytics.Since = since

	// Get session statistics
	var sessionStats struct {
		TotalSessions      int64   `json:"total_sessions"`
		ActiveSessions     int64   `json:"active_sessions"`
		AvgSessionDuration float64 `json:"avg_session_duration"`
		AvgBufferHealth    float64 `json:"avg_buffer_health"`
	}

	err := r.db.Model(&models.PlaybackSession{}).
		Select(`
			COUNT(*) as total_sessions,
			SUM(CASE WHEN is_active THEN 1 ELSE 0 END) as active_sessions,
			AVG(EXTRACT(EPOCH FROM (COALESCE(session_end, NOW()) - session_start))) as avg_session_duration,
			AVG(buffer_status_buffer_health) as avg_buffer_health
		`).
		Where("platform_context = ? AND session_start > ?", platformType, since).
		Scan(&sessionStats).Error

	if err != nil {
		return nil, err
	}

	analytics.TotalSessions = sessionStats.TotalSessions
	analytics.ActiveSessions = sessionStats.ActiveSessions
	analytics.AvgSessionDuration = sessionStats.AvgSessionDuration
	analytics.AvgBufferHealth = sessionStats.AvgBufferHealth

	return &analytics, nil
}

func (r *VideoRepository) GetPopularVideos(limit int, timeframe time.Duration) ([]*models.VideoContent, error) {
	since := time.Now().Add(-timeframe)

	var videoIDs []uuid.UUID
	if err := r.db.Model(&models.ViewingHistory{}).
		Select("video_id").
		Where("watched_at > ?", since).
		Group("video_id").
		Order("COUNT(*) DESC, AVG(engagement_score) DESC").
		Limit(limit).
		Pluck("video_id", &videoIDs).Error; err != nil {
		return nil, err
	}

	var videos []*models.VideoContent
	if err := r.db.Where("id IN ?", videoIDs).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func (r *VideoRepository) GetTrendingVideos(limit int) ([]*models.VideoContent, error) {
	// Videos with high recent engagement and growth rate
	since := time.Now().Add(-24 * time.Hour) // Last 24 hours

	var videoIDs []uuid.UUID
	if err := r.db.Model(&models.ViewingHistory{}).
		Select("video_id").
		Where("watched_at > ?", since).
		Group("video_id").
		Having("COUNT(*) > ? AND AVG(engagement_score) > ?", 10, 0.7). // Minimum views and engagement
		Order("COUNT(*) DESC, AVG(completion_rate) DESC").
		Limit(limit).
		Pluck("video_id", &videoIDs).Error; err != nil {
		return nil, err
	}

	var videos []*models.VideoContent
	if err := r.db.Where("id IN ?", videoIDs).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// Search and discovery

func (r *VideoRepository) SearchVideos(query string, filters VideoSearchFilters, limit, offset int) ([]*models.VideoContent, error) {
	db := r.db.Model(&models.VideoContent{})

	// Text search
	if query != "" {
		db = db.Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// Apply filters
	if len(filters.Tags) > 0 {
		db = db.Where("tags @> ?", filters.Tags)
	}

	if filters.CreatorID != uuid.Nil {
		db = db.Where("creator_id = ?", filters.CreatorID)
	}

	if filters.Status != "" {
		db = db.Where("availability_status = ?", filters.Status)
	}

	if filters.ContentRating != "" {
		db = db.Where("content_rating = ?", filters.ContentRating)
	}

	if filters.MinDuration > 0 {
		db = db.Where("duration >= ?", filters.MinDuration)
	}

	if filters.MaxDuration > 0 {
		db = db.Where("duration <= ?", filters.MaxDuration)
	}

	if !filters.UploadedAfter.IsZero() {
		db = db.Where("upload_timestamp >= ?", filters.UploadedAfter)
	}

	if !filters.UploadedBefore.IsZero() {
		db = db.Where("upload_timestamp <= ?", filters.UploadedBefore)
	}

	var videos []*models.VideoContent
	err := db.Order("upload_timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error

	return videos, err
}

func (r *VideoRepository) GetRecommendations(userID uuid.UUID, limit int) ([]*models.VideoContent, error) {
	// Simple recommendation based on user's viewing history and similar users
	// In production, this would use machine learning algorithms

	// Get user's preferred tags and creators
	var userTags []string
	r.db.Model(&models.ViewingHistory{}).
		Joins("JOIN video_contents ON viewing_histories.video_id = video_contents.id").
		Where("viewing_histories.user_id = ? AND viewing_histories.completion_rate > ?", userID, 70).
		Pluck("DISTINCT UNNEST(video_contents.tags)", &userTags)

	var videos []*models.VideoContent
	if len(userTags) > 0 {
		err := r.db.Where("tags && ? AND id NOT IN (SELECT video_id FROM viewing_histories WHERE user_id = ?)",
			userTags, userID).
			Order("social_metrics_view_count DESC, upload_timestamp DESC").
			Limit(limit).
			Find(&videos).Error
		return videos, err
	}

	// Fallback to popular videos
	return r.GetPopularVideos(limit, 7*24*time.Hour)
}

func (r *VideoRepository) GetSimilarVideos(videoID uuid.UUID, limit int) ([]*models.VideoContent, error) {
	// Get the source video's tags and creator
	var sourceVideo models.VideoContent
	if err := r.db.Select("tags, creator_id").First(&sourceVideo, "id = ?", videoID).Error; err != nil {
		return nil, err
	}

	var videos []*models.VideoContent
	err := r.db.Where("(tags && ? OR creator_id = ?) AND id != ?",
		sourceVideo.Tags, sourceVideo.CreatorID, videoID).
		Order("social_metrics_view_count DESC").
		Limit(limit).
		Find(&videos).Error

	return videos, err
}

// Health and maintenance

func (r *VideoRepository) GetRepositoryHealth() (*RepositoryHealth, error) {
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

	// Check table counts
	var counts struct {
		Videos      int64 `json:"videos"`
		Sessions    int64 `json:"sessions"`
		History     int64 `json:"history"`
		Configs     int64 `json:"configs"`
		SyncStates  int64 `json:"sync_states"`
	}

	r.db.Model(&models.VideoContent{}).Count(&counts.Videos)
	r.db.Model(&models.PlaybackSession{}).Count(&counts.Sessions)
	r.db.Model(&models.ViewingHistory{}).Count(&counts.History)
	r.db.Model(&models.PlatformConfiguration{}).Count(&counts.Configs)
	r.db.Model(&models.SynchronizationState{}).Count(&counts.SyncStates)

	health.TableCounts = map[string]int64{
		"videos":      counts.Videos,
		"sessions":    counts.Sessions,
		"history":     counts.History,
		"configs":     counts.Configs,
		"sync_states": counts.SyncStates,
	}

	health.Healthy = health.DatabaseConnected && len(health.Errors) == 0
	health.CheckedAt = time.Now()

	return &health, nil
}

func (r *VideoRepository) OptimizeDatabase() error {
	// Analyze tables for better query performance
	tables := []string{
		"video_contents",
		"playback_sessions",
		"viewing_histories",
		"platform_configurations",
		"synchronization_states",
	}

	for _, table := range tables {
		if err := r.db.Exec(fmt.Sprintf("ANALYZE %s", table)).Error; err != nil {
			return fmt.Errorf("failed to analyze table %s: %w", table, err)
		}
	}

	return nil
}

func (r *VideoRepository) GetDatabaseStats() (*DatabaseStats, error) {
	var stats DatabaseStats

	// Get table sizes
	query := `
		SELECT
			schemaname,
			tablename,
			attname,
			n_distinct,
			correlation
		FROM pg_stats
		WHERE schemaname = 'public'
		AND tablename IN ('video_contents', 'playback_sessions', 'viewing_histories', 'platform_configurations', 'synchronization_states')
	`

	rows, err := r.db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.TableStats = make(map[string]interface{})
	for rows.Next() {
		var schema, table, column string
		var nDistinct, correlation float64

		if err := rows.Scan(&schema, &table, &column, &nDistinct, &correlation); err != nil {
			continue
		}

		if stats.TableStats[table] == nil {
			stats.TableStats[table] = make(map[string]interface{})
		}

		tableStats := stats.TableStats[table].(map[string]interface{})
		tableStats[column] = map[string]float64{
			"n_distinct":  nDistinct,
			"correlation": correlation,
		}
	}

	stats.GeneratedAt = time.Now()
	return &stats, nil
}

// Supporting types for analytics and search

type VideoAnalytics struct {
	Video               models.VideoContent `json:"video"`
	TotalViews          int64               `json:"total_views"`
	UniqueViewers       int64               `json:"unique_viewers"`
	AvgCompletionRate   float64             `json:"avg_completion_rate"`
	AvgEngagementScore  float64             `json:"avg_engagement_score"`
	TotalWatchTime      int64               `json:"total_watch_time"`
	PopularSegments     []TimeSegment       `json:"popular_segments"`
	PlatformBreakdown   map[string]int64    `json:"platform_breakdown"`
}

type UserAnalytics struct {
	UserID              uuid.UUID           `json:"user_id"`
	TotalVideosWatched  int64               `json:"total_videos_watched"`
	TotalWatchTime      int64               `json:"total_watch_time"`
	AvgSessionDuration  float64             `json:"avg_session_duration"`
	PreferredQuality    string              `json:"preferred_quality"`
	PreferredPlatform   string              `json:"preferred_platform"`
	ViewingPatterns     map[string]int64    `json:"viewing_patterns"`
	TopCategories       []string            `json:"top_categories"`
}

type PlatformAnalytics struct {
	Platform            string              `json:"platform"`
	Since               time.Time           `json:"since"`
	TotalSessions       int64               `json:"total_sessions"`
	ActiveSessions      int64               `json:"active_sessions"`
	AvgSessionDuration  float64             `json:"avg_session_duration"`
	AvgBufferHealth     float64             `json:"avg_buffer_health"`
	QualityDistribution map[string]int64    `json:"quality_distribution"`
	DeviceBreakdown     map[string]int64    `json:"device_breakdown"`
}

type VideoSearchFilters struct {
	Tags            []string                     `json:"tags"`
	CreatorID       uuid.UUID                    `json:"creator_id"`
	Status          models.AvailabilityStatus    `json:"status"`
	ContentRating   models.ContentRating         `json:"content_rating"`
	MinDuration     int                          `json:"min_duration"`
	MaxDuration     int                          `json:"max_duration"`
	UploadedAfter   time.Time                    `json:"uploaded_after"`
	UploadedBefore  time.Time                    `json:"uploaded_before"`
}

type TimeSegment struct {
	Start     int    `json:"start"`
	End       int    `json:"end"`
	ViewCount int64  `json:"view_count"`
}

type RepositoryHealth struct {
	Healthy             bool                `json:"healthy"`
	DatabaseConnected   bool                `json:"database_connected"`
	TableCounts         map[string]int64    `json:"table_counts"`
	Errors              []string            `json:"errors"`
	CheckedAt           time.Time           `json:"checked_at"`
}

type DatabaseStats struct {
	TableStats          map[string]interface{} `json:"table_stats"`
	IndexStats          map[string]interface{} `json:"index_stats"`
	GeneratedAt         time.Time              `json:"generated_at"`
}