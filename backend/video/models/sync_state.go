// backend/video/models/sync_state.go
// Synchronization State model - Cross-platform video sync coordination
// Implements T029: Synchronization State model for real-time cross-platform synchronization

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

// SynchronizationState manages cross-platform video synchronization
type SynchronizationState struct {
	// Primary identifier
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Session relationship
	SessionID uuid.UUID `gorm:"type:uuid;not null;index" json:"session_id" validate:"required"`

	// Platform identification
	PlatformType    PlatformType `gorm:"type:varchar(20);not null" json:"platform_type" validate:"required"`
	DeviceID        string       `gorm:"not null;index" json:"device_id" validate:"required"`
	UserAgent       string       `json:"user_agent"`
	AppVersion      string       `json:"app_version"`

	// Synchronization data
	SyncData        SyncData        `gorm:"embedded" json:"sync_data"`
	ConflictState   ConflictState   `gorm:"embedded" json:"conflict_state"`
	SyncMetrics     SyncMetrics     `gorm:"embedded" json:"sync_metrics"`

	// State management
	SyncStatus      SyncStatus      `gorm:"type:varchar(20);default:'pending'" json:"sync_status"`
	LastSyncTime    *time.Time      `json:"last_sync_time,omitempty"`
	NextSyncTime    *time.Time      `json:"next_sync_time,omitempty"`
	SyncFrequency   int             `gorm:"default:5" json:"sync_frequency"` // seconds

	// Error tracking
	ErrorHistory    []SyncError     `gorm:"type:jsonb" json:"error_history,omitempty"`
	RetryCount      int             `gorm:"default:0" json:"retry_count"`
	MaxRetries      int             `gorm:"default:3" json:"max_retries"`

	// Relationships
	PlaybackSession *PlaybackSession `gorm:"foreignKey:SessionID" json:"playback_session,omitempty"`

	// Soft delete
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Audit fields
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// SyncData contains the actual synchronization information
type SyncData struct {
	// Playback synchronization
	Position            int     `json:"position"`               // Current playback position in seconds
	PlaybackSpeed       float64 `json:"playback_speed"`         // Current playback speed
	Volume              float64 `json:"volume"`                 // Volume level (0.0-1.0)
	QualitySetting      string  `json:"quality_setting"`        // Current video quality
	IsPlaying           bool    `json:"is_playing"`             // Play/pause state
	IsMuted             bool    `json:"is_muted"`               // Mute state

	// UI synchronization
	FullscreenMode      bool    `json:"fullscreen_mode"`        // Fullscreen state
	ControlsVisible     bool    `json:"controls_visible"`       // UI controls visibility
	SubtitlesEnabled    bool    `json:"subtitles_enabled"`      // Subtitle state
	SubtitleLanguage    string  `json:"subtitle_language"`      // Selected subtitle language

	// Progress tracking
	WatchedSegments     []TimeSegment   `gorm:"type:jsonb" json:"watched_segments,omitempty"`
	BookmarkedTimes     []int           `gorm:"type:jsonb" json:"bookmarked_times,omitempty"`
	LastWatchedTime     *time.Time      `json:"last_watched_time,omitempty"`

	// Platform-specific data
	PlatformSpecific    map[string]interface{} `gorm:"type:jsonb" json:"platform_specific,omitempty"`
}

// ConflictState manages synchronization conflicts between platforms
type ConflictState struct {
	HasConflicts        bool            `json:"has_conflicts"`
	ConflictType        ConflictType    `gorm:"type:varchar(30)" json:"conflict_type"`
	ConflictResolution  ResolutionStrategy `gorm:"type:varchar(30)" json:"conflict_resolution"`
	ConflictData        ConflictData    `gorm:"embedded" json:"conflict_data"`
	ResolvedAt          *time.Time      `json:"resolved_at,omitempty"`
	ResolutionMethod    string          `json:"resolution_method,omitempty"`
}

// ConflictData stores detailed conflict information
type ConflictData struct {
	ConflictingPlatforms    []string        `gorm:"type:jsonb" json:"conflicting_platforms,omitempty"`
	ConflictingValues       map[string]interface{} `gorm:"type:jsonb" json:"conflicting_values,omitempty"`
	TimestampDifference     int             `json:"timestamp_difference"` // milliseconds
	PositionDifference      int             `json:"position_difference"`  // seconds
	AuthorityPlatform       string          `json:"authority_platform"`   // Platform with authority to resolve
	ConflictSeverity        SeverityLevel   `gorm:"type:varchar(20)" json:"conflict_severity"`
}

// SyncMetrics tracks synchronization performance
type SyncMetrics struct {
	// Performance metrics
	LastSyncLatency     int     `json:"last_sync_latency"`      // milliseconds
	AverageSyncLatency  float64 `json:"average_sync_latency"`   // milliseconds
	SyncSuccessRate     float64 `json:"sync_success_rate"`      // 0.0-1.0
	TotalSyncAttempts   int     `json:"total_sync_attempts"`
	SuccessfulSyncs     int     `json:"successful_syncs"`
	FailedSyncs         int     `json:"failed_syncs"`

	// Network metrics
	DataTransferred     int64   `json:"data_transferred"`       // bytes
	CompressionRatio    float64 `json:"compression_ratio"`      // data reduction ratio
	NetworkEfficiency   float64 `json:"network_efficiency"`     // 0.0-1.0

	// Quality metrics
	SyncAccuracy        float64 `json:"sync_accuracy"`          // position accuracy in seconds
	StateConsistency    float64 `json:"state_consistency"`      // 0.0-1.0
	ConflictRate        float64 `json:"conflict_rate"`          // conflicts per sync
	ResolutionTime      float64 `json:"resolution_time"`        // average conflict resolution time

	// User experience metrics
	UserSatisfaction    float64 `json:"user_satisfaction"`      // 0.0-1.0 (from feedback)
	SeamlessTransitions int     `json:"seamless_transitions"`   // count of smooth platform switches
	InterruptedSessions int     `json:"interrupted_sessions"`   // count of sync interruptions
}

// TimeSegment represents a watched segment of video
type TimeSegment struct {
	Start       int     `json:"start"`      // Start time in seconds
	End         int     `json:"end"`        // End time in seconds
	WatchedAt   time.Time `json:"watched_at"` // When this segment was watched
	Platform    string  `json:"platform"`   // Platform where it was watched
	Quality     string  `json:"quality"`    // Video quality during watching
}

// SyncError represents a synchronization error
type SyncError struct {
	Timestamp   time.Time   `json:"timestamp"`
	ErrorType   string      `json:"error_type"`
	ErrorCode   string      `json:"error_code"`
	Message     string      `json:"message"`
	Platform    string      `json:"platform"`
	Severity    SeverityLevel `json:"severity"`
	Resolved    bool        `json:"resolved"`
	ResolvedAt  *time.Time  `json:"resolved_at,omitempty"`
}

// Enums for synchronization management

type SyncStatus string

const (
	SyncPending     SyncStatus = "pending"
	SyncInProgress  SyncStatus = "in_progress"
	SyncCompleted   SyncStatus = "completed"
	SyncFailed      SyncStatus = "failed"
	SyncConflicted  SyncStatus = "conflicted"
	SyncSuspended   SyncStatus = "suspended"
)

type ConflictType string

const (
	ConflictPosition    ConflictType = "position_mismatch"
	ConflictPlayState   ConflictType = "play_state_mismatch"
	ConflictQuality     ConflictType = "quality_mismatch"
	ConflictVolume      ConflictType = "volume_mismatch"
	ConflictUI          ConflictType = "ui_state_mismatch"
	ConflictTimestamp   ConflictType = "timestamp_conflict"
	ConflictPlatform    ConflictType = "platform_conflict"
)

type ResolutionStrategy string

const (
	ResolutionLatest        ResolutionStrategy = "latest_timestamp"
	ResolutionAuthority     ResolutionStrategy = "authority_platform"
	ResolutionMajority      ResolutionStrategy = "majority_vote"
	ResolutionUserChoice    ResolutionStrategy = "user_choice"
	ResolutionAverage       ResolutionStrategy = "average_value"
	ResolutionManual        ResolutionStrategy = "manual_resolution"
)

type SeverityLevel string

const (
	SeverityLow         SeverityLevel = "low"
	SeverityMedium      SeverityLevel = "medium"
	SeverityHigh        SeverityLevel = "high"
	SeverityCritical    SeverityLevel = "critical"
)

// Business logic methods

// UpdateSyncData updates the synchronization data with new values
func (ss *SynchronizationState) UpdateSyncData(newData SyncData) {
	// Check for conflicts before updating
	if ss.HasPositionConflict(newData.Position) {
		ss.RecordConflict(ConflictPosition, map[string]interface{}{
			"current_position": ss.SyncData.Position,
			"new_position":     newData.Position,
		})
	}

	if ss.HasPlayStateConflict(newData.IsPlaying) {
		ss.RecordConflict(ConflictPlayState, map[string]interface{}{
			"current_state": ss.SyncData.IsPlaying,
			"new_state":     newData.IsPlaying,
		})
	}

	// Update sync data
	ss.SyncData = newData
	ss.LastSyncTime = timePtr(time.Now())
	ss.ScheduleNextSync()
}

// HasPositionConflict checks for position synchronization conflicts
func (ss *SynchronizationState) HasPositionConflict(newPosition int) bool {
	positionDiff := abs(ss.SyncData.Position - newPosition)
	return positionDiff > 5 // More than 5 seconds difference is a conflict
}

// HasPlayStateConflict checks for play state conflicts
func (ss *SynchronizationState) HasPlayStateConflict(newPlayState bool) bool {
	return ss.SyncData.IsPlaying != newPlayState
}

// RecordConflict records a synchronization conflict
func (ss *SynchronizationState) RecordConflict(conflictType ConflictType, conflictingValues map[string]interface{}) {
	ss.ConflictState.HasConflicts = true
	ss.ConflictState.ConflictType = conflictType
	ss.ConflictState.ConflictData.ConflictingValues = conflictingValues
	ss.ConflictState.ConflictData.ConflictSeverity = ss.calculateSeverity(conflictType)
	ss.SyncStatus = SyncConflicted
}

// calculateSeverity determines the severity of a conflict
func (ss *SynchronizationState) calculateSeverity(conflictType ConflictType) SeverityLevel {
	switch conflictType {
	case ConflictPosition:
		return SeverityHigh
	case ConflictPlayState:
		return SeverityMedium
	case ConflictQuality, ConflictVolume:
		return SeverityLow
	case ConflictTimestamp, ConflictPlatform:
		return SeverityCritical
	default:
		return SeverityMedium
	}
}

// ResolveConflict resolves synchronization conflicts using the specified strategy
func (ss *SynchronizationState) ResolveConflict(strategy ResolutionStrategy) error {
	if !ss.ConflictState.HasConflicts {
		return fmt.Errorf("no conflicts to resolve")
	}

	switch strategy {
	case ResolutionLatest:
		return ss.resolveByLatestTimestamp()
	case ResolutionAuthority:
		return ss.resolveByAuthority()
	case ResolutionMajority:
		return ss.resolveByMajority()
	case ResolutionAverage:
		return ss.resolveByAverage()
	default:
		return fmt.Errorf("unsupported resolution strategy: %s", strategy)
	}
}

// resolveByLatestTimestamp resolves conflicts by using the latest timestamp
func (ss *SynchronizationState) resolveByLatestTimestamp() error {
	// Implementation would compare timestamps and use the most recent data
	ss.ConflictState.HasConflicts = false
	ss.ConflictState.ResolvedAt = timePtr(time.Now())
	ss.ConflictState.ResolutionMethod = "latest_timestamp"
	ss.SyncStatus = SyncCompleted
	return nil
}

// resolveByAuthority resolves conflicts using a designated authority platform
func (ss *SynchronizationState) resolveByAuthority() error {
	if ss.ConflictState.ConflictData.AuthorityPlatform == "" {
		return fmt.Errorf("no authority platform specified")
	}

	// Implementation would use data from the authority platform
	ss.ConflictState.HasConflicts = false
	ss.ConflictState.ResolvedAt = timePtr(time.Now())
	ss.ConflictState.ResolutionMethod = "authority_platform"
	ss.SyncStatus = SyncCompleted
	return nil
}

// resolveByMajority resolves conflicts using majority consensus
func (ss *SynchronizationState) resolveByMajority() error {
	// Implementation would analyze data from all platforms and use majority value
	ss.ConflictState.HasConflicts = false
	ss.ConflictState.ResolvedAt = timePtr(time.Now())
	ss.ConflictState.ResolutionMethod = "majority_vote"
	ss.SyncStatus = SyncCompleted
	return nil
}

// resolveByAverage resolves conflicts by averaging conflicting values
func (ss *SynchronizationState) resolveByAverage() error {
	// Implementation would calculate average of conflicting numerical values
	ss.ConflictState.HasConflicts = false
	ss.ConflictState.ResolvedAt = timePtr(time.Now())
	ss.ConflictState.ResolutionMethod = "average_value"
	ss.SyncStatus = SyncCompleted
	return nil
}

// ScheduleNextSync schedules the next synchronization attempt
func (ss *SynchronizationState) ScheduleNextSync() {
	nextSync := time.Now().Add(time.Duration(ss.SyncFrequency) * time.Second)
	ss.NextSyncTime = &nextSync
}

// UpdateSyncMetrics updates synchronization performance metrics
func (ss *SynchronizationState) UpdateSyncMetrics(latency int, success bool, dataSize int64) {
	ss.SyncMetrics.TotalSyncAttempts++
	ss.SyncMetrics.LastSyncLatency = latency

	if success {
		ss.SyncMetrics.SuccessfulSyncs++
		ss.SyncMetrics.DataTransferred += dataSize
	} else {
		ss.SyncMetrics.FailedSyncs++
	}

	// Calculate running averages
	ss.SyncMetrics.SyncSuccessRate = float64(ss.SyncMetrics.SuccessfulSyncs) / float64(ss.SyncMetrics.TotalSyncAttempts)

	if ss.SyncMetrics.SuccessfulSyncs > 0 {
		ss.SyncMetrics.AverageSyncLatency = (ss.SyncMetrics.AverageSyncLatency*float64(ss.SyncMetrics.SuccessfulSyncs-1) + float64(latency)) / float64(ss.SyncMetrics.SuccessfulSyncs)
	}
}

// RecordError records a synchronization error
func (ss *SynchronizationState) RecordError(errorType, errorCode, message string, severity SeverityLevel) {
	syncError := SyncError{
		Timestamp: time.Now(),
		ErrorType: errorType,
		ErrorCode: errorCode,
		Message:   message,
		Platform:  string(ss.PlatformType),
		Severity:  severity,
		Resolved:  false,
	}

	ss.ErrorHistory = append(ss.ErrorHistory, syncError)
	ss.RetryCount++

	if ss.RetryCount >= ss.MaxRetries {
		ss.SyncStatus = SyncFailed
	}
}

// IsHealthy checks if the synchronization state is healthy
func (ss *SynchronizationState) IsHealthy() bool {
	return ss.SyncStatus != SyncFailed &&
		   ss.SyncMetrics.SyncSuccessRate > 0.8 &&
		   !ss.ConflictState.HasConflicts &&
		   ss.RetryCount < ss.MaxRetries
}

// GetSyncPerformanceScore calculates overall sync performance score (0.0-1.0)
func (ss *SynchronizationState) GetSyncPerformanceScore() float64 {
	score := 0.0

	// Success rate component (40% weight)
	score += ss.SyncMetrics.SyncSuccessRate * 0.4

	// Latency component (30% weight)
	if ss.SyncMetrics.AverageSyncLatency > 0 {
		latencyScore := 1.0 - (ss.SyncMetrics.AverageSyncLatency / 1000.0) // Normalize to 1 second
		if latencyScore < 0 {
			latencyScore = 0
		}
		score += latencyScore * 0.3
	}

	// Conflict rate component (20% weight)
	conflictScore := 1.0 - ss.SyncMetrics.ConflictRate
	if conflictScore < 0 {
		conflictScore = 0
	}
	score += conflictScore * 0.2

	// State consistency component (10% weight)
	score += ss.SyncMetrics.StateConsistency * 0.1

	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	return score
}

// Utility functions

func timePtr(t time.Time) *time.Time {
	return &t
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GORM hooks

// BeforeCreate sets up defaults before creating synchronization state record
func (ss *SynchronizationState) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}

	if ss.SyncStatus == "" {
		ss.SyncStatus = SyncPending
	}

	if ss.SyncFrequency <= 0 {
		ss.SyncFrequency = 5 // Default 5 seconds
	}

	if ss.MaxRetries <= 0 {
		ss.MaxRetries = 3
	}

	// Initialize metrics
	if ss.SyncMetrics.SyncSuccessRate == 0 {
		ss.SyncMetrics.SyncSuccessRate = 1.0
	}

	ss.ScheduleNextSync()
	return nil
}

// BeforeUpdate validates state before updating
func (ss *SynchronizationState) BeforeUpdate(tx *gorm.DB) error {
	// Validate state transitions
	if ss.SyncStatus == SyncCompleted && ss.ConflictState.HasConflicts {
		return fmt.Errorf("cannot mark sync as completed while conflicts exist")
	}

	return nil
}

// TableName returns the table name for GORM
func (SynchronizationState) TableName() string {
	return "synchronization_states"
}

// Database indexes for performance
/*
CREATE INDEX CONCURRENTLY idx_synchronization_states_session_id ON synchronization_states(session_id);
CREATE INDEX CONCURRENTLY idx_synchronization_states_device_id ON synchronization_states(device_id);
CREATE INDEX CONCURRENTLY idx_synchronization_states_platform_type ON synchronization_states(platform_type);
CREATE INDEX CONCURRENTLY idx_synchronization_states_sync_status ON synchronization_states(sync_status);
CREATE INDEX CONCURRENTLY idx_synchronization_states_last_sync_time ON synchronization_states(last_sync_time DESC);
CREATE INDEX CONCURRENTLY idx_synchronization_states_next_sync_time ON synchronization_states(next_sync_time ASC);
CREATE INDEX CONCURRENTLY idx_synchronization_states_has_conflicts ON synchronization_states(conflict_state_has_conflicts);
CREATE INDEX CONCURRENTLY idx_synchronization_states_device_platform ON synchronization_states(device_id, platform_type);
*/