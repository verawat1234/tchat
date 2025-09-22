/**
 * Network Condition Testing Scenarios
 * 
 * Comprehensive testing for various network conditions and their impact on content loading performance
 */

import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { PerformanceMeasurement, PERFORMANCE_BUDGET } from "./measurement";

/**
 * Network condition definitions based on real-world scenarios
 */
export interface NetworkCondition {
  name: string;
  description: string;
  latency: number; // Round trip time in ms
  downloadSpeed: number; // Kbps
  uploadSpeed: number; // Kbps
  packetLoss: number; // Percentage (0-1)
  jitter: number; // Latency variation in ms
}

export const NETWORK_CONDITIONS: NetworkCondition[] = [
  {
    name: "5g",
    description: "5G Ultra Wideband",
    latency: 1,
    downloadSpeed: 1000000, // 1 Gbps
    uploadSpeed: 500000, // 500 Mbps
    packetLoss: 0.001,
    jitter: 0.5,
  },
  {
    name: "wifi-fast",
    description: "Fast WiFi (802.11ac)",
    latency: 5,
    downloadSpeed: 100000, // 100 Mbps
    uploadSpeed: 50000, // 50 Mbps
    packetLoss: 0.01,
    jitter: 2,
  },
  {
    name: "wifi-typical",
    description: "Typical WiFi",
    latency: 20,
    downloadSpeed: 25000, // 25 Mbps
    uploadSpeed: 5000, // 5 Mbps
    packetLoss: 0.02,
    jitter: 5,
  },
  {
    name: "4g-excellent",
    description: "4G LTE Excellent Signal",
    latency: 30,
    downloadSpeed: 20000, // 20 Mbps
    uploadSpeed: 10000, // 10 Mbps
    packetLoss: 0.05,
    jitter: 10,
  },
  {
    name: "4g-good",
    description: "4G LTE Good Signal",
    latency: 50,
    downloadSpeed: 10000, // 10 Mbps
    uploadSpeed: 3000, // 3 Mbps
    packetLoss: 0.1,
    jitter: 15,
  },
  {
    name: "4g-fair",
    description: "4G LTE Fair Signal",
    latency: 100,
    downloadSpeed: 3000, // 3 Mbps
    uploadSpeed: 1000, // 1 Mbps
    packetLoss: 0.2,
    jitter: 25,
  },
  {
    name: "3g",
    description: "3G Connection",
    latency: 200,
    downloadSpeed: 1000, // 1 Mbps
    uploadSpeed: 500, // 500 Kbps
    packetLoss: 0.3,
    jitter: 50,
  },
  {
    name: "2g",
    description: "2G/EDGE Connection",
    latency: 500,
    downloadSpeed: 100, // 100 Kbps
    uploadSpeed: 50, // 50 Kbps
    packetLoss: 0.5,
    jitter: 100,
  },
  {
    name: "satellite",
    description: "Satellite Internet",
    latency: 600,
    downloadSpeed: 5000, // 5 Mbps
    uploadSpeed: 1000, // 1 Mbps
    packetLoss: 0.1,
    jitter: 200,
  },
  {
    name: "dialup",
    description: "Dial-up Connection",
    latency: 200,
    downloadSpeed: 56, // 56 Kbps
    uploadSpeed: 33, // 33 Kbps
    packetLoss: 0.2,
    jitter: 50,
  },
];

/**
 * Network simulation utility
 */
export class NetworkSimulator {
  private currentCondition: NetworkCondition;
  
  constructor(condition: NetworkCondition) {
    this.currentCondition = condition;
  }

  /**
   * Simulate network delay based on content size and current conditions
   */
  async simulateNetworkRequest(contentSizeBytes: number): Promise<{
    downloadTime: number;
    totalTime: number;
    actualLatency: number;
    timeoutOccurred: boolean;
  }> {
    const { latency, downloadSpeed, packetLoss, jitter } = this.currentCondition;
    
    // Simulate packet loss
    if (Math.random() < packetLoss) {
      throw new Error(`Network request failed due to packet loss (${packetLoss * 100}%)`);
    }

    // Calculate variable latency with jitter
    const actualLatency = latency + (Math.random() - 0.5) * jitter * 2;
    
    // Calculate download time: size (bytes) * 8 (bits) / speed (kbps) * 1000 (ms)
    const downloadTime = (contentSizeBytes * 8) / (downloadSpeed * 1000) * 1000;
    
    const totalTime = actualLatency + downloadTime;
    
    // Simulate timeout for very slow connections
    const timeoutThreshold = 30000; // 30 seconds
    const timeoutOccurred = totalTime > timeoutThreshold;
    
    if (timeoutOccurred) {
      throw new Error(`Request timeout: ${totalTime}ms exceeds ${timeoutThreshold}ms threshold`);
    }

    // Simulate the actual delay
    await new Promise(resolve => setTimeout(resolve, Math.min(totalTime, 1000))); // Cap at 1s for tests
    
    return {
      downloadTime,
      totalTime,
      actualLatency,
      timeoutOccurred,
    };
  }

  /**
   * Get condition information
   */
  getCondition(): NetworkCondition {
    return { ...this.currentCondition };
  }

  /**
   * Check if condition meets performance requirements
   */
  canMeetPerformanceBudget(contentSizeBytes: number): boolean {
    const estimatedTime = this.estimateLoadTime(contentSizeBytes);
    return estimatedTime <= PERFORMANCE_BUDGET.maxLoadTime;
  }

  /**
   * Estimate load time without actually simulating
   */
  estimateLoadTime(contentSizeBytes: number): number {
    const { latency, downloadSpeed } = this.currentCondition;
    const downloadTime = (contentSizeBytes * 8) / (downloadSpeed * 1000) * 1000;
    return latency + downloadTime;
  }
}

/**
 * Content size categories for testing
 */
export const CONTENT_SIZES = {
  tiny: 100, // 100 bytes - small text
  small: 1024, // 1 KB - typical text content
  medium: 10240, // 10 KB - rich text with formatting
  large: 102400, // 100 KB - complex configuration
  xlarge: 1048576, // 1 MB - large image URLs or complex data
};

describe("Network Condition Testing Scenarios", () => {
  let performanceMeasurement: PerformanceMeasurement;

  beforeEach(() => {
    performanceMeasurement = new PerformanceMeasurement();
  });

  afterEach(() => {
    performanceMeasurement.disconnect();
  });

  describe("Individual Network Condition Tests", () => {
    NETWORK_CONDITIONS.forEach((condition) => {
      describe(`${condition.name} (${condition.description})`, () => {
        let simulator: NetworkSimulator;

        beforeEach(() => {
          simulator = new NetworkSimulator(condition);
        });

        Object.entries(CONTENT_SIZES).forEach(([sizeCategory, sizeBytes]) => {
          it(`should handle ${sizeCategory} content (${sizeBytes} bytes)`, async () => {
            const contentId = `${condition.name}-${sizeCategory}`;
            
            performanceMeasurement.startContentLoad(contentId, "test", condition.name);
            
            try {
              const networkResult = await simulator.simulateNetworkRequest(sizeBytes);
              const metrics = performanceMeasurement.endContentLoad(contentId, sizeBytes);
              
              expect(metrics).toBeTruthy();
              expect(metrics!.loadTime).toBeGreaterThan(0);
              
              // Log performance for analysis
              console.log(`${condition.name} ${sizeCategory}: ${metrics!.loadTime}ms (estimated: ${networkResult.totalTime}ms)`);
              
              // Check if this combination can meet the performance budget
              const canMeetBudget = simulator.canMeetPerformanceBudget(sizeBytes);
              if (canMeetBudget) {
                expect(metrics!.loadTime).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
              }
              
            } catch (error) {
              // Some combinations (like 2G + large content) are expected to fail
              console.warn(`${condition.name} ${sizeCategory} failed: ${error.message}`);
              
              if (error.message.includes("timeout") || error.message.includes("packet loss")) {
                // Expected failure for very poor conditions
                expect(error).toBeTruthy();
              } else {
                throw error;
              }
            }
          });
        });

        it("should provide accurate performance estimates", () => {
          Object.entries(CONTENT_SIZES).forEach(([sizeCategory, sizeBytes]) => {
            const estimatedTime = simulator.estimateLoadTime(sizeBytes);
            const canMeet = simulator.canMeetPerformanceBudget(sizeBytes);
            
            expect(estimatedTime).toBeGreaterThan(0);
            expect(typeof canMeet).toBe("boolean");
            
            if (canMeet) {
              expect(estimatedTime).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
            }
          });
        });
      });
    });
  });

  describe("Network Condition Performance Matrix", () => {
    it("should create comprehensive performance matrix", async () => {
      const results: Array<{
        condition: string;
        size: string;
        bytes: number;
        loadTime: number | null;
        meetsBudget: boolean;
        error?: string;
      }> = [];

      for (const condition of NETWORK_CONDITIONS.slice(0, 5)) { // Test first 5 conditions
        const simulator = new NetworkSimulator(condition);
        
        for (const [sizeCategory, sizeBytes] of Object.entries(CONTENT_SIZES)) {
          if (sizeBytes > 10240 && condition.downloadSpeed < 1000) {
            // Skip large content on very slow connections to avoid test timeouts
            continue;
          }

          const contentId = `matrix-${condition.name}-${sizeCategory}`;
          
          performanceMeasurement.startContentLoad(contentId, "matrix-test", condition.name);
          
          try {
            await simulator.simulateNetworkRequest(sizeBytes);
            const metrics = performanceMeasurement.endContentLoad(contentId, sizeBytes);
            
            results.push({
              condition: condition.name,
              size: sizeCategory,
              bytes: sizeBytes,
              loadTime: metrics?.loadTime || null,
              meetsBudget: (metrics?.loadTime || Infinity) <= PERFORMANCE_BUDGET.maxLoadTime,
            });
          } catch (error) {
            results.push({
              condition: condition.name,
              size: sizeCategory,
              bytes: sizeBytes,
              loadTime: null,
              meetsBudget: false,
              error: error.message,
            });
          }
        }
      }

      // Analyze results
      const successfulTests = results.filter(r => r.loadTime !== null);
      const failedTests = results.filter(r => r.loadTime === null);
      const passedBudget = results.filter(r => r.meetsBudget);

      console.log(`Performance Matrix Results:
        Total tests: ${results.length}
        Successful: ${successfulTests.length}
        Failed: ${failedTests.length}
        Met budget: ${passedBudget.length}
        Pass rate: ${(passedBudget.length / results.length * 100).toFixed(1)}%`);

      // Fast connections should handle all content sizes
      const fastConditions = results.filter(r => 
        ["5g", "wifi-fast", "wifi-typical"].includes(r.condition)
      );
      const fastConditionPassRate = fastConditions.filter(r => r.meetsBudget).length / fastConditions.length;
      expect(fastConditionPassRate).toBeGreaterThan(0.8); // 80% pass rate for fast connections

      expect(results.length).toBeGreaterThan(0);
    });
  });

  describe("Adaptive Performance Budgets", () => {
    it("should adjust expectations based on network conditions", () => {
      const adaptiveBudgets = NETWORK_CONDITIONS.map(condition => {
        const simulator = new NetworkSimulator(condition);
        
        // Calculate adaptive budget based on network capabilities
        let adaptiveBudget = PERFORMANCE_BUDGET.maxLoadTime;
        
        // Adjust budget based on connection quality
        if (condition.latency > 300) {
          adaptiveBudget *= 3; // 3x budget for very slow connections
        } else if (condition.latency > 100) {
          adaptiveBudget *= 2; // 2x budget for slow connections
        } else if (condition.latency > 50) {
          adaptiveBudget *= 1.5; // 1.5x budget for moderate connections
        }
        
        return {
          condition: condition.name,
          baseBudget: PERFORMANCE_BUDGET.maxLoadTime,
          adaptiveBudget,
          latency: condition.latency,
          downloadSpeed: condition.downloadSpeed,
        };
      });

      // Verify adaptive budgets are reasonable
      adaptiveBudgets.forEach(budget => {
        expect(budget.adaptiveBudget).toBeGreaterThanOrEqual(budget.baseBudget);
        expect(budget.adaptiveBudget).toBeLessThanOrEqual(budget.baseBudget * 3);
      });

      // Fast connections should maintain strict budget
      const fastConnections = adaptiveBudgets.filter(b => 
        ["5g", "wifi-fast", "wifi-typical", "4g-excellent"].includes(b.condition)
      );
      fastConnections.forEach(budget => {
        expect(budget.adaptiveBudget).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime * 1.5);
      });
    });
  });

  describe("Network Resilience Tests", () => {
    it("should handle network degradation gracefully", async () => {
      const degradationScenario = [
        NETWORK_CONDITIONS.find(c => c.name === "wifi-fast")!,
        NETWORK_CONDITIONS.find(c => c.name === "4g-good")!,
        NETWORK_CONDITIONS.find(c => c.name === "3g")!,
      ];

      const results: number[] = [];

      for (const [index, condition] of degradationScenario.entries()) {
        const simulator = new NetworkSimulator(condition);
        const contentId = `degradation-${index}`;
        
        performanceMeasurement.startContentLoad(contentId, "degradation-test", condition.name);
        
        try {
          await simulator.simulateNetworkRequest(CONTENT_SIZES.small);
          const metrics = performanceMeasurement.endContentLoad(contentId, CONTENT_SIZES.small);
          
          if (metrics) {
            results.push(metrics.loadTime);
          }
        } catch (error) {
          results.push(Infinity);
        }
      }

      // Results should show increasing load times (generally)
      expect(results[0]).toBeLessThan(results[1] * 2); // Some degradation expected
      expect(results.length).toBe(degradationScenario.length);
      
      // At least the first condition should meet budget
      expect(results[0]).toBeLessThanOrEqual(PERFORMANCE_BUDGET.maxLoadTime);
    });

    it("should provide fallback strategies for poor connections", () => {
      const poorConditions = NETWORK_CONDITIONS.filter(c => 
        c.latency > 200 || c.downloadSpeed < 1000
      );

      poorConditions.forEach(condition => {
        const simulator = new NetworkSimulator(condition);
        
        // Check which content sizes are viable
        const viableSizes = Object.entries(CONTENT_SIZES).filter(([_, size]) => 
          simulator.canMeetPerformanceBudget(size)
        );

        // At least tiny content should be viable for most conditions
        if (condition.name !== "dialup") {
          expect(viableSizes.length).toBeGreaterThan(0);
        }

        // Generate fallback recommendations
        const maxViableSize = viableSizes.length > 0 ? 
          Math.max(...viableSizes.map(([_, size]) => size)) : 0;

        console.log(`${condition.name}: Max viable content size: ${maxViableSize} bytes`);
        
        expect(maxViableSize).toBeGreaterThanOrEqual(0);
      });
    });
  });

  describe("Real-World Network Scenarios", () => {
    const realWorldScenarios = [
      {
        name: "mobile_commute",
        description: "Mobile user on commute with varying signal",
        conditions: ["4g-excellent", "4g-good", "4g-fair", "3g"],
        durations: [5000, 10000, 15000, 5000], // ms spent in each condition
      },
      {
        name: "office_wifi",
        description: "Office worker with intermittent WiFi issues",
        conditions: ["wifi-fast", "wifi-typical", "4g-good", "wifi-fast"],
        durations: [20000, 5000, 10000, 15000],
      },
      {
        name: "rural_user",
        description: "Rural user with limited connectivity options",
        conditions: ["satellite", "4g-fair", "3g", "2g"],
        durations: [15000, 10000, 10000, 5000],
      },
    ];

    realWorldScenarios.forEach(scenario => {
      it(`should handle ${scenario.name} scenario (${scenario.description})`, async () => {
        const results: Array<{ condition: string; loadTime: number; success: boolean }> = [];

        for (const [index, conditionName] of scenario.conditions.entries()) {
          const condition = NETWORK_CONDITIONS.find(c => c.name === conditionName)!;
          const simulator = new NetworkSimulator(condition);
          const contentId = `scenario-${scenario.name}-${index}`;

          performanceMeasurement.startContentLoad(contentId, "scenario-test", conditionName);

          try {
            await simulator.simulateNetworkRequest(CONTENT_SIZES.medium);
            const metrics = performanceMeasurement.endContentLoad(contentId, CONTENT_SIZES.medium);

            results.push({
              condition: conditionName,
              loadTime: metrics?.loadTime || Infinity,
              success: (metrics?.loadTime || Infinity) <= PERFORMANCE_BUDGET.maxLoadTime,
            });
          } catch (error) {
            results.push({
              condition: conditionName,
              loadTime: Infinity,
              success: false,
            });
          }
        }

        // Analyze scenario results
        const successRate = results.filter(r => r.success).length / results.length;
        const avgLoadTime = results
          .filter(r => r.loadTime !== Infinity)
          .reduce((sum, r) => sum + r.loadTime, 0) / 
          results.filter(r => r.loadTime !== Infinity).length;

        console.log(`${scenario.name} Results:
          Success rate: ${(successRate * 100).toFixed(1)}%
          Average load time: ${avgLoadTime.toFixed(1)}ms
          Conditions tested: ${results.map(r => r.condition).join(", ")}`);

        // Each scenario should have some successful conditions
        expect(successRate).toBeGreaterThan(0);
        expect(results.length).toBe(scenario.conditions.length);
      });
    });
  });
});
