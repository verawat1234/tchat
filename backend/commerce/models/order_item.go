package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderItem represents an item within an order
type OrderItem struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null" json:"orderId" validate:"required"`
	ProductID      uuid.UUID `gorm:"type:uuid;not null" json:"productId" validate:"required"`
	Quantity       int       `gorm:"not null" json:"quantity" validate:"required,gt=0"`
	Price          float64   `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gte=0"`
	Currency       string    `gorm:"type:varchar(3);not null;default:'USD'" json:"currency" validate:"required,len=3"`
	MediaLicense   *string   `gorm:"type:varchar(50);default:null" json:"mediaLicense,omitempty"`
	DownloadFormat *string   `gorm:"type:varchar(50);default:null" json:"downloadFormat,omitempty"`
	DownloadURL    *string   `gorm:"type:text;default:null" json:"downloadUrl,omitempty"`
	LicenseKey     *string   `gorm:"type:text;default:null" json:"licenseKey,omitempty"`
	ExpiryDate     *time.Time `gorm:"default:null" json:"expiryDate,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName specifies the table name for OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeCreate hook to validate before creating
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return oi.validate()
}

// BeforeUpdate hook to validate before updating
func (oi *OrderItem) BeforeUpdate(tx *gorm.DB) error {
	return oi.validate()
}

// validate performs model validation
func (oi *OrderItem) validate() error {
	if oi.OrderID == uuid.Nil {
		return gorm.ErrInvalidData
	}
	if oi.ProductID == uuid.Nil {
		return gorm.ErrInvalidData
	}
	if oi.Quantity <= 0 {
		return gorm.ErrInvalidData
	}
	if oi.Price < 0 {
		return gorm.ErrInvalidData
	}
	if oi.Currency == "" {
		return gorm.ErrInvalidData
	}
	return nil
}

// GetTotalPrice returns the total price for this order item (price * quantity)
func (oi *OrderItem) GetTotalPrice() float64 {
	return oi.Price * float64(oi.Quantity)
}

// IsMediaItem returns true if this order item is for a media product
func (oi *OrderItem) IsMediaItem() bool {
	return oi.Product.IsMedia()
}

// IsPhysicalItem returns true if this order item is for a physical product
func (oi *OrderItem) IsPhysicalItem() bool {
	return oi.Product.IsPhysical()
}

// HasDownloadURL returns true if a download URL is available
func (oi *OrderItem) HasDownloadURL() bool {
	return oi.DownloadURL != nil && *oi.DownloadURL != ""
}

// HasLicenseKey returns true if a license key is available
func (oi *OrderItem) HasLicenseKey() bool {
	return oi.LicenseKey != nil && *oi.LicenseKey != ""
}

// IsExpired returns true if the license has expired
func (oi *OrderItem) IsExpired() bool {
	if oi.ExpiryDate == nil {
		return false // No expiry date means never expires
	}
	return time.Now().After(*oi.ExpiryDate)
}

// GetDownloadURL returns the download URL or empty string if not set
func (oi *OrderItem) GetDownloadURL() string {
	if oi.DownloadURL == nil {
		return ""
	}
	return *oi.DownloadURL
}

// GetLicenseKey returns the license key or empty string if not set
func (oi *OrderItem) GetLicenseKey() string {
	if oi.LicenseKey == nil {
		return ""
	}
	return *oi.LicenseKey
}

// SetDownloadURL sets the download URL
func (oi *OrderItem) SetDownloadURL(url string) {
	if url == "" {
		oi.DownloadURL = nil
	} else {
		oi.DownloadURL = &url
	}
}

// SetLicenseKey sets the license key
func (oi *OrderItem) SetLicenseKey(key string) {
	if key == "" {
		oi.LicenseKey = nil
	} else {
		oi.LicenseKey = &key
	}
}

// SetExpiryDate sets the license expiry date
func (oi *OrderItem) SetExpiryDate(date *time.Time) {
	oi.ExpiryDate = date
}