package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/content/models"
)

// PostgreSQLContentRepository implements ContentRepository using PostgreSQL
type PostgreSQLContentRepository struct {
	db *gorm.DB
}

// NewPostgreSQLContentRepository creates a new PostgreSQL content repository
func NewPostgreSQLContentRepository(db *gorm.DB) *PostgreSQLContentRepository {
	return &PostgreSQLContentRepository{db: db}
}

// CreateContent creates a new content item
func (r *PostgreSQLContentRepository) CreateContent(ctx context.Context, content *models.ContentItem) error {
	if content.ID == uuid.Nil {
		content.ID = uuid.New()
	}
	if content.Status == "" {
		content.Status = models.ContentStatusDraft
	}
	return r.db.WithContext(ctx).Create(content).Error
}

// GetContentByID retrieves a content item by ID
func (r *PostgreSQLContentRepository) GetContentByID(ctx context.Context, id uuid.UUID) (*models.ContentItem, error) {
	var content models.ContentItem
	err := r.db.WithContext(ctx).First(&content, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// GetContentByKey retrieves a content item by key (category)
func (r *PostgreSQLContentRepository) GetContentByKey(ctx context.Context, key string) (*models.ContentItem, error) {
	var content models.ContentItem
	err := r.db.WithContext(ctx).First(&content, "category = ?", key).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// GetContentItems retrieves content items with filters, pagination, and sorting
func (r *PostgreSQLContentRepository) GetContentItems(ctx context.Context, filters models.ContentFilters, pagination models.Pagination, sort models.SortOptions) (*models.ContentResponse, error) {
	var items []models.ContentItem
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ContentItem{})

	// Apply filters
	query = r.applyContentFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if sort.Field != "" {
		order := "ASC"
		if strings.ToUpper(sort.Order) == "DESC" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", sort.Field, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Offset(offset).Limit(pagination.PageSize)

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	return &models.ContentResponse{
		Items: items,
		Pagination: models.Pagination{
			Page:     pagination.Page,
			PageSize: pagination.PageSize,
			Offset:   offset,
			Total:    total,
		},
		Total: total,
	}, nil
}

// GetContentByCategory retrieves content items by category
func (r *PostgreSQLContentRepository) GetContentByCategory(ctx context.Context, category string, pagination models.Pagination, sort models.SortOptions) (*models.ContentResponse, error) {
	filters := models.ContentFilters{
		Category: &category,
	}
	return r.GetContentItems(ctx, filters, pagination, sort)
}

// UpdateContent updates an existing content item
func (r *PostgreSQLContentRepository) UpdateContent(ctx context.Context, content *models.ContentItem) error {
	return r.db.WithContext(ctx).Save(content).Error
}

// DeleteContent deletes a content item (soft delete by changing status)
func (r *PostgreSQLContentRepository) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.ContentItem{}).
		Where("id = ?", id).
		Update("status", models.ContentStatusDeleted).Error
}

// BulkUpdateContent updates multiple content items
func (r *PostgreSQLContentRepository) BulkUpdateContent(ctx context.Context, ids []uuid.UUID, updates models.UpdateContentRequest) ([]models.ContentItem, error) {
	updateData := make(map[string]interface{})

	if updates.Category != nil {
		updateData["category"] = *updates.Category
	}
	if updates.Type != nil {
		updateData["type"] = *updates.Type
	}
	if updates.Value != nil {
		updateData["value"] = *updates.Value
	}
	if updates.Metadata != nil {
		updateData["metadata"] = *updates.Metadata
	}
	if updates.Status != nil {
		updateData["status"] = *updates.Status
	}
	if updates.Tags != nil {
		updateData["tags"] = *updates.Tags
	}
	if updates.Notes != nil {
		updateData["notes"] = *updates.Notes
	}

	updateData["updated_at"] = time.Now()

	err := r.db.WithContext(ctx).Model(&models.ContentItem{}).
		Where("id IN ?", ids).
		Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	// Fetch updated items
	var items []models.ContentItem
	err = r.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error
	return items, err
}

// PublishContent publishes a content item
func (r *PostgreSQLContentRepository) PublishContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error) {
	err := r.db.WithContext(ctx).Model(&models.ContentItem{}).
		Where("id = ?", id).
		Update("status", models.ContentStatusPublished).Error
	if err != nil {
		return nil, err
	}

	return r.GetContentByID(ctx, id)
}

// ArchiveContent archives a content item
func (r *PostgreSQLContentRepository) ArchiveContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error) {
	err := r.db.WithContext(ctx).Model(&models.ContentItem{}).
		Where("id = ?", id).
		Update("status", models.ContentStatusArchived).Error
	if err != nil {
		return nil, err
	}

	return r.GetContentByID(ctx, id)
}

// SyncContent handles content synchronization
func (r *PostgreSQLContentRepository) SyncContent(ctx context.Context, req models.SyncContentRequest) (*models.SyncContentResponse, error) {
	var items []models.ContentItem
	var deletedIDs []uuid.UUID

	lastSync := time.Now().Add(-24 * time.Hour) // Default to 24 hours ago
	if req.LastSyncTime != nil {
		lastSync = *req.LastSyncTime
	}

	// Get modified content since last sync
	query := r.db.WithContext(ctx).Model(&models.ContentItem{}).
		Where("updated_at > ?", lastSync).
		Where("status != ?", models.ContentStatusDeleted)

	if len(req.Categories) > 0 {
		query = query.Where("category IN ?", req.Categories)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	// Get deleted content IDs since last sync
	deletedIDs, err := r.GetDeletedContentSince(ctx, lastSync)
	if err != nil {
		return nil, err
	}

	return &models.SyncContentResponse{
		Items:      items,
		DeletedIDs: deletedIDs,
		SyncTime:   time.Now(),
		HasMore:    len(items) >= 100, // Simple pagination check
	}, nil
}

// GetContentModifiedSince retrieves content modified since a specific time
func (r *PostgreSQLContentRepository) GetContentModifiedSince(ctx context.Context, since time.Time, categories []string) ([]models.ContentItem, error) {
	var items []models.ContentItem
	query := r.db.WithContext(ctx).Where("updated_at > ?", since)

	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
	}

	err := query.Find(&items).Error
	return items, err
}

// GetDeletedContentSince retrieves deleted content IDs since a specific time
func (r *PostgreSQLContentRepository) GetDeletedContentSince(ctx context.Context, since time.Time) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).Model(&models.ContentItem{}).
		Where("updated_at > ? AND status = ?", since, models.ContentStatusDeleted).
		Pluck("id", &ids).Error
	return ids, err
}

// applyContentFilters applies filters to a GORM query
func (r *PostgreSQLContentRepository) applyContentFilters(query *gorm.DB, filters models.ContentFilters) *gorm.DB {
	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if len(filters.Tags) > 0 {
		query = query.Where("tags && ?", filters.Tags)
	}

	if filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where("category ILIKE ? OR notes ILIKE ? OR value::text ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if filters.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *filters.CreatedFrom)
	}

	if filters.CreatedTo != nil {
		query = query.Where("created_at <= ?", *filters.CreatedTo)
	}

	return query
}

// PostgreSQLCategoryRepository implements CategoryRepository using PostgreSQL
type PostgreSQLCategoryRepository struct {
	db *gorm.DB
}

// NewPostgreSQLCategoryRepository creates a new PostgreSQL category repository
func NewPostgreSQLCategoryRepository(db *gorm.DB) *PostgreSQLCategoryRepository {
	return &PostgreSQLCategoryRepository{db: db}
}

// CreateCategory creates a new content category
func (r *PostgreSQLCategoryRepository) CreateCategory(ctx context.Context, category *models.ContentCategory) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(category).Error
}

// GetCategoryByID retrieves a category by ID
func (r *PostgreSQLCategoryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.ContentCategory, error) {
	var category models.ContentCategory
	err := r.db.WithContext(ctx).Preload("Parent").Preload("Children").First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategories retrieves categories with pagination and sorting
func (r *PostgreSQLCategoryRepository) GetCategories(ctx context.Context, pagination models.Pagination, sort models.SortOptions) (*models.ContentCategoryResponse, error) {
	var categories []models.ContentCategory
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ContentCategory{}).Where("is_active = ?", true)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if sort.Field != "" {
		order := "ASC"
		if strings.ToUpper(sort.Order) == "DESC" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", sort.Field, order))
	} else {
		query = query.Order("sort_order ASC, name ASC")
	}

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Offset(offset).Limit(pagination.PageSize)

	if err := query.Preload("Parent").Find(&categories).Error; err != nil {
		return nil, err
	}

	return &models.ContentCategoryResponse{
		Items: categories,
		Pagination: models.Pagination{
			Page:     pagination.Page,
			PageSize: pagination.PageSize,
			Offset:   offset,
			Total:    total,
		},
		Total: total,
	}, nil
}

// UpdateCategory updates an existing category
func (r *PostgreSQLCategoryRepository) UpdateCategory(ctx context.Context, category *models.ContentCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// DeleteCategory deletes a category (soft delete by setting is_active to false)
func (r *PostgreSQLCategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.ContentCategory{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// GetCategoryChildren retrieves child categories
func (r *PostgreSQLCategoryRepository) GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]models.ContentCategory, error) {
	var categories []models.ContentCategory
	err := r.db.WithContext(ctx).Where("parent_id = ? AND is_active = ?", parentID, true).
		Order("sort_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

// GetCategoryPath retrieves the full path of a category
func (r *PostgreSQLCategoryRepository) GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]models.ContentCategory, error) {
	var path []models.ContentCategory
	currentID := categoryID

	for currentID != uuid.Nil {
		var category models.ContentCategory
		err := r.db.WithContext(ctx).First(&category, "id = ?", currentID).Error
		if err != nil {
			break
		}

		path = append([]models.ContentCategory{category}, path...)
		if category.ParentID == nil {
			break
		}
		currentID = *category.ParentID
	}

	return path, nil
}

// PostgreSQLVersionRepository implements VersionRepository using PostgreSQL
type PostgreSQLVersionRepository struct {
	db *gorm.DB
}

// NewPostgreSQLVersionRepository creates a new PostgreSQL version repository
func NewPostgreSQLVersionRepository(db *gorm.DB) *PostgreSQLVersionRepository {
	return &PostgreSQLVersionRepository{db: db}
}

// CreateVersion creates a new content version
func (r *PostgreSQLVersionRepository) CreateVersion(ctx context.Context, version *models.ContentVersion) error {
	if version.ID == uuid.Nil {
		version.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(version).Error
}

// GetVersionByID retrieves a version by ID
func (r *PostgreSQLVersionRepository) GetVersionByID(ctx context.Context, id uuid.UUID) (*models.ContentVersion, error) {
	var version models.ContentVersion
	err := r.db.WithContext(ctx).Preload("Content").First(&version, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetContentVersions retrieves versions for a content item
func (r *PostgreSQLVersionRepository) GetContentVersions(ctx context.Context, contentID uuid.UUID, pagination models.Pagination) (*models.ContentVersionResponse, error) {
	var versions []models.ContentVersion
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ContentVersion{}).Where("content_id = ?", contentID)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination and sorting
	offset := pagination.Page * pagination.PageSize
	if err := query.Order("version DESC").Offset(offset).Limit(pagination.PageSize).Find(&versions).Error; err != nil {
		return nil, err
	}

	return &models.ContentVersionResponse{
		Items: versions,
		Pagination: models.Pagination{
			Page:     pagination.Page,
			PageSize: pagination.PageSize,
			Offset:   offset,
			Total:    total,
		},
		Total: total,
	}, nil
}

// RevertContentVersion reverts content to a specific version
func (r *PostgreSQLVersionRepository) RevertContentVersion(ctx context.Context, contentID uuid.UUID, version int) (*models.ContentItem, error) {
	// Get the version to revert to
	var contentVersion models.ContentVersion
	err := r.db.WithContext(ctx).Where("content_id = ? AND version = ?", contentID, version).First(&contentVersion).Error
	if err != nil {
		return nil, err
	}

	// Update the current content with the version data
	updateData := map[string]interface{}{
		"value":      contentVersion.Value,
		"metadata":   contentVersion.Metadata,
		"status":     contentVersion.Status,
		"updated_at": time.Now(),
	}

	err = r.db.WithContext(ctx).Model(&models.ContentItem{}).
		Where("id = ?", contentID).
		Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated content
	var content models.ContentItem
	err = r.db.WithContext(ctx).First(&content, "id = ?", contentID).Error
	return &content, err
}

// GetLatestVersion retrieves the latest version for a content item
func (r *PostgreSQLVersionRepository) GetLatestVersion(ctx context.Context, contentID uuid.UUID) (*models.ContentVersion, error) {
	var version models.ContentVersion
	err := r.db.WithContext(ctx).Where("content_id = ?", contentID).
		Order("version DESC").First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// CleanupOldVersions removes old versions, keeping only the specified count
func (r *PostgreSQLVersionRepository) CleanupOldVersions(ctx context.Context, contentID uuid.UUID, keepCount int) error {
	// Get version numbers to keep
	var versionsToKeep []int
	err := r.db.WithContext(ctx).Model(&models.ContentVersion{}).
		Where("content_id = ?", contentID).
		Order("version DESC").
		Limit(keepCount).
		Pluck("version", &versionsToKeep).Error
	if err != nil {
		return err
	}

	if len(versionsToKeep) == 0 {
		return nil // No versions to clean up
	}

	// Delete old versions
	return r.db.WithContext(ctx).Where("content_id = ? AND version NOT IN ?", contentID, versionsToKeep).
		Delete(&models.ContentVersion{}).Error
}