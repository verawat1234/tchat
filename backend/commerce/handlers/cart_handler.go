package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/services"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// GetCart retrieves user's current cart
// @Summary Get user cart
// @Description Get current cart for user or session
// @Tags carts
// @Produce json
// @Param userId query string false "User ID"
// @Param sessionId query string false "Session ID"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /carts [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userIDStr := c.Query("userId")
	sessionID := c.Query("sessionId")

	if userIDStr == "" && sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either userId or sessionId is required"})
		return
	}

	var cart *models.Cart
	var err error

	if userIDStr != "" {
		userID, parseErr := uuid.Parse(userIDStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		cart, err = h.cartService.GetCartByUserID(c.Request.Context(), userID)
	} else {
		cart, err = h.cartService.GetCartBySessionID(c.Request.Context(), sessionID)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddToCart adds an item to the cart
// @Summary Add item to cart
// @Description Add a product to user's cart
// @Tags carts
// @Accept json
// @Produce json
// @Param cartId path string true "Cart ID"
// @Param item body models.AddToCartRequest true "Item to add"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/{cartId}/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	cartIDStr := c.Param("cartId")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.AddToCart(c.Request.Context(), cartID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated cart to return
	cart, err := h.cartService.GetCart(c.Request.Context(), nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// UpdateCartItem updates quantity of an item in cart
// @Summary Update cart item
// @Description Update quantity of an item in the cart
// @Tags carts
// @Accept json
// @Produce json
// @Param cartId path string true "Cart ID"
// @Param productId path string true "Product ID"
// @Param update body models.UpdateCartItemRequest true "Update data"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/{cartId}/items/{productId} [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	cartIDStr := c.Param("cartId")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	productIDStr := c.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.UpdateCartItem(c.Request.Context(), cartID, productID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated cart to return
	cart, err := h.cartService.GetCart(c.Request.Context(), nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// RemoveFromCart removes an item from the cart
// @Summary Remove item from cart
// @Description Remove a product from user's cart
// @Tags carts
// @Param cartId path string true "Cart ID"
// @Param productId path string true "Product ID"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/{cartId}/items/{productId} [delete]
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	cartIDStr := c.Param("cartId")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	productIDStr := c.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.cartService.RemoveFromCart(c.Request.Context(), cartID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated cart to return
	cart, err := h.cartService.GetCart(c.Request.Context(), nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// ClearCart removes all items from the cart
// @Summary Clear cart
// @Description Remove all items from user's cart
// @Tags carts
// @Param cartId path string true "Cart ID"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/{cartId}/clear [post]
func (h *CartHandler) ClearCart(c *gin.Context) {
	cartIDStr := c.Param("cartId")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	if err := h.cartService.ClearCart(c.Request.Context(), cartID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated cart to return
	cart, err := h.cartService.GetCart(c.Request.Context(), nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// MergeCart merges guest cart with user cart
// @Summary Merge carts
// @Description Merge guest session cart with user cart after login
// @Tags carts
// @Accept json
// @Param merge body models.MergeCartRequest true "Merge data"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/merge [post]
func (h *CartHandler) MergeCart(c *gin.Context) {
	var req models.MergeCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get guest cart by session ID
	guestCart, err := h.cartService.GetCartBySessionID(c.Request.Context(), req.SessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guest cart not found"})
		return
	}

	// Get user cart by user ID
	userCart, err := h.cartService.GetCartByUserID(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User cart not found"})
		return
	}

	// Merge carts using their IDs
	if err := h.cartService.MergeCart(c.Request.Context(), guestCart.ID, userCart.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated user cart
	updatedCart, err := h.cartService.GetCartByUserID(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}

// GetAbandonedCarts gets abandoned carts for recovery campaigns
// @Summary Get abandoned carts
// @Description Get abandoned carts with filtering for recovery campaigns
// @Tags carts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param minValue query number false "Minimum cart value"
// @Param dateFrom query string false "Date from (YYYY-MM-DD)"
// @Param dateTo query string false "Date to (YYYY-MM-DD)"
// @Success 200 {object} models.CartResponse
// @Failure 500 {object} map[string]interface{}
// @Router /carts/abandoned [get]
func (h *CartHandler) GetAbandonedCarts(c *gin.Context) {
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

	if minValue := c.Query("minValue"); minValue != "" {
		if value, err := strconv.ParseFloat(minValue, 64); err == nil {
			filters["min_value"] = value
		}
	}

	if dateFrom := c.Query("dateFrom"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}

	if dateTo := c.Query("dateTo"); dateTo != "" {
		filters["date_to"] = dateTo
	}

	response, err := h.cartService.GetAbandonedCarts(c.Request.Context(), filters, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateAbandonmentTracking starts tracking cart abandonment
// @Summary Track cart abandonment
// @Description Start tracking cart abandonment for recovery
// @Tags carts
// @Accept json
// @Param tracking body models.CreateAbandonmentTrackingRequest true "Tracking data"
// @Success 201 {object} models.CartAbandonmentTracking
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/abandonment [post]
func (h *CartHandler) CreateAbandonmentTracking(c *gin.Context) {
	var req models.CreateAbandonmentTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tracking, err := h.cartService.CreateAbandonmentTracking(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tracking)
}

// GetAbandonmentAnalytics gets cart abandonment analytics
// @Summary Get abandonment analytics
// @Description Get cart abandonment analytics and recovery data
// @Tags carts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param stage query string false "Abandonment stage"
// @Param dateFrom query string false "Date from (YYYY-MM-DD)"
// @Param dateTo query string false "Date to (YYYY-MM-DD)"
// @Param isRecovered query bool false "Filter by recovery status"
// @Success 200 {object} models.AbandonmentTrackingResponse
// @Failure 500 {object} map[string]interface{}
// @Router /carts/abandonment/analytics [get]
func (h *CartHandler) GetAbandonmentAnalytics(c *gin.Context) {
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

	if stage := c.Query("stage"); stage != "" {
		filters["stage"] = stage
	}

	if dateFrom := c.Query("dateFrom"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}

	if dateTo := c.Query("dateTo"); dateTo != "" {
		filters["date_to"] = dateTo
	}

	if isRecovered := c.Query("isRecovered"); isRecovered != "" {
		if recovered, err := strconv.ParseBool(isRecovered); err == nil {
			filters["is_recovered"] = recovered
		}
	}

	response, err := h.cartService.GetAbandonmentAnalytics(c.Request.Context(), filters, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ApplyCoupon applies a coupon to the cart
// @Summary Apply coupon to cart
// @Description Apply a coupon code to reduce cart total
// @Tags carts
// @Accept json
// @Produce json
// @Param cartId path string true "Cart ID"
// @Param coupon body models.ApplyCouponRequest true "Coupon data"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/{cartId}/coupons [post]
func (h *CartHandler) ApplyCoupon(c *gin.Context) {
	cartIDStr := c.Param("cartId")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	var req models.ApplyCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.ApplyCoupon(c.Request.Context(), cartID, req.CouponCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated cart to return
	cart, err := h.cartService.GetCart(c.Request.Context(), nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// RemoveCoupon removes a coupon from the cart
// @Summary Remove coupon from cart
// @Description Remove the applied coupon from cart
// @Tags carts
// @Param cartId path string true "Cart ID"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/{cartId}/coupons [delete]
func (h *CartHandler) RemoveCoupon(c *gin.Context) {
	cartIDStr := c.Param("cartId")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	if err := h.cartService.RemoveCoupon(c.Request.Context(), cartID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated cart to return
	cart, err := h.cartService.GetCart(c.Request.Context(), nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// ValidateCart validates cart contents and pricing
// @Summary Validate cart
// @Description Validate cart items, pricing, and availability
// @Tags carts
// @Produce json
// @Param cartId path string true "Cart ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /carts/{cartId}/validate [get]
func (h *CartHandler) ValidateCart(c *gin.Context) {
	cartIDStr := c.Param("cartId")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	validation, err := h.cartService.ValidateCart(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, validation)
}