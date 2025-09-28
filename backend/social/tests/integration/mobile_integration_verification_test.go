package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// Mock implementations for testing
type mockRepositoryManager struct{}

func (m *mockRepositoryManager) Users() repository.UserRepository {
	return &mockUserRepository{}
}

func (m *mockRepositoryManager) Posts() repository.PostRepository {
	return &mockPostRepository{}
}

func (m *mockRepositoryManager) Comments() repository.CommentRepository {
	return &mockCommentRepository{}
}

func (m *mockRepositoryManager) Reactions() repository.ReactionRepository {
	return &mockReactionRepository{}
}

func (m *mockRepositoryManager) Communities() repository.CommunityRepository {
	return &mockCommunityRepository{}
}

func (m *mockRepositoryManager) Shares() repository.ShareRepository {
	return &mockShareRepository{}
}

func (m *mockRepositoryManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, rm repository.RepositoryManager) error) error {
	return fn(ctx, m)
}

func (m *mockRepositoryManager) Close() error {
	return nil
}

// Mock repository implementations (minimal for testing)
type mockUserRepository struct{}

func (r *mockUserRepository) GetSocialProfile(ctx context.Context, userID uuid.UUID) (*models.SocialProfile, error) {
	// Import the shared models for the embedded User
	sharedModels := &sharedModels.User{
		ID:          userID,
		Username:    "testuser",
		Name:        "Test User",
		DisplayName: "Test User",
		Email:       "test@example.com",
		Country:     "TH",
		Locale:      "th-TH",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	return &models.SocialProfile{
		User:               *sharedModels,
		FollowersCount:     100,
		FollowingCount:     50,
		PostsCount:         25,
		IsSocialVerified:   false,
		SocialCreatedAt:    time.Now().Add(-24 * time.Hour),
		SocialUpdatedAt:    time.Now(),
	}, nil
}

func (r *mockUserRepository) CreateSocialProfile(ctx context.Context, profile *models.SocialProfile) error {
	return nil
}

func (r *mockUserRepository) UpdateSocialProfile(ctx context.Context, userID uuid.UUID, updates *models.UpdateSocialProfileRequest) error {
	return nil
}

func (r *mockUserRepository) GetUserFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error) {
	return []*models.Follow{}, nil
}

func (r *mockUserRepository) GetUserFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error) {
	return []*models.Follow{}, nil
}

func (r *mockUserRepository) CreateFollow(ctx context.Context, follow *models.Follow) error {
	return nil
}

func (r *mockUserRepository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return nil
}

func (r *mockUserRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return false, nil
}

func (r *mockUserRepository) GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.UserActivity, error) {
	return []*models.UserActivity{}, nil
}

func (r *mockUserRepository) CreateUserActivity(ctx context.Context, activity *models.UserActivity) error {
	return nil
}

func (r *mockUserRepository) DiscoverUsers(ctx context.Context, req *models.UserDiscoveryRequest) ([]*models.SocialProfile, error) {
	return []*models.SocialProfile{}, nil
}

// Mock implementations for other repositories (minimal stubs)
type mockPostRepository struct{}

func (r *mockPostRepository) GetPost(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	return &models.Post{ID: postID}, nil
}

func (r *mockPostRepository) CreatePost(ctx context.Context, post *models.Post) error { return nil }
func (r *mockPostRepository) UpdatePost(ctx context.Context, postID uuid.UUID, updates *models.UpdatePostRequest) error { return nil }
func (r *mockPostRepository) DeletePost(ctx context.Context, postID uuid.UUID) error { return nil }
func (r *mockPostRepository) GetPostsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error) { return []*models.Post{}, nil }
func (r *mockPostRepository) GetPostsByCommunity(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]*models.Post, error) { return []*models.Post{}, nil }
func (r *mockPostRepository) GetSocialFeed(ctx context.Context, userID uuid.UUID, req *models.SocialFeedRequest) (*models.SocialFeed, error) { return &models.SocialFeed{}, nil }
func (r *mockPostRepository) GetTrendingPosts(ctx context.Context, req *models.TrendingRequest) (*models.TrendingContent, error) { return &models.TrendingContent{}, nil }
func (r *mockPostRepository) IncrementViewCount(ctx context.Context, postID uuid.UUID) error { return nil }
func (r *mockPostRepository) UpdateInteractionCounts(ctx context.Context, postID uuid.UUID, likes, comments, shares, reactions int) error { return nil }

type mockCommentRepository struct{}
func (r *mockCommentRepository) GetComment(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) { return &models.Comment{}, nil }
func (r *mockCommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error { return nil }
func (r *mockCommentRepository) UpdateComment(ctx context.Context, commentID uuid.UUID, content string, metadata map[string]interface{}) error { return nil }
func (r *mockCommentRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error { return nil }
func (r *mockCommentRepository) GetCommentsByPost(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*models.Comment, error) { return []*models.Comment{}, nil }
func (r *mockCommentRepository) GetCommentReplies(ctx context.Context, parentID uuid.UUID, limit, offset int) ([]*models.Comment, error) { return []*models.Comment{}, nil }
func (r *mockCommentRepository) UpdateInteractionCounts(ctx context.Context, commentID uuid.UUID, likes, replies, reactions int) error { return nil }

type mockReactionRepository struct{}
func (r *mockReactionRepository) GetReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) (*models.Reaction, error) { return nil, nil }
func (r *mockReactionRepository) CreateReaction(ctx context.Context, reaction *models.Reaction) error { return nil }
func (r *mockReactionRepository) UpdateReaction(ctx context.Context, reactionID uuid.UUID, reactionType string) error { return nil }
func (r *mockReactionRepository) DeleteReaction(ctx context.Context, userID, targetID uuid.UUID, targetType string) error { return nil }
func (r *mockReactionRepository) GetReactionsByTarget(ctx context.Context, targetID uuid.UUID, targetType string) ([]*models.Reaction, error) { return []*models.Reaction{}, nil }
func (r *mockReactionRepository) GetReactionCounts(ctx context.Context, targetID uuid.UUID, targetType string) (map[string]int, error) { return map[string]int{}, nil }

type mockCommunityRepository struct{}
func (r *mockCommunityRepository) GetCommunity(ctx context.Context, communityID uuid.UUID) (*models.Community, error) { return &models.Community{}, nil }
func (r *mockCommunityRepository) CreateCommunity(ctx context.Context, community *models.Community) error { return nil }
func (r *mockCommunityRepository) UpdateCommunity(ctx context.Context, communityID uuid.UUID, updates *models.UpdateCommunityRequest) error { return nil }
func (r *mockCommunityRepository) DeleteCommunity(ctx context.Context, communityID uuid.UUID) error { return nil }
func (r *mockCommunityRepository) GetCommunitiesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Community, error) { return []*models.Community{}, nil }
func (r *mockCommunityRepository) DiscoverCommunities(ctx context.Context, req *models.CommunityDiscoveryRequest) ([]*models.Community, error) { return []*models.Community{}, nil }
func (r *mockCommunityRepository) JoinCommunity(ctx context.Context, member *models.CommunityMember) error { return nil }
func (r *mockCommunityRepository) LeaveCommunity(ctx context.Context, communityID, userID uuid.UUID) error { return nil }
func (r *mockCommunityRepository) GetCommunityMembers(ctx context.Context, communityID uuid.UUID, limit, offset int) ([]*models.CommunityMember, error) { return []*models.CommunityMember{}, nil }
func (r *mockCommunityRepository) UpdateMemberRole(ctx context.Context, communityID, userID uuid.UUID, role string) error { return nil }
func (r *mockCommunityRepository) GetMembershipStatus(ctx context.Context, communityID, userID uuid.UUID) (*models.CommunityMember, error) { return nil, nil }

type mockShareRepository struct{}
func (r *mockShareRepository) CreateShare(ctx context.Context, share *models.Share) error { return nil }
func (r *mockShareRepository) GetSharesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Share, error) { return []*models.Share{}, nil }
func (r *mockShareRepository) GetSharesByContent(ctx context.Context, contentID uuid.UUID, contentType string) ([]*models.Share, error) { return []*models.Share{}, nil }
func (r *mockShareRepository) UpdateShareStatus(ctx context.Context, shareID uuid.UUID, status string) error { return nil }
func (r *mockShareRepository) DeleteShare(ctx context.Context, shareID uuid.UUID) error { return nil }

// MobileIntegrationVerificationTest provides comprehensive verification of mobile integration patterns
type MobileIntegrationVerificationTest struct {
	router       *gin.Engine
	mobileHandler *handlers.MobileHandler
	testUserID   uuid.UUID
	testContext  context.Context
}

// setupMobileIntegrationTest initializes the test environment
func setupMobileIntegrationTest(t *testing.T) *MobileIntegrationVerificationTest {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock services for testing
	mockRepo := &mockRepositoryManager{}
	userService := services.NewUserService(mockRepo)
	syncService := services.NewMobileSyncService(mockRepo)

	// Create mobile handler
	mobileHandler := handlers.NewMobileHandler(syncService, userService)

	// Setup router with mobile routes
	router := gin.New()
	api := router.Group("/api/v1")
	mobileHandler.RegisterMobileRoutes(api)

	return &MobileIntegrationVerificationTest{
		router:        router,
		mobileHandler: mobileHandler,
		testUserID:    uuid.New(),
		testContext:   context.Background(),
	}
}

// TestComprehensiveMobileIntegrationVerification runs complete mobile integration verification
func TestComprehensiveMobileIntegrationVerification(t *testing.T) {
	mvt := setupMobileIntegrationTest(t)

	t.Run("Mobile API Endpoint Verification", func(t *testing.T) {
		mvt.verifyMobileAPIEndpoints(t)
	})

	t.Run("KMP Data Serialization Verification", func(t *testing.T) {
		mvt.verifyKMPDataSerialization(t)
	})

	t.Run("Mobile Sync Pattern Verification", func(t *testing.T) {
		mvt.verifyMobileSyncPatterns(t)
	})

	t.Run("Offline-First Pattern Verification", func(t *testing.T) {
		mvt.verifyOfflineFirstPatterns(t)
	})

	t.Run("Southeast Asian Regional Compliance", func(t *testing.T) {
		mvt.verifySEARegionalCompliance(t)
	})

	t.Run("Mobile Performance Pattern Verification", func(t *testing.T) {
		mvt.verifyMobilePerformancePatterns(t)
	})

	t.Run("Cross-Platform Compatibility Verification", func(t *testing.T) {
		mvt.verifyCrossPlatformCompatibility(t)
	})

	t.Run("Error Handling Pattern Verification", func(t *testing.T) {
		mvt.verifyErrorHandlingPatterns(t)
	})
}

// verifyMobileAPIEndpoints tests all mobile-specific API endpoints
func (mvt *MobileIntegrationVerificationTest) verifyMobileAPIEndpoints(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		requiresBody   bool
		body           interface{}
	}{
		{
			name:           "Mobile Health Check",
			method:         "GET",
			path:           "/api/v1/mobile/health",
			expectedStatus: http.StatusOK,
			requiresBody:   false,
		},
		{
			name:           "Get Mobile Profile",
			method:         "GET",
			path:           fmt.Sprintf("/api/v1/mobile/profile/%s", mvt.testUserID),
			expectedStatus: http.StatusOK,
			requiresBody:   false,
		},
		{
			name:           "Initial User Data Load",
			method:         "GET",
			path:           fmt.Sprintf("/api/v1/mobile/init/%s", mvt.testUserID),
			expectedStatus: http.StatusOK,
			requiresBody:   false,
		},
		{
			name:           "User Feed Request",
			method:         "GET",
			path:           fmt.Sprintf("/api/v1/mobile/feed/%s?limit=20&offset=0", mvt.testUserID),
			expectedStatus: http.StatusOK,
			requiresBody:   false,
		},
		{
			name:           "Discovery Feed Request",
			method:         "GET",
			path:           fmt.Sprintf("/api/v1/mobile/discover/%s?region=TH&limit=10", mvt.testUserID),
			expectedStatus: http.StatusOK,
			requiresBody:   false,
		},
		{
			name:           "Trending Content Request",
			method:         "GET",
			path:           "/api/v1/mobile/trending?region=TH&limit=10",
			expectedStatus: http.StatusOK,
			requiresBody:   false,
		},
		{
			name:           "Profile Sync Changes",
			method:         "GET",
			path:           fmt.Sprintf("/api/v1/mobile/sync/profile/%s?since=%s", mvt.testUserID, url.QueryEscape(time.Now().Add(-time.Hour).Format(time.RFC3339))),
			expectedStatus: http.StatusOK,
			requiresBody:   false,
		},
		{
			name:           "Mobile Follow Request",
			method:         "POST",
			path:           "/api/v1/mobile/follow",
			expectedStatus: http.StatusOK,
			requiresBody:   true,
			body: map[string]interface{}{
				"followerId":  mvt.testUserID,
				"followingId": uuid.New(),
				"source":      "mobile_app",
			},
		},
		{
			name:           "Apply Client Changes",
			method:         "POST",
			path:           fmt.Sprintf("/api/v1/mobile/apply/%s", mvt.testUserID),
			expectedStatus: http.StatusOK,
			requiresBody:   true,
			body: services.ClientChanges{
				UserID:   mvt.testUserID,
				LastSync: time.Now().Add(-time.Hour),
				ProfileChanges: &services.ClientProfileChanges{
					DisplayName: stringPtr("Mobile Updated Name"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tc.requiresBody {
				bodyBytes, _ := json.Marshal(tc.body)
				req, err = http.NewRequest(tc.method, tc.path, bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tc.method, tc.path, nil)
			}

			require.NoError(t, err, "Failed to create request")

			// Execute request
			w := httptest.NewRecorder()
			mvt.router.ServeHTTP(w, req)

			// Verify response
			assert.Equal(t, tc.expectedStatus, w.Code, "Unexpected status code for %s", tc.name)

			// Verify response is valid JSON
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON for %s", tc.name)

			// Verify response has expected structure for successful requests
			if w.Code == http.StatusOK {
				assert.NotEmpty(t, response, "Successful responses should not be empty")
			}
		})
	}
}

// verifyKMPDataSerialization tests KMP-specific data serialization requirements
func (mvt *MobileIntegrationVerificationTest) verifyKMPDataSerialization(t *testing.T) {
	t.Run("Mobile Sync Response KMP Compatibility", func(t *testing.T) {
		syncResponse := &services.MobileSyncResponse{
			UserID:     mvt.testUserID,
			SyncTime:   time.Now(),
			LastSync:   time.Now().Add(-time.Hour),
			HasChanges: true,
			ChangeType: "profile_updated",
			Changes: []services.ChangeItem{
				{
					Type:      "profile",
					Action:    "update",
					EntityID:  mvt.testUserID.String(),
					Timestamp: time.Now(),
					Data:      map[string]interface{}{"displayName": "Updated Name"},
				},
			},
		}

		// Test JSON serialization for KMP
		jsonData, err := json.Marshal(syncResponse)
		require.NoError(t, err, "MobileSyncResponse must serialize for KMP")

		// Verify KMP-critical fields
		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, "\"userId\":", "UserID required for KMP")
		assert.Contains(t, jsonStr, "\"syncTime\":", "SyncTime required for KMP")
		assert.Contains(t, jsonStr, "\"hasChanges\":", "HasChanges required for KMP")
		assert.Contains(t, jsonStr, "\"changes\":", "Changes array required for KMP")

		// Test deserialization
		var deserializedResponse services.MobileSyncResponse
		err = json.Unmarshal(jsonData, &deserializedResponse)
		require.NoError(t, err, "MobileSyncResponse must deserialize for KMP")

		// Verify critical fields survive round-trip
		assert.Equal(t, syncResponse.UserID, deserializedResponse.UserID, "UserID must survive serialization")
		assert.Equal(t, syncResponse.HasChanges, deserializedResponse.HasChanges, "HasChanges must survive")
		assert.Equal(t, len(syncResponse.Changes), len(deserializedResponse.Changes), "Changes count must survive")
	})

	t.Run("Initial Sync Response KMP Compatibility", func(t *testing.T) {
		initialResponse := &services.InitialSyncResponse{
			UserID:   mvt.testUserID,
			SyncTime: time.Now(),
			Profile: &models.SocialProfile{
				FollowersCount: 42,
				FollowingCount: 24,
				PostsCount:     8,
			},
			Followers: []*models.Follow{},
			Following: []*models.Follow{},
			MetaData: services.InitialSyncMetadata{
				TotalFollowers: 42,
				TotalFollowing: 24,
				TotalPosts:     8,
				LastActivity:   time.Now(),
			},
		}

		// Test JSON serialization
		jsonData, err := json.Marshal(initialResponse)
		require.NoError(t, err, "InitialSyncResponse must serialize for KMP")

		// Verify mobile-critical fields
		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, "\"profile\":", "Profile data required for mobile UI")
		assert.Contains(t, jsonStr, "\"followers\":", "Followers list required for mobile")
		assert.Contains(t, jsonStr, "\"following\":", "Following list required for mobile")
		assert.Contains(t, jsonStr, "\"metadata\":", "Metadata required for mobile caching")

		// Test deserialization
		var deserializedResponse services.InitialSyncResponse
		err = json.Unmarshal(jsonData, &deserializedResponse)
		require.NoError(t, err, "InitialSyncResponse must deserialize for KMP")

		assert.Equal(t, initialResponse.UserID, deserializedResponse.UserID, "UserID must survive")
		assert.NotNil(t, deserializedResponse.Profile, "Profile must not be nil after deserialization")
	})
}

// verifyMobileSyncPatterns tests mobile synchronization patterns
func (mvt *MobileIntegrationVerificationTest) verifyMobileSyncPatterns(t *testing.T) {
	t.Run("Incremental Sync Pattern", func(t *testing.T) {
		// Test incremental sync request pattern
		since := time.Now().Add(-time.Hour)
		path := fmt.Sprintf("/api/v1/mobile/sync/profile/%s?since=%s", mvt.testUserID, url.QueryEscape(since.Format(time.RFC3339)))

		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Incremental sync should succeed")

		var response services.MobileSyncResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should deserialize to MobileSyncResponse")

		// Verify sync response structure
		assert.Equal(t, mvt.testUserID, response.UserID, "Response should contain correct user ID")
		assert.False(t, response.SyncTime.IsZero(), "Response should contain sync timestamp")
		assert.NotNil(t, response.Changes, "Changes array should not be nil")
	})

	t.Run("Initial Data Load Pattern", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/mobile/init/%s", mvt.testUserID)

		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Initial data load should succeed")

		var response services.InitialSyncResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should deserialize to InitialSyncResponse")

		// Verify initial data structure for mobile needs
		assert.Equal(t, mvt.testUserID, response.UserID, "Response should contain correct user ID")
		assert.NotNil(t, response.Profile, "Profile should be included in initial data")
		assert.NotNil(t, response.Followers, "Followers should be included")
		assert.NotNil(t, response.Following, "Following should be included")
	})

	t.Run("Feed Pagination Pattern", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/mobile/feed/%s?limit=10&offset=0", mvt.testUserID)

		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Feed request should succeed")

		var response services.FeedSyncResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should deserialize to FeedSyncResponse")

		// Verify pagination structure for mobile
		assert.Equal(t, mvt.testUserID, response.UserID, "Response should contain correct user ID")
		assert.NotNil(t, response.Posts, "Posts array should not be nil")
		assert.GreaterOrEqual(t, response.NextOffset, 0, "NextOffset should be non-negative")
	})
}

// verifyOfflineFirstPatterns tests offline-first mobile patterns
func (mvt *MobileIntegrationVerificationTest) verifyOfflineFirstPatterns(t *testing.T) {
	t.Run("Client Changes Application Pattern", func(t *testing.T) {
		changes := services.ClientChanges{
			UserID:   mvt.testUserID,
			LastSync: time.Now().Add(-time.Hour),
			ProfileChanges: &services.ClientProfileChanges{
				DisplayName: stringPtr("Offline Updated Name"),
				Bio:         stringPtr("Updated offline"),
			},
			FollowChanges: []services.ClientFollowChange{
				{
					Action:       "follow",
					TargetUserID: uuid.New(),
					Timestamp:    time.Now(),
				},
			},
		}

		bodyBytes, err := json.Marshal(changes)
		require.NoError(t, err)

		path := fmt.Sprintf("/api/v1/mobile/apply/%s", mvt.testUserID)
		req, err := http.NewRequest("POST", path, bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Client changes application should succeed")

		var response services.SyncResult
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should deserialize to SyncResult")

		// Verify offline changes were processed
		assert.Equal(t, mvt.testUserID, response.UserID, "Response should contain correct user ID")
		assert.True(t, response.Success, "Changes should be applied successfully")
		assert.GreaterOrEqual(t, response.AppliedChanges, 0, "Applied changes count should be non-negative")
	})

	t.Run("Conflict Resolution Pattern", func(t *testing.T) {
		conflictData := map[string]interface{}{
			"clientData": map[string]interface{}{
				"displayName": "Client Version",
				"bio":         "Client bio",
			},
			"lastSync": time.Now().Add(-time.Hour * 2).Format(time.RFC3339),
		}

		bodyBytes, err := json.Marshal(conflictData)
		require.NoError(t, err)

		path := fmt.Sprintf("/api/v1/mobile/resolve/%s", mvt.testUserID)
		req, err := http.NewRequest("POST", path, bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Conflict resolution should succeed")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should be valid JSON")

		// Verify conflict resolution structure
		assert.Contains(t, response, "userId", "Response should contain user ID")
		assert.Contains(t, response, "resolvedAt", "Response should contain resolution timestamp")
		assert.Contains(t, response, "hasConflicts", "Response should indicate conflict status")
	})
}

// verifySEARegionalCompliance tests Southeast Asian regional compliance
func (mvt *MobileIntegrationVerificationTest) verifySEARegionalCompliance(t *testing.T) {
	supportedRegions := []string{"TH", "SG", "ID", "MY", "PH", "VN"}

	for _, region := range supportedRegions {
		t.Run(fmt.Sprintf("Region %s Compliance", region), func(t *testing.T) {
			// Test discovery feed for region
			path := fmt.Sprintf("/api/v1/mobile/discover/%s?region=%s&limit=5", mvt.testUserID, region)
			req, err := http.NewRequest("GET", path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			mvt.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Discovery should work for region %s", region)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify region is preserved in response
			assert.Equal(t, region, response["region"], "Response should preserve region %s", region)
			assert.Contains(t, response, "profiles", "Response should contain profiles")

			// Test trending content for region
			trendingPath := fmt.Sprintf("/api/v1/mobile/trending?region=%s&limit=5", region)
			trendingReq, err := http.NewRequest("GET", trendingPath, nil)
			require.NoError(t, err)

			trendingW := httptest.NewRecorder()
			mvt.router.ServeHTTP(trendingW, trendingReq)

			assert.Equal(t, http.StatusOK, trendingW.Code, "Trending should work for region %s", region)
		})
	}

	t.Run("Unsupported Region Handling", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/mobile/discover/%s?region=XX&limit=5", mvt.testUserID)
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Unsupported region should return error")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error", "Error response should contain error field")
		assert.Equal(t, "unsupported_region", response["error"], "Should return unsupported_region error")
	})
}

// verifyMobilePerformancePatterns tests mobile-specific performance optimizations
func (mvt *MobileIntegrationVerificationTest) verifyMobilePerformancePatterns(t *testing.T) {
	t.Run("Mobile-Optimized Pagination", func(t *testing.T) {
		// Test that mobile pagination limits are enforced
		path := fmt.Sprintf("/api/v1/mobile/feed/%s?limit=100&offset=0", mvt.testUserID)
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Feed request should succeed")

		// Verify that excessive limits are capped for mobile performance
		var response services.FeedSyncResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			// Response should handle large limit requests gracefully
			assert.LessOrEqual(t, len(response.Posts), 50, "Mobile responses should limit data size")
		}
	})

	t.Run("Discovery Feed Size Limits", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/mobile/discover/%s?region=TH&limit=30", mvt.testUserID)
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Discovery should succeed")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify mobile-appropriate response size
		if profiles, ok := response["profiles"].([]interface{}); ok {
			assert.LessOrEqual(t, len(profiles), 20, "Discovery should limit results for mobile")
		}
	})
}

// verifyCrossPlatformCompatibility tests cross-platform compatibility patterns
func (mvt *MobileIntegrationVerificationTest) verifyCrossPlatformCompatibility(t *testing.T) {
	t.Run("UUID Field Compatibility", func(t *testing.T) {
		// Test that UUID fields are properly serialized for KMP
		followRequest := map[string]interface{}{
			"followerId":  mvt.testUserID.String(),
			"followingId": uuid.New().String(),
			"source":      "kmp_test",
		}

		bodyBytes, err := json.Marshal(followRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/mobile/follow", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Follow request should succeed")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool), "Follow should be successful")
		assert.Contains(t, response, "timestamp", "Response should contain timestamp")
	})

	t.Run("Timestamp Format Compatibility", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/mobile/sync/profile/%s?since=%s", mvt.testUserID, url.QueryEscape(time.Now().Add(-time.Hour).Format(time.RFC3339)))
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Sync request should succeed")

		var response services.MobileSyncResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify timestamps are in KMP-compatible format
		assert.False(t, response.SyncTime.IsZero(), "SyncTime should be set")
		assert.False(t, response.LastSync.IsZero(), "LastSync should be set")
	})
}

// verifyErrorHandlingPatterns tests mobile-specific error handling
func (mvt *MobileIntegrationVerificationTest) verifyErrorHandlingPatterns(t *testing.T) {
	t.Run("Invalid UUID Error Handling", func(t *testing.T) {
		path := "/api/v1/mobile/profile/invalid-uuid"
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid UUID should return bad request")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error", "Error response should contain error field")
		assert.Equal(t, "invalid_user_id", response["error"], "Should return invalid_user_id error")
	})

	t.Run("Invalid Sync Parameter Error Handling", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/mobile/sync/profile/%s?since=invalid-time", mvt.testUserID)
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid time format should return bad request")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error", "Error response should contain error field")
		assert.Contains(t, response, "message", "Error response should contain message field")
	})

	t.Run("Missing Required Parameters", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/mobile/sync/profile/%s", mvt.testUserID) // Missing 'since' parameter
		req, err := http.NewRequest("GET", path, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		mvt.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Missing required parameter should return bad request")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error", "Error response should contain error field")
		assert.Equal(t, "missing_since_parameter", response["error"], "Should return missing_since_parameter error")
	})
}

// Helper functions and mock implementations

func stringPtr(s string) *string {
	return &s
}

// Helper types for testing