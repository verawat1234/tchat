/**
 * Content Prefetching Usage Examples
 *
 * This file demonstrates how to use the useContentPrefetch hook and its utilities
 * for implementing intelligent content prefetching strategies in your application.
 */

import React, { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  useContentPrefetch,
  useRoutePrefetch,
  useCategoryPrefetch,
  useSmartPrefetch,
  PrefetchPriority,
  PrefetchStrategy,
  type PrefetchConfig,
} from './useContentPrefetch';

// =============================================================================
// Basic Usage Examples
// =============================================================================

/**
 * Example 1: Basic Content Prefetching
 * Simple integration for manual prefetching
 */
export const BasicPrefetchExample: React.FC = () => {
  const {
    prefetchContent,
    prefetchCategoryContent,
    getAnalytics,
    getQueueStatus,
    isNetworkSuitable,
  } = useContentPrefetch();

  // Prefetch critical content on component mount
  useEffect(() => {
    if (isNetworkSuitable) {
      // Prefetch high-priority navigation content
      prefetchContent('navigation.header.app_title', PrefetchPriority.HIGH);
      prefetchContent('navigation.tabs.chat', PrefetchPriority.HIGH);
      
      // Prefetch entire navigation category
      prefetchCategoryContent('navigation', PrefetchPriority.MEDIUM);
    }
  }, [isNetworkSuitable, prefetchContent, prefetchCategoryContent]);

  // Monitor prefetch performance
  const analytics = getAnalytics();
  const queueStatus = getQueueStatus();

  return (
    <div className="prefetch-monitor">
      <h3>Prefetch Status</h3>
      <p>Network Suitable: {isNetworkSuitable ? 'Yes' : 'No'}</p>
      <p>Queue Size: {queueStatus.size}</p>
      <p>Active Requests: {queueStatus.activeRequests}</p>
      <p>Hit Rate: {(analytics.hitRate * 100).toFixed(1)}%</p>
      <p>Data Usage: {queueStatus.sessionDataUsage}KB</p>
    </div>
  );
};

/**
 * Example 2: Advanced Configuration
 * Custom configuration for specific performance requirements
 */
export const AdvancedPrefetchExample: React.FC = () => {
  // Custom configuration for mobile-optimized prefetching
  const mobileConfig: Partial<PrefetchConfig> = {
    budget: {
      maxMemoryMB: 20, // Lower memory usage for mobile
      maxOperationsPerMinute: 15, // Reduced operations
      maxConcurrentRequests: 2, // Fewer concurrent requests
      maxDataPerSessionMB: 5, // Strict data limit
    },
    enabledStrategies: [
      PrefetchStrategy.ROUTE_BASED,
      PrefetchStrategy.BEHAVIOR_BASED,
      // Exclude category bulk prefetching on mobile
    ],
    minProbabilityThreshold: 0.5, // Higher threshold for mobile
    processingIntervalMs: 3000, // Less frequent processing
  };

  const {
    prefetchContent,
    prefetchAdjacentContent,
    updateConfig,
    networkCondition,
    config,
  } = useContentPrefetch(mobileConfig);

  // Dynamically adjust configuration based on network conditions
  useEffect(() => {
    if (networkCondition.effectiveType === '2g' || networkCondition.saveData) {
      updateConfig({
        enabled: false, // Disable on very slow connections
      });
    } else if (networkCondition.effectiveType === 'wifi') {
      updateConfig({
        budget: {
          ...config.budget,
          maxConcurrentRequests: 5, // More aggressive on WiFi
          maxDataPerSessionMB: 15,
        },
      });
    }
  }, [networkCondition, updateConfig, config.budget]);

  const handleContentView = (contentId: string, relatedIds: string[]) => {
    // Prefetch related content when user views content
    prefetchAdjacentContent(contentId, relatedIds);
  };

  return (
    <div className="advanced-prefetch">
      <h3>Network-Adaptive Prefetching</h3>
      <p>Connection: {networkCondition.effectiveType}</p>
      <p>RTT: {networkCondition.rtt}ms</p>
      <p>Downlink: {networkCondition.downlink}Mbps</p>
      <p>Data Saver: {networkCondition.saveData ? 'On' : 'Off'}</p>
      <p>Prefetching: {config.enabled ? 'Enabled' : 'Disabled'}</p>
      
      <button 
        onClick={() => handleContentView('article-123', ['article-124', 'article-125'])}
      >
        View Article (Triggers Adjacent Prefetch)
      </button>
    </div>
  );
};

/**
 * Example 3: Route-Based Prefetching
 * Prefetch content based on likely navigation patterns
 */
export const RoutePrefetchExample: React.FC = () => {
  const location = useLocation();
  
  // Define route-based prefetching strategies
  const routePrefetchMap: Record<string, string[]> = {
    '/': ['navigation.tabs.chat', 'navigation.tabs.store', 'navigation.tabs.social'],
    '/chat': ['chat.recent_messages', 'chat.contacts', 'chat.new_message'],
    '/store': ['store.categories', 'store.featured', 'store.cart'],
    '/social': ['social.feed', 'social.notifications', 'social.friends'],
    '/profile': ['profile.settings', 'profile.preferences', 'profile.security'],
  };

  // Prefetch routes that are commonly accessed after current route
  const nextLikelyRoutes = routePrefetchMap[location.pathname] || [];
  useRoutePrefetch(nextLikelyRoutes, PrefetchPriority.MEDIUM);

  return (
    <div className="route-prefetch">
      <h3>Current Route: {location.pathname}</h3>
      <p>Prefetching: {nextLikelyRoutes.join(', ')}</p>
    </div>
  );
};

/**
 * Example 4: Category-Based Prefetching
 * Bulk prefetch content categories based on user context
 */
export const CategoryPrefetchExample: React.FC = () => {
  const location = useLocation();
  
  // Define categories to prefetch based on current section
  const categoryMap: Record<string, string[]> = {
    '/chat': ['navigation', 'messages', 'contacts'],
    '/store': ['navigation', 'products', 'commerce'],
    '/social': ['navigation', 'social', 'notifications'],
    '/settings': ['navigation', 'settings', 'user'],
  };

  const currentSection = '/' + location.pathname.split('/')[1];
  const categoriesToPrefetch = categoryMap[currentSection] || ['navigation'];

  useCategoryPrefetch(categoriesToPrefetch);

  return (
    <div className="category-prefetch">
      <h3>Section: {currentSection}</h3>
      <p>Prefetching Categories: {categoriesToPrefetch.join(', ')}</p>
    </div>
  );
};

/**
 * Example 5: Smart Prefetching with Performance Monitoring
 * Advanced prefetching with real-time performance metrics
 */
export const SmartPrefetchExample: React.FC = () => {
  const {
    prefetchContent,
    prefetchForOffline,
    clearPrefetchCache,
    performanceMetrics,
    getAnalytics,
    getQueueStatus,
    networkCondition,
  } = useSmartPrefetch();

  const analytics = getAnalytics();
  const queueStatus = getQueueStatus();

  // Handle offline preparation
  const handlePrepareOffline = () => {
    const essentialCategories = ['navigation', 'core', 'error_messages'];
    prefetchForOffline(essentialCategories);
  };

  // Handle cache management
  const handleOptimizeCache = () => {
    if (performanceMetrics.networkEfficiency < 50) {
      // Clear cache if efficiency is poor
      clearPrefetchCache();
    }
  };

  return (
    <div className="smart-prefetch-dashboard">
      <h3>Smart Prefetch Dashboard</h3>
      
      <div className="metrics-grid">
        <div className="metric-card">
          <h4>Performance</h4>
          <p>Load Time Improvement: {performanceMetrics.loadTime.toFixed(1)}ms</p>
          <p>Cache Hit Rate: {performanceMetrics.cacheHitRate.toFixed(1)}%</p>
          <p>Network Efficiency: {performanceMetrics.networkEfficiency.toFixed(1)}%</p>
        </div>
        
        <div className="metric-card">
          <h4>Queue Status</h4>
          <p>Pending Tasks: {queueStatus.size}</p>
          <p>Active Requests: {queueStatus.activeRequests}</p>
          <p>Operations/Min: {queueStatus.operationsThisMinute}</p>
        </div>
        
        <div className="metric-card">
          <h4>Analytics</h4>
          <p>Total Operations: {analytics.totalOperations}</p>
          <p>Successful Hits: {analytics.successfulHits}</p>
          <p>Data Successful: {analytics.dataUsage.successful}KB</p>
          <p>Data Wasted: {analytics.dataUsage.wasted}KB</p>
        </div>
        
        <div className="metric-card">
          <h4>Network</h4>
          <p>Type: {networkCondition.effectiveType}</p>
          <p>RTT: {networkCondition.rtt}ms</p>
          <p>Downlink: {networkCondition.downlink}Mbps</p>
          <p>Data Saver: {networkCondition.saveData ? 'On' : 'Off'}</p>
        </div>
      </div>

      <div className="strategy-stats">
        <h4>Strategy Effectiveness</h4>
        {Object.entries(analytics.strategyStats).map(([strategy, stats]) => (
          <div key={strategy} className="strategy-stat">
            <span>{strategy}:</span>
            <span>{stats.operations} ops</span>
            <span>{(stats.hitRate * 100).toFixed(1)}% hit rate</span>
          </div>
        ))}
      </div>

      <div className="control-buttons">
        <button onClick={handlePrepareOffline}>
          Prepare for Offline
        </button>
        <button onClick={handleOptimizeCache}>
          Optimize Cache
        </button>
        <button onClick={() => clearPrefetchCache()}>
          Clear Cache
        </button>
      </div>
    </div>
  );
};

/**
 * Example 6: E-commerce Product Prefetching
 * Specialized prefetching for product browsing scenarios
 */
export const EcommercePrefetchExample: React.FC = () => {
  const { prefetchContent, prefetchAdjacentContent } = useContentPrefetch();

  // Prefetch product details when user hovers over product cards
  const handleProductHover = (productId: string) => {
    prefetchContent(`product.${productId}.details`, PrefetchPriority.HIGH);
    prefetchContent(`product.${productId}.images`, PrefetchPriority.MEDIUM);
    prefetchContent(`product.${productId}.reviews`, PrefetchPriority.LOW);
  };

  // Prefetch related products when viewing a product
  const handleProductView = (productId: string, categoryId: string) => {
    // Prefetch similar products in the same category
    const relatedProductIds = [
      `product.${categoryId}.similar.1`,
      `product.${categoryId}.similar.2`,
      `product.${categoryId}.similar.3`,
    ];
    
    prefetchAdjacentContent(productId, relatedProductIds);
    
    // Prefetch shopping cart and checkout flows
    prefetchContent('cart.summary', PrefetchPriority.MEDIUM);
    prefetchContent('checkout.shipping', PrefetchPriority.LOW);
  };

  // Prefetch user's likely next actions in shopping funnel
  const handleAddToCart = (productId: string) => {
    // High probability user will view cart or proceed to checkout
    prefetchContent('cart.items', PrefetchPriority.HIGH);
    prefetchContent('cart.total', PrefetchPriority.HIGH);
    prefetchContent('checkout.form', PrefetchPriority.MEDIUM);
    prefetchContent('checkout.payment_methods', PrefetchPriority.MEDIUM);
  };

  return (
    <div className="ecommerce-prefetch">
      <h3>E-commerce Prefetching</h3>
      <div className="product-grid">
        {/* Product cards with hover prefetching */}
        <div 
          className="product-card"
          onMouseEnter={() => handleProductHover('laptop-123')}
          onClick={() => handleProductView('laptop-123', 'electronics')}
        >
          <h4>Gaming Laptop</h4>
          <p>Hover to prefetch details</p>
        </div>
        
        <div 
          className="product-card"
          onMouseEnter={() => handleProductHover('phone-456')}
          onClick={() => handleProductView('phone-456', 'electronics')}
        >
          <h4>Smartphone</h4>
          <p>Hover to prefetch details</p>
        </div>
      </div>
      
      <button onClick={() => handleAddToCart('laptop-123')}>
        Add to Cart (Triggers Checkout Prefetch)
      </button>
    </div>
  );
};

/**
 * Example 7: Content Management Dashboard
 * Prefetching for content editing and management workflows
 */
export const ContentManagementPrefetchExample: React.FC = () => {
  const { 
    prefetchContent, 
    prefetchCategoryContent, 
    getAnalytics,
    clearPrefetchCache,
    updateConfig,
  } = useContentPrefetch();

  // Prefetch content editing dependencies
  const handleEditContent = (contentId: string) => {
    // Prefetch edit form, validation rules, and preview content
    prefetchContent(`content.${contentId}.edit_form`, PrefetchPriority.HIGH);
    prefetchContent(`content.${contentId}.validation_rules`, PrefetchPriority.HIGH);
    prefetchContent(`content.${contentId}.preview`, PrefetchPriority.MEDIUM);
    prefetchContent('content.templates', PrefetchPriority.LOW);
  };

  // Bulk prefetch for content category management
  const handleCategoryManagement = (category: string) => {
    prefetchCategoryContent(category, PrefetchPriority.HIGH);
    prefetchContent(`category.${category}.permissions`, PrefetchPriority.MEDIUM);
    prefetchContent(`category.${category}.analytics`, PrefetchPriority.LOW);
  };

  // Performance-aware mode for content managers
  const handleTogglePerformanceMode = () => {
    const analytics = getAnalytics();
    
    if (analytics.hitRate < 0.5) {
      // Low hit rate - reduce prefetching
      updateConfig({
        minProbabilityThreshold: 0.8,
        enabledStrategies: [PrefetchStrategy.ROUTE_BASED],
      });
    } else {
      // Good hit rate - enable more aggressive prefetching
      updateConfig({
        minProbabilityThreshold: 0.3,
        enabledStrategies: [
          PrefetchStrategy.ROUTE_BASED,
          PrefetchStrategy.CATEGORY_BULK,
          PrefetchStrategy.BEHAVIOR_BASED,
        ],
      });
    }
  };

  return (
    <div className="content-management-prefetch">
      <h3>Content Management Prefetching</h3>
      
      <div className="content-actions">
        <button onClick={() => handleEditContent('article-123')}>
          Edit Article (Prefetch Edit Tools)
        </button>
        
        <button onClick={() => handleCategoryManagement('blog')}>
          Manage Blog Category
        </button>
        
        <button onClick={handleTogglePerformanceMode}>
          Toggle Performance Mode
        </button>
        
        <button onClick={() => clearPrefetchCache()}>
          Clear Cache
        </button>
      </div>
      
      <div className="prefetch-insights">
        <h4>Prefetch Insights</h4>
        <p>Use these patterns to optimize content management workflows:</p>
        <ul>
          <li>Prefetch edit forms when hovering over content items</li>
          <li>Bulk prefetch categories when entering management sections</li>
          <li>Prefetch validation rules and templates for content creation</li>
          <li>Adjust prefetch aggressiveness based on hit rate performance</li>
        </ul>
      </div>
    </div>
  );
};

// =============================================================================
// Integration Example Component
// =============================================================================

/**
 * Complete Integration Example
 * Demonstrates how to integrate all prefetching strategies in a real application
 */
export const CompletePrefetchIntegration: React.FC = () => {
  return (
    <div className="prefetch-integration-demo">
      <h2>Content Prefetching Integration Demo</h2>
      
      <div className="examples-grid">
        <div className="example-section">
          <h3>Basic Usage</h3>
          <BasicPrefetchExample />
        </div>
        
        <div className="example-section">
          <h3>Advanced Configuration</h3>
          <AdvancedPrefetchExample />
        </div>
        
        <div className="example-section">
          <h3>Route-Based Prefetching</h3>
          <RoutePrefetchExample />
        </div>
        
        <div className="example-section">
          <h3>Category Prefetching</h3>
          <CategoryPrefetchExample />
        </div>
        
        <div className="example-section">
          <h3>Smart Prefetching Dashboard</h3>
          <SmartPrefetchExample />
        </div>
        
        <div className="example-section">
          <h3>E-commerce Integration</h3>
          <EcommercePrefetchExample />
        </div>
        
        <div className="example-section">
          <h3>Content Management</h3>
          <ContentManagementPrefetchExample />
        </div>
      </div>
      
      <div className="integration-notes">
        <h3>Integration Best Practices</h3>
        <ol>
          <li><strong>Start Simple:</strong> Begin with basic manual prefetching for critical content</li>
          <li><strong>Monitor Performance:</strong> Use analytics to optimize prefetch strategies</li>
          <li><strong>Respect Network Conditions:</strong> Adapt behavior based on connection quality</li>
          <li><strong>Set Performance Budgets:</strong> Configure limits based on your application's requirements</li>
          <li><strong>Use Priority Wisely:</strong> Reserve high priority for content likely to be accessed immediately</li>
          <li><strong>Test Different Strategies:</strong> Experiment with different prefetch strategies to find what works best</li>
          <li><strong>Handle Failures Gracefully:</strong> Ensure your application works well even if prefetching fails</li>
        </ol>
      </div>
    </div>
  );
};

export default CompletePrefetchIntegration;
