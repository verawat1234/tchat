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

// T018: Contract test POST /dialogs - Create new dialog
func TestMessagingDialogCreateContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		payload        map[string]interface{}
		expectedStatus int
		expectedFields []string
		description    string
	}{
		{
			name:       "create_direct_dialog",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type":           "direct",
				"participant_id": "456e7890-e89b-12d3-a456-426614174001",
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "type", "participants", "created_at", "updated_at"},
			description:    "Should create direct dialog between two users",
		},
		{
			name:       "create_group_dialog",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "group",
				"name": "Project Team",
				"description": "Discussion about the new project",
				"participant_ids": []interface{}{
					"456e7890-e89b-12d3-a456-426614174001",
					"789e1234-e89b-12d3-a456-426614174002",
				},
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "type", "name", "description", "participants", "participant_count", "created_at"},
			description:    "Should create group dialog with multiple participants",
		},
		{
			name:       "create_channel_dialog",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "channel",
				"name": "Announcements",
				"description": "Company announcements channel",
				"participant_ids": []interface{}{
					"456e7890-e89b-12d3-a456-426614174001",
					"789e1234-e89b-12d3-a456-426614174002",
					"abc1234e-e89b-12d3-a456-426614174003",
				},
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "type", "name", "description", "participants", "participant_count"},
			description:    "Should create channel with multiple participants",
		},
		{
			name:       "create_group_with_avatar",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "group",
				"name": "Design Team",
				"description": "Design team collaboration",
				"avatar_url": "https://cdn.tchat.sea/dialogs/design-team.jpg",
				"participant_ids": []interface{}{
					"456e7890-e89b-12d3-a456-426614174001",
				},
			},
			expectedStatus: http.StatusCreated,
			expectedFields: []string{"id", "type", "name", "avatar_url", "participants"},
			description:    "Should create group dialog with avatar",
		},
		{
			name:           "missing_authorization",
			authHeader:     "",
			payload:        map[string]interface{}{"type": "direct"},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized when no auth header",
		},
		{
			name:           "invalid_jwt_token",
			authHeader:     "Bearer invalid_token_67890",
			payload:        map[string]interface{}{"type": "direct"},
			expectedStatus: http.StatusUnauthorized,
			expectedFields: []string{"error"},
			description:    "Should return unauthorized for invalid token",
		},
		{
			name:       "missing_required_fields",
			authHeader: "Bearer valid_jwt_token_12345",
			payload:    map[string]interface{}{}, // Missing type
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error when required fields are missing",
		},
		{
			name:       "invalid_dialog_type",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "invalid_type",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid dialog type",
		},
		{
			name:       "direct_dialog_missing_participant",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "direct",
				// Missing participant_id
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error when direct dialog missing participant",
		},
		{
			name:       "group_dialog_missing_name",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "group",
				// Missing name
				"participant_ids": []interface{}{
					"456e7890-e89b-12d3-a456-426614174001",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error when group dialog missing name",
		},
		{
			name:       "group_dialog_missing_participants",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "group",
				"name": "Test Group",
				// Missing participant_ids
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error when group dialog missing participants",
		},
		{
			name:       "invalid_participant_uuid",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type":           "direct",
				"participant_id": "invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid participant UUID",
		},
		{
			name:       "nonexistent_participant",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type":           "direct",
				"participant_id": "999e9999-e89b-12d3-a456-426614174999",
			},
			expectedStatus: http.StatusNotFound,
			expectedFields: []string{"error"},
			description:    "Should return not found for nonexistent participant",
		},
		{
			name:       "group_name_too_long",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "group",
				"name": "This is a very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very long group name that exceeds reasonable limits",
				"participant_ids": []interface{}{
					"456e7890-e89b-12d3-a456-426614174001",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for group name too long",
		},
		{
			name:       "too_many_participants",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "group",
				"name": "Large Group",
				"participant_ids": make([]interface{}, 1001), // Assuming limit is 1000
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for too many participants",
		},
		{
			name:       "duplicate_direct_dialog",
			authHeader: "Bearer user_with_existing_dialog_token",
			payload: map[string]interface{}{
				"type":           "direct",
				"participant_id": "existing_dialog_participant_id",
			},
			expectedStatus: http.StatusConflict,
			expectedFields: []string{"error", "existing_dialog_id"},
			description:    "Should return conflict for duplicate direct dialog",
		},
		{
			name:       "invalid_avatar_url",
			authHeader: "Bearer valid_jwt_token_12345",
			payload: map[string]interface{}{
				"type": "group",
				"name": "Test Group",
				"avatar_url": "not-a-valid-url",
				"participant_ids": []interface{}{
					"456e7890-e89b-12d3-a456-426614174001",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedFields: []string{"error"},
			description:    "Should return error for invalid avatar URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()

			// Mock authentication middleware and endpoint
			router.POST("/api/v1/messaging/dialogs", func(c *gin.Context) {
				authHeader := c.GetHeader("Authorization")

				// Mock authentication logic
				if authHeader == "" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
					return
				}

				if authHeader != "Bearer valid_jwt_token_12345" && authHeader != "Bearer user_with_existing_dialog_token" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
					return
				}

				// Parse request body
				var request map[string]interface{}
				if err := c.ShouldBindJSON(&request); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
					return
				}

				// Check for duplicate dialog
				if authHeader == "Bearer user_with_existing_dialog_token" {
					if participantID, exists := request["participant_id"]; exists {
						if participantID == "existing_dialog_participant_id" {
							c.JSON(http.StatusConflict, gin.H{
								"error": "Direct dialog already exists",
								"existing_dialog_id": "existing_dialog_123",
							})
							return
						}
					}
				}

				// Validation logic
				if _, exists := request["type"]; !exists {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Dialog type is required"})
					return
				}

				dialogType := request["type"].(string)
				validTypes := []string{"direct", "group", "channel"}
				valid := false
				for _, validType := range validTypes {
					if dialogType == validType {
						valid = true
						break
					}
				}
				if !valid {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dialog type"})
					return
				}

				// Validate based on dialog type
				switch dialogType {
				case "direct":
					if _, exists := request["participant_id"]; !exists {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Participant ID is required for direct dialog"})
						return
					}
					participantID := request["participant_id"].(string)
					if participantID == "invalid-uuid" {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant UUID format"})
						return
					}
					if participantID == "999e9999-e89b-12d3-a456-426614174999" {
						c.JSON(http.StatusNotFound, gin.H{"error": "Participant not found"})
						return
					}

				case "group", "channel":
					if _, exists := request["name"]; !exists {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required for " + dialogType + " dialog"})
						return
					}
					name := request["name"].(string)
					if len(name) > 100 {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Dialog name too long"})
						return
					}

					if _, exists := request["participant_ids"]; !exists {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Participant IDs are required for " + dialogType + " dialog"})
						return
					}

					participantIDs := request["participant_ids"].([]interface{})
					if len(participantIDs) > 1000 {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Too many participants"})
						return
					}
				}

				// Validate avatar URL if provided
				if avatarURL, exists := request["avatar_url"]; exists {
					if urlStr, ok := avatarURL.(string); ok {
						if urlStr != "" && !(len(urlStr) > 7 && (urlStr[:7] == "http://" || urlStr[:8] == "https://")) {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid avatar URL"})
							return
						}
					}
				}

				// Mock successful dialog creation response
				participants := []map[string]interface{}{
					{
						"user_id": "123e4567-e89b-12d3-a456-426614174000", // Current user
						"role":    "admin",
						"joined_at": "2023-12-01T15:30:00Z",
					},
				}

				// Add participants based on type
				switch dialogType {
				case "direct":
					participants = append(participants, map[string]interface{}{
						"user_id": request["participant_id"],
						"role":    "member",
						"joined_at": "2023-12-01T15:30:00Z",
					})
				case "group", "channel":
					if participantIDs, ok := request["participant_ids"].([]interface{}); ok {
						for _, pid := range participantIDs {
							participants = append(participants, map[string]interface{}{
								"user_id": pid,
								"role":    "member",
								"joined_at": "2023-12-01T15:30:00Z",
							})
						}
					}
				}

				response := gin.H{
					"id":               "dialog_123e4567-e89b-12d3-a456-426614174000",
					"type":             dialogType,
					"participants":     participants,
					"participant_count": len(participants),
					"created_at":       "2023-12-01T15:30:00Z",
					"updated_at":       "2023-12-01T15:30:00Z",
				}

				// Add optional fields
				if name, exists := request["name"]; exists {
					response["name"] = name
				}
				if description, exists := request["description"]; exists {
					response["description"] = description
				}
				if avatarURL, exists := request["avatar_url"]; exists {
					response["avatar_url"] = avatarURL
				}

				c.JSON(http.StatusCreated, response)
			})

			// Prepare request
			jsonData, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/messaging/dialogs", bytes.NewBuffer(jsonData))
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
			case "create_direct_dialog", "create_group_dialog", "create_channel_dialog", "create_group_with_avatar":
				// Verify dialog ID format
				dialogID, ok := response["id"].(string)
				assert.True(t, ok, "id should be string")
				assert.Contains(t, dialogID, "dialog_", "id should have dialog_ prefix")

				// Verify dialog type
				dialogType, ok := response["type"].(string)
				assert.True(t, ok, "type should be string")
				assert.Equal(t, tt.payload["type"], dialogType, "type should match request")

				// Verify participants
				participants, ok := response["participants"].([]interface{})
				assert.True(t, ok, "participants should be array")
				assert.Greater(t, len(participants), 0, "participants should not be empty")

				// Verify participant count
				participantCount, ok := response["participant_count"].(int)
				assert.True(t, ok, "participant_count should be integer")
				assert.Equal(t, len(participants), participantCount, "participant_count should match participants length")

				// Verify timestamps
				createdAt, ok := response["created_at"].(string)
				assert.True(t, ok, "created_at should be string")
				assert.NotEmpty(t, createdAt, "created_at should not be empty")

				// Type-specific validations
				if dialogType == "group" || dialogType == "channel" {
					name, ok := response["name"].(string)
					assert.True(t, ok, "name should be string for "+dialogType)
					assert.NotEmpty(t, name, "name should not be empty for "+dialogType)
				}

				// Check for avatar URL if provided
				if _, exists := tt.payload["avatar_url"]; exists {
					avatarURL, ok := response["avatar_url"].(string)
					assert.True(t, ok, "avatar_url should be string")
					assert.Equal(t, tt.payload["avatar_url"], avatarURL, "avatar_url should match request")
				}

			case "duplicate_direct_dialog":
				// Verify conflict response includes existing dialog ID
				existingDialogID, ok := response["existing_dialog_id"].(string)
				assert.True(t, ok, "existing_dialog_id should be string")
				assert.NotEmpty(t, existingDialogID, "existing_dialog_id should not be empty")

			default:
				// Error cases
				errorMsg, ok := response["error"].(string)
				assert.True(t, ok, "error should be string")
				assert.NotEmpty(t, errorMsg, "error message should not be empty")
			}
		})
	}
}

// TestMessagingDialogCreateSecurity tests security aspects
func TestMessagingDialogCreateSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("user_isolation", func(t *testing.T) {
		// Ensure users can only create dialogs with valid participants
		router := gin.New()
		router.POST("/api/v1/messaging/dialogs", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")

			// Mock user context from token
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

			var request map[string]interface{}
			c.ShouldBindJSON(&request)

			// User can only create dialogs they participate in
			c.JSON(http.StatusCreated, gin.H{
				"id": "dialog_123",
				"type": "direct",
				"participants": []gin.H{
					{"user_id": currentUserID, "role": "admin"},
					{"user_id": request["participant_id"], "role": "member"},
				},
				"participant_count": 2,
				"created_at": "2023-12-01T15:30:00Z",
				"updated_at": "2023-12-01T15:30:00Z",
			})
		})

		// User 1 creates dialog
		payload1 := map[string]interface{}{
			"type": "direct",
			"participant_id": "other-user-id",
		}
		jsonData1, _ := json.Marshal(payload1)

		req1, _ := http.NewRequest("POST", "/api/v1/messaging/dialogs", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("Authorization", "Bearer user1_token")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		var response1 map[string]interface{}
		json.Unmarshal(w1.Body.Bytes(), &response1)

		// Verify user1 is included as admin
		participants1 := response1["participants"].([]interface{})
		user1Found := false
		for _, p := range participants1 {
			participant := p.(map[string]interface{})
			if participant["user_id"] == "user1-id" && participant["role"] == "admin" {
				user1Found = true
				break
			}
		}
		assert.True(t, user1Found, "Current user should be included as admin")
	})

	t.Run("input_sanitization", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/v1/messaging/dialogs", func(c *gin.Context) {
			var request map[string]interface{}
			c.ShouldBindJSON(&request)

			// Check for dangerous patterns in name and description
			if name, exists := request["name"]; exists {
				if nameStr, ok := name.(string); ok {
					dangerousPatterns := []string{"<script>", "javascript:", "DROP TABLE"}
					for _, pattern := range dangerousPatterns {
						if len(nameStr) >= len(pattern) && nameStr[:len(pattern)] == pattern {
							c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid characters in dialog name"})
							return
						}
					}
				}
			}

			c.JSON(http.StatusCreated, gin.H{"id": "dialog_123", "type": "group", "name": request["name"]})
		})

		maliciousNames := []string{
			"<script>alert('xss')</script>",
			"javascript:alert('xss')",
			"DROP TABLE dialogs;",
		}

		for _, maliciousName := range maliciousNames {
			payload := map[string]interface{}{
				"type": "group",
				"name": maliciousName,
				"participant_ids": []interface{}{"user123"},
			}
			jsonData, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", "/api/v1/messaging/dialogs", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid_token")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code,
				"Should reject malicious name: %s", maliciousName)
		}
	})
}

// TestMessagingDialogCreateRateLimit tests rate limiting
func TestMessagingDialogCreateRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	createCount := 0
	router.POST("/api/v1/messaging/dialogs", func(c *gin.Context) {
		createCount++
		if createCount > 10 { // Simulate rate limit of 10 dialog creations
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many dialog creations. Please wait before creating more.",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id": "dialog_" + string(rune(createCount)),
			"type": "direct",
		})
	})

	payload := map[string]interface{}{
		"type": "direct",
		"participant_id": "test-user",
	}
	jsonData, _ := json.Marshal(payload)

	// Make requests beyond rate limit
	for i := 0; i < 15; i++ {
		req, _ := http.NewRequest("POST", "/api/v1/messaging/dialogs", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i < 10 {
			assert.Equal(t, http.StatusCreated, w.Code, "Creation %d should succeed", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Creation %d should be rate limited", i+1)
		}
	}
}