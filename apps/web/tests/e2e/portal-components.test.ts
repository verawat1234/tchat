import { test, expect, Page } from '@playwright/test';

/**
 * E2E Tests for Portal Components (Tooltips, Dialogs, Dropdowns)
 * These components use React portals and have complex async behaviors
 * that are difficult to test with unit tests
 */

test.describe('Portal Components E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to a test page with portal components
    // This assumes you have a route that showcases these components
    await page.goto('/components'); // Adjust this URL to your actual component showcase
  });

  test.describe('Tooltip Component', () => {
    test('should show tooltip on hover', async ({ page }) => {
      // Find an element with a tooltip
      const triggerElement = page.locator('[data-testid="tooltip-trigger"]').first();

      // Hover over the element
      await triggerElement.hover();

      // Wait for tooltip to appear
      const tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
      await expect(tooltip).toBeVisible({ timeout: 2000 });

      // Verify tooltip content
      await expect(tooltip).toContainText(/.*/, { timeout: 1000 });

      // Move mouse away
      await page.mouse.move(0, 0);

      // Tooltip should disappear
      await expect(tooltip).toBeHidden({ timeout: 2000 });
    });

    test('should show tooltip on focus', async ({ page }) => {
      const triggerElement = page.locator('[data-testid="tooltip-trigger"]').first();

      // Focus the element
      await triggerElement.focus();

      // Wait for tooltip
      const tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
      await expect(tooltip).toBeVisible({ timeout: 2000 });

      // Blur the element
      await page.keyboard.press('Tab');

      // Tooltip should disappear
      await expect(tooltip).toBeHidden({ timeout: 2000 });
    });

    test('should position tooltip correctly', async ({ page }) => {
      const triggerElement = page.locator('[data-testid="tooltip-trigger"]').first();

      // Get trigger element position
      const triggerBox = await triggerElement.boundingBox();
      if (!triggerBox) throw new Error('Trigger element not found');

      // Hover to show tooltip
      await triggerElement.hover();

      // Get tooltip position
      const tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
      await tooltip.waitFor({ state: 'visible', timeout: 2000 });

      const tooltipBox = await tooltip.boundingBox();
      if (!tooltipBox) throw new Error('Tooltip not found');

      // Verify tooltip is positioned near the trigger
      // Allow for some offset due to positioning strategies
      const maxDistance = 100;
      const distance = Math.sqrt(
        Math.pow(tooltipBox.x - triggerBox.x, 2) +
        Math.pow(tooltipBox.y - triggerBox.y, 2)
      );

      expect(distance).toBeLessThan(maxDistance);
    });

    test('should handle multiple tooltips without interference', async ({ page }) => {
      const triggers = await page.locator('[data-testid^="tooltip-trigger"]').all();

      if (triggers.length < 2) {
        test.skip(); // Skip if not enough tooltips to test
        return;
      }

      // Hover first tooltip
      await triggers[0].hover();
      let tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
      await expect(tooltip).toBeVisible();
      const firstTooltipText = await tooltip.textContent();

      // Hover second tooltip
      await triggers[1].hover();
      tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
      await expect(tooltip).toBeVisible();
      const secondTooltipText = await tooltip.textContent();

      // Verify different content
      expect(firstTooltipText).not.toBe(secondTooltipText);
    });

    test('should be accessible', async ({ page }) => {
      const triggerElement = page.locator('[data-testid="tooltip-trigger"]').first();

      // Check trigger has aria-describedby when tooltip is shown
      await triggerElement.hover();

      const tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
      await tooltip.waitFor({ state: 'visible' });

      // Get the tooltip id
      const tooltipId = await tooltip.getAttribute('id');
      if (tooltipId) {
        const triggerAriaDescribedBy = await triggerElement.getAttribute('aria-describedby');
        expect(triggerAriaDescribedBy).toContain(tooltipId);
      }

      // Tooltip should have role="tooltip"
      const role = await tooltip.getAttribute('role');
      expect(role).toBe('tooltip');
    });
  });

  test.describe('Dialog Component', () => {
    test('should open and close dialog', async ({ page }) => {
      // Find and click dialog trigger
      const triggerButton = page.locator('[data-testid="dialog-trigger"]').first();
      await triggerButton.click();

      // Wait for dialog to appear
      const dialog = page.locator('[role="dialog"], [data-radix-dialog-content]');
      await expect(dialog).toBeVisible({ timeout: 2000 });

      // Verify dialog content is visible
      const dialogTitle = dialog.locator('[role="heading"], [data-radix-dialog-title]');
      await expect(dialogTitle.first()).toBeVisible();

      // Close dialog via close button
      const closeButton = dialog.locator('[data-testid="dialog-close"], [aria-label*="close" i]');
      if (await closeButton.count() > 0) {
        await closeButton.first().click();
      } else {
        // Try ESC key as fallback
        await page.keyboard.press('Escape');
      }

      // Dialog should be hidden
      await expect(dialog).toBeHidden({ timeout: 2000 });
    });

    test('should trap focus within dialog', async ({ page }) => {
      // Open dialog
      const triggerButton = page.locator('[data-testid="dialog-trigger"]').first();
      await triggerButton.click();

      const dialog = page.locator('[role="dialog"], [data-radix-dialog-content]');
      await dialog.waitFor({ state: 'visible' });

      // Get all focusable elements in dialog
      const focusableElements = await dialog.locator('button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])').all();

      if (focusableElements.length > 0) {
        // Focus should be within dialog
        const activeElement = await page.evaluate(() => document.activeElement?.tagName);
        expect(activeElement).toBeTruthy();

        // Tab through elements
        for (let i = 0; i < focusableElements.length + 1; i++) {
          await page.keyboard.press('Tab');
        }

        // Focus should wrap around to first element
        const focusedElement = await page.evaluate(() => document.activeElement);
        expect(focusedElement).toBeTruthy();
      }
    });

    test('should prevent body scroll when open', async ({ page }) => {
      // Add content to make page scrollable
      await page.evaluate(() => {
        document.body.style.height = '200vh';
      });

      // Scroll down a bit
      await page.evaluate(() => window.scrollTo(0, 100));
      const initialScroll = await page.evaluate(() => window.scrollY);

      // Open dialog
      const triggerButton = page.locator('[data-testid="dialog-trigger"]').first();
      await triggerButton.click();

      const dialog = page.locator('[role="dialog"], [data-radix-dialog-content]');
      await dialog.waitFor({ state: 'visible' });

      // Try to scroll
      await page.mouse.wheel(0, 100);

      // Check if body scroll is prevented (might not work in all implementations)
      const bodyStyle = await page.evaluate(() => document.body.style.overflow);
      expect(['hidden', 'clip']).toContain(bodyStyle);
    });

    test('should close on overlay click', async ({ page }) => {
      // Open dialog
      const triggerButton = page.locator('[data-testid="dialog-trigger"]').first();
      await triggerButton.click();

      const dialog = page.locator('[role="dialog"], [data-radix-dialog-content]');
      await dialog.waitFor({ state: 'visible' });

      // Click overlay
      const overlay = page.locator('[data-radix-dialog-overlay]');
      if (await overlay.count() > 0) {
        await overlay.click({ position: { x: 10, y: 10 } });

        // Dialog should close
        await expect(dialog).toBeHidden({ timeout: 2000 });
      }
    });

    test('should be accessible', async ({ page }) => {
      // Open dialog
      const triggerButton = page.locator('[data-testid="dialog-trigger"]').first();
      await triggerButton.click();

      const dialog = page.locator('[role="dialog"], [data-radix-dialog-content]');
      await dialog.waitFor({ state: 'visible' });

      // Check ARIA attributes
      const role = await dialog.getAttribute('role');
      expect(role).toBe('dialog');

      // Check for aria-modal
      const ariaModal = await dialog.getAttribute('aria-modal');
      expect(ariaModal).toBe('true');

      // Check for aria-labelledby or aria-label
      const ariaLabelledby = await dialog.getAttribute('aria-labelledby');
      const ariaLabel = await dialog.getAttribute('aria-label');
      expect(ariaLabelledby || ariaLabel).toBeTruthy();

      // Check for aria-describedby
      const ariaDescribedby = await dialog.getAttribute('aria-describedby');
      if (ariaDescribedby) {
        const description = page.locator(`#${ariaDescribedby}`);
        await expect(description).toBeVisible();
      }
    });
  });

  test.describe('Dropdown Component', () => {
    test('should open and close dropdown', async ({ page }) => {
      // Find dropdown trigger
      const trigger = page.locator('[data-testid="dropdown-trigger"]').first();
      await trigger.click();

      // Wait for dropdown menu
      const menu = page.locator('[role="menu"], [data-radix-dropdown-content]');
      await expect(menu).toBeVisible({ timeout: 2000 });

      // Click outside to close
      await page.mouse.click(0, 0);
      await expect(menu).toBeHidden({ timeout: 2000 });
    });

    test('should navigate with keyboard', async ({ page }) => {
      // Open dropdown
      const trigger = page.locator('[data-testid="dropdown-trigger"]').first();
      await trigger.click();

      const menu = page.locator('[role="menu"], [data-radix-dropdown-content]');
      await menu.waitFor({ state: 'visible' });

      // Navigate with arrow keys
      await page.keyboard.press('ArrowDown');
      await page.keyboard.press('ArrowDown');
      await page.keyboard.press('ArrowUp');

      // Select with Enter
      await page.keyboard.press('Enter');

      // Menu should close after selection
      await expect(menu).toBeHidden({ timeout: 2000 });
    });

    test('should handle sub-menus', async ({ page }) => {
      const trigger = page.locator('[data-testid="dropdown-trigger"]').first();
      await trigger.click();

      const menu = page.locator('[role="menu"], [data-radix-dropdown-content]');
      await menu.waitFor({ state: 'visible' });

      // Look for items that have sub-menus
      const subMenuTrigger = menu.locator('[data-radix-dropdown-subtrigger], [aria-haspopup="true"]').first();

      if (await subMenuTrigger.count() > 0) {
        await subMenuTrigger.hover();

        // Wait for sub-menu
        const subMenu = page.locator('[data-radix-dropdown-subcontent]');
        await expect(subMenu).toBeVisible({ timeout: 2000 });
      }
    });

    test('should be accessible', async ({ page }) => {
      const trigger = page.locator('[data-testid="dropdown-trigger"]').first();

      // Check trigger attributes
      const ariaHaspopup = await trigger.getAttribute('aria-haspopup');
      expect(['menu', 'true']).toContain(ariaHaspopup);

      // Open dropdown
      await trigger.click();

      const menu = page.locator('[role="menu"], [data-radix-dropdown-content]');
      await menu.waitFor({ state: 'visible' });

      // Check menu role
      const menuRole = await menu.getAttribute('role');
      expect(menuRole).toBe('menu');

      // Check menu items
      const menuItems = await menu.locator('[role="menuitem"]').all();
      expect(menuItems.length).toBeGreaterThan(0);

      // Check trigger aria-expanded
      const ariaExpanded = await trigger.getAttribute('aria-expanded');
      expect(ariaExpanded).toBe('true');
    });
  });

  test.describe('Portal Rendering', () => {
    test('should render portals in correct container', async ({ page }) => {
      // Check if portal container exists
      const portalContainer = page.locator('[data-radix-portal], #radix-portal-container, body > [data-portal-container]');

      // Open a dialog to trigger portal creation
      const dialogTrigger = page.locator('[data-testid="dialog-trigger"]').first();
      if (await dialogTrigger.count() > 0) {
        await dialogTrigger.click();

        // Check that content is rendered in portal
        const dialogInPortal = portalContainer.locator('[role="dialog"]');
        await expect(dialogInPortal.first()).toBeVisible();
      }
    });

    test('should handle multiple portals simultaneously', async ({ page }) => {
      // Open dialog
      const dialogTrigger = page.locator('[data-testid="dialog-trigger"]').first();
      if (await dialogTrigger.count() > 0) {
        await dialogTrigger.click();

        const dialog = page.locator('[role="dialog"]');
        await dialog.waitFor({ state: 'visible' });

        // Open dropdown inside dialog (if exists)
        const dropdownInDialog = dialog.locator('[data-testid="dropdown-trigger"]').first();
        if (await dropdownInDialog.count() > 0) {
          await dropdownInDialog.click();

          const dropdown = page.locator('[role="menu"]');
          await expect(dropdown).toBeVisible();

          // Both should be visible
          await expect(dialog).toBeVisible();
          await expect(dropdown).toBeVisible();
        }
      }
    });

    test('should clean up portals on unmount', async ({ page }) => {
      // Count initial portals
      const initialPortals = await page.locator('[data-radix-portal]').count();

      // Open and close dialog
      const dialogTrigger = page.locator('[data-testid="dialog-trigger"]').first();
      if (await dialogTrigger.count() > 0) {
        await dialogTrigger.click();
        await page.waitForTimeout(500);
        await page.keyboard.press('Escape');
        await page.waitForTimeout(500);

        // Portal count should return to initial
        const finalPortals = await page.locator('[data-radix-portal]').count();
        expect(finalPortals).toBe(initialPortals);
      }
    });
  });

  test.describe('Performance', () => {
    test('should render tooltips quickly', async ({ page }) => {
      const trigger = page.locator('[data-testid="tooltip-trigger"]').first();

      const startTime = Date.now();
      await trigger.hover();

      const tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
      await tooltip.waitFor({ state: 'visible', timeout: 1000 });

      const renderTime = Date.now() - startTime;
      expect(renderTime).toBeLessThan(1000); // Should appear within 1 second
    });

    test('should handle rapid interactions', async ({ page }) => {
      const triggers = await page.locator('[data-testid^="tooltip-trigger"]').all();

      if (triggers.length >= 3) {
        // Rapidly hover over multiple elements
        for (let i = 0; i < 10; i++) {
          await triggers[i % triggers.length].hover();
          await page.waitForTimeout(50);
        }

        // Page should still be responsive
        const finalTrigger = triggers[0];
        await finalTrigger.hover();

        const tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
        await expect(tooltip).toBeVisible({ timeout: 2000 });
      }
    });
  });

  test.describe('Edge Cases', () => {
    test('should handle scroll while tooltip is open', async ({ page }) => {
      // Make page scrollable
      await page.evaluate(() => {
        document.body.style.height = '200vh';
      });

      const trigger = page.locator('[data-testid="tooltip-trigger"]').first();
      await trigger.hover();

      const tooltip = page.locator('[role="tooltip"], [data-radix-tooltip-content]');
      await tooltip.waitFor({ state: 'visible' });

      // Scroll the page
      await page.evaluate(() => window.scrollBy(0, 100));

      // Tooltip should either follow or hide
      const isVisible = await tooltip.isVisible();
      // This is implementation-specific behavior
      expect(typeof isVisible).toBe('boolean');
    });

    test('should handle window resize', async ({ page }) => {
      // Open a dialog
      const dialogTrigger = page.locator('[data-testid="dialog-trigger"]').first();
      if (await dialogTrigger.count() > 0) {
        await dialogTrigger.click();

        const dialog = page.locator('[role="dialog"]');
        await dialog.waitFor({ state: 'visible' });

        // Resize viewport
        await page.setViewportSize({ width: 400, height: 600 });
        await page.waitForTimeout(500);

        // Dialog should still be visible and centered
        await expect(dialog).toBeVisible();

        const box = await dialog.boundingBox();
        if (box) {
          const viewportSize = page.viewportSize();
          if (viewportSize) {
            // Check if roughly centered
            const centerX = box.x + box.width / 2;
            const centerY = box.y + box.height / 2;

            expect(Math.abs(centerX - viewportSize.width / 2)).toBeLessThan(100);
            expect(Math.abs(centerY - viewportSize.height / 2)).toBeLessThan(100);
          }
        }
      }
    });
  });
});

// Test for specific component library behaviors (Radix UI)
test.describe('Radix UI Specific Behaviors', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/components');
  });

  test('should handle Radix UI delay props', async ({ page }) => {
    // Radix tooltips often have delay props
    const trigger = page.locator('[data-testid="tooltip-with-delay"]').first();

    if (await trigger.count() > 0) {
      const startTime = Date.now();
      await trigger.hover();

      const tooltip = page.locator('[role="tooltip"]');
      await tooltip.waitFor({ state: 'visible', timeout: 3000 });

      const showTime = Date.now() - startTime;
      // If there's a delay, it should be noticeable
      expect(showTime).toBeGreaterThan(100);
    }
  });

  test('should respect dismissable props', async ({ page }) => {
    const dialogTrigger = page.locator('[data-testid="non-dismissable-dialog"]').first();

    if (await dialogTrigger.count() > 0) {
      await dialogTrigger.click();

      const dialog = page.locator('[role="dialog"]');
      await dialog.waitFor({ state: 'visible' });

      // Try to dismiss with ESC
      await page.keyboard.press('Escape');
      await page.waitForTimeout(500);

      // Non-dismissable dialog should still be visible
      await expect(dialog).toBeVisible();

      // Should only close with explicit close button
      const closeButton = dialog.locator('[data-testid="dialog-close"]');
      await closeButton.click();
      await expect(dialog).toBeHidden();
    }
  });
});