package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/services"
	"tchat.dev/shared/utils"
)

// CommerceHandler handles commerce HTTP requests
type CommerceHandler struct {
	commerceService services.CommerceService
	validator       *utils.Validator
}

// NewCommerceHandler creates a new commerce handler
func NewCommerceHandler(commerceService services.CommerceService) *CommerceHandler {
	return &CommerceHandler{
		commerceService: commerceService,
		validator:       utils.NewValidator(),
	}
}

// RegisterRoutes registers commerce routes
func (h *CommerceHandler) RegisterRoutes(router *mux.Router) {
	// Commerce routes
	commerce := router.PathPrefix("/commerce").Subrouter()

	// Shop endpoints
	commerce.HandleFunc("/shops", h.CreateShop).Methods("POST")
	commerce.HandleFunc("/shops", h.SearchShops).Methods("GET")
	commerce.HandleFunc("/shops/{id}", h.GetShop).Methods("GET")
	commerce.HandleFunc("/shops/{id}", h.UpdateShop).Methods("PUT")
	commerce.HandleFunc("/shops/{id}", h.DeleteShop).Methods("DELETE")
	commerce.HandleFunc("/shops/my", h.GetMyShops).Methods("GET")
	commerce.HandleFunc("/shops/featured", h.GetFeaturedShops).Methods("GET")

	// Product endpoints
	commerce.HandleFunc("/products", h.SearchProducts).Methods("GET")
	commerce.HandleFunc("/products", h.CreateProduct).Methods("POST")
	commerce.HandleFunc("/products/{id}", h.GetProduct).Methods("GET")
	commerce.HandleFunc("/products/{id}", h.UpdateProduct).Methods("PUT")
	commerce.HandleFunc("/products/{id}", h.DeleteProduct).Methods("DELETE")
	commerce.HandleFunc("/shops/{shop_id}/products", h.GetShopProducts).Methods("GET")
	commerce.HandleFunc("/categories/{category}/products", h.GetProductsByCategory).Methods("GET")

	// Order endpoints
	commerce.HandleFunc("/orders", h.CreateOrder).Methods("POST")
	commerce.HandleFunc("/orders", h.GetUserOrders).Methods("GET")
	commerce.HandleFunc("/orders/{id}", h.GetOrder).Methods("GET")
	commerce.HandleFunc("/orders/{id}/cancel", h.CancelOrder).Methods("POST")
	commerce.HandleFunc("/orders/{id}/payment", h.ProcessPayment).Methods("POST")
	commerce.HandleFunc("/orders/{id}/confirm", h.ConfirmPayment).Methods("POST")
	commerce.HandleFunc("/shops/{shop_id}/orders", h.GetShopOrders).Methods("GET")

	// Public endpoints (no auth required)
	public := router.PathPrefix("/public/commerce").Subrouter()
	public.HandleFunc("/shops", h.GetPublicShops).Methods("GET")
	public.HandleFunc("/shops/{id}", h.GetPublicShop).Methods("GET")
	public.HandleFunc("/products", h.GetPublicProducts).Methods("GET")
	public.HandleFunc("/products/{id}", h.GetPublicProduct).Methods("GET")
}

// Shop endpoints

// CreateShop handles shop creation
func (h *CommerceHandler) CreateShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req CreateShopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateCreateShopRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Create service request
	serviceReq := &services.CreateShopRequest{
		OwnerID:     user.ID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Logo:        req.Logo,
		Banner:      req.Banner,
		Address:     req.Address,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
	}

	// Create shop
	shop, err := h.commerceService.CreateShop(ctx, serviceReq)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to create shop", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusCreated, "Shop created successfully", h.sanitizeShop(shop))
}

// GetShop handles getting a specific shop
func (h *CommerceHandler) GetShop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get shop ID from URL
	vars := mux.Vars(r)
	shopIDStr := vars["id"]
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid shop ID", err)
		return
	}

	// Get shop
	shop, err := h.commerceService.GetShop(ctx, shopID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Shop not found", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Shop retrieved successfully", h.sanitizeShop(shop))
}

// SearchShops handles shop search with filters
func (h *CommerceHandler) SearchShops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get search parameters
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	location := r.URL.Query().Get("location")
	verified := r.URL.Query().Get("verified")

	// Get pagination parameters
	limit, offset := h.getPaginationParams(r)

	// Build filters
	filters := make(map[string]interface{})
	if category != "" {
		filters["category"] = category
	}
	if location != "" {
		filters["location"] = location
	}
	if verified != "" {
		if isVerified, err := strconv.ParseBool(verified); err == nil {
			filters["is_verified"] = isVerified
		}
	}

	// Search shops (this would need to be implemented in the service)
	shops := []*models.Shop{} // Placeholder

	// Sanitize shops
	var sanitizedShops []map[string]interface{}
	for _, shop := range shops {
		sanitizedShops = append(sanitizedShops, h.sanitizeShop(shop))
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Shops retrieved successfully", map[string]interface{}{
		"shops":   sanitizedShops,
		"count":   len(sanitizedShops),
		"query":   query,
		"filters": filters,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetMyShops handles getting user's shops
func (h *CommerceHandler) GetMyShops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get user's shops (this would need to be implemented in the service)
	shops := []*models.Shop{} // Placeholder

	// Sanitize shops
	var sanitizedShops []map[string]interface{}
	for _, shop := range shops {
		sanitizedShops = append(sanitizedShops, h.sanitizeShop(shop))
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Shops retrieved successfully", map[string]interface{}{
		"shops": sanitizedShops,
		"count": len(sanitizedShops),
	})
}

// Product endpoints

// CreateProduct handles product creation
func (h *CommerceHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateCreateProductRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Verify shop ownership (this would need to be implemented)
	shop, err := h.commerceService.GetShop(ctx, req.ShopID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Shop not found", err)
		return
	}

	if shop.OwnerID != user.ID {
		h.respondError(w, http.StatusForbidden, "Access denied", nil)
		return
	}

	// Create service request
	serviceReq := &services.CreateProductRequest{
		ShopID:      req.ShopID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Currency:    req.Currency,
		Inventory:   req.Inventory,
		Images:      req.Images,
		Variants:    req.Variants,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
	}

	// Create product
	product, err := h.commerceService.CreateProduct(ctx, serviceReq)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to create product", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusCreated, "Product created successfully", h.sanitizeProduct(product))
}

// SearchProducts handles product search with filters
func (h *CommerceHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get search parameters
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	shopID := r.URL.Query().Get("shop_id")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	sortBy := r.URL.Query().Get("sort")
	currency := r.URL.Query().Get("currency")

	// Get pagination parameters
	limit, offset := h.getPaginationParams(r)

	// Build filters
	filters := make(map[string]interface{})
	if category != "" {
		filters["category"] = category
	}
	if shopID != "" {
		if shopUUID, err := uuid.Parse(shopID); err == nil {
			filters["shop_id"] = shopUUID
		}
	}
	if minPrice != "" {
		if price, err := strconv.ParseInt(minPrice, 10, 64); err == nil {
			filters["min_price"] = price
		}
	}
	if maxPrice != "" {
		if price, err := strconv.ParseInt(maxPrice, 10, 64); err == nil {
			filters["max_price"] = price
		}
	}
	if currency != "" {
		filters["currency"] = currency
	}
	if sortBy != "" {
		filters["sort"] = sortBy
	}

	// Search products
	products, err := h.commerceService.SearchProducts(ctx, query, filters, limit, offset)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to search products", err)
		return
	}

	// Sanitize products
	var sanitizedProducts []map[string]interface{}
	for _, product := range products {
		sanitizedProducts = append(sanitizedProducts, h.sanitizeProduct(product))
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Products retrieved successfully", map[string]interface{}{
		"products": sanitizedProducts,
		"count":    len(sanitizedProducts),
		"query":    query,
		"filters":  filters,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetShopProducts handles getting products for a specific shop
func (h *CommerceHandler) GetShopProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get shop ID from URL
	vars := mux.Vars(r)
	shopIDStr := vars["shop_id"]
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid shop ID", err)
		return
	}

	// Get pagination parameters
	limit, offset := h.getPaginationParams(r)

	// Get shop products
	products, err := h.commerceService.GetShopProducts(ctx, shopID, limit, offset)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get shop products", err)
		return
	}

	// Sanitize products
	var sanitizedProducts []map[string]interface{}
	for _, product := range products {
		sanitizedProducts = append(sanitizedProducts, h.sanitizeProduct(product))
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Products retrieved successfully", map[string]interface{}{
		"products": sanitizedProducts,
		"count":    len(sanitizedProducts),
		"shop_id":  shopID,
		"limit":    limit,
		"offset":   offset,
	})
}

// Order endpoints

// CreateOrder handles order creation
func (h *CommerceHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateCreateOrderRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Create service request
	serviceReq := &services.CreateOrderRequest{
		UserID:          user.ID,
		ShopID:          req.ShopID,
		Items:           req.Items,
		Currency:        req.Currency,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Notes:           req.Notes,
	}

	// Create order
	order, err := h.commerceService.CreateOrder(ctx, serviceReq)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to create order", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusCreated, "Order created successfully", h.sanitizeOrder(order))
}

// ProcessPayment handles order payment processing
func (h *CommerceHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var req ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateProcessPaymentRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Verify order ownership
	order, err := h.commerceService.GetOrder(ctx, orderID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Order not found", err)
		return
	}

	if order.UserID != user.ID {
		h.respondError(w, http.StatusForbidden, "Access denied", nil)
		return
	}

	// Process payment
	paymentResp, err := h.commerceService.ProcessPayment(ctx, orderID, req.PaymentMethod)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Payment processing failed", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Payment initiated successfully", map[string]interface{}{
		"payment_id":   paymentResp.PaymentID,
		"status":       paymentResp.Status,
		"redirect_url": paymentResp.RedirectURL,
		"expires_at":   paymentResp.ExpiresAt,
	})
}

// CancelOrder handles order cancellation
func (h *CommerceHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var req CancelOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := h.validateCancelOrderRequest(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Cancel order
	if err := h.commerceService.CancelOrder(ctx, orderID, user.ID, req.Reason); err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to cancel order", err)
		return
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Order cancelled successfully", map[string]interface{}{
		"order_id":  orderID,
		"cancelled": true,
		"reason":    req.Reason,
	})
}

// GetUserOrders handles getting user orders
func (h *CommerceHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := h.getUserFromContext(ctx)
	if user == nil {
		h.respondError(w, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	// Get pagination parameters
	limit, offset := h.getPaginationParams(r)

	// Get filter parameters
	status := r.URL.Query().Get("status")
	shopID := r.URL.Query().Get("shop_id")

	// Get orders
	orders, err := h.commerceService.GetUserOrders(ctx, user.ID, limit, offset)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get orders", err)
		return
	}

	// Apply filters
	var filteredOrders []*models.Order
	for _, order := range orders {
		include := true

		if status != "" && string(order.Status) != status {
			include = false
		}

		if shopID != "" {
			if shopUUID, err := uuid.Parse(shopID); err == nil && order.ShopID != shopUUID {
				include = false
			}
		}

		if include {
			filteredOrders = append(filteredOrders, order)
		}
	}

	// Sanitize orders
	var sanitizedOrders []map[string]interface{}
	for _, order := range filteredOrders {
		sanitizedOrders = append(sanitizedOrders, h.sanitizeOrder(order))
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Orders retrieved successfully", map[string]interface{}{
		"orders": sanitizedOrders,
		"count":  len(sanitizedOrders),
		"limit":  limit,
		"offset": offset,
	})
}

// Public endpoints (no authentication required)

// GetPublicShops handles getting public shops
func (h *CommerceHandler) GetPublicShops(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get pagination parameters
	limit, offset := h.getPaginationParams(r)

	// Get public shops (this would need to be implemented in the service)
	shops := []*models.Shop{} // Placeholder

	// Sanitize shops (public view)
	var sanitizedShops []map[string]interface{}
	for _, shop := range shops {
		sanitizedShops = append(sanitizedShops, h.sanitizeShopPublic(shop))
	}

	// Success response
	h.respondSuccess(w, http.StatusOK, "Public shops retrieved successfully", map[string]interface{}{
		"shops":  sanitizedShops,
		"count":  len(sanitizedShops),
		"limit":  limit,
		"offset": offset,
	})
}

// Request/Response types
type CreateShopRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Category    string  `json:"category"`
	Logo        *string `json:"logo,omitempty"`
	Banner      *string `json:"banner,omitempty"`
	Address     *string `json:"address,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Email       *string `json:"email,omitempty"`
	Website     *string `json:"website,omitempty"`
}

type CreateProductRequest struct {
	ShopID      uuid.UUID                  `json:"shop_id"`
	Name        string                     `json:"name"`
	Description *string                    `json:"description,omitempty"`
	Category    string                     `json:"category"`
	Price       int64                      `json:"price"`
	Currency    models.Currency            `json:"currency"`
	Inventory   int                        `json:"inventory"`
	Images      []string                   `json:"images,omitempty"`
	Variants    []models.ProductVariant    `json:"variants,omitempty"`
	Tags        []string                   `json:"tags,omitempty"`
	Metadata    map[string]string          `json:"metadata,omitempty"`
}

type CreateOrderRequest struct {
	ShopID          uuid.UUID              `json:"shop_id"`
	Items           []services.OrderItem   `json:"items"`
	Currency        models.Currency        `json:"currency"`
	ShippingAddress *models.Address        `json:"shipping_address,omitempty"`
	BillingAddress  *models.Address        `json:"billing_address,omitempty"`
	Notes           *string                `json:"notes,omitempty"`
}

type ProcessPaymentRequest struct {
	PaymentMethod string `json:"payment_method"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

// Standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Validation methods
func (h *CommerceHandler) validateCreateShopRequest(req *CreateShopRequest) error {
	h.validator.Reset()

	h.validator.Required("name", req.Name).MinLength("name", req.Name, 2).MaxLength("name", req.Name, 100)
	h.validator.Required("category", req.Category).MinLength("category", req.Category, 2)

	if req.Description != nil {
		h.validator.MaxLength("description", *req.Description, 1000)
	}

	if req.Email != nil {
		h.validator.Email("email", *req.Email)
	}

	if req.Phone != nil {
		h.validator.Phone("phone", *req.Phone)
	}

	if req.Website != nil {
		h.validator.URL("website", *req.Website)
	}

	return h.validator.GetError()
}

func (h *CommerceHandler) validateCreateProductRequest(req *CreateProductRequest) error {
	h.validator.Reset()

	if req.ShopID == uuid.Nil {
		h.validator.AddError("shop_id", "shop ID is required")
	}

	h.validator.Required("name", req.Name).MinLength("name", req.Name, 2).MaxLength("name", req.Name, 200)
	h.validator.Required("category", req.Category).MinLength("category", req.Category, 2)
	h.validator.Positive("price", req.Price)
	h.validator.NonNegative("inventory", int64(req.Inventory))

	if !req.Currency.IsValid() {
		h.validator.AddError("currency", "invalid currency")
	}

	if req.Description != nil {
		h.validator.MaxLength("description", *req.Description, 2000)
	}

	return h.validator.GetError()
}

func (h *CommerceHandler) validateCreateOrderRequest(req *CreateOrderRequest) error {
	h.validator.Reset()

	if req.ShopID == uuid.Nil {
		h.validator.AddError("shop_id", "shop ID is required")
	}

	if len(req.Items) == 0 {
		h.validator.AddError("items", "order must have at least one item")
	}

	for i, item := range req.Items {
		if item.ProductID == uuid.Nil {
			h.validator.AddErrorf("items[%d]", "product ID is required", i)
		}
		if item.Quantity <= 0 {
			h.validator.AddErrorf("items[%d]", "quantity must be positive", i)
		}
		if item.Price <= 0 {
			h.validator.AddErrorf("items[%d]", "price must be positive", i)
		}
	}

	if !req.Currency.IsValid() {
		h.validator.AddError("currency", "invalid currency")
	}

	return h.validator.GetError()
}

func (h *CommerceHandler) validateProcessPaymentRequest(req *ProcessPaymentRequest) error {
	h.validator.Reset()

	h.validator.Required("payment_method", req.PaymentMethod)

	return h.validator.GetError()
}

func (h *CommerceHandler) validateCancelOrderRequest(req *CancelOrderRequest) error {
	h.validator.Reset()

	h.validator.Required("reason", req.Reason).MinLength("reason", req.Reason, 5).MaxLength("reason", req.Reason, 500)

	return h.validator.GetError()
}

// Utility methods
func (h *CommerceHandler) respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *CommerceHandler) respondError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var errorData interface{}
	if err != nil {
		errorData = err.Error()
	}

	response := APIResponse{
		Success: false,
		Message: message,
		Error:   errorData,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *CommerceHandler) getUserFromContext(ctx context.Context) *models.User {
	if user, ok := ctx.Value("user").(*models.User); ok {
		return user
	}
	return nil
}

func (h *CommerceHandler) sanitizeShop(shop *models.Shop) map[string]interface{} {
	if shop == nil {
		return nil
	}

	return map[string]interface{}{
		"id":            shop.ID,
		"owner_id":      shop.OwnerID,
		"name":          shop.Name,
		"slug":          shop.Slug,
		"description":   shop.Description,
		"category":      shop.Category,
		"logo":          shop.Logo,
		"banner":        shop.Banner,
		"address":       shop.Address,
		"phone":         shop.Phone,
		"email":         shop.Email,
		"website":       shop.Website,
		"status":        shop.Status,
		"is_verified":   shop.IsVerified,
		"product_count": shop.ProductCount,
		"order_count":   shop.OrderCount,
		"rating":        shop.Rating,
		"created_at":    shop.CreatedAt,
		"updated_at":    shop.UpdatedAt,
	}
}

func (h *CommerceHandler) sanitizeShopPublic(shop *models.Shop) map[string]interface{} {
	if shop == nil {
		return nil
	}

	return map[string]interface{}{
		"id":            shop.ID,
		"name":          shop.Name,
		"slug":          shop.Slug,
		"description":   shop.Description,
		"category":      shop.Category,
		"logo":          shop.Logo,
		"banner":        shop.Banner,
		"is_verified":   shop.IsVerified,
		"product_count": shop.ProductCount,
		"rating":        shop.Rating,
		"created_at":    shop.CreatedAt,
	}
}

func (h *CommerceHandler) sanitizeProduct(product *models.Product) map[string]interface{} {
	if product == nil {
		return nil
	}

	return map[string]interface{}{
		"id":               product.ID,
		"shop_id":          product.ShopID,
		"name":             product.Name,
		"slug":             product.Slug,
		"description":      product.Description,
		"category":         product.Category,
		"price":            product.Price,
		"formatted_price":  product.GetFormattedPrice(),
		"currency":         product.Currency,
		"inventory":        product.Inventory,
		"status":           product.Status,
		"images":           product.Images,
		"variants":         product.Variants,
		"tags":             product.Tags,
		"metadata":         product.Metadata,
		"rating":           product.Rating,
		"review_count":     product.ReviewCount,
		"created_at":       product.CreatedAt,
		"updated_at":       product.UpdatedAt,
	}
}

func (h *CommerceHandler) sanitizeOrder(order *models.Order) map[string]interface{} {
	if order == nil {
		return nil
	}

	return map[string]interface{}{
		"id":               order.ID,
		"user_id":          order.UserID,
		"shop_id":          order.ShopID,
		"items":            order.Items,
		"status":           order.Status,
		"currency":         order.Currency,
		"subtotal":         order.Subtotal,
		"tax":              order.Tax,
		"shipping_cost":    order.ShippingCost,
		"commission":       order.Commission,
		"total":            order.Total,
		"formatted_total":  order.GetFormattedTotal(),
		"payment_id":       order.PaymentID,
		"payment_method":   order.PaymentMethod,
		"shipping_address": order.ShippingAddress,
		"billing_address":  order.BillingAddress,
		"notes":            order.Notes,
		"paid_at":          order.PaidAt,
		"shipped_at":       order.ShippedAt,
		"delivered_at":     order.DeliveredAt,
		"cancelled_at":     order.CancelledAt,
		"cancellation_reason": order.CancellationReason,
		"created_at":       order.CreatedAt,
		"updated_at":       order.UpdatedAt,
	}
}

func (h *CommerceHandler) getPaginationParams(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 20 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset = 0 // default
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}

// Placeholder methods that would need full implementation
func (h *CommerceHandler) UpdateShop(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) DeleteShop(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) GetFeaturedShops(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) GetShopOrders(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) GetPublicShop(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) GetPublicProducts(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}

func (h *CommerceHandler) GetPublicProduct(w http.ResponseWriter, r *http.Request) {
	h.respondError(w, http.StatusNotImplemented, "Not implemented", nil)
}