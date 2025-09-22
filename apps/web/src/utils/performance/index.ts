/**
 * Performance Validation System - Main Export
 * 
 * Comprehensive performance validation system for ensuring <200ms content load times
 * Includes measurement, monitoring, caching validation, network testing, and reporting
 */

// Core measurement utilities
export {
  PerformanceMeasurement,
  getCoreWebVitals,
  measureExecutionTime,
  withPerformanceMonitoring,
  useContentPerformance,
  PERFORMANCE_BUDGET,
  type ContentLoadMetrics,
  type PerformanceMetrics,
  type PerformanceBudget,
} from "./measurement";

// Monitoring and tracking
export {
  PerformanceMonitor,
  createPerformanceMonitor,
  usePerformanceMonitoring,
  DEFAULT_MONITORING_CONFIG,
  DEFAULT_ALERTS,
  type MonitoringConfig,
  type PerformanceAlert,
  type MonitoringDashboard,
  type PerformanceTrend,
} from "./monitoring";

// Cache validation
export {
  CacheValidator,
  MockCache,
  type CacheConfig,
  type CacheStats,
  type CacheValidationResult,
} from "./cache-validation.test";

// Network testing
export {
  NetworkSimulator,
  NETWORK_CONDITIONS,
  CONTENT_SIZES,
  type NetworkCondition,
} from "./network-testing.test";

// Reporting and analysis
export {
  PerformanceReportGenerator,
  createPerformanceReportGenerator,
  type PerformanceReport,
  type PerformanceReportSummary,
  type PerformanceRecommendation,
  type PerformanceChart,
  type PerformanceTable,
} from "./reporting";

/**
 * Quick performance validation for React components
 */
export function usePerformanceValidation(componentName: string) {
  const React = require("react");
  const { useContentPerformance } = require("./measurement");
  const { usePerformanceMonitoring } = require("./monitoring");

  const {
    metrics: contentMetrics,
    startMeasurement,
    endMeasurement,
    isWithinBudget,
  } = useContentPerformance(componentName, "component", "ui");

  const {
    monitor,
    dashboard,
    alerts,
    startMonitoring,
    endMonitoring,
  } = usePerformanceMonitoring();

  React.useEffect(() => {
    startMeasurement();
    startMonitoring(componentName, "component", "ui");

    return () => {
      endMeasurement();
      endMonitoring(componentName);
    };
  }, [componentName, startMeasurement, endMeasurement, startMonitoring, endMonitoring]);

  return {
    contentMetrics,
    dashboard,
    alerts,
    isWithinBudget,
    monitor,
  };
}

/**
 * Performance validation configuration for different environments
 */
export const PERFORMANCE_CONFIGS = {
  development: {
    enabled: true,
    sampleRate: 1.0, // Monitor all requests in development
    alerts: {
      loadTimeWarning: 300, // More lenient in development
      cacheHitRatio: 0.5,
    },
  },
  staging: {
    enabled: true,
    sampleRate: 0.5, // Monitor 50% of requests in staging
    alerts: {
      loadTimeWarning: 250,
      cacheHitRatio: 0.7,
    },
  },
  production: {
    enabled: true,
    sampleRate: 0.1, // Monitor 10% of requests in production
    alerts: {
      loadTimeWarning: 200, // Strict budget enforcement
      cacheHitRatio: 0.8,
    },
  },
} as const;

/**
 * Get performance configuration for current environment
 */
export function getPerformanceConfig(environment: keyof typeof PERFORMANCE_CONFIGS = "production") {
  return PERFORMANCE_CONFIGS[environment];
}

/**
 * Performance validation presets for common use cases
 */
export const VALIDATION_PRESETS = {
  contentManagement: {
    contentTypes: ["text", "richText", "config", "translation"],
    networkConditions: ["wifi-fast", "4g-good", "3g"],
    cacheStrategies: ["memory", "localStorage", "sessionStorage"],
    performanceBudget: {
      maxLoadTime: 200,
      maxResourceLoadTime: 150,
      minCacheHitRatio: 0.8,
    },
  },
  ecommerce: {
    contentTypes: ["text", "config", "imageUrl"],
    networkConditions: ["wifi-fast", "4g-excellent", "4g-good", "3g"],
    cacheStrategies: ["memory", "localStorage", "indexedDB"],
    performanceBudget: {
      maxLoadTime: 150, // Stricter for e-commerce
      maxResourceLoadTime: 100,
      minCacheHitRatio: 0.9,
    },
  },
  mobile: {
    contentTypes: ["text", "config"],
    networkConditions: ["4g-excellent", "4g-good", "4g-fair", "3g", "2g"],
    cacheStrategies: ["memory", "localStorage"],
    performanceBudget: {
      maxLoadTime: 300, // More lenient for mobile
      maxResourceLoadTime: 200,
      minCacheHitRatio: 0.7,
    },
  },
} as const;

/**
 * Get validation preset for specific use case
 */
export function getValidationPreset(useCase: keyof typeof VALIDATION_PRESETS) {
  return VALIDATION_PRESETS[useCase];
}

/**
 * Summary of the performance validation system
 */
export const PERFORMANCE_VALIDATION_SUMMARY = {
  purpose: "Comprehensive validation of <200ms content load time requirement",
  features: [
    "Browser Performance API measurement utilities",
    "Performance benchmarking for different content types",
    "Network condition testing scenarios",
    "Caching effectiveness validation",
    "Real-time performance monitoring and alerting",
    "Comprehensive reporting and analysis tools",
    "Regression testing and validation",
  ],
  testCoverage: [
    "Content load performance across all content types",
    "Cache hit ratios and effectiveness",
    "Network resilience across different conditions",
    "Real-world usage pattern simulation",
    "Performance regression detection",
    "Cross-browser compatibility",
  ],
  requirements: {
    primaryBudget: "200ms maximum content load time",
    budgetCompliance: "80% of requests must meet budget",
    cacheHitRatio: "80% minimum cache hit ratio",
    performanceScore: "70+ overall performance score",
  },
  integration: [
    "React hooks for component-level monitoring",
    "RTK Query integration for API performance",
    "Automated CI/CD pipeline validation",
    "Real-time dashboard monitoring",
    "Alert system for performance degradation",
  ],
} as const;
