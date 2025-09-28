package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat/social/models"
)

// GormReactionRepository implements ReactionRepository using GORM
type GormReactionRepository struct {
	db *gorm.DB
}

// NewGormReactionRepository creates a new GORM reaction repository
func NewGormReactionRepository(db *gorm.DB) ReactionRepository {
	return &GormReactionRepository{db: db}
}

func (r *GormReactionRepository) GetReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) (*models.Reaction, error) {
	var reaction models.Reaction
	err := r.db.WithContext(ctx).Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).First(&reaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get reaction: %w", err)
	}
	return &reaction, nil
}

func (r *GormReactionRepository) CreateReaction(ctx context.Context, reaction *models.Reaction) error {
	err := r.db.WithContext(ctx).Create(reaction).Error
	if err != nil {
		return fmt.Errorf("failed to create reaction: %w", err)
	}
	return nil
}

func (r *GormReactionRepository) UpdateReaction(ctx context.Context, reactionID uuid.UUID, reactionType string) error {
	err := r.db.WithContext(ctx).Model(&models.Reaction{}).Where("id = ?", reactionID).Update("type", reactionType).Error
	if err != nil {
		return fmt.Errorf("failed to update reaction: %w", err)
	}
	return nil
}

func (r *GormReactionRepository) DeleteReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) error {
	err := r.db.WithContext(ctx).Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).Delete(&models.Reaction{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	return nil
}

func (r *GormReactionRepository) GetReactionsByTarget(ctx context.Context, targetID uuid.UUID, targetType string) ([]*models.Reaction, error) {
	var reactions []*models.Reaction
	err := r.db.WithContext(ctx).Where("target_id = ? AND target_type = ?", targetID, targetType).Find(&reactions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reactions by target: %w", err)
	}
	return reactions, nil
}

func (r *GormReactionRepository) GetReactionCounts(ctx context.Context, targetID uuid.UUID, targetType string) (map[string]int, error) {
	var results []struct {
		Type  string
		Count int
	}

	err := r.db.WithContext(ctx).Model(&models.Reaction{}).
		Select("type, COUNT(*) as count").
		Where("target_id = ? AND target_type = ?", targetID, targetType).
		Group("type").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get reaction counts: %w", err)
	}

	counts := make(map[string]int)
	for _, result := range results {
		counts[result.Type] = result.Count
	}

	return counts, nil
}