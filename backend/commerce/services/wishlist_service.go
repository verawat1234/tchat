package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
	sharedModels "tchat.dev/shared/models"
)

type WishlistService interface {
	CreateWishlist(ctx context.Context, req models.CreateWishlistRequest, userID uuid.UUID) (*models.Wishlist, error)
	GetWishlist(ctx context.Context, id uuid.UUID) (*models.Wishlist, error)
	GetWishlistByShareToken(ctx context.Context, shareToken string) (*models.Wishlist, error)
	UpdateWishlist(ctx context.Context, id uuid.UUID, req models.UpdateWishlistRequest) (*models.Wishlist, error)
	DeleteWishlist(ctx context.Context, id uuid.UUID) error
	ListUserWishlists(ctx context.Context, userID uuid.UUID, pagination models.Pagination) (*models.WishlistResponse, error)
	GetDefaultWishlist(ctx context.Context, userID uuid.UUID) (*models.Wishlist, error)

	AddToWishlist(ctx context.Context, wishlistID uuid.UUID, req models.AddToWishlistRequest) error
	RemoveFromWishlist(ctx context.Context, wishlistID, productID uuid.UUID) error
	UpdateWishlistItem(ctx context.Context, wishlistID, productID uuid.UUID, quantity int, note string, priority int) error
	MoveToWishlist(ctx context.Context, fromWishlistID, toWishlistID, productID uuid.UUID) error

	ShareWishlist(ctx context.Context, wishlistID uuid.UUID, sharedWithUserID uuid.UUID, permission string) error
	UnshareWishlist(ctx context.Context, wishlistID uuid.UUID, sharedWithUserID uuid.UUID) error
	GetSharedWishlists(ctx context.Context, userID uuid.UUID, pagination models.Pagination) (*models.WishlistResponse, error)

	FollowProduct(ctx context.Context, userID, productID, businessID uuid.UUID, preferences ProductFollowPreferences) error
	UnfollowProduct(ctx context.Context, userID, productID uuid.UUID) error
	ListFollowedProducts(ctx context.Context, userID uuid.UUID, pagination models.Pagination) ([]models.ProductFollow, int64, error)
	UpdateFollowPreferences(ctx context.Context, userID, productID uuid.UUID, preferences ProductFollowPreferences) error
}

type ProductFollowPreferences struct {
	NotifyPriceChange bool `json:"notifyPriceChange"`
	NotifyBackInStock bool `json:"notifyBackInStock"`
	NotifyPromotion   bool `json:"notifyPromotion"`
	NotifyNewVariant  bool `json:"notifyNewVariant"`
}

type wishlistService struct {
	wishlistRepo      repository.WishlistRepository
	productFollowRepo repository.ProductFollowRepository
	wishlistShareRepo repository.WishlistShareRepository
	db                *gorm.DB
}

func NewWishlistService(wishlistRepo repository.WishlistRepository, productFollowRepo repository.ProductFollowRepository, wishlistShareRepo repository.WishlistShareRepository, db *gorm.DB) WishlistService {
	return &wishlistService{
		wishlistRepo:      wishlistRepo,
		productFollowRepo: productFollowRepo,
		wishlistShareRepo: wishlistShareRepo,
		db:                db,
	}
}

func (s *wishlistService) CreateWishlist(ctx context.Context, req models.CreateWishlistRequest, userID uuid.UUID) (*models.Wishlist, error) {
	// Check if this is the first wishlist for the user
	var count int64
	if err := s.db.WithContext(ctx).Model(&models.Wishlist{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to count user wishlists: %w", err)
	}

	wishlist := &models.Wishlist{
		UserID:      userID,
		Type:        req.Type,
		Privacy:     req.Privacy,
		Name:        req.Name,
		Description: req.Description,
		IsDefault:   count == 0 || req.Type == models.WishlistTypeDefault, // First wishlist or explicit default
		Items:       []models.WishlistItem{},
		ItemCount:   0,
	}

	// Generate share token if not private
	if req.Privacy != models.WishlistPrivacyPrivate {
		shareToken, err := s.generateShareToken()
		if err != nil {
			return nil, fmt.Errorf("failed to generate share token: %w", err)
		}
		wishlist.ShareToken = shareToken
	}

	if err := s.db.WithContext(ctx).Create(wishlist).Error; err != nil {
		return nil, fmt.Errorf("failed to create wishlist: %w", err)
	}

	return wishlist, nil
}

func (s *wishlistService) GetWishlist(ctx context.Context, id uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	if err := s.db.WithContext(ctx).First(&wishlist, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wishlist not found")
		}
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	return &wishlist, nil
}

func (s *wishlistService) GetWishlistByShareToken(ctx context.Context, shareToken string) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	if err := s.db.WithContext(ctx).Where("share_token = ?", shareToken).First(&wishlist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wishlist not found")
		}
		return nil, fmt.Errorf("failed to get wishlist by share token: %w", err)
	}

	return &wishlist, nil
}

func (s *wishlistService) UpdateWishlist(ctx context.Context, id uuid.UUID, req models.UpdateWishlistRequest) (*models.Wishlist, error) {
	wishlist, err := s.GetWishlist(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Privacy != nil {
		updates["privacy"] = *req.Privacy

		// Generate or remove share token based on privacy
		if *req.Privacy == models.WishlistPrivacyPrivate {
			updates["share_token"] = ""
		} else if wishlist.ShareToken == "" {
			shareToken, err := s.generateShareToken()
			if err != nil {
				return nil, fmt.Errorf("failed to generate share token: %w", err)
			}
			updates["share_token"] = shareToken
		}
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(wishlist).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update wishlist: %w", err)
		}
	}

	return wishlist, nil
}

func (s *wishlistService) DeleteWishlist(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).Delete(&models.Wishlist{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete wishlist: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("wishlist not found")
	}
	return nil
}

func (s *wishlistService) ListUserWishlists(ctx context.Context, userID uuid.UUID, pagination models.Pagination) (*models.WishlistResponse, error) {
	var total int64
	if err := s.db.WithContext(ctx).Model(&models.Wishlist{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count wishlists: %w", err)
	}

	var wishlists []*models.Wishlist
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Offset(offset).Limit(pagination.PageSize).Find(&wishlists).Error; err != nil {
		return nil, fmt.Errorf("failed to list wishlists: %w", err)
	}

	totalPages := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	return &models.WishlistResponse{
		Wishlists:  wishlists,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *wishlistService) GetDefaultWishlist(ctx context.Context, userID uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := s.db.WithContext(ctx).Where("user_id = ? AND is_default = ?", userID, true).First(&wishlist).Error

	if err == gorm.ErrRecordNotFound {
		// Create default wishlist if none exists
		req := models.CreateWishlistRequest{
			Name:        "My Wishlist",
			Description: "Default wishlist",
			Type:        models.WishlistTypeDefault,
			Privacy:     models.WishlistPrivacyPrivate,
		}
		return s.CreateWishlist(ctx, req, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get default wishlist: %w", err)
	}

	return &wishlist, nil
}

func (s *wishlistService) AddToWishlist(ctx context.Context, wishlistID uuid.UUID, req models.AddToWishlistRequest) error {
	wishlist, err := s.GetWishlist(ctx, wishlistID)
	if err != nil {
		return err
	}

	// Check if product already exists in wishlist
	for _, item := range wishlist.Items {
		if item.ProductID == req.ProductID {
			if req.VariantID == nil && item.VariantID == nil {
				return fmt.Errorf("product already in wishlist")
			}
			if req.VariantID != nil && item.VariantID != nil && *req.VariantID == *item.VariantID {
				return fmt.Errorf("product variant already in wishlist")
			}
		}
	}

	// Create new wishlist item
	newItem := models.WishlistItem{
		ProductID: req.ProductID,
		VariantID: req.VariantID,
		Quantity:  req.Quantity,
		Note:      req.Note,
		Priority:  req.Priority,
	}

	// Add to items array
	wishlist.Items = append(wishlist.Items, newItem)
	wishlist.ItemCount = len(wishlist.Items)

	if err := s.db.WithContext(ctx).Model(wishlist).Updates(map[string]interface{}{
		"items":      wishlist.Items,
		"item_count": wishlist.ItemCount,
	}).Error; err != nil {
		return fmt.Errorf("failed to add item to wishlist: %w", err)
	}

	// Create event
	eventData := map[string]interface{}{
		"wishlist_id": wishlistID,
		"product_id":  req.ProductID,
		"variant_id":  req.VariantID,
	}

	event := &sharedModels.Event{
		Type:          sharedModels.EventTypeWishlistItemAdded,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       "Wishlist Item Added",
		Description:   "An item has been added to the wishlist",
		AggregateType: "wishlist",
		AggregateID:   wishlistID.String(),
	}

	if err := event.MarshalData(eventData); err != nil {
		fmt.Printf("Failed to marshal event data: %v\n", err)
	}

	if err := s.db.WithContext(ctx).Create(event).Error; err != nil {
		fmt.Printf("Failed to create wishlist event: %v\n", err)
	}

	return nil
}

func (s *wishlistService) RemoveFromWishlist(ctx context.Context, wishlistID, productID uuid.UUID) error {
	wishlist, err := s.GetWishlist(ctx, wishlistID)
	if err != nil {
		return err
	}

	// Remove item from slice
	var updatedItems []models.WishlistItem
	found := false
	for _, item := range wishlist.Items {
		if item.ProductID != productID {
			updatedItems = append(updatedItems, item)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("product not found in wishlist")
	}

	// Update wishlist
	if err := s.db.WithContext(ctx).Model(wishlist).Updates(map[string]interface{}{
		"items":      updatedItems,
		"item_count": len(updatedItems),
	}).Error; err != nil {
		return fmt.Errorf("failed to remove item from wishlist: %w", err)
	}

	return nil
}

func (s *wishlistService) UpdateWishlistItem(ctx context.Context, wishlistID, productID uuid.UUID, quantity int, note string, priority int) error {
	wishlist, err := s.GetWishlist(ctx, wishlistID)
	if err != nil {
		return err
	}

	// Find and update item
	found := false
	for i, item := range wishlist.Items {
		if item.ProductID == productID {
			wishlist.Items[i].Quantity = quantity
			wishlist.Items[i].Note = note
			wishlist.Items[i].Priority = priority
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("product not found in wishlist")
	}

	if err := s.db.WithContext(ctx).Model(wishlist).Update("items", wishlist.Items).Error; err != nil {
		return fmt.Errorf("failed to update wishlist item: %w", err)
	}

	return nil
}

func (s *wishlistService) MoveToWishlist(ctx context.Context, fromWishlistID, toWishlistID, productID uuid.UUID) error {
	// Get both wishlists
	fromWishlist, err := s.GetWishlist(ctx, fromWishlistID)
	if err != nil {
		return err
	}

	toWishlist, err := s.GetWishlist(ctx, toWishlistID)
	if err != nil {
		return err
	}

	// Find item in source wishlist
	var itemToMove models.WishlistItem
	var updatedFromItems []models.WishlistItem
	found := false

	for _, item := range fromWishlist.Items {
		if item.ProductID == productID {
			itemToMove = item
			found = true
		} else {
			updatedFromItems = append(updatedFromItems, item)
		}
	}

	if !found {
		return fmt.Errorf("product not found in source wishlist")
	}

	// Add to destination wishlist
	toWishlist.Items = append(toWishlist.Items, itemToMove)

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()

	// Update source wishlist
	if err := tx.Model(fromWishlist).Updates(map[string]interface{}{
		"items":      updatedFromItems,
		"item_count": len(updatedFromItems),
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update source wishlist: %w", err)
	}

	// Update destination wishlist
	if err := tx.Model(toWishlist).Updates(map[string]interface{}{
		"items":      toWishlist.Items,
		"item_count": len(toWishlist.Items),
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update destination wishlist: %w", err)
	}

	return tx.Commit().Error
}

func (s *wishlistService) ShareWishlist(ctx context.Context, wishlistID uuid.UUID, sharedWithUserID uuid.UUID, permission string) error {
	share := &models.WishlistShare{
		WishlistID: wishlistID,
		SharedBy:   uuid.Nil, // This should be set from context
		SharedWith: sharedWithUserID,
		Permission: permission,
	}

	if err := s.db.WithContext(ctx).Create(share).Error; err != nil {
		return fmt.Errorf("failed to share wishlist: %w", err)
	}

	// Update share count
	if err := s.db.WithContext(ctx).Model(&models.Wishlist{}).Where("id = ?", wishlistID).
		UpdateColumn("share_count", gorm.Expr("share_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to update share count: %w", err)
	}

	return nil
}

func (s *wishlistService) UnshareWishlist(ctx context.Context, wishlistID uuid.UUID, sharedWithUserID uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("wishlist_id = ? AND shared_with = ?", wishlistID, sharedWithUserID).
		Delete(&models.WishlistShare{})

	if result.Error != nil {
		return fmt.Errorf("failed to unshare wishlist: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		// Update share count
		if err := s.db.WithContext(ctx).Model(&models.Wishlist{}).Where("id = ?", wishlistID).
			UpdateColumn("share_count", gorm.Expr("share_count - 1")).Error; err != nil {
			return fmt.Errorf("failed to update share count: %w", err)
		}
	}

	return nil
}

func (s *wishlistService) GetSharedWishlists(ctx context.Context, userID uuid.UUID, pagination models.Pagination) (*models.WishlistResponse, error) {
	var total int64
	shareQuery := s.db.WithContext(ctx).Model(&models.WishlistShare{}).Where("shared_with = ?", userID)
	if err := shareQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count shared wishlists: %w", err)
	}

	var wishlists []*models.Wishlist
	offset := (pagination.Page - 1) * pagination.PageSize

	if err := s.db.WithContext(ctx).
		Joins("JOIN wishlist_shares ON wishlists.id = wishlist_shares.wishlist_id").
		Where("wishlist_shares.shared_with = ?", userID).
		Order("wishlist_shares.created_at DESC").
		Offset(offset).Limit(pagination.PageSize).
		Find(&wishlists).Error; err != nil {
		return nil, fmt.Errorf("failed to list shared wishlists: %w", err)
	}

	totalPages := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	return &models.WishlistResponse{
		Wishlists:  wishlists,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *wishlistService) FollowProduct(ctx context.Context, userID, productID, businessID uuid.UUID, preferences ProductFollowPreferences) error {
	// Check if already following
	var existing models.ProductFollow
	err := s.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).First(&existing).Error

	if err == nil {
		return fmt.Errorf("already following this product")
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing follow: %w", err)
	}

	follow := &models.ProductFollow{
		UserID:              userID,
		ProductID:           productID,
		BusinessID:          businessID,
		NotifyPriceChange:   preferences.NotifyPriceChange,
		NotifyBackInStock:   preferences.NotifyBackInStock,
		NotifyPromotion:     preferences.NotifyPromotion,
		NotifyNewVariant:    preferences.NotifyNewVariant,
		IsActive:            true,
	}

	if err := s.db.WithContext(ctx).Create(follow).Error; err != nil {
		return fmt.Errorf("failed to follow product: %w", err)
	}

	return nil
}

func (s *wishlistService) UnfollowProduct(ctx context.Context, userID, productID uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.ProductFollow{})

	if result.Error != nil {
		return fmt.Errorf("failed to unfollow product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product follow not found")
	}

	return nil
}

func (s *wishlistService) ListFollowedProducts(ctx context.Context, userID uuid.UUID, pagination models.Pagination) ([]models.ProductFollow, int64, error) {
	var total int64
	if err := s.db.WithContext(ctx).Model(&models.ProductFollow{}).
		Where("user_id = ? AND is_active = ?", userID, true).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count followed products: %w", err)
	}

	var follows []models.ProductFollow
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := s.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Offset(offset).Limit(pagination.PageSize).Find(&follows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list followed products: %w", err)
	}

	return follows, total, nil
}

func (s *wishlistService) UpdateFollowPreferences(ctx context.Context, userID, productID uuid.UUID, preferences ProductFollowPreferences) error {
	updates := map[string]interface{}{
		"notify_price_change": preferences.NotifyPriceChange,
		"notify_back_in_stock": preferences.NotifyBackInStock,
		"notify_promotion":     preferences.NotifyPromotion,
		"notify_new_variant":   preferences.NotifyNewVariant,
	}

	result := s.db.WithContext(ctx).Model(&models.ProductFollow{}).
		Where("user_id = ? AND product_id = ?", userID, productID).Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update follow preferences: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product follow not found")
	}

	return nil
}

func (s *wishlistService) generateShareToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}