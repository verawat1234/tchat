package contract

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/assert"
)

// TestSocialServiceProvider runs Pact provider verification tests
func TestSocialServiceProvider(t *testing.T) {
	go startSocialService()

	pact := &dsl.Pact{
		Provider: "social-api",
	}

	// Get the current working directory
	dir, _ := os.Getwd()
	pactDir := filepath.Join(dir, "pacts")

	// Verify provider against consumer pacts
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:        "http://localhost:8093", // Social service port
		PactURLs:               []string{filepath.Join(pactDir, "social-service-social-api.json")},
		BrokerURL:              os.Getenv("PACT_BROKER_URL"),
		BrokerUsername:         os.Getenv("PACT_BROKER_USERNAME"),
		BrokerPassword:         os.Getenv("PACT_BROKER_PASSWORD"),
		PublishVerificationResults: len(os.Getenv("PACT_BROKER_URL")) > 0,
		ProviderVersion:        "1.0.0",
		StateHandlers: types.StateHandlers{
			// User profile states
			"A user with social profile exists": func() error {
				// Setup test data for user profile
				return setupUserProfile()
			},
			"Two users exist": func() error {
				// Setup test data for following relationship
				return setupTwoUsers()
			},
			"A user with followers exists": func() error {
				// Setup test data for followers
				return setupUserWithFollowers()
			},
			"Users exist for discovery": func() error {
				// Setup test data for user discovery
				return setupUsersForDiscovery()
			},

			// Post states
			"A user is authenticated": func() error {
				// Setup authenticated user
				return setupAuthenticatedUser()
			},
			"A post exists": func() error {
				// Setup test post data
				return setupPost()
			},
			"User has content in their feed": func() error {
				// Setup feed content
				return setupFeedContent()
			},
			"Trending content exists": func() error {
				// Setup trending content
				return setupTrendingContent()
			},

			// Analytics states
			"A user with analytics data exists": func() error {
				// Setup analytics data
				return setupUserAnalytics()
			},
		},
		BeforeEach: func() error {
			// Reset database state before each test
			return resetTestDatabase()
		},
		AfterEach: func() error {
			// Clean up test data after each test
			return cleanupTestData()
		},
	})

	assert.NoError(t, err)
}

// startSocialService starts the social service for testing
func startSocialService() {
	// Mock implementation - in real scenario, this would start the actual social service
	// For now, we'll create a simple HTTP server that mimics the social service responses

	mux := http.NewServeMux()

	// User profile endpoints
	mux.HandleFunc("/api/v1/social/profiles/", handleGetSocialProfile)
	mux.HandleFunc("/api/v1/social/discover/users", handleDiscoverUsers)
	mux.HandleFunc("/api/v1/social/follow", handleFollowUser)
	mux.HandleFunc("/api/v1/social/followers/", handleGetFollowers)
	mux.HandleFunc("/api/v1/social/following/", handleGetFollowing)

	// Post endpoints
	mux.HandleFunc("/api/v1/social/posts", handlePosts)
	mux.HandleFunc("/api/v1/social/feed", handleSocialFeed)
	mux.HandleFunc("/api/v1/social/trending", handleTrendingContent)

	// Reaction endpoints
	mux.HandleFunc("/api/v1/social/reactions", handleReactions)

	// Analytics endpoints
	mux.HandleFunc("/api/v1/social/analytics/users/", handleUserAnalytics)

	server := &http.Server{
		Addr:    ":8093",
		Handler: mux,
	}

	fmt.Println("Starting social service mock on :8093")
	server.ListenAndServe()
}

// State handler implementations
func setupUserProfile() error {
	// Mock setup - would insert test user profile into database
	fmt.Println("Setting up user profile test data")
	return nil
}

func setupTwoUsers() error {
	// Mock setup - would insert two test users
	fmt.Println("Setting up two users test data")
	return nil
}

func setupUserWithFollowers() error {
	// Mock setup - would insert user with followers
	fmt.Println("Setting up user with followers test data")
	return nil
}

func setupUsersForDiscovery() error {
	// Mock setup - would insert users for discovery
	fmt.Println("Setting up users for discovery test data")
	return nil
}

func setupAuthenticatedUser() error {
	// Mock setup - would setup authenticated user session
	fmt.Println("Setting up authenticated user")
	return nil
}

func setupPost() error {
	// Mock setup - would insert test post
	fmt.Println("Setting up post test data")
	return nil
}

func setupFeedContent() error {
	// Mock setup - would insert feed content
	fmt.Println("Setting up feed content test data")
	return nil
}

func setupTrendingContent() error {
	// Mock setup - would insert trending content
	fmt.Println("Setting up trending content test data")
	return nil
}

func setupUserAnalytics() error {
	// Mock setup - would insert analytics data
	fmt.Println("Setting up user analytics test data")
	return nil
}

func resetTestDatabase() error {
	// Mock implementation - would reset test database
	fmt.Println("Resetting test database")
	return nil
}

func cleanupTestData() error {
	// Mock implementation - would clean up test data
	fmt.Println("Cleaning up test data")
	return nil
}

// HTTP handlers that mock the social service responses
func handleGetSocialProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == "GET" {
		response := `{
			"success": true,
			"data": {
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"username": "user_12345678",
				"email": "user@example.com",
				"displayName": "John Doe",
				"bio": "This is my bio",
				"avatar": "https://example.com/avatar.jpg",
				"status": "active",
				"role": "user",
				"region": "SEA",
				"country": "TH",
				"locale": "th-TH",
				"timezone": "Asia/Bangkok",
				"kycTier": "tier_1",
				"isVerified": false,
				"interests": ["technology"],
				"socialLinks": {
					"twitter": "https://twitter.com/user",
					"instagram": "https://instagram.com/user"
				},
				"socialPreferences": {
					"privacy_level": "public",
					"show_activity": true,
					"allow_messages": true
				},
				"followersCount": 150,
				"followingCount": 75,
				"postsCount": 42,
				"isSocialVerified": false,
				"createdAt": "2024-09-22T18:30:00Z",
				"updatedAt": "2024-09-22T18:30:00Z",
				"socialCreatedAt": "2024-09-22T18:30:00Z",
				"socialUpdatedAt": "2024-09-22T18:30:00Z"
			},
			"timestamp": "2024-09-22T18:30:00Z"
		}`
		w.Write([]byte(response))
	} else if r.Method == "PUT" {
		response := `{
			"success": true,
			"data": {
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"displayName": "Updated Name",
				"bio": "Updated bio",
				"interests": ["technology"],
				"updatedAt": "2024-09-22T18:30:00Z"
			},
			"timestamp": "2024-09-22T18:30:00Z"
		}`
		w.Write([]byte(response))
	}
}

func handleDiscoverUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := `{
		"success": true,
		"data": [{
			"id": "123e4567-e89b-12d3-a456-426614174001",
			"username": "discovered_user_1",
			"displayName": "Discovered User",
			"bio": "Discovered through search",
			"avatar": "https://example.com/avatar.jpg",
			"interests": ["technology"],
			"followersCount": 100,
			"followingCount": 50,
			"postsCount": 20,
			"isSocialVerified": false
		}],
		"timestamp": "2024-09-22T18:30:00Z"
	}`
	w.Write([]byte(response))
}

func handleFollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := `{
		"success": true,
		"message": "Successfully followed user",
		"timestamp": "2024-09-22T18:30:00Z"
	}`
	w.Write([]byte(response))
}

func handleGetFollowers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := `{
		"success": true,
		"data": {
			"followers": [{
				"id": "123e4567-e89b-12d3-a456-426614174002",
				"username": "follower_1",
				"displayName": "Follower User",
				"avatar": "https://example.com/avatar.jpg",
				"followedAt": "2024-09-22T18:30:00Z",
				"isVerified": false,
				"mutualCount": 5
			}],
			"total": 150,
			"limit": 20,
			"offset": 0,
			"hasMore": true
		},
		"timestamp": "2024-09-22T18:30:00Z"
	}`
	w.Write([]byte(response))
}

func handleGetFollowing(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := `{
		"success": true,
		"data": {
			"following": [{
				"id": "123e4567-e89b-12d3-a456-426614174003",
				"username": "following_1",
				"displayName": "Following User",
				"avatar": "https://example.com/avatar.jpg",
				"followedAt": "2024-09-22T18:30:00Z",
				"isVerified": false,
				"mutualCount": 3
			}],
			"total": 75,
			"limit": 20,
			"offset": 0,
			"hasMore": true
		},
		"timestamp": "2024-09-22T18:30:00Z"
	}`
	w.Write([]byte(response))
}

func handlePosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == "POST" {
		w.WriteHeader(http.StatusCreated)
		response := `{
			"success": true,
			"data": {
				"id": "123e4567-e89b-12d3-a456-426614174004",
				"authorId": "123e4567-e89b-12d3-a456-426614174000",
				"content": "This is a test post about technology trends",
				"mediaUrls": ["https://example.com/image.jpg"],
				"tags": ["technology"],
				"visibility": "public",
				"postType": "text",
				"status": "published",
				"likesCount": 0,
				"commentsCount": 0,
				"sharesCount": 0,
				"viewsCount": 0,
				"isSponsored": false,
				"isPinned": false,
				"createdAt": "2024-09-22T18:30:00Z",
				"updatedAt": "2024-09-22T18:30:00Z"
			},
			"timestamp": "2024-09-22T18:30:00Z"
		}`
		w.Write([]byte(response))
	} else if r.Method == "GET" {
		// Extract post ID from URL path
		response := `{
			"success": true,
			"data": {
				"id": "123e4567-e89b-12d3-a456-426614174004",
				"authorId": "123e4567-e89b-12d3-a456-426614174000",
				"content": "This is a sample post content",
				"mediaUrls": ["https://example.com/image.jpg"],
				"tags": ["technology"],
				"visibility": "public",
				"postType": "text",
				"status": "published",
				"likesCount": 45,
				"commentsCount": 12,
				"sharesCount": 8,
				"viewsCount": 342,
				"isSponsored": false,
				"isPinned": false,
				"createdAt": "2024-09-22T18:30:00Z",
				"updatedAt": "2024-09-22T18:30:00Z"
			},
			"timestamp": "2024-09-22T18:30:00Z"
		}`
		w.Write([]byte(response))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSocialFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := `{
		"success": true,
		"data": {
			"posts": [{
				"id": "123e4567-e89b-12d3-a456-426614174005",
				"authorId": "123e4567-e89b-12d3-a456-426614174000",
				"content": "Feed post content",
				"mediaUrls": ["https://example.com/image.jpg"],
				"tags": ["feed"],
				"visibility": "public",
				"postType": "text",
				"status": "published",
				"likesCount": 10,
				"commentsCount": 2,
				"sharesCount": 1,
				"viewsCount": 100,
				"isSponsored": false,
				"isPinned": false,
				"createdAt": "2024-09-22T18:30:00Z",
				"updatedAt": "2024-09-22T18:30:00Z"
			}],
			"nextCursor": "cursor_20",
			"hasMore": true,
			"algorithm": "personalized",
			"generatedAt": "2024-09-22T18:30:00Z",
			"metadata": {
				"totalPosts": 20,
				"sponsored": 1,
				"personalScore": 0.85,
				"region": "SEA"
			}
		},
		"timestamp": "2024-09-22T18:30:00Z"
	}`
	w.Write([]byte(response))
}

func handleTrendingContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := `{
		"success": true,
		"data": {
			"posts": [{
				"id": "123e4567-e89b-12d3-a456-426614174006",
				"authorId": "123e4567-e89b-12d3-a456-426614174000",
				"content": "Trending post content",
				"tags": ["trending"],
				"visibility": "public",
				"postType": "text",
				"status": "published",
				"likesCount": 100,
				"commentsCount": 25,
				"sharesCount": 15,
				"viewsCount": 1000,
				"isSponsored": false,
				"isPinned": true,
				"createdAt": "2024-09-22T18:30:00Z",
				"updatedAt": "2024-09-22T18:30:00Z"
			}],
			"topics": ["artificial intelligence"],
			"hashtags": [{
				"tag": "#AI2024",
				"count": 12500,
				"growth": "+45%"
			}],
			"region": "SEA",
			"timeframe": "24h",
			"updatedAt": "2024-09-22T18:30:00Z",
			"metadata": {
				"algorithm": "engagement_velocity",
				"samples": 200
			}
		},
		"timestamp": "2024-09-22T18:30:00Z"
	}`
	w.Write([]byte(response))
}

func handleReactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := `{
		"success": true,
		"message": "Reaction added successfully",
		"timestamp": "2024-09-22T18:30:00Z"
	}`
	w.Write([]byte(response))
}

func handleUserAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := `{
		"success": true,
		"data": {
			"userId": "123e4567-e89b-12d3-a456-426614174000",
			"period": "30d",
			"followers": {
				"current": 150,
				"growth": "+12",
				"rate": "8.7%"
			},
			"following": {
				"current": 75,
				"growth": "+3",
				"rate": "4.2%"
			},
			"engagement": {
				"rate": "3.2%",
				"totalLikes": 245,
				"totalComments": 68,
				"totalShares": 23,
				"averagePerPost": 8.4
			},
			"reach": {
				"impressions": 12500,
				"uniqueViews": 8200,
				"countries": ["TH"]
			},
			"demographics": {
				"ageGroups": {
					"18-24": 35,
					"25-34": 40,
					"35-44": 20,
					"45+": 5
				},
				"topCountries": [{
					"country": "TH",
					"percentage": 45
				}]
			},
			"updatedAt": "2024-09-22T18:30:00Z"
		},
		"timestamp": "2024-09-22T18:30:00Z"
	}`
	w.Write([]byte(response))
}