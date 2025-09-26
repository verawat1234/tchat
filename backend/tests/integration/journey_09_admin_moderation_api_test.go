// Journey 9: Admin & Moderation API Integration Tests
// Comprehensive testing of admin tools, content moderation, user management,
// safety systems, compliance features, and automated moderation workflows

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

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

// AdminUser represents an administrative user with elevated permissions
type AdminUser struct {
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	Region      string    `json:"region"`
	CreatedAt   time.Time `json:"created_at"`
	LastActive  time.Time `json:"last_active"`
	Status      string    `json:"status"`
}

// ModerationAction represents a moderation action taken by admin
type ModerationAction struct {
	ID           string                 `json:"id,omitempty"`
	ActionType   string                 `json:"action_type"`
	TargetType   string                 `json:"target_type"`
	TargetID     string                 `json:"target_id"`
	Reason       string                 `json:"reason"`
	Severity     string                 `json:"severity"`
	ModeratorID  string                 `json:"moderator_id"`
	Automated    bool                   `json:"automated"`
	Timestamp    time.Time             `json:"timestamp"`
	ExpiresAt    time.Time             `json:"expires_at,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Status       string                 `json:"status"`
	ReviewedBy   string                 `json:"reviewed_by,omitempty"`
	ReviewedAt   time.Time             `json:"reviewed_at,omitempty"`
}

// ContentReview represents a content moderation review
type ContentReview struct {
	ID              string                 `json:"id,omitempty"`
	ContentID       string                 `json:"content_id"`
	ContentType     string                 `json:"content_type"`
	ReporterID      string                 `json:"reporter_id"`
	ReportReason    string                 `json:"report_reason"`
	ReviewerID      string                 `json:"reviewer_id,omitempty"`
	Status          string                 `json:"status"`
	Priority        string                 `json:"priority"`
	AutoModerated   bool                   `json:"auto_moderated"`
	ToxicityScore   float64               `json:"toxicity_score,omitempty"`
	ViolationTypes  []string               `json:"violation_types,omitempty"`
	Evidence        map[string]interface{} `json:"evidence,omitempty"`
	Resolution      string                 `json:"resolution,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	ResolvedAt      time.Time             `json:"resolved_at,omitempty"`
	EscalatedAt     time.Time             `json:"escalated_at,omitempty"`
}

// UserSafetyProfile represents user safety and compliance data
type UserSafetyProfile struct {
	UserID              string               `json:"user_id"`
	TrustScore          float64              `json:"trust_score"`
	ViolationHistory    []ViolationRecord    `json:"violation_history"`
	RestrictionStatus   string               `json:"restriction_status"`
	RestrictedUntil     time.Time           `json:"restricted_until,omitempty"`
	WatchlistStatus     string               `json:"watchlist_status"`
	VerificationLevel   string               `json:"verification_level"`
	RiskAssessment      map[string]float64   `json:"risk_assessment"`
	ComplianceFlags     []string             `json:"compliance_flags"`
	RecentReports       int                  `json:"recent_reports"`
	AccountAge          time.Duration        `json:"account_age"`
	ActivityPattern     map[string]interface{} `json:"activity_pattern"`
	LastReview          time.Time           `json:"last_review"`
}

// ViolationRecord represents a single policy violation
type ViolationRecord struct {
	ID              string                 `json:"id"`
	ViolationType   string                 `json:"violation_type"`
	Severity        string                 `json:"severity"`
	Description     string                 `json:"description"`
	Evidence        map[string]interface{} `json:"evidence"`
	ActionTaken     string                 `json:"action_taken"`
	CreatedAt       time.Time             `json:"created_at"`
	Status          string                 `json:"status"`
}

// ComplianceReport represents a compliance reporting document
type ComplianceReport struct {
	ID              string                 `json:"id,omitempty"`
	ReportType      string                 `json:"report_type"`
	Period          string                 `json:"period"`
	Region          string                 `json:"region"`
	GeneratedAt     time.Time             `json:"generated_at"`
	Data            map[string]interface{} `json:"data"`
	Summary         string                 `json:"summary"`
	TotalViolations int                    `json:"total_violations"`
	ActionsTaken    int                    `json:"actions_taken"`
	Metrics         map[string]float64     `json:"metrics"`
	Status          string                 `json:"status"`
	GeneratedBy     string                 `json:"generated_by"`
}

// AutoModRule represents an automated moderation rule
type AutoModRule struct {
	ID              string                 `json:"id,omitempty"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	RuleType        string                 `json:"rule_type"`
	Conditions      map[string]interface{} `json:"conditions"`
	Actions         []string               `json:"actions"`
	Enabled         bool                   `json:"enabled"`
	Severity        string                 `json:"severity"`
	Confidence      float64               `json:"confidence_threshold"`
	CreatedBy       string                 `json:"created_by"`
	UpdatedAt       time.Time             `json:"updated_at"`
	TriggeredCount  int                    `json:"triggered_count"`
	FalsePositives  int                    `json:"false_positives"`
	Effectiveness   float64               `json:"effectiveness"`
}

// SystemAlert represents a system-generated alert for admins
type SystemAlert struct {
	ID          string                 `json:"id,omitempty"`
	AlertType   string                 `json:"alert_type"`
	Severity    string                 `json:"severity"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	ResolvedAt  time.Time             `json:"resolved_at,omitempty"`
	ResolvedBy  string                 `json:"resolved_by,omitempty"`
	Region      string                 `json:"region"`
	Priority    int                    `json:"priority"`
}

// Journey09AdminModerationAPISuite tests comprehensive admin and moderation systems
type Journey09AdminModerationAPISuite struct {
	suite.Suite
	baseURL        string
	httpClient     *http.Client
	superAdmin     *AuthenticatedUser
	moderator      *AuthenticatedUser
	regionalAdmin  *AuthenticatedUser
	testUser       *AuthenticatedUser
	violatingUser  *AuthenticatedUser
	testContent    []string
	moderationRules []AutoModRule
	pendingReviews  []ContentReview
}

func (suite *Journey09AdminModerationAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081"
	suite.httpClient = &http.Client{Timeout: 30 * time.Second}

	// Create test users with different admin roles
	suite.superAdmin = suite.createAdminUser("super@tchat.com", "superadmin", "super_admin", []string{"*"})
	suite.moderator = suite.createAdminUser("moderator@tchat.com", "moderator123", "content_moderator", []string{"moderate_content", "review_reports"})
	suite.regionalAdmin = suite.createAdminUser("admin_sea@tchat.com", "admin789", "regional_admin", []string{"manage_users", "view_analytics", "moderate_content"})

	// Create regular test users
	suite.testUser = suite.createTestUser("testuser@tchat.com", "password123")
	suite.violatingUser = suite.createTestUser("violator@tchat.com", "password456")

	// Set up moderation rules
	suite.setupAutoModerationRules()

	// Create test content for moderation
	suite.createTestContent()

	// Generate sample violation data
	suite.generateTestViolations()
}

func (suite *Journey09AdminModerationAPISuite) TearDownSuite() {
	suite.cleanupTestData()
}

func (suite *Journey09AdminModerationAPISuite) createAdminUser(email, password, role string, permissions []string) *AuthenticatedUser {
	registerData := map[string]interface{}{
		"email":       email,
		"password":    password,
		"role":        role,
		"permissions": permissions,
		"region":      "SEA",
	}

	jsonData, _ := json.Marshal(registerData)
	resp, err := suite.httpClient.Post(
		fmt.Sprintf("%s/api/v1/admin/users", suite.baseURL),
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

func (suite *Journey09AdminModerationAPISuite) createTestUser(email, password string) *AuthenticatedUser {
	registerData := map[string]interface{}{
		"email":    email,
		"password": password,
		"role":     "user",
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

func (suite *Journey09AdminModerationAPISuite) setupAutoModerationRules() {
	rules := []AutoModRule{
		{
			Name:        "Toxic Language Filter",
			Description: "Detects toxic and harmful language in content",
			RuleType:    "text_analysis",
			Conditions: map[string]interface{}{
				"toxicity_threshold": 0.8,
				"categories": []string{"toxicity", "severe_toxicity", "identity_attack", "insult"},
			},
			Actions:            []string{"flag_for_review", "hide_content", "warn_user"},
			Enabled:            true,
			Severity:           "high",
			Confidence:         0.85,
			CreatedBy:          suite.superAdmin.UserID,
		},
		{
			Name:        "Spam Detection",
			Description: "Identifies spam and promotional content",
			RuleType:    "behavior_analysis",
			Conditions: map[string]interface{}{
				"repetition_threshold": 5,
				"link_density":         0.3,
				"keyword_matching":     true,
			},
			Actions:   []string{"rate_limit", "flag_for_review"},
			Enabled:   true,
			Severity:  "medium",
			Confidence: 0.75,
			CreatedBy: suite.superAdmin.UserID,
		},
		{
			Name:        "NSFW Content Filter",
			Description: "Detects not-safe-for-work content",
			RuleType:    "image_analysis",
			Conditions: map[string]interface{}{
				"nsfw_threshold": 0.7,
				"violence_threshold": 0.6,
			},
			Actions:   []string{"blur_content", "age_restrict", "flag_for_review"},
			Enabled:   true,
			Severity:  "high",
			Confidence: 0.9,
			CreatedBy: suite.superAdmin.UserID,
		},
	}

	for _, rule := range rules {
		rule.UpdatedAt = time.Now()
		suite.moderationRules = append(suite.moderationRules, rule)
	}
}

func (suite *Journey09AdminModerationAPISuite) createTestContent() {
	contents := []map[string]interface{}{
		{
			"type":        "text",
			"content":     "This is normal, appropriate content for testing",
			"user_id":     suite.testUser.UserID,
		},
		{
			"type":        "text",
			"content":     "This content contains some potentially problematic language",
			"user_id":     suite.violatingUser.UserID,
		},
		{
			"type":        "image",
			"content":     "image_url_placeholder.jpg",
			"description": "Test image content",
			"user_id":     suite.testUser.UserID,
		},
	}

	for _, content := range contents {
		jsonData, _ := json.Marshal(content)
		req, _ := http.NewRequest("POST",
			fmt.Sprintf("%s/api/v1/content", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.testUser.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			var contentResp map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&contentResp)
			if contentID, ok := contentResp["id"].(string); ok {
				suite.testContent = append(suite.testContent, contentID)
			}
		}
	}
}

func (suite *Journey09AdminModerationAPISuite) generateTestViolations() {
	// Simulate some violation reports
	if len(suite.testContent) > 0 {
		review := ContentReview{
			ContentID:      suite.testContent[0],
			ContentType:    "text",
			ReporterID:     suite.testUser.UserID,
			ReportReason:   "inappropriate_content",
			Status:         "pending",
			Priority:       "medium",
			AutoModerated:  false,
			ToxicityScore:  0.65,
			ViolationTypes: []string{"harassment"},
			CreatedAt:      time.Now(),
		}
		suite.pendingReviews = append(suite.pendingReviews, review)
	}
}

// Test admin user management
func (suite *Journey09AdminModerationAPISuite) TestAdminUserManagement() {
	// List all users (admin only)
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/users", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	var usersResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&usersResponse)

	suite.Contains(usersResponse, "users")
	users := usersResponse["users"].([]interface{})
	suite.Greater(len(users), 0)

	// Get specific user details
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/users/%s", suite.baseURL, suite.testUser.UserID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var userResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userResponse)

	suite.Equal(suite.testUser.UserID, userResponse["user_id"])
	suite.Contains(userResponse, "email")
	suite.Contains(userResponse, "status")

	// Update user status
	statusUpdate := map[string]interface{}{
		"status": "restricted",
		"reason": "Policy violation testing",
		"duration": "24h",
	}

	jsonData, _ := json.Marshal(statusUpdate)
	req, _ = http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/admin/users/%s/status", suite.baseURL, suite.violatingUser.UserID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.moderator.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

// Test content moderation workflows
func (suite *Journey09AdminModerationAPISuite) TestContentModerationWorkflows() {
	// Create moderation rules
	for _, rule := range suite.moderationRules {
		jsonData, _ := json.Marshal(rule)
		req, _ := http.NewRequest("POST",
			fmt.Sprintf("%s/api/v1/admin/moderation/rules", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusCreated, resp.StatusCode)
	}

	// Get pending content reviews
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/moderation/reviews?status=pending", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.moderator.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var reviewsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&reviewsResponse)

	suite.Contains(reviewsResponse, "reviews")

	// Submit content report
	if len(suite.testContent) > 0 {
		report := map[string]interface{}{
			"content_id":    suite.testContent[0],
			"content_type":  "text",
			"reason":        "inappropriate_content",
			"description":   "This content violates community guidelines",
			"evidence_urls": []string{"screenshot1.jpg"},
		}

		jsonData, _ := json.Marshal(report)
		req, _ := http.NewRequest("POST",
			fmt.Sprintf("%s/api/v1/admin/moderation/reports", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.testUser.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusCreated, resp.StatusCode)

		var reportResponse ContentReview
		json.NewDecoder(resp.Body).Decode(&reportResponse)
		suite.NotEmpty(reportResponse.ID)
		suite.Equal("pending", reportResponse.Status)

		// Review and resolve the report
		time.Sleep(2 * time.Second)

		resolution := map[string]interface{}{
			"action":     "warn_user",
			"resolution": "content_warning_applied",
			"notes":      "Applied content warning, user notified",
		}

		jsonData, _ = json.Marshal(resolution)
		req, _ = http.NewRequest("PUT",
			fmt.Sprintf("%s/api/v1/admin/moderation/reviews/%s", suite.baseURL, reportResponse.ID),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.moderator.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err = suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusOK, resp.StatusCode)
	}
}

// Test automated moderation systems
func (suite *Journey09AdminModerationAPISuite) TestAutomatedModerationSystems() {
	// Test content analysis API
	testContent := map[string]interface{}{
		"text":    "This is a test message to analyze for moderation",
		"type":    "text",
		"user_id": suite.testUser.UserID,
	}

	jsonData, _ := json.Marshal(testContent)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/admin/moderation/analyze", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.moderator.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var analysisResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&analysisResponse)

	suite.Contains(analysisResponse, "toxicity_score")
	suite.Contains(analysisResponse, "violation_probability")
	suite.Contains(analysisResponse, "recommended_actions")

	// Test rule trigger simulation
	triggerData := map[string]interface{}{
		"rule_id": "test_rule_001",
		"content": map[string]interface{}{
			"text":    "Simulated toxic content for rule testing",
			"user_id": suite.violatingUser.UserID,
		},
		"simulate": true,
	}

	jsonData, _ = json.Marshal(triggerData)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/admin/moderation/rules/test", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var testResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&testResponse)

	suite.Contains(testResponse, "would_trigger")
	suite.Contains(testResponse, "confidence_score")
	suite.Contains(testResponse, "suggested_actions")
}

// Test user safety profiles and trust scoring
func (suite *Journey09AdminModerationAPISuite) TestUserSafetyProfiles() {
	// Get user safety profile
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/safety/users/%s", suite.baseURL, suite.violatingUser.UserID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.moderator.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var safetyProfile UserSafetyProfile
	json.NewDecoder(resp.Body).Decode(&safetyProfile)

	suite.Equal(suite.violatingUser.UserID, safetyProfile.UserID)
	suite.GreaterOrEqual(safetyProfile.TrustScore, 0.0)
	suite.LessOrEqual(safetyProfile.TrustScore, 1.0)
	suite.Contains(safetyProfile.RiskAssessment, "spam_risk")
	suite.Contains(safetyProfile.RiskAssessment, "abuse_risk")

	// Update trust score
	scoreUpdate := map[string]interface{}{
		"trust_score_delta": -0.1,
		"reason":           "Policy violation during testing",
		"reviewer_id":      suite.moderator.UserID,
	}

	jsonData, _ := json.Marshal(scoreUpdate)
	req, _ = http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/admin/safety/users/%s/trust-score", suite.baseURL, suite.violatingUser.UserID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.moderator.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)

	// Add user to watchlist
	watchlistEntry := map[string]interface{}{
		"reason":      "Multiple policy violations",
		"priority":    "medium",
		"alert_level": "standard",
		"duration":    "30d",
		"notes":       "Monitor for continued violations",
	}

	jsonData, _ = json.Marshal(watchlistEntry)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/admin/safety/watchlist", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.moderator.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)
}

// Test compliance reporting and audit trails
func (suite *Journey09AdminModerationAPISuite) TestComplianceReportingAndAuditTrails() {
	// Generate compliance report
	reportRequest := map[string]interface{}{
		"report_type": "transparency_report",
		"period":      "monthly",
		"region":      "SEA",
		"include_sections": []string{
			"content_removals",
			"user_suspensions",
			"policy_violations",
			"appeals_processed",
		},
	}

	jsonData, _ := json.Marshal(reportRequest)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/admin/compliance/reports", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var reportResponse ComplianceReport
	json.NewDecoder(resp.Body).Decode(&reportResponse)
	suite.NotEmpty(reportResponse.ID)
	suite.Equal("transparency_report", reportResponse.ReportType)

	// Get audit trail for user actions
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/audit/users/%s?limit=50", suite.baseURL, suite.violatingUser.UserID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var auditResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&auditResponse)

	suite.Contains(auditResponse, "audit_logs")
	suite.Contains(auditResponse, "total_count")

	// Get moderation action history
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/moderation/actions?moderator=%s", suite.baseURL, suite.moderator.UserID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var actionsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&actionsResponse)

	suite.Contains(actionsResponse, "actions")
	suite.Contains(actionsResponse, "moderator_stats")
}

// Test system alerts and monitoring
func (suite *Journey09AdminModerationAPISuite) TestSystemAlertsAndMonitoring() {
	// Get active system alerts
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/alerts?status=active", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var alertsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&alertsResponse)

	suite.Contains(alertsResponse, "alerts")

	// Create test system alert
	testAlert := SystemAlert{
		AlertType: "content_spike",
		Severity:  "warning",
		Title:     "Unusual Content Volume",
		Message:   "Detected 50% increase in content reports in the last hour",
		Data: map[string]interface{}{
			"current_rate":  150,
			"baseline_rate": 100,
			"percentage_increase": 50,
		},
		Status:   "active",
		Region:   "SEA",
		Priority: 2,
	}

	jsonData, _ := json.Marshal(testAlert)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/admin/alerts", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var alertResponse SystemAlert
	json.NewDecoder(resp.Body).Decode(&alertResponse)
	suite.NotEmpty(alertResponse.ID)

	// Mark alert as resolved
	resolution := map[string]interface{}{
		"resolution_notes": "Content volume normalized after temporary spike",
		"resolved_by":      suite.superAdmin.UserID,
	}

	jsonData, _ = json.Marshal(resolution)
	req, _ = http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/admin/alerts/%s/resolve", suite.baseURL, alertResponse.ID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

// Test admin dashboard and analytics
func (suite *Journey09AdminModerationAPISuite) TestAdminDashboardAndAnalytics() {
	// Get moderation dashboard data
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/dashboard/moderation", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var dashboardResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&dashboardResponse)

	suite.Contains(dashboardResponse, "pending_reviews")
	suite.Contains(dashboardResponse, "auto_mod_actions")
	suite.Contains(dashboardResponse, "moderator_workload")
	suite.Contains(dashboardResponse, "violation_trends")

	// Get regional moderation stats
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/stats/moderation?region=SEA&period=7d", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.regionalAdmin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var statsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&statsResponse)

	suite.Contains(statsResponse, "total_actions")
	suite.Contains(statsResponse, "content_removals")
	suite.Contains(statsResponse, "user_sanctions")
	suite.Contains(statsResponse, "response_times")

	// Get effectiveness metrics for auto-moderation rules
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/admin/moderation/rules/metrics", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var metricsResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&metricsResponse)

	suite.Contains(metricsResponse, "rule_effectiveness")
	suite.Contains(metricsResponse, "false_positive_rates")
	suite.Contains(metricsResponse, "coverage_analysis")
}

// Test policy management and appeals process
func (suite *Journey09AdminModerationAPISuite) TestPolicyManagementAndAppeals() {
	// Create policy violation record
	violation := map[string]interface{}{
		"user_id":        suite.violatingUser.UserID,
		"violation_type": "harassment",
		"severity":       "medium",
		"description":    "Inappropriate comment towards another user",
		"evidence": map[string]interface{}{
			"content_id": "content_123",
			"report_id":  "report_456",
		},
		"action_taken": "content_removal",
	}

	jsonData, _ := json.Marshal(violation)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/admin/violations", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.moderator.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var violationResponse ViolationRecord
	json.NewDecoder(resp.Body).Decode(&violationResponse)
	suite.NotEmpty(violationResponse.ID)

	// Submit appeal
	appeal := map[string]interface{}{
		"violation_id": violationResponse.ID,
		"reason":       "False positive - content was taken out of context",
		"explanation":  "The flagged content was a quote from a movie, not harassment",
		"evidence_urls": []string{"context_screenshot.jpg"},
	}

	jsonData, _ = json.Marshal(appeal)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/admin/appeals", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.violatingUser.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var appealResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&appealResponse)
	suite.NoError(err)

	appealID, ok := appealResponse["appeal_id"].(string)
	suite.True(ok, "Expected appeal_id in response: %+v", appealResponse)

	// Review appeal
	time.Sleep(1 * time.Second)

	appealDecision := map[string]interface{}{
		"decision":   "upheld",
		"reasoning":  "Upon review, the original decision was correct",
		"reviewer_id": suite.superAdmin.UserID,
	}

	jsonData, _ = json.Marshal(appealDecision)
	req, _ = http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/admin/appeals/%s", suite.baseURL, appealID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *Journey09AdminModerationAPISuite) cleanupTestData() {
	// Clean up test moderation data, users, and rules
	req, _ := http.NewRequest("DELETE",
		fmt.Sprintf("%s/api/v1/admin/cleanup/test-data", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.superAdmin.AccessToken)

	_, _ = suite.httpClient.Do(req)
}

func TestJourney09AdminModerationAPISuite(t *testing.T) {
	suite.Run(t, new(Journey09AdminModerationAPISuite))
}