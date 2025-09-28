package coverage_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"tchat.dev/tests/fixtures"
)

// CoverageConfigTestSuite provides comprehensive coverage reporting configuration testing
// for monitoring and enforcing test coverage across all Tchat microservices
type CoverageConfigTestSuite struct {
	suite.Suite
	fixtures    *fixtures.MasterFixtures
	ctx         context.Context
	projectRoot string
}

// CoverageReport represents a coverage report structure
type CoverageReport struct {
	Service         string                         `json:"service"`
	Timestamp       time.Time                      `json:"timestamp"`
	TotalCoverage   float64                        `json:"total_coverage"`
	PackageCoverage map[string]PackageCoverageInfo `json:"package_coverage"`
	Thresholds      CoverageThresholds             `json:"thresholds"`
	Status          string                         `json:"status"`
	Violations      []CoverageViolation            `json:"violations"`
	Metadata        CoverageMetadata               `json:"metadata"`
}

// PackageCoverageInfo contains coverage information for a package
type PackageCoverageInfo struct {
	Package           string                 `json:"package"`
	Coverage          float64                `json:"coverage"`
	TotalStatements   int                    `json:"total_statements"`
	CoveredStatements int                    `json:"covered_statements"`
	UncoveredLines    []int                  `json:"uncovered_lines"`
	Functions         []FunctionCoverageInfo `json:"functions"`
}

// FunctionCoverageInfo contains coverage information for a function
type FunctionCoverageInfo struct {
	Name              string  `json:"name"`
	File              string  `json:"file"`
	StartLine         int     `json:"start_line"`
	EndLine           int     `json:"end_line"`
	Coverage          float64 `json:"coverage"`
	CoveredStatements int     `json:"covered_statements"`
	TotalStatements   int     `json:"total_statements"`
}

// CoverageThresholds defines coverage requirements
type CoverageThresholds struct {
	TotalCoverage     float64            `json:"total_coverage"`
	PackageCoverage   float64            `json:"package_coverage"`
	FunctionCoverage  float64            `json:"function_coverage"`
	ServiceThresholds map[string]float64 `json:"service_thresholds"`
}

// CoverageViolation represents a coverage threshold violation
type CoverageViolation struct {
	Type        string  `json:"type"`
	Package     string  `json:"package,omitempty"`
	Function    string  `json:"function,omitempty"`
	Current     float64 `json:"current"`
	Required    float64 `json:"required"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
}

// CoverageMetadata contains additional coverage metadata
type CoverageMetadata struct {
	GoVersion     string            `json:"go_version"`
	TestSuite     string            `json:"test_suite"`
	Environment   string            `json:"environment"`
	BuildNumber   string            `json:"build_number,omitempty"`
	CommitHash    string            `json:"commit_hash,omitempty"`
	TestDuration  time.Duration     `json:"test_duration"`
	TestCount     int               `json:"test_count"`
	ExcludedFiles []string          `json:"excluded_files"`
	Tags          map[string]string `json:"tags,omitempty"`
}

// CoverageConfig defines coverage configuration
type CoverageConfig struct {
	Enabled         bool               `json:"enabled"`
	OutputFormat    []string           `json:"output_format"` // html, json, text, lcov
	OutputDir       string             `json:"output_dir"`
	Thresholds      CoverageThresholds `json:"thresholds"`
	ExcludePatterns []string           `json:"exclude_patterns"`
	IncludePatterns []string           `json:"include_patterns"`
	FailOnViolation bool               `json:"fail_on_violation"`
	ReportFormats   []string           `json:"report_formats"`
	Integration     IntegrationConfig  `json:"integration"`
}

// IntegrationConfig defines external integration configuration
type IntegrationConfig struct {
	Codecov     CodecovConfig      `json:"codecov"`
	Coveralls   CoverallsConfig    `json:"coveralls"`
	SonarQube   SonarQubeConfig    `json:"sonarqube"`
	CustomHooks []CustomHookConfig `json:"custom_hooks"`
}

// CodecovConfig defines Codecov integration
type CodecovConfig struct {
	Enabled     bool   `json:"enabled"`
	Token       string `json:"token,omitempty"`
	UploadURL   string `json:"upload_url,omitempty"`
	ProjectName string `json:"project_name"`
}

// CoverallsConfig defines Coveralls integration
type CoverallsConfig struct {
	Enabled     bool   `json:"enabled"`
	Token       string `json:"token,omitempty"`
	ServiceName string `json:"service_name"`
}

// SonarQubeConfig defines SonarQube integration
type SonarQubeConfig struct {
	Enabled    bool   `json:"enabled"`
	ServerURL  string `json:"server_url"`
	ProjectKey string `json:"project_key"`
	Token      string `json:"token,omitempty"`
}

// CustomHookConfig defines custom webhook integration
type CustomHookConfig struct {
	Name            string            `json:"name"`
	URL             string            `json:"url"`
	Method          string            `json:"method"`
	Headers         map[string]string `json:"headers"`
	PayloadTemplate string            `json:"payload_template"`
}

// SetupSuite initializes the test suite
func (suite *CoverageConfigTestSuite) SetupSuite() {
	suite.fixtures = fixtures.NewMasterFixtures(12345)
	suite.ctx = context.Background()
	suite.projectRoot = "/Users/weerawat/Tchat/backend"
}

// TestCoverageConfigurationValidation tests coverage configuration validation
func (suite *CoverageConfigTestSuite) TestCoverageConfigurationValidation() {
	testCases := []struct {
		name        string
		config      CoverageConfig
		expectValid bool
		description string
	}{
		{
			name: "Valid basic configuration",
			config: CoverageConfig{
				Enabled:      true,
				OutputFormat: []string{"html", "json"},
				OutputDir:    "coverage",
				Thresholds: CoverageThresholds{
					TotalCoverage:    80.0,
					PackageCoverage:  70.0,
					FunctionCoverage: 60.0,
				},
				FailOnViolation: true,
			},
			expectValid: true,
			description: "Basic configuration should be valid",
		},
		{
			name: "Configuration with service-specific thresholds",
			config: CoverageConfig{
				Enabled:      true,
				OutputFormat: []string{"html", "json", "lcov"},
				OutputDir:    "coverage",
				Thresholds: CoverageThresholds{
					TotalCoverage:    85.0,
					PackageCoverage:  75.0,
					FunctionCoverage: 65.0,
					ServiceThresholds: map[string]float64{
						"auth":      90.0,
						"content":   85.0,
						"messaging": 80.0,
						"payment":   95.0,
					},
				},
				ExcludePatterns: []string{
					"*_test.go",
					"*/mocks/*",
					"*/vendor/*",
					"main.go",
				},
				IncludePatterns: []string{
					"*/handlers/*",
					"*/services/*",
					"*/repositories/*",
				},
				FailOnViolation: true,
			},
			expectValid: true,
			description: "Configuration with service thresholds should be valid",
		},
		{
			name: "Configuration with integrations",
			config: CoverageConfig{
				Enabled:      true,
				OutputFormat: []string{"lcov", "json"},
				OutputDir:    "coverage",
				Thresholds: CoverageThresholds{
					TotalCoverage:    80.0,
					PackageCoverage:  70.0,
					FunctionCoverage: 60.0,
				},
				Integration: IntegrationConfig{
					Codecov: CodecovConfig{
						Enabled:     true,
						ProjectName: "tchat.dev",
					},
					Coveralls: CoverallsConfig{
						Enabled:     true,
						ServiceName: "github",
					},
					SonarQube: SonarQubeConfig{
						Enabled:    true,
						ServerURL:  "https://sonarqube.tchat.dev",
						ProjectKey: "tchat.dev",
					},
				},
			},
			expectValid: true,
			description: "Configuration with integrations should be valid",
		},
		{
			name: "Invalid threshold values",
			config: CoverageConfig{
				Enabled:      true,
				OutputFormat: []string{"html"},
				OutputDir:    "coverage",
				Thresholds: CoverageThresholds{
					TotalCoverage:    150.0, // Invalid: > 100
					PackageCoverage:  -10.0, // Invalid: < 0
					FunctionCoverage: 60.0,
				},
			},
			expectValid: false,
			description: "Invalid threshold values should be rejected",
		},
		{
			name: "Empty output format",
			config: CoverageConfig{
				Enabled:      true,
				OutputFormat: []string{}, // Invalid: empty
				OutputDir:    "coverage",
				Thresholds: CoverageThresholds{
					TotalCoverage: 80.0,
				},
			},
			expectValid: false,
			description: "Empty output format should be invalid",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isValid := suite.validateCoverageConfig(tc.config)
			suite.Equal(tc.expectValid, isValid, tc.description)

			if tc.expectValid {
				// Test JSON serialization for valid configs
				jsonData, err := json.Marshal(tc.config)
				suite.NoError(err, "Valid config should serialize to JSON")

				var deserializedConfig CoverageConfig
				err = json.Unmarshal(jsonData, &deserializedConfig)
				suite.NoError(err, "Valid config should deserialize from JSON")

				suite.Equal(tc.config.Enabled, deserializedConfig.Enabled)
				suite.Equal(tc.config.OutputFormat, deserializedConfig.OutputFormat)
			}
		})
	}
}

// TestCoverageThresholdEnforcement tests coverage threshold enforcement
func (suite *CoverageConfigTestSuite) TestCoverageThresholdEnforcement() {
	testCases := []struct {
		name             string
		coverageReport   CoverageReport
		thresholds       CoverageThresholds
		expectViolations []CoverageViolation
		description      string
	}{
		{
			name: "All thresholds met",
			coverageReport: CoverageReport{
				Service:       "auth",
				TotalCoverage: 85.0,
				PackageCoverage: map[string]PackageCoverageInfo{
					"handlers":     {Coverage: 90.0},
					"services":     {Coverage: 85.0},
					"repositories": {Coverage: 80.0},
				},
			},
			thresholds: CoverageThresholds{
				TotalCoverage:   80.0,
				PackageCoverage: 75.0,
			},
			expectViolations: []CoverageViolation{},
			description:      "No violations when all thresholds are met",
		},
		{
			name: "Total coverage violation",
			coverageReport: CoverageReport{
				Service:       "content",
				TotalCoverage: 70.0,
				PackageCoverage: map[string]PackageCoverageInfo{
					"handlers": {Coverage: 75.0},
					"services": {Coverage: 80.0},
				},
			},
			thresholds: CoverageThresholds{
				TotalCoverage:   80.0,
				PackageCoverage: 70.0,
			},
			expectViolations: []CoverageViolation{
				{
					Type:        "total_coverage",
					Current:     70.0,
					Required:    80.0,
					Severity:    "high",
					Description: "Total coverage below threshold",
				},
			},
			description: "Violation when total coverage is below threshold",
		},
		{
			name: "Package coverage violations",
			coverageReport: CoverageReport{
				Service:       "messaging",
				TotalCoverage: 85.0,
				PackageCoverage: map[string]PackageCoverageInfo{
					"handlers":     {Coverage: 65.0}, // Below threshold
					"services":     {Coverage: 80.0}, // OK
					"repositories": {Coverage: 60.0}, // Below threshold
				},
			},
			thresholds: CoverageThresholds{
				TotalCoverage:   75.0,
				PackageCoverage: 70.0,
			},
			expectViolations: []CoverageViolation{
				{
					Type:        "package_coverage",
					Package:     "handlers",
					Current:     65.0,
					Required:    70.0,
					Severity:    "medium",
					Description: "Package coverage below threshold",
				},
				{
					Type:        "package_coverage",
					Package:     "repositories",
					Current:     60.0,
					Required:    70.0,
					Severity:    "medium",
					Description: "Package coverage below threshold",
				},
			},
			description: "Violations for packages below threshold",
		},
		{
			name: "Service-specific threshold violation",
			coverageReport: CoverageReport{
				Service:       "payment",
				TotalCoverage: 85.0,
				PackageCoverage: map[string]PackageCoverageInfo{
					"handlers": {Coverage: 90.0},
				},
			},
			thresholds: CoverageThresholds{
				TotalCoverage:   80.0,
				PackageCoverage: 75.0,
				ServiceThresholds: map[string]float64{
					"payment": 95.0, // Higher threshold for payment service
				},
			},
			expectViolations: []CoverageViolation{
				{
					Type:        "service_coverage",
					Package:     "payment",
					Current:     85.0,
					Required:    95.0,
					Severity:    "high",
					Description: "Service coverage below service-specific threshold",
				},
			},
			description: "Violation for service-specific threshold",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			violations := suite.evaluateCoverageThresholds(tc.coverageReport, tc.thresholds)

			suite.Len(violations, len(tc.expectViolations), tc.description)

			for i, expectedViolation := range tc.expectViolations {
				if i < len(violations) {
					suite.Equal(expectedViolation.Type, violations[i].Type)
					suite.Equal(expectedViolation.Current, violations[i].Current)
					suite.Equal(expectedViolation.Required, violations[i].Required)
					suite.Equal(expectedViolation.Severity, violations[i].Severity)
				}
			}
		})
	}
}

// TestCoverageReportGeneration tests coverage report generation
func (suite *CoverageConfigTestSuite) TestCoverageReportGeneration() {
	testCases := []struct {
		name        string
		service     string
		format      string
		expectValid bool
		description string
	}{
		{
			name:        "HTML report generation",
			service:     "auth",
			format:      "html",
			expectValid: true,
			description: "HTML coverage reports should be generated",
		},
		{
			name:        "JSON report generation",
			service:     "content",
			format:      "json",
			expectValid: true,
			description: "JSON coverage reports should be generated",
		},
		{
			name:        "LCOV report generation",
			service:     "messaging",
			format:      "lcov",
			expectValid: true,
			description: "LCOV coverage reports should be generated",
		},
		{
			name:        "Text report generation",
			service:     "payment",
			format:      "text",
			expectValid: true,
			description: "Text coverage reports should be generated",
		},
		{
			name:        "Invalid format",
			service:     "notification",
			format:      "invalid",
			expectValid: false,
			description: "Invalid formats should be rejected",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			report := suite.generateMockCoverageReport(tc.service)

			if tc.expectValid {
				// Test report generation
				generatedReport := suite.formatCoverageReport(report, tc.format)
				suite.NotEmpty(generatedReport, tc.description)

				// Validate report content based on format
				switch tc.format {
				case "json":
					var jsonReport CoverageReport
					err := json.Unmarshal([]byte(generatedReport), &jsonReport)
					suite.NoError(err, "JSON report should be valid JSON")
					suite.Equal(tc.service, jsonReport.Service)

				case "html":
					suite.Contains(generatedReport, "<html>", "HTML report should contain HTML tags")
					suite.Contains(generatedReport, tc.service, "HTML report should contain service name")

				case "lcov":
					suite.Contains(generatedReport, "TN:", "LCOV report should contain test name")
					suite.Contains(generatedReport, "SF:", "LCOV report should contain source file info")

				case "text":
					suite.Contains(generatedReport, "Coverage Report", "Text report should contain header")
					suite.Contains(generatedReport, tc.service, "Text report should contain service name")
				}
			} else {
				// Test invalid format handling
				generatedReport := suite.formatCoverageReport(report, tc.format)
				suite.Empty(generatedReport, tc.description)
			}
		})
	}
}

// TestCoverageIntegrationConfiguration tests external integration configuration
func (suite *CoverageConfigTestSuite) TestCoverageIntegrationConfiguration() {
	integrationTests := []struct {
		name        string
		integration IntegrationConfig
		expectValid bool
		description string
	}{
		{
			name: "Codecov integration",
			integration: IntegrationConfig{
				Codecov: CodecovConfig{
					Enabled:     true,
					ProjectName: "tchat.dev",
					UploadURL:   "https://codecov.io/upload",
				},
			},
			expectValid: true,
			description: "Codecov integration should be configurable",
		},
		{
			name: "Coveralls integration",
			integration: IntegrationConfig{
				Coveralls: CoverallsConfig{
					Enabled:     true,
					ServiceName: "github",
				},
			},
			expectValid: true,
			description: "Coveralls integration should be configurable",
		},
		{
			name: "SonarQube integration",
			integration: IntegrationConfig{
				SonarQube: SonarQubeConfig{
					Enabled:    true,
					ServerURL:  "https://sonarqube.tchat.dev",
					ProjectKey: "tchat.dev",
				},
			},
			expectValid: true,
			description: "SonarQube integration should be configurable",
		},
		{
			name: "Custom webhook integration",
			integration: IntegrationConfig{
				CustomHooks: []CustomHookConfig{
					{
						Name:   "slack-notification",
						URL:    "https://hooks.slack.com/services/...",
						Method: "POST",
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						PayloadTemplate: `{"text": "Coverage: {{.TotalCoverage}}%"}`,
					},
				},
			},
			expectValid: true,
			description: "Custom webhook integration should be configurable",
		},
		{
			name: "Multiple integrations",
			integration: IntegrationConfig{
				Codecov: CodecovConfig{
					Enabled:     true,
					ProjectName: "tchat.dev",
				},
				SonarQube: SonarQubeConfig{
					Enabled:    true,
					ServerURL:  "https://sonarqube.tchat.dev",
					ProjectKey: "tchat.dev",
				},
				CustomHooks: []CustomHookConfig{
					{
						Name:   "teams-notification",
						URL:    "https://outlook.office.com/webhook/...",
						Method: "POST",
					},
				},
			},
			expectValid: true,
			description: "Multiple integrations should be supported",
		},
	}

	for _, tt := range integrationTests {
		suite.Run(tt.name, func() {
			isValid := suite.validateIntegrationConfig(tt.integration)
			suite.Equal(tt.expectValid, isValid, tt.description)

			if tt.expectValid {
				// Test integration payload generation
				mockReport := suite.generateMockCoverageReport("test-service")
				payload := suite.generateIntegrationPayload(mockReport, tt.integration)
				suite.NotEmpty(payload, "Integration payload should be generated")
			}
		})
	}
}

// TestCoverageFilePatternMatching tests file inclusion/exclusion patterns
func (suite *CoverageConfigTestSuite) TestCoverageFilePatternMatching() {
	testFiles := []string{
		"auth/handlers/auth_handler.go",
		"auth/handlers/auth_handler_test.go",
		"auth/services/auth_service.go",
		"auth/services/auth_service_test.go",
		"auth/repositories/user_repository.go",
		"auth/mocks/mock_service.go",
		"content/handlers/content_handler.go",
		"vendor/github.com/some/package.go",
		"main.go",
		"cmd/server/main.go",
	}

	patternTests := []struct {
		name            string
		includePatterns []string
		excludePatterns []string
		expectedFiles   []string
		description     string
	}{
		{
			name: "Include handlers only",
			includePatterns: []string{
				"*/handlers/*",
			},
			excludePatterns: []string{
				"*_test.go",
			},
			expectedFiles: []string{
				"auth/handlers/auth_handler.go",
				"content/handlers/content_handler.go",
			},
			description: "Should include only handler files, excluding tests",
		},
		{
			name: "Exclude test files and mocks",
			includePatterns: []string{
				"*/*.go",
			},
			excludePatterns: []string{
				"*_test.go",
				"*/mocks/*",
				"*/vendor/*",
				"main.go",
			},
			expectedFiles: []string{
				"auth/handlers/auth_handler.go",
				"auth/services/auth_service.go",
				"auth/repositories/user_repository.go",
				"content/handlers/content_handler.go",
				"cmd/server/main.go",
			},
			description: "Should exclude tests, mocks, vendor, and main.go",
		},
		{
			name: "Include specific services",
			includePatterns: []string{
				"auth/services/*",
				"auth/repositories/*",
			},
			excludePatterns: []string{
				"*_test.go",
			},
			expectedFiles: []string{
				"auth/services/auth_service.go",
				"auth/repositories/user_repository.go",
			},
			description: "Should include only auth services and repositories",
		},
	}

	for _, pt := range patternTests {
		suite.Run(pt.name, func() {
			matchedFiles := suite.applyFilePatterns(testFiles, pt.includePatterns, pt.excludePatterns)

			suite.ElementsMatch(pt.expectedFiles, matchedFiles, pt.description)

			// Verify that excluded patterns are actually excluded
			for _, excludePattern := range pt.excludePatterns {
				for _, file := range matchedFiles {
					matched, _ := filepath.Match(excludePattern, file)
					suite.False(matched, "File %s should not match exclude pattern %s", file, excludePattern)
				}
			}
		})
	}
}

// TestCoverageQualityGates tests coverage quality gates enforcement
func (suite *CoverageConfigTestSuite) TestCoverageQualityGates() {
	qualityGateTests := []struct {
		name           string
		coverageReport CoverageReport
		config         CoverageConfig
		shouldPass     bool
		description    string
	}{
		{
			name: "Quality gate passes",
			coverageReport: CoverageReport{
				Service:       "auth",
				TotalCoverage: 90.0,
				PackageCoverage: map[string]PackageCoverageInfo{
					"handlers": {Coverage: 95.0},
					"services": {Coverage: 90.0},
				},
			},
			config: CoverageConfig{
				Thresholds: CoverageThresholds{
					TotalCoverage:   80.0,
					PackageCoverage: 75.0,
				},
				FailOnViolation: true,
			},
			shouldPass:  true,
			description: "Quality gate should pass when all thresholds are met",
		},
		{
			name: "Quality gate fails on total coverage",
			coverageReport: CoverageReport{
				Service:       "content",
				TotalCoverage: 70.0,
				PackageCoverage: map[string]PackageCoverageInfo{
					"handlers": {Coverage: 80.0},
				},
			},
			config: CoverageConfig{
				Thresholds: CoverageThresholds{
					TotalCoverage: 80.0,
				},
				FailOnViolation: true,
			},
			shouldPass:  false,
			description: "Quality gate should fail when total coverage is below threshold",
		},
		{
			name: "Quality gate warning only",
			coverageReport: CoverageReport{
				Service:       "messaging",
				TotalCoverage: 70.0,
				PackageCoverage: map[string]PackageCoverageInfo{
					"handlers": {Coverage: 75.0},
				},
			},
			config: CoverageConfig{
				Thresholds: CoverageThresholds{
					TotalCoverage: 80.0,
				},
				FailOnViolation: false, // Warning only
			},
			shouldPass:  true,
			description: "Quality gate should pass with warning when FailOnViolation is false",
		},
	}

	for _, qgt := range qualityGateTests {
		suite.Run(qgt.name, func() {
			passed := suite.evaluateQualityGate(qgt.coverageReport, qgt.config)
			suite.Equal(qgt.shouldPass, passed, qgt.description)

			// Test quality gate status
			status := suite.getQualityGateStatus(qgt.coverageReport, qgt.config)
			if qgt.shouldPass {
				suite.Contains([]string{"passed", "warning"}, status)
			} else {
				suite.Equal("failed", status)
			}
		})
	}
}

// Helper methods

// validateCoverageConfig validates coverage configuration
func (suite *CoverageConfigTestSuite) validateCoverageConfig(config CoverageConfig) bool {
	// Check required fields
	if len(config.OutputFormat) == 0 {
		return false
	}

	// Validate threshold ranges
	if config.Thresholds.TotalCoverage < 0 || config.Thresholds.TotalCoverage > 100 {
		return false
	}
	if config.Thresholds.PackageCoverage < 0 || config.Thresholds.PackageCoverage > 100 {
		return false
	}
	if config.Thresholds.FunctionCoverage < 0 || config.Thresholds.FunctionCoverage > 100 {
		return false
	}

	// Validate service thresholds
	for _, threshold := range config.Thresholds.ServiceThresholds {
		if threshold < 0 || threshold > 100 {
			return false
		}
	}

	// Validate output formats
	validFormats := map[string]bool{
		"html": true, "json": true, "lcov": true, "text": true,
	}
	for _, format := range config.OutputFormat {
		if !validFormats[format] {
			return false
		}
	}

	return true
}

// evaluateCoverageThresholds evaluates coverage against thresholds
func (suite *CoverageConfigTestSuite) evaluateCoverageThresholds(report CoverageReport, thresholds CoverageThresholds) []CoverageViolation {
	var violations []CoverageViolation

	// Check total coverage
	if report.TotalCoverage < thresholds.TotalCoverage {
		violations = append(violations, CoverageViolation{
			Type:        "total_coverage",
			Current:     report.TotalCoverage,
			Required:    thresholds.TotalCoverage,
			Severity:    "high",
			Description: "Total coverage below threshold",
		})
	}

	// Check service-specific thresholds
	if serviceThreshold, exists := thresholds.ServiceThresholds[report.Service]; exists {
		if report.TotalCoverage < serviceThreshold {
			violations = append(violations, CoverageViolation{
				Type:        "service_coverage",
				Package:     report.Service,
				Current:     report.TotalCoverage,
				Required:    serviceThreshold,
				Severity:    "high",
				Description: "Service coverage below service-specific threshold",
			})
		}
	}

	// Check package coverage
	for packageName, packageInfo := range report.PackageCoverage {
		if packageInfo.Coverage < thresholds.PackageCoverage {
			violations = append(violations, CoverageViolation{
				Type:        "package_coverage",
				Package:     packageName,
				Current:     packageInfo.Coverage,
				Required:    thresholds.PackageCoverage,
				Severity:    "medium",
				Description: "Package coverage below threshold",
			})
		}
	}

	return violations
}

// generateMockCoverageReport generates a mock coverage report for testing
func (suite *CoverageConfigTestSuite) generateMockCoverageReport(service string) CoverageReport {
	return CoverageReport{
		Service:       service,
		Timestamp:     time.Now(),
		TotalCoverage: 82.5,
		PackageCoverage: map[string]PackageCoverageInfo{
			"handlers": {
				Package:           fmt.Sprintf("%s/handlers", service),
				Coverage:          85.0,
				TotalStatements:   100,
				CoveredStatements: 85,
				UncoveredLines:    []int{23, 45, 67},
				Functions: []FunctionCoverageInfo{
					{
						Name:              "HandleRequest",
						File:              fmt.Sprintf("%s/handlers/handler.go", service),
						StartLine:         10,
						EndLine:           50,
						Coverage:          90.0,
						CoveredStatements: 18,
						TotalStatements:   20,
					},
				},
			},
			"services": {
				Package:           fmt.Sprintf("%s/services", service),
				Coverage:          80.0,
				TotalStatements:   150,
				CoveredStatements: 120,
				UncoveredLines:    []int{78, 89, 134, 145},
			},
		},
		Thresholds: CoverageThresholds{
			TotalCoverage:   80.0,
			PackageCoverage: 75.0,
		},
		Status: "passed",
		Metadata: CoverageMetadata{
			GoVersion:    "1.21.0",
			TestSuite:    "tchat-test-suite",
			Environment:  "test",
			TestDuration: 30 * time.Second,
			TestCount:    150,
		},
	}
}

// formatCoverageReport formats coverage report in specified format
func (suite *CoverageConfigTestSuite) formatCoverageReport(report CoverageReport, format string) string {
	switch format {
	case "json":
		jsonData, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return ""
		}
		return string(jsonData)

	case "html":
		return suite.generateHTMLReport(report)

	case "lcov":
		return suite.generateLCOVReport(report)

	case "text":
		return suite.generateTextReport(report)

	default:
		return ""
	}
}

// generateHTMLReport generates HTML coverage report
func (suite *CoverageConfigTestSuite) generateHTMLReport(report CoverageReport) string {
	html := fmt.Sprintf(`
<html>
<head>
	<title>Coverage Report - %s</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 20px; }
		.header { background-color: #f5f5f5; padding: 10px; }
		.coverage { margin: 10px 0; }
		.high { color: green; }
		.medium { color: orange; }
		.low { color: red; }
	</style>
</head>
<body>
	<div class="header">
		<h1>Coverage Report</h1>
		<p>Service: %s</p>
		<p>Generated: %s</p>
	</div>
	<div class="coverage">
		<h2>Total Coverage: %.2f%%</h2>
	</div>
	<div class="packages">
		<h3>Package Coverage</h3>`,
		report.Service, report.Service, report.Timestamp.Format(time.RFC3339), report.TotalCoverage)

	for packageName, packageInfo := range report.PackageCoverage {
		coverageClass := "low"
		if packageInfo.Coverage >= 80 {
			coverageClass = "high"
		} else if packageInfo.Coverage >= 60 {
			coverageClass = "medium"
		}

		html += fmt.Sprintf(`
		<div class="package">
			<h4>%s</h4>
			<p class="%s">Coverage: %.2f%% (%d/%d statements)</p>
		</div>`, packageName, coverageClass, packageInfo.Coverage, packageInfo.CoveredStatements, packageInfo.TotalStatements)
	}

	html += `
	</div>
</body>
</html>`

	return html
}

// generateLCOVReport generates LCOV format coverage report
func (suite *CoverageConfigTestSuite) generateLCOVReport(report CoverageReport) string {
	lcov := fmt.Sprintf("TN:%s\n", report.Service)

	for packageName, packageInfo := range report.PackageCoverage {
		for _, function := range packageInfo.Functions {
			lcov += fmt.Sprintf("SF:%s\n", function.File)
			lcov += fmt.Sprintf("FN:%d,%s\n", function.StartLine, function.Name)
			lcov += fmt.Sprintf("FNDA:%d,%s\n", function.CoveredStatements, function.Name)
			lcov += fmt.Sprintf("FNF:%d\n", function.TotalStatements)
			lcov += fmt.Sprintf("FNH:%d\n", function.CoveredStatements)

			// Add line coverage data
			for line := function.StartLine; line <= function.EndLine; line++ {
				covered := 1
				for _, uncoveredLine := range packageInfo.UncoveredLines {
					if uncoveredLine == line {
						covered = 0
						break
					}
				}
				lcov += fmt.Sprintf("DA:%d,%d\n", line, covered)
			}

			lcov += fmt.Sprintf("LF:%d\n", function.TotalStatements)
			lcov += fmt.Sprintf("LH:%d\n", function.CoveredStatements)
			lcov += "end_of_record\n"
		}
	}

	return lcov
}

// generateTextReport generates text format coverage report
func (suite *CoverageConfigTestSuite) generateTextReport(report CoverageReport) string {
	text := fmt.Sprintf(`
Coverage Report
===============

Service: %s
Generated: %s
Total Coverage: %.2f%%

Package Coverage:
`, report.Service, report.Timestamp.Format(time.RFC3339), report.TotalCoverage)

	for packageName, packageInfo := range report.PackageCoverage {
		text += fmt.Sprintf("  %s: %.2f%% (%d/%d statements)\n",
			packageName, packageInfo.Coverage, packageInfo.CoveredStatements, packageInfo.TotalStatements)
	}

	if len(report.Violations) > 0 {
		text += "\nThreshold Violations:\n"
		for _, violation := range report.Violations {
			text += fmt.Sprintf("  [%s] %s: %.2f%% (required: %.2f%%)\n",
				violation.Severity, violation.Description, violation.Current, violation.Required)
		}
	}

	return text
}

// validateIntegrationConfig validates integration configuration
func (suite *CoverageConfigTestSuite) validateIntegrationConfig(config IntegrationConfig) bool {
	// Validate Codecov config
	if config.Codecov.Enabled && config.Codecov.ProjectName == "" {
		return false
	}

	// Validate SonarQube config
	if config.SonarQube.Enabled && (config.SonarQube.ServerURL == "" || config.SonarQube.ProjectKey == "") {
		return false
	}

	// Validate custom hooks
	for _, hook := range config.CustomHooks {
		if hook.Name == "" || hook.URL == "" || hook.Method == "" {
			return false
		}
	}

	return true
}

// generateIntegrationPayload generates payload for external integrations
func (suite *CoverageConfigTestSuite) generateIntegrationPayload(report CoverageReport, config IntegrationConfig) string {
	// Simple JSON payload for testing
	payload := map[string]interface{}{
		"service":   report.Service,
		"coverage":  report.TotalCoverage,
		"timestamp": report.Timestamp,
		"status":    report.Status,
	}

	jsonData, _ := json.Marshal(payload)
	return string(jsonData)
}

// applyFilePatterns applies include/exclude patterns to file list
func (suite *CoverageConfigTestSuite) applyFilePatterns(files []string, includePatterns []string, excludePatterns []string) []string {
	var matchedFiles []string

	for _, file := range files {
		// Check include patterns
		included := len(includePatterns) == 0 // Include all if no patterns specified
		for _, pattern := range includePatterns {
			if matched, _ := filepath.Match(pattern, file); matched {
				included = true
				break
			}
		}

		if !included {
			continue
		}

		// Check exclude patterns
		excluded := false
		for _, pattern := range excludePatterns {
			if matched, _ := filepath.Match(pattern, file); matched {
				excluded = true
				break
			}
		}

		if !excluded {
			matchedFiles = append(matchedFiles, file)
		}
	}

	return matchedFiles
}

// evaluateQualityGate evaluates coverage quality gate
func (suite *CoverageConfigTestSuite) evaluateQualityGate(report CoverageReport, config CoverageConfig) bool {
	violations := suite.evaluateCoverageThresholds(report, config.Thresholds)

	if len(violations) == 0 {
		return true
	}

	return !config.FailOnViolation
}

// getQualityGateStatus gets quality gate status
func (suite *CoverageConfigTestSuite) getQualityGateStatus(report CoverageReport, config CoverageConfig) string {
	violations := suite.evaluateCoverageThresholds(report, config.Thresholds)

	if len(violations) == 0 {
		return "passed"
	}

	if config.FailOnViolation {
		return "failed"
	}

	return "warning"
}

func TestCoverageConfigSuite(t *testing.T) {
	suite.Run(t, new(CoverageConfigTestSuite))
}
