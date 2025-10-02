package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tchat.dev/shared/config"
)

// DefaultDatabaseManager implements DatabaseManager
type DefaultDatabaseManager struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

// NewDefaultDatabaseManager creates a new database manager
func NewDefaultDatabaseManager() DatabaseManager {
	return &DefaultDatabaseManager{}
}

// Connect establishes a database connection
func (m *DefaultDatabaseManager) Connect(cfg *config.Config) (*gorm.DB, *sql.DB, error) {
	// Create database connection
	dsn := cfg.GetDatabaseURL()

	var gormLogger logger.Interface
	if cfg.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Store references
	m.db = db
	m.sqlDB = sqlDB

	return db, sqlDB, nil
}

// ConfigureConnection configures the database connection pool
func (m *DefaultDatabaseManager) ConfigureConnection(db *sql.DB, cfg *config.Config) error {
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

	return nil
}

// Migrate runs database migrations
func (m *DefaultDatabaseManager) Migrate(db *gorm.DB, models []interface{}) error {
	if len(models) == 0 {
		return nil
	}

	return db.AutoMigrate(models...)
}

// Close closes the database connection
func (m *DefaultDatabaseManager) Close() error {
	if m.sqlDB != nil {
		return m.sqlDB.Close()
	}
	return nil
}

// Ping checks if the database connection is alive
func (m *DefaultDatabaseManager) Ping() error {
	if m.sqlDB != nil {
		return m.sqlDB.Ping()
	}
	return fmt.Errorf("database connection not established")
}

// GetConnection returns the GORM database instance
func (m *DefaultDatabaseManager) GetConnection() *gorm.DB {
	return m.db
}

// DefaultDatabaseInitializer provides a basic database initializer
type DefaultDatabaseInitializer struct {
	models  []interface{}
	manager DatabaseManager
}

// NewDefaultDatabaseInitializer creates a new database initializer
func NewDefaultDatabaseInitializer(models []interface{}) DatabaseInitializer {
	return &DefaultDatabaseInitializer{
		models:  models,
		manager: NewDefaultDatabaseManager(),
	}
}

// InitializeDatabase initializes the database connection
func (d *DefaultDatabaseInitializer) InitializeDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, sqlDB, err := d.manager.Connect(cfg)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	if err := d.manager.ConfigureConnection(sqlDB, cfg); err != nil {
		return nil, fmt.Errorf("failed to configure connection: %w", err)
	}

	return db, nil
}

// RunMigrations runs database migrations
func (d *DefaultDatabaseInitializer) RunMigrations(db *gorm.DB) error {
	return d.manager.Migrate(db, d.models)
}

// GetModels returns the models to migrate
func (d *DefaultDatabaseInitializer) GetModels() []interface{} {
	return d.models
}

// ServiceDatabaseInitializer allows services to customize database initialization
type ServiceDatabaseInitializer struct {
	*DefaultDatabaseInitializer
	customInitFunc func(*gorm.DB, *config.Config) error
	customModels   []interface{}
}

// NewServiceDatabaseInitializer creates a service-specific database initializer
func NewServiceDatabaseInitializer(models []interface{}, customInit func(*gorm.DB, *config.Config) error) DatabaseInitializer {
	return &ServiceDatabaseInitializer{
		DefaultDatabaseInitializer: NewDefaultDatabaseInitializer(models).(*DefaultDatabaseInitializer),
		customInitFunc:            customInit,
		customModels:              models,
	}
}

// InitializeDatabase initializes the database with custom logic
func (s *ServiceDatabaseInitializer) InitializeDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Use default initialization
	db, err := s.DefaultDatabaseInitializer.InitializeDatabase(cfg)
	if err != nil {
		return nil, err
	}

	// Run custom initialization if provided
	if s.customInitFunc != nil {
		if err := s.customInitFunc(db, cfg); err != nil {
			return nil, fmt.Errorf("custom database initialization failed: %w", err)
		}
	}

	return db, nil
}

// GetModels returns the custom models
func (s *ServiceDatabaseInitializer) GetModels() []interface{} {
	if s.customModels != nil {
		return s.customModels
	}
	return s.DefaultDatabaseInitializer.GetModels()
}

// DatabaseHealthComponent provides database health checking
type DatabaseHealthComponent struct {
	*BaseServiceComponent
	dbManager DatabaseManager
}

// NewDatabaseHealthComponent creates a new database health component
func NewDatabaseHealthComponent(dbManager DatabaseManager) ServiceComponent {
	return &DatabaseHealthComponent{
		BaseServiceComponent: NewBaseServiceComponent("database-health"),
		dbManager:           dbManager,
	}
}

// Initialize initializes the database health component
func (c *DatabaseHealthComponent) Initialize(ctx context.Context, cfg *config.Config, db *gorm.DB) error {
	return c.BaseServiceComponent.Initialize(ctx, cfg, db)
}

// Start starts the database health monitoring
func (c *DatabaseHealthComponent) Start(ctx context.Context) error {
	if err := c.BaseServiceComponent.Start(ctx); err != nil {
		return err
	}

	// Check database health
	if err := c.dbManager.Ping(); err != nil {
		c.SetHealthy(false)
		return fmt.Errorf("database health check failed: %w", err)
	}

	c.SetHealthy(true)
	log.Printf("Database health component started")
	return nil
}

// IsHealthy returns true if database is healthy
func (c *DatabaseHealthComponent) IsHealthy() bool {
	if err := c.dbManager.Ping(); err != nil {
		c.SetHealthy(false)
		return false
	}

	c.SetHealthy(true)
	return c.BaseServiceComponent.IsHealthy()
}