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

// TestStartStreamPactConsumer runs consumer Pact tests for POST /api/v1/streams/{streamId}/start endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for starting scheduled streams with WebRTC negotiation.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for POST /streams/{streamId}/start endpoint
// 2. No WebRTC SDP negotiation logic is implemented
// 3. No ICE server configuration management exists
// 4. No stream status validation (scheduled -> live transition) is implemented
// 5. No broadcaster ownership verification exists
func TestStartStreamPactConsumer(t *testing.T) {
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

	// Test Case 1: Successful Stream Start with WebRTC Negotiation
	// Verifies that broadcasters can start scheduled streams with valid WebRTC SDP offer
	// and receive proper SDP answer with ICE server configuration
	t.Run("StartStream_Success", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 exists and is scheduled").
			UponReceiving("a request to start stream with valid WebRTC offer").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/start", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/start$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-broadcaster", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					// WebRTC SDP offer from broadcaster
					"sdp_offer": matchers.Regex(
						"v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0 1\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:abc123\r\na=ice-pwd:def456xyz789\r\na=fingerprint:sha-256 AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99\r\na=setup:actpass\r\na=mid:0\r\na=sendonly\r\na=rtpmap:96 VP8/90000\r\n",
						`^v=0\r\n.*`,
					),
					// Quality layers for simulcast (adaptive bitrate)
					"quality_layers": matchers.EachLike("720p", 1),
				})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with WebRTC answer and ICE servers
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Stream started successfully"),
					"data": map[string]interface{}{
						// WebRTC SDP answer from server
						"sdp_answer": matchers.Regex(
							"v=0\r\no=- 987654321 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0 1\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:xyz789\r\na=ice-pwd:abc123def456\r\na=fingerprint:sha-256 99:88:77:66:55:44:33:22:11:00:FF:EE:DD:CC:BB:AA:99:88:77:66:55:44:33:22:11:00:FF:EE:DD:CC:BB:AA\r\na=setup:active\r\na=mid:0\r\na=recvonly\r\na=rtpmap:96 VP8/90000\r\n",
							`^v=0\r\n.*`,
						),

						// ICE servers for WebRTC connection
						"ice_servers": matchers.EachLike(map[string]interface{}{
							"urls": matchers.EachLike("stun:stun.tchat.dev:3478", 1),
						}, 1),

						// WebRTC session identifier for tracking
						"webrtc_session_id": matchers.Regex("wrtc_session_abc123xyz789", `^wrtc_session_[a-zA-Z0-9]+$`),

						// Updated stream object with live status
						"stream": map[string]interface{}{
							"id":                  matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"status":              matchers.Like("live"),
							"actual_start_time":   matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
							"viewer_count":        matchers.Integer(0),
							"peak_viewer_count":   matchers.Integer(0),
							"stream_key":          matchers.Regex("rtmp://stream.tchat.dev/live/sk_abc123xyz789", `^rtmp://[a-z0-9\.\-]+/live/sk_[a-zA-Z0-9]+$`),
							"webrtc_session_id":   matchers.Regex("wrtc_session_abc123xyz789", `^wrtc_session_[a-zA-Z0-9]+$`),
							"broadcaster_id":      matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
							"stream_type":         matchers.Like("store"),
							"title":               matchers.Like("Product Launch Event"),
							"privacy_setting":     matchers.Like("public"),
							"scheduled_start_time": matchers.Regex("2025-09-30T15:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
						},
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL with stream ID
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/start", config.Host, config.Port, streamID)

				// Prepare request payload with WebRTC offer
				startStreamReq := map[string]interface{}{
					"sdp_offer":       "v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0 1\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:abc123\r\na=ice-pwd:def456xyz789\r\na=fingerprint:sha-256 AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99\r\na=setup:actpass\r\na=mid:0\r\na=sendonly\r\na=rtpmap:96 VP8/90000\r\n",
					"quality_layers": []string{"360p", "720p", "1080p"},
				}

				// Marshal request body to JSON
				jsonData, err := json.Marshal(startStreamReq)
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

				// Verify WebRTC negotiation fields are present
				assert.NotEmpty(t, data["sdp_answer"], "SDP answer should be present")
				assert.NotEmpty(t, data["ice_servers"], "ICE servers should be present")
				assert.NotEmpty(t, data["webrtc_session_id"], "WebRTC session ID should be present")

				// Verify stream status transition
				stream := data["stream"].(map[string]interface{})
				assert.Equal(t, "live", stream["status"], "Stream status should be 'live'")
				assert.NotEmpty(t, stream["actual_start_time"], "Actual start time should be set")
				assert.Equal(t, float64(0), stream["viewer_count"], "Initial viewer count should be 0")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Invalid WebRTC SDP Offer
	// Verifies that malformed SDP offers are rejected with proper error messages
	t.Run("StartStream_InvalidOffer", func(t *testing.T) {
		// Define the expected interaction for invalid SDP
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 exists and is scheduled").
			UponReceiving("a request to start stream with invalid WebRTC offer").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/start", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/start$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-broadcaster", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					// Invalid/malformed SDP offer
					"sdp_offer":       "invalid sdp format without proper headers",
					"quality_layers": []string{"720p"},
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Invalid WebRTC SDP offer format"),
					"details": map[string]interface{}{
						"sdp_offer": matchers.Like("SDP offer must be valid Session Description Protocol format"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/start", config.Host, config.Port, streamID)

				invalidReq := map[string]interface{}{
					"sdp_offer":       "invalid sdp format without proper headers",
					"quality_layers": []string{"720p"},
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
				req.Header.Set("Authorization", "Bearer test-token-broadcaster")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 400 Bad Request
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for invalid SDP")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "Invalid", "Error should mention invalid SDP")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Stream Already Live
	// Verifies that attempting to start an already-live stream returns 409 Conflict
	t.Run("StartStream_AlreadyLive", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440010 exists and is already live").
			UponReceiving("a request to start an already-live stream").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440010/start", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/start$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-broadcaster", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					"sdp_offer":       "v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n",
					"quality_layers": []string{"720p"},
				})
			}).
			WillRespondWith(409, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success":        false,
					"error":          matchers.Like("Stream is already live"),
					"current_status": matchers.Like("live"),
					"stream_id":      matchers.Regex("550e8400-e29b-41d4-a716-446655440010", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"actual_start_time": matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440010"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/start", config.Host, config.Port, streamID)

				startReq := map[string]interface{}{
					"sdp_offer":       "v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n",
					"quality_layers": []string{"720p"},
				}

				jsonData, err := json.Marshal(startReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
				assert.Equal(t, 409, resp.StatusCode, "Expected 409 Conflict for already-live stream")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Equal(t, "live", responseBody["current_status"], "Should indicate stream is live")
				assert.NotEmpty(t, responseBody["actual_start_time"], "Should show when stream started")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Non-Broadcaster Attempting to Start Stream
	// Verifies that only the broadcaster can start their own stream (403 Forbidden)
	t.Run("StartStream_NotOwner", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-446655440000 exists and user is not broadcaster").
			UponReceiving("a request to start stream from non-broadcaster").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-446655440000/start", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/start$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-other-user", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					"sdp_offer":       "v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n",
					"quality_layers": []string{"720p"},
				})
			}).
			WillRespondWith(403, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Only the broadcaster can start their stream"),
					"stream_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440000", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"broadcaster_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-446655440000"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/start", config.Host, config.Port, streamID)

				startReq := map[string]interface{}{
					"sdp_offer":       "v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n",
					"quality_layers": []string{"720p"},
				}

				jsonData, err := json.Marshal(startReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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

	// Test Case 5: Stream Not Found
	// Verifies that attempting to start a non-existent stream returns 404 Not Found
	t.Run("StartStream_NotFound", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("stream 550e8400-e29b-41d4-a716-999999999999 does not exist").
			UponReceiving("a request to start non-existent stream").
			WithRequest("POST", matchers.Regex("/api/v1/streams/550e8400-e29b-41d4-a716-999999999999/start", `^/api/v1/streams/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/start$`), func(b *consumer.V2RequestBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.Header("Authorization", matchers.Regex("Bearer test-token-broadcaster", `^Bearer .+$`))
				b.JSONBody(map[string]interface{}{
					"sdp_offer":       "v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n",
					"quality_layers": []string{"720p"},
				})
			}).
			WillRespondWith(404, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Stream not found"),
					"stream_id": matchers.Regex("550e8400-e29b-41d4-a716-999999999999", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				streamID := "550e8400-e29b-41d4-a716-999999999999"
				url := fmt.Sprintf("http://%s:%d/api/v1/streams/%s/start", config.Host, config.Port, streamID)

				startReq := map[string]interface{}{
					"sdp_offer":       "v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n",
					"quality_layers": []string{"720p"},
				}

				jsonData, err := json.Marshal(startReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}