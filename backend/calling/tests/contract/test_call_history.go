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

// CallHistoryContractTestSuite defines contract tests for GET /history
type CallHistoryContractTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *CallHistoryContractTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	// Note: Endpoints not implemented yet - tests should fail
}

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_Success() {
	// Test case: Valid call history request
	req, _ := http.NewRequest("GET", "/api/v1/history?user_id=123e4567-e89b-12d3-a456-426614174000", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: GET /history should return 200 OK
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	// Contract: Response should contain call history items
	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.Equal(suite.T(), 20, data["limit"])  // Default limit
	assert.Equal(suite.T(), 0, data["offset"])  // Default offset
	assert.GreaterOrEqual(suite.T(), len(items), 0)

	// Verify structure of history items
	if len(items) > 0 {
		item := items[0].(map[string]interface{})
		assert.NotEmpty(suite.T(), item["id"])
		assert.NotEmpty(suite.T(), item["call_session_id"])
		assert.NotEmpty(suite.T(), item["user_id"])
		assert.NotEmpty(suite.T(), item["other_participant_id"])
		assert.Contains(suite.T(), []string{"voice", "video"}, item["call_type"])
		assert.Contains(suite.T(), []string{"connecting", "active", "ended", "failed"}, item["call_status"])
		assert.NotNil(suite.T(), item["initiated_by_me"])
		assert.GreaterOrEqual(suite.T(), item["duration"], 0)
		assert.Contains(suite.T(), []string{"incoming", "outgoing"}, item["direction"])
		assert.Contains(suite.T(), []string{"completed", "missed", "failed", "declined"}, item["outcome"])
		assert.NotEmpty(suite.T(), item["created_at"])
		assert.NotEmpty(suite.T(), item["updated_at"])
	}
}

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_WithPagination() {
	// Test case: Call history with custom pagination
	req, _ := http.NewRequest("GET", "/api/v1/history?user_id=123e4567-e89b-12d3-a456-426614174000&limit=10&offset=5", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), 10, data["limit"])
	assert.Equal(suite.T(), 5, data["offset"])
}

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_WithCallTypeFilter() {
	// Test case: Filter by call type
	req, _ := http.NewRequest("GET", "/api/v1/history?user_id=123e4567-e89b-12d3-a456-426614174000&call_type=voice", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})

	// All items should be voice calls
	for _, item := range items {
		callItem := item.(map[string]interface{})
		assert.Equal(suite.T(), "voice", callItem["call_type"])
	}
}

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_MissingUserID() {
	// Test case: Missing user_id parameter
	req, _ := http.NewRequest("GET", "/api/v1/history", nil)
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

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_InvalidUserID() {
	// Test case: Invalid UUID format
	req, _ := http.NewRequest("GET", "/api/v1/history?user_id=invalid-uuid", nil)
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

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_InvalidPagination() {
	// Test case: Invalid pagination parameters
	req, _ := http.NewRequest("GET", "/api/v1/history?user_id=123e4567-e89b-12d3-a456-426614174000&limit=invalid", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Invalid pagination should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "limit")
}

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_ExceedsLimitBounds() {
	// Test case: Limit exceeds maximum allowed (100)
	req, _ := http.NewRequest("GET", "/api/v1/history?user_id=123e4567-e89b-12d3-a456-426614174000&limit=150", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Limit exceeding bounds should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "limit must be")
}

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_InvalidCallType() {
	// Test case: Invalid call_type filter
	req, _ := http.NewRequest("GET", "/api/v1/history?user_id=123e4567-e89b-12d3-a456-426614174000&call_type=invalid", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Invalid call_type should return 400 Bad Request
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	assert.Contains(suite.T(), response["error"], "call_type must be")
}

func (suite *CallHistoryContractTestSuite) TestGetCallHistory_Unauthorized() {
	// Test case: Missing authorization header
	req, _ := http.NewRequest("GET", "/api/v1/history?user_id=123e4567-e89b-12d3-a456-426614174000", nil)
	// Missing Authorization header
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Contract: Unauthorized request should return 401 Unauthorized
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *CallHistoryContractTestSuite) TestGetCallHistoryStats_Success() {
	// Test case: Get call statistics
	req, _ := http.NewRequest("GET", "/api/v1/history/stats?user_id=123e4567-e89b-12d3-a456-426614174000", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	// Contract: Response should contain call statistics
	data := response["data"].(map[string]interface{})
	assert.GreaterOrEqual(suite.T(), data["total_calls"], 0)
	assert.GreaterOrEqual(suite.T(), data["outgoing_calls"], 0)
	assert.GreaterOrEqual(suite.T(), data["incoming_calls"], 0)
	assert.GreaterOrEqual(suite.T(), data["successful_calls"], 0)
	assert.GreaterOrEqual(suite.T(), data["missed_calls"], 0)
	assert.GreaterOrEqual(suite.T(), data["voice_calls"], 0)
	assert.GreaterOrEqual(suite.T(), data["video_calls"], 0)
	assert.GreaterOrEqual(suite.T(), data["total_duration"], 0)
	assert.GreaterOrEqual(suite.T(), data["average_duration"], 0.0)
}

func (suite *CallHistoryContractTestSuite) TestGetRecentCallHistory_Success() {
	// Test case: Get recent call history
	req, _ := http.NewRequest("GET", "/api/v1/history/recent?user_id=123e4567-e89b-12d3-a456-426614174000&days=7&limit=5", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), 7, data["days"])
	assert.Equal(suite.T(), 5, data["limit"])

	items := data["items"].([]interface{})
	assert.LessOrEqual(suite.T(), len(items), 5)
}

// TestCallHistoryContractTestSuite runs the contract test suite
func TestCallHistoryContractTestSuite(t *testing.T) {
	suite.Run(t, new(CallHistoryContractTestSuite))
}