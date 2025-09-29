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

// PresenceGetContractTestSuite defines contract tests for GET /presence/status
type PresenceGetContractTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *PresenceGetContractTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	// Note: Endpoints not implemented yet - tests should fail
}

func (suite *PresenceGetContractTestSuite) TestGetPresenceStatus_Success() {
	// Test case: Valid user presence status request
	req, _ := http.NewRequest("GET", "/api/v1/presence/status?user_id=123e4567-e89b-12d3-a456-426614174000", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: GET /presence/status should return 200 OK
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	// Contract: Response should contain presence status details
	data := response["data"].(map[string]interface{})
	assert.NotEmpty(suite.T(), data["user_id"])
	assert.Contains(suite.T(), []string{"online", "offline", "busy", "in_call"}, data["status"])
	assert.NotEmpty(suite.T(), data["last_seen"])

	// Optional fields based on status
	if data["status"] == "in_call" {
		assert.NotEmpty(suite.T(), data["call_id"])
	}
}

func (suite *PresenceGetContractTestSuite) TestGetPresenceStatus_MissingUserID() {
	// Test case: Missing user_id parameter
	req, _ := http.NewRequest("GET", "/api/v1/presence/status", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Missing user_id should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "user_id")
}

func (suite *PresenceGetContractTestSuite) TestGetPresenceStatus_InvalidUserID() {
	// Test case: Invalid UUID format
	req, _ := http.NewRequest("GET", "/api/v1/presence/status?user_id=invalid-uuid", nil)
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

func (suite *PresenceGetContractTestSuite) TestGetPresenceStatus_UserNotFound() {
	// Test case: Non-existent user
	req, _ := http.NewRequest("GET", "/api/v1/presence/status?user_id=nonexistent-e89b-12d3-a456-426614174000", nil)
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

func (suite *PresenceGetContractTestSuite) TestGetPresenceStatus_Unauthorized() {
	// Test case: Missing authorization header
	req, _ := http.NewRequest("GET", "/api/v1/presence/status?user_id=123e4567-e89b-12d3-a456-426614174000", nil)
	// Missing Authorization header
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Unauthorized request should return 401 Unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *PresenceGetContractTestSuite) TestGetPresenceStatus_OnlineUser() {
	// Test case: User with online status
	req, _ := http.NewRequest("GET", "/api/v1/presence/status?user_id=online-user-e89b-12d3-a456-426614174000", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), "online", data["status"])
	assert.Nil(suite.T(), data["call_id"]) // Should be null when not in call
}

func (suite *PresenceGetContractTestSuite) TestGetPresenceStatus_InCallUser() {
	// Test case: User currently in a call
	req, _ := http.NewRequest("GET", "/api/v1/presence/status?user_id=incall-user-e89b-12d3-a456-426614174000", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), "in_call", data["status"])
	assert.NotEmpty(suite.T(), data["call_id"]) // Should have call_id when in call
}

func (suite *PresenceGetContractTestSuite) TestGetPresenceStatus_OfflineUser() {
	// Test case: User with offline status
	req, _ := http.NewRequest("GET", "/api/v1/presence/status?user_id=offline-user-e89b-12d3-a456-426614174000", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), "offline", data["status"])
	assert.NotEmpty(suite.T(), data["last_seen"]) // Should have last_seen timestamp
}

// TestPresenceGetContractTestSuite runs the contract test suite
func TestPresenceGetContractTestSuite(t *testing.T) {
	suite.Run(t, new(PresenceGetContractTestSuite))
}