package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MediaDeliveryStatus represents the status of media delivery
type MediaDeliveryStatus string

const (
	MediaDeliveryPending   MediaDeliveryStatus = "pending"
	MediaDeliveryDelivered MediaDeliveryStatus = "delivered"
	MediaDeliveryFailed    MediaDeliveryStatus = "failed"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents an order containing items (physical and/or media)
type Order struct {
	ID                  uuid.UUID           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderNumber         string              `gorm:"type:varchar(50);not null;uniqueIndex" json:"orderNumber" validate:"required"`
	UserID              uuid.UUID           `gorm:"type:uuid;not null" json:"userId" validate:"required"`
	Status              OrderStatus         `gorm:"type:varchar(20);not null;default:'pending'" json:"status" validate:"required"`
	TotalAmount         float64             `gorm:"type:decimal(10,2);not null" json:"totalAmount" validate:"required,gte=0"`
	Currency            string              `gorm:"type:varchar(3);not null;default:'USD'" json:"currency" validate:"required,len=3"`
	ContainsMediaItems  bool                `gorm:"not null;default:false" json:"containsMediaItems"`
	MediaDeliveryStatus MediaDeliveryStatus `gorm:"type:varchar(20);default:'pending'" json:"mediaDeliveryStatus,omitempty"`
	ShippingAddress     *string             `gorm:"type:text;default:null" json:"shippingAddress,omitempty"`
	BillingAddress      *string             `gorm:"type:text;default:null" json:"billingAddress,omitempty"`
	PaymentMethod       string              `gorm:"type:varchar(50);not null" json:"paymentMethod" validate:"required"`
	PaymentStatus       string              `gorm:"type:varchar(20);not null;default:'pending'" json:"paymentStatus" validate:"required"`
	CreatedAt           time.Time           `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time           `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate hook to validate before creating
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.OrderNumber == "" {
		o.OrderNumber = o.generateOrderNumber()
	}
	return o.validate()
}

// BeforeUpdate hook to validate before updating
func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	return o.validate()
}

// validate performs model validation
func (o *Order) validate() error {
	if o.UserID == uuid.Nil {
		return gorm.ErrInvalidData
	}
	if o.OrderNumber == "" {
		return gorm.ErrInvalidData
	}
	if o.TotalAmount < 0 {
		return gorm.ErrInvalidData
	}
	if o.Currency == "" {
		return gorm.ErrInvalidData
	}
	if o.PaymentMethod == "" {
		return gorm.ErrInvalidData
	}
	if o.PaymentStatus == "" {
		return gorm.ErrInvalidData
	}

	// Validate status
	validStatuses := map[OrderStatus]bool{
		OrderStatusPending:   true,
		OrderStatusConfirmed: true,
		OrderStatusShipped:   true,
		OrderStatusDelivered: true,
		OrderStatusCancelled: true,
	}
	if !validStatuses[o.Status] {
		return gorm.ErrInvalidData
	}

	// Validate media delivery status if contains media items
	if o.ContainsMediaItems {
		validMediaStatuses := map[MediaDeliveryStatus]bool{
			MediaDeliveryPending:   true,
			MediaDeliveryDelivered: true,
			MediaDeliveryFailed:    true,
		}
		if !validMediaStatuses[o.MediaDeliveryStatus] {
			return gorm.ErrInvalidData
		}
	}

	return nil
}

// generateOrderNumber generates a unique order number
func (o *Order) generateOrderNumber() string {
	// Generate order number based on timestamp and random component
	return fmt.Sprintf("ORD-%d-%s", time.Now().Unix(), uuid.New().String()[:8])
}

// GetMediaItems returns only the media items in this order
func (o *Order) GetMediaItems() []OrderItem {
	var mediaItems []OrderItem
	for _, item := range o.Items {
		if item.IsMediaItem() {
			mediaItems = append(mediaItems, item)
		}
	}
	return mediaItems
}

// GetPhysicalItems returns only the physical items in this order
func (o *Order) GetPhysicalItems() []OrderItem {
	var physicalItems []OrderItem
	for _, item := range o.Items {
		if item.IsPhysicalItem() {
			physicalItems = append(physicalItems, item)
		}
	}
	return physicalItems
}

// GetTotalMediaAmount returns the total amount for media items only
func (o *Order) GetTotalMediaAmount() float64 {
	total := 0.0
	for _, item := range o.GetMediaItems() {
		total += item.GetTotalPrice()
	}
	return total
}

// GetTotalPhysicalAmount returns the total amount for physical items only
func (o *Order) GetTotalPhysicalAmount() float64 {
	total := 0.0
	for _, item := range o.GetPhysicalItems() {
		total += item.GetTotalPrice()
	}
	return total
}

// HasMediaItems returns true if the order contains any media items
func (o *Order) HasMediaItems() bool {
	return o.ContainsMediaItems && len(o.GetMediaItems()) > 0
}

// HasPhysicalItems returns true if the order contains any physical items
func (o *Order) HasPhysicalItems() bool {
	return len(o.GetPhysicalItems()) > 0
}

// IsPending returns true if the order is pending
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsConfirmed returns true if the order is confirmed
func (o *Order) IsConfirmed() bool {
	return o.Status == OrderStatusConfirmed
}

// IsDelivered returns true if the order is delivered
func (o *Order) IsDelivered() bool {
	return o.Status == OrderStatusDelivered
}

// IsCancelled returns true if the order is cancelled
func (o *Order) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

// IsMediaDelivered returns true if media content has been delivered
func (o *Order) IsMediaDelivered() bool {
	return o.MediaDeliveryStatus == MediaDeliveryDelivered
}