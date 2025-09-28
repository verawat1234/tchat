package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tchat/social/models"
	"tchat/social/repository"
)

// userService implements the UserService interface
type userService struct {
	repo repository.RepositoryManager
}

// NewUserService creates a new user service
func NewUserService(repo repository.RepositoryManager) UserService {
	return &userService{
		repo: repo,
	}
}

// GetSocialProfile retrieves a user's social profile
func (s *userService) GetSocialProfile(ctx context.Context, userID uuid.UUID) (*models.SocialProfile, error) {
	profile, err := s.repo.Users().GetSocialProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get social profile: %w", err)
	}
	return profile, nil
}

// UpdateSocialProfile updates a user's social profile
func (s *userService) UpdateSocialProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateSocialProfileRequest) (*models.SocialProfile, error) {
	err := s.repo.Users().UpdateSocialProfile(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update social profile: %w", err)
	}

	// Return updated profile
	return s.GetSocialProfile(ctx, userID)
}

// DiscoverUsers finds users based on criteria
func (s *userService) DiscoverUsers(ctx context.Context, req *models.UserDiscoveryRequest) ([]*models.SocialProfile, error) {
	users, err := s.repo.Users().DiscoverUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to discover users: %w", err)
	}
	return users, nil
}

// FollowUser creates a following relationship
func (s *userService) FollowUser(ctx context.Context, req *models.FollowRequest) error {
	if req.FollowerID == req.FollowingID {
		return fmt.Errorf("cannot follow yourself")
	}

	// Check if already following
	isFollowing, err := s.repo.Users().IsFollowing(ctx, req.FollowerID, req.FollowingID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}
	if isFollowing {
		return fmt.Errorf("already following this user")
	}

	// Create follow relationship
	follow := &models.Follow{
		ID:          uuid.New(),
		FollowerID:  req.FollowerID,
		FollowingID: req.FollowingID,
		CreatedAt:   time.Now(),
	}

	err = s.repo.Users().CreateFollow(ctx, follow)
	if err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}

	return nil
}

// UnfollowUser removes a following relationship
func (s *userService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return fmt.Errorf("invalid relationship")
	}

	err := s.repo.Users().DeleteFollow(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to unfollow user: %w", err)
	}

	return nil
}

// GetFollowers retrieves a user's followers
func (s *userService) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) (map[string]interface{}, error) {
	follows, err := s.repo.Users().GetUserFollowers(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	followers := make([]map[string]interface{}, 0, len(follows))
	for _, follow := range follows {
		// Get follower profile info
		followerProfile, err := s.repo.Users().GetSocialProfile(ctx, follow.FollowerID)
		if err != nil {
			continue // Skip if profile not found
		}

		follower := map[string]interface{}{
			"id":          followerProfile.ID.String(),
			"username":    followerProfile.Username,
			"displayName": followerProfile.DisplayName,
			"avatar":      followerProfile.Avatar,
			"followedAt":  follow.CreatedAt,
			"isVerified":  followerProfile.IsSocialVerified,
		}
		followers = append(followers, follower)
	}

	return map[string]interface{}{
		"followers": followers,
		"limit":     limit,
		"offset":    offset,
		"hasMore":   len(follows) == limit,
	}, nil
}

// GetFollowing retrieves users that a user is following
func (s *userService) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) (map[string]interface{}, error) {
	follows, err := s.repo.Users().GetUserFollowing(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	following := make([]map[string]interface{}, 0, len(follows))
	for _, follow := range follows {
		// Get following profile info
		followingProfile, err := s.repo.Users().GetSocialProfile(ctx, follow.FollowingID)
		if err != nil {
			continue // Skip if profile not found
		}

		followedUser := map[string]interface{}{
			"id":          followingProfile.ID.String(),
			"username":    followingProfile.Username,
			"displayName": followingProfile.DisplayName,
			"avatar":      followingProfile.Avatar,
			"followedAt":  follow.CreatedAt,
			"isVerified":  followingProfile.IsSocialVerified,
		}
		following = append(following, followedUser)
	}

	return map[string]interface{}{
		"following": following,
		"limit":     limit,
		"offset":    offset,
		"hasMore":   len(follows) == limit,
	}, nil
}

// GetConnections retrieves mutual connections
func (s *userService) GetConnections(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Get user's followers and following
	followers, err := s.repo.Users().GetUserFollowers(ctx, userID, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	following, err := s.repo.Users().GetUserFollowing(ctx, userID, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	// Find mutual connections (users who are both followers and following)
	connections := make([]map[string]interface{}, 0)
	mutualCount := 0

	followingMap := make(map[uuid.UUID]bool)
	for _, follow := range following {
		followingMap[follow.FollowingID] = true
	}

	for _, follower := range followers {
		if followingMap[follower.FollowerID] {
			// This is a mutual connection
			profile, err := s.repo.Users().GetSocialProfile(ctx, follower.FollowerID)
			if err != nil {
				continue // Skip if profile not found
			}

			connection := map[string]interface{}{
				"id":              profile.ID.String(),
				"username":        profile.Username,
				"displayName":     profile.DisplayName,
				"avatar":          profile.Avatar,
				"mutualFollowers": 0, // Would need additional query to calculate
				"connectionType":  []string{"follower", "following"},
			}
			connections = append(connections, connection)
			mutualCount++
		}
	}

	return map[string]interface{}{
		"connections":     connections,
		"mutualCount":     mutualCount,
		"suggestionCount": len(followers) + len(following) - mutualCount,
	}, nil
}

// GetUserAnalytics retrieves user analytics
func (s *userService) GetUserAnalytics(ctx context.Context, userID uuid.UUID, period string) (*models.UserAnalyticsResponse, error) {
	// Get current follower/following counts
	followers, err := s.repo.Users().GetUserFollowers(ctx, userID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	following, err := s.repo.Users().GetUserFollowing(ctx, userID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	// Get user's posts for engagement calculations
	posts, err := s.repo.Posts().GetPostsByUser(ctx, userID, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}

	// Calculate engagement metrics
	totalLikes := 0
	totalComments := 0
	totalShares := 0
	totalViews := 0

	for _, post := range posts {
		totalLikes += post.LikesCount
		totalComments += post.CommentsCount
		totalShares += post.SharesCount
		totalViews += post.ViewsCount
	}

	averagePerPost := 0.0
	if len(posts) > 0 {
		averagePerPost = float64(totalLikes+totalComments+totalShares) / float64(len(posts))
	}

	engagementRate := 0.0
	if totalViews > 0 {
		engagementRate = float64(totalLikes+totalComments+totalShares) / float64(totalViews) * 100
	}

	analytics := &models.UserAnalyticsResponse{
		UserID: userID,
		Period: period,
		Followers: map[string]interface{}{
			"current": len(followers),
			"growth":  "+0", // Would need historical data
			"rate":    "0%",
		},
		Following: map[string]interface{}{
			"current": len(following),
			"growth":  "+0", // Would need historical data
			"rate":    "0%",
		},
		Engagement: map[string]interface{}{
			"rate":           fmt.Sprintf("%.1f%%", engagementRate),
			"totalLikes":     totalLikes,
			"totalComments":  totalComments,
			"totalShares":    totalShares,
			"averagePerPost": averagePerPost,
		},
		Reach: map[string]interface{}{
			"impressions": totalViews,
			"uniqueViews": int(float64(totalViews) * 0.65), // Estimate unique views
			"countries":   []string{"TH", "SG", "MY", "ID"}, // Would need real geo data
		},
		Growth: map[string]interface{}{
			"newFollowers":    0, // Would need historical data
			"unfollowers":     0, // Would need historical data
			"netGrowth":       0,
			"growthRate":      "0%",
			"peakDay":         "N/A",
		},
		Demographics: map[string]interface{}{
			"ageGroups": map[string]int{
				"18-24": 0, // Would need user demographic data
				"25-34": 0,
				"35-44": 0,
				"45+":   0,
			},
			"topCountries": []map[string]interface{}{
				{"country": "TH", "percentage": 0}, // Would need real geo data
				{"country": "SG", "percentage": 0},
				{"country": "MY", "percentage": 0},
				{"country": "ID", "percentage": 0},
			},
		},
		UpdatedAt: time.Now(),
	}

	return analytics, nil
}