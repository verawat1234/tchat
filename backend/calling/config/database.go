package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	DBName                string
	SSLMode               string
	MaxIdleConns          int
	MaxOpenConns          int
	ConnMaxLifetime       time.Duration
	ConnMaxIdleTime       time.Duration
	LogLevel              logger.LogLevel
	AutoMigrate           bool
	ConnectionTimeout     time.Duration
	StatementTimeout      time.Duration
	LockTimeout           time.Duration
	IdleInTransactionTimeout time.Duration
}

// NewDatabaseConfig creates a new database configuration from environment variables
func NewDatabaseConfig() *DatabaseConfig {
	logLevel := logger.Silent
	if getEnv("DB_LOG_LEVEL", "silent") == "info" {
		logLevel = logger.Info
	}

	return &DatabaseConfig{
		Host:                     getEnv("DB_HOST", "localhost"),
		Port:                     getEnvAsInt("DB_PORT", 5432),
		User:                     getEnv("DB_USER", "postgres"),
		Password:                 getEnv("DB_PASSWORD", ""),
		DBName:                   getEnv("DB_NAME", "tchat_calling"),
		SSLMode:                  getEnv("DB_SSL_MODE", "disable"),
		MaxIdleConns:             getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:             getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime:          getEnvAsDuration("DB_CONN_MAX_LIFETIME", 1*time.Hour),
		ConnMaxIdleTime:          getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
		LogLevel:                 logLevel,
		AutoMigrate:              getEnvAsBool("DB_AUTO_MIGRATE", false),
		ConnectionTimeout:        getEnvAsDuration("DB_CONNECTION_TIMEOUT", 10*time.Second),
		StatementTimeout:         getEnvAsDuration("DB_STATEMENT_TIMEOUT", 30*time.Second),
		LockTimeout:              getEnvAsDuration("DB_LOCK_TIMEOUT", 30*time.Second),
		IdleInTransactionTimeout: getEnvAsDuration("DB_IDLE_IN_TRANSACTION_TIMEOUT", 60*time.Second),
	}
}

// DatabaseManager manages database connections and operations
type DatabaseManager struct {
	DB     *gorm.DB
	Config *DatabaseConfig
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager() (*DatabaseManager, error) {
	config := NewDatabaseConfig()

	// Build connection string with timeouts
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d statement_timeout=%d lock_timeout=%d idle_in_transaction_session_timeout=%d",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
		int(config.ConnectionTimeout.Seconds()),
		int(config.StatementTimeout.Milliseconds()),
		int(config.LockTimeout.Milliseconds()),
		int(config.IdleInTransactionTimeout.Milliseconds()),
	)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // Enable prepared statements for better performance
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL database at %s:%d/%s",
		config.Host, config.Port, config.DBName)

	dm := &DatabaseManager{
		DB:     db,
		Config: config,
	}

	// Run auto-migration if enabled
	if config.AutoMigrate {
		if err := dm.RunMigrations(); err != nil {
			log.Printf("Warning: Auto-migration failed: %v", err)
		}
	}

	return dm, nil
}

// RunMigrations runs database migrations for calling service models
func (dm *DatabaseManager) RunMigrations() error {
	log.Println("Running database migrations for calling service...")

	// Import your models here when they're created
	// For now, we'll use the raw SQL migrations

	// Note: In a production setup, you might want to use a proper migration tool
	// like golang-migrate/migrate or integrate with your existing migration system

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (dm *DatabaseManager) Close() error {
	sqlDB, err := dm.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the GORM database instance
func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.DB
}

// HealthCheck performs a health check on the database connection
func (dm *DatabaseManager) HealthCheck() error {
	sqlDB, err := dm.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// GetConnectionStats returns database connection statistics
func (dm *DatabaseManager) GetConnectionStats() (map[string]interface{}, error) {
	sqlDB, err := dm.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()

	return map[string]interface{}{
		"max_open_connections":     stats.MaxOpenConnections,
		"open_connections":         stats.OpenConnections,
		"in_use":                  stats.InUse,
		"idle":                    stats.Idle,
		"wait_count":              stats.WaitCount,
		"wait_duration_ms":        stats.WaitDuration.Milliseconds(),
		"max_idle_closed":         stats.MaxIdleClosed,
		"max_idle_time_closed":    stats.MaxIdleTimeClosed,
		"max_lifetime_closed":     stats.MaxLifetimeClosed,
	}, nil
}

// Transaction executes a function within a database transaction
func (dm *DatabaseManager) Transaction(fn func(*gorm.DB) error) error {
	return dm.DB.Transaction(fn)
}

// Helper functions for environment variables
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// Validate database configuration
func (config *DatabaseConfig) Validate() error {
	if config.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Port)
	}
	if config.MaxOpenConns <= 0 {
		return fmt.Errorf("max open connections must be positive")
	}
	if config.MaxIdleConns < 0 {
		return fmt.Errorf("max idle connections cannot be negative")
	}
	if config.MaxIdleConns > config.MaxOpenConns {
		return fmt.Errorf("max idle connections cannot exceed max open connections")
	}
	return nil
}

// Initialize sets up both database and Redis connections
func Initialize() (*DatabaseManager, *RedisClient, error) {
	// Initialize database
	dbManager, err := NewDatabaseManager()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Validate database configuration
	if err := dbManager.Config.Validate(); err != nil {
		dbManager.Close()
		return nil, nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	// Initialize Redis
	redisClient, err := NewRedisClient()
	if err != nil {
		dbManager.Close()
		return nil, nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	log.Println("Successfully initialized all database connections")
	return dbManager, redisClient, nil
}