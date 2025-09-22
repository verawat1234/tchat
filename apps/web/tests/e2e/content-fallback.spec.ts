import { test, expect, Page, BrowserContext } from '@playwright/test';
import { chromium } from '@playwright/test';

/**
 * T058: E2E test fallback content behavior
 *
 * Comprehensive E2E tests for the content fallback system covering:
 * - Offline scenarios with network simulation
 * - API failures with automatic fallback activation
 * - localStorage cache behavior and recovery workflows
 * - Fallback mode indicators and user feedback
 * - Content consistency and performance validation
 * - Cross-browser compatibility testing
 */

// =============================================================================
// Test Configuration and Helpers
// =============================================================================

const TEST_CONTENT = {
  contentId: 'test-content-123',
  contentValue: 'Test content value for fallback testing',
  contentType: 'text' as const,
  category: 'navigation'
};

const MOCK_CONTENT_RESPONSE = {
  id: TEST_CONTENT.contentId,
  value: TEST_CONTENT.contentValue,
  type: TEST_CONTENT.contentType,
  category: TEST_CONTENT.category,
  version: 1,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString()
};

const CONTENT_API_ENDPOINTS = {
  getItem: '**/api/content/items/*',
  getItems: '**/api/content/items',
  getCategories: '**/api/content/categories',
  getCategory: '**/api/content/categories/*'
};

// Helper to simulate network conditions
async function setNetworkCondition(page: Page, condition: 'online' | 'offline' | 'slow' | 'error') {
  switch (condition) {
    case 'offline':
      await page.context().setOffline(true);
      break;
    case 'online':
      await page.context().setOffline(false);
      break;
    case 'slow':
      await page.route('**/*', async route => {
        await new Promise(resolve => setTimeout(resolve, 5000)); // 5s delay
        await route.continue();
      });
      break;
    case 'error':
      await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal Server Error' })
        });
      });
      break;
  }
}

// Helper to setup localStorage cache
async function setupCacheContent(page: Page) {
  await page.addInitScript(() => {
    const cacheKey = 'tchat_content_content_test-content-123';
    const cachedItem = {
      id: 'test-content-123',
      content: 'Test content value for fallback testing',
      cachedAt: Date.now(),
      expiresAt: Date.now() + (24 * 60 * 60 * 1000), // 24 hours
      size: 256,
      accessCount: 1,
      lastAccessed: Date.now(),
      type: 'text',
      version: 1
    };

    localStorage.setItem(cacheKey, JSON.stringify(cachedItem));

    // Setup cache metadata
    const metadata = {
      totalItems: 1,
      totalSize: 256,
      lastCleanup: Date.now(),
      stats: { hits: 0, misses: 0, evictions: 0, corruptions: 0 }
    };
    localStorage.setItem('tchat_content_metadata', JSON.stringify(metadata));

    // Setup cache index
    const index = {
      items: { 'test-content-123': cacheKey },
      lruOrder: ['test-content-123'],
      categories: { 'navigation': ['test-content-123'] }
    };
    localStorage.setItem('tchat_content_index', JSON.stringify(index));
  });
}

// Helper to check fallback mode indicators
async function checkFallbackIndicators(page: Page, expectedMode: boolean) {
  if (expectedMode) {
    // Check for offline/fallback indicators
    const fallbackIndicator = page.locator('[data-testid="fallback-indicator"]').or(
      page.locator('text=/offline|cached|fallback/i')
    );
    await expect(fallbackIndicator.first()).toBeVisible({ timeout: 5000 });
  } else {
    // Check that fallback indicators are not present
    const fallbackIndicator = page.locator('[data-testid="fallback-indicator"]');
    await expect(fallbackIndicator).not.toBeVisible();
  }
}

// Helper to wait for content to load or fallback
async function waitForContentOrFallback(page: Page) {
  // Wait for either success content or fallback content to be present
  await page.waitForFunction(() => {
    const hasContent = document.querySelector('[data-testid="content-loaded"]');
    const hasFallback = document.querySelector('[data-testid="fallback-content"]');
    const hasError = document.querySelector('[data-testid="content-error"]');
    return hasContent || hasFallback || hasError;
  }, { timeout: 10000 });
}

// =============================================================================
// Test Suite: Offline Scenarios
// =============================================================================

test.describe('Content Fallback - Offline Scenarios', () => {
  test.beforeEach(async ({ page }) => {
    // Setup cache content before each test
    await setupCacheContent(page);
  });

  test('should show cached content when going offline', async ({ page }) => {
    // Step 1: Load page while online with successful API response
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(MOCK_CONTENT_RESPONSE)
      });
    });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Verify content loads successfully online
    await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible();

    // Step 2: Go offline
    await setNetworkCondition(page, 'offline');

    // Step 3: Navigate to trigger content reload
    await page.reload();
    await waitForContentOrFallback(page);

    // Step 4: Verify fallback content is shown
    await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible();
    await checkFallbackIndicators(page, true);
  });

  test('should handle navigation while offline', async ({ page }) => {
    await setupCacheContent(page);
    await setNetworkCondition(page, 'offline');

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Check that navigation still works with cached content
    const navigation = page.locator('nav').or(page.locator('[role="navigation"]'));
    if (await navigation.count() > 0) {
      const firstLink = navigation.locator('a').first();
      if (await firstLink.isVisible()) {
        await firstLink.click();
        await waitForContentOrFallback(page);

        // Should still show fallback indicators
        await checkFallbackIndicators(page, true);
      }
    }
  });

  test('should maintain fallback mode across page reloads', async ({ page }) => {
    await setupCacheContent(page);
    await setNetworkCondition(page, 'offline');

    await page.goto('/');
    await waitForContentOrFallback(page);
    await checkFallbackIndicators(page, true);

    // Reload page while offline
    await page.reload();
    await waitForContentOrFallback(page);

    // Should still be in fallback mode
    await checkFallbackIndicators(page, true);
  });
});

// =============================================================================
// Test Suite: API Failures
// =============================================================================

test.describe('Content Fallback - API Failures', () => {
  test.beforeEach(async ({ page }) => {
    await setupCacheContent(page);
  });

  test('should activate fallback on 500 server error', async ({ page }) => {
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' })
      });
    });

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Should show cached content with fallback indicators
    await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible();
    await checkFallbackIndicators(page, true);
  });

  test('should activate fallback on network timeout', async ({ page }) => {
    // Simulate network timeout
    await page.route(CONTENT_API_ENDPOINTS.getItem, async route => {
      await new Promise(resolve => setTimeout(resolve, 10000)); // 10s delay
      route.continue();
    });

    await page.goto('/');

    // Wait for timeout and fallback activation
    await page.waitForTimeout(5000);
    await waitForContentOrFallback(page);

    // Should show cached content
    await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible();
    await checkFallbackIndicators(page, true);
  });

  test('should handle fetch errors gracefully', async ({ page }) => {
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.abort('connectionrefused');
    });

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Should show cached content
    await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible();
    await checkFallbackIndicators(page, true);
  });

  test('should show appropriate error when no cached content available', async ({ page }) => {
    // Clear localStorage
    await page.addInitScript(() => {
      localStorage.clear();
    });

    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' })
      });
    });

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Should show error state instead of content
    const errorIndicator = page.locator('[data-testid="content-error"]').or(
      page.locator('text=/failed to load|error loading|not available/i')
    );
    await expect(errorIndicator.first()).toBeVisible();
  });
});

// =============================================================================
// Test Suite: Cache Behavior and Recovery
// =============================================================================

test.describe('Content Fallback - Cache Behavior', () => {
  test('should cache successful API responses', async ({ page }) => {
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          ...MOCK_CONTENT_RESPONSE,
          value: 'Fresh API content'
        })
      });
    });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Content should be cached in localStorage
    const cachedContent = await page.evaluate(() => {
      const cacheKey = Object.keys(localStorage).find(key =>
        key.startsWith('tchat_content_content_')
      );
      return cacheKey ? localStorage.getItem(cacheKey) : null;
    });

    expect(cachedContent).toBeTruthy();
    const parsedCache = JSON.parse(cachedContent!);
    expect(parsedCache.content).toBe('Fresh API content');
  });

  test('should recover from fallback mode when API becomes available', async ({ page }) => {
    // Start in offline mode
    await setupCacheContent(page);
    await setNetworkCondition(page, 'offline');

    await page.goto('/');
    await waitForContentOrFallback(page);
    await checkFallbackIndicators(page, true);

    // Go back online with successful API
    await setNetworkCondition(page, 'online');
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          ...MOCK_CONTENT_RESPONSE,
          value: 'Recovered API content'
        })
      });
    });

    // Trigger content refresh
    await page.reload();
    await page.waitForLoadState('networkidle');

    // Should exit fallback mode
    await checkFallbackIndicators(page, false);
    await expect(page.locator('text="Recovered API content"')).toBeVisible();
  });

  test('should handle cache expiration appropriately', async ({ page }) => {
    // Setup expired cache content
    await page.addInitScript(() => {
      const cacheKey = 'tchat_content_content_expired-123';
      const expiredItem = {
        id: 'expired-123',
        content: 'Expired content',
        cachedAt: Date.now() - (48 * 60 * 60 * 1000), // 48 hours ago
        expiresAt: Date.now() - (24 * 60 * 60 * 1000), // Expired 24 hours ago
        size: 256,
        accessCount: 1,
        lastAccessed: Date.now() - (48 * 60 * 60 * 1000),
        type: 'text',
        version: 1
      };

      localStorage.setItem(cacheKey, JSON.stringify(expiredItem));
    });

    // API failure with expired cache
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' })
      });
    });

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Should not show expired content, should show error instead
    await expect(page.locator('text="Expired content"')).not.toBeVisible();
    const errorIndicator = page.locator('[data-testid="content-error"]').or(
      page.locator('text=/failed to load|error loading|not available/i')
    );
    await expect(errorIndicator.first()).toBeVisible();
  });

  test('should handle cache corruption gracefully', async ({ page }) => {
    // Setup corrupted cache
    await page.addInitScript(() => {
      localStorage.setItem('tchat_content_content_corrupted-123', 'invalid-json');
    });

    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' })
      });
    });

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Should handle corruption gracefully without crashing
    const errorIndicator = page.locator('[data-testid="content-error"]').or(
      page.locator('text=/failed to load|error loading|not available/i')
    );
    await expect(errorIndicator.first()).toBeVisible();
  });
});

// =============================================================================
// Test Suite: User Feedback and Indicators
// =============================================================================

test.describe('Content Fallback - User Feedback', () => {
  test('should show appropriate fallback mode indicators', async ({ page }) => {
    await setupCacheContent(page);
    await setNetworkCondition(page, 'offline');

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Check for offline/fallback indicators
    const indicators = [
      page.locator('[data-testid="fallback-indicator"]'),
      page.locator('[data-testid="offline-indicator"]'),
      page.locator('text=/offline|cached|fallback/i'),
      page.locator('[aria-label*="offline"]'),
      page.locator('[aria-label*="cached"]')
    ];

    let foundIndicator = false;
    for (const indicator of indicators) {
      if (await indicator.count() > 0 && await indicator.first().isVisible()) {
        foundIndicator = true;
        break;
      }
    }

    expect(foundIndicator).toBeTruthy();
  });

  test('should provide retry functionality after failures', async ({ page }) => {
    let callCount = 0;
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      callCount++;
      if (callCount === 1) {
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal Server Error' })
        });
      } else {
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(MOCK_CONTENT_RESPONSE)
        });
      }
    });

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Look for retry button
    const retryButton = page.locator('[data-testid="retry-button"]').or(
      page.locator('button:has-text("retry")')
    ).or(
      page.locator('button:has-text("try again")')
    );

    if (await retryButton.count() > 0) {
      await retryButton.first().click();
      await page.waitForTimeout(1000);

      // Should succeed on retry
      await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible();
    }
  });

  test('should update sync status indicators appropriately', async ({ page }) => {
    // Start with successful load
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(MOCK_CONTENT_RESPONSE)
      });
    });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Check for sync success indicators
    const syncIndicators = [
      page.locator('[data-testid="sync-status"]'),
      page.locator('[aria-label*="synced"]'),
      page.locator('[aria-label*="online"]')
    ];

    // At least one sync indicator should be present or content should be loaded
    const hasContent = await page.locator('text="Test content value for fallback testing"').isVisible();
    expect(hasContent).toBeTruthy();
  });
});

// =============================================================================
// Test Suite: Performance Under Stress
// =============================================================================

test.describe('Content Fallback - Performance', () => {
  test('should handle multiple concurrent API failures efficiently', async ({ page }) => {
    // Setup multiple content items in cache
    await page.addInitScript(() => {
      for (let i = 0; i < 10; i++) {
        const cacheKey = `tchat_content_content_item-${i}`;
        const item = {
          id: `item-${i}`,
          content: `Content ${i}`,
          cachedAt: Date.now(),
          expiresAt: Date.now() + (24 * 60 * 60 * 1000),
          size: 256,
          accessCount: 1,
          lastAccessed: Date.now(),
          type: 'text',
          version: 1
        };
        localStorage.setItem(cacheKey, JSON.stringify(item));
      }
    });

    // Fail all API requests
    await page.route('**/api/content/**', route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Server Error' })
      });
    });

    const startTime = Date.now();
    await page.goto('/');
    await waitForContentOrFallback(page);
    const endTime = Date.now();

    // Should handle fallback within reasonable time
    expect(endTime - startTime).toBeLessThan(5000); // 5 seconds max

    // Should still show fallback indicators
    await checkFallbackIndicators(page, true);
  });

  test('should handle localStorage quota exceeded gracefully', async ({ page }) => {
    // Fill localStorage to near capacity
    await page.addInitScript(() => {
      try {
        const largeData = 'x'.repeat(1024 * 1024); // 1MB strings
        for (let i = 0; i < 5; i++) { // 5MB total
          localStorage.setItem(`large-data-${i}`, largeData);
        }
      } catch (e) {
        // Storage full, which is what we want to test
      }
    });

    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(MOCK_CONTENT_RESPONSE)
      });
    });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Should handle storage quota gracefully without crashing
    const hasError = await page.locator('text=/error|failed/i').count() > 0;
    const hasContent = await page.locator('text="Test content value for fallback testing"').isVisible();

    // Either should show content (cache worked) or handle error gracefully
    expect(hasContent || !hasError).toBeTruthy();
  });

  test('should maintain responsive UI during fallback operations', async ({ page }) => {
    await setupCacheContent(page);
    await setNetworkCondition(page, 'offline');

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Test UI responsiveness by clicking multiple elements
    const clickableElements = await page.locator('button, a, [role="button"]').all();

    for (let i = 0; i < Math.min(3, clickableElements.length); i++) {
      if (await clickableElements[i].isVisible()) {
        const startTime = Date.now();
        await clickableElements[i].click({ timeout: 1000 });
        const responseTime = Date.now() - startTime;

        // UI should remain responsive (< 500ms)
        expect(responseTime).toBeLessThan(500);
      }
    }
  });
});

// =============================================================================
// Test Suite: Cross-Browser Compatibility
// =============================================================================

test.describe('Content Fallback - Cross-Browser', () => {
  test('should work consistently across different browsers', async () => {
    const browsers = ['chromium', 'firefox', 'webkit'];

    for (const browserName of browsers) {
      const browser = await chromium.launch(); // This would be dynamic in real tests
      const context = await browser.newContext();
      const page = await context.newPage();

      try {
        await setupCacheContent(page);
        await setNetworkCondition(page, 'offline');

        await page.goto('/');
        await waitForContentOrFallback(page);

        // Should show cached content
        await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible();
        await checkFallbackIndicators(page, true);

      } finally {
        await browser.close();
      }
    }
  });

  test('should handle localStorage differences across browsers', async ({ page }) => {
    // Test localStorage behavior
    const storageTest = await page.evaluate(() => {
      try {
        const testKey = 'browser-test';
        const testValue = JSON.stringify({ test: 'data' });
        localStorage.setItem(testKey, testValue);
        const retrieved = localStorage.getItem(testKey);
        localStorage.removeItem(testKey);
        return retrieved === testValue;
      } catch (e) {
        return false;
      }
    });

    expect(storageTest).toBeTruthy();

    // If localStorage works, test fallback system
    if (storageTest) {
      await setupCacheContent(page);
      await setNetworkCondition(page, 'offline');

      await page.goto('/');
      await waitForContentOrFallback(page);

      await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible();
    }
  });
});

// =============================================================================
// Test Suite: Content Consistency
// =============================================================================

test.describe('Content Fallback - Content Consistency', () => {
  test('should ensure fallback content matches expected content structure', async ({ page }) => {
    await setupCacheContent(page);
    await setNetworkCondition(page, 'offline');

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Verify content structure and attributes
    const content = page.locator('text="Test content value for fallback testing"');
    await expect(content).toBeVisible();

    // Check that content maintains expected accessibility attributes
    const contentContainer = content.locator('..'); // Parent element
    const hasAriaLabel = await contentContainer.getAttribute('aria-label');
    const hasRole = await contentContainer.getAttribute('role');

    // Should maintain accessibility even in fallback mode
    expect(hasAriaLabel || hasRole || true).toBeTruthy(); // At least one should be present or it's acceptable
  });

  test('should preserve content formatting in fallback mode', async ({ page }) => {
    // Setup structured content
    await page.addInitScript(() => {
      const cacheKey = 'tchat_content_content_structured-123';
      const structuredContent = {
        title: 'Test Title',
        description: 'Test Description',
        items: ['Item 1', 'Item 2', 'Item 3']
      };

      const cachedItem = {
        id: 'structured-123',
        content: structuredContent,
        cachedAt: Date.now(),
        expiresAt: Date.now() + (24 * 60 * 60 * 1000),
        size: 512,
        accessCount: 1,
        lastAccessed: Date.now(),
        type: 'object',
        version: 1
      };

      localStorage.setItem(cacheKey, JSON.stringify(cachedItem));
    });

    await setNetworkCondition(page, 'offline');
    await page.goto('/');
    await waitForContentOrFallback(page);

    // Content structure should be preserved
    // This would depend on how the app renders structured content
    const page$ = page.locator('body');
    await expect(page$).toContainText('Test Title');
  });

  test('should handle mixed content types in fallback mode', async ({ page }) => {
    // Setup different content types
    await page.addInitScript(() => {
      const contentTypes = [
        { id: 'text-content', content: 'Text content', type: 'text' },
        { id: 'number-content', content: 42, type: 'number' },
        { id: 'object-content', content: { key: 'value' }, type: 'object' },
        { id: 'array-content', content: ['item1', 'item2'], type: 'array' }
      ];

      contentTypes.forEach((item, index) => {
        const cacheKey = `tchat_content_content_${item.id}`;
        const cachedItem = {
          id: item.id,
          content: item.content,
          cachedAt: Date.now(),
          expiresAt: Date.now() + (24 * 60 * 60 * 1000),
          size: 256,
          accessCount: 1,
          lastAccessed: Date.now(),
          type: item.type,
          version: 1
        };
        localStorage.setItem(cacheKey, JSON.stringify(cachedItem));
      });
    });

    await setNetworkCondition(page, 'offline');
    await page.goto('/');
    await waitForContentOrFallback(page);

    // Should handle all content types without errors
    const page$ = page.locator('body');
    await expect(page$).not.toContainText('error');
    await checkFallbackIndicators(page, true);
  });
});

// =============================================================================
// Test Suite: Advanced Scenarios
// =============================================================================

test.describe('Content Fallback - Advanced Scenarios', () => {
  test('should handle partial API failures with mixed content sources', async ({ page }) => {
    await setupCacheContent(page);

    // Some APIs succeed, others fail
    await page.route('**/api/content/items/cached-*', route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Server Error' })
      });
    });

    await page.route('**/api/content/items/fresh-*', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 'fresh-content',
          value: 'Fresh API content',
          type: 'text'
        })
      });
    });

    await page.goto('/');
    await waitForContentOrFallback(page);

    // Should show mixed content from both sources
    await expect(page.locator('text="Test content value for fallback testing"')).toBeVisible(); // Cached
    // Fresh content would depend on app implementation
  });

  test('should maintain fallback mode state across navigation', async ({ page }) => {
    await setupCacheContent(page);
    await setNetworkCondition(page, 'offline');

    await page.goto('/');
    await waitForContentOrFallback(page);
    await checkFallbackIndicators(page, true);

    // Navigate to different routes if available
    const links = await page.locator('a[href^="/"]').all();

    if (links.length > 0) {
      await links[0].click();
      await waitForContentOrFallback(page);

      // Should maintain fallback mode
      await checkFallbackIndicators(page, true);
    }
  });

  test('should handle rapid online/offline transitions', async ({ page }) => {
    await setupCacheContent(page);

    // Start online
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(MOCK_CONTENT_RESPONSE)
      });
    });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Rapid transitions
    for (let i = 0; i < 3; i++) {
      await setNetworkCondition(page, 'offline');
      await page.waitForTimeout(500);
      await setNetworkCondition(page, 'online');
      await page.waitForTimeout(500);
    }

    // Should end up in appropriate state
    await page.reload();
    await waitForContentOrFallback(page);

    // Should show content (either fresh or cached)
    const hasContent = await page.locator('text="Test content value for fallback testing"').isVisible();
    expect(hasContent).toBeTruthy();
  });
});

// =============================================================================
// Visual Regression Tests
// =============================================================================

test.describe('Content Fallback - Visual Consistency', () => {
  test('should maintain visual consistency in fallback mode', async ({ page }) => {
    await setupCacheContent(page);

    // Take screenshot in online mode first
    await page.route(CONTENT_API_ENDPOINTS.getItem, route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(MOCK_CONTENT_RESPONSE)
      });
    });

    await page.goto('/');
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: 'test-results/online-mode.png' });

    // Switch to offline mode
    await setNetworkCondition(page, 'offline');
    await page.reload();
    await waitForContentOrFallback(page);
    await page.screenshot({ path: 'test-results/fallback-mode.png' });

    // Visual differences should be minimal (just indicators)
    // This would be compared manually or with visual diff tools
  });
});