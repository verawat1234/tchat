package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tchat/social/models"
)

// BatchOperations provides efficient bulk database operations for social media workloads
type BatchOperations struct {
	db        *gorm.DB
	batchSize int
}

// NewBatchOperations creates a new batch operations instance
func NewBatchOperations(db *gorm.DB) *BatchOperations {
	return &BatchOperations{
		db:        db,
		batchSize: 100, // Default batch size optimized for social media operations
	}
}

// SetBatchSize configures the batch size for operations
func (b *BatchOperations) SetBatchSize(size int) {
	if size > 0 && size <= 1000 {
		b.batchSize = size
	}
}

// BatchCreatePosts efficiently creates multiple posts in batches
func (b *BatchOperations) BatchCreatePosts(ctx context.Context, posts []*models.Post) error {
	if len(posts) == 0 {
		return nil
	}

	return b.batchCreate(ctx, posts, "posts")
}

// BatchCreateComments efficiently creates multiple comments in batches
func (b *BatchOperations) BatchCreateComments(ctx context.Context, comments []*models.Comment) error {
	if len(comments) == 0 {
		return nil
	}

	return b.batchCreate(ctx, comments, "comments")
}

// BatchCreateReactions efficiently creates multiple reactions in batches
func (b *BatchOperations) BatchCreateReactions(ctx context.Context, reactions []*models.Reaction) error {
	if len(reactions) == 0 {
		return nil
	}

	return b.batchCreate(ctx, reactions, "reactions")
}

// BatchCreateFollows efficiently creates multiple follow relationships in batches
func (b *BatchOperations) BatchCreateFollows(ctx context.Context, follows []*models.Follow) error {
	if len(follows) == 0 {
		return nil
	}

	return b.batchCreate(ctx, follows, "follows")
}

// BatchUpsertSocialProfiles efficiently upserts multiple social profiles
func (b *BatchOperations) BatchUpsertSocialProfiles(ctx context.Context, profiles []*models.SocialProfile) error {
	if len(profiles) == 0 {
		return nil
	}

	return b.processBatches(ctx, profiles, func(batch []interface{}) error {
		return b.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"display_name", "bio", "interests", "social_links", "updated_at"}),
		}).Create(batch).Error
	})
}

// BatchUpdatePostMetrics efficiently updates post engagement metrics
func (b *BatchOperations) BatchUpdatePostMetrics(ctx context.Context, updates []PostMetricUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	return b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			result := tx.Model(&models.Post{}).
				Where("id = ?", update.PostID).
				Updates(map[string]interface{}{
					"likes_count":    update.LikesCount,
					"comments_count": update.CommentsCount,
					"shares_count":   update.SharesCount,
					"updated_at":     time.Now(),
				})

			if result.Error != nil {
				return fmt.Errorf("failed to update post %s metrics: %w", update.PostID, result.Error)
			}
		}
		return nil
	})
}

// BatchDeleteSoftly efficiently soft deletes multiple records
func (b *BatchOperations) BatchDeleteSoftly(ctx context.Context, tableName string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	return b.processBatches(ctx, stringSliceToInterface(ids), func(batch []interface{}) error {
		query := fmt.Sprintf("UPDATE %s SET is_deleted = true, updated_at = CURRENT_TIMESTAMP WHERE id = ANY(?)", tableName)
		return b.db.WithContext(ctx).Exec(query, batch).Error
	})
}

// BatchAnalyzeEngagement analyzes engagement patterns for posts in batches
func (b *BatchOperations) BatchAnalyzeEngagement(ctx context.Context, postIDs []string, since time.Time) ([]EngagementAnalysis, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	var results []EngagementAnalysis

	err := b.processBatches(ctx, stringSliceToInterface(postIDs), func(batch []interface{}) error {
		query := `
			WITH post_engagement AS (
				SELECT
					p.id,
					p.likes_count,
					p.comments_count,
					p.shares_count,
					COUNT(DISTINCT r.user_id) as unique_reactors,
					COUNT(DISTINCT c.author_id) as unique_commenters,
					COUNT(DISTINCT s.user_id) as unique_sharers,
					AVG(CASE WHEN r.type = 'like' THEN 1.0 ELSE 0.0 END) as like_ratio,
					COUNT(CASE WHEN r.created_at > ? THEN 1 END) as recent_reactions
				FROM posts p
				LEFT JOIN reactions r ON p.id = r.target_id AND r.target_type = 'post'
				LEFT JOIN comments c ON p.id = c.post_id
				LEFT JOIN shares s ON p.id = s.content_id AND s.content_type = 'post'
				WHERE p.id = ANY(?)
				GROUP BY p.id, p.likes_count, p.comments_count, p.shares_count
			)
			SELECT
				id,
				likes_count,
				comments_count,
				shares_count,
				unique_reactors,
				unique_commenters,
				unique_sharers,
				like_ratio,
				recent_reactions,
				(likes_count + comments_count * 2 + shares_count * 3) as engagement_score
			FROM post_engagement
		`

		var batchResults []EngagementAnalysis
		if err := b.db.WithContext(ctx).Raw(query, since, batch).Scan(&batchResults).Error; err != nil {
			return fmt.Errorf("failed to analyze engagement for batch: %w", err)
		}

		results = append(results, batchResults...)
		return nil
	})

	return results, err
}

// BatchOptimizeQueries executes query optimization for social data
func (b *BatchOperations) BatchOptimizeQueries(ctx context.Context) error {
	optimizations := []string{
		"ANALYZE posts",
		"ANALYZE comments",
		"ANALYZE reactions",
		"ANALYZE follows",
		"ANALYZE social_profiles",
		"ANALYZE communities",
		"ANALYZE community_members",
		"ANALYZE shares",
		"REINDEX INDEX CONCURRENTLY idx_posts_author_created_visibility",
		"REINDEX INDEX CONCURRENTLY idx_reactions_popular_content",
		"REINDEX INDEX CONCURRENTLY idx_follows_mutual_connections",
	}

	for _, optimization := range optimizations {
		if err := b.db.WithContext(ctx).Exec(optimization).Error; err != nil {
			// Log but continue on optimization errors
			fmt.Printf("Optimization warning: %s - %v\n", optimization, err)
		}
	}

	return nil
}

// Generic batch processing helper
func (b *BatchOperations) batchCreate(ctx context.Context, items interface{}, tableName string) error {
	return b.processBatches(ctx, items, func(batch []interface{}) error {
		return b.db.WithContext(ctx).CreateInBatches(batch, b.batchSize).Error
	})
}

// processBatches processes items in configurable batch sizes
func (b *BatchOperations) processBatches(ctx context.Context, items interface{}, processor func([]interface{}) error) error {
	var itemSlice []interface{}

	// Convert different slice types to interface slice
	switch v := items.(type) {
	case []*models.Post:
		for _, item := range v {
			itemSlice = append(itemSlice, item)
		}
	case []*models.Comment:
		for _, item := range v {
			itemSlice = append(itemSlice, item)
		}
	case []*models.Reaction:
		for _, item := range v {
			itemSlice = append(itemSlice, item)
		}
	case []*models.Follow:
		for _, item := range v {
			itemSlice = append(itemSlice, item)
		}
	case []*models.SocialProfile:
		for _, item := range v {
			itemSlice = append(itemSlice, item)
		}
	case []interface{}:
		itemSlice = v
	default:
		return fmt.Errorf("unsupported item type for batch processing")
	}

	// Process in batches
	for i := 0; i < len(itemSlice); i += b.batchSize {
		end := i + b.batchSize
		if end > len(itemSlice) {
			end = len(itemSlice)
		}

		batch := itemSlice[i:end]
		if err := processor(batch); err != nil {
			return fmt.Errorf("failed to process batch %d-%d: %w", i, end, err)
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

// Helper function to convert string slice to interface slice
func stringSliceToInterface(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

// Data structures for batch operations

// PostMetricUpdate represents a post metrics update operation
type PostMetricUpdate struct {
	PostID        string `json:"post_id"`
	LikesCount    int    `json:"likes_count"`
	CommentsCount int    `json:"comments_count"`
	SharesCount   int    `json:"shares_count"`
}

// EngagementAnalysis represents engagement analysis results
type EngagementAnalysis struct {
	ID               string  `json:"id"`
	LikesCount       int     `json:"likes_count"`
	CommentsCount    int     `json:"comments_count"`
	SharesCount      int     `json:"shares_count"`
	UniqueReactors   int     `json:"unique_reactors"`
	UniqueCommenters int     `json:"unique_commenters"`
	UniqueSharers    int     `json:"unique_sharers"`
	LikeRatio        float64 `json:"like_ratio"`
	RecentReactions  int     `json:"recent_reactions"`
	EngagementScore  int     `json:"engagement_score"`
}