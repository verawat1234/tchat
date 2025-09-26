package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"tchat.dev/shared/cache"
	"tchat.dev/shared/database"
	"tchat.dev/shared/external"
	"tchat.dev/shared/messaging"
)

// InfrastructureConfig holds all infrastructure configurations
type InfrastructureConfig struct {
	Database    *InfraDatabaseConfig    `mapstructure:"database"`
	Cache       *CacheConfig       `mapstructure:"cache"`
	Messaging   *MessagingConfig   `mapstructure:"messaging"`
	Storage     *InfraStorageConfig     `mapstructure:"storage"`
	SMS         *InfraSMSConfig         `mapstructure:"sms"`
	Payment     *InfraPaymentConfig     `mapstructure:"payment"`
	Environment string             `mapstructure:"environment"`
	ServiceName string             `mapstructure:"service_name"`
}

// InfraDatabaseConfig holds database configurations
type InfraDatabaseConfig struct {
	Postgres *database.PostgresConfig `mapstructure:"postgres"`
	Scylla   *database.ScyllaConfig   `mapstructure:"scylla"`
}

// CacheConfig holds cache configurations
type CacheConfig struct {
	Redis *database.RedisConfig `mapstructure:"redis"`
}

// MessagingConfig holds messaging configurations
type MessagingConfig struct {
	Kafka *messaging.KafkaConfig `mapstructure:"kafka"`
}

// InfraStorageConfig holds storage configurations
type InfraStorageConfig struct {
	Primary   *external.StorageConfig `mapstructure:"primary"`
	Secondary *external.StorageConfig `mapstructure:"secondary,omitempty"`
}

// InfraSMSConfig holds SMS configurations
type InfraSMSConfig struct {
	Primary   *external.SMSConfig `mapstructure:"primary"`
	Secondary *external.SMSConfig `mapstructure:"secondary,omitempty"`
}

// InfraPaymentConfig holds payment configurations
type InfraPaymentConfig struct {
	Primary   *external.PaymentConfig `mapstructure:"primary"`
	Secondary *external.PaymentConfig `mapstructure:"secondary,omitempty"`
}

// Infrastructure manages all infrastructure services
type Infrastructure struct {
	// Databases
	PostgresDB *database.PostgresDB
	ScyllaDB   *database.ScyllaDB

	// Cache
	Redis        *database.RedisDB
	CacheManager *cache.CacheManager

	// Messaging
	Kafka    *messaging.KafkaClient
	EventBus *messaging.EventBus

	// External Services
	StorageManager *external.StorageManager
	SMSManager     *external.SMSManager
	PaymentManager *external.PaymentGatewayManager

	// Configuration
	Config *InfrastructureConfig
}

// NewInfrastructure creates a new infrastructure instance
func NewInfrastructure(config *InfrastructureConfig) *Infrastructure {
	return &Infrastructure{
		Config: config,
	}
}

// Initialize initializes all infrastructure services
func (i *Infrastructure) Initialize(ctx context.Context) error {
	log.Printf("Initializing infrastructure for service: %s in environment: %s",
		i.Config.ServiceName, i.Config.Environment)

	// Initialize databases
	if err := i.initializeDatabases(ctx); err != nil {
		return fmt.Errorf("failed to initialize databases: %w", err)
	}

	// Initialize cache
	if err := i.initializeCache(ctx); err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Initialize messaging
	if err := i.initializeMessaging(ctx); err != nil {
		return fmt.Errorf("failed to initialize messaging: %w", err)
	}

	// Initialize external services
	if err := i.initializeExternalServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize external services: %w", err)
	}

	log.Println("Infrastructure initialization completed successfully")
	return nil
}

// initializeDatabases initializes database connections
func (i *Infrastructure) initializeDatabases(ctx context.Context) error {
	// Initialize PostgreSQL
	if i.Config.Database.Postgres != nil {
		postgres, err := database.NewPostgresDB(i.Config.Database.Postgres)
		if err != nil {
			return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
		i.PostgresDB = postgres

		// Run health check
		if err := postgres.HealthCheck(ctx); err != nil {
			return fmt.Errorf("PostgreSQL health check failed: %w", err)
		}
	}

	// Initialize ScyllaDB
	if i.Config.Database.Scylla != nil {
		scylla, err := database.NewScyllaDB(i.Config.Database.Scylla)
		if err != nil {
			return fmt.Errorf("failed to connect to ScyllaDB: %w", err)
		}
		i.ScyllaDB = scylla

		// Create keyspace if needed
		if err := scylla.CreateKeyspace(); err != nil {
			log.Printf("Warning: Failed to create keyspace: %v", err)
		}

		// Run health check
		if err := scylla.HealthCheck(ctx); err != nil {
			return fmt.Errorf("ScyllaDB health check failed: %w", err)
		}
	}

	return nil
}

// initializeCache initializes cache services
func (i *Infrastructure) initializeCache(ctx context.Context) error {
	// Initialize Redis
	if i.Config.Cache.Redis != nil {
		redis, err := database.NewRedisDB(i.Config.Cache.Redis)
		if err != nil {
			return fmt.Errorf("failed to connect to Redis: %w", err)
		}
		i.Redis = redis

		// Initialize cache manager
		prefix := fmt.Sprintf("%s:%s", i.Config.ServiceName, i.Config.Environment)
		i.CacheManager = cache.NewCacheManager(redis, prefix)

		// Run health check
		if err := redis.HealthCheck(ctx); err != nil {
			return fmt.Errorf("Redis health check failed: %w", err)
		}
	}

	return nil
}

// initializeMessaging initializes messaging services
func (i *Infrastructure) initializeMessaging(ctx context.Context) error {
	// Initialize Kafka
	if i.Config.Messaging.Kafka != nil {
		kafka, err := messaging.NewKafkaClient(i.Config.Messaging.Kafka)
		if err != nil {
			return fmt.Errorf("failed to connect to Kafka: %w", err)
		}
		i.Kafka = kafka

		// Initialize event bus
		i.EventBus = messaging.NewEventBus(kafka)

		// Create topics
		if err := i.EventBus.CreateAllTopics(ctx); err != nil {
			log.Printf("Warning: Failed to create Kafka topics: %v", err)
		}

		// Run health check
		if err := kafka.HealthCheck(ctx); err != nil {
			return fmt.Errorf("Kafka health check failed: %w", err)
		}
	}

	return nil
}

// initializeExternalServices initializes external services
func (i *Infrastructure) initializeExternalServices(ctx context.Context) error {
	// Initialize storage services
	if i.Config.Storage.Primary != nil {
		i.StorageManager = external.NewStorageManager(i.Config.Storage.Primary)
		// TODO: Register specific storage providers based on configuration
	}

	// Initialize SMS services
	if i.Config.SMS.Primary != nil {
		i.SMSManager = external.NewSMSManager(i.Config.SMS.Primary)
		// TODO: Register specific SMS providers based on configuration
	}

	// Initialize payment services
	if i.Config.Payment.Primary != nil {
		i.PaymentManager = external.NewPaymentGatewayManager(i.Config.Payment.Primary)
		// TODO: Register specific payment gateways based on configuration
	}

	return nil
}

// HealthCheck performs health checks on all infrastructure services
func (i *Infrastructure) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)

	// Database health checks
	if i.PostgresDB != nil {
		results["postgres"] = i.PostgresDB.HealthCheck(ctx)
	}
	if i.ScyllaDB != nil {
		results["scylla"] = i.ScyllaDB.HealthCheck(ctx)
	}

	// Cache health checks
	if i.Redis != nil {
		results["redis"] = i.Redis.HealthCheck(ctx)
	}

	// Messaging health checks
	if i.Kafka != nil {
		results["kafka"] = i.Kafka.HealthCheck(ctx)
	}

	// External service health checks
	if i.StorageManager != nil {
		results["storage"] = i.StorageManager.HealthCheck(ctx)
	}
	if i.SMSManager != nil {
		results["sms"] = i.SMSManager.HealthCheck(ctx)
	}
	if i.PaymentManager != nil {
		results["payment"] = i.PaymentManager.HealthCheck(ctx)
	}

	return results
}

// GetStats returns statistics from all infrastructure services
func (i *Infrastructure) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// Database stats
	if i.PostgresDB != nil {
		stats["postgres"] = i.PostgresDB.GetStats()
	}
	if i.ScyllaDB != nil {
		// GetMetrics is disabled in ScyllaDB implementation
		stats["scylla"] = map[string]interface{}{
			"status": "connected",
			"metrics": "unavailable - gocql.Metrics not supported",
		}
	}

	// Cache stats
	if i.Redis != nil {
		stats["redis"] = i.Redis.GetStats()
	}

	// External service stats
	if i.StorageManager != nil {
		stats["storage"] = i.StorageManager.GetStats()
	}
	if i.SMSManager != nil {
		stats["sms"] = i.SMSManager.GetStats()
	}
	if i.PaymentManager != nil {
		stats["payment"] = i.PaymentManager.GetStats()
	}

	return stats
}

// Shutdown gracefully shuts down all infrastructure services
func (i *Infrastructure) Shutdown(ctx context.Context) error {
	log.Println("Shutting down infrastructure services...")

	var errors []error

	// Stop event bus
	if i.EventBus != nil {
		if err := i.EventBus.Stop(); err != nil {
			errors = append(errors, fmt.Errorf("failed to stop event bus: %w", err))
		}
	}

	// Close Kafka
	if i.Kafka != nil {
		if err := i.Kafka.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Kafka: %w", err))
		}
	}

	// Close Redis
	if i.Redis != nil {
		if err := i.Redis.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Redis: %w", err))
		}
	}

	// Close ScyllaDB
	if i.ScyllaDB != nil {
		i.ScyllaDB.Close()
	}

	// Close PostgreSQL
	if i.PostgresDB != nil {
		if err := i.PostgresDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close PostgreSQL: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("shutdown errors: %v", errors)
	}

	log.Println("Infrastructure shutdown completed")
	return nil
}

// RunMigrations runs database migrations
func (i *Infrastructure) RunMigrations(ctx context.Context, migrationsPath string) error {
	if i.PostgresDB != nil {
		if err := i.PostgresDB.RunMigrations(migrationsPath); err != nil {
			return fmt.Errorf("failed to run PostgreSQL migrations: %w", err)
		}
	}

	// ScyllaDB migrations are handled by the messaging service's migration manager

	return nil
}

// StartEventBus starts the event bus for consuming events
func (i *Infrastructure) StartEventBus() error {
	if i.EventBus != nil {
		return i.EventBus.Start()
	}
	return nil
}

// Default configurations for different environments

// DefaultDevelopmentConfig returns default development configuration
func DefaultDevelopmentConfig(serviceName string) *InfrastructureConfig {
	return &InfrastructureConfig{
		ServiceName: serviceName,
		Environment: "development",
		Database: &InfraDatabaseConfig{
			Postgres: &database.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "tchat_dev",
				Password: "dev_password",
				Database: fmt.Sprintf("tchat_%s_dev", serviceName),
				SSLMode:  "disable",
			},
			Scylla: &database.ScyllaConfig{
				Hosts:    []string{"127.0.0.1:9042"},
				Keyspace: fmt.Sprintf("tchat_%s_dev", serviceName),
			},
		},
		Cache: &CacheConfig{
			Redis: &database.RedisConfig{
				Host:     "localhost",
				Port:     6379,
				Database: 0,
			},
		},
		Messaging: &MessagingConfig{
			Kafka: &messaging.KafkaConfig{
				Brokers:     []string{"localhost:9092"},
				GroupID:     fmt.Sprintf("tchat-%s-dev", serviceName),
				TopicPrefix: "tchat-dev",
			},
		},
		Storage: &InfraStorageConfig{
			Primary: &external.StorageConfig{
				Provider:    external.LocalProvider,
				MaxFileSize: 50 * 1024 * 1024, // 50MB
			},
		},
		SMS: &InfraSMSConfig{
			Primary: &external.SMSConfig{
				Provider: external.TwilioProvider,
			},
		},
		Payment: &InfraPaymentConfig{
			Primary: &external.PaymentConfig{
				Gateway:     external.StripeGateway,
				Environment: "sandbox",
			},
		},
	}
}

// DefaultProductionConfig returns default production configuration
func DefaultProductionConfig(serviceName string) *InfrastructureConfig {
	return &InfrastructureConfig{
		ServiceName: serviceName,
		Environment: "production",
		Database: &InfraDatabaseConfig{
			Postgres: &database.PostgresConfig{
				Host:            "postgres.tchat.prod",
				Port:            5432,
				Database:        fmt.Sprintf("tchat_%s", serviceName),
				SSLMode:         "require",
				MaxIdleConns:    10,
				MaxOpenConns:    100,
				ConnMaxLifetime: time.Hour,
			},
			Scylla: &database.ScyllaConfig{
				Hosts:             []string{"scylla-1.tchat.prod:9042", "scylla-2.tchat.prod:9042", "scylla-3.tchat.prod:9042"},
				Keyspace:          fmt.Sprintf("tchat_%s", serviceName),
				ReplicationFactor: 3,
				Consistency:       "quorum",
			},
		},
		Cache: &CacheConfig{
			Redis: &database.RedisConfig{
				Host:         "redis.tchat.prod",
				Port:         6379,
				PoolSize:     50,
				MinIdleConns: 10,
				EnableTLS:    true,
			},
		},
		Messaging: &MessagingConfig{
			Kafka: &messaging.KafkaConfig{
				Brokers:          []string{"kafka-1.tchat.prod:9092", "kafka-2.tchat.prod:9092", "kafka-3.tchat.prod:9092"},
				GroupID:          fmt.Sprintf("tchat-%s", serviceName),
				TopicPrefix:      "tchat",
				SecurityProtocol: "sasl_ssl",
				SASLMechanism:    "scram-sha-256",
			},
		},
		Storage: &InfraStorageConfig{
			Primary: &external.StorageConfig{
				Provider:    external.AWSS3Provider,
				Region:      "ap-southeast-1",
				MaxFileSize: 100 * 1024 * 1024, // 100MB
				UseSSL:      true,
				PublicRead:  true,
			},
		},
		SMS: &InfraSMSConfig{
			Primary: &external.SMSConfig{
				Provider: external.TwilioProvider,
			},
		},
		Payment: &InfraPaymentConfig{
			Primary: &external.PaymentConfig{
				Gateway:     external.OmiseGateway,
				Environment: "production",
			},
		},
	}
}