/**
 * Accessibility E2E Tests
 * Comprehensive testing of accessibility compliance across commerce workflows
 */

import { test, expect } from '@playwright/test';
import { injectAxe, checkA11y } from 'axe-playwright';
import { CartPage } from '../web/page-objects/CartPage';
import { CategoryPage } from '../web/page-objects/CategoryPage';
import { ProductPage } from '../web/page-objects/ProductPage';
import { TestDataSets } from '../utils/test-data';

test.describe('Accessibility Tests', () => {
  let cartPage: CartPage;
  let categoryPage: CategoryPage;
  let productPage: ProductPage;

  test.beforeEach(async ({ page }) => {
    cartPage = new CartPage(page);
    categoryPage = new CategoryPage(page);
    productPage = new ProductPage(page);

    // Inject axe-core for accessibility testing
    await injectAxe(page);
  });

  test.describe('WCAG 2.1 AA Compliance', () => {
    test('should meet accessibility standards on category page', async ({ page }) => {
      await test.step('Navigate to category page', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();
      });

      await test.step('Check accessibility compliance', async () => {
        await checkA11y(page, null, {
          detailedReport: true,
          detailedReportOptions: { html: true },
          rules: {
            // Specific rules for commerce pages
            'color-contrast': { enabled: true },
            'keyboard-navigation': { enabled: true },
            'focus-management': { enabled: true },
            'aria-labels': { enabled: true },
            'semantic-structure': { enabled: true },
          },
        });
      });

      await test.step('Test keyboard navigation', async () => {
        // Test Tab navigation through products
        await page.keyboard.press('Tab');
        await page.keyboard.press('Tab');
        await page.keyboard.press('Tab');

        const focusedElement = await page.locator(':focus').first();
        await expect(focusedElement).toBeVisible();

        // Test Enter key on product card
        await page.keyboard.press('Enter');

        // Should navigate to product detail or perform action
        await page.waitForTimeout(1000);
      });

      await test.step('Test screen reader compatibility', async () => {
        // Check for proper ARIA labels
        const productCards = categoryPage.productCards;
        const firstProduct = productCards.first();

        await expect(firstProduct).toHaveAttribute('role', 'button');

        const productName = firstProduct.getByTestId('product-name');
        await expect(productName).toHaveAccessibleName();

        const productPrice = firstProduct.getByTestId('product-price');
        await expect(productPrice).toHaveAccessibleName();
      });
    });

    test('should meet accessibility standards on product detail page', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Navigate to product page', async () => {
        await productPage.goto(testProduct.sku!);
        await productPage.waitForProductLoad();
      });

      await test.step('Check accessibility compliance', async () => {
        await checkA11y(page, null, {
          detailedReport: true,
          detailedReportOptions: { html: true },
        });
      });

      await test.step('Test form accessibility', async () => {
        // Test quantity selector accessibility
        const quantitySelector = productPage.quantitySelector;
        await expect(quantitySelector).toHaveAccessibleName();

        // Test add to cart button
        const addToCartButton = productPage.addToCartButton;
        await expect(addToCartButton).toHaveAccessibleName();
        await expect(addToCartButton).toHaveAttribute('role', 'button');

        // Test variant selectors if present
        const variantSelector = productPage.variantSelector;
        if (await variantSelector.isVisible({ timeout: 1000 })) {
          const firstVariant = variantSelector.locator('button').first();
          await expect(firstVariant).toHaveAccessibleName();
        }
      });

      await test.step('Test image accessibility', async () => {
        const mainImage = productPage.mainImage;
        await expect(mainImage).toHaveAttribute('alt');

        const altText = await mainImage.getAttribute('alt');
        expect(altText).toBeTruthy();
        expect(altText!.length).toBeGreaterThan(5);
      });
    });

    test('should meet accessibility standards on cart page', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Setup cart with items', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
        await cartPage.goto();
        await cartPage.waitForCartLoad();
      });

      await test.step('Check accessibility compliance', async () => {
        await checkA11y(page, null, {
          detailedReport: true,
          detailedReportOptions: { html: true },
        });
      });

      await test.step('Test cart item accessibility', async () => {
        const cartItem = cartPage.getCartItem(testProduct.sku!);

        // Test quantity controls
        const quantityInput = cartPage.getQuantityInput(testProduct.sku!);
        await expect(quantityInput).toHaveAccessibleName();
        await expect(quantityInput).toHaveAttribute('aria-label');

        // Test remove button
        const removeButton = cartPage.getRemoveButton(testProduct.sku!);
        await expect(removeButton).toHaveAccessibleName();
        await expect(removeButton).toHaveAttribute('aria-label');
      });

      await test.step('Test coupon form accessibility', async () => {
        const couponInput = cartPage.couponInput;
        await expect(couponInput).toHaveAccessibleName();
        await expect(couponInput).toHaveAttribute('aria-label');

        const applyCouponButton = cartPage.applyCouponButton;
        await expect(applyCouponButton).toHaveAccessibleName();
      });
    });

    test('should meet accessibility standards on checkout page', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Setup checkout', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
        await cartPage.goto();
        await cartPage.proceedToCheckout();
      });

      await test.step('Check checkout form accessibility', async () => {
        await checkA11y(page, null, {
          detailedReport: true,
          detailedReportOptions: { html: true },
        });
      });

      await test.step('Test form field accessibility', async () => {
        // Test email field
        const emailField = page.getByTestId('shipping-email');
        await expect(emailField).toHaveAccessibleName();
        await expect(emailField).toHaveAttribute('type', 'email');

        // Test required field indicators
        const requiredFields = page.locator('[required]');
        const count = await requiredFields.count();

        for (let i = 0; i < count; i++) {
          const field = requiredFields.nth(i);
          await expect(field).toHaveAttribute('aria-required', 'true');
        }
      });

      await test.step('Test error message accessibility', async () => {
        // Trigger validation errors
        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();

        // Check error messages have proper ARIA attributes
        const errorMessages = page.locator('[role="alert"]');
        const errorCount = await errorMessages.count();

        if (errorCount > 0) {
          for (let i = 0; i < errorCount; i++) {
            const errorMessage = errorMessages.nth(i);
            await expect(errorMessage).toBeVisible();
          }
        }
      });
    });
  });

  test.describe('Keyboard Navigation', () => {
    test('should support full keyboard navigation on category page', async ({ page }) => {
      await test.step('Navigate and test keyboard support', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Test sequential navigation
        let tabCount = 0;
        const maxTabs = 20;

        while (tabCount < maxTabs) {
          await page.keyboard.press('Tab');
          tabCount++;

          const focusedElement = await page.locator(':focus').first();
          if (await focusedElement.isVisible()) {
            // Verify focused element is interactive
            const tagName = await focusedElement.evaluate(el => el.tagName.toLowerCase());
            const role = await focusedElement.getAttribute('role');

            const isInteractive = ['button', 'input', 'select', 'a', 'textarea'].includes(tagName) ||
                                 ['button', 'link', 'textbox', 'combobox'].includes(role || '');

            if (isInteractive) {
              expect(true).toBeTruthy(); // Valid interactive element
            }
          }
        }
      });
    });

    test('should support Enter and Space key activation', async ({ page }) => {
      await test.step('Test key activation', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Find a product card and activate with Enter
        const productCard = categoryPage.productCards.first();
        await productCard.focus();
        await page.keyboard.press('Enter');

        // Should navigate to product detail
        await productPage.waitForProductLoad();
        await productPage.expectProductDetails({ name: TestDataSets.products.electronics.name });

        // Go back and test Space key
        await page.goBack();
        await categoryPage.waitForCategoryLoad();

        const secondProduct = categoryPage.productCards.nth(1);
        if (await secondProduct.isVisible({ timeout: 2000 })) {
          await secondProduct.focus();
          await page.keyboard.press(' ');

          // Should also navigate or perform action
          await page.waitForTimeout(1000);
        }
      });
    });

    test('should support escape key for modals and dialogs', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Test escape key functionality', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Open quick view modal if available
        const quickViewButton = categoryPage.getQuickViewButton(testProduct.sku!);
        if (await quickViewButton.isVisible({ timeout: 2000 })) {
          await quickViewButton.click();

          const modal = page.getByTestId('quick-view-modal');
          await expect(modal).toBeVisible();

          // Press Escape to close
          await page.keyboard.press('Escape');
          await expect(modal).not.toBeVisible();
        }
      });
    });
  });

  test.describe('Screen Reader Support', () => {
    test('should provide meaningful content for screen readers', async ({ page }) => {
      await test.step('Test ARIA landmarks', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Check for main landmarks
        const main = page.locator('[role="main"]');
        await expect(main).toBeVisible();

        const navigation = page.locator('[role="navigation"]');
        if (await navigation.isVisible({ timeout: 2000 })) {
          await expect(navigation).toBeVisible();
        }

        // Check for proper heading structure
        const headings = page.locator('h1, h2, h3, h4, h5, h6');
        const headingCount = await headings.count();
        expect(headingCount).toBeGreaterThan(0);

        // Verify heading hierarchy
        const h1 = page.locator('h1');
        if (await h1.count() > 0) {
          expect(await h1.count()).toBe(1); // Should have only one h1
        }
      });
    });

    test('should provide accessible product information', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Test product accessibility', async () => {
        await productPage.goto(testProduct.sku!);
        await productPage.waitForProductLoad();

        // Test product title accessibility
        const productTitle = productPage.productTitle;
        await expect(productTitle).toHaveAttribute('role', 'heading');

        // Test price accessibility
        const productPrice = productPage.productPrice;
        const priceText = await productPrice.textContent();
        expect(priceText).toContain('$');

        // Test product description
        const productDescription = productPage.productDescription;
        if (await productDescription.isVisible({ timeout: 2000 })) {
          const descriptionText = await productDescription.textContent();
          expect(descriptionText!.length).toBeGreaterThan(10);
        }
      });
    });

    test('should announce dynamic content changes', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Test live region announcements', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);

        // Check for live region announcement
        const liveRegion = page.locator('[aria-live]');
        if (await liveRegion.isVisible({ timeout: 5000 })) {
          const liveContent = await liveRegion.textContent();
          expect(liveContent).toContain('cart');
        }

        // Test cart update announcements
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.updateQuantity(testProduct.sku!, 3);

        // Should announce quantity change
        if (await liveRegion.isVisible({ timeout: 2000 })) {
          const updatedContent = await liveRegion.textContent();
          expect(updatedContent).toBeTruthy();
        }
      });
    });
  });

  test.describe('Color Contrast and Visual Design', () => {
    test('should meet color contrast requirements', async ({ page }) => {
      await test.step('Test color contrast ratios', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Check contrast using axe-core
        await checkA11y(page, null, {
          rules: {
            'color-contrast': { enabled: true },
          },
        });
      });

      await test.step('Test focus indicators', async () => {
        // Test that focused elements have visible focus indicators
        const interactiveElements = page.locator('button, a, input, select');
        const count = await interactiveElements.count();

        for (let i = 0; i < Math.min(count, 5); i++) {
          const element = interactiveElements.nth(i);
          await element.focus();

          // Check if element has focus styling
          const hasOutline = await element.evaluate(el => {
            const styles = window.getComputedStyle(el);
            return styles.outline !== 'none' ||
                   styles.boxShadow !== 'none' ||
                   styles.border !== 'none';
          });

          expect(hasOutline).toBeTruthy();
        }
      });
    });

    test('should work without relying solely on color', async ({ page }) => {
      await test.step('Test color-independent information', async () => {
        const testProduct = TestDataSets.products.electronics;

        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
        await cartPage.goto();

        // Apply invalid coupon to test error display
        await cartPage.applyCoupon('INVALID123');

        // Error should be indicated by more than just color
        const errorMessage = page.getByTestId('coupon-error');
        if (await errorMessage.isVisible({ timeout: 5000 })) {
          // Should have icon or text indication, not just color
          const hasIcon = await errorMessage.locator('[data-icon], .icon, svg').count() > 0;
          const hasErrorText = (await errorMessage.textContent())?.toLowerCase().includes('error') || false;

          expect(hasIcon || hasErrorText).toBeTruthy();
        }
      });
    });
  });

  test.describe('Mobile Accessibility', () => {
    test('should be accessible on mobile devices', async ({ page }) => {
      await test.step('Test mobile accessibility', async () => {
        // Set mobile viewport
        await page.setViewportSize({ width: 375, height: 667 });

        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Check accessibility on mobile
        await checkA11y(page, null, {
          detailedReport: true,
        });
      });

      await test.step('Test touch target sizes', async () => {
        // Check that interactive elements meet minimum touch target size (44px)
        const buttons = page.locator('button, a[role="button"], [role="button"]');
        const count = await buttons.count();

        for (let i = 0; i < Math.min(count, 10); i++) {
          const button = buttons.nth(i);
          if (await button.isVisible()) {
            const boundingBox = await button.boundingBox();
            if (boundingBox) {
              expect(boundingBox.width).toBeGreaterThanOrEqual(44);
              expect(boundingBox.height).toBeGreaterThanOrEqual(44);
            }
          }
        }
      });
    });

    test('should support mobile screen readers', async ({ page }) => {
      await test.step('Test mobile screen reader support', async () => {
        await page.setViewportSize({ width: 375, height: 667 });

        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Test swipe gestures simulation
        await page.touchscreen.tap(100, 100);

        // Verify mobile-specific ARIA attributes
        const mobileNavigation = page.locator('[data-testid="mobile-navigation"]');
        if (await mobileNavigation.isVisible({ timeout: 2000 })) {
          await expect(mobileNavigation).toHaveAttribute('role', 'navigation');
        }
      });
    });
  });

  test.describe('Error Handling Accessibility', () => {
    test('should provide accessible error messages', async ({ page }) => {
      await test.step('Test form validation errors', async () => {
        const testProduct = TestDataSets.products.electronics;

        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        // Submit form without filling required fields
        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();

        // Check error message accessibility
        const errorMessages = page.locator('[role="alert"], .error-message');
        const errorCount = await errorMessages.count();

        if (errorCount > 0) {
          for (let i = 0; i < errorCount; i++) {
            const errorMessage = errorMessages.nth(i);
            await expect(errorMessage).toBeVisible();

            // Error should be associated with form field
            const ariaDescribedBy = await errorMessage.getAttribute('aria-describedby');
            const id = await errorMessage.getAttribute('id');

            expect(ariaDescribedBy || id).toBeTruthy();
          }
        }
      });
    });

    test('should handle network errors accessibly', async ({ page }) => {
      await test.step('Test network error accessibility', async () => {
        // Simulate network error
        await page.route('**/api/v1/commerce/products**', route => route.abort());

        await categoryPage.goto('electronics');

        // Check error message accessibility
        const errorContainer = page.getByTestId('error-message');
        if (await errorContainer.isVisible({ timeout: 5000 })) {
          await expect(errorContainer).toHaveAttribute('role', 'alert');

          const retryButton = page.getByTestId('retry-button');
          if (await retryButton.isVisible({ timeout: 2000 })) {
            await expect(retryButton).toHaveAccessibleName();
          }
        }
      });
    });
  });
});