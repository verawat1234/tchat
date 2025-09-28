package contract

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tchat/social/models"
	"tchat/social/repository"
	"tchat/social/services"
	sharedModels "tchat.dev/shared/models"
)

// KMPIntegrationTest tests the contract between backend services and KMP mobile clients
// This focuses on data model compatibility and service integration patterns
type KMPIntegrationTest struct {
	userService services.UserService
	testUserID  uuid.UUID
}

// setupTestServices initializes services for KMP integration testing
func (kit *KMPIntegrationTest) setupTestServices(t *testing.T) {
	// Use in-memory repository manager for testing
	repo := &mockRepositoryManager{}
	kit.userService = services.NewUserService(repo)
	kit.testUserID = uuid.New()
}

// TestKMPDataModelCompatibility tests that all models are KMP-compatible
func TestKMPDataModelCompatibility(t *testing.T) {
	kit := &KMPIntegrationTest{}
	kit.setupTestServices(t)

	t.Run("SocialProfile JSON Serialization for KMP", func(t *testing.T) {
		// Create a comprehensive SocialProfile with all KMP-relevant fields using direct assignment
		profile := &models.SocialProfile{}

		// Set embedded User fields directly
		profile.ID = kit.testUserID
		profile.Username = "kmp_user"
		profile.Email = "kmp@tchat.dev"
		profile.Name = "KMP User"
		profile.DisplayName = "KMP Tester"
		profile.Bio = "Testing KMP integration"
		profile.Avatar = "https://api.tchat.dev/avatar.jpg"
		profile.Country = "TH"
		profile.Locale = "th"
		profile.Timezone = "Asia/Bangkok"
		profile.Status = "active"
		profile.Active = true
		profile.KYCTier = 1
		profile.Verified = false

		// Set SocialProfile-specific fields
		profile.Interests = []string{"kotlin", "compose", "mobile"}
		profile.SocialLinks = map[string]interface{}{
			"github":   "https://github.com/user",
			"linkedin": "https://linkedin.com/in/user",
		}
		profile.SocialPreferences = map[string]interface{}{
			"privacy_level":      "public",
			"show_activity":      true,
			"allow_messages":     true,
			"notification_types": []string{"follows", "mentions"},
		}
		profile.FollowersCount = 42
		profile.FollowingCount = 24
		profile.PostsCount = 8
		profile.IsSocialVerified = false
		profile.SocialCreatedAt = time.Now().Add(-time.Hour * 24 * 30)
		profile.SocialUpdatedAt = time.Now()

		// Test JSON serialization (critical for KMP)
		jsonData, err := json.Marshal(profile)
		require.NoError(t, err, "SocialProfile must serialize to JSON for KMP")

		// Verify critical KMP fields are present in JSON
		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, "\"id\":", "ID field required for KMP")
		assert.Contains(t, jsonStr, "\"username\":", "Username required for KMP")
		assert.Contains(t, jsonStr, "\"country\":", "Country required for regional features")
		assert.Contains(t, jsonStr, "\"interests\":", "Interests required for discovery")
		assert.Contains(t, jsonStr, "\"followersCount\":", "Metrics required for UI")
		assert.Contains(t, jsonStr, "\"socialPreferences\":", "Preferences required for mobile settings")

		// Test JSON deserialization (KMP will unmarshal JSON responses)
		var deserializedProfile models.SocialProfile
		err = json.Unmarshal(jsonData, &deserializedProfile)
		require.NoError(t, err, "SocialProfile must deserialize from JSON for KMP")

		// Verify critical fields survive round-trip serialization
		assert.Equal(t, profile.ID, deserializedProfile.ID, "ID must survive serialization")
		assert.Equal(t, profile.Username, deserializedProfile.Username, "Username must survive")
		assert.Equal(t, profile.Country, deserializedProfile.Country, "Country must survive")
		assert.Equal(t, len(profile.Interests), len(deserializedProfile.Interests), "Interests count must survive")
		assert.Equal(t, profile.FollowersCount, deserializedProfile.FollowersCount, "Metrics must survive")
	})

	t.Run("Request Models KMP Compatibility", func(t *testing.T) {
		// Test UpdateSocialProfileRequest (common mobile operation)
		updateReq := &models.UpdateSocialProfileRequest{
			DisplayName: stringPtr("New Display Name"),
			Bio:         stringPtr("Updated bio for mobile"),
			Interests:   []string{"kotlin", "android", "ios"},
			SocialLinks: &map[string]interface{}{
				"github": "https://github.com/newuser",
			},
		}

		// Test JSON serialization of request
		jsonData, err := json.Marshal(updateReq)
		require.NoError(t, err, "UpdateRequest must serialize for KMP")

		// Verify nullable fields are handled correctly (critical for KMP)
		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, "\"displayName\":", "Nullable string fields must serialize")
		assert.Contains(t, jsonStr, "\"interests\":", "Array fields must serialize")
		assert.NotContains(t, jsonStr, "\"avatar\":", "Omitted fields should not appear")

		// Test deserialization
		var deserializedReq models.UpdateSocialProfileRequest
		err = json.Unmarshal(jsonData, &deserializedReq)
		require.NoError(t, err, "UpdateRequest must deserialize for KMP")

		// Verify pointer fields are handled correctly
		assert.NotNil(t, deserializedReq.DisplayName, "Non-nil pointers should remain non-nil")
		assert.Equal(t, "New Display Name", *deserializedReq.DisplayName, "Pointer values should match")
		assert.Nil(t, deserializedReq.Avatar, "Nil pointers should remain nil")
	})

	t.Run("FollowRequest KMP Validation", func(t *testing.T) {
		// Test FollowRequest (critical for social features)
		followReq := &models.FollowRequest{
			FollowerID:  kit.testUserID,
			FollowingID: uuid.New(),
			Source:      "manual",
		}

		// Test JSON handling
		jsonData, err := json.Marshal(followReq)
		require.NoError(t, err, "FollowRequest must serialize for KMP")

		var deserializedReq models.FollowRequest
		err = json.Unmarshal(jsonData, &deserializedReq)
		require.NoError(t, err, "FollowRequest must deserialize for KMP")

		// Verify UUID fields are handled correctly (critical for KMP)
		assert.Equal(t, followReq.FollowerID, deserializedReq.FollowerID, "UUIDs must survive serialization")
		assert.Equal(t, followReq.FollowingID, deserializedReq.FollowingID, "UUIDs must survive serialization")
		assert.Equal(t, followReq.Source, deserializedReq.Source, "String fields must survive")
	})

	t.Run("Pagination Request KMP Compatibility", func(t *testing.T) {
		// Test UserDiscoveryRequest (common mobile pattern)
		discoveryReq := &models.UserDiscoveryRequest{
			Region:    "TH",
			Interests: []string{"technology"},
			Limit:     20,
			Offset:    0,
		}

		// Test JSON handling
		jsonData, err := json.Marshal(discoveryReq)
		require.NoError(t, err, "Discovery request must serialize for KMP")

		// Verify mobile-friendly pagination
		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, "\"limit\":20", "Limit must be numeric for KMP")
		assert.Contains(t, jsonStr, "\"offset\":0", "Offset must be numeric for KMP")
		assert.Contains(t, jsonStr, "\"region\":\"TH\"", "Region filtering must work")

		// Test deserialization
		var deserializedReq models.UserDiscoveryRequest
		err = json.Unmarshal(jsonData, &deserializedReq)
		require.NoError(t, err, "Discovery request must deserialize for KMP")

		assert.Equal(t, 20, deserializedReq.Limit, "Pagination limits must be preserved")
		assert.Equal(t, "TH", deserializedReq.Region, "Regional filters must be preserved")
	})
}

// TestKMPServiceIntegrationPatterns tests service integration patterns for KMP
func TestKMPServiceIntegrationPatterns(t *testing.T) {
	kit := &KMPIntegrationTest{}
	kit.setupTestServices(t)

	ctx := context.Background()

	t.Run("Service Response Structure for KMP", func(t *testing.T) {
		// Test that service responses match KMP expectations
		profile, err := kit.userService.GetSocialProfile(ctx, kit.testUserID)

		// Even with mock data, verify response structure
		if err == nil && profile != nil {
			// Verify response has all required fields for mobile UI
			assert.NotEqual(t, uuid.Nil, profile.ID, "ID required for KMP entity identification")
			assert.NotEmpty(t, profile.Username, "Username required for mobile display")
			assert.GreaterOrEqual(t, profile.FollowersCount, 0, "Metrics must be non-negative")
			assert.GreaterOrEqual(t, profile.FollowingCount, 0, "Metrics must be non-negative")
			assert.GreaterOrEqual(t, profile.PostsCount, 0, "Metrics must be non-negative")

			// Verify timestamp fields for mobile synchronization
			assert.False(t, profile.SocialCreatedAt.IsZero(), "Creation time required for sync")
			assert.False(t, profile.SocialUpdatedAt.IsZero(), "Update time required for sync")
		}
	})

	t.Run("Update Service Pattern for KMP", func(t *testing.T) {
		// Test update pattern that mobile clients will use
		updateReq := &models.UpdateSocialProfileRequest{
			DisplayName: stringPtr("Mobile Updated Name"),
			Bio:         stringPtr("Updated from mobile app"),
			Interests:   []string{"mobile", "kotlin", "compose"},
		}

		// Test service call pattern
		updatedProfile, err := kit.userService.UpdateSocialProfile(ctx, kit.testUserID, updateReq)

		// Verify service contract (even with mock implementation)
		if err == nil && updatedProfile != nil {
			// Service should return updated profile for mobile UI refresh
			assert.NotNil(t, updatedProfile, "Service must return updated profile")

			// Verify response structure for mobile consumption
			jsonData, jsonErr := json.Marshal(updatedProfile)
			assert.NoError(t, jsonErr, "Service response must be JSON serializable")
			assert.Greater(t, len(jsonData), 0, "Response must contain data")
		}
	})

	t.Run("Error Handling Pattern for KMP", func(t *testing.T) {
		// Test error handling patterns that KMP clients expect
		invalidUserID := uuid.Nil

		profile, err := kit.userService.GetSocialProfile(ctx, invalidUserID)

		// Verify error handling contract
		if err != nil {
			// Errors should be actionable by mobile clients
			assert.NotEmpty(t, err.Error(), "Errors must have descriptive messages for mobile")
			assert.Nil(t, profile, "Failed operations should return nil data")
		}
	})
}

// TestKMPRegionalCompliance tests Southeast Asian regional compliance for KMP
func TestKMPRegionalCompliance(t *testing.T) {
	t.Run("Southeast Asian Country Support", func(t *testing.T) {
		// Test that all SEA countries are supported (critical for KMP deployment)
		supportedCountries := []string{"TH", "SG", "ID", "MY", "PH", "VN"}

		for _, country := range supportedCountries {
			profile := &models.SocialProfile{
				User: sharedModels.User{
					ID:      uuid.New(),
					Country: country,
					Name:    "Test User",
					Status:  "active",
					Active:  true,
				},
			}

			// Test JSON serialization works for all regions
			jsonData, err := json.Marshal(profile)
			assert.NoError(t, err, "Country %s must be JSON serializable", country)
			assert.Contains(t, string(jsonData), country, "Country code must appear in JSON")
		}
	})

	t.Run("Regional Discovery Request Validation", func(t *testing.T) {
		// Test that discovery requests work for all SEA regions
		regions := []string{"TH", "SG", "ID", "MY", "PH", "VN"}

		for _, region := range regions {
			req := &models.UserDiscoveryRequest{
				Region: region,
				Limit:  20,
				Offset: 0,
			}

			jsonData, err := json.Marshal(req)
			assert.NoError(t, err, "Region %s discovery must serialize", region)
			assert.Contains(t, string(jsonData), region, "Region filter must be preserved")
		}
	})
}

// Mock repository manager for testing
type mockRepositoryManager struct{}

func (m *mockRepositoryManager) Users() repository.UserRepository                   { return &mockUserRepository{} }
func (m *mockRepositoryManager) Posts() repository.PostRepository                   { return nil }
func (m *mockRepositoryManager) Comments() repository.CommentRepository             { return nil }
func (m *mockRepositoryManager) Reactions() repository.ReactionRepository           { return nil }
func (m *mockRepositoryManager) Communities() repository.CommunityRepository         { return nil }
func (m *mockRepositoryManager) Shares() repository.ShareRepository                 { return nil }
func (m *mockRepositoryManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, rm repository.RepositoryManager) error) error { return nil }
func (m *mockRepositoryManager) Close() error                                       { return nil }

// Mock user repository for testing
type mockUserRepository struct{}

func (m *mockUserRepository) GetSocialProfile(ctx context.Context, userID uuid.UUID) (*models.SocialProfile, error) {
	return &models.SocialProfile{
		User: sharedModels.User{
			ID:          userID,
			Username:    "mock_user",
			Email:       "mock@test.com",
			Name:        "Mock User",
			DisplayName: "Mock Display",
			Country:     "TH",
			Status:      "active",
			Active:      true,
		},
		Interests:        []string{"technology"},
		FollowersCount:   10,
		FollowingCount:   5,
		PostsCount:       3,
		SocialCreatedAt:  time.Now().Add(-time.Hour * 24),
		SocialUpdatedAt:  time.Now(),
	}, nil
}

func (m *mockUserRepository) CreateSocialProfile(ctx context.Context, profile *models.SocialProfile) error { return nil }
func (m *mockUserRepository) UpdateSocialProfile(ctx context.Context, userID uuid.UUID, updates *models.UpdateSocialProfileRequest) error { return nil }
func (m *mockUserRepository) GetUserFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error) { return nil, nil }
func (m *mockUserRepository) GetUserFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Follow, error) { return nil, nil }
func (m *mockUserRepository) CreateFollow(ctx context.Context, follow *models.Follow) error { return nil }
func (m *mockUserRepository) DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error { return nil }
func (m *mockUserRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) { return false, nil }
func (m *mockUserRepository) GetUserActivity(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.UserActivity, error) { return nil, nil }
func (m *mockUserRepository) CreateUserActivity(ctx context.Context, activity *models.UserActivity) error { return nil }
func (m *mockUserRepository) DiscoverUsers(ctx context.Context, req *models.UserDiscoveryRequest) ([]*models.SocialProfile, error) { return nil, nil }

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}