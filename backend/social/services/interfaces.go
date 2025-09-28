package services

import (
	"context"

	"github.com/google/uuid"
	"tchat/social/models"
)

// UserService defines the interface for user-related social operations
type UserService interface {
	// Profile management
	GetSocialProfile(ctx context.Context, userID uuid.UUID) (*models.SocialProfile, error)
	UpdateSocialProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateSocialProfileRequest) (*models.SocialProfile, error)

	// User discovery and relationships
	DiscoverUsers(ctx context.Context, req *models.UserDiscoveryRequest) ([]*models.SocialProfile, error)
	FollowUser(ctx context.Context, req *models.FollowRequest) error
	UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error
	GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) (map[string]interface{}, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) (map[string]interface{}, error)
	GetConnections(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)

	// Analytics
	GetUserAnalytics(ctx context.Context, userID uuid.UUID, period string) (*models.UserAnalyticsResponse, error)
}

// PostService defines the interface for post-related social operations
type PostService interface {
	// Post management
	CreatePost(ctx context.Context, authorID uuid.UUID, req *models.CreatePostRequest) (*models.Post, error)
	GetPost(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	UpdatePost(ctx context.Context, postID, authorID uuid.UUID, req *models.UpdatePostRequest) (*models.Post, error)
	DeletePost(ctx context.Context, postID, authorID uuid.UUID) error

	// Interactions
	AddReaction(ctx context.Context, userID uuid.UUID, req *models.CreateReactionRequest) error
	RemoveReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) error
	CreateComment(ctx context.Context, authorID uuid.UUID, req *models.CreateCommentRequest) (*models.Comment, error)

	// Feed and discovery
	GetSocialFeed(ctx context.Context, userID uuid.UUID, req *models.SocialFeedRequest) (*models.SocialFeed, error)
	GetTrendingContent(ctx context.Context, req *models.TrendingRequest) (*models.TrendingContent, error)
	ShareContent(ctx context.Context, userID uuid.UUID, req *models.ShareRequest) error
}

// CommunityService defines the interface for community-related social operations
type CommunityService interface {
	// Community management
	CreateCommunity(ctx context.Context, creatorID uuid.UUID, req *models.CreateCommunityRequest) (*models.Community, error)
	GetCommunity(ctx context.Context, communityID uuid.UUID) (*models.Community, error)
	UpdateCommunity(ctx context.Context, communityID, userID uuid.UUID, req *models.UpdateCommunityRequest) (*models.Community, error)
	DeleteCommunity(ctx context.Context, communityID, userID uuid.UUID) error

	// Membership management
	JoinCommunity(ctx context.Context, communityID, userID uuid.UUID, req *models.JoinCommunityRequest) error
	LeaveCommunity(ctx context.Context, communityID, userID uuid.UUID) error
	GetCommunityMembers(ctx context.Context, communityID uuid.UUID, limit, offset int) (map[string]interface{}, error)
	AssignRole(ctx context.Context, communityID, assignerID uuid.UUID, req *models.CommunityRoleRequest) error

	// Community discovery and content
	DiscoverCommunities(ctx context.Context, req *models.CommunityDiscoveryRequest) ([]*models.Community, error)
	GetCommunityFeed(ctx context.Context, communityID uuid.UUID, limit, offset int) (map[string]interface{}, error)

	// Analytics and insights
	GetCommunityAnalytics(ctx context.Context, communityID uuid.UUID, period string) (*models.CommunityAnalytics, error)
	GetCommunityInsights(ctx context.Context, communityID uuid.UUID) (*models.CommunityInsights, error)
}

// ModerationService defines the interface for moderation and safety operations
type ModerationService interface {
	// Report management
	CreateReport(ctx context.Context, reporterID uuid.UUID, req *models.CreateReportRequest) (*models.Report, error)
	GetReports(ctx context.Context, communityID *uuid.UUID, status string, limit, offset int) ([]*models.Report, error)
	UpdateReport(ctx context.Context, reportID, moderatorID uuid.UUID, req *models.UpdateReportRequest) (*models.Report, error)

	// Moderation actions
	TakeModerationAction(ctx context.Context, moderatorID uuid.UUID, req *models.CreateModerationActionRequest) (*models.ModerationAction, error)
	GetModerationActions(ctx context.Context, targetID uuid.UUID, targetType string) ([]*models.ModerationAction, error)

	// Safety guidelines
	GetSafetyGuidelines(ctx context.Context, communityID *uuid.UUID) (*models.SafetyGuidelines, error)
	UpdateSafetyGuidelines(ctx context.Context, communityID *uuid.UUID, guidelines *models.SafetyGuidelines) error

	// Moderation stats
	GetModerationStats(ctx context.Context, communityID *uuid.UUID, period string) (*models.ModerationStats, error)
}

// AnalyticsService defines the interface for social analytics operations
type AnalyticsService interface {
	// User analytics
	TrackUserActivity(ctx context.Context, activity *models.UserActivity) error
	GetUserEngagementMetrics(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error)

	// Content analytics
	GetPostPerformance(ctx context.Context, postID uuid.UUID) (map[string]interface{}, error)
	GetTrendingTopics(ctx context.Context, region string, timeframe string) ([]string, error)

	// Community analytics
	GetCommunityGrowthMetrics(ctx context.Context, communityID uuid.UUID, period string) (map[string]interface{}, error)
	GetRegionalInsights(ctx context.Context, region string, period string) (map[string]interface{}, error)
}