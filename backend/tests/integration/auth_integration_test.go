package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// AuthIntegrationSuite tests the Auth service endpoints
type AuthIntegrationSuite struct {
	APIIntegrationSuite
	ports ServicePort
}

// AuthRequest represents authentication request payload
type AuthRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Country     string `json:"country"`
	Language    string `json:"language"`
}

// VerifyOTPRequest represents OTP verification payload
type VerifyOTPRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	OTP         string `json:"otp"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success      bool   `json:"success"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	UserID       string `json:"userId,omitempty"`
	Data         struct {
		User struct {
			ID          string `json:"id"`
			PhoneNumber string `json:"phoneNumber"`
			Country     string `json:"country"`
			Language    string `json:"language"`
		} `json:"user,omitempty"`
		Tokens struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		} `json:"tokens,omitempty"`
	} `json:"data,omitempty"`
}

// SetupSuite initializes the auth integration test suite
func (suite *AuthIntegrationSuite) SetupSuite() {
	suite.APIIntegrationSuite.SetupSuite()
	suite.ports = DefaultServicePorts()

	// Wait for auth service to be available
	err := suite.waitForService(suite.ports.Auth, 30*time.Second)
	if err != nil {
		suite.T().Fatalf("Auth service not available: %v", err)
	}
}

// TestAuthServiceHealth verifies auth service health endpoint
func (suite *AuthIntegrationSuite) TestAuthServiceHealth() {
	healthCheck, err := suite.checkServiceHealth(suite.ports.Auth)
	require.NoError(suite.T(), err, "Health check should succeed")

	assert.Equal(suite.T(), "healthy", healthCheck.Status)
	assert.Equal(suite.T(), "auth-service", healthCheck.Service)
	assert.NotEmpty(suite.T(), healthCheck.Timestamp)
}

// TestRegisterNewUser tests user registration flow
func (suite *AuthIntegrationSuite) TestRegisterNewUser() {
	url := fmt.Sprintf("%s:%d/api/v1/auth/register", suite.baseURL, suite.ports.Auth)

	// Generate unique phone number for test
	testPhone := fmt.Sprintf("+1555%07d", time.Now().Unix()%10000000)

	authReq := AuthRequest{
		PhoneNumber: testPhone,
		Country:     "US",
		Language:    "en",
	}

	resp, err := suite.makeRequest("POST", url, authReq, nil)
	require.NoError(suite.T(), err, "Registration request should succeed")
	defer resp.Body.Close()

	// Should return 200 for successful registration initiation
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var authResp AuthResponse
	err = suite.parseResponse(resp, &authResp)
	require.NoError(suite.T(), err, "Should parse registration response")

	assert.True(suite.T(), authResp.Success)
	assert.Equal(suite.T(), "success", authResp.Status)
	assert.Contains(suite.T(), authResp.Message, "OTP")
}

// TestRegisterDuplicateUser tests duplicate user registration
func (suite *AuthIntegrationSuite) TestRegisterDuplicateUser() {
	url := fmt.Sprintf("%s:%d/api/v1/auth/register", suite.baseURL, suite.ports.Auth)

	// Use the same phone number as test user
	authReq := AuthRequest{
		PhoneNumber: suite.testUser.PhoneNumber,
		Country:     suite.testUser.Country,
		Language:    suite.testUser.Language,
	}

	resp, err := suite.makeRequest("POST", url, authReq, nil)
	require.NoError(suite.T(), err, "Duplicate registration request should complete")
	defer resp.Body.Close()

	// Should handle duplicate registration gracefully
	assert.True(suite.T(), resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusConflict)
}

// TestRegisterInvalidPhoneNumber tests registration with invalid phone
func (suite *AuthIntegrationSuite) TestRegisterInvalidPhoneNumber() {
	url := fmt.Sprintf("%s:%d/api/v1/auth/register", suite.baseURL, suite.ports.Auth)

	authReq := AuthRequest{
		PhoneNumber: "invalid-phone",
		Country:     "US",
		Language:    "en",
	}

	resp, err := suite.makeRequest("POST", url, authReq, nil)
	require.NoError(suite.T(), err, "Invalid phone registration should complete")
	defer resp.Body.Close()

	// Should return validation error
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var authResp AuthResponse
	err = suite.parseResponse(resp, &authResp)
	require.NoError(suite.T(), err, "Should parse error response")

	assert.False(suite.T(), authResp.Success)
	assert.Contains(suite.T(), authResp.Message, "phone")
}

// TestLogin tests user login flow
func (suite *AuthIntegrationSuite) TestLogin() {
	url := fmt.Sprintf("%s:%d/api/v1/auth/login", suite.baseURL, suite.ports.Auth)

	authReq := AuthRequest{
		PhoneNumber: suite.testUser.PhoneNumber,
		Country:     suite.testUser.Country,
		Language:    suite.testUser.Language,
	}

	resp, err := suite.makeRequest("POST", url, authReq, nil)
	require.NoError(suite.T(), err, "Login request should succeed")
	defer resp.Body.Close()

	// Should return 200 for successful login initiation
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var authResp AuthResponse
	err = suite.parseResponse(resp, &authResp)
	require.NoError(suite.T(), err, "Should parse login response")

	assert.True(suite.T(), authResp.Success)
	assert.Equal(suite.T(), "success", authResp.Status)
	assert.Contains(suite.T(), authResp.Message, "OTP")
}

// TestLoginNonExistentUser tests login with non-existent user
func (suite *AuthIntegrationSuite) TestLoginNonExistentUser() {
	url := fmt.Sprintf("%s:%d/api/v1/auth/login", suite.baseURL, suite.ports.Auth)

	authReq := AuthRequest{
		PhoneNumber: "+1555999999999", // Non-existent phone
		Country:     "US",
		Language:    "en",
	}

	resp, err := suite.makeRequest("POST", url, authReq, nil)
	require.NoError(suite.T(), err, "Non-existent user login should complete")
	defer resp.Body.Close()

	// Should return 404 or 400 for non-existent user
	assert.True(suite.T(), resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest)
}

// TestVerifyOTP tests OTP verification (with mock OTP)
func (suite *AuthIntegrationSuite) TestVerifyOTP() {
	url := fmt.Sprintf("%s:%d/api/v1/auth/verify-otp", suite.baseURL, suite.ports.Auth)

	// Use test OTP (in real implementation, this would be obtained from registration/login)
	verifyReq := VerifyOTPRequest{
		PhoneNumber: suite.testUser.PhoneNumber,
		OTP:         "123456", // Test OTP
	}

	resp, err := suite.makeRequest("POST", url, verifyReq, nil)
	require.NoError(suite.T(), err, "OTP verification request should succeed")
	defer resp.Body.Close()

	// Response varies based on OTP validity
	assert.True(suite.T(), resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest)

	var authResp AuthResponse
	err = suite.parseResponse(resp, &authResp)
	require.NoError(suite.T(), err, "Should parse OTP verification response")

	// If successful, should contain tokens
	if resp.StatusCode == http.StatusOK && authResp.Success {
		assert.NotEmpty(suite.T(), authResp.Data.Tokens.AccessToken)
		assert.NotEmpty(suite.T(), authResp.Data.Tokens.RefreshToken)
		assert.NotEmpty(suite.T(), authResp.Data.User.ID)
	}
}

// TestVerifyInvalidOTP tests OTP verification with invalid OTP
func (suite *AuthIntegrationSuite) TestVerifyInvalidOTP() {
	url := fmt.Sprintf("%s:%d/api/v1/auth/verify-otp", suite.baseURL, suite.ports.Auth)

	verifyReq := VerifyOTPRequest{
		PhoneNumber: suite.testUser.PhoneNumber,
		OTP:         "000000", // Invalid OTP
	}

	resp, err := suite.makeRequest("POST", url, verifyReq, nil)
	require.NoError(suite.T(), err, "Invalid OTP verification should complete")
	defer resp.Body.Close()

	// Should return 400 for invalid OTP
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var authResp AuthResponse
	err = suite.parseResponse(resp, &authResp)
	require.NoError(suite.T(), err, "Should parse invalid OTP response")

	assert.False(suite.T(), authResp.Success)
	assert.Contains(suite.T(), authResp.Message, "OTP")
}

// TestMissingRequestBody tests endpoints with missing request body
func (suite *AuthIntegrationSuite) TestMissingRequestBody() {
	endpoints := []string{
		fmt.Sprintf("%s:%d/api/v1/auth/register", suite.baseURL, suite.ports.Auth),
		fmt.Sprintf("%s:%d/api/v1/auth/login", suite.baseURL, suite.ports.Auth),
		fmt.Sprintf("%s:%d/api/v1/auth/verify-otp", suite.baseURL, suite.ports.Auth),
	}

	for _, url := range endpoints {
		resp, err := suite.makeRequest("POST", url, nil, nil)
		require.NoError(suite.T(), err, "Empty body request should complete")
		defer resp.Body.Close()

		// Should return 400 for missing body
		assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode, "URL: %s", url)
	}
}

// TestInvalidHTTPMethods tests endpoints with invalid HTTP methods
func (suite *AuthIntegrationSuite) TestInvalidHTTPMethods() {
	endpoints := []string{
		fmt.Sprintf("%s:%d/api/v1/auth/register", suite.baseURL, suite.ports.Auth),
		fmt.Sprintf("%s:%d/api/v1/auth/login", suite.baseURL, suite.ports.Auth),
		fmt.Sprintf("%s:%d/api/v1/auth/verify-otp", suite.baseURL, suite.ports.Auth),
	}

	invalidMethods := []string{"GET", "PUT", "DELETE", "PATCH"}

	for _, url := range endpoints {
		for _, method := range invalidMethods {
			resp, err := suite.makeRequest(method, url, nil, nil)
			require.NoError(suite.T(), err, "Invalid method request should complete")
			defer resp.Body.Close()

			// Should return 405 Method Not Allowed
			assert.Equal(suite.T(), http.StatusMethodNotAllowed, resp.StatusCode,
				"URL: %s, Method: %s", url, method)
		}
	}
}

// TestAuthServiceConcurrency tests concurrent requests to auth service
func (suite *AuthIntegrationSuite) TestAuthServiceConcurrency() {
	url := fmt.Sprintf("%s:%d/api/v1/auth/register", suite.baseURL, suite.ports.Auth)

	// Create 10 concurrent registration requests
	concurrency := 10
	results := make(chan int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			testPhone := fmt.Sprintf("+1555%03d%04d", index, time.Now().Unix()%10000)
			authReq := AuthRequest{
				PhoneNumber: testPhone,
				Country:     "US",
				Language:    "en",
			}

			resp, err := suite.makeRequest("POST", url, authReq, nil)
			if err != nil {
				results <- 0
				return
			}
			defer resp.Body.Close()

			results <- resp.StatusCode
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		statusCode := <-results
		if statusCode == http.StatusOK {
			successCount++
		}
	}

	// At least 70% of concurrent requests should succeed
	assert.GreaterOrEqual(suite.T(), successCount, 7, "Concurrent requests should mostly succeed")
}

// RunAuthIntegrationTests runs the auth integration test suite
func RunAuthIntegrationTests(t *testing.T) {
	suite.Run(t, new(AuthIntegrationSuite))
}