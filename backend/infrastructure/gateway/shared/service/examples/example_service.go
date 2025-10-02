package examples

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"tchat.dev/shared/config"
	"tchat.dev/shared/service"
)

// Example models for demonstration
type ExampleModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255;not null"`
}

type ExampleRepository interface {
	GetAll() ([]ExampleModel, error)
	GetByID(id uint) (*ExampleModel, error)
	Create(model *ExampleModel) error
}

type ExampleService struct {
	repo ExampleRepository
}

type ExampleHandler struct {
	service *ExampleService
}

// Example repository implementation
type postgresExampleRepository struct {
	db *gorm.DB
}

func NewPostgresExampleRepository(db *gorm.DB) ExampleRepository {
	return &postgresExampleRepository{db: db}
}

func (r *postgresExampleRepository) GetAll() ([]ExampleModel, error) {
	var models []ExampleModel
	err := r.db.Find(&models).Error
	return models, err
}

func (r *postgresExampleRepository) GetByID(id uint) (*ExampleModel, error) {
	var model ExampleModel
	err := r.db.First(&model, id).Error
	return &model, err
}

func (r *postgresExampleRepository) Create(model *ExampleModel) error {
	return r.db.Create(model).Error
}

// Example service implementation
func NewExampleService(repo ExampleRepository) *ExampleService {
	return &ExampleService{repo: repo}
}

func (s *ExampleService) GetAll() ([]ExampleModel, error) {
	return s.repo.GetAll()
}

func (s *ExampleService) Create(model *ExampleModel) error {
	return s.repo.Create(model)
}

// Example handler implementation
func NewExampleHandler(service *ExampleService) *ExampleHandler {
	return &ExampleHandler{service: service}
}

func (h *ExampleHandler) GetAll(c *gin.Context) {
	models, err := h.service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, models)
}

func (h *ExampleHandler) Create(c *gin.Context) {
	var model ExampleModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&model); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, model)
}

// Example service initializers using the framework

// ExampleRepositoryInitializer implements service.RepositoryInitializer
type ExampleRepositoryInitializer struct {
	repositories map[string]interface{}
}

func NewExampleRepositoryInitializer() service.RepositoryInitializer {
	return &ExampleRepositoryInitializer{
		repositories: make(map[string]interface{}),
	}
}

func (r *ExampleRepositoryInitializer) InitializeRepositories(db *gorm.DB) error {
	exampleRepo := NewPostgresExampleRepository(db)
	r.repositories["example"] = exampleRepo
	return nil
}

func (r *ExampleRepositoryInitializer) GetRepositories() map[string]interface{} {
	return r.repositories
}

// ExampleServiceInitializer implements service.ServiceInitializer
type ExampleServiceInitializer struct {
	services map[string]interface{}
}

func NewExampleServiceInitializer() service.ServiceInitializer {
	return &ExampleServiceInitializer{
		services: make(map[string]interface{}),
	}
}

func (s *ExampleServiceInitializer) InitializeServices(repos map[string]interface{}, db *gorm.DB) error {
	exampleRepo, ok := repos["example"].(ExampleRepository)
	if !ok {
		return fmt.Errorf("example repository not found or invalid type")
	}

	exampleService := NewExampleService(exampleRepo)
	s.services["example"] = exampleService
	return nil
}

func (s *ExampleServiceInitializer) GetServices() map[string]interface{} {
	return s.services
}

// ExampleHandlerInitializer implements service.HandlerInitializer
type ExampleHandlerInitializer struct {
	handlers map[string]interface{}
}

func NewExampleHandlerInitializer() service.HandlerInitializer {
	return &ExampleHandlerInitializer{
		handlers: make(map[string]interface{}),
	}
}

func (h *ExampleHandlerInitializer) InitializeHandlers(services map[string]interface{}) error {
	exampleService, ok := services["example"].(*ExampleService)
	if !ok {
		return fmt.Errorf("example service not found or invalid type")
	}

	exampleHandler := NewExampleHandler(exampleService)
	h.handlers["example"] = exampleHandler
	return nil
}

func (h *ExampleHandlerInitializer) GetHandlers() map[string]interface{} {
	return h.handlers
}

// ExampleRouteRegistrar implements service.RouteRegistrar
type ExampleRouteRegistrar struct{}

func NewExampleRouteRegistrar() service.RouteRegistrar {
	return &ExampleRouteRegistrar{}
}

func (r *ExampleRouteRegistrar) RegisterRoutes(router *gin.Engine, handlers map[string]interface{}) {
	exampleHandler, ok := handlers["example"].(*ExampleHandler)
	if !ok {
		log.Printf("Warning: example handler not found or invalid type")
		return
	}

	v1 := router.Group("/api/v1")
	{
		examples := v1.Group("/examples")
		{
			examples.GET("", exampleHandler.GetAll)
			examples.POST("", exampleHandler.Create)
		}
	}
}

// Example of how to create a service using the framework
func CreateExampleService() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set service port
	cfg.Server.Port = 8092

	// Create the service using the builder pattern
	app, err := service.NewServiceBuilder("example-service", cfg).
		WithModels(&ExampleModel{}).
		WithDefaultPort(8092).
		WithServiceInfo(service.ServiceInfo{
			Name:        "example-service",
			Version:     "1.0.0",
			Description: "Example service using the framework",
			Port:        8092,
		}).
		WithRepositoryInitializer(NewExampleRepositoryInitializer()).
		WithServiceInitializer(NewExampleServiceInitializer()).
		WithHandlerInitializer(NewExampleHandlerInitializer()).
		WithRouteRegistrar(NewExampleRouteRegistrar()).
		WithConfigurableMiddleware().
		WithHealthChecks(true).
		Build()

	if err != nil {
		log.Fatalf("Failed to build application: %v", err)
	}

	// Initialize and run
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.RunWithGracefulShutdown(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

// Alternative way using the fluent builder
func CreateExampleServiceWithFluentBuilder() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg.Server.Port = 8092

	// Use the fluent builder to create and run the service in one go
	err = service.NewServiceBuilder("example-service", cfg).
		WithModels(&ExampleModel{}).
		WithServiceInfo(service.ServiceInfo{
			Name:        "example-service",
			Version:     "1.0.0",
			Description: "Example service using fluent builder",
			Port:        8092,
		}).
		WithRepositoryInitializer(NewExampleRepositoryInitializer()).
		WithServiceInitializer(NewExampleServiceInitializer()).
		WithHandlerInitializer(NewExampleHandlerInitializer()).
		WithRouteRegistrar(NewExampleRouteRegistrar()).
		WithDefaultMiddleware().
		WithHealthChecks(true).
		BuildAndRun()

	if err != nil {
		log.Fatalf("Service failed: %v", err)
	}
}