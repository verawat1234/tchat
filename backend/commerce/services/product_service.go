package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
	sharedModels "tchat.dev/shared/models"
)

// productService implements the ProductService interface
type productService struct {
	productRepo  repository.ProductRepository
	businessRepo repository.BusinessRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo repository.ProductRepository, businessRepo repository.BusinessRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		businessRepo: businessRepo,
	}
}

// GetProducts retrieves products with filters and pagination
func (s *productService) GetProducts(ctx context.Context, filters models.ProductFilters, pagination models.Pagination, sort models.SortOptions) (*models.ProductResponse, error) {
	products, total, err := s.productRepo.FindProducts(ctx, filters, pagination, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}

	// Products are already using shared models (via type alias)
	return &models.ProductResponse{
		Products:   products,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize),
	}, nil
}

// GetProduct retrieves a single product by ID
func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product, err := s.productRepo.FindProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}

	if req.BusinessID == uuid.Nil {
		return nil, fmt.Errorf("business ID is required")
	}

	if req.Price.LessThan(decimal.Zero) {
		return nil, fmt.Errorf("price must be non-negative")
	}

	// Validate business exists
	_, err := s.businessRepo.FindBusinessByID(ctx, req.BusinessID)
	if err != nil {
		return nil, fmt.Errorf("business not found")
	}

	// Set default values
	status := req.Status
	if status == "" {
		status = models.ProductStatusDraft
	}

	productType := req.Type
	if productType == "" {
		productType = models.ProductTypePhysical
	}

	// Convert simple price/currency to ProductPricing array
	pricing := []models.ProductPricing{
		{
			Currency: req.Currency,
			Price:    req.Price,
		},
	}

	// Convert string images to ProductImage array
	var images []models.ProductImage
	for i, imageURL := range req.Images {
		images = append(images, models.ProductImage{
			URL:       imageURL,
			IsPrimary: i == 0, // First image is primary
			SortOrder: i,
		})
	}

	// Convert variants map to ProductVariant array (simplified for now)
	var variants []models.ProductVariant
	// Note: This is a simplified conversion - in a real implementation,
	// you'd properly convert the map structure to ProductVariant objects

	// Create product entity
	product := &sharedModels.Product{
		ID:          uuid.New(),
		BusinessID:  req.BusinessID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Type:        productType,
		Status:      status,
		SKU:         req.SKU,
		Tags:        req.Tags,
		Pricing:     pricing,
		Images:      images,
		HasVariants: len(variants) > 0,
		Variants:    variants,
		Inventory: models.ProductInventory{
			Stock:             req.Inventory.Stock,
			LowStockThreshold: req.Inventory.LowStockThreshold,
			TrackQuantity:     req.Inventory.TrackQuantity,
			AllowBackorder:    req.Inventory.AllowBackorder,
			ManageStock:       req.Inventory.ManageStock,
		},
		Shipping: models.ProductShipping{
			Weight:           req.Shipping.Weight,
			Length:           req.Shipping.Length,
			Width:            req.Shipping.Width,
			Height:           req.Shipping.Height,
			RequiresShipping: req.Shipping.RequiresShipping,
			ShippingClass:    req.Shipping.ShippingClass,
			FreeShipping:     req.Shipping.FreeShipping,
		},
		SEO: models.ProductSEO{
			MetaTitle:       req.SEO.MetaTitle,
			MetaDescription: req.SEO.MetaDescription,
			URLSlug:         req.SEO.URLSlug,
			Keywords:        req.SEO.Keywords,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	err = s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error) {
	// Check if product exists
	_, err := s.productRepo.FindProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Price != nil && req.Currency != nil {
		if req.Price.LessThan(decimal.Zero) {
			return nil, fmt.Errorf("price must be non-negative")
		}
		// Update pricing array with new price/currency
		pricing := []sharedModels.ProductPricing{
			{
				Currency: *req.Currency,
				Price:    *req.Price,
			},
		}
		updates["pricing"] = pricing
	}
	if req.SKU != nil {
		updates["sku"] = *req.SKU
	}
	if req.Images != nil {
		// Convert string images to ProductImage array
		var images []sharedModels.ProductImage
		for i, imageURL := range req.Images {
			images = append(images, sharedModels.ProductImage{
				URL:       imageURL,
				IsPrimary: i == 0,
				SortOrder: i,
			})
		}
		updates["images"] = images
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Variants != nil {
		// Convert variants map to ProductVariant array (simplified)
		var variants []sharedModels.ProductVariant
		// Note: This is a simplified conversion
		updates["variants"] = variants
		updates["has_variants"] = len(variants) > 0
	}
	if req.Inventory != nil {
		updates["inventory_stock"] = req.Inventory.Stock
		updates["inventory_low_stock_threshold"] = req.Inventory.LowStockThreshold
		updates["inventory_track_quantity"] = req.Inventory.TrackQuantity
		updates["inventory_allow_backorder"] = req.Inventory.AllowBackorder
		updates["inventory_manage_stock"] = req.Inventory.ManageStock
	}
	if req.Shipping != nil {
		updates["shipping_weight"] = req.Shipping.Weight
		updates["shipping_length"] = req.Shipping.Length
		updates["shipping_width"] = req.Shipping.Width
		updates["shipping_height"] = req.Shipping.Height
		updates["shipping_requires_shipping"] = req.Shipping.RequiresShipping
		updates["shipping_class"] = req.Shipping.ShippingClass
		updates["shipping_free_shipping"] = req.Shipping.FreeShipping
	}
	if req.SEO != nil {
		updates["seo_meta_title"] = req.SEO.MetaTitle
		updates["seo_meta_description"] = req.SEO.MetaDescription
		updates["seo_url_slug"] = req.SEO.URLSlug
		updates["seo_keywords"] = req.SEO.Keywords
	}

	updates["updated_at"] = time.Now()

	// Update product
	err = s.productRepo.UpdateProduct(ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Return updated product
	return s.productRepo.FindProductByID(ctx, id)
}

// DeleteProduct deletes a product (soft delete)
func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	// Check if product exists
	_, err := s.productRepo.FindProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	// Soft delete
	err = s.productRepo.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}