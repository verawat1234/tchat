package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

// TestListProductsPactConsumer runs consumer Pact tests for GET /api/v1/streams/{streamId}/products endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for retrieving featured products in store live streams with optional analytics data.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for GET /streams/{streamId}/products endpoint
// 2. No database query logic for retrieving featured products by stream_id
// 3. No filtering logic for include_analytics and sort_by query parameters
// 4. No analytics aggregation logic for product performance metrics
// 5. No validation for stream_type (products only applicable to store streams)
func TestListProductsPactConsumer(t *testing.T) {
	// Create Pact consumer with mock provider configuration
	// Consumer: streaming-web-client (represents web frontend)
	// Provider: streaming-service (represents backend streaming service)
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "streaming-web-client",
		Provider: "streaming-service",
		Host:     "localhost",
		Port:     9001, // Same port as other contract tests
		PactDir:  "./pacts", // Output directory for generated pact files
	})
	assert.NoError(t, err, "Failed to create Pact mock provider")

	// Test Case 1: List Products for Store Stream
	// Verifies that clients can retrieve featured products sorted by display_order
	// for store context streams with proper product metadata
	t.Run("ListProducts_StoreStream", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440100"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("a store stream with multiple featured products exists").
			UponReceiving("a request to list featured products sorted by display order").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/products", streamID), func(b *consumer.V2RequestBuilder) {
				// No query parameters in this test case - default sorting by display_order
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with matchers
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Products retrieved successfully"),
					"data": map[string]interface{}{
						"products": matchers.EachLike(map[string]interface{}{
							// Core identifiers - use UUID regex for validation
							"id":         matchers.Regex("650e8400-e29b-41d4-a716-446655440200", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_id":  matchers.Like(streamID),
							"product_id": matchers.Regex("750e8400-e29b-41d4-a716-446655440300", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Feature timing metadata
							"featured_at":              matchers.Regex("2025-09-30T15:10:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"display_duration_seconds": matchers.Integer(120),

							// Display configuration
							"display_position": matchers.Like("overlay"),
							"display_priority": matchers.Integer(1),

							// Basic analytics (always included)
							"view_count":     matchers.Integer(450),
							"click_count":    matchers.Integer(85),
							"purchase_count": matchers.Integer(23),
							"revenue_generated": matchers.Decimal(1250.75),
						}, 3), // Expect at least 3 featured products
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products", config.Host, config.Port, streamID)

				// Create HTTP GET request
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				// Set required headers
				req.Header.Set("Content-Type", "application/json")

				// Execute the request
				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify response status code
				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				// Parse response body
				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Additional contract assertions
				assert.True(t, responseBody["success"].(bool), "Response should indicate success")

				data, ok := responseBody["data"].(map[string]interface{})
				assert.True(t, ok, "Response should contain data object")

				// Verify products array structure
				products := data["products"].([]interface{})
				assert.GreaterOrEqual(t, len(products), 3, "Should return at least 3 featured products")

				// Verify first product structure
				firstProduct := products[0].(map[string]interface{})
				assert.NotEmpty(t, firstProduct["id"], "Product feature ID should be present")
				assert.Equal(t, streamID, firstProduct["stream_id"], "Stream ID should match request")
				assert.NotEmpty(t, firstProduct["product_id"], "Product ID should be present")
				assert.NotEmpty(t, firstProduct["featured_at"], "Featured timestamp should be present")
				assert.NotEmpty(t, firstProduct["display_position"], "Display position should be present")
				assert.NotNil(t, firstProduct["view_count"], "View count should be present")
				assert.NotNil(t, firstProduct["click_count"], "Click count should be present")
				assert.NotNil(t, firstProduct["purchase_count"], "Purchase count should be present")
				assert.NotNil(t, firstProduct["revenue_generated"], "Revenue should be present")

				// Verify products are sorted by display_priority (ascending)
				if len(products) >= 2 {
					firstPriority := products[0].(map[string]interface{})["display_priority"].(float64)
					secondPriority := products[1].(map[string]interface{})["display_priority"].(float64)
					assert.LessOrEqual(t, firstPriority, secondPriority, "Products should be sorted by display_priority")
				}

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Include Analytics with Enhanced Metrics
	// Verifies that clients can request enhanced analytics data
	// including detailed performance metrics for each featured product
	t.Run("ListProducts_IncludeAnalytics", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440110"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("a store stream with products that have analytics data").
			UponReceiving("a request to list products with analytics included").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/products", streamID), func(b *consumer.V2RequestBuilder) {
				b.Query("include_analytics", matchers.String("true"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Products retrieved successfully"),
					"data": map[string]interface{}{
						"products": matchers.EachLike(map[string]interface{}{
							"id":         matchers.Regex("650e8400-e29b-41d4-a716-446655440210", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_id":  matchers.Like(streamID),
							"product_id": matchers.Regex("750e8400-e29b-41d4-a716-446655440310", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							"featured_at":              matchers.Regex("2025-09-30T16:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"display_duration_seconds": matchers.Integer(300),
							"display_position":         matchers.Like("fullscreen"),
							"display_priority":         matchers.Integer(0), // Highest priority

							// Enhanced analytics data
							"view_count":     matchers.Integer(1250),
							"click_count":    matchers.Integer(320),
							"purchase_count": matchers.Integer(85),
							"revenue_generated": matchers.Decimal(4567.89),

							// Additional analytics fields when include_analytics=true
							"analytics": map[string]interface{}{
								"conversion_rate":     matchers.Decimal(26.56), // (purchase_count / click_count) * 100
								"average_order_value": matchers.Decimal(53.74), // revenue_generated / purchase_count
								"click_through_rate":  matchers.Decimal(25.60), // (click_count / view_count) * 100
								"total_impressions":   matchers.Integer(1250),
								"unique_viewers":      matchers.Integer(980),
							},
						}, 2),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products?include_analytics=true", config.Host, config.Port, streamID)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.True(t, responseBody["success"].(bool), "Response should indicate success")

				data := responseBody["data"].(map[string]interface{})
				products := data["products"].([]interface{})

				// Verify analytics data is included
				for _, product := range products {
					productData := product.(map[string]interface{})
					analytics, ok := productData["analytics"].(map[string]interface{})
					assert.True(t, ok, "Analytics object should be present when include_analytics=true")

					// Verify analytics fields
					assert.NotNil(t, analytics["conversion_rate"], "Conversion rate should be present")
					assert.NotNil(t, analytics["average_order_value"], "Average order value should be present")
					assert.NotNil(t, analytics["click_through_rate"], "Click through rate should be present")
					assert.NotNil(t, analytics["total_impressions"], "Total impressions should be present")
					assert.NotNil(t, analytics["unique_viewers"], "Unique viewers should be present")
				}

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Video Stream Returns Empty Products
	// Verifies that video context streams (non-store) return empty products array
	// since product featuring is only applicable to store streams
	t.Run("ListProducts_VideoStream", func(t *testing.T) {
		videoStreamID := "550e8400-e29b-41d4-a716-446655440120"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("a video stream with no featured products").
			UponReceiving("a request to list products for a video stream").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/products", videoStreamID), func(b *consumer.V2RequestBuilder) {
				// No query parameters
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Products retrieved successfully"),
					"data": map[string]interface{}{
						"products": []interface{}{}, // Empty array for video streams
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products", config.Host, config.Port, videoStreamID)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.True(t, responseBody["success"].(bool), "Response should indicate success")

				data := responseBody["data"].(map[string]interface{})
				products := data["products"].([]interface{})

				// Verify empty products array for video streams
				assert.Empty(t, products, "Video streams should have no featured products")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Empty Result - Store Stream with No Products
	// Verifies that store streams with no featured products return empty array
	// maintaining consistent response structure
	t.Run("ListProducts_EmptyResult", func(t *testing.T) {
		emptyStreamID := "550e8400-e29b-41d4-a716-446655440130"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("a store stream with no featured products").
			UponReceiving("a request to list products for a store stream with no products").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/products", emptyStreamID), func(b *consumer.V2RequestBuilder) {
				// No query parameters
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Products retrieved successfully"),
					"data": map[string]interface{}{
						"products": []interface{}{}, // Empty array - no products featured yet
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products", config.Host, config.Port, emptyStreamID)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.True(t, responseBody["success"].(bool), "Response should indicate success even with empty results")

				data := responseBody["data"].(map[string]interface{})
				products := data["products"].([]interface{})

				// Verify empty products array
				assert.Empty(t, products, "Products array should be empty when no products featured")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Sort by Revenue
	// Verifies that clients can sort featured products by revenue_generated
	// to identify top-performing products during the stream
	t.Run("ListProducts_SortByRevenue", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440140"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("a store stream with multiple products having different revenues").
			UponReceiving("a request to list products sorted by revenue").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/products", streamID), func(b *consumer.V2RequestBuilder) {
				b.Query("sort_by", matchers.String("revenue"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Products retrieved successfully"),
					"data": map[string]interface{}{
						"products": []interface{}{
							// Products sorted by revenue_generated (descending)
							map[string]interface{}{
								"id":                       matchers.Regex("650e8400-e29b-41d4-a716-446655440220", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"stream_id":                matchers.Like(streamID),
								"product_id":               matchers.Regex("750e8400-e29b-41d4-a716-446655440320", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"featured_at":              matchers.Regex("2025-09-30T17:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
								"display_duration_seconds": matchers.Integer(180),
								"display_position":         matchers.Like("fullscreen"),
								"display_priority":         matchers.Integer(2),
								"view_count":               matchers.Integer(2000),
								"click_count":              matchers.Integer(450),
								"purchase_count":           matchers.Integer(120),
								"revenue_generated":        matchers.Decimal(8950.00), // Highest revenue
							},
							map[string]interface{}{
								"id":                       matchers.Regex("650e8400-e29b-41d4-a716-446655440221", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"stream_id":                matchers.Like(streamID),
								"product_id":               matchers.Regex("750e8400-e29b-41d4-a716-446655440321", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"featured_at":              matchers.Regex("2025-09-30T17:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
								"display_duration_seconds": matchers.Integer(150),
								"display_position":         matchers.Like("overlay"),
								"display_priority":         matchers.Integer(1),
								"view_count":               matchers.Integer(1500),
								"click_count":              matchers.Integer(280),
								"purchase_count":           matchers.Integer(65),
								"revenue_generated":        matchers.Decimal(4250.50), // Medium revenue
							},
							map[string]interface{}{
								"id":                       matchers.Regex("650e8400-e29b-41d4-a716-446655440222", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"stream_id":                matchers.Like(streamID),
								"product_id":               matchers.Regex("750e8400-e29b-41d4-a716-446655440322", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"featured_at":              matchers.Regex("2025-09-30T17:10:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
								"display_duration_seconds": matchers.Integer(120),
								"display_position":         matchers.Like("sidebar"),
								"display_priority":         matchers.Integer(3),
								"view_count":               matchers.Integer(800),
								"click_count":              matchers.Integer(120),
								"purchase_count":           matchers.Integer(28),
								"revenue_generated":        matchers.Decimal(1890.25), // Lowest revenue
							},
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/products?sort_by=revenue", config.Host, config.Port, streamID)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.True(t, responseBody["success"].(bool), "Response should indicate success")

				data := responseBody["data"].(map[string]interface{})
				products := data["products"].([]interface{})

				// Verify products are sorted by revenue (descending)
				assert.GreaterOrEqual(t, len(products), 3, "Should return at least 3 products")

				for i := 0; i < len(products)-1; i++ {
					currentRevenue := products[i].(map[string]interface{})["revenue_generated"].(float64)
					nextRevenue := products[i+1].(map[string]interface{})["revenue_generated"].(float64)
					assert.GreaterOrEqual(t, currentRevenue, nextRevenue, "Products should be sorted by revenue (descending)")
				}

				// Verify highest revenue product is first
				firstProduct := products[0].(map[string]interface{})
				assert.Greater(t, firstProduct["revenue_generated"].(float64), 8000.0, "First product should have highest revenue")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}