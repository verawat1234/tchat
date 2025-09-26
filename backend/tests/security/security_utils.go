package security

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

// SecurityUtils provides common security testing utilities
type SecurityUtils struct{}

// NewSecurityUtils creates a new security utils instance
func NewSecurityUtils() *SecurityUtils {
	return &SecurityUtils{}
}

// GenerateSecureToken generates a cryptographically secure token
func (s *SecurityUtils) GenerateSecureToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("fallback-token-%d", length)
	}
	return hex.EncodeToString(bytes)
}

// ValidateEmail validates email format according to RFC 5322
func (s *SecurityUtils) ValidateEmail(email string) bool {
	// Basic email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Check basic format
	if !emailRegex.MatchString(email) {
		return false
	}

	// Additional security checks
	if len(email) > 254 { // RFC 5321 limit
		return false
	}

	// Check for dangerous characters
	dangerousChars := []string{"<", ">", "\"", "'", "&", "%", "$", ";", "(", ")", "[", "]", "{", "}"}
	for _, char := range dangerousChars {
		if strings.Contains(email, char) {
			return false
		}
	}

	return true
}

// ValidatePhone validates phone number format for Southeast Asian countries
func (s *SecurityUtils) ValidatePhone(phone string) bool {
	// Remove all non-digit characters except +
	cleanPhone := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Check if it starts with + and has correct length
	if !strings.HasPrefix(cleanPhone, "+") {
		return false
	}

	// Southeast Asian country codes and their expected lengths
	seaPatterns := map[string][]int{
		"+66": {11, 12}, // Thailand: +66812345678 (11) or +66812345678 (12)
		"+65": {10, 11}, // Singapore: +6512345678 (10) or +6561234567 (11)
		"+62": {11, 13}, // Indonesia: +628123456789 (11-13)
		"+60": {10, 12}, // Malaysia: +60123456789 (10-12)
		"+84": {11, 12}, // Vietnam: +84123456789 (11-12)
		"+63": {12, 13}, // Philippines: +639123456789 (12-13)
	}

	for prefix, lengths := range seaPatterns {
		if strings.HasPrefix(cleanPhone, prefix) {
			phoneLength := len(cleanPhone)
			for _, validLength := range lengths {
				if phoneLength == validLength {
					return true
				}
			}
		}
	}

	return false
}

// ValidatePhoneNumber validates phone number for specific country
func (s *SecurityUtils) ValidatePhoneNumber(phone, country string) bool {
	// Remove all non-digit characters except +
	cleanPhone := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Country-specific validation
	switch country {
	case "TH":
		return strings.HasPrefix(cleanPhone, "+66") && (len(cleanPhone) == 11 || len(cleanPhone) == 12)
	case "SG":
		return strings.HasPrefix(cleanPhone, "+65") && (len(cleanPhone) == 10 || len(cleanPhone) == 11)
	case "ID":
		return strings.HasPrefix(cleanPhone, "+62") && len(cleanPhone) >= 11 && len(cleanPhone) <= 13
	case "MY":
		return strings.HasPrefix(cleanPhone, "+60") && len(cleanPhone) >= 10 && len(cleanPhone) <= 12
	case "VN":
		return strings.HasPrefix(cleanPhone, "+84") && (len(cleanPhone) == 11 || len(cleanPhone) == 12)
	case "PH":
		return strings.HasPrefix(cleanPhone, "+63") && (len(cleanPhone) == 12 || len(cleanPhone) == 13)
	default:
		return s.ValidatePhone(phone)
	}
}

// ValidateUUID validates UUID format
func (s *SecurityUtils) ValidateUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}

// ValidatePassword validates password strength
func (s *SecurityUtils) ValidatePassword(password string) (bool, []string) {
	var errors []string

	// Minimum length
	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}

	// Maximum length (to prevent DoS)
	if len(password) > 128 {
		errors = append(errors, "Password must not exceed 128 characters")
	}

	// Check for required character types
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}
	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}
	if !hasDigit {
		errors = append(errors, "Password must contain at least one digit")
	}
	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character")
	}

	// Check for common patterns
	if s.containsCommonPatterns(password) {
		errors = append(errors, "Password contains common patterns and is not secure")
	}

	return len(errors) == 0, errors
}

// containsCommonPatterns checks for common weak password patterns
func (s *SecurityUtils) containsCommonPatterns(password string) bool {
	lowerPassword := strings.ToLower(password)

	// Common weak patterns
	weakPatterns := []string{
		"123456", "password", "admin", "test", "user",
		"abcdef", "qwerty", "111111", "000000", "letmein",
		"welcome", "monkey", "dragon", "master", "shadow",
		"tchat", "chat", "app", "mobile", "api",
	}

	for _, pattern := range weakPatterns {
		if strings.Contains(lowerPassword, pattern) {
			return true
		}
	}

	// Sequential characters
	if s.containsSequentialChars(password) {
		return true
	}

	// Repeated characters
	if s.containsRepeatedChars(password) {
		return true
	}

	return false
}

// containsSequentialChars checks for sequential character patterns
func (s *SecurityUtils) containsSequentialChars(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		// Check for ascending sequence
		if password[i]+1 == password[i+1] && password[i+1]+1 == password[i+2] {
			return true
		}
		// Check for descending sequence
		if password[i]-1 == password[i+1] && password[i+1]-1 == password[i+2] {
			return true
		}
	}

	return false
}

// containsRepeatedChars checks for repeated character patterns
func (s *SecurityUtils) containsRepeatedChars(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}

	return false
}

// SanitizeInput sanitizes input to prevent various injection attacks
func (s *SecurityUtils) SanitizeInput(input string, allowHTML bool) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove other control characters except allowed whitespace
	var result strings.Builder
	for _, r := range input {
		if r >= 32 || r == '\t' || r == '\n' || r == '\r' {
			result.WriteRune(r)
		}
	}
	input = result.String()

	if !allowHTML {
		// Escape HTML characters
		input = strings.ReplaceAll(input, "<", "&lt;")
		input = strings.ReplaceAll(input, ">", "&gt;")
		input = strings.ReplaceAll(input, "\"", "&quot;")
		input = strings.ReplaceAll(input, "'", "&#x27;")
		input = strings.ReplaceAll(input, "&", "&amp;")
	}

	return input
}

// ValidateFilename validates filename for security
func (s *SecurityUtils) ValidateFilename(filename string) (bool, []string) {
	var errors []string

	// Check for empty filename
	if strings.TrimSpace(filename) == "" {
		errors = append(errors, "Filename cannot be empty")
		return false, errors
	}

	// Check for path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		errors = append(errors, "Filename contains invalid characters")
	}

	// Check for dangerous extensions
	dangerousExts := []string{
		".exe", ".bat", ".cmd", ".com", ".pif", ".scr", ".vbs", ".js",
		".jar", ".php", ".asp", ".jsp", ".py", ".rb", ".pl", ".sh",
		".ps1", ".vb", ".reg", ".msi", ".dll", ".so", ".dylib",
	}

	lowerFilename := strings.ToLower(filename)
	for _, ext := range dangerousExts {
		if strings.HasSuffix(lowerFilename, ext) {
			errors = append(errors, fmt.Sprintf("File extension %s is not allowed", ext))
		}
	}

	// Check for reserved names (Windows)
	reservedNames := []string{
		"con", "prn", "aux", "nul", "com1", "com2", "com3", "com4",
		"com5", "com6", "com7", "com8", "com9", "lpt1", "lpt2",
		"lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9",
	}

	baseFilename := strings.TrimSuffix(lowerFilename, filepath.Ext(lowerFilename))
	for _, reserved := range reservedNames {
		if baseFilename == reserved {
			errors = append(errors, "Filename uses a reserved system name")
		}
	}

	// Check length
	if len(filename) > 255 {
		errors = append(errors, "Filename is too long (maximum 255 characters)")
	}

	return len(errors) == 0, errors
}

// ValidateJSONInput validates JSON input for security
func (s *SecurityUtils) ValidateJSONInput(jsonStr string) (bool, []string) {
	var errors []string

	// Check for maximum size (prevent DoS)
	if len(jsonStr) > 1024*1024 { // 1MB limit
		errors = append(errors, "JSON input is too large")
	}

	// Check for deeply nested structures (prevent DoS)
	nestingLevel := 0
	maxNesting := 20

	for _, char := range jsonStr {
		if char == '{' || char == '[' {
			nestingLevel++
			if nestingLevel > maxNesting {
				errors = append(errors, "JSON structure is too deeply nested")
				break
			}
		} else if char == '}' || char == ']' {
			nestingLevel--
		}
	}

	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"__proto__", "constructor", "prototype",
		"function", "eval", "require", "import",
		"$where", "$ne", "$gt", "$regex",
	}

	lowerJSON := strings.ToLower(jsonStr)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerJSON, pattern) {
			errors = append(errors, fmt.Sprintf("JSON contains suspicious pattern: %s", pattern))
		}
	}

	return len(errors) == 0, errors
}

// RateLimitTracker tracks rate limiting for security testing
type RateLimitTracker struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
}

// NewRateLimitTracker creates a new rate limit tracker
func NewRateLimitTracker() *RateLimitTracker {
	return &RateLimitTracker{
		requests: make(map[string][]time.Time),
	}
}

// CheckRateLimit checks if a request exceeds rate limits
func (r *RateLimitTracker) CheckRateLimit(identifier string, maxRequests int, window time.Duration) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-window)

	// Get existing requests for this identifier
	requests := r.requests[identifier]

	// Filter requests within the time window
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if adding this request would exceed the limit
	if len(validRequests) >= maxRequests {
		return false // Rate limit exceeded
	}

	// Add current request
	validRequests = append(validRequests, now)
	r.requests[identifier] = validRequests

	return true // Within rate limit
}

// SecurityTestConfig holds configuration for security tests
type SecurityTestConfig struct {
	EnableSQLInjectionTests    bool
	EnableXSSTests            bool
	EnableNoSQLInjectionTests bool
	EnableCommandInjectionTests bool
	EnablePathTraversalTests  bool
	EnableLDAPInjectionTests  bool
	EnableFileUploadTests     bool
	EnableRateLimitTests      bool
	MaxRequestsPerMinute      int
	MaxPasswordLength         int
	AllowedFileExtensions     []string
	MaxFileSize               int64
	EnableCORSTests           bool
	EnableCSRFTests           bool
}

// DefaultSecurityTestConfig returns a default security test configuration
func DefaultSecurityTestConfig() *SecurityTestConfig {
	return &SecurityTestConfig{
		EnableSQLInjectionTests:     true,
		EnableXSSTests:             true,
		EnableNoSQLInjectionTests:  true,
		EnableCommandInjectionTests: true,
		EnablePathTraversalTests:   true,
		EnableLDAPInjectionTests:   true,
		EnableFileUploadTests:      true,
		EnableRateLimitTests:       true,
		MaxRequestsPerMinute:       60,
		MaxPasswordLength:          128,
		AllowedFileExtensions:      []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".txt", ".doc", ".docx"},
		MaxFileSize:                10 * 1024 * 1024, // 10MB
		EnableCORSTests:            true,
		EnableCSRFTests:            true,
	}
}

// SecurityAuditReport represents a security audit report
type SecurityAuditReport struct {
	TestName           string                 `json:"test_name"`
	Timestamp          time.Time              `json:"timestamp"`
	TotalTests         int                    `json:"total_tests"`
	PassedTests        int                    `json:"passed_tests"`
	FailedTests        int                    `json:"failed_tests"`
	VulnerabilitiesFound []VulnerabilityReport `json:"vulnerabilities_found"`
	Recommendations    []string               `json:"recommendations"`
	RiskScore          int                    `json:"risk_score"` // 0-100
	ComplianceStatus   map[string]bool        `json:"compliance_status"`
}

// VulnerabilityReport represents a single vulnerability
type VulnerabilityReport struct {
	Type        string            `json:"type"`
	Severity    string            `json:"severity"` // Low, Medium, High, Critical
	Endpoint    string            `json:"endpoint"`
	Field       string            `json:"field"`
	Payload     string            `json:"payload,omitempty"`
	Description string            `json:"description"`
	Impact      string            `json:"impact"`
	Remediation string            `json:"remediation"`
	CVSS        float64           `json:"cvss,omitempty"`
	References  []string          `json:"references,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// GenerateSecurityAuditReport generates a comprehensive security audit report
func (s *SecurityUtils) GenerateSecurityAuditReport(testResults []VulnerabilityReport) *SecurityAuditReport {
	report := &SecurityAuditReport{
		TestName:           "Tchat Input Validation Security Audit",
		Timestamp:          time.Now(),
		VulnerabilitiesFound: testResults,
		ComplianceStatus:   make(map[string]bool),
	}

	// Calculate statistics
	report.TotalTests = len(testResults)
	for _, vuln := range testResults {
		if vuln.Severity == "Critical" || vuln.Severity == "High" {
			report.FailedTests++
		} else {
			report.PassedTests++
		}
	}

	// Calculate risk score
	report.RiskScore = s.calculateRiskScore(testResults)

	// Generate recommendations
	report.Recommendations = s.generateRecommendations(testResults)

	// Check compliance status
	report.ComplianceStatus["OWASP_Top_10"] = s.checkOWASPCompliance(testResults)
	report.ComplianceStatus["Input_Validation"] = report.FailedTests == 0
	report.ComplianceStatus["Data_Sanitization"] = s.checkSanitizationCompliance(testResults)

	return report
}

// calculateRiskScore calculates overall risk score based on vulnerabilities
func (s *SecurityUtils) calculateRiskScore(vulnerabilities []VulnerabilityReport) int {
	if len(vulnerabilities) == 0 {
		return 0
	}

	totalScore := 0
	for _, vuln := range vulnerabilities {
		switch vuln.Severity {
		case "Critical":
			totalScore += 40
		case "High":
			totalScore += 25
		case "Medium":
			totalScore += 10
		case "Low":
			totalScore += 5
		}
	}

	// Cap at 100
	if totalScore > 100 {
		totalScore = 100
	}

	return totalScore
}

// generateRecommendations generates security recommendations
func (s *SecurityUtils) generateRecommendations(vulnerabilities []VulnerabilityReport) []string {
	recommendations := []string{
		"Implement comprehensive input validation for all user inputs",
		"Use parameterized queries to prevent SQL injection",
		"Sanitize all output to prevent XSS attacks",
		"Implement proper file upload validation and restrictions",
		"Use rate limiting to prevent brute force attacks",
		"Implement proper error handling to avoid information disclosure",
		"Regular security testing and code reviews",
		"Keep all dependencies and frameworks updated",
	}

	// Add specific recommendations based on found vulnerabilities
	vulnTypes := make(map[string]bool)
	for _, vuln := range vulnerabilities {
		vulnTypes[vuln.Type] = true
	}

	if vulnTypes["SQL Injection"] {
		recommendations = append(recommendations, "Implement prepared statements and stored procedures")
	}
	if vulnTypes["XSS"] {
		recommendations = append(recommendations, "Implement Content Security Policy (CSP) headers")
	}
	if vulnTypes["File Upload"] {
		recommendations = append(recommendations, "Implement file type validation and virus scanning")
	}

	return recommendations
}

// checkOWASPCompliance checks compliance with OWASP Top 10
func (s *SecurityUtils) checkOWASPCompliance(vulnerabilities []VulnerabilityReport) bool {
	// Check for OWASP Top 10 vulnerabilities
	owaspVulns := []string{
		"SQL Injection", "XSS", "Command Injection", "Path Traversal",
		"NoSQL Injection", "LDAP Injection",
	}

	for _, vuln := range vulnerabilities {
		for _, owaspVuln := range owaspVulns {
			if strings.Contains(vuln.Type, owaspVuln) &&
			   (vuln.Severity == "Critical" || vuln.Severity == "High") {
				return false
			}
		}
	}

	return true
}

// checkSanitizationCompliance checks data sanitization compliance
func (s *SecurityUtils) checkSanitizationCompliance(vulnerabilities []VulnerabilityReport) bool {
	sanitizationVulns := []string{"XSS", "HTML Injection", "Script Injection"}

	for _, vuln := range vulnerabilities {
		for _, sanitizationVuln := range sanitizationVulns {
			if strings.Contains(vuln.Type, sanitizationVuln) {
				return false
			}
		}
	}

	return true
}

// HTML Sanitization Methods

// SanitizeHTML sanitizes HTML content to prevent XSS
func (s *SecurityUtils) SanitizeHTML(input string) string {
	// Escape HTML characters
	return html.EscapeString(input)
}

// IsHTMLSafe checks if HTML content is safe
func (s *SecurityUtils) IsHTMLSafe(input string) bool {
	// Check for dangerous HTML patterns
	dangerousPatterns := []string{
		"<script", "</script>", "javascript:", "onload=", "onerror=",
		"onclick=", "onmouseover=", "<iframe", "</iframe>",
		"<object", "</object>", "<embed", "</embed>",
		"<style", "</style>", "expression(", "data:text/html",
		"vbscript:", "data:application/",
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return false
		}
	}

	return true
}

// JSON Sanitization Methods

// IsJSONSafe validates if JSON content is safe
func (s *SecurityUtils) IsJSONSafe(jsonStr string) bool {
	// Check size limit
	if len(jsonStr) > 1024*1024 { // 1MB
		return false
	}

	// Check if it's valid JSON
	var temp interface{}
	if err := json.Unmarshal([]byte(jsonStr), &temp); err != nil {
		return false
	}

	// Check for dangerous patterns
	dangerousPatterns := []string{
		"<script", "javascript:", "__proto__", "constructor",
		"prototype", "function", "eval", "require",
	}

	lowerJSON := strings.ToLower(jsonStr)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerJSON, pattern) {
			return false
		}
	}

	// Check nesting depth
	return s.checkJSONNestingDepth(jsonStr, 20)
}

// checkJSONNestingDepth checks if JSON nesting is within limits
func (s *SecurityUtils) checkJSONNestingDepth(jsonStr string, maxDepth int) bool {
	depth := 0
	for _, char := range jsonStr {
		if char == '{' || char == '[' {
			depth++
			if depth > maxDepth {
				return false
			}
		} else if char == '}' || char == ']' {
			depth--
		}
	}
	return true
}

// URL Sanitization Methods

// IsURLSafe validates if URL is safe
func (s *SecurityUtils) IsURLSafe(urlStr string) bool {
	// Parse URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check scheme
	allowedSchemes := []string{"http", "https"}
	schemeAllowed := false
	for _, scheme := range allowedSchemes {
		if parsedURL.Scheme == scheme {
			schemeAllowed = true
			break
		}
	}
	if !schemeAllowed {
		return false
	}

	// Check for dangerous patterns
	if strings.Contains(strings.ToLower(urlStr), "javascript:") ||
		strings.Contains(strings.ToLower(urlStr), "data:text/html") ||
		strings.Contains(strings.ToLower(urlStr), "vbscript:") {
		return false
	}

	// Check for localhost/private IPs
	host := parsedURL.Host
	if strings.Contains(host, "localhost") ||
		strings.Contains(host, "127.0.0.1") ||
		strings.Contains(host, "192.168.") ||
		strings.Contains(host, "10.") ||
		strings.Contains(host, "172.16.") {
		return false
	}

	// Check for credentials in URL
	if parsedURL.User != nil {
		return false
	}

	return true
}

// Filename Sanitization Methods

// IsFilenameSafe checks if filename is safe
func (s *SecurityUtils) IsFilenameSafe(filename string) bool {
	valid, _ := s.ValidateFilename(filename)
	return valid
}

// SanitizeFilename sanitizes filename
func (s *SecurityUtils) SanitizeFilename(filename string) string {
	// Remove path components
	filename = filepath.Base(filename)

	// Replace dangerous characters
	filename = regexp.MustCompile(`[<>:"|?*\\]`).ReplaceAllString(filename, "_")

	// Replace path traversal
	filename = strings.ReplaceAll(filename, "..", "_")

	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// Remove null bytes
	filename = strings.ReplaceAll(filename, "\x00", "_")

	// Handle hidden files
	if strings.HasPrefix(filename, ".") && len(filename) > 1 {
		filename = "_" + filename[1:]
	}

	// Limit length
	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		base := filename[:255-len(ext)]
		filename = base + ext
	}

	return filename
}

// Data Type Validation Methods

// ValidateDataType validates data based on expected type
func (s *SecurityUtils) ValidateDataType(dataType string, value interface{}) bool {
	switch dataType {
	case "integer":
		switch v := value.(type) {
		case int, int32, int64:
			return true
		case string:
			_, err := strconv.Atoi(v)
			return err == nil
		default:
			return false
		}
	case "email":
		if str, ok := value.(string); ok {
			return s.ValidateEmail(str)
		}
		return false
	case "uuid":
		if str, ok := value.(string); ok {
			return s.ValidateUUID(str)
		}
		return false
	case "phone":
		if str, ok := value.(string); ok {
			return s.ValidatePhone(str)
		}
		return false
	default:
		return true
	}
}

// Southeast Asian Content Methods

// IsSEAContentSafe validates Southeast Asian content
func (s *SecurityUtils) IsSEAContentSafe(content, country string) bool {
	// Check basic safety first
	if !s.IsHTMLSafe(content) {
		return false
	}

	// Check length limits
	if len(content) > 10000 {
		return false
	}

	// Check for valid UTF-8
	if !utf8.ValidString(content) {
		return false
	}

	// Country-specific validation could be added here
	// For now, basic validation is sufficient

	return true
}

// SanitizeSEAContent sanitizes Southeast Asian content
func (s *SecurityUtils) SanitizeSEAContent(content, country string) string {
	// Normalize Unicode
	// Basic sanitization - remove dangerous HTML
	content = s.SanitizeHTML(content)

	// Trim to reasonable length
	if len(content) > 5000 {
		content = content[:5000] + "..."
	}

	return content
}

// Input Length Validation Methods

// ValidateInputLength validates input length
func (s *SecurityUtils) ValidateInputLength(input string, maxLength int) bool {
	// Count Unicode characters, not bytes
	return utf8.RuneCountInString(input) <= maxLength
}

// Special Character Validation Methods

// ValidateSpecialCharacters validates special characters in input
func (s *SecurityUtils) ValidateSpecialCharacters(input, mode string) bool {
	switch mode {
	case "alphanumeric":
		for _, r := range input {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return false
			}
		}
		return true
	case "text":
		// Allow basic punctuation but not dangerous characters
		dangerousChars := []string{"<", ">", "\"", "'", "&", "$", ";", "(", ")", "[", "]", "{", "}"}
		for _, char := range dangerousChars {
			if strings.Contains(input, char) {
				return false
			}
		}

		// Check for control characters
		for _, r := range input {
			if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
				return false
			}
		}
		return true
	default:
		return true
	}
}

// SanitizeSpecialCharacters sanitizes special characters
func (s *SecurityUtils) SanitizeSpecialCharacters(input, mode string) string {
	switch mode {
	case "alphanumeric":
		var result strings.Builder
		for _, r := range input {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				result.WriteRune(r)
			}
		}
		return result.String()
	case "text":
		// Remove dangerous characters
		input = regexp.MustCompile(`[<>"'&$;()\[\]{}]`).ReplaceAllString(input, "")

		// Remove control characters except allowed whitespace
		var result strings.Builder
		for _, r := range input {
			if !unicode.IsControl(r) || r == '\t' || r == '\n' || r == '\r' {
				result.WriteRune(r)
			}
		}
		return result.String()
	default:
		return input
	}
}

// Encoding Validation Methods

// ValidateUTF8Encoding validates UTF-8 encoding
func (s *SecurityUtils) ValidateUTF8Encoding(data []byte) bool {
	// Check for valid UTF-8
	if !utf8.Valid(data) {
		return false
	}

	// Check for null bytes
	for _, b := range data {
		if b == 0 {
			return false
		}
	}

	// Check for overlong encodings and other dangerous patterns
	// This is a simplified check - production code should use more comprehensive validation
	str := string(data)

	// Check for dangerous control characters
	for _, r := range str {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}

	return true
}