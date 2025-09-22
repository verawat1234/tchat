package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tchat-backend/tests/testutil"
)

// TestOTPSendContract tests the OTP send endpoint contract
func TestOTPSendContract(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	defer helper.Cleanup()

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock OTP handler (this would be the actual handler in implementation)
	router.POST("/auth/otp/send", func(c *gin.Context) {
		var req struct {
			PhoneNumber string `json:"phone_number" validate:"required"`
			CountryCode string `json:"country_code" validate:"required"`
			Purpose     string `json:"purpose" validate:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate phone number format
		if len(req.PhoneNumber) < 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
			return
		}

		// Validate country code
		validCountries := []string{"+66", "+65", "+62", "+60", "+63", "+84"}
		isValidCountry := false
		for _, code := range validCountries {
			if req.CountryCode == code {
				isValidCountry = true
				break
			}
		}
		if !isValidCountry {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported country code"})
			return
		}

		// Validate purpose
		validPurposes := []string{"registration", "login", "phone_change", "password_reset"}
		isValidPurpose := false
		for _, purpose := range validPurposes {
			if req.Purpose == purpose {
				isValidPurpose = true
				break
			}
		}
		if !isValidPurpose {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purpose"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"message":    "OTP sent successfully",
			"expires_in": 300,
			"retry_after": 60,
		})
	})

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedFields []string
		errorFields    []string
	}{
		{
			name: "Valid OTP send request for Thailand",
			requestBody: map[string]string{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"success", "message", "expires_in", "retry_after"},
		},
		{
			name: "Valid OTP send request for Singapore",
			requestBody: map[string]string{
				"phone_number": "+6591234567",
				"country_code": "+65",
				"purpose":      "login",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"success", "message", "expires_in", "retry_after"},
		},
		{
			name: "Valid OTP send request for Indonesia",
			requestBody: map[string]string{
				"phone_number": "+628123456789",
				"country_code": "+62",
				"purpose":      "phone_change",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"success", "message", "expires_in", "retry_after"},
		},
		{
			name: "Invalid phone number - too short",
			requestBody: map[string]string{
				"phone_number": "+66123",
				"country_code": "+66",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Invalid country code",
			requestBody: map[string]string{
				"phone_number": "+1234567890",
				"country_code": "+1",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Invalid purpose",
			requestBody: map[string]string{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"purpose":      "invalid_purpose",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing phone number",
			requestBody: map[string]string{
				"country_code": "+66",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing country code",
			requestBody: map[string]string{
				"phone_number": "+66812345678",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing purpose",
			requestBody: map[string]string{
				"phone_number": "+66812345678",
				"country_code": "+66",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name:           "Empty request body",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonData, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/auth/otp/send", bytes.NewReader(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Assert expected fields
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Assert error fields
			for _, field := range tt.errorFields {
				assert.Contains(t, response, field, "Error response should contain field: %s", field)
			}

			// Specific assertions for success responses
			if tt.expectedStatus == http.StatusOK {
				assert.True(t, response["success"].(bool))
				assert.Equal(t, float64(300), response["expires_in"].(float64))
				assert.Equal(t, float64(60), response["retry_after"].(float64))
			}
		})
	}
}

// TestOTPVerifyContract tests the OTP verify endpoint contract
func TestOTPVerifyContract(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	defer helper.Cleanup()

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock OTP verify handler
	router.POST("/auth/otp/verify", func(c *gin.Context) {
		var req struct {
			PhoneNumber string `json:"phone_number" validate:"required"`
			CountryCode string `json:"country_code" validate:"required"`
			Code        string `json:"code" validate:"required"`
			Purpose     string `json:"purpose" validate:"required"`
			DeviceInfo  struct {
				DeviceID   string `json:"device_id"`
				DeviceName string `json:"device_name"`
				DeviceType string `json:"device_type"`
				AppVersion string `json:"app_version"`
			} `json:"device_info,omitempty"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate required fields
		if req.PhoneNumber == "" || req.CountryCode == "" || req.Code == "" || req.Purpose == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// Validate OTP code format (6 digits)
		if len(req.Code) != 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP code format"})
			return
		}

		// Mock OTP verification (in real implementation, this would check against database)
		if req.Code == "123456" {
			// Valid OTP - return success with tokens
			c.JSON(http.StatusOK, gin.H{
				"success":       true,
				"access_token":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				"refresh_token": "refresh_token_example",
				"expires_in":    900, // 15 minutes
				"token_type":    "Bearer",
				"user": gin.H{
					"id":           "user_123",
					"phone_number": req.PhoneNumber,
					"country_code": req.CountryCode,
					"created_at":   time.Now().UTC().Format(time.RFC3339),
				},
			})
		} else if req.Code == "000000" {
			// Expired OTP
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP code expired"})
		} else {
			// Invalid OTP
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP code"})
		}
	})

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedFields []string
		errorFields    []string
	}{
		{
			name: "Valid OTP verification",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"code":         "123456",
				"purpose":      "registration",
				"device_info": map[string]string{
					"device_id":   "device_123",
					"device_name": "iPhone 14",
					"device_type": "mobile",
					"app_version": "1.0.0",
				},
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"success", "access_token", "refresh_token", "expires_in", "token_type", "user"},
		},
		{
			name: "Valid OTP verification without device info",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"code":         "123456",
				"purpose":      "login",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"success", "access_token", "refresh_token", "expires_in", "token_type", "user"},
		},
		{
			name: "Invalid OTP code",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"code":         "999999",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Expired OTP code",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"code":         "000000",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Invalid OTP code format - too short",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"code":         "123",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Invalid OTP code format - too long",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"code":         "1234567",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing phone number",
			requestBody: map[string]interface{}{
				"country_code": "+66",
				"code":         "123456",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing country code",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"code":         "123456",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing OTP code",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing purpose",
			requestBody: map[string]interface{}{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"code":         "123456",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonData, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/auth/otp/verify", bytes.NewReader(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Assert expected fields
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Assert error fields
			for _, field := range tt.errorFields {
				assert.Contains(t, response, field, "Error response should contain field: %s", field)
			}

			// Specific assertions for success responses
			if tt.expectedStatus == http.StatusOK {
				assert.True(t, response["success"].(bool))
				assert.NotEmpty(t, response["access_token"])
				assert.NotEmpty(t, response["refresh_token"])
				assert.Equal(t, float64(900), response["expires_in"].(float64))
				assert.Equal(t, "Bearer", response["token_type"])

				// Validate user object
				user := response["user"].(map[string]interface{})
				assert.NotEmpty(t, user["id"])
				assert.NotEmpty(t, user["phone_number"])
				assert.NotEmpty(t, user["country_code"])
				assert.NotEmpty(t, user["created_at"])
			}
		})
	}
}

// TestOTPResendContract tests the OTP resend endpoint contract
func TestOTPResendContract(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	defer helper.Cleanup()

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock OTP resend handler
	router.POST("/auth/otp/resend", func(c *gin.Context) {
		var req struct {
			PhoneNumber string `json:"phone_number" validate:"required"`
			CountryCode string `json:"country_code" validate:"required"`
			Purpose     string `json:"purpose" validate:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate required fields
		if req.PhoneNumber == "" || req.CountryCode == "" || req.Purpose == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// Mock rate limiting check
		if req.PhoneNumber == "+66812345678" && req.Purpose == "registration" {
			// Simulate rate limit exceeded
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": 60,
			})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"message":     "OTP resent successfully",
			"expires_in":  300,
			"retry_after": 60,
		})
	})

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedFields []string
		errorFields    []string
	}{
		{
			name: "Valid OTP resend request",
			requestBody: map[string]string{
				"phone_number": "+66887654321",
				"country_code": "+66",
				"purpose":      "login",
			},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"success", "message", "expires_in", "retry_after"},
		},
		{
			name: "Rate limit exceeded",
			requestBody: map[string]string{
				"phone_number": "+66812345678",
				"country_code": "+66",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusTooManyRequests,
			errorFields:    []string{"error", "retry_after"},
		},
		{
			name: "Missing phone number",
			requestBody: map[string]string{
				"country_code": "+66",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing country code",
			requestBody: map[string]string{
				"phone_number": "+66812345678",
				"purpose":      "registration",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
		{
			name: "Missing purpose",
			requestBody: map[string]string{
				"phone_number": "+66812345678",
				"country_code": "+66",
			},
			expectedStatus: http.StatusBadRequest,
			errorFields:    []string{"error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonData, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/auth/otp/resend", bytes.NewReader(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Assert expected fields
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Assert error fields
			for _, field := range tt.errorFields {
				assert.Contains(t, response, field, "Error response should contain field: %s", field)
			}

			// Specific assertions for success responses
			if tt.expectedStatus == http.StatusOK {
				assert.True(t, response["success"].(bool))
				assert.Equal(t, float64(300), response["expires_in"].(float64))
				assert.Equal(t, float64(60), response["retry_after"].(float64))
			}

			// Specific assertions for rate limit responses
			if tt.expectedStatus == http.StatusTooManyRequests {
				assert.Equal(t, float64(60), response["retry_after"].(float64))
			}
		})
	}
}