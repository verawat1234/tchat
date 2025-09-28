package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Review, int64, error)
	GetByProductID(ctx context.Context, productID uuid.UUID, offset, limit int) ([]*models.Review, int64, error)
	GetByBusinessID(ctx context.Context, businessID uuid.UUID, offset, limit int) ([]*models.Review, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Review, int64, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).First(&review, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *reviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, "id = ?", id).Error
}

func (r *reviewRepository) List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Review, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Review{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var reviews []*models.Review
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&reviews).Error
	return reviews, total, err
}

func (r *reviewRepository) GetByProductID(ctx context.Context, productID uuid.UUID, offset, limit int) ([]*models.Review, int64, error) {
	filters := map[string]interface{}{
		"product_id": productID,
		"type":       models.ReviewTypeProduct,
		"status":     models.ReviewStatusApproved,
	}
	return r.List(ctx, filters, offset, limit)
}

func (r *reviewRepository) GetByBusinessID(ctx context.Context, businessID uuid.UUID, offset, limit int) ([]*models.Review, int64, error) {
	filters := map[string]interface{}{
		"business_id": businessID,
		"type":        models.ReviewTypeBusiness,
		"status":      models.ReviewStatusApproved,
	}
	return r.List(ctx, filters, offset, limit)
}

func (r *reviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Review, int64, error) {
	filters := map[string]interface{}{
		"user_id": userID,
	}
	return r.List(ctx, filters, offset, limit)
}