package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination models.Pagination, sort models.SortOptions) ([]*models.Order, int64, error)
	Update(ctx context.Context, order *models.Order) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Status operations
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	GetOrdersByStatus(ctx context.Context, status string, pagination models.Pagination) ([]*models.Order, int64, error)

	// Media-specific operations
	GetMediaOrders(ctx context.Context, filters MediaOrderFilters, pagination models.Pagination) ([]*models.Order, int64, error)
	GetOrdersWithMediaContent(ctx context.Context, mediaContentID uuid.UUID, pagination models.Pagination) ([]*models.Order, int64, error)
	GetUserMediaPurchases(ctx context.Context, userID uuid.UUID, contentType string) ([]*models.Order, error)
	GetDigitalDownloads(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
}

// MediaOrderFilters represents filters for media order queries
type MediaOrderFilters struct {
	UserID        *uuid.UUID
	ContentType   *string // book, podcast, cartoon, short_movie, long_movie
	CategoryID    *string
	DateFrom      *time.Time
	DateTo        *time.Time
	MinAmount     *float64
	MaxAmount     *float64
	Status        *string
	HasMediaItems bool
}

// orderRepository implements the OrderRepository interface
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// Create creates a new order
func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// GetByID retrieves an order by its ID
func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order

	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("id = ?", id).
		First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	return &order, nil
}

// GetByUserID retrieves orders for a specific user
func (r *orderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, pagination models.Pagination, sort models.SortOptions) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if sort.Field != "" {
		order := "DESC"
		if sort.Order == "asc" {
			order = "ASC"
		}
		orderBy = fmt.Sprintf("%s %s", sort.Field, order)
	}
	query = query.Order(orderBy)

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Offset(offset).Limit(pagination.PageSize)

	// Execute query with preloads
	if err := query.Preload("Items").Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find orders: %w", err)
	}

	return orders, total, nil
}

// Update updates an order
func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	if err := r.db.WithContext(ctx).Save(order).Error; err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

// Delete soft deletes an order
func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Order{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// UpdateStatus updates an order's status
func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// GetOrdersByStatus retrieves orders by status
func (r *orderRepository) GetOrdersByStatus(ctx context.Context, status string, pagination models.Pagination) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Order("created_at DESC").Offset(offset).Limit(pagination.PageSize)

	// Execute query with preloads
	if err := query.Preload("Items").Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find orders: %w", err)
	}

	return orders, total, nil
}

// GetMediaOrders retrieves orders containing media products with filters
func (r *orderRepository) GetMediaOrders(ctx context.Context, filters MediaOrderFilters, pagination models.Pagination) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	// Build base query for orders with media items
	query := r.db.WithContext(ctx).Model(&models.Order{}).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Where("products.product_type = ?", models.ProductTypeMedia)

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("orders.user_id = ?", *filters.UserID)
	}

	if filters.Status != nil {
		query = query.Where("orders.status = ?", *filters.Status)
	}

	if filters.DateFrom != nil {
		query = query.Where("orders.created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where("orders.created_at <= ?", *filters.DateTo)
	}

	if filters.MinAmount != nil {
		query = query.Where("orders.total_amount >= ?", *filters.MinAmount)
	}

	if filters.MaxAmount != nil {
		query = query.Where("orders.total_amount <= ?", *filters.MaxAmount)
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
	query = query.Distinct("orders.id")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count media orders: %w", err)
	}

	// Apply pagination and execute query
	offset := pagination.Page * pagination.PageSize
	if err := query.Order("orders.created_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Preload("Items").
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find media orders: %w", err)
	}

	return orders, total, nil
}

// GetOrdersWithMediaContent retrieves orders containing a specific media content item
func (r *orderRepository) GetOrdersWithMediaContent(ctx context.Context, mediaContentID uuid.UUID, pagination models.Pagination) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{}).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Where("products.media_content_id = ?", mediaContentID).
		Distinct("orders.id")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders with media content: %w", err)
	}

	// Apply pagination and execute query
	offset := pagination.Page * pagination.PageSize
	if err := query.Order("orders.created_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Preload("Items").
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find orders with media content: %w", err)
	}

	return orders, total, nil
}

// GetUserMediaPurchases retrieves all media purchases for a user, optionally filtered by content type
func (r *orderRepository) GetUserMediaPurchases(ctx context.Context, userID uuid.UUID, contentType string) ([]*models.Order, error) {
	var orders []*models.Order

	query := r.db.WithContext(ctx).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Where("orders.user_id = ? AND products.product_type = ? AND orders.status = ?",
			userID, models.ProductTypeMedia, "completed")

	if contentType != "" {
		query = query.Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
			Where("media_content_items.content_type = ?", contentType)
	}

	if err := query.Distinct("orders.id").
		Order("orders.created_at DESC").
		Preload("Items").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find user media purchases: %w", err)
	}

	return orders, nil
}

// GetDigitalDownloads retrieves all completed media orders for a user (their digital library)
func (r *orderRepository) GetDigitalDownloads(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	var orders []*models.Order

	if err := r.db.WithContext(ctx).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Where("orders.user_id = ? AND products.product_type = ? AND orders.status = ?",
			userID, models.ProductTypeMedia, "completed").
		Distinct("orders.id").
		Order("orders.created_at DESC").
		Preload("Items").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find digital downloads: %w", err)
	}

	return orders, nil
}