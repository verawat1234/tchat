/**
 * Content Prefetching Strategy Hook
 *
 * Implements intelligent content prefetching with smart prediction, priority management,
 * and performance optimization. Provides measurable improvements to user experience
 * while respecting bandwidth and memory constraints.
 *
 * Features:
 * - Smart prediction based on user behavior patterns
 * - Route-based prefetching for upcoming navigation
 * - Category bulk prefetching for related content
 * - Performance budgets and network condition awareness
 * - Priority queue system based on access likelihood
 * - Background processing with non-blocking operations
 * - Analytics integration for effectiveness tracking
 */

import { useEffect, useRef, useCallback, useMemo, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { 
  useLazyGetContentItemQuery,
  useLazyGetContentByCategoryQuery,
  useLazyGetContentItemsQuery,
  contentApi 
} from '../services/contentApi';
import type { RootState } from '../store';

// =============================================================================
// Type Definitions
// =============================================================================

/**
 * Network condition assessment for adaptive prefetching
 */
export interface NetworkCondition {
  /** Connection type (4g, 3g, 2g, wifi, unknown) */
  effectiveType: string;
  /** Round-trip time in milliseconds */
  rtt: number;
  /** Downlink speed in Mbps */
  downlink: number;
  /** Data saver mode enabled */
  saveData: boolean;
}

/**
 * Performance budget configuration
 */
export interface PerformanceBudget {
  /** Maximum memory usage in MB */
  maxMemoryMB: number;
  /** Maximum prefetch operations per minute */
  maxOperationsPerMinute: number;
  /** Maximum concurrent prefetch requests */
  maxConcurrentRequests: number;
  /** Maximum prefetch data per session in MB */
  maxDataPerSessionMB: number;
}

/**
 * Prefetch priority levels
 */
export enum PrefetchPriority {
  /** Critical for immediate use */
  CRITICAL = 0,
  /** High likelihood of access */
  HIGH = 1,
  /** Medium likelihood of access */
  MEDIUM = 2,
  /** Low likelihood of access */
  LOW = 3,
  /** Background prefetch when idle */
  IDLE = 4
}

/**
 * Prefetch strategy types
 */
export enum PrefetchStrategy {
  /** Prefetch based on route prediction */
  ROUTE_BASED = 'route_based',
  /** Prefetch entire content categories */
  CATEGORY_BULK = 'category_bulk',
  /** Prefetch based on user behavior patterns */
  BEHAVIOR_BASED = 'behavior_based',
  /** Prefetch adjacent/related content */
  ADJACENT_CONTENT = 'adjacent_content',
  /** Prefetch for offline support */
  OFFLINE_SUPPORT = 'offline_support'
}

/**
 * Prefetch task definition
 */
export interface PrefetchTask {
  /** Unique task identifier */
  id: string;
  /** Content or category to prefetch */
  target: string;
  /** Prefetch strategy type */
  strategy: PrefetchStrategy;
  /** Task priority level */
  priority: PrefetchPriority;
  /** Estimated access probability (0-1) */
  probability: number;
  /** Estimated data size in KB */
  estimatedSizeKB: number;
  /** Task creation timestamp */
  createdAt: number;
  /** Retry count for failed tasks */
  retryCount: number;
  /** Maximum retry attempts */
  maxRetries: number;
}

/**
 * User behavior pattern for prediction
 */
export interface BehaviorPattern {
  /** Route or content access patterns */
  accessPath: string[];
  /** Time spent on each step */
  dwellTimes: number[];
  /** Access frequency */
  frequency: number;
  /** Last access timestamp */
  lastAccess: number;
  /** Success rate of predictions */
  predictionAccuracy: number;
}

/**
 * Prefetch analytics data
 */
export interface PrefetchAnalytics {
  /** Total prefetch operations */
  totalOperations: number;
  /** Successful prefetch hits */
  successfulHits: number;
  /** Cache hit rate percentage */
  hitRate: number;
  /** Average response time improvement */
  avgResponseImprovement: number;
  /** Data usage metrics */
  dataUsage: {
    total: number;
    successful: number;
    wasted: number;
  };
  /** Strategy effectiveness */
  strategyStats: Record<PrefetchStrategy, {
    operations: number;
    hits: number;
    hitRate: number;
  }>;
}

/**
 * Prefetch configuration options
 */
export interface PrefetchConfig {
  /** Enable/disable prefetching */
  enabled: boolean;
  /** Performance budget settings */
  budget: PerformanceBudget;
  /** Strategies to enable */
  enabledStrategies: PrefetchStrategy[];
  /** Minimum probability threshold for prefetching */
  minProbabilityThreshold: number;
  /** Maximum task queue size */
  maxQueueSize: number;
  /** Background processing interval in ms */
  processingIntervalMs: number;
  /** Analytics collection enabled */
  analyticsEnabled: boolean;
}

// =============================================================================
// Default Configuration
// =============================================================================

const DEFAULT_CONFIG: PrefetchConfig = {
  enabled: true,
  budget: {
    maxMemoryMB: 50,
    maxOperationsPerMinute: 30,
    maxConcurrentRequests: 3,
    maxDataPerSessionMB: 10,
  },
  enabledStrategies: [
    PrefetchStrategy.ROUTE_BASED,
    PrefetchStrategy.CATEGORY_BULK,
    PrefetchStrategy.BEHAVIOR_BASED,
    PrefetchStrategy.ADJACENT_CONTENT,
  ],
  minProbabilityThreshold: 0.3,
  maxQueueSize: 100,
  processingIntervalMs: 2000,
  analyticsEnabled: true,
};

// =============================================================================
// Network Condition Detection
// =============================================================================

/**
 * Get current network condition information
 */
const getNetworkCondition = (): NetworkCondition => {
  const connection = (navigator as any).connection || (navigator as any).mozConnection || (navigator as any).webkitConnection;
  
  return {
    effectiveType: connection?.effectiveType || 'unknown',
    rtt: connection?.rtt || 100,
    downlink: connection?.downlink || 10,
    saveData: connection?.saveData || false,
  };
};

/**
 * Determine if network conditions are suitable for prefetching
 */
const isNetworkSuitableForPrefetching = (condition: NetworkCondition): boolean => {
  // Don't prefetch in data saver mode
  if (condition.saveData) return false;
  
  // Limit prefetching on slow connections
  if (condition.effectiveType === '2g' || condition.rtt > 2000) return false;
  
  // Allow prefetching on good connections
  return true;
};

// =============================================================================
// Priority Queue Implementation
// =============================================================================

class PrefetchQueue {
  private queue: PrefetchTask[] = [];
  private maxSize: number;

  constructor(maxSize: number = 100) {
    this.maxSize = maxSize;
  }

  /**
   * Add task to queue with priority ordering
   */
  enqueue(task: PrefetchTask): void {
    // Remove existing task for same target if present
    this.queue = this.queue.filter(t => t.target !== task.target);
    
    // Insert task in priority order
    let insertIndex = 0;
    while (insertIndex < this.queue.length && 
           this.comparePriority(this.queue[insertIndex], task) <= 0) {
      insertIndex++;
    }
    
    this.queue.splice(insertIndex, 0, task);
    
    // Maintain queue size limit
    if (this.queue.length > this.maxSize) {
      this.queue = this.queue.slice(0, this.maxSize);
    }
  }

  /**
   * Get next highest priority task
   */
  dequeue(): PrefetchTask | undefined {
    return this.queue.shift();
  }

  /**
   * Get all tasks with given priority or higher
   */
  getTasksByPriority(maxPriority: PrefetchPriority): PrefetchTask[] {
    return this.queue.filter(task => task.priority <= maxPriority);
  }

  /**
   * Compare task priorities (lower number = higher priority)
   */
  private comparePriority(a: PrefetchTask, b: PrefetchTask): number {
    // First by priority level
    if (a.priority !== b.priority) {
      return a.priority - b.priority;
    }
    
    // Then by probability (higher = better)
    if (Math.abs(a.probability - b.probability) > 0.1) {
      return b.probability - a.probability;
    }
    
    // Finally by creation time (older = higher priority)
    return a.createdAt - b.createdAt;
  }

  /**
   * Get queue size
   */
  size(): number {
    return this.queue.length;
  }

  /**
   * Clear the queue
   */
  clear(): void {
    this.queue = [];
  }

  /**
   * Get all tasks (for debugging)
   */
  getAllTasks(): PrefetchTask[] {
    return [...this.queue];
  }
}

// =============================================================================
// Behavior Pattern Analysis
// =============================================================================

class BehaviorAnalyzer {
  private patterns: Map<string, BehaviorPattern> = new Map();
  private accessHistory: Array<{ path: string; timestamp: number }> = [];
  private maxHistorySize = 1000;

  /**
   * Record a route/content access
   */
  recordAccess(path: string): void {
    const timestamp = Date.now();
    
    // Add to history
    this.accessHistory.push({ path, timestamp });
    
    // Maintain history size
    if (this.accessHistory.length > this.maxHistorySize) {
      this.accessHistory = this.accessHistory.slice(-this.maxHistorySize);
    }
    
    // Update patterns
    this.updatePatterns(path, timestamp);
  }

  /**
   * Predict next likely access targets
   */
  predictNextAccess(currentPath: string): Array<{ target: string; probability: number }> {
    const predictions: Array<{ target: string; probability: number }> = [];
    
    // Analyze historical patterns
    const recentHistory = this.accessHistory.slice(-50); // Last 50 accesses
    const currentIndex = recentHistory.findLastIndex(h => h.path === currentPath);
    
    if (currentIndex >= 0 && currentIndex < recentHistory.length - 1) {
      // Find what typically comes after current path
      const followingPaths: Record<string, number> = {};
      
      for (let i = 0; i < recentHistory.length - 1; i++) {
        if (recentHistory[i].path === currentPath) {
          const nextPath = recentHistory[i + 1].path;
          followingPaths[nextPath] = (followingPaths[nextPath] || 0) + 1;
        }
      }
      
      // Convert to probabilities
      const totalFollowing = Object.values(followingPaths).reduce((sum, count) => sum + count, 0);
      
      for (const [path, count] of Object.entries(followingPaths)) {
        predictions.push({
          target: path,
          probability: count / totalFollowing,
        });
      }
    }
    
    // Sort by probability and return top predictions
    return predictions
      .sort((a, b) => b.probability - a.probability)
      .slice(0, 5);
  }

  /**
   * Update behavior patterns
   */
  private updatePatterns(path: string, timestamp: number): void {
    const existing = this.patterns.get(path);
    
    if (existing) {
      existing.frequency += 1;
      existing.lastAccess = timestamp;
      // Update prediction accuracy based on whether this was predicted
      // (This would be enhanced with actual prediction tracking)
    } else {
      this.patterns.set(path, {
        accessPath: [path],
        dwellTimes: [],
        frequency: 1,
        lastAccess: timestamp,
        predictionAccuracy: 0.5, // Start with neutral accuracy
      });
    }
  }

  /**
   * Get pattern data for analytics
   */
  getPatternStats(): { totalPatterns: number; avgAccuracy: number } {
    const patterns = Array.from(this.patterns.values());
    const avgAccuracy = patterns.length > 0 
      ? patterns.reduce((sum, p) => sum + p.predictionAccuracy, 0) / patterns.length
      : 0;
    
    return {
      totalPatterns: patterns.length,
      avgAccuracy,
    };
  }
}

// =============================================================================
// Main Hook Implementation
// =============================================================================

/**
 * Content prefetching hook with intelligent prediction and performance optimization
 */
export const useContentPrefetch = (config: Partial<PrefetchConfig> = {}) => {
  const finalConfig = useMemo(() => ({ ...DEFAULT_CONFIG, ...config }), [config]);
  
  // Redux state and hooks
  const dispatch = useDispatch();
  
  
  
  // API hooks for prefetching
  const [prefetchContentItem] = useLazyGetContentItemQuery();
  const [prefetchContentByCategory] = useLazyGetContentByCategoryQuery();
  const [prefetchContentItems] = useLazyGetContentItemsQuery();
  
  // Internal state management
  const queueRef = useRef(new PrefetchQueue(finalConfig.maxQueueSize));
  const behaviorAnalyzerRef = useRef(new BehaviorAnalyzer());
  const analyticsRef = useRef<PrefetchAnalytics>({
    totalOperations: 0,
    successfulHits: 0,
    hitRate: 0,
    avgResponseImprovement: 0,
    dataUsage: { total: 0, successful: 0, wasted: 0 },
    strategyStats: Object.values(PrefetchStrategy).reduce((acc, strategy) => {
      acc[strategy] = { operations: 0, hits: 0, hitRate: 0 };
      return acc;
    }, {} as any),
  });
  
  const processingRef = useRef<NodeJS.Timeout | null>(null);
  const activeRequestsRef = useRef(new Set<string>());
  const sessionDataUsageRef = useRef(0);
  const operationCountRef = useRef({ count: 0, windowStart: Date.now() });

  // =============================================================================
  // Core Prefetching Functions
  // =============================================================================

  /**
   * Execute a prefetch task
   */
  const executePrefetchTask = useCallback(async (task: PrefetchTask): Promise<boolean> => {
    if (!finalConfig.enabled) return false;
    
    // Check if already processing this target
    if (activeRequestsRef.current.has(task.target)) return false;
    
    // Check network conditions
    const networkCondition = getNetworkCondition();
    if (!isNetworkSuitableForPrefetching(networkCondition)) return false;
    
    // Check performance budgets
    if (activeRequestsRef.current.size >= finalConfig.budget.maxConcurrentRequests) return false;
    if (sessionDataUsageRef.current >= finalConfig.budget.maxDataPerSessionMB * 1024) return false;
    
    // Check rate limits
    const now = Date.now();
    const windowDuration = 60000; // 1 minute
    if (now - operationCountRef.current.windowStart > windowDuration) {
      operationCountRef.current = { count: 0, windowStart: now };
    }
    if (operationCountRef.current.count >= finalConfig.budget.maxOperationsPerMinute) return false;
    
    try {
      activeRequestsRef.current.add(task.target);
      operationCountRef.current.count++;
      
      let success = false;
      const startTime = performance.now();
      
      // Execute based on strategy
      switch (task.strategy) {
        case PrefetchStrategy.ROUTE_BASED:
        case PrefetchStrategy.BEHAVIOR_BASED:
        case PrefetchStrategy.ADJACENT_CONTENT:
          // Prefetch individual content item
          const itemResult = await prefetchContentItem(task.target);
          success = !itemResult.isError;
          break;
          
        case PrefetchStrategy.CATEGORY_BULK:
          // Prefetch entire category
          const categoryResult = await prefetchContentByCategory(task.target);
          success = !categoryResult.isError;
          break;
          
        case PrefetchStrategy.OFFLINE_SUPPORT:
          // Prefetch with specific parameters for offline support
          const offlineResult = await prefetchContentItems({ 
            category: task.target, 
            limit: 50,
            status: 'published' 
          });
          success = !offlineResult.isError;
          break;
      }
      
      const endTime = performance.now();
      const responseTime = endTime - startTime;
      
      // Update analytics
      if (finalConfig.analyticsEnabled) {
        analyticsRef.current.totalOperations++;
        if (success) {
          analyticsRef.current.successfulHits++;
          analyticsRef.current.dataUsage.successful += task.estimatedSizeKB;
          sessionDataUsageRef.current += task.estimatedSizeKB;
        } else {
          analyticsRef.current.dataUsage.wasted += task.estimatedSizeKB;
        }
        
        // Update strategy stats
        const strategyStats = analyticsRef.current.strategyStats[task.strategy];
        strategyStats.operations++;
        if (success) {
          strategyStats.hits++;
        }
        strategyStats.hitRate = strategyStats.hits / strategyStats.operations;
        
        // Update overall hit rate
        analyticsRef.current.hitRate = 
          analyticsRef.current.successfulHits / analyticsRef.current.totalOperations;
      }
      
      return success;
      
    } catch (error) {
      console.warn('Prefetch task failed:', task.target, error);
      return false;
    } finally {
      activeRequestsRef.current.delete(task.target);
    }
  }, [finalConfig, prefetchContentItem, prefetchContentByCategory, prefetchContentItems]);

  /**
   * Process the prefetch queue
   */
  const processQueue = useCallback(async () => {
    if (!finalConfig.enabled) return;
    
    const queue = queueRef.current;
    const networkCondition = getNetworkCondition();
    
    // Determine how many tasks to process based on network and current load
    let maxTasksToProcess = 1;
    if (networkCondition.effectiveType === 'wifi' || networkCondition.effectiveType === '4g') {
      maxTasksToProcess = Math.min(3, finalConfig.budget.maxConcurrentRequests - activeRequestsRef.current.size);
    }
    
    const tasksToProcess: PrefetchTask[] = [];
    
    // Get tasks based on current conditions
    if (networkCondition.effectiveType === 'wifi') {
      // On WiFi, process all high-priority tasks
      tasksToProcess.push(...queue.getTasksByPriority(PrefetchPriority.HIGH).slice(0, maxTasksToProcess));
    } else if (networkCondition.effectiveType === '4g') {
      // On 4G, focus on critical and high-priority tasks
      tasksToProcess.push(...queue.getTasksByPriority(PrefetchPriority.HIGH).slice(0, maxTasksToProcess));
    } else {
      // On slower connections, only critical tasks
      tasksToProcess.push(...queue.getTasksByPriority(PrefetchPriority.CRITICAL).slice(0, maxTasksToProcess));
    }
    
    // Execute tasks
    for (const task of tasksToProcess) {
      const success = await executePrefetchTask(task);
      
      if (!success && task.retryCount < task.maxRetries) {
        // Re-queue failed task with lower priority
        const retryTask: PrefetchTask = {
          ...task,
          retryCount: task.retryCount + 1,
          priority: Math.min(task.priority + 1, PrefetchPriority.IDLE) as PrefetchPriority,
          probability: task.probability * 0.8, // Reduce probability on retry
        };
        queue.enqueue(retryTask);
      }
    }
  }, [finalConfig, executePrefetchTask]);

  // =============================================================================
  // Prefetch Strategy Implementations
  // =============================================================================

  /**
   * Route-based prefetching: predict and prefetch likely next routes
   */
  const prefetchRouteBasedContent = useCallback((currentRoute: string) => {
    if (!finalConfig.enabledStrategies.includes(PrefetchStrategy.ROUTE_BASED)) return;
    
    // Analyze route patterns and predict next routes
    const predictions = behaviorAnalyzerRef.current.predictNextAccess(currentRoute);
    
    predictions.forEach(({ target, probability }) => {
      if (probability >= finalConfig.minProbabilityThreshold) {
        const task: PrefetchTask = {
          id: `route-${target}-${Date.now()}`,
          target,
          strategy: PrefetchStrategy.ROUTE_BASED,
          priority: probability > 0.7 ? PrefetchPriority.HIGH : PrefetchPriority.MEDIUM,
          probability,
          estimatedSizeKB: 5, // Estimate for route content
          createdAt: Date.now(),
          retryCount: 0,
          maxRetries: 2,
        };
        
        queueRef.current.enqueue(task);
      }
    });
  }, [finalConfig]);

  /**
   * Category bulk prefetching: prefetch all content in related categories
   */
  const prefetchCategoryContent = useCallback((category: string, priority: PrefetchPriority = PrefetchPriority.MEDIUM) => {
    if (!finalConfig.enabledStrategies.includes(PrefetchStrategy.CATEGORY_BULK)) return;
    
    const task: PrefetchTask = {
      id: `category-${category}-${Date.now()}`,
      target: category,
      strategy: PrefetchStrategy.CATEGORY_BULK,
      priority,
      probability: 0.8, // High probability for category content
      estimatedSizeKB: 50, // Estimate for category content
      createdAt: Date.now(),
      retryCount: 0,
      maxRetries: 1,
    };
    
    queueRef.current.enqueue(task);
  }, [finalConfig]);

  /**
   * Behavior-based prefetching: use ML-style pattern recognition
   */
  const prefetchBehaviorBasedContent = useCallback(() => {
    if (!finalConfig.enabledStrategies.includes(PrefetchStrategy.BEHAVIOR_BASED)) return;
    
    // This would be enhanced with more sophisticated ML algorithms
    const currentPath = window.location.pathname || "/";
    const predictions = behaviorAnalyzerRef.current.predictNextAccess(currentPath);
    
    predictions.slice(0, 3).forEach(({ target, probability }) => {
      if (probability >= finalConfig.minProbabilityThreshold) {
        const task: PrefetchTask = {
          id: `behavior-${target}-${Date.now()}`,
          target,
          strategy: PrefetchStrategy.BEHAVIOR_BASED,
          priority: probability > 0.6 ? PrefetchPriority.HIGH : PrefetchPriority.MEDIUM,
          probability,
          estimatedSizeKB: 10,
          createdAt: Date.now(),
          retryCount: 0,
          maxRetries: 2,
        };
        
        queueRef.current.enqueue(task);
      }
    });
  }, [finalConfig, window.location.pathname || "/"]);

  /**
   * Adjacent content prefetching: prefetch related/nearby content
   */
  const prefetchAdjacentContent = useCallback((contentId: string, relatedIds: string[]) => {
    if (!finalConfig.enabledStrategies.includes(PrefetchStrategy.ADJACENT_CONTENT)) return;
    
    relatedIds.forEach((relatedId, index) => {
      const task: PrefetchTask = {
        id: `adjacent-${relatedId}-${Date.now()}`,
        target: relatedId,
        strategy: PrefetchStrategy.ADJACENT_CONTENT,
        priority: index < 2 ? PrefetchPriority.HIGH : PrefetchPriority.MEDIUM,
        probability: Math.max(0.5 - (index * 0.1), 0.2),
        estimatedSizeKB: 8,
        createdAt: Date.now(),
        retryCount: 0,
        maxRetries: 1,
      };
      
      queueRef.current.enqueue(task);
    });
  }, [finalConfig]);

  // =============================================================================
  // Effect Hooks
  // =============================================================================

  /**
   * Initialize background processing
   */
  useEffect(() => {
    if (!finalConfig.enabled) return;
    
    // Start background queue processing
    processingRef.current = setInterval(processQueue, finalConfig.processingIntervalMs);
    
    return () => {
      if (processingRef.current) {
        clearInterval(processingRef.current);
        processingRef.current = null;
      }
    };
  }, [finalConfig.enabled, finalConfig.processingIntervalMs, processQueue]);

  /**
   * Track route changes for behavior analysis
   */
  useEffect(() => {
    if (!finalConfig.analyticsEnabled) return;
    
    // Record route access
    behaviorAnalyzerRef.current.recordAccess(window.location.pathname || "/");
    
    // Trigger route-based prefetching
    prefetchRouteBasedContent(window.location.pathname || "/");
    
    // Trigger behavior-based prefetching with debouncing
    const timer = setTimeout(() => {
      prefetchBehaviorBasedContent();
    }, 1000);
    
    return () => clearTimeout(timer);
  }, [window.location.pathname || "/", finalConfig.analyticsEnabled, prefetchRouteBasedContent, prefetchBehaviorBasedContent]);

  // =============================================================================
  // Public API
  // =============================================================================

  /**
   * Manually trigger prefetch for specific content
   */
  const prefetchContent = useCallback((contentId: string, priority: PrefetchPriority = PrefetchPriority.MEDIUM) => {
    const task: PrefetchTask = {
      id: `manual-${contentId}-${Date.now()}`,
      target: contentId,
      strategy: PrefetchStrategy.BEHAVIOR_BASED,
      priority,
      probability: 1.0, // Manual requests have high probability
      estimatedSizeKB: 5,
      createdAt: Date.now(),
      retryCount: 0,
      maxRetries: 3,
    };
    
    queueRef.current.enqueue(task);
  }, []);

  /**
   * Prefetch content for offline support
   */
  const prefetchForOffline = useCallback((categories: string[]) => {
    categories.forEach(category => {
      const task: PrefetchTask = {
        id: `offline-${category}-${Date.now()}`,
        target: category,
        strategy: PrefetchStrategy.OFFLINE_SUPPORT,
        priority: PrefetchPriority.LOW,
        probability: 0.9,
        estimatedSizeKB: 100,
        createdAt: Date.now(),
        retryCount: 0,
        maxRetries: 1,
      };
      
      queueRef.current.enqueue(task);
    });
  }, []);

  /**
   * Get current analytics data
   */
  const getAnalytics = useCallback((): PrefetchAnalytics => {
    return { ...analyticsRef.current };
  }, []);

  /**
   * Get current queue status
   */
  const getQueueStatus = useCallback(() => {
    const queue = queueRef.current;
    return {
      size: queue.size(),
      activeRequests: activeRequestsRef.current.size,
      sessionDataUsage: sessionDataUsageRef.current,
      operationsThisMinute: operationCountRef.current.count,
    };
  }, []);

  /**
   * Clear prefetch queue and cache
   */
  const clearPrefetchCache = useCallback(() => {
    queueRef.current.clear();
    activeRequestsRef.current.clear();
    sessionDataUsageRef.current = 0;
    operationCountRef.current = { count: 0, windowStart: Date.now() };
    
    // Reset analytics
    if (finalConfig.analyticsEnabled) {
      analyticsRef.current = {
        totalOperations: 0,
        successfulHits: 0,
        hitRate: 0,
        avgResponseImprovement: 0,
        dataUsage: { total: 0, successful: 0, wasted: 0 },
        strategyStats: Object.values(PrefetchStrategy).reduce((acc, strategy) => {
          acc[strategy] = { operations: 0, hits: 0, hitRate: 0 };
          return acc;
        }, {} as any),
      };
    }
  }, [finalConfig.analyticsEnabled]);

  /**
   * Update configuration
   */
  const updateConfig = useCallback((newConfig: Partial<PrefetchConfig>) => {
    Object.assign(finalConfig, newConfig);
  }, [finalConfig]);

  return {
    // Core prefetching functions
    prefetchContent,
    prefetchCategoryContent,
    prefetchAdjacentContent,
    prefetchForOffline,
    
    // Analytics and monitoring
    getAnalytics,
    getQueueStatus,
    clearPrefetchCache,
    
    // Configuration
    updateConfig,
    config: finalConfig,
    
    // Network status
    isNetworkSuitable: isNetworkSuitableForPrefetching(getNetworkCondition()),
    networkCondition: getNetworkCondition(),
  };
};

// =============================================================================
// Utility Hooks and Components
// =============================================================================

/**
 * Hook for route-specific prefetching
 */
export const useRoutePrefetch = (routes: string[], priority: PrefetchPriority = PrefetchPriority.MEDIUM) => {
  const { prefetchContent } = useContentPrefetch();
  
  useEffect(() => {
    routes.forEach(route => {
      prefetchContent(route, priority);
    });
  }, [routes, priority, prefetchContent]);
};

/**
 * Hook for category prefetching with navigation awareness
 */
export const useCategoryPrefetch = (categories: string[]) => {
  const { prefetchCategoryContent } = useContentPrefetch();
  
  
  useEffect(() => {
    // Determine priority based on current route
    const currentCategory = window.location.pathname || "/".split('/')[1];
    
    categories.forEach(category => {
      const priority = category === currentCategory 
        ? PrefetchPriority.HIGH 
        : PrefetchPriority.MEDIUM;
      
      prefetchCategoryContent(category, priority);
    });
  }, [categories, window.location.pathname || "/", prefetchCategoryContent]);
};

/**
 * Hook for smart prefetching with performance monitoring
 */
export const useSmartPrefetch = () => {
  const prefetch = useContentPrefetch();
  const [performanceMetrics, setPerformanceMetrics] = useState({
    loadTime: 0,
    cacheHitRate: 0,
    networkEfficiency: 0,
  });
  
  useEffect(() => {
    const updateMetrics = () => {
      const analytics = prefetch.getAnalytics();
      const queueStatus = prefetch.getQueueStatus();
      
      setPerformanceMetrics({
        loadTime: analytics.avgResponseImprovement,
        cacheHitRate: analytics.hitRate * 100,
        networkEfficiency: analytics.dataUsage.successful / (analytics.dataUsage.total || 1) * 100,
      });
    };
    
    const interval = setInterval(updateMetrics, 30000); // Update every 30 seconds
    updateMetrics(); // Initial update
    
    return () => clearInterval(interval);
  }, [prefetch]);
  
  return {
    ...prefetch,
    performanceMetrics,
  };
};

export default useContentPrefetch;
