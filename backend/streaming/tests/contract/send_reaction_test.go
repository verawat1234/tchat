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

// TestSendReactionPactConsumer runs consumer Pact tests for POST /api/v1/streams/{streamId}/react endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for sending emoji reactions during video live streams.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for POST /streams/{streamId}/react endpoint
// 2. No ScyllaDB stream_reactions table integration exists
// 3. No rate limiting (10 reactions/second) is implemented
// 4. No emoji unicode validation is implemented
// 5. No video-stream-only restriction enforcement (store streams cannot receive reactions)
// 6. No 30-day TTL enforcement is implemented
// 7. No anonymous user support (nullable viewer_id) is implemented
// 8. No real-time aggregation via Redis is implemented
func TestSendReactionPactConsumer(t *testing.T) {
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

	// Test Case 1: Successful Emoji Reaction Send (Video Stream)
	// Verifies that authenticated users can send emoji reactions to video streams
	// with proper ScyllaDB storage and Redis aggregation
	t.Run("SendReaction_Success", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440001 is a live video stream").
			UponReceiving("a request to send heart emoji reaction from authenticated user").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440001/react", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/react$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					// Emoji unicode - supports heart, thumbs up, laughing, party emojis
					"reaction": matchers.Regex("❤️", `^[\x{1F300}-\x{1F9FF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}]+$`),
				})
			}).
			WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with complete StreamReaction entity
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Reaction sent successfully"),
					"data": map[string]interface{}{
						"stream_reaction": map[string]interface{}{
							// Reaction identity (ScyllaDB UUID)
							"reaction_id": matchers.Regex("d1e2f3a4-b5c6-4d5e-9f0a-1b2c3d4e5f6a", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Stream reference
							"stream_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Viewer information (authenticated user)
							"viewer_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440003", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Emoji reaction type
							"reaction_type": matchers.Like("❤️"),

							// Timestamp (ScyllaDB clustering key for sorting)
							"timestamp": matchers.Regex("2025-09-30T15:35:20Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),

							// 30-day TTL indication (not exposed in response but validated in backend)
							"ttl_seconds": matchers.Integer(2592000), // 30 days in seconds
						},

						// Real-time aggregation data from Redis
						"aggregation": map[string]interface{}{
							"total_reactions":   matchers.Integer(1523),
							"reaction_counts": map[string]interface{}{
								"❤️": matchers.Integer(850),
								"👍": matchers.Integer(420),
								"😂": matchers.Integer(180),
								"🎉": matchers.Integer(73),
							},
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL with stream ID
				streamID := "550e8400-e29b-41d4-a716-446655440001"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/react", config.Host, config.Port, streamID)

				// Prepare request payload with emoji reaction
				sendReactionReq := map[string]interface{}{
					"reaction": "❤️",
				}

				// Marshal request body to JSON
				jsonData, err := json.Marshal(sendReactionReq)
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
				req.Header.Set("Authorization", "Bearer test-token-viewer")

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

				streamReaction := data["stream_reaction"].(map[string]interface{})

				// Verify required fields are present
				assert.NotEmpty(t, streamReaction["reaction_id"], "Reaction ID should be present")
				assert.NotEmpty(t, streamReaction["stream_id"], "Stream ID should be present")
				assert.NotEmpty(t, streamReaction["viewer_id"], "Viewer ID should be present")
				assert.NotEmpty(t, streamReaction["reaction_type"], "Reaction type should be present")
				assert.NotEmpty(t, streamReaction["timestamp"], "Timestamp should be present")

				// Verify reaction type is preserved
				assert.Equal(t, "❤️", streamReaction["reaction_type"], "Reaction type should be heart emoji")

				// Verify 30-day TTL is set (2592000 seconds)
				ttl, ok := streamReaction["ttl_seconds"].(float64)
				assert.True(t, ok, "TTL should be present and numeric")
				assert.Equal(t, float64(2592000), ttl, "TTL should be 30 days (2592000 seconds)")

				// Verify aggregation data is included
				aggregation, ok := data["aggregation"].(map[string]interface{})
				assert.True(t, ok, "Aggregation data should be present")
				assert.NotEmpty(t, aggregation["total_reactions"], "Total reactions count should be present")
				assert.NotEmpty(t, aggregation["reaction_counts"], "Reaction counts breakdown should be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Store Stream Cannot Receive Reactions
	// Verifies that reactions can only be sent to video streams, not store streams (400 Bad Request)
	t.Run("SendReaction_VideoStreamOnly", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440020 is a live store stream").
			UponReceiving("a request to send reaction to store stream").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440020/react", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/react$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					"reaction": matchers.Like("👍"),
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Reactions are only available for video streams"),
					"details": map[string]interface{}{
						"stream_id":   matchers.Regex("550e8400-e29b-41d4-a716-446655440020", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"stream_type": matchers.Like("store"),
						"allowed_types": []interface{}{
							matchers.Like("video"),
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440020"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/react", config.Host, config.Port, streamID)

				sendReactionReq := map[string]interface{}{
					"reaction": "👍",
				}

				jsonData, err := json.Marshal(sendReactionReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-viewer")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 400 Bad Request
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for store stream")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "video streams", "Error should mention video streams only")

				// Verify details include stream type
				details := responseBody["details"].(map[string]interface{})
				assert.Equal(t, "store", details["stream_type"], "Should indicate stream type is store")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Invalid Emoji Unicode Validation
	// Verifies that invalid emoji unicode strings are rejected (400 Bad Request)
	t.Run("SendReaction_InvalidEmoji", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440001 is a live video stream").
			UponReceiving("a request to send invalid emoji unicode").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440001/react", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/react$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					// Invalid emoji - regular text instead of emoji unicode
					"reaction": "invalid",
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Invalid emoji unicode"),
					"details": map[string]interface{}{
						"reaction": matchers.Like("invalid"),
						"message":  matchers.Like("Reaction must be a valid emoji unicode character"),
						"allowed_emojis": []interface{}{
							matchers.Like("❤️"),
							matchers.Like("👍"),
							matchers.Like("😂"),
							matchers.Like("🎉"),
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440001"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/react", config.Host, config.Port, streamID)

				sendReactionReq := map[string]interface{}{
					"reaction": "invalid",
				}

				jsonData, err := json.Marshal(sendReactionReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-viewer")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 400 Bad Request
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for invalid emoji")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "Invalid emoji", "Error should mention invalid emoji")

				// Verify validation details
				details := responseBody["details"].(map[string]interface{})
				assert.NotEmpty(t, details["message"], "Should provide validation details")
				assert.NotEmpty(t, details["allowed_emojis"], "Should provide list of allowed emojis")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Rate Limit Exceeded
	// Verifies that sending more than 10 reactions/second is rejected (429 Too Many Requests)
	t.Run("SendReaction_RateLimit", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("user has sent 10 reactions in the last second to stream 550e8400-e29b-41d4-a716-446655440001").
			UponReceiving("a request to send 11th reaction within 1 second").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440001/react", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/react$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					"reaction": matchers.Like("🎉"),
				})
			}).
			WillRespondWith(429, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Retry-After", matchers.String("1")) // Retry after 1 second
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Rate limit exceeded"),
					"details": map[string]interface{}{
						"limit":         matchers.Integer(10),
						"window":        matchers.Like("1 second"),
						"retry_after":   matchers.Integer(1),
						"current_count": matchers.Integer(11),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440001"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/react", config.Host, config.Port, streamID)

				sendReactionReq := map[string]interface{}{
					"reaction": "🎉",
				}

				jsonData, err := json.Marshal(sendReactionReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-viewer")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 429 Too Many Requests
				assert.Equal(t, 429, resp.StatusCode, "Expected 429 Too Many Requests for rate limit")

				// Verify Retry-After header
				retryAfter := resp.Header.Get("Retry-After")
				assert.Equal(t, "1", retryAfter, "Should include Retry-After header")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "Rate limit", "Error should mention rate limit")

				// Verify rate limit details
				details := responseBody["details"].(map[string]interface{})
				limit, ok := details["limit"].(float64)
				assert.True(t, ok, "Should provide rate limit")
				assert.Equal(t, float64(10), limit, "Rate limit should be 10 reactions per second")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Anonymous User Sends Reaction
	// Verifies that anonymous viewers can send reactions with null viewer_id
	// (nullable viewer_id field in ScyllaDB)
	t.Run("SendReaction_AnonymousUser", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440001 is a live video stream").
			UponReceiving("a request to send reaction from anonymous user").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440001/react", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/react$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				// No Authorization header for anonymous users
				b.JSONBody(map[string]interface{}{
					"reaction": matchers.Like("😂"),
				})
			}).
			WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Reaction sent successfully"),
					"data": map[string]interface{}{
						"stream_reaction": map[string]interface{}{
							"reaction_id": matchers.Regex("e2f3a4b5-c6d7-4e5f-0a1b-2c3d4e5f6a7b", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_id":   matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Anonymous user: viewer_id is null
							"viewer_id": nil,

							"reaction_type": matchers.Like("😂"),
							"timestamp":     matchers.Regex("2025-09-30T15:36:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"ttl_seconds":   matchers.Integer(2592000),
						},
						"aggregation": map[string]interface{}{
							"total_reactions": matchers.Integer(1524),
							"reaction_counts": map[string]interface{}{
								"❤️": matchers.Integer(850),
								"👍": matchers.Integer(420),
								"😂": matchers.Integer(181),
								"🎉": matchers.Integer(73),
							},
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440001"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/react", config.Host, config.Port, streamID)

				sendReactionReq := map[string]interface{}{
					"reaction": "😂",
				}

				jsonData, err := json.Marshal(sendReactionReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				// No Authorization header for anonymous users

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 201 Created
				assert.Equal(t, 201, resp.StatusCode, "Expected 201 Created for anonymous user")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify anonymous user handling
				data := responseBody["data"].(map[string]interface{})
				streamReaction := data["stream_reaction"].(map[string]interface{})

				// viewer_id should be null for anonymous users
				assert.Nil(t, streamReaction["viewer_id"], "Anonymous users should have null viewer_id")
				assert.Equal(t, "😂", streamReaction["reaction_type"], "Reaction type should be laughing emoji")

				// Verify aggregation data is still included
				aggregation := data["aggregation"].(map[string]interface{})
				assert.NotEmpty(t, aggregation["total_reactions"], "Total reactions count should be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 6: Empty Reaction Validation Error
	// Verifies that empty reaction strings are rejected (400 Bad Request)
	t.Run("SendReaction_EmptyReaction", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440001 is a live video stream").
			UponReceiving("a request to send empty reaction").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440001/react", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/react$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					// Empty reaction
					"reaction": "",
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Invalid reaction content"),
					"details": map[string]interface{}{
						"reaction": matchers.Like("Reaction cannot be empty and must be a valid emoji unicode"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440001"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/react", config.Host, config.Port, streamID)

				sendReactionReq := map[string]interface{}{
					"reaction": "",
				}

				jsonData, err := json.Marshal(sendReactionReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-viewer")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 400 Bad Request
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for empty reaction")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "Invalid", "Error should mention invalid reaction")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}