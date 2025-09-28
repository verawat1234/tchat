package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat/social/models"
)

// GormCommunityRepository implements CommunityRepository using GORM
type GormCommunityRepository struct {
	db *gorm.DB
}

// NewGormCommunityRepository creates a new GORM community repository
func NewGormCommunityRepository(db *gorm.DB) CommunityRepository {
	return &GormCommunityRepository{db: db}
}

func (r *GormCommunityRepository) GetCommunity(ctx context.Context, communityID uuid.UUID) (*models.Community, error) {
	var community models.Community
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = false", communityID).First(&community).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get community: %w", err)
	}
	return &community, nil
}

func (r *GormCommunityRepository) CreateCommunity(ctx context.Context, community *models.Community) error {
	err := r.db.WithContext(ctx).Create(community).Error
	if err != nil {
		return fmt.Errorf("failed to create community: %w", err)
	}
	return nil
}

func (r *GormCommunityRepository) UpdateCommunity(ctx context.Context, communityID uuid.UUID, updates *models.UpdateCommunityRequest) error {
	err := r.db.WithContext(ctx).Model(&models.Community{}).Where("id = ?", communityID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update community: %w", err)
	}
	return nil
}

func (r *GormCommunityRepository) DeleteCommunity(ctx context.Context, communityID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&models.Community{}).Where("id = ?", communityID).Update("is_deleted", true).Error
	if err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}
	return nil
}

func (r *GormCommunityRepository) GetCommunitiesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Community, error) {
	var communities []*models.Community
	err := r.db.WithContext(ctx).
		Joins("JOIN community_members ON communities.id = community_members.community_id").
		Where("community_members.user_id = ? AND communities.is_deleted = false", userID).
		Order("communities.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&communities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get communities by user: %w", err)
	}
	return communities, nil
}

func (r *GormCommunityRepository) DiscoverCommunities(ctx context.Context, req *models.CommunityDiscoveryRequest) ([]*models.Community, error) {
	query := r.db.WithContext(ctx).Where("is_deleted = false AND is_public = true")

	if req.Region != "" {
		query = query.Where("region = ?", req.Region)
	}

	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	if len(req.Tags) > 0 {
		query = query.Where("tags && ?", req.Tags)
	}

	var communities []*models.Community
	err := query.Order("members_count DESC, created_at DESC").
		Limit(req.Limit).Offset(req.Offset).
		Find(&communities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to discover communities: %w", err)
	}
	return communities, nil
}

func (r *GormCommunityRepository) JoinCommunity(ctx context.Context, member *models.CommunityMember) error {
	err := r.db.WithContext(ctx).Create(member).Error
	if err != nil {
		return fmt.Errorf("failed to join community: %w", err)
	}

	// Increment member count
	err = r.db.WithContext(ctx).Model(&models.Community{}).
		Where("id = ?", member.CommunityID).
		UpdateColumn("members_count", gorm.Expr("members_count + 1")).Error
	if err != nil {
		return fmt.Errorf("failed to update member count: %w", err)
	}

	return nil
}

func (r *GormCommunityRepository) LeaveCommunity(ctx context.Context, communityID, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("community_id = ? AND user_id = ?", communityID, userID).Delete(&models.CommunityMember{}).Error
	if err != nil {
		return fmt.Errorf("failed to leave community: %w", err)
	}

	// Decrement member count
	err = r.db.WithContext(ctx).Model(&models.Community{}).
		Where("id = ?", communityID).
		UpdateColumn("members_count", gorm.Expr("members_count - 1")).Error
	if err != nil {
		return fmt.Errorf("failed to update member count: %w", err)
	}

	return nil
}

func (r *GormCommunityRepository) GetCommunityMembers(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]*models.CommunityMember, error) {
	var members []*models.CommunityMember
	err := r.db.WithContext(ctx).Where("community_id = ?", communityID).
		Order("joined_at ASC").Limit(limit).Offset(offset).Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get community members: %w", err)
	}
	return members, nil
}

func (r *GormCommunityRepository) UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role string) error {
	err := r.db.WithContext(ctx).Model(&models.CommunityMember{}).
		Where("community_id = ? AND user_id = ?", communityID, userID).
		Update("role", role).Error
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}
	return nil
}

func (r *GormCommunityRepository) GetMembershipStatus(ctx context.Context, communityID, userID uuid.UUID) (*models.CommunityMember, error) {
	var member models.CommunityMember
	err := r.db.WithContext(ctx).Where("community_id = ? AND user_id = ?", communityID, userID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get membership status: %w", err)
	}
	return &member, nil
}