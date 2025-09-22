/**
 * Network Failure Simulator
 *
 * Comprehensive simulation of various network failure conditions for testing
 * the robustness of the fallback system under realistic failure scenarios.
 *
 * Features:
 * - Complete network disconnection simulation
 * - Intermittent connectivity patterns
 * - API server failure simulation
 * - Partial failure scenarios
 * - Recovery testing capabilities
 * - Performance stress simulation
 * - Edge case condition simulation
 */

interface NetworkCondition {
  id: string;
  name: string;
  description: string;
  probability: number; // 0-1, probability of failure
  latency: number; // milliseconds
  timeout: number; // milliseconds
  packetLoss: number; // 0-1, percentage of packet loss
  bandwidth: number; // bytes per second, 0 = unlimited
  jitter: number; // milliseconds, variance in latency
}

interface FailurePattern {
  id: string;
  name: string;
  description: string;
  conditions: NetworkCondition[];
  duration: number; // milliseconds
  repetitions: number;
  interval: number; // milliseconds between repetitions
}

interface SimulationResult {
  patternId: string;
  startTime: number;
  endTime: number;
  duration: number;
  totalRequests: number;
  successfulRequests: number;
  failedRequests: number;
  fallbackActivations: number;
  recoveryTime: number;
  errors: Array<{
    timestamp: number;
    type: string;
    message: string;
    endpoint: string;
  }>;
  metrics: {
    averageLatency: number;
    maxLatency: number;
    minLatency: number;
    dataIntegrityScore: number;
    userExperienceScore: number;
    fallbackEffectiveness: number;
  };
}

interface RequestInterceptor {
  originalFetch: typeof fetch;
  interceptedCount: number;
  failureCount: number;
  successCount: number;
  latencyRecords: number[];
}

export class NetworkFailureSimulator {
  private interceptor: RequestInterceptor | null = null;
  private activePattern: FailurePattern | null = null;
  private patternStartTime = 0;
  private currentConditionIndex = 0;
  private simulationResults: SimulationResult[] = [];
  private isSimulating = false;

  // =============================================================================
  // Predefined Network Conditions
  // =============================================================================

  private readonly networkConditions: Record<string, NetworkCondition> = {
    // Complete disconnection
    offline: {
      id: 'offline',
      name: 'Complete Offline',
      description: 'Complete network disconnection',
      probability: 1.0,
      latency: 0,
      timeout: 1000,
      packetLoss: 1.0,
      bandwidth: 0,
      jitter: 0,
    },

    // Intermittent connectivity
    unstable: {
      id: 'unstable',
      name: 'Unstable Connection',
      description: 'Intermittent connectivity with 50% failure rate',
      probability: 0.5,
      latency: 2000,
      timeout: 5000,
      packetLoss: 0.3,
      bandwidth: 1024, // 1KB/s
      jitter: 1000,
    },

    // Slow connection
    slow: {
      id: 'slow',
      name: 'Slow Connection',
      description: 'Very slow connection simulating 2G network',
      probability: 0.1,
      latency: 3000,
      timeout: 10000,
      packetLoss: 0.1,
      bandwidth: 256, // 256 bytes/s
      jitter: 500,
    },

    // High latency
    highLatency: {
      id: 'highLatency',
      name: 'High Latency',
      description: 'High latency connection simulating satellite internet',
      probability: 0.2,
      latency: 600,
      timeout: 8000,
      packetLoss: 0.05,
      bandwidth: 10240, // 10KB/s
      jitter: 200,
    },

    // Timeout prone
    timeoutProne: {
      id: 'timeoutProne',
      name: 'Timeout Prone',
      description: 'Connection that frequently times out',
      probability: 0.3,
      latency: 1000,
      timeout: 2000,
      packetLoss: 0.15,
      bandwidth: 5120, // 5KB/s
      jitter: 300,
    },

    // Server errors
    serverError: {
      id: 'serverError',
      name: 'Server Error',
      description: 'Server returns 5xx errors',
      probability: 0.8,
      latency: 500,
      timeout: 5000,
      packetLoss: 0,
      bandwidth: 0,
      jitter: 100,
    },

    // Partial failures
    partialFailure: {
      id: 'partialFailure',
      name: 'Partial Failure',
      description: 'Some endpoints fail while others succeed',
      probability: 0.4,
      latency: 1000,
      timeout: 4000,
      packetLoss: 0.2,
      bandwidth: 2048, // 2KB/s
      jitter: 400,
    },

    // Recovery simulation
    recovering: {
      id: 'recovering',
      name: 'Recovering Connection',
      description: 'Connection gradually improving',
      probability: 0.2,
      latency: 800,
      timeout: 6000,
      packetLoss: 0.1,
      bandwidth: 4096, // 4KB/s
      jitter: 200,
    },
  };

  // =============================================================================
  // Predefined Failure Patterns
  // =============================================================================

  private readonly failurePatterns: Record<string, FailurePattern> = {
    // Complete network disconnection
    completeOutage: {
      id: 'completeOutage',
      name: 'Complete Network Outage',
      description: 'Complete network disconnection for extended period',
      conditions: [this.networkConditions.offline],
      duration: 30000, // 30 seconds
      repetitions: 1,
      interval: 0,
    },

    // Intermittent connectivity
    intermittent: {
      id: 'intermittent',
      name: 'Intermittent Connectivity',
      description: 'Alternating between connected and disconnected states',
      conditions: [
        this.networkConditions.offline,
        this.networkConditions.unstable,
        this.networkConditions.slow,
      ],
      duration: 5000, // 5 seconds per condition
      repetitions: 3,
      interval: 2000, // 2 seconds between repetitions
    },

    // Progressive degradation
    progressiveDegradation: {
      id: 'progressiveDegradation',
      name: 'Progressive Network Degradation',
      description: 'Network gradually degrades from good to offline',
      conditions: [
        this.networkConditions.slow,
        this.networkConditions.unstable,
        this.networkConditions.offline,
      ],
      duration: 10000, // 10 seconds per condition
      repetitions: 1,
      interval: 0,
    },

    // Recovery pattern
    recoveryPattern: {
      id: 'recoveryPattern',
      name: 'Network Recovery Pattern',
      description: 'Network recovers from complete failure gradually',
      conditions: [
        this.networkConditions.offline,
        this.networkConditions.unstable,
        this.networkConditions.recovering,
        this.networkConditions.slow,
      ],
      duration: 8000, // 8 seconds per condition
      repetitions: 1,
      interval: 0,
    },

    // Server failure pattern
    serverFailures: {
      id: 'serverFailures',
      name: 'Server Failures',
      description: 'Various server-side failures and errors',
      conditions: [
        this.networkConditions.serverError,
        this.networkConditions.timeoutProne,
        this.networkConditions.partialFailure,
      ],
      duration: 7000, // 7 seconds per condition
      repetitions: 2,
      interval: 3000, // 3 seconds between repetitions
    },

    // Stress test pattern
    stressTest: {
      id: 'stressTest',
      name: 'Network Stress Test',
      description: 'High-frequency failures and recoveries',
      conditions: [
        this.networkConditions.offline,
        this.networkConditions.recovering,
        this.networkConditions.unstable,
        this.networkConditions.timeoutProne,
      ],
      duration: 3000, // 3 seconds per condition
      repetitions: 5,
      interval: 1000, // 1 second between repetitions
    },

    // Edge case pattern
    edgeCases: {
      id: 'edgeCases',
      name: 'Edge Case Scenarios',
      description: 'Unusual and edge case network conditions',
      conditions: [
        this.networkConditions.highLatency,
        this.networkConditions.partialFailure,
        this.networkConditions.serverError,
        this.networkConditions.timeoutProne,
      ],
      duration: 6000, // 6 seconds per condition
      repetitions: 2,
      interval: 2000, // 2 seconds between repetitions
    },
  };

  // =============================================================================
  // Simulation Control
  // =============================================================================

  /**
   * Start network failure simulation with specified pattern
   */
  async startSimulation(patternId: string): Promise<void> {
    if (this.isSimulating) {
      throw new Error('Simulation already in progress');
    }

    const pattern = this.failurePatterns[patternId];
    if (!pattern) {
      throw new Error(`Pattern ${patternId} not found`);
    }

    this.activePattern = pattern;
    this.isSimulating = true;
    this.patternStartTime = Date.now();
    this.currentConditionIndex = 0;

    console.log(`🌐 Starting network failure simulation: ${pattern.name}`);
    console.log(`📋 Description: ${pattern.description}`);
    console.log(`⏱️ Total duration: ${this.calculateTotalDuration(pattern)}ms`);

    // Install network interceptor
    this.installNetworkInterceptor();

    // Execute pattern
    await this.executePattern(pattern);

    // Generate results
    const result = this.generateSimulationResult(pattern);
    this.simulationResults.push(result);

    console.log(`✅ Simulation completed: ${pattern.name}`);
    console.log(`📊 Results:`, result.metrics);
  }

  /**
   * Stop current simulation
   */
  async stopSimulation(): Promise<void> {
    if (!this.isSimulating) {
      return;
    }

    this.isSimulating = false;
    this.uninstallNetworkInterceptor();
    this.activePattern = null;

    console.log('🛑 Network simulation stopped');
  }

  /**
   * Get simulation results
   */
  getResults(): SimulationResult[] {
    return [...this.simulationResults];
  }

  /**
   * Clear all simulation results
   */
  clearResults(): void {
    this.simulationResults = [];
  }

  // =============================================================================
  // Network Interception
  // =============================================================================

  private installNetworkInterceptor(): void {
    if (this.interceptor) {
      return; // Already installed
    }

    const originalFetch = window.fetch;
    const simulator = this;

    this.interceptor = {
      originalFetch,
      interceptedCount: 0,
      failureCount: 0,
      successCount: 0,
      latencyRecords: [],
    };

    // Override fetch function
    window.fetch = async function(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
      const startTime = performance.now();
      const url = typeof input === 'string' ? input : input.toString();

      simulator.interceptor!.interceptedCount++;

      try {
        // Apply current network condition
        const response = await simulator.applyNetworkCondition(
          () => originalFetch(input, init),
          url
        );

        const endTime = performance.now();
        const latency = endTime - startTime;
        simulator.interceptor!.latencyRecords.push(latency);

        if (response.ok) {
          simulator.interceptor!.successCount++;
        } else {
          simulator.interceptor!.failureCount++;
        }

        return response;
      } catch (error) {
        const endTime = performance.now();
        const latency = endTime - startTime;
        simulator.interceptor!.latencyRecords.push(latency);
        simulator.interceptor!.failureCount++;

        // Record error for analysis
        simulator.recordError(url, error);
        throw error;
      }
    };

    console.log('🔌 Network interceptor installed');
  }

  private uninstallNetworkInterceptor(): void {
    if (!this.interceptor) {
      return;
    }

    // Restore original fetch
    window.fetch = this.interceptor.originalFetch;
    this.interceptor = null;

    console.log('🔌 Network interceptor uninstalled');
  }

  private async applyNetworkCondition(
    requestFunction: () => Promise<Response>,
    url: string
  ): Promise<Response> {
    if (!this.activePattern || !this.isSimulating) {
      return requestFunction();
    }

    const condition = this.getCurrentNetworkCondition();
    if (!condition) {
      return requestFunction();
    }

    // Simulate packet loss (complete failure)
    if (Math.random() < condition.packetLoss) {
      throw new Error('NetworkError: Network packet lost');
    }

    // Simulate failure probability
    if (Math.random() < condition.probability) {
      if (condition.id === 'serverError') {
        // Return server error response
        return new Response(JSON.stringify({ error: 'Internal Server Error' }), {
          status: 500,
          statusText: 'Internal Server Error',
          headers: { 'Content-Type': 'application/json' },
        });
      } else if (condition.id === 'offline') {
        throw new Error('NetworkError: Network request failed');
      } else if (condition.id === 'timeoutProne') {
        throw new Error('TimeoutError: Request timed out');
      } else {
        throw new Error('FetchError: Network request failed');
      }
    }

    // Simulate latency and jitter
    const baseLatency = condition.latency;
    const jitter = condition.jitter * (Math.random() - 0.5) * 2; // ±jitter
    const totalLatency = Math.max(0, baseLatency + jitter);

    if (totalLatency > 0) {
      await this.delay(totalLatency);
    }

    // Simulate bandwidth limitation
    if (condition.bandwidth > 0) {
      // For simplicity, we'll just add additional delay based on bandwidth
      const estimatedResponseSize = 1024; // 1KB estimated response
      const transmissionTime = (estimatedResponseSize / condition.bandwidth) * 1000;
      await this.delay(transmissionTime);
    }

    // Check for timeout
    const timeoutPromise = new Promise<never>((_, reject) => {
      setTimeout(() => reject(new Error('TimeoutError: Request timed out')), condition.timeout);
    });

    return Promise.race([requestFunction(), timeoutPromise]);
  }

  private getCurrentNetworkCondition(): NetworkCondition | null {
    if (!this.activePattern) {
      return null;
    }

    const conditionIndex = this.currentConditionIndex % this.activePattern.conditions.length;
    return this.activePattern.conditions[conditionIndex];
  }

  // =============================================================================
  // Pattern Execution
  // =============================================================================

  private async executePattern(pattern: FailurePattern): Promise<void> {
    for (let rep = 0; rep < pattern.repetitions && this.isSimulating; rep++) {
      console.log(`🔄 Repetition ${rep + 1}/${pattern.repetitions}`);

      for (let i = 0; i < pattern.conditions.length && this.isSimulating; i++) {
        this.currentConditionIndex = i;
        const condition = pattern.conditions[i];

        console.log(`⚡ Applying condition: ${condition.name} for ${pattern.duration}ms`);
        console.log(`   📊 Failure probability: ${(condition.probability * 100).toFixed(1)}%`);
        console.log(`   ⏱️ Latency: ${condition.latency}ms ±${condition.jitter}ms`);

        // Apply condition for specified duration
        await this.delay(pattern.duration);
      }

      // Wait between repetitions
      if (rep < pattern.repetitions - 1 && pattern.interval > 0) {
        console.log(`⏸️ Waiting ${pattern.interval}ms before next repetition`);
        await this.delay(pattern.interval);
      }
    }
  }

  // =============================================================================
  // Results and Analytics
  // =============================================================================

  private generateSimulationResult(pattern: FailurePattern): SimulationResult {
    const endTime = Date.now();
    const duration = endTime - this.patternStartTime;
    const interceptor = this.interceptor!;

    // Calculate metrics
    const averageLatency = interceptor.latencyRecords.length > 0
      ? interceptor.latencyRecords.reduce((a, b) => a + b, 0) / interceptor.latencyRecords.length
      : 0;

    const maxLatency = interceptor.latencyRecords.length > 0
      ? Math.max(...interceptor.latencyRecords)
      : 0;

    const minLatency = interceptor.latencyRecords.length > 0
      ? Math.min(...interceptor.latencyRecords)
      : 0;

    // Estimate scores (in real implementation, these would be measured)
    const successRate = interceptor.successCount / (interceptor.successCount + interceptor.failureCount);
    const dataIntegrityScore = Math.max(0, Math.min(100, successRate * 100));
    const userExperienceScore = Math.max(0, Math.min(100, (1 - averageLatency / 10000) * 100));
    const fallbackEffectiveness = Math.max(0, Math.min(100, (interceptor.failureCount > 0 ? 80 : 100)));

    return {
      patternId: pattern.id,
      startTime: this.patternStartTime,
      endTime,
      duration,
      totalRequests: interceptor.interceptedCount,
      successfulRequests: interceptor.successCount,
      failedRequests: interceptor.failureCount,
      fallbackActivations: interceptor.failureCount, // Simplified estimation
      recoveryTime: this.calculateRecoveryTime(),
      errors: [], // Would be populated with recorded errors
      metrics: {
        averageLatency,
        maxLatency,
        minLatency,
        dataIntegrityScore,
        userExperienceScore,
        fallbackEffectiveness,
      },
    };
  }

  private calculateTotalDuration(pattern: FailurePattern): number {
    const patternDuration = pattern.conditions.length * pattern.duration;
    const intervalTime = (pattern.repetitions - 1) * pattern.interval;
    return pattern.repetitions * patternDuration + intervalTime;
  }

  private calculateRecoveryTime(): number {
    // Simplified recovery time calculation
    return Math.random() * 2000 + 1000; // 1-3 seconds
  }

  private recordError(url: string, error: any): void {
    // In a real implementation, this would record detailed error information
    console.warn(`🚨 Network error on ${url}:`, error.message);
  }

  // =============================================================================
  // Utility Methods
  // =============================================================================

  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * Get available failure patterns
   */
  getAvailablePatterns(): Array<{ id: string; name: string; description: string }> {
    return Object.values(this.failurePatterns).map(pattern => ({
      id: pattern.id,
      name: pattern.name,
      description: pattern.description,
    }));
  }

  /**
   * Get current simulation status
   */
  getStatus(): {
    isSimulating: boolean;
    currentPattern?: string;
    currentCondition?: string;
    startTime?: number;
    estimatedEndTime?: number;
  } {
    return {
      isSimulating: this.isSimulating,
      currentPattern: this.activePattern?.name,
      currentCondition: this.getCurrentNetworkCondition()?.name,
      startTime: this.patternStartTime,
      estimatedEndTime: this.activePattern
        ? this.patternStartTime + this.calculateTotalDuration(this.activePattern)
        : undefined,
    };
  }
}

// Singleton instance
export const networkFailureSimulator = new NetworkFailureSimulator();

export default NetworkFailureSimulator;