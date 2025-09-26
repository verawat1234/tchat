/**
 * Visual Regression Test: TchatDialog Web vs iOS
 * Validates cross-platform consistency for dialog components
 */

import { test, expect } from '@playwright/test';
import { VisualTestRunner } from '../visual-test-utils';
import { VISUAL_THRESHOLDS, STANDARD_VARIANTS, STANDARD_SIZES } from '../visual-config';

test.describe('TchatDialog Cross-Platform Consistency', () => {
  let visualTester: VisualTestRunner;

  test.beforeEach(async ({ page }) => {
    visualTester = new VisualTestRunner(page);

    // Set iOS-compatible viewport
    await visualTester.setIOSViewport('iPhone 12');

    // Navigate to Storybook
    await page.goto('http://localhost:6006');
  });

  test('should match iOS dialog appearance - all variants', async () => {
    const results = await visualTester.testComponentVariants(
      'TchatDialog',
      ['default', 'confirmation', 'destructive'],
      ['medium'], // Dialogs typically have fixed sizes
      {
        threshold: VISUAL_THRESHOLDS.STRICT, // Strict threshold for dialogs
        animations: 'disabled',
        fullPage: true // Capture overlay and backdrop
      }
    );

    // Validate results
    for (const result of results) {
      expect(result.passed, `${result.component} ${result.variant} ${result.size} failed consistency check`).toBe(true);
      expect(result.consistencyScore).toBeGreaterThan(0.95);
    }

    // Generate report
    const report = visualTester.generateConsistencyReport(results);
    console.log('TchatDialog Consistency Report:', report);

    expect(report.overallScore).toBeGreaterThan(0.95);
    expect(report.componentsPassingThreshold).toBe(report.componentsAnalyzed);
  });

  test('should handle modal interactions consistently', async ({ page }) => {
    await visualTester.navigateToComponent('TchatDialog', 'default');

    // Test open state
    await page.click('[data-testid="dialog-trigger"]');
    await expect(page.locator('[data-testid="dialog-content"]')).toBeVisible();

    // Capture modal state
    await expect(page).toHaveScreenshot('dialog-open-state.png', {
      threshold: VISUAL_THRESHOLDS.DEFAULT,
      fullPage: true,
      animations: 'disabled'
    });

    // Test close interaction
    await page.keyboard.press('Escape');
    await expect(page.locator('[data-testid="dialog-content"]')).not.toBeVisible();
  });

  test('should validate design token accuracy', async ({ page }) => {
    await visualTester.navigateToComponent('TchatDialog', 'default');
    await page.click('[data-testid="dialog-trigger"]');

    // Validate design tokens match iOS specifications
    const dialogStyles = await page.evaluate(() => {
      const dialog = document.querySelector('[data-testid="dialog-content"]');
      if (!dialog) return {};

      const computed = window.getComputedStyle(dialog);
      return {
        backgroundColor: computed.backgroundColor,
        borderRadius: computed.borderRadius,
        boxShadow: computed.boxShadow,
        padding: computed.padding,
        maxWidth: computed.maxWidth
      };
    });

    // These values should match iOS TchatDialog specifications
    expect(dialogStyles.backgroundColor).toBe('rgb(255, 255, 255)'); // White background
    expect(dialogStyles.borderRadius).toBe('12px'); // iOS-style rounded corners
  });

  test('should validate accessibility compliance', async ({ page }) => {
    await visualTester.navigateToComponent('TchatDialog', 'default');

    // Test keyboard navigation
    await page.keyboard.press('Tab');
    await page.keyboard.press('Enter'); // Open dialog

    await expect(page.locator('[data-testid="dialog-content"]')).toBeVisible();
    await expect(page.locator('[data-testid="dialog-content"]')).toBeFocused();

    // Test ARIA attributes
    const dialogContent = page.locator('[data-testid="dialog-content"]');
    await expect(dialogContent).toHaveAttribute('role', 'dialog');
    await expect(dialogContent).toHaveAttribute('aria-modal', 'true');
  });

  test('should handle responsive behavior', async ({ page }) => {
    const devices = ['iPhone 12', 'iPad Air'] as const;

    for (const device of devices) {
      await visualTester.setIOSViewport(device);
      await visualTester.navigateToComponent('TchatDialog', 'default');

      await page.click('[data-testid="dialog-trigger"]');
      await expect(page.locator('[data-testid="dialog-content"]')).toBeVisible();

      // Capture responsive behavior
      await expect(page).toHaveScreenshot(`dialog-${device.toLowerCase().replace(' ', '-')}.png`, {
        threshold: VISUAL_THRESHOLDS.DEFAULT,
        fullPage: true,
        animations: 'disabled'
      });
    }
  });
});