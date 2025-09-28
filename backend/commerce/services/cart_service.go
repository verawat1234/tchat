package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
	sharedModels "tchat.dev/shared/models"
)

type CartService interface {
	GetCart(ctx context.Context, userID *uuid.UUID, sessionID string) (*models.Cart, error)
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	GetCartBySessionID(ctx context.Context, sessionID string) (*models.Cart, error)
	CreateCart(ctx context.Context, userID *uuid.UUID, sessionID string) (*models.Cart, error)
	AddToCart(ctx context.Context, cartID uuid.UUID, req models.AddToCartRequest) error
	UpdateCartItem(ctx context.Context, cartID, productID uuid.UUID, req models.UpdateCartItemRequest) error
	RemoveFromCart(ctx context.Context, cartID, productID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
	MergeCart(ctx context.Context, guestCartID, userCartID uuid.UUID) error
	ConvertToOrder(ctx context.Context, cartID uuid.UUID, orderID uuid.UUID) error

	ApplyCoupon(ctx context.Context, cartID uuid.UUID, couponCode string) error
	RemoveCoupon(ctx context.Context, cartID uuid.UUID) error

	GetCartSummary(ctx context.Context, cartID uuid.UUID) (*CartSummary, error)
	ValidateCart(ctx context.Context, cartID uuid.UUID) (*CartValidation, error)

	SaveForLater(ctx context.Context, cartID, productID uuid.UUID) error
	MoveToCart(ctx context.Context, cartID, productID uuid.UUID) error
	GetSavedItems(ctx context.Context, cartID uuid.UUID) ([]models.CartItem, error)

	TrackAbandonment(ctx context.Context, cartID uuid.UUID, stage string, lastPage string) error
	GetAbandonedCarts(ctx context.Context, filters map[string]interface{}, pagination models.Pagination) (*models.CartResponse, error)
	CreateAbandonmentTracking(ctx context.Context, req models.CreateAbandonmentTrackingRequest) (*models.CartAbandonmentTracking, error)
	GetAbandonmentAnalytics(ctx context.Context, filters map[string]interface{}, pagination models.Pagination) (*models.AbandonmentTrackingResponse, error)
	MarkCartRecovered(ctx context.Context, cartID, orderID uuid.UUID) error

	CleanupExpiredCarts(ctx context.Context) error
}

type CartSummary struct {
	SubtotalAmount decimal.Decimal `json:"subtotalAmount"`
	TaxAmount      decimal.Decimal `json:"taxAmount"`
	ShippingAmount decimal.Decimal `json:"shippingAmount"`
	DiscountAmount decimal.Decimal `json:"discountAmount"`
	TotalAmount    decimal.Decimal `json:"totalAmount"`
	ItemCount      int             `json:"itemCount"`
	BusinessCount  int             `json:"businessCount"`
	Currency       string          `json:"currency"`
	CouponCode     string          `json:"couponCode,omitempty"`
	EstimatedTax   decimal.Decimal `json:"estimatedTax"`
}

type CartValidation struct {
	IsValid        bool                   `json:"isValid"`
	Issues         []CartValidationIssue  `json:"issues"`
	UpdatedItems   []models.CartItem      `json:"updatedItems,omitempty"`
	RemovedItems   []models.CartItem      `json:"removedItems,omitempty"`
}

type CartValidationIssue struct {
	Type        string      `json:"type"`
	ProductID   uuid.UUID   `json:"productId"`
	ProductName string      `json:"productName"`
	Message     string      `json:"message"`
	Severity    string      `json:"severity"` // warning, error
}

type AbandonmentFilters struct {
	Stage         *string    `json:"stage,omitempty"`
	MinValue      *decimal.Decimal `json:"minValue,omitempty"`
	DateFrom      *time.Time `json:"dateFrom,omitempty"`
	DateTo        *time.Time `json:"dateTo,omitempty"`
	IsRecovered   *bool      `json:"isRecovered,omitempty"`
	EmailsSentMin *int       `json:"emailsSentMin,omitempty"`
	EmailsSentMax *int       `json:"emailsSentMax,omitempty"`
}

type cartService struct {
	cartRepo            repository.CartRepository
	cartAbandonmentRepo repository.CartAbandonmentRepository
	db                  *gorm.DB
}

func NewCartService(cartRepo repository.CartRepository, cartAbandonmentRepo repository.CartAbandonmentRepository, db *gorm.DB) CartService {
	return &cartService{
		cartRepo:            cartRepo,
		cartAbandonmentRepo: cartAbandonmentRepo,
		db:                  db,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID *uuid.UUID, sessionID string) (*models.Cart, error) {
	var cart models.Cart
	query := s.db.WithContext(ctx)

	if userID != nil {
		// Try to find user cart first
		err := query.Where("user_id = ? AND status = ?", *userID, models.CartStatusActive).First(&cart).Error
		if err == nil {
			return &cart, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get user cart: %w", err)
		}
	}

	// Try to find session cart
	if sessionID != "" {
		err := query.Where("session_id = ? AND status = ?", sessionID, models.CartStatusActive).First(&cart).Error
		if err == nil {
			return &cart, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to get session cart: %w", err)
		}
	}

	// Create new cart if none found
	return s.CreateCart(ctx, userID, sessionID)
}

func (s *cartService) CreateCart(ctx context.Context, userID *uuid.UUID, sessionID string) (*models.Cart, error) {
	cart := &models.Cart{
		UserID:         userID,
		SessionID:      sessionID,
		Status:         models.CartStatusActive,
		Items:          []models.CartItem{},
		ItemCount:      0,
		BusinessCount:  0,
		SubtotalAmount: decimal.Zero,
		TaxAmount:      decimal.Zero,
		ShippingAmount: decimal.Zero,
		DiscountAmount: decimal.Zero,
		TotalAmount:    decimal.Zero,
		Currency:       "USD",
		LastActivity:   time.Now(),
	}

	// Set expiry for guest carts
	if userID == nil {
		expiryTime := time.Now().Add(7 * 24 * time.Hour) // 7 days
		cart.ExpiresAt = &expiryTime
	}

	if err := s.db.WithContext(ctx).Create(cart).Error; err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}

	return cart, nil
}

func (s *cartService) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	return s.GetCart(ctx, &userID, "")
}

func (s *cartService) GetCartBySessionID(ctx context.Context, sessionID string) (*models.Cart, error) {
	return s.GetCart(ctx, nil, sessionID)
}

func (s *cartService) AddToCart(ctx context.Context, cartID uuid.UUID, req models.AddToCartRequest) error {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Check if item already exists in cart
	itemIndex := -1
	for i, item := range cart.Items {
		if item.ProductID == req.ProductID {
			if req.VariantID == nil && item.VariantID == nil {
				itemIndex = i
				break
			}
			if req.VariantID != nil && item.VariantID != nil && *req.VariantID == *item.VariantID {
				itemIndex = i
				break
			}
		}
	}

	if itemIndex >= 0 {
		// Update existing item quantity
		cart.Items[itemIndex].Quantity += req.Quantity
		cart.Items[itemIndex].IsGift = req.IsGift
		cart.Items[itemIndex].GiftMessage = req.GiftMessage

		// Recalculate total price for the item
		cart.Items[itemIndex].TotalPrice = cart.Items[itemIndex].UnitPrice.Mul(decimal.NewFromInt(int64(cart.Items[itemIndex].Quantity)))
	} else {
		// Add new item - Note: In a real implementation, you'd fetch product details
		newItem := models.CartItem{
			ProductID:    req.ProductID,
			VariantID:    req.VariantID,
			BusinessID:   uuid.Nil, // Should be fetched from product
			Type:         models.CartItemTypeProduct,
			Quantity:     req.Quantity,
			UnitPrice:    decimal.Zero, // Should be fetched from product
			TotalPrice:   decimal.Zero, // Should be calculated
			Currency:     cart.Currency,
			ProductName:  "Product Name", // Should be fetched from product
			IsGift:       req.IsGift,
			GiftMessage:  req.GiftMessage,
			IsAvailable:  true,
		}

		cart.Items = append(cart.Items, newItem)
	}

	// Update cart totals
	s.calculateCartTotals(cart)

	// Save cart
	if err := s.db.WithContext(ctx).Model(cart).Updates(map[string]interface{}{
		"items":           cart.Items,
		"item_count":      cart.ItemCount,
		"business_count":  cart.BusinessCount,
		"subtotal_amount": cart.SubtotalAmount,
		"total_amount":    cart.TotalAmount,
		"last_activity":   time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	// Create event
	eventData := map[string]interface{}{
		"cart_id":    cartID,
		"product_id": req.ProductID,
		"quantity":   req.Quantity,
	}

	event := &sharedModels.Event{
		Type:          sharedModels.EventTypeCartItemAdded,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       "Cart Item Added",
		Description:   "An item has been added to the cart",
		AggregateType: "cart",
		AggregateID:   cartID.String(),
	}

	if err := event.MarshalData(eventData); err != nil {
		fmt.Printf("Failed to marshal event data: %v\n", err)
	}

	if err := s.db.WithContext(ctx).Create(event).Error; err != nil {
		fmt.Printf("Failed to create cart event: %v\n", err)
	}

	return nil
}

func (s *cartService) UpdateCartItem(ctx context.Context, cartID, productID uuid.UUID, req models.UpdateCartItemRequest) error {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Find item
	itemIndex := -1
	for i, item := range cart.Items {
		if item.ProductID == productID {
			itemIndex = i
			break
		}
	}

	if itemIndex == -1 {
		return fmt.Errorf("item not found in cart")
	}

	// Update item
	if req.Quantity != nil {
		if *req.Quantity <= 0 {
			return s.RemoveFromCart(ctx, cartID, productID)
		}
		cart.Items[itemIndex].Quantity = *req.Quantity
		cart.Items[itemIndex].TotalPrice = cart.Items[itemIndex].UnitPrice.Mul(decimal.NewFromInt(int64(*req.Quantity)))
	}
	if req.IsGift != nil {
		cart.Items[itemIndex].IsGift = *req.IsGift
	}
	if req.GiftMessage != nil {
		cart.Items[itemIndex].GiftMessage = *req.GiftMessage
	}

	// Update cart totals
	s.calculateCartTotals(cart)

	// Save cart
	if err := s.db.WithContext(ctx).Model(cart).Updates(map[string]interface{}{
		"items":           cart.Items,
		"item_count":      cart.ItemCount,
		"subtotal_amount": cart.SubtotalAmount,
		"total_amount":    cart.TotalAmount,
		"last_activity":   time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}

func (s *cartService) RemoveFromCart(ctx context.Context, cartID, productID uuid.UUID) error {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Remove item
	var updatedItems []models.CartItem
	found := false
	for _, item := range cart.Items {
		if item.ProductID != productID {
			updatedItems = append(updatedItems, item)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("item not found in cart")
	}

	cart.Items = updatedItems

	// Update cart totals
	s.calculateCartTotals(cart)

	// Save cart
	if err := s.db.WithContext(ctx).Model(cart).Updates(map[string]interface{}{
		"items":           cart.Items,
		"item_count":      cart.ItemCount,
		"business_count":  cart.BusinessCount,
		"subtotal_amount": cart.SubtotalAmount,
		"total_amount":    cart.TotalAmount,
		"last_activity":   time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}

func (s *cartService) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Clear all items and reset totals
	cart.Items = []models.CartItem{}
	cart.ItemCount = 0
	cart.BusinessCount = 0
	cart.SubtotalAmount = decimal.Zero
	cart.TaxAmount = decimal.Zero
	cart.ShippingAmount = decimal.Zero
	cart.DiscountAmount = decimal.Zero
	cart.TotalAmount = decimal.Zero
	cart.CouponCode = ""

	if err := s.db.WithContext(ctx).Model(cart).Updates(map[string]interface{}{
		"items":           cart.Items,
		"item_count":      cart.ItemCount,
		"business_count":  cart.BusinessCount,
		"subtotal_amount": cart.SubtotalAmount,
		"tax_amount":      cart.TaxAmount,
		"shipping_amount": cart.ShippingAmount,
		"discount_amount": cart.DiscountAmount,
		"total_amount":    cart.TotalAmount,
		"coupon_code":     cart.CouponCode,
		"last_activity":   time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

func (s *cartService) MergeCart(ctx context.Context, guestCartID, userCartID uuid.UUID) error {
	// Get both carts
	guestCart, err := s.getCartByID(ctx, guestCartID)
	if err != nil {
		return fmt.Errorf("failed to get guest cart: %w", err)
	}

	userCart, err := s.getCartByID(ctx, userCartID)
	if err != nil {
		return fmt.Errorf("failed to get user cart: %w", err)
	}

	// Merge items from guest cart to user cart
	for _, guestItem := range guestCart.Items {
		// Check if item exists in user cart
		found := false
		for i, userItem := range userCart.Items {
			if userItem.ProductID == guestItem.ProductID &&
				((userItem.VariantID == nil && guestItem.VariantID == nil) ||
					(userItem.VariantID != nil && guestItem.VariantID != nil && *userItem.VariantID == *guestItem.VariantID)) {
				// Merge quantities
				userCart.Items[i].Quantity += guestItem.Quantity
				userCart.Items[i].TotalPrice = userCart.Items[i].UnitPrice.Mul(decimal.NewFromInt(int64(userCart.Items[i].Quantity)))
				found = true
				break
			}
		}

		if !found {
			// Add guest item to user cart
			userCart.Items = append(userCart.Items, guestItem)
		}
	}

	// Update user cart totals
	s.calculateCartTotals(userCart)

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()

	// Update user cart
	if err := tx.Model(userCart).Updates(map[string]interface{}{
		"items":           userCart.Items,
		"item_count":      userCart.ItemCount,
		"business_count":  userCart.BusinessCount,
		"subtotal_amount": userCart.SubtotalAmount,
		"total_amount":    userCart.TotalAmount,
		"last_activity":   time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user cart: %w", err)
	}

	// Delete guest cart
	if err := tx.Delete(guestCart).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete guest cart: %w", err)
	}

	return tx.Commit().Error
}

func (s *cartService) ConvertToOrder(ctx context.Context, cartID uuid.UUID, orderID uuid.UUID) error {
	updates := map[string]interface{}{
		"status":                  models.CartStatusConverted,
		"converted_to_order_id":   orderID,
	}

	if err := s.db.WithContext(ctx).Model(&models.Cart{}).Where("id = ?", cartID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to convert cart to order: %w", err)
	}

	return nil
}

func (s *cartService) ApplyCoupon(ctx context.Context, cartID uuid.UUID, couponCode string) error {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Validate coupon code
	discount, err := s.validateAndCalculateCouponDiscount(ctx, couponCode, cart)
	if err != nil {
		return fmt.Errorf("invalid coupon: %w", err)
	}

	// Apply coupon and recalculate totals
	cart.CouponCode = couponCode
	cart.DiscountAmount = discount

	// Recalculate total amount
	cart.TotalAmount = cart.SubtotalAmount.Add(cart.TaxAmount).Add(cart.ShippingAmount).Sub(cart.DiscountAmount)

	if err := s.db.WithContext(ctx).Model(cart).Updates(map[string]interface{}{
		"coupon_code":     couponCode,
		"discount_amount": discount,
		"total_amount":    cart.TotalAmount,
		"last_activity":   time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to apply coupon: %w", err)
	}

	return nil
}

func (s *cartService) RemoveCoupon(ctx context.Context, cartID uuid.UUID) error {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return err
	}

	cart.CouponCode = ""
	cart.DiscountAmount = decimal.Zero

	// Recalculate totals
	s.calculateCartTotals(cart)

	if err := s.db.WithContext(ctx).Model(cart).Updates(map[string]interface{}{
		"coupon_code":     "",
		"discount_amount": cart.DiscountAmount,
		"total_amount":    cart.TotalAmount,
		"last_activity":   time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to remove coupon: %w", err)
	}

	return nil
}

func (s *cartService) GetCartSummary(ctx context.Context, cartID uuid.UUID) (*CartSummary, error) {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	return &CartSummary{
		SubtotalAmount: cart.SubtotalAmount,
		TaxAmount:      cart.TaxAmount,
		ShippingAmount: cart.ShippingAmount,
		DiscountAmount: cart.DiscountAmount,
		TotalAmount:    cart.TotalAmount,
		ItemCount:      cart.ItemCount,
		BusinessCount:  cart.BusinessCount,
		Currency:       cart.Currency,
		CouponCode:     cart.CouponCode,
		EstimatedTax:   s.calculateTaxAmount(cart), // Tax based on shipping address
	}, nil
}

func (s *cartService) ValidateCart(ctx context.Context, cartID uuid.UUID) (*CartValidation, error) {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	validation := &CartValidation{
		IsValid: true,
		Issues:  []CartValidationIssue{},
	}

	// Validate cart is not empty
	if len(cart.Items) == 0 {
		validation.Issues = append(validation.Issues, CartValidationIssue{
			Type:     "empty_cart",
			Message:  "Cart is empty",
			Severity: "error",
		})
		validation.IsValid = false
	}

	// Validate each cart item
	for _, item := range cart.Items {
		// Check product availability
		if err := s.validateProductAvailability(ctx, item.ProductID); err != nil {
			validation.Issues = append(validation.Issues, CartValidationIssue{
				Type:      "product_unavailable",
				Message:   fmt.Sprintf("Product %s is no longer available", item.ProductID),
				Severity:  "error",
				ProductID: item.ProductID,
			})
			validation.IsValid = false
		}

		// Check stock levels
		if err := s.validateStockLevel(ctx, item.ProductID, item.Quantity); err != nil {
			validation.Issues = append(validation.Issues, CartValidationIssue{
				Type:      "insufficient_stock",
				Message:   fmt.Sprintf("Insufficient stock for product %s", item.ProductID),
				Severity:  "error",
				ProductID: item.ProductID,
			})
			validation.IsValid = false
		}

		// Validate prices
		if err := s.validateProductPrice(ctx, item); err != nil {
			validation.Issues = append(validation.Issues, CartValidationIssue{
				Type:      "price_changed",
				Message:   fmt.Sprintf("Price has changed for product %s", item.ProductID),
				Severity:  "warning",
				ProductID: item.ProductID,
			})
			// Price changes are warnings, not errors
		}
	}

	// Check shipping restrictions
	if err := s.validateShippingRestrictions(ctx, cart); err != nil {
		validation.Issues = append(validation.Issues, CartValidationIssue{
			Type:     "shipping_restriction",
			Message:  err.Error(),
			Severity: "error",
		})
		validation.IsValid = false
	}

	// Check cart total minimum
	if cart.SubtotalAmount.LessThan(decimal.NewFromFloat(1.0)) {
		validation.Issues = append(validation.Issues, CartValidationIssue{
			Type:     "minimum_order",
			Message:  "Order total must be at least $1.00",
			Severity: "error",
		})
		validation.IsValid = false
	}

	return validation, nil
}

func (s *cartService) SaveForLater(ctx context.Context, cartID, productID uuid.UUID) error {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Find and update item
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].IsSavedForLater = true
			break
		}
	}

	if err := s.db.WithContext(ctx).Model(cart).Update("items", cart.Items).Error; err != nil {
		return fmt.Errorf("failed to save item for later: %w", err)
	}

	return nil
}

func (s *cartService) MoveToCart(ctx context.Context, cartID, productID uuid.UUID) error {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return err
	}

	// Find and update item
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].IsSavedForLater = false
			break
		}
	}

	if err := s.db.WithContext(ctx).Model(cart).Update("items", cart.Items).Error; err != nil {
		return fmt.Errorf("failed to move item to cart: %w", err)
	}

	return nil
}

func (s *cartService) GetSavedItems(ctx context.Context, cartID uuid.UUID) ([]models.CartItem, error) {
	cart, err := s.getCartByID(ctx, cartID)
	if err != nil {
		return nil, err
	}

	var savedItems []models.CartItem
	for _, item := range cart.Items {
		if item.IsSavedForLater {
			savedItems = append(savedItems, item)
		}
	}

	return savedItems, nil
}

func (s *cartService) TrackAbandonment(ctx context.Context, cartID uuid.UUID, stage string, lastPage string) error {
	// Check if abandonment already tracked
	var existing models.CartAbandonmentTracking
	err := s.db.WithContext(ctx).Where("cart_id = ?", cartID).First(&existing).Error

	if err == nil {
		// Update existing record
		updates := map[string]interface{}{
			"abandonment_stage":  stage,
			"last_page_visited":  lastPage,
		}
		return s.db.WithContext(ctx).Model(&existing).Updates(updates).Error
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing abandonment: %w", err)
	}

	// Create new abandonment record
	abandonment := &models.CartAbandonmentTracking{
		CartID:            cartID,
		AbandonedAt:       time.Now(),
		AbandonmentStage:  stage,
		LastPageVisited:   lastPage,
	}

	if err := s.db.WithContext(ctx).Create(abandonment).Error; err != nil {
		return fmt.Errorf("failed to track cart abandonment: %w", err)
	}

	return nil
}


func (s *cartService) MarkCartRecovered(ctx context.Context, cartID, orderID uuid.UUID) error {
	updates := map[string]interface{}{
		"is_recovered":         true,
		"recovered_at":         time.Now(),
		"recovered_order_id":   orderID,
	}

	if err := s.db.WithContext(ctx).Model(&models.CartAbandonmentTracking{}).
		Where("cart_id = ?", cartID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to mark cart as recovered: %w", err)
	}

	return nil
}

func (s *cartService) CleanupExpiredCarts(ctx context.Context) error {
	// Delete expired guest carts
	cutoff := time.Now().Add(-30 * 24 * time.Hour) // 30 days old

	result := s.db.WithContext(ctx).Where("expires_at < ? OR (user_id IS NULL AND updated_at < ?)",
		time.Now(), cutoff).Delete(&models.Cart{})

	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired carts: %w", result.Error)
	}

	fmt.Printf("Cleaned up %d expired carts\n", result.RowsAffected)
	return nil
}

// Helper methods

func (s *cartService) getCartByID(ctx context.Context, cartID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	if err := s.db.WithContext(ctx).First(&cart, "id = ?", cartID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	return &cart, nil
}

func (s *cartService) calculateCartTotals(cart *models.Cart) {
	cart.SubtotalAmount = decimal.Zero
	cart.ItemCount = 0
	businessMap := make(map[uuid.UUID]bool)

	// Calculate subtotal and count items
	for _, item := range cart.Items {
		if !item.IsSavedForLater {
			cart.SubtotalAmount = cart.SubtotalAmount.Add(item.TotalPrice)
			cart.ItemCount += item.Quantity
			businessMap[item.BusinessID] = true
		}
	}

	cart.BusinessCount = len(businessMap)

	// Calculate tax based on shipping address
	cart.TaxAmount = s.calculateTaxAmount(cart)

	// Calculate shipping based on cart contents and destination
	cart.ShippingAmount = s.calculateShippingAmount(cart)

	// Calculate total
	cart.TotalAmount = cart.SubtotalAmount.Add(cart.TaxAmount).Add(cart.ShippingAmount).Sub(cart.DiscountAmount)
}

// Coupon validation and discount calculation
func (s *cartService) validateAndCalculateCouponDiscount(ctx context.Context, couponCode string, cart *models.Cart) (decimal.Decimal, error) {
	// Define some predefined coupon codes with their rules
	coupons := map[string]struct {
		DiscountType string  // "percentage" or "fixed"
		Value        float64 // percentage (0.1 = 10%) or fixed amount
		MinOrder     float64 // minimum order amount
		MaxDiscount  float64 // maximum discount amount for percentage
		IsActive     bool
	}{
		"SAVE10": {
			DiscountType: "percentage",
			Value:        0.10, // 10%
			MinOrder:     50.0,
			MaxDiscount:  20.0,
			IsActive:     true,
		},
		"SAVE20": {
			DiscountType: "percentage",
			Value:        0.20, // 20%
			MinOrder:     100.0,
			MaxDiscount:  50.0,
			IsActive:     true,
		},
		"FLAT15": {
			DiscountType: "fixed",
			Value:        15.0,
			MinOrder:     75.0,
			MaxDiscount:  0, // Not applicable for fixed
			IsActive:     true,
		},
		"WELCOME5": {
			DiscountType: "fixed",
			Value:        5.0,
			MinOrder:     25.0,
			MaxDiscount:  0,
			IsActive:     true,
		},
	}

	coupon, exists := coupons[couponCode]
	if !exists {
		return decimal.Zero, fmt.Errorf("coupon code '%s' not found", couponCode)
	}

	if !coupon.IsActive {
		return decimal.Zero, fmt.Errorf("coupon code '%s' is expired or inactive", couponCode)
	}

	subtotal, _ := cart.SubtotalAmount.Float64()
	if subtotal < coupon.MinOrder {
		return decimal.Zero, fmt.Errorf("minimum order amount of $%.2f required for coupon '%s'", coupon.MinOrder, couponCode)
	}

	var discount decimal.Decimal

	switch coupon.DiscountType {
	case "percentage":
		// Calculate percentage discount
		discountAmount := subtotal * coupon.Value

		// Apply maximum discount limit
		if coupon.MaxDiscount > 0 && discountAmount > coupon.MaxDiscount {
			discountAmount = coupon.MaxDiscount
		}

		discount = decimal.NewFromFloat(discountAmount)

	case "fixed":
		// Fixed amount discount
		discount = decimal.NewFromFloat(coupon.Value)

		// Don't let discount exceed subtotal
		if discount.GreaterThan(cart.SubtotalAmount) {
			discount = cart.SubtotalAmount
		}

	default:
		return decimal.Zero, fmt.Errorf("invalid discount type for coupon '%s'", couponCode)
	}

	return discount, nil
}

// Tax calculation based on shipping address and Southeast Asian regions
func (s *cartService) calculateTaxAmount(cart *models.Cart) decimal.Decimal {
	if cart.SubtotalAmount.IsZero() {
		return decimal.Zero
	}

	// Southeast Asian tax rates by country/region
	taxRates := map[string]float64{
		"TH":   0.07,  // Thailand VAT 7%
		"SG":   0.08,  // Singapore GST 8%
		"MY":   0.06,  // Malaysia SST 6%
		"ID":   0.11,  // Indonesia VAT 11%
		"VN":   0.10,  // Vietnam VAT 10%
		"PH":   0.12,  // Philippines VAT 12%
		"US":   0.08,  // US average sales tax 8%
		"DEFAULT": 0.05, // Default tax rate 5%
	}

	// Determine tax rate based on shipping address
	var taxRate float64

	// If shipping address is available, use country-specific rate
	if cart.ShippingAddress != nil && cart.ShippingAddress.Country != "" {
		if rate, exists := taxRates[cart.ShippingAddress.Country]; exists {
			taxRate = rate
		} else {
			taxRate = taxRates["DEFAULT"]
		}
	} else {
		// Use default tax rate if no shipping address
		taxRate = taxRates["DEFAULT"]
	}

	// Calculate tax amount
	taxAmount := cart.SubtotalAmount.Mul(decimal.NewFromFloat(taxRate))

	return taxAmount
}

// Shipping calculation based on cart contents and destination
func (s *cartService) calculateShippingAmount(cart *models.Cart) decimal.Decimal {
	if cart.SubtotalAmount.IsZero() {
		return decimal.Zero
	}

	// Free shipping threshold
	freeShippingThreshold := decimal.NewFromFloat(100.0)
	if cart.SubtotalAmount.GreaterThanOrEqual(freeShippingThreshold) {
		return decimal.Zero
	}

	// Base shipping rates by country/region
	baseShippingRates := map[string]float64{
		"TH": 3.0,  // Thailand
		"SG": 5.0,  // Singapore
		"MY": 4.0,  // Malaysia
		"ID": 6.0,  // Indonesia
		"VN": 5.0,  // Vietnam
		"PH": 5.0,  // Philippines
		"US": 10.0, // US
		"DEFAULT": 8.0,
	}

	// Shipping rate based on destination
	var baseRate float64
	if cart.ShippingAddress != nil && cart.ShippingAddress.Country != "" {
		if rate, exists := baseShippingRates[cart.ShippingAddress.Country]; exists {
			baseRate = rate
		} else {
			baseRate = baseShippingRates["DEFAULT"]
		}
	} else {
		baseRate = baseShippingRates["DEFAULT"]
	}

	// Calculate shipping based on number of businesses (multiple vendors)
	multiplier := float64(cart.BusinessCount)
	if multiplier < 1 {
		multiplier = 1
	}

	// Apply business multiplier with diminishing returns
	if multiplier > 1 {
		// Each additional business adds 50% of base rate
		additionalRate := (multiplier - 1) * baseRate * 0.5
		baseRate += additionalRate
	}

	// Weight-based shipping (estimate based on item count)
	if cart.ItemCount > 5 {
		// Add extra shipping for heavy/bulky orders
		weightSurcharge := float64(cart.ItemCount-5) * 0.5
		baseRate += weightSurcharge
	}

	return decimal.NewFromFloat(baseRate)
}

// Cart validation helper methods

func (s *cartService) validateProductAvailability(ctx context.Context, productID uuid.UUID) error {
	// Check if product exists and is active
	var product struct {
		ID     uuid.UUID
		Status string
	}

	err := s.db.WithContext(ctx).
		Table("products").
		Select("id, status").
		Where("id = ?", productID).
		First(&product).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("failed to check product availability: %w", err)
	}

	if product.Status != "active" {
		return fmt.Errorf("product is not available (status: %s)", product.Status)
	}

	return nil
}

func (s *cartService) validateStockLevel(ctx context.Context, productID uuid.UUID, requiredQuantity int) error {
	// Check available stock
	var product struct {
		ID               uuid.UUID
		StockQuantity    int
		TrackInventory   bool
		AllowBackorders  bool
	}

	err := s.db.WithContext(ctx).
		Table("products").
		Select("id, stock_quantity, track_inventory, allow_backorders").
		Where("id = ?", productID).
		First(&product).Error

	if err != nil {
		return fmt.Errorf("failed to check stock level: %w", err)
	}

	// Skip validation if inventory tracking is disabled
	if !product.TrackInventory {
		return nil
	}

	// Check if sufficient stock is available
	if product.StockQuantity < requiredQuantity {
		if !product.AllowBackorders {
			return fmt.Errorf("insufficient stock: %d available, %d requested",
				product.StockQuantity, requiredQuantity)
		}
		// Allow backorders but could show warning
	}

	return nil
}

func (s *cartService) validateProductPrice(ctx context.Context, item models.CartItem) error {
	// Get current product price
	var currentPrice decimal.Decimal
	err := s.db.WithContext(ctx).
		Table("products").
		Select("price").
		Where("id = ?", item.ProductID).
		Scan(&currentPrice).Error

	if err != nil {
		return fmt.Errorf("failed to check product price: %w", err)
	}

	// Compare with cart item price (allow small differences due to precision)
	difference := currentPrice.Sub(item.UnitPrice).Abs()
	tolerance := decimal.NewFromFloat(0.01) // 1 cent tolerance

	if difference.GreaterThan(tolerance) {
		return fmt.Errorf("price changed from %s to %s",
			item.UnitPrice.String(), currentPrice.String())
	}

	return nil
}

func (s *cartService) validateShippingRestrictions(ctx context.Context, cart *models.Cart) error {
	// Basic shipping restrictions validation

	// Check if shipping address is required
	if cart.ShippingAddress == nil {
		return fmt.Errorf("shipping address is required")
	}

	// Check if we ship to the destination country
	allowedCountries := map[string]bool{
		"TH": true, // Thailand
		"SG": true, // Singapore
		"MY": true, // Malaysia
		"ID": true, // Indonesia
		"VN": true, // Vietnam
		"PH": true, // Philippines
		"US": true, // United States
	}

	if !allowedCountries[cart.ShippingAddress.Country] {
		return fmt.Errorf("shipping not available to %s", cart.ShippingAddress.Country)
	}

	// Check for restricted items (example: dangerous goods)
	for _, item := range cart.Items {
		// This would typically check product categories or tags
		// For demo, we'll simulate some restrictions
		var product struct {
			Name        string
			Category    string
			IsRestricted bool
		}

		err := s.db.WithContext(ctx).
			Table("products").
			Select("name, category").
			Where("id = ?", item.ProductID).
			First(&product).Error

		if err != nil {
			continue // Skip validation if product not found
		}

		// Example restriction: electronics to certain countries
		if product.Category == "electronics" && cart.ShippingAddress.Country == "ID" {
			return fmt.Errorf("electronics cannot be shipped to Indonesia due to import restrictions")
		}
	}

	return nil
}

// GetAbandonedCarts gets abandoned carts with filtering
func (s *cartService) GetAbandonedCarts(ctx context.Context, filters map[string]interface{}, pagination models.Pagination) (*models.CartResponse, error) {
	query := s.db.WithContext(ctx).Model(&models.Cart{}).Where("status = ?", models.CartStatusAbandoned)

	// Apply filters
	if minValue, ok := filters["min_value"]; ok {
		query = query.Where("total_amount >= ?", minValue)
	}
	if dateFrom, ok := filters["date_from"]; ok {
		query = query.Where("updated_at >= ?", dateFrom)
	}
	if dateTo, ok := filters["date_to"]; ok {
		query = query.Where("updated_at <= ?", dateTo)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count abandoned carts: %w", err)
	}

	// Get results with pagination
	var carts []*models.Cart
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Order("updated_at DESC").Offset(offset).Limit(pagination.PageSize).Find(&carts).Error; err != nil {
		return nil, fmt.Errorf("failed to get abandoned carts: %w", err)
	}

	totalPages := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	return &models.CartResponse{
		Carts:      carts,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

// CreateAbandonmentTracking creates cart abandonment tracking
func (s *cartService) CreateAbandonmentTracking(ctx context.Context, req models.CreateAbandonmentTrackingRequest) (*models.CartAbandonmentTracking, error) {
	tracking := &models.CartAbandonmentTracking{
		CartID:           req.CartID,
		AbandonmentStage: req.AbandonmentStage,
		LastPageVisited:  req.LastPageVisited,
		AbandonedAt:      time.Now(),
		EmailsSent:       0,
		RecoveryClicks:   0,
		IsRecovered:      false,
	}

	if err := s.db.WithContext(ctx).Create(tracking).Error; err != nil {
		return nil, fmt.Errorf("failed to create abandonment tracking: %w", err)
	}

	return tracking, nil
}

// GetAbandonmentAnalytics gets cart abandonment analytics
func (s *cartService) GetAbandonmentAnalytics(ctx context.Context, filters map[string]interface{}, pagination models.Pagination) (*models.AbandonmentTrackingResponse, error) {
	query := s.db.WithContext(ctx).Model(&models.CartAbandonmentTracking{})

	// Apply filters
	if stage, ok := filters["stage"]; ok {
		query = query.Where("abandonment_stage = ?", stage)
	}
	if dateFrom, ok := filters["date_from"]; ok {
		query = query.Where("abandoned_at >= ?", dateFrom)
	}
	if dateTo, ok := filters["date_to"]; ok {
		query = query.Where("abandoned_at <= ?", dateTo)
	}
	if isRecovered, ok := filters["is_recovered"]; ok {
		query = query.Where("is_recovered = ?", isRecovered)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count abandonment analytics: %w", err)
	}

	// Get results with pagination
	var tracking []*models.CartAbandonmentTracking
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Order("abandoned_at DESC").Offset(offset).Limit(pagination.PageSize).Find(&tracking).Error; err != nil {
		return nil, fmt.Errorf("failed to get abandonment analytics: %w", err)
	}

	totalPages := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	return &models.AbandonmentTrackingResponse{
		Tracking:   tracking,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}