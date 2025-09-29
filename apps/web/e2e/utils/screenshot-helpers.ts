/**
 * Screenshot and Visual Testing Utilities
 * Provides utilities for visual regression testing and screenshot management
 */

import { Page, expect, Locator } from '@playwright/test';

export interface ScreenshotOptions {
  name: string;
  fullPage?: boolean;
  mask?: Locator[];
  clip?: { x: number; y: number; width: number; height: number };
  threshold?: number;
  maxDiffPixels?: number;
}

export interface VisualTestOptions {
  threshold?: number;
  maxDiffPixels?: number;
  animations?: 'disabled' | 'allow';
  caret?: 'hide' | 'initial';
}

export class ScreenshotHelpers {
  constructor(private page: Page) {}

  /**
   * Take a screenshot and compare with baseline
   */
  async expectScreenshot(options: ScreenshotOptions): Promise<void> {
    // Wait for page to be stable
    await this.waitForStableState();

    // Disable animations for consistent screenshots
    await this.disableAnimations();

    // Hide dynamic content if specified
    if (options.mask) {
      for (const element of options.mask) {
        await element.evaluate(el => {
          el.style.opacity = '0';
        });
      }
    }

    // Take screenshot and compare
    await expect(this.page).toHaveScreenshot(options.name, {
      fullPage: options.fullPage || false,
      clip: options.clip,
      threshold: options.threshold || 0.1,
      maxDiffPixels: options.maxDiffPixels || 100,
    });

    // Restore masked elements
    if (options.mask) {
      for (const element of options.mask) {
        await element.evaluate(el => {
          el.style.opacity = '';
        });
      }
    }
  }

  /**
   * Take a screenshot of a specific element
   */
  async expectElementScreenshot(
    element: Locator,
    name: string,
    options: VisualTestOptions = {}
  ): Promise<void> {
    await this.waitForStableState();
    await this.disableAnimations();

    await expect(element).toHaveScreenshot(name, {
      threshold: options.threshold || 0.1,
      maxDiffPixels: options.maxDiffPixels || 50,
      animations: options.animations || 'disabled',
      caret: options.caret || 'hide',
    });
  }

  /**
   * Compare multiple viewport sizes
   */
  async expectResponsiveScreenshots(
    name: string,
    viewports: { width: number; height: number; name: string }[]
  ): Promise<void> {
    for (const viewport of viewports) {
      await this.page.setViewportSize(viewport);
      await this.waitForStableState();
      await this.expectScreenshot({
        name: `${name}-${viewport.name}`,
        fullPage: true,
      });
    }
  }

  /**
   * Test component in different states
   */
  async expectComponentStates(
    component: Locator,
    name: string,
    states: { action: () => Promise<void>; suffix: string }[]
  ): Promise<void> {
    for (const state of states) {
      await state.action();
      await this.waitForStableState();
      await this.expectElementScreenshot(component, `${name}-${state.suffix}`);
    }
  }

  /**
   * Wait for page to be in a stable state for screenshots
   */
  private async waitForStableState(): Promise<void> {
    // Wait for network to be idle
    await this.page.waitForLoadState('networkidle');

    // Wait for fonts to load
    await this.page.evaluate(() => document.fonts.ready);

    // Wait for any pending animations to complete
    await this.page.waitForFunction(() => {
      const animations = document.getAnimations();
      return animations.every(animation => animation.playState === 'finished');
    });

    // Additional wait for React state updates
    await this.page.waitForTimeout(100);
  }

  /**
   * Disable animations for consistent screenshots
   */
  private async disableAnimations(): Promise<void> {
    await this.page.addStyleTag({
      content: `
        *, *::before, *::after {
          animation-duration: 0s !important;
          animation-delay: 0s !important;
          transition-duration: 0s !important;
          transition-delay: 0s !important;
        }
      `,
    });
  }

  /**
   * Hide dynamic content for screenshots
   */
  async hideDynamicContent(): Promise<void> {
    await this.page.addStyleTag({
      content: `
        /* Hide timestamps and dynamic content */
        [data-testid*="timestamp"],
        [data-testid*="time"],
        [data-testid*="date"],
        .timestamp,
        .time,
        .loading-spinner {
          opacity: 0 !important;
        }
      `,
    });
  }

  /**
   * Mock images for consistent screenshots
   */
  async mockImages(): Promise<void> {
    await this.page.route('**/*.{png,jpg,jpeg,gif,webp}', route => {
      route.fulfill({
        status: 200,
        contentType: 'image/png',
        body: Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==', 'base64'),
      });
    });
  }

  /**
   * Test cart visual states
   */
  async testCartVisualStates(cartLocator: Locator): Promise<void> {
    await this.expectComponentStates(cartLocator, 'cart', [
      {
        action: async () => {
          // Empty cart state
          await this.clearCart();
        },
        suffix: 'empty',
      },
      {
        action: async () => {
          // Single item state
          await this.addItemToCart();
        },
        suffix: 'single-item',
      },
      {
        action: async () => {
          // Multiple items state
          await this.addMultipleItemsToCart();
        },
        suffix: 'multiple-items',
      },
      {
        action: async () => {
          // With coupon state
          await this.applyCouponToCart();
        },
        suffix: 'with-coupon',
      },
    ]);
  }

  /**
   * Test product card visual states
   */
  async testProductCardStates(productCard: Locator): Promise<void> {
    await this.expectComponentStates(productCard, 'product-card', [
      {
        action: async () => {
          // Default state - no action needed
        },
        suffix: 'default',
      },
      {
        action: async () => {
          await productCard.hover();
        },
        suffix: 'hover',
      },
      {
        action: async () => {
          await productCard.focus();
        },
        suffix: 'focus',
      },
      {
        action: async () => {
          // Simulate out of stock
          await this.setProductOutOfStock(productCard);
        },
        suffix: 'out-of-stock',
      },
    ]);
  }

  /**
   * Test form visual states
   */
  async testFormStates(form: Locator): Promise<void> {
    const submitButton = form.locator('[type="submit"]');
    const firstInput = form.locator('input').first();

    await this.expectComponentStates(form, 'checkout-form', [
      {
        action: async () => {
          // Empty form state
        },
        suffix: 'empty',
      },
      {
        action: async () => {
          // Filled form state
          await this.fillCheckoutForm(form);
        },
        suffix: 'filled',
      },
      {
        action: async () => {
          // Form with validation errors
          await submitButton.click();
          await this.page.waitForSelector('[data-testid*="error"]');
        },
        suffix: 'validation-errors',
      },
      {
        action: async () => {
          // Loading state
          await this.triggerFormLoading(form);
        },
        suffix: 'loading',
      },
    ]);
  }

  /**
   * Helper methods for cart operations
   */
  private async clearCart(): Promise<void> {
    // Implementation would clear cart via API or UI
    await this.page.evaluate(() => {
      // Clear cart logic
    });
  }

  private async addItemToCart(): Promise<void> {
    // Implementation would add single item to cart
  }

  private async addMultipleItemsToCart(): Promise<void> {
    // Implementation would add multiple items to cart
  }

  private async applyCouponToCart(): Promise<void> {
    // Implementation would apply coupon to cart
  }

  private async setProductOutOfStock(productCard: Locator): Promise<void> {
    // Implementation would set product as out of stock
    await productCard.evaluate(el => {
      el.setAttribute('data-out-of-stock', 'true');
    });
  }

  private async fillCheckoutForm(form: Locator): Promise<void> {
    // Implementation would fill checkout form with valid data
    await form.locator('[name="email"]').fill('test@example.com');
    await form.locator('[name="firstName"]').fill('John');
    await form.locator('[name="lastName"]').fill('Doe');
  }

  private async triggerFormLoading(form: Locator): Promise<void> {
    // Implementation would trigger loading state
    await form.evaluate(el => {
      el.setAttribute('data-loading', 'true');
    });
  }
}

/**
 * Visual testing presets for common commerce scenarios
 */
export const VisualTestPresets = {
  cart: {
    threshold: 0.05,
    maxDiffPixels: 100,
  },
  productGrid: {
    threshold: 0.1,
    maxDiffPixels: 200,
  },
  checkout: {
    threshold: 0.03,
    maxDiffPixels: 50,
  },
  mobile: {
    threshold: 0.08,
    maxDiffPixels: 150,
  },
};

/**
 * Common viewport sizes for responsive testing
 */
export const CommonViewports = {
  mobile: { width: 375, height: 667, name: 'mobile' },
  tablet: { width: 768, height: 1024, name: 'tablet' },
  desktop: { width: 1280, height: 720, name: 'desktop' },
  ultrawide: { width: 1920, height: 1080, name: 'ultrawide' },
};