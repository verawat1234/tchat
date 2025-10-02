// backend/video/services/sync_service.go
// Synchronization Service - Cross-platform video playback synchronization
// Implements T032: SyncService for real-time cross-platform sync

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tchat.dev/video/models"
	"tchat.dev/video/repository"
)

// SyncServiceInterface defines the contract for synchronization logic
type SyncServiceInterface interface {
	// Synchronization state management
	InitializeSync(ctx context.Context, sessionID uuid.UUID, deviceID string, platform models.PlatformType) (*models.SynchronizationState, error)
	UpdateSyncState(ctx context.Context, syncID uuid.UUID, syncData models.SyncData) error
	GetSyncState(ctx context.Context, syncID uuid.UUID) (*models.SynchronizationState, error)
	GetSyncStatesBySession(ctx context.Context, sessionID uuid.UUID) ([]*models.SynchronizationState, error)

	// Cross-platform synchronization
	SyncPlaybackPosition(ctx context.Context, sessionID uuid.UUID, position int, platform models.PlatformType) error
	SyncPlaybackState(ctx context.Context, sessionID uuid.UUID, isPlaying bool, platform models.PlatformType) error
	SyncQualitySettings(ctx context.Context, sessionID uuid.UUID, quality string, platform models.PlatformType) error
	SyncUIState(ctx context.Context, sessionID uuid.UUID, uiState UIState, platform models.PlatformType) error

	// Conflict resolution
	DetectConflicts(ctx context.Context, sessionID uuid.UUID) ([]*models.SynchronizationState, error)
	ResolveConflict(ctx context.Context, syncID uuid.UUID, strategy models.ResolutionStrategy) error
	ResolveAllConflicts(ctx context.Context, sessionID uuid.UUID) error

	// Sync scheduling and execution
	ScheduleSync(ctx context.Context, syncID uuid.UUID, frequency int) error
	ExecutePendingSync(ctx context.Context) error
	ForceSyncAll(ctx context.Context, sessionID uuid.UUID) error

	// Platform-specific operations
	GetPlatformSyncState(ctx context.Context, deviceID string) (*models.SynchronizationState, error)
	SuspendPlatformSync(ctx context.Context, deviceID string) error
	ResumePlatformSync(ctx context.Context, deviceID string) error

	// Metrics and monitoring
	GetSyncMetrics(ctx context.Context, sessionID uuid.UUID) (*SyncMetrics, error)
	GetSyncHealth(ctx context.Context) (*SyncHealth, error)
	CleanupOldSyncStates(ctx context.Context, retentionDays int) (int64, error)

	// Bulk operations
	SyncSessionToAllPlatforms(ctx context.Context, sessionID uuid.UUID) error
	PropagateChanges(ctx context.Context, sessionID uuid.UUID, changes map[string]interface{}) error
}

// SyncService implements SyncServiceInterface
type SyncService struct {
	repo            repository.VideoRepositoryInterface
	cacheService    CacheServiceInterface
	messagingService MessagingServiceInterface
	configService   ConfigServiceInterface
}

// NewSyncService creates a new sync service instance
func NewSyncService(
	repo repository.VideoRepositoryInterface,
	cache CacheServiceInterface,
	messaging MessagingServiceInterface,
	config ConfigServiceInterface,
) SyncServiceInterface {
	return &SyncService{
		repo:            repo,
		cacheService:    cache,
		messagingService: messaging,
		configService:   config,
	}
}

// Synchronization state management

func (s *SyncService) InitializeSync(ctx context.Context, sessionID uuid.UUID, deviceID string, platform models.PlatformType) (*models.SynchronizationState, error) {
	// Get session details
	session, err := s.repo.GetPlaybackSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Create sync state
	syncState := &models.SynchronizationState{
		ID:           uuid.New(),
		SessionID:    sessionID,
		PlatformType: platform,
		DeviceID:     deviceID,
		SyncData: models.SyncData{
			Position:      session.CurrentPosition,
			PlaybackSpeed: session.PlaybackSpeed,
			IsPlaying:     session.IsActive,
			QualitySetting: session.QualitySetting,
		},
		SyncStatus:    models.SyncPending,
		SyncFrequency: 5, // Default 5 seconds
		MaxRetries:    3,
	}

	if err := s.repo.CreateSyncState(syncState); err != nil {
		return nil, fmt.Errorf("failed to create sync state: %w", err)
	}

	// Schedule initial sync
	s.ScheduleSync(ctx, syncState.ID, syncState.SyncFrequency)

	return syncState, nil
}

func (s *SyncService) UpdateSyncState(ctx context.Context, syncID uuid.UUID, syncData models.SyncData) error {
	syncState, err := s.repo.GetSyncState(syncID)
	if err != nil {
		return err
	}

	// Record metrics before update
	startTime := time.Now()

	// Update sync data
	syncState.UpdateSyncData(syncData)

	// Calculate sync latency
	latency := int(time.Since(startTime).Milliseconds())

	// Update metrics
	syncState.UpdateSyncMetrics(latency, true, 0)

	if err := s.repo.UpdateSyncState(syncState); err != nil {
		return err
	}

	// Broadcast changes to other platforms
	go s.broadcastSyncUpdate(ctx, syncState)

	return nil
}

func (s *SyncService) GetSyncState(ctx context.Context, syncID uuid.UUID) (*models.SynchronizationState, error) {
	return s.repo.GetSyncState(syncID)
}

func (s *SyncService) GetSyncStatesBySession(ctx context.Context, sessionID uuid.UUID) ([]*models.SynchronizationState, error) {
	return s.repo.GetSyncStateBySession(sessionID)
}

// Cross-platform synchronization

func (s *SyncService) SyncPlaybackPosition(ctx context.Context, sessionID uuid.UUID, position int, platform models.PlatformType) error {
	// Get all sync states for session
	syncStates, err := s.repo.GetSyncStateBySession(sessionID)
	if err != nil {
		return err
	}

	// Update position for all platforms except the source
	for _, syncState := range syncStates {
		if syncState.PlatformType == platform {
			continue // Skip source platform
		}

		syncState.SyncData.Position = position
		syncState.SyncData.LastWatchedTime = timePtr(time.Now())

		if err := s.repo.UpdateSyncState(syncState); err != nil {
			// Log error but continue with other platforms
			fmt.Printf("Warning: Failed to sync position for platform %s: %v\n", syncState.PlatformType, err)
			continue
		}

		// Send real-time update
		s.sendSyncMessage(ctx, syncState, "position_update")
	}

	return nil
}

func (s *SyncService) SyncPlaybackState(ctx context.Context, sessionID uuid.UUID, isPlaying bool, platform models.PlatformType) error {
	syncStates, err := s.repo.GetSyncStateBySession(sessionID)
	if err != nil {
		return err
	}

	for _, syncState := range syncStates {
		if syncState.PlatformType == platform {
			continue
		}

		syncState.SyncData.IsPlaying = isPlaying

		if err := s.repo.UpdateSyncState(syncState); err != nil {
			fmt.Printf("Warning: Failed to sync playback state for platform %s: %v\n", syncState.PlatformType, err)
			continue
		}

		s.sendSyncMessage(ctx, syncState, "playback_state_update")
	}

	return nil
}

func (s *SyncService) SyncQualitySettings(ctx context.Context, sessionID uuid.UUID, quality string, platform models.PlatformType) error {
	syncStates, err := s.repo.GetSyncStateBySession(sessionID)
	if err != nil {
		return err
	}

	for _, syncState := range syncStates {
		if syncState.PlatformType == platform {
			continue
		}

		// Check if platform supports the quality
		platformConfig, err := s.repo.GetPlatformConfig(syncState.PlatformType, "mobile")
		if err == nil && platformConfig.CanPlayQuality(quality) {
			syncState.SyncData.QualitySetting = quality

			if err := s.repo.UpdateSyncState(syncState); err != nil {
				fmt.Printf("Warning: Failed to sync quality for platform %s: %v\n", syncState.PlatformType, err)
				continue
			}

			s.sendSyncMessage(ctx, syncState, "quality_update")
		}
	}

	return nil
}

func (s *SyncService) SyncUIState(ctx context.Context, sessionID uuid.UUID, uiState UIState, platform models.PlatformType) error {
	syncStates, err := s.repo.GetSyncStateBySession(sessionID)
	if err != nil {
		return err
	}

	for _, syncState := range syncStates {
		if syncState.PlatformType == platform {
			continue
		}

		// Update UI state
		syncState.SyncData.FullscreenMode = uiState.FullscreenMode
		syncState.SyncData.ControlsVisible = uiState.ControlsVisible
		syncState.SyncData.SubtitlesEnabled = uiState.SubtitlesEnabled
		syncState.SyncData.SubtitleLanguage = uiState.SubtitleLanguage
		syncState.SyncData.Volume = uiState.Volume
		syncState.SyncData.IsMuted = uiState.IsMuted

		if err := s.repo.UpdateSyncState(syncState); err != nil {
			fmt.Printf("Warning: Failed to sync UI state for platform %s: %v\n", syncState.PlatformType, err)
			continue
		}

		s.sendSyncMessage(ctx, syncState, "ui_state_update")
	}

	return nil
}

// Conflict resolution

func (s *SyncService) DetectConflicts(ctx context.Context, sessionID uuid.UUID) ([]*models.SynchronizationState, error) {
	syncStates, err := s.repo.GetSyncStateBySession(sessionID)
	if err != nil {
		return nil, err
	}

	conflictedStates := make([]*models.SynchronizationState, 0)

	// Check for position conflicts (>5 seconds difference)
	if len(syncStates) > 1 {
		basePosition := syncStates[0].SyncData.Position

		for i := 1; i < len(syncStates); i++ {
			positionDiff := abs(syncStates[i].SyncData.Position - basePosition)

			if positionDiff > 5 {
				// Mark as conflicted
				syncStates[i].RecordConflict(models.ConflictPosition, map[string]interface{}{
					"base_position":    basePosition,
					"current_position": syncStates[i].SyncData.Position,
					"difference":       positionDiff,
				})

				s.repo.UpdateSyncState(syncStates[i])
				conflictedStates = append(conflictedStates, syncStates[i])
			}
		}
	}

	return conflictedStates, nil
}

func (s *SyncService) ResolveConflict(ctx context.Context, syncID uuid.UUID, strategy models.ResolutionStrategy) error {
	syncState, err := s.repo.GetSyncState(syncID)
	if err != nil {
		return err
	}

	if !syncState.ConflictState.HasConflicts {
		return fmt.Errorf("no conflicts to resolve for sync state %s", syncID)
	}

	// Resolve conflict using strategy
	if err := syncState.ResolveConflict(strategy); err != nil {
		return fmt.Errorf("failed to resolve conflict: %w", err)
	}

	// Update sync state
	if err := s.repo.UpdateSyncState(syncState); err != nil {
		return err
	}

	// Broadcast resolution
	s.sendSyncMessage(ctx, syncState, "conflict_resolved")

	return nil
}

func (s *SyncService) ResolveAllConflicts(ctx context.Context, sessionID uuid.UUID) error {
	conflictedStates, err := s.DetectConflicts(ctx, sessionID)
	if err != nil {
		return err
	}

	for _, syncState := range conflictedStates {
		// Use latest timestamp strategy by default
		if err := s.ResolveConflict(ctx, syncState.ID, models.ResolutionLatest); err != nil {
			fmt.Printf("Warning: Failed to resolve conflict for sync state %s: %v\n", syncState.ID, err)
		}
	}

	return nil
}

// Sync scheduling and execution

func (s *SyncService) ScheduleSync(ctx context.Context, syncID uuid.UUID, frequency int) error {
	syncState, err := s.repo.GetSyncState(syncID)
	if err != nil {
		return err
	}

	syncState.SyncFrequency = frequency
	syncState.ScheduleNextSync()

	return s.repo.UpdateSyncState(syncState)
}

func (s *SyncService) ExecutePendingSync(ctx context.Context) error {
	// Get all pending sync states
	pendingSyncStates, err := s.repo.GetPendingSyncStates()
	if err != nil {
		return err
	}

	successCount := 0
	failureCount := 0

	for _, syncState := range pendingSyncStates {
		// Check if it's time to sync
		if syncState.NextSyncTime != nil && time.Now().Before(*syncState.NextSyncTime) {
			continue
		}

		// Mark as in progress
		syncState.SyncStatus = models.SyncInProgress
		s.repo.UpdateSyncState(syncState)

		// Execute sync
		startTime := time.Now()

		// Get latest session data
		session, err := s.repo.GetPlaybackSession(syncState.SessionID)
		if err != nil {
			syncState.RecordError("sync_error", "SESSION_NOT_FOUND", err.Error(), models.SeverityHigh)
			syncState.SyncStatus = models.SyncFailed
			s.repo.UpdateSyncState(syncState)
			failureCount++
			continue
		}

		// Update sync data
		syncState.SyncData.Position = session.CurrentPosition
		syncState.SyncData.PlaybackSpeed = session.PlaybackSpeed
		syncState.SyncData.IsPlaying = session.IsActive
		syncState.SyncData.QualitySetting = session.QualitySetting

		// Calculate latency
		latency := int(time.Since(startTime).Milliseconds())

		// Update metrics
		syncState.UpdateSyncMetrics(latency, true, 0)

		// Mark as completed
		syncState.SyncStatus = models.SyncCompleted
		syncState.LastSyncTime = timePtr(time.Now())
		syncState.ScheduleNextSync()

		if err := s.repo.UpdateSyncState(syncState); err != nil {
			failureCount++
			continue
		}

		successCount++

		// Broadcast update
		s.sendSyncMessage(ctx, syncState, "sync_update")
	}

	fmt.Printf("Sync execution complete: %d succeeded, %d failed\n", successCount, failureCount)
	return nil
}

func (s *SyncService) ForceSyncAll(ctx context.Context, sessionID uuid.UUID) error {
	syncStates, err := s.repo.GetSyncStateBySession(sessionID)
	if err != nil {
		return err
	}

	// Get latest session data
	session, err := s.repo.GetPlaybackSession(sessionID)
	if err != nil {
		return err
	}

	for _, syncState := range syncStates {
		// Force immediate sync
		syncState.SyncData.Position = session.CurrentPosition
		syncState.SyncData.PlaybackSpeed = session.PlaybackSpeed
		syncState.SyncData.IsPlaying = session.IsActive
		syncState.SyncData.QualitySetting = session.QualitySetting
		syncState.LastSyncTime = timePtr(time.Now())

		if err := s.repo.UpdateSyncState(syncState); err != nil {
			fmt.Printf("Warning: Failed to force sync for platform %s: %v\n", syncState.PlatformType, err)
			continue
		}

		s.sendSyncMessage(ctx, syncState, "force_sync")
	}

	return nil
}

// Platform-specific operations

func (s *SyncService) GetPlatformSyncState(ctx context.Context, deviceID string) (*models.SynchronizationState, error) {
	syncStates, err := s.repo.GetSyncStateByDevice(deviceID)
	if err != nil {
		return nil, err
	}

	if len(syncStates) == 0 {
		return nil, fmt.Errorf("no sync state found for device %s", deviceID)
	}

	// Return the most recent sync state
	return syncStates[0], nil
}

func (s *SyncService) SuspendPlatformSync(ctx context.Context, deviceID string) error {
	syncStates, err := s.repo.GetSyncStateByDevice(deviceID)
	if err != nil {
		return err
	}

	for _, syncState := range syncStates {
		syncState.SyncStatus = models.SyncSuspended
		if err := s.repo.UpdateSyncState(syncState); err != nil {
			return err
		}
	}

	return nil
}

func (s *SyncService) ResumePlatformSync(ctx context.Context, deviceID string) error {
	syncStates, err := s.repo.GetSyncStateByDevice(deviceID)
	if err != nil {
		return err
	}

	for _, syncState := range syncStates {
		if syncState.SyncStatus == models.SyncSuspended {
			syncState.SyncStatus = models.SyncPending
			syncState.ScheduleNextSync()
			if err := s.repo.UpdateSyncState(syncState); err != nil {
				return err
			}
		}
	}

	return nil
}

// Metrics and monitoring

func (s *SyncService) GetSyncMetrics(ctx context.Context, sessionID uuid.UUID) (*SyncMetrics, error) {
	syncStates, err := s.repo.GetSyncStateBySession(sessionID)
	if err != nil {
		return nil, err
	}

	metrics := &SyncMetrics{
		SessionID:       sessionID,
		TotalPlatforms:  len(syncStates),
		ActivePlatforms: 0,
		ConflictCount:   0,
	}

	var totalLatency float64
	var totalSuccessRate float64

	for _, syncState := range syncStates {
		if syncState.SyncStatus != models.SyncFailed && syncState.SyncStatus != models.SyncSuspended {
			metrics.ActivePlatforms++
		}

		if syncState.ConflictState.HasConflicts {
			metrics.ConflictCount++
		}

		totalLatency += syncState.SyncMetrics.AverageSyncLatency
		totalSuccessRate += syncState.SyncMetrics.SyncSuccessRate

		// Aggregate platform-specific metrics
		if metrics.PlatformMetrics == nil {
			metrics.PlatformMetrics = make(map[string]PlatformSyncMetrics)
		}

		metrics.PlatformMetrics[string(syncState.PlatformType)] = PlatformSyncMetrics{
			Platform:      string(syncState.PlatformType),
			AvgLatency:    syncState.SyncMetrics.AverageSyncLatency,
			SuccessRate:   syncState.SyncMetrics.SyncSuccessRate,
			LastSyncTime:  syncState.LastSyncTime,
			ConflictCount: 0, // TODO: Track per-platform conflicts
		}
	}

	if len(syncStates) > 0 {
		metrics.AvgSyncLatency = totalLatency / float64(len(syncStates))
		metrics.OverallSuccessRate = totalSuccessRate / float64(len(syncStates))
	}

	metrics.PerformanceScore = s.calculatePerformanceScore(syncStates)
	metrics.GeneratedAt = time.Now()

	return metrics, nil
}

func (s *SyncService) GetSyncHealth(ctx context.Context) (*SyncHealth, error) {
	// Get all pending sync states
	pendingStates, err := s.repo.GetPendingSyncStates()
	if err != nil {
		return nil, err
	}

	health := &SyncHealth{
		Healthy:        true,
		PendingCount:   len(pendingStates),
		FailedCount:    0,
		ConflictCount:  0,
		CheckedAt:      time.Now(),
	}

	for _, state := range pendingStates {
		if state.SyncStatus == models.SyncFailed {
			health.FailedCount++
		}

		if state.ConflictState.HasConflicts {
			health.ConflictCount++
		}

		// Check if sync is overdue
		if state.NextSyncTime != nil && time.Since(*state.NextSyncTime) > 30*time.Second {
			health.Warnings = append(health.Warnings, fmt.Sprintf("Sync state %s is overdue", state.ID))
		}
	}

	// Determine health status
	if health.FailedCount > 10 || health.ConflictCount > 5 {
		health.Healthy = false
		health.Status = "unhealthy"
	} else if health.FailedCount > 5 || health.ConflictCount > 2 {
		health.Status = "degraded"
	} else {
		health.Status = "healthy"
	}

	return health, nil
}

func (s *SyncService) CleanupOldSyncStates(ctx context.Context, retentionDays int) (int64, error) {
	return s.repo.CleanupOldSyncStates(retentionDays)
}

// Bulk operations

func (s *SyncService) SyncSessionToAllPlatforms(ctx context.Context, sessionID uuid.UUID) error {
	return s.ForceSyncAll(ctx, sessionID)
}

func (s *SyncService) PropagateChanges(ctx context.Context, sessionID uuid.UUID, changes map[string]interface{}) error {
	syncStates, err := s.repo.GetSyncStateBySession(sessionID)
	if err != nil {
		return err
	}

	for _, syncState := range syncStates {
		// Apply changes
		for key, value := range changes {
			switch key {
			case "position":
				if pos, ok := value.(int); ok {
					syncState.SyncData.Position = pos
				}
			case "isPlaying":
				if playing, ok := value.(bool); ok {
					syncState.SyncData.IsPlaying = playing
				}
			case "quality":
				if quality, ok := value.(string); ok {
					syncState.SyncData.QualitySetting = quality
				}
			case "volume":
				if volume, ok := value.(float64); ok {
					syncState.SyncData.Volume = volume
				}
			}
		}

		if err := s.repo.UpdateSyncState(syncState); err != nil {
			fmt.Printf("Warning: Failed to propagate changes for platform %s: %v\n", syncState.PlatformType, err)
			continue
		}

		s.sendSyncMessage(ctx, syncState, "propagate_changes")
	}

	return nil
}

// Helper methods

func (s *SyncService) broadcastSyncUpdate(ctx context.Context, syncState *models.SynchronizationState) {
	// Send real-time update to all connected clients
	message := SyncMessage{
		Type:       "sync_update",
		SessionID:  syncState.SessionID,
		Platform:   string(syncState.PlatformType),
		DeviceID:   syncState.DeviceID,
		SyncData:   syncState.SyncData,
		Timestamp:  time.Now(),
	}

	s.messagingService.Broadcast(message)
}

func (s *SyncService) sendSyncMessage(ctx context.Context, syncState *models.SynchronizationState, messageType string) {
	message := SyncMessage{
		Type:       messageType,
		SessionID:  syncState.SessionID,
		Platform:   string(syncState.PlatformType),
		DeviceID:   syncState.DeviceID,
		SyncData:   syncState.SyncData,
		Timestamp:  time.Now(),
	}

	s.messagingService.SendToDevice(syncState.DeviceID, message)
}

func (s *SyncService) calculatePerformanceScore(syncStates []*models.SynchronizationState) float64 {
	if len(syncStates) == 0 {
		return 0.0
	}

	totalScore := 0.0
	for _, state := range syncStates {
		totalScore += state.GetSyncPerformanceScore()
	}

	return totalScore / float64(len(syncStates))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// Supporting types

type UIState struct {
	FullscreenMode   bool    `json:"fullscreen_mode"`
	ControlsVisible  bool    `json:"controls_visible"`
	SubtitlesEnabled bool    `json:"subtitles_enabled"`
	SubtitleLanguage string  `json:"subtitle_language"`
	Volume           float64 `json:"volume"`
	IsMuted          bool    `json:"is_muted"`
}

type SyncMetrics struct {
	SessionID          uuid.UUID                       `json:"session_id"`
	TotalPlatforms     int                             `json:"total_platforms"`
	ActivePlatforms    int                             `json:"active_platforms"`
	ConflictCount      int                             `json:"conflict_count"`
	AvgSyncLatency     float64                         `json:"avg_sync_latency"`
	OverallSuccessRate float64                         `json:"overall_success_rate"`
	PerformanceScore   float64                         `json:"performance_score"`
	PlatformMetrics    map[string]PlatformSyncMetrics  `json:"platform_metrics"`
	GeneratedAt        time.Time                       `json:"generated_at"`
}

type PlatformSyncMetrics struct {
	Platform      string     `json:"platform"`
	AvgLatency    float64    `json:"avg_latency"`
	SuccessRate   float64    `json:"success_rate"`
	LastSyncTime  *time.Time `json:"last_sync_time"`
	ConflictCount int        `json:"conflict_count"`
}

type SyncHealth struct {
	Healthy       bool      `json:"healthy"`
	Status        string    `json:"status"`
	PendingCount  int       `json:"pending_count"`
	FailedCount   int       `json:"failed_count"`
	ConflictCount int       `json:"conflict_count"`
	Warnings      []string  `json:"warnings"`
	CheckedAt     time.Time `json:"checked_at"`
}

type SyncMessage struct {
	Type      string           `json:"type"`
	SessionID uuid.UUID        `json:"session_id"`
	Platform  string           `json:"platform"`
	DeviceID  string           `json:"device_id"`
	SyncData  models.SyncData  `json:"sync_data"`
	Timestamp time.Time        `json:"timestamp"`
}

// External service interfaces

type MessagingServiceInterface interface {
	Broadcast(message interface{})
	SendToDevice(deviceID string, message interface{})
	SendToSession(sessionID uuid.UUID, message interface{})
}

type ConfigServiceInterface interface {
	GetSyncFrequency() int
	GetConflictResolutionStrategy() models.ResolutionStrategy
	GetMaxRetries() int
}