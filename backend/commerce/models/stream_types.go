package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// StreamCategory represents different types of Stream content categories
type StreamCategory struct {
	ID                      string          `gorm:"type:varchar(50);primaryKey" json:"id" validate:"required"`
	Name                    string          `gorm:"type:varchar(50);not null" json:"name" validate:"required,max=50"`
	DisplayOrder            int             `gorm:"not null" json:"displayOrder" validate:"gte=0"`
	IconName                string          `gorm:"type:varchar(100);not null" json:"iconName" validate:"required,max=100"`
	IsActive                bool            `gorm:"not null;default:true" json:"isActive"`
	FeaturedContentEnabled  bool            `gorm:"not null;default:true" json:"featuredContentEnabled"`
	CreatedAt               time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt               time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	ContentItems []StreamContentItem `gorm:"foreignKey:CategoryID" json:"contentItems,omitempty"`
	Subtabs      []StreamSubtab      `gorm:"foreignKey:CategoryID" json:"subtabs,omitempty"`
}

// TableName specifies the table name for StreamCategory
func (StreamCategory) TableName() string {
	return "stream_categories"
}

// BeforeCreate hook to validate before creating
func (sc *StreamCategory) BeforeCreate(tx *gorm.DB) error {
	return sc.validate()
}

// BeforeUpdate hook to validate before updating
func (sc *StreamCategory) BeforeUpdate(tx *gorm.DB) error {
	return sc.validate()
}

// validate performs model validation
func (sc *StreamCategory) validate() error {
	if sc.ID == "" {
		return gorm.ErrInvalidData
	}
	if sc.Name == "" {
		return gorm.ErrInvalidData
	}
	if sc.IconName == "" {
		return gorm.ErrInvalidData
	}
	if sc.DisplayOrder < 0 {
		return gorm.ErrInvalidData
	}
	return nil
}

// StreamSubtab represents subcategories within main media categories
type StreamSubtab struct {
	ID             string          `gorm:"type:varchar(50);primaryKey" json:"id" validate:"required"`
	CategoryID     string          `gorm:"type:varchar(50);not null" json:"categoryId" validate:"required"`
	Name           string          `gorm:"type:varchar(30);not null" json:"name" validate:"required,max=30"`
	DisplayOrder   int             `gorm:"not null" json:"displayOrder" validate:"gte=0"`
	FilterCriteria json.RawMessage `gorm:"type:jsonb;not null" json:"filterCriteria" validate:"required"`
	IsActive       bool            `gorm:"not null;default:true" json:"isActive"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Category StreamCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

// TableName specifies the table name for StreamSubtab
func (StreamSubtab) TableName() string {
	return "stream_subtabs"
}

// BeforeCreate hook to validate before creating
func (ss *StreamSubtab) BeforeCreate(tx *gorm.DB) error {
	return ss.validate()
}

// BeforeUpdate hook to validate before updating
func (ss *StreamSubtab) BeforeUpdate(tx *gorm.DB) error {
	return ss.validate()
}

// validate performs model validation
func (ss *StreamSubtab) validate() error {
	if ss.ID == "" {
		return gorm.ErrInvalidData
	}
	if ss.CategoryID == "" {
		return gorm.ErrInvalidData
	}
	if ss.Name == "" {
		return gorm.ErrInvalidData
	}
	if ss.DisplayOrder < 0 {
		return gorm.ErrInvalidData
	}

	// Validate filter criteria is valid JSON
	if len(ss.FilterCriteria) > 0 {
		var temp interface{}
		if err := json.Unmarshal(ss.FilterCriteria, &temp); err != nil {
			return gorm.ErrInvalidData
		}
	}

	return nil
}

// GetFilterCriteria returns the filter criteria as a map
func (ss *StreamSubtab) GetFilterCriteria() (map[string]interface{}, error) {
	var criteria map[string]interface{}
	if len(ss.FilterCriteria) == 0 {
		return criteria, nil
	}

	err := json.Unmarshal(ss.FilterCriteria, &criteria)
	return criteria, err
}

// SetFilterCriteria sets the filter criteria from a map
func (ss *StreamSubtab) SetFilterCriteria(criteria map[string]interface{}) error {
	if criteria == nil {
		ss.FilterCriteria = nil
		return nil
	}

	data, err := json.Marshal(criteria)
	if err != nil {
		return err
	}

	ss.FilterCriteria = json.RawMessage(data)
	return nil
}

// StreamContentType represents the type of stream content
type StreamContentType string

const (
	StreamContentTypeBook       StreamContentType = "book"
	StreamContentTypePodcast    StreamContentType = "podcast"
	StreamContentTypeCartoon    StreamContentType = "cartoon"
	StreamContentTypeShortMovie StreamContentType = "short_movie"
	StreamContentTypeLongMovie  StreamContentType = "long_movie"
	StreamContentTypeMusic      StreamContentType = "music"
	StreamContentTypeArt        StreamContentType = "art"
)

// StreamAvailabilityStatus represents the availability status of content
type StreamAvailabilityStatus string

const (
	StreamAvailabilityAvailable   StreamAvailabilityStatus = "available"
	StreamAvailabilityComingSoon  StreamAvailabilityStatus = "coming_soon"
	StreamAvailabilityUnavailable StreamAvailabilityStatus = "unavailable"
)

// StreamContentItem represents individual media items within categories
type StreamContentItem struct {
	ID                 uuid.UUID                 `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CategoryID         string                    `gorm:"type:varchar(50);not null" json:"categoryId" validate:"required"`
	Title              string                    `gorm:"type:varchar(200);not null" json:"title" validate:"required,max=200"`
	Description        string                    `gorm:"type:text;not null" json:"description" validate:"required"`
	ThumbnailURL       string                    `gorm:"type:text;not null" json:"thumbnailUrl" validate:"required"`
	ContentType        StreamContentType         `gorm:"type:varchar(20);not null" json:"contentType" validate:"required"`
	Duration           *int                      `gorm:"default:null" json:"duration,omitempty"` // in seconds, null for books
	Price              decimal.Decimal           `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gte=0"`
	Currency           string                    `gorm:"type:varchar(3);not null;default:'USD'" json:"currency" validate:"required,len=3"`
	AvailabilityStatus StreamAvailabilityStatus  `gorm:"type:varchar(20);not null;default:'available'" json:"availabilityStatus" validate:"required"`
	IsFeatured         bool                      `gorm:"not null;default:false" json:"isFeatured"`
	FeaturedOrder      *int                      `gorm:"default:null" json:"featuredOrder,omitempty"`
	Metadata           json.RawMessage           `gorm:"type:jsonb;default:null" json:"metadata,omitempty"`
	CreatedAt          time.Time                 `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time                 `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Category StreamCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

// TableName specifies the table name for StreamContentItem
func (StreamContentItem) TableName() string {
	return "stream_content_items"
}

// BeforeCreate hook to validate before creating
func (sci *StreamContentItem) BeforeCreate(tx *gorm.DB) error {
	if sci.ID == uuid.Nil {
		sci.ID = uuid.New()
	}
	return sci.validate()
}

// BeforeUpdate hook to validate before updating
func (sci *StreamContentItem) BeforeUpdate(tx *gorm.DB) error {
	return sci.validate()
}

// validate performs model validation
func (sci *StreamContentItem) validate() error {
	if sci.CategoryID == "" {
		return gorm.ErrInvalidData
	}
	if sci.Title == "" {
		return gorm.ErrInvalidData
	}
	if sci.Description == "" {
		return gorm.ErrInvalidData
	}
	if sci.ThumbnailURL == "" {
		return gorm.ErrInvalidData
	}
	if sci.Currency == "" {
		return gorm.ErrInvalidData
	}

	// Validate content type
	validContentTypes := []StreamContentType{
		StreamContentTypeBook,
		StreamContentTypePodcast,
		StreamContentTypeCartoon,
		StreamContentTypeShortMovie,
		StreamContentTypeLongMovie,
		StreamContentTypeMusic,
		StreamContentTypeArt,
	}

	isValidType := false
	for _, validType := range validContentTypes {
		if sci.ContentType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return gorm.ErrInvalidData
	}

	// Validate availability status
	validStatuses := []StreamAvailabilityStatus{
		StreamAvailabilityAvailable,
		StreamAvailabilityComingSoon,
		StreamAvailabilityUnavailable,
	}

	isValidStatus := false
	for _, validStatus := range validStatuses {
		if sci.AvailabilityStatus == validStatus {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return gorm.ErrInvalidData
	}

	// Duration must be > 0 for non-book content
	if sci.ContentType != StreamContentTypeBook && (sci.Duration == nil || *sci.Duration <= 0) {
		return gorm.ErrInvalidData
	}

	// Featured order required when isFeatured is true
	if sci.IsFeatured && sci.FeaturedOrder == nil {
		return gorm.ErrInvalidData
	}

	// Validate metadata is valid JSON if present
	if len(sci.Metadata) > 0 {
		var temp interface{}
		if err := json.Unmarshal(sci.Metadata, &temp); err != nil {
			return gorm.ErrInvalidData
		}
	}

	return nil
}

// GetMetadata returns the metadata as a map
func (sci *StreamContentItem) GetMetadata() (map[string]interface{}, error) {
	var metadata map[string]interface{}
	if len(sci.Metadata) == 0 {
		return metadata, nil
	}

	err := json.Unmarshal(sci.Metadata, &metadata)
	return metadata, err
}

// SetMetadata sets the metadata from a map
func (sci *StreamContentItem) SetMetadata(metadata map[string]interface{}) error {
	if metadata == nil {
		sci.Metadata = nil
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	sci.Metadata = json.RawMessage(data)
	return nil
}

// IsBook returns true if this is a book content item
func (sci *StreamContentItem) IsBook() bool {
	return sci.ContentType == StreamContentTypeBook
}

// IsVideo returns true if this is a video content item
func (sci *StreamContentItem) IsVideo() bool {
	return sci.ContentType == StreamContentTypeShortMovie ||
		   sci.ContentType == StreamContentTypeLongMovie ||
		   sci.ContentType == StreamContentTypeCartoon
}

// IsAudio returns true if this is an audio content item
func (sci *StreamContentItem) IsAudio() bool {
	return sci.ContentType == StreamContentTypePodcast ||
		   sci.ContentType == StreamContentTypeMusic
}

// IsAvailable returns true if the content is available for purchase
func (sci *StreamContentItem) IsAvailable() bool {
	return sci.AvailabilityStatus == StreamAvailabilityAvailable
}

// CanPurchase returns true if the content can be purchased
func (sci *StreamContentItem) CanPurchase() bool {
	return sci.IsAvailable()
}

// GetDurationString returns a human-readable duration string
func (sci *StreamContentItem) GetDurationString() string {
	if sci.Duration == nil {
		return ""
	}

	duration := *sci.Duration
	hours := duration / 3600
	minutes := (duration % 3600) / 60
	seconds := duration % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// ContentResponse represents paginated content response
type ContentResponse struct {
	Items   []StreamContentItem `json:"items"`
	Total   int64               `json:"total"`
	HasMore bool                `json:"hasMore"`
}

// FeaturedResponse represents featured content response
type FeaturedResponse struct {
	Items   []StreamContentItem `json:"items"`
	Total   int                 `json:"total"`
	HasMore bool                `json:"hasMore"`
}

// StreamPurchaseRequest represents a content purchase request
type StreamPurchaseRequest struct {
	MediaContentID string `json:"mediaContentId" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	MediaLicense   string `json:"mediaLicense" binding:"required"`
	DownloadFormat string `json:"downloadFormat,omitempty"`
	CartID         string `json:"cartId,omitempty"`
}

// StreamPurchaseResponse represents a content purchase response
type StreamPurchaseResponse struct {
	OrderID     string  `json:"orderId"`
	TotalAmount float64 `json:"totalAmount"`
	Currency    string  `json:"currency"`
	Success     bool    `json:"success"`
	Message     string  `json:"message"`
}