package contract

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// CrossPlatformCompatibilityTestSuite validates contract compatibility between web and mobile platforms
type CrossPlatformCompatibilityTestSuite struct {
	suite.Suite
	webContract    PactContract
	mobileContract PactContract
}

// PactContract represents a Pact contract structure
type PactContract struct {
	Consumer    ConsumerInfo    `json:"consumer"`
	Provider    ProviderInfo    `json:"provider"`
	Interactions []Interaction  `json:"interactions"`
	Metadata    Metadata       `json:"metadata"`
}

// ConsumerInfo represents consumer information
type ConsumerInfo struct {
	Name string `json:"name"`
}

// ProviderInfo represents provider information
type ProviderInfo struct {
	Name string `json:"name"`
}

// Interaction represents a contract interaction
type Interaction struct {
	Description   string                 `json:"description"`
	ProviderState string                 `json:"providerState"`
	Request       Request                `json:"request"`
	Response      Response               `json:"response"`
}

// Request represents HTTP request specification
type Request struct {
	Method  string                 `json:"method"`
	Path    string                 `json:"path"`
	Query   map[string]string      `json:"query,omitempty"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body,omitempty"`
}

// Response represents HTTP response specification
type Response struct {
	Status        int                    `json:"status"`
	Headers       map[string]string      `json:"headers"`
	Body          map[string]interface{} `json:"body"`
	MatchingRules map[string]MatchRule   `json:"matchingRules,omitempty"`
}

// MatchRule represents matching rule for response validation
type MatchRule struct {
	Match string `json:"match"`
	Regex string `json:"regex,omitempty"`
	Min   int    `json:"min,omitempty"`
}

// Metadata represents contract metadata
type Metadata struct {
	PactSpecification PactSpecInfo           `json:"pactSpecification"`
	Client            ClientInfo             `json:"client"`
	Platform          map[string]string      `json:"platform,omitempty"`
}

// PactSpecInfo represents Pact specification version
type PactSpecInfo struct {
	Version string `json:"version"`
}

// ClientInfo represents client library information
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// CompatibilityResult represents validation result
type CompatibilityResult struct {
	Score         float64                      `json:"score"`
	Issues        []CompatibilityIssue         `json:"issues"`
	Summary       CompatibilitySummary         `json:"summary"`
	Recommendations []string                   `json:"recommendations"`
	Timestamp     time.Time                    `json:"timestamp"`
}

// CompatibilityIssue represents a specific compatibility issue
type CompatibilityIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	WebValue    string `json:"web_value,omitempty"`
	MobileValue string `json:"mobile_value,omitempty"`
	Path        string `json:"path,omitempty"`
}

// CompatibilitySummary provides overall compatibility statistics
type CompatibilitySummary struct {
	TotalInteractions    int `json:"total_interactions"`
	CompatibleEndpoints  int `json:"compatible_endpoints"`
	CriticalIssues      int `json:"critical_issues"`
	MajorIssues         int `json:"major_issues"`
	MinorIssues         int `json:"minor_issues"`
}

// SetupSuite loads contract files before running tests
func (suite *CrossPlatformCompatibilityTestSuite) SetupSuite() {
	// Load web contract
	webContractPath := filepath.Join("..", "..", "specs", "021-implement-pact-contract", "contracts", "pact-consumer-web.json")
	webData, err := ioutil.ReadFile(webContractPath)
	require.NoError(suite.T(), err, "Failed to load web contract file")

	err = json.Unmarshal(webData, &suite.webContract)
	require.NoError(suite.T(), err, "Failed to parse web contract JSON")

	// Load mobile contract
	mobileContractPath := filepath.Join("..", "..", "specs", "021-implement-pact-contract", "contracts", "pact-consumer-mobile.json")
	mobileData, err := ioutil.ReadFile(mobileContractPath)
	require.NoError(suite.T(), err, "Failed to load mobile contract file")

	err = json.Unmarshal(mobileData, &suite.mobileContract)
	require.NoError(suite.T(), err, "Failed to parse mobile contract JSON")
}

// TestContractStructureCompatibility validates basic contract structure consistency
func (suite *CrossPlatformCompatibilityTestSuite) TestContractStructureCompatibility() {
	suite.T().Log("Validating basic contract structure compatibility")

	// Validate Pact specification versions
	webVersion := suite.webContract.Metadata.PactSpecification.Version
	mobileVersion := suite.mobileContract.Metadata.PactSpecification.Version

	assert.Equal(suite.T(), webVersion, mobileVersion,
		"Pact specification versions must match between platforms")

	// Validate both contracts have interactions
	assert.Greater(suite.T(), len(suite.webContract.Interactions), 0,
		"Web contract must have interactions")
	assert.Greater(suite.T(), len(suite.mobileContract.Interactions), 0,
		"Mobile contract must have interactions")

	suite.T().Logf("Structure validation passed - Web: %d interactions, Mobile: %d interactions",
		len(suite.webContract.Interactions), len(suite.mobileContract.Interactions))
}

// TestAuthenticationCompatibility validates authentication patterns between platforms
func (suite *CrossPlatformCompatibilityTestSuite) TestAuthenticationCompatibility() {
	suite.T().Log("Validating authentication compatibility across platforms")

	issues := []CompatibilityIssue{}

	// Check web auth patterns
	webAuthPatterns := suite.extractAuthenticationPatterns(suite.webContract)
	mobileAuthPatterns := suite.extractAuthenticationPatterns(suite.mobileContract)

	// Validate token format consistency
	if webAuthPatterns.TokenFormat != "" && mobileAuthPatterns.TokenFormat != "" {
		if webAuthPatterns.TokenFormat != mobileAuthPatterns.TokenFormat {
			issues = append(issues, CompatibilityIssue{
				Type:        "authentication",
				Severity:    "critical",
				Description: "Token format mismatch between platforms",
				WebValue:    webAuthPatterns.TokenFormat,
				MobileValue: mobileAuthPatterns.TokenFormat,
				Path:        "$.headers.Authorization",
			})
		}
	}

	// Validate auth header patterns
	if webAuthPatterns.HeaderPattern != mobileAuthPatterns.HeaderPattern {
		issues = append(issues, CompatibilityIssue{
			Type:        "authentication",
			Severity:    "major",
			Description: "Authorization header pattern inconsistency",
			WebValue:    webAuthPatterns.HeaderPattern,
			MobileValue: mobileAuthPatterns.HeaderPattern,
			Path:        "$.headers.Authorization",
		})
	}

	// MUST FAIL if critical authentication issues exist
	for _, issue := range issues {
		if issue.Severity == "critical" {
			suite.T().Errorf("CRITICAL AUTH INCOMPATIBILITY: %s", issue.Description)
		}
	}

	assert.Empty(suite.T(), issues, "Authentication compatibility issues detected")
	suite.T().Log("Authentication compatibility validation passed")
}

// TestResponseSchemaCompatibility validates response schema consistency
func (suite *CrossPlatformCompatibilityTestSuite) TestResponseSchemaCompatibility() {
	suite.T().Log("Validating response schema compatibility")

	incompatibilities := []CompatibilityIssue{}

	// Analyze response schemas for common endpoint patterns
	webResponseSchemas := suite.extractResponseSchemas(suite.webContract)
	mobileResponseSchemas := suite.extractResponseSchemas(suite.mobileContract)

	// Check for schema field type mismatches
	for endpoint, webSchema := range webResponseSchemas {
		if mobileSchema, exists := mobileResponseSchemas[endpoint]; exists {
			issues := suite.compareSchemaStructures(webSchema, mobileSchema, endpoint)
			incompatibilities = append(incompatibilities, issues...)
		}
	}

	// MUST FAIL if critical schema incompatibilities exist
	criticalIssues := 0
	for _, issue := range incompatibilities {
		if issue.Severity == "critical" {
			criticalIssues++
			suite.T().Errorf("CRITICAL SCHEMA INCOMPATIBILITY: %s at %s",
				issue.Description, issue.Path)
		}
	}

	assert.Zero(suite.T(), criticalIssues,
		"Critical response schema incompatibilities detected")

	suite.T().Logf("Response schema validation completed with %d total issues",
		len(incompatibilities))
}

// TestDataTypeConsistency validates data type consistency across platforms
func (suite *CrossPlatformCompatibilityTestSuite) TestDataTypeConsistency() {
	suite.T().Log("Validating data type consistency")

	typeInconsistencies := []CompatibilityIssue{}

	// Extract and compare data types from both contracts
	webTypes := suite.extractDataTypes(suite.webContract)
	mobileTypes := suite.extractDataTypes(suite.mobileContract)

	// Compare common field types
	for fieldPath, webType := range webTypes {
		if mobileType, exists := mobileTypes[fieldPath]; exists {
			if !suite.areTypesCompatible(webType, mobileType) {
				typeInconsistencies = append(typeInconsistencies, CompatibilityIssue{
					Type:        "data_type",
					Severity:    "major",
					Description: fmt.Sprintf("Data type mismatch for field '%s'", fieldPath),
					WebValue:    webType,
					MobileValue: mobileType,
					Path:        fieldPath,
				})
			}
		}
	}

	// MUST FAIL if major type inconsistencies exist
	majorTypeIssues := 0
	for _, issue := range typeInconsistencies {
		if issue.Severity == "major" {
			majorTypeIssues++
			suite.T().Errorf("MAJOR TYPE INCONSISTENCY: %s", issue.Description)
		}
	}

	assert.Zero(suite.T(), majorTypeIssues,
		"Major data type inconsistencies detected")

	suite.T().Log("Data type consistency validation passed")
}

// TestErrorHandlingConsistency validates error response consistency
func (suite *CrossPlatformCompatibilityTestSuite) TestErrorHandlingConsistency() {
	suite.T().Log("Validating error handling consistency")

	webErrorPatterns := suite.extractErrorPatterns(suite.webContract)
	mobileErrorPatterns := suite.extractErrorPatterns(suite.mobileContract)

	inconsistencies := []CompatibilityIssue{}

	// Compare error status codes
	if len(webErrorPatterns.StatusCodes) > 0 && len(mobileErrorPatterns.StatusCodes) > 0 {
		for statusCode := range webErrorPatterns.StatusCodes {
			if _, exists := mobileErrorPatterns.StatusCodes[statusCode]; !exists {
				inconsistencies = append(inconsistencies, CompatibilityIssue{
					Type:        "error_handling",
					Severity:    "minor",
					Description: fmt.Sprintf("Status code %d only handled in web platform", statusCode),
					Path:        "$.response.status",
				})
			}
		}
	}

	// Compare error response structures
	if webErrorPatterns.ErrorStructure != nil && mobileErrorPatterns.ErrorStructure != nil {
		structuralDiffs := suite.compareErrorStructures(
			webErrorPatterns.ErrorStructure,
			mobileErrorPatterns.ErrorStructure)
		inconsistencies = append(inconsistencies, structuralDiffs...)
	}

	assert.Less(suite.T(), len(inconsistencies), 3,
		"Too many error handling inconsistencies detected")

	suite.T().Logf("Error handling validation completed with %d inconsistencies",
		len(inconsistencies))
}

// TestVersionCompatibility validates API versioning consistency
func (suite *CrossPlatformCompatibilityTestSuite) TestVersionCompatibility() {
	suite.T().Log("Validating API versioning compatibility")

	webVersions := suite.extractAPIVersions(suite.webContract)
	mobileVersions := suite.extractAPIVersions(suite.mobileContract)

	// Check version consistency
	versionIssues := []string{}

	for endpoint, webVersion := range webVersions {
		if mobileVersion, exists := mobileVersions[endpoint]; exists {
			if webVersion != mobileVersion {
				versionIssues = append(versionIssues,
					fmt.Sprintf("Version mismatch for %s: web=%s, mobile=%s",
						endpoint, webVersion, mobileVersion))
			}
		}
	}

	// MUST FAIL if version inconsistencies exist
	assert.Empty(suite.T(), versionIssues, "API version inconsistencies detected")
	suite.T().Log("API versioning compatibility validated")
}

// TestComprehensiveCompatibilityScore generates overall compatibility assessment
func (suite *CrossPlatformCompatibilityTestSuite) TestComprehensiveCompatibilityScore() {
	suite.T().Log("Generating comprehensive compatibility assessment")

	result := suite.generateCompatibilityReport()

	// Log detailed compatibility report
	suite.T().Logf("Compatibility Assessment:")
	suite.T().Logf("- Overall Score: %.2f%%", result.Score)
	suite.T().Logf("- Total Interactions: %d", result.Summary.TotalInteractions)
	suite.T().Logf("- Compatible Endpoints: %d", result.Summary.CompatibleEndpoints)
	suite.T().Logf("- Critical Issues: %d", result.Summary.CriticalIssues)
	suite.T().Logf("- Major Issues: %d", result.Summary.MajorIssues)
	suite.T().Logf("- Minor Issues: %d", result.Summary.MinorIssues)

	// Print recommendations
	if len(result.Recommendations) > 0 {
		suite.T().Log("Recommendations:")
		for i, rec := range result.Recommendations {
			suite.T().Logf("%d. %s", i+1, rec)
		}
	}

	// MUST FAIL if compatibility score is below acceptable threshold
	minimumCompatibilityScore := 85.0
	assert.GreaterOrEqual(suite.T(), result.Score, minimumCompatibilityScore,
		"Cross-platform compatibility score below acceptable threshold")

	// MUST FAIL if critical issues exist
	assert.Zero(suite.T(), result.Summary.CriticalIssues,
		"Critical compatibility issues must be resolved")

	// Save detailed report for CI/CD integration
	suite.saveCompatibilityReport(result)

	suite.T().Log("Comprehensive compatibility assessment completed")
}

// Helper method implementations

// AuthenticationPattern represents authentication patterns
type AuthenticationPattern struct {
	TokenFormat    string
	HeaderPattern  string
	TokenPrefix    string
}

func (suite *CrossPlatformCompatibilityTestSuite) extractAuthenticationPatterns(contract PactContract) AuthenticationPattern {
	pattern := AuthenticationPattern{}

	for _, interaction := range contract.Interactions {
		if authHeader, exists := interaction.Request.Headers["Authorization"]; exists {
			// Extract token prefix (Bearer, Basic, etc.)
			parts := strings.Split(authHeader, " ")
			if len(parts) >= 1 {
				pattern.TokenPrefix = parts[0]
				pattern.HeaderPattern = authHeader
			}

			// Extract token format from response if login interaction
			if strings.Contains(strings.ToLower(interaction.Description), "login") {
				if token, exists := interaction.Response.Body["access_token"]; exists {
					pattern.TokenFormat = fmt.Sprintf("%T", token)
				}
			}
		}
	}

	return pattern
}

func (suite *CrossPlatformCompatibilityTestSuite) extractResponseSchemas(contract PactContract) map[string]map[string]interface{} {
	schemas := make(map[string]map[string]interface{})

	for _, interaction := range contract.Interactions {
		endpoint := fmt.Sprintf("%s %s", interaction.Request.Method, interaction.Request.Path)
		schemas[endpoint] = interaction.Response.Body
	}

	return schemas
}

func (suite *CrossPlatformCompatibilityTestSuite) compareSchemaStructures(webSchema, mobileSchema map[string]interface{}, endpoint string) []CompatibilityIssue {
	issues := []CompatibilityIssue{}

	// Compare top-level fields
	for field, webValue := range webSchema {
		if mobileValue, exists := mobileSchema[field]; exists {
			webType := reflect.TypeOf(webValue)
			mobileType := reflect.TypeOf(mobileValue)

			if webType != mobileType {
				issues = append(issues, CompatibilityIssue{
					Type:        "schema_structure",
					Severity:    "critical",
					Description: fmt.Sprintf("Field type mismatch for '%s'", field),
					WebValue:    webType.String(),
					MobileValue: mobileType.String(),
					Path:        fmt.Sprintf("$.%s.body.%s", endpoint, field),
				})
			}
		} else {
			issues = append(issues, CompatibilityIssue{
				Type:        "schema_structure",
				Severity:    "major",
				Description: fmt.Sprintf("Field '%s' missing in mobile schema", field),
				Path:        fmt.Sprintf("$.%s.body.%s", endpoint, field),
			})
		}
	}

	return issues
}

func (suite *CrossPlatformCompatibilityTestSuite) extractDataTypes(contract PactContract) map[string]string {
	types := make(map[string]string)

	for _, interaction := range contract.Interactions {
		suite.extractTypesFromMap(interaction.Response.Body, "", types)
	}

	return types
}

func (suite *CrossPlatformCompatibilityTestSuite) extractTypesFromMap(data map[string]interface{}, prefix string, types map[string]string) {
	for key, value := range data {
		fieldPath := key
		if prefix != "" {
			fieldPath = prefix + "." + key
		}

		types[fieldPath] = reflect.TypeOf(value).String()

		// Recursively extract from nested maps
		if nested, ok := value.(map[string]interface{}); ok {
			suite.extractTypesFromMap(nested, fieldPath, types)
		}
	}
}

func (suite *CrossPlatformCompatibilityTestSuite) areTypesCompatible(webType, mobileType string) bool {
	// Define type compatibility rules
	compatibleTypes := map[string][]string{
		"string": {"string"},
		"float64": {"float64", "int"},
		"int": {"int", "float64"},
		"bool": {"bool"},
		"[]interface {}": {"[]interface {}"},
		"map[string]interface {}": {"map[string]interface {}"},
	}

	if compatible, exists := compatibleTypes[webType]; exists {
		for _, compatType := range compatible {
			if compatType == mobileType {
				return true
			}
		}
	}

	return webType == mobileType
}

// ErrorPattern represents error handling patterns
type ErrorPattern struct {
	StatusCodes    map[int]bool
	ErrorStructure map[string]interface{}
}

func (suite *CrossPlatformCompatibilityTestSuite) extractErrorPatterns(contract PactContract) ErrorPattern {
	pattern := ErrorPattern{
		StatusCodes: make(map[int]bool),
	}

	for _, interaction := range contract.Interactions {
		if interaction.Response.Status >= 400 {
			pattern.StatusCodes[interaction.Response.Status] = true
			pattern.ErrorStructure = interaction.Response.Body
		}
	}

	return pattern
}

func (suite *CrossPlatformCompatibilityTestSuite) compareErrorStructures(webStructure, mobileStructure map[string]interface{}) []CompatibilityIssue {
	issues := []CompatibilityIssue{}

	// Compare error response fields
	for field := range webStructure {
		if _, exists := mobileStructure[field]; !exists {
			issues = append(issues, CompatibilityIssue{
				Type:        "error_structure",
				Severity:    "minor",
				Description: fmt.Sprintf("Error field '%s' missing in mobile responses", field),
				Path:        fmt.Sprintf("$.error.%s", field),
			})
		}
	}

	return issues
}

func (suite *CrossPlatformCompatibilityTestSuite) extractAPIVersions(contract PactContract) map[string]string {
	versions := make(map[string]string)

	for _, interaction := range contract.Interactions {
		// Extract version from path (e.g., /api/v1/auth/login -> v1)
		pathParts := strings.Split(interaction.Request.Path, "/")
		for _, part := range pathParts {
			if strings.HasPrefix(part, "v") && len(part) > 1 {
				endpoint := fmt.Sprintf("%s %s", interaction.Request.Method, interaction.Request.Path)
				versions[endpoint] = part
				break
			}
		}
	}

	return versions
}

func (suite *CrossPlatformCompatibilityTestSuite) generateCompatibilityReport() CompatibilityResult {
	result := CompatibilityResult{
		Issues:      []CompatibilityIssue{},
		Timestamp:   time.Now(),
	}

	// Calculate compatibility metrics
	totalInteractions := len(suite.webContract.Interactions) + len(suite.mobileContract.Interactions)
	result.Summary.TotalInteractions = totalInteractions

	// Run all compatibility checks and aggregate issues
	issues := []CompatibilityIssue{}

	// Collect issues from different validation areas
	authIssues := suite.validateAuthenticationCompatibility()
	schemaIssues := suite.validateSchemaCompatibility()
	typeIssues := suite.validateTypeCompatibility()
	errorIssues := suite.validateErrorCompatibility()

	issues = append(issues, authIssues...)
	issues = append(issues, schemaIssues...)
	issues = append(issues, typeIssues...)
	issues = append(issues, errorIssues...)

	result.Issues = issues

	// Count issues by severity
	for _, issue := range issues {
		switch issue.Severity {
		case "critical":
			result.Summary.CriticalIssues++
		case "major":
			result.Summary.MajorIssues++
		case "minor":
			result.Summary.MinorIssues++
		}
	}

	// Calculate compatibility score
	totalIssues := result.Summary.CriticalIssues + result.Summary.MajorIssues + result.Summary.MinorIssues
	maxPossibleIssues := totalInteractions * 5 // Assume 5 potential issues per interaction

	if maxPossibleIssues > 0 {
		issueWeight := float64(result.Summary.CriticalIssues*3 + result.Summary.MajorIssues*2 + result.Summary.MinorIssues*1)
		result.Score = math.Max(0, 100.0 - (issueWeight / float64(maxPossibleIssues) * 100.0))
	} else {
		result.Score = 100.0
	}

	// Calculate compatible endpoints
	result.Summary.CompatibleEndpoints = totalInteractions - totalIssues
	if result.Summary.CompatibleEndpoints < 0 {
		result.Summary.CompatibleEndpoints = 0
	}

	// Generate recommendations
	result.Recommendations = suite.generateRecommendations(result)

	return result
}

// Additional validation methods for comprehensive report
func (suite *CrossPlatformCompatibilityTestSuite) validateAuthenticationCompatibility() []CompatibilityIssue {
	issues := []CompatibilityIssue{}

	// Check web auth patterns
	webAuthPatterns := suite.extractAuthenticationPatterns(suite.webContract)
	mobileAuthPatterns := suite.extractAuthenticationPatterns(suite.mobileContract)

	// Validate token format consistency
	if webAuthPatterns.TokenFormat != "" && mobileAuthPatterns.TokenFormat != "" {
		if webAuthPatterns.TokenFormat != mobileAuthPatterns.TokenFormat {
			issues = append(issues, CompatibilityIssue{
				Type:        "authentication",
				Severity:    "critical",
				Description: "Token format mismatch between platforms",
				WebValue:    webAuthPatterns.TokenFormat,
				MobileValue: mobileAuthPatterns.TokenFormat,
				Path:        "$.headers.Authorization",
			})
		}
	}

	// Validate auth header patterns
	if webAuthPatterns.HeaderPattern != mobileAuthPatterns.HeaderPattern {
		issues = append(issues, CompatibilityIssue{
			Type:        "authentication",
			Severity:    "major",
			Description: "Authorization header pattern inconsistency",
			WebValue:    webAuthPatterns.HeaderPattern,
			MobileValue: mobileAuthPatterns.HeaderPattern,
			Path:        "$.headers.Authorization",
		})
	}

	return issues
}

func (suite *CrossPlatformCompatibilityTestSuite) validateSchemaCompatibility() []CompatibilityIssue {
	issues := []CompatibilityIssue{}

	// Analyze response schemas for common endpoint patterns
	webResponseSchemas := suite.extractResponseSchemas(suite.webContract)
	mobileResponseSchemas := suite.extractResponseSchemas(suite.mobileContract)

	// Check for schema field type mismatches across services
	// Even though they're different services, check for common field structures
	commonFields := map[string]bool{
		"id": true, "status": true, "created_at": true, "updated_at": true,
	}

	for webEndpoint, webSchema := range webResponseSchemas {
		for mobileEndpoint, mobileSchema := range mobileResponseSchemas {
			// Compare common fields between different services
			for field := range commonFields {
				if webField, webExists := webSchema[field]; webExists {
					if mobileField, mobileExists := mobileSchema[field]; mobileExists {
						webType := reflect.TypeOf(webField)
						mobileType := reflect.TypeOf(mobileField)

						if webType != mobileType {
							issues = append(issues, CompatibilityIssue{
								Type:        "schema_structure",
								Severity:    "major",
								Description: fmt.Sprintf("Common field '%s' has different types", field),
								WebValue:    fmt.Sprintf("%s: %s", webEndpoint, webType.String()),
								MobileValue: fmt.Sprintf("%s: %s", mobileEndpoint, mobileType.String()),
								Path:        fmt.Sprintf("$.%s", field),
							})
						}
					}
				}
			}
		}
	}

	return issues
}

func (suite *CrossPlatformCompatibilityTestSuite) validateTypeCompatibility() []CompatibilityIssue {
	issues := []CompatibilityIssue{}

	// Extract and compare data types from both contracts
	webTypes := suite.extractDataTypes(suite.webContract)
	mobileTypes := suite.extractDataTypes(suite.mobileContract)

	// Compare common field types across different services
	commonFieldPatterns := []string{"id", "status", "created_at", "updated_at"}

	for _, pattern := range commonFieldPatterns {
		var webType, mobileType string
		var webFound, mobileFound bool

		// Find fields matching pattern in web contract
		for fieldPath, fieldType := range webTypes {
			if strings.Contains(fieldPath, pattern) {
				webType = fieldType
				webFound = true
				break
			}
		}

		// Find fields matching pattern in mobile contract
		for fieldPath, fieldType := range mobileTypes {
			if strings.Contains(fieldPath, pattern) {
				mobileType = fieldType
				mobileFound = true
				break
			}
		}

		// Compare if both found
		if webFound && mobileFound && webType != mobileType {
			if !suite.areTypesCompatible(webType, mobileType) {
				issues = append(issues, CompatibilityIssue{
					Type:        "data_type",
					Severity:    "major",
					Description: fmt.Sprintf("Data type mismatch for common field pattern '%s'", pattern),
					WebValue:    webType,
					MobileValue: mobileType,
					Path:        pattern,
				})
			}
		}
	}

	return issues
}

func (suite *CrossPlatformCompatibilityTestSuite) validateErrorCompatibility() []CompatibilityIssue {
	issues := []CompatibilityIssue{}

	// Since contracts don't have explicit error cases, we check response status patterns
	// Look for consistency in success status codes (200, 201, etc.)
	webStatusCodes := make(map[int]bool)
	mobileStatusCodes := make(map[int]bool)

	for _, interaction := range suite.webContract.Interactions {
		webStatusCodes[interaction.Response.Status] = true
	}

	for _, interaction := range suite.mobileContract.Interactions {
		mobileStatusCodes[interaction.Response.Status] = true
	}

	// Check if platforms use different status codes for similar operations
	if len(webStatusCodes) > 0 && len(mobileStatusCodes) > 0 {
		hasCreate := false
		webCreateStatus := 0
		mobileCreateStatus := 0

		// Check for create operation status codes
		for _, interaction := range suite.webContract.Interactions {
			if interaction.Request.Method == "POST" {
				hasCreate = true
				webCreateStatus = interaction.Response.Status
				break
			}
		}

		for _, interaction := range suite.mobileContract.Interactions {
			if interaction.Request.Method == "POST" {
				mobileCreateStatus = interaction.Response.Status
				break
			}
		}

		if hasCreate && webCreateStatus != 0 && mobileCreateStatus != 0 &&
		   webCreateStatus != mobileCreateStatus {
			issues = append(issues, CompatibilityIssue{
				Type:        "error_handling",
				Severity:    "minor",
				Description: "Different status codes for create operations",
				WebValue:    fmt.Sprintf("%d", webCreateStatus),
				MobileValue: fmt.Sprintf("%d", mobileCreateStatus),
				Path:        "$.response.status",
			})
		}
	}

	return issues
}

func (suite *CrossPlatformCompatibilityTestSuite) generateRecommendations(result CompatibilityResult) []string {
	recommendations := []string{}

	if result.Summary.CriticalIssues > 0 {
		recommendations = append(recommendations,
			"Address all critical compatibility issues before deployment")
	}

	if result.Summary.MajorIssues > 0 {
		recommendations = append(recommendations,
			"Review and resolve major compatibility issues to improve cross-platform consistency")
	}

	if result.Score < 90.0 {
		recommendations = append(recommendations,
			"Consider implementing contract-first development approach")
	}

	return recommendations
}

func (suite *CrossPlatformCompatibilityTestSuite) saveCompatibilityReport(result CompatibilityResult) {
	// Save report for CI/CD integration
	reportJSON, _ := json.MarshalIndent(result, "", "  ")

	reportPath := filepath.Join(".", "compatibility-report.json")
	err := ioutil.WriteFile(reportPath, reportJSON, 0644)
	if err != nil {
		suite.T().Logf("Warning: Could not save compatibility report: %v", err)
	} else {
		suite.T().Logf("Compatibility report saved to %s", reportPath)
	}
}

// TestCrossPlatformCompatibility runs the complete test suite
func TestCrossPlatformCompatibility(t *testing.T) {
	suite.Run(t, new(CrossPlatformCompatibilityTestSuite))
}