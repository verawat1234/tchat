package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	sharedModels "tchat.dev/shared/models"
)

// Business-related types
type BusinessVerificationStatus = sharedModels.BusinessVerificationStatus
type BusinessContactInfo = sharedModels.BusinessContactInfo
type BusinessAddress = sharedModels.BusinessAddress
type BusinessSettings = sharedModels.BusinessSettings

// Product-related types
type ProductStatus = sharedModels.ProductStatus
type ProductInventory = sharedModels.ProductInventory
type ProductShipping = sharedModels.ProductShipping
type ProductSEO = sharedModels.ProductSEO

// Payment and fulfillment status from shared models
type PaymentStatus = sharedModels.PaymentStatus
type FulfillmentStatus = sharedModels.FulfillmentStatus
type ShippingAddress = sharedModels.ShippingAddress
type BillingAddress = sharedModels.BillingAddress
type ShippingMethod = sharedModels.ShippingMethod
type Event = sharedModels.Event

// Request/Response types
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type SortOptions struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// Business filters
type BusinessFilters struct {
	Country  *string                         `json:"country,omitempty"`
	Category *string                         `json:"category,omitempty"`
	Status   *BusinessVerificationStatus    `json:"status,omitempty"`
	Search   *string                         `json:"search,omitempty"`
}

// Product filters
type ProductFilters struct {
	BusinessID *uuid.UUID     `json:"businessId,omitempty"`
	Category   *string        `json:"category,omitempty"`
	Type       *ProductType   `json:"type,omitempty"`
	Status     *ProductStatus `json:"status,omitempty"`
	Search     *string        `json:"search,omitempty"`
}

// Business request types
type CreateBusinessRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Address     BusinessAddress `json:"address" validate:"required"`
	Contact     BusinessContactInfo `json:"contact" validate:"required"`
	Settings    BusinessSettings    `json:"settings"`
}

type UpdateBusinessRequest struct {
	Name        *string              `json:"name,omitempty"`
	Description *string              `json:"description,omitempty"`
	Category    *string              `json:"category,omitempty"`
	Address     *BusinessAddress     `json:"address,omitempty"`
	Contact     *BusinessContactInfo `json:"contact,omitempty"`
	Settings    *BusinessSettings    `json:"settings,omitempty"`
	IsActive    *bool                `json:"isActive,omitempty"`
}

// Product request types
type CreateProductRequest struct {
	BusinessID  uuid.UUID       `json:"businessId" validate:"required"`
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
	Type        ProductType     `json:"type"`
	Status      ProductStatus   `json:"status"`
	Price       decimal.Decimal `json:"price" validate:"required"`
	Currency    string          `json:"currency" validate:"required"`
	SKU         string          `json:"sku"`
	Images      []string        `json:"images"`
	Tags        []string        `json:"tags"`
	Variants    map[string]interface{} `json:"variants"`
	Inventory   ProductInventory       `json:"inventory"`
	Shipping    ProductShipping        `json:"shipping"`
	SEO         ProductSEO             `json:"seo"`
}

type UpdateProductRequest struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Category    *string                 `json:"category,omitempty"`
	Type        *ProductType            `json:"type,omitempty"`
	Status      *ProductStatus          `json:"status,omitempty"`
	Price       *decimal.Decimal        `json:"price,omitempty"`
	Currency    *string                 `json:"currency,omitempty"`
	SKU         *string                 `json:"sku,omitempty"`
	Images      []string                `json:"images,omitempty"`
	Tags        []string                `json:"tags,omitempty"`
	Variants    map[string]interface{}  `json:"variants,omitempty"`
	Inventory   *ProductInventory       `json:"inventory,omitempty"`
	Shipping    *ProductShipping        `json:"shipping,omitempty"`
	SEO         *ProductSEO             `json:"seo,omitempty"`
	IsActive    *bool                   `json:"isActive,omitempty"`
}

// Response types
type BusinessResponse struct {
	Businesses []*sharedModels.Business `json:"businesses"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"pageSize"`
	TotalPages int64                    `json:"totalPages"`
}

type ProductResponse struct {
	Products   []*sharedModels.Product `json:"products"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"pageSize"`
	TotalPages int64                   `json:"totalPages"`
}

// Review request types
type CreateReviewRequest struct {
	Type       ReviewType      `json:"type" validate:"required"`
	ProductID  *uuid.UUID      `json:"productId,omitempty"`
	BusinessID *uuid.UUID      `json:"businessId,omitempty"`
	OrderID    *uuid.UUID      `json:"orderId,omitempty"`
	UserID     uuid.UUID       `json:"userId" validate:"required"`
	UserName   string          `json:"userName" validate:"required"`
	UserEmail  string          `json:"userEmail" validate:"required,email"`
	Rating     decimal.Decimal `json:"rating" validate:"required"`
	Title      string          `json:"title"`
	Content    string          `json:"content" validate:"required"`
	Images     []ReviewImage   `json:"images,omitempty"`
}

type UpdateReviewRequest struct {
	Rating  *decimal.Decimal `json:"rating,omitempty"`
	Title   *string          `json:"title,omitempty"`
	Content *string          `json:"content,omitempty"`
	Images  []ReviewImage    `json:"images,omitempty"`
}

type ReviewResponse struct {
	Reviews    []*Review `json:"reviews"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"pageSize"`
	TotalPages int64     `json:"totalPages"`
}

// Wishlist request types
type CreateWishlistRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Type        WishlistType    `json:"type"`
	Privacy     WishlistPrivacy `json:"privacy"`
}

type UpdateWishlistRequest struct {
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Privacy     *WishlistPrivacy `json:"privacy,omitempty"`
}

type AddToWishlistRequest struct {
	ProductID uuid.UUID  `json:"productId" validate:"required"`
	VariantID *uuid.UUID `json:"variantId,omitempty"`
	Quantity  int        `json:"quantity"`
	Note      string     `json:"note"`
	Priority  int        `json:"priority"`
}

type WishlistResponse struct {
	Wishlists  []*Wishlist `json:"wishlists"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int64       `json:"totalPages"`
}

// Cart request types
type AddToCartRequest struct {
	ProductID    uuid.UUID  `json:"productId" validate:"required"`
	VariantID    *uuid.UUID `json:"variantId,omitempty"`
	Quantity     int        `json:"quantity" validate:"required,min=1"`
	IsGift       bool       `json:"isGift"`
	GiftMessage  string     `json:"giftMessage"`
}

type UpdateCartItemRequest struct {
	Quantity    *int    `json:"quantity,omitempty"`
	IsGift      *bool   `json:"isGift,omitempty"`
	GiftMessage *string `json:"giftMessage,omitempty"`
}

type ApplyCouponRequest struct {
	CouponCode string `json:"couponCode" validate:"required"`
}

type MergeCartRequest struct {
	UserID    uuid.UUID `json:"userId" validate:"required"`
	SessionID string    `json:"sessionId" validate:"required"`
}

type CreateAbandonmentTrackingRequest struct {
	CartID              uuid.UUID `json:"cartId" validate:"required"`
	AbandonmentStage    string    `json:"abandonmentStage"`
	LastPageVisited     string    `json:"lastPageVisited"`
}

// CartValidationIssue represents an issue found during cart validation
type CartValidationIssue struct {
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	ProductID   uuid.UUID `json:"productId"`
	Severity    string    `json:"severity"`  // "error", "warning", "info"
}

// CartValidation represents the result of cart validation
type CartValidation struct {
	IsValid           bool                   `json:"isValid"`
	Issues            []CartValidationIssue  `json:"issues"`
	TotalItems        int                    `json:"totalItems"`
	TotalValue        decimal.Decimal        `json:"totalValue"`
	Currency          string                 `json:"currency"`
	EstimatedShipping decimal.Decimal        `json:"estimatedShipping"`
	EstimatedTax      decimal.Decimal        `json:"estimatedTax"`
	EstimatedTotal    decimal.Decimal        `json:"estimatedTotal"`
	ValidatedAt       time.Time              `json:"validatedAt"`
}

type CartResponse struct {
	Carts      []*Cart `json:"carts"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"pageSize"`
	TotalPages int64   `json:"totalPages"`
}

type AbandonmentTrackingResponse struct {
	Tracking   []*CartAbandonmentTracking `json:"tracking"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	PageSize   int                        `json:"pageSize"`
	TotalPages int64                      `json:"totalPages"`
}

// Category request types
type CreateCategoryRequest struct {
	BusinessID       *uuid.UUID       `json:"businessId,omitempty"`
	ParentID         *uuid.UUID       `json:"parentId,omitempty"`
	Name             string           `json:"name" validate:"required"`
	Description      string           `json:"description"`
	ShortDescription string           `json:"shortDescription"`
	Type             CategoryType     `json:"type"`
	Icon             string           `json:"icon"`
	Image            *CategoryImage   `json:"image,omitempty"`
	Color            string           `json:"color"`
	SortOrder        int              `json:"sortOrder"`
	IsVisible        bool             `json:"isVisible"`
	IsFeatured       bool             `json:"isFeatured"`
	AllowProducts    bool             `json:"allowProducts"`
	SEO              CategorySEO      `json:"seo"`
	Attributes       []CategoryAttribute `json:"attributes,omitempty"`
}

type UpdateCategoryRequest struct {
	Name             *string          `json:"name,omitempty"`
	Description      *string          `json:"description,omitempty"`
	ShortDescription *string          `json:"shortDescription,omitempty"`
	Status           *CategoryStatus  `json:"status,omitempty"`
	Type             *CategoryType    `json:"type,omitempty"`
	Icon             *string          `json:"icon,omitempty"`
	Image            *CategoryImage   `json:"image,omitempty"`
	Color            *string          `json:"color,omitempty"`
	SortOrder        *int             `json:"sortOrder,omitempty"`
	IsVisible        *bool            `json:"isVisible,omitempty"`
	IsFeatured       *bool            `json:"isFeatured,omitempty"`
	AllowProducts    *bool            `json:"allowProducts,omitempty"`
	SEO              *CategorySEO     `json:"seo,omitempty"`
	Attributes       []CategoryAttribute `json:"attributes,omitempty"`
}

type CategoryResponse struct {
	Categories []*Category `json:"categories"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int64       `json:"totalPages"`
}

// Order filters and request types
type OrderFilters struct {
	BusinessID        *uuid.UUID     `json:"businessId,omitempty"`
	CustomerID        *uuid.UUID     `json:"customerId,omitempty"`
	Status            *OrderStatus   `json:"status,omitempty"`
	PaymentStatus     *PaymentStatus `json:"paymentStatus,omitempty"`
	FulfillmentStatus *FulfillmentStatus `json:"fulfillmentStatus,omitempty"`
	DateFrom          *time.Time     `json:"dateFrom,omitempty"`
	DateTo            *time.Time     `json:"dateTo,omitempty"`
	Search            *string        `json:"search,omitempty"`
}

type CreateOrderRequest struct {
	BusinessID       uuid.UUID        `json:"businessId" validate:"required"`
	Items            []CartItem       `json:"items" validate:"required"`
	ShippingAddress  ShippingAddress  `json:"shippingAddress" validate:"required"`
	BillingAddress   BillingAddress   `json:"billingAddress" validate:"required"`
	ShippingMethod   ShippingMethod   `json:"shippingMethod"`
	CustomerEmail    string           `json:"customerEmail" validate:"required,email"`
	CustomerPhone    string           `json:"customerPhone"`
	CustomerNotes    string           `json:"customerNotes"`
	CouponCode       string           `json:"couponCode"`
	PaymentGateway   string           `json:"paymentGateway"`
}

type OrderResponse struct {
	Orders     []*Order `json:"orders"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
	TotalPages int64    `json:"totalPages"`
}

// AddProductToCategoryRequest represents request to add product to category
type AddProductToCategoryRequest struct {
	ProductID  uuid.UUID `json:"productId" validate:"required"`
	CategoryID uuid.UUID `json:"categoryId" validate:"required"`
	IsPrimary  bool      `json:"isPrimary"`
	SortOrder  int       `json:"sortOrder"`
}

// TrackCategoryViewRequest represents request to track category view
type TrackCategoryViewRequest struct {
	CategoryID uuid.UUID  `json:"categoryId" validate:"required"`
	UserID     *uuid.UUID `json:"userId,omitempty"`
	SessionID  string     `json:"sessionId"`
	IPAddress  string     `json:"ipAddress"`
	UserAgent  string     `json:"userAgent"`
	Referrer   string     `json:"referrer"`
}