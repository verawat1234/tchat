package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MediaContentItem represents individual media items within categories
type MediaContentItem struct {
	ID                 uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CategoryID         string          `gorm:"type:varchar(50);not null" json:"categoryId" validate:"required"`
	Title              string          `gorm:"type:varchar(200);not null" json:"title" validate:"required,max=200"`
	Description        string          `gorm:"type:text;not null" json:"description" validate:"required,max=1000"`
	ThumbnailURL       string          `gorm:"type:text;not null" json:"thumbnailUrl" validate:"required,url"`
	ContentType        string          `gorm:"type:varchar(20);not null" json:"contentType" validate:"required,oneof=book podcast cartoon short_movie long_movie"`
	Duration           *int            `gorm:"default:null" json:"duration,omitempty"` // in seconds, null for books
	Price              float64         `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gte=0"`
	Currency           string          `gorm:"type:varchar(3);not null;default:'USD'" json:"currency" validate:"required,len=3"`
	AvailabilityStatus string          `gorm:"type:varchar(20);not null;default:'available'" json:"availabilityStatus" validate:"required,oneof=available coming_soon unavailable"`
	IsFeatured         bool            `gorm:"not null;default:false" json:"isFeatured"`
	FeaturedOrder      *int            `gorm:"default:null" json:"featuredOrder,omitempty"`
	Metadata           json.RawMessage `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Category MediaCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

// TableName specifies the table name for MediaContentItem
func (MediaContentItem) TableName() string {
	return "media_content_items"
}

// BeforeCreate hook to validate before creating
func (mci *MediaContentItem) BeforeCreate(tx *gorm.DB) error {
	if mci.ID == uuid.Nil {
		mci.ID = uuid.New()
	}
	return mci.validate()
}

// BeforeUpdate hook to validate before updating
func (mci *MediaContentItem) BeforeUpdate(tx *gorm.DB) error {
	return mci.validate()
}

// validate performs model validation
func (mci *MediaContentItem) validate() error {
	if mci.CategoryID == "" {
		return gorm.ErrInvalidData
	}
	if mci.Title == "" {
		return gorm.ErrInvalidData
	}
	if mci.Description == "" {
		return gorm.ErrInvalidData
	}
	if mci.ThumbnailURL == "" {
		return gorm.ErrInvalidData
	}
	if mci.ContentType == "" {
		return gorm.ErrInvalidData
	}
	if mci.Price < 0 {
		return gorm.ErrInvalidData
	}
	if mci.Currency == "" {
		return gorm.ErrInvalidData
	}

	// Validate content type
	validContentTypes := map[string]bool{
		"book":        true,
		"podcast":     true,
		"cartoon":     true,
		"short_movie": true,
		"long_movie":  true,
	}
	if !validContentTypes[mci.ContentType] {
		return gorm.ErrInvalidData
	}

	// Validate availability status
	validStatuses := map[string]bool{
		"available":    true,
		"coming_soon":  true,
		"unavailable":  true,
	}
	if !validStatuses[mci.AvailabilityStatus] {
		return gorm.ErrInvalidData
	}

	// Duration validation - must be positive for non-book content
	if mci.ContentType != "book" {
		if mci.Duration == nil || *mci.Duration <= 0 {
			return gorm.ErrInvalidData
		}
	}

	// Featured order required when featured
	if mci.IsFeatured && mci.FeaturedOrder == nil {
		return gorm.ErrInvalidData
	}

	// Validate metadata is valid JSON
	if len(mci.Metadata) > 0 {
		var temp interface{}
		if err := json.Unmarshal(mci.Metadata, &temp); err != nil {
			return gorm.ErrInvalidData
		}
	}

	return nil
}

// MetadataMap represents the metadata as a map
type MetadataMap map[string]interface{}

// GetMetadata returns the metadata as a map
func (mci *MediaContentItem) GetMetadata() (MetadataMap, error) {
	var metadata MetadataMap
	if len(mci.Metadata) == 0 {
		return metadata, nil
	}

	err := json.Unmarshal(mci.Metadata, &metadata)
	return metadata, err
}

// SetMetadata sets the metadata from a map
func (mci *MediaContentItem) SetMetadata(metadata MetadataMap) error {
	if metadata == nil {
		mci.Metadata = json.RawMessage("{}")
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	mci.Metadata = json.RawMessage(data)
	return nil
}

// IsBook returns true if this is a book content item
func (mci *MediaContentItem) IsBook() bool {
	return mci.ContentType == "book"
}

// IsPodcast returns true if this is a podcast content item
func (mci *MediaContentItem) IsPodcast() bool {
	return mci.ContentType == "podcast"
}

// IsCartoon returns true if this is a cartoon content item
func (mci *MediaContentItem) IsCartoon() bool {
	return mci.ContentType == "cartoon"
}

// IsMovie returns true if this is a movie content item (short or long)
func (mci *MediaContentItem) IsMovie() bool {
	return mci.ContentType == "short_movie" || mci.ContentType == "long_movie"
}

// IsShortMovie returns true if this is a short movie (≤30 minutes)
func (mci *MediaContentItem) IsShortMovie() bool {
	return mci.ContentType == "short_movie"
}

// IsLongMovie returns true if this is a long movie (>30 minutes)
func (mci *MediaContentItem) IsLongMovie() bool {
	return mci.ContentType == "long_movie"
}

// GetDurationString returns a formatted duration string
func (mci *MediaContentItem) GetDurationString() string {
	if mci.Duration == nil {
		return ""
	}

	duration := *mci.Duration
	hours := duration / 3600
	minutes := (duration % 3600) / 60
	seconds := duration % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// IsAvailable returns true if the content is currently available
func (mci *MediaContentItem) IsAvailable() bool {
	return mci.AvailabilityStatus == "available"
}