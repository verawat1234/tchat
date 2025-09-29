package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"tchat.dev/commerce/models"
)

// MockCategoryRepository for testing
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByPath(ctx context.Context, path string) (*models.Category, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]*models.Category, int64, error) {
	args := m.Called(ctx, filters, offset, limit)
	return args.Get(0).([]*models.Category), args.Get(1).(int64), args.Error(2)
}

func (m *MockCategoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*models.Category, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByBusinessID(ctx context.Context, businessID uuid.UUID, offset, limit int) ([]*models.Category, int64, error) {
	args := m.Called(ctx, businessID, offset, limit)
	return args.Get(0).([]*models.Category), args.Get(1).(int64), args.Error(2)
}

func (m *MockCategoryRepository) GetGlobalCategories(ctx context.Context, offset, limit int) ([]*models.Category, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*models.Category), args.Get(1).(int64), args.Error(2)
}

func (m *MockCategoryRepository) GetFeaturedCategories(ctx context.Context, businessID *uuid.UUID, limit int) ([]*models.Category, error) {
	args := m.Called(ctx, businessID, limit)
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetRootCategories(ctx context.Context, businessID *uuid.UUID) ([]*models.Category, error) {
	args := m.Called(ctx, businessID)
	return args.Get(0).([]*models.Category), args.Error(1)
}

// MockProductCategoryRepository for testing
type MockProductCategoryRepository struct {
	mock.Mock
}

func (m *MockProductCategoryRepository) Create(ctx context.Context, productCategory *models.ProductCategory) error {
	args := m.Called(ctx, productCategory)
	return args.Error(0)
}

func (m *MockProductCategoryRepository) Delete(ctx context.Context, productID, categoryID uuid.UUID) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockProductCategoryRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*models.ProductCategory, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]*models.ProductCategory), args.Error(1)
}

func (m *MockProductCategoryRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, offset, limit int) ([]*models.ProductCategory, int64, error) {
	args := m.Called(ctx, categoryID, offset, limit)
	return args.Get(0).([]*models.ProductCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductCategoryRepository) Update(ctx context.Context, productCategory *models.ProductCategory) error {
	args := m.Called(ctx, productCategory)
	return args.Error(0)
}

func (m *MockProductCategoryRepository) SetPrimary(ctx context.Context, productID, categoryID uuid.UUID) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockProductCategoryRepository) UnsetPrimary(ctx context.Context, productID uuid.UUID) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

// MockCategoryViewRepository for testing
type MockCategoryViewRepository struct {
	mock.Mock
}

func (m *MockCategoryViewRepository) Create(ctx context.Context, view *models.CategoryView) error {
	args := m.Called(ctx, view)
	return args.Error(0)
}

func (m *MockCategoryViewRepository) GetAnalytics(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (map[string]interface{}, error) {
	args := m.Called(ctx, categoryID, dateFrom, dateTo)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCategoryViewRepository) CountViews(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (int64, error) {
	args := m.Called(ctx, categoryID, dateFrom, dateTo)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCategoryViewRepository) CountUniqueVisitors(ctx context.Context, categoryID uuid.UUID, dateFrom, dateTo string) (int64, error) {
	args := m.Called(ctx, categoryID, dateFrom, dateTo)
	return args.Get(0).(int64), args.Error(1)
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Suppress SQL logs in tests
	})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Create simple tables manually for testing (avoiding complex GORM migrations)
	db.Exec(`CREATE TABLE categories (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		parent_id TEXT,
		level INTEGER DEFAULT 0,
		product_count INTEGER DEFAULT 0,
		active_product_count INTEGER DEFAULT 0,
		children_count INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME
	)`)

	db.Exec(`CREATE TABLE category_views (
		id TEXT PRIMARY KEY,
		category_id TEXT NOT NULL,
		user_id TEXT,
		session_id TEXT NOT NULL,
		ip_address TEXT,
		user_agent TEXT,
		referrer TEXT,
		viewed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)

	db.Exec(`CREATE TABLE product_categories (
		id TEXT PRIMARY KEY,
		product_id TEXT NOT NULL,
		category_id TEXT NOT NULL,
		is_primary BOOLEAN DEFAULT FALSE,
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)

	return db
}

// setupCategoryService creates a category service with mock repositories
func setupCategoryService() (*categoryService, *MockCategoryRepository, *MockProductCategoryRepository, *MockCategoryViewRepository) {
	mockCategoryRepo := &MockCategoryRepository{}
	mockProductCategoryRepo := &MockProductCategoryRepository{}
	mockCategoryViewRepo := &MockCategoryViewRepository{}

	service := &categoryService{
		categoryRepo:        mockCategoryRepo,
		productCategoryRepo: mockProductCategoryRepo,
		categoryViewRepo:    mockCategoryViewRepo,
		db:                  nil, // Don't test database-dependent methods
	}

	return service, mockCategoryRepo, mockProductCategoryRepo, mockCategoryViewRepo
}

func TestCategoryService_GetRootCategories(t *testing.T) {
	service, mockCategoryRepo, _, _ := setupCategoryService()

	businessID := uuid.New()
	expectedCategories := []*models.Category{
		{
			ID:   uuid.New(),
			Name: "Root Category 1",
		},
		{
			ID:   uuid.New(),
			Name: "Root Category 2",
		},
	}

	mockCategoryRepo.On("GetRootCategories", mock.Anything, &businessID).Return(expectedCategories, nil)

	categories, err := service.GetRootCategories(context.Background(), &businessID)

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Root Category 1", categories[0].Name)
	assert.Equal(t, "Root Category 2", categories[1].Name)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryService_TrackCategoryView(t *testing.T) {
	service, _, _, _ := setupCategoryService()

	// Set up a proper test database for this method
	testDB := setupTestDB()
	service.db = testDB

	categoryID := uuid.New()
	userID := uuid.New()

	err := service.TrackCategoryView(context.Background(), categoryID, &userID, "session123", "192.168.1.1", "Mozilla/5.0", "https://example.com")

	assert.NoError(t, err)
}

func TestCategoryService_GetCategoryAnalytics(t *testing.T) {
	service, _, _, mockCategoryViewRepo := setupCategoryService()

	// Set up a proper test database for this method since it calls getProductClicksForCategory
	testDB := setupTestDB()
	service.db = testDB

	categoryID := uuid.New()
	dateFrom := "2023-01-01"
	dateTo := "2023-12-31"

	// Mock analytics data
	analyticsData := map[string]interface{}{
		"total_views":     int64(1000),
		"unique_visitors": int64(500),
		"top_referrers": []map[string]interface{}{
			{"referrer": "google.com", "count": int64(100)},
		},
		"views_by_day": []map[string]interface{}{
			{"date": "2023-01-01", "views": int64(50)},
		},
	}

	mockCategoryViewRepo.On("GetAnalytics", mock.Anything, categoryID, dateFrom, dateTo).Return(analyticsData, nil)

	analytics, err := service.GetCategoryAnalytics(context.Background(), categoryID, dateFrom, dateTo)

	assert.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, int64(1000), analytics.TotalViews)
	assert.Equal(t, int64(500), analytics.UniqueVisitors)
	assert.Len(t, analytics.TopReferrers, 1)
	assert.Equal(t, "google.com", analytics.TopReferrers[0].Referrer)
	assert.Len(t, analytics.ViewsByDay, 1)

	mockCategoryViewRepo.AssertExpectations(t)
}

func TestCategoryService_AddProductToCategory(t *testing.T) {
	service, _, _, _ := setupCategoryService()

	// Set up a proper test database for this method
	testDB := setupTestDB()
	service.db = testDB

	productID := uuid.New()
	categoryID := uuid.New()

	// Insert test category to satisfy the stats update
	testDB.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", categoryID.String(), "Test Category")

	// This method uses direct database operations, not repository calls
	err := service.AddProductToCategory(context.Background(), productID, categoryID, true)

	assert.NoError(t, err)

	// Verify the relationship was created in the database
	var count int64
	testDB.Raw("SELECT COUNT(*) FROM product_categories WHERE product_id = ? AND category_id = ?",
		productID.String(), categoryID.String()).Scan(&count)
	assert.Equal(t, int64(1), count)
}

func TestCategoryService_RemoveProductFromCategory(t *testing.T) {
	service, _, _, _ := setupCategoryService()

	// Set up a proper test database for this method
	testDB := setupTestDB()
	service.db = testDB

	productID := uuid.New()
	categoryID := uuid.New()

	// Insert test category and product relationship to satisfy the existence check
	testDB.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", categoryID.String(), "Test Category")
	testDB.Exec("INSERT INTO product_categories (id, product_id, category_id) VALUES (?, ?, ?)",
		uuid.New().String(), productID.String(), categoryID.String())

	// This method uses direct database operations, not repository calls
	err := service.RemoveProductFromCategory(context.Background(), productID, categoryID)

	assert.NoError(t, err)

	// Verify the relationship was removed from the database
	var count int64
	testDB.Raw("SELECT COUNT(*) FROM product_categories WHERE product_id = ? AND category_id = ?",
		productID.String(), categoryID.String()).Scan(&count)
	assert.Equal(t, int64(0), count)
}

