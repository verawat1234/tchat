package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db            *PostgresDB
	migrationsDir string
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *PostgresDB, migrationsDir string) *MigrationManager {
	return &MigrationManager{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// CreateMigration creates a new migration file
func (m *MigrationManager) CreateMigration(name string) error {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s", timestamp, name)

	upFile := filepath.Join(m.migrationsDir, fmt.Sprintf("%s.up.sql", filename))
	downFile := filepath.Join(m.migrationsDir, fmt.Sprintf("%s.down.sql", filename))

	// Ensure migrations directory exists
	if err := os.MkdirAll(m.migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create up migration file
	upTemplate := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Add your up migration here\n",
		name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(upFile, []byte(upTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	// Create down migration file
	downTemplate := fmt.Sprintf("-- Rollback migration: %s\n-- Created at: %s\n\n-- Add your down migration here\n",
		name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(downFile, []byte(downTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	log.Printf("Created migration files:\n  %s\n  %s", upFile, downFile)
	return nil
}

// Up runs all pending migrations
func (m *MigrationManager) Up() error {
	migration, err := m.createMigrate()
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run up migrations: %w", err)
	}

	log.Println("Successfully applied all pending migrations")
	return nil
}

// Down rolls back all migrations
func (m *MigrationManager) Down() error {
	migration, err := m.createMigrate()
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run down migrations: %w", err)
	}

	log.Println("Successfully rolled back all migrations")
	return nil
}

// Steps runs n migration steps (positive for up, negative for down)
func (m *MigrationManager) Steps(n int) error {
	migration, err := m.createMigrate()
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Steps(n); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run %d migration steps: %w", n, err)
	}

	direction := "up"
	if n < 0 {
		direction = "down"
		n = -n
	}

	log.Printf("Successfully ran %d migration steps %s", n, direction)
	return nil
}

// Version returns the current migration version
func (m *MigrationManager) Version() (uint, bool, error) {
	migration, err := m.createMigrate()
	if err != nil {
		return 0, false, err
	}
	defer migration.Close()

	return migration.Version()
}

// Force sets the migration version without running migrations
func (m *MigrationManager) Force(version int) error {
	migration, err := m.createMigrate()
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version to %d: %w", version, err)
	}

	log.Printf("Successfully forced migration version to %d", version)
	return nil
}

// Drop drops all tables and removes migration history
func (m *MigrationManager) Drop() error {
	migration, err := m.createMigrate()
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Drop(); err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	log.Println("Successfully dropped all tables and migration history")
	return nil
}

// createMigrate creates a new migrate instance
func (m *MigrationManager) createMigrate() (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(m.db.DB.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.migrationsDir),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return migration, nil
}

// SeedDatabase runs initial data seeding
func (m *MigrationManager) SeedDatabase(ctx context.Context) error {
	seedsDir := filepath.Join(filepath.Dir(m.migrationsDir), "seeds")

	// Check if seeds directory exists
	if _, err := os.Stat(seedsDir); os.IsNotExist(err) {
		log.Println("No seeds directory found, skipping database seeding")
		return nil
	}

	files, err := filepath.Glob(filepath.Join(seedsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to find seed files: %w", err)
	}

	if len(files) == 0 {
		log.Println("No seed files found, skipping database seeding")
		return nil
	}

	return m.db.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		for _, file := range files {
			log.Printf("Running seed file: %s", filepath.Base(file))

			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read seed file %s: %w", file, err)
			}

			if _, err := tx.ExecContext(ctx, string(content)); err != nil {
				return fmt.Errorf("failed to execute seed file %s: %w", file, err)
			}
		}

		log.Printf("Successfully ran %d seed files", len(files))
		return nil
	})
}

// CheckMigrationStatus checks if there are pending migrations
func (m *MigrationManager) CheckMigrationStatus() (bool, error) {
	migration, err := m.createMigrate()
	if err != nil {
		return false, err
	}
	defer migration.Close()

	currentVersion, dirty, err := migration.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return false, fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		return false, fmt.Errorf("database is in dirty state at version %d", currentVersion)
	}

	// Try to step up one migration to check if there are pending ones
	tempMigrate, err := m.createMigrate()
	if err != nil {
		return false, err
	}
	defer tempMigrate.Close()

	if err := tempMigrate.Steps(1); err != nil {
		if err == migrate.ErrNoChange {
			return false, nil // No pending migrations
		}
		return false, fmt.Errorf("failed to check for pending migrations: %w", err)
	}

	// Roll back the test migration
	if err := tempMigrate.Steps(-1); err != nil {
		return false, fmt.Errorf("failed to rollback test migration: %w", err)
	}

	return true, nil // There are pending migrations
}