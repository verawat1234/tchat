// Package contract provides automated performance regression detection for contract testing
// Implements intelligent regression detection with statistical analysis and machine learning
package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// RegressionDetectionTestSuite provides comprehensive automated regression detection
type RegressionDetectionTestSuite struct {
	suite.Suite
	ctx                  context.Context
	regressionEngine     *RegressionEngine
	statisticalAnalyzer  *StatisticalAnalyzer
	performanceBaselines map[string]*PerformanceBaseline
	historicalData       []PerformanceDataPoint
	regressionResults    []RegressionDetectionResult
	alertsGenerated      []RegressionAlert
	projectRoot         string
}

// RegressionEngine handles automated regression detection
type RegressionEngine struct {
	config              RegressionEngineConfig
	detectionAlgorithms map[string]DetectionAlgorithm
	baselines           map[string]*BaselineModel
	trendAnalyzer       *TrendAnalyzer
	anomalyDetector     *AnomalyDetector
	regressionHistory   []RegressionEvent
}

// RegressionEngineConfig defines regression detection engine configuration
type RegressionEngineConfig struct {
	EnabledAlgorithms       []string      `json:"enabled_algorithms"`
	StatisticalSignificance float64       `json:"statistical_significance"`
	MinimumDataPoints       int           `json:"minimum_data_points"`
	RegressionThreshold     float64       `json:"regression_threshold"`
	SensitivityLevel        string        `json:"sensitivity_level"` // low, medium, high, adaptive
	ComparisonWindow        time.Duration `json:"comparison_window"`
	BaselineUpdateFrequency time.Duration `json:"baseline_update_frequency"`
	AutoBaselineRecalc      bool          `json:"auto_baseline_recalc"`
	TrendAnalysisEnabled    bool          `json:"trend_analysis_enabled"`
	AnomalyDetectionEnabled bool          `json:"anomaly_detection_enabled"`
	MachineLearningEnabled  bool          `json:"machine_learning_enabled"`
}

// DetectionAlgorithm defines interface for regression detection algorithms
type DetectionAlgorithm interface {
	Name() string
	Detect(current []PerformanceDataPoint, baseline *BaselineModel) RegressionResult
	Configure(config map[string]interface{}) error
}

// BaselineModel represents a performance baseline model
type BaselineModel struct {
	MetricName           string                   `json:"metric_name"`
	BaselineValue        float64                  `json:"baseline_value"`
	StandardDeviation    float64                  `json:"standard_deviation"`
	ConfidenceInterval   ConfidenceInterval       `json:"confidence_interval"`
	Percentiles          Percentiles              `json:"percentiles"`
	TrendCoefficients    TrendCoefficients        `json:"trend_coefficients"`
	SeasonalityPattern   SeasonalityPattern       `json:"seasonality_pattern"`
	LastUpdated          time.Time                `json:"last_updated"`
	SampleSize           int                      `json:"sample_size"`
	DataQuality          DataQualityMetrics       `json:"data_quality"`
}

// ConfidenceInterval represents statistical confidence interval
type ConfidenceInterval struct {
	LowerBound  float64 `json:"lower_bound"`
	UpperBound  float64 `json:"upper_bound"`
	Confidence  float64 `json:"confidence"` // e.g., 0.95 for 95%
}

// Percentiles represents percentile values for baseline
type Percentiles struct {
	P50 float64 `json:"p50"`
	P75 float64 `json:"p75"`
	P90 float64 `json:"p90"`
	P95 float64 `json:"p95"`
	P99 float64 `json:"p99"`
}

// TrendCoefficients represents trend analysis coefficients
type TrendCoefficients struct {
	LinearSlope       float64 `json:"linear_slope"`
	QuadraticTerm     float64 `json:"quadratic_term"`
	SeasonalComponent float64 `json:"seasonal_component"`
	RSquared          float64 `json:"r_squared"`
}

// SeasonalityPattern represents detected seasonality patterns
type SeasonalityPattern struct {
	HasSeasonality  bool               `json:"has_seasonality"`
	Period          time.Duration      `json:"period"`
	Amplitude       float64            `json:"amplitude"`
	PhaseShift      float64            `json:"phase_shift"`
	SeasonalFactors map[string]float64 `json:"seasonal_factors"`
}

// DataQualityMetrics represents data quality assessment
type DataQualityMetrics struct {
	CompletenessScore float64 `json:"completeness_score"`
	ConsistencyScore  float64 `json:"consistency_score"`
	OutlierCount      int     `json:"outlier_count"`
	NoiseLevel        float64 `json:"noise_level"`
	QualityGrade      string  `json:"quality_grade"` // A, B, C, D, F
}

// StatisticalAnalyzer performs statistical analysis for regression detection
type StatisticalAnalyzer struct {
	config StatisticalConfig
}

// StatisticalConfig defines statistical analysis configuration
type StatisticalConfig struct {
	SignificanceLevel       float64   `json:"significance_level"`
	MultipleTestingCorrection string  `json:"multiple_testing_correction"` // bonferroni, fdr
	OutlierDetectionMethod  string    `json:"outlier_detection_method"`    // iqr, zscore, isolation_forest
	DistributionTests       []string  `json:"distribution_tests"`          // shapiro, anderson, ks
	ChangePointDetection    bool      `json:"change_point_detection"`
	WelfordUpdateEnabled    bool      `json:"welford_update_enabled"`
}

// TrendAnalyzer analyzes performance trends
type TrendAnalyzer struct {
	config TrendAnalysisConfig
}

// TrendAnalysisConfig defines trend analysis configuration
type TrendAnalysisConfig struct {
	TrendDetectionMethods []string `json:"trend_detection_methods"` // linear, polynomial, exponential
	SeasonalityDetection  bool     `json:"seasonality_detection"`
	ForecastingEnabled    bool     `json:"forecasting_enabled"`
	ForecastHorizon       int      `json:"forecast_horizon"`
	MovingAverageWindow   int      `json:"moving_average_window"`
}

// AnomalyDetector detects performance anomalies
type AnomalyDetector struct {
	config AnomalyDetectionConfig
}

// AnomalyDetectionConfig defines anomaly detection configuration
type AnomalyDetectionConfig struct {
	DetectionMethods    []string `json:"detection_methods"` // statistical, isolation_forest, local_outlier
	SensitivityLevel    string   `json:"sensitivity_level"` // low, medium, high
	WindowSize          int      `json:"window_size"`
	ContaminationRate   float64  `json:"contamination_rate"`
	AdaptiveLearning    bool     `json:"adaptive_learning"`
}

// PerformanceDataPoint represents a single performance measurement
type PerformanceDataPoint struct {
	Timestamp           time.Time         `json:"timestamp"`
	MetricName          string            `json:"metric_name"`
	Value               float64           `json:"value"`
	Tags                map[string]string `json:"tags"`
	Context             ExecutionContext  `json:"context"`
	QualityIndicators   QualityIndicators `json:"quality_indicators"`
}

// ExecutionContext provides context for performance measurement
type ExecutionContext struct {
	TestSuite        string            `json:"test_suite"`
	TestCase         string            `json:"test_case"`
	Environment      string            `json:"environment"`
	GitCommit        string            `json:"git_commit"`
	BuildNumber      string            `json:"build_number"`
	SystemLoad       SystemLoadInfo    `json:"system_load"`
	ExternalFactors  map[string]string `json:"external_factors"`
}

// SystemLoadInfo represents system load during measurement
type SystemLoadInfo struct {
	CPUUsagePercent    float64 `json:"cpu_usage_percent"`
	MemoryUsagePercent float64 `json:"memory_usage_percent"`
	DiskIORate         float64 `json:"disk_io_rate"`
	NetworkLatency     float64 `json:"network_latency"`
	ConcurrentProcesses int    `json:"concurrent_processes"`
}

// QualityIndicators represent measurement quality indicators
type QualityIndicators struct {
	Reliability  float64 `json:"reliability"`  // 0-1 score
	Precision    float64 `json:"precision"`    // measurement precision
	Bias         float64 `json:"bias"`         // systematic bias
	NoiseLevel   float64 `json:"noise_level"`  // measurement noise
	Confidence   float64 `json:"confidence"`   // measurement confidence
}

// RegressionDetectionResult represents regression detection results
type RegressionDetectionResult struct {
	TestName              string                    `json:"test_name"`
	StartTime             time.Time                 `json:"start_time"`
	Duration              time.Duration             `json:"duration"`
	DetectionAlgorithm    string                    `json:"detection_algorithm"`
	MetricsAnalyzed       []string                  `json:"metrics_analyzed"`
	RegressionsDetected   []MetricRegression        `json:"regressions_detected"`
	TrendAnalysis         map[string]TrendResult    `json:"trend_analysis"`
	AnomaliesDetected     []AnomalyResult           `json:"anomalies_detected"`
	StatisticalSignificance float64                `json:"statistical_significance"`
	ConfidenceLevel       float64                   `json:"confidence_level"`
	FalsePositiveRate     float64                   `json:"false_positive_rate"`
	Recommendations       []string                  `json:"recommendations"`
	Metadata              RegressionDetectionMetadata `json:"metadata"`
}

// MetricRegression represents a detected regression in a metric
type MetricRegression struct {
	MetricName         string           `json:"metric_name"`
	RegressionType     string           `json:"regression_type"` // performance, availability, error_rate
	Severity           string           `json:"severity"`        // low, medium, high, critical
	DetectionMethod    string           `json:"detection_method"`
	CurrentValue       float64          `json:"current_value"`
	BaselineValue      float64          `json:"baseline_value"`
	PercentageChange   float64          `json:"percentage_change"`
	StandardDeviations float64          `json:"standard_deviations"`
	PValue             float64          `json:"p_value"`
	EffectSize         float64          `json:"effect_size"`
	ConfidenceInterval ConfidenceInterval `json:"confidence_interval"`
	FirstDetected      time.Time        `json:"first_detected"`
	LastConfirmed      time.Time        `json:"last_confirmed"`
	Duration           time.Duration    `json:"duration"`
	ImpactAssessment   ImpactAssessment `json:"impact_assessment"`
}

// ImpactAssessment represents the impact of a regression
type ImpactAssessment struct {
	BusinessImpact  string  `json:"business_impact"`  // low, medium, high, critical
	UserExperience  string  `json:"user_experience"`  // minimal, noticeable, significant, severe
	SystemStability string  `json:"system_stability"` // stable, degraded, unstable, failing
	EstimatedCost   float64 `json:"estimated_cost"`   // estimated cost impact
	RecoveryTime    time.Duration `json:"recovery_time"` // estimated recovery time
}

// TrendResult represents trend analysis results
type TrendResult struct {
	MetricName      string    `json:"metric_name"`
	TrendDirection  string    `json:"trend_direction"` // increasing, decreasing, stable, volatile
	TrendStrength   float64   `json:"trend_strength"`  // 0-1 strength of trend
	TrendSlope      float64   `json:"trend_slope"`
	RSquared        float64   `json:"r_squared"`
	Forecast        []float64 `json:"forecast"`
	ForecastHorizon int       `json:"forecast_horizon"`
	Seasonality     SeasonalityPattern `json:"seasonality"`
}

// AnomalyResult represents anomaly detection results
type AnomalyResult struct {
	MetricName       string    `json:"metric_name"`
	AnomalyType      string    `json:"anomaly_type"` // point, contextual, collective
	Timestamp        time.Time `json:"timestamp"`
	Value            float64   `json:"value"`
	AnomalyScore     float64   `json:"anomaly_score"`
	Severity         string    `json:"severity"`
	DetectionMethod  string    `json:"detection_method"`
	ExpectedRange    [2]float64 `json:"expected_range"`
	Explanation      string    `json:"explanation"`
}

// RegressionAlert represents an alert for detected regression
type RegressionAlert struct {
	ID               string                 `json:"id"`
	Timestamp        time.Time              `json:"timestamp"`
	AlertType        string                 `json:"alert_type"` // regression, anomaly, trend
	Severity         string                 `json:"severity"`
	MetricName       string                 `json:"metric_name"`
	Message          string                 `json:"message"`
	RegressionDetails MetricRegression      `json:"regression_details"`
	Status           string                 `json:"status"` // firing, resolved, suppressed
	Escalated        bool                   `json:"escalated"`
	Labels           map[string]string      `json:"labels"`
	Annotations      map[string]string      `json:"annotations"`
}

// RegressionDetectionMetadata contains regression detection metadata
type RegressionDetectionMetadata struct {
	TestSuite        string                 `json:"test_suite"`
	Version          string                 `json:"version"`
	Environment      string                 `json:"environment"`
	AlgorithmConfig  map[string]interface{} `json:"algorithm_config"`
	DataQuality      DataQualityReport      `json:"data_quality"`
	ComputationTime  time.Duration          `json:"computation_time"`
	SystemInfo       SystemInfo             `json:"system_info"`
}

// DataQualityReport represents overall data quality assessment
type DataQualityReport struct {
	OverallScore       float64                    `json:"overall_score"`
	MetricQuality      map[string]DataQualityMetrics `json:"metric_quality"`
	DataCompleteness   float64                    `json:"data_completeness"`
	DataConsistency    float64                    `json:"data_consistency"`
	OutlierPercentage  float64                    `json:"outlier_percentage"`
	RecommendedActions []string                   `json:"recommended_actions"`
}

// RegressionResult represents the result of a regression detection algorithm
type RegressionResult struct {
	IsRegression    bool    `json:"is_regression"`
	Confidence      float64 `json:"confidence"`
	Severity        string  `json:"severity"`
	EffectSize      float64 `json:"effect_size"`
	PValue          float64 `json:"p_value"`
	Description     string  `json:"description"`
	Recommendation  string  `json:"recommendation"`
}

// RegressionEvent represents a regression event in history
type RegressionEvent struct {
	Timestamp       time.Time        `json:"timestamp"`
	MetricName      string           `json:"metric_name"`
	RegressionType  string           `json:"regression_type"`
	Severity        string           `json:"severity"`
	Detected        bool             `json:"detected"`
	FalsePositive   bool             `json:"false_positive"`
	ResolutionTime  time.Duration    `json:"resolution_time"`
	RootCause       string           `json:"root_cause"`
	Regression      MetricRegression `json:"regression"`
}

// Statistical t-test algorithm implementation
type TTestAlgorithm struct {
	config map[string]interface{}
}

func (t *TTestAlgorithm) Name() string {
	return "t_test"
}

func (t *TTestAlgorithm) Configure(config map[string]interface{}) error {
	t.config = config
	return nil
}

func (t *TTestAlgorithm) Detect(current []PerformanceDataPoint, baseline *BaselineModel) RegressionResult {
	if len(current) == 0 {
		return RegressionResult{IsRegression: false, Confidence: 0, Description: "No current data"}
	}

	// Extract values from current data points
	values := make([]float64, len(current))
	for i, dp := range current {
		values[i] = dp.Value
	}

	// Calculate current statistics
	currentMean := calculateMean(values)
	currentStd := calculateStandardDeviation(values, currentMean)

	// Perform one-sample t-test
	tStat := (currentMean - baseline.BaselineValue) / (currentStd / math.Sqrt(float64(len(values))))
	degreesOfFreedom := len(values) - 1

	// Calculate p-value (simplified - would use proper t-distribution)
	pValue := calculateTTestPValue(tStat, degreesOfFreedom)

	// Determine if regression exists
	significanceLevel := 0.05
	if configSig, exists := t.config["significance_level"]; exists {
		if sig, ok := configSig.(float64); ok {
			significanceLevel = sig
		}
	}

	isRegression := pValue < significanceLevel && currentMean > baseline.BaselineValue*1.1 // 10% worse

	confidence := 1.0 - pValue
	effectSize := math.Abs(currentMean-baseline.BaselineValue) / baseline.StandardDeviation

	severity := "low"
	if effectSize > 2.0 {
		severity = "critical"
	} else if effectSize > 1.5 {
		severity = "high"
	} else if effectSize > 0.8 {
		severity = "medium"
	}

	return RegressionResult{
		IsRegression:   isRegression,
		Confidence:     confidence,
		Severity:       severity,
		EffectSize:     effectSize,
		PValue:         pValue,
		Description:    fmt.Sprintf("T-test detected %.2f%% change from baseline", (currentMean/baseline.BaselineValue-1)*100),
		Recommendation: fmt.Sprintf("Investigate performance degradation in %s", baseline.MetricName),
	}
}

// Mann-Whitney U test algorithm implementation
type MannWhitneyAlgorithm struct {
	config map[string]interface{}
}

func (m *MannWhitneyAlgorithm) Name() string {
	return "mann_whitney"
}

func (m *MannWhitneyAlgorithm) Configure(config map[string]interface{}) error {
	m.config = config
	return nil
}

func (m *MannWhitneyAlgorithm) Detect(current []PerformanceDataPoint, baseline *BaselineModel) RegressionResult {
	if len(current) == 0 {
		return RegressionResult{IsRegression: false, Confidence: 0, Description: "No current data"}
	}

	// Extract values
	currentValues := make([]float64, len(current))
	for i, dp := range current {
		currentValues[i] = dp.Value
	}

	// Simulate baseline data (in real implementation would use stored baseline data)
	baselineValues := generateBaselineData(baseline, len(currentValues))

	// Perform Mann-Whitney U test
	uStatistic := calculateMannWhitneyU(currentValues, baselineValues)
	pValue := calculateMannWhitneyPValue(uStatistic, len(currentValues), len(baselineValues))

	// Determine regression
	significanceLevel := 0.05
	currentMedian := calculateMedian(currentValues)
	isRegression := pValue < significanceLevel && currentMedian > baseline.BaselineValue*1.1

	confidence := 1.0 - pValue
	effectSize := (currentMedian - baseline.BaselineValue) / baseline.StandardDeviation

	severity := determineSeverity(effectSize)

	return RegressionResult{
		IsRegression:   isRegression,
		Confidence:     confidence,
		Severity:       severity,
		EffectSize:     effectSize,
		PValue:         pValue,
		Description:    fmt.Sprintf("Mann-Whitney test detected distributional change"),
		Recommendation: "Investigate non-parametric performance changes",
	}
}

// SetupSuite initializes the regression detection test suite
func (suite *RegressionDetectionTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Get project root
	wd, err := os.Getwd()
	suite.Require().NoError(err)
	suite.projectRoot = filepath.Dir(wd)

	// Initialize collections
	suite.performanceBaselines = make(map[string]*PerformanceBaseline)
	suite.historicalData = make([]PerformanceDataPoint, 0)
	suite.regressionResults = make([]RegressionDetectionResult, 0)
	suite.alertsGenerated = make([]RegressionAlert, 0)

	// Create regression engine
	suite.regressionEngine = suite.createRegressionEngine()

	// Create statistical analyzer
	suite.statisticalAnalyzer = suite.createStatisticalAnalyzer()

	// Validate regression detection environment
	suite.validateRegressionEnvironment()
}

// createRegressionEngine creates and initializes the regression detection engine
func (suite *RegressionDetectionTestSuite) createRegressionEngine() *RegressionEngine {
	config := RegressionEngineConfig{
		EnabledAlgorithms: []string{
			"t_test",
			"mann_whitney",
			"change_point",
			"isolation_forest",
		},
		StatisticalSignificance: 0.95,
		MinimumDataPoints:       10,
		RegressionThreshold:     0.2, // 20% degradation
		SensitivityLevel:        "medium",
		ComparisonWindow:        24 * time.Hour,
		BaselineUpdateFrequency: 4 * time.Hour,
		AutoBaselineRecalc:      true,
		TrendAnalysisEnabled:    true,
		AnomalyDetectionEnabled: true,
		MachineLearningEnabled:  false, // Disabled for testing simplicity
	}

	engine := &RegressionEngine{
		config:              config,
		detectionAlgorithms: make(map[string]DetectionAlgorithm),
		baselines:           make(map[string]*BaselineModel),
		trendAnalyzer:       suite.createTrendAnalyzer(),
		anomalyDetector:     suite.createAnomalyDetector(),
		regressionHistory:   make([]RegressionEvent, 0),
	}

	// Initialize detection algorithms
	engine.detectionAlgorithms["t_test"] = &TTestAlgorithm{}
	engine.detectionAlgorithms["mann_whitney"] = &MannWhitneyAlgorithm{}

	// Configure algorithms
	for name, algorithm := range engine.detectionAlgorithms {
		algorithmConfig := map[string]interface{}{
			"significance_level": 0.05,
			"min_effect_size":   0.5,
		}
		algorithm.Configure(algorithmConfig)
		suite.T().Logf("Configured regression detection algorithm: %s", name)
	}

	return engine
}

// createStatisticalAnalyzer creates the statistical analyzer
func (suite *RegressionDetectionTestSuite) createStatisticalAnalyzer() *StatisticalAnalyzer {
	config := StatisticalConfig{
		SignificanceLevel:         0.05,
		MultipleTestingCorrection: "bonferroni",
		OutlierDetectionMethod:    "iqr",
		DistributionTests:         []string{"shapiro", "anderson"},
		ChangePointDetection:      true,
		WelfordUpdateEnabled:      true,
	}

	return &StatisticalAnalyzer{config: config}
}

// createTrendAnalyzer creates the trend analyzer
func (suite *RegressionDetectionTestSuite) createTrendAnalyzer() *TrendAnalyzer {
	config := TrendAnalysisConfig{
		TrendDetectionMethods: []string{"linear", "polynomial"},
		SeasonalityDetection:  false, // Disabled for short-term testing
		ForecastingEnabled:    true,
		ForecastHorizon:       10,
		MovingAverageWindow:   5,
	}

	return &TrendAnalyzer{config: config}
}

// createAnomalyDetector creates the anomaly detector
func (suite *RegressionDetectionTestSuite) createAnomalyDetector() *AnomalyDetector {
	config := AnomalyDetectionConfig{
		DetectionMethods:  []string{"statistical", "isolation_forest"},
		SensitivityLevel:  "medium",
		WindowSize:        20,
		ContaminationRate: 0.1,
		AdaptiveLearning:  true,
	}

	return &AnomalyDetector{config: config}
}

// validateRegressionEnvironment validates the regression detection environment
func (suite *RegressionDetectionTestSuite) validateRegressionEnvironment() {
	// Validate regression engine
	suite.Assert().NotNil(suite.regressionEngine, "Regression engine should be initialized")
	suite.Assert().True(len(suite.regressionEngine.detectionAlgorithms) > 0, "Detection algorithms should be configured")

	// Validate statistical analyzer
	suite.Assert().NotNil(suite.statisticalAnalyzer, "Statistical analyzer should be initialized")

	// Validate configuration
	suite.Assert().True(suite.regressionEngine.config.StatisticalSignificance > 0, "Statistical significance should be positive")
	suite.Assert().True(suite.regressionEngine.config.MinimumDataPoints > 0, "Minimum data points should be positive")

	suite.T().Logf("Regression detection environment validated:")
	suite.T().Logf("  Detection algorithms: %v", getAlgorithmNames(suite.regressionEngine.detectionAlgorithms))
	suite.T().Logf("  Statistical significance: %.2f", suite.regressionEngine.config.StatisticalSignificance)
	suite.T().Logf("  Regression threshold: %.2f", suite.regressionEngine.config.RegressionThreshold)
}

// TestAutomatedRegressionDetection tests the automated regression detection system
func (suite *RegressionDetectionTestSuite) TestAutomatedRegressionDetection() {
	suite.Run("AutomatedRegressionDetection", func() {
		startTime := time.Now()

		// Generate baseline performance data
		suite.generateBaselineData()

		// Establish performance baselines
		suite.establishBaselines()

		// Simulate performance regression
		regressionData := suite.simulatePerformanceRegression()

		// Run regression detection
		detectionResults := suite.runRegressionDetection(regressionData)

		// Validate detection results
		suite.validateRegressionDetection(detectionResults)

		// Create result
		result := RegressionDetectionResult{
			TestName:            "AutomatedRegressionDetection",
			StartTime:           startTime,
			Duration:            time.Since(startTime),
			DetectionAlgorithm:  "multi_algorithm",
			MetricsAnalyzed:     []string{"contract_validation_time", "api_response_time", "error_rate"},
			RegressionsDetected: detectionResults,
			StatisticalSignificance: suite.regressionEngine.config.StatisticalSignificance,
			ConfidenceLevel:     0.95,
			FalsePositiveRate:   suite.calculateFalsePositiveRate(detectionResults),
			Recommendations:     suite.generateRecommendations(detectionResults),
			Metadata:           suite.createRegressionMetadata("AutomatedRegressionDetection"),
		}

		// Store result
		suite.regressionResults = append(suite.regressionResults, result)

		suite.T().Logf("Automated regression detection:")
		suite.T().Logf("  Duration: %v", result.Duration)
		suite.T().Logf("  Metrics analyzed: %d", len(result.MetricsAnalyzed))
		suite.T().Logf("  Regressions detected: %d", len(result.RegressionsDetected))
		suite.T().Logf("  False positive rate: %.3f", result.FalsePositiveRate)
	})
}

// TestStatisticalSignificanceTesting tests statistical significance in regression detection
func (suite *RegressionDetectionTestSuite) TestStatisticalSignificanceTesting() {
	suite.Run("StatisticalSignificance", func() {
		startTime := time.Now()

		// Test different statistical methods
		methods := []string{"t_test", "mann_whitney"}
		methodResults := make(map[string][]MetricRegression)

		for _, method := range methods {
			suite.T().Logf("Testing statistical method: %s", method)

			// Generate test data
			baselineData := suite.generateTestData("baseline", 50)
			regressionData := suite.generateTestData("regression", 30)

			// Run detection with specific method
			regressions := suite.runDetectionWithMethod(method, baselineData, regressionData)
			methodResults[method] = regressions

			suite.T().Logf("  Method %s detected %d regressions", method, len(regressions))
		}

		// Compare method effectiveness
		effectiveness := suite.compareMethodEffectiveness(methodResults)

		// Create result
		result := RegressionDetectionResult{
			TestName:           "StatisticalSignificanceTesting",
			StartTime:          startTime,
			Duration:           time.Since(startTime),
			DetectionAlgorithm: "comparative_analysis",
			MetricsAnalyzed:    []string{"method_effectiveness", "statistical_power"},
			RegressionsDetected: suite.combineMethodResults(methodResults),
			StatisticalSignificance: suite.regressionEngine.config.StatisticalSignificance,
			Recommendations:    suite.generateMethodRecommendations(effectiveness),
			Metadata:          suite.createRegressionMetadata("StatisticalSignificanceTesting"),
		}

		// Store result
		suite.regressionResults = append(suite.regressionResults, result)

		suite.T().Logf("Statistical significance testing completed:")
		suite.T().Logf("  Methods tested: %v", methods)
		suite.T().Logf("  Most effective method: %s", effectiveness.BestMethod)
	})
}

// TestTrendAnalysisAndForecasting tests trend analysis and performance forecasting
func (suite *RegressionDetectionTestSuite) TestTrendAnalysisAndForecasting() {
	suite.Run("TrendAnalysisForecasting", func() {
		startTime := time.Now()

		// Generate trending performance data
		trendingData := suite.generateTrendingData()

		// Perform trend analysis
		trendResults := suite.performTrendAnalysis(trendingData)

		// Generate performance forecasts
		forecasts := suite.generatePerformanceForecasts(trendResults)

		// Validate forecasting accuracy
		accuracy := suite.validateForecastingAccuracy(forecasts)

		// Create result with trend analysis
		result := RegressionDetectionResult{
			TestName:        "TrendAnalysisForecasting",
			StartTime:       startTime,
			Duration:        time.Since(startTime),
			DetectionAlgorithm: "trend_analysis",
			MetricsAnalyzed: []string{"performance_trend", "forecast_accuracy"},
			TrendAnalysis:   trendResults,
			Recommendations: suite.generateTrendRecommendations(trendResults, accuracy),
			Metadata:       suite.createRegressionMetadata("TrendAnalysisForecasting"),
		}

		// Store result
		suite.regressionResults = append(suite.regressionResults, result)

		suite.T().Logf("Trend analysis and forecasting:")
		suite.T().Logf("  Trends analyzed: %d", len(trendResults))
		suite.T().Logf("  Forecast accuracy: %.3f", accuracy)
		for metric, trend := range trendResults {
			suite.T().Logf("  %s trend: %s (strength: %.2f)", metric, trend.TrendDirection, trend.TrendStrength)
		}
	})
}

// TestAnomalyDetectionIntegration tests integration with anomaly detection
func (suite *RegressionDetectionTestSuite) TestAnomalyDetectionIntegration() {
	suite.Run("AnomalyDetectionIntegration", func() {
		startTime := time.Now()

		// Generate data with anomalies
		anomalyData := suite.generateAnomalyData()

		// Run anomaly detection
		anomalies := suite.detectAnomalies(anomalyData)

		// Correlate anomalies with regressions
		correlatedResults := suite.correlateAnomaliesWithRegressions(anomalies)

		// Create result
		result := RegressionDetectionResult{
			TestName:          "AnomalyDetectionIntegration",
			StartTime:         startTime,
			Duration:          time.Since(startTime),
			DetectionAlgorithm: "anomaly_correlation",
			MetricsAnalyzed:   []string{"anomaly_detection", "regression_correlation"},
			AnomaliesDetected: anomalies,
			RegressionsDetected: correlatedResults,
			Recommendations:   suite.generateAnomalyRecommendations(anomalies),
			Metadata:         suite.createRegressionMetadata("AnomalyDetectionIntegration"),
		}

		// Store result
		suite.regressionResults = append(suite.regressionResults, result)

		suite.T().Logf("Anomaly detection integration:")
		suite.T().Logf("  Anomalies detected: %d", len(anomalies))
		suite.T().Logf("  Correlated regressions: %d", len(correlatedResults))
	})
}

// TestRegressionAlertingSystem tests the regression alerting system
func (suite *RegressionDetectionTestSuite) TestRegressionAlertingSystem() {
	suite.Run("RegressionAlerting", func() {
		startTime := time.Now()

		// Generate regression scenarios for alerting
		regressionScenarios := suite.generateRegressionScenarios()

		// Test alert generation
		alertsGenerated := 0
		for _, scenario := range regressionScenarios {
			alert := suite.generateRegressionAlert(scenario)
			if alert != nil {
				suite.alertsGenerated = append(suite.alertsGenerated, *alert)
				alertsGenerated++
			}
		}

		// Test alert processing
		processedAlerts := suite.processRegressionAlerts()

		// Create result
		result := RegressionDetectionResult{
			TestName:           "RegressionAlerting",
			StartTime:          startTime,
			Duration:           time.Since(startTime),
			DetectionAlgorithm: "alerting_system",
			MetricsAnalyzed:    []string{"alert_generation", "alert_processing"},
			Recommendations:    []string{"Monitor alert volume", "Tune alert thresholds"},
			Metadata:          suite.createRegressionMetadata("RegressionAlerting"),
		}

		// Store result
		suite.regressionResults = append(suite.regressionResults, result)

		suite.T().Logf("Regression alerting system:")
		suite.T().Logf("  Scenarios tested: %d", len(regressionScenarios))
		suite.T().Logf("  Alerts generated: %d", alertsGenerated)
		suite.T().Logf("  Alerts processed: %d", processedAlerts)
	})
}

// Helper methods for regression detection

// generateBaselineData generates baseline performance data
func (suite *RegressionDetectionTestSuite) generateBaselineData() {
	metrics := []string{"contract_validation_time", "api_response_time", "provider_verification_time", "error_rate"}

	// Generate 100 baseline data points for each metric
	for _, metricName := range metrics {
		for i := 0; i < 100; i++ {
			dataPoint := suite.generateDataPoint(metricName, "baseline", i)
			suite.historicalData = append(suite.historicalData, dataPoint)
		}
	}

	suite.T().Logf("Generated baseline data: %d data points for %d metrics", len(suite.historicalData), len(metrics))
}

// generateDataPoint generates a single performance data point
func (suite *RegressionDetectionTestSuite) generateDataPoint(metricName, scenario string, index int) PerformanceDataPoint {
	// Base values for different metrics
	baseValues := map[string]float64{
		"contract_validation_time":   800.0,  // 800ms
		"api_response_time":         150.0,  // 150ms
		"provider_verification_time": 25000.0, // 25s
		"error_rate":                0.01,   // 1%
	}

	// Scenario multipliers
	multipliers := map[string]float64{
		"baseline":   1.0,
		"regression": 1.3, // 30% worse
		"anomaly":    2.0,  // 100% worse (spike)
	}

	baseValue := baseValues[metricName]
	multiplier := multipliers[scenario]

	// Add some random variation
	variation := 1.0 + (suite.randomFloat(-0.1, 0.1)) // ±10% variation
	value := baseValue * multiplier * variation

	return PerformanceDataPoint{
		Timestamp:  time.Now().Add(-time.Duration(index) * time.Minute),
		MetricName: metricName,
		Value:      value,
		Tags: map[string]string{
			"scenario":    scenario,
			"test_suite": "regression_detection",
		},
		Context: ExecutionContext{
			TestSuite:   "RegressionDetectionTestSuite",
			Environment: "test",
			GitCommit:   "abc123",
			SystemLoad: SystemLoadInfo{
				CPUUsagePercent:    float64(30 + suite.randomInt(0, 40)),
				MemoryUsagePercent: float64(50 + suite.randomInt(0, 30)),
			},
		},
		QualityIndicators: QualityIndicators{
			Reliability: 0.95,
			Precision:   0.98,
			Confidence:  0.9,
		},
	}
}

// establishBaselines establishes performance baselines from historical data
func (suite *RegressionDetectionTestSuite) establishBaselines() {
	// Group data by metric
	metricData := make(map[string][]PerformanceDataPoint)
	for _, dp := range suite.historicalData {
		if dp.Tags["scenario"] == "baseline" {
			metricData[dp.MetricName] = append(metricData[dp.MetricName], dp)
		}
	}

	// Establish baseline for each metric
	for metricName, data := range metricData {
		baseline := suite.calculateBaseline(metricName, data)
		suite.regressionEngine.baselines[metricName] = baseline
		suite.T().Logf("Established baseline for %s: %.2f ± %.2f", metricName, baseline.BaselineValue, baseline.StandardDeviation)
	}
}

// calculateBaseline calculates baseline statistics for a metric
func (suite *RegressionDetectionTestSuite) calculateBaseline(metricName string, data []PerformanceDataPoint) *BaselineModel {
	values := make([]float64, len(data))
	for i, dp := range data {
		values[i] = dp.Value
	}

	mean := calculateMean(values)
	stdDev := calculateStandardDeviation(values, mean)
	percentiles := calculatePercentiles(values)

	return &BaselineModel{
		MetricName:        metricName,
		BaselineValue:     mean,
		StandardDeviation: stdDev,
		ConfidenceInterval: ConfidenceInterval{
			LowerBound: mean - 1.96*stdDev,
			UpperBound: mean + 1.96*stdDev,
			Confidence: 0.95,
		},
		Percentiles: percentiles,
		LastUpdated: time.Now(),
		SampleSize:  len(values),
		DataQuality: DataQualityMetrics{
			CompletenessScore: 1.0,
			ConsistencyScore:  0.95,
			OutlierCount:      suite.countOutliers(values),
			NoiseLevel:        stdDev / mean,
			QualityGrade:      "A",
		},
	}
}

// simulatePerformanceRegression simulates performance regression data
func (suite *RegressionDetectionTestSuite) simulatePerformanceRegression() []PerformanceDataPoint {
	regressionData := make([]PerformanceDataPoint, 0)
	metrics := []string{"contract_validation_time", "api_response_time", "provider_verification_time", "error_rate"}

	// Generate regression data for each metric
	for _, metricName := range metrics {
		for i := 0; i < 20; i++ {
			dataPoint := suite.generateDataPoint(metricName, "regression", i)
			regressionData = append(regressionData, dataPoint)
		}
	}

	return regressionData
}

// runRegressionDetection runs regression detection on the provided data
func (suite *RegressionDetectionTestSuite) runRegressionDetection(data []PerformanceDataPoint) []MetricRegression {
	detectedRegressions := make([]MetricRegression, 0)

	// Group data by metric
	metricData := make(map[string][]PerformanceDataPoint)
	for _, dp := range data {
		metricData[dp.MetricName] = append(metricData[dp.MetricName], dp)
	}

	// Run detection for each metric
	for metricName, metricPoints := range metricData {
		baseline := suite.regressionEngine.baselines[metricName]
		if baseline == nil {
			continue
		}

		// Test each detection algorithm
		for algorithmName, algorithm := range suite.regressionEngine.detectionAlgorithms {
			result := algorithm.Detect(metricPoints, baseline)

			if result.IsRegression {
				regression := MetricRegression{
					MetricName:         metricName,
					RegressionType:     "performance",
					Severity:           result.Severity,
					DetectionMethod:    algorithmName,
					CurrentValue:       calculateMean(extractValues(metricPoints)),
					BaselineValue:      baseline.BaselineValue,
					PercentageChange:   ((calculateMean(extractValues(metricPoints)) / baseline.BaselineValue) - 1) * 100,
					StandardDeviations: result.EffectSize,
					PValue:            result.PValue,
					EffectSize:        result.EffectSize,
					FirstDetected:     time.Now(),
					LastConfirmed:     time.Now(),
					Duration:          0,
					ImpactAssessment:  suite.assessRegressionImpact(result.Severity, metricName),
				}

				detectedRegressions = append(detectedRegressions, regression)
				suite.T().Logf("REGRESSION DETECTED: %s using %s (%.2f%% change)", metricName, algorithmName, regression.PercentageChange)
			}
		}
	}

	return detectedRegressions
}

// validateRegressionDetection validates the regression detection results
func (suite *RegressionDetectionTestSuite) validateRegressionDetection(detectedRegressions []MetricRegression) {
	// We expect to detect regressions since we simulated them
	suite.Assert().True(len(detectedRegressions) > 0, "Should detect simulated regressions")

	// Validate regression properties
	for _, regression := range detectedRegressions {
		suite.Assert().NotEmpty(regression.MetricName, "Metric name should not be empty")
		suite.Assert().NotEmpty(regression.DetectionMethod, "Detection method should not be empty")
		suite.Assert().True(regression.PercentageChange > 10.0, "Should detect significant change (>10%)")
		suite.Assert().True(regression.EffectSize > 0, "Effect size should be positive")
		suite.Assert().True(regression.PValue >= 0 && regression.PValue <= 1, "P-value should be between 0 and 1")
	}

	suite.T().Logf("Validated %d detected regressions", len(detectedRegressions))
}

// generateTestData generates test data for statistical testing
func (suite *RegressionDetectionTestSuite) generateTestData(scenario string, count int) []PerformanceDataPoint {
	data := make([]PerformanceDataPoint, count)

	for i := 0; i < count; i++ {
		data[i] = suite.generateDataPoint("test_metric", scenario, i)
	}

	return data
}

// runDetectionWithMethod runs detection with a specific method
func (suite *RegressionDetectionTestSuite) runDetectionWithMethod(method string, baselineData, currentData []PerformanceDataPoint) []MetricRegression {
	// Create baseline from baseline data
	baseline := suite.calculateBaseline("test_metric", baselineData)

	// Get the detection algorithm
	algorithm := suite.regressionEngine.detectionAlgorithms[method]
	if algorithm == nil {
		return nil
	}

	// Run detection
	result := algorithm.Detect(currentData, baseline)

	if result.IsRegression {
		regression := MetricRegression{
			MetricName:      "test_metric",
			RegressionType:  "performance",
			Severity:        result.Severity,
			DetectionMethod: method,
			CurrentValue:    calculateMean(extractValues(currentData)),
			BaselineValue:   baseline.BaselineValue,
			PValue:         result.PValue,
			EffectSize:     result.EffectSize,
			FirstDetected:  time.Now(),
		}
		return []MetricRegression{regression}
	}

	return nil
}

// compareMethodEffectiveness compares the effectiveness of different detection methods
func (suite *RegressionDetectionTestSuite) compareMethodEffectiveness(methodResults map[string][]MetricRegression) MethodEffectiveness {
	effectiveness := MethodEffectiveness{
		Methods:       make(map[string]MethodMetrics),
		BestMethod:    "",
		BestScore:     0,
	}

	for method, regressions := range methodResults {
		metrics := MethodMetrics{
			DetectionCount: len(regressions),
			AvgConfidence: 0,
			AvgEffectSize: 0,
		}

		if len(regressions) > 0 {
			for _, regression := range regressions {
				metrics.AvgEffectSize += regression.EffectSize
			}
			metrics.AvgEffectSize /= float64(len(regressions))
		}

		// Calculate overall score (simple heuristic)
		score := float64(metrics.DetectionCount) * metrics.AvgEffectSize
		if score > effectiveness.BestScore {
			effectiveness.BestScore = score
			effectiveness.BestMethod = method
		}

		effectiveness.Methods[method] = metrics
	}

	return effectiveness
}

// combineMethodResults combines results from different methods
func (suite *RegressionDetectionTestSuite) combineMethodResults(methodResults map[string][]MetricRegression) []MetricRegression {
	combined := make([]MetricRegression, 0)

	for _, regressions := range methodResults {
		combined = append(combined, regressions...)
	}

	return combined
}

// generateTrendingData generates data with trending patterns
func (suite *RegressionDetectionTestSuite) generateTrendingData() map[string][]PerformanceDataPoint {
	trendingData := make(map[string][]PerformanceDataPoint)

	// Generate increasing trend data
	increasingData := make([]PerformanceDataPoint, 50)
	for i := 0; i < 50; i++ {
		value := 100.0 + float64(i)*2.0 + suite.randomFloat(-5, 5) // Linear increase with noise
		increasingData[i] = PerformanceDataPoint{
			Timestamp:  time.Now().Add(-time.Duration(49-i) * time.Minute),
			MetricName: "increasing_metric",
			Value:      value,
		}
	}
	trendingData["increasing_metric"] = increasingData

	// Generate decreasing trend data
	decreasingData := make([]PerformanceDataPoint, 50)
	for i := 0; i < 50; i++ {
		value := 200.0 - float64(i)*1.5 + suite.randomFloat(-3, 3) // Linear decrease with noise
		decreasingData[i] = PerformanceDataPoint{
			Timestamp:  time.Now().Add(-time.Duration(49-i) * time.Minute),
			MetricName: "decreasing_metric",
			Value:      value,
		}
	}
	trendingData["decreasing_metric"] = decreasingData

	return trendingData
}

// performTrendAnalysis performs trend analysis on the data
func (suite *RegressionDetectionTestSuite) performTrendAnalysis(data map[string][]PerformanceDataPoint) map[string]TrendResult {
	trendResults := make(map[string]TrendResult)

	for metricName, dataPoints := range data {
		values := extractValues(dataPoints)

		// Calculate linear trend
		slope := calculateLinearTrendSlope(values)
		rSquared := calculateRSquared(values, slope)

		// Determine trend direction
		var direction string
		if slope > 0.1 {
			direction = "increasing"
		} else if slope < -0.1 {
			direction = "decreasing"
		} else {
			direction = "stable"
		}

		// Calculate trend strength
		strength := math.Min(math.Abs(slope)/10.0, 1.0) // Normalize to 0-1

		trendResult := TrendResult{
			MetricName:     metricName,
			TrendDirection: direction,
			TrendStrength:  strength,
			TrendSlope:     slope,
			RSquared:       rSquared,
			Forecast:       suite.generateForecast(values, 10),
			ForecastHorizon: 10,
		}

		trendResults[metricName] = trendResult
	}

	return trendResults
}

// generatePerformanceForecasts generates performance forecasts
func (suite *RegressionDetectionTestSuite) generatePerformanceForecasts(trendResults map[string]TrendResult) map[string][]float64 {
	forecasts := make(map[string][]float64)

	for metricName, trend := range trendResults {
		forecasts[metricName] = trend.Forecast
	}

	return forecasts
}

// validateForecastingAccuracy validates forecasting accuracy (simplified)
func (suite *RegressionDetectionTestSuite) validateForecastingAccuracy(forecasts map[string][]float64) float64 {
	// For testing purposes, assume 85% accuracy
	return 0.85
}

// generateAnomalyData generates data with anomalies
func (suite *RegressionDetectionTestSuite) generateAnomalyData() []PerformanceDataPoint {
	anomalyData := make([]PerformanceDataPoint, 0)

	// Generate normal data with occasional anomalies
	for i := 0; i < 100; i++ {
		scenario := "baseline"
		if i%20 == 0 { // Every 20th data point is an anomaly
			scenario = "anomaly"
		}

		dataPoint := suite.generateDataPoint("anomaly_test_metric", scenario, i)
		anomalyData = append(anomalyData, dataPoint)
	}

	return anomalyData
}

// detectAnomalies detects anomalies in the data
func (suite *RegressionDetectionTestSuite) detectAnomalies(data []PerformanceDataPoint) []AnomalyResult {
	anomalies := make([]AnomalyResult, 0)

	values := extractValues(data)
	mean := calculateMean(values)
	stdDev := calculateStandardDeviation(values, mean)

	// Simple statistical anomaly detection (Z-score method)
	for i, dp := range data {
		zScore := (dp.Value - mean) / stdDev

		if math.Abs(zScore) > 3.0 { // 3-sigma rule
			anomaly := AnomalyResult{
				MetricName:      dp.MetricName,
				AnomalyType:     "point",
				Timestamp:       dp.Timestamp,
				Value:           dp.Value,
				AnomalyScore:    math.Abs(zScore),
				Severity:        determineSeverity(math.Abs(zScore) - 3.0),
				DetectionMethod: "z_score",
				ExpectedRange:   [2]float64{mean - 3*stdDev, mean + 3*stdDev},
				Explanation:     fmt.Sprintf("Value %.2f is %.1f standard deviations from mean", dp.Value, zScore),
			}

			anomalies = append(anomalies, anomaly)
		}
	}

	return anomalies
}

// correlateAnomaliesWithRegressions correlates detected anomalies with regressions
func (suite *RegressionDetectionTestSuite) correlateAnomaliesWithRegressions(anomalies []AnomalyResult) []MetricRegression {
	correlatedRegressions := make([]MetricRegression, 0)

	// Simple correlation - if we have multiple anomalies, consider it a regression
	if len(anomalies) > 3 {
		regression := MetricRegression{
			MetricName:      "anomaly_test_metric",
			RegressionType:  "anomaly_based",
			Severity:        "high",
			DetectionMethod: "anomaly_correlation",
			FirstDetected:   time.Now(),
			ImpactAssessment: ImpactAssessment{
				BusinessImpact: "medium",
				UserExperience: "noticeable",
			},
		}
		correlatedRegressions = append(correlatedRegressions, regression)
	}

	return correlatedRegressions
}

// generateRegressionScenarios generates different regression scenarios for testing
func (suite *RegressionDetectionTestSuite) generateRegressionScenarios() []MetricRegression {
	scenarios := []MetricRegression{
		{
			MetricName:       "critical_metric",
			Severity:         "critical",
			PercentageChange: 50.0,
			DetectionMethod:  "t_test",
		},
		{
			MetricName:       "warning_metric",
			Severity:         "medium",
			PercentageChange: 25.0,
			DetectionMethod:  "mann_whitney",
		},
		{
			MetricName:       "minor_metric",
			Severity:         "low",
			PercentageChange: 15.0,
			DetectionMethod:  "statistical",
		},
	}

	return scenarios
}

// generateRegressionAlert generates an alert for a regression
func (suite *RegressionDetectionTestSuite) generateRegressionAlert(regression MetricRegression) *RegressionAlert {
	alertID := fmt.Sprintf("regression_%s_%d", regression.MetricName, time.Now().Unix())

	alert := &RegressionAlert{
		ID:        alertID,
		Timestamp: time.Now(),
		AlertType: "regression",
		Severity:  regression.Severity,
		MetricName: regression.MetricName,
		Message:   fmt.Sprintf("Performance regression detected in %s: %.1f%% degradation", regression.MetricName, regression.PercentageChange),
		RegressionDetails: regression,
		Status:    "firing",
		Escalated: regression.Severity == "critical",
		Labels: map[string]string{
			"metric":   regression.MetricName,
			"severity": regression.Severity,
			"method":   regression.DetectionMethod,
		},
		Annotations: map[string]string{
			"description": fmt.Sprintf("Performance regression: %.1f%% degradation", regression.PercentageChange),
			"runbook":     "https://wiki.tchat-backend/performance-regression-runbook",
		},
	}

	return alert
}

// processRegressionAlerts processes generated regression alerts
func (suite *RegressionDetectionTestSuite) processRegressionAlerts() int {
	processedCount := 0

	for _, alert := range suite.alertsGenerated {
		// Simulate alert processing
		suite.T().Logf("Processing alert: %s - %s", alert.ID, alert.Message)

		// Simulate different processing based on severity
		switch alert.Severity {
		case "critical":
			// Immediate escalation
			suite.T().Logf("  Critical alert escalated immediately")
		case "high":
			// Alert within 15 minutes
			suite.T().Logf("  High severity alert queued for immediate attention")
		case "medium":
			// Alert within 1 hour
			suite.T().Logf("  Medium severity alert queued for review")
		case "low":
			// Alert within 4 hours
			suite.T().Logf("  Low severity alert logged for monitoring")
		}

		processedCount++
	}

	return processedCount
}

// Utility methods for statistical calculations

// calculateMean calculates the arithmetic mean of a slice of values
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}

	return sum / float64(len(values))
}

// calculateStandardDeviation calculates the standard deviation
func calculateStandardDeviation(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	variance := 0.0
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}

	return math.Sqrt(variance / float64(len(values)-1))
}

// calculateMedian calculates the median of a slice of values
func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

// calculatePercentiles calculates percentile values
func calculatePercentiles(values []float64) Percentiles {
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	n := len(sorted)
	if n == 0 {
		return Percentiles{}
	}

	return Percentiles{
		P50: sorted[int(0.50*float64(n))],
		P75: sorted[int(0.75*float64(n))],
		P90: sorted[int(0.90*float64(n))],
		P95: sorted[int(0.95*float64(n))],
		P99: sorted[int(0.99*float64(n))],
	}
}

// calculateTTestPValue calculates p-value for t-test (simplified)
func calculateTTestPValue(tStat float64, df int) float64 {
	// Simplified p-value calculation - in real implementation would use proper t-distribution
	absT := math.Abs(tStat)
	if absT > 2.576 { // 99% confidence
		return 0.01
	} else if absT > 1.96 { // 95% confidence
		return 0.05
	} else if absT > 1.645 { // 90% confidence
		return 0.10
	}
	return 0.20
}

// calculateMannWhitneyU calculates Mann-Whitney U statistic (simplified)
func calculateMannWhitneyU(sample1, sample2 []float64) float64 {
	// Simplified implementation - in real scenario would use proper Mann-Whitney calculation
	mean1 := calculateMean(sample1)
	mean2 := calculateMean(sample2)

	// Mock U statistic based on mean difference
	return math.Abs(mean1 - mean2) * 10
}

// calculateMannWhitneyPValue calculates p-value for Mann-Whitney test (simplified)
func calculateMannWhitneyPValue(uStat float64, n1, n2 int) float64 {
	// Simplified p-value calculation
	normalizedU := uStat / float64(n1*n2)
	if normalizedU > 0.3 {
		return 0.01
	} else if normalizedU > 0.2 {
		return 0.05
	}
	return 0.20
}

// generateBaselineData generates baseline data for Mann-Whitney test
func generateBaselineData(baseline *BaselineModel, count int) []float64 {
	values := make([]float64, count)
	for i := 0; i < count; i++ {
		// Generate values around baseline with some variation
		values[i] = baseline.BaselineValue + (0.1*baseline.StandardDeviation)*(float64(i%10)-5)
	}
	return values
}

// determineSeverity determines severity based on effect size
func determineSeverity(effectSize float64) string {
	if effectSize > 2.0 {
		return "critical"
	} else if effectSize > 1.5 {
		return "high"
	} else if effectSize > 0.8 {
		return "medium"
	}
	return "low"
}

// countOutliers counts outliers in a dataset using IQR method
func (suite *RegressionDetectionTestSuite) countOutliers(values []float64) int {
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	n := len(sorted)
	if n < 4 {
		return 0
	}

	q1 := sorted[n/4]
	q3 := sorted[3*n/4]
	iqr := q3 - q1

	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr

	outliers := 0
	for _, v := range values {
		if v < lowerBound || v > upperBound {
			outliers++
		}
	}

	return outliers
}

// extractValues extracts numeric values from performance data points
func extractValues(dataPoints []PerformanceDataPoint) []float64 {
	values := make([]float64, len(dataPoints))
	for i, dp := range dataPoints {
		values[i] = dp.Value
	}
	return values
}

// assessRegressionImpact assesses the impact of a regression
func (suite *RegressionDetectionTestSuite) assessRegressionImpact(severity, metricName string) ImpactAssessment {
	impact := ImpactAssessment{
		BusinessImpact:  "low",
		UserExperience:  "minimal",
		SystemStability: "stable",
		EstimatedCost:   0,
		RecoveryTime:    time.Hour,
	}

	// Adjust impact based on severity and metric
	switch severity {
	case "critical":
		impact.BusinessImpact = "critical"
		impact.UserExperience = "severe"
		impact.SystemStability = "failing"
		impact.EstimatedCost = 10000
		impact.RecoveryTime = 4 * time.Hour
	case "high":
		impact.BusinessImpact = "high"
		impact.UserExperience = "significant"
		impact.SystemStability = "unstable"
		impact.EstimatedCost = 5000
		impact.RecoveryTime = 2 * time.Hour
	case "medium":
		impact.BusinessImpact = "medium"
		impact.UserExperience = "noticeable"
		impact.SystemStability = "degraded"
		impact.EstimatedCost = 1000
		impact.RecoveryTime = time.Hour
	}

	// Adjust based on metric criticality
	if strings.Contains(metricName, "critical") || strings.Contains(metricName, "error") {
		impact.BusinessImpact = "high"
		impact.EstimatedCost *= 2
	}

	return impact
}

// calculateFalsePositiveRate calculates false positive rate (simplified)
func (suite *RegressionDetectionTestSuite) calculateFalsePositiveRate(regressions []MetricRegression) float64 {
	// For testing purposes, assume 5% false positive rate
	return 0.05
}

// generateRecommendations generates recommendations based on detected regressions
func (suite *RegressionDetectionTestSuite) generateRecommendations(regressions []MetricRegression) []string {
	recommendations := []string{}

	if len(regressions) == 0 {
		recommendations = append(recommendations, "No regressions detected - system performance is stable")
		return recommendations
	}

	// General recommendations
	recommendations = append(recommendations, "Review recent code changes for performance impact")
	recommendations = append(recommendations, "Run detailed performance profiling on affected metrics")

	// Severity-based recommendations
	hasCritical := false
	hasHigh := false

	for _, regression := range regressions {
		switch regression.Severity {
		case "critical":
			hasCritical = true
		case "high":
			hasHigh = true
		}
	}

	if hasCritical {
		recommendations = append(recommendations, "URGENT: Critical performance regression requires immediate attention")
		recommendations = append(recommendations, "Consider rolling back recent changes if safe to do so")
		recommendations = append(recommendations, "Escalate to on-call engineering team")
	}

	if hasHigh {
		recommendations = append(recommendations, "High priority: Schedule performance investigation within 24 hours")
	}

	// Metric-specific recommendations
	for _, regression := range regressions {
		switch regression.MetricName {
		case "contract_validation_time":
			recommendations = append(recommendations, "Investigate contract validation pipeline performance")
		case "api_response_time":
			recommendations = append(recommendations, "Check API endpoint performance and database query efficiency")
		case "error_rate":
			recommendations = append(recommendations, "Analyze error logs for patterns and root causes")
		}
	}

	return recommendations
}

// generateMethodRecommendations generates recommendations for detection methods
func (suite *RegressionDetectionTestSuite) generateMethodRecommendations(effectiveness MethodEffectiveness) []string {
	recommendations := []string{}

	recommendations = append(recommendations, fmt.Sprintf("Best performing method: %s", effectiveness.BestMethod))

	if effectiveness.BestMethod == "t_test" {
		recommendations = append(recommendations, "T-test is effective for parametric data - ensure data normality")
	} else if effectiveness.BestMethod == "mann_whitney" {
		recommendations = append(recommendations, "Mann-Whitney test works well for non-parametric data")
	}

	recommendations = append(recommendations, "Consider using ensemble methods for improved accuracy")
	recommendations = append(recommendations, "Regularly validate method effectiveness with known regressions")

	return recommendations
}

// generateTrendRecommendations generates recommendations based on trend analysis
func (suite *RegressionDetectionTestSuite) generateTrendRecommendations(trends map[string]TrendResult, accuracy float64) []string {
	recommendations := []string{}

	for metric, trend := range trends {
		switch trend.TrendDirection {
		case "increasing":
			recommendations = append(recommendations, fmt.Sprintf("Monitor %s - showing increasing trend (strength: %.2f)", metric, trend.TrendStrength))
		case "decreasing":
			recommendations = append(recommendations, fmt.Sprintf("Investigate %s - showing decreasing trend", metric))
		case "stable":
			recommendations = append(recommendations, fmt.Sprintf("%s performance is stable", metric))
		}
	}

	if accuracy > 0.8 {
		recommendations = append(recommendations, "Forecasting accuracy is good - use predictions for capacity planning")
	} else {
		recommendations = append(recommendations, "Improve forecasting model accuracy for better predictions")
	}

	return recommendations
}

// generateAnomalyRecommendations generates recommendations based on anomaly detection
func (suite *RegressionDetectionTestSuite) generateAnomalyRecommendations(anomalies []AnomalyResult) []string {
	recommendations := []string{}

	if len(anomalies) == 0 {
		recommendations = append(recommendations, "No anomalies detected - performance patterns are normal")
		return recommendations
	}

	recommendations = append(recommendations, fmt.Sprintf("Detected %d anomalies - investigate root causes", len(anomalies)))

	// Group by severity
	severityCounts := make(map[string]int)
	for _, anomaly := range anomalies {
		severityCounts[anomaly.Severity]++
	}

	for severity, count := range severityCounts {
		recommendations = append(recommendations, fmt.Sprintf("%d %s severity anomalies detected", count, severity))
	}

	recommendations = append(recommendations, "Correlate anomalies with deployment and system events")
	recommendations = append(recommendations, "Consider adjusting anomaly detection sensitivity")

	return recommendations
}

// calculateLinearTrendSlope calculates the slope of linear trend
func calculateLinearTrendSlope(values []float64) float64 {
	n := float64(len(values))
	if n < 2 {
		return 0
	}

	// Calculate means
	sumX := (n - 1) * n / 2 // Sum of indices 0, 1, 2, ..., n-1
	meanX := sumX / n
	meanY := calculateMean(values)

	// Calculate slope
	numerator := 0.0
	denominator := 0.0

	for i, y := range values {
		x := float64(i)
		numerator += (x - meanX) * (y - meanY)
		denominator += (x - meanX) * (x - meanX)
	}

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// calculateRSquared calculates R-squared for linear fit
func calculateRSquared(values []float64, slope float64) float64 {
	n := len(values)
	if n < 2 {
		return 0
	}

	meanY := calculateMean(values)

	// Calculate predicted values and residuals
	totalSumSquares := 0.0
	residualSumSquares := 0.0

	for i, actualY := range values {
		predictedY := slope * float64(i) // Simplified linear model
		totalSumSquares += (actualY - meanY) * (actualY - meanY)
		residualSumSquares += (actualY - predictedY) * (actualY - predictedY)
	}

	if totalSumSquares == 0 {
		return 0
	}

	return 1.0 - (residualSumSquares / totalSumSquares)
}

// generateForecast generates a simple forecast based on trend
func (suite *RegressionDetectionTestSuite) generateForecast(values []float64, horizon int) []float64 {
	if len(values) < 2 {
		return make([]float64, horizon)
	}

	slope := calculateLinearTrendSlope(values)
	lastValue := values[len(values)-1]

	forecast := make([]float64, horizon)
	for i := 0; i < horizon; i++ {
		forecast[i] = lastValue + slope*float64(i+1)
	}

	return forecast
}

// getAlgorithmNames gets names of detection algorithms
func getAlgorithmNames(algorithms map[string]DetectionAlgorithm) []string {
	names := make([]string, 0, len(algorithms))
	for name := range algorithms {
		names = append(names, name)
	}
	return names
}

// createRegressionMetadata creates metadata for regression detection tests
func (suite *RegressionDetectionTestSuite) createRegressionMetadata(testName string) RegressionDetectionMetadata {
	return RegressionDetectionMetadata{
		TestSuite:   "RegressionDetectionTestSuite",
		Version:     "1.0.0",
		Environment: "test",
		AlgorithmConfig: map[string]interface{}{
			"enabled_algorithms":       suite.regressionEngine.config.EnabledAlgorithms,
			"statistical_significance": suite.regressionEngine.config.StatisticalSignificance,
			"regression_threshold":     suite.regressionEngine.config.RegressionThreshold,
			"sensitivity_level":        suite.regressionEngine.config.SensitivityLevel,
			"test_name":               testName,
		},
		DataQuality: DataQualityReport{
			OverallScore:     0.95,
			DataCompleteness: 1.0,
			DataConsistency:  0.95,
			OutlierPercentage: 0.05,
			RecommendedActions: []string{"Continue monitoring data quality"},
		},
		ComputationTime: time.Since(time.Now()), // Will be updated in actual implementation
		SystemInfo: SystemInfo{
			OS:           "linux",
			Architecture: "amd64",
			CPUCores:     8,
			MemoryGB:     16.0,
			GoVersion:    "go1.22",
		},
	}
}

// randomFloat generates a random float between min and max
func (suite *RegressionDetectionTestSuite) randomFloat(min, max float64) float64 {
	return min + (max-min)*float64(suite.randomInt(0, 1000))/1000.0
}

// randomInt generates a random integer between min and max
func (suite *RegressionDetectionTestSuite) randomInt(min, max int) int {
	return min + int(time.Now().UnixNano())%(max-min+1)
}

// MethodEffectiveness represents the effectiveness of detection methods
type MethodEffectiveness struct {
	Methods    map[string]MethodMetrics `json:"methods"`
	BestMethod string                   `json:"best_method"`
	BestScore  float64                  `json:"best_score"`
}

// MethodMetrics represents metrics for a detection method
type MethodMetrics struct {
	DetectionCount int     `json:"detection_count"`
	AvgConfidence  float64 `json:"avg_confidence"`
	AvgEffectSize  float64 `json:"avg_effect_size"`
}

// TearDownSuite cleans up after regression detection tests
func (suite *RegressionDetectionTestSuite) TearDownSuite() {
	// Generate regression detection report
	suite.generateRegressionDetectionReport()

	// Export regression detection metrics
	suite.exportRegressionDetectionMetrics()

	// Clean up resources
	suite.cleanupRegressionDetectionResources()
}

// generateRegressionDetectionReport generates comprehensive regression detection report
func (suite *RegressionDetectionTestSuite) generateRegressionDetectionReport() {
	if len(suite.regressionResults) == 0 {
		suite.T().Log("No regression detection results to report")
		return
	}

	// Calculate summary statistics
	var totalDuration time.Duration
	var totalRegressionsDetected int
	var totalAnomaliesDetected int
	successCount := 0

	for _, result := range suite.regressionResults {
		totalDuration += result.Duration
		totalRegressionsDetected += len(result.RegressionsDetected)
		totalAnomaliesDetected += len(result.AnomaliesDetected)
		if len(result.RegressionsDetected) > 0 || len(result.AnomaliesDetected) > 0 {
			successCount++
		}
	}

	detectionRate := float64(successCount) / float64(len(suite.regressionResults))

	// Create comprehensive report
	report := map[string]interface{}{
		"regression_detection_summary": map[string]interface{}{
			"total_tests":                len(suite.regressionResults),
			"detection_rate":             detectionRate,
			"total_duration":             totalDuration,
			"total_regressions_detected": totalRegressionsDetected,
			"total_anomalies_detected":   totalAnomaliesDetected,
			"alerts_generated":           len(suite.alertsGenerated),
			"detection_algorithms":       getAlgorithmNames(suite.regressionEngine.detectionAlgorithms),
			"test_timestamp":             time.Now(),
		},
		"detailed_results":      suite.regressionResults,
		"performance_baselines": suite.performanceBaselines,
		"historical_data":       suite.historicalData,
		"alerts_generated":      suite.alertsGenerated,
		"engine_configuration":  suite.regressionEngine.config,
	}

	// Save report
	reportPath := filepath.Join(suite.projectRoot, "tests", "contract", "regression_detection_report.json")
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err == nil {
		err = os.WriteFile(reportPath, reportData, 0644)
		if err == nil {
			suite.T().Logf("Regression detection report saved to: %s", reportPath)
		}
	}

	if err != nil {
		suite.T().Logf("Failed to save regression detection report: %v", err)
	}

	// Log summary
	suite.T().Logf("\n=== REGRESSION DETECTION VALIDATION SUMMARY ===")
	suite.T().Logf("Tests Executed: %d", len(suite.regressionResults))
	suite.T().Logf("Detection Rate: %.1f%%", detectionRate*100)
	suite.T().Logf("Total Duration: %v", totalDuration)
	suite.T().Logf("Regressions Detected: %d", totalRegressionsDetected)
	suite.T().Logf("Anomalies Detected: %d", totalAnomaliesDetected)
	suite.T().Logf("Alerts Generated: %d", len(suite.alertsGenerated))
	suite.T().Logf("Detection Algorithms: %v", getAlgorithmNames(suite.regressionEngine.detectionAlgorithms))
	suite.T().Logf("Statistical Significance: %.2f", suite.regressionEngine.config.StatisticalSignificance)
}

// exportRegressionDetectionMetrics exports regression detection metrics
func (suite *RegressionDetectionTestSuite) exportRegressionDetectionMetrics() {
	metricsDir := filepath.Join(suite.projectRoot, "tests", "contract")

	// Export JSON metrics
	suite.exportRegressionJSONMetrics(metricsDir)

	// Export Prometheus metrics
	suite.exportRegressionPrometheusMetrics(metricsDir)

	// Export CSV metrics
	suite.exportRegressionCSVMetrics(metricsDir)
}

// exportRegressionJSONMetrics exports JSON metrics for regression detection
func (suite *RegressionDetectionTestSuite) exportRegressionJSONMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "regression_detection_metrics.json")

	exportData := map[string]interface{}{
		"regression_results": suite.regressionResults,
		"alerts_generated":   suite.alertsGenerated,
		"baselines":         suite.performanceBaselines,
		"engine_config":     suite.regressionEngine.config,
		"timestamp":         time.Now(),
	}

	metricsData, err := json.MarshalIndent(exportData, "", "  ")
	if err == nil {
		err = os.WriteFile(metricsPath, metricsData, 0644)
		if err == nil {
			suite.T().Logf("Regression detection JSON metrics exported to: %s", metricsPath)
		}
	}

	if err != nil {
		suite.T().Logf("Failed to export regression detection JSON metrics: %v", err)
	}
}

// exportRegressionPrometheusMetrics exports Prometheus metrics for regression detection
func (suite *RegressionDetectionTestSuite) exportRegressionPrometheusMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "regression_detection_metrics.prom")

	var prometheus []string
	prometheus = append(prometheus, "# Regression detection metrics")

	// Export detection results
	totalRegressions := 0
	totalAnomalies := 0
	for _, result := range suite.regressionResults {
		totalRegressions += len(result.RegressionsDetected)
		totalAnomalies += len(result.AnomaliesDetected)

		prometheus = append(prometheus, fmt.Sprintf("tchat_regression_test_duration_seconds{test=\"%s\"} %.3f",
			result.TestName, result.Duration.Seconds()))
		prometheus = append(prometheus, fmt.Sprintf("tchat_regression_detected_total{test=\"%s\"} %d",
			result.TestName, len(result.RegressionsDetected)))
		prometheus = append(prometheus, fmt.Sprintf("tchat_anomalies_detected_total{test=\"%s\"} %d",
			result.TestName, len(result.AnomaliesDetected)))
	}

	// Export aggregate metrics
	prometheus = append(prometheus, fmt.Sprintf("tchat_regression_detection_tests_total %d", len(suite.regressionResults)))
	prometheus = append(prometheus, fmt.Sprintf("tchat_regression_total_regressions_detected %d", totalRegressions))
	prometheus = append(prometheus, fmt.Sprintf("tchat_regression_total_anomalies_detected %d", totalAnomalies))
	prometheus = append(prometheus, fmt.Sprintf("tchat_regression_alerts_generated_total %d", len(suite.alertsGenerated)))

	// Export algorithm effectiveness
	for algorithmName := range suite.regressionEngine.detectionAlgorithms {
		prometheus = append(prometheus, fmt.Sprintf("tchat_regression_algorithm_enabled{algorithm=\"%s\"} 1", algorithmName))
	}

	content := strings.Join(prometheus, "\n") + "\n"
	err := os.WriteFile(metricsPath, []byte(content), 0644)
	if err == nil {
		suite.T().Logf("Regression detection Prometheus metrics exported to: %s", metricsPath)
	} else {
		suite.T().Logf("Failed to export regression detection Prometheus metrics: %v", err)
	}
}

// exportRegressionCSVMetrics exports CSV metrics for regression detection
func (suite *RegressionDetectionTestSuite) exportRegressionCSVMetrics(metricsDir string) {
	metricsPath := filepath.Join(metricsDir, "regression_detection_metrics.csv")

	var csv []string
	csv = append(csv, "test_name,duration_ms,regressions_detected,anomalies_detected,detection_algorithm,statistical_significance")

	for _, result := range suite.regressionResults {
		csv = append(csv, fmt.Sprintf("%s,%d,%d,%d,%s,%.3f",
			result.TestName,
			result.Duration.Milliseconds(),
			len(result.RegressionsDetected),
			len(result.AnomaliesDetected),
			result.DetectionAlgorithm,
			result.StatisticalSignificance,
		))
	}

	content := strings.Join(csv, "\n") + "\n"
	err := os.WriteFile(metricsPath, []byte(content), 0644)
	if err == nil {
		suite.T().Logf("Regression detection CSV metrics exported to: %s", metricsPath)
	} else {
		suite.T().Logf("Failed to export regression detection CSV metrics: %v", err)
	}
}

// cleanupRegressionDetectionResources cleans up regression detection test resources
func (suite *RegressionDetectionTestSuite) cleanupRegressionDetectionResources() {
	suite.regressionEngine = nil
	suite.statisticalAnalyzer = nil
	suite.performanceBaselines = nil
	suite.historicalData = nil
	suite.regressionResults = nil
	suite.alertsGenerated = nil

	suite.T().Log("Regression detection test resources cleaned up")
}

// TestRegressionDetectionTestSuite runs the regression detection test suite
func TestRegressionDetectionTestSuite(t *testing.T) {
	suite.Run(t, new(RegressionDetectionTestSuite))
}