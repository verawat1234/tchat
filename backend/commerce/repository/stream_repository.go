package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"tchat.dev/commerce/models"
)

// StreamCategoryRepository handles database operations for stream categories
type StreamCategoryRepository struct {
	db *gorm.DB
}

// NewStreamCategoryRepository creates a new stream category repository
func NewStreamCategoryRepository(db *gorm.DB) *StreamCategoryRepository {
	return &StreamCategoryRepository{
		db: db,
	}
}

// GetAllCategories retrieves all active stream categories ordered by display order
func (r *StreamCategoryRepository) GetAllCategories() ([]models.StreamCategory, error) {
	var categories []models.StreamCategory
	err := r.db.Where("is_active = ?", true).
		Order("display_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// GetCategoryByID retrieves a stream category by ID
func (r *StreamCategoryRepository) GetCategoryByID(id string) (*models.StreamCategory, error) {
	var category models.StreamCategory
	err := r.db.Where("id = ? AND is_active = ?", id, true).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategoriesWithSubtabs retrieves categories with their associated subtabs
func (r *StreamCategoryRepository) GetCategoriesWithSubtabs() ([]models.StreamCategory, error) {
	var categories []models.StreamCategory
	err := r.db.Preload("Subtabs", "is_active = ?", true).
		Where("is_active = ?", true).
		Order("display_order ASC").
		Find(&categories).Error
	return categories, err
}

// GetSubtabsByCategoryID retrieves all subtabs for a specific category
func (r *StreamCategoryRepository) GetSubtabsByCategoryID(categoryID string) ([]models.StreamSubtab, error) {
	var subtabs []models.StreamSubtab
	err := r.db.Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("display_order ASC, name ASC").
		Find(&subtabs).Error
	return subtabs, err
}

// GetSubtabByID retrieves a subtab by ID
func (r *StreamCategoryRepository) GetSubtabByID(id string) (*models.StreamSubtab, error) {
	var subtab models.StreamSubtab
	err := r.db.Where("id = ? AND is_active = ?", id, true).
		First(&subtab).Error
	if err != nil {
		return nil, err
	}
	return &subtab, nil
}

// CreateCategory creates a new stream category
func (r *StreamCategoryRepository) CreateCategory(category *models.StreamCategory) error {
	return r.db.Create(category).Error
}

// UpdateCategory updates an existing stream category
func (r *StreamCategoryRepository) UpdateCategory(category *models.StreamCategory) error {
	return r.db.Save(category).Error
}

// DeleteCategory soft deletes a stream category
func (r *StreamCategoryRepository) DeleteCategory(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.StreamCategory{}).Error
}

// CreateSubtab creates a new stream subtab
func (r *StreamCategoryRepository) CreateSubtab(subtab *models.StreamSubtab) error {
	return r.db.Create(subtab).Error
}

// UpdateSubtab updates an existing stream subtab
func (r *StreamCategoryRepository) UpdateSubtab(subtab *models.StreamSubtab) error {
	return r.db.Save(subtab).Error
}

// DeleteSubtab soft deletes a stream subtab
func (r *StreamCategoryRepository) DeleteSubtab(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.StreamSubtab{}).Error
}

// StreamContentRepository handles database operations for stream content
type StreamContentRepository struct {
	db *gorm.DB
}

// NewStreamContentRepository creates a new stream content repository
func NewStreamContentRepository(db *gorm.DB) *StreamContentRepository {
	return &StreamContentRepository{
		db: db,
	}
}

// GetContentByCategory retrieves content items for a specific category with pagination
func (r *StreamContentRepository) GetContentByCategory(categoryID string, limit, offset int, subtabID *string) ([]models.StreamContentItem, int64, error) {
	query := r.db.Model(&models.StreamContentItem{}).
		Where("category_id = ? AND availability_status = ?", categoryID, "available")

	// Filter by subtab if specified
	if subtabID != nil && *subtabID != "" {
		query = query.Where("subtab_id = ?", *subtabID)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	var content []models.StreamContentItem
	err := query.Order("featured_order ASC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&content).Error

	return content, total, err
}

// GetFeaturedContent retrieves featured content for a category
func (r *StreamContentRepository) GetFeaturedContent(categoryID string, limit int) ([]models.StreamContentItem, error) {
	var content []models.StreamContentItem
	err := r.db.Where("category_id = ? AND is_featured = ? AND availability_status = ?",
		categoryID, true, "available").
		Order("featured_order ASC, created_at DESC").
		Limit(limit).
		Find(&content).Error
	return content, err
}

// GetContentByID retrieves a specific content item by ID
func (r *StreamContentRepository) GetContentByID(id string) (*models.StreamContentItem, error) {
	var content models.StreamContentItem
	err := r.db.Where("id = ? AND availability_status = ?", id, "available").
		First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// SearchContent searches for content by title, description, or tags
func (r *StreamContentRepository) SearchContent(query string, categoryID *string, limit, offset int) ([]models.StreamContentItem, int64, error) {
	dbQuery := r.db.Model(&models.StreamContentItem{}).
		Where("availability_status = ?", "available").
		Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")

	if categoryID != nil && *categoryID != "" {
		dbQuery = dbQuery.Where("category_id = ?", *categoryID)
	}

	// Get total count
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	var content []models.StreamContentItem
	err := dbQuery.Order("featured_order ASC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&content).Error

	return content, total, err
}

// GetContentByIDs retrieves multiple content items by their IDs
func (r *StreamContentRepository) GetContentByIDs(ids []string) ([]models.StreamContentItem, error) {
	var content []models.StreamContentItem
	err := r.db.Where("id IN ? AND availability_status = ?", ids, "available").
		Find(&content).Error
	return content, err
}

// CreateContent creates a new stream content item
func (r *StreamContentRepository) CreateContent(content *models.StreamContentItem) error {
	return r.db.Create(content).Error
}

// UpdateContent updates an existing stream content item
func (r *StreamContentRepository) UpdateContent(content *models.StreamContentItem) error {
	return r.db.Save(content).Error
}

// DeleteContent soft deletes a stream content item
func (r *StreamContentRepository) DeleteContent(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.StreamContentItem{}).Error
}

// GetPopularContent retrieves popular content based on view counts
func (r *StreamContentRepository) GetPopularContent(categoryID *string, limit int, timeRange string) ([]models.StreamContentItem, error) {
	query := r.db.Model(&models.StreamContentItem{}).
		Where("availability_status = ?", "available")

	if categoryID != nil && *categoryID != "" {
		query = query.Where("category_id = ?", *categoryID)
	}

	// Add time range filter if specified
	if timeRange != "" {
		var since time.Time
		switch timeRange {
		case "day":
			since = time.Now().AddDate(0, 0, -1)
		case "week":
			since = time.Now().AddDate(0, 0, -7)
		case "month":
			since = time.Now().AddDate(0, -1, 0)
		default:
			since = time.Now().AddDate(0, 0, -7) // Default to week
		}
		query = query.Where("created_at >= ?", since)
	}

	var content []models.StreamContentItem
	err := query.Order("view_count DESC, rating DESC, created_at DESC").
		Limit(limit).
		Find(&content).Error

	return content, err
}

// StreamSessionRepository handles database operations for user sessions and navigation state
type StreamSessionRepository struct {
	db *gorm.DB
}

// NewStreamSessionRepository creates a new stream session repository
func NewStreamSessionRepository(db *gorm.DB) *StreamSessionRepository {
	return &StreamSessionRepository{
		db: db,
	}
}

// GetUserNavigationState retrieves the current navigation state for a user
func (r *StreamSessionRepository) GetUserNavigationState(userID string) (*models.TabNavigationState, error) {
	return models.GetUserNavigationState(r.db, userID)
}

// GetOrCreateUserNavigationState gets existing or creates new navigation state
func (r *StreamSessionRepository) GetOrCreateUserNavigationState(userID, sessionID string) (*models.TabNavigationState, error) {
	return models.GetOrCreateUserNavigationState(r.db, userID, sessionID)
}

// UpdateNavigationState updates user's navigation state
func (r *StreamSessionRepository) UpdateNavigationState(state *models.TabNavigationState) error {
	state.LastVisitedAt = time.Now()
	return r.db.Save(state).Error
}

// CreateUserSession creates a new user session
func (r *StreamSessionRepository) CreateUserSession(session *models.StreamUserSession) error {
	return r.db.Create(session).Error
}

// GetActiveUserSession retrieves active session for a user
func (r *StreamSessionRepository) GetActiveUserSession(userID string) (*models.StreamUserSession, error) {
	var session models.StreamUserSession
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
		Order("last_activity_at DESC").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateUserSessionActivity updates session's last activity timestamp
func (r *StreamSessionRepository) UpdateUserSessionActivity(sessionToken string) error {
	return r.db.Model(&models.StreamUserSession{}).
		Where("session_token = ?", sessionToken).
		Update("last_activity_at", time.Now()).Error
}

// DeactivateUserSession deactivates a user session
func (r *StreamSessionRepository) DeactivateUserSession(sessionToken string) error {
	return r.db.Model(&models.StreamUserSession{}).
		Where("session_token = ?", sessionToken).
		Update("is_active", false).Error
}

// CreateContentView creates a new content view record
func (r *StreamSessionRepository) CreateContentView(view *models.StreamContentView) error {
	return r.db.Create(view).Error
}

// UpdateContentView updates a content view record
func (r *StreamSessionRepository) UpdateContentView(view *models.StreamContentView) error {
	return r.db.Save(view).Error
}

// GetContentView retrieves a content view by user and content ID
func (r *StreamSessionRepository) GetContentView(userID, contentID, sessionID string) (*models.StreamContentView, error) {
	var view models.StreamContentView
	err := r.db.Where("user_id = ? AND content_id = ? AND session_id = ?", userID, contentID, sessionID).
		Order("created_at DESC").
		First(&view).Error
	if err != nil {
		return nil, err
	}
	return &view, nil
}

// GetUserPreferences retrieves user preferences
func (r *StreamSessionRepository) GetUserPreferences(userID string) (*models.StreamUserPreference, error) {
	return models.GetUserPreferences(r.db, userID)
}

// UpdateUserPreferences updates user preferences
func (r *StreamSessionRepository) UpdateUserPreferences(pref *models.StreamUserPreference) error {
	return r.db.Save(pref).Error
}

// GetUserViewHistory retrieves user's content view history with pagination
func (r *StreamSessionRepository) GetUserViewHistory(userID string, limit, offset int) ([]models.StreamContentView, int64, error) {
	// Get total count
	var total int64
	if err := r.db.Model(&models.StreamContentView{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with content preloaded
	var views []models.StreamContentView
	err := r.db.Preload("Content").
		Where("user_id = ?", userID).
		Order("view_started_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&views).Error

	return views, total, err
}

// StreamRepository combines all stream-related repositories
type StreamRepository struct {
	Category *StreamCategoryRepository
	Content  *StreamContentRepository
	Session  *StreamSessionRepository
	db       *gorm.DB
}

// NewStreamRepository creates a new combined stream repository
func NewStreamRepository(db *gorm.DB) *StreamRepository {
	return &StreamRepository{
		Category: NewStreamCategoryRepository(db),
		Content:  NewStreamContentRepository(db),
		Session:  NewStreamSessionRepository(db),
		db:       db,
	}
}

// BeginTransaction starts a new database transaction
func (r *StreamRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

// CommitTransaction commits a database transaction
func (r *StreamRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

// RollbackTransaction rolls back a database transaction
func (r *StreamRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

// GetCategoryStats retrieves statistics for a category (content count, views, etc.)
func (r *StreamRepository) GetCategoryStats(categoryID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count total content items
	var contentCount int64
	if err := r.db.Model(&models.StreamContentItem{}).
		Where("category_id = ? AND availability_status = ?", categoryID, "available").
		Count(&contentCount).Error; err != nil {
		return nil, err
	}
	stats["content_count"] = contentCount

	// Count featured content
	var featuredCount int64
	if err := r.db.Model(&models.StreamContentItem{}).
		Where("category_id = ? AND is_featured = ? AND availability_status = ?",
			categoryID, true, "available").
		Count(&featuredCount).Error; err != nil {
		return nil, err
	}
	stats["featured_count"] = featuredCount

	// Count total views for this category
	var totalViews int64
	if err := r.db.Table("stream_content_views").
		Joins("JOIN stream_content_items ON stream_content_views.content_id = stream_content_items.id").
		Where("stream_content_items.category_id = ?", categoryID).
		Count(&totalViews).Error; err != nil {
		return nil, err
	}
	stats["total_views"] = totalViews

	// Get average rating for category content
	var avgRating float64
	if err := r.db.Model(&models.StreamContentItem{}).
		Where("category_id = ? AND availability_status = ?", categoryID, "available").
		Select("AVG(rating)").Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	stats["average_rating"] = avgRating

	return stats, nil
}

// BulkUpdateContentOrder updates display order for multiple content items
func (r *StreamRepository) BulkUpdateContentOrder(updates []struct {
	ID           string
	FeaturedOrder *int
}) error {
	tx := r.BeginTransaction()
	defer func() {
		if rec := recover(); rec != nil {
			r.RollbackTransaction(tx)
		}
	}()

	for _, update := range updates {
		if err := tx.Model(&models.StreamContentItem{}).
			Where("id = ?", update.ID).
			Update("featured_order", update.FeaturedOrder).Error; err != nil {
			r.RollbackTransaction(tx)
			return fmt.Errorf("failed to update content order for ID %s: %w", update.ID, err)
		}
	}

	return r.CommitTransaction(tx)
}