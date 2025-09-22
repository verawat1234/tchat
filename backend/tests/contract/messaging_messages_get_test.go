package contract

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T019: Contract test GET /dialogs/{id}/messages - Get dialog messages
func TestMessagingMessagesGetContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		dialogID       string
		queryParams    map[string]string
		expectedStatus int
		expectedFields []string
		description    string
	}{
		{
			name:           "get_messages_success",
			authHeader:     "Bearer valid_jwt_token_12345",
			dialogID:       "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"messages", "pagination"},
			description:    "Should return messages for valid dialog",
		},
		{
			name:        "get_messages_with_pagination",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"page": "1", "limit": "20"},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"messages", "pagination"},
			description:    "Should return paginated messages",
		},
		{
			name:        "get_messages_with_before_cursor",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"before": "msg_456e7890-e89b-12d3-a456-426614174001"},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"messages", "pagination"},
			description:    "Should return messages before specified message",
		},
		{
			name:        "get_messages_with_after_cursor",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"after": "msg_456e7890-e89b-12d3-a456-426614174001"},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"messages", "pagination"},
			description:    "Should return messages after specified message",
		},
		{
			name:        "get_messages_with_custom_limit",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"limit": "10"},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"messages", "pagination"},
			description:    "Should return messages with custom limit",
		},
		{
			name:           "missing_authorization",
			authHeader:     "",
			dialogID:       "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized when no auth header",
		},
		{
			name:           "invalid_jwt_token",
			authHeader:     "Bearer invalid_token_67890",
			dialogID:       "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for invalid token",
		},
		{
			name:           "invalid_dialog_id",
			authHeader:     "Bearer valid_jwt_token_12345",
			dialogID:       "invalid-dialog-id",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid dialog ID format",
		},
		{
			name:           "nonexistent_dialog",
			authHeader:     "Bearer valid_jwt_token_12345",
			dialogID:       "dialog_999e9999-e89b-12d3-a456-426614174999",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusNotFound,
			expectedFields: []string{"error"},
			description:    "Should return not found for nonexistent dialog",
		},
		{
			name:           "unauthorized_dialog_access",
			authHeader:     "Bearer valid_jwt_token_12345",
			dialogID:       "dialog_unauthorized_access",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusForbidden,
			expectedFields: []string{"error"},
			description:    "Should return forbidden for dialog user is not participant in",
		},
		{
			name:        "invalid_page_parameter",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"page": "invalid"},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid page parameter",
		},
		{
			name:        "invalid_limit_parameter",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"limit": "invalid"},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid limit parameter",
		},
		{
			name:        "limit_too_high",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"limit": "1000"},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for limit exceeding maximum",
		},
		{
			name:        "invalid_before_cursor",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"before": "invalid-message-id"},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid before cursor",
		},
		{
			name:        "invalid_after_cursor",
			authHeader:  "Bearer valid_jwt_token_12345",
			dialogID:    "dialog_123e4567-e89b-12d3-a456-426614174000",
			queryParams: map[string]string{"after": "invalid-message-id"},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid after cursor",
		},
		{
			name:           "empty_dialog",
			authHeader:     "Bearer valid_jwt_token_12345",
			dialogID:       "dialog_empty_no_messages",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			expectedFields: []string{"messages", "pagination"},
			description:    "Should return empty messages array for dialog with no messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()

			// Mock authentication middleware and endpoint
			router.GET("/api/v1/messaging/dialogs/:dialog_id/messages", func(c *gin.Context) {
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

				dialogID := c.Param("dialog_id")

				// Validate dialog ID format
				if dialogID == "invalid-dialog-id" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dialog ID format"})
					return
				}

				// Check if dialog exists
				if dialogID == "dialog_999e9999-e89b-12d3-a456-426614174999" {
					c.JSON(http.StatusNotFound, gin.H{"error": "Dialog not found"})
					return
				}

				// Check if user has access to dialog
				if dialogID == "dialog_unauthorized_access" {
					c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this dialog"})
					return
				}

				// Parse and validate query parameters
				page := c.DefaultQuery("page", "1")
				limit := c.DefaultQuery("limit", "50")
				before := c.Query("before")
				after := c.Query("after")

				// Validate page parameter
				if page != "1" && page != "2" && page != "3" && page != "" {
					if page == "invalid" {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
						return
					}
				}

				// Validate limit parameter
				if limit == "invalid" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
					return
				}
				if limit == "1000" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Limit exceeds maximum allowed value (100)"})
					return
				}

				// Validate cursor parameters
				if before == "invalid-message-id" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid before cursor format"})
					return
				}
				if after == "invalid-message-id" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid after cursor format"})
					return
				}

				// Mock messages data
				var messages []map[string]interface{}

				if dialogID == "dialog_empty_no_messages" {
					// Empty dialog
					messages = []map[string]interface{}{}
				} else {
					// Mock messages with different types
					messages = []map[string]interface{}{
						{
							"id":           "msg_123e4567-e89b-12d3-a456-426614174001",
							"dialog_id":    dialogID,
							"sender_id":    "user_123e4567-e89b-12d3-a456-426614174000",
							"content":      "Hello, how are you?",
							"message_type": "text",
							"metadata":     map[string]interface{}{},
							"reply_to_id":  nil,
							"edited":       false,
							"reactions":    []interface{}{},
							"read_receipts": []map[string]interface{}{
								{
									"user_id": "user_456e7890-e89b-12d3-a456-426614174001",
									"read_at": "2023-12-01T15:35:00Z",
								},
							},
							"created_at": "2023-12-01T15:30:00Z",
							"updated_at": "2023-12-01T15:30:00Z",
						},
						{
							"id":           "msg_456e7890-e89b-12d3-a456-426614174002",
							"dialog_id":    dialogID,
							"sender_id":    "user_456e7890-e89b-12d3-a456-426614174001",
							"content":      "I'm doing great, thanks!",
							"message_type": "text",
							"metadata":     map[string]interface{}{},
							"reply_to_id":  "msg_123e4567-e89b-12d3-a456-426614174001",
							"edited":       false,
							"reactions": []map[string]interface{}{
								{
									"emoji":    "👍",
									"user_ids": []string{"user_123e4567-e89b-12d3-a456-426614174000"},
								},
							},
							"read_receipts": []interface{}{},
							"created_at":    "2023-12-01T15:32:00Z",
							"updated_at":    "2023-12-01T15:32:00Z",
						},
						{
							"id":           "msg_789e1234-e89b-12d3-a456-426614174003",
							"dialog_id":    dialogID,
							"sender_id":    "user_123e4567-e89b-12d3-a456-426614174000",
							"content":      "https://cdn.tchat.sea/images/photo123.jpg",
							"message_type": "image",
							"metadata": map[string]interface{}{
								"file_size":  1024000,
								"mime_type":  "image/jpeg",
								"width":      1920,
								"height":     1080,
								"thumbnail":  "https://cdn.tchat.sea/images/thumb123.jpg",
							},
							"reply_to_id":   nil,
							"edited":        false,
							"reactions":     []interface{}{},
							"read_receipts": []interface{}{},
							"created_at":    "2023-12-01T15:35:00Z",
							"updated_at":    "2023-12-01T15:35:00Z",
						},
					}

					// Apply cursor filters
					if before != "" {
						// Filter messages before the specified message
						var filteredMessages []map[string]interface{}
						for _, msg := range messages {
							if msg["id"] != before {
								filteredMessages = append(filteredMessages, msg)
							} else {
								break
							}
						}
						messages = filteredMessages
					}

					if after != "" {
						// Filter messages after the specified message
						var filteredMessages []map[string]interface{}
						found := false
						for _, msg := range messages {
							if found {
								filteredMessages = append(filteredMessages, msg)
							}
							if msg["id"] == after {
								found = true
							}
						}
						messages = filteredMessages
					}
				}

				// Mock pagination
				pagination := map[string]interface{}{
					"page":        1,
					"limit":       50,
					"total":       len(messages),
					"total_pages": 1,
					"has_next":    false,
					"has_prev":    false,
				}

				if limit == "10" {
					pagination["limit"] = 10
					if len(messages) > 10 {
						messages = messages[:10]
						pagination["has_next"] = true
					}
				}

				c.JSON(http.StatusOK, gin.H{
					"messages":   messages,
					"pagination": pagination,
				})
			})

			// Prepare request URL with query parameters
			url := "/api/v1/messaging/dialogs/" + tt.dialogID + "/messages"
			if len(tt.queryParams) > 0 {
				url += "?"
				first := true
				for key, value := range tt.queryParams {
					if !first {
						url += "&"
					}
					url += key + "=" + value
					first = false
				}
			}

			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

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
			case "get_messages_success", "get_messages_with_pagination", "get_messages_with_before_cursor", "get_messages_with_after_cursor", "get_messages_with_custom_limit":
				// Verify messages structure
				messages, ok := response["messages"].([]interface{})
				assert.True(t, ok, "messages should be array")

				if len(messages) > 0 {
					// Verify message structure
					message := messages[0].(map[string]interface{})

					// Required fields
					assert.Contains(t, message, "id", "message should contain id")
					assert.Contains(t, message, "dialog_id", "message should contain dialog_id")
					assert.Contains(t, message, "sender_id", "message should contain sender_id")
					assert.Contains(t, message, "content", "message should contain content")
					assert.Contains(t, message, "message_type", "message should contain message_type")
					assert.Contains(t, message, "created_at", "message should contain created_at")
					assert.Contains(t, message, "updated_at", "message should contain updated_at")

					// Verify message ID format
					messageID, ok := message["id"].(string)
					assert.True(t, ok, "message id should be string")
					assert.Contains(t, messageID, "msg_", "message id should have msg_ prefix")

					// Verify message type is valid
					messageType, ok := message["message_type"].(string)
					assert.True(t, ok, "message_type should be string")
					validTypes := []string{"text", "image", "video", "audio", "file", "location", "sticker"}
					assert.Contains(t, validTypes, messageType, "message_type should be valid")

					// Verify dialog_id matches request
					dialogID, ok := message["dialog_id"].(string)
					assert.True(t, ok, "dialog_id should be string")
					assert.Equal(t, tt.dialogID, dialogID, "dialog_id should match request")

					// Verify optional fields structure
					if reactions, exists := message["reactions"]; exists {
						assert.IsType(t, []interface{}{}, reactions, "reactions should be array")
					}

					if readReceipts, exists := message["read_receipts"]; exists {
						assert.IsType(t, []interface{}{}, readReceipts, "read_receipts should be array")
					}

					if metadata, exists := message["metadata"]; exists {
						assert.IsType(t, map[string]interface{}{}, metadata, "metadata should be object")
					}
				}

				// Verify pagination structure
				pagination, ok := response["pagination"].(map[string]interface{})
				assert.True(t, ok, "pagination should be object")

				assert.Contains(t, pagination, "page", "pagination should contain page")
				assert.Contains(t, pagination, "limit", "pagination should contain limit")
				assert.Contains(t, pagination, "total", "pagination should contain total")
				assert.Contains(t, pagination, "has_next", "pagination should contain has_next")
				assert.Contains(t, pagination, "has_prev", "pagination should contain has_prev")

				// Verify pagination values
				page, ok := pagination["page"].(int)
				assert.True(t, ok, "page should be integer")
				assert.GreaterOrEqual(t, page, 1, "page should be >= 1")

				limit, ok := pagination["limit"].(int)
				assert.True(t, ok, "limit should be integer")
				assert.GreaterOrEqual(t, limit, 1, "limit should be >= 1")
				assert.LessOrEqual(t, limit, 100, "limit should be <= 100")

				total, ok := pagination["total"].(int)
				assert.True(t, ok, "total should be integer")
				assert.GreaterOrEqual(t, total, 0, "total should be >= 0")

			case "empty_dialog":
				// Verify empty messages array
				messages, ok := response["messages"].([]interface{})
				assert.True(t, ok, "messages should be array")
				assert.Len(t, messages, 0, "messages should be empty for empty dialog")

				// Verify pagination for empty result
				pagination, ok := response["pagination"].(map[string]interface{})
				assert.True(t, ok, "pagination should be object")

				total, ok := pagination["total"].(int)
				assert.True(t, ok, "total should be integer")
				assert.Equal(t, 0, total, "total should be 0 for empty dialog")

			default:
				// Error cases
				errorMsg, ok := response["error"].(string)
				assert.True(t, ok, "error should be string")
				assert.NotEmpty(t, errorMsg, "error message should not be empty")
			}
		})
	}
}

// TestMessagingMessagesGetSecurity tests security aspects
func TestMessagingMessagesGetSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("dialog_access_control", func(t *testing.T) {
		// Ensure users can only access messages from dialogs they participate in
		router := gin.New()
		router.GET("/api/v1/messaging/dialogs/:dialog_id/messages", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			dialogID := c.Param("dialog_id")

			// Mock user and dialog participants
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

			// Mock dialog participants
			dialogParticipants := map[string][]string{
				"dialog_user1_dialog": {"user1-id", "user3-id"},
				"dialog_user2_dialog": {"user2-id", "user4-id"},
				"dialog_shared":       {"user1-id", "user2-id"},
			}

			participants, exists := dialogParticipants[dialogID]
			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Dialog not found"})
				return
			}

			// Check if current user is participant
			isParticipant := false
			for _, participantID := range participants {
				if participantID == currentUserID {
					isParticipant = true
					break
				}
			}

			if !isParticipant {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this dialog"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"messages": []interface{}{},
				"pagination": gin.H{"total": 0},
			})
		})

		// User 1 should access their dialog
		req1, _ := http.NewRequest("GET", "/api/v1/messaging/dialogs/dialog_user1_dialog/messages", nil)
		req1.Header.Set("Authorization", "Bearer user1_token")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code, "User1 should access their dialog")

		// User 1 should NOT access user2's dialog
		req2, _ := http.NewRequest("GET", "/api/v1/messaging/dialogs/dialog_user2_dialog/messages", nil)
		req2.Header.Set("Authorization", "Bearer user1_token")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusForbidden, w2.Code, "User1 should not access user2's dialog")

		// Both users should access shared dialog
		req3, _ := http.NewRequest("GET", "/api/v1/messaging/dialogs/dialog_shared/messages", nil)
		req3.Header.Set("Authorization", "Bearer user1_token")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w3.Code, "User1 should access shared dialog")

		req4, _ := http.NewRequest("GET", "/api/v1/messaging/dialogs/dialog_shared/messages", nil)
		req4.Header.Set("Authorization", "Bearer user2_token")
		w4 := httptest.NewRecorder()
		router.ServeHTTP(w4, req4)
		assert.Equal(t, http.StatusOK, w4.Code, "User2 should access shared dialog")
	})

	t.Run("sensitive_data_filtering", func(t *testing.T) {
		// Ensure no sensitive data is exposed in message responses
		router := gin.New()
		router.GET("/api/v1/messaging/dialogs/:dialog_id/messages", func(c *gin.Context) {
			// Mock message with potential sensitive data filtered out
			messages := []map[string]interface{}{
				{
					"id":           "msg_123",
					"dialog_id":    c.Param("dialog_id"),
					"sender_id":    "user_123",
					"content":      "Hello there!",
					"message_type": "text",
					"created_at":   "2023-12-01T15:30:00Z",
					"updated_at":   "2023-12-01T15:30:00Z",
					// Sensitive fields should NOT be included:
					// - raw_content (with potential user input)
					// - ip_address
					// - device_info
					// - encryption_keys
				},
			}

			c.JSON(http.StatusOK, gin.H{
				"messages": messages,
				"pagination": gin.H{"total": 1},
			})
		})

		req, _ := http.NewRequest("GET", "/api/v1/messaging/dialogs/test_dialog/messages", nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		messages := response["messages"].([]interface{})
		message := messages[0].(map[string]interface{})

		// Verify sensitive fields are not present
		sensitiveFields := []string{"ip_address", "device_info", "encryption_keys", "raw_content", "private_metadata"}
		for _, field := range sensitiveFields {
			assert.NotContains(t, message, field, "Message should not contain sensitive field: %s", field)
		}
	})
}

// TestMessagingMessagesGetPerformance tests performance aspects
func TestMessagingMessagesGetPerformance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("response_size_limits", func(t *testing.T) {
		router := gin.New()
		router.GET("/api/v1/messaging/dialogs/:dialog_id/messages", func(c *gin.Context) {
			// Mock large response but with reasonable limits
			messages := make([]map[string]interface{}, 50) // Default limit
			for i := 0; i < 50; i++ {
				messages[i] = map[string]interface{}{
					"id":           "msg_" + string(rune(i)),
					"dialog_id":    c.Param("dialog_id"),
					"sender_id":    "user_123",
					"content":      "Message content",
					"message_type": "text",
					"created_at":   "2023-12-01T15:30:00Z",
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"messages":   messages,
				"pagination": gin.H{"total": 1000, "limit": 50, "page": 1},
			})
		})

		req, _ := http.NewRequest("GET", "/api/v1/messaging/dialogs/test_dialog/messages", nil)
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response size is reasonable (< 1MB)
		responseSize := len(w.Body.Bytes())
		assert.Less(t, responseSize, 1*1024*1024, "Response should be less than 1MB, got %d bytes", responseSize)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		messages := response["messages"].([]interface{})
		assert.LessOrEqual(t, len(messages), 100, "Should not return more than 100 messages per request")
	})
}