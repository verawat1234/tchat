// backend/video/models/platform_config.go
// Platform Configuration model - Cross-platform video settings and capabilities
// Implements T028: Platform Configuration model for cross-platform optimization

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

// PlatformConfiguration manages platform-specific video settings and capabilities
type PlatformConfiguration struct {
	// Primary identifier
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Platform identification
	PlatformType    PlatformType `gorm:"type:varchar(20);not null;index" json:"platform_type" validate:"required"`
	PlatformVersion string       `gorm:"not null" json:"platform_version" validate:"required"`
	DeviceCategory  string       `gorm:"not null" json:"device_category" validate:"oneof=desktop mobile tablet tv"`

	// Video capabilities
	VideoCapabilities VideoCapabilities `gorm:"embedded" json:"video_capabilities"`

	// Quality configurations
	QualitySettings QualityConfiguration `gorm:"embedded" json:"quality_settings"`

	// Performance optimization
	PerformanceSettings PerformanceConfiguration `gorm:"embedded" json:"performance_settings"`

	// Network optimization
	NetworkOptimization NetworkOptimization `gorm:"embedded" json:"network_optimization"`

	// UI/UX configuration
	UIConfiguration UIConfiguration `gorm:"embedded" json:"ui_configuration"`

	// Storage and caching
	StorageConfiguration StorageConfiguration `gorm:"embedded" json:"storage_configuration"`

	// Configuration metadata
	ConfigVersion string    `gorm:"not null;default:'1.0'" json:"config_version"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	LastUpdated   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"last_updated"`

	// Relationships
	PlaybackSessions []PlaybackSession `gorm:"foreignKey:PlatformContext;references:PlatformType" json:"playback_sessions,omitempty"`

	// Audit fields
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// VideoCapabilities defines what video features the platform supports
type VideoCapabilities struct {
	SupportedFormats    []string `gorm:"type:jsonb" json:"supported_formats" validate:"required,min=1"`
	SupportedCodecs     []string `gorm:"type:jsonb" json:"supported_codecs" validate:"required,min=1"`
	MaxResolution       string   `json:"max_resolution" validate:"oneof=360p 480p 720p 1080p 4K 8K"`
	MaxFrameRate        int      `json:"max_frame_rate" validate:"gte=24,lte=120"`
	SupportsHDR         bool     `json:"supports_hdr"`
	HardwareAccelerated bool     `json:"hardware_accelerated"`
	SupportsVR          bool     `json:"supports_vr"`
	MaxBitrate          int      `json:"max_bitrate"` // kbps
}

// QualityConfiguration defines quality settings and adaptive streaming
type QualityConfiguration struct {
	AutoQualityEnabled   bool     `json:"auto_quality_enabled"`
	DefaultQuality       string   `json:"default_quality" validate:"oneof=auto 360p 480p 720p 1080p 4K"`
	AdaptiveStreaming    bool     `json:"adaptive_streaming"`
	QualityLadder        []string `gorm:"type:jsonb" json:"quality_ladder"`
	BufferConfiguration  BufferConfiguration `gorm:"embedded" json:"buffer_configuration"`
	BitrateConfiguration BitrateConfiguration `gorm:"embedded" json:"bitrate_configuration"`
}

// BufferConfiguration manages buffering behavior
type BufferConfiguration struct {
	InitialBufferSize    int     `json:"initial_buffer_size"`    // seconds
	MaxBufferSize        int     `json:"max_buffer_size"`        // seconds
	RebufferThreshold    float64 `json:"rebuffer_threshold"`     // seconds
	SeekBufferSize       int     `json:"seek_buffer_size"`       // seconds
	LiveBufferSize       int     `json:"live_buffer_size"`       // seconds
	BufferHealthMinimum  float64 `json:"buffer_health_minimum"`  // 0.0-1.0
}

// BitrateConfiguration manages adaptive bitrate settings
type BitrateConfiguration struct {
	StartupBitrate       int     `json:"startup_bitrate"`        // kbps
	MaxBitrate           int     `json:"max_bitrate"`            // kbps
	MinBitrate           int     `json:"min_bitrate"`            // kbps
	BitrateStep          int     `json:"bitrate_step"`           // kbps increase/decrease
	AdaptationThreshold  float64 `json:"adaptation_threshold"`   // bandwidth utilization threshold
	FastSwitchThreshold  float64 `json:"fast_switch_threshold"`  // quick adaptation threshold
}

// PerformanceConfiguration optimizes performance per platform
type PerformanceConfiguration struct {
	PreloadingEnabled     bool `json:"preloading_enabled"`
	PreloadDuration       int  `json:"preload_duration"`       // seconds
	MaxConcurrentStreams  int  `json:"max_concurrent_streams"`
	ThumbnailPreloading   bool `json:"thumbnail_preloading"`
	MetadataCaching       bool `json:"metadata_caching"`
	ProgressiveDownload   bool `json:"progressive_download"`
	ChunkSize             int  `json:"chunk_size"`             // bytes
	ParallelChunks        int  `json:"parallel_chunks"`
}

// NetworkOptimization handles network-specific optimizations
type NetworkOptimization struct {
	NetworkAwareQuality   bool                    `json:"network_aware_quality"`
	ConnectionTypeMapping map[string]string       `gorm:"type:jsonb" json:"connection_type_mapping,omitempty"`
	LatencyOptimization   LatencyOptimization     `gorm:"embedded" json:"latency_optimization"`
	BandwidthEstimation   BandwidthEstimation     `gorm:"embedded" json:"bandwidth_estimation"`
	OfflineConfiguration  OfflineConfiguration    `gorm:"embedded" json:"offline_configuration"`
}

// LatencyOptimization reduces startup and seeking delays
type LatencyOptimization struct {
	FastStart            bool `json:"fast_start"`
	LowLatencyMode       bool `json:"low_latency_mode"`
	PredictiveCaching    bool `json:"predictive_caching"`
	EdgeCaching          bool `json:"edge_caching"`
	DNSPrefetching       bool `json:"dns_prefetching"`
	ConnectionReuse      bool `json:"connection_reuse"`
}

// BandwidthEstimation manages adaptive streaming decisions
type BandwidthEstimation struct {
	EstimationMethod     string  `json:"estimation_method" validate:"oneof=throughput latency hybrid"`
	SampleWindow         int     `json:"sample_window"`         // seconds
	EstimationFrequency  int     `json:"estimation_frequency"`  // seconds
	SafetyMargin         float64 `json:"safety_margin"`         // percentage (0.0-1.0)
	ReactionTime         int     `json:"reaction_time"`         // milliseconds
}

// OfflineConfiguration manages offline video capabilities
type OfflineConfiguration struct {
	OfflineSupported     bool `json:"offline_supported"`
	MaxOfflineVideos     int  `json:"max_offline_videos"`
	MaxOfflineStorage    int  `json:"max_offline_storage"`    // MB
	OfflineQuality       string `json:"offline_quality" validate:"oneof=360p 480p 720p 1080p"`
	ExpirationDays       int  `json:"expiration_days"`
	BackgroundDownload   bool `json:"background_download"`
}

// UIConfiguration manages platform-specific UI behavior
type UIConfiguration struct {
	ControlsTimeout      int                     `json:"controls_timeout"`       // seconds
	GestureControls      bool                    `json:"gesture_controls"`
	PictureInPicture     bool                    `json:"picture_in_picture"`
	FullscreenBehavior   FullscreenBehavior      `gorm:"embedded" json:"fullscreen_behavior"`
	PlayerSkin           PlayerSkin              `gorm:"embedded" json:"player_skin"`
	MobileOptimizations  MobileOptimizations     `gorm:"embedded" json:"mobile_optimizations"`
}

// FullscreenBehavior controls fullscreen video behavior
type FullscreenBehavior struct {
	AutoFullscreen       bool   `json:"auto_fullscreen"`
	OrientationLock      bool   `json:"orientation_lock"`
	PreferredOrientation string `json:"preferred_orientation" validate:"oneof=portrait landscape auto"`
	ImmersiveMode        bool   `json:"immersive_mode"`
	StatusBarHiding      bool   `json:"status_bar_hiding"`
}

// PlayerSkin defines visual appearance and controls
type PlayerSkin struct {
	Theme                string    `json:"theme" validate:"oneof=light dark auto"`
	AccentColor          string    `json:"accent_color"`
	ControlsStyle        string    `json:"controls_style" validate:"oneof=modern classic minimal"`
	ShowProgressBar      bool      `json:"show_progress_bar"`
	ShowVolumeControl    bool      `json:"show_volume_control"`
	ShowPlaybackSpeed    bool      `json:"show_playback_speed"`
	ShowQualitySelector  bool      `json:"show_quality_selector"`
	CustomControls       []string  `gorm:"type:jsonb" json:"custom_controls,omitempty"`
}

// MobileOptimizations for mobile-specific features
type MobileOptimizations struct {
	TouchControls        bool `json:"touch_controls"`
	SwipeGestures        bool `json:"swipe_gestures"`
	PinchToZoom          bool `json:"pinch_to_zoom"`
	DoubleTapToSeek      bool `json:"double_tap_to_seek"`
	VolumeGestures       bool `json:"volume_gestures"`
	BrightnessGestures   bool `json:"brightness_gestures"`
	HapticFeedback       bool `json:"haptic_feedback"`
	BatteryOptimization  bool `json:"battery_optimization"`
}

// StorageConfiguration manages local storage and caching
type StorageConfiguration struct {
	CacheEnabled         bool `json:"cache_enabled"`
	MaxCacheSize         int  `json:"max_cache_size"`         // MB
	CacheEvictionPolicy  string `json:"cache_eviction_policy" validate:"oneof=lru lfu fifo"`
	MetadataCache        int  `json:"metadata_cache"`         // MB
	ThumbnailCache       int  `json:"thumbnail_cache"`        // MB
	VideoSegmentCache    int  `json:"video_segment_cache"`    // MB
	CacheEncryption      bool `json:"cache_encryption"`
	AutoCleanup          bool `json:"auto_cleanup"`
	CleanupFrequency     int  `json:"cleanup_frequency"`      // hours
}

// PlatformType represents supported platforms
type PlatformType string

const (
	PlatformWeb     PlatformType = "web"
	PlatformAndroid PlatformType = "android"
	PlatformIOS     PlatformType = "ios"
	PlatformDesktop PlatformType = "desktop"
	PlatformTV      PlatformType = "tv"
	PlatformConsole PlatformType = "console"
)

// Business logic methods

// GetOptimalQuality returns the best quality for current network conditions
func (pc *PlatformConfiguration) GetOptimalQuality(networkSpeed float64, deviceCapability string) string {
	if !pc.QualitySettings.AutoQualityEnabled {
		return pc.QualitySettings.DefaultQuality
	}

	// Simple quality selection based on network speed (Mbps)
	switch {
	case networkSpeed >= 25 && deviceCapability == "4K":
		return "4K"
	case networkSpeed >= 8 && deviceCapability >= "1080p":
		return "1080p"
	case networkSpeed >= 5:
		return "720p"
	case networkSpeed >= 2:
		return "480p"
	default:
		return "360p"
	}
}

// SupportsFormat checks if the platform supports a specific video format
func (pc *PlatformConfiguration) SupportsFormat(format string) bool {
	for _, supportedFormat := range pc.VideoCapabilities.SupportedFormats {
		if supportedFormat == format {
			return true
		}
	}
	return false
}

// SupportsCodec checks if the platform supports a specific codec
func (pc *PlatformConfiguration) SupportsCodec(codec string) bool {
	for _, supportedCodec := range pc.VideoCapabilities.SupportedCodecs {
		if supportedCodec == codec {
			return true
		}
	}
	return false
}

// CanPlayQuality determines if the platform can play the specified quality
func (pc *PlatformConfiguration) CanPlayQuality(quality string) bool {
	for _, availableQuality := range pc.QualitySettings.QualityLadder {
		if availableQuality == quality {
			return true
		}
	}
	return false
}

// GetBufferConfiguration returns platform-optimized buffer settings
func (pc *PlatformConfiguration) GetBufferConfiguration() BufferConfiguration {
	return pc.QualitySettings.BufferConfiguration
}

// ShouldUseHardwareAcceleration determines if hardware acceleration should be enabled
func (pc *PlatformConfiguration) ShouldUseHardwareAcceleration() bool {
	return pc.VideoCapabilities.HardwareAccelerated && pc.IsActive
}

// GetPreloadDuration returns optimal preload duration for the platform
func (pc *PlatformConfiguration) GetPreloadDuration() int {
	if !pc.PerformanceSettings.PreloadingEnabled {
		return 0
	}
	return pc.PerformanceSettings.PreloadDuration
}

// IsOfflineCapable checks if the platform supports offline video
func (pc *PlatformConfiguration) IsOfflineCapable() bool {
	return pc.NetworkOptimization.OfflineConfiguration.OfflineSupported
}

// GetStorageLimit returns maximum storage allowed for video content
func (pc *PlatformConfiguration) GetStorageLimit() int {
	if pc.IsOfflineCapable() {
		return pc.NetworkOptimization.OfflineConfiguration.MaxOfflineStorage
	}
	return pc.StorageConfiguration.MaxCacheSize
}

// ApplyNetworkOptimizations returns optimized settings based on network conditions
func (pc *PlatformConfiguration) ApplyNetworkOptimizations(connectionType string, bandwidth float64) map[string]interface{} {
	optimizations := make(map[string]interface{})

	// Apply connection-specific optimizations
	if mapping, exists := pc.NetworkOptimization.ConnectionTypeMapping[connectionType]; exists {
		optimizations["quality_override"] = mapping
	}

	// Apply bandwidth-based optimizations
	optimizations["estimated_bandwidth"] = bandwidth
	optimizations["buffer_health_minimum"] = pc.QualitySettings.BufferConfiguration.BufferHealthMinimum
	optimizations["preload_enabled"] = pc.PerformanceSettings.PreloadingEnabled

	// Low bandwidth optimizations
	if bandwidth < 1.0 { // Less than 1 Mbps
		optimizations["quality_override"] = "360p"
		optimizations["preload_enabled"] = false
		optimizations["parallel_chunks"] = 1
	}

	return optimizations
}

// ValidateConfiguration checks if the platform configuration is valid
func (pc *PlatformConfiguration) ValidateConfiguration() []string {
	var errors []string

	// Validate required fields
	if len(pc.VideoCapabilities.SupportedFormats) == 0 {
		errors = append(errors, "at least one supported format is required")
	}

	if len(pc.VideoCapabilities.SupportedCodecs) == 0 {
		errors = append(errors, "at least one supported codec is required")
	}

	if len(pc.QualitySettings.QualityLadder) == 0 {
		errors = append(errors, "quality ladder cannot be empty")
	}

	// Validate buffer configuration
	if pc.QualitySettings.BufferConfiguration.InitialBufferSize > pc.QualitySettings.BufferConfiguration.MaxBufferSize {
		errors = append(errors, "initial buffer size cannot exceed max buffer size")
	}

	// Validate bitrate configuration
	if pc.QualitySettings.BitrateConfiguration.MinBitrate > pc.QualitySettings.BitrateConfiguration.MaxBitrate {
		errors = append(errors, "minimum bitrate cannot exceed maximum bitrate")
	}

	return errors
}

// GORM hooks

// BeforeCreate sets up defaults before creating platform configuration record
func (pc *PlatformConfiguration) BeforeCreate(tx *gorm.DB) error {
	if pc.ID == uuid.Nil {
		pc.ID = uuid.New()
	}

	if pc.ConfigVersion == "" {
		pc.ConfigVersion = "1.0"
	}

	// Set default video capabilities if empty
	if len(pc.VideoCapabilities.SupportedFormats) == 0 {
		switch pc.PlatformType {
		case PlatformWeb:
			pc.VideoCapabilities.SupportedFormats = []string{"mp4", "webm", "hls", "dash"}
			pc.VideoCapabilities.SupportedCodecs = []string{"h264", "vp8", "vp9", "av1"}
		case PlatformAndroid:
			pc.VideoCapabilities.SupportedFormats = []string{"mp4", "webm", "hls"}
			pc.VideoCapabilities.SupportedCodecs = []string{"h264", "h265", "vp8", "vp9"}
		case PlatformIOS:
			pc.VideoCapabilities.SupportedFormats = []string{"mp4", "hls"}
			pc.VideoCapabilities.SupportedCodecs = []string{"h264", "h265"}
		}
	}

	// Set default quality ladder if empty
	if len(pc.QualitySettings.QualityLadder) == 0 {
		pc.QualitySettings.QualityLadder = []string{"360p", "480p", "720p", "1080p"}
	}

	return nil
}

// BeforeUpdate validates configuration before updating
func (pc *PlatformConfiguration) BeforeUpdate(tx *gorm.DB) error {
	errors := pc.ValidateConfiguration()
	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %v", errors)
	}
	return nil
}

// TableName returns the table name for GORM
func (PlatformConfiguration) TableName() string {
	return "platform_configurations"
}

// Database indexes for performance and querying
/*
CREATE INDEX CONCURRENTLY idx_platform_configurations_platform_type ON platform_configurations(platform_type);
CREATE INDEX CONCURRENTLY idx_platform_configurations_device_category ON platform_configurations(device_category);
CREATE INDEX CONCURRENTLY idx_platform_configurations_is_active ON platform_configurations(is_active);
CREATE INDEX CONCURRENTLY idx_platform_configurations_platform_version ON platform_configurations(platform_type, platform_version);
CREATE INDEX CONCURRENTLY idx_platform_configurations_config_version ON platform_configurations(config_version);
CREATE INDEX CONCURRENTLY idx_platform_configurations_last_updated ON platform_configurations(last_updated DESC);
*/