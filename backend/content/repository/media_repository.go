package repository

import (
	"context"

	"tchat.dev/content/models"
)

// MediaRepository defines the interface for media data access
type MediaRepository interface {
	// Category operations
	GetCategories(ctx context.Context) ([]*models.MediaCategory, error)
	GetCategoryByID(ctx context.Context, id string) (*models.MediaCategory, error)
	CreateCategory(ctx context.Context, category *models.MediaCategory) error
	UpdateCategory(ctx context.Context, category *models.MediaCategory) error
	DeleteCategory(ctx context.Context, id string) error

	// Subtab operations
	GetSubtabsByCategory(ctx context.Context, categoryID string) ([]*models.MediaSubtab, error)
	GetSubtabByID(ctx context.Context, id string) (*models.MediaSubtab, error)
	CreateSubtab(ctx context.Context, subtab *models.MediaSubtab) error
	UpdateSubtab(ctx context.Context, subtab *models.MediaSubtab) error
	DeleteSubtab(ctx context.Context, id string) error

	// Content operations
	GetContentByCategory(ctx context.Context, categoryID string, page, limit int, subtab string) ([]*models.MediaContentItem, int64, error)
	GetContentByID(ctx context.Context, id string) (*models.MediaContentItem, error)
	GetFeaturedContent(ctx context.Context, categoryID string, limit int) ([]*models.MediaContentItem, error)
	SearchContent(ctx context.Context, query, categoryID string, page, limit int) ([]*models.MediaContentItem, int64, error)
	CreateContent(ctx context.Context, content *models.MediaContentItem) error
	UpdateContent(ctx context.Context, content *models.MediaContentItem) error
	DeleteContent(ctx context.Context, id string) error

	// Utility operations
	GetTotalContentCount(ctx context.Context, categoryID string, subtab string) (int64, error)
	GetTotalSearchCount(ctx context.Context, query, categoryID string) (int64, error)
}