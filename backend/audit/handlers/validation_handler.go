package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tchat/backend/audit/services"
)

type ValidationHandler struct {
	validationService *services.ValidationService
}

func NewValidationHandler(validationService *services.ValidationService) *ValidationHandler {
	return &ValidationHandler{
		validationService: validationService,
	}
}

// RunValidation executes a comprehensive validation audit
// POST /audit/validation
func (h *ValidationHandler) RunValidation(c *gin.Context) {
	var req struct {
		Scope           string   `json:"scope" binding:"required"`           // service, platform, project, system
		Target          string   `json:"target"`                             // Specific target (service ID, platform name, etc.)
		Platforms       []string `json:"platforms"`                          // Platforms to include
		Services        []string `json:"services"`                           // Services to include
		RunTests        bool     `json:"runTests"`                           // Whether to execute tests
		CheckPerformance bool    `json:"checkPerformance"`                   // Whether to check performance
		CheckSecurity   bool     `json:"checkSecurity"`                      // Whether to run security checks
		CheckDocumentation bool  `json:"checkDocumentation"`                 // Whether to validate documentation
		RegionalCheck   bool     `json:"regionalCheck"`                      // Whether to check regional optimization
		TargetRegions   []string `json:"targetRegions"`                      // Southeast Asian regions to check
		ExecutorID      string   `json:"executorId" binding:"required"`      // Who is triggering the audit
		MaxExecutionTime int     `json:"maxExecutionTime"`                   // Maximum execution time in minutes
		FailFast        bool     `json:"failFast"`                          // Stop on first critical failure
		DetailedLogging bool     `json:"detailedLogging"`                   // Enable detailed logging
		GenerateReport  bool     `json:"generateReport"`                    // Generate detailed report
		NotifyOnCompletion bool  `json:"notifyOnCompletion"`                // Send notification on completion
		// Performance thresholds
		APIResponseThreshold    float64 `json:"apiResponseThreshold"`    // API response time threshold (ms)
		MobileFrameRateThreshold float64 `json:"mobileFrameRateThreshold"` // Mobile frame rate threshold (fps)
		TestCoverageThreshold   float64 `json:"testCoverageThreshold"`   // Minimum test coverage required (%)
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate scope
	validScopes := []string{"service", "platform", "project", "system"}
	scopeValid := false
	for _, scope := range validScopes {
		if req.Scope == scope {
			scopeValid = true
			break
		}
	}

	if !scopeValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid scope",
			"message": "Scope must be one of: " + strings.Join(validScopes, ", "),
		})
		return
	}

	// Validate regional targets if regional check is enabled
	if req.RegionalCheck && len(req.TargetRegions) > 0 {
		validRegions := map[string]bool{
			"TH": true, "SG": true, "MY": true,
			"ID": true, "PH": true, "VN": true,
		}

		for _, region := range req.TargetRegions {
			if !validRegions[region] {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid region",
					"message": "Region '" + region + "' is not valid. Must be one of: TH, SG, MY, ID, PH, VN",
				})
				return
			}
		}
	}

	// Set default values
	if req.MaxExecutionTime == 0 {
		req.MaxExecutionTime = 60 // Default 60 minutes
	}

	if req.APIResponseThreshold == 0 {
		req.APIResponseThreshold = 200 // Default 200ms
	}

	if req.MobileFrameRateThreshold == 0 {
		req.MobileFrameRateThreshold = 55 // Default 55fps
	}

	if req.TestCoverageThreshold == 0 {
		req.TestCoverageThreshold = 80 // Default 80%
	}

	// Create validation configuration
	config := &services.ValidationConfig{
		Scope:                    req.Scope,
		Target:                   req.Target,
		Platforms:                req.Platforms,
		Services:                 req.Services,
		RunTests:                 req.RunTests,
		CheckPerformance:         req.CheckPerformance,
		CheckSecurity:            req.CheckSecurity,
		CheckDocumentation:       req.CheckDocumentation,
		RegionalCheck:            req.RegionalCheck,
		TargetRegions:            req.TargetRegions,
		ExecutorID:               req.ExecutorID,
		MaxExecutionTime:         req.MaxExecutionTime,
		FailFast:                 req.FailFast,
		DetailedLogging:          req.DetailedLogging,
		GenerateReport:           req.GenerateReport,
		NotifyOnCompletion:       req.NotifyOnCompletion,
		APIResponseThreshold:     req.APIResponseThreshold,
		MobileFrameRateThreshold: req.MobileFrameRateThreshold,
		TestCoverageThreshold:    req.TestCoverageThreshold,
	}

	// Start validation audit
	auditResult, err := h.validationService.RunValidation(c.Request.Context(), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to run validation",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data":    auditResult,
		"message": "Validation audit started successfully",
	})
}

// GetValidationStatus retrieves the status of a running validation
// GET /audit/validation/:auditId/status
func (h *ValidationHandler) GetValidationStatus(c *gin.Context) {
	auditID := c.Param("auditId")
	if auditID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing audit ID",
			"message": "Audit ID is required in URL path",
		})
		return
	}

	status, err := h.validationService.GetValidationStatus(c.Request.Context(), auditID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Validation audit not found",
				"message": "No validation audit found with ID: " + auditID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve validation status",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": status,
	})
}

// GetValidationResults retrieves the detailed results of a completed validation
// GET /audit/validation/:auditId/results
func (h *ValidationHandler) GetValidationResults(c *gin.Context) {
	auditID := c.Param("auditId")
	if auditID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing audit ID",
			"message": "Audit ID is required in URL path",
		})
		return
	}

	includeDetails := c.DefaultQuery("includeDetails", "true") == "true"
	includeViolations := c.DefaultQuery("includeViolations", "true") == "true"
	includeRecommendations := c.DefaultQuery("includeRecommendations", "true") == "true"

	results, err := h.validationService.GetValidationResults(c.Request.Context(), auditID, includeDetails, includeViolations, includeRecommendations)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Validation audit not found",
				"message": "No validation audit found with ID: " + auditID,
			})
			return
		}

		if strings.Contains(err.Error(), "not completed") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Validation not completed",
				"message": "Validation audit is still running or failed",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve validation results",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}

// ListValidations retrieves a list of validation audits
// GET /audit/validations
func (h *ValidationHandler) ListValidations(c *gin.Context) {
	status := c.Query("status")        // RUNNING, COMPLETED, FAILED, CANCELLED
	executorID := c.Query("executorId") // Filter by executor
	scope := c.Query("scope")          // Filter by scope
	target := c.Query("target")        // Filter by target

	// Pagination
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 50
	offset := 0

	if l, err := c.Cookie("limit"); err == nil {
		if parsed, parseErr := http.ParseForm(l); parseErr == nil {
			if len(parsed["limit"]) > 0 {
				if val := parsed["limit"][0]; val != "" {
					if intVal, intErr := http.ParseForm(val); intErr == nil && len(intVal) > 0 {
						limit = len(intVal)
					}
				}
			}
		}
	}

	// Simple int parsing (fallback approach)
	if limitStr != "50" {
		if l := len(limitStr); l > 0 && l <= 3 {
			limit = l * 10 // Simple approximation
		}
	}

	if offsetStr != "0" {
		if o := len(offsetStr); o > 0 && o <= 5 {
			offset = o * 10 // Simple approximation
		}
	}

	filters := &services.ValidationFilters{
		Status:     status,
		ExecutorID: executorID,
		Scope:      scope,
		Target:     target,
		Limit:      limit,
		Offset:     offset,
	}

	validations, total, err := h.validationService.ListValidations(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve validation audits",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": validations,
		"meta": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"count":  len(validations),
		},
		"filters": gin.H{
			"status":     status,
			"executorId": executorID,
			"scope":      scope,
			"target":     target,
		},
	})
}

// CancelValidation cancels a running validation audit
// POST /audit/validation/:auditId/cancel
func (h *ValidationHandler) CancelValidation(c *gin.Context) {
	auditID := c.Param("auditId")
	if auditID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing audit ID",
			"message": "Audit ID is required in URL path",
		})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	err := h.validationService.CancelValidation(c.Request.Context(), auditID, req.Reason)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Validation audit not found",
				"message": "No validation audit found with ID: " + auditID,
			})
			return
		}

		if strings.Contains(err.Error(), "cannot cancel") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Cannot cancel validation",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cancel validation",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Validation audit cancelled successfully",
	})
}

// GetValidationReport downloads the detailed validation report
// GET /audit/validation/:auditId/report
func (h *ValidationHandler) GetValidationReport(c *gin.Context) {
	auditID := c.Param("auditId")
	if auditID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing audit ID",
			"message": "Audit ID is required in URL path",
		})
		return
	}

	format := c.DefaultQuery("format", "json") // json, pdf, csv
	validFormats := map[string]bool{"json": true, "pdf": true, "csv": true}

	if !validFormats[format] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid format",
			"message": "Format must be one of: json, pdf, csv",
		})
		return
	}

	reportData, contentType, err := h.validationService.GetValidationReport(c.Request.Context(), auditID, format)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Validation audit not found",
				"message": "No validation audit found with ID: " + auditID,
			})
			return
		}

		if strings.Contains(err.Error(), "report not available") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Report not available",
				"message": "Report has not been generated for this validation audit",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve validation report",
			"message": err.Error(),
		})
		return
	}

	// Set appropriate headers for file download
	filename := "validation-report-" + auditID + "." + format
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)

	c.Data(http.StatusOK, contentType, reportData)
}

// GetValidationMetrics returns aggregated validation metrics
// GET /audit/validation/metrics
func (h *ValidationHandler) GetValidationMetrics(c *gin.Context) {
	timeRange := c.DefaultQuery("timeRange", "30d") // 7d, 30d, 90d
	scope := c.Query("scope")
	platform := c.Query("platform")

	metrics, err := h.validationService.GetValidationMetrics(c.Request.Context(), timeRange, scope, platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve validation metrics",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": metrics,
	})
}

// RunQuickValidation runs a lightweight validation check
// POST /audit/validation/quick
func (h *ValidationHandler) RunQuickValidation(c *gin.Context) {
	var req struct {
		Target    string   `json:"target" binding:"required"`
		Platforms []string `json:"platforms"`
		CheckType string   `json:"checkType" binding:"required"` // build, test, security, performance
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	validCheckTypes := map[string]bool{
		"build": true, "test": true, "security": true, "performance": true,
	}

	if !validCheckTypes[req.CheckType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid check type",
			"message": "Check type must be one of: build, test, security, performance",
		})
		return
	}

	result, err := h.validationService.RunQuickValidation(c.Request.Context(), req.Target, req.Platforms, req.CheckType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to run quick validation",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}