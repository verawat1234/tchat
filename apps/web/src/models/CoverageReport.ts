import { CoverageMetrics, CoverageData } from './TestSuite';

export interface CoverageReport {
  id: string;                      // Unique coverage report identifier
  projectId: string;               // Project or component identifier
  timestamp: Date;                 // When coverage was measured
  environment: CoverageEnvironment;
  metrics: CoverageMetrics;        // Overall coverage metrics
  files: FileCoverage[];          // Per-file coverage data
  uncoveredCode: UncoveredCode[]; // Specific uncovered code sections
  thresholds: CoverageThresholds;  // Coverage requirements
  violations: ThresholdViolation[]; // Threshold violations
  comparison?: CoverageComparison; // Comparison with previous report
  summary: CoverageSummary;        // Human-readable summary
}

export interface CoverageEnvironment {
  testFramework: string;           // Vitest, Jest, etc.
  coverageProvider: string;        // v8, istanbul, etc.
  nodeVersion: string;
  testCommand: string;
  configFile?: string;
}

export interface FileCoverage {
  filePath: string;
  relativePath: string;
  metrics: CoverageMetrics;
  uncoveredLines: number[];        // Line numbers not covered
  uncoveredBranches: BranchCoverage[];
  uncoveredFunctions: string[];    // Function names not covered
  source?: SourceCoverage;         // Source code with coverage annotations
  complexity?: ComplexityMetrics;   // Cyclomatic complexity
}

export interface BranchCoverage {
  line: number;
  branch: number;
  taken: boolean;
  condition?: string;              // The actual condition text
}

export interface SourceCoverage {
  lines: SourceLine[];
}

export interface SourceLine {
  number: number;
  content: string;
  covered: boolean | null;         // null for non-executable lines
  hits?: number;                   // Number of times line was executed
  branch?: BranchInfo;
}

export interface BranchInfo {
  taken: number;                   // Times true branch taken
  notTaken: number;               // Times false branch taken
}

export interface ComplexityMetrics {
  cyclomatic: number;              // Cyclomatic complexity
  cognitive: number;               // Cognitive complexity
  halstead: HalsteadMetrics;
}

export interface HalsteadMetrics {
  difficulty: number;
  volume: number;
  effort: number;
  bugs: number;
  time: number;
}

export interface UncoveredCode {
  filePath: string;
  startLine: number;
  endLine: number;
  type: 'function' | 'branch' | 'statement';
  code: string;
  reason?: string;                 // Why it might be uncovered
  priority: 'critical' | 'high' | 'medium' | 'low';
  suggestion?: string;             // How to cover it
}

export interface CoverageThresholds {
  global: ThresholdConfig;
  perFile?: ThresholdConfig;
  customPaths?: CustomThreshold[];
}

export interface ThresholdConfig {
  lines: number;
  branches: number;
  functions: number;
  statements: number;
}

export interface CustomThreshold {
  pattern: string;                 // Glob pattern for files
  thresholds: ThresholdConfig;
  reason?: string;                 // Why custom threshold
}

export interface ThresholdViolation {
  type: 'lines' | 'branches' | 'functions' | 'statements';
  scope: 'global' | 'file';
  path?: string;                   // File path if file-level violation
  actual: number;
  threshold: number;
  difference: number;
  severity: 'error' | 'warning';
}

export interface CoverageComparison {
  previousReportId: string;
  previousTimestamp: Date;
  delta: CoverageDelta;
  newUncoveredCode: UncoveredCode[];
  fixedCode: UncoveredCode[];      // Previously uncovered, now covered
  trend: 'improving' | 'degrading' | 'stable';
}

export interface CoverageDelta {
  lines: number;                   // Change in line coverage %
  branches: number;                // Change in branch coverage %
  functions: number;               // Change in function coverage %
  statements: number;              // Change in statement coverage %
  overall: number;                 // Change in overall coverage %
}

export interface CoverageSummary {
  totalFiles: number;
  coveredFiles: number;            // Files with 100% coverage
  partialFiles: number;            // Files with partial coverage
  uncoveredFiles: number;          // Files with 0% coverage
  criticalGaps: string[];          // Critical areas lacking coverage
  recommendations: string[];        // Actionable recommendations
  riskLevel: 'low' | 'medium' | 'high' | 'critical';
}

// Coverage analysis helpers
export interface CoverageHeatmap {
  files: HeatmapEntry[];
  scale: HeatmapScale;
}

export interface HeatmapEntry {
  path: string;
  coverage: number;
  complexity: number;
  risk: number;                   // Coverage * Complexity factor
  color: string;                  // Heatmap color code
}

export interface HeatmapScale {
  excellent: { min: 90, color: '#00ff00' };
  good: { min: 80, max: 90, color: '#90ee90' };
  acceptable: { min: 70, max: 80, color: '#ffff00' };
  poor: { min: 50, max: 70, color: '#ffa500' };
  critical: { max: 50, color: '#ff0000' };
}

// Factory functions
export const createCoverageReport = (
  projectId: string,
  overrides?: Partial<CoverageReport>
): CoverageReport => ({
  id: `coverage-${Date.now()}`,
  projectId,
  timestamp: new Date(),
  environment: {
    testFramework: 'vitest',
    coverageProvider: 'v8',
    nodeVersion: process.version,
    testCommand: 'npm test',
  },
  metrics: {
    lines: { total: 0, covered: 0, percentage: 0, uncovered: [] },
    branches: { total: 0, covered: 0, percentage: 0, uncovered: [] },
    functions: { total: 0, covered: 0, percentage: 0, uncovered: [] },
    statements: { total: 0, covered: 0, percentage: 0, uncovered: [] },
    overall: 0,
    threshold: 90,
    trend: [],
  },
  files: [],
  uncoveredCode: [],
  thresholds: {
    global: {
      lines: 90,
      branches: 90,
      functions: 90,
      statements: 90,
    },
  },
  violations: [],
  summary: {
    totalFiles: 0,
    coveredFiles: 0,
    partialFiles: 0,
    uncoveredFiles: 0,
    criticalGaps: [],
    recommendations: [],
    riskLevel: 'low',
  },
  ...overrides,
});

// Validation functions
export const validateCoverageThresholds = (
  metrics: CoverageMetrics,
  thresholds: ThresholdConfig
): ThresholdViolation[] => {
  const violations: ThresholdViolation[] = [];

  const checks = [
    { type: 'lines' as const, actual: metrics.lines.percentage },
    { type: 'branches' as const, actual: metrics.branches.percentage },
    { type: 'functions' as const, actual: metrics.functions.percentage },
    { type: 'statements' as const, actual: metrics.statements.percentage },
  ];

  checks.forEach(({ type, actual }) => {
    const threshold = thresholds[type];
    if (actual < threshold) {
      violations.push({
        type,
        scope: 'global',
        actual,
        threshold,
        difference: threshold - actual,
        severity: threshold - actual > 10 ? 'error' : 'warning',
      });
    }
  });

  return violations;
};

// Risk assessment
export const assessCoverageRisk = (
  coverage: number,
  complexity: number
): 'low' | 'medium' | 'high' | 'critical' => {
  const riskScore = (100 - coverage) * (complexity / 10);

  if (riskScore < 10) return 'low';
  if (riskScore < 30) return 'medium';
  if (riskScore < 60) return 'high';
  return 'critical';
};

// Generate recommendations
export const generateCoverageRecommendations = (
  report: CoverageReport
): string[] => {
  const recommendations: string[] = [];

  // Check overall coverage
  if (report.metrics.overall < 80) {
    recommendations.push('Increase overall test coverage to at least 80% for production readiness');
  }

  // Check branch coverage specifically
  if (report.metrics.branches.percentage < report.metrics.lines.percentage - 10) {
    recommendations.push('Focus on branch coverage - many conditional paths are untested');
  }

  // Check for files with zero coverage
  const zeroCoverageFiles = report.files.filter(f => f.metrics.overall === 0);
  if (zeroCoverageFiles.length > 0) {
    recommendations.push(`Add tests for ${zeroCoverageFiles.length} files with no coverage`);
  }

  // Check for complex files with low coverage
  const riskyFiles = report.files.filter(f =>
    f.complexity && f.complexity.cyclomatic > 10 && f.metrics.overall < 70
  );
  if (riskyFiles.length > 0) {
    recommendations.push(`Prioritize testing complex files: ${riskyFiles.map(f => f.relativePath).join(', ')}`);
  }

  return recommendations;
};