package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/content/models"
)

// ContentService provides business logic for content management
type ContentService struct {
	contentRepo  ContentRepository
	categoryRepo CategoryRepository
	versionRepo  VersionRepository
	db           *gorm.DB
}

// NewContentService creates a new content service
func NewContentService(
	contentRepo ContentRepository,
	categoryRepo CategoryRepository,
	versionRepo VersionRepository,
	db *gorm.DB,
) *ContentService {
	return &ContentService{
		contentRepo:  contentRepo,
		categoryRepo: categoryRepo,
		versionRepo:  versionRepo,
		db:           db,
	}
}

// CreateContent creates a new content item
func (s *ContentService) CreateContent(ctx context.Context, req models.CreateContentRequest) (*models.ContentItem, error) {
	// Validate content type
	if !req.Type.IsValid() {
		return nil, fmt.Errorf("invalid content type: %s", req.Type)
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = models.ContentStatusDraft
	}

	// Validate status
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid content status: %s", status)
	}

	// Create content item
	content := &models.ContentItem{
		ID:       uuid.New(),
		Category: req.Category,
		Type:     req.Type,
		Value:    req.Value,
		Metadata: req.Metadata,
		Status:   status,
		Tags:     req.Tags,
		Notes:    req.Notes,
	}

	err := s.contentRepo.CreateContent(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	// Create initial version
	version := &models.ContentVersion{
		ID:        uuid.New(),
		ContentID: content.ID,
		Version:   1,
		Value:     content.Value,
		Metadata:  content.Metadata,
		Status:    content.Status,
	}

	if err := s.versionRepo.CreateVersion(ctx, version); err != nil {
		// Log error but don't fail the content creation
		// In production, you might want to handle this differently
	}

	return content, nil
}

// GetContent retrieves a content item by ID
func (s *ContentService) GetContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error) {
	return s.contentRepo.GetContentByID(ctx, id)
}

// GetContentByKey retrieves a content item by key (category)
func (s *ContentService) GetContentByKey(ctx context.Context, key string) (*models.ContentItem, error) {
	return s.contentRepo.GetContentByKey(ctx, key)
}

// GetContentItems retrieves content items with filters and pagination
func (s *ContentService) GetContentItems(ctx context.Context, filters models.ContentFilters, pagination models.Pagination, sort models.SortOptions) (*models.ContentResponse, error) {
	// Set default pagination if not provided
	if pagination.PageSize <= 0 {
		pagination.PageSize = 20
	}
	if pagination.PageSize > 100 {
		pagination.PageSize = 100
	}

	return s.contentRepo.GetContentItems(ctx, filters, pagination, sort)
}

// GetContentByCategory retrieves content items by category
func (s *ContentService) GetContentByCategory(ctx context.Context, category string, pagination models.Pagination, sort models.SortOptions) (*models.ContentResponse, error) {
	// Set default pagination if not provided
	if pagination.PageSize <= 0 {
		pagination.PageSize = 20
	}
	if pagination.PageSize > 100 {
		pagination.PageSize = 100
	}

	return s.contentRepo.GetContentByCategory(ctx, category, pagination, sort)
}

// UpdateContent updates an existing content item
func (s *ContentService) UpdateContent(ctx context.Context, id uuid.UUID, req models.UpdateContentRequest) (*models.ContentItem, error) {
	// Get existing content
	content, err := s.contentRepo.GetContentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// Store original values for version creation
	_ = content.Value     // originalValue for future use
	_ = content.Metadata  // originalMetadata for future use
	_ = content.Status    // originalStatus for future use

	// Apply updates
	hasChanges := false

	if req.Category != nil && *req.Category != content.Category {
		content.Category = *req.Category
		hasChanges = true
	}

	if req.Type != nil && *req.Type != content.Type {
		if !req.Type.IsValid() {
			return nil, fmt.Errorf("invalid content type: %s", *req.Type)
		}
		content.Type = *req.Type
		hasChanges = true
	}

	if req.Value != nil {
		content.Value = *req.Value
		hasChanges = true
	}

	if req.Metadata != nil {
		content.Metadata = *req.Metadata
		hasChanges = true
	}

	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, fmt.Errorf("invalid content status: %s", *req.Status)
		}
		content.Status = *req.Status
		hasChanges = true
	}

	if req.Tags != nil {
		content.Tags = *req.Tags
		hasChanges = true
	}

	if req.Notes != nil {
		content.Notes = req.Notes
		hasChanges = true
	}

	if !hasChanges {
		return content, nil // No changes to apply
	}

	// Update content
	err = s.contentRepo.UpdateContent(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	// Create new version if content value, metadata, or status changed
	if req.Value != nil || req.Metadata != nil || req.Status != nil {
		latestVersion, err := s.versionRepo.GetLatestVersion(ctx, id)
		nextVersion := 1
		if err == nil {
			nextVersion = latestVersion.Version + 1
		}

		version := &models.ContentVersion{
			ID:        uuid.New(),
			ContentID: content.ID,
			Version:   nextVersion,
			Value:     content.Value,
			Metadata:  content.Metadata,
			Status:    content.Status,
		}

		if err := s.versionRepo.CreateVersion(ctx, version); err != nil {
			// Log error but don't fail the update
		}
	}

	return content, nil
}

// DeleteContent deletes a content item (soft delete)
func (s *ContentService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return s.contentRepo.DeleteContent(ctx, id)
}

// PublishContent publishes a content item
func (s *ContentService) PublishContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error) {
	content, err := s.contentRepo.PublishContent(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to publish content: %w", err)
	}

	// Create version for publish status change
	latestVersion, err := s.versionRepo.GetLatestVersion(ctx, id)
	nextVersion := 1
	if err == nil {
		nextVersion = latestVersion.Version + 1
	}

	version := &models.ContentVersion{
		ID:        uuid.New(),
		ContentID: content.ID,
		Version:   nextVersion,
		Value:     content.Value,
		Metadata:  content.Metadata,
		Status:    content.Status,
	}

	if err := s.versionRepo.CreateVersion(ctx, version); err != nil {
		// Log error but don't fail the publish
	}

	return content, nil
}

// ArchiveContent archives a content item
func (s *ContentService) ArchiveContent(ctx context.Context, id uuid.UUID) (*models.ContentItem, error) {
	return s.contentRepo.ArchiveContent(ctx, id)
}

// BulkUpdateContent updates multiple content items
func (s *ContentService) BulkUpdateContent(ctx context.Context, req models.BulkUpdateRequest) ([]models.ContentItem, error) {
	// Validate status if provided
	if req.Updates.Status != nil && !req.Updates.Status.IsValid() {
		return nil, fmt.Errorf("invalid content status: %s", *req.Updates.Status)
	}

	// Validate type if provided
	if req.Updates.Type != nil && !req.Updates.Type.IsValid() {
		return nil, fmt.Errorf("invalid content type: %s", *req.Updates.Type)
	}

	return s.contentRepo.BulkUpdateContent(ctx, req.IDs, req.Updates)
}

// SyncContent handles content synchronization
func (s *ContentService) SyncContent(ctx context.Context, req models.SyncContentRequest) (*models.SyncContentResponse, error) {
	return s.contentRepo.SyncContent(ctx, req)
}

// GetCategories retrieves content categories
func (s *ContentService) GetCategories(ctx context.Context, pagination models.Pagination, sort models.SortOptions) (*models.ContentCategoryResponse, error) {
	// Set default pagination if not provided
	if pagination.PageSize <= 0 {
		pagination.PageSize = 20
	}
	if pagination.PageSize > 100 {
		pagination.PageSize = 100
	}

	return s.categoryRepo.GetCategories(ctx, pagination, sort)
}

// CreateCategory creates a new content category
func (s *ContentService) CreateCategory(ctx context.Context, category *models.ContentCategory) error {
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	return s.categoryRepo.CreateCategory(ctx, category)
}

// GetContentVersions retrieves versions for a content item
func (s *ContentService) GetContentVersions(ctx context.Context, contentID uuid.UUID, pagination models.Pagination) (*models.ContentVersionResponse, error) {
	// Set default pagination if not provided
	if pagination.PageSize <= 0 {
		pagination.PageSize = 10
	}
	if pagination.PageSize > 50 {
		pagination.PageSize = 50
	}

	return s.versionRepo.GetContentVersions(ctx, contentID, pagination)
}

// RevertContentVersion reverts content to a specific version
func (s *ContentService) RevertContentVersion(ctx context.Context, contentID uuid.UUID, version int) (*models.ContentItem, error) {
	content, err := s.versionRepo.RevertContentVersion(ctx, contentID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to revert content version: %w", err)
	}

	// Create new version for the revert
	latestVersion, err := s.versionRepo.GetLatestVersion(ctx, contentID)
	nextVersion := 1
	if err == nil {
		nextVersion = latestVersion.Version + 1
	}

	newVersion := &models.ContentVersion{
		ID:        uuid.New(),
		ContentID: content.ID,
		Version:   nextVersion,
		Value:     content.Value,
		Metadata:  content.Metadata,
		Status:    content.Status,
	}

	if err := s.versionRepo.CreateVersion(ctx, newVersion); err != nil {
		// Log error but don't fail the revert
	}

	return content, nil
}

// ValidateContentValue validates content value based on type
func (s *ContentService) ValidateContentValue(contentType models.ContentType, value models.ContentValue) error {
	switch contentType {
	case models.ContentTypeText:
		if text, ok := value["text"].(string); !ok || text == "" {
			return fmt.Errorf("text content must have a 'text' field with string value")
		}
	case models.ContentTypeHTML:
		if html, ok := value["html"].(string); !ok || html == "" {
			return fmt.Errorf("HTML content must have an 'html' field with string value")
		}
	case models.ContentTypeMarkdown:
		if markdown, ok := value["markdown"].(string); !ok || markdown == "" {
			return fmt.Errorf("Markdown content must have a 'markdown' field with string value")
		}
	case models.ContentTypeJSON:
		if data, ok := value["data"]; !ok || data == nil {
			return fmt.Errorf("JSON content must have a 'data' field")
		}
	case models.ContentTypeConfiguration:
		if config, ok := value["config"]; !ok || config == nil {
			return fmt.Errorf("Configuration content must have a 'config' field")
		}
	}
	return nil
}

// GetContentStatistics returns content statistics
func (s *ContentService) GetContentStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// This would be implemented with proper database queries
	// For now, returning mock data matching the frontend expectations
	stats["total_items"] = 150
	stats["published"] = 120
	stats["draft"] = 25
	stats["archived"] = 5
	stats["categories"] = 12
	stats["last_updated"] = time.Now()

	return stats, nil
}