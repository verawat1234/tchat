import {
  TestCategory,
  TestPriority,
  TestCaseStatus,
  TestInput,
  TestOutput,
  TestAssertion,
  TestSetup,
  TestTeardown,
  TestExecution,
} from './TestSuite';

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
  dependencies?: string[];     // Other test case dependencies
  tags?: string[];            // Tags for grouping and filtering
  timeout?: number;           // Custom timeout in ms
  retry?: RetryConfig;        // Retry configuration for flaky tests
}

export interface RetryConfig {
  maxAttempts: number;
  delay: number;              // Delay between retries in ms
  backoff?: 'linear' | 'exponential';
}

// Test case templates for different component types
export interface AtomTestCase extends TestCase {
  atomSpecific: {
    visualStates: string[];   // hover, focus, disabled, active
    propCombinations: PropCombination[];
    a11yChecks: AccessibilityCheck[];
  };
}

export interface MoleculeTestCase extends TestCase {
  moleculeSpecific: {
    childInteractions: ChildInteraction[];
    stateManagement: StateTransition[];
    eventPropagation: EventFlow[];
  };
}

export interface OrganismTestCase extends TestCase {
  organismSpecific: {
    businessRules: BusinessRuleValidation[];
    apiIntegrations: ApiIntegration[];
    userWorkflows: UserWorkflow[];
    performanceTargets: PerformanceTarget[];
  };
}

export interface PropCombination {
  props: Record<string, any>;
  expectedBehavior: string;
  isValid: boolean;
}

export interface AccessibilityCheck {
  rule: string;               // WCAG rule
  level: 'A' | 'AA' | 'AAA';
  element: string;
  criteria: string;
}

export interface ChildInteraction {
  childId: string;
  action: string;
  expectedParentResponse: string;
  expectedSiblingEffects?: string[];
}

export interface StateTransition {
  from: Record<string, any>;
  action: string;
  to: Record<string, any>;
  sideEffects?: string[];
}

export interface EventFlow {
  source: string;
  event: string;
  propagationPath: string[];
  handlers: string[];
  preventDefault?: boolean;
  stopPropagation?: boolean;
}

export interface BusinessRuleValidation {
  rule: string;
  input: any;
  expectedOutput: any;
  errorScenarios?: ErrorScenario[];
}

export interface ErrorScenario {
  condition: string;
  expectedError: string;
  recovery?: string;
}

export interface ApiIntegration {
  endpoint: string;
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  mockResponse: any;
  expectedHandling: string;
  errorHandling?: ErrorHandling;
}

export interface ErrorHandling {
  statusCode: number;
  expectedBehavior: string;
  userFeedback: string;
}

export interface UserWorkflow {
  name: string;
  steps: WorkflowStep[];
  expectedOutcome: string;
  criticalPath: boolean;
}

export interface WorkflowStep {
  action: string;
  target: string;
  input?: any;
  validation: string;
  screenshot?: boolean;
}

export interface PerformanceTarget {
  metric: 'renderTime' | 'interactionLatency' | 'memoryUsage';
  threshold: number;
  condition: string;
}

// Factory functions for creating test cases
export const createAtomTestCase = (
  name: string,
  componentName: string,
  overrides?: Partial<AtomTestCase>
): AtomTestCase => ({
  id: `atom-${componentName}-${Date.now()}`,
  name,
  category: 'rendering',
  priority: 'high',
  input: {},
  expectedOutput: { rendered: true },
  assertions: [],
  status: 'pending',
  executionHistory: [],
  atomSpecific: {
    visualStates: ['default', 'hover', 'focus', 'disabled'],
    propCombinations: [],
    a11yChecks: [],
  },
  ...overrides,
});

export const createMoleculeTestCase = (
  name: string,
  componentName: string,
  overrides?: Partial<MoleculeTestCase>
): MoleculeTestCase => ({
  id: `molecule-${componentName}-${Date.now()}`,
  name,
  category: 'interactions',
  priority: 'medium',
  input: {},
  expectedOutput: {},
  assertions: [],
  status: 'pending',
  executionHistory: [],
  moleculeSpecific: {
    childInteractions: [],
    stateManagement: [],
    eventPropagation: [],
  },
  ...overrides,
});

export const createOrganismTestCase = (
  name: string,
  componentName: string,
  overrides?: Partial<OrganismTestCase>
): OrganismTestCase => ({
  id: `organism-${componentName}-${Date.now()}`,
  name,
  category: 'interactions',
  priority: 'critical',
  input: {},
  expectedOutput: {},
  assertions: [],
  status: 'pending',
  executionHistory: [],
  organismSpecific: {
    businessRules: [],
    apiIntegrations: [],
    userWorkflows: [],
    performanceTargets: [],
  },
  ...overrides,
});

// Validation helpers
export const validateTestCase = (testCase: TestCase): string[] => {
  const errors: string[] = [];

  if (!testCase.id) errors.push('Test case must have an ID');
  if (!testCase.name) errors.push('Test case must have a name');
  if (testCase.assertions.length === 0) errors.push('Test case must have at least one assertion');

  return errors;
};

export const isTestCaseFlaky = (testCase: TestCase): boolean => {
  if (testCase.executionHistory.length < 5) return false;

  const recentExecutions = testCase.executionHistory.slice(-5);
  const failureCount = recentExecutions.filter(e => e.status === 'failed').length;

  return failureCount >= 2 && failureCount < 5;
};