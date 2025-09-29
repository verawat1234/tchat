package repository

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
	// Basic CRUD operations
	FindProducts(ctx context.Context, filters models.ProductFilters, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error)
	FindProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindProductsByBusinessID(ctx context.Context, businessID uuid.UUID, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	// Media-specific operations
	FindMediaProducts(ctx context.Context, filters MediaProductFilters, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error)
	FindProductByMediaContentID(ctx context.Context, mediaContentID uuid.UUID) (*models.Product, error)
	CreateMediaProduct(ctx context.Context, mediaContentID uuid.UUID, price float64, currency string) (*models.Product, error)
	GetMediaProductsByCategory(ctx context.Context, categoryID string, pagination models.Pagination) ([]*models.Product, int64, error)
	GetMediaProductsByContentType(ctx context.Context, contentType string, pagination models.Pagination) ([]*models.Product, int64, error)
	SearchMediaProducts(ctx context.Context, query string, filters MediaProductFilters, pagination models.Pagination) ([]*models.Product, int64, error)
}

// MediaProductFilters represents filters specifically for media products
type MediaProductFilters struct {
	CategoryID    *string
	ContentType   *string // book, podcast, cartoon, short_movie, long_movie
	MinPrice      *float64
	MaxPrice      *float64
	IsFeatured    *bool
	IsAvailable   *bool
	Search        *string
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

// FindMediaProducts retrieves media products with specific filters
func (r *productRepository) FindMediaProducts(ctx context.Context, filters MediaProductFilters, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{}).
		Where("product_type = ? AND deleted_at IS NULL", models.ProductTypeMedia)

	// Join with media content for filtering
	if filters.CategoryID != nil || filters.ContentType != nil || filters.IsFeatured != nil || filters.IsAvailable != nil {
		query = query.Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id")
	}

	// Apply filters
	if filters.CategoryID != nil {
		query = query.Where("media_content_items.category_id = ?", *filters.CategoryID)
	}

	if filters.ContentType != nil {
		query = query.Where("media_content_items.content_type = ?", *filters.ContentType)
	}

	if filters.MinPrice != nil {
		query = query.Where("products.price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("products.price <= ?", *filters.MaxPrice)
	}

	if filters.IsFeatured != nil {
		query = query.Where("media_content_items.is_featured = ?", *filters.IsFeatured)
	}

	if filters.IsAvailable != nil {
		if *filters.IsAvailable {
			query = query.Where("media_content_items.availability_status = ?", "available")
		} else {
			query = query.Where("media_content_items.availability_status != ?", "available")
		}
	}

	if filters.Search != nil {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ? OR media_content_items.title ILIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count media products: %w", err)
	}

	// Apply sorting
	orderBy := "products.created_at DESC"
	if sort.Field != "" {
		order := "DESC"
		if strings.ToLower(sort.Order) == "asc" {
			order = "ASC"
		}
		orderBy = fmt.Sprintf("products.%s %s", sort.Field, order)
	}
	query = query.Order(orderBy)

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Offset(offset).Limit(pagination.PageSize)

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find media products: %w", err)
	}

	return products, total, nil
}

// FindProductByMediaContentID retrieves a product by its media content ID
func (r *productRepository) FindProductByMediaContentID(ctx context.Context, mediaContentID uuid.UUID) (*models.Product, error) {
	var product models.Product

	if err := r.db.WithContext(ctx).
		Where("media_content_id = ? AND product_type = ? AND deleted_at IS NULL",
			mediaContentID, models.ProductTypeMedia).
		First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media product not found")
		}
		return nil, fmt.Errorf("failed to find media product: %w", err)
	}

	return &product, nil
}

// CreateMediaProduct creates a new media product from media content
func (r *productRepository) CreateMediaProduct(ctx context.Context, mediaContentID uuid.UUID, price float64, currency string) (*models.Product, error) {
	// First, get the media content to populate product fields
	var mediaContent struct {
		ID           uuid.UUID `gorm:"column:id"`
		Title        string    `gorm:"column:title"`
		Description  string    `gorm:"column:description"`
		ThumbnailURL string    `gorm:"column:thumbnail_url"`
	}

	if err := r.db.WithContext(ctx).
		Table("media_content_items").
		Select("id, title, description, thumbnail_url").
		Where("id = ?", mediaContentID).
		First(&mediaContent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media content not found")
		}
		return nil, fmt.Errorf("failed to find media content: %w", err)
	}

	// Create the product
	product := &models.Product{
		ID:             uuid.New(),
		Name:           mediaContent.Title,
		Description:    mediaContent.Description,
		Price:          price,
		Currency:       currency,
		ProductType:    models.ProductTypeMedia,
		MediaContentID: &mediaContentID,
		ThumbnailURL:   mediaContent.ThumbnailURL,
		IsActive:       true,
		StockQuantity:  nil, // Media products don't have stock
	}

	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, fmt.Errorf("failed to create media product: %w", err)
	}

	return product, nil
}

// GetMediaProductsByCategory retrieves media products by category
func (r *productRepository) GetMediaProductsByCategory(ctx context.Context, categoryID string, pagination models.Pagination) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{}).
		Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
		Where("products.product_type = ? AND products.deleted_at IS NULL AND media_content_items.category_id = ?",
			models.ProductTypeMedia, categoryID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count media products by category: %w", err)
	}

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Order("products.created_at DESC").Offset(offset).Limit(pagination.PageSize)

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find media products by category: %w", err)
	}

	return products, total, nil
}

// GetMediaProductsByContentType retrieves media products by content type
func (r *productRepository) GetMediaProductsByContentType(ctx context.Context, contentType string, pagination models.Pagination) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{}).
		Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
		Where("products.product_type = ? AND products.deleted_at IS NULL AND media_content_items.content_type = ?",
			models.ProductTypeMedia, contentType)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count media products by content type: %w", err)
	}

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Order("products.created_at DESC").Offset(offset).Limit(pagination.PageSize)

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find media products by content type: %w", err)
	}

	return products, total, nil
}

// SearchMediaProducts performs text search on media products
func (r *productRepository) SearchMediaProducts(ctx context.Context, query string, filters MediaProductFilters, pagination models.Pagination) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&models.Product{}).
		Joins("JOIN media_content_items ON products.media_content_id = media_content_items.id").
		Where("products.product_type = ? AND products.deleted_at IS NULL", models.ProductTypeMedia)

	// Apply search term
	if query != "" {
		searchTerm := "%" + strings.ToLower(query) + "%"
		dbQuery = dbQuery.Where(
			"LOWER(products.name) LIKE ? OR LOWER(products.description) LIKE ? OR LOWER(media_content_items.title) LIKE ? OR LOWER(media_content_items.description) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// Apply additional filters
	if filters.CategoryID != nil {
		dbQuery = dbQuery.Where("media_content_items.category_id = ?", *filters.CategoryID)
	}

	if filters.ContentType != nil {
		dbQuery = dbQuery.Where("media_content_items.content_type = ?", *filters.ContentType)
	}

	if filters.MinPrice != nil {
		dbQuery = dbQuery.Where("products.price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		dbQuery = dbQuery.Where("products.price <= ?", *filters.MaxPrice)
	}

	if filters.IsFeatured != nil {
		dbQuery = dbQuery.Where("media_content_items.is_featured = ?", *filters.IsFeatured)
	}

	if filters.IsAvailable != nil {
		if *filters.IsAvailable {
			dbQuery = dbQuery.Where("media_content_items.availability_status = ?", "available")
		} else {
			dbQuery = dbQuery.Where("media_content_items.availability_status != ?", "available")
		}
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	dbQuery = dbQuery.Order("products.created_at DESC").Offset(offset).Limit(pagination.PageSize)

	// Execute query
	if err := dbQuery.Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search media products: %w", err)
	}

	return products, total, nil
}