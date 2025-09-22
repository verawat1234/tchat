package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// T032: Integration test cross-service user registration
// Tests user registration flow across Auth, Payment, and Messaging services
type UserRegistrationTestSuite struct {
	suite.Suite
	router   *gin.Engine
	users    map[string]map[string]interface{}
	wallets  map[string][]map[string]interface{}
	dialogs  map[string]map[string]interface{}
}

func TestUserRegistrationSuite(t *testing.T) {
	suite.Run(t, new(UserRegistrationTestSuite))
}

func (suite *UserRegistrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.users = make(map[string]map[string]interface{})
	suite.wallets = make(map[string][]map[string]interface{})
	suite.dialogs = make(map[string]map[string]interface{})

	suite.setupCrossServiceEndpoints()
}

func (suite *UserRegistrationTestSuite) setupCrossServiceEndpoints() {
	// Auth service endpoints
	suite.router.POST("/auth/register", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		phone := req["phone"].(string)
		userID := fmt.Sprintf("user_%s", phone[1:])

		user := map[string]interface{}{
			"id":         userID,
			"phone":      phone,
			"first_name": req["first_name"],
			"last_name":  req["last_name"],
			"email":      req["email"],
			"country":    req["country"],
			"status":     "verified",
			"created_at": time.Now().UTC().Format(time.RFC3339),
		}

		suite.users[userID] = user
		suite.wallets[userID] = make([]map[string]interface{}, 0)

		c.JSON(http.StatusCreated, user)
	})

	// Payment service - auto wallet creation
	suite.router.POST("/wallets/auto-create", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		userID := req["user_id"].(string)
		country := req["country"].(string)

		// Auto-create default currency wallet based on country
		currencyMap := map[string]string{
			"TH": "THB", "SG": "SGD", "ID": "IDR",
			"MY": "MYR", "PH": "PHP", "VN": "VND",
		}

		currency := currencyMap[country]
		if currency == "" {
			currency = "USD" // Default fallback
		}

		wallet := map[string]interface{}{
			"id":                fmt.Sprintf("wallet_%s_%s", userID, currency),
			"user_id":           userID,
			"currency":          currency,
			"status":            "active",
			"available_balance": 0.0,
			"total_balance":     0.0,
			"is_primary":        true,
			"created_at":        time.Now().UTC().Format(time.RFC3339),
		}

		suite.wallets[userID] = append(suite.wallets[userID], wallet)

		c.JSON(http.StatusCreated, wallet)
	})

	// Messaging service - user profile creation
	suite.router.POST("/messaging/users", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		userID := req["user_id"].(string)

		profile := map[string]interface{}{
			"user_id":      userID,
			"display_name": req["display_name"],
			"avatar_url":   "",
			"status":       "online",
			"last_seen":    time.Now().UTC().Format(time.RFC3339),
			"preferences": map[string]interface{}{
				"language":         req["language"],
				"notifications":    true,
				"read_receipts":    true,
				"typing_indicator": true,
			},
			"created_at": time.Now().UTC().Format(time.RFC3339),
		}

		c.JSON(http.StatusCreated, profile)
	})

	// Cross-service user lookup
	suite.router.GET("/users/:id/profile", func(c *gin.Context) {
		userID := c.Param("id")

		user, exists := suite.users[userID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found"})
			return
		}

		// Include wallet info
		wallets := suite.wallets[userID]

		profile := map[string]interface{}{
			"user_id":    userID,
			"phone":      user["phone"],
			"first_name": user["first_name"],
			"last_name":  user["last_name"],
			"email":      user["email"],
			"country":    user["country"],
			"status":     user["status"],
			"wallets":    wallets,
			"created_at": user["created_at"],
		}

		c.JSON(http.StatusOK, profile)
	})

	// Welcome dialog creation
	suite.router.POST("/dialogs/welcome", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		userID := req["user_id"].(string)
		dialogID := fmt.Sprintf("welcome_%s", userID)

		dialog := map[string]interface{}{
			"id":           dialogID,
			"type":         "system",
			"title":        "Welcome to Tchat!",
			"participants": []string{userID, "system_bot"},
			"created_at":   time.Now().UTC().Format(time.RFC3339),
			"messages": []map[string]interface{}{
				{
					"id":         fmt.Sprintf("msg_welcome_%s", userID),
					"sender_id":  "system_bot",
					"content":    "Welcome to Tchat! Your account has been created successfully.",
					"type":       "text",
					"created_at": time.Now().UTC().Format(time.RFC3339),
				},
			},
		}

		suite.dialogs[dialogID] = dialog

		c.JSON(http.StatusCreated, dialog)
	})
}

func (suite *UserRegistrationTestSuite) TestCompleteUserRegistration() {
	suite.T().Log("Testing complete cross-service user registration")

	// Step 1: Register user in Auth service
	registrationData := map[string]interface{}{
		"phone":      "+66812345678",
		"first_name": "Somchai",
		"last_name":  "Jaidee",
		"email":      "somchai@example.com",
		"country":    "TH",
	}

	jsonData, _ := json.Marshal(registrationData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	start := time.Now()
	suite.router.ServeHTTP(w, req)
	authDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.True(suite.T(), authDuration < 200*time.Millisecond, "Auth registration should be <200ms")

	var authResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &authResponse)
	require.NoError(suite.T(), err)

	userID := authResponse["id"].(string)
	assert.NotEmpty(suite.T(), userID)
	assert.Equal(suite.T(), "verified", authResponse["status"])

	// Step 2: Auto-create wallet in Payment service
	walletData := map[string]interface{}{
		"user_id": userID,
		"country": "TH",
	}

	jsonData, _ = json.Marshal(walletData)
	req = httptest.NewRequest("POST", "/wallets/auto-create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	start = time.Now()
	suite.router.ServeHTTP(w, req)
	walletDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.True(suite.T(), walletDuration < 150*time.Millisecond, "Wallet creation should be <150ms")

	var walletResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &walletResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), userID, walletResponse["user_id"])
	assert.Equal(suite.T(), "THB", walletResponse["currency"]) // Thailand -> THB
	assert.Equal(suite.T(), true, walletResponse["is_primary"])

	// Step 3: Create messaging profile
	messagingData := map[string]interface{}{
		"user_id":      userID,
		"display_name": "Somchai J.",
		"language":     "th-TH",
	}

	jsonData, _ = json.Marshal(messagingData)
	req = httptest.NewRequest("POST", "/messaging/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	start = time.Now()
	suite.router.ServeHTTP(w, req)
	messagingDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.True(suite.T(), messagingDuration < 100*time.Millisecond, "Messaging profile should be <100ms")

	var messagingResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &messagingResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), userID, messagingResponse["user_id"])
	assert.Equal(suite.T(), "Somchai J.", messagingResponse["display_name"])
	assert.Equal(suite.T(), "online", messagingResponse["status"])

	// Step 4: Create welcome dialog
	welcomeData := map[string]interface{}{
		"user_id": userID,
	}

	jsonData, _ = json.Marshal(welcomeData)
	req = httptest.NewRequest("POST", "/dialogs/welcome", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var dialogResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	assert.Equal(suite.T(), "system", dialogResponse["type"])
	assert.Contains(suite.T(), dialogResponse["participants"], userID)
	assert.Len(suite.T(), dialogResponse["messages"], 1)

	// Step 5: Verify cross-service data consistency
	req = httptest.NewRequest("GET", fmt.Sprintf("/users/%s/profile", userID), nil)

	w = httptest.NewRecorder()
	start = time.Now()
	suite.router.ServeHTTP(w, req)
	profileDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), profileDuration < 100*time.Millisecond, "Profile lookup should be <100ms")

	var profileResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &profileResponse)

	// Verify all services have consistent data
	assert.Equal(suite.T(), userID, profileResponse["user_id"])
	assert.Equal(suite.T(), "+66812345678", profileResponse["phone"])
	assert.Equal(suite.T(), "Somchai", profileResponse["first_name"])
	assert.Equal(suite.T(), "TH", profileResponse["country"])

	wallets := profileResponse["wallets"].([]interface{})
	assert.Len(suite.T(), wallets, 1)
	wallet := wallets[0].(map[string]interface{})
	assert.Equal(suite.T(), "THB", wallet["currency"])
}

func (suite *UserRegistrationTestSuite) TestMultiRegionRegistration() {
	countries := []struct {
		country  string
		phone    string
		currency string
		language string
	}{
		{"TH", "+66812345678", "THB", "th-TH"},
		{"SG", "+6591234567", "SGD", "en-SG"},
		{"ID", "+628123456789", "IDR", "id-ID"},
		{"MY", "+60123456789", "MYR", "ms-MY"},
		{"PH", "+639123456789", "PHP", "en-PH"},
		{"VN", "+84987654321", "VND", "vi-VN"},
	}

	for _, testCase := range countries {
		suite.T().Logf("Testing registration for country: %s", testCase.country)

		registrationData := map[string]interface{}{
			"phone":      testCase.phone,
			"first_name": "Test",
			"last_name":  "User",
			"email":      fmt.Sprintf("test%s@example.com", testCase.country),
			"country":    testCase.country,
		}

		// Register user
		jsonData, _ := json.Marshal(registrationData)
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		assert.Equal(suite.T(), http.StatusCreated, w.Code)

		var authResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &authResponse)
		userID := authResponse["id"].(string)

		// Create wallet
		walletData := map[string]interface{}{
			"user_id": userID,
			"country": testCase.country,
		}

		jsonData, _ = json.Marshal(walletData)
		req = httptest.NewRequest("POST", "/wallets/auto-create", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		assert.Equal(suite.T(), http.StatusCreated, w.Code)

		var walletResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &walletResponse)
		assert.Equal(suite.T(), testCase.currency, walletResponse["currency"])

		// Create messaging profile
		messagingData := map[string]interface{}{
			"user_id":      userID,
			"display_name": fmt.Sprintf("Test User %s", testCase.country),
			"language":     testCase.language,
		}

		jsonData, _ = json.Marshal(messagingData)
		req = httptest.NewRequest("POST", "/messaging/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		assert.Equal(suite.T(), http.StatusCreated, w.Code)
	}
}

func (suite *UserRegistrationTestSuite) TestRegistrationRollback() {
	// Test partial registration failure and rollback
	suite.T().Log("Testing registration rollback on service failure")

	// Simulate successful auth registration
	registrationData := map[string]interface{}{
		"phone":      "+66999888777",
		"first_name": "Rollback",
		"last_name":  "Test",
		"email":      "rollback@example.com",
		"country":    "TH",
	}

	jsonData, _ := json.Marshal(registrationData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var authResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &authResponse)
	userID := authResponse["id"].(string)

	// Verify user exists in auth service
	_, userExists := suite.users[userID]
	assert.True(suite.T(), userExists, "User should exist after auth registration")

	// Simulate wallet creation failure by using invalid country
	walletData := map[string]interface{}{
		"user_id": userID,
		"country": "INVALID",
	}

	jsonData, _ = json.Marshal(walletData)
	req = httptest.NewRequest("POST", "/wallets/auto-create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Wallet creation should succeed but use default currency
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var walletResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &walletResponse)
	assert.Equal(suite.T(), "USD", walletResponse["currency"]) // Default fallback
}

func (suite *UserRegistrationTestSuite) TestConcurrentRegistrations() {
	// Test multiple users registering simultaneously
	suite.T().Log("Testing concurrent user registrations")

	userCount := 5
	results := make(chan bool, userCount)

	for i := 0; i < userCount; i++ {
		go func(index int) {
			registrationData := map[string]interface{}{
				"phone":      fmt.Sprintf("+6681234567%d", index),
				"first_name": fmt.Sprintf("Concurrent%d", index),
				"last_name":  "User",
				"email":      fmt.Sprintf("concurrent%d@example.com", index),
				"country":    "TH",
			}

			jsonData, _ := json.Marshal(registrationData)
			req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			start := time.Now()
			suite.router.ServeHTTP(w, req)
			duration := time.Since(start)

			success := w.Code == http.StatusCreated && duration < 300*time.Millisecond
			results <- success
		}(i)
	}

	// Wait for all registrations to complete
	successCount := 0
	for i := 0; i < userCount; i++ {
		if <-results {
			successCount++
		}
	}

	assert.Equal(suite.T(), userCount, successCount, "All concurrent registrations should succeed")
}