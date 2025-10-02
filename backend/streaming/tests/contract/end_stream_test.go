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

// TestEndStreamPactConsumer runs consumer Pact tests for POST /api/v1/streams/{streamId}/end endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for ending live streams and initiating recording processing.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for POST /streams/{streamId}/end endpoint
// 2. No stream status validation (live -> ended transition) is implemented
// 3. No recording processing initiation logic exists
// 4. No 30-day recording expiry calculation is implemented
// 5. No broadcaster ownership verification exists
func TestEndStreamPactConsumer(t *testing.T) {
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

	// Test Case 1: Successful Stream End with Recording Processing
	// Verifies that broadcasters can terminate live streams, triggering recording
	// processing and 30-day retention period calculation
	t.Run("EndStream_Success", func(t *testing.T) {
		// Calculate expected recording expiry (30 days from now)
		// This matches the OpenAPI specification requirement for 30-day retention
		expectedExpiryTime := time.Now().Add(30 * 24 * time.Hour)
		expiryDatePattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`

		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 exists and is live").
			UponReceiving("a request to end live stream").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/end", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/end$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-broadcaster", `^Bearer .+$`))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with ended stream data
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Stream ended successfully"),
					"data": map[string]interface{}{
						// Updated stream object with ended status
						"stream": map[string]interface{}{
							"id":                  matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"status":              matchers.Like("ended"),
							"broadcaster_id":      matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_type":         matchers.Like("store"),
							"title":               matchers.Like("Product Launch Event"),
							"privacy_setting":     matchers.Like("public"),

							// Timing fields
							"scheduled_start_time": matchers.Regex("2025-09-30T15:00:00Z", expiryDatePattern),
							"actual_start_time":   matchers.Regex("2025-09-30T15:05:00Z", expiryDatePattern),
							"end_time":            matchers.Regex("2025-09-30T16:30:00Z", expiryDatePattern),

							// Recording information
							"recording_url":        matchers.Regex("https://cdn.tchat.dev/recordings/550e8400-e29b-41d4-a716-446655440000.mp4", `^https://[a-z0-9\.\-]+/recordings/[a-zA-Z0-9\-]+\.mp4$`),
							"recording_expiry_date": matchers.Regex(expectedExpiryTime.Format(time.RFC3339), expiryDatePattern),

							// Stream statistics
							"viewer_count":        matchers.Integer(0), // Should be 0 after stream ends
							"peak_viewer_count":   matchers.Integer(1234),
						},

						// Recording processing metadata
						"recording_metadata": map[string]interface{}{
							"processing_status": matchers.Like("processing"),
							"duration_seconds":  matchers.Integer(5100), // 85 minutes
							"file_size_bytes":   matchers.Integer(0), // Not yet available during processing
							"retention_days":    matchers.Integer(30),
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL with stream ID
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/end", config.Host, config.Port, streamID)

				// Create HTTP POST request (no request body needed)
				req, err := http.NewRequest("POST", url, nil)
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

				// Verify stream status transition
				stream := data["stream"].(map[string]interface{})
				assert.Equal(t, "ended", stream["status"], "Stream status should be 'ended'")
				assert.NotEmpty(t, stream["end_time"], "End time should be set")
				assert.NotEmpty(t, stream["recording_url"], "Recording URL should be set")
				assert.NotEmpty(t, stream["recording_expiry_date"], "Recording expiry date should be set")
				assert.Equal(t, float64(0), stream["viewer_count"], "Viewer count should be 0 after stream ends")
				assert.Greater(t, stream["peak_viewer_count"].(float64), float64(0), "Peak viewer count should be greater than 0")

				// Verify recording metadata
				recordingMeta := data["recording_metadata"].(map[string]interface{})
				assert.Equal(t, "processing", recordingMeta["processing_status"], "Recording should be processing")
				assert.Equal(t, float64(30), recordingMeta["retention_days"], "Should have 30-day retention")
				assert.Greater(t, recordingMeta["duration_seconds"].(float64), float64(0), "Duration should be positive")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Stream Not Live - Cannot End
	// Verifies that attempting to end a stream that isn't live returns 409 Conflict
	t.Run("EndStream_NotLive", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440020 exists and is scheduled but not live").
			UponReceiving("a request to end a stream that is not live").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440020/end", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/end$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-broadcaster", `^Bearer .+$`))
			}).
			WillRespondWith(409, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success":        false,
					"error":          matchers.Like("Cannot end stream that is not live"),
					"current_status": matchers.Like("scheduled"),
					"stream_id":      matchers.Regex("550e8400-e29b-41d4-a716-446655440020", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"details": map[string]interface{}{
						"message": matchers.Like("Stream must be in 'live' status to be ended"),
						"allowed_statuses": matchers.EachLike("live", 1),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440020"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/end", config.Host, config.Port, streamID)

				req, err := http.NewRequest("POST", url, nil)
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

				// Verify 409 Conflict
				assert.Equal(t, 409, resp.StatusCode, "Expected 409 Conflict for non-live stream")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Equal(t, "scheduled", responseBody["current_status"], "Should indicate stream is scheduled")
				assert.Contains(t, responseBody["error"].(string), "not live", "Error should mention stream is not live")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Non-Broadcaster Attempting to End Stream
	// Verifies that only the broadcaster can end their own stream (403 Forbidden)
	t.Run("EndStream_NotOwner", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 exists and is live but user is not broadcaster").
			UponReceiving("a request to end stream from non-broadcaster").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/end", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/end$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-other-user", `^Bearer .+$`))
			}).
			WillRespondWith(403, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Only the broadcaster can end their stream"),
					"stream_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"broadcaster_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"details": map[string]interface{}{
						"message": matchers.Like("This stream belongs to another user"),
						"permission": matchers.Like("broadcaster_only"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/end", config.Host, config.Port, streamID)

				req, err := http.NewRequest("POST", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-other-user")

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

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "broadcaster", "Error should mention broadcaster requirement")
				assert.NotEmpty(t, responseBody["stream_id"], "Should include stream ID")
				assert.NotEmpty(t, responseBody["broadcaster_id"], "Should include actual broadcaster ID")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Recording Retention Calculation Verification
	// Verifies that the recording_expiry_date is exactly 30 days from the end_time
	// This ensures compliance with the OpenAPI specification for 30-day retention
	t.Run("EndStream_RecordingRetention", func(t *testing.T) {
		// Use a fixed end time for predictable expiry calculation
		endTime := "2025-09-30T16:30:00Z"
		expectedExpiryDate := "2025-10-30T16:30:00Z" // Exactly 30 days later

		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440030 exists and is live").
			UponReceiving("a request to end stream with recording retention validation").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440030/end", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/end$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-broadcaster", `^Bearer .+$`))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Stream ended successfully"),
					"data": map[string]interface{}{
						"stream": map[string]interface{}{
							"id":                  matchers.Regex("550e8400-e29b-41d4-a716-446655440030", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"status":              matchers.Like("ended"),
							"end_time":            matchers.Regex(endTime, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"recording_url":        matchers.Regex("https://cdn.tchat.dev/recordings/550e8400-e29b-41d4-a716-446655440030.mp4", `^https://[a-z0-9\.\-]+/recordings/[a-zA-Z0-9\-]+\.mp4$`),
							"recording_expiry_date": matchers.Regex(expectedExpiryDate, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"broadcaster_id":      matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_type":         matchers.Like("video"),
							"title":               matchers.Like("Live Gaming Session"),
						},
						"recording_metadata": map[string]interface{}{
							"processing_status": matchers.Like("processing"),
							"retention_days":    matchers.Integer(30),
							"expiry_reminder_sent": matchers.Boolean(false),
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440030"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/end", config.Host, config.Port, streamID)

				req, err := http.NewRequest("POST", url, nil)
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

				// Verify 200 OK
				assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status code")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify recording retention
				data := responseBody["data"].(map[string]interface{})
				stream := data["stream"].(map[string]interface{})

				// Parse end_time and recording_expiry_date
				endTimeStr, ok := stream["end_time"].(string)
				assert.True(t, ok, "end_time should be a string")

				expiryDateStr, ok := stream["recording_expiry_date"].(string)
				assert.True(t, ok, "recording_expiry_date should be a string")

				// Verify expiry is exactly 30 days after end time
				parsedEndTime, err := time.Parse(time.RFC3339, endTimeStr)
				if err != nil {
					return fmt.Errorf("failed to parse end_time: %w", err)
				}

				parsedExpiryDate, err := time.Parse(time.RFC3339, expiryDateStr)
				if err != nil {
					return fmt.Errorf("failed to parse recording_expiry_date: %w", err)
				}

				// Calculate expected expiry (30 days after end time)
				expectedExpiry := parsedEndTime.Add(30 * 24 * time.Hour)

				// Allow small tolerance for time comparison (1 minute)
				timeDiff := parsedExpiryDate.Sub(expectedExpiry).Abs()
				assert.LessOrEqual(t, timeDiff, time.Minute, "Recording expiry should be exactly 30 days after end time")

				// Verify recording metadata
				recordingMeta := data["recording_metadata"].(map[string]interface{})
				assert.Equal(t, float64(30), recordingMeta["retention_days"], "Should have 30-day retention")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Stream Not Found
	// Verifies that attempting to end a non-existent stream returns 404 Not Found
	t.Run("EndStream_NotFound", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-999999999999 does not exist").
			UponReceiving("a request to end non-existent stream").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-999999999999/end", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/end$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-broadcaster", `^Bearer .+$`))
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Stream not found"),
					"stream_id": matchers.Regex("550e8400-e29b-41d4-a716-999999999999", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"details": map[string]interface{}{
						"message": matchers.Like("The requested stream does not exist or has been deleted"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-999999999999"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/end", config.Host, config.Port, streamID)

				req, err := http.NewRequest("POST", url, nil)
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

				// Verify 404 Not Found
				assert.Equal(t, 404, resp.StatusCode, "Expected 404 Not Found for non-existent stream")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Equal(t, "Stream not found", responseBody["error"], "Error should indicate stream not found")
				assert.NotEmpty(t, responseBody["stream_id"], "Should include stream ID")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}