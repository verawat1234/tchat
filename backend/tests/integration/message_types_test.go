package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// T037: Integration test for all 9 message types
// Tests comprehensive message type support, validation, and processing
type MessageTypesTestSuite struct {
	suite.Suite
	router         *gin.Engine
	testDialogID   string
	testSenderID   string
	messageStorage map[string]interface{} // Simulated message storage
}

func TestMessageTypesSuite(t *testing.T) {
	suite.Run(t, new(MessageTypesTestSuite))
}

func (suite *MessageTypesTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.messageStorage = make(map[string]interface{})
	suite.testDialogID = uuid.New().String()
	suite.testSenderID = uuid.New().String()

	suite.setupMessageTypeEndpoints()
}

func (suite *MessageTypesTestSuite) setupMessageTypeEndpoints() {
	// Message sending endpoint with comprehensive validation
	suite.router.POST("/dialogs/:dialogId/messages", func(c *gin.Context) {
		dialogId := c.Param("dialogId")

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		messageType, exists := req["type"]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message type is required"})
			return
		}

		content, contentExists := req["content"]
		if !contentExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
			return
		}

		// Validate message type and content
		if err := suite.validateMessageTypeAndContent(messageType.(string), content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create message ID and store
		messageId := uuid.New().String()
		messageData := map[string]interface{}{
			"id":         messageId,
			"dialog_id":  dialogId,
			"sender_id":  suite.testSenderID,
			"type":       messageType,
			"content":    content,
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"is_edited":  false,
			"is_deleted": false,
		}

		suite.messageStorage[messageId] = messageData

		c.JSON(http.StatusCreated, gin.H{
			"message_id": messageId,
			"status":     "sent",
			"type":       messageType,
		})
	})

	// Message retrieval endpoint
	suite.router.GET("/dialogs/:dialogId/messages/:messageId", func(c *gin.Context) {
		messageId := c.Param("messageId")

		message, exists := suite.messageStorage[messageId]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
			return
		}

		c.JSON(http.StatusOK, message)
	})

	// Message validation endpoint
	suite.router.POST("/messages/validate", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		messageType := req["type"].(string)
		content := req["content"]

		if err := suite.validateMessageTypeAndContent(messageType, content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"valid": false,
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"valid": true,
			"type":  messageType,
		})
	})
}

// validateMessageTypeAndContent performs comprehensive validation for all 9 message types
func (suite *MessageTypesTestSuite) validateMessageTypeAndContent(messageType string, content interface{}) error {
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return fmt.Errorf("content must be an object")
	}

	switch messageType {
	case "text":
		return suite.validateTextContent(contentMap)
	case "voice":
		return suite.validateVoiceContent(contentMap)
	case "file":
		return suite.validateFileContent(contentMap)
	case "image":
		return suite.validateImageContent(contentMap)
	case "video":
		return suite.validateVideoContent(contentMap)
	case "payment":
		return suite.validatePaymentContent(contentMap)
	case "location":
		return suite.validateLocationContent(contentMap)
	case "sticker":
		return suite.validateStickerContent(contentMap)
	case "system":
		return suite.validateSystemContent(contentMap)
	default:
		return fmt.Errorf("unsupported message type: %s", messageType)
	}
}

// Text message validation
func (suite *MessageTypesTestSuite) validateTextContent(content map[string]interface{}) error {
	text, exists := content["text"]
	if !exists {
		return fmt.Errorf("text content must have 'text' field")
	}

	textStr, ok := text.(string)
	if !ok {
		return fmt.Errorf("text field must be a string")
	}

	if len(textStr) == 0 {
		return fmt.Errorf("text content cannot be empty")
	}

	if len(textStr) > 4096 {
		return fmt.Errorf("text content cannot exceed 4096 characters")
	}

	return nil
}

// Voice message validation
func (suite *MessageTypesTestSuite) validateVoiceContent(content map[string]interface{}) error {
	requiredFields := []string{"url", "duration", "file_size"}
	for _, field := range requiredFields {
		if _, exists := content[field]; !exists {
			return fmt.Errorf("voice content must have '%s' field", field)
		}
	}

	if duration, ok := content["duration"].(float64); ok {
		if duration <= 0 || duration > 300000 { // Max 5 minutes
			return fmt.Errorf("voice duration must be between 1ms and 5 minutes")
		}
	}

	if fileSize, ok := content["file_size"].(float64); ok {
		if fileSize <= 0 || fileSize > 50*1024*1024 { // Max 50MB
			return fmt.Errorf("voice file size must be between 1 byte and 50MB")
		}
	}

	return nil
}

// File message validation
func (suite *MessageTypesTestSuite) validateFileContent(content map[string]interface{}) error {
	requiredFields := []string{"url", "filename", "file_size", "mime_type"}
	for _, field := range requiredFields {
		if _, exists := content[field]; !exists {
			return fmt.Errorf("file content must have '%s' field", field)
		}
	}

	if fileSize, ok := content["file_size"].(float64); ok {
		if fileSize <= 0 || fileSize > 100*1024*1024 { // Max 100MB
			return fmt.Errorf("file size must be between 1 byte and 100MB")
		}
	}

	return nil
}

// Image message validation
func (suite *MessageTypesTestSuite) validateImageContent(content map[string]interface{}) error {
	requiredFields := []string{"url", "width", "height", "file_size"}
	for _, field := range requiredFields {
		if _, exists := content[field]; !exists {
			return fmt.Errorf("image content must have '%s' field", field)
		}
	}

	if width, ok := content["width"].(float64); ok {
		if width <= 0 || width > 8192 {
			return fmt.Errorf("image width must be between 1 and 8192 pixels")
		}
	}

	if height, ok := content["height"].(float64); ok {
		if height <= 0 || height > 8192 {
			return fmt.Errorf("image height must be between 1 and 8192 pixels")
		}
	}

	return nil
}

// Video message validation
func (suite *MessageTypesTestSuite) validateVideoContent(content map[string]interface{}) error {
	requiredFields := []string{"url", "duration", "width", "height", "file_size"}
	for _, field := range requiredFields {
		if _, exists := content[field]; !exists {
			return fmt.Errorf("video content must have '%s' field", field)
		}
	}

	if duration, ok := content["duration"].(float64); ok {
		if duration <= 0 || duration > 1800000 { // Max 30 minutes
			return fmt.Errorf("video duration must be between 1ms and 30 minutes")
		}
	}

	if fileSize, ok := content["file_size"].(float64); ok {
		if fileSize <= 0 || fileSize > 500*1024*1024 { // Max 500MB
			return fmt.Errorf("video file size must be between 1 byte and 500MB")
		}
	}

	return nil
}

// Payment message validation
func (suite *MessageTypesTestSuite) validatePaymentContent(content map[string]interface{}) error {
	requiredFields := []string{"amount", "currency", "description", "status"}
	for _, field := range requiredFields {
		if _, exists := content[field]; !exists {
			return fmt.Errorf("payment content must have '%s' field", field)
		}
	}

	if amount, ok := content["amount"].(float64); ok {
		if amount <= 0 {
			return fmt.Errorf("payment amount must be positive")
		}
	}

	validCurrencies := []string{"THB", "SGD", "IDR", "MYR", "PHP", "VND", "USD"}
	if currency, ok := content["currency"].(string); ok {
		valid := false
		for _, validCurrency := range validCurrencies {
			if currency == validCurrency {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid currency: %s", currency)
		}
	}

	validStatuses := []string{"pending", "completed", "failed", "cancelled"}
	if status, ok := content["status"].(string); ok {
		valid := false
		for _, validStatus := range validStatuses {
			if status == validStatus {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid payment status: %s", status)
		}
	}

	return nil
}

// Location message validation
func (suite *MessageTypesTestSuite) validateLocationContent(content map[string]interface{}) error {
	requiredFields := []string{"latitude", "longitude"}
	for _, field := range requiredFields {
		if _, exists := content[field]; !exists {
			return fmt.Errorf("location content must have '%s' field", field)
		}
	}

	if lat, ok := content["latitude"].(float64); ok {
		if lat < -90 || lat > 90 {
			return fmt.Errorf("latitude must be between -90 and 90")
		}
	}

	if lng, ok := content["longitude"].(float64); ok {
		if lng < -180 || lng > 180 {
			return fmt.Errorf("longitude must be between -180 and 180")
		}
	}

	return nil
}

// Sticker message validation
func (suite *MessageTypesTestSuite) validateStickerContent(content map[string]interface{}) error {
	requiredFields := []string{"sticker_id", "pack_id", "url"}
	for _, field := range requiredFields {
		if _, exists := content[field]; !exists {
			return fmt.Errorf("sticker content must have '%s' field", field)
		}
	}

	return nil
}

// System message validation
func (suite *MessageTypesTestSuite) validateSystemContent(content map[string]interface{}) error {
	requiredFields := []string{"type", "message"}
	for _, field := range requiredFields {
		if _, exists := content[field]; !exists {
			return fmt.Errorf("system content must have '%s' field", field)
		}
	}

	return nil
}

// TEST 1: Text Message Type
func (suite *MessageTypesTestSuite) TestTextMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name:         "Valid text message",
			content:      map[string]interface{}{"text": "Hello, world!"},
			expectStatus: http.StatusCreated,
		},
		{
			name:         "Text with entities",
			content:      map[string]interface{}{"text": "Check this @user", "entities": []map[string]interface{}{{"type": "mention", "offset": 11, "length": 5}}},
			expectStatus: http.StatusCreated,
		},
		{
			name:         "Empty text should fail",
			content:      map[string]interface{}{"text": ""},
			expectStatus: http.StatusBadRequest,
			expectError:  "text content cannot be empty",
		},
		{
			name:         "Missing text field should fail",
			content:      map[string]interface{}{"message": "Hello"},
			expectStatus: http.StatusBadRequest,
			expectError:  "text content must have 'text' field",
		},
		{
			name:         "Long text should fail (>4096 chars)",
			content:      map[string]interface{}{"text": string(make([]byte, 4097))},
			expectStatus: http.StatusBadRequest,
			expectError:  "text content cannot exceed 4096 characters",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "text",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 2: Voice Message Type
func (suite *MessageTypesTestSuite) TestVoiceMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name: "Valid voice message",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/voice/audio.ogg",
				"duration":  5000.0,
				"file_size": 1024000.0,
				"mime_type": "audio/ogg",
				"waveform":  []int{1, 2, 3, 4, 5},
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Voice with long duration should fail",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/voice/long.ogg",
				"duration":  400000.0, // > 5 minutes
				"file_size": 1024000.0,
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "voice duration must be between 1ms and 5 minutes",
		},
		{
			name: "Voice with large file should fail",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/voice/large.ogg",
				"duration":  5000.0,
				"file_size": 60 * 1024 * 1024.0, // > 50MB
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "voice file size must be between 1 byte and 50MB",
		},
		{
			name: "Missing required field should fail",
			content: map[string]interface{}{
				"url":      "https://cdn.tchat.sea/voice/incomplete.ogg",
				"duration": 5000.0,
				// missing file_size
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "voice content must have 'file_size' field",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "voice",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 3: File Message Type
func (suite *MessageTypesTestSuite) TestFileMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name: "Valid file message",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/files/document.pdf",
				"filename":  "document.pdf",
				"file_size": 2048000.0,
				"mime_type": "application/pdf",
				"caption":   "Important document",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "File too large should fail",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/files/large.zip",
				"filename":  "large.zip",
				"file_size": 150 * 1024 * 1024.0, // > 100MB
				"mime_type": "application/zip",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "file size must be between 1 byte and 100MB",
		},
		{
			name: "Missing filename should fail",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/files/unnamed",
				"file_size": 1024.0,
				"mime_type": "text/plain",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "file content must have 'filename' field",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "file",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 4: Image Message Type
func (suite *MessageTypesTestSuite) TestImageMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name: "Valid image message",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/images/photo.jpg",
				"thumbnail": "https://cdn.tchat.sea/thumbs/photo_thumb.jpg",
				"width":     1920.0,
				"height":    1080.0,
				"file_size": 5242880.0, // 5MB
				"caption":   "Beautiful sunset",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Image too wide should fail",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/images/wide.jpg",
				"width":     10000.0, // > 8192
				"height":    1000.0,
				"file_size": 1024000.0,
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "image width must be between 1 and 8192 pixels",
		},
		{
			name: "Image too tall should fail",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/images/tall.jpg",
				"width":     1000.0,
				"height":    10000.0, // > 8192
				"file_size": 1024000.0,
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "image height must be between 1 and 8192 pixels",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "image",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 5: Video Message Type
func (suite *MessageTypesTestSuite) TestVideoMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name: "Valid video message",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/videos/clip.mp4",
				"thumbnail": "https://cdn.tchat.sea/thumbs/clip_thumb.jpg",
				"duration":  30000.0, // 30 seconds
				"width":     1280.0,
				"height":    720.0,
				"file_size": 52428800.0, // 50MB
				"caption":   "Funny video",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Video too long should fail",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/videos/long.mp4",
				"duration":  2000000.0, // > 30 minutes
				"width":     1280.0,
				"height":    720.0,
				"file_size": 100000000.0,
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "video duration must be between 1ms and 30 minutes",
		},
		{
			name: "Video too large should fail",
			content: map[string]interface{}{
				"url":       "https://cdn.tchat.sea/videos/huge.mp4",
				"duration":  30000.0,
				"width":     1920.0,
				"height":    1080.0,
				"file_size": 600 * 1024 * 1024.0, // > 500MB
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "video file size must be between 1 byte and 500MB",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "video",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 6: Payment Message Type
func (suite *MessageTypesTestSuite) TestPaymentMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name: "Valid payment message - THB",
			content: map[string]interface{}{
				"amount":      10000.0, // 100 THB
				"currency":    "THB",
				"description": "Coffee payment",
				"status":      "pending",
				"reference":   "PAY123456",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Valid payment message - USD",
			content: map[string]interface{}{
				"amount":      2500.0, // $25.00
				"currency":    "USD",
				"description": "Lunch payment",
				"status":      "completed",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Negative amount should fail",
			content: map[string]interface{}{
				"amount":      -1000.0,
				"currency":    "THB",
				"description": "Invalid payment",
				"status":      "pending",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "payment amount must be positive",
		},
		{
			name: "Invalid currency should fail",
			content: map[string]interface{}{
				"amount":      1000.0,
				"currency":    "INVALID",
				"description": "Test payment",
				"status":      "pending",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "invalid currency: INVALID",
		},
		{
			name: "Invalid status should fail",
			content: map[string]interface{}{
				"amount":      1000.0,
				"currency":    "THB",
				"description": "Test payment",
				"status":      "invalid_status",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "invalid payment status: invalid_status",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "payment",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 7: Location Message Type
func (suite *MessageTypesTestSuite) TestLocationMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name: "Valid location message - Bangkok",
			content: map[string]interface{}{
				"latitude":  13.7563,
				"longitude": 100.5018,
				"address":   "Bangkok, Thailand",
				"venue":     "Siam Paragon",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Valid location message - Singapore",
			content: map[string]interface{}{
				"latitude":  1.3521,
				"longitude": 103.8198,
				"address":   "Singapore",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Invalid latitude should fail",
			content: map[string]interface{}{
				"latitude":  95.0, // > 90
				"longitude": 100.0,
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "latitude must be between -90 and 90",
		},
		{
			name: "Invalid longitude should fail",
			content: map[string]interface{}{
				"latitude":  45.0,
				"longitude": 200.0, // > 180
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "longitude must be between -180 and 180",
		},
		{
			name: "Missing coordinates should fail",
			content: map[string]interface{}{
				"address": "Bangkok, Thailand",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "location content must have 'latitude' field",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "location",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 8: Sticker Message Type
func (suite *MessageTypesTestSuite) TestStickerMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name: "Valid sticker message",
			content: map[string]interface{}{
				"sticker_id": "sticker_001",
				"pack_id":    "pack_emoji",
				"url":        "https://cdn.tchat.sea/stickers/emoji_001.webp",
				"width":      128.0,
				"height":     128.0,
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Valid animated sticker",
			content: map[string]interface{}{
				"sticker_id": "animated_001",
				"pack_id":    "pack_animated",
				"url":        "https://cdn.tchat.sea/stickers/animated_001.tgs",
				"width":      256.0,
				"height":     256.0,
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Missing sticker_id should fail",
			content: map[string]interface{}{
				"pack_id": "pack_emoji",
				"url":     "https://cdn.tchat.sea/stickers/missing_id.webp",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "sticker content must have 'sticker_id' field",
		},
		{
			name: "Missing pack_id should fail",
			content: map[string]interface{}{
				"sticker_id": "sticker_001",
				"url":        "https://cdn.tchat.sea/stickers/no_pack.webp",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "sticker content must have 'pack_id' field",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "sticker",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 9: System Message Type
func (suite *MessageTypesTestSuite) TestSystemMessageType() {
	testCases := []struct {
		name           string
		content        map[string]interface{}
		expectStatus   int
		expectError    string
	}{
		{
			name: "Valid system message - user joined",
			content: map[string]interface{}{
				"type":    "user_joined",
				"message": "Alice joined the group",
				"data": map[string]interface{}{
					"user_id":   "user_123",
					"user_name": "Alice",
				},
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Valid system message - user left",
			content: map[string]interface{}{
				"type":    "user_left",
				"message": "Bob left the group",
				"data": map[string]interface{}{
					"user_id":   "user_456",
					"user_name": "Bob",
				},
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Valid system message - name changed",
			content: map[string]interface{}{
				"type":    "name_changed",
				"message": "Group name changed to 'New Name'",
				"data": map[string]interface{}{
					"old_name": "Old Name",
					"new_name": "New Name",
				},
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Missing type should fail",
			content: map[string]interface{}{
				"message": "System notification",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "system content must have 'type' field",
		},
		{
			name: "Missing message should fail",
			content: map[string]interface{}{
				"type": "user_joined",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  "system content must have 'message' field",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    "system",
				"content": tc.content,
			}

			suite.sendMessageAndAssert(payload, tc.expectStatus, tc.expectError)
		})
	}
}

// TEST 10: Invalid Message Types
func (suite *MessageTypesTestSuite) TestInvalidMessageTypes() {
	invalidTypes := []string{
		"unknown",
		"invalid",
		"custom",
		"",
		"TEXT", // Case sensitive
		"123",
	}

	for _, invalidType := range invalidTypes {
		suite.Run("Invalid type: "+invalidType, func() {
			payload := map[string]interface{}{
				"type": invalidType,
				"content": map[string]interface{}{
					"text": "Test content",
				},
			}

			jsonData, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", "/dialogs/"+suite.testDialogID+"/messages", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Contains(suite.T(), response["error"], "unsupported message type")
		})
	}
}

// TEST 11: Content Validation Endpoint
func (suite *MessageTypesTestSuite) TestMessageValidationEndpoint() {
	testCases := []struct {
		name         string
		messageType  string
		content      interface{}
		expectValid  bool
	}{
		{
			name:        "Valid text validation",
			messageType: "text",
			content:     map[string]interface{}{"text": "Hello"},
			expectValid: true,
		},
		{
			name:        "Invalid text validation",
			messageType: "text",
			content:     map[string]interface{}{"text": ""},
			expectValid: false,
		},
		{
			name:        "Valid payment validation",
			messageType: "payment",
			content: map[string]interface{}{
				"amount":      1000.0,
				"currency":    "THB",
				"description": "Test",
				"status":      "pending",
			},
			expectValid: true,
		},
		{
			name:        "Invalid payment validation",
			messageType: "payment",
			content: map[string]interface{}{
				"amount":   -1000.0, // Invalid negative amount
				"currency": "THB",
			},
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			payload := map[string]interface{}{
				"type":    tc.messageType,
				"content": tc.content,
			}

			jsonData, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", "/messages/validate", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tc.expectValid {
				assert.Equal(suite.T(), http.StatusOK, w.Code)
				assert.Equal(suite.T(), true, response["valid"])
				assert.Equal(suite.T(), tc.messageType, response["type"])
			} else {
				assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
				assert.Equal(suite.T(), false, response["valid"])
				assert.NotEmpty(suite.T(), response["error"])
			}
		})
	}
}

// TEST 12: Message Storage and Retrieval
func (suite *MessageTypesTestSuite) TestMessageStorageAndRetrieval() {
	// Send a message of each type and verify it can be retrieved
	messageTypes := []map[string]interface{}{
		{"type": "text", "content": map[string]interface{}{"text": "Test text message"}},
		{"type": "voice", "content": map[string]interface{}{"url": "https://example.com/voice.ogg", "duration": 5000.0, "file_size": 1024.0}},
		{"type": "file", "content": map[string]interface{}{"url": "https://example.com/file.pdf", "filename": "test.pdf", "file_size": 1024.0, "mime_type": "application/pdf"}},
		{"type": "image", "content": map[string]interface{}{"url": "https://example.com/image.jpg", "width": 800.0, "height": 600.0, "file_size": 1024.0}},
		{"type": "video", "content": map[string]interface{}{"url": "https://example.com/video.mp4", "duration": 30000.0, "width": 1280.0, "height": 720.0, "file_size": 5120.0}},
		{"type": "payment", "content": map[string]interface{}{"amount": 1000.0, "currency": "THB", "description": "Test payment", "status": "pending"}},
		{"type": "location", "content": map[string]interface{}{"latitude": 13.7563, "longitude": 100.5018, "address": "Bangkok"}},
		{"type": "sticker", "content": map[string]interface{}{"sticker_id": "001", "pack_id": "emoji", "url": "https://example.com/sticker.webp"}},
		{"type": "system", "content": map[string]interface{}{"type": "user_joined", "message": "User joined"}},
	}

	messageIds := make([]string, 0, len(messageTypes))

	// Send all message types
	for _, msgType := range messageTypes {
		suite.Run(fmt.Sprintf("Send and store %s message", msgType["type"]), func() {
			jsonData, _ := json.Marshal(msgType)
			req := httptest.NewRequest("POST", "/dialogs/"+suite.testDialogID+"/messages", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusCreated, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			messageId := response["message_id"].(string)
			messageIds = append(messageIds, messageId)

			assert.NotEmpty(suite.T(), messageId)
			assert.Equal(suite.T(), msgType["type"], response["type"])
			assert.Equal(suite.T(), "sent", response["status"])
		})
	}

	// Retrieve all messages and verify content
	for i, messageId := range messageIds {
		expectedType := messageTypes[i]["type"].(string)
		suite.Run(fmt.Sprintf("Retrieve %s message", expectedType), func() {
			req := httptest.NewRequest("GET", "/dialogs/"+suite.testDialogID+"/messages/"+messageId, nil)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusOK, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(suite.T(), messageId, response["id"])
			assert.Equal(suite.T(), suite.testDialogID, response["dialog_id"])
			assert.Equal(suite.T(), expectedType, response["type"])
			assert.NotEmpty(suite.T(), response["created_at"])
			assert.NotNil(suite.T(), response["content"])
		})
	}
}

// Helper function to send message and assert response
func (suite *MessageTypesTestSuite) sendMessageAndAssert(payload map[string]interface{}, expectStatus int, expectError string) {
	jsonData, err := json.Marshal(payload)
	assert.NoError(suite.T(), err)

	req := httptest.NewRequest("POST", "/dialogs/"+suite.testDialogID+"/messages", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), expectStatus, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if expectStatus == http.StatusCreated {
		assert.NotEmpty(suite.T(), response["message_id"])
		assert.Equal(suite.T(), "sent", response["status"])
		assert.Equal(suite.T(), payload["type"], response["type"])
	} else {
		assert.Contains(suite.T(), response["error"].(string), expectError)
	}
}