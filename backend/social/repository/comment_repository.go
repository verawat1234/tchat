package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat/social/models"
)

// GormCommentRepository implements CommentRepository using GORM
type GormCommentRepository struct {
	db *gorm.DB
}

// NewGormCommentRepository creates a new GORM comment repository
func NewGormCommentRepository(db *gorm.DB) CommentRepository {
	return &GormCommentRepository{db: db}
}

func (r *GormCommentRepository) GetComment(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = false", commentID).First(&comment).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return &comment, nil
}

func (r *GormCommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	err := r.db.WithContext(ctx).Create(comment).Error
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

func (r *GormCommentRepository) UpdateComment(ctx context.Context, commentID uuid.UUID, content string, metadata map[string]interface{}) error {
	updates := map[string]interface{}{
		"content":  content,
		"metadata": metadata,
	}
	err := r.db.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", commentID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

func (r *GormCommentRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", commentID).Update("is_deleted", true).Error
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

func (r *GormCommentRepository) GetCommentsByPost(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := r.db.WithContext(ctx).Where("post_id = ? AND is_deleted = false", postID).
		Order("created_at ASC").Limit(limit).Offset(offset).Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post: %w", err)
	}
	return comments, nil
}

func (r *GormCommentRepository) GetCommentReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := r.db.WithContext(ctx).Where("parent_id = ? AND is_deleted = false", parentID).
		Order("created_at ASC").Limit(limit).Offset(offset).Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get comment replies: %w", err)
	}
	return comments, nil
}

func (r *GormCommentRepository) UpdateInteractionCounts(ctx context.Context, commentID uuid.UUID, likes, replies, reactions int) error {
	updates := map[string]interface{}{
		"likes_count":   likes,
		"replies_count": replies,
	}
	err := r.db.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", commentID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update interaction counts: %w", err)
	}
	return nil
}