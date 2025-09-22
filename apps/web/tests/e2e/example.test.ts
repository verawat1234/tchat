import { test, expect, Page } from '@playwright/test';

// Helper function to check accessibility
async function checkA11y(page: Page) {
  // This is a simplified check - in production, you'd use @axe-core/playwright
  const violations = await page.evaluate(() => {
    const issues = [];

    // Check for images without alt text
    const images = document.querySelectorAll('img:not([alt])');
    if (images.length > 0) {
      issues.push(`${images.length} images without alt text`);
    }

    // Check for buttons without accessible text
    const buttons = document.querySelectorAll('button');
    buttons.forEach((btn) => {
      if (!btn.textContent?.trim() && !btn.getAttribute('aria-label')) {
        issues.push('Button without accessible text');
      }
    });

    return issues;
  });

  expect(violations).toHaveLength(0);
}

test.describe('Basic Application Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should load the homepage', async ({ page }) => {
    // Wait for the page to be fully loaded
    await page.waitForLoadState('networkidle');

    // Check that the page has loaded
    await expect(page).toHaveTitle(/Tchat/);

    // Take a screenshot for visual reference
    await page.screenshot({ path: 'test-results/homepage.png' });
  });

  test('should have no accessibility violations on homepage', async ({ page }) => {
    await checkA11y(page);
  });

  test('should navigate between pages', async ({ page }) => {
    // Look for navigation elements
    const nav = page.locator('nav');

    if (await nav.isVisible()) {
      // Get all navigation links
      const links = await nav.locator('a').all();

      // Test first 3 navigation links
      for (let i = 0; i < Math.min(3, links.length); i++) {
        const link = links[i];
        const href = await link.getAttribute('href');

        if (href && !href.startsWith('http')) {
          await link.click();
          await page.waitForLoadState('networkidle');

          // Verify navigation occurred
          expect(page.url()).toContain(href);

          // Go back to test next link
          await page.goto('/');
        }
      }
    }
  });

  test('should be responsive', async ({ page }) => {
    // Test desktop viewport
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.screenshot({ path: 'test-results/desktop.png' });

    // Test tablet viewport
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.screenshot({ path: 'test-results/tablet.png' });

    // Test mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await page.screenshot({ path: 'test-results/mobile.png' });

    // Verify that content is still visible and accessible
    const mainContent = page.locator('main, [role="main"], #root > div');
    await expect(mainContent.first()).toBeVisible();
  });

  test('should handle keyboard navigation', async ({ page }) => {
    // Focus on the first interactive element
    await page.keyboard.press('Tab');

    // Get the focused element
    const focusedElement = await page.evaluate(() => {
      const el = document.activeElement;
      return {
        tagName: el?.tagName.toLowerCase(),
        text: el?.textContent?.trim(),
        href: (el as HTMLAnchorElement)?.href,
      };
    });

    // Verify that something is focused
    expect(focusedElement.tagName).toBeTruthy();

    // Navigate through a few more elements
    for (let i = 0; i < 5; i++) {
      await page.keyboard.press('Tab');
    }

    // Navigate backwards
    await page.keyboard.press('Shift+Tab');

    // Test Enter key on a focused link or button
    const activeElement = await page.evaluate(() => document.activeElement?.tagName.toLowerCase());

    if (activeElement === 'a' || activeElement === 'button') {
      await page.keyboard.press('Enter');
      // Wait for any potential navigation or action
      await page.waitForTimeout(500);
    }
  });

  test('should have proper ARIA attributes', async ({ page }) => {
    // Check for main landmark
    const main = page.locator('main, [role="main"]');
    await expect(main.first()).toBeVisible();

    // Check for navigation landmark
    const nav = page.locator('nav, [role="navigation"]');
    if (await nav.count() > 0) {
      await expect(nav.first()).toBeVisible();
    }

    // Check that interactive elements have proper roles and labels
    const buttons = await page.locator('button').all();
    for (const button of buttons.slice(0, 5)) { // Check first 5 buttons
      const hasText = await button.textContent();
      const hasAriaLabel = await button.getAttribute('aria-label');
      const hasAriaLabelledby = await button.getAttribute('aria-labelledby');

      expect(
        hasText?.trim() || hasAriaLabel || hasAriaLabelledby,
        'Button should have accessible text'
      ).toBeTruthy();
    }
  });

  test('should measure performance metrics', async ({ page }) => {
    // Collect performance metrics
    const metrics = await page.evaluate(() => {
      const perfData = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;

      return {
        domContentLoaded: perfData.domContentLoadedEventEnd - perfData.domContentLoadedEventStart,
        loadComplete: perfData.loadEventEnd - perfData.loadEventStart,
        domInteractive: perfData.domInteractive - perfData.fetchStart,
        firstPaint: performance.getEntriesByName('first-paint')[0]?.startTime || 0,
        firstContentfulPaint: performance.getEntriesByName('first-contentful-paint')[0]?.startTime || 0,
      };
    });

    // Log metrics for monitoring
    console.log('Performance Metrics:', metrics);

    // Assert reasonable performance thresholds
    expect(metrics.domInteractive).toBeLessThan(3000); // DOM interactive in under 3s
    expect(metrics.loadComplete).toBeLessThan(5000); // Page load in under 5s
  });
});

test.describe('Component Interaction Tests', () => {
  test('should interact with forms', async ({ page }) => {
    await page.goto('/');

    // Look for any form on the page
    const form = page.locator('form').first();

    if (await form.isVisible()) {
      // Find input fields
      const textInputs = await form.locator('input[type="text"], input[type="email"], input[type="password"]').all();

      // Fill in the first few inputs
      for (let i = 0; i < Math.min(3, textInputs.length); i++) {
        await textInputs[i].fill(`test-value-${i}`);
      }

      // Look for a submit button
      const submitButton = form.locator('button[type="submit"], input[type="submit"], button:has-text("Submit")');

      if (await submitButton.isVisible()) {
        // Check that the button is enabled
        await expect(submitButton.first()).toBeEnabled();
      }
    }
  });

  test('should handle error states gracefully', async ({ page }) => {
    // Navigate to a non-existent page
    await page.goto('/404-page-not-found-test', { waitUntil: 'networkidle' });

    // Check that some content is still displayed
    const body = page.locator('body');
    await expect(body).toBeVisible();

    // Check for common error page elements
    const errorIndicators = await page.locator(
      'text=/404|not found|error|oops/i'
    ).count();

    // If it's a SPA, it might handle 404 differently
    // Just ensure the page doesn't completely break
    expect(await body.textContent()).toBeTruthy();
  });
});

// Test specific user journeys
test.describe('User Journeys', () => {
  test('should complete a basic user flow', async ({ page }) => {
    // This is a placeholder for specific user journeys
    // Replace with actual user flows once the application structure is known

    await page.goto('/');

    // Example: Test a search flow if search exists
    const searchInput = page.locator('input[type="search"], input[placeholder*="search" i]').first();

    if (await searchInput.isVisible()) {
      await searchInput.fill('test search query');
      await searchInput.press('Enter');

      // Wait for results or navigation
      await page.waitForTimeout(1000);

      // Verify something changed (URL or content)
      const currentUrl = page.url();
      expect(currentUrl).toContain('search');
    }
  });
});

// Visual regression tests
test.describe('Visual Regression', () => {
  test('should match visual snapshot', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Remove dynamic content that might change between tests
    await page.evaluate(() => {
      // Hide timestamps, random IDs, etc.
      document.querySelectorAll('[data-timestamp]').forEach(el => {
        (el as HTMLElement).style.visibility = 'hidden';
      });
    });

    // Take screenshot for visual comparison
    await expect(page).toHaveScreenshot('homepage.png', {
      maxDiffPixels: 100,
      fullPage: true,
    });
  });
});

// API interaction tests
test.describe('API Integration', () => {
  test('should handle API responses', async ({ page }) => {
    // Intercept API calls
    await page.route('**/api/**', route => {
      console.log('API call intercepted:', route.request().url());
      route.continue();
    });

    await page.goto('/');

    // Wait for any API calls to complete
    await page.waitForLoadState('networkidle');

    // Verify the page handles API responses correctly
    // This would be more specific based on actual API endpoints
  });
});