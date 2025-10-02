package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StreamProduct represents a product featured during a live stream
type StreamProduct struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StreamID               uuid.UUID      `gorm:"type:uuid;not null;index:idx_stream_id" json:"stream_id"`
	ProductID              uuid.UUID      `gorm:"type:uuid;not null;index:idx_product_id" json:"product_id"` // References commerce.products (foreign service)
	FeaturedAt             time.Time      `gorm:"not null;index:idx_featured_at" json:"featured_at"`
	DisplayDurationSeconds *int           `gorm:"" json:"display_duration_seconds,omitempty"`
	ViewCount              int            `gorm:"default:0" json:"view_count"`
	ClickCount             int            `gorm:"default:0" json:"click_count"`
	PurchaseCount          int            `gorm:"default:0" json:"purchase_count"`
	RevenueGenerated       float64        `gorm:"type:decimal(10,2);default:0.00" json:"revenue_generated"`
	DisplayPosition        string         `gorm:"type:varchar(20);default:'overlay'" json:"display_position"`
	DisplayPriority        int            `gorm:"default:0" json:"display_priority"`
	CreatedAt              time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName returns the table name for StreamProduct
func (StreamProduct) TableName() string {
	return "stream_products"
}