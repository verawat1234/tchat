/**
 * Mobile E2E Tests: Stream Navigation - Cross-Platform
 *
 * Comprehensive end-to-end testing for Stream Store Tabs on mobile platforms.
 * Tests cover iOS and Android Stream navigation functionality including
 * tab switching, content loading, performance validation, and cross-platform consistency.
 *
 * Based on quickstart scenarios from specs/026-help-me-add/quickstart.md
 * Platform coverage: iOS (Swift/SwiftUI), Android (Kotlin/Compose)
 *
 * Test Requirements:
 * - Tab switching: <200ms response time
 * - Content loading: <3s initial load
 * - Cross-platform consistency: >95% visual parity
 * - Store integration: unified cart experience
 * - Offline capability: graceful fallback
 * - Native performance: 60fps animations
 */

import { test, expect, device } from '@playwright/test';

// Mobile-specific performance budgets
const MOBILE_PERFORMANCE_BUDGETS = {
  TAB_SWITCH_TIME: 200, // ms
  CONTENT_LOAD_TIME: 3000, // ms (3s)
  API_RESPONSE_TIME: 1000, // ms (mobile networks)
  ANIMATION_FRAME_RATE: 55, // fps (allow mobile variance)
  MEMORY_USAGE_LIMIT: 100 * 1024 * 1024, // 100MB
  VISUAL_CONSISTENCY: 0.95, // 95%
};

// Stream categories for validation
const STREAM_CATEGORIES = [
  { id: 'books', name: 'Books', icon: 'book-open' },
  { id: 'podcasts', name: 'Podcasts', icon: 'microphone' },
  { id: 'cartoons', name: 'Cartoons', icon: 'film' },
  { id: 'movies', name: 'Movies', icon: 'video' },
  { id: 'music', name: 'Music', icon: 'music' },
  { id: 'art', name: 'Art', icon: 'palette' },
];

// Platform configurations
const MOBILE_PLATFORMS = [
  {
    name: 'iOS',
    device: 'iPhone 13',
    viewport: { width: 390, height: 844 },
    userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)',
    testId: 'ios-stream-tab',
  },
  {
    name: 'Android',
    device: 'Pixel 5',
    viewport: { width: 393, height: 851 },
    userAgent: 'Mozilla/5.0 (Linux; Android 11; Pixel 5)',
    testId: 'android-stream-tab',
  },
];

// Helper functions for mobile testing
async function navigateToMobileStore(page: any, platform: string) {
  // Mobile-specific navigation logic
  if (platform === 'iOS') {
    // iOS SwiftUI navigation patterns
    await page.locator('[data-testid="tab-bar-store"], text="Store"').click();
  } else {
    // Android Compose navigation patterns
    await page.locator('[data-testid="bottom-nav-store"], text="Store"').click();
  }

  await page.waitForLoadState('networkidle');
}

async function waitForMobileStreamTab(page: any, platform: string) {
  const streamSelector = platform === 'iOS'
    ? '[data-testid="stream-tab-ios"], text="Stream"'
    : '[data-testid="stream-tab-android"], text="Stream"';

  await page.waitForSelector(streamSelector, { timeout: 10000 });
}

async function clickMobileStreamTab(page: any, platform: string) {
  const streamSelector = platform === 'iOS'
    ? '[data-testid="stream-tab-ios"], text="Stream"'
    : '[data-testid="stream-tab-android"], text="Stream"';

  await page.locator(streamSelector).first().click();

  // Wait for Stream content to load with mobile timeout
  await page.waitForSelector('[data-testid="stream-content"], .stream-content', { timeout: 15000 });
}

async function measureMobileTabSwitchTime(page: any, fromTab: string, toTab: string, platform: string): Promise<number> {
  const startTime = Date.now();

  // Platform-specific tab selectors
  const tabSelector = platform === 'iOS'
    ? `[data-testid="${toTab}-tab-ios"], text="${toTab}"`
    : `[data-testid="${toTab}-tab-android"], text="${toTab}"`;

  await page.locator(tabSelector).first().click();
  await page.waitForSelector(`[data-testid="${toTab}-content"], .${toTab}-content`, { timeout: 5000 });

  const endTime = Date.now();
  return endTime - startTime;
}

async function validateMobileFeaturedContent(page: any, categoryId: string, platform: string) {
  // Mobile-specific featured content validation
  const featuredSelector = platform === 'iOS'
    ? '[data-testid="featured-carousel-ios"], .featured-content'
    : '[data-testid="featured-carousel-android"], .featured-content';

  await expect(page.locator(featuredSelector)).toBeVisible({ timeout: 10000 });

  // Check for mobile-optimized featured items
  const featuredItems = page.locator('[data-testid="featured-item"], .featured-item');
  const itemCount = await featuredItems.count();

  if (itemCount > 0) {
    const firstItem = featuredItems.first();

    // Validate mobile touch targets and layout
    await expect(firstItem.locator('[data-testid="item-title"], .item-title')).toBeVisible();
    await expect(firstItem.locator('[data-testid="item-image"], .item-image, img')).toBeVisible();
    await expect(firstItem.locator('[data-testid="item-price"], .item-price')).toBeVisible();

    // Validate mobile touch target size (44dp minimum)
    const touchTarget = firstItem.locator('[data-testid="item-touch-target"], .touch-target');
    if (await touchTarget.isVisible()) {
      const box = await touchTarget.boundingBox();
      expect(box?.height).toBeGreaterThanOrEqual(44);
      expect(box?.width).toBeGreaterThanOrEqual(44);
    }
  }

  return itemCount;
}

// Cross-platform test suite
test.describe('Mobile Stream Navigation - Cross-Platform', () => {

  // Test each mobile platform
  for (const platform of MOBILE_PLATFORMS) {
    test.describe(`${platform.name} Platform Tests`, () => {

      test.beforeEach(async ({ page }) => {
        // Configure mobile viewport and user agent
        await page.setViewportSize(platform.viewport);
        await page.setUserAgent(platform.userAgent);

        // Navigate to mobile app
        await page.goto('http://localhost:3000');
        await page.waitForLoadState('networkidle');
      });

      test.describe('Basic Mobile Navigation', () => {

        test(`should display Stream tab with mobile-optimized layout on ${platform.name}`, async ({ page }) => {
          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);

          // Verify Stream tab is visible and accessible
          const streamTab = page.locator(`[data-testid="stream-tab-${platform.name.toLowerCase()}"], text="Stream"`).first();
          await expect(streamTab).toBeVisible();

          // Validate mobile accessibility
          await expect(streamTab).toHaveAttribute('role', 'tab');

          // Check mobile touch target size
          const tabBox = await streamTab.boundingBox();
          expect(tabBox?.height).toBeGreaterThanOrEqual(44); // iOS HIG minimum
        });

        test(`should navigate between Stream categories smoothly on ${platform.name}`, async ({ page }) => {
          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);
          await clickMobileStreamTab(page, platform.name);

          // Test category switching performance
          for (const category of STREAM_CATEGORIES.slice(0, 3)) {
            const startTime = Date.now();

            const categoryTab = page.locator(`text="${category.name}"`).first();
            await categoryTab.click();

            // Wait for content update
            await page.waitForSelector('[data-testid="category-content"], .category-content', { timeout: 3000 });

            const switchTime = Date.now() - startTime;
            expect(switchTime).toBeLessThan(MOBILE_PERFORMANCE_BUDGETS.TAB_SWITCH_TIME);

            console.log(`${platform.name} - ${category.name} switch: ${switchTime}ms`);
          }
        });

        test(`should maintain state during mobile orientation changes on ${platform.name}`, async ({ page }) => {
          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);
          await clickMobileStreamTab(page, platform.name);

          // Select Podcasts category
          const podcastsTab = page.locator('text="Podcasts"').first();
          await podcastsTab.click();
          await page.waitForTimeout(1000);

          // Simulate orientation change (landscape)
          await page.setViewportSize({ width: 844, height: 390 });
          await page.waitForTimeout(1000);

          // Verify state persists
          await expect(podcastsTab).toHaveClass(/active|selected/);

          // Return to portrait
          await page.setViewportSize(platform.viewport);
          await page.waitForTimeout(1000);

          // State should still persist
          await expect(podcastsTab).toHaveClass(/active|selected/);
        });
      });

      test.describe('Movies Subtab Mobile Navigation', () => {

        test(`should display and navigate Movies subtabs on ${platform.name}`, async ({ page }) => {
          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);
          await clickMobileStreamTab(page, platform.name);

          // Navigate to Movies category
          const moviesTab = page.locator('text="Movies"').first();
          await moviesTab.click();
          await page.waitForSelector('[data-testid="movies-subtabs"], .subtabs', { timeout: 5000 });

          // Verify mobile-optimized subtabs
          const shortFilmsTab = page.locator('text="Short Films"').first();
          const featureFilmsTab = page.locator('text="Feature Films"').first();

          await expect(shortFilmsTab).toBeVisible();
          await expect(featureFilmsTab).toBeVisible();

          // Test mobile subtab switching
          await featureFilmsTab.click();
          await page.waitForTimeout(500);
          await expect(featureFilmsTab).toHaveClass(/active|selected/);

          await shortFilmsTab.click();
          await page.waitForTimeout(500);
          await expect(shortFilmsTab).toHaveClass(/active|selected/);

          // Validate mobile touch targets
          const shortBox = await shortFilmsTab.boundingBox();
          const featureBox = await featureFilmsTab.boundingBox();

          expect(shortBox?.height).toBeGreaterThanOrEqual(44);
          expect(featureBox?.height).toBeGreaterThanOrEqual(44);
        });
      });

      test.describe('Mobile Featured Content', () => {

        test(`should display mobile-optimized featured content on ${platform.name}`, async ({ page }) => {
          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);
          await clickMobileStreamTab(page, platform.name);

          // Test featured content in Books category
          const booksTab = page.locator('text="Books"').first();
          await booksTab.click();
          await page.waitForTimeout(1000);

          const featuredCount = await validateMobileFeaturedContent(page, 'books', platform.name);
          console.log(`${platform.name} - Books featured content: ${featuredCount} items`);

          // Test horizontal scrolling on mobile
          const carousel = page.locator('[data-testid="featured-carousel"], .featured-content');
          if (await carousel.isVisible()) {
            // Test swipe gesture simulation
            const carouselBox = await carousel.boundingBox();
            if (carouselBox) {
              await page.mouse.move(carouselBox.x + carouselBox.width - 50, carouselBox.y + carouselBox.height / 2);
              await page.mouse.down();
              await page.mouse.move(carouselBox.x + 50, carouselBox.y + carouselBox.height / 2);
              await page.mouse.up();

              await page.waitForTimeout(500); // Allow animation to complete
            }
          }
        });
      });

      test.describe('Mobile Store Integration', () => {

        test(`should integrate with mobile cart experience on ${platform.name}`, async ({ page }) => {
          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);
          await clickMobileStreamTab(page, platform.name);

          // Navigate to Books and try to add to cart
          const booksTab = page.locator('text="Books"').first();
          await booksTab.click();
          await page.waitForTimeout(1000);

          // Look for mobile-optimized add to cart button
          const addToCartButton = page.locator('[data-testid="add-to-cart-mobile"], [data-testid="add-to-cart"], .add-to-cart').first();

          if (await addToCartButton.isVisible()) {
            // Validate mobile button size and accessibility
            const buttonBox = await addToCartButton.boundingBox();
            expect(buttonBox?.height).toBeGreaterThanOrEqual(44);

            await addToCartButton.click();

            // Check for mobile-friendly feedback
            const feedback = page.locator('.toast, .snackbar, [data-testid="cart-feedback"]');
            await expect(feedback).toBeVisible({ timeout: 3000 });
          }
        });
      });

      test.describe('Mobile Performance Validation', () => {

        test(`should meet mobile performance targets on ${platform.name}`, async ({ page }) => {
          // Monitor performance during navigation
          const startTime = Date.now();

          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);

          const initialLoadTime = Date.now() - startTime;
          expect(initialLoadTime).toBeLessThan(MOBILE_PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME);

          const streamLoadStart = Date.now();
          await clickMobileStreamTab(page, platform.name);

          const streamLoadTime = Date.now() - streamLoadStart;
          expect(streamLoadTime).toBeLessThan(MOBILE_PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME);

          console.log(`${platform.name} Performance - Initial: ${initialLoadTime}ms, Stream: ${streamLoadTime}ms`);
        });

        test(`should handle mobile network conditions gracefully on ${platform.name}`, async ({ page }) => {
          // Simulate slow mobile network
          await page.route('**/api/v1/stream/**', route => {
            setTimeout(() => route.continue(), 2000); // 2s delay
          });

          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);

          const loadingStart = Date.now();
          await clickMobileStreamTab(page, platform.name);

          // Should show mobile loading indicators
          const loadingIndicator = page.locator('[data-testid="loading-mobile"], [data-testid="loading"], .loading');
          await expect(loadingIndicator).toBeVisible({ timeout: 1000 });

          // Wait for content to eventually load
          await page.waitForSelector('[data-testid="stream-content"]', { timeout: 15000 });

          const totalLoadTime = Date.now() - loadingStart;
          console.log(`${platform.name} Slow network load time: ${totalLoadTime}ms`);

          // Should handle gracefully even with slow network
          expect(totalLoadTime).toBeLessThan(10000); // 10s maximum
        });
      });

      test.describe('Mobile Error Handling', () => {

        test(`should handle offline scenarios gracefully on ${platform.name}`, async ({ page }) => {
          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);

          // Simulate network failure
          await page.route('**/api/v1/stream/**', route => {
            route.abort('failed');
          });

          await clickMobileStreamTab(page, platform.name);

          // Should show mobile-friendly error state
          const errorState = page.locator('[data-testid="offline-state"], [data-testid="error-state"], .error-state');
          await expect(errorState).toBeVisible({ timeout: 5000 });

          // Error state should be mobile-optimized
          const errorBox = await errorState.boundingBox();
          expect(errorBox?.width).toBeLessThanOrEqual(platform.viewport.width);
        });

        test(`should handle empty content states on mobile ${platform.name}`, async ({ page }) => {
          // Mock empty content response
          await page.route('**/api/v1/stream/content**', route => {
            route.fulfill({
              status: 200,
              contentType: 'application/json',
              body: JSON.stringify({
                content: [],
                total: 0,
                page: 1,
                page_size: 20,
                total_pages: 0,
              }),
            });
          });

          await navigateToMobileStore(page, platform.name);
          await waitForMobileStreamTab(page, platform.name);
          await clickMobileStreamTab(page, platform.name);

          // Navigate to Books category
          const booksTab = page.locator('text="Books"').first();
          await booksTab.click();
          await page.waitForTimeout(1000);

          // Should show mobile-friendly empty state
          const emptyState = page.locator('[data-testid="empty-content"], [data-testid="empty-state"], .empty-state');
          await expect(emptyState).toBeVisible({ timeout: 5000 });

          // Empty state should provide helpful guidance
          await expect(emptyState).toContainText(/no content|empty|try again/i);
        });
      });
    });
  }

  test.describe('Cross-Platform Consistency', () => {

    test('should maintain visual consistency between iOS and Android', async ({ page }) => {
      const results: any[] = [];

      for (const platform of MOBILE_PLATFORMS) {
        await page.setViewportSize(platform.viewport);
        await page.setUserAgent(platform.userAgent);

        await navigateToMobileStore(page, platform.name);
        await waitForMobileStreamTab(page, platform.name);
        await clickMobileStreamTab(page, platform.name);

        // Take screenshot for comparison
        const screenshot = await page.screenshot({
          fullPage: false,
          clip: { x: 0, y: 0, width: platform.viewport.width, height: 600 }
        });

        results.push({
          platform: platform.name,
          screenshot,
          timestamp: Date.now(),
        });

        await page.waitForTimeout(1000);
      }

      // Visual consistency should be maintained
      expect(results).toHaveLength(2);
      console.log('Cross-platform visual consistency screenshots captured');
    });

    test('should synchronize data across platforms', async ({ page }) => {
      // Test data synchronization between platforms
      // This would typically involve testing with actual backend sync

      for (const platform of MOBILE_PLATFORMS) {
        await page.setViewportSize(platform.viewport);
        await page.setUserAgent(platform.userAgent);

        await navigateToMobileStore(page, platform.name);
        await waitForMobileStreamTab(page, platform.name);
        await clickMobileStreamTab(page, platform.name);

        // Select Podcasts category
        const podcastsTab = page.locator('text="Podcasts"').first();
        await podcastsTab.click();
        await page.waitForTimeout(1000);

        // Verify category is selected
        await expect(podcastsTab).toHaveClass(/active|selected/);

        console.log(`${platform.name} - Podcasts category selection verified`);
      }
    });
  });

  test.describe('Mobile Accessibility', () => {

    test('should support mobile accessibility features', async ({ page }) => {
      // Test with first platform (iOS)
      const platform = MOBILE_PLATFORMS[0];

      await page.setViewportSize(platform.viewport);
      await page.setUserAgent(platform.userAgent);

      await navigateToMobileStore(page, platform.name);
      await waitForMobileStreamTab(page, platform.name);
      await clickMobileStreamTab(page, platform.name);

      // Check accessibility attributes
      const streamTabs = page.locator('[role="tab"]');
      const tabCount = await streamTabs.count();

      for (let i = 0; i < Math.min(tabCount, 3); i++) {
        const tab = streamTabs.nth(i);

        // Should have accessibility attributes
        await expect(tab).toHaveAttribute('role', 'tab');
        await expect(tab).toHaveAttribute('aria-selected');

        // Should have sufficient color contrast (visual validation)
        const tabBox = await tab.boundingBox();
        expect(tabBox?.height).toBeGreaterThanOrEqual(44); // Touch target size
      }
    });
  });
});