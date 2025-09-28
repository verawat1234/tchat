package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tchat.dev/payment/services"
	"tchat.dev/shared/responses"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService *services.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreateWallet handles wallet creation requests
func (h *PaymentHandler) CreateWallet(c *gin.Context) {
	var req services.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	wallet, err := h.paymentService.CreateWallet(c.Request.Context(), &req)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "wallet_creation_failed", "Failed to create wallet")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Wallet created successfully",
		"data":    wallet,
	})
}

// GetWallet handles wallet retrieval requests
func (h *PaymentHandler) GetWallet(c *gin.Context) {
	walletIDStr := c.Param("id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "invalid_wallet_id", "Invalid wallet ID")
		return
	}

	wallet, err := h.paymentService.GetWallet(c.Request.Context(), walletID)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusNotFound, "wallet_not_found", "Wallet not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Wallet retrieved successfully",
		"data":    wallet,
	})
}

// GetUserWallets handles user wallet list requests
func (h *PaymentHandler) GetUserWallets(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "invalid_user_id", "Invalid user ID")
		return
	}

	wallets, err := h.paymentService.GetUserWallets(c.Request.Context(), userID)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "wallets_retrieval_failed", "Failed to retrieve wallets")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Wallets retrieved successfully",
		"data":    wallets,
	})
}

// CreateTransaction handles transaction creation requests
func (h *PaymentHandler) CreateTransaction(c *gin.Context) {
	var req services.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	transaction, err := h.paymentService.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "transaction_creation_failed", "Failed to create transaction")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Transaction created successfully",
		"data":    transaction,
	})
}

// GetTransaction handles transaction retrieval requests
func (h *PaymentHandler) GetTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "invalid_transaction_id", "Invalid transaction ID")
		return
	}

	transaction, err := h.paymentService.GetTransaction(c.Request.Context(), transactionID)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusNotFound, "transaction_not_found", "Transaction not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transaction retrieved successfully",
		"data":    transaction,
	})
}

// GetWalletTransactions handles wallet transaction list requests
func (h *PaymentHandler) GetWalletTransactions(c *gin.Context) {
	walletIDStr := c.Param("walletId")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusBadRequest, "invalid_wallet_id", "Invalid wallet ID")
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	transactions, err := h.paymentService.GetWalletTransactions(c.Request.Context(), walletID, limit, offset)
	if err != nil {
		responses.SendErrorResponse(c, http.StatusInternalServerError, "transactions_retrieval_failed", "Failed to retrieve transactions")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transactions retrieved successfully",
		"data":    transactions,
	})
}

// RegisterRoutes registers payment routes
func (h *PaymentHandler) RegisterRoutes(router gin.IRouter) {
	payment := router.Group("/payment")
	{
		// Wallet routes
		payment.POST("/wallets", h.CreateWallet)
		payment.GET("/wallets/:id", h.GetWallet)
		payment.GET("/users/:userId/wallets", h.GetUserWallets)

		// Transaction routes
		payment.POST("/transactions", h.CreateTransaction)
		payment.GET("/transactions/:id", h.GetTransaction)
		payment.GET("/wallet/:walletId/transactions", h.GetWalletTransactions)
	}
}