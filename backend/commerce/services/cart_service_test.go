package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"tchat.dev/commerce/models"
)

// MockCartRepository for testing
type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) Create(ctx context.Context, cart *models.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *MockCartRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Cart, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Cart), args.Error(1)
}

func (m *MockCartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.Cart), args.Error(1)
}

func (m *MockCartRepository) GetBySessionID(ctx context.Context, sessionID string) (*models.Cart, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(*models.Cart), args.Error(1)
}

func (m *MockCartRepository) Update(ctx context.Context, cart *models.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *MockCartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCartRepository) GetExpiredCarts(ctx context.Context, cutoffTime time.Time) ([]*models.Cart, error) {
	args := m.Called(ctx, cutoffTime)
	return args.Get(0).([]*models.Cart), args.Error(1)
}

func (m *MockCartRepository) GetAbandonedCarts(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Cart, int64, error) {
	args := m.Called(ctx, filters, offset, limit)
	return args.Get(0).([]*models.Cart), args.Get(1).(int64), args.Error(2)
}

// MockCartAbandonmentRepository for testing
type MockCartAbandonmentRepository struct {
	mock.Mock
}

func (m *MockCartAbandonmentRepository) Create(ctx context.Context, abandonment *models.CartAbandonmentTracking) error {
	args := m.Called(ctx, abandonment)
	return args.Error(0)
}

func (m *MockCartAbandonmentRepository) GetByCartID(ctx context.Context, cartID uuid.UUID) (*models.CartAbandonmentTracking, error) {
	args := m.Called(ctx, cartID)
	return args.Get(0).(*models.CartAbandonmentTracking), args.Error(1)
}

func (m *MockCartAbandonmentRepository) Update(ctx context.Context, abandonment *models.CartAbandonmentTracking) error {
	args := m.Called(ctx, abandonment)
	return args.Error(0)
}

func (m *MockCartAbandonmentRepository) List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.CartAbandonmentTracking, int64, error) {
	args := m.Called(ctx, filters, offset, limit)
	return args.Get(0).([]*models.CartAbandonmentTracking), args.Get(1).(int64), args.Error(2)
}

func (m *MockCartAbandonmentRepository) GetUnrecoveredAbandoned(ctx context.Context, olderThan time.Time) ([]*models.CartAbandonmentTracking, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).([]*models.CartAbandonmentTracking), args.Error(1)
}

func TestCartService_ApplyCoupon(t *testing.T) {
	mockCartRepo := &MockCartRepository{}
	mockAbandonmentRepo := &MockCartAbandonmentRepository{}

	service := NewCartService(mockCartRepo, mockAbandonmentRepo, nil)

	cart := &models.Cart{
		ID:             uuid.New(),
		SubtotalAmount: decimal.NewFromFloat(100.0),
		Currency:       "USD",
	}

	// Test applying valid coupon
	mockCartRepo.On("GetByID", mock.Anything, cart.ID).Return(cart, nil)
	mockCartRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Cart")).Return(nil)

	err := service.ApplyCoupon(context.Background(), cart.ID, "SAVE10")
	assert.NoError(t, err)

	mockCartRepo.AssertExpectations(t)
}

func TestCartService_RemoveCoupon(t *testing.T) {
	mockCartRepo := &MockCartRepository{}
	mockAbandonmentRepo := &MockCartAbandonmentRepository{}

	service := NewCartService(mockCartRepo, mockAbandonmentRepo, nil)

	cart := &models.Cart{
		ID:         uuid.New(),
		CouponCode: "SAVE10",
	}

	mockCartRepo.On("GetByID", mock.Anything, cart.ID).Return(cart, nil)
	mockCartRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Cart")).Return(nil)

	err := service.RemoveCoupon(context.Background(), cart.ID)
	assert.NoError(t, err)

	mockCartRepo.AssertExpectations(t)
}

func TestCartService_ValidateCart(t *testing.T) {
	mockCartRepo := &MockCartRepository{}
	mockAbandonmentRepo := &MockCartAbandonmentRepository{}

	service := NewCartService(mockCartRepo, mockAbandonmentRepo, nil)

	// Test cart with valid items
	cart := &models.Cart{
		ID: uuid.New(),
		Items: []models.CartItem{
			{
				ID:            uuid.New(),
				ProductID:     uuid.New(),
				Quantity:      2,
				UnitPrice:     decimal.NewFromFloat(50.0),
				TotalPrice:    decimal.NewFromFloat(100.0),
				ProductName:   "Test Product",
				IsAvailable:   true,
				StockQuantity: 10,
				MaxQuantity:   5,
			},
		},
		SubtotalAmount: decimal.NewFromFloat(100.0),
		Currency:       "USD",
	}

	mockCartRepo.On("GetByID", mock.Anything, cart.ID).Return(cart, nil)

	validation, err := service.ValidateCart(context.Background(), cart.ID)

	assert.NoError(t, err)
	assert.NotNil(t, validation)
	assert.True(t, validation.IsValid)
	assert.Empty(t, validation.Issues)

	mockCartRepo.AssertExpectations(t)
}

func TestCartService_GetAbandonedCarts(t *testing.T) {
	mockCartRepo := &MockCartRepository{}
	mockAbandonmentRepo := &MockCartAbandonmentRepository{}

	service := NewCartService(mockCartRepo, mockAbandonmentRepo, nil)

	filters := map[string]interface{}{
		"min_value": 50.0,
	}
	pagination := models.Pagination{
		Page:     1,
		PageSize: 20,
	}

	expectedCarts := []*models.Cart{
		{
			ID:             uuid.New(),
			Status:         models.CartStatusAbandoned,
			SubtotalAmount: decimal.NewFromFloat(100.0),
		},
	}

	mockCartRepo.On("GetAbandonedCarts", mock.Anything, filters, 0, 20).Return(expectedCarts, int64(1), nil)

	response, err := service.GetAbandonedCarts(context.Background(), filters, pagination)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(1), response.Total)
	assert.Len(t, response.Carts, 1)

	mockCartRepo.AssertExpectations(t)
}

func TestCartService_CreateAbandonmentTracking(t *testing.T) {
	mockCartRepo := &MockCartRepository{}
	mockAbandonmentRepo := &MockCartAbandonmentRepository{}

	service := NewCartService(mockCartRepo, mockAbandonmentRepo, nil)

	req := models.CreateAbandonmentTrackingRequest{
		CartID:           uuid.New(),
		AbandonmentStage: "checkout",
		LastPageVisited:  "/checkout",
	}

	mockAbandonmentRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.CartAbandonmentTracking")).Return(nil)

	tracking, err := service.CreateAbandonmentTracking(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, tracking)
	assert.Equal(t, req.CartID, tracking.CartID)
	assert.Equal(t, req.AbandonmentStage, tracking.AbandonmentStage)
	assert.Equal(t, req.LastPageVisited, tracking.LastPageVisited)
	assert.False(t, tracking.IsRecovered)

	mockAbandonmentRepo.AssertExpectations(t)
}