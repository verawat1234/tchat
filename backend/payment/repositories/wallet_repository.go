package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	sharedModels "tchat.dev/shared/models"
)

type PostgreSQLWalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *PostgreSQLWalletRepository {
	return &PostgreSQLWalletRepository{DB: db}
}

func (r *PostgreSQLWalletRepository) Create(ctx context.Context, wallet *sharedModels.Wallet) error {
	return r.DB.WithContext(ctx).Create(wallet).Error
}

func (r *PostgreSQLWalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*sharedModels.Wallet, error) {
	var wallet sharedModels.Wallet
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *PostgreSQLWalletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*sharedModels.Wallet, error) {
	var wallets []*sharedModels.Wallet
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&wallets).Error
	return wallets, err
}

func (r *PostgreSQLWalletRepository) GetByUserIDAndCurrency(ctx context.Context, userID uuid.UUID, currency string) (*sharedModels.Wallet, error) {
	var wallet sharedModels.Wallet
	err := r.DB.WithContext(ctx).Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *PostgreSQLWalletRepository) Update(ctx context.Context, wallet *sharedModels.Wallet) error {
	return r.DB.WithContext(ctx).Save(wallet).Error
}

func (r *PostgreSQLWalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal, txID uuid.UUID) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet sharedModels.Wallet
		if err := tx.Where("id = ?", walletID).First(&wallet).Error; err != nil {
			return err
		}

		// Update the default currency balance (assuming USD if not specified)
		currency := wallet.Settings.DefaultCurrency
		if currency == "" {
			currency = "USD"
		}

		// Add balance using the wallet's built-in method
		if amount.IsPositive() {
			if err := wallet.AddBalance(currency, amount, "available"); err != nil {
				return err
			}
		} else {
			// For negative amounts, deduct from available balance
			if err := wallet.DeductBalance(currency, amount.Abs(), "available"); err != nil {
				return err
			}
		}

		// Update activity tracking
		wallet.UpdateActivity()

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		return nil
	})
}