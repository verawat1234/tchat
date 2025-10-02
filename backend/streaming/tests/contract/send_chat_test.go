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

// TestSendChatPactConsumer runs consumer Pact tests for POST /api/v1/streams/{streamId}/chat endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for sending real-time chat messages during live streams.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for POST /streams/{streamId}/chat endpoint
// 2. No ScyllaDB chat_messages table integration exists
// 3. No rate limiting (5 messages/second) is implemented
// 4. No message validation (max 500 characters) is implemented
// 5. No moderation status management exists
// 6. No anonymous user support (nullable user_id) is implemented
// 7. No 30-day TTL enforcement is implemented
func TestSendChatPactConsumer(t *testing.T) {
	// Create Pact consumer with mock provider configuration
	// Consumer: streaming-web-client (represents web frontend)
	// Provider: streaming-service (represents backend streaming service)
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "streaming-web-client",
		Provider: "streaming-service",
		Host:     "localhost",
		Port:     9001, // Same port as start_stream_test.go
		PactDir:  "./pacts", // Output directory for generated pact files
	})
	assert.NoError(t, err, "Failed to create Pact mock provider")

	// Test Case 1: Successful Chat Message Send
	// Verifies that authenticated users can send chat messages to live streams
	// with proper ScyllaDB storage and moderation status
	t.Run("SendChat_Success", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 is live and accepting chat").
			UponReceiving("a request to send chat message from authenticated user").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/chat", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/chat$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					// Chat message content with maximum 500 characters validation
					"message": matchers.Regex(
						"Great product! Love this live stream!",
						`^.{1,500}$`,
					),
					// Optional message type (text/emoji)
					"message_type": matchers.Like("text"),
				})
			}).
			WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with complete ChatMessage entity
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Message sent successfully"),
					"data": map[string]interface{}{
						"chat_message": map[string]interface{}{
							// Message identity (ScyllaDB UUID)
							"message_id": matchers.Regex("a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Stream reference
							"stream_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Sender information (authenticated user)
							"sender_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440002", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"sender_display_name": matchers.Like("JohnDoe"),

							// Message content
							"message_text": matchers.Like("Great product! Love this live stream!"),
							"message_type":  matchers.Like("text"),

							// Moderation status (visible/removed/flagged)
							"moderation_status": matchers.Like("visible"),

							// Timestamp (ScyllaDB clustering key for sorting)
							"timestamp": matchers.Regex("2025-09-30T15:30:45Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),

							// 30-day TTL indication (not exposed in response but validated in backend)
							"ttl_seconds": matchers.Integer(2592000), // 30 days in seconds
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL with stream ID
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat", config.Host, config.Port, streamID)

				// Prepare request payload with chat message
				sendChatReq := map[string]interface{}{
					"message":      "Great product! Love this live stream!",
					"message_type": "text",
				}

				// Marshal request body to JSON
				jsonData, err := json.Marshal(sendChatReq)
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

				chatMessage := data["chat_message"].(map[string]interface{})

				// Verify required fields are present
				assert.NotEmpty(t, chatMessage["message_id"], "Message ID should be present")
				assert.NotEmpty(t, chatMessage["stream_id"], "Stream ID should be present")
				assert.NotEmpty(t, chatMessage["sender_id"], "Sender ID should be present")
				assert.NotEmpty(t, chatMessage["sender_display_name"], "Sender display name should be present")
				assert.NotEmpty(t, chatMessage["message_text"], "Message text should be present")
				assert.NotEmpty(t, chatMessage["timestamp"], "Timestamp should be present")

				// Verify moderation status is set correctly
				assert.Equal(t, "visible", chatMessage["moderation_status"], "New messages should have 'visible' moderation status")

				// Verify message type is preserved
				assert.Equal(t, "text", chatMessage["message_type"], "Message type should be 'text'")

				// Verify 30-day TTL is set (2592000 seconds)
				ttl, ok := chatMessage["ttl_seconds"].(float64)
				assert.True(t, ok, "TTL should be present and numeric")
				assert.Equal(t, float64(2592000), ttl, "TTL should be 30 days (2592000 seconds)")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Anonymous User Sends Chat Message
	// Verifies that anonymous viewers can send messages with null user_id
	// (nullable sender_id field in ScyllaDB)
	t.Run("SendChat_AnonymousUser", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 is live and allows anonymous chat").
			UponReceiving("a request to send chat message from anonymous user").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/chat", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/chat$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				// No Authorization header for anonymous users
				b.JSONBody(map[string]interface{}{
					"message":      matchers.Like("This is amazing!"),
					"message_type": matchers.Like("text"),
				})
			}).
			WillRespondWith(201, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Message sent successfully"),
					"data": map[string]interface{}{
						"chat_message": map[string]interface{}{
							"message_id": matchers.Regex("b2c3d4e5-f6a7-4b5c-9d0e-1f2a3b4c5d6e", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_id":  matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Anonymous user: sender_id is null
							"sender_id":           nil,
							"sender_display_name": matchers.Like("Anonymous"),

							"message_text":      matchers.Like("This is amazing!"),
							"message_type":      matchers.Like("text"),
							"moderation_status": matchers.Like("visible"),
							"timestamp":         matchers.Regex("2025-09-30T15:31:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"ttl_seconds":       matchers.Integer(2592000),
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat", config.Host, config.Port, streamID)

				sendChatReq := map[string]interface{}{
					"message":      "This is amazing!",
					"message_type": "text",
				}

				jsonData, err := json.Marshal(sendChatReq)
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
				chatMessage := data["chat_message"].(map[string]interface{})

				// sender_id should be null for anonymous users
				assert.Nil(t, chatMessage["sender_id"], "Anonymous users should have null sender_id")
				assert.Equal(t, "Anonymous", chatMessage["sender_display_name"], "Display name should be 'Anonymous'")
				assert.Equal(t, "visible", chatMessage["moderation_status"], "Anonymous messages should be visible by default")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Stream Not Live - Cannot Send Messages
	// Verifies that messages cannot be sent to ended streams (400 Bad Request)
	t.Run("SendChat_StreamNotLive", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440010 has ended").
			UponReceiving("a request to send chat message to ended stream").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440010/chat", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/chat$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					"message":      matchers.Like("Hello!"),
					"message_type": matchers.Like("text"),
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Stream is not live"),
					"details": map[string]interface{}{
						"stream_id":      matchers.Regex("550e8400-e29b-41d4-a716-446655440010", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"current_status": matchers.Like("ended"),
						"end_time":       matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440010"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat", config.Host, config.Port, streamID)

				sendChatReq := map[string]interface{}{
					"message":      "Hello!",
					"message_type": "text",
				}

				jsonData, err := json.Marshal(sendChatReq)
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
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for ended stream")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "not live", "Error should mention stream is not live")

				// Verify details include current status
				details := responseBody["details"].(map[string]interface{})
				assert.Equal(t, "ended", details["current_status"], "Should indicate stream has ended")
				assert.NotEmpty(t, details["end_time"], "Should show when stream ended")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Empty Message Validation Error
	// Verifies that empty messages are rejected (400 Bad Request)
	t.Run("SendChat_EmptyMessage", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 is live").
			UponReceiving("a request to send empty chat message").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/chat", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/chat$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					// Empty message content
					"message":      "",
					"message_type": matchers.Like("text"),
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Invalid message content"),
					"details": map[string]interface{}{
						"message": matchers.Like("Message cannot be empty and must be between 1 and 500 characters"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat", config.Host, config.Port, streamID)

				sendChatReq := map[string]interface{}{
					"message":      "",
					"message_type": "text",
				}

				jsonData, err := json.Marshal(sendChatReq)
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
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for empty message")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "Invalid", "Error should mention invalid message")

				// Verify validation details
				details := responseBody["details"].(map[string]interface{})
				assert.NotEmpty(t, details["message"], "Should provide validation details")
				assert.Contains(t, details["message"].(string), "empty", "Should mention message cannot be empty")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Message Too Long Validation Error
	// Verifies that messages exceeding 500 characters are rejected (400 Bad Request)
	t.Run("SendChat_MessageTooLong", func(t *testing.T) {
		// Generate a message with 501 characters (exceeds 500 character limit)
		longMessage := string(make([]byte, 501))
		for i := range longMessage {
			longMessage = longMessage[:i] + "a" + longMessage[i+1:]
		}

		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 is live").
			UponReceiving("a request to send chat message exceeding 500 characters").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/chat", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/chat$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					// Message with 501 characters
					"message":      longMessage,
					"message_type": matchers.Like("text"),
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Invalid message content"),
					"details": map[string]interface{}{
						"message":        matchers.Like("Message cannot be empty and must be between 1 and 500 characters"),
						"message_length": matchers.Integer(501),
						"max_length":     matchers.Integer(500),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat", config.Host, config.Port, streamID)

				sendChatReq := map[string]interface{}{
					"message":      longMessage,
					"message_type": "text",
				}

				jsonData, err := json.Marshal(sendChatReq)
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
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for message too long")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")

				// Verify validation details include length information
				details := responseBody["details"].(map[string]interface{})
				messageLength, ok := details["message_length"].(float64)
				assert.True(t, ok, "Should provide message length")
				assert.Equal(t, float64(501), messageLength, "Should indicate actual message length")

				maxLength, ok := details["max_length"].(float64)
				assert.True(t, ok, "Should provide max length")
				assert.Equal(t, float64(500), maxLength, "Should indicate max allowed length")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 6: Rate Limit Exceeded
	// Verifies that sending more than 5 messages/second is rejected (429 Too Many Requests)
	t.Run("SendChat_RateLimitExceeded", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("user has sent 5 messages in the last second to stream 550e8400-e29b-41d4-a716-446655440000").
			UponReceiving("a request to send 6th message within 1 second").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/chat", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/chat$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-viewer", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					"message":      matchers.Like("Too fast!"),
					"message_type": matchers.Like("text"),
				})
			}).
			WillRespondWith(429, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Retry-After", matchers.String("1")) // Retry after 1 second
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Rate limit exceeded"),
					"details": map[string]interface{}{
						"limit":         matchers.Integer(5),
						"window":        matchers.Like("1 second"),
						"retry_after":   matchers.Integer(1),
						"current_count": matchers.Integer(6),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat", config.Host, config.Port, streamID)

				sendChatReq := map[string]interface{}{
					"message":      "Too fast!",
					"message_type": "text",
				}

				jsonData, err := json.Marshal(sendChatReq)
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
				assert.Equal(t, float64(5), limit, "Rate limit should be 5 messages per second")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}