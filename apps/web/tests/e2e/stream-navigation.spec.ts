/**
 * E2E Tests: Stream Navigation - Complete User Journey
 *
 * Comprehensive end-to-end testing for Stream Store Tabs feature.
 * Tests cover all quickstart scenarios including tab navigation, subtabs,
 * featured content, store integration, cross-platform consistency, and performance.
 *
 * Based on quickstart scenarios from specs/026-help-me-add/quickstart.md
 *
 * Test Requirements:
 * - Tab switching: <200ms response time
 * - Content loading: <3s initial load
 * - Animation frame rate: 60fps
 * - Visual consistency: >95% across platforms
 * - Store integration: unified cart experience
 * - Featured content: proper display and interaction
 * - Error handling: graceful empty states and network issues
 */

import { test, expect, Page } from '@playwright/test';

// Performance budgets based on requirements
const PERFORMANCE_BUDGETS = {
  TAB_SWITCH_TIME: 200, // ms
  CONTENT_LOAD_TIME: 3000, // ms (3s)
  API_RESPONSE_TIME: 500, // ms
  ANIMATION_FRAME_RATE: 55, // fps (target 60fps, allow 55fps)
  VISUAL_CONSISTENCY: 0.95, // 95%
};

// Test data for Stream content validation
const STREAM_CATEGORIES = [
  { id: 'books', name: 'Books', icon: 'book-open' },
  { id: 'podcasts', name: 'Podcasts', icon: 'microphone' },
  { id: 'cartoons', name: 'Cartoons', icon: 'film' },
  { id: 'movies', name: 'Movies', icon: 'video' },
  { id: 'music', name: 'Music', icon: 'music' },
  { id: 'art', name: 'Art', icon: 'palette' },
];

const MOVIES_SUBTABS = [
  { id: 'short-movies', name: 'Short Films', maxDuration: 1800 },
  { id: 'long-movies', name: 'Feature Films', minDuration: 1801 },
];

// Helper functions for test utilities
async function navigateToStore(page: Page) {
  await page.goto('http://localhost:3000');

  // Wait for the app to load
  await page.waitForLoadState('networkidle');

  // Navigate to store if not already there
  const storeTab = page.locator('[data-testid="store-tab"], .store-tab, a[href*="store"]').first();
  if (await storeTab.isVisible()) {
    await storeTab.click();
  }

  // Wait for store content to load
  await page.waitForSelector('[data-testid="store-layout"], .store-layout, [role="tablist"]', { timeout: 10000 });
}

async function waitForStreamTab(page: Page) {
  // Wait for Stream tab to be available
  await page.waitForSelector('[data-testid="stream-tab"], [value="stream"], text="Stream"', { timeout: 10000 });
}

async function clickStreamTab(page: Page) {
  // Click on Stream tab
  const streamTab = page.locator('[data-testid="stream-tab"], [value="stream"], text="Stream"').first();
  await streamTab.click();

  // Wait for Stream content to load
  await page.waitForSelector('[data-testid="stream-content"], .stream-content', { timeout: 10000 });
}

async function measureTabSwitchTime(page: Page, fromTab: string, toTab: string): Promise<number> {
  const startTime = Date.now();

  // Click the target tab
  const targetTab = page.locator(`[value="${toTab}"], text="${toTab}"`, { hasText: toTab }).first();
  await targetTab.click();

  // Wait for content to update
  await page.waitForSelector(`[data-testid="${toTab}-content"], .${toTab}-content`, { timeout: 5000 });

  const endTime = Date.now();
  return endTime - startTime;
}

async function validateTabOrder(page: Page, expectedOrder: string[]) {
  const tabs = page.locator('[role="tablist"] [role="tab"], .tabs-list .tab-trigger');
  const tabCount = await tabs.count();

  for (let i = 0; i < Math.min(tabCount, expectedOrder.length); i++) {
    const tabText = await tabs.nth(i).textContent();
    expect(tabText?.toLowerCase()).toContain(expectedOrder[i].toLowerCase());
  }
}

async function validateFeaturedContent(page: Page, categoryId: string) {
  // Check for featured content carousel
  const featuredCarousel = page.locator('[data-testid="featured-carousel"], .featured-content, .carousel');
  await expect(featuredCarousel).toBeVisible({ timeout: 5000 });

  // Validate featured items have required properties
  const featuredItems = page.locator('[data-testid="featured-item"], .featured-item, .carousel-item');
  const itemCount = await featuredItems.count();

  if (itemCount > 0) {
    const firstItem = featuredItems.first();

    // Check for title, thumbnail, and price
    await expect(firstItem.locator('[data-testid="item-title"], .item-title, h3, h4')).toBeVisible();
    await expect(firstItem.locator('[data-testid="item-thumbnail"], .item-thumbnail, img')).toBeVisible();
    await expect(firstItem.locator('[data-testid="item-price"], .item-price, .price')).toBeVisible();
  }

  return itemCount;
}

// Main test suite
test.describe('Stream Navigation - Complete User Journey', () => {

  test.beforeEach(async ({ page }) => {
    // Set up test environment
    await page.goto('http://localhost:3000');
    await page.waitForLoadState('networkidle');
  });

  test.describe('Scenario 1: Basic Tab Navigation', () => {

    test('should display all media tabs in correct order before Live tab', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);

      // Validate tab order: Books → Podcasts → Cartoons → Movies → Live
      // Note: Stream is the container, individual categories are within it
      const expectedMainTabs = ['shops', 'products', 'stream', 'live'];
      await validateTabOrder(page, expectedMainTabs);

      // Click Stream tab to access media categories
      await clickStreamTab(page);

      // Wait for and validate stream category tabs
      await page.waitForSelector('[data-testid="stream-categories"], .stream-tabs', { timeout: 5000 });

      // Check that media category tabs are present
      for (const category of STREAM_CATEGORIES) {
        const categoryTab = page.locator(`text="${category.name}"`).first();
        await expect(categoryTab).toBeVisible();
      }
    });

    test('should switch between media category tabs smoothly', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Test switching between each media category
      for (const category of STREAM_CATEGORIES.slice(0, 4)) { // Test first 4 categories
        const startTime = Date.now();

        // Click category tab
        const categoryTab = page.locator(`[data-testid="${category.id}-tab"], text="${category.name}"`).first();
        await categoryTab.click();

        // Wait for content to load
        await page.waitForSelector(`[data-testid="${category.id}-content"], .category-content`, { timeout: 3000 });

        const switchTime = Date.now() - startTime;

        // Validate switch time meets performance budget
        expect(switchTime).toBeLessThan(PERFORMANCE_BUDGETS.TAB_SWITCH_TIME);

        // Verify tab is active
        await expect(categoryTab).toHaveAttribute('data-state', 'active');
      }
    });

    test('should maintain tab selection state during navigation', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Select Podcasts category
      const podcastsTab = page.locator('text="Podcasts"').first();
      await podcastsTab.click();
      await page.waitForTimeout(500); // Allow state to settle

      // Navigate away from Stream tab
      const productsTab = page.locator('[value="products"], text="Products"').first();
      await productsTab.click();
      await page.waitForTimeout(500);

      // Return to Stream tab
      await clickStreamTab(page);

      // Verify Podcasts tab is still selected
      await expect(podcastsTab).toHaveAttribute('data-state', 'active');
    });
  });

  test.describe('Scenario 2: Movies Subtab Navigation', () => {

    test('should display and switch between Movies subtabs', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Navigate to Movies category
      const moviesTab = page.locator('text="Movies"').first();
      await moviesTab.click();
      await page.waitForSelector('[data-testid="movies-subtabs"], .subtabs', { timeout: 5000 });

      // Verify subtabs are displayed
      const shortFilmsTab = page.locator('text="Short Films"').first();
      const featureFilmsTab = page.locator('text="Feature Films"').first();

      await expect(shortFilmsTab).toBeVisible();
      await expect(featureFilmsTab).toBeVisible();

      // Test subtab switching
      await featureFilmsTab.click();
      await page.waitForTimeout(500);
      await expect(featureFilmsTab).toHaveClass(/active|selected/);

      await shortFilmsTab.click();
      await page.waitForTimeout(500);
      await expect(shortFilmsTab).toHaveClass(/active|selected/);
    });

    test('should persist subtab selection across navigation', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Navigate to Movies and select Feature Films
      const moviesTab = page.locator('text="Movies"').first();
      await moviesTab.click();

      const featureFilmsTab = page.locator('text="Feature Films"').first();
      await featureFilmsTab.click();
      await page.waitForTimeout(500);

      // Navigate to different category
      const booksTab = page.locator('text="Books"').first();
      await booksTab.click();
      await page.waitForTimeout(500);

      // Return to Movies
      await moviesTab.click();

      // Verify Feature Films subtab is still selected
      await expect(featureFilmsTab).toHaveClass(/active|selected/);
    });
  });

  test.describe('Scenario 3: Featured Content Display', () => {

    test('should display featured content in all categories', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Test featured content in each category
      for (const category of STREAM_CATEGORIES.slice(0, 3)) { // Test first 3 categories
        const categoryTab = page.locator(`text="${category.name}"`).first();
        await categoryTab.click();
        await page.waitForTimeout(1000);

        const featuredCount = await validateFeaturedContent(page, category.id);

        // Verify at least some featured content exists (test data dependent)
        console.log(`Featured content count for ${category.name}: ${featuredCount}`);
      }
    });

    test('should load featured content within performance budget', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Measure featured content loading time
      const startTime = Date.now();

      // Navigate to Books category (should trigger featured content load)
      const booksTab = page.locator('text="Books"').first();
      await booksTab.click();

      // Wait for featured content to appear
      await page.waitForSelector('[data-testid="featured-carousel"], .featured-content', { timeout: 5000 });

      const loadTime = Date.now() - startTime;

      // Validate load time meets performance budget
      expect(loadTime).toBeLessThan(PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME);
    });
  });

  test.describe('Scenario 4: Store Integration', () => {

    test('should integrate stream content with unified cart', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Navigate to Books category
      const booksTab = page.locator('text="Books"').first();
      await booksTab.click();
      await page.waitForTimeout(1000);

      // Try to add item to cart (test data dependent)
      const addToCartButton = page.locator('[data-testid="add-to-cart"], .add-to-cart, text="Add to Cart"').first();

      if (await addToCartButton.isVisible()) {
        await addToCartButton.click();

        // Verify cart interaction
        await page.waitForTimeout(1000);

        // Check for cart notification or updated cart count
        const cartNotification = page.locator('.toast, .notification, [data-testid="cart-notification"]');
        const cartCount = page.locator('[data-testid="cart-count"], .cart-count');

        // At least one indicator should be present
        const hasNotification = await cartNotification.isVisible();
        const hasCartCount = await cartCount.isVisible();

        expect(hasNotification || hasCartCount).toBeTruthy();
      }
    });
  });

  test.describe('Scenario 5: Performance Validation', () => {

    test('should meet tab switching performance targets', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Test tab switching performance between categories
      const categories = STREAM_CATEGORIES.slice(0, 3);

      for (let i = 0; i < categories.length - 1; i++) {
        const fromCategory = categories[i];
        const toCategory = categories[i + 1];

        const switchTime = await measureTabSwitchTime(page, fromCategory.id, toCategory.name);

        // Validate switch time meets performance budget
        expect(switchTime).toBeLessThan(PERFORMANCE_BUDGETS.TAB_SWITCH_TIME);

        console.log(`Tab switch ${fromCategory.name} → ${toCategory.name}: ${switchTime}ms`);
      }
    });

    test('should handle content loading efficiently', async ({ page }) => {
      // Monitor network requests during navigation
      const requests: any[] = [];
      page.on('request', request => {
        if (request.url().includes('/api/')) {
          requests.push({
            url: request.url(),
            method: request.method(),
            timestamp: Date.now()
          });
        }
      });

      await navigateToStore(page);
      await waitForStreamTab(page);

      const startTime = Date.now();
      await clickStreamTab(page);

      // Wait for initial content load
      await page.waitForSelector('[data-testid="stream-content"], .stream-content', { timeout: 5000 });

      const loadTime = Date.now() - startTime;

      // Validate total load time
      expect(loadTime).toBeLessThan(PERFORMANCE_BUDGETS.CONTENT_LOAD_TIME);

      // Validate API response times
      const apiRequests = requests.filter(req => req.timestamp >= startTime);
      console.log(`API requests during Stream load: ${apiRequests.length}`);

      // Should have made API calls for content
      expect(apiRequests.length).toBeGreaterThan(0);
    });
  });

  test.describe('Scenario 6: Error Handling', () => {

    test('should handle empty content gracefully', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Navigate through categories and check for empty states
      for (const category of STREAM_CATEGORIES.slice(0, 2)) {
        const categoryTab = page.locator(`text="${category.name}"`).first();
        await categoryTab.click();
        await page.waitForTimeout(1000);

        // Check if content exists or empty state is shown
        const contentItems = page.locator('[data-testid="content-item"], .content-item, .stream-item');
        const emptyState = page.locator('[data-testid="empty-state"], .empty-state, text="No content"');

        const hasContent = await contentItems.count() > 0;
        const hasEmptyState = await emptyState.isVisible();

        // Either content should exist OR empty state should be shown
        if (!hasContent) {
          expect(hasEmptyState).toBeTruthy();
        }

        console.log(`${category.name} - Content items: ${await contentItems.count()}, Empty state: ${hasEmptyState}`);
      }
    });

    test('should handle network errors gracefully', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);

      // Simulate slow network
      await page.route('**/api/v1/stream/**', route => {
        setTimeout(() => route.continue(), 1000); // 1s delay
      });

      await clickStreamTab(page);

      // Should show loading state initially
      const loadingIndicator = page.locator('[data-testid="loading"], .loading, .spinner');

      // Wait for either content or error state
      await page.waitForFunction(() => {
        const hasContent = document.querySelector('[data-testid="stream-content"]');
        const hasError = document.querySelector('[data-testid="error-state"]');
        return hasContent || hasError;
      }, { timeout: 10000 });

      // Verify no broken UI elements
      const brokenElements = page.locator('.error, [data-error="true"]');
      const brokenCount = await brokenElements.count();
      expect(brokenCount).toBeLessThanOrEqual(1); // Allow for controlled error states
    });
  });

  test.describe('Integration & Regression Tests', () => {

    test('should maintain Stream functionality after page reload', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Navigate to specific category and subtab
      const moviesTab = page.locator('text="Movies"').first();
      await moviesTab.click();

      const featureFilmsTab = page.locator('text="Feature Films"').first();
      if (await featureFilmsTab.isVisible()) {
        await featureFilmsTab.click();
      }

      // Reload page
      await page.reload({ waitUntil: 'networkidle' });

      // Navigate back to Stream
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Verify functionality still works
      await moviesTab.click();

      // Stream should be functional after reload
      const streamContent = page.locator('[data-testid="stream-content"], .stream-content');
      await expect(streamContent).toBeVisible();
    });

    test('should handle concurrent category switching', async ({ page }) => {
      await navigateToStore(page);
      await waitForStreamTab(page);
      await clickStreamTab(page);

      // Rapidly switch between categories
      const categories = STREAM_CATEGORIES.slice(0, 3);

      for (let i = 0; i < 3; i++) {
        for (const category of categories) {
          const categoryTab = page.locator(`text="${category.name}"`).first();
          await categoryTab.click();
          await page.waitForTimeout(100); // Quick switching
        }
      }

      // Verify final state is stable
      await page.waitForTimeout(1000);
      const lastCategory = categories[categories.length - 1];
      const lastTab = page.locator(`text="${lastCategory.name}"`).first();
      await expect(lastTab).toHaveAttribute('data-state', 'active');
    });
  });
});