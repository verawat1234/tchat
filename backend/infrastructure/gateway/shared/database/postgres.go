package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresConfig holds PostgreSQL database configuration
type PostgresConfig struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	User            string        `mapstructure:"user" validate:"required"`
	Password        string        `mapstructure:"password" validate:"required"`
	Database        string        `mapstructure:"database" validate:"required"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	Timezone        string        `mapstructure:"timezone"`
}

// PostgresDB wraps sqlx.DB with additional functionality
type PostgresDB struct {
	*sqlx.DB
	Config *PostgresConfig
}

// DefaultPostgresConfig returns default PostgreSQL configuration
func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:            "localhost",
		Port:            5432,
		SSLMode:         "disable",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Minute * 30,
		Timezone:        "UTC",
	}
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(config *PostgresConfig) (*PostgresDB, error) {
	if config == nil {
		config = DefaultPostgresConfig()
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
		config.Timezone,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL database: %s@%s:%d/%s",
		config.User, config.Host, config.Port, config.Database)

	return &PostgresDB{
		DB:     db,
		Config: config,
	}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	if p.DB != nil {
		log.Println("Closing PostgreSQL database connection")
		return p.DB.Close()
	}
	return nil
}

// HealthCheck performs a health check on the database
func (p *PostgresDB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result int
	if err := p.GetContext(ctx, &result, "SELECT 1"); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetStats returns database connection statistics
func (p *PostgresDB) GetStats() sql.DBStats {
	return p.DB.Stats()
}

// RunMigrations runs database migrations from the specified directory
func (p *PostgresDB) RunMigrations(migrationsPath string) error {
	driver, err := postgres.WithInstance(p.DB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Database migrations completed successfully")
	return nil
}

// BeginTx starts a new transaction
func (p *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return p.DB.BeginTxx(ctx, opts)
}

// WithTransaction executes a function within a database transaction
func (p *PostgresDB) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := p.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Failed to rollback transaction: %v (original error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ExecInTransaction executes a query within a transaction
func (p *PostgresDB) ExecInTransaction(ctx context.Context, query string, args ...interface{}) error {
	return p.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	})
}

// GetInTransaction executes a query and scans the result within a transaction
func (p *PostgresDB) GetInTransaction(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, dest, query, args...)
	})
}

// SelectInTransaction executes a query and scans multiple results within a transaction
func (p *PostgresDB) SelectInTransaction(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, dest, query, args...)
	})
}