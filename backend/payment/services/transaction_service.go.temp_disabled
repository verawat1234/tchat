package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/payment/models"
	"tchat.dev/payment/repositories"
	"tchat.dev/shared/events"
)

type TransactionService struct {
	transactionRepo repositories.TransactionRepository
	walletRepo      repositories.WalletRepository
	eventPublisher  events.EventPublisher
	db              *gorm.DB
}

func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	walletRepo repositories.WalletRepository,
	eventPublisher events.EventPublisher,
	db *gorm.DB,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		eventPublisher:  eventPublisher,
		db:              db,
	}
}

// Transaction processing requests
type ProcessTransactionRequest struct {
	FromWalletID uuid.UUID `json:"from_wallet_id" validate:"required"`
	ToWalletID   uuid.UUID `json:"to_wallet_id" validate:"required"`
	Amount       float64   `json:"amount" validate:"required,gt=0"`
	Currency     string    `json:"currency" validate:"required,len=3"`
	Type         string    `json:"type" validate:"required,oneof=transfer payment settlement"`
	Description  string    `json:"description" validate:"max=500"`
	Metadata     map[string]interface{} `json:"metadata"`
	ExternalID   string    `json:"external_id" validate:"max=100"`
}

type CreateTransactionRequest struct {
	WalletID    uuid.UUID `json:"wallet_id" validate:"required"`
	Amount      float64   `json:"amount" validate:"required"`
	Currency    string    `json:"currency" validate:"required,len=3"`
	Type        string    `json:"type" validate:"required,oneof=credit debit"`
	Status      string    `json:"status" validate:"required,oneof=pending processing completed failed"`
	Direction   string    `json:"direction" validate:"required,oneof=incoming outgoing"`
	Description string    `json:"description" validate:"max=500"`
	Metadata    map[string]interface{} `json:"metadata"`
	ExternalID  string    `json:"external_id" validate:"max=100"`
}

type UpdateTransactionStatusRequest struct {
	TransactionID uuid.UUID `json:"transaction_id" validate:"required"`
	Status        string    `json:"status" validate:"required,oneof=pending processing completed failed cancelled"`
	FailureReason string    `json:"failure_reason,omitempty" validate:"max=500"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty"`
}

type GetTransactionHistoryRequest struct {
	WalletID  uuid.UUID `json:"wallet_id" validate:"required"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Type      string    `json:"type,omitempty" validate:"omitempty,oneof=credit debit"`
	Status    string    `json:"status,omitempty" validate:"omitempty,oneof=pending processing completed failed cancelled"`
	Limit     int       `json:"limit" validate:"min=1,max=100"`
	Offset    int       `json:"offset" validate:"min=0"`
}

// ProcessTransaction handles distributed transaction processing
func (ts *TransactionService) ProcessTransaction(ctx context.Context, req *ProcessTransactionRequest) (*models.Transaction, error) {
	// Validate request
	if err := ts.validateProcessTransactionRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start database transaction
	tx := ts.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify both wallets exist and are active
	fromWallet, err := ts.walletRepo.GetByID(ctx, req.FromWalletID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get from wallet: %w", err)
	}

	toWallet, err := ts.walletRepo.GetByID(ctx, req.ToWalletID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get to wallet: %w", err)
	}

	// Validate transaction
	if err := ts.validateWalletTransaction(fromWallet, toWallet, req.Amount, req.Currency); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("transaction validation failed: %w", err)
	}

	// Check sufficient balance
	if fromWallet.AvailableBalance < req.Amount {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance: available %.2f, required %.2f", fromWallet.AvailableBalance, req.Amount)
	}

	// Create transaction record
	transaction := &models.Transaction{
		ID:          uuid.New(),
		WalletID:    req.FromWalletID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Type:        models.TransactionTypeDebit,
		Status:      models.TransactionStatusProcessing,
		Direction:   models.TransactionDirectionOutgoing,
		Description: req.Description,
		ExternalID:  req.ExternalID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Add metadata
	if req.Metadata != nil {
		metadataBytes, _ := json.Marshal(req.Metadata)
		transaction.Metadata = metadataBytes
	}

	// Save initial transaction
	if err := ts.transactionRepo.Create(ctx, transaction); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Process the transfer
	if err := ts.processWalletTransfer(ctx, fromWallet, toWallet, req.Amount, transaction.ID); err != nil {
		// Update transaction status to failed
		transaction.Status = models.TransactionStatusFailed
		ts.transactionRepo.Update(ctx, transaction)
		tx.Rollback()
		return nil, fmt.Errorf("transfer processing failed: %w", err)
	}

	// Update transaction status to completed
	transaction.Status = models.TransactionStatusCompleted
	transaction.ProcessedAt = &time.Time{}
	*transaction.ProcessedAt = time.Now()

	if err := ts.transactionRepo.Update(ctx, transaction); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Publish events
	ts.publishTransactionEvents(ctx, transaction, req.ToWalletID)

	return transaction, nil
}

// CreateTransaction creates a new transaction record
func (ts *TransactionService) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*models.Transaction, error) {
	// Validate request
	if err := ts.validateCreateTransactionRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify wallet exists
	wallet, err := ts.walletRepo.GetByID(ctx, req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	if wallet.Status != models.WalletStatusActive {
		return nil, fmt.Errorf("wallet is not active: %s", wallet.Status)
	}

	// Check for duplicate external ID
	if req.ExternalID != "" {
		existing, err := ts.transactionRepo.GetByExternalID(ctx, req.ExternalID)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("transaction with external ID already exists: %s", req.ExternalID)
		}
	}

	// Create transaction
	transaction := &models.Transaction{
		ID:          uuid.New(),
		WalletID:    req.WalletID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Type:        models.TransactionType(req.Type),
		Status:      models.TransactionStatus(req.Status),
		Direction:   models.TransactionDirection(req.Direction),
		Description: req.Description,
		ExternalID:  req.ExternalID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Add metadata
	if req.Metadata != nil {
		metadataBytes, _ := json.Marshal(req.Metadata)
		transaction.Metadata = metadataBytes
	}

	// Save transaction
	if err := ts.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Publish event
	ts.publishTransactionCreatedEvent(ctx, transaction)

	return transaction, nil
}

// UpdateTransactionStatus updates the status of a transaction
func (ts *TransactionService) UpdateTransactionStatus(ctx context.Context, req *UpdateTransactionStatusRequest) (*models.Transaction, error) {
	// Validate request
	if err := ts.validateUpdateTransactionStatusRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing transaction
	transaction, err := ts.transactionRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Validate status transition
	if err := ts.validateStatusTransition(transaction.Status, models.TransactionStatus(req.Status)); err != nil {
		return nil, fmt.Errorf("invalid status transition: %w", err)
	}

	// Update transaction
	transaction.Status = models.TransactionStatus(req.Status)
	transaction.FailureReason = req.FailureReason
	transaction.UpdatedAt = time.Now()

	if req.ProcessedAt != nil {
		transaction.ProcessedAt = req.ProcessedAt
	} else if req.Status == "completed" || req.Status == "failed" {
		now := time.Now()
		transaction.ProcessedAt = &now
	}

	// Save updated transaction
	if err := ts.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Publish event
	ts.publishTransactionUpdatedEvent(ctx, transaction)

	return transaction, nil
}

// GetTransactionHistory retrieves transaction history for a wallet
func (ts *TransactionService) GetTransactionHistory(ctx context.Context, req *GetTransactionHistoryRequest) ([]*models.Transaction, int64, error) {
	// Validate request
	if err := ts.validateGetTransactionHistoryRequest(req); err != nil {
		return nil, 0, fmt.Errorf("validation failed: %w", err)
	}

	// Verify wallet exists
	_, err := ts.walletRepo.GetByID(ctx, req.WalletID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Get transactions
	transactions, total, err := ts.transactionRepo.GetByWallet(ctx, req.WalletID, repositories.TransactionFilter{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Type:      req.Type,
		Status:    req.Status,
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, total, nil
}

// GetTransactionByID retrieves a transaction by ID
func (ts *TransactionService) GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error) {
	transaction, err := ts.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

// GetTransactionByExternalID retrieves a transaction by external ID
func (ts *TransactionService) GetTransactionByExternalID(ctx context.Context, externalID string) (*models.Transaction, error) {
	if externalID == "" {
		return nil, fmt.Errorf("external ID is required")
	}

	transaction, err := ts.transactionRepo.GetByExternalID(ctx, externalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

// processWalletTransfer handles the actual wallet balance transfer
func (ts *TransactionService) processWalletTransfer(ctx context.Context, fromWallet, toWallet *models.Wallet, amount float64, transactionID uuid.UUID) error {
	// Deduct from source wallet
	fromWallet.AvailableBalance -= amount
	fromWallet.UpdatedAt = time.Now()

	if err := ts.walletRepo.Update(ctx, fromWallet); err != nil {
		return fmt.Errorf("failed to update from wallet: %w", err)
	}

	// Add to destination wallet
	toWallet.AvailableBalance += amount
	toWallet.UpdatedAt = time.Now()

	if err := ts.walletRepo.Update(ctx, toWallet); err != nil {
		// Rollback from wallet update
		fromWallet.AvailableBalance += amount
		ts.walletRepo.Update(ctx, fromWallet)
		return fmt.Errorf("failed to update to wallet: %w", err)
	}

	// Create corresponding credit transaction for destination wallet
	creditTransaction := &models.Transaction{
		ID:          uuid.New(),
		WalletID:    toWallet.ID,
		Amount:      amount,
		Currency:    toWallet.Currency,
		Type:        models.TransactionTypeCredit,
		Status:      models.TransactionStatusCompleted,
		Direction:   models.TransactionDirectionIncoming,
		Description: fmt.Sprintf("Transfer from wallet %s", fromWallet.ID),
		ExternalID:  fmt.Sprintf("transfer_%s", transactionID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	now := time.Now()
	creditTransaction.ProcessedAt = &now

	if err := ts.transactionRepo.Create(ctx, creditTransaction); err != nil {
		return fmt.Errorf("failed to create credit transaction: %w", err)
	}

	return nil
}

// Validation functions
func (ts *TransactionService) validateProcessTransactionRequest(req *ProcessTransactionRequest) error {
	if req.FromWalletID == uuid.Nil {
		return fmt.Errorf("from wallet ID is required")
	}
	if req.ToWalletID == uuid.Nil {
		return fmt.Errorf("to wallet ID is required")
	}
	if req.FromWalletID == req.ToWalletID {
		return fmt.Errorf("from and to wallet cannot be the same")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters")
	}
	return nil
}

func (ts *TransactionService) validateCreateTransactionRequest(req *CreateTransactionRequest) error {
	if req.WalletID == uuid.Nil {
		return fmt.Errorf("wallet ID is required")
	}
	if req.Amount == 0 {
		return fmt.Errorf("amount cannot be zero")
	}
	if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters")
	}
	return nil
}

func (ts *TransactionService) validateUpdateTransactionStatusRequest(req *UpdateTransactionStatusRequest) error {
	if req.TransactionID == uuid.Nil {
		return fmt.Errorf("transaction ID is required")
	}
	validStatuses := []string{"pending", "processing", "completed", "failed", "cancelled"}
	for _, status := range validStatuses {
		if req.Status == status {
			return nil
		}
	}
	return fmt.Errorf("invalid status: %s", req.Status)
}

func (ts *TransactionService) validateGetTransactionHistoryRequest(req *GetTransactionHistoryRequest) error {
	if req.WalletID == uuid.Nil {
		return fmt.Errorf("wallet ID is required")
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20 // Default limit
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	return nil
}

func (ts *TransactionService) validateWalletTransaction(fromWallet, toWallet *models.Wallet, amount float64, currency string) error {
	if fromWallet.Status != models.WalletStatusActive {
		return fmt.Errorf("from wallet is not active: %s", fromWallet.Status)
	}
	if toWallet.Status != models.WalletStatusActive {
		return fmt.Errorf("to wallet is not active: %s", toWallet.Status)
	}
	if fromWallet.Currency != currency {
		return fmt.Errorf("currency mismatch: from wallet %s, transaction %s", fromWallet.Currency, currency)
	}
	if toWallet.Currency != currency {
		return fmt.Errorf("currency mismatch: to wallet %s, transaction %s", toWallet.Currency, currency)
	}
	return nil
}

func (ts *TransactionService) validateStatusTransition(currentStatus, newStatus models.TransactionStatus) error {
	// Define valid transitions
	validTransitions := map[models.TransactionStatus][]models.TransactionStatus{
		models.TransactionStatusPending:    {models.TransactionStatusProcessing, models.TransactionStatusCancelled},
		models.TransactionStatusProcessing: {models.TransactionStatusCompleted, models.TransactionStatusFailed},
		models.TransactionStatusCompleted:  {}, // No transitions from completed
		models.TransactionStatusFailed:     {}, // No transitions from failed
		models.TransactionStatusCancelled:  {}, // No transitions from cancelled
	}

	allowedTransitions, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("unknown current status: %s", currentStatus)
	}

	for _, allowed := range allowedTransitions {
		if newStatus == allowed {
			return nil
		}
	}

	return fmt.Errorf("invalid transition from %s to %s", currentStatus, newStatus)
}

// Event publishing functions
func (ts *TransactionService) publishTransactionEvents(ctx context.Context, transaction *models.Transaction, toWalletID uuid.UUID) {
	// Publish transaction completed event
	ts.eventPublisher.Publish(ctx, events.Event{
		Type: "transaction.completed",
		Payload: map[string]interface{}{
			"transaction_id": transaction.ID,
			"from_wallet_id": transaction.WalletID,
			"to_wallet_id":   toWalletID,
			"amount":         transaction.Amount,
			"currency":       transaction.Currency,
			"processed_at":   transaction.ProcessedAt,
		},
	})

	// Publish wallet balance updated events
	ts.eventPublisher.Publish(ctx, events.Event{
		Type: "wallet.balance_updated",
		Payload: map[string]interface{}{
			"wallet_id":      transaction.WalletID,
			"transaction_id": transaction.ID,
			"amount":         -transaction.Amount,
			"type":           "debit",
		},
	})

	ts.eventPublisher.Publish(ctx, events.Event{
		Type: "wallet.balance_updated",
		Payload: map[string]interface{}{
			"wallet_id":      toWalletID,
			"transaction_id": transaction.ID,
			"amount":         transaction.Amount,
			"type":           "credit",
		},
	})
}

func (ts *TransactionService) publishTransactionCreatedEvent(ctx context.Context, transaction *models.Transaction) {
	ts.eventPublisher.Publish(ctx, events.Event{
		Type: "transaction.created",
		Payload: map[string]interface{}{
			"transaction_id": transaction.ID,
			"wallet_id":      transaction.WalletID,
			"amount":         transaction.Amount,
			"currency":       transaction.Currency,
			"type":           transaction.Type,
			"status":         transaction.Status,
			"created_at":     transaction.CreatedAt,
		},
	})
}

func (ts *TransactionService) publishTransactionUpdatedEvent(ctx context.Context, transaction *models.Transaction) {
	ts.eventPublisher.Publish(ctx, events.Event{
		Type: "transaction.status_updated",
		Payload: map[string]interface{}{
			"transaction_id": transaction.ID,
			"wallet_id":      transaction.WalletID,
			"status":         transaction.Status,
			"failure_reason": transaction.FailureReason,
			"updated_at":     transaction.UpdatedAt,
		},
	})
}