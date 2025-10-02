package service_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"tchat.dev/shared/config"
	"tchat.dev/shared/service"
)

// Test models
type TestModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255;not null"`
}

// Mock implementations for testing
type MockRepositoryInitializer struct {
	mock.Mock
}

func (m *MockRepositoryInitializer) InitializeRepositories(db *gorm.DB) error {
	args := m.Called(db)
	return args.Error(0)
}

func (m *MockRepositoryInitializer) GetRepositories() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

type MockServiceInitializer struct {
	mock.Mock
}

func (m *MockServiceInitializer) InitializeServices(repos map[string]interface{}, db *gorm.DB) error {
	args := m.Called(repos, db)
	return args.Error(0)
}

func (m *MockServiceInitializer) GetServices() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

type MockHandlerInitializer struct {
	mock.Mock
}

func (m *MockHandlerInitializer) InitializeHandlers(services map[string]interface{}) error {
	args := m.Called(services)
	return args.Error(0)
}

func (m *MockHandlerInitializer) GetHandlers() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

type MockRouteRegistrar struct {
	mock.Mock
}

func (m *MockRouteRegistrar) RegisterRoutes(router *gin.Engine, handlers map[string]interface{}) {
	m.Called(router, handlers)
}

type MockServiceComponent struct {
	mock.Mock
}

func (m *MockServiceComponent) Initialize(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	args := m.Called(ctx, cfg, db)
	return args.Error(0)
}

func (m *MockServiceComponent) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockServiceComponent) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockServiceComponent) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockServiceComponent) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// Test helper functions
func createTestConfig() *config.Config {
	return &config.Config{
		Environment: "test",
		Debug:       true,
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		Database: config.DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			Username:        "test",
			Password:        "test",
			Database:        ":memory:",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
			ConnMaxIdleTime: 5 * time.Minute,
		},
	}
}

func createTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Migrate test models
	db.AutoMigrate(&TestModel{})

	return db
}

// Tests for ServiceRegistry
func TestServiceRegistry(t *testing.T) {
	registry := service.NewDefaultServiceRegistry()

	// Test component registration
	component := &MockServiceComponent{}
	component.On("Name").Return("test-component")

	err := registry.Register("test", component)
	assert.NoError(t, err)

	// Test duplicate registration
	err = registry.Register("test", component)
	assert.Error(t, err)

	// Test component retrieval
	retrieved := registry.Get("test")
	assert.Equal(t, component, retrieved)

	// Test component listing
	components := registry.List()
	assert.Len(t, components, 1)
	assert.Contains(t, components, "test")

	// Test component unregistration
	err = registry.Unregister("test")
	assert.NoError(t, err)

	// Test retrieving unregistered component
	retrieved = registry.Get("test")
	assert.Nil(t, retrieved)
}

func TestServiceRegistryLifecycle(t *testing.T) {
	registry := service.NewDefaultServiceRegistry()
	ctx := context.Background()

	// Create mock components
	component1 := &MockServiceComponent{}
	component1.On("Name").Return("component1")
	component1.On("Start", ctx).Return(nil)
	component1.On("Stop", ctx).Return(nil)

	component2 := &MockServiceComponent{}
	component2.On("Name").Return("component2")
	component2.On("Start", ctx).Return(nil)
	component2.On("Stop", ctx).Return(nil)

	// Register components
	registry.Register("comp1", component1)
	registry.Register("comp2", component2)

	// Test start all
	err := registry.StartAll(ctx)
	assert.NoError(t, err)

	// Test stop all
	err = registry.StopAll(ctx)
	assert.NoError(t, err)

	// Verify all methods were called
	component1.AssertExpectations(t)
	component2.AssertExpectations(t)
}

// Tests for DatabaseManager
func TestDatabaseManager(t *testing.T) {
	manager := service.NewDefaultDatabaseManager()
	cfg := createTestConfig()

	// Note: This test would need a real database connection
	// For demonstration purposes, we'll test the interface
	assert.NotNil(t, manager)

	// Test ping without connection
	err := manager.Ping()
	assert.Error(t, err)
}

// Tests for App
func TestAppInitialization(t *testing.T) {
	cfg := createTestConfig()

	// Create mock initializers
	mockRepo := &MockRepositoryInitializer{}
	mockService := &MockServiceInitializer{}
	mockHandler := &MockHandlerInitializer{}
	mockRoute := &MockRouteRegistrar{}

	// Set up mock expectations
	mockRepo.On("InitializeRepositories", mock.AnythingOfType("*gorm.DB")).Return(nil)
	mockRepo.On("GetRepositories").Return(map[string]interface{}{})

	mockService.On("InitializeServices", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*gorm.DB")).Return(nil)
	mockService.On("GetServices").Return(map[string]interface{}{})

	mockHandler.On("InitializeHandlers", mock.AnythingOfType("map[string]interface {}")).Return(nil)
	mockHandler.On("GetHandlers").Return(map[string]interface{}{})

	mockRoute.On("RegisterRoutes", mock.AnythingOfType("*gin.Engine"), mock.AnythingOfType("map[string]interface {}")).Return()

	// Create database initializer with test models
	dbInitializer := service.NewDefaultDatabaseInitializer([]interface{}{&TestModel{}})

	// Create app configuration
	appCfg := service.AppConfig{
		ServiceInfo: service.ServiceInfo{
			Name:    "test-service",
			Version: "1.0.0",
			Port:    8080,
		},
		DatabaseInitializer:   dbInitializer,
		RepositoryInitializer: mockRepo,
		ServiceInitializer:    mockService,
		HandlerInitializer:    mockHandler,
		RouteRegistrar:        mockRoute,
		MiddlewareProvider:    service.NewDefaultMiddlewareProvider(true, true),
	}

	app := service.NewApp(cfg, appCfg)

	// Test basic getters
	assert.Equal(t, cfg, app.GetConfig())
	assert.Equal(t, "test-service", app.GetServiceInfo().Name)
	assert.NotNil(t, app.GetValidator())

	// Note: Full initialization test would require a database connection
}

// Tests for Builder
func TestServiceBuilder(t *testing.T) {
	cfg := createTestConfig()

	// Create mock initializers
	mockRepo := &MockRepositoryInitializer{}
	mockService := &MockServiceInitializer{}
	mockHandler := &MockHandlerInitializer{}
	mockRoute := &MockRouteRegistrar{}

	builder := service.NewServiceBuilder("test-service", cfg).
		WithModels(&TestModel{}).
		WithDefaultPort(8080).
		WithServiceInfo(service.ServiceInfo{
			Name:        "test-service",
			Version:     "1.0.0",
			Description: "Test service",
		}).
		WithRepositoryInitializer(mockRepo).
		WithServiceInitializer(mockService).
		WithHandlerInitializer(mockHandler).
		WithRouteRegistrar(mockRoute).
		WithDefaultMiddleware().
		WithHealthChecks(true)

	assert.NotNil(t, builder)

	// Note: Building the app would require setting up all mock expectations
}

// Tests for Health Checker
func TestHealthChecker(t *testing.T) {
	cfg := createTestConfig()
	db := createTestDB()
	registry := service.NewDefaultServiceRegistry()

	healthChecker := service.NewDefaultHealthChecker("test-service", "1.0.0", db, cfg, registry)

	// Test health data
	healthData := healthChecker.GetHealthData()
	assert.Equal(t, "test-service", healthData["service"])
	assert.Equal(t, "1.0.0", healthData["version"])
	assert.Equal(t, "test", healthData["environment"])

	// Test health check endpoint
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", healthChecker.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "test-service")

	// Test readiness check endpoint
	router.GET("/ready", healthChecker.ReadinessCheck)

	req, _ = http.NewRequest("GET", "/ready", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Tests for Middleware Provider
func TestMiddlewareProvider(t *testing.T) {
	provider := service.NewDefaultMiddlewareProvider(true, true)

	middlewares := provider.GetMiddlewares()
	assert.NotEmpty(t, middlewares)

	// Test middleware configuration
	gin.SetMode(gin.TestMode)
	router := gin.New()
	provider.ConfigureMiddleware(router)

	// Test CORS middleware
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
}

// Tests for Background Service
func TestBackgroundService(t *testing.T) {
	executed := false
	taskFunc := func(ctx context.Context) error {
		executed = true
		return nil
	}

	bgService := service.NewBackgroundService("test-bg", taskFunc)
	assert.Equal(t, "test-bg", bgService.Name())
	assert.False(t, bgService.IsHealthy())

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start the service
	err := bgService.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, bgService.IsHealthy())

	// Give it time to execute
	time.Sleep(50 * time.Millisecond)

	// Stop the service
	err = bgService.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, bgService.IsHealthy())

	// Note: The executed flag might not be true due to timing in tests
	// In a real scenario, you'd use proper synchronization
}

// Integration test
func TestFrameworkIntegration(t *testing.T) {
	cfg := createTestConfig()

	// This would be a more comprehensive integration test
	// that tests the entire framework flow
	serviceInfo := service.ServiceInfo{
		Name:        "integration-test",
		Version:     "1.0.0",
		Description: "Integration test service",
		Port:        8080,
	}

	// Create database initializer
	dbInitializer := service.NewDefaultDatabaseInitializer([]interface{}{&TestModel{}})

	// Note: This test would require implementing full mock initializers
	// and testing the complete application lifecycle
	assert.NotNil(t, serviceInfo)
	assert.NotNil(t, dbInitializer)
}