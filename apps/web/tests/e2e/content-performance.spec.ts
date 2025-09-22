/**
 * E2E Content Performance Tests
 *
 * Comprehensive performance validation for the content management system.
 * Tests validate load times, Core Web Vitals, bundle impact, memory usage,
 * network efficiency, concurrent operations, and cross-device performance.
 *
 * Performance Requirements:
 * - Content load time: <200ms budget
 * - LCP: <2.5s, FID: <100ms, CLS: <0.1
 * - Bundle size: <500KB initial, <2MB total
 * - Memory usage: <100MB mobile, <500MB desktop
 * - Network efficiency: optimal caching and request minimization
 * - Concurrent operations: handle 10+ simultaneous requests
 * - Cross-device: consistent performance across devices/networks
 */

import { test, expect, Page, Browser } from '@playwright/test';

// Performance budgets and thresholds
const PERFORMANCE_BUDGETS = {
  CONTENT_LOAD_TIME: 200, // ms
  LCP_THRESHOLD: 2500, // ms
  FID_THRESHOLD: 100, // ms
  CLS_THRESHOLD: 0.1,
  BUNDLE_SIZE_INITIAL: 500 * 1024, // 500KB
  BUNDLE_SIZE_TOTAL: 2 * 1024 * 1024, // 2MB
  MEMORY_MOBILE: 100 * 1024 * 1024, // 100MB
  MEMORY_DESKTOP: 500 * 1024 * 1024, // 500MB
  CONCURRENT_REQUESTS: 10,
  API_RESPONSE_TIME: 200, // ms
  CACHE_HIT_RATIO: 0.8, // 80%
};

// Test data for content operations
const TEST_CONTENT = {
  category: 'test-category',
  items: [
    { id: 'content-1', type: 'text', value: 'Test content 1' },
    { id: 'content-2', type: 'rich_text', value: '<h1>Rich Test Content</h1>' },
    { id: 'content-3', type: 'config', value: { setting: 'value' } },
  ],
};

// Device configurations for testing
const DEVICE_CONFIGS = [
  { name: 'Desktop High-End', viewport: { width: 1920, height: 1080 }, cpu: 4, memory: 8 },
  { name: 'Desktop Standard', viewport: { width: 1366, height: 768 }, cpu: 2, memory: 4 },
  { name: 'Tablet', viewport: { width: 768, height: 1024 }, cpu: 2, memory: 2 },
  { name: 'Mobile High-End', viewport: { width: 390, height: 844 }, cpu: 2, memory: 2 },
  { name: 'Mobile Low-End', viewport: { width: 375, height: 667 }, cpu: 1, memory: 1 },
];

// Network conditions
const NETWORK_CONDITIONS = [
  { name: 'Fast 3G', download: 1600, upload: 750, latency: 150 },
  { name: 'Slow 3G', download: 500, upload: 500, latency: 300 },
  { name: 'Offline', download: 0, upload: 0, latency: 0 },
];

/**
 * Helper function to measure Core Web Vitals
 */
async function measureCoreWebVitals(page: Page) {
  return await page.evaluate(() => {
    return new Promise((resolve) => {
      const metrics = {
        LCP: 0,
        FID: 0,
        CLS: 0,
        TTFB: 0,
        FCP: 0,
      };

      // Largest Contentful Paint
      new PerformanceObserver((list) => {
        const entries = list.getEntries();
        const lastEntry = entries[entries.length - 1];
        metrics.LCP = lastEntry?.startTime || 0;
      }).observe({ entryTypes: ['largest-contentful-paint'] });

      // First Input Delay
      new PerformanceObserver((list) => {
        const entries = list.getEntries();
        entries.forEach((entry: any) => {
          metrics.FID = entry.processingStart - entry.startTime;
        });
      }).observe({ entryTypes: ['first-input'] });

      // Cumulative Layout Shift
      let clsValue = 0;
      new PerformanceObserver((list) => {
        for (const entry of list.getEntries() as any[]) {
          if (!entry.hadRecentInput) {
            clsValue += entry.value;
          }
        }
        metrics.CLS = clsValue;
      }).observe({ entryTypes: ['layout-shift'] });

      // First Contentful Paint
      const fcpEntry = performance.getEntriesByName('first-contentful-paint')[0];
      metrics.FCP = fcpEntry?.startTime || 0;

      // Time to First Byte
      const navEntry = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
      metrics.TTFB = navEntry?.responseStart - navEntry?.requestStart || 0;

      // Return metrics after a delay to allow measurements to complete
      setTimeout(() => resolve(metrics), 2000);
    });
  });
}

/**
 * Helper function to measure memory usage
 */
async function measureMemoryUsage(page: Page) {
  return await page.evaluate(() => {
    // @ts-ignore - performance.memory is available in Chrome
    const memory = (performance as any).memory;
    if (memory) {
      return {
        usedJSHeapSize: memory.usedJSHeapSize,
        totalJSHeapSize: memory.totalJSHeapSize,
        jsHeapSizeLimit: memory.jsHeapSizeLimit,
      };
    }
    return null;
  });
}

/**
 * Helper function to measure network activity
 */
async function measureNetworkActivity(page: Page) {
  const requests: Array<{ url: string; size: number; time: number; cached: boolean }> = [];

  page.on('response', async (response) => {
    const request = response.request();
    const timing = await response.allHeaders();
    const cached = response.fromCache();

    try {
      const body = await response.body();
      requests.push({
        url: request.url(),
        size: body.length,
        time: Date.now(),
        cached,
      });
    } catch {
      // Handle cases where response body is not available
    }
  });

  return requests;
}

/**
 * Helper function to simulate content operations
 */
async function simulateContentOperations(page: Page, operations: number = 5) {
  const results: Array<{ operation: string; time: number; success: boolean }> = [];

  for (let i = 0; i < operations; i++) {
    const startTime = performance.now();

    try {
      // Simulate different content operations
      const operation = i % 4;
      switch (operation) {
        case 0: // Load content list
          await page.click('[data-testid="content-list-button"]');
          await page.waitForSelector('[data-testid="content-item"]', { timeout: PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME });
          break;
        case 1: // Create content
          await page.click('[data-testid="create-content-button"]');
          await page.fill('[data-testid="content-id-input"]', `test-content-${i}`);
          await page.fill('[data-testid="content-value-input"]', `Test value ${i}`);
          await page.click('[data-testid="save-content-button"]');
          break;
        case 2: // Edit content
          await page.click('[data-testid="edit-content-button"]');
          await page.fill('[data-testid="content-value-input"]', `Updated value ${i}`);
          await page.click('[data-testid="save-content-button"]');
          break;
        case 3: // Delete content
          await page.click('[data-testid="delete-content-button"]');
          await page.click('[data-testid="confirm-delete-button"]');
          break;
      }

      const endTime = performance.now();
      results.push({
        operation: `operation-${operation}`,
        time: endTime - startTime,
        success: true,
      });
    } catch (error) {
      const endTime = performance.now();
      results.push({
        operation: `operation-${operation}`,
        time: endTime - startTime,
        success: false,
      });
    }
  }

  return results;
}

test.describe('Content Performance Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Set up performance monitoring
    await page.goto('/content');
    await page.waitForLoadState('networkidle');
  });

  test.describe('1. Load Time Validation (<200ms budget)', () => {
    test('should load content list within performance budget', async ({ page }) => {
      const startTime = performance.now();

      await page.click('[data-testid="content-list-button"]');
      await page.waitForSelector('[data-testid="content-item"]');

      const endTime = performance.now();
      const loadTime = endTime - startTime;

      expect(loadTime).toBeLessThan(PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME);
      console.log(`Content list load time: ${loadTime.toFixed(2)}ms`);
    });

    test('should load individual content items within budget', async ({ page }) => {
      // First load the content list
      await page.click('[data-testid="content-list-button"]');
      await page.waitForSelector('[data-testid="content-item"]');

      const contentItems = await page.locator('[data-testid="content-item"]').all();

      for (let i = 0; i < Math.min(5, contentItems.length); i++) {
        const startTime = performance.now();

        await contentItems[i].click();
        await page.waitForSelector('[data-testid="content-detail"]');

        const endTime = performance.now();
        const loadTime = endTime - startTime;

        expect(loadTime).toBeLessThan(PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME);
        console.log(`Content item ${i} load time: ${loadTime.toFixed(2)}ms`);

        // Navigate back
        await page.click('[data-testid="back-button"]');
      }
    });

    test('should load content creation form within budget', async ({ page }) => {
      const startTime = performance.now();

      await page.click('[data-testid="create-content-button"]');
      await page.waitForSelector('[data-testid="content-form"]');

      const endTime = performance.now();
      const loadTime = endTime - startTime;

      expect(loadTime).toBeLessThan(PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME);
      console.log(`Content form load time: ${loadTime.toFixed(2)}ms`);
    });
  });

  test.describe('2. Core Web Vitals Testing', () => {
    test('should meet Core Web Vitals thresholds', async ({ page }) => {
      // Navigate to content-heavy page
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      const metrics = await measureCoreWebVitals(page);

      console.log('Core Web Vitals:', metrics);

      expect(metrics.LCP).toBeLessThan(PERFORMANCE_BUDGETS.LCP_THRESHOLD);
      expect(metrics.FID).toBeLessThan(PERFORMANCE_BUDGETS.FID_THRESHOLD);
      expect(metrics.CLS).toBeLessThan(PERFORMANCE_BUDGETS.CLS_THRESHOLD);

      // Additional metrics for monitoring
      expect(metrics.FCP).toBeLessThan(1800); // First Contentful Paint
      expect(metrics.TTFB).toBeLessThan(600); // Time to First Byte
    });

    test('should maintain Core Web Vitals during content interactions', async ({ page }) => {
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      // Perform various content interactions
      await page.click('[data-testid="content-list-button"]');
      await page.waitForSelector('[data-testid="content-item"]');

      await page.click('[data-testid="filter-button"]');
      await page.selectOption('[data-testid="category-filter"]', 'test-category');

      await page.fill('[data-testid="search-input"]', 'test');
      await page.press('[data-testid="search-input"]', 'Enter');

      const metrics = await measureCoreWebVitals(page);

      expect(metrics.CLS).toBeLessThan(PERFORMANCE_BUDGETS.CLS_THRESHOLD);
      console.log('CLS during interactions:', metrics.CLS);
    });
  });

  test.describe('3. Bundle Impact Validation', () => {
    test('should not exceed bundle size budgets', async ({ page }) => {
      const responses: Array<{ url: string; size: number }> = [];

      page.on('response', async (response) => {
        if (response.url().includes('.js') || response.url().includes('.css')) {
          try {
            const body = await response.body();
            responses.push({
              url: response.url(),
              size: body.length,
            });
          } catch {
            // Handle cases where response body is not available
          }
        }
      });

      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      const initialBundle = responses.filter(r =>
        r.url.includes('index') || r.url.includes('main') || r.url.includes('app')
      );
      const totalBundle = responses;

      const initialSize = initialBundle.reduce((sum, r) => sum + r.size, 0);
      const totalSize = totalBundle.reduce((sum, r) => sum + r.size, 0);

      console.log(`Initial bundle size: ${(initialSize / 1024).toFixed(2)}KB`);
      console.log(`Total bundle size: ${(totalSize / 1024 / 1024).toFixed(2)}MB`);

      expect(initialSize).toBeLessThan(PERFORMANCE_BUDGETS.BUNDLE_SIZE_INITIAL);
      expect(totalSize).toBeLessThan(PERFORMANCE_BUDGETS.BUNDLE_SIZE_TOTAL);
    });

    test('should load content features without increasing bundle size significantly', async ({ page }) => {
      // Measure baseline bundle size
      let baselineSize = 0;
      let contentFeatureSize = 0;

      page.on('response', async (response) => {
        if (response.url().includes('.js')) {
          try {
            const body = await response.body();
            baselineSize += body.length;
          } catch {
            // Handle cases where response body is not available
          }
        }
      });

      await page.goto('/');
      await page.waitForLoadState('networkidle');

      // Reset counters and navigate to content features
      page.removeAllListeners('response');

      page.on('response', async (response) => {
        if (response.url().includes('.js')) {
          try {
            const body = await response.body();
            contentFeatureSize += body.length;
          } catch {
            // Handle cases where response body is not available
          }
        }
      });

      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      const additionalSize = contentFeatureSize - baselineSize;

      console.log(`Additional content bundle size: ${(additionalSize / 1024).toFixed(2)}KB`);

      // Content features should not add more than 200KB
      expect(additionalSize).toBeLessThan(200 * 1024);
    });
  });

  test.describe('4. Memory Usage Testing', () => {
    test('should maintain memory usage within limits on desktop', async ({ page }) => {
      await page.setViewportSize({ width: 1920, height: 1080 });
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      // Perform memory-intensive operations
      for (let i = 0; i < 10; i++) {
        await page.click('[data-testid="content-list-button"]');
        await page.waitForSelector('[data-testid="content-item"]');

        if (i % 3 === 0) {
          // Force garbage collection if available
          await page.evaluate(() => {
            // @ts-ignore
            if (window.gc) window.gc();
          });
        }
      }

      const memory = await measureMemoryUsage(page);

      if (memory) {
        console.log(`Memory usage: ${(memory.usedJSHeapSize / 1024 / 1024).toFixed(2)}MB`);
        expect(memory.usedJSHeapSize).toBeLessThan(PERFORMANCE_BUDGETS.MEMORY_DESKTOP);
      }
    });

    test('should maintain memory usage within limits on mobile', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      // Simulate mobile usage patterns
      for (let i = 0; i < 5; i++) {
        await page.click('[data-testid="content-list-button"]');
        await page.waitForSelector('[data-testid="content-item"]');

        // Simulate scrolling
        await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
        await page.evaluate(() => window.scrollTo(0, 0));
      }

      const memory = await measureMemoryUsage(page);

      if (memory) {
        console.log(`Mobile memory usage: ${(memory.usedJSHeapSize / 1024 / 1024).toFixed(2)}MB`);
        expect(memory.usedJSHeapSize).toBeLessThan(PERFORMANCE_BUDGETS.MEMORY_MOBILE);
      }
    });

    test('should handle memory cleanup after content operations', async ({ page }) => {
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      const initialMemory = await measureMemoryUsage(page);

      // Perform many content operations
      for (let i = 0; i < 20; i++) {
        await page.click('[data-testid="create-content-button"]');
        await page.fill('[data-testid="content-id-input"]', `temp-content-${i}`);
        await page.fill('[data-testid="content-value-input"]', `Temporary content ${i}`);
        await page.press('Escape'); // Cancel creation
      }

      // Force cleanup
      await page.evaluate(() => {
        // @ts-ignore
        if (window.gc) window.gc();
      });

      const finalMemory = await measureMemoryUsage(page);

      if (initialMemory && finalMemory) {
        const memoryIncrease = finalMemory.usedJSHeapSize - initialMemory.usedJSHeapSize;
        console.log(`Memory increase after operations: ${(memoryIncrease / 1024 / 1024).toFixed(2)}MB`);

        // Memory increase should be minimal after cleanup
        expect(memoryIncrease).toBeLessThan(50 * 1024 * 1024); // 50MB
      }
    });
  });

  test.describe('5. Network Efficiency Validation', () => {
    test('should optimize network requests and caching', async ({ page }) => {
      const requests = await measureNetworkActivity(page);

      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      // Navigate around to trigger various requests
      await page.click('[data-testid="content-list-button"]');
      await page.waitForSelector('[data-testid="content-item"]');

      await page.click('[data-testid="categories-button"]');
      await page.waitForSelector('[data-testid="category-item"]');

      // Go back to content list (should use cache)
      await page.click('[data-testid="content-list-button"]');
      await page.waitForSelector('[data-testid="content-item"]');

      const apiRequests = requests.filter(r => r.url.includes('/api/'));
      const cachedRequests = apiRequests.filter(r => r.cached);
      const cacheHitRatio = cachedRequests.length / apiRequests.length;

      console.log(`API requests: ${apiRequests.length}, Cached: ${cachedRequests.length}`);
      console.log(`Cache hit ratio: ${(cacheHitRatio * 100).toFixed(1)}%`);

      expect(cacheHitRatio).toBeGreaterThan(PERFORMANCE_BUDGETS.CACHE_HIT_RATIO);
    });

    test('should minimize request payload sizes', async ({ page }) => {
      const requests: Array<{ url: string; size: number; type: string }> = [];

      page.on('request', (request) => {
        if (request.url().includes('/api/')) {
          const postData = request.postData();
          requests.push({
            url: request.url(),
            size: postData ? postData.length : 0,
            type: request.method(),
          });
        }
      });

      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      // Perform various content operations
      await simulateContentOperations(page, 5);

      const avgPayloadSize = requests.reduce((sum, r) => sum + r.size, 0) / requests.length;

      console.log(`Average request payload size: ${avgPayloadSize.toFixed(2)} bytes`);

      // Average payload should be reasonable
      expect(avgPayloadSize).toBeLessThan(10 * 1024); // 10KB
    });

    test('should handle network failures gracefully', async ({ page, context }) => {
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      // Simulate network failure
      await context.setOffline(true);

      const startTime = performance.now();

      try {
        await page.click('[data-testid="content-list-button"]');
        await page.waitForSelector('[data-testid="offline-message"]', { timeout: 5000 });

        const endTime = performance.now();
        const responseTime = endTime - startTime;

        // Should detect offline state quickly
        expect(responseTime).toBeLessThan(2000);

        // Should show appropriate offline message
        await expect(page.locator('[data-testid="offline-message"]')).toBeVisible();
      } finally {
        await context.setOffline(false);
      }
    });
  });

  test.describe('6. Concurrent Operations Performance', () => {
    test('should handle multiple simultaneous content requests', async ({ page }) => {
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      const operations = [];
      const startTime = performance.now();

      // Start multiple concurrent operations
      for (let i = 0; i < PERFORMANCE_BUDGETS.CONCURRENT_REQUESTS; i++) {
        operations.push(
          page.locator('[data-testid="content-list-button"]').click()
        );
      }

      await Promise.all(operations);

      const endTime = performance.now();
      const totalTime = endTime - startTime;

      console.log(`${PERFORMANCE_BUDGETS.CONCURRENT_REQUESTS} concurrent operations took: ${totalTime.toFixed(2)}ms`);

      // Should complete within reasonable time
      expect(totalTime).toBeLessThan(PERFORMANCE_BUDGETS.CONCURRENT_REQUESTS * 100);
    });

    test('should maintain performance under load', async ({ page }) => {
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      const results = await simulateContentOperations(page, 15);

      const avgTime = results.reduce((sum, r) => sum + r.time, 0) / results.length;
      const successRate = results.filter(r => r.success).length / results.length;

      console.log(`Average operation time under load: ${avgTime.toFixed(2)}ms`);
      console.log(`Success rate: ${(successRate * 100).toFixed(1)}%`);

      expect(avgTime).toBeLessThan(PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME * 2);
      expect(successRate).toBeGreaterThan(0.9); // 90% success rate
    });

    test('should handle concurrent user interactions', async ({ browser }) => {
      const contexts = await Promise.all([
        browser.newContext(),
        browser.newContext(),
        browser.newContext(),
      ]);

      const pages = await Promise.all(
        contexts.map(context => context.newPage())
      );

      try {
        // Start concurrent user sessions
        const sessions = pages.map(async (page, index) => {
          await page.goto('/content');
          await page.waitForLoadState('networkidle');

          const startTime = performance.now();

          // Simulate different user behaviors
          for (let i = 0; i < 5; i++) {
            switch (index) {
              case 0: // User 1: Browsing content
                await page.click('[data-testid="content-list-button"]');
                await page.waitForSelector('[data-testid="content-item"]');
                break;
              case 1: // User 2: Creating content
                await page.click('[data-testid="create-content-button"]');
                await page.fill('[data-testid="content-id-input"]', `user2-content-${i}`);
                await page.press('Escape');
                break;
              case 2: // User 3: Searching content
                await page.fill('[data-testid="search-input"]', `search-${i}`);
                await page.press('[data-testid="search-input"]', 'Enter');
                break;
            }
          }

          const endTime = performance.now();
          return endTime - startTime;
        });

        const sessionTimes = await Promise.all(sessions);
        const maxTime = Math.max(...sessionTimes);

        console.log(`Concurrent session times: ${sessionTimes.map(t => t.toFixed(2)).join(', ')}ms`);

        // No session should take excessively long
        expect(maxTime).toBeLessThan(5000);
      } finally {
        await Promise.all(contexts.map(context => context.close()));
      }
    });
  });

  test.describe('7. Device Performance Testing', () => {
    DEVICE_CONFIGS.forEach(device => {
      test(`should perform well on ${device.name}`, async ({ page }) => {
        await page.setViewportSize(device.viewport);

        // Simulate device constraints
        await page.evaluate((constraints) => {
          // @ts-ignore - performance.memory
          if ((performance as any).memory) {
            // Simulate memory constraints (this is illustrative)
            console.log(`Simulating ${constraints.memory}GB memory`);
          }
        }, device);

        const startTime = performance.now();

        await page.goto('/content');
        await page.waitForLoadState('networkidle');

        const loadTime = performance.now() - startTime;

        // Adjust expectations based on device tier
        const expectedLoadTime = device.cpu >= 4 ? 1000 : device.cpu >= 2 ? 2000 : 3000;

        console.log(`${device.name} load time: ${loadTime.toFixed(2)}ms`);
        expect(loadTime).toBeLessThan(expectedLoadTime);

        // Test content interactions
        const operationStart = performance.now();
        await page.click('[data-testid="content-list-button"]');
        await page.waitForSelector('[data-testid="content-item"]');
        const operationTime = performance.now() - operationStart;

        const expectedOperationTime = device.cpu >= 4 ? 200 : device.cpu >= 2 ? 400 : 800;

        console.log(`${device.name} operation time: ${operationTime.toFixed(2)}ms`);
        expect(operationTime).toBeLessThan(expectedOperationTime);
      });
    });

    test('should adapt to different network conditions', async ({ page, context }) => {
      for (const network of NETWORK_CONDITIONS) {
        if (network.name === 'Offline') continue; // Skip offline test here

        // Simulate network conditions
        await context.route('**/*', route => {
          setTimeout(() => route.continue(), network.latency);
        });

        const startTime = performance.now();

        await page.goto('/content');
        await page.waitForLoadState('networkidle');

        const loadTime = performance.now() - startTime;

        // Adjust expectations based on network
        const expectedTime = network.name === 'Fast 3G' ? 3000 : 6000;

        console.log(`${network.name} load time: ${loadTime.toFixed(2)}ms`);
        expect(loadTime).toBeLessThan(expectedTime);

        // Clear routes for next iteration
        await context.unroute('**/*');
      }
    });

    test('should handle touch interactions efficiently on mobile', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.goto('/content');
      await page.waitForLoadState('networkidle');

      // Test touch interactions
      const touchTargets = [
        '[data-testid="content-list-button"]',
        '[data-testid="create-content-button"]',
        '[data-testid="search-input"]',
        '[data-testid="filter-button"]',
      ];

      for (const target of touchTargets) {
        const element = page.locator(target);
        if (await element.isVisible()) {
          const startTime = performance.now();

          await element.tap();

          const responseTime = performance.now() - startTime;

          console.log(`Touch response time for ${target}: ${responseTime.toFixed(2)}ms`);
          expect(responseTime).toBeLessThan(100); // Should respond to touch quickly
        }
      }
    });
  });

  test.describe('Performance Regression Detection', () => {
    test('should establish performance baseline', async ({ page }) => {
      const metrics = {
        pageLoad: 0,
        contentLoad: 0,
        memoryUsage: 0,
        bundleSize: 0,
        apiResponseTime: 0,
      };

      // Measure page load time
      const pageStart = performance.now();
      await page.goto('/content');
      await page.waitForLoadState('networkidle');
      metrics.pageLoad = performance.now() - pageStart;

      // Measure content load time
      const contentStart = performance.now();
      await page.click('[data-testid="content-list-button"]');
      await page.waitForSelector('[data-testid="content-item"]');
      metrics.contentLoad = performance.now() - contentStart;

      // Measure memory usage
      const memory = await measureMemoryUsage(page);
      if (memory) {
        metrics.memoryUsage = memory.usedJSHeapSize;
      }

      // Measure Core Web Vitals
      const vitals = await measureCoreWebVitals(page);

      // Store baseline metrics (in real implementation, this would go to a database)
      console.log('Performance Baseline:', {
        ...metrics,
        vitals,
        timestamp: new Date().toISOString(),
      });

      // Validate against budgets
      expect(metrics.pageLoad).toBeLessThan(3000);
      expect(metrics.contentLoad).toBeLessThan(PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME);
      expect(vitals.LCP).toBeLessThan(PERFORMANCE_BUDGETS.LCP_THRESHOLD);
      expect(vitals.CLS).toBeLessThan(PERFORMANCE_BUDGETS.CLS_THRESHOLD);
    });
  });
});

test.describe('Content Performance - Edge Cases', () => {
  test('should handle large content datasets efficiently', async ({ page }) => {
    // Mock large dataset
    await page.route('**/api/content/items*', route => {
      const largeDataset = Array.from({ length: 1000 }, (_, i) => ({
        id: `content-${i}`,
        category: `category-${i % 10}`,
        type: 'text',
        value: `Large content item ${i} with substantial text content to simulate real-world usage patterns and test performance under load`,
        tags: [`tag-${i % 5}`, `tag-${i % 3}`],
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      }));

      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items: largeDataset,
          pagination: {
            total: 1000,
            page: 1,
            limit: 1000,
            hasNext: false,
            hasPrev: false,
          },
        }),
      });
    });

    await page.goto('/content');
    await page.waitForLoadState('networkidle');

    const startTime = performance.now();
    await page.click('[data-testid="content-list-button"]');
    await page.waitForSelector('[data-testid="content-item"]');
    const loadTime = performance.now() - startTime;

    console.log(`Large dataset load time: ${loadTime.toFixed(2)}ms`);

    // Should handle large datasets within reasonable time
    expect(loadTime).toBeLessThan(1000);

    // Test scrolling performance
    const scrollStart = performance.now();
    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
    const scrollTime = performance.now() - scrollStart;

    console.log(`Scroll performance: ${scrollTime.toFixed(2)}ms`);
    expect(scrollTime).toBeLessThan(200);
  });

  test('should handle rapid user interactions without degradation', async ({ page }) => {
    await page.goto('/content');
    await page.waitForLoadState('networkidle');

    const interactions = [];

    // Rapid fire interactions
    for (let i = 0; i < 20; i++) {
      const startTime = performance.now();

      // Alternate between different interactions
      if (i % 2 === 0) {
        await page.click('[data-testid="content-list-button"]');
      } else {
        await page.click('[data-testid="create-content-button"]');
        await page.press('Escape');
      }

      const endTime = performance.now();
      interactions.push(endTime - startTime);
    }

    const avgTime = interactions.reduce((sum, time) => sum + time, 0) / interactions.length;
    const maxTime = Math.max(...interactions);

    console.log(`Rapid interactions - Average: ${avgTime.toFixed(2)}ms, Max: ${maxTime.toFixed(2)}ms`);

    // Performance should not degrade significantly
    expect(avgTime).toBeLessThan(150);
    expect(maxTime).toBeLessThan(300);
  });
});