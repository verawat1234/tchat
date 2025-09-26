// Journey 2: Real-Time Messaging API Integration Tests
// Tests all API endpoints involved in real-time messaging functionality

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Journey02MessagingAPISuite struct {
	suite.Suite
	baseURL    string
	wsURL      string
	httpClient *http.Client
	ctx        context.Context
	user1      *AuthenticatedUser // Maya (Thailand)
	user2      *AuthenticatedUser // Li Wei (Singapore)
	wsConn1    *websocket.Conn
	wsConn2    *websocket.Conn
}

type AuthenticatedUser struct {
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Country      string `json:"country"`
	Language     string `json:"language"`
}

type CreateConversationRequest struct {
	Type         string                 `json:"type"`         // "direct", "group", "channel", "broadcast"
	Title        string                 `json:"title,omitempty"`
	Participants []string               `json:"participants"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type ConversationResponse struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Title        string                 `json:"title"`
	Participants []string               `json:"participants"`
	CreatedAt    string                 `json:"createdAt"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type SendMessageRequest struct {
	ConversationID string                 `json:"conversationId"`
	Type           string                 `json:"type"` // "text", "image", "file", "voice", "video", "location", "payment", "sticker"
	Content        string                 `json:"content"`
	ReplyTo        string                 `json:"replyTo,omitempty"`
	Attachments    []MessageAttachment    `json:"attachments,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type MessageAttachment struct {
	Type         string                 `json:"type"`
	URL          string                 `json:"url"`
	ThumbnailURL string                 `json:"thumbnailUrl,omitempty"`
	Filename     string                 `json:"filename,omitempty"`
	Size         int64                  `json:"size,omitempty"`
	Duration     int                    `json:"duration,omitempty"` // for audio/video
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type MessageResponse struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversationId"`
	SenderID       string                 `json:"senderId"`
	Type           string                 `json:"type"`
	Content        string                 `json:"content"`
	Attachments    []MessageAttachment    `json:"attachments,omitempty"`
	SentAt         string                 `json:"sentAt"`
	ReadBy         []string               `json:"readBy,omitempty"`
	Reactions      []MessageReaction      `json:"reactions,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type MessageReaction struct {
	UserID string `json:"userId"`
	Emoji  string `json:"emoji"`
	SentAt string `json:"sentAt"`
}

type WebSocketMessage struct {
	Type           string                 `json:"type"`
	MessageID      string                 `json:"messageId,omitempty"`
	ConversationID string                 `json:"conversationId,omitempty"`
	UserID         string                 `json:"userId,omitempty"`
	Content        string                 `json:"content,omitempty"`
	IsTyping       bool                   `json:"isTyping,omitempty"`
	ReadBy         string                 `json:"readBy,omitempty"`
	Emoji          string                 `json:"emoji,omitempty"`
	Timestamp      string                 `json:"timestamp,omitempty"`
}

type TypingIndicatorRequest struct {
	ConversationID string `json:"conversationId"`
	IsTyping       bool   `json:"isTyping"`
}

type ReadReceiptRequest struct {
	MessageID string    `json:"messageId"`
	ReadAt    time.Time `json:"readAt"`
}

type MessageReactionRequest struct {
	MessageID string `json:"messageId"`
	Emoji     string `json:"emoji"`
	Action    string `json:"action"` // "add" or "remove"
}

type FileUploadRequest struct {
	ConversationID string `json:"conversationId"`
	FileType       string `json:"fileType"`
	FileName       string `json:"fileName"`
	FileSize       int64  `json:"fileSize"`
}

type FileUploadResponse struct {
	FileURL      string `json:"fileUrl"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	FileID       string `json:"fileId"`
}

func (suite *Journey02MessagingAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081"   // Auth Service Direct
	suite.wsURL = "ws://localhost:8081/ws"    // WebSocket endpoint
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.ctx = context.Background()

	// Create authenticated test users
	suite.user1 = suite.createAuthenticatedUser("maya@test.com", "TH", "th")
	suite.user2 = suite.createAuthenticatedUser("liwei@test.com", "SG", "en")

	// Establish WebSocket connections
	suite.wsConn1 = suite.establishWebSocketConnection(suite.user1.AccessToken)
	suite.wsConn2 = suite.establishWebSocketConnection(suite.user2.AccessToken)
}

func (suite *Journey02MessagingAPISuite) TearDownSuite() {
	if suite.wsConn1 != nil {
		suite.wsConn1.Close()
	}
	if suite.wsConn2 != nil {
		suite.wsConn2.Close()
	}
}

// Test 2.1: Conversation Creation API
func (suite *Journey02MessagingAPISuite) TestConversationCreationAPI() {
	// Step 1: POST /api/v1/conversations - Create direct conversation
	convReq := CreateConversationRequest{
		Type:         "direct",
		Participants: []string{suite.user1.UserID, suite.user2.UserID},
		Metadata: map[string]interface{}{
			"initiatedBy": suite.user1.UserID,
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}

	convResp, statusCode := suite.makeAPICall("POST", "/api/v1/conversations", convReq, headers)
	assert.Equal(suite.T(), 201, statusCode, "Conversation creation should succeed")

	var conversation ConversationResponse
	err := json.Unmarshal(convResp, &conversation)
	require.NoError(suite.T(), err, "Should parse conversation response")

	assert.NotEmpty(suite.T(), conversation.ID, "Should return conversation ID")
	assert.Equal(suite.T(), "direct", conversation.Type, "Type should be direct")
	assert.Len(suite.T(), conversation.Participants, 2, "Should have 2 participants")

	conversationID := conversation.ID

	// Step 2: GET /api/v1/conversations/{id} - Retrieve conversation details
	getConvResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/conversations/%s", conversationID), nil, headers)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve conversation")

	var retrievedConv ConversationResponse
	err = json.Unmarshal(getConvResp, &retrievedConv)
	require.NoError(suite.T(), err, "Should parse retrieved conversation")
	assert.Equal(suite.T(), conversationID, retrievedConv.ID, "IDs should match")

	// Step 3: GET /api/v1/conversations/user/{userId} - List user conversations
	listConvResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/conversations/user/%s", suite.user1.UserID), nil, headers)
	assert.Equal(suite.T(), 200, statusCode, "Should list user conversations")

	var conversations []ConversationResponse
	err = json.Unmarshal(listConvResp, &conversations)
	require.NoError(suite.T(), err, "Should parse conversations list")
	assert.Greater(suite.T(), len(conversations), 0, "Should have conversations")

	// Find our conversation
	found := false
	for _, conv := range conversations {
		if conv.ID == conversationID {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "Should find created conversation in list")
}

// Test 2.2: Real-Time Message Sending API
func (suite *Journey02MessagingAPISuite) TestRealTimeMessagingAPI() {
	// Create conversation first
	conversationID := suite.createTestConversation()

	// Step 1: POST /api/v1/messages - Send message
	msgReq := SendMessageRequest{
		ConversationID: conversationID,
		Type:          "text",
		Content:       "Hello from Thailand! How's Singapore? 🇹🇭",
		Metadata: map[string]interface{}{
			"platform": "ios",
			"language": "en",
		},
	}

	headers1 := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}

	msgResp, statusCode := suite.makeAPICall("POST", "/api/v1/messages", msgReq, headers1)
	assert.Equal(suite.T(), 201, statusCode, "Message sending should succeed")

	var message MessageResponse
	err := json.Unmarshal(msgResp, &message)
	require.NoError(suite.T(), err, "Should parse message response")

	assert.NotEmpty(suite.T(), message.ID, "Should return message ID")
	assert.Equal(suite.T(), conversationID, message.ConversationID, "Conversation ID should match")
	assert.Equal(suite.T(), suite.user1.UserID, message.SenderID, "Sender ID should match")
	assert.Equal(suite.T(), "Hello from Thailand! How's Singapore? 🇹🇭", message.Content)

	messageID := message.ID

	// Step 2: Verify WebSocket message received by user2
	wsMessage := suite.readWebSocketMessage(suite.wsConn2, 5*time.Second)
	assert.Equal(suite.T(), "new_message", wsMessage.Type, "Should receive new_message event")
	assert.Equal(suite.T(), messageID, wsMessage.MessageID, "Message ID should match")
	assert.Equal(suite.T(), conversationID, wsMessage.ConversationID, "Conversation ID should match")

	// Step 3: GET /api/v1/conversations/{id}/messages - Retrieve messages
	headers2 := map[string]string{
		"Authorization": "Bearer " + suite.user2.AccessToken,
	}

	messagesResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/conversations/%s/messages", conversationID), nil, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve messages")

	var messages []MessageResponse
	err = json.Unmarshal(messagesResp, &messages)
	require.NoError(suite.T(), err, "Should parse messages list")
	assert.Len(suite.T(), messages, 1, "Should have 1 message")
	assert.Equal(suite.T(), messageID, messages[0].ID, "Message ID should match")

	// Step 4: POST /api/v1/messages/{id}/read - Mark message as read
	readReq := ReadReceiptRequest{
		MessageID: messageID,
		ReadAt:    time.Now(),
	}

	_, statusCode = suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/messages/%s/read", messageID), readReq, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Read receipt should succeed")

	// Step 5: Verify read receipt via WebSocket
	readReceiptMsg := suite.readWebSocketMessage(suite.wsConn1, 3*time.Second)
	assert.Equal(suite.T(), "read_receipt", readReceiptMsg.Type, "Should receive read receipt")
	assert.Equal(suite.T(), messageID, readReceiptMsg.MessageID, "Message ID should match")
	assert.Equal(suite.T(), suite.user2.UserID, readReceiptMsg.ReadBy, "Read by should match user2")
}

// Test 2.3: Typing Indicators API
func (suite *Journey02MessagingAPISuite) TestTypingIndicatorsAPI() {
	conversationID := suite.createTestConversation()

	// Step 1: POST /api/v1/conversations/{id}/typing - Start typing
	typingReq := TypingIndicatorRequest{
		ConversationID: conversationID,
		IsTyping:      true,
	}

	headers2 := map[string]string{
		"Authorization": "Bearer " + suite.user2.AccessToken,
	}

	_, statusCode := suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/conversations/%s/typing", conversationID), typingReq, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Typing indicator should succeed")

	// Step 2: Verify typing indicator via WebSocket
	typingMessage := suite.readWebSocketMessage(suite.wsConn1, 3*time.Second)
	assert.Equal(suite.T(), "typing_indicator", typingMessage.Type, "Should receive typing indicator")
	assert.Equal(suite.T(), suite.user2.UserID, typingMessage.UserID, "User ID should match")
	assert.True(suite.T(), typingMessage.IsTyping, "Should be typing")

	// Step 3: POST /api/v1/conversations/{id}/typing - Stop typing
	typingReq.IsTyping = false
	_, statusCode = suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/conversations/%s/typing", conversationID), typingReq, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Stop typing should succeed")

	// Step 4: Verify stop typing via WebSocket
	stopTypingMessage := suite.readWebSocketMessage(suite.wsConn1, 3*time.Second)
	assert.Equal(suite.T(), "typing_indicator", stopTypingMessage.Type, "Should receive typing indicator")
	assert.False(suite.T(), stopTypingMessage.IsTyping, "Should not be typing")
}

// Test 2.4: Message Reactions API
func (suite *Journey02MessagingAPISuite) TestMessageReactionsAPI() {
	conversationID := suite.createTestConversation()
	messageID := suite.sendTestMessage(conversationID, "React to this message! 😄")

	// Step 1: POST /api/v1/messages/{id}/reactions - Add reaction
	reactionReq := MessageReactionRequest{
		MessageID: messageID,
		Emoji:     "👍",
		Action:    "add",
	}

	headers2 := map[string]string{
		"Authorization": "Bearer " + suite.user2.AccessToken,
	}

	_, statusCode := suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/messages/%s/reactions", messageID), reactionReq, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Adding reaction should succeed")

	// Step 2: Verify reaction via WebSocket
	reactionMessage := suite.readWebSocketMessage(suite.wsConn1, 3*time.Second)
	assert.Equal(suite.T(), "message_reaction", reactionMessage.Type, "Should receive reaction event")
	assert.Equal(suite.T(), messageID, reactionMessage.MessageID, "Message ID should match")
	assert.Equal(suite.T(), "👍", reactionMessage.Emoji, "Emoji should match")
	assert.Equal(suite.T(), suite.user2.UserID, reactionMessage.UserID, "User ID should match")

	// Step 3: GET /api/v1/messages/{id} - Verify reaction stored
	messageResp, statusCode := suite.makeAPICall("GET",
		fmt.Sprintf("/api/v1/messages/%s", messageID), nil, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Should retrieve message with reactions")

	var messageWithReactions MessageResponse
	err := json.Unmarshal(messageResp, &messageWithReactions)
	require.NoError(suite.T(), err, "Should parse message with reactions")

	assert.Len(suite.T(), messageWithReactions.Reactions, 1, "Should have 1 reaction")
	assert.Equal(suite.T(), "👍", messageWithReactions.Reactions[0].Emoji, "Emoji should match")
	assert.Equal(suite.T(), suite.user2.UserID, messageWithReactions.Reactions[0].UserID, "User ID should match")

	// Step 4: POST /api/v1/messages/{id}/reactions - Remove reaction
	reactionReq.Action = "remove"
	_, statusCode = suite.makeAPICall("POST",
		fmt.Sprintf("/api/v1/messages/%s/reactions", messageID), reactionReq, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Removing reaction should succeed")
}

// Test 2.5: File Upload and Media Sharing API
func (suite *Journey02MessagingAPISuite) TestFileUploadAPI() {
	conversationID := suite.createTestConversation()

	// Step 1: POST /api/v1/files/upload - Upload file
	uploadReq := FileUploadRequest{
		ConversationID: conversationID,
		FileType:       "image/jpeg",
		FileName:       "test_image.jpg",
		FileSize:       1024000, // 1MB
	}

	headers1 := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
		"Content-Type":  "multipart/form-data",
	}

	// Create test image data
	testImageData := suite.generateTestImageData()
	uploadResp, statusCode := suite.uploadFile("/api/v1/files/upload", uploadReq, testImageData, headers1)
	assert.Equal(suite.T(), 200, statusCode, "File upload should succeed")

	var uploadResult FileUploadResponse
	err := json.Unmarshal(uploadResp, &uploadResult)
	require.NoError(suite.T(), err, "Should parse upload response")

	assert.NotEmpty(suite.T(), uploadResult.FileURL, "Should return file URL")
	assert.NotEmpty(suite.T(), uploadResult.FileID, "Should return file ID")

	// Step 2: POST /api/v1/messages - Send message with attachment
	msgReq := SendMessageRequest{
		ConversationID: conversationID,
		Type:          "image",
		Content:       "Check out this beautiful sunset in Bangkok! 🌅",
		Attachments: []MessageAttachment{
			{
				Type:         "image",
				URL:          uploadResult.FileURL,
				ThumbnailURL: uploadResult.ThumbnailURL,
				Filename:     uploadReq.FileName,
				Size:         uploadReq.FileSize,
				Metadata: map[string]interface{}{
					"width":  1920,
					"height": 1080,
				},
			},
		},
	}

	msgResp, statusCode := suite.makeAPICall("POST", "/api/v1/messages", msgReq, headers1)
	assert.Equal(suite.T(), 201, statusCode, "Message with attachment should succeed")

	var message MessageResponse
	err = json.Unmarshal(msgResp, &message)
	require.NoError(suite.T(), err, "Should parse message response")

	assert.Len(suite.T(), message.Attachments, 1, "Should have 1 attachment")
	assert.Equal(suite.T(), "image", message.Attachments[0].Type, "Attachment type should be image")
	assert.Equal(suite.T(), uploadResult.FileURL, message.Attachments[0].URL, "File URL should match")

	// Step 3: Verify media message via WebSocket
	wsMessage := suite.readWebSocketMessage(suite.wsConn2, 5*time.Second)
	assert.Equal(suite.T(), "new_message", wsMessage.Type, "Should receive new message event")
	assert.Contains(suite.T(), wsMessage.Content, "sunset", "Content should match")
}

// Test 2.6: Group Chat Management API
func (suite *Journey02MessagingAPISuite) TestGroupChatManagementAPI() {
	// Create third user for group testing
	user3 := suite.createAuthenticatedUser("arif@test.com", "ID", "id")
	wsConn3 := suite.establishWebSocketConnection(user3.AccessToken)
	defer wsConn3.Close()

	// Step 1: POST /api/v1/conversations - Create group chat
	groupReq := CreateConversationRequest{
		Type:  "group",
		Title: "SEA Tech Friends",
		Participants: []string{
			suite.user1.UserID, // Maya (TH)
			suite.user2.UserID, // Li Wei (SG)
			user3.UserID,       // Arif (ID)
		},
		Metadata: map[string]interface{}{
			"description": "Southeast Asian tech professionals",
			"avatar":      "https://example.com/group-avatar.jpg",
			"adminId":     suite.user2.UserID, // Li Wei is admin
		},
	}

	headers2 := map[string]string{
		"Authorization": "Bearer " + suite.user2.AccessToken,
	}

	groupResp, statusCode := suite.makeAPICall("POST", "/api/v1/conversations", groupReq, headers2)
	assert.Equal(suite.T(), 201, statusCode, "Group creation should succeed")

	var group ConversationResponse
	err := json.Unmarshal(groupResp, &group)
	require.NoError(suite.T(), err, "Should parse group response")

	groupID := group.ID
	assert.NotEmpty(suite.T(), groupID, "Should return group ID")
	assert.Equal(suite.T(), "group", group.Type, "Type should be group")
	assert.Len(suite.T(), group.Participants, 3, "Should have 3 participants")

	// Step 2: Verify all members receive group creation notification
	wsMessage1 := suite.readWebSocketMessage(suite.wsConn1, 5*time.Second)
	wsMessage2 := suite.readWebSocketMessage(suite.wsConn2, 5*time.Second)
	wsMessage3 := suite.readWebSocketMessage(wsConn3, 5*time.Second)

	for _, msg := range []*WebSocketMessage{&wsMessage1, &wsMessage2, &wsMessage3} {
		assert.Equal(suite.T(), "group_created", msg.Type, "Should receive group creation event")
		assert.Equal(suite.T(), groupID, msg.ConversationID, "Group ID should match")
	}

	// Step 3: Send group message
	groupMsgReq := SendMessageRequest{
		ConversationID: groupID,
		Type:          "text",
		Content:       "Welcome everyone to our tech group! 🚀 Let's share knowledge across SEA!",
		Metadata: map[string]interface{}{
			"platform": "web",
		},
	}

	_, statusCode = suite.makeAPICall("POST", "/api/v1/messages", groupMsgReq, headers2)
	assert.Equal(suite.T(), 201, statusCode, "Group message should succeed")

	// Step 4: Verify all members receive group message
	groupMsg1 := suite.readWebSocketMessage(suite.wsConn1, 5*time.Second) // Maya
	groupMsg3 := suite.readWebSocketMessage(wsConn3, 5*time.Second)       // Arif

	for _, msg := range []*WebSocketMessage{&groupMsg1, &groupMsg3} {
		assert.Equal(suite.T(), "new_message", msg.Type, "Should receive new message")
		assert.Equal(suite.T(), groupID, msg.ConversationID, "Group ID should match")
		assert.Contains(suite.T(), msg.Content, "Welcome", "Content should match")
	}

	// Step 5: PUT /api/v1/conversations/{id}/participants - Add participant
	addParticipantReq := map[string]interface{}{
		"action":     "add",
		"userIds":    []string{"new-user-id"},
		"addedBy":    suite.user2.UserID,
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/conversations/%s/participants", groupID), addParticipantReq, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Adding participant should succeed")

	// Step 6: PUT /api/v1/conversations/{id} - Update group info
	updateGroupReq := map[string]interface{}{
		"title":       "SEA Tech Friends - Updated",
		"description": "Southeast Asian tech professionals - Active community",
		"updatedBy":   suite.user2.UserID,
	}

	_, statusCode = suite.makeAPICall("PUT",
		fmt.Sprintf("/api/v1/conversations/%s", groupID), updateGroupReq, headers2)
	assert.Equal(suite.T(), 200, statusCode, "Updating group should succeed")
}

// Helper methods
func (suite *Journey02MessagingAPISuite) makeAPICall(method, endpoint string, body interface{}, headers map[string]string) ([]byte, int) {
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
			if key != "Content-Type" { // Don't override content type for multipart
				req.Header.Set(key, value)
			}
		}
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	return respBody, resp.StatusCode
}

func (suite *Journey02MessagingAPISuite) uploadFile(endpoint string, uploadReq FileUploadRequest, fileData []byte, headers map[string]string) ([]byte, int) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add form fields
	writer.WriteField("conversationId", uploadReq.ConversationID)
	writer.WriteField("fileType", uploadReq.FileType)
	writer.WriteField("fileName", uploadReq.FileName)

	// Add file
	part, err := writer.CreateFormFile("file", uploadReq.FileName)
	require.NoError(suite.T(), err)
	part.Write(fileData)

	writer.Close()

	req, err := http.NewRequestWithContext(suite.ctx, "POST", suite.baseURL+endpoint, &b)
	require.NoError(suite.T(), err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if headers != nil {
		for key, value := range headers {
			if key != "Content-Type" {
				req.Header.Set(key, value)
			}
		}
	}

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	return respBody, resp.StatusCode
}

func (suite *Journey02MessagingAPISuite) createAuthenticatedUser(email, country, language string) *AuthenticatedUser {
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

func (suite *Journey02MessagingAPISuite) establishWebSocketConnection(accessToken string) *websocket.Conn {
	dialer := websocket.DefaultDialer
	header := http.Header{}
	header.Add("Authorization", "Bearer "+accessToken)

	conn, resp, err := dialer.Dial(suite.wsURL, header)
	if err != nil {
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			suite.T().Fatalf("WebSocket connection failed: %v, response: %s", err, string(body))
		}
		suite.T().Fatalf("WebSocket connection failed: %v", err)
	}

	// Wait for connection confirmation
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, _, err = conn.ReadMessage()
	if err != nil {
		suite.T().Logf("Warning: Did not receive connection confirmation: %v", err)
	}
	conn.SetReadDeadline(time.Time{}) // Clear deadline

	return conn
}

func (suite *Journey02MessagingAPISuite) readWebSocketMessage(conn *websocket.Conn, timeout time.Duration) WebSocketMessage {
	conn.SetReadDeadline(time.Now().Add(timeout))
	defer conn.SetReadDeadline(time.Time{}) // Clear deadline

	_, message, err := conn.ReadMessage()
	require.NoError(suite.T(), err, "Should read WebSocket message")

	var wsMsg WebSocketMessage
	err = json.Unmarshal(message, &wsMsg)
	require.NoError(suite.T(), err, "Should parse WebSocket message")

	return wsMsg
}

func (suite *Journey02MessagingAPISuite) createTestConversation() string {
	convReq := CreateConversationRequest{
		Type:         "direct",
		Participants: []string{suite.user1.UserID, suite.user2.UserID},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}

	convResp, statusCode := suite.makeAPICall("POST", "/api/v1/conversations", convReq, headers)
	require.Equal(suite.T(), 201, statusCode)

	var conversation ConversationResponse
	err := json.Unmarshal(convResp, &conversation)
	require.NoError(suite.T(), err)

	return conversation.ID
}

func (suite *Journey02MessagingAPISuite) sendTestMessage(conversationID, content string) string {
	msgReq := SendMessageRequest{
		ConversationID: conversationID,
		Type:          "text",
		Content:       content,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + suite.user1.AccessToken,
	}

	msgResp, statusCode := suite.makeAPICall("POST", "/api/v1/messages", msgReq, headers)
	require.Equal(suite.T(), 201, statusCode)

	var message MessageResponse
	err := json.Unmarshal(msgResp, &message)
	require.NoError(suite.T(), err)

	return message.ID
}

func (suite *Journey02MessagingAPISuite) generateTestImageData() []byte {
	// Generate simple test image data (minimal JPEG header)
	return []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
}

func TestJourney02MessagingAPISuite(t *testing.T) {
	suite.Run(t, new(Journey02MessagingAPISuite))
}