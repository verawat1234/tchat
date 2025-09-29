package integration

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

// StreamIntegrationTestSuite tests complete Stream Store Tabs flows
type StreamIntegrationTestSuite struct {
	suite.Suite
	db             *gorm.DB
	router         *gin.Engine
	streamHandler  *handlers.StreamHandler
	testUserID     string
	testSessionID  string

	// Test data
	testCategories []models.StreamCategory
	testContent    []models.StreamContentItem
}

func (suite *StreamIntegrationTestSuite) SetupSuite() {
	// Setup test database
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://localhost/tchat_test?sslmode=disable"
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

	// Setup routes
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

	// Setup test data
	suite.testUserID = "test-user-123"
	suite.testSessionID = "test-session-456"
	suite.setupTestData()
}

func (suite *StreamIntegrationTestSuite) setupTestData() {
	// Create test categories
	suite.testCategories = []models.StreamCategory{
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
			ID:                      "movies",
			Name:                    "Movies",
			DisplayOrder:            4,
			IconName:                "video",
			IsActive:                true,
			FeaturedContentEnabled:  true,
		},
	}

	for _, category := range suite.testCategories {
		suite.db.Create(&category)
	}

	// Create test subtabs for movies
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

	// Create test content
	suite.testContent = []models.StreamContentItem{
		{
			ID:                 "book-1",
			CategoryID:         "books",
			Title:             "Test Book 1",
			Description:       "A fascinating test book",
			ThumbnailURL:      "https://example.com/book1.jpg",
			ContentType:       "BOOK",
			Price:             9.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        true,
			FeaturedOrder:     &[]int{1}[0],
			Metadata:          `{"author": "Test Author", "genre": "fiction"}`,
		},
		{
			ID:                 "podcast-1",
			CategoryID:         "podcasts",
			Title:             "Test Podcast Episode",
			Description:       "An interesting podcast episode",
			ThumbnailURL:      "https://example.com/podcast1.jpg",
			ContentType:       "PODCAST",
			Duration:          &[]int{2400}[0], // 40 minutes
			Price:             0.00,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        false,
			Metadata:          `{"host": "Test Host", "series": "Test Series"}`,
		},
		{
			ID:                 "movie-short-1",
			CategoryID:         "movies",
			Title:             "Test Short Film",
			Description:       "A compelling short film",
			ThumbnailURL:      "https://example.com/short1.jpg",
			ContentType:       "SHORT_MOVIE",
			Duration:          &[]int{1200}[0], // 20 minutes
			Price:             2.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        true,
			FeaturedOrder:     &[]int{1}[0],
			Metadata:          `{"director": "Test Director", "year": 2024}`,
		},
		{
			ID:                 "movie-feature-1",
			CategoryID:         "movies",
			Title:             "Test Feature Film",
			Description:       "An epic feature film",
			ThumbnailURL:      "https://example.com/feature1.jpg",
			ContentType:       "LONG_MOVIE",
			Duration:          &[]int{7200}[0], // 2 hours
			Price:             12.99,
			Currency:          "USD",
			AvailabilityStatus: "AVAILABLE",
			IsFeatured:        false,
			Metadata:          `{"director": "Test Director", "year": 2024, "rating": "PG-13"}`,
		},
	}

	for _, content := range suite.testContent {
		suite.db.Create(&content)
	}
}

func (suite *StreamIntegrationTestSuite) TearDownSuite() {
	// Clean up test data
	suite.db.Exec("DELETE FROM stream_content_views WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_user_sessions WHERE 1=1")
	suite.db.Exec("DELETE FROM tab_navigation_states WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_user_preferences WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_content_items WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_subtabs WHERE 1=1")
	suite.db.Exec("DELETE FROM stream_categories WHERE 1=1")
}

// FLOW TEST 1: Complete Navigation Flow
func (suite *StreamIntegrationTestSuite) TestCompleteNavigationFlow() {
	// Step 1: Get all categories
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
	assert.True(suite.T(), categoriesResponse.Success)
	assert.Equal(suite.T(), 3, categoriesResponse.Total)
	assert.Equal(suite.T(), "Books", categoriesResponse.Categories[0].Name)

	// Step 2: Navigate to specific category (Movies)
	req = httptest.NewRequest("GET", "/api/v1/stream/categories/movies", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var categoryResponse struct {
		Category models.StreamCategory   `json:"category"`
		Subtabs  []models.StreamSubtab   `json:"subtabs"`
		Stats    map[string]interface{}  `json:"stats"`
		Success  bool                   `json:"success"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &categoryResponse)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), categoryResponse.Success)
	assert.Equal(suite.T(), "Movies", categoryResponse.Category.Name)
	assert.Equal(suite.T(), 2, len(categoryResponse.Subtabs))

	// Step 3: Get content for category
	req = httptest.NewRequest("GET", "/api/v1/stream/content?categoryId=movies&page=1&limit=10", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var contentResponse struct {
		Items   []models.StreamContentItem `json:"items"`
		Total   int64                     `json:"total"`
		HasMore bool                      `json:"hasMore"`
		Success bool                      `json:"success"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &contentResponse)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), contentResponse.Success)
	assert.Equal(suite.T(), int64(2), contentResponse.Total) // 2 movies
}

// FLOW TEST 2: Featured Content Discovery Flow
func (suite *StreamIntegrationTestSuite) TestFeaturedContentFlow() {
	// Get featured content for books category
	req := httptest.NewRequest("GET", "/api/v1/stream/featured?categoryId=books&limit=5", nil)
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
	assert.Equal(suite.T(), int64(1), featuredResponse.Total) // 1 featured book
	assert.True(suite.T(), featuredResponse.Items[0].IsFeatured)
	assert.Equal(suite.T(), "Test Book 1", featuredResponse.Items[0].Title)
}

// FLOW TEST 3: Content Purchase Flow
func (suite *StreamIntegrationTestSuite) TestContentPurchaseFlow() {
	// Mock authentication by setting user header
	purchaseData := map[string]interface{}{
		"mediaContentId": "book-1",
		"quantity":       1,
		"mediaLicense":   "personal",
	}

	jsonData, _ := json.Marshal(purchaseData)
	req := httptest.NewRequest("POST", "/api/v1/stream/content/purchase", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("X-User-ID", suite.testUserID) // Mock user context

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Note: Without real auth middleware, this will return 401
	// In real integration tests, you'd setup proper auth context
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

// FLOW TEST 4: User Navigation State Flow
func (suite *StreamIntegrationTestSuite) TestNavigationStateFlow() {
	// Update navigation state
	navigationData := map[string]interface{}{
		"userId":         suite.testUserID,
		"sessionId":      suite.testSessionID,
		"categoryId":     "movies",
		"subtabId":       "movies_short",
	}

	jsonData, _ := json.Marshal(navigationData)
	req := httptest.NewRequest("PUT", "/api/v1/stream/navigation", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", suite.testUserID)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code) // Auth required

	// Get navigation state
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/stream/navigation?userId=%s", suite.testUserID), nil)
	req.Header.Set("X-User-ID", suite.testUserID)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code) // Auth required
}

// FLOW TEST 5: Content Search Flow
func (suite *StreamIntegrationTestSuite) TestContentSearchFlow() {
	// Search for content across categories
	req := httptest.NewRequest("GET", "/api/v1/stream/search?q=test&page=1&limit=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var searchResponse struct {
		Items   []models.StreamContentItem `json:"items"`
		Total   int64                     `json:"total"`
		HasMore bool                      `json:"hasMore"`
		Success bool                      `json:"success"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &searchResponse)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), searchResponse.Success)
	assert.Greater(suite.T(), searchResponse.Total, int64(0)) // Should find test content
}

// FLOW TEST 6: Subtab Filtering Flow
func (suite *StreamIntegrationTestSuite) TestSubtabFilteringFlow() {
	// Get short movies only
	req := httptest.NewRequest("GET", "/api/v1/stream/content?categoryId=movies&subtabId=movies_short&page=1&limit=10", nil)
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

	// Should only return short movies
	for _, item := range contentResponse.Items {
		assert.Equal(suite.T(), "SHORT_MOVIE", string(item.ContentType))
		if item.Duration != nil {
			assert.LessOrEqual(suite.T(), *item.Duration, 1800) // ≤ 30 minutes
		}
	}
}

// FLOW TEST 7: Content Detail Flow
func (suite *StreamIntegrationTestSuite) TestContentDetailFlow() {
	// Get specific content detail
	req := httptest.NewRequest("GET", "/api/v1/stream/content/book-1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var detailResponse struct {
		Content models.StreamContentItem `json:"content"`
		Success bool                    `json:"success"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &detailResponse)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), detailResponse.Success)
	assert.Equal(suite.T(), "Test Book 1", detailResponse.Content.Title)
	assert.Equal(suite.T(), "books", detailResponse.Content.CategoryID)
}

// FLOW TEST 8: Pagination Flow
func (suite *StreamIntegrationTestSuite) TestPaginationFlow() {
	// Test pagination with small limit
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
	assert.True(suite.T(), page1Response.HasMore) // Should have more items

	// Get page 2
	req = httptest.NewRequest("GET", "/api/v1/stream/content?categoryId=movies&page=2&limit=1", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var page2Response struct {
		Items   []models.StreamContentItem `json:"items"`
		Total   int64                     `json:"total"`
		HasMore bool                      `json:"hasMore"`
		Success bool                      `json:"success"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &page2Response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), page2Response.Success)
	assert.Equal(suite.T(), 1, len(page2Response.Items))
	assert.False(suite.T(), page2Response.HasMore) // Should be last page
}

// FLOW TEST 9: Error Handling Flow
func (suite *StreamIntegrationTestSuite) TestErrorHandlingFlow() {
	// Test invalid category ID
	req := httptest.NewRequest("GET", "/api/v1/stream/categories/invalid", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Test invalid content ID
	req = httptest.NewRequest("GET", "/api/v1/stream/content/invalid", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Test missing required parameters
	req = httptest.NewRequest("GET", "/api/v1/stream/content", nil) // Missing categoryId
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// FLOW TEST 10: Performance Flow (Response Time)
func (suite *StreamIntegrationTestSuite) TestPerformanceFlow() {
	start := time.Now()

	// Test API response time
	req := httptest.NewRequest("GET", "/api/v1/stream/categories", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	duration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Less(suite.T(), duration, 200*time.Millisecond) // <200ms requirement
}

func TestStreamIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(StreamIntegrationTestSuite))
}