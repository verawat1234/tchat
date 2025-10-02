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

// TestGetStreamAnalyticsPactConsumer runs consumer Pact tests for GET /api/v1/streams/{streamId}/analytics endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for retrieving comprehensive stream analytics including viewer metrics, engagement, geography, and commerce data.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for GET /streams/{streamId}/analytics endpoint
// 2. No analytics aggregation service logic exists
// 3. No analytics calculation for viewer metrics, engagement, and geographic distribution
// 4. No commerce-specific analytics (products, revenue) for store streams
// 5. No broadcaster authorization validation for analytics access
func TestGetStreamAnalyticsPactConsumer(t *testing.T) {
	// Create Pact consumer with mock provider configuration
	// Consumer: streaming-web-client (represents web frontend)
	// Provider: streaming-service (represents backend streaming service)
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "streaming-web-client",
		Provider: "streaming-service",
		Host:     "localhost",
		Port:     9001, // Same port as other streaming contract tests
		PactDir:  "./pacts", // Output directory for generated pact files
	})
	assert.NoError(t, err, "Failed to create Pact mock provider")

	// Test Case 1: Retrieve Store Stream Analytics with Full Commerce Data
	// Verifies that store streams return complete analytics including:
	// - Viewer metrics (unique viewers, peak concurrent, watch duration)
	// - Engagement metrics (chat messages, reactions, unique chatters)
	// - Commerce metrics (products featured, views, clicks, purchases, revenue)
	// - Quality metrics (average quality, rebuffer events)
	// - Geographic distribution (viewer countries)
	t.Run("GetAnalytics_StoreStream", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("store stream exists with complete analytics data").
			UponReceiving("a request to retrieve store stream analytics").
			WithRequest("GET", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/analytics", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/analytics$`), func(b *consumer.V2RequestBuilder) {
				// Authorization required - only broadcaster can access analytics
				b.Header("Authorization", matchers.String("Bearer test-token-broadcaster"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with complete StreamAnalytics schema
				b.JSONBody(map[string]interface{}{
					// Core identifier
					"stream_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

					// Viewer metrics
					"total_unique_viewers":         matchers.Integer(12500),
					"peak_concurrent_viewers":      matchers.Integer(15000),
					"average_watch_duration_seconds": matchers.Integer(3600), // 1 hour average

					// Engagement metrics
					"total_chat_messages": matchers.Integer(45000),
					"total_reactions":     matchers.Integer(120000),
					"unique_chatters":     matchers.Integer(8500),

					// Store-specific commerce metrics
					"products_featured":    matchers.Integer(8),
					"total_product_views":  matchers.Integer(95000),
					"total_product_clicks": matchers.Integer(12000),
					"total_purchases":      matchers.Integer(2500),
					"total_revenue":        matchers.Decimal(125000.50),

					// Quality metrics
					"average_viewer_quality": matchers.Like("high"),
					"total_rebuffer_events":  matchers.Integer(150),

					// Geographic distribution (Southeast Asian focus)
					"viewer_countries": map[string]interface{}{
						"TH": matchers.Integer(5000), // Thailand
						"SG": matchers.Integer(3500), // Singapore
						"ID": matchers.Integer(2000), // Indonesia
						"MY": matchers.Integer(1500), // Malaysia
						"VN": matchers.Integer(500),  // Vietnam
					},

					// Timestamp
					"calculated_at": matchers.Regex("2025-09-30T12:10:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL with stream ID
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/analytics", config.Host, config.Port)

				// Create HTTP GET request
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				// Authorization required - only broadcaster can access analytics
				req.Header.Set("Authorization", "Bearer test-token-broadcaster")
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
				var analytics map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Contract assertions for store stream analytics
				assert.NotEmpty(t, analytics["stream_id"], "Stream ID should be present")

				// Verify viewer metrics
				assert.NotNil(t, analytics["total_unique_viewers"], "Total unique viewers should be present")
				assert.NotNil(t, analytics["peak_concurrent_viewers"], "Peak concurrent viewers should be present")
				assert.NotNil(t, analytics["average_watch_duration_seconds"], "Average watch duration should be present")

				// Verify engagement metrics
				assert.NotNil(t, analytics["total_chat_messages"], "Total chat messages should be present")
				assert.NotNil(t, analytics["total_reactions"], "Total reactions should be present")
				assert.NotNil(t, analytics["unique_chatters"], "Unique chatters should be present")

				// Verify store-specific commerce metrics
				assert.NotNil(t, analytics["products_featured"], "Products featured count should be present")
				assert.NotNil(t, analytics["total_product_views"], "Total product views should be present")
				assert.NotNil(t, analytics["total_product_clicks"], "Total product clicks should be present")
				assert.NotNil(t, analytics["total_purchases"], "Total purchases should be present")
				assert.NotNil(t, analytics["total_revenue"], "Total revenue should be present")

				// Verify quality metrics
				assert.NotEmpty(t, analytics["average_viewer_quality"], "Average viewer quality should be present")
				assert.NotNil(t, analytics["total_rebuffer_events"], "Total rebuffer events should be present")

				// Verify geographic distribution
				viewerCountries, ok := analytics["viewer_countries"].(map[string]interface{})
				assert.True(t, ok, "Viewer countries should be present")
				assert.NotEmpty(t, viewerCountries, "Viewer countries should have data")
				assert.NotNil(t, viewerCountries["TH"], "Thailand viewer count should be present")
				assert.NotNil(t, viewerCountries["SG"], "Singapore viewer count should be present")

				// Verify calculated timestamp
				assert.NotEmpty(t, analytics["calculated_at"], "Calculated timestamp should be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Retrieve Video Stream Analytics (No Commerce Fields)
	// Verifies that video streams return analytics without commerce-specific data
	// Commerce fields should be null for video streams
	t.Run("GetAnalytics_VideoStream", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("video stream exists with analytics but no commerce data").
			UponReceiving("a request to retrieve video stream analytics").
			WithRequest("GET", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440002/analytics", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/analytics$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-broadcaster-2"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				// Video stream analytics without commerce fields
				b.JSONBody(map[string]interface{}{
					"stream_id":                    matchers.Regex("550e8400-e29b-41d4-a716-446655440002", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"total_unique_viewers":         matchers.Integer(8900),
					"peak_concurrent_viewers":      matchers.Integer(10500),
					"average_watch_duration_seconds": matchers.Integer(2400), // 40 minutes average
					"total_chat_messages":          matchers.Integer(28000),
					"total_reactions":              matchers.Integer(85000),
					"unique_chatters":              matchers.Integer(5200),

					// Commerce fields should be null for video streams
					"products_featured":    nil,
					"total_product_views":  nil,
					"total_product_clicks": nil,
					"total_purchases":      nil,
					"total_revenue":        nil,

					"average_viewer_quality": matchers.Like("medium"),
					"total_rebuffer_events":  matchers.Integer(220),

					"viewer_countries": map[string]interface{}{
						"TH": matchers.Integer(3200),
						"SG": matchers.Integer(2400),
						"ID": matchers.Integer(1800),
						"MY": matchers.Integer(1000),
						"VN": matchers.Integer(500),
					},

					"calculated_at": matchers.Regex("2025-09-30T13:15:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/550e8400-e29b-41d4-a716-446655440002/analytics", config.Host, config.Port)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Authorization", "Bearer test-token-broadcaster-2")
				req.Header.Set("Accept", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var analytics map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Contract assertions for video stream analytics
				assert.NotEmpty(t, analytics["stream_id"], "Stream ID should be present")

				// Verify viewer and engagement metrics still present
				assert.NotNil(t, analytics["total_unique_viewers"], "Total unique viewers should be present")
				assert.NotNil(t, analytics["peak_concurrent_viewers"], "Peak concurrent viewers should be present")
				assert.NotNil(t, analytics["total_chat_messages"], "Total chat messages should be present")
				assert.NotNil(t, analytics["total_reactions"], "Total reactions should be present")

				// Verify commerce fields are null for video streams
				assert.Nil(t, analytics["products_featured"], "Products featured should be null for video streams")
				assert.Nil(t, analytics["total_product_views"], "Total product views should be null for video streams")
				assert.Nil(t, analytics["total_product_clicks"], "Total product clicks should be null for video streams")
				assert.Nil(t, analytics["total_purchases"], "Total purchases should be null for video streams")
				assert.Nil(t, analytics["total_revenue"], "Total revenue should be null for video streams")

				// Verify quality and geographic data still present
				assert.NotEmpty(t, analytics["average_viewer_quality"], "Average viewer quality should be present")
				assert.NotNil(t, analytics["total_rebuffer_events"], "Total rebuffer events should be present")

				viewerCountries, ok := analytics["viewer_countries"].(map[string]interface{})
				assert.True(t, ok, "Viewer countries should be present")
				assert.NotEmpty(t, viewerCountries, "Viewer countries should have data")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Analytics Access Forbidden (Non-Broadcaster)
	// Verifies that only the stream broadcaster can access analytics (403 Forbidden)
	t.Run("GetAnalytics_OnlyBroadcaster", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream exists but requester is not the broadcaster").
			UponReceiving("a request from non-broadcaster to access analytics").
			WithRequest("GET", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440003/analytics", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/analytics$`), func(b *consumer.V2RequestBuilder) {
				// Different user attempting to access analytics
				b.Header("Authorization", matchers.String("Bearer test-token-unauthorized-user"))
			}).
			WillRespondWith(403, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"error": matchers.Like("Only broadcaster can view analytics"),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/550e8400-e29b-41d4-a716-446655440003/analytics", config.Host, config.Port)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				// Unauthorized user token
				req.Header.Set("Authorization", "Bearer test-token-unauthorized-user")
				req.Header.Set("Accept", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 403 Forbidden
				assert.Equal(t, 403, resp.StatusCode, "Expected 403 Forbidden for non-broadcaster")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error message
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"], "broadcaster", "Error should mention broadcaster restriction")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Live Stream Analytics (Partial Data)
	// Verifies that analytics for currently live streams return partial/real-time data
	// Some metrics may be incomplete until stream ends
	t.Run("GetAnalytics_LiveStream", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("store stream is currently live with partial analytics").
			UponReceiving("a request to retrieve analytics for live stream").
			WithRequest("GET", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440004/analytics", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/analytics$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-broadcaster-3"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				// Live stream analytics with real-time data
				b.JSONBody(map[string]interface{}{
					"stream_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440004", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

					// Current real-time metrics
					"total_unique_viewers":    matchers.Integer(3200),
					"peak_concurrent_viewers": matchers.Integer(4500),
					// Average watch duration not yet meaningful for live stream
					"average_watch_duration_seconds": matchers.Integer(0),

					"total_chat_messages": matchers.Integer(12000),
					"total_reactions":     matchers.Integer(35000),
					"unique_chatters":     matchers.Integer(2800),

					// Commerce metrics for live store stream
					"products_featured":    matchers.Integer(3),
					"total_product_views":  matchers.Integer(28000),
					"total_product_clicks": matchers.Integer(3500),
					"total_purchases":      matchers.Integer(450),
					"total_revenue":        matchers.Decimal(22500.75),

					"average_viewer_quality": matchers.Like("high"),
					"total_rebuffer_events":  matchers.Integer(45),

					"viewer_countries": map[string]interface{}{
						"TH": matchers.Integer(1500),
						"SG": matchers.Integer(900),
						"ID": matchers.Integer(500),
						"MY": matchers.Integer(200),
						"VN": matchers.Integer(100),
					},

					// Real-time calculated timestamp
					"calculated_at": matchers.Regex("2025-09-30T15:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/550e8400-e29b-41d4-a716-446655440004/analytics", config.Host, config.Port)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Authorization", "Bearer test-token-broadcaster-3")
				req.Header.Set("Accept", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var analytics map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Contract assertions for live stream analytics
				assert.NotEmpty(t, analytics["stream_id"], "Stream ID should be present")

				// Verify real-time metrics are present
				assert.NotNil(t, analytics["total_unique_viewers"], "Real-time unique viewers should be present")
				assert.NotNil(t, analytics["peak_concurrent_viewers"], "Peak concurrent viewers should be present")

				// Average watch duration may be 0 or low for live streams
				watchDuration, ok := analytics["average_watch_duration_seconds"].(float64)
				assert.True(t, ok, "Average watch duration should be numeric")
				assert.True(t, watchDuration >= 0, "Average watch duration should be non-negative")

				// Verify engagement metrics are accumulating
				assert.NotNil(t, analytics["total_chat_messages"], "Chat messages should be accumulating")
				assert.NotNil(t, analytics["total_reactions"], "Reactions should be accumulating")

				// Verify commerce metrics are tracking in real-time
				assert.NotNil(t, analytics["products_featured"], "Products featured should be tracked")
				assert.NotNil(t, analytics["total_purchases"], "Purchases should be tracked in real-time")
				assert.NotNil(t, analytics["total_revenue"], "Revenue should be tracked in real-time")

				// Verify calculated timestamp is recent
				assert.NotEmpty(t, analytics["calculated_at"], "Calculated timestamp should be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Analytics Not Found
	// Verifies proper 404 error handling for non-existent streams
	t.Run("GetAnalytics_NotFound", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream does not exist").
			UponReceiving("a request for analytics of non-existent stream").
			WithRequest("GET", matchers.Regex("/api/v1/streams/00000000-0000-0000-0000-000000000000/analytics", `/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/analytics$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-broadcaster"))
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"error": matchers.Like("Stream not found"),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/00000000-0000-0000-0000-000000000000/analytics", config.Host, config.Port)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Authorization", "Bearer test-token-broadcaster")
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
}