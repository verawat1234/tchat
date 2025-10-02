package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusInactive  ProductStatus = "inactive"
	ProductStatusArchived  ProductStatus = "archived"
	ProductStatusDeleted   ProductStatus = "deleted"
)

// IsValid checks if the product status is valid
func (ps ProductStatus) IsValid() bool {
	switch ps {
	case ProductStatusDraft, ProductStatusActive, ProductStatusInactive, ProductStatusArchived, ProductStatusDeleted:
		return true
	default:
		return false
	}
}

// IsAvailable checks if the product is available for purchase
func (ps ProductStatus) IsAvailable() bool {
	return ps == ProductStatusActive
}

// ProductType represents the type of product
type ProductType string

const (
	ProductTypePhysical ProductType = "physical"
	ProductTypeDigital  ProductType = "digital"
	ProductTypeService  ProductType = "service"
)

// IsValid checks if the product type is valid
func (pt ProductType) IsValid() bool {
	switch pt {
	case ProductTypePhysical, ProductTypeDigital, ProductTypeService:
		return true
	default:
		return false
	}
}

// ProductCondition represents the condition of a product
type ProductCondition string

const (
	ProductConditionNew         ProductCondition = "new"
	ProductConditionUsed        ProductCondition = "used"
	ProductConditionRefurbished ProductCondition = "refurbished"
)

// IsValid checks if the product condition is valid
func (pc ProductCondition) IsValid() bool {
	switch pc {
	case ProductConditionNew, ProductConditionUsed, ProductConditionRefurbished:
		return true
	default:
		return false
	}
}

// ProductImage represents a product image
type ProductImage struct {
	URL         string `json:"url" gorm:"column:url;size:500;not null"`
	AltText     string `json:"alt_text,omitempty" gorm:"column:alt_text;size:200"`
	IsPrimary   bool   `json:"is_primary" gorm:"column:is_primary;default:false"`
	SortOrder   int    `json:"sort_order" gorm:"column:sort_order;default:0"`
	Width       int    `json:"width,omitempty" gorm:"column:width"`
	Height      int    `json:"height,omitempty" gorm:"column:height"`
	Size        int64  `json:"size,omitempty" gorm:"column:size"` // File size in bytes
}

// ProductVariant represents a product variant (size, color, etc.)
type ProductVariant struct {
	ID          uuid.UUID `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"column:name;size:100;not null"`
	SKU         string    `json:"sku,omitempty" gorm:"column:sku;size:100"`
	Price       decimal.Decimal `json:"price" gorm:"column:price;type:decimal(20,8);not null"`
	ComparePrice *decimal.Decimal `json:"compare_price,omitempty" gorm:"column:compare_price;type:decimal(20,8)"`
	Stock       int       `json:"stock" gorm:"column:stock;default:0"`
	Weight      *decimal.Decimal `json:"weight,omitempty" gorm:"column:weight;type:decimal(10,3)"`
	Dimensions  map[string]interface{} `json:"dimensions,omitempty" gorm:"column:dimensions;type:jsonb"`
	Attributes  map[string]interface{} `json:"attributes,omitempty" gorm:"column:attributes;type:jsonb"`
	IsDefault   bool      `json:"is_default" gorm:"column:is_default;default:false"`
	IsActive    bool      `json:"is_active" gorm:"column:is_active;default:true"`
}

// ProductPricing represents pricing information for different regions
type ProductPricing struct {
	Currency    string          `json:"currency" gorm:"column:currency;size:3;not null"`
	Price       decimal.Decimal `json:"price" gorm:"column:price;type:decimal(20,8);not null"`
	ComparePrice *decimal.Decimal `json:"compare_price,omitempty" gorm:"column:compare_price;type:decimal(20,8)"`
	CostPrice   *decimal.Decimal `json:"cost_price,omitempty" gorm:"column:cost_price;type:decimal(20,8)"`
	TaxRate     decimal.Decimal `json:"tax_rate" gorm:"column:tax_rate;type:decimal(5,4);default:0"`
	TaxInclusive bool           `json:"tax_inclusive" gorm:"column:tax_inclusive;default:false"`
}

// ProductShipping represents shipping information
type ProductShipping struct {
	Weight      *decimal.Decimal `json:"weight,omitempty" gorm:"column:weight;type:decimal(10,3)"`
	Length      *decimal.Decimal `json:"length,omitempty" gorm:"column:length;type:decimal(10,2)"`
	Width       *decimal.Decimal `json:"width,omitempty" gorm:"column:width;type:decimal(10,2)"`
	Height      *decimal.Decimal `json:"height,omitempty" gorm:"column:height;type:decimal(10,2)"`
	RequiresShipping bool         `json:"requires_shipping" gorm:"column:requires_shipping;default:true"`
	ShippingClass    string       `json:"shipping_class,omitempty" gorm:"column:shipping_class;size:50"`
	FreeShipping     bool         `json:"free_shipping" gorm:"column:free_shipping;default:false"`
}

// ProductSEO represents SEO information
type ProductSEO struct {
	MetaTitle       string `json:"meta_title,omitempty" gorm:"column:meta_title;size:200"`
	MetaDescription string `json:"meta_description,omitempty" gorm:"column:meta_description;size:500"`
	Keywords        []string `json:"keywords,omitempty" gorm:"column:keywords;type:jsonb"`
	URLSlug         string `json:"url_slug,omitempty" gorm:"column:url_slug;size:200"`
}

// ProductInventory represents inventory tracking information
type ProductInventory struct {
	TrackQuantity    bool `json:"track_quantity" gorm:"column:track_quantity;default:true"`
	Stock           int  `json:"stock" gorm:"column:stock;default:0"`
	LowStockThreshold int `json:"low_stock_threshold" gorm:"column:low_stock_threshold;default:10"`
	AllowBackorder  bool `json:"allow_backorder" gorm:"column:allow_backorder;default:false"`
	ManageStock     bool `json:"manage_stock" gorm:"column:manage_stock;default:true"`
}

// Product represents a product in the system
type Product struct {
	ID         uuid.UUID `json:"id" gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID uuid.UUID `json:"business_id" gorm:"column:business_id;type:uuid;not null;index"`

	// Basic product information
	Name        string           `json:"name" gorm:"column:name;size:200;not null"`
	Description string           `json:"description" gorm:"column:description;size:5000"`
	ShortDescription string      `json:"short_description,omitempty" gorm:"column:short_description;size:500"`
	Type        ProductType      `json:"type" gorm:"column:type;type:varchar(20);not null;default:'physical'"`
	Status      ProductStatus    `json:"status" gorm:"column:status;type:varchar(20);not null;default:'draft'"`
	Condition   ProductCondition `json:"condition" gorm:"column:condition;type:varchar(20);not null;default:'new'"`

	// Identification and organization
	SKU      string   `json:"sku,omitempty" gorm:"column:sku;size:100;uniqueIndex"`
	Category string   `json:"category" gorm:"column:category;size:100;not null"`
	Tags     []string `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`
	Brand    string   `json:"brand,omitempty" gorm:"column:brand;size:100"`

	// Pricing (multi-currency support)
	Pricing []ProductPricing `json:"pricing" gorm:"column:pricing;type:jsonb"`

	// Media
	Images []ProductImage `json:"images,omitempty" gorm:"column:images;type:jsonb"`
	Videos []string       `json:"videos,omitempty" gorm:"column:videos;type:jsonb"`

	// Variants
	HasVariants bool             `json:"has_variants" gorm:"column:has_variants;default:false"`
	Variants    []ProductVariant `json:"variants,omitempty" gorm:"column:variants;type:jsonb"`

	// Inventory management
	Inventory ProductInventory `json:"inventory" gorm:"embedded;embeddedPrefix:inventory_"`

	// Shipping information
	Shipping ProductShipping `json:"shipping" gorm:"embedded;embeddedPrefix:shipping_"`

	// SEO and search
	SEO            ProductSEO `json:"seo" gorm:"embedded;embeddedPrefix:seo_"`
	SearchKeywords []string   `json:"search_keywords,omitempty" gorm:"column:search_keywords;type:jsonb"`

	// Product specifications and attributes
	Specifications map[string]interface{} `json:"specifications,omitempty" gorm:"column:specifications;type:jsonb"`
	Attributes     map[string]interface{} `json:"attributes,omitempty" gorm:"column:attributes;type:jsonb"`

	// Regional compliance and restrictions
	DataRegion         string   `json:"data_region" gorm:"column:data_region;size:20"`
	AllowedCountries   []string `json:"allowed_countries,omitempty" gorm:"column:allowed_countries;type:jsonb"`
	RestrictedCountries []string `json:"restricted_countries,omitempty" gorm:"column:restricted_countries;type:jsonb"`
	ComplianceData     map[string]interface{} `json:"compliance_data,omitempty" gorm:"column:compliance_data;type:jsonb"`

	// Sales and performance metrics
	SalesCount    int64           `json:"sales_count" gorm:"column:sales_count;default:0"`
	ViewCount     int64           `json:"view_count" gorm:"column:view_count;default:0"`
	AverageRating decimal.Decimal `json:"average_rating" gorm:"column:average_rating;type:decimal(3,2);default:0"`
	ReviewCount   int             `json:"review_count" gorm:"column:review_count;default:0"`

	// Promotion and visibility
	IsFeatured  bool `json:"is_featured" gorm:"column:is_featured;default:false"`
	IsPromoted  bool `json:"is_promoted" gorm:"column:is_promoted;default:false"`
	SortOrder   int  `json:"sort_order" gorm:"column:sort_order;default:0"`

	// Timestamps and scheduling
	PublishedAt *time.Time `json:"published_at,omitempty" gorm:"column:published_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" gorm:"column:expires_at"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`

	// Relationships
	Business *Business `json:"business,omitempty" gorm:"foreignKey:BusinessID;references:ID"`
}

// TableName returns the table name for the Product model
func (Product) TableName() string {
	return "products"
}

// BeforeCreate sets up the product before creation
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	// Generate SKU if not provided
	if p.SKU == "" {
		p.SKU = p.generateSKU()
	}

	// Set data region based on business
	if p.DataRegion == "" {
		var business Business
		if err := tx.First(&business, p.BusinessID).Error; err == nil {
			p.DataRegion = GetDataRegionForCountry(business.Address.Country)
		} else {
			p.DataRegion = "sea-central" // Default region
		}
	}

	// Initialize default pricing if empty
	if len(p.Pricing) == 0 && len(p.Variants) == 0 {
		return fmt.Errorf("product must have either pricing or variants")
	}

	// Set URL slug if not provided
	if p.SEO.URLSlug == "" {
		p.SEO.URLSlug = p.generateURLSlug()
	}

	// Generate search keywords
	p.SearchKeywords = p.GenerateSearchKeywords()

	// Set published date for active products
	if p.Status == ProductStatusActive && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}

	// Initialize default inventory settings
	if !p.HasVariants {
		if p.Inventory.LowStockThreshold == 0 {
			p.Inventory.LowStockThreshold = 10
		}
	}

	// Validate the product
	if err := p.Validate(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate validates the product before updating
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	// Update search keywords
	p.SearchKeywords = p.GenerateSearchKeywords()

	// Update published date when status changes to active
	if p.Status == ProductStatusActive && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}

	return p.Validate()
}

// Validate validates the product data
func (p *Product) Validate() error {
	// Validate UUIDs
	if p.ID == uuid.Nil {
		return fmt.Errorf("product ID cannot be nil")
	}
	if p.BusinessID == uuid.Nil {
		return fmt.Errorf("business ID cannot be nil")
	}

	// Validate basic fields
	if len(p.Name) == 0 || len(p.Name) > 200 {
		return fmt.Errorf("product name must be between 1 and 200 characters")
	}

	if !p.Type.IsValid() {
		return fmt.Errorf("invalid product type: %s", p.Type)
	}

	if !p.Status.IsValid() {
		return fmt.Errorf("invalid product status: %s", p.Status)
	}

	if !p.Condition.IsValid() {
		return fmt.Errorf("invalid product condition: %s", p.Condition)
	}

	if p.Category == "" {
		return fmt.Errorf("product category is required")
	}

	// Validate pricing
	if err := p.validatePricing(); err != nil {
		return err
	}

	// Validate variants
	if err := p.validateVariants(); err != nil {
		return err
	}

	// Validate inventory
	if err := p.validateInventory(); err != nil {
		return err
	}

	// Validate shipping for physical products
	if p.Type == ProductTypePhysical {
		if err := p.validateShipping(); err != nil {
			return err
		}
	}

	return nil
}

// validatePricing validates the pricing information
func (p *Product) validatePricing() error {
	if !p.HasVariants && len(p.Pricing) == 0 {
		return fmt.Errorf("product without variants must have pricing information")
	}

	currencySeen := make(map[string]bool)
	for i, pricing := range p.Pricing {
		if !IsValidCurrency(pricing.Currency) {
			return fmt.Errorf("invalid currency at pricing index %d: %s", i, pricing.Currency)
		}

		if currencySeen[pricing.Currency] {
			return fmt.Errorf("duplicate currency in pricing: %s", pricing.Currency)
		}
		currencySeen[pricing.Currency] = true

		if pricing.Price.IsNegative() || pricing.Price.IsZero() {
			return fmt.Errorf("price must be positive for currency %s", pricing.Currency)
		}

		if pricing.ComparePrice != nil && pricing.ComparePrice.LessThan(pricing.Price) {
			return fmt.Errorf("compare price must be greater than or equal to price for currency %s", pricing.Currency)
		}
	}

	return nil
}

// validateVariants validates the product variants
func (p *Product) validateVariants() error {
	if p.HasVariants && len(p.Variants) == 0 {
		return fmt.Errorf("product with variants must have at least one variant")
	}

	if !p.HasVariants && len(p.Variants) > 0 {
		return fmt.Errorf("product without variants should not have variant data")
	}

	skuSeen := make(map[string]bool)
	defaultCount := 0

	for i, variant := range p.Variants {
		if variant.Name == "" {
			return fmt.Errorf("variant name is required at index %d", i)
		}

		if variant.SKU != "" {
			if skuSeen[variant.SKU] {
				return fmt.Errorf("duplicate variant SKU: %s", variant.SKU)
			}
			skuSeen[variant.SKU] = true
		}

		if variant.Price.IsNegative() || variant.Price.IsZero() {
			return fmt.Errorf("variant price must be positive at index %d", i)
		}

		if variant.IsDefault {
			defaultCount++
		}
	}

	if p.HasVariants && defaultCount != 1 {
		return fmt.Errorf("product with variants must have exactly one default variant")
	}

	return nil
}

// validateInventory validates inventory settings
func (p *Product) validateInventory() error {
	if p.Inventory.TrackQuantity && p.Inventory.Stock < 0 {
		return fmt.Errorf("stock quantity cannot be negative")
	}

	if p.Inventory.LowStockThreshold < 0 {
		return fmt.Errorf("low stock threshold cannot be negative")
	}

	return nil
}

// validateShipping validates shipping information
func (p *Product) validateShipping() error {
	if p.Shipping.Weight != nil && p.Shipping.Weight.IsNegative() {
		return fmt.Errorf("shipping weight cannot be negative")
	}

	dimensions := []*decimal.Decimal{p.Shipping.Length, p.Shipping.Width, p.Shipping.Height}
	dimensionNames := []string{"length", "width", "height"}

	for i, dimension := range dimensions {
		if dimension != nil && dimension.IsNegative() {
			return fmt.Errorf("shipping %s cannot be negative", dimensionNames[i])
		}
	}

	return nil
}

// generateSKU generates a SKU for the product
func (p *Product) generateSKU() string {
	// Simple SKU generation: BUSINESS_ID[:8]-PRODUCT_ID[:8]
	return fmt.Sprintf("%s-%s", p.BusinessID.String()[:8], p.ID.String()[:8])
}

// generateURLSlug generates a URL slug from the product name
func (p *Product) generateURLSlug() string {
	slug := strings.ToLower(p.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (simplified)
	slug = strings.ReplaceAll(slug, ".", "")
	slug = strings.ReplaceAll(slug, ",", "")
	slug = strings.ReplaceAll(slug, "!", "")
	slug = strings.ReplaceAll(slug, "?", "")
	return slug
}

// GetPriceForCurrency returns the price for a specific currency
func (p *Product) GetPriceForCurrency(currency string) (*ProductPricing, error) {
	for _, pricing := range p.Pricing {
		if pricing.Currency == currency {
			return &pricing, nil
		}
	}
	return nil, fmt.Errorf("price not found for currency: %s", currency)
}

// GetDefaultVariant returns the default variant
func (p *Product) GetDefaultVariant() *ProductVariant {
	for i, variant := range p.Variants {
		if variant.IsDefault {
			return &p.Variants[i]
		}
	}
	return nil
}

// GetVariantByID returns a variant by its ID
func (p *Product) GetVariantByID(variantID uuid.UUID) *ProductVariant {
	for i, variant := range p.Variants {
		if variant.ID == variantID {
			return &p.Variants[i]
		}
	}
	return nil
}

// IsAvailable checks if the product is available for purchase
func (p *Product) IsAvailable() bool {
	if !p.Status.IsAvailable() {
		return false
	}

	if p.ExpiresAt != nil && p.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}

// IsInStock checks if the product is in stock
func (p *Product) IsInStock() bool {
	if !p.Inventory.TrackQuantity {
		return true
	}

	if p.HasVariants {
		for _, variant := range p.Variants {
			if variant.IsActive && variant.Stock > 0 {
				return true
			}
		}
		return false
	}

	return p.Inventory.Stock > 0 || p.Inventory.AllowBackorder
}

// IsLowStock checks if the product is low in stock
func (p *Product) IsLowStock() bool {
	if !p.Inventory.TrackQuantity {
		return false
	}

	if p.HasVariants {
		for _, variant := range p.Variants {
			if variant.IsActive && variant.Stock <= p.Inventory.LowStockThreshold {
				return true
			}
		}
		return false
	}

	return p.Inventory.Stock <= p.Inventory.LowStockThreshold
}

// GetPrimaryImage returns the primary product image
func (p *Product) GetPrimaryImage() *ProductImage {
	for i, image := range p.Images {
		if image.IsPrimary {
			return &p.Images[i]
		}
	}
	// Return first image if no primary is set
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

// IncrementViewCount increments the view count
func (p *Product) IncrementViewCount() {
	p.ViewCount++
	p.UpdatedAt = time.Now()
}

// IncrementSalesCount increments the sales count
func (p *Product) IncrementSalesCount() {
	p.SalesCount++
	p.UpdatedAt = time.Now()
}

// UpdateRating updates the average rating and review count
func (p *Product) UpdateRating(newRating decimal.Decimal) {
	totalRating := p.AverageRating.Mul(decimal.NewFromInt(int64(p.ReviewCount)))
	totalRating = totalRating.Add(newRating)
	p.ReviewCount++
	p.AverageRating = totalRating.Div(decimal.NewFromInt(int64(p.ReviewCount)))
	p.UpdatedAt = time.Now()
}

// GenerateSearchKeywords generates search keywords for the product
func (p *Product) GenerateSearchKeywords() []string {
	keywords := []string{
		p.Name,
		p.Category,
		p.Brand,
		string(p.Type),
		string(p.Condition),
	}

	// Add description keywords
	if p.Description != "" {
		descWords := strings.Fields(strings.ToLower(p.Description))
		keywords = append(keywords, descWords...)
	}

	// Add short description keywords
	if p.ShortDescription != "" {
		shortDescWords := strings.Fields(strings.ToLower(p.ShortDescription))
		keywords = append(keywords, shortDescWords...)
	}

	// Add tags
	keywords = append(keywords, p.Tags...)

	// Add SEO keywords
	keywords = append(keywords, p.SEO.Keywords...)

	// Add variant names
	for _, variant := range p.Variants {
		keywords = append(keywords, strings.ToLower(variant.Name))
	}

	// Remove duplicates and empty strings
	seen := make(map[string]bool)
	var unique []string
	for _, keyword := range keywords {
		cleaned := strings.ToLower(strings.TrimSpace(keyword))
		if cleaned != "" && len(cleaned) > 2 && !seen[cleaned] {
			seen[cleaned] = true
			unique = append(unique, cleaned)
		}
	}

	return unique
}

// GetProductSummary returns a summary of product information
func (p *Product) GetProductSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"id":              p.ID,
		"name":            p.Name,
		"sku":             p.SKU,
		"type":            p.Type,
		"status":          p.Status,
		"condition":       p.Condition,
		"category":        p.Category,
		"brand":           p.Brand,
		"is_available":    p.IsAvailable(),
		"is_in_stock":     p.IsInStock(),
		"is_low_stock":    p.IsLowStock(),
		"has_variants":    p.HasVariants,
		"is_featured":     p.IsFeatured,
		"is_promoted":     p.IsPromoted,
		"sales_count":     p.SalesCount,
		"view_count":      p.ViewCount,
		"average_rating":  p.AverageRating,
		"review_count":    p.ReviewCount,
		"created_at":      p.CreatedAt,
		"published_at":    p.PublishedAt,
	}

	// Add primary image
	if primaryImage := p.GetPrimaryImage(); primaryImage != nil {
		summary["primary_image"] = primaryImage.URL
	}

	return summary
}

// MarshalJSON customizes JSON serialization
func (p *Product) MarshalJSON() ([]byte, error) {
	type Alias Product
	return json.Marshal(&struct {
		*Alias
		IsAvailable    bool                   `json:"is_available"`
		IsInStock      bool                   `json:"is_in_stock"`
		IsLowStock     bool                   `json:"is_low_stock"`
		PrimaryImage   *ProductImage          `json:"primary_image,omitempty"`
		DefaultVariant *ProductVariant        `json:"default_variant,omitempty"`
		ProductSummary map[string]interface{} `json:"product_summary"`
	}{
		Alias:          (*Alias)(p),
		IsAvailable:    p.IsAvailable(),
		IsInStock:      p.IsInStock(),
		IsLowStock:     p.IsLowStock(),
		PrimaryImage:   p.GetPrimaryImage(),
		DefaultVariant: p.GetDefaultVariant(),
		ProductSummary: p.GetProductSummary(),
	})
}

// Helper functions for product categories

// GetValidProductCategories returns valid product categories for Southeast Asian markets
func GetValidProductCategories() []string {
	return []string{
		"electronics",
		"fashion",
		"food",
		"health",
		"beauty",
		"home",
		"sports",
		"automotive",
		"books",
		"toys",
		"services",
		"digital",
		"agriculture",
		"crafts",
		"jewelry",
		"travel",
		"education",
		"finance",
		"real_estate",
		"entertainment",
	}
}

// IsValidProductCategory checks if a category is valid
func IsValidProductCategory(category string) bool {
	validCategories := GetValidProductCategories()
	for _, valid := range validCategories {
		if category == valid {
			return true
		}
	}
	return false
}