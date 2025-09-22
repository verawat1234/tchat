package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ContentItem represents a content item in the system
type ContentItem struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Category  string          `json:"category" gorm:"not null;index"`
	Type      ContentType     `json:"type" gorm:"type:varchar(50);not null"`
	Value     ContentValue    `json:"value" gorm:"type:jsonb;not null"`
	Metadata  ContentMetadata `json:"metadata" gorm:"type:jsonb"`
	Status    ContentStatus   `json:"status" gorm:"type:varchar(20);not null;default:'draft';index"`
	Tags      []string        `json:"tags,omitempty" gorm:"type:text[]"`
	Notes     *string         `json:"notes,omitempty"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ContentItem
func (ContentItem) TableName() string {
	return "content_items"
}

// ContentCategory represents a content category
type ContentCategory struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description *string   `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" gorm:"type:uuid"`
	Parent      *ContentCategory `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []ContentCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for ContentCategory
func (ContentCategory) TableName() string {
	return "content_categories"
}

// ContentVersion represents a version of content
type ContentVersion struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContentID uuid.UUID       `json:"content_id" gorm:"type:uuid;not null;index"`
	Content   ContentItem     `json:"content" gorm:"foreignKey:ContentID"`
	Version   int             `json:"version" gorm:"not null"`
	Value     ContentValue    `json:"value" gorm:"type:jsonb;not null"`
	Metadata  ContentMetadata `json:"metadata" gorm:"type:jsonb"`
	Status    ContentStatus   `json:"status" gorm:"type:varchar(20);not null"`
	CreatedBy *uuid.UUID      `json:"created_by,omitempty" gorm:"type:uuid"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for ContentVersion
func (ContentVersion) TableName() string {
	return "content_versions"
}

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeText         ContentType = "text"
	ContentTypeHTML         ContentType = "html"
	ContentTypeMarkdown     ContentType = "markdown"
	ContentTypeJSON         ContentType = "json"
	ContentTypeConfiguration ContentType = "configuration"
)

// ValidContentTypes returns all supported content types
func ValidContentTypes() []ContentType {
	return []ContentType{
		ContentTypeText,
		ContentTypeHTML,
		ContentTypeMarkdown,
		ContentTypeJSON,
		ContentTypeConfiguration,
	}
}

// IsValid validates if the content type is supported
func (ct ContentType) IsValid() bool {
	for _, valid := range ValidContentTypes() {
		if ct == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of ContentType
func (ct ContentType) String() string {
	return string(ct)
}

// Value implements the driver.Valuer interface for database storage
func (ct ContentType) Value() (driver.Value, error) {
	return string(ct), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (ct *ContentType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if str, ok := value.(string); ok {
		*ct = ContentType(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into ContentType", value)
}

// ContentStatus represents the status of content
type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusPublished ContentStatus = "published"
	ContentStatusArchived  ContentStatus = "archived"
	ContentStatusDeleted   ContentStatus = "deleted"
)

// ValidContentStatuses returns all supported content statuses
func ValidContentStatuses() []ContentStatus {
	return []ContentStatus{
		ContentStatusDraft,
		ContentStatusPublished,
		ContentStatusArchived,
		ContentStatusDeleted,
	}
}

// IsValid validates if the content status is supported
func (cs ContentStatus) IsValid() bool {
	for _, valid := range ValidContentStatuses() {
		if cs == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of ContentStatus
func (cs ContentStatus) String() string {
	return string(cs)
}

// Value implements the driver.Valuer interface for database storage
func (cs ContentStatus) Value() (driver.Value, error) {
	return string(cs), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (cs *ContentStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if str, ok := value.(string); ok {
		*cs = ContentStatus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into ContentStatus", value)
}

// ContentValue represents the value of content (flexible JSON structure)
type ContentValue map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (cv ContentValue) Value() (driver.Value, error) {
	if cv == nil {
		return nil, nil
	}
	return json.Marshal(cv)
}

// Scan implements the sql.Scanner interface for database retrieval
func (cv *ContentValue) Scan(value interface{}) error {
	if value == nil {
		*cv = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, cv)
}

// ContentMetadata represents metadata for content
type ContentMetadata map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (cm ContentMetadata) Value() (driver.Value, error) {
	if cm == nil {
		return nil, nil
	}
	return json.Marshal(cm)
}

// Scan implements the sql.Scanner interface for database retrieval
func (cm *ContentMetadata) Scan(value interface{}) error {
	if value == nil {
		*cm = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, cm)
}

// ContentFilters represents filters for querying content
type ContentFilters struct {
	Category    *string        `json:"category,omitempty"`
	Type        *ContentType   `json:"type,omitempty"`
	Status      *ContentStatus `json:"status,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
	Search      *string        `json:"search,omitempty"`
	CreatedFrom *time.Time     `json:"created_from,omitempty"`
	CreatedTo   *time.Time     `json:"created_to,omitempty"`
}

// Pagination represents pagination parameters
type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Offset   int `json:"offset"`
	Total    int64 `json:"total,omitempty"`
}

// SortOptions represents sorting options
type SortOptions struct {
	Field string `json:"field" form:"sort_field"`
	Order string `json:"order" form:"sort_order"` // "asc" or "desc"
}

// ContentResponse represents a paginated content response
type ContentResponse struct {
	Items      []ContentItem `json:"items"`
	Pagination Pagination   `json:"pagination"`
	Total      int64         `json:"total"`
}

// ContentCategoryResponse represents a paginated category response
type ContentCategoryResponse struct {
	Items      []ContentCategory `json:"items"`
	Pagination Pagination        `json:"pagination"`
	Total      int64             `json:"total"`
}

// ContentVersionResponse represents a paginated version response
type ContentVersionResponse struct {
	Items      []ContentVersion `json:"items"`
	Pagination Pagination       `json:"pagination"`
	Total      int64            `json:"total"`
}

// CreateContentRequest represents a request to create content
type CreateContentRequest struct {
	Category string          `json:"category" binding:"required"`
	Type     ContentType     `json:"type" binding:"required"`
	Value    ContentValue    `json:"value" binding:"required"`
	Metadata ContentMetadata `json:"metadata,omitempty"`
	Status   ContentStatus   `json:"status,omitempty"`
	Tags     []string        `json:"tags,omitempty"`
	Notes    *string         `json:"notes,omitempty"`
}

// UpdateContentRequest represents a request to update content
type UpdateContentRequest struct {
	Category *string          `json:"category,omitempty"`
	Type     *ContentType     `json:"type,omitempty"`
	Value    *ContentValue    `json:"value,omitempty"`
	Metadata *ContentMetadata `json:"metadata,omitempty"`
	Status   *ContentStatus   `json:"status,omitempty"`
	Tags     *[]string        `json:"tags,omitempty"`
	Notes    *string          `json:"notes,omitempty"`
}

// BulkUpdateRequest represents a request to bulk update content
type BulkUpdateRequest struct {
	IDs     []uuid.UUID        `json:"ids" binding:"required"`
	Updates UpdateContentRequest `json:"updates" binding:"required"`
}

// SyncContentRequest represents a request to sync content
type SyncContentRequest struct {
	LastSyncTime *time.Time `json:"last_sync_time,omitempty"`
	Categories   []string   `json:"categories,omitempty"`
}

// SyncContentResponse represents a response from content sync
type SyncContentResponse struct {
	Items        []ContentItem `json:"items"`
	DeletedIDs   []uuid.UUID   `json:"deleted_ids,omitempty"`
	SyncTime     time.Time     `json:"sync_time"`
	HasMore      bool          `json:"has_more"`
	NextSyncTime *time.Time    `json:"next_sync_time,omitempty"`
}