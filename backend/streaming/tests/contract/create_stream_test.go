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

// TestCreateStreamPactConsumer runs consumer Pact tests for POST /api/v1/streams endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for creating new live streams with KYC validation requirements.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for POST /streams endpoint
// 2. No KYC validation logic is implemented
// 3. No stream_key generation service exists
// 4. No WebRTC session management is implemented
func TestCreateStreamPactConsumer(t *testing.T) {
	// Create Pact consumer with mock provider configuration
	// Consumer: streaming-web-client (represents web frontend)
	// Provider: streaming-service (represents backend streaming service)
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "streaming-web-client",
		Provider: "streaming-service",
		Host:     "localhost",
		Port:     9001, // Different port from video service (9000)
		PactDir:  "./pacts", // Output directory for generated pact files
	})
	assert.NoError(t, err, "Failed to create Pact mock provider")

	// Test Case 1: Successful Store Stream Creation with KYC Tier 1
	// Verifies that users with Standard KYC (Tier 1) can create store streams
	// and receive proper RTMP stream credentials
	t.Run("CreateStoreStream_Success", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("broadcaster has KYC Tier 1").
			UponReceiving("a request to create a store stream with valid KYC").
			WithRequest("POST", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.JSONBody(map[string]interface{}{
					"stream_type":          "store",
					"title":                "Product Launch Event",
					"description":          "Join us for an exclusive look at our new collection",
					"privacy_setting":      "public",
					"scheduled_start_time": "2025-09-30T15:00:00Z",
					"language":             "en",
					"tags":                 []string{"commerce", "product-launch"},
				})
			}).
			WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with matchers
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Stream created successfully"),
					"data": map[string]interface{}{
						// Core identifiers - use UUID regex for validation
						"id":             matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"broadcaster_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						// KYC and stream metadata
						"broadcaster_kyc_tier": matchers.Integer(1),
						"stream_type":          matchers.Like("store"),
						"title":                matchers.Like("Product Launch Event"),
						"description":          matchers.Like("Join us for an exclusive look at our new collection"),
						"privacy_setting":      matchers.Like("public"),
						"status":               matchers.Like("scheduled"),

						// Critical streaming credentials
						"stream_key":        matchers.Regex("rtmp://stream.tchat.dev/live/sk_abc123xyz789", `^rtmp://[a-z0-9\.\-]+/live/sk_[a-zA-Z0-9]+$`),
						"webrtc_session_id": matchers.Regex("wrtc_session_abc123", `^wrtc_session_[a-zA-Z0-9]+$`),

						// Stream configuration
						"scheduled_start_time": matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"viewer_count":         matchers.Integer(0),
						"peak_viewer_count":    matchers.Integer(0),
						"max_capacity":         matchers.Integer(50000),
						"language":             matchers.Like("en"),
						"tags":                 matchers.EachLike("commerce", 1),

						// Timestamps
						"created_at": matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"updated_at": matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL
				url := fmt.Sprintf("http://%s:%d/api/v1/streams", config.Host, config.Port)

				// Prepare request payload
				createStreamReq := map[string]interface{}{
					"stream_type":          "store",
					"title":                "Product Launch Event",
					"description":          "Join us for an exclusive look at our new collection",
					"privacy_setting":      "public",
					"scheduled_start_time": "2025-09-30T15:00:00Z",
					"language":             "en",
					"tags":                 []string{"commerce", "product-launch"},
				}

				// Marshal request body to JSON
				jsonData, err := json.Marshal(createStreamReq)
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
				req.Header.Set("Authorization", "Bearer test-token-kyc-tier-1")

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

				// Verify critical streaming fields are present
				assert.NotEmpty(t, data["id"], "Stream ID should be present")
				assert.NotEmpty(t, data["stream_key"], "Stream key should be present for RTMP streaming")
				assert.NotEmpty(t, data["webrtc_session_id"], "WebRTC session ID should be present")
				assert.Equal(t, "store", data["stream_type"], "Stream type should be 'store'")
				assert.Equal(t, "scheduled", data["status"], "Initial status should be 'scheduled'")
				assert.Equal(t, float64(50000), data["max_capacity"], "Max capacity should be 50000")
				assert.Equal(t, float64(1), data["broadcaster_kyc_tier"], "KYC tier should be 1")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Insufficient KYC Verification (Tier 0)
	// Verifies that users without Standard KYC (Tier < 1) cannot create store streams
	// and receive proper error response with actionable guidance
	t.Run("CreateStoreStream_InsufficientKYC", func(t *testing.T) {
		// Define the expected interaction for KYC rejection
		err := mockProvider.
			AddInteraction().
			Given("broadcaster has KYC Tier 0").
			UponReceiving("a request to create a store stream without sufficient KYC").
			WithRequest("POST", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.JSONBody(map[string]interface{}{
					"stream_type":          "store",
					"title":                "Product Launch Event",
					"description":          "Join us for an exclusive look at our new collection",
					"privacy_setting":      "public",
					"scheduled_start_time": "2025-09-30T15:00:00Z",
				})
			}).
			WillRespondWith(403, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success":           false,
					"error":             matchers.Like("Store sellers require Standard KYC (Tier 1) verification"),
					"required_kyc_tier": matchers.Integer(1),
					"current_kyc_tier":  matchers.Integer(0),
					"verification_url":  matchers.Like("https://tchat.dev/verify/kyc"),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams", config.Host, config.Port)

				createStreamReq := map[string]interface{}{
					"stream_type":          "store",
					"title":                "Product Launch Event",
					"description":          "Join us for an exclusive look at our new collection",
					"privacy_setting":      "public",
					"scheduled_start_time": "2025-09-30T15:00:00Z",
				}

				jsonData, err := json.Marshal(createStreamReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-kyc-tier-0")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 403 Forbidden
				assert.Equal(t, 403, resp.StatusCode, "Expected 403 Forbidden for insufficient KYC")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Equal(t, float64(1), responseBody["required_kyc_tier"], "Should indicate required KYC tier is 1")
				assert.Equal(t, float64(0), responseBody["current_kyc_tier"], "Should indicate current KYC tier is 0")
				assert.NotEmpty(t, responseBody["verification_url"], "Should provide verification URL")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Create Video Stream (No KYC Required)
	// Verifies that video streams can be created without KYC restrictions
	// demonstrating the business rule difference between store and video stream types
	t.Run("CreateVideoStream_NoKYCRequired", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("broadcaster has KYC Tier 0").
			UponReceiving("a request to create a video stream without KYC requirement").
			WithRequest("POST", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.JSONBody(map[string]interface{}{
					"stream_type":          "video",
					"title":                "Gaming Session",
					"description":          "Playing the latest game release",
					"privacy_setting":      "public",
					"scheduled_start_time": "2025-09-30T18:00:00Z",
					"language":             "en",
					"tags":                 []string{"gaming", "entertainment"},
				})
			}).
			WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Stream created successfully"),
					"data": map[string]interface{}{
						"id":                   matchers.Regex("550e8400-e29b-41d4-a716-446655440002", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"broadcaster_id":       matchers.Regex("550e8400-e29b-41d4-a716-446655440003", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"broadcaster_kyc_tier": matchers.Integer(0), // Video streams don't require KYC
						"stream_type":          matchers.Like("video"),
						"title":                matchers.Like("Gaming Session"),
						"description":          matchers.Like("Playing the latest game release"),
						"privacy_setting":      matchers.Like("public"),
						"status":               matchers.Like("scheduled"),
						"stream_key":           matchers.Regex("rtmp://stream.tchat.dev/live/sk_video123xyz", `^rtmp://[a-z0-9\.\-]+/live/sk_[a-zA-Z0-9]+$`),
						"webrtc_session_id":    matchers.Regex("wrtc_session_video456", `^wrtc_session_[a-zA-Z0-9]+$`),
						"scheduled_start_time": matchers.Regex("2025-09-30T18:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"viewer_count":         matchers.Integer(0),
						"peak_viewer_count":    matchers.Integer(0),
						"max_capacity":         matchers.Integer(50000),
						"language":             matchers.Like("en"),
						"tags":                 matchers.EachLike("gaming", 1),
						"created_at":           matchers.Regex("2025-09-30T17:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						"updated_at":           matchers.Regex("2025-09-30T17:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams", config.Host, config.Port)

				createStreamReq := map[string]interface{}{
					"stream_type":          "video",
					"title":                "Gaming Session",
					"description":          "Playing the latest game release",
					"privacy_setting":      "public",
					"scheduled_start_time": "2025-09-30T18:00:00Z",
					"language":             "en",
					"tags":                 []string{"gaming", "entertainment"},
				}

				jsonData, err := json.Marshal(createStreamReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-kyc-tier-0")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify successful creation despite KYC Tier 0
				assert.Equal(t, 201, resp.StatusCode, "Expected 201 Created for video stream")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.True(t, responseBody["success"].(bool), "Response should indicate success")

				data := responseBody["data"].(map[string]interface{})
				assert.Equal(t, "video", data["stream_type"], "Stream type should be 'video'")
				assert.Equal(t, float64(0), data["broadcaster_kyc_tier"], "Video streams should work with KYC Tier 0")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Invalid Request - Missing Required Fields
	// Verifies proper validation error handling for incomplete requests
	t.Run("CreateStream_MissingRequiredFields", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("no preconditions").
			UponReceiving("a request to create stream with missing required fields").
			WithRequest("POST", "/api/v1/streams", func(b *consumer.V2RequestBuilder) {
				b.JSONBody(map[string]interface{}{
					// Missing required fields: stream_type and title
					"description": "Test stream",
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Validation error: missing required fields"),
					"validation_errors": map[string]interface{}{
						"stream_type": matchers.Like("stream_type is required"),
						"title":       matchers.Like("title is required"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams", config.Host, config.Port)

				// Invalid request with missing fields
				invalidReq := map[string]interface{}{
					"description": "Test stream",
				}

				jsonData, err := json.Marshal(invalidReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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

				// Verify validation error
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for validation error")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.NotNil(t, responseBody["validation_errors"], "Validation errors should be detailed")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}