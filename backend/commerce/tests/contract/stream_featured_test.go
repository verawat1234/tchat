package contract

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// StreamFeaturedTestSuite tests the featured Stream content API
type StreamFeaturedTestSuite struct {
	suite.Suite
	client  *http.Client
	baseURL string
}

// StreamFeaturedResponse represents the API response for featured content
type StreamFeaturedResponse struct {
	Items   []StreamContentItem `json:"items"`
	Total   int                 `json:"total"`
	HasMore bool                `json:"hasMore"`
	Success bool                `json:"success"`
}

func (suite *StreamFeaturedTestSuite) SetupSuite() {
	suite.client = &http.Client{}
	suite.baseURL = "http://localhost:8083" // Commerce service port
}

// TestGetStreamFeatured tests GET /api/v1/stream/featured
func (suite *StreamFeaturedTestSuite) TestGetStreamFeatured() {
	// This test MUST FAIL until the endpoint is implemented
	url := suite.baseURL + "/api/v1/stream/featured"

	resp, err := suite.client.Get(url)

	// Expected behavior: endpoint should exist and return proper response
	suite.Require().NoError(err, "Request should not fail")
	defer resp.Body.Close()

	// Expected: 200 OK status
	suite.Equal(http.StatusOK, resp.StatusCode, "Should return 200 OK")

	// Expected: proper Content-Type
	suite.Equal("application/json", resp.Header.Get("Content-Type"), "Should return JSON")

	// Expected: valid response structure
	var response StreamFeaturedResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: response should contain featured content
	suite.True(response.Success, "Response should indicate success")
	suite.GreaterOrEqual(response.Total, 0, "Total should be non-negative")

	// Expected: all items should be featured
	for _, item := range response.Items {
		suite.True(item.IsFeatured, "All returned items should be featured")
		suite.NotNil(item.FeaturedOrder, "Featured items should have featured order")
		suite.GreaterOrEqual(*item.FeaturedOrder, 0, "Featured order should be non-negative")
	}

	// Expected: items should be ordered by featured order
	for i := 1; i < len(response.Items); i++ {
		suite.LessOrEqual(*response.Items[i-1].FeaturedOrder, *response.Items[i].FeaturedOrder,
			"Featured items should be ordered by featured order")
	}
}

// TestGetStreamFeaturedByCategory tests category-specific featured content
func (suite *StreamFeaturedTestSuite) TestGetStreamFeaturedByCategory() {
	// This test MUST FAIL until category filtering is implemented
	url := suite.baseURL + "/api/v1/stream/featured?categoryId=books"

	resp, err := suite.client.Get(url)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: 200 OK status
	suite.Equal(http.StatusOK, resp.StatusCode, "Should return 200 OK")

	var response StreamFeaturedResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: all items should belong to the books category
	for _, item := range response.Items {
		suite.Equal("books", item.CategoryID, "All featured items should belong to the books category")
		suite.True(item.IsFeatured, "All returned items should be featured")
	}
}

// TestGetStreamFeaturedLimit tests limit parameter
func (suite *StreamFeaturedTestSuite) TestGetStreamFeaturedLimit() {
	// This test MUST FAIL until limit parameter is implemented
	url := suite.baseURL + "/api/v1/stream/featured?limit=3"

	resp, err := suite.client.Get(url)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var response StreamFeaturedResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: should respect limit parameter
	suite.LessOrEqual(len(response.Items), 3, "Should not return more items than limit")
}

// TestGetStreamFeaturedPerformance tests response time requirements
func (suite *StreamFeaturedTestSuite) TestGetStreamFeaturedPerformance() {
	// This test MUST FAIL until performance requirements are met
	url := suite.baseURL + "/api/v1/stream/featured"

	start := time.Now()
	resp, err := suite.client.Get(url)
	elapsed := time.Since(start)

	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: <200ms response time
	suite.Less(elapsed.Milliseconds(), int64(200), "Response time should be less than 200ms")
}

// TestGetStreamFeaturedAvailability tests that only available content is featured
func (suite *StreamFeaturedTestSuite) TestGetStreamFeaturedAvailability() {
	// This test MUST FAIL until availability filtering is implemented
	url := suite.baseURL + "/api/v1/stream/featured"

	resp, err := suite.client.Get(url)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var response StreamFeaturedResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: only available content should be featured
	for _, item := range response.Items {
		suite.Equal("available", item.AvailabilityStatus, "Only available content should be featured")
	}
}

func TestStreamFeaturedTestSuite(t *testing.T) {
	suite.Run(t, new(StreamFeaturedTestSuite))
}