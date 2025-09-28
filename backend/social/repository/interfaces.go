package repository

import (
	"context"

	"github.com/google/uuid"
	"tchat/social/models"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	GetSocialProfile(ctx context.Context, userID uuid.UUID) (*models.SocialProfile, error)
	CreateSocialProfile(ctx context.Context, profile *models.SocialProfile) error
	UpdateSocialProfile(ctx context.Context, userID uuid.UUID, updates *models.UpdateSocialProfileRequest) error
	GetUserFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error)
	GetUserFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error)
	CreateFollow(ctx context.Context, follow *models.Follow) error
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.UserActivity, error)
	CreateUserActivity(ctx context.Context, activity *models.UserActivity) error
	DiscoverUsers(ctx context.Context, req *models.UserDiscoveryRequest) ([]*models.SocialProfile, error)
}

// PostRepository defines the interface for post-related database operations
type PostRepository interface {
	GetPost(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	CreatePost(ctx context.Context, post *models.Post) error
	UpdatePost(ctx context.Context, postID uuid.UUID, updates *models.UpdatePostRequest) error
	DeletePost(ctx context.Context, postID uuid.UUID) error
	GetPostsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error)
	GetPostsByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]*models.Post, error)
	GetSocialFeed(ctx context.Context, userID uuid.UUID, req *models.SocialFeedRequest) (*models.SocialFeed, error)
	GetTrendingPosts(ctx context.Context, req *models.TrendingRequest) (*models.TrendingContent, error)
	IncrementViewCount(ctx context.Context, postID uuid.UUID) error
	UpdateInteractionCounts(ctx context.Context, postID uuid.UUID, likes, comments, shares, reactions int) error
}

// CommentRepository defines the interface for comment-related database operations
type CommentRepository interface {
	GetComment(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
	CreateComment(ctx context.Context, comment *models.Comment) error
	UpdateComment(ctx context.Context, commentID uuid.UUID, content string, metadata map[string]interface{}) error
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
	GetCommentsByPost(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*models.Comment, error)
	GetCommentReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*models.Comment, error)
	UpdateInteractionCounts(ctx context.Context, commentID uuid.UUID, likes, replies, reactions int) error
}

// ReactionRepository defines the interface for reaction-related database operations
type ReactionRepository interface {
	GetReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) (*models.Reaction, error)
	CreateReaction(ctx context.Context, reaction *models.Reaction) error
	UpdateReaction(ctx context.Context, reactionID uuid.UUID, reactionType string) error
	DeleteReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) error
	GetReactionsByTarget(ctx context.Context, targetID uuid.UUID, targetType string) ([]*models.Reaction, error)
	GetReactionCounts(ctx context.Context, targetID uuid.UUID, targetType string) (map[string]int, error)
}

// CommunityRepository defines the interface for community-related database operations
type CommunityRepository interface {
	GetCommunity(ctx context.Context, communityID uuid.UUID) (*models.Community, error)
	CreateCommunity(ctx context.Context, community *models.Community) error
	UpdateCommunity(ctx context.Context, communityID uuid.UUID, updates *models.UpdateCommunityRequest) error
	DeleteCommunity(ctx context.Context, communityID uuid.UUID) error
	GetCommunitiesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Community, error)
	DiscoverCommunities(ctx context.Context, req *models.CommunityDiscoveryRequest) ([]*models.Community, error)
	JoinCommunity(ctx context.Context, member *models.CommunityMember) error
	LeaveCommunity(ctx context.Context, communityID, userID uuid.UUID) error
	GetCommunityMembers(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]*models.CommunityMember, error)
	UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role string) error
	GetMembershipStatus(ctx context.Context, communityID, userID uuid.UUID) (*models.CommunityMember, error)
}

// ShareRepository defines the interface for share-related database operations
type ShareRepository interface {
	CreateShare(ctx context.Context, share *models.Share) error
	GetSharesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Share, error)
	GetSharesByContent(ctx context.Context, contentID uuid.UUID, contentType string) ([]*models.Share, error)
	UpdateShareStatus(ctx context.Context, shareID uuid.UUID, status string) error
	DeleteShare(ctx context.Context, shareID uuid.UUID) error
}

// RepositoryManager coordinates all repositories and provides transaction support
type RepositoryManager interface {
	Users() UserRepository
	Posts() PostRepository
	Comments() CommentRepository
	Reactions() ReactionRepository
	Communities() CommunityRepository
	Shares() ShareRepository

	// Transaction support
	WithTransaction(ctx context.Context, fn func(ctx context.Context, rm RepositoryManager) error) error
	Close() error
}

// BaseRepository provides common database operations
type BaseRepository interface {
	BeginTransaction(ctx context.Context) (context.Context, error)
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
	GetHealth(ctx context.Context) error
}