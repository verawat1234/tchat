package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Product represents a product available for purchase in the commerce system
type Product struct {
	ID               uuid.UUID       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ShopID           uuid.UUID       `json:"shop_id" gorm:"type:varchar(36);not null;index"`
	SKU              string          `json:"sku" gorm:"type:varchar(100);uniqueIndex;not null"`
	Name             string          `json:"name" gorm:"type:varchar(255);not null"`
	Description      *string         `json:"description,omitempty" gorm:"type:text"`
	ShortDescription *string         `json:"short_description,omitempty" gorm:"type:varchar(500)"`
	Category         ProductCategory `json:"category" gorm:"type:varchar(50);not null"`
	Tags             StringSlice     `json:"tags" gorm:"type:json"`
	Brand            *string         `json:"brand,omitempty" gorm:"type:varchar(100)"`
	Model            *string         `json:"model,omitempty" gorm:"type:varchar(100)"`
	Price            int64           `json:"price" gorm:"not null"`                    // Price in cents
	Currency         Currency        `json:"currency" gorm:"type:varchar(3);not null"`
	ComparePrice     *int64          `json:"compare_price,omitempty"`                 // Original price for discounts
	CostPrice        *int64          `json:"cost_price,omitempty"`                    // Cost to merchant
	Weight           *float64        `json:"weight,omitempty"`                        // Weight in grams
	Dimensions       *Dimensions     `json:"dimensions,omitempty" gorm:"type:json"`
	Images           ImageSlice      `json:"images" gorm:"type:json"`
	Variants         VariantSlice    `json:"variants,omitempty" gorm:"type:json"`
	Inventory        Inventory       `json:"inventory" gorm:"type:json"`
	SEO              SEOInfo         `json:"seo" gorm:"type:json"`
	Shipping         ShippingInfo    `json:"shipping" gorm:"type:json"`
	Attributes       ProductAttrs    `json:"attributes" gorm:"type:json"`
	Metadata         ProductMetadata `json:"metadata" gorm:"type:json"`
	Status           ProductStatus   `json:"status" gorm:"type:varchar(20);default:'draft'"`
	Visibility       ProductVisibility `json:"visibility" gorm:"type:varchar(20);default:'public'"`
	IsDigital        bool            `json:"is_digital" gorm:"default:false"`
	RequiresShipping bool            `json:"requires_shipping" gorm:"default:true"`
	IsTaxable        bool            `json:"is_taxable" gorm:"default:true"`
	FeaturedAt       *time.Time      `json:"featured_at,omitempty"`
	PublishedAt      *time.Time      `json:"published_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time       `json:"updated_at" gorm:"not null"`
}

// ProductCategory represents product categories for Southeast Asian markets
type ProductCategory string

const (
	CategoryElectronics      ProductCategory = "electronics"
	CategoryFashion          ProductCategory = "fashion"
	CategoryHome             ProductCategory = "home"
	CategoryBeauty           ProductCategory = "beauty"
	CategorySports           ProductCategory = "sports"
	CategoryBooks            ProductCategory = "books"
	CategoryFood             ProductCategory = "food"
	CategoryHealth           ProductCategory = "health"
	CategoryToys             ProductCategory = "toys"
	CategoryAutomotive       ProductCategory = "automotive"
	CategoryServices         ProductCategory = "services"
	CategoryDigital          ProductCategory = "digital"
	CategoryTraditionalCrafts ProductCategory = "traditional_crafts"
	CategoryLocalProducts    ProductCategory = "local_products"
	CategoryHalal            ProductCategory = "halal"
	CategoryOrganic          ProductCategory = "organic"
)

// ProductStatus represents the current status of a product
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusInactive  ProductStatus = "inactive"
	ProductStatusArchived  ProductStatus = "archived"
	ProductStatusDeleted   ProductStatus = "deleted"
	ProductStatusPending   ProductStatus = "pending"
	ProductStatusRejected  ProductStatus = "rejected"
)

// ProductVisibility represents product visibility settings
type ProductVisibility string

const (
	VisibilityPublic    ProductVisibility = "public"
	VisibilityPrivate   ProductVisibility = "private"
	VisibilityHidden    ProductVisibility = "hidden"
	VisibilityScheduled ProductVisibility = "scheduled"
)

// Currency represents supported currencies (reusing from payment models)
type Currency string

const (
	CurrencyTHB Currency = "THB" // Thai Baht
	CurrencySGD Currency = "SGD" // Singapore Dollar
	CurrencyIDR Currency = "IDR" // Indonesian Rupiah
	CurrencyMYR Currency = "MYR" // Malaysian Ringgit
	CurrencyPHP Currency = "PHP" // Philippine Peso
	CurrencyVND Currency = "VND" // Vietnamese Dong
	CurrencyUSD Currency = "USD" // US Dollar
)

// StringSlice represents a slice of strings stored as JSON
type StringSlice []string

// ImageSlice represents a slice of product images
type ImageSlice []ProductImage

// VariantSlice represents a slice of product variants
type VariantSlice []ProductVariant

// ProductAttrs represents product attributes
type ProductAttrs map[string]interface{}

// Dimensions represents product dimensions
type Dimensions struct {
	Length float64 `json:"length"` // in cm
	Width  float64 `json:"width"`  // in cm
	Height float64 `json:"height"` // in cm
	Unit   string  `json:"unit"`   // measurement unit
}

// ProductImage represents a product image
type ProductImage struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	AltText  string `json:"alt_text,omitempty"`
	Position int    `json:"position"`
	IsMain   bool   `json:"is_main"`
	Width    *int   `json:"width,omitempty"`
	Height   *int   `json:"height,omitempty"`
	FileSize *int64 `json:"file_size,omitempty"`
}

// ProductVariant represents a product variant (size, color, etc.)
type ProductVariant struct {
	ID           string            `json:"id"`
	SKU          string            `json:"sku"`
	Name         string            `json:"name"`
	Options      map[string]string `json:"options"` // color: red, size: L
	Price        *int64            `json:"price,omitempty"`
	ComparePrice *int64            `json:"compare_price,omitempty"`
	Weight       *float64          `json:"weight,omitempty"`
	Image        *ProductImage     `json:"image,omitempty"`
	Inventory    Inventory         `json:"inventory"`
	IsDefault    bool              `json:"is_default"`
	Position     int               `json:"position"`
}

// Inventory represents inventory information
type Inventory struct {
	TrackQuantity bool   `json:"track_quantity"`
	Quantity      int    `json:"quantity"`
	ReservedQty   int    `json:"reserved_qty"`
	LowStockThreshold *int `json:"low_stock_threshold,omitempty"`
	AllowBackorder    bool `json:"allow_backorder"`
	StockStatus       StockStatus `json:"stock_status"`
	WarehouseLocation *string     `json:"warehouse_location,omitempty"`
}

// StockStatus represents inventory stock status
type StockStatus string

const (
	StockStatusInStock    StockStatus = "in_stock"
	StockStatusOutOfStock StockStatus = "out_of_stock"
	StockStatusLowStock   StockStatus = "low_stock"
	StockStatusBackorder  StockStatus = "backorder"
	StockStatusDiscontinued StockStatus = "discontinued"
)

// SEOInfo represents SEO information
type SEOInfo struct {
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	MetaKeywords    []string `json:"meta_keywords,omitempty"`
	SlugURL         string `json:"slug_url,omitempty"`
	CanonicalURL    string `json:"canonical_url,omitempty"`
	OpenGraphTitle  string `json:"og_title,omitempty"`
	OpenGraphDesc   string `json:"og_description,omitempty"`
	OpenGraphImage  string `json:"og_image,omitempty"`
}

// ShippingInfo represents shipping information
type ShippingInfo struct {
	FreeShipping      bool    `json:"free_shipping"`
	ShippingClass     string  `json:"shipping_class,omitempty"`
	ShippingCost      *int64  `json:"shipping_cost,omitempty"`
	ProcessingTime    *int    `json:"processing_time,omitempty"` // days
	ShippingWeightKg  *float64 `json:"shipping_weight_kg,omitempty"`
	RequiresSignature bool    `json:"requires_signature"`
	FragileItem       bool    `json:"fragile_item"`
	HazardousMaterial bool    `json:"hazardous_material"`
	RestrictedRegions []string `json:"restricted_regions,omitempty"`
}

// ProductMetadata represents additional product metadata
type ProductMetadata struct {
	Manufacturer    string     `json:"manufacturer,omitempty"`
	OriginCountry   string     `json:"origin_country,omitempty"`
	Certification   []string   `json:"certification,omitempty"` // halal, organic, etc.
	Warranty        *Warranty  `json:"warranty,omitempty"`
	ReturnPolicy    *ReturnPolicy `json:"return_policy,omitempty"`
	AgeRestriction  *int       `json:"age_restriction,omitempty"`
	ReviewsSummary  ReviewsSummary `json:"reviews_summary"`
	ViewCount       int        `json:"view_count"`
	PurchaseCount   int        `json:"purchase_count"`
	WishlistCount   int        `json:"wishlist_count"`
	ShareCount      int        `json:"share_count"`
	LastViewedAt    *time.Time `json:"last_viewed_at,omitempty"`
	TrendingScore   float64    `json:"trending_score"`
	QualityScore    float64    `json:"quality_score"`
	LocalizedNames  map[string]string `json:"localized_names,omitempty"` // language -> name
}

// Warranty represents warranty information
type Warranty struct {
	Duration     int    `json:"duration"`      // in months
	Type         string `json:"type"`          // manufacturer, seller, extended
	Description  string `json:"description,omitempty"`
	ContactInfo  string `json:"contact_info,omitempty"`
	CoveredParts []string `json:"covered_parts,omitempty"`
}

// ReturnPolicy represents return policy information
type ReturnPolicy struct {
	Returnable    bool   `json:"returnable"`
	ReturnWindow  int    `json:"return_window"`  // days
	ReturnCost    string `json:"return_cost"`    // free, buyer_pays, seller_pays
	RestockingFee *float64 `json:"restocking_fee,omitempty"`
	Conditions    []string `json:"conditions,omitempty"`
}

// ReviewsSummary represents aggregated review data
type ReviewsSummary struct {
	TotalReviews int     `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
	RatingDistribution map[int]int `json:"rating_distribution"` // rating -> count
	RecentReviews []ReviewSummary `json:"recent_reviews,omitempty"`
}

// ReviewSummary represents a summary of a review
type ReviewSummary struct {
	ID        string    `json:"id"`
	Rating    int       `json:"rating"`
	Title     string    `json:"title,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	Verified  bool      `json:"verified"`
}

// Value implementations for custom types

func (ss StringSlice) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

func (ss *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*ss = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into StringSlice", value)
	}

	return json.Unmarshal(jsonData, ss)
}

func (is ImageSlice) Value() (driver.Value, error) {
	return json.Marshal(is)
}

func (is *ImageSlice) Scan(value interface{}) error {
	if value == nil {
		*is = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into ImageSlice", value)
	}

	return json.Unmarshal(jsonData, is)
}

func (vs VariantSlice) Value() (driver.Value, error) {
	return json.Marshal(vs)
}

func (vs *VariantSlice) Scan(value interface{}) error {
	if value == nil {
		*vs = nil
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into VariantSlice", value)
	}

	return json.Unmarshal(jsonData, vs)
}

func (d Dimensions) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Dimensions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into Dimensions", value)
	}

	return json.Unmarshal(jsonData, d)
}

func (i Inventory) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *Inventory) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into Inventory", value)
	}

	return json.Unmarshal(jsonData, i)
}

func (s SEOInfo) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SEOInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into SEOInfo", value)
	}

	return json.Unmarshal(jsonData, s)
}

func (si ShippingInfo) Value() (driver.Value, error) {
	return json.Marshal(si)
}

func (si *ShippingInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into ShippingInfo", value)
	}

	return json.Unmarshal(jsonData, si)
}

func (pa ProductAttrs) Value() (driver.Value, error) {
	return json.Marshal(pa)
}

func (pa *ProductAttrs) Scan(value interface{}) error {
	if value == nil {
		*pa = make(ProductAttrs)
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into ProductAttrs", value)
	}

	return json.Unmarshal(jsonData, pa)
}

func (pm ProductMetadata) Value() (driver.Value, error) {
	return json.Marshal(pm)
}

func (pm *ProductMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var jsonData []byte
	switch v := value.(type) {
	case []byte:
		jsonData = v
	case string:
		jsonData = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into ProductMetadata", value)
	}

	return json.Unmarshal(jsonData, pm)
}

// Helper functions for enums

func ValidProductCategories() []ProductCategory {
	return []ProductCategory{
		CategoryElectronics, CategoryFashion, CategoryHome, CategoryBeauty,
		CategorySports, CategoryBooks, CategoryFood, CategoryHealth,
		CategoryToys, CategoryAutomotive, CategoryServices, CategoryDigital,
		CategoryTraditionalCrafts, CategoryLocalProducts, CategoryHalal, CategoryOrganic,
	}
}

func ValidProductStatuses() []ProductStatus {
	return []ProductStatus{
		ProductStatusDraft, ProductStatusActive, ProductStatusInactive,
		ProductStatusArchived, ProductStatusDeleted, ProductStatusPending, ProductStatusRejected,
	}
}

func ValidProductVisibilities() []ProductVisibility {
	return []ProductVisibility{
		VisibilityPublic, VisibilityPrivate, VisibilityHidden, VisibilityScheduled,
	}
}

func ValidCurrencies() []Currency {
	return []Currency{
		CurrencyTHB, CurrencySGD, CurrencyIDR, CurrencyMYR,
		CurrencyPHP, CurrencyVND, CurrencyUSD,
	}
}

func ValidStockStatuses() []StockStatus {
	return []StockStatus{
		StockStatusInStock, StockStatusOutOfStock, StockStatusLowStock,
		StockStatusBackorder, StockStatusDiscontinued,
	}
}

// Validation methods

func (pc ProductCategory) IsValid() bool {
	for _, valid := range ValidProductCategories() {
		if pc == valid {
			return true
		}
	}
	return false
}

func (ps ProductStatus) IsValid() bool {
	for _, valid := range ValidProductStatuses() {
		if ps == valid {
			return true
		}
	}
	return false
}

func (pv ProductVisibility) IsValid() bool {
	for _, valid := range ValidProductVisibilities() {
		if pv == valid {
			return true
		}
	}
	return false
}

func (c Currency) IsValid() bool {
	for _, valid := range ValidCurrencies() {
		if c == valid {
			return true
		}
	}
	return false
}

func (ss StockStatus) IsValid() bool {
	for _, valid := range ValidStockStatuses() {
		if ss == valid {
			return true
		}
	}
	return false
}

// Product business logic methods

func (ps ProductStatus) CanBePurchased() bool {
	return ps == ProductStatusActive
}

func (pv ProductVisibility) IsPubliclyVisible() bool {
	return pv == VisibilityPublic
}

func (ss StockStatus) IsAvailable() bool {
	return ss == StockStatusInStock || (ss == StockStatusBackorder)
}

// Main validation method
func (p *Product) Validate() error {
	var errs []string

	// Shop ID validation
	if p.ShopID == uuid.Nil {
		errs = append(errs, "shop_id is required")
	}

	// SKU validation
	if strings.TrimSpace(p.SKU) == "" {
		errs = append(errs, "sku is required")
	}
	if len(p.SKU) > 100 {
		errs = append(errs, "sku cannot exceed 100 characters")
	}

	// Name validation
	if strings.TrimSpace(p.Name) == "" {
		errs = append(errs, "name is required")
	}
	if len(p.Name) > 255 {
		errs = append(errs, "name cannot exceed 255 characters")
	}

	// Category validation
	if !p.Category.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid category: %s", p.Category))
	}

	// Price validation
	if p.Price < 0 {
		errs = append(errs, "price cannot be negative")
	}

	// Currency validation
	if !p.Currency.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid currency: %s", p.Currency))
	}

	// Status validation
	if !p.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid status: %s", p.Status))
	}

	// Visibility validation
	if !p.Visibility.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid visibility: %s", p.Visibility))
	}

	// Compare price validation
	if p.ComparePrice != nil && *p.ComparePrice < p.Price {
		errs = append(errs, "compare_price must be greater than or equal to price")
	}

	// Weight validation
	if p.Weight != nil && *p.Weight < 0 {
		errs = append(errs, "weight cannot be negative")
	}

	// Images validation
	if err := p.validateImages(); err != nil {
		errs = append(errs, err.Error())
	}

	// Variants validation
	if err := p.validateVariants(); err != nil {
		errs = append(errs, err.Error())
	}

	// Inventory validation
	if err := p.validateInventory(); err != nil {
		errs = append(errs, err.Error())
	}

	// SEO validation
	if err := p.validateSEO(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func (p *Product) validateImages() error {
	var mainImageCount int
	imageIDs := make(map[string]bool)

	for _, img := range p.Images {
		// Check for duplicate IDs
		if imageIDs[img.ID] {
			return fmt.Errorf("duplicate image ID: %s", img.ID)
		}
		imageIDs[img.ID] = true

		// Count main images
		if img.IsMain {
			mainImageCount++
		}

		// Validate URL
		if strings.TrimSpace(img.URL) == "" {
			return errors.New("image URL cannot be empty")
		}

		// Validate position
		if img.Position < 0 {
			return errors.New("image position cannot be negative")
		}
	}

	// Should have exactly one main image
	if len(p.Images) > 0 && mainImageCount != 1 {
		return errors.New("exactly one image must be marked as main")
	}

	return nil
}

func (p *Product) validateVariants() error {
	variantSKUs := make(map[string]bool)
	variantIDs := make(map[string]bool)
	var defaultVariantCount int

	for _, variant := range p.Variants {
		// Check for duplicate SKUs
		if variantSKUs[variant.SKU] {
			return fmt.Errorf("duplicate variant SKU: %s", variant.SKU)
		}
		variantSKUs[variant.SKU] = true

		// Check for duplicate IDs
		if variantIDs[variant.ID] {
			return fmt.Errorf("duplicate variant ID: %s", variant.ID)
		}
		variantIDs[variant.ID] = true

		// Count default variants
		if variant.IsDefault {
			defaultVariantCount++
		}

		// Validate variant inventory
		if !variant.Inventory.StockStatus.IsValid() {
			return fmt.Errorf("invalid stock status for variant %s: %s", variant.ID, variant.Inventory.StockStatus)
		}
	}

	// Should have at most one default variant
	if defaultVariantCount > 1 {
		return errors.New("only one variant can be marked as default")
	}

	return nil
}

func (p *Product) validateInventory() error {
	if !p.Inventory.StockStatus.IsValid() {
		return fmt.Errorf("invalid stock status: %s", p.Inventory.StockStatus)
	}

	if p.Inventory.Quantity < 0 {
		return errors.New("inventory quantity cannot be negative")
	}

	if p.Inventory.ReservedQty < 0 {
		return errors.New("reserved quantity cannot be negative")
	}

	if p.Inventory.ReservedQty > p.Inventory.Quantity {
		return errors.New("reserved quantity cannot exceed total quantity")
	}

	if p.Inventory.LowStockThreshold != nil && *p.Inventory.LowStockThreshold < 0 {
		return errors.New("low stock threshold cannot be negative")
	}

	return nil
}

func (p *Product) validateSEO() error {
	seo := p.SEO

	// Meta title length validation
	if len(seo.MetaTitle) > 60 {
		return errors.New("meta title should not exceed 60 characters")
	}

	// Meta description length validation
	if len(seo.MetaDescription) > 160 {
		return errors.New("meta description should not exceed 160 characters")
	}

	// Slug URL validation
	if seo.SlugURL != "" {
		if strings.Contains(seo.SlugURL, " ") {
			return errors.New("slug URL cannot contain spaces")
		}
		if len(seo.SlugURL) > 100 {
			return errors.New("slug URL cannot exceed 100 characters")
		}
	}

	return nil
}

// Product lifecycle methods

func (p *Product) BeforeCreate() error {
	// Generate UUID if not set
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	// Set default values
	if p.Status == "" {
		p.Status = ProductStatusDraft
	}

	if p.Visibility == "" {
		p.Visibility = VisibilityPublic
	}

	// Initialize metadata if empty
	if p.Metadata == (ProductMetadata{}) {
		p.Metadata = ProductMetadata{
			ReviewsSummary: ReviewsSummary{
				RatingDistribution: make(map[int]int),
			},
		}
	}

	// Initialize inventory if empty
	if p.Inventory == (Inventory{}) {
		p.Inventory = Inventory{
			TrackQuantity: true,
			StockStatus:   StockStatusInStock,
		}
	}

	// Generate SEO slug if not provided
	if p.SEO.SlugURL == "" {
		p.SEO.SlugURL = p.generateSlug()
	}

	// Validate before creation
	return p.Validate()
}

func (p *Product) BeforeUpdate() error {
	// Update timestamp
	p.UpdatedAt = time.Now().UTC()

	// Update stock status based on inventory
	p.updateStockStatus()

	// Update published timestamp if becoming active
	if p.Status == ProductStatusActive && p.PublishedAt == nil {
		now := time.Now().UTC()
		p.PublishedAt = &now
	}

	// Validate before update
	return p.Validate()
}

func (p *Product) generateSlug() string {
	// Simple slug generation from product name
	slug := strings.ToLower(p.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (simplified)
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	return slug
}

func (p *Product) updateStockStatus() {
	if !p.Inventory.TrackQuantity {
		return
	}

	availableQty := p.Inventory.Quantity - p.Inventory.ReservedQty

	if availableQty <= 0 {
		if p.Inventory.AllowBackorder {
			p.Inventory.StockStatus = StockStatusBackorder
		} else {
			p.Inventory.StockStatus = StockStatusOutOfStock
		}
	} else if p.Inventory.LowStockThreshold != nil && availableQty <= *p.Inventory.LowStockThreshold {
		p.Inventory.StockStatus = StockStatusLowStock
	} else {
		p.Inventory.StockStatus = StockStatusInStock
	}
}

// Business logic methods

func (p *Product) IsAvailableForPurchase() bool {
	return p.Status.CanBePurchased() &&
		p.Visibility.IsPubliclyVisible() &&
		p.Inventory.StockStatus.IsAvailable()
}

func (p *Product) GetAvailableQuantity() int {
	if !p.Inventory.TrackQuantity {
		return 999999 // Assume unlimited if not tracking
	}
	available := p.Inventory.Quantity - p.Inventory.ReservedQty
	if available < 0 {
		return 0
	}
	return available
}

func (p *Product) HasVariants() bool {
	return len(p.Variants) > 0
}

func (p *Product) GetMainImage() *ProductImage {
	for _, img := range p.Images {
		if img.IsMain {
			return &img
		}
	}
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

func (p *Product) GetDefaultVariant() *ProductVariant {
	for _, variant := range p.Variants {
		if variant.IsDefault {
			return &variant
		}
	}
	if len(p.Variants) > 0 {
		return &p.Variants[0]
	}
	return nil
}

func (p *Product) GetEffectivePrice() int64 {
	if defaultVariant := p.GetDefaultVariant(); defaultVariant != nil && defaultVariant.Price != nil {
		return *defaultVariant.Price
	}
	return p.Price
}

func (p *Product) HasDiscount() bool {
	return p.ComparePrice != nil && *p.ComparePrice > p.Price
}

func (p *Product) GetDiscountPercent() float64 {
	if !p.HasDiscount() {
		return 0
	}
	return float64(*p.ComparePrice-p.Price) / float64(*p.ComparePrice) * 100
}

func (p *Product) IsDigitalProduct() bool {
	return p.IsDigital
}

func (p *Product) RequiresPhysicalShipping() bool {
	return p.RequiresShipping && !p.IsDigital
}

func (p *Product) IsFeatured() bool {
	return p.FeaturedAt != nil && p.FeaturedAt.Before(time.Now().UTC())
}

func (p *Product) IsPublished() bool {
	return p.PublishedAt != nil && p.PublishedAt.Before(time.Now().UTC())
}

// Currency formatting methods
func (c Currency) FormatAmount(amountInCents int64) string {
	symbols := map[Currency]string{
		CurrencyTHB: "฿",
		CurrencySGD: "S$",
		CurrencyIDR: "Rp",
		CurrencyMYR: "RM",
		CurrencyPHP: "₱",
		CurrencyVND: "₫",
		CurrencyUSD: "$",
	}

	symbol := symbols[c]

	// Handle currencies without decimal places
	if c == CurrencyIDR || c == CurrencyVND {
		return fmt.Sprintf("%s %d", symbol, amountInCents)
	}

	// Standard decimal currencies
	whole := amountInCents / 100
	fraction := amountInCents % 100
	return fmt.Sprintf("%s %d.%02d", symbol, whole, fraction)
}

func (p *Product) GetFormattedPrice() string {
	return p.Currency.FormatAmount(p.Price)
}

func (p *Product) GetFormattedComparePrice() string {
	if p.ComparePrice == nil {
		return ""
	}
	return p.Currency.FormatAmount(*p.ComparePrice)
}

// Analytics and metrics methods

func (p *Product) IncrementViewCount() {
	p.Metadata.ViewCount++
	now := time.Now().UTC()
	p.Metadata.LastViewedAt = &now
	p.UpdatedAt = now
}

func (p *Product) IncrementPurchaseCount() {
	p.Metadata.PurchaseCount++
	p.UpdatedAt = time.Now().UTC()
}

func (p *Product) IncrementWishlistCount() {
	p.Metadata.WishlistCount++
	p.UpdatedAt = time.Now().UTC()
}

func (p *Product) UpdateTrendingScore(score float64) {
	p.Metadata.TrendingScore = score
	p.UpdatedAt = time.Now().UTC()
}

// Status management methods

func (p *Product) Publish() error {
	if p.Status == ProductStatusActive {
		return errors.New("product is already published")
	}

	now := time.Now().UTC()
	p.Status = ProductStatusActive
	p.PublishedAt = &now
	p.UpdatedAt = now

	return p.Validate()
}

func (p *Product) Unpublish() error {
	p.Status = ProductStatusInactive
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Product) Archive() error {
	p.Status = ProductStatusArchived
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Product) SetFeatured() error {
	now := time.Now().UTC()
	p.FeaturedAt = &now
	p.UpdatedAt = now
	return nil
}

func (p *Product) UnsetFeatured() error {
	p.FeaturedAt = nil
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// Public API response methods

func (p *Product) ToPublicProduct() map[string]interface{} {
	response := map[string]interface{}{
		"id":                p.ID,
		"sku":               p.SKU,
		"name":              p.Name,
		"description":       p.Description,
		"short_description": p.ShortDescription,
		"category":          p.Category,
		"tags":              p.Tags,
		"brand":             p.Brand,
		"model":             p.Model,
		"price":             p.Price,
		"formatted_price":   p.GetFormattedPrice(),
		"currency":          p.Currency,
		"images":            p.Images,
		"status":            p.Status,
		"is_digital":        p.IsDigital,
		"requires_shipping": p.RequiresShipping,
		"created_at":        p.CreatedAt,
		"updated_at":        p.UpdatedAt,
	}

	// Add discount information
	if p.HasDiscount() {
		response["compare_price"] = p.ComparePrice
		response["formatted_compare_price"] = p.GetFormattedComparePrice()
		response["discount_percent"] = p.GetDiscountPercent()
	}

	// Add inventory information
	response["inventory"] = map[string]interface{}{
		"stock_status":      p.Inventory.StockStatus,
		"available_quantity": p.GetAvailableQuantity(),
		"allow_backorder":   p.Inventory.AllowBackorder,
	}

	// Add variants if present
	if p.HasVariants() {
		response["variants"] = p.Variants
		response["has_variants"] = true
	}

	// Add main image
	if mainImage := p.GetMainImage(); mainImage != nil {
		response["main_image"] = mainImage
	}

	// Add review summary
	response["reviews"] = p.Metadata.ReviewsSummary

	// Add SEO information
	response["seo"] = p.SEO

	return response
}

func (p *Product) ToAdminProduct() map[string]interface{} {
	response := p.ToPublicProduct()

	// Add admin-only fields
	response["shop_id"] = p.ShopID
	response["cost_price"] = p.CostPrice
	response["weight"] = p.Weight
	response["dimensions"] = p.Dimensions
	response["visibility"] = p.Visibility
	response["is_taxable"] = p.IsTaxable
	response["shipping"] = p.Shipping
	response["attributes"] = p.Attributes
	response["metadata"] = p.Metadata
	response["featured_at"] = p.FeaturedAt
	response["published_at"] = p.PublishedAt

	// Add detailed inventory
	response["inventory"] = p.Inventory

	return response
}

// Request/Response structures

type ProductCreateRequest struct {
	ShopID           uuid.UUID         `json:"shop_id" validate:"required"`
	SKU              string            `json:"sku" validate:"required,max=100"`
	Name             string            `json:"name" validate:"required,max=255"`
	Description      *string           `json:"description,omitempty"`
	ShortDescription *string           `json:"short_description,omitempty"`
	Category         ProductCategory   `json:"category" validate:"required"`
	Tags             []string          `json:"tags,omitempty"`
	Brand            *string           `json:"brand,omitempty"`
	Model            *string           `json:"model,omitempty"`
	Price            int64             `json:"price" validate:"required,min=0"`
	Currency         Currency          `json:"currency" validate:"required"`
	ComparePrice     *int64            `json:"compare_price,omitempty"`
	Weight           *float64          `json:"weight,omitempty"`
	Dimensions       *Dimensions       `json:"dimensions,omitempty"`
	Images           []ProductImage    `json:"images,omitempty"`
	Variants         []ProductVariant  `json:"variants,omitempty"`
	Shipping         ShippingInfo      `json:"shipping,omitempty"`
	Attributes       ProductAttrs      `json:"attributes,omitempty"`
	IsDigital        bool              `json:"is_digital"`
	RequiresShipping bool              `json:"requires_shipping"`
	IsTaxable        bool              `json:"is_taxable"`
}

func (req *ProductCreateRequest) ToProduct() *Product {
	return &Product{
		ShopID:           req.ShopID,
		SKU:              req.SKU,
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Category:         req.Category,
		Tags:             req.Tags,
		Brand:            req.Brand,
		Model:            req.Model,
		Price:            req.Price,
		Currency:         req.Currency,
		ComparePrice:     req.ComparePrice,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		Images:           req.Images,
		Variants:         req.Variants,
		Shipping:         req.Shipping,
		Attributes:       req.Attributes,
		IsDigital:        req.IsDigital,
		RequiresShipping: req.RequiresShipping,
		IsTaxable:        req.IsTaxable,
	}
}

// Product manager utility

type ProductManager struct {
	// Add dependencies like database, search engine, etc.
}

func NewProductManager() *ProductManager {
	return &ProductManager{}
}

func (pm *ProductManager) CreateProduct(req *ProductCreateRequest) (*Product, error) {
	product := req.ToProduct()

	if err := product.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("product creation failed: %v", err)
	}

	return product, nil
}