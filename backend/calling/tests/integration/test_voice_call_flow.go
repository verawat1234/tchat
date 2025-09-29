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

// VoiceCallFlowIntegrationTestSuite tests complete voice call flows
type VoiceCallFlowIntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	server *httptest.Server
}

func (suite *VoiceCallFlowIntegrationTestSuite) SetupSuite() {
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

func (suite *VoiceCallFlowIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *VoiceCallFlowIntegrationTestSuite) TestVoiceCallSuccessFlow() {
	// Integration Test: Complete voice call flow from initiation to completion
	callerID := uuid.New()
	calleeID := uuid.New()

	// Step 1: Set both users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Step 2: Initiate voice call
	callID := suite.initiateCall(callerID, calleeID, "voice")
	suite.Require().NotEmpty(callID)

	// Step 3: Verify call status is connecting
	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "connecting", callStatus["status"])
	assert.Equal(suite.T(), "voice", callStatus["type"])
	assert.Equal(suite.T(), callerID.String(), callStatus["initiated_by"])

	// Step 4: Verify users are now in call status
	callerPresence := suite.getUserPresence(callerID)
	calleePresence := suite.getUserPresence(calleeID)
	assert.Equal(suite.T(), "in_call", callerPresence["status"])
	assert.Equal(suite.T(), "in_call", calleePresence["status"])
	assert.Equal(suite.T(), callID, callerPresence["call_id"])
	assert.Equal(suite.T(), callID, calleePresence["call_id"])

	// Step 5: Callee answers the call
	suite.answerCall(callID, calleeID, true)

	// Step 6: Verify call status is now active
	callStatus = suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "active", callStatus["status"])
	assert.Nil(suite.T(), callStatus["ended_at"])
	assert.Equal(suite.T(), 0, callStatus["duration"])

	// Step 7: Simulate call duration
	time.Sleep(100 * time.Millisecond) // Simulate brief call

	// Step 8: End the call
	suite.endCall(callID, callerID)

	// Step 9: Verify call is ended
	callStatus = suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "ended", callStatus["status"])
	assert.NotNil(suite.T(), callStatus["ended_at"])
	assert.Greater(suite.T(), callStatus["duration"], 0)

	// Step 10: Verify users are back online
	callerPresence = suite.getUserPresence(callerID)
	calleePresence = suite.getUserPresence(calleeID)
	assert.Equal(suite.T(), "online", callerPresence["status"])
	assert.Equal(suite.T(), "online", calleePresence["status"])
	assert.Nil(suite.T(), callerPresence["call_id"])
	assert.Nil(suite.T(), calleePresence["call_id"])

	// Step 11: Verify call history is created
	callerHistory := suite.getCallHistory(callerID)
	calleeHistory := suite.getCallHistory(calleeID)

	assert.GreaterOrEqual(suite.T(), len(callerHistory), 1)
	assert.GreaterOrEqual(suite.T(), len(calleeHistory), 1)

	// Find the call in history
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

	// Verify history details
	assert.NotNil(suite.T(), callerHistoryItem)
	assert.NotNil(suite.T(), calleeHistoryItem)
	assert.Equal(suite.T(), "voice", callerHistoryItem["call_type"])
	assert.Equal(suite.T(), "voice", calleeHistoryItem["call_type"])
	assert.Equal(suite.T(), "ended", callerHistoryItem["call_status"])
	assert.Equal(suite.T(), "ended", calleeHistoryItem["call_status"])
	assert.Equal(suite.T(), true, callerHistoryItem["initiated_by_me"])
	assert.Equal(suite.T(), false, calleeHistoryItem["initiated_by_me"])
	assert.Equal(suite.T(), "outgoing", callerHistoryItem["direction"])
	assert.Equal(suite.T(), "incoming", calleeHistoryItem["direction"])
	assert.Equal(suite.T(), "completed", callerHistoryItem["outcome"])
	assert.Equal(suite.T(), "completed", calleeHistoryItem["outcome"])
}

func (suite *VoiceCallFlowIntegrationTestSuite) TestVoiceCallDeclinedFlow() {
	// Integration Test: Voice call declined by callee
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set both users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Initiate voice call
	callID := suite.initiateCall(callerID, calleeID, "voice")

	// Callee declines the call
	suite.answerCall(callID, calleeID, false)

	// Verify call status is failed with decline reason
	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "failed", callStatus["status"])
	assert.Equal(suite.T(), "declined", callStatus["failure_reason"])

	// Verify users are back online
	callerPresence := suite.getUserPresence(callerID)
	calleePresence := suite.getUserPresence(calleeID)
	assert.Equal(suite.T(), "online", callerPresence["status"])
	assert.Equal(suite.T(), "online", calleePresence["status"])

	// Verify call history shows declined call
	callerHistory := suite.getCallHistory(callerID)
	calleeHistory := suite.getCallHistory(calleeID)

	// Find call in history and verify outcome
	for _, item := range callerHistory {
		if item.(map[string]interface{})["call_session_id"] == callID {
			assert.Equal(suite.T(), "declined", item.(map[string]interface{})["outcome"])
			break
		}
	}
	for _, item := range calleeHistory {
		if item.(map[string]interface{})["call_session_id"] == callID {
			assert.Equal(suite.T(), "declined", item.(map[string]interface{})["outcome"])
			break
		}
	}
}

func (suite *VoiceCallFlowIntegrationTestSuite) TestVoiceCallBusyFlow() {
	// Integration Test: Call to busy user
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set caller online, callee busy
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "busy")

	// Attempt to initiate call - should fail
	response := suite.attemptCallInitiation(callerID, calleeID, "voice")
	assert.Equal(suite.T(), http.StatusConflict, response.Code) // 409 Conflict

	var errorResponse map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &errorResponse)
	suite.Require().NoError(err)
	assert.Contains(suite.T(), errorResponse["error"], "not available")
}

// Helper methods

func (suite *VoiceCallFlowIntegrationTestSuite) setUserPresence(userID uuid.UUID, status string) {
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
			// Log error but don't fail test
			_ = err // Avoid unused variable warning
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *VoiceCallFlowIntegrationTestSuite) getUserPresence(userID uuid.UUID) map[string]interface{} {
	req, _ := http.NewRequest("GET", suite.server.URL+"/api/v1/presence/status?user_id="+userID.String(), nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't fail test
			_ = err // Avoid unused variable warning
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	return response["data"].(map[string]interface{})
}

func (suite *VoiceCallFlowIntegrationTestSuite) initiateCall(callerID, calleeID uuid.UUID, callType string) string {
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
			// Log error but don't fail test
			_ = err // Avoid unused variable warning
		}
	}()
	suite.Require().Equal(http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	return data["id"].(string)
}

func (suite *VoiceCallFlowIntegrationTestSuite) attemptCallInitiation(callerID, calleeID uuid.UUID, callType string) *httptest.ResponseRecorder {
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

func (suite *VoiceCallFlowIntegrationTestSuite) getCallStatus(callID string, userID uuid.UUID) map[string]interface{} {
	req, _ := http.NewRequest("GET", suite.server.URL+"/api/v1/calls/"+callID+"/status", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't fail test
			_ = err // Avoid unused variable warning
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	return response["data"].(map[string]interface{})
}

func (suite *VoiceCallFlowIntegrationTestSuite) answerCall(callID string, userID uuid.UUID, accept bool) {
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
			// Log error but don't fail test
			_ = err // Avoid unused variable warning
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *VoiceCallFlowIntegrationTestSuite) endCall(callID string, userID uuid.UUID) {
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
			// Log error but don't fail test
			_ = err // Avoid unused variable warning
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (suite *VoiceCallFlowIntegrationTestSuite) getCallHistory(userID uuid.UUID) []interface{} {
	req, _ := http.NewRequest("GET", suite.server.URL+"/api/v1/history?user_id="+userID.String(), nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't fail test
			_ = err // Avoid unused variable warning
		}
	}()
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	data := response["data"].(map[string]interface{})
	return data["items"].([]interface{})
}

// TestVoiceCallFlowIntegrationTestSuite runs the integration test suite
func TestVoiceCallFlowIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(VoiceCallFlowIntegrationTestSuite))
}