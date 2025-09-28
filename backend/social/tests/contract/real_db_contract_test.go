package contract

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tchat/social/database"
	"tchat/social/models"
	"tchat/social/repository"
	"tchat/social/services"
	sharedModels "tchat.dev/shared/models"
	"tchat.dev/shared/config"
)

// RealDBContractTest tests API contracts with real database data for KMP mobile integration
type RealDBContractTest struct {
	db           *gorm.DB
	repo         repository.RepositoryManager
	userService  services.UserService
	postService  services.PostService
	testUserID   uuid.UUID
	testUser     *models.SocialProfile
}

// SetupTestSuite initializes the test environment with real database
func (suite *RealDBContractTest) SetupTestSuite(t *testing.T) {
	// Connect to test database
	dsn := "host=localhost user=postgres password=postgres dbname=tchat_social_test port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
		return
	}

	suite.db = db

	// Initialize database schema
	dbConfig := database.NewSocialDatabaseConfig()
	err = dbConfig.Initialize(&config.Config{Debug: true})
	require.NoError(t, err, "Failed to initialize database schema")

	// Setup repository and services
	suite.repo = repository.NewManager(db)
	suite.userService = services.NewUserService(suite.repo)

	// Create test user with real data
	suite.createTestUser(t)
}

// TeardownTestSuite cleans up test data
func (suite *RealDBContractTest) TeardownTestSuite(t *testing.T) {
	if suite.db != nil {
		// Clean up test data
		suite.db.Where("id = ?", suite.testUserID).Delete(&models.SocialProfile{})
		suite.db.Where("follower_id = ? OR following_id = ?", suite.testUserID, suite.testUserID).Delete(&models.Follow{})
	}
}

// createTestUser creates a real test user in the database
func (suite *RealDBContractTest) createTestUser(t *testing.T) {
	suite.testUserID = uuid.New()

	// Create test user with KMP-compatible data
	testUser := &models.SocialProfile{
		User: sharedModels.User{
			ID:          suite.testUserID,
			Username:    "kmp_test_user",
			Email:       "kmp.test@tchat.dev",
			Name:        "KMP Test User",
			DisplayName: "KMP Tester",
			Bio:         "Testing KMP mobile integration",
			Avatar:      "https://api.tchat.dev/avatars/test.jpg",
			Country:     "TH",
			Locale:      "th",
			Timezone:    "Asia/Bangkok",
			Status:      "active",
			Active:      true,
			KYCTier:     1,
			Verified:    false,
		},
		Interests: []string{"technology", "mobile", "kotlin"},
		SocialLinks: map[string]interface{}{
			"github":   "https://github.com/kmp-tester",
			"linkedin": "https://linkedin.com/in/kmp-tester",
		},
		SocialPreferences: map[string]interface{}{
			"privacy_level":      "public",
			"show_activity":      true,
			"allow_messages":     true,
			"notification_types": []string{"follows", "mentions", "comments"},
		},
		FollowersCount:    25,
		FollowingCount:    15,
		PostsCount:        8,
		IsSocialVerified:  false,
		SocialCreatedAt:   time.Now().Add(-time.Hour * 24 * 30),
		SocialUpdatedAt:   time.Now(),
	}

	err := suite.repo.Users().CreateSocialProfile(context.Background(), testUser)
	require.NoError(t, err, "Failed to create test user")

	suite.testUser = testUser
}

// TestKMPUserProfileContract tests user profile contract for KMP mobile clients
func TestKMPUserProfileContract(t *testing.T) {
	suite := &RealDBContractTest{}
	suite.SetupTestSuite(t)
	defer suite.TeardownTestSuite(t)

	ctx := context.Background()

	t.Run("KMP Contract: Get Social Profile", func(t *testing.T) {
		// Test the actual service with real database data
		profile, err := suite.userService.GetSocialProfile(ctx, suite.testUserID)

		// Contract assertions for KMP compatibility
		assert.NoError(t, err, "Should successfully get social profile")
		assert.NotNil(t, profile, "Profile should not be nil")

		// Verify KMP-essential fields are present and properly typed
		assert.Equal(t, suite.testUserID, profile.ID, "User ID should match")
		assert.Equal(t, "kmp_test_user", profile.Username, "Username should match")
		assert.Equal(t, "KMP Test User", profile.Name, "Name should match")
		assert.Equal(t, "TH", profile.Country, "Country should be Southeast Asian")
		assert.True(t, profile.Active, "User should be active")

		// Verify social-specific fields for mobile UI
		assert.NotNil(t, profile.Interests, "Interests should not be nil for mobile")
		assert.Len(t, profile.Interests, 3, "Should have test interests")
		assert.Contains(t, profile.Interests, "kotlin", "Should contain mobile-relevant interests")

		// Verify social metrics for mobile dashboards
		assert.GreaterOrEqual(t, profile.FollowersCount, 0, "Followers count should be non-negative")
		assert.GreaterOrEqual(t, profile.FollowingCount, 0, "Following count should be non-negative")
		assert.GreaterOrEqual(t, profile.PostsCount, 0, "Posts count should be non-negative")

		// Verify timestamps are properly formatted for KMP
		assert.False(t, profile.SocialCreatedAt.IsZero(), "Created timestamp should be set")
		assert.False(t, profile.SocialUpdatedAt.IsZero(), "Updated timestamp should be set")

		// Verify JSON serialization works for KMP
		jsonData, err := json.Marshal(profile)
		assert.NoError(t, err, "Should serialize to JSON for mobile")
		assert.Contains(t, string(jsonData), "socialPreferences", "Should include social preferences")
		assert.Contains(t, string(jsonData), "followersCount", "Should include followers count")
	})

	t.Run("KMP Contract: Update Social Profile", func(t *testing.T) {
		// Test update with KMP-compatible request
		updateReq := &models.UpdateSocialProfileRequest{
			DisplayName: stringPtr("Updated KMP User"),
			Bio:         stringPtr("Updated bio for mobile testing"),
			Interests:   []string{"kotlin", "compose", "mobile"},
		}

		updatedProfile, err := suite.userService.UpdateSocialProfile(ctx, suite.testUserID, updateReq)

		// Contract assertions
		assert.NoError(t, err, "Should successfully update profile")
		assert.NotNil(t, updatedProfile, "Updated profile should not be nil")
		assert.Equal(t, "Updated KMP User", updatedProfile.DisplayName, "Display name should be updated")
		assert.Equal(t, "Updated bio for mobile testing", updatedProfile.Bio, "Bio should be updated")
		assert.Len(t, updatedProfile.Interests, 3, "Should have updated interests")
		assert.Contains(t, updatedProfile.Interests, "compose", "Should contain new interests")
	})
}

// TestKMPFollowContract tests follow/unfollow contracts for KMP mobile clients
func TestKMPFollowContract(t *testing.T) {
	suite := &RealDBContractTest{}
	suite.SetupTestSuite(t)
	defer suite.TeardownTestSuite(t)

	ctx := context.Background()

	// Create a second test user to follow
	targetUserID := uuid.New()
	targetUser := &models.SocialProfile{
		User: sharedModels.User{
			ID:          targetUserID,
			Username:    "kmp_target_user",
			Email:       "target@tchat.dev",
			Name:        "Target User",
			DisplayName: "Target",
			Country:     "SG",
			Locale:      "en",
			Status:      "active",
			Active:      true,
		},
		FollowersCount: 0,
		FollowingCount: 0,
		PostsCount:     0,
	}

	err := suite.repo.Users().CreateSocialProfile(ctx, targetUser)
	require.NoError(t, err, "Failed to create target user")
	defer suite.db.Where("id = ?", targetUserID).Delete(&models.SocialProfile{})

	t.Run("KMP Contract: Follow User", func(t *testing.T) {
		// Test follow functionality with real database
		followReq := &models.FollowRequest{
			FollowerID:  suite.testUserID,
			FollowingID: targetUserID,
			Source:      "manual",
		}

		err := suite.userService.FollowUser(ctx, followReq)
		assert.NoError(t, err, "Should successfully follow user")

		// Verify follow relationship in database
		isFollowing, err := suite.repo.Users().IsFollowing(ctx, suite.testUserID, targetUserID)
		assert.NoError(t, err, "Should check follow status")
		assert.True(t, isFollowing, "Should be following the target user")
	})

	t.Run("KMP Contract: Get Followers", func(t *testing.T) {
		// Test get followers with real data
		followers, err := suite.userService.GetFollowers(ctx, targetUserID, 20, 0)

		assert.NoError(t, err, "Should get followers list")
		assert.NotNil(t, followers, "Followers response should not be nil")

		// Verify response structure for mobile
		followersData, ok := followers["followers"]
		assert.True(t, ok, "Should have followers array")
		assert.NotNil(t, followersData, "Followers data should not be nil")

		limit, ok := followers["limit"]
		assert.True(t, ok, "Should have limit field")
		assert.Equal(t, 20, limit, "Limit should match request")

		hasMore, ok := followers["hasMore"]
		assert.True(t, ok, "Should have hasMore field for pagination")
		assert.NotNil(t, hasMore, "HasMore should not be nil")
	})

	t.Run("KMP Contract: Unfollow User", func(t *testing.T) {
		// Test unfollow functionality
		err := suite.userService.UnfollowUser(ctx, suite.testUserID, targetUserID)
		assert.NoError(t, err, "Should successfully unfollow user")

		// Verify unfollow in database
		isFollowing, err := suite.repo.Users().IsFollowing(ctx, suite.testUserID, targetUserID)
		assert.NoError(t, err, "Should check follow status")
		assert.False(t, isFollowing, "Should not be following after unfollow")
	})
}

// TestKMPDiscoveryContract tests user discovery contracts for KMP mobile clients
func TestKMPDiscoveryContract(t *testing.T) {
	suite := &RealDBContractTest{}
	suite.SetupTestSuite(t)
	defer suite.TeardownTestSuite(t)

	ctx := context.Background()

	t.Run("KMP Contract: Discover Users", func(t *testing.T) {
		// Test user discovery with mobile-friendly parameters
		discoveryReq := &models.UserDiscoveryRequest{
			Region:    "TH",
			Interests: []string{"technology"},
			Limit:     10,
			Offset:    0,
		}

		users, err := suite.userService.DiscoverUsers(ctx, discoveryReq)

		// Contract assertions for mobile discovery
		assert.NoError(t, err, "Should discover users successfully")
		assert.NotNil(t, users, "Users list should not be nil")
		assert.LessOrEqual(t, len(users), 10, "Should respect limit parameter")

		// Verify returned users have required fields for mobile UI
		if len(users) > 0 {
			user := users[0]
			assert.NotEqual(t, uuid.Nil, user.ID, "User should have valid ID")
			assert.NotEmpty(t, user.Username, "User should have username")
			assert.NotEmpty(t, user.DisplayName, "User should have display name")
			assert.Equal(t, "TH", user.Country, "User should match region filter")
			assert.NotNil(t, user.Interests, "User should have interests for filtering")
		}
	})
}

// stringPtr helper function is defined in kmp_integration_test.go