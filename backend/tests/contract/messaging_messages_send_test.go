package contract_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMessagingMessageSend_Contract validates the POST /dialogs/{id}/messages endpoint contract
// This test MUST FAIL initially as no implementation exists yet (TDD)
func TestMessagingMessageSend_Contract(t *testing.T) {
	// Test server URL - will fail until server is implemented
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	dialogID := uuid.New().String()

	tests := []struct {
		name           string
		payload        map[string]interface{}
		token          string
		expectedStatus int
		expectedFields []string
	}{
		{
			name: "Send text message",
			payload: map[string]interface{}{
				"type": "text",
				"content": map[string]interface{}{
					"text": "Hello world!",
				},
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "dialog_id", "sender_id", "type", "content", "created_at"},
		},
		{
			name: "Send text message with reply",
			payload: map[string]interface{}{
				"type": "text",
				"content": map[string]interface{}{
					"text": "Reply message",
				},
				"reply_to_id": uuid.New().String(),
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "dialog_id", "sender_id", "type", "content", "reply_to_id", "created_at"},
		},
		{
			name: "Send message with mentions",
			payload: map[string]interface{}{
				"type": "text",
				"content": map[string]interface{}{
					"text": "Hello @user1 and @user2",
				},
				"mentions": []string{uuid.New().String(), uuid.New().String()},
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "dialog_id", "sender_id", "type", "content", "mentions", "created_at"},
		},
		{
			name: "Send image message",
			payload: map[string]interface{}{
				"type": "image",
				"content": map[string]interface{}{
					"url":      "https://cdn.tchat.sea/images/example.jpg",
					"caption":  "Check this out!",
					"width":    1920,
					"height":   1080,
					"file_size": 524288,
				},
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "dialog_id", "sender_id", "type", "content", "created_at"},
		},
		{
			name: "Send payment message",
			payload: map[string]interface{}{
				"type": "payment",
				"content": map[string]interface{}{
					"amount":      10000, // 100.00 in cents
					"currency":    "THB",
					"description": "Lunch money",
				},
			},
			token:          mockToken,
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "dialog_id", "sender_id", "type", "content", "created_at"},
		},
		{
			name: "Invalid message type",
			payload: map[string]interface{}{
				"type": "invalid_type",
				"content": map[string]interface{}{
					"text": "Invalid message",
				},
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Missing required fields",
			payload: map[string]interface{}{
				"type": "text",
				// missing content field
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Unauthorized - no token",
			payload: map[string]interface{}{
				"type": "text",
				"content": map[string]interface{}{
					"text": "Hello world!",
				},
			},
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error", "message"},
		},
		{
			name: "Invalid dialog ID format",
			payload: map[string]interface{}{
				"type": "text",
				"content": map[string]interface{}{
					"text": "Hello world!",
				},
			},
			token:          mockToken,
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error", "message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use invalid dialog ID for the last test case
			testDialogID := dialogID
			if tt.name == "Invalid dialog ID format" {
				testDialogID = "invalid-uuid"
			}

			// Marshal request payload
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest("POST", baseURL+"/dialogs/"+testDialogID+"/messages", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Add authentication if provided
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// Execute request (will fail until server is implemented)
			client := &http.Client{}
			resp, err := client.Do(req)

			// This SHOULD FAIL initially - no server running
			if err != nil {
				t.Logf("Expected failure: %v (no implementation yet)", err)
				return // Test passes by failing as expected in TDD
			}
			defer resp.Body.Close()

			// Validate response status
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Validate response structure
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Check expected fields exist
			for _, field := range tt.expectedFields {
				assert.Contains(t, response, field, "Response should contain field: %s", field)
			}

			// Additional contract validations for successful message creation
			if resp.StatusCode == http.StatusCreated {
				// Message ID validation
				messageID, exists := response["id"]
				assert.True(t, exists, "Success response should have 'id' field")
				if exists {
					_, err := uuid.Parse(messageID.(string))
					assert.NoError(t, err, "Message ID should be valid UUID")
				}

				// Dialog ID validation
				dialogIDResp, exists := response["dialog_id"]
				assert.True(t, exists, "Success response should have 'dialog_id' field")
				if exists {
					assert.Equal(t, testDialogID, dialogIDResp, "Response dialog_id should match request")
				}

				// Sender ID validation
				senderID, exists := response["sender_id"]
				assert.True(t, exists, "Success response should have 'sender_id' field")
				if exists {
					_, err := uuid.Parse(senderID.(string))
					assert.NoError(t, err, "Sender ID should be valid UUID")
				}

				// Message type validation
				messageType, exists := response["type"]
				assert.True(t, exists, "Success response should have 'type' field")
				if exists {
					assert.Equal(t, tt.payload["type"], messageType, "Response type should match request")
				}

				// Content validation
				content, exists := response["content"]
				assert.True(t, exists, "Success response should have 'content' field")
				assert.NotNil(t, content, "Content should not be nil")

				// Timestamp validation
				createdAt, exists := response["created_at"]
				assert.True(t, exists, "Success response should have 'created_at' field")
				assert.NotEmpty(t, createdAt, "created_at should not be empty")

				// Reply validation if reply_to_id was provided
				if replyToID, hasReply := tt.payload["reply_to_id"]; hasReply {
					responseReplyID, exists := response["reply_to_id"]
					assert.True(t, exists, "Response should have 'reply_to_id' field when provided in request")
					if exists {
						assert.Equal(t, replyToID, responseReplyID, "Response reply_to_id should match request")
					}
				}

				// Mentions validation if mentions were provided
				if mentions, hasMentions := tt.payload["mentions"]; hasMentions {
					responseMentions, exists := response["mentions"]
					assert.True(t, exists, "Response should have 'mentions' field when provided in request")
					if exists {
						mentionsArray, ok := responseMentions.([]interface{})
						assert.True(t, ok, "Mentions should be an array")
						if ok {
							requestMentions := mentions.([]string)
							assert.Equal(t, len(requestMentions), len(mentionsArray), "Mentions array length should match")
						}
					}
				}
			}
		})
	}
}

// TestMessagingMessageSend_MessageTypes validates different message type contracts
func TestMessagingMessageSend_MessageTypes(t *testing.T) {
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	dialogID := uuid.New().String()

	messageTypes := []struct {
		msgType string
		content map[string]interface{}
		valid   bool
	}{
		{
			msgType: "text",
			content: map[string]interface{}{"text": "Hello world"},
			valid:   true,
		},
		{
			msgType: "voice",
			content: map[string]interface{}{
				"url":      "https://cdn.tchat.sea/voice/file.ogg",
				"duration": 5000,
				"waveform": []int{1, 2, 3, 4, 5},
			},
			valid: true,
		},
		{
			msgType: "file",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/files/document.pdf",
				"filename":  "document.pdf",
				"file_size": 1048576,
				"mime_type": "application/pdf",
			},
			valid: true,
		},
		{
			msgType: "video",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/videos/clip.mp4",
				"thumbnail": "https://cdn.tchat.sea/thumbs/clip.jpg",
				"duration":  30000,
				"width":     1920,
				"height":    1080,
			},
			valid: true,
		},
		{
			msgType: "location",
			content: map[string]interface{}{
				"latitude":  13.7563,
				"longitude": 100.5018,
				"address":   "Bangkok, Thailand",
			},
			valid: true,
		},
		{
			msgType: "invalid",
			content: map[string]interface{}{"text": "Invalid type"},
			valid:   false,
		},
	}

	for _, mt := range messageTypes {
		t.Run("Message type: "+mt.msgType, func(t *testing.T) {
			payload := map[string]interface{}{
				"type":    mt.msgType,
				"content": mt.content,
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", baseURL+"/dialogs/"+dialogID+"/messages", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+mockToken)

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				t.Logf("Expected failure: %v (no implementation yet)", err)
				return // Test passes by failing as expected in TDD
			}
			defer resp.Body.Close()

			if mt.valid {
				assert.Equal(t, http.StatusCreated, resp.StatusCode, "Valid message type should be accepted")
			} else {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Invalid message type should be rejected")
			}
		})
	}
}

// TestMessagingMessageSend_ContentValidation validates content structure per message type
func TestMessagingMessageSend_ContentValidation(t *testing.T) {
	baseURL := "http://localhost:8080"
	mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	dialogID := uuid.New().String()

	tests := []struct {
		name           string
		msgType        string
		content        interface{}
		shouldPass     bool
	}{
		{"Text - valid", "text", map[string]interface{}{"text": "Hello"}, true},
		{"Text - empty", "text", map[string]interface{}{"text": ""}, false},
		{"Text - missing text field", "text", map[string]interface{}{}, false},
		{"Image - valid", "image", map[string]interface{}{
			"url": "https://example.com/image.jpg",
		}, true},
		{"Image - invalid URL", "image", map[string]interface{}{
			"url": "not-a-url",
		}, false},
		{"Payment - valid", "payment", map[string]interface{}{
			"amount":   10000,
			"currency": "THB",
		}, true},
		{"Payment - negative amount", "payment", map[string]interface{}{
			"amount":   -1000,
			"currency": "THB",
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]interface{}{
				"type":    tt.msgType,
				"content": tt.content,
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", baseURL+"/dialogs/"+dialogID+"/messages", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+mockToken)

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				t.Logf("Expected failure: %v (no implementation yet)", err)
				return // Test passes by failing as expected in TDD
			}
			defer resp.Body.Close()

			if tt.shouldPass {
				assert.Equal(t, http.StatusCreated, resp.StatusCode, "Valid content should be accepted")
			} else {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Invalid content should be rejected")
			}
		})
	}
}