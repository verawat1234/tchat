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

// TestGetChatPactConsumer runs consumer Pact tests for GET /api/v1/streams/{streamId}/chat endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for retrieving chat message history with pagination and moderation filtering capabilities.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for GET /streams/{streamId}/chat endpoint
// 2. No chat message retrieval logic is implemented
// 3. No pagination logic (limit, before_timestamp, after_timestamp) exists
// 4. No moderation filtering (moderation_status) is implemented
// 5. No ScyllaDB query layer exists for chat message history
func TestGetChatPactConsumer(t *testing.T) {
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

	// Test Case 1: Get Recent Chat Messages with Pagination
	// Verifies that clients can retrieve the most recent chat messages
	// with proper pagination metadata (limit, has_more)
	t.Run("GetChat_RecentMessages", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440000"

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("stream has multiple chat messages").
			UponReceiving("a request to retrieve recent chat messages").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/chat", streamID), func(b *consumer.V2RequestBuilder) {
				b.Query("limit", matchers.String("50"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with matchers
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Chat messages retrieved successfully"),
					"data": map[string]interface{}{
						"messages": matchers.EachLike(map[string]interface{}{
							// Core identifiers - use UUID regex for validation
							"message_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_id":  matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"sender_id":  matchers.Regex("550e8400-e29b-41d4-a716-446655440002", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

							// Message content
							"sender_display_name": matchers.Like("JohnDoe"),
							"message_text":        matchers.Like("Great product!"),
							"message_type":        matchers.Like("text"),

							// Moderation metadata
							"moderation_status": matchers.Like("visible"),

							// Timestamps - ISO 8601 format
							"timestamp": matchers.Regex("2025-09-30T15:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						}, 50), // Expect 50 messages (default limit)

						// Pagination metadata
						"has_more": matchers.Like(true),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL with query parameters
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat?limit=50",
					config.Host, config.Port, streamID)

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
				assert.NotNil(t, data["messages"], "Messages array should be present")
				assert.NotNil(t, data["has_more"], "has_more flag should be present")

				// Verify messages array structure
				messages := data["messages"].([]interface{})
				assert.Len(t, messages, 50, "Should return exactly 50 messages (default limit)")

				// Verify first message structure
				firstMessage := messages[0].(map[string]interface{})
				assert.NotEmpty(t, firstMessage["message_id"], "Message ID should be present")
				assert.NotEmpty(t, firstMessage["stream_id"], "Stream ID should be present")
				assert.NotEmpty(t, firstMessage["sender_id"], "Sender ID should be present")
				assert.NotEmpty(t, firstMessage["sender_display_name"], "Sender display name should be present")
				assert.NotEmpty(t, firstMessage["message_text"], "Message text should be present")
				assert.Equal(t, "visible", firstMessage["moderation_status"], "Moderation status should default to visible")
				assert.NotEmpty(t, firstMessage["timestamp"], "Timestamp should be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Filter Chat Messages by Moderation Status
	// Verifies that clients can filter messages by moderation_status
	// to retrieve only visible, flagged, or removed messages
	t.Run("GetChat_FilterByModeration", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440010"

		// Define the expected interaction for filtering by moderation status
		err := mockProvider.
			AddInteraction().
			Given("stream has messages with different moderation statuses").
			UponReceiving("a request to filter messages by moderation status").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/chat", streamID), func(b *consumer.V2RequestBuilder) {
				b.Query("moderation_status", matchers.String("visible"))
				b.Query("limit", matchers.String("50"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Chat messages retrieved successfully"),
					"data": map[string]interface{}{
						"messages": matchers.EachLike(map[string]interface{}{
							"message_id":          matchers.Regex("550e8400-e29b-41d4-a716-446655440011", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_id":           matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"sender_id":           matchers.Regex("550e8400-e29b-41d4-a716-446655440012", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"sender_display_name": matchers.Like("JaneDoe"),
							"message_text":        matchers.Like("This looks amazing!"),
							"message_type":        matchers.Like("text"),
							// All messages should have moderation_status = "visible"
							"moderation_status": matchers.Like("visible"),
							"timestamp":         matchers.Regex("2025-09-30T16:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						}, 30), // Return 30 visible messages

						"has_more": matchers.Like(false),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat?moderation_status=visible&limit=50",
					config.Host, config.Port, streamID)

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
				messages := data["messages"].([]interface{})

				// Verify all returned messages have moderation_status = "visible"
				for _, msg := range messages {
					messageData := msg.(map[string]interface{})
					assert.Equal(t, "visible", messageData["moderation_status"], "All messages should have moderation_status 'visible'")
				}

				assert.Equal(t, 30, len(messages), "Should return 30 visible messages")
				assert.False(t, data["has_more"].(bool), "has_more should be false when no more messages")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Query Messages by Timestamp Range
	// Verifies that clients can query messages between specific timestamps
	// using before_timestamp and after_timestamp parameters
	t.Run("GetChat_TimestampRange", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440020"
		beforeTimestamp := "2025-09-30T17:00:00Z"
		afterTimestamp := "2025-09-30T16:00:00Z"

		// Define the expected interaction for timestamp range queries
		err := mockProvider.
			AddInteraction().
			Given("stream has messages across different time ranges").
			UponReceiving("a request to query messages by timestamp range").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/chat", streamID), func(b *consumer.V2RequestBuilder) {
				b.Query("before_timestamp", matchers.String(beforeTimestamp))
				b.Query("after_timestamp", matchers.String(afterTimestamp))
				b.Query("limit", matchers.String("100"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Chat messages retrieved successfully"),
					"data": map[string]interface{}{
						"messages": matchers.EachLike(map[string]interface{}{
							"message_id":          matchers.Regex("550e8400-e29b-41d4-a716-446655440021", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_id":           matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"sender_id":           matchers.Regex("550e8400-e29b-41d4-a716-446655440022", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"sender_display_name": matchers.Like("BobSmith"),
							"message_text":        matchers.Like("When does the sale end?"),
							"message_type":        matchers.Like("text"),
							"moderation_status":   matchers.Like("visible"),
							// All timestamps should be within the specified range
							"timestamp": matchers.Regex("2025-09-30T16:30:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						}, 45), // Return 45 messages within the time range

						"has_more": matchers.Like(false),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat?before_timestamp=%s&after_timestamp=%s&limit=100",
					config.Host, config.Port, streamID, beforeTimestamp, afterTimestamp)

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
				messages := data["messages"].([]interface{})

				// Verify messages are within the timestamp range
				beforeTime, _ := time.Parse(time.RFC3339, beforeTimestamp)
				afterTime, _ := time.Parse(time.RFC3339, afterTimestamp)

				for _, msg := range messages {
					messageData := msg.(map[string]interface{})
					msgTimestamp, _ := time.Parse(time.RFC3339, messageData["timestamp"].(string))

					assert.True(t, msgTimestamp.Before(beforeTime) || msgTimestamp.Equal(beforeTime),
						"Message timestamp should be before or equal to before_timestamp")
					assert.True(t, msgTimestamp.After(afterTime) || msgTimestamp.Equal(afterTime),
						"Message timestamp should be after or equal to after_timestamp")
				}

				assert.Equal(t, 45, len(messages), "Should return 45 messages within time range")
				assert.False(t, data["has_more"].(bool), "has_more should be false when no more messages")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Empty Chat History
	// Verifies that the API properly handles streams with no chat messages
	// returning an empty array with valid pagination metadata
	t.Run("GetChat_EmptyResult", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440030"

		// Define the expected interaction for empty chat history
		err := mockProvider.
			AddInteraction().
			Given("stream has no chat messages").
			UponReceiving("a request that returns no chat messages").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/chat", streamID), func(b *consumer.V2RequestBuilder) {
				b.Query("limit", matchers.String("50"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Chat messages retrieved successfully"),
					"data": map[string]interface{}{
						"messages": []interface{}{}, // Empty array
						"has_more": matchers.Like(false),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat?limit=50",
					config.Host, config.Port, streamID)

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
				messages := data["messages"].([]interface{})

				// Verify empty results
				assert.Empty(t, messages, "Messages array should be empty")
				assert.False(t, data["has_more"].(bool), "has_more should be false for empty results")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Pagination with Different Message Types
	// Verifies that pagination works correctly with different message types
	// (text, emoji, system messages) and properly handles has_more flag
	t.Run("GetChat_MixedMessageTypes", func(t *testing.T) {
		streamID := "550e8400-e29b-41d4-a716-446655440040"

		// Define the expected interaction for mixed message types
		err := mockProvider.
			AddInteraction().
			Given("stream has mixed message types").
			UponReceiving("a request to retrieve messages with different types").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/chat", streamID), func(b *consumer.V2RequestBuilder) {
				b.Query("limit", matchers.String("20"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Chat messages retrieved successfully"),
					"data": map[string]interface{}{
						"messages": []interface{}{
							// Text message
							map[string]interface{}{
								"message_id":          matchers.Regex("550e8400-e29b-41d4-a716-446655440041", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"stream_id":           matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"sender_id":           matchers.Regex("550e8400-e29b-41d4-a716-446655440042", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"sender_display_name": matchers.Like("AliceWonder"),
								"message_text":        matchers.Like("I love this product!"),
								"message_type":        matchers.Like("text"),
								"moderation_status":   matchers.Like("visible"),
								"timestamp":           matchers.Regex("2025-09-30T18:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							},
							// Emoji message
							map[string]interface{}{
								"message_id":          matchers.Regex("550e8400-e29b-41d4-a716-446655440043", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"stream_id":           matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"sender_id":           matchers.Regex("550e8400-e29b-41d4-a716-446655440044", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"sender_display_name": matchers.Like("CharlieRose"),
								"message_text":        matchers.Like("👍"),
								"message_type":        matchers.Like("emoji"),
								"moderation_status":   matchers.Like("visible"),
								"timestamp":           matchers.Regex("2025-09-30T18:01:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							},
							// System message
							map[string]interface{}{
								"message_id":          matchers.Regex("550e8400-e29b-41d4-a716-446655440045", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"stream_id":           matchers.Regex(streamID, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"sender_id":           matchers.Regex("00000000-0000-0000-0000-000000000000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
								"sender_display_name": matchers.Like("System"),
								"message_text":        matchers.Like("Stream has started"),
								"message_type":        matchers.Like("system"),
								"moderation_status":   matchers.Like("visible"),
								"timestamp":           matchers.Regex("2025-09-30T17:59:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							},
						},

						"has_more": matchers.Like(true),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat?limit=20",
					config.Host, config.Port, streamID)

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
				messages := data["messages"].([]interface{})

				// Verify different message types are present
				assert.Len(t, messages, 3, "Should return 3 messages of different types")

				// Verify message type variety
				messageTypes := make(map[string]bool)
				for _, msg := range messages {
					messageData := msg.(map[string]interface{})
					msgType := messageData["message_type"].(string)
					messageTypes[msgType] = true
				}

				assert.True(t, messageTypes["text"], "Should include text message")
				assert.True(t, messageTypes["emoji"], "Should include emoji message")
				assert.True(t, messageTypes["system"], "Should include system message")

				assert.True(t, data["has_more"].(bool), "has_more should be true when more messages exist")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 6: Invalid Stream ID (404 Error Handling)
	// Verifies that the API properly handles requests for non-existent streams
	// returning appropriate 404 error responses
	t.Run("GetChat_InvalidStreamID", func(t *testing.T) {
		invalidStreamID := "00000000-0000-0000-0000-000000000000"

		// Define the expected interaction for invalid stream ID
		err := mockProvider.
			AddInteraction().
			Given("stream does not exist").
			UponReceiving("a request with invalid stream ID").
			WithRequest("GET", fmt.Sprintf("/api/v1/streams/%s/chat", invalidStreamID), func(b *consumer.V2RequestBuilder) {
				b.Query("limit", matchers.String("50"))
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))

				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Stream not found"),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/chat?limit=50",
					config.Host, config.Port, invalidStreamID)

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

				assert.Equal(t, 404, resp.StatusCode, "Expected 404 Not Found status code")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}