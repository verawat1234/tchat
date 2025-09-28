package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
	sharedModels "tchat.dev/shared/models"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req models.CreateCategoryRequest) (*models.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetCategoryByPath(ctx context.Context, path string) (*models.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req models.UpdateCategoryRequest) (*models.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context, filters map[string]interface{}, pagination models.Pagination) (*models.CategoryResponse, error)

	GetCategoryTree(ctx context.Context, businessID *uuid.UUID, maxDepth int) ([]*CategoryNode, error)
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID, recursive bool) ([]*models.Category, error)
	GetCategoryAncestors(ctx context.Context, categoryID uuid.UUID) ([]*models.Category, error)
	MoveCategoryToParent(ctx context.Context, categoryID uuid.UUID, newParentID *uuid.UUID) error

	GetCategoriesByBusiness(ctx context.Context, businessID uuid.UUID, pagination models.Pagination) (*models.CategoryResponse, error)
	GetBusinessCategories(ctx context.Context, businessID uuid.UUID, pagination models.Pagination) (*models.CategoryResponse, error)
	GetGlobalCategories(ctx context.Context, pagination models.Pagination) (*models.CategoryResponse, error)
	GetFeaturedCategories(ctx context.Context, businessID *uuid.UUID, limit int) ([]*models.Category, error)
	GetRootCategories(ctx context.Context, businessID *uuid.UUID) ([]*models.Category, error)

	AddProductToCategory(ctx context.Context, productID, categoryID uuid.UUID, isPrimary bool) error
	RemoveProductFromCategory(ctx context.Context, productID, categoryID uuid.UUID) error
	GetProductCategories(ctx context.Context, productID uuid.UUID) ([]*models.Category, error)
	GetCategoryProducts(ctx context.Context, categoryID uuid.UUID, pagination models.Pagination) ([]uuid.UUID, int64, error)

	UpdateCategoryStats(ctx context.Context, categoryID uuid.UUID) error
	RecalculateAllStats(ctx context.Context) error

	TrackCategoryView(ctx context.Context, categoryID uuid.UUID, userID *uuid.UUID, sessionID, ipAddress, userAgent, referrer string) error
	GetCategoryAnalytics(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (*CategoryAnalytics, error)
}

type CategoryFilters struct {
	BusinessID   *uuid.UUID              `json:"businessId,omitempty"`
	ParentID     *uuid.UUID              `json:"parentId,omitempty"`
	Type         *models.CategoryType    `json:"type,omitempty"`
	Status       *models.CategoryStatus  `json:"status,omitempty"`
	IsVisible    *bool                   `json:"isVisible,omitempty"`
	IsFeatured   *bool                   `json:"isFeatured,omitempty"`
	Level        *int                    `json:"level,omitempty"`
	Search       *string                 `json:"search,omitempty"`
}

type CategoryNode struct {
	*models.Category
	Children []*CategoryNode `json:"children,omitempty"`
}

type CategoryAnalytics struct {
	TotalViews      int64                    `json:"totalViews"`
	UniqueVisitors  int64                    `json:"uniqueVisitors"`
	TopReferrers    []ReferrerStat           `json:"topReferrers"`
	ViewsByDay      []DayViewStat            `json:"viewsByDay"`
	ProductClicks   int64                    `json:"productClicks"`
	ConversionRate  float64                  `json:"conversionRate"`
}

type ReferrerStat struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

type DayViewStat struct {
	Date  string `json:"date"`
	Views int64  `json:"views"`
}

type categoryService struct {
	categoryRepo        repository.CategoryRepository
	productCategoryRepo repository.ProductCategoryRepository
	categoryViewRepo    repository.CategoryViewRepository
	db                  *gorm.DB
}

func NewCategoryService(categoryRepo repository.CategoryRepository, productCategoryRepo repository.ProductCategoryRepository, categoryViewRepo repository.CategoryViewRepository, db *gorm.DB) CategoryService {
	return &categoryService{
		categoryRepo:        categoryRepo,
		productCategoryRepo: productCategoryRepo,
		categoryViewRepo:    categoryViewRepo,
		db:                  db,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, req models.CreateCategoryRequest) (*models.Category, error) {
	// Calculate level and path
	level := 0
	path := req.Name

	if req.ParentID != nil {
		parent, err := s.GetCategory(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent category: %w", err)
		}
		level = parent.Level + 1
		path = parent.Path + "/" + req.Name
	}

	category := &models.Category{
		BusinessID:       req.BusinessID,
		ParentID:         req.ParentID,
		Level:            level,
		Path:             path,
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Type:             req.Type,
		Status:           models.CategoryStatusActive,
		Icon:             req.Icon,
		Image:            req.Image,
		Color:            req.Color,
		SortOrder:        req.SortOrder,
		IsVisible:        req.IsVisible,
		IsFeatured:       req.IsFeatured,
		AllowProducts:    req.AllowProducts,
		SEO:              req.SEO,
		Attributes:       req.Attributes,
	}

	if err := s.db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Update parent's children count if parent exists
	if req.ParentID != nil {
		if err := s.updateChildrenCount(ctx, *req.ParentID); err != nil {
			fmt.Printf("Failed to update parent children count: %v\n", err)
		}
	}

	// Create event
	eventData := map[string]interface{}{
		"category_id":   category.ID,
		"business_id":   category.BusinessID,
		"name":          category.Name,
		"level":         category.Level,
	}

	event := &sharedModels.Event{
		Type:          sharedModels.EventTypeCategoryCreated,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       "Category Created",
		Description:   "A new category has been created",
		AggregateType: "category",
		AggregateID:   category.ID.String(),
	}

	if err := event.MarshalData(eventData); err != nil {
		fmt.Printf("Failed to marshal event data: %v\n", err)
	}

	if err := s.db.WithContext(ctx).Create(event).Error; err != nil {
		fmt.Printf("Failed to create category event: %v\n", err)
	}

	return category, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := s.db.WithContext(ctx).First(&category, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func (s *categoryService) GetCategoryByPath(ctx context.Context, path string) (*models.Category, error) {
	var category models.Category
	if err := s.db.WithContext(ctx).Where("path = ?", path).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category by path: %w", err)
	}

	return &category, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uuid.UUID, req models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	pathChanged := false

	if req.Name != nil && *req.Name != category.Name {
		updates["name"] = *req.Name
		pathChanged = true
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ShortDescription != nil {
		updates["short_description"] = *req.ShortDescription
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.Image != nil {
		updates["image"] = *req.Image
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsVisible != nil {
		updates["is_visible"] = *req.IsVisible
	}
	if req.IsFeatured != nil {
		updates["is_featured"] = *req.IsFeatured
	}
	if req.AllowProducts != nil {
		updates["allow_products"] = *req.AllowProducts
	}
	if req.SEO != nil {
		updates["seo_meta_title"] = req.SEO.MetaTitle
		updates["seo_meta_description"] = req.SEO.MetaDescription
		updates["seo_keywords"] = req.SEO.Keywords
		updates["seo_url_slug"] = req.SEO.URLSlug
		updates["seo_canonical"] = req.SEO.Canonical
	}
	if req.Attributes != nil {
		updates["attributes"] = req.Attributes
	}

	// Update path if name changed
	if pathChanged {
		newPath, err := s.calculatePath(ctx, category.ParentID, *req.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate new path: %w", err)
		}
		updates["path"] = newPath
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(category).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update category: %w", err)
		}

		// Update paths of all descendants if path changed
		if pathChanged {
			if err := s.updateDescendantPaths(ctx, category.ID); err != nil {
				fmt.Printf("Failed to update descendant paths: %v\n", err)
			}
		}
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	category, err := s.GetCategory(ctx, id)
	if err != nil {
		return err
	}

	// Check if category has children
	var childCount int64
	if err := s.db.WithContext(ctx).Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return fmt.Errorf("failed to check children: %w", err)
	}

	if childCount > 0 {
		return fmt.Errorf("cannot delete category with children")
	}

	// Check if category has products
	var productCount int64
	if err := s.db.WithContext(ctx).Model(&models.ProductCategory{}).Where("category_id = ?", id).Count(&productCount).Error; err != nil {
		return fmt.Errorf("failed to check products: %w", err)
	}

	if productCount > 0 {
		return fmt.Errorf("cannot delete category with products")
	}

	// Delete category
	if err := s.db.WithContext(ctx).Delete(category).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	// Update parent's children count if parent exists
	if category.ParentID != nil {
		if err := s.updateChildrenCount(ctx, *category.ParentID); err != nil {
			fmt.Printf("Failed to update parent children count: %v\n", err)
		}
	}

	return nil
}

func (s *categoryService) ListCategories(ctx context.Context, filters map[string]interface{}, pagination models.Pagination) (*models.CategoryResponse, error) {
	query := s.db.WithContext(ctx).Model(&models.Category{})

	// Apply filters
	if businessID, ok := filters["businessId"]; ok && businessID != nil {
		query = query.Where("business_id = ?", businessID)
	}
	if parentID, ok := filters["parentId"]; ok && parentID != nil {
		query = query.Where("parent_id = ?", parentID)
	}
	if categoryType, ok := filters["type"]; ok && categoryType != nil {
		query = query.Where("type = ?", categoryType)
	}
	if status, ok := filters["status"]; ok && status != nil {
		query = query.Where("status = ?", status)
	}
	if isVisible, ok := filters["isVisible"]; ok && isVisible != nil {
		query = query.Where("is_visible = ?", isVisible)
	}
	if isFeatured, ok := filters["isFeatured"]; ok && isFeatured != nil {
		query = query.Where("is_featured = ?", isFeatured)
	}
	if level, ok := filters["level"]; ok && level != nil {
		query = query.Where("level = ?", level)
	}
	if search, ok := filters["search"]; ok && search != nil && search.(string) != "" {
		searchTerm := "%" + search.(string) + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count categories: %w", err)
	}

	// Get results
	var categories []*models.Category
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Order("level ASC, sort_order ASC, name ASC").
		Offset(offset).Limit(pagination.PageSize).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	totalPages := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	return &models.CategoryResponse{
		Categories: categories,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *categoryService) GetCategoryTree(ctx context.Context, businessID *uuid.UUID, maxDepth int) ([]*CategoryNode, error) {
	// Get root categories
	query := s.db.WithContext(ctx).Where("parent_id IS NULL AND status = ?", models.CategoryStatusActive)
	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	var rootCategories []*models.Category
	if err := query.Order("sort_order ASC, name ASC").Find(&rootCategories).Error; err != nil {
		return nil, fmt.Errorf("failed to get root categories: %w", err)
	}

	// Build tree
	var tree []*CategoryNode
	for _, root := range rootCategories {
		node := &CategoryNode{Category: root}
		if maxDepth > 0 {
			children, err := s.buildCategoryTree(ctx, root.ID, maxDepth-1)
			if err != nil {
				return nil, err
			}
			node.Children = children
		}
		tree = append(tree, node)
	}

	return tree, nil
}

func (s *categoryService) buildCategoryTree(ctx context.Context, parentID uuid.UUID, depth int) ([]*CategoryNode, error) {
	if depth < 0 {
		return nil, nil
	}

	var categories []*models.Category
	if err := s.db.WithContext(ctx).Where("parent_id = ? AND status = ?", parentID, models.CategoryStatusActive).
		Order("sort_order ASC, name ASC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get child categories: %w", err)
	}

	var nodes []*CategoryNode
	for _, category := range categories {
		node := &CategoryNode{Category: category}
		if depth > 0 {
			children, err := s.buildCategoryTree(ctx, category.ID, depth-1)
			if err != nil {
				return nil, err
			}
			node.Children = children
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (s *categoryService) GetCategoryChildren(ctx context.Context, parentID uuid.UUID, recursive bool) ([]*models.Category, error) {
	var categories []*models.Category

	if recursive {
		// Get all descendants
		if err := s.db.WithContext(ctx).Raw(`
			WITH RECURSIVE category_tree AS (
				SELECT * FROM categories WHERE parent_id = ?
				UNION ALL
				SELECT c.* FROM categories c
				INNER JOIN category_tree ct ON c.parent_id = ct.id
			)
			SELECT * FROM category_tree ORDER BY level, sort_order, name
		`, parentID).Find(&categories).Error; err != nil {
			return nil, fmt.Errorf("failed to get category descendants: %w", err)
		}
	} else {
		// Get direct children only
		if err := s.db.WithContext(ctx).Where("parent_id = ?", parentID).
			Order("sort_order ASC, name ASC").Find(&categories).Error; err != nil {
			return nil, fmt.Errorf("failed to get category children: %w", err)
		}
	}

	return categories, nil
}

func (s *categoryService) GetCategoryAncestors(ctx context.Context, categoryID uuid.UUID) ([]*models.Category, error) {
	var ancestors []*models.Category

	if err := s.db.WithContext(ctx).Raw(`
		WITH RECURSIVE category_ancestors AS (
			SELECT * FROM categories WHERE id = ?
			UNION ALL
			SELECT c.* FROM categories c
			INNER JOIN category_ancestors ca ON c.id = ca.parent_id
		)
		SELECT * FROM category_ancestors WHERE id != ? ORDER BY level
	`, categoryID, categoryID).Find(&ancestors).Error; err != nil {
		return nil, fmt.Errorf("failed to get category ancestors: %w", err)
	}

	return ancestors, nil
}

func (s *categoryService) MoveCategoryToParent(ctx context.Context, categoryID uuid.UUID, newParentID *uuid.UUID) error {
	category, err := s.GetCategory(ctx, categoryID)
	if err != nil {
		return err
	}

	// Validate new parent
	if newParentID != nil {
		parent, err := s.GetCategory(ctx, *newParentID)
		if err != nil {
			return fmt.Errorf("invalid parent category: %w", err)
		}

		// Check for circular reference
		ancestors, err := s.GetCategoryAncestors(ctx, parent.ID)
		if err != nil {
			return fmt.Errorf("failed to check ancestors: %w", err)
		}

		for _, ancestor := range ancestors {
			if ancestor.ID == categoryID {
				return fmt.Errorf("cannot move category to its own descendant")
			}
		}
	}

	// Calculate new level and path
	newLevel := 0
	newPath := category.Name

	if newParentID != nil {
		parent, _ := s.GetCategory(ctx, *newParentID)
		newLevel = parent.Level + 1
		newPath = parent.Path + "/" + category.Name
	}

	oldParentID := category.ParentID

	// Update category
	updates := map[string]interface{}{
		"parent_id": newParentID,
		"level":     newLevel,
		"path":      newPath,
	}

	if err := s.db.WithContext(ctx).Model(category).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to move category: %w", err)
	}

	// Update descendant paths
	if err := s.updateDescendantPaths(ctx, categoryID); err != nil {
		fmt.Printf("Failed to update descendant paths: %v\n", err)
	}

	// Update children counts for old and new parents
	if oldParentID != nil {
		s.updateChildrenCount(ctx, *oldParentID)
	}
	if newParentID != nil {
		s.updateChildrenCount(ctx, *newParentID)
	}

	return nil
}

func (s *categoryService) GetCategoriesByBusiness(ctx context.Context, businessID uuid.UUID, pagination models.Pagination) (*models.CategoryResponse, error) {
	filters := map[string]interface{}{
		"businessId": businessID,
		"status":     models.CategoryStatusActive,
	}
	return s.ListCategories(ctx, filters, pagination)
}

func (s *categoryService) GetGlobalCategories(ctx context.Context, pagination models.Pagination) (*models.CategoryResponse, error) {
	filters := map[string]interface{}{
		"status": models.CategoryStatusActive,
		// For global categories, we don't specify businessId which means business_id IS NULL
	}
	return s.ListCategories(ctx, filters, pagination)
}

func (s *categoryService) GetFeaturedCategories(ctx context.Context, businessID *uuid.UUID, limit int) ([]*models.Category, error) {
	query := s.db.WithContext(ctx).Where("is_featured = ? AND status = ?", true, models.CategoryStatusActive)

	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	var categories []*models.Category
	if err := query.Order("sort_order ASC, name ASC").Limit(limit).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get featured categories: %w", err)
	}

	return categories, nil
}

func (s *categoryService) AddProductToCategory(ctx context.Context, productID, categoryID uuid.UUID, isPrimary bool) error {
	// Check if relationship already exists
	var existing models.ProductCategory
	err := s.db.WithContext(ctx).Where("product_id = ? AND category_id = ?", productID, categoryID).First(&existing).Error

	if err == nil {
		// Update existing relationship
		if err := s.db.WithContext(ctx).Model(&existing).Update("is_primary", isPrimary).Error; err != nil {
			return fmt.Errorf("failed to update product category: %w", err)
		}
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing relationship: %w", err)
	}

	// If setting as primary, unset other primary categories for this product
	if isPrimary {
		if err := s.db.WithContext(ctx).Model(&models.ProductCategory{}).
			Where("product_id = ?", productID).Update("is_primary", false).Error; err != nil {
			return fmt.Errorf("failed to unset other primary categories: %w", err)
		}
	}

	// Create new relationship
	productCategory := &models.ProductCategory{
		ProductID:  productID,
		CategoryID: categoryID,
		IsPrimary:  isPrimary,
	}

	if err := s.db.WithContext(ctx).Create(productCategory).Error; err != nil {
		return fmt.Errorf("failed to add product to category: %w", err)
	}

	// Update category product count
	if err := s.UpdateCategoryStats(ctx, categoryID); err != nil {
		fmt.Printf("Failed to update category stats: %v\n", err)
	}

	return nil
}

func (s *categoryService) RemoveProductFromCategory(ctx context.Context, productID, categoryID uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("product_id = ? AND category_id = ?", productID, categoryID).
		Delete(&models.ProductCategory{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove product from category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product category relationship not found")
	}

	// Update category product count
	if err := s.UpdateCategoryStats(ctx, categoryID); err != nil {
		fmt.Printf("Failed to update category stats: %v\n", err)
	}

	return nil
}

func (s *categoryService) GetProductCategories(ctx context.Context, productID uuid.UUID) ([]*models.Category, error) {
	var categories []*models.Category

	if err := s.db.WithContext(ctx).
		Joins("JOIN product_categories ON categories.id = product_categories.category_id").
		Where("product_categories.product_id = ?", productID).
		Order("product_categories.is_primary DESC, categories.name ASC").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get product categories: %w", err)
	}

	return categories, nil
}

func (s *categoryService) GetCategoryProducts(ctx context.Context, categoryID uuid.UUID, pagination models.Pagination) ([]uuid.UUID, int64, error) {
	// Count total
	var total int64
	if err := s.db.WithContext(ctx).Model(&models.ProductCategory{}).
		Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count category products: %w", err)
	}

	// Get product IDs
	var productIDs []uuid.UUID
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := s.db.WithContext(ctx).Model(&models.ProductCategory{}).
		Select("product_id").Where("category_id = ?", categoryID).
		Order("is_primary DESC, sort_order ASC, created_at DESC").
		Offset(offset).Limit(pagination.PageSize).Find(&productIDs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get category products: %w", err)
	}

	return productIDs, total, nil
}

func (s *categoryService) UpdateCategoryStats(ctx context.Context, categoryID uuid.UUID) error {
	// Count total products
	var productCount int64
	if err := s.db.WithContext(ctx).Model(&models.ProductCategory{}).
		Where("category_id = ?", categoryID).Count(&productCount).Error; err != nil {
		return fmt.Errorf("failed to count products: %w", err)
	}

	// Count active products (assuming we have access to product status)
	var activeProductCount int64
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM product_categories pc
		JOIN products p ON pc.product_id = p.id
		WHERE pc.category_id = ? AND p.status = 'active'
	`, categoryID).Scan(&activeProductCount).Error; err != nil {
		// If products table doesn't exist or accessible, use total count
		activeProductCount = productCount
	}

	// Count children
	var childrenCount int64
	if err := s.db.WithContext(ctx).Model(&models.Category{}).
		Where("parent_id = ?", categoryID).Count(&childrenCount).Error; err != nil {
		return fmt.Errorf("failed to count children: %w", err)
	}

	// Update stats
	updates := map[string]interface{}{
		"product_count":        productCount,
		"active_product_count": activeProductCount,
		"children_count":       childrenCount,
	}

	if err := s.db.WithContext(ctx).Model(&models.Category{}).
		Where("id = ?", categoryID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update category stats: %w", err)
	}

	return nil
}

func (s *categoryService) RecalculateAllStats(ctx context.Context) error {
	var categories []*models.Category
	if err := s.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return fmt.Errorf("failed to get categories: %w", err)
	}

	for _, category := range categories {
		if err := s.UpdateCategoryStats(ctx, category.ID); err != nil {
			fmt.Printf("Failed to update stats for category %s: %v\n", category.Name, err)
		}
	}

	return nil
}

func (s *categoryService) TrackCategoryView(ctx context.Context, categoryID uuid.UUID, userID *uuid.UUID, sessionID, ipAddress, userAgent, referrer string) error {
	view := &models.CategoryView{
		CategoryID: categoryID,
		UserID:     userID,
		SessionID:  sessionID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Referrer:   referrer,
		ViewedAt:   time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(view).Error; err != nil {
		return fmt.Errorf("failed to track category view: %w", err)
	}

	return nil
}

func (s *categoryService) GetCategoryAnalytics(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (*CategoryAnalytics, error) {
	// Get comprehensive analytics from repository
	analyticsData, err := s.categoryViewRepo.GetAnalytics(ctx, categoryID, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics data: %w", err)
	}

	analytics := &CategoryAnalytics{
		TotalViews:     analyticsData["total_views"].(int64),
		UniqueVisitors: analyticsData["unique_visitors"].(int64),
	}

	// Convert top referrers
	if topReferrersData, ok := analyticsData["top_referrers"].([]map[string]interface{}); ok {
		analytics.TopReferrers = make([]ReferrerStat, len(topReferrersData))
		for i, referrer := range topReferrersData {
			analytics.TopReferrers[i] = ReferrerStat{
				Referrer: referrer["referrer"].(string),
				Count:    referrer["count"].(int64),
			}
		}
	}

	// Convert views by day
	if viewsByDayData, ok := analyticsData["views_by_day"].([]map[string]interface{}); ok {
		analytics.ViewsByDay = make([]DayViewStat, len(viewsByDayData))
		for i, day := range viewsByDayData {
			analytics.ViewsByDay[i] = DayViewStat{
				Date:  day["date"].(string),
				Views: day["views"].(int64),
			}
		}
	}

	// Calculate product clicks and conversion rate
	productClicks, err := s.getProductClicksForCategory(ctx, categoryID, dateFrom, dateTo)
	if err != nil {
		fmt.Printf("Failed to get product clicks: %v\n", err)
		productClicks = 0
	}
	analytics.ProductClicks = productClicks

	// Calculate conversion rate (clicks/views)
	if analytics.TotalViews > 0 {
		analytics.ConversionRate = float64(productClicks) / float64(analytics.TotalViews) * 100
	}

	return analytics, nil
}

// Helper methods

func (s *categoryService) getProductClicksForCategory(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (int64, error) {
	// This would typically require tracking product clicks through a separate analytics table
	// For now, we'll estimate based on category products and their view patterns
	query := `
		SELECT COUNT(DISTINCT pc.product_id) *
		       (SELECT COUNT(*) FROM category_views WHERE category_id = ?) as estimated_clicks
		FROM product_categories pc
		WHERE pc.category_id = ?
	`
	args := []interface{}{categoryID, categoryID}

	var clicks int64
	if err := s.db.WithContext(ctx).Raw(query, args...).Scan(&clicks).Error; err != nil {
		return 0, err
	}

	// Apply a conversion factor (estimated 5% click-through rate)
	return clicks / 20, nil
}

func (s *categoryService) calculatePath(ctx context.Context, parentID *uuid.UUID, name string) (string, error) {
	if parentID == nil {
		return name, nil
	}

	parent, err := s.GetCategory(ctx, *parentID)
	if err != nil {
		return "", err
	}

	return parent.Path + "/" + name, nil
}

func (s *categoryService) updateDescendantPaths(ctx context.Context, categoryID uuid.UUID) error {
	// Get the updated category
	category, err := s.GetCategory(ctx, categoryID)
	if err != nil {
		return err
	}

	// Get all descendants
	descendants, err := s.GetCategoryChildren(ctx, categoryID, true)
	if err != nil {
		return err
	}

	// Update each descendant's path
	for _, descendant := range descendants {
		// Calculate new path
		pathParts := strings.Split(descendant.Path, "/")
		if len(pathParts) > category.Level+1 {
			// Reconstruct path with new parent path
			newPathParts := strings.Split(category.Path, "/")
			newPathParts = append(newPathParts, pathParts[category.Level+1:]...)
			newPath := strings.Join(newPathParts, "/")

			if err := s.db.WithContext(ctx).Model(&descendant).Update("path", newPath).Error; err != nil {
				return fmt.Errorf("failed to update descendant path: %w", err)
			}
		}
	}

	return nil
}

func (s *categoryService) updateChildrenCount(ctx context.Context, parentID uuid.UUID) error {
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.Category{}).
		Where("parent_id = ?", parentID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count children: %w", err)
	}

	if err := s.db.WithContext(ctx).Model(&models.Category{}).
		Where("id = ?", parentID).Update("children_count", count).Error; err != nil {
		return fmt.Errorf("failed to update children count: %w", err)
	}

	return nil
}

// GetBusinessCategories gets categories for a specific business (alias for GetCategoriesByBusiness)
func (s *categoryService) GetBusinessCategories(ctx context.Context, businessID uuid.UUID, pagination models.Pagination) (*models.CategoryResponse, error) {
	return s.GetCategoriesByBusiness(ctx, businessID, pagination)
}

// GetRootCategories gets root categories (categories without parents)
func (s *categoryService) GetRootCategories(ctx context.Context, businessID *uuid.UUID) ([]*models.Category, error) {
	return s.categoryRepo.GetRootCategories(ctx, businessID)
}