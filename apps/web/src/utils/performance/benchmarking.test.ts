/**
 * Performance Benchmarking Test Suite
 * 
 * Comprehensive testing framework for content load performance validation
 * Tests various content types, sizes, and network conditions
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { 
  PerformanceMeasurement, 
  ContentLoadMetrics, 
  PERFORMANCE_BUDGET,
  getCoreWebVitals,
  measureExecutionTime 
} from "./measurement";

// Mock content data for testing
const mockContentTypes = {
  text: {
    small: "Short text content",
    medium: "Medium length text content ".repeat(50),
    large: "Large text content ".repeat(500),
  },
  richText: {
    small: "<p>Simple HTML content</p>",
    medium: `<div>${"<p>Rich text content with formatting</p>".repeat(20)}</div>`,
    large: `<div>${"<p>Large rich text content with complex formatting</p>".repeat(100)}</div>`,
  },
  config: {
    small: { setting: "value" },
    medium: { settings: new Array(50).fill({ key: "value", nested: { data: true } }) },
    large: { settings: new Array(500).fill({ key: "value", nested: { data: true, more: "data" } }) },
  },
  translation: {
    small: { en: "Hello", es: "Hola" },
    medium: Object.fromEntries(
      Array.from({ length: 50 }, (_, i) => [`key${i}`, { en: `English ${i}`, es: `Spanish ${i}` }])
    ),
    large: Object.fromEntries(
      Array.from({ length: 500 }, (_, i) => [`key${i}`, { en: `English ${i}`, es: `Spanish ${i}`, fr: `French ${i}` }])
    ),
  },
};

const mockNetworkConditions = {
  "4g": { latency: 20, bandwidth: 10000 },
  "3g": { latency: 100, bandwidth: 1000 },
  "2g": { latency: 300, bandwidth: 100 },
  "wifi": { latency: 5, bandwidth: 50000 },
  "offline": { latency: 0, bandwidth: 0 },
};

/**
 * Mock network simulation
 */
function simulateNetworkDelay(condition: keyof typeof mockNetworkConditions, size: number): Promise<void> {
  const { latency, bandwidth } = mockNetworkConditions[condition];
  
  if (condition === "offline") {
    return Promise.reject(new Error("Network unavailable"));
  }
  
  // Calculate transmission time based on size and bandwidth
  const transmissionTime = (size / bandwidth) * 1000; // Convert to ms
  const totalDelay = latency + transmissionTime;
  
  return new Promise((resolve) => {
    setTimeout(resolve, totalDelay);
  });
}

/**
 * Calculate content size in bytes
 */
function calculateContentSize(content: any): number {
  return new Blob([JSON.stringify(content)]).size;
}

/**
 * Content load simulation with realistic network behavior
 */
async function simulateContentLoad(
  contentType: string,
  size: "small" | "medium" | "large",
  networkCondition: keyof typeof mockNetworkConditions,
  cacheStatus: "hit" | "miss" | "stale" = "miss"
): Promise<{ content: any; loadTime: number; size: number }> {
  const content = mockContentTypes[contentType as keyof typeof mockContentTypes][size];
  const contentSize = calculateContentSize(content);
  
  const startTime = performance.now();
  
  try {
    // Simulate cache behavior
    if (cacheStatus === "hit") {
      // Cache hit - minimal delay
      await new Promise(resolve => setTimeout(resolve, 1));
    } else if (cacheStatus === "stale") {
      // Stale cache - small delay for validation
      await simulateNetworkDelay(networkCondition, 100);
    } else {
      // Cache miss - full network request
      await simulateNetworkDelay(networkCondition, contentSize);
    }
    
    const endTime = performance.now();
    const loadTime = endTime - startTime;
    
    return { content, loadTime, size: contentSize };
  } catch (error) {
    throw new Error(`Failed to load content: ${error.message}`);
  }
}

describe("Performance Benchmarking Test Suite", () => {
  let performanceMeasurement: PerformanceMeasurement;
  let metricsCollected: ContentLoadMetrics[] = [];

  beforeEach(() => {
    metricsCollected = [];
    performanceMeasurement = new PerformanceMeasurement((metrics) => {
      metricsCollected.push(metrics);
    });
    
    // Clear any existing performance marks
    try {
      performance.clearMarks();
      performance.clearMeasures();
    } catch (error) {
      // Ignore errors in test environment
    }
  });

  afterEach(() => {
    performanceMeasurement.disconnect();
  });

  describe("Content Type Performance Tests", () => {
    const contentTypes = ["text", "richText", "config", "translation"];
    const sizes = ["small", "medium", "large"] as const;
    
    contentTypes.forEach((contentType) => {
      sizes.forEach((size) => {
        it(`should load ${contentType} ${size} content within 200ms budget`, async () => {
          const contentId = `${contentType}-${size}-test`;
          
          performanceMeasurement.startContentLoad(contentId, contentType, "benchmark");
          
          const { content, loadTime, size: contentSize } = await simulateContentLoad(
            contentType,
            size,
            "4g",
            "miss"
          );
          
          const metrics = performanceMeasurement.endContentLoad(contentId, contentSize, "miss");
          
          expect(metrics).toBeTruthy();
          expect(metrics!.contentType).toBe(contentType);
          expect(metrics!.loadTime).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
          expect(content).toBeTruthy();
        });
      });
    });
  });

  describe("Network Condition Performance Tests", () => {
    const networkConditions = ["4g", "3g", "wifi"] as const;
    
    networkConditions.forEach((condition) => {
      it(`should handle ${condition} network conditions efficiently`, async () => {
        const contentId = `network-${condition}-test`;
        
        performanceMeasurement.startContentLoad(contentId, "text", "network-test");
        
        const { loadTime } = await simulateContentLoad("text", "medium", condition, "miss");
        
        const metrics = performanceMeasurement.endContentLoad(contentId, 1000, "miss");
        
        expect(metrics).toBeTruthy();
        
        // Adjust expectations based on network condition
        const expectedMaxTime = condition === "wifi" ? 50 : condition === "4g" ? 150 : 300;
        expect(metrics!.loadTime).toBeLessThanOrEqual(expectedMaxTime);
      });
    });
  });

  describe("Cache Performance Tests", () => {
    const cacheStatuses = ["hit", "miss", "stale"] as const;
    
    cacheStatuses.forEach((cacheStatus) => {
      it(`should handle cache ${cacheStatus} efficiently`, async () => {
        const contentId = `cache-${cacheStatus}-test`;
        
        performanceMeasurement.startContentLoad(contentId, "text", "cache-test");
        
        const { loadTime } = await simulateContentLoad("text", "medium", "4g", cacheStatus);
        
        const metrics = performanceMeasurement.endContentLoad(contentId, 1000, cacheStatus);
        
        expect(metrics).toBeTruthy();
        expect(metrics!.cacheStatus).toBe(cacheStatus);
        
        // Cache hits should be very fast
        if (cacheStatus === "hit") {
          expect(metrics!.loadTime).toBeLessThanOrEqual(10);
        } else if (cacheStatus === "stale") {
          expect(metrics!.loadTime).toBeLessThanOrEqual(50);
        } else {
          expect(metrics!.loadTime).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
        }
      });
    });
  });

  describe("Content Size Impact Tests", () => {
    it("should demonstrate performance impact of content size", async () => {
      const sizes = ["small", "medium", "large"] as const;
      const results: { size: string; loadTime: number; contentSize: number }[] = [];
      
      for (const size of sizes) {
        const contentId = `size-${size}-test`;
        
        performanceMeasurement.startContentLoad(contentId, "text", "size-test");
        
        const { loadTime, size: contentSize } = await simulateContentLoad("text", size, "4g", "miss");
        
        const metrics = performanceMeasurement.endContentLoad(contentId, contentSize, "miss");
        
        expect(metrics).toBeTruthy();
        results.push({ size, loadTime: metrics!.loadTime, contentSize });
      }
      
      // Verify that larger content takes more time (generally)
      expect(results[0].contentSize).toBeLessThan(results[1].contentSize);
      expect(results[1].contentSize).toBeLessThan(results[2].contentSize);
      
      // All should still meet performance budget
      results.forEach(result => {
        expect(result.loadTime).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
      });
    });
  });

  describe("Bulk Content Loading Tests", () => {
    it("should handle concurrent content loading efficiently", async () => {
      const contentIds = Array.from({ length: 10 }, (_, i) => `bulk-content-${i}`);
      const loadPromises: Promise<ContentLoadMetrics | null>[] = [];
      
      // Start all measurements
      contentIds.forEach((contentId, index) => {
        performanceMeasurement.startContentLoad(contentId, "text", "bulk-test");
        
        const loadPromise = simulateContentLoad("text", "small", "4g", "miss")
          .then(({ size }) => performanceMeasurement.endContentLoad(contentId, size, "miss"));
        
        loadPromises.push(loadPromise);
      });
      
      const results = await Promise.all(loadPromises);
      
      // All loads should complete successfully
      results.forEach((metrics, index) => {
        expect(metrics).toBeTruthy();
        expect(metrics!.contentId).toBe(contentIds[index]);
        expect(metrics!.loadTime).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
      });
      
      // Check that bulk loading didnt significantly impact individual load times
      const avgLoadTime = results.reduce((sum, m) => sum + m!.loadTime, 0) / results.length;
      expect(avgLoadTime).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
    });
  });

  describe("Performance Measurement API Tests", () => {
    it("should correctly measure execution time", async () => {
      const testFunction = async () => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return "test result";
      };
      
      const { result, duration } = await measureExecutionTime(testFunction, "test-function");
      
      expect(result).toBe("test result");
      expect(duration).toBeGreaterThanOrEqual(45);
      expect(duration).toBeLessThanOrEqual(100);
    });

    it("should collect comprehensive performance summary", async () => {
      // Generate multiple measurements
      const testCount = 5;
      
      for (let i = 0; i < testCount; i++) {
        const contentId = `summary-test-${i}`;
        performanceMeasurement.startContentLoad(contentId, "text", "summary-test");
        
        const { size } = await simulateContentLoad("text", "small", "4g", "miss");
        performanceMeasurement.endContentLoad(contentId, size, "miss");
      }
      
      const summary = performanceMeasurement.getPerformanceSummary();
      
      expect(summary.totalMeasurements).toBe(testCount);
      expect(summary.averageLoadTime).toBeGreaterThan(0);
      expect(summary.maxLoadTime).toBeGreaterThanOrEqual(summary.minLoadTime);
      expect(summary.passRate).toBeGreaterThanOrEqual(0);
      expect(summary.passRate).toBeLessThanOrEqual(1);
      expect(summary.failedCount).toBe(testCount - Math.floor(summary.passRate * testCount));
    });
  });

  describe("Core Web Vitals Integration", () => {
    it("should measure core web vitals", async () => {
      const metrics = await getCoreWebVitals();
      
      expect(metrics).toBeTruthy();
      expect(typeof metrics.loadTime).toBe("number");
      expect(typeof metrics.firstContentfulPaint).toBe("number");
      expect(typeof metrics.timeToInteractive).toBe("number");
      expect(typeof metrics.resourceLoadTime).toBe("number");
      expect(typeof metrics.networkLatency).toBe("number");
    });
  });

  describe("Error Handling Tests", () => {
    it("should handle network failures gracefully", async () => {
      const contentId = "error-test";
      
      performanceMeasurement.startContentLoad(contentId, "text", "error-test");
      
      try {
        await simulateContentLoad("text", "small", "offline", "miss");
        expect.fail("Should have thrown an error");
      } catch (error) {
        expect(error.message).toContain("Network unavailable");
      }
      
      // Measurement should still be cleanable
      const measurements = performanceMeasurement.getAllMeasurements();
      expect(measurements.some(m => m.contentId === contentId)).toBe(true);
    });
  });

  describe("Performance Regression Tests", () => {
    it("should detect performance regressions", async () => {
      const baselineResults: number[] = [];
      const testResults: number[] = [];
      
      // Simulate baseline performance
      for (let i = 0; i < 5; i++) {
        const contentId = `baseline-${i}`;
        performanceMeasurement.startContentLoad(contentId, "text", "regression-test");
        
        const { size } = await simulateContentLoad("text", "medium", "4g", "miss");
        const metrics = performanceMeasurement.endContentLoad(contentId, size, "miss");
        
        if (metrics) {
          baselineResults.push(metrics.loadTime);
        }
      }
      
      // Simulate potentially regressed performance
      for (let i = 0; i < 5; i++) {
        const contentId = `test-${i}`;
        performanceMeasurement.startContentLoad(contentId, "text", "regression-test");
        
        // Add slight delay to simulate regression
        await new Promise(resolve => setTimeout(resolve, 10));
        
        const { size } = await simulateContentLoad("text", "medium", "4g", "miss");
        const metrics = performanceMeasurement.endContentLoad(contentId, size, "miss");
        
        if (metrics) {
          testResults.push(metrics.loadTime);
        }
      }
      
      const baselineAvg = baselineResults.reduce((a, b) => a + b, 0) / baselineResults.length;
      const testAvg = testResults.reduce((a, b) => a + b, 0) / testResults.length;
      
      // Both should still meet budget
      expect(baselineAvg).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
      expect(testAvg).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
      
      // Regression detection (test avg should not be significantly higher)
      const regressionThreshold = 1.5; // 50% degradation threshold
      expect(testAvg).toBeLessThanOrEqual(baselineAvg * regressionThreshold);
    });
  });
});

describe("Performance Budget Validation", () => {
  it("should validate all performance budgets are realistic", () => {
    expect(PERFORMANCE_BUDGET.maxLoadTime).toBe(200);
    expect(PERFORMANCE_BUDGET.maxFirstContentfulPaint).toBeGreaterThan(0);
    expect(PERFORMANCE_BUDGET.maxLargestContentfulPaint).toBeGreaterThan(0);
    expect(PERFORMANCE_BUDGET.maxFirstInputDelay).toBeGreaterThan(0);
    expect(PERFORMANCE_BUDGET.maxCumulativeLayoutShift).toBeGreaterThan(0);
    expect(PERFORMANCE_BUDGET.maxResourceLoadTime).toBeGreaterThan(0);
    expect(PERFORMANCE_BUDGET.minCacheHitRatio).toBeGreaterThan(0);
    expect(PERFORMANCE_BUDGET.minCacheHitRatio).toBeLessThanOrEqual(1);
  });

  it("should enforce strict 200ms requirement for content loading", () => {
    // This is the core requirement - must always be 200ms
    expect(PERFORMANCE_BUDGET.maxLoadTime).toBe(200);
    
    // Other budgets should be reasonable but secondary
    expect(PERFORMANCE_BUDGET.maxResourceLoadTime).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
  });
});
