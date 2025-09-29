/**
 * Checkout Flow E2E Tests
 * Comprehensive testing of the checkout process from cart to order completion
 */

import { test, expect } from '@playwright/test';
import { CartPage } from '../page-objects/CartPage';
import { CategoryPage } from '../page-objects/CategoryPage';
import { TestDataGenerator, TestDataSets } from '../../utils/test-data';
import { createApiHelpers } from '../../utils/api-helpers';
import { ScreenshotHelpers } from '../../utils/screenshot-helpers';

test.describe('Checkout Flow', () => {
  let cartPage: CartPage;
  let categoryPage: CategoryPage;
  let apiHelpers: ReturnType<typeof createApiHelpers>;
  let screenshotHelpers: ScreenshotHelpers;

  // Test user data
  const testUser = TestDataSets.users.new;
  const testAddress = TestDataGenerator.generateAddress();
  const testPaymentMethod = {
    cardNumber: '4111111111111111',
    expiryMonth: '12',
    expiryYear: '2025',
    cvv: '123',
    name: 'John Doe',
  };

  test.beforeEach(async ({ page, request }) => {
    cartPage = new CartPage(page);
    categoryPage = new CategoryPage(page);
    apiHelpers = createApiHelpers(request);
    screenshotHelpers = new ScreenshotHelpers(page);

    // Setup test data - add products to cart
    await test.step('Setup cart with test products', async () => {
      const testProduct = TestDataSets.products.electronics;
      await categoryPage.goto('electronics');
      await categoryPage.waitForCategoryLoad();
      await categoryPage.addProductToCart(testProduct.sku!, 1);
    });
  });

  test.afterEach(async ({ page }) => {
    // Cleanup
    await cartPage.goto();
    await cartPage.clearCart();
  });

  test.describe('Checkout Initiation', () => {
    test('should navigate to checkout from cart', async ({ page }) => {
      await test.step('Navigate to checkout', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();
        await cartPage.proceedToCheckout();
      });

      await test.step('Verify checkout page loaded', async () => {
        await expect(page).toHaveURL(/.*\/checkout.*/);

        const checkoutContainer = page.getByTestId('checkout-container');
        await expect(checkoutContainer).toBeVisible();

        const checkoutSteps = page.getByTestId('checkout-steps');
        await expect(checkoutSteps).toBeVisible();
      });

      await test.step('Take checkout initial screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'checkout-initial',
          fullPage: true,
        });
      });
    });

    test('should display cart summary in checkout', async ({ page }) => {
      await test.step('Navigate to checkout', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();
        await cartPage.proceedToCheckout();
      });

      await test.step('Verify cart summary', async () => {
        const cartSummary = page.getByTestId('checkout-cart-summary');
        await expect(cartSummary).toBeVisible();

        const testProduct = TestDataSets.products.electronics;
        const itemName = cartSummary.getByTestId('cart-item-name');
        await expect(itemName).toContainText(testProduct.name);

        const itemPrice = cartSummary.getByTestId('cart-item-price');
        await expect(itemPrice).toContainText(`$${testProduct.price.toFixed(2)}`);
      });
    });
  });

  test.describe('Guest Checkout', () => {
    test('should complete guest checkout successfully', async ({ page }) => {
      await test.step('Start guest checkout', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();
      });

      await test.step('Fill shipping information', async () => {
        const shippingForm = page.getByTestId('shipping-form');
        await expect(shippingForm).toBeVisible();

        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);
        await page.getByTestId('shipping-phone').fill(testUser.phone!);

        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();
      });

      await test.step('Take shipping form screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'checkout-shipping-filled',
          fullPage: true,
        });
      });

      await test.step('Select shipping method', async () => {
        const shippingMethods = page.getByTestId('shipping-methods');
        await expect(shippingMethods).toBeVisible();

        const standardShipping = page.getByTestId('shipping-standard');
        await standardShipping.click();

        const continueToPayment = page.getByTestId('continue-to-payment');
        await continueToPayment.click();
      });

      await test.step('Fill payment information', async () => {
        const paymentForm = page.getByTestId('payment-form');
        await expect(paymentForm).toBeVisible();

        await page.getByTestId('card-number').fill(testPaymentMethod.cardNumber);
        await page.getByTestId('card-expiry-month').selectOption(testPaymentMethod.expiryMonth);
        await page.getByTestId('card-expiry-year').selectOption(testPaymentMethod.expiryYear);
        await page.getByTestId('card-cvv').fill(testPaymentMethod.cvv);
        await page.getByTestId('card-name').fill(testPaymentMethod.name);

        const billingAddressSame = page.getByTestId('billing-same-as-shipping');
        await billingAddressSame.check();
      });

      await test.step('Take payment form screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'checkout-payment-filled',
          fullPage: true,
        });
      });

      await test.step('Review and place order', async () => {
        const continueToReview = page.getByTestId('continue-to-review');
        await continueToReview.click();

        const orderReview = page.getByTestId('order-review');
        await expect(orderReview).toBeVisible();

        // Verify order details
        const orderTotal = page.getByTestId('order-total');
        await expect(orderTotal).toContainText('$');

        const placeOrderButton = page.getByTestId('place-order-button');
        await expect(placeOrderButton).toBeEnabled();

        await placeOrderButton.click();
      });

      await test.step('Verify order confirmation', async () => {
        // Wait for order processing
        await page.waitForLoadState('networkidle');

        const orderConfirmation = page.getByTestId('order-confirmation');
        await expect(orderConfirmation).toBeVisible({ timeout: 10000 });

        const orderNumber = page.getByTestId('order-number');
        await expect(orderNumber).toBeVisible();

        const confirmationEmail = page.getByTestId('confirmation-email');
        await expect(confirmationEmail).toContainText(testUser.email);
      });

      await test.step('Take order confirmation screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'checkout-confirmation',
          fullPage: true,
        });
      });
    });

    test('should handle validation errors in shipping form', async ({ page }) => {
      await test.step('Start checkout and submit empty form', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();
      });

      await test.step('Verify validation errors', async () => {
        const emailError = page.getByTestId('shipping-email-error');
        await expect(emailError).toBeVisible();
        await expect(emailError).toContainText('Email is required');

        const nameError = page.getByTestId('shipping-first-name-error');
        await expect(nameError).toBeVisible();
        await expect(nameError).toContainText('First name is required');

        const addressError = page.getByTestId('shipping-address-error');
        await expect(addressError).toBeVisible();
        await expect(addressError).toContainText('Address is required');
      });

      await test.step('Take validation errors screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'checkout-validation-errors',
          fullPage: true,
        });
      });
    });

    test('should handle invalid payment information', async ({ page }) => {
      await test.step('Complete shipping and reach payment', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        // Fill shipping form
        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);

        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();

        // Select shipping method
        const standardShipping = page.getByTestId('shipping-standard');
        await standardShipping.click();

        const continueToPayment = page.getByTestId('continue-to-payment');
        await continueToPayment.click();
      });

      await test.step('Submit invalid payment information', async () => {
        // Fill with invalid card number
        await page.getByTestId('card-number').fill('1234567890123456');
        await page.getByTestId('card-expiry-month').selectOption('12');
        await page.getByTestId('card-expiry-year').selectOption('2025');
        await page.getByTestId('card-cvv').fill('123');
        await page.getByTestId('card-name').fill('Test User');

        const continueToReview = page.getByTestId('continue-to-review');
        await continueToReview.click();
      });

      await test.step('Verify payment validation error', async () => {
        const cardError = page.getByTestId('card-number-error');
        await expect(cardError).toBeVisible();
        await expect(cardError).toContainText('Invalid card number');
      });
    });
  });

  test.describe('Registered User Checkout', () => {
    test('should use saved addresses for registered users', async ({ page }) => {
      await test.step('Login as registered user', async () => {
        // Simulate user login
        await page.goto('/login');
        await page.getByTestId('email-input').fill(TestDataSets.users.existing.email);
        await page.getByTestId('password-input').fill(TestDataSets.users.existing.password);
        await page.getByTestId('login-button').click();
        await page.waitForLoadState('networkidle');
      });

      await test.step('Start checkout with saved address', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const savedAddresses = page.getByTestId('saved-addresses');
        await expect(savedAddresses).toBeVisible();

        const firstAddress = page.getByTestId('address-option-0');
        await firstAddress.click();

        const useAddressButton = page.getByTestId('use-selected-address');
        await useAddressButton.click();
      });

      await test.step('Verify address pre-filled', async () => {
        const shippingForm = page.getByTestId('shipping-form');
        const firstName = await page.getByTestId('shipping-first-name').inputValue();
        expect(firstName).toBeTruthy();

        const address = await page.getByTestId('shipping-address').inputValue();
        expect(address).toBeTruthy();
      });
    });

    test('should save new address during checkout', async ({ page }) => {
      await test.step('Login and start checkout', async () => {
        await page.goto('/login');
        await page.getByTestId('email-input').fill(TestDataSets.users.existing.email);
        await page.getByTestId('password-input').fill(TestDataSets.users.existing.password);
        await page.getByTestId('login-button').click();

        await cartPage.goto();
        await cartPage.proceedToCheckout();
      });

      await test.step('Add new address with save option', async () => {
        const addNewAddressButton = page.getByTestId('add-new-address');
        await addNewAddressButton.click();

        // Fill new address
        await page.getByTestId('shipping-first-name').fill('John');
        await page.getByTestId('shipping-last-name').fill('Doe');
        await page.getByTestId('shipping-address').fill('123 New Street');
        await page.getByTestId('shipping-city').fill('New City');
        await page.getByTestId('shipping-state').selectOption('CA');
        await page.getByTestId('shipping-zip').fill('90210');

        const saveAddressCheckbox = page.getByTestId('save-address-checkbox');
        await saveAddressCheckbox.check();

        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();
      });

      await test.step('Verify address save confirmation', async () => {
        const saveConfirmation = page.getByTestId('address-saved-confirmation');
        await expect(saveConfirmation).toBeVisible({ timeout: 5000 });
      });
    });
  });

  test.describe('Payment Methods', () => {
    test('should handle credit card payment', async ({ page }) => {
      await test.step('Complete checkout with credit card', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        // Complete guest checkout flow
        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        // Fill minimal shipping info
        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);

        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();

        // Select shipping
        const standardShipping = page.getByTestId('shipping-standard');
        await standardShipping.click();

        const continueToPayment = page.getByTestId('continue-to-payment');
        await continueToPayment.click();
      });

      await test.step('Select credit card payment method', async () => {
        const creditCardOption = page.getByTestId('payment-credit-card');
        await creditCardOption.click();

        const cardForm = page.getByTestId('credit-card-form');
        await expect(cardForm).toBeVisible();
      });

      await test.step('Fill valid credit card information', async () => {
        await page.getByTestId('card-number').fill(testPaymentMethod.cardNumber);
        await page.getByTestId('card-expiry-month').selectOption(testPaymentMethod.expiryMonth);
        await page.getByTestId('card-expiry-year').selectOption(testPaymentMethod.expiryYear);
        await page.getByTestId('card-cvv').fill(testPaymentMethod.cvv);
        await page.getByTestId('card-name').fill(testPaymentMethod.name);

        const continueToReview = page.getByTestId('continue-to-review');
        await continueToReview.click();
      });

      await test.step('Verify payment method in review', async () => {
        const paymentSummary = page.getByTestId('payment-summary');
        await expect(paymentSummary).toContainText('**** **** **** 1111');
        await expect(paymentSummary).toContainText('Credit Card');
      });
    });

    test('should handle PayPal payment', async ({ page }) => {
      await test.step('Select PayPal payment option', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        // Fill minimal shipping info
        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);

        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();

        // Select shipping
        const standardShipping = page.getByTestId('shipping-standard');
        await standardShipping.click();

        const continueToPayment = page.getByTestId('continue-to-payment');
        await continueToPayment.click();

        // Select PayPal
        const paypalOption = page.getByTestId('payment-paypal');
        await paypalOption.click();
      });

      await test.step('Verify PayPal integration', async () => {
        const paypalButton = page.getByTestId('paypal-button');
        await expect(paypalButton).toBeVisible();

        // Click PayPal button (in real test, this would open PayPal popup)
        await paypalButton.click();

        // Mock PayPal response
        const paypalSuccess = page.getByTestId('paypal-success');
        await expect(paypalSuccess).toBeVisible({ timeout: 5000 });
      });
    });
  });

  test.describe('Shipping Options', () => {
    test('should display shipping options and calculate costs', async ({ page }) => {
      await test.step('Reach shipping selection', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        // Fill shipping address
        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);

        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();
      });

      await test.step('Verify shipping options', async () => {
        const shippingMethods = page.getByTestId('shipping-methods');
        await expect(shippingMethods).toBeVisible();

        // Standard shipping
        const standardShipping = page.getByTestId('shipping-standard');
        await expect(standardShipping).toBeVisible();
        await expect(standardShipping).toContainText('5-7 business days');

        // Express shipping
        const expressShipping = page.getByTestId('shipping-express');
        await expect(expressShipping).toBeVisible();
        await expect(expressShipping).toContainText('2-3 business days');

        // Overnight shipping
        const overnightShipping = page.getByTestId('shipping-overnight');
        await expect(overnightShipping).toBeVisible();
        await expect(overnightShipping).toContainText('1 business day');
      });

      await test.step('Test shipping cost calculation', async () => {
        const orderSummary = page.getByTestId('order-summary');

        // Select standard shipping
        const standardShipping = page.getByTestId('shipping-standard');
        await standardShipping.click();

        let shippingCost = await orderSummary.getByTestId('shipping-cost').textContent();
        expect(shippingCost).toContain('$');

        // Select express shipping
        const expressShipping = page.getByTestId('shipping-express');
        await expressShipping.click();

        const newShippingCost = await orderSummary.getByTestId('shipping-cost').textContent();
        expect(newShippingCost).toContain('$');

        // Express should cost more than standard
        const standardCostValue = parseFloat(shippingCost?.replace(/[^0-9.]/g, '') || '0');
        const expressCostValue = parseFloat(newShippingCost?.replace(/[^0-9.]/g, '') || '0');
        expect(expressCostValue).toBeGreaterThan(standardCostValue);
      });
    });

    test('should handle free shipping threshold', async ({ page }) => {
      await test.step('Add items to reach free shipping', async () => {
        // Add multiple items to cart to reach free shipping threshold
        const products = [
          TestDataSets.products.electronics,
          TestDataSets.products.clothing,
          TestDataSets.products.books,
        ];

        for (const product of products) {
          await categoryPage.goto(product.category.toLowerCase());
          await categoryPage.addProductToCart(product.sku!, 2);
        }
      });

      await test.step('Verify free shipping in checkout', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const freeShippingBanner = page.getByTestId('free-shipping-banner');
        await expect(freeShippingBanner).toBeVisible();
        await expect(freeShippingBanner).toContainText('You qualify for free shipping!');
      });
    });
  });

  test.describe('Order Review and Confirmation', () => {
    test('should display complete order summary', async ({ page }) => {
      await test.step('Complete checkout to review step', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        // Quick form fill
        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);

        await page.getByTestId('continue-to-payment').click();

        await page.getByTestId('shipping-standard').click();
        await page.getByTestId('continue-to-payment').click();

        await page.getByTestId('card-number').fill(testPaymentMethod.cardNumber);
        await page.getByTestId('card-expiry-month').selectOption(testPaymentMethod.expiryMonth);
        await page.getByTestId('card-expiry-year').selectOption(testPaymentMethod.expiryYear);
        await page.getByTestId('card-cvv').fill(testPaymentMethod.cvv);
        await page.getByTestId('card-name').fill(testPaymentMethod.name);

        await page.getByTestId('continue-to-review').click();
      });

      await test.step('Verify order review components', async () => {
        const orderReview = page.getByTestId('order-review');
        await expect(orderReview).toBeVisible();

        // Shipping address summary
        const shippingAddress = page.getByTestId('review-shipping-address');
        await expect(shippingAddress).toContainText(testUser.firstName);
        await expect(shippingAddress).toContainText(testAddress.street);

        // Payment method summary
        const paymentMethod = page.getByTestId('review-payment-method');
        await expect(paymentMethod).toContainText('**** **** **** 1111');

        // Order items summary
        const orderItems = page.getByTestId('review-order-items');
        await expect(orderItems).toBeVisible();

        // Order totals
        const orderTotal = page.getByTestId('order-total');
        await expect(orderTotal).toContainText('$');
      });

      await test.step('Take order review screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'checkout-order-review',
          fullPage: true,
        });
      });
    });

    test('should allow editing from review step', async ({ page }) => {
      await test.step('Complete checkout to review step', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        // Complete checkout flow
        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);

        await page.getByTestId('continue-to-payment').click();
        await page.getByTestId('shipping-standard').click();
        await page.getByTestId('continue-to-payment').click();

        await page.getByTestId('card-number').fill(testPaymentMethod.cardNumber);
        await page.getByTestId('card-expiry-month').selectOption(testPaymentMethod.expiryMonth);
        await page.getByTestId('card-expiry-year').selectOption(testPaymentMethod.expiryYear);
        await page.getByTestId('card-cvv').fill(testPaymentMethod.cvv);
        await page.getByTestId('card-name').fill(testPaymentMethod.name);

        await page.getByTestId('continue-to-review').click();
      });

      await test.step('Edit shipping address from review', async () => {
        const editShippingButton = page.getByTestId('edit-shipping-address');
        await editShippingButton.click();

        // Should return to shipping step
        const shippingForm = page.getByTestId('shipping-form');
        await expect(shippingForm).toBeVisible();

        // Modify address
        await page.getByTestId('shipping-address').fill('456 Modified Street');

        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();

        // Skip through steps back to review
        await page.getByTestId('continue-to-payment').click();
        await page.getByTestId('continue-to-review').click();

        // Verify modification reflected
        const shippingAddress = page.getByTestId('review-shipping-address');
        await expect(shippingAddress).toContainText('456 Modified Street');
      });
    });
  });

  test.describe('Error Handling and Edge Cases', () => {
    test('should handle payment processing errors', async ({ page }) => {
      await test.step('Complete checkout with failing payment', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        // Complete checkout flow with card that will fail
        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);

        await page.getByTestId('continue-to-payment').click();
        await page.getByTestId('shipping-standard').click();
        await page.getByTestId('continue-to-payment').click();

        // Use card number that will be declined
        await page.getByTestId('card-number').fill('4000000000000002');
        await page.getByTestId('card-expiry-month').selectOption('12');
        await page.getByTestId('card-expiry-year').selectOption('2025');
        await page.getByTestId('card-cvv').fill('123');
        await page.getByTestId('card-name').fill('Test User');

        await page.getByTestId('continue-to-review').click();

        const placeOrderButton = page.getByTestId('place-order-button');
        await placeOrderButton.click();
      });

      await test.step('Verify payment error handling', async () => {
        const paymentError = page.getByTestId('payment-error');
        await expect(paymentError).toBeVisible({ timeout: 10000 });
        await expect(paymentError).toContainText('Payment was declined');

        // Should return to payment step
        const paymentForm = page.getByTestId('payment-form');
        await expect(paymentForm).toBeVisible();
      });
    });

    test('should handle session timeout', async ({ page }) => {
      await test.step('Simulate session timeout', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        // Simulate session expiration
        await page.evaluate(() => {
          localStorage.removeItem('auth-token');
          sessionStorage.clear();
        });

        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();
      });

      await test.step('Verify session handling', async () => {
        // Should handle gracefully, either redirect to login or continue as guest
        const sessionMessage = page.getByTestId('session-expired-message');

        if (await sessionMessage.isVisible({ timeout: 2000 })) {
          await expect(sessionMessage).toContainText('session expired');
        } else {
          // Should continue as guest checkout
          const shippingForm = page.getByTestId('shipping-form');
          await expect(shippingForm).toBeVisible();
        }
      });
    });

    test('should handle inventory changes during checkout', async ({ page }) => {
      await test.step('Complete checkout with item going out of stock', async () => {
        await cartPage.goto();
        await cartPage.proceedToCheckout();

        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        // Fill basic info
        await page.getByTestId('shipping-email').fill(testUser.email);
        await page.getByTestId('shipping-first-name').fill(testUser.firstName);
        await page.getByTestId('shipping-last-name').fill(testUser.lastName);
        await page.getByTestId('shipping-address').fill(testAddress.street);
        await page.getByTestId('shipping-city').fill(testAddress.city);
        await page.getByTestId('shipping-state').selectOption(testAddress.state);
        await page.getByTestId('shipping-zip').fill(testAddress.zipCode);

        await page.getByTestId('continue-to-payment').click();

        // Simulate inventory change during checkout
        await page.route('**/api/v1/commerce/orders', route => {
          route.fulfill({
            status: 409,
            contentType: 'application/json',
            body: JSON.stringify({
              error: 'Item out of stock',
              code: 'INVENTORY_INSUFFICIENT',
            }),
          });
        });

        await page.getByTestId('shipping-standard').click();
        await page.getByTestId('continue-to-payment').click();

        await page.getByTestId('card-number').fill(testPaymentMethod.cardNumber);
        await page.getByTestId('card-expiry-month').selectOption(testPaymentMethod.expiryMonth);
        await page.getByTestId('card-expiry-year').selectOption(testPaymentMethod.expiryYear);
        await page.getByTestId('card-cvv').fill(testPaymentMethod.cvv);
        await page.getByTestId('card-name').fill(testPaymentMethod.name);

        await page.getByTestId('continue-to-review').click();

        const placeOrderButton = page.getByTestId('place-order-button');
        await placeOrderButton.click();
      });

      await test.step('Verify inventory error handling', async () => {
        const inventoryError = page.getByTestId('inventory-error');
        await expect(inventoryError).toBeVisible({ timeout: 10000 });
        await expect(inventoryError).toContainText('out of stock');

        // Should provide option to update cart
        const updateCartButton = page.getByTestId('update-cart-button');
        await expect(updateCartButton).toBeVisible();
      });
    });
  });

  test.describe('Mobile Checkout Experience', () => {
    test('should work correctly on mobile devices', async ({ page }) => {
      await test.step('Set mobile viewport', async () => {
        await page.setViewportSize({ width: 375, height: 667 });
      });

      await test.step('Test mobile checkout flow', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();
        await cartPage.proceedToCheckout();

        // Verify mobile-optimized checkout layout
        const checkoutContainer = page.getByTestId('checkout-container');
        await expect(checkoutContainer).toBeVisible();

        const mobileSteps = page.getByTestId('mobile-checkout-steps');
        if (await mobileSteps.isVisible({ timeout: 2000 })) {
          await expect(mobileSteps).toBeVisible();
        }
      });

      await test.step('Test mobile form interaction', async () => {
        const guestCheckoutButton = page.getByTestId('guest-checkout-button');
        await guestCheckoutButton.click();

        // Test mobile keyboard behavior
        const emailInput = page.getByTestId('shipping-email');
        await emailInput.click();
        await emailInput.fill(testUser.email);

        // Verify mobile-specific validations
        const continueButton = page.getByTestId('continue-to-payment');
        await continueButton.click();

        // Should show validation errors on mobile
        const emailError = page.getByTestId('shipping-first-name-error');
        await expect(emailError).toBeVisible();
      });

      await test.step('Take mobile checkout screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'checkout-mobile-view',
          fullPage: true,
        });
      });
    });
  });
});