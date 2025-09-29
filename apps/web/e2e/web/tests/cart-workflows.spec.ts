/**
 * Cart Workflows E2E Tests
 * Comprehensive testing of cart functionality and user workflows
 */

import { test, expect } from '@playwright/test';
import { CartPage } from '../page-objects/CartPage';
import { CategoryPage } from '../page-objects/CategoryPage';
import { ProductPage } from '../page-objects/ProductPage';
import { TestDataGenerator, TestDataSets } from '../../utils/test-data';
import { createApiHelpers } from '../../utils/api-helpers';
import { ScreenshotHelpers } from '../../utils/screenshot-helpers';

test.describe('Cart Workflows', () => {
  let cartPage: CartPage;
  let categoryPage: CategoryPage;
  let productPage: ProductPage;
  let apiHelpers: ReturnType<typeof createApiHelpers>;
  let screenshotHelpers: ScreenshotHelpers;

  test.beforeEach(async ({ page, request }) => {
    cartPage = new CartPage(page);
    categoryPage = new CategoryPage(page);
    productPage = new ProductPage(page);
    apiHelpers = createApiHelpers(request);
    screenshotHelpers = new ScreenshotHelpers(page);

    // Setup test data
    await test.step('Setup test environment', async () => {
      // Clear any existing cart
      await page.goto('/');
      await page.waitForLoadState('networkidle');
    });
  });

  test.afterEach(async ({ page }) => {
    // Cleanup: Clear cart after each test
    await cartPage.goto();
    await cartPage.clearCart();
  });

  test.describe('Basic Cart Operations', () => {
    test('should display empty cart message when no items', async () => {
      await test.step('Navigate to empty cart', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();
      });

      await test.step('Verify empty cart state', async () => {
        await cartPage.expectCartEmpty();
        await screenshotHelpers.expectScreenshot({
          name: 'empty-cart',
          fullPage: true,
        });
      });
    });

    test('should add single item to cart from category page', async () => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Browse to category and add item', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Mock product in category if needed
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Verify item in cart', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.expectCartHasItems(1);
        await cartPage.expectCartItemExists(testProduct.sku!);
        await cartPage.expectCartItemDetails(testProduct.sku!, {
          name: testProduct.name,
          quantity: 1,
        });
      });

      await test.step('Take screenshot of cart with item', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'cart-with-single-item',
          fullPage: true,
        });
      });
    });

    test('should add multiple different items to cart', async () => {
      const products = [
        TestDataSets.products.electronics,
        TestDataSets.products.clothing,
        TestDataSets.products.books,
      ];

      await test.step('Add multiple products to cart', async () => {
        for (const product of products) {
          await categoryPage.goto(product.category.toLowerCase());
          await categoryPage.addProductToCart(product.sku!, 1);
        }
      });

      await test.step('Verify all items in cart', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.expectCartHasItems(products.length);

        for (const product of products) {
          await cartPage.expectCartItemExists(product.sku!);
          await cartPage.expectCartItemDetails(product.sku!, {
            name: product.name,
            quantity: 1,
          });
        }
      });

      await test.step('Take screenshot of cart with multiple items', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'cart-with-multiple-items',
          fullPage: true,
        });
      });
    });

    test('should update item quantity in cart', async () => {
      const testProduct = TestDataSets.products.electronics;
      const initialQuantity = 1;
      const updatedQuantity = 3;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, initialQuantity);
      });

      await test.step('Update quantity in cart', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.updateQuantity(testProduct.sku!, updatedQuantity);
      });

      await test.step('Verify updated quantity', async () => {
        await cartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: updatedQuantity,
        });

        // Verify total price updated
        const expectedTotal = (testProduct.price * updatedQuantity).toFixed(2);
        await cartPage.expectCartTotals({
          total: `$${expectedTotal}`,
        });
      });

      await test.step('Take screenshot of updated cart', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'cart-updated-quantity',
          fullPage: true,
        });
      });
    });

    test('should remove item from cart', async () => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Remove item from cart', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.removeItem(testProduct.sku!);
      });

      await test.step('Verify item removed', async () => {
        await cartPage.expectCartEmpty();
      });
    });

    test('should clear entire cart', async () => {
      const products = [
        TestDataSets.products.electronics,
        TestDataSets.products.clothing,
      ];

      await test.step('Add multiple items to cart', async () => {
        for (const product of products) {
          await categoryPage.goto(product.category.toLowerCase());
          await categoryPage.addProductToCart(product.sku!, 1);
        }
      });

      await test.step('Clear entire cart', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.clearCart();
      });

      await test.step('Verify cart is empty', async () => {
        await cartPage.expectCartEmpty();
      });
    });
  });

  test.describe('Coupon Functionality', () => {
    test('should apply valid percentage coupon', async () => {
      const testProduct = TestDataSets.products.electronics;
      const coupon = TestDataSets.coupons.percentage;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Apply coupon', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.applyCoupon(coupon.code);
      });

      await test.step('Verify coupon applied', async () => {
        await cartPage.expectCouponApplied(coupon.code);

        // Calculate expected discount
        const subtotal = testProduct.price;
        const discountAmount = Math.min(
          (subtotal * coupon.value) / 100,
          coupon.maxDiscount || Infinity
        );
        const total = subtotal - discountAmount;

        await cartPage.expectCartTotals({
          subtotal: `$${subtotal.toFixed(2)}`,
          discount: `-$${discountAmount.toFixed(2)}`,
          total: `$${total.toFixed(2)}`,
        });
      });

      await test.step('Take screenshot with coupon applied', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'cart-with-percentage-coupon',
          fullPage: true,
        });
      });
    });

    test('should apply valid fixed amount coupon', async () => {
      const testProduct = TestDataSets.products.electronics;
      const coupon = TestDataSets.coupons.fixed;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Apply coupon', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.applyCoupon(coupon.code);
      });

      await test.step('Verify coupon applied', async () => {
        await cartPage.expectCouponApplied(coupon.code);

        const subtotal = testProduct.price;
        const discountAmount = coupon.value;
        const total = subtotal - discountAmount;

        await cartPage.expectCartTotals({
          subtotal: `$${subtotal.toFixed(2)}`,
          discount: `-$${discountAmount.toFixed(2)}`,
          total: `$${total.toFixed(2)}`,
        });
      });
    });

    test('should reject invalid coupon code', async () => {
      const testProduct = TestDataSets.products.electronics;
      const invalidCoupon = 'INVALID123';

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Try to apply invalid coupon', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.applyCoupon(invalidCoupon);
      });

      await test.step('Verify coupon error', async () => {
        await cartPage.expectCouponError('Invalid coupon code');
      });
    });

    test('should reject expired coupon', async () => {
      const testProduct = TestDataSets.products.electronics;
      const expiredCoupon = TestDataSets.coupons.expired;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Try to apply expired coupon', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.applyCoupon(expiredCoupon.code);
      });

      await test.step('Verify coupon error', async () => {
        await cartPage.expectCouponError('Coupon has expired');
      });
    });

    test('should remove applied coupon', async () => {
      const testProduct = TestDataSets.products.electronics;
      const coupon = TestDataSets.coupons.percentage;

      await test.step('Add item and apply coupon', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);

        await cartPage.goto();
        await cartPage.waitForCartLoad();
        await cartPage.applyCoupon(coupon.code);
        await cartPage.expectCouponApplied(coupon.code);
      });

      await test.step('Remove coupon', async () => {
        await cartPage.removeCoupon();
      });

      await test.step('Verify coupon removed', async () => {
        const removeCouponButton = cartPage.removeCouponButton;
        await expect(removeCouponButton).not.toBeVisible();

        // Verify totals are back to original
        await cartPage.expectCartTotals({
          total: `$${testProduct.price.toFixed(2)}`,
        });
      });
    });
  });

  test.describe('Cart Persistence', () => {
    test('should persist cart across page reloads', async () => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 2);
      });

      await test.step('Test cart persistence', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();
        await cartPage.testCartPersistence();
      });
    });

    test('should persist cart across browser sessions', async ({ context }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Close and reopen browser', async () => {
        // Close current page
        await cartPage.page.close();

        // Create new page in same context (preserves localStorage)
        const newPage = await context.newPage();
        cartPage = new CartPage(newPage);
      });

      await test.step('Verify cart persisted', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.expectCartHasItems(1);
        await cartPage.expectCartItemExists(testProduct.sku!);
      });
    });
  });

  test.describe('Cart Edge Cases', () => {
    test('should handle out of stock items', async () => {
      const outOfStockProduct = TestDataSets.products.outOfStock;

      await test.step('Try to add out of stock item', async () => {
        await productPage.goto(outOfStockProduct.sku!);
        await productPage.waitForProductLoad();
      });

      await test.step('Verify add to cart is disabled', async () => {
        await productPage.expectProductDetails({
          availability: 'out-of-stock',
        });

        await expect(productPage.addToCartButton).toBeDisabled();
      });
    });

    test('should handle maximum quantity limits', async () => {
      const testProduct = TestDataSets.products.electronics;
      const maxQuantity = 5;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Try to exceed maximum quantity', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        // Try to set quantity beyond limit
        await cartPage.updateQuantity(testProduct.sku!, maxQuantity + 1);
      });

      await test.step('Verify quantity is limited', async () => {
        // Verify quantity is capped at maximum
        await cartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: maxQuantity,
        });

        // Verify error message
        const quantityError = cartPage.page.getByTestId('quantity-error');
        await expect(quantityError).toBeVisible();
        await expect(quantityError).toContainText('Maximum quantity exceeded');
      });
    });

    test('should handle network errors gracefully', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Add item to cart normally', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Simulate network error', async () => {
        // Intercept and fail cart update requests
        await page.route('**/api/v1/commerce/cart/**', route => {
          route.abort('failed');
        });

        await cartPage.goto();
        await cartPage.waitForCartLoad();

        // Try to update quantity
        await cartPage.updateQuantity(testProduct.sku!, 3);
      });

      await test.step('Verify error handling', async () => {
        const errorMessage = page.getByTestId('cart-error');
        await expect(errorMessage).toBeVisible();
        await expect(errorMessage).toContainText('Unable to update cart');

        // Verify quantity hasn't changed
        await cartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 1,
        });
      });
    });
  });

  test.describe('Checkout Initiation', () => {
    test('should proceed to checkout with items in cart', async () => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Add item to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Proceed to checkout', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.proceedToCheckout();
      });

      await test.step('Verify checkout page loaded', async () => {
        await expect(cartPage.page).toHaveURL(/.*\/checkout.*/);

        const checkoutContainer = cartPage.page.getByTestId('checkout-container');
        await expect(checkoutContainer).toBeVisible();
      });
    });

    test('should disable checkout button for empty cart', async () => {
      await test.step('Go to empty cart', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();
      });

      await test.step('Verify checkout is disabled', async () => {
        await cartPage.expectCartEmpty();
        await expect(cartPage.checkoutButton).not.toBeVisible();
      });
    });
  });

  test.describe('Mobile Cart Experience', () => {
    test('should work correctly on mobile devices', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Set mobile viewport', async () => {
        await page.setViewportSize({ width: 375, height: 667 });
      });

      await test.step('Add item to cart on mobile', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();
        await categoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Test mobile cart functionality', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();
        await cartPage.testResponsiveLayout();

        await cartPage.expectCartHasItems(1);
        await cartPage.expectCartItemExists(testProduct.sku!);
      });

      await test.step('Take mobile cart screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'cart-mobile-view',
          fullPage: true,
        });
      });
    });
  });
});