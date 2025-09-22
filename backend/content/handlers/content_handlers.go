package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat.dev/content/models"
	"tchat.dev/content/services"
	"tchat.dev/content/utils"
)

// ContentHandlers provides HTTP handlers for content management
type ContentHandlers struct {
	contentService *services.ContentService
	validator      *validator.Validate
}

// NewContentHandlers creates a new content handlers instance
func NewContentHandlers(contentService *services.ContentService) *ContentHandlers {
	return &ContentHandlers{
		contentService: contentService,
		validator:      validator.New(),
	}
}

// GetContentItems retrieves content items with filters and pagination
// @Summary Get content items
// @Description Retrieve content items with optional filters, pagination, and sorting
// @Tags content
// @Accept json
// @Produce json
// @Param category query string false "Category filter"
// @Param type query string false "Content type filter"
// @Param status query string false "Status filter"
// @Param search query string false "Search term"
// @Param page query int false "Page number (default: 0)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Param sort_field query string false "Sort field (default: created_at)"
// @Param sort_order query string false "Sort order: asc|desc (default: desc)"
// @Success 200 {object} utils.DataResponse{data=models.ContentResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content [get]
func (h *ContentHandlers) GetContentItems(c *gin.Context) {
	// Parse filters
	filters := models.ContentFilters{}

	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}

	if typeStr := c.Query("type"); typeStr != "" {
		contentType := models.ContentType(typeStr)
		if contentType.IsValid() {
			filters.Type = &contentType
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.ContentStatus(statusStr)
		if status.IsValid() {
			filters.Status = &status
		}
	}

	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	// Parse pagination
	pagination := models.Pagination{
		Page:     0,
		PageSize: 20,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			pagination.PageSize = pageSize
		}
	}

	// Parse sorting
	sort := models.SortOptions{
		Field: c.DefaultQuery("sort_field", "created_at"),
		Order: c.DefaultQuery("sort_order", "desc"),
	}

	result, err := h.contentService.GetContentItems(c.Request.Context(), filters, pagination, sort)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get content items", err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// GetContent retrieves a single content item by ID
// @Summary Get content item
// @Description Retrieve a content item by its ID
// @Tags content
// @Accept json
// @Produce json
// @Param id path string true "Content ID"
// @Success 200 {object} utils.DataResponse{data=models.ContentItem}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/{id} [get]
func (h *ContentHandlers) GetContent(c *gin.Context) {
	idStr := c.Param("id")

	// Try parsing as UUID first for backwards compatibility
	if id, err := uuid.Parse(idStr); err == nil {
		content, err := h.contentService.GetContent(c.Request.Context(), id)
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Content not found", err.Error())
			return
		}
		utils.SuccessResponse(c, content)
		return
	}

	// If not a UUID, treat as content key and look it up
	content, err := h.contentService.GetContentByKey(c.Request.Context(), idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Content not found", err.Error())
		return
	}

	utils.SuccessResponse(c, content)
}

// GetContentByCategory retrieves content items by category
// @Summary Get content by category
// @Description Retrieve content items filtered by category
// @Tags content
// @Accept json
// @Produce json
// @Param category path string true "Category name"
// @Param page query int false "Page number (default: 0)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Param sort_field query string false "Sort field (default: created_at)"
// @Param sort_order query string false "Sort order: asc|desc (default: desc)"
// @Success 200 {object} utils.DataResponse{data=models.ContentResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/category/{category} [get]
func (h *ContentHandlers) GetContentByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Category is required", "")
		return
	}

	// Parse pagination
	pagination := models.Pagination{
		Page:     0,
		PageSize: 20,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			pagination.PageSize = pageSize
		}
	}

	// Parse sorting
	sort := models.SortOptions{
		Field: c.DefaultQuery("sort_field", "created_at"),
		Order: c.DefaultQuery("sort_order", "desc"),
	}

	result, err := h.contentService.GetContentByCategory(c.Request.Context(), category, pagination, sort)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get content by category", err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// CreateContent creates a new content item
// @Summary Create content item
// @Description Create a new content item
// @Tags content
// @Accept json
// @Produce json
// @Param request body models.CreateContentRequest true "Create content request"
// @Success 201 {object} utils.DataResponse{data=models.ContentItem}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content [post]
func (h *ContentHandlers) CreateContent(c *gin.Context) {
	var req models.CreateContentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate content value based on type
	if err := h.contentService.ValidateContentValue(req.Type, req.Value); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content value", err.Error())
		return
	}

	content, err := h.contentService.CreateContent(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create content", err.Error())
		return
	}

	utils.CreatedResponse(c, content)
}

// UpdateContent updates an existing content item
// @Summary Update content item
// @Description Update an existing content item
// @Tags content
// @Accept json
// @Produce json
// @Param id path string true "Content ID"
// @Param request body models.UpdateContentRequest true "Update content request"
// @Success 200 {object} utils.DataResponse{data=models.ContentItem}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/{id} [put]
func (h *ContentHandlers) UpdateContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err.Error())
		return
	}

	var req models.UpdateContentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate content value if provided
	if req.Value != nil && req.Type != nil {
		if err := h.contentService.ValidateContentValue(*req.Type, *req.Value); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content value", err.Error())
			return
		}
	}

	content, err := h.contentService.UpdateContent(c.Request.Context(), id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update content", err.Error())
		return
	}

	utils.SuccessResponse(c, content)
}

// DeleteContent deletes a content item
// @Summary Delete content item
// @Description Delete a content item (soft delete)
// @Tags content
// @Accept json
// @Produce json
// @Param id path string true "Content ID"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/{id} [delete]
func (h *ContentHandlers) DeleteContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err.Error())
		return
	}

	err = h.contentService.DeleteContent(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete content", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// PublishContent publishes a content item
// @Summary Publish content item
// @Description Publish a content item
// @Tags content
// @Accept json
// @Produce json
// @Param id path string true "Content ID"
// @Success 200 {object} utils.DataResponse{data=models.ContentItem}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/{id}/publish [post]
func (h *ContentHandlers) PublishContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err.Error())
		return
	}

	content, err := h.contentService.PublishContent(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to publish content", err.Error())
		return
	}

	utils.SuccessResponse(c, content)
}

// ArchiveContent archives a content item
// @Summary Archive content item
// @Description Archive a content item
// @Tags content
// @Accept json
// @Produce json
// @Param id path string true "Content ID"
// @Success 200 {object} utils.DataResponse{data=models.ContentItem}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/{id}/archive [post]
func (h *ContentHandlers) ArchiveContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err.Error())
		return
	}

	content, err := h.contentService.ArchiveContent(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to archive content", err.Error())
		return
	}

	utils.SuccessResponse(c, content)
}

// BulkUpdateContent updates multiple content items
// @Summary Bulk update content items
// @Description Update multiple content items at once
// @Tags content
// @Accept json
// @Produce json
// @Param request body models.BulkUpdateRequest true "Bulk update request"
// @Success 200 {object} utils.DataResponse{data=[]models.ContentItem}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/bulk [put]
func (h *ContentHandlers) BulkUpdateContent(c *gin.Context) {
	var req models.BulkUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if len(req.IDs) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "No content IDs provided", "")
		return
	}

	items, err := h.contentService.BulkUpdateContent(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to bulk update content", err.Error())
		return
	}

	utils.SuccessResponse(c, items)
}

// SyncContent handles content synchronization
// @Summary Sync content
// @Description Synchronize content based on last sync time
// @Tags content
// @Accept json
// @Produce json
// @Param request body models.SyncContentRequest true "Sync request"
// @Success 200 {object} utils.DataResponse{data=models.SyncContentResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/sync [post]
func (h *ContentHandlers) SyncContent(c *gin.Context) {
	var req models.SyncContentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.contentService.SyncContent(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to sync content", err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// GetContentCategories retrieves content categories
// @Summary Get content categories
// @Description Retrieve content categories with pagination and sorting
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 0)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Param sort_field query string false "Sort field (default: sort_order)"
// @Param sort_order query string false "Sort order: asc|desc (default: asc)"
// @Success 200 {object} utils.DataResponse{data=models.ContentCategoryResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/categories [get]
func (h *ContentHandlers) GetContentCategories(c *gin.Context) {
	// Parse pagination
	pagination := models.Pagination{
		Page:     0,
		PageSize: 20,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			pagination.PageSize = pageSize
		}
	}

	// Parse sorting
	sort := models.SortOptions{
		Field: c.DefaultQuery("sort_field", "sort_order"),
		Order: c.DefaultQuery("sort_order", "asc"),
	}

	result, err := h.contentService.GetCategories(c.Request.Context(), pagination, sort)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get categories", err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// GetContentVersions retrieves versions for a content item
// @Summary Get content versions
// @Description Retrieve version history for a content item
// @Tags versions
// @Accept json
// @Produce json
// @Param id path string true "Content ID"
// @Param page query int false "Page number (default: 0)"
// @Param page_size query int false "Page size (default: 10, max: 50)"
// @Success 200 {object} utils.DataResponse{data=models.ContentVersionResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/{id}/versions [get]
func (h *ContentHandlers) GetContentVersions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err.Error())
		return
	}

	// Parse pagination
	pagination := models.Pagination{
		Page:     0,
		PageSize: 10,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 0 {
			pagination.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			pagination.PageSize = pageSize
		}
	}

	result, err := h.contentService.GetContentVersions(c.Request.Context(), id, pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get content versions", err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// RevertContentVersion reverts content to a specific version
// @Summary Revert content version
// @Description Revert content to a specific version
// @Tags versions
// @Accept json
// @Produce json
// @Param id path string true "Content ID"
// @Param version path int true "Version number"
// @Success 200 {object} utils.DataResponse{data=models.ContentItem}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /content/{id}/versions/{version}/revert [post]
func (h *ContentHandlers) RevertContentVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid content ID", err.Error())
		return
	}

	versionStr := c.Param("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil || version <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid version number", err.Error())
		return
	}

	content, err := h.contentService.RevertContentVersion(c.Request.Context(), id, version)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to revert content version", err.Error())
		return
	}

	utils.SuccessResponse(c, content)
}

// RegisterContentRoutes registers all content-related routes
func RegisterContentRoutes(router *gin.RouterGroup, handlers *ContentHandlers) {
	content := router.Group("/content")
	{
		// Content CRUD operations
		content.GET("", handlers.GetContentItems)
		content.POST("", handlers.CreateContent)
		content.GET("/:id", handlers.GetContent)
		content.PUT("/:id", handlers.UpdateContent)
		content.DELETE("/:id", handlers.DeleteContent)

		// Content operations
		content.POST("/:id/publish", handlers.PublishContent)
		content.POST("/:id/archive", handlers.ArchiveContent)
		content.PUT("/bulk", handlers.BulkUpdateContent)
		content.POST("/sync", handlers.SyncContent)

		// Category operations
		content.GET("/categories", handlers.GetContentCategories)
		content.GET("/category/:category", handlers.GetContentByCategory)

		// Version operations
		content.GET("/:id/versions", handlers.GetContentVersions)
		content.POST("/:id/versions/:version/revert", handlers.RevertContentVersion)
	}
}