package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductType represents the type of product
type ProductType string

const (
	ProductTypePhysical ProductType = "physical"
	ProductTypeMedia    ProductType = "media"
)

// Product represents a product in the store (physical or media)
type Product struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string          `gorm:"type:varchar(200);not null" json:"name" validate:"required,max=200"`
	Description     string          `gorm:"type:text;not null" json:"description" validate:"required"`
	Price           float64         `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gte=0"`
	Currency        string          `gorm:"type:varchar(3);not null;default:'USD'" json:"currency" validate:"required,len=3"`
	ProductType     ProductType     `gorm:"type:varchar(20);not null;default:'physical'" json:"productType" validate:"required,oneof=physical media"`
	MediaContentID  *uuid.UUID      `gorm:"type:uuid;default:null" json:"mediaContentId,omitempty"`
	MediaMetadata   json.RawMessage `gorm:"type:jsonb;default:null" json:"mediaMetadata,omitempty"`
	ThumbnailURL    string          `gorm:"type:text" json:"thumbnailUrl,omitempty"`
	IsActive        bool            `gorm:"not null;default:true" json:"isActive"`
	StockQuantity   *int            `gorm:"default:null" json:"stockQuantity,omitempty"` // null for digital/media products
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName specifies the table name for Product
func (Product) TableName() string {
	return "products"
}

// BeforeCreate hook to validate before creating
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return p.validate()
}

// BeforeUpdate hook to validate before updating
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	return p.validate()
}

// validate performs model validation
func (p *Product) validate() error {
	if p.Name == "" {
		return gorm.ErrInvalidData
	}
	if p.Description == "" {
		return gorm.ErrInvalidData
	}
	if p.Price < 0 {
		return gorm.ErrInvalidData
	}
	if p.Currency == "" {
		return gorm.ErrInvalidData
	}

	// Validate product type
	if p.ProductType != ProductTypePhysical && p.ProductType != ProductTypeMedia {
		return gorm.ErrInvalidData
	}

	// Media products must have media content ID
	if p.ProductType == ProductTypeMedia && p.MediaContentID == nil {
		return gorm.ErrInvalidData
	}

	// Physical products should have stock quantity
	if p.ProductType == ProductTypePhysical && p.StockQuantity == nil {
		return gorm.ErrInvalidData
	}

	// Validate media metadata is valid JSON if present
	if len(p.MediaMetadata) > 0 {
		var temp interface{}
		if err := json.Unmarshal(p.MediaMetadata, &temp); err != nil {
			return gorm.ErrInvalidData
		}
	}

	return nil
}

// MediaMetadataMap represents the media metadata as a map
type MediaMetadataMap map[string]interface{}

// GetMediaMetadata returns the media metadata as a map
func (p *Product) GetMediaMetadata() (MediaMetadataMap, error) {
	var metadata MediaMetadataMap
	if len(p.MediaMetadata) == 0 {
		return metadata, nil
	}

	err := json.Unmarshal(p.MediaMetadata, &metadata)
	return metadata, err
}

// SetMediaMetadata sets the media metadata from a map
func (p *Product) SetMediaMetadata(metadata MediaMetadataMap) error {
	if metadata == nil {
		p.MediaMetadata = nil
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	p.MediaMetadata = json.RawMessage(data)
	return nil
}

// IsPhysical returns true if this is a physical product
func (p *Product) IsPhysical() bool {
	return p.ProductType == ProductTypePhysical
}

// IsMedia returns true if this is a media product
func (p *Product) IsMedia() bool {
	return p.ProductType == ProductTypeMedia
}

// IsInStock returns true if the product is in stock (always true for media)
func (p *Product) IsInStock() bool {
	if p.IsMedia() {
		return true // Digital products are always "in stock"
	}
	return p.StockQuantity != nil && *p.StockQuantity > 0
}

// CanPurchase returns true if the product can be purchased
func (p *Product) CanPurchase() bool {
	return p.IsActive && p.IsInStock()
}