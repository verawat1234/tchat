package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryStatus represents the status of a category
type CategoryStatus string

const (
	CategoryStatusActive   CategoryStatus = "active"
	CategoryStatusInactive CategoryStatus = "inactive"
	CategoryStatusDraft    CategoryStatus = "draft"
)

// CategoryType represents the type of category
type CategoryType string

const (
	CategoryTypeProduct CategoryType = "product"
	CategoryTypeService CategoryType = "service"
	CategoryTypeBoth    CategoryType = "both"
)

// CategoryImage represents an image for a category
type CategoryImage struct {
	URL         string `json:"url" gorm:"column:url;size:500;not null"`
	AltText     string `json:"alt_text,omitempty" gorm:"column:alt_text;size:200"`
	Width       int    `json:"width,omitempty" gorm:"column:width"`
	Height      int    `json:"height,omitempty" gorm:"column:height"`
}

// CategorySEO represents SEO information for a category
type CategorySEO struct {
	MetaTitle       string   `json:"meta_title,omitempty" gorm:"column:meta_title;size:200"`
	MetaDescription string   `json:"meta_description,omitempty" gorm:"column:meta_description;size:500"`
	Keywords        []string `json:"keywords,omitempty" gorm:"column:keywords;type:jsonb"`
	URLSlug         string   `json:"url_slug,omitempty" gorm:"column:url_slug;size:200"`
	Canonical       string   `json:"canonical,omitempty" gorm:"column:canonical;size:500"`
}

// Category represents a product/service category
type Category struct {
	ID              uuid.UUID       `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID      *uuid.UUID      `json:"business_id,omitempty" gorm:"column:business_id;type:uuid"` // null for global categories

	// Category hierarchy
	ParentID        *uuid.UUID      `json:"parent_id,omitempty" gorm:"column:parent_id;type:uuid"`
	Level           int             `json:"level" gorm:"column:level;default:0"`
	Path            string          `json:"path" gorm:"column:path;size:500"` // e.g. "electronics/phones/smartphones"

	// Basic info
	Name            string          `json:"name" gorm:"column:name;size:100;not null"`
	Description     string          `json:"description,omitempty" gorm:"column:description;size:1000"`
	ShortDescription string         `json:"short_description,omitempty" gorm:"column:short_description;size:255"`
	Status          CategoryStatus  `json:"status" gorm:"column:status;type:varchar(20);not null;default:'active'"`
	Type            CategoryType    `json:"type" gorm:"column:type;type:varchar(20);not null;default:'product'"`

	// Display
	Icon            string          `json:"icon,omitempty" gorm:"column:icon;size:100"`
	Image           *CategoryImage  `json:"image,omitempty" gorm:"column:image;type:jsonb"`
	Color           string          `json:"color,omitempty" gorm:"column:color;size:7"` // hex color
	SortOrder       int             `json:"sort_order" gorm:"column:sort_order;default:0"`

	// Behavior
	IsVisible       bool            `json:"is_visible" gorm:"column:is_visible;default:true"`
	IsFeatured      bool            `json:"is_featured" gorm:"column:is_featured;default:false"`
	AllowProducts   bool            `json:"allow_products" gorm:"column:allow_products;default:true"`

	// SEO
	SEO             CategorySEO     `json:"seo" gorm:"embedded;embeddedPrefix:seo_"`

	// Statistics
	ProductCount    int             `json:"product_count" gorm:"column:product_count;default:0"`
	ActiveProductCount int          `json:"active_product_count" gorm:"column:active_product_count;default:0"`
	ChildrenCount   int             `json:"children_count" gorm:"column:children_count;default:0"`

	// Attributes that can be used for products in this category
	Attributes      []CategoryAttribute `json:"attributes,omitempty" gorm:"column:attributes;type:jsonb"`

	// Regional compliance
	DataRegion      string          `json:"data_region" gorm:"column:data_region;size:20"`
	RestrictedCountries []string    `json:"restricted_countries,omitempty" gorm:"column:restricted_countries;type:jsonb"`

	// Metadata
	Metadata        map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	Tags            []string               `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`

	// Timestamps
	CreatedAt       time.Time       `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	Parent          *Category       `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Children        []Category      `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
}

// TableName returns the table name for the Category model
func (Category) TableName() string {
	return "categories"
}

// CategoryAttribute represents an attribute that can be used for products in a category
type CategoryAttribute struct {
	Name        string                 `json:"name" gorm:"column:name;size:100;not null"`
	Type        string                 `json:"type" gorm:"column:type;size:20;not null"` // text, number, boolean, select, multiselect
	Required    bool                   `json:"required" gorm:"column:required;default:false"`
	Options     []string               `json:"options,omitempty" gorm:"column:options;type:jsonb"`
	DefaultValue interface{}           `json:"default_value,omitempty" gorm:"column:default_value;type:jsonb"`
	Validation  map[string]interface{} `json:"validation,omitempty" gorm:"column:validation;type:jsonb"`
	SortOrder   int                    `json:"sort_order" gorm:"column:sort_order;default:0"`
}

// ProductCategory represents the many-to-many relationship between products and categories
type ProductCategory struct {
	ID         uuid.UUID `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID  uuid.UUID `json:"product_id" gorm:"column:product_id;type:uuid;not null;index"`
	CategoryID uuid.UUID `json:"category_id" gorm:"column:category_id;type:uuid;not null;index"`
	IsPrimary  bool      `json:"is_primary" gorm:"column:is_primary;default:false"`
	SortOrder  int       `json:"sort_order" gorm:"column:sort_order;default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;not null"`

	// Unique constraint on product_id + category_id
}

// TableName returns the table name for the ProductCategory model
func (ProductCategory) TableName() string {
	return "product_categories"
}

// CategoryView represents a category view for analytics
type CategoryView struct {
	ID         uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	CategoryID uuid.UUID  `json:"category_id" gorm:"column:category_id;type:uuid;not null;index"`
	UserID     *uuid.UUID `json:"user_id,omitempty" gorm:"column:user_id;type:uuid"`
	SessionID  string     `json:"session_id" gorm:"column:session_id;size:255;not null"`
	IPAddress  string     `json:"ip_address" gorm:"column:ip_address;size:45"`
	UserAgent  string     `json:"user_agent,omitempty" gorm:"column:user_agent;size:500"`
	Referrer   string     `json:"referrer,omitempty" gorm:"column:referrer;size:500"`
	ViewedAt   time.Time  `json:"viewed_at" gorm:"column:viewed_at;not null"`
}

// TableName returns the table name for the CategoryView model
func (CategoryView) TableName() string {
	return "category_views"
}