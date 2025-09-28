package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tchat/backend/audit/models"
	"github.com/tchat/backend/audit/services"
)

type PlaceholdersHandler struct {
	placeholderService *services.PlaceholderService
}

func NewPlaceholdersHandler(placeholderService *services.PlaceholderService) *PlaceholdersHandler {
	return &PlaceholdersHandler{
		placeholderService: placeholderService,
	}
}

// GetPlaceholders retrieves placeholder items with filtering and pagination
// GET /audit/placeholders
func (h *PlaceholdersHandler) GetPlaceholders(c *gin.Context) {
	// Parse query parameters
	serviceID := c.Query("serviceId")
	platform := c.Query("platform")
	status := c.Query("status")
	priority := c.Query("priority")
	assignedTo := c.Query("assignedTo")
	regional := c.Query("regional")

	// Pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid limit parameter",
			"message": "Limit must be between 1 and 1000",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid offset parameter",
			"message": "Offset must be non-negative",
		})
		return
	}

	// Build filter criteria
	filters := &services.PlaceholderFilters{
		ServiceID:  serviceID,
		Platform:   platform,
		Status:     status,
		Priority:   priority,
		AssignedTo: assignedTo,
		Regional:   regional,
		Limit:      limit,
		Offset:     offset,
	}

	// Sorting parameters
	sortBy := c.DefaultQuery("sortBy", "priority")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	filters.SortBy = sortBy
	filters.SortOrder = sortOrder

	// Call service to get placeholders
	placeholders, total, err := h.placeholderService.GetPlaceholders(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve placeholders",
			"message": err.Error(),
		})
		return
	}

	// Prepare response with metadata
	response := gin.H{
		"data": placeholders,
		"meta": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"count":  len(placeholders),
		},
		"filters": gin.H{
			"serviceId":  serviceID,
			"platform":   platform,
			"status":     status,
			"priority":   priority,
			"assignedTo": assignedTo,
			"regional":   regional,
			"sortBy":     sortBy,
			"sortOrder":  sortOrder,
		},
	}

	c.JSON(http.StatusOK, response)
}

// CreatePlaceholder creates a new placeholder item
// POST /audit/placeholders
func (h *PlaceholdersHandler) CreatePlaceholder(c *gin.Context) {
	var req models.PlaceholderItem

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Validate required fields
	if req.ServiceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required field",
			"message": "ServiceID is required",
		})
		return
	}

	if req.Platform == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required field",
			"message": "Platform is required",
		})
		return
	}

	if req.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing required field",
			"message": "FilePath is required",
		})
		return
	}

	if req.Type == "" {
		req.Type = string(models.PlaceholderTypeTODO)
	}

	if req.Priority == "" {
		req.Priority = string(models.PriorityMedium)
	}

	if req.Status == "" {
		req.Status = string(models.StatusOpen)
	}

	if req.DetectedBy == "" {
		req.DetectedBy = "manual"
	}

	if req.PerformanceImpact == "" {
		req.PerformanceImpact = string(models.ImpactUnknown)
	}

	// Create placeholder through service
	placeholder, err := h.placeholderService.CreatePlaceholder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create placeholder",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    placeholder,
		"message": "Placeholder created successfully",
	})
}

// UpdatePlaceholder updates an existing placeholder item
// PATCH /audit/placeholders/:id
func (h *PlaceholdersHandler) UpdatePlaceholder(c *gin.Context) {
	placeholderID := c.Param("id")
	if placeholderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing placeholder ID",
			"message": "Placeholder ID is required in URL path",
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
		"status":                true,
		"priority":              true,
		"assignedTo":            true,
		"estimatedHours":        true,
		"description":           true,
		"performanceImpact":     true,
		"resolutionNotes":       true,
		"reviewedBy":            true,
		"technicalDebtScore":    true,
		"complexityScore":       true,
		"riskScore":            true,
		"regionalContext":       true,
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

	// Update placeholder through service
	placeholder, err := h.placeholderService.UpdatePlaceholder(c.Request.Context(), placeholderID, updates)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Placeholder not found",
				"message": "No placeholder found with ID: " + placeholderID,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update placeholder",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    placeholder,
		"message": "Placeholder updated successfully",
	})
}

// GetPlaceholderStats returns statistics about placeholder items
// GET /audit/placeholders/stats
func (h *PlaceholdersHandler) GetPlaceholderStats(c *gin.Context) {
	serviceID := c.Query("serviceId")
	platform := c.Query("platform")

	stats, err := h.placeholderService.GetPlaceholderStats(c.Request.Context(), serviceID, platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve placeholder statistics",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// BulkUpdatePlaceholders updates multiple placeholders at once
// POST /audit/placeholders/bulk-update
func (h *PlaceholdersHandler) BulkUpdatePlaceholders(c *gin.Context) {
	var req struct {
		PlaceholderIDs []string               `json:"placeholderIds" binding:"required"`
		Updates        map[string]interface{} `json:"updates" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	if len(req.PlaceholderIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No placeholders specified",
			"message": "PlaceholderIDs array cannot be empty",
		})
		return
	}

	if len(req.PlaceholderIDs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Too many placeholders",
			"message": "Cannot update more than 100 placeholders at once",
		})
		return
	}

	// Validate allowed update fields
	allowedFields := map[string]bool{
		"status":            true,
		"priority":          true,
		"assignedTo":        true,
		"performanceImpact": true,
		"regionalContext":   true,
	}

	for field := range req.Updates {
		if !allowedFields[field] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid field",
				"message": "Field '" + field + "' is not allowed for bulk updates",
			})
			return
		}
	}

	// Perform bulk update through service
	updatedCount, err := h.placeholderService.BulkUpdatePlaceholders(c.Request.Context(), req.PlaceholderIDs, req.Updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to bulk update placeholders",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bulk update completed successfully",
		"updated": updatedCount,
		"total":   len(req.PlaceholderIDs),
	})
}

// GetPlaceholdersByRegion returns placeholders filtered by Southeast Asian regions
// GET /audit/placeholders/region/:region
func (h *PlaceholdersHandler) GetPlaceholdersByRegion(c *gin.Context) {
	region := c.Param("region")
	validRegions := map[string]bool{
		"TH": true, // Thailand
		"SG": true, // Singapore
		"MY": true, // Malaysia
		"ID": true, // Indonesia
		"PH": true, // Philippines
		"VN": true, // Vietnam
	}

	if !validRegions[region] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid region",
			"message": "Region must be one of: TH, SG, MY, ID, PH, VN",
		})
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	placeholders, total, err := h.placeholderService.GetPlaceholdersByRegion(c.Request.Context(), region, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve regional placeholders",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": placeholders,
		"meta": gin.H{
			"region": region,
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"count":  len(placeholders),
		},
	})
}

// ArchivePlaceholder archives a completed placeholder
// POST /audit/placeholders/:id/archive
func (h *PlaceholdersHandler) ArchivePlaceholder(c *gin.Context) {
	placeholderID := c.Param("id")
	if placeholderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing placeholder ID",
			"message": "Placeholder ID is required in URL path",
		})
		return
	}

	var req struct {
		ArchiveReason string `json:"archiveReason"`
		ReviewedBy    string `json:"reviewedBy"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	err := h.placeholderService.ArchivePlaceholder(c.Request.Context(), placeholderID, req.ArchiveReason, req.ReviewedBy)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Placeholder not found",
				"message": "No placeholder found with ID: " + placeholderID,
			})
			return
		}

		if strings.Contains(err.Error(), "not completed") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Cannot archive",
				"message": "Only completed placeholders can be archived",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to archive placeholder",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Placeholder archived successfully",
	})
}