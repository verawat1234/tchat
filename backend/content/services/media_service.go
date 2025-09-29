package services

import (
	"context"

	"tchat.dev/content/models"
)

// MediaService defines the interface for media business logic
type MediaService interface {
	// Category operations
	GetCategories(ctx context.Context) ([]*models.MediaCategory, error)
	GetCategoryByID(ctx context.Context, id string) (*models.MediaCategory, error)

	// Subtab operations
	GetMovieSubtabs(ctx context.Context) (*MovieSubtabsResponse, error)

	// Content operations
	GetContentByCategory(ctx context.Context, req *GetContentByCategoryRequest) (*PaginatedContentResponse, error)
	GetFeaturedContent(ctx context.Context, req *GetFeaturedContentRequest) (*FeaturedContentResponse, error)
	SearchContent(ctx context.Context, req *SearchContentRequest) (*SearchContentResponse, error)

	// Admin operations (future use)
	CreateCategory(ctx context.Context, category *models.MediaCategory) error
	CreateContent(ctx context.Context, content *models.MediaContentItem) error
}

// Request/Response types

type GetContentByCategoryRequest struct {
	CategoryID string
	Page       int
	Limit      int
	Subtab     string
}

type PaginatedContentResponse struct {
	Items   []*models.MediaContentItem `json:"items"`
	Page    int                        `json:"page"`
	Limit   int                        `json:"limit"`
	Total   int64                      `json:"total"`
	HasMore bool                       `json:"hasMore"`
}

type GetFeaturedContentRequest struct {
	CategoryID string
	Limit      int
}

type FeaturedContentResponse struct {
	Items   []*models.MediaContentItem `json:"items"`
	Total   int                        `json:"total"`
	HasMore bool                       `json:"hasMore"`
}

type SearchContentRequest struct {
	Query      string
	CategoryID string
	Page       int
	Limit      int
}

type SearchContentResponse struct {
	Items []*models.MediaContentItem `json:"items"`
	Query string                     `json:"query"`
	Total int64                      `json:"total"`
	Page  int                        `json:"page"`
}

type MovieSubtabsResponse struct {
	Subtabs       []*models.MediaSubtab `json:"subtabs"`
	DefaultSubtab string                `json:"defaultSubtab"`
}