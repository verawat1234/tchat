package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
)

type CartRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, cart *models.Cart) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	GetBySessionID(ctx context.Context, sessionID string) (*models.Cart, error)
	Update(ctx context.Context, cart *models.Cart) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetExpiredCarts(ctx context.Context, cutoffTime time.Time) ([]*models.Cart, error)
	GetAbandonedCarts(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Cart, int64, error)

	// Media-specific operations
	AddMediaItemToCart(ctx context.Context, cartID uuid.UUID, mediaProductID uuid.UUID, quantity int) error
	RemoveMediaItemFromCart(ctx context.Context, cartID uuid.UUID, mediaProductID uuid.UUID) error
	GetMediaItemsInCart(ctx context.Context, cartID uuid.UUID) ([]*models.CartItem, error)
	CheckMediaItemInCart(ctx context.Context, cartID uuid.UUID, mediaProductID uuid.UUID) (bool, error)
	GetCartsWithMediaItems(ctx context.Context, filters MediaCartFilters, pagination models.Pagination) ([]*models.Cart, int64, error)
	GetMediaCartValue(ctx context.Context, cartID uuid.UUID) (float64, error)
	GetMediaCartsByContentType(ctx context.Context, contentType string, pagination models.Pagination) ([]*models.Cart, int64, error)
	GetAbandonedMediaCarts(ctx context.Context, filters MediaCartFilters, pagination models.Pagination) ([]*models.Cart, int64, error)
	GetMediaCartAnalytics(ctx context.Context, dateFrom, dateTo time.Time) (*MediaCartAnalytics, error)
}

// MediaCartFilters represents filters for carts containing media items
type MediaCartFilters struct {
	UserID      *uuid.UUID
	ContentType *string
	CategoryID  *string
	MinValue    *float64
	MaxValue    *float64
	DateFrom    *time.Time
	DateTo      *time.Time
	Status      *models.CartStatus
	HasDigitalItems bool
}

// MediaCartAnalytics represents analytics data for media carts
type MediaCartAnalytics struct {
	TotalCarts          int64                    `json:"totalCarts"`
	TotalValue          float64                  `json:"totalValue"`
	AverageCartValue    float64                  `json:"averageCartValue"`
	ConversionRate      float64                  `json:"conversionRate"`
	AbandonmentRate     float64                  `json:"abandonmentRate"`
	TopContentTypes     []ContentTypeStats       `json:"topContentTypes"`
	TopCategories       []CategoryStats          `json:"topCategories"`
	DailyStats          []DailyCartStats         `json:"dailyStats"`
}

// ContentTypeStats represents statistics by content type
type ContentTypeStats struct {
	ContentType string  `json:"contentType"`
	CartCount   int64   `json:"cartCount"`
	TotalValue  float64 `json:"totalValue"`
	Percentage  float64 `json:"percentage"`
}

// CategoryStats represents statistics by category
type CategoryStats struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	CartCount    int64   `json:"cartCount"`
	TotalValue   float64 `json:"totalValue"`
	Percentage   float64 `json:"percentage"`
}

// DailyCartStats represents daily cart statistics
type DailyCartStats struct {
	Date       time.Time `json:"date"`
	CartCount  int64     `json:"cartCount"`
	TotalValue float64   `json:"totalValue"`
}

type CartAbandonmentRepository interface {
	Create(ctx context.Context, abandonment *models.CartAbandonmentTracking) error
	GetByCartID(ctx context.Context, cartID uuid.UUID) (*models.CartAbandonmentTracking, error)
	Update(ctx context.Context, abandonment *models.CartAbandonmentTracking) error
	List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.CartAbandonmentTracking, int64, error)
	GetUnrecoveredAbandoned(ctx context.Context, olderThan time.Time) ([]*models.CartAbandonmentTracking, error)
}

type cartRepository struct {
	db *gorm.DB
}

type cartAbandonmentRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func NewCartAbandonmentRepository(db *gorm.DB) CartAbandonmentRepository {
	return &cartAbandonmentRepository{db: db}
}

// Cart Repository Implementation
func (r *cartRepository) Create(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *cartRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).First(&cart, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, models.CartStatusActive).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetBySessionID(ctx context.Context, sessionID string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).Where("session_id = ? AND status = ?", sessionID, models.CartStatusActive).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Update(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

func (r *cartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Cart{}, "id = ?", id).Error
}

func (r *cartRepository) GetExpiredCarts(ctx context.Context, cutoffTime time.Time) ([]*models.Cart, error) {
	var carts []*models.Cart
	err := r.db.WithContext(ctx).Where("expires_at < ? OR (user_id IS NULL AND updated_at < ?)",
		time.Now(), cutoffTime).Find(&carts).Error
	return carts, err
}

func (r *cartRepository) GetAbandonedCarts(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Cart, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Cart{}).Where("status = ?", models.CartStatusAbandoned)

	// Apply filters
	for key, value := range filters {
		switch key {
		case "min_value":
			query = query.Where("total_amount >= ?", value)
		case "date_from":
			query = query.Where("updated_at >= ?", value)
		case "date_to":
			query = query.Where("updated_at <= ?", value)
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var carts []*models.Cart
	err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&carts).Error
	return carts, total, err
}

// Cart Abandonment Repository Implementation
func (r *cartAbandonmentRepository) Create(ctx context.Context, abandonment *models.CartAbandonmentTracking) error {
	return r.db.WithContext(ctx).Create(abandonment).Error
}

func (r *cartAbandonmentRepository) GetByCartID(ctx context.Context, cartID uuid.UUID) (*models.CartAbandonmentTracking, error) {
	var abandonment models.CartAbandonmentTracking
	err := r.db.WithContext(ctx).Where("cart_id = ?", cartID).First(&abandonment).Error
	if err != nil {
		return nil, err
	}
	return &abandonment, nil
}

func (r *cartAbandonmentRepository) Update(ctx context.Context, abandonment *models.CartAbandonmentTracking) error {
	return r.db.WithContext(ctx).Save(abandonment).Error
}

func (r *cartAbandonmentRepository) List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.CartAbandonmentTracking, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.CartAbandonmentTracking{})

	// Apply filters
	for key, value := range filters {
		switch key {
		case "stage":
			if value != nil {
				query = query.Where("abandonment_stage = ?", value)
			}
		case "date_from":
			if value != nil {
				query = query.Where("abandoned_at >= ?", value)
			}
		case "date_to":
			if value != nil {
				query = query.Where("abandoned_at <= ?", value)
			}
		case "is_recovered":
			if value != nil {
				query = query.Where("is_recovered = ?", value)
			}
		case "emails_sent_min":
			if value != nil {
				query = query.Where("emails_sent >= ?", value)
			}
		case "emails_sent_max":
			if value != nil {
				query = query.Where("emails_sent <= ?", value)
			}
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var abandonments []*models.CartAbandonmentTracking
	err := query.Order("abandoned_at DESC").Offset(offset).Limit(limit).Find(&abandonments).Error
	return abandonments, total, err
}

func (r *cartAbandonmentRepository) GetUnrecoveredAbandoned(ctx context.Context, olderThan time.Time) ([]*models.CartAbandonmentTracking, error) {
	var abandonments []*models.CartAbandonmentTracking
	err := r.db.WithContext(ctx).Where("is_recovered = ? AND abandoned_at < ?", false, olderThan).Find(&abandonments).Error
	return abandonments, err
}

// AddMediaItemToCart adds a media product to a cart
func (r *cartRepository) AddMediaItemToCart(ctx context.Context, cartID uuid.UUID, mediaProductID uuid.UUID, quantity int) error {
	// First, verify the product is a media product
	var product models.Product
	if err := r.db.WithContext(ctx).
		Where("id = ? AND product_type = ?", mediaProductID, models.ProductTypeMedia).
		First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("media product not found")
		}
		return fmt.Errorf("failed to verify media product: %w", err)
	}

	// Verify cart exists
	var cart models.Cart
	if err := r.db.WithContext(ctx).Where("id = ?", cartID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("cart not found")
		}
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if item already exists in cart
	var existingItem models.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, mediaProductID).
		First(&existingItem).Error

	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += quantity
		if err := r.db.WithContext(ctx).Save(&existingItem).Error; err != nil {
			return fmt.Errorf("failed to update cart item quantity: %w", err)
		}
	} else if err == gorm.ErrRecordNotFound {
		// Create new cart item
		cartItem := &models.CartItem{
			ID:        uuid.New(),
			CartID:    cartID,
			ProductID: mediaProductID,
			Quantity:  quantity,
			Price:     decimal.NewFromFloat(product.Price),
			Currency:  product.Currency,
		}

		// Set media-specific properties
		if product.IsMedia() {
			license := "personal"
			cartItem.SetMediaLicense(license)

			// Set download format based on content type from media content
			var mediaContent struct {
				ContentType string `gorm:"column:content_type"`
			}
			if err := r.db.WithContext(ctx).
				Table("media_content_items").
				Select("content_type").
				Where("id = ?", product.MediaContentID).
				First(&mediaContent).Error; err == nil {

				format := getDefaultDownloadFormat(mediaContent.ContentType)
				cartItem.SetDownloadFormat(format)
			}
		}

		if err := r.db.WithContext(ctx).Create(cartItem).Error; err != nil {
			return fmt.Errorf("failed to add media item to cart: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check existing cart item: %w", err)
	}

	// Update cart totals
	return r.updateCartTotals(ctx, cartID)
}

// getDefaultDownloadFormat returns default download format for content type
func getDefaultDownloadFormat(contentType string) string {
	switch contentType {
	case "book":
		return "PDF"
	case "podcast":
		return "MP3"
	case "cartoon", "short_movie", "long_movie":
		return "MP4"
	default:
		return "PDF"
	}
}

// RemoveMediaItemFromCart removes a media product from a cart
func (r *cartRepository) RemoveMediaItemFromCart(ctx context.Context, cartID uuid.UUID, mediaProductID uuid.UUID) error {
	// Delete the cart item
	result := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, mediaProductID).
		Delete(&models.CartItem{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove media item from cart: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("media item not found in cart")
	}

	// Update cart totals
	return r.updateCartTotals(ctx, cartID)
}

// GetMediaItemsInCart retrieves all media items in a cart
func (r *cartRepository) GetMediaItemsInCart(ctx context.Context, cartID uuid.UUID) ([]*models.CartItem, error) {
	var cartItems []*models.CartItem

	if err := r.db.WithContext(ctx).
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("cart_items.cart_id = ? AND products.product_type = ?", cartID, models.ProductTypeMedia).
		Preload("Product").
		Find(&cartItems).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items in cart: %w", err)
	}

	return cartItems, nil
}

// CheckMediaItemInCart checks if a media product is already in the cart
func (r *cartRepository) CheckMediaItemInCart(ctx context.Context, cartID uuid.UUID, mediaProductID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("cart_items.cart_id = ? AND cart_items.product_id = ? AND products.product_type = ?",
			cartID, mediaProductID, models.ProductTypeMedia).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check media item in cart: %w", err)
	}

	return count > 0, nil
}

// GetMediaCartValue calculates the total value of media items in a cart
func (r *cartRepository) GetMediaCartValue(ctx context.Context, cartID uuid.UUID) (float64, error) {
	var totalValue float64

	err := r.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Select("COALESCE(SUM(cart_items.price * cart_items.quantity), 0)").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("cart_items.cart_id = ? AND products.product_type = ?", cartID, models.ProductTypeMedia).
		Scan(&totalValue).Error

	if err != nil {
		return 0, fmt.Errorf("failed to calculate media cart value: %w", err)
	}

	return totalValue, nil
}

// GetCartsWithMediaItems retrieves carts containing media items with filters
func (r *cartRepository) GetCartsWithMediaItems(ctx context.Context, filters MediaCartFilters, pagination models.Pagination) ([]*models.Cart, int64, error) {
	var carts []*models.Cart
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Cart{}).
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("products.product_type = ?", models.ProductTypeMedia)

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("carts.user_id = ?", *filters.UserID)
	}

	if filters.Status != nil {
		query = query.Where("carts.status = ?", *filters.Status)
	}

	if filters.MinValue != nil {
		query = query.Where("carts.total_amount >= ?", *filters.MinValue)
	}

	if filters.MaxValue != nil {
		query = query.Where("carts.total_amount <= ?", *filters.MaxValue)
	}

	if filters.DateFrom != nil {
		query = query.Where("carts.created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("carts.created_at <= ?", *filters.DateTo)
	}

	if filters.ContentType != nil {
		query = query.Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
			Where("media_content_items.content_type = ?", *filters.ContentType)
	}

	if filters.CategoryID != nil {
		query = query.Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
			Where("media_content_items.category_id = ?", *filters.CategoryID)
	}

	// Remove duplicates
	query = query.Distinct("carts.id")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count carts with media items: %w", err)
	}

	// Apply pagination and execute query
	offset := pagination.Page * pagination.PageSize
	if err := query.Order("carts.updated_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&carts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find carts with media items: %w", err)
	}

	return carts, total, nil
}

// GetMediaCartsByContentType retrieves carts containing specific content type
func (r *cartRepository) GetMediaCartsByContentType(ctx context.Context, contentType string, pagination models.Pagination) ([]*models.Cart, int64, error) {
	var carts []*models.Cart
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Cart{}).
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
		Where("products.product_type = ? AND media_content_items.content_type = ?",
			models.ProductTypeMedia, contentType).
		Distinct("carts.id")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count media carts by content type: %w", err)
	}

	// Apply pagination and execute query
	offset := pagination.Page * pagination.PageSize
	if err := query.Order("carts.updated_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&carts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find media carts by content type: %w", err)
	}

	return carts, total, nil
}

// GetAbandonedMediaCarts retrieves abandoned carts containing media items
func (r *cartRepository) GetAbandonedMediaCarts(ctx context.Context, filters MediaCartFilters, pagination models.Pagination) ([]*models.Cart, int64, error) {
	var carts []*models.Cart
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Cart{}).
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("products.product_type = ? AND carts.status = ?", models.ProductTypeMedia, models.CartStatusAbandoned)

	// Apply additional filters
	if filters.UserID != nil {
		query = query.Where("carts.user_id = ?", *filters.UserID)
	}

	if filters.MinValue != nil {
		query = query.Where("carts.total_amount >= ?", *filters.MinValue)
	}

	if filters.MaxValue != nil {
		query = query.Where("carts.total_amount <= ?", *filters.MaxValue)
	}

	if filters.DateFrom != nil {
		query = query.Where("carts.created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("carts.created_at <= ?", *filters.DateTo)
	}

	if filters.ContentType != nil {
		query = query.Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
			Where("media_content_items.content_type = ?", *filters.ContentType)
	}

	if filters.CategoryID != nil {
		query = query.Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
			Where("media_content_items.category_id = ?", *filters.CategoryID)
	}

	// Remove duplicates
	query = query.Distinct("carts.id")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count abandoned media carts: %w", err)
	}

	// Apply pagination and execute query
	offset := pagination.Page * pagination.PageSize
	if err := query.Order("carts.updated_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&carts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find abandoned media carts: %w", err)
	}

	return carts, total, nil
}

// GetMediaCartAnalytics generates comprehensive analytics for media carts
func (r *cartRepository) GetMediaCartAnalytics(ctx context.Context, dateFrom, dateTo time.Time) (*MediaCartAnalytics, error) {
	analytics := &MediaCartAnalytics{}

	// Get total carts and value
	err := r.db.WithContext(ctx).Model(&models.Cart{}).
		Select("COUNT(DISTINCT carts.id) as total_carts, COALESCE(SUM(carts.total_amount), 0) as total_value").
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("products.product_type = ? AND carts.created_at BETWEEN ? AND ?",
			models.ProductTypeMedia, dateFrom, dateTo).
		Scan(&struct {
			TotalCarts int64   `gorm:"column:total_carts"`
			TotalValue float64 `gorm:"column:total_value"`
		}{
			TotalCarts: analytics.TotalCarts,
			TotalValue: analytics.TotalValue,
		}).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get basic analytics: %w", err)
	}

	// Calculate average cart value
	if analytics.TotalCarts > 0 {
		analytics.AverageCartValue = analytics.TotalValue / float64(analytics.TotalCarts)
	}

	// Get conversion rate (completed vs total)
	var completedCarts int64
	err = r.db.WithContext(ctx).Model(&models.Cart{}).
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("products.product_type = ? AND carts.status = ? AND carts.created_at BETWEEN ? AND ?",
			models.ProductTypeMedia, models.CartStatusConverted, dateFrom, dateTo).
		Count(&completedCarts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get conversion analytics: %w", err)
	}

	if analytics.TotalCarts > 0 {
		analytics.ConversionRate = float64(completedCarts) / float64(analytics.TotalCarts) * 100
	}

	// Get abandonment rate
	var abandonedCarts int64
	err = r.db.WithContext(ctx).Model(&models.Cart{}).
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("products.product_type = ? AND carts.status = ? AND carts.created_at BETWEEN ? AND ?",
			models.ProductTypeMedia, models.CartStatusAbandoned, dateFrom, dateTo).
		Count(&abandonedCarts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get abandonment analytics: %w", err)
	}

	if analytics.TotalCarts > 0 {
		analytics.AbandonmentRate = float64(abandonedCarts) / float64(analytics.TotalCarts) * 100
	}

	// Get top content types
	err = r.db.WithContext(ctx).
		Select("media_content_items.content_type, COUNT(DISTINCT carts.id) as cart_count, COALESCE(SUM(carts.total_amount), 0) as total_value").
		Table("carts").
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
		Where("products.product_type = ? AND carts.created_at BETWEEN ? AND ?",
			models.ProductTypeMedia, dateFrom, dateTo).
		Group("media_content_items.content_type").
		Order("cart_count DESC").
		Limit(10).
		Scan(&analytics.TopContentTypes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get content type analytics: %w", err)
	}

	// Calculate percentages for content types
	for i := range analytics.TopContentTypes {
		if analytics.TotalCarts > 0 {
			analytics.TopContentTypes[i].Percentage = float64(analytics.TopContentTypes[i].CartCount) / float64(analytics.TotalCarts) * 100
		}
	}

	// Get top categories
	err = r.db.WithContext(ctx).
		Select("media_content_items.category_id, media_categories.name as category_name, COUNT(DISTINCT carts.id) as cart_count, COALESCE(SUM(carts.total_amount), 0) as total_value").
		Table("carts").
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
		Joins("JOIN media_categories ON media_content_items.category_id = media_categories.id").
		Where("products.product_type = ? AND carts.created_at BETWEEN ? AND ?",
			models.ProductTypeMedia, dateFrom, dateTo).
		Group("media_content_items.category_id, media_categories.name").
		Order("cart_count DESC").
		Limit(10).
		Scan(&analytics.TopCategories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get category analytics: %w", err)
	}

	// Calculate percentages for categories
	for i := range analytics.TopCategories {
		if analytics.TotalCarts > 0 {
			analytics.TopCategories[i].Percentage = float64(analytics.TopCategories[i].CartCount) / float64(analytics.TotalCarts) * 100
		}
	}

	// Get daily stats
	err = r.db.WithContext(ctx).
		Select("DATE(carts.created_at) as date, COUNT(DISTINCT carts.id) as cart_count, COALESCE(SUM(carts.total_amount), 0) as total_value").
		Table("carts").
		Joins("JOIN cart_items ON carts.id = cart_items.cart_id").
		Joins("JOIN products ON cart_items.product_id = products.id").
		Where("products.product_type = ? AND carts.created_at BETWEEN ? AND ?",
			models.ProductTypeMedia, dateFrom, dateTo).
		Group("DATE(carts.created_at)").
		Order("date ASC").
		Scan(&analytics.DailyStats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get daily analytics: %w", err)
	}

	return analytics, nil
}

// updateCartTotals recalculates and updates cart totals
func (r *cartRepository) updateCartTotals(ctx context.Context, cartID uuid.UUID) error {
	var totals struct {
		TotalAmount float64
		ItemCount   int
	}

	err := r.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Select("COALESCE(SUM(price * quantity), 0) as total_amount, COALESCE(SUM(quantity), 0) as item_count").
		Where("cart_id = ?", cartID).
		Scan(&totals).Error

	if err != nil {
		return fmt.Errorf("failed to calculate cart totals: %w", err)
	}

	// Update cart
	result := r.db.WithContext(ctx).
		Model(&models.Cart{}).
		Where("id = ?", cartID).
		Updates(map[string]interface{}{
			"total_amount": totals.TotalAmount,
			"item_count":   totals.ItemCount,
			"updated_at":   time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update cart totals: %w", result.Error)
	}

	return nil
}