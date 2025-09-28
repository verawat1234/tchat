package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetByPath(ctx context.Context, path string) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Category, int64, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*models.Category, error)
	GetByBusinessID(ctx context.Context, businessID uuid.UUID, offset, limit int) ([]*models.Category, int64, error)
	GetGlobalCategories(ctx context.Context, offset, limit int) ([]*models.Category, int64, error)
	GetFeaturedCategories(ctx context.Context, businessID *uuid.UUID, limit int) ([]*models.Category, error)
	GetRootCategories(ctx context.Context, businessID *uuid.UUID) ([]*models.Category, error)
}

type ProductCategoryRepository interface {
	Create(ctx context.Context, productCategory *models.ProductCategory) error
	Delete(ctx context.Context, productID, categoryID uuid.UUID) error
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*models.ProductCategory, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID, offset, limit int) ([]*models.ProductCategory, int64, error)
	Update(ctx context.Context, productCategory *models.ProductCategory) error
	SetPrimary(ctx context.Context, productID, categoryID uuid.UUID) error
	UnsetPrimary(ctx context.Context, productID uuid.UUID) error
}

type CategoryViewRepository interface {
	Create(ctx context.Context, view *models.CategoryView) error
	GetAnalytics(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (map[string]interface{}, error)
	CountViews(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (int64, error)
	CountUniqueVisitors(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (int64, error)
}

type categoryRepository struct {
	db *gorm.DB
}

type productCategoryRepository struct {
	db *gorm.DB
}

type categoryViewRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func NewProductCategoryRepository(db *gorm.DB) ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

func NewCategoryViewRepository(db *gorm.DB) CategoryViewRepository {
	return &categoryViewRepository{db: db}
}

// Category Repository Implementation
func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetByPath(ctx context.Context, path string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Where("path = ?", path).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id).Error
}

func (r *categoryRepository) List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Category, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Category{})

	// Apply filters
	for key, value := range filters {
		switch key {
		case "business_id":
			if value != nil {
				query = query.Where("business_id = ?", value)
			} else {
				query = query.Where("business_id IS NULL")
			}
		case "search":
			if value != nil && value.(string) != "" {
				searchTerm := "%" + value.(string) + "%"
				query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
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
	var categories []*models.Category
	err := query.Order("level ASC, sort_order ASC, name ASC").Offset(offset).Limit(limit).Find(&categories).Error
	return categories, total, err
}

func (r *categoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.WithContext(ctx).Where("parent_id = ? AND status = ?", parentID, models.CategoryStatusActive).
		Order("sort_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByBusinessID(ctx context.Context, businessID uuid.UUID, offset, limit int) ([]*models.Category, int64, error) {
	filters := map[string]interface{}{
		"business_id": businessID,
		"status":      models.CategoryStatusActive,
	}
	return r.List(ctx, filters, offset, limit)
}

func (r *categoryRepository) GetGlobalCategories(ctx context.Context, offset, limit int) ([]*models.Category, int64, error) {
	filters := map[string]interface{}{
		"business_id": nil,
		"status":      models.CategoryStatusActive,
	}
	return r.List(ctx, filters, offset, limit)
}

func (r *categoryRepository) GetFeaturedCategories(ctx context.Context, businessID *uuid.UUID, limit int) ([]*models.Category, error) {
	query := r.db.WithContext(ctx).Where("is_featured = ? AND status = ?", true, models.CategoryStatusActive)

	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	var categories []*models.Category
	err := query.Order("sort_order ASC, name ASC").Limit(limit).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetRootCategories(ctx context.Context, businessID *uuid.UUID) ([]*models.Category, error) {
	query := r.db.WithContext(ctx).Where("parent_id IS NULL AND status = ?", models.CategoryStatusActive)

	if businessID != nil {
		query = query.Where("business_id = ?", *businessID)
	} else {
		query = query.Where("business_id IS NULL")
	}

	var categories []*models.Category
	err := query.Order("sort_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

// Product Category Repository Implementation
func (r *productCategoryRepository) Create(ctx context.Context, productCategory *models.ProductCategory) error {
	return r.db.WithContext(ctx).Create(productCategory).Error
}

func (r *productCategoryRepository) Delete(ctx context.Context, productID, categoryID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("product_id = ? AND category_id = ?", productID, categoryID).
		Delete(&models.ProductCategory{}).Error
}

func (r *productCategoryRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*models.ProductCategory, error) {
	var productCategories []*models.ProductCategory
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).
		Order("is_primary DESC, sort_order ASC, created_at ASC").Find(&productCategories).Error
	return productCategories, err
}

func (r *productCategoryRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, offset, limit int) ([]*models.ProductCategory, int64, error) {
	query := r.db.WithContext(ctx).Where("category_id = ?", categoryID)

	// Count total
	var total int64
	if err := query.Model(&models.ProductCategory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var productCategories []*models.ProductCategory
	err := query.Order("is_primary DESC, sort_order ASC, created_at DESC").
		Offset(offset).Limit(limit).Find(&productCategories).Error
	return productCategories, total, err
}

func (r *productCategoryRepository) Update(ctx context.Context, productCategory *models.ProductCategory) error {
	return r.db.WithContext(ctx).Save(productCategory).Error
}

func (r *productCategoryRepository) SetPrimary(ctx context.Context, productID, categoryID uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()

	// Unset all primary categories for the product
	if err := tx.Model(&models.ProductCategory{}).Where("product_id = ?", productID).
		Update("is_primary", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the specified category as primary
	if err := tx.Model(&models.ProductCategory{}).
		Where("product_id = ? AND category_id = ?", productID, categoryID).
		Update("is_primary", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *productCategoryRepository) UnsetPrimary(ctx context.Context, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.ProductCategory{}).
		Where("product_id = ?", productID).Update("is_primary", false).Error
}

// Category View Repository Implementation
func (r *categoryViewRepository) Create(ctx context.Context, view *models.CategoryView) error {
	return r.db.WithContext(ctx).Create(view).Error
}

func (r *categoryViewRepository) GetAnalytics(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Total views
	totalViews, err := r.CountViews(ctx, categoryID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	analytics["total_views"] = totalViews

	// Unique visitors
	uniqueVisitors, err := r.CountUniqueVisitors(ctx, categoryID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	analytics["unique_visitors"] = uniqueVisitors

	// Top referrers
	topReferrers, err := r.getTopReferrers(ctx, categoryID, dateFrom, dateTo, 10)
	if err != nil {
		return nil, err
	}
	analytics["top_referrers"] = topReferrers

	// Views by day
	viewsByDay, err := r.getViewsByDay(ctx, categoryID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	analytics["views_by_day"] = viewsByDay

	// Peak hour analysis
	peakHours, err := r.getPeakHours(ctx, categoryID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	analytics["peak_hours"] = peakHours

	return analytics, nil
}

func (r *categoryViewRepository) CountViews(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (int64, error) {
	query := r.db.WithContext(ctx).Model(&models.CategoryView{}).Where("category_id = ?", categoryID)

	if dateFrom != "" {
		query = query.Where("viewed_at >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("viewed_at <= ?", dateTo)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (r *categoryViewRepository) CountUniqueVisitors(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (int64, error) {
	query := r.db.WithContext(ctx).Model(&models.CategoryView{}).Where("category_id = ?", categoryID)

	if dateFrom != "" {
		query = query.Where("viewed_at >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("viewed_at <= ?", dateTo)
	}

	var count int64
	err := query.Distinct("COALESCE(user_id::text, session_id)").Count(&count).Error
	return count, err
}

func (r *categoryViewRepository) getTopReferrers(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT referrer, COUNT(*) as count
		FROM category_views
		WHERE category_id = ? AND referrer != '' AND referrer IS NOT NULL
	`
	args := []interface{}{categoryID}

	if dateFrom != "" {
		query += " AND viewed_at >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND viewed_at <= ?"
		args = append(args, dateTo)
	}

	query += " GROUP BY referrer ORDER BY count DESC LIMIT ?"
	args = append(args, limit)

	var results []map[string]interface{}
	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var referrer string
		var count int64
		if err := rows.Scan(&referrer, &count); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"referrer": referrer,
			"count":    count,
		})
	}

	return results, nil
}

func (r *categoryViewRepository) getViewsByDay(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) ([]map[string]interface{}, error) {
	query := `
		SELECT DATE(viewed_at) as date, COUNT(*) as views
		FROM category_views
		WHERE category_id = ?
	`
	args := []interface{}{categoryID}

	if dateFrom != "" {
		query += " AND viewed_at >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND viewed_at <= ?"
		args = append(args, dateTo)
	}

	query += " GROUP BY DATE(viewed_at) ORDER BY date"

	var results []map[string]interface{}
	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var views int64
		if err := rows.Scan(&date, &views); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"date":  date,
			"views": views,
		})
	}

	return results, nil
}

func (r *categoryViewRepository) getPeakHours(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) ([]map[string]interface{}, error) {
	query := `
		SELECT EXTRACT(HOUR FROM viewed_at) as hour, COUNT(*) as views
		FROM category_views
		WHERE category_id = ?
	`
	args := []interface{}{categoryID}

	if dateFrom != "" {
		query += " AND viewed_at >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND viewed_at <= ?"
		args = append(args, dateTo)
	}

	query += " GROUP BY EXTRACT(HOUR FROM viewed_at) ORDER BY views DESC"

	var results []map[string]interface{}
	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hour int
		var views int64
		if err := rows.Scan(&hour, &views); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"hour":  hour,
			"views": views,
		})
	}

	return results, nil
}