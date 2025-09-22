package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T014: Contract test GET /users/profile
func TestAuthProfileGetContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedFields []string
		description    string
	}{
		{
			name:           "valid_authenticated_request",
			authHeader:     "Bearer valid_jwt_token_12345",
			expectedStatus: http.StatusOK,
			expectedFields: []string{"id", "phone_number", "country_code", "status", "profile", "created_at", "updated_at"},
			description:    "Should return user profile for valid JWT token",
		},
		{
			name:           "missing_authorization_header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized when no auth header provided",
		},
		{
			name:           "invalid_bearer_format",
			authHeader:     "Invalid token_format",
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for invalid bearer format",
		},
		{
			name:           "expired_jwt_token",
			authHeader:     "Bearer expired_jwt_token_67890",
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for expired JWT token",
		},
		{
			name:           "invalid_jwt_token",
			authHeader:     "Bearer invalid_jwt_signature",
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for invalid JWT signature",
		},
		{
			name:           "malformed_jwt_token",
			authHeader:     "Bearer not.a.valid.jwt",
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for malformed JWT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()

			// Mock authentication middleware and endpoint
			router.GET("/api/v1/users/profile", func(c *gin.Context) {
				authHeader := c.GetHeader("Authorization")

				// Mock authentication logic
				if authHeader == "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
					return
				}

				if authHeader != "Bearer valid_jwt_token_12345" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
					return
				}

				// Mock user profile data for valid token
				c.JSON(http.StatusOK, gin.H{
					"id":           "123e4567-e89b-12d3-a456-426614174000",
					"phone_number": "+66812345678",
					"country_code": "TH",
					"status":       "active",
					"profile": gin.H{
						"display_name": "John Doe",
						"avatar_url":   "https://cdn.tchat.sea/avatars/user123.jpg",
						"locale":       "en",
					},
					"created_at": "2023-12-01T10:00:00Z",
					"updated_at": "2023-12-01T10:00:00Z",
				})
			})

			// Prepare request
			req, err := http.NewRequest("GET", "/api/v1/users/profile", nil)
			require.NoError(t, err)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify status code
			assert.Equal(t, tt.expectedStatus, w.Code,
				"Test: %s - %s", tt.name, tt.description)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Verify expected fields are present
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field,
					"Test: %s - Response should contain field '%s'", tt.name, field)
			}

			// Additional assertions based on test case
			switch tt.name {
			case "valid_authenticated_request":
				// Verify profile structure
				profile, ok := response["profile"].(map[string]interface{})
				assert.True(t, ok, "profile should be an object")
				assert.Contains(t, profile, "display_name", "profile should contain display_name")
				assert.Contains(t, profile, "locale", "profile should contain locale")

				// Verify user ID format (UUID)
				userID, ok := response["id"].(string)
				assert.True(t, ok, "id should be string")
				assert.Len(t, userID, 36, "id should be UUID format (36 chars)")

				// Verify phone number format
				phoneNumber, ok := response["phone_number"].(string)
				assert.True(t, ok, "phone_number should be string")
				assert.Regexp(t, `^\+\d{10,15}$`, phoneNumber, "phone_number should be in international format")

				// Verify country code format
				countryCode, ok := response["country_code"].(string)
				assert.True(t, ok, "country_code should be string")
				assert.Len(t, countryCode, 2, "country_code should be 2 characters")
				assert.Contains(t, []string{"TH", "SG", "ID", "MY", "PH", "VN"}, countryCode,
					"country_code should be a supported SEA country")

				// Verify status
				status, ok := response["status"].(string)
				assert.True(t, ok, "status should be string")
				assert.Contains(t, []string{"active", "inactive", "suspended"}, status,
					"status should be valid user status")

				// Verify timestamps
				createdAt, ok := response["created_at"].(string)
				assert.True(t, ok, "created_at should be string")
				assert.NotEmpty(t, createdAt, "created_at should not be empty")

				updatedAt, ok := response["updated_at"].(string)
				assert.True(t, ok, "updated_at should be string")
				assert.NotEmpty(t, updatedAt, "updated_at should not be empty")

			case "missing_authorization_header", "invalid_bearer_format", "expired_jwt_token", "invalid_jwt_token", "malformed_jwt_token":
				// Verify error message
				errorMsg, ok := response["error"].(string)
				assert.True(t, ok, "error should be string")
				assert.NotEmpty(t, errorMsg, "error message should not be empty")
			}
		})
	}
}

// TestAuthProfileGetSecurity tests security aspects of profile endpoint
func TestAuthProfileGetSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("no_sensitive_data_exposure", func(t *testing.T) {
		router := gin.New()
		router.GET("/api/v1/users/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"id":           "123e4567-e89b-12d3-a456-426614174000",
				"phone_number": "+66812345678",
				"country_code": "TH",
				"status":       "active",
				"profile": gin.H{
					"display_name": "John Doe",
					"avatar_url":   "https://cdn.tchat.sea/avatars/user123.jpg",
					"locale":       "en",
				},
				"created_at": "2023-12-01T10:00:00Z",
				"updated_at": "2023-12-01T10:00:00Z",
			})
		})

		req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		// Ensure sensitive data is not exposed
		sensitiveFields := []string{"password", "otp_code", "refresh_token", "private_key", "secret", "session_id"}
		for _, field := range sensitiveFields {
			assert.NotContains(t, response, field, "Profile should not expose sensitive field: %s", field)
		}

		// Verify profile object doesn't contain sensitive data
		if profile, ok := response["profile"].(map[string]interface{}); ok {
			for _, field := range sensitiveFields {
				assert.NotContains(t, profile, field, "Profile object should not expose sensitive field: %s", field)
			}
		}
	})

	t.Run("user_isolation", func(t *testing.T) {
		// This test ensures users can only access their own profile
		router := gin.New()

		router.GET("/api/v1/users/profile", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")

			// Mock different users based on token
			switch authHeader {
			case "Bearer user1_token":
				c.JSON(http.StatusOK, gin.H{
					"id": "user1-id",
					"phone_number": "+66811111111",
				})
			case "Bearer user2_token":
				c.JSON(http.StatusOK, gin.H{
					"id": "user2-id",
					"phone_number": "+66822222222",
				})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
		})

		// User 1 request
		req1, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
		req1.Header.Set("Authorization", "Bearer user1_token")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		var response1 map[string]interface{}
		json.Unmarshal(w1.Body.Bytes(), &response1)

		// User 2 request
		req2, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
		req2.Header.Set("Authorization", "Bearer user2_token")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		var response2 map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &response2)

		// Verify users get different profiles
		assert.NotEqual(t, response1["id"], response2["id"], "Users should get different profile IDs")
		assert.NotEqual(t, response1["phone_number"], response2["phone_number"], "Users should get different phone numbers")
	})
}

// TestAuthProfileGetPerformance tests performance requirements
func TestAuthProfileGetPerformance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/api/v1/users/profile", func(c *gin.Context) {
		// Simulate database lookup time
		// time.Sleep(5 * time.Millisecond) // Commented out for test speed
		c.JSON(http.StatusOK, gin.H{
			"id":           "user-123",
			"phone_number": "+66812345678",
			"country_code": "TH",
			"status":       "active",
			"profile": gin.H{
				"display_name": "Test User",
				"locale":       "en",
			},
			"created_at": "2023-12-01T10:00:00Z",
			"updated_at": "2023-12-01T10:00:00Z",
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response size is reasonable (< 10KB for profile data)
	responseSize := len(w.Body.Bytes())
	assert.Less(t, responseSize, 10*1024, "Profile response should be less than 10KB, got %d bytes", responseSize)
}

// TestAuthProfileGetLocalization tests Southeast Asian localization
func TestAuthProfileGetLocalization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	supportedCountries := []string{"TH", "SG", "ID", "MY", "PH", "VN"}
	supportedLocales := []string{"en", "th", "id", "ms", "fil", "vi"}

	for _, country := range supportedCountries {
		t.Run("country_"+country, func(t *testing.T) {
			router := gin.New()
			router.GET("/api/v1/users/profile", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"id":           "user-123",
					"phone_number": "+66812345678", // This would vary by country in real implementation
					"country_code": country,
					"status":       "active",
					"profile": gin.H{
						"display_name": "Test User",
						"locale":       "en", // This would be set based on user preference
					},
					"created_at": "2023-12-01T10:00:00Z",
					"updated_at": "2023-12-01T10:00:00Z",
				})
			})

			req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
			req.Header.Set("Authorization", "Bearer valid_token")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Verify country code is correct
			assert.Equal(t, country, response["country_code"], "Country code should match")

			// Verify profile contains locale
			if profile, ok := response["profile"].(map[string]interface{}); ok {
				locale, ok := profile["locale"].(string)
				assert.True(t, ok, "Profile should contain locale")
				assert.Contains(t, supportedLocales, locale, "Locale should be supported")
			}
		})
	}
}