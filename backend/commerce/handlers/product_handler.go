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

// ProductHandler handles product-related commerce endpoints
type ProductHandler struct {
	productService services.ProductService
	validator      *validator.Validate
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		validator:      validator.New(),
	}
}

// GetProducts retrieves products with filters and pagination
// @Summary Get products
// @Description Retrieve products with optional filters, pagination, and sorting
// @Tags commerce, products
// @Accept json
// @Produce json
// @Param business_id query string false "Business ID filter"
// @Param category query string false "Product category filter"
// @Param type query string false "Product type filter"
// @Param status query string false "Product status filter"
// @Param search query string false "Search term"
// @Param page query int false "Page number (default: 0)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Param sort_field query string false "Sort field (default: created_at)"
// @Param sort_order query string false "Sort order: asc|desc (default: desc)"
// @Success 200 {object} responses.DataResponse{data=models.ProductResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse filters
	filters := models.ProductFilters{}

	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		if businessID, err := uuid.Parse(businessIDStr); err == nil {
			filters.BusinessID = &businessID
		}
	}

	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}

	if typeStr := c.Query("type"); typeStr != "" {
		productType := models.ProductType(typeStr)
		if productType.IsValid() {
			filters.Type = &productType
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.ProductStatus(statusStr)
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

	result, err := h.productService.GetProducts(c.Request.Context(), filters, pagination, sort)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get products", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, result)
}

// GetProduct retrieves a single product by ID
// @Summary Get product
// @Description Retrieve a product by its ID
// @Tags commerce, products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} responses.DataResponse{data=models.Product}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid product ID format", "VALIDATION_ERROR")
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			responses.SendErrorResponse(c, http.StatusNotFound, "Product not found", "NOT_FOUND")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve product", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, product)
}

// CreateProduct creates a new product
// @Summary Create product
// @Description Create a new product
// @Tags commerce, products
// @Accept json
// @Produce json
// @Param request body models.CreateProductRequest true "Create product request"
// @Success 201 {object} responses.DataResponse{data=models.Product}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "VALIDATION_ERROR")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR")
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "business not found" {
			responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid business ID", "BUSINESS_NOT_FOUND")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create product", "DATABASE_ERROR")
		return
	}

	c.JSON(http.StatusCreated, responses.DataResponse{
		Success: true,
		Data:    product,
	})
}

// UpdateProduct updates an existing product
// @Summary Update product
// @Description Update an existing product
// @Tags commerce, products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body models.UpdateProductRequest true "Update product request"
// @Success 200 {object} responses.DataResponse{data=models.Product}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid product ID format", "VALIDATION_ERROR")
		return
	}

	var req models.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", "VALIDATION_ERROR")
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "product not found" {
			responses.SendErrorResponse(c, http.StatusNotFound, "Product not found", "NOT_FOUND")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update product", "DATABASE_ERROR")
		return
	}

	responses.SendSuccessResponse(c, product)
}

// DeleteProduct deletes a product
// @Summary Delete product
// @Description Delete a product (soft delete)
// @Tags commerce, products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "Invalid product ID format", "VALIDATION_ERROR")
		return
	}

	err = h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			responses.SendErrorResponse(c, http.StatusNotFound, "Product not found", "NOT_FOUND")
			return
		}
		responses.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete product", "DATABASE_ERROR")
		return
	}

	c.Status(http.StatusNoContent)
}