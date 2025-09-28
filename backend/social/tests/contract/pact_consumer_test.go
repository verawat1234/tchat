package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"
	"tchat/social/models"
)

var pact dsl.Pact

func init() {
	// Initialize Pact
	pact = dsl.Pact{
		Consumer: "social-service",
		Provider: "social-api",
		LogDir:   "./logs",
		PactDir:  "./pacts",
	}
}

// getBaseURL constructs the proper base URL with HTTP scheme
func getBaseURL() string {
	if pact.Server != nil {
		return fmt.Sprintf("http://%s:%d", pact.Host, pact.Server.Port)
	}
	return fmt.Sprintf("http://%s", pact.Host)
}

// TestSocialUserProfileContract tests the social user profile endpoints
func TestSocialUserProfileContract(t *testing.T) {
	// Note: Pact setup/teardown handled by TestMain

	userID := uuid.New().String()

	// Test GET /api/v1/social/profiles/{userId}
	t.Run("GET social profile", func(t *testing.T) {
		pact.
			AddInteraction().
			Given("A user with social profile exists").
			UponReceiving("A request to get social profile").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.String(fmt.Sprintf("/api/v1/social/profiles/%s", userID)),
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success":   true,
					"data": map[string]interface{}{
						"id":            dsl.Like(userID),
						"username":      dsl.Like("user_12345678"),
						"email":         dsl.Like("user@example.com"),
						"displayName":   dsl.Like("John Doe"),
						"bio":           dsl.Like("This is my bio"),
						"avatar":        dsl.Like("https://example.com/avatar.jpg"),
						"status":        dsl.Like("active"),
						"role":          dsl.Like("user"),
						"region":        dsl.Like("SEA"),
						"country":       dsl.Like("TH"),
						"locale":        dsl.Like("th-TH"),
						"timezone":      dsl.Like("Asia/Bangkok"),
						"kycTier":       dsl.Like("tier_1"),
						"isVerified":    dsl.Like(false),
						"interests":     dsl.EachLike("technology", 1),
						"socialLinks": map[string]interface{}{
							"twitter":   dsl.Like("https://twitter.com/user"),
							"instagram": dsl.Like("https://instagram.com/user"),
						},
						"socialPreferences": map[string]interface{}{
							"privacy_level":  dsl.Like("public"),
							"show_activity":  dsl.Like(true),
							"allow_messages": dsl.Like(true),
						},
						"followersCount":     dsl.Like(150),
						"followingCount":     dsl.Like(75),
						"postsCount":         dsl.Like(42),
						"isSocialVerified":   dsl.Like(false),
						"createdAt":          dsl.Like("2024-09-22T18:30:00Z"),
						"updatedAt":          dsl.Like("2024-09-22T18:30:00Z"),
						"socialCreatedAt":    dsl.Like("2024-09-22T18:30:00Z"),
						"socialUpdatedAt":    dsl.Like("2024-09-22T18:30:00Z"),
					},
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		// Execute test
		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/social/profiles/%s", getBaseURL(), userID), nil)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})

	// Test PUT /api/v1/social/profiles/{userId}
	t.Run("PUT update social profile", func(t *testing.T) {
		updateRequest := models.UpdateSocialProfileRequest{
			DisplayName: stringPtr("Updated Name"),
			Bio:         stringPtr("Updated bio"),
			Interests:   []string{"technology", "travel"},
		}

		requestBody, _ := json.Marshal(updateRequest)

		pact.
			AddInteraction().
			Given("A user with social profile exists").
			UponReceiving("A request to update social profile").
			WithRequest(dsl.Request{
				Method: "PUT",
				Path:   dsl.String(fmt.Sprintf("/api/v1/social/profiles/%s", userID)),
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
				Body: dsl.Like(map[string]interface{}{
					"displayName": "Updated Name",
					"bio":         "Updated bio",
					"interests":   []string{"technology", "travel"},
				}),
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"id":          dsl.Like(userID),
						"displayName": dsl.Like("Updated Name"),
						"bio":         dsl.Like("Updated bio"),
						"interests":   dsl.EachLike("technology", 1),
						"updatedAt":   dsl.Like("2024-09-22T18:30:00Z"),
					},
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/social/profiles/%s", getBaseURL(), userID), bytes.NewBuffer(requestBody))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})
}

// TestSocialUserDiscoveryContract tests user discovery endpoints
func TestSocialUserDiscoveryContract(t *testing.T) {
	// Note: Pact setup/teardown handled by TestMain

	// Test GET /api/v1/social/discover/users
	t.Run("GET discover users", func(t *testing.T) {
		pact.
			AddInteraction().
			Given("Users exist for discovery").
			UponReceiving("A request to discover users").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.String("/api/v1/social/discover/users"),
				Query: dsl.MapMatcher{
					"region":    dsl.String("SEA"),
					"interests": dsl.String("technology"),
					"limit":     dsl.String("10"),
					"offset":    dsl.String("0"),
				},
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"data": dsl.EachLike(map[string]interface{}{
						"id":               dsl.Like(uuid.New().String()),
						"username":         dsl.Like("discovered_user_1"),
						"displayName":      dsl.Like("Discovered User"),
						"bio":              dsl.Like("Discovered through search"),
						"avatar":           dsl.Like("https://example.com/avatar.jpg"),
						"interests":        dsl.EachLike("technology", 1),
						"followersCount":   dsl.Like(100),
						"followingCount":   dsl.Like(50),
						"postsCount":       dsl.Like(20),
						"isSocialVerified": dsl.Like(false),
					}, 1),
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/social/discover/users?region=SEA&interests=technology&limit=10&offset=0", getBaseURL()), nil)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})
}

// TestSocialFollowingContract tests following/follower endpoints
func TestSocialFollowingContract(t *testing.T) {
	// Note: Pact setup/teardown handled by TestMain

	// Test POST /api/v1/social/follow
	t.Run("POST follow user", func(t *testing.T) {
		followRequest := models.FollowRequest{
			FollowerID:  uuid.New(),
			FollowingID: uuid.New(),
			Source:      "discovery",
		}

		requestBody, _ := json.Marshal(followRequest)

		pact.
			AddInteraction().
			Given("Two users exist").
			UponReceiving("A request to follow a user").
			WithRequest(dsl.Request{
				Method: "POST",
				Path:   dsl.String("/api/v1/social/follow"),
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
				Body: dsl.Like(map[string]interface{}{
					"followerId":  dsl.Like(followRequest.FollowerID.String()),
					"followingId": dsl.Like(followRequest.FollowingID.String()),
					"source":      dsl.Like("discovery"),
				}),
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"message": dsl.Like("Successfully followed user"),
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/social/follow", getBaseURL()), bytes.NewBuffer(requestBody))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})

	// Test GET /api/v1/social/followers/{userId}
	t.Run("GET user followers", func(t *testing.T) {
		userID := uuid.New().String()

		pact.
			AddInteraction().
			Given("A user with followers exists").
			UponReceiving("A request to get user followers").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.String(fmt.Sprintf("/api/v1/social/followers/%s", userID)),
				Query: dsl.MapMatcher{
					"limit":  dsl.String("20"),
					"offset": dsl.String("0"),
				},
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"followers": dsl.EachLike(map[string]interface{}{
							"id":          dsl.Like(uuid.New().String()),
							"username":    dsl.Like("follower_1"),
							"displayName": dsl.Like("Follower User"),
							"avatar":      dsl.Like("https://example.com/avatar.jpg"),
							"followedAt":  dsl.Like("2024-09-22T18:30:00Z"),
							"isVerified":  dsl.Like(false),
							"mutualCount": dsl.Like(5),
						}, 1),
						"total":   dsl.Like(150),
						"limit":   dsl.Like(20),
						"offset":  dsl.Like(0),
						"hasMore": dsl.Like(true),
					},
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/social/followers/%s?limit=20&offset=0", getBaseURL(), userID), nil)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})
}

// TestSocialPostsContract tests post-related endpoints
func TestSocialPostsContract(t *testing.T) {
	// Note: Pact setup/teardown handled by TestMain

	// Test POST /api/v1/social/posts
	t.Run("POST create post", func(t *testing.T) {
		createRequest := models.CreatePostRequest{
			Content:    "This is a test post about technology trends",
			MediaURLs:  []string{"https://example.com/image.jpg"},
			Tags:       []string{"technology", "trends"},
			Visibility: "public",
			Type:       "text",
		}

		requestBody, _ := json.Marshal(createRequest)

		pact.
			AddInteraction().
			Given("A user is authenticated").
			UponReceiving("A request to create a post").
			WithRequest(dsl.Request{
				Method: "POST",
				Path:   dsl.String("/api/v1/social/posts"),
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
				Body: dsl.Like(map[string]interface{}{
					"content":    "This is a test post about technology trends",
					"mediaUrls":  []string{"https://example.com/image.jpg"},
					"tags":       []string{"technology", "trends"},
					"visibility": "public",
					"postType":   "text",
				}),
			}).
			WillRespondWith(dsl.Response{
				Status: 201,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"id":            dsl.Like(uuid.New().String()),
						"authorId":      dsl.Like(uuid.New().String()),
						"content":       dsl.Like("This is a test post about technology trends"),
						"mediaUrls":     dsl.EachLike("https://example.com/image.jpg", 1),
						"tags":          dsl.EachLike("technology", 1),
						"visibility":    dsl.Like("public"),
						"postType":      dsl.Like("text"),
						"status":        dsl.Like("published"),
						"likesCount":    dsl.Like(0),
						"commentsCount": dsl.Like(0),
						"sharesCount":   dsl.Like(0),
						"viewsCount":    dsl.Like(0),
						"isSponsored":   dsl.Like(false),
						"isPinned":      dsl.Like(false),
						"createdAt":     dsl.Like("2024-09-22T18:30:00Z"),
						"updatedAt":     dsl.Like("2024-09-22T18:30:00Z"),
					},
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/social/posts", getBaseURL()), bytes.NewBuffer(requestBody))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 201, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})

	// Test GET /api/v1/social/posts/{postId}
	t.Run("GET post by ID", func(t *testing.T) {
		postID := uuid.New().String()

		pact.
			AddInteraction().
			Given("A post exists").
			UponReceiving("A request to get a post by ID").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.String(fmt.Sprintf("/api/v1/social/posts/%s", postID)),
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"id":            dsl.Like(postID),
						"authorId":      dsl.Like(uuid.New().String()),
						"content":       dsl.Like("This is a sample post content"),
						"mediaUrls":     dsl.EachLike("https://example.com/image.jpg", 1),
						"tags":          dsl.EachLike("technology", 1),
						"visibility":    dsl.Like("public"),
						"postType":      dsl.Like("text"),
						"status":        dsl.Like("published"),
						"likesCount":    dsl.Like(45),
						"commentsCount": dsl.Like(12),
						"sharesCount":   dsl.Like(8),
						"viewsCount":    dsl.Like(342),
						"isSponsored":   dsl.Like(false),
						"isPinned":      dsl.Like(false),
						"createdAt":     dsl.Like("2024-09-22T18:30:00Z"),
						"updatedAt":     dsl.Like("2024-09-22T18:30:00Z"),
					},
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/social/posts/%s", getBaseURL(), postID), nil)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})
}

// TestSocialFeedContract tests social feed endpoints
func TestSocialFeedContract(t *testing.T) {
	// Note: Pact setup/teardown handled by TestMain

	// Test GET /api/v1/social/feed
	t.Run("GET social feed", func(t *testing.T) {
		pact.
			AddInteraction().
			Given("User has content in their feed").
			UponReceiving("A request to get social feed").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.String("/api/v1/social/feed"),
				Query: dsl.MapMatcher{
					"algorithm": dsl.String("personalized"),
					"limit":     dsl.String("20"),
					"region":    dsl.String("SEA"),
				},
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"posts": dsl.EachLike(map[string]interface{}{
							"id":            dsl.Like(uuid.New().String()),
							"authorId":      dsl.Like(uuid.New().String()),
							"content":       dsl.Like("Feed post content"),
							"mediaUrls":     dsl.EachLike("https://example.com/image.jpg", 1),
							"tags":          dsl.EachLike("feed", 1),
							"visibility":    dsl.Like("public"),
							"postType":      dsl.Like("text"),
							"status":        dsl.Like("published"),
							"likesCount":    dsl.Like(10),
							"commentsCount": dsl.Like(2),
							"sharesCount":   dsl.Like(1),
							"viewsCount":    dsl.Like(100),
							"isSponsored":   dsl.Like(false),
							"isPinned":      dsl.Like(false),
							"createdAt":     dsl.Like("2024-09-22T18:30:00Z"),
							"updatedAt":     dsl.Like("2024-09-22T18:30:00Z"),
						}, 1),
						"nextCursor":   dsl.Like("cursor_20"),
						"hasMore":      dsl.Like(true),
						"algorithm":    dsl.Like("personalized"),
						"generatedAt":  dsl.Like("2024-09-22T18:30:00Z"),
						"metadata": map[string]interface{}{
							"totalPosts":    dsl.Like(20),
							"sponsored":     dsl.Like(1),
							"personalScore": dsl.Like(0.85),
							"region":        dsl.Like("SEA"),
						},
					},
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/social/feed?algorithm=personalized&limit=20&region=SEA", getBaseURL()), nil)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})

	// Test GET /api/v1/social/trending
	t.Run("GET trending content", func(t *testing.T) {
		pact.
			AddInteraction().
			Given("Trending content exists").
			UponReceiving("A request to get trending content").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.String("/api/v1/social/trending"),
				Query: dsl.MapMatcher{
					"region":    dsl.String("SEA"),
					"timeframe": dsl.String("24h"),
					"limit":     dsl.String("20"),
				},
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"posts": dsl.EachLike(map[string]interface{}{
							"id":            dsl.Like(uuid.New().String()),
							"authorId":      dsl.Like(uuid.New().String()),
							"content":       dsl.Like("Trending post content"),
							"tags":          dsl.EachLike("trending", 1),
							"visibility":    dsl.Like("public"),
							"postType":      dsl.Like("text"),
							"status":        dsl.Like("published"),
							"likesCount":    dsl.Like(100),
							"commentsCount": dsl.Like(25),
							"sharesCount":   dsl.Like(15),
							"viewsCount":    dsl.Like(1000),
							"isSponsored":   dsl.Like(false),
							"isPinned":      dsl.Like(true),
							"createdAt":     dsl.Like("2024-09-22T18:30:00Z"),
							"updatedAt":     dsl.Like("2024-09-22T18:30:00Z"),
						}, 1),
						"topics":    dsl.EachLike("artificial intelligence", 1),
						"hashtags":  dsl.EachLike(map[string]interface{}{
							"tag":    dsl.Like("#AI2024"),
							"count":  dsl.Like(12500),
							"growth": dsl.Like("+45%"),
						}, 1),
						"region":    dsl.Like("SEA"),
						"timeframe": dsl.Like("24h"),
						"updatedAt": dsl.Like("2024-09-22T18:30:00Z"),
						"metadata": map[string]interface{}{
							"algorithm": dsl.Like("engagement_velocity"),
							"samples":   dsl.Like(200),
						},
					},
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/social/trending?region=SEA&timeframe=24h&limit=20", getBaseURL()), nil)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})
}

// TestSocialReactionsContract tests reaction endpoints
func TestSocialReactionsContract(t *testing.T) {
	// Note: Pact setup/teardown handled by TestMain

	// Test POST /api/v1/social/reactions
	t.Run("POST add reaction", func(t *testing.T) {
		reactionRequest := models.CreateReactionRequest{
			TargetID:     uuid.New(),
			TargetType:   "post",
			Type:         "like",
		}

		requestBody, _ := json.Marshal(reactionRequest)

		pact.
			AddInteraction().
			Given("A post exists").
			UponReceiving("A request to add a reaction").
			WithRequest(dsl.Request{
				Method: "POST",
				Path:   dsl.String("/api/v1/social/reactions"),
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
				Body: dsl.Like(map[string]interface{}{
					"targetId":     dsl.Like(reactionRequest.TargetID.String()),
					"targetType":   dsl.Like("post"),
					"reactionType": dsl.Like("like"),
				}),
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success":   true,
					"message":   dsl.Like("Reaction added successfully"),
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/social/reactions", getBaseURL()), bytes.NewBuffer(requestBody))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})
}

// TestSocialAnalyticsContract tests analytics endpoints
func TestSocialAnalyticsContract(t *testing.T) {
	// Note: Pact setup/teardown handled by TestMain

	// Test GET /api/v1/social/analytics/users/{userId}
	t.Run("GET user analytics", func(t *testing.T) {
		userID := uuid.New().String()

		pact.
			AddInteraction().
			Given("A user with analytics data exists").
			UponReceiving("A request to get user analytics").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.String(fmt.Sprintf("/api/v1/social/analytics/users/%s", userID)),
				Query: dsl.MapMatcher{
					"period": dsl.String("30d"),
				},
				Headers: dsl.MapMatcher{
					"Content-Type":  dsl.String("application/json"),
					"Authorization": dsl.Like("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json; charset=utf-8"),
				},
				Body: dsl.Like(map[string]interface{}{
					"success": true,
					"data": map[string]interface{}{
						"userId": dsl.Like(userID),
						"period": dsl.Like("30d"),
						"followers": map[string]interface{}{
							"current": dsl.Like(150),
							"growth":  dsl.Like("+12"),
							"rate":    dsl.Like("8.7%"),
						},
						"following": map[string]interface{}{
							"current": dsl.Like(75),
							"growth":  dsl.Like("+3"),
							"rate":    dsl.Like("4.2%"),
						},
						"engagement": map[string]interface{}{
							"rate":           dsl.Like("3.2%"),
							"totalLikes":     dsl.Like(245),
							"totalComments":  dsl.Like(68),
							"totalShares":    dsl.Like(23),
							"averagePerPost": dsl.Like(8.4),
						},
						"reach": map[string]interface{}{
							"impressions": dsl.Like(12500),
							"uniqueViews": dsl.Like(8200),
							"countries":   dsl.EachLike("TH", 1),
						},
						"demographics": map[string]interface{}{
							"ageGroups": map[string]interface{}{
								"18-24": dsl.Like(35),
								"25-34": dsl.Like(40),
								"35-44": dsl.Like(20),
								"45+":   dsl.Like(5),
							},
							"topCountries": dsl.EachLike(map[string]interface{}{
								"country":    dsl.Like("TH"),
								"percentage": dsl.Like(45),
							}, 1),
						},
						"updatedAt": dsl.Like("2024-09-22T18:30:00Z"),
					},
					"timestamp": dsl.Like("2024-09-22T18:30:00Z"),
				}),
			})

		test := func() error {
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/social/analytics/users/%s?period=30d", getBaseURL(), userID), nil)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
			return nil
		}

		err := pact.Verify(test)
		assert.NoError(t, err)
	})
}

// stringPtr helper function is defined in kmp_integration_test.go

// Test runner
func TestMain(m *testing.M) {
	// Setup Pact
	pact.Setup(true)

	// Run tests
	code := m.Run()

	// Write Pact files
	pact.WritePact()
	pact.Teardown()

	// Exit with test result code
	os.Exit(code)
}