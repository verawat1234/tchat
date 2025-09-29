package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
	sharedModels "tchat.dev/shared/models"
)

// businessRepository implements the BusinessRepository interface
type businessRepository struct {
	db *gorm.DB
}

// BusinessRepository defines the interface for business data access
type BusinessRepository interface {
	FindBusinesses(ctx context.Context, filters models.BusinessFilters, pagination models.Pagination, sort models.SortOptions) ([]*sharedModels.Business, int64, error)
	FindBusinessByID(ctx context.Context, id uuid.UUID) (*sharedModels.Business, error)
	CreateBusiness(ctx context.Context, business *sharedModels.Business) error
	UpdateBusiness(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteBusiness(ctx context.Context, id uuid.UUID) error
}

// NewBusinessRepository creates a new business repository
func NewBusinessRepository(db *gorm.DB) BusinessRepository {
	return &businessRepository{
		db: db,
	}
}

// FindBusinesses retrieves businesses with filters and pagination
func (r *businessRepository) FindBusinesses(ctx context.Context, filters models.BusinessFilters, pagination models.Pagination, sort models.SortOptions) ([]*sharedModels.Business, int64, error) {
	var businesses []*sharedModels.Business
	var total int64

	query := r.db.WithContext(ctx).Model(&sharedModels.Business{})

	// Apply filters
	if filters.Country != nil {
		query = query.Where("address_country = ?", *filters.Country)
	}

	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Search != nil {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	// Only get active businesses by default
	query = query.Where("deleted_at IS NULL")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count businesses: %w", err)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if sort.Field != "" {
		order := "DESC"
		if strings.ToLower(sort.Order) == "asc" {
			order = "ASC"
		}
		orderBy = fmt.Sprintf("%s %s", sort.Field, order)
	}
	query = query.Order(orderBy)

	// Apply pagination
	offset := pagination.Page * pagination.PageSize
	query = query.Offset(offset).Limit(pagination.PageSize)

	// Execute query
	if err := query.Find(&businesses).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find businesses: %w", err)
	}

	return businesses, total, nil
}

// FindBusinessByID retrieves a business by its ID
func (r *businessRepository) FindBusinessByID(ctx context.Context, id uuid.UUID) (*sharedModels.Business, error) {
	var business sharedModels.Business

	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&business).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("business not found")
		}
		return nil, fmt.Errorf("failed to find business: %w", err)
	}

	return &business, nil
}

// CreateBusiness creates a new business
func (r *businessRepository) CreateBusiness(ctx context.Context, business *sharedModels.Business) error {
	if err := r.db.WithContext(ctx).Create(business).Error; err != nil {
		return fmt.Errorf("failed to create business: %w", err)
	}
	return nil
}

// UpdateBusiness updates a business
func (r *businessRepository) UpdateBusiness(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&sharedModels.Business{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update business: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("business not found")
	}
	return nil
}

// DeleteBusiness soft deletes a business
func (r *businessRepository) DeleteBusiness(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&sharedModels.Business{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete business: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("business not found")
	}
	return nil
}