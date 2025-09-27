package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// MessagingIntegrationSuite tests the Messaging service endpoints
type MessagingIntegrationSuite struct {
	APIIntegrationSuite
	ports ServicePort
}

// Dialog represents a chat dialog/conversation
type Dialog struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Type          string               `json:"type"`
	Description   string               `json:"description"`
	CreatedBy     string               `json:"createdBy"`
	Participants  []DialogParticipant  `json:"participants"`
	Settings      map[string]string    `json:"settings"`
	LastMessage   *Message             `json:"lastMessage,omitempty"`
	MessageCount  int                  `json:"messageCount"`
	UnreadCount   int                  `json:"unreadCount"`
	IsArchived    bool                 `json:"isArchived"`
	IsMuted       bool                 `json:"isMuted"`
	CreatedAt     string               `json:"createdAt"`
	UpdatedAt     string               `json:"updatedAt"`
}

// DialogParticipant represents a participant in a dialog
type DialogParticipant struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	DialogID string `json:"dialogId"`
	Role     string `json:"role"`
	JoinedAt string `json:"joinedAt"`
	LeftAt   *string `json:"leftAt,omitempty"`
	IsActive bool   `json:"isActive"`
}

// Message represents a chat message
type Message struct {
	ID        string            `json:"id"`
	DialogID  string            `json:"dialogId"`
	SenderID  string            `json:"senderId"`
	Type      string            `json:"type"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata"`
	ReplyToID *string           `json:"replyToId,omitempty"`
	EditedAt  *string           `json:"editedAt,omitempty"`
	Status    string            `json:"status"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
}

// CreateDialogRequest represents dialog creation request
type CreateDialogRequest struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Description  string            `json:"description"`
	Participants []string          `json:"participants"`
	Settings     map[string]string `json:"settings"`
}

// SendMessageRequest represents message sending request
type SendMessageRequest struct {
	Type      string            `json:"type"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata"`
	ReplyToID *string           `json:"replyToId,omitempty"`
}

// UpdateMessageRequest represents message update request
type UpdateMessageRequest struct {
	Content  *string           `json:"content,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// MessagingResponse represents messaging API response
type MessagingResponse struct {
	Success     bool        `json:"success"`
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Dialog      *Dialog     `json:"dialog,omitempty"`
	Dialogs     []Dialog    `json:"dialogs,omitempty"`
	ChatMessage *Message    `json:"chatMessage,omitempty"`
	Messages    []Message   `json:"messages,omitempty"`
	Total       int         `json:"total,omitempty"`
	Page        int         `json:"page,omitempty"`
	Limit       int         `json:"limit,omitempty"`
	Timestamp   string      `json:"timestamp"`
}

// SetupSuite initializes the messaging integration test suite
func (suite *MessagingIntegrationSuite) SetupSuite() {
	suite.APIIntegrationSuite.SetupSuite()
	suite.ports = DefaultServicePorts()

	// Wait for messaging service to be available
	err := suite.waitForService(suite.ports.Messaging, 30*time.Second)
	if err != nil {
		suite.T().Fatalf("Messaging service not available: %v", err)
	}
}

// TestMessagingServiceHealth verifies messaging service health endpoint
func (suite *MessagingIntegrationSuite) TestMessagingServiceHealth() {
	healthCheck, err := suite.checkServiceHealth(suite.ports.Messaging)
	require.NoError(suite.T(), err, "Health check should succeed")

	assert.Equal(suite.T(), "healthy", healthCheck.Status)
	assert.Equal(suite.T(), "messaging-service", healthCheck.Service)
	assert.NotEmpty(suite.T(), healthCheck.Timestamp)
}

// TestCreateDialog tests dialog creation
func (suite *MessagingIntegrationSuite) TestCreateDialog() {
	url := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)

	createReq := CreateDialogRequest{
		Name:        "Test Group Chat",
		Type:        "group",
		Description: "A test group chat for integration testing",
		Participants: []string{"user1", "user2", "user3"},
		Settings: map[string]string{
			"privacy":           "private",
			"message_retention": "30d",
			"max_participants":  "100",
		},
	}

	resp, err := suite.makeRequest("POST", url, createReq, nil)
	require.NoError(suite.T(), err, "Create dialog request should succeed")
	defer resp.Body.Close()

	// Should return 201 for successful creation
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var messagingResp MessagingResponse
	err = suite.parseResponse(resp, &messagingResp)
	require.NoError(suite.T(), err, "Should parse dialog creation response")

	assert.True(suite.T(), messagingResp.Success)
	assert.Equal(suite.T(), "success", messagingResp.Status)
	assert.NotNil(suite.T(), messagingResp.Dialog)
	assert.NotEmpty(suite.T(), messagingResp.Dialog.ID)
	assert.Equal(suite.T(), createReq.Name, messagingResp.Dialog.Name)
	assert.Equal(suite.T(), createReq.Type, messagingResp.Dialog.Type)
	assert.Equal(suite.T(), len(createReq.Participants), len(messagingResp.Dialog.Participants))
}

// TestGetDialog tests retrieving a specific dialog
func (suite *MessagingIntegrationSuite) TestGetDialog() {
	// First create a dialog
	createURL := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)
	createReq := CreateDialogRequest{
		Name:        "Get Test Dialog",
		Type:        "private",
		Description: "Dialog for get test",
		Participants: []string{"user1", "user2"},
		Settings:    map[string]string{"privacy": "private"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create dialog for get test should succeed")
	defer createResp.Body.Close()

	var createMessagingResp MessagingResponse
	err = suite.parseResponse(createResp, &createMessagingResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createMessagingResp.Dialog)

	dialogID := createMessagingResp.Dialog.ID

	// Now get the dialog
	getURL := fmt.Sprintf("%s:%d/api/v1/dialogs/%s", suite.baseURL, suite.ports.Messaging, dialogID)
	resp, err := suite.makeRequest("GET", getURL, nil, nil)
	require.NoError(suite.T(), err, "Get dialog request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var messagingResp MessagingResponse
	err = suite.parseResponse(resp, &messagingResp)
	require.NoError(suite.T(), err, "Should parse get dialog response")

	assert.True(suite.T(), messagingResp.Success)
	assert.NotNil(suite.T(), messagingResp.Dialog)
	assert.Equal(suite.T(), dialogID, messagingResp.Dialog.ID)
	assert.Equal(suite.T(), createReq.Name, messagingResp.Dialog.Name)
}

// TestListDialogs tests listing dialogs
func (suite *MessagingIntegrationSuite) TestListDialogs() {
	url := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List dialogs request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var messagingResp MessagingResponse
	err = suite.parseResponse(resp, &messagingResp)
	require.NoError(suite.T(), err, "Should parse list dialogs response")

	assert.True(suite.T(), messagingResp.Success)
	assert.NotNil(suite.T(), messagingResp.Dialogs)
	assert.GreaterOrEqual(suite.T(), messagingResp.Total, 0)
}

// TestSendMessage tests sending a message to a dialog
func (suite *MessagingIntegrationSuite) TestSendMessage() {
	// First create a dialog
	createDialogURL := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)
	createDialogReq := CreateDialogRequest{
		Name:        "Message Test Dialog",
		Type:        "group",
		Description: "Dialog for message testing",
		Participants: []string{"user1", "user2"},
		Settings:    map[string]string{"privacy": "private"},
	}

	createDialogResp, err := suite.makeRequest("POST", createDialogURL, createDialogReq, nil)
	require.NoError(suite.T(), err, "Create dialog for message test should succeed")
	defer createDialogResp.Body.Close()

	var createDialogMessagingResp MessagingResponse
	err = suite.parseResponse(createDialogResp, &createDialogMessagingResp)
	require.NoError(suite.T(), err, "Should parse create dialog response")
	require.NotNil(suite.T(), createDialogMessagingResp.Dialog)

	dialogID := createDialogMessagingResp.Dialog.ID

	// Now send a message
	sendMessageURL := fmt.Sprintf("%s:%d/api/v1/dialogs/%s/messages", suite.baseURL, suite.ports.Messaging, dialogID)
	sendMessageReq := SendMessageRequest{
		Type:    "text",
		Content: "Hello, this is a test message!",
		Metadata: map[string]string{
			"client_version": "1.0.0",
			"device_type":    "mobile",
		},
	}

	// Add sender authentication header
	headers := map[string]string{
		"X-User-ID": "user1",
	}

	resp, err := suite.makeRequest("POST", sendMessageURL, sendMessageReq, headers)
	require.NoError(suite.T(), err, "Send message request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var messagingResp MessagingResponse
	err = suite.parseResponse(resp, &messagingResp)
	require.NoError(suite.T(), err, "Should parse send message response")

	assert.True(suite.T(), messagingResp.Success)
	assert.NotNil(suite.T(), messagingResp.ChatMessage)
	assert.NotEmpty(suite.T(), messagingResp.ChatMessage.ID)
	assert.Equal(suite.T(), dialogID, messagingResp.ChatMessage.DialogID)
	assert.Equal(suite.T(), sendMessageReq.Content, messagingResp.ChatMessage.Content)
	assert.Equal(suite.T(), "sent", messagingResp.ChatMessage.Status) // Default status
}

// TestGetMessages tests retrieving messages from a dialog
func (suite *MessagingIntegrationSuite) TestGetMessages() {
	// First create a dialog and send a message
	createDialogURL := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)
	createDialogReq := CreateDialogRequest{
		Name:        "Get Messages Test Dialog",
		Type:        "private",
		Description: "Dialog for get messages test",
		Participants: []string{"user1", "user2"},
	}

	createDialogResp, err := suite.makeRequest("POST", createDialogURL, createDialogReq, nil)
	require.NoError(suite.T(), err, "Create dialog should succeed")
	defer createDialogResp.Body.Close()

	var createDialogMessagingResp MessagingResponse
	err = suite.parseResponse(createDialogResp, &createDialogMessagingResp)
	require.NoError(suite.T(), err, "Should parse create dialog response")
	require.NotNil(suite.T(), createDialogMessagingResp.Dialog)

	dialogID := createDialogMessagingResp.Dialog.ID

	// Send a test message
	sendMessageURL := fmt.Sprintf("%s:%d/api/v1/dialogs/%s/messages", suite.baseURL, suite.ports.Messaging, dialogID)
	sendMessageReq := SendMessageRequest{
		Type:    "text",
		Content: "Test message for retrieval",
	}

	headers := map[string]string{"X-User-ID": "user1"}
	sendResp, err := suite.makeRequest("POST", sendMessageURL, sendMessageReq, headers)
	require.NoError(suite.T(), err, "Send message should succeed")
	sendResp.Body.Close()

	// Now get messages
	getMessagesURL := fmt.Sprintf("%s:%d/api/v1/dialogs/%s/messages", suite.baseURL, suite.ports.Messaging, dialogID)
	resp, err := suite.makeRequest("GET", getMessagesURL, nil, nil)
	require.NoError(suite.T(), err, "Get messages request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var messagingResp MessagingResponse
	err = suite.parseResponse(resp, &messagingResp)
	require.NoError(suite.T(), err, "Should parse get messages response")

	assert.True(suite.T(), messagingResp.Success)
	assert.NotNil(suite.T(), messagingResp.Messages)
	assert.GreaterOrEqual(suite.T(), len(messagingResp.Messages), 1)

	// Check that at least one message has the content we sent
	found := false
	for _, msg := range messagingResp.Messages {
		if msg.Content == sendMessageReq.Content {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "Should find the sent message in the response")
}

// TestUpdateMessage tests updating a message
func (suite *MessagingIntegrationSuite) TestUpdateMessage() {
	// First create a dialog and send a message
	createDialogURL := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)
	createDialogReq := CreateDialogRequest{
		Name:        "Update Message Test Dialog",
		Type:        "private",
		Participants: []string{"user1", "user2"},
	}

	createDialogResp, err := suite.makeRequest("POST", createDialogURL, createDialogReq, nil)
	require.NoError(suite.T(), err, "Create dialog should succeed")
	defer createDialogResp.Body.Close()

	var createDialogMessagingResp MessagingResponse
	err = suite.parseResponse(createDialogResp, &createDialogMessagingResp)
	require.NoError(suite.T(), err, "Should parse create dialog response")
	require.NotNil(suite.T(), createDialogMessagingResp.Dialog)

	dialogID := createDialogMessagingResp.Dialog.ID

	// Send a message
	sendMessageURL := fmt.Sprintf("%s:%d/api/v1/dialogs/%s/messages", suite.baseURL, suite.ports.Messaging, dialogID)
	sendMessageReq := SendMessageRequest{
		Type:    "text",
		Content: "Original message content",
	}

	headers := map[string]string{"X-User-ID": "user1"}
	sendResp, err := suite.makeRequest("POST", sendMessageURL, sendMessageReq, headers)
	require.NoError(suite.T(), err, "Send message should succeed")
	defer sendResp.Body.Close()

	var sendMessagingResp MessagingResponse
	err = suite.parseResponse(sendResp, &sendMessagingResp)
	require.NoError(suite.T(), err, "Should parse send message response")
	require.NotNil(suite.T(), sendMessagingResp.ChatMessage)

	messageID := sendMessagingResp.ChatMessage.ID

	// Update the message
	updateMessageURL := fmt.Sprintf("%s:%d/api/v1/messages/%s", suite.baseURL, suite.ports.Messaging, messageID)
	newContent := "Updated message content"
	updateMessageReq := UpdateMessageRequest{
		Content: &newContent,
		Metadata: map[string]string{
			"edited": "true",
			"reason": "typo correction",
		},
	}

	resp, err := suite.makeRequest("PUT", updateMessageURL, updateMessageReq, headers)
	require.NoError(suite.T(), err, "Update message request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var messagingResp MessagingResponse
	err = suite.parseResponse(resp, &messagingResp)
	require.NoError(suite.T(), err, "Should parse update message response")

	assert.True(suite.T(), messagingResp.Success)
	assert.NotNil(suite.T(), messagingResp.ChatMessage)
	assert.Equal(suite.T(), messageID, messagingResp.ChatMessage.ID)
	assert.Equal(suite.T(), newContent, messagingResp.ChatMessage.Content)
	assert.NotNil(suite.T(), messagingResp.ChatMessage.EditedAt)
}

// TestDeleteMessage tests deleting a message
func (suite *MessagingIntegrationSuite) TestDeleteMessage() {
	// First create a dialog and send a message
	createDialogURL := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)
	createDialogReq := CreateDialogRequest{
		Name:        "Delete Message Test Dialog",
		Type:        "private",
		Participants: []string{"user1", "user2"},
	}

	createDialogResp, err := suite.makeRequest("POST", createDialogURL, createDialogReq, nil)
	require.NoError(suite.T(), err, "Create dialog should succeed")
	defer createDialogResp.Body.Close()

	var createDialogMessagingResp MessagingResponse
	err = suite.parseResponse(createDialogResp, &createDialogMessagingResp)
	require.NoError(suite.T(), err, "Should parse create dialog response")
	require.NotNil(suite.T(), createDialogMessagingResp.Dialog)

	dialogID := createDialogMessagingResp.Dialog.ID

	// Send a message
	sendMessageURL := fmt.Sprintf("%s:%d/api/v1/dialogs/%s/messages", suite.baseURL, suite.ports.Messaging, dialogID)
	sendMessageReq := SendMessageRequest{
		Type:    "text",
		Content: "Message to be deleted",
	}

	headers := map[string]string{"X-User-ID": "user1"}
	sendResp, err := suite.makeRequest("POST", sendMessageURL, sendMessageReq, headers)
	require.NoError(suite.T(), err, "Send message should succeed")
	defer sendResp.Body.Close()

	var sendMessagingResp MessagingResponse
	err = suite.parseResponse(sendResp, &sendMessagingResp)
	require.NoError(suite.T(), err, "Should parse send message response")
	require.NotNil(suite.T(), sendMessagingResp.ChatMessage)

	messageID := sendMessagingResp.ChatMessage.ID

	// Delete the message
	deleteMessageURL := fmt.Sprintf("%s:%d/api/v1/messages/%s", suite.baseURL, suite.ports.Messaging, messageID)
	resp, err := suite.makeRequest("DELETE", deleteMessageURL, nil, headers)
	require.NoError(suite.T(), err, "Delete message request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify message is deleted by trying to get it
	getMessageURL := fmt.Sprintf("%s:%d/api/v1/messages/%s", suite.baseURL, suite.ports.Messaging, messageID)
	getResp, err := suite.makeRequest("GET", getMessageURL, nil, nil)
	require.NoError(suite.T(), err, "Get deleted message should complete")
	defer getResp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, getResp.StatusCode)
}

// TestGetNonExistentDialog tests retrieving non-existent dialog
func (suite *MessagingIntegrationSuite) TestGetNonExistentDialog() {
	url := fmt.Sprintf("%s:%d/api/v1/dialogs/non-existent-id", suite.baseURL, suite.ports.Messaging)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Get non-existent dialog should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// TestCreateDialogInvalidData tests dialog creation with invalid data
func (suite *MessagingIntegrationSuite) TestCreateDialogInvalidData() {
	url := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)

	// Test with missing required fields
	invalidReq := CreateDialogRequest{
		// Missing name and participants
		Type:        "group",
		Description: "Dialog without name",
	}

	resp, err := suite.makeRequest("POST", url, invalidReq, nil)
	require.NoError(suite.T(), err, "Invalid dialog creation should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// TestDialogsByType tests filtering dialogs by type
func (suite *MessagingIntegrationSuite) TestDialogsByType() {
	// First create a dialog with specific type
	createURL := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)
	createReq := CreateDialogRequest{
		Name:        "Type Filter Dialog",
		Type:        "test-type-filter",
		Description: "Dialog for type filtering test",
		Participants: []string{"user1", "user2"},
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create dialog for type test should succeed")
	createResp.Body.Close()

	// Now filter by type
	url := fmt.Sprintf("%s:%d/api/v1/dialogs?type=test-type-filter", suite.baseURL, suite.ports.Messaging)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Type filter request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var messagingResp MessagingResponse
	err = suite.parseResponse(resp, &messagingResp)
	require.NoError(suite.T(), err, "Should parse type filter response")

	assert.True(suite.T(), messagingResp.Success)

	// All returned dialogs should have the specified type
	for _, dialog := range messagingResp.Dialogs {
		assert.Equal(suite.T(), "test-type-filter", dialog.Type)
	}
}

// TestMessagePagination tests message pagination
func (suite *MessagingIntegrationSuite) TestMessagePagination() {
	// First create a dialog
	createDialogURL := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)
	createDialogReq := CreateDialogRequest{
		Name:        "Pagination Test Dialog",
		Type:        "private",
		Participants: []string{"user1", "user2"},
	}

	createDialogResp, err := suite.makeRequest("POST", createDialogURL, createDialogReq, nil)
	require.NoError(suite.T(), err, "Create dialog should succeed")
	defer createDialogResp.Body.Close()

	var createDialogMessagingResp MessagingResponse
	err = suite.parseResponse(createDialogResp, &createDialogMessagingResp)
	require.NoError(suite.T(), err, "Should parse create dialog response")
	require.NotNil(suite.T(), createDialogMessagingResp.Dialog)

	dialogID := createDialogMessagingResp.Dialog.ID

	// Test pagination (may be empty initially)
	url := fmt.Sprintf("%s:%d/api/v1/dialogs/%s/messages?page=1&limit=10", suite.baseURL, suite.ports.Messaging, dialogID)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Paginated messages request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var messagingResp MessagingResponse
	err = suite.parseResponse(resp, &messagingResp)
	require.NoError(suite.T(), err, "Should parse paginated response")

	assert.True(suite.T(), messagingResp.Success)
	assert.NotNil(suite.T(), messagingResp.Messages)
	assert.GreaterOrEqual(suite.T(), messagingResp.Total, 0)
}

// TestInvalidHTTPMethods tests endpoints with invalid HTTP methods
func (suite *MessagingIntegrationSuite) TestInvalidHTTPMethods() {
	baseURL := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)

	testCases := []struct {
		url    string
		method string
	}{
		{baseURL, "PATCH"},           // List endpoint with invalid method
		{baseURL + "/123", "POST"},   // Get endpoint with invalid method
		{baseURL + "/123", "PATCH"},  // Update endpoint with invalid method
	}

	for _, tc := range testCases {
		resp, err := suite.makeRequest(tc.method, tc.url, nil, nil)
		require.NoError(suite.T(), err, "Invalid method request should complete")
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusMethodNotAllowed, resp.StatusCode,
			"URL: %s, Method: %s", tc.url, tc.method)
	}
}

// TestMessagingServiceConcurrency tests concurrent requests to messaging service
func (suite *MessagingIntegrationSuite) TestMessagingServiceConcurrency() {
	url := fmt.Sprintf("%s:%d/api/v1/dialogs", suite.baseURL, suite.ports.Messaging)

	// Create 5 concurrent dialog creation requests
	concurrency := 5
	results := make(chan int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			createReq := CreateDialogRequest{
				Name:        fmt.Sprintf("Concurrent Dialog %d", index),
				Type:        "concurrency-test",
				Description: fmt.Sprintf("Dialog created concurrently #%d", index),
				Participants: []string{fmt.Sprintf("user%d", index), "user0"},
			}

			resp, err := suite.makeRequest("POST", url, createReq, nil)
			if err != nil {
				results <- 0
				return
			}
			defer resp.Body.Close()

			results <- resp.StatusCode
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		statusCode := <-results
		if statusCode == http.StatusCreated {
			successCount++
		}
	}

	// At least 80% of concurrent requests should succeed
	assert.GreaterOrEqual(suite.T(), successCount, 4, "Concurrent requests should mostly succeed")
}

// RunMessagingIntegrationTests runs the messaging integration test suite
func RunMessagingIntegrationTests(t *testing.T) {
	suite.Run(t, new(MessagingIntegrationSuite))
}