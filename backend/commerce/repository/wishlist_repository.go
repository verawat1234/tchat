package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
)

type WishlistRepository interface {
	Create(ctx context.Context, wishlist *models.Wishlist) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Wishlist, error)
	GetByShareToken(ctx context.Context, shareToken string) (*models.Wishlist, error)
	Update(ctx context.Context, wishlist *models.Wishlist) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Wishlist, int64, error)
	GetDefaultByUserID(ctx context.Context, userID uuid.UUID) (*models.Wishlist, error)
}

type ProductFollowRepository interface {
	Create(ctx context.Context, follow *models.ProductFollow) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ProductFollow, error)
	Update(ctx context.Context, follow *models.ProductFollow) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.ProductFollow, int64, error)
	GetByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (*models.ProductFollow, error)
}

type WishlistShareRepository interface {
	Create(ctx context.Context, share *models.WishlistShare) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.WishlistShare, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByWishlistAndUser(ctx context.Context, wishlistID, userID uuid.UUID) (*models.WishlistShare, error)
	GetSharedWithUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.WishlistShare, int64, error)
}

type wishlistRepository struct {
	db *gorm.DB
}

type productFollowRepository struct {
	db *gorm.DB
}

type wishlistShareRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{db: db}
}

func NewProductFollowRepository(db *gorm.DB) ProductFollowRepository {
	return &productFollowRepository{db: db}
}

func NewWishlistShareRepository(db *gorm.DB) WishlistShareRepository {
	return &wishlistShareRepository{db: db}
}

// Wishlist Repository Implementation
func (r *wishlistRepository) Create(ctx context.Context, wishlist *models.Wishlist) error {
	return r.db.WithContext(ctx).Create(wishlist).Error
}

func (r *wishlistRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.db.WithContext(ctx).First(&wishlist, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

func (r *wishlistRepository) GetByShareToken(ctx context.Context, shareToken string) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.db.WithContext(ctx).Where("share_token = ?", shareToken).First(&wishlist).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

func (r *wishlistRepository) Update(ctx context.Context, wishlist *models.Wishlist) error {
	return r.db.WithContext(ctx).Save(wishlist).Error
}

func (r *wishlistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Wishlist{}, "id = ?", id).Error
}

func (r *wishlistRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Wishlist, int64, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Count total
	var total int64
	if err := query.Model(&models.Wishlist{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var wishlists []*models.Wishlist
	err := query.Order("is_default DESC, created_at DESC").Offset(offset).Limit(limit).Find(&wishlists).Error
	return wishlists, total, err
}

func (r *wishlistRepository) GetDefaultByUserID(ctx context.Context, userID uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_default = ?", userID, true).First(&wishlist).Error
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

// Product Follow Repository Implementation
func (r *productFollowRepository) Create(ctx context.Context, follow *models.ProductFollow) error {
	return r.db.WithContext(ctx).Create(follow).Error
}

func (r *productFollowRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ProductFollow, error) {
	var follow models.ProductFollow
	err := r.db.WithContext(ctx).First(&follow, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

func (r *productFollowRepository) Update(ctx context.Context, follow *models.ProductFollow) error {
	return r.db.WithContext(ctx).Save(follow).Error
}

func (r *productFollowRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ProductFollow{}, "id = ?", id).Error
}

func (r *productFollowRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.ProductFollow, int64, error) {
	query := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true)

	// Count total
	var total int64
	if err := query.Model(&models.ProductFollow{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var follows []*models.ProductFollow
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&follows).Error
	return follows, total, err
}

func (r *productFollowRepository) GetByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (*models.ProductFollow, error) {
	var follow models.ProductFollow
	err := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).First(&follow).Error
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

// Wishlist Share Repository Implementation
func (r *wishlistShareRepository) Create(ctx context.Context, share *models.WishlistShare) error {
	return r.db.WithContext(ctx).Create(share).Error
}

func (r *wishlistShareRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.WishlistShare, error) {
	var share models.WishlistShare
	err := r.db.WithContext(ctx).First(&share, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *wishlistShareRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.WishlistShare{}, "id = ?", id).Error
}

func (r *wishlistShareRepository) GetByWishlistAndUser(ctx context.Context, wishlistID, userID uuid.UUID) (*models.WishlistShare, error) {
	var share models.WishlistShare
	err := r.db.WithContext(ctx).Where("wishlist_id = ? AND shared_with = ?", wishlistID, userID).First(&share).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *wishlistShareRepository) GetSharedWithUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.WishlistShare, int64, error) {
	query := r.db.WithContext(ctx).Where("shared_with = ?", userID)

	// Count total
	var total int64
	if err := query.Model(&models.WishlistShare{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	var shares []*models.WishlistShare
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&shares).Error
	return shares, total, err
}