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

type WishlistHandler struct {
	wishlistService services.WishlistService
}

func NewWishlistHandler(wishlistService services.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// CreateWishlist creates a new wishlist
// @Summary Create a new wishlist
// @Description Create a new wishlist for the user
// @Tags wishlists
// @Accept json
// @Produce json
// @Param wishlist body models.CreateWishlistRequest true "Wishlist data"
// @Success 201 {object} models.Wishlist
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /wishlists [post]
func (h *WishlistHandler) CreateWishlist(c *gin.Context) {
	var req models.CreateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	wishlist, err := h.wishlistService.CreateWishlist(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wishlist)
}

// GetWishlist retrieves a wishlist by ID
// @Summary Get a wishlist by ID
// @Description Get wishlist details by ID
// @Tags wishlists
// @Produce json
// @Param id path string true "Wishlist ID"
// @Success 200 {object} models.Wishlist
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /wishlists/{id} [get]
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	wishlist, err := h.wishlistService.GetWishlist(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wishlist)
}

// GetWishlistByShareToken retrieves a wishlist by share token
// @Summary Get a wishlist by share token
// @Description Get wishlist details by share token
// @Tags wishlists
// @Produce json
// @Param token path string true "Share token"
// @Success 200 {object} models.Wishlist
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /wishlists/shared/{token} [get]
func (h *WishlistHandler) GetWishlistByShareToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Share token is required"})
		return
	}

	wishlist, err := h.wishlistService.GetWishlistByShareToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wishlist)
}

// UpdateWishlist updates an existing wishlist
// @Summary Update a wishlist
// @Description Update wishlist details
// @Tags wishlists
// @Accept json
// @Produce json
// @Param id path string true "Wishlist ID"
// @Param wishlist body models.UpdateWishlistRequest true "Updated wishlist data"
// @Success 200 {object} models.Wishlist
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /wishlists/{id} [put]
func (h *WishlistHandler) UpdateWishlist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	var req models.UpdateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wishlist, err := h.wishlistService.UpdateWishlist(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wishlist)
}

// DeleteWishlist deletes a wishlist
// @Summary Delete a wishlist
// @Description Delete a wishlist by ID
// @Tags wishlists
// @Param id path string true "Wishlist ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /wishlists/{id} [delete]
func (h *WishlistHandler) DeleteWishlist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	if err := h.wishlistService.DeleteWishlist(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUserWishlists lists user's wishlists with pagination
// @Summary List user wishlists
// @Description Get a paginated list of user's wishlists
// @Tags wishlists
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} models.WishlistResponse
// @Failure 500 {object} map[string]interface{}
// @Router /wishlists [get]
func (h *WishlistHandler) ListUserWishlists(c *gin.Context) {
	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	response, err := h.wishlistService.ListUserWishlists(c.Request.Context(), userID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetDefaultWishlist gets user's default wishlist
// @Summary Get default wishlist
// @Description Get user's default wishlist
// @Tags wishlists
// @Produce json
// @Success 200 {object} models.Wishlist
// @Failure 500 {object} map[string]interface{}
// @Router /wishlists/default [get]
func (h *WishlistHandler) GetDefaultWishlist(c *gin.Context) {
	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	wishlist, err := h.wishlistService.GetDefaultWishlist(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wishlist)
}

// AddToWishlist adds a product to a wishlist
// @Summary Add product to wishlist
// @Description Add a product to a specific wishlist
// @Tags wishlists
// @Accept json
// @Param id path string true "Wishlist ID"
// @Param item body models.AddToWishlistRequest true "Product to add"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /wishlists/{id}/items [post]
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	idStr := c.Param("id")
	wishlistID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	var req models.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.wishlistService.AddToWishlist(c.Request.Context(), wishlistID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveFromWishlist removes a product from a wishlist
// @Summary Remove product from wishlist
// @Description Remove a product from a specific wishlist
// @Tags wishlists
// @Param id path string true "Wishlist ID"
// @Param productId path string true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /wishlists/{id}/items/{productId} [delete]
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	idStr := c.Param("id")
	wishlistID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	productIDStr := c.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.wishlistService.RemoveFromWishlist(c.Request.Context(), wishlistID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ShareWishlist shares a wishlist with another user
// @Summary Share wishlist
// @Description Share a wishlist with another user
// @Tags wishlists
// @Accept json
// @Param id path string true "Wishlist ID"
// @Param share body map[string]interface{} true "Share details"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /wishlists/{id}/share [post]
func (h *WishlistHandler) ShareWishlist(c *gin.Context) {
	idStr := c.Param("id")
	wishlistID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	var req struct {
		SharedWith uuid.UUID `json:"sharedWith" binding:"required"`
		Permission string    `json:"permission"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Permission == "" {
		req.Permission = "view"
	}

	if err := h.wishlistService.ShareWishlist(c.Request.Context(), wishlistID, req.SharedWith, req.Permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetSharedWishlists gets wishlists shared with the user
// @Summary Get shared wishlists
// @Description Get wishlists that have been shared with the user
// @Tags wishlists
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} models.WishlistResponse
// @Failure 500 {object} map[string]interface{}
// @Router /wishlists/shared [get]
func (h *WishlistHandler) GetSharedWishlists(c *gin.Context) {
	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	response, err := h.wishlistService.GetSharedWishlists(c.Request.Context(), userID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// FollowProduct starts following a product for notifications
// @Summary Follow product
// @Description Start following a product for price changes and stock updates
// @Tags product-follows
// @Accept json
// @Param productId path string true "Product ID"
// @Param preferences body services.ProductFollowPreferences true "Follow preferences"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{productId}/follow [post]
func (h *WishlistHandler) FollowProduct(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req services.ProductFollowPreferences
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get business ID from product
	businessID := uuid.New() // Placeholder

	if err := h.wishlistService.FollowProduct(c.Request.Context(), userID, productID, businessID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UnfollowProduct stops following a product
// @Summary Unfollow product
// @Description Stop following a product
// @Tags product-follows
// @Param productId path string true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{productId}/follow [delete]
func (h *WishlistHandler) UnfollowProduct(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := h.wishlistService.UnfollowProduct(c.Request.Context(), userID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListFollowedProducts lists products the user is following
// @Summary List followed products
// @Description Get a list of products the user is following
// @Tags product-follows
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/following [get]
func (h *WishlistHandler) ListFollowedProducts(c *gin.Context) {
	// Get user ID from authentication context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	follows, total, err := h.wishlistService.ListFollowedProducts(c.Request.Context(), userID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	c.JSON(http.StatusOK, gin.H{
		"follows":     follows,
		"total":       total,
		"page":        pagination.Page,
		"pageSize":    pagination.PageSize,
		"totalPages":  totalPages,
	})
}