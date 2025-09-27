// Journey 5: Cross-Platform Continuity API Integration Tests
// Tests all API endpoints involved in cross-platform state synchronization and continuity

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// AuthenticatedUser represents an authenticated user session
// Note: AuthenticatedUser is now defined in types.go

// LocationData represents location information for content or user activities
// Note: LocationData is now defined in types.go

// Coordinates represents geographic coordinates
// Note: Coordinates is now defined in types.go

type Journey05CrossPlatformAPISuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	ctx        context.Context
	user       *AuthenticatedUser // Li Wei (Singapore) - Cross-platform user
	webSession    *SessionContext
	mobileSession *SessionContext
}

type SessionContext struct {
	SessionID     string                 `json:"sessionId"`
	Platform      string                 `json:"platform"` // "web", "ios", "android"
	DeviceID      string                 `json:"deviceId"`
	UserAgent     string                 `json:"userAgent"`
	AppVersion    string                 `json:"appVersion,omitempty"`
	LastActive    string                 `json:"lastActive"`
	SyncState     map[string]interface{} `json:"syncState"`
	PushToken     string                 `json:"pushToken,omitempty"`
	Preferences   UserPreferences        `json:"preferences"`
}

type UserPreferences struct {
	Theme           string                 `json:"theme"` // "light", "dark", "auto"
	Language        string                 `json:"language"`
	Notifications   NotificationSettings   `json:"notifications"`
	Privacy         PrivacySettings        `json:"privacy"`
	Accessibility   AccessibilitySettings  `json:"accessibility"`
	CustomSettings  map[string]interface{} `json:"customSettings,omitempty"`
}

type NotificationSettings struct {
	Push        bool `json:"push"`
	Email       bool `json:"email"`
	SMS         bool `json:"sms"`
	InApp       bool `json:"inApp"`
	Sound       bool `json:"sound"`
	Vibration   bool `json:"vibration"`
	DND         bool `json:"doNotDisturb"`
	DNDStart    string `json:"dndStart,omitempty"` // "22:00"
	DNDEnd      string `json:"dndEnd,omitempty"`   // "07:00"
}

type PrivacySettings struct {
	ProfileVisibility  string `json:"profileVisibility"` // "public", "friends", "private"
	OnlineStatus       bool   `json:"onlineStatus"`
	ReadReceipts       bool   `json:"readReceipts"`
	TypingIndicators   bool   `json:"typingIndicators"`
	LastSeen           bool   `json:"lastSeen"`
	LocationSharing    bool   `json:"locationSharing"`
	DataCollection     bool   `json:"dataCollection"`
}

type AccessibilitySettings struct {
	HighContrast      bool    `json:"highContrast"`
	LargeText         bool    `json:"largeText"`
	ScreenReader      bool    `json:"screenReader"`
	VoiceCommands     bool    `json:"voiceCommands"`
	ReducedMotion     bool    `json:"reducedMotion"`
	TextToSpeech      bool    `json:"textToSpeech"`
	FontSize          float64 `json:"fontSize"` // Multiplier: 1.0 = normal, 1.5 = large
	ColorBlindSupport string  `json:"colorBlindSupport"` // "none", "protanopia", "deuteranopia", "tritanopia"
}

type CreateSessionRequest struct {
	Platform     string          `json:"platform"`
	DeviceID     string          `json:"deviceId"`
	UserAgent    string          `json:"userAgent"`
	AppVersion   string          `json:"appVersion,omitempty"`
	PushToken    string          `json:"pushToken,omitempty"`
	Location     LocationData    `json:"location,omitempty"`
	Preferences  UserPreferences `json:"preferences,omitempty"`
}

type DraftMessage struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversationId"`
	Content        string                 `json:"content"`
	Type           string                 `json:"type"` // "text", "voice", "media"
	Language       string                 `json:"language,omitempty"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
	Platform       string                 `json:"platform"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type SaveDraftRequest struct {
	ConversationID string                 `json:"conversationId"`
	Content        string                 `json:"content"`
	Type           string                 `json:"type,omitempty"`
	Language       string                 `json:"language,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type SyncStateRequest struct {
	StateData      map[string]interface{} `json:"stateData"`
	LastSyncTime   string                 `json:"lastSyncTime"`
	Platform       string                 `json:"platform"`
	ConflictPolicy string                 `json:"conflictPolicy"` // "client_wins", "server_wins", "merge", "prompt"
}

type SyncStateResponse struct {
	StateData       map[string]interface{} `json:"stateData"`
	LastSyncTime    string                 `json:"lastSyncTime"`
	ConflictsFound  bool                   `json:"conflictsFound"`
	ConflictDetails []SyncConflict         `json:"conflictDetails,omitempty"`
	SyncStatus      string                 `json:"syncStatus"` // "success", "partial", "failed"
}

type SyncConflict struct {
	Key          string      `json:"key"`
	ClientValue  interface{} `json:"clientValue"`
	ServerValue  interface{} `json:"serverValue"`
	Resolution   string      `json:"resolution"` // "client", "server", "merged"
	ResolvedBy   string      `json:"resolvedBy"` // "system", "user"
}

type OfflineAction struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // "send_message", "update_profile", "like_content", etc.
	Action      string                 `json:"action"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   string                 `json:"timestamp"`
	RetryCount  int                    `json:"retryCount"`
	MaxRetries  int                    `json:"maxRetries"`
	Status      string                 `json:"status"` // "pending", "processing", "completed", "failed"
}

type QueueOfflineActionRequest struct {
	Type        string                 `json:"type"`
	Action      string                 `json:"action"`
	Data        map[string]interface{} `json:"data"`
	MaxRetries  int                    `json:"maxRetries,omitempty"`
}

type CrossPlatformSyncRequest struct {
	TargetPlatforms []string               `json:"targetPlatforms"` // ["web", "ios", "android"]
	SyncData        map[string]interface{} `json:"syncData"`
	Priority        string                 `json:"priority"` // "high", "normal", "low"
	Immediate       bool                   `json:"immediate"`
}

func (suite *Journey05CrossPlatformAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081" // API Gateway
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.ctx = context.Background()

	// Create authenticated test user
	suite.user = suite.createAuthenticatedUser("liwei@test.com", "SG", "en")

	// Initialize session contexts
	suite.webSession = &SessionContext{
		Platform:   "web",
		DeviceID:   "web-device-12345",
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		AppVersion: "1.0.0",
	}

	suite.mobileSession = &SessionContext{
		Platform:   "ios",
		DeviceID:   "ios-device-67890",
		UserAgent:  "TchatApp/1.0.0 (iPhone; iOS 17.0)",
		AppVersion: "1.0.0",
		PushToken:  "apns-token-abc123def456",
	}
}

// Test 5.1: Session Management API
func (suite *Journey05CrossPlatformAPISuite) TestSessionManagementAPI() {
	userHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
	}

	// Step 1: POST /api/v1/sessions - Create web session
	webSessionReq := CreateSessionRequest{
		Platform:   suite.webSession.Platform,
		DeviceID:   suite.webSession.DeviceID,
		UserAgent:  suite.webSession.UserAgent,
		AppVersion: suite.webSession.AppVersion,
		Location: LocationData{
			Country: "SG",
			City:    "Singapore",
			Coordinates: Coordinates{
				Latitude:  1.3521,
				Longitude: 103.8198,
			},
		},
		Preferences: UserPreferences{
			Theme:    "dark",
			Language: "en",
			Notifications: NotificationSettings{
				Push:      true,
				Email:     true,
				SMS:       false,
				InApp:     true,
				Sound:     true,
				Vibration: false,
			},
			Privacy: PrivacySettings{
				ProfileVisibility: "public",
				OnlineStatus:      true,
				ReadReceipts:      true,
				TypingIndicators:  true,
				LastSeen:          true,
			},
			Accessibility: AccessibilitySettings{
				HighContrast:  false,
				LargeText:     false,
				ScreenReader:  false,
				FontSize:      1.0,
			},
		},
	}

	webSessionResp, statusCode := suite.makeAPICall("POST", "/api/v1/sessions", webSessionReq, userHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Web session creation should succeed")

	var webSession SessionContext
	err := json.Unmarshal(webSessionResp, &webSession)
	require.NoError(suite.T(), err, "Should parse web session response")

	assert.NotEmpty(suite.T(), webSession.SessionID, "Should return session ID")
	assert.Equal(suite.T(), "web", webSession.Platform, "Platform should match")
	assert.Equal(suite.T(), "dark", webSession.Preferences.Theme, "Theme should match")

	suite.webSession.SessionID = webSession.SessionID

	// Step 2: POST /api/v1/sessions - Create mobile session
	mobileSessionReq := CreateSessionRequest{
		Platform:   suite.mobileSession.Platform,
		DeviceID:   suite.mobileSession.DeviceID,
		UserAgent:  suite.mobileSession.UserAgent,
		AppVersion: suite.mobileSession.AppVersion,
		PushToken:  suite.mobileSession.PushToken,
		Preferences: UserPreferences{
			Theme:    "auto",
			Language: "en",
			Notifications: NotificationSettings{
				Push:      true,
				Email:     false,
				SMS:       true,
				InApp:     true,
				Sound:     true,
				Vibration: true,
			},
		},
	}

	mobileSessionResp, statusCode := suite.makeAPICall("POST", "/api/v1/sessions", mobileSessionReq, userHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Mobile session creation should succeed")

	var mobileSession SessionContext
	err = json.Unmarshal(mobileSessionResp, &mobileSession)
	require.NoError(suite.T(), err, "Should parse mobile session response")

	suite.mobileSession.SessionID = mobileSession.SessionID

	// Step 3: GET /api/v1/sessions - List user sessions
	listSessionsResp, statusCode := suite.makeAPICall("GET", "/api/v1/sessions", nil, userHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should list user sessions")

	var sessions []SessionContext
	err = json.Unmarshal(listSessionsResp, &sessions)
	require.NoError(suite.T(), err, "Should parse sessions list")

	assert.Len(suite.T(), sessions, 2, "Should have 2 active sessions")

	// Find both sessions
	var foundWeb, foundMobile bool
	for _, session := range sessions {
		if session.Platform == "web" && session.SessionID == suite.webSession.SessionID {
			foundWeb = true
		} else if session.Platform == "ios" && session.SessionID == suite.mobileSession.SessionID {
			foundMobile = true
		}
	}
	assert.True(suite.T(), foundWeb, "Should find web session")
	assert.True(suite.T(), foundMobile, "Should find mobile session")

	// Step 4: PUT /api/v1/sessions/{id} - Update session
	updateSessionReq := map[string]interface{}{
		"preferences": map[string]interface{}{
			"theme": "light",
			"notifications": map[string]interface{}{
				"push":  false,
				"email": false,
			},
		},
		"lastActive": time.Now().Format(time.RFC3339),
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/sessions/%s", suite.webSession.SessionID), updateSessionReq, userHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Session update should succeed")

	// Step 5: GET /api/v1/sessions/{id} - Verify update
	getUpdatedResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/sessions/%s", suite.webSession.SessionID), nil, userHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve updated session")

	var updatedSession SessionContext
	err = json.Unmarshal(getUpdatedResp, &updatedSession)
	require.NoError(suite.T(), err, "Should parse updated session")

	assert.Equal(suite.T(), "light", updatedSession.Preferences.Theme, "Theme should be updated")
	assert.False(suite.T(), updatedSession.Preferences.Notifications.Push, "Push notifications should be disabled")
}

// Test 5.2: Draft Message Synchronization API
func (suite *Journey05CrossPlatformAPISuite) TestDraftSynchronizationAPI() {
	// Create test conversation first
	conversationID := suite.createTestConversation()

	userHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.webSession.SessionID,
	}

	// Step 1: POST /api/v1/drafts - Save draft on web
	draftReq := SaveDraftRequest{
		ConversationID: conversationID,
		Content:        "I was just reviewing the project documentation and I think we should consider implementing",
		Type:          "text",
		Language:      "en",
		Metadata: map[string]interface{}{
			"platform":     "web",
			"cursorPos":    97,
			"wordCount":    15,
		},
	}

	draftResp, statusCode := suite.makeAPICall("POST", "/api/v1/drafts", draftReq, userHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Draft saving should succeed")

	var draft DraftMessage
	err := json.Unmarshal(draftResp, &draft)
	require.NoError(suite.T(), err, "Should parse draft response")

	draftID := draft.ID
	assert.NotEmpty(suite.T(), draftID, "Should return draft ID")
	assert.Equal(suite.T(), conversationID, draft.ConversationID, "Conversation ID should match")
	assert.Contains(suite.T(), draft.Content, "project documentation", "Content should match")

	// Step 2: Switch to mobile session and retrieve drafts
	mobileHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.mobileSession.SessionID,
		"X-Platform":    "ios",
	}

	getDraftsResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/conversations/%s/drafts", conversationID), nil, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve drafts")

	var drafts []DraftMessage
	err = json.Unmarshal(getDraftsResp, &drafts)
	require.NoError(suite.T(), err, "Should parse drafts list")

	assert.Len(suite.T(), drafts, 1, "Should have 1 draft")
	assert.Equal(suite.T(), draftID, drafts[0].ID, "Draft ID should match")
	assert.Contains(suite.T(), drafts[0].Content, "project documentation", "Content should sync across platforms")

	// Step 3: PUT /api/v1/drafts/{id} - Continue draft on mobile
	continueDraftReq := map[string]interface{}{
		"content": "I was just reviewing the project documentation and I think we should consider implementing the new authentication flow discussed yesterday. This would improve security significantly.",
		"metadata": map[string]interface{}{
			"platform":     "ios",
			"completed":    true,
			"originalFrom": "web",
		},
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/drafts/%s", draftID), continueDraftReq, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Draft continuation should succeed")

	// Step 4: Switch back to web and verify sync
	time.Sleep(1 * time.Second) // Allow sync time

	webSyncHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.webSession.SessionID,
	}

	syncedDraftResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/drafts/%s", draftID), nil, webSyncHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve synced draft")

	var syncedDraft DraftMessage
	err = json.Unmarshal(syncedDraftResp, &syncedDraft)
	require.NoError(suite.T(), err, "Should parse synced draft")

	assert.Contains(suite.T(), syncedDraft.Content, "authentication flow", "Should contain continued content")
	assert.Contains(suite.T(), syncedDraft.Content, "improve security", "Should contain mobile additions")
	assert.Equal(suite.T(), "ios", syncedDraft.Metadata["originalFrom"], "Should track platform origin")

	// Step 5: POST /api/v1/messages - Send completed draft from web
	sendMessageReq := map[string]interface{}{
		"conversationId": conversationID,
		"type":          "text",
		"content":       syncedDraft.Content,
		"metadata": map[string]interface{}{
			"fromDraft":  true,
			"draftId":    draftID,
			"platforms":  []string{"web", "ios"},
		},
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/messages", sendMessageReq, webSyncHeaders)
	assert.Equal(suite.T(), 201, statusCode, "Sending draft as message should succeed")

	// Step 6: Verify draft is deleted after sending
	_, statusCode = suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/drafts/%s", draftID), nil, webSyncHeaders)
	assert.Equal(suite.T(), 404, statusCode, "Draft should be deleted after sending")
}

// Test 5.3: State Synchronization API
func (suite *Journey05CrossPlatformAPISuite) TestStateSynchronizationAPI() {
	userHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.webSession.SessionID,
	}

	// Step 1: POST /api/v1/sync/state - Save state from web
	webStateReq := SyncStateRequest{
		StateData: map[string]interface{}{
			"activeConversations": []string{"conv1", "conv2", "conv3"},
			"unreadCounts": map[string]interface{}{
				"conv1": 5,
				"conv2": 0,
				"conv3": 12,
			},
			"currentTab": "chat",
			"scrollPositions": map[string]interface{}{
				"conv1": 245,
				"conv2": 0,
				"conv3": 890,
			},
			"preferences": map[string]interface{}{
				"chatBubbleStyle": "modern",
				"compactMode":     false,
				"mediaAutoplay":   true,
			},
			"filters": map[string]interface{}{
				"showOnlineOnly": false,
				"mutedChats":     []string{"conv4", "conv5"},
			},
		},
		LastSyncTime:   time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		Platform:       "web",
		ConflictPolicy: "client_wins",
	}

	webSyncResp, statusCode := suite.makeAPICall("POST", "/api/v1/sync/state", webStateReq, userHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Web state sync should succeed")

	var webSyncResult SyncStateResponse
	err := json.Unmarshal(webSyncResp, &webSyncResult)
	require.NoError(suite.T(), err, "Should parse web sync response")

	assert.Equal(suite.T(), "success", webSyncResult.SyncStatus, "Sync status should be success")
	assert.False(suite.T(), webSyncResult.ConflictsFound, "Should not have conflicts initially")

	// Step 2: Switch to mobile and sync state
	mobileHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.mobileSession.SessionID,
	}

	mobileStateReq := SyncStateRequest{
		StateData: map[string]interface{}{
			"activeConversations": []string{"conv1", "conv2", "conv3", "conv6"}, // Added conv6
			"unreadCounts": map[string]interface{}{
				"conv1": 5,    // Same
				"conv2": 2,    // Different (conflict)
				"conv3": 12,   // Same
				"conv6": 3,    // New
			},
			"currentTab": "store", // Different
			"preferences": map[string]interface{}{
				"chatBubbleStyle": "classic",  // Different (conflict)
				"compactMode":     true,       // Different (conflict)
				"mediaAutoplay":   true,       // Same
				"vibrationEnabled": true,      // Mobile-only
			},
			"mobileSpecific": map[string]interface{}{
				"pushEnabled":   true,
				"batteryMode":   "normal",
				"biometricAuth": true,
			},
		},
		LastSyncTime:   time.Now().Add(-2 * time.Minute).Format(time.RFC3339),
		Platform:       "ios",
		ConflictPolicy: "merge",
	}

	mobileSyncResp, statusCode := suite.makeAPICall("POST", "/api/v1/sync/state", mobileStateReq, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Mobile state sync should succeed")

	var mobileSyncResult SyncStateResponse
	err = json.Unmarshal(mobileSyncResp, &mobileSyncResult)
	require.NoError(suite.T(), err, "Should parse mobile sync response")

	assert.Equal(suite.T(), "partial", mobileSyncResult.SyncStatus, "Sync status should be partial due to conflicts")
	assert.True(suite.T(), mobileSyncResult.ConflictsFound, "Should detect conflicts")
	assert.Greater(suite.T(), len(mobileSyncResult.ConflictDetails), 0, "Should have conflict details")

	// Verify specific conflicts
	conflictKeys := make([]string, len(mobileSyncResult.ConflictDetails))
	for i, conflict := range mobileSyncResult.ConflictDetails {
		conflictKeys[i] = conflict.Key
	}
	assert.Contains(suite.T(), conflictKeys, "unreadCounts.conv2", "Should detect unread count conflict")
	assert.Contains(suite.T(), conflictKeys, "preferences.chatBubbleStyle", "Should detect style conflict")

	// Step 3: GET /api/v1/sync/state - Retrieve synchronized state
	getStateResp, statusCode := suite.makeAPICall("GET", "/api/v1/sync/state", nil, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve synchronized state")

	var currentState SyncStateResponse
	err = json.Unmarshal(getStateResp, &currentState)
	require.NoError(suite.T(), err, "Should parse current state")

	// Verify merged state
	activeConvs := currentState.StateData["activeConversations"].([]interface{})
	assert.Len(suite.T(), activeConvs, 4, "Should have 4 active conversations (merged)")

	unreadCounts := currentState.StateData["unreadCounts"].(map[string]interface{})
	assert.Equal(suite.T(), float64(2), unreadCounts["conv2"], "Should use mobile value for conv2 (newer)")
	assert.Equal(suite.T(), float64(3), unreadCounts["conv6"], "Should include new conversation from mobile")

	// Step 4: Switch back to web and verify sync
	webSyncHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.webSession.SessionID,
	}

	webGetStateResp, statusCode := suite.makeAPICall("GET", "/api/v1/sync/state", nil, webSyncHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve state on web")

	var webCurrentState SyncStateResponse
	err = json.Unmarshal(webGetStateResp, &webCurrentState)
	require.NoError(suite.T(), err, "Should parse web current state")

	// Verify state is synchronized across platforms
	webActiveConvs := webCurrentState.StateData["activeConversations"].([]interface{})
	assert.Len(suite.T(), webActiveConvs, 4, "Web should have synchronized conversation list")

	webUnreadCounts := webCurrentState.StateData["unreadCounts"].(map[string]interface{})
	assert.Equal(suite.T(), float64(2), webUnreadCounts["conv2"], "Web should have updated unread count")
	assert.Equal(suite.T(), float64(3), webUnreadCounts["conv6"], "Web should have new conversation")
}

// Test 5.4: Offline Actions and Queue Management API
func (suite *Journey05CrossPlatformAPISuite) TestOfflineActionsAPI() {
	conversationID := suite.createTestConversation()

	mobileHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.mobileSession.SessionID,
		"X-Platform":    "ios",
	}

	// Step 1: POST /api/v1/offline/actions - Queue offline actions
	offlineActions := []QueueOfflineActionRequest{
		{
			Type:   "send_message",
			Action: "POST /api/v1/messages",
			Data: map[string]interface{}{
				"conversationId": conversationID,
				"type":          "text",
				"content":       "This message was queued while offline",
				"metadata": map[string]interface{}{
					"offline":   true,
					"timestamp": time.Now().Format(time.RFC3339),
				},
			},
			MaxRetries: 3,
		},
		{
			Type:   "update_profile",
			Action: "PUT /api/v1/profile",
			Data: map[string]interface{}{
				"status":      "Available",
				"customStatus": "Working from Singapore office",
			},
			MaxRetries: 2,
		},
		{
			Type:   "mark_read",
			Action: "POST /api/v1/messages/{messageId}/read",
			Data: map[string]interface{}{
				"messageId": "msg123",
				"readAt":    time.Now().Format(time.RFC3339),
			},
			MaxRetries: 1,
		},
	}

	// Queue each action
	queuedActionIDs := make([]string, len(offlineActions))
	for i, action := range offlineActions {
		queueResp, statusCode := suite.makeAPICall("POST", "/api/v1/offline/actions", action, mobileHeaders)
		assert.Equal(suite.T(), 201, statusCode, fmt.Sprintf("Queuing action %d should succeed", i+1))

		var queuedAction OfflineAction
		err := json.Unmarshal(queueResp, &queuedAction)
		require.NoError(suite.T(), err, "Should parse queued action response")

		queuedActionIDs[i] = queuedAction.ID
		assert.Equal(suite.T(), "pending", queuedAction.Status, "Status should be pending")
		assert.Equal(suite.T(), action.Type, queuedAction.Type, "Type should match")
	}

	// Step 2: GET /api/v1/offline/actions - List queued actions
	listActionsResp, statusCode := suite.makeAPICall("GET", "/api/v1/offline/actions", nil, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should list offline actions")

	var actionsList []OfflineAction
	err := json.Unmarshal(listActionsResp, &actionsList)
	require.NoError(suite.T(), err, "Should parse actions list")

	assert.Len(suite.T(), actionsList, 3, "Should have 3 queued actions")

	// Step 3: POST /api/v1/offline/actions/process - Process queued actions (simulate coming online)
	processReq := map[string]interface{}{
		"batchSize": 10,
		"priority":  "normal",
	}

	processResp, statusCode := suite.makeAPICall("POST", "/api/v1/offline/actions/process", processReq, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Action processing should succeed")

	var processResult map[string]interface{}
	err = json.Unmarshal(processResp, &processResult)
	require.NoError(suite.T(), err, "Should parse process result")

	assert.Equal(suite.T(), float64(3), processResult["processed"], "Should process 3 actions")
	assert.Equal(suite.T(), float64(0), processResult["failed"], "Should have no failed actions")

	// Step 4: Wait for processing and verify action status
	time.Sleep(3 * time.Second) // Allow processing time

	processedActionsResp, statusCode := suite.makeAPICall("GET", "/api/v1/offline/actions", nil, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve processed actions")

	var processedActions []OfflineAction
	err = json.Unmarshal(processedActionsResp, &processedActions)
	require.NoError(suite.T(), err, "Should parse processed actions")

	// Verify all actions are completed
	for _, action := range processedActions {
		assert.Equal(suite.T(), "completed", action.Status, fmt.Sprintf("Action %s should be completed", action.ID))
	}

	// Step 5: Verify the message was actually sent
	messagesResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/conversations/%s/messages", conversationID), nil, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve conversation messages")

	var messages []map[string]interface{}
	err = json.Unmarshal(messagesResp, &messages)
	require.NoError(suite.T(), err, "Should parse messages")

	// Find the offline message
	found := false
	for _, message := range messages {
		if content, ok := message["content"].(string); ok &&
		   strings.Contains(content, "queued while offline") {
			found = true
			metadata := message["metadata"].(map[string]interface{})
			assert.True(suite.T(), metadata["offline"].(bool), "Should have offline metadata")
			break
		}
	}
	assert.True(suite.T(), found, "Should find the offline-queued message")

	// Step 6: DELETE /api/v1/offline/actions/{id} - Clean up completed actions
	for _, actionID := range queuedActionIDs {
		_, statusCode := suite.makeAPICall("DELETE",
			fmt.Sprintf("/api/v1/offline/actions/%s", actionID), nil, mobileHeaders)
		assert.Equal(suite.T(), 204, statusCode, "Action deletion should succeed")
	}

	// Verify actions are deleted
	finalListResp, statusCode := suite.makeAPICall("GET", "/api/v1/offline/actions", nil, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve final actions list")

	var finalActions []OfflineAction
	err = json.Unmarshal(finalListResp, &finalActions)
	require.NoError(suite.T(), err, "Should parse final actions")
	assert.Len(suite.T(), finalActions, 0, "Should have no actions after cleanup")
}

// Test 5.5: Cross-Platform Real-Time Sync API
func (suite *Journey05CrossPlatformAPISuite) TestCrossPlatformRealTimeSyncAPI() {
	conversationID := suite.createTestConversation()

	webHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.webSession.SessionID,
	}

	mobileHeaders := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
		"X-Session-ID":  suite.mobileSession.SessionID,
	}

	// Step 1: POST /api/v1/sync/cross-platform - Initiate cross-platform sync
	crossPlatformReq := CrossPlatformSyncRequest{
		TargetPlatforms: []string{"web", "ios"},
		SyncData: map[string]interface{}{
			"action": "conversation_opened",
			"conversationId": conversationID,
			"timestamp": time.Now().Format(time.RFC3339),
			"metadata": map[string]interface{}{
				"source":   "ios",
				"priority": "high",
			},
		},
		Priority:  "high",
		Immediate: true,
	}

	syncResp, statusCode := suite.makeAPICall("POST", "/api/v1/sync/cross-platform", crossPlatformReq, mobileHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Cross-platform sync should succeed")

	var syncResult map[string]interface{}
	err := json.Unmarshal(syncResp, &syncResult)
	require.NoError(suite.T(), err, "Should parse sync result")

	syncID := syncResult["syncId"].(string)
	assert.NotEmpty(suite.T(), syncID, "Should return sync ID")
	assert.Equal(suite.T(), "initiated", syncResult["status"], "Status should be initiated")

	// Step 2: GET /api/v1/sync/cross-platform/{id} - Check sync status
	time.Sleep(1 * time.Second) // Allow sync propagation

	checkSyncResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/sync/cross-platform/%s", syncID), nil, webHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should check sync status")

	var syncStatus map[string]interface{}
	err = json.Unmarshal(checkSyncResp, &syncStatus)
	require.NoError(suite.T(), err, "Should parse sync status")

	assert.Equal(suite.T(), "completed", syncStatus["status"], "Sync should be completed")
	assert.Equal(suite.T(), float64(2), syncStatus["targetCount"], "Should have 2 target platforms")
	assert.Equal(suite.T(), float64(2), syncStatus["successCount"], "Should have 2 successful syncs")

	// Step 3: Verify sync data received on web platform
	getSyncDataResp, statusCode := suite.makeAPICall("GET", "/api/v1/sync/pending", nil, webHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve pending sync data")

	var pendingSyncs []map[string]interface{}
	err = json.Unmarshal(getSyncDataResp, &pendingSyncs)
	require.NoError(suite.T(), err, "Should parse pending syncs")

	// Find our sync data
	var foundSync map[string]interface{}
	for _, sync := range pendingSyncs {
		if sync["syncId"] == syncID {
			foundSync = sync
			break
		}
	}

	require.NotNil(suite.T(), foundSync, "Should find our sync data")
	syncData := foundSync["data"].(map[string]interface{})
	assert.Equal(suite.T(), "conversation_opened", syncData["action"], "Action should match")
	assert.Equal(suite.T(), conversationID, syncData["conversationId"], "Conversation ID should match")

	// Step 4: POST /api/v1/sync/acknowledge - Acknowledge sync receipt
	ackReq := map[string]interface{}{
		"syncId": syncID,
		"status": "processed",
		"platform": "web",
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/sync/acknowledge", ackReq, webHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Sync acknowledgment should succeed")

	// Step 5: POST /api/v1/sync/preferences - Update sync preferences
	prefsReq := map[string]interface{}{
		"autoSync": map[string]interface{}{
			"enabled": true,
			"interval": 30, // seconds
			"priority": "normal",
		},
		"syncScope": []string{"conversations", "drafts", "preferences", "state"},
		"conflictResolution": map[string]interface{}{
			"default": "server_wins",
			"preferences": "client_wins",
			"drafts": "merge",
		},
		"platforms": map[string]interface{}{
			"web": map[string]interface{}{
				"enabled": true,
				"priority": "high",
			},
			"ios": map[string]interface{}{
				"enabled": true,
				"priority": "high",
			},
		},
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/sync/preferences", prefsReq, webHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Sync preferences update should succeed")

	// Step 6: GET /api/v1/sync/status - Get overall sync status
	statusResp, statusCode := suite.makeAPICall("GET", "/api/v1/sync/status", nil, webHeaders)
	assert.Equal(suite.T(), 200, statusCode, "Should get sync status")

	var overallStatus map[string]interface{}
	err = json.Unmarshal(statusResp, &overallStatus)
	require.NoError(suite.T(), err, "Should parse overall status")

	assert.True(suite.T(), overallStatus["enabled"].(bool), "Sync should be enabled")
	assert.Greater(suite.T(), overallStatus["lastSync"], "", "Should have last sync timestamp")

	activeSessions := overallStatus["activeSessions"].([]interface{})
	assert.Len(suite.T(), activeSessions, 2, "Should show 2 active sessions")
}

// Helper methods
func (suite *Journey05CrossPlatformAPISuite) makeAPICall(method, endpoint string, body interface{}, headers map[string]string) ([]byte, int) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(suite.ctx, method, suite.baseURL+endpoint, reqBody)
	require.NoError(suite.T(), err)

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	return respBody, resp.StatusCode
}

func (suite *Journey05CrossPlatformAPISuite) createAuthenticatedUser(email, country, language string) *AuthenticatedUser {
	regReq := map[string]interface{}{
		"email":     email,
		"password":  "SecurePass123!",
		"firstName": "Test",
		"lastName":  "User",
		"country":   country,
		"language":  language,
	}

	regResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/register", regReq, nil)
	require.Equal(suite.T(), 201, statusCode)

	var regResult map[string]interface{}
	err := json.Unmarshal(regResp, &regResult)
	require.NoError(suite.T(), err)

	verifyReq := map[string]string{
		"userId": regResult["userId"].(string),
		"code":   regResult["verifyCode"].(string),
	}

	verifyResp, statusCode := suite.makeAPICall("POST", "/api/v1/auth/verify", verifyReq, nil)
	require.Equal(suite.T(), 200, statusCode)

	var verifyResult map[string]interface{}
	err = json.Unmarshal(verifyResp, &verifyResult)
	require.NoError(suite.T(), err)

	return &AuthenticatedUser{
		UserID:       regResult["userId"].(string),
		Email:        email,
		AccessToken:  verifyResult["accessToken"].(string),
		RefreshToken: verifyResult["refreshToken"].(string),
		Country:      country,
		Language:     language,
	}
}

func (suite *Journey05CrossPlatformAPISuite) createTestConversation() string {
	// Create a second user for conversation
	user2 := suite.createAuthenticatedUser("test.user2@test.com", "SG", "en")

	convReq := map[string]interface{}{
		"type":         "direct",
		"participants": []string{suite.user.UserID, user2.UserID},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + suite.user.AccessToken,
	}

	convResp, statusCode := suite.makeAPICall("POST", "/api/v1/conversations", convReq, headers)
	require.Equal(suite.T(), 201, statusCode)

	var conversation map[string]interface{}
	err := json.Unmarshal(convResp, &conversation)
	require.NoError(suite.T(), err)

	return conversation["id"].(string)
}

func TestJourney05CrossPlatformAPISuite(t *testing.T) {
	suite.Run(t, new(Journey05CrossPlatformAPISuite))
}
