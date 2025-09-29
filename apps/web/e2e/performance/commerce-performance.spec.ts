/**
 * Commerce Performance E2E Tests
 * Tests for performance benchmarks, Core Web Vitals, and load testing
 */

import { test, expect } from '@playwright/test';
import { CartPage } from '../web/page-objects/CartPage';
import { CategoryPage } from '../web/page-objects/CategoryPage';
import { ProductPage } from '../web/page-objects/ProductPage';
import { TestDataSets } from '../utils/test-data';

test.describe('Commerce Performance Tests', () => {
  let cartPage: CartPage;
  let categoryPage: CategoryPage;
  let productPage: ProductPage;

  test.beforeEach(async ({ page }) => {
    cartPage = new CartPage(page);
    categoryPage = new CategoryPage(page);
    productPage = new ProductPage(page);
  });

  test.describe('Page Load Performance', () => {
    test('should load category page within performance budget', async ({ page }) => {
      const startTime = Date.now();

      await test.step('Navigate to category page', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();
      });

      const loadTime = Date.now() - startTime;

      await test.step('Verify performance metrics', async () => {
        // Page should load within 3 seconds
        expect(loadTime).toBeLessThan(3000);

        // Verify products are visible
        await categoryPage.expectProductsDisplayed();

        // Check Core Web Vitals
        const webVitals = await page.evaluate(() => {
          return new Promise((resolve) => {
            new PerformanceObserver((list) => {
              const entries = list.getEntries();
              const vitals: Record<string, number> = {};

              entries.forEach((entry) => {
                if (entry.name === 'largest-contentful-paint') {
                  vitals.lcp = entry.startTime;
                }
                if (entry.name === 'first-input-delay') {
                  vitals.fid = (entry as any).processingStart - entry.startTime;
                }
                if (entry.name === 'cumulative-layout-shift') {
                  vitals.cls = (entry as any).value;
                }
              });

              if (Object.keys(vitals).length > 0) {
                resolve(vitals);
              }
            }).observe({ entryTypes: ['largest-contentful-paint', 'first-input', 'layout-shift'] });

            // Fallback timeout
            setTimeout(() => resolve({}), 5000);
          });
        });

        console.log('Core Web Vitals:', webVitals);

        // LCP should be less than 2.5 seconds
        if ((webVitals as any).lcp) {
          expect((webVitals as any).lcp).toBeLessThan(2500);
        }

        // FID should be less than 100ms
        if ((webVitals as any).fid) {
          expect((webVitals as any).fid).toBeLessThan(100);
        }

        // CLS should be less than 0.1
        if ((webVitals as any).cls) {
          expect((webVitals as any).cls).toBeLessThan(0.1);
        }
      });
    });

    test('should load product detail page quickly', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Navigate to product page and measure load time', async () => {
        const startTime = performance.now();

        await productPage.goto(testProduct.sku!);
        await productPage.waitForProductLoad();

        const loadTime = performance.now() - startTime;

        // Product page should load within 2 seconds
        expect(loadTime).toBeLessThan(2000);
      });

      await test.step('Verify page content loaded', async () => {
        await productPage.expectProductDetails({
          name: testProduct.name,
        });

        // Verify images loaded
        await productPage.expectProductImages(1);
      });
    });

    test('should handle large product catalogs efficiently', async ({ page }) => {
      await test.step('Load category with many products', async () => {
        // Simulate category with 100+ products
        await page.route('**/api/v1/commerce/products**', async route => {
          const response = await route.fetch();
          const data = await response.json();

          // Expand product list to simulate large catalog
          const expandedProducts = [];
          for (let i = 0; i < 50; i++) {
            expandedProducts.push(...data.products);
          }

          await route.fulfill({
            response,
            json: {
              ...data,
              products: expandedProducts.slice(0, 100),
              total: 100,
            },
          });
        });

        const startTime = performance.now();

        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        const loadTime = performance.now() - startTime;

        // Should still load within reasonable time even with many products
        expect(loadTime).toBeLessThan(5000);
      });

      await test.step('Test scroll performance', async () => {
        const startTime = performance.now();

        // Test scroll performance
        await page.evaluate(() => {
          const productGrid = document.querySelector('[data-testid="product-grid"]');
          if (productGrid) {
            for (let i = 0; i < 10; i++) {
              productGrid.scrollTop += 300;
            }
          }
        });

        const scrollTime = performance.now() - startTime;

        // Scrolling should be smooth (less than 100ms for 10 scrolls)
        expect(scrollTime).toBeLessThan(100);
      });
    });
  });

  test.describe('Cart Performance', () => {
    test('should handle cart operations efficiently', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Add item to cart and measure performance', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        const startTime = performance.now();

        await categoryPage.addProductToCart(testProduct.sku!, 1);

        const addToCartTime = performance.now() - startTime;

        // Add to cart should be fast (less than 1 second)
        expect(addToCartTime).toBeLessThan(1000);
      });

      await test.step('Test cart update performance', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        const updateStartTime = performance.now();

        await cartPage.updateQuantity(testProduct.sku!, 5);

        const updateTime = performance.now() - updateStartTime;

        // Cart update should be fast (less than 500ms)
        expect(updateTime).toBeLessThan(500);
      });

      await test.step('Test cart load performance with multiple items', async () => {
        // Add more items to cart
        const products = [
          TestDataSets.products.clothing,
          TestDataSets.products.books,
        ];

        for (const product of products) {
          await categoryPage.goto(product.category.toLowerCase());
          await categoryPage.addProductToCart(product.sku!, 1);
        }

        const loadStartTime = performance.now();

        await cartPage.goto();
        await cartPage.waitForCartLoad();

        const loadTime = performance.now() - loadStartTime;

        // Cart with multiple items should load quickly
        expect(loadTime).toBeLessThan(1500);

        await cartPage.expectCartHasItems(3);
      });
    });

    test('should handle rapid cart operations', async ({ page }) => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Setup cart with item', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.addProductToCart(testProduct.sku!, 1);
        await cartPage.goto();
        await cartPage.waitForCartLoad();
      });

      await test.step('Perform rapid quantity updates', async () => {
        const startTime = performance.now();

        // Rapidly update quantity multiple times
        for (let i = 2; i <= 10; i++) {
          await cartPage.updateQuantity(testProduct.sku!, i);
          await page.waitForTimeout(100); // Small delay to allow update
        }

        const totalTime = performance.now() - startTime;

        // All updates should complete reasonably quickly
        expect(totalTime).toBeLessThan(5000);

        // Verify final quantity
        await cartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 10,
        });
      });
    });
  });

  test.describe('Search Performance', () => {
    test('should perform search quickly', async ({ page }) => {
      await test.step('Perform product search and measure time', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        const startTime = performance.now();

        await categoryPage.searchProducts('smartphone');

        const searchTime = performance.now() - startTime;

        // Search should complete within 1 second
        expect(searchTime).toBeLessThan(1000);

        await categoryPage.expectProductsDisplayed();
      });
    });

    test('should handle complex searches efficiently', async ({ page }) => {
      await test.step('Perform search with filters', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        const startTime = performance.now();

        // Apply multiple filters
        await categoryPage.searchProducts('phone');
        await categoryPage.applyPriceFilter(100, 500);
        await categoryPage.applyRatingFilter(4);

        const searchTime = performance.now() - startTime;

        // Complex search should complete within 2 seconds
        expect(searchTime).toBeLessThan(2000);
      });
    });
  });

  test.describe('Image Loading Performance', () => {
    test('should load product images efficiently', async ({ page }) => {
      await test.step('Test image loading performance', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Measure time for all images to load
        const startTime = performance.now();

        await page.waitForFunction(() => {
          const images = document.querySelectorAll('img');
          return Array.from(images).every(img => img.complete);
        }, { timeout: 10000 });

        const imageLoadTime = performance.now() - startTime;

        // Images should load within 5 seconds
        expect(imageLoadTime).toBeLessThan(5000);
      });
    });

    test('should handle lazy loading properly', async ({ page }) => {
      await test.step('Test lazy loading behavior', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Check that images are lazy loaded
        const visibleImages = await page.locator('img[loading="lazy"]').count();

        if (visibleImages > 0) {
          // Scroll down to trigger lazy loading
          await page.evaluate(() => {
            window.scrollTo(0, document.body.scrollHeight);
          });

          // Wait for lazy images to load
          await page.waitForTimeout(2000);

          const loadedImages = await page.locator('img[loading="lazy"][src]').count();
          expect(loadedImages).toBeGreaterThan(0);
        }
      });
    });
  });

  test.describe('Network Performance', () => {
    test('should handle slow network conditions', async ({ page }) => {
      await test.step('Simulate slow 3G network', async () => {
        // Simulate slow network
        const cdp = await page.context().newCDPSession(page);
        await cdp.send('Network.enable');
        await cdp.send('Network.emulateNetworkConditions', {
          offline: false,
          downloadThroughput: 500 * 1024, // 500 KB/s
          uploadThroughput: 500 * 1024,
          latency: 2000, // 2 second latency
        });

        const startTime = performance.now();

        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        const loadTime = performance.now() - startTime;

        // Should load within reasonable time even on slow network
        expect(loadTime).toBeLessThan(10000); // 10 seconds max for slow network

        await categoryPage.expectProductsDisplayed();
      });
    });

    test('should gracefully handle network errors', async ({ page }) => {
      await test.step('Test network error handling', async () => {
        // Intercept and fail some requests
        await page.route('**/api/v1/commerce/products**', route => {
          // Fail 50% of requests
          if (Math.random() > 0.5) {
            route.abort();
          } else {
            route.continue();
          }
        });

        await categoryPage.goto('electronics');

        // Should handle errors gracefully
        const errorMessage = page.getByTestId('error-message');
        const productGrid = page.getByTestId('product-grid');

        // Either show error message or products should load
        await expect(
          errorMessage.or(productGrid)
        ).toBeVisible({ timeout: 10000 });
      });
    });
  });

  test.describe('Memory Performance', () => {
    test('should not have memory leaks during navigation', async ({ page }) => {
      await test.step('Navigate between pages and monitor memory', async () => {
        const initialMemory = await page.evaluate(() => {
          return (performance as any).memory?.usedJSHeapSize || 0;
        });

        // Navigate between different pages
        const pages = ['electronics', 'clothing', 'books'];

        for (let i = 0; i < 3; i++) {
          for (const category of pages) {
            await categoryPage.goto(category);
            await categoryPage.waitForCategoryLoad();
          }
        }

        const finalMemory = await page.evaluate(() => {
          return (performance as any).memory?.usedJSHeapSize || 0;
        });

        if (initialMemory && finalMemory) {
          const memoryIncrease = finalMemory - initialMemory;
          const memoryIncreasePercent = (memoryIncrease / initialMemory) * 100;

          // Memory shouldn't increase by more than 50%
          expect(memoryIncreasePercent).toBeLessThan(50);
        }
      });
    });

    test('should handle large data sets efficiently', async ({ page }) => {
      await test.step('Load large product list and monitor memory', async () => {
        // Mock large product response
        await page.route('**/api/v1/commerce/products**', async route => {
          const largeProductList = Array(1000).fill(null).map((_, index) => ({
            id: `product-${index}`,
            name: `Product ${index}`,
            price: Math.random() * 1000,
            category: 'Electronics',
            image: `https://example.com/image-${index}.jpg`,
          }));

          await route.fulfill({
            json: {
              products: largeProductList,
              total: 1000,
            },
          });
        });

        const startMemory = await page.evaluate(() => {
          return (performance as any).memory?.usedJSHeapSize || 0;
        });

        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        const endMemory = await page.evaluate(() => {
          return (performance as any).memory?.usedJSHeapSize || 0;
        });

        if (startMemory && endMemory) {
          const memoryUsed = endMemory - startMemory;
          // Should not use more than 50MB for large product list
          expect(memoryUsed).toBeLessThan(50 * 1024 * 1024);
        }
      });
    });
  });

  test.describe('Mobile Performance', () => {
    test('should perform well on mobile devices', async ({ page, browser }) => {
      await test.step('Test mobile performance', async () => {
        // Use mobile context
        const mobileContext = await browser.newContext({
          userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15',
          viewport: { width: 375, height: 667 },
          deviceScaleFactor: 2,
        });

        const mobilePage = await mobileContext.newPage();
        const mobileCategoryPage = new CategoryPage(mobilePage);

        const startTime = performance.now();

        await mobileCategoryPage.goto('electronics');
        await mobileCategoryPage.waitForCategoryLoad();

        const loadTime = performance.now() - startTime;

        // Mobile should load within 4 seconds (allowing for slower processing)
        expect(loadTime).toBeLessThan(4000);

        await mobileCategoryPage.expectProductsDisplayed();

        await mobileContext.close();
      });
    });

    test('should handle touch interactions smoothly', async ({ page }) => {
      await test.step('Test touch interaction performance', async () => {
        // Set mobile viewport
        await page.setViewportSize({ width: 375, height: 667 });

        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        const startTime = performance.now();

        // Simulate touch scrolling
        await page.touchscreen.tap(200, 300);

        for (let i = 0; i < 5; i++) {
          await page.mouse.wheel(0, 300);
          await page.waitForTimeout(50);
        }

        const touchTime = performance.now() - startTime;

        // Touch interactions should be responsive
        expect(touchTime).toBeLessThan(1000);
      });
    });
  });
});