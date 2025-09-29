package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"tchat.dev/calling/handlers"
	"tchat.dev/calling/services"
)

// CallFailureIntegrationTestSuite tests various call failure scenarios and recovery
type CallFailureIntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	server *httptest.Server
}

func (suite *CallFailureIntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// This will fail until all services are properly implemented and integrated
	callService := &services.CallService{}
	presenceService := &services.PresenceService{}
	signalingService := &services.SignalingService{}

	callHandlers := handlers.NewCallHandlers(callService)
	presenceHandlers := handlers.NewPresenceHandlers(presenceService)
	wsHandler := handlers.NewWebSocketHandler(signalingService)

	// Setup router with all handlers
	suite.router = gin.New()
	api := suite.router.Group("/api/v1")

	callHandlers.RegisterRoutes(api)
	presenceHandlers.RegisterRoutes(api)
	wsHandler.RegisterRoutes(api)

	suite.server = httptest.NewServer(suite.router)
}

func (suite *CallFailureIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *CallFailureIntegrationTestSuite) TestCallFailure_UserOffline() {
	// Integration Test: Call failure when callee is offline
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set caller online, leave callee offline
	suite.setUserPresence(callerID, "online")
	// callee is offline (not set)

	// Attempt to initiate call - should fail
	response := suite.attemptCallInitiation(callerID, calleeID, "voice")
	assert.Equal(suite.T(), http.StatusConflict, response.Code) // 409 Conflict

	var errorResponse map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &errorResponse)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), errorResponse["error"], "offline")
	assert.Equal(suite.T(), "USER_OFFLINE", errorResponse["error_code"])

	// Verify caller presence is still online (unchanged)
	callerPresence := suite.getUserPresence(callerID)
	assert.Equal(suite.T(), "online", callerPresence["status"])

	// Verify no call history is created for failed initiation
	callerHistory := suite.getCallHistory(callerID)
	assert.Equal(suite.T(), 0, len(callerHistory))
}

func (suite *CallFailureIntegrationTestSuite) TestCallFailure_UserBusy() {
	// Integration Test: Call failure when callee is busy
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set caller online, callee busy (in another call)
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "busy")

	// Attempt to initiate call - should fail
	response := suite.attemptCallInitiation(callerID, calleeID, "video")
	assert.Equal(suite.T(), http.StatusConflict, response.Code) // 409 Conflict

	var errorResponse map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &errorResponse)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), errorResponse["error"], "busy")
	assert.Equal(suite.T(), "USER_BUSY", errorResponse["error_code"])

	// Verify presence states are unchanged
	callerPresence := suite.getUserPresence(callerID)
	calleePresence := suite.getUserPresence(calleeID)
	assert.Equal(suite.T(), "online", callerPresence["status"])
	assert.Equal(suite.T(), "busy", calleePresence["status"])
}

func (suite *CallFailureIntegrationTestSuite) TestCallFailure_CallDeclined() {
	// Integration Test: Call declined by callee
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set both users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Initiate call successfully
	callID := suite.initiateCall(callerID, calleeID, "voice")
	suite.Require().NotEmpty(callID)

	// Verify call is in connecting state
	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "connecting", callStatus["status"])

	// Callee declines the call
	suite.answerCall(callID, calleeID, false) // false = decline

	// Verify call status is failed with decline reason
	callStatus = suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "failed", callStatus["status"])
	assert.Equal(suite.T(), "declined", callStatus["failure_reason"])
	assert.NotNil(suite.T(), callStatus["ended_at"])

	// Verify users are back to online status
	callerPresence := suite.getUserPresence(callerID)
	calleePresence := suite.getUserPresence(calleeID)
	assert.Equal(suite.T(), "online", callerPresence["status"])
	assert.Equal(suite.T(), "online", calleePresence["status"])
	assert.Nil(suite.T(), callerPresence["call_id"])
	assert.Nil(suite.T(), calleePresence["call_id"])

	// Verify call history shows declined call
	callerHistory := suite.getCallHistory(callerID)
	calleeHistory := suite.getCallHistory(calleeID)

	var callerHistoryItem, calleeHistoryItem map[string]interface{}
	for _, item := range callerHistory {
		if item.(map[string]interface{})["call_session_id"] == callID {
			callerHistoryItem = item.(map[string]interface{})
			break
		}
	}
	for _, item := range calleeHistory {
		if item.(map[string]interface{})["call_session_id"] == callID {
			calleeHistoryItem = item.(map[string]interface{})
			break
		}
	}

	assert.NotNil(suite.T(), callerHistoryItem)
	assert.NotNil(suite.T(), calleeHistoryItem)
	assert.Equal(suite.T(), "declined", callerHistoryItem["outcome"])
	assert.Equal(suite.T(), "declined", calleeHistoryItem["outcome"])
	assert.Equal(suite.T(), "outgoing", callerHistoryItem["direction"])
	assert.Equal(suite.T(), "incoming", calleeHistoryItem["direction"])
}

func (suite *CallFailureIntegrationTestSuite) TestCallFailure_CallTimeout() {
	// Integration Test: Call times out without answer
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set both users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Initiate call
	callID := suite.initiateCall(callerID, calleeID, "voice")

	// Simulate timeout by waiting and then checking status
	// In real implementation, the service would automatically timeout after configured duration
	time.Sleep(100 * time.Millisecond) // Simulate timeout period

	// Force timeout via API call for testing
	suite.forceCallTimeout(callID)

	// Verify call status is failed with timeout reason
	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "failed", callStatus["status"])
	assert.Equal(suite.T(), "timeout", callStatus["failure_reason"])

	// Verify users are back online
	callerPresence := suite.getUserPresence(callerID)
	calleePresence := suite.getUserPresence(calleeID)
	assert.Equal(suite.T(), "online", callerPresence["status"])
	assert.Equal(suite.T(), "online", calleePresence["status"])

	// Verify timeout is recorded in history
	callerHistory := suite.getCallHistory(callerID)
	var historyItem map[string]interface{}
	for _, item := range callerHistory {
		if item.(map[string]interface{})["call_session_id"] == callID {
			historyItem = item.(map[string]interface{})
			break
		}
	}
	assert.Equal(suite.T(), "timeout", historyItem["outcome"])
}

func (suite *CallFailureIntegrationTestSuite) TestCallFailure_NetworkIssues() {
	// Integration Test: Call failure due to network connectivity issues
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set users online and establish call
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")
	callID := suite.initiateCall(callerID, calleeID, "video")
	suite.answerCall(callID, calleeID, true)

	// Verify call is active
	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "active", callStatus["status"])

	// Simulate network failure
	suite.simulateNetworkFailure(callID, callerID, "connection_lost")

	// Verify call status shows network failure
	callStatus = suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "failed", callStatus["status"])
	assert.Equal(suite.T(), "network_failure", callStatus["failure_reason"])
	assert.Equal(suite.T(), "connection_lost", callStatus["network_error"])

	// Verify users are set back to online
	callerPresence := suite.getUserPresence(callerID)
	calleePresence := suite.getUserPresence(calleeID)
	assert.Equal(suite.T(), "online", callerPresence["status"])
	assert.Equal(suite.T(), "online", calleePresence["status"])

	// Verify network failure is recorded in history
	callerHistory := suite.getCallHistory(callerID)
	var historyItem map[string]interface{}
	for _, item := range callerHistory {
		if item.(map[string]interface{})["call_session_id"] == callID {
			historyItem = item.(map[string]interface{})
			break
		}
	}
	assert.Equal(suite.T(), "network_failure", historyItem["outcome"])
	assert.Greater(suite.T(), historyItem["duration"], 0) // Some duration before failure
}

func (suite *CallFailureIntegrationTestSuite) TestCallFailure_ServiceUnavailable() {
	// Integration Test: Call failure when service dependencies are unavailable
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Simulate service unavailability (e.g., Redis down, DB connection lost)
	suite.simulateServiceFailure("redis_unavailable")

	// Attempt to initiate call - should fail with service error
	response := suite.attemptCallInitiation(callerID, calleeID, "voice")
	assert.Equal(suite.T(), http.StatusServiceUnavailable, response.Code) // 503 Service Unavailable

	var errorResponse map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &errorResponse)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), errorResponse["error"], "service unavailable")
	assert.Equal(suite.T(), "SERVICE_UNAVAILABLE", errorResponse["error_code"])

	// Restore service
	suite.restoreService()

	// Verify that call can now be initiated successfully
	callID := suite.initiateCall(callerID, calleeID, "voice")
	suite.Require().NotEmpty(callID)

	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "connecting", callStatus["status"])
}

func (suite *CallFailureIntegrationTestSuite) TestCallFailure_InvalidParameters() {
	// Integration Test: Call failure with invalid parameters
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Test invalid call type
	response := suite.attemptCallInitiation(callerID, calleeID, "invalid_type")
	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &errorResponse)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), errorResponse["error"], "invalid call type")

	// Test calling self
	response = suite.attemptCallInitiation(callerID, callerID, "voice")
	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)

	err = json.Unmarshal(response.Body.Bytes(), &errorResponse)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), errorResponse["error"], "cannot call yourself")

	// Test invalid UUID format
	invalidID := "invalid-uuid-format"
	requestBody := map[string]interface{}{
		"caller_id": callerID.String(),
		"callee_id": invalidID,
		"call_type": "voice",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/calls/initiate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), errorResponse["error"], "invalid UUID")
}

func (suite *CallFailureIntegrationTestSuite) TestCallFailure_RecoveryScenarios() {
	// Integration Test: Call recovery after transient failures
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Establish call
	callID := suite.initiateCall(callerID, calleeID, "video")
	suite.answerCall(callID, calleeID, true)

	// Simulate temporary network issue
	suite.simulateTemporaryNetworkIssue(callID, callerID)

	// Verify call status shows reconnecting
	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "reconnecting", callStatus["status"])

	// Simulate recovery
	suite.simulateNetworkRecovery(callID, callerID)

	// Verify call is back to active
	callStatus = suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "active", callStatus["status"])

	// End call normally
	suite.endCall(callID, callerID)

	// Verify history shows successful completion despite temporary issues
	callerHistory := suite.getCallHistory(callerID)
	var historyItem map[string]interface{}
	for _, item := range callerHistory {
		if item.(map[string]interface{})["call_session_id"] == callID {
			historyItem = item.(map[string]interface{})
			break
		}
	}
	assert.Equal(suite.T(), "completed", historyItem["outcome"])
	assert.Greater(suite.T(), historyItem["reconnection_count"], 0)
}

// Helper methods for call failure testing

func (suite *CallFailureIntegrationTestSuite) setUserPresence(userID uuid.UUID, status string) {
	requestBody := map[string]interface{}{
		"user_id": userID.String(),
		"status":  status,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", suite.server.URL+"/api/v1/presence/status", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CallFailureIntegrationTestSuite) getUserPresence(userID uuid.UUID) map[string]interface{} {
	req, _ := http.NewRequest("GET", suite.server.URL+"/api/v1/presence/status?user_id="+userID.String(), nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	return response["data"].(map[string]interface{})
}

func (suite *CallFailureIntegrationTestSuite) attemptCallInitiation(callerID, calleeID uuid.UUID, callType string) *httptest.ResponseRecorder {
	requestBody := map[string]interface{}{
		"caller_id": callerID.String(),
		"callee_id": calleeID.String(),
		"call_type": callType,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/calls/initiate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *CallFailureIntegrationTestSuite) initiateCall(callerID, calleeID uuid.UUID, callType string) string {
	requestBody := map[string]interface{}{
		"caller_id": callerID.String(),
		"callee_id": calleeID.String(),
		"call_type": callType,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/initiate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	return data["id"].(string)
}

func (suite *CallFailureIntegrationTestSuite) getCallStatus(callID string, userID uuid.UUID) map[string]interface{} {
	req, _ := http.NewRequest("GET", suite.server.URL+"/api/v1/calls/"+callID+"/status", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	return response["data"].(map[string]interface{})
}

func (suite *CallFailureIntegrationTestSuite) answerCall(callID string, userID uuid.UUID, accept bool) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
		"accept":  accept,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/answer", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CallFailureIntegrationTestSuite) endCall(callID string, userID uuid.UUID) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/end", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CallFailureIntegrationTestSuite) getCallHistory(userID uuid.UUID) []interface{} {
	req, _ := http.NewRequest("GET", suite.server.URL+"/api/v1/history?user_id="+userID.String(), nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	return data["items"].([]interface{})
}

func (suite *CallFailureIntegrationTestSuite) forceCallTimeout(callID string) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"reason":  "timeout",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/force-timeout", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CallFailureIntegrationTestSuite) simulateNetworkFailure(callID string, userID uuid.UUID, errorType string) {
	requestBody := map[string]interface{}{
		"call_id":    callID,
		"user_id":    userID.String(),
		"error_type": errorType,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/simulate-network-failure", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CallFailureIntegrationTestSuite) simulateServiceFailure(serviceType string) {
	requestBody := map[string]interface{}{
		"service_type": serviceType,
		"action":       "fail",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/admin/simulate-service-failure", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CallFailureIntegrationTestSuite) restoreService() {
	requestBody := map[string]interface{}{
		"action": "restore",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/admin/restore-services", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer admin-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CallFailureIntegrationTestSuite) simulateTemporaryNetworkIssue(callID string, userID uuid.UUID) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
		"type":    "temporary",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/simulate-network-issue", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *CallFailureIntegrationTestSuite) simulateNetworkRecovery(callID string, userID uuid.UUID) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/simulate-network-recovery", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			_ = err
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

// TestCallFailureIntegrationTestSuite runs the call failure integration test suite
func TestCallFailureIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CallFailureIntegrationTestSuite))
}