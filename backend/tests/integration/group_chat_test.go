package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// T042: Comprehensive Group Chat Testing Suite
// Tests all aspects of group chat functionality including creation, participant management,
// permissions, messaging, and edge cases for scalable multi-user conversations.

type GroupChatTestSuite struct {
	suite.Suite
	router      *gin.Engine
	users       map[string]string // userID -> name
	groups      map[string]GroupInfo
	mu          sync.RWMutex
}

type GroupInfo struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	CreatorID    string   `json:"creator_id"`
	Participants []string `json:"participants"`
	AdminIDs     []string `json:"admin_ids"`
	Settings     GroupSettings `json:"settings"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type GroupSettings struct {
	IsPublic        bool   `json:"is_public"`
	JoinByLink      bool   `json:"join_by_link"`
	MaxParticipants int    `json:"max_participants"`
	MessageHistory  string `json:"message_history"`
	WhoCanInvite    string `json:"who_can_invite"`
	WhoCanMessage   string `json:"who_can_message"`
}

type GroupMessage struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

func TestGroupChatSuite(t *testing.T) {
	suite.Run(t, new(GroupChatTestSuite))
}

func (suite *GroupChatTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.groups = make(map[string]GroupInfo)

	// Initialize test users
	suite.users = map[string]string{
		"user_admin":     "Alice Admin",
		"user_moderator": "Bob Moderator",
		"user_member1":   "Charlie Member",
		"user_member2":   "Diana Member",
		"user_member3":   "Eve Member",
		"user_external":  "Frank External",
	}

	suite.setupGroupChatEndpoints()
}

func (suite *GroupChatTestSuite) SetupTest() {
	// Clear groups for each test to ensure isolation
	suite.mu.Lock()
	suite.groups = make(map[string]GroupInfo)
	suite.mu.Unlock()
}

func (suite *GroupChatTestSuite) setupGroupChatEndpoints() {
	// Authentication middleware
	suite.router.Use(func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userID := auth[7:] // Remove "Bearer " prefix
		if _, exists := suite.users[userID]; !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_user"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	})

	// POST /dialogs - Create group
	suite.router.POST("/dialogs", func(c *gin.Context) {
		userID := c.GetString("user_id")

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		// Validate required fields
		dialogType, ok := req["type"].(string)
		if !ok || dialogType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "type_required"})
			return
		}

		// Validate dialog type
		validTypes := []string{"user", "group", "channel", "business"}
		isValidType := false
		for _, validType := range validTypes {
			if dialogType == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_dialog_type"})
			return
		}

		// For groups, channels, and business chats, require name
		var name string
		if dialogType == "group" || dialogType == "channel" || dialogType == "business" {
			if n, exists := req["title"].(string); exists && n != "" {
				name = n
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "group_name_required"})
				return
			}
		}

		// Get participants
		participants := []string{userID} // Creator is always a participant
		if participantIDs, exists := req["participant_ids"].([]interface{}); exists {
			for _, pid := range participantIDs {
				if participantID, ok := pid.(string); ok && participantID != userID {
					if _, userExists := suite.users[participantID]; userExists {
						participants = append(participants, participantID)
					}
				}
			}
		}

		// Validate participant limits
		maxParticipants := suite.getMaxParticipants(dialogType)
		if len(participants) > maxParticipants {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "participant_limit_exceeded",
				"max":   maxParticipants,
			})
			return
		}

		// Create group
		groupID := fmt.Sprintf("group_%d", time.Now().UnixNano())
		description := ""
		if desc, exists := req["description"].(string); exists {
			description = desc
		}

		group := GroupInfo{
			ID:           groupID,
			Type:         dialogType,
			Name:         name,
			Description:  description,
			CreatorID:    userID,
			Participants: participants,
			AdminIDs:     []string{userID}, // Creator is admin by default
			Settings:     suite.getDefaultSettings(dialogType),
			CreatedAt:    time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
		}

		suite.mu.Lock()
		suite.groups[groupID] = group
		suite.mu.Unlock()

		c.JSON(http.StatusCreated, gin.H{
			"id":                group.ID,
			"type":              group.Type,
			"name":              group.Name,
			"description":       group.Description,
			"creator_id":        group.CreatorID,
			"participant_count": len(group.Participants),
			"settings":          group.Settings,
			"created_at":        group.CreatedAt,
		})
	})

	// GET /dialogs - List user's groups
	suite.router.GET("/dialogs", func(c *gin.Context) {
		userID := c.GetString("user_id")
		dialogType := c.Query("type")

		var userGroups []GroupInfo
		suite.mu.RLock()
		for _, group := range suite.groups {
			// Check if user is participant
			isParticipant := false
			for _, pid := range group.Participants {
				if pid == userID {
					isParticipant = true
					break
				}
			}

			if isParticipant && (dialogType == "" || group.Type == dialogType) {
				userGroups = append(userGroups, group)
			}
		}
		suite.mu.RUnlock()

		c.JSON(http.StatusOK, gin.H{
			"dialogs": userGroups,
			"total":   len(userGroups),
		})
	})

	// GET /dialogs/:id - Get group details
	suite.router.GET("/dialogs/:id", func(c *gin.Context) {
		groupID := c.Param("id")
		userID := c.GetString("user_id")

		suite.mu.RLock()
		group, exists := suite.groups[groupID]
		suite.mu.RUnlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "group_not_found"})
			return
		}

		// Check if user is participant
		isParticipant := false
		for _, pid := range group.Participants {
			if pid == userID {
				isParticipant = true
				break
			}
		}

		if !isParticipant {
			c.JSON(http.StatusForbidden, gin.H{"error": "access_denied"})
			return
		}

		c.JSON(http.StatusOK, group)
	})

	// POST /dialogs/:id/participants - Add participants
	suite.router.POST("/dialogs/:id/participants", func(c *gin.Context) {
		groupID := c.Param("id")
		userID := c.GetString("user_id")

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		suite.mu.Lock()
		group, exists := suite.groups[groupID]
		if !exists {
			suite.mu.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "group_not_found"})
			return
		}

		// Check if user can invite
		if !suite.canUserInvite(group, userID) {
			suite.mu.Unlock()
			c.JSON(http.StatusForbidden, gin.H{"error": "permission_denied"})
			return
		}

		// Get new participants
		participantIDs, ok := req["participant_ids"].([]interface{})
		if !ok {
			suite.mu.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "participant_ids_required"})
			return
		}

		var newParticipants []string
		for _, pid := range participantIDs {
			if participantID, ok := pid.(string); ok {
				// Check if user exists
				if _, userExists := suite.users[participantID]; !userExists {
					suite.mu.Unlock()
					c.JSON(http.StatusBadRequest, gin.H{"error": "user_not_found", "user_id": participantID})
					return
				}

				// Check if already participant
				isAlready := false
				for _, existing := range group.Participants {
					if existing == participantID {
						isAlready = true
						break
					}
				}

				if !isAlready {
					newParticipants = append(newParticipants, participantID)
				}
			}
		}

		// Check participant limit
		maxParticipants := suite.getMaxParticipants(group.Type)
		if len(group.Participants)+len(newParticipants) > maxParticipants {
			suite.mu.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "participant_limit_exceeded",
				"max":   maxParticipants,
			})
			return
		}

		// Add participants
		group.Participants = append(group.Participants, newParticipants...)
		group.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		suite.groups[groupID] = group
		suite.mu.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"message":           "participants_added",
			"added_count":       len(newParticipants),
			"total_participants": len(group.Participants),
		})
	})

	// POST /dialogs/:id/participants/remove - Remove participant
	suite.router.POST("/dialogs/:id/participants/remove", func(c *gin.Context) {
		groupID := c.Param("id")
		userID := c.GetString("user_id")

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		participantID, ok := req["participant_id"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "participant_id_required"})
			return
		}

		suite.mu.Lock()
		group, exists := suite.groups[groupID]
		if !exists {
			suite.mu.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "group_not_found"})
			return
		}

		// Check permissions - only admins can remove others, users can remove themselves
		canRemove := false
		if participantID == userID {
			canRemove = true // Can always remove yourself (leave group)
		} else {
			// Check if user is admin
			for _, adminID := range group.AdminIDs {
				if adminID == userID {
					canRemove = true
					break
				}
			}
		}

		if !canRemove {
			suite.mu.Unlock()
			c.JSON(http.StatusForbidden, gin.H{"error": "permission_denied"})
			return
		}

		// Cannot remove group creator unless they're leaving themselves
		if participantID == group.CreatorID && userID != participantID {
			suite.mu.Unlock()
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot_remove_creator"})
			return
		}

		// Find and remove participant
		found := false
		newParticipants := make([]string, 0, len(group.Participants))
		for _, pid := range group.Participants {
			if pid != participantID {
				newParticipants = append(newParticipants, pid)
			} else {
				found = true
			}
		}

		if !found {
			suite.mu.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "participant_not_found"})
			return
		}

		// Also remove from admins if present
		newAdminIDs := make([]string, 0, len(group.AdminIDs))
		for _, adminID := range group.AdminIDs {
			if adminID != participantID {
				newAdminIDs = append(newAdminIDs, adminID)
			}
		}

		group.Participants = newParticipants
		group.AdminIDs = newAdminIDs
		group.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		suite.groups[groupID] = group
		suite.mu.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"message":            "participant_removed",
			"total_participants": len(group.Participants),
		})
	})

	// POST /dialogs/:id/messages - Send message to group
	suite.router.POST("/dialogs/:id/messages", func(c *gin.Context) {
		groupID := c.Param("id")
		userID := c.GetString("user_id")

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		suite.mu.RLock()
		group, exists := suite.groups[groupID]
		suite.mu.RUnlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "group_not_found"})
			return
		}

		// Check if user can message
		if !suite.canUserMessage(group, userID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission_denied"})
			return
		}

		content, ok := req["content"].(string)
		if !ok || content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content_required"})
			return
		}

		messageType := "text"
		if t, exists := req["type"].(string); exists {
			messageType = t
		}

		// Create message
		message := GroupMessage{
			ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
			GroupID:   groupID,
			SenderID:  userID,
			Content:   content,
			Type:      messageType,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		}

		c.JSON(http.StatusCreated, message)
	})

	// PUT /dialogs/:id - Update group settings
	suite.router.PUT("/dialogs/:id", func(c *gin.Context) {
		groupID := c.Param("id")
		userID := c.GetString("user_id")

		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		suite.mu.Lock()
		group, exists := suite.groups[groupID]
		if !exists {
			suite.mu.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "group_not_found"})
			return
		}

		// Check if user is admin
		isAdmin := false
		for _, adminID := range group.AdminIDs {
			if adminID == userID {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			suite.mu.Unlock()
			c.JSON(http.StatusForbidden, gin.H{"error": "permission_denied"})
			return
		}

		// Update fields
		if title, exists := req["title"].(string); exists {
			group.Name = title
		}
		if description, exists := req["description"].(string); exists {
			group.Description = description
		}

		group.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		suite.groups[groupID] = group
		suite.mu.Unlock()

		c.JSON(http.StatusOK, group)
	})
}

// Helper methods
func (suite *GroupChatTestSuite) getMaxParticipants(dialogType string) int {
	switch dialogType {
	case "user":
		return 2
	case "group":
		return 5000
	case "channel":
		return 200000
	case "business":
		return 100
	default:
		return 2
	}
}

func (suite *GroupChatTestSuite) getDefaultSettings(dialogType string) GroupSettings {
	switch dialogType {
	case "group":
		return GroupSettings{
			IsPublic:        false,
			JoinByLink:      false,
			MaxParticipants: 5000,
			MessageHistory:  "members",
			WhoCanInvite:    "admins",
			WhoCanMessage:   "members",
		}
	case "channel":
		return GroupSettings{
			IsPublic:        true,
			JoinByLink:      true,
			MaxParticipants: 200000,
			MessageHistory:  "everyone",
			WhoCanInvite:    "admins",
			WhoCanMessage:   "admins",
		}
	case "business":
		return GroupSettings{
			IsPublic:        false,
			JoinByLink:      false,
			MaxParticipants: 100,
			MessageHistory:  "members",
			WhoCanInvite:    "admins",
			WhoCanMessage:   "members",
		}
	default:
		return GroupSettings{
			IsPublic:        false,
			JoinByLink:      false,
			MaxParticipants: 2,
			MessageHistory:  "members",
			WhoCanInvite:    "members",
			WhoCanMessage:   "members",
		}
	}
}

func (suite *GroupChatTestSuite) canUserInvite(group GroupInfo, userID string) bool {
	switch group.Settings.WhoCanInvite {
	case "everyone":
		return true
	case "members":
		for _, pid := range group.Participants {
			if pid == userID {
				return true
			}
		}
	case "admins":
		for _, adminID := range group.AdminIDs {
			if adminID == userID {
				return true
			}
		}
	}
	return false
}

func (suite *GroupChatTestSuite) canUserMessage(group GroupInfo, userID string) bool {
	switch group.Settings.WhoCanMessage {
	case "everyone":
		return true
	case "members":
		for _, pid := range group.Participants {
			if pid == userID {
				return true
			}
		}
	case "admins":
		for _, adminID := range group.AdminIDs {
			if adminID == userID {
				return true
			}
		}
	}
	return false
}

// Test Cases

func (suite *GroupChatTestSuite) TestCreateGroup() {
	suite.T().Log("Testing group creation with different dialog types")

	testCases := []struct {
		name           string
		dialogType     string
		title          string
		participants   []string
		expectStatus   int
		expectError    string
	}{
		{
			name:         "Create basic group",
			dialogType:   "group",
			title:        "Test Group",
			participants: []string{"user_member1", "user_member2"},
			expectStatus: http.StatusCreated,
		},
		{
			name:         "Create channel",
			dialogType:   "channel",
			title:        "Test Channel",
			participants: []string{"user_member1"},
			expectStatus: http.StatusCreated,
		},
		{
			name:         "Create business chat",
			dialogType:   "business",
			title:        "Support Chat",
			participants: []string{"user_member1"},
			expectStatus: http.StatusCreated,
		},
		{
			name:         "Group without name should fail",
			dialogType:   "group",
			title:        "",
			participants: []string{"user_member1"},
			expectStatus: http.StatusBadRequest,
			expectError:  "group_name_required",
		},
		{
			name:         "Invalid dialog type",
			dialogType:   "invalid",
			title:        "Test",
			participants: []string{"user_member1"},
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			reqBody := map[string]interface{}{
				"type":           tc.dialogType,
				"title":          tc.title,
				"participant_ids": tc.participants,
			}

			jsonData, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
			req.Header.Set("Authorization", "Bearer user_admin")
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectStatus, w.Code)

			if tc.expectStatus == http.StatusCreated {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tc.dialogType, response["type"])
				assert.Equal(t, tc.title, response["name"])
				assert.Contains(t, response, "id")
			} else if tc.expectError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tc.expectError, response["error"])
			}
		})
	}
}

func (suite *GroupChatTestSuite) TestParticipantManagement() {
	suite.T().Log("Testing group participant management")

	// First create a group
	reqBody := map[string]interface{}{
		"type":           "group",
		"title":          "Test Group",
		"participant_ids": []string{"user_member1"},
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_admin")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	groupID := createResponse["id"].(string)

	// Test adding participants
	suite.T().Run("Add participants as admin", func(t *testing.T) {
		addReq := map[string]interface{}{
			"participant_ids": []string{"user_member2", "user_member3"},
		}

		jsonData, _ := json.Marshal(addReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/participants", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_admin")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(2), response["added_count"])
		assert.Equal(t, float64(4), response["total_participants"]) // admin + member1 + member2 + member3
	})

	// Test adding participants as non-admin (should fail)
	suite.T().Run("Add participants as non-admin should fail", func(t *testing.T) {
		addReq := map[string]interface{}{
			"participant_ids": []string{"user_external"},
		}

		jsonData, _ := json.Marshal(addReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/participants", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_member1") // Not admin
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Test removing participant as admin
	suite.T().Run("Remove participant as admin", func(t *testing.T) {
		removeReq := map[string]interface{}{
			"participant_id": "user_member2",
		}

		jsonData, _ := json.Marshal(removeReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/participants/remove", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_admin")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(3), response["total_participants"])
	})

	// Test self-removal (leaving group)
	suite.T().Run("Member can leave group", func(t *testing.T) {
		removeReq := map[string]interface{}{
			"participant_id": "user_member1",
		}

		jsonData, _ := json.Marshal(removeReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/participants/remove", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_member1") // Self-removal
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test cannot remove creator
	suite.T().Run("Cannot remove group creator", func(t *testing.T) {
		removeReq := map[string]interface{}{
			"participant_id": "user_admin", // Creator
		}

		jsonData, _ := json.Marshal(removeReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/participants/remove", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_member3") // Not creator
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func (suite *GroupChatTestSuite) TestGroupMessaging() {
	suite.T().Log("Testing group messaging functionality")

	// Create a group first
	reqBody := map[string]interface{}{
		"type":           "group",
		"title":          "Chat Group",
		"participant_ids": []string{"user_member1", "user_member2"},
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_admin")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	groupID := createResponse["id"].(string)

	// Test sending message as participant
	suite.T().Run("Send message as participant", func(t *testing.T) {
		msgReq := map[string]interface{}{
			"content": "Hello group!",
			"type":    "text",
		}

		jsonData, _ := json.Marshal(msgReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_member1")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response GroupMessage
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Hello group!", response.Content)
		assert.Equal(t, "user_member1", response.SenderID)
		assert.Equal(t, groupID, response.GroupID)
	})

	// Test sending message as non-participant should fail
	suite.T().Run("Send message as non-participant should fail", func(t *testing.T) {
		msgReq := map[string]interface{}{
			"content": "Hello from outside!",
		}

		jsonData, _ := json.Marshal(msgReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_external")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func (suite *GroupChatTestSuite) TestGroupPermissions() {
	suite.T().Log("Testing group permissions and access control")

	// Create a channel (admins-only messaging)
	reqBody := map[string]interface{}{
		"type":           "channel",
		"title":          "Test Channel",
		"participant_ids": []string{"user_member1", "user_member2"},
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_admin")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	channelID := createResponse["id"].(string)

	// Test member cannot message in channel (admin-only)
	suite.T().Run("Member cannot message in admin-only channel", func(t *testing.T) {
		msgReq := map[string]interface{}{
			"content": "Member message",
		}

		jsonData, _ := json.Marshal(msgReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", channelID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_member1")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Test admin can message in channel
	suite.T().Run("Admin can message in channel", func(t *testing.T) {
		msgReq := map[string]interface{}{
			"content": "Admin announcement",
		}

		jsonData, _ := json.Marshal(msgReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", channelID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_admin")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func (suite *GroupChatTestSuite) TestGroupListing() {
	suite.T().Log("Testing group listing and filtering")

	// Create multiple groups and channels
	groups := []struct {
		dialogType string
		title      string
		user       string
	}{
		{"group", "Group 1", "user_admin"},
		{"group", "Group 2", "user_admin"},
		{"channel", "Channel 1", "user_admin"},
		{"business", "Business Chat", "user_admin"},
	}

	for _, g := range groups {
		reqBody := map[string]interface{}{
			"type":           g.dialogType,
			"title":          g.title,
			"participant_ids": []string{"user_member1"},
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.user))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		assert.Equal(suite.T(), http.StatusCreated, w.Code)
	}

	// Test listing all groups for user
	suite.T().Run("List all dialogs for user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dialogs", nil)
		req.Header.Set("Authorization", "Bearer user_admin")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		dialogs := response["dialogs"].([]interface{})
		assert.Equal(t, 4, len(dialogs))
		assert.Equal(t, float64(4), response["total"])
	})

	// Test filtering by type
	suite.T().Run("Filter groups by type", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dialogs?type=group", nil)
		req.Header.Set("Authorization", "Bearer user_admin")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		dialogs := response["dialogs"].([]interface{})
		assert.Equal(t, 2, len(dialogs)) // Only groups

		// Verify all are groups
		for _, dialog := range dialogs {
			d := dialog.(map[string]interface{})
			assert.Equal(t, "group", d["type"])
		}
	})
}

func (suite *GroupChatTestSuite) TestGroupUpdates() {
	suite.T().Log("Testing group updates and settings changes")

	// Create a group
	reqBody := map[string]interface{}{
		"type":           "group",
		"title":          "Original Name",
		"description":    "Original description",
		"participant_ids": []string{"user_member1"},
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_admin")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	groupID := createResponse["id"].(string)

	// Test updating group as admin
	suite.T().Run("Update group as admin", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"title":       "Updated Name",
			"description": "Updated description",
		}

		jsonData, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/dialogs/%s", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_admin")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Updated Name", response["name"])
		assert.Equal(t, "Updated description", response["description"])
	})

	// Test updating group as non-admin should fail
	suite.T().Run("Update group as non-admin should fail", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"title": "Unauthorized Update",
		}

		jsonData, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/dialogs/%s", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_member1") // Not admin
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func (suite *GroupChatTestSuite) TestParticipantLimits() {
	suite.T().Log("Testing group participant limits")

	// Test group limit (5000)
	suite.T().Run("Group supports large participant count", func(t *testing.T) {
		// Create many participant IDs (simulated)
		participants := make([]string, 4999) // Plus creator = 5000 total
		for i := 0; i < 4999; i++ {
			userID := fmt.Sprintf("user_%d", i)
			// Add to users map for validation
			suite.users[userID] = fmt.Sprintf("User %d", i)
			participants[i] = userID
		}

		reqBody := map[string]interface{}{
			"type":           "group",
			"title":          "Large Group",
			"participant_ids": participants,
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_admin")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(5000), response["participant_count"])
	})

	// Test exceeding group limit
	suite.T().Run("Exceed group participant limit", func(t *testing.T) {
		participants := make([]string, 5000) // Plus creator = 5001 total (exceeds limit)
		for i := 0; i < 5000; i++ {
			userID := fmt.Sprintf("user_large_%d", i)
			suite.users[userID] = fmt.Sprintf("Large User %d", i)
			participants[i] = userID
		}

		reqBody := map[string]interface{}{
			"type":           "group",
			"title":          "Too Large Group",
			"participant_ids": participants,
		}

		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_admin")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "participant_limit_exceeded", response["error"])
		assert.Equal(t, float64(5000), response["max"])
	})
}

func (suite *GroupChatTestSuite) TestConcurrentGroupOperations() {
	suite.T().Log("Testing concurrent group operations")

	// Create a group first
	reqBody := map[string]interface{}{
		"type":           "group",
		"title":          "Concurrent Test Group",
		"participant_ids": []string{"user_member1"},
	}

	jsonData, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/dialogs", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer user_admin")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	groupID := createResponse["id"].(string)

	// Test concurrent message sending
	suite.T().Run("Concurrent message sending", func(t *testing.T) {
		const numMessages = 10
		const numUsers = 3

		var wg sync.WaitGroup
		successCount := int64(0)
		var mu sync.Mutex

		users := []string{"user_admin", "user_member1", "user_member2"}

		// First add user_member2 to the group
		addReq := map[string]interface{}{
			"participant_ids": []string{"user_member2"},
		}
		jsonData, _ := json.Marshal(addReq)
		req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/participants", groupID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer user_admin")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		for userIdx := 0; userIdx < numUsers; userIdx++ {
			wg.Add(1)
			go func(userID string) {
				defer wg.Done()

				for i := 0; i < numMessages; i++ {
					msgReq := map[string]interface{}{
						"content": fmt.Sprintf("Concurrent message %d from %s", i, userID),
						"type":    "text",
					}

					jsonData, _ := json.Marshal(msgReq)
					req := httptest.NewRequest("POST", fmt.Sprintf("/dialogs/%s/messages", groupID), bytes.NewBuffer(jsonData))
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userID))
					req.Header.Set("Content-Type", "application/json")

					w := httptest.NewRecorder()
					suite.router.ServeHTTP(w, req)

					if w.Code == http.StatusCreated {
						mu.Lock()
						successCount++
						mu.Unlock()
					}
				}
			}(users[userIdx])
		}

		wg.Wait()

		// All messages should be sent successfully
		assert.Equal(t, int64(numMessages*numUsers), successCount)
	})
}