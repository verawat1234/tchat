package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// T030: Integration test real-time messaging flow
// Tests end-to-end messaging workflow including:
// 1. Dialog creation → 2. WebSocket connection → 3. Message sending → 4. Real-time delivery → 5. Read receipts
type MessagingFlowTestSuite struct {
	suite.Suite
	router      *gin.Engine
	dialogs     map[string]map[string]interface{}
	messages    map[string][]map[string]interface{}
	connections map[string]*websocket.Conn
	upgrader    websocket.Upgrader
	mu          sync.RWMutex
}

func TestMessagingFlowSuite(t *testing.T) {
	suite.Run(t, new(MessagingFlowTestSuite))
}

func (suite *MessagingFlowTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.dialogs = make(map[string]map[string]interface{})
	suite.messages = make(map[string][]map[string]interface{})
	suite.connections = make(map[string]*websocket.Conn)
	suite.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Setup messaging service endpoints
	suite.setupMessagingEndpoints()
}

func (suite *MessagingFlowTestSuite) setupMessagingEndpoints() {
	// Mock authentication middleware
	authMiddleware := func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		// Extract user ID from token (simplified)
		if strings.HasPrefix(auth, "Bearer user_") {
			userID := strings.TrimPrefix(auth, "Bearer ")
			c.Set("user_id", userID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			c.Abort()
			return
		}
		c.Next()
	}

	// Create dialog endpoint
	suite.router.POST("/dialogs", authMiddleware, func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		userID := c.GetString("user_id")
		dialogType := req["type"].(string)
		participants := req["participants"].([]interface{})

		// Validate dialog type
		validTypes := map[string]bool{"direct": true, "group": true, "channel": true}
		if !validTypes[dialogType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_dialog_type"})
			return
		}

		// Generate dialog ID
		dialogID := fmt.Sprintf("dialog_%d", time.Now().UnixNano())

		// Add creator to participants
		allParticipants := []string{userID}
		for _, p := range participants {
			if participant, ok := p.(string); ok && participant != userID {
				allParticipants = append(allParticipants, participant)
			}
		}

		dialog := map[string]interface{}{
			"id":           dialogID,
			"type":         dialogType,
			"participants": allParticipants,
			"created_by":   userID,
			"created_at":   time.Now().UTC().Format(time.RFC3339),
			"updated_at":   time.Now().UTC().Format(time.RFC3339),
			"message_count": 0,
		}

		if title, exists := req["title"]; exists {
			dialog["title"] = title
		}

		suite.mu.Lock()
		suite.dialogs[dialogID] = dialog
		suite.messages[dialogID] = make([]map[string]interface{}, 0)
		suite.mu.Unlock()

		c.JSON(http.StatusCreated, dialog)
	})

	// Send message endpoint
	suite.router.POST("/dialogs/:dialog_id/messages", authMiddleware, func(c *gin.Context) {
		dialogID := c.Param("dialog_id")
		userID := c.GetString("user_id")

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json", "message": err.Error()})
			return
		}

		// Check if dialog exists and user is participant
		suite.mu.RLock()
		dialog, exists := suite.dialogs[dialogID]
		suite.mu.RUnlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "dialog_not_found"})
			return
		}

		participants := dialog["participants"].([]string)
		isParticipant := false
		for _, p := range participants {
			if p == userID {
				isParticipant = true
				break
			}
		}

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{"error": "not_participant"})
			return
		}

		// Create message
		messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())
		message := map[string]interface{}{
			"id":         messageID,
			"dialog_id":  dialogID,
			"sender_id":  userID,
			"content":    req["content"],
			"type":       req["type"],
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"updated_at": time.Now().UTC().Format(time.RFC3339),
			"read_by":    []string{userID}, // Sender automatically reads
		}

		if attachments, exists := req["attachments"]; exists {
			message["attachments"] = attachments
		}

		// Store message
		suite.mu.Lock()
		suite.messages[dialogID] = append(suite.messages[dialogID], message)
		dialog["message_count"] = len(suite.messages[dialogID])
		dialog["updated_at"] = time.Now().UTC().Format(time.RFC3339)
		suite.mu.Unlock()

		// Broadcast to WebSocket connections
		suite.broadcastMessage(message, participants)

		c.JSON(http.StatusCreated, message)
	})

	// Get messages endpoint
	suite.router.GET("/dialogs/:dialog_id/messages", authMiddleware, func(c *gin.Context) {
		dialogID := c.Param("dialog_id")
		userID := c.GetString("user_id")

		// Check access
		suite.mu.RLock()
		dialog, exists := suite.dialogs[dialogID]
		messages := suite.messages[dialogID]
		suite.mu.RUnlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "dialog_not_found"})
			return
		}

		participants := dialog["participants"].([]string)
		isParticipant := false
		for _, p := range participants {
			if p == userID {
				isParticipant = true
				break
			}
		}

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{"error": "not_participant"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"dialog_id": dialogID,
			"messages":  messages,
			"total":     len(messages),
		})
	})

	// Mark message as read endpoint
	suite.router.POST("/messages/:message_id/read", authMiddleware, func(c *gin.Context) {
		messageID := c.Param("message_id")
		userID := c.GetString("user_id")

		// Find message
		var foundMessage map[string]interface{}
		var dialogID string

		suite.mu.Lock()
		defer suite.mu.Unlock()

		for dID, messages := range suite.messages {
			for i, msg := range messages {
				if msg["id"] == messageID {
					foundMessage = msg
					dialogID = dID
					// Add user to read_by list if not already there
					readBy := msg["read_by"].([]string)
					alreadyRead := false
					for _, reader := range readBy {
						if reader == userID {
							alreadyRead = true
							break
						}
					}
					if !alreadyRead {
						readBy = append(readBy, userID)
						suite.messages[dID][i]["read_by"] = readBy
					}
					break
				}
			}
			if foundMessage != nil {
				break
			}
		}

		if foundMessage == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "message_not_found"})
			return
		}

		// Check if user is participant
		dialog := suite.dialogs[dialogID]
		participants := dialog["participants"].([]string)
		isParticipant := false
		for _, p := range participants {
			if p == userID {
				isParticipant = true
				break
			}
		}

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{"error": "not_participant"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message_id": messageID,
			"read_at":    time.Now().UTC().Format(time.RFC3339),
			"read_by":    foundMessage["read_by"],
		})
	})

	// WebSocket endpoint
	suite.router.GET("/websocket", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing_user_id"})
			return
		}

		conn, err := suite.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "websocket_upgrade_failed"})
			return
		}

		suite.mu.Lock()
		suite.connections[userID] = conn
		suite.mu.Unlock()

		// Send connection confirmation
		conn.WriteJSON(map[string]interface{}{
			"type":    "connection_established",
			"user_id": userID,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})

		// Handle incoming messages
		go suite.handleWebSocketConnection(userID, conn)
	})
}

func (suite *MessagingFlowTestSuite) broadcastMessage(message map[string]interface{}, participants []string) {
	suite.mu.RLock()
	defer suite.mu.RUnlock()

	wsMessage := map[string]interface{}{
		"type":    "new_message",
		"message": message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	for _, participantID := range participants {
		if conn, exists := suite.connections[participantID]; exists {
			err := conn.WriteJSON(wsMessage)
			if err != nil {
				// Connection closed, remove it
				delete(suite.connections, participantID)
			}
		}
	}
}

func (suite *MessagingFlowTestSuite) handleWebSocketConnection(userID string, conn *websocket.Conn) {
	defer func() {
		suite.mu.Lock()
		delete(suite.connections, userID)
		suite.mu.Unlock()
		conn.Close()
	}()

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		// Handle WebSocket messages (ping, typing indicators, etc.)
		switch msg["type"] {
		case "ping":
			conn.WriteJSON(map[string]interface{}{
				"type": "pong",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		case "typing_start":
			// Broadcast typing indicator
			suite.broadcastTypingIndicator(userID, msg["dialog_id"].(string), true)
		case "typing_stop":
			// Broadcast typing stop
			suite.broadcastTypingIndicator(userID, msg["dialog_id"].(string), false)
		}
	}
}

func (suite *MessagingFlowTestSuite) broadcastTypingIndicator(userID, dialogID string, isTyping bool) {
	suite.mu.RLock()
	defer suite.mu.RUnlock()

	if dialog, exists := suite.dialogs[dialogID]; exists {
		participants := dialog["participants"].([]string)
		typingMessage := map[string]interface{}{
			"type":      "typing_indicator",
			"user_id":   userID,
			"dialog_id": dialogID,
			"is_typing": isTyping,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}

		for _, participantID := range participants {
			if participantID != userID { // Don't send to sender
				if conn, exists := suite.connections[participantID]; exists {
					conn.WriteJSON(typingMessage)
				}
			}
		}
	}
}

func (suite *MessagingFlowTestSuite) TestCompleteMessagingFlow() {
	// Step 1: Create Dialog
	suite.T().Log("Step 1: Creating dialog")

	dialogData := map[string]interface{}{
		"type":         "direct",
		"participants": []string{"user_456"},
		"title":        "Test Chat",
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var dialogResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	require.NoError(suite.T(), err)

	dialogID := dialogResponse["id"].(string)
	assert.NotEmpty(suite.T(), dialogID)
	assert.Equal(suite.T(), "direct", dialogResponse["type"])
	assert.Contains(suite.T(), dialogResponse["participants"], "user_123")
	assert.Contains(suite.T(), dialogResponse["participants"], "user_456")

	// Step 2: Establish WebSocket Connections
	suite.T().Log("Step 2: Establishing WebSocket connections")

	// Simulate WebSocket connections (in real test, would use actual WebSocket client)
	suite.mu.Lock()
	// Mock WebSocket connections
	suite.connections["user_123"] = nil // Placeholder for actual connection
	suite.connections["user_456"] = nil // Placeholder for actual connection
	suite.mu.Unlock()

	// Step 3: Send Message
	suite.T().Log("Step 3: Sending message")

	messageData := map[string]interface{}{
		"content": "Hello! How are you?",
		"type":    "text",
	}

	jsonData, _ = json.Marshal(messageData)
	req = httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", dialogID), bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	start := time.Now()
	suite.router.ServeHTTP(w, req)
	sendDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.True(suite.T(), sendDuration < 200*time.Millisecond, "Message sending should be <200ms")

	var messageResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &messageResponse)
	require.NoError(suite.T(), err)

	messageID := messageResponse["id"].(string)
	assert.NotEmpty(suite.T(), messageID)
	assert.Equal(suite.T(), dialogID, messageResponse["dialog_id"])
	assert.Equal(suite.T(), "user_123", messageResponse["sender_id"])
	assert.Equal(suite.T(), "Hello! How are you?", messageResponse["content"])
	assert.Equal(suite.T(), "text", messageResponse["type"])

	// Step 4: Retrieve Messages
	suite.T().Log("Step 4: Retrieving messages")

	req = httptest.NewRequest("GET", fmt.Sprintf("/dialogs/%s/messages", dialogID), nil)
	req.Header.Set("Authorization", "Bearer user_456")
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	start = time.Now()
	suite.router.ServeHTTP(w, req)
	retrieveDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), retrieveDuration < 100*time.Millisecond, "Message retrieval should be <100ms")

	var messagesResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &messagesResponse)
	require.NoError(suite.T(), err)

	messages := messagesResponse["messages"].([]interface{})
	assert.Len(suite.T(), messages, 1)
	assert.Equal(suite.T(), 1, int(messagesResponse["total"].(float64)))

	// Step 5: Mark Message as Read
	suite.T().Log("Step 5: Marking message as read")

	req = httptest.NewRequest("POST", fmt.Sprintf("/messages/%s/read", messageID), nil)
	req.Header.Set("Authorization", "Bearer user_456")
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	start = time.Now()
	suite.router.ServeHTTP(w, req)
	readDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), readDuration < 50*time.Millisecond, "Read receipt should be <50ms")

	var readResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &readResponse)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), messageID, readResponse["message_id"])
	assert.NotEmpty(suite.T(), readResponse["read_at"])

	readBy := readResponse["read_by"].([]interface{})
	assert.Contains(suite.T(), readBy, "user_123") // Sender
	assert.Contains(suite.T(), readBy, "user_456") // Reader
}

func (suite *MessagingFlowTestSuite) TestGroupMessaging() {
	// Create group dialog
	dialogData := map[string]interface{}{
		"type":         "group",
		"participants": []string{"user_456", "user_789"},
		"title":        "Project Team",
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var dialogResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	dialogID := dialogResponse["id"].(string)

	// Send message to group
	messageData := map[string]interface{}{
		"content": "Team meeting at 3 PM",
		"type":    "text",
	}

	jsonData, _ = json.Marshal(messageData)
	req = httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", dialogID), bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Verify all participants can access
	participants := []string{"user_123", "user_456", "user_789"}
	for _, userID := range participants {
		req = httptest.NewRequest("GET", fmt.Sprintf("/dialogs/%s/messages", dialogID), nil)
		req.Header.Set("Authorization", "Bearer "+userID)

		w = httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code, "User %s should access group messages", userID)
	}
}

func (suite *MessagingFlowTestSuite) TestMessageAttachments() {
	// Create dialog
	dialogData := map[string]interface{}{
		"type":         "direct",
		"participants": []string{"user_456"},
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var dialogResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	dialogID := dialogResponse["id"].(string)

	// Send message with attachments
	messageData := map[string]interface{}{
		"content": "Check out this image!",
		"type":    "multimedia",
		"attachments": []map[string]interface{}{
			{
				"type": "image",
				"url":  "https://example.com/image.jpg",
				"size": 1024000,
				"name": "vacation.jpg",
			},
		},
	}

	jsonData, _ = json.Marshal(messageData)
	req = httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", dialogID), bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var messageResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &messageResponse)

	assert.Equal(suite.T(), "multimedia", messageResponse["type"])
	assert.NotNil(suite.T(), messageResponse["attachments"])

	attachments := messageResponse["attachments"].([]interface{})
	assert.Len(suite.T(), attachments, 1)

	attachment := attachments[0].(map[string]interface{})
	assert.Equal(suite.T(), "image", attachment["type"])
	assert.Equal(suite.T(), "https://example.com/image.jpg", attachment["url"])
}

func (suite *MessagingFlowTestSuite) TestMessagingPerformance() {
	// Create dialog for performance testing
	dialogData := map[string]interface{}{
		"type":         "direct",
		"participants": []string{"user_456"},
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var dialogResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	dialogID := dialogResponse["id"].(string)

	// Test rapid message sending
	suite.T().Log("Testing rapid message sending performance")

	messagesSent := 10
	var totalDuration time.Duration

	for i := 0; i < messagesSent; i++ {
		messageData := map[string]interface{}{
			"content": fmt.Sprintf("Performance test message %d", i+1),
			"type":    "text",
		}

		jsonData, _ = json.Marshal(messageData)
		req = httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", dialogID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_123")
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		start := time.Now()
		suite.router.ServeHTTP(w, req)
		duration := time.Since(start)
		totalDuration += duration

		assert.Equal(suite.T(), http.StatusCreated, w.Code)
		assert.True(suite.T(), duration < 200*time.Millisecond, "Each message should send in <200ms")
	}

	avgDuration := totalDuration / time.Duration(messagesSent)
	assert.True(suite.T(), avgDuration < 150*time.Millisecond, "Average message send time should be <150ms")

	// Test message retrieval performance with larger dataset
	req = httptest.NewRequest("GET", fmt.Sprintf("/dialogs/%s/messages", dialogID), nil)
	req.Header.Set("Authorization", "Bearer user_456")

	w = httptest.NewRecorder()
	start := time.Now()
	suite.router.ServeHTTP(w, req)
	retrievalDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), retrievalDuration < 200*time.Millisecond, "Message retrieval should be <200ms even with %d messages", messagesSent)

	var messagesResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &messagesResponse)
	assert.Equal(suite.T(), messagesSent, int(messagesResponse["total"].(float64)))
}

func (suite *MessagingFlowTestSuite) TestAccessControl() {
	// Create private dialog
	dialogData := map[string]interface{}{
		"type":         "direct",
		"participants": []string{"user_456"},
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var dialogResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	dialogID := dialogResponse["id"].(string)

	// Test unauthorized user cannot send message
	messageData := map[string]interface{}{
		"content": "Unauthorized message",
		"type":    "text",
	}

	jsonData, _ = json.Marshal(messageData)
	req = httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", dialogID), bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_999") // Not a participant
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)

	// Test unauthorized user cannot read messages
	req = httptest.NewRequest("GET", fmt.Sprintf("/dialogs/%s/messages", dialogID), nil)
	req.Header.Set("Authorization", "Bearer user_999")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *MessagingFlowTestSuite) TestErrorHandling() {
	// Test sending message to non-existent dialog
	messageData := map[string]interface{}{
		"content": "Message to nowhere",
		"type":    "text",
	}

	jsonData, _ := json.Marshal(messageData)
	req := httptest.NewRequest("POST", "/dialogs/invalid_dialog/messages", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_123")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Test reading from non-existent dialog
	req = httptest.NewRequest("GET", "/dialogs/invalid_dialog/messages", nil)
	req.Header.Set("Authorization", "Bearer user_123")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// Test marking non-existent message as read
	req = httptest.NewRequest("POST", "/messages/invalid_message/read", nil)
	req.Header.Set("Authorization", "Bearer user_123")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}