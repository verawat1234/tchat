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

// Comprehensive messaging stress testing suite
// Tests edge cases, concurrency, performance, and reliability
type MessagingStressTestSuite struct {
	suite.Suite
	router      *gin.Engine
	dialogs     map[string]map[string]interface{}
	messages    map[string][]map[string]interface{}
	connections map[string]*websocket.Conn
	upgrader    websocket.Upgrader
	mu          sync.RWMutex
	metrics     map[string]interface{}
}

func TestMessagingStressTestSuite(t *testing.T) {
	suite.Run(t, new(MessagingStressTestSuite))
}

func (suite *MessagingStressTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.dialogs = make(map[string]map[string]interface{})
	suite.messages = make(map[string][]map[string]interface{})
	suite.connections = make(map[string]*websocket.Conn)
	suite.metrics = make(map[string]interface{})
	suite.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	suite.setupMessagingEndpoints()
}

func (suite *MessagingStressTestSuite) setupMessagingEndpoints() {
	// Auth middleware
	authMiddleware := func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing_auth"})
			c.Abort()
			return
		}
		userID := strings.TrimPrefix(auth, "Bearer ")
		c.Set("user_id", userID)
		c.Next()
	}

	// Create dialog endpoint
	suite.router.POST("/dialogs", authMiddleware, func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json"})
			return
		}

		userID := c.GetString("user_id")
		participants := req["participants"].([]interface{})
		dialogType := req["type"].(string)

		dialogID := fmt.Sprintf("dialog_%d", time.Now().UnixNano())
		participantList := make([]string, len(participants))
		for i, p := range participants {
			participantList[i] = p.(string)
		}

		// Ensure creator is in participants
		found := false
		for _, p := range participantList {
			if p == userID {
				found = true
				break
			}
		}
		if !found {
			participantList = append(participantList, userID)
		}

		dialog := map[string]interface{}{
			"id":            dialogID,
			"type":          dialogType,
			"participants":  participantList,
			"creator_id":    userID,
			"message_count": 0,
			"created_at":    time.Now().UTC().Format(time.RFC3339),
			"updated_at":    time.Now().UTC().Format(time.RFC3339),
		}

		suite.mu.Lock()
		suite.dialogs[dialogID] = dialog
		suite.messages[dialogID] = make([]map[string]interface{}, 0)
		suite.mu.Unlock()

		c.JSON(http.StatusCreated, dialog)
	})

	// Send message endpoint with enhanced validation
	suite.router.POST("/dialogs/:dialog_id/messages", authMiddleware, func(c *gin.Context) {
		dialogID := c.Param("dialog_id")
		userID := c.GetString("user_id")

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json"})
			return
		}

		// Validate message content
		content, exists := req["content"]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content_required"})
			return
		}

		contentStr := content.(string)
		// Test message size limits
		if len(contentStr) > 10000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message_too_large"})
			return
		}

		suite.mu.RLock()
		dialog, exists := suite.dialogs[dialogID]
		if !exists {
			suite.mu.RUnlock()
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
		suite.mu.RUnlock()

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{"error": "not_participant"})
			return
		}

		// Create message with enhanced metadata
		messageID := fmt.Sprintf("msg_%d", time.Now().UnixNano())
		// Store message with thread safety
		suite.mu.Lock()

		message := map[string]interface{}{
			"id":         messageID,
			"dialog_id":  dialogID,
			"sender_id":  userID,
			"content":    contentStr,
			"type":       req["type"],
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"updated_at": time.Now().UTC().Format(time.RFC3339),
			"read_by":    []string{userID},
			"message_number": len(suite.messages[dialogID]) + 1, // Track ordering (now thread-safe)
		}

		if attachments, exists := req["attachments"]; exists {
			message["attachments"] = attachments
		}

		suite.messages[dialogID] = append(suite.messages[dialogID], message)
		dialog["message_count"] = len(suite.messages[dialogID])
		dialog["updated_at"] = time.Now().UTC().Format(time.RFC3339)
		suite.mu.Unlock()

		// Simulate message broadcasting (simplified)
		suite.broadcastMessage(message, participants)

		c.JSON(http.StatusCreated, message)
	})

	// Get messages endpoint
	suite.router.GET("/dialogs/:dialog_id/messages", authMiddleware, func(c *gin.Context) {
		dialogID := c.Param("dialog_id")
		userID := c.GetString("user_id")

		suite.mu.RLock()
		dialog, exists := suite.dialogs[dialogID]
		if !exists {
			suite.mu.RUnlock()
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

		messages := make([]map[string]interface{}, len(suite.messages[dialogID]))
		copy(messages, suite.messages[dialogID])
		suite.mu.RUnlock()

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{"error": "not_participant"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"messages": messages,
			"total":    len(messages),
		})
	})

	// WebSocket endpoint with enhanced error handling
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

func (suite *MessagingStressTestSuite) broadcastMessage(message map[string]interface{}, participants []string) {
	suite.mu.RLock()
	defer suite.mu.RUnlock()

	wsMessage := map[string]interface{}{
		"type":    "new_message",
		"message": message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	for _, participantID := range participants {
		if conn, exists := suite.connections[participantID]; exists && conn != nil {
			err := conn.WriteJSON(wsMessage)
			if err != nil {
				// Connection closed, remove it
				suite.mu.RUnlock()
				suite.mu.Lock()
				delete(suite.connections, participantID)
				suite.mu.Unlock()
				suite.mu.RLock()
			}
		}
	}
}

func (suite *MessagingStressTestSuite) handleWebSocketConnection(userID string, conn *websocket.Conn) {
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

		switch msg["type"] {
		case "ping":
			conn.WriteJSON(map[string]interface{}{
				"type": "pong",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		}
	}
}

// Test 1: Concurrent Message Sending
func (suite *MessagingStressTestSuite) TestConcurrentMessageSending() {
	suite.T().Log("Testing concurrent message sending with race condition detection")

	// Create dialog
	dialogData := map[string]interface{}{
		"type":         "group",
		"participants": []interface{}{"user_1", "user_2", "user_3"},
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusCreated, w.Code)

	var dialogResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	require.NoError(suite.T(), err)

	dialogID := dialogResponse["id"].(string)

	// Concurrent message sending
	const numMessages = 50
	const numUsers = 3

	var wg sync.WaitGroup
	messagesSent := make([][]int, numUsers)

	for i := 0; i < numUsers; i++ {
		messagesSent[i] = make([]int, 0)
	}

	start := time.Now()

	// Send messages concurrently from multiple users
	for userIdx := 0; userIdx < numUsers; userIdx++ {
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()
			userID := fmt.Sprintf("user_%d", userIndex+1)

			for msgIdx := 0; msgIdx < numMessages; msgIdx++ {
				messageData := map[string]interface{}{
					"type":    "text",
					"content": fmt.Sprintf("Concurrent message %d from %s", msgIdx, userID),
				}

				jsonData, _ := json.Marshal(messageData)
				req := httptest.NewRequest("POST", "/dialogs/"+dialogID+"/messages", bytes.NewBuffer(jsonData))
				req.Header.Set("Authorization", "Bearer "+userID)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				suite.router.ServeHTTP(w, req)

				if w.Code == http.StatusCreated {
					messagesSent[userIndex] = append(messagesSent[userIndex], msgIdx)
				}
			}
		}(userIdx)
	}

	wg.Wait()
	duration := time.Since(start)

	suite.T().Logf("Concurrent messaging completed in %v", duration)

	// Verify all messages were sent
	totalSent := 0
	for i, sent := range messagesSent {
		suite.T().Logf("User %d sent %d messages", i+1, len(sent))
		totalSent += len(sent)
	}

	assert.Equal(suite.T(), numMessages*numUsers, totalSent, "Not all messages were sent successfully")

	// Verify message ordering and integrity
	suite.verifyMessageIntegrity(dialogID, numMessages*numUsers)
}

// Test 2: Large Message Handling
func (suite *MessagingStressTestSuite) TestLargeMessageHandling() {
	suite.T().Log("Testing large message handling and size limits")

	// Create dialog
	dialogData := map[string]interface{}{
		"type":         "direct",
		"participants": []interface{}{"user_1", "user_2"},
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var dialogResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	dialogID := dialogResponse["id"].(string)

	// Test various message sizes
	testCases := []struct {
		name        string
		size        int
		shouldPass  bool
		expectedCode int
	}{
		{"Small message", 100, true, http.StatusCreated},
		{"Medium message", 1000, true, http.StatusCreated},
		{"Large message", 5000, true, http.StatusCreated},
		{"Very large message", 9000, true, http.StatusCreated},
		{"Oversized message", 15000, false, http.StatusBadRequest},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			content := strings.Repeat("A", tc.size)
			messageData := map[string]interface{}{
				"type":    "text",
				"content": content,
			}

			jsonData, _ := json.Marshal(messageData)
			req := httptest.NewRequest("POST", "/dialogs/"+dialogID+"/messages", bytes.NewBuffer(jsonData))
			req.Header.Set("Authorization", "Bearer user_1")
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code, "Unexpected response code for %s", tc.name)

			if tc.shouldPass {
				var messageResponse map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &messageResponse)
				assert.Equal(t, content, messageResponse["content"], "Message content mismatch")
			}
		})
	}
}

// Test 3: Message Ordering Validation
func (suite *MessagingStressTestSuite) TestMessageOrdering() {
	suite.T().Log("Testing message ordering under concurrent load")

	// Create dialog
	dialogData := map[string]interface{}{
		"type":         "direct",
		"participants": []interface{}{"user_1", "user_2"},
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var dialogResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	dialogID := dialogResponse["id"].(string)

	// Send ordered messages rapidly
	const numMessages = 100
	sentOrder := make([]int, 0, numMessages)
	var orderMux sync.Mutex

	start := time.Now()

	for i := 0; i < numMessages; i++ {
		go func(index int) {
			messageData := map[string]interface{}{
				"type":    "text",
				"content": fmt.Sprintf("Ordered message %d", index),
			}

			jsonData, _ := json.Marshal(messageData)
			req := httptest.NewRequest("POST", "/dialogs/"+dialogID+"/messages", bytes.NewBuffer(jsonData))
			req.Header.Set("Authorization", "Bearer user_1")
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			if w.Code == http.StatusCreated {
				orderMux.Lock()
				sentOrder = append(sentOrder, index)
				orderMux.Unlock()
			}
		}(i)
	}

	// Wait for all messages
	for len(sentOrder) < numMessages {
		time.Sleep(10 * time.Millisecond)
	}

	duration := time.Since(start)
	suite.T().Logf("Ordered messaging completed in %v", duration)

	// Retrieve messages and verify storage order
	req = httptest.NewRequest("GET", "/dialogs/"+dialogID+"/messages", nil)
	req.Header.Set("Authorization", "Bearer user_1")

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	messages := response["messages"].([]interface{})

	assert.Equal(suite.T(), numMessages, len(messages), "Not all messages stored")

	// Verify messages have sequential message numbers
	for i, msg := range messages {
		msgMap := msg.(map[string]interface{})
		expectedNumber := i + 1
		actualNumber := int(msgMap["message_number"].(float64))
		assert.Equal(suite.T(), expectedNumber, actualNumber, "Message number mismatch at position %d", i)
	}
}

// Test 4: Error Handling and Validation
func (suite *MessagingStressTestSuite) TestErrorHandling() {
	suite.T().Log("Testing comprehensive error handling scenarios")

	// Create dialog for testing
	dialogData := map[string]interface{}{
		"type":         "direct",
		"participants": []interface{}{"user_1", "user_2"},
	}

	jsonData, _ := json.Marshal(dialogData)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_1")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var dialogResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dialogResponse)
	dialogID := dialogResponse["id"].(string)

	testCases := []struct {
		name         string
		endpoint     string
		method       string
		data         map[string]interface{}
		auth         string
		expectedCode int
		expectedError string
	}{
		{
			name:          "Missing content",
			endpoint:      "/dialogs/" + dialogID + "/messages",
			method:        "POST",
			data:          map[string]interface{}{"type": "text"},
			auth:          "Bearer user_1",
			expectedCode:  http.StatusBadRequest,
			expectedError: "content_required",
		},
		{
			name:         "Unauthorized access",
			endpoint:     "/dialogs/" + dialogID + "/messages",
			method:       "POST",
			data:         map[string]interface{}{"type": "text", "content": "test"},
			auth:         "",
			expectedCode: http.StatusUnauthorized,
			expectedError: "missing_auth",
		},
		{
			name:         "Non-participant access",
			endpoint:     "/dialogs/" + dialogID + "/messages",
			method:       "POST",
			data:         map[string]interface{}{"type": "text", "content": "test"},
			auth:         "Bearer user_3",
			expectedCode: http.StatusForbidden,
			expectedError: "not_participant",
		},
		{
			name:         "Non-existent dialog",
			endpoint:     "/dialogs/nonexistent/messages",
			method:       "POST",
			data:         map[string]interface{}{"type": "text", "content": "test"},
			auth:         "Bearer user_1",
			expectedCode: http.StatusNotFound,
			expectedError: "dialog_not_found",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.method == "POST" {
				jsonData, _ := json.Marshal(tc.data)
				req = httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.endpoint, nil)
			}

			if tc.auth != "" {
				req.Header.Set("Authorization", tc.auth)
			}

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code, "Unexpected response code for %s", tc.name)

			if tc.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tc.expectedError, response["error"], "Unexpected error message for %s", tc.name)
			}
		})
	}
}

// Test 5: Performance Under Load
func (suite *MessagingStressTestSuite) TestPerformanceUnderLoad() {
	suite.T().Log("Testing messaging system performance under sustained load")

	const numDialogs = 10
	const messagesPerDialog = 20
	const concurrentUsers = 5

	dialogIDs := make([]string, numDialogs)

	// Create multiple dialogs with all test users as participants
	allUsers := make([]interface{}, concurrentUsers)
	for i := 1; i <= concurrentUsers; i++ {
		allUsers[i-1] = fmt.Sprintf("user_%d", i)
	}

	// Create multiple dialogs
	for i := 0; i < numDialogs; i++ {
		dialogData := map[string]interface{}{
			"type":         "group",
			"participants": allUsers,
		}

		jsonData, _ := json.Marshal(dialogData)
		req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_1")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		var dialogResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &dialogResponse)
		dialogIDs[i] = dialogResponse["id"].(string)
	}

	suite.T().Logf("Created %d dialogs for load testing", numDialogs)

	// Performance metrics
	var totalMessages int64
	var totalDuration time.Duration
	var successCount int64
	var errorCount int64
	var mu sync.Mutex

	start := time.Now()

	// Generate sustained load
	var wg sync.WaitGroup
	for userIdx := 1; userIdx <= concurrentUsers; userIdx++ {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()

			for dialogIdx := 0; dialogIdx < numDialogs; dialogIdx++ {
				for msgIdx := 0; msgIdx < messagesPerDialog; msgIdx++ {
					msgStart := time.Now()

					messageData := map[string]interface{}{
						"type":    "text",
						"content": fmt.Sprintf("Load test message %d from %s in dialog %d", msgIdx, userID, dialogIdx),
					}

					jsonData, _ := json.Marshal(messageData)
					req := httptest.NewRequest("POST", "/dialogs/"+dialogIDs[dialogIdx]+"/messages", bytes.NewBuffer(jsonData))
					req.Header.Set("Authorization", "Bearer "+userID)
					req.Header.Set("Content-Type", "application/json")

					w := httptest.NewRecorder()
					suite.router.ServeHTTP(w, req)

					msgDuration := time.Since(msgStart)

					mu.Lock()
					totalMessages++
					totalDuration += msgDuration
					if w.Code == http.StatusCreated {
						successCount++
					} else {
						errorCount++
						// Log first few errors for debugging
						if errorCount <= 5 {
							suite.T().Logf("Error %d: HTTP %d - %s", errorCount, w.Code, w.Body.String())
						}
					}
					mu.Unlock()
				}
			}
		}(fmt.Sprintf("user_%d", userIdx))
	}

	wg.Wait()
	totalTestDuration := time.Since(start)

	// Calculate performance metrics
	avgMessageTime := totalDuration / time.Duration(totalMessages)
	messagesPerSecond := float64(totalMessages) / totalTestDuration.Seconds()
	successRate := float64(successCount) / float64(totalMessages) * 100

	suite.T().Logf("Performance Results:")
	suite.T().Logf("- Total messages: %d", totalMessages)
	suite.T().Logf("- Total duration: %v", totalTestDuration)
	suite.T().Logf("- Average message processing time: %v", avgMessageTime)
	suite.T().Logf("- Messages per second: %.2f", messagesPerSecond)
	suite.T().Logf("- Success rate: %.2f%%", successRate)
	suite.T().Logf("- Errors: %d", errorCount)

	// Store metrics for reporting
	suite.mu.Lock()
	suite.metrics["load_test"] = map[string]interface{}{
		"total_messages":      totalMessages,
		"total_duration_ms":   totalTestDuration.Milliseconds(),
		"avg_message_time_ms": avgMessageTime.Milliseconds(),
		"messages_per_second": messagesPerSecond,
		"success_rate":        successRate,
		"error_count":         errorCount,
	}
	suite.mu.Unlock()

	// Performance assertions
	assert.True(suite.T(), avgMessageTime < 100*time.Millisecond, "Average message processing too slow: %v", avgMessageTime)
	assert.True(suite.T(), messagesPerSecond > 10, "Messages per second too low: %.2f", messagesPerSecond)
	assert.True(suite.T(), successRate > 95, "Success rate too low: %.2f%%", successRate)
	assert.Equal(suite.T(), int64(0), errorCount, "Unexpected errors during load test")
}

// Helper function to verify message integrity
func (suite *MessagingStressTestSuite) verifyMessageIntegrity(dialogID string, expectedCount int) {
	suite.mu.RLock()
	messages := suite.messages[dialogID]
	suite.mu.RUnlock()

	assert.Equal(suite.T(), expectedCount, len(messages), "Message count mismatch")

	// Verify each message has required fields
	for i, msg := range messages {
		assert.NotEmpty(suite.T(), msg["id"], "Message %d missing ID", i)
		assert.NotEmpty(suite.T(), msg["content"], "Message %d missing content", i)
		assert.NotEmpty(suite.T(), msg["sender_id"], "Message %d missing sender_id", i)
		assert.NotEmpty(suite.T(), msg["created_at"], "Message %d missing created_at", i)
		assert.Contains(suite.T(), msg, "message_number", "Message %d missing message_number", i)
	}
}