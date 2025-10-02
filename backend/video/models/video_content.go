// backend/video/models/video_content.go
// Video Content model - Primary entity for user-uploaded videos
// Implements T025: Video Content model with comprehensive metadata support

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

// VideoContent represents the primary video entity with comprehensive metadata
type VideoContent struct {
	// Primary identifier
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Basic metadata
	Title       string `gorm:"not null;size:200" json:"title" validate:"required,min=1,max=200"`
	Description string `gorm:"type:text" json:"description"`
	Duration    int    `gorm:"not null" json:"duration" validate:"required,gt=0"`

	// Quality and format specifications
	QualityOptions      []string             `gorm:"type:jsonb" json:"quality_options" validate:"required,min=1"`
	FormatSpecification FormatSpecification  `gorm:"embedded" json:"format_specifications"`

	// Timing information
	UploadTimestamp time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"upload_timestamp"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Availability and monetization
	AvailabilityStatus AvailabilityStatus    `gorm:"type:varchar(20);not null;default:'public'" json:"availability_status"`
	PricingInformation *PricingInformation   `gorm:"embedded" json:"pricing_information,omitempty"`

	// Social engagement metrics
	SocialMetrics SocialEngagementMetrics `gorm:"embedded" json:"social_engagement_metrics"`

	// Creator and content information
	CreatorID    uuid.UUID `gorm:"type:uuid;not null;index" json:"creator_id" validate:"required"`
	FileURL      string    `gorm:"not null" json:"file_url" validate:"required,url"`
	ThumbnailURL string    `json:"thumbnail_url" validate:"url"`

	// Searchable content
	Tags          []string    `gorm:"type:jsonb" json:"tags"`
	ContentRating ContentRating `gorm:"type:varchar(10);default:'G'" json:"content_rating"`

	// Processing status
	ProcessingStatus ProcessingStatus `gorm:"type:varchar(20);default:'pending'" json:"processing_status"`
	ProcessingLog    string          `gorm:"type:text" json:"processing_log,omitempty"`

	// Relationships will be loaded via preloading
	PlaybackSessions []PlaybackSession `gorm:"foreignKey:VideoID" json:"playback_sessions,omitempty"`
	ViewingHistory   []ViewingHistory  `gorm:"foreignKey:VideoID" json:"viewing_history,omitempty"`
}

// FormatSpecification contains video format details
type FormatSpecification struct {
	Codec     string `json:"codec"`
	Container string `json:"container"`
	Bitrate   int    `json:"bitrate"` // in kbps
	FPS       int    `json:"fps"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// PricingInformation contains monetization details
type PricingInformation struct {
	Price            float64         `json:"price" validate:"gte=0"`
	Currency         string          `json:"currency" validate:"required,len=3"`
	MonetizationType MonetizationType `json:"monetization_type"`
	RegionalPricing  map[string]float64 `gorm:"type:jsonb" json:"regional_pricing,omitempty"`
}

// SocialEngagementMetrics tracks social interaction counts
type SocialEngagementMetrics struct {
	ViewCount    int64 `json:"view_count" gorm:"default:0"`
	LikeCount    int64 `json:"like_count" gorm:"default:0"`
	CommentCount int64 `json:"comment_count" gorm:"default:0"`
	ShareCount   int64 `json:"share_count" gorm:"default:0"`
	BookmarkCount int64 `json:"bookmark_count" gorm:"default:0"`
}

// Enums for video status and classification

type AvailabilityStatus string

const (
	StatusPublic     AvailabilityStatus = "public"
	StatusPrivate    AvailabilityStatus = "private"
	StatusUnlisted   AvailabilityStatus = "unlisted"
	StatusMonetized  AvailabilityStatus = "monetized"
	StatusArchived   AvailabilityStatus = "archived"
	StatusDeleted    AvailabilityStatus = "deleted"
)

type ContentRating string

const (
	RatingG    ContentRating = "G"        // General audiences
	RatingPG   ContentRating = "PG"       // Parental guidance
	RatingPG13 ContentRating = "PG-13"    // Parents strongly cautioned
	RatingR    ContentRating = "R"        // Restricted
	RatingNC17 ContentRating = "NC-17"    // Adults only
)

type MonetizationType string

const (
	MonetizationFree        MonetizationType = "free"
	MonetizationPayPerView  MonetizationType = "pay_per_view"
	MonetizationSubscription MonetizationType = "subscription"
	MonetizationRental      MonetizationType = "rental"
	MonetizationPurchase    MonetizationType = "purchase"
)

type ProcessingStatus string

const (
	ProcessingPending   ProcessingStatus = "pending"
	ProcessingQueued    ProcessingStatus = "queued"
	ProcessingActive    ProcessingStatus = "processing"
	ProcessingCompleted ProcessingStatus = "completed"
	ProcessingFailed    ProcessingStatus = "failed"
	ProcessingCanceled  ProcessingStatus = "canceled"
)

// Business logic methods

// IsAccessible checks if video is accessible to users
func (v *VideoContent) IsAccessible() bool {
	return v.AvailabilityStatus == StatusPublic ||
		   v.AvailabilityStatus == StatusUnlisted ||
		   v.AvailabilityStatus == StatusMonetized
}

// IsProcessingComplete checks if video processing is finished
func (v *VideoContent) IsProcessingComplete() bool {
	return v.ProcessingStatus == ProcessingCompleted
}

// GetQualityURL returns streaming URL for specific quality
func (v *VideoContent) GetQualityURL(quality string) string {
	// This would integrate with CDN/streaming service
	// For now, return base URL with quality parameter
	return v.FileURL + "?quality=" + quality
}

// UpdateSocialMetrics updates engagement counters
func (v *VideoContent) UpdateSocialMetrics(metricType string, delta int64) {
	switch metricType {
	case "view":
		v.SocialMetrics.ViewCount += delta
	case "like":
		v.SocialMetrics.LikeCount += delta
	case "comment":
		v.SocialMetrics.CommentCount += delta
	case "share":
		v.SocialMetrics.ShareCount += delta
	case "bookmark":
		v.SocialMetrics.BookmarkCount += delta
	}
}

// IsMonetized checks if video has pricing configuration
func (v *VideoContent) IsMonetized() bool {
	return v.AvailabilityStatus == StatusMonetized && v.PricingInformation != nil
}

// GetPrice returns the price for a specific currency/region
func (v *VideoContent) GetPrice(currency string) float64 {
	if v.PricingInformation == nil {
		return 0
	}

	// Check regional pricing first
	if regionalPrice, exists := v.PricingInformation.RegionalPricing[currency]; exists {
		return regionalPrice
	}

	// Fall back to default price
	return v.PricingInformation.Price
}

// Validation methods

// ValidateAvailabilityTransition checks if status transition is allowed
func (v *VideoContent) ValidateAvailabilityTransition(newStatus AvailabilityStatus) error {
	validTransitions := map[AvailabilityStatus][]AvailabilityStatus{
		StatusPublic:    {StatusPrivate, StatusUnlisted, StatusMonetized, StatusArchived, StatusDeleted},
		StatusPrivate:   {StatusPublic, StatusUnlisted, StatusMonetized, StatusDeleted},
		StatusUnlisted:  {StatusPublic, StatusPrivate, StatusMonetized, StatusArchived, StatusDeleted},
		StatusMonetized: {StatusPublic, StatusPrivate, StatusUnlisted, StatusArchived, StatusDeleted},
		StatusArchived:  {StatusPublic, StatusPrivate, StatusUnlisted, StatusMonetized, StatusDeleted},
		StatusDeleted:   {}, // No transitions from deleted state
	}

	allowedTransitions := validTransitions[v.AvailabilityStatus]
	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return nil
		}
	}

	return ErrInvalidStatusTransition
}

// Custom errors
var (
	ErrInvalidStatusTransition = fmt.Errorf("invalid availability status transition")
	ErrVideoNotProcessed      = fmt.Errorf("video processing not completed")
	ErrVideoNotAccessible     = fmt.Errorf("video is not accessible")
)

// GORM hooks

// BeforeCreate sets up defaults before creating video record
func (v *VideoContent) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}

	if v.AvailabilityStatus == "" {
		v.AvailabilityStatus = StatusPublic
	}

	if v.ProcessingStatus == "" {
		v.ProcessingStatus = ProcessingPending
	}

	if v.ContentRating == "" {
		v.ContentRating = RatingG
	}

	return nil
}

// TableName returns the table name for GORM
func (VideoContent) TableName() string {
	return "video_contents"
}

// Database indexes for performance
// These would be created via migrations
/*
CREATE INDEX CONCURRENTLY idx_video_contents_creator_id ON video_contents(creator_id);
CREATE INDEX CONCURRENTLY idx_video_contents_availability_status ON video_contents(availability_status);
CREATE INDEX CONCURRENTLY idx_video_contents_upload_timestamp ON video_contents(upload_timestamp DESC);
CREATE INDEX CONCURRENTLY idx_video_contents_processing_status ON video_contents(processing_status);
CREATE INDEX CONCURRENTLY idx_video_contents_content_rating ON video_contents(content_rating);
CREATE INDEX CONCURRENTLY idx_video_contents_tags ON video_contents USING GIN(tags);
*/