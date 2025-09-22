export type TestType = 'unit' | 'integration' | 'e2e' | 'visual' | 'performance';
export type TestStatus = 'pending' | 'running' | 'passed' | 'failed' | 'skipped';
export type TestEnvironment = 'browser' | 'node' | 'storybook';

export interface PerformanceMetrics {
  renderTime: number;        // Initial render time in ms
  rerenderTime: number;      // Re-render time in ms
  bundleSize: number;        // Component bundle size in bytes
  memoryUsage: number;       // Memory usage in bytes
  interactionLatency: number; // User interaction response time
  threshold: PerformanceThreshold;
  history: PerformanceHistory[];
}

export interface PerformanceThreshold {
  renderTime: number;        // Maximum acceptable render time
  bundleSize: number;        // Maximum acceptable bundle size
  memoryUsage: number;       // Maximum acceptable memory usage
  interactionLatency: number; // Maximum acceptable interaction latency
}

export interface PerformanceHistory {
  timestamp: Date;
  metrics: Omit<PerformanceMetrics, 'threshold' | 'history'>;
}

export interface CoverageMetrics {
  lines: CoverageData;         // Line coverage
  branches: CoverageData;      // Branch coverage
  functions: CoverageData;     // Function coverage
  statements: CoverageData;    // Statement coverage
  overall: number;             // Overall coverage percentage
  threshold: number;           // Minimum required coverage
  trend: CoverageTrend[];      // Coverage history
}

export interface CoverageData {
  total: number;               // Total items
  covered: number;             // Covered items
  percentage: number;          // Coverage percentage
  uncovered: string[];         // List of uncovered items
}

export interface CoverageTrend {
  timestamp: Date;
  percentage: number;
  delta: number;              // Change from previous
}

export interface TestExecution {
  id: string;
  status: TestStatus;
  executionTime: number;
  timestamp: Date;
  error?: TestError;
  retryCount: number;
}

export interface TestError {
  type: string;              // Error type
  message: string;           // Error message
  stack: string;            // Stack trace
  line?: number;            // Line number where error occurred
  suggestions: string[];     // Suggested fixes
}

export interface TestCase {
  id: string;                   // Unique test case identifier
  name: string;                // Test case description
  category: TestCategory;      // rendering | props | interactions | accessibility | performance
  priority: TestPriority;      // critical | high | medium | low
  input: TestInput;            // Test input data
  expectedOutput: TestOutput;  // Expected results
  assertions: TestAssertion[]; // Test assertions
  setup?: TestSetup;           // Pre-test setup requirements
  teardown?: TestTeardown;     // Post-test cleanup
  status: TestCaseStatus;      // pending | passed | failed | skipped | flaky
  executionHistory: TestExecution[];
}

export type TestCategory = 'rendering' | 'props' | 'interactions' | 'accessibility' | 'performance';
export type TestPriority = 'critical' | 'high' | 'medium' | 'low';
export type TestCaseStatus = 'pending' | 'passed' | 'failed' | 'skipped' | 'flaky';

export interface TestInput {
  props?: Record<string, any>;
  state?: Record<string, any>;
  context?: Record<string, any>;
  userEvents?: UserEvent[];
}

export interface UserEvent {
  type: 'click' | 'type' | 'focus' | 'hover' | 'keypress' | 'drag';
  target: string;
  value?: any;
  delay?: number;
}

export interface TestOutput {
  rendered?: boolean;
  snapshot?: string;
  dom?: string;
  state?: Record<string, any>;
  callbacks?: CallbackAssertion[];
  metrics?: PerformanceMetrics;
}

export interface CallbackAssertion {
  name: string;
  calledTimes: number;
  arguments?: any[];
}

export interface TestAssertion {
  type: 'toBe' | 'toEqual' | 'toContain' | 'toHaveBeenCalled' | 'toBeVisible' | 'custom';
  target: string;
  expected: any;
  actual?: any;
  passed?: boolean;
  message?: string;
}

export interface TestSetup {
  mocks?: MockDefinition[];
  fixtures?: Record<string, any>;
  environment?: Record<string, any>;
}

export interface MockDefinition {
  type: 'function' | 'module' | 'api';
  target: string;
  implementation: any;
}

export interface TestTeardown {
  cleanup?: string[];
  restores?: string[];
}

export interface TestSuite {
  id: string;                   // Unique test suite identifier
  componentId: string;          // Reference to Component
  type: TestType;              // unit | integration | e2e | visual | performance
  filePath: string;            // Path to test file
  testCases: TestCase[];       // Individual test cases
  coverage: CoverageMetrics;   // Coverage information
  performance: PerformanceMetrics;
  status: TestStatus;          // pending | running | passed | failed | skipped
  lastRun?: Date;              // Last execution timestamp
  executionTime?: number;      // Last execution time in ms
  environment: TestEnvironment; // browser | node | storybook
}

// Factory functions for creating test entities
export const createTestSuite = (overrides?: Partial<TestSuite>): TestSuite => ({
  id: '',
  componentId: '',
  type: 'unit',
  filePath: '',
  testCases: [],
  coverage: {
    lines: { total: 0, covered: 0, percentage: 0, uncovered: [] },
    branches: { total: 0, covered: 0, percentage: 0, uncovered: [] },
    functions: { total: 0, covered: 0, percentage: 0, uncovered: [] },
    statements: { total: 0, covered: 0, percentage: 0, uncovered: [] },
    overall: 0,
    threshold: 90,
    trend: [],
  },
  performance: {
    renderTime: 0,
    rerenderTime: 0,
    bundleSize: 0,
    memoryUsage: 0,
    interactionLatency: 0,
    threshold: {
      renderTime: 100,
      bundleSize: 50000,
      memoryUsage: 10485760,
      interactionLatency: 50,
    },
    history: [],
  },
  status: 'pending',
  environment: 'node',
  ...overrides,
});

export const createTestCase = (overrides?: Partial<TestCase>): TestCase => ({
  id: '',
  name: '',
  category: 'rendering',
  priority: 'medium',
  input: {},
  expectedOutput: {},
  assertions: [],
  status: 'pending',
  executionHistory: [],
  ...overrides,
});