package contract

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// StreamCategoryDetailTestSuite tests the specific Stream category retrieval API
type StreamCategoryDetailTestSuite struct {
	suite.Suite
	client  *http.Client
	baseURL string
}

// StreamCategoryDetailResponse represents the API response for a specific category
type StreamCategoryDetailResponse struct {
	Category StreamCategory `json:"category"`
	Success  bool           `json:"success"`
}

func (suite *StreamCategoryDetailTestSuite) SetupSuite() {
	suite.client = &http.Client{}
	suite.baseURL = "http://localhost:8083" // Commerce service port
}

// TestGetStreamCategoryDetail tests GET /api/v1/stream/categories/{categoryId}
func (suite *StreamCategoryDetailTestSuite) TestGetStreamCategoryDetail() {
	// This test MUST FAIL until the endpoint is implemented
	categoryID := "books"
	url := suite.baseURL + "/api/v1/stream/categories/" + categoryID

	resp, err := suite.client.Get(url)

	// Expected behavior: endpoint should exist and return proper response
	suite.Require().NoError(err, "Request should not fail")
	defer resp.Body.Close()

	// Expected: 200 OK status
	suite.Equal(http.StatusOK, resp.StatusCode, "Should return 200 OK")

	// Expected: proper Content-Type
	suite.Equal("application/json", resp.Header.Get("Content-Type"), "Should return JSON")

	// Expected: valid response structure
	var response StreamCategoryDetailResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: response should contain the specific category
	suite.True(response.Success, "Response should indicate success")
	suite.Equal(categoryID, response.Category.ID, "Should return the requested category")
	suite.NotEmpty(response.Category.Name, "Category name should not be empty")
	suite.GreaterOrEqual(response.Category.DisplayOrder, 0, "Display order should be non-negative")
}

// TestGetStreamCategoryDetailNotFound tests 404 handling
func (suite *StreamCategoryDetailTestSuite) TestGetStreamCategoryDetailNotFound() {
	// This test MUST FAIL until proper error handling is implemented
	url := suite.baseURL + "/api/v1/stream/categories/nonexistent"

	resp, err := suite.client.Get(url)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: 404 Not Found
	suite.Equal(http.StatusNotFound, resp.StatusCode, "Should return 404 for nonexistent category")
}

// TestGetStreamCategoryDetailPerformance tests response time requirements
func (suite *StreamCategoryDetailTestSuite) TestGetStreamCategoryDetailPerformance() {
	// This test MUST FAIL until performance requirements are met
	url := suite.baseURL + "/api/v1/stream/categories/books"

	start := time.Now()
	resp, err := suite.client.Get(url)
	elapsed := time.Since(start)

	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: <200ms response time
	suite.Less(elapsed.Milliseconds(), int64(200), "Response time should be less than 200ms")
}

func TestStreamCategoryDetailTestSuite(t *testing.T) {
	suite.Run(t, new(StreamCategoryDetailTestSuite))
}