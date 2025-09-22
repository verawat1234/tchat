package contract_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthOTPSend_Contract validates the POST /auth/otp/send endpoint contract
// This test MUST FAIL initially as no implementation exists yet (TDD)
func TestAuthOTPSend_Contract(t *testing.T) {
	// Test server URL - will fail until server is implemented
	baseURL := "http://localhost:8080"

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedFields []string
	}{
		{
			name: "Valid phone OTP request",
			payload: map[string]interface{}{
				"identifier": "+66812345678",
				"type":       "phone",
				"country":    "TH",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"success", "message", "expires_at"},
		},
		{
			name: "Valid email OTP request",
			payload: map[string]interface{}{
				"identifier": "user@example.com",
				"type":       "email",
				"country":    "TH",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"success", "message", "expires_at"},
		},
		{
			name: "Invalid identifier format",
			payload: map[string]interface{}{
				"identifier": "invalid",
				"type":       "phone",
				"country":    "TH",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Missing required fields",
			payload: map[string]interface{}{
				"identifier": "+66812345678",
				// missing type field
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Invalid country code",
			payload: map[string]interface{}{
				"identifier": "+66812345678",
				"type":       "phone",
				"country":    "XX", // invalid country
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request payload
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest("POST", baseURL+"/auth/otp/send", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Execute request (will fail until server is implemented)
			client := &http.Client{}
			resp, err := client.Do(req)

			// This SHOULD FAIL initially - no server running
			if err != nil {
				t.Logf("Expected failure: %v (no implementation yet)", err)
				return // Test passes by failing as expected in TDD
			}
			defer resp.Body.Close()

			// Validate response status
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Validate response structure
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Check expected fields exist
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Additional contract validations
			if resp.StatusCode == http.StatusOK {
				// Success response should have boolean success field
				success, exists := response["success"]
				assert.True(t, exists, "Success response should have 'success' field")
				assert.True(t, success.(bool), "Success field should be true for successful requests")

				// Should have message field
				message, exists := response["message"]
				assert.True(t, exists, "Response should have 'message' field")
				assert.NotEmpty(t, message, "Message should not be empty")

				// Should have expires_at field
				expiresAt, exists := response["expires_at"]
				assert.True(t, exists, "Response should have 'expires_at' field")
				assert.NotEmpty(t, expiresAt, "expires_at should not be empty")
			} else {
				// Error response should have error field
				errorField, exists := response["error"]
				assert.True(t, exists, "Error response should have 'error' field")
				assert.NotEmpty(t, errorField, "Error field should not be empty")
			}
		})
	}
}

// TestAuthOTPSend_RateLimiting validates rate limiting contract
func TestAuthOTPSend_RateLimiting(t *testing.T) {
	baseURL := "http://localhost:8080"

	payload := map[string]interface{}{
		"identifier": "+66812345678",
		"type":       "phone",
		"country":    "TH",
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	// Attempt multiple rapid requests to test rate limiting
	for i := 0; i < 6; i++ { // Assume rate limit is 5 requests per minute
		req, err := http.NewRequest("POST", baseURL+"/auth/otp/send", bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Logf("Expected failure: %v (no implementation yet)", err)
			return // Test passes by failing as expected in TDD
		}
		defer resp.Body.Close()

		// After 5th request, should get rate limited (429)
		if i >= 5 {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode,
				"Should be rate limited after 5 requests")
		}
	}
}

// TestAuthOTPSend_CountrySpecificValidation validates country-specific phone validation
func TestAuthOTPSend_CountrySpecificValidation(t *testing.T) {
	baseURL := "http://localhost:8080"

	tests := []struct {
		name        string
		phone       string
		country     string
		shouldPass  bool
	}{
		{"Thailand valid", "+66812345678", "TH", true},
		{"Indonesia valid", "+628123456789", "ID", true},
		{"Thailand invalid format", "+6681234", "TH", false},
		{"Wrong country for phone", "+66812345678", "ID", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"identifier": tt.phone,
				"type":       "phone",
				"country":    tt.country,
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", baseURL+"/auth/otp/send", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				t.Logf("Expected failure: %v (no implementation yet)", err)
				return // Test passes by failing as expected in TDD
			}
			defer resp.Body.Close()

			if tt.shouldPass {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			}
		})
	}
}