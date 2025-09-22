import React from "react";
/**
 * Performance Measurement Utilities
 * 
 * Browser Performance API wrapper for content load time validation
 * Validates <200ms content load time requirement
 */

export interface PerformanceMetrics {
  loadTime: number;
  firstContentfulPaint: number;
  largestContentfulPaint: number;
  firstInputDelay: number;
  cumulativeLayoutShift: number;
  timeToInteractive: number;
  resourceLoadTime: number;
  cacheHitRatio: number;
  networkLatency: number;
  contentSize: number;
}

export interface ContentLoadMetrics {
  contentId: string;
  contentType: string;
  category: string;
  loadTime: number;
  cacheStatus: "hit" | "miss" | "stale";
  networkCondition: string;
  timestamp: number;
  size: number;
  renderTime: number;
}

export interface PerformanceBudget {
  maxLoadTime: number; // 200ms requirement
  maxFirstContentfulPaint: number;
  maxLargestContentfulPaint: number;
  maxFirstInputDelay: number;
  maxCumulativeLayoutShift: number;
  maxResourceLoadTime: number;
  minCacheHitRatio: number;
}

export const PERFORMANCE_BUDGET: PerformanceBudget = {
  maxLoadTime: 200, // Primary requirement
  maxFirstContentfulPaint: 1000,
  maxLargestContentfulPaint: 2500,
  maxFirstInputDelay: 100,
  maxCumulativeLayoutShift: 0.1,
  maxResourceLoadTime: 150,
  minCacheHitRatio: 0.8,
};

/**
 * Performance measurement utility class
 */
export class PerformanceMeasurement {
  private observer: PerformanceObserver | null = null;
  private measurements: Map<string, ContentLoadMetrics> = new Map();
  private onMetricsCallback?: (metrics: ContentLoadMetrics) => void;

  constructor(onMetrics?: (metrics: ContentLoadMetrics) => void) {
    this.onMetricsCallback = onMetrics;
    this.initializeObserver();
  }

  /**
   * Initialize Performance Observer for monitoring
   */
  private initializeObserver(): void {
    if (typeof window === "undefined" || !("PerformanceObserver" in window)) {
      console.warn("PerformanceObserver not supported");
      return;
    }

    this.observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        this.processPerformanceEntry(entry);
      });
    });

    // Observe all performance entry types
    try {
      this.observer.observe({ 
        entryTypes: ["measure", "navigation", "resource", "paint", "layout-shift", "first-input"] 
      });
    } catch (error) {
      console.warn("Performance observation failed:", error);
    }
  }

  /**
   * Process performance entries
   */
  private processPerformanceEntry(entry: PerformanceEntry): void {
    if (entry.name.includes("content-load")) {
      const contentId = this.extractContentId(entry.name);
      if (contentId) {
        this.updateContentMetrics(contentId, entry);
      }
    }
  }

  /**
   * Extract content ID from performance mark name
   */
  private extractContentId(markName: string): string | null {
    const match = markName.match(/content-load-(.+)-(start|end)/);
    return match ? match[1] : null;
  }

  /**
   * Update content metrics
   */
  private updateContentMetrics(contentId: string, entry: PerformanceEntry): void {
    const existing = this.measurements.get(contentId);
    if (existing) {
      existing.loadTime = entry.duration || 0;
      existing.timestamp = Date.now();
    }
  }

  /**
   * Start measuring content load time
   */
  startContentLoad(contentId: string, contentType: string, category: string): void {
    const markName = `content-load-${contentId}-start`;
    
    try {
      performance.mark(markName);
      
      const metrics: ContentLoadMetrics = {
        contentId,
        contentType,
        category,
        loadTime: 0,
        cacheStatus: "miss",
        networkCondition: this.getNetworkCondition(),
        timestamp: Date.now(),
        size: 0,
        renderTime: 0,
      };
      
      this.measurements.set(contentId, metrics);
    } catch (error) {
      console.warn("Failed to start performance measurement:", error);
    }
  }

  /**
   * End measuring content load time
   */
  endContentLoad(contentId: string, size?: number, cacheStatus?: "hit" | "miss" | "stale"): ContentLoadMetrics | null {
    const startMark = `content-load-${contentId}-start`;
    const endMark = `content-load-${contentId}-end`;
    const measureName = `content-load-${contentId}`;

    try {
      performance.mark(endMark);
      performance.measure(measureName, startMark, endMark);

      const measure = performance.getEntriesByName(measureName, "measure")[0];
      const metrics = this.measurements.get(contentId);

      if (metrics && measure) {
        metrics.loadTime = measure.duration;
        metrics.size = size || 0;
        metrics.cacheStatus = cacheStatus || "miss";
        metrics.renderTime = performance.now();

        // Clean up performance marks
        performance.clearMarks(startMark);
        performance.clearMarks(endMark);
        performance.clearMeasures(measureName);

        // Notify callback
        if (this.onMetricsCallback) {
          this.onMetricsCallback(metrics);
        }

        return metrics;
      }
    } catch (error) {
      console.warn("Failed to end performance measurement:", error);
    }

    return null;
  }

  /**
   * Get current network condition
   */
  private getNetworkCondition(): string {
    if ("connection" in navigator) {
      const connection = (navigator as any).connection;
      return connection.effectiveType || "unknown";
    }
    return "unknown";
  }

  /**
   * Get all measurements
   */
  getAllMeasurements(): ContentLoadMetrics[] {
    return Array.from(this.measurements.values());
  }

  /**
   * Get measurements by content type
   */
  getMeasurementsByType(contentType: string): ContentLoadMetrics[] {
    return this.getAllMeasurements().filter(m => m.contentType === contentType);
  }

  /**
   * Get measurements by category
   */
  getMeasurementsByCategory(category: string): ContentLoadMetrics[] {
    return this.getAllMeasurements().filter(m => m.category === category);
  }

  /**
   * Check if metrics meet performance budget
   */
  validateMetrics(metrics: ContentLoadMetrics): boolean {
    return metrics.loadTime <= PERFORMANCE_BUDGET.maxLoadTime;
  }

  /**
   * Get performance summary
   */
  getPerformanceSummary(): {
    totalMeasurements: number;
    averageLoadTime: number;
    maxLoadTime: number;
    minLoadTime: number;
    passRate: number;
    failedCount: number;
  } {
    const measurements = this.getAllMeasurements();
    
    if (measurements.length === 0) {
      return {
        totalMeasurements: 0,
        averageLoadTime: 0,
        maxLoadTime: 0,
        minLoadTime: 0,
        passRate: 0,
        failedCount: 0,
      };
    }

    const loadTimes = measurements.map(m => m.loadTime);
    const passedCount = measurements.filter(m => this.validateMetrics(m)).length;

    return {
      totalMeasurements: measurements.length,
      averageLoadTime: loadTimes.reduce((a, b) => a + b, 0) / loadTimes.length,
      maxLoadTime: Math.max(...loadTimes),
      minLoadTime: Math.min(...loadTimes),
      passRate: passedCount / measurements.length,
      failedCount: measurements.length - passedCount,
    };
  }

  /**
   * Clear all measurements
   */
  clearMeasurements(): void {
    this.measurements.clear();
  }

  /**
   * Stop performance observation
   */
  disconnect(): void {
    if (this.observer) {
      this.observer.disconnect();
      this.observer = null;
    }
  }
}

/**
 * Get Core Web Vitals metrics
 */
export function getCoreWebVitals(): Promise<PerformanceMetrics> {
  return new Promise((resolve) => {
    const metrics: Partial<PerformanceMetrics> = {};

    // Get navigation timing
    const navigation = performance.getEntriesByType("navigation")[0] as PerformanceNavigationTiming;
    if (navigation) {
      metrics.loadTime = navigation.loadEventEnd - navigation.loadEventStart;
      metrics.timeToInteractive = navigation.domInteractive - navigation.navigationStart;
      metrics.networkLatency = navigation.responseStart - navigation.requestStart;
    }

    // Get paint metrics
    const paintEntries = performance.getEntriesByType("paint");
    paintEntries.forEach((entry) => {
      if (entry.name === "first-contentful-paint") {
        metrics.firstContentfulPaint = entry.startTime;
      }
    });

    // Get resource metrics
    const resourceEntries = performance.getEntriesByType("resource");
    if (resourceEntries.length > 0) {
      const totalResourceTime = resourceEntries.reduce((sum, entry) => sum + entry.duration, 0);
      metrics.resourceLoadTime = totalResourceTime / resourceEntries.length;
    }

    // Set defaults for missing metrics
    const completeMetrics: PerformanceMetrics = {
      loadTime: metrics.loadTime || 0,
      firstContentfulPaint: metrics.firstContentfulPaint || 0,
      largestContentfulPaint: 0, // Will be updated by observer
      firstInputDelay: 0, // Will be updated by observer
      cumulativeLayoutShift: 0, // Will be updated by observer
      timeToInteractive: metrics.timeToInteractive || 0,
      resourceLoadTime: metrics.resourceLoadTime || 0,
      cacheHitRatio: 0.8, // Default assumption
      networkLatency: metrics.networkLatency || 0,
      contentSize: 0,
    };

    resolve(completeMetrics);
  });
}

/**
 * Measure specific function execution time
 */
export function measureExecutionTime<T>(
  fn: () => Promise<T> | T,
  label: string
): Promise<{ result: T; duration: number }> {
  return new Promise(async (resolve, reject) => {
    const startTime = performance.now();
    
    try {
      const result = await fn();
      const endTime = performance.now();
      const duration = endTime - startTime;

      // Log if exceeds budget
      if (duration > PERFORMANCE_BUDGET.maxLoadTime) {
        console.warn(`Performance budget exceeded for ${label}: ${duration}ms > ${PERFORMANCE_BUDGET.maxLoadTime}ms`);
      }

      resolve({ result, duration });
    } catch (error) {
      reject(error);
    }
  });
}

/**
 * Create performance HOC for React components
 */
export function withPerformanceMonitoring<P extends object>(
  Component: React.ComponentType<P>,
  componentName: string
): React.ComponentType<P> {
  return function PerformanceMonitoredComponent(props: P) {
    const measurement = new PerformanceMeasurement();
    
    React.useEffect(() => {
      measurement.startContentLoad(componentName, "component", "ui");
      
      return () => {
        measurement.endContentLoad(componentName);
        measurement.disconnect();
      };
    }, []);

    return React.createElement(Component, props);
  };
}

/**
 * Hook for measuring content load performance
 */
export function useContentPerformance(contentId: string, contentType: string, category: string) {
  const [metrics, setMetrics] = React.useState<ContentLoadMetrics | null>(null);
  const measurementRef = React.useRef<PerformanceMeasurement | null>(null);

  React.useEffect(() => {
    measurementRef.current = new PerformanceMeasurement((m) => setMetrics(m));
    return () => {
      measurementRef.current?.disconnect();
    };
  }, []);

  const startMeasurement = React.useCallback(() => {
    measurementRef.current?.startContentLoad(contentId, contentType, category);
  }, [contentId, contentType, category]);

  const endMeasurement = React.useCallback((size?: number, cacheStatus?: "hit" | "miss" | "stale") => {
    const result = measurementRef.current?.endContentLoad(contentId, size, cacheStatus);
    return result;
  }, [contentId]);

  return {
    metrics,
    startMeasurement,
    endMeasurement,
    isWithinBudget: metrics ? metrics.loadTime <= PERFORMANCE_BUDGET.maxLoadTime : null,
  };
}
