package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
)

type CartRepository interface {
	Create(ctx context.Context, cart *models.Cart) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	GetBySessionID(ctx context.Context, sessionID string) (*models.Cart, error)
	Update(ctx context.Context, cart *models.Cart) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetExpiredCarts(ctx context.Context, cutoffTime time.Time) ([]*models.Cart, error)
	GetAbandonedCarts(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Cart, int64, error)
}

type CartAbandonmentRepository interface {
	Create(ctx context.Context, abandonment *models.CartAbandonmentTracking) error
	GetByCartID(ctx context.Context, cartID uuid.UUID) (*models.CartAbandonmentTracking, error)
	Update(ctx context.Context, abandonment *models.CartAbandonmentTracking) error
	List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.CartAbandonmentTracking, int64, error)
	GetUnrecoveredAbandoned(ctx context.Context, olderThan time.Time) ([]*models.CartAbandonmentTracking, error)
}

type cartRepository struct {
	db *gorm.DB
}

type cartAbandonmentRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func NewCartAbandonmentRepository(db *gorm.DB) CartAbandonmentRepository {
	return &cartAbandonmentRepository{db: db}
}

// Cart Repository Implementation
func (r *cartRepository) Create(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *cartRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).First(&cart, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, models.CartStatusActive).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetBySessionID(ctx context.Context, sessionID string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).Where("session_id = ? AND status = ?", sessionID, models.CartStatusActive).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Update(ctx context.Context, cart *models.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

func (r *cartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Cart{}, "id = ?", id).Error
}

func (r *cartRepository) GetExpiredCarts(ctx context.Context, cutoffTime time.Time) ([]*models.Cart, error) {
	var carts []*models.Cart
	err := r.db.WithContext(ctx).Where("expires_at < ? OR (user_id IS NULL AND updated_at < ?)",
		time.Now(), cutoffTime).Find(&carts).Error
	return carts, err
}

func (r *cartRepository) GetAbandonedCarts(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Cart, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Cart{}).Where("status = ?", models.CartStatusAbandoned)

	// Apply filters
	for key, value := range filters {
		switch key {
		case "min_value":
			query = query.Where("total_amount >= ?", value)
		case "date_from":
			query = query.Where("updated_at >= ?", value)
		case "date_to":
			query = query.Where("updated_at <= ?", value)
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var carts []*models.Cart
	err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&carts).Error
	return carts, total, err
}

// Cart Abandonment Repository Implementation
func (r *cartAbandonmentRepository) Create(ctx context.Context, abandonment *models.CartAbandonmentTracking) error {
	return r.db.WithContext(ctx).Create(abandonment).Error
}

func (r *cartAbandonmentRepository) GetByCartID(ctx context.Context, cartID uuid.UUID) (*models.CartAbandonmentTracking, error) {
	var abandonment models.CartAbandonmentTracking
	err := r.db.WithContext(ctx).Where("cart_id = ?", cartID).First(&abandonment).Error
	if err != nil {
		return nil, err
	}
	return &abandonment, nil
}

func (r *cartAbandonmentRepository) Update(ctx context.Context, abandonment *models.CartAbandonmentTracking) error {
	return r.db.WithContext(ctx).Save(abandonment).Error
}

func (r *cartAbandonmentRepository) List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.CartAbandonmentTracking, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.CartAbandonmentTracking{})

	// Apply filters
	for key, value := range filters {
		switch key {
		case "stage":
			if value != nil {
				query = query.Where("abandonment_stage = ?", value)
			}
		case "date_from":
			if value != nil {
				query = query.Where("abandoned_at >= ?", value)
			}
		case "date_to":
			if value != nil {
				query = query.Where("abandoned_at <= ?", value)
			}
		case "is_recovered":
			if value != nil {
				query = query.Where("is_recovered = ?", value)
			}
		case "emails_sent_min":
			if value != nil {
				query = query.Where("emails_sent >= ?", value)
			}
		case "emails_sent_max":
			if value != nil {
				query = query.Where("emails_sent <= ?", value)
			}
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var abandonments []*models.CartAbandonmentTracking
	err := query.Order("abandoned_at DESC").Offset(offset).Limit(limit).Find(&abandonments).Error
	return abandonments, total, err
}

func (r *cartAbandonmentRepository) GetUnrecoveredAbandoned(ctx context.Context, olderThan time.Time) ([]*models.CartAbandonmentTracking, error) {
	var abandonments []*models.CartAbandonmentTracking
	err := r.db.WithContext(ctx).Where("is_recovered = ? AND abandoned_at < ?", false, olderThan).Find(&abandonments).Error
	return abandonments, err
}