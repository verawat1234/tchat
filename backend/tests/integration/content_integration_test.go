package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ContentIntegrationSuite tests the Content service endpoints
type ContentIntegrationSuite struct {
	APIIntegrationSuite
	ports ServicePort
}

// ContentItem represents a content item
type ContentItem struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Type        string            `json:"type"`
	Category    string            `json:"category"`
	Status      string            `json:"status"`
	AuthorID    string            `json:"authorId"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
	PublishedAt *string           `json:"publishedAt,omitempty"`
}

// Note: CreateContentRequest is now defined in types.go

// UpdateContentRequest represents content update request
type UpdateContentRequest struct {
	Title    *string           `json:"title,omitempty"`
	Content  *string           `json:"content,omitempty"`
	Category *string           `json:"category,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Note: ContentResponse is now defined in types.go

// SetupSuite initializes the content integration test suite
func (suite *ContentIntegrationSuite) SetupSuite() {
	suite.APIIntegrationSuite.SetupSuite()
	suite.ports = DefaultServicePorts()

	// Wait for content service to be available
	err := suite.waitForService(suite.ports.Content, 30*time.Second)
	if err != nil {
		suite.T().Fatalf("Content service not available: %v", err)
	}
}

// TestContentServiceHealth verifies content service health endpoint
func (suite *ContentIntegrationSuite) TestContentServiceHealth() {
	healthCheck, err := suite.checkServiceHealth(suite.ports.Content)
	require.NoError(suite.T(), err, "Health check should succeed")

	assert.Equal(suite.T(), "healthy", healthCheck.Status)
	assert.Equal(suite.T(), "content-service", healthCheck.Service)
	assert.NotEmpty(suite.T(), healthCheck.Timestamp)
}

// TestCreateContentItem tests content creation
func (suite *ContentIntegrationSuite) TestCreateContentItem() {
	url := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)

	createReq := CreateContentRequest{
		Title:    "Test Article",
		Content:  "This is a test article content",
		Type:     "article",
		Category: "technology",
		Tags:     []string{"test", "technology", "api"},
		Metadata: map[string]string{
			"author": "test-user",
			"source": "integration-test",
		},
	}

	resp, err := suite.makeRequest("POST", url, createReq, nil)
	require.NoError(suite.T(), err, "Create content request should succeed")
	defer resp.Body.Close()

	// Should return 201 for successful creation
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var contentResp ContentResponse
	err = suite.parseResponse(resp, &contentResp)
	require.NoError(suite.T(), err, "Should parse content creation response")

	assert.True(suite.T(), contentResp.Success)
	assert.Equal(suite.T(), "success", contentResp.Status)
	assert.NotNil(suite.T(), contentResp.Data)
	assert.NotEmpty(suite.T(), contentResp.Data.ID)
	assert.Equal(suite.T(), createReq.Title, contentResp.Data.Title)
	assert.Equal(suite.T(), createReq.Content, contentResp.Data.Content)
	assert.Equal(suite.T(), "draft", contentResp.Data.Status) // Default status
}

// TestGetContentItem tests retrieving a specific content item
func (suite *ContentIntegrationSuite) TestGetContentItem() {
	// First create a content item
	createURL := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)
	createReq := CreateContentRequest{
		Title:    "Test Get Content",
		Content:  "Content for get test",
		Type:     "article",
		Category: "test",
		Tags:     []string{"get-test"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create content for get test should succeed")
	defer createResp.Body.Close()

	var createContentResp ContentResponse
	err = suite.parseResponse(createResp, &createContentResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createContentResp.Data)

	contentID := createContentResp.Data.ID

	// Now get the content item
	getURL := fmt.Sprintf("%s:%d/api/v1/content/%s", suite.baseURL, suite.ports.Content, contentID)
	resp, err := suite.makeRequest("GET", getURL, nil, nil)
	require.NoError(suite.T(), err, "Get content request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var contentResp ContentResponse
	err = suite.parseResponse(resp, &contentResp)
	require.NoError(suite.T(), err, "Should parse get content response")

	assert.True(suite.T(), contentResp.Success)
	assert.NotNil(suite.T(), contentResp.Data)
	assert.Equal(suite.T(), contentID, contentResp.Data.ID)
	assert.Equal(suite.T(), createReq.Title, contentResp.Data.Title)
}

// TestGetNonExistentContent tests retrieving non-existent content
func (suite *ContentIntegrationSuite) TestGetNonExistentContent() {
	url := fmt.Sprintf("%s:%d/api/v1/content/non-existent-id", suite.baseURL, suite.ports.Content)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Get non-existent content should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// TestListContentItems tests listing content items
func (suite *ContentIntegrationSuite) TestListContentItems() {
	url := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List content request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var contentResp ContentResponse
	err = suite.parseResponse(resp, &contentResp)
	require.NoError(suite.T(), err, "Should parse list content response")

	assert.True(suite.T(), contentResp.Success)
	assert.NotNil(suite.T(), contentResp.Items)
	assert.GreaterOrEqual(suite.T(), contentResp.Total, 0)
}

// TestListContentWithPagination tests content listing with pagination
func (suite *ContentIntegrationSuite) TestListContentWithPagination() {
	url := fmt.Sprintf("%s:%d/api/v1/content?page=1&limit=5", suite.baseURL, suite.ports.Content)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Paginated list request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var contentResp ContentResponse
	err = suite.parseResponse(resp, &contentResp)
	require.NoError(suite.T(), err, "Should parse paginated response")

	assert.True(suite.T(), contentResp.Success)
	assert.Equal(suite.T(), 1, contentResp.Page)
	assert.Equal(suite.T(), 5, contentResp.Limit)
	assert.LessOrEqual(suite.T(), len(contentResp.Items), 5)
}

// TestListContentByCategory tests filtering content by category
func (suite *ContentIntegrationSuite) TestListContentByCategory() {
	// First create content with specific category
	createURL := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)
	createReq := CreateContentRequest{
		Title:    "Category Test Content",
		Content:  "Content for category filtering test",
		Type:     "article",
		Category: "test-category",
		Tags:     []string{"category-test"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create content for category test should succeed")
	createResp.Body.Close()

	// Now filter by category
	url := fmt.Sprintf("%s:%d/api/v1/content?category=test-category", suite.baseURL, suite.ports.Content)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Category filter request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var contentResp ContentResponse
	err = suite.parseResponse(resp, &contentResp)
	require.NoError(suite.T(), err, "Should parse category filter response")

	assert.True(suite.T(), contentResp.Success)

	// All returned items should have the specified category
	for _, item := range contentResp.Items {
		assert.Equal(suite.T(), "test-category", item.Category)
	}
}

// TestUpdateContentItem tests updating content
func (suite *ContentIntegrationSuite) TestUpdateContentItem() {
	// First create a content item
	createURL := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)
	createReq := CreateContentRequest{
		Title:    "Original Title",
		Content:  "Original content",
		Type:     "article",
		Category: "original",
		Tags:     []string{"original"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create content for update test should succeed")
	defer createResp.Body.Close()

	var createContentResp ContentResponse
	err = suite.parseResponse(createResp, &createContentResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createContentResp.Data)

	contentID := createContentResp.Data.ID

	// Now update the content
	updateURL := fmt.Sprintf("%s:%d/api/v1/content/%s", suite.baseURL, suite.ports.Content, contentID)
	newTitle := "Updated Title"
	newCategory := "updated"
	updateReq := UpdateContentRequest{
		Title:    &newTitle,
		Category: &newCategory,
		Tags:     []string{"updated", "test"},
	}

	resp, err := suite.makeRequest("PUT", updateURL, updateReq, nil)
	require.NoError(suite.T(), err, "Update content request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var contentResp ContentResponse
	err = suite.parseResponse(resp, &contentResp)
	require.NoError(suite.T(), err, "Should parse update response")

	assert.True(suite.T(), contentResp.Success)
	assert.NotNil(suite.T(), contentResp.Data)
	assert.Equal(suite.T(), contentID, contentResp.Data.ID)
	assert.Equal(suite.T(), newTitle, contentResp.Data.Title)
	assert.Equal(suite.T(), newCategory, contentResp.Data.Category)
	assert.Contains(suite.T(), contentResp.Data.Tags, "updated")
}

// TestDeleteContentItem tests content deletion
func (suite *ContentIntegrationSuite) TestDeleteContentItem() {
	// First create a content item
	createURL := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)
	createReq := CreateContentRequest{
		Title:    "Content to Delete",
		Content:  "This content will be deleted",
		Type:     "article",
		Category: "delete-test",
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create content for delete test should succeed")
	defer createResp.Body.Close()

	var createContentResp ContentResponse
	err = suite.parseResponse(createResp, &createContentResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createContentResp.Data)

	contentID := createContentResp.Data.ID

	// Now delete the content
	deleteURL := fmt.Sprintf("%s:%d/api/v1/content/%s", suite.baseURL, suite.ports.Content, contentID)
	resp, err := suite.makeRequest("DELETE", deleteURL, nil, nil)
	require.NoError(suite.T(), err, "Delete content request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify content is deleted by trying to get it
	getURL := fmt.Sprintf("%s:%d/api/v1/content/%s", suite.baseURL, suite.ports.Content, contentID)
	getResp, err := suite.makeRequest("GET", getURL, nil, nil)
	require.NoError(suite.T(), err, "Get deleted content should complete")
	defer getResp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, getResp.StatusCode)
}

// TestCreateContentInvalidData tests content creation with invalid data
func (suite *ContentIntegrationSuite) TestCreateContentInvalidData() {
	url := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)

	// Test with missing required fields
	invalidReq := CreateContentRequest{
		// Missing title and content
		Type:     "article",
		Category: "test",
	}

	resp, err := suite.makeRequest("POST", url, invalidReq, nil)
	require.NoError(suite.T(), err, "Invalid content creation should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// TestContentSearchByTags tests searching content by tags
func (suite *ContentIntegrationSuite) TestContentSearchByTags() {
	// First create content with specific tags
	createURL := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)
	createReq := CreateContentRequest{
		Title:    "Tagged Content",
		Content:  "Content with specific tags for search test",
		Type:     "article",
		Category: "search-test",
		Tags:     []string{"search", "test", "unique-tag"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create tagged content should succeed")
	createResp.Body.Close()

	// Search by tags
	url := fmt.Sprintf("%s:%d/api/v1/content?tags=unique-tag", suite.baseURL, suite.ports.Content)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Tag search request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var contentResp ContentResponse
	err = suite.parseResponse(resp, &contentResp)
	require.NoError(suite.T(), err, "Should parse tag search response")

	assert.True(suite.T(), contentResp.Success)

	// All returned items should contain the searched tag
	for _, item := range contentResp.Items {
		assert.Contains(suite.T(), item.Tags, "unique-tag")
	}
}

// TestInvalidHTTPMethods tests endpoints with invalid HTTP methods
func (suite *ContentIntegrationSuite) TestInvalidHTTPMethods() {
	baseURL := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)

	testCases := []struct {
		url    string
		method string
	}{
		{baseURL, "PATCH"},           // List endpoint with invalid method
		{baseURL + "/123", "POST"},   // Get endpoint with invalid method
		{baseURL + "/123", "PATCH"},  // Update endpoint with invalid method
	}

	for _, tc := range testCases {
		resp, err := suite.makeRequest(tc.method, tc.url, nil, nil)
		require.NoError(suite.T(), err, "Invalid method request should complete")
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusMethodNotAllowed, resp.StatusCode,
			"URL: %s, Method: %s", tc.url, tc.method)
	}
}

// TestContentServiceConcurrency tests concurrent requests to content service
func (suite *ContentIntegrationSuite) TestContentServiceConcurrency() {
	url := fmt.Sprintf("%s:%d/api/v1/content", suite.baseURL, suite.ports.Content)

	// Create 5 concurrent content creation requests
	concurrency := 5
	results := make(chan int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			createReq := CreateContentRequest{
				Title:    fmt.Sprintf("Concurrent Content %d", index),
				Content:  fmt.Sprintf("Content created concurrently #%d", index),
				Type:     "article",
				Category: "concurrency-test",
				Tags:     []string{"concurrent", fmt.Sprintf("test-%d", index)},
			}

			resp, err := suite.makeRequest("POST", url, createReq, nil)
			if err != nil {
				results <- 0
				return
			}
			defer resp.Body.Close()

			results <- resp.StatusCode
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		statusCode := <-results
		if statusCode == http.StatusCreated {
			successCount++
		}
	}

	// At least 80% of concurrent requests should succeed
	assert.GreaterOrEqual(suite.T(), successCount, 4, "Concurrent requests should mostly succeed")
}

// RunContentIntegrationTests runs the content integration test suite
func RunContentIntegrationTests(t *testing.T) {
	suite.Run(t, new(ContentIntegrationSuite))
}