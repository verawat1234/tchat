// backend/video/models/viewing_history.go
// Viewing History model - Historical record of user's video consumption patterns
// Implements T027: Viewing History model for analytics and recommendations

package models

import (
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

// ViewingHistory records user's video consumption patterns for analytics and recommendations
type ViewingHistory struct {
	// Primary identifier
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// References
	UserID  uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	VideoID uuid.UUID `gorm:"type:uuid;not null;index" json:"video_id" validate:"required"`

	// Viewing details
	WatchedAt       time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"watched_at"`
	CompletionRate  float64   `gorm:"not null;default:0" json:"completion_rate" validate:"gte=0,lte=100"`
	WatchedDuration int       `gorm:"not null;default:0" json:"watched_duration" validate:"gte=0"` // seconds
	TotalDuration   int       `gorm:"not null" json:"total_duration" validate:"gt=0"`

	// Session context
	PlatformUsed    string `gorm:"not null" json:"platform_used" validate:"oneof=web android ios"`
	QualityWatched  string `gorm:"not null;default:'auto'" json:"quality_watched"`
	DeviceCategory  string `gorm:"not null;default:'unknown'" json:"device_category" validate:"oneof=desktop mobile tablet tv unknown"`

	// Engagement metrics
	InteractionData InteractionData `gorm:"embedded" json:"interaction_data"`

	// Behavioral analysis
	ViewingPattern  ViewingPattern  `gorm:"embedded" json:"viewing_pattern"`
	EngagementScore float64        `gorm:"default:0" json:"engagement_score"` // 0.0-1.0

	// Recommendation signals
	RecommendationContext RecommendationContext `gorm:"embedded" json:"recommendation_context"`

	// Session metrics
	StartupLatency    int `json:"startup_latency,omitempty"`    // ms to start playing
	BufferingTime     int `json:"buffering_time,omitempty"`     // total buffering in ms
	QualitySwitches   int `gorm:"default:0" json:"quality_switches"`
	NetworkConditions NetworkConditions `gorm:"embedded" json:"network_conditions"`

	// Relationships
	VideoContent *VideoContent `gorm:"foreignKey:VideoID" json:"video_content,omitempty"`

	// Audit fields
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// InteractionData tracks user interactions during viewing
type InteractionData struct {
	SeekCount        int     `json:"seek_count" gorm:"default:0"`
	PauseCount       int     `json:"pause_count" gorm:"default:0"`
	VolumeAdjustments int    `json:"volume_adjustments" gorm:"default:0"`
	FullscreenToggle int     `json:"fullscreen_toggle" gorm:"default:0"`
	SpeedChanges     int     `json:"speed_changes" gorm:"default:0"`
	LikedVideo       bool    `json:"liked_video" gorm:"default:false"`
	SharedVideo      bool    `json:"shared_video" gorm:"default:false"`
	CommentedOnVideo bool    `json:"commented_on_video" gorm:"default:false"`
	BookmarkedVideo  bool    `json:"bookmarked_video" gorm:"default:false"`
	AveragePlaybackSpeed float64 `json:"average_playback_speed" gorm:"default:1.0"`
}

// ViewingPattern analyzes how the user consumed the content
type ViewingPattern struct {
	ViewingStyle    ViewingStyle `gorm:"type:varchar(20);default:'continuous'" json:"viewing_style"`
	AttentionSpans  []int        `gorm:"type:jsonb" json:"attention_spans,omitempty"` // seconds of continuous watching
	SkippedSegments []TimeRange  `gorm:"type:jsonb" json:"skipped_segments,omitempty"`
	ReplayedSegments []TimeRange `gorm:"type:jsonb" json:"replayed_segments,omitempty"`
	MostWatchedSegment *TimeRange `gorm:"embedded;embeddedPrefix:most_watched_" json:"most_watched_segment,omitempty"`
	WatchingVelocity float64     `json:"watching_velocity"` // actual time / video time
}

// RecommendationContext provides signals for recommendation algorithms
type RecommendationContext struct {
	DiscoveryMethod   string    `json:"discovery_method"`   // search, recommendation, trending, etc.
	CategoryInterest  []string  `gorm:"type:jsonb" json:"category_interest,omitempty"`
	CreatorAffinity   float64   `json:"creator_affinity"`   // 0.0-1.0 interest in this creator
	TopicRelevance    float64   `json:"topic_relevance"`    // 0.0-1.0 relevance to user interests
	TimeOfDay         string    `json:"time_of_day"`        // morning, afternoon, evening, night
	DayOfWeek         string    `json:"day_of_week"`
	SeasonalContext   string    `json:"seasonal_context"`   // holiday, season, etc.
	MoodIndicator     string    `json:"mood_indicator"`     // entertainment, education, relaxation
}

// NetworkConditions captures network quality during viewing
type NetworkConditions struct {
	ConnectionType   string  `json:"connection_type"`   // wifi, cellular, ethernet
	Bandwidth        float64 `json:"bandwidth"`         // Mbps
	Latency         int     `json:"latency"`           // ms
	PacketLoss      float64 `json:"packet_loss"`       // percentage
	ConnectionStability float64 `json:"connection_stability"` // 0.0-1.0
}

// TimeRange represents a time segment in the video
type TimeRange struct {
	Start int `json:"start"` // seconds
	End   int `json:"end"`   // seconds
}

// ViewingStyle categorizes how the user watched the content
type ViewingStyle string

const (
	ViewingContinuous ViewingStyle = "continuous"    // Watched without major interruptions
	ViewingFragmented ViewingStyle = "fragmented"    // Multiple pauses and seeks
	ViewingSkimming   ViewingStyle = "skimming"      // Fast-forwarded through content
	ViewingFocused    ViewingStyle = "focused"       // Minimal interactions, high attention
	ViewingCasual     ViewingStyle = "casual"        // Background viewing with interruptions
	ViewingAnalytical ViewingStyle = "analytical"    // Multiple replays and detailed viewing
)

// Business logic methods

// CalculateEngagementScore computes an engagement score based on viewing behavior
func (vh *ViewingHistory) CalculateEngagementScore() float64 {
	score := 0.0

	// Base score from completion rate (40% weight)
	score += (vh.CompletionRate / 100) * 0.4

	// Interaction bonus (30% weight)
	interactionScore := 0.0
	if vh.InteractionData.LikedVideo {
		interactionScore += 0.3
	}
	if vh.InteractionData.CommentedOnVideo {
		interactionScore += 0.3
	}
	if vh.InteractionData.SharedVideo {
		interactionScore += 0.2
	}
	if vh.InteractionData.BookmarkedVideo {
		interactionScore += 0.2
	}

	score += interactionScore * 0.3

	// Viewing pattern bonus (20% weight)
	patternScore := 0.0
	switch vh.ViewingPattern.ViewingStyle {
	case ViewingFocused:
		patternScore = 1.0
	case ViewingContinuous:
		patternScore = 0.8
	case ViewingAnalytical:
		patternScore = 0.9
	case ViewingCasual:
		patternScore = 0.5
	case ViewingFragmented:
		patternScore = 0.3
	case ViewingSkimming:
		patternScore = 0.2
	}
	score += patternScore * 0.2

	// Quality and performance bonus (10% weight)
	performanceScore := 1.0
	if vh.BufferingTime > 0 {
		bufferRatio := float64(vh.BufferingTime) / float64(vh.WatchedDuration*1000)
		performanceScore -= bufferRatio * 0.5 // Reduce score for buffering
	}
	score += performanceScore * 0.1

	// Ensure score is between 0 and 1
	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	vh.EngagementScore = score
	return score
}

// IsCompletedViewing checks if the video was substantially watched
func (vh *ViewingHistory) IsCompletedViewing() bool {
	return vh.CompletionRate >= 80.0 // 80% completion threshold
}

// GetWatchingEfficiency calculates how efficiently the user watched
func (vh *ViewingHistory) GetWatchingEfficiency() float64 {
	if vh.TotalDuration <= 0 {
		return 0
	}

	actualWatchTime := float64(vh.WatchedDuration)
	videoLength := float64(vh.TotalDuration)

	return actualWatchTime / videoLength
}

// GetViewingQuality returns a quality score for the viewing session
func (vh *ViewingHistory) GetViewingQuality() float64 {
	qualityScore := 1.0

	// Reduce score for excessive seeking
	if vh.InteractionData.SeekCount > 10 {
		qualityScore -= 0.2
	}

	// Reduce score for poor network conditions
	if vh.NetworkConditions.PacketLoss > 1.0 {
		qualityScore -= 0.1
	}

	// Reduce score for excessive buffering
	if vh.BufferingTime > vh.WatchedDuration*100 { // More than 10% buffering
		qualityScore -= 0.3
	}

	// Reduce score for quality switches (network instability)
	if vh.QualitySwitches > 5 {
		qualityScore -= 0.1
	}

	if qualityScore < 0 {
		qualityScore = 0
	}

	return qualityScore
}

// UpdateFromPlaybackSession updates history from a completed playback session
func (vh *ViewingHistory) UpdateFromPlaybackSession(session *PlaybackSession) {
	vh.WatchedDuration = session.CurrentPosition
	vh.PlatformUsed = session.PlatformContext
	vh.QualityWatched = session.QualitySetting
	vh.StartupLatency = session.StartupTime
	vh.BufferingTime = session.TotalBufferTime
	vh.QualitySwitches = session.QualitySwitches

	// Calculate completion rate
	if vh.TotalDuration > 0 {
		vh.CompletionRate = float64(vh.WatchedDuration) / float64(vh.TotalDuration) * 100
	}

	// Set device category based on platform and screen info
	switch session.PlatformContext {
	case "web":
		if session.DeviceInfo.ScreenResolution != "" {
			vh.DeviceCategory = "desktop"
		} else {
			vh.DeviceCategory = "desktop"
		}
	case "android", "ios":
		vh.DeviceCategory = "mobile"
	}

	// Update network conditions
	vh.NetworkConditions.ConnectionType = session.DeviceInfo.NetworkType
	vh.NetworkConditions.ConnectionStability = session.GetBufferHealthScore()

	// Calculate engagement score
	vh.CalculateEngagementScore()
}

// GetRecommendationSignals returns signals useful for recommendation algorithms
func (vh *ViewingHistory) GetRecommendationSignals() map[string]interface{} {
	return map[string]interface{}{
		"engagement_score":     vh.EngagementScore,
		"completion_rate":      vh.CompletionRate,
		"viewing_style":        vh.ViewingPattern.ViewingStyle,
		"platform_preference":  vh.PlatformUsed,
		"interaction_level":    vh.GetInteractionLevel(),
		"quality_preference":   vh.QualityWatched,
		"device_category":      vh.DeviceCategory,
		"discovery_method":     vh.RecommendationContext.DiscoveryMethod,
		"creator_affinity":     vh.RecommendationContext.CreatorAffinity,
		"time_context":         vh.RecommendationContext.TimeOfDay,
		"viewing_quality":      vh.GetViewingQuality(),
		"watching_efficiency":  vh.GetWatchingEfficiency(),
	}
}

// GetInteractionLevel calculates the level of user interaction
func (vh *ViewingHistory) GetInteractionLevel() string {
	totalInteractions := vh.InteractionData.SeekCount +
		vh.InteractionData.PauseCount +
		vh.InteractionData.VolumeAdjustments +
		vh.InteractionData.SpeedChanges

	socialInteractions := 0
	if vh.InteractionData.LikedVideo {
		socialInteractions++
	}
	if vh.InteractionData.SharedVideo {
		socialInteractions++
	}
	if vh.InteractionData.CommentedOnVideo {
		socialInteractions++
	}
	if vh.InteractionData.BookmarkedVideo {
		socialInteractions++
	}

	if socialInteractions >= 2 {
		return "highly_engaged"
	} else if socialInteractions >= 1 || totalInteractions >= 5 {
		return "engaged"
	} else if totalInteractions >= 2 {
		return "moderately_engaged"
	} else {
		return "passive"
	}
}

// GORM hooks

// BeforeCreate sets up defaults before creating viewing history record
func (vh *ViewingHistory) BeforeCreate(tx *gorm.DB) error {
	if vh.ID == uuid.Nil {
		vh.ID = uuid.New()
	}

	if vh.ViewingPattern.ViewingStyle == "" {
		vh.ViewingPattern.ViewingStyle = ViewingContinuous
	}

	if vh.DeviceCategory == "" {
		vh.DeviceCategory = "unknown"
	}

	if vh.QualityWatched == "" {
		vh.QualityWatched = "auto"
	}

	return nil
}

// BeforeUpdate calculates engagement score before updating
func (vh *ViewingHistory) BeforeUpdate(tx *gorm.DB) error {
	vh.CalculateEngagementScore()
	return nil
}

// TableName returns the table name for GORM
func (ViewingHistory) TableName() string {
	return "viewing_histories"
}

// Database indexes for performance and analytics
/*
CREATE INDEX CONCURRENTLY idx_viewing_histories_user_id ON viewing_histories(user_id);
CREATE INDEX CONCURRENTLY idx_viewing_histories_video_id ON viewing_histories(video_id);
CREATE INDEX CONCURRENTLY idx_viewing_histories_user_watched_at ON viewing_histories(user_id, watched_at DESC);
CREATE INDEX CONCURRENTLY idx_viewing_histories_completion_rate ON viewing_histories(completion_rate);
CREATE INDEX CONCURRENTLY idx_viewing_histories_engagement_score ON viewing_histories(engagement_score DESC);
CREATE INDEX CONCURRENTLY idx_viewing_histories_platform ON viewing_histories(platform_used);
CREATE INDEX CONCURRENTLY idx_viewing_histories_viewing_style ON viewing_histories(viewing_pattern_viewing_style);
CREATE INDEX CONCURRENTLY idx_viewing_histories_device_category ON viewing_histories(device_category);
*/