package interfaces

import (
	"context"

	"github.com/google/uuid"
	"tchat.dev/content/models"
)

// ContentServiceInterface defines the contract for content service operations
// This interface enables dependency injection and mocking for testing
type ContentServiceInterface interface {
	// Content CRUD operations
	CreateContent(ctx context.Context, req models.CreateContentRequest) (*models.ContentItem, error)
	GetContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error)
	GetContentByKey(ctx context.Context, key string) (*models.ContentItem, error)
	GetContentItems(ctx context.Context, filters models.ContentFilters, pagination models.Pagination, sort models.SortOptions) (*models.ContentResponse, error)
	GetContentByCategory(ctx context.Context, category string, pagination models.Pagination, sort models.SortOptions) (*models.ContentResponse, error)
	UpdateContent(ctx context.Context, id uuid.UUID, req models.UpdateContentRequest) (*models.ContentItem, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error

	// Content operations
	PublishContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error)
	ArchiveContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error)
	BulkUpdateContent(ctx context.Context, req models.BulkUpdateRequest) ([]models.ContentItem, error)
	SyncContent(ctx context.Context, req models.SyncContentRequest) (*models.SyncContentResponse, error)

	// Category operations
	GetCategories(ctx context.Context, pagination models.Pagination, sort models.SortOptions) (*models.ContentCategoryResponse, error)
	CreateCategory(ctx context.Context, category *models.ContentCategory) error

	// Version operations
	GetContentVersions(ctx context.Context, contentID uuid.UUID, pagination models.Pagination) (*models.ContentVersionResponse, error)
	RevertContentVersion(ctx context.Context, contentID uuid.UUID, version int) (*models.ContentItem, error)

	// Utility operations
	ValidateContentValue(contentType models.ContentType, value models.ContentValue) error
	GetContentStatistics(ctx context.Context) (map[string]interface{}, error)
}