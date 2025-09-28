package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tchat/social/models"
	"tchat/social/repository"
)

// postService implements the PostService interface
type postService struct {
	repo repository.RepositoryManager
}

// NewPostService creates a new post service
func NewPostService(repo repository.RepositoryManager) PostService {
	return &postService{
		repo: repo,
	}
}

// CreatePost creates a new social post
func (s *postService) CreatePost(ctx context.Context, authorID uuid.UUID, req *models.CreatePostRequest) (*models.Post, error) {
	post := &models.Post{
		ID:           uuid.New(),
		AuthorID:     authorID,
		CommunityID:  req.CommunityID,
		Content:      req.Content,
		MediaURLs:    req.MediaURLs,
		Tags:         req.Tags,
		Visibility:   req.Visibility,
		Type:         req.Type,
		LikesCount:   0,
		CommentsCount: 0,
		SharesCount:  0,
		ViewsCount:   0,
		IsPinned:     false,
		Metadata:     req.Metadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Validate community exists if specified
	if req.CommunityID != nil {
		_, err := s.repo.Communities().GetCommunity(ctx, *req.CommunityID)
		if err != nil {
			return nil, fmt.Errorf("community not found")
		}
	}

	// Create post in database
	err := s.repo.Posts().CreatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetPost retrieves a specific post
func (s *postService) GetPost(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	post, err := s.repo.Posts().GetPost(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return post, nil
}

// UpdatePost updates an existing post
func (s *postService) UpdatePost(ctx context.Context, postID, authorID uuid.UUID, req *models.UpdatePostRequest) (*models.Post, error) {
	// Get existing post
	post, err := s.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check authorization
	if post.AuthorID != authorID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Update post in database
	err = s.repo.Posts().UpdatePost(ctx, postID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// Return updated post
	return s.GetPost(ctx, postID)
}

// DeletePost soft deletes a post
func (s *postService) DeletePost(ctx context.Context, postID, authorID uuid.UUID) error {
	// Get existing post
	post, err := s.GetPost(ctx, postID)
	if err != nil {
		return err
	}

	// Check authorization
	if post.AuthorID != authorID {
		return fmt.Errorf("unauthorized")
	}

	// Delete post in database
	err = s.repo.Posts().DeletePost(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// AddReaction adds a reaction to a post or comment
func (s *postService) AddReaction(ctx context.Context, userID uuid.UUID, req *models.CreateReactionRequest) error {
	// Validate target exists
	if req.TargetType == "post" {
		_, err := s.GetPost(ctx, req.TargetID)
		if err != nil {
			return fmt.Errorf("target not found")
		}
	}

	// Check if user already reacted
	existingReaction, err := s.repo.Reactions().GetReaction(ctx, userID, req.TargetID, req.TargetType)
	if err == nil && existingReaction != nil {
		return fmt.Errorf("already reacted")
	}

	// Create reaction
	reaction := &models.Reaction{
		ID:         uuid.New(),
		UserID:     userID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		Type:       req.Type,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.repo.Reactions().CreateReaction(ctx, reaction)
	if err != nil {
		return fmt.Errorf("failed to create reaction: %w", err)
	}

	return nil
}

// RemoveReaction removes a user's reaction
func (s *postService) RemoveReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) error {
	// Check if reaction exists
	_, err := s.repo.Reactions().GetReaction(ctx, userID, targetID, targetType)
	if err != nil {
		return fmt.Errorf("reaction not found")
	}

	// Delete reaction from database
	err = s.repo.Reactions().DeleteReaction(ctx, userID, targetID, targetType)
	if err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	return nil
}

// CreateComment creates a comment on a post
func (s *postService) CreateComment(ctx context.Context, authorID uuid.UUID, req *models.CreateCommentRequest) (*models.Comment, error) {
	// Validate post exists
	_, err := s.GetPost(ctx, req.PostID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	comment := &models.Comment{
		ID:           uuid.New(),
		PostID:       req.PostID,
		AuthorID:     authorID,
		ParentID:     req.ParentID,
		Content:      req.Content,
		LikesCount:   0,
		RepliesCount: 0,
		Metadata:     req.Metadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create comment in database
	err = s.repo.Comments().CreateComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetSocialFeed retrieves a personalized social feed
func (s *postService) GetSocialFeed(ctx context.Context, userID uuid.UUID, req *models.SocialFeedRequest) (*models.SocialFeed, error) {
	feed, err := s.repo.Posts().GetSocialFeed(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get social feed: %w", err)
	}
	return feed, nil
}

// GetTrendingContent retrieves trending posts and topics
func (s *postService) GetTrendingContent(ctx context.Context, req *models.TrendingRequest) (*models.TrendingContent, error) {
	trending, err := s.repo.Posts().GetTrendingPosts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending content: %w", err)
	}
	return trending, nil
}

// ShareContent shares content to external platforms
func (s *postService) ShareContent(ctx context.Context, userID uuid.UUID, req *models.ShareRequest) error {
	// Validate content exists
	if req.ContentType == "post" {
		_, err := s.GetPost(ctx, req.ContentID)
		if err != nil {
			return fmt.Errorf("content not found")
		}
	}

	// Create share record
	share := &models.Share{
		ID:          uuid.New(),
		UserID:      userID,
		ContentID:   req.ContentID,
		ContentType: req.ContentType,
		Platform:    req.Platform,
		Message:     req.Message,
		Privacy:     req.Privacy,
		Metadata:    req.Metadata,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	err := s.repo.Shares().CreateShare(ctx, share)
	if err != nil {
		return fmt.Errorf("failed to create share: %w", err)
	}

	// External sharing logic would be implemented here
	// For now, just mark as completed
	err = s.repo.Shares().UpdateShareStatus(ctx, share.ID, "completed")
	if err != nil {
		return fmt.Errorf("failed to update share status: %w", err)
	}

	return nil
}