/**
 * Cart Page Object Model
 * Encapsulates cart-related interactions and selectors
 */

import { Page, Locator, expect } from '@playwright/test';
import { TestProduct } from '../../utils/test-data';

export class CartPage {
  readonly page: Page;
  readonly cartContainer: Locator;
  readonly cartItems: Locator;
  readonly emptyCartMessage: Locator;
  readonly cartTotal: Locator;
  readonly subtotal: Locator;
  readonly tax: Locator;
  readonly shipping: Locator;
  readonly discount: Locator;
  readonly couponInput: Locator;
  readonly applyCouponButton: Locator;
  readonly removeCouponButton: Locator;
  readonly checkoutButton: Locator;
  readonly continueShoppingButton: Locator;
  readonly clearCartButton: Locator;

  constructor(page: Page) {
    this.page = page;
    this.cartContainer = page.getByTestId('cart-container');
    this.cartItems = page.getByTestId('cart-items');
    this.emptyCartMessage = page.getByTestId('empty-cart-message');
    this.cartTotal = page.getByTestId('cart-total');
    this.subtotal = page.getByTestId('cart-subtotal');
    this.tax = page.getByTestId('cart-tax');
    this.shipping = page.getByTestId('cart-shipping');
    this.discount = page.getByTestId('cart-discount');
    this.couponInput = page.getByTestId('coupon-input');
    this.applyCouponButton = page.getByTestId('apply-coupon-button');
    this.removeCouponButton = page.getByTestId('remove-coupon-button');
    this.checkoutButton = page.getByTestId('checkout-button');
    this.continueShoppingButton = page.getByTestId('continue-shopping-button');
    this.clearCartButton = page.getByTestId('clear-cart-button');
  }

  /**
   * Navigate to cart page
   */
  async goto(): Promise<void> {
    await this.page.goto('/cart');
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Wait for cart to load
   */
  async waitForCartLoad(): Promise<void> {
    await this.cartContainer.waitFor({ state: 'visible' });
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Get cart item by product ID
   */
  getCartItem(productId: string): Locator {
    return this.cartItems.locator(`[data-product-id="${productId}"]`);
  }

  /**
   * Get cart item by index
   */
  getCartItemByIndex(index: number): Locator {
    return this.cartItems.locator('.cart-item').nth(index);
  }

  /**
   * Get quantity input for a cart item
   */
  getQuantityInput(productId: string): Locator {
    return this.getCartItem(productId).getByTestId('quantity-input');
  }

  /**
   * Get remove button for a cart item
   */
  getRemoveButton(productId: string): Locator {
    return this.getCartItem(productId).getByTestId('remove-item-button');
  }

  /**
   * Get item total for a cart item
   */
  getItemTotal(productId: string): Locator {
    return this.getCartItem(productId).getByTestId('item-total');
  }

  /**
   * Verify cart is empty
   */
  async expectCartEmpty(): Promise<void> {
    await expect(this.emptyCartMessage).toBeVisible();
    await expect(this.cartItems.locator('.cart-item')).toHaveCount(0);
    await expect(this.checkoutButton).not.toBeVisible();
  }

  /**
   * Verify cart has items
   */
  async expectCartHasItems(count?: number): Promise<void> {
    await expect(this.emptyCartMessage).not.toBeVisible();

    if (count !== undefined) {
      await expect(this.cartItems.locator('.cart-item')).toHaveCount(count);
    } else {
      await expect(this.cartItems.locator('.cart-item').first()).toBeVisible();
    }

    await expect(this.checkoutButton).toBeVisible();
  }

  /**
   * Verify cart item exists
   */
  async expectCartItemExists(productId: string, quantity?: number): Promise<void> {
    const item = this.getCartItem(productId);
    await expect(item).toBeVisible();

    if (quantity !== undefined) {
      const quantityInput = this.getQuantityInput(productId);
      await expect(quantityInput).toHaveValue(quantity.toString());
    }
  }

  /**
   * Verify cart totals
   */
  async expectCartTotals(expectedTotals: {
    subtotal?: string;
    tax?: string;
    shipping?: string;
    discount?: string;
    total: string;
  }): Promise<void> {
    if (expectedTotals.subtotal) {
      await expect(this.subtotal).toContainText(expectedTotals.subtotal);
    }

    if (expectedTotals.tax) {
      await expect(this.tax).toContainText(expectedTotals.tax);
    }

    if (expectedTotals.shipping) {
      await expect(this.shipping).toContainText(expectedTotals.shipping);
    }

    if (expectedTotals.discount) {
      await expect(this.discount).toContainText(expectedTotals.discount);
    }

    await expect(this.cartTotal).toContainText(expectedTotals.total);
  }

  /**
   * Update item quantity
   */
  async updateQuantity(productId: string, quantity: number): Promise<void> {
    const quantityInput = this.getQuantityInput(productId);
    await quantityInput.fill(quantity.toString());
    await quantityInput.press('Enter');

    // Wait for cart to update
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Remove item from cart
   */
  async removeItem(productId: string): Promise<void> {
    const removeButton = this.getRemoveButton(productId);
    await removeButton.click();

    // Wait for cart to update
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Apply coupon code
   */
  async applyCoupon(couponCode: string): Promise<void> {
    await this.couponInput.fill(couponCode);
    await this.applyCouponButton.click();

    // Wait for cart to update
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Remove applied coupon
   */
  async removeCoupon(): Promise<void> {
    await this.removeCouponButton.click();

    // Wait for cart to update
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Verify coupon is applied
   */
  async expectCouponApplied(couponCode: string): Promise<void> {
    const couponIndicator = this.page.getByTestId('applied-coupon');
    await expect(couponIndicator).toBeVisible();
    await expect(couponIndicator).toContainText(couponCode);
    await expect(this.discount).toBeVisible();
    await expect(this.removeCouponButton).toBeVisible();
  }

  /**
   * Verify coupon error
   */
  async expectCouponError(errorMessage: string): Promise<void> {
    const errorElement = this.page.getByTestId('coupon-error');
    await expect(errorElement).toBeVisible();
    await expect(errorElement).toContainText(errorMessage);
  }

  /**
   * Clear entire cart
   */
  async clearCart(): Promise<void> {
    if (await this.clearCartButton.isVisible()) {
      await this.clearCartButton.click();

      // Confirm clear cart if modal appears
      const confirmButton = this.page.getByTestId('confirm-clear-cart');
      if (await confirmButton.isVisible({ timeout: 1000 })) {
        await confirmButton.click();
      }

      await this.page.waitForLoadState('networkidle');
    }
  }

  /**
   * Proceed to checkout
   */
  async proceedToCheckout(): Promise<void> {
    await expect(this.checkoutButton).toBeEnabled();
    await this.checkoutButton.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Continue shopping
   */
  async continueShopping(): Promise<void> {
    await this.continueShoppingButton.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Get cart summary for verification
   */
  async getCartSummary(): Promise<{
    itemCount: number;
    subtotal: string;
    total: string;
    hasCoupon: boolean;
  }> {
    const itemCount = await this.cartItems.locator('.cart-item').count();
    const subtotal = await this.subtotal.textContent() || '';
    const total = await this.cartTotal.textContent() || '';
    const hasCoupon = await this.removeCouponButton.isVisible();

    return {
      itemCount,
      subtotal,
      total,
      hasCoupon,
    };
  }

  /**
   * Verify cart item details
   */
  async expectCartItemDetails(productId: string, details: {
    name?: string;
    price?: string;
    quantity?: number;
    total?: string;
    variant?: string;
  }): Promise<void> {
    const item = this.getCartItem(productId);

    if (details.name) {
      const nameElement = item.getByTestId('item-name');
      await expect(nameElement).toContainText(details.name);
    }

    if (details.price) {
      const priceElement = item.getByTestId('item-price');
      await expect(priceElement).toContainText(details.price);
    }

    if (details.quantity !== undefined) {
      const quantityInput = this.getQuantityInput(productId);
      await expect(quantityInput).toHaveValue(details.quantity.toString());
    }

    if (details.total) {
      const totalElement = this.getItemTotal(productId);
      await expect(totalElement).toContainText(details.total);
    }

    if (details.variant) {
      const variantElement = item.getByTestId('item-variant');
      await expect(variantElement).toContainText(details.variant);
    }
  }

  /**
   * Test cart persistence across page reloads
   */
  async testCartPersistence(): Promise<void> {
    const cartSummaryBefore = await this.getCartSummary();

    await this.page.reload();
    await this.waitForCartLoad();

    const cartSummaryAfter = await this.getCartSummary();

    expect(cartSummaryAfter.itemCount).toBe(cartSummaryBefore.itemCount);
    expect(cartSummaryAfter.total).toBe(cartSummaryBefore.total);
  }

  /**
   * Test cart responsiveness
   */
  async testResponsiveLayout(): Promise<void> {
    // Test mobile layout
    await this.page.setViewportSize({ width: 375, height: 667 });
    await this.waitForCartLoad();
    await expect(this.cartContainer).toBeVisible();

    // Test tablet layout
    await this.page.setViewportSize({ width: 768, height: 1024 });
    await this.waitForCartLoad();
    await expect(this.cartContainer).toBeVisible();

    // Test desktop layout
    await this.page.setViewportSize({ width: 1280, height: 720 });
    await this.waitForCartLoad();
    await expect(this.cartContainer).toBeVisible();
  }

  /**
   * Verify cart loading states
   */
  async expectLoadingState(): Promise<void> {
    const loadingIndicator = this.page.getByTestId('cart-loading');
    await expect(loadingIndicator).toBeVisible();
  }

  /**
   * Verify cart loaded state
   */
  async expectLoadedState(): Promise<void> {
    const loadingIndicator = this.page.getByTestId('cart-loading');
    await expect(loadingIndicator).not.toBeVisible();
    await expect(this.cartContainer).toBeVisible();
  }
}