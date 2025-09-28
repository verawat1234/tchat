package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/services"
	"tchat.dev/shared/middleware"
)

type ReviewHandler struct {
	reviewService services.ReviewService
}

func NewReviewHandler(reviewService services.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// CreateReview creates a new review
// @Summary Create a new review
// @Description Create a review for a product, business, or order
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.CreateReviewRequest true "Review data"
// @Success 201 {object} models.Review
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewService.CreateReview(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetReview retrieves a review by ID
// @Summary Get a review by ID
// @Description Get review details by ID
// @Tags reviews
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} models.Review
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.reviewService.GetReview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// UpdateReview updates an existing review
// @Summary Update a review
// @Description Update review details
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Param review body models.UpdateReviewRequest true "Updated review data"
// @Success 200 {object} models.Review
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req models.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewService.UpdateReview(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// DeleteReview deletes a review
// @Summary Delete a review
// @Description Delete a review by ID
// @Tags reviews
// @Param id path string true "Review ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.reviewService.DeleteReview(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListReviews lists reviews with filters and pagination
// @Summary List reviews
// @Description Get a paginated list of reviews with optional filters
// @Tags reviews
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param type query string false "Review type (product, business, order)"
// @Param status query string false "Review status"
// @Param productId query string false "Product ID"
// @Param businessId query string false "Business ID"
// @Param userId query string false "User ID"
// @Param search query string false "Search term"
// @Success 200 {object} models.ReviewResponse
// @Failure 500 {object} map[string]interface{}
// @Router /reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
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

	filters := services.ReviewFilters{}

	if reviewType := c.Query("type"); reviewType != "" {
		rt := models.ReviewType(reviewType)
		filters.Type = &rt
	}

	if status := c.Query("status"); status != "" {
		rs := models.ReviewStatus(status)
		filters.Status = &rs
	}

	if productID := c.Query("productId"); productID != "" {
		if id, err := uuid.Parse(productID); err == nil {
			filters.ProductID = &id
		}
	}

	if businessID := c.Query("businessId"); businessID != "" {
		if id, err := uuid.Parse(businessID); err == nil {
			filters.BusinessID = &id
		}
	}

	if userID := c.Query("userId"); userID != "" {
		if id, err := uuid.Parse(userID); err == nil {
			filters.UserID = &id
		}
	}

	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}

	response, err := h.reviewService.ListReviews(c.Request.Context(), filters, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProductReviews gets reviews for a specific product
// @Summary Get product reviews
// @Description Get reviews for a specific product
// @Tags reviews
// @Produce json
// @Param productId path string true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} models.ReviewResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{productId}/reviews [get]
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
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

	response, err := h.reviewService.GetReviewsByProduct(c.Request.Context(), productID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetBusinessReviews gets reviews for a specific business
// @Summary Get business reviews
// @Description Get reviews for a specific business
// @Tags reviews
// @Produce json
// @Param businessId path string true "Business ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} models.ReviewResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /businesses/{businessId}/reviews [get]
func (h *ReviewHandler) GetBusinessReviews(c *gin.Context) {
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

	response, err := h.reviewService.GetReviewsByBusiness(c.Request.Context(), businessID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// MarkReviewHelpful marks a review as helpful or not helpful
// @Summary Mark review as helpful
// @Description Mark a review as helpful or not helpful
// @Tags reviews
// @Accept json
// @Param reviewId path string true "Review ID"
// @Param helpful body map[string]bool true "Helpful flag" example({"helpful": true})
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reviews/{reviewId}/helpful [post]
func (h *ReviewHandler) MarkReviewHelpful(c *gin.Context) {
	reviewIDStr := c.Param("reviewId")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Helpful bool `json:"helpful"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reviewService.MarkReviewHelpful(c.Request.Context(), reviewID, userID, req.Helpful); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ReportReview reports a review for inappropriate content
// @Summary Report a review
// @Description Report a review for inappropriate content
// @Tags reviews
// @Accept json
// @Param reviewId path string true "Review ID"
// @Param report body map[string]string true "Report details"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reviews/{reviewId}/report [post]
func (h *ReviewHandler) ReportReview(c *gin.Context) {
	reviewIDStr := c.Param("reviewId")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Reason  string `json:"reason" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reviewService.ReportReview(c.Request.Context(), reviewID, userID, req.Reason, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ModerateReview moderates a review (admin only)
// @Summary Moderate a review
// @Description Moderate a review by changing its status
// @Tags reviews
// @Accept json
// @Param reviewId path string true "Review ID"
// @Param moderation body map[string]interface{} true "Moderation details"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reviews/{reviewId}/moderate [post]
func (h *ReviewHandler) ModerateReview(c *gin.Context) {
	reviewIDStr := c.Param("reviewId")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// Get moderator ID from authentication context and verify admin role
	moderatorID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Add admin role verification
	// claims, _ := middleware.GetUserClaims(c)
	// if !claims.HasRole("admin") {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
	//     return
	// }

	var req struct {
		Status string `json:"status" binding:"required"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := models.ReviewStatus(req.Status)
	if err := h.reviewService.ModerateReview(c.Request.Context(), reviewID, status, moderatorID, req.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAverageRating gets the average rating for a target (product, business, order)
// @Summary Get average rating
// @Description Get average rating and review count for a target
// @Tags reviews
// @Produce json
// @Param type query string true "Target type (product, business, order)"
// @Param targetId query string true "Target ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reviews/average-rating [get]
func (h *ReviewHandler) GetAverageRating(c *gin.Context) {
	targetType := c.Query("type")
	targetIDStr := c.Query("targetId")

	if targetType == "" || targetIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type and targetId are required"})
		return
	}

	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	reviewType := models.ReviewType(targetType)
	avgRating, count, err := h.reviewService.GetAverageRating(c.Request.Context(), reviewType, targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": avgRating,
		"review_count":   count,
		"target_type":    targetType,
		"target_id":      targetID,
	})
}