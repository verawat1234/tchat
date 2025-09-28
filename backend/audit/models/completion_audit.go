package models

import (
	"time"
)

// CompletionAudit represents a comprehensive audit of placeholder completion status
type CompletionAudit struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	AuditName       string    `json:"auditName" gorm:"not null"`
	Description     string    `json:"description" gorm:"type:text"`
	Status          string    `json:"status" gorm:"not null;default:'RUNNING'"` // RUNNING, COMPLETED, FAILED, CANCELLED
	OverallStatus   string    `json:"overallStatus" gorm:"not null;default:'UNKNOWN'"` // PASS, FAIL, WARNING, UNKNOWN
	Platforms       []string  `json:"platforms" gorm:"type:json"` // Platforms included in this audit
	RunTests        bool      `json:"runTests" gorm:"default:true"`
	CheckPerformance bool     `json:"checkPerformance" gorm:"default:true"`
	CreatedAt       time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	StartedAt       *time.Time `json:"startedAt"`
	CompletedAt     *time.Time `json:"completedAt"`

	// Execution metadata
	ExecutorID      string  `json:"executorId" gorm:"not null"`        // Who triggered the audit
	ExecutionMode   string  `json:"executionMode" gorm:"not null"`     // MANUAL, SCHEDULED, TRIGGERED
	TriggerEvent    *string `json:"triggerEvent"`                      // What triggered the audit
	ExecutionTimeMS *int64  `json:"executionTimeMs"`                   // Total execution time in milliseconds

	// Results aggregation
	BuildResults         BuildResults         `json:"buildResults" gorm:"embedded"`
	TestResults          TestResults          `json:"testResults" gorm:"embedded"`
	PerformanceResults   PerformanceResults   `json:"performanceResults" gorm:"embedded"`
	PlaceholderResults   PlaceholderResults   `json:"placeholderResults" gorm:"embedded"`
	QualityGateResults   QualityGateResults   `json:"qualityGateResults" gorm:"embedded"`

	// Violation tracking
	Violations []AuditViolation `json:"violations" gorm:"type:json"`

	// Regional optimization tracking (Southeast Asian markets)
	RegionalResults map[string]RegionalAuditResult `json:"regionalResults" gorm:"type:json"`

	// Comparison with previous audits
	PreviousAuditID   *string            `json:"previousAuditId"`
	ComparisonMetrics *ComparisonMetrics `json:"comparisonMetrics" gorm:"embedded"`

	// Configuration and scope
	AuditScope      AuditScope      `json:"auditScope" gorm:"embedded"`
	AuditConfig     AuditConfig     `json:"auditConfig" gorm:"embedded"`

	// Output and artifacts
	ReportPath      *string `json:"reportPath"`      // Path to detailed audit report
	ArtifactsPath   *string `json:"artifactsPath"`   // Path to audit artifacts
	LogsPath        *string `json:"logsPath"`        // Path to audit execution logs
}

// BuildResults contains build status for each platform
type BuildResults struct {
	Backend bool `json:"backend"`
	Web     bool `json:"web"`
	iOS     bool `json:"ios"`
	Android bool `json:"android"`
	KMP     bool `json:"kmp"`
}

// TestResults contains test execution results for each platform
type TestResults struct {
	Backend  TestPlatformResult `json:"backend" gorm:"embedded;embeddedPrefix:backend_"`
	Web      TestPlatformResult `json:"web" gorm:"embedded;embeddedPrefix:web_"`
	iOS      TestPlatformResult `json:"ios" gorm:"embedded;embeddedPrefix:ios_"`
	Android  TestPlatformResult `json:"android" gorm:"embedded;embeddedPrefix:android_"`
	KMP      TestPlatformResult `json:"kmp" gorm:"embedded;embeddedPrefix:kmp_"`
}

// TestPlatformResult contains test results for a specific platform
type TestPlatformResult struct {
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
	Total   int `json:"total"`
	Coverage *float64 `json:"coverage"` // Test coverage percentage
	Duration *int64   `json:"duration"` // Test execution time in milliseconds
}

// PerformanceResults contains performance metrics
type PerformanceResults struct {
	APIResponseTime  *float64 `json:"apiResponseTime"`  // Average API response time in ms
	MobileFrameRate  *float64 `json:"mobileFrameRate"`  // Mobile frame rate in fps
	BuildTime        *float64 `json:"buildTime"`        // Build time in seconds
	MemoryUsage      *float64 `json:"memoryUsage"`      // Memory usage in MB
	BundleSize       *float64 `json:"bundleSize"`       // Bundle size in KB
	LoadTime         *float64 `json:"loadTime"`         // Application load time in seconds

	// Performance thresholds compliance
	APIResponseOK    bool `json:"apiResponseOk"`    // <200ms requirement
	MobileFrameRateOK bool `json:"mobileFrameRateOk"` // >55fps requirement
	BuildTimeOK      bool `json:"buildTimeOk"`      // Build time threshold
	MemoryUsageOK    bool `json:"memoryUsageOk"`    // Memory usage threshold
	BundleSizeOK     bool `json:"bundleSizeOk"`     // Bundle size threshold
	LoadTimeOK       bool `json:"loadTimeOk"`       // Load time threshold
}

// PlaceholderResults contains placeholder audit results
type PlaceholderResults struct {
	TotalFound        int `json:"totalFound"`
	TotalCompleted    int `json:"totalCompleted"`
	TotalRemaining    int `json:"totalRemaining"`
	CriticalRemaining int `json:"criticalRemaining"`
	HighRemaining     int `json:"highRemaining"`
	MediumRemaining   int `json:"mediumRemaining"`
	LowRemaining      int `json:"lowRemaining"`

	// Progress tracking
	CompletionPercentage float64 `json:"completionPercentage"`
	EstimatedHoursRemaining *float64 `json:"estimatedHoursRemaining"`
}

// QualityGateResults contains quality gate validation results
type QualityGateResults struct {
	CodeQualityPassed    bool `json:"codeQualityPassed"`
	SecurityScanPassed   bool `json:"securityScanPassed"`
	PerformancePassed    bool `json:"performancePassed"`
	TestCoveragePassed   bool `json:"testCoveragePassed"`
	DocumentationPassed  bool `json:"documentationPassed"`

	// Quality metrics
	CodeQualityScore    *float64 `json:"codeQualityScore"`
	SecurityScore       *float64 `json:"securityScore"`
	PerformanceScore    *float64 `json:"performanceScore"`
	TestCoverageScore   *float64 `json:"testCoverageScore"`
	DocumentationScore  *float64 `json:"documentationScore"`
}

// AuditViolation represents a validation failure or warning
type AuditViolation struct {
	Type        string `json:"type"`        // ERROR, WARNING, INFO
	Description string `json:"description"`
	Severity    string `json:"severity"`    // CRITICAL, HIGH, MEDIUM, LOW
	Platform    string `json:"platform"`
	Component   string `json:"component"`
	FilePath    *string `json:"filePath"`
	LineNumber  *int   `json:"lineNumber"`

	// Remediation guidance
	Remediation     *string `json:"remediation"`
	DocumentationURL *string `json:"documentationUrl"`
	EstimatedFixTime *float64 `json:"estimatedFixTime"` // Hours
}

// RegionalAuditResult contains audit results specific to a regional market
type RegionalAuditResult struct {
	Region           string  `json:"region"` // TH, SG, MY, ID, PH, VN
	LocalizationScore *float64 `json:"localizationScore"`
	PerformanceScore  *float64 `json:"performanceScore"`
	ComplianceScore   *float64 `json:"complianceScore"`
	Issues           []string `json:"issues"`
	Recommendations  []string `json:"recommendations"`
}

// ComparisonMetrics contains comparison with previous audit results
type ComparisonMetrics struct {
	PlaceholderChange       int     `json:"placeholderChange"`       // Net change in placeholder count
	CompletionPercentageChange float64 `json:"completionPercentageChange"` // Change in completion percentage
	PerformanceChange       *float64 `json:"performanceChange"`       // Change in overall performance score
	QualityChange          *float64 `json:"qualityChange"`           // Change in quality score
	TestCoverageChange     *float64 `json:"testCoverageChange"`      // Change in test coverage

	// Trend indicators
	Trend                   string   `json:"trend"`                   // IMPROVING, DECLINING, STABLE
	TrendConfidence         *float64 `json:"trendConfidence"`         // Confidence in trend analysis (0-1)
}

// AuditScope defines what is included in the audit
type AuditScope struct {
	IncludeTests        bool     `json:"includeTests"`
	IncludePerformance  bool     `json:"includePerformance"`
	IncludeSecurity     bool     `json:"includeSecurity"`
	IncludeDocumentation bool    `json:"includeDocumentation"`
	PlatformFilter      []string `json:"platformFilter"`      // Empty means all platforms
	ServiceFilter       []string `json:"serviceFilter"`       // Empty means all services
	PathFilter          []string `json:"pathFilter"`          // Path patterns to include/exclude
}

// AuditConfig contains configuration parameters for the audit
type AuditConfig struct {
	MaxExecutionTimeMinutes int      `json:"maxExecutionTimeMinutes"` // Maximum allowed execution time
	ParallelExecution      bool      `json:"parallelExecution"`       // Enable parallel execution
	FailFast               bool      `json:"failFast"`                // Stop on first critical failure
	DetailedLogging        bool      `json:"detailedLogging"`         // Enable detailed logging
	GenerateReport         bool      `json:"generateReport"`          // Generate detailed report
	NotifyOnCompletion     bool      `json:"notifyOnCompletion"`      // Send notification on completion

	// Performance thresholds
	APIResponseThresholdMS  *float64 `json:"apiResponseThresholdMs"`  // API response time threshold
	MobileFrameRateThreshold *float64 `json:"mobileFrameRateThreshold"` // Mobile frame rate threshold
	TestCoverageThreshold   *float64 `json:"testCoverageThreshold"`   // Minimum test coverage required

	// Regional configuration
	EnableRegionalOptimization bool     `json:"enableRegionalOptimization"`
	TargetRegions             []string `json:"targetRegions"` // SEA markets to optimize for
}

// AuditStatus defines the possible audit statuses
type AuditStatus string

const (
	AuditStatusRunning   AuditStatus = "RUNNING"
	AuditStatusCompleted AuditStatus = "COMPLETED"
	AuditStatusFailed    AuditStatus = "FAILED"
	AuditStatusCancelled AuditStatus = "CANCELLED"
	AuditStatusPending   AuditStatus = "PENDING"
)

// OverallStatus defines the overall audit outcome
type OverallStatus string

const (
	OverallStatusPass    OverallStatus = "PASS"
	OverallStatusFail    OverallStatus = "FAIL"
	OverallStatusWarning OverallStatus = "WARNING"
	OverallStatusUnknown OverallStatus = "UNKNOWN"
)

// ExecutionMode defines how the audit was triggered
type ExecutionMode string

const (
	ExecutionModeManual    ExecutionMode = "MANUAL"
	ExecutionModeScheduled ExecutionMode = "SCHEDULED"
	ExecutionModeTriggered ExecutionMode = "TRIGGERED"
	ExecutionModeCI        ExecutionMode = "CI"
)

// GetExecutionDuration returns the audit execution duration in milliseconds
func (ca *CompletionAudit) GetExecutionDuration() int64 {
	if ca.ExecutionTimeMS != nil {
		return *ca.ExecutionTimeMS
	}
	if ca.StartedAt != nil && ca.CompletedAt != nil {
		return ca.CompletedAt.Sub(*ca.StartedAt).Milliseconds()
	}
	return 0
}

// IsCompleted returns true if the audit has completed (successfully or with failures)
func (ca *CompletionAudit) IsCompleted() bool {
	return ca.Status == string(AuditStatusCompleted) && ca.CompletedAt != nil
}

// HasCriticalViolations returns true if there are any critical violations
func (ca *CompletionAudit) HasCriticalViolations() bool {
	for _, violation := range ca.Violations {
		if violation.Severity == "CRITICAL" {
			return true
		}
	}
	return false
}

// GetOverallCompletionPercentage calculates the overall completion percentage
func (ca *CompletionAudit) GetOverallCompletionPercentage() float64 {
	return ca.PlaceholderResults.CompletionPercentage
}

// GetQualityScore calculates an overall quality score based on all metrics
func (ca *CompletionAudit) GetQualityScore() float64 {
	scores := []float64{}

	if ca.QualityGateResults.CodeQualityScore != nil {
		scores = append(scores, *ca.QualityGateResults.CodeQualityScore)
	}
	if ca.QualityGateResults.SecurityScore != nil {
		scores = append(scores, *ca.QualityGateResults.SecurityScore)
	}
	if ca.QualityGateResults.PerformanceScore != nil {
		scores = append(scores, *ca.QualityGateResults.PerformanceScore)
	}
	if ca.QualityGateResults.TestCoverageScore != nil {
		scores = append(scores, *ca.QualityGateResults.TestCoverageScore)
	}

	if len(scores) == 0 {
		return 0
	}

	total := 0.0
	for _, score := range scores {
		total += score
	}
	return total / float64(len(scores))
}

// GetRegionalOptimizationScore calculates the regional optimization score for SEA markets
func (ca *CompletionAudit) GetRegionalOptimizationScore() float64 {
	if len(ca.RegionalResults) == 0 {
		return 0
	}

	seaMarkets := []string{"TH", "SG", "MY", "ID", "PH", "VN"}
	totalScore := 0.0
	count := 0

	for _, market := range seaMarkets {
		if result, exists := ca.RegionalResults[market]; exists {
			if result.PerformanceScore != nil {
				totalScore += *result.PerformanceScore
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}
	return totalScore / float64(count)
}