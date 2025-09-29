package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// PresenceUpdateContractTestSuite defines contract tests for PUT /presence/status
type PresenceUpdateContractTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *PresenceUpdateContractTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	// Note: Endpoints not implemented yet - tests should fail
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_ToOnline() {
	// Test case: Update user status to online
	requestBody := map[string]interface{}{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"status":  "online",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: PUT /presence/status should return 200 OK
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	// Contract: Response should contain updated presence status
	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), "123e4567-e89b-12d3-a456-426614174000", data["user_id"])
	assert.Equal(suite.T(), "online", data["status"])
	assert.NotEmpty(suite.T(), data["last_seen"])
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_ToBusy() {
	// Test case: Update user status to busy
	requestBody := map[string]interface{}{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"status":  "busy",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), "busy", data["status"])
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_ToOffline() {
	// Test case: Update user status to offline
	requestBody := map[string]interface{}{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"status":  "offline",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), "offline", data["status"])
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_InvalidStatus() {
	// Test case: Invalid status value
	requestBody := map[string]interface{}{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"status":  "invalid_status",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Invalid status should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "Invalid status")
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_MissingUserID() {
	// Test case: Missing user_id field
	requestBody := map[string]interface{}{
		"status": "online",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Missing user_id should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "user_id")
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_InvalidUserID() {
	// Test case: Invalid UUID format
	requestBody := map[string]interface{}{
		"user_id": "invalid-uuid-format",
		"status":  "online",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Invalid UUID should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "Invalid user ID")
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_Unauthorized() {
	// Test case: Missing authorization header
	requestBody := map[string]interface{}{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"status":  "online",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	// Missing Authorization header
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Unauthorized request should return 401 Unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_UserNotFound() {
	// Test case: Non-existent user
	requestBody := map[string]interface{}{
		"user_id": "nonexistent-e89b-12d3-a456-426614174000",
		"status":  "online",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Non-existent user should return 404 Not Found
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "User not found", response["error"])
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_CannotSetInCall() {
	// Test case: Attempt to manually set in_call status
	requestBody := map[string]interface{}{
		"user_id": "123e4567-e89b-12d3-a456-426614174000",
		"status":  "in_call",
	}

	bodyBytes, err := json.Marshal(requestBody)
	suite.Require().NoError(err)

	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Cannot manually set in_call status - should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "in_call status is managed automatically")
}

func (suite *PresenceUpdateContractTestSuite) TestUpdatePresenceStatus_InvalidJSON() {
	// Test case: Invalid JSON payload
	req, _ := http.NewRequest("PUT", "/api/v1/presence/status", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Invalid JSON should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "Invalid JSON")
}

// TestPresenceUpdateContractTestSuite runs the contract test suite
func TestPresenceUpdateContractTestSuite(t *testing.T) {
	suite.Run(t, new(PresenceUpdateContractTestSuite))
}