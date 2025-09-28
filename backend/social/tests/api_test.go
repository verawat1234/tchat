package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tchat/social/handlers"
	"tchat/social/models"
	"tchat/social/repository"
	"tchat/social/services"
	sharedModels "tchat.dev/shared/models"
)

// TestSocialAPIEndpoints tests core social API functionality
func TestSocialAPIEndpoints(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create mock repository
	mockRepo := &mockRepositoryManager{}

	// Create services
	userService := services.NewUserService(mockRepo)
	syncService := services.NewMobileSyncService(mockRepo)

	// Create handlers
	userHandler := handlers.NewUserHandler(userService)
	mobileHandler := handlers.NewMobileHandler(syncService, userService)

	// Set up router
	router := gin.New()
	v1 := router.Group("/api/v1")

	// Register routes
	social := v1.Group("/social")
	{
		social.GET("/profiles/:userId", userHandler.GetSocialProfile)
		social.PUT("/profiles/:userId", userHandler.UpdateSocialProfile)
		social.POST("/follow", userHandler.FollowUser)
		social.DELETE("/follow/:followerId/:followingId", userHandler.UnfollowUser)
		social.GET("/followers/:userId", userHandler.GetFollowers)
		social.GET("/following/:userId", userHandler.GetFollowing)
		social.GET("/discover/users", userHandler.DiscoverUsers)
	}

	// Register mobile routes
	mobileHandler.RegisterMobileRoutes(v1)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	testUserID := uuid.New()

	t.Run("Health Check", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	})

	t.Run("Get Social Profile", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/social/profiles/"+testUserID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "data")
	})

	t.Run("Mobile Health Check", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/mobile/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "healthy", response["status"])
		assert.Contains(t, response, "features")
	})

	t.Run("Mobile Discovery Feed", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/mobile/discover/"+testUserID.String()+"?region=TH&limit=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, testUserID.String(), response["userId"])
		assert.Equal(t, "TH", response["region"])
		assert.Contains(t, response, "profiles")
	})

	t.Run("Mobile Sync Profile", func(t *testing.T) {
		since := time.Now().Add(-time.Hour).Format(time.RFC3339)
		// URL encode the since parameter properly
		req, _ := http.NewRequest("GET", "/api/v1/mobile/sync/profile/"+testUserID.String()+"?since="+url.QueryEscape(since), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "userId")
		assert.Contains(t, response, "hasChanges")
		assert.Contains(t, response, "syncTime")
	})

	t.Run("Follow User Mobile API", func(t *testing.T) {
		followReq := map[string]interface{}{
			"followerId":  testUserID,
			"followingId": uuid.New(),
			"source":      "mobile_test",
		}

		body, _ := json.Marshal(followReq)
		req, _ := http.NewRequest("POST", "/api/v1/mobile/follow", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["success"].(bool))
		assert.True(t, response["followed"].(bool))
	})

	t.Run("Invalid User ID Error Handling", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/mobile/profile/invalid-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "invalid_user_id", response["error"])
	})

	t.Run("Unsupported Region Error", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/mobile/discover/"+testUserID.String()+"?region=XX", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "unsupported_region", response["error"])
	})
}

// Mock implementations for testing
type mockRepositoryManager struct{}

func (m *mockRepositoryManager) Users() repository.UserRepository     { return &mockUserRepository{} }
func (m *mockRepositoryManager) Posts() repository.PostRepository     { return &mockPostRepository{} }
func (m *mockRepositoryManager) Comments() repository.CommentRepository { return &mockCommentRepository{} }
func (m *mockRepositoryManager) Reactions() repository.ReactionRepository { return &mockReactionRepository{} }
func (m *mockRepositoryManager) Communities() repository.CommunityRepository { return &mockCommunityRepository{} }
func (m *mockRepositoryManager) Shares() repository.ShareRepository   { return &mockShareRepository{} }
func (m *mockRepositoryManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, rm repository.RepositoryManager) error) error { return nil }
func (m *mockRepositoryManager) Close() error { return nil }

type mockUserRepository struct{}

func (m *mockUserRepository) GetSocialProfile(ctx context.Context, userID uuid.UUID) (*models.SocialProfile, error) {
	return &models.SocialProfile{
		User: sharedModels.User{
			ID:          userID,
			Username:    "test_user",
			Email:       "test@example.com",
			Name:        "Test User",
			DisplayName: "Test Display",
			Country:     "TH",
			Status:      "active",
			Active:      true,
		},
		Interests:      []string{"testing"},
		FollowersCount: 10,
		FollowingCount: 5,
		PostsCount:     3,
		SocialCreatedAt: time.Now().Add(-time.Hour * 24),
		SocialUpdatedAt: time.Now(),
	}, nil
}

func (m *mockUserRepository) CreateSocialProfile(ctx context.Context, profile *models.SocialProfile) error { return nil }
func (m *mockUserRepository) UpdateSocialProfile(ctx context.Context, userID uuid.UUID, updates *models.UpdateSocialProfileRequest) error { return nil }
func (m *mockUserRepository) GetUserFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error) { return []*models.Follow{}, nil }
func (m *mockUserRepository) GetUserFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error) { return []*models.Follow{}, nil }
func (m *mockUserRepository) CreateFollow(ctx context.Context, follow *models.Follow) error { return nil }
func (m *mockUserRepository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error { return nil }
func (m *mockUserRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) { return false, nil }
func (m *mockUserRepository) GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.UserActivity, error) { return []*models.UserActivity{}, nil }
func (m *mockUserRepository) CreateUserActivity(ctx context.Context, activity *models.UserActivity) error { return nil }
func (m *mockUserRepository) DiscoverUsers(ctx context.Context, req *models.UserDiscoveryRequest) ([]*models.SocialProfile, error) { return []*models.SocialProfile{}, nil }

// Empty mock implementations for other repositories
type mockPostRepository struct{}
func (m *mockPostRepository) CreatePost(ctx context.Context, post *models.Post) error { return nil }
func (m *mockPostRepository) GetPost(ctx context.Context, postID uuid.UUID) (*models.Post, error) { return nil, nil }
func (m *mockPostRepository) UpdatePost(ctx context.Context, postID uuid.UUID, updates *models.UpdatePostRequest) error { return nil }
func (m *mockPostRepository) DeletePost(ctx context.Context, postID uuid.UUID) error { return nil }
func (m *mockPostRepository) GetPostsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error) { return []*models.Post{}, nil }
func (m *mockPostRepository) GetPostsByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]*models.Post, error) { return []*models.Post{}, nil }
func (m *mockPostRepository) GetSocialFeed(ctx context.Context, userID uuid.UUID, req *models.SocialFeedRequest) (*models.SocialFeed, error) { return &models.SocialFeed{}, nil }
func (m *mockPostRepository) GetTrendingPosts(ctx context.Context, req *models.TrendingRequest) (*models.TrendingContent, error) { return &models.TrendingContent{}, nil }
func (m *mockPostRepository) IncrementViewCount(ctx context.Context, postID uuid.UUID) error { return nil }
func (m *mockPostRepository) UpdateInteractionCounts(ctx context.Context, postID uuid.UUID, likes, comments, shares, reactions int) error { return nil }

type mockCommentRepository struct{}
func (m *mockCommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error { return nil }
func (m *mockCommentRepository) GetComment(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) { return nil, nil }
func (m *mockCommentRepository) UpdateComment(ctx context.Context, commentID uuid.UUID, content string, metadata map[string]interface{}) error { return nil }
func (m *mockCommentRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error { return nil }
func (m *mockCommentRepository) GetCommentsByPost(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*models.Comment, error) { return []*models.Comment{}, nil }
func (m *mockCommentRepository) GetCommentReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*models.Comment, error) { return []*models.Comment{}, nil }
func (m *mockCommentRepository) UpdateInteractionCounts(ctx context.Context, commentID uuid.UUID, likes, replies, reactions int) error { return nil }

type mockReactionRepository struct{}
func (m *mockReactionRepository) CreateReaction(ctx context.Context, reaction *models.Reaction) error { return nil }
func (m *mockReactionRepository) UpdateReaction(ctx context.Context, reactionID uuid.UUID, reactionType string) error { return nil }
func (m *mockReactionRepository) DeleteReaction(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, targetType string) error { return nil }
func (m *mockReactionRepository) GetReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) (*models.Reaction, error) { return nil, nil }
func (m *mockReactionRepository) GetReactionsByTarget(ctx context.Context, targetID uuid.UUID, targetType string) ([]*models.Reaction, error) { return []*models.Reaction{}, nil }
func (m *mockReactionRepository) GetReactionCounts(ctx context.Context, targetID uuid.UUID, targetType string) (map[string]int, error) { return map[string]int{}, nil }

type mockCommunityRepository struct{}
func (m *mockCommunityRepository) CreateCommunity(ctx context.Context, community *models.Community) error { return nil }
func (m *mockCommunityRepository) GetCommunity(ctx context.Context, communityID uuid.UUID) (*models.Community, error) { return nil, nil }
func (m *mockCommunityRepository) UpdateCommunity(ctx context.Context, communityID uuid.UUID, updates *models.UpdateCommunityRequest) error { return nil }
func (m *mockCommunityRepository) DeleteCommunity(ctx context.Context, communityID uuid.UUID) error { return nil }
func (m *mockCommunityRepository) GetCommunitiesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Community, error) { return []*models.Community{}, nil }
func (m *mockCommunityRepository) DiscoverCommunities(ctx context.Context, req *models.CommunityDiscoveryRequest) ([]*models.Community, error) { return []*models.Community{}, nil }
func (m *mockCommunityRepository) JoinCommunity(ctx context.Context, member *models.CommunityMember) error { return nil }
func (m *mockCommunityRepository) LeaveCommunity(ctx context.Context, communityID, userID uuid.UUID) error { return nil }
func (m *mockCommunityRepository) GetCommunityMembers(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]*models.CommunityMember, error) { return []*models.CommunityMember{}, nil }
func (m *mockCommunityRepository) UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role string) error { return nil }
func (m *mockCommunityRepository) GetMembershipStatus(ctx context.Context, communityID, userID uuid.UUID) (*models.CommunityMember, error) { return nil, nil }

type mockShareRepository struct{}
func (m *mockShareRepository) CreateShare(ctx context.Context, share *models.Share) error { return nil }
func (m *mockShareRepository) GetShare(ctx context.Context, shareID uuid.UUID) (*models.Share, error) { return nil, nil }
func (m *mockShareRepository) GetSharesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Share, error) { return []*models.Share{}, nil }
func (m *mockShareRepository) GetSharesByContent(ctx context.Context, contentID uuid.UUID, contentType string) ([]*models.Share, error) { return []*models.Share{}, nil }
func (m *mockShareRepository) UpdateShareStatus(ctx context.Context, shareID uuid.UUID, status string) error { return nil }
func (m *mockShareRepository) DeleteShare(ctx context.Context, shareID uuid.UUID) error { return nil }