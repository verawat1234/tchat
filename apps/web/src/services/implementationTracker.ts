/**
 * Platform Implementation Status Tracker (T053)
 * Track component implementation status across all platforms
 * Integration with existing Redux store and RTK Query
 */

import { api } from './api';
import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import type { ComponentDefinition, PlatformComponentInfo } from './componentSync';

// Implementation tracking interfaces
export interface ImplementationStatus {
  id: string;
  componentId: string;
  componentName: string;
  platform: 'web' | 'ios' | 'android';
  status: ImplementationState;
  version: string;
  lastUpdated: string;
  updatedBy: string;
  filePath: string;
  lineCount: number;
  dependencies: string[];
  testCoverage: number; // 0.0 - 1.0
  performanceScore: number; // 0.0 - 1.0
  accessibilityScore: number; // 0.0 - 1.0
  designTokenCompliance: number; // 0.0 - 1.0
  buildStatus: BuildStatus;
  deploymentStatus: DeploymentStatus;
  issues: ImplementationIssue[];
  metrics: ImplementationMetrics;
}

export type ImplementationState =
  | 'not_started'
  | 'planning'
  | 'in_progress'
  | 'review'
  | 'testing'
  | 'completed'
  | 'deployed'
  | 'deprecated'
  | 'blocked'
  | 'failed';

export interface BuildStatus {
  status: 'success' | 'failed' | 'in_progress' | 'not_built';
  lastBuild: string;
  buildDuration: number; // seconds
  buildSize: number; // bytes
  errors: string[];
  warnings: string[];
  buildId?: string;
}

export interface DeploymentStatus {
  status: 'deployed' | 'deploying' | 'failed' | 'not_deployed';
  environment: 'development' | 'staging' | 'production';
  lastDeployment: string;
  deploymentDuration: number; // seconds
  url?: string;
  version?: string;
}

export interface ImplementationIssue {
  id: string;
  type: 'bug' | 'performance' | 'accessibility' | 'design_inconsistency' | 'build_error' | 'test_failure';
  severity: 'low' | 'medium' | 'high' | 'critical' | 'constitutional_violation';
  title: string;
  description: string;
  affectedPlatforms: string[];
  createdAt: string;
  resolvedAt?: string;
  assignee?: string;
  estimatedFixTime: number; // minutes
  blocksRelease: boolean;
}

export interface ImplementationMetrics {
  codeQuality: {
    maintainabilityIndex: number;
    cyclomaticComplexity: number;
    linesOfCode: number;
    duplicateCodePercentage: number;
  };
  performance: {
    bundleSize: number; // bytes
    renderTime: number; // ms
    memoryUsage: number; // MB
    loadTime: number; // ms
    frameRate: number; // fps for animations
  };
  accessibility: {
    wcagCompliance: 'A' | 'AA' | 'AAA' | 'non_compliant';
    screenReaderScore: number;
    keyboardNavigationScore: number;
    colorContrastScore: number;
    focusManagementScore: number;
  };
  testing: {
    unitTestCoverage: number; // 0.0 - 1.0
    integrationTestCoverage: number; // 0.0 - 1.0
    e2eTestCoverage: number; // 0.0 - 1.0
    testPassRate: number; // 0.0 - 1.0
    testExecutionTime: number; // seconds
  };
  designSystem: {
    tokenUsageCompliance: number; // 0.0 - 1.0
    crossPlatformConsistency: number; // 0.0 - 1.0
    constitutionalCompliance: boolean;
    variantCoverage: number; // 0.0 - 1.0
  };
}

export interface ImplementationTrackingState {
  implementations: Record<string, ImplementationStatus>;
  platformSummaries: Record<string, PlatformSummary>;
  overallHealth: OverallHealthMetrics;
  realtimeStatus: RealtimeTrackingStatus;
  lastSync: string;
  trackingEnabled: boolean;
  alertsEnabled: boolean;
}

export interface PlatformSummary {
  platform: 'web' | 'ios' | 'android';
  totalComponents: number;
  completedComponents: number;
  inProgressComponents: number;
  blockedComponents: number;
  overallProgress: number; // 0.0 - 1.0
  healthScore: number; // 0.0 - 1.0
  constitutionalCompliance: boolean;
  averageTestCoverage: number;
  averagePerformanceScore: number;
  criticalIssues: number;
  lastUpdated: string;
}

export interface OverallHealthMetrics {
  overallProgress: number; // 0.0 - 1.0
  crossPlatformConsistency: number; // 0.0 - 1.0
  constitutionalCompliance: boolean;
  totalIssues: number;
  criticalIssues: number;
  blockedComponents: string[];
  healthTrend: number[]; // Last 30 days health scores
  performanceTrend: number[]; // Last 30 days performance scores
  testCoverageTrend: number[]; // Last 30 days test coverage
}

export interface RealtimeTrackingStatus {
  isActive: boolean;
  watchedDirectories: string[];
  lastFileChange: string;
  pendingUpdates: string[];
  updateQueue: ImplementationUpdate[];
  syncInProgress: boolean;
  errorsSinceLastSync: number;
}

export interface ImplementationUpdate {
  componentId: string;
  platform: string;
  changeType: 'file_modified' | 'build_completed' | 'test_run' | 'deployment' | 'manual_update';
  timestamp: string;
  metadata: Record<string, any>;
}

export interface TrackingConfiguration {
  platforms: string[];
  trackingInterval: number; // seconds
  enableRealtime: boolean;
  enableAlerts: boolean;
  constitutionalComplianceThreshold: number; // 0.97
  performanceThresholds: {
    loadTime: number; // ms
    bundleSize: number; // bytes
    memoryUsage: number; // MB
    testCoverage: number; // 0.0 - 1.0
  };
  alertThresholds: {
    criticalIssues: number;
    blockedComponents: number;
    constitutionalViolations: number;
    performanceDegradation: number; // percentage
  };
}

// Default tracking configuration
export const DEFAULT_TRACKING_CONFIG: TrackingConfiguration = {
  platforms: ['web', 'ios', 'android'],
  trackingInterval: 300, // 5 minutes
  enableRealtime: true,
  enableAlerts: true,
  constitutionalComplianceThreshold: 0.97,
  performanceThresholds: {
    loadTime: 200, // 200ms constitutional requirement
    bundleSize: 500 * 1024, // 500KB
    memoryUsage: 100, // 100MB mobile, would be adjusted per platform
    testCoverage: 0.8 // 80%
  },
  alertThresholds: {
    criticalIssues: 0,
    blockedComponents: 1,
    constitutionalViolations: 0,
    performanceDegradation: 10 // 10% degradation triggers alert
  }
};

// Redux slice for implementation tracking state
const implementationTrackingSlice = createSlice({
  name: 'implementationTracking',
  initialState: {
    implementations: {},
    platformSummaries: {},
    overallHealth: {
      overallProgress: 0,
      crossPlatformConsistency: 0,
      constitutionalCompliance: false,
      totalIssues: 0,
      criticalIssues: 0,
      blockedComponents: [],
      healthTrend: [],
      performanceTrend: [],
      testCoverageTrend: []
    },
    realtimeStatus: {
      isActive: false,
      watchedDirectories: [],
      lastFileChange: '',
      pendingUpdates: [],
      updateQueue: [],
      syncInProgress: false,
      errorsSinceLastSync: 0
    },
    lastSync: '',
    trackingEnabled: true,
    alertsEnabled: true
  } as ImplementationTrackingState,
  reducers: {
    updateImplementationStatus: (state, action: PayloadAction<ImplementationStatus>) => {
      state.implementations[action.payload.id] = action.payload;
      state.lastSync = new Date().toISOString();
    },
    updatePlatformSummary: (state, action: PayloadAction<PlatformSummary>) => {
      state.platformSummaries[action.payload.platform] = action.payload;
    },
    updateOverallHealth: (state, action: PayloadAction<OverallHealthMetrics>) => {
      state.overallHealth = action.payload;
    },
    updateRealtimeStatus: (state, action: PayloadAction<Partial<RealtimeTrackingStatus>>) => {
      state.realtimeStatus = { ...state.realtimeStatus, ...action.payload };
    },
    addImplementationUpdate: (state, action: PayloadAction<ImplementationUpdate>) => {
      state.realtimeStatus.updateQueue.push(action.payload);
      // Keep only last 100 updates
      if (state.realtimeStatus.updateQueue.length > 100) {
        state.realtimeStatus.updateQueue = state.realtimeStatus.updateQueue.slice(-100);
      }
    },
    setTrackingEnabled: (state, action: PayloadAction<boolean>) => {
      state.trackingEnabled = action.payload;
    },
    setAlertsEnabled: (state, action: PayloadAction<boolean>) => {
      state.alertsEnabled = action.payload;
    },
    clearImplementationUpdates: (state) => {
      state.realtimeStatus.updateQueue = [];
      state.realtimeStatus.pendingUpdates = [];
    }
  }
});

export const {
  updateImplementationStatus,
  updatePlatformSummary,
  updateOverallHealth,
  updateRealtimeStatus,
  addImplementationUpdate,
  setTrackingEnabled,
  setAlertsEnabled,
  clearImplementationUpdates
} = implementationTrackingSlice.actions;

export const implementationTrackingReducer = implementationTrackingSlice.reducer;

// RTK Query API endpoints for implementation tracking
export const implementationTrackingApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // Get all implementation statuses
    getAllImplementationStatuses: builder.query<ImplementationStatus[], { platform?: string }>({
      query: ({ platform }) => ({
        url: '/implementation/statuses',
        params: platform ? { platform } : {}
      }),
      providesTags: ['ImplementationStatus'],
      transformResponse: (response: ImplementationStatus[]) => {
        // Ensure constitutional compliance is properly calculated
        return response.map(impl => ({
          ...impl,
          metrics: {
            ...impl.metrics,
            designSystem: {
              ...impl.metrics.designSystem,
              constitutionalCompliance: impl.metrics.designSystem.crossPlatformConsistency >= 0.97
            }
          }
        }));
      }
    }),

    // Get specific implementation status
    getImplementationStatus: builder.query<ImplementationStatus, { componentId: string; platform: string }>({
      query: ({ componentId, platform }) => `/implementation/status/${componentId}/${platform}`,
      providesTags: (_result, _error, { componentId, platform }) => [
        { type: 'ImplementationStatus', id: `${componentId}-${platform}` }
      ]
    }),

    // Update implementation status
    updateImplementationStatus: builder.mutation<
      ImplementationStatus,
      { componentId: string; platform: string; updates: Partial<ImplementationStatus> }
    >({
      query: ({ componentId, platform, updates }) => ({
        url: `/implementation/status/${componentId}/${platform}`,
        method: 'PUT',
        body: updates
      }),
      invalidatesTags: (_result, _error, { componentId, platform }) => [
        { type: 'ImplementationStatus', id: `${componentId}-${platform}` },
        'ImplementationStatus',
        'PlatformSummary'
      ]
    }),

    // Get platform summaries
    getPlatformSummaries: builder.query<PlatformSummary[], void>({
      query: () => '/implementation/platform-summaries',
      providesTags: ['PlatformSummary'],
      transformResponse: (response: PlatformSummary[]) => {
        // Calculate constitutional compliance for each platform
        return response.map(summary => ({
          ...summary,
          constitutionalCompliance: summary.overallProgress >= 0.97 && summary.healthScore >= 0.97
        }));
      }
    }),

    // Get overall health metrics
    getOverallHealthMetrics: builder.query<OverallHealthMetrics, { timeRange?: string }>({
      query: ({ timeRange = '30d' }) => ({
        url: '/implementation/health',
        params: { timeRange }
      }),
      providesTags: ['OverallHealth'],
      transformResponse: (response: OverallHealthMetrics) => ({
        ...response,
        constitutionalCompliance: response.crossPlatformConsistency >= 0.97
      })
    }),

    // Get implementation issues
    getImplementationIssues: builder.query<
      ImplementationIssue[],
      { platform?: string; severity?: string; status?: 'open' | 'resolved' }
    >({
      query: (params) => ({
        url: '/implementation/issues',
        params
      }),
      providesTags: ['ImplementationIssues']
    }),

    // Create or update implementation issue
    upsertImplementationIssue: builder.mutation<
      ImplementationIssue,
      { issue: Omit<ImplementationIssue, 'id' | 'createdAt'> }
    >({
      query: ({ issue }) => ({
        url: '/implementation/issues',
        method: 'POST',
        body: issue
      }),
      invalidatesTags: ['ImplementationIssues', 'OverallHealth', 'PlatformSummary']
    }),

    // Resolve implementation issue
    resolveImplementationIssue: builder.mutation<
      { success: boolean },
      { issueId: string; resolution: string }
    >({
      query: ({ issueId, resolution }) => ({
        url: `/implementation/issues/${issueId}/resolve`,
        method: 'POST',
        body: { resolution }
      }),
      invalidatesTags: ['ImplementationIssues', 'OverallHealth', 'PlatformSummary']
    }),

    // Start realtime tracking
    startRealtimeTracking: builder.mutation<
      RealtimeTrackingStatus,
      { config: Partial<TrackingConfiguration> }
    >({
      query: ({ config }) => ({
        url: '/implementation/realtime/start',
        method: 'POST',
        body: config
      })
    }),

    // Stop realtime tracking
    stopRealtimeTracking: builder.mutation<{ success: boolean }, void>({
      query: () => ({
        url: '/implementation/realtime/stop',
        method: 'POST'
      })
    }),

    // Get realtime tracking status
    getRealtimeTrackingStatus: builder.query<RealtimeTrackingStatus, void>({
      query: () => '/implementation/realtime/status',
      pollingInterval: 10000 // Poll every 10 seconds
    }),

    // Trigger sync for specific component/platform
    triggerComponentSync: builder.mutation<
      { success: boolean; syncId: string },
      { componentId: string; platform: string; syncType: 'full' | 'metadata' | 'metrics' }
    >({
      query: ({ componentId, platform, syncType }) => ({
        url: '/implementation/sync',
        method: 'POST',
        body: { componentId, platform, syncType }
      }),
      invalidatesTags: (_result, _error, { componentId, platform }) => [
        { type: 'ImplementationStatus', id: `${componentId}-${platform}` },
        'ImplementationStatus'
      ]
    }),

    // Batch update multiple implementations
    batchUpdateImplementations: builder.mutation<
      { successful: string[]; failed: string[] },
      { updates: Array<{ componentId: string; platform: string; updates: Partial<ImplementationStatus> }> }
    >({
      query: ({ updates }) => ({
        url: '/implementation/batch-update',
        method: 'POST',
        body: { updates }
      }),
      invalidatesTags: ['ImplementationStatus', 'PlatformSummary', 'OverallHealth']
    }),

    // Generate implementation report
    generateImplementationReport: builder.mutation<
      {
        reportId: string;
        url: string;
        format: 'json' | 'pdf' | 'csv';
        generatedAt: string;
      },
      {
        platforms: string[];
        includeMetrics: boolean;
        includeIssues: boolean;
        format: 'json' | 'pdf' | 'csv';
        timeRange: string;
      }
    >({
      query: (params) => ({
        url: '/implementation/reports/generate',
        method: 'POST',
        body: params
      })
    }),

    // Get component dependencies and impact analysis
    getComponentDependencies: builder.query<
      {
        componentId: string;
        dependencies: string[];
        dependents: string[];
        impactAnalysis: {
          affectedComponents: string[];
          estimatedUpdateTime: number;
          riskLevel: 'low' | 'medium' | 'high';
        };
      },
      { componentId: string; platform: string }
    >({
      query: ({ componentId, platform }) => ({
        url: `/implementation/dependencies/${componentId}`,
        params: { platform }
      }),
      providesTags: (_result, _error, { componentId }) => [
        { type: 'ComponentDependencies', id: componentId }
      ]
    })
  })
});

// Export hooks
export const {
  useGetAllImplementationStatusesQuery,
  useGetImplementationStatusQuery,
  useUpdateImplementationStatusMutation,
  useGetPlatformSummariesQuery,
  useGetOverallHealthMetricsQuery,
  useGetImplementationIssuesQuery,
  useUpsertImplementationIssueMutation,
  useResolveImplementationIssueMutation,
  useStartRealtimeTrackingMutation,
  useStopRealtimeTrackingMutation,
  useGetRealtimeTrackingStatusQuery,
  useTriggerComponentSyncMutation,
  useBatchUpdateImplementationsMutation,
  useGenerateImplementationReportMutation,
  useGetComponentDependenciesQuery
} = implementationTrackingApi;

// Implementation Tracker Service Class
export class ImplementationTracker {
  private config: TrackingConfiguration;
  private updateCallbacks: Map<string, (update: ImplementationUpdate) => void> = new Map();
  private healthCallbacks: Map<string, (health: OverallHealthMetrics) => void> = new Map();

  constructor(config: TrackingConfiguration = DEFAULT_TRACKING_CONFIG) {
    this.config = config;
  }

  /**
   * Subscribe to implementation updates for a component/platform
   */
  subscribeToUpdates(
    key: string,
    callback: (update: ImplementationUpdate) => void
  ): () => void {
    this.updateCallbacks.set(key, callback);
    return () => this.updateCallbacks.delete(key);
  }

  /**
   * Subscribe to overall health metric changes
   */
  subscribeToHealthUpdates(
    key: string,
    callback: (health: OverallHealthMetrics) => void
  ): () => void {
    this.healthCallbacks.set(key, callback);
    return () => this.healthCallbacks.delete(key);
  }

  /**
   * Calculate implementation progress for a platform
   */
  calculatePlatformProgress(implementations: ImplementationStatus[]): PlatformSummary {
    const platformImpls = implementations.filter(impl =>
      impl.platform === implementations[0]?.platform
    );

    const total = platformImpls.length;
    const completed = platformImpls.filter(impl => impl.status === 'completed' || impl.status === 'deployed').length;
    const inProgress = platformImpls.filter(impl => impl.status === 'in_progress' || impl.status === 'review' || impl.status === 'testing').length;
    const blocked = platformImpls.filter(impl => impl.status === 'blocked' || impl.status === 'failed').length;

    const overallProgress = total > 0 ? completed / total : 0;
    const avgTestCoverage = platformImpls.reduce((sum, impl) => sum + impl.testCoverage, 0) / total;
    const avgPerformanceScore = platformImpls.reduce((sum, impl) => sum + impl.performanceScore, 0) / total;
    const criticalIssues = platformImpls.reduce((sum, impl) =>
      sum + impl.issues.filter(issue => issue.severity === 'critical' || issue.severity === 'constitutional_violation').length, 0
    );

    const healthScore = this.calculateHealthScore({
      progress: overallProgress,
      testCoverage: avgTestCoverage,
      performanceScore: avgPerformanceScore,
      criticalIssues,
      blockedComponents: blocked
    });

    return {
      platform: platformImpls[0]?.platform || 'web',
      totalComponents: total,
      completedComponents: completed,
      inProgressComponents: inProgress,
      blockedComponents: blocked,
      overallProgress,
      healthScore,
      constitutionalCompliance: healthScore >= 0.97 && avgPerformanceScore >= 0.90,
      averageTestCoverage: avgTestCoverage,
      averagePerformanceScore: avgPerformanceScore,
      criticalIssues,
      lastUpdated: new Date().toISOString()
    };
  }

  /**
   * Calculate overall health score
   */
  private calculateHealthScore(metrics: {
    progress: number;
    testCoverage: number;
    performanceScore: number;
    criticalIssues: number;
    blockedComponents: number;
  }): number {
    const progressWeight = 0.3;
    const testCoverageWeight = 0.25;
    const performanceWeight = 0.25;
    const qualityWeight = 0.2;

    const qualityScore = Math.max(0, 1.0 - (metrics.criticalIssues * 0.1) - (metrics.blockedComponents * 0.05));

    return Math.min(1.0,
      (metrics.progress * progressWeight) +
      (metrics.testCoverage * testCoverageWeight) +
      (metrics.performanceScore * performanceWeight) +
      (qualityScore * qualityWeight)
    );
  }

  /**
   * Generate constitutional compliance report
   */
  async generateConstitutionalComplianceReport(implementations: ImplementationStatus[]): Promise<{
    overallCompliance: boolean;
    complianceScore: number;
    violations: ImplementationIssue[];
    recommendations: string[];
    platformBreakdown: Record<string, { compliant: boolean; score: number; issues: number }>;
  }> {
    const violations: ImplementationIssue[] = [];
    const platformBreakdown: Record<string, { compliant: boolean; score: number; issues: number }> = {};

    let totalScore = 0;
    let totalImplementations = 0;

    // Group by platform for analysis
    const byPlatform = implementations.reduce((acc, impl) => {
      if (!acc[impl.platform]) acc[impl.platform] = [];
      acc[impl.platform].push(impl);
      return acc;
    }, {} as Record<string, ImplementationStatus[]>);

    for (const [platform, impls] of Object.entries(byPlatform)) {
      let platformScore = 0;
      let platformIssues = 0;

      for (const impl of impls) {
        // Calculate constitutional compliance score
        const consistencyScore = impl.designTokenCompliance;
        const performanceScore = impl.performanceScore;
        const accessibilityScore = impl.accessibilityScore;

        const implScore = (consistencyScore + performanceScore + accessibilityScore) / 3;
        platformScore += implScore;
        totalScore += implScore;
        totalImplementations++;

        // Check for constitutional violations
        const constitutionalViolations = impl.issues.filter(
          issue => issue.severity === 'constitutional_violation'
        );

        violations.push(...constitutionalViolations);
        platformIssues += constitutionalViolations.length;
      }

      const avgPlatformScore = impls.length > 0 ? platformScore / impls.length : 0;
      platformBreakdown[platform] = {
        compliant: avgPlatformScore >= 0.97 && platformIssues === 0,
        score: avgPlatformScore,
        issues: platformIssues
      };
    }

    const overallScore = totalImplementations > 0 ? totalScore / totalImplementations : 0;
    const overallCompliance = overallScore >= 0.97 && violations.length === 0;

    const recommendations: string[] = [];
    if (!overallCompliance) {
      recommendations.push('Achieve 97% cross-platform consistency to meet constitutional requirements');

      if (violations.length > 0) {
        recommendations.push(`Resolve ${violations.length} constitutional violations immediately`);
      }

      for (const [platform, breakdown] of Object.entries(platformBreakdown)) {
        if (!breakdown.compliant) {
          recommendations.push(`Improve ${platform} implementation to achieve constitutional compliance`);
        }
      }
    }

    return {
      overallCompliance,
      complianceScore: overallScore,
      violations,
      recommendations,
      platformBreakdown
    };
  }

  /**
   * Predict implementation risks based on current status
   */
  predictImplementationRisks(implementations: ImplementationStatus[]): {
    highRiskComponents: string[];
    potentialDelays: Array<{ componentId: string; estimatedDelay: number; reason: string }>;
    resourceBottlenecks: string[];
    recommendations: string[];
  } {
    const highRiskComponents: string[] = [];
    const potentialDelays: Array<{ componentId: string; estimatedDelay: number; reason: string }> = [];
    const resourceBottlenecks: string[] = [];
    const recommendations: string[] = [];

    for (const impl of implementations) {
      // Identify high-risk components
      const criticalIssueCount = impl.issues.filter(i => i.severity === 'critical').length;
      const constitutionalViolations = impl.issues.filter(i => i.severity === 'constitutional_violation').length;

      if (criticalIssueCount > 0 || constitutionalViolations > 0 || impl.status === 'blocked') {
        highRiskComponents.push(impl.componentId);
      }

      // Predict potential delays
      if (impl.testCoverage < 0.8) {
        potentialDelays.push({
          componentId: impl.componentId,
          estimatedDelay: Math.ceil((0.8 - impl.testCoverage) * 100), // Rough estimate in hours
          reason: 'Insufficient test coverage'
        });
      }

      if (impl.performanceScore < 0.7) {
        potentialDelays.push({
          componentId: impl.componentId,
          estimatedDelay: 40, // Estimate 40 hours for performance optimization
          reason: 'Performance optimization required'
        });
      }

      // Identify resource bottlenecks
      if (impl.buildStatus.status === 'failed') {
        resourceBottlenecks.push(`Build system issues for ${impl.componentId}`);
      }
    }

    // Generate recommendations
    if (highRiskComponents.length > 0) {
      recommendations.push(`Focus on ${highRiskComponents.length} high-risk components immediately`);
    }

    if (potentialDelays.length > 0) {
      const totalDelay = potentialDelays.reduce((sum, delay) => sum + delay.estimatedDelay, 0);
      recommendations.push(`Address potential delays totaling ${totalDelay} hours`);
    }

    return {
      highRiskComponents,
      potentialDelays,
      resourceBottlenecks,
      recommendations
    };
  }
}

// Singleton instance
export const implementationTracker = new ImplementationTracker();

// Utility functions
export const ImplementationTrackerUtils = {
  formatImplementationStatus: (status: ImplementationState): string => {
    const statusMap: Record<ImplementationState, string> = {
      'not_started': '📝 Not Started',
      'planning': '📋 Planning',
      'in_progress': '🔄 In Progress',
      'review': '👀 Review',
      'testing': '🧪 Testing',
      'completed': '✅ Completed',
      'deployed': '🚀 Deployed',
      'deprecated': '🗑️ Deprecated',
      'blocked': '🚫 Blocked',
      'failed': '❌ Failed'
    };
    return statusMap[status] || status;
  },

  getStatusColor: (status: ImplementationState): string => {
    const colorMap: Record<ImplementationState, string> = {
      'not_started': 'text-gray-500',
      'planning': 'text-blue-500',
      'in_progress': 'text-yellow-500',
      'review': 'text-purple-500',
      'testing': 'text-orange-500',
      'completed': 'text-green-500',
      'deployed': 'text-green-600',
      'deprecated': 'text-gray-400',
      'blocked': 'text-red-500',
      'failed': 'text-red-600'
    };
    return colorMap[status] || 'text-gray-500';
  },

  calculateProgressPercentage: (status: ImplementationState): number => {
    const progressMap: Record<ImplementationState, number> = {
      'not_started': 0,
      'planning': 10,
      'in_progress': 50,
      'review': 80,
      'testing': 90,
      'completed': 100,
      'deployed': 100,
      'deprecated': 100,
      'blocked': 0,
      'failed': 0
    };
    return progressMap[status] || 0;
  },

  formatHealthScore: (score: number): string => {
    const percentage = (score * 100).toFixed(1);
    const emoji = score >= 0.9 ? '💚' : score >= 0.7 ? '💛' : '❤️';
    return `${emoji} ${percentage}%`;
  },

  isConstitutionalViolation: (issue: ImplementationIssue): boolean => {
    return issue.severity === 'constitutional_violation';
  }
};

export default implementationTracker;