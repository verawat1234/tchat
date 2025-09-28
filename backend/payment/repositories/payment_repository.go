package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tchat.dev/payment/models"
)

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	Update(ctx context.Context, payment *models.Payment) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByWallet(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*models.Payment, int64, error)
	GetByStatus(ctx context.Context, status models.PaymentStatus) ([]*models.Payment, error)
}

// GormPaymentRepository implements PaymentRepository using GORM
type GormPaymentRepository struct {
	db *gorm.DB
}

// NewGormPaymentRepository creates a new GORM payment repository
func NewGormPaymentRepository(db *gorm.DB) PaymentRepository {
	return &GormPaymentRepository{
		db: db,
	}
}

// Create creates a new payment record
func (r *GormPaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetByID retrieves a payment by ID
func (r *GormPaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// Update updates an existing payment
func (r *GormPaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// Delete removes a payment (soft delete)
func (r *GormPaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Payment{}).Error
}

// GetByWallet retrieves payments for a specific wallet with pagination
func (r *GormPaymentRepository) GetByWallet(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*models.Payment, int64, error) {
	var payments []*models.Payment
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&models.Payment{}).Where("wallet_id = ?", walletID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&payments).Error

	return payments, total, err
}

// GetByStatus retrieves payments by status
func (r *GormPaymentRepository) GetByStatus(ctx context.Context, status models.PaymentStatus) ([]*models.Payment, error) {
	var payments []*models.Payment
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&payments).Error
	return payments, err
}