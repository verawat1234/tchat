package security_test

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"tchat.dev/tests/fixtures"
)

// DataSanitizationTestSuite provides comprehensive data sanitization testing
// for all input types across Tchat microservices
type DataSanitizationTestSuite struct {
	suite.Suite
	fixtures *fixtures.MasterFixtures
	ctx      context.Context
	utils    *SecurityUtils
}

// SetupSuite initializes the test suite
func (suite *DataSanitizationTestSuite) SetupSuite() {
	suite.fixtures = fixtures.NewMasterFixtures(12345)
	suite.ctx = context.Background()
	suite.utils = NewSecurityUtils()
}

// TestHTMLSanitization tests HTML content sanitization
func (suite *DataSanitizationTestSuite) TestHTMLSanitization() {
	testCases := []struct {
		name            string
		input           string
		expectedSafe    bool
		sanitizedOutput string
		description     string
	}{
		{
			name:            "Clean text content",
			input:           "Hello world, this is a normal message",
			expectedSafe:    true,
			sanitizedOutput: "Hello world, this is a normal message",
			description:     "Plain text should pass through unchanged",
		},
		{
			name:            "Script tag injection",
			input:           "<script>alert('XSS')</script>",
			expectedSafe:    false,
			sanitizedOutput: "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
			description:     "Script tags should be escaped",
		},
		{
			name:            "Image with javascript src",
			input:           "<img src='javascript:alert(1)'>",
			expectedSafe:    false,
			sanitizedOutput: "&lt;img src=&#39;javascript:alert(1)&#39;&gt;",
			description:     "JavaScript URLs should be escaped",
		},
		{
			name:            "iframe injection",
			input:           "<iframe src='http://evil.com'></iframe>",
			expectedSafe:    false,
			sanitizedOutput: "&lt;iframe src=&#39;http://evil.com&#39;&gt;&lt;/iframe&gt;",
			description:     "iframe tags should be escaped",
		},
		{
			name:            "Event handler injection",
			input:           "<div onclick='alert(1)'>Click me</div>",
			expectedSafe:    false,
			sanitizedOutput: "&lt;div onclick=&#39;alert(1)&#39;&gt;Click me&lt;/div&gt;",
			description:     "Event handlers should be escaped",
		},
		{
			name:            "Style injection",
			input:           "<style>body{background:url('javascript:alert(1)')}</style>",
			expectedSafe:    false,
			sanitizedOutput: "&lt;style&gt;body{background:url(&#39;javascript:alert(1)&#39;)}&lt;/style&gt;",
			description:     "Style tags with JavaScript should be escaped",
		},
		{
			name:            "Data URL injection",
			input:           "<img src='data:text/html,<script>alert(1)</script>'>",
			expectedSafe:    false,
			sanitizedOutput: "&lt;img src=&#39;data:text/html,&lt;script&gt;alert(1)&lt;/script&gt;&#39;&gt;",
			description:     "Data URLs with HTML should be escaped",
		},
		{
			name:            "Unicode normalization attack",
			input:           "<scri\u0070t>alert('XSS')</script>",
			expectedSafe:    false,
			sanitizedOutput: "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
			description:     "Unicode variations should be normalized and escaped",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Test HTML sanitization
			sanitized := suite.utils.SanitizeHTML(tc.input)
			isSafe := suite.utils.IsHTMLSafe(tc.input)

			suite.Equal(tc.expectedSafe, isSafe, tc.description)
			suite.Equal(tc.sanitizedOutput, sanitized, "Sanitized output should match expected")

			// Ensure no dangerous patterns remain
			suite.False(strings.Contains(strings.ToLower(sanitized), "<script"), "No script tags should remain")
			suite.False(strings.Contains(strings.ToLower(sanitized), "javascript:"), "No javascript: URLs should remain")
			suite.False(strings.Contains(strings.ToLower(sanitized), "onload="), "No event handlers should remain")
			suite.False(strings.Contains(strings.ToLower(sanitized), "onerror="), "No error handlers should remain")
		})
	}
}

// TestJSONSanitization tests JSON content sanitization
func (suite *DataSanitizationTestSuite) TestJSONSanitization() {
	testCases := []struct {
		name         string
		input        string
		expectedSafe bool
		description  string
	}{
		{
			name:         "Valid JSON object",
			input:        `{"name": "John", "age": 30}`,
			expectedSafe: true,
			description:  "Valid JSON should be safe",
		},
		{
			name:         "JSON with XSS in string",
			input:        `{"message": "<script>alert('XSS')</script>"}`,
			expectedSafe: false,
			description:  "JSON containing script tags should be flagged",
		},
		{
			name:         "JSON with JavaScript URL",
			input:        `{"url": "javascript:alert(1)"}`,
			expectedSafe: false,
			description:  "JSON containing JavaScript URLs should be flagged",
		},
		{
			name:         "Malformed JSON",
			input:        `{"name": "John", "age":}`,
			expectedSafe: false,
			description:  "Malformed JSON should be rejected",
		},
		{
			name:         "JSON with prototype pollution",
			input:        `{"__proto__": {"admin": true}}`,
			expectedSafe: false,
			description:  "Prototype pollution attempts should be flagged",
		},
		{
			name:         "Deeply nested JSON",
			input:        strings.Repeat(`{"level":`, 1000) + `"deep"` + strings.Repeat(`}`, 1000),
			expectedSafe: false,
			description:  "Excessively nested JSON should be rejected",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isSafe := suite.utils.IsJSONSafe(tc.input)
			suite.Equal(tc.expectedSafe, isSafe, tc.description)

			if tc.expectedSafe {
				// Ensure valid JSON can be parsed
				var parsed interface{}
				err := json.Unmarshal([]byte(tc.input), &parsed)
				suite.NoError(err, "Safe JSON should be parseable")
			}
		})
	}
}

// TestURLSanitization tests URL sanitization and validation
func (suite *DataSanitizationTestSuite) TestURLSanitization() {
	testCases := []struct {
		name         string
		input        string
		expectedSafe bool
		description  string
	}{
		{
			name:         "Valid HTTPS URL",
			input:        "https://example.com/path",
			expectedSafe: true,
			description:  "Valid HTTPS URLs should be safe",
		},
		{
			name:         "Valid HTTP URL",
			input:        "http://example.com/path",
			expectedSafe: true,
			description:  "Valid HTTP URLs should be safe",
		},
		{
			name:         "JavaScript URL",
			input:        "javascript:alert('XSS')",
			expectedSafe: false,
			description:  "JavaScript URLs should be rejected",
		},
		{
			name:         "Data URL with HTML",
			input:        "data:text/html,<script>alert(1)</script>",
			expectedSafe: false,
			description:  "Data URLs with HTML should be rejected",
		},
		{
			name:         "FTP URL",
			input:        "ftp://files.example.com/file.txt",
			expectedSafe: false,
			description:  "FTP URLs should be rejected for security",
		},
		{
			name:         "File URL",
			input:        "file:///etc/passwd",
			expectedSafe: false,
			description:  "File URLs should be rejected",
		},
		{
			name:         "URL with credentials",
			input:        "https://user:pass@example.com/",
			expectedSafe: false,
			description:  "URLs with credentials should be rejected",
		},
		{
			name:         "Localhost URL",
			input:        "http://localhost:8080/admin",
			expectedSafe: false,
			description:  "Localhost URLs should be rejected",
		},
		{
			name:         "Private IP URL",
			input:        "http://192.168.1.1/admin",
			expectedSafe: false,
			description:  "Private IP URLs should be rejected",
		},
		{
			name:         "URL with encoded JavaScript",
			input:        "http://example.com/%6A%61%76%61%73%63%72%69%70%74%3A%61%6C%65%72%74%28%31%29",
			expectedSafe: false,
			description:  "URL-encoded JavaScript should be detected and rejected",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isSafe := suite.utils.IsURLSafe(tc.input)
			suite.Equal(tc.expectedSafe, isSafe, tc.description)

			if tc.expectedSafe {
				// Ensure valid URLs can be parsed
				parsed, err := url.Parse(tc.input)
				suite.NoError(err, "Safe URLs should be parseable")
				suite.NotEmpty(parsed.Host, "Safe URLs should have a host")
			}
		})
	}
}

// TestFilenameSanitization tests filename sanitization
func (suite *DataSanitizationTestSuite) TestFilenameSanitization() {
	testCases := []struct {
		name         string
		input        string
		expectedSafe bool
		sanitized    string
		description  string
	}{
		{
			name:         "Normal filename",
			input:        "document.pdf",
			expectedSafe: true,
			sanitized:    "document.pdf",
			description:  "Normal filenames should be safe",
		},
		{
			name:         "Filename with spaces",
			input:        "my document.pdf",
			expectedSafe: true,
			sanitized:    "my_document.pdf",
			description:  "Spaces should be replaced with underscores",
		},
		{
			name:         "Path traversal attempt",
			input:        "../../../etc/passwd",
			expectedSafe: false,
			sanitized:    "passwd",
			description:  "Path traversal should be removed",
		},
		{
			name:         "Windows path traversal",
			input:        "..\\..\\windows\\system32\\config\\sam",
			expectedSafe: false,
			sanitized:    "sam",
			description:  "Windows path traversal should be removed",
		},
		{
			name:         "Filename with null bytes",
			input:        "file.pdf\x00.exe",
			expectedSafe: false,
			sanitized:    "file.pdf_.exe",
			description:  "Null bytes should be replaced",
		},
		{
			name:         "Long filename",
			input:        strings.Repeat("a", 300) + ".txt",
			expectedSafe: false,
			sanitized:    strings.Repeat("a", 251) + ".txt",
			description:  "Long filenames should be truncated",
		},
		{
			name:         "Filename with special characters",
			input:        "file<>|:*?\"\\/.txt",
			expectedSafe: false,
			sanitized:    "file_________.txt",
			description:  "Special characters should be replaced",
		},
		{
			name:         "Hidden file",
			input:        ".hidden",
			expectedSafe: false,
			sanitized:    "_hidden",
			description:  "Hidden files should be prefixed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isSafe := suite.utils.IsFilenameSafe(tc.input)
			sanitized := suite.utils.SanitizeFilename(tc.input)

			suite.Equal(tc.expectedSafe, isSafe, tc.description)
			suite.Equal(tc.sanitized, sanitized, "Sanitized filename should match expected")

			// Ensure sanitized filename is safe
			suite.True(suite.utils.IsFilenameSafe(sanitized), "Sanitized filename should be safe")
			suite.False(strings.Contains(sanitized, ".."), "No path traversal in sanitized filename")
			suite.False(strings.Contains(sanitized, "\x00"), "No null bytes in sanitized filename")
			suite.True(len(sanitized) <= 255, "Sanitized filename should not exceed 255 characters")
		})
	}
}

// TestDataTypeSanitization tests data type validation and sanitization
func (suite *DataSanitizationTestSuite) TestDataTypeSanitization() {
	testCases := []struct {
		name         string
		dataType     string
		input        interface{}
		expectedSafe bool
		description  string
	}{
		{
			name:         "Valid integer",
			dataType:     "integer",
			input:        42,
			expectedSafe: true,
			description:  "Valid integers should be safe",
		},
		{
			name:         "String as integer",
			dataType:     "integer",
			input:        "42",
			expectedSafe: true,
			description:  "String representations of integers should be convertible",
		},
		{
			name:         "Invalid integer",
			dataType:     "integer",
			input:        "not a number",
			expectedSafe: false,
			description:  "Invalid integer strings should be rejected",
		},
		{
			name:         "Valid email",
			dataType:     "email",
			input:        "user@example.com",
			expectedSafe: true,
			description:  "Valid email addresses should be safe",
		},
		{
			name:         "Invalid email",
			dataType:     "email",
			input:        "not-an-email",
			expectedSafe: false,
			description:  "Invalid email formats should be rejected",
		},
		{
			name:         "Email with script",
			dataType:     "email",
			input:        "<script>alert('xss')</script>@example.com",
			expectedSafe: false,
			description:  "Emails with scripts should be rejected",
		},
		{
			name:         "Valid UUID",
			dataType:     "uuid",
			input:        "550e8400-e29b-41d4-a716-446655440000",
			expectedSafe: true,
			description:  "Valid UUIDs should be safe",
		},
		{
			name:         "Invalid UUID",
			dataType:     "uuid",
			input:        "not-a-uuid",
			expectedSafe: false,
			description:  "Invalid UUID formats should be rejected",
		},
		{
			name:         "Valid phone number",
			dataType:     "phone",
			input:        "+66812345678",
			expectedSafe: true,
			description:  "Valid phone numbers should be safe",
		},
		{
			name:         "Invalid phone number",
			dataType:     "phone",
			input:        "123-ABC-DEFG",
			expectedSafe: false,
			description:  "Invalid phone formats should be rejected",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isSafe := suite.utils.ValidateDataType(tc.dataType, tc.input)
			suite.Equal(tc.expectedSafe, isSafe, tc.description)
		})
	}
}

// TestContentFilteringSEA tests Southeast Asian content filtering
func (suite *DataSanitizationTestSuite) TestContentFilteringSEA() {
	testCases := []struct {
		name         string
		content      string
		country      string
		expectedSafe bool
		description  string
	}{
		{
			name:         "Normal Thai content",
			content:      "สวัสดีครับ ยินดีที่ได้รู้จัก",
			country:      "TH",
			expectedSafe: true,
			description:  "Normal Thai greetings should be safe",
		},
		{
			name:         "Normal Vietnamese content",
			content:      "Xin chào, rất vui được gặp bạn",
			country:      "VN",
			expectedSafe: true,
			description:  "Normal Vietnamese greetings should be safe",
		},
		{
			name:         "Normal Indonesian content",
			content:      "Halo, senang berkenalan dengan Anda",
			country:      "ID",
			expectedSafe: true,
			description:  "Normal Indonesian greetings should be safe",
		},
		{
			name:         "Mixed script with injection",
			content:      "สวัสดี <script>alert('XSS')</script>",
			country:      "TH",
			expectedSafe: false,
			description:  "Mixed script content should be flagged",
		},
		{
			name:         "Unicode normalization test",
			content:      "Café vs Café", // Different Unicode representations
			country:      "SG",
			expectedSafe: true,
			description:  "Unicode variations should be normalized safely",
		},
		{
			name:         "Extremely long content",
			content:      strings.Repeat("ก", 10000), // Thai character repeated
			country:      "TH",
			expectedSafe: false,
			description:  "Extremely long content should be flagged",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isSafe := suite.utils.IsSEAContentSafe(tc.content, tc.country)
			suite.Equal(tc.expectedSafe, isSafe, tc.description)

			if !tc.expectedSafe {
				sanitized := suite.utils.SanitizeSEAContent(tc.content, tc.country)
				suite.True(suite.utils.IsSEAContentSafe(sanitized, tc.country), "Sanitized content should be safe")
			}
		})
	}
}

// TestInputLengthValidation tests input length validation
func (suite *DataSanitizationTestSuite) TestInputLengthValidation() {
	testCases := []struct {
		name         string
		input        string
		maxLength    int
		expectedSafe bool
		description  string
	}{
		{
			name:         "Normal length input",
			input:        "Hello World",
			maxLength:    100,
			expectedSafe: true,
			description:  "Normal length inputs should be safe",
		},
		{
			name:         "Exactly at limit",
			input:        strings.Repeat("a", 100),
			maxLength:    100,
			expectedSafe: true,
			description:  "Inputs at exact limit should be safe",
		},
		{
			name:         "Exceeds limit",
			input:        strings.Repeat("a", 101),
			maxLength:    100,
			expectedSafe: false,
			description:  "Inputs exceeding limit should be rejected",
		},
		{
			name:         "Empty input",
			input:        "",
			maxLength:    100,
			expectedSafe: true,
			description:  "Empty inputs should be safe",
		},
		{
			name:         "Unicode character length",
			input:        strings.Repeat("🌟", 50), // Each emoji is multiple bytes
			maxLength:    50,
			expectedSafe: true,
			description:  "Unicode character count should be accurate",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isSafe := suite.utils.ValidateInputLength(tc.input, tc.maxLength)
			suite.Equal(tc.expectedSafe, isSafe, tc.description)
		})
	}
}

// TestSpecialCharacterHandling tests special character validation
func (suite *DataSanitizationTestSuite) TestSpecialCharacterHandling() {
	testCases := []struct {
		name         string
		input        string
		allowedChars string
		expectedSafe bool
		description  string
	}{
		{
			name:         "Alphanumeric only",
			input:        "abc123",
			allowedChars: "alphanumeric",
			expectedSafe: true,
			description:  "Alphanumeric characters should be safe",
		},
		{
			name:         "Special characters in alphanumeric",
			input:        "abc123!@#",
			allowedChars: "alphanumeric",
			expectedSafe: false,
			description:  "Special characters should be rejected in alphanumeric mode",
		},
		{
			name:         "Basic punctuation allowed",
			input:        "Hello, world! How are you?",
			allowedChars: "text",
			expectedSafe: true,
			description:  "Basic punctuation should be allowed in text mode",
		},
		{
			name:         "Dangerous characters in text",
			input:        "Hello <script>alert('xss')</script>",
			allowedChars: "text",
			expectedSafe: false,
			description:  "Dangerous characters should be rejected even in text mode",
		},
		{
			name:         "Control characters",
			input:        "Hello\x00\x01\x02World",
			allowedChars: "text",
			expectedSafe: false,
			description:  "Control characters should be rejected",
		},
		{
			name:         "Unicode normalization",
			input:        "café", // Composed vs decomposed forms
			allowedChars: "text",
			expectedSafe: true,
			description:  "Unicode characters should be normalized properly",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isSafe := suite.utils.ValidateSpecialCharacters(tc.input, tc.allowedChars)
			suite.Equal(tc.expectedSafe, isSafe, tc.description)

			if !tc.expectedSafe {
				sanitized := suite.utils.SanitizeSpecialCharacters(tc.input, tc.allowedChars)
				suite.True(suite.utils.ValidateSpecialCharacters(sanitized, tc.allowedChars), "Sanitized input should be safe")
			}
		})
	}
}

// TestEncodingValidation tests character encoding validation
func (suite *DataSanitizationTestSuite) TestEncodingValidation() {
	testCases := []struct {
		name         string
		input        []byte
		expectedSafe bool
		description  string
	}{
		{
			name:         "Valid UTF-8",
			input:        []byte("Hello, 世界!"),
			expectedSafe: true,
			description:  "Valid UTF-8 should be safe",
		},
		{
			name:         "Invalid UTF-8 sequence",
			input:        []byte{0xff, 0xfe, 0xfd},
			expectedSafe: false,
			description:  "Invalid UTF-8 sequences should be rejected",
		},
		{
			name:         "UTF-8 BOM",
			input:        []byte{0xef, 0xbb, 0xbf, 'H', 'e', 'l', 'l', 'o'},
			expectedSafe: true,
			description:  "UTF-8 with BOM should be handled correctly",
		},
		{
			name:         "Null bytes in UTF-8",
			input:        []byte{'H', 'e', 'l', 'l', 'o', 0x00, 'W', 'o', 'r', 'l', 'd'},
			expectedSafe: false,
			description:  "Null bytes should be rejected",
		},
		{
			name:         "Overlong UTF-8 encoding",
			input:        []byte{0xc0, 0xaf}, // Overlong encoding of '/'
			expectedSafe: false,
			description:  "Overlong UTF-8 encodings should be rejected",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			isSafe := suite.utils.ValidateUTF8Encoding(tc.input)
			suite.Equal(tc.expectedSafe, isSafe, tc.description)
		})
	}
}

// TestDataSanitizationIntegration tests integration with existing fixtures
func (suite *DataSanitizationTestSuite) TestDataSanitizationIntegration() {
	countries := []string{"TH", "SG", "ID", "MY", "VN", "PH"}

	for _, country := range countries {
		suite.Run(fmt.Sprintf("Country_%s_Integration", country), func() {
			// Create test data using fixtures
			user := suite.fixtures.Users.BasicUser(country)
			content := suite.fixtures.Content.BasicContent("message", country)

			// Test user data sanitization
			suite.True(suite.utils.ValidateEmail(*user.Email), "User email should be valid")
			suite.True(suite.utils.ValidatePhoneNumber(*user.Phone, country), "User phone should be valid")
			suite.True(suite.utils.IsFilenameSafe(user.Name), "User name should be filename-safe")

			// Test content sanitization
			suite.True(suite.utils.IsSEAContentSafe(content.Title, country), "Content title should be safe")
			suite.True(suite.utils.IsHTMLSafe(content.Body), "Content body should be HTML-safe")

			// Test against common attack vectors
			attackVectors := []string{
				"<script>alert('XSS')</script>",
				"'; DROP TABLE users; --",
				"javascript:alert(1)",
				"data:text/html,<script>alert(1)</script>",
			}

			for _, attack := range attackVectors {
				suite.False(suite.utils.IsHTMLSafe(attack), "Attack vector should be detected: %s", attack)
				suite.False(suite.utils.IsURLSafe(attack), "Attack vector should be rejected as URL: %s", attack)
				suite.False(suite.utils.IsSEAContentSafe(attack, country), "Attack vector should be rejected as content: %s", attack)
			}
		})
	}
}

// TestSanitizationPerformance tests sanitization performance with large inputs
func (suite *DataSanitizationTestSuite) TestSanitizationPerformance() {
	// Large input for performance testing
	largeInput := strings.Repeat("Hello World! ", 10000) // ~120KB

	suite.Run("HTML_Sanitization_Performance", func() {
		// Should complete within reasonable time
		sanitized := suite.utils.SanitizeHTML(largeInput)
		suite.NotEmpty(sanitized, "Large input should be sanitized")
		suite.True(suite.utils.IsHTMLSafe(sanitized), "Large sanitized input should be safe")
	})

	suite.Run("JSON_Validation_Performance", func() {
		// Create large valid JSON
		largeJSON := `{"data": "` + strings.Repeat("a", 50000) + `"}`
		isSafe := suite.utils.IsJSONSafe(largeJSON)
		suite.True(isSafe, "Large valid JSON should be processed efficiently")
	})

	suite.Run("URL_Validation_Performance", func() {
		// Test with various URL lengths
		baseURL := "https://example.com/"
		longPath := strings.Repeat("path/", 1000)
		longURL := baseURL + longPath

		isSafe := suite.utils.IsURLSafe(longURL)
		suite.False(isSafe, "Excessively long URLs should be rejected")
	})
}

func TestDataSanitizationSuite(t *testing.T) {
	suite.Run(t, new(DataSanitizationTestSuite))
}