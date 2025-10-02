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

// TestUpdateStreamPactConsumer runs consumer Pact tests for PATCH /api/v1/streams/{streamId} endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for updating existing live stream metadata.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for PATCH /streams/{streamId} endpoint
// 2. No ownership validation logic is implemented
// 3. No stream status checks (live stream update restrictions) exist
// 4. No partial update logic with field validation is implemented
func TestUpdateStreamPactConsumer(t *testing.T) {
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

	// Test Case 1: Update Stream Title and Description
	// Verifies that stream owners can update basic metadata fields
	// for scheduled streams before they go live
	t.Run("UpdateStream_Title", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440000"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("stream exists with ID 550e8400-e29b-41d4-a716-446655440000 and status is scheduled").
			UponReceiving("a request to update stream title and description").
			WithRequest("PATCH", fmt.Sprintf("/api/v1/streams/%s", streamID), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-owner"))
				b.JSONBody(map[string]interface{}{
					"title":       "Updated Product Launch Event",
					"description": "Updated description with more details about our new collection",
				})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with matchers
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Stream updated successfully"),
					"data": map[string]interface{}{
						// Core identifiers - use UUID regex for validation
						"id":             matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"broadcaster_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						// Updated fields
						"title":       matchers.Like("Updated Product Launch Event"),
						"description": matchers.Like("Updated description with more details about our new collection"),

						// Unchanged fields
						"broadcaster_kyc_tier": matchers.Integer(1),
						"stream_type":          matchers.Like("store"),
						"privacy_setting":      matchers.Like("public"),
						"status":               matchers.Like("scheduled"),
						"scheduled_start_time": matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"viewer_count":         matchers.Integer(0),
						"peak_viewer_count":    matchers.Integer(0),
						"max_capacity":         matchers.Integer(50000),
						"language":             matchers.Like("en"),

						// Timestamps - updated_at should be newer than created_at
						"created_at": matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"updated_at": matchers.Regex("2025-09-30T14:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s", config.Host, config.Port, streamID)

				// Prepare request payload
				updateStreamReq := map[string]interface{}{
					"title":       "Updated Product Launch Event",
					"description": "Updated description with more details about our new collection",
				}

				// Marshal request body to JSON
				jsonData, err := json.Marshal(updateStreamReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				// Create HTTP PATCH request
				req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				// Set required headers
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-owner")

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

				// Verify updated fields
				assert.Equal(t, "Updated Product Launch Event", data["title"], "Title should be updated")
				assert.Equal(t, "Updated description with more details about our new collection", data["description"], "Description should be updated")

				// Verify unchanged fields
				assert.Equal(t, streamID, data["id"], "Stream ID should remain unchanged")
				assert.Equal(t, "scheduled", data["status"], "Status should remain unchanged")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Update Privacy Setting
	// Verifies that stream owners can change privacy settings from public to private/unlisted
	t.Run("UpdateStream_Privacy", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440002"

		err := mockProvider.
			AddInteraction().
			Given("stream exists with ID 550e8400-e29b-41d4-a716-446655440002 and privacy is public").
			UponReceiving("a request to change privacy setting to private").
			WithRequest("PATCH", fmt.Sprintf("/api/v1/streams/%s", streamID), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-owner"))
				b.JSONBody(map[string]interface{}{
					"privacy_setting": "private",
				})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Stream updated successfully"),
					"data": map[string]interface{}{
						"id":                   matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"broadcaster_id":       matchers.Regex("550e8400-e29b-41d4-a716-446655440003", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"broadcaster_kyc_tier": matchers.Integer(1),
						"stream_type":          matchers.Like("store"),
						"title":                matchers.Like("Product Launch Event"),
						"description":          matchers.Like("Join us for an exclusive look"),

						// Updated field
						"privacy_setting": matchers.Like("private"),

						"status":               matchers.Like("scheduled"),
						"scheduled_start_time": matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"viewer_count":         matchers.Integer(0),
						"peak_viewer_count":    matchers.Integer(0),
						"max_capacity":         matchers.Integer(50000),
						"language":             matchers.Like("en"),
						"created_at":           matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"updated_at":           matchers.Regex("2025-09-30T14:35:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s", config.Host, config.Port, streamID)

				updateStreamReq := map[string]interface{}{
					"privacy_setting": "private",
				}

				jsonData, err := json.Marshal(updateStreamReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-owner")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify response status code
				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.True(t, responseBody["success"].(bool), "Response should indicate success")

				data := responseBody["data"].(map[string]interface{})
				assert.Equal(t, "private", data["privacy_setting"], "Privacy setting should be updated to private")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Non-Owner Attempts Update (403 Forbidden)
	// Verifies that only the stream owner can update stream metadata
	t.Run("UpdateStream_NotOwner", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440004"

		err := mockProvider.
			AddInteraction().
			Given("stream exists with ID 550e8400-e29b-41d4-a716-446655440004 owned by different user").
			UponReceiving("a request from non-owner to update stream").
			WithRequest("PATCH", fmt.Sprintf("/api/v1/streams/%s", streamID), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-not-owner"))
				b.JSONBody(map[string]interface{}{
					"title": "Unauthorized Update Attempt",
				})
			}).
			WillRespondWith(403, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Only stream broadcaster can perform this action"),
					"details": map[string]interface{}{
						"stream_id":     matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"required_role": matchers.Like("owner"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s", config.Host, config.Port, streamID)

				updateStreamReq := map[string]interface{}{
					"title": "Unauthorized Update Attempt",
				}

				jsonData, err := json.Marshal(updateStreamReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-not-owner")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 403 Forbidden
				assert.Equal(t, 403, resp.StatusCode, "Expected 403 Forbidden for non-owner")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "broadcaster", "Error should mention broadcaster requirement")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Update Active Live Stream (Restricted Updates)
	// Verifies that certain fields cannot be updated while stream is live
	t.Run("UpdateStream_ActiveStream", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440005"

		err := mockProvider.
			AddInteraction().
			Given("stream exists with ID 550e8400-e29b-41d4-a716-446655440005 and status is live").
			UponReceiving("a request to update restricted fields on live stream").
			WithRequest("PATCH", fmt.Sprintf("/api/v1/streams/%s", streamID), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-owner"))
				b.JSONBody(map[string]interface{}{
					"privacy_setting":      "private",
					"scheduled_start_time": "2025-09-30T20:00:00Z",
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Cannot update certain fields while stream is live"),
					"details": map[string]interface{}{
						"stream_id":     matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"current_status": matchers.Like("live"),
						"restricted_fields": matchers.EachLike("privacy_setting", 1),
						"allowed_fields":    matchers.EachLike("title", 1),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s", config.Host, config.Port, streamID)

				updateStreamReq := map[string]interface{}{
					"privacy_setting":      "private",
					"scheduled_start_time": "2025-09-30T20:00:00Z",
				}

				jsonData, err := json.Marshal(updateStreamReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-owner")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 400 Bad Request
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for restricted field update")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "live", "Error should mention stream is live")

				details := responseBody["details"].(map[string]interface{})
				assert.Equal(t, "live", details["current_status"], "Should indicate current status is live")
				assert.NotNil(t, details["restricted_fields"], "Should list restricted fields")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Update Thumbnail URL
	// Verifies that stream thumbnail can be updated with valid URL
	t.Run("UpdateStream_Thumbnail", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440006"

		err := mockProvider.
			AddInteraction().
			Given("stream exists with ID 550e8400-e29b-41d4-a716-446655440006 and status is scheduled").
			UponReceiving("a request to update stream thumbnail URL").
			WithRequest("PATCH", fmt.Sprintf("/api/v1/streams/%s", streamID), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-owner"))
				b.JSONBody(map[string]interface{}{
					"thumbnail_url": "https://cdn.tchat.dev/thumbnails/stream-thumb-12345.jpg",
				})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Stream updated successfully"),
					"data": map[string]interface{}{
						"id":                   matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"broadcaster_id":       matchers.Regex("550e8400-e29b-41d4-a716-446655440007", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"broadcaster_kyc_tier": matchers.Integer(1),
						"stream_type":          matchers.Like("store"),
						"title":                matchers.Like("Product Launch Event"),
						"description":          matchers.Like("Join us for an exclusive look"),
						"privacy_setting":      matchers.Like("public"),
						"status":               matchers.Like("scheduled"),

						// Updated field - thumbnail URL with regex validation
						"thumbnail_url": matchers.Regex("https://cdn.tchat.dev/thumbnails/stream-thumb-12345.jpg", `^https?://[a-z0-9\.\-]+/[a-zA-Z0-9\-\_/]+\.(jpg|jpeg|png|webp)$`),

						"scheduled_start_time": matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"viewer_count":         matchers.Integer(0),
						"peak_viewer_count":    matchers.Integer(0),
						"max_capacity":         matchers.Integer(50000),
						"language":             matchers.Like("en"),
						"created_at":           matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"updated_at":           matchers.Regex("2025-09-30T14:40:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s", config.Host, config.Port, streamID)

				updateStreamReq := map[string]interface{}{
					"thumbnail_url": "https://cdn.tchat.dev/thumbnails/stream-thumb-12345.jpg",
				}

				jsonData, err := json.Marshal(updateStreamReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-owner")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify response status code
				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.True(t, responseBody["success"].(bool), "Response should indicate success")

				data := responseBody["data"].(map[string]interface{})
				assert.Equal(t, "https://cdn.tchat.dev/thumbnails/stream-thumb-12345.jpg", data["thumbnail_url"], "Thumbnail URL should be updated")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 6: Stream Not Found (404)
	// Verifies proper error handling when updating non-existent stream
	t.Run("UpdateStream_NotFound", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-000000000000"

		err := mockProvider.
			AddInteraction().
			Given("no stream exists with ID 550e8400-e29b-41d4-a716-000000000000").
			UponReceiving("a request to update non-existent stream").
			WithRequest("PATCH", fmt.Sprintf("/api/v1/streams/%s", streamID), func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token"))
				b.JSONBody(map[string]interface{}{
					"title": "Update Non-Existent Stream",
				})
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Stream not found"),
					"details": map[string]interface{}{
						"stream_id": matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s", config.Host, config.Port, streamID)

				updateStreamReq := map[string]interface{}{
					"title": "Update Non-Existent Stream",
				}

				jsonData, err := json.Marshal(updateStreamReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token")

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

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "not found", "Error should mention stream not found")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}