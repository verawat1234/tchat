package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"tchat.dev/content/models"
)

// ContentRepository defines the interface for content data access
type ContentRepository interface {
	// Content CRUD operations
	CreateContent(ctx context.Context, content *models.ContentItem) error
	GetContentByID(ctx context.Context, id uuid.UUID) (*models.ContentItem, error)
	GetContentByKey(ctx context.Context, key string) (*models.ContentItem, error)
	GetContentItems(ctx context.Context, filters models.ContentFilters, pagination models.Pagination, sort models.SortOptions) (*models.ContentResponse, error)
	GetContentByCategory(ctx context.Context, category string, pagination models.Pagination, sort models.SortOptions) (*models.ContentResponse, error)
	UpdateContent(ctx context.Context, content *models.ContentItem) error
	DeleteContent(ctx context.Context, id uuid.UUID) error
	BulkUpdateContent(ctx context.Context, ids []uuid.UUID, updates models.UpdateContentRequest) ([]models.ContentItem, error)

	// Content status operations
	PublishContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error)
	ArchiveContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error)

	// Content sync operations
	SyncContent(ctx context.Context, req models.SyncContentRequest) (*models.SyncContentResponse, error)
	GetContentModifiedSince(ctx context.Context, since time.Time, categories []string) ([]models.ContentItem, error)
	GetDeletedContentSince(ctx context.Context, since time.Time) ([]uuid.UUID, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	// Category CRUD operations
	CreateCategory(ctx context.Context, category *models.ContentCategory) error
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.ContentCategory, error)
	GetCategories(ctx context.Context, pagination models.Pagination, sort models.SortOptions) (*models.ContentCategoryResponse, error)
	UpdateCategory(ctx context.Context, category *models.ContentCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Category hierarchy operations
	GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]models.ContentCategory, error)
	GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]models.ContentCategory, error)
}

// VersionRepository defines the interface for version data access
type VersionRepository interface {
	// Version operations
	CreateVersion(ctx context.Context, version *models.ContentVersion) error
	GetVersionByID(ctx context.Context, id uuid.UUID) (*models.ContentVersion, error)
	GetContentVersions(ctx context.Context, contentID uuid.UUID, pagination models.Pagination) (*models.ContentVersionResponse, error)
	RevertContentVersion(ctx context.Context, contentID uuid.UUID, version int) (*models.ContentItem, error)

	// Version management
	GetLatestVersion(ctx context.Context, contentID uuid.UUID) (*models.ContentVersion, error)
	CleanupOldVersions(ctx context.Context, contentID uuid.UUID, keepCount int) error
}