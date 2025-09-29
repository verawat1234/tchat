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

// VideoCallFlowIntegrationTestSuite tests complete video call flows with media controls
type VideoCallFlowIntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	server *httptest.Server
}

func (suite *VideoCallFlowIntegrationTestSuite) SetupSuite() {
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

func (suite *VideoCallFlowIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *VideoCallFlowIntegrationTestSuite) TestVideoCallSuccessFlowWithMediaControls() {
	// Integration Test: Complete video call flow with media control toggles
	callerID := uuid.New()
	calleeID := uuid.New()

	// Step 1: Set both users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Step 2: Initiate video call
	callID := suite.initiateCall(callerID, calleeID, "video")
	suite.Require().NotEmpty(callID)

	// Step 3: Verify call status is connecting with video type
	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "connecting", callStatus["status"])
	assert.Equal(suite.T(), "video", callStatus["type"])
	assert.Equal(suite.T(), callerID.String(), callStatus["initiated_by"])

	// Step 4: Verify media settings are enabled by default for video calls
	mediaSettings := callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), true, mediaSettings["video_enabled"])
	assert.Equal(suite.T(), true, mediaSettings["audio_enabled"])
	assert.Equal(suite.T(), "camera", mediaSettings["camera_facing"]) // front camera default
	assert.Equal(suite.T(), "720p", mediaSettings["video_quality"])

	// Step 5: Callee answers the video call
	suite.answerCall(callID, calleeID, true)

	// Step 6: Verify call status is now active
	callStatus = suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "active", callStatus["status"])

	// Step 7: Test media control toggles during active call

	// Toggle video off
	suite.toggleMedia(callID, callerID, "video", false)
	callStatus = suite.getCallStatus(callID, callerID)
	mediaSettings = callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), false, mediaSettings["video_enabled"])
	assert.Equal(suite.T(), true, mediaSettings["audio_enabled"]) // audio still on

	// Toggle video back on
	suite.toggleMedia(callID, callerID, "video", true)
	callStatus = suite.getCallStatus(callID, callerID)
	mediaSettings = callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), true, mediaSettings["video_enabled"])

	// Toggle audio off (mute)
	suite.toggleMedia(callID, callerID, "audio", false)
	callStatus = suite.getCallStatus(callID, callerID)
	mediaSettings = callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), false, mediaSettings["audio_enabled"])
	assert.Equal(suite.T(), true, mediaSettings["video_enabled"]) // video still on

	// Switch camera from front to back
	suite.switchCamera(callID, callerID, "back")
	callStatus = suite.getCallStatus(callID, callerID)
	mediaSettings = callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), "back", mediaSettings["camera_facing"])

	// Change video quality
	suite.changeVideoQuality(callID, callerID, "1080p")
	callStatus = suite.getCallStatus(callID, callerID)
	mediaSettings = callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), "1080p", mediaSettings["video_quality"])

	// Step 8: Toggle audio back on
	suite.toggleMedia(callID, callerID, "audio", true)
	callStatus = suite.getCallStatus(callID, callerID)
	mediaSettings = callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), true, mediaSettings["audio_enabled"])

	// Step 9: Simulate call duration
	time.Sleep(100 * time.Millisecond)

	// Step 10: End the video call
	suite.endCall(callID, callerID)

	// Step 11: Verify call is ended
	callStatus = suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "ended", callStatus["status"])
	assert.NotNil(suite.T(), callStatus["ended_at"])
	assert.Greater(suite.T(), callStatus["duration"], 0)

	// Step 12: Verify call history contains video call details
	callerHistory := suite.getCallHistory(callerID)
	calleeHistory := suite.getCallHistory(calleeID)

	// Find the video call in history
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

	// Verify video call history details
	assert.NotNil(suite.T(), callerHistoryItem)
	assert.NotNil(suite.T(), calleeHistoryItem)
	assert.Equal(suite.T(), "video", callerHistoryItem["call_type"])
	assert.Equal(suite.T(), "video", calleeHistoryItem["call_type"])
	assert.Equal(suite.T(), "completed", callerHistoryItem["outcome"])
	assert.Equal(suite.T(), "completed", calleeHistoryItem["outcome"])

	// Verify media usage statistics were recorded
	mediaStats := callerHistoryItem["media_stats"].(map[string]interface{})
	assert.Greater(suite.T(), mediaStats["video_time"], 0) // Video was on for some time
	assert.Greater(suite.T(), mediaStats["audio_time"], 0) // Audio was on for some time
	assert.Equal(suite.T(), 2, mediaStats["video_toggles"])   // Toggled off then on
	assert.Equal(suite.T(), 2, mediaStats["audio_toggles"])   // Toggled off then on
	assert.Equal(suite.T(), 1, mediaStats["camera_switches"]) // Switched front to back
	assert.Equal(suite.T(), 1, mediaStats["quality_changes"]) // Changed to 1080p
}

func (suite *VideoCallFlowIntegrationTestSuite) TestVideoCallDowngradeToAudio() {
	// Integration Test: Video call downgraded to audio-only due to bandwidth
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Initiate video call
	callID := suite.initiateCall(callerID, calleeID, "video")
	suite.answerCall(callID, calleeID, true)

	// Verify video call is active
	callStatus := suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "active", callStatus["status"])
	assert.Equal(suite.T(), "video", callStatus["type"])

	// Simulate bandwidth issues and automatic downgrade
	suite.downgradeCall(callID, callerID, "bandwidth_low")

	// Verify call is now audio-only but still active
	callStatus = suite.getCallStatus(callID, callerID)
	assert.Equal(suite.T(), "active", callStatus["status"])
	assert.Equal(suite.T(), "audio", callStatus["type"]) // Downgraded to audio

	mediaSettings := callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), false, mediaSettings["video_enabled"]) // Video disabled
	assert.Equal(suite.T(), true, mediaSettings["audio_enabled"])  // Audio still active
	assert.Equal(suite.T(), "bandwidth_low", callStatus["downgrade_reason"])

	// End call
	suite.endCall(callID, callerID)

	// Verify history shows original type was video but ended as audio
	callerHistory := suite.getCallHistory(callerID)
	var historyItem map[string]interface{}
	for _, item := range callerHistory {
		if item.(map[string]interface{})["call_session_id"] == callID {
			historyItem = item.(map[string]interface{})
			break
		}
	}

	assert.Equal(suite.T(), "video", historyItem["original_call_type"])
	assert.Equal(suite.T(), "audio", historyItem["final_call_type"])
	assert.Equal(suite.T(), "downgraded", historyItem["outcome"])
}

func (suite *VideoCallFlowIntegrationTestSuite) TestVideoCallWithScreenShare() {
	// Integration Test: Video call with screen sharing capability
	callerID := uuid.New()
	calleeID := uuid.New()

	// Set users online
	suite.setUserPresence(callerID, "online")
	suite.setUserPresence(calleeID, "online")

	// Initiate and answer video call
	callID := suite.initiateCall(callerID, calleeID, "video")
	suite.answerCall(callID, calleeID, true)

	// Start screen sharing
	suite.startScreenShare(callID, callerID)

	// Verify screen sharing is active
	callStatus := suite.getCallStatus(callID, callerID)
	mediaSettings := callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), true, mediaSettings["screen_sharing"])
	assert.Equal(suite.T(), callerID.String(), mediaSettings["screen_share_by"])

	// Stop screen sharing
	suite.stopScreenShare(callID, callerID)

	// Verify screen sharing is stopped
	callStatus = suite.getCallStatus(callID, callerID)
	mediaSettings = callStatus["media_settings"].(map[string]interface{})
	assert.Equal(suite.T(), false, mediaSettings["screen_sharing"])

	// End call
	suite.endCall(callID, callerID)
}

// Helper methods for video call testing

func (suite *VideoCallFlowIntegrationTestSuite) setUserPresence(userID uuid.UUID, status string) {
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

func (suite *VideoCallFlowIntegrationTestSuite) initiateCall(callerID, calleeID uuid.UUID, callType string) string {
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

func (suite *VideoCallFlowIntegrationTestSuite) getCallStatus(callID string, userID uuid.UUID) map[string]interface{} {
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

func (suite *VideoCallFlowIntegrationTestSuite) answerCall(callID string, userID uuid.UUID, accept bool) {
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

func (suite *VideoCallFlowIntegrationTestSuite) toggleMedia(callID string, userID uuid.UUID, mediaType string, enabled bool) {
	requestBody := map[string]interface{}{
		"call_id":    callID,
		"user_id":    userID.String(),
		"media_type": mediaType,
		"enabled":    enabled,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/media/toggle", bytes.NewBuffer(bodyBytes))
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

func (suite *VideoCallFlowIntegrationTestSuite) switchCamera(callID string, userID uuid.UUID, facing string) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
		"facing":  facing,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/camera/switch", bytes.NewBuffer(bodyBytes))
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

func (suite *VideoCallFlowIntegrationTestSuite) changeVideoQuality(callID string, userID uuid.UUID, quality string) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
		"quality": quality,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/video/quality", bytes.NewBuffer(bodyBytes))
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

func (suite *VideoCallFlowIntegrationTestSuite) downgradeCall(callID string, userID uuid.UUID, reason string) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
		"reason":  reason,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/downgrade", bytes.NewBuffer(bodyBytes))
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

func (suite *VideoCallFlowIntegrationTestSuite) startScreenShare(callID string, userID uuid.UUID) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/screen/start", bytes.NewBuffer(bodyBytes))
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

func (suite *VideoCallFlowIntegrationTestSuite) stopScreenShare(callID string, userID uuid.UUID) {
	requestBody := map[string]interface{}{
		"call_id": callID,
		"user_id": userID.String(),
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/calls/"+callID+"/screen/stop", bytes.NewBuffer(bodyBytes))
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

func (suite *VideoCallFlowIntegrationTestSuite) endCall(callID string, userID uuid.UUID) {
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

func (suite *VideoCallFlowIntegrationTestSuite) getCallHistory(userID uuid.UUID) []interface{} {
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

// TestVideoCallFlowIntegrationTestSuite runs the video call integration test suite
func TestVideoCallFlowIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(VideoCallFlowIntegrationTestSuite))
}