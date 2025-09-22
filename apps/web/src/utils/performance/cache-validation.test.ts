/**
 * Caching Effectiveness Validation Tests
 * 
 * Comprehensive testing for cache performance and effectiveness validation
 * Ensures cached content loads within performance budgets
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { PerformanceMeasurement, PERFORMANCE_BUDGET, ContentLoadMetrics } from "./measurement";

/**
 * Cache configuration types
 */
export interface CacheConfig {
  strategy: "memory" | "localStorage" | "sessionStorage" | "indexedDB" | "serviceWorker";
  maxSize: number; // bytes
  ttl: number; // milliseconds
  compression: boolean;
  encryption: boolean;
}

/**
 * Cache hit statistics
 */
export interface CacheStats {
  totalRequests: number;
  cacheHits: number;
  cacheMisses: number;
  hitRatio: number;
  averageHitTime: number;
  averageMissTime: number;
  totalSize: number;
  evictions: number;
}

/**
 * Cache validation result
 */
export interface CacheValidationResult {
  strategy: string;
  config: CacheConfig;
  stats: CacheStats;
  performanceMetrics: {
    hitTimeP50: number;
    hitTimeP95: number;
    missTimeP50: number;
    missTimeP95: number;
  };
  meetsBudget: boolean;
  recommendations: string[];
}

/**
 * Mock cache implementation for testing
 */
export class MockCache {
  private storage = new Map<string, { data: any; timestamp: number; size: number; accessCount: number }>();
  private config: CacheConfig;
  private stats: CacheStats;
  private hitTimes: number[] = [];
  private missTimes: number[] = [];

  constructor(config: CacheConfig) {
    this.config = config;
    this.stats = {
      totalRequests: 0,
      cacheHits: 0,
      cacheMisses: 0,
      hitRatio: 0,
      averageHitTime: 0,
      averageMissTime: 0,
      totalSize: 0,
      evictions: 0,
    };
  }

  /**
   * Simulate cache get operation
   */
  async get(key: string): Promise<{ data: any; cacheStatus: "hit" | "miss" | "stale" }> {
    const startTime = performance.now();
    this.stats.totalRequests++;

    const entry = this.storage.get(key);
    
    if (entry) {
      const age = Date.now() - entry.timestamp;
      
      // Check TTL
      if (age > this.config.ttl) {
        // Cache entry is stale
        this.storage.delete(key);
        this.stats.totalSize -= entry.size;
        
        const endTime = performance.now();
        const duration = endTime - startTime + 50; // Add slight delay for stale handling
        this.missTimes.push(duration);
        this.stats.cacheMisses++;
        this.updateStats();
        
        return { data: null, cacheStatus: "stale" };
      }
      
      // Cache hit
      entry.accessCount++;
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Simulate cache access delay based on strategy
      const accessDelay = this.getCacheAccessDelay("read");
      await new Promise(resolve => setTimeout(resolve, accessDelay));
      
      this.hitTimes.push(duration + accessDelay);
      this.stats.cacheHits++;
      this.updateStats();
      
      return { data: entry.data, cacheStatus: "hit" };
    }
    
    // Cache miss
    const endTime = performance.now();
    const duration = endTime - startTime;
    this.missTimes.push(duration);
    this.stats.cacheMisses++;
    this.updateStats();
    
    return { data: null, cacheStatus: "miss" };
  }

  /**
   * Simulate cache set operation
   */
  async set(key: string, data: any): Promise<void> {
    const size = this.calculateSize(data);
    const accessDelay = this.getCacheAccessDelay("write");
    
    // Check cache size limits
    while (this.stats.totalSize + size > this.config.maxSize && this.storage.size > 0) {
      this.evictLeastRecentlyUsed();
    }
    
    // Simulate write delay
    await new Promise(resolve => setTimeout(resolve, accessDelay));
    
    this.storage.set(key, {
      data,
      timestamp: Date.now(),
      size,
      accessCount: 0,
    });
    
    this.stats.totalSize += size;
  }

  /**
   * Get cache access delay based on strategy
   */
  private getCacheAccessDelay(operation: "read" | "write"): number {
    const delays = {
      memory: { read: 0.1, write: 0.2 },
      localStorage: { read: 1, write: 2 },
      sessionStorage: { read: 0.8, write: 1.5 },
      indexedDB: { read: 2, write: 5 },
      serviceWorker: { read: 5, write: 10 },
    };
    
    return delays[this.config.strategy][operation];
  }

  /**
   * Calculate data size
   */
  private calculateSize(data: any): number {
    const json = JSON.stringify(data);
    const size = new Blob([json]).size;
    
    // Apply compression if enabled
    if (this.config.compression) {
      return Math.floor(size * 0.7); // Assume 30% compression
    }
    
    return size;
  }

  /**
   * Evict least recently used entry
   */
  private evictLeastRecentlyUsed(): void {
    let lruKey: string | null = null;
    let lruAccessCount = Infinity;
    
    for (const [key, entry] of this.storage.entries()) {
      if (entry.accessCount < lruAccessCount) {
        lruAccessCount = entry.accessCount;
        lruKey = key;
      }
    }
    
    if (lruKey) {
      const entry = this.storage.get(lruKey)!;
      this.storage.delete(lruKey);
      this.stats.totalSize -= entry.size;
      this.stats.evictions++;
    }
  }

  /**
   * Update cache statistics
   */
  private updateStats(): void {
    this.stats.hitRatio = this.stats.totalRequests > 0 ? 
      this.stats.cacheHits / this.stats.totalRequests : 0;
    
    this.stats.averageHitTime = this.hitTimes.length > 0 ?
      this.hitTimes.reduce((a, b) => a + b, 0) / this.hitTimes.length : 0;
    
    this.stats.averageMissTime = this.missTimes.length > 0 ?
      this.missTimes.reduce((a, b) => a + b, 0) / this.missTimes.length : 0;
  }

  /**
   * Get current statistics
   */
  getStats(): CacheStats {
    return { ...this.stats };
  }

  /**
   * Get performance percentiles
   */
  getPerformanceMetrics() {
    const sortedHitTimes = [...this.hitTimes].sort((a, b) => a - b);
    const sortedMissTimes = [...this.missTimes].sort((a, b) => a - b);
    
    return {
      hitTimeP50: this.getPercentile(sortedHitTimes, 50),
      hitTimeP95: this.getPercentile(sortedHitTimes, 95),
      missTimeP50: this.getPercentile(sortedMissTimes, 50),
      missTimeP95: this.getPercentile(sortedMissTimes, 95),
    };
  }

  /**
   * Calculate percentile
   */
  private getPercentile(sortedArray: number[], percentile: number): number {
    if (sortedArray.length === 0) return 0;
    
    const index = Math.floor((percentile / 100) * sortedArray.length);
    return sortedArray[Math.min(index, sortedArray.length - 1)];
  }

  /**
   * Clear cache
   */
  clear(): void {
    this.storage.clear();
    this.stats.totalSize = 0;
    this.hitTimes = [];
    this.missTimes = [];
  }
}

/**
 * Cache validation utility
 */
export class CacheValidator {
  private caches = new Map<string, MockCache>();
  private performanceMeasurement: PerformanceMeasurement;

  constructor() {
    this.performanceMeasurement = new PerformanceMeasurement();
  }

  /**
   * Create cache with configuration
   */
  createCache(name: string, config: CacheConfig): MockCache {
    const cache = new MockCache(config);
    this.caches.set(name, cache);
    return cache;
  }

  /**
   * Validate cache effectiveness
   */
  async validateCacheEffectiveness(
    cacheName: string,
    contentItems: Array<{ id: string; data: any; type: string; category: string }>
  ): Promise<CacheValidationResult> {
    const cache = this.caches.get(cacheName);
    if (!cache) {
      throw new Error(`Cache ${cacheName} not found`);
    }

    const config = (cache as any).config;
    
    // Simulate real-world usage pattern
    await this.simulateUsagePattern(cache, contentItems);
    
    const stats = cache.getStats();
    const performanceMetrics = cache.getPerformanceMetrics();
    
    // Check if cache meets performance budget
    const meetsBudget = this.evaluatePerformanceBudget(stats, performanceMetrics);
    
    // Generate recommendations
    const recommendations = this.generateRecommendations(config, stats, performanceMetrics);
    
    return {
      strategy: config.strategy,
      config,
      stats,
      performanceMetrics,
      meetsBudget,
      recommendations,
    };
  }

  /**
   * Simulate realistic usage pattern
   */
  private async simulateUsagePattern(
    cache: MockCache,
    contentItems: Array<{ id: string; data: any; type: string; category: string }>
  ): Promise<void> {
    // Phase 1: Initial cache population (cold start)
    for (const item of contentItems) {
      await cache.set(item.id, item.data);
    }

    // Phase 2: Mixed read/write pattern (80% reads, 20% writes)
    const operations = 100;
    
    for (let i = 0; i < operations; i++) {
      const isRead = Math.random() < 0.8;
      
      if (isRead) {
        // Read operation - simulate content access patterns
        const item = this.selectItemByAccessPattern(contentItems);
        const contentId = `cache-test-${item.id}-${i}`;
        
        this.performanceMeasurement.startContentLoad(contentId, item.type, item.category);
        
        const result = await cache.get(item.id);
        
        // If cache miss, simulate fetch and cache set
        if (result.cacheStatus === "miss" || result.cacheStatus === "stale") {
          // Simulate network fetch delay
          await new Promise(resolve => setTimeout(resolve, 50));
          await cache.set(item.id, item.data);
        }
        
        this.performanceMeasurement.endContentLoad(
          contentId, 
          JSON.stringify(item.data).length,
          result.cacheStatus
        );
      } else {
        // Write operation - update existing content
        const item = contentItems[Math.floor(Math.random() * contentItems.length)];
        const updatedData = { ...item.data, updated: Date.now() };
        await cache.set(item.id, updatedData);
      }
    }
  }

  /**
   * Select item based on realistic access patterns (80/20 rule)
   */
  private selectItemByAccessPattern(items: any[]): any {
    // 80% of accesses go to 20% of content (popular content)
    const isPopularAccess = Math.random() < 0.8;
    
    if (isPopularAccess) {
      // Access popular content (first 20% of items)
      const popularCount = Math.max(1, Math.floor(items.length * 0.2));
      return items[Math.floor(Math.random() * popularCount)];
    } else {
      // Access less popular content
      return items[Math.floor(Math.random() * items.length)];
    }
  }

  /**
   * Evaluate if cache meets performance budget
   */
  private evaluatePerformanceBudget(stats: CacheStats, metrics: any): boolean {
    // Cache hits should be very fast
    const hitTimeBudget = 10; // 10ms for cache hits
    const hitTimeOk = metrics.hitTimeP95 <= hitTimeBudget;
    
    // Hit ratio should be high
    const hitRatioOk = stats.hitRatio >= PERFORMANCE_BUDGET.minCacheHitRatio;
    
    // Average hit time should be much faster than miss time
    const speedupOk = stats.averageMissTime === 0 || 
      (stats.averageHitTime / stats.averageMissTime) < 0.2; // 80% faster
    
    return hitTimeOk && hitRatioOk && speedupOk;
  }

  /**
   * Generate optimization recommendations
   */
  private generateRecommendations(
    config: CacheConfig,
    stats: CacheStats,
    metrics: any
  ): string[] {
    const recommendations: string[] = [];

    // Hit ratio recommendations
    if (stats.hitRatio < PERFORMANCE_BUDGET.minCacheHitRatio) {
      recommendations.push(`Increase cache size or TTL (current hit ratio: ${(stats.hitRatio * 100).toFixed(1)}%)`);
    }

    // Performance recommendations
    if (metrics.hitTimeP95 > 10) {
      recommendations.push(`Optimize cache strategy for faster access (P95 hit time: ${metrics.hitTimeP95.toFixed(1)}ms)`);
      
      if (config.strategy === "indexedDB" || config.strategy === "serviceWorker") {
        recommendations.push("Consider memory or localStorage for better performance");
      }
    }

    // Size recommendations
    if (stats.evictions > stats.totalRequests * 0.1) {
      recommendations.push(`Increase cache size to reduce evictions (${stats.evictions} evictions)`);
    }

    // Compression recommendations
    if (!config.compression && config.maxSize < 1000000) {
      recommendations.push("Enable compression to store more content");
    }

    // Strategy recommendations
    if (config.strategy === "sessionStorage" && stats.totalSize > 5000000) {
      recommendations.push("Consider IndexedDB for larger cache sizes");
    }

    return recommendations;
  }

  /**
   * Get all validation results
   */
  getAllResults(): CacheValidationResult[] {
    // Implementation would return all cached validation results
    return [];
  }

  /**
   * Cleanup
   */
  cleanup(): void {
    this.performanceMeasurement.disconnect();
    this.caches.clear();
  }
}

// Test content for cache validation
const testContent = [
  { id: "text-1", data: "Simple text content", type: "text", category: "ui" },
  { id: "text-2", data: "Medium length text content".repeat(10), type: "text", category: "ui" },
  { id: "config-1", data: { theme: "dark", language: "en" }, type: "config", category: "settings" },
  { id: "config-2", data: { api: { baseUrl: "https://api.example.com", timeout: 5000 } }, type: "config", category: "api" },
  { id: "rich-text-1", data: "<div><h1>Title</h1><p>Content</p></div>", type: "richText", category: "content" },
  { id: "translation-1", data: { en: "Hello", es: "Hola", fr: "Bonjour" }, type: "translation", category: "i18n" },
];

describe("Caching Effectiveness Validation Tests", () => {
  let validator: CacheValidator;

  beforeEach(() => {
    validator = new CacheValidator();
  });

  afterEach(() => {
    validator.cleanup();
  });

  describe("Cache Strategy Performance", () => {
    const cacheConfigs: Array<{ name: string; config: CacheConfig }> = [
      {
        name: "memory-fast",
        config: {
          strategy: "memory",
          maxSize: 1000000, // 1MB
          ttl: 300000, // 5 minutes
          compression: false,
          encryption: false,
        },
      },
      {
        name: "localStorage-balanced",
        config: {
          strategy: "localStorage",
          maxSize: 5000000, // 5MB
          ttl: 600000, // 10 minutes
          compression: true,
          encryption: false,
        },
      },
      {
        name: "indexedDB-large",
        config: {
          strategy: "indexedDB",
          maxSize: 50000000, // 50MB
          ttl: 1800000, // 30 minutes
          compression: true,
          encryption: true,
        },
      },
    ];

    cacheConfigs.forEach(({ name, config }) => {
      it(`should validate ${name} cache strategy effectiveness`, async () => {
        validator.createCache(name, config);
        
        const result = await validator.validateCacheEffectiveness(name, testContent);
        
        expect(result).toBeTruthy();
        expect(result.strategy).toBe(config.strategy);
        expect(result.stats.totalRequests).toBeGreaterThan(0);
        expect(result.stats.hitRatio).toBeGreaterThanOrEqual(0);
        expect(result.stats.hitRatio).toBeLessThanOrEqual(1);
        
        // Log results for analysis
        console.log(`${name} Results:
          Hit Ratio: ${(result.stats.hitRatio * 100).toFixed(1)}%
          Avg Hit Time: ${result.stats.averageHitTime.toFixed(2)}ms
          Avg Miss Time: ${result.stats.averageMissTime.toFixed(2)}ms
          P95 Hit Time: ${result.performanceMetrics.hitTimeP95.toFixed(2)}ms
          Meets Budget: ${result.meetsBudget}
          Recommendations: ${result.recommendations.length}`);
        
        // Memory cache should have the best performance
        if (config.strategy === "memory") {
          expect(result.performanceMetrics.hitTimeP95).toBeLessThanOrEqual(5);
        }
        
        // All caches should meet minimum hit ratio after warm-up
        expect(result.stats.hitRatio).toBeGreaterThan(0.5); // 50% minimum after usage pattern
      });
    });
  });

  describe("Cache Hit Ratio Optimization", () => {
    it("should achieve high hit ratios with proper configuration", async () => {
      const optimizedConfig: CacheConfig = {
        strategy: "memory",
        maxSize: 2000000, // 2MB - larger cache
        ttl: 900000, // 15 minutes - longer TTL
        compression: true,
        encryption: false,
      };

      validator.createCache("optimized", optimizedConfig);
      
      const result = await validator.validateCacheEffectiveness("optimized", testContent);
      
      expect(result.stats.hitRatio).toBeGreaterThanOrEqual(PERFORMANCE_BUDGET.minCacheHitRatio);
      expect(result.meetsBudget).toBe(true);
      
      // Should have minimal evictions
      expect(result.stats.evictions).toBeLessThanOrEqual(result.stats.totalRequests * 0.05);
    });

    it("should detect poor hit ratios and provide recommendations", async () => {
      const poorConfig: CacheConfig = {
        strategy: "memory",
        maxSize: 1000, // Very small cache
        ttl: 1000, // Very short TTL
        compression: false,
        encryption: false,
      };

      validator.createCache("poor", poorConfig);
      
      const result = await validator.validateCacheEffectiveness("poor", testContent);
      
      expect(result.stats.hitRatio).toBeLessThan(PERFORMANCE_BUDGET.minCacheHitRatio);
      expect(result.meetsBudget).toBe(false);
      expect(result.recommendations.length).toBeGreaterThan(0);
      
      // Should recommend increasing cache size and TTL
      const hasRecommendations = result.recommendations.some(r => 
        r.includes("cache size") || r.includes("TTL")
      );
      expect(hasRecommendations).toBe(true);
    });
  });

  describe("Cache Performance Budgets", () => {
    it("should enforce strict performance budgets for cache operations", async () => {
      const config: CacheConfig = {
        strategy: "memory",
        maxSize: 1000000,
        ttl: 600000,
        compression: false,
        encryption: false,
      };

      validator.createCache("budget-test", config);
      
      const result = await validator.validateCacheEffectiveness("budget-test", testContent);
      
      // Cache hits should be very fast
      expect(result.performanceMetrics.hitTimeP95).toBeLessThanOrEqual(10);
      
      // Hit ratio should meet minimum requirement
      expect(result.stats.hitRatio).toBeGreaterThanOrEqual(PERFORMANCE_BUDGET.minCacheHitRatio);
      
      // Cache hits should be significantly faster than misses
      if (result.stats.averageMissTime > 0) {
        const speedup = result.stats.averageMissTime / result.stats.averageHitTime;
        expect(speedup).toBeGreaterThan(2); // At least 2x faster
      }
    });
  });

  describe("Cache Strategy Comparison", () => {
    it("should compare different cache strategies", async () => {
      const strategies = ["memory", "localStorage", "indexedDB"] as const;
      const results: CacheValidationResult[] = [];

      for (const strategy of strategies) {
        const config: CacheConfig = {
          strategy,
          maxSize: 1000000,
          ttl: 600000,
          compression: false,
          encryption: false,
        };

        validator.createCache(strategy, config);
        const result = await validator.validateCacheEffectiveness(strategy, testContent);
        results.push(result);
      }

      // Memory should be fastest
      const memoryResult = results.find(r => r.strategy === "memory")!;
      const localStorageResult = results.find(r => r.strategy === "localStorage")!;
      const indexedDBResult = results.find(r => r.strategy === "indexedDB")!;

      expect(memoryResult.performanceMetrics.hitTimeP95)
        .toBeLessThan(localStorageResult.performanceMetrics.hitTimeP95);
      
      expect(localStorageResult.performanceMetrics.hitTimeP95)
        .toBeLessThan(indexedDBResult.performanceMetrics.hitTimeP95);

      // All should meet basic requirements
      results.forEach(result => {
        expect(result.stats.hitRatio).toBeGreaterThan(0);
        expect(result.stats.totalRequests).toBeGreaterThan(0);
      });
    });
  });

  describe("Content Type Caching Patterns", () => {
    it("should optimize caching for different content types", async () => {
      const contentTypes = ["text", "config", "richText", "translation"];
      const results = new Map<string, { hitRatio: number; avgHitTime: number }>();

      for (const contentType of contentTypes) {
        const filteredContent = testContent.filter(item => item.type === contentType);
        
        if (filteredContent.length === 0) continue;

        const config: CacheConfig = {
          strategy: "memory",
          maxSize: 1000000,
          ttl: 600000,
          compression: contentType === "richText", // Compress rich text
          encryption: contentType === "config", // Encrypt sensitive config
        };

        const cacheName = `${contentType}-cache`;
        validator.createCache(cacheName, config);
        
        const result = await validator.validateCacheEffectiveness(cacheName, filteredContent);
        
        results.set(contentType, {
          hitRatio: result.stats.hitRatio,
          avgHitTime: result.stats.averageHitTime,
        });
      }

      // All content types should achieve reasonable hit ratios
      results.forEach((metrics, contentType) => {
        expect(metrics.hitRatio).toBeGreaterThan(0.3); // 30% minimum
        expect(metrics.avgHitTime).toBeLessThanOrEqual(10); // 10ms maximum
        
        console.log(`${contentType}: Hit ratio ${(metrics.hitRatio * 100).toFixed(1)}%, Avg hit time ${metrics.avgHitTime.toFixed(2)}ms`);
      });
    });
  });

  describe("Cache Eviction and Memory Management", () => {
    it("should handle cache eviction efficiently", async () => {
      const smallCacheConfig: CacheConfig = {
        strategy: "memory",
        maxSize: 5000, // Very small cache to force evictions
        ttl: 600000,
        compression: false,
        encryption: false,
      };

      validator.createCache("eviction-test", smallCacheConfig);
      
      const result = await validator.validateCacheEffectiveness("eviction-test", testContent);
      
      // Should have evictions due to small cache size
      expect(result.stats.evictions).toBeGreaterThan(0);
      
      // Despite evictions, should still maintain some hit ratio
      expect(result.stats.hitRatio).toBeGreaterThan(0.2);
      
      // Should recommend increasing cache size
      const hasEvictionRecommendation = result.recommendations.some(r => 
        r.includes("eviction") || r.includes("cache size")
      );
      expect(hasEvictionRecommendation).toBe(true);
    });
  });

  describe("TTL and Staleness Handling", () => {
    it("should handle TTL expiration correctly", async () => {
      const shortTTLConfig: CacheConfig = {
        strategy: "memory",
        maxSize: 1000000,
        ttl: 100, // Very short TTL
        compression: false,
        encryption: false,
      };

      const cache = validator.createCache("ttl-test", shortTTLConfig);
      
      // Set some data
      await cache.set("test-key", { data: "test" });
      
      // Immediate access should hit
      let result = await cache.get("test-key");
      expect(result.cacheStatus).toBe("hit");
      
      // Wait for TTL expiration
      await new Promise(resolve => setTimeout(resolve, 150));
      
      // Should now be stale/miss
      result = await cache.get("test-key");
      expect(result.cacheStatus).toBe("stale");
    });
  });
});
