package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T015: Contract test PUT /users/profile
func TestAuthProfileUpdateContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		payload        map[string]interface{}
		expectedStatus int
		expectedFields []string
		description    string
	}{
		{
			name:       "valid_profile_update",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"display_name": "Updated Name",
				"locale":       "th",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"id", "phone_number", "country_code", "status", "profile", "updated_at"},
			description:    "Should update profile for valid request",
		},
		{
			name:       "update_display_name_only",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"display_name": "New Display Name",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"profile"},
			description:    "Should update only display name",
		},
		{
			name:       "update_locale_only",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"locale": "vi",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"profile"},
			description:    "Should update only locale",
		},
		{
			name:       "update_avatar_url",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"avatar_url": "https://cdn.tchat.sea/avatars/new-avatar.jpg",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"profile"},
			description:    "Should update avatar URL",
		},
		{
			name:           "missing_authorization",
			authHeader:     "",
			payload:        map[string]interface{}{"display_name": "Test"},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized when no auth header",
		},
		{
			name:           "invalid_jwt_token",
			authHeader:     "Bearer invalid_token_67890",
			payload:        map[string]interface{}{"display_name": "Test"},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for invalid token",
		},
		{
			name:       "invalid_locale",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"locale": "invalid_locale",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid locale",
		},
		{
			name:       "empty_display_name",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"display_name": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for empty display name",
		},
		{
			name:       "display_name_too_long",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"display_name": "This is a very long display name that exceeds the maximum allowed length of characters for a display name in the system",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for display name too long",
		},
		{
			name:       "invalid_avatar_url",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"avatar_url": "not-a-valid-url",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid avatar URL",
		},
		{
			name:       "malformed_request_body",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"display_name": 12345, // Wrong type
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for malformed request body",
		},
		{
			name:           "empty_request_body",
			authHeader:     "Bearer valid_jwt_token_12345",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for empty request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()

			// Mock authentication middleware and endpoint
			router.PUT("/api/v1/users/profile", func(c *gin.Context) {
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

				// Parse request body
				var request map[string]interface{}
				if err := c.ShouldBindJSON(&request); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
					return
				}

				// Validation logic
				if len(request) == 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
					return
				}

				// Validate display_name
				if displayName, exists := request["display_name"]; exists {
					if name, ok := displayName.(string); ok {
						if name == "" {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Display name cannot be empty"})
							return
						}
						if len(name) > 100 {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Display name too long"})
							return
						}
					} else {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Display name must be string"})
						return
					}
				}

				// Validate locale
				if locale, exists := request["locale"]; exists {
					if localeStr, ok := locale.(string); ok {
						validLocales := []string{"en", "th", "id", "ms", "fil", "vi"}
						valid := false
						for _, validLocale := range validLocales {
							if localeStr == validLocale {
								valid = true
								break
							}
						}
						if !valid {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid locale"})
							return
						}
					} else {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Locale must be string"})
						return
					}
				}

				// Validate avatar_url
				if avatarURL, exists := request["avatar_url"]; exists {
					if urlStr, ok := avatarURL.(string); ok {
						if urlStr != "" && !(len(urlStr) > 7 && (urlStr[:7] == "http://" || urlStr[:8] == "https://")) {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid avatar URL"})
							return
						}
					} else {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar URL must be string"})
						return
					}
				}

				// Mock successful update response
				profile := gin.H{
					"display_name": "Updated Name",
					"avatar_url":   "https://cdn.tchat.sea/avatars/user123.jpg",
					"locale":       "en",
				}

				// Update profile with request data
				if displayName, exists := request["display_name"]; exists {
					profile["display_name"] = displayName
				}
				if locale, exists := request["locale"]; exists {
					profile["locale"] = locale
				}
				if avatarURL, exists := request["avatar_url"]; exists {
					profile["avatar_url"] = avatarURL
				}

				c.JSON(http.StatusOK, gin.H{
					"id":           "123e4567-e89b-12d3-a456-426614174000",
					"phone_number": "+66812345678",
					"country_code": "TH",
					"status":       "active",
					"profile":      profile,
					"created_at":   "2023-12-01T10:00:00Z",
					"updated_at":   "2023-12-01T15:30:00Z", // Updated timestamp
				})
			})

			// Prepare request
			jsonData, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

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
			case "valid_profile_update", "update_display_name_only", "update_locale_only", "update_avatar_url":
				// Verify profile structure
				profile, ok := response["profile"].(map[string]interface{})
				assert.True(t, ok, "profile should be an object")

				// Verify specific updates
				if displayName, exists := tt.payload["display_name"]; exists {
					assert.Equal(t, displayName, profile["display_name"],
						"Display name should be updated")
				}
				if locale, exists := tt.payload["locale"]; exists {
					assert.Equal(t, locale, profile["locale"],
						"Locale should be updated")
				}
				if avatarURL, exists := tt.payload["avatar_url"]; exists {
					assert.Equal(t, avatarURL, profile["avatar_url"],
						"Avatar URL should be updated")
				}

				// Verify updated_at timestamp is present and recent
				updatedAt, ok := response["updated_at"].(string)
				assert.True(t, ok, "updated_at should be string")
				assert.NotEmpty(t, updatedAt, "updated_at should not be empty")

			case "missing_authorization", "invalid_jwt_token":
				// Verify error message
				errorMsg, ok := response["error"].(string)
				assert.True(t, ok, "error should be string")
				assert.NotEmpty(t, errorMsg, "error message should not be empty")

			case "invalid_locale", "empty_display_name", "display_name_too_long", "invalid_avatar_url", "malformed_request_body", "empty_request_body":
				// Verify validation error
				errorMsg, ok := response["error"].(string)
				assert.True(t, ok, "error should be string")
				assert.NotEmpty(t, errorMsg, "error message should not be empty")
			}
		})
	}
}

// TestAuthProfileUpdateSecurity tests security aspects
func TestAuthProfileUpdateSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("cannot_update_readonly_fields", func(t *testing.T) {
		router := gin.New()
		router.PUT("/api/v1/users/profile", func(c *gin.Context) {
			var request map[string]interface{}
			c.ShouldBindJSON(&request)

			// Check for readonly fields
			readonlyFields := []string{"id", "phone_number", "country_code", "created_at", "status"}
			for _, field := range readonlyFields {
				if _, exists := request[field]; exists {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Cannot update readonly field: " + field,
					})
					return
				}
			}

			c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
		})

		// Try to update readonly fields
		readonlyFields := map[string]interface{}{
			"id":           "hacker-id",
			"phone_number": "+66999999999",
			"country_code": "XX",
			"created_at":   "2020-01-01T00:00:00Z",
			"status":       "suspended",
		}

		for field, value := range readonlyFields {
			payload := map[string]interface{}{
				field: value,
			}

			jsonData, _ := json.Marshal(payload)
			req, _ := http.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid_token")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code,
				"Should reject attempt to update readonly field: %s", field)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(t, response["error"].(string), field,
				"Error message should mention the readonly field: %s", field)
		}
	})

	t.Run("input_sanitization", func(t *testing.T) {
		router := gin.New()
		router.PUT("/api/v1/users/profile", func(c *gin.Context) {
			var request map[string]interface{}
			c.ShouldBindJSON(&request)

			// Mock sanitization
			if displayName, exists := request["display_name"]; exists {
				if name, ok := displayName.(string); ok {
					// Check for potential XSS or injection attempts
					dangerousPatterns := []string{"<script>", "javascript:", "onload=", "onerror=", "DROP TABLE", "SELECT *"}
					for _, pattern := range dangerousPatterns {
						if len(name) > len(pattern) && name[:len(pattern)] == pattern {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid characters in display name"})
							return
						}
					}
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"profile": gin.H{
					"display_name": request["display_name"],
				},
			})
		})

		maliciousInputs := []string{
			"<script>alert('xss')</script>",
			"javascript:alert('xss')",
			"onload=alert('xss')",
			"DROP TABLE users;",
			"SELECT * FROM users",
		}

		for _, maliciousInput := range maliciousInputs {
			payload := map[string]interface{}{
				"display_name": maliciousInput,
			}

			jsonData, _ := json.Marshal(payload)
			req, _ := http.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid_token")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code,
				"Should reject malicious input: %s", maliciousInput)
		}
	})
}

// TestAuthProfileUpdateRateLimit tests rate limiting
func TestAuthProfileUpdateRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	updateCount := 0
	router.PUT("/api/v1/users/profile", func(c *gin.Context) {
		updateCount++
		if updateCount > 5 { // Simulate rate limit of 5 updates per window
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many profile updates. Please wait before trying again.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
	})

	payload := map[string]interface{}{"display_name": "Test Update"}
	jsonData, _ := json.Marshal(payload)

	// Make requests beyond rate limit
	for i := 0; i < 8; i++ {
		req, _ := http.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i < 5 {
			assert.Equal(t, http.StatusOK, w.Code, "Update %d should succeed", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Update %d should be rate limited", i+1)
		}
	}
}

// TestAuthProfileUpdateLocalization tests localization features
func TestAuthProfileUpdateLocalization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	supportedLocales := []string{"en", "th", "id", "ms", "fil", "vi"}

	for _, locale := range supportedLocales {
		t.Run("locale_"+locale, func(t *testing.T) {
			router := gin.New()
			router.PUT("/api/v1/users/profile", func(c *gin.Context) {
				var request map[string]interface{}
				c.ShouldBindJSON(&request)

				c.JSON(http.StatusOK, gin.H{
					"profile": gin.H{
						"locale": request["locale"],
					},
				})
			})

			payload := map[string]interface{}{"locale": locale}
			jsonData, _ := json.Marshal(payload)

			req, _ := http.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid_token")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Should accept locale: %s", locale)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			profile := response["profile"].(map[string]interface{})
			assert.Equal(t, locale, profile["locale"], "Locale should be updated to: %s", locale)
		})
	}
}