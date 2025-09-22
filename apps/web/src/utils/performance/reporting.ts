/**
 * Performance Test Reports and Analysis Tools
 * 
 * Comprehensive reporting and analysis tools for performance validation results
 * Generates detailed reports with recommendations and trend analysis
 */

import { 
  ContentLoadMetrics, 
  PerformanceMetrics, 
  PERFORMANCE_BUDGET 
} from "./measurement";
import { MonitoringDashboard, PerformanceTrend } from "./monitoring";
import { CacheValidationResult } from "./cache-validation.test";

/**
 * Performance test report structure
 */
export interface PerformanceReport {
  id: string;
  title: string;
  generatedAt: string;
  summary: PerformanceReportSummary;
  sections: PerformanceReportSection[];
  recommendations: PerformanceRecommendation[];
  metadata: PerformanceReportMetadata;
}

/**
 * Report summary section
 */
export interface PerformanceReportSummary {
  overallScore: number; // 0-100
  budgetCompliance: number; // Percentage meeting <200ms requirement
  totalTests: number;
  passedTests: number;
  failedTests: number;
  criticalIssues: number;
  status: "pass" | "fail" | "warning";
  keyFindings: string[];
}

/**
 * Individual report section
 */
export interface PerformanceReportSection {
  id: string;
  title: string;
  type: "overview" | "detailed" | "trends" | "breakdown" | "recommendations";
  content: any;
  charts?: PerformanceChart[];
  tables?: PerformanceTable[];
  insights: string[];
}

/**
 * Performance recommendations
 */
export interface PerformanceRecommendation {
  id: string;
  title: string;
  description: string;
  severity: "low" | "medium" | "high" | "critical";
  category: "content" | "network" | "cache" | "system" | "architecture";
  impact: "low" | "medium" | "high";
  effort: "low" | "medium" | "high";
  expectedImprovement: string;
  implementation: string[];
  relatedMetrics: string[];
}

/**
 * Report metadata
 */
export interface PerformanceReportMetadata {
  testEnvironment: {
    userAgent: string;
    networkConditions: string[];
    contentTypes: string[];
    testDuration: number;
  };
  configuration: {
    performanceBudget: typeof PERFORMANCE_BUDGET;
    sampleSize: number;
    testScenarios: string[];
  };
  version: string;
  generatedBy: string;
}

/**
 * Chart configuration for reports
 */
export interface PerformanceChart {
  id: string;
  title: string;
  type: "line" | "bar" | "histogram" | "scatter" | "pie";
  data: ChartDataPoint[];
  xAxis: string;
  yAxis: string;
  threshold?: number; // Performance budget line
}

/**
 * Chart data point
 */
export interface ChartDataPoint {
  x: string | number;
  y: number;
  category?: string;
  metadata?: any;
}

/**
 * Table configuration for reports
 */
export interface PerformanceTable {
  id: string;
  title: string;
  headers: string[];
  rows: TableRow[];
  sortBy?: string;
  highlightCondition?: (row: TableRow) => boolean;
}

/**
 * Table row
 */
export interface TableRow {
  [key: string]: string | number | boolean;
}

/**
 * Performance report generator
 */
export class PerformanceReportGenerator {
  private reportHistory: PerformanceReport[] = [];

  /**
   * Generate comprehensive performance report
   */
  generateReport(
    metrics: ContentLoadMetrics[],
    cacheResults: CacheValidationResult[],
    dashboard: MonitoringDashboard,
    title: string = "Performance Validation Report"
  ): PerformanceReport {
    const reportId = `report-${Date.now()}`;
    
    const summary = this.generateSummary(metrics);
    const sections = this.generateSections(metrics, cacheResults, dashboard);
    const recommendations = this.generateRecommendations(metrics, cacheResults);
    const metadata = this.generateMetadata(metrics);

    const report: PerformanceReport = {
      id: reportId,
      title,
      generatedAt: new Date().toISOString(),
      summary,
      sections,
      recommendations,
      metadata,
    };

    this.reportHistory.push(report);
    return report;
  }

  /**
   * Generate report summary
   */
  private generateSummary(metrics: ContentLoadMetrics[]): PerformanceReportSummary {
    if (metrics.length === 0) {
      return {
        overallScore: 0,
        budgetCompliance: 0,
        totalTests: 0,
        passedTests: 0,
        failedTests: 0,
        criticalIssues: 0,
        status: "fail",
        keyFindings: ["No metrics available for analysis"],
      };
    }

    const passedTests = metrics.filter(m => m.loadTime <= PERFORMANCE_BUDGET.maxLoadTime).length;
    const failedTests = metrics.length - passedTests;
    const criticalIssues = metrics.filter(m => m.loadTime > PERFORMANCE_BUDGET.maxLoadTime * 2).length;
    const budgetCompliance = (passedTests / metrics.length) * 100;
    
    // Calculate overall score based on multiple factors
    const overallScore = this.calculateOverallScore(metrics);
    
    const status = overallScore >= 80 ? "pass" : overallScore >= 60 ? "warning" : "fail";
    
    const keyFindings = this.generateKeyFindings(metrics);

    return {
      overallScore,
      budgetCompliance,
      totalTests: metrics.length,
      passedTests,
      failedTests,
      criticalIssues,
      status,
      keyFindings,
    };
  }

  /**
   * Calculate overall performance score
   */
  private calculateOverallScore(metrics: ContentLoadMetrics[]): number {
    if (metrics.length === 0) return 0;

    const budgetScore = (metrics.filter(m => m.loadTime <= PERFORMANCE_BUDGET.maxLoadTime).length / metrics.length) * 40;
    
    const avgLoadTime = metrics.reduce((sum, m) => sum + m.loadTime, 0) / metrics.length;
    const speedScore = Math.max(0, 40 - (avgLoadTime / PERFORMANCE_BUDGET.maxLoadTime) * 20);
    
    const cacheHitRatio = metrics.filter(m => m.cacheStatus === "hit").length / metrics.length;
    const cacheScore = cacheHitRatio * 20;

    return Math.round(budgetScore + speedScore + cacheScore);
  }

  /**
   * Generate key findings
   */
  private generateKeyFindings(metrics: ContentLoadMetrics[]): string[] {
    const findings: string[] = [];
    
    const avgLoadTime = metrics.reduce((sum, m) => sum + m.loadTime, 0) / metrics.length;
    const budgetCompliance = (metrics.filter(m => m.loadTime <= PERFORMANCE_BUDGET.maxLoadTime).length / metrics.length) * 100;
    const cacheHitRatio = (metrics.filter(m => m.cacheStatus === "hit").length / metrics.length) * 100;
    
    findings.push(`Average load time: ${avgLoadTime.toFixed(1)}ms (Budget: ${PERFORMANCE_BUDGET.maxLoadTime}ms)`);
    findings.push(`Budget compliance: ${budgetCompliance.toFixed(1)}% of requests`);
    findings.push(`Cache hit ratio: ${cacheHitRatio.toFixed(1)}%`);
    
    if (budgetCompliance < 80) {
      findings.push(`⚠️ Low budget compliance - ${(100 - budgetCompliance).toFixed(1)}% of requests exceed 200ms`);
    }
    
    if (cacheHitRatio < 70) {
      findings.push(`⚠️ Low cache efficiency - Only ${cacheHitRatio.toFixed(1)}% cache hits`);
    }
    
    const slowestMetric = metrics.reduce((slowest, current) => 
      current.loadTime > slowest.loadTime ? current : slowest
    );
    findings.push(`Slowest request: ${slowestMetric.loadTime.toFixed(1)}ms (${slowestMetric.contentType})`);

    return findings;
  }

  /**
   * Generate report sections
   */
  private generateSections(
    metrics: ContentLoadMetrics[],
    cacheResults: CacheValidationResult[],
    dashboard: MonitoringDashboard
  ): PerformanceReportSection[] {
    return [
      this.generateOverviewSection(metrics),
      this.generateDetailedAnalysisSection(metrics),
      this.generateCacheAnalysisSection(cacheResults),
      this.generateTrendsSection(dashboard.trends),
      this.generateBreakdownSection(dashboard.breakdown),
      this.generateNetworkAnalysisSection(metrics),
    ];
  }

  /**
   * Generate overview section
   */
  private generateOverviewSection(metrics: ContentLoadMetrics[]): PerformanceReportSection {
    const loadTimes = metrics.map(m => m.loadTime);
    const p50 = this.calculatePercentile(loadTimes, 50);
    const p95 = this.calculatePercentile(loadTimes, 95);
    const p99 = this.calculatePercentile(loadTimes, 99);

    const chart: PerformanceChart = {
      id: "load-time-distribution",
      title: "Load Time Distribution",
      type: "histogram",
      data: this.createHistogramData(loadTimes),
      xAxis: "Load Time (ms)",
      yAxis: "Frequency",
      threshold: PERFORMANCE_BUDGET.maxLoadTime,
    };

    const table: PerformanceTable = {
      id: "performance-percentiles",
      title: "Performance Percentiles",
      headers: ["Metric", "Value", "Budget", "Status"],
      rows: [
        { metric: "P50", value: `${p50.toFixed(1)}ms`, budget: `${PERFORMANCE_BUDGET.maxLoadTime}ms`, status: p50 <= PERFORMANCE_BUDGET.maxLoadTime ? "✅" : "❌" },
        { metric: "P95", value: `${p95.toFixed(1)}ms`, budget: `${PERFORMANCE_BUDGET.maxLoadTime}ms`, status: p95 <= PERFORMANCE_BUDGET.maxLoadTime ? "✅" : "❌" },
        { metric: "P99", value: `${p99.toFixed(1)}ms`, budget: `${PERFORMANCE_BUDGET.maxLoadTime}ms`, status: p99 <= PERFORMANCE_BUDGET.maxLoadTime ? "✅" : "❌" },
      ],
    };

    const insights = [
      `${metrics.length} total content load operations measured`,
      `P95 load time: ${p95.toFixed(1)}ms ${p95 <= PERFORMANCE_BUDGET.maxLoadTime ? "(within budget)" : "(exceeds budget)"}`,
      `${loadTimes.filter(t => t <= PERFORMANCE_BUDGET.maxLoadTime).length} requests met the 200ms budget`,
    ];

    return {
      id: "overview",
      title: "Performance Overview",
      type: "overview",
      content: {
        totalRequests: metrics.length,
        averageLoadTime: loadTimes.reduce((a, b) => a + b, 0) / loadTimes.length,
        percentiles: { p50, p95, p99 },
      },
      charts: [chart],
      tables: [table],
      insights,
    };
  }

  /**
   * Generate detailed analysis section
   */
  private generateDetailedAnalysisSection(metrics: ContentLoadMetrics[]): PerformanceReportSection {
    const byContentType = this.groupBy(metrics, "contentType");
    
    const chart: PerformanceChart = {
      id: "content-type-performance",
      title: "Performance by Content Type",
      type: "bar",
      data: Object.entries(byContentType).map(([type, typeMetrics]) => ({
        x: type,
        y: typeMetrics.reduce((sum, m) => sum + m.loadTime, 0) / typeMetrics.length,
        category: type,
      })),
      xAxis: "Content Type",
      yAxis: "Average Load Time (ms)",
      threshold: PERFORMANCE_BUDGET.maxLoadTime,
    };

    const table: PerformanceTable = {
      id: "content-type-breakdown",
      title: "Content Type Performance Breakdown",
      headers: ["Content Type", "Count", "Avg Time", "P95 Time", "Budget Compliance"],
      rows: Object.entries(byContentType).map(([type, typeMetrics]) => {
        const loadTimes = typeMetrics.map(m => m.loadTime);
        const avgTime = loadTimes.reduce((a, b) => a + b, 0) / loadTimes.length;
        const p95Time = this.calculatePercentile(loadTimes, 95);
        const compliance = (typeMetrics.filter(m => m.loadTime <= PERFORMANCE_BUDGET.maxLoadTime).length / typeMetrics.length) * 100;
        
        return {
          contentType: type,
          count: typeMetrics.length,
          avgTime: `${avgTime.toFixed(1)}ms`,
          p95Time: `${p95Time.toFixed(1)}ms`,
          budgetCompliance: `${compliance.toFixed(1)}%`,
        };
      }),
      highlightCondition: (row) => parseFloat(row.budgetCompliance as string) < 80,
    };

    const insights = [
      "Content type performance analysis shows variation across different content types",
      Object.entries(byContentType).length > 0 ? 
        `Best performing: ${Object.entries(byContentType).reduce((best, [type, typeMetrics]) => {
          const avgTime = typeMetrics.reduce((sum, m) => sum + m.loadTime, 0) / typeMetrics.length;
          return !best.avgTime || avgTime < best.avgTime ? { type, avgTime } : best;
        }, { type: "", avgTime: 0 }).type}` : "No data available",
    ];

    return {
      id: "detailed-analysis",
      title: "Detailed Performance Analysis",
      type: "detailed",
      content: { byContentType },
      charts: [chart],
      tables: [table],
      insights,
    };
  }

  /**
   * Generate cache analysis section
   */
  private generateCacheAnalysisSection(cacheResults: CacheValidationResult[]): PerformanceReportSection {
    if (cacheResults.length === 0) {
      return {
        id: "cache-analysis",
        title: "Cache Performance Analysis",
        type: "detailed",
        content: {},
        insights: ["No cache validation results available"],
      };
    }

    const chart: PerformanceChart = {
      id: "cache-hit-ratios",
      title: "Cache Hit Ratios by Strategy",
      type: "bar",
      data: cacheResults.map(result => ({
        x: result.strategy,
        y: result.stats.hitRatio * 100,
        category: result.strategy,
      })),
      xAxis: "Cache Strategy",
      yAxis: "Hit Ratio (%)",
      threshold: PERFORMANCE_BUDGET.minCacheHitRatio * 100,
    };

    const table: PerformanceTable = {
      id: "cache-performance-table",
      title: "Cache Strategy Performance",
      headers: ["Strategy", "Hit Ratio", "Avg Hit Time", "Avg Miss Time", "Meets Budget"],
      rows: cacheResults.map(result => ({
        strategy: result.strategy,
        hitRatio: `${(result.stats.hitRatio * 100).toFixed(1)}%`,
        avgHitTime: `${result.stats.averageHitTime.toFixed(1)}ms`,
        avgMissTime: `${result.stats.averageMissTime.toFixed(1)}ms`,
        meetsBudget: result.meetsBudget ? "✅" : "❌",
      })),
      highlightCondition: (row) => row.meetsBudget === "❌",
    };

    const insights = cacheResults.map(result => 
      `${result.strategy}: ${(result.stats.hitRatio * 100).toFixed(1)}% hit ratio, ${result.recommendations.length} recommendations`
    );

    return {
      id: "cache-analysis",
      title: "Cache Performance Analysis",
      type: "detailed",
      content: { cacheResults },
      charts: [chart],
      tables: [table],
      insights,
    };
  }

  /**
   * Generate trends section
   */
  private generateTrendsSection(trends: PerformanceTrend[]): PerformanceReportSection {
    const insights = trends.map(trend => 
      `${trend.metric} is ${trend.direction} by ${trend.changePercent.toFixed(1)}% over the last ${trend.timeframe}`
    );

    return {
      id: "trends",
      title: "Performance Trends",
      type: "trends",
      content: { trends },
      insights: insights.length > 0 ? insights : ["No trend data available"],
    };
  }

  /**
   * Generate breakdown section
   */
  private generateBreakdownSection(breakdown: any): PerformanceReportSection {
    const categoryChart: PerformanceChart = {
      id: "category-breakdown",
      title: "Performance by Category",
      type: "pie",
      data: Object.entries(breakdown.byCategory).map(([category, data]: [string, any]) => ({
        x: category,
        y: data.count,
        category,
      })),
      xAxis: "Category",
      yAxis: "Count",
    };

    const insights = [
      `Content distributed across ${Object.keys(breakdown.byCategory).length} categories`,
      `Network conditions tested: ${Object.keys(breakdown.byNetworkCondition).join(", ")}`,
      `Cache status distribution: ${Object.entries(breakdown.byCacheStatus).map(([status, data]: [string, any]) => `${status}: ${data.count}`).join(", ")}`,
    ];

    return {
      id: "breakdown",
      title: "Performance Breakdown",
      type: "breakdown",
      content: { breakdown },
      charts: [categoryChart],
      insights,
    };
  }

  /**
   * Generate network analysis section
   */
  private generateNetworkAnalysisSection(metrics: ContentLoadMetrics[]): PerformanceReportSection {
    const byNetwork = this.groupBy(metrics, "networkCondition");
    
    const chart: PerformanceChart = {
      id: "network-performance",
      title: "Performance by Network Condition",
      type: "bar",
      data: Object.entries(byNetwork).map(([condition, conditionMetrics]) => ({
        x: condition,
        y: conditionMetrics.reduce((sum, m) => sum + m.loadTime, 0) / conditionMetrics.length,
        category: condition,
      })),
      xAxis: "Network Condition",
      yAxis: "Average Load Time (ms)",
      threshold: PERFORMANCE_BUDGET.maxLoadTime,
    };

    const insights = [
      `Network conditions tested: ${Object.keys(byNetwork).join(", ")}`,
      Object.keys(byNetwork).length > 0 ? 
        `Best network performance: ${Object.entries(byNetwork).reduce((best, [condition, conditionMetrics]) => {
          const avgTime = conditionMetrics.reduce((sum, m) => sum + m.loadTime, 0) / conditionMetrics.length;
          return !best.avgTime || avgTime < best.avgTime ? { condition, avgTime } : best;
        }, { condition: "", avgTime: 0 }).condition}` : "No network data available",
    ];

    return {
      id: "network-analysis",
      title: "Network Performance Analysis",
      type: "detailed",
      content: { byNetwork },
      charts: [chart],
      insights,
    };
  }

  /**
   * Generate recommendations
   */
  private generateRecommendations(
    metrics: ContentLoadMetrics[],
    cacheResults: CacheValidationResult[]
  ): PerformanceRecommendation[] {
    const recommendations: PerformanceRecommendation[] = [];

    // Budget compliance recommendations
    const budgetCompliance = (metrics.filter(m => m.loadTime <= PERFORMANCE_BUDGET.maxLoadTime).length / metrics.length) * 100;
    if (budgetCompliance < 80) {
      recommendations.push({
        id: "budget-compliance",
        title: "Improve Budget Compliance",
        description: `Only ${budgetCompliance.toFixed(1)}% of requests meet the 200ms budget requirement`,
        severity: "high",
        category: "system",
        impact: "high",
        effort: "medium",
        expectedImprovement: "20-40% reduction in load times",
        implementation: [
          "Optimize content delivery pipeline",
          "Implement aggressive caching strategies",
          "Reduce content payload sizes",
          "Optimize network requests"
        ],
        relatedMetrics: ["loadTime", "budgetCompliance"],
      });
    }

    // Cache recommendations
    const cacheHitRatio = (metrics.filter(m => m.cacheStatus === "hit").length / metrics.length) * 100;
    if (cacheHitRatio < 70) {
      recommendations.push({
        id: "cache-optimization",
        title: "Optimize Cache Performance",
        description: `Cache hit ratio is ${cacheHitRatio.toFixed(1)}%, below the recommended 70%`,
        severity: "medium",
        category: "cache",
        impact: "high",
        effort: "low",
        expectedImprovement: "50-80% faster load times for cached content",
        implementation: [
          "Increase cache TTL for stable content",
          "Implement cache warming strategies",
          "Optimize cache key strategies",
          "Consider cache compression"
        ],
        relatedMetrics: ["cacheHitRatio", "cacheStatus"],
      });
    }

    // Content-specific recommendations
    const slowContentTypes = this.identifySlowContentTypes(metrics);
    slowContentTypes.forEach(({ type, avgTime }) => {
      recommendations.push({
        id: `content-optimization-${type}`,
        title: `Optimize ${type} Content Performance`,
        description: `${type} content has average load time of ${avgTime.toFixed(1)}ms`,
        severity: avgTime > PERFORMANCE_BUDGET.maxLoadTime * 2 ? "high" : "medium",
        category: "content",
        impact: "medium",
        effort: "low",
        expectedImprovement: "30-50% reduction in load times",
        implementation: [
          `Optimize ${type} content size`,
          "Implement content-specific caching",
          "Consider lazy loading for non-critical content",
          "Compress content payloads"
        ],
        relatedMetrics: ["loadTime", "contentType"],
      });
    });

    // Add cache-specific recommendations
    cacheResults.forEach(result => {
      result.recommendations.forEach(rec => {
        recommendations.push({
          id: `cache-${result.strategy}-${recommendations.length}`,
          title: `${result.strategy} Cache: ${rec}`,
          description: rec,
          severity: "medium",
          category: "cache",
          impact: "medium",
          effort: "low",
          expectedImprovement: "10-30% improvement in cache performance",
          implementation: [rec],
          relatedMetrics: ["cacheHitRatio", "cacheStatus"],
        });
      });
    });

    return recommendations;
  }

  /**
   * Identify slow content types
   */
  private identifySlowContentTypes(metrics: ContentLoadMetrics[]): Array<{ type: string; avgTime: number }> {
    const byType = this.groupBy(metrics, "contentType");
    
    return Object.entries(byType)
      .map(([type, typeMetrics]) => ({
        type,
        avgTime: typeMetrics.reduce((sum, m) => sum + m.loadTime, 0) / typeMetrics.length,
      }))
      .filter(({ avgTime }) => avgTime > PERFORMANCE_BUDGET.maxLoadTime)
      .sort((a, b) => b.avgTime - a.avgTime);
  }

  /**
   * Generate metadata
   */
  private generateMetadata(metrics: ContentLoadMetrics[]): PerformanceReportMetadata {
    const contentTypes = [...new Set(metrics.map(m => m.contentType))];
    const networkConditions = [...new Set(metrics.map(m => m.networkCondition))];
    
    return {
      testEnvironment: {
        userAgent: typeof navigator !== "undefined" ? navigator.userAgent : "Node.js",
        networkConditions,
        contentTypes,
        testDuration: metrics.length > 0 ? 
          Math.max(...metrics.map(m => m.timestamp)) - Math.min(...metrics.map(m => m.timestamp)) : 0,
      },
      configuration: {
        performanceBudget: PERFORMANCE_BUDGET,
        sampleSize: metrics.length,
        testScenarios: ["load-time", "cache-validation", "network-conditions"],
      },
      version: "1.0.0",
      generatedBy: "Performance Validation System",
    };
  }

  /**
   * Export report to JSON
   */
  exportToJSON(report: PerformanceReport): string {
    return JSON.stringify(report, null, 2);
  }

  /**
   * Export report to HTML
   */
  exportToHTML(report: PerformanceReport): string {
    return `
<!DOCTYPE html>
<html>
<head>
    <title>${report.title}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { border-bottom: 2px solid #333; padding-bottom: 20px; }
        .summary { background: #f5f5f5; padding: 20px; margin: 20px 0; }
        .status-pass { color: green; }
        .status-fail { color: red; }
        .status-warning { color: orange; }
        .section { margin: 30px 0; }
        .recommendations { background: #fff3cd; padding: 15px; margin: 10px 0; }
        .critical { border-left: 4px solid red; }
        .high { border-left: 4px solid orange; }
        .medium { border-left: 4px solid yellow; }
        .low { border-left: 4px solid green; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>${report.title}</h1>
        <p>Generated: ${new Date(report.generatedAt).toLocaleString()}</p>
    </div>
    
    <div class="summary">
        <h2>Summary</h2>
        <p>Status: <span class="status-${report.summary.status}">${report.summary.status.toUpperCase()}</span></p>
        <p>Overall Score: ${report.summary.overallScore}/100</p>
        <p>Budget Compliance: ${report.summary.budgetCompliance.toFixed(1)}%</p>
        <p>Tests: ${report.summary.passedTests}/${report.summary.totalTests} passed</p>
        <ul>
            ${report.summary.keyFindings.map(finding => `<li>${finding}</li>`).join("")}
        </ul>
    </div>
    
    <div class="recommendations">
        <h2>Recommendations</h2>
        ${report.recommendations.map(rec => `
            <div class="${rec.severity}">
                <h3>${rec.title}</h3>
                <p>${rec.description}</p>
                <p><strong>Expected Improvement:</strong> ${rec.expectedImprovement}</p>
                <ul>
                    ${rec.implementation.map(impl => `<li>${impl}</li>`).join("")}
                </ul>
            </div>
        `).join("")}
    </div>
    
    <div class="metadata">
        <h2>Test Configuration</h2>
        <p>Sample Size: ${report.metadata.configuration.sampleSize}</p>
        <p>Content Types: ${report.metadata.testEnvironment.contentTypes.join(", ")}</p>
        <p>Network Conditions: ${report.metadata.testEnvironment.networkConditions.join(", ")}</p>
        <p>Performance Budget: ${report.metadata.configuration.performanceBudget.maxLoadTime}ms</p>
    </div>
</body>
</html>`;
  }

  /**
   * Utility methods
   */
  private groupBy<T>(array: T[], key: keyof T): Record<string, T[]> {
    return array.reduce((groups, item) => {
      const groupKey = String(item[key]);
      groups[groupKey] = groups[groupKey] || [];
      groups[groupKey].push(item);
      return groups;
    }, {} as Record<string, T[]>);
  }

  private calculatePercentile(values: number[], percentile: number): number {
    if (values.length === 0) return 0;
    const sorted = [...values].sort((a, b) => a - b);
    const index = Math.floor((percentile / 100) * sorted.length);
    return sorted[Math.min(index, sorted.length - 1)];
  }

  private createHistogramData(values: number[], buckets: number = 10): ChartDataPoint[] {
    if (values.length === 0) return [];
    
    const min = Math.min(...values);
    const max = Math.max(...values);
    const bucketSize = (max - min) / buckets;
    
    const histogram = new Array(buckets).fill(0);
    
    values.forEach(value => {
      const bucketIndex = Math.min(Math.floor((value - min) / bucketSize), buckets - 1);
      histogram[bucketIndex]++;
    });
    
    return histogram.map((count, index) => ({
      x: min + (index * bucketSize),
      y: count,
    }));
  }

  /**
   * Get report history
   */
  getReportHistory(): PerformanceReport[] {
    return [...this.reportHistory];
  }

  /**
   * Compare reports
   */
  compareReports(reportId1: string, reportId2: string): any {
    const report1 = this.reportHistory.find(r => r.id === reportId1);
    const report2 = this.reportHistory.find(r => r.id === reportId2);
    
    if (!report1 || !report2) {
      throw new Error("One or both reports not found");
    }
    
    return {
      scoreDifference: report2.summary.overallScore - report1.summary.overallScore,
      complianceDifference: report2.summary.budgetCompliance - report1.summary.budgetCompliance,
      recommendations: {
        added: report2.recommendations.filter(r2 => 
          !report1.recommendations.some(r1 => r1.title === r2.title)
        ),
        resolved: report1.recommendations.filter(r1 => 
          !report2.recommendations.some(r2 => r2.title === r1.title)
        ),
      },
    };
  }
}

/**
 * Create performance report generator
 */
export function createPerformanceReportGenerator(): PerformanceReportGenerator {
  return new PerformanceReportGenerator();
}
