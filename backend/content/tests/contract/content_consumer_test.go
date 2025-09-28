package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

// ContentItem represents the content data structure for Pact tests
type ContentItem struct {
	ID        string                 `json:"id"`
	Category  string                 `json:"category"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Value     map[string]interface{} `json:"value"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	Notes     *string                `json:"notes,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// ContentCategory represents the content category structure
type ContentCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ContentVersion represents content version structure
type ContentVersion struct {
	ID        string                 `json:"id"`
	ContentID string                 `json:"content_id"`
	Version   int                    `json:"version"`
	Value     map[string]interface{} `json:"value"`
	CreatedAt time.Time              `json:"created_at"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse struct {
	Items      []ContentItem `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

// Test suite for Content Service Consumer Pact tests
func TestContentServiceConsumer(t *testing.T) {
	// Load test configuration
	testConfig, err := LoadTestConfig()
	assert.NoError(t, err)

	// Create Pact consumer with configured port
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: testConfig.WebClientName,
		Provider: testConfig.ContentServiceName,
		Host:     testConfig.MockServerHost,
		Port:     testConfig.ContentServiceMockPort,
		PactDir:  testConfig.PactDir,
	})
	assert.NoError(t, err)

	t.Run("Get Content Items - Success", func(t *testing.T) {
		// Define the expected interaction with correct v2 API
		err := mockProvider.
			AddInteraction().
			Given("content items exist").
			UponReceiving("a request for content items").
			WithRequest("GET", "/api/v1/content", func(b *consumer.V2RequestBuilder) {
				b.Query("page", matchers.String("1"))
				b.Query("per_page", matchers.String("10"))
				b.Header("Content-Type", matchers.String("application/json"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json; charset=utf-8"))
				b.JSONBody(matchers.StructMatcher{
					"items": matchers.EachLike(map[string]interface{}{
						"id":         matchers.String("550e8400-e29b-41d4-a716-446655440000"),
						"category":   matchers.String("test-category"),
						"type":       matchers.String("text"),
						"status":     matchers.String("published"),
						"value":      matchers.Like(map[string]interface{}{"content": "Sample content value"}),
						"created_at": matchers.Timestamp(),
						"updated_at": matchers.Timestamp(),
					}, 1),
					"total":       matchers.Integer(1),
					"page":        matchers.Integer(1),
					"per_page":    matchers.Integer(10),
					"total_pages": matchers.Integer(1),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Make the actual HTTP request using the correct host:port format
				client := &http.Client{}
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/api/v1/content?page=1&per_page=10", config.Host, config.Port), nil)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)
				return nil
			})

		assert.NoError(t, err)
	})

	t.Run("Get Single Content Item - Success", func(t *testing.T) {
		contentID := "550e8400-e29b-41d4-a716-446655440000"

		err := mockProvider.
			AddInteraction().
			Given("content item exists").
			UponReceiving("a request for a specific content item").
			WithRequest("GET", fmt.Sprintf("/api/v1/content/%s", contentID), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json; charset=utf-8"))
				b.JSONBody(matchers.StructMatcher{
					"id":         matchers.String(contentID),
					"category":   matchers.String("test-category"),
					"type":       matchers.String("text"),
					"status":     matchers.String("published"),
					"value":      matchers.Like(map[string]interface{}{"content": "Sample content value"}),
					"created_at": matchers.Timestamp(),
					"updated_at": matchers.Timestamp(),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := &http.Client{}
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/api/v1/content/%s", config.Host, config.Port, contentID), nil)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)
				return nil
			})
		assert.NoError(t, err)
	})

	t.Run("Create Content Item - Success", func(t *testing.T) {
		newContent := map[string]interface{}{
			"category": "test-category",
			"type":     "text",
			"value": map[string]interface{}{
				"content": "New content value",
			},
		}

		err := mockProvider.
			AddInteraction().
			Given("valid content data").
			UponReceiving("a request to create content").
			WithRequest("POST", "/api/v1/content", func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(matchers.Like(newContent))
			}).
			WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json; charset=utf-8"))
				b.JSONBody(matchers.StructMatcher{
					"id":         matchers.String("550e8400-e29b-41d4-a716-446655440001"),
					"category":   matchers.String("test-category"),
					"type":       matchers.String("text"),
					"status":     matchers.String("draft"),
					"value":      matchers.Like(map[string]interface{}{"content": "New content value"}),
					"created_at": matchers.Timestamp(),
					"updated_at": matchers.Timestamp(),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				body, _ := json.Marshal(newContent)
				client := &http.Client{}
				req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%d/api/v1/content", config.Host, config.Port), strings.NewReader(string(body)))
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, http.StatusCreated, resp.StatusCode)
				return nil
			})
		assert.NoError(t, err)
	})

	// Pact file writing is handled automatically in v2
}

// Test error scenarios
func TestContentServiceConsumerErrors(t *testing.T) {
	// Load test configuration
	testConfig, err := LoadTestConfig()
	assert.NoError(t, err)

	// Create Pact consumer for error scenarios
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: testConfig.WebClientName,
		Provider: testConfig.ContentServiceName,
		Host:     testConfig.MockServerHost,
		Port:     testConfig.ContentServiceMockPort,
		PactDir:  testConfig.PactDir,
	})
	assert.NoError(t, err)

	t.Run("Get Content Item - Not Found", func(t *testing.T) {
		contentID := "non-existent-id"

		err := mockProvider.
			AddInteraction().
			Given("content item does not exist").
			UponReceiving("a request for non-existent content").
			WithRequest("GET", fmt.Sprintf("/api/v1/content/%s", contentID), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json; charset=utf-8"))
				b.JSONBody(matchers.StructMatcher{
					"error":   matchers.String("Content not found"),
					"message": matchers.String("Content item with the specified ID does not exist"),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client := &http.Client{}
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/api/v1/content/%s", config.Host, config.Port, contentID), nil)
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
				return nil
			})
		assert.NoError(t, err)
	})

	t.Run("Create Content - Validation Error", func(t *testing.T) {
		invalidContent := map[string]interface{}{
			"category": "", // Empty category should cause validation error
			"type":     "invalid_type",
		}

		err := mockProvider.
			AddInteraction().
			Given("content validation rules are enforced").
			UponReceiving("a request to create invalid content").
			WithRequest("POST", "/api/v1/content", func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(matchers.Like(invalidContent))
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json; charset=utf-8"))
				b.JSONBody(matchers.StructMatcher{
					"error":   matchers.String("Validation failed"),
					"message": matchers.String("Invalid content data provided"),
					"details": matchers.EachLike(map[string]interface{}{
						"field":   matchers.String("category"),
						"message": matchers.String("Category is required"),
					}, 1),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				body, _ := json.Marshal(invalidContent)
				client := &http.Client{}
				req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%d/api/v1/content", config.Host, config.Port), strings.NewReader(string(body)))
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				return nil
			})
		assert.NoError(t, err)
	})

	// Pact file writing for error scenarios is handled automatically in v2
}