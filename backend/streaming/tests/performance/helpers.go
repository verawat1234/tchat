package performance

import (
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
)

// ChatPercentile calculates the specified percentile for time.Duration slices
func ChatPercentile(durations []time.Duration, percentile int) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Sort durations
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)

	// Simple bubble sort (sufficient for test data)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	index := int(float64(len(sorted)-1) * float64(percentile) / 100.0)
	return sorted[index]
}

// ChatAverage calculates the average latency
func ChatAverage(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var sum time.Duration
	for _, latency := range durations {
		sum += latency
	}
	return sum / time.Duration(len(durations))
}

// TestConfig holds configuration for performance tests
type TestConfig struct {
	ScyllaDBHost     string
	ScyllaDBKeyspace string
	NumWriters       int
	MessagesPerSec   int
	TestDuration     time.Duration
}

// LoadTestConfig loads configuration from environment with defaults
func LoadTestConfig() *TestConfig {
	config := &TestConfig{
		ScyllaDBHost:     getEnv("SCYLLADB_HOST", "localhost:9042"),
		ScyllaDBKeyspace: getEnv("SCYLLADB_KEYSPACE", "tchat"),
		NumWriters:       getEnvInt("NUM_WRITERS", 100),
		MessagesPerSec:   getEnvInt("MESSAGES_PER_SEC", 1000),
		TestDuration:     getEnvDuration("TEST_DURATION", 60*time.Second),
	}
	return config
}

// SetupChatScyllaDB creates a ScyllaDB session for testing
func SetupChatScyllaDB(config *TestConfig) (*gocql.Session, error) {
	cluster := gocql.NewCluster(config.ScyllaDBHost)
	cluster.Keyspace = config.ScyllaDBKeyspace
	cluster.Consistency = gocql.LocalQuorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 5 * time.Second
	cluster.NumConns = 4
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	createTableCQL := `
		CREATE TABLE IF NOT EXISTS chat_messages (
			stream_id uuid,
			timestamp timestamp,
			message_id uuid,
			sender_id uuid,
			sender_display_name text,
			message_text text,
			moderation_status text,
			message_type text,
			PRIMARY KEY ((stream_id), timestamp, message_id)
		) WITH CLUSTERING ORDER BY (timestamp DESC)
		  AND default_time_to_live = 2592000
		  AND compaction = {'class': 'TimeWindowCompactionStrategy'}
	`
	if err := session.Query(createTableCQL).Exec(); err != nil {
		session.Close()
		return nil, err
	}

	return session, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}