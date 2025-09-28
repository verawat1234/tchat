package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat.dev/commerce/models"
	"tchat.dev/commerce/services"
	"tchat.dev/shared/responses"
)

// BusinessHandler handles business-related commerce endpoints
type BusinessHandler struct {
	businessService services.BusinessService
	validator       *validator.Validate
}

// NewBusinessHandler creates a new business handler
func NewBusinessHandler(businessService services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
		validator:       validator.New(),
	}
}

// GetBusinesses retrieves businesses with filters and pagination
// @Summary Get businesses (shops)
// @Description Retrieve businesses with optional filters, pagination, and sorting
// @Tags commerce, businesses
// @Accept json
// @Produce json
// @Param country query string false "Country filter (TH, SG, MY, etc.)"
// @Param category query string false "Business category filter"
// @Param status query string false "Verification status filter"
// @Param search query string false "Search term"
// @Param page query int false "Page number (default: 0)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Param sort_field query string false "Sort field (default: created_at)"
// @Param sort_order query string false "Sort order: asc|desc (default: desc)"
// @Success 200 {object} responses.DataResponse{data=models.BusinessResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/shops [get]
func (h *BusinessHandler) GetBusinesses(c *gin.Context) {
	// Parse filters
	filters := models.BusinessFilters{}

	if country := c.Query("country"); country != "" {
		filters.Country = &country
	}

	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.BusinessVerificationStatus(statusStr)
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
			if pageSize > 100 {
				pageSize = 100
			}
			pagination.PageSize = pageSize
		}
	}

	// Parse sorting
	sort := models.SortOptions{
		Field: c.DefaultQuery("sort_field", "created_at"),
		Order: c.DefaultQuery("sort_order", "desc"),
	}

	result, err := h.businessService.GetBusinesses(c.Request.Context(), filters, pagination, sort)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get businesses", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, result)
}

// GetBusiness retrieves a single business by ID
// @Summary Get business details
// @Description Retrieve a business by its ID
// @Tags commerce, businesses
// @Accept json
// @Produce json
// @Param id path string true "Business ID"
// @Success 200 {object} responses.DataResponse{data=models.Business}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/shops/{id} [get]
func (h *BusinessHandler) GetBusiness(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid business ID format", "VALIDATION_ERROR")
		return
	}

	business, err := h.businessService.GetBusiness(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "business not found" {
			responses.SendErrorResponse(c, http.StatusNotFound, "Business not found", "NOT_FOUND")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve business", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, business)
}

// CreateBusiness creates a new business
// @Summary Create business (shop)
// @Description Create a new business
// @Tags commerce, businesses
// @Accept json
// @Produce json
// @Param request body models.CreateBusinessRequest true "Create business request"
// @Success 201 {object} responses.DataResponse{data=models.Business}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/shops [post]
func (h *BusinessHandler) CreateBusiness(c *gin.Context) {
	var req models.CreateBusinessRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "VALIDATION_ERROR")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		// For now, generate a user ID or use anonymous
		userID = uuid.New()
	}

	ownerID := userID.(uuid.UUID)

	business, err := h.businessService.CreateBusiness(c.Request.Context(), ownerID, &req)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create business", "DATABASE_ERROR")
		return
	}

	c.JSON(http.StatusCreated, responses.DataResponse{
		Success: true,
		Data:    business,
	})
}

// UpdateBusiness updates an existing business
// @Summary Update business
// @Description Update an existing business
// @Tags commerce, businesses
// @Accept json
// @Produce json
// @Param id path string true "Business ID"
// @Param request body models.UpdateBusinessRequest true "Update business request"
// @Success 200 {object} responses.DataResponse{data=models.Business}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/shops/{id} [put]
func (h *BusinessHandler) UpdateBusiness(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid business ID format", "VALIDATION_ERROR")
		return
	}

	var req models.UpdateBusinessRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "VALIDATION_ERROR")
		return
	}

	business, err := h.businessService.UpdateBusiness(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "business not found" {
			responses.SendErrorResponse(c, http.StatusNotFound, "Business not found", "NOT_FOUND")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update business", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, business)
}

// DeleteBusiness deletes a business
// @Summary Delete business
// @Description Delete a business (soft delete)
// @Tags commerce, businesses
// @Accept json
// @Produce json
// @Param id path string true "Business ID"
// @Success 204 "No Content"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/shops/{id} [delete]
func (h *BusinessHandler) DeleteBusiness(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid business ID format", "VALIDATION_ERROR")
		return
	}

	err = h.businessService.DeleteBusiness(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "business not found" {
			responses.SendErrorResponse(c, http.StatusNotFound, "Business not found", "NOT_FOUND")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete business", "DATABASE_ERROR")
		return
	}

	c.Status(http.StatusNoContent)
}

// GetBusinessProducts retrieves products for a specific business
// @Summary Get business products
// @Description Retrieve products belonging to a specific business
// @Tags commerce, businesses, products
// @Accept json
// @Produce json
// @Param id path string true "Business ID"
// @Param page query int false "Page number (default: 0)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Param sort_field query string false "Sort field (default: created_at)"
// @Param sort_order query string false "Sort order: asc|desc (default: desc)"
// @Success 200 {object} responses.DataResponse{data=models.ProductResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/shops/{id}/products [get]
func (h *BusinessHandler) GetBusinessProducts(c *gin.Context) {
	idStr := c.Param("id")
	businessID, err := uuid.Parse(idStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid business ID format", "VALIDATION_ERROR")
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
			if pageSize > 100 {
				pageSize = 100
			}
			pagination.PageSize = pageSize
		}
	}

	// Parse sorting
	sort := models.SortOptions{
		Field: c.DefaultQuery("sort_field", "created_at"),
		Order: c.DefaultQuery("sort_order", "desc"),
	}

	result, err := h.businessService.GetBusinessProducts(c.Request.Context(), businessID, pagination, sort)
	if err != nil {
		if err.Error() == "business not found" {
			responses.SendErrorResponse(c, http.StatusNotFound, "Business not found", "NOT_FOUND")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get business products", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, result)
}