// Journey 7: Notifications & Alerts API Integration Tests
// Comprehensive testing of notification services, push notifications, alert systems,
// real-time updates, and notification preferences management across platforms

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
)

// AuthenticatedUser represents an authenticated user session
type AuthenticatedUser struct {
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Country      string `json:"country"`
	Language     string `json:"language"`
}

// NotificationInfo represents notification data structure
type NotificationInfo struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Type            string    `json:"type"`
	Title           string    `json:"title"`
	Message         string    `json:"message"`
	Category        string    `json:"category"`
	Priority        string    `json:"priority"`
	Status          string    `json:"status"`
	ReadAt          time.Time `json:"read_at,omitempty"`
	ActionURL       string    `json:"action_url,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at,omitempty"`
	PlatformTargets []string  `json:"platform_targets"`
	DeliveryStatus  map[string]string `json:"delivery_status,omitempty"`
}

// NotificationPreferences represents user notification settings
type NotificationPreferences struct {
	UserID           string            `json:"user_id"`
	EmailEnabled     bool              `json:"email_enabled"`
	PushEnabled      bool              `json:"push_enabled"`
	SMSEnabled       bool              `json:"sms_enabled"`
	InAppEnabled     bool              `json:"in_app_enabled"`
	Categories       map[string]bool   `json:"categories"`
	Schedules        map[string]string `json:"schedules"`
	Platforms        map[string]bool   `json:"platforms"`
	Languages        []string          `json:"languages"`
	TimeZone         string            `json:"timezone"`
	QuietHours       map[string]string `json:"quiet_hours"`
	FrequencyLimits  map[string]int    `json:"frequency_limits"`
}

// PushNotificationInfo represents push notification data
type PushNotificationInfo struct {
	ID           string                 `json:"id"`
	DeviceTokens []string               `json:"device_tokens"`
	Title        string                 `json:"title"`
	Body         string                 `json:"body"`
	Data         map[string]interface{} `json:"data"`
	Badge        int                    `json:"badge"`
	Sound        string                 `json:"sound"`
	Priority     string                 `json:"priority"`
	TTL          int                    `json:"ttl"`
	CollapseKey  string                 `json:"collapse_key,omitempty"`
	Platform     string                 `json:"platform"`
	Status       string                 `json:"status"`
	DeliveredAt  time.Time             `json:"delivered_at,omitempty"`
	ClickedAt    time.Time             `json:"clicked_at,omitempty"`
}

// AlertRule represents automated alert configuration
type AlertRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Conditions  map[string]interface{} `json:"conditions"`
	Actions     []AlertAction          `json:"actions"`
	Enabled     bool                   `json:"enabled"`
	Priority    string                 `json:"priority"`
	Cooldown    time.Duration          `json:"cooldown"`
	Recipients  []string               `json:"recipients"`
	CreatedBy   string                 `json:"created_by"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// AlertAction represents action to take when alert triggers
type AlertAction struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
	Delay      time.Duration          `json:"delay,omitempty"`
}

// NotificationTemplate represents reusable notification template
type NotificationTemplate struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Subject      string                 `json:"subject"`
	Body         string                 `json:"body"`
	Variables    []string               `json:"variables"`
	Localization map[string]interface{} `json:"localization"`
	Platforms    []string               `json:"platforms"`
	IsDefault    bool                   `json:"is_default"`
	CreatedBy    string                 `json:"created_by"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// Journey07NotificationsAPISuite tests comprehensive notification and alert systems
type Journey07NotificationsAPISuite struct {
	suite.Suite
	baseURL    string
	httpClient *http.Client
	user1      *AuthenticatedUser
	user2      *AuthenticatedUser
	admin      *AuthenticatedUser
	wsConn     *websocket.Conn
	deviceTokens []string
	templates  []NotificationTemplate
	alertRules []AlertRule
}

func (suite *Journey07NotificationsAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081"
	suite.httpClient = &http.Client{Timeout: 30 * time.Second}

	// Initialize test data
	suite.deviceTokens = []string{
		"mock_ios_token_12345",
		"mock_android_token_67890",
		"mock_web_token_abcdef",
	}

	// Create test users with different notification preferences
	suite.user1 = suite.createTestUser("notifications_user1@tchat.com", "+66812345001", "Asia/Bangkok")
	suite.user2 = suite.createTestUser("notifications_user2@tchat.com", "+66812345002", "Asia/Bangkok")
	suite.admin = suite.createTestUser("notifications_admin@tchat.com", "+66812345003", "Asia/Bangkok")

	// Set up notification preferences for test users
	suite.setupNotificationPreferences()

	// Create notification templates
	suite.createNotificationTemplates()

	// Set up alert rules
	suite.setupAlertRules()

	// Establish WebSocket connection for real-time notifications
	suite.connectWebSocket()
}

func (suite *Journey07NotificationsAPISuite) TearDownSuite() {
	if suite.wsConn != nil {
		suite.wsConn.Close()
	}
	suite.cleanupTestData()
}

func (suite *Journey07NotificationsAPISuite) createTestUser(email, phoneNumber, timezone string) *AuthenticatedUser {
	registerData := map[string]interface{}{
		"phone_number": phoneNumber,
		"email":        email,
		"password":     "SecurePassword123!",
		"firstName":    "Test",
		"lastName":     "User",
		"country":      "TH",
		"language":     "en",
		"timezone":     timezone,
	}

	jsonData, _ := json.Marshal(registerData)
	resp, err := suite.httpClient.Post(
		fmt.Sprintf("%s/api/v1/auth/register", suite.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	suite.NoError(err)
	defer resp.Body.Close()

	var registerResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&registerResponse)
	suite.NoError(err)

	// Check if registration was successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		suite.FailNow("Registration failed", "Status: %d, Response: %+v", resp.StatusCode, registerResponse)
	}

	// Safely extract user_id and token with validation
	userID, userIDOk := registerResponse["user_id"].(string)
	token, tokenOk := registerResponse["token"].(string)

	if !userIDOk || !tokenOk {
		suite.FailNow("Invalid registration response", "Expected user_id and token, got: %+v", registerResponse)
	}

	return &AuthenticatedUser{
		UserID:      userID,
		AccessToken: token,
		Email:       email,
	}
}

func (suite *Journey07NotificationsAPISuite) setupNotificationPreferences() {
	// Configure comprehensive notification preferences for user1
	preferences := NotificationPreferences{
		UserID:       suite.user1.UserID,
		EmailEnabled: true,
		PushEnabled:  true,
		SMSEnabled:   false,
		InAppEnabled: true,
		Categories: map[string]bool{
			"messages":     true,
			"social":       true,
			"commerce":     false,
			"content":      true,
			"security":     true,
			"system":       false,
			"marketing":    false,
		},
		Schedules: map[string]string{
			"immediate": "0 * * * *",
			"digest":    "0 9 * * *",
		},
		Platforms: map[string]bool{
			"web":     true,
			"ios":     true,
			"android": false,
		},
		Languages:   []string{"en", "th", "id"},
		TimeZone:    "Asia/Jakarta",
		QuietHours: map[string]string{
			"start": "22:00",
			"end":   "07:00",
		},
		FrequencyLimits: map[string]int{
			"marketing": 2,
			"social":    10,
			"messages":  100,
		},
	}

	jsonData, _ := json.Marshal(preferences)
	req, _ := http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/notifications/preferences", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	_, err := suite.httpClient.Do(req)
	suite.NoError(err)
}

func (suite *Journey07NotificationsAPISuite) createNotificationTemplates() {
	templates := []NotificationTemplate{
		{
			Name:    "Welcome Message",
			Type:    "system",
			Subject: "Welcome to Tchat!",
			Body:    "Welcome {{name}}! Your account is ready to use.",
			Variables: []string{"name", "email"},
			Localization: map[string]interface{}{
				"th": map[string]string{
					"subject": "ยินดีต้อนรับสู่ Tchat!",
					"body":    "ยินดีต้อนรับ {{name}}! บัญชีของคุณพร้อมใช้งานแล้ว",
				},
				"id": map[string]string{
					"subject": "Selamat datang di Tchat!",
					"body":    "Selamat datang {{name}}! Akun Anda siap digunakan.",
				},
			},
			Platforms: []string{"web", "ios", "android", "email"},
			IsDefault: true,
		},
		{
			Name:    "New Message Alert",
			Type:    "message",
			Subject: "New message from {{sender_name}}",
			Body:    "{{sender_name}}: {{message_preview}}",
			Variables: []string{"sender_name", "message_preview", "chat_id"},
			Platforms: []string{"push", "email", "sms"},
		},
		{
			Name:    "Security Alert",
			Type:    "security",
			Subject: "Security Alert: {{alert_type}}",
			Body:    "Security event detected: {{description}} at {{timestamp}}",
			Variables: []string{"alert_type", "description", "timestamp", "location"},
			Platforms: []string{"email", "sms", "push"},
		},
	}

	for _, template := range templates {
		jsonData, _ := json.Marshal(template)
		req, _ := http.NewRequest("POST",
			fmt.Sprintf("%s/api/v1/notifications/templates", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		_, err := suite.httpClient.Do(req)
		suite.NoError(err)
		suite.templates = append(suite.templates, template)
	}
}

func (suite *Journey07NotificationsAPISuite) setupAlertRules() {
	rules := []AlertRule{
		{
			Name: "Failed Login Attempts",
			Type: "security",
			Conditions: map[string]interface{}{
				"event_type": "failed_login",
				"threshold":  5,
				"timeframe":  "5m",
			},
			Actions: []AlertAction{
				{
					Type: "notification",
					Parameters: map[string]interface{}{
						"template": "security_alert",
						"priority": "high",
						"channels": []string{"email", "sms"},
					},
				},
				{
					Type: "block_ip",
					Parameters: map[string]interface{}{
						"duration": "1h",
					},
					Delay: 5 * time.Minute,
				},
			},
			Enabled:   true,
			Priority:  "high",
			Cooldown:  15 * time.Minute,
			Recipients: []string{suite.admin.UserID},
		},
		{
			Name: "High Error Rate",
			Type: "system",
			Conditions: map[string]interface{}{
				"metric":    "error_rate",
				"threshold": 5.0,
				"timeframe": "10m",
			},
			Actions: []AlertAction{
				{
					Type: "notification",
					Parameters: map[string]interface{}{
						"template": "system_alert",
						"priority": "critical",
						"channels": []string{"slack", "email"},
					},
				},
			},
			Enabled:   true,
			Priority:  "critical",
			Cooldown:  5 * time.Minute,
		},
	}

	for _, rule := range rules {
		jsonData, _ := json.Marshal(rule)
		req, _ := http.NewRequest("POST",
			fmt.Sprintf("%s/api/v1/notifications/alerts/rules", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		_, err := suite.httpClient.Do(req)
		suite.NoError(err)
		suite.alertRules = append(suite.alertRules, rule)
	}
}

func (suite *Journey07NotificationsAPISuite) connectWebSocket() {
	wsURL := strings.Replace(suite.baseURL, "http://", "ws://", 1)
	wsURL += "/api/v1/notifications/ws?token=" + suite.user1.AccessToken

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		// WebSocket connection optional for testing
		return
	}
	suite.wsConn = conn
}

// Test notification creation and delivery
func (suite *Journey07NotificationsAPISuite) TestCreateAndDeliverNotification() {
	notification := NotificationInfo{
		UserID:   suite.user2.UserID,
		Type:     "message",
		Title:    "New Message",
		Message:  "You have a new message from " + suite.user1.Email,
		Category: "messages",
		Priority: "normal",
		PlatformTargets: []string{"web", "push"},
		Metadata: map[string]interface{}{
			"sender_id":    suite.user1.UserID,
			"message_type": "text",
		},
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Create notification
	jsonData, _ := json.Marshal(notification)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/notifications", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var createResponse NotificationInfo
	err = json.NewDecoder(resp.Body).Decode(&createResponse)
	suite.NoError(err)
	suite.NotEmpty(createResponse.ID)
	suite.Equal("pending", createResponse.Status)

	// Verify notification appears in user's inbox
	time.Sleep(2 * time.Second) // Allow processing time

	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/notifications/inbox", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user2.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var inboxResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&inboxResponse)

	notifications := inboxResponse["notifications"].([]interface{})
	suite.Greater(len(notifications), 0)

	foundNotification := false
	for _, notif := range notifications {
		notifMap := notif.(map[string]interface{})
		if notifID, ok := notifMap["id"].(string); ok && notifID == createResponse.ID {
			foundNotification = true
			suite.Equal("delivered", notifMap["status"])
			break
		}
	}
	suite.True(foundNotification)
}

// Test push notification delivery
func (suite *Journey07NotificationsAPISuite) TestPushNotificationDelivery() {
	// Register device token
	deviceData := map[string]interface{}{
		"token":    suite.deviceTokens[0],
		"platform": "ios",
		"app_version": "1.0.0",
		"os_version":  "16.0",
	}

	jsonData, _ := json.Marshal(deviceData)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/notifications/devices", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	_, err := suite.httpClient.Do(req)
	suite.NoError(err)

	// Send push notification
	pushNotification := PushNotificationInfo{
		DeviceTokens: []string{suite.deviceTokens[0]},
		Title:        "Test Push Notification",
		Body:         "This is a test push notification",
		Data: map[string]interface{}{
			"type":    "test",
			"user_id": suite.user1.UserID,
		},
		Badge:    1,
		Sound:    "default",
		Priority: "high",
		TTL:      3600,
		Platform: "ios",
	}

	jsonData, _ = json.Marshal(pushNotification)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/notifications/push", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusAccepted, resp.StatusCode)

	var pushResponse PushNotificationInfo
	json.NewDecoder(resp.Body).Decode(&pushResponse)
	suite.NotEmpty(pushResponse.ID)
	suite.Equal("sent", pushResponse.Status)
}

// Test notification preferences management
func (suite *Journey07NotificationsAPISuite) TestNotificationPreferencesManagement() {
	// Get current preferences
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/notifications/preferences", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var preferences NotificationPreferences
	json.NewDecoder(resp.Body).Decode(&preferences)

	suite.Equal(suite.user1.UserID, preferences.UserID)
	suite.True(preferences.EmailEnabled)
	suite.True(preferences.PushEnabled)

	// Update preferences
	preferences.SMSEnabled = true
	preferences.Categories["commerce"] = true
	preferences.FrequencyLimits["marketing"] = 5

	jsonData, _ := json.Marshal(preferences)
	req, _ = http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/notifications/preferences", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	// Verify updates
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/notifications/preferences", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var updatedPreferences NotificationPreferences
	json.NewDecoder(resp.Body).Decode(&updatedPreferences)

	suite.True(updatedPreferences.SMSEnabled)
	suite.True(updatedPreferences.Categories["commerce"])
	suite.Equal(5, updatedPreferences.FrequencyLimits["marketing"])
}

// Test notification template system
func (suite *Journey07NotificationsAPISuite) TestNotificationTemplates() {
	// List available templates
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/notifications/templates", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var templatesResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&templatesResponse)

	templates := templatesResponse["templates"].([]interface{})
	suite.Greater(len(templates), 0)

	// Send notification using template
	templateNotification := map[string]interface{}{
		"user_id":     suite.user2.UserID,
		"template":    "welcome_message",
		"variables": map[string]interface{}{
			"name":  "Test User",
			"email": suite.user2.Email,
		},
		"platforms": []string{"email", "push"},
		"priority":  "normal",
	}

	jsonData, _ := json.Marshal(templateNotification)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/notifications/template", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)
}

// Test alert system functionality
func (suite *Journey07NotificationsAPISuite) TestAlertSystemFunctionality() {
	// Trigger a security alert by simulating failed login attempts
	for i := 0; i < 6; i++ {
		loginData := map[string]string{
			"email":    suite.user1.Email,
			"password": "wrong_password",
		}

		jsonData, _ := json.Marshal(loginData)
		_, err := suite.httpClient.Post(
			fmt.Sprintf("%s/api/v1/auth/login", suite.baseURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		suite.NoError(err)
	}

	// Wait for alert processing
	time.Sleep(5 * time.Second)

	// Check alert history
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/notifications/alerts/history", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var alertResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&alertResponse)

	alerts := alertResponse["alerts"].([]interface{})
	suite.Greater(len(alerts), 0)

	// Verify security alert was triggered
	foundSecurityAlert := false
	for _, alert := range alerts {
		alertMap := alert.(map[string]interface{})
		if alertMap["type"].(string) == "security" {
			foundSecurityAlert = true
			suite.Equal("triggered", alertMap["status"])
			break
		}
	}
	suite.True(foundSecurityAlert)
}

// Test real-time notification delivery via WebSocket
func (suite *Journey07NotificationsAPISuite) TestRealTimeNotificationDelivery() {
	if suite.wsConn == nil {
		suite.T().Skip("WebSocket connection not available")
	}

	// Send notification that should trigger real-time delivery
	notification := NotificationInfo{
		UserID:   suite.user1.UserID,
		Type:     "system",
		Title:    "Real-time Test",
		Message:  "This is a real-time notification test",
		Category: "system",
		Priority: "high",
		PlatformTargets: []string{"websocket"},
	}

	jsonData, _ := json.Marshal(notification)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/notifications", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	_, err := suite.httpClient.Do(req)
	suite.NoError(err)

	// Listen for WebSocket message
	suite.wsConn.SetReadDeadline(time.Now().Add(10 * time.Second))

	var wsMessage map[string]interface{}
	err = suite.wsConn.ReadJSON(&wsMessage)
	suite.NoError(err)

	suite.Equal("notification", wsMessage["type"])
	suite.Equal("Real-time Test", wsMessage["title"])
}

// Test notification analytics and reporting
func (suite *Journey07NotificationsAPISuite) TestNotificationAnalytics() {
	// Get notification statistics
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/notifications/analytics", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var analyticsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&analyticsResponse)

	suite.Contains(analyticsResponse, "total_sent")
	suite.Contains(analyticsResponse, "delivery_rate")
	suite.Contains(analyticsResponse, "open_rate")
	suite.Contains(analyticsResponse, "platform_breakdown")

	// Get user-specific analytics
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/notifications/analytics/user/%s", suite.baseURL, suite.user1.UserID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var userAnalytics map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userAnalytics)

	suite.Contains(userAnalytics, "received")
	suite.Contains(userAnalytics, "read")
	suite.Contains(userAnalytics, "categories")
}

// Test notification cleanup and archival
func (suite *Journey07NotificationsAPISuite) TestNotificationCleanupAndArchival() {
	// Mark old notifications as read
	req, _ := http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/notifications/mark-all-read", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	// Archive old notifications
	archiveData := map[string]interface{}{
		"older_than": "30d",
		"categories": []string{"marketing", "social"},
	}

	jsonData, _ := json.Marshal(archiveData)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/notifications/archive", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var archiveResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&archiveResponse)

	suite.Contains(archiveResponse, "archived_count")
}

func (suite *Journey07NotificationsAPISuite) cleanupTestData() {
	// Clean up test notifications, templates, and alert rules
	req, _ := http.NewRequest("DELETE",
		fmt.Sprintf("%s/api/v1/notifications/cleanup/test-data", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	_, _ = suite.httpClient.Do(req)
}

func TestJourney07NotificationsAPISuite(t *testing.T) {
	suite.Run(t, new(Journey07NotificationsAPISuite))
}