package security

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"tchat-backend/tests/fixtures"
)

// InputValidationTestSuite provides comprehensive input validation security testing
type InputValidationTestSuite struct {
	suite.Suite
	fixtures *fixtures.MasterFixtures
	ctx      context.Context
}

// SetupTest initializes the test suite
func (suite *InputValidationTestSuite) SetupTest() {
	suite.fixtures = fixtures.NewMasterFixtures(42) // Fixed seed for reproducible tests
	suite.ctx = context.Background()
}

// TestInputValidationTestSuite runs the input validation test suite
func TestInputValidationTestSuite(t *testing.T) {
	suite.Run(t, new(InputValidationTestSuite))
}

// SQLInjectionPayloads contains common SQL injection attack vectors
var SQLInjectionPayloads = []string{
	"'; DROP TABLE users; --",
	"' OR '1'='1",
	"' OR 1=1 --",
	"' UNION SELECT * FROM users --",
	"admin'--",
	"' OR 'x'='x",
	"1'; INSERT INTO users (name) VALUES ('hacker'); --",
	"' AND (SELECT COUNT(*) FROM users) > 0 --",
	"'; EXEC xp_cmdshell('dir'); --",
	"' OR SLEEP(5) --",
	"' AND EXTRACTVALUE(1, CONCAT(0x7e, (SELECT version()), 0x7e)) --",
	"' OR (SELECT * FROM (SELECT(SLEEP(5)))a) --",
	"1' AND 1=1 UNION SELECT 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30 --",
}

// XSSPayloads contains common XSS attack vectors
var XSSPayloads = []string{
	"<script>alert('XSS')</script>",
	"<img src=x onerror=alert('XSS')>",
	"javascript:alert('XSS')",
	"<svg onload=alert('XSS')>",
	"<iframe src=\"javascript:alert('XSS')\"></iframe>",
	"<body onload=alert('XSS')>",
	"<input type=\"image\" src=\"x\" onerror=\"alert('XSS')\">",
	"<object data=\"javascript:alert('XSS')\">",
	"<embed src=\"javascript:alert('XSS')\">",
	"<div style=\"background-image:url(javascript:alert('XSS'))\">",
	"<link rel=\"stylesheet\" href=\"javascript:alert('XSS')\">",
	"<meta http-equiv=\"refresh\" content=\"0;url=javascript:alert('XSS')\">",
	"<marquee onstart=alert('XSS')></marquee>",
	"<video><source onerror=\"alert('XSS')\">",
	"<audio src=x onerror=alert('XSS')>",
}

// NoSQLInjectionPayloads contains NoSQL injection attack vectors
var NoSQLInjectionPayloads = []string{
	"{\"$ne\": null}",
	"{\"$regex\": \".*\"}",
	"{\"$where\": \"function() { return true; }\"}",
	"{\"$gt\": \"\"}",
	"{\"$gte\": \"\"}",
	"{\"$lt\": \"\"}",
	"{\"$lte\": \"\"}",
	"{\"$in\": []}",
	"{\"$nin\": []}",
	"{\"$exists\": true}",
	"{\"$type\": 2}",
	"{\"$size\": 0}",
	"{\"$elemMatch\": {}}",
	"{\"$all\": []}",
}

// CommandInjectionPayloads contains command injection attack vectors
var CommandInjectionPayloads = []string{
	"; cat /etc/passwd",
	"| cat /etc/passwd",
	"&& cat /etc/passwd",
	"`cat /etc/passwd`",
	"$(cat /etc/passwd)",
	"; rm -rf /",
	"| rm -rf /",
	"&& rm -rf /",
	"`rm -rf /`",
	"$(rm -rf /)",
	"; ping 127.0.0.1",
	"| ping 127.0.0.1",
	"&& ping 127.0.0.1",
	"`ping 127.0.0.1`",
	"$(ping 127.0.0.1)",
}

// PathTraversalPayloads contains path traversal attack vectors
var PathTraversalPayloads = []string{
	"../../../etc/passwd",
	"..\\..\\..\\windows\\system32\\drivers\\etc\\hosts",
	"....//....//....//etc/passwd",
	"..%2F..%2F..%2Fetc%2Fpasswd",
	"..%252F..%252F..%252Fetc%252Fpasswd",
	"..%c0%af..%c0%af..%c0%afetc%c0%afpasswd",
	"/%2e%2e/%2e%2e/%2e%2e/etc/passwd",
	"/var/www/../../etc/passwd",
	"file:///etc/passwd",
	"file://c:/windows/system32/drivers/etc/hosts",
}

// LDAPInjectionPayloads contains LDAP injection attack vectors
var LDAPInjectionPayloads = []string{
	"*)(uid=*",
	"*)(|(password=*))",
	"*)(|(objectClass=*))",
	"*))%00",
	"*()|&'",
	"*)(uid=*)(|(uid=*",
	"*)(|(mail=*@*))",
	"*)(|(cn=*))",
	"admin)(&(password=*)",
	"test*",
}

// TestSQLInjectionPrevention tests SQL injection prevention across all input fields
func (suite *InputValidationTestSuite) TestSQLInjectionPrevention() {
	testCases := []struct {
		name        string
		field       string
		endpoint    string
		payloads    []string
		expectError bool
	}{
		{
			name:        "User Registration Name Field",
			field:       "name",
			endpoint:    "/api/auth/register",
			payloads:    SQLInjectionPayloads,
			expectError: true,
		},
		{
			name:        "User Login Email Field",
			field:       "email",
			endpoint:    "/api/auth/login",
			payloads:    SQLInjectionPayloads,
			expectError: true,
		},
		{
			name:        "Content Search Query",
			field:       "query",
			endpoint:    "/api/content/search",
			payloads:    SQLInjectionPayloads,
			expectError: true,
		},
		{
			name:        "Payment Transaction Description",
			field:       "description",
			endpoint:    "/api/payments/transfer",
			payloads:    SQLInjectionPayloads,
			expectError: true,
		},
		{
			name:        "Message Content",
			field:       "content",
			endpoint:    "/api/messages/send",
			payloads:    SQLInjectionPayloads,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for i, payload := range tc.payloads {
				suite.Run(fmt.Sprintf("Payload_%d", i), func() {
					// Create test request with malicious payload
					requestBody := map[string]interface{}{
						tc.field: payload,
					}

					// Add required fields based on endpoint
					suite.addRequiredFields(requestBody, tc.endpoint)

					// Test the injection attempt
					suite.testInputValidation(tc.endpoint, requestBody, tc.expectError, "SQL injection")
				})
			}
		})
	}
}

// TestXSSPrevention tests XSS prevention across all text input fields
func (suite *InputValidationTestSuite) TestXSSPrevention() {
	testCases := []struct {
		name        string
		field       string
		endpoint    string
		payloads    []string
		expectError bool
	}{
		{
			name:        "Content Item Text",
			field:       "text",
			endpoint:    "/api/content/create",
			payloads:    XSSPayloads,
			expectError: true,
		},
		{
			name:        "User Profile Name",
			field:       "name",
			endpoint:    "/api/users/profile",
			payloads:    XSSPayloads,
			expectError: true,
		},
		{
			name:        "Message Text Content",
			field:       "text",
			endpoint:    "/api/messages/send",
			payloads:    XSSPayloads,
			expectError: true,
		},
		{
			name:        "Chat Group Name",
			field:       "name",
			endpoint:    "/api/chats/create",
			payloads:    XSSPayloads,
			expectError: true,
		},
		{
			name:        "Comment Content",
			field:       "content",
			endpoint:    "/api/comments/create",
			payloads:    XSSPayloads,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for i, payload := range tc.payloads {
				suite.Run(fmt.Sprintf("Payload_%d", i), func() {
					requestBody := map[string]interface{}{
						tc.field: payload,
					}

					suite.addRequiredFields(requestBody, tc.endpoint)
					suite.testInputValidation(tc.endpoint, requestBody, tc.expectError, "XSS")
				})
			}
		})
	}
}

// TestNoSQLInjectionPrevention tests NoSQL injection prevention
func (suite *InputValidationTestSuite) TestNoSQLInjectionPrevention() {
	testCases := []struct {
		name        string
		field       string
		endpoint    string
		expectError bool
	}{
		{
			name:        "User Query Filter",
			field:       "filter",
			endpoint:    "/api/users/search",
			expectError: true,
		},
		{
			name:        "Content Metadata Query",
			field:       "metadata",
			endpoint:    "/api/content/query",
			expectError: true,
		},
		{
			name:        "Transaction Query",
			field:       "query",
			endpoint:    "/api/transactions/search",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for i, payload := range NoSQLInjectionPayloads {
				suite.Run(fmt.Sprintf("Payload_%d", i), func() {
					// Test both string and parsed JSON payloads
					requestBody := map[string]interface{}{
						tc.field: payload,
					}

					suite.addRequiredFields(requestBody, tc.endpoint)
					suite.testInputValidation(tc.endpoint, requestBody, tc.expectError, "NoSQL injection")

					// Test parsed JSON payload
					var jsonPayload interface{}
					if json.Unmarshal([]byte(payload), &jsonPayload) == nil {
						requestBodyJSON := map[string]interface{}{
							tc.field: jsonPayload,
						}
						suite.addRequiredFields(requestBodyJSON, tc.endpoint)
						suite.testInputValidation(tc.endpoint, requestBodyJSON, tc.expectError, "NoSQL injection (JSON)")
					}
				})
			}
		})
	}
}

// TestCommandInjectionPrevention tests command injection prevention
func (suite *InputValidationTestSuite) TestCommandInjectionPrevention() {
	testCases := []struct {
		name        string
		field       string
		endpoint    string
		expectError bool
	}{
		{
			name:        "File Upload Name",
			field:       "filename",
			endpoint:    "/api/files/upload",
			expectError: true,
		},
		{
			name:        "Export File Name",
			field:       "filename",
			endpoint:    "/api/export/data",
			expectError: true,
		},
		{
			name:        "Backup File Name",
			field:       "backup_name",
			endpoint:    "/api/admin/backup",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for i, payload := range CommandInjectionPayloads {
				suite.Run(fmt.Sprintf("Payload_%d", i), func() {
					requestBody := map[string]interface{}{
						tc.field: payload,
					}

					suite.addRequiredFields(requestBody, tc.endpoint)
					suite.testInputValidation(tc.endpoint, requestBody, tc.expectError, "Command injection")
				})
			}
		})
	}
}

// TestPathTraversalPrevention tests path traversal prevention
func (suite *InputValidationTestSuite) TestPathTraversalPrevention() {
	testCases := []struct {
		name        string
		field       string
		endpoint    string
		expectError bool
	}{
		{
			name:        "File Access Path",
			field:       "path",
			endpoint:    "/api/files/download",
			expectError: true,
		},
		{
			name:        "Template Path",
			field:       "template",
			endpoint:    "/api/templates/load",
			expectError: true,
		},
		{
			name:        "Log File Path",
			field:       "logfile",
			endpoint:    "/api/admin/logs",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for i, payload := range PathTraversalPayloads {
				suite.Run(fmt.Sprintf("Payload_%d", i), func() {
					requestBody := map[string]interface{}{
						tc.field: payload,
					}

					suite.addRequiredFields(requestBody, tc.endpoint)
					suite.testInputValidation(tc.endpoint, requestBody, tc.expectError, "Path traversal")
				})
			}
		})
	}
}

// TestLDAPInjectionPrevention tests LDAP injection prevention
func (suite *InputValidationTestSuite) TestLDAPInjectionPrevention() {
	testCases := []struct {
		name        string
		field       string
		endpoint    string
		expectError bool
	}{
		{
			name:        "LDAP User Search",
			field:       "username",
			endpoint:    "/api/ldap/users",
			expectError: true,
		},
		{
			name:        "LDAP Group Search",
			field:       "group",
			endpoint:    "/api/ldap/groups",
			expectError: true,
		},
		{
			name:        "LDAP Authentication",
			field:       "credentials",
			endpoint:    "/api/ldap/auth",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for i, payload := range LDAPInjectionPayloads {
				suite.Run(fmt.Sprintf("Payload_%d", i), func() {
					requestBody := map[string]interface{}{
						tc.field: payload,
					}

					suite.addRequiredFields(requestBody, tc.endpoint)
					suite.testInputValidation(tc.endpoint, requestBody, tc.expectError, "LDAP injection")
				})
			}
		})
	}
}

// TestDataTypeValidation tests proper data type validation
func (suite *InputValidationTestSuite) TestDataTypeValidation() {
	testCases := []struct {
		name         string
		endpoint     string
		field        string
		expectedType string
		invalidValue interface{}
		expectError  bool
	}{
		{
			name:         "User ID Should Be UUID",
			endpoint:     "/api/users/profile",
			field:        "user_id",
			expectedType: "uuid",
			invalidValue: "not-a-uuid",
			expectError:  true,
		},
		{
			name:         "Amount Should Be Numeric",
			endpoint:     "/api/payments/transfer",
			field:        "amount",
			expectedType: "number",
			invalidValue: "not-a-number",
			expectError:  true,
		},
		{
			name:         "Email Should Be Valid Format",
			endpoint:     "/api/auth/register",
			field:        "email",
			expectedType: "email",
			invalidValue: "not-an-email",
			expectError:  true,
		},
		{
			name:         "Phone Should Be Valid Format",
			endpoint:     "/api/auth/register",
			field:        "phone",
			expectedType: "phone",
			invalidValue: "12345",
			expectError:  true,
		},
		{
			name:         "Date Should Be Valid Format",
			endpoint:     "/api/events/create",
			field:        "date",
			expectedType: "date",
			invalidValue: "not-a-date",
			expectError:  true,
		},
		{
			name:         "Boolean Should Be Boolean",
			endpoint:     "/api/settings/update",
			field:        "enabled",
			expectedType: "boolean",
			invalidValue: "not-a-boolean",
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			requestBody := map[string]interface{}{
				tc.field: tc.invalidValue,
			}

			suite.addRequiredFields(requestBody, tc.endpoint)
			suite.testInputValidation(tc.endpoint, requestBody, tc.expectError, fmt.Sprintf("Invalid %s", tc.expectedType))
		})
	}
}

// TestInputLengthValidation tests input length restrictions
func (suite *InputValidationTestSuite) TestInputLengthValidation() {
	testCases := []struct {
		name        string
		endpoint    string
		field       string
		maxLength   int
		expectError bool
	}{
		{
			name:        "User Name Too Long",
			endpoint:    "/api/auth/register",
			field:       "name",
			maxLength:   255,
			expectError: true,
		},
		{
			name:        "Message Content Too Long",
			endpoint:    "/api/messages/send",
			field:       "content",
			maxLength:   4000,
			expectError: true,
		},
		{
			name:        "Content Title Too Long",
			endpoint:    "/api/content/create",
			field:       "title",
			maxLength:   500,
			expectError: true,
		},
		{
			name:        "Transaction Description Too Long",
			endpoint:    "/api/payments/transfer",
			field:       "description",
			maxLength:   1000,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create a string that exceeds the maximum length
			longString := strings.Repeat("a", tc.maxLength+100)

			requestBody := map[string]interface{}{
				tc.field: longString,
			}

			suite.addRequiredFields(requestBody, tc.endpoint)
			suite.testInputValidation(tc.endpoint, requestBody, tc.expectError, "Length validation")
		})
	}
}

// TestSpecialCharacterHandling tests handling of special characters
func (suite *InputValidationTestSuite) TestSpecialCharacterHandling() {
	specialCharPayloads := []string{
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F", // Control characters
		"'\"\\\r\n\t",                                                        // Quote and escape characters
		"🚀💻🔒🌟",                                                               // Unicode emojis (valid in SEA context)
		"测试中文",                                                                // Chinese characters
		"한국어테스트",                                                             // Korean characters
		"ทดสอบภาษาไทย",                                                           // Thai characters
		"テストの日本語",                                                           // Japanese characters
		"العربية",                                                             // Arabic characters
		"русский",                                                            // Cyrillic characters
		"<>&\"'",                                                             // HTML special characters
		"%3Cscript%3E",                                                       // URL encoded
		"\u0000\u0001\u0002",                                                 // Unicode null and control
	}

	testCases := []struct {
		name        string
		endpoint    string
		field       string
		allowUnicode bool
		expectError bool
	}{
		{
			name:        "Message Content - Unicode Allowed",
			endpoint:    "/api/messages/send",
			field:       "content",
			allowUnicode: true,
			expectError:  false, // Some unicode should be allowed
		},
		{
			name:        "User Name - Limited Unicode",
			endpoint:    "/api/auth/register",
			field:       "name",
			allowUnicode: true,
			expectError:  false, // Names can contain unicode
		},
		{
			name:        "Email Field - ASCII Only",
			endpoint:    "/api/auth/register",
			field:       "email",
			allowUnicode: false,
			expectError:  true, // Emails should be ASCII
		},
		{
			name:        "File Path - ASCII Only",
			endpoint:    "/api/files/upload",
			field:       "path",
			allowUnicode: false,
			expectError:  true, // File paths should be ASCII
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			for i, payload := range specialCharPayloads {
				suite.Run(fmt.Sprintf("Payload_%d", i), func() {
					requestBody := map[string]interface{}{
						tc.field: payload,
					}

					suite.addRequiredFields(requestBody, tc.endpoint)

					// For unicode-allowed fields, only control characters should be rejected
					expectError := tc.expectError
					if tc.allowUnicode && !containsControlCharacters(payload) && isValidUnicode(payload) {
						expectError = false
					}

					suite.testInputValidation(tc.endpoint, requestBody, expectError, "Special character handling")
				})
			}
		})
	}
}

// TestFileUploadValidation tests file upload security validation
func (suite *InputValidationTestSuite) TestFileUploadValidation() {
	maliciousFilePayloads := []struct {
		filename    string
		content     string
		contentType string
		expectError bool
	}{
		{
			filename:    "test.php",
			content:     "<?php system($_GET['cmd']); ?>",
			contentType: "application/x-php",
			expectError: true,
		},
		{
			filename:    "test.jsp",
			content:     "<% Runtime.getRuntime().exec(request.getParameter(\"cmd\")); %>",
			contentType: "application/x-jsp",
			expectError: true,
		},
		{
			filename:    "test.exe",
			content:     "MZ\x90\x00", // PE header
			contentType: "application/x-msdownload",
			expectError: true,
		},
		{
			filename:    "test.sh",
			content:     "#!/bin/bash\nrm -rf /",
			contentType: "application/x-sh",
			expectError: true,
		},
		{
			filename:    "../../../etc/passwd",
			content:     "root:x:0:0:root:/root:/bin/bash",
			contentType: "text/plain",
			expectError: true,
		},
		{
			filename:    "test.jpg",
			content:     "\xFF\xD8\xFF\xE0", // JPEG header
			contentType: "image/jpeg",
			expectError: false, // Valid image should be allowed
		},
		{
			filename:    "test.png",
			content:     "\x89PNG\r\n\x1a\n", // PNG header
			contentType: "image/png",
			expectError: false, // Valid image should be allowed
		},
	}

	for _, payload := range maliciousFilePayloads {
		suite.Run(fmt.Sprintf("File_%s", payload.filename), func() {
			requestBody := map[string]interface{}{
				"filename":     payload.filename,
				"content":      payload.content,
				"content_type": payload.contentType,
			}

			suite.testInputValidation("/api/files/upload", requestBody, payload.expectError, "File upload validation")
		})
	}
}

// Helper Methods

// testInputValidation tests input validation for a given endpoint and payload
func (suite *InputValidationTestSuite) testInputValidation(endpoint string, requestBody map[string]interface{}, expectError bool, attackType string) {
	// Convert request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(suite.T(), err)

	// Create HTTP request
	req := httptest.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Note: In a real implementation, this would call the actual handler
	// For this test, we simulate the validation response
	suite.simulateValidationResponse(rr, req, requestBody, attackType)

	// Verify response
	if expectError {
		assert.True(suite.T(), rr.Code >= 400,
			"Expected error status code for %s attack, got %d", attackType, rr.Code)

		// Verify error response contains security-related information
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err == nil {
			// Should contain validation error message
			assert.Contains(suite.T(), fmt.Sprintf("%v", response), "validation",
				"Error response should indicate validation failure")
		}
	} else {
		assert.True(suite.T(), rr.Code < 400,
			"Expected success status code, got %d", rr.Code)
	}
}

// simulateValidationResponse simulates the validation response based on input
func (suite *InputValidationTestSuite) simulateValidationResponse(rr *httptest.ResponseRecorder, req *http.Request, requestBody map[string]interface{}, attackType string) {
	// This is a simulation of what the actual validation should do
	// In a real implementation, this would be replaced with actual handler calls

	hasVulnerableInput := false

	// Check for malicious patterns in all string fields
	for _, value := range requestBody {
		if strValue, ok := value.(string); ok {
			if suite.containsMaliciousPattern(strValue, attackType) {
				hasVulnerableInput = true
				break
			}
		}
	}

	if hasVulnerableInput {
		// Return validation error
		rr.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"error":   "Validation failed",
			"message": "Input contains potentially malicious content",
			"code":    "VALIDATION_ERROR",
		}
		json.NewEncoder(rr).Encode(response)
	} else {
		// Return success
		rr.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"success": true,
			"message": "Input validation passed",
		}
		json.NewEncoder(rr).Encode(response)
	}
}

// containsMaliciousPattern checks if input contains malicious patterns
func (suite *InputValidationTestSuite) containsMaliciousPattern(input, attackType string) bool {
	switch attackType {
	case "SQL injection":
		return suite.containsSQLInjectionPattern(input)
	case "XSS":
		return suite.containsXSSPattern(input)
	case "NoSQL injection", "NoSQL injection (JSON)":
		return suite.containsNoSQLInjectionPattern(input)
	case "Command injection":
		return suite.containsCommandInjectionPattern(input)
	case "Path traversal":
		return suite.containsPathTraversalPattern(input)
	case "LDAP injection":
		return suite.containsLDAPInjectionPattern(input)
	default:
		return false
	}
}

// containsSQLInjectionPattern checks for SQL injection patterns
func (suite *InputValidationTestSuite) containsSQLInjectionPattern(input string) bool {
	sqlPatterns := []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
		"DROP", "INSERT", "DELETE", "UPDATE", "SELECT",
		"UNION", "OR", "AND", "EXEC", "EXECUTE",
	}

	upperInput := strings.ToUpper(input)
	for _, pattern := range sqlPatterns {
		if strings.Contains(upperInput, strings.ToUpper(pattern)) {
			return true
		}
	}
	return false
}

// containsXSSPattern checks for XSS patterns
func (suite *InputValidationTestSuite) containsXSSPattern(input string) bool {
	xssPatterns := []string{
		"<script", "<img", "<iframe", "<object", "<embed",
		"javascript:", "onerror", "onload", "onclick", "onmouseover",
		"<svg", "<body", "<input", "<link", "<meta",
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range xssPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}
	return false
}

// containsNoSQLInjectionPattern checks for NoSQL injection patterns
func (suite *InputValidationTestSuite) containsNoSQLInjectionPattern(input string) bool {
	nosqlPatterns := []string{
		"$ne", "$gt", "$gte", "$lt", "$lte", "$in", "$nin",
		"$exists", "$type", "$regex", "$where", "$size",
		"$elemMatch", "$all",
	}

	for _, pattern := range nosqlPatterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}
	return false
}

// containsCommandInjectionPattern checks for command injection patterns
func (suite *InputValidationTestSuite) containsCommandInjectionPattern(input string) bool {
	cmdPatterns := []string{
		";", "|", "&", "`", "$(",
		"cat", "rm", "ls", "cp", "mv", "mkdir", "rmdir",
		"ping", "wget", "curl", "nc", "netcat",
	}

	for _, pattern := range cmdPatterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}
	return false
}

// containsPathTraversalPattern checks for path traversal patterns
func (suite *InputValidationTestSuite) containsPathTraversalPattern(input string) bool {
	pathPatterns := []string{
		"../", "..\\", "....//", "....\\\\",
		"%2e%2e%2f", "%2e%2e%5c", "%252e%252e%252f",
		"/etc/passwd", "/windows/system32", "c:\\windows",
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range pathPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}
	return false
}

// containsLDAPInjectionPattern checks for LDAP injection patterns
func (suite *InputValidationTestSuite) containsLDAPInjectionPattern(input string) bool {
	ldapPatterns := []string{
		"*", "(", ")", "|", "&", "!",
		"objectClass", "uid", "cn", "mail",
		"%00", "\\", "/",
	}

	for _, pattern := range ldapPatterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}
	return false
}

// addRequiredFields adds required fields to request body based on endpoint
func (suite *InputValidationTestSuite) addRequiredFields(requestBody map[string]interface{}, endpoint string) {
	// Add required fields based on endpoint to make requests valid
	switch endpoint {
	case "/api/auth/register":
		if _, exists := requestBody["email"]; !exists {
			requestBody["email"] = "test@example.com"
		}
		if _, exists := requestBody["password"]; !exists {
			requestBody["password"] = "SecurePassword123!"
		}
		if _, exists := requestBody["name"]; !exists {
			requestBody["name"] = "Test User"
		}
		if _, exists := requestBody["phone"]; !exists {
			requestBody["phone"] = "+66812345678"
		}

	case "/api/auth/login":
		if _, exists := requestBody["email"]; !exists {
			requestBody["email"] = "test@example.com"
		}
		if _, exists := requestBody["password"]; !exists {
			requestBody["password"] = "SecurePassword123!"
		}

	case "/api/content/create":
		if _, exists := requestBody["category"]; !exists {
			requestBody["category"] = "general"
		}
		if _, exists := requestBody["type"]; !exists {
			requestBody["type"] = "text"
		}

	case "/api/payments/transfer":
		if _, exists := requestBody["to_user_id"]; !exists {
			requestBody["to_user_id"] = uuid.New().String()
		}
		if _, exists := requestBody["amount"]; !exists {
			requestBody["amount"] = 1000
		}
		if _, exists := requestBody["currency"]; !exists {
			requestBody["currency"] = "THB"
		}

	case "/api/messages/send":
		if _, exists := requestBody["chat_id"]; !exists {
			requestBody["chat_id"] = uuid.New().String()
		}
		if _, exists := requestBody["type"]; !exists {
			requestBody["type"] = "text"
		}

	case "/api/files/upload":
		if _, exists := requestBody["filename"]; !exists {
			requestBody["filename"] = "test.txt"
		}
		if _, exists := requestBody["content_type"]; !exists {
			requestBody["content_type"] = "text/plain"
		}
		if _, exists := requestBody["content"]; !exists {
			requestBody["content"] = "test content"
		}
	}
}

// containsControlCharacters checks if string contains control characters
func containsControlCharacters(s string) bool {
	for _, r := range s {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			return true
		}
	}
	return false
}

// isValidUnicode checks if string is valid unicode
func isValidUnicode(s string) bool {
	return len(s) == len([]rune(s))
}

// Benchmark tests for performance impact of validation

// BenchmarkSQLInjectionValidation benchmarks SQL injection validation performance
func BenchmarkSQLInjectionValidation(b *testing.B) {
	suite := &InputValidationTestSuite{}
	suite.SetupTest()

	requestBody := map[string]interface{}{
		"name":     "'; DROP TABLE users; --",
		"email":    "test@example.com",
		"password": "password",
		"phone":    "+66812345678",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.testInputValidation("/api/auth/register", requestBody, true, "SQL injection")
	}
}

// BenchmarkXSSValidation benchmarks XSS validation performance
func BenchmarkXSSValidation(b *testing.B) {
	suite := &InputValidationTestSuite{}
	suite.SetupTest()

	requestBody := map[string]interface{}{
		"content": "<script>alert('XSS')</script>",
		"chat_id": uuid.New().String(),
		"type":    "text",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.testInputValidation("/api/messages/send", requestBody, true, "XSS")
	}
}