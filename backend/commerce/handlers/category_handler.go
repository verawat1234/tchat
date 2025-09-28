package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/services"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory creates a new category
// @Summary Create a new category
// @Description Create a new product category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category data"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategory retrieves a category by ID
// @Summary Get a category by ID
// @Description Get category details by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.categoryService.GetCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryByPath retrieves a category by path
// @Summary Get a category by path
// @Description Get category details by URL path
// @Tags categories
// @Produce json
// @Param path path string true "Category path"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/path/{path} [get]
func (h *CategoryHandler) GetCategoryByPath(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category path is required"})
		return
	}

	category, err := h.categoryService.GetCategoryByPath(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory updates an existing category
// @Summary Update a category
// @Description Update category details
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.UpdateCategoryRequest true "Updated category data"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Param id path string true "Category ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListCategories lists categories with filters and pagination
// @Summary List categories
// @Description Get a paginated list of categories with optional filters
// @Tags categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param businessId query string false "Business ID"
// @Param parentId query string false "Parent category ID"
// @Param level query int false "Category level"
// @Param status query string false "Category status"
// @Param search query string false "Search term"
// @Success 200 {object} models.CategoryResponse
// @Failure 500 {object} map[string]interface{}
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	pagination := models.Pagination{
		Page:     1,
		PageSize: 20,
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		pagination.Page = page
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20")); err == nil && pageSize > 0 && pageSize <= 100 {
		pagination.PageSize = pageSize
	}

	filters := make(map[string]interface{})

	if businessID := c.Query("businessId"); businessID != "" {
		if id, err := uuid.Parse(businessID); err == nil {
			filters["business_id"] = id
		}
	}

	if parentID := c.Query("parentId"); parentID != "" {
		if id, err := uuid.Parse(parentID); err == nil {
			filters["parent_id"] = id
		}
	}

	if level := c.Query("level"); level != "" {
		if lvl, err := strconv.Atoi(level); err == nil {
			filters["level"] = lvl
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = models.CategoryStatus(status)
	}

	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	response, err := h.categoryService.ListCategories(c.Request.Context(), filters, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCategoryChildren gets child categories
// @Summary Get category children
// @Description Get child categories of a parent category
// @Tags categories
// @Produce json
// @Param id path string true "Parent Category ID"
// @Success 200 {array} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetCategoryChildren(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Get recursive parameter from query (default false)
	recursive := c.DefaultQuery("recursive", "false") == "true"

	children, err := h.categoryService.GetCategoryChildren(c.Request.Context(), id, recursive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, children)
}

// GetBusinessCategories gets categories for a business
// @Summary Get business categories
// @Description Get categories for a specific business
// @Tags categories
// @Produce json
// @Param businessId path string true "Business ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} models.CategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /businesses/{businessId}/categories [get]
func (h *CategoryHandler) GetBusinessCategories(c *gin.Context) {
	businessIDStr := c.Param("businessId")
	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	pagination := models.Pagination{
		Page:     1,
		PageSize: 20,
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		pagination.Page = page
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20")); err == nil && pageSize > 0 && pageSize <= 100 {
		pagination.PageSize = pageSize
	}

	response, err := h.categoryService.GetBusinessCategories(c.Request.Context(), businessID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetGlobalCategories gets global categories
// @Summary Get global categories
// @Description Get global categories not tied to specific businesses
// @Tags categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} models.CategoryResponse
// @Failure 500 {object} map[string]interface{}
// @Router /categories/global [get]
func (h *CategoryHandler) GetGlobalCategories(c *gin.Context) {
	pagination := models.Pagination{
		Page:     1,
		PageSize: 20,
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		pagination.Page = page
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20")); err == nil && pageSize > 0 && pageSize <= 100 {
		pagination.PageSize = pageSize
	}

	response, err := h.categoryService.GetGlobalCategories(c.Request.Context(), pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFeaturedCategories gets featured categories
// @Summary Get featured categories
// @Description Get featured categories for homepage or promotions
// @Tags categories
// @Produce json
// @Param businessId query string false "Business ID"
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]interface{}
// @Router /categories/featured [get]
func (h *CategoryHandler) GetFeaturedCategories(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	var businessID *uuid.UUID
	if businessIDStr := c.Query("businessId"); businessIDStr != "" {
		if id, err := uuid.Parse(businessIDStr); err == nil {
			businessID = &id
		}
	}

	categories, err := h.categoryService.GetFeaturedCategories(c.Request.Context(), businessID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetRootCategories gets root categories
// @Summary Get root categories
// @Description Get top-level categories without parents
// @Tags categories
// @Produce json
// @Param businessId query string false "Business ID"
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]interface{}
// @Router /categories/root [get]
func (h *CategoryHandler) GetRootCategories(c *gin.Context) {
	var businessID *uuid.UUID
	if businessIDStr := c.Query("businessId"); businessIDStr != "" {
		if id, err := uuid.Parse(businessIDStr); err == nil {
			businessID = &id
		}
	}

	categories, err := h.categoryService.GetRootCategories(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// AddProductToCategory adds a product to a category
// @Summary Add product to category
// @Description Add a product to a specific category
// @Tags categories
// @Accept json
// @Param categoryId path string true "Category ID"
// @Param assignment body models.AddProductToCategoryRequest true "Product assignment"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{categoryId}/products [post]
func (h *CategoryHandler) AddProductToCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req models.AddProductToCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoryService.AddProductToCategory(c.Request.Context(), req.ProductID, categoryID, req.IsPrimary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveProductFromCategory removes a product from a category
// @Summary Remove product from category
// @Description Remove a product from a specific category
// @Tags categories
// @Param categoryId path string true "Category ID"
// @Param productId path string true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{categoryId}/products/{productId} [delete]
func (h *CategoryHandler) RemoveProductFromCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	productIDStr := c.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.categoryService.RemoveProductFromCategory(c.Request.Context(), productID, categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCategoryProducts gets products in a category
// @Summary Get category products
// @Description Get products assigned to a category
// @Tags categories
// @Produce json
// @Param categoryId path string true "Category ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} models.ProductCategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{categoryId}/products [get]
func (h *CategoryHandler) GetCategoryProducts(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	pagination := models.Pagination{
		Page:     1,
		PageSize: 20,
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		pagination.Page = page
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20")); err == nil && pageSize > 0 && pageSize <= 100 {
		pagination.PageSize = pageSize
	}

	productIDs, total, err := h.categoryService.GetCategoryProducts(c.Request.Context(), categoryID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	response := map[string]interface{}{
		"productIds":  productIDs,
		"total":       total,
		"page":        pagination.Page,
		"pageSize":    pagination.PageSize,
		"totalPages":  totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// TrackCategoryView tracks category view for analytics
// @Summary Track category view
// @Description Track a category view for analytics
// @Tags categories
// @Accept json
// @Param categoryId path string true "Category ID"
// @Param view body models.TrackCategoryViewRequest true "View data"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{categoryId}/views [post]
func (h *CategoryHandler) TrackCategoryView(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req models.TrackCategoryViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoryService.TrackCategoryView(c.Request.Context(), categoryID, req.UserID, req.SessionID, req.IPAddress, req.UserAgent, req.Referrer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCategoryAnalytics gets category analytics
// @Summary Get category analytics
// @Description Get analytics data for a category
// @Tags categories
// @Produce json
// @Param categoryId path string true "Category ID"
// @Param dateFrom query string false "Date from (YYYY-MM-DD)"
// @Param dateTo query string false "Date to (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /categories/{categoryId}/analytics [get]
func (h *CategoryHandler) GetCategoryAnalytics(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	analytics, err := h.categoryService.GetCategoryAnalytics(c.Request.Context(), categoryID, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}