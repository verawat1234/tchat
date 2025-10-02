// backend/video/services/video_service.go
// Video Service - Business logic for video upload, processing, and streaming
// Implements T031: VideoService with upload/streaming logic

package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"tchat.dev/video/models"
	"tchat.dev/video/repository"
)

// VideoServiceInterface defines the contract for video business logic
type VideoServiceInterface interface {
	// Upload and processing
	UploadVideo(ctx context.Context, file multipart.File, header *multipart.FileHeader, metadata VideoMetadata) (*models.VideoContent, error)
	ProcessVideoUpload(ctx context.Context, videoID uuid.UUID) error
	GenerateThumbnails(ctx context.Context, videoID uuid.UUID) error
	TranscodeVideo(ctx context.Context, videoID uuid.UUID, qualities []string) error

	// Video management
	GetVideo(ctx context.Context, videoID uuid.UUID) (*models.VideoContent, error)
	GetVideosByCreator(ctx context.Context, creatorID uuid.UUID, pagination Pagination) ([]*models.VideoContent, error)
	UpdateVideoMetadata(ctx context.Context, videoID uuid.UUID, metadata VideoMetadata) error
	DeleteVideo(ctx context.Context, videoID uuid.UUID) error
	UpdateVideoStatus(ctx context.Context, videoID uuid.UUID, status models.AvailabilityStatus) error

	// Streaming
	GetStreamingURL(ctx context.Context, videoID uuid.UUID, userID uuid.UUID, quality string, platform models.PlatformType) (*StreamingURLResponse, error)
	GetAdaptiveStreamManifest(ctx context.Context, videoID uuid.UUID, userID uuid.UUID) (*StreamManifest, error)
	ValidateStreamAccess(ctx context.Context, videoID uuid.UUID, userID uuid.UUID) error

	// Playback sessions
	CreatePlaybackSession(ctx context.Context, req CreateSessionRequest) (*models.PlaybackSession, error)
	UpdatePlaybackProgress(ctx context.Context, sessionID uuid.UUID, position int, metrics PlaybackMetrics) error
	EndPlaybackSession(ctx context.Context, sessionID uuid.UUID) error
	GetActiveUserSessions(ctx context.Context, userID uuid.UUID) ([]*models.PlaybackSession, error)

	// Analytics
	GetVideoAnalytics(ctx context.Context, videoID uuid.UUID) (*repository.VideoAnalytics, error)
	GetUserViewingHistory(ctx context.Context, userID uuid.UUID, pagination Pagination) ([]*models.ViewingHistory, error)
	RecordVideoView(ctx context.Context, view VideoViewEvent) error

	// Search and discovery
	SearchVideos(ctx context.Context, query string, filters repository.VideoSearchFilters, pagination Pagination) (*SearchResults, error)
	GetRecommendations(ctx context.Context, userID uuid.UUID, limit int) ([]*models.VideoContent, error)
	GetTrendingVideos(ctx context.Context, limit int) ([]*models.VideoContent, error)
	GetPopularVideos(ctx context.Context, limit int, timeframe time.Duration) ([]*models.VideoContent, error)

	// Content moderation
	FlagVideoForReview(ctx context.Context, videoID uuid.UUID, reason string, reportedBy uuid.UUID) error
	ModerateVideo(ctx context.Context, videoID uuid.UUID, decision ModerationDecision) error

	// Health and metrics
	GetServiceHealth(ctx context.Context) (*ServiceHealth, error)
	GetServiceMetrics(ctx context.Context) (*ServiceMetrics, error)
}

// VideoService implements VideoServiceInterface
type VideoService struct {
	repo            repository.VideoRepositoryInterface
	storageService  StorageServiceInterface
	transcoder      TranscoderInterface
	cdnService      CDNServiceInterface
	cacheService    CacheServiceInterface
	notificationSvc NotificationServiceInterface
}

// NewVideoService creates a new video service instance
func NewVideoService(
	repo repository.VideoRepositoryInterface,
	storage StorageServiceInterface,
	transcoder TranscoderInterface,
	cdn CDNServiceInterface,
	cache CacheServiceInterface,
	notification NotificationServiceInterface,
) VideoServiceInterface {
	return &VideoService{
		repo:            repo,
		storageService:  storage,
		transcoder:      transcoder,
		cdnService:      cdn,
		cacheService:    cache,
		notificationSvc: notification,
	}
}

// Upload and processing

func (s *VideoService) UploadVideo(ctx context.Context, file multipart.File, header *multipart.FileHeader, metadata VideoMetadata) (*models.VideoContent, error) {
	// Validate file
	if err := s.validateVideoFile(header); err != nil {
		return nil, fmt.Errorf("invalid video file: %w", err)
	}

	// Create video record
	video := &models.VideoContent{
		ID:                  uuid.New(),
		Title:               metadata.Title,
		Description:         metadata.Description,
		CreatorID:           metadata.CreatorID,
		Duration:            0, // Will be set after processing
		ProcessingStatus:    models.ProcessingActive,
		ContentRating:       metadata.ContentRating,
		Tags:                metadata.Tags,
		UploadTimestamp:     time.Now(),
	}

	// Upload to storage
	storageURL, err := s.storageService.UploadVideo(ctx, file, video.ID.String(), header.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to upload video: %w", err)
	}

	video.FileURL = storageURL

	// Save to database
	if err := s.repo.CreateVideo(video); err != nil {
		// Cleanup uploaded file
		s.storageService.DeleteVideo(ctx, video.ID.String())
		return nil, fmt.Errorf("failed to save video record: %w", err)
	}

	// Trigger async processing
	go s.ProcessVideoUpload(context.Background(), video.ID)

	// Send notification
	s.notificationSvc.NotifyVideoUploadStarted(video.CreatorID, video.ID)

	return video, nil
}

func (s *VideoService) ProcessVideoUpload(ctx context.Context, videoID uuid.UUID) error {
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	// Extract video metadata (duration, resolution, etc.)
	metadata, err := s.transcoder.ExtractMetadata(ctx, video.FileURL)
	if err != nil {
		video.ProcessingStatus = models.ProcessingFailed
		s.repo.UpdateVideo(video)
		return fmt.Errorf("failed to extract metadata: %w", err)
	}

	video.Duration = metadata.Duration
	video.FormatSpecification.FPS = metadata.FrameRate
	video.FormatSpecification.Bitrate = metadata.Bitrate
	// Note: Width, Height, Codec would need to be extracted from Resolution and stored separately

	// Generate thumbnails
	if err := s.GenerateThumbnails(ctx, videoID); err != nil {
		// Non-fatal error, log and continue
		fmt.Printf("Warning: Failed to generate thumbnails for video %s: %v\n", videoID, err)
	}

	// Transcode to multiple qualities
	qualities := []string{"360p", "480p", "720p", "1080p"}
	if err := s.TranscodeVideo(ctx, videoID, qualities); err != nil {
		video.ProcessingStatus = models.ProcessingFailed
		s.repo.UpdateVideo(video)
		return fmt.Errorf("failed to transcode video: %w", err)
	}

	// Update status to published and mark processing complete
	video.AvailabilityStatus = models.StatusPublic
	video.ProcessingStatus = models.ProcessingCompleted

	if err := s.repo.UpdateVideo(video); err != nil {
		return err
	}

	// Notify creator
	s.notificationSvc.NotifyVideoProcessingComplete(video.CreatorID, videoID)

	return nil
}

func (s *VideoService) GenerateThumbnails(ctx context.Context, videoID uuid.UUID) error {
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	// Generate thumbnails at different timestamps
	timestamps := []int{0, video.Duration / 4, video.Duration / 2, 3 * video.Duration / 4}
	thumbnails := make([]string, 0, len(timestamps))

	for _, timestamp := range timestamps {
		thumbnailURL, err := s.transcoder.GenerateThumbnail(ctx, video.FileURL, timestamp)
		if err != nil {
			return fmt.Errorf("failed to generate thumbnail at %ds: %w", timestamp, err)
		}
		thumbnails = append(thumbnails, thumbnailURL)
	}

	video.ThumbnailURL = thumbnails[0] // Use first thumbnail as primary
	// Note: PreviewImages would need to be stored in a separate model or as JSONB field

	return s.repo.UpdateVideo(video)
}

func (s *VideoService) TranscodeVideo(ctx context.Context, videoID uuid.UUID, qualities []string) error {
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	qualityOptions := make([]string, 0, len(qualities))
	// Note: StreamingURLs would need to be stored in a separate model or as JSONB field

	for _, quality := range qualities {
		// Transcode video to specific quality
		_, err := s.transcoder.TranscodeVideo(ctx, video.FileURL, quality)
		if err != nil {
			fmt.Printf("Warning: Failed to transcode video %s to %s: %v\n", videoID, quality, err)
			continue
		}

		qualityOptions = append(qualityOptions, quality)
		// transcodedURL would be stored in a StreamingURLs table/JSONB
	}

	if len(qualityOptions) == 0 {
		return fmt.Errorf("failed to transcode video to any quality")
	}

	video.QualityOptions = qualityOptions
	// StreamingURLs would be managed separately - not in VideoContent model

	return s.repo.UpdateVideo(video)
}

// Video management

func (s *VideoService) GetVideo(ctx context.Context, videoID uuid.UUID) (*models.VideoContent, error) {
	// Try cache first
	if cachedVideo := s.cacheService.GetVideo(videoID); cachedVideo != nil {
		return cachedVideo, nil
	}

	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return nil, err
	}

	// Cache for future requests
	s.cacheService.SetVideo(videoID, video, 5*time.Minute)

	return video, nil
}

func (s *VideoService) GetVideosByCreator(ctx context.Context, creatorID uuid.UUID, pagination Pagination) ([]*models.VideoContent, error) {
	return s.repo.GetVideosByCreator(creatorID, pagination.Limit, pagination.Offset)
}

func (s *VideoService) UpdateVideoMetadata(ctx context.Context, videoID uuid.UUID, metadata VideoMetadata) error {
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	// Update allowed fields
	if metadata.Title != "" {
		video.Title = metadata.Title
	}
	if metadata.Description != "" {
		video.Description = metadata.Description
	}
	if len(metadata.Tags) > 0 {
		video.Tags = metadata.Tags
	}
	if metadata.ContentRating != "" {
		video.ContentRating = metadata.ContentRating
	}

	// Invalidate cache
	s.cacheService.DeleteVideo(videoID)

	return s.repo.UpdateVideo(video)
}

func (s *VideoService) DeleteVideo(ctx context.Context, videoID uuid.UUID) error {
	_, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	// Delete from storage
	if err := s.storageService.DeleteVideo(ctx, videoID.String()); err != nil {
		fmt.Printf("Warning: Failed to delete video files for %s: %v\n", videoID, err)
	}

	// Delete from CDN cache
	if err := s.cdnService.PurgeVideo(ctx, videoID.String()); err != nil {
		fmt.Printf("Warning: Failed to purge CDN cache for %s: %v\n", videoID, err)
	}

	// Delete from cache
	s.cacheService.DeleteVideo(videoID)

	// Soft delete from database
	return s.repo.DeleteVideo(videoID)
}

func (s *VideoService) UpdateVideoStatus(ctx context.Context, videoID uuid.UUID, status models.AvailabilityStatus) error {
	video, err := s.repo.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	// Validate status transition
	if err := video.ValidateAvailabilityTransition(status); err != nil {
		return err
	}

	video.AvailabilityStatus = status

	// Note: PublishTimestamp tracking would require adding this field to VideoContent model
	// if status == models.StatusPublic { ... }

	// Invalidate cache
	s.cacheService.DeleteVideo(videoID)

	return s.repo.UpdateVideo(video)
}

// Streaming

func (s *VideoService) GetStreamingURL(ctx context.Context, videoID uuid.UUID, userID uuid.UUID, quality string, platform models.PlatformType) (*StreamingURLResponse, error) {
	// Validate access
	if err := s.ValidateStreamAccess(ctx, videoID, userID); err != nil {
		return nil, err
	}

	video, err := s.GetVideo(ctx, videoID)
	if err != nil {
		return nil, err
	}

	// Get platform configuration for optimal quality
	platformConfig, err := s.repo.GetPlatformConfig(platform, "mobile") // TODO: Determine device category
	if err != nil {
		platformConfig = nil // Use defaults
	}

	// Select best quality based on requested quality and platform capabilities
	selectedQuality := quality
	if quality == "auto" || quality == "" {
		selectedQuality = s.selectOptimalQuality(video, platformConfig)
	}

	// Note: StreamingURLs would need to be stored in a separate model or retrieved from CDN
	// Using FileURL as fallback streaming URL
	streamingURL := video.FileURL
	if len(video.QualityOptions) > 0 {
		selectedQuality = video.QualityOptions[0]
	}

	// Generate signed URL with expiration
	signedURL, err := s.cdnService.GenerateSignedURL(ctx, streamingURL, 3*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return &StreamingURLResponse{
		URL:              signedURL,
		Quality:          selectedQuality,
		Format:           "mp4", // Default format since OriginalFormat field doesn't exist
		Duration:         video.Duration,
		ExpiresAt:        time.Now().Add(3 * time.Hour),
		SupportedQualities: video.QualityOptions,
	}, nil
}

func (s *VideoService) GetAdaptiveStreamManifest(ctx context.Context, videoID uuid.UUID, userID uuid.UUID) (*StreamManifest, error) {
	// Validate access
	if err := s.ValidateStreamAccess(ctx, videoID, userID); err != nil {
		return nil, err
	}

	video, err := s.GetVideo(ctx, videoID)
	if err != nil {
		return nil, err
	}

	// Note: StreamingURLs would need to be stored in a separate model or retrieved from CDN
	// Using FileURL as fallback for HLS manifest
	hlsManifestURL := video.FileURL

	// Generate signed manifest URL
	signedManifestURL, err := s.cdnService.GenerateSignedURL(ctx, hlsManifestURL, 3*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signed manifest URL: %w", err)
	}

	return &StreamManifest{
		Type:             "hls",
		ManifestURL:      signedManifestURL,
		AvailableQualities: video.QualityOptions,
		Duration:         video.Duration,
		ExpiresAt:        time.Now().Add(3 * time.Hour),
	}, nil
}

func (s *VideoService) ValidateStreamAccess(ctx context.Context, videoID uuid.UUID, userID uuid.UUID) error {
	video, err := s.GetVideo(ctx, videoID)
	if err != nil {
		return err
	}

	// Check if video is accessible (IsAccessible takes no parameters)
	if !video.IsAccessible() {
		return fmt.Errorf("access denied to video %s", videoID)
	}

	// Check subscription/payment requirements
	if video.IsMonetized() && video.GetPrice("USD") > 0 {
		// TODO: Check if user has purchased or subscribed
		// For now, allow access
	}

	return nil
}

// Playback sessions

func (s *VideoService) CreatePlaybackSession(ctx context.Context, req CreateSessionRequest) (*models.PlaybackSession, error) {
	// Validate video access
	if err := s.ValidateStreamAccess(ctx, req.VideoID, req.UserID); err != nil {
		return nil, err
	}

	// Note: req.DeviceInfo is a string but PlaybackSession expects DeviceInformation struct
	// Parse or convert as needed
	session := &models.PlaybackSession{
		ID:              uuid.New(),
		UserID:          req.UserID,
		VideoID:         req.VideoID,
		PlatformContext: string(req.Platform),
		DeviceInfo:      models.DeviceInformation{UserAgent: req.DeviceInfo}, // Convert string to struct
		SessionStart:    time.Now(),
		IsActive:        true,
		SessionState:    models.SessionActive,
		CurrentPosition: req.StartPosition,
		PlaybackSpeed:   1.0,
	}

	if err := s.repo.CreatePlaybackSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *VideoService) UpdatePlaybackProgress(ctx context.Context, sessionID uuid.UUID, position int, metrics PlaybackMetrics) error {
	session, err := s.repo.GetPlaybackSession(sessionID)
	if err != nil {
		return err
	}

	// Update session
	session.CurrentPosition = position
	session.BufferStatus.BufferHealth = metrics.BufferHealth
	// Note: BufferStatus doesn't have DroppedFrames field
	// Note: PlaybackSession doesn't have PerformanceMetrics field - metrics tracked differently

	if err := s.repo.UpdatePlaybackSession(session); err != nil {
		return err
	}

	return nil
}

func (s *VideoService) EndPlaybackSession(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.repo.GetPlaybackSession(sessionID)
	if err != nil {
		return err
	}

	// Complete session
	if err := s.repo.CompletePlaybackSession(sessionID); err != nil {
		return err
	}

	// Create viewing history record
	completionRate := float64(session.CurrentPosition) / float64(session.VideoContent.Duration) * 100
	if completionRate > 100 {
		completionRate = 100
	}

	viewingHistory := &models.ViewingHistory{
		UserID:          session.UserID,
		VideoID:         session.VideoID,
		WatchedAt:       session.SessionStart,
		WatchedDuration: session.CurrentPosition,
		CompletionRate:  completionRate,
		QualityWatched:  session.QualitySetting, // Use QualitySetting instead of QualityLevel
		PlatformUsed:    session.PlatformContext,
	}

	return s.repo.CreateViewingHistory(viewingHistory)
}

func (s *VideoService) GetActiveUserSessions(ctx context.Context, userID uuid.UUID) ([]*models.PlaybackSession, error) {
	return s.repo.GetActiveSessionsByUser(userID)
}

// Analytics

func (s *VideoService) GetVideoAnalytics(ctx context.Context, videoID uuid.UUID) (*repository.VideoAnalytics, error) {
	return s.repo.GetVideoAnalytics(videoID)
}

func (s *VideoService) GetUserViewingHistory(ctx context.Context, userID uuid.UUID, pagination Pagination) ([]*models.ViewingHistory, error) {
	return s.repo.GetViewingHistoryByUser(userID, pagination.Limit, pagination.Offset)
}

func (s *VideoService) RecordVideoView(ctx context.Context, view VideoViewEvent) error {
	// Increment view count
	video, err := s.GetVideo(ctx, view.VideoID)
	if err != nil {
		return err
	}

	video.SocialMetrics.ViewCount++
	// Note: SocialEngagementMetrics doesn't have LastViewedAt field

	// Invalidate cache
	s.cacheService.DeleteVideo(view.VideoID)

	return s.repo.UpdateVideo(video)
}

// Search and discovery

func (s *VideoService) SearchVideos(ctx context.Context, query string, filters repository.VideoSearchFilters, pagination Pagination) (*SearchResults, error) {
	videos, err := s.repo.SearchVideos(query, filters, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	// TODO: Implement count query in repository
	totalCount := int64(len(videos))

	return &SearchResults{
		Videos:     videos,
		Total:      totalCount,
		Page:       pagination.Offset/pagination.Limit + 1,
		PageSize:   pagination.Limit,
		HasMore:    totalCount > int64(pagination.Offset+pagination.Limit),
	}, nil
}

func (s *VideoService) GetRecommendations(ctx context.Context, userID uuid.UUID, limit int) ([]*models.VideoContent, error) {
	return s.repo.GetRecommendations(userID, limit)
}

func (s *VideoService) GetTrendingVideos(ctx context.Context, limit int) ([]*models.VideoContent, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("trending:videos:%d", limit)
	if cachedVideos := s.cacheService.Get(cacheKey); cachedVideos != nil {
		if videos, ok := cachedVideos.([]*models.VideoContent); ok {
			return videos, nil
		}
	}

	videos, err := s.repo.GetTrendingVideos(limit)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	s.cacheService.Set(cacheKey, videos, 5*time.Minute)

	return videos, nil
}

func (s *VideoService) GetPopularVideos(ctx context.Context, limit int, timeframe time.Duration) ([]*models.VideoContent, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("popular:videos:%d:%s", limit, timeframe)
	if cachedVideos := s.cacheService.Get(cacheKey); cachedVideos != nil {
		if videos, ok := cachedVideos.([]*models.VideoContent); ok {
			return videos, nil
		}
	}

	videos, err := s.repo.GetPopularVideos(limit, timeframe)
	if err != nil {
		return nil, err
	}

	// Cache for 10 minutes
	s.cacheService.Set(cacheKey, videos, 10*time.Minute)

	return videos, nil
}

// Content moderation

func (s *VideoService) FlagVideoForReview(ctx context.Context, videoID uuid.UUID, reason string, reportedBy uuid.UUID) error {
	video, err := s.GetVideo(ctx, videoID)
	if err != nil {
		return err
	}

	// Update status to private for review (StatusUnderReview doesn't exist in model)
	video.AvailabilityStatus = models.StatusPrivate

	// TODO: Create moderation record
	// For now, just update video status
	s.cacheService.DeleteVideo(videoID)

	if err := s.repo.UpdateVideo(video); err != nil {
		return err
	}

	// Notify moderation team
	s.notificationSvc.NotifyModerationRequired(videoID, reason, reportedBy)

	return nil
}

func (s *VideoService) ModerateVideo(ctx context.Context, videoID uuid.UUID, decision ModerationDecision) error {
	video, err := s.GetVideo(ctx, videoID)
	if err != nil {
		return err
	}

	switch decision.Action {
	case "approve":
		video.AvailabilityStatus = models.StatusPublic
	case "reject":
		video.AvailabilityStatus = models.StatusPrivate // Use StatusPrivate since StatusRestricted doesn't exist
	case "remove":
		return s.DeleteVideo(ctx, videoID)
	default:
		return fmt.Errorf("invalid moderation action: %s", decision.Action)
	}

	// Invalidate cache
	s.cacheService.DeleteVideo(videoID)

	if err := s.repo.UpdateVideo(video); err != nil {
		return err
	}

	// Notify creator
	s.notificationSvc.NotifyModerationDecision(video.CreatorID, videoID, decision.Action, decision.Reason)

	return nil
}

// Health and metrics

func (s *VideoService) GetServiceHealth(ctx context.Context) (*ServiceHealth, error) {
	repoHealth, err := s.repo.GetRepositoryHealth()
	if err != nil {
		return &ServiceHealth{
			Status:  "unhealthy",
			Message: fmt.Sprintf("repository health check failed: %v", err),
		}, nil
	}

	return &ServiceHealth{
		Status:            "healthy",
		RepositoryHealthy: repoHealth.Healthy,
		StorageHealthy:    s.storageService.IsHealthy(),
		CDNHealthy:        s.cdnService.IsHealthy(),
		CacheHealthy:      s.cacheService.IsHealthy(),
		CheckedAt:         time.Now(),
	}, nil
}

func (s *VideoService) GetServiceMetrics(ctx context.Context) (*ServiceMetrics, error) {
	dbStats, err := s.repo.GetDatabaseStats()
	if err != nil {
		return nil, err
	}

	return &ServiceMetrics{
		DatabaseStats:    dbStats,
		ActiveSessions:   0, // TODO: Get from session tracking
		TotalVideos:      0, // TODO: Get from repository
		TotalViews:       0, // TODO: Get from analytics
		StorageUsed:      0, // TODO: Get from storage service
		BandwidthUsed:    0, // TODO: Get from CDN
		GeneratedAt:      time.Now(),
	}, nil
}

// Helper methods

func (s *VideoService) validateVideoFile(header *multipart.FileHeader) error {
	// Check file size (max 5GB)
	maxSize := int64(5 * 1024 * 1024 * 1024) // 5GB
	if header.Size > maxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", header.Size, maxSize)
	}

	// Check file extension
	ext := filepath.Ext(header.Filename)
	allowedExtensions := []string{".mp4", ".mov", ".avi", ".mkv", ".webm"}
	valid := false
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	return nil
}

func (s *VideoService) selectOptimalQuality(video *models.VideoContent, config *models.PlatformConfiguration) string {
	if config == nil || len(video.QualityOptions) == 0 {
		// Default to 720p if available
		for _, quality := range video.QualityOptions {
			if quality == "720p" {
				return quality
			}
		}
		return video.QualityOptions[0]
	}

	// Use platform configuration to select optimal quality
	return config.GetOptimalQuality(1.0, config.VideoCapabilities.MaxResolution)
}

// Supporting types

type VideoMetadata struct {
	Title         string                    `json:"title" validate:"required,min=1,max=200"`
	Description   string                    `json:"description" validate:"max=5000"`
	CreatorID     uuid.UUID                 `json:"creator_id" validate:"required"`
	ContentRating models.ContentRating      `json:"content_rating"`
	Tags          []string                  `json:"tags" validate:"max=20"`
}

type Pagination struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

type StreamingURLResponse struct {
	URL                string    `json:"url"`
	Quality            string    `json:"quality"`
	Format             string    `json:"format"`
	Duration           int       `json:"duration"`
	ExpiresAt          time.Time `json:"expires_at"`
	SupportedQualities []string  `json:"supported_qualities"`
}

type StreamManifest struct {
	Type               string    `json:"type"`
	ManifestURL        string    `json:"manifest_url"`
	AvailableQualities []string  `json:"available_qualities"`
	Duration           int       `json:"duration"`
	ExpiresAt          time.Time `json:"expires_at"`
}

type CreateSessionRequest struct {
	UserID        uuid.UUID            `json:"user_id" validate:"required"`
	VideoID       uuid.UUID            `json:"video_id" validate:"required"`
	Platform      models.PlatformType  `json:"platform" validate:"required"`
	DeviceInfo    string               `json:"device_info"`
	StartPosition int                  `json:"start_position"`
}

type PlaybackMetrics struct {
	BufferHealth   float64 `json:"buffer_health"`
	DroppedFrames  int     `json:"dropped_frames"`
	ActualBitrate  int     `json:"actual_bitrate"`
	LatencyMS      int     `json:"latency_ms"`
}

type VideoViewEvent struct {
	VideoID   uuid.UUID `json:"video_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ViewedAt  time.Time `json:"viewed_at"`
	Platform  string    `json:"platform"`
	Duration  int       `json:"duration"`
}

type SearchResults struct {
	Videos   []*models.VideoContent `json:"videos"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	HasMore  bool                   `json:"has_more"`
}

type ModerationDecision struct {
	Action       string    `json:"action" validate:"required,oneof=approve reject remove"`
	Reason       string    `json:"reason" validate:"required"`
	ModeratedBy  uuid.UUID `json:"moderated_by" validate:"required"`
	ModeratedAt  time.Time `json:"moderated_at"`
}

type ServiceHealth struct {
	Status            string    `json:"status"`
	Message           string    `json:"message,omitempty"`
	RepositoryHealthy bool      `json:"repository_healthy"`
	StorageHealthy    bool      `json:"storage_healthy"`
	CDNHealthy        bool      `json:"cdn_healthy"`
	CacheHealthy      bool      `json:"cache_healthy"`
	CheckedAt         time.Time `json:"checked_at"`
}

type ServiceMetrics struct {
	DatabaseStats    *repository.DatabaseStats `json:"database_stats"`
	ActiveSessions   int64                     `json:"active_sessions"`
	TotalVideos      int64                     `json:"total_videos"`
	TotalViews       int64                     `json:"total_views"`
	StorageUsed      int64                     `json:"storage_used"`
	BandwidthUsed    int64                     `json:"bandwidth_used"`
	GeneratedAt      time.Time                 `json:"generated_at"`
}

// External service interfaces (to be implemented separately)

type StorageServiceInterface interface {
	UploadVideo(ctx context.Context, file io.Reader, videoID string, filename string) (string, error)
	DeleteVideo(ctx context.Context, videoID string) error
	IsHealthy() bool
}

type TranscoderInterface interface {
	ExtractMetadata(ctx context.Context, videoURL string) (*VideoMetadataExtracted, error)
	TranscodeVideo(ctx context.Context, sourceURL string, quality string) (string, error)
	GenerateThumbnail(ctx context.Context, videoURL string, timestamp int) (string, error)
	GenerateHLSManifest(ctx context.Context, streamingURLs map[string]string) (string, error)
}

type CDNServiceInterface interface {
	GenerateSignedURL(ctx context.Context, url string, expiration time.Duration) (string, error)
	PurgeVideo(ctx context.Context, videoID string) error
	IsHealthy() bool
}

type CacheServiceInterface interface {
	GetVideo(videoID uuid.UUID) *models.VideoContent
	SetVideo(videoID uuid.UUID, video *models.VideoContent, ttl time.Duration)
	DeleteVideo(videoID uuid.UUID)
	Get(key string) interface{}
	Set(key string, value interface{}, ttl time.Duration)
	IsHealthy() bool
}

type NotificationServiceInterface interface {
	NotifyVideoUploadStarted(creatorID uuid.UUID, videoID uuid.UUID)
	NotifyVideoProcessingComplete(creatorID uuid.UUID, videoID uuid.UUID)
	NotifyModerationRequired(videoID uuid.UUID, reason string, reportedBy uuid.UUID)
	NotifyModerationDecision(creatorID uuid.UUID, videoID uuid.UUID, action string, reason string)
}

type VideoMetadataExtracted struct {
	Duration    int     `json:"duration"`
	AspectRatio string  `json:"aspect_ratio"`
	FrameRate   int     `json:"frame_rate"`
	Resolution  string  `json:"resolution"`
	Bitrate     int     `json:"bitrate"`
}