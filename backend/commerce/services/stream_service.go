package services

import (
	"fmt"
	"time"

	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
)

// StreamCategoryService handles business logic for stream categories
type StreamCategoryService struct {
	repo *repository.StreamRepository
}

// NewStreamCategoryService creates a new stream category service
func NewStreamCategoryService(repo *repository.StreamRepository) *StreamCategoryService {
	return &StreamCategoryService{
		repo: repo,
	}
}

// GetCategories retrieves all active stream categories
func (s *StreamCategoryService) GetCategories() ([]models.StreamCategory, error) {
	categories, err := s.repo.Category.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Initialize default categories if none exist
	if len(categories) == 0 {
		return s.initializeDefaultCategories()
	}

	return categories, nil
}

// GetCategoriesWithSubtabs retrieves categories with their subtabs
func (s *StreamCategoryService) GetCategoriesWithSubtabs() ([]models.StreamCategory, error) {
	return s.repo.Category.GetCategoriesWithSubtabs()
}

// GetCategoryByID retrieves a specific category by ID
func (s *StreamCategoryService) GetCategoryByID(id string) (*models.StreamCategory, error) {
	return s.repo.Category.GetCategoryByID(id)
}

// GetSubtabsByCategoryID retrieves subtabs for a category
func (s *StreamCategoryService) GetSubtabsByCategoryID(categoryID string) ([]models.StreamSubtab, error) {
	return s.repo.Category.GetSubtabsByCategoryID(categoryID)
}

// GetCategoryStats retrieves statistics for a category
func (s *StreamCategoryService) GetCategoryStats(categoryID string) (map[string]interface{}, error) {
	return s.repo.GetCategoryStats(categoryID)
}

// initializeDefaultCategories creates the default 6 categories
func (s *StreamCategoryService) initializeDefaultCategories() ([]models.StreamCategory, error) {
	defaultCategories := []models.StreamCategory{
		{
			ID:                     "books",
			Name:                   "Books",
			DisplayOrder:           1,
			IconName:               "book-open",
			IsActive:               true,
			FeaturedContentEnabled: true,
		},
		{
			ID:                     "podcasts",
			Name:                   "Podcasts",
			DisplayOrder:           2,
			IconName:               "microphone",
			IsActive:               true,
			FeaturedContentEnabled: true,
		},
		{
			ID:                     "cartoons",
			Name:                   "Cartoons",
			DisplayOrder:           3,
			IconName:               "film",
			IsActive:               true,
			FeaturedContentEnabled: true,
		},
		{
			ID:                     "movies",
			Name:                   "Movies",
			DisplayOrder:           4,
			IconName:               "video",
			IsActive:               true,
			FeaturedContentEnabled: true,
		},
		{
			ID:                     "music",
			Name:                   "Music",
			DisplayOrder:           5,
			IconName:               "music",
			IsActive:               true,
			FeaturedContentEnabled: true,
		},
		{
			ID:                     "art",
			Name:                   "Art",
			DisplayOrder:           6,
			IconName:               "palette",
			IsActive:               true,
			FeaturedContentEnabled: true,
		},
	}

	// Create categories in database
	for i := range defaultCategories {
		if err := s.repo.Category.CreateCategory(&defaultCategories[i]); err != nil {
			return nil, fmt.Errorf("failed to create default category %s: %w", defaultCategories[i].ID, err)
		}
	}

	// Initialize default subtabs for movies category
	movieSubtabs := []models.StreamSubtab{
		{
			ID:           "short-movies",
			CategoryID:   "movies",
			Name:         "Short Films",
			DisplayOrder: 1,
			IsActive:     true,
		},
		{
			ID:           "long-movies",
			CategoryID:   "movies",
			Name:         "Feature Films",
			DisplayOrder: 2,
			IsActive:     true,
		},
	}

	// Set filter criteria for subtabs
	movieSubtabs[0].SetFilterCriteria(map[string]interface{}{"maxDuration": 1800})
	movieSubtabs[1].SetFilterCriteria(map[string]interface{}{"minDuration": 1801})

	for i := range movieSubtabs {
		if err := s.repo.Category.CreateSubtab(&movieSubtabs[i]); err != nil {
			return nil, fmt.Errorf("failed to create movie subtab %s: %w", movieSubtabs[i].ID, err)
		}
	}

	return defaultCategories, nil
}

// StreamContentService handles business logic for stream content
type StreamContentService struct {
	repo *repository.StreamRepository
}

// NewStreamContentService creates a new stream content service
func NewStreamContentService(repo *repository.StreamRepository) *StreamContentService {
	return &StreamContentService{
		repo: repo,
	}
}

// GetContentByCategory retrieves content for a specific category
func (s *StreamContentService) GetContentByCategory(categoryID string, page, limit int, subtabID *string) (*models.ContentResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Default limit
	}

	offset := (page - 1) * limit

	items, total, err := s.repo.Content.GetContentByCategory(categoryID, limit, offset, subtabID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content for category %s: %w", categoryID, err)
	}

	hasMore := int64(offset+limit) < total

	return &models.ContentResponse{
		Items:   items,
		Total:   total,
		HasMore: hasMore,
	}, nil
}

// GetFeaturedContent retrieves featured content for a category
func (s *StreamContentService) GetFeaturedContent(categoryID string, limit int) (*models.FeaturedResponse, error) {
	if limit < 1 || limit > 50 {
		limit = 10 // Default limit for featured content
	}

	items, err := s.repo.Content.GetFeaturedContent(categoryID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured content for category %s: %w", categoryID, err)
	}

	return &models.FeaturedResponse{
		Items:   items,
		Total:   len(items),
		HasMore: false, // Featured content is typically limited
	}, nil
}

// GetContentByID retrieves a specific content item
func (s *StreamContentService) GetContentByID(id string) (*models.StreamContentItem, error) {
	content, err := s.repo.Content.GetContentByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get content %s: %w", id, err)
	}

	// Note: ViewCount tracking would be implemented here
	// The current model doesn't have ViewCount field, so we skip this for now
	// This can be added to the model if needed

	return content, nil
}

// SearchContent searches for content
func (s *StreamContentService) SearchContent(query string, categoryID *string, page, limit int) (*models.ContentResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	items, total, err := s.repo.Content.SearchContent(query, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search content: %w", err)
	}

	hasMore := int64(offset+limit) < total

	return &models.ContentResponse{
		Items:   items,
		Total:   total,
		HasMore: hasMore,
	}, nil
}

// GetPopularContent retrieves popular content
func (s *StreamContentService) GetPopularContent(categoryID *string, limit int, timeRange string) ([]models.StreamContentItem, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.repo.Content.GetPopularContent(categoryID, limit, timeRange)
}

// StreamSessionService handles business logic for user sessions and navigation
type StreamSessionService struct {
	repo *repository.StreamRepository
}

// NewStreamSessionService creates a new stream session service
func NewStreamSessionService(repo *repository.StreamRepository) *StreamSessionService {
	return &StreamSessionService{
		repo: repo,
	}
}

// GetUserNavigationState retrieves user's current navigation state
func (s *StreamSessionService) GetUserNavigationState(userID string) (*models.TabNavigationState, error) {
	return s.repo.Session.GetUserNavigationState(userID)
}

// UpdateUserNavigationState updates user's navigation state
func (s *StreamSessionService) UpdateUserNavigationState(userID, sessionID, categoryID string, subtabID *string) (*models.TabNavigationState, error) {
	state, err := s.repo.Session.GetOrCreateUserNavigationState(userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get navigation state: %w", err)
	}

	// Update navigation state
	state.CurrentCategoryID = categoryID
	state.CurrentSubtabID = subtabID
	state.SessionID = sessionID

	if err := s.repo.Session.UpdateNavigationState(state); err != nil {
		return nil, fmt.Errorf("failed to update navigation state: %w", err)
	}

	return state, nil
}

// TrackContentView tracks user's content viewing
func (s *StreamSessionService) TrackContentView(userID, contentID, sessionID string) (*models.StreamContentView, error) {
	// Check if view already exists
	existingView, err := s.repo.Session.GetContentView(userID, contentID, sessionID)
	if err == nil {
		// Update existing view
		existingView.UpdateActivity()
		if err := s.repo.Session.UpdateContentView(existingView); err != nil {
			return nil, fmt.Errorf("failed to update content view: %w", err)
		}
		return existingView, nil
	}

	// Create new view
	view := &models.StreamContentView{
		UserID:         userID,
		ContentID:      contentID,
		SessionID:      sessionID,
		ViewStartedAt:  time.Now(),
		DevicePlatform: "web", // Default to web, can be updated based on request headers
	}

	if err := s.repo.Session.CreateContentView(view); err != nil {
		return nil, fmt.Errorf("failed to create content view: %w", err)
	}

	return view, nil
}

// UpdateContentViewProgress updates the progress of content viewing
func (s *StreamSessionService) UpdateContentViewProgress(userID, contentID, sessionID string, progress float64) error {
	view, err := s.repo.Session.GetContentView(userID, contentID, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get content view: %w", err)
	}

	view.UpdateProgress(progress)

	if err := s.repo.Session.UpdateContentView(view); err != nil {
		return fmt.Errorf("failed to update content view progress: %w", err)
	}

	return nil
}

// GetUserPreferences retrieves user's stream preferences
func (s *StreamSessionService) GetUserPreferences(userID string) (*models.StreamUserPreference, error) {
	return s.repo.Session.GetUserPreferences(userID)
}

// UpdateUserPreferences updates user's stream preferences
func (s *StreamSessionService) UpdateUserPreferences(userID string, prefs *models.StreamUserPreference) error {
	prefs.UserID = userID
	return s.repo.Session.UpdateUserPreferences(prefs)
}

// GetUserViewHistory retrieves user's content viewing history
func (s *StreamSessionService) GetUserViewHistory(userID string, page, limit int) ([]models.StreamContentView, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.repo.Session.GetUserViewHistory(userID, limit, offset)
}

// StreamPurchaseService handles business logic for stream content purchases
type StreamPurchaseService struct {
	repo *repository.StreamRepository
}

// NewStreamPurchaseService creates a new stream purchase service
func NewStreamPurchaseService(repo *repository.StreamRepository) *StreamPurchaseService {
	return &StreamPurchaseService{
		repo: repo,
	}
}

// ProcessContentPurchase processes a content purchase
func (s *StreamPurchaseService) ProcessContentPurchase(req *models.StreamPurchaseRequest) (*models.StreamPurchaseResponse, error) {
	// Validate request
	if err := s.validatePurchaseRequest(req); err != nil {
		return nil, fmt.Errorf("invalid purchase request: %w", err)
	}

	// Get content item
	content, err := s.repo.Content.GetContentByID(req.MediaContentID)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}

	// Check availability
	if content.AvailabilityStatus != "available" {
		return nil, fmt.Errorf("content is not available for purchase")
	}

	// Calculate total amount - convert decimal.Decimal to float64
	priceFloat, _ := content.Price.Float64()
	totalAmount := priceFloat * float64(req.Quantity)

	// Create order ID (in real implementation, this would integrate with payment service)
	orderID := fmt.Sprintf("order_%d_%s", time.Now().Unix(), req.MediaContentID)

	// For now, return success response
	// In real implementation, this would integrate with payment processing
	return &models.StreamPurchaseResponse{
		OrderID:     orderID,
		TotalAmount: totalAmount,
		Currency:    content.Currency,
		Success:     true,
		Message:     "Purchase completed successfully",
	}, nil
}

// validatePurchaseRequest validates the purchase request
func (s *StreamPurchaseService) validatePurchaseRequest(req *models.StreamPurchaseRequest) error {
	if req.MediaContentID == "" {
		return fmt.Errorf("media content ID is required")
	}
	if req.Quantity < 1 {
		return fmt.Errorf("quantity must be at least 1")
	}
	if req.MediaLicense == "" {
		return fmt.Errorf("media license is required")
	}

	// Validate license type
	validLicenses := map[string]bool{
		"personal": true,
		"family":   true,
		"commercial": true,
	}
	if !validLicenses[req.MediaLicense] {
		return fmt.Errorf("invalid media license type")
	}

	// Validate download format if specified
	if req.DownloadFormat != "" {
		validFormats := map[string]bool{
			"PDF":  true,
			"EPUB": true,
			"MP3":  true,
			"MP4":  true,
			"FLAC": true,
		}
		if !validFormats[req.DownloadFormat] {
			return fmt.Errorf("invalid download format")
		}
	}

	return nil
}