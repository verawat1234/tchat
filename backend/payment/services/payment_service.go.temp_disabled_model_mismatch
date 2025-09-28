package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"tchat.dev/payment/models"
)

// PaymentService provides payment and wallet functionality
type PaymentService struct {
	walletRepo      WalletRepository
	transactionRepo TransactionRepository
	cache           CacheService
	events          EventService
	processor       PaymentProcessor
	config          *PaymentConfig
}

// PaymentConfig holds payment service configuration
type PaymentConfig struct {
	DefaultCurrency        models.Currency
	MaxTransactionAmount   int64
	MinTransactionAmount   int64
	FeePercentage         float64
	FixedFee              int64
	MaxDailyTransactions  int
	MaxMonthlyTransactions int
	EnableFraudDetection  bool
	TransactionTimeout    time.Duration
	ProcessingRetries     int
}

// DefaultPaymentConfig returns default payment configuration
func DefaultPaymentConfig() *PaymentConfig {
	return &PaymentConfig{
		DefaultCurrency:        models.CurrencyTHB,
		MaxTransactionAmount:   100000000, // 1M THB in cents
		MinTransactionAmount:   100,       // 1 THB in cents
		FeePercentage:         0.025,     // 2.5%
		FixedFee:              1000,      // 10 THB in cents
		MaxDailyTransactions:  100,
		MaxMonthlyTransactions: 1000,
		EnableFraudDetection:  true,
		TransactionTimeout:    30 * time.Minute,
		ProcessingRetries:     3,
	}
}

// WalletRepository interface for wallet data access
type WalletRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, currency models.Currency) (*models.Wallet, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64, operation string) error
	LockForUpdate(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
}

// TransactionRepository interface for transaction data access
type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	GetByReference(ctx context.Context, reference string) (*models.Transaction, error)
	GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*models.Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Transaction, error)
	Update(ctx context.Context, transaction *models.Transaction) error
	GetDailyTotal(ctx context.Context, walletID uuid.UUID, date time.Time) (int64, error)
	GetMonthlyTotal(ctx context.Context, walletID uuid.UUID, year int, month time.Month) (int64, error)
	GetPendingTransactions(ctx context.Context) ([]*models.Transaction, error)
	UpdateStatus(ctx context.Context, transactionID uuid.UUID, status models.TransactionStatus) error
}

// CacheService interface for caching operations
type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
	Lock(ctx context.Context, key string, expiry time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
}

// EventService interface for event publishing
type EventService interface {
	PublishTransaction(ctx context.Context, event *TransactionEvent) error
	PublishWallet(ctx context.Context, event *WalletEvent) error
}

// PaymentProcessor interface for external payment processing
type PaymentProcessor interface {
	ProcessDeposit(ctx context.Context, req *DepositRequest) (*ProcessorResponse, error)
	ProcessWithdrawal(ctx context.Context, req *WithdrawalRequest) (*ProcessorResponse, error)
	ProcessTransfer(ctx context.Context, req *TransferRequest) (*ProcessorResponse, error)
	GetStatus(ctx context.Context, externalID string) (*ProcessorStatus, error)
}

// Events
type TransactionEvent struct {
	Type         string                `json:"type"`
	TransactionID uuid.UUID            `json:"transaction_id"`
	WalletID     uuid.UUID             `json:"wallet_id"`
	UserID       uuid.UUID             `json:"user_id"`
	Transaction  *models.Transaction   `json:"transaction,omitempty"`
	Timestamp    time.Time             `json:"timestamp"`
}

type WalletEvent struct {
	Type      string         `json:"type"`
	WalletID  uuid.UUID      `json:"wallet_id"`
	UserID    uuid.UUID      `json:"user_id"`
	Wallet    *models.Wallet `json:"wallet,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// Request/Response types
type CreateWalletRequest struct {
	UserID   uuid.UUID       `json:"user_id"`
	Currency models.Currency `json:"currency"`
	Name     *string         `json:"name,omitempty"`
}

type DepositRequest struct {
	WalletID     uuid.UUID       `json:"wallet_id"`
	Amount       int64           `json:"amount"`
	Currency     models.Currency `json:"currency"`
	PaymentMethod string         `json:"payment_method"`
	ExternalID   *string         `json:"external_id,omitempty"`
	Description  *string         `json:"description,omitempty"`
}

type WithdrawalRequest struct {
	WalletID      uuid.UUID       `json:"wallet_id"`
	Amount        int64           `json:"amount"`
	Currency      models.Currency `json:"currency"`
	Destination   string          `json:"destination"`
	ExternalID    *string         `json:"external_id,omitempty"`
	Description   *string         `json:"description,omitempty"`
}

type TransferRequest struct {
	FromWalletID uuid.UUID       `json:"from_wallet_id"`
	ToWalletID   uuid.UUID       `json:"to_wallet_id"`
	Amount       int64           `json:"amount"`
	Currency     models.Currency `json:"currency"`
	Description  *string         `json:"description,omitempty"`
}

type ProcessorResponse struct {
	ExternalID string                 `json:"external_id"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

type ProcessorStatus struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	walletRepo WalletRepository,
	transactionRepo TransactionRepository,
	cache CacheService,
	events EventService,
	processor PaymentProcessor,
	config *PaymentConfig,
) *PaymentService {
	if config == nil {
		config = DefaultPaymentConfig()
	}

	return &PaymentService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		cache:           cache,
		events:          events,
		processor:       processor,
		config:          config,
	}
}

// CreateWallet creates a new wallet for a user
func (p *PaymentService) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*models.Wallet, error) {
	// Validate request
	if err := p.validateCreateWalletRequest(req); err != nil {
		return nil, fmt.Errorf("invalid create wallet request: %v", err)
	}

	// Check if wallet already exists for this user and currency
	existingWallet, err := p.walletRepo.GetByUserID(ctx, req.UserID, req.Currency)
	if err == nil && existingWallet != nil {
		return nil, fmt.Errorf("wallet already exists for user %s and currency %s", req.UserID, req.Currency)
	}

	// Create wallet
	walletManager := models.NewWalletManager()
	wallet, err := walletManager.CreateWallet(&models.WalletCreateRequest{
		UserID:   req.UserID,
		Currency: req.Currency,
		Name:     req.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %v", err)
	}

	// Save wallet
	if err := p.walletRepo.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to save wallet: %v", err)
	}

	// Publish wallet creation event
	event := &WalletEvent{
		Type:      "wallet_created",
		WalletID:  wallet.ID,
		UserID:    req.UserID,
		Wallet:    wallet,
		Timestamp: time.Now().UTC(),
	}
	p.events.PublishWallet(ctx, event)

	return wallet, nil
}

// Deposit processes a deposit transaction
func (p *PaymentService) Deposit(ctx context.Context, req *DepositRequest) (*models.Transaction, error) {
	// Validate request
	if err := p.validateDepositRequest(req); err != nil {
		return nil, fmt.Errorf("invalid deposit request: %v", err)
	}

	// Get wallet
	wallet, err := p.walletRepo.GetByID(ctx, req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %v", err)
	}

	// Validate currency match
	if wallet.Currency != req.Currency {
		return nil, fmt.Errorf("currency mismatch: wallet currency %s, request currency %s", wallet.Currency, req.Currency)
	}

	// Calculate fees
	feeAmount := p.calculateFee(req.Amount)

	// Create transaction
	transactionManager := models.NewTransactionManager()
	transaction, err := transactionManager.CreateTransaction(&models.TransactionCreateRequest{
		WalletID:   req.WalletID,
		Type:       models.TransactionTypeDeposit,
		Currency:   req.Currency,
		Amount:     req.Amount,
		FeeAmount:  feeAmount,
		ExternalID: req.ExternalID,
		Description: req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Validate transaction limits
	if err := transactionManager.ValidateTransactionLimits(transaction, wallet); err != nil {
		return nil, fmt.Errorf("transaction validation failed: %v", err)
	}

	// Save transaction
	if err := p.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %v", err)
	}

	// Process with external payment processor
	if err := p.processDeposit(ctx, transaction, req); err != nil {
		// Mark transaction as failed
		transaction.Fail(err.Error())
		p.transactionRepo.Update(ctx, transaction)
		return nil, fmt.Errorf("deposit processing failed: %v", err)
	}

	return transaction, nil
}

// Withdraw processes a withdrawal transaction
func (p *PaymentService) Withdraw(ctx context.Context, req *WithdrawalRequest) (*models.Transaction, error) {
	// Validate request
	if err := p.validateWithdrawalRequest(req); err != nil {
		return nil, fmt.Errorf("invalid withdrawal request: %v", err)
	}

	// Lock wallet for update
	lockKey := fmt.Sprintf("wallet_lock:%s", req.WalletID)
	locked, err := p.cache.Lock(ctx, lockKey, 30*time.Second)
	if !locked || err != nil {
		return nil, fmt.Errorf("failed to acquire wallet lock")
	}
	defer p.cache.Unlock(ctx, lockKey)

	// Get wallet
	wallet, err := p.walletRepo.LockForUpdate(ctx, req.WalletID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %v", err)
	}

	// Validate currency match
	if wallet.Currency != req.Currency {
		return nil, fmt.Errorf("currency mismatch: wallet currency %s, request currency %s", wallet.Currency, req.Currency)
	}

	// Calculate fees
	feeAmount := p.calculateFee(req.Amount)
	totalAmount := req.Amount + feeAmount

	// Check available balance
	availableBalance := wallet.Balance - wallet.FrozenBalance
	if totalAmount > availableBalance {
		return nil, fmt.Errorf("insufficient balance: required %d, available %d", totalAmount, availableBalance)
	}

	// Check daily limit
	today := time.Now().UTC().Truncate(24 * time.Hour)
	dailyTotal, err := p.transactionRepo.GetDailyTotal(ctx, req.WalletID, today)
	if err != nil {
		return nil, fmt.Errorf("failed to check daily limits: %v", err)
	}

	if dailyTotal+totalAmount > wallet.DailyLimit {
		return nil, fmt.Errorf("daily limit exceeded: current %d, limit %d, requested %d",
			dailyTotal, wallet.DailyLimit, totalAmount)
	}

	// Create transaction
	transactionManager := models.NewTransactionManager()
	transaction, err := transactionManager.CreateTransaction(&models.TransactionCreateRequest{
		WalletID:   req.WalletID,
		Type:       models.TransactionTypeWithdrawal,
		Currency:   req.Currency,
		Amount:     req.Amount,
		FeeAmount:  feeAmount,
		ExternalID: req.ExternalID,
		Description: req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Freeze balance for withdrawal
	wallet.FrozenBalance += totalAmount
	if err := p.walletRepo.Update(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to freeze balance: %v", err)
	}

	// Save transaction
	if err := p.transactionRepo.Create(ctx, transaction); err != nil {
		// Rollback frozen balance
		wallet.FrozenBalance -= totalAmount
		p.walletRepo.Update(ctx, wallet)
		return nil, fmt.Errorf("failed to save transaction: %v", err)
	}

	// Process withdrawal
	if err := p.processWithdrawal(ctx, transaction, req); err != nil {
		// Mark transaction as failed and unfreeze balance
		transaction.Fail(err.Error())
		p.transactionRepo.Update(ctx, transaction)

		wallet.FrozenBalance -= totalAmount
		p.walletRepo.Update(ctx, wallet)

		return nil, fmt.Errorf("withdrawal processing failed: %v", err)
	}

	return transaction, nil
}

// Transfer processes a transfer between wallets
func (p *PaymentService) Transfer(ctx context.Context, req *TransferRequest) (*models.Transaction, error) {
	// Validate request
	if err := p.validateTransferRequest(req); err != nil {
		return nil, fmt.Errorf("invalid transfer request: %v", err)
	}

	// Prevent self-transfer
	if req.FromWalletID == req.ToWalletID {
		return nil, fmt.Errorf("cannot transfer to the same wallet")
	}

	// Lock both wallets in a consistent order to prevent deadlock
	walletIDs := []uuid.UUID{req.FromWalletID, req.ToWalletID}
	if req.FromWalletID.String() > req.ToWalletID.String() {
		walletIDs[0], walletIDs[1] = walletIDs[1], walletIDs[0]
	}

	lockKey1 := fmt.Sprintf("wallet_lock:%s", walletIDs[0])
	lockKey2 := fmt.Sprintf("wallet_lock:%s", walletIDs[1])

	locked1, err1 := p.cache.Lock(ctx, lockKey1, 30*time.Second)
	locked2, err2 := p.cache.Lock(ctx, lockKey2, 30*time.Second)

	if !locked1 || !locked2 || err1 != nil || err2 != nil {
		p.cache.Unlock(ctx, lockKey1)
		p.cache.Unlock(ctx, lockKey2)
		return nil, fmt.Errorf("failed to acquire wallet locks")
	}
	defer func() {
		p.cache.Unlock(ctx, lockKey1)
		p.cache.Unlock(ctx, lockKey2)
	}()

	// Get both wallets
	fromWallet, err := p.walletRepo.LockForUpdate(ctx, req.FromWalletID)
	if err != nil {
		return nil, fmt.Errorf("source wallet not found: %v", err)
	}

	toWallet, err := p.walletRepo.LockForUpdate(ctx, req.ToWalletID)
	if err != nil {
		return nil, fmt.Errorf("destination wallet not found: %v", err)
	}

	// Validate currency match
	if fromWallet.Currency != req.Currency || toWallet.Currency != req.Currency {
		return nil, fmt.Errorf("currency mismatch")
	}

	// Calculate fees
	feeAmount := p.calculateFee(req.Amount)
	totalAmount := req.Amount + feeAmount

	// Check available balance
	availableBalance := fromWallet.Balance - fromWallet.FrozenBalance
	if totalAmount > availableBalance {
		return nil, fmt.Errorf("insufficient balance: required %d, available %d", totalAmount, availableBalance)
	}

	// Create transaction
	transactionManager := models.NewTransactionManager()
	transaction, err := transactionManager.CreateTransaction(&models.TransactionCreateRequest{
		WalletID:       req.FromWalletID,
		CounterpartyID: &req.ToWalletID,
		Type:           models.TransactionTypeTransfer,
		Currency:       req.Currency,
		Amount:         req.Amount,
		FeeAmount:      feeAmount,
		Description:    req.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Save transaction
	if err := p.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %v", err)
	}

	// Start processing
	if err := transaction.StartProcessing(); err != nil {
		return nil, fmt.Errorf("failed to start processing: %v", err)
	}
	p.transactionRepo.Update(ctx, transaction)

	// Execute transfer
	if err := p.executeTransfer(ctx, transaction, fromWallet, toWallet); err != nil {
		// Mark transaction as failed
		transaction.Fail(err.Error())
		p.transactionRepo.Update(ctx, transaction)
		return nil, fmt.Errorf("transfer execution failed: %v", err)
	}

	// Complete transaction
	transaction.Complete()
	p.transactionRepo.Update(ctx, transaction)

	// Publish transaction event
	event := &TransactionEvent{
		Type:          "transfer_completed",
		TransactionID: transaction.ID,
		WalletID:      req.FromWalletID,
		UserID:        fromWallet.UserID,
		Transaction:   transaction,
		Timestamp:     time.Now().UTC(),
	}
	p.events.PublishTransaction(ctx, event)

	return transaction, nil
}

// GetWallet retrieves a wallet by ID
func (p *PaymentService) GetWallet(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	return p.walletRepo.GetByID(ctx, walletID)
}

// GetUserWallets retrieves all wallets for a user
func (p *PaymentService) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*models.Wallet, error) {
	return p.walletRepo.GetAllByUserID(ctx, userID)
}

// GetTransaction retrieves a transaction by ID
func (p *PaymentService) GetTransaction(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error) {
	return p.transactionRepo.GetByID(ctx, transactionID)
}

// GetWalletTransactions retrieves transactions for a wallet
func (p *PaymentService) GetWalletTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*models.Transaction, error) {
	return p.transactionRepo.GetByWalletID(ctx, walletID, limit, offset)
}

// ProcessPendingTransactions processes all pending transactions
func (p *PaymentService) ProcessPendingTransactions(ctx context.Context) error {
	transactions, err := p.transactionRepo.GetPendingTransactions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending transactions: %v", err)
	}

	for _, transaction := range transactions {
		// Check if transaction is expired
		if transaction.IsExpired() {
			transaction.Expire()
			p.transactionRepo.Update(ctx, transaction)
			continue
		}

		// Process based on transaction type
		switch transaction.Type {
		case models.TransactionTypeDeposit:
			p.processPendingDeposit(ctx, transaction)
		case models.TransactionTypeWithdrawal:
			p.processPendingWithdrawal(ctx, transaction)
		}
	}

	return nil
}

// Helper methods

func (p *PaymentService) validateCreateWalletRequest(req *CreateWalletRequest) error {
	if req.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if !req.Currency.IsValid() {
		return fmt.Errorf("invalid currency: %s", req.Currency)
	}

	return nil
}

func (p *PaymentService) validateDepositRequest(req *DepositRequest) error {
	if req.WalletID == uuid.Nil {
		return errors.New("wallet_id is required")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if req.Amount < p.config.MinTransactionAmount {
		return fmt.Errorf("amount below minimum: %d", p.config.MinTransactionAmount)
	}

	if req.Amount > p.config.MaxTransactionAmount {
		return fmt.Errorf("amount exceeds maximum: %d", p.config.MaxTransactionAmount)
	}

	if !req.Currency.IsValid() {
		return fmt.Errorf("invalid currency: %s", req.Currency)
	}

	return nil
}

func (p *PaymentService) validateWithdrawalRequest(req *WithdrawalRequest) error {
	if req.WalletID == uuid.Nil {
		return errors.New("wallet_id is required")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if req.Amount < p.config.MinTransactionAmount {
		return fmt.Errorf("amount below minimum: %d", p.config.MinTransactionAmount)
	}

	if req.Amount > p.config.MaxTransactionAmount {
		return fmt.Errorf("amount exceeds maximum: %d", p.config.MaxTransactionAmount)
	}

	if !req.Currency.IsValid() {
		return fmt.Errorf("invalid currency: %s", req.Currency)
	}

	if strings.TrimSpace(req.Destination) == "" {
		return errors.New("destination is required")
	}

	return nil
}

func (p *PaymentService) validateTransferRequest(req *TransferRequest) error {
	if req.FromWalletID == uuid.Nil {
		return errors.New("from_wallet_id is required")
	}

	if req.ToWalletID == uuid.Nil {
		return errors.New("to_wallet_id is required")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if req.Amount < p.config.MinTransactionAmount {
		return fmt.Errorf("amount below minimum: %d", p.config.MinTransactionAmount)
	}

	if req.Amount > p.config.MaxTransactionAmount {
		return fmt.Errorf("amount exceeds maximum: %d", p.config.MaxTransactionAmount)
	}

	if !req.Currency.IsValid() {
		return fmt.Errorf("invalid currency: %s", req.Currency)
	}

	return nil
}

func (p *PaymentService) calculateFee(amount int64) int64 {
	percentageFee := int64(float64(amount) * p.config.FeePercentage)
	totalFee := percentageFee + p.config.FixedFee

	// Ensure fee doesn't exceed 50% of amount
	maxFee := amount / 2
	if totalFee > maxFee {
		totalFee = maxFee
	}

	return totalFee
}

func (p *PaymentService) processDeposit(ctx context.Context, transaction *models.Transaction, req *DepositRequest) error {
	// Start processing
	if err := transaction.StartProcessing(); err != nil {
		return err
	}
	p.transactionRepo.Update(ctx, transaction)

	// Process with external processor
	processorReq := &DepositRequest{
		WalletID:      req.WalletID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		ExternalID:    &transaction.ID.String(),
		Description:   req.Description,
	}

	response, err := p.processor.ProcessDeposit(ctx, processorReq)
	if err != nil {
		return err
	}

	// Update transaction with processor response
	transaction.ExternalID = &response.ExternalID
	processorResponseStr := fmt.Sprintf("Status: %s, Message: %s", response.Status, response.Message)
	transaction.ProcessorResponse = &processorResponseStr

	if response.Status == "completed" {
		// Complete the deposit
		return p.completeDeposit(ctx, transaction)
	}

	return nil // Will be processed later by background job
}

func (p *PaymentService) processWithdrawal(ctx context.Context, transaction *models.Transaction, req *WithdrawalRequest) error {
	// Start processing
	if err := transaction.StartProcessing(); err != nil {
		return err
	}
	p.transactionRepo.Update(ctx, transaction)

	// Process with external processor
	processorReq := &WithdrawalRequest{
		WalletID:    req.WalletID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Destination: req.Destination,
		ExternalID:  &transaction.ID.String(),
		Description: req.Description,
	}

	response, err := p.processor.ProcessWithdrawal(ctx, processorReq)
	if err != nil {
		return err
	}

	// Update transaction with processor response
	transaction.ExternalID = &response.ExternalID
	processorResponseStr := fmt.Sprintf("Status: %s, Message: %s", response.Status, response.Message)
	transaction.ProcessorResponse = &processorResponseStr

	if response.Status == "completed" {
		// Complete the withdrawal
		return p.completeWithdrawal(ctx, transaction)
	}

	return nil // Will be processed later by background job
}

func (p *PaymentService) executeTransfer(ctx context.Context, transaction *models.Transaction, fromWallet, toWallet *models.Wallet) error {
	totalAmount := transaction.Amount + transaction.FeeAmount

	// Deduct from source wallet
	fromWallet.Balance -= totalAmount
	if err := p.walletRepo.Update(ctx, fromWallet); err != nil {
		return fmt.Errorf("failed to update source wallet: %v", err)
	}

	// Add to destination wallet (excluding fee)
	toWallet.Balance += transaction.Amount
	if err := p.walletRepo.Update(ctx, toWallet); err != nil {
		// Rollback source wallet
		fromWallet.Balance += totalAmount
		p.walletRepo.Update(ctx, fromWallet)
		return fmt.Errorf("failed to update destination wallet: %v", err)
	}

	return nil
}

func (p *PaymentService) completeDeposit(ctx context.Context, transaction *models.Transaction) error {
	// Get wallet
	wallet, err := p.walletRepo.LockForUpdate(ctx, transaction.WalletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %v", err)
	}

	// Add balance (excluding fee)
	wallet.Balance += transaction.NetAmount
	if err := p.walletRepo.Update(ctx, wallet); err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Complete transaction
	if err := transaction.Complete(); err != nil {
		return fmt.Errorf("failed to complete transaction: %v", err)
	}

	return p.transactionRepo.Update(ctx, transaction)
}

func (p *PaymentService) completeWithdrawal(ctx context.Context, transaction *models.Transaction) error {
	// Get wallet
	wallet, err := p.walletRepo.LockForUpdate(ctx, transaction.WalletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %v", err)
	}

	totalAmount := transaction.Amount + transaction.FeeAmount

	// Unfreeze and deduct balance
	wallet.FrozenBalance -= totalAmount
	wallet.Balance -= totalAmount
	if err := p.walletRepo.Update(ctx, wallet); err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Complete transaction
	if err := transaction.Complete(); err != nil {
		return fmt.Errorf("failed to complete transaction: %v", err)
	}

	return p.transactionRepo.Update(ctx, transaction)
}

func (p *PaymentService) processPendingDeposit(ctx context.Context, transaction *models.Transaction) {
	if transaction.ExternalID == nil {
		return
	}

	// Check status with processor
	status, err := p.processor.GetStatus(ctx, *transaction.ExternalID)
	if err != nil {
		return
	}

	switch status.Status {
	case "completed":
		p.completeDeposit(ctx, transaction)
	case "failed":
		transaction.Fail(status.Message)
		p.transactionRepo.Update(ctx, transaction)
	}
}

func (p *PaymentService) processPendingWithdrawal(ctx context.Context, transaction *models.Transaction) {
	if transaction.ExternalID == nil {
		return
	}

	// Check status with processor
	status, err := p.processor.GetStatus(ctx, *transaction.ExternalID)
	if err != nil {
		return
	}

	switch status.Status {
	case "completed":
		p.completeWithdrawal(ctx, transaction)
	case "failed":
		// Unfreeze balance and mark as failed
		wallet, err := p.walletRepo.LockForUpdate(ctx, transaction.WalletID)
		if err == nil {
			totalAmount := transaction.Amount + transaction.FeeAmount
			wallet.FrozenBalance -= totalAmount
			p.walletRepo.Update(ctx, wallet)
		}

		transaction.Fail(status.Message)
		p.transactionRepo.Update(ctx, transaction)
	}
}