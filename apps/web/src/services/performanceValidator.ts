/**
 * Performance & Accessibility Validator (T054-T056)
 * Component performance optimization (<200ms load, 60fps animations)
 * Accessibility compliance validation (WCAG 2.1 AA across platforms)
 * Cross-platform visual consistency validation (97% target)
 */

import { api } from './api';
import type { ComponentDefinition } from './componentSync';

// Performance validation interfaces
export interface PerformanceValidationConfig {
  thresholds: {
    loadTime: number; // 200ms constitutional requirement
    renderTime: number; // 16ms for 60fps
    bundleSize: number; // 500KB target
    memoryUsage: number; // 100MB mobile, 500MB desktop
    animationFrameRate: number; // 60fps
    timeToInteractive: number; // 3s on 3G
    firstContentfulPaint: number; // 1.5s
    largestContentfulPaint: number; // 2.5s
    cumulativeLayoutShift: number; // 0.1
    firstInputDelay: number; // 100ms
  };
  platforms: Platform[];
  networkConditions: NetworkCondition[];
  devices: DeviceProfile[];
}

export interface NetworkCondition {
  name: string;
  downloadSpeed: number; // Mbps
  uploadSpeed: number; // Mbps
  latency: number; // ms
  packetLoss: number; // percentage
}

export interface DeviceProfile {
  name: string;
  category: 'mobile' | 'tablet' | 'desktop';
  cpu: 'low' | 'mid' | 'high';
  memory: number; // GB
  screenSize: { width: number; height: number };
  pixelRatio: number;
}

export interface PerformanceMetrics {
  componentId: string;
  platform: Platform;
  device: string;
  network: string;
  timestamp: string;
  loadTime: number; // ms
  renderTime: number; // ms
  bundleSize: number; // bytes
  memoryUsage: number; // MB
  animationFrameRate: number; // fps
  coreWebVitals: {
    firstContentfulPaint: number;
    largestContentfulPaint: number;
    cumulativeLayoutShift: number;
    firstInputDelay: number;
    totalBlockingTime: number;
    timeToInteractive: number;
  };
  resourceMetrics: {
    jsSize: number;
    cssSize: number;
    imageSize: number;
    totalRequests: number;
    cachedResources: number;
    renderBlockingResources: number;
  };
  runtimeMetrics: {
    heapUsedSize: number;
    heapTotalSize: number;
    scriptDuration: number;
    layoutDuration: number;
    paintDuration: number;
  };
}

export interface PerformanceValidationResult {
  componentId: string;
  overallScore: number; // 0.0 - 1.0
  meetsConstitutionalRequirements: boolean; // <200ms load time
  platformResults: Record<Platform, PlatformPerformanceResult>;
  issues: PerformanceIssue[];
  recommendations: PerformanceRecommendation[];
  validationMetadata: {
    validatedAt: string;
    validationDuration: number;
    metricsCollected: number;
    benchmarkVersion: string;
  };
}

export interface PlatformPerformanceResult {
  platform: Platform;
  overallScore: number;
  meetsThresholds: boolean;
  deviceResults: Record<string, DevicePerformanceResult>;
  averageMetrics: PerformanceMetrics;
  bestMetrics: PerformanceMetrics;
  worstMetrics: PerformanceMetrics;
}

export interface DevicePerformanceResult {
  device: string;
  networkResults: Record<string, PerformanceMetrics>;
  averageScore: number;
  meetsThresholds: boolean;
}

export interface PerformanceIssue {
  id: string;
  severity: 'low' | 'medium' | 'high' | 'critical' | 'constitutional_violation';
  type: 'load_time' | 'render_performance' | 'bundle_size' | 'memory_usage' | 'animation_performance' | 'core_web_vitals';
  title: string;
  description: string;
  affectedPlatforms: Platform[];
  affectedDevices: string[];
  threshold: number;
  actualValue: number;
  impact: 'user_experience' | 'constitutional_compliance' | 'performance_budget' | 'seo_ranking';
  estimatedFixTime: number; // hours
}

export interface PerformanceRecommendation {
  id: string;
  priority: 'low' | 'medium' | 'high' | 'critical';
  category: 'code_optimization' | 'asset_optimization' | 'caching' | 'lazy_loading' | 'bundling' | 'runtime_optimization';
  title: string;
  description: string;
  implementation: string;
  estimatedImpact: {
    loadTimeImprovement: number; // ms
    bundleSizeReduction: number; // bytes
    memoryReduction: number; // MB
    performanceScoreIncrease: number; // 0.0 - 1.0
  };
  platforms: Platform[];
  estimatedEffort: number; // hours
  dependencies: string[];
}

// Accessibility validation interfaces
export interface AccessibilityValidationConfig {
  wcagLevel: 'A' | 'AA' | 'AAA';
  guidelines: AccessibilityGuideline[];
  platforms: Platform[];
  testingTools: AccessibilityTestingTool[];
  automatedChecks: boolean;
  manualTestingRequired: boolean;
}

export interface AccessibilityGuideline {
  id: string;
  wcagReference: string;
  level: 'A' | 'AA' | 'AAA';
  principle: 'perceivable' | 'operable' | 'understandable' | 'robust';
  title: string;
  description: string;
  testCriteria: AccessibilityTestCriterion[];
}

export interface AccessibilityTestCriterion {
  id: string;
  description: string;
  automated: boolean;
  testMethod: string;
  expectedOutcome: string;
  priority: 'required' | 'recommended' | 'optional';
}

export interface AccessibilityTestingTool {
  name: string;
  type: 'automated' | 'manual' | 'screen_reader' | 'keyboard_navigation';
  platform: Platform;
  coverage: string[];
}

export interface AccessibilityValidationResult {
  componentId: string;
  overallScore: number; // 0.0 - 1.0
  wcagLevel: 'A' | 'AA' | 'AAA' | 'non_compliant';
  meetsConstitutionalRequirements: boolean; // WCAG 2.1 AA compliance
  platformResults: Record<Platform, PlatformAccessibilityResult>;
  issues: AccessibilityIssue[];
  recommendations: AccessibilityRecommendation[];
  testResults: AccessibilityTestResult[];
  validationMetadata: {
    validatedAt: string;
    toolsUsed: string[];
    manualTestingCompleted: boolean;
    validationDuration: number;
  };
}

export interface PlatformAccessibilityResult {
  platform: Platform;
  overallScore: number;
  wcagLevel: 'A' | 'AA' | 'AAA' | 'non_compliant';
  meetsRequirements: boolean;
  guidelineResults: Record<string, GuidelineResult>;
  testCoverage: number; // 0.0 - 1.0
  automatedTestsPass: number;
  manualTestsPass: number;
  totalTests: number;
}

export interface GuidelineResult {
  guideline: string;
  level: 'A' | 'AA' | 'AAA';
  compliance: 'pass' | 'fail' | 'partial' | 'not_tested';
  score: number; // 0.0 - 1.0
  testResults: AccessibilityTestResult[];
  issues: AccessibilityIssue[];
}

export interface AccessibilityTestResult {
  testId: string;
  guideline: string;
  criterion: string;
  result: 'pass' | 'fail' | 'warning' | 'manual_review_required';
  details: string;
  toolUsed: string;
  platform: Platform;
  timestamp: string;
  evidence?: {
    screenshots?: string[];
    codeSnippets?: string[];
    userFlows?: string[];
  };
}

export interface AccessibilityIssue {
  id: string;
  severity: 'low' | 'medium' | 'high' | 'critical' | 'constitutional_violation';
  type: 'color_contrast' | 'keyboard_navigation' | 'screen_reader' | 'focus_management' | 'semantic_markup' | 'alternative_text' | 'form_labels' | 'heading_structure';
  wcagReference: string;
  title: string;
  description: string;
  affectedPlatforms: Platform[];
  impact: 'barrier' | 'difficulty' | 'inconvenience' | 'constitutional_violation';
  userGroups: string[];
  estimatedUsers: number;
  estimatedFixTime: number; // hours
}

export interface AccessibilityRecommendation {
  id: string;
  priority: 'low' | 'medium' | 'high' | 'critical';
  category: 'color_contrast' | 'keyboard_accessibility' | 'screen_reader' | 'semantic_markup' | 'focus_management' | 'responsive_design';
  title: string;
  description: string;
  implementation: string;
  wcagReference: string;
  platforms: Platform[];
  estimatedImpact: {
    usersAffected: number;
    complianceImprovement: number; // 0.0 - 1.0
    accessibilityScoreIncrease: number; // 0.0 - 1.0
  };
  estimatedEffort: number; // hours
  dependencies: string[];
  testingRequired: string[];
}

// Visual consistency validation interfaces
export interface VisualConsistencyConfig {
  constitutionalThreshold: number; // 0.97
  platforms: Platform[];
  comparisonMethods: ComparisonMethod[];
  visualDifferenceThreshold: number; // pixel difference threshold
  colorAccuracyTolerance: number; // OKLCH tolerance
  spatialConsistencyTolerance: number; // spacing/sizing tolerance
}

export interface ComparisonMethod {
  name: string;
  type: 'pixel_comparison' | 'structural_comparison' | 'perceptual_comparison' | 'semantic_comparison';
  weight: number; // 0.0 - 1.0
  configuration: Record<string, any>;
}

export interface VisualConsistencyResult {
  componentId: string;
  overallConsistencyScore: number; // 0.0 - 1.0
  meetsConstitutionalRequirement: boolean; // >= 97%
  platformComparisons: PlatformComparison[];
  issues: VisualConsistencyIssue[];
  recommendations: VisualConsistencyRecommendation[];
  validationMetadata: {
    validatedAt: string;
    methodsUsed: string[];
    screenshotsTaken: number;
    comparisonsDone: number;
  };
}

export interface PlatformComparison {
  platformA: Platform;
  platformB: Platform;
  consistencyScore: number; // 0.0 - 1.0
  differences: VisualDifference[];
  screenshots: {
    platformA: string;
    platformB: string;
    diffImage: string;
  };
}

export interface VisualDifference {
  type: 'color' | 'spacing' | 'typography' | 'layout' | 'interaction_state';
  severity: 'minor' | 'moderate' | 'major' | 'constitutional_violation';
  description: string;
  location: { x: number; y: number; width: number; height: number };
  expectedValue: any;
  actualValue: any;
  pixelDifference: number;
  impact: string;
}

export interface VisualConsistencyIssue {
  id: string;
  severity: 'low' | 'medium' | 'high' | 'constitutional_violation';
  type: 'color_mismatch' | 'spacing_inconsistency' | 'typography_difference' | 'layout_variation' | 'interaction_inconsistency';
  title: string;
  description: string;
  affectedPlatforms: Platform[];
  consistencyScore: number;
  visualEvidence: {
    screenshots: Record<Platform, string>;
    diffImages: string[];
    annotations: string[];
  };
  estimatedFixTime: number; // hours
}

export interface VisualConsistencyRecommendation {
  id: string;
  priority: 'low' | 'medium' | 'high' | 'critical';
  category: 'design_tokens' | 'component_implementation' | 'platform_conventions' | 'responsive_design';
  title: string;
  description: string;
  implementation: string;
  platforms: Platform[];
  estimatedImpact: {
    consistencyImprovement: number; // 0.0 - 1.0
    affectedComponents: string[];
  };
  estimatedEffort: number; // hours
}

export type Platform = 'web' | 'ios' | 'android';

// Default configurations
export const DEFAULT_PERFORMANCE_CONFIG: PerformanceValidationConfig = {
  thresholds: {
    loadTime: 200, // Constitutional requirement
    renderTime: 16, // 60fps
    bundleSize: 500 * 1024, // 500KB
    memoryUsage: 100, // 100MB mobile baseline
    animationFrameRate: 60,
    timeToInteractive: 3000, // 3s on 3G
    firstContentfulPaint: 1500, // 1.5s
    largestContentfulPaint: 2500, // 2.5s
    cumulativeLayoutShift: 0.1,
    firstInputDelay: 100
  },
  platforms: ['web', 'ios', 'android'],
  networkConditions: [
    { name: '3G Fast', downloadSpeed: 1.6, uploadSpeed: 0.75, latency: 150, packetLoss: 0 },
    { name: '4G', downloadSpeed: 9, uploadSpeed: 9, latency: 50, packetLoss: 0 },
    { name: 'WiFi', downloadSpeed: 30, uploadSpeed: 15, latency: 10, packetLoss: 0 }
  ],
  devices: [
    { name: 'iPhone 12', category: 'mobile', cpu: 'high', memory: 4, screenSize: { width: 390, height: 844 }, pixelRatio: 3 },
    { name: 'Pixel 5', category: 'mobile', cpu: 'mid', memory: 8, screenSize: { width: 393, height: 851 }, pixelRatio: 2.75 },
    { name: 'iPad Air', category: 'tablet', cpu: 'high', memory: 4, screenSize: { width: 820, height: 1180 }, pixelRatio: 2 },
    { name: 'Desktop 1080p', category: 'desktop', cpu: 'high', memory: 16, screenSize: { width: 1920, height: 1080 }, pixelRatio: 1 }
  ]
};

export const DEFAULT_ACCESSIBILITY_CONFIG: AccessibilityValidationConfig = {
  wcagLevel: 'AA',
  guidelines: [], // Would be populated with WCAG guidelines
  platforms: ['web', 'ios', 'android'],
  testingTools: [], // Would be populated with testing tools
  automatedChecks: true,
  manualTestingRequired: true
};

export const DEFAULT_VISUAL_CONSISTENCY_CONFIG: VisualConsistencyConfig = {
  constitutionalThreshold: 0.97,
  platforms: ['web', 'ios', 'android'],
  comparisonMethods: [
    { name: 'pixel_diff', type: 'pixel_comparison', weight: 0.3, configuration: {} },
    { name: 'structural_diff', type: 'structural_comparison', weight: 0.3, configuration: {} },
    { name: 'perceptual_diff', type: 'perceptual_comparison', weight: 0.4, configuration: {} }
  ],
  visualDifferenceThreshold: 0.05, // 5% pixel difference threshold
  colorAccuracyTolerance: 0.01, // 1% OKLCH tolerance
  spatialConsistencyTolerance: 0.02 // 2% spacing tolerance
};

// RTK Query API endpoints
export const performanceValidationApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // Performance validation endpoints
    validateComponentPerformance: builder.mutation<
      PerformanceValidationResult,
      { componentId: string; config?: Partial<PerformanceValidationConfig> }
    >({
      query: ({ componentId, config }) => ({
        url: `/validation/performance/${componentId}`,
        method: 'POST',
        body: { config: config || DEFAULT_PERFORMANCE_CONFIG }
      }),
      transformResponse: (response: PerformanceValidationResult) => ({
        ...response,
        meetsConstitutionalRequirements: response.overallScore >= 0.9 &&
          Object.values(response.platformResults).every(platform =>
            platform.averageMetrics.loadTime <= DEFAULT_PERFORMANCE_CONFIG.thresholds.loadTime
          )
      })
    }),

    // Get performance metrics for component
    getPerformanceMetrics: builder.query<
      PerformanceMetrics[],
      { componentId: string; platform?: Platform; timeRange?: string }
    >({
      query: ({ componentId, platform, timeRange = '7d' }) => ({
        url: `/validation/performance/${componentId}/metrics`,
        params: { platform, timeRange }
      }),
      providesTags: (_result, _error, { componentId }) => [
        { type: 'PerformanceMetrics', id: componentId }
      ]
    }),

    // Accessibility validation endpoints
    validateComponentAccessibility: builder.mutation<
      AccessibilityValidationResult,
      { componentId: string; config?: Partial<AccessibilityValidationConfig> }
    >({
      query: ({ componentId, config }) => ({
        url: `/validation/accessibility/${componentId}`,
        method: 'POST',
        body: { config: config || DEFAULT_ACCESSIBILITY_CONFIG }
      }),
      transformResponse: (response: AccessibilityValidationResult) => ({
        ...response,
        meetsConstitutionalRequirements: response.wcagLevel === 'AA' || response.wcagLevel === 'AAA'
      })
    }),

    // Get accessibility test results
    getAccessibilityTestResults: builder.query<
      AccessibilityTestResult[],
      { componentId: string; platform?: Platform; guideline?: string }
    >({
      query: ({ componentId, platform, guideline }) => ({
        url: `/validation/accessibility/${componentId}/tests`,
        params: { platform, guideline }
      }),
      providesTags: (_result, _error, { componentId }) => [
        { type: 'AccessibilityTests', id: componentId }
      ]
    }),

    // Visual consistency validation endpoints
    validateVisualConsistency: builder.mutation<
      VisualConsistencyResult,
      { componentId: string; config?: Partial<VisualConsistencyConfig> }
    >({
      query: ({ componentId, config }) => ({
        url: `/validation/visual-consistency/${componentId}`,
        method: 'POST',
        body: { config: config || DEFAULT_VISUAL_CONSISTENCY_CONFIG }
      }),
      transformResponse: (response: VisualConsistencyResult) => ({
        ...response,
        meetsConstitutionalRequirement: response.overallConsistencyScore >= 0.97
      })
    }),

    // Get visual comparison screenshots
    getVisualComparisons: builder.query<
      { screenshots: Record<Platform, string>; diffImages: string[] },
      { componentId: string; variant?: string; state?: string }
    >({
      query: ({ componentId, variant, state }) => ({
        url: `/validation/visual-consistency/${componentId}/screenshots`,
        params: { variant, state }
      }),
      providesTags: (_result, _error, { componentId }) => [
        { type: 'VisualComparisons', id: componentId }
      ]
    }),

    // Comprehensive validation (all three areas)
    validateComponentComprehensive: builder.mutation<
      {
        performance: PerformanceValidationResult;
        accessibility: AccessibilityValidationResult;
        visualConsistency: VisualConsistencyResult;
        overallScore: number;
        constitutionalCompliance: boolean;
        summary: {
          totalIssues: number;
          criticalIssues: number;
          constitutionalViolations: number;
          estimatedFixTime: number;
        };
      },
      {
        componentId: string;
        performanceConfig?: Partial<PerformanceValidationConfig>;
        accessibilityConfig?: Partial<AccessibilityValidationConfig>;
        visualConsistencyConfig?: Partial<VisualConsistencyConfig>;
      }
    >({
      query: ({ componentId, performanceConfig, accessibilityConfig, visualConsistencyConfig }) => ({
        url: `/validation/comprehensive/${componentId}`,
        method: 'POST',
        body: {
          performanceConfig: performanceConfig || DEFAULT_PERFORMANCE_CONFIG,
          accessibilityConfig: accessibilityConfig || DEFAULT_ACCESSIBILITY_CONFIG,
          visualConsistencyConfig: visualConsistencyConfig || DEFAULT_VISUAL_CONSISTENCY_CONFIG
        }
      }),
      invalidatesTags: (_result, _error, { componentId }) => [
        { type: 'PerformanceMetrics', id: componentId },
        { type: 'AccessibilityTests', id: componentId },
        { type: 'VisualComparisons', id: componentId }
      ]
    }),

    // Batch validate multiple components
    validateComponentsBatch: builder.mutation<
      Array<{
        componentId: string;
        performance: PerformanceValidationResult;
        accessibility: AccessibilityValidationResult;
        visualConsistency: VisualConsistencyResult;
      }>,
      { componentIds: string[]; parallel: boolean }
    >({
      query: ({ componentIds, parallel }) => ({
        url: '/validation/batch',
        method: 'POST',
        body: { componentIds, parallel }
      })
    }),

    // Generate validation report
    generateValidationReport: builder.mutation<
      {
        reportId: string;
        url: string;
        format: 'json' | 'pdf' | 'html';
        generatedAt: string;
      },
      {
        componentIds: string[];
        validationTypes: ('performance' | 'accessibility' | 'visual_consistency')[];
        format: 'json' | 'pdf' | 'html';
        includeScreenshots: boolean;
      }
    >({
      query: (params) => ({
        url: '/validation/reports/generate',
        method: 'POST',
        body: params
      })
    })
  })
});

// Export hooks
export const {
  useValidateComponentPerformanceMutation,
  useGetPerformanceMetricsQuery,
  useValidateComponentAccessibilityMutation,
  useGetAccessibilityTestResultsQuery,
  useValidateVisualConsistencyMutation,
  useGetVisualComparisonsQuery,
  useValidateComponentComprehensiveMutation,
  useValidateComponentsBatchMutation,
  useGenerateValidationReportMutation
} = performanceValidationApi;

// Validation Service Class
export class ValidationService {
  /**
   * Validate constitutional performance requirements
   */
  async validateConstitutionalPerformance(
    componentId: string,
    metrics: PerformanceMetrics[]
  ): Promise<{
    compliant: boolean;
    violations: PerformanceIssue[];
    overallScore: number;
    platformBreakdown: Record<Platform, { compliant: boolean; score: number }>;
  }> {
    const violations: PerformanceIssue[] = [];
    const platformBreakdown: Record<Platform, { compliant: boolean; score: number }> = {} as any;
    let totalScore = 0;
    let platformCount = 0;

    const byPlatform = metrics.reduce((acc, metric) => {
      if (!acc[metric.platform]) acc[metric.platform] = [];
      acc[metric.platform].push(metric);
      return acc;
    }, {} as Record<Platform, PerformanceMetrics[]>);

    for (const [platform, platformMetrics] of Object.entries(byPlatform)) {
      const avgLoadTime = platformMetrics.reduce((sum, m) => sum + m.loadTime, 0) / platformMetrics.length;
      const avgFrameRate = platformMetrics.reduce((sum, m) => sum + m.animationFrameRate, 0) / platformMetrics.length;

      let platformScore = 1.0;
      let isCompliant = true;

      // Check load time constitutional requirement (200ms)
      if (avgLoadTime > DEFAULT_PERFORMANCE_CONFIG.thresholds.loadTime) {
        violations.push({
          id: `load-time-${platform}`,
          severity: 'constitutional_violation',
          type: 'load_time',
          title: `Constitutional load time violation on ${platform}`,
          description: `Average load time of ${avgLoadTime.toFixed(1)}ms exceeds 200ms constitutional requirement`,
          affectedPlatforms: [platform as Platform],
          affectedDevices: platformMetrics.map(m => m.device),
          threshold: DEFAULT_PERFORMANCE_CONFIG.thresholds.loadTime,
          actualValue: avgLoadTime,
          impact: 'constitutional_compliance',
          estimatedFixTime: 8
        });
        platformScore -= 0.3;
        isCompliant = false;
      }

      // Check animation frame rate (60fps)
      if (avgFrameRate < DEFAULT_PERFORMANCE_CONFIG.thresholds.animationFrameRate) {
        violations.push({
          id: `frame-rate-${platform}`,
          severity: 'high',
          type: 'animation_performance',
          title: `Animation performance issue on ${platform}`,
          description: `Average frame rate of ${avgFrameRate.toFixed(1)}fps is below 60fps requirement`,
          affectedPlatforms: [platform as Platform],
          affectedDevices: platformMetrics.map(m => m.device),
          threshold: DEFAULT_PERFORMANCE_CONFIG.thresholds.animationFrameRate,
          actualValue: avgFrameRate,
          impact: 'user_experience',
          estimatedFixTime: 4
        });
        platformScore -= 0.2;
      }

      platformBreakdown[platform as Platform] = { compliant: isCompliant, score: Math.max(platformScore, 0) };
      totalScore += platformScore;
      platformCount++;
    }

    const overallScore = platformCount > 0 ? totalScore / platformCount : 0;
    const overallCompliant = overallScore >= 0.9 && violations.filter(v => v.severity === 'constitutional_violation').length === 0;

    return {
      compliant: overallCompliant,
      violations,
      overallScore,
      platformBreakdown
    };
  }

  /**
   * Validate WCAG 2.1 AA accessibility requirements
   */
  async validateAccessibilityCompliance(
    componentId: string,
    testResults: AccessibilityTestResult[]
  ): Promise<{
    compliant: boolean;
    wcagLevel: 'A' | 'AA' | 'AAA' | 'non_compliant';
    violations: AccessibilityIssue[];
    overallScore: number;
    platformBreakdown: Record<Platform, { compliant: boolean; score: number; level: string }>;
  }> {
    const violations: AccessibilityIssue[] = [];
    const platformBreakdown: Record<Platform, { compliant: boolean; score: number; level: string }> = {} as any;

    const byPlatform = testResults.reduce((acc, result) => {
      if (!acc[result.platform]) acc[result.platform] = [];
      acc[result.platform].push(result);
      return acc;
    }, {} as Record<Platform, AccessibilityTestResult[]>);

    let totalScore = 0;
    let platformCount = 0;

    for (const [platform, results] of Object.entries(byPlatform)) {
      const totalTests = results.length;
      const passedTests = results.filter(r => r.result === 'pass').length;
      const failedTests = results.filter(r => r.result === 'fail').length;

      const platformScore = totalTests > 0 ? passedTests / totalTests : 0;
      const wcagLevel = this.calculateWCAGLevel(results);
      const isCompliant = wcagLevel === 'AA' || wcagLevel === 'AAA';

      if (!isCompliant) {
        violations.push({
          id: `wcag-${platform}`,
          severity: 'constitutional_violation',
          type: 'screen_reader',
          wcagReference: 'WCAG 2.1',
          title: `WCAG 2.1 AA compliance violation on ${platform}`,
          description: `Platform does not meet constitutional WCAG 2.1 AA requirements (${failedTests} failed tests)`,
          affectedPlatforms: [platform as Platform],
          impact: 'constitutional_violation',
          userGroups: ['screen_reader_users', 'keyboard_users', 'low_vision_users'],
          estimatedUsers: 1000000, // Estimated users affected
          estimatedFixTime: failedTests * 2 // 2 hours per failed test
        });
      }

      platformBreakdown[platform as Platform] = {
        compliant: isCompliant,
        score: platformScore,
        level: wcagLevel
      };

      totalScore += platformScore;
      platformCount++;
    }

    const overallScore = platformCount > 0 ? totalScore / platformCount : 0;
    const overallWCAGLevel = this.calculateOverallWCAGLevel(Object.values(platformBreakdown).map(p => p.level));
    const overallCompliant = overallWCAGLevel === 'AA' || overallWCAGLevel === 'AAA';

    return {
      compliant: overallCompliant,
      wcagLevel: overallWCAGLevel,
      violations,
      overallScore,
      platformBreakdown
    };
  }

  /**
   * Validate constitutional visual consistency (97% requirement)
   */
  async validateVisualConsistency(
    componentId: string,
    comparisons: PlatformComparison[]
  ): Promise<{
    compliant: boolean;
    consistencyScore: number;
    violations: VisualConsistencyIssue[];
    platformScores: Record<string, number>;
  }> {
    const violations: VisualConsistencyIssue[] = [];
    const platformScores: Record<string, number> = {};

    let totalConsistency = 0;
    let comparisonCount = 0;

    for (const comparison of comparisons) {
      const { platformA, platformB, consistencyScore, differences } = comparison;

      totalConsistency += consistencyScore;
      comparisonCount++;

      // Track individual platform scores
      if (!platformScores[platformA]) platformScores[platformA] = 0;
      if (!platformScores[platformB]) platformScores[platformB] = 0;
      platformScores[platformA] += consistencyScore;
      platformScores[platformB] += consistencyScore;

      // Check for constitutional violations (below 97%)
      if (consistencyScore < 0.97) {
        const constitutionalDiffs = differences.filter(d => d.severity === 'constitutional_violation');

        violations.push({
          id: `consistency-${platformA}-${platformB}`,
          severity: 'constitutional_violation',
          type: 'color_mismatch', // Would be determined by difference types
          title: `Constitutional visual consistency violation between ${platformA} and ${platformB}`,
          description: `Consistency score of ${(consistencyScore * 100).toFixed(1)}% is below 97% constitutional requirement`,
          affectedPlatforms: [platformA, platformB],
          consistencyScore,
          visualEvidence: {
            screenshots: comparison.screenshots ? {
              [platformA]: comparison.screenshots.platformA,
              [platformB]: comparison.screenshots.platformB
            } as Record<Platform, string> : {} as Record<Platform, string>,
            diffImages: comparison.screenshots ? [comparison.screenshots.diffImage] : [],
            annotations: constitutionalDiffs.map(d => d.description)
          },
          estimatedFixTime: constitutionalDiffs.length * 3 // 3 hours per major difference
        });
      }
    }

    const overallConsistency = comparisonCount > 0 ? totalConsistency / comparisonCount : 1.0;
    const overallCompliant = overallConsistency >= 0.97;

    // Normalize platform scores
    const platformCount = Object.keys(platformScores).length;
    Object.keys(platformScores).forEach(platform => {
      platformScores[platform] /= Math.max(comparisonCount / platformCount, 1);
    });

    return {
      compliant: overallCompliant,
      consistencyScore: overallConsistency,
      violations,
      platformScores
    };
  }

  // Private helper methods
  private calculateWCAGLevel(results: AccessibilityTestResult[]): 'A' | 'AA' | 'AAA' | 'non_compliant' {
    const testsByLevel = results.reduce((acc, result) => {
      // Extract WCAG level from guideline (mock implementation)
      const level = result.guideline.includes('AA') ? 'AA' : result.guideline.includes('AAA') ? 'AAA' : 'A';
      if (!acc[level]) acc[level] = [];
      acc[level].push(result);
      return acc;
    }, {} as Record<string, AccessibilityTestResult[]>);

    // Check AA compliance first (constitutional requirement)
    const aaTests = testsByLevel['AA'] || [];
    const aaPassed = aaTests.filter(r => r.result === 'pass').length;
    const aaTotal = aaTests.length;

    if (aaTotal > 0 && aaPassed / aaTotal >= 0.95) { // 95% pass rate for AA
      // Check AAA compliance
      const aaaTests = testsByLevel['AAA'] || [];
      const aaaPassed = aaaTests.filter(r => r.result === 'pass').length;
      const aaaTotal = aaaTests.length;

      if (aaaTotal > 0 && aaaPassed / aaaTotal >= 0.95) {
        return 'AAA';
      }
      return 'AA';
    }

    // Check A compliance
    const aTests = testsByLevel['A'] || [];
    const aPassed = aTests.filter(r => r.result === 'pass').length;
    const aTotal = aTests.length;

    if (aTotal > 0 && aPassed / aTotal >= 0.95) {
      return 'A';
    }

    return 'non_compliant';
  }

  private calculateOverallWCAGLevel(platformLevels: string[]): 'A' | 'AA' | 'AAA' | 'non_compliant' {
    if (platformLevels.some(level => level === 'non_compliant')) {
      return 'non_compliant';
    }

    if (platformLevels.every(level => level === 'AAA')) {
      return 'AAA';
    }

    if (platformLevels.every(level => level === 'AA' || level === 'AAA')) {
      return 'AA';
    }

    if (platformLevels.every(level => level === 'A' || level === 'AA' || level === 'AAA')) {
      return 'A';
    }

    return 'non_compliant';
  }
}

// Singleton instance
export const validationService = new ValidationService();

// Utility functions
export const ValidationUtils = {
  formatPerformanceScore: (score: number): string => {
    const percentage = (score * 100).toFixed(1);
    const emoji = score >= 0.9 ? '🚀' : score >= 0.7 ? '⚡' : '🐌';
    return `${emoji} ${percentage}%`;
  },

  formatLoadTime: (ms: number): string => {
    if (ms <= 200) return `✅ ${ms}ms`;
    if (ms <= 500) return `⚠️ ${ms}ms`;
    return `❌ ${ms}ms`;
  },

  formatAccessibilityLevel: (level: string): string => {
    const levelMap = {
      'AAA': '🥇 WCAG AAA',
      'AA': '✅ WCAG AA',
      'A': '⚠️ WCAG A',
      'non_compliant': '❌ Non-compliant'
    };
    return levelMap[level as keyof typeof levelMap] || level;
  },

  formatConsistencyScore: (score: number): string => {
    const percentage = (score * 100).toFixed(1);
    const emoji = score >= 0.97 ? '✅' : score >= 0.90 ? '⚠️' : '❌';
    return `${emoji} ${percentage}%`;
  },

  isConstitutionalViolation: (
    issue: PerformanceIssue | AccessibilityIssue | VisualConsistencyIssue
  ): boolean => {
    return issue.severity === 'constitutional_violation';
  },

  calculateOverallValidationScore: (
    performanceScore: number,
    accessibilityScore: number,
    consistencyScore: number
  ): number => {
    // Weighted average: Performance 30%, Accessibility 35%, Consistency 35%
    return (performanceScore * 0.3) + (accessibilityScore * 0.35) + (consistencyScore * 0.35);
  }
};

export default ValidationService;