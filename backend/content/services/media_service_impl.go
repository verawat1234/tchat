package services

import (
	"context"
	"fmt"
	"math"

	"tchat.dev/content/models"
	"tchat.dev/content/repository"
)

// MediaServiceImpl implements the MediaService interface
type MediaServiceImpl struct {
	repo repository.MediaRepository
}

// NewMediaService creates a new MediaService instance
func NewMediaService(repo repository.MediaRepository) MediaService {
	return &MediaServiceImpl{
		repo: repo,
	}
}

// Category operations

func (s *MediaServiceImpl) GetCategories(ctx context.Context) ([]*models.MediaCategory, error) {
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	// Apply business logic: ensure categories are sorted by display order
	// and only active categories are returned (already handled in repository)
	return categories, nil
}

func (s *MediaServiceImpl) GetCategoryByID(ctx context.Context, id string) (*models.MediaCategory, error) {
	if id == "" {
		return nil, fmt.Errorf("category ID is required")
	}

	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// Subtab operations

func (s *MediaServiceImpl) GetMovieSubtabs(ctx context.Context) (*MovieSubtabsResponse, error) {
	// Get subtabs for movies category
	subtabs, err := s.repo.GetSubtabsByCategory(ctx, "movies")
	if err != nil {
		return nil, fmt.Errorf("failed to get movie subtabs: %w", err)
	}

	// Business logic: default to "short" if available, otherwise first subtab
	defaultSubtab := "short"
	if len(subtabs) > 0 {
		// Check if "short" subtab exists
		found := false
		for _, subtab := range subtabs {
			if subtab.ID == "short" {
				found = true
				break
			}
		}
		if !found {
			defaultSubtab = subtabs[0].ID
		}
	}

	return &MovieSubtabsResponse{
		Subtabs:       subtabs,
		DefaultSubtab: defaultSubtab,
	}, nil
}

// Content operations

func (s *MediaServiceImpl) GetContentByCategory(ctx context.Context, req *GetContentByCategoryRequest) (*PaginatedContentResponse, error) {
	// Validate request
	if req.CategoryID == "" {
		return nil, fmt.Errorf("category ID is required")
	}

	// Set default pagination values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit for performance
	}

	// Verify category exists
	_, err := s.repo.GetCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Get content with pagination
	content, total, err := s.repo.GetContentByCategory(ctx, req.CategoryID, req.Page, req.Limit, req.Subtab)
	if err != nil {
		return nil, fmt.Errorf("failed to get content for category: %w", err)
	}

	// Calculate if there are more pages
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))
	hasMore := req.Page < totalPages

	return &PaginatedContentResponse{
		Items:   content,
		Page:    req.Page,
		Limit:   req.Limit,
		Total:   total,
		HasMore: hasMore,
	}, nil
}

func (s *MediaServiceImpl) GetFeaturedContent(ctx context.Context, req *GetFeaturedContentRequest) (*FeaturedContentResponse, error) {
	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 50 {
		req.Limit = 50 // Max limit for featured content
	}

	// Validate category if specified
	if req.CategoryID != "" {
		_, err := s.repo.GetCategoryByID(ctx, req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
	}

	// Get featured content
	content, err := s.repo.GetFeaturedContent(ctx, req.CategoryID, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get featured content: %w", err)
	}

	// Business logic: featured content is always limited and sorted by featured_order
	// HasMore is always false for featured content (it's a curated list)
	return &FeaturedContentResponse{
		Items:   content,
		Total:   len(content),
		HasMore: false,
	}, nil
}

func (s *MediaServiceImpl) SearchContent(ctx context.Context, req *SearchContentRequest) (*SearchContentResponse, error) {
	// Validate request
	if req.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	// Set default pagination values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit for performance
	}

	// Validate category if specified
	if req.CategoryID != "" {
		_, err := s.repo.GetCategoryByID(ctx, req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
	}

	// Perform search
	content, total, err := s.repo.SearchContent(ctx, req.Query, req.CategoryID, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search content: %w", err)
	}

	return &SearchContentResponse{
		Items: content,
		Query: req.Query,
		Total: total,
		Page:  req.Page,
	}, nil
}

// Admin operations

func (s *MediaServiceImpl) CreateCategory(ctx context.Context, category *models.MediaCategory) error {
	// Business validation
	if category.ID == "" {
		return fmt.Errorf("category ID is required")
	}
	if category.Name == "" {
		return fmt.Errorf("category name is required")
	}
	if category.DisplayOrder <= 0 {
		return fmt.Errorf("display order must be positive")
	}

	// Check if category already exists
	existing, err := s.repo.GetCategoryByID(ctx, category.ID)
	if err == nil && existing != nil {
		return fmt.Errorf("category with ID %s already exists", category.ID)
	}

	err = s.repo.CreateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

func (s *MediaServiceImpl) CreateContent(ctx context.Context, content *models.MediaContentItem) error {
	// Business validation
	if content.CategoryID == "" {
		return fmt.Errorf("category ID is required")
	}
	if content.Title == "" {
		return fmt.Errorf("content title is required")
	}
	if content.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}

	// Verify category exists
	_, err := s.repo.GetCategoryByID(ctx, content.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category: %w", err)
	}

	err = s.repo.CreateContent(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}

	return nil
}