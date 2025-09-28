package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
)

// productRepository implements the ProductRepository interface
type productRepository struct {
	db *gorm.DB
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	FindProducts(ctx context.Context, filters models.ProductFilters, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error)
	FindProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindProductsByBusinessID(ctx context.Context, businessID uuid.UUID, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

// FindProducts retrieves products with filters and pagination
func (r *productRepository) FindProducts(ctx context.Context, filters models.ProductFilters, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{})

	// Apply filters
	if filters.BusinessID != nil {
		query = query.Where("business_id = ?", *filters.BusinessID)
	}

	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Search != nil {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Only get active products by default
	query = query.Where("deleted_at IS NULL")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if sort.Field != "" {
		order := "DESC"
		if strings.ToLower(sort.Order) == "asc" {
			order = "ASC"
		}
		orderBy = fmt.Sprintf("%s %s", sort.Field, order)
	}
	query = query.Order(orderBy)

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Offset(offset).Limit(pagination.PageSize)

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	return products, total, nil
}

// FindProductByID retrieves a product by its ID
func (r *productRepository) FindProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product

	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	return &product, nil
}

// FindProductsByBusinessID retrieves products for a specific business
func (r *productRepository) FindProductsByBusinessID(ctx context.Context, businessID uuid.UUID, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{}).Where("business_id = ? AND deleted_at IS NULL", businessID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if sort.Field != "" {
		order := "DESC"
		if strings.ToLower(sort.Order) == "asc" {
			order = "ASC"
		}
		orderBy = fmt.Sprintf("%s %s", sort.Field, order)
	}
	query = query.Order(orderBy)

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Offset(offset).Limit(pagination.PageSize)

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	return products, total, nil
}

// CreateProduct creates a new product
func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// UpdateProduct updates a product
func (r *productRepository) UpdateProduct(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// DeleteProduct soft deletes a product
func (r *productRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}