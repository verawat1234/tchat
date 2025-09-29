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

// CallStatusContractTestSuite defines contract tests for GET /calls/{id}/status
type CallStatusContractTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *CallStatusContractTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	// Note: Endpoints not implemented yet - tests should fail
}

func (suite *CallStatusContractTestSuite) TestGetCallStatus_Success() {
	// Test case: Valid call ID returns call status
	callID := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/calls/"+callID+"/status", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: GET /calls/{id}/status should return 200 OK
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	// Contract: Response should contain call status details
	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), callID, data["id"])
	assert.Contains(suite.T(), []string{"connecting", "active", "ended", "failed"}, data["status"])
	assert.NotEmpty(suite.T(), data["type"])
	assert.NotEmpty(suite.T(), data["initiated_by"])
	assert.NotEmpty(suite.T(), data["started_at"])
	assert.NotEmpty(suite.T(), data["participants"])

	// Verify participants structure
	participants := data["participants"].([]interface{})
	assert.GreaterOrEqual(suite.T(), len(participants), 1)

	for _, p := range participants {
		participant := p.(map[string]interface{})
		assert.NotEmpty(suite.T(), participant["user_id"])
		assert.Contains(suite.T(), []string{"caller", "callee"}, participant["role"])
		assert.NotNil(suite.T(), participant["audio_enabled"])
		assert.NotNil(suite.T(), participant["video_enabled"])
		assert.Contains(suite.T(), []string{"good", "poor", "disconnected"}, participant["connection_quality"])
	}
}

func (suite *CallStatusContractTestSuite) TestGetCallStatus_CallNotFound() {
	// Test case: Non-existent call ID
	callID := "nonexistent-call-id"

	req, _ := http.NewRequest("GET", "/api/v1/calls/"+callID+"/status", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Non-existent call should return 404 Not Found
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), "Call not found", response["error"])
}

func (suite *CallStatusContractTestSuite) TestGetCallStatus_InvalidCallID() {
	// Test case: Invalid UUID format
	callID := "invalid-uuid-format"

	req, _ := http.NewRequest("GET", "/api/v1/calls/"+callID+"/status", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Invalid UUID should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "Invalid call ID")
}

func (suite *CallStatusContractTestSuite) TestGetCallStatus_Unauthorized() {
	// Test case: Missing or invalid authorization
	callID := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/calls/"+callID+"/status", nil)
	// Missing Authorization header
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Unauthorized request should return 401 Unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *CallStatusContractTestSuite) TestGetCallStatus_ForbiddenUser() {
	// Test case: User not participant in the call
	callID := "123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/calls/"+callID+"/status", nil)
	req.Header.Set("Authorization", "Bearer token-for-non-participant")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Non-participant should return 403 Forbidden
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "not authorized")
}

func (suite *CallStatusContractTestSuite) TestGetCallStatus_ActiveCall() {
	// Test case: Verify active call status details
	callID := "active-call-123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/calls/"+callID+"/status", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), "active", data["status"])
	assert.Nil(suite.T(), data["ended_at"])   // Should be null for active call
	assert.Zero(suite.T(), data["duration"])  // Should be 0 for active call
}

func (suite *CallStatusContractTestSuite) TestGetCallStatus_EndedCall() {
	// Test case: Verify ended call status details
	callID := "ended-call-123e4567-e89b-12d3-a456-426614174000"

	req, _ := http.NewRequest("GET", "/api/v1/calls/"+callID+"/status", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), "ended", data["status"])
	assert.NotEmpty(suite.T(), data["ended_at"])    // Should have end time
	assert.Greater(suite.T(), data["duration"], 0) // Should have positive duration
}

// TestCallStatusContractTestSuite runs the contract test suite
func TestCallStatusContractTestSuite(t *testing.T) {
	suite.Run(t, new(CallStatusContractTestSuite))
}