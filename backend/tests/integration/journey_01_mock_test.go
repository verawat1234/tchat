// Journey 1: Mock-based Registration API Tests
// Tests registration endpoints using mock HTTP handlers (similar to auth_flow_test.go pattern)

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Journey01MockSuite struct {
	suite.Suite
	router *gin.Engine
	users  map[string]interface{}
}

type MockRegistrationRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Country   string `json:"country"`
	Language  string `json:"language"`
}

type MockRegistrationResponse struct {
	UserID     string `json:"userId"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	NextStep   string `json:"nextStep"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Country    string `json:"country"`
	Language   string `json:"language"`
}

type MockPhoneRegistrationRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Country     string `json:"country"`
	Language    string `json:"language"`
}

type MockPhoneRegistrationResponse struct {
	UserID      string `json:"userId"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	NextStep    string `json:"nextStep"`
	PhoneNumber string `json:"phoneNumber"`
	Country     string `json:"country"`
	Language    string `json:"language"`
}

type MockWelcomeNotificationResponse struct {
	NotificationID string `json:"notificationId"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	UserID        string `json:"userId"`
}

func (suite *Journey01MockSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.users = make(map[string]interface{})
	suite.setupRegistrationEndpoints()
}

func (suite *Journey01MockSuite) setupRegistrationEndpoints() {
	v1 := suite.router.Group("/api/v1")

	// Email registration endpoint
	v1.POST("/auth/register", func(c *gin.Context) {
		var req MockRegistrationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate required fields
		if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// Validate email format
		if !strings.Contains(req.Email, "@") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		// Check for duplicate email
		if _, exists := suite.users[req.Email]; exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}

		// Create user
		userID := "user_" + strings.ReplaceAll(req.Email, "@", "_")
		userID = strings.ReplaceAll(userID, ".", "_")

		user := MockRegistrationResponse{
			UserID:    userID,
			Status:    "pending_verification",
			Message:   "Registration successful. Please check your email to verify your account.",
			NextStep:  "email_verification",
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Country:   req.Country,
			Language:  req.Language,
		}

		suite.users[req.Email] = user

		c.JSON(http.StatusCreated, user)
	})

	// Phone registration endpoint
	v1.POST("/auth/register-phone", func(c *gin.Context) {
		var req MockPhoneRegistrationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate required fields
		if req.PhoneNumber == "" || req.Country == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// Validate phone number format (basic validation)
		if !strings.HasPrefix(req.PhoneNumber, "+") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must start with country code"})
			return
		}

		// Southeast Asian phone validation
		validCountries := map[string][]int{
			"+66": {11, 12}, // Thailand
			"+65": {8, 8},   // Singapore
			"+62": {11, 14}, // Indonesia
			"+60": {9, 11},  // Malaysia
			"+63": {10, 10}, // Philippines
			"+84": {9, 12},  // Vietnam
		}

		var isValid bool
		for prefix, lengths := range validCountries {
			if strings.HasPrefix(req.PhoneNumber, prefix) {
				phoneLen := len(req.PhoneNumber)
				if phoneLen >= lengths[0] && phoneLen <= lengths[1] {
					isValid = true
					break
				}
			}
		}

		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number for Southeast Asian countries"})
			return
		}

		// Check for duplicate phone
		if _, exists := suite.users[req.PhoneNumber]; exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
			return
		}

		// Create user
		userID := "user_phone_" + strings.ReplaceAll(req.PhoneNumber, "+", "")

		user := MockPhoneRegistrationResponse{
			UserID:      userID,
			Status:      "pending_verification",
			Message:     "Registration successful. Please check your SMS for verification code.",
			NextStep:    "sms_verification",
			PhoneNumber: req.PhoneNumber,
			Country:     req.Country,
			Language:    req.Language,
		}

		suite.users[req.PhoneNumber] = user

		c.JSON(http.StatusCreated, user)
	})

	// Email verification endpoint (simplified)
	v1.POST("/auth/verify-email", func(c *gin.Context) {
		var req map[string]string
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		email := req["email"]
		code := req["verificationCode"]

		if email == "" || code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing email or verification code"})
			return
		}

		// Mock verification (always accept "123456")
		if code != "123456" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "verified",
			"message": "Email verified successfully",
			"email":   email,
		})
	})

	// Welcome notification endpoint
	v1.POST("/notifications/welcome", func(c *gin.Context) {
		var req map[string]string
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userID := req["userId"]
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId"})
			return
		}

		response := MockWelcomeNotificationResponse{
			NotificationID: "welcome_" + userID,
			Type:          "welcome",
			Status:        "sent",
			Message:       "Welcome notification sent successfully",
			UserID:        userID,
		}

		c.JSON(http.StatusOK, response)
	})
}

func (suite *Journey01MockSuite) TestEmailRegistrationFlow() {
	// Test successful email registration
	registrationData := MockRegistrationRequest{
		Email:     "testuser@example.com",
		Password:  "SecurePass123!",
		FirstName: "John",
		LastName:  "Doe",
		Country:   "TH",
		Language:  "en",
	}

	jsonData, _ := json.Marshal(registrationData)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response MockRegistrationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "pending_verification", response.Status)
	assert.Equal(suite.T(), "email_verification", response.NextStep)
	assert.Contains(suite.T(), response.UserID, "user_testuser")
	assert.Equal(suite.T(), "John", response.FirstName)
	assert.Equal(suite.T(), "Doe", response.LastName)
}

func (suite *Journey01MockSuite) TestPhoneRegistrationFlow() {
	// Test successful phone registration for Southeast Asian countries
	testCases := []struct {
		name        string
		phoneNumber string
		country     string
		shouldPass  bool
	}{
		{"Thailand valid", "+66812345678", "TH", true},
		{"Singapore valid", "+6598765432", "SG", true},
		{"Indonesia valid", "+6281234567890", "ID", true},
		{"Malaysia valid", "+60123456789", "MY", true},
		{"Philippines valid", "+639123456789", "PH", true},
		{"Vietnam valid", "+84987654321", "VN", true},
		{"Invalid country code", "+1234567890", "US", false},
		{"Invalid Thailand phone", "+6681234567", "TH", false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			registrationData := MockPhoneRegistrationRequest{
				PhoneNumber: tc.phoneNumber,
				Country:     tc.country,
				Language:    "en",
			}

			jsonData, _ := json.Marshal(registrationData)
			req := httptest.NewRequest("POST", "/api/v1/auth/register-phone", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			if tc.shouldPass {
				assert.Equal(suite.T(), http.StatusCreated, w.Code)

				var response MockPhoneRegistrationResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), "pending_verification", response.Status)
				assert.Equal(suite.T(), "sms_verification", response.NextStep)
				assert.Equal(suite.T(), tc.phoneNumber, response.PhoneNumber)
				assert.Equal(suite.T(), tc.country, response.Country)
			} else {
				assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
			}
		})
	}
}

func (suite *Journey01MockSuite) TestRegistrationErrorHandling() {
	// Test missing required fields
	invalidData := MockRegistrationRequest{
		Email: "test@example.com",
		// Missing required fields
	}

	jsonData, _ := json.Marshal(invalidData)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Test invalid email format
	invalidEmailData := MockRegistrationRequest{
		Email:     "invalid-email",
		Password:  "password",
		FirstName: "John",
		LastName:  "Doe",
		Country:   "TH",
		Language:  "en",
	}

	jsonData, _ = json.Marshal(invalidEmailData)
	req = httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *Journey01MockSuite) TestDuplicateRegistration() {
	// Register a user first
	registrationData := MockRegistrationRequest{
		Email:     "duplicate@example.com",
		Password:  "SecurePass123!",
		FirstName: "Jane",
		LastName:  "Smith",
		Country:   "SG",
		Language:  "en",
	}

	// First registration - should succeed
	jsonData, _ := json.Marshal(registrationData)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Second registration with same email - should fail
	req = httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusConflict, w.Code)
}

func (suite *Journey01MockSuite) TestEmailVerification() {
	// Test email verification with correct code
	verificationData := map[string]string{
		"email":            "testuser@example.com",
		"verificationCode": "123456",
	}

	jsonData, _ := json.Marshal(verificationData)
	req := httptest.NewRequest("POST", "/api/v1/auth/verify-email", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "verified", response["status"])

	// Test with incorrect code
	verificationData["verificationCode"] = "wrong"
	jsonData, _ = json.Marshal(verificationData)
	req = httptest.NewRequest("POST", "/api/v1/auth/verify-email", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *Journey01MockSuite) TestWelcomeNotification() {
	// Test welcome notification
	notificationData := map[string]string{
		"userId": "user_testuser_example_com",
	}

	jsonData, _ := json.Marshal(notificationData)
	req := httptest.NewRequest("POST", "/api/v1/notifications/welcome", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response MockWelcomeNotificationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "welcome", response.Type)
	assert.Equal(suite.T(), "sent", response.Status)
	assert.Contains(suite.T(), response.NotificationID, "welcome_")
}

func (suite *Journey01MockSuite) TestRegionalLocalization() {
	// Test registration with different Southeast Asian countries
	countries := []struct {
		countryCode string
		language    string
		phonePrefix string
		phoneNumber string
	}{
		{"TH", "th", "+66", "+66812345678"},
		{"SG", "en", "+65", "+6598765432"},
		{"ID", "id", "+62", "+6281234567890"},
		{"MY", "en", "+60", "+60123456789"},
		{"PH", "en", "+63", "+639123456789"},
		{"VN", "vi", "+84", "+84987654321"},
	}

	for _, country := range countries {
		suite.Run("Country_"+country.countryCode, func() {
			registrationData := MockPhoneRegistrationRequest{
				PhoneNumber: country.phoneNumber,
				Country:     country.countryCode,
				Language:    country.language,
			}

			jsonData, _ := json.Marshal(registrationData)
			req := httptest.NewRequest("POST", "/api/v1/auth/register-phone", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusCreated, w.Code)

			var response MockPhoneRegistrationResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), country.countryCode, response.Country)
			assert.Equal(suite.T(), country.language, response.Language)
			assert.Equal(suite.T(), country.phoneNumber, response.PhoneNumber)
		})
	}
}

// Test runner
func TestJourney01MockSuite(t *testing.T) {
	suite.Run(t, new(Journey01MockSuite))
}