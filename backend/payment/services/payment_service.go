package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	sharedModels "tchat.dev/shared/models"
)

// PaymentService handles payment operations
type PaymentService struct {
	walletRepo    WalletRepository
	txRepo        TransactionRepository
	eventPublisher EventPublisher
	db            *gorm.DB
}

// Repositories
type WalletRepository interface {
	Create(ctx context.Context, wallet *sharedModels.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*sharedModels.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*sharedModels.Wallet, error)
	GetByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency string) (*sharedModels.Wallet, error)
	Update(ctx context.Context, wallet *sharedModels.Wallet) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal, txID uuid.UUID) error
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *sharedModels.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*sharedModels.Transaction, error)
	GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*sharedModels.Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*sharedModels.Transaction, error)
	Update(ctx context.Context, tx *sharedModels.Transaction) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event *sharedModels.Event) error
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(
	walletRepo WalletRepository,
	txRepo TransactionRepository,
	eventPublisher EventPublisher,
	db *gorm.DB,
) *PaymentService {
	return &PaymentService{
		walletRepo:     walletRepo,
		txRepo:         txRepo,
		eventPublisher: eventPublisher,
		db:             db,
	}
}

// Request/Response types
type CreateWalletRequest struct {
	UserID            uuid.UUID              `json:"user_id" binding:"required"`
	Currency          string                 `json:"currency" binding:"required"`
	DailySpendLimit   decimal.Decimal        `json:"daily_spend_limit"`
	MonthlySpendLimit decimal.Decimal        `json:"monthly_spend_limit"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type CreateTransactionRequest struct {
	WalletID    uuid.UUID              `json:"wallet_id" binding:"required"`
	Amount      decimal.Decimal        `json:"amount" binding:"required"`
	Currency    string                 `json:"currency" binding:"required"`
	Type        string                 `json:"type" binding:"required"`
	Description string                 `json:"description"`
	Gateway     string                 `json:"gateway"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type WalletResponse struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	Currency    string                 `json:"currency"`
	Balance     decimal.Decimal        `json:"balance"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type TransactionResponse struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	WalletID    uuid.UUID              `json:"wallet_id"`
	Amount      decimal.Decimal        `json:"amount"`
	Currency    string                 `json:"currency"`
	Type        string                 `json:"type"`
	Status      string                 `json:"status"`
	Description string                 `json:"description"`
	Gateway     string                 `json:"gateway"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CreateWallet creates a new wallet for a user
func (ps *PaymentService) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*WalletResponse, error) {
	wallet := &sharedModels.Wallet{
		ID:     uuid.New(),
		UserID: &req.UserID,
		Name:   fmt.Sprintf("%s Wallet", req.Currency),
		Type:   sharedModels.WalletTypePersonal,
		Status: sharedModels.WalletStatusActive,
		Settings: sharedModels.WalletSettings{
			DefaultCurrency:   req.Currency,
			AllowedCurrencies: []string{req.Currency},
		},
		Security: sharedModels.WalletSecurity{
			SecurityLevel:          "basic",
			DailyWithdrawalLimit:   req.DailySpendLimit,
			MonthlyWithdrawalLimit: req.MonthlySpendLimit,
		},
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ps.walletRepo.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Publish wallet created event
	eventData, _ := json.Marshal(map[string]interface{}{
		"user_id":  wallet.UserID,
		"currency": wallet.Settings.DefaultCurrency,
	})

	event := &sharedModels.Event{
		ID:           uuid.New(),
		Type:         sharedModels.EventType("wallet.created"),
		Category:     sharedModels.EventCategoryDomain,
		Severity:     sharedModels.SeverityInfo,
		Subject:      "Wallet Created",
		Description:  fmt.Sprintf("New wallet created for user %s", wallet.UserID),
		Data:         json.RawMessage(eventData),
		AggregateID:  wallet.ID.String(),
		OccurredAt:   time.Now(),
	}

	if err := ps.eventPublisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish wallet created event: %v\n", err)
	}

	return &WalletResponse{
		ID:        wallet.ID,
		UserID:    *wallet.UserID,
		Currency:  wallet.Settings.DefaultCurrency,
		Balance:   wallet.GetTotalBalance(wallet.Settings.DefaultCurrency),
		Status:    string(wallet.Status),
		Metadata:  wallet.Metadata,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}, nil
}

// GetWallet retrieves a wallet by ID
func (ps *PaymentService) GetWallet(ctx context.Context, walletID uuid.UUID) (*WalletResponse, error) {
	wallet, err := ps.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	var userID uuid.UUID
	if wallet.UserID != nil {
		userID = *wallet.UserID
	}

	return &WalletResponse{
		ID:        wallet.ID,
		UserID:    userID,
		Currency:  wallet.Settings.DefaultCurrency,
		Balance:   wallet.GetTotalBalance(wallet.Settings.DefaultCurrency),
		Status:    string(wallet.Status),
		Metadata:  wallet.Metadata,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}, nil
}

// GetUserWallets retrieves all wallets for a user
func (ps *PaymentService) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*WalletResponse, error) {
	wallets, err := ps.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}

	responses := make([]*WalletResponse, len(wallets))
	for i, wallet := range wallets {
		var walletUserID uuid.UUID
		if wallet.UserID != nil {
			walletUserID = *wallet.UserID
		}

		responses[i] = &WalletResponse{
			ID:        wallet.ID,
			UserID:    walletUserID,
			Currency:  wallet.Settings.DefaultCurrency,
			Balance:   wallet.GetTotalBalance(wallet.Settings.DefaultCurrency),
			Status:    string(wallet.Status),
			Metadata:  wallet.Metadata,
			CreatedAt: wallet.CreatedAt,
			UpdatedAt: wallet.UpdatedAt,
		}
	}

	return responses, nil
}

// CreateTransaction creates a new transaction
func (ps *PaymentService) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*TransactionResponse, error) {
	// Get wallet to extract UserID
	wallet, err := ps.walletRepo.GetByID(ctx, req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	var userID uuid.UUID
	if wallet.UserID != nil {
		userID = *wallet.UserID
	}

	// Start database transaction
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	gateway := req.Gateway
	if gateway == "" {
		gateway = "internal"
	}

	transaction := &sharedModels.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		WalletID:    req.WalletID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Type:        sharedModels.TransactionType(req.Type),
		Status:      sharedModels.TransactionStatusPending,
		Description: req.Description,
		Gateway:     sharedModels.PaymentGateway(gateway),
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := ps.txRepo.Create(ctx, transaction); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update wallet balance
	if err := ps.walletRepo.UpdateBalance(ctx, req.WalletID, req.Amount, transaction.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Update transaction status to completed
	transaction.Status = sharedModels.TransactionStatusCompleted
	transaction.UpdatedAt = time.Now()

	if err := ps.txRepo.Update(ctx, transaction); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Publish transaction completed event
	eventData, _ := json.Marshal(map[string]interface{}{
		"amount":   transaction.Amount,
		"currency": transaction.Currency,
		"type":     transaction.Type,
	})

	event := &sharedModels.Event{
		ID:           uuid.New(),
		Type:         sharedModels.EventTypeTransactionCreated,
		Category:     sharedModels.EventCategoryDomain,
		Severity:     sharedModels.SeverityInfo,
		Subject:      "Transaction Completed",
		Description:  fmt.Sprintf("Transaction %s completed", transaction.ID),
		Data:         json.RawMessage(eventData),
		AggregateID:  transaction.ID.String(),
		OccurredAt:   time.Now(),
	}

	if err := ps.eventPublisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish transaction completed event: %v\n", err)
	}

	return &TransactionResponse{
		ID:          transaction.ID,
		UserID:      transaction.UserID,
		WalletID:    transaction.WalletID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Type:        string(transaction.Type),
		Status:      string(transaction.Status),
		Description: transaction.Description,
		Gateway:     string(transaction.Gateway),
		Metadata:    transaction.Metadata,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}, nil
}

// GetTransaction retrieves a transaction by ID
func (ps *PaymentService) GetTransaction(ctx context.Context, transactionID uuid.UUID) (*TransactionResponse, error) {
	transaction, err := ps.txRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &TransactionResponse{
		ID:          transaction.ID,
		UserID:      transaction.UserID,
		WalletID:    transaction.WalletID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Type:        string(transaction.Type),
		Status:      string(transaction.Status),
		Description: transaction.Description,
		Gateway:     string(transaction.Gateway),
		Metadata:    transaction.Metadata,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}, nil
}

// GetWalletTransactions retrieves transactions for a wallet
func (ps *PaymentService) GetWalletTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*TransactionResponse, error) {
	transactions, err := ps.txRepo.GetByWalletID(ctx, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet transactions: %w", err)
	}

	responses := make([]*TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = &TransactionResponse{
			ID:          tx.ID,
			UserID:      tx.UserID,
			WalletID:    tx.WalletID,
			Amount:      tx.Amount,
			Currency:    tx.Currency,
			Type:        string(tx.Type),
			Status:      string(tx.Status),
			Description: tx.Description,
			Gateway:     string(tx.Gateway),
			Metadata:    tx.Metadata,
			CreatedAt:   tx.CreatedAt,
			UpdatedAt:   tx.UpdatedAt,
		}
	}

	return responses, nil
}