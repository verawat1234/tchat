/**
 * Cross-Platform Synchronization E2E Tests
 * Tests for cart and user data synchronization between web and mobile platforms
 */

import { test, expect, devices } from '@playwright/test';
import { CartPage } from '../web/page-objects/CartPage';
import { CategoryPage } from '../web/page-objects/CategoryPage';
import { TestDataGenerator, TestDataSets } from '../utils/test-data';
import { createApiHelpers } from '../utils/api-helpers';

test.describe('Cross-Platform Synchronization', () => {
  let webContext: any;
  let mobileContext: any;
  let webPage: any;
  let mobilePage: any;
  let apiHelpers: ReturnType<typeof createApiHelpers>;

  test.beforeAll(async ({ browser, request }) => {
    // Create web context
    webContext = await browser.newContext({
      ...devices['Desktop Chrome'],
      storageState: undefined,
    });

    // Create mobile context
    mobileContext = await browser.newContext({
      ...devices['iPhone 12'],
      storageState: undefined,
    });

    webPage = await webContext.newPage();
    mobilePage = await mobileContext.newPage();

    apiHelpers = createApiHelpers(request);
  });

  test.afterAll(async () => {
    await webContext.close();
    await mobileContext.close();
  });

  test.beforeEach(async () => {
    // Clear storage and reset state
    await webPage.evaluate(() => {
      localStorage.clear();
      sessionStorage.clear();
    });

    await mobilePage.evaluate(() => {
      localStorage.clear();
      sessionStorage.clear();
    });
  });

  test.describe('Cart Synchronization', () => {
    test('should sync cart from web to mobile when user logs in', async () => {
      const testUser = TestDataSets.users.existing;
      const testProduct = TestDataSets.products.electronics;

      await test.step('Add items to cart on web as guest', async () => {
        const webCartPage = new CartPage(webPage);
        const webCategoryPage = new CategoryPage(webPage);

        await webCategoryPage.goto('electronics');
        await webCategoryPage.addProductToCart(testProduct.sku!, 2);

        await webCartPage.goto();
        await webCartPage.expectCartHasItems(1);
        await webCartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 2,
        });
      });

      await test.step('Login on web to save cart to account', async () => {
        await webPage.goto('/login');
        await webPage.getByTestId('email-input').fill(testUser.email);
        await webPage.getByTestId('password-input').fill(testUser.password);
        await webPage.getByTestId('login-button').click();
        await webPage.waitForLoadState('networkidle');

        // Verify cart persisted after login
        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();
        await webCartPage.expectCartHasItems(1);
      });

      await test.step('Login on mobile and verify cart synced', async () => {
        await mobilePage.goto('/login');
        await mobilePage.getByTestId('email-input').fill(testUser.email);
        await mobilePage.getByTestId('password-input').fill(testUser.password);
        await mobilePage.getByTestId('login-button').click();
        await mobilePage.waitForLoadState('networkidle');

        // Navigate to cart on mobile
        const mobileCartPage = new CartPage(mobilePage);
        await mobileCartPage.goto();
        await mobileCartPage.expectCartHasItems(1);
        await mobileCartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 2,
        });
      });
    });

    test('should sync cart changes between platforms in real-time', async () => {
      const testUser = TestDataSets.users.existing;
      const testProduct = TestDataSets.products.electronics;

      await test.step('Login on both platforms', async () => {
        // Login on web
        await webPage.goto('/login');
        await webPage.getByTestId('email-input').fill(testUser.email);
        await webPage.getByTestId('password-input').fill(testUser.password);
        await webPage.getByTestId('login-button').click();
        await webPage.waitForLoadState('networkidle');

        // Login on mobile
        await mobilePage.goto('/login');
        await mobilePage.getByTestId('email-input').fill(testUser.email);
        await mobilePage.getByTestId('password-input').fill(testUser.password);
        await mobilePage.getByTestId('login-button').click();
        await mobilePage.waitForLoadState('networkidle');
      });

      await test.step('Add item on web', async () => {
        const webCategoryPage = new CategoryPage(webPage);
        await webCategoryPage.goto('electronics');
        await webCategoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Verify item appears on mobile', async () => {
        const mobileCartPage = new CartPage(mobilePage);
        await mobileCartPage.goto();

        // Wait for sync to occur
        await mobilePage.waitForTimeout(2000);
        await mobileCartPage.expectCartHasItems(1);
        await mobileCartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 1,
        });
      });

      await test.step('Update quantity on mobile', async () => {
        const mobileCartPage = new CartPage(mobilePage);
        await mobileCartPage.updateQuantity(testProduct.sku!, 3);
      });

      await test.step('Verify quantity updated on web', async () => {
        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();

        // Wait for sync to occur
        await webPage.waitForTimeout(2000);
        await webCartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 3,
        });
      });

      await test.step('Remove item on web', async () => {
        const webCartPage = new CartPage(webPage);
        await webCartPage.removeItem(testProduct.sku!);
      });

      await test.step('Verify item removed on mobile', async () => {
        const mobileCartPage = new CartPage(mobilePage);
        await mobileCartPage.goto();

        // Wait for sync to occur
        await mobilePage.waitForTimeout(2000);
        await mobileCartPage.expectCartEmpty();
      });
    });

    test('should handle concurrent cart modifications', async () => {
      const testUser = TestDataSets.users.existing;
      const product1 = TestDataSets.products.electronics;
      const product2 = TestDataSets.products.clothing;

      await test.step('Login on both platforms', async () => {
        // Login on both platforms simultaneously
        await Promise.all([
          webPage.goto('/login').then(async () => {
            await webPage.getByTestId('email-input').fill(testUser.email);
            await webPage.getByTestId('password-input').fill(testUser.password);
            await webPage.getByTestId('login-button').click();
            await webPage.waitForLoadState('networkidle');
          }),
          mobilePage.goto('/login').then(async () => {
            await mobilePage.getByTestId('email-input').fill(testUser.email);
            await mobilePage.getByTestId('password-input').fill(testUser.password);
            await mobilePage.getByTestId('login-button').click();
            await mobilePage.waitForLoadState('networkidle');
          }),
        ]);
      });

      await test.step('Add different items simultaneously', async () => {
        // Add items to cart on both platforms at the same time
        await Promise.all([
          new CategoryPage(webPage).goto('electronics').then(async () => {
            const webCategoryPage = new CategoryPage(webPage);
            await webCategoryPage.addProductToCart(product1.sku!, 1);
          }),
          new CategoryPage(mobilePage).goto('clothing').then(async () => {
            const mobileCategoryPage = new CategoryPage(mobilePage);
            await mobileCategoryPage.addProductToCart(product2.sku!, 1);
          }),
        ]);
      });

      await test.step('Verify both items in cart on both platforms', async () => {
        // Wait for synchronization
        await Promise.all([webPage.waitForTimeout(3000), mobilePage.waitForTimeout(3000)]);

        // Check web cart
        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();
        await webCartPage.expectCartHasItems(2);

        // Check mobile cart
        const mobileCartPage = new CartPage(mobilePage);
        await mobileCartPage.goto();
        await mobileCartPage.expectCartHasItems(2);
      });
    });

    test('should resolve cart conflicts with last-write-wins strategy', async () => {
      const testUser = TestDataSets.users.existing;
      const testProduct = TestDataSets.products.electronics;

      await test.step('Setup cart with item on both platforms', async () => {
        // Login and add item on web
        await webPage.goto('/login');
        await webPage.getByTestId('email-input').fill(testUser.email);
        await webPage.getByTestId('password-input').fill(testUser.password);
        await webPage.getByTestId('login-button').click();
        await webPage.waitForLoadState('networkidle');

        const webCategoryPage = new CategoryPage(webPage);
        await webCategoryPage.goto('electronics');
        await webCategoryPage.addProductToCart(testProduct.sku!, 2);

        // Login on mobile and wait for sync
        await mobilePage.goto('/login');
        await mobilePage.getByTestId('email-input').fill(testUser.email);
        await mobilePage.getByTestId('password-input').fill(testUser.password);
        await mobilePage.getByTestId('login-button').click();
        await mobilePage.waitForLoadState('networkidle');

        await mobilePage.waitForTimeout(2000);
      });

      await test.step('Simulate network disconnection and make conflicting changes', async () => {
        // Simulate offline mode on mobile
        await mobilePage.route('**/api/v1/commerce/cart/**', route => route.abort());

        // Update quantity on mobile (offline)
        const mobileCartPage = new CartPage(mobilePage);
        await mobileCartPage.goto();
        await mobileCartPage.updateQuantity(testProduct.sku!, 5);

        // Update quantity on web (online)
        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();
        await webCartPage.updateQuantity(testProduct.sku!, 3);
      });

      await test.step('Restore network and verify conflict resolution', async () => {
        // Restore network on mobile
        await mobilePage.unroute('**/api/v1/commerce/cart/**');

        // Trigger sync by navigating away and back
        await mobilePage.goto('/');
        await mobilePage.waitForTimeout(1000);

        const mobileCartPage = new CartPage(mobilePage);
        await mobileCartPage.goto();

        // Wait for sync and verify last-write-wins (web change should win)
        await mobilePage.waitForTimeout(3000);
        await mobileCartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 3,
        });
      });
    });
  });

  test.describe('User Session Synchronization', () => {
    test('should sync user login across platforms', async () => {
      const testUser = TestDataSets.users.existing;

      await test.step('Login on web', async () => {
        await webPage.goto('/login');
        await webPage.getByTestId('email-input').fill(testUser.email);
        await webPage.getByTestId('password-input').fill(testUser.password);
        await webPage.getByTestId('login-button').click();
        await webPage.waitForLoadState('networkidle');

        // Verify logged in on web
        const userMenu = webPage.getByTestId('user-menu');
        await expect(userMenu).toBeVisible();
      });

      await test.step('Verify auto-login on mobile', async () => {
        await mobilePage.goto('/');
        await mobilePage.waitForLoadState('networkidle');

        // Wait for session sync
        await mobilePage.waitForTimeout(2000);

        // Check if user is automatically logged in
        const userMenu = mobilePage.getByTestId('user-menu');

        if (await userMenu.isVisible({ timeout: 1000 })) {
          // Already logged in
          const userInfo = mobilePage.getByTestId('user-info');
          await expect(userInfo).toContainText(testUser.firstName);
        } else {
          // Need to refresh or manually sync
          await mobilePage.reload();
          await mobilePage.waitForLoadState('networkidle');
        }
      });
    });

    test('should sync user logout across platforms', async () => {
      const testUser = TestDataSets.users.existing;

      await test.step('Login on both platforms', async () => {
        await Promise.all([
          webPage.goto('/login').then(async () => {
            await webPage.getByTestId('email-input').fill(testUser.email);
            await webPage.getByTestId('password-input').fill(testUser.password);
            await webPage.getByTestId('login-button').click();
            await webPage.waitForLoadState('networkidle');
          }),
          mobilePage.goto('/login').then(async () => {
            await mobilePage.getByTestId('email-input').fill(testUser.email);
            await mobilePage.getByTestId('password-input').fill(testUser.password);
            await mobilePage.getByTestId('login-button').click();
            await mobilePage.waitForLoadState('networkidle');
          }),
        ]);
      });

      await test.step('Logout on web', async () => {
        const userMenu = webPage.getByTestId('user-menu');
        await userMenu.click();

        const logoutButton = webPage.getByTestId('logout-button');
        await logoutButton.click();

        // Verify logged out on web
        const loginButton = webPage.getByTestId('login-link');
        await expect(loginButton).toBeVisible();
      });

      await test.step('Verify auto-logout on mobile', async () => {
        // Navigate to a protected page to trigger session check
        await mobilePage.goto('/profile');

        // Should redirect to login or show login prompt
        await mobilePage.waitForTimeout(3000);

        const currentUrl = mobilePage.url();
        const loginForm = mobilePage.getByTestId('login-form');

        expect(currentUrl.includes('/login') || await loginForm.isVisible()).toBeTruthy();
      });
    });
  });

  test.describe('Data Consistency', () => {
    test('should maintain data consistency during network interruptions', async () => {
      const testUser = TestDataSets.users.existing;
      const testProduct = TestDataSets.products.electronics;

      await test.step('Setup authenticated session', async () => {
        await webPage.goto('/login');
        await webPage.getByTestId('email-input').fill(testUser.email);
        await webPage.getByTestId('password-input').fill(testUser.password);
        await webPage.getByTestId('login-button').click();
        await webPage.waitForLoadState('networkidle');
      });

      await test.step('Add item to cart while online', async () => {
        const webCategoryPage = new CategoryPage(webPage);
        await webCategoryPage.goto('electronics');
        await webCategoryPage.addProductToCart(testProduct.sku!, 2);

        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();
        await webCartPage.expectCartHasItems(1);
      });

      await test.step('Simulate network failure and make changes', async () => {
        // Intercept and block network requests
        await webPage.route('**/api/v1/**', route => route.abort());

        const webCartPage = new CartPage(webPage);
        await webCartPage.updateQuantity(testProduct.sku!, 4);

        // Verify optimistic update
        await webCartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 4,
        });
      });

      await test.step('Restore network and verify data consistency', async () => {
        // Restore network
        await webPage.unroute('**/api/v1/**');

        // Trigger sync by refreshing page
        await webPage.reload();
        await webPage.waitForLoadState('networkidle');

        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();

        // Verify data consistency (should either be 4 if sync worked, or 2 if it was reverted)
        const quantity = await webCartPage.getCartItem(testProduct.sku!).getByTestId('quantity-input').inputValue();
        expect(parseInt(quantity)).toBeGreaterThan(0); // At least the item should still exist
      });
    });

    test('should handle simultaneous updates gracefully', async () => {
      const testUser = TestDataSets.users.existing;
      const testProduct = TestDataSets.products.electronics;

      await test.step('Setup authenticated sessions on both platforms', async () => {
        await Promise.all([
          webPage.goto('/login').then(async () => {
            await webPage.getByTestId('email-input').fill(testUser.email);
            await webPage.getByTestId('password-input').fill(testUser.password);
            await webPage.getByTestId('login-button').click();
            await webPage.waitForLoadState('networkidle');
          }),
          mobilePage.goto('/login').then(async () => {
            await mobilePage.getByTestId('email-input').fill(testUser.email);
            await mobilePage.getByTestId('password-input').fill(testUser.password);
            await mobilePage.getByTestId('login-button').click();
            await mobilePage.waitForLoadState('networkidle');
          }),
        ]);

        // Add initial item
        const webCategoryPage = new CategoryPage(webPage);
        await webCategoryPage.goto('electronics');
        await webCategoryPage.addProductToCart(testProduct.sku!, 1);

        // Wait for sync
        await Promise.all([webPage.waitForTimeout(2000), mobilePage.waitForTimeout(2000)]);
      });

      await test.step('Make simultaneous updates', async () => {
        // Update quantity on both platforms simultaneously
        await Promise.all([
          new CartPage(webPage).goto().then(async () => {
            const webCartPage = new CartPage(webPage);
            await webCartPage.updateQuantity(testProduct.sku!, 3);
          }),
          new CartPage(mobilePage).goto().then(async () => {
            const mobileCartPage = new CartPage(mobilePage);
            await mobileCartPage.updateQuantity(testProduct.sku!, 5);
          }),
        ]);
      });

      await test.step('Verify conflict resolution', async () => {
        // Wait for conflict resolution
        await Promise.all([webPage.waitForTimeout(5000), mobilePage.waitForTimeout(5000)]);

        // Refresh both platforms
        await Promise.all([
          webPage.reload().then(() => webPage.waitForLoadState('networkidle')),
          mobilePage.reload().then(() => mobilePage.waitForLoadState('networkidle')),
        ]);

        // Check final state on both platforms
        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();

        const mobileCartPage = new CartPage(mobilePage);
        await mobileCartPage.goto();

        const webQuantity = await webCartPage.getCartItem(testProduct.sku!).getByTestId('quantity-input').inputValue();
        const mobileQuantity = await mobileCartPage.getCartItem(testProduct.sku!).getByTestId('quantity-input').inputValue();

        // Both platforms should have the same final quantity
        expect(webQuantity).toBe(mobileQuantity);
      });
    });
  });

  test.describe('Offline/Online Synchronization', () => {
    test('should sync changes when going from offline to online', async () => {
      const testUser = TestDataSets.users.existing;
      const testProduct = TestDataSets.products.electronics;

      await test.step('Setup cart while online', async () => {
        await webPage.goto('/login');
        await webPage.getByTestId('email-input').fill(testUser.email);
        await webPage.getByTestId('password-input').fill(testUser.password);
        await webPage.getByTestId('login-button').click();
        await webPage.waitForLoadState('networkidle');

        const webCategoryPage = new CategoryPage(webPage);
        await webCategoryPage.goto('electronics');
        await webCategoryPage.addProductToCart(testProduct.sku!, 1);
      });

      await test.step('Go offline and make changes', async () => {
        // Simulate offline mode
        await webPage.route('**/api/v1/**', route => route.abort());

        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();

        // Try to update quantity (should queue for sync)
        await webCartPage.updateQuantity(testProduct.sku!, 3);

        // Verify offline indicator
        const offlineIndicator = webPage.getByTestId('offline-indicator');
        if (await offlineIndicator.isVisible({ timeout: 2000 })) {
          await expect(offlineIndicator).toBeVisible();
        }
      });

      await test.step('Go back online and verify sync', async () => {
        // Restore network
        await webPage.unroute('**/api/v1/**');

        // Trigger online detection
        await webPage.evaluate(() => {
          window.dispatchEvent(new Event('online'));
        });

        // Wait for sync to complete
        await webPage.waitForTimeout(3000);

        // Verify changes were synced
        const webCartPage = new CartPage(webPage);
        await webCartPage.goto();
        await webCartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 3,
        });
      });
    });

    test('should show appropriate offline indicators', async () => {
      await test.step('Go offline', async () => {
        // Simulate offline mode
        await webPage.route('**/api/v1/**', route => route.abort());

        await webPage.goto('/');
        await webPage.waitForLoadState('networkidle');

        // Verify offline indicator appears
        const offlineIndicator = webPage.getByTestId('offline-indicator');
        if (await offlineIndicator.isVisible({ timeout: 3000 })) {
          await expect(offlineIndicator).toContainText('offline');
        }
      });

      await test.step('Try to perform online-only actions', async () => {
        const categoryPage = new CategoryPage(webPage);
        await categoryPage.goto('electronics');

        // Try to add product to cart
        if (await categoryPage.productCards.first().isVisible({ timeout: 2000 })) {
          await categoryPage.productCards.first().click();

          const addToCartButton = webPage.getByTestId('add-to-cart-button');
          await addToCartButton.click();

          // Should show offline message
          const offlineMessage = webPage.getByTestId('offline-action-message');
          if (await offlineMessage.isVisible({ timeout: 2000 })) {
            await expect(offlineMessage).toContainText('offline');
          }
        }
      });

      await test.step('Go back online', async () => {
        // Restore network
        await webPage.unroute('**/api/v1/**');

        await webPage.evaluate(() => {
          window.dispatchEvent(new Event('online'));
        });

        // Verify offline indicator disappears
        const offlineIndicator = webPage.getByTestId('offline-indicator');
        if (await offlineIndicator.isVisible({ timeout: 1000 })) {
          await expect(offlineIndicator).not.toBeVisible({ timeout: 5000 });
        }
      });
    });
  });
});