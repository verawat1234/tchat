package performance

import (
	"time"
)

// BenchmarkConfig holds configuration for performance benchmark tests
type BenchmarkConfig struct {
	Enabled             bool                          `json:"enabled"`
	ConcurrencyLevels   []int                         `json:"concurrency_levels"` // [1, 10, 50, 100, 500]
	DurationSeconds     int                           `json:"duration_seconds"`   // 30 seconds default
	WarmupSeconds       int                           `json:"warmup_seconds"`     // 5 seconds default
	MaxRPS              int                           `json:"max_rps"`            // 10000 requests per second
	MemoryLimitMB       int                           `json:"memory_limit_mb"`    // 512MB default
	CPULimitPercent     float64                       `json:"cpu_limit_percent"`  // 80% default
	ResponseTimeTargets ResponseTimeTargets           `json:"response_time_targets"`
	ThroughputTargets   ThroughputTargets             `json:"throughput_targets"`
	ResourceTargets     ResourceTargets               `json:"resource_targets"`
	ServiceConfigs      map[string]ServiceBenchConfig `json:"service_configs"`
	ReportFormats       []string                      `json:"report_formats"`      // html, json, csv, prometheus
	ComparisonBaseline  string                        `json:"comparison_baseline"` // baseline results file
	ContinuousMode      bool                          `json:"continuous_mode"`     // continuous performance monitoring
}

// ResponseTimeTargets defines response time performance targets
type ResponseTimeTargets struct {
	P50        time.Duration `json:"p50"`         // 50ms median
	P95        time.Duration `json:"p95"`         // 200ms 95th percentile
	P99        time.Duration `json:"p99"`         // 500ms 99th percentile
	Max        time.Duration `json:"max"`         // 2s maximum
}

// ThroughputTargets defines throughput performance targets
type ThroughputTargets struct {
	MinRPS          float64 `json:"min_rps"`           // 100 minimum RPS
	TargetRPS       float64 `json:"target_rps"`        // 2000 target RPS
	MaxRPS          float64 `json:"max_rps"`           // 10000 maximum RPS
	ScalabilityGoal float64 `json:"scalability_goal"`  // 1.2 (20% annual growth)
}

// ResourceTargets defines resource utilization targets
type ResourceTargets struct {
	MaxCPUPercent  float64 `json:"max_cpu_percent"`   // 80%
	MaxMemoryMB    int     `json:"max_memory_mb"`     // 512MB
	MaxDiskIOKBps  int     `json:"max_disk_io_kbps"`  // 10MB/s
	MaxNetworkKBps int     `json:"max_network_kbps"`  // 100MB/s
	MaxOpenFiles   int     `json:"max_open_files"`    // 1000
	MaxGoroutines  int     `json:"max_goroutines"`    // 10000
}

// ServiceBenchConfig defines service-specific benchmark configuration
type ServiceBenchConfig struct {
	Enabled         bool              `json:"enabled"`
	EndpointConfigs map[string]string `json:"endpoint_configs"`
	ConnectionPool  int               `json:"connection_pool"`  // 100 connections
	Timeout         time.Duration     `json:"timeout"`          // 5s timeout
	Retries         int               `json:"retries"`          // 3 retries
}

// BenchmarkMetadata contains metadata about benchmark execution
type BenchmarkMetadata struct {
	Environment   string            `json:"environment"`    // development, staging, production
	Version       string            `json:"version"`        // application version
	GitCommit     string            `json:"git_commit"`     // git commit hash
	TestConfig    string            `json:"test_config"`    // test configuration details
	Infrastructure map[string]string `json:"infrastructure"` // infrastructure details
}

// BenchmarkResult contains results from a performance benchmark test
type BenchmarkResult struct {
	TestName      string                 `json:"test_name"`
	Service       string                 `json:"service"`
	Endpoint      string                 `json:"endpoint,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Duration      time.Duration          `json:"duration"`
	Concurrency   int                    `json:"concurrency"`
	TotalRequests int64                  `json:"total_requests"`
	TotalErrors   int64                  `json:"total_errors"`
	RPS           float64                `json:"rps"`
	ResponseTimes ResponseTimeMetrics    `json:"response_times"`
	Resources     ResourceMetrics        `json:"resources"`
	Errors        map[string]int64       `json:"errors"`
	Status        string                 `json:"status"` // PASS, FAIL, WARNING
	Violations    []PerformanceViolation `json:"violations"`
	Tags          []string               `json:"tags"`
	Metadata      BenchmarkMetadata      `json:"metadata"`
}

// PerformanceViolation represents a performance threshold violation
type PerformanceViolation struct {
	Type        string   `json:"type"`   // response_time, throughput, resource
	Metric      string   `json:"metric"` // p95, rps, cpu_percent
	Expected    float64  `json:"expected"`
	Actual      float64  `json:"actual"`
	Severity    string   `json:"severity"`    // critical, warning, info
	Description string   `json:"description"` // human-readable description
	Suggestions []string `json:"suggestions"` // improvement suggestions
}