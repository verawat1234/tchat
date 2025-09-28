package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BaseResponse represents the base response structure
type BaseResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// DataResponse represents a successful response with data
type DataResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool       `json:"success"`
	Error     *ErrorInfo `json:"error"`
	Timestamp string     `json:"timestamp"`
}

// SuccessResponse represents a simple success response
type SuccessResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// SendSuccessResponse sends a success response
func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success:   true,
		Message:   "Success",
		Data:      data,
		Timestamp: getCurrentTimestamp(),
	})
}

// SuccessMessageResponse sends a success response with custom message
func SuccessMessageResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success:   true,
		Message:   message,
		Timestamp: getCurrentTimestamp(),
	})
}

// DataResponse sends a data response
func SendDataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, DataResponse{
		Success:   true,
		Data:      data,
		Timestamp: getCurrentTimestamp(),
	})
}

// ErrorResponse sends an error response
func SendErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		Timestamp: getCurrentTimestamp(),
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	var validationErrors []map[string]string

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErr {
			validationErrors = append(validationErrors, map[string]string{
				"field":   fieldError.Field(),
				"tag":     fieldError.Tag(),
				"value":   fieldError.Param(),
				"message": getValidationErrorMessage(fieldError),
			})
		}
	}

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    "validation_error",
			Message: "Validation failed",
			Details: validationErrors,
		},
		Timestamp: getCurrentTimestamp(),
	})
}

// InternalErrorResponse sends an internal server error response
func InternalErrorResponse(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusInternalServerError, "internal_error", message)
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusUnauthorized, "unauthorized", message)
}

// ForbiddenResponse sends a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusForbidden, "forbidden", message)
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusNotFound, "not_found", message)
}

// ConflictResponse sends a conflict response
func ConflictResponse(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusConflict, "conflict", message)
}

// BadRequestResponse sends a bad request response
func BadRequestResponse(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusBadRequest, "bad_request", message)
}

// Helper functions

func getCurrentTimestamp() string {
	return "2024-09-22T18:30:00Z" // Simplified for now
}

func getValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "len":
		return "Invalid length"
	case "numeric":
		return "Must be numeric"
	case "e164":
		return "Invalid phone number format"
	default:
		return "Invalid value"
	}
}