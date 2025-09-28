package models

import (
	"time"
)

// ServiceCompletion represents the completion status of a specific service
type ServiceCompletion struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	ServiceID        string    `json:"serviceId" gorm:"not null;uniqueIndex:idx_service_platform"`
	Platform         string    `json:"platform" gorm:"not null;uniqueIndex:idx_service_platform"` // BACKEND, WEB, IOS, ANDROID, KMP
	ServiceName      string    `json:"serviceName" gorm:"not null"`
	ServiceType      string    `json:"serviceType" gorm:"not null"` // MICROSERVICE, FRONTEND, MOBILE, SHARED
	PlaceholderCount int       `json:"placeholderCount" gorm:"not null;default:0"`
	CompletedCount   int       `json:"completedCount" gorm:"not null;default:0"`
	TestsPassing     bool      `json:"testsPassing" gorm:"not null;default:false"`
	BuildSuccessful  bool      `json:"buildSuccessful" gorm:"not null;default:false"`
	LastUpdated      time.Time `json:"lastUpdated" gorm:"autoUpdateTime"`
	CreatedAt        time.Time `json:"createdAt" gorm:"autoCreateTime"`

	// Detailed metrics
	CriticalPlaceholders int `json:"criticalPlaceholders" gorm:"not null;default:0"`
	HighPlaceholders     int `json:"highPlaceholders" gorm:"not null;default:0"`
	MediumPlaceholders   int `json:"mediumPlaceholders" gorm:"not null;default:0"`
	LowPlaceholders      int `json:"lowPlaceholders" gorm:"not null;default:0"`

	// Progress tracking
	CompletionPercentage    float64  `json:"completionPercentage" gorm:"not null;default:0"`
	EstimatedHoursRemaining *float64 `json:"estimatedHoursRemaining"`
	EstimatedCompletionDate *time.Time `json:"estimatedCompletionDate"`

	// Quality metrics
	CodeQualityScore    *float64 `json:"codeQualityScore"`    // 0-100 score
	TestCoverage        *float64 `json:"testCoverage"`        // 0-100 percentage
	SecurityScore       *float64 `json:"securityScore"`       // 0-100 score
	PerformanceScore    *float64 `json:"performanceScore"`    // 0-100 score
	DocumentationScore  *float64 `json:"documentationScore"`  // 0-100 score

	// Performance metrics
	PerformanceMetrics *PerformanceMetrics `json:"performanceMetrics" gorm:"embedded"`

	// Build and deployment info
	BuildInfo     *BuildInfo     `json:"buildInfo" gorm:"embedded"`
	DeploymentInfo *DeploymentInfo `json:"deploymentInfo" gorm:"embedded"`

	// Testing information
	TestInfo *TestInfo `json:"testInfo" gorm:"embedded"`

	// Dependencies and relationships
	Dependencies     []string `json:"dependencies" gorm:"type:json"`     // Services this service depends on
	Dependents       []string `json:"dependents" gorm:"type:json"`       // Services that depend on this service
	BlockingServices []string `json:"blockingServices" gorm:"type:json"` // Services blocking this service's completion

	// Regional optimization (Southeast Asian markets)
	RegionalOptimization *RegionalOptimization `json:"regionalOptimization" gorm:"embedded"`

	// Assignee and ownership
	TeamOwner       *string `json:"teamOwner"`       // Team responsible for this service
	TechnicalLead   *string `json:"technicalLead"`   // Technical lead for this service
	ProductOwner    *string `json:"productOwner"`    // Product owner for this service

	// Status and health
	HealthStatus        string    `json:"healthStatus" gorm:"not null;default:'UNKNOWN'"` // HEALTHY, DEGRADED, UNHEALTHY, UNKNOWN
	LastHealthCheck     *time.Time `json:"lastHealthCheck"`
	ServiceVersion      *string   `json:"serviceVersion"`      // Current deployed version
	LastDeployment      *time.Time `json:"lastDeployment"`      // Last deployment timestamp
	UpstreamStatus      string    `json:"upstreamStatus" gorm:"default:'UNKNOWN'"` // Status of upstream dependencies

	// Quality gates
	QualityGates *QualityGates `json:"qualityGates" gorm:"embedded"`

	// Trend analysis
	TrendData *TrendData `json:"trendData" gorm:"embedded"`

	// Risk assessment
	RiskAssessment *RiskAssessment `json:"riskAssessment" gorm:"embedded"`
}

// PerformanceMetrics contains performance-related metrics for a service
type PerformanceMetrics struct {
	AverageResponseTime  *float64 `json:"averageResponseTime"`  // Average response time in ms
	P95ResponseTime      *float64 `json:"p95ResponseTime"`      // 95th percentile response time in ms
	ThroughputPerSecond  *float64 `json:"throughputPerSecond"`  // Requests per second
	ErrorRate            *float64 `json:"errorRate"`            // Error rate percentage
	MemoryUsageMB        *float64 `json:"memoryUsageMb"`        // Memory usage in MB
	CPUUtilization       *float64 `json:"cpuUtilization"`       // CPU utilization percentage
	DiskUsageGB          *float64 `json:"diskUsageGb"`          // Disk usage in GB
	NetworkBandwidthMbps *float64 `json:"networkBandwidthMbps"` // Network bandwidth in Mbps

	// Mobile-specific metrics
	AppLaunchTime        *float64 `json:"appLaunchTime"`        // App launch time in seconds
	FrameRate            *float64 `json:"frameRate"`            // Frames per second (mobile)
	BatteryUsage         *float64 `json:"batteryUsage"`         // Battery usage percentage per hour
	CrashRate            *float64 `json:"crashRate"`            // Crash rate percentage

	// Web-specific metrics
	FirstContentfulPaint *float64 `json:"firstContentfulPaint"` // FCP in ms
	LargestContentfulPaint *float64 `json:"largestContentfulPaint"` // LCP in ms
	CumulativeLayoutShift *float64 `json:"cumulativeLayoutShift"` // CLS score
	FirstInputDelay      *float64 `json:"firstInputDelay"`      // FID in ms
	BundleSize           *float64 `json:"bundleSize"`           // Bundle size in KB
}

// BuildInfo contains build-related information
type BuildInfo struct {
	LastBuildTime      *time.Time `json:"lastBuildTime"`
	BuildDuration      *float64   `json:"buildDuration"`      // Build duration in seconds
	BuildSuccessRate   *float64   `json:"buildSuccessRate"`   // Success rate over last 10 builds
	ArtifactSize       *float64   `json:"artifactSize"`       // Build artifact size in MB
	BuildNumber        *string    `json:"buildNumber"`        // Latest build number
	BuildBranch        *string    `json:"buildBranch"`        // Branch used for build
	CompilerWarnings   *int       `json:"compilerWarnings"`   // Number of compiler warnings
	LintingIssues      *int       `json:"lintingIssues"`      // Number of linting issues
}

// DeploymentInfo contains deployment-related information
type DeploymentInfo struct {
	LastDeployment        *time.Time `json:"lastDeployment"`
	DeploymentEnvironment *string    `json:"deploymentEnvironment"` // DEVELOPMENT, STAGING, PRODUCTION
	DeploymentStatus      *string    `json:"deploymentStatus"`      // SUCCESS, FAILED, IN_PROGRESS
	RollbackCapable       *bool      `json:"rollbackCapable"`       // Whether rollback is possible
	DeploymentDuration    *float64   `json:"deploymentDuration"`    // Deployment duration in seconds
	AutoDeployEnabled     *bool      `json:"autoDeployEnabled"`     // Whether auto-deployment is enabled
	CanaryDeployment      *bool      `json:"canaryDeployment"`      // Whether using canary deployment
}

// TestInfo contains testing-related information
type TestInfo struct {
	UnitTestsPassing      *bool      `json:"unitTestsPassing"`
	IntegrationTestsPassing *bool    `json:"integrationTestsPassing"`
	E2ETestsPassing       *bool      `json:"e2eTestsPassing"`
	UnitTestCoverage      *float64   `json:"unitTestCoverage"`      // Unit test coverage percentage
	IntegrationTestCoverage *float64 `json:"integrationTestCoverage"` // Integration test coverage percentage
	TotalTests            *int       `json:"totalTests"`            // Total number of tests
	PassingTests          *int       `json:"passingTests"`          // Number of passing tests
	FailingTests          *int       `json:"failingTests"`          // Number of failing tests
	SkippedTests          *int       `json:"skippedTests"`          // Number of skipped tests
	TestExecutionTime     *float64   `json:"testExecutionTime"`     // Test execution time in seconds
	FlakynessScore        *float64   `json:"flakynessScore"`        // Test flakiness score (0-100)
}

// RegionalOptimization contains regional optimization information for SEA markets
type RegionalOptimization struct {
	ThailandOptimized   bool     `json:"thailandOptimized"`   // TH market optimization
	SingaporeOptimized  bool     `json:"singaporeOptimized"`  // SG market optimization
	MalaysiaOptimized   bool     `json:"malaysiaOptimized"`   // MY market optimization
	IndonesiaOptimized  bool     `json:"indonesiaOptimized"`  // ID market optimization
	PhilippinesOptimized bool    `json:"philippinesOptimized"` // PH market optimization
	VietnamOptimized    bool     `json:"vietnamOptimized"`    // VN market optimization

	LocalizationScore   *float64 `json:"localizationScore"`   // 0-100 localization completeness
	RegionalPerformance map[string]float64 `json:"regionalPerformance" gorm:"type:json"` // Performance by region
	CulturalAdaptation  *float64 `json:"culturalAdaptation"`  // 0-100 cultural adaptation score
	ComplianceScore     *float64 `json:"complianceScore"`     // 0-100 regional compliance score
}

// QualityGates contains quality gate status
type QualityGates struct {
	CodeQualityPassed    bool `json:"codeQualityPassed"`
	SecurityScanPassed   bool `json:"securityScanPassed"`
	PerformancePassed    bool `json:"performancePassed"`
	AccessibilityPassed  bool `json:"accessibilityPassed"`
	DocumentationPassed  bool `json:"documentationPassed"`
	TestCoveragePassed   bool `json:"testCoveragePassed"`
	LicenseScanPassed    bool `json:"licenseScanPassed"`
	VulnerabilityScanPassed bool `json:"vulnerabilityScanPassed"`
}

// TrendData contains trend analysis information
type TrendData struct {
	CompletionTrend      *string  `json:"completionTrend"`      // IMPROVING, DECLINING, STABLE
	QualityTrend         *string  `json:"qualityTrend"`         // IMPROVING, DECLINING, STABLE
	PerformanceTrend     *string  `json:"performanceTrend"`     // IMPROVING, DECLINING, STABLE
	TrendConfidence      *float64 `json:"trendConfidence"`      // 0-1 confidence in trend analysis
	VelocityScore        *float64 `json:"velocityScore"`        // Development velocity score
	StabilityScore       *float64 `json:"stabilityScore"`       // Service stability score
}

// RiskAssessment contains risk assessment information
type RiskAssessment struct {
	OverallRiskScore     *float64 `json:"overallRiskScore"`     // 0-100 overall risk score
	TechnicalDebtRisk    *float64 `json:"technicalDebtRisk"`    // 0-100 technical debt risk
	SecurityRisk         *float64 `json:"securityRisk"`         // 0-100 security risk
	PerformanceRisk      *float64 `json:"performanceRisk"`      // 0-100 performance risk
	DependencyRisk       *float64 `json:"dependencyRisk"`       // 0-100 dependency risk
	MaintenanceRisk      *float64 `json:"maintenanceRisk"`      // 0-100 maintenance risk
	ComplianceRisk       *float64 `json:"complianceRisk"`       // 0-100 compliance risk

	RiskFactors          []string `json:"riskFactors" gorm:"type:json"`          // List of identified risk factors
	MitigationStrategies []string `json:"mitigationStrategies" gorm:"type:json"` // Risk mitigation strategies
	RiskOwner            *string  `json:"riskOwner"`            // Person responsible for risk management
}

// ServiceType defines the types of services
type ServiceType string

const (
	ServiceTypeMicroservice ServiceType = "MICROSERVICE"
	ServiceTypeFrontend     ServiceType = "FRONTEND"
	ServiceTypeMobile       ServiceType = "MOBILE"
	ServiceTypeShared       ServiceType = "SHARED"
	ServiceTypeLibrary      ServiceType = "LIBRARY"
	ServiceTypeTools        ServiceType = "TOOLS"
)

// HealthStatus defines the health status of a service
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "HEALTHY"
	HealthStatusDegraded  HealthStatus = "DEGRADED"
	HealthStatusUnhealthy HealthStatus = "UNHEALTHY"
	HealthStatusUnknown   HealthStatus = "UNKNOWN"
)

// GetCompletionPercentage calculates the completion percentage
func (sc *ServiceCompletion) GetCompletionPercentage() float64 {
	if sc.PlaceholderCount == 0 {
		return 100.0
	}
	return (float64(sc.CompletedCount) / float64(sc.PlaceholderCount)) * 100.0
}

// GetRemainingPlaceholders returns the number of remaining placeholders
func (sc *ServiceCompletion) GetRemainingPlaceholders() int {
	return sc.PlaceholderCount - sc.CompletedCount
}

// IsCompleted returns true if all placeholders are completed
func (sc *ServiceCompletion) IsCompleted() bool {
	return sc.GetRemainingPlaceholders() == 0
}

// GetPriorityScore calculates a priority score based on critical and high priority placeholders
func (sc *ServiceCompletion) GetPriorityScore() int {
	return (sc.CriticalPlaceholders * 4) + (sc.HighPlaceholders * 3) + (sc.MediumPlaceholders * 2) + sc.LowPlaceholders
}

// GetHealthScore calculates an overall health score (0-100) based on multiple factors
func (sc *ServiceCompletion) GetHealthScore() float64 {
	score := 0.0
	factors := 0

	// Completion percentage (30% weight)
	score += sc.CompletionPercentage * 0.3
	factors++

	// Build and test status (25% weight)
	if sc.BuildSuccessful {
		score += 25.0
	}
	if sc.TestsPassing {
		score += 25.0
	}
	factors++

	// Quality metrics (20% weight)
	if sc.CodeQualityScore != nil {
		score += (*sc.CodeQualityScore) * 0.2
		factors++
	}

	// Performance metrics (15% weight)
	if sc.PerformanceScore != nil {
		score += (*sc.PerformanceScore) * 0.15
		factors++
	}

	// Test coverage (10% weight)
	if sc.TestCoverage != nil {
		score += (*sc.TestCoverage) * 0.1
		factors++
	}

	return score
}

// GetRegionalOptimizationScore calculates the regional optimization score for SEA markets
func (sc *ServiceCompletion) GetRegionalOptimizationScore() float64 {
	if sc.RegionalOptimization == nil {
		return 0.0
	}

	// Count optimized markets
	optimizedCount := 0
	totalMarkets := 6 // TH, SG, MY, ID, PH, VN

	if sc.RegionalOptimization.ThailandOptimized {
		optimizedCount++
	}
	if sc.RegionalOptimization.SingaporeOptimized {
		optimizedCount++
	}
	if sc.RegionalOptimization.MalaysiaOptimized {
		optimizedCount++
	}
	if sc.RegionalOptimization.IndonesiaOptimized {
		optimizedCount++
	}
	if sc.RegionalOptimization.PhilippinesOptimized {
		optimizedCount++
	}
	if sc.RegionalOptimization.VietnamOptimized {
		optimizedCount++
	}

	marketScore := (float64(optimizedCount) / float64(totalMarkets)) * 100.0

	// Include localization score if available
	if sc.RegionalOptimization.LocalizationScore != nil {
		return (marketScore + *sc.RegionalOptimization.LocalizationScore) / 2.0
	}

	return marketScore
}

// IsBlockedByDependencies returns true if the service is blocked by dependencies
func (sc *ServiceCompletion) IsBlockedByDependencies() bool {
	return len(sc.BlockingServices) > 0
}

// HasCriticalIssues returns true if the service has critical issues that need immediate attention
func (sc *ServiceCompletion) HasCriticalIssues() bool {
	return sc.CriticalPlaceholders > 0 || !sc.BuildSuccessful || sc.HealthStatus == string(HealthStatusUnhealthy)
}

// GetEstimatedCompletionDays returns estimated days to completion based on remaining work
func (sc *ServiceCompletion) GetEstimatedCompletionDays() *float64 {
	if sc.EstimatedHoursRemaining == nil || *sc.EstimatedHoursRemaining <= 0 {
		return nil
	}

	// Assuming 8 hours per working day
	days := *sc.EstimatedHoursRemaining / 8.0
	return &days
}

// GetEfficiency calculates development efficiency based on completion rate and time
func (sc *ServiceCompletion) GetEfficiency() float64 {
	// Simple efficiency calculation: completion percentage / days since creation
	daysSinceCreation := time.Since(sc.CreatedAt).Hours() / 24.0
	if daysSinceCreation <= 0 {
		return 0.0
	}

	return sc.CompletionPercentage / daysSinceCreation
}