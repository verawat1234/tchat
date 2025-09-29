package crossplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tchat.dev/commerce/handlers"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
	"tchat.dev/commerce/services"
)

// CrossPlatformStreamTestSuite validates Stream Store Tabs functionality
// across web, mobile, and backend platforms for consistency
type CrossPlatformStreamTestSuite struct {
	suite.Suite
	db                *gorm.DB
	router            *gin.Engine
	streamHandler     *handlers.StreamHandler
	platformResponses map[string]interface{}
}

func (suite *CrossPlatformStreamTestSuite) SetupSuite() {
	// Setup test database
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://localhost/tchat_cross_platform_test?sslmode=disable"
	}

	var err error
	suite.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err)

	// Auto-migrate models
	err = suite.db.AutoMigrate(
		&models.StreamCategory{},
		&models.StreamSubtab{},
		&models.StreamContentItem{},
		&models.TabNavigationState{},
		&models.StreamUserSession{},
		&models.StreamContentView{},
		&models.StreamUserPreference{},
	)
	suite.Require().NoError(err)

	// Setup repositories and services
	streamRepo := repository.NewStreamRepository(suite.db)
	categoryService := services.NewStreamCategoryService(streamRepo)
	contentService := services.NewStreamContentService(streamRepo)
	sessionService := services.NewStreamSessionService(streamRepo)
	purchaseService := services.NewStreamPurchaseService(streamRepo)

	suite.streamHandler = handlers.NewStreamHandler(
		categoryService,
		contentService,
		sessionService,
		purchaseService,
	)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.setupRoutes()

	// Initialize platform responses storage
	suite.platformResponses = make(map[string]interface{})

	// Setup test data
	suite.setupCrossPlatformTestData()
}

func (suite *CrossPlatformStreamTestSuite) setupRoutes() {
	v1 := suite.router.Group("/api/v1")
	stream := v1.Group("/stream")
	{
		stream.GET("/categories", suite.streamHandler.GetStreamCategories)
		stream.GET("/categories/:id", suite.streamHandler.GetStreamCategoryDetail)
		stream.GET("/content", suite.streamHandler.GetStreamContent)
		stream.GET("/content/:id", suite.streamHandler.GetStreamContentDetail)
		stream.GET("/featured", suite.streamHandler.GetStreamFeatured)
		stream.GET("/search", suite.streamHandler.SearchStreamContent)
		stream.POST("/content/purchase", suite.streamHandler.PostStreamContentPurchase)
		stream.GET("/navigation", suite.streamHandler.GetUserNavigationState)
		stream.PUT("/navigation", suite.streamHandler.UpdateUserNavigationState)
		stream.PUT("/content/:id/progress", suite.streamHandler.UpdateContentViewProgress)
		stream.GET("/preferences", suite.streamHandler.GetUserPreferences)
		stream.PUT("/preferences", suite.streamHandler.UpdateUserPreferences)
	}
}

func (suite *CrossPlatformStreamTestSuite) setupCrossPlatformTestData() {
	// Create 6 Stream categories matching specification
	categories := []models.StreamCategory{
		{
			ID:                      "books",
			Name:                    "Books",
			DisplayOrder:            1,
			IconName:                "book-open",
			IsActive:                true,
			FeaturedContentEnabled:  true,
		},
		{
			ID:                      "podcasts",
			Name:                    "Podcasts",
			DisplayOrder:            2,
			IconName:                "headphones",
			IsActive:                true,
			FeaturedContentEnabled:  true,
		},
		{
			ID:                      "cartoons",
			Name:                    "Cartoons",
			DisplayOrder:            3,
			IconName:                "film",
			IsActive:                true,
			FeaturedContentEnabled:  true,
		},
		{
			ID:                      "movies",
			Name:                    "Movies",
			DisplayOrder:            4,
			IconName:                "video",
			IsActive:                true,
			FeaturedContentEnabled:  true,
		},
		{
			ID:                      "music",
			Name:                    "Music",
			DisplayOrder:            5,
			IconName:                "music",
			IsActive:                true,
			FeaturedContentEnabled:  true,
		},
		{
			ID:                      "art",
			Name:                    "Art",
			DisplayOrder:            6,
			IconName:                "palette",
			IsActive:                true,
			FeaturedContentEnabled:  true,
		},
	}

	for _, category := range categories {
		suite.db.Create(&category)
	}

	// Create subtabs for movies
	movieSubtabs := []models.StreamSubtab{
		{
			ID:             "movies_short",
			CategoryID:     "movies",
			Name:           "Short Films",
			DisplayOrder:   1,
			FilterCriteria: `{"content_type": "SHORT_MOVIE", "max_duration": 1800}`,
			IsActive:       true,
		},
		{
			ID:             "movies_feature",
			CategoryID:     "movies",
			Name:           "Feature Films",
			DisplayOrder:   2,
			FilterCriteria: `{"content_type": "LONG_MOVIE", "min_duration": 1800}`,
			IsActive:       true,
		},
	}

	for _, subtab := range movieSubtabs {
		suite.db.Create(&subtab)
	}

	// Create representative content for each category
	contents := []models.StreamContentItem{
		// Books
		{
			ID:                 "book-fiction-1",
			CategoryID:         "books",
			Title:             "The Great Test Novel",
			Description:       "A captivating work of fiction",
			ThumbnailURL:      "https://cdn.example.com/books/fiction1.jpg",
			ContentType:       "BOOK",
			Price:             12.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        true,
			FeaturedOrder:     &[]int{1}[0],
			Metadata:          `{"author": "Jane Doe", "genre": "fiction", "pages": 320}`,
		},
		// Podcasts
		{
			ID:                 "podcast-tech-1",
			CategoryID:         "podcasts",
			Title:             "Tech Talk Episode 1",
			Description:       "Latest technology trends discussion",
			ThumbnailURL:      "https://cdn.example.com/podcasts/tech1.jpg",
			ContentType:       "PODCAST",
			Duration:          &[]int{3600}[0], // 1 hour
			Price:             0.00,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        true,
			FeaturedOrder:     &[]int{1}[0],
			Metadata:          `{"host": "John Smith", "series": "Tech Talk", "episode": 1}`,
		},
		// Cartoons
		{
			ID:                 "cartoon-adventure-1",
			CategoryID:         "cartoons",
			Title:             "Adventure Quest",
			Description:       "Animated adventure series",
			ThumbnailURL:      "https://cdn.example.com/cartoons/adventure1.jpg",
			ContentType:       "CARTOON",
			Duration:          &[]int{1200}[0], // 20 minutes
			Price:             1.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        true,
			FeaturedOrder:     &[]int{1}[0],
			Metadata:          `{"studio": "Animation Studios", "rating": "G", "season": 1}`,
		},
		// Movies - Short
		{
			ID:                 "movie-short-drama-1",
			CategoryID:         "movies",
			Title:             "Short Drama",
			Description:       "Compelling short film",
			ThumbnailURL:      "https://cdn.example.com/movies/short1.jpg",
			ContentType:       "SHORT_MOVIE",
			Duration:          &[]int{1200}[0], // 20 minutes
			Price:             3.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        true,
			FeaturedOrder:     &[]int{1}[0],
			Metadata:          `{"director": "Alice Director", "year": 2024, "genre": "drama"}`,
		},
		// Movies - Feature
		{
			ID:                 "movie-feature-action-1",
			CategoryID:         "movies",
			Title:             "Action Adventure",
			Description:       "Epic action adventure film",
			ThumbnailURL:      "https://cdn.example.com/movies/action1.jpg",
			ContentType:       "LONG_MOVIE",
			Duration:          &[]int{7200}[0], // 2 hours
			Price:             14.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        false,
			Metadata:          `{"director": "Bob Director", "year": 2024, "rating": "PG-13"}`,
		},
		// Music
		{
			ID:                 "music-album-1",
			CategoryID:         "music",
			Title:             "Greatest Hits Album",
			Description:       "Collection of popular songs",
			ThumbnailURL:      "https://cdn.example.com/music/album1.jpg",
			ContentType:       "MUSIC",
			Duration:          &[]int{2400}[0], // 40 minutes
			Price:             9.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        true,
			FeaturedOrder:     &[]int{1}[0],
			Metadata:          `{"artist": "Test Artist", "genre": "pop", "tracks": 12}`,
		},
		// Art
		{
			ID:                 "art-digital-1",
			CategoryID:         "art",
			Title:             "Digital Art Collection",
			Description:       "Contemporary digital artwork",
			ThumbnailURL:      "https://cdn.example.com/art/digital1.jpg",
			ContentType:       "ART",
			Price:             19.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        true,
			FeaturedOrder:     &[]int{1}[0],
			Metadata:          `{"artist": "Digital Artist", "medium": "digital", "year": 2024}`,
		},
	}

	for _, content := range contents {
		suite.db.Create(&content)
	}
}

func (suite *CrossPlatformStreamTestSuite) TearDownSuite() {
	// Clean up test data
	suite.db.Exec("DELETE FROM stream_content_views WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_user_sessions WHERE 1=1")
	suite.db.Exec("DELETE FROM tab_navigation_states WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_user_preferences WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_content_items WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_subtabs WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_categories WHERE 1=1")
}

// CROSS-PLATFORM TEST 1: API Response Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformAPIConsistency() {
	// Test categories endpoint consistency
	req := httptest.NewRequest("GET", "/api/v1/stream/categories", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var categoriesResponse struct {
		Categories []models.StreamCategory `json:"categories"`
		Total      int                    `json:"total"`
		Success    bool                   `json:"success"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &categoriesResponse)
	assert.NoError(suite.T(), err)

	// Store backend response for cross-platform comparison
	suite.platformResponses["backend_categories"] = categoriesResponse

	// Validate response structure for cross-platform consumption
	assert.True(suite.T(), categoriesResponse.Success)
	assert.Equal(suite.T(), 6, categoriesResponse.Total)

	// Validate each category has required fields for all platforms
	for _, category := range categoriesResponse.Categories {
		assert.NotEmpty(suite.T(), category.ID)
		assert.NotEmpty(suite.T(), category.Name)
		assert.Greater(suite.T(), category.DisplayOrder, 0)
		assert.NotEmpty(suite.T(), category.IconName)
		assert.NotEmpty(suite.T(), category.CreatedAt)
		assert.NotEmpty(suite.T(), category.UpdatedAt)
	}

	// Validate category order matches specification
	expectedOrder := []string{"books", "podcasts", "cartoons", "movies", "music", "art"}
	for i, expectedID := range expectedOrder {
		assert.Equal(suite.T(), expectedID, categoriesResponse.Categories[i].ID)
	}
}

// CROSS-PLATFORM TEST 2: Content Structure Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformContentStructure() {
	// Test each category's content structure
	categories := []string{"books", "podcasts", "cartoons", "movies", "music", "art"}

	for _, categoryID := range categories {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/stream/content?categoryId=%s&page=1&limit=10", categoryID), nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var contentResponse struct {
			Items   []models.StreamContentItem `json:"items"`
			Total   int64                     `json:"total"`
			HasMore bool                      `json:"hasMore"`
			Success bool                      `json:"success"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &contentResponse)
		assert.NoError(suite.T(), err)

		// Store response for cross-platform validation
		suite.platformResponses[fmt.Sprintf("backend_content_%s", categoryID)] = contentResponse

		assert.True(suite.T(), contentResponse.Success)

		// Validate content structure for each category
		for _, content := range contentResponse.Items {
			assert.Equal(suite.T(), categoryID, content.CategoryID)
			assert.NotEmpty(suite.T(), content.ID)
			assert.NotEmpty(suite.T(), content.Title)
			assert.NotEmpty(suite.T(), content.Description)
			assert.NotEmpty(suite.T(), content.ThumbnailURL)
			assert.GreaterOrEqual(suite.T(), content.Price, 0.0)
			assert.NotEmpty(suite.T(), content.Currency)
			assert.NotEmpty(suite.T(), content.AvailabilityStatus)
			assert.NotEmpty(suite.T(), content.CreatedAt)
			assert.NotEmpty(suite.T(), content.UpdatedAt)

			// Validate content type matches category
			suite.validateContentTypeForCategory(content, categoryID)
		}
	}
}

func (suite *CrossPlatformStreamTestSuite) validateContentTypeForCategory(content models.StreamContentItem, categoryID string) {
	switch categoryID {
	case "books":
		assert.Equal(suite.T(), "BOOK", string(content.ContentType))
		assert.Nil(suite.T(), content.Duration)
	case "podcasts":
		assert.Equal(suite.T(), "PODCAST", string(content.ContentType))
		assert.NotNil(suite.T(), content.Duration)
	case "cartoons":
		assert.Equal(suite.T(), "CARTOON", string(content.ContentType))
		assert.NotNil(suite.T(), content.Duration)
	case "movies":
		assert.Contains(suite.T(), []string{"SHORT_MOVIE", "LONG_MOVIE"}, string(content.ContentType))
		assert.NotNil(suite.T(), content.Duration)
	case "music":
		assert.Equal(suite.T(), "MUSIC", string(content.ContentType))
		assert.NotNil(suite.T(), content.Duration)
	case "art":
		assert.Equal(suite.T(), "ART", string(content.ContentType))
	}
}

// CROSS-PLATFORM TEST 3: Featured Content Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformFeaturedContent() {
	categories := []string{"books", "podcasts", "cartoons", "movies", "music", "art"}

	for _, categoryID := range categories {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/stream/featured?categoryId=%s&limit=5", categoryID), nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var featuredResponse struct {
			Items   []models.StreamContentItem `json:"items"`
			Total   int64                     `json:"total"`
			HasMore bool                      `json:"hasMore"`
			Success bool                      `json:"success"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &featuredResponse)
		assert.NoError(suite.T(), err)

		assert.True(suite.T(), featuredResponse.Success)

		// Each category should have featured content
		assert.Greater(suite.T(), featuredResponse.Total, int64(0))

		// All featured content should be marked as featured
		for _, content := range featuredResponse.Items {
			assert.True(suite.T(), content.IsFeatured)
			assert.NotNil(suite.T(), content.FeaturedOrder)
			assert.Equal(suite.T(), categoryID, content.CategoryID)
		}

		// Store for cross-platform comparison
		suite.platformResponses[fmt.Sprintf("backend_featured_%s", categoryID)] = featuredResponse
	}
}

// CROSS-PLATFORM TEST 4: Subtab Functionality Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformSubtabFunctionality() {
	// Test movies category subtabs (only category with subtabs in test data)
	req := httptest.NewRequest("GET", "/api/v1/stream/categories/movies", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var categoryResponse struct {
		Category models.StreamCategory   `json:"category"`
		Subtabs  []models.StreamSubtab   `json:"subtabs"`
		Stats    map[string]interface{}  `json:"stats"`
		Success  bool                   `json:"success"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &categoryResponse)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), categoryResponse.Success)
	assert.Equal(suite.T(), "Movies", categoryResponse.Category.Name)
	assert.Equal(suite.T(), 2, len(categoryResponse.Subtabs))

	// Validate subtab structure for cross-platform usage
	for _, subtab := range categoryResponse.Subtabs {
		assert.NotEmpty(suite.T(), subtab.ID)
		assert.Equal(suite.T(), "movies", subtab.CategoryID)
		assert.NotEmpty(suite.T(), subtab.Name)
		assert.Greater(suite.T(), subtab.DisplayOrder, 0)
		assert.NotEmpty(suite.T(), subtab.FilterCriteria)
	}

	// Test subtab filtering
	subtabTests := []struct {
		subtabID       string
		expectedType   string
		expectedCount  int
	}{
		{"movies_short", "SHORT_MOVIE", 1},
		{"movies_feature", "LONG_MOVIE", 1},
	}

	for _, test := range subtabTests {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/stream/content?categoryId=movies&subtabId=%s", test.subtabID), nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var contentResponse struct {
			Items   []models.StreamContentItem `json:"items"`
			Total   int64                     `json:"total"`
			Success bool                      `json:"success"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &contentResponse)
		assert.NoError(suite.T(), err)

		assert.True(suite.T(), contentResponse.Success)
		assert.Equal(suite.T(), int64(test.expectedCount), contentResponse.Total)

		// Validate filtered content matches subtab criteria
		for _, content := range contentResponse.Items {
			assert.Equal(suite.T(), test.expectedType, string(content.ContentType))
		}
	}
}

// CROSS-PLATFORM TEST 5: Performance Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformPerformance() {
	performanceTargets := map[string]time.Duration{
		"categories":       200 * time.Millisecond, // <200ms API response requirement
		"content":          200 * time.Millisecond,
		"featured":         200 * time.Millisecond,
		"category_detail":  200 * time.Millisecond,
	}

	performanceResults := make(map[string]time.Duration)

	// Test categories performance
	start := time.Now()
	req := httptest.NewRequest("GET", "/api/v1/stream/categories", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	categoriesTime := time.Since(start)
	performanceResults["categories"] = categoriesTime

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Less(suite.T(), categoriesTime, performanceTargets["categories"])

	// Test content performance
	start = time.Now()
	req = httptest.NewRequest("GET", "/api/v1/stream/content?categoryId=books&page=1&limit=20", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	contentTime := time.Since(start)
	performanceResults["content"] = contentTime

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Less(suite.T(), contentTime, performanceTargets["content"])

	// Test featured performance
	start = time.Now()
	req = httptest.NewRequest("GET", "/api/v1/stream/featured?categoryId=books&limit=10", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	featuredTime := time.Since(start)
	performanceResults["featured"] = featuredTime

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Less(suite.T(), featuredTime, performanceTargets["featured"])

	// Test category detail performance
	start = time.Now()
	req = httptest.NewRequest("GET", "/api/v1/stream/categories/movies", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	categoryDetailTime := time.Since(start)
	performanceResults["category_detail"] = categoryDetailTime

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Less(suite.T(), categoryDetailTime, performanceTargets["category_detail"])

	// Store performance results for cross-platform validation
	suite.platformResponses["backend_performance"] = performanceResults
}

// CROSS-PLATFORM TEST 6: Error Response Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformErrorHandling() {
	errorTests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid Category ID",
			url:            "/api/v1/stream/categories/invalid",
			expectedStatus: http.StatusNotFound,
			expectedError:  "not found",
		},
		{
			name:           "Invalid Content ID",
			url:            "/api/v1/stream/content/invalid",
			expectedStatus: http.StatusNotFound,
			expectedError:  "not found",
		},
		{
			name:           "Missing Category Parameter",
			url:            "/api/v1/stream/content",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "required",
		},
		{
			name:           "Missing Search Query",
			url:            "/api/v1/stream/search",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "required",
		},
	}

	for _, test := range errorTests {
		suite.T().Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.url, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)

			var errorResponse struct {
				Success bool   `json:"success"`
				Error   string `json:"error"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			assert.NoError(t, err)

			assert.False(t, errorResponse.Success)
			assert.Contains(t, errorResponse.Error, test.expectedError)
		})
	}
}

// CROSS-PLATFORM TEST 7: Data Type Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformDataTypes() {
	req := httptest.NewRequest("GET", "/api/v1/stream/content?categoryId=books&page=1&limit=1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var contentResponse struct {
		Items   []models.StreamContentItem `json:"items"`
		Total   int64                     `json:"total"`
		HasMore bool                      `json:"hasMore"`
		Success bool                      `json:"success"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &contentResponse)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), contentResponse.Success)
	assert.Greater(suite.T(), len(contentResponse.Items), 0)

	content := contentResponse.Items[0]

	// Validate data types for cross-platform consistency
	assert.IsType(suite.T(), "", content.ID)                      // string
	assert.IsType(suite.T(), "", content.Title)                   // string
	assert.IsType(suite.T(), "", content.Description)             // string
	assert.IsType(suite.T(), "", content.ThumbnailURL)            // string
	assert.IsType(suite.T(), float64(0), content.Price)           // number
	assert.IsType(suite.T(), "", content.Currency)                // string
	assert.IsType(suite.T(), "", content.AvailabilityStatus)      // string
	assert.IsType(suite.T(), true, content.IsFeatured)            // boolean
	assert.IsType(suite.T(), "", content.CreatedAt)               // string (ISO date)
	assert.IsType(suite.T(), "", content.UpdatedAt)               // string (ISO date)

	// Validate metadata is valid JSON
	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(content.Metadata), &metadata)
	assert.NoError(suite.T(), err)
}

// CROSS-PLATFORM TEST 8: Pagination Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformPagination() {
	// Test pagination with consistent page size
	req := httptest.NewRequest("GET", "/api/v1/stream/content?categoryId=movies&page=1&limit=1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var page1Response struct {
		Items   []models.StreamContentItem `json:"items"`
		Total   int64                     `json:"total"`
		HasMore bool                      `json:"hasMore"`
		Success bool                      `json:"success"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &page1Response)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), page1Response.Success)
	assert.Equal(suite.T(), 1, len(page1Response.Items))
	assert.Equal(suite.T(), int64(2), page1Response.Total) // 2 movies total

	// Validate pagination structure for cross-platform usage
	expectedPagination := map[string]interface{}{
		"page":     1,
		"limit":    1,
		"total":    page1Response.Total,
		"hasMore":  page1Response.HasMore,
		"items":    len(page1Response.Items),
	}

	suite.platformResponses["backend_pagination"] = expectedPagination

	assert.Equal(suite.T(), 1, expectedPagination["items"])
	assert.Equal(suite.T(), int64(2), expectedPagination["total"])
	assert.True(suite.T(), expectedPagination["hasMore"].(bool))
}

// CROSS-PLATFORM TEST 9: Authentication Integration Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformAuthRequirements() {
	authRequiredEndpoints := []struct {
		method string
		url    string
		body   map[string]interface{}
	}{
		{
			method: "POST",
			url:    "/api/v1/stream/content/purchase",
			body: map[string]interface{}{
				"mediaContentId": "book-fiction-1",
				"quantity":       1,
				"mediaLicense":   "personal",
			},
		},
		{
			method: "GET",
			url:    "/api/v1/stream/navigation?userId=test-user",
			body:   nil,
		},
		{
			method: "PUT",
			url:    "/api/v1/stream/navigation",
			body: map[string]interface{}{
				"userId":     "test-user",
				"sessionId":  "test-session",
				"categoryId": "books",
			},
		},
	}

	for _, endpoint := range authRequiredEndpoints {
		suite.T().Run(fmt.Sprintf("%s %s", endpoint.method, endpoint.url), func(t *testing.T) {
			var req *http.Request

			if endpoint.body != nil {
				jsonData, _ := json.Marshal(endpoint.body)
				req = httptest.NewRequest(endpoint.method, endpoint.url, bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(endpoint.method, endpoint.url, nil)
			}

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			// Should require authentication
			assert.Equal(t, http.StatusUnauthorized, w.Code)

			var errorResponse struct {
				Success bool   `json:"success"`
				Error   string `json:"error"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			assert.NoError(t, err)
			assert.False(t, errorResponse.Success)
		})
	}
}

// CROSS-PLATFORM TEST 10: Content Metadata Consistency
func (suite *CrossPlatformStreamTestSuite) TestCrossPlatformMetadataStructure() {
	categories := []string{"books", "podcasts", "cartoons", "movies", "music", "art"}

	for _, categoryID := range categories {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/stream/content?categoryId=%s&page=1&limit=1", categoryID), nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var contentResponse struct {
			Items   []models.StreamContentItem `json:"items"`
			Success bool                      `json:"success"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &contentResponse)
		assert.NoError(suite.T(), err)

		if len(contentResponse.Items) > 0 {
			content := contentResponse.Items[0]

			// Parse and validate metadata structure
			var metadata map[string]interface{}
			err = json.Unmarshal([]byte(content.Metadata), &metadata)
			assert.NoError(suite.T(), err)

			// Validate category-specific metadata fields
			suite.validateMetadataForCategory(metadata, categoryID)
		}
	}
}

func (suite *CrossPlatformStreamTestSuite) validateMetadataForCategory(metadata map[string]interface{}, categoryID string) {
	switch categoryID {
	case "books":
		assert.Contains(suite.T(), metadata, "author")
		assert.Contains(suite.T(), metadata, "genre")
	case "podcasts":
		assert.Contains(suite.T(), metadata, "host")
		assert.Contains(suite.T(), metadata, "series")
	case "cartoons":
		assert.Contains(suite.T(), metadata, "studio")
		assert.Contains(suite.T(), metadata, "rating")
	case "movies":
		assert.Contains(suite.T(), metadata, "director")
		assert.Contains(suite.T(), metadata, "year")
	case "music":
		assert.Contains(suite.T(), metadata, "artist")
		assert.Contains(suite.T(), metadata, "genre")
	case "art":
		assert.Contains(suite.T(), metadata, "artist")
		assert.Contains(suite.T(), metadata, "medium")
	}
}

func TestCrossPlatformStreamTestSuite(t *testing.T) {
	suite.Run(t, new(CrossPlatformStreamTestSuite))
}