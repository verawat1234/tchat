package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tchat/backend/audit/services"
)

type ServiceCompletionHandler struct {
	serviceCompletionService *services.ServiceCompletionService
}

func NewServiceCompletionHandler(serviceCompletionService *services.ServiceCompletionService) *ServiceCompletionHandler {
	return &ServiceCompletionHandler{
		serviceCompletionService: serviceCompletionService,
	}
}

// GetServiceCompletion retrieves completion status for a specific service
// GET /audit/services/:serviceId/completion
func (h *ServiceCompletionHandler) GetServiceCompletion(c *gin.Context) {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing service ID",
			"message": "Service ID is required in URL path",
		})
		return
	}

	platform := c.Query("platform")
	includeDetails := c.DefaultQuery("includeDetails", "true") == "true"

	completion, err := h.serviceCompletionService.GetServiceCompletion(c.Request.Context(), serviceID, platform, includeDetails)
	if err != nil {
		if err.Error() == "service completion not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Service completion not found",
				"message": "No completion data found for service: " + serviceID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve service completion",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": completion,
	})
}

// GetAllServiceCompletions retrieves completion status for all services
// GET /audit/services/completion
func (h *ServiceCompletionHandler) GetAllServiceCompletions(c *gin.Context) {
	platform := c.Query("platform")
	status := c.Query("status") // HEALTHY, DEGRADED, UNHEALTHY, UNKNOWN
	includeDetails := c.DefaultQuery("includeDetails", "false") == "true"

	// Pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 200 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	filters := &services.ServiceCompletionFilters{
		Platform:       platform,
		Status:         status,
		IncludeDetails: includeDetails,
		Limit:          limit,
		Offset:         offset,
	}

	completions, total, err := h.serviceCompletionService.GetAllServiceCompletions(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve service completions",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": completions,
		"meta": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"count":  len(completions),
		},
		"filters": gin.H{
			"platform":       platform,
			"status":         status,
			"includeDetails": includeDetails,
		},
	})
}

// UpdateServiceCompletion updates service completion metrics
// PUT /audit/services/:serviceId/completion
func (h *ServiceCompletionHandler) UpdateServiceCompletion(c *gin.Context) {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing service ID",
			"message": "Service ID is required in URL path",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate allowed update fields
	allowedFields := map[string]bool{
		"placeholderCount":         true,
		"completedCount":           true,
		"testsPassing":             true,
		"buildSuccessful":          true,
		"criticalPlaceholders":     true,
		"highPlaceholders":         true,
		"mediumPlaceholders":       true,
		"lowPlaceholders":          true,
		"estimatedHoursRemaining":  true,
		"estimatedCompletionDate":  true,
		"codeQualityScore":         true,
		"testCoverage":             true,
		"securityScore":            true,
		"performanceScore":         true,
		"documentationScore":       true,
		"healthStatus":             true,
		"serviceVersion":           true,
		"teamOwner":                true,
		"technicalLead":            true,
		"productOwner":             true,
		"dependencies":             true,
		"dependents":               true,
		"blockingServices":         true,
	}

	for field := range updates {
		if !allowedFields[field] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid field",
				"message": "Field '" + field + "' is not allowed for updates",
			})
			return
		}
	}

	completion, err := h.serviceCompletionService.UpdateServiceCompletion(c.Request.Context(), serviceID, updates)
	if err != nil {
		if err.Error() == "service completion not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Service completion not found",
				"message": "No completion data found for service: " + serviceID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update service completion",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    completion,
		"message": "Service completion updated successfully",
	})
}

// GetServiceHealth returns health overview for all services
// GET /audit/services/health
func (h *ServiceCompletionHandler) GetServiceHealth(c *gin.Context) {
	platform := c.Query("platform")
	region := c.Query("region")

	healthOverview, err := h.serviceCompletionService.GetServiceHealth(c.Request.Context(), platform, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve service health",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": healthOverview,
	})
}

// GetServiceMetrics returns aggregated metrics for services
// GET /audit/services/metrics
func (h *ServiceCompletionHandler) GetServiceMetrics(c *gin.Context) {
	platform := c.Query("platform")
	timeRange := c.DefaultQuery("timeRange", "24h") // 24h, 7d, 30d

	metrics, err := h.serviceCompletionService.GetServiceMetrics(c.Request.Context(), platform, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve service metrics",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": metrics,
	})
}

// TriggerServiceRefresh manually triggers a service completion refresh
// POST /audit/services/:serviceId/refresh
func (h *ServiceCompletionHandler) TriggerServiceRefresh(c *gin.Context) {
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing service ID",
			"message": "Service ID is required in URL path",
		})
		return
	}

	var req struct {
		Platform      string `json:"platform"`
		FullRefresh   bool   `json:"fullRefresh"`
		UpdateMetrics bool   `json:"updateMetrics"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	refreshResult, err := h.serviceCompletionService.TriggerServiceRefresh(c.Request.Context(), serviceID, req.Platform, req.FullRefresh, req.UpdateMetrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to refresh service",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    refreshResult,
		"message": "Service refresh completed successfully",
	})
}

// GetRegionalOptimization returns regional optimization status for Southeast Asian markets
// GET /audit/services/regional
func (h *ServiceCompletionHandler) GetRegionalOptimization(c *gin.Context) {
	region := c.Query("region") // TH, SG, MY, ID, PH, VN
	serviceType := c.Query("serviceType")

	validRegions := map[string]bool{
		"TH": true, "SG": true, "MY": true,
		"ID": true, "PH": true, "VN": true,
	}

	if region != "" && !validRegions[region] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid region",
			"message": "Region must be one of: TH, SG, MY, ID, PH, VN",
		})
		return
	}

	regionalData, err := h.serviceCompletionService.GetRegionalOptimization(c.Request.Context(), region, serviceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve regional optimization data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": regionalData,
	})
}

// GetDependencyGraph returns service dependency relationships
// GET /audit/services/dependencies
func (h *ServiceCompletionHandler) GetDependencyGraph(c *gin.Context) {
	serviceID := c.Query("serviceId")
	depth := c.DefaultQuery("depth", "2") // How many levels deep to traverse

	depthInt, err := strconv.Atoi(depth)
	if err != nil || depthInt < 1 || depthInt > 5 {
		depthInt = 2
	}

	dependencyGraph, err := h.serviceCompletionService.GetDependencyGraph(c.Request.Context(), serviceID, depthInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve dependency graph",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dependencyGraph,
	})
}

// GetCompletionTrends returns completion trends over time
// GET /audit/services/trends
func (h *ServiceCompletionHandler) GetCompletionTrends(c *gin.Context) {
	serviceID := c.Query("serviceId")
	timeRange := c.DefaultQuery("timeRange", "30d") // 7d, 30d, 90d
	granularity := c.DefaultQuery("granularity", "daily") // hourly, daily, weekly

	trends, err := h.serviceCompletionService.GetCompletionTrends(c.Request.Context(), serviceID, timeRange, granularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve completion trends",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trends,
	})
}