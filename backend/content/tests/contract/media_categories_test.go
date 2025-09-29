package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MediaCategoriesContractTestSuite defines the test suite for media categories contract
type MediaCategoriesContractTestSuite struct {
	suite.Suite
	router *gin.Engine
}

// SetupSuite runs before all tests in the suite
func (suite *MediaCategoriesContractTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	// Note: These tests are intentionally failing until implementation is complete
	// This follows TDD approach - write tests first, then implement
}

// TestGetMediaCategories tests GET /media/categories endpoint
func (suite *MediaCategoriesContractTestSuite) TestGetMediaCategories() {
	// Expected response structure based on OpenAPI spec
	expectedResponse := map[string]interface{}{
		"categories": []interface{}{},
		"total":      0,
	}

	// Make request to endpoint (this will fail until implementation exists)
	req, _ := http.NewRequest("GET", "/api/v1/media/categories", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract test: Verify response structure
	suite.T().Run("should return 200 status", func(t *testing.T) {
		// This will fail until handler is implemented
		assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK response")
	})

	suite.T().Run("should return categories array", func(t *testing.T) {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		// This will fail until handler is implemented
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Contains(t, response, "categories", "Response should contain categories field")
		assert.Contains(t, response, "total", "Response should contain total field")

		// Verify categories is an array
		categories, ok := response["categories"].([]interface{})
		assert.True(t, ok, "Categories should be an array")
		assert.NotNil(t, categories, "Categories array should not be nil")

		// Expected structure should match OpenAPI spec
		assert.Equal(t, expectedResponse["categories"], categories)
	})

	suite.T().Run("should return correct content type", func(t *testing.T) {
		// This will fail until handler is implemented
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})
}

// TestGetMediaCategory tests GET /media/categories/{categoryId} endpoint
func (suite *MediaCategoriesContractTestSuite) TestGetMediaCategory() {
	categoryID := "books"

	// Make request to endpoint (this will fail until implementation exists)
	req, _ := http.NewRequest("GET", "/api/v1/media/categories/"+categoryID, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.T().Run("should return 200 for valid category", func(t *testing.T) {
		// This will fail until handler is implemented
		assert.Equal(t, http.StatusOK, w.Code)
	})

	suite.T().Run("should return category object", func(t *testing.T) {
		var category map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &category)

		// This will fail until handler is implemented
		assert.NoError(t, err, "Response should be valid JSON")

		// Verify required fields according to OpenAPI spec
		assert.Contains(t, category, "id")
		assert.Contains(t, category, "name")
		assert.Contains(t, category, "displayOrder")
		assert.Contains(t, category, "iconName")
		assert.Contains(t, category, "isActive")
		assert.Contains(t, category, "featuredContentEnabled")
	})
}

// TestGetMediaCategoryNotFound tests 404 response for non-existent category
func (suite *MediaCategoriesContractTestSuite) TestGetMediaCategoryNotFound() {
	categoryID := "nonexistent"

	// Make request to endpoint
	req, _ := http.NewRequest("GET", "/api/v1/media/categories/"+categoryID, nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.T().Run("should return 404 for non-existent category", func(t *testing.T) {
		// This will fail until handler is implemented
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	suite.T().Run("should return error object", func(t *testing.T) {
		var errorResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)

		// This will fail until handler is implemented
		assert.NoError(t, err)
		assert.Contains(t, errorResponse, "error")
		assert.Contains(t, errorResponse, "message")
	})
}

// TestGetMovieSubtabs tests GET /media/movies/subtabs endpoint
func (suite *MediaCategoriesContractTestSuite) TestGetMovieSubtabs() {
	// Make request to endpoint
	req, _ := http.NewRequest("GET", "/api/v1/media/movies/subtabs", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.T().Run("should return 200 status", func(t *testing.T) {
		// This will fail until handler is implemented
		assert.Equal(t, http.StatusOK, w.Code)
	})

	suite.T().Run("should return subtabs array", func(t *testing.T) {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		// This will fail until handler is implemented
		assert.NoError(t, err)
		assert.Contains(t, response, "subtabs")
		assert.Contains(t, response, "defaultSubtab")

		subtabs, ok := response["subtabs"].([]interface{})
		assert.True(t, ok, "Subtabs should be an array")
		assert.NotNil(t, subtabs, "Subtabs array should not be nil")
	})
}

// TestInSuite runs all tests in the suite
func TestMediaCategoriesContractSuite(t *testing.T) {
	suite.Run(t, new(MediaCategoriesContractTestSuite))
}