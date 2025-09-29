package contract

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// StreamContentTestSuite tests the Stream content listing API
type StreamContentTestSuite struct {
	suite.Suite
	client  *http.Client
	baseURL string
}

// StreamContentItem represents content items from the API
type StreamContentItem struct {
	ID                 string  `json:"id"`
	CategoryID         string  `json:"categoryId"`
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	ThumbnailURL       string  `json:"thumbnailUrl"`
	ContentType        string  `json:"contentType"`
	Duration           *int    `json:"duration,omitempty"`
	Price              float64 `json:"price"`
	Currency           string  `json:"currency"`
	AvailabilityStatus string  `json:"availabilityStatus"`
	IsFeatured         bool    `json:"isFeatured"`
	FeaturedOrder      *int    `json:"featuredOrder,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

// StreamContentResponse represents the API response for content listing
type StreamContentResponse struct {
	Items   []StreamContentItem `json:"items"`
	Page    int                 `json:"page"`
	Limit   int                 `json:"limit"`
	Total   int                 `json:"total"`
	HasMore bool                `json:"hasMore"`
	Success bool                `json:"success"`
}

func (suite *StreamContentTestSuite) SetupSuite() {
	suite.client = &http.Client{}
	suite.baseURL = "http://localhost:8083" // Commerce service port
}

// TestGetStreamContent tests GET /api/v1/stream/content
func (suite *StreamContentTestSuite) TestGetStreamContent() {
	// This test MUST FAIL until the endpoint is implemented
	url := suite.baseURL + "/api/v1/stream/content"

	resp, err := suite.client.Get(url)

	// Expected behavior: endpoint should exist and return proper response
	suite.Require().NoError(err, "Request should not fail")
	defer resp.Body.Close()

	// Expected: 200 OK status
	suite.Equal(http.StatusOK, resp.StatusCode, "Should return 200 OK")

	// Expected: proper Content-Type
	suite.Equal("application/json", resp.Header.Get("Content-Type"), "Should return JSON")

	// Expected: valid response structure
	var response StreamContentResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: response should contain content items
	suite.True(response.Success, "Response should indicate success")
	suite.GreaterOrEqual(response.Total, 0, "Total should be non-negative")
	suite.GreaterOrEqual(response.Page, 1, "Page should be at least 1")
	suite.GreaterOrEqual(response.Limit, 1, "Limit should be at least 1")

	// Expected: each content item should have required fields
	for _, item := range response.Items {
		suite.NotEmpty(item.ID, "Content ID should not be empty")
		suite.NotEmpty(item.CategoryID, "Category ID should not be empty")
		suite.NotEmpty(item.Title, "Title should not be empty")
		suite.NotEmpty(item.Description, "Description should not be empty")
		suite.NotEmpty(item.ThumbnailURL, "Thumbnail URL should not be empty")
		suite.NotEmpty(item.ContentType, "Content type should not be empty")
		suite.GreaterOrEqual(item.Price, float64(0), "Price should be non-negative")
		suite.NotEmpty(item.Currency, "Currency should not be empty")
		suite.NotEmpty(item.AvailabilityStatus, "Availability status should not be empty")
	}
}

// TestGetStreamContentWithCategoryFilter tests category filtering
func (suite *StreamContentTestSuite) TestGetStreamContentWithCategoryFilter() {
	// This test MUST FAIL until filtering is implemented
	url := suite.baseURL + "/api/v1/stream/content?categoryId=books"

	resp, err := suite.client.Get(url)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: 200 OK status
	suite.Equal(http.StatusOK, resp.StatusCode, "Should return 200 OK")

	var response StreamContentResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: all items should belong to the books category
	for _, item := range response.Items {
		suite.Equal("books", item.CategoryID, "All items should belong to the books category")
	}
}

// TestGetStreamContentPagination tests pagination functionality
func (suite *StreamContentTestSuite) TestGetStreamContentPagination() {
	// This test MUST FAIL until pagination is implemented
	url := suite.baseURL + "/api/v1/stream/content?page=1&limit=5"

	resp, err := suite.client.Get(url)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var response StreamContentResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err, "Response should be valid JSON")

	// Expected: should respect pagination parameters
	suite.Equal(1, response.Page, "Should return requested page")
	suite.Equal(5, response.Limit, "Should respect limit parameter")
	suite.LessOrEqual(len(response.Items), 5, "Should not return more items than limit")
}

// TestGetStreamContentPerformance tests response time requirements
func (suite *StreamContentTestSuite) TestGetStreamContentPerformance() {
	// This test MUST FAIL until performance requirements are met
	url := suite.baseURL + "/api/v1/stream/content"

	start := time.Now()
	resp, err := suite.client.Get(url)
	elapsed := time.Since(start)

	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Expected: <200ms response time
	suite.Less(elapsed.Milliseconds(), int64(200), "Response time should be less than 200ms")
}

func TestStreamContentTestSuite(t *testing.T) {
	suite.Run(t, new(StreamContentTestSuite))
}