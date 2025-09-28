package contract_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"tchat.dev/commerce/handlers"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/services"
	authModels "tchat.dev/auth/models"
)

// TestCommerceProviderContract runs Pact provider verification tests for the Commerce service
// This validates that our Commerce service meets the contract expectations from consumers
func TestCommerceProviderContract(t *testing.T) {
	// Create mock dependencies for the commerce service
	mockShopRepo := &MockShopRepository{shops: make(map[uuid.UUID]*models.Shop)}
	mockProductRepo := &MockProductRepository{products: make(map[uuid.UUID]*models.Product)}
	mockOrderRepo := &MockOrderRepository{orders: make(map[uuid.UUID]*models.Order)}
	mockCache := &MockCacheService{}
	mockEvents := &MockEventService{}
	mockPayment := &MockPaymentService{}

	// Create commerce service with mock dependencies
	commerceService := services.NewCommerceService(
		mockShopRepo,
		mockProductRepo,
		mockOrderRepo,
		mockCache,
		mockEvents,
		mockPayment,
		services.DefaultCommerceConfig(),
	)

	// Create commerce handler - pass the service pointer
	commerceHandler := handlers.NewCommerceHandler(*commerceService)

	// Set up Mux router (not Gin) for commerce routes
	router := mux.NewRouter()
	commerceHandler.RegisterRoutes(router)

	// Create test server
	testServer := httptest.NewServer(mockAuthMiddleware(router))
	defer testServer.Close()

	// Set up provider states for different test scenarios
	stateHandlers := map[string]func(setup bool, state map[string]interface{}) error{
		// Product catalog states
		"products exist in catalog": func(setup bool, state map[string]interface{}) error {
			return setupProductCatalogState(mockProductRepo, setup)
		},
		"product with ID exists": func(setup bool, state map[string]interface{}) error {
			if productIDValue, exists := state["product_id"]; exists {
				productID := productIDValue.(string)
				return setupSingleProductState(mockProductRepo, productID, setup)
			}
			return setupSingleProductState(mockProductRepo, "11111111-1111-1111-1111-111111111111", setup)
		},
		"shop has products": func(setup bool, state map[string]interface{}) error {
			if shopIDValue, exists := state["shop_id"]; exists {
				shopID := shopIDValue.(string)
				return setupShopProductsState(mockProductRepo, shopID, setup)
			}
			return setupShopProductsState(mockProductRepo, "22222222-2222-2222-2222-222222222222", setup)
		},

		// Shopping cart states
		"user has items in cart": func(setup bool, state map[string]interface{}) error {
			if userIDValue, exists := state["user_id"]; exists {
				userID := userIDValue.(string)
				return setupUserCartState(mockProductRepo, userID, setup)
			}
			return setupUserCartState(mockProductRepo, "33333333-3333-3333-3333-333333333333", setup)
		},

		// Authentication states
		"user is authenticated for checkout": func(setup bool, state map[string]interface{}) error {
			return setupAuthenticatedUserState(setup)
		},
		"user has valid JWT token": func(setup bool, state map[string]interface{}) error {
			return setupValidJWTState(setup)
		},

		// Order states
		"order exists for user": func(setup bool, state map[string]interface{}) error {
			orderID := "44444444-4444-4444-4444-444444444444"
			userID := "55555555-5555-5555-5555-555555555555"
			if orderIDValue, exists := state["order_id"]; exists {
				orderID = orderIDValue.(string)
			}
			if userIDValue, exists := state["user_id"]; exists {
				userID = userIDValue.(string)
			}
			return setupOrderExistsState(mockOrderRepo, orderID, userID, setup)
		},
		"pending order exists": func(setup bool, state map[string]interface{}) error {
			orderID := "66666666-6666-6666-6666-666666666666"
			if orderIDValue, exists := state["order_id"]; exists {
				orderID = orderIDValue.(string)
			}
			return setupPendingOrderState(mockOrderRepo, orderID, setup)
		},

		// Shop states
		"shop exists with products": func(setup bool, state map[string]interface{}) error {
			shopID := "77777777-7777-7777-7777-777777777777"
			if shopIDValue, exists := state["shop_id"]; exists {
				shopID = shopIDValue.(string)
			}
			return setupShopExistsState(mockShopRepo, shopID, setup)
		},
		"user owns shop": func(setup bool, state map[string]interface{}) error {
			userID := "88888888-8888-8888-8888-888888888888"
			shopID := "99999999-9999-9999-9999-999999999999"
			if userIDValue, exists := state["user_id"]; exists {
				userID = userIDValue.(string)
			}
			if shopIDValue, exists := state["shop_id"]; exists {
				shopID = shopIDValue.(string)
			}
			return setupUserOwnsShopState(mockShopRepo, userID, shopID, setup)
		},
	}

	// Configure Pact provider verification
	// For now, we'll implement a simplified test that validates the service setup
	// In a real implementation, this would use the actual Pact verification process

	// Test 1: Validate that we can set up all provider states
	t.Run("validate_provider_states", func(t *testing.T) {
		for stateName, stateHandler := range stateHandlers {
			t.Run(stateName, func(t *testing.T) {
				// Set up the state
				err := stateHandler(true, make(map[string]interface{}))
				assert.NoError(t, err, "Should be able to set up state: %s", stateName)

				// Clean up the state
				err = stateHandler(false, make(map[string]interface{}))
				assert.NoError(t, err, "Should be able to clean up state: %s", stateName)
			})
		}
	})

	// Test 2: Validate that the service responds correctly to basic requests
	t.Run("validate_service_endpoints", func(t *testing.T) {
		// Set up products in catalog state
		err := stateHandlers["products exist in catalog"](true, make(map[string]interface{}))
		assert.NoError(t, err)

		// Make a request to the test server
		resp, err := http.Get(testServer.URL + "/commerce/products")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Clean up
		stateHandlers["products exist in catalog"](false, make(map[string]interface{}))
	})

	// Test 3: Validate authentication handling
	t.Run("validate_authentication", func(t *testing.T) {
		// Test without authentication
		resp, err := http.Post(testServer.URL+"/commerce/shops", "application/json", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Test with mock authentication (would be handled by middleware)
		client := &http.Client{}
		req, _ := http.NewRequest("GET", testServer.URL+"/commerce/orders", nil)
		req.Header.Set("Authorization", "Bearer mock-token")
		resp, err = client.Do(req)
		assert.NoError(t, err)
		// Should not be unauthorized (might be other status based on implementation)
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// Provider state setup functions

// setupProductCatalogState sets up a product catalog with multiple products
func setupProductCatalogState(mockRepo *MockProductRepository, setup bool) error {
	if !setup {
		mockRepo.Clear()
		return nil
	}

	// Create test products for catalog
	products := []*models.Product{
		{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			ShopID:      uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Name:        "Premium Wireless Headphones",
			Description: stringPtr("High-quality wireless headphones with noise cancellation"),
			Category:    models.CategoryElectronics,
			Price:       299900, // $2999.00 in cents
			Currency:    models.CurrencyUSD,
			Status:      models.ProductStatusActive,
			Visibility:  models.VisibilityPublic,
			Inventory: models.Inventory{
				TrackQuantity: true,
				Quantity:      50,
				StockStatus:   models.StockStatusInStock,
			},
			Images: models.ImageSlice{
				{
					ID:       "img-001",
					URL:      "https://cdn.tchat.dev/products/headphones-main.jpg",
					AltText:  "Premium wireless headphones",
					Position: 1,
					IsMain:   true,
				},
			},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111112"),
			ShopID:      uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			Name:        "Organic Cotton T-Shirt",
			Description: stringPtr("100% organic cotton t-shirt, ethically sourced"),
			Category:    models.CategoryFashion,
			Price:       4990, // $49.90 in cents
			Currency:    models.CurrencyUSD,
			Status:      models.ProductStatusActive,
			Visibility:  models.VisibilityPublic,
			Inventory: models.Inventory{
				TrackQuantity: true,
				Quantity:      100,
				StockStatus:   models.StockStatusInStock,
			},
			Variants: models.VariantSlice{
				{
					ID:   "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					SKU:  "TSHIRT-S-WHITE",
					Name: "Small White",
					Options: map[string]string{
						"size":  "S",
						"color": "white",
					},
					Inventory: models.Inventory{
						TrackQuantity: true,
						Quantity:      25,
						StockStatus:   models.StockStatusInStock,
					},
				},
				{
					ID:   "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					SKU:  "TSHIRT-M-WHITE",
					Name: "Medium White",
					Options: map[string]string{
						"size":  "M",
						"color": "white",
					},
					Inventory: models.Inventory{
						TrackQuantity: true,
						Quantity:      30,
						StockStatus:   models.StockStatusInStock,
					},
				},
			},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111113"),
			ShopID:      uuid.MustParse("22222222-2222-2222-2222-222222222223"),
			Name:        "Thai Jasmine Rice",
			Description: stringPtr("Premium quality jasmine rice from Thailand"),
			Category:    models.CategoryFood,
			Price:       1590, // $15.90 in cents
			Currency:    models.CurrencyUSD,
			Status:      models.ProductStatusActive,
			Visibility:  models.VisibilityPublic,
			Inventory: models.Inventory{
				TrackQuantity: true,
				Quantity:      200,
				StockStatus:   models.StockStatusInStock,
			},
			Tags: []string{"organic", "thai", "premium", "rice"},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	for _, product := range products {
		mockRepo.products[product.ID] = product
	}

	return nil
}

// setupSingleProductState sets up a single product by ID
func setupSingleProductState(mockRepo *MockProductRepository, productIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return fmt.Errorf("invalid product ID: %v", err)
	}

	product := &models.Product{
		ID:          productID,
		ShopID:      uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Name:        "Test Product",
		Description: stringPtr("Test product description"),
		Category:    models.CategoryElectronics,
		Price:       9999, // $99.99 in cents
		Currency:    models.CurrencyUSD,
		Status:      models.ProductStatusActive,
		Visibility:  models.VisibilityPublic,
		Inventory: models.Inventory{
			TrackQuantity: true,
			Quantity:      10,
			StockStatus:   models.StockStatusInStock,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mockRepo.products[productID] = product
	return nil
}

// setupShopProductsState sets up products for a specific shop
func setupShopProductsState(mockRepo *MockProductRepository, shopIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return fmt.Errorf("invalid shop ID: %v", err)
	}

	products := []*models.Product{
		{
			ID:          uuid.New(),
			ShopID:      shopID,
			Name:        "Shop Product 1",
			Description: stringPtr("First product from the shop"),
			Category:    models.CategoryElectronics,
			Price:       5000,
			Currency:    models.CurrencyUSD,
			Status:      models.ProductStatusActive,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			ShopID:      shopID,
			Name:        "Shop Product 2",
			Description: stringPtr("Second product from the shop"),
			Category:    models.CategoryFashion,
			Price:       7500,
			Currency:    models.CurrencyUSD,
			Status:      models.ProductStatusActive,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		},
	}

	for _, product := range products {
		mockRepo.products[product.ID] = product
	}

	return nil
}

// setupUserCartState sets up cart items for a user
func setupUserCartState(mockRepo *MockProductRepository, userIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	// Set up products that will be in the user's cart
	return setupProductCatalogState(mockRepo, true)
}

// setupAuthenticatedUserState sets up an authenticated user
func setupAuthenticatedUserState(setup bool) error {
	// This is handled by the request filter middleware
	return nil
}

// setupValidJWTState sets up valid JWT authentication
func setupValidJWTState(setup bool) error {
	// This is handled by the request filter middleware
	return nil
}

// setupOrderExistsState sets up an existing order for a user
func setupOrderExistsState(mockRepo *MockOrderRepository, orderIDStr, userIDStr string, setup bool) error {
	if !setup {
		mockRepo.Clear()
		return nil
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return fmt.Errorf("invalid order ID: %v", err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	order := &models.Order{
		ID:              orderID,
		CustomerID:      userID,
		ShopID:          uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		OrderNumber:     "ORD-20240101-123456",
		Status:          models.OrderStatusConfirmed,
		PaymentStatus:   models.PaymentStatusPaid,
		FulfillmentStatus: models.FulfillmentStatusUnfulfilled,
		Currency:        models.CurrencyUSD,
		SubtotalAmount:  9999,
		TaxAmount:       1000,
		ShippingAmount:  500,
		TotalAmount:     11499,
		Items: []models.OrderItem{
			{
				ID:          "cccccccc-cccc-cccc-cccc-cccccccccccc",
				ProductID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				Quantity:    1,
				UnitPrice:   9999,
				TotalPrice:  9999,
				Name:        "Test Product",
				FulfillmentStatus: models.FulfillmentStatusUnfulfilled,
			},
		},
		ShippingAddress: models.Address{
			FirstName:    "John",
			LastName:     "Doe",
			AddressLine1: "123 Main St",
			City:         "Bangkok",
			State:        "Bangkok",
			PostalCode:   "10110",
			Country:      "TH",
		},
		BillingAddress: models.Address{
			FirstName:    "John",
			LastName:     "Doe",
			AddressLine1: "123 Main St",
			City:         "Bangkok",
			State:        "Bangkok",
			PostalCode:   "10110",
			Country:      "TH",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mockRepo.orders[orderID] = order
	return nil
}

// setupPendingOrderState sets up a pending order
func setupPendingOrderState(mockRepo *MockOrderRepository, orderIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return fmt.Errorf("invalid order ID: %v", err)
	}

	order := &models.Order{
		ID:              orderID,
		CustomerID:      uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		ShopID:          uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		OrderNumber:     "ORD-20240101-PENDING",
		Status:          models.OrderStatusPending,
		PaymentStatus:   models.PaymentStatusPending,
		FulfillmentStatus: models.FulfillmentStatusUnfulfilled,
		Currency:        models.CurrencyUSD,
		TotalAmount:     5000,
		SubtotalAmount:  4500,
		TaxAmount:       450,
		ShippingAmount:  50,
		Items: []models.OrderItem{
			{
				ID:          "dddddddd-dddd-dddd-dddd-dddddddddddd",
				ProductID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				Quantity:    1,
				UnitPrice:   4500,
				TotalPrice:  4500,
				Name:        "Pending Product",
				FulfillmentStatus: models.FulfillmentStatusUnfulfilled,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mockRepo.orders[orderID] = order
	return nil
}

// setupShopExistsState sets up an existing shop
func setupShopExistsState(mockRepo *MockShopRepository, shopIDStr string, setup bool) error {
	if !setup {
		mockRepo.Clear()
		return nil
	}

	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return fmt.Errorf("invalid shop ID: %v", err)
	}

	shop := &models.Shop{
		ID:          shopID,
		OwnerID:     uuid.MustParse("88888888-8888-8888-8888-888888888888"),
		Name:        "Test Electronics Store",
		Description: stringPtr("Premium electronics and gadgets"),
		Slug:        "test-electronics-store",
		Category:    models.CategoryElectronicsShop,
		Email:       "shop@test.com",
		Phone:       stringPtr("+66812345678"),
		Status:      models.ShopStatusActive,
		Visibility:  models.ShopVisibilityPublic,
		Country:     "TH",
		Currency:    models.CurrencyTHB,
		Language:    "th",
		Timezone:    "Asia/Bangkok",
		IsVerified:  true,
		Analytics: models.ShopAnalytics{
			ActiveProducts: 25,
			TotalOrders:    150,
			ConversionRate: 4.5,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	mockRepo.shops[shopID] = shop
	return nil
}

// setupUserOwnsShopState sets up a user as owner of a shop
func setupUserOwnsShopState(mockRepo *MockShopRepository, userIDStr, shopIDStr string, setup bool) error {
	if !setup {
		return nil
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return fmt.Errorf("invalid shop ID: %v", err)
	}

	shop := &models.Shop{
		ID:          shopID,
		OwnerID:     userID,
		Name:        "User's Shop",
		Description: stringPtr("Shop owned by the authenticated user"),
		Slug:        "users-shop",
		Category:    models.CategoryRetailShop,
		Email:       "owner@shop.com",
		Status:      models.ShopStatusActive,
		Visibility:  models.ShopVisibilityPublic,
		Country:     "TH",
		Currency:    models.CurrencyTHB,
		Language:    "th",
		Timezone:    "Asia/Bangkok",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	mockRepo.shops[shopID] = shop
	return nil
}

// Mock authentication middleware for Pact tests
func mockAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock authentication for Pact provider tests
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && authHeader != "Bearer invalid-token" {
			// Set up mock user context for valid tokens
			user := &authModels.User{
				ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				PhoneNumber: stringPtr("+66812345678"),
				CountryCode: stringPtr("TH"),
				Status:      "active",
			}
			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// Utility functions

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

// Mock implementations for testing

type MockShopRepository struct {
	shops map[uuid.UUID]*models.Shop
}

func (m *MockShopRepository) Create(ctx context.Context, shop *models.Shop) error {
	if m.shops == nil {
		m.shops = make(map[uuid.UUID]*models.Shop)
	}
	m.shops[shop.ID] = shop
	return nil
}

func (m *MockShopRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Shop, error) {
	if shop, exists := m.shops[id]; exists {
		return shop, nil
	}
	return nil, fmt.Errorf("shop not found")
}

func (m *MockShopRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*models.Shop, error) {
	var shops []*models.Shop
	for _, shop := range m.shops {
		if shop.OwnerID == ownerID {
			shops = append(shops, shop)
		}
	}
	return shops, nil
}

func (m *MockShopRepository) GetBySlug(ctx context.Context, slug string) (*models.Shop, error) {
	for _, shop := range m.shops {
		if shop.Slug == slug {
			return shop, nil
		}
	}
	return nil, fmt.Errorf("shop not found")
}

func (m *MockShopRepository) Update(ctx context.Context, shop *models.Shop) error {
	m.shops[shop.ID] = shop
	return nil
}

func (m *MockShopRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.shops, id)
	return nil
}

func (m *MockShopRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.Shop, error) {
	var shops []*models.Shop
	for _, shop := range m.shops {
		shops = append(shops, shop)
	}
	return shops, nil
}

func (m *MockShopRepository) GetFeatured(ctx context.Context, limit int) ([]*models.Shop, error) {
	var shops []*models.Shop
	count := 0
	for _, shop := range m.shops {
		if shop.IsFeatured && count < limit {
			shops = append(shops, shop)
			count++
		}
	}
	return shops, nil
}

func (m *MockShopRepository) Clear() {
	m.shops = make(map[uuid.UUID]*models.Shop)
}

type MockProductRepository struct {
	products map[uuid.UUID]*models.Product
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	if m.products == nil {
		m.products = make(map[uuid.UUID]*models.Product)
	}
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	if product, exists := m.products[id]; exists {
		return product, nil
	}
	return nil, fmt.Errorf("product not found")
}

func (m *MockProductRepository) GetByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	for _, product := range m.products {
		if product.ShopID == shopID {
			products = append(products, product)
		}
	}
	return products, nil
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	for _, product := range m.products {
		if product.SKU == sku {
			return product, nil
		}
	}
	return nil, fmt.Errorf("product not found")
}

func (m *MockProductRepository) Update(ctx context.Context, product *models.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.products, id)
	return nil
}

func (m *MockProductRepository) Search(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	count := 0
	skip := 0

	for _, product := range m.products {
		if skip < offset {
			skip++
			continue
		}
		if count >= limit {
			break
		}

		// Apply filters
		if category, exists := filters["category"]; exists {
			if string(product.Category) != category.(string) {
				continue
			}
		}
		if shopID, exists := filters["shop_id"]; exists {
			if product.ShopID != shopID.(uuid.UUID) {
				continue
			}
		}
		if minPrice, exists := filters["min_price"]; exists {
			if product.Price < minPrice.(int64) {
				continue
			}
		}
		if maxPrice, exists := filters["max_price"]; exists {
			if product.Price > maxPrice.(int64) {
				continue
			}
		}

		products = append(products, product)
		count++
	}

	return products, nil
}

func (m *MockProductRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	for _, product := range m.products {
		if string(product.Category) == category {
			products = append(products, product)
		}
	}
	return products, nil
}

func (m *MockProductRepository) UpdateInventory(ctx context.Context, productID uuid.UUID, quantity int) error {
	if product, exists := m.products[productID]; exists {
		product.Inventory.Quantity = quantity
		return nil
	}
	return fmt.Errorf("product not found")
}

func (m *MockProductRepository) LockForUpdate(ctx context.Context, productID uuid.UUID) (*models.Product, error) {
	return m.GetByID(ctx, productID)
}

func (m *MockProductRepository) Clear() {
	m.products = make(map[uuid.UUID]*models.Product)
}

type MockOrderRepository struct {
	orders map[uuid.UUID]*models.Order
}

func (m *MockOrderRepository) Create(ctx context.Context, order *models.Order) error {
	if m.orders == nil {
		m.orders = make(map[uuid.UUID]*models.Order)
	}
	m.orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	if order, exists := m.orders[id]; exists {
		return order, nil
	}
	return nil, fmt.Errorf("order not found")
}

func (m *MockOrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	for _, order := range m.orders {
		if order.CustomerID == userID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (m *MockOrderRepository) GetByShopID(ctx context.Context, shopID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	for _, order := range m.orders {
		if order.ShopID == shopID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (m *MockOrderRepository) Update(ctx context.Context, order *models.Order) error {
	m.orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.orders, id)
	return nil
}

func (m *MockOrderRepository) GetPendingOrders(ctx context.Context) ([]*models.Order, error) {
	var orders []*models.Order
	for _, order := range m.orders {
		if order.Status == models.OrderStatusPending {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	if order, exists := m.orders[orderID]; exists {
		order.Status = status
		return nil
	}
	return fmt.Errorf("order not found")
}

func (m *MockOrderRepository) Clear() {
	m.orders = make(map[uuid.UUID]*models.Order)
}

type MockCacheService struct{}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	return nil
}

func (m *MockCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	return nil, nil
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	return nil
}

func (m *MockCacheService) Lock(ctx context.Context, key string, expiry time.Duration) (bool, error) {
	return true, nil
}

func (m *MockCacheService) Unlock(ctx context.Context, key string) error {
	return nil
}

type MockEventService struct{}

func (m *MockEventService) PublishShop(ctx context.Context, event *services.ShopEvent) error {
	return nil
}

func (m *MockEventService) PublishProduct(ctx context.Context, event *services.ProductEvent) error {
	return nil
}

func (m *MockEventService) PublishOrder(ctx context.Context, event *services.OrderEvent) error {
	return nil
}

type MockPaymentService struct{}

func (m *MockPaymentService) CreatePayment(ctx context.Context, req *services.PaymentRequest) (*services.PaymentResponse, error) {
	return &services.PaymentResponse{
		PaymentID:   fmt.Sprintf("pay_%s", uuid.New().String()[:8]),
		Status:      "pending",
		RedirectURL: stringPtr("https://payment.example.com/checkout"),
		ExpiresAt:   time.Now().Add(30 * time.Minute),
	}, nil
}

func (m *MockPaymentService) ProcessRefund(ctx context.Context, req *services.RefundRequest) (*services.RefundResponse, error) {
	return &services.RefundResponse{
		RefundID:    fmt.Sprintf("ref_%s", uuid.New().String()[:8]),
		Status:      "processed",
		ProcessedAt: time.Now(),
	}, nil
}

func (m *MockPaymentService) ValidatePayment(ctx context.Context, paymentID string) (*services.PaymentStatus, error) {
	return &services.PaymentStatus{
		PaymentID: paymentID,
		Status:    "completed",
		Amount:    10000,
		PaidAt:    timePtr(time.Now()),
	}, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}