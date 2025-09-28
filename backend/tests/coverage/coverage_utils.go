package coverage

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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

// CoverageThresholds defines coverage requirements
type CoverageThresholds struct {
	TotalCoverage     float64            `json:"total_coverage"`
	PackageCoverage   float64            `json:"package_coverage"`
	FunctionCoverage  float64            `json:"function_coverage"`
	ServiceThresholds map[string]float64 `json:"service_thresholds"`
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

// CoverageViolation represents a coverage threshold violation
type CoverageViolation struct {
	Type        string   `json:"type"`
	Package     string   `json:"package,omitempty"`
	Function    string   `json:"function,omitempty"`
	Current     float64  `json:"current"`
	Expected    float64  `json:"expected"`
	Severity    string   `json:"severity"`
	Message     string   `json:"message"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// CoverageMetadata contains additional report metadata
type CoverageMetadata struct {
	GoVersion     string            `json:"go_version"`
	Timestamp     time.Time         `json:"timestamp"`
	ExecutionTime string            `json:"execution_time"`
	TotalPackages int               `json:"total_packages"`
	ExcludedFiles []string          `json:"excluded_files"`
	ReportFormat  string            `json:"report_format"`
	Environment   string            `json:"environment,omitempty"`
	CommitHash    string            `json:"commit_hash,omitempty"`
	Branch        string            `json:"branch,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}

// CoverageAnalyzer provides coverage analysis utilities
type CoverageAnalyzer struct {
	projectRoot string
	config      CoverageConfig
}

// NewCoverageAnalyzer creates a new coverage analyzer
func NewCoverageAnalyzer(projectRoot string, config CoverageConfig) *CoverageAnalyzer {
	return &CoverageAnalyzer{
		projectRoot: projectRoot,
		config:      config,
	}
}

// RunCoverageAnalysis runs comprehensive coverage analysis
func (ca *CoverageAnalyzer) RunCoverageAnalysis(service string) (*CoverageReport, error) {
	// Run go test with coverage
	coverageData, err := ca.runGoTestCoverage(service)
	if err != nil {
		return nil, fmt.Errorf("failed to run coverage: %w", err)
	}

	// Parse coverage data
	report, err := ca.parseCoverageData(service, coverageData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coverage data: %w", err)
	}

	// Evaluate thresholds
	violations := ca.evaluateThresholds(report)
	report.Violations = violations
	report.Status = ca.determineStatus(violations)

	// Add metadata
	report.Metadata = ca.generateMetadata(service)

	return report, nil
}

// runGoTestCoverage runs go test with coverage collection
func (ca *CoverageAnalyzer) runGoTestCoverage(service string) (string, error) {
	servicePath := filepath.Join(ca.projectRoot, service)

	// Create coverage directory
	coverageDir := filepath.Join(ca.projectRoot, ca.config.OutputDir)
	if err := os.MkdirAll(coverageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create coverage directory: %w", err)
	}

	// Run go test with coverage
	coverageFile := filepath.Join(coverageDir, fmt.Sprintf("%s.out", service))
	cmd := exec.Command("go", "test", "-coverprofile="+coverageFile, "-covermode=atomic", "./...")
	cmd.Dir = servicePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go test failed: %w, output: %s", err, string(output))
	}

	// Read coverage file
	coverageData, err := os.ReadFile(coverageFile)
	if err != nil {
		return "", fmt.Errorf("failed to read coverage file: %w", err)
	}

	return string(coverageData), nil
}

// parseCoverageData parses go coverage data
func (ca *CoverageAnalyzer) parseCoverageData(service string, coverageData string) (*CoverageReport, error) {
	report := &CoverageReport{
		Service:         service,
		Timestamp:       time.Now(),
		PackageCoverage: make(map[string]PackageCoverageInfo),
		Thresholds:      ca.config.Thresholds,
	}

	// Parse coverage lines
	lines := strings.Split(coverageData, "\n")
	packageStats := make(map[string]*packageCoverageStats)

	for _, line := range lines[1:] { // Skip header line
		if strings.TrimSpace(line) == "" {
			continue
		}

		coverage, err := ca.parseCoverageLine(line)
		if err != nil {
			continue // Skip invalid lines
		}

		packageName := ca.extractPackageName(coverage.File)

		if packageStats[packageName] == nil {
			packageStats[packageName] = &packageCoverageStats{}
		}

		packageStats[packageName].addCoverage(coverage)
	}

	// Calculate package coverage
	totalStatements := 0
	totalCovered := 0

	for packageName, stats := range packageStats {
		packageInfo := PackageCoverageInfo{
			Package:           packageName,
			Coverage:          stats.calculateCoverage(),
			TotalStatements:   stats.totalStatements,
			CoveredStatements: stats.coveredStatements,
			UncoveredLines:    stats.uncoveredLines,
			Functions:         stats.functions,
		}

		report.PackageCoverage[packageName] = packageInfo
		totalStatements += stats.totalStatements
		totalCovered += stats.coveredStatements
	}

	// Calculate total coverage
	if totalStatements > 0 {
		report.TotalCoverage = float64(totalCovered) / float64(totalStatements) * 100
	}

	return report, nil
}

// coverageLineData represents parsed coverage line data
type coverageLineData struct {
	File       string
	StartLine  int
	StartCol   int
	EndLine    int
	EndCol     int
	Statements int
	Count      int
}

// packageCoverageStats accumulates coverage statistics for a package
type packageCoverageStats struct {
	totalStatements   int
	coveredStatements int
	uncoveredLines    []int
	functions         []FunctionCoverageInfo
}

// addCoverage adds coverage data to package stats
func (pcs *packageCoverageStats) addCoverage(coverage coverageLineData) {
	pcs.totalStatements += coverage.Statements
	if coverage.Count > 0 {
		pcs.coveredStatements += coverage.Statements
	} else {
		// Add uncovered lines
		for line := coverage.StartLine; line <= coverage.EndLine; line++ {
			pcs.uncoveredLines = append(pcs.uncoveredLines, line)
		}
	}
}

// calculateCoverage calculates coverage percentage for package
func (pcs *packageCoverageStats) calculateCoverage() float64 {
	if pcs.totalStatements == 0 {
		return 0
	}
	return float64(pcs.coveredStatements) / float64(pcs.totalStatements) * 100
}

// parseCoverageLine parses a single coverage line
func (ca *CoverageAnalyzer) parseCoverageLine(line string) (coverageLineData, error) {
	// Coverage line format: file.go:startLine.startCol,endLine.endCol statements count
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return coverageLineData{}, fmt.Errorf("invalid coverage line format")
	}

	// Parse file and position
	filePos := parts[0]
	fileParts := strings.Split(filePos, ":")
	if len(fileParts) != 2 {
		return coverageLineData{}, fmt.Errorf("invalid file position format")
	}

	file := fileParts[0]
	position := fileParts[1]

	// Parse position (startLine.startCol,endLine.endCol)
	posParts := strings.Split(position, ",")
	if len(posParts) != 2 {
		return coverageLineData{}, fmt.Errorf("invalid position format")
	}

	startParts := strings.Split(posParts[0], ".")
	endParts := strings.Split(posParts[1], ".")
	if len(startParts) != 2 || len(endParts) != 2 {
		return coverageLineData{}, fmt.Errorf("invalid line.col format")
	}

	startLine, _ := strconv.Atoi(startParts[0])
	startCol, _ := strconv.Atoi(startParts[1])
	endLine, _ := strconv.Atoi(endParts[0])
	endCol, _ := strconv.Atoi(endParts[1])

	// Parse statements and count
	statements, _ := strconv.Atoi(parts[1])
	count, _ := strconv.Atoi(parts[2])

	return coverageLineData{
		File:       file,
		StartLine:  startLine,
		StartCol:   startCol,
		EndLine:    endLine,
		EndCol:     endCol,
		Statements: statements,
		Count:      count,
	}, nil
}

// extractPackageName extracts package name from file path
func (ca *CoverageAnalyzer) extractPackageName(filePath string) string {
	// Remove project root and extract package name
	relativePath := strings.TrimPrefix(filePath, ca.projectRoot+"/")
	parts := strings.Split(relativePath, "/")

	if len(parts) > 1 {
		return parts[len(parts)-2] // Parent directory name
	}

	return "main"
}

// evaluateThresholds evaluates coverage against configured thresholds
func (ca *CoverageAnalyzer) evaluateThresholds(report *CoverageReport) []CoverageViolation {
	var violations []CoverageViolation

	// Check total coverage threshold
	if report.TotalCoverage < ca.config.Thresholds.TotalCoverage {
		violations = append(violations, CoverageViolation{
			Type:     "total_coverage",
			Current:  report.TotalCoverage,
			Expected: ca.config.Thresholds.TotalCoverage,
			Severity: ca.determineSeverity(report.TotalCoverage, ca.config.Thresholds.TotalCoverage),
			Message:  "Total coverage below threshold",
		})
	}

	// Check service-specific threshold
	if serviceThreshold, exists := ca.config.Thresholds.ServiceThresholds[report.Service]; exists {
		if report.TotalCoverage < serviceThreshold {
			violations = append(violations, CoverageViolation{
				Type:     "service_coverage",
				Package:  report.Service,
				Current:  report.TotalCoverage,
				Expected: serviceThreshold,
				Severity: ca.determineSeverity(report.TotalCoverage, serviceThreshold),
				Message:  fmt.Sprintf("Service %s coverage below service-specific threshold", report.Service),
			})
		}
	}

	// Check package coverage thresholds
	for packageName, packageInfo := range report.PackageCoverage {
		if packageInfo.Coverage < ca.config.Thresholds.PackageCoverage {
			violations = append(violations, CoverageViolation{
				Type:     "package_coverage",
				Package:  packageName,
				Current:  packageInfo.Coverage,
				Expected: ca.config.Thresholds.PackageCoverage,
				Severity: ca.determineSeverity(packageInfo.Coverage, ca.config.Thresholds.PackageCoverage),
				Message:  fmt.Sprintf("Package %s coverage below threshold", packageName),
			})
		}
	}

	return violations
}

// determineSeverity determines violation severity based on gap
func (ca *CoverageAnalyzer) determineSeverity(current, required float64) string {
	gap := required - current

	if gap >= 20 {
		return "critical"
	} else if gap >= 10 {
		return "high"
	} else if gap >= 5 {
		return "medium"
	}

	return "low"
}

// determineStatus determines overall status based on violations
func (ca *CoverageAnalyzer) determineStatus(violations []CoverageViolation) string {
	if len(violations) == 0 {
		return "passed"
	}

	if ca.config.FailOnViolation {
		for _, violation := range violations {
			if violation.Severity == "critical" || violation.Severity == "high" {
				return "failed"
			}
		}
		return "warning"
	}

	return "warning"
}

// generateMetadata generates coverage metadata
func (ca *CoverageAnalyzer) generateMetadata(service string) CoverageMetadata {
	return CoverageMetadata{
		GoVersion:     ca.getGoVersion(),
		Timestamp:     time.Now().UTC(),
		ExecutionTime: "0s", // Will be calculated during actual test run
		TotalPackages: 0,    // Will be calculated during actual test run
		ExcludedFiles: ca.config.ExcludePatterns,
		ReportFormat:  "json",
		Environment:   "test",
		Tags: map[string]string{
			"service": service,
			"tool":    "go-coverage",
		},
	}
}

// getGoVersion gets current Go version
func (ca *CoverageAnalyzer) getGoVersion() string {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	// Parse "go version go1.21.0 ..."
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return strings.TrimPrefix(parts[2], "go")
	}

	return "unknown"
}

// CoverageReporter handles coverage report generation and distribution
type CoverageReporter struct {
	config CoverageConfig
}

// NewCoverageReporter creates a new coverage reporter
func NewCoverageReporter(config CoverageConfig) *CoverageReporter {
	return &CoverageReporter{config: config}
}

// GenerateReports generates coverage reports in all configured formats
func (cr *CoverageReporter) GenerateReports(report *CoverageReport) error {
	for _, format := range cr.config.OutputFormat {
		if err := cr.generateReport(report, format); err != nil {
			return fmt.Errorf("failed to generate %s report: %w", format, err)
		}
	}

	return nil
}

// generateReport generates a single coverage report
func (cr *CoverageReporter) generateReport(report *CoverageReport, format string) error {
	outputPath := filepath.Join(cr.config.OutputDir, fmt.Sprintf("%s-coverage.%s", report.Service, format))

	switch format {
	case "json":
		return cr.generateJSONReport(report, outputPath)
	case "html":
		return cr.generateHTMLReport(report, outputPath)
	case "lcov":
		return cr.generateLCOVReport(report, outputPath)
	case "text":
		return cr.generateTextReport(report, outputPath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// generateJSONReport generates JSON coverage report
func (cr *CoverageReporter) generateJSONReport(report *CoverageReport, outputPath string) error {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}

// generateHTMLReport generates HTML coverage report
func (cr *CoverageReporter) generateHTMLReport(report *CoverageReport, outputPath string) error {
	html := cr.buildHTMLReport(report)
	return os.WriteFile(outputPath, []byte(html), 0644)
}

// generateLCOVReport generates LCOV coverage report
func (cr *CoverageReporter) generateLCOVReport(report *CoverageReport, outputPath string) error {
	lcov := cr.buildLCOVReport(report)
	return os.WriteFile(outputPath, []byte(lcov), 0644)
}

// generateTextReport generates text coverage report
func (cr *CoverageReporter) generateTextReport(report *CoverageReport, outputPath string) error {
	text := cr.buildTextReport(report)
	return os.WriteFile(outputPath, []byte(text), 0644)
}

// buildHTMLReport builds HTML report content
func (cr *CoverageReporter) buildHTMLReport(report *CoverageReport) string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Coverage Report - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f9f9f9; }
        .container { max-width: 1200px; margin: 0 auto; background-color: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .metric { background-color: #f8f9fa; padding: 15px; border-radius: 8px; text-align: center; border-left: 4px solid #007bff; }
        .metric-value { font-size: 2em; font-weight: bold; color: #333; }
        .metric-label { color: #666; margin-top: 5px; }
        .packages { margin-bottom: 30px; }
        .package { background-color: #f8f9fa; margin: 10px 0; padding: 15px; border-radius: 8px; border-left: 4px solid #28a745; }
        .package-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
        .package-name { font-weight: bold; font-size: 1.1em; }
        .coverage-bar { background-color: #e9ecef; height: 20px; border-radius: 10px; overflow: hidden; }
        .coverage-fill { height: 100%%; background: linear-gradient(90deg, #28a745, #20c997); }
        .violations { margin-top: 20px; }
        .violation { background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; margin: 10px 0; border-radius: 4px; border-left: 4px solid #f39c12; }
        .violation.high { background-color: #f8d7da; border-color: #f5c6cb; border-left-color: #dc3545; }
        .violation.critical { background-color: #f8d7da; border-color: #f5c6cb; border-left-color: #721c24; }
        .status-passed { color: #28a745; }
        .status-warning { color: #ffc107; }
        .status-failed { color: #dc3545; }
        .timestamp { color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Coverage Report</h1>
            <h2>Service: %s</h2>
            <p class="timestamp">Generated: %s</p>
            <p class="status status-%s">Status: %s</p>
        </div>

        <div class="summary">
            <div class="metric">
                <div class="metric-value">%.1f%%</div>
                <div class="metric-label">Total Coverage</div>
            </div>
            <div class="metric">
                <div class="metric-value">%d</div>
                <div class="metric-label">Packages</div>
            </div>
            <div class="metric">
                <div class="metric-value">%d</div>
                <div class="metric-label">Violations</div>
            </div>
        </div>`,
		report.Service, report.Service, report.Timestamp.Format("2006-01-02 15:04:05"),
		report.Status, strings.Title(report.Status), report.TotalCoverage,
		len(report.PackageCoverage), len(report.Violations))

	// Add package details
	html += `<div class="packages"><h3>Package Coverage</h3>`
	for packageName, packageInfo := range report.PackageCoverage {
		coverageClass := "low"
		if packageInfo.Coverage >= 80 {
			coverageClass = "high"
		} else if packageInfo.Coverage >= 60 {
			coverageClass = "medium"
		}

		html += fmt.Sprintf(`
        <div class="package %s">
            <div class="package-header">
                <span class="package-name">%s</span>
                <span class="coverage-percent">%.1f%%</span>
            </div>
            <div class="coverage-bar">
                <div class="coverage-fill" style="width: %.1f%%"></div>
            </div>
            <p>%d/%d statements covered</p>
        </div>`, coverageClass, packageName, packageInfo.Coverage, packageInfo.Coverage,
			packageInfo.CoveredStatements, packageInfo.TotalStatements)
	}
	html += `</div>`

	// Add violations
	if len(report.Violations) > 0 {
		html += `<div class="violations"><h3>Threshold Violations</h3>`
		for _, violation := range report.Violations {
			html += fmt.Sprintf(`
            <div class="violation %s">
                <strong>%s</strong>: %s<br>
                Current: %.1f%% | Expected: %.1f%% | Gap: %.1f%%
            </div>`, violation.Severity, strings.Title(violation.Severity), violation.Message,
				violation.Current, violation.Expected, violation.Expected-violation.Current)
		}
		html += `</div>`
	}

	html += `</div></body></html>`
	return html
}

// buildLCOVReport builds LCOV report content
func (cr *CoverageReporter) buildLCOVReport(report *CoverageReport) string {
	lcov := fmt.Sprintf("TN:%s\n", report.Service)

	for _, packageInfo := range report.PackageCoverage {
		for _, function := range packageInfo.Functions {
			lcov += fmt.Sprintf("SF:%s\n", function.File)
			lcov += fmt.Sprintf("FN:%d,%s\n", function.StartLine, function.Name)
			lcov += fmt.Sprintf("FNDA:%d,%s\n", function.CoveredStatements, function.Name)
			lcov += fmt.Sprintf("FNF:%d\n", function.TotalStatements)
			lcov += fmt.Sprintf("FNH:%d\n", function.CoveredStatements)

			// Add line coverage
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

// buildTextReport builds text report content
func (cr *CoverageReporter) buildTextReport(report *CoverageReport) string {
	text := fmt.Sprintf(`
Tchat Coverage Report
====================

Service: %s
Status: %s
Generated: %s
Total Coverage: %.2f%%

Thresholds:
  Total Coverage: %.1f%%
  Package Coverage: %.1f%%

Package Coverage:
`, report.Service, strings.ToUpper(report.Status), report.Timestamp.Format("2006-01-02 15:04:05"),
		report.TotalCoverage, report.Thresholds.TotalCoverage, report.Thresholds.PackageCoverage)

	for packageName, packageInfo := range report.PackageCoverage {
		status := "✓"
		if packageInfo.Coverage < report.Thresholds.PackageCoverage {
			status = "✗"
		}

		text += fmt.Sprintf("  %s %-20s %.2f%% (%d/%d statements)\n",
			status, packageName, packageInfo.Coverage,
			packageInfo.CoveredStatements, packageInfo.TotalStatements)
	}

	if len(report.Violations) > 0 {
		text += "\nThreshold Violations:\n"
		for _, violation := range report.Violations {
			text += fmt.Sprintf("  [%s] %s: %.2f%% (expected: %.2f%%)\n",
				strings.ToUpper(violation.Severity), violation.Message,
				violation.Current, violation.Expected)
		}
	}

	text += fmt.Sprintf(`
Summary:
  Total Packages: %d
  Violations: %d
  Status: %s
`, len(report.PackageCoverage), len(report.Violations), strings.ToUpper(report.Status))

	return text
}

// CoverageIntegrator handles external integrations
type CoverageIntegrator struct {
	config IntegrationConfig
}

// NewCoverageIntegrator creates a new coverage integrator
func NewCoverageIntegrator(config IntegrationConfig) *CoverageIntegrator {
	return &CoverageIntegrator{config: config}
}

// PublishReports publishes coverage reports to configured integrations
func (ci *CoverageIntegrator) PublishReports(report *CoverageReport) error {
	var errors []string

	// Publish to Codecov
	if ci.config.Codecov.Enabled {
		if err := ci.publishToCodecov(report); err != nil {
			errors = append(errors, fmt.Sprintf("Codecov: %v", err))
		}
	}

	// Publish to Coveralls
	if ci.config.Coveralls.Enabled {
		if err := ci.publishToCoveralls(report); err != nil {
			errors = append(errors, fmt.Sprintf("Coveralls: %v", err))
		}
	}

	// Publish to SonarQube
	if ci.config.SonarQube.Enabled {
		if err := ci.publishToSonarQube(report); err != nil {
			errors = append(errors, fmt.Sprintf("SonarQube: %v", err))
		}
	}

	// Publish to custom hooks
	for _, hook := range ci.config.CustomHooks {
		if err := ci.publishToCustomHook(report, hook); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", hook.Name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("integration errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// publishToCodecov publishes to Codecov (mock implementation for testing)
func (ci *CoverageIntegrator) publishToCodecov(report *CoverageReport) error {
	// Mock implementation - in real scenario would use Codecov API
	fmt.Printf("Publishing to Codecov: %s (%.2f%%)\n", report.Service, report.TotalCoverage)
	return nil
}

// publishToCoveralls publishes to Coveralls (mock implementation for testing)
func (ci *CoverageIntegrator) publishToCoveralls(report *CoverageReport) error {
	// Mock implementation - in real scenario would use Coveralls API
	fmt.Printf("Publishing to Coveralls: %s (%.2f%%)\n", report.Service, report.TotalCoverage)
	return nil
}

// publishToSonarQube publishes to SonarQube (mock implementation for testing)
func (ci *CoverageIntegrator) publishToSonarQube(report *CoverageReport) error {
	// Mock implementation - in real scenario would use SonarQube API
	fmt.Printf("Publishing to SonarQube: %s (%.2f%%)\n", report.Service, report.TotalCoverage)
	return nil
}

// publishToCustomHook publishes to custom webhook (mock implementation for testing)
func (ci *CoverageIntegrator) publishToCustomHook(report *CoverageReport, hook CustomHookConfig) error {
	// Mock implementation - in real scenario would make HTTP request
	fmt.Printf("Publishing to %s: %s (%.2f%%)\n", hook.Name, report.Service, report.TotalCoverage)
	return nil
}

// CoverageConfigManager manages coverage configuration
type CoverageConfigManager struct {
	configPath string
}

// NewCoverageConfigManager creates a new config manager
func NewCoverageConfigManager(configPath string) *CoverageConfigManager {
	return &CoverageConfigManager{configPath: configPath}
}

// LoadConfig loads coverage configuration from file
func (ccm *CoverageConfigManager) LoadConfig() (*CoverageConfig, error) {
	if _, err := os.Stat(ccm.configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return ccm.getDefaultConfig(), nil
	}

	data, err := os.ReadFile(ccm.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config CoverageConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves coverage configuration to file
func (ccm *CoverageConfigManager) SaveConfig(config *CoverageConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(ccm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getDefaultConfig returns default coverage configuration
func (ccm *CoverageConfigManager) getDefaultConfig() *CoverageConfig {
	return &CoverageConfig{
		Enabled:      true,
		OutputFormat: []string{"html", "json"},
		OutputDir:    "coverage",
		Thresholds: CoverageThresholds{
			TotalCoverage:    80.0,
			PackageCoverage:  70.0,
			FunctionCoverage: 60.0,
			ServiceThresholds: map[string]float64{
				"auth":         90.0,
				"payment":      95.0,
				"content":      85.0,
				"messaging":    80.0,
				"notification": 75.0,
				"commerce":     85.0,
			},
		},
		ExcludePatterns: []string{
			"*_test.go",
			"*/mocks/*",
			"*/vendor/*",
			"main.go",
			"*/testutils/*",
		},
		IncludePatterns: []string{
			"*/handlers/*",
			"*/services/*",
			"*/repositories/*",
			"*/middleware/*",
		},
		FailOnViolation: true,
		Integration: IntegrationConfig{
			Codecov: CodecovConfig{
				Enabled:     false,
				ProjectName: "tchat.dev",
			},
			Coveralls: CoverallsConfig{
				Enabled: false,
			},
			SonarQube: SonarQubeConfig{
				Enabled: false,
			},
		},
	}
}
