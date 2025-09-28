package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat/social/models"
)

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GORM user repository
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetSocialProfile(ctx context.Context, userID uuid.UUID) (*models.SocialProfile, error) {
	var profile models.SocialProfile
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get social profile: %w", err)
	}
	return &profile, nil
}

func (r *GormUserRepository) CreateSocialProfile(ctx context.Context, profile *models.SocialProfile) error {
	err := r.db.WithContext(ctx).Create(profile).Error
	if err != nil {
		return fmt.Errorf("failed to create social profile: %w", err)
	}
	return nil
}

func (r *GormUserRepository) UpdateSocialProfile(ctx context.Context, userID uuid.UUID, updates *models.UpdateSocialProfileRequest) error {
	err := r.db.WithContext(ctx).Model(&models.SocialProfile{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update social profile: %w", err)
	}
	return nil
}

func (r *GormUserRepository) GetUserFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error) {
	var follows []*models.Follow
	err := r.db.WithContext(ctx).Where("following_id = ?", userID).Limit(limit).Offset(offset).Find(&follows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	return follows, nil
}

func (r *GormUserRepository) GetUserFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error) {
	var follows []*models.Follow
	err := r.db.WithContext(ctx).Where("follower_id = ?", userID).Limit(limit).Offset(offset).Find(&follows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	return follows, nil
}

func (r *GormUserRepository) CreateFollow(ctx context.Context, follow *models.Follow) error {
	err := r.db.WithContext(ctx).Create(follow).Error
	if err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}
	return nil
}

func (r *GormUserRepository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follow{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete follow: %w", err)
	}
	return nil
}

func (r *GormUserRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %w", err)
	}
	return count > 0, nil
}

func (r *GormUserRepository) GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.UserActivity, error) {
	var activities []*models.UserActivity
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&activities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user activity: %w", err)
	}
	return activities, nil
}

func (r *GormUserRepository) CreateUserActivity(ctx context.Context, activity *models.UserActivity) error {
	err := r.db.WithContext(ctx).Create(activity).Error
	if err != nil {
		return fmt.Errorf("failed to create user activity: %w", err)
	}
	return nil
}

func (r *GormUserRepository) DiscoverUsers(ctx context.Context, req *models.UserDiscoveryRequest) ([]*models.SocialProfile, error) {
	query := r.db.WithContext(ctx).Model(&models.SocialProfile{})

	if req.Region != "" {
		// Use shared User model region field
		query = query.Where("region = ?", req.Region)
	}

	if len(req.Interests) > 0 {
		query = query.Where("interests && ?", req.Interests)
	}

	var profiles []*models.SocialProfile
	err := query.Limit(req.Limit).Offset(req.Offset).Find(&profiles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to discover users: %w", err)
	}
	return profiles, nil
}