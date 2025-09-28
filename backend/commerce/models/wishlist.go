package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WishlistType represents the type of wishlist
type WishlistType string

const (
	WishlistTypeDefault   WishlistType = "default"
	WishlistTypeCustom    WishlistType = "custom"
	WishlistTypeShared    WishlistType = "shared"
	WishlistTypeFavorites WishlistType = "favorites"
)

// WishlistPrivacy represents the privacy setting of a wishlist
type WishlistPrivacy string

const (
	WishlistPrivacyPrivate WishlistPrivacy = "private"
	WishlistPrivacyShared  WishlistPrivacy = "shared"
	WishlistPrivacyPublic  WishlistPrivacy = "public"
)

// WishlistItem represents an item in a wishlist
type WishlistItem struct {
	ID          uuid.UUID  `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	ProductID   uuid.UUID  `json:"product_id" gorm:"column:product_id;type:uuid;not null"`
	VariantID   *uuid.UUID `json:"variant_id,omitempty" gorm:"column:variant_id;type:uuid"`
	Quantity    int        `json:"quantity" gorm:"column:quantity;default:1"`
	Note        string     `json:"note,omitempty" gorm:"column:note;size:500"`
	Priority    int        `json:"priority" gorm:"column:priority;default:0"`
	AddedAt     time.Time  `json:"added_at" gorm:"column:added_at;not null"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;not null"`

	// Cached product info for performance
	ProductName  string `json:"product_name" gorm:"column:product_name;size:200"`
	ProductPrice string `json:"product_price" gorm:"column:product_price;size:50"`
	ProductImage string `json:"product_image,omitempty" gorm:"column:product_image;size:500"`
}

// Wishlist represents a user's wishlist
type Wishlist struct {
	ID          uuid.UUID       `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID       `json:"user_id" gorm:"column:user_id;type:uuid;not null;index"`
	Type        WishlistType    `json:"type" gorm:"column:type;type:varchar(20);not null;default:'default'"`
	Privacy     WishlistPrivacy `json:"privacy" gorm:"column:privacy;type:varchar(20);not null;default:'private'"`

	// Wishlist details
	Name        string `json:"name" gorm:"column:name;size:100;not null"`
	Description string `json:"description,omitempty" gorm:"column:description;size:500"`
	IsDefault   bool   `json:"is_default" gorm:"column:is_default;default:false"`

	// Items
	Items     []WishlistItem `json:"items" gorm:"column:items;type:jsonb"`
	ItemCount int            `json:"item_count" gorm:"column:item_count;default:0"`

	// Sharing
	ShareToken  string     `json:"share_token,omitempty" gorm:"column:share_token;size:64;uniqueIndex"`
	SharedWith  []uuid.UUID `json:"shared_with,omitempty" gorm:"column:shared_with;type:jsonb"`
	ShareCount  int        `json:"share_count" gorm:"column:share_count;default:0"`

	// Settings
	AllowDuplicates    bool `json:"allow_duplicates" gorm:"column:allow_duplicates;default:false"`
	AutoRemoveOnPurchase bool `json:"auto_remove_on_purchase" gorm:"column:auto_remove_on_purchase;default:true"`

	// Regional compliance
	DataRegion string `json:"data_region" gorm:"column:data_region;size:20"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	Tags     []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

// TableName returns the table name for the Wishlist model
func (Wishlist) TableName() string {
	return "wishlists"
}

// WishlistShare represents a wishlist share record
type WishlistShare struct {
	ID         uuid.UUID `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	WishlistID uuid.UUID `json:"wishlist_id" gorm:"column:wishlist_id;type:uuid;not null;index"`
	SharedBy   uuid.UUID `json:"shared_by" gorm:"column:shared_by;type:uuid;not null"`
	SharedWith uuid.UUID `json:"shared_with" gorm:"column:shared_with;type:uuid;not null"`
	Permission string    `json:"permission" gorm:"column:permission;size:20;not null;default:'view'"` // view, edit
	ViewCount  int       `json:"view_count" gorm:"column:view_count;default:0"`
	LastViewed *time.Time `json:"last_viewed,omitempty" gorm:"column:last_viewed"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;not null"`

	// Unique constraint on wishlist_id + shared_with
}

// TableName returns the table name for the WishlistShare model
func (WishlistShare) TableName() string {
	return "wishlist_shares"
}

// ProductFollow represents a user following a product for updates
type ProductFollow struct {
	ID              uuid.UUID `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null;index"`
	ProductID       uuid.UUID `json:"product_id" gorm:"column:product_id;type:uuid;not null;index"`
	BusinessID      uuid.UUID `json:"business_id" gorm:"column:business_id;type:uuid;not null;index"`

	// Notification preferences
	NotifyPriceChange   bool `json:"notify_price_change" gorm:"column:notify_price_change;default:true"`
	NotifyBackInStock   bool `json:"notify_back_in_stock" gorm:"column:notify_back_in_stock;default:true"`
	NotifyPromotion     bool `json:"notify_promotion" gorm:"column:notify_promotion;default:true"`
	NotifyNewVariant    bool `json:"notify_new_variant" gorm:"column:notify_new_variant;default:false"`

	// Cached product info
	ProductName  string `json:"product_name" gorm:"column:product_name;size:200"`
	ProductPrice string `json:"product_price" gorm:"column:product_price;size:50"`
	ProductImage string `json:"product_image,omitempty" gorm:"column:product_image;size:500"`

	// Follow status
	IsActive     bool       `json:"is_active" gorm:"column:is_active;default:true"`
	LastNotified *time.Time `json:"last_notified,omitempty" gorm:"column:last_notified"`

	// Regional compliance
	DataRegion string `json:"data_region" gorm:"column:data_region;size:20"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Unique constraint on user_id + product_id
}

// TableName returns the table name for the ProductFollow model
func (ProductFollow) TableName() string {
	return "product_follows"
}