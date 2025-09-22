/**
 * T056: E2E Test - Dynamic Content Loading Flow
 *
 * Comprehensive E2E tests using Playwright for the dynamic content management system.
 * Tests complete user journeys from page load to content display with real browser automation.
 *
 * Test Coverage:
 * - Complete content loading workflow from user perspective
 * - Loading states: spinners, skeleton loaders, progress indicators
 * - Content rendering: dynamic content replacing hardcoded text
 * - Error handling: network failures, content loading errors, recovery
 * - Performance validation: <200ms content load time requirements
 * - Multiple content types: text, rich text, images, config
 * - User interactions: navigation, refresh, content triggers
 * - Cross-browser compatibility and responsive behavior
 * - Accessibility compliance and keyboard navigation
 */

import { test, expect, Page, BrowserContext } from '@playwright/test';

// Test data constants
const PERFORMANCE_THRESHOLD = 200; // <200ms requirement
const CONTENT_LOAD_TIMEOUT = 5000;
const NAVIGATION_TIMEOUT = 3000;

// Mock content data for testing
const mockContentItems = {
  welcomeMessage: {
    id: 'welcome.message',
    type: 'text',
    value: 'Welcome to Tchat - your dynamic chat platform!'
  },
  headerTitle: {
    id: 'header.title',
    type: 'rich_text',
    value: '<h1>Dynamic Header Title</h1>'
  },
  profileImage: {
    id: 'profile.image',
    type: 'image_url',
    value: 'https://example.com/profile.jpg'
  },
  featureConfig: {
    id: 'features.notifications',
    type: 'config',
    value: { enabled: true, sound: false }
  }
};

// Helper functions for common operations
class ContentLoadingHelper {
  constructor(private page: Page) {}

  /**
   * Wait for content loading state to appear and disappear
   */
  async waitForContentLoading(testId: string, timeout = CONTENT_LOAD_TIMEOUT) {
    // Wait for loading state to appear
    await this.page.waitForSelector(`[data-testid="${testId}-loading"]`, {
      state: 'visible',
      timeout
    });

    // Wait for loading state to disappear (content loaded)
    await this.page.waitForSelector(`[data-testid="${testId}-loading"]`, {
      state: 'hidden',
      timeout
    });
  }

  /**
   * Measure content load performance
   */
  async measureContentLoadTime(testId: string): Promise<number> {
    const startTime = Date.now();
    await this.waitForContentLoading(testId);
    const endTime = Date.now();
    return endTime - startTime;
  }

  /**
   * Check if skeleton loader is displayed correctly
   */
  async checkSkeletonLoader(testId: string) {
    const skeleton = this.page.locator(`[data-testid="${testId}-skeleton"]`);
    await expect(skeleton).toBeVisible();
    await expect(skeleton).toHaveClass(/skeleton|loading/);
  }

  /**
   * Verify content has been populated
   */
  async verifyContentPopulated(testId: string, expectedContent?: string) {
    const contentElement = this.page.locator(`[data-testid="${testId}-content"]`);
    await expect(contentElement).toBeVisible();

    if (expectedContent) {
      await expect(contentElement).toContainText(expectedContent);
    }

    // Ensure it's not showing loading or error state
    await expect(this.page.locator(`[data-testid="${testId}-loading"]`)).not.toBeVisible();
    await expect(this.page.locator(`[data-testid="${testId}-error"]`)).not.toBeVisible();
  }

  /**
   * Trigger content refresh and verify reload
   */
  async triggerContentRefresh(testId: string) {
    const refreshButton = this.page.locator(`[data-testid="${testId}-refresh"]`);
    await refreshButton.click();
    await this.waitForContentLoading(testId);
  }

  /**
   * Check error state display
   */
  async verifyErrorState(testId: string, expectedErrorMessage?: string) {
    const errorElement = this.page.locator(`[data-testid="${testId}-error"]`);
    await expect(errorElement).toBeVisible();
    await expect(errorElement).toHaveAttribute('role', 'alert');

    if (expectedErrorMessage) {
      await expect(errorElement).toContainText(expectedErrorMessage);
    }
  }
}

// Setup mock API responses
async function setupMockAPIResponses(context: BrowserContext) {
  // Mock successful content API responses
  await context.route('**/api/content/**', async (route) => {
    const url = route.request().url();
    const contentId = url.split('/').pop();

    if (mockContentItems[contentId as keyof typeof mockContentItems]) {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          success: true,
          data: mockContentItems[contentId as keyof typeof mockContentItems]
        })
      });
    } else {
      await route.fulfill({
        status: 404,
        contentType: 'application/json',
        body: JSON.stringify({
          success: false,
          error: 'Content not found'
        })
      });
    }
  });

  // Mock content list endpoint
  await context.route('**/api/content', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data: {
          items: Object.values(mockContentItems),
          pagination: {
            page: 1,
            limit: 20,
            total: Object.keys(mockContentItems).length,
            totalPages: 1
          }
        }
      })
    });
  });
}

// Setup slow network simulation
async function setupSlowNetwork(context: BrowserContext) {
  await context.route('**/api/content/**', async (route) => {
    // Simulate 1 second delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    await route.continue();
  });
}

// Setup network failure simulation
async function setupNetworkFailure(context: BrowserContext) {
  await context.route('**/api/content/**', async (route) => {
    await route.fulfill({
      status: 500,
      contentType: 'application/json',
      body: JSON.stringify({
        success: false,
        error: 'Internal server error'
      })
    });
  });
}

test.describe('Dynamic Content Loading Flow', () => {
  let helper: ContentLoadingHelper;

  test.beforeEach(async ({ page, context }) => {
    helper = new ContentLoadingHelper(page);
    await setupMockAPIResponses(context);
    await page.goto('/');
  });

  test.describe('Content Loading Workflow', () => {
    test('should load page and display initial loading states', async ({ page }) => {
      // Verify page loads with loading indicators
      await expect(page.locator('[data-testid="main-content-loading"]')).toBeVisible();

      // Check for skeleton loaders
      await helper.checkSkeletonLoader('welcome-message');
      await helper.checkSkeletonLoader('header-title');

      // Verify loading indicators have proper accessibility attributes
      const loadingElement = page.locator('[data-testid="main-content-loading"]');
      await expect(loadingElement).toHaveAttribute('role', 'status');
      await expect(loadingElement).toHaveAttribute('aria-label', /loading/i);
    });

    test('should complete content loading workflow within performance threshold', async ({ page }) => {
      // Measure welcome message load time
      const welcomeLoadTime = await helper.measureContentLoadTime('welcome-message');
      expect(welcomeLoadTime).toBeLessThan(PERFORMANCE_THRESHOLD);

      // Verify content is populated
      await helper.verifyContentPopulated('welcome-message', 'Welcome to Tchat');

      // Measure header title load time
      const headerLoadTime = await helper.measureContentLoadTime('header-title');
      expect(headerLoadTime).toBeLessThan(PERFORMANCE_THRESHOLD);

      // Verify dynamic header is displayed
      await helper.verifyContentPopulated('header-title');
      await expect(page.locator('[data-testid="header-title-content"] h1')).toContainText('Dynamic Header Title');
    });

    test('should replace hardcoded text with dynamic content', async ({ page }) => {
      // Initially should show loading or placeholder
      const contentElement = page.locator('[data-testid="welcome-message-content"]');

      // Wait for content to load
      await helper.waitForContentLoading('welcome-message');

      // Verify hardcoded text is replaced with dynamic content
      await expect(contentElement).toContainText('Welcome to Tchat - your dynamic chat platform!');
      await expect(contentElement).not.toContainText('Loading...');
      await expect(contentElement).not.toContainText('Placeholder text');
    });
  });

  test.describe('Loading States and UI Feedback', () => {
    test('should display loading spinners during content fetch', async ({ page }) => {
      // Check loading spinner is visible initially
      const spinner = page.locator('[data-testid="content-loading-spinner"]');
      await expect(spinner).toBeVisible();

      // Verify spinner has proper animation
      await expect(spinner).toHaveClass(/animate-spin|spinning/);

      // Wait for content to load and spinner to disappear
      await helper.waitForContentLoading('welcome-message');
      await expect(spinner).not.toBeVisible();
    });

    test('should show skeleton loaders for different content types', async ({ page }) => {
      // Text content skeleton
      await helper.checkSkeletonLoader('welcome-message');

      // Rich text content skeleton
      await helper.checkSkeletonLoader('header-title');

      // Image content skeleton
      await helper.checkSkeletonLoader('profile-image');

      // Config content skeleton
      await helper.checkSkeletonLoader('feature-config');

      // Wait for all skeletons to be replaced with content
      await Promise.all([
        helper.waitForContentLoading('welcome-message'),
        helper.waitForContentLoading('header-title'),
        helper.waitForContentLoading('profile-image'),
        helper.waitForContentLoading('feature-config')
      ]);
    });

    test('should display progress indicators for large content loads', async ({ page, context }) => {
      // Setup slower loading for progress testing
      await setupSlowNetwork(context);

      // Navigate to page with progress indicators
      await page.goto('/content-heavy');

      // Check progress bar appears
      const progressBar = page.locator('[data-testid="content-progress-bar"]');
      await expect(progressBar).toBeVisible();

      // Verify progress updates
      await expect(progressBar).toHaveAttribute('aria-valuenow');

      // Wait for completion
      await helper.waitForContentLoading('content-heavy');
      await expect(progressBar).not.toBeVisible();
    });
  });

  test.describe('Multiple Content Types', () => {
    test('should handle text content loading', async ({ page }) => {
      await helper.waitForContentLoading('welcome-message');

      const textContent = page.locator('[data-testid="welcome-message-content"]');
      await expect(textContent).toContainText('Welcome to Tchat');
      await expect(textContent).toHaveAttribute('data-content-type', 'text');
    });

    test('should handle rich text content loading', async ({ page }) => {
      await helper.waitForContentLoading('header-title');

      const richTextContent = page.locator('[data-testid="header-title-content"]');
      await expect(richTextContent).toContainText('Dynamic Header Title');
      await expect(richTextContent.locator('h1')).toBeVisible();
      await expect(richTextContent).toHaveAttribute('data-content-type', 'rich_text');
    });

    test('should handle image content loading', async ({ page }) => {
      await helper.waitForContentLoading('profile-image');

      const imageContent = page.locator('[data-testid="profile-image-content"] img');
      await expect(imageContent).toBeVisible();
      await expect(imageContent).toHaveAttribute('src', 'https://example.com/profile.jpg');
      await expect(imageContent).toHaveAttribute('alt');
    });

    test('should handle configuration content loading', async ({ page }) => {
      await helper.waitForContentLoading('feature-config');

      const configContent = page.locator('[data-testid="feature-config-content"]');
      await expect(configContent).toBeVisible();
      await expect(configContent).toHaveAttribute('data-config-enabled', 'true');
      await expect(configContent).toHaveAttribute('data-config-sound', 'false');
    });
  });

  test.describe('Error Handling and Recovery', () => {
    test('should handle network failures gracefully', async ({ page, context }) => {
      // Setup network failure
      await setupNetworkFailure(context);

      // Navigate to trigger content loading
      await page.goto('/');

      // Verify error state is displayed
      await helper.verifyErrorState('welcome-message', 'Failed to load content');

      // Check retry button is available
      const retryButton = page.locator('[data-testid="welcome-message-retry"]');
      await expect(retryButton).toBeVisible();
      await expect(retryButton).toContainText(/retry|reload/i);
    });

    test('should provide content recovery mechanisms', async ({ page, context }) => {
      // Initial network failure
      await setupNetworkFailure(context);
      await page.goto('/');

      // Verify error state
      await helper.verifyErrorState('welcome-message');

      // Fix network and retry
      await setupMockAPIResponses(context);
      await helper.triggerContentRefresh('welcome-message');

      // Verify content loads successfully after retry
      await helper.verifyContentPopulated('welcome-message');
    });

    test('should handle partial content loading failures', async ({ page, context }) => {
      // Mock mixed success/failure responses
      await context.route('**/api/content/welcomeMessage', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: mockContentItems.welcomeMessage
          })
        });
      });

      await context.route('**/api/content/headerTitle', async (route) => {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({
            success: false,
            error: 'Server error'
          })
        });
      });

      await page.goto('/');

      // Verify partial success
      await helper.verifyContentPopulated('welcome-message');
      await helper.verifyErrorState('header-title');
    });

    test('should display fallback content when available', async ({ page, context }) => {
      // Setup network failure
      await setupNetworkFailure(context);

      // Navigate to page with fallback content
      await page.goto('/content-with-fallback');

      // Verify fallback content is displayed
      const fallbackContent = page.locator('[data-testid="welcome-message-fallback"]');
      await expect(fallbackContent).toBeVisible();
      await expect(fallbackContent).toContainText('Default welcome message');

      // Verify fallback indicator
      const fallbackIndicator = page.locator('[data-testid="content-fallback-indicator"]');
      await expect(fallbackIndicator).toBeVisible();
    });
  });

  test.describe('User Interactions and Navigation', () => {
    test('should trigger content loading on navigation', async ({ page }) => {
      // Navigate to different sections
      await page.click('[data-testid="nav-chat-tab"]');
      await helper.waitForContentLoading('chat-content');
      await helper.verifyContentPopulated('chat-content');

      await page.click('[data-testid="nav-social-tab"]');
      await helper.waitForContentLoading('social-content');
      await helper.verifyContentPopulated('social-content');
    });

    test('should refresh content on user action', async ({ page }) => {
      // Wait for initial load
      await helper.waitForContentLoading('welcome-message');

      // Trigger manual refresh
      await page.click('[data-testid="content-refresh-button"]');

      // Verify content reloads
      await helper.waitForContentLoading('welcome-message');
      await helper.verifyContentPopulated('welcome-message');
    });

    test('should handle page refresh and maintain content state', async ({ page }) => {
      // Wait for initial content load
      await helper.waitForContentLoading('welcome-message');
      await helper.verifyContentPopulated('welcome-message');

      // Refresh the page
      await page.reload();

      // Verify content loads again
      await helper.waitForContentLoading('welcome-message');
      await helper.verifyContentPopulated('welcome-message');
    });

    test('should support keyboard navigation for content interactions', async ({ page }) => {
      // Focus on retry button using keyboard
      await page.keyboard.press('Tab');
      const retryButton = page.locator('[data-testid="welcome-message-retry"]');
      await expect(retryButton).toBeFocused();

      // Activate retry with keyboard
      await page.keyboard.press('Enter');
      await helper.waitForContentLoading('welcome-message');
    });
  });

  test.describe('Performance Validation', () => {
    test('should meet content load time requirements', async ({ page }) => {
      const startTime = Date.now();

      // Wait for all critical content to load
      await Promise.all([
        helper.waitForContentLoading('welcome-message'),
        helper.waitForContentLoading('header-title'),
        helper.waitForContentLoading('navigation-menu')
      ]);

      const totalLoadTime = Date.now() - startTime;

      // Verify performance requirement
      expect(totalLoadTime).toBeLessThan(PERFORMANCE_THRESHOLD);
    });

    test('should optimize content loading with caching', async ({ page }) => {
      // First load - measure time
      const firstLoadStart = Date.now();
      await helper.waitForContentLoading('welcome-message');
      const firstLoadTime = Date.now() - firstLoadStart;

      // Navigate away and back
      await page.click('[data-testid="nav-chat-tab"]');
      await page.click('[data-testid="nav-home-tab"]');

      // Second load - should be faster due to caching
      const secondLoadStart = Date.now();
      await helper.waitForContentLoading('welcome-message');
      const secondLoadTime = Date.now() - secondLoadStart;

      // Cached load should be significantly faster
      expect(secondLoadTime).toBeLessThan(firstLoadTime * 0.5);
    });

    test('should handle concurrent content loading efficiently', async ({ page }) => {
      const loadPromises = [
        helper.measureContentLoadTime('welcome-message'),
        helper.measureContentLoadTime('header-title'),
        helper.measureContentLoadTime('profile-image'),
        helper.measureContentLoadTime('feature-config')
      ];

      const loadTimes = await Promise.all(loadPromises);

      // All content should load within threshold
      loadTimes.forEach(time => {
        expect(time).toBeLessThan(PERFORMANCE_THRESHOLD);
      });
    });
  });

  test.describe('Accessibility and Cross-Browser', () => {
    test('should maintain accessibility during content loading', async ({ page }) => {
      // Check loading states have proper ARIA attributes
      const loadingElement = page.locator('[data-testid="content-loading"]');
      await expect(loadingElement).toHaveAttribute('role', 'status');
      await expect(loadingElement).toHaveAttribute('aria-live', 'polite');

      // Check error states have proper ARIA attributes
      await setupNetworkFailure(page.context());
      await page.goto('/');

      const errorElement = page.locator('[data-testid="welcome-message-error"]');
      await expect(errorElement).toHaveAttribute('role', 'alert');
      await expect(errorElement).toHaveAttribute('aria-live', 'assertive');
    });

    test('should announce content loading progress to screen readers', async ({ page }) => {
      // Check for announcement regions
      const announcements = page.locator('[data-testid="content-announcements"]');
      await expect(announcements).toHaveAttribute('aria-live', 'polite');

      // Verify loading announcements
      await expect(announcements).toContainText(/loading/i);

      // Wait for content load and verify completion announcement
      await helper.waitForContentLoading('welcome-message');
      await expect(announcements).toContainText(/loaded|complete/i);
    });

    test('should work across different viewport sizes', async ({ page }) => {
      // Test mobile viewport
      await page.setViewportSize({ width: 375, height: 667 });
      await helper.waitForContentLoading('welcome-message');
      await helper.verifyContentPopulated('welcome-message');

      // Test tablet viewport
      await page.setViewportSize({ width: 768, height: 1024 });
      await helper.waitForContentLoading('welcome-message');
      await helper.verifyContentPopulated('welcome-message');

      // Test desktop viewport
      await page.setViewportSize({ width: 1280, height: 720 });
      await helper.waitForContentLoading('welcome-message');
      await helper.verifyContentPopulated('welcome-message');
    });

    test('should handle reduced motion preferences', async ({ page }) => {
      // Set reduced motion preference
      await page.emulateMedia({ reducedMotion: 'reduce' });

      // Verify animations are disabled/reduced
      const loadingSpinner = page.locator('[data-testid="content-loading-spinner"]');
      await expect(loadingSpinner).not.toHaveClass(/animate-spin/);

      // Content should still load properly
      await helper.waitForContentLoading('welcome-message');
      await helper.verifyContentPopulated('welcome-message');
    });
  });
});

test.describe('Cross-Browser Content Loading', () => {
  ['chromium', 'firefox', 'webkit'].forEach(browserName => {
    test(`should work correctly in ${browserName}`, async ({ page, browserName: currentBrowser }) => {
      // Skip if not the current browser being tested
      test.skip(currentBrowser !== browserName, `Skipping ${browserName} specific test`);

      const helper = new ContentLoadingHelper(page);

      await page.goto('/');
      await helper.waitForContentLoading('welcome-message');
      await helper.verifyContentPopulated('welcome-message');

      // Verify performance across browsers
      const loadTime = await helper.measureContentLoadTime('header-title');
      expect(loadTime).toBeLessThan(PERFORMANCE_THRESHOLD);
    });
  });
});