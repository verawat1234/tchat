/**
 * Accessibility API Service (T040)
 *
 * RTK Query service for accessibility auditing and WCAG 2.1 AA compliance validation.
 * Provides comprehensive accessibility testing, automated scanning, and compliance
 * reporting across Web, iOS, and Android platforms with detailed remediation guidance.
 *
 * Features:
 * - WCAG 2.1 AA/AAA compliance auditing
 * - Cross-platform accessibility validation
 * - Automated and manual testing workflows
 * - Real-time accessibility monitoring
 * - Detailed remediation recommendations
 * - Compliance reporting and analytics
 * - Integration with assistive technology testing
 */

import { api } from './api';
import type {
  AccessibilityAudit,
  CreateAuditRequest,
  AuditListQuery,
  AccessibilityFinding,
  ComplianceResult,
  AuditResult,
  Platform,
  AccessibilityStandard,
  AuditStatus,
  BulkAuditRequest,
  BulkAuditResponse,
  AccessibilityMonitoringConfig,
  AccessibilityReport,
  ReportType,
  NotificationConfig,
} from '../types/accessibility';
import { PaginatedResponse } from '../types/api';

/**
 * Extended accessibility API types
 */
export interface AccessibilityTestSuite {
  id: string;
  name: string;
  description: string;
  tests: AccessibilityTest[];
  platform: Platform;
  standard: AccessibilityStandard;
  createdAt: string;
  updatedAt: string;
}

export interface AccessibilityTest {
  id: string;
  name: string;
  description: string;
  category: string;
  wcagCriteria: string[];
  automatable: boolean;
  severity: 'critical' | 'major' | 'moderate' | 'minor';
  testSteps?: string[];
  expectedResult: string;
}

export interface AccessibilityMetrics {
  totalElements: number;
  testedElements: number;
  passedTests: number;
  failedTests: number;
  coverage: number; // percentage
  complianceScore: number; // 0-100
  avgFixTime: number; // minutes
  trendsData: MetricTrend[];
}

export interface MetricTrend {
  date: string;
  complianceScore: number;
  issueCount: number;
  fixedIssues: number;
}

export interface RemediationPlan {
  auditId: string;
  findings: AccessibilityFinding[];
  prioritizedActions: RemediationAction[];
  estimatedEffort: string;
  timeline: string;
  resources: string[];
  dependencies: string[];
}

export interface RemediationAction {
  findingId: string;
  priority: 'immediate' | 'high' | 'medium' | 'low';
  action: string;
  implementation: string;
  estimatedHours: number;
  assignee?: string;
  dueDate?: string;
  status: 'pending' | 'in_progress' | 'completed' | 'deferred';
}

export interface AccessibilityTraining {
  id: string;
  title: string;
  description: string;
  topics: string[];
  duration: number; // minutes
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  platform?: Platform;
  completionUrl: string;
}

/**
 * Accessibility API Service
 *
 * Comprehensive accessibility audit and compliance management system.
 * Supports WCAG 2.1 A/AA/AAA standards across multiple platforms.
 */
export const accessibilityApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Retrieves paginated list of accessibility audits
     */
    getAccessibilityAudits: builder.query<PaginatedResponse<AccessibilityAudit>, AuditListQuery | void>({
      query: (params = {}) => ({
        url: '/accessibility/audits',
        method: 'GET',
        params,
      }),
      providesTags: (result) =>
        result?.items
          ? [
              ...result.items.map(({ id }) => ({ type: 'AccessibilityAudit' as const, id })),
              { type: 'AccessibilityAudit', id: 'LIST' },
            ]
          : [{ type: 'AccessibilityAudit', id: 'LIST' }],
      transformResponse: (response: any): PaginatedResponse<AccessibilityAudit> => ({
        items: response.data?.items || [],
        pagination: response.data?.pagination || {
          hasMore: false,
          total: 0,
        },
      }),
    }),

    /**
     * Retrieves a single accessibility audit with full details
     */
    getAccessibilityAudit: builder.query<AccessibilityAudit, string>({
      query: (auditId) => ({
        url: `/accessibility/audits/${auditId}`,
        method: 'GET',
      }),
      providesTags: (result, error, auditId) => [
        { type: 'AccessibilityAudit', id: auditId },
      ],
      transformResponse: (response: any): AccessibilityAudit => response.data,
    }),

    /**
     * Creates a new accessibility audit
     */
    createAccessibilityAudit: builder.mutation<AccessibilityAudit, CreateAuditRequest>({
      query: (auditData) => ({
        url: '/accessibility/audits',
        method: 'POST',
        body: auditData,
      }),
      invalidatesTags: [{ type: 'AccessibilityAudit', id: 'LIST' }],
      transformResponse: (response: any): AccessibilityAudit => response.data,
    }),

    /**
     * Executes an accessibility audit
     */
    executeAccessibilityAudit: builder.mutation<
      AuditResult,
      { auditId: string; options?: { includeScreenshots?: boolean; deepScan?: boolean } }
    >({
      query: ({ auditId, options }) => ({
        url: `/accessibility/audits/${auditId}/execute`,
        method: 'POST',
        body: options || {},
      }),
      invalidatesTags: (result, error, { auditId }) => [
        { type: 'AccessibilityAudit', id: auditId },
      ],
      transformResponse: (response: any): AuditResult => response.data,
    }),

    /**
     * Bulk executes multiple accessibility audits
     */
    bulkExecuteAudits: builder.mutation<BulkAuditResponse, BulkAuditRequest>({
      query: (bulkRequest) => ({
        url: '/accessibility/audits/bulk-execute',
        method: 'POST',
        body: bulkRequest,
      }),
      invalidatesTags: [{ type: 'AccessibilityAudit', id: 'LIST' }],
      transformResponse: (response: any): BulkAuditResponse => response.data,
    }),

    /**
     * Gets audit results for a specific audit
     */
    getAuditResults: builder.query<
      AuditResult[],
      { auditId: string; platform?: Platform; standard?: AccessibilityStandard }
    >({
      query: ({ auditId, platform, standard }) => ({
        url: `/accessibility/audits/${auditId}/results`,
        method: 'GET',
        params: { platform, standard },
      }),
      providesTags: (result, error, { auditId }) => [
        { type: 'AccessibilityAudit', id: `RESULTS_${auditId}` },
      ],
      transformResponse: (response: any): AuditResult[] => response.data,
    }),

    /**
     * Gets accessibility findings with filtering
     */
    getAccessibilityFindings: builder.query<
      PaginatedResponse<AccessibilityFinding>,
      {
        auditId?: string;
        platform?: Platform;
        severity?: AccessibilityFinding['severity'];
        category?: AccessibilityFinding['category'];
        resolved?: boolean;
        page?: number;
        limit?: number;
      }
    >({
      query: (params) => ({
        url: '/accessibility/findings',
        method: 'GET',
        params,
      }),
      providesTags: [{ type: 'AccessibilityFinding', id: 'LIST' }],
      transformResponse: (response: any): PaginatedResponse<AccessibilityFinding> => ({
        items: response.data?.items || [],
        pagination: response.data?.pagination || {
          hasMore: false,
          total: 0,
        },
      }),
    }),

    /**
     * Updates finding status (resolved, deferred, etc.)
     */
    updateFindingStatus: builder.mutation<
      AccessibilityFinding,
      {
        findingId: string;
        status: 'open' | 'resolved' | 'deferred' | 'false_positive';
        resolution?: string;
        assignee?: string;
      }
    >({
      query: ({ findingId, ...data }) => ({
        url: `/accessibility/findings/${findingId}/status`,
        method: 'PATCH',
        body: data,
      }),
      invalidatesTags: [{ type: 'AccessibilityFinding', id: 'LIST' }],
      transformResponse: (response: any): AccessibilityFinding => response.data,
    }),

    /**
     * Gets compliance report for specific standards
     */
    getComplianceReport: builder.query<
      ComplianceResult,
      {
        auditId?: string;
        componentId?: string;
        platform?: Platform;
        standard: AccessibilityStandard;
        includeDetails?: boolean;
      }
    >({
      query: (params) => ({
        url: '/accessibility/compliance',
        method: 'GET',
        params,
      }),
      transformResponse: (response: any): ComplianceResult => response.data,
    }),

    /**
     * Gets accessibility test suites
     */
    getAccessibilityTestSuites: builder.query<
      PaginatedResponse<AccessibilityTestSuite>,
      { platform?: Platform; standard?: AccessibilityStandard }
    >({
      query: (params) => ({
        url: '/accessibility/test-suites',
        method: 'GET',
        params,
      }),
      providesTags: [{ type: 'AccessibilityTestSuite', id: 'LIST' }],
      transformResponse: (response: any): PaginatedResponse<AccessibilityTestSuite> => ({
        items: response.data?.items || [],
        pagination: response.data?.pagination || {
          hasMore: false,
          total: 0,
        },
      }),
    }),

    /**
     * Creates custom accessibility test suite
     */
    createTestSuite: builder.mutation<AccessibilityTestSuite, Omit<AccessibilityTestSuite, 'id' | 'createdAt' | 'updatedAt'>>({
      query: (suiteData) => ({
        url: '/accessibility/test-suites',
        method: 'POST',
        body: suiteData,
      }),
      invalidatesTags: [{ type: 'AccessibilityTestSuite', id: 'LIST' }],
      transformResponse: (response: any): AccessibilityTestSuite => response.data,
    }),

    /**
     * Gets accessibility metrics and analytics
     */
    getAccessibilityMetrics: builder.query<
      AccessibilityMetrics,
      {
        auditIds?: string[];
        platform?: Platform;
        dateRange?: { from: string; to: string };
        groupBy?: 'day' | 'week' | 'month';
      }
    >({
      query: (params) => ({
        url: '/accessibility/metrics',
        method: 'GET',
        params,
      }),
      // Cache metrics for 10 minutes
      keepUnusedDataFor: 600,
      transformResponse: (response: any): AccessibilityMetrics => response.data,
    }),

    /**
     * Generates remediation plan for audit findings
     */
    generateRemediationPlan: builder.mutation<
      RemediationPlan,
      {
        auditId: string;
        prioritization?: 'severity' | 'impact' | 'effort' | 'wcag_level';
        includeTraining?: boolean;
        teamSize?: number;
      }
    >({
      query: ({ auditId, ...options }) => ({
        url: `/accessibility/audits/${auditId}/remediation-plan`,
        method: 'POST',
        body: options,
      }),
      transformResponse: (response: any): RemediationPlan => response.data,
    }),

    /**
     * Updates remediation action status
     */
    updateRemediationAction: builder.mutation<
      RemediationAction,
      {
        planId: string;
        actionId: string;
        updates: Partial<RemediationAction>;
      }
    >({
      query: ({ planId, actionId, updates }) => ({
        url: `/accessibility/remediation-plans/${planId}/actions/${actionId}`,
        method: 'PATCH',
        body: updates,
      }),
      transformResponse: (response: any): RemediationAction => response.data,
    }),

    /**
     * Configures accessibility monitoring
     */
    configureAccessibilityMonitoring: builder.mutation<
      AccessibilityMonitoringConfig,
      Partial<AccessibilityMonitoringConfig>
    >({
      query: (config) => ({
        url: '/accessibility/monitoring/config',
        method: 'PUT',
        body: config,
      }),
      transformResponse: (response: any): AccessibilityMonitoringConfig => response.data,
    }),

    /**
     * Gets accessibility monitoring configuration
     */
    getAccessibilityMonitoringConfig: builder.query<AccessibilityMonitoringConfig, void>({
      query: () => ({
        url: '/accessibility/monitoring/config',
        method: 'GET',
      }),
      transformResponse: (response: any): AccessibilityMonitoringConfig => response.data,
    }),

    /**
     * Generates accessibility report
     */
    generateAccessibilityReport: builder.mutation<
      AccessibilityReport,
      {
        type: ReportType;
        auditIds?: string[];
        dateRange?: { from: string; to: string };
        format: 'html' | 'pdf' | 'json' | 'csv';
        includeScreenshots?: boolean;
      }
    >({
      query: (reportRequest) => ({
        url: '/accessibility/reports',
        method: 'POST',
        body: reportRequest,
      }),
      transformResponse: (response: any): AccessibilityReport => response.data,
    }),

    /**
     * Gets available accessibility training resources
     */
    getAccessibilityTraining: builder.query<
      PaginatedResponse<AccessibilityTraining>,
      {
        platform?: Platform;
        difficulty?: AccessibilityTraining['difficulty'];
        topic?: string;
      }
    >({
      query: (params) => ({
        url: '/accessibility/training',
        method: 'GET',
        params,
      }),
      transformResponse: (response: any): PaginatedResponse<AccessibilityTraining> => ({
        items: response.data?.items || [],
        pagination: response.data?.pagination || {
          hasMore: false,
          total: 0,
        },
      }),
    }),

    /**
     * Validates accessibility compliance in real-time
     */
    validateAccessibilityRealtime: builder.mutation<
      {
        isCompliant: boolean;
        issues: AccessibilityFinding[];
        quickFixes: string[];
        score: number;
      },
      {
        html?: string;
        url?: string;
        platform: Platform;
        standard: AccessibilityStandard;
        includeWarnings?: boolean;
      }
    >({
      query: (validationData) => ({
        url: '/accessibility/validate/realtime',
        method: 'POST',
        body: validationData,
      }),
      transformResponse: (response: any) => response.data,
    }),

    /**
     * Gets accessibility best practices and guidelines
     */
    getAccessibilityGuidelines: builder.query<
      {
        guidelines: {
          category: string;
          principles: string[];
          techniques: string[];
          examples: { platform: Platform; code: string; description: string }[];
        }[];
        standards: {
          [key in AccessibilityStandard]: {
            description: string;
            requirements: string[];
            testingApproach: string[];
          };
        };
      },
      { platform?: Platform; standard?: AccessibilityStandard }
    >({
      query: (params) => ({
        url: '/accessibility/guidelines',
        method: 'GET',
        params,
      }),
      // Cache guidelines for 1 hour
      keepUnusedDataFor: 3600,
      transformResponse: (response: any) => response.data,
    }),
  }),
  overrideExisting: false,
});

// Export hooks for use in React components
export const {
  // Audit management
  useGetAccessibilityAuditsQuery,
  useGetAccessibilityAuditQuery,
  useCreateAccessibilityAuditMutation,
  useExecuteAccessibilityAuditMutation,
  useBulkExecuteAuditsMutation,

  // Results and findings
  useGetAuditResultsQuery,
  useGetAccessibilityFindingsQuery,
  useUpdateFindingStatusMutation,

  // Compliance and reporting
  useGetComplianceReportQuery,
  useGenerateAccessibilityReportMutation,

  // Test suites
  useGetAccessibilityTestSuitesQuery,
  useCreateTestSuiteMutation,

  // Analytics and metrics
  useGetAccessibilityMetricsQuery,

  // Remediation
  useGenerateRemediationPlanMutation,
  useUpdateRemediationActionMutation,

  // Monitoring
  useConfigureAccessibilityMonitoringMutation,
  useGetAccessibilityMonitoringConfigQuery,

  // Training and guidelines
  useGetAccessibilityTrainingQuery,
  useGetAccessibilityGuidelinesQuery,

  // Real-time validation
  useValidateAccessibilityRealtimeMutation,
} = accessibilityApi;

/**
 * Accessibility utilities
 */
export const accessibilityUtils = {
  /**
   * Calculates overall accessibility score based on WCAG criteria
   */
  calculateAccessibilityScore: (findings: AccessibilityFinding[]): {
    score: number;
    level: 'A' | 'AA' | 'AAA' | 'FAIL';
    breakdown: Record<string, number>;
  } => {
    const totalFindings = findings.length;
    const criticalFindings = findings.filter(f => f.severity === 'critical').length;
    const majorFindings = findings.filter(f => f.severity === 'major').length;
    const moderateFindings = findings.filter(f => f.severity === 'moderate').length;

    // Calculate score based on weighted severity
    const weightedScore = Math.max(0, 100 - (
      criticalFindings * 25 +
      majorFindings * 10 +
      moderateFindings * 5 +
      (totalFindings - criticalFindings - majorFindings - moderateFindings) * 2
    ));

    // Determine WCAG level
    let level: 'A' | 'AA' | 'AAA' | 'FAIL' = 'FAIL';
    if (criticalFindings === 0 && majorFindings === 0) {
      if (moderateFindings === 0) level = 'AAA';
      else if (moderateFindings <= 2) level = 'AA';
      else level = 'A';
    } else if (criticalFindings === 0 && majorFindings <= 1) {
      level = 'A';
    }

    return {
      score: Math.round(weightedScore),
      level,
      breakdown: {
        critical: criticalFindings,
        major: majorFindings,
        moderate: moderateFindings,
        minor: totalFindings - criticalFindings - majorFindings - moderateFindings,
      },
    };
  },

  /**
   * Gets severity color class for UI display
   */
  getSeverityColor: (severity: AccessibilityFinding['severity']): string => {
    switch (severity) {
      case 'critical': return 'bg-red-100 text-red-800';
      case 'major': return 'bg-orange-100 text-orange-800';
      case 'moderate': return 'bg-yellow-100 text-yellow-800';
      case 'minor': return 'bg-blue-100 text-blue-800';
      case 'info': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  },

  /**
   * Gets WCAG level badge color
   */
  getWCAGLevelColor: (level: 'A' | 'AA' | 'AAA' | 'FAIL'): string => {
    switch (level) {
      case 'AAA': return 'bg-green-100 text-green-800';
      case 'AA': return 'bg-blue-100 text-blue-800';
      case 'A': return 'bg-yellow-100 text-yellow-800';
      case 'FAIL': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  },

  /**
   * Formats WCAG criteria for display
   */
  formatWCAGCriteria: (criteria: string[]): string => {
    return criteria.map(c => `WCAG ${c}`).join(', ');
  },

  /**
   * Determines priority level for remediation
   */
  getRemediationPriority: (finding: AccessibilityFinding): 'immediate' | 'high' | 'medium' | 'low' => {
    if (finding.isBlocking || finding.severity === 'critical') return 'immediate';
    if (finding.severity === 'major' || finding.impact === 'blocker') return 'high';
    if (finding.severity === 'moderate' || finding.impact === 'critical') return 'medium';
    return 'low';
  },

  /**
   * Estimates fix effort based on finding type and complexity
   */
  estimateFixEffort: (finding: AccessibilityFinding): { hours: number; difficulty: 'easy' | 'medium' | 'hard' } => {
    // Simplified effort estimation - would use ML/historical data in production
    const baseHours = {
      'color_contrast': 1,
      'text_alternatives': 2,
      'keyboard': 4,
      'focus': 3,
      'aria': 2,
      'structure': 6,
      'forms': 3,
      'navigation': 5,
    };

    const categoryHours = baseHours[finding.category as keyof typeof baseHours] || 3;

    const severityMultiplier = {
      'critical': 1.5,
      'major': 1.2,
      'moderate': 1.0,
      'minor': 0.8,
      'info': 0.5,
    };

    const hours = Math.ceil(categoryHours * severityMultiplier[finding.severity]);
    const difficulty = hours <= 2 ? 'easy' : hours <= 5 ? 'medium' : 'hard';

    return { hours, difficulty };
  },

  /**
   * Groups findings by category for better organization
   */
  groupFindingsByCategory: (findings: AccessibilityFinding[]): Record<string, AccessibilityFinding[]> => {
    return findings.reduce((groups, finding) => {
      const category = finding.category;
      if (!groups[category]) {
        groups[category] = [];
      }
      groups[category].push(finding);
      return groups;
    }, {} as Record<string, AccessibilityFinding[]>);
  },

  /**
   * Validates if audit meets specific compliance standard
   */
  validateCompliance: (
    findings: AccessibilityFinding[],
    standard: AccessibilityStandard
  ): { isCompliant: boolean; missingCriteria: string[]; score: number } => {
    const scoreResult = accessibilityUtils.calculateAccessibilityScore(findings);

    const complianceThresholds = {
      wcag_2_1_a: { level: 'A', minScore: 80 },
      wcag_2_1_aa: { level: 'AA', minScore: 85 },
      wcag_2_1_aaa: { level: 'AAA', minScore: 95 },
      section_508: { level: 'AA', minScore: 85 },
      en_301_549: { level: 'AA', minScore: 85 },
      ada: { level: 'AA', minScore: 85 },
    };

    const threshold = complianceThresholds[standard];
    const isCompliant = scoreResult.score >= threshold.minScore &&
                       (scoreResult.level === threshold.level || scoreResult.level === 'AAA');

    // Extract missing criteria from findings
    const missingCriteria = findings
      .filter(f => f.severity === 'critical' || f.severity === 'major')
      .flatMap(f => f.wcagCriteria.map(c => c.criterion))
      .filter((criterion, index, array) => array.indexOf(criterion) === index);

    return {
      isCompliant,
      missingCriteria,
      score: scoreResult.score,
    };
  },

  /**
   * Generates accessibility checklist for manual testing
   */
  generateManualChecklist: (platform: Platform, standard: AccessibilityStandard): {
    category: string;
    items: { criterion: string; description: string; testSteps: string[] }[];
  }[] => {
    // Simplified checklist generation - would be more comprehensive in production
    const baseChecklist = [
      {
        category: 'Keyboard Navigation',
        items: [
          {
            criterion: '2.1.1',
            description: 'All functionality available from keyboard',
            testSteps: ['Tab through all interactive elements', 'Verify all functions work with keyboard only'],
          },
          {
            criterion: '2.4.7',
            description: 'Focus visible',
            testSteps: ['Tab through interface', 'Verify focus indicator is always visible'],
          },
        ],
      },
      {
        category: 'Color and Contrast',
        items: [
          {
            criterion: '1.4.3',
            description: 'Contrast (Minimum)',
            testSteps: ['Check color contrast ratios', 'Verify 4.5:1 for normal text, 3:1 for large text'],
          },
        ],
      },
    ];

    return baseChecklist;
  },
};