package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

// ScyllaConfig holds ScyllaDB configuration
type ScyllaConfig struct {
	Hosts              []string      `mapstructure:"hosts" validate:"required"`
	Keyspace           string        `mapstructure:"keyspace" validate:"required"`
	Username           string        `mapstructure:"username"`
	Password           string        `mapstructure:"password"`
	Consistency        string        `mapstructure:"consistency"`
	Timeout            time.Duration `mapstructure:"timeout"`
	ConnectTimeout     time.Duration `mapstructure:"connect_timeout"`
	NumConns           int           `mapstructure:"num_conns"`
	ReplicationFactor  int           `mapstructure:"replication_factor"`
	ReplicationClass   string        `mapstructure:"replication_class"`
	EnableHostVerify   bool          `mapstructure:"enable_host_verify"`
	EnableCompression  bool          `mapstructure:"enable_compression"`
	RetryPolicy        string        `mapstructure:"retry_policy"`
	MaxRetries         int           `mapstructure:"max_retries"`
	RetryBackoff       time.Duration `mapstructure:"retry_backoff"`
}

// ScyllaDB wraps gocql.Session with additional functionality
type ScyllaDB struct {
	Session   *gocql.Session
	Config    *ScyllaConfig
	cluster   *gocql.ClusterConfig
}

// DefaultScyllaConfig returns default ScyllaDB configuration
func DefaultScyllaConfig() *ScyllaConfig {
	return &ScyllaConfig{
		Hosts:             []string{"127.0.0.1:9042"},
		Keyspace:          "tchat",
		Consistency:       "quorum",
		Timeout:           5 * time.Second,
		ConnectTimeout:    10 * time.Second,
		NumConns:          2,
		ReplicationFactor: 1,
		ReplicationClass:  "SimpleStrategy",
		EnableHostVerify:  false,
		EnableCompression: true,
		RetryPolicy:       "exponential",
		MaxRetries:        3,
		RetryBackoff:      time.Second,
	}
}

// NewScyllaDB creates a new ScyllaDB connection
func NewScyllaDB(config *ScyllaConfig) (*ScyllaDB, error) {
	if config == nil {
		config = DefaultScyllaConfig()
	}

	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Keyspace = config.Keyspace
	cluster.Timeout = config.Timeout
	cluster.ConnectTimeout = config.ConnectTimeout
	cluster.NumConns = config.NumConns
	cluster.DisableInitialHostLookup = !config.EnableHostVerify

	// Set consistency level
	consistency, err := parseConsistency(config.Consistency)
	if err != nil {
		return nil, fmt.Errorf("invalid consistency level %s: %w", config.Consistency, err)
	}
	cluster.Consistency = consistency

	// Set authentication if provided
	if config.Username != "" && config.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: config.Username,
			Password: config.Password,
		}
	}

	// Set compression
	if config.EnableCompression {
		cluster.Compressor = &gocql.SnappyCompressor{}
	}

	// Set retry policy
	switch config.RetryPolicy {
	case "exponential":
		cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
			NumRetries: config.MaxRetries,
			Min:        config.RetryBackoff,
			Max:        config.RetryBackoff * 10,
		}
	case "simple":
		cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
			NumRetries: config.MaxRetries,
		}
	default:
		cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
			NumRetries: 3,
			Min:        time.Second,
			Max:        10 * time.Second,
		}
	}

	// Create session
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create ScyllaDB session: %w", err)
	}

	log.Printf("Successfully connected to ScyllaDB cluster: %v, keyspace: %s",
		config.Hosts, config.Keyspace)

	return &ScyllaDB{
		Session: session,
		Config:  config,
		cluster: cluster,
	}, nil
}

// Close closes the ScyllaDB session
func (s *ScyllaDB) Close() {
	if s.Session != nil {
		log.Println("Closing ScyllaDB session")
		s.Session.Close()
	}
}

// HealthCheck performs a health check on the database
func (s *ScyllaDB) HealthCheck(ctx context.Context) error {
	query := s.Session.Query("SELECT now() FROM system.local")

	done := make(chan error, 1)
	go func() {
		var timestamp time.Time
		done <- query.Scan(&timestamp)
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("ScyllaDB health check failed: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("ScyllaDB health check timeout: %w", ctx.Err())
	}
}

// CreateKeyspace creates the keyspace if it doesn't exist
func (s *ScyllaDB) CreateKeyspace() error {
	// Connect without keyspace to create it
	cluster := gocql.NewCluster(s.Config.Hosts...)
	cluster.Timeout = s.Config.ConnectTimeout
	cluster.DisableInitialHostLookup = !s.Config.EnableHostVerify

	if s.Config.Username != "" && s.Config.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: s.Config.Username,
			Password: s.Config.Password,
		}
	}

	tempSession, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create temporary session: %w", err)
	}
	defer tempSession.Close()

	query := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
		WITH REPLICATION = {
			'class': '%s',
			'replication_factor': %d
		}`,
		s.Config.Keyspace,
		s.Config.ReplicationClass,
		s.Config.ReplicationFactor,
	)

	if err := tempSession.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace %s: %w", s.Config.Keyspace, err)
	}

	log.Printf("Created keyspace: %s", s.Config.Keyspace)
	return nil
}

// ExecuteQuery executes a CQL query
func (s *ScyllaDB) ExecuteQuery(query string, values ...interface{}) error {
	return s.Session.Query(query, values...).Exec()
}

// QueryRow executes a query and scans the first row
func (s *ScyllaDB) QueryRow(dest []interface{}, query string, values ...interface{}) error {
	return s.Session.Query(query, values...).Scan(dest...)
}

// QueryRows executes a query and returns an iterator
func (s *ScyllaDB) QueryRows(query string, values ...interface{}) *gocql.Iter {
	return s.Session.Query(query, values...).Iter()
}

// ExecuteBatch executes a batch of queries
func (s *ScyllaDB) ExecuteBatch(batch *gocql.Batch) error {
	return s.Session.ExecuteBatch(batch)
}

// CreateBatch creates a new batch query
func (s *ScyllaDB) CreateBatch(batchType gocql.BatchType) *gocql.Batch {
	return s.Session.NewBatch(batchType)
}

// RunMigrations runs CQL migration scripts
func (s *ScyllaDB) RunMigrations(migrationQueries []string) error {
	for i, query := range migrationQueries {
		if err := s.ExecuteQuery(query); err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", i+1, err)
		}
		log.Printf("Executed migration %d", i+1)
	}

	log.Printf("Successfully executed %d migrations", len(migrationQueries))
	return nil
}

// GetMetrics returns session metrics (disabled - gocql.Metrics not available)
// func (s *ScyllaDB) GetMetrics() gocql.Metrics {
//	return s.Session.Metrics()
// }

// parseConsistency converts string to gocql.Consistency
func parseConsistency(consistency string) (gocql.Consistency, error) {
	switch consistency {
	case "any":
		return gocql.Any, nil
	case "one":
		return gocql.One, nil
	case "two":
		return gocql.Two, nil
	case "three":
		return gocql.Three, nil
	case "quorum":
		return gocql.Quorum, nil
	case "all":
		return gocql.All, nil
	case "local_quorum":
		return gocql.LocalQuorum, nil
	case "each_quorum":
		return gocql.EachQuorum, nil
	case "local_one":
		return gocql.LocalOne, nil
	default:
		return gocql.Quorum, fmt.Errorf("unknown consistency level: %s", consistency)
	}
}