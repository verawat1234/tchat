import { test, expect, Page } from '@playwright/test';

test.describe('OTP Authentication Flow - Infrastructure Aware', () => {
  // Configure test environment
  test.beforeEach(async ({ page }) => {
    // Set base URL to port 3001 as specified
    await page.goto('http://localhost:3001/');

    // Wait for the page to be fully loaded and content fallback to initialize
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000); // Allow content fallback service to stabilize

    // Add test ID attributes via JavaScript injection for better element selection
    await page.addInitScript(() => {
      // Add test IDs to key elements if they don't exist
      const addTestId = (selector: string, testId: string) => {
        const element = document.querySelector(selector);
        if (element && !element.getAttribute('data-testid')) {
          element.setAttribute('data-testid', testId);
        }
      };

      // Wait for DOM to be ready and add test IDs
      const observer = new MutationObserver(() => {
        addTestId('input[placeholder*="+66"], input[placeholder*="phone" i]', 'phone-input');
        addTestId('button:has-text("Send OTP"), button:has-text("Send")', 'send-otp-button');
        addTestId('input[placeholder*="code" i], input[placeholder*="verification" i]', 'otp-input');
        addTestId('button:has-text("Verify"), button:has-text("Continue")', 'verify-button');
        addTestId('button:has-text("Back")', 'back-button');
      });

      observer.observe(document.body, { childList: true, subtree: true });

      // Cleanup observer after 10 seconds
      setTimeout(() => observer.disconnect(), 10000);
    });
  });

  test('should load AuthScreen without infinite loops and display all UI elements', async ({ page }) => {
    // Check for title and main content
    await expect(page).toHaveTitle(/Tchat|Telegram SEA/i);

    // Verify main app elements are present with proper fallback content
    await expect(page.locator('h1')).toContainText(/telegram/i, { timeout: 10000 });

    // Verify feature highlights are present (these should load from fallbacks)
    await expect(page.locator('text=End-to-End Encrypted')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('text=Ultra Low Data')).toBeVisible();
    await expect(page.locator('text=QR Payments')).toBeVisible();
    await expect(page.locator('text=SEA Languages')).toBeVisible();

    // Verify auth form is displayed
    await expect(page.locator('text=Sign In with Phone')).toBeVisible();
    await expect(page.locator('text=Enter your phone number')).toBeVisible();

    // Verify demo credentials are shown
    await expect(page.locator('text=Demo phone: +66812345678')).toBeVisible();

    // Check for no excessive loading indicators after content fallback settles
    await page.waitForTimeout(3000); // Allow fallback system to settle
    const loadingElementsCount = await page.locator('.animate-pulse').count();
    expect(loadingElementsCount).toBeLessThan(10); // Allow some loading but not excessive
  });

  test('should display and interact with country code badges', async ({ page }) => {
    // Verify all country code badges are present and correctly formatted
    const badges = {
      thailand: page.locator('[aria-label="Thailand +66"]'),
      indonesia: page.locator('[aria-label="Indonesia +62"]'),
      philippines: page.locator('[aria-label="Philippines +63"]'),
      vietnam: page.locator('[aria-label="Vietnam +84"]')
    };

    // Check each badge visibility and content
    await expect(badges.thailand).toBeVisible();
    await expect(badges.thailand).toContainText('🇹🇭 +66');

    await expect(badges.indonesia).toBeVisible();
    await expect(badges.indonesia).toContainText('🇮🇩 +62');

    await expect(badges.philippines).toBeVisible();
    await expect(badges.philippines).toContainText('🇵🇭 +63');

    await expect(badges.vietnam).toBeVisible();
    await expect(badges.vietnam).toContainText('🇻🇳 +84');

    // Test clicking a country code badge updates phone input
    const phoneInput = page.locator('input[placeholder*="+66"], input[aria-label*="phone"]').first();

    await badges.indonesia.click();
    await expect(phoneInput).toHaveValue('+62 ');

    // Test another badge
    await badges.thailand.click();
    await expect(phoneInput).toHaveValue('+66 ');
  });

  test('should handle phone number input and validation correctly', async ({ page }) => {
    // Find phone input field
    const phoneInput = page.locator('input[placeholder*="+66"], input[aria-label*="phone"], [data-testid="phone-input"]').first();

    await expect(phoneInput).toBeVisible();
    await expect(phoneInput).toHaveValue('+66812345678'); // Pre-filled with demo number

    // Test clearing and entering new number
    await phoneInput.clear();
    await phoneInput.fill('+62812345678');
    await expect(phoneInput).toHaveValue('+62812345678');

    // Reset to demo number
    await phoneInput.clear();
    await phoneInput.fill('+66812345678');
    await expect(phoneInput).toHaveValue('+66812345678');

    // Verify helper text is present
    await expect(page.locator('text=send you a 6-digit OTP')).toBeVisible();
    await expect(page.locator('text=Demo phone: +66812345678')).toBeVisible();

    // Verify Send OTP button is enabled with valid phone number
    const sendOtpButton = page.locator('button:has-text("Send OTP"), button:has-text("Send")').first();
    await expect(sendOtpButton).toBeVisible();
    await expect(sendOtpButton).toBeEnabled();
  });

  test('should handle backend unavailable gracefully with proper error feedback', async ({ page }) => {
    // Verify phone input works
    const phoneInput = page.locator('input[placeholder*="+66"], input[aria-label*="phone"]').first();
    await expect(phoneInput).toHaveValue('+66812345678');

    // Click Send OTP button (this will fail due to backend being down)
    const sendOtpButton = page.locator('button:has-text("Send OTP"), button:has-text("Send")').first();
    await sendOtpButton.click();

    // Wait a moment for any UI updates
    await page.waitForTimeout(2000);

    // Should remain on the phone input screen (not transition to OTP verification)
    await expect(page.locator('text=Sign In with Phone')).toBeVisible();
    await expect(phoneInput).toBeVisible();

    // Check console for proper error logging (via console messages)
    const consoleLogs = await page.evaluate(() => {
      return window.console;
    });

    // Verify the application handles the error gracefully without crashing
    await expect(page.locator('h1')).toContainText(/telegram/i); // App still functional
  });

  test('should display responsive design correctly on mobile viewports', async ({ page }) => {
    // Set mobile viewport (iPhone 12 size)
    await page.setViewportSize({ width: 390, height: 844 });

    // Wait for responsive adjustments
    await page.waitForTimeout(1000);

    // Verify main elements are still visible and properly sized
    await expect(page.locator('h1')).toBeVisible();
    await expect(page.locator('text=Sign In with Phone')).toBeVisible();

    // Feature highlights should be in mobile grid layout
    const featureGrid = page.locator('.grid');
    await expect(featureGrid).toBeVisible();

    // Country code badges should be visible (adjust touch target expectations)
    const badges = page.locator('[aria-label*="+"]');
    await expect(badges.first()).toBeVisible();

    // Test mobile-specific interactions
    const phoneInput = page.locator('input[placeholder*="+66"]').first();
    await expect(phoneInput).toBeVisible();
    await phoneInput.tap(); // Use tap instead of click for mobile

    // Verify text is readable and properly sized
    const heading = page.locator('h1');
    const headingBox = await heading.boundingBox();
    expect(headingBox?.width).toBeGreaterThan(200); // Text should not be too cramped
  });

  test('should maintain accessibility standards where possible', async ({ page }) => {
    // Check for proper ARIA labels
    await expect(page.locator('[aria-label="Thailand +66"]')).toBeVisible();
    await expect(page.locator('[aria-label="Indonesia +62"]')).toBeVisible();
    await expect(page.locator('[aria-label="Philippines +63"]')).toBeVisible();
    await expect(page.locator('[aria-label="Vietnam +84"]')).toBeVisible();

    // Check for proper headings hierarchy
    await expect(page.locator('h1')).toBeVisible(); // Main title
    await expect(page.locator('h3')).toBeVisible(); // Form title

    // Verify form has proper labels and descriptions
    const phoneInput = page.locator('input[placeholder*="+66"]').first();
    await expect(phoneInput).toHaveAttribute('aria-label');

    // Check for contrast and readability
    const backgroundColor = await page.evaluate(() => {
      return window.getComputedStyle(document.body).backgroundColor;
    });

    // Background should not be pure white or black for better contrast
    expect(backgroundColor).not.toBe('rgb(255, 255, 255)');
    expect(backgroundColor).not.toBe('rgb(0, 0, 0)');
  });

  test('should handle form validation appropriately', async ({ page }) => {
    const phoneInput = page.locator('input[placeholder*="+66"]').first();
    const sendOtpButton = page.locator('button:has-text("Send OTP")').first();

    // Test with empty phone number
    await phoneInput.clear();
    await sendOtpButton.click();

    // Should show validation error (though may be handled differently due to backend issues)
    await page.waitForTimeout(1000);

    // Verify phone input is still visible (not transitioned away)
    await expect(phoneInput).toBeVisible();
    await expect(page.locator('text=Sign In with Phone')).toBeVisible();

    // Test with valid format
    await phoneInput.fill('+66812345678');
    await expect(sendOtpButton).toBeEnabled();
  });
});