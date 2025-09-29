package repository

import (
	"context"
	"fmt"
	"strings"

	"tchat.dev/content/models"
	"gorm.io/gorm"
)

// MediaRepositoryImpl implements the MediaRepository interface
type MediaRepositoryImpl struct {
	db *gorm.DB
}

// NewMediaRepository creates a new MediaRepository instance
func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &MediaRepositoryImpl{
		db: db,
	}
}

// Category operations

func (r *MediaRepositoryImpl) GetCategories(ctx context.Context) ([]*models.MediaCategory, error) {
	var categories []*models.MediaCategory
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("display_order ASC").
		Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

func (r *MediaRepositoryImpl) GetCategoryByID(ctx context.Context, id string) (*models.MediaCategory, error) {
	var category models.MediaCategory
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", id, true).
		First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

func (r *MediaRepositoryImpl) CreateCategory(ctx context.Context, category *models.MediaCategory) error {
	err := r.db.WithContext(ctx).Create(category).Error
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *MediaRepositoryImpl) UpdateCategory(ctx context.Context, category *models.MediaCategory) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", category.ID).
		Updates(category).Error
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *MediaRepositoryImpl) DeleteCategory(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).
		Model(&models.MediaCategory{}).
		Where("id = ?", id).
		Update("is_active", false).Error
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// Subtab operations

func (r *MediaRepositoryImpl) GetSubtabsByCategory(ctx context.Context, categoryID string) ([]*models.MediaSubtab, error) {
	var subtabs []*models.MediaSubtab
	err := r.db.WithContext(ctx).
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("display_order ASC").
		Find(&subtabs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get subtabs for category %s: %w", categoryID, err)
	}
	return subtabs, nil
}

func (r *MediaRepositoryImpl) GetSubtabByID(ctx context.Context, id string) (*models.MediaSubtab, error) {
	var subtab models.MediaSubtab
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", id, true).
		First(&subtab).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subtab not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get subtab: %w", err)
	}
	return &subtab, nil
}

func (r *MediaRepositoryImpl) CreateSubtab(ctx context.Context, subtab *models.MediaSubtab) error {
	err := r.db.WithContext(ctx).Create(subtab).Error
	if err != nil {
		return fmt.Errorf("failed to create subtab: %w", err)
	}
	return nil
}

func (r *MediaRepositoryImpl) UpdateSubtab(ctx context.Context, subtab *models.MediaSubtab) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", subtab.ID).
		Updates(subtab).Error
	if err != nil {
		return fmt.Errorf("failed to update subtab: %w", err)
	}
	return nil
}

func (r *MediaRepositoryImpl) DeleteSubtab(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).
		Model(&models.MediaSubtab{}).
		Where("id = ?", id).
		Update("is_active", false).Error
	if err != nil {
		return fmt.Errorf("failed to delete subtab: %w", err)
	}
	return nil
}

// Content operations

func (r *MediaRepositoryImpl) GetContentByCategory(ctx context.Context, categoryID string, page, limit int, subtab string) ([]*models.MediaContentItem, int64, error) {
	var content []*models.MediaContentItem
	var total int64

	query := r.db.WithContext(ctx).
		Where("category_id = ? AND availability_status = ?", categoryID, "available")

	// Apply subtab filtering if specified
	if subtab != "" {
		switch subtab {
		case "short":
			// Short films: duration < 30 minutes (1800 seconds)
			query = query.Where("duration < ?", 1800)
		case "long":
			// Feature films: duration >= 30 minutes
			query = query.Where("duration >= ?", 1800)
		}
	}

	// Get total count
	err := query.Model(&models.MediaContentItem{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count content: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	err = query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&content).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get content for category %s: %w", categoryID, err)
	}

	return content, total, nil
}

func (r *MediaRepositoryImpl) GetContentByID(ctx context.Context, id string) (*models.MediaContentItem, error) {
	var content models.MediaContentItem
	err := r.db.WithContext(ctx).
		Where("id = ? AND availability_status = ?", id, "available").
		First(&content).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("content not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get content: %w", err)
	}
	return &content, nil
}

func (r *MediaRepositoryImpl) GetFeaturedContent(ctx context.Context, categoryID string, limit int) ([]*models.MediaContentItem, error) {
	var content []*models.MediaContentItem

	query := r.db.WithContext(ctx).
		Where("is_featured = ? AND availability_status = ?", true, "available")

	// Filter by category if specified
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	err := query.
		Order("featured_order ASC, created_at DESC").
		Limit(limit).
		Find(&content).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get featured content: %w", err)
	}

	return content, nil
}

func (r *MediaRepositoryImpl) SearchContent(ctx context.Context, query, categoryID string, page, limit int) ([]*models.MediaContentItem, int64, error) {
	var content []*models.MediaContentItem
	var total int64

	// Build search query
	searchQuery := r.db.WithContext(ctx).
		Where("availability_status = ?", "available")

	// Add text search
	if query != "" {
		searchTerm := "%" + strings.ToLower(query) + "%"
		searchQuery = searchQuery.Where(
			"LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm,
		)
	}

	// Filter by category if specified
	if categoryID != "" {
		searchQuery = searchQuery.Where("category_id = ?", categoryID)
	}

	// Get total count
	err := searchQuery.Model(&models.MediaContentItem{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	err = searchQuery.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&content).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search content: %w", err)
	}

	return content, total, nil
}

func (r *MediaRepositoryImpl) CreateContent(ctx context.Context, content *models.MediaContentItem) error {
	err := r.db.WithContext(ctx).Create(content).Error
	if err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}
	return nil
}

func (r *MediaRepositoryImpl) UpdateContent(ctx context.Context, content *models.MediaContentItem) error {
	err := r.db.WithContext(ctx).
		Where("id = ?", content.ID).
		Updates(content).Error
	if err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}
	return nil
}

func (r *MediaRepositoryImpl) DeleteContent(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).
		Model(&models.MediaContentItem{}).
		Where("id = ?", id).
		Update("availability_status", "unavailable").Error
	if err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}
	return nil
}

// Utility operations

func (r *MediaRepositoryImpl) GetTotalContentCount(ctx context.Context, categoryID string, subtab string) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Model(&models.MediaContentItem{}).
		Where("category_id = ? AND availability_status = ?", categoryID, "available")

	// Apply subtab filtering if specified
	if subtab != "" {
		switch subtab {
		case "short":
			query = query.Where("duration < ?", 1800)
		case "long":
			query = query.Where("duration >= ?", 1800)
		}
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count content: %w", err)
	}

	return count, nil
}

func (r *MediaRepositoryImpl) GetTotalSearchCount(ctx context.Context, query, categoryID string) (int64, error) {
	var count int64

	searchQuery := r.db.WithContext(ctx).
		Model(&models.MediaContentItem{}).
		Where("availability_status = ?", "available")

	// Add text search
	if query != "" {
		searchTerm := "%" + strings.ToLower(query) + "%"
		searchQuery = searchQuery.Where(
			"LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm,
		)
	}

	// Filter by category if specified
	if categoryID != "" {
		searchQuery = searchQuery.Where("category_id = ?", categoryID)
	}

	err := searchQuery.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count search results: %w", err)
	}

	return count, nil
}