package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"tchat-backend/payment/services"
	"tchat-backend/shared/middleware"
	"tchat-backend/shared/responses"
)

type WalletHandler struct {
	walletService *services.WalletService
	validator     *validator.Validate
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
		validator:     validator.New(),
	}
}

// CreateWalletRequest represents the request to create a new wallet
type CreateWalletRequest struct {
	Currency string `json:"currency" validate:"required,len=3" example:"THB"`
	Name     string `json:"name,omitempty" validate:"omitempty,min=1,max=50" example:"Main Wallet"`
}

// WalletResponse represents a wallet in API responses
type WalletResponse struct {
	ID               uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID           uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	Currency         string    `json:"currency" example:"THB"`
	Name             string    `json:"name" example:"Main Wallet"`
	AvailableBalance float64   `json:"available_balance" example:"1250.75"`
	FrozenBalance    float64   `json:"frozen_balance" example:"50.25"`
	TotalBalance     float64   `json:"total_balance" example:"1301.00"`
	Status           string    `json:"status" example:"active"`
	IsDefault        bool      `json:"is_default" example:"true"`
	CreatedAt        string    `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt        string    `json:"updated_at" example:"2024-01-20T15:45:00Z"`
	Limits           WalletLimits `json:"limits"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type WalletLimits struct {
	DailyTransactionLimit  float64 `json:"daily_transaction_limit" example:"50000.00"`
	MonthlyTransactionLimit float64 `json:"monthly_transaction_limit" example:"500000.00"`
	SingleTransactionLimit float64 `json:"single_transaction_limit" example:"10000.00"`
	DailyWithdrawalLimit   float64 `json:"daily_withdrawal_limit" example:"20000.00"`
	MonthlyWithdrawalLimit float64 `json:"monthly_withdrawal_limit" example:"200000.00"`
}

// WalletStatsResponse represents wallet statistics
type WalletStatsResponse struct {
	TotalBalance       float64 `json:"total_balance" example:"1301.00"`
	TotalTransactions  int     `json:"total_transactions" example:"47"`
	MonthlyInflow      float64 `json:"monthly_inflow" example:"3500.50"`
	MonthlyOutflow     float64 `json:"monthly_outflow" example:"2249.75"`
	DailyTransactionCount int  `json:"daily_transaction_count" example:"5"`
	DailyTransactionAmount float64 `json:"daily_transaction_amount" example:"875.25"`
	LastTransactionAt  string  `json:"last_transaction_at" example:"2024-01-20T14:30:00Z"`
}

// FreezeBalanceRequest represents the request to freeze wallet balance
type FreezeBalanceRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0" example:"100.00"`
	Reason      string  `json:"reason" validate:"required,max=200" example:"Payment authorization"`
	ExternalRef string  `json:"external_ref,omitempty" validate:"omitempty,max=100" example:"payment_12345"`
}

// UnfreezeBalanceRequest represents the request to unfreeze wallet balance
type UnfreezeBalanceRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0" example:"100.00"`
	Reason      string  `json:"reason" validate:"required,max=200" example:"Payment completed"`
	ExternalRef string  `json:"external_ref,omitempty" validate:"omitempty,max=100" example:"payment_12345"`
}

// SetDefaultWalletRequest represents the request to set default wallet
type SetDefaultWalletRequest struct {
	WalletID uuid.UUID `json:"wallet_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// @Summary Get user wallets
// @Description Get list of wallets for the authenticated user
// @Tags wallets
// @Produce json
// @Security BearerAuth
// @Param currency query string false "Filter by currency" example:"THB"
// @Param status query string false "Filter by status" Enums(active,suspended,closed)
// @Param include_stats query bool false "Include wallet statistics" default(false)
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse
// @Router /wallets [get]
func (h *WalletHandler) GetWallets(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Parse query parameters
	currency := c.Query("currency")
	status := c.Query("status")
	includeStats, _ := strconv.ParseBool(c.Query("include_stats"))

	// Build service request
	req := &services.GetWalletsRequest{
		UserID:       userUUID,
		Currency:     currency,
		Status:       status,
		IncludeStats: includeStats,
	}

	// Get wallets from service
	result, err := h.walletService.GetUserWallets(c.Request.Context(), req)
	if err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get wallets", "Unable to retrieve wallets.")
		return
	}

	// Convert to response format
	walletResponses := make([]WalletResponse, len(result.Wallets))
	for i, wallet := range result.Wallets {
		walletResponses[i] = h.convertToWalletResponse(wallet)
	}

	// Build response data
	data := gin.H{
		"wallets":       walletResponses,
		"total_balance": result.TotalBalance,
		"currencies":    result.SupportedCurrencies,
		"default_wallet_id": result.DefaultWalletID,
	}

	// Add stats if requested
	if includeStats && result.Stats != nil {
		data["stats"] = WalletStatsResponse{
			TotalBalance:          result.Stats.TotalBalance,
			TotalTransactions:     result.Stats.TotalTransactions,
			MonthlyInflow:         result.Stats.MonthlyInflow,
			MonthlyOutflow:        result.Stats.MonthlyOutflow,
			DailyTransactionCount: result.Stats.DailyTransactionCount,
			DailyTransactionAmount: result.Stats.DailyTransactionAmount,
			LastTransactionAt:     result.Stats.LastTransactionAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	responses.DataResponse(c, data)
}

// @Summary Create a new wallet
// @Description Create a new wallet for the authenticated user
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateWalletRequest true "Wallet creation request"
// @Success 201 {object} responses.DataResponse{data=WalletResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 409 {object} responses.ErrorResponse "Wallet already exists"
// @Failure 500 {object} responses.ErrorResponse
// @Router /wallets [post]
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req CreateWalletRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Normalize currency
	req.Currency = strings.ToUpper(strings.TrimSpace(req.Currency))

	// Build service request
	serviceReq := &services.CreateWalletRequest{
		UserID:   userUUID,
		Currency: req.Currency,
		Name:     req.Name,
	}

	// Create wallet
	wallet, err := h.walletService.CreateWallet(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "wallet already exists"):
			responses.ErrorResponse(c, http.StatusConflict, "Wallet already exists", "A wallet with this currency already exists for this user.")
			return
		case strings.Contains(err.Error(), "unsupported currency"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Unsupported currency", "The specified currency is not supported.")
			return
		case strings.Contains(err.Error(), "wallet limit exceeded"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Wallet limit exceeded", "Maximum number of wallets reached for this user.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Wallet creation failed", "Failed to create wallet.")
			return
		}
	}

	// Convert to response format
	walletResponse := h.convertToWalletResponse(wallet)

	// Log wallet creation
	middleware.LogInfo(c, "Wallet created", gin.H{
		"wallet_id": wallet.ID,
		"user_id":   userUUID,
		"currency":  req.Currency,
		"name":      req.Name,
	})

	c.JSON(http.StatusCreated, responses.DataResponse{
		Success: true,
		Data:    walletResponse,
	})
}

// @Summary Get wallet by ID
// @Description Get detailed information about a specific wallet
// @Tags wallets
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Param include_transactions query bool false "Include recent transactions" default(false)
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Wallet not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /wallets/{id} [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {
	// Parse wallet ID
	walletIDStr := c.Param("id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid wallet ID", "Wallet ID must be a valid UUID.")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Parse query parameters
	includeTransactions, _ := strconv.ParseBool(c.Query("include_transactions"))

	// Get wallet from service
	wallet, err := h.walletService.GetWalletByID(c.Request.Context(), walletID, userUUID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Wallet not found", "The specified wallet does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have access to this wallet.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get wallet", "Unable to retrieve wallet.")
			return
		}
	}

	// Convert to response format
	walletResponse := h.convertToWalletResponse(wallet)

	// Build response data
	data := gin.H{
		"wallet": walletResponse,
	}

	// Include recent transactions if requested
	if includeTransactions {
		transactions, err := h.walletService.GetWalletTransactions(c.Request.Context(), walletID, userUUID, 10, 0)
		if err == nil {
			data["recent_transactions"] = transactions
		}
	}

	responses.DataResponse(c, data)
}

// @Summary Freeze wallet balance
// @Description Freeze a specific amount in the wallet balance
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Param request body FreezeBalanceRequest true "Freeze balance request"
// @Success 200 {object} responses.DataResponse{data=WalletResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Wallet not found"
// @Failure 422 {object} responses.ErrorResponse "Insufficient balance"
// @Failure 500 {object} responses.ErrorResponse
// @Router /wallets/{id}/freeze [post]
func (h *WalletHandler) FreezeBalance(c *gin.Context) {
	var req FreezeBalanceRequest

	// Parse wallet ID
	walletIDStr := c.Param("id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid wallet ID", "Wallet ID must be a valid UUID.")
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.FreezeBalanceRequest{
		WalletID:    walletID,
		UserID:      userUUID,
		Amount:      req.Amount,
		Reason:      req.Reason,
		ExternalRef: req.ExternalRef,
	}

	// Freeze balance
	wallet, err := h.walletService.FreezeBalance(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Wallet not found", "The specified wallet does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have access to this wallet.")
			return
		case strings.Contains(err.Error(), "insufficient balance"):
			responses.ErrorResponse(c, http.StatusUnprocessableEntity, "Insufficient balance", "Not enough available balance to freeze the requested amount.")
			return
		case strings.Contains(err.Error(), "wallet suspended"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Wallet suspended", "This wallet is currently suspended.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to freeze balance", "Unable to freeze balance.")
			return
		}
	}

	// Convert to response format
	walletResponse := h.convertToWalletResponse(wallet)

	// Log balance freeze
	middleware.LogInfo(c, "Balance frozen", gin.H{
		"wallet_id":    walletID,
		"user_id":      userUUID,
		"amount":       req.Amount,
		"reason":       req.Reason,
		"external_ref": req.ExternalRef,
	})

	responses.DataResponse(c, walletResponse)
}

// @Summary Unfreeze wallet balance
// @Description Unfreeze a specific amount in the wallet balance
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Param request body UnfreezeBalanceRequest true "Unfreeze balance request"
// @Success 200 {object} responses.DataResponse{data=WalletResponse}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Wallet not found"
// @Failure 422 {object} responses.ErrorResponse "Insufficient frozen balance"
// @Failure 500 {object} responses.ErrorResponse
// @Router /wallets/{id}/unfreeze [post]
func (h *WalletHandler) UnfreezeBalance(c *gin.Context) {
	var req UnfreezeBalanceRequest

	// Parse wallet ID
	walletIDStr := c.Param("id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid wallet ID", "Wallet ID must be a valid UUID.")
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Build service request
	serviceReq := &services.UnfreezeBalanceRequest{
		WalletID:    walletID,
		UserID:      userUUID,
		Amount:      req.Amount,
		Reason:      req.Reason,
		ExternalRef: req.ExternalRef,
	}

	// Unfreeze balance
	wallet, err := h.walletService.UnfreezeBalance(c.Request.Context(), serviceReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Wallet not found", "The specified wallet does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have access to this wallet.")
			return
		case strings.Contains(err.Error(), "insufficient frozen balance"):
			responses.ErrorResponse(c, http.StatusUnprocessableEntity, "Insufficient frozen balance", "Not enough frozen balance to unfreeze the requested amount.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to unfreeze balance", "Unable to unfreeze balance.")
			return
		}
	}

	// Convert to response format
	walletResponse := h.convertToWalletResponse(wallet)

	// Log balance unfreeze
	middleware.LogInfo(c, "Balance unfrozen", gin.H{
		"wallet_id":    walletID,
		"user_id":      userUUID,
		"amount":       req.Amount,
		"reason":       req.Reason,
		"external_ref": req.ExternalRef,
	})

	responses.DataResponse(c, walletResponse)
}

// @Summary Set default wallet
// @Description Set a wallet as the default for the user
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SetDefaultWalletRequest true "Set default wallet request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Wallet not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /wallets/default [post]
func (h *WalletHandler) SetDefaultWallet(c *gin.Context) {
	var req SetDefaultWalletRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		responses.ValidationErrorResponse(c, err)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Set default wallet
	err := h.walletService.SetDefaultWallet(c.Request.Context(), req.WalletID, userUUID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Wallet not found", "The specified wallet does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have access to this wallet.")
			return
		case strings.Contains(err.Error(), "wallet not active"):
			responses.ErrorResponse(c, http.StatusBadRequest, "Wallet not active", "Only active wallets can be set as default.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to set default wallet", "Unable to set default wallet.")
			return
		}
	}

	// Log default wallet change
	middleware.LogInfo(c, "Default wallet set", gin.H{
		"wallet_id": req.WalletID,
		"user_id":   userUUID,
	})

	responses.SuccessMessageResponse(c, "Default wallet set successfully")
}

// @Summary Get wallet statistics
// @Description Get detailed statistics for a specific wallet
// @Tags wallets
// @Produce json
// @Security BearerAuth
// @Param id path string true "Wallet ID"
// @Param period query string false "Statistics period" Enums(day,week,month,year) default(month)
// @Success 200 {object} responses.DataResponse{data=map[string]interface{}}
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Access denied"
// @Failure 404 {object} responses.ErrorResponse "Wallet not found"
// @Failure 500 {object} responses.ErrorResponse
// @Router /wallets/{id}/stats [get]
func (h *WalletHandler) GetWalletStats(c *gin.Context) {
	// Parse wallet ID
	walletIDStr := c.Param("id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		responses.ErrorResponse(c, http.StatusBadRequest, "Invalid wallet ID", "Wallet ID must be a valid UUID.")
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Authentication required.")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		responses.ErrorResponse(c, http.StatusInternalServerError, "Invalid user context", "Invalid user context.")
		return
	}

	// Parse query parameters
	period := c.DefaultQuery("period", "month")

	// Get wallet statistics
	stats, err := h.walletService.GetWalletStats(c.Request.Context(), walletID, userUUID, period)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			responses.ErrorResponse(c, http.StatusNotFound, "Wallet not found", "The specified wallet does not exist.")
			return
		case strings.Contains(err.Error(), "access denied"):
			responses.ErrorResponse(c, http.StatusForbidden, "Access denied", "You don't have access to this wallet.")
			return
		default:
			responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get statistics", "Unable to retrieve wallet statistics.")
			return
		}
	}

	// Convert to response format
	statsResponse := WalletStatsResponse{
		TotalBalance:           stats.TotalBalance,
		TotalTransactions:      stats.TotalTransactions,
		MonthlyInflow:          stats.MonthlyInflow,
		MonthlyOutflow:         stats.MonthlyOutflow,
		DailyTransactionCount:  stats.DailyTransactionCount,
		DailyTransactionAmount: stats.DailyTransactionAmount,
		LastTransactionAt:      stats.LastTransactionAt.Format("2006-01-02T15:04:05Z"),
	}

	responses.DataResponse(c, statsResponse)
}

// Helper functions

func (h *WalletHandler) convertToWalletResponse(wallet interface{}) WalletResponse {
	// This is a simplified conversion - in a real implementation,
	// you would properly convert from your wallet model to the response
	return WalletResponse{
		// Populate fields from wallet model
		// This is just a placeholder structure
		Limits: WalletLimits{
			DailyTransactionLimit:   50000.00,
			MonthlyTransactionLimit: 500000.00,
			SingleTransactionLimit:  10000.00,
			DailyWithdrawalLimit:    20000.00,
			MonthlyWithdrawalLimit:  200000.00,
		},
	}
}

// RegisterWalletRoutes registers all wallet-related routes
func RegisterWalletRoutes(router *gin.RouterGroup, walletService *services.WalletService) {
	handler := NewWalletHandler(walletService)

	// Protected wallet routes
	wallets := router.Group("/wallets")
	wallets.Use(middleware.AuthRequired())
	{
		wallets.GET("", handler.GetWallets)
		wallets.POST("", handler.CreateWallet)
		wallets.GET("/:id", handler.GetWallet)
		wallets.POST("/:id/freeze", handler.FreezeBalance)
		wallets.POST("/:id/unfreeze", handler.UnfreezeBalance)
		wallets.POST("/default", handler.SetDefaultWallet)
		wallets.GET("/:id/stats", handler.GetWalletStats)
	}
}

// RegisterWalletRoutesWithMiddleware registers wallet routes with custom middleware
func RegisterWalletRoutesWithMiddleware(
	router *gin.RouterGroup,
	walletService *services.WalletService,
	middlewares ...gin.HandlerFunc,
) {
	handler := NewWalletHandler(walletService)

	// Protected wallet routes with middleware
	wallets := router.Group("/wallets")
	allMiddlewares := append(middlewares, middleware.AuthRequired())
	wallets.Use(allMiddlewares...)
	{
		wallets.GET("", handler.GetWallets)
		wallets.POST("", handler.CreateWallet)
		wallets.GET("/:id", handler.GetWallet)
		wallets.POST("/:id/freeze", handler.FreezeBalance)
		wallets.POST("/:id/unfreeze", handler.UnfreezeBalance)
		wallets.POST("/default", handler.SetDefaultWallet)
		wallets.GET("/:id/stats", handler.GetWalletStats)
	}
}