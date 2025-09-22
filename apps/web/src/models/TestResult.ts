import { TestEnvironment, TestError, PerformanceMetrics, CoverageMetrics } from './TestSuite';

export type TestResultStatus = 'passed' | 'failed' | 'skipped' | 'error';
export type ArtifactType = 'screenshot' | 'video' | 'log' | 'coverage' | 'performance';

export interface BrowserInfo {
  name: string;
  version: string;
  platform: string;
  userAgent: string;
  viewport?: {
    width: number;
    height: number;
  };
}

export interface TestMetrics {
  performance?: PerformanceMetrics;
  coverage?: CoverageMetrics;
  assertions: AssertionMetrics;
  timing: TimingMetrics;
  memory?: MemoryMetrics;
}

export interface AssertionMetrics {
  total: number;
  passed: number;
  failed: number;
  skipped: number;
  successRate: number;
}

export interface TimingMetrics {
  setup: number;               // Setup time in ms
  execution: number;           // Main execution time in ms
  teardown: number;           // Teardown time in ms
  total: number;              // Total time in ms
  breakdown?: TimingBreakdown[];
}

export interface TimingBreakdown {
  phase: string;
  duration: number;
  percentage: number;
}

export interface MemoryMetrics {
  heapUsedBefore: number;      // Heap used before test in bytes
  heapUsedAfter: number;       // Heap used after test in bytes
  heapDelta: number;           // Memory leaked/freed
  peakUsage: number;           // Peak memory usage during test
}

export interface TestArtifact {
  id: string;                // Unique artifact identifier
  type: ArtifactType;       // screenshot | video | log | coverage | performance
  filePath: string;         // Path to artifact file
  size: number;             // File size in bytes
  description: string;      // Human-readable description
  metadata: Record<string, any>; // Additional metadata
  createdAt: Date;
  url?: string;             // URL if artifact is hosted
  thumbnail?: string;       // Thumbnail URL for images/videos
}

export interface TestResult {
  id: string;                 // Unique result identifier
  testSuiteId: string;       // Reference to TestSuite
  testCaseId: string;        // Reference to TestCase
  status: TestResultStatus;  // passed | failed | skipped | error
  executionTime: number;     // Execution time in milliseconds
  timestamp: Date;           // When test was executed
  environment: TestEnvironment;
  browser?: BrowserInfo;     // Browser information (for E2E tests)
  error?: TestError;         // Error information if failed
  artifacts: TestArtifact[]; // Screenshots, logs, etc.
  metrics: TestMetrics;      // Performance and coverage data
  retryCount?: number;       // Number of retries if test was retried
  annotations?: TestAnnotation[]; // Additional test metadata
  comparison?: TestComparison; // Comparison with previous runs
}

export interface TestAnnotation {
  key: string;
  value: any;
  type: 'info' | 'warning' | 'error';
  message?: string;
}

export interface TestComparison {
  previousResultId: string;
  executionTimeDelta: number;  // Change in execution time
  performanceDelta?: Partial<PerformanceMetrics>;
  coverageDelta?: number;      // Change in coverage percentage
  isRegression: boolean;        // Whether this is a regression
  improvements: string[];       // List of improvements
  degradations: string[];       // List of degradations
}

// Aggregated results for test suites
export interface TestSuiteResult {
  id: string;
  testSuiteId: string;
  results: TestResult[];
  summary: TestSuiteSummary;
  startTime: Date;
  endTime: Date;
  duration: number;
  artifacts: TestArtifact[];
}

export interface TestSuiteSummary {
  total: number;
  passed: number;
  failed: number;
  skipped: number;
  error: number;
  successRate: number;
  averageExecutionTime: number;
  totalExecutionTime: number;
  coverage?: CoverageMetrics;
  failedTests: FailedTestInfo[];
  slowestTests: SlowTestInfo[];
  flakyTests: string[];
}

export interface FailedTestInfo {
  testCaseId: string;
  testCaseName: string;
  error: TestError;
  artifacts: string[];         // Artifact IDs
}

export interface SlowTestInfo {
  testCaseId: string;
  testCaseName: string;
  executionTime: number;
  threshold: number;           // Expected max time
}

// Result analysis helpers
export interface ResultTrend {
  period: 'daily' | 'weekly' | 'monthly';
  data: TrendDataPoint[];
  summary: TrendSummary;
}

export interface TrendDataPoint {
  date: Date;
  successRate: number;
  averageExecutionTime: number;
  testCount: number;
  coverage?: number;
}

export interface TrendSummary {
  improving: boolean;
  successRateTrend: 'up' | 'down' | 'stable';
  performanceTrend: 'faster' | 'slower' | 'stable';
  coverageTrend?: 'up' | 'down' | 'stable';
  recommendations: string[];
}

// Factory functions
export const createTestResult = (
  testSuiteId: string,
  testCaseId: string,
  status: TestResultStatus,
  overrides?: Partial<TestResult>
): TestResult => ({
  id: `result-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
  testSuiteId,
  testCaseId,
  status,
  executionTime: 0,
  timestamp: new Date(),
  environment: 'node',
  artifacts: [],
  metrics: {
    assertions: {
      total: 0,
      passed: 0,
      failed: 0,
      skipped: 0,
      successRate: 0,
    },
    timing: {
      setup: 0,
      execution: 0,
      teardown: 0,
      total: 0,
    },
  },
  ...overrides,
});

// Analysis functions
export const analyzeTestResults = (results: TestResult[]): TestSuiteSummary => {
  const total = results.length;
  const passed = results.filter(r => r.status === 'passed').length;
  const failed = results.filter(r => r.status === 'failed').length;
  const skipped = results.filter(r => r.status === 'skipped').length;
  const error = results.filter(r => r.status === 'error').length;

  const executionTimes = results.map(r => r.executionTime);
  const totalExecutionTime = executionTimes.reduce((sum, time) => sum + time, 0);
  const averageExecutionTime = total > 0 ? totalExecutionTime / total : 0;

  const failedTests = results
    .filter(r => r.status === 'failed' && r.error)
    .map(r => ({
      testCaseId: r.testCaseId,
      testCaseName: r.testCaseId, // Would need to look up actual name
      error: r.error!,
      artifacts: r.artifacts.map(a => a.id),
    }));

  const slowestTests = results
    .sort((a, b) => b.executionTime - a.executionTime)
    .slice(0, 5)
    .map(r => ({
      testCaseId: r.testCaseId,
      testCaseName: r.testCaseId,
      executionTime: r.executionTime,
      threshold: 1000, // Default 1 second threshold
    }));

  return {
    total,
    passed,
    failed,
    skipped,
    error,
    successRate: total > 0 ? (passed / total) * 100 : 0,
    averageExecutionTime,
    totalExecutionTime,
    failedTests,
    slowestTests,
    flakyTests: [],
  };
};

export const isResultRegression = (
  current: TestResult,
  previous: TestResult
): boolean => {
  // Check if test that was passing is now failing
  if (previous.status === 'passed' && current.status === 'failed') {
    return true;
  }

  // Check if performance has significantly degraded (>20% slower)
  if (current.executionTime > previous.executionTime * 1.2) {
    return true;
  }

  // Check if coverage has dropped
  if (
    current.metrics.coverage &&
    previous.metrics.coverage &&
    current.metrics.coverage.overall < previous.metrics.coverage.overall
  ) {
    return true;
  }

  return false;
};