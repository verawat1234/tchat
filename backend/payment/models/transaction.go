package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Transaction represents a payment transaction in the system
type Transaction struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	WalletID          uuid.UUID         `json:"wallet_id" db:"wallet_id"`
	CounterpartyID    *uuid.UUID        `json:"counterparty_id,omitempty" db:"counterparty_id"`
	Type              TransactionType   `json:"type" db:"type"`
	Status            TransactionStatus `json:"status" db:"status"`
	Currency          Currency          `json:"currency" db:"currency"`
	Amount            int64             `json:"amount" db:"amount"` // Amount in currency's smallest unit
	FeeAmount         int64             `json:"fee_amount" db:"fee_amount"`
	NetAmount         int64             `json:"net_amount" db:"net_amount"` // Amount after fees
	ExchangeRate      *float64          `json:"exchange_rate,omitempty" db:"exchange_rate"`
	Reference         string            `json:"reference" db:"reference"`
	Description       *string           `json:"description,omitempty" db:"description"`
	Metadata          map[string]string `json:"metadata,omitempty" db:"metadata"`
	ExternalID        *string           `json:"external_id,omitempty" db:"external_id"`
	ProcessorResponse *string           `json:"processor_response,omitempty" db:"processor_response"`
	FailureReason     *string           `json:"failure_reason,omitempty" db:"failure_reason"`
	ProcessedAt       *time.Time        `json:"processed_at,omitempty" db:"processed_at"`
	CompletedAt       *time.Time        `json:"completed_at,omitempty" db:"completed_at"`
	ExpiresAt         *time.Time        `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeRefund     TransactionType = "refund"
	TransactionTypeFee        TransactionType = "fee"
	TransactionTypeExchange   TransactionType = "exchange"
	TransactionTypeReward     TransactionType = "reward"
	TransactionTypePenalty    TransactionType = "penalty"
)

// ValidTransactionTypes returns all supported transaction types
func ValidTransactionTypes() []TransactionType {
	return []TransactionType{
		TransactionTypeDeposit,
		TransactionTypeWithdrawal,
		TransactionTypeTransfer,
		TransactionTypePayment,
		TransactionTypeRefund,
		TransactionTypeFee,
		TransactionTypeExchange,
		TransactionTypeReward,
		TransactionTypePenalty,
	}
}

// IsValid validates if the transaction type is supported
func (t TransactionType) IsValid() bool {
	for _, valid := range ValidTransactionTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of TransactionType
func (t TransactionType) String() string {
	return string(t)
}

// TransactionStatus represents the current state of a transaction
type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusProcessing TransactionStatus = "processing"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
	TransactionStatusExpired    TransactionStatus = "expired"
	TransactionStatusRefunded   TransactionStatus = "refunded"
)

// ValidTransactionStatuses returns all supported transaction statuses
func ValidTransactionStatuses() []TransactionStatus {
	return []TransactionStatus{
		TransactionStatusPending,
		TransactionStatusProcessing,
		TransactionStatusCompleted,
		TransactionStatusFailed,
		TransactionStatusCancelled,
		TransactionStatusExpired,
		TransactionStatusRefunded,
	}
}

// IsValid validates if the transaction status is supported
func (s TransactionStatus) IsValid() bool {
	for _, valid := range ValidTransactionStatuses() {
		if s == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of TransactionStatus
func (s TransactionStatus) String() string {
	return string(s)
}

// Transaction configuration constants
const (
	MaxTransactionAmount   = 100000000000 // 1 billion in smallest unit
	DefaultTransactionTTL  = 24 * time.Hour
	MaxDescriptionLength   = 500
	MaxReferenceLength     = 100
	MaxExternalIDLength    = 255
	MaxFailureReasonLength = 1000
	MaxMetadataKeys        = 20
	MaxMetadataValueLength = 500
)

// TransactionValidationError represents transaction validation errors
type TransactionValidationError struct {
	Field   string
	Message string
}

func (e TransactionValidationError) Error() string {
	return fmt.Sprintf("transaction validation error - %s: %s", e.Field, e.Message)
}

// Validate performs comprehensive validation on the Transaction model
func (t *Transaction) Validate() error {
	var errs []string

	// Wallet ID validation
	if t.WalletID == uuid.Nil {
		errs = append(errs, "wallet_id is required")
	}

	// Transaction type validation
	if !t.Type.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid transaction type: %s", t.Type))
	}

	// Transaction status validation
	if !t.Status.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid transaction status: %s", t.Status))
	}

	// Currency validation
	if !t.Currency.IsValid() {
		errs = append(errs, fmt.Sprintf("invalid currency: %s", t.Currency))
	}

	// Amount validation
	if t.Amount <= 0 {
		errs = append(errs, "amount must be positive")
	}
	if t.Amount > MaxTransactionAmount {
		errs = append(errs, fmt.Sprintf("amount exceeds maximum limit: %d", MaxTransactionAmount))
	}

	// Fee amount validation
	if t.FeeAmount < 0 {
		errs = append(errs, "fee_amount cannot be negative")
	}
	if t.FeeAmount >= t.Amount {
		errs = append(errs, "fee_amount cannot exceed transaction amount")
	}

	// Net amount validation
	expectedNetAmount := t.Amount - t.FeeAmount
	if t.NetAmount != expectedNetAmount {
		errs = append(errs, fmt.Sprintf("net_amount mismatch: expected %d, got %d", expectedNetAmount, t.NetAmount))
	}

	// Exchange rate validation
	if t.ExchangeRate != nil && *t.ExchangeRate <= 0 {
		errs = append(errs, "exchange_rate must be positive when provided")
	}

	// Reference validation
	if strings.TrimSpace(t.Reference) == "" {
		errs = append(errs, "reference is required")
	}
	if len(t.Reference) > MaxReferenceLength {
		errs = append(errs, fmt.Sprintf("reference must not exceed %d characters", MaxReferenceLength))
	}

	// Description validation
	if t.Description != nil && len(*t.Description) > MaxDescriptionLength {
		errs = append(errs, fmt.Sprintf("description must not exceed %d characters", MaxDescriptionLength))
	}

	// External ID validation
	if t.ExternalID != nil && len(*t.ExternalID) > MaxExternalIDLength {
		errs = append(errs, fmt.Sprintf("external_id must not exceed %d characters", MaxExternalIDLength))
	}

	// Failure reason validation
	if t.FailureReason != nil && len(*t.FailureReason) > MaxFailureReasonLength {
		errs = append(errs, fmt.Sprintf("failure_reason must not exceed %d characters", MaxFailureReasonLength))
	}

	// Metadata validation
	if t.Metadata != nil {
		if len(t.Metadata) > MaxMetadataKeys {
			errs = append(errs, fmt.Sprintf("metadata cannot have more than %d keys", MaxMetadataKeys))
		}
		for key, value := range t.Metadata {
			if len(key) == 0 {
				errs = append(errs, "metadata keys cannot be empty")
			}
			if len(value) > MaxMetadataValueLength {
				errs = append(errs, fmt.Sprintf("metadata value for key '%s' exceeds maximum length", key))
			}
		}
	}

	// Counterparty validation for specific transaction types
	if err := t.validateCounterpartyRequirements(); err != nil {
		errs = append(errs, err.Error())
	}

	// Status-specific validations
	if err := t.validateStatusConsistency(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// validateCounterpartyRequirements validates counterparty requirements based on transaction type
func (t *Transaction) validateCounterpartyRequirements() error {
	requiresCounterparty := []TransactionType{
		TransactionTypeTransfer,
		TransactionTypePayment,
		TransactionTypeRefund,
	}

	for _, txType := range requiresCounterparty {
		if t.Type == txType && t.CounterpartyID == nil {
			return fmt.Errorf("counterparty_id is required for %s transactions", txType)
		}
	}

	prohibitsCounterparty := []TransactionType{
		TransactionTypeDeposit,
		TransactionTypeWithdrawal,
		TransactionTypeFee,
		TransactionTypeReward,
		TransactionTypePenalty,
	}

	for _, txType := range prohibitsCounterparty {
		if t.Type == txType && t.CounterpartyID != nil {
			return fmt.Errorf("counterparty_id is not allowed for %s transactions", txType)
		}
	}

	return nil
}

// validateStatusConsistency validates status-specific field consistency
func (t *Transaction) validateStatusConsistency() error {
	switch t.Status {
	case TransactionStatusCompleted:
		if t.CompletedAt == nil {
			return errors.New("completed_at is required for completed transactions")
		}
		if t.ProcessedAt == nil {
			return errors.New("processed_at is required for completed transactions")
		}
	case TransactionStatusFailed:
		if t.FailureReason == nil || strings.TrimSpace(*t.FailureReason) == "" {
			return errors.New("failure_reason is required for failed transactions")
		}
		if t.ProcessedAt == nil {
			return errors.New("processed_at is required for failed transactions")
		}
	case TransactionStatusExpired:
		if t.ExpiresAt == nil {
			return errors.New("expires_at is required for expired transactions")
		}
		if !time.Now().UTC().After(*t.ExpiresAt) {
			return errors.New("transaction marked as expired but expiry time has not passed")
		}
	case TransactionStatusRefunded:
		if t.CompletedAt == nil {
			return errors.New("completed_at is required for refunded transactions")
		}
	}

	return nil
}

// BeforeCreate sets up the transaction before database creation
func (t *Transaction) BeforeCreate() error {
	// Generate UUID if not set
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	// Set default status
	if t.Status == "" {
		t.Status = TransactionStatusPending
	}

	// Generate reference if not provided
	if t.Reference == "" {
		t.Reference = t.generateReference()
	}

	// Calculate net amount
	t.NetAmount = t.Amount - t.FeeAmount

	// Set default expiry for pending transactions
	if t.Status == TransactionStatusPending && t.ExpiresAt == nil {
		expiryTime := now.Add(DefaultTransactionTTL)
		t.ExpiresAt = &expiryTime
	}

	// Initialize metadata if nil
	if t.Metadata == nil {
		t.Metadata = make(map[string]string)
	}

	// Validate before creation
	return t.Validate()
}

// BeforeUpdate sets up the transaction before database update
func (t *Transaction) BeforeUpdate() error {
	// Update timestamp
	t.UpdatedAt = time.Now().UTC()

	// Recalculate net amount
	t.NetAmount = t.Amount - t.FeeAmount

	// Validate before update
	return t.Validate()
}

// generateReference generates a unique transaction reference
func (t *Transaction) generateReference() string {
	timestamp := time.Now().UTC().Format("20060102150405")
	shortID := t.ID.String()[:8]
	typePrefix := strings.ToUpper(string(t.Type)[:3])
	return fmt.Sprintf("%s-%s-%s", typePrefix, timestamp, shortID)
}

// State machine methods for transaction status transitions

// CanTransitionTo checks if transition to target status is allowed
func (t *Transaction) CanTransitionTo(targetStatus TransactionStatus) bool {
	allowedTransitions := map[TransactionStatus][]TransactionStatus{
		TransactionStatusPending: {
			TransactionStatusProcessing,
			TransactionStatusCancelled,
			TransactionStatusExpired,
		},
		TransactionStatusProcessing: {
			TransactionStatusCompleted,
			TransactionStatusFailed,
		},
		TransactionStatusCompleted: {
			TransactionStatusRefunded,
		},
		TransactionStatusFailed:    {},
		TransactionStatusCancelled: {},
		TransactionStatusExpired:   {},
		TransactionStatusRefunded:  {},
	}

	allowed, exists := allowedTransitions[t.Status]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == targetStatus {
			return true
		}
	}
	return false
}

// StartProcessing transitions transaction to processing status
func (t *Transaction) StartProcessing() error {
	if !t.CanTransitionTo(TransactionStatusProcessing) {
		return fmt.Errorf("cannot transition from %s to processing", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TransactionStatusProcessing
	t.ProcessedAt = &now
	t.UpdatedAt = now

	return t.Validate()
}

// Complete transitions transaction to completed status
func (t *Transaction) Complete() error {
	if !t.CanTransitionTo(TransactionStatusCompleted) {
		return fmt.Errorf("cannot transition from %s to completed", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TransactionStatusCompleted
	t.CompletedAt = &now
	t.UpdatedAt = now

	// Clear expiry time for completed transactions
	t.ExpiresAt = nil

	return t.Validate()
}

// Fail transitions transaction to failed status
func (t *Transaction) Fail(reason string) error {
	if !t.CanTransitionTo(TransactionStatusFailed) {
		return fmt.Errorf("cannot transition from %s to failed", t.Status)
	}

	if strings.TrimSpace(reason) == "" {
		return errors.New("failure reason is required")
	}

	now := time.Now().UTC()
	t.Status = TransactionStatusFailed
	t.FailureReason = &reason
	t.UpdatedAt = now

	// Ensure processed timestamp is set
	if t.ProcessedAt == nil {
		t.ProcessedAt = &now
	}

	return t.Validate()
}

// Cancel transitions transaction to cancelled status
func (t *Transaction) Cancel() error {
	if !t.CanTransitionTo(TransactionStatusCancelled) {
		return fmt.Errorf("cannot transition from %s to cancelled", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TransactionStatusCancelled
	t.UpdatedAt = now

	return t.Validate()
}

// Expire transitions transaction to expired status
func (t *Transaction) Expire() error {
	if !t.CanTransitionTo(TransactionStatusExpired) {
		return fmt.Errorf("cannot transition from %s to expired", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TransactionStatusExpired
	t.UpdatedAt = now

	// Ensure expiry time is set and in the past
	if t.ExpiresAt == nil {
		t.ExpiresAt = &now
	}

	return t.Validate()
}

// Refund transitions transaction to refunded status
func (t *Transaction) Refund() error {
	if !t.CanTransitionTo(TransactionStatusRefunded) {
		return fmt.Errorf("cannot transition from %s to refunded", t.Status)
	}

	now := time.Now().UTC()
	t.Status = TransactionStatusRefunded
	t.UpdatedAt = now

	return t.Validate()
}

// Utility methods

// IsCompleted checks if the transaction is in a completed state
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// IsFailed checks if the transaction is in a failed state
func (t *Transaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed
}

// IsPending checks if the transaction is pending
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// IsProcessing checks if the transaction is currently processing
func (t *Transaction) IsProcessing() bool {
	return t.Status == TransactionStatusProcessing
}

// IsCancelled checks if the transaction was cancelled
func (t *Transaction) IsCancelled() bool {
	return t.Status == TransactionStatusCancelled
}

// IsExpired checks if the transaction has expired
func (t *Transaction) IsExpired() bool {
	if t.Status == TransactionStatusExpired {
		return true
	}
	if t.ExpiresAt != nil && time.Now().UTC().After(*t.ExpiresAt) {
		return true
	}
	return false
}

// IsRefunded checks if the transaction was refunded
func (t *Transaction) IsRefunded() bool {
	return t.Status == TransactionStatusRefunded
}

// IsFinalState checks if the transaction is in a final state (cannot be changed)
func (t *Transaction) IsFinalState() bool {
	finalStates := []TransactionStatus{
		TransactionStatusCompleted,
		TransactionStatusFailed,
		TransactionStatusCancelled,
		TransactionStatusExpired,
		TransactionStatusRefunded,
	}

	for _, state := range finalStates {
		if t.Status == state {
			return true
		}
	}
	return false
}

// GetFormattedAmount returns the amount formatted according to currency specifications
func (t *Transaction) GetFormattedAmount() string {
	return t.Currency.FormatAmount(t.Amount)
}

// GetFormattedFee returns the fee amount formatted according to currency specifications
func (t *Transaction) GetFormattedFee() string {
	return t.Currency.FormatAmount(t.FeeAmount)
}

// GetFormattedNetAmount returns the net amount formatted according to currency specifications
func (t *Transaction) GetFormattedNetAmount() string {
	return t.Currency.FormatAmount(t.NetAmount)
}

// GetDuration returns the duration of the transaction processing
func (t *Transaction) GetDuration() *time.Duration {
	if t.ProcessedAt == nil {
		return nil
	}
	duration := t.ProcessedAt.Sub(t.CreatedAt)
	return &duration
}

// AddMetadata adds or updates metadata key-value pair
func (t *Transaction) AddMetadata(key, value string) error {
	if t.Metadata == nil {
		t.Metadata = make(map[string]string)
	}

	if len(t.Metadata) >= MaxMetadataKeys && t.Metadata[key] == "" {
		return fmt.Errorf("cannot add more than %d metadata keys", MaxMetadataKeys)
	}

	if len(value) > MaxMetadataValueLength {
		return fmt.Errorf("metadata value exceeds maximum length of %d characters", MaxMetadataValueLength)
	}

	t.Metadata[key] = value
	return nil
}

// RemoveMetadata removes a metadata key
func (t *Transaction) RemoveMetadata(key string) {
	if t.Metadata != nil {
		delete(t.Metadata, key)
	}
}

// ToAuditLog returns an audit log representation of the transaction
func (t *Transaction) ToAuditLog() map[string]interface{} {
	return map[string]interface{}{
		"id":              t.ID,
		"wallet_id":       t.WalletID,
		"counterparty_id": t.CounterpartyID,
		"type":            t.Type,
		"status":          t.Status,
		"currency":        t.Currency,
		"amount":          t.Amount,
		"fee_amount":      t.FeeAmount,
		"net_amount":      t.NetAmount,
		"reference":       t.Reference,
		"external_id":     t.ExternalID,
		"created_at":      t.CreatedAt,
		"updated_at":      t.UpdatedAt,
		"processed_at":    t.ProcessedAt,
		"completed_at":    t.CompletedAt,
		"expires_at":      t.ExpiresAt,
	}
}

// ToPublicTransaction returns a sanitized version for public API responses
func (t *Transaction) ToPublicTransaction() map[string]interface{} {
	return map[string]interface{}{
		"id":               t.ID,
		"type":             t.Type,
		"status":           t.Status,
		"currency":         t.Currency,
		"amount":           t.Amount,
		"formatted_amount": t.GetFormattedAmount(),
		"fee_amount":       t.FeeAmount,
		"formatted_fee":    t.GetFormattedFee(),
		"net_amount":       t.NetAmount,
		"formatted_net":    t.GetFormattedNetAmount(),
		"reference":        t.Reference,
		"description":      t.Description,
		"created_at":       t.CreatedAt,
		"completed_at":     t.CompletedAt,
		"expires_at":       t.ExpiresAt,
	}
}

// TransactionCreateRequest represents a request to create a new transaction
type TransactionCreateRequest struct {
	WalletID       uuid.UUID         `json:"wallet_id" validate:"required"`
	CounterpartyID *uuid.UUID        `json:"counterparty_id,omitempty"`
	Type           TransactionType   `json:"type" validate:"required"`
	Currency       Currency          `json:"currency" validate:"required"`
	Amount         int64             `json:"amount" validate:"required,min=1"`
	FeeAmount      int64             `json:"fee_amount,omitempty"`
	Description    *string           `json:"description,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	ExternalID     *string           `json:"external_id,omitempty"`
	ExpiresAt      *time.Time        `json:"expires_at,omitempty"`
}

// ToTransaction converts a create request to a Transaction model
func (req *TransactionCreateRequest) ToTransaction() *Transaction {
	return &Transaction{
		WalletID:       req.WalletID,
		CounterpartyID: req.CounterpartyID,
		Type:           req.Type,
		Currency:       req.Currency,
		Amount:         req.Amount,
		FeeAmount:      req.FeeAmount,
		Description:    req.Description,
		Metadata:       req.Metadata,
		ExternalID:     req.ExternalID,
		ExpiresAt:      req.ExpiresAt,
	}
}

// TransactionManager provides transaction management utilities
type TransactionManager struct {
	// Add dependencies like database, payment processors, etc.
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

// CreateTransaction creates a new transaction with proper validation
func (tm *TransactionManager) CreateTransaction(req *TransactionCreateRequest) (*Transaction, error) {
	transaction := req.ToTransaction()

	if err := transaction.BeforeCreate(); err != nil {
		return nil, fmt.Errorf("transaction creation failed: %v", err)
	}

	return transaction, nil
}

// ProcessTransaction handles transaction processing workflow
func (tm *TransactionManager) ProcessTransaction(transaction *Transaction) error {
	if !transaction.IsPending() {
		return fmt.Errorf("transaction must be pending to process")
	}

	// Start processing
	if err := transaction.StartProcessing(); err != nil {
		return fmt.Errorf("failed to start processing: %v", err)
	}

	// Additional processing logic would go here
	// (payment processor integration, wallet balance checks, etc.)

	return nil
}

// GetExpiredTransactions returns filter criteria for expired transaction cleanup
func (tm *TransactionManager) GetExpiredTransactions() map[string]interface{} {
	now := time.Now().UTC()
	return map[string]interface{}{
		"status":      TransactionStatusPending,
		"expires_at":  now,
		"operator":    "lt", // less than
	}
}

// ValidateTransactionLimits validates transaction against wallet limits
func (tm *TransactionManager) ValidateTransactionLimits(transaction *Transaction, wallet *Wallet) error {
	if wallet == nil {
		return errors.New("wallet is required")
	}

	// Check currency compatibility
	if transaction.Currency != wallet.Currency {
		return fmt.Errorf("transaction currency %s does not match wallet currency %s",
			transaction.Currency, wallet.Currency)
	}

	// Check daily limits for outgoing transactions
	if transaction.Type == TransactionTypeWithdrawal ||
		transaction.Type == TransactionTypeTransfer ||
		transaction.Type == TransactionTypePayment {

		if transaction.Amount > wallet.DailyLimit {
			return fmt.Errorf("transaction amount %d exceeds daily limit %d",
				transaction.Amount, wallet.DailyLimit)
		}
	}

	// Check sufficient balance for outgoing transactions
	requiredAmount := transaction.Amount
	if transaction.Type == TransactionTypeWithdrawal ||
		transaction.Type == TransactionTypeTransfer ||
		transaction.Type == TransactionTypePayment {

		availableBalance := wallet.Balance - wallet.FrozenBalance
		if requiredAmount > availableBalance {
			return fmt.Errorf("insufficient balance: required %d, available %d",
				requiredAmount, availableBalance)
		}
	}

	return nil
}