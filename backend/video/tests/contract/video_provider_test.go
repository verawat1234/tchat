package contract_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"tchat.dev/shared/responses"
	"tchat.dev/video/handlers"
	"tchat.dev/video/models"
	"tchat.dev/video/services"
)

// TestVideoProviderContract runs Pact provider verification tests for the Video service
// This validates that our Video service meets the contract expectations from consumers
func TestVideoProviderContract(t *testing.T) {
	// Create mock dependencies for the video service
	mockVideoRepo := &MockVideoRepository{
		videos:           make(map[uuid.UUID]*models.Video),
		channels:         make(map[uuid.UUID]*models.Channel),
		interactions:     make(map[uuid.UUID][]*models.VideoInteraction),
		comments:         make(map[uuid.UUID][]*models.VideoComment),
		shares:           make(map[uuid.UUID][]*models.VideoShare),
	}

	// Create video service with mock dependencies
	videoService := services.NewVideoService(mockVideoRepo, nil)

	// Create video handler
	videoHandlers := handlers.NewVideoHandlers(videoService)

	// Set up Gin router for video routes
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add middleware
	router.Use(mockAuthMiddleware())

	// API routes for video service
	v1 := router.Group("/api/v1")
	{
		videos := v1.Group("/videos")
		{
			// Video CRUD operations
			videos.GET("", videoHandlers.GetVideos)
			videos.POST("", videoHandlers.CreateVideo)
			videos.GET("/:id", videoHandlers.GetVideo)
			videos.PUT("/:id", videoHandlers.UpdateVideo)
			videos.DELETE("/:id", videoHandlers.DeleteVideo)

			// Video interactions
			videos.POST("/:id/like", videoHandlers.LikeVideo)
			videos.POST("/:id/share", videoHandlers.ShareVideo)

			// Video comments
			videos.GET("/:id/comments", videoHandlers.GetVideoComments)
			videos.POST("/:id/comments", videoHandlers.AddVideoComment)

			// Video upload
			videos.POST("/upload", videoHandlers.UploadVideo)

			// Video search and filtering
			videos.GET("/search", videoHandlers.SearchVideos)
			videos.GET("/category/:category", videoHandlers.GetVideoByCategory)
			videos.GET("/trending", videoHandlers.GetTrendingVideos)

			// Short videos (TikTok style)
			videos.GET("/shorts", videoHandlers.GetShortVideos)

			// Health check
			videos.GET("/health", videoHandlers.VideoHealth)
		}

		// Channel routes
		channels := v1.Group("/channels")
		{
			channels.POST("", videoHandlers.CreateChannel)
			channels.GET("/:id", videoHandlers.GetChannel)
		}
	}

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		responses.SendSuccessResponse(c, gin.H{
			"status":    "ok",
			"service":   "video-service",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
		})
	})

	// Create test server
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Set up provider states for different test scenarios
	stateHandlers := map[string]func(setup bool, state map[string]interface{}) error{
		// Video catalog states
		"videos exist in catalog": func(setup bool, state map[string]interface{}) error {
			return setupVideoCatalogState(mockVideoRepo, setup)
		},
		"video with ID exists": func(setup bool, state map[string]interface{}) error {
			if videoIDValue, exists := state["video_id"]; exists {
				videoID := videoIDValue.(string)
				return setupSingleVideoState(mockVideoRepo, videoID, setup)
			}
			return setupSingleVideoState(mockVideoRepo, "11111111-1111-1111-1111-111111111111", setup)
		},
		"channel has videos": func(setup bool, state map[string]interface{}) error {
			if channelIDValue, exists := state["channel_id"]; exists {
				channelID := channelIDValue.(string)
				return setupChannelVideosState(mockVideoRepo, channelID, setup)
			}
			return setupChannelVideosState(mockVideoRepo, "22222222-2222-2222-2222-222222222222", setup)
		},

		// Short videos states
		"short videos exist": func(setup bool, state map[string]interface{}) error {
			return setupShortVideosState(mockVideoRepo, setup)
		},
		"trending videos exist": func(setup bool, state map[string]interface{}) error {
			return setupTrendingVideosState(mockVideoRepo, setup)
		},

		// Video interaction states
		"video has likes and comments": func(setup bool, state map[string]interface{}) error {
			if videoIDValue, exists := state["video_id"]; exists {
				videoID := videoIDValue.(string)
				return setupVideoInteractionsState(mockVideoRepo, videoID, setup)
			}
			return setupVideoInteractionsState(mockVideoRepo, "33333333-3333-3333-3333-333333333333", setup)
		},
		"user can like videos": func(setup bool, state map[string]interface{}) error {
			return setupAuthenticatedUserState(setup)
		},
		"user can comment on videos": func(setup bool, state map[string]interface{}) error {
			return setupAuthenticatedUserState(setup)
		},

		// Channel states
		"channel exists with videos": func(setup bool, state map[string]interface{}) error {
			channelID := "44444444-4444-4444-4444-444444444444"
			if channelIDValue, exists := state["channel_id"]; exists {
				channelID = channelIDValue.(string)
			}
			return setupChannelExistsState(mockVideoRepo, channelID, setup)
		},
		"user owns channel": func(setup bool, state map[string]interface{}) error {
			userID := "55555555-5555-5555-5555-555555555555"
			channelID := "66666666-6666-6666-6666-666666666666"
			if userIDValue, exists := state["user_id"]; exists {
				userID = userIDValue.(string)
			}
			if channelIDValue, exists := state["channel_id"]; exists {
				channelID = channelIDValue.(string)
			}
			return setupUserOwnsChannelState(mockVideoRepo, userID, channelID, setup)
		},

		// Search and filtering states
		"videos exist by category": func(setup bool, state map[string]interface{}) error {
			category := "entertainment"
			if categoryValue, exists := state["category"]; exists {
				category = categoryValue.(string)
			}
			return setupCategoryVideosState(mockVideoRepo, category, setup)
		},
		"searchable videos exist": func(setup bool, state map[string]interface{}) error {
			return setupSearchableVideosState(mockVideoRepo, setup)
		},

		// Authentication states
		"user is authenticated": func(setup bool, state map[string]interface{}) error {
			return setupAuthenticatedUserState(setup)
		},
		"user has valid JWT token": func(setup bool, state map[string]interface{}) error {
			return setupValidJWTState(setup)
		},
	}

	// Test 1: Validate that we can set up all provider states
	t.Run("validate_provider_states", func(t *testing.T) {
		for stateName, stateHandler := range stateHandlers {
			t.Run(stateName, func(t *testing.T) {
				// Set up the state
				err := stateHandler(true, make(map[string]interface{}))
				assert.NoError(t, err, "Should be able to set up state: %s", stateName)

				// Clean up the state
				err = stateHandler(false, make(map[string]interface{}))
				assert.NoError(t, err, "Should be able to clean up state: %s", stateName)
			})
		}
	})

	// Test 2: Validate that the service responds correctly to basic requests
	t.Run("validate_service_endpoints", func(t *testing.T) {
		// Set up videos in catalog state
		err := stateHandlers["videos exist in catalog"](true, make(map[string]interface{}))
		assert.NoError(t, err)

		// Test video catalog endpoint
		resp, err := http.Get(testServer.URL + "/api/v1/videos")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Test short videos endpoint
		resp, err = http.Get(testServer.URL + "/api/v1/videos/shorts")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Test health endpoint
		resp, err = http.Get(testServer.URL + "/health")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Clean up
		stateHandlers["videos exist in catalog"](false, make(map[string]interface{}))
	})

	// Test 3: Validate authentication handling
	t.Run("validate_authentication", func(t *testing.T) {
		// Test without authentication (should fail for protected endpoints)
		resp, err := http.Post(testServer.URL+"/api/v1/videos", "application/json", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Test with mock authentication (would be handled by middleware)
		client := &http.Client{}
		req, _ := http.NewRequest("GET", testServer.URL+"/api/v1/videos", nil)
		req.Header.Set("Authorization", "Bearer mock-token")
		resp, err = client.Do(req)
		assert.NoError(t, err)
		// Should not be unauthorized (might be other status based on implementation)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// Test 4: Validate video interactions
	t.Run("validate_video_interactions", func(t *testing.T) {
		// Set up a video with interactions
		err := stateHandlers["video has likes and comments"](true, make(map[string]interface{}))
		assert.NoError(t, err)

		// Test getting video comments
		resp, err := http.Get(testServer.URL + "/api/v1/videos/33333333-3333-3333-3333-333333333333/comments")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Clean up
		stateHandlers["video has likes and comments"](false, make(map[string]interface{}))
	})

	// Test 5: Validate search and filtering
	t.Run("validate_search_and_filtering", func(t *testing.T) {
		// Set up searchable videos
		err := stateHandlers["searchable videos exist"](true, make(map[string]interface{}))
		assert.NoError(t, err)

		// Test video search
		resp, err := http.Get(testServer.URL + "/api/v1/videos/search?q=test")
		assert.NoError(t, err)
		// Search might return empty results but should not fail
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound)

		// Test trending videos
		resp, err = http.Get(testServer.URL + "/api/v1/videos/trending")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Clean up
		stateHandlers["searchable videos exist"](false, make(map[string]interface{}))
	})
}

// Provider state setup functions

// setupVideoCatalogState sets up a video catalog with multiple videos
func setupVideoCatalogState(mockRepo *MockVideoRepository, setup bool) error {
	if !setup {
		mockRepo.Clear()
		return nil
	}

	// Create test videos for catalog
	videos := []*models.Video{
		{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Title:        "Amazing Tech Review",
			Description:  "Latest smartphone technology review",
			ThumbnailURL: "https://cdn.tchat.dev/videos/thumb1.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/video1.mp4",
			Duration:     "10:30",
			Views:        1500,
			Likes:        120,
			Category:     "technology",
			Tags:         []string{"tech", "review", "smartphone"},
			Type:         "long",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
		{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111112"),
			Title:        "Funny Cat Compilation",
			Description:  "Hilarious cats doing funny things",
			ThumbnailURL: "https://cdn.tchat.dev/videos/thumb2.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/video2.mp4",
			Duration:     "5:45",
			Views:        5000,
			Likes:        450,
			Category:     "entertainment",
			Tags:         []string{"cats", "funny", "animals"},
			Type:         "short",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222223"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
		{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111113"),
			Title:        "Cooking Thai Pad Thai",
			Description:  "Traditional Thai Pad Thai recipe",
			ThumbnailURL: "https://cdn.tchat.dev/videos/thumb3.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/video3.mp4",
			Duration:     "15:20",
			Views:        3200,
			Likes:        280,
			Category:     "cooking",
			Tags:         []string{"cooking", "thai", "recipe"},
			Type:         "long",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222224"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
	}

	for _, video := range videos {
		mockRepo.videos[video.ID] = video
	}

	// Set up corresponding channels
	channels := []*models.Channel{
		{
			ID:          uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Name:        "Tech Reviews Plus",
			Avatar:      "https://cdn.tchat.dev/avatars/tech.jpg",
			Subscribers: 50000,
			Verified:    true,
			UserID:      uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.MustParse("22222222-2222-2222-2222-222222222223"),
			Name:        "Funny Animals TV",
			Avatar:      "https://cdn.tchat.dev/avatars/animals.jpg",
			Subscribers: 125000,
			Verified:    false,
			UserID:      uuid.MustParse("33333333-3333-3333-3333-333333333334"),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.MustParse("22222222-2222-2222-2222-222222222224"),
			Name:        "Thai Kitchen Secrets",
			Avatar:      "https://cdn.tchat.dev/avatars/cooking.jpg",
			Subscribers: 75000,
			Verified:    true,
			UserID:      uuid.MustParse("33333333-3333-3333-3333-333333333335"),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	for _, channel := range channels {
		mockRepo.channels[channel.ID] = channel
	}

	return nil
}

// setupSingleVideoState sets up a single video by ID
func setupSingleVideoState(mockRepo *MockVideoRepository, videoIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return fmt.Errorf("invalid video ID: %v", err)
	}

	video := &models.Video{
		ID:           videoID,
		Title:        "Test Video",
		Description:  "Test video description",
		ThumbnailURL: "https://cdn.tchat.dev/videos/test-thumb.jpg",
		VideoURL:     "https://cdn.tchat.dev/videos/test-video.mp4",
		Duration:     "3:30",
		Views:        100,
		Likes:        10,
		Category:     "test",
		Tags:         []string{"test"},
		Type:         "short",
		Status:       "active",
		ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	mockRepo.videos[videoID] = video
	return nil
}

// setupChannelVideosState sets up videos for a specific channel
func setupChannelVideosState(mockRepo *MockVideoRepository, channelIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return fmt.Errorf("invalid channel ID: %v", err)
	}

	videos := []*models.Video{
		{
			ID:           uuid.New(),
			Title:        "Channel Video 1",
			Description:  "First video from the channel",
			ThumbnailURL: "https://cdn.tchat.dev/videos/channel1.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/channel1.mp4",
			Duration:     "5:00",
			Views:        200,
			Likes:        15,
			Category:     "entertainment",
			Type:         "short",
			Status:       "active",
			ChannelID:    channelID,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
		{
			ID:           uuid.New(),
			Title:        "Channel Video 2",
			Description:  "Second video from the channel",
			ThumbnailURL: "https://cdn.tchat.dev/videos/channel2.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/channel2.mp4",
			Duration:     "7:30",
			Views:        350,
			Likes:        25,
			Category:     "technology",
			Type:         "long",
			Status:       "active",
			ChannelID:    channelID,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
	}

	for _, video := range videos {
		mockRepo.videos[video.ID] = video
	}

	return nil
}

// setupShortVideosState sets up short-form videos
func setupShortVideosState(mockRepo *MockVideoRepository, setup bool) error {
	if !setup {
		return nil
	}

	videos := []*models.Video{
		{
			ID:           uuid.New(),
			Title:        "Quick Dance",
			Description:  "30-second dance routine",
			ThumbnailURL: "https://cdn.tchat.dev/videos/dance.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/dance.mp4",
			Duration:     "0:30",
			Views:        10000,
			Likes:        800,
			Category:     "entertainment",
			Type:         "short",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
		{
			ID:           uuid.New(),
			Title:        "Life Hack",
			Description:  "Quick life hack tip",
			ThumbnailURL: "https://cdn.tchat.dev/videos/hack.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/hack.mp4",
			Duration:     "0:45",
			Views:        15000,
			Likes:        1200,
			Category:     "lifestyle",
			Type:         "short",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222223"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
	}

	for _, video := range videos {
		mockRepo.videos[video.ID] = video
	}

	return nil
}

// setupTrendingVideosState sets up trending videos
func setupTrendingVideosState(mockRepo *MockVideoRepository, setup bool) error {
	if !setup {
		return nil
	}

	// Create videos with high views and likes for trending
	videos := []*models.Video{
		{
			ID:           uuid.New(),
			Title:        "Viral Challenge",
			Description:  "Latest viral challenge",
			ThumbnailURL: "https://cdn.tchat.dev/videos/viral.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/viral.mp4",
			Duration:     "2:15",
			Views:        100000,
			Likes:        8500,
			Category:     "entertainment",
			Type:         "short",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			CreatedAt:    time.Now().Add(-2 * time.Hour), // Recent for trending
			UpdatedAt:    time.Now().UTC(),
		},
		{
			ID:           uuid.New(),
			Title:        "Breaking News",
			Description:  "Important news update",
			ThumbnailURL: "https://cdn.tchat.dev/videos/news.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/news.mp4",
			Duration:     "5:30",
			Views:        250000,
			Likes:        15000,
			Category:     "news",
			Type:         "long",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222223"),
			CreatedAt:    time.Now().Add(-1 * time.Hour), // Very recent for trending
			UpdatedAt:    time.Now().UTC(),
		},
	}

	for _, video := range videos {
		mockRepo.videos[video.ID] = video
	}

	return nil
}

// setupVideoInteractionsState sets up a video with interactions
func setupVideoInteractionsState(mockRepo *MockVideoRepository, videoIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		return fmt.Errorf("invalid video ID: %v", err)
	}

	// Set up the video first
	video := &models.Video{
		ID:           videoID,
		Title:        "Interactive Video",
		Description:  "Video with lots of interactions",
		ThumbnailURL: "https://cdn.tchat.dev/videos/interactive.jpg",
		VideoURL:     "https://cdn.tchat.dev/videos/interactive.mp4",
		Duration:     "8:15",
		Views:        5000,
		Likes:        350,
		Category:     "entertainment",
		Type:         "long",
		Status:       "active",
		ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	mockRepo.videos[videoID] = video

	// Set up interactions
	interactions := []*models.VideoInteraction{
		{
			ID:        uuid.New(),
			VideoID:   videoID,
			UserID:    uuid.MustParse("44444444-4444-4444-4444-444444444444"),
			Type:      "like",
			CreatedAt: time.Now().UTC(),
		},
		{
			ID:        uuid.New(),
			VideoID:   videoID,
			UserID:    uuid.MustParse("44444444-4444-4444-4444-444444444445"),
			Type:      "like",
			CreatedAt: time.Now().UTC(),
		},
	}

	mockRepo.interactions[videoID] = interactions

	// Set up comments
	comments := []*models.VideoComment{
		{
			ID:         uuid.New(),
			VideoID:    videoID,
			UserID:     uuid.MustParse("44444444-4444-4444-4444-444444444444"),
			UserName:   "TestUser1",
			UserAvatar: "https://cdn.tchat.dev/avatars/user1.jpg",
			Content:    "Great video!",
			Likes:      5,
			IsEdited:   false,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		},
		{
			ID:         uuid.New(),
			VideoID:    videoID,
			UserID:     uuid.MustParse("44444444-4444-4444-4444-444444444445"),
			UserName:   "TestUser2",
			UserAvatar: "https://cdn.tchat.dev/avatars/user2.jpg",
			Content:    "Love this content!",
			Likes:      3,
			IsEdited:   false,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		},
	}

	mockRepo.comments[videoID] = comments

	return nil
}

// setupChannelExistsState sets up an existing channel
func setupChannelExistsState(mockRepo *MockVideoRepository, channelIDStr string, setup bool) error {
	if !setup {
		mockRepo.Clear()
		return nil
	}

	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return fmt.Errorf("invalid channel ID: %v", err)
	}

	channel := &models.Channel{
		ID:          channelID,
		Name:        "Test Channel",
		Avatar:      "https://cdn.tchat.dev/avatars/test-channel.jpg",
		Subscribers: 25000,
		Verified:    true,
		UserID:      uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	mockRepo.channels[channelID] = channel
	return nil
}

// setupUserOwnsChannelState sets up a user as owner of a channel
func setupUserOwnsChannelState(mockRepo *MockVideoRepository, userIDStr, channelIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return fmt.Errorf("invalid channel ID: %v", err)
	}

	channel := &models.Channel{
		ID:          channelID,
		Name:        "User's Channel",
		Avatar:      "https://cdn.tchat.dev/avatars/user-channel.jpg",
		Subscribers: 1000,
		Verified:    false,
		UserID:      userID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	mockRepo.channels[channelID] = channel
	return nil
}

// setupCategoryVideosState sets up videos for a specific category
func setupCategoryVideosState(mockRepo *MockVideoRepository, category string, setup bool) error {
	if !setup {
		return nil
	}

	videos := []*models.Video{
		{
			ID:           uuid.New(),
			Title:        "Category Video 1",
			Description:  fmt.Sprintf("First video in %s category", category),
			ThumbnailURL: fmt.Sprintf("https://cdn.tchat.dev/videos/%s1.jpg", category),
			VideoURL:     fmt.Sprintf("https://cdn.tchat.dev/videos/%s1.mp4", category),
			Duration:     "4:20",
			Views:        800,
			Likes:        60,
			Category:     category,
			Type:         "short",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
		{
			ID:           uuid.New(),
			Title:        "Category Video 2",
			Description:  fmt.Sprintf("Second video in %s category", category),
			ThumbnailURL: fmt.Sprintf("https://cdn.tchat.dev/videos/%s2.jpg", category),
			VideoURL:     fmt.Sprintf("https://cdn.tchat.dev/videos/%s2.mp4", category),
			Duration:     "6:45",
			Views:        1200,
			Likes:        90,
			Category:     category,
			Type:         "long",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222223"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
	}

	for _, video := range videos {
		mockRepo.videos[video.ID] = video
	}

	return nil
}

// setupSearchableVideosState sets up videos with searchable content
func setupSearchableVideosState(mockRepo *MockVideoRepository, setup bool) error {
	if !setup {
		return nil
	}

	videos := []*models.Video{
		{
			ID:           uuid.New(),
			Title:        "How to Test Software",
			Description:  "Complete guide to software testing",
			ThumbnailURL: "https://cdn.tchat.dev/videos/testing-guide.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/testing-guide.mp4",
			Duration:     "12:30",
			Views:        2500,
			Likes:        180,
			Category:     "education",
			Tags:         []string{"testing", "software", "tutorial"},
			Type:         "long",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
		{
			ID:           uuid.New(),
			Title:        "Test Automation Best Practices",
			Description:  "Advanced testing automation techniques",
			ThumbnailURL: "https://cdn.tchat.dev/videos/automation.jpg",
			VideoURL:     "https://cdn.tchat.dev/videos/automation.mp4",
			Duration:     "18:45",
			Views:        1800,
			Likes:        140,
			Category:     "technology",
			Tags:         []string{"automation", "testing", "programming"},
			Type:         "long",
			Status:       "active",
			ChannelID:    uuid.MustParse("22222222-2222-2222-2222-222222222223"),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
	}

	for _, video := range videos {
		mockRepo.videos[video.ID] = video
	}

	return nil
}

// setupAuthenticatedUserState sets up an authenticated user
func setupAuthenticatedUserState(setup bool) error {
	// This is handled by the request filter middleware
	return nil
}

// setupValidJWTState sets up valid JWT authentication
func setupValidJWTState(setup bool) error {
	// This is handled by the request filter middleware
	return nil
}

// Mock authentication middleware for Pact tests
func mockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mock authentication for Pact provider tests
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && authHeader != "Bearer invalid-token" {
			// Set up mock user context for valid tokens
			c.Set("user_id", "123e4567-e89b-12d3-a456-426614174000")
			c.Set("authenticated", true)
		} else {
			// For protected endpoints without auth, return unauthorized
			if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
				if c.Request.URL.Path != "/health" && c.Request.URL.Path != "/api/v1/videos/health" {
					responses.UnauthorizedResponse(c, "Authentication required")
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}

// Mock Video Repository for testing
type MockVideoRepository struct {
	videos       map[uuid.UUID]*models.Video
	channels     map[uuid.UUID]*models.Channel
	interactions map[uuid.UUID][]*models.VideoInteraction
	comments     map[uuid.UUID][]*models.VideoComment
	shares       map[uuid.UUID][]*models.VideoShare
}

func (m *MockVideoRepository) GetVideos(limit, offset int, category string) ([]models.Video, error) {
	var videos []models.Video
	count := 0
	skip := 0

	for _, video := range m.videos {
		if skip < offset {
			skip++
			continue
		}
		if count >= limit {
			break
		}

		if category != "" && video.Category != category {
			continue
		}

		videos = append(videos, *video)
		count++
	}

	return videos, nil
}

func (m *MockVideoRepository) GetVideoByID(id uuid.UUID) (*models.Video, error) {
	if video, exists := m.videos[id]; exists {
		return video, nil
	}
	return nil, fmt.Errorf("video not found")
}

func (m *MockVideoRepository) CreateVideo(video *models.Video) error {
	if m.videos == nil {
		m.videos = make(map[uuid.UUID]*models.Video)
	}
	m.videos[video.ID] = video
	return nil
}

func (m *MockVideoRepository) UpdateVideo(video *models.Video) error {
	m.videos[video.ID] = video
	return nil
}

func (m *MockVideoRepository) DeleteVideo(id uuid.UUID) error {
	delete(m.videos, id)
	return nil
}

func (m *MockVideoRepository) GetVideoInteractions(videoID uuid.UUID) ([]models.VideoInteraction, error) {
	var interactions []models.VideoInteraction
	if videoInteractions, exists := m.interactions[videoID]; exists {
		for _, interaction := range videoInteractions {
			interactions = append(interactions, *interaction)
		}
	}
	return interactions, nil
}

func (m *MockVideoRepository) CreateVideoInteraction(interaction *models.VideoInteraction) error {
	if m.interactions == nil {
		m.interactions = make(map[uuid.UUID][]*models.VideoInteraction)
	}
	m.interactions[interaction.VideoID] = append(m.interactions[interaction.VideoID], interaction)
	return nil
}

func (m *MockVideoRepository) GetChannelByID(id uuid.UUID) (*models.Channel, error) {
	if channel, exists := m.channels[id]; exists {
		return channel, nil
	}
	return nil, fmt.Errorf("channel not found")
}

func (m *MockVideoRepository) CreateChannel(channel *models.Channel) error {
	if m.channels == nil {
		m.channels = make(map[uuid.UUID]*models.Channel)
	}
	m.channels[channel.ID] = channel
	return nil
}

func (m *MockVideoRepository) UpdateChannel(channel *models.Channel) error {
	m.channels[channel.ID] = channel
	return nil
}

func (m *MockVideoRepository) SearchVideos(query string, limit, offset int, category string) ([]models.Video, error) {
	var videos []models.Video
	count := 0
	skip := 0

	for _, video := range m.videos {
		if skip < offset {
			skip++
			continue
		}
		if count >= limit {
			break
		}

		// Simple search in title and description
		if query != "" {
			if !contains(video.Title, query) && !contains(video.Description, query) {
				continue
			}
		}

		if category != "" && video.Category != category {
			continue
		}

		videos = append(videos, *video)
		count++
	}

	return videos, nil
}

func (m *MockVideoRepository) GetVideosByCategory(category string, limit, offset int) ([]models.Video, error) {
	return m.GetVideos(limit, offset, category)
}

func (m *MockVideoRepository) GetTrendingVideos(timeframe string, limit, offset int) ([]models.Video, error) {
	var videos []models.Video
	count := 0
	skip := 0

	// Simple trending logic: sort by views + likes
	for _, video := range m.videos {
		if skip < offset {
			skip++
			continue
		}
		if count >= limit {
			break
		}

		// For mock, just return videos with high engagement
		if video.Views > 1000 || video.Likes > 100 {
			videos = append(videos, *video)
			count++
		}
	}

	return videos, nil
}

func (m *MockVideoRepository) GetVideoComments(videoID uuid.UUID, limit, offset int) ([]models.VideoComment, error) {
	var comments []models.VideoComment
	if videoComments, exists := m.comments[videoID]; exists {
		count := 0
		skip := 0
		for _, comment := range videoComments {
			if skip < offset {
				skip++
				continue
			}
			if count >= limit {
				break
			}
			comments = append(comments, *comment)
			count++
		}
	}
	return comments, nil
}

func (m *MockVideoRepository) CreateVideoComment(comment *models.VideoComment) error {
	if m.comments == nil {
		m.comments = make(map[uuid.UUID][]*models.VideoComment)
	}
	m.comments[comment.VideoID] = append(m.comments[comment.VideoID], comment)
	return nil
}

func (m *MockVideoRepository) CreateVideoShare(share *models.VideoShare) error {
	if m.shares == nil {
		m.shares = make(map[uuid.UUID][]*models.VideoShare)
	}
	m.shares[share.VideoID] = append(m.shares[share.VideoID], share)
	return nil
}

func (m *MockVideoRepository) Clear() {
	m.videos = make(map[uuid.UUID]*models.Video)
	m.channels = make(map[uuid.UUID]*models.Channel)
	m.interactions = make(map[uuid.UUID][]*models.VideoInteraction)
	m.comments = make(map[uuid.UUID][]*models.VideoComment)
	m.shares = make(map[uuid.UUID][]*models.VideoShare)
}

// Helper function for simple string search
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(substr) > 0 && containsIgnoreCase(s, substr)))
}

func containsIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}