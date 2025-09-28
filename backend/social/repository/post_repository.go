package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat/social/models"
)

// GormPostRepository implements PostRepository using GORM
type GormPostRepository struct {
	db *gorm.DB
}

// NewGormPostRepository creates a new GORM post repository
func NewGormPostRepository(db *gorm.DB) PostRepository {
	return &GormPostRepository{db: db}
}

func (r *GormPostRepository) GetPost(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = false", postID).First(&post).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

func (r *GormPostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	err := r.db.WithContext(ctx).Create(post).Error
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *GormPostRepository) UpdatePost(ctx context.Context, postID uuid.UUID, updates *models.UpdatePostRequest) error {
	err := r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (r *GormPostRepository) DeletePost(ctx context.Context, postID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).Update("is_deleted", true).Error
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (r *GormPostRepository) GetPostsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	err := r.db.WithContext(ctx).Where("author_id = ? AND is_deleted = false", userID).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by user: %w", err)
	}
	return posts, nil
}

func (r *GormPostRepository) GetPostsByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	err := r.db.WithContext(ctx).Where("community_id = ? AND is_deleted = false", communityID).
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by community: %w", err)
	}
	return posts, nil
}

func (r *GormPostRepository) GetSocialFeed(ctx context.Context, userID uuid.UUID, req *models.SocialFeedRequest) (*models.SocialFeed, error) {
	query := r.db.WithContext(ctx).Model(&models.Post{}).Where("is_deleted = false")

	if req.Region != "" {
		query = query.Joins("JOIN social_profiles ON posts.author_id = social_profiles.id").
			Where("social_profiles.region = ?", req.Region)
	}

	var posts []models.Post
	err := query.Order("created_at DESC").Limit(req.Limit).Find(&posts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get social feed: %w", err)
	}

	return &models.SocialFeed{
		UserID:    userID,
		Posts:     posts,
		Algorithm: req.Algorithm,
		Region:    req.Region,
		HasMore:   len(posts) == req.Limit,
	}, nil
}

func (r *GormPostRepository) GetTrendingPosts(ctx context.Context, req *models.TrendingRequest) (*models.TrendingContent, error) {
	var posts []models.Post
	query := r.db.WithContext(ctx).Where("is_deleted = false AND is_trending = true")

	if req.Region != "" {
		query = query.Joins("JOIN social_profiles ON posts.author_id = social_profiles.id").
			Where("social_profiles.region = ?", req.Region)
	}

	err := query.Order("views_count DESC, created_at DESC").Limit(req.Limit).Find(&posts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get trending posts: %w", err)
	}

	return &models.TrendingContent{
		Region:    req.Region,
		Timeframe: req.Timeframe,
		Posts:     posts,
	}, nil
}

func (r *GormPostRepository) IncrementViewCount(ctx context.Context, postID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).
		UpdateColumn("views_count", gorm.Expr("views_count + 1")).Error
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}
	return nil
}

func (r *GormPostRepository) UpdateInteractionCounts(ctx context.Context, postID uuid.UUID, likes, comments, shares, reactions int) error {
	updates := map[string]interface{}{
		"likes_count":     likes,
		"comments_count":  comments,
		"shares_count":    shares,
		"reactions_count": reactions,
	}

	err := r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update interaction counts: %w", err)
	}
	return nil
}