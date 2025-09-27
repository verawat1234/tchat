// Journey 8: Analytics & Insights API Integration Tests
// Comprehensive testing of analytics services, business intelligence, user behavior tracking,
// performance metrics, revenue analytics, and real-time dashboard data

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
// Note: AuthenticatedUser is now defined in types.go

// AnalyticsEvent represents a trackable user event
type AnalyticsEvent struct {
	ID           string                 `json:"id,omitempty"`
	UserID       string                 `json:"user_id"`
	SessionID    string                 `json:"session_id"`
	EventType    string                 `json:"event_type"`
	EventName    string                 `json:"event_name"`
	Properties   map[string]interface{} `json:"properties"`
	Timestamp    time.Time             `json:"timestamp"`
	Platform     string                 `json:"platform"`
	AppVersion   string                 `json:"app_version"`
	DeviceType   string                 `json:"device_type"`
	Location     string                 `json:"location,omitempty"`
	Referrer     string                 `json:"referrer,omitempty"`
	UTMSource    string                 `json:"utm_source,omitempty"`
	UTMCampaign  string                 `json:"utm_campaign,omitempty"`
	Value        float64               `json:"value,omitempty"`
}

// UserBehaviorMetrics represents user engagement metrics
type UserBehaviorMetrics struct {
	UserID              string    `json:"user_id"`
	SessionCount        int       `json:"session_count"`
	TotalSessionTime    int       `json:"total_session_time"`
	AverageSessionTime  float64   `json:"average_session_time"`
	PageViews           int       `json:"page_views"`
	UniquePages         int       `json:"unique_pages"`
	BounceRate          float64   `json:"bounce_rate"`
	RetentionRate       float64   `json:"retention_rate"`
	LastActiveAt        time.Time `json:"last_active_at"`
	LifetimeValue       float64   `json:"lifetime_value"`
	ConversionEvents    int       `json:"conversion_events"`
	EngagementScore     float64   `json:"engagement_score"`
	ChurnProbability    float64   `json:"churn_probability"`
	PreferredPlatform   string    `json:"preferred_platform"`
	ActivityHeatmap     map[string]int `json:"activity_heatmap"`
}

// BusinessMetrics represents key business intelligence data
type BusinessMetrics struct {
	Period              string  `json:"period"`
	ActiveUsers         int     `json:"active_users"`
	NewUsers            int     `json:"new_users"`
	RetainedUsers       int     `json:"retained_users"`
	ChurnedUsers        int     `json:"churned_users"`
	Revenue             float64 `json:"revenue"`
	AverageOrderValue   float64 `json:"average_order_value"`
	ConversionRate      float64 `json:"conversion_rate"`
	CustomerAcquisitionCost float64 `json:"customer_acquisition_cost"`
	LifetimeValue       float64 `json:"lifetime_value"`
	GrowthRate          float64 `json:"growth_rate"`
	EngagementRate      float64 `json:"engagement_rate"`
	RetentionRate       float64 `json:"retention_rate"`
	ChurnRate           float64 `json:"churn_rate"`
}

// PerformanceMetrics represents system performance analytics
// Note: PerformanceMetrics is now defined in types.go

// RevenueAnalytics represents financial performance data
type RevenueAnalytics struct {
	Period                string                    `json:"period"`
	TotalRevenue          float64                   `json:"total_revenue"`
	RecurringRevenue      float64                   `json:"recurring_revenue"`
	OneTimeRevenue        float64                   `json:"one_time_revenue"`
	RevenueByRegion       map[string]float64        `json:"revenue_by_region"`
	RevenueByProduct      map[string]float64        `json:"revenue_by_product"`
	RevenueByPlatform     map[string]float64        `json:"revenue_by_platform"`
	PaymentMethodBreakdown map[string]float64       `json:"payment_method_breakdown"`
	RefundRate            float64                   `json:"refund_rate"`
	ChargebackRate        float64                   `json:"chargeback_rate"`
	AverageTransactionValue float64                 `json:"average_transaction_value"`
	TransactionCount      int                       `json:"transaction_count"`
	GrossMargin           float64                   `json:"gross_margin"`
	NetMargin             float64                   `json:"net_margin"`
}

// ContentPerformanceAnalytics represents content performance metrics
type ContentPerformanceAnalytics struct {
	ContentID       string                 `json:"content_id"`
	ContentType     string                 `json:"content_type"`
	Views           int                    `json:"views"`
	UniqueViews     int                    `json:"unique_views"`
	Likes           int                    `json:"likes"`
	Shares          int                    `json:"shares"`
	Comments        int                    `json:"comments"`
	EngagementRate  float64               `json:"engagement_rate"`
	AverageWatchTime float64              `json:"average_watch_time"`
	CompletionRate  float64               `json:"completion_rate"`
	ClickThroughRate float64              `json:"click_through_rate"`
	ConversionRate  float64               `json:"conversion_rate"`
	RevenueGenerated float64              `json:"revenue_generated"`
	DemographicBreakdown map[string]interface{} `json:"demographic_breakdown"`
	GeographicBreakdown map[string]int     `json:"geographic_breakdown"`
	PlatformBreakdown map[string]int       `json:"platform_breakdown"`
	CreatedAt       time.Time             `json:"created_at"`
	LastUpdated     time.Time             `json:"last_updated"`
}

// DashboardWidget represents real-time dashboard component
type DashboardWidget struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	DataSource  string                 `json:"data_source"`
	Config      map[string]interface{} `json:"config"`
	RefreshRate int                    `json:"refresh_rate_seconds"`
	Position    map[string]int         `json:"position"`
	Size        map[string]int         `json:"size"`
	Filters     map[string]interface{} `json:"filters"`
	LastUpdated time.Time             `json:"last_updated"`
}

// AnalyticsQuery represents a custom analytics query
type AnalyticsQuery struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Query       string                 `json:"query"`
	Parameters  map[string]interface{} `json:"parameters"`
	Schedule    string                 `json:"schedule,omitempty"`
	OutputFormat string                `json:"output_format"`
	Recipients  []string               `json:"recipients,omitempty"`
	CreatedBy   string                 `json:"created_by"`
	CreatedAt   time.Time             `json:"created_at,omitempty"`
}

// Journey08AnalyticsAPISuite tests comprehensive analytics and insights systems
type Journey08AnalyticsAPISuite struct {
	suite.Suite
	baseURL      string
	httpClient   *http.Client
	user1        *AuthenticatedUser
	user2        *AuthenticatedUser
	admin        *AuthenticatedUser
	analyst      *AuthenticatedUser
	testEvents   []AnalyticsEvent
	testQueries  []AnalyticsQuery
	dashboardWidgets []DashboardWidget
}

func (suite *Journey08AnalyticsAPISuite) SetupSuite() {
	suite.baseURL = "http://localhost:8081"
	suite.httpClient = &http.Client{Timeout: 30 * time.Second}

	// Create test users with different roles
	suite.user1 = suite.createTestUser("analytics_user1@tchat.com", "password123", "user")
	suite.user2 = suite.createTestUser("analytics_user2@tchat.com", "password456", "user")
	suite.admin = suite.createTestUser("analytics_admin@tchat.com", "admin789", "admin")
	suite.analyst = suite.createTestUser("analytics_analyst@tchat.com", "analyst123", "analyst")

	// Generate sample analytics events
	suite.generateSampleEvents()

	// Create custom analytics queries
	suite.createAnalyticsQueries()

	// Set up dashboard widgets
	suite.setupDashboardWidgets()
}

func (suite *Journey08AnalyticsAPISuite) TearDownSuite() {
	suite.cleanupTestData()
}

func (suite *Journey08AnalyticsAPISuite) createTestUser(email, password, role string) *AuthenticatedUser {
	registerData := map[string]interface{}{
		"email":             email,
		"password":          password,
		"firstName":         "Test",
		"lastName":          "User",
		"country":           "TH",
		"language":          "en",
		"role":              role,
		"analytics_consent": true,
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

func (suite *Journey08AnalyticsAPISuite) generateSampleEvents() {
	events := []AnalyticsEvent{
		{
			UserID:      suite.user1.UserID,
			SessionID:   "session_001",
			EventType:   "page_view",
			EventName:   "homepage_view",
			Platform:    "web",
			AppVersion:  "1.0.0",
			DeviceType:  "desktop",
			Location:    "Singapore",
			Properties: map[string]interface{}{
				"page_url":    "/home",
				"referrer":    "google.com",
				"load_time":   1.2,
			},
		},
		{
			UserID:      suite.user1.UserID,
			SessionID:   "session_001",
			EventType:   "user_action",
			EventName:   "button_click",
			Platform:    "web",
			AppVersion:  "1.0.0",
			DeviceType:  "desktop",
			Properties: map[string]interface{}{
				"button_name": "subscribe",
				"page_url":    "/pricing",
				"position":    "header",
			},
		},
		{
			UserID:      suite.user2.UserID,
			SessionID:   "session_002",
			EventType:   "commerce",
			EventName:   "purchase_completed",
			Platform:    "mobile",
			AppVersion:  "1.0.0",
			DeviceType:  "smartphone",
			Location:    "Thailand",
			Value:       99.99,
			Properties: map[string]interface{}{
				"product_id":    "prod_123",
				"category":      "premium",
				"payment_method": "credit_card",
				"currency":      "THB",
				"amount":        2999.70,
			},
		},
	}

	for _, event := range events {
		event.Timestamp = time.Now()
		suite.testEvents = append(suite.testEvents, event)
	}
}

func (suite *Journey08AnalyticsAPISuite) createAnalyticsQueries() {
	queries := []AnalyticsQuery{
		{
			Name:  "Daily Active Users",
			Query: "SELECT COUNT(DISTINCT user_id) as dau FROM events WHERE event_type = 'page_view' AND date = ?",
			Parameters: map[string]interface{}{
				"date": time.Now().Format("2006-01-02"),
			},
			OutputFormat: "json",
			CreatedBy:    suite.analyst.UserID,
		},
		{
			Name:  "Revenue by Region",
			Query: "SELECT location, SUM(value) as revenue FROM events WHERE event_name = 'purchase_completed' GROUP BY location",
			Parameters: map[string]interface{}{
				"start_date": time.Now().AddDate(0, 0, -30).Format("2006-01-02"),
				"end_date":   time.Now().Format("2006-01-02"),
			},
			OutputFormat: "json",
			CreatedBy:    suite.analyst.UserID,
		},
	}

	for _, query := range queries {
		query.CreatedAt = time.Now()
		suite.testQueries = append(suite.testQueries, query)
	}
}

func (suite *Journey08AnalyticsAPISuite) setupDashboardWidgets() {
	widgets := []DashboardWidget{
		{
			Type:        "metric_card",
			Title:       "Active Users",
			DataSource:  "analytics.daily_active_users",
			RefreshRate: 300, // 5 minutes
			Position:    map[string]int{"x": 0, "y": 0},
			Size:        map[string]int{"width": 4, "height": 2},
			Config: map[string]interface{}{
				"metric_type": "count",
				"time_range": "24h",
			},
		},
		{
			Type:        "line_chart",
			Title:       "Revenue Trend",
			DataSource:  "analytics.revenue_over_time",
			RefreshRate: 600, // 10 minutes
			Position:    map[string]int{"x": 4, "y": 0},
			Size:        map[string]int{"width": 8, "height": 4},
			Config: map[string]interface{}{
				"chart_type": "line",
				"time_range": "7d",
				"group_by":   "day",
			},
		},
	}

	for _, widget := range widgets {
		widget.LastUpdated = time.Now()
		suite.dashboardWidgets = append(suite.dashboardWidgets, widget)
	}
}

// Test event tracking and ingestion
func (suite *Journey08AnalyticsAPISuite) TestEventTrackingAndIngestion() {
	for _, event := range suite.testEvents {
		jsonData, _ := json.Marshal(event)
		req, _ := http.NewRequest("POST",
			fmt.Sprintf("%s/api/v1/analytics/events", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusCreated, resp.StatusCode)

		var eventResponse AnalyticsEvent
		json.NewDecoder(resp.Body).Decode(&eventResponse)
		suite.NotEmpty(eventResponse.ID)
	}

	// Test batch event ingestion
	batchEvents := map[string]interface{}{
		"events": suite.testEvents,
		"batch_id": "batch_001",
	}

	jsonData, _ := json.Marshal(batchEvents)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/analytics/events/batch", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var batchResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&batchResponse)
	suite.Contains(batchResponse, "processed_count")
	suite.Contains(batchResponse, "batch_id")
}

// Test user behavior analytics
func (suite *Journey08AnalyticsAPISuite) TestUserBehaviorAnalytics() {
	// Get user behavior metrics
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/users/%s/behavior", suite.baseURL, suite.user1.UserID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var behaviorMetrics UserBehaviorMetrics
	json.NewDecoder(resp.Body).Decode(&behaviorMetrics)

	suite.Equal(suite.user1.UserID, behaviorMetrics.UserID)
	suite.Greater(behaviorMetrics.EngagementScore, 0.0)
	suite.NotEmpty(behaviorMetrics.ActivityHeatmap)

	// Get user segment analytics
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/users/segments", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var segmentResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&segmentResponse)

	suite.Contains(segmentResponse, "segments")
	segments := segmentResponse["segments"].([]interface{})
	suite.Greater(len(segments), 0)
}

// Test business metrics and KPIs
func (suite *Journey08AnalyticsAPISuite) TestBusinessMetricsAndKPIs() {
	// Get business metrics for different periods
	periods := []string{"today", "week", "month", "quarter"}

	for _, period := range periods {
		req, _ := http.NewRequest("GET",
			fmt.Sprintf("%s/api/v1/analytics/business/metrics?period=%s", suite.baseURL, period), nil)
		req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

		resp, err := suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		var businessMetrics BusinessMetrics
		json.NewDecoder(resp.Body).Decode(&businessMetrics)

		suite.Equal(period, businessMetrics.Period)
		suite.GreaterOrEqual(businessMetrics.ActiveUsers, 0)
		suite.GreaterOrEqual(businessMetrics.ConversionRate, 0.0)
	}

	// Get KPI dashboard
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/kpis/dashboard", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var kpiDashboard map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&kpiDashboard)

	suite.Contains(kpiDashboard, "revenue")
	suite.Contains(kpiDashboard, "users")
	suite.Contains(kpiDashboard, "engagement")
	suite.Contains(kpiDashboard, "performance")
}

// Test performance metrics and monitoring
func (suite *Journey08AnalyticsAPISuite) TestPerformanceMetricsAndMonitoring() {
	// Get current performance metrics
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/performance/current", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var perfMetrics PerformanceMetrics
	json.NewDecoder(resp.Body).Decode(&perfMetrics)

	suite.Greater(perfMetrics.ThroughputRPS, 0.0)
	suite.GreaterOrEqual(perfMetrics.ErrorRate, 0.0)
	suite.Greater(perfMetrics.ServiceAvailability, 90.0)

	// Get historical performance data
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/performance/history?hours=24", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var historyResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&historyResponse)

	suite.Contains(historyResponse, "metrics")
	metrics := historyResponse["metrics"].([]interface{})
	suite.Greater(len(metrics), 0)

	// Test performance alerts
	alertRule := map[string]interface{}{
		"name": "High Response Time Alert",
		"metric": "response_time",
		"threshold": 1000.0,
		"operator": "greater_than",
		"duration": "5m",
		"severity": "warning",
	}

	jsonData, _ := json.Marshal(alertRule)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/analytics/performance/alerts", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)
}

// Test revenue analytics and financial metrics
func (suite *Journey08AnalyticsAPISuite) TestRevenueAnalyticsAndFinancialMetrics() {
	// Get revenue analytics for current month
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/revenue/monthly", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var revenueAnalytics RevenueAnalytics
	json.NewDecoder(resp.Body).Decode(&revenueAnalytics)

	suite.Contains(revenueAnalytics.RevenueByRegion, "Singapore")
	suite.Contains(revenueAnalytics.RevenueByRegion, "Thailand")
	suite.GreaterOrEqual(revenueAnalytics.TotalRevenue, 0.0)

	// Get revenue breakdown by product
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/revenue/products", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var productRevenue map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&productRevenue)

	suite.Contains(productRevenue, "products")
	suite.Contains(productRevenue, "total_revenue")

	// Test financial forecasting
	forecastParams := map[string]interface{}{
		"horizon_months": 6,
		"confidence_interval": 0.95,
		"include_seasonal": true,
	}

	jsonData, _ := json.Marshal(forecastParams)
	req, _ = http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/analytics/revenue/forecast", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var forecastResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&forecastResponse)

	suite.Contains(forecastResponse, "forecast")
	suite.Contains(forecastResponse, "confidence_bounds")
}

// Test content analytics and engagement metrics
func (suite *Journey08AnalyticsAPISuite) TestContentPerformanceAnalyticsAndEngagementMetrics() {
	// Create sample content for analytics
	contentData := map[string]interface{}{
		"title":       "Test Content for Analytics",
		"type":        "video",
		"description": "Content created for analytics testing",
		"tags":        []string{"test", "analytics"},
	}

	jsonData, _ := json.Marshal(contentData)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/content", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var contentResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&contentResponse)
	contentID, ok := contentResponse["id"].(string)
	suite.True(ok, "Expected id in content response: %+v", contentResponse)

	// Get content analytics
	time.Sleep(2 * time.Second) // Allow analytics processing

	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/content/%s", suite.baseURL, contentID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.user1.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var contentAnalytics ContentPerformanceAnalytics
	json.NewDecoder(resp.Body).Decode(&contentAnalytics)

	suite.Equal(contentID, contentAnalytics.ContentID)
	suite.GreaterOrEqual(contentAnalytics.Views, 0)
	suite.GreaterOrEqual(contentAnalytics.EngagementRate, 0.0)

	// Get top performing content
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/content/top?metric=engagement&period=week", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var topContentResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&topContentResponse)

	suite.Contains(topContentResponse, "content")
	suite.Contains(topContentResponse, "period")
}

// Test custom analytics queries
func (suite *Journey08AnalyticsAPISuite) TestCustomAnalyticsQueries() {
	// Create custom analytics query
	for _, query := range suite.testQueries {
		jsonData, _ := json.Marshal(query)
		req, _ := http.NewRequest("POST",
			fmt.Sprintf("%s/api/v1/analytics/queries", suite.baseURL),
			bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.httpClient.Do(req)
		suite.NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusCreated, resp.StatusCode)

		var queryResponse AnalyticsQuery
		json.NewDecoder(resp.Body).Decode(&queryResponse)
		suite.NotEmpty(queryResponse.ID)
	}

	// Execute custom query
	executeData := map[string]interface{}{
		"query_id": "test_query_001",
		"parameters": map[string]interface{}{
			"date": time.Now().Format("2006-01-02"),
		},
		"output_format": "json",
	}

	jsonData, _ := json.Marshal(executeData)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/analytics/queries/execute", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var queryResults map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&queryResults)

	suite.Contains(queryResults, "results")
	suite.Contains(queryResults, "execution_time")
}

// Test real-time dashboard functionality
func (suite *Journey08AnalyticsAPISuite) TestRealTimeDashboardFunctionality() {
	// Create dashboard
	dashboardData := map[string]interface{}{
		"name":        "Test Analytics Dashboard",
		"description": "Dashboard for testing analytics functionality",
		"widgets":     suite.dashboardWidgets,
		"layout":      "grid",
		"refresh_rate": 300,
		"filters": map[string]interface{}{
			"date_range": "7d",
			"platform":   "all",
		},
	}

	jsonData, _ := json.Marshal(dashboardData)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/analytics/dashboards", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusCreated, resp.StatusCode)

	var dashboardResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&dashboardResponse)
	dashboardID, ok := dashboardResponse["id"].(string)
	suite.True(ok, "Expected id in dashboard response: %+v", dashboardResponse)

	// Get dashboard data
	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/dashboards/%s/data", suite.baseURL, dashboardID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var dashboardData2 map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&dashboardData2)

	suite.Contains(dashboardData2, "widgets")
	suite.Contains(dashboardData2, "last_updated")

	// Update dashboard widget
	widgetUpdate := map[string]interface{}{
		"widget_id": "widget_001",
		"config": map[string]interface{}{
			"time_range": "30d",
			"metric_type": "sum",
		},
	}

	jsonData, _ = json.Marshal(widgetUpdate)
	req, _ = http.NewRequest("PUT",
		fmt.Sprintf("%s/api/v1/analytics/dashboards/%s/widgets", suite.baseURL, dashboardID),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
}

// Test analytics data export functionality
func (suite *Journey08AnalyticsAPISuite) TestAnalyticsDataExport() {
	// Export user behavior data
	exportRequest := map[string]interface{}{
		"data_type": "user_behavior",
		"format":    "csv",
		"filters": map[string]interface{}{
			"date_range": "30d",
			"user_segment": "active_users",
		},
		"fields": []string{
			"user_id", "session_count", "total_session_time",
			"page_views", "engagement_score",
		},
	}

	jsonData, _ := json.Marshal(exportRequest)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/api/v1/analytics/export", suite.baseURL),
		bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	suite.Equal(http.StatusAccepted, resp.StatusCode)

	var exportResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&exportResponse)

	suite.Contains(exportResponse, "export_id")
	suite.Contains(exportResponse, "status")

	// Check export status
	exportID, ok := exportResponse["export_id"].(string)
	suite.True(ok, "Expected export_id in export response: %+v", exportResponse)
	time.Sleep(2 * time.Second) // Allow processing time

	req, _ = http.NewRequest("GET",
		fmt.Sprintf("%s/api/v1/analytics/export/%s/status", suite.baseURL, exportID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.analyst.AccessToken)

	resp, err = suite.httpClient.Do(req)
	suite.NoError(err)
	defer resp.Body.Close()

	var statusResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&statusResponse)

	suite.Contains(statusResponse, "status")
	suite.Contains(statusResponse, "progress")
}

func (suite *Journey08AnalyticsAPISuite) cleanupTestData() {
	// Clean up test analytics data
	req, _ := http.NewRequest("DELETE",
		fmt.Sprintf("%s/api/v1/analytics/cleanup/test-data", suite.baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.admin.AccessToken)

	_, _ = suite.httpClient.Do(req)
}

func TestJourney08AnalyticsAPISuite(t *testing.T) {
	suite.Run(t, new(Journey08AnalyticsAPISuite))
}
