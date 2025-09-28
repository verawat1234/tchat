package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	sharedModels "tchat.dev/shared/models"
)

type PostgreSQLTransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *PostgreSQLTransactionRepository {
	return &PostgreSQLTransactionRepository{DB: db}
}

func (r *PostgreSQLTransactionRepository) Create(ctx context.Context, tx *sharedModels.Transaction) error {
	return r.DB.WithContext(ctx).Create(tx).Error
}

func (r *PostgreSQLTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*sharedModels.Transaction, error) {
	var transaction sharedModels.Transaction
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *PostgreSQLTransactionRepository) GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*sharedModels.Transaction, error) {
	var transactions []*sharedModels.Transaction
	err := r.DB.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *PostgreSQLTransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*sharedModels.Transaction, error) {
	var transactions []*sharedModels.Transaction
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *PostgreSQLTransactionRepository) Update(ctx context.Context, tx *sharedModels.Transaction) error {
	return r.DB.WithContext(ctx).Save(tx).Error
}