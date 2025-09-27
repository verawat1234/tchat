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

// NotificationIntegrationSuite tests the Notification service endpoints
type NotificationIntegrationSuite struct {
	APIIntegrationSuite
	ports ServicePort
}

// Notification represents a notification
type Notification struct {
	ID        string            `json:"id"`
	UserID    string            `json:"userId"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Data      map[string]string `json:"data"`
	Status    string            `json:"status"`
	Priority  string            `json:"priority"`
	Channel   string            `json:"channel"`
	ScheduledAt *string         `json:"scheduledAt,omitempty"`
	SentAt    *string           `json:"sentAt,omitempty"`
	ReadAt    *string           `json:"readAt,omitempty"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
}

// NotificationTemplate represents a notification template
type NotificationTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Variables   []string          `json:"variables"`
	Channels    []string          `json:"channels"`
	Settings    map[string]string `json:"settings"`
	IsActive    bool              `json:"isActive"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
}

// NotificationPreference represents user notification preferences
type NotificationPreference struct {
	ID       string            `json:"id"`
	UserID   string            `json:"userId"`
	Type     string            `json:"type"`
	Channel  string            `json:"channel"`
	Enabled  bool              `json:"enabled"`
	Settings map[string]string `json:"settings"`
	CreatedAt string           `json:"createdAt"`
	UpdatedAt string           `json:"updatedAt"`
}

// CreateNotificationRequest represents notification creation request
type CreateNotificationRequest struct {
	UserID    string            `json:"userId"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Message   string            `json:"message"`
	Data      map[string]string `json:"data"`
	Priority  string            `json:"priority"`
	Channel   string            `json:"channel"`
	ScheduledAt *string         `json:"scheduledAt,omitempty"`
}

// CreateTemplateRequest represents template creation request
type CreateTemplateRequest struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Variables []string          `json:"variables"`
	Channels  []string          `json:"channels"`
	Settings  map[string]string `json:"settings"`
}

// UpdatePreferenceRequest represents preference update request
type UpdatePreferenceRequest struct {
	Enabled  *bool             `json:"enabled,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
}

// NotificationResponse represents notification API response
type NotificationResponse struct {
	Success      bool                       `json:"success"`
	Status       string                     `json:"status"`
	Message      string                     `json:"message"`
	Notification *Notification              `json:"notification,omitempty"`
	Notifications []Notification            `json:"notifications,omitempty"`
	Template     *NotificationTemplate      `json:"template,omitempty"`
	Templates    []NotificationTemplate     `json:"templates,omitempty"`
	Preference   *NotificationPreference    `json:"preference,omitempty"`
	Preferences  []NotificationPreference   `json:"preferences,omitempty"`
	Total        int                        `json:"total,omitempty"`
	Page         int                        `json:"page,omitempty"`
	Limit        int                        `json:"limit,omitempty"`
	Timestamp    string                     `json:"timestamp"`
}

// SetupSuite initializes the notification integration test suite
func (suite *NotificationIntegrationSuite) SetupSuite() {
	suite.APIIntegrationSuite.SetupSuite()
	suite.ports = DefaultServicePorts()

	// Wait for notification service to be available
	err := suite.waitForService(suite.ports.Notification, 30*time.Second)
	if err != nil {
		suite.T().Fatalf("Notification service not available: %v", err)
	}
}

// TestNotificationServiceHealth verifies notification service health endpoint
func (suite *NotificationIntegrationSuite) TestNotificationServiceHealth() {
	healthCheck, err := suite.checkServiceHealth(suite.ports.Notification)
	require.NoError(suite.T(), err, "Health check should succeed")

	assert.Equal(suite.T(), "healthy", healthCheck.Status)
	assert.Equal(suite.T(), "notification-service", healthCheck.Service)
	assert.NotEmpty(suite.T(), healthCheck.Timestamp)
}

// TestCreateNotification tests notification creation
func (suite *NotificationIntegrationSuite) TestCreateNotification() {
	url := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)

	createReq := CreateNotificationRequest{
		UserID:   "test-user-123",
		Type:     "welcome",
		Title:    "Welcome to Tchat!",
		Message:  "Thank you for joining our platform",
		Priority: "medium",
		Channel:  "push",
		Data: map[string]string{
			"action":      "onboarding",
			"screen":      "welcome",
			"user_level":  "new",
		},
	}

	resp, err := suite.makeRequest("POST", url, createReq, nil)
	require.NoError(suite.T(), err, "Create notification request should succeed")
	defer resp.Body.Close()

	// Should return 201 for successful creation
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse notification creation response")

	assert.True(suite.T(), notificationResp.Success)
	assert.Equal(suite.T(), "success", notificationResp.Status)
	assert.NotNil(suite.T(), notificationResp.Notification)
	assert.NotEmpty(suite.T(), notificationResp.Notification.ID)
	assert.Equal(suite.T(), createReq.UserID, notificationResp.Notification.UserID)
	assert.Equal(suite.T(), createReq.Type, notificationResp.Notification.Type)
	assert.Equal(suite.T(), createReq.Title, notificationResp.Notification.Title)
	assert.Equal(suite.T(), "pending", notificationResp.Notification.Status) // Default status
}

// TestGetNotification tests retrieving a specific notification
func (suite *NotificationIntegrationSuite) TestGetNotification() {
	// First create a notification
	createURL := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)
	createReq := CreateNotificationRequest{
		UserID:   "test-user-get",
		Type:     "info",
		Title:    "Get Test Notification",
		Message:  "Notification for get test",
		Priority: "low",
		Channel:  "email",
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create notification for get test should succeed")
	defer createResp.Body.Close()

	var createNotificationResp NotificationResponse
	err = suite.parseResponse(createResp, &createNotificationResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createNotificationResp.Notification)

	notificationID := createNotificationResp.Notification.ID

	// Now get the notification
	getURL := fmt.Sprintf("%s:%d/api/v1/notifications/%s", suite.baseURL, suite.ports.Notification, notificationID)
	resp, err := suite.makeRequest("GET", getURL, nil, nil)
	require.NoError(suite.T(), err, "Get notification request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse get notification response")

	assert.True(suite.T(), notificationResp.Success)
	assert.NotNil(suite.T(), notificationResp.Notification)
	assert.Equal(suite.T(), notificationID, notificationResp.Notification.ID)
	assert.Equal(suite.T(), createReq.UserID, notificationResp.Notification.UserID)
}

// TestListNotifications tests listing notifications
func (suite *NotificationIntegrationSuite) TestListNotifications() {
	url := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List notifications request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse list notifications response")

	assert.True(suite.T(), notificationResp.Success)
	assert.NotNil(suite.T(), notificationResp.Notifications)
	assert.GreaterOrEqual(suite.T(), notificationResp.Total, 0)
}

// TestListUserNotifications tests listing notifications for a specific user
func (suite *NotificationIntegrationSuite) TestListUserNotifications() {
	userID := "test-user-list"

	// First create a notification for this user
	createURL := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)
	createReq := CreateNotificationRequest{
		UserID:   userID,
		Type:     "reminder",
		Title:    "User List Test",
		Message:  "Notification for user list test",
		Priority: "medium",
		Channel:  "push",
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create notification for user list test should succeed")
	createResp.Body.Close()

	// Now list notifications for this user
	url := fmt.Sprintf("%s:%d/api/v1/notifications?userId=%s", suite.baseURL, suite.ports.Notification, userID)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List user notifications request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse user notifications response")

	assert.True(suite.T(), notificationResp.Success)

	// All returned notifications should belong to the specified user
	for _, notification := range notificationResp.Notifications {
		assert.Equal(suite.T(), userID, notification.UserID)
	}
}

// TestMarkNotificationRead tests marking a notification as read
func (suite *NotificationIntegrationSuite) TestMarkNotificationRead() {
	// First create a notification
	createURL := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)
	createReq := CreateNotificationRequest{
		UserID:   "test-user-read",
		Type:     "alert",
		Title:    "Read Test Notification",
		Message:  "Notification for read test",
		Priority: "high",
		Channel:  "push",
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create notification for read test should succeed")
	defer createResp.Body.Close()

	var createNotificationResp NotificationResponse
	err = suite.parseResponse(createResp, &createNotificationResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createNotificationResp.Notification)

	notificationID := createNotificationResp.Notification.ID

	// Mark as read
	readURL := fmt.Sprintf("%s:%d/api/v1/notifications/%s/read", suite.baseURL, suite.ports.Notification, notificationID)
	resp, err := suite.makeRequest("PUT", readURL, nil, nil)
	require.NoError(suite.T(), err, "Mark notification read should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse mark read response")

	assert.True(suite.T(), notificationResp.Success)
	assert.NotNil(suite.T(), notificationResp.Notification)
	assert.NotNil(suite.T(), notificationResp.Notification.ReadAt)
}

// TestDeleteNotification tests notification deletion
func (suite *NotificationIntegrationSuite) TestDeleteNotification() {
	// First create a notification
	createURL := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)
	createReq := CreateNotificationRequest{
		UserID:   "test-user-delete",
		Type:     "temp",
		Title:    "Delete Test Notification",
		Message:  "Notification for delete test",
		Priority: "low",
		Channel:  "email",
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create notification for delete test should succeed")
	defer createResp.Body.Close()

	var createNotificationResp NotificationResponse
	err = suite.parseResponse(createResp, &createNotificationResp)
	require.NoError(suite.T(), err, "Should parse create response")
	require.NotNil(suite.T(), createNotificationResp.Notification)

	notificationID := createNotificationResp.Notification.ID

	// Delete the notification
	deleteURL := fmt.Sprintf("%s:%d/api/v1/notifications/%s", suite.baseURL, suite.ports.Notification, notificationID)
	resp, err := suite.makeRequest("DELETE", deleteURL, nil, nil)
	require.NoError(suite.T(), err, "Delete notification should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify notification is deleted
	getURL := fmt.Sprintf("%s:%d/api/v1/notifications/%s", suite.baseURL, suite.ports.Notification, notificationID)
	getResp, err := suite.makeRequest("GET", getURL, nil, nil)
	require.NoError(suite.T(), err, "Get deleted notification should complete")
	defer getResp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, getResp.StatusCode)
}

// TestCreateNotificationTemplate tests template creation
func (suite *NotificationIntegrationSuite) TestCreateNotificationTemplate() {
	url := fmt.Sprintf("%s:%d/api/v1/notification-templates", suite.baseURL, suite.ports.Notification)

	createReq := CreateTemplateRequest{
		Name:     "Welcome Template",
		Type:     "welcome",
		Title:    "Welcome {{user_name}}!",
		Body:     "Hello {{user_name}}, welcome to {{app_name}}! Your journey starts now.",
		Variables: []string{"user_name", "app_name"},
		Channels: []string{"push", "email"},
		Settings: map[string]string{
			"priority":      "medium",
			"retry_count":   "3",
			"retry_delay":   "300",
		},
	}

	resp, err := suite.makeRequest("POST", url, createReq, nil)
	require.NoError(suite.T(), err, "Create template request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse template creation response")

	assert.True(suite.T(), notificationResp.Success)
	assert.NotNil(suite.T(), notificationResp.Template)
	assert.NotEmpty(suite.T(), notificationResp.Template.ID)
	assert.Equal(suite.T(), createReq.Name, notificationResp.Template.Name)
	assert.Equal(suite.T(), createReq.Type, notificationResp.Template.Type)
	assert.True(suite.T(), notificationResp.Template.IsActive) // Default active
}

// TestListNotificationTemplates tests listing templates
func (suite *NotificationIntegrationSuite) TestListNotificationTemplates() {
	url := fmt.Sprintf("%s:%d/api/v1/notification-templates", suite.baseURL, suite.ports.Notification)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "List templates request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse list templates response")

	assert.True(suite.T(), notificationResp.Success)
	assert.NotNil(suite.T(), notificationResp.Templates)
	assert.GreaterOrEqual(suite.T(), notificationResp.Total, 0)
}

// TestGetNotificationPreferences tests getting user preferences
func (suite *NotificationIntegrationSuite) TestGetNotificationPreferences() {
	userID := "test-user-preferences"
	url := fmt.Sprintf("%s:%d/api/v1/users/%s/notification-preferences", suite.baseURL, suite.ports.Notification, userID)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Get preferences request should succeed")
	defer resp.Body.Close()

	// Should return either 200 with preferences or 404 if not set
	assert.True(suite.T(), resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound)

	if resp.StatusCode == http.StatusOK {
		var notificationResp NotificationResponse
		err = suite.parseResponse(resp, &notificationResp)
		require.NoError(suite.T(), err, "Should parse preferences response")

		assert.True(suite.T(), notificationResp.Success)
		assert.NotNil(suite.T(), notificationResp.Preferences)
	}
}

// TestUpdateNotificationPreference tests updating user preference
func (suite *NotificationIntegrationSuite) TestUpdateNotificationPreference() {
	userID := "test-user-update-pref"
	notificationType := "welcome"
	channel := "push"

	url := fmt.Sprintf("%s:%d/api/v1/users/%s/notification-preferences/%s/%s",
		suite.baseURL, suite.ports.Notification, userID, notificationType, channel)

	updateReq := UpdatePreferenceRequest{
		Enabled: func() *bool { b := false; return &b }(), // Disable
		Settings: map[string]string{
			"frequency": "weekly",
			"time":      "morning",
		},
	}

	resp, err := suite.makeRequest("PUT", url, updateReq, nil)
	require.NoError(suite.T(), err, "Update preference request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse update preference response")

	assert.True(suite.T(), notificationResp.Success)
	assert.NotNil(suite.T(), notificationResp.Preference)
	assert.Equal(suite.T(), userID, notificationResp.Preference.UserID)
	assert.Equal(suite.T(), notificationType, notificationResp.Preference.Type)
	assert.Equal(suite.T(), channel, notificationResp.Preference.Channel)
	assert.False(suite.T(), notificationResp.Preference.Enabled)
}

// TestGetNonExistentNotification tests retrieving non-existent notification
func (suite *NotificationIntegrationSuite) TestGetNonExistentNotification() {
	url := fmt.Sprintf("%s:%d/api/v1/notifications/non-existent-id", suite.baseURL, suite.ports.Notification)

	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Get non-existent notification should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// TestCreateNotificationInvalidData tests notification creation with invalid data
func (suite *NotificationIntegrationSuite) TestCreateNotificationInvalidData() {
	url := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)

	// Test with missing required fields
	invalidReq := CreateNotificationRequest{
		// Missing userID, type, title, message
		Priority: "low",
		Channel:  "push",
	}

	resp, err := suite.makeRequest("POST", url, invalidReq, nil)
	require.NoError(suite.T(), err, "Invalid notification creation should complete")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// TestNotificationsByType tests filtering notifications by type
func (suite *NotificationIntegrationSuite) TestNotificationsByType() {
	// First create a notification with specific type
	createURL := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)
	createReq := CreateNotificationRequest{
		UserID:   "test-user-type-filter",
		Type:     "test-type-filter",
		Title:    "Type Filter Test",
		Message:  "Notification for type filtering test",
		Priority: "medium",
		Channel:  "push",
	}

	createResp, err := suite.makeRequest("POST", createURL, createReq, nil)
	require.NoError(suite.T(), err, "Create notification for type test should succeed")
	createResp.Body.Close()

	// Now filter by type
	url := fmt.Sprintf("%s:%d/api/v1/notifications?type=test-type-filter", suite.baseURL, suite.ports.Notification)
	resp, err := suite.makeRequest("GET", url, nil, nil)
	require.NoError(suite.T(), err, "Type filter request should succeed")
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var notificationResp NotificationResponse
	err = suite.parseResponse(resp, &notificationResp)
	require.NoError(suite.T(), err, "Should parse type filter response")

	assert.True(suite.T(), notificationResp.Success)

	// All returned notifications should have the specified type
	for _, notification := range notificationResp.Notifications {
		assert.Equal(suite.T(), "test-type-filter", notification.Type)
	}
}

// TestInvalidHTTPMethods tests endpoints with invalid HTTP methods
func (suite *NotificationIntegrationSuite) TestInvalidHTTPMethods() {
	baseURL := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)

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

// TestNotificationServiceConcurrency tests concurrent requests to notification service
func (suite *NotificationIntegrationSuite) TestNotificationServiceConcurrency() {
	url := fmt.Sprintf("%s:%d/api/v1/notifications", suite.baseURL, suite.ports.Notification)

	// Create 5 concurrent notification creation requests
	concurrency := 5
	results := make(chan int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			createReq := CreateNotificationRequest{
				UserID:   fmt.Sprintf("concurrent-user-%d", index),
				Type:     "concurrency-test",
				Title:    fmt.Sprintf("Concurrent Notification %d", index),
				Message:  fmt.Sprintf("Notification created concurrently #%d", index),
				Priority: "medium",
				Channel:  "push",
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

// RunNotificationIntegrationTests runs the notification integration test suite
func RunNotificationIntegrationTests(t *testing.T) {
	suite.Run(t, new(NotificationIntegrationSuite))
}