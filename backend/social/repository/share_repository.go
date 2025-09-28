package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat/social/models"
)

// GormShareRepository implements ShareRepository using GORM
type GormShareRepository struct {
	db *gorm.DB
}

// NewGormShareRepository creates a new GORM share repository
func NewGormShareRepository(db *gorm.DB) ShareRepository {
	return &GormShareRepository{db: db}
}

func (r *GormShareRepository) CreateShare(ctx context.Context, share *models.Share) error {
	err := r.db.WithContext(ctx).Create(share).Error
	if err != nil {
		return fmt.Errorf("failed to create share: %w", err)
	}
	return nil
}

func (r *GormShareRepository) GetShare(ctx context.Context, shareID uuid.UUID) (*models.Share, error) {
	var share models.Share
	err := r.db.WithContext(ctx).Where("id = ?", shareID).First(&share).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get share: %w", err)
	}
	return &share, nil
}

func (r *GormShareRepository) GetSharesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Share, error) {
	var shares []*models.Share
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&shares).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get shares by user: %w", err)
	}
	return shares, nil
}

func (r *GormShareRepository) GetSharesByContent(ctx context.Context, contentID uuid.UUID, contentType string) ([]*models.Share, error) {
	var shares []*models.Share
	err := r.db.WithContext(ctx).Where("content_id = ? AND content_type = ?", contentID, contentType).
		Order("created_at DESC").Find(&shares).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get shares by content: %w", err)
	}
	return shares, nil
}

func (r *GormShareRepository) UpdateShareStatus(ctx context.Context, shareID uuid.UUID, status string) error {
	err := r.db.WithContext(ctx).Model(&models.Share{}).Where("id = ?", shareID).Update("status", status).Error
	if err != nil {
		return fmt.Errorf("failed to update share status: %w", err)
	}
	return nil
}

func (r *GormShareRepository) DeleteShare(ctx context.Context, shareID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("id = ?", shareID).Delete(&models.Share{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete share: %w", err)
	}
	return nil
}