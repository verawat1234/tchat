package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/payment/models"
	sharedModels "tchat.dev/shared/models"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Wallet, error)
	GetByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency models.Currency) (*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, amount float64, transactionID uuid.UUID) error
	FreezeAmount(ctx context.Context, walletID uuid.UUID, amount float64, reason string) error
	UnfreezeAmount(ctx context.Context, walletID uuid.UUID, amount float64) error
	GetWalletHistory(ctx context.Context, walletID uuid.UUID, limit int) ([]*WalletBalanceHistory, error)
	GetWalletStats(ctx context.Context, userID uuid.UUID) (*WalletStats, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event *sharedModels.Event) error
}

type ComplianceChecker interface {
	CheckTransactionLimits(ctx context.Context, userID uuid.UUID, amount float64, currency models.Currency, transactionType string) error
	CheckAMLCompliance(ctx context.Context, userID uuid.UUID, amount float64, currency models.Currency) error
	ReportSuspiciousActivity(ctx context.Context, userID uuid.UUID, activity string, metadata map[string]interface{}) error
}

type ExchangeRateService interface {
	GetExchangeRate(ctx context.Context, fromCurrency, toCurrency models.Currency) (float64, error)
	ConvertAmount(ctx context.Context, amount float64, fromCurrency, toCurrency models.Currency) (float64, error)
}

type WalletBalanceHistory struct {
	ID            uuid.UUID `json:"id"`
	WalletID      uuid.UUID `json:"wallet_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	PreviousBalance float64 `json:"previous_balance"`
	NewBalance    float64   `json:"new_balance"`
	AmountChanged float64   `json:"amount_changed"`
	Timestamp     time.Time `json:"timestamp"`
	Reason        string    `json:"reason"`
}

type WalletStats struct {
	TotalWallets      int                        `json:"total_wallets"`
	TotalBalance      map[models.Currency]float64 `json:"total_balance"`
	FrozenBalance     map[models.Currency]float64 `json:"frozen_balance"`
	AvailableBalance  map[models.Currency]float64 `json:"available_balance"`
	TransactionCount  int64                      `json:"transaction_count"`
	LastTransactionAt *time.Time                 `json:"last_transaction_at"`
}

type WalletService struct {
	walletRepo       WalletRepository
	eventPublisher   EventPublisher
	complianceChecker ComplianceChecker
	exchangeRateService ExchangeRateService
	db               *gorm.DB
}

func NewWalletService(
	walletRepo WalletRepository,
	eventPublisher EventPublisher,
	complianceChecker ComplianceChecker,
	exchangeRateService ExchangeRateService,
	db *gorm.DB,
) *WalletService {
	return &WalletService{
		walletRepo:          walletRepo,
		eventPublisher:      eventPublisher,
		complianceChecker:   complianceChecker,
		exchangeRateService: exchangeRateService,
		db:                  db,
	}
}

func (ws *WalletService) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*models.Wallet, error) {
	// Validate request
	if err := ws.validateCreateWalletRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if wallet already exists for this user and currency
	existingWallet, err := ws.walletRepo.GetByUserIDAndCurrency(ctx, req.UserID, req.Currency)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing wallet: %w", err)
	}
	if existingWallet != nil {
		return nil, fmt.Errorf("wallet already exists for user %s and currency %s", req.UserID, req.Currency)
	}

	// Create wallet
	wallet := &models.Wallet{
		ID:                uuid.New(),
		UserID:            req.UserID,
		Currency:          req.Currency,
		Balance:           0.0,
		FrozenBalance:     0.0,
		AvailableBalance:  0.0,
		Status:            models.WalletStatusActive,
		Limits: models.WalletLimits{
			DailySpendLimit:    req.DailySpendLimit,
			MonthlySpendLimit:  req.MonthlySpendLimit,
			DailyTopupLimit:    req.DailyTopupLimit,
			MonthlyTopupLimit:  req.MonthlyTopupLimit,
			MaxBalance:         req.MaxBalance,
		},
		Metadata:          req.Metadata,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Set default limits if not provided
	wallet.SetDefaultLimits()

	// Validate wallet
	if err := wallet.Validate(); err != nil {
		return nil, fmt.Errorf("wallet validation failed: %w", err)
	}

	// Save wallet
	if err := ws.walletRepo.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Publish wallet created event
	if err := ws.publishWalletEvent(ctx, "wallet.created", wallet.ID, wallet.UserID, map[string]interface{}{
		"wallet_id": wallet.ID,
		"currency":  wallet.Currency,
		"limits":    wallet.Limits,
	}); err != nil {
		fmt.Printf("Failed to publish wallet created event: %v\n", err)
	}

	return wallet, nil
}

func (ws *WalletService) GetWalletByID(ctx context.Context, walletID uuid.UUID, userID uuid.UUID) (*models.Wallet, error) {
	if walletID == uuid.Nil {
		return nil, fmt.Errorf("wallet ID is required")
	}

	wallet, err := ws.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Check if user owns this wallet
	if wallet.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return wallet, nil
}

func (ws *WalletService) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*models.Wallet, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	wallets, err := ws.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}

	return wallets, nil
}

func (ws *WalletService) GetWalletByCurrency(ctx context.Context, userID uuid.UUID, currency models.Currency) (*models.Wallet, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	wallet, err := ws.walletRepo.GetByUserIDAndCurrency(ctx, userID, currency)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet not found for currency %s", currency)
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return wallet, nil
}

func (ws *WalletService) GetWalletBalance(ctx context.Context, walletID uuid.UUID, userID uuid.UUID) (*WalletBalanceResponse, error) {
	wallet, err := ws.GetWalletByID(ctx, walletID, userID)
	if err != nil {
		return nil, err
	}

	return &WalletBalanceResponse{
		WalletID:         wallet.ID,
		Currency:         wallet.Currency,
		Balance:          wallet.Balance,
		AvailableBalance: wallet.AvailableBalance,
		FrozenBalance:    wallet.FrozenBalance,
		FormattedBalance: wallet.FormatBalance(),
		LastUpdated:      wallet.UpdatedAt,
	}, nil
}

func (ws *WalletService) UpdateWalletLimits(ctx context.Context, walletID uuid.UUID, userID uuid.UUID, req *UpdateWalletLimitsRequest) (*models.Wallet, error) {
	wallet, err := ws.GetWalletByID(ctx, walletID, userID)
	if err != nil {
		return nil, err
	}

	// Validate new limits
	if err := ws.validateLimits(req.Limits); err != nil {
		return nil, fmt.Errorf("invalid limits: %w", err)
	}

	// Track changes for event
	oldLimits := wallet.Limits

	// Update limits
	wallet.Limits = req.Limits
	wallet.UpdatedAt = time.Now()

	// Save updated wallet
	if err := ws.walletRepo.Update(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to update wallet limits: %w", err)
	}

	// Publish wallet limits updated event
	if err := ws.publishWalletEvent(ctx, "wallet.limits_updated", walletID, userID, map[string]interface{}{
		"old_limits": oldLimits,
		"new_limits": req.Limits,
	}); err != nil {
		fmt.Printf("Failed to publish wallet limits updated event: %v\n", err)
	}

	return wallet, nil
}

func (ws *WalletService) FreezeWallet(ctx context.Context, walletID uuid.UUID, userID uuid.UUID, reason string) error {
	wallet, err := ws.GetWalletByID(ctx, walletID, userID)
	if err != nil {
		return err
	}

	if !wallet.CanFreeze() {
		return fmt.Errorf("wallet cannot be frozen in current status: %s", wallet.Status)
	}

	wallet.Status = models.WalletStatusFrozen
	wallet.UpdatedAt = time.Now()

	if err := ws.walletRepo.Update(ctx, wallet); err != nil {
		return fmt.Errorf("failed to freeze wallet: %w", err)
	}

	// Publish wallet frozen event
	if err := ws.publishWalletEvent(ctx, "wallet.frozen", walletID, userID, map[string]interface{}{
		"reason": reason,
	}); err != nil {
		fmt.Printf("Failed to publish wallet frozen event: %v\n", err)
	}

	return nil
}

func (ws *WalletService) UnfreezeWallet(ctx context.Context, walletID uuid.UUID, userID uuid.UUID, reason string) error {
	wallet, err := ws.GetWalletByID(ctx, walletID, userID)
	if err != nil {
		return err
	}

	if wallet.Status != models.WalletStatusFrozen {
		return fmt.Errorf("wallet is not frozen")
	}

	wallet.Status = models.WalletStatusActive
	wallet.UpdatedAt = time.Now()

	if err := ws.walletRepo.Update(ctx, wallet); err != nil {
		return fmt.Errorf("failed to unfreeze wallet: %w", err)
	}

	// Publish wallet unfrozen event
	if err := ws.publishWalletEvent(ctx, "wallet.unfrozen", walletID, userID, map[string]interface{}{
		"reason": reason,
	}); err != nil {
		fmt.Printf("Failed to publish wallet unfrozen event: %v\n", err)
	}

	return nil
}

func (ws *WalletService) AddFunds(ctx context.Context, walletID uuid.UUID, amount float64, transactionID uuid.UUID, source string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	wallet, err := ws.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	if !wallet.CanAddFunds() {
		return fmt.Errorf("cannot add funds to wallet in status: %s", wallet.Status)
	}

	// Check if adding funds would exceed max balance
	if wallet.Balance+amount > wallet.Limits.MaxBalance {
		return fmt.Errorf("adding funds would exceed maximum balance limit")
	}

	// Check compliance
	if err := ws.complianceChecker.CheckAMLCompliance(ctx, wallet.UserID, amount, wallet.Currency); err != nil {
		return fmt.Errorf("AML compliance check failed: %w", err)
	}

	// Update balance atomically
	if err := ws.walletRepo.UpdateBalance(ctx, walletID, amount, transactionID); err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Publish wallet balance changed event
	if err := ws.publishWalletEvent(ctx, sharedModels.EventTypeWalletBalanceChanged, walletID, wallet.UserID, map[string]interface{}{
		"amount":         amount,
		"transaction_id": transactionID,
		"source":         source,
		"operation":      "add_funds",
		"new_balance":    wallet.Balance + amount,
	}); err != nil {
		fmt.Printf("Failed to publish wallet balance changed event: %v\n", err)
	}

	return nil
}

func (ws *WalletService) DeductFunds(ctx context.Context, walletID uuid.UUID, amount float64, transactionID uuid.UUID, purpose string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	wallet, err := ws.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	if !wallet.CanDeductFunds() {
		return fmt.Errorf("cannot deduct funds from wallet in status: %s", wallet.Status)
	}

	// Check if sufficient balance
	if wallet.AvailableBalance < amount {
		return fmt.Errorf("insufficient balance: available %.2f, required %.2f", wallet.AvailableBalance, amount)
	}

	// Check compliance
	if err := ws.complianceChecker.CheckTransactionLimits(ctx, wallet.UserID, amount, wallet.Currency, "debit"); err != nil {
		return fmt.Errorf("transaction limits check failed: %w", err)
	}

	// Update balance atomically (negative amount for deduction)
	if err := ws.walletRepo.UpdateBalance(ctx, walletID, -amount, transactionID); err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Publish wallet balance changed event
	if err := ws.publishWalletEvent(ctx, sharedModels.EventTypeWalletBalanceChanged, walletID, wallet.UserID, map[string]interface{}{
		"amount":         -amount,
		"transaction_id": transactionID,
		"purpose":        purpose,
		"operation":      "deduct_funds",
		"new_balance":    wallet.Balance - amount,
	}); err != nil {
		fmt.Printf("Failed to publish wallet balance changed event: %v\n", err)
	}

	return nil
}

func (ws *WalletService) FreezeAmount(ctx context.Context, walletID uuid.UUID, amount float64, reason string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	wallet, err := ws.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	if wallet.AvailableBalance < amount {
		return fmt.Errorf("insufficient available balance to freeze")
	}

	if err := ws.walletRepo.FreezeAmount(ctx, walletID, amount, reason); err != nil {
		return fmt.Errorf("failed to freeze amount: %w", err)
	}

	// Publish amount frozen event
	if err := ws.publishWalletEvent(ctx, "wallet.amount_frozen", walletID, wallet.UserID, map[string]interface{}{
		"amount": amount,
		"reason": reason,
	}); err != nil {
		fmt.Printf("Failed to publish amount frozen event: %v\n", err)
	}

	return nil
}

func (ws *WalletService) UnfreezeAmount(ctx context.Context, walletID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	wallet, err := ws.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	if wallet.FrozenBalance < amount {
		return fmt.Errorf("insufficient frozen balance to unfreeze")
	}

	if err := ws.walletRepo.UnfreezeAmount(ctx, walletID, amount); err != nil {
		return fmt.Errorf("failed to unfreeze amount: %w", err)
	}

	// Publish amount unfrozen event
	if err := ws.publishWalletEvent(ctx, "wallet.amount_unfrozen", walletID, wallet.UserID, map[string]interface{}{
		"amount": amount,
	}); err != nil {
		fmt.Printf("Failed to publish amount unfrozen event: %v\n", err)
	}

	return nil
}

func (ws *WalletService) TransferFunds(ctx context.Context, req *TransferFundsRequest) (*TransferResult, error) {
	// Validate request
	if err := ws.validateTransferRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get source and destination wallets
	sourceWallet, err := ws.walletRepo.GetByID(ctx, req.FromWalletID)
	if err != nil {
		return nil, fmt.Errorf("source wallet not found: %w", err)
	}

	destWallet, err := ws.walletRepo.GetByID(ctx, req.ToWalletID)
	if err != nil {
		return nil, fmt.Errorf("destination wallet not found: %w", err)
	}

	// Check if user owns source wallet
	if sourceWallet.UserID != req.UserID {
		return nil, fmt.Errorf("access denied: not owner of source wallet")
	}

	transferAmount := req.Amount

	// Handle currency conversion if needed
	if sourceWallet.Currency != destWallet.Currency {
		convertedAmount, err := ws.exchangeRateService.ConvertAmount(ctx, req.Amount, sourceWallet.Currency, destWallet.Currency)
		if err != nil {
			return nil, fmt.Errorf("currency conversion failed: %w", err)
		}
		transferAmount = convertedAmount
	}

	// Check source wallet balance
	if sourceWallet.AvailableBalance < req.Amount {
		return nil, fmt.Errorf("insufficient balance in source wallet")
	}

	// Check destination wallet limits
	if destWallet.Balance+transferAmount > destWallet.Limits.MaxBalance {
		return nil, fmt.Errorf("transfer would exceed destination wallet's maximum balance")
	}

	// Check compliance for both wallets
	if err := ws.complianceChecker.CheckTransactionLimits(ctx, sourceWallet.UserID, req.Amount, sourceWallet.Currency, "transfer_out"); err != nil {
		return nil, fmt.Errorf("source wallet transaction limits check failed: %w", err)
	}

	if err := ws.complianceChecker.CheckTransactionLimits(ctx, destWallet.UserID, transferAmount, destWallet.Currency, "transfer_in"); err != nil {
		return nil, fmt.Errorf("destination wallet transaction limits check failed: %w", err)
	}

	// Create transfer transaction ID
	transferID := uuid.New()

	// Begin database transaction for atomic transfer
	tx := ws.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deduct from source wallet
	if err := ws.walletRepo.UpdateBalance(ctx, req.FromWalletID, -req.Amount, transferID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to deduct from source wallet: %w", err)
	}

	// Add to destination wallet
	if err := ws.walletRepo.UpdateBalance(ctx, req.ToWalletID, transferAmount, transferID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to add to destination wallet: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transfer: %w", err)
	}

	result := &TransferResult{
		TransferID:       transferID,
		FromWalletID:     req.FromWalletID,
		ToWalletID:       req.ToWalletID,
		OriginalAmount:   req.Amount,
		TransferredAmount: transferAmount,
		SourceCurrency:   sourceWallet.Currency,
		DestCurrency:     destWallet.Currency,
		TransferredAt:    time.Now(),
	}

	// Publish transfer events
	if err := ws.publishWalletEvent(ctx, "wallet.transfer_completed", req.FromWalletID, req.UserID, map[string]interface{}{
		"transfer_id":        transferID,
		"from_wallet":        req.FromWalletID,
		"to_wallet":          req.ToWalletID,
		"original_amount":    req.Amount,
		"transferred_amount": transferAmount,
		"source_currency":    sourceWallet.Currency,
		"dest_currency":      destWallet.Currency,
	}); err != nil {
		fmt.Printf("Failed to publish wallet transfer event: %v\n", err)
	}

	return result, nil
}

func (ws *WalletService) GetWalletHistory(ctx context.Context, walletID uuid.UUID, userID uuid.UUID, limit int) ([]*WalletBalanceHistory, error) {
	// Verify wallet ownership
	_, err := ws.GetWalletByID(ctx, walletID, userID)
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	history, err := ws.walletRepo.GetWalletHistory(ctx, walletID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet history: %w", err)
	}

	return history, nil
}

func (ws *WalletService) GetWalletStats(ctx context.Context, userID uuid.UUID) (*WalletStats, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	stats, err := ws.walletRepo.GetWalletStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet stats: %w", err)
	}

	return stats, nil
}

// Private helper methods

func (ws *WalletService) validateCreateWalletRequest(req *CreateWalletRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	// Validate currency
	if !models.IsValidCurrency(req.Currency) {
		return fmt.Errorf("invalid currency: %s", req.Currency)
	}

	// Validate limits
	return ws.validateLimits(models.WalletLimits{
		DailySpendLimit:   req.DailySpendLimit,
		MonthlySpendLimit: req.MonthlySpendLimit,
		DailyTopupLimit:   req.DailyTopupLimit,
		MonthlyTopupLimit: req.MonthlyTopupLimit,
		MaxBalance:        req.MaxBalance,
	})
}

func (ws *WalletService) validateLimits(limits models.WalletLimits) error {
	if limits.DailySpendLimit < 0 {
		return fmt.Errorf("daily spend limit cannot be negative")
	}

	if limits.MonthlySpendLimit < 0 {
		return fmt.Errorf("monthly spend limit cannot be negative")
	}

	if limits.DailyTopupLimit < 0 {
		return fmt.Errorf("daily topup limit cannot be negative")
	}

	if limits.MonthlyTopupLimit < 0 {
		return fmt.Errorf("monthly topup limit cannot be negative")
	}

	if limits.MaxBalance < 0 {
		return fmt.Errorf("max balance cannot be negative")
	}

	if limits.DailySpendLimit > limits.MonthlySpendLimit {
		return fmt.Errorf("daily spend limit cannot exceed monthly spend limit")
	}

	if limits.DailyTopupLimit > limits.MonthlyTopupLimit {
		return fmt.Errorf("daily topup limit cannot exceed monthly topup limit")
	}

	return nil
}

func (ws *WalletService) validateTransferRequest(req *TransferFundsRequest) error {
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if req.FromWalletID == uuid.Nil {
		return fmt.Errorf("source wallet ID is required")
	}

	if req.ToWalletID == uuid.Nil {
		return fmt.Errorf("destination wallet ID is required")
	}

	if req.FromWalletID == req.ToWalletID {
		return fmt.Errorf("source and destination wallets cannot be the same")
	}

	if req.Amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	return nil
}

func (ws *WalletService) publishWalletEvent(ctx context.Context, eventType sharedModels.EventType, walletID uuid.UUID, userID uuid.UUID, data map[string]interface{}) error {
	event := &sharedModels.Event{
		ID:            uuid.New(),
		Type:          eventType,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       fmt.Sprintf("Wallet event: %s", eventType),
		AggregateID:   walletID.String(),
		AggregateType: "wallet",
		EventVersion:  1,
		OccurredAt:    time.Now(),
		Status:        sharedModels.EventStatusPending,
		Metadata: sharedModels.EventMetadata{
			Source:      "payment-service",
			Environment: "production",
			Region:      "sea",
		},
	}

	// Add user context to data
	data["user_id"] = userID

	if err := event.MarshalData(data); err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	return ws.eventPublisher.Publish(ctx, event)
}

// Request/Response structures

type CreateWalletRequest struct {
	UserID            uuid.UUID              `json:"user_id" binding:"required"`
	Currency          models.Currency        `json:"currency" binding:"required"`
	DailySpendLimit   float64                `json:"daily_spend_limit"`
	MonthlySpendLimit float64                `json:"monthly_spend_limit"`
	DailyTopupLimit   float64                `json:"daily_topup_limit"`
	MonthlyTopupLimit float64                `json:"monthly_topup_limit"`
	MaxBalance        float64                `json:"max_balance"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type UpdateWalletLimitsRequest struct {
	Limits models.WalletLimits `json:"limits" binding:"required"`
}

type TransferFundsRequest struct {
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	FromWalletID uuid.UUID `json:"from_wallet_id" binding:"required"`
	ToWalletID   uuid.UUID `json:"to_wallet_id" binding:"required"`
	Amount       float64   `json:"amount" binding:"required"`
	Description  string    `json:"description"`
}

type TransferResult struct {
	TransferID        uuid.UUID       `json:"transfer_id"`
	FromWalletID      uuid.UUID       `json:"from_wallet_id"`
	ToWalletID        uuid.UUID       `json:"to_wallet_id"`
	OriginalAmount    float64         `json:"original_amount"`
	TransferredAmount float64         `json:"transferred_amount"`
	SourceCurrency    models.Currency `json:"source_currency"`
	DestCurrency      models.Currency `json:"dest_currency"`
	TransferredAt     time.Time       `json:"transferred_at"`
}

type WalletResponse struct {
	ID               uuid.UUID              `json:"id"`
	UserID           uuid.UUID              `json:"user_id"`
	Currency         models.Currency        `json:"currency"`
	Balance          float64                `json:"balance"`
	AvailableBalance float64                `json:"available_balance"`
	FrozenBalance    float64                `json:"frozen_balance"`
	FormattedBalance string                 `json:"formatted_balance"`
	Status           models.WalletStatus    `json:"status"`
	Limits           models.WalletLimits    `json:"limits"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type WalletBalanceResponse struct {
	WalletID         uuid.UUID       `json:"wallet_id"`
	Currency         models.Currency `json:"currency"`
	Balance          float64         `json:"balance"`
	AvailableBalance float64         `json:"available_balance"`
	FrozenBalance    float64         `json:"frozen_balance"`
	FormattedBalance string          `json:"formatted_balance"`
	LastUpdated      time.Time       `json:"last_updated"`
}

type WalletListResponse struct {
	Wallets []*WalletResponse `json:"wallets"`
	Total   int               `json:"total"`
}

func (wallet *models.Wallet) ToResponse() *WalletResponse {
	return &WalletResponse{
		ID:               wallet.ID,
		UserID:           wallet.UserID,
		Currency:         wallet.Currency,
		Balance:          wallet.Balance,
		AvailableBalance: wallet.AvailableBalance,
		FrozenBalance:    wallet.FrozenBalance,
		FormattedBalance: wallet.FormatBalance(),
		Status:           wallet.Status,
		Limits:           wallet.Limits,
		CreatedAt:        wallet.CreatedAt,
		UpdatedAt:        wallet.UpdatedAt,
	}
}