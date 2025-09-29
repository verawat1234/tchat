package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// PresenceCheckContractTestSuite defines contract tests for GET /presence/check/{user_id}
type PresenceCheckContractTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *PresenceCheckContractTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	// Note: Endpoints not implemented yet - tests should fail
}

func (suite *PresenceCheckContractTestSuite) TestCheckPresence_Available() {
	// Test case: User is available for calls
	userID := "available-123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/presence/check/"+userID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: GET /presence/check/{user_id} should return 200 OK
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	// Contract: Response should indicate availability
	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), userID, data["user_id"])
	assert.Equal(suite.T(), true, data["available"])
	assert.Equal(suite.T(), "online", data["status"])
	assert.Equal(suite.T(), "User is available for calls", data["message"])
}

func (suite *PresenceCheckContractTestSuite) TestCheckPresence_Busy() {
	// Test case: User is busy
	userID := "busy-user-123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/presence/check/"+userID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), false, data["available"])
	assert.Equal(suite.T(), "busy", data["status"])
	assert.Equal(suite.T(), "User is busy", data["message"])
}

func (suite *PresenceCheckContractTestSuite) TestCheckPresence_InCall() {
	// Test case: User is in an active call
	userID := "incall-user-123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/presence/check/"+userID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), false, data["available"])
	assert.Equal(suite.T(), "in_call", data["status"])
	assert.Equal(suite.T(), "User is currently in a call", data["message"])
	assert.NotEmpty(suite.T(), data["call_id"])
}

func (suite *PresenceCheckContractTestSuite) TestCheckPresence_Offline() {
	// Test case: User is offline
	userID := "offline-user-123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/presence/check/"+userID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), false, data["available"])
	assert.Equal(suite.T(), "offline", data["status"])
	assert.Equal(suite.T(), "User is offline", data["message"])
	assert.NotEmpty(suite.T(), data["last_seen"])
}

func (suite *PresenceCheckContractTestSuite) TestCheckPresence_InvalidUserID() {
	// Test case: Invalid UUID format
	userID := "invalid-uuid-format"

	req, _ := http.NewRequest("GET", "/api/v1/presence/check/"+userID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Invalid UUID should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "Invalid user ID")
}

func (suite *PresenceCheckContractTestSuite) TestCheckPresence_UserNotFound() {
	// Test case: Non-existent user
	userID := "nonexistent-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/presence/check/"+userID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Non-existent user should return 404 Not Found
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "User not found", response["error"])
}

func (suite *PresenceCheckContractTestSuite) TestCheckPresence_Unauthorized() {
	// Test case: Missing authorization header
	userID := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/presence/check/"+userID, nil)
	// Missing Authorization header
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Unauthorized request should return 401 Unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

// TestPresenceCheckContractTestSuite runs the contract test suite
func TestPresenceCheckContractTestSuite(t *testing.T) {
	suite.Run(t, new(PresenceCheckContractTestSuite))
}