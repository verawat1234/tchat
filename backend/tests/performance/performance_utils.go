package performance

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// PerformanceAnalyzer provides comprehensive performance analysis utilities
type PerformanceAnalyzer struct {
	config    BenchmarkConfig
	results   []BenchmarkResult
	baseline  *BaselineMetrics
	mu        sync.RWMutex
	reportDir string
	sessionID string
}

// BaselineMetrics represents baseline performance metrics for comparison
type BaselineMetrics struct {
	Version     string                     `json:"version"`
	Timestamp   time.Time                  `json:"timestamp"`
	Services    map[string]ServiceBaseline `json:"services"`
	Environment string                     `json:"environment"`
	TestSuite   string                     `json:"test_suite"`
}

// ServiceBaseline contains baseline metrics for a service
type ServiceBaseline struct {
	ServiceName  string                      `json:"service_name"`
	Endpoints    map[string]EndpointBaseline `json:"endpoints"`
	Database     DatabaseBaseline            `json:"database"`
	Cache        CacheBaseline               `json:"cache"`
	MessageQueue MessageQueueBaseline        `json:"message_queue"`
}

// EndpointBaseline contains baseline metrics for an endpoint
type EndpointBaseline struct {
	Path          string        `json:"path"`
	Method        string        `json:"method"`
	ExpectedRPS   float64       `json:"expected_rps"`
	ExpectedP95   time.Duration `json:"expected_p95"`
	ExpectedP99   time.Duration `json:"expected_p99"`
	MaxCPUPercent float64       `json:"max_cpu_percent"`
	MaxMemoryMB   float64       `json:"max_memory_mb"`
	MaxErrorRate  float64       `json:"max_error_rate"`
	Tags          []string      `json:"tags"`
}

// DatabaseBaseline contains baseline metrics for database operations
type DatabaseBaseline struct {
	Operations     map[string]OperationBaseline `json:"operations"`
	ConnectionPool int                          `json:"connection_pool"`
	MaxLatency     time.Duration                `json:"max_latency"`
}

// CacheBaseline contains baseline metrics for cache operations
type CacheBaseline struct {
	Operations map[string]OperationBaseline `json:"operations"`
	HitRate    float64                      `json:"hit_rate"`
	MaxLatency time.Duration                `json:"max_latency"`
}

// MessageQueueBaseline contains baseline metrics for message queue operations
type MessageQueueBaseline struct {
	ProducerRPS   float64       `json:"producer_rps"`
	ConsumerRPS   float64       `json:"consumer_rps"`
	MaxLatency    time.Duration `json:"max_latency"`
	MaxQueueDepth int           `json:"max_queue_depth"`
}

// OperationBaseline contains baseline metrics for individual operations
type OperationBaseline struct {
	Operation    string        `json:"operation"`
	ExpectedRPS  float64       `json:"expected_rps"`
	ExpectedP95  time.Duration `json:"expected_p95"`
	MaxErrorRate float64       `json:"max_error_rate"`
}

// PerformanceRegression represents a performance regression detection
type PerformanceRegression struct {
	TestName        string    `json:"test_name"`
	Service         string    `json:"service"`
	Metric          string    `json:"metric"`
	Current         float64   `json:"current"`
	Baseline        float64   `json:"baseline"`
	RegressionPct   float64   `json:"regression_pct"`
	Threshold       float64   `json:"threshold"`
	Severity        string    `json:"severity"`
	Impact          string    `json:"impact"`
	Recommendations []string  `json:"recommendations"`
	Timestamp       time.Time `json:"timestamp"`
}

// PerformanceTrend represents performance trend analysis
type PerformanceTrend struct {
	TestName    string             `json:"test_name"`
	Service     string             `json:"service"`
	Metric      string             `json:"metric"`
	Timeframe   string             `json:"timeframe"`
	DataPoints  []TrendDataPoint   `json:"data_points"`
	TrendLine   TrendAnalysis      `json:"trend_line"`
	Predictions []PredictionPoint  `json:"predictions"`
	Anomalies   []AnomalyDetection `json:"anomalies"`
}

// TrendDataPoint represents a single data point in trend analysis
type TrendDataPoint struct {
	Timestamp time.Time         `json:"timestamp"`
	Value     float64           `json:"value"`
	Metadata  map[string]string `json:"metadata"`
}

// TrendAnalysis contains trend analysis results
type TrendAnalysis struct {
	Slope        float64 `json:"slope"`
	Intercept    float64 `json:"intercept"`
	Correlation  float64 `json:"correlation"`
	Confidence   float64 `json:"confidence"`
	Direction    string  `json:"direction"`    // improving, degrading, stable
	Significance string  `json:"significance"` // high, medium, low
}

// PredictionPoint represents a predicted performance point
type PredictionPoint struct {
	Timestamp      time.Time `json:"timestamp"`
	PredictedValue float64   `json:"predicted_value"`
	Confidence     float64   `json:"confidence"`
	LowerBound     float64   `json:"lower_bound"`
	UpperBound     float64   `json:"upper_bound"`
}

// AnomalyDetection represents detected performance anomalies
type AnomalyDetection struct {
	Timestamp     time.Time `json:"timestamp"`
	Value         float64   `json:"value"`
	ExpectedValue float64   `json:"expected_value"`
	Deviation     float64   `json:"deviation"`
	Severity      string    `json:"severity"`
	Type          string    `json:"type"` // spike, drop, outlier
	Cause         string    `json:"cause,omitempty"`
}

// PerformanceProfiler provides detailed performance profiling
type PerformanceProfiler struct {
	enabled      bool
	samplingRate float64
	profiles     map[string]*ProfileData
	mu           sync.RWMutex
}

// ProfileData contains detailed profiling information
type ProfileData struct {
	FunctionName    string            `json:"function_name"`
	CallCount       int64             `json:"call_count"`
	TotalDuration   time.Duration     `json:"total_duration"`
	AverageDuration time.Duration     `json:"average_duration"`
	MinDuration     time.Duration     `json:"min_duration"`
	MaxDuration     time.Duration     `json:"max_duration"`
	MemoryAllocated int64             `json:"memory_allocated"`
	MemoryFreed     int64             `json:"memory_freed"`
	CPUUsage        float64           `json:"cpu_usage"`
	Hotspots        []string          `json:"hotspots"`
	CallStack       []string          `json:"call_stack"`
	Tags            map[string]string `json:"tags"`
}

// LoadGenerator provides configurable load generation for performance testing
type LoadGenerator struct {
	config     LoadGeneratorConfig
	httpClient *http.Client
	workers    []*LoadWorker
	results    chan WorkerResult
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// LoadGeneratorConfig defines load generation configuration
type LoadGeneratorConfig struct {
	TargetURL       string            `json:"target_url"`
	Method          string            `json:"method"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	Concurrency     int               `json:"concurrency"`
	RPS             float64           `json:"rps"`
	Duration        time.Duration     `json:"duration"`
	WarmupDuration  time.Duration     `json:"warmup_duration"`
	TimeoutDuration time.Duration     `json:"timeout_duration"`
	KeepAlive       bool              `json:"keep_alive"`
	TLSEnabled      bool              `json:"tls_enabled"`
	FollowRedirects bool              `json:"follow_redirects"`
}

// LoadWorker represents a single load generation worker
type LoadWorker struct {
	id       int
	client   *http.Client
	config   LoadGeneratorConfig
	results  chan<- WorkerResult
	stopChan <-chan struct{}
}

// WorkerResult contains results from a load worker
type WorkerResult struct {
	WorkerID      int           `json:"worker_id"`
	RequestID     string        `json:"request_id"`
	Timestamp     time.Time     `json:"timestamp"`
	Duration      time.Duration `json:"duration"`
	StatusCode    int           `json:"status_code"`
	ResponseSize  int64         `json:"response_size"`
	Error         string        `json:"error,omitempty"`
	DNSTime       time.Duration `json:"dns_time"`
	ConnectTime   time.Duration `json:"connect_time"`
	TLSTime       time.Duration `json:"tls_time"`
	FirstByteTime time.Duration `json:"first_byte_time"`
	DownloadTime  time.Duration `json:"download_time"`
}

// ResourceMonitor monitors system resource utilization during performance tests
type ResourceMonitor struct {
	enabled  bool
	interval time.Duration
	metrics  []ResourceSnapshot
	stopChan chan struct{}
	mu       sync.RWMutex
}

// ResourceSnapshot represents a snapshot of resource utilization
type ResourceSnapshot struct {
	Timestamp         time.Time       `json:"timestamp"`
	CPUPercent        float64         `json:"cpu_percent"`
	MemoryUsedMB      float64         `json:"memory_used_mb"`
	MemoryTotalMB     float64         `json:"memory_total_mb"`
	MemoryPercent     float64         `json:"memory_percent"`
	DiskReadKBps      float64         `json:"disk_read_kbps"`
	DiskWriteKBps     float64         `json:"disk_write_kbps"`
	NetworkInKBps     float64         `json:"network_in_kbps"`
	NetworkOutKBps    float64         `json:"network_out_kbps"`
	OpenFiles         int             `json:"open_files"`
	ActiveConnections int             `json:"active_connections"`
	Goroutines        int             `json:"goroutines"`
	GCPauses          []time.Duration `json:"gc_pauses"`
	HeapSizeMB        float64         `json:"heap_size_mb"`
	AllocRateMBps     float64         `json:"alloc_rate_mbps"`
}

// NewPerformanceAnalyzer creates a new performance analyzer
func NewPerformanceAnalyzer(config BenchmarkConfig, reportDir string) *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		config:    config,
		results:   make([]BenchmarkResult, 0),
		reportDir: reportDir,
		sessionID: uuid.New().String(),
	}
}

// AddResult adds a benchmark result to the analyzer
func (pa *PerformanceAnalyzer) AddResult(result BenchmarkResult) {
	pa.mu.Lock()
	defer pa.mu.Unlock()
	pa.results = append(pa.results, result)
}

// LoadBaseline loads baseline metrics from file
func (pa *PerformanceAnalyzer) LoadBaseline(baselineFile string) error {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	if baselineFile == "" {
		return nil // No baseline configured
	}

	data, err := ioutil.ReadFile(baselineFile)
	if err != nil {
		return fmt.Errorf("failed to read baseline file: %w", err)
	}

	var baseline BaselineMetrics
	if err := json.Unmarshal(data, &baseline); err != nil {
		return fmt.Errorf("failed to parse baseline metrics: %w", err)
	}

	pa.baseline = &baseline
	return nil
}

// SaveBaseline saves current results as baseline metrics
func (pa *PerformanceAnalyzer) SaveBaseline(baselineFile string, version string) error {
	pa.mu.RLock()
	defer pa.mu.RUnlock()

	baseline := BaselineMetrics{
		Version:     version,
		Timestamp:   time.Now().UTC(),
		Services:    make(map[string]ServiceBaseline),
		Environment: "production",
		TestSuite:   "performance_benchmark",
	}

	// Convert results to baseline format
	serviceMap := make(map[string]map[string][]BenchmarkResult)

	for _, result := range pa.results {
		if serviceMap[result.Service] == nil {
			serviceMap[result.Service] = make(map[string][]BenchmarkResult)
		}

		testType := "endpoint"
		if result.Endpoint == "MessageQueue Producer" || result.Endpoint == "MessageQueue Consumer" {
			testType = "messagequeue"
		} else if len(result.Tags) > 0 {
			for _, tag := range result.Tags {
				if tag == "database" {
					testType = "database"
					break
				} else if tag == "cache" {
					testType = "cache"
					break
				}
			}
		}

		serviceMap[result.Service][testType] = append(serviceMap[result.Service][testType], result)
	}

	// Build service baselines
	for serviceName, typeMap := range serviceMap {
		serviceBaseline := ServiceBaseline{
			ServiceName: serviceName,
			Endpoints:   make(map[string]EndpointBaseline),
		}

		if endpoints, exists := typeMap["endpoint"]; exists {
			for _, result := range endpoints {
				endpointBaseline := EndpointBaseline{
					Path:          result.Endpoint,
					ExpectedRPS:   result.RPS * 0.95, // 95% of measured performance
					ExpectedP95:   result.ResponseTimes.P95,
					ExpectedP99:   result.ResponseTimes.P99,
					MaxCPUPercent: result.Resources.CPUUsage.Average * 1.1, // 10% margin
					MaxMemoryMB:   result.Resources.MemoryUsage.PeakGB * 1024 * 1.1,
					MaxErrorRate:  float64(result.TotalErrors) / float64(result.TotalRequests) * 100.0 * 2.0, // Double error rate as max
					Tags:          result.Tags,
				}
				serviceBaseline.Endpoints[result.TestName] = endpointBaseline
			}
		}

		baseline.Services[serviceName] = serviceBaseline
	}

	// Save to file
	data, err := json.MarshalIndent(baseline, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal baseline: %w", err)
	}

	return ioutil.WriteFile(baselineFile, data, 0644)
}

// DetectRegressions compares current results against baseline to detect regressions
func (pa *PerformanceAnalyzer) DetectRegressions(threshold float64) []PerformanceRegression {
	pa.mu.RLock()
	defer pa.mu.RUnlock()

	if pa.baseline == nil {
		return nil
	}

	var regressions []PerformanceRegression

	for _, result := range pa.results {
		serviceBaseline, serviceExists := pa.baseline.Services[result.Service]
		if !serviceExists {
			continue
		}

		endpointBaseline, endpointExists := serviceBaseline.Endpoints[result.TestName]
		if !endpointExists {
			continue
		}

		// Check RPS regression
		if result.RPS < endpointBaseline.ExpectedRPS {
			regressionPct := (endpointBaseline.ExpectedRPS - result.RPS) / endpointBaseline.ExpectedRPS * 100.0
			if regressionPct > threshold {
				regressions = append(regressions, PerformanceRegression{
					TestName:      result.TestName,
					Service:       result.Service,
					Metric:        "rps",
					Current:       result.RPS,
					Baseline:      endpointBaseline.ExpectedRPS,
					RegressionPct: regressionPct,
					Threshold:     threshold,
					Severity:      pa.determineSeverity(regressionPct),
					Impact:        "throughput",
					Recommendations: []string{
						"Profile application performance",
						"Check for resource bottlenecks",
						"Review recent code changes",
					},
					Timestamp: time.Now().UTC(),
				})
			}
		}

		// Check P95 latency regression
		if result.ResponseTimes.P95 > endpointBaseline.ExpectedP95 {
			regressionPct := float64(result.ResponseTimes.P95-endpointBaseline.ExpectedP95) / float64(endpointBaseline.ExpectedP95) * 100.0
			if regressionPct > threshold {
				regressions = append(regressions, PerformanceRegression{
					TestName:      result.TestName,
					Service:       result.Service,
					Metric:        "p95_latency",
					Current:       float64(result.ResponseTimes.P95.Milliseconds()),
					Baseline:      float64(endpointBaseline.ExpectedP95.Milliseconds()),
					RegressionPct: regressionPct,
					Threshold:     threshold,
					Severity:      pa.determineSeverity(regressionPct),
					Impact:        "user_experience",
					Recommendations: []string{
						"Optimize slow database queries",
						"Add response caching",
						"Check network latency",
					},
					Timestamp: time.Now().UTC(),
				})
			}
		}

		// Check resource utilization regression
		if result.Resources.CPUUsage.Average > endpointBaseline.MaxCPUPercent {
			regressionPct := (result.Resources.CPUUsage.Average - endpointBaseline.MaxCPUPercent) / endpointBaseline.MaxCPUPercent * 100.0
			if regressionPct > threshold {
				regressions = append(regressions, PerformanceRegression{
					TestName:      result.TestName,
					Service:       result.Service,
					Metric:        "cpu_usage",
					Current:       result.Resources.CPUUsage.Average,
					Baseline:      endpointBaseline.MaxCPUPercent,
					RegressionPct: regressionPct,
					Threshold:     threshold,
					Severity:      pa.determineSeverity(regressionPct),
					Impact:        "resource_efficiency",
					Recommendations: []string{
						"Profile CPU hotspots",
						"Optimize algorithms",
						"Consider horizontal scaling",
					},
					Timestamp: time.Now().UTC(),
				})
			}
		}
	}

	return regressions
}

// GeneratePerformanceReport generates comprehensive performance analysis report
func (pa *PerformanceAnalyzer) GeneratePerformanceReport() (*PerformanceReport, error) {
	pa.mu.RLock()
	defer pa.mu.RUnlock()

	report := &PerformanceReport{
		SessionID:       pa.sessionID,
		Timestamp:       time.Now().UTC(),
		TotalTests:      len(pa.results),
		Summary:         pa.generateSummaryStats(),
		ServiceStats:    pa.generateServiceStats(),
		Regressions:     pa.DetectRegressions(10.0), // 10% regression threshold
		Trends:          pa.generateTrendAnalysis(),
		Recommendations: pa.generateRecommendations(),
		Metadata: map[string]interface{}{
			"go_version": runtime.Version(),
			"num_cpu":    runtime.NumCPU(),
			"config":     pa.config,
		},
	}

	return report, nil
}

// PerformanceReport represents a comprehensive performance analysis report
type PerformanceReport struct {
	SessionID       string                  `json:"session_id"`
	Timestamp       time.Time               `json:"timestamp"`
	TotalTests      int                     `json:"total_tests"`
	Summary         SummaryStats            `json:"summary"`
	ServiceStats    map[string]ServiceStats `json:"service_stats"`
	Regressions     []PerformanceRegression `json:"regressions"`
	Trends          []PerformanceTrend      `json:"trends"`
	Recommendations []string                `json:"recommendations"`
	Metadata        map[string]interface{}  `json:"metadata"`
}

// SummaryStats contains overall performance summary statistics
type SummaryStats struct {
	TotalRequests   int64         `json:"total_requests"`
	TotalErrors     int64         `json:"total_errors"`
	ErrorRate       float64       `json:"error_rate"`
	AverageRPS      float64       `json:"average_rps"`
	MedianLatency   time.Duration `json:"median_latency"`
	P95Latency      time.Duration `json:"p95_latency"`
	P99Latency      time.Duration `json:"p99_latency"`
	AverageCPU      float64       `json:"average_cpu"`
	AverageMemoryMB float64       `json:"average_memory_mb"`
	PassedTests     int           `json:"passed_tests"`
	WarningTests    int           `json:"warning_tests"`
	FailedTests     int           `json:"failed_tests"`
	SuccessRate     float64       `json:"success_rate"`
}

// ServiceStats contains service-specific performance statistics
type ServiceStats struct {
	ServiceName     string                   `json:"service_name"`
	EndpointStats   map[string]EndpointStats `json:"endpoint_stats"`
	TotalRequests   int64                    `json:"total_requests"`
	TotalErrors     int64                    `json:"total_errors"`
	ErrorRate       float64                  `json:"error_rate"`
	AverageRPS      float64                  `json:"average_rps"`
	MedianLatency   time.Duration            `json:"median_latency"`
	P95Latency      time.Duration            `json:"p95_latency"`
	AverageCPU      float64                  `json:"average_cpu"`
	AverageMemoryMB float64                  `json:"average_memory_mb"`
	ViolationCount  int                      `json:"violation_count"`
}

// EndpointStats contains endpoint-specific performance statistics
type EndpointStats struct {
	Endpoint        string        `json:"endpoint"`
	Method          string        `json:"method"`
	TotalRequests   int64         `json:"total_requests"`
	TotalErrors     int64         `json:"total_errors"`
	ErrorRate       float64       `json:"error_rate"`
	RPS             float64       `json:"rps"`
	MedianLatency   time.Duration `json:"median_latency"`
	P95Latency      time.Duration `json:"p95_latency"`
	P99Latency      time.Duration `json:"p99_latency"`
	MinLatency      time.Duration `json:"min_latency"`
	MaxLatency      time.Duration `json:"max_latency"`
	AverageCPU      float64       `json:"average_cpu"`
	AverageMemoryMB float64       `json:"average_memory_mb"`
	Tags            []string      `json:"tags"`
}

// generateSummaryStats generates overall summary statistics
func (pa *PerformanceAnalyzer) generateSummaryStats() SummaryStats {
	var totalRequests, totalErrors int64
	var totalRPS, totalCPU, totalMemory float64
	var latencies []time.Duration
	var passedTests, warningTests, failedTests int

	for _, result := range pa.results {
		totalRequests += result.TotalRequests
		totalErrors += result.TotalErrors
		totalRPS += result.RPS
		totalCPU += result.Resources.CPUUsage.Average
		totalMemory += result.Resources.MemoryUsage.PeakGB * 1024
		latencies = append(latencies, result.ResponseTimes.Median, result.ResponseTimes.P95, result.ResponseTimes.P99)

		switch result.Status {
		case "PASS":
			passedTests++
		case "WARNING":
			warningTests++
		case "FAIL":
			failedTests++
		}
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	var medianLatency, p95Latency, p99Latency time.Duration
	if len(latencies) > 0 {
		medianLatency = latencies[len(latencies)/2]
		p95Latency = latencies[int(float64(len(latencies))*0.95)]
		p99Latency = latencies[int(float64(len(latencies))*0.99)]
	}

	errorRate := 0.0
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests) * 100.0
	}

	successRate := 0.0
	if len(pa.results) > 0 {
		successRate = float64(passedTests) / float64(len(pa.results)) * 100.0
	}

	return SummaryStats{
		TotalRequests:   totalRequests,
		TotalErrors:     totalErrors,
		ErrorRate:       errorRate,
		AverageRPS:      totalRPS / float64(len(pa.results)),
		MedianLatency:   medianLatency,
		P95Latency:      p95Latency,
		P99Latency:      p99Latency,
		AverageCPU:      totalCPU / float64(len(pa.results)),
		AverageMemoryMB: totalMemory / float64(len(pa.results)),
		PassedTests:     passedTests,
		WarningTests:    warningTests,
		FailedTests:     failedTests,
		SuccessRate:     successRate,
	}
}

// generateServiceStats generates service-specific statistics
func (pa *PerformanceAnalyzer) generateServiceStats() map[string]ServiceStats {
	serviceMap := make(map[string][]BenchmarkResult)

	// Group results by service
	for _, result := range pa.results {
		serviceMap[result.Service] = append(serviceMap[result.Service], result)
	}

	serviceStats := make(map[string]ServiceStats)

	for serviceName, results := range serviceMap {
		var totalRequests, totalErrors int64
		var totalRPS, totalCPU, totalMemory float64
		var latencies []time.Duration
		var violationCount int

		endpointMap := make(map[string][]BenchmarkResult)

		for _, result := range results {
			totalRequests += result.TotalRequests
			totalErrors += result.TotalErrors
			totalRPS += result.RPS
			totalCPU += result.Resources.CPUUsage.Average
			totalMemory += result.Resources.MemoryUsage.PeakGB * 1024
			latencies = append(latencies, result.ResponseTimes.Median, result.ResponseTimes.P95)
			violationCount += len(result.Violations)

			endpointMap[result.Endpoint] = append(endpointMap[result.Endpoint], result)
		}

		sort.Slice(latencies, func(i, j int) bool {
			return latencies[i] < latencies[j]
		})

		var medianLatency, p95Latency time.Duration
		if len(latencies) > 0 {
			medianLatency = latencies[len(latencies)/2]
			p95Latency = latencies[int(float64(len(latencies))*0.95)]
		}

		errorRate := 0.0
		if totalRequests > 0 {
			errorRate = float64(totalErrors) / float64(totalRequests) * 100.0
		}

		// Generate endpoint stats
		endpointStats := make(map[string]EndpointStats)
		for endpoint, endpointResults := range endpointMap {
			endpointStats[endpoint] = pa.generateEndpointStats(endpoint, endpointResults)
		}

		serviceStats[serviceName] = ServiceStats{
			ServiceName:     serviceName,
			EndpointStats:   endpointStats,
			TotalRequests:   totalRequests,
			TotalErrors:     totalErrors,
			ErrorRate:       errorRate,
			AverageRPS:      totalRPS / float64(len(results)),
			MedianLatency:   medianLatency,
			P95Latency:      p95Latency,
			AverageCPU:      totalCPU / float64(len(results)),
			AverageMemoryMB: totalMemory / float64(len(results)),
			ViolationCount:  violationCount,
		}
	}

	return serviceStats
}

// generateEndpointStats generates endpoint-specific statistics
func (pa *PerformanceAnalyzer) generateEndpointStats(endpoint string, results []BenchmarkResult) EndpointStats {
	if len(results) == 0 {
		return EndpointStats{}
	}

	var totalRequests, totalErrors int64
	var totalRPS, totalCPU, totalMemory float64
	var latencies []time.Duration
	var minLatency, maxLatency time.Duration = time.Hour, 0

	firstResult := results[0]

	for _, result := range results {
		totalRequests += result.TotalRequests
		totalErrors += result.TotalErrors
		totalRPS += result.RPS
		totalCPU += result.Resources.CPUUsage.Average
		totalMemory += result.Resources.MemoryUsage.PeakGB * 1024

		latencies = append(latencies, result.ResponseTimes.Median, result.ResponseTimes.P95, result.ResponseTimes.P99)

		if result.ResponseTimes.Min < minLatency {
			minLatency = result.ResponseTimes.Min
		}
		if result.ResponseTimes.Max > maxLatency {
			maxLatency = result.ResponseTimes.Max
		}
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	var medianLatency, p95Latency, p99Latency time.Duration
	if len(latencies) > 0 {
		medianLatency = latencies[len(latencies)/2]
		p95Latency = latencies[int(float64(len(latencies))*0.95)]
		p99Latency = latencies[int(float64(len(latencies))*0.99)]
	}

	errorRate := 0.0
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests) * 100.0
	}

	return EndpointStats{
		Endpoint:        endpoint,
		TotalRequests:   totalRequests,
		TotalErrors:     totalErrors,
		ErrorRate:       errorRate,
		RPS:             totalRPS / float64(len(results)),
		MedianLatency:   medianLatency,
		P95Latency:      p95Latency,
		P99Latency:      p99Latency,
		MinLatency:      minLatency,
		MaxLatency:      maxLatency,
		AverageCPU:      totalCPU / float64(len(results)),
		AverageMemoryMB: totalMemory / float64(len(results)),
		Tags:            firstResult.Tags,
	}
}

// generateTrendAnalysis generates performance trend analysis
func (pa *PerformanceAnalyzer) generateTrendAnalysis() []PerformanceTrend {
	// This would typically analyze historical data
	// For now, return placeholder trends
	return []PerformanceTrend{
		{
			TestName:  "performance_trends",
			Service:   "overall",
			Metric:    "rps",
			Timeframe: "last_30_days",
			TrendLine: TrendAnalysis{
				Direction:    "stable",
				Significance: "medium",
				Confidence:   0.85,
			},
		},
	}
}

// generateRecommendations generates performance optimization recommendations
func (pa *PerformanceAnalyzer) generateRecommendations() []string {
	recommendations := []string{}

	// Analyze results and generate recommendations
	highCPUServices := []string{}
	highMemoryServices := []string{}
	highLatencyServices := []string{}
	lowThroughputServices := []string{}

	for _, result := range pa.results {
		if result.Resources.CPUUsage.Average > 70.0 {
			highCPUServices = append(highCPUServices, result.Service)
		}
		if result.Resources.MemoryUsage.PeakGB * 1024 > 400.0 {
			highMemoryServices = append(highMemoryServices, result.Service)
		}
		if result.ResponseTimes.P95 > 200*time.Millisecond {
			highLatencyServices = append(highLatencyServices, result.Service)
		}
		if result.RPS < 500.0 {
			lowThroughputServices = append(lowThroughputServices, result.Service)
		}
	}

	if len(highCPUServices) > 0 {
		recommendations = append(recommendations, "Optimize CPU-intensive operations in services: "+joinUnique(highCPUServices))
	}
	if len(highMemoryServices) > 0 {
		recommendations = append(recommendations, "Investigate memory usage in services: "+joinUnique(highMemoryServices))
	}
	if len(highLatencyServices) > 0 {
		recommendations = append(recommendations, "Reduce response times for services: "+joinUnique(highLatencyServices))
	}
	if len(lowThroughputServices) > 0 {
		recommendations = append(recommendations, "Improve throughput for services: "+joinUnique(lowThroughputServices))
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Performance is within acceptable ranges")
	}

	return recommendations
}

// determineSeverity determines regression severity based on percentage
func (pa *PerformanceAnalyzer) determineSeverity(regressionPct float64) string {
	if regressionPct > 50.0 {
		return "CRITICAL"
	} else if regressionPct > 25.0 {
		return "HIGH"
	} else if regressionPct > 10.0 {
		return "MEDIUM"
	}
	return "LOW"
}

// SaveReportToFile saves performance report to file
func (pa *PerformanceAnalyzer) SaveReportToFile(report *PerformanceReport, format string) error {
	if err := os.MkdirAll(pa.reportDir, 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	filename := fmt.Sprintf("performance_report_%s_%s.%s",
		pa.sessionID,
		time.Now().Format("20060102_150405"),
		format)
	filepath := filepath.Join(pa.reportDir, filename)

	switch format {
	case "json":
		return pa.saveJSONReport(report, filepath)
	case "html":
		return pa.saveHTMLReport(report, filepath)
	case "csv":
		return pa.saveCSVReport(report, filepath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// saveJSONReport saves report in JSON format
func (pa *PerformanceAnalyzer) saveJSONReport(report *PerformanceReport, filepath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON report: %w", err)
	}
	return ioutil.WriteFile(filepath, data, 0644)
}

// saveHTMLReport saves report in HTML format
func (pa *PerformanceAnalyzer) saveHTMLReport(report *PerformanceReport, filepath string) error {
	html := pa.generateHTMLReport(report)
	return ioutil.WriteFile(filepath, []byte(html), 0644)
}

// saveCSVReport saves report in CSV format
func (pa *PerformanceAnalyzer) saveCSVReport(report *PerformanceReport, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV headers
	headers := []string{
		"TestName", "Service", "Endpoint", "RPS", "P95Latency", "P99Latency",
		"ErrorRate", "CPUPercent", "MemoryMB", "Status", "ViolationCount",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Write data rows
	for _, result := range pa.results {
		row := []string{
			result.TestName,
			result.Service,
			result.Endpoint,
			fmt.Sprintf("%.2f", result.RPS),
			fmt.Sprintf("%d", result.ResponseTimes.P95.Milliseconds()),
			fmt.Sprintf("%d", result.ResponseTimes.P99.Milliseconds()),
			fmt.Sprintf("%.2f", float64(result.TotalErrors)/float64(result.TotalRequests)*100.0),
			fmt.Sprintf("%.2f", result.Resources.CPUUsage.Average),
			fmt.Sprintf("%.2f", result.Resources.MemoryUsage.PeakGB * 1024),
			result.Status,
			fmt.Sprintf("%d", len(result.Violations)),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// generateHTMLReport generates HTML performance report
func (pa *PerformanceAnalyzer) generateHTMLReport(report *PerformanceReport) string {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Performance Test Report - Session %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 20px; border-radius: 5px; }
        .summary { margin: 20px 0; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; }
        .stat-card { background: #f9f9f9; padding: 15px; border-radius: 5px; border: 1px solid #ddd; }
        .stat-value { font-size: 24px; font-weight: bold; color: #2196F3; }
        .stat-label { color: #666; font-size: 14px; }
        .service-section { margin: 30px 0; }
        .service-title { background: #e3f2fd; padding: 10px; border-radius: 5px; }
        table { width: 100%%; border-collapse: collapse; margin: 10px 0; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f5f5f5; }
        .status-pass { color: green; font-weight: bold; }
        .status-warning { color: orange; font-weight: bold; }
        .status-fail { color: red; font-weight: bold; }
        .recommendations { background: #fff3cd; padding: 15px; border-radius: 5px; border: 1px solid #ffeaa7; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Performance Test Report</h1>
        <p>Session ID: %s</p>
        <p>Generated: %s</p>
        <p>Total Tests: %d</p>
    </div>

    <div class="summary">
        <h2>Summary Statistics</h2>
        <div class="stats">
            <div class="stat-card">
                <div class="stat-value">%.1f%%</div>
                <div class="stat-label">Success Rate</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">%.0f</div>
                <div class="stat-label">Average RPS</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">%dms</div>
                <div class="stat-label">P95 Latency</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">%.1f%%</div>
                <div class="stat-label">Error Rate</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">%.1f%%</div>
                <div class="stat-label">Average CPU</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">%.0fMB</div>
                <div class="stat-label">Average Memory</div>
            </div>
        </div>
    </div>
`,
		report.SessionID,
		report.SessionID,
		report.Timestamp.Format("2006-01-02 15:04:05 UTC"),
		report.TotalTests,
		report.Summary.SuccessRate,
		report.Summary.AverageRPS,
		report.Summary.P95Latency.Milliseconds(),
		report.Summary.ErrorRate,
		report.Summary.AverageCPU,
		report.Summary.AverageMemoryMB,
	)

	// Add service details
	for serviceName, serviceStats := range report.ServiceStats {
		html += fmt.Sprintf(`
    <div class="service-section">
        <div class="service-title">
            <h3>Service: %s</h3>
        </div>
        <table>
            <tr>
                <th>Metric</th>
                <th>Value</th>
            </tr>
            <tr><td>Total Requests</td><td>%d</td></tr>
            <tr><td>Error Rate</td><td>%.2f%%</td></tr>
            <tr><td>Average RPS</td><td>%.1f</td></tr>
            <tr><td>P95 Latency</td><td>%dms</td></tr>
            <tr><td>Average CPU</td><td>%.1f%%</td></tr>
            <tr><td>Average Memory</td><td>%.1fMB</td></tr>
            <tr><td>Violations</td><td>%d</td></tr>
        </table>
    </div>`,
			serviceName,
			serviceStats.TotalRequests,
			serviceStats.ErrorRate,
			serviceStats.AverageRPS,
			serviceStats.P95Latency.Milliseconds(),
			serviceStats.AverageCPU,
			serviceStats.AverageMemoryMB,
			serviceStats.ViolationCount,
		)
	}

	// Add recommendations
	if len(report.Recommendations) > 0 {
		html += `
    <div class="recommendations">
        <h3>Recommendations</h3>
        <ul>`
		for _, rec := range report.Recommendations {
			html += fmt.Sprintf("<li>%s</li>", rec)
		}
		html += `
        </ul>
    </div>`
	}

	html += `
</body>
</html>`

	return html
}

// joinUnique joins unique strings with commas
func joinUnique(strings []string) string {
	unique := make(map[string]bool)
	var result []string

	for _, s := range strings {
		if !unique[s] {
			unique[s] = true
			result = append(result, s)
		}
	}

	output := ""
	for i, s := range result {
		if i > 0 {
			output += ", "
		}
		output += s
	}

	return output
}

// NewLoadGenerator creates a new load generator
func NewLoadGenerator(config LoadGeneratorConfig) *LoadGenerator {
	return &LoadGenerator{
		config:     config,
		httpClient: &http.Client{Timeout: config.TimeoutDuration},
		results:    make(chan WorkerResult, config.Concurrency*10),
		stopChan:   make(chan struct{}),
	}
}

// NewResourceMonitor creates a new resource monitor
func NewResourceMonitor(interval time.Duration) *ResourceMonitor {
	return &ResourceMonitor{
		enabled:  true,
		interval: interval,
		metrics:  make([]ResourceSnapshot, 0),
		stopChan: make(chan struct{}),
	}
}

// Start starts the resource monitor
func (rm *ResourceMonitor) Start() {
	if !rm.enabled {
		return
	}

	go func() {
		ticker := time.NewTicker(rm.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				snapshot := rm.captureSnapshot()
				rm.mu.Lock()
				rm.metrics = append(rm.metrics, snapshot)
				rm.mu.Unlock()
			case <-rm.stopChan:
				return
			}
		}
	}()
}

// Stop stops the resource monitor
func (rm *ResourceMonitor) Stop() {
	if rm.enabled {
		close(rm.stopChan)
		rm.enabled = false
	}
}

// GetMetrics returns captured resource metrics
func (rm *ResourceMonitor) GetMetrics() []ResourceSnapshot {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	result := make([]ResourceSnapshot, len(rm.metrics))
	copy(result, rm.metrics)
	return result
}

// captureSnapshot captures current resource utilization snapshot
func (rm *ResourceMonitor) captureSnapshot() ResourceSnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return ResourceSnapshot{
		Timestamp:         time.Now().UTC(),
		CPUPercent:        0.0, // Would use actual CPU monitoring in real implementation
		MemoryUsedMB:      float64(m.Alloc) / 1024 / 1024,
		MemoryTotalMB:     float64(m.Sys) / 1024 / 1024,
		MemoryPercent:     float64(m.Alloc) / float64(m.Sys) * 100.0,
		Goroutines:        runtime.NumGoroutine(),
		HeapSizeMB:        float64(m.HeapSys) / 1024 / 1024,
		AllocRateMBps:     0.0, // Would calculate based on allocation rate
		OpenFiles:         0,   // Would use actual file descriptor count
		ActiveConnections: 0,   // Would use actual connection count
		DiskReadKBps:      0.0, // Would use actual disk I/O monitoring
		DiskWriteKBps:     0.0,
		NetworkInKBps:     0.0, // Would use actual network monitoring
		NetworkOutKBps:    0.0,
		GCPauses:          []time.Duration{time.Duration(m.PauseNs[(m.NumGC+255)%256])},
	}
}
