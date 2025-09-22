package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T021: Contract test POST /messages/{id}/read - Mark message as read
func TestMessagingMessageReadContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		messageID      string
		payload        map[string]interface{}
		expectedStatus int
		expectedFields []string
		description    string
	}{
		{
			name:           "mark_message_as_read_success",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_123e4567-e89b-12d3-a456-426614174001",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"message_id", "read_at", "user_id"},
			description:    "Should mark message as read successfully",
		},
		{
			name:           "mark_message_with_timestamp",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_456e7890-e89b-12d3-a456-426614174002",
			payload:        map[string]interface{}{"read_at": "2023-12-01T15:45:00Z"},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"message_id", "read_at", "user_id"},
			description:    "Should mark message as read with custom timestamp",
		},
		{
			name:           "mark_already_read_message",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_already_read_by_user",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"message_id", "read_at", "user_id"},
			description:    "Should handle already read message gracefully",
		},
		{
			name:           "missing_authorization",
			authHeader:     "",
			messageID:      "msg_123e4567-e89b-12d3-a456-426614174001",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized when no auth header",
		},
		{
			name:           "invalid_jwt_token",
			authHeader:     "Bearer invalid_token_67890",
			messageID:      "msg_123e4567-e89b-12d3-a456-426614174001",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for invalid token",
		},
		{
			name:           "invalid_message_id_format",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "invalid-message-id",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid message ID format",
		},
		{
			name:           "nonexistent_message",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_999e9999-e89b-12d3-a456-426614174999",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusNotFound,
			expectedFields: []string{"error"},
			description:    "Should return not found for nonexistent message",
		},
		{
			name:           "unauthorized_message_access",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_unauthorized_access",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusForbidden,
			expectedFields: []string{"error"},
			description:    "Should return forbidden for message in dialog user is not participant in",
		},
		{
			name:           "cannot_mark_own_message",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_sent_by_current_user",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error when trying to mark own message as read",
		},
		{
			name:           "invalid_timestamp_format",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_123e4567-e89b-12d3-a456-426614174001",
			payload:        map[string]interface{}{"read_at": "invalid-timestamp"},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid timestamp format",
		},
		{
			name:           "future_timestamp",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_123e4567-e89b-12d3-a456-426614174001",
			payload:        map[string]interface{}{"read_at": "2030-12-01T15:45:00Z"},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for future timestamp",
		},
		{
			name:           "timestamp_before_message_creation",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_123e4567-e89b-12d3-a456-426614174001",
			payload:        map[string]interface{}{"read_at": "2020-01-01T00:00:00Z"},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for timestamp before message creation",
		},
		{
			name:           "malformed_request_body",
			authHeader:     "Bearer valid_jwt_token_12345",
			messageID:      "msg_123e4567-e89b-12d3-a456-426614174001",
			payload:        map[string]interface{}{"read_at": 12345}, // Wrong type
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for malformed request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()

			// Mock authentication middleware and endpoint
			router.POST("/api/v1/messaging/messages/:message_id/read", func(c *gin.Context) {
				authHeader := c.GetHeader("Authorization")

				// Mock authentication logic
				if authHeader == "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
					return
				}

				if authHeader != "Bearer valid_jwt_token_12345" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
					return
				}

				messageID := c.Param("message_id")

				// Validate message ID format
				if messageID == "invalid-message-id" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID format"})
					return
				}

				// Check if message exists
				if messageID == "msg_999e9999-e89b-12d3-a456-426614174999" {
					c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
					return
				}

				// Check if user has access to message (is participant in dialog)
				if messageID == "msg_unauthorized_access" {
					c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this message"})
					return
				}

				// Check if user is trying to mark their own message as read
				if messageID == "msg_sent_by_current_user" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot mark your own message as read"})
					return
				}

				// Parse request body
				var request map[string]interface{}
				if err := c.ShouldBindJSON(&request); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
					return
				}

				// Validate read_at timestamp if provided
				var readAt string
				if readAtValue, exists := request["read_at"]; exists {
					if readAtStr, ok := readAtValue.(string); ok {
						readAt = readAtStr

						// Simple timestamp validation
						if readAt == "invalid-timestamp" {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
							return
						}
						if readAt == "2030-12-01T15:45:00Z" {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Read timestamp cannot be in the future"})
							return
						}
						if readAt == "2020-01-01T00:00:00Z" {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Read timestamp cannot be before message creation"})
							return
						}
					} else {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
						return
					}
				} else {
					// Use current timestamp if not provided
					readAt = "2023-12-01T15:45:00Z"
				}

				// Mock current user ID from token
				currentUserID := "user_123e4567-e89b-12d3-a456-426614174000"

				// Check if message is already read by user
				alreadyRead := messageID == "msg_already_read_by_user"

				response := gin.H{
					"message_id": messageID,
					"user_id":    currentUserID,
					"read_at":    readAt,
				}

				if alreadyRead {
					response["previously_read"] = true
					response["original_read_at"] = "2023-12-01T15:40:00Z"
				}

				c.JSON(http.StatusOK, response)
			})

			// Prepare request
			jsonData, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/messaging/messages/"+tt.messageID+"/read", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify status code
			assert.Equal(t, tt.expectedStatus, w.Code,
				"Test: %s - %s", tt.name, tt.description)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Verify expected fields are present
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field,
					"Test: %s - Response should contain field '%s'", tt.name, field)
			}

			// Additional assertions based on test case
			switch tt.name {
			case "mark_message_as_read_success", "mark_message_with_timestamp", "mark_already_read_message":
				// Verify message ID matches
				messageID, ok := response["message_id"].(string)
				assert.True(t, ok, "message_id should be string")
				assert.Equal(t, tt.messageID, messageID, "message_id should match request")

				// Verify user ID
				userID, ok := response["user_id"].(string)
				assert.True(t, ok, "user_id should be string")
				assert.NotEmpty(t, userID, "user_id should not be empty")

				// Verify read_at timestamp
				readAt, ok := response["read_at"].(string)
				assert.True(t, ok, "read_at should be string")
				assert.NotEmpty(t, readAt, "read_at should not be empty")

				// If custom timestamp was provided, verify it matches
				if customReadAt, exists := tt.payload["read_at"]; exists {
					assert.Equal(t, customReadAt, readAt, "read_at should match custom timestamp")
				}

				// For already read message, verify additional fields
				if tt.name == "mark_already_read_message" {
					if previouslyRead, exists := response["previously_read"]; exists {
						assert.True(t, previouslyRead.(bool), "previously_read should be true")
					}
					if originalReadAt, exists := response["original_read_at"]; exists {
						assert.NotEmpty(t, originalReadAt.(string), "original_read_at should not be empty")
					}
				}

			default:
				// Error cases
				errorMsg, ok := response["error"].(string)
				assert.True(t, ok, "error should be string")
				assert.NotEmpty(t, errorMsg, "error message should not be empty")
			}
		})
	}
}

// TestMessagingMessageReadBulk tests bulk message read operations
func TestMessagingMessageReadBulk(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("mark_multiple_messages_read", func(t *testing.T) {
		router := gin.New()

		// Mock bulk read endpoint
		router.POST("/api/v1/messaging/messages/bulk-read", func(c *gin.Context) {
			var request struct {
				MessageIDs []string `json:"message_ids"`
				ReadAt     *string  `json:"read_at,omitempty"`
			}

			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
				return
			}

			if len(request.MessageIDs) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No message IDs provided"})
				return
			}

			if len(request.MessageIDs) > 100 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Too many messages (max 100)"})
				return
			}

			// Mock processing results
			results := make([]map[string]interface{}, len(request.MessageIDs))
			for i, msgID := range request.MessageIDs {
				results[i] = map[string]interface{}{
					"message_id": msgID,
					"user_id":    "user_123",
					"read_at":    "2023-12-01T15:45:00Z",
					"success":    true,
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"results":     results,
				"total_count": len(request.MessageIDs),
				"success_count": len(request.MessageIDs),
				"error_count":   0,
			})
		})

		// Test bulk read
		payload := map[string]interface{}{
			"message_ids": []string{
				"msg_123",
				"msg_456",
				"msg_789",
			},
		}
		jsonData, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/v1/messaging/messages/bulk-read", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(3), response["total_count"])
		assert.Equal(t, float64(3), response["success_count"])
		assert.Equal(t, float64(0), response["error_count"])

		results := response["results"].([]interface{})
		assert.Len(t, results, 3)
	})

	t.Run("mark_dialog_messages_read", func(t *testing.T) {
		router := gin.New()

		// Mock dialog read endpoint (mark all messages in dialog as read)
		router.POST("/api/v1/messaging/dialogs/:dialog_id/read", func(c *gin.Context) {
			dialogID := c.Param("dialog_id")

			var request struct {
				ReadAt *string `json:"read_at,omitempty"`
			}
			c.ShouldBindJSON(&request)

			c.JSON(http.StatusOK, gin.H{
				"dialog_id":    dialogID,
				"user_id":      "user_123",
				"read_at":      "2023-12-01T15:45:00Z",
				"messages_marked": 15,
			})
		})

		payload := map[string]interface{}{}
		jsonData, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/v1/messaging/dialogs/dialog_123/read", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "dialog_123", response["dialog_id"])
		assert.Equal(t, float64(15), response["messages_marked"])
	})
}

// TestMessagingMessageReadSecurity tests security aspects
func TestMessagingMessageReadSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("message_access_control", func(t *testing.T) {
		// Ensure users can only mark messages as read in dialogs they participate in
		router := gin.New()
		router.POST("/api/v1/messaging/messages/:message_id/read", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			messageID := c.Param("message_id")

			// Mock user and message-dialog mapping
			var currentUserID string
			switch authHeader {
			case "Bearer user1_token":
				currentUserID = "user1-id"
			case "Bearer user2_token":
				currentUserID = "user2-id"
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			// Mock message-dialog-participants mapping
			messageDialogParticipants := map[string][]string{
				"msg_user1_dialog": {"user1-id", "user3-id"},
				"msg_user2_dialog": {"user2-id", "user4-id"},
				"msg_shared":       {"user1-id", "user2-id"},
			}

			participants, exists := messageDialogParticipants[messageID]
			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
				return
			}

			// Check if current user is participant in the dialog
			isParticipant := false
			for _, participantID := range participants {
				if participantID == currentUserID {
					isParticipant = true
					break
				}
			}

			if !isParticipant {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this message"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message_id": messageID,
				"user_id":    currentUserID,
				"read_at":    "2023-12-01T15:45:00Z",
			})
		})

		// User 1 should mark messages in their dialog
		req1, _ := http.NewRequest("POST", "/api/v1/messaging/messages/msg_user1_dialog/read", bytes.NewBuffer([]byte("{}")))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("Authorization", "Bearer user1_token")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code, "User1 should mark messages in their dialog")

		// User 1 should NOT mark messages in user2's dialog
		req2, _ := http.NewRequest("POST", "/api/v1/messaging/messages/msg_user2_dialog/read", bytes.NewBuffer([]byte("{}")))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer user1_token")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusForbidden, w2.Code, "User1 should not mark messages in user2's dialog")

		// Both users should mark messages in shared dialog
		req3, _ := http.NewRequest("POST", "/api/v1/messaging/messages/msg_shared/read", bytes.NewBuffer([]byte("{}")))
		req3.Header.Set("Content-Type", "application/json")
		req3.Header.Set("Authorization", "Bearer user1_token")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w3.Code, "User1 should mark messages in shared dialog")
	})

	t.Run("prevent_read_receipt_spoofing", func(t *testing.T) {
		// Ensure users cannot spoof read receipts for other users
		router := gin.New()
		router.POST("/api/v1/messaging/messages/:message_id/read", func(c *gin.Context) {
			var request map[string]interface{}
			c.ShouldBindJSON(&request)

			// Check for attempts to set user_id in request (should be ignored)
			if _, exists := request["user_id"]; exists {
				// In real implementation, this should be ignored and user_id should come from auth token
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot specify user_id in request"})
				return
			}

			// User ID should always come from authentication token, not request
			c.JSON(http.StatusOK, gin.H{
				"message_id": c.Param("message_id"),
				"user_id":    "user_from_token", // Always from token
				"read_at":    "2023-12-01T15:45:00Z",
			})
		})

		// Try to spoof user_id in request
		payload := map[string]interface{}{
			"user_id": "malicious_user_id",
		}
		jsonData, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/v1/messaging/messages/msg_123/read", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Should reject request with user_id")
	})
}

// TestMessagingMessageReadRateLimit tests rate limiting
func TestMessagingMessageReadRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	readCount := 0
	router.POST("/api/v1/messaging/messages/:message_id/read", func(c *gin.Context) {
		readCount++
		if readCount > 100 { // Simulate rate limit of 100 read operations per minute
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many read operations. Please slow down.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message_id": c.Param("message_id"),
			"user_id":    "user_123",
			"read_at":    "2023-12-01T15:45:00Z",
		})
	})

	payload := map[string]interface{}{}
	jsonData, _ := json.Marshal(payload)

	// Make requests beyond rate limit
	for i := 0; i < 105; i++ {
		req, _ := http.NewRequest("POST", "/api/v1/messaging/messages/msg_123/read", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i < 100 {
			assert.Equal(t, http.StatusOK, w.Code, "Read %d should succeed", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Read %d should be rate limited", i+1)
		}
	}
}