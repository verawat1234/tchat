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

// TestListStreamsPactConsumer runs consumer Pact tests for GET /api/v1/streams endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for listing live streams with filtering and pagination capabilities.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for GET /streams endpoint
// 2. No filtering logic for status, stream_type, broadcaster_id is implemented
// 3. No pagination logic (limit, offset) is implemented
// 4. No database query layer exists for stream listing
func TestListStreamsPactConsumer(t *testing.T) {
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

	// Test Case 1: List All Active Live Streams
	// Verifies that clients can retrieve all currently live streams
	// with proper pagination and stream metadata
	t.Run("ListStreams_AllActive", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("multiple live streams exist").
			UponReceiving("a request to list all active live streams").
			WithRequest("GET", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.Query("status", matchers.String("live"))
				b.Query("limit", matchers.String("20"))
				b.Query("offset", matchers.String("0"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with matchers
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Streams retrieved successfully"),
					"data": map[string]interface{}{
						"streams": matchers.EachLike(map[string]interface{}{
							// Core identifiers - use UUID regex for validation
							"id":             matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"broadcaster_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// KYC and stream metadata
							"broadcaster_kyc_tier": matchers.Integer(1),
							"stream_type":          matchers.Like("store"),
							"title":                matchers.Like("Product Launch Event"),
							"description":          matchers.Like("Join us for an exclusive look at our new collection"),
							"privacy_setting":      matchers.Like("public"),
							"status":               matchers.Like("live"),

							// Stream metrics
							"viewer_count":      matchers.Integer(125),
							"peak_viewer_count": matchers.Integer(250),
							"max_capacity":      matchers.Integer(50000),

							// Timestamps
							"scheduled_start_time": matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"actual_start_time":    matchers.Regex("2025-09-30T15:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"created_at":           matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"updated_at":           matchers.Regex("2025-09-30T15:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),

							// Optional fields
							"thumbnail_url": matchers.Like("https://cdn.tchat.dev/thumbnails/stream123.jpg"),
							"language":      matchers.Like("en"),
							"tags":          matchers.EachLike("commerce", 1),
						}, 2), // Expect at least 2 streams in the list

						// Pagination metadata
						"total":  matchers.Integer(45),
						"limit":  matchers.Integer(20),
						"offset": matchers.Integer(0),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL with query parameters
				url := fmt.Sprintf("http://%s:%d/api/v1/streams?status=live&limit=20&offset=0", config.Host, config.Port)

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

				// Verify pagination metadata
				assert.NotNil(t, data["streams"], "Streams array should be present")
				assert.NotNil(t, data["total"], "Total count should be present")
				assert.Equal(t, float64(20), data["limit"], "Limit should match request")
				assert.Equal(t, float64(0), data["offset"], "Offset should match request")

				// Verify streams array structure
				streams := data["streams"].([]interface{})
				assert.GreaterOrEqual(t, len(streams), 2, "Should return at least 2 streams")

				// Verify first stream structure
				firstStream := streams[0].(map[string]interface{})
				assert.NotEmpty(t, firstStream["id"], "Stream ID should be present")
				assert.NotEmpty(t, firstStream["broadcaster_id"], "Broadcaster ID should be present")
				assert.Equal(t, "live", firstStream["status"], "All streams should have live status")
				assert.NotEmpty(t, firstStream["viewer_count"], "Viewer count should be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Filter Streams by Type (Store)
	// Verifies that clients can filter streams by stream_type
	// to retrieve only store/commerce streams
	t.Run("ListStreams_ByType", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("multiple store streams exist").
			UponReceiving("a request to list streams filtered by stream type").
			WithRequest("GET", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.Query("stream_type", matchers.String("store"))
				b.Query("limit", matchers.String("20"))
				b.Query("offset", matchers.String("0"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Streams retrieved successfully"),
					"data": map[string]interface{}{
						"streams": matchers.EachLike(map[string]interface{}{
							"id":                   matchers.Regex("550e8400-e29b-41d4-a716-446655440010", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"broadcaster_id":       matchers.Regex("550e8400-e29b-41d4-a716-446655440011", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"broadcaster_kyc_tier": matchers.Integer(1), // Store streams require KYC Tier 1
							"stream_type":          matchers.Like("store"),
							"title":                matchers.Like("Flash Sale Event"),
							"description":          matchers.Like("Limited time offers on best sellers"),
							"privacy_setting":      matchers.Like("public"),
							"status":               matchers.Like("live"),
							"viewer_count":         matchers.Integer(500),
							"peak_viewer_count":    matchers.Integer(750),
							"max_capacity":         matchers.Integer(50000),
							"scheduled_start_time": matchers.Regex("2025-09-30T16:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"actual_start_time":    matchers.Regex("2025-09-30T16:02:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"created_at":           matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"updated_at":           matchers.Regex("2025-09-30T16:02:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"thumbnail_url":        matchers.Like("https://cdn.tchat.dev/thumbnails/store456.jpg"),
							"language":             matchers.Like("en"),
							"tags":                 matchers.EachLike("commerce", 1),
						}, 1),

						"total":  matchers.Integer(15),
						"limit":  matchers.Integer(20),
						"offset": matchers.Integer(0),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams?stream_type=store&limit=20&offset=0", config.Host, config.Port)

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
				streams := data["streams"].([]interface{})

				// Verify all returned streams are of type "store"
				for _, stream := range streams {
					streamData := stream.(map[string]interface{})
					assert.Equal(t, "store", streamData["stream_type"], "All streams should be of type 'store'")
					assert.GreaterOrEqual(t, streamData["broadcaster_kyc_tier"].(float64), float64(1), "Store streams should have KYC Tier >= 1")
				}

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Filter Streams by Broadcaster ID
	// Verifies that clients can retrieve all streams from a specific broadcaster
	// useful for user profile pages and broadcaster dashboards
	t.Run("ListStreams_ByBroadcaster", func(t *testing.T) {
		broadcasterID := "550e8400-e29b-41d4-a716-446655440020"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("broadcaster has multiple streams").
			UponReceiving("a request to list streams by broadcaster ID").
			WithRequest("GET", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.Query("broadcaster_id", matchers.String(broadcasterID))
				b.Query("limit", matchers.String("20"))
				b.Query("offset", matchers.String("0"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Streams retrieved successfully"),
					"data": map[string]interface{}{
						"streams": matchers.EachLike(map[string]interface{}{
							"id": matchers.Regex("550e8400-e29b-41d4-a716-446655440021", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							// All streams should have the same broadcaster_id
							"broadcaster_id":       matchers.Like(broadcasterID),
							"broadcaster_kyc_tier": matchers.Integer(2),
							"stream_type":          matchers.Like("video"),
							"title":                matchers.Like("Weekly Gaming Stream"),
							"description":          matchers.Like("Join me for gaming and giveaways"),
							"privacy_setting":      matchers.Like("public"),
							"status":               matchers.Like("scheduled"),
							"viewer_count":         matchers.Integer(0),
							"peak_viewer_count":    matchers.Integer(0),
							"max_capacity":         matchers.Integer(50000),
							"scheduled_start_time": matchers.Regex("2025-10-01T18:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"created_at":           matchers.Regex("2025-09-30T10:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"updated_at":           matchers.Regex("2025-09-30T10:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"thumbnail_url":        matchers.Like("https://cdn.tchat.dev/thumbnails/gaming789.jpg"),
							"language":             matchers.Like("en"),
							"tags":                 matchers.EachLike("gaming", 1),
						}, 3), // Expect at least 3 streams from this broadcaster

						"total":  matchers.Integer(3),
						"limit":  matchers.Integer(20),
						"offset": matchers.Integer(0),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams?broadcaster_id=%s&limit=20&offset=0",
					config.Host, config.Port, broadcasterID)

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
				streams := data["streams"].([]interface{})

				// Verify all returned streams belong to the specified broadcaster
				for _, stream := range streams {
					streamData := stream.(map[string]interface{})
					assert.Equal(t, broadcasterID, streamData["broadcaster_id"], "All streams should belong to the specified broadcaster")
				}

				assert.Equal(t, float64(3), data["total"], "Total should match expected stream count")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Empty Result Set
	// Verifies that the API properly handles cases where no streams match the filter criteria
	// returning an empty array with valid pagination metadata
	t.Run("ListStreams_EmptyResult", func(t *testing.T) {
		// Define the expected interaction for empty results
		err := mockProvider.
			AddInteraction().
			Given("no streams match the filter criteria").
			UponReceiving("a request that returns no matching streams").
			WithRequest("GET", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.Query("status", matchers.String("live"))
				b.Query("stream_type", matchers.String("store"))
				b.Query("broadcaster_id", matchers.String("00000000-0000-0000-0000-000000000000"))
				b.Query("limit", matchers.String("20"))
				b.Query("offset", matchers.String("0"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Streams retrieved successfully"),
					"data": map[string]interface{}{
						"streams": []interface{}{}, // Empty array
						"total":   matchers.Integer(0),
						"limit":   matchers.Integer(20),
						"offset":  matchers.Integer(0),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams?status=live&stream_type=store&broadcaster_id=00000000-0000-0000-0000-000000000000&limit=20&offset=0",
					config.Host, config.Port)

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

				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK for empty results")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.True(t, responseBody["success"].(bool), "Response should indicate success even with empty results")

				data := responseBody["data"].(map[string]interface{})
				streams := data["streams"].([]interface{})

				// Verify empty results
				assert.Empty(t, streams, "Streams array should be empty")
				assert.Equal(t, float64(0), data["total"], "Total should be 0 for empty results")
				assert.Equal(t, float64(20), data["limit"], "Limit should still be present")
				assert.Equal(t, float64(0), data["offset"], "Offset should still be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Pagination - Second Page
	// Verifies that pagination works correctly with offset parameter
	// ensuring proper navigation through large result sets
	t.Run("ListStreams_Pagination", func(t *testing.T) {
		// Define the expected interaction for second page
		err := mockProvider.
			AddInteraction().
			Given("more than 20 streams exist").
			UponReceiving("a request for the second page of streams").
			WithRequest("GET", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.Query("limit", matchers.String("20"))
				b.Query("offset", matchers.String("20"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Streams retrieved successfully"),
					"data": map[string]interface{}{
						"streams": matchers.EachLike(map[string]interface{}{
							"id":                   matchers.Regex("550e8400-e29b-41d4-a716-446655440030", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"broadcaster_id":       matchers.Regex("550e8400-e29b-41d4-a716-446655440031", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"broadcaster_kyc_tier": matchers.Integer(1),
							"stream_type":          matchers.Like("video"),
							"title":                matchers.Like("Music Performance"),
							"description":          matchers.Like("Live concert streaming"),
							"privacy_setting":      matchers.Like("public"),
							"status":               matchers.Like("ended"),
							"viewer_count":         matchers.Integer(0),
							"peak_viewer_count":    matchers.Integer(3500),
							"max_capacity":         matchers.Integer(50000),
							"scheduled_start_time": matchers.Regex("2025-09-29T20:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"actual_start_time":    matchers.Regex("2025-09-29T20:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"end_time":             matchers.Regex("2025-09-29T22:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"recording_url":        matchers.Like("https://cdn.tchat.dev/recordings/stream456.mp4"),
							"recording_expiry_date": matchers.Regex("2025-10-29T22:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"created_at":           matchers.Regex("2025-09-29T19:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"updated_at":           matchers.Regex("2025-09-29T22:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"thumbnail_url":        matchers.Like("https://cdn.tchat.dev/thumbnails/concert123.jpg"),
							"language":             matchers.Like("en"),
							"tags":                 matchers.EachLike("music", 1),
						}, 10), // Second page has 10 streams remaining

						"total":  matchers.Integer(30),
						"limit":  matchers.Integer(20),
						"offset": matchers.Integer(20),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams?limit=20&offset=20", config.Host, config.Port)

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

				// Verify pagination metadata
				assert.Equal(t, float64(30), data["total"], "Total should be 30")
				assert.Equal(t, float64(20), data["limit"], "Limit should be 20")
				assert.Equal(t, float64(20), data["offset"], "Offset should be 20 (second page)")

				streams := data["streams"].([]interface{})
				assert.Equal(t, 10, len(streams), "Second page should have 10 streams (30 total - 20 offset)")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}