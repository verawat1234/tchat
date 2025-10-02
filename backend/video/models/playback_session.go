// backend/video/models/playback_session.go
// Playback Session model - Real-time viewing state across platforms
// Implements T026: Playback Session model with cross-platform synchronization

package models

import (
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

// PlaybackSession tracks user's current viewing state across platforms
type PlaybackSession struct {
	// Primary identifier
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// References
	UserID  uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	VideoID uuid.UUID `gorm:"type:uuid;not null;index" json:"video_id" validate:"required"`

	// Playback state
	CurrentPosition int     `gorm:"not null;default:0" json:"current_position" validate:"gte=0"`
	PlaybackSpeed   float64 `gorm:"not null;default:1.0" json:"playback_speed" validate:"oneof=0.25 0.5 0.75 1.0 1.25 1.5 2.0"`
	QualitySetting  string  `gorm:"not null;default:'auto'" json:"quality_setting" validate:"oneof=auto 360p 480p 720p 1080p 4K"`
	VolumeLevel     float64 `gorm:"not null;default:1.0" json:"volume_level" validate:"gte=0,lte=1"`

	// Platform context
	PlatformContext string `gorm:"not null" json:"platform_context" validate:"oneof=web android ios"`
	DeviceInfo      DeviceInformation `gorm:"embedded" json:"device_info"`

	// Session timing
	SessionStart time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"session_start"`
	LastUpdated  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"last_updated"`
	SessionEnd   *time.Time `json:"session_end,omitempty"`

	// Session state
	IsActive     bool          `gorm:"not null;default:true" json:"is_active"`
	BufferStatus BufferStatus  `gorm:"embedded" json:"buffer_status"`
	SessionState SessionState  `gorm:"type:varchar(20);default:'active'" json:"session_state"`

	// Performance metrics
	StartupTime      int `json:"startup_time,omitempty"`    // Time to first frame in ms
	BufferingEvents  int `gorm:"default:0" json:"buffering_events"`
	QualitySwitches  int `gorm:"default:0" json:"quality_switches"`
	TotalBufferTime  int `json:"total_buffer_time,omitempty"` // Total buffering time in ms

	// Relationships
	VideoContent        *VideoContent        `gorm:"foreignKey:VideoID" json:"video_content,omitempty"`
	SynchronizationStates []SynchronizationState `gorm:"foreignKey:SessionID" json:"synchronization_states,omitempty"`

	// Soft delete
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// DeviceInformation contains device specifications and capabilities
type DeviceInformation struct {
	UserAgent        string `json:"user_agent"`
	Platform         string `json:"platform"`
	PlatformVersion  string `json:"platform_version"`
	AppVersion       string `json:"app_version"`
	ScreenResolution string `json:"screen_resolution"`
	NetworkType      string `json:"network_type"`
	HardwareDecoding bool   `json:"hardware_decoding"`
	MaxResolution    string `json:"max_resolution"`
}

// BufferStatus contains buffering state and network conditions
type BufferStatus struct {
	BufferedAhead    int     `json:"buffered_ahead"`     // Seconds buffered ahead
	DownloadSpeed    float64 `json:"download_speed"`     // Mbps
	NetworkLatency   int     `json:"network_latency"`    // ms
	BufferHealth     float64 `json:"buffer_health"`      // 0.0-1.0
	IsBuffering      bool    `json:"is_buffering"`
	LastBufferEvent  *time.Time `json:"last_buffer_event,omitempty"`
}

// SessionState represents the current state of the playback session
type SessionState string

const (
	SessionActive    SessionState = "active"
	SessionPaused    SessionState = "paused"
	SessionBuffering SessionState = "buffering"
	SessionCompleted SessionState = "completed"
	SessionAbandoned SessionState = "abandoned"
	SessionError     SessionState = "error"
)

// Business logic methods

// GetWatchedPercentage calculates how much of the video has been watched
func (ps *PlaybackSession) GetWatchedPercentage(videoDuration int) float64 {
	if videoDuration <= 0 {
		return 0
	}

	percentage := float64(ps.CurrentPosition) / float64(videoDuration) * 100
	if percentage > 100 {
		percentage = 100
	}

	return percentage
}

// IsRecentlyActive checks if session was active within the last interval
func (ps *PlaybackSession) IsRecentlyActive(interval time.Duration) bool {
	return ps.IsActive && time.Since(ps.LastUpdated) < interval
}

// GetSessionDuration returns the total session duration
func (ps *PlaybackSession) GetSessionDuration() time.Duration {
	endTime := time.Now()
	if ps.SessionEnd != nil {
		endTime = *ps.SessionEnd
	}

	return endTime.Sub(ps.SessionStart)
}

// UpdatePosition updates the current playback position and session state
func (ps *PlaybackSession) UpdatePosition(position int) {
	ps.CurrentPosition = position
	ps.LastUpdated = time.Now()
	ps.IsActive = true
}

// ChangeQuality updates the quality setting and tracks the change
func (ps *PlaybackSession) ChangeQuality(newQuality string) {
	if ps.QualitySetting != newQuality {
		ps.QualitySetting = newQuality
		ps.QualitySwitches++
		ps.LastUpdated = time.Now()
	}
}

// RecordBufferEvent logs a buffering event
func (ps *PlaybackSession) RecordBufferEvent() {
	ps.BufferingEvents++
	now := time.Now()
	ps.BufferStatus.LastBufferEvent = &now
	ps.BufferStatus.IsBuffering = true
}

// EndBuffering marks the end of a buffering event
func (ps *PlaybackSession) EndBuffering(bufferDuration int) {
	ps.BufferStatus.IsBuffering = false
	ps.TotalBufferTime += bufferDuration
}

// CompleteSession marks the session as completed
func (ps *PlaybackSession) CompleteSession() {
	ps.IsActive = false
	ps.SessionState = SessionCompleted
	now := time.Now()
	ps.SessionEnd = &now
}

// AbandonSession marks the session as abandoned
func (ps *PlaybackSession) AbandonSession() {
	ps.IsActive = false
	ps.SessionState = SessionAbandoned
	now := time.Now()
	ps.SessionEnd = &now
}

// PauseSession sets the session to paused state
func (ps *PlaybackSession) PauseSession() {
	ps.SessionState = SessionPaused
	ps.LastUpdated = time.Now()
}

// ResumeSession sets the session back to active state
func (ps *PlaybackSession) ResumeSession() {
	ps.SessionState = SessionActive
	ps.IsActive = true
	ps.LastUpdated = time.Now()
}

// Validation methods

// IsValidQualitySetting checks if the quality setting is supported
func (ps *PlaybackSession) IsValidQualitySetting(quality string) bool {
	validQualities := []string{"auto", "360p", "480p", "720p", "1080p", "4K"}
	for _, valid := range validQualities {
		if valid == quality {
			return true
		}
	}
	return false
}

// IsValidPlaybackSpeed checks if the playback speed is supported
func (ps *PlaybackSession) IsValidPlaybackSpeed(speed float64) bool {
	validSpeeds := []float64{0.25, 0.5, 0.75, 1.0, 1.25, 1.5, 2.0}
	for _, valid := range validSpeeds {
		if valid == speed {
			return true
		}
	}
	return false
}

// Performance analysis methods

// GetBufferHealthScore calculates overall buffer health (0.0-1.0)
func (ps *PlaybackSession) GetBufferHealthScore() float64 {
	if ps.BufferingEvents == 0 {
		return 1.0
	}

	sessionDuration := ps.GetSessionDuration().Seconds()
	if sessionDuration <= 0 {
		return 1.0
	}

	// Calculate buffer health based on buffering frequency and total buffer time
	bufferRatio := float64(ps.TotalBufferTime) / 1000 / sessionDuration
	eventFrequency := float64(ps.BufferingEvents) / (sessionDuration / 60) // events per minute

	// Lower score for more buffering
	healthScore := 1.0 - (bufferRatio*0.5 + eventFrequency*0.1)
	if healthScore < 0 {
		healthScore = 0
	}

	return healthScore
}

// GetPlaybackQualityMetrics returns comprehensive playback quality metrics
func (ps *PlaybackSession) GetPlaybackQualityMetrics() map[string]interface{} {
	return map[string]interface{}{
		"startup_time":        ps.StartupTime,
		"buffering_events":    ps.BufferingEvents,
		"total_buffer_time":   ps.TotalBufferTime,
		"quality_switches":    ps.QualitySwitches,
		"buffer_health_score": ps.GetBufferHealthScore(),
		"session_duration":    ps.GetSessionDuration().Seconds(),
		"average_quality":     ps.QualitySetting,
		"platform":           ps.PlatformContext,
	}
}

// GORM hooks

// BeforeCreate sets up defaults before creating session record
func (ps *PlaybackSession) BeforeCreate(tx *gorm.DB) error {
	if ps.ID == uuid.Nil {
		ps.ID = uuid.New()
	}

	if ps.SessionState == "" {
		ps.SessionState = SessionActive
	}

	if ps.QualitySetting == "" {
		ps.QualitySetting = "auto"
	}

	if ps.PlaybackSpeed == 0 {
		ps.PlaybackSpeed = 1.0
	}

	if ps.VolumeLevel == 0 {
		ps.VolumeLevel = 1.0
	}

	return nil
}

// BeforeUpdate validates state transitions before updating
func (ps *PlaybackSession) BeforeUpdate(tx *gorm.DB) error {
	ps.LastUpdated = time.Now()
	return nil
}

// TableName returns the table name for GORM
func (PlaybackSession) TableName() string {
	return "playback_sessions"
}

// Database indexes for performance
/*
CREATE INDEX CONCURRENTLY idx_playback_sessions_user_id ON playback_sessions(user_id);
CREATE INDEX CONCURRENTLY idx_playback_sessions_video_id ON playback_sessions(video_id);
CREATE INDEX CONCURRENTLY idx_playback_sessions_user_video ON playback_sessions(user_id, video_id);
CREATE INDEX CONCURRENTLY idx_playback_sessions_is_active ON playback_sessions(is_active);
CREATE INDEX CONCURRENTLY idx_playback_sessions_last_updated ON playback_sessions(last_updated DESC);
CREATE INDEX CONCURRENTLY idx_playback_sessions_platform ON playback_sessions(platform_context);
CREATE INDEX CONCURRENTLY idx_playback_sessions_session_state ON playback_sessions(session_state);
*/