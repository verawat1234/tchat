package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// CartStatus represents the status of a cart
type CartStatus string

const (
	CartStatusActive    CartStatus = "active"
	CartStatusAbandoned CartStatus = "abandoned"
	CartStatusConverted CartStatus = "converted"
	CartStatusExpired   CartStatus = "expired"
)

// CartItemType represents the type of cart item
type CartItemType string

const (
	CartItemTypeProduct CartItemType = "product"
	CartItemTypeService CartItemType = "service"
	CartItemTypeDigital CartItemType = "digital"
)

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID              uuid.UUID       `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	ProductID       uuid.UUID       `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	VariantID       *uuid.UUID      `json:"variant_id,omitempty" gorm:"column:variant_id;type:uuid"`
	BusinessID      uuid.UUID       `json:"business_id" gorm:"column:business_id;type:uuid;not null"`
	Type            CartItemType    `json:"type" gorm:"column:type;type:varchar(20);not null;default:'product'"`

	// Quantity and pricing
	Quantity        int             `json:"quantity" gorm:"column:quantity;not null"`
	UnitPrice       decimal.Decimal `json:"unit_price" gorm:"column:unit_price;type:decimal(20,8);not null"`
	TotalPrice      decimal.Decimal `json:"total_price" gorm:"column:total_price;type:decimal(20,8);not null"`
	Currency        string          `json:"currency" gorm:"column:currency;size:3;not null"`

	// Discounts and taxes
	DiscountAmount  decimal.Decimal `json:"discount_amount" gorm:"column:discount_amount;type:decimal(20,8);default:0"`
	TaxAmount       decimal.Decimal `json:"tax_amount" gorm:"column:tax_amount;type:decimal(20,8);default:0"`

	// Cached product info for performance
	ProductName     string          `json:"product_name" gorm:"column:product_name;size:200;not null"`
	ProductSKU      string          `json:"product_sku,omitempty" gorm:"column:product_sku;size:100"`
	ProductImage    string          `json:"product_image,omitempty" gorm:"column:product_image;size:500"`
	BusinessName    string          `json:"business_name" gorm:"column:business_name;size:100;not null"`

	// Variant info
	VariantName     string          `json:"variant_name,omitempty" gorm:"column:variant_name;size:100"`
	VariantOptions  map[string]interface{} `json:"variant_options,omitempty" gorm:"column:variant_options;type:jsonb"`

	// Item settings
	IsSavedForLater bool            `json:"is_saved_for_later" gorm:"column:is_saved_for_later;default:false"`
	IsGift          bool            `json:"is_gift" gorm:"column:is_gift;default:false"`
	GiftMessage     string          `json:"gift_message,omitempty" gorm:"column:gift_message;size:500"`

	// Availability
	IsAvailable     bool            `json:"is_available" gorm:"column:is_available;default:true"`
	StockQuantity   int             `json:"stock_quantity" gorm:"column:stock_quantity;default:0"`
	MaxQuantity     int             `json:"max_quantity" gorm:"column:max_quantity;default:0"`

	// Timestamps
	AddedAt         time.Time       `json:"added_at" gorm:"column:added_at;not null"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"column:updated_at;not null"`

	// Metadata
	Metadata        map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
}

// Cart represents a shopping cart
type Cart struct {
	ID              uuid.UUID       `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          *uuid.UUID      `json:"user_id,omitempty" gorm:"column:user_id;type:uuid;index"`  // null for guest carts
	SessionID       string          `json:"session_id" gorm:"column:session_id;size:255;index"`       // for guest carts
	Status          CartStatus      `json:"status" gorm:"column:status;type:varchar(20);not null;default:'active'"`

	// Cart items grouped by business
	Items           []CartItem      `json:"items" gorm:"column:items;type:jsonb"`
	ItemCount       int             `json:"item_count" gorm:"column:item_count;default:0"`
	BusinessCount   int             `json:"business_count" gorm:"column:business_count;default:0"`

	// Totals
	SubtotalAmount  decimal.Decimal `json:"subtotal_amount" gorm:"column:subtotal_amount;type:decimal(20,8);default:0"`
	TaxAmount       decimal.Decimal `json:"tax_amount" gorm:"column:tax_amount;type:decimal(20,8);default:0"`
	ShippingAmount  decimal.Decimal `json:"shipping_amount" gorm:"column:shipping_amount;type:decimal(20,8);default:0"`
	DiscountAmount  decimal.Decimal `json:"discount_amount" gorm:"column:discount_amount;type:decimal(20,8);default:0"`
	TotalAmount     decimal.Decimal `json:"total_amount" gorm:"column:total_amount;type:decimal(20,8);default:0"`
	Currency        string          `json:"currency" gorm:"column:currency;size:3;not null;default:'USD'"`

	// Cart behavior
	ExpiresAt       *time.Time      `json:"expires_at,omitempty" gorm:"column:expires_at"`
	LastActivity    time.Time       `json:"last_activity" gorm:"column:last_activity;not null"`
	ConvertedToOrderID *uuid.UUID   `json:"converted_to_order_id,omitempty" gorm:"column:converted_to_order_id;type:uuid"`

	// Applied promotions
	CouponCode      string          `json:"coupon_code,omitempty" gorm:"column:coupon_code;size:50"`
	AppliedPromotions []string      `json:"applied_promotions,omitempty" gorm:"column:applied_promotions;type:jsonb"`

	// Shipping information
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty" gorm:"column:shipping_address;type:jsonb"`

	// Guest cart info
	GuestEmail      string          `json:"guest_email,omitempty" gorm:"column:guest_email;size:255"`
	GuestPhone      string          `json:"guest_phone,omitempty" gorm:"column:guest_phone;size:20"`

	// Regional compliance
	DataRegion      string          `json:"data_region" gorm:"column:data_region;size:20"`
	ShippingCountry string          `json:"shipping_country,omitempty" gorm:"column:shipping_country;size:2"`

	// Metadata
	Metadata        map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	Tags            []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`

	// Timestamps
	CreatedAt       time.Time       `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

// TableName returns the table name for the Cart model
func (Cart) TableName() string {
	return "carts"
}

// CartAbandonmentTracking represents cart abandonment tracking
type CartAbandonmentTracking struct {
	ID                  uuid.UUID   `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	CartID              uuid.UUID   `json:"cart_id" gorm:"column:cart_id;type:uuid;not null;index"`
	UserID              *uuid.UUID  `json:"user_id,omitempty" gorm:"column:user_id;type:uuid"`
	SessionID           string      `json:"session_id" gorm:"column:session_id;size:255"`

	// Abandonment details
	AbandonedAt         time.Time   `json:"abandoned_at" gorm:"column:abandoned_at;not null"`
	AbandonmentStage    string      `json:"abandonment_stage" gorm:"column:abandonment_stage;size:50"` // cart, checkout, payment
	LastPageVisited     string      `json:"last_page_visited,omitempty" gorm:"column:last_page_visited;size:500"`

	// Recovery attempts
	EmailsSent          int         `json:"emails_sent" gorm:"column:emails_sent;default:0"`
	LastEmailSent       *time.Time  `json:"last_email_sent,omitempty" gorm:"column:last_email_sent"`
	RecoveryClicks      int         `json:"recovery_clicks" gorm:"column:recovery_clicks;default:0"`
	LastRecoveryClick   *time.Time  `json:"last_recovery_click,omitempty" gorm:"column:last_recovery_click"`

	// Recovery status
	IsRecovered         bool        `json:"is_recovered" gorm:"column:is_recovered;default:false"`
	RecoveredAt         *time.Time  `json:"recovered_at,omitempty" gorm:"column:recovered_at"`
	RecoveredOrderID    *uuid.UUID  `json:"recovered_order_id,omitempty" gorm:"column:recovered_order_id;type:uuid"`

	// Cart snapshot at abandonment
	CartSnapshot        map[string]interface{} `json:"cart_snapshot" gorm:"column:cart_snapshot;type:jsonb"`

	// Timestamps
	CreatedAt           time.Time   `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt           time.Time   `json:"updated_at" gorm:"column:updated_at;not null"`
}

// TableName returns the table name for the CartAbandonmentTracking model
func (CartAbandonmentTracking) TableName() string {
	return "cart_abandonment_tracking"
}