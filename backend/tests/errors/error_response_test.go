package errors_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"tchat.dev/tests/fixtures"
)

// ErrorResponseTestSuite provides comprehensive error response standardization testing
// for consistent error handling across all Tchat microservices
type ErrorResponseTestSuite struct {
	suite.Suite
	fixtures *fixtures.MasterFixtures
	ctx      context.Context
}

// StandardErrorResponse represents the standardized error format
type StandardErrorResponse struct {
	Error   ErrorDetails `json:"error"`
	TraceID string       `json:"trace_id,omitempty"`
	Path    string       `json:"path,omitempty"`
	Method  string       `json:"method,omitempty"`
	Time    time.Time    `json:"timestamp"`
}

// ErrorDetails contains detailed error information
type ErrorDetails struct {
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	Type        string                 `json:"type"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Validation  []ValidationError      `json:"validation,omitempty"`
	Suggestion  string                 `json:"suggestion,omitempty"`
	HelpURL     string                 `json:"help_url,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	ServiceName string                 `json:"service_name,omitempty"`
}

// ValidationError represents field-specific validation errors
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value,omitempty"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Path    string      `json:"path,omitempty"`
}

// HTTPErrorCode represents standard HTTP error codes with business context
type HTTPErrorCode struct {
	HTTPStatus int
	Code       string
	Type       string
	Message    string
	Category   string
}

// Southeast Asian localized error messages
type LocalizedErrorMessages struct {
	Code     string            `json:"code"`
	Messages map[string]string `json:"messages"` // country code -> localized message
}

// SetupSuite initializes the test suite
func (suite *ErrorResponseTestSuite) SetupSuite() {
	suite.fixtures = fixtures.NewMasterFixtures(12345)
	suite.ctx = context.Background()
}

// TestStandardErrorResponseStructure tests the standard error response format
func (suite *ErrorResponseTestSuite) TestStandardErrorResponseStructure() {
	testCases := []struct {
		name           string
		errorCode      string
		httpStatus     int
		errorType      string
		message        string
		hasDetails     bool
		hasValidation  bool
		hasSuggestion  bool
		description    string
	}{
		{
			name:          "Bad Request with validation",
			errorCode:     "VALIDATION_ERROR",
			httpStatus:    400,
			errorType:     "validation_error",
			message:       "Input validation failed",
			hasDetails:    true,
			hasValidation: true,
			hasSuggestion: true,
			description:   "Standard validation error format",
		},
		{
			name:          "Authentication error",
			errorCode:     "AUTHENTICATION_FAILED",
			httpStatus:    401,
			errorType:     "authentication_error",
			message:       "Authentication credentials are invalid",
			hasDetails:    false,
			hasValidation: false,
			hasSuggestion: true,
			description:   "Authentication failure format",
		},
		{
			name:          "Authorization error",
			errorCode:     "AUTHORIZATION_FAILED",
			httpStatus:    403,
			errorType:     "authorization_error",
			message:       "Insufficient permissions for this resource",
			hasDetails:    true,
			hasValidation: false,
			hasSuggestion: true,
			description:   "Authorization failure format",
		},
		{
			name:          "Resource not found",
			errorCode:     "RESOURCE_NOT_FOUND",
			httpStatus:    404,
			errorType:     "not_found_error",
			message:       "The requested resource was not found",
			hasDetails:    true,
			hasValidation: false,
			hasSuggestion: true,
			description:   "Not found error format",
		},
		{
			name:          "Conflict error",
			errorCode:     "RESOURCE_CONFLICT",
			httpStatus:    409,
			errorType:     "conflict_error",
			message:       "Resource conflict detected",
			hasDetails:    true,
			hasValidation: false,
			hasSuggestion: true,
			description:   "Conflict error format",
		},
		{
			name:          "Rate limit exceeded",
			errorCode:     "RATE_LIMIT_EXCEEDED",
			httpStatus:    429,
			errorType:     "rate_limit_error",
			message:       "Rate limit exceeded",
			hasDetails:    true,
			hasValidation: false,
			hasSuggestion: true,
			description:   "Rate limiting error format",
		},
		{
			name:          "Internal server error",
			errorCode:     "INTERNAL_SERVER_ERROR",
			httpStatus:    500,
			errorType:     "internal_error",
			message:       "An internal server error occurred",
			hasDetails:    false,
			hasValidation: false,
			hasSuggestion: false,
			description:   "Internal error format",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create standard error response
			errorResponse := suite.createStandardErrorResponse(
				tc.errorCode,
				tc.httpStatus,
				tc.errorType,
				tc.message,
				tc.hasDetails,
				tc.hasValidation,
				tc.hasSuggestion,
			)

			// Validate structure
			suite.validateErrorResponseStructure(errorResponse, tc)

			// Test JSON serialization/deserialization
			suite.testJSONSerialization(errorResponse)

			// Validate HTTP status mapping
			suite.Equal(tc.httpStatus, suite.getHTTPStatusFromErrorCode(tc.errorCode), tc.description)
		})
	}
}

// TestValidationErrorDetails tests detailed validation error format
func (suite *ErrorResponseTestSuite) TestValidationErrorDetails() {
	testCases := []struct {
		name        string
		field       string
		value       interface{}
		code        string
		message     string
		path        string
		description string
	}{
		{
			name:        "Required field missing",
			field:       "email",
			value:       nil,
			code:        "REQUIRED",
			message:     "Email is required",
			path:        "user.email",
			description: "Missing required field validation",
		},
		{
			name:        "Invalid email format",
			field:       "email",
			value:       "invalid-email",
			code:        "INVALID_FORMAT",
			message:     "Email format is invalid",
			path:        "user.email",
			description: "Email format validation",
		},
		{
			name:        "String too long",
			field:       "name",
			value:       strings.Repeat("a", 256),
			code:        "TOO_LONG",
			message:     "Name must be less than 255 characters",
			path:        "user.name",
			description: "String length validation",
		},
		{
			name:        "Invalid phone number",
			field:       "phone",
			value:       "123-invalid",
			code:        "INVALID_PHONE",
			message:     "Phone number format is invalid for this country",
			path:        "user.phone",
			description: "Phone number validation",
		},
		{
			name:        "Invalid UUID",
			field:       "id",
			value:       "not-a-uuid",
			code:        "INVALID_UUID",
			message:     "ID must be a valid UUID",
			path:        "user.id",
			description: "UUID format validation",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			validationError := ValidationError{
				Field:   tc.field,
				Value:   tc.value,
				Code:    tc.code,
				Message: tc.message,
				Path:    tc.path,
			}

			// Validate validation error structure
			suite.Equal(tc.field, validationError.Field, "Field should match")
			suite.Equal(tc.value, validationError.Value, "Value should match")
			suite.Equal(tc.code, validationError.Code, "Code should match")
			suite.Equal(tc.message, validationError.Message, "Message should match")
			suite.Equal(tc.path, validationError.Path, "Path should match")

			// Test validation error in error response
			errorResponse := StandardErrorResponse{
				Error: ErrorDetails{
					Code:    "VALIDATION_ERROR",
					Message: "Input validation failed",
					Type:    "validation_error",
					Validation: []ValidationError{validationError},
				},
				Time: time.Now(),
			}

			suite.Len(errorResponse.Error.Validation, 1, "Should have one validation error")
			suite.Equal(tc.field, errorResponse.Error.Validation[0].Field, tc.description)
		})
	}
}

// TestHTTPStatusCodeMapping tests proper HTTP status code mapping
func (suite *ErrorResponseTestSuite) TestHTTPStatusCodeMapping() {
	statusMappings := map[string]HTTPErrorCode{
		"VALIDATION_ERROR": {
			HTTPStatus: 400,
			Code:       "VALIDATION_ERROR",
			Type:       "validation_error",
			Message:    "Input validation failed",
			Category:   "client_error",
		},
		"AUTHENTICATION_FAILED": {
			HTTPStatus: 401,
			Code:       "AUTHENTICATION_FAILED",
			Type:       "authentication_error",
			Message:    "Authentication credentials are invalid",
			Category:   "client_error",
		},
		"AUTHORIZATION_FAILED": {
			HTTPStatus: 403,
			Code:       "AUTHORIZATION_FAILED",
			Type:       "authorization_error",
			Message:    "Insufficient permissions for this resource",
			Category:   "client_error",
		},
		"RESOURCE_NOT_FOUND": {
			HTTPStatus: 404,
			Code:       "RESOURCE_NOT_FOUND",
			Type:       "not_found_error",
			Message:    "The requested resource was not found",
			Category:   "client_error",
		},
		"METHOD_NOT_ALLOWED": {
			HTTPStatus: 405,
			Code:       "METHOD_NOT_ALLOWED",
			Type:       "method_error",
			Message:    "HTTP method not allowed for this endpoint",
			Category:   "client_error",
		},
		"RESOURCE_CONFLICT": {
			HTTPStatus: 409,
			Code:       "RESOURCE_CONFLICT",
			Type:       "conflict_error",
			Message:    "Resource conflict detected",
			Category:   "client_error",
		},
		"UNPROCESSABLE_ENTITY": {
			HTTPStatus: 422,
			Code:       "UNPROCESSABLE_ENTITY",
			Type:       "processing_error",
			Message:    "Request cannot be processed",
			Category:   "client_error",
		},
		"RATE_LIMIT_EXCEEDED": {
			HTTPStatus: 429,
			Code:       "RATE_LIMIT_EXCEEDED",
			Type:       "rate_limit_error",
			Message:    "Rate limit exceeded",
			Category:   "client_error",
		},
		"INTERNAL_SERVER_ERROR": {
			HTTPStatus: 500,
			Code:       "INTERNAL_SERVER_ERROR",
			Type:       "internal_error",
			Message:    "An internal server error occurred",
			Category:   "server_error",
		},
		"SERVICE_UNAVAILABLE": {
			HTTPStatus: 503,
			Code:       "SERVICE_UNAVAILABLE",
			Type:       "availability_error",
			Message:    "Service temporarily unavailable",
			Category:   "server_error",
		},
		"GATEWAY_TIMEOUT": {
			HTTPStatus: 504,
			Code:       "GATEWAY_TIMEOUT",
			Type:       "timeout_error",
			Message:    "Gateway timeout occurred",
			Category:   "server_error",
		},
	}

	for code, expectedMapping := range statusMappings {
		suite.Run(fmt.Sprintf("HTTP_Status_%s", code), func() {
			// Test HTTP status mapping
			actualStatus := suite.getHTTPStatusFromErrorCode(code)
			suite.Equal(expectedMapping.HTTPStatus, actualStatus, "HTTP status should match for code %s", code)

			// Test error categorization
			category := suite.getErrorCategory(expectedMapping.HTTPStatus)
			suite.Equal(expectedMapping.Category, category, "Error category should match for status %d", expectedMapping.HTTPStatus)

			// Test error response creation
			errorResponse := suite.createStandardErrorResponse(
				code,
				expectedMapping.HTTPStatus,
				expectedMapping.Type,
				expectedMapping.Message,
				false,
				false,
				false,
			)

			suite.Equal(code, errorResponse.Error.Code)
			suite.Equal(expectedMapping.Type, errorResponse.Error.Type)
			suite.Equal(expectedMapping.Message, errorResponse.Error.Message)
		})
	}
}

// TestSoutheastAsianLocalization tests localized error messages
func (suite *ErrorResponseTestSuite) TestSoutheastAsianLocalization() {
	localizedMessages := map[string]LocalizedErrorMessages{
		"VALIDATION_ERROR": {
			Code: "VALIDATION_ERROR",
			Messages: map[string]string{
				"TH": "ข้อมูลที่ป้อนไม่ถูกต้อง",
				"SG": "Input validation failed",
				"ID": "Validasi input gagal",
				"MY": "Pengesahan input gagal",
				"VN": "Xác thực đầu vào thất bại",
				"PH": "Input validation failed",
				"EN": "Input validation failed",
			},
		},
		"AUTHENTICATION_FAILED": {
			Code: "AUTHENTICATION_FAILED",
			Messages: map[string]string{
				"TH": "การยืนยันตัวตนล้มเหลว",
				"SG": "Authentication failed",
				"ID": "Autentikasi gagal",
				"MY": "Pengesahan gagal",
				"VN": "Xác thực thất bại",
				"PH": "Authentication failed",
				"EN": "Authentication failed",
			},
		},
		"RESOURCE_NOT_FOUND": {
			Code: "RESOURCE_NOT_FOUND",
			Messages: map[string]string{
				"TH": "ไม่พบข้อมูลที่ร้องขอ",
				"SG": "Resource not found",
				"ID": "Sumber daya tidak ditemukan",
				"MY": "Sumber tidak dijumpai",
				"VN": "Không tìm thấy tài nguyên",
				"PH": "Resource not found",
				"EN": "Resource not found",
			},
		},
		"RATE_LIMIT_EXCEEDED": {
			Code: "RATE_LIMIT_EXCEEDED",
			Messages: map[string]string{
				"TH": "เกินขดกำหนดการเรียกใช้งาน",
				"SG": "Rate limit exceeded",
				"ID": "Batas tingkat terlampaui",
				"MY": "Had kadar melebihi",
				"VN": "Vượt quá giới hạn tốc độ",
				"PH": "Rate limit exceeded",
				"EN": "Rate limit exceeded",
			},
		},
	}

	countries := []string{"TH", "SG", "ID", "MY", "VN", "PH"}

	for errorCode, localizedMsg := range localizedMessages {
		for _, country := range countries {
			suite.Run(fmt.Sprintf("%s_%s", errorCode, country), func() {
				// Test localized message exists
				message, exists := localizedMsg.Messages[country]
				suite.True(exists, "Localized message should exist for country %s", country)
				suite.NotEmpty(message, "Localized message should not be empty for country %s", country)

				// Test localized error response
				errorResponse := suite.createLocalizedErrorResponse(errorCode, country, message)

				suite.Equal(errorCode, errorResponse.Error.Code)
				suite.Equal(message, errorResponse.Error.Message)
				suite.NotEmpty(errorResponse.Error.ServiceName)

				// Test JSON serialization with localized content
				jsonBytes, err := json.Marshal(errorResponse)
				suite.NoError(err, "Should serialize localized error response")

				var deserializedResponse StandardErrorResponse
				err = json.Unmarshal(jsonBytes, &deserializedResponse)
				suite.NoError(err, "Should deserialize localized error response")

				suite.Equal(message, deserializedResponse.Error.Message, "Localized message should be preserved")
			})
		}
	}
}

// TestErrorResponseWithDetails tests error responses with additional details
func (suite *ErrorResponseTestSuite) TestErrorResponseWithDetails() {
	testCases := []struct {
		name        string
		errorCode   string
		details     map[string]interface{}
		description string
	}{
		{
			name:      "Rate limit with details",
			errorCode: "RATE_LIMIT_EXCEEDED",
			details: map[string]interface{}{
				"limit":         100,
				"window":        "1m",
				"retry_after":   60,
				"current_usage": 101,
			},
			description: "Rate limit error with timing details",
		},
		{
			name:      "Validation with field details",
			errorCode: "VALIDATION_ERROR",
			details: map[string]interface{}{
				"failed_fields": []string{"email", "phone"},
				"total_errors":  2,
				"error_count":   map[string]int{"email": 1, "phone": 1},
			},
			description: "Validation error with field breakdown",
		},
		{
			name:      "Authentication with context",
			errorCode: "AUTHENTICATION_FAILED",
			details: map[string]interface{}{
				"reason":        "token_expired",
				"expired_at":    time.Now().Add(-time.Hour).Unix(),
				"token_type":    "jwt",
				"refresh_token": true,
			},
			description: "Authentication error with context details",
		},
		{
			name:      "Resource conflict with resolution",
			errorCode: "RESOURCE_CONFLICT",
			details: map[string]interface{}{
				"conflict_type":   "unique_constraint",
				"conflicting_field": "email",
				"existing_id":     "550e8400-e29b-41d4-a716-446655440000",
				"resolution_hint": "use_different_email",
			},
			description: "Conflict error with resolution guidance",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			errorResponse := StandardErrorResponse{
				Error: ErrorDetails{
					Code:        tc.errorCode,
					Type:        suite.getErrorTypeFromCode(tc.errorCode),
					Message:     suite.getDefaultMessageFromCode(tc.errorCode),
					Details:     tc.details,
					ServiceName: "tchat-test-service",
				},
				TraceID: suite.generateTraceID(),
				Path:    "/api/test",
				Method:  "POST",
				Time:    time.Now(),
			}

			// Validate details are present
			suite.NotNil(errorResponse.Error.Details, "Details should be present")
			suite.Equal(tc.details, errorResponse.Error.Details, tc.description)

			// Test JSON serialization preserves details
			jsonBytes, err := json.Marshal(errorResponse)
			suite.NoError(err, "Should serialize error response with details")

			var deserializedResponse StandardErrorResponse
			err = json.Unmarshal(jsonBytes, &deserializedResponse)
			suite.NoError(err, "Should deserialize error response with details")

			suite.Equal(tc.details, deserializedResponse.Error.Details, "Details should be preserved in JSON")
		})
	}
}

// TestErrorResponseConsistency tests consistency across different services
func (suite *ErrorResponseTestSuite) TestErrorResponseConsistency() {
	services := []string{
		"tchat-auth-service",
		"tchat-content-service",
		"tchat-messaging-service",
		"tchat-payment-service",
		"tchat-notification-service",
		"tchat-commerce-service",
	}

	for _, service := range services {
		suite.Run(fmt.Sprintf("Service_%s", service), func() {
			// Test that each service returns consistent error format
			errorResponse := StandardErrorResponse{
				Error: ErrorDetails{
					Code:        "VALIDATION_ERROR",
					Type:        "validation_error",
					Message:     "Input validation failed",
					ServiceName: service,
					RequestID:   suite.generateRequestID(),
				},
				TraceID: suite.generateTraceID(),
				Path:    "/api/test",
				Method:  "POST",
				Time:    time.Now(),
			}

			// Validate service consistency
			suite.Equal(service, errorResponse.Error.ServiceName, "Service name should match")
			suite.NotEmpty(errorResponse.Error.RequestID, "Request ID should be present")
			suite.NotEmpty(errorResponse.TraceID, "Trace ID should be present")

			// Test standard fields are present
			suite.NotEmpty(errorResponse.Error.Code, "Error code should be present")
			suite.NotEmpty(errorResponse.Error.Type, "Error type should be present")
			suite.NotEmpty(errorResponse.Error.Message, "Error message should be present")
			suite.False(errorResponse.Time.IsZero(), "Timestamp should be present")

			// Test JSON structure consistency
			jsonBytes, err := json.Marshal(errorResponse)
			suite.NoError(err, "Should serialize consistently")

			var deserializedResponse StandardErrorResponse
			err = json.Unmarshal(jsonBytes, &deserializedResponse)
			suite.NoError(err, "Should deserialize consistently")

			suite.Equal(errorResponse.Error.ServiceName, deserializedResponse.Error.ServiceName)
		})
	}
}

// TestErrorResponseSuggestions tests error responses with helpful suggestions
func (suite *ErrorResponseTestSuite) TestErrorResponseSuggestions() {
	testCases := []struct {
		name        string
		errorCode   string
		suggestion  string
		helpURL     string
		description string
	}{
		{
			name:        "Authentication failure suggestion",
			errorCode:   "AUTHENTICATION_FAILED",
			suggestion:  "Please check your credentials and try again. If you forgot your password, use the password reset feature.",
			helpURL:     "https://docs.tchat.dev/authentication",
			description: "Authentication error with helpful suggestion",
		},
		{
			name:        "Rate limit suggestion",
			errorCode:   "RATE_LIMIT_EXCEEDED",
			suggestion:  "Please wait before making another request. Consider upgrading your plan for higher limits.",
			helpURL:     "https://docs.tchat.dev/rate-limits",
			description: "Rate limit error with upgrade suggestion",
		},
		{
			name:        "Validation error suggestion",
			errorCode:   "VALIDATION_ERROR",
			suggestion:  "Please check the highlighted fields and correct any errors before resubmitting.",
			helpURL:     "https://docs.tchat.dev/validation",
			description: "Validation error with correction guidance",
		},
		{
			name:        "Resource not found suggestion",
			errorCode:   "RESOURCE_NOT_FOUND",
			suggestion:  "The resource may have been deleted or moved. Please check the ID and try again.",
			helpURL:     "https://docs.tchat.dev/resources",
			description: "Not found error with troubleshooting help",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			errorResponse := StandardErrorResponse{
				Error: ErrorDetails{
					Code:        tc.errorCode,
					Type:        suite.getErrorTypeFromCode(tc.errorCode),
					Message:     suite.getDefaultMessageFromCode(tc.errorCode),
					Suggestion:  tc.suggestion,
					HelpURL:     tc.helpURL,
					ServiceName: "tchat-test-service",
				},
				Time: time.Now(),
			}

			// Validate suggestions are present
			suite.Equal(tc.suggestion, errorResponse.Error.Suggestion, tc.description)
			suite.Equal(tc.helpURL, errorResponse.Error.HelpURL, "Help URL should match")

			// Test that suggestions are helpful and actionable
			suite.True(len(tc.suggestion) > 20, "Suggestion should be meaningful")
			suite.True(strings.Contains(tc.helpURL, "https://"), "Help URL should be valid HTTPS")
			suite.True(strings.Contains(tc.helpURL, "tchat.dev"), "Help URL should point to official docs")

			// Test JSON serialization includes suggestions
			jsonBytes, err := json.Marshal(errorResponse)
			suite.NoError(err, "Should serialize error response with suggestions")

			var deserializedResponse StandardErrorResponse
			err = json.Unmarshal(jsonBytes, &deserializedResponse)
			suite.NoError(err, "Should deserialize error response with suggestions")

			suite.Equal(tc.suggestion, deserializedResponse.Error.Suggestion, "Suggestion should be preserved")
			suite.Equal(tc.helpURL, deserializedResponse.Error.HelpURL, "Help URL should be preserved")
		})
	}
}

// TestErrorResponseTraceability tests error tracing and debugging support
func (suite *ErrorResponseTestSuite) TestErrorResponseTraceability() {
	suite.Run("Trace_ID_Generation", func() {
		errorResponse1 := suite.createStandardErrorResponse("VALIDATION_ERROR", 400, "validation_error", "Test error", false, false, false)
		errorResponse2 := suite.createStandardErrorResponse("VALIDATION_ERROR", 400, "validation_error", "Test error", false, false, false)

		// Trace IDs should be unique
		suite.NotEqual(errorResponse1.TraceID, errorResponse2.TraceID, "Trace IDs should be unique")
		suite.NotEmpty(errorResponse1.TraceID, "Trace ID should not be empty")
		suite.NotEmpty(errorResponse2.TraceID, "Trace ID should not be empty")

		// Trace IDs should follow expected format (UUID-like)
		suite.Len(errorResponse1.TraceID, 36, "Trace ID should be UUID format")
		suite.Contains(errorResponse1.TraceID, "-", "Trace ID should contain hyphens")
	})

	suite.Run("Request_ID_Generation", func() {
		errorResponse := StandardErrorResponse{
			Error: ErrorDetails{
				Code:      "TEST_ERROR",
				Type:      "test_error",
				Message:   "Test error message",
				RequestID: suite.generateRequestID(),
			},
			TraceID: suite.generateTraceID(),
			Time:    time.Now(),
		}

		suite.NotEmpty(errorResponse.Error.RequestID, "Request ID should be present")
		suite.NotEqual(errorResponse.TraceID, errorResponse.Error.RequestID, "Trace ID and Request ID should be different")
	})

	suite.Run("Timestamp_Consistency", func() {
		beforeTime := time.Now()
		errorResponse := suite.createStandardErrorResponse("TEST_ERROR", 500, "test_error", "Test error", false, false, false)
		afterTime := time.Now()

		suite.True(errorResponse.Time.After(beforeTime) || errorResponse.Time.Equal(beforeTime), "Error timestamp should be after start time")
		suite.True(errorResponse.Time.Before(afterTime) || errorResponse.Time.Equal(afterTime), "Error timestamp should be before end time")
	})
}

// Helper methods

// createStandardErrorResponse creates a standard error response for testing
func (suite *ErrorResponseTestSuite) createStandardErrorResponse(
	code string,
	httpStatus int,
	errorType string,
	message string,
	hasDetails bool,
	hasValidation bool,
	hasSuggestion bool,
) StandardErrorResponse {
	errorDetails := ErrorDetails{
		Code:        code,
		Type:        errorType,
		Message:     message,
		ServiceName: "tchat-test-service",
		RequestID:   suite.generateRequestID(),
	}

	if hasDetails {
		errorDetails.Details = map[string]interface{}{
			"additional_info": "Test details",
			"error_source":    "test_suite",
		}
	}

	if hasValidation {
		errorDetails.Validation = []ValidationError{
			{
				Field:   "test_field",
				Value:   "invalid_value",
				Code:    "INVALID",
				Message: "Test field is invalid",
				Path:    "test.field",
			},
		}
	}

	if hasSuggestion {
		errorDetails.Suggestion = "This is a test suggestion for error resolution"
		errorDetails.HelpURL = "https://docs.tchat.dev/errors/" + strings.ToLower(code)
	}

	return StandardErrorResponse{
		Error:   errorDetails,
		TraceID: suite.generateTraceID(),
		Path:    "/api/test",
		Method:  "POST",
		Time:    time.Now(),
	}
}

// createLocalizedErrorResponse creates a localized error response
func (suite *ErrorResponseTestSuite) createLocalizedErrorResponse(code, country, message string) StandardErrorResponse {
	return StandardErrorResponse{
		Error: ErrorDetails{
			Code:        code,
			Type:        suite.getErrorTypeFromCode(code),
			Message:     message,
			ServiceName: "tchat-test-service",
			RequestID:   suite.generateRequestID(),
			Details: map[string]interface{}{
				"locale": country,
			},
		},
		TraceID: suite.generateTraceID(),
		Path:    "/api/test",
		Method:  "POST",
		Time:    time.Now(),
	}
}

// validateErrorResponseStructure validates the error response structure
func (suite *ErrorResponseTestSuite) validateErrorResponseStructure(errorResponse StandardErrorResponse, tc struct {
	name           string
	errorCode      string
	httpStatus     int
	errorType      string
	message        string
	hasDetails     bool
	hasValidation  bool
	hasSuggestion  bool
	description    string
}) {
	// Validate required fields
	suite.Equal(tc.errorCode, errorResponse.Error.Code, "Error code should match")
	suite.Equal(tc.errorType, errorResponse.Error.Type, "Error type should match")
	suite.Equal(tc.message, errorResponse.Error.Message, "Error message should match")
	suite.NotEmpty(errorResponse.TraceID, "Trace ID should be present")
	suite.False(errorResponse.Time.IsZero(), "Timestamp should be present")

	// Validate conditional fields
	if tc.hasDetails {
		suite.NotNil(errorResponse.Error.Details, "Details should be present when expected")
	}

	if tc.hasValidation {
		suite.NotEmpty(errorResponse.Error.Validation, "Validation errors should be present when expected")
	}

	if tc.hasSuggestion {
		suite.NotEmpty(errorResponse.Error.Suggestion, "Suggestion should be present when expected")
	}
}

// testJSONSerialization tests JSON serialization and deserialization
func (suite *ErrorResponseTestSuite) testJSONSerialization(errorResponse StandardErrorResponse) {
	// Test serialization
	jsonBytes, err := json.Marshal(errorResponse)
	suite.NoError(err, "Should serialize to JSON")
	suite.NotEmpty(jsonBytes, "JSON should not be empty")

	// Test deserialization
	var deserializedResponse StandardErrorResponse
	err = json.Unmarshal(jsonBytes, &deserializedResponse)
	suite.NoError(err, "Should deserialize from JSON")

	// Validate preserved fields
	suite.Equal(errorResponse.Error.Code, deserializedResponse.Error.Code)
	suite.Equal(errorResponse.Error.Type, deserializedResponse.Error.Type)
	suite.Equal(errorResponse.Error.Message, deserializedResponse.Error.Message)
	suite.Equal(errorResponse.TraceID, deserializedResponse.TraceID)
}

// getHTTPStatusFromErrorCode maps error codes to HTTP status codes
func (suite *ErrorResponseTestSuite) getHTTPStatusFromErrorCode(code string) int {
	statusMap := map[string]int{
		"VALIDATION_ERROR":        400,
		"AUTHENTICATION_FAILED":   401,
		"AUTHORIZATION_FAILED":    403,
		"RESOURCE_NOT_FOUND":      404,
		"METHOD_NOT_ALLOWED":      405,
		"RESOURCE_CONFLICT":       409,
		"UNPROCESSABLE_ENTITY":    422,
		"RATE_LIMIT_EXCEEDED":     429,
		"INTERNAL_SERVER_ERROR":   500,
		"SERVICE_UNAVAILABLE":     503,
		"GATEWAY_TIMEOUT":         504,
	}

	if status, exists := statusMap[code]; exists {
		return status
	}
	return 500 // Default to internal server error
}

// getErrorCategory categorizes errors based on HTTP status
func (suite *ErrorResponseTestSuite) getErrorCategory(httpStatus int) string {
	if httpStatus >= 400 && httpStatus < 500 {
		return "client_error"
	}
	if httpStatus >= 500 && httpStatus < 600 {
		return "server_error"
	}
	return "unknown"
}

// getErrorTypeFromCode maps error codes to error types
func (suite *ErrorResponseTestSuite) getErrorTypeFromCode(code string) string {
	typeMap := map[string]string{
		"VALIDATION_ERROR":        "validation_error",
		"AUTHENTICATION_FAILED":   "authentication_error",
		"AUTHORIZATION_FAILED":    "authorization_error",
		"RESOURCE_NOT_FOUND":      "not_found_error",
		"METHOD_NOT_ALLOWED":      "method_error",
		"RESOURCE_CONFLICT":       "conflict_error",
		"UNPROCESSABLE_ENTITY":    "processing_error",
		"RATE_LIMIT_EXCEEDED":     "rate_limit_error",
		"INTERNAL_SERVER_ERROR":   "internal_error",
		"SERVICE_UNAVAILABLE":     "availability_error",
		"GATEWAY_TIMEOUT":         "timeout_error",
	}

	if errorType, exists := typeMap[code]; exists {
		return errorType
	}
	return "unknown_error"
}

// getDefaultMessageFromCode provides default messages for error codes
func (suite *ErrorResponseTestSuite) getDefaultMessageFromCode(code string) string {
	messageMap := map[string]string{
		"VALIDATION_ERROR":        "Input validation failed",
		"AUTHENTICATION_FAILED":   "Authentication credentials are invalid",
		"AUTHORIZATION_FAILED":    "Insufficient permissions for this resource",
		"RESOURCE_NOT_FOUND":      "The requested resource was not found",
		"METHOD_NOT_ALLOWED":      "HTTP method not allowed for this endpoint",
		"RESOURCE_CONFLICT":       "Resource conflict detected",
		"UNPROCESSABLE_ENTITY":    "Request cannot be processed",
		"RATE_LIMIT_EXCEEDED":     "Rate limit exceeded",
		"INTERNAL_SERVER_ERROR":   "An internal server error occurred",
		"SERVICE_UNAVAILABLE":     "Service temporarily unavailable",
		"GATEWAY_TIMEOUT":         "Gateway timeout occurred",
	}

	if message, exists := messageMap[code]; exists {
		return message
	}
	return "An error occurred"
}

// generateTraceID generates a unique trace ID
func (suite *ErrorResponseTestSuite) generateTraceID() string {
	return fmt.Sprintf("trace-%d-%d", time.Now().UnixNano(), suite.fixtures.RandomInt(1000, 9999))
}

// generateRequestID generates a unique request ID
func (suite *ErrorResponseTestSuite) generateRequestID() string {
	return fmt.Sprintf("req-%d-%d", time.Now().UnixNano(), suite.fixtures.RandomInt(100, 999))
}

func TestErrorResponseSuite(t *testing.T) {
	suite.Run(t, new(ErrorResponseTestSuite))
}