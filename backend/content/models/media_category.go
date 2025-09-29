package models

import (
	"time"

	"gorm.io/gorm"
)

// MediaCategory represents different types of media content categories
type MediaCategory struct {
	ID                     string    `gorm:"primaryKey;type:varchar(50)" json:"id" validate:"required,max=50"`
	Name                   string    `gorm:"type:varchar(50);not null" json:"name" validate:"required,max=50"`
	DisplayOrder           int       `gorm:"not null" json:"displayOrder" validate:"required,gt=0"`
	IconName               string    `gorm:"type:varchar(50);not null" json:"iconName" validate:"required,max=50"`
	IsActive               bool      `gorm:"not null;default:true" json:"isActive"`
	FeaturedContentEnabled bool      `gorm:"not null;default:true" json:"featuredContentEnabled"`
	CreatedAt              time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Subtabs      []MediaSubtab      `gorm:"foreignKey:CategoryID" json:"subtabs,omitempty"`
	ContentItems []MediaContentItem `gorm:"foreignKey:CategoryID" json:"contentItems,omitempty"`
}

// TableName specifies the table name for MediaCategory
func (MediaCategory) TableName() string {
	return "media_categories"
}

// BeforeCreate hook to validate before creating
func (mc *MediaCategory) BeforeCreate(tx *gorm.DB) error {
	return mc.validate()
}

// BeforeUpdate hook to validate before updating
func (mc *MediaCategory) BeforeUpdate(tx *gorm.DB) error {
	return mc.validate()
}

// validate performs model validation
func (mc *MediaCategory) validate() error {
	if mc.ID == "" {
		return gorm.ErrInvalidData
	}
	if mc.Name == "" {
		return gorm.ErrInvalidData
	}
	if mc.DisplayOrder <= 0 {
		return gorm.ErrInvalidData
	}
	if mc.IconName == "" {
		return gorm.ErrInvalidData
	}
	return nil
}

// IsValidID checks if the category ID is URL-safe
func (mc *MediaCategory) IsValidID() bool {
	if len(mc.ID) == 0 || len(mc.ID) > 50 {
		return false
	}

	// Check for URL-safe characters (alphanumeric, hyphens, underscores)
	for _, char := range mc.ID {
		if !((char >= 'a' && char <= 'z') ||
			 (char >= 'A' && char <= 'Z') ||
			 (char >= '0' && char <= '9') ||
			 char == '-' || char == '_') {
			return false
		}
	}
	return true
}

// HasSubtabs returns true if this category has subtabs configured
func (mc *MediaCategory) HasSubtabs() bool {
	return len(mc.Subtabs) > 0
}

// GetActiveSubtabs returns only active subtabs ordered by display order
func (mc *MediaCategory) GetActiveSubtabs() []MediaSubtab {
	var activeSubtabs []MediaSubtab
	for _, subtab := range mc.Subtabs {
		if subtab.IsActive {
			activeSubtabs = append(activeSubtabs, subtab)
		}
	}
	return activeSubtabs
}