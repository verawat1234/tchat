/**
 * Performance Monitor Utility
 *
 * Comprehensive performance monitoring system for Stream content loading.
 * Tracks metrics, identifies bottlenecks, and provides real-time performance insights.
 *
 * Metrics Tracked:
 * - Content load times by category and type
 * - UI transition performance and frame rates
 * - Memory usage and cache efficiency
 * - Network request timing and success rates
 * - User interaction response times
 * - Error rates and recovery times
 *
 * Features:
 * - Real-time performance dashboard
 * - Automatic bottleneck detection
 * - Performance budgets and alerting
 * - Session-based analytics
 * - Export capabilities for analysis
 */

interface PerformanceMetric {
  name: string;
  value: number;
  timestamp: number;
  category: 'load' | 'transition' | 'interaction' | 'memory' | 'network' | 'error';
  metadata?: Record<string, any>;
}

interface PerformanceSession {
  sessionId: string;
  startTime: number;
  endTime?: number;
  metrics: PerformanceMetric[];
  budgetViolations: BudgetViolation[];
  summary: PerformanceSummary;
}

interface BudgetViolation {
  metric: string;
  expected: number;
  actual: number;
  severity: 'warning' | 'critical';
  timestamp: number;
}

interface PerformanceSummary {
  avgLoadTime: number;
  avgTransitionTime: number;
  errorRate: number;
  cacheHitRate: number;
  memoryUsage: number;
  totalOperations: number;
}

interface PerformanceBudgets {
  contentLoad: number;      // 1000ms
  tabTransition: number;    // 200ms
  imageLoad: number;        // 500ms
  apiResponse: number;      // 300ms
  memoryUsage: number;      // 100MB mobile, 500MB desktop
  errorRate: number;        // 1%
  cacheHitRate: number;     // 95%
}

export class PerformanceMonitor {
  private session: PerformanceSession;
  private budgets: PerformanceBudgets;
  private observers: Map<string, PerformanceObserver>;
  private isMonitoring: boolean = false;

  // Performance budgets based on target requirements
  private static readonly DEFAULT_BUDGETS: PerformanceBudgets = {
    contentLoad: 1000,    // 1s content load target
    tabTransition: 200,   // 200ms transition target
    imageLoad: 500,       // 500ms image load target
    apiResponse: 300,     // 300ms API response target
    memoryUsage: window.innerWidth <= 768 ? 100 : 500, // MB
    errorRate: 0.01,      // 1%
    cacheHitRate: 0.95,   // 95%
  };

  constructor(customBudgets?: Partial<PerformanceBudgets>) {
    this.budgets = { ...PerformanceMonitor.DEFAULT_BUDGETS, ...customBudgets };
    this.observers = new Map();
    this.session = this.createNewSession();

    this.initializeObservers();
  }

  /**
   * Start a new performance monitoring session
   */
  startSession(): void {
    if (this.isMonitoring) {
      this.endSession();
    }

    this.session = this.createNewSession();
    this.isMonitoring = true;

    // Start Core Web Vitals monitoring
    this.startCoreWebVitalsMonitoring();

    // Start memory monitoring
    this.startMemoryMonitoring();

    console.log(`Performance monitoring started - Session: ${this.session.sessionId}`);
  }

  /**
   * End the current performance monitoring session
   */
  endSession(): PerformanceSummary {
    if (!this.isMonitoring) return this.session.summary;

    this.session.endTime = performance.now();
    this.isMonitoring = false;

    // Generate session summary
    this.session.summary = this.generateSummary();

    // Stop all observers
    this.observers.forEach(observer => observer.disconnect());
    this.observers.clear();

    // Log session results
    this.logSessionResults();

    return this.session.summary;
  }

  /**
   * Record a performance operation
   */
  recordOperation(
    name: string,
    duration: number,
    category: PerformanceMetric['category'] = 'interaction',
    metadata?: Record<string, any>
  ): void {
    const metric: PerformanceMetric = {
      name,
      value: duration,
      timestamp: performance.now(),
      category,
      metadata
    };

    this.session.metrics.push(metric);

    // Check against budgets
    this.checkPerformanceBudget(metric);

    // Real-time performance alerts
    if (this.shouldAlert(metric)) {
      this.emitPerformanceAlert(metric);
    }
  }

  /**
   * Record an error with timing
   */
  recordError(operation: string, duration: number, error?: any): void {
    this.recordOperation(operation, duration, 'error', {
      error: error?.message || 'Unknown error',
      stack: error?.stack
    });
  }

  /**
   * Get current session metrics
   */
  getSessionMetrics(): PerformanceMetric[] {
    return [...this.session.metrics];
  }

  /**
   * Get performance summary
   */
  getPerformanceSummary(): PerformanceSummary {
    return this.generateSummary();
  }

  /**
   * Get budget violations
   */
  getBudgetViolations(): BudgetViolation[] {
    return [...this.session.budgetViolations];
  }

  /**
   * Export session data for analysis
   */
  exportSessionData(): string {
    return JSON.stringify({
      session: this.session,
      budgets: this.budgets,
      userAgent: navigator.userAgent,
      viewport: {
        width: window.innerWidth,
        height: window.innerHeight
      },
      timestamp: new Date().toISOString()
    }, null, 2);
  }

  /**
   * Get real-time performance dashboard data
   */
  getDashboardData() {
    const summary = this.generateSummary();
    const recentMetrics = this.session.metrics.slice(-20); // Last 20 operations

    return {
      summary,
      recentMetrics,
      budgetViolations: this.session.budgetViolations.slice(-10),
      charts: {
        loadTimes: this.getMetricsByCategory('load'),
        transitions: this.getMetricsByCategory('transition'),
        errors: this.getMetricsByCategory('error')
      },
      status: this.getOverallStatus()
    };
  }

  // Private methods

  private createNewSession(): PerformanceSession {
    return {
      sessionId: `perf_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      startTime: performance.now(),
      metrics: [],
      budgetViolations: [],
      summary: {
        avgLoadTime: 0,
        avgTransitionTime: 0,
        errorRate: 0,
        cacheHitRate: 0,
        memoryUsage: 0,
        totalOperations: 0
      }
    };
  }

  private initializeObservers(): void {
    // Navigation timing observer
    if ('PerformanceObserver' in window) {
      try {
        const navObserver = new PerformanceObserver((list) => {
          list.getEntries().forEach((entry) => {
            if (entry.entryType === 'navigation') {
              this.recordNavigationMetrics(entry as PerformanceNavigationTiming);
            }
          });
        });

        navObserver.observe({ entryTypes: ['navigation'] });
        this.observers.set('navigation', navObserver);
      } catch (error) {
        console.warn('Navigation observer failed:', error);
      }
    }

    // Resource timing observer
    if ('PerformanceObserver' in window) {
      try {
        const resourceObserver = new PerformanceObserver((list) => {
          list.getEntries().forEach((entry) => {
            if (entry.entryType === 'resource') {
              this.recordResourceMetrics(entry as PerformanceResourceTiming);
            }
          });
        });

        resourceObserver.observe({ entryTypes: ['resource'] });
        this.observers.set('resource', resourceObserver);
      } catch (error) {
        console.warn('Resource observer failed:', error);
      }
    }
  }

  private startCoreWebVitalsMonitoring(): void {
    // Monitor Core Web Vitals if available
    if ('web-vitals' in window || typeof webVitals !== 'undefined') {
      // This would integrate with web-vitals library if available
      console.log('Core Web Vitals monitoring would be initialized here');
    }

    // Manual FCP monitoring
    if ('PerformanceObserver' in window) {
      try {
        const paintObserver = new PerformanceObserver((list) => {
          list.getEntries().forEach((entry) => {
            if (entry.name === 'first-contentful-paint') {
              this.recordOperation('first_contentful_paint', entry.startTime, 'load');
            }
          });
        });

        paintObserver.observe({ entryTypes: ['paint'] });
        this.observers.set('paint', paintObserver);
      } catch (error) {
        console.warn('Paint observer failed:', error);
      }
    }
  }

  private startMemoryMonitoring(): void {
    // Memory monitoring for supported browsers
    if ('memory' in performance) {
      const checkMemory = () => {
        if (!this.isMonitoring) return;

        const memory = (performance as any).memory;
        const memoryUsageMB = memory.usedJSHeapSize / (1024 * 1024);

        this.recordOperation('memory_usage', memoryUsageMB, 'memory');

        // Schedule next check
        setTimeout(checkMemory, 5000); // Every 5 seconds
      };

      checkMemory();
    }
  }

  private recordNavigationMetrics(entry: PerformanceNavigationTiming): void {
    // Record various navigation timing metrics
    this.recordOperation('dns_lookup', entry.domainLookupEnd - entry.domainLookupStart, 'network');
    this.recordOperation('tcp_connect', entry.connectEnd - entry.connectStart, 'network');
    this.recordOperation('request_response', entry.responseEnd - entry.requestStart, 'network');
    this.recordOperation('dom_content_loaded', entry.domContentLoadedEventEnd - entry.domContentLoadedEventStart, 'load');
    this.recordOperation('page_load', entry.loadEventEnd - entry.loadEventStart, 'load');
  }

  private recordResourceMetrics(entry: PerformanceResourceTiming): void {
    const duration = entry.responseEnd - entry.requestStart;

    // Categorize resources
    if (entry.name.includes('/api/')) {
      this.recordOperation('api_request', duration, 'network', {
        url: entry.name,
        size: entry.transferSize
      });
    } else if (entry.name.match(/\.(jpg|jpeg|png|gif|webp|svg)$/i)) {
      this.recordOperation('image_load', duration, 'load', {
        url: entry.name,
        size: entry.transferSize
      });
    }
  }

  private checkPerformanceBudget(metric: PerformanceMetric): void {
    let budgetKey: keyof PerformanceBudgets | null = null;
    let budgetValue: number | null = null;

    // Map metrics to budgets
    switch (metric.name) {
      case 'category_transition':
      case 'subtab_transition':
        budgetKey = 'tabTransition';
        budgetValue = this.budgets.tabTransition;
        break;
      case 'content_load':
      case 'featured_load':
        budgetKey = 'contentLoad';
        budgetValue = this.budgets.contentLoad;
        break;
      case 'image_load':
        budgetKey = 'imageLoad';
        budgetValue = this.budgets.imageLoad;
        break;
      case 'api_request':
        budgetKey = 'apiResponse';
        budgetValue = this.budgets.apiResponse;
        break;
      case 'memory_usage':
        budgetKey = 'memoryUsage';
        budgetValue = this.budgets.memoryUsage;
        break;
    }

    if (budgetKey && budgetValue && metric.value > budgetValue) {
      const violation: BudgetViolation = {
        metric: metric.name,
        expected: budgetValue,
        actual: metric.value,
        severity: metric.value > budgetValue * 1.5 ? 'critical' : 'warning',
        timestamp: metric.timestamp
      };

      this.session.budgetViolations.push(violation);
    }
  }

  private shouldAlert(metric: PerformanceMetric): boolean {
    // Alert on critical performance issues
    return (
      (metric.name === 'category_transition' && metric.value > 500) ||
      (metric.name === 'content_load' && metric.value > 2000) ||
      (metric.name === 'memory_usage' && metric.value > this.budgets.memoryUsage * 1.2) ||
      (metric.category === 'error')
    );
  }

  private emitPerformanceAlert(metric: PerformanceMetric): void {
    console.warn(`Performance Alert: ${metric.name} took ${metric.value}ms`, metric);

    // Could emit custom events for external monitoring
    window.dispatchEvent(new CustomEvent('performance-alert', {
      detail: metric
    }));
  }

  private generateSummary(): PerformanceSummary {
    const loadMetrics = this.getMetricsByCategory('load');
    const transitionMetrics = this.getMetricsByCategory('transition');
    const errorMetrics = this.getMetricsByCategory('error');
    const memoryMetrics = this.getMetricsByCategory('memory');

    return {
      avgLoadTime: this.calculateAverage(loadMetrics),
      avgTransitionTime: this.calculateAverage(transitionMetrics),
      errorRate: errorMetrics.length / Math.max(this.session.metrics.length, 1),
      cacheHitRate: this.calculateCacheHitRate(),
      memoryUsage: memoryMetrics.length > 0 ? memoryMetrics[memoryMetrics.length - 1].value : 0,
      totalOperations: this.session.metrics.length
    };
  }

  private getMetricsByCategory(category: PerformanceMetric['category']): PerformanceMetric[] {
    return this.session.metrics.filter(metric => metric.category === category);
  }

  private calculateAverage(metrics: PerformanceMetric[]): number {
    if (metrics.length === 0) return 0;
    return metrics.reduce((sum, metric) => sum + metric.value, 0) / metrics.length;
  }

  private calculateCacheHitRate(): number {
    // Estimate cache hit rate based on API request patterns
    const apiMetrics = this.session.metrics.filter(m => m.name === 'api_request');
    if (apiMetrics.length === 0) return 1;

    // Simple heuristic: fast requests likely cache hits
    const fastRequests = apiMetrics.filter(m => m.value < 100);
    return fastRequests.length / apiMetrics.length;
  }

  private getOverallStatus(): 'excellent' | 'good' | 'poor' | 'critical' {
    const summary = this.generateSummary();
    const violations = this.session.budgetViolations;

    const criticalViolations = violations.filter(v => v.severity === 'critical').length;
    const warningViolations = violations.filter(v => v.severity === 'warning').length;

    if (criticalViolations > 0) return 'critical';
    if (warningViolations > 3 || summary.avgLoadTime > 1500) return 'poor';
    if (warningViolations > 0 || summary.avgLoadTime > 800) return 'good';
    return 'excellent';
  }

  private logSessionResults(): void {
    const summary = this.session.summary;
    const violations = this.session.budgetViolations;

    console.group(`Performance Session Results - ${this.session.sessionId}`);
    console.log('Summary:', summary);
    console.log('Budget Violations:', violations);
    console.log('Status:', this.getOverallStatus());

    if (violations.length > 0) {
      console.warn(`${violations.length} budget violations detected`);
      violations.forEach(violation => {
        console.warn(`${violation.metric}: ${violation.actual}ms > ${violation.expected}ms (${violation.severity})`);
      });
    }

    console.groupEnd();
  }
}