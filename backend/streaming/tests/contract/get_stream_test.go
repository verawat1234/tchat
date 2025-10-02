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

// TestGetStreamPactConsumer runs consumer Pact tests for GET /api/v1/streams/{streamId} endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for retrieving detailed stream information including metadata, viewer counts, and product features.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for GET /streams/{streamId} endpoint
// 2. No stream retrieval service logic exists
// 3. No product feature association is implemented for store streams
// 4. No analytics data aggregation exists
func TestGetStreamPactConsumer(t *testing.T) {
	// Create Pact consumer with mock provider configuration
	// Consumer: streaming-web-client (represents web frontend)
	// Provider: streaming-service (represents backend streaming service)
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "streaming-web-client",
		Provider: "streaming-service",
		Host:     "localhost",
		Port:     9001, // Same port as create_stream_test.go
		PactDir:  "./pacts", // Output directory for generated pact files
	})
	assert.NoError(t, err, "Failed to create Pact mock provider")

	// Test Case 1: Retrieve Store Stream with Product Details
	// Verifies that store streams return complete metadata including featured products,
	// viewer metrics, and commerce-specific fields
	t.Run("GetStoreStream_Success", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("store stream exists with featured products").
			UponReceiving("a request to retrieve store stream details").
			WithRequest("GET", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`), func(b *consumer.V2RequestBuilder) {
				// No request body for GET request
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with complete LiveStream schema
				b.JSONBody(map[string]interface{}{
					// Core identifiers
					"id":             matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"broadcaster_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

					// KYC and stream type
					"broadcaster_kyc_tier": matchers.Integer(1),
					"stream_type":          matchers.Like("store"),

					// Stream metadata
					"title":           matchers.Like("Product Launch Event"),
					"description":     matchers.Like("Join us for an exclusive look at our new collection"),
					"privacy_setting": matchers.Like("public"),
					"status":          matchers.Like("live"),

					// Scheduling and timing
					"scheduled_start_time": matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"actual_start_time":    matchers.Regex("2025-09-30T15:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"end_time":             nil, // Stream is still live

					// Recording details
					"recording_url":         nil, // Not available while live
					"recording_expiry_date": nil,

					// Viewer metrics
					"viewer_count":      matchers.Integer(1234),
					"peak_viewer_count": matchers.Integer(2500),
					"max_capacity":      matchers.Integer(50000),

					// Media
					"thumbnail_url": matchers.Regex("https://cdn.tchat.dev/thumbnails/550e8400-e29b-41d4-a716-446655440000.jpg", `^https://[a-z0-9\.\-]+/thumbnails/[0-9a-f\-]+\.(jpg|png|webp)$`),

					// Language and tags
					"language": matchers.Like("en"),
					"tags":     matchers.EachLike("commerce", 1),

					// Store-specific: Featured products (for store streams)
					"featured_products": matchers.EachLike(map[string]interface{}{
						"id":                   matchers.Regex("650e8400-e29b-41d4-a716-446655440100", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"product_id":           matchers.Regex("750e8400-e29b-41d4-a716-446655440200", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"product_name":         matchers.Like("Premium Widget"),
						"product_price":        matchers.Decimal(29.99),
						"product_image_url":    matchers.Regex("https://cdn.tchat.dev/products/750e8400.jpg", `^https://[a-z0-9\.\-]+/products/[0-9a-f\-]+\.(jpg|png|webp)$`),
						"display_position":     matchers.Like("overlay"),
						"display_priority":     matchers.Integer(5),
						"featured_at":          matchers.Regex("2025-09-30T15:10:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"view_count":           matchers.Integer(850),
						"click_count":          matchers.Integer(120),
						"purchase_count":       matchers.Integer(15),
						"revenue_generated":    matchers.Decimal(449.85),
					}, 1),

					// Timestamps
					"created_at": matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"updated_at": matchers.Regex("2025-09-30T15:20:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL with stream ID
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/550e8400-e29b-41d4-a716-446655440000", config.Host, config.Port)

				// Create HTTP GET request
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				// No authentication required for public stream retrieval
				req.Header.Set("Accept", "application/json")

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

				// Contract assertions for store stream
				assert.NotEmpty(t, responseBody["id"], "Stream ID should be present")
				assert.Equal(t, "store", responseBody["stream_type"], "Stream type should be 'store'")
				assert.Equal(t, "live", responseBody["status"], "Stream should be live")
				assert.Equal(t, float64(1), responseBody["broadcaster_kyc_tier"], "Store stream should have KYC Tier 1")

				// Verify viewer metrics
				assert.NotNil(t, responseBody["viewer_count"], "Viewer count should be present")
				assert.NotNil(t, responseBody["peak_viewer_count"], "Peak viewer count should be present")
				assert.Equal(t, float64(50000), responseBody["max_capacity"], "Max capacity should be 50000")

				// Verify store-specific featured products
				featuredProducts, ok := responseBody["featured_products"].([]interface{})
				assert.True(t, ok, "Featured products should be present for store streams")
				assert.NotEmpty(t, featuredProducts, "Store stream should have at least one featured product")

				if len(featuredProducts) > 0 {
					product := featuredProducts[0].(map[string]interface{})
					assert.NotEmpty(t, product["product_id"], "Product should have ID")
					assert.NotEmpty(t, product["product_name"], "Product should have name")
					assert.NotNil(t, product["product_price"], "Product should have price")
					assert.NotEmpty(t, product["display_position"], "Product should have display position")
					assert.NotNil(t, product["view_count"], "Product should track view count")
					assert.NotNil(t, product["revenue_generated"], "Product should track revenue")
				}

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Retrieve Video Stream (No Store-Specific Fields)
	// Verifies that video streams return standard fields without commerce data
	t.Run("GetVideoStream_Success", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("video stream exists without products").
			UponReceiving("a request to retrieve video stream details").
			WithRequest("GET", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440002", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`), func(b *consumer.V2RequestBuilder) {
				// No request body for GET request
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				// Video stream response without featured_products
				b.JSONBody(map[string]interface{}{
					"id":                   matchers.Regex("550e8400-e29b-41d4-a716-446655440002", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"broadcaster_id":       matchers.Regex("550e8400-e29b-41d4-a716-446655440003", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"broadcaster_kyc_tier": matchers.Integer(0), // Video streams don't require KYC
					"stream_type":          matchers.Like("video"),
					"title":                matchers.Like("Gaming Session"),
					"description":          matchers.Like("Playing the latest game release"),
					"privacy_setting":      matchers.Like("public"),
					"status":               matchers.Like("live"),
					"scheduled_start_time": matchers.Regex("2025-09-30T18:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"actual_start_time":    matchers.Regex("2025-09-30T18:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"end_time":             nil,
					"recording_url":        nil,
					"recording_expiry_date": nil,
					"viewer_count":         matchers.Integer(5678),
					"peak_viewer_count":    matchers.Integer(8900),
					"max_capacity":         matchers.Integer(50000),
					"thumbnail_url":        matchers.Regex("https://cdn.tchat.dev/thumbnails/550e8400-e29b-41d4-a716-446655440002.jpg", `^https://[a-z0-9\.\-]+/thumbnails/[0-9a-f\-]+\.(jpg|png|webp)$`),
					"language":             matchers.Like("en"),
					"tags":                 matchers.EachLike("gaming", 1),

					// No featured_products for video streams
					"featured_products": nil,

					"created_at": matchers.Regex("2025-09-30T17:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"updated_at": matchers.Regex("2025-09-30T18:15:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/550e8400-e29b-41d4-a716-446655440002", config.Host, config.Port)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Accept", "application/json")

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

				// Contract assertions for video stream
				assert.Equal(t, "video", responseBody["stream_type"], "Stream type should be 'video'")
				assert.Equal(t, float64(0), responseBody["broadcaster_kyc_tier"], "Video streams can have KYC Tier 0")
				assert.Equal(t, "live", responseBody["status"], "Stream should be live")

				// Verify no store-specific fields
				featuredProducts, _ := responseBody["featured_products"].([]interface{})
				assert.Empty(t, featuredProducts, "Video streams should not have featured products")

				// Verify viewer metrics still present
				assert.NotNil(t, responseBody["viewer_count"], "Viewer count should be present")
				assert.NotNil(t, responseBody["peak_viewer_count"], "Peak viewer count should be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Stream Not Found
	// Verifies proper 404 error handling for invalid stream IDs
	t.Run("GetStream_NotFound", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream does not exist").
			UponReceiving("a request for non-existent stream").
			WithRequest("GET", matchers.Regex("/api/v1/streams/00000000-0000-0000-0000-000000000000", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`), func(b *consumer.V2RequestBuilder) {
				// No request body for GET request
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"error": matchers.Like("Stream not found"),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Use invalid/non-existent stream ID
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/00000000-0000-0000-0000-000000000000", config.Host, config.Port)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Accept", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 404 Not Found
				assert.Equal(t, 404, resp.StatusCode, "Expected 404 Not Found for non-existent stream")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error message
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"], "not found", "Error should indicate stream not found")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Include Analytics Query Parameter
	// Verifies that analytics data can be optionally included with ?include_analytics=true
	t.Run("GetStream_IncludeAnalytics", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream exists with analytics data").
			UponReceiving("a request to retrieve stream with analytics included").
			WithRequest("GET", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440004?include_analytics=true", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\?include_analytics=true`), func(b *consumer.V2RequestBuilder) {
				b.Query("include_analytics", []string{"true"})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				// Response with additional analytics object
				b.JSONBody(map[string]interface{}{
					// Standard stream fields
					"id":                   matchers.Regex("550e8400-e29b-41d4-a716-446655440004", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"broadcaster_id":       matchers.Regex("550e8400-e29b-41d4-a716-446655440005", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"broadcaster_kyc_tier": matchers.Integer(1),
					"stream_type":          matchers.Like("store"),
					"title":                matchers.Like("Flash Sale Event"),
					"description":          matchers.Like("Limited time offers"),
					"privacy_setting":      matchers.Like("public"),
					"status":               matchers.Like("ended"),
					"scheduled_start_time": matchers.Regex("2025-09-30T10:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"actual_start_time":    matchers.Regex("2025-09-30T10:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"end_time":             matchers.Regex("2025-09-30T12:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"recording_url":        matchers.Regex("https://cdn.tchat.dev/recordings/550e8400-e29b-41d4-a716-446655440004.mp4", `^https://[a-z0-9\.\-]+/recordings/[0-9a-f\-]+\.mp4$`),
					"recording_expiry_date": matchers.Regex("2025-10-30T12:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"viewer_count":         matchers.Integer(0), // Ended stream
					"peak_viewer_count":    matchers.Integer(15000),
					"max_capacity":         matchers.Integer(50000),
					"thumbnail_url":        matchers.Regex("https://cdn.tchat.dev/thumbnails/550e8400-e29b-41d4-a716-446655440004.jpg", `^https://[a-z0-9\.\-]+/thumbnails/[0-9a-f\-]+\.(jpg|png|webp)$`),
					"language":             matchers.Like("en"),
					"tags":                 matchers.EachLike("flash-sale", 1),
					"created_at":           matchers.Regex("2025-09-30T09:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					"updated_at":           matchers.Regex("2025-09-30T12:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),

					// Additional analytics object (from OpenAPI StreamAnalytics schema)
					"analytics": map[string]interface{}{
						"stream_id":                    matchers.Regex("550e8400-e29b-41d4-a716-446655440004", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"total_unique_viewers":         matchers.Integer(12500),
						"peak_concurrent_viewers":      matchers.Integer(15000),
						"average_watch_duration_seconds": matchers.Integer(3600), // 1 hour average
						"total_chat_messages":          matchers.Integer(45000),
						"total_reactions":              matchers.Integer(120000),
						"unique_chatters":              matchers.Integer(8500),
						"products_featured":            matchers.Integer(8),
						"total_product_views":          matchers.Integer(95000),
						"total_product_clicks":         matchers.Integer(12000),
						"total_purchases":              matchers.Integer(2500),
						"total_revenue":                matchers.Decimal(125000.50),
						"average_viewer_quality":       matchers.Like("high"),
						"total_rebuffer_events":        matchers.Integer(150),
						"viewer_countries": map[string]interface{}{
							"TH": matchers.Integer(5000), // Thailand
							"SG": matchers.Integer(3500), // Singapore
							"ID": matchers.Integer(2000), // Indonesia
							"MY": matchers.Integer(1500), // Malaysia
							"VN": matchers.Integer(500),  // Vietnam
						},
						"calculated_at": matchers.Regex("2025-09-30T12:10:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Request with query parameter
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/550e8400-e29b-41d4-a716-446655440004?include_analytics=true", config.Host, config.Port)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Accept", "application/json")
				// Authenticated request for analytics
				req.Header.Set("Authorization", "Bearer test-token-broadcaster")

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

				// Verify standard stream fields
				assert.Equal(t, "store", responseBody["stream_type"], "Stream type should be 'store'")
				assert.Equal(t, "ended", responseBody["status"], "Stream should be ended")

				// Verify recording available for ended stream
				assert.NotNil(t, responseBody["recording_url"], "Recording should be available for ended stream")
				assert.NotNil(t, responseBody["recording_expiry_date"], "Recording expiry should be set")

				// Verify analytics object presence
				analytics, ok := responseBody["analytics"].(map[string]interface{})
				assert.True(t, ok, "Analytics object should be present when include_analytics=true")
				assert.NotNil(t, analytics, "Analytics should not be nil")

				// Verify key analytics metrics
				assert.NotNil(t, analytics["total_unique_viewers"], "Total unique viewers should be present")
				assert.NotNil(t, analytics["peak_concurrent_viewers"], "Peak concurrent viewers should be present")
				assert.NotNil(t, analytics["average_watch_duration_seconds"], "Average watch duration should be present")
				assert.NotNil(t, analytics["total_chat_messages"], "Total chat messages should be present")
				assert.NotNil(t, analytics["total_reactions"], "Total reactions should be present")

				// Verify store-specific analytics
				assert.NotNil(t, analytics["products_featured"], "Products featured count should be present")
				assert.NotNil(t, analytics["total_product_views"], "Total product views should be present")
				assert.NotNil(t, analytics["total_purchases"], "Total purchases should be present")
				assert.NotNil(t, analytics["total_revenue"], "Total revenue should be present")

				// Verify geographic distribution
				viewerCountries, ok := analytics["viewer_countries"].(map[string]interface{})
				assert.True(t, ok, "Viewer countries should be present")
				assert.NotEmpty(t, viewerCountries, "Viewer countries should have data")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}