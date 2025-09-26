// Journey 1: New User Registration & Onboarding API Integration Tests
// Tests all API endpoints involved in user registration and onboarding flow

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Journey01RegistrationAPISuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	testUser   *TestUser
	ctx        context.Context
}

type TestUser struct {
	Email     string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Country   string `json:"country"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
}

type RegistrationRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Country   string `json:"country"`
	Language  string `json:"language"`
}

type RegistrationResponse struct {
	UserID     string `json:"userId"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	VerifyCode string `json:"verifyCode,omitempty"`
}

type PhoneRegistrationRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Country  string `json:"country"`
	Language string `json:"language"`
}

type PhoneRegistrationResponse struct {
	VerificationID string `json:"verificationId"`
	Status         string `json:"status"`
}

type OTPVerificationRequest struct {
	VerificationID string `json:"verificationId"`
	OTP           string `json:"otp"`
}

type OTPVerificationResponse struct {
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Status       string `json:"status"`
}

type ProfileSetupRequest struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Country   string `json:"country"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	Avatar    string `json:"avatar,omitempty"`
}

type ProfileSetupResponse struct {
	ProfileID string `json:"profileId"`
	Status    string `json:"status"`
}

type WelcomeNotificationRequest struct {
	UserID   string `json:"userId"`
	Language string `json:"language"`
	Country  string `json:"country"`
}

type NotificationResponse struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Status   string `json:"status"`
	SentAt   string `json:"sentAt"`
}

func (suite *Journey01RegistrationAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081" // Auth Service Direct
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.ctx = context.Background()

	// Initialize test user for Singapore
	suite.testUser = &TestUser{
		Email:     fmt.Sprintf("li.wei.test.%d@tchat.com", time.Now().Unix()),
		FirstName: "Li Wei",
		LastName:  "Tan",
		Country:   "SG",
		Language:  "en",
		Timezone:  "Asia/Singapore",
	}
}

// Test 1.1: Email Registration API Flow
func (suite *Journey01RegistrationAPISuite) TestEmailRegistrationFlow() {
	// Step 1: POST /api/v1/auth/register
	regReq := RegistrationRequest{
		Email:     suite.testUser.Email,
		Password:  "SecurePassword123",
		FirstName: suite.testUser.FirstName,
		LastName:  suite.testUser.LastName,
		Country:   suite.testUser.Country,
		Language:  suite.testUser.Language,
	}

	regResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register", regReq, nil)
	assert.Equal(suite.T(), 201, statusCode, "Registration should return 201 Created")

	var regResult RegistrationResponse
	err := json.Unmarshal(regResp, &regResult)
	require.NoError(suite.T(), err, "Should parse registration response")

	assert.NotEmpty(suite.T(), regResult.UserID, "Should return user ID")
	assert.Equal(suite.T(), "pending_verification", regResult.Status, "Status should be pending verification")

	userID := regResult.UserID

	// Step 2: Verify user in auth service - GET /api/v1/users/{userId}
	userResp, statusCode := suite.makeAPICall("GET", fmt.Sprintf("/api/v1/users/%s", userID), nil, nil)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve user info")

	var userInfo map[string]interface{}
	err = json.Unmarshal(userResp, &userInfo)
	require.NoError(suite.T(), err, "Should parse user info")
	assert.Equal(suite.T(), suite.testUser.Email, userInfo["email"])
	assert.Equal(suite.T(), false, userInfo["verified"], "User should not be verified yet")

	// Step 3: Simulate email verification - POST /api/v1/auth/verify
	verifyReq := map[string]string{
		"userId": userID,
		"code":   regResult.VerifyCode,
	}

	verifyResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/verify", verifyReq, nil)
	assert.Equal(suite.T(), 200, statusCode, "Verification should succeed")

	var verifyResult map[string]interface{}
	err = json.Unmarshal(verifyResp, &verifyResult)
	require.NoError(suite.T(), err, "Should parse verification response")

	assert.NotEmpty(suite.T(), verifyResult["accessToken"], "Should return access token")
	assert.NotEmpty(suite.T(), verifyResult["refreshToken"], "Should return refresh token")
}

// Test 1.2: Phone Registration API Flow
func (suite *Journey01RegistrationAPISuite) TestPhoneRegistrationFlow() {
	// Step 1: POST /api/v1/auth/register-phone
	phoneReq := PhoneRegistrationRequest{
		PhoneNumber: "+66812345678", // Thailand number
		Country:     "TH",
		Language:    "th",
	}

	phoneResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register-phone", phoneReq, nil)
	assert.Equal(suite.T(), 200, statusCode, "Phone registration should succeed")

	var phoneResult PhoneRegistrationResponse
	err := json.Unmarshal(phoneResp, &phoneResult)
	require.NoError(suite.T(), err, "Should parse phone registration response")

	assert.NotEmpty(suite.T(), phoneResult.VerificationID, "Should return verification ID")
	assert.Equal(suite.T(), "otp_sent", phoneResult.Status, "Status should be otp_sent")

	// Step 2: POST /api/v1/auth/verify-otp
	otpReq := OTPVerificationRequest{
		VerificationID: phoneResult.VerificationID,
		OTP:           "123456", // Test OTP
	}

	otpResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/verify-otp", otpReq, nil)
	assert.Equal(suite.T(), 200, statusCode, "OTP verification should succeed")

	var otpResult OTPVerificationResponse
	err = json.Unmarshal(otpResp, &otpResult)
	require.NoError(suite.T(), err, "Should parse OTP verification response")

	assert.NotEmpty(suite.T(), otpResult.UserID, "Should return user ID")
	assert.NotEmpty(suite.T(), otpResult.AccessToken, "Should return access token")
	assert.Equal(suite.T(), "verified", otpResult.Status, "Status should be verified")
}

// Test 1.3: Profile Setup API Flow
func (suite *Journey01RegistrationAPISuite) TestProfileSetupFlow() {
	// First register and verify a user
	userID, accessToken := suite.createVerifiedUser()

	// Step 1: POST /api/v1/profiles
	profileReq := ProfileSetupRequest{
		UserID:    userID,
		FirstName: suite.testUser.FirstName,
		LastName:  suite.testUser.LastName,
		Country:   suite.testUser.Country,
		Language:  suite.testUser.Language,
		Timezone:  suite.testUser.Timezone,
		Avatar:    "https://example.com/avatar.jpg",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	profileResp, statusCode := suite.makeAPICall("POST", "/api/v1/profiles", profileReq, headers)
	assert.Equal(suite.T(), 201, statusCode, "Profile creation should succeed")

	var profileResult ProfileSetupResponse
	err := json.Unmarshal(profileResp, &profileResult)
	require.NoError(suite.T(), err, "Should parse profile response")

	assert.NotEmpty(suite.T(), profileResult.ProfileID, "Should return profile ID")
	assert.Equal(suite.T(), "created", profileResult.Status, "Status should be created")

	// Step 2: GET /api/v1/profiles/{profileId} - Verify profile was created
	getProfileResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/profiles/%s", profileResult.ProfileID), nil, headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve profile")

	var profileInfo map[string]interface{}
	err = json.Unmarshal(getProfileResp, &profileInfo)
	require.NoError(suite.T(), err, "Should parse profile info")

	assert.Equal(suite.T(), suite.testUser.Country, profileInfo["country"])
	assert.Equal(suite.T(), suite.testUser.Language, profileInfo["language"])
}

// Test 1.4: Welcome Notification API Flow
func (suite *Journey01RegistrationAPISuite) TestWelcomeNotificationFlow() {
	userID, accessToken := suite.createVerifiedUser()

	// Step 1: POST /api/v1/notifications/welcome
	welcomeReq := WelcomeNotificationRequest{
		UserID:   userID,
		Language: suite.testUser.Language,
		Country:  suite.testUser.Country,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	_, statusCode := suite.makeAPICall("POST", "/api/v1/notifications/welcome", welcomeReq, headers)
	assert.Equal(suite.T(), 200, statusCode, "Welcome notification should succeed")

	// Step 2: GET /api/v1/notifications/user/{userId} - Check notification was sent
	time.Sleep(2 * time.Second) // Allow notification processing

	notifResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/notifications/user/%s", userID), nil, headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve notifications")

	var notifications []NotificationResponse
	err := json.Unmarshal(notifResp, &notifications)
	require.NoError(suite.T(), err, "Should parse notifications")

	assert.Greater(suite.T(), len(notifications), 0, "Should have notifications")

	// Find welcome notification
	var welcomeNotif *NotificationResponse
	for _, notif := range notifications {
		if notif.Type == "welcome" {
			welcomeNotif = &notif
			break
		}
	}

	require.NotNil(suite.T(), welcomeNotif, "Should find welcome notification")
	assert.Contains(suite.T(), welcomeNotif.Title, "Welcome", "Title should contain Welcome")
	assert.Equal(suite.T(), "sent", welcomeNotif.Status, "Status should be sent")
}

// Test 1.5: Regional Localization API Testing
func (suite *Journey01RegistrationAPISuite) TestRegionalLocalization() {
	countries := []struct {
		code     string
		language string
		phone    string
		currency string
	}{
		{"SG", "en", "+6591234567", "SGD"},
		{"TH", "th", "+66812345678", "THB"},
		{"ID", "id", "+628123456789", "IDR"},
		{"PH", "tl", "+639171234567", "PHP"},
		{"MY", "ms", "+60123456789", "MYR"},
		{"VN", "vi", "+84901234567", "VND"},
	}

	for _, country := range countries {
		suite.T().Run(fmt.Sprintf("Country_%s", country.code), func(t *testing.T) {
			// Test country-specific registration
			phoneReq := PhoneRegistrationRequest{
				PhoneNumber: country.phone,
				Country:     country.code,
				Language:    country.language,
			}

			phoneResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register-phone", phoneReq, nil)
			assert.Equal(t, 200, statusCode, fmt.Sprintf("Phone registration should work for %s", country.code))

			var phoneResult PhoneRegistrationResponse
			err := json.Unmarshal(phoneResp, &phoneResult)
			require.NoError(t, err, "Should parse phone registration response")
			assert.Equal(t, "otp_sent", phoneResult.Status, "Should send OTP")

			// Test country-specific currency API - GET /api/v1/commerce/currencies/{country}
			currencyResp, statusCode := suite.makeAPICall("GET",
				fmt.Sprintf("/api/v1/commerce/currencies/%s", country.code), nil, nil)
			assert.Equal(t, 200, statusCode, fmt.Sprintf("Should get currency for %s", country.code))

			var currencyInfo map[string]interface{}
			err = json.Unmarshal(currencyResp, &currencyInfo)
			require.NoError(t, err, "Should parse currency info")
			assert.Equal(t, country.currency, currencyInfo["code"],
				fmt.Sprintf("Currency should be %s for %s", country.currency, country.code))
		})
	}
}

// Test 1.6: Registration Error Handling
func (suite *Journey01RegistrationAPISuite) TestRegistrationErrorHandling() {
	// Test duplicate email registration
	uniqueEmail := fmt.Sprintf("duplicate.%d@test.com", time.Now().UnixNano())
	regReq := RegistrationRequest{
		Email:     uniqueEmail,
		Password:  "SecurePassword123",
		FirstName: "Test",
		LastName:  "User",
		Country:   "SG",
		Language:  "en",
	}

	// First registration should succeed
	_, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register", regReq, nil)
	assert.Equal(suite.T(), 201, statusCode, "First registration should succeed")

	// Second registration should fail
	_, statusCode = suite.makeAPICall("POST", "/api/v1/auth/register", regReq, nil)
	assert.Equal(suite.T(), 409, statusCode, "Duplicate registration should return 409 Conflict")

	// Test invalid email format
	invalidEmailReq := regReq
	invalidEmailReq.Email = "invalid-email"

	_, statusCode = suite.makeAPICall("POST", "/api/v1/auth/register", invalidEmailReq, nil)
	assert.Equal(suite.T(), 400, statusCode, "Invalid email should return 400 Bad Request")

	// Test weak password
	weakPasswordReq := regReq
	weakPasswordReq.Email = fmt.Sprintf("weak.%d@test.com", time.Now().UnixNano())
	weakPasswordReq.Password = "123"

	_, statusCode = suite.makeAPICall("POST", "/api/v1/auth/register", weakPasswordReq, nil)
	assert.Equal(suite.T(), 400, statusCode, "Weak password should return 400 Bad Request")
}

// Helper methods
func (suite *Journey01RegistrationAPISuite) makeAPICall(method, endpoint string, body interface{}, headers map[string]string) ([]byte, int) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(suite.ctx, method, suite.baseURL+endpoint, reqBody)
	require.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	return respBody, resp.StatusCode
}

func (suite *Journey01RegistrationAPISuite) createVerifiedUser() (string, string) {
	regReq := RegistrationRequest{
		Email:     fmt.Sprintf("verified.user.%d@test.com", time.Now().UnixNano()),
		Password:  "SecurePassword123",
		FirstName: "Test",
		LastName:  "User",
		Country:   "SG",
		Language:  "en",
	}

	regResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register", regReq, nil)
	require.Equal(suite.T(), 201, statusCode)

	var regResult RegistrationResponse
	err := json.Unmarshal(regResp, &regResult)
	require.NoError(suite.T(), err)

	verifyReq := map[string]string{
		"userId": regResult.UserID,
		"code":   regResult.VerifyCode,
	}

	verifyResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/verify", verifyReq, nil)
	require.Equal(suite.T(), 200, statusCode)

	var verifyResult map[string]interface{}
	err = json.Unmarshal(verifyResp, &verifyResult)
	require.NoError(suite.T(), err)

	return regResult.UserID, verifyResult["accessToken"].(string)
}

func TestJourney01RegistrationAPISuite(t *testing.T) {
	suite.Run(t, new(Journey01RegistrationAPISuite))
}