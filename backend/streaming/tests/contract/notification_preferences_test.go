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

// TestNotificationPreferencesPactConsumer runs consumer Pact tests for GET/PUT /api/v1/notification-preferences endpoint
// This test defines the contract between streaming clients (web/mobile) and the streaming service
// for managing user notification preferences for live stream notifications.
//
// This test follows TDD principles and MUST FAIL initially because:
// 1. No handler implementation exists for GET/PUT /notification-preferences endpoint
// 2. No NotificationPreference model with helper methods (IsQuietHoursActive, ShouldNotify) exists
// 3. No default preference logic for first-time users is implemented
// 4. No quiet hours validation logic (start/end time format and logic) is implemented
// 5. No channel toggle logic (push/email/sms/in_app) is implemented
func TestNotificationPreferencesPactConsumer(t *testing.T) {
	// Create Pact consumer with mock provider configuration
	// Consumer: streaming-web-client (represents web frontend)
	// Provider: streaming-service (represents backend streaming service)
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "streaming-web-client",
		Provider: "streaming-service",
		Host:     "localhost",
		Port:     9001, // Same port as other streaming service tests
		PactDir:  "./pacts", // Output directory for generated pact files
	})
	assert.NoError(t, err, "Failed to create Pact mock provider")

	// Test Case 1: Get Preferences - Success
	// Verifies that authenticated users can retrieve their notification preferences
	// Uses JWT token to identify user (no path parameter needed)
	t.Run("GetPreferences_Success", func(t *testing.T) {
		// Define the expected interaction
		err := mockProvider.
			AddInteraction().
			Given("user exists with ID 550e8400-e29b-41d4-a716-446655440001 with configured preferences").
			UponReceiving("a request to get notification preferences").
			WithRequest("GET", "/api/v1/notification-preferences", func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-user"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				// Response headers
				b.Header("Content-Type", matchers.String("application/json"))

				// Response body with matchers - NotificationPreference entity
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Notification preferences retrieved successfully"),
					"data": map[string]interface{}{
						// User identity from JWT
						"user_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440001", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						// Channel preferences - all boolean toggles
						"push_enabled":   matchers.Like(true),
						"email_enabled":  matchers.Like(true),
						"sms_enabled":    matchers.Like(false),
						"in_app_enabled": matchers.Like(true),

						// Context-specific preferences
						"follow_only_mode": matchers.Like(false), // If true, only notify for followed broadcasters

						// Quiet hours configuration - TIME format (HH:MM:SS)
						"quiet_hours_start": matchers.Regex("22:00:00", `^([01]\d|2[0-3]):([0-5]\d):([0-5]\d)$`),
						"quiet_hours_end":   matchers.Regex("08:00:00", `^([01]\d|2[0-3]):([0-5]\d):([0-5]\d)$`),

						// Timezone for quiet hours calculation
						"timezone": matchers.Like("Asia/Bangkok"),

						// Timestamp
						"updated_at": matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Construct the API endpoint URL
				url := fmt.Sprintf("http://%s:%d/api/v1/notification-preferences", config.Host, config.Port)

				// Create HTTP GET request
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				// Set required headers - JWT contains user_id
				req.Header.Set("Authorization", "Bearer test-token-user")

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

				// Verify preference fields
				assert.NotNil(t, data["user_id"], "Should have user_id")
				assert.NotNil(t, data["push_enabled"], "Should have push_enabled")
				assert.NotNil(t, data["quiet_hours_start"], "Should have quiet_hours_start")
				assert.Equal(t, "Asia/Bangkok", data["timezone"], "Should have timezone")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 2: Get Preferences - Default Settings (First-Time User)
	// Verifies that users without configured preferences receive default values
	// All channels enabled by default, no quiet hours
	t.Run("GetPreferences_DefaultSettings", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("user exists with ID 550e8400-e29b-41d4-a716-446655440002 without configured preferences").
			UponReceiving("a request to get notification preferences for first-time user").
			WithRequest("GET", "/api/v1/notification-preferences", func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-new-user"))
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Default notification preferences created"),
					"data": map[string]interface{}{
						"user_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440002", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						// Default values - all channels enabled
						"push_enabled":   matchers.Like(true),
						"email_enabled":  matchers.Like(true),
						"sms_enabled":    matchers.Like(true),
						"in_app_enabled": matchers.Like(true),

						// Default - notify for all broadcasters
						"follow_only_mode": matchers.Like(false),

						// No quiet hours by default (null values)
						"quiet_hours_start": nil,
						"quiet_hours_end":   nil,

						// Default timezone from user profile
						"timezone": matchers.Like("UTC"),

						"updated_at": matchers.Regex("2025-09-30T14:00:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/notification-preferences", config.Host, config.Port)

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Authorization", "Bearer test-token-new-user")

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
				// Verify all channels enabled by default
				assert.True(t, data["push_enabled"].(bool), "Push should be enabled by default")
				assert.True(t, data["email_enabled"].(bool), "Email should be enabled by default")
				assert.True(t, data["in_app_enabled"].(bool), "In-app should be enabled by default")

				// Verify no quiet hours by default
				assert.Nil(t, data["quiet_hours_start"], "Should have no quiet hours start by default")
				assert.Nil(t, data["quiet_hours_end"], "Should have no quiet hours end by default")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 3: Update Preferences - Disable Channels
	// Verifies that users can selectively disable notification channels
	t.Run("UpdatePreferences_DisableChannels", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("user exists with ID 550e8400-e29b-41d4-a716-446655440003 with default preferences").
			UponReceiving("a request to disable push and email notifications").
			WithRequest("PUT", "/api/v1/notification-preferences", func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-user"))
				b.JSONBody(map[string]interface{}{
					"push_enabled":  false,
					"email_enabled": false,
					"sms_enabled":   false,
					"in_app_enabled": true, // Keep in-app enabled
				})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Notification preferences updated successfully"),
					"data": map[string]interface{}{
						"user_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440003", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						// Updated channel preferences
						"push_enabled":   matchers.Like(false),
						"email_enabled":  matchers.Like(false),
						"sms_enabled":    matchers.Like(false),
						"in_app_enabled": matchers.Like(true),

						"follow_only_mode": matchers.Like(false),

						"quiet_hours_start": nil,
						"quiet_hours_end":   nil,
						"timezone":          matchers.Like("UTC"),

						"updated_at": matchers.Regex("2025-09-30T14:05:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/notification-preferences", config.Host, config.Port)

				updatePrefsReq := map[string]interface{}{
					"push_enabled":   false,
					"email_enabled":  false,
					"sms_enabled":    false,
					"in_app_enabled": true,
				}

				jsonData, err := json.Marshal(updatePrefsReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-user")

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
				// Verify channels were disabled
				assert.False(t, data["push_enabled"].(bool), "Push should be disabled")
				assert.False(t, data["email_enabled"].(bool), "Email should be disabled")
				assert.False(t, data["sms_enabled"].(bool), "SMS should be disabled")
				assert.True(t, data["in_app_enabled"].(bool), "In-app should remain enabled")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 4: Update Preferences - Set Quiet Hours
	// Verifies that users can set quiet hours to suppress notifications during sleep time
	t.Run("UpdatePreferences_QuietHours", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("user exists with ID 550e8400-e29b-41d4-a716-446655440004 with default preferences").
			UponReceiving("a request to set quiet hours from 22:00 to 08:00").
			WithRequest("PUT", "/api/v1/notification-preferences", func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-user"))
				b.JSONBody(map[string]interface{}{
					"quiet_hours_start": "22:00:00",
					"quiet_hours_end":   "08:00:00",
					"timezone":          "Asia/Bangkok",
				})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Notification preferences updated successfully"),
					"data": map[string]interface{}{
						"user_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440004", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						"push_enabled":   matchers.Like(true),
						"email_enabled":  matchers.Like(true),
						"sms_enabled":    matchers.Like(true),
						"in_app_enabled": matchers.Like(true),

						"follow_only_mode": matchers.Like(false),

						// Updated quiet hours
						"quiet_hours_start": matchers.Regex("22:00:00", `^([01]\d|2[0-3]):([0-5]\d):([0-5]\d)$`),
						"quiet_hours_end":   matchers.Regex("08:00:00", `^([01]\d|2[0-3]):([0-5]\d):([0-5]\d)$`),
						"timezone":          matchers.Like("Asia/Bangkok"),

						"updated_at": matchers.Regex("2025-09-30T14:10:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/notification-preferences", config.Host, config.Port)

				updatePrefsReq := map[string]interface{}{
					"quiet_hours_start": "22:00:00",
					"quiet_hours_end":   "08:00:00",
					"timezone":          "Asia/Bangkok",
				}

				jsonData, err := json.Marshal(updatePrefsReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-user")

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
				// Verify quiet hours were set
				assert.Equal(t, "22:00:00", data["quiet_hours_start"], "Quiet hours start should be set")
				assert.Equal(t, "08:00:00", data["quiet_hours_end"], "Quiet hours end should be set")
				assert.Equal(t, "Asia/Bangkok", data["timezone"], "Timezone should be set")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 5: Update Preferences - Enable Follow Only Mode
	// Verifies that users can enable follow_only_mode to only receive notifications from followed broadcasters
	t.Run("UpdatePreferences_FollowOnly", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("user exists with ID 550e8400-e29b-41d4-a716-446655440005 with default preferences").
			UponReceiving("a request to enable follow_only_mode").
			WithRequest("PUT", "/api/v1/notification-preferences", func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-user"))
				b.JSONBody(map[string]interface{}{
					"follow_only_mode": true,
				})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Notification preferences updated successfully"),
					"data": map[string]interface{}{
						"user_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440005", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						"push_enabled":   matchers.Like(true),
						"email_enabled":  matchers.Like(true),
						"sms_enabled":    matchers.Like(true),
						"in_app_enabled": matchers.Like(true),

						// Updated field - only notify for followed broadcasters
						"follow_only_mode": matchers.Like(true),

						"quiet_hours_start": nil,
						"quiet_hours_end":   nil,
						"timezone":          matchers.Like("UTC"),

						"updated_at": matchers.Regex("2025-09-30T14:15:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/notification-preferences", config.Host, config.Port)

				updatePrefsReq := map[string]interface{}{
					"follow_only_mode": true,
				}

				jsonData, err := json.Marshal(updatePrefsReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-user")

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
				// Verify follow_only_mode was enabled
				assert.True(t, data["follow_only_mode"].(bool), "Follow only mode should be enabled")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 6: Update Preferences - Invalid Quiet Hours (400 Bad Request)
	// Verifies validation of quiet hours time format and logical consistency
	t.Run("UpdatePreferences_InvalidTimes", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("user exists with ID 550e8400-e29b-41d4-a716-446655440006 with default preferences").
			UponReceiving("a request to set invalid quiet hours format").
			WithRequest("PUT", "/api/v1/notification-preferences", func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-user"))
				b.JSONBody(map[string]interface{}{
					"quiet_hours_start": "25:00:00", // Invalid hour (>23)
					"quiet_hours_end":   "08:00:00",
				})
			}).
			WillRespondWith(400, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Invalid quiet hours time format"),
					"details": map[string]interface{}{
						"field":          matchers.Like("quiet_hours_start"),
						"invalid_value":  matchers.Like("25:00:00"),
						"expected_format": matchers.Like("HH:MM:SS (00:00:00 - 23:59:59)"),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/notification-preferences", config.Host, config.Port)

				updatePrefsReq := map[string]interface{}{
					"quiet_hours_start": "25:00:00",
					"quiet_hours_end":   "08:00:00",
				}

				jsonData, err := json.Marshal(updatePrefsReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-user")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 400 Bad Request
				assert.Equal(t, 400, resp.StatusCode, "Expected 400 Bad Request for invalid time format")

				var responseBody map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
					return fmt.Errorf("failed to decode response body: %w", err)
				}

				// Verify error response structure
				assert.False(t, responseBody["success"].(bool), "Response should indicate failure")
				assert.NotEmpty(t, responseBody["error"], "Error message should be present")
				assert.Contains(t, responseBody["error"].(string), "Invalid", "Error should mention invalid format")

				details := responseBody["details"].(map[string]interface{})
				assert.Equal(t, "quiet_hours_start", details["field"], "Should identify problematic field")
				assert.NotNil(t, details["expected_format"], "Should provide expected format guidance")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})

	// Test Case 7: Update Preferences - Unauthorized (401)
	// Verifies that requests without valid JWT are rejected
	t.Run("UpdatePreferences_Unauthorized", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("no authentication token provided").
			UponReceiving("a request to update preferences without authentication").
			WithRequest("PUT", "/api/v1/notification-preferences", func(b *consumer.V2RequestBuilder) {
				// No Authorization header
				b.JSONBody(map[string]interface{}{
					"push_enabled": false,
				})
			}).
			WillRespondWith(401, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": false,
					"error":   matchers.Like("Missing or invalid JWT token"),
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/notification-preferences", config.Host, config.Port)

				updatePrefsReq := map[string]interface{}{
					"push_enabled": false,
				}

				jsonData, err := json.Marshal(updatePrefsReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				// No Authorization header

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return fmt.Errorf("failed to execute HTTP request: %w", err)
				}
				defer resp.Body.Close()

				// Verify 401 Unauthorized
				assert.Equal(t, 401, resp.StatusCode, "Expected 401 Unauthorized without authentication")

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

	// Test Case 8: Update Preferences - Clear Quiet Hours
	// Verifies that users can remove quiet hours by setting them to null
	t.Run("UpdatePreferences_ClearQuietHours", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("user exists with ID 550e8400-e29b-41d4-a716-446655440007 with quiet hours configured").
			UponReceiving("a request to clear quiet hours by setting to null").
			WithRequest("PUT", "/api/v1/notification-preferences", func(b *consumer.V2RequestBuilder) {
				b.Header("Authorization", matchers.String("Bearer test-token-user"))
				b.JSONBody(map[string]interface{}{
					"quiet_hours_start": nil,
					"quiet_hours_end":   nil,
				})
			}).
			WillRespondWith(200, func(b *consumer.V2ResponseBuilder) {
				b.Header("Content-Type", matchers.String("application/json"))
				b.JSONBody(map[string]interface{}{
					"success": true,
					"message": matchers.Like("Notification preferences updated successfully"),
					"data": map[string]interface{}{
						"user_id": matchers.Regex("550e8400-e29b-41d4-a716-446655440007", `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),

						"push_enabled":   matchers.Like(true),
						"email_enabled":  matchers.Like(true),
						"sms_enabled":    matchers.Like(true),
						"in_app_enabled": matchers.Like(true),

						"follow_only_mode": matchers.Like(false),

						// Cleared quiet hours
						"quiet_hours_start": nil,
						"quiet_hours_end":   nil,
						"timezone":          matchers.Like("Asia/Bangkok"),

						"updated_at": matchers.Regex("2025-09-30T14:20:00Z", `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`),
					},
				})
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/notification-preferences", config.Host, config.Port)

				updatePrefsReq := map[string]interface{}{
					"quiet_hours_start": nil,
					"quiet_hours_end":   nil,
				}

				jsonData, err := json.Marshal(updatePrefsReq)
				if err != nil {
					return fmt.Errorf("failed to marshal request body: %w", err)
				}

				req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return fmt.Errorf("failed to create HTTP request: %w", err)
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token-user")

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
				// Verify quiet hours were cleared
				assert.Nil(t, data["quiet_hours_start"], "Quiet hours start should be cleared")
				assert.Nil(t, data["quiet_hours_end"], "Quiet hours end should be cleared")

				return nil
			})

		assert.NoError(t, err, "Pact interaction should complete successfully")
	})
}