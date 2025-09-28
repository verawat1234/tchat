package services

import (
	"context"
	"fmt"
	"time"

	"github.com/tchat/backend/audit/models"
	"gorm.io/gorm"
)

type ServiceCompletionService struct {
	db *gorm.DB
}

type ServiceCompletionFilters struct {
	Platform       string
	Status         string
	IncludeDetails bool
	Limit          int
	Offset         int
}

type ServiceHealthOverview struct {
	TotalServices    int64                    `json:"totalServices"`
	HealthyServices  int64                    `json:"healthyServices"`
	DegradedServices int64                    `json:"degradedServices"`
	UnhealthyServices int64                   `json:"unhealthyServices"`
	UnknownServices  int64                    `json:"unknownServices"`
	ByPlatform      map[string]ServiceHealth `json:"byPlatform"`
	ByRegion        map[string]ServiceHealth `json:"byRegion"`
	LastUpdated     time.Time                `json:"lastUpdated"`
}

type ServiceHealth struct {
	Total     int64   `json:"total"`
	Healthy   int64   `json:"healthy"`
	Degraded  int64   `json:"degraded"`
	Unhealthy int64   `json:"unhealthy"`
	Unknown   int64   `json:"unknown"`
	HealthPct float64 `json:"healthPct"`
}

type ServiceMetrics struct {
	AverageCompletion        float64            `json:"averageCompletion"`
	AverageQualityScore      float64            `json:"averageQualityScore"`
	AveragePerformanceScore  float64            `json:"averagePerformanceScore"`
	AverageSecurityScore     float64            `json:"averageSecurityScore"`
	TotalPlaceholders        int64              `json:"totalPlaceholders"`
	CompletedPlaceholders    int64              `json:"completedPlaceholders"`
	CriticalPlaceholders     int64              `json:"criticalPlaceholders"`
	TestCoverageAverage      float64            `json:"testCoverageAverage"`
	ServicesWithFailingTests int64              `json:"servicesWithFailingTests"`
	ServicesWithFailingBuilds int64             `json:"servicesWithFailingBuilds"`
	RegionalOptimization     map[string]float64 `json:"regionalOptimization"`
	LastCalculated          time.Time          `json:"lastCalculated"`
}

type ServiceRefreshResult struct {
	ServiceID            string    `json:"serviceId"`
	Platform             string    `json:"platform"`
	RefreshStarted       time.Time `json:"refreshStarted"`
	RefreshCompleted     time.Time `json:"refreshCompleted"`
	PlaceholdersScanned  int       `json:"placeholdersScanned"`
	PlaceholdersUpdated  int       `json:"placeholdersUpdated"`
	TestsExecuted        bool      `json:"testsExecuted"`
	BuildStatus          bool      `json:"buildStatus"`
	MetricsUpdated       bool      `json:"metricsUpdated"`
	PreviousCompletion   float64   `json:"previousCompletion"`
	CurrentCompletion    float64   `json:"currentCompletion"`
	CompletionChange     float64   `json:"completionChange"`
}

type RegionalOptimizationData struct {
	Region                   string                 `json:"region"`
	Services                []models.ServiceCompletion `json:"services"`
	TotalServices           int                    `json:"totalServices"`
	OptimizedServices       int                    `json:"optimizedServices"`
	OptimizationPercentage  float64                `json:"optimizationPercentage"`
	AverageLocalizationScore float64               `json:"averageLocalizationScore"`
	CulturalAdaptationScore float64                `json:"culturalAdaptationScore"`
	ComplianceScore         float64                `json:"complianceScore"`
	PerformanceByRegion     map[string]float64     `json:"performanceByRegion"`
	RecommendedActions      []string               `json:"recommendedActions"`
}

type DependencyGraph struct {
	Services     []ServiceNode    `json:"services"`
	Dependencies []DependencyEdge `json:"dependencies"`
}

type ServiceNode struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Platform     string  `json:"platform"`
	Status       string  `json:"status"`
	Completion   float64 `json:"completion"`
	HealthStatus string  `json:"healthStatus"`
}

type DependencyEdge struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Type         string `json:"type"` // depends_on, blocks, provides
	Critical     bool   `json:"critical"`
}

type CompletionTrend struct {
	Timestamp        time.Time `json:"timestamp"`
	CompletionPct    float64   `json:"completionPct"`
	PlaceholderCount int       `json:"placeholderCount"`
	QualityScore     float64   `json:"qualityScore"`
	TestCoverage     float64   `json:"testCoverage"`
}

func NewServiceCompletionService(db *gorm.DB) *ServiceCompletionService {
	return &ServiceCompletionService{db: db}
}

func (s *ServiceCompletionService) GetServiceCompletion(ctx context.Context, serviceID, platform string, includeDetails bool) (*models.ServiceCompletion, error) {
	var completion models.ServiceCompletion

	query := s.db.WithContext(ctx).Where("service_id = ?", serviceID)
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	if err := query.First(&completion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("service completion not found")
		}
		return nil, fmt.Errorf("failed to retrieve service completion: %w", err)
	}

	// Calculate dynamic fields
	completion.CompletionPercentage = completion.GetCompletionPercentage()

	return &completion, nil
}

func (s *ServiceCompletionService) GetAllServiceCompletions(ctx context.Context, filters *ServiceCompletionFilters) ([]models.ServiceCompletion, int64, error) {
	var completions []models.ServiceCompletion
	var total int64

	query := s.db.WithContext(ctx).Model(&models.ServiceCompletion{})

	// Apply filters
	if filters.Platform != "" {
		query = query.Where("platform = ?", filters.Platform)
	}
	if filters.Status != "" {
		query = query.Where("health_status = ?", filters.Status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count service completions: %w", err)
	}

	// Apply pagination and ordering
	query = query.Order("completion_percentage DESC, service_name ASC").
		Limit(filters.Limit).Offset(filters.Offset)

	if err := query.Find(&completions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve service completions: %w", err)
	}

	// Calculate dynamic fields for each completion
	for i := range completions {
		completions[i].CompletionPercentage = completions[i].GetCompletionPercentage()
	}

	return completions, total, nil
}

func (s *ServiceCompletionService) UpdateServiceCompletion(ctx context.Context, serviceID string, updates map[string]interface{}) (*models.ServiceCompletion, error) {
	var completion models.ServiceCompletion

	// First, get the existing service completion
	if err := s.db.WithContext(ctx).Where("service_id = ?", serviceID).First(&completion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("service completion not found")
		}
		return nil, fmt.Errorf("failed to retrieve service completion: %w", err)
	}

	// Add timestamp
	updates["last_updated"] = time.Now()

	// Recalculate completion percentage if placeholder counts changed
	if placeholderCount, exists := updates["placeholderCount"]; exists {
		if completedCount, hasCompleted := updates["completedCount"]; hasCompleted {
			if pc, ok := placeholderCount.(int); ok && pc > 0 {
				if cc, ok := completedCount.(int); ok {
					updates["completion_percentage"] = float64(cc) / float64(pc) * 100.0
				}
			}
		}
	}

	// Update the completion
	if err := s.db.WithContext(ctx).Model(&completion).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update service completion: %w", err)
	}

	// Retrieve updated completion
	if err := s.db.WithContext(ctx).Where("service_id = ?", serviceID).First(&completion).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve updated service completion: %w", err)
	}

	// Calculate dynamic fields
	completion.CompletionPercentage = completion.GetCompletionPercentage()

	return &completion, nil
}

func (s *ServiceCompletionService) GetServiceHealth(ctx context.Context, platform, region string) (*ServiceHealthOverview, error) {
	var overview ServiceHealthOverview

	// Base query
	query := s.db.WithContext(ctx).Model(&models.ServiceCompletion{})
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	// Get total count
	if err := query.Count(&overview.TotalServices).Error; err != nil {
		return nil, fmt.Errorf("failed to get total services count: %w", err)
	}

	// Get health status counts
	var healthCounts []struct {
		HealthStatus string
		Count        int64
	}

	if err := query.Session(&gorm.Session{}).Select("health_status, COUNT(*) as count").Group("health_status").Scan(&healthCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get health status counts: %w", err)
	}

	for _, count := range healthCounts {
		switch count.HealthStatus {
		case string(models.HealthStatusHealthy):
			overview.HealthyServices = count.Count
		case string(models.HealthStatusDegraded):
			overview.DegradedServices = count.Count
		case string(models.HealthStatusUnhealthy):
			overview.UnhealthyServices = count.Count
		case string(models.HealthStatusUnknown):
			overview.UnknownServices = count.Count
		}
	}

	// Get health by platform
	overview.ByPlatform = make(map[string]ServiceHealth)
	var platformHealth []struct {
		Platform     string
		HealthStatus string
		Count        int64
	}

	platformQuery := s.db.WithContext(ctx).Model(&models.ServiceCompletion{}).
		Select("platform, health_status, COUNT(*) as count").
		Group("platform, health_status")

	if err := platformQuery.Scan(&platformHealth).Error; err != nil {
		return nil, fmt.Errorf("failed to get platform health: %w", err)
	}

	for _, ph := range platformHealth {
		if _, exists := overview.ByPlatform[ph.Platform]; !exists {
			overview.ByPlatform[ph.Platform] = ServiceHealth{}
		}
		health := overview.ByPlatform[ph.Platform]
		health.Total += ph.Count

		switch ph.HealthStatus {
		case string(models.HealthStatusHealthy):
			health.Healthy = ph.Count
		case string(models.HealthStatusDegraded):
			health.Degraded = ph.Count
		case string(models.HealthStatusUnhealthy):
			health.Unhealthy = ph.Count
		case string(models.HealthStatusUnknown):
			health.Unknown = ph.Count
		}

		if health.Total > 0 {
			health.HealthPct = float64(health.Healthy) / float64(health.Total) * 100.0
		}
		overview.ByPlatform[ph.Platform] = health
	}

	overview.LastUpdated = time.Now()
	return &overview, nil
}

func (s *ServiceCompletionService) GetServiceMetrics(ctx context.Context, platform, timeRange string) (*ServiceMetrics, error) {
	var metrics ServiceMetrics

	query := s.db.WithContext(ctx).Model(&models.ServiceCompletion{})
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	// Calculate average completion percentage
	var avgCompletion struct {
		AvgCompletion float64
	}
	if err := query.Session(&gorm.Session{}).Select("AVG(completion_percentage) as avg_completion").Scan(&avgCompletion).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average completion: %w", err)
	}
	metrics.AverageCompletion = avgCompletion.AvgCompletion

	// Calculate other averages
	var avgScores struct {
		AvgQuality     float64
		AvgPerformance float64
		AvgSecurity    float64
		AvgTestCoverage float64
	}

	avgQuery := query.Session(&gorm.Session{}).Select(`
		AVG(CASE WHEN code_quality_score IS NOT NULL THEN code_quality_score END) as avg_quality,
		AVG(CASE WHEN performance_score IS NOT NULL THEN performance_score END) as avg_performance,
		AVG(CASE WHEN security_score IS NOT NULL THEN security_score END) as avg_security,
		AVG(CASE WHEN test_coverage IS NOT NULL THEN test_coverage END) as avg_test_coverage
	`)

	if err := avgQuery.Scan(&avgScores).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average scores: %w", err)
	}

	metrics.AverageQualityScore = avgScores.AvgQuality
	metrics.AveragePerformanceScore = avgScores.AvgPerformance
	metrics.AverageSecurityScore = avgScores.AvgSecurity
	metrics.TestCoverageAverage = avgScores.AvgTestCoverage

	// Calculate placeholder totals
	var placeholderTotals struct {
		TotalPlaceholders     int64
		CompletedPlaceholders int64
		CriticalPlaceholders  int64
	}

	placeholderQuery := query.Session(&gorm.Session{}).Select(`
		SUM(placeholder_count) as total_placeholders,
		SUM(completed_count) as completed_placeholders,
		SUM(critical_placeholders) as critical_placeholders
	`)

	if err := placeholderQuery.Scan(&placeholderTotals).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate placeholder totals: %w", err)
	}

	metrics.TotalPlaceholders = placeholderTotals.TotalPlaceholders
	metrics.CompletedPlaceholders = placeholderTotals.CompletedPlaceholders
	metrics.CriticalPlaceholders = placeholderTotals.CriticalPlaceholders

	// Count services with failing tests and builds
	var failingCounts struct {
		FailingTests  int64
		FailingBuilds int64
	}

	failingQuery := query.Session(&gorm.Session{}).Select(`
		COUNT(CASE WHEN tests_passing = false THEN 1 END) as failing_tests,
		COUNT(CASE WHEN build_successful = false THEN 1 END) as failing_builds
	`)

	if err := failingQuery.Scan(&failingCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count failing services: %w", err)
	}

	metrics.ServicesWithFailingTests = failingCounts.FailingTests
	metrics.ServicesWithFailingBuilds = failingCounts.FailingBuilds

	// Calculate regional optimization scores (Southeast Asian markets)
	metrics.RegionalOptimization = make(map[string]float64)
	regions := []string{"TH", "SG", "MY", "ID", "PH", "VN"}

	for _, region := range regions {
		var avgOptimization struct {
			AvgScore float64
		}

		// This would need to be expanded based on actual regional optimization data
		// For now, we'll use a placeholder calculation
		regionQuery := query.Session(&gorm.Session{}).
			Joins("LEFT JOIN regional_optimization ON service_completions.id = regional_optimization.service_completion_id").
			Where("regional_optimization.region = ?", region).
			Select("AVG(regional_optimization.localization_score) as avg_score")

		if err := regionQuery.Scan(&avgOptimization).Error; err != nil {
			// If table doesn't exist yet, default to 0
			metrics.RegionalOptimization[region] = 0.0
		} else {
			metrics.RegionalOptimization[region] = avgOptimization.AvgScore
		}
	}

	metrics.LastCalculated = time.Now()
	return &metrics, nil
}

func (s *ServiceCompletionService) TriggerServiceRefresh(ctx context.Context, serviceID, platform string, fullRefresh, updateMetrics bool) (*ServiceRefreshResult, error) {
	result := &ServiceRefreshResult{
		ServiceID:       serviceID,
		Platform:        platform,
		RefreshStarted:  time.Now(),
		TestsExecuted:   updateMetrics,
		MetricsUpdated:  updateMetrics,
	}

	// Get current completion for comparison
	var currentCompletion models.ServiceCompletion
	if err := s.db.WithContext(ctx).Where("service_id = ? AND platform = ?", serviceID, platform).First(&currentCompletion).Error; err == nil {
		result.PreviousCompletion = currentCompletion.GetCompletionPercentage()
	}

	// Simulate placeholder scanning and updates
	// In a real implementation, this would scan the actual codebase
	result.PlaceholdersScanned = 50 // Mock value
	result.PlaceholdersUpdated = 5  // Mock value

	if fullRefresh {
		// Simulate full refresh operations
		result.BuildStatus = true // Mock successful build
	}

	// Update service completion record
	updates := map[string]interface{}{
		"last_updated": time.Now(),
	}

	if updateMetrics {
		// Simulate metrics updates
		updates["last_health_check"] = time.Now()
	}

	if err := s.db.WithContext(ctx).Model(&currentCompletion).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update service completion: %w", err)
	}

	// Get updated completion for comparison
	if err := s.db.WithContext(ctx).Where("service_id = ? AND platform = ?", serviceID, platform).First(&currentCompletion).Error; err == nil {
		result.CurrentCompletion = currentCompletion.GetCompletionPercentage()
		result.CompletionChange = result.CurrentCompletion - result.PreviousCompletion
	}

	result.RefreshCompleted = time.Now()
	return result, nil
}

func (s *ServiceCompletionService) GetRegionalOptimization(ctx context.Context, region, serviceType string) (*RegionalOptimizationData, error) {
	data := &RegionalOptimizationData{
		Region: region,
		PerformanceByRegion: make(map[string]float64),
	}

	query := s.db.WithContext(ctx).Model(&models.ServiceCompletion{})
	if serviceType != "" {
		query = query.Where("service_type = ?", serviceType)
	}

	// Get all services for the region
	if err := query.Find(&data.Services).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve services: %w", err)
	}

	data.TotalServices = len(data.Services)

	// Calculate optimization metrics
	var totalLocalization, totalCultural, totalCompliance float64
	optimizedCount := 0

	for _, service := range data.Services {
		if service.RegionalOptimization != nil {
			score := service.GetRegionalOptimizationScore()
			if score > 70.0 { // Consider 70%+ as optimized
				optimizedCount++
			}

			if service.RegionalOptimization.LocalizationScore != nil {
				totalLocalization += *service.RegionalOptimization.LocalizationScore
			}
			if service.RegionalOptimization.CulturalAdaptation != nil {
				totalCultural += *service.RegionalOptimization.CulturalAdaptation
			}
			if service.RegionalOptimization.ComplianceScore != nil {
				totalCompliance += *service.RegionalOptimization.ComplianceScore
			}
		}
	}

	data.OptimizedServices = optimizedCount
	if data.TotalServices > 0 {
		data.OptimizationPercentage = float64(optimizedCount) / float64(data.TotalServices) * 100.0
		data.AverageLocalizationScore = totalLocalization / float64(data.TotalServices)
		data.CulturalAdaptationScore = totalCultural / float64(data.TotalServices)
		data.ComplianceScore = totalCompliance / float64(data.TotalServices)
	}

	// Generate recommendations based on optimization status
	if data.OptimizationPercentage < 80.0 {
		data.RecommendedActions = append(data.RecommendedActions, "Increase localization coverage for "+region+" market")
	}
	if data.AverageLocalizationScore < 70.0 {
		data.RecommendedActions = append(data.RecommendedActions, "Improve translation quality and cultural adaptation")
	}
	if data.ComplianceScore < 85.0 {
		data.RecommendedActions = append(data.RecommendedActions, "Review regulatory compliance for "+region+" region")
	}

	return data, nil
}

func (s *ServiceCompletionService) GetDependencyGraph(ctx context.Context, serviceID string, depth int) (*DependencyGraph, error) {
	graph := &DependencyGraph{
		Services:     []ServiceNode{},
		Dependencies: []DependencyEdge{},
	}

	// Get all services for building the graph
	var services []models.ServiceCompletion
	if err := s.db.WithContext(ctx).Find(&services).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve services: %w", err)
	}

	// Build service nodes
	for _, service := range services {
		node := ServiceNode{
			ID:           service.ServiceID,
			Name:         service.ServiceName,
			Platform:     service.Platform,
			Status:       service.HealthStatus,
			Completion:   service.GetCompletionPercentage(),
			HealthStatus: service.HealthStatus,
		}
		graph.Services = append(graph.Services, node)

		// Build dependency edges
		for _, dep := range service.Dependencies {
			edge := DependencyEdge{
				From:     service.ServiceID,
				To:       dep,
				Type:     "depends_on",
				Critical: false, // This could be determined by priority or impact
			}
			graph.Dependencies = append(graph.Dependencies, edge)
		}

		for _, dependent := range service.Dependents {
			edge := DependencyEdge{
				From:     dependent,
				To:       service.ServiceID,
				Type:     "depends_on",
				Critical: false,
			}
			graph.Dependencies = append(graph.Dependencies, edge)
		}

		for _, blocking := range service.BlockingServices {
			edge := DependencyEdge{
				From:     blocking,
				To:       service.ServiceID,
				Type:     "blocks",
				Critical: true,
			}
			graph.Dependencies = append(graph.Dependencies, edge)
		}
	}

	return graph, nil
}

func (s *ServiceCompletionService) GetCompletionTrends(ctx context.Context, serviceID, timeRange, granularity string) ([]CompletionTrend, error) {
	// This would typically query a time-series table or audit log
	// For now, we'll return mock trend data
	trends := []CompletionTrend{}

	// Generate mock trend data based on current completion
	var currentCompletion models.ServiceCompletion
	if serviceID != "" {
		if err := s.db.WithContext(ctx).Where("service_id = ?", serviceID).First(&currentCompletion).Error; err != nil {
			return nil, fmt.Errorf("service not found: %w", err)
		}
	}

	// Generate trends for the past period based on timeRange
	days := 7
	if timeRange == "30d" {
		days = 30
	} else if timeRange == "90d" {
		days = 90
	}

	for i := days; i >= 0; i-- {
		trend := CompletionTrend{
			Timestamp:        time.Now().AddDate(0, 0, -i),
			CompletionPct:    float64(70 + i*2), // Mock progression
			PlaceholderCount: 100 - i*2,         // Mock decreasing placeholders
			QualityScore:     75.0 + float64(i)*0.5,
			TestCoverage:     80.0 + float64(i)*0.3,
		}
		trends = append(trends, trend)
	}

	return trends, nil
}