package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// MediaSubtab represents subcategories within main media categories
type MediaSubtab struct {
	ID             string          `gorm:"primaryKey;type:varchar(50)" json:"id" validate:"required,max=50"`
	CategoryID     string          `gorm:"type:varchar(50);not null" json:"categoryId" validate:"required"`
	Name           string          `gorm:"type:varchar(30);not null" json:"name" validate:"required,max=30"`
	DisplayOrder   int             `gorm:"not null" json:"displayOrder" validate:"required,gt=0"`
	FilterCriteria json.RawMessage `gorm:"type:jsonb;not null;default:'{}'" json:"filterCriteria"`
	IsActive       bool            `gorm:"not null;default:true" json:"isActive"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Category MediaCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

// TableName specifies the table name for MediaSubtab
func (MediaSubtab) TableName() string {
	return "media_subtabs"
}

// BeforeCreate hook to validate before creating
func (ms *MediaSubtab) BeforeCreate(tx *gorm.DB) error {
	return ms.validate()
}

// BeforeUpdate hook to validate before updating
func (ms *MediaSubtab) BeforeUpdate(tx *gorm.DB) error {
	return ms.validate()
}

// validate performs model validation
func (ms *MediaSubtab) validate() error {
	if ms.ID == "" {
		return gorm.ErrInvalidData
	}
	if ms.CategoryID == "" {
		return gorm.ErrInvalidData
	}
	if ms.Name == "" {
		return gorm.ErrInvalidData
	}
	if ms.DisplayOrder <= 0 {
		return gorm.ErrInvalidData
	}

	// Validate filter criteria is valid JSON
	if len(ms.FilterCriteria) > 0 {
		var temp interface{}
		if err := json.Unmarshal(ms.FilterCriteria, &temp); err != nil {
			return gorm.ErrInvalidData
		}
	}

	return nil
}

// FilterCriteriaMap represents the filter criteria as a map
type FilterCriteriaMap map[string]interface{}

// GetFilterCriteria returns the filter criteria as a map
func (ms *MediaSubtab) GetFilterCriteria() (FilterCriteriaMap, error) {
	var criteria FilterCriteriaMap
	if len(ms.FilterCriteria) == 0 {
		return criteria, nil
	}

	err := json.Unmarshal(ms.FilterCriteria, &criteria)
	return criteria, err
}

// SetFilterCriteria sets the filter criteria from a map
func (ms *MediaSubtab) SetFilterCriteria(criteria FilterCriteriaMap) error {
	if criteria == nil {
		ms.FilterCriteria = json.RawMessage("{}")
		return nil
	}

	data, err := json.Marshal(criteria)
	if err != nil {
		return err
	}

	ms.FilterCriteria = json.RawMessage(data)
	return nil
}

// MatchesContent checks if a content item matches this subtab's filter criteria
func (ms *MediaSubtab) MatchesContent(contentItem *MediaContentItem) bool {
	criteria, err := ms.GetFilterCriteria()
	if err != nil || len(criteria) == 0 {
		return true // No criteria means all content matches
	}

	// Check duration-based filtering (for movies)
	if maxDuration, exists := criteria["maxDuration"]; exists {
		if maxDur, ok := maxDuration.(float64); ok {
			if contentItem.Duration != nil && *contentItem.Duration > int(maxDur) {
				return false
			}
		}
	}

	if minDuration, exists := criteria["minDuration"]; exists {
		if minDur, ok := minDuration.(float64); ok {
			if contentItem.Duration == nil || *contentItem.Duration < int(minDur) {
				return false
			}
		}
	}

	// Add more filter criteria as needed
	return true
}

// IsValidID checks if the subtab ID is URL-safe
func (ms *MediaSubtab) IsValidID() bool {
	if len(ms.ID) == 0 || len(ms.ID) > 50 {
		return false
	}

	// Check for URL-safe characters (alphanumeric, hyphens, underscores)
	for _, char := range ms.ID {
		if !((char >= 'a' && char <= 'z') ||
			 (char >= 'A' && char <= 'Z') ||
			 (char >= '0' && char <= '9') ||
			 char == '-' || char == '_') {
			return false
		}
	}
	return true
}