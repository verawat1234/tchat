package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}

// CreatedResponse sends a created response
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, status int, message string, details string) {
	c.JSON(status, gin.H{
		"status":  "error",
		"message": message,
		"details": details,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": "Validation failed",
		"details": err.Error(),
	})
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  "error",
		"message": message,
	})
}