/**
 * Performance Monitoring and Tracking System
 * 
 * Real-time performance monitoring, alerting, and continuous tracking
 * for content load times and system performance metrics
 */

import { 
  PerformanceMeasurement, 
  ContentLoadMetrics, 
  PerformanceMetrics,
  PERFORMANCE_BUDGET 
} from "./measurement";

/**
 * Performance alert configuration
 */
export interface PerformanceAlert {
  id: string;
  name: string;
  description: string;
  condition: (metrics: ContentLoadMetrics | PerformanceMetrics) => boolean;
  severity: "low" | "medium" | "high" | "critical";
  threshold: number;
  enabled: boolean;
  cooldownMs: number; // Minimum time between alerts
  lastTriggered?: number;
}

/**
 * Performance monitoring configuration
 */
export interface MonitoringConfig {
  enabled: boolean;
  sampleRate: number; // 0-1, percentage of requests to monitor
  batchSize: number; // Number of metrics to batch before sending
  flushIntervalMs: number; // How often to flush batched metrics
  retentionDays: number; // How long to keep historical data
  alerts: PerformanceAlert[];
  endpoints: {
    metrics: string; // Where to send metrics
    alerts: string; // Where to send alerts
  };
}

/**
 * Performance trend analysis
 */
export interface PerformanceTrend {
  metric: string;
  timeframe: "hour" | "day" | "week" | "month";
  direction: "improving" | "degrading" | "stable";
  changePercent: number;
  confidence: number; // 0-1
  significance: "high" | "medium" | "low";
}

/**
 * Performance monitoring dashboard data
 */
export interface MonitoringDashboard {
  overview: {
    totalRequests: number;
    averageLoadTime: number;
    p95LoadTime: number;
    errorRate: number;
    cacheHitRatio: number;
    budgetCompliance: number; // Percentage meeting budget
  };
  trends: PerformanceTrend[];
  alerts: {
    active: number;
    resolved: number;
    critical: number;
  };
  breakdown: {
    byContentType: Record<string, { count: number; avgTime: number }>;
    byCategory: Record<string, { count: number; avgTime: number }>;
    byNetworkCondition: Record<string, { count: number; avgTime: number }>;
    byCacheStatus: Record<string, { count: number; avgTime: number }>;
  };
  recentMetrics: ContentLoadMetrics[];
}

/**
 * Performance monitoring system
 */
export class PerformanceMonitor {
  private config: MonitoringConfig;
  private measurement: PerformanceMeasurement;
  private metricsBuffer: ContentLoadMetrics[] = [];
  private historicalData: Map<string, ContentLoadMetrics[]> = new Map();
  private activeAlerts = new Set<string>();
  private flushTimer?: NodeJS.Timeout;
  private onAlert?: (alert: PerformanceAlert, metrics: ContentLoadMetrics) => void;
  private onDashboardUpdate?: (dashboard: MonitoringDashboard) => void;

  constructor(config: MonitoringConfig) {
    this.config = config;
    this.measurement = new PerformanceMeasurement((metrics) => {
      this.processMetrics(metrics);
    });
    
    this.startPeriodicFlush();
  }

  /**
   * Set alert callback
   */
  onAlertTriggered(callback: (alert: PerformanceAlert, metrics: ContentLoadMetrics) => void): void {
    this.onAlert = callback;
  }

  /**
   * Set dashboard update callback
   */
  onDashboardUpdated(callback: (dashboard: MonitoringDashboard) => void): void {
    this.onDashboardUpdate = callback;
  }

  /**
   * Start monitoring a content load operation
   */
  startMonitoring(contentId: string, contentType: string, category: string): void {
    if (!this.config.enabled || Math.random() > this.config.sampleRate) {
      return; // Skip monitoring based on sample rate
    }

    this.measurement.startContentLoad(contentId, contentType, category);
  }

  /**
   * End monitoring a content load operation
   */
  endMonitoring(
    contentId: string, 
    size?: number, 
    cacheStatus?: "hit" | "miss" | "stale"
  ): ContentLoadMetrics | null {
    if (!this.config.enabled) {
      return null;
    }

    return this.measurement.endContentLoad(contentId, size, cacheStatus);
  }

  /**
   * Process collected metrics
   */
  private processMetrics(metrics: ContentLoadMetrics): void {
    // Add to buffer
    this.metricsBuffer.push(metrics);
    
    // Store in historical data
    const dateKey = new Date().toISOString().split("T")[0]; // YYYY-MM-DD
    if (!this.historicalData.has(dateKey)) {
      this.historicalData.set(dateKey, []);
    }
    this.historicalData.get(dateKey)!.push(metrics);

    // Check alerts
    this.checkAlerts(metrics);

    // Flush if buffer is full
    if (this.metricsBuffer.length >= this.config.batchSize) {
      this.flushMetrics();
    }

    // Update dashboard
    this.updateDashboard();
  }

  /**
   * Check performance alerts
   */
  private checkAlerts(metrics: ContentLoadMetrics): void {
    this.config.alerts.forEach(alert => {
      if (!alert.enabled) return;

      // Check cooldown
      if (alert.lastTriggered && 
          (Date.now() - alert.lastTriggered) < alert.cooldownMs) {
        return;
      }

      // Check condition
      if (alert.condition(metrics)) {
        this.triggerAlert(alert, metrics);
      }
    });
  }

  /**
   * Trigger a performance alert
   */
  private triggerAlert(alert: PerformanceAlert, metrics: ContentLoadMetrics): void {
    alert.lastTriggered = Date.now();
    this.activeAlerts.add(alert.id);

    console.warn(`Performance Alert: ${alert.name}`, {
      alert,
      metrics,
      severity: alert.severity,
    });

    // Send alert to external system
    this.sendAlert(alert, metrics);

    // Notify callback
    if (this.onAlert) {
      this.onAlert(alert, metrics);
    }
  }

  /**
   * Send alert to external system
   */
  private async sendAlert(alert: PerformanceAlert, metrics: ContentLoadMetrics): Promise<void> {
    if (!this.config.endpoints.alerts) return;

    try {
      const alertData = {
        alertId: alert.id,
        name: alert.name,
        severity: alert.severity,
        timestamp: Date.now(),
        metrics,
        threshold: alert.threshold,
      };

      // In a real implementation, this would send to an external alerting service
      console.log("Sending alert to:", this.config.endpoints.alerts, alertData);
      
      // Mock API call
      // await fetch(this.config.endpoints.alerts, {
      //   method: "POST",
      //   headers: { "Content-Type": "application/json" },
      //   body: JSON.stringify(alertData),
      // });
    } catch (error) {
      console.error("Failed to send alert:", error);
    }
  }

  /**
   * Start periodic metrics flushing
   */
  private startPeriodicFlush(): void {
    this.flushTimer = setInterval(() => {
      this.flushMetrics();
    }, this.config.flushIntervalMs);
  }

  /**
   * Flush metrics buffer to external system
   */
  private async flushMetrics(): Promise<void> {
    if (this.metricsBuffer.length === 0 || !this.config.endpoints.metrics) {
      return;
    }

    const metricsToSend = [...this.metricsBuffer];
    this.metricsBuffer = [];

    try {
      const payload = {
        timestamp: Date.now(),
        metrics: metricsToSend,
        metadata: {
          sampleRate: this.config.sampleRate,
          budgetThreshold: PERFORMANCE_BUDGET.maxLoadTime,
        },
      };

      console.log("Flushing metrics to:", this.config.endpoints.metrics, {
        count: metricsToSend.length,
        timeRange: {
          start: Math.min(...metricsToSend.map(m => m.timestamp)),
          end: Math.max(...metricsToSend.map(m => m.timestamp)),
        },
      });

      // In a real implementation, this would send to a metrics collection service
      // await fetch(this.config.endpoints.metrics, {
      //   method: "POST",
      //   headers: { "Content-Type": "application/json" },
      //   body: JSON.stringify(payload),
      // });
    } catch (error) {
      console.error("Failed to flush metrics:", error);
      // Re-add metrics to buffer for retry
      this.metricsBuffer.unshift(...metricsToSend);
    }
  }

  /**
   * Update monitoring dashboard
   */
  private updateDashboard(): void {
    const dashboard = this.generateDashboard();
    
    if (this.onDashboardUpdate) {
      this.onDashboardUpdate(dashboard);
    }
  }

  /**
   * Generate dashboard data
   */
  generateDashboard(): MonitoringDashboard {
    const allMetrics = this.getAllRecentMetrics();
    
    if (allMetrics.length === 0) {
      return this.getEmptyDashboard();
    }

    const overview = this.calculateOverview(allMetrics);
    const trends = this.calculateTrends();
    const breakdown = this.calculateBreakdown(allMetrics);
    
    return {
      overview,
      trends,
      alerts: {
        active: this.activeAlerts.size,
        resolved: 0, // Would be calculated from historical alert data
        critical: this.config.alerts.filter(a => a.severity === "critical").length,
      },
      breakdown,
      recentMetrics: allMetrics.slice(-20), // Last 20 metrics
    };
  }

  /**
   * Get all recent metrics from historical data
   */
  private getAllRecentMetrics(): ContentLoadMetrics[] {
    const recent: ContentLoadMetrics[] = [];
    const cutoffTime = Date.now() - (24 * 60 * 60 * 1000); // Last 24 hours

    this.historicalData.forEach(metrics => {
      const recentMetrics = metrics.filter(m => m.timestamp >= cutoffTime);
      recent.push(...recentMetrics);
    });

    return recent.sort((a, b) => a.timestamp - b.timestamp);
  }

  /**
   * Calculate overview metrics
   */
  private calculateOverview(metrics: ContentLoadMetrics[]) {
    const loadTimes = metrics.map(m => m.loadTime);
    const errors = metrics.filter(m => m.loadTime > PERFORMANCE_BUDGET.maxLoadTime * 2);
    const cacheHits = metrics.filter(m => m.cacheStatus === "hit");
    const budgetCompliant = metrics.filter(m => m.loadTime <= PERFORMANCE_BUDGET.maxLoadTime);

    return {
      totalRequests: metrics.length,
      averageLoadTime: loadTimes.reduce((a, b) => a + b, 0) / loadTimes.length || 0,
      p95LoadTime: this.calculatePercentile(loadTimes, 95),
      errorRate: errors.length / metrics.length || 0,
      cacheHitRatio: cacheHits.length / metrics.length || 0,
      budgetCompliance: budgetCompliant.length / metrics.length || 0,
    };
  }

  /**
   * Calculate performance trends
   */
  private calculateTrends(): PerformanceTrend[] {
    // Simplified trend calculation
    // In a real implementation, this would analyze historical data over time
    const trends: PerformanceTrend[] = [
      {
        metric: "loadTime",
        timeframe: "hour",
        direction: "stable",
        changePercent: 0,
        confidence: 0.8,
        significance: "medium",
      },
      {
        metric: "cacheHitRatio",
        timeframe: "day",
        direction: "improving",
        changePercent: 5.2,
        confidence: 0.9,
        significance: "high",
      },
    ];

    return trends;
  }

  /**
   * Calculate breakdown metrics
   */
  private calculateBreakdown(metrics: ContentLoadMetrics[]) {
    const byContentType: Record<string, { count: number; avgTime: number }> = {};
    const byCategory: Record<string, { count: number; avgTime: number }> = {};
    const byNetworkCondition: Record<string, { count: number; avgTime: number }> = {};
    const byCacheStatus: Record<string, { count: number; avgTime: number }> = {};

    metrics.forEach(metric => {
      // By content type
      if (!byContentType[metric.contentType]) {
        byContentType[metric.contentType] = { count: 0, avgTime: 0 };
      }
      byContentType[metric.contentType].count++;
      byContentType[metric.contentType].avgTime = 
        (byContentType[metric.contentType].avgTime * (byContentType[metric.contentType].count - 1) + metric.loadTime) / 
        byContentType[metric.contentType].count;

      // By category
      if (!byCategory[metric.category]) {
        byCategory[metric.category] = { count: 0, avgTime: 0 };
      }
      byCategory[metric.category].count++;
      byCategory[metric.category].avgTime = 
        (byCategory[metric.category].avgTime * (byCategory[metric.category].count - 1) + metric.loadTime) / 
        byCategory[metric.category].count;

      // By network condition
      if (!byNetworkCondition[metric.networkCondition]) {
        byNetworkCondition[metric.networkCondition] = { count: 0, avgTime: 0 };
      }
      byNetworkCondition[metric.networkCondition].count++;
      byNetworkCondition[metric.networkCondition].avgTime = 
        (byNetworkCondition[metric.networkCondition].avgTime * (byNetworkCondition[metric.networkCondition].count - 1) + metric.loadTime) / 
        byNetworkCondition[metric.networkCondition].count;

      // By cache status
      if (!byCacheStatus[metric.cacheStatus]) {
        byCacheStatus[metric.cacheStatus] = { count: 0, avgTime: 0 };
      }
      byCacheStatus[metric.cacheStatus].count++;
      byCacheStatus[metric.cacheStatus].avgTime = 
        (byCacheStatus[metric.cacheStatus].avgTime * (byCacheStatus[metric.cacheStatus].count - 1) + metric.loadTime) / 
        byCacheStatus[metric.cacheStatus].count;
    });

    return {
      byContentType,
      byCategory,
      byNetworkCondition,
      byCacheStatus,
    };
  }

  /**
   * Calculate percentile
   */
  private calculatePercentile(values: number[], percentile: number): number {
    if (values.length === 0) return 0;
    
    const sorted = [...values].sort((a, b) => a - b);
    const index = Math.floor((percentile / 100) * sorted.length);
    return sorted[Math.min(index, sorted.length - 1)];
  }

  /**
   * Get empty dashboard
   */
  private getEmptyDashboard(): MonitoringDashboard {
    return {
      overview: {
        totalRequests: 0,
        averageLoadTime: 0,
        p95LoadTime: 0,
        errorRate: 0,
        cacheHitRatio: 0,
        budgetCompliance: 0,
      },
      trends: [],
      alerts: {
        active: 0,
        resolved: 0,
        critical: 0,
      },
      breakdown: {
        byContentType: {},
        byCategory: {},
        byNetworkCondition: {},
        byCacheStatus: {},
      },
      recentMetrics: [],
    };
  }

  /**
   * Add custom alert
   */
  addAlert(alert: PerformanceAlert): void {
    this.config.alerts.push(alert);
  }

  /**
   * Remove alert
   */
  removeAlert(alertId: string): void {
    this.config.alerts = this.config.alerts.filter(a => a.id !== alertId);
    this.activeAlerts.delete(alertId);
  }

  /**
   * Update monitoring configuration
   */
  updateConfig(newConfig: Partial<MonitoringConfig>): void {
    this.config = { ...this.config, ...newConfig };
  }

  /**
   * Get current configuration
   */
  getConfig(): MonitoringConfig {
    return { ...this.config };
  }

  /**
   * Cleanup and stop monitoring
   */
  stop(): void {
    if (this.flushTimer) {
      clearInterval(this.flushTimer);
    }
    
    // Flush remaining metrics
    this.flushMetrics();
    
    this.measurement.disconnect();
  }
}

/**
 * Default performance alerts
 */
export const DEFAULT_ALERTS: PerformanceAlert[] = [
  {
    id: "load-time-critical",
    name: "Critical Load Time",
    description: "Content load time exceeds 500ms (2.5x budget)",
    condition: (metrics) => (metrics as ContentLoadMetrics).loadTime > 500,
    severity: "critical",
    threshold: 500,
    enabled: true,
    cooldownMs: 60000, // 1 minute
  },
  {
    id: "load-time-warning",
    name: "Load Time Warning",
    description: "Content load time exceeds 200ms budget",
    condition: (metrics) => (metrics as ContentLoadMetrics).loadTime > PERFORMANCE_BUDGET.maxLoadTime,
    severity: "medium",
    threshold: PERFORMANCE_BUDGET.maxLoadTime,
    enabled: true,
    cooldownMs: 300000, // 5 minutes
  },
  {
    id: "cache-miss-rate",
    name: "High Cache Miss Rate",
    description: "Cache miss rate exceeds 50%",
    condition: (metrics) => (metrics as ContentLoadMetrics).cacheStatus === "miss",
    severity: "medium",
    threshold: 0.5,
    enabled: true,
    cooldownMs: 600000, // 10 minutes
  },
  {
    id: "error-rate-high",
    name: "High Error Rate",
    description: "Error rate exceeds 5%",
    condition: (metrics) => (metrics as ContentLoadMetrics).loadTime > 1000, // Consider >1s as error
    severity: "high",
    threshold: 0.05,
    enabled: true,
    cooldownMs: 180000, // 3 minutes
  },
];

/**
 * Default monitoring configuration
 */
export const DEFAULT_MONITORING_CONFIG: MonitoringConfig = {
  enabled: true,
  sampleRate: 0.1, // Monitor 10% of requests
  batchSize: 50,
  flushIntervalMs: 30000, // 30 seconds
  retentionDays: 30,
  alerts: DEFAULT_ALERTS,
  endpoints: {
    metrics: process.env.PERFORMANCE_METRICS_ENDPOINT || "",
    alerts: process.env.PERFORMANCE_ALERTS_ENDPOINT || "",
  },
};

/**
 * Create and configure performance monitor
 */
export function createPerformanceMonitor(
  config: Partial<MonitoringConfig> = {}
): PerformanceMonitor {
  const finalConfig = { ...DEFAULT_MONITORING_CONFIG, ...config };
  return new PerformanceMonitor(finalConfig);
}

/**
 * React hook for performance monitoring
 */
export function usePerformanceMonitoring(config?: Partial<MonitoringConfig>) {
  const [monitor] = React.useState(() => createPerformanceMonitor(config));
  const [dashboard, setDashboard] = React.useState<MonitoringDashboard | null>(null);
  const [alerts, setAlerts] = React.useState<Array<{ alert: PerformanceAlert; metrics: ContentLoadMetrics }>>([]);

  React.useEffect(() => {
    monitor.onDashboardUpdated(setDashboard);
    monitor.onAlertTriggered((alert, metrics) => {
      setAlerts(prev => [...prev, { alert, metrics }]);
    });

    return () => {
      monitor.stop();
    };
  }, [monitor]);

  const startMonitoring = React.useCallback((contentId: string, contentType: string, category: string) => {
    monitor.startMonitoring(contentId, contentType, category);
  }, [monitor]);

  const endMonitoring = React.useCallback((contentId: string, size?: number, cacheStatus?: "hit" | "miss" | "stale") => {
    return monitor.endMonitoring(contentId, size, cacheStatus);
  }, [monitor]);

  const clearAlerts = React.useCallback(() => {
    setAlerts([]);
  }, []);

  return {
    monitor,
    dashboard,
    alerts,
    startMonitoring,
    endMonitoring,
    clearAlerts,
  };
}
