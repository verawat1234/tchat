/**
 * Category Browsing E2E Tests
 * Tests for product discovery, filtering, and category navigation
 */

import { test, expect } from '@playwright/test';
import { CategoryPage } from '../page-objects/CategoryPage';
import { ProductPage } from '../page-objects/ProductPage';
import { CartPage } from '../page-objects/CartPage';
import { TestDataGenerator, TestDataSets } from '../../utils/test-data';
import { createApiHelpers } from '../../utils/api-helpers';
import { ScreenshotHelpers } from '../../utils/screenshot-helpers';

test.describe('Category Browsing', () => {
  let categoryPage: CategoryPage;
  let productPage: ProductPage;
  let cartPage: CartPage;
  let apiHelpers: ReturnType<typeof createApiHelpers>;
  let screenshotHelpers: ScreenshotHelpers;

  test.beforeEach(async ({ page, request }) => {
    categoryPage = new CategoryPage(page);
    productPage = new ProductPage(page);
    cartPage = new CartPage(page);
    apiHelpers = createApiHelpers(request);
    screenshotHelpers = new ScreenshotHelpers(page);
  });

  test.describe('Category Navigation', () => {
    test('should load category page with products', async () => {
      await test.step('Navigate to electronics category', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();
      });

      await test.step('Verify category page loaded', async () => {
        await categoryPage.expectCategoryPageLoaded('Electronics');
        await categoryPage.expectProductsDisplayed();
      });

      await test.step('Take category page screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'category-electronics',
          fullPage: true,
        });
      });
    });

    test('should display breadcrumbs correctly', async () => {
      await test.step('Navigate to electronics category', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();
      });

      await test.step('Verify breadcrumbs', async () => {
        await categoryPage.expectBreadcrumbs(['Home', 'Categories', 'Electronics']);
      });
    });

    test('should handle category with no products', async () => {
      await test.step('Navigate to empty category', async () => {
        await categoryPage.goto('empty-category');
        await categoryPage.waitForCategoryLoad();
      });

      await test.step('Verify no products message', async () => {
        await categoryPage.expectNoProductsFound();
      });

      await test.step('Take empty category screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'category-empty',
          fullPage: true,
        });
      });
    });
  });

  test.describe('Product Search', () => {
    test('should search products within category', async () => {
      const searchQuery = 'smartphone';

      await test.step('Navigate to electronics and search', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.searchProducts(searchQuery);
      });

      await test.step('Verify search results', async () => {
        await categoryPage.expectProductsDisplayed();

        // Verify search results contain query term
        const productCards = categoryPage.productCards;
        const firstProductName = await productCards.first().getByTestId('product-name').textContent();
        expect(firstProductName?.toLowerCase()).toContain(searchQuery.toLowerCase());
      });

      await test.step('Take search results screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'category-search-results',
          fullPage: true,
        });
      });
    });

    test('should handle search with no results', async () => {
      const searchQuery = 'nonexistentproduct123';

      await test.step('Search for non-existent product', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.searchProducts(searchQuery);
      });

      await test.step('Verify no results', async () => {
        await categoryPage.expectNoProductsFound();
      });
    });

    test('should clear search and show all products', async () => {
      await test.step('Search and then clear', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Get initial product count
        const initialCount = await categoryPage.productCards.count();

        // Search to reduce results
        await categoryPage.searchProducts('smartphone');
        const searchResultCount = await categoryPage.productCards.count();
        expect(searchResultCount).toBeLessThan(initialCount);

        // Clear search
        await categoryPage.clearSearch();
      });

      await test.step('Verify all products displayed again', async () => {
        await categoryPage.expectProductsDisplayed();
        const finalCount = await categoryPage.productCards.count();
        expect(finalCount).toBeGreaterThan(5); // Should show more products
      });
    });
  });

  test.describe('Product Sorting', () => {
    test('should sort products by price ascending', async () => {
      await test.step('Navigate and sort by price low to high', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.sortBy('price-asc');
      });

      await test.step('Verify price sorting', async () => {
        await categoryPage.expectSortOrder('price-asc');
      });

      await test.step('Take price sorted screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'category-sorted-price-asc',
          fullPage: true,
        });
      });
    });

    test('should sort products by price descending', async () => {
      await test.step('Navigate and sort by price high to low', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.sortBy('price-desc');
      });

      await test.step('Verify price sorting', async () => {
        await categoryPage.expectSortOrder('price-desc');
      });
    });

    test('should sort products by name', async () => {
      await test.step('Navigate and sort by name', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.sortBy('name-asc');
      });

      await test.step('Verify name sorting', async () => {
        await categoryPage.expectSortOrder('name-asc');
      });
    });

    test('should sort products by rating', async () => {
      await test.step('Navigate and sort by rating', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.sortBy('rating');
      });

      await test.step('Verify products sorted by rating', async () => {
        // Verify first few products have high ratings
        const productCards = categoryPage.productCards;
        const count = Math.min(await productCards.count(), 3);

        for (let i = 0; i < count; i++) {
          const ratingElement = productCards.nth(i).getByTestId('product-rating');
          const ratingText = await ratingElement.textContent();
          const rating = parseFloat(ratingText?.replace(/[^0-9.]/g, '') || '0');
          expect(rating).toBeGreaterThanOrEqual(4.0);
        }
      });
    });
  });

  test.describe('Product Filtering', () => {
    test('should filter products by price range', async () => {
      const minPrice = 100;
      const maxPrice = 300;

      await test.step('Apply price filter', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.applyPriceFilter(minPrice, maxPrice);
      });

      await test.step('Verify filtered results', async () => {
        await categoryPage.expectFiltersApplied({
          price: { min: minPrice, max: maxPrice },
        });

        // Verify product prices are within range
        const productCards = categoryPage.productCards;
        const count = Math.min(await productCards.count(), 5);

        for (let i = 0; i < count; i++) {
          const priceElement = productCards.nth(i).getByTestId('product-price');
          const priceText = await priceElement.textContent();
          const price = parseFloat(priceText?.replace(/[^0-9.]/g, '') || '0');
          expect(price).toBeGreaterThanOrEqual(minPrice);
          expect(price).toBeLessThanOrEqual(maxPrice);
        }
      });

      await test.step('Take filtered results screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'category-price-filtered',
          fullPage: true,
        });
      });
    });

    test('should filter products by availability', async () => {
      await test.step('Apply availability filter', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.applyAvailabilityFilter(true);
      });

      await test.step('Verify only in-stock products shown', async () => {
        await categoryPage.expectFiltersApplied({
          inStockOnly: true,
        });

        // Verify all products are in stock
        const productCards = categoryPage.productCards;
        const count = await productCards.count();

        for (let i = 0; i < count; i++) {
          const availabilityElement = productCards.nth(i).getByTestId('product-availability');
          await expect(availabilityElement).toContainText('In Stock');
        }
      });
    });

    test('should filter products by rating', async () => {
      const minRating = 4;

      await test.step('Apply rating filter', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.applyRatingFilter(minRating);
      });

      await test.step('Verify filtered results', async () => {
        await categoryPage.expectFiltersApplied({
          rating: minRating,
        });

        // Verify product ratings meet minimum
        const productCards = categoryPage.productCards;
        const count = Math.min(await productCards.count(), 5);

        for (let i = 0; i < count; i++) {
          const ratingElement = productCards.nth(i).getByTestId('product-rating');
          const ratingText = await ratingElement.textContent();
          const rating = parseFloat(ratingText?.replace(/[^0-9.]/g, '') || '0');
          expect(rating).toBeGreaterThanOrEqual(minRating);
        }
      });
    });

    test('should apply multiple filters simultaneously', async () => {
      const minPrice = 50;
      const maxPrice = 200;
      const minRating = 3;

      await test.step('Apply multiple filters', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.applyPriceFilter(minPrice, maxPrice);
        await categoryPage.applyRatingFilter(minRating);
        await categoryPage.applyAvailabilityFilter(true);
      });

      await test.step('Verify all filters applied', async () => {
        await categoryPage.expectFiltersApplied({
          price: { min: minPrice, max: maxPrice },
          rating: minRating,
          inStockOnly: true,
        });
      });

      await test.step('Take multi-filtered screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'category-multi-filtered',
          fullPage: true,
        });
      });
    });

    test('should clear all filters', async () => {
      await test.step('Apply filters then clear', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Apply some filters
        await categoryPage.applyPriceFilter(100, 300);
        await categoryPage.applyRatingFilter(4);

        // Get filtered count
        const filteredCount = await categoryPage.productCards.count();

        // Clear all filters
        await categoryPage.clearAllFilters();

        // Get final count
        const finalCount = await categoryPage.productCards.count();
        expect(finalCount).toBeGreaterThan(filteredCount);
      });
    });
  });

  test.describe('Product Grid Views', () => {
    test('should switch between grid and list views', async () => {
      await test.step('Navigate to category', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();
      });

      await test.step('Test grid view', async () => {
        await categoryPage.switchToGridView();

        // Take grid view screenshot
        await screenshotHelpers.expectScreenshot({
          name: 'category-grid-view',
        });
      });

      await test.step('Test list view', async () => {
        await categoryPage.switchToListView();

        // Take list view screenshot
        await screenshotHelpers.expectScreenshot({
          name: 'category-list-view',
        });
      });
    });

    test('should handle responsive grid layout', async () => {
      await test.step('Test responsive grid', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.testResponsiveGrid();
      });
    });
  });

  test.describe('Product Interactions', () => {
    test('should add product to cart from category page', async () => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Add product to cart', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.addProductToCart(testProduct.sku!, 2);
      });

      await test.step('Verify cart updated', async () => {
        await cartPage.goto();
        await cartPage.waitForCartLoad();

        await cartPage.expectCartHasItems(1);
        await cartPage.expectCartItemDetails(testProduct.sku!, {
          quantity: 2,
        });
      });
    });

    test('should open product quick view', async () => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Open quick view', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.openQuickView(testProduct.sku!);
      });

      await test.step('Verify quick view modal', async () => {
        const quickViewModal = categoryPage.page.getByTestId('quick-view-modal');
        await expect(quickViewModal).toBeVisible();

        // Verify product details in quick view
        const modalTitle = quickViewModal.getByTestId('product-title');
        await expect(modalTitle).toContainText(testProduct.name);
      });

      await test.step('Take quick view screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'product-quick-view',
        });
      });
    });

    test('should navigate to product detail page', async () => {
      const testProduct = TestDataSets.products.electronics;

      await test.step('Navigate to product detail', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.goToProductDetail(testProduct.sku!);
      });

      await test.step('Verify product page loaded', async () => {
        await productPage.waitForProductLoad();
        await productPage.expectProductDetails({
          name: testProduct.name,
        });
      });
    });
  });

  test.describe('Pagination and Loading', () => {
    test('should handle pagination', async () => {
      await test.step('Navigate to category with pagination', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();
      });

      await test.step('Test pagination', async () => {
        // Check if pagination is available
        const paginationContainer = categoryPage.paginationContainer;

        if (await paginationContainer.isVisible({ timeout: 2000 })) {
          // Go to page 2
          await categoryPage.goToPage(2);

          // Verify page changed
          const currentPageIndicator = paginationContainer.getByTestId('current-page');
          await expect(currentPageIndicator).toContainText('2');
        }
      });
    });

    test('should handle infinite scroll', async () => {
      await test.step('Test infinite scroll', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        // Only test if load more button exists
        const loadMoreButton = categoryPage.loadMoreButton;

        if (await loadMoreButton.isVisible({ timeout: 2000 })) {
          await categoryPage.testInfiniteScroll();
        }
      });
    });

    test('should display loading states', async () => {
      await test.step('Navigate and verify loading', async () => {
        await categoryPage.goto('electronics');

        // Check for loading spinner during initial load
        const loadingSpinner = categoryPage.loadingSpinner;

        // Loading spinner should appear briefly
        if (await loadingSpinner.isVisible({ timeout: 1000 })) {
          await expect(loadingSpinner).not.toBeVisible({ timeout: 5000 });
        }

        await categoryPage.expectProductsDisplayed();
      });
    });
  });

  test.describe('Mobile Category Experience', () => {
    test('should work correctly on mobile devices', async ({ page }) => {
      await test.step('Set mobile viewport', async () => {
        await page.setViewportSize({ width: 375, height: 667 });
      });

      await test.step('Test mobile category browsing', async () => {
        await categoryPage.goto('electronics');
        await categoryPage.waitForCategoryLoad();

        await categoryPage.expectCategoryPageLoaded();
        await categoryPage.expectProductsDisplayed();
      });

      await test.step('Test mobile filtering', async () => {
        // Mobile filter menu
        const filterButton = page.getByTestId('mobile-filter-button');

        if (await filterButton.isVisible({ timeout: 2000 })) {
          await filterButton.click();

          const filterModal = page.getByTestId('mobile-filter-modal');
          await expect(filterModal).toBeVisible();

          // Close filter modal
          const closeButton = filterModal.getByTestId('close-filter-modal');
          await closeButton.click();
        }
      });

      await test.step('Take mobile category screenshot', async () => {
        await screenshotHelpers.expectScreenshot({
          name: 'category-mobile-view',
          fullPage: true,
        });
      });
    });
  });

  test.describe('Error Handling', () => {
    test('should handle network errors gracefully', async ({ page }) => {
      await test.step('Simulate network error', async () => {
        // Intercept and fail product requests
        await page.route('**/api/v1/commerce/products**', route => {
          route.abort('failed');
        });

        await categoryPage.goto('electronics');
      });

      await test.step('Verify error handling', async () => {
        const errorMessage = page.getByTestId('category-error');
        await expect(errorMessage).toBeVisible();
        await expect(errorMessage).toContainText('Unable to load products');
      });
    });

    test('should handle slow loading gracefully', async ({ page }) => {
      await test.step('Simulate slow network', async () => {
        // Add delay to product requests
        await page.route('**/api/v1/commerce/products**', async route => {
          await new Promise(resolve => setTimeout(resolve, 2000));
          await route.continue();
        });

        await categoryPage.goto('electronics');
      });

      await test.step('Verify loading indicators', async () => {
        const loadingSpinner = categoryPage.loadingSpinner;
        await expect(loadingSpinner).toBeVisible();

        // Eventually products should load
        await categoryPage.expectProductsDisplayed();
        await expect(loadingSpinner).not.toBeVisible();
      });
    });
  });
});