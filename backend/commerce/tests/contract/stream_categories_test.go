package contract

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// StreamCategoriesTestSuite tests the Stream categories API endpoint
type StreamCategoriesTestSuite struct {
	suite.Suite
	client *http.Client
	baseURL string
}

// StreamCategory represents the expected structure from backend
type StreamCategory struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	DisplayOrder            int    `json:"displayOrder"`
	IconName                string `json:"iconName"`
	IsActive                bool   `json:"isActive"`
	FeaturedContentEnabled  bool   `json:"featuredContentEnabled"`
	CreatedAt               string `json:"createdAt"`
	UpdatedAt               string `json:"updatedAt"`
}

// StreamCategoriesResponse represents the API response structure
type StreamCategoriesResponse struct {
	Categories []StreamCategory `json:"categories"`
	Total      int              `json:"total"`
	Success    bool             `json:"success"`
}

func (suite *StreamCategoriesTestSuite) SetupSuite() {
	suite.client = &http.Client{}
	suite.baseURL = "http://localhost:8083" // Commerce service port
}

// TestGetStreamCategories tests GET /api/v1/stream/categories
func (suite *StreamCategoriesTestSuite) TestGetStreamCategories() {
	// This test MUST FAIL until the endpoint is implemented
	url := suite.baseURL + "/api/v1/stream/categories"

	resp, err := suite.client.Get(url)

	// Expected behavior: endpoint should exist and return proper response
	suite.Require().NoError(err, "Request should not fail")
	defer resp.Body.Close()

	// Expected: 200 OK status
	suite.Equal(http.StatusOK, resp.Status, "Should return 200 OK")

	// Expected: proper Content-Type
	suite.Equal("application/json", resp.Header.Get("Content-Type"), "Should return JSON")

	// Expected: valid response structure
	var response StreamCategoriesResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: response should contain categories
	suite.True(response.Success, "Response should indicate success")
	suite.GreaterOrEqual(response.Total, 6, "Should return at least 6 categories (Books, Podcasts, Cartoons, Movies, Music, Art)")
	suite.Len(response.Categories, response.Total, "Categories count should match total")

	// Expected: categories should have required fields
	for _, category := range response.Categories {
		suite.NotEmpty(category.ID, "Category ID should not be empty")
		suite.NotEmpty(category.Name, "Category name should not be empty")
		suite.GreaterOrEqual(category.DisplayOrder, 0, "Display order should be non-negative")
		suite.NotEmpty(category.IconName, "Icon name should not be empty")
		suite.NotEmpty(category.CreatedAt, "Created at should not be empty")
		suite.NotEmpty(category.UpdatedAt, "Updated at should not be empty")
	}

	// Expected: categories should be ordered by displayOrder
	for i := 1; i < len(response.Categories); i++ {
		suite.LessOrEqual(response.Categories[i-1].DisplayOrder, response.Categories[i].DisplayOrder,
			"Categories should be ordered by display order")
	}
}

// TestGetStreamCategoriesAuth tests authentication requirements
func (suite *StreamCategoriesTestSuite) TestGetStreamCategoriesAuth() {
	// This test MUST FAIL until proper auth is implemented
	url := suite.baseURL + "/api/v1/stream/categories"

	req, err := http.NewRequest("GET", url, nil)
	suite.Require().NoError(err)

	// Request without authentication header
	resp, err := suite.client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: should work without auth (public endpoint) OR require auth
	// This will be determined during implementation
	suite.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized,
		"Should either allow public access or require authentication")
}

// TestGetStreamCategoriesPerformance tests response time requirements
func (suite *StreamCategoriesTestSuite) TestGetStreamCategoriesPerformance() {
	// This test MUST FAIL until performance requirements are met
	url := suite.baseURL + "/api/v1/stream/categories"

	// Measure response time
	start := time.Now()
	resp, err := suite.client.Get(url)
	elapsed := time.Since(start)

	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: <200ms response time
	suite.Less(elapsed.Milliseconds(), int64(200), "Response time should be less than 200ms")
}

// TestGetStreamCategoriesErrorHandling tests error scenarios
func (suite *StreamCategoriesTestSuite) TestGetStreamCategoriesErrorHandling() {
	// Test with invalid parameters
	url := suite.baseURL + "/api/v1/stream/categories?invalid=true"

	resp, err := suite.client.Get(url)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Should handle invalid parameters gracefully
	suite.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest,
		"Should handle invalid parameters gracefully")
}

func TestStreamCategoriesTestSuite(t *testing.T) {
	suite.Run(t, new(StreamCategoriesTestSuite))
}