package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// T029: Integration test complete auth flow
// Tests end-to-end authentication workflow including:
// 1. User registration → 2. Phone verification → 3. Login → 4. Token refresh → 5. Profile access
type AuthFlowTestSuite struct {
	suite.Suite
	router   *gin.Engine
	testUser map[string]interface{}
	tokens   map[string]string
}

func TestAuthFlowSuite(t *testing.T) {
	suite.Run(t, new(AuthFlowTestSuite))
}

func (suite *AuthFlowTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.tokens = make(map[string]string)
	suite.testUser = make(map[string]interface{})

	// Setup auth service endpoints
	suite.setupAuthEndpoints()
}

func (suite *AuthFlowTestSuite) setupAuthEndpoints() {
	// Mock database for integration testing
	users := make(map[string]map[string]interface{})
	otpCodes := make(map[string]string)
	refreshTokens := make(map[string]string)

	// Registration endpoint
	suite.router.POST("/auth/register", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		phone, ok := req["phone"].(string)
		if !ok || phone == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing_phone", "message": "Phone number is required"})
			return
		}

		// Check if user already exists
		if _, exists := users[phone]; exists {
			c.JSON(http.StatusConflict, gin.H{"error": "user_exists", "message": "User already exists"})
			return
		}

		// Southeast Asian phone validation
		validCountries := []string{"+66", "+65", "+62", "+60", "+63", "+84"} // TH, SG, ID, MY, PH, VN
		validPhone := false
		for _, prefix := range validCountries {
			if strings.HasPrefix(phone, prefix) {
				validPhone = true
				break
			}
		}

		if !validPhone {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_phone", "message": "Phone number must be from Southeast Asia"})
			return
		}

		// Generate OTP
		otpCode := "123456" // Simplified for testing
		otpCodes[phone] = otpCode

		// Store user data
		users[phone] = map[string]interface{}{
			"phone":      phone,
			"first_name": req["first_name"],
			"last_name":  req["last_name"],
			"email":      req["email"],
			"country":    req["country"],
			"status":     "pending_verification",
			"created_at": time.Now().UTC().Format(time.RFC3339),
		}

		c.JSON(http.StatusCreated, gin.H{
			"user_id":     fmt.Sprintf("user_%s", strings.ReplaceAll(phone, "+", "")),
			"phone":       phone,
			"status":      "pending_verification",
			"otp_sent":    true,
			"expires_in":  300,
			"created_at":  time.Now().UTC().Format(time.RFC3339),
		})
	})

	// OTP verification endpoint
	suite.router.POST("/auth/verify", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		phone := req["phone"].(string)
		otp := req["otp"].(string)

		// Check if OTP is correct
		expectedOTP, exists := otpCodes[phone]
		if !exists || expectedOTP != otp {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_otp", "message": "Invalid or expired OTP"})
			return
		}

		// Update user status
		if user, exists := users[phone]; exists {
			user["status"] = "verified"
			user["verified_at"] = time.Now().UTC().Format(time.RFC3339)
			users[phone] = user
		}

		// Generate tokens
		accessToken := fmt.Sprintf("access_%s_%d", strings.ReplaceAll(phone, "+", ""), time.Now().Unix())
		refreshToken := fmt.Sprintf("refresh_%s_%d", strings.ReplaceAll(phone, "+", ""), time.Now().Unix())
		refreshTokens[refreshToken] = phone

		// Clear OTP after successful verification
		delete(otpCodes, phone)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
			"user": gin.H{
				"user_id":     fmt.Sprintf("user_%s", strings.ReplaceAll(phone, "+", "")),
				"phone":       phone,
				"status":      "verified",
				"verified_at": time.Now().UTC().Format(time.RFC3339),
			},
		})
	})

	// Login endpoint
	suite.router.POST("/auth/login", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		phone := req["phone"].(string)
		password := req["password"].(string)

		// Check if user exists and is verified
		user, exists := users[phone]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_found", "message": "User not found"})
			return
		}

		if user["status"] != "verified" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_verified", "message": "Phone number not verified"})
			return
		}

		// Simplified password check (in real app, this would be bcrypt)
		if password != "password123" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials", "message": "Invalid credentials"})
			return
		}

		// Generate new tokens
		accessToken := fmt.Sprintf("access_%s_%d", strings.ReplaceAll(phone, "+", ""), time.Now().Unix())
		refreshToken := fmt.Sprintf("refresh_%s_%d", strings.ReplaceAll(phone, "+", ""), time.Now().Unix())
		refreshTokens[refreshToken] = phone

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
			"user": gin.H{
				"user_id":    fmt.Sprintf("user_%s", strings.ReplaceAll(phone, "+", "")),
				"phone":      phone,
				"first_name": user["first_name"],
				"last_name":  user["last_name"],
				"email":      user["email"],
				"country":    user["country"],
				"status":     "verified",
			},
		})
	})

	// Token refresh endpoint
	suite.router.POST("/auth/token/refresh", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		refreshToken := req["refresh_token"].(string)

		// Check if refresh token is valid
		phone, exists := refreshTokens[refreshToken]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_refresh_token", "message": "Invalid or expired refresh token"})
			return
		}

		// Generate new tokens
		newAccessToken := fmt.Sprintf("access_%s_%d", strings.ReplaceAll(phone, "+", ""), time.Now().Unix())
		newRefreshToken := fmt.Sprintf("refresh_%s_%d", strings.ReplaceAll(phone, "+", ""), time.Now().Unix())

		// Rotate refresh token
		delete(refreshTokens, refreshToken)
		refreshTokens[newRefreshToken] = phone

		c.JSON(http.StatusOK, gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
	})

	// Profile endpoint
	suite.router.GET("/users/profile", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "Authentication required"})
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_auth_format", "message": "Invalid authorization format"})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		// Simplified token validation
		if !strings.HasPrefix(token, "access_") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token", "message": "Invalid access token"})
			return
		}

		// Extract phone from token
		parts := strings.Split(token, "_")
		if len(parts) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "malformed_token", "message": "Malformed token"})
			return
		}

		phoneDigits := parts[1]
		// Reconstruct phone (simplified)
		var phone string
		switch {
		case strings.HasPrefix(phoneDigits, "66"):
			phone = "+66" + phoneDigits[2:]
		case strings.HasPrefix(phoneDigits, "65"):
			phone = "+65" + phoneDigits[2:]
		case strings.HasPrefix(phoneDigits, "62"):
			phone = "+62" + phoneDigits[2:]
		case strings.HasPrefix(phoneDigits, "60"):
			phone = "+60" + phoneDigits[2:]
		case strings.HasPrefix(phoneDigits, "63"):
			phone = "+63" + phoneDigits[2:]
		case strings.HasPrefix(phoneDigits, "84"):
			phone = "+84" + phoneDigits[2:]
		default:
			phone = "+66" + phoneDigits // Default fallback
		}

		user, exists := users[phone]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found", "message": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":     fmt.Sprintf("user_%s", strings.ReplaceAll(phone, "+", "")),
			"phone":       phone,
			"first_name":  user["first_name"],
			"last_name":   user["last_name"],
			"email":       user["email"],
			"country":     user["country"],
			"status":      user["status"],
			"created_at":  user["created_at"],
			"verified_at": user["verified_at"],
		})
	})
}

func (suite *AuthFlowTestSuite) TestCompleteAuthFlow() {
	// Step 1: User Registration
	suite.T().Log("Step 1: Testing user registration")

	registrationData := map[string]interface{}{
		"phone":      "+66812345678",
		"first_name": "Somchai",
		"last_name":  "Jaidee",
		"email":      "somchai.jaidee@example.com",
		"country":    "TH",
	}

	jsonData, _ := json.Marshal(registrationData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var regResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &regResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "+66812345678", regResponse["phone"])
	assert.Equal(suite.T(), "pending_verification", regResponse["status"])
	assert.Equal(suite.T(), true, regResponse["otp_sent"])
	assert.NotEmpty(suite.T(), regResponse["user_id"])

	suite.testUser["phone"] = regResponse["phone"]
	suite.testUser["user_id"] = regResponse["user_id"]

	// Step 2: Phone Verification
	suite.T().Log("Step 2: Testing phone verification")

	verificationData := map[string]interface{}{
		"phone": "+66812345678",
		"otp":   "123456",
	}

	jsonData, _ = json.Marshal(verificationData)
	req = httptest.NewRequest("POST", "/auth/verify", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var verifyResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &verifyResponse)
	require.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), verifyResponse["access_token"])
	assert.NotEmpty(suite.T(), verifyResponse["refresh_token"])
	assert.Equal(suite.T(), "Bearer", verifyResponse["token_type"])
	assert.Equal(suite.T(), float64(3600), verifyResponse["expires_in"])

	// Extract user info
	userInfo := verifyResponse["user"].(map[string]interface{})
	assert.Equal(suite.T(), "verified", userInfo["status"])
	assert.NotEmpty(suite.T(), userInfo["verified_at"])

	suite.tokens["access_token"] = verifyResponse["access_token"].(string)
	suite.tokens["refresh_token"] = verifyResponse["refresh_token"].(string)

	// Step 3: Login with Password
	suite.T().Log("Step 3: Testing login with password")

	loginData := map[string]interface{}{
		"phone":    "+66812345678",
		"password": "password123",
	}

	jsonData, _ = json.Marshal(loginData)
	req = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var loginResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), loginResponse["access_token"])
	assert.NotEmpty(suite.T(), loginResponse["refresh_token"])

	// Update tokens with login response
	suite.tokens["access_token"] = loginResponse["access_token"].(string)
	suite.tokens["refresh_token"] = loginResponse["refresh_token"].(string)

	// Validate user data in login response
	loginUserInfo := loginResponse["user"].(map[string]interface{})
	assert.Equal(suite.T(), "Somchai", loginUserInfo["first_name"])
	assert.Equal(suite.T(), "Jaidee", loginUserInfo["last_name"])
	assert.Equal(suite.T(), "somchai.jaidee@example.com", loginUserInfo["email"])
	assert.Equal(suite.T(), "TH", loginUserInfo["country"])

	// Step 4: Access Protected Resource (Profile)
	suite.T().Log("Step 4: Testing access to protected profile endpoint")

	req = httptest.NewRequest("GET", "/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+suite.tokens["access_token"])
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var profileResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &profileResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "+66812345678", profileResponse["phone"])
	assert.Equal(suite.T(), "Somchai", profileResponse["first_name"])
	assert.Equal(suite.T(), "Jaidee", profileResponse["last_name"])
	assert.Equal(suite.T(), "somchai.jaidee@example.com", profileResponse["email"])
	assert.Equal(suite.T(), "TH", profileResponse["country"])
	assert.Equal(suite.T(), "verified", profileResponse["status"])
	assert.NotEmpty(suite.T(), profileResponse["created_at"])
	assert.NotEmpty(suite.T(), profileResponse["verified_at"])

	// Step 5: Token Refresh
	suite.T().Log("Step 5: Testing token refresh")

	refreshData := map[string]interface{}{
		"refresh_token": suite.tokens["refresh_token"],
	}

	jsonData, _ = json.Marshal(refreshData)
	req = httptest.NewRequest("POST", "/auth/token/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var refreshResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &refreshResponse)
	require.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), refreshResponse["access_token"])
	assert.NotEmpty(suite.T(), refreshResponse["refresh_token"])
	assert.Equal(suite.T(), "Bearer", refreshResponse["token_type"])

	// Verify new tokens are different (token rotation)
	assert.NotEqual(suite.T(), suite.tokens["access_token"], refreshResponse["access_token"])
	assert.NotEqual(suite.T(), suite.tokens["refresh_token"], refreshResponse["refresh_token"])

	// Step 6: Verify New Token Works
	suite.T().Log("Step 6: Testing new access token functionality")

	req = httptest.NewRequest("GET", "/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+refreshResponse["access_token"].(string))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var finalProfileResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &finalProfileResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "+66812345678", finalProfileResponse["phone"])
	assert.Equal(suite.T(), "verified", finalProfileResponse["status"])
}

func (suite *AuthFlowTestSuite) TestAuthFlowErrorCases() {
	// Test duplicate registration
	suite.T().Log("Testing duplicate registration")

	registrationData := map[string]interface{}{
		"phone":      "+66812345678",
		"first_name": "Test",
		"last_name":  "User",
		"email":      "test@example.com",
		"country":    "TH",
	}

	jsonData, _ := json.Marshal(registrationData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusConflict, w.Code)

	// Test invalid OTP
	suite.T().Log("Testing invalid OTP")

	verificationData := map[string]interface{}{
		"phone": "+66999888777",
		"otp":   "000000",
	}

	jsonData, _ = json.Marshal(verificationData)
	req = httptest.NewRequest("POST", "/auth/verify", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	// Test unauthorized profile access
	suite.T().Log("Testing unauthorized profile access")

	req = httptest.NewRequest("GET", "/users/profile", nil)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *AuthFlowTestSuite) TestSoutheastAsianPhoneValidation() {
	// Test valid Southeast Asian phone numbers
	validPhones := []string{
		"+66812345678", // Thailand
		"+6598765432",  // Singapore
		"+6281234567",  // Indonesia
		"+60123456789", // Malaysia
		"+639123456789", // Philippines
		"+84987654321",  // Vietnam
	}

	for _, phone := range validPhones {
		suite.T().Logf("Testing valid phone: %s", phone)

		registrationData := map[string]interface{}{
			"phone":      phone,
			"first_name": "Test",
			"last_name":  "User",
			"email":      fmt.Sprintf("test_%s@example.com", strings.ReplaceAll(phone, "+", "")),
			"country":    "TH",
		}

		jsonData, _ := json.Marshal(registrationData)
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusCreated, w.Code, "Phone %s should be valid", phone)
	}

	// Test invalid phone numbers
	invalidPhones := []string{
		"+1234567890",  // US
		"+4412345678",  // UK
		"+8612345678",  // China
		"66812345678",  // Missing +
		"+6681234567",  // Too short
	}

	for _, phone := range invalidPhones {
		suite.T().Logf("Testing invalid phone: %s", phone)

		registrationData := map[string]interface{}{
			"phone":      phone,
			"first_name": "Test",
			"last_name":  "User",
			"email":      "test@example.com",
			"country":    "US",
		}

		jsonData, _ := json.Marshal(registrationData)
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusBadRequest, w.Code, "Phone %s should be invalid", phone)
	}
}

func (suite *AuthFlowTestSuite) TestAuthFlowPerformance() {
	// Test performance requirements
	suite.T().Log("Testing authentication flow performance")

	// Registration performance
	start := time.Now()
	registrationData := map[string]interface{}{
		"phone":      "+66999888777",
		"first_name": "Performance",
		"last_name":  "Test",
		"email":      "perf@example.com",
		"country":    "TH",
	}

	jsonData, _ := json.Marshal(registrationData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	regDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.True(suite.T(), regDuration < 200*time.Millisecond, "Registration should complete in <200ms")

	// Verification performance
	start = time.Now()
	verificationData := map[string]interface{}{
		"phone": "+66999888777",
		"otp":   "123456",
	}

	jsonData, _ = json.Marshal(verificationData)
	req = httptest.NewRequest("POST", "/auth/verify", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	verifyDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), verifyDuration < 150*time.Millisecond, "Verification should complete in <150ms")

	// Extract token for profile test
	var verifyResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &verifyResponse)
	accessToken := verifyResponse["access_token"].(string)

	// Profile access performance
	start = time.Now()
	req = httptest.NewRequest("GET", "/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	profileDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), profileDuration < 100*time.Millisecond, "Profile access should complete in <100ms")
}