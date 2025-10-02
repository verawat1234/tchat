package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"gorm.io/gorm"
)

// StreamProductRepository defines the interface for stream product data access
type StreamProductRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, product *models.StreamProduct) error
	GetByID(ctx context.Context, productID uuid.UUID) (*models.StreamProduct, error)
	ListByStream(ctx context.Context, streamID uuid.UUID) ([]*models.StreamProduct, error)
	Update(ctx context.Context, productID uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, productID uuid.UUID) error

	// Analytics and metrics
	CountByStream(ctx context.Context, streamID uuid.UUID) (int, error)
	UpdateAnalytics(ctx context.Context, productID uuid.UUID, viewCount, clickCount, purchaseCount int, revenue float64) error
	GetTopProducts(ctx context.Context, streamID uuid.UUID, sortBy string, limit int) ([]*models.StreamProduct, error)

	// Batch operations
	CreateBatch(ctx context.Context, products []*models.StreamProduct) error
	DeleteByStream(ctx context.Context, streamID uuid.UUID) error
}

// streamProductRepository implements the StreamProductRepository interface
type streamProductRepository struct {
	db *gorm.DB
}

// NewStreamProductRepository creates a new stream product repository
func NewStreamProductRepository(db *gorm.DB) StreamProductRepository {
	return &streamProductRepository{
		db: db,
	}
}

// Create creates a new stream product
// Business Rule: Max 10 products per stream
func (r *streamProductRepository) Create(ctx context.Context, product *models.StreamProduct) error {
	// Validate max products per stream
	count, err := r.CountByStream(ctx, product.StreamID)
	if err != nil {
		return fmt.Errorf("failed to count existing products: %w", err)
	}

	if count >= 10 {
		return fmt.Errorf("maximum 10 products per stream exceeded (current: %d)", count)
	}

	// Create the product
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create stream product: %w", err)
	}

	return nil
}

// GetByID retrieves a stream product by its ID
func (r *streamProductRepository) GetByID(ctx context.Context, productID uuid.UUID) (*models.StreamProduct, error) {
	var product models.StreamProduct

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("stream product not found")
		}
		return nil, fmt.Errorf("failed to find stream product: %w", err)
	}

	return &product, nil
}

// ListByStream retrieves all products for a specific stream
// Ordered by display_priority DESC, featured_at ASC
func (r *streamProductRepository) ListByStream(ctx context.Context, streamID uuid.UUID) ([]*models.StreamProduct, error) {
	var products []*models.StreamProduct

	if err := r.db.WithContext(ctx).
		Where("stream_id = ? AND deleted_at IS NULL", streamID).
		Order("display_priority DESC, featured_at ASC").
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to list stream products: %w", err)
	}

	return products, nil
}

// Update updates a stream product
func (r *streamProductRepository) Update(ctx context.Context, productID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&models.StreamProduct{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update stream product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream product not found")
	}

	return nil
}

// Delete soft deletes a stream product
func (r *streamProductRepository) Delete(ctx context.Context, productID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", productID).
		Delete(&models.StreamProduct{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete stream product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream product not found")
	}

	return nil
}

// CountByStream counts active products for a stream
func (r *streamProductRepository) CountByStream(ctx context.Context, streamID uuid.UUID) (int, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.StreamProduct{}).
		Where("stream_id = ? AND deleted_at IS NULL", streamID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count stream products: %w", err)
	}

	return int(count), nil
}

// UpdateAnalytics updates analytics metrics for a stream product
// This method is called by analytics service to update engagement metrics
func (r *streamProductRepository) UpdateAnalytics(ctx context.Context, productID uuid.UUID, viewCount, clickCount, purchaseCount int, revenue float64) error {
	updates := map[string]interface{}{
		"view_count":       gorm.Expr("view_count + ?", viewCount),
		"click_count":      gorm.Expr("click_count + ?", clickCount),
		"purchase_count":   gorm.Expr("purchase_count + ?", purchaseCount),
		"revenue_generated": gorm.Expr("revenue_generated + ?", revenue),
	}

	result := r.db.WithContext(ctx).
		Model(&models.StreamProduct{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update stream product analytics: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("stream product not found")
	}

	return nil
}

// GetTopProducts retrieves top-performing products for a stream
// sortBy options: "views", "clicks", "purchases", "revenue"
func (r *streamProductRepository) GetTopProducts(ctx context.Context, streamID uuid.UUID, sortBy string, limit int) ([]*models.StreamProduct, error) {
	var products []*models.StreamProduct

	// Validate sortBy parameter
	validSortFields := map[string]string{
		"views":     "view_count",
		"clicks":    "click_count",
		"purchases": "purchase_count",
		"revenue":   "revenue_generated",
	}

	sortField, ok := validSortFields[strings.ToLower(sortBy)]
	if !ok {
		sortField = "view_count" // Default to views
	}

	// Validate limit
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 10 {
		limit = 10 // Maximum limit (all products)
	}

	query := r.db.WithContext(ctx).
		Where("stream_id = ? AND deleted_at IS NULL", streamID).
		Order(fmt.Sprintf("%s DESC, featured_at ASC", sortField)).
		Limit(limit)

	if err := query.Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get top stream products: %w", err)
	}

	return products, nil
}

// CreateBatch creates multiple stream products in a single transaction
// Business Rule: Total products (existing + new) must not exceed 10
func (r *streamProductRepository) CreateBatch(ctx context.Context, products []*models.StreamProduct) error {
	if len(products) == 0 {
		return nil
	}

	// Validate all products belong to the same stream
	streamID := products[0].StreamID
	for _, p := range products {
		if p.StreamID != streamID {
			return fmt.Errorf("all products in batch must belong to the same stream")
		}
	}

	// Validate max products per stream
	existingCount, err := r.CountByStream(ctx, streamID)
	if err != nil {
		return fmt.Errorf("failed to count existing products: %w", err)
	}

	if existingCount+len(products) > 10 {
		return fmt.Errorf("maximum 10 products per stream exceeded (existing: %d, new: %d)", existingCount, len(products))
	}

	// Create products in transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, product := range products {
			if err := tx.Create(product).Error; err != nil {
				return fmt.Errorf("failed to create stream product in batch: %w", err)
			}
		}
		return nil
	})
}

// DeleteByStream deletes all products for a specific stream
// Used when a stream is deleted or ended
func (r *streamProductRepository) DeleteByStream(ctx context.Context, streamID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("stream_id = ?", streamID).
		Delete(&models.StreamProduct{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete stream products: %w", result.Error)
	}

	return nil
}