package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

// TestFeatureProductPactConsumer runs consumer Pact tests for POST /api/v1/streams/{streamId}/products endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for featuring products in store live streams with real-time analytics tracking.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for POST /streams/{streamId}/products endpoint
// 2. No stream type validation logic (store vs video) is implemented
// 3. No product validation with commerce service integration exists
// 4. No featured product limit validation (max 10 products per stream)
// 5. No analytics tracking fields (view_count, click_count, purchase_count, revenue_generated) are implemented
func TestFeatureProductPactConsumer(t *testing.T) {
	// Create Pact consumer with mock provider configuration
	// Consumer: streaming-web-client (represents web frontend)
	// Provider: streaming-service (represents backend streaming service)
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "streaming-web-client",
		Provider: "streaming-service",
		Host:     "localhost",
		Port:     9001, // Same port as create_stream_test.go for provider consistency
		PactDir:  "./pacts", // Output directory for generated pact files
	})
	assert.NoError(t, err, "Failed to create Pact mock provider")

	// Test Case 1: Successful Product Feature in Store Stream
	// Verifies that products can be featured in store streams and receive proper analytics tracking
	// Reference: OpenAPI spec line 356-395, Data model line 154-191
	t.Run("FeatureProduct_Success", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440010"
		productID := "650e8400-e29b-41d4-a716-446655440020"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("stream exists with type 'store' and broadcaster owns the stream").
			UponReceiving("a request to feature a product in a store stream").
			WithRequest("POST", fmt.Sprintf("/api/v1/streams/%s/products", streamID), func(b *consumer.V2RequestBuilder) {
				b.JSONBody(map[string]interface{}{
					"product_id":        productID,
					"display_position":  "overlay",
					"display_priority":  5,
				})
			}).
			WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with matchers matching OpenAPI StreamProduct schema
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Product featured successfully"),
					"data": map[string]interface{}{
						// Core identifiers - UUID format validation
						"id":         matchers.Regex("750e8400-e29b-41d4-a716-446655440030", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"stream_id":  matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"product_id": matchers.Regex(productID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						// Feature metadata
						"featured_at":              matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"display_duration_seconds": nil, // NULL initially as still being featured
						"display_position":         matchers.Like("overlay"),
						"display_priority":         matchers.Integer(5),

						// Analytics tracking fields - initialized to 0 for new featured product
						"view_count":       matchers.Integer(0),
						"click_count":      matchers.Integer(0),
						"purchase_count":   matchers.Integer(0),
						"revenue_generated": matchers.Decimal(0.00),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products", config.Host, config.Port, streamID)

				// Prepare request payload
				featureProductReq := map[string]interface{}{
					"product_id":        productID,
					"display_position":  "overlay",
					"display_priority":  5,
				}

				// Marshal request body to JSON
				jsonData, err := json.Marshal(featureProductReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				// Create HTTP POST request
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				// Set required headers
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-broadcaster")

				// Execute the request
				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify response status code
				assert.Equal(t, 201, resp.StatusCode, "Expected 201 Created status code")

				// Parse response body
				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Additional contract assertions
				assert.True(t, responseBody["success"].(bool), "Response should indicate success")

				data, ok := responseBody["data"].(map[string]interface{})
				assert.True(t, ok, "Response should contain data object")

				// Verify critical StreamProduct fields are present
				assert.NotEmpty(t, data["id"], "Feature ID should be present")
				assert.Equal(t, streamID, data["stream_id"], "Stream ID should match request parameter")
				assert.Equal(t, productID, data["product_id"], "Product ID should match request body")
				assert.NotEmpty(t, data["featured_at"], "Featured timestamp should be present")
				assert.Equal(t, "overlay", data["display_position"], "Display position should match request")
				assert.Equal(t, float64(5), data["display_priority"], "Display priority should match request")

				// Verify analytics fields are initialized
				assert.Equal(t, float64(0), data["view_count"], "View count should be initialized to 0")
				assert.Equal(t, float64(0), data["click_count"], "Click count should be initialized to 0")
				assert.Equal(t, float64(0), data["purchase_count"], "Purchase count should be initialized to 0")
				assert.Equal(t, float64(0), data["revenue_generated"], "Revenue should be initialized to 0.00")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Cannot Feature Products in Video Streams
	// Verifies that the stream type validation prevents featuring products in video context
	// Reference: OpenAPI spec line 391 "Not a store stream", Data model line 155 "[store context only]"
	t.Run("FeatureProduct_VideoStream", func(t *testing.T) {
		videoStreamID := "550e8400-e29b-41d4-a716-446655440011"
		productID := "650e8400-e29b-41d4-a716-446655440021"

		// Define the expected interaction for video stream rejection
		err := mockProvider.
			AddInteraction().
			Given("stream exists with type 'video'").
			UponReceiving("a request to feature a product in a video stream").
			WithRequest("POST", fmt.Sprintf("/api/v1/streams/%s/products", videoStreamID), func(b *consumer.V2RequestBuilder) {
				b.JSONBody(map[string]interface{}{
					"product_id":        productID,
					"display_position":  "overlay",
					"display_priority":  3,
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Products can only be featured in store streams"),
					"details": map[string]interface{}{
						"stream_id":   matchers.Regex(videoStreamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"stream_type": matchers.Like("video"),
						"allowed_stream_types": []string{"store"},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products", config.Host, config.Port, videoStreamID)

				featureProductReq := map[string]interface{}{
					"product_id":        productID,
					"display_position":  "overlay",
					"display_priority":  3,
				}

				jsonData, err := json.Marshal(featureProductReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-broadcaster")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 400 Bad Request for video stream
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for video stream")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")

				details, ok := responseBody["details"].(map[string]interface{})
				assert.True(t, ok, "Details object should be present")
				assert.Equal(t, "video", details["stream_type"], "Should indicate stream is of type 'video'")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Product Not Found
	// Verifies that invalid product IDs are rejected with proper error response
	// This validates integration with commerce service for product existence
	// Reference: OpenAPI spec line 391 "product not found"
	t.Run("FeatureProduct_ProductNotFound", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440012"
		invalidProductID := "999e8400-e29b-41d4-a716-446655440000"

		err := mockProvider.
			AddInteraction().
			Given("stream exists with type 'store' but product does not exist in commerce service").
			UponReceiving("a request to feature a non-existent product").
			WithRequest("POST", fmt.Sprintf("/api/v1/streams/%s/products", streamID), func(b *consumer.V2RequestBuilder) {
				b.JSONBody(map[string]interface{}{
					"product_id":        invalidProductID,
					"display_position":  "sidebar",
					"display_priority":  8,
				})
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Product not found"),
					"details": map[string]interface{}{
						"product_id": matchers.Regex(invalidProductID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"message":    matchers.Like("Product does not exist or is not available"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products", config.Host, config.Port, streamID)

				featureProductReq := map[string]interface{}{
					"product_id":        invalidProductID,
					"display_position":  "sidebar",
					"display_priority":  8,
				}

				jsonData, err := json.Marshal(featureProductReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-broadcaster")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 404 Not Found for invalid product
				assert.Equal(t, 404, resp.StatusCode, "Expected 404 Not Found for invalid product")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.NotNil(t, responseBody["details"], "Details should be present for debugging")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Maximum Featured Products Limit
	// Verifies that streams cannot exceed 10 featured products limit
	// This is a business rule enforcement to prevent UI overload and maintain stream quality
	// Reference: Business logic requirement (not explicitly in OpenAPI but derived from commerce context)
	t.Run("FeatureProduct_MaxProducts", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440013"
		productID := "650e8400-e29b-41d4-a716-446655440023"

		err := mockProvider.
			AddInteraction().
			Given("stream already has 10 featured products").
			UponReceiving("a request to feature an 11th product").
			WithRequest("POST", fmt.Sprintf("/api/v1/streams/%s/products", streamID), func(b *consumer.V2RequestBuilder) {
				b.JSONBody(map[string]interface{}{
					"product_id":        productID,
					"display_position":  "fullscreen",
					"display_priority":  10,
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Maximum number of featured products reached"),
					"details": map[string]interface{}{
						"stream_id":             matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"current_product_count": matchers.Integer(10),
						"max_allowed":           matchers.Integer(10),
						"message":               matchers.Like("Remove a featured product before adding a new one"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products", config.Host, config.Port, streamID)

				featureProductReq := map[string]interface{}{
					"product_id":        productID,
					"display_position":  "fullscreen",
					"display_priority":  10,
				}

				jsonData, err := json.Marshal(featureProductReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-broadcaster")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 400 Bad Request for exceeding limit
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for exceeding product limit")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")

				details, ok := responseBody["details"].(map[string]interface{})
				assert.True(t, ok, "Details object should be present")
				assert.Equal(t, float64(10), details["current_product_count"], "Should indicate current count is 10")
				assert.Equal(t, float64(10), details["max_allowed"], "Should indicate maximum allowed is 10")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}