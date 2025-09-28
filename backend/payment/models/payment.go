package models

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// Payment represents a payment transaction
type Payment struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	WalletID    uuid.UUID     `json:"wallet_id" db:"wallet_id"`
	Amount      int64         `json:"amount" db:"amount"`
	Currency    Currency      `json:"currency" db:"currency"`
	Gateway     string        `json:"gateway" db:"gateway"`
	Method      string        `json:"method" db:"method"`
	Status      PaymentStatus `json:"status" db:"status"`
	Description string        `json:"description" db:"description"`
	ReturnURL   string        `json:"return_url" db:"return_url"`
	Metadata    []byte        `json:"metadata" db:"metadata"`
	ExternalID  *string       `json:"external_id,omitempty" db:"external_id"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

// IsValid validates the payment
func (p *Payment) IsValid() bool {
	return p.ID != uuid.Nil &&
		p.WalletID != uuid.Nil &&
		p.Amount > 0 &&
		p.Gateway != "" &&
		p.Method != ""
}