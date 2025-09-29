/**
 * Category Page Object Model
 * Encapsulates category browsing and product discovery interactions
 */

import { Page, Locator, expect } from '@playwright/test';
import { TestProduct } from '../../utils/test-data';

export class CategoryPage {
  readonly page: Page;
  readonly categoryContainer: Locator;
  readonly categoryTitle: Locator;
  readonly categoryDescription: Locator;
  readonly breadcrumbs: Locator;
  readonly productGrid: Locator;
  readonly productCards: Locator;
  readonly filterSidebar: Locator;
  readonly sortDropdown: Locator;
  readonly searchInput: Locator;
  readonly searchButton: Locator;
  readonly loadMoreButton: Locator;
  readonly paginationContainer: Locator;
  readonly resultCount: Locator;
  readonly noResultsMessage: Locator;
  readonly loadingSpinner: Locator;

  // Filter elements
  readonly priceFilter: Locator;
  readonly categoryFilter: Locator;
  readonly brandFilter: Locator;
  readonly ratingFilter: Locator;
  readonly availabilityFilter: Locator;
  readonly clearFiltersButton: Locator;

  // View toggle elements
  readonly gridViewButton: Locator;
  readonly listViewButton: Locator;

  constructor(page: Page) {
    this.page = page;
    this.categoryContainer = page.getByTestId('category-container');
    this.categoryTitle = page.getByTestId('category-title');
    this.categoryDescription = page.getByTestId('category-description');
    this.breadcrumbs = page.getByTestId('breadcrumbs');
    this.productGrid = page.getByTestId('product-grid');
    this.productCards = page.getByTestId('product-card');
    this.filterSidebar = page.getByTestId('filter-sidebar');
    this.sortDropdown = page.getByTestId('sort-dropdown');
    this.searchInput = page.getByTestId('category-search-input');
    this.searchButton = page.getByTestId('category-search-button');
    this.loadMoreButton = page.getByTestId('load-more-button');
    this.paginationContainer = page.getByTestId('pagination');
    this.resultCount = page.getByTestId('result-count');
    this.noResultsMessage = page.getByTestId('no-results-message');
    this.loadingSpinner = page.getByTestId('loading-spinner');

    // Filter elements
    this.priceFilter = page.getByTestId('price-filter');
    this.categoryFilter = page.getByTestId('category-filter');
    this.brandFilter = page.getByTestId('brand-filter');
    this.ratingFilter = page.getByTestId('rating-filter');
    this.availabilityFilter = page.getByTestId('availability-filter');
    this.clearFiltersButton = page.getByTestId('clear-filters-button');

    // View toggle elements
    this.gridViewButton = page.getByTestId('grid-view-button');
    this.listViewButton = page.getByTestId('list-view-button');
  }

  /**
   * Navigate to category page
   */
  async goto(categorySlug?: string): Promise<void> {
    const url = categorySlug ? `/categories/${categorySlug}` : '/categories';
    await this.page.goto(url);
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Wait for category page to load
   */
  async waitForCategoryLoad(): Promise<void> {
    await this.categoryContainer.waitFor({ state: 'visible' });
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Get product card by product ID
   */
  getProductCard(productId: string): Locator {
    return this.productGrid.locator(`[data-product-id="${productId}"]`);
  }

  /**
   * Get product card by index
   */
  getProductCardByIndex(index: number): Locator {
    return this.productCards.nth(index);
  }

  /**
   * Get add to cart button for a product
   */
  getAddToCartButton(productId: string): Locator {
    return this.getProductCard(productId).getByTestId('add-to-cart-button');
  }

  /**
   * Get product quick view button
   */
  getQuickViewButton(productId: string): Locator {
    return this.getProductCard(productId).getByTestId('quick-view-button');
  }

  /**
   * Search for products within category
   */
  async searchProducts(query: string): Promise<void> {
    await this.searchInput.fill(query);
    await this.searchButton.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Clear search query
   */
  async clearSearch(): Promise<void> {
    await this.searchInput.clear();
    await this.searchButton.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Sort products
   */
  async sortBy(sortOption: string): Promise<void> {
    await this.sortDropdown.click();
    await this.page.getByTestId(`sort-option-${sortOption}`).click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Apply price filter
   */
  async applyPriceFilter(minPrice: number, maxPrice: number): Promise<void> {
    const minPriceInput = this.priceFilter.getByTestId('min-price-input');
    const maxPriceInput = this.priceFilter.getByTestId('max-price-input');
    const applyButton = this.priceFilter.getByTestId('apply-price-filter');

    await minPriceInput.fill(minPrice.toString());
    await maxPriceInput.fill(maxPrice.toString());
    await applyButton.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Apply category filter
   */
  async applyCategoryFilter(categories: string[]): Promise<void> {
    for (const category of categories) {
      const categoryCheckbox = this.categoryFilter.getByTestId(`category-${category}`);
      await categoryCheckbox.check();
    }
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Apply brand filter
   */
  async applyBrandFilter(brands: string[]): Promise<void> {
    for (const brand of brands) {
      const brandCheckbox = this.brandFilter.getByTestId(`brand-${brand}`);
      await brandCheckbox.check();
    }
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Apply rating filter
   */
  async applyRatingFilter(minRating: number): Promise<void> {
    const ratingOption = this.ratingFilter.getByTestId(`rating-${minRating}`);
    await ratingOption.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Apply availability filter
   */
  async applyAvailabilityFilter(inStockOnly: boolean): Promise<void> {
    const availabilityCheckbox = this.availabilityFilter.getByTestId('in-stock-only');

    if (inStockOnly) {
      await availabilityCheckbox.check();
    } else {
      await availabilityCheckbox.uncheck();
    }

    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Clear all filters
   */
  async clearAllFilters(): Promise<void> {
    await this.clearFiltersButton.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Switch to grid view
   */
  async switchToGridView(): Promise<void> {
    await this.gridViewButton.click();
    await this.page.waitForTimeout(500); // Wait for view change animation
  }

  /**
   * Switch to list view
   */
  async switchToListView(): Promise<void> {
    await this.listViewButton.click();
    await this.page.waitForTimeout(500); // Wait for view change animation
  }

  /**
   * Load more products (infinite scroll)
   */
  async loadMoreProducts(): Promise<void> {
    if (await this.loadMoreButton.isVisible()) {
      await this.loadMoreButton.click();
      await this.page.waitForLoadState('networkidle');
    }
  }

  /**
   * Navigate to specific page
   */
  async goToPage(pageNumber: number): Promise<void> {
    const pageButton = this.paginationContainer.getByTestId(`page-${pageNumber}`);
    await pageButton.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Add product to cart from category page
   */
  async addProductToCart(productId: string, quantity: number = 1): Promise<void> {
    const addToCartButton = this.getAddToCartButton(productId);

    // Handle quantity selection if available
    const quantitySelector = this.getProductCard(productId).getByTestId('quantity-selector');
    if (await quantitySelector.isVisible({ timeout: 1000 })) {
      await quantitySelector.selectOption(quantity.toString());
    }

    await addToCartButton.click();

    // Wait for add to cart confirmation
    const confirmation = this.page.getByTestId('add-to-cart-confirmation');
    await expect(confirmation).toBeVisible({ timeout: 5000 });
  }

  /**
   * Open product quick view
   */
  async openQuickView(productId: string): Promise<void> {
    const quickViewButton = this.getQuickViewButton(productId);
    await quickViewButton.click();

    // Wait for quick view modal to open
    const quickViewModal = this.page.getByTestId('quick-view-modal');
    await expect(quickViewModal).toBeVisible();
  }

  /**
   * Navigate to product detail page
   */
  async goToProductDetail(productId: string): Promise<void> {
    const productCard = this.getProductCard(productId);
    const productLink = productCard.getByTestId('product-link');
    await productLink.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Verify category page elements
   */
  async expectCategoryPageLoaded(categoryName?: string): Promise<void> {
    await expect(this.categoryContainer).toBeVisible();
    await expect(this.productGrid).toBeVisible();

    if (categoryName) {
      await expect(this.categoryTitle).toContainText(categoryName);
    }
  }

  /**
   * Verify product count
   */
  async expectProductCount(count: number): Promise<void> {
    await expect(this.productCards).toHaveCount(count);
  }

  /**
   * Verify products are displayed
   */
  async expectProductsDisplayed(): Promise<void> {
    await expect(this.productCards.first()).toBeVisible();
    await expect(this.noResultsMessage).not.toBeVisible();
  }

  /**
   * Verify no products found
   */
  async expectNoProductsFound(): Promise<void> {
    await expect(this.productCards).toHaveCount(0);
    await expect(this.noResultsMessage).toBeVisible();
  }

  /**
   * Verify product card details
   */
  async expectProductCardDetails(productId: string, details: {
    name?: string;
    price?: string;
    rating?: string;
    availability?: 'in-stock' | 'out-of-stock';
  }): Promise<void> {
    const productCard = this.getProductCard(productId);

    if (details.name) {
      const nameElement = productCard.getByTestId('product-name');
      await expect(nameElement).toContainText(details.name);
    }

    if (details.price) {
      const priceElement = productCard.getByTestId('product-price');
      await expect(priceElement).toContainText(details.price);
    }

    if (details.rating) {
      const ratingElement = productCard.getByTestId('product-rating');
      await expect(ratingElement).toContainText(details.rating);
    }

    if (details.availability) {
      const availabilityElement = productCard.getByTestId('product-availability');

      if (details.availability === 'in-stock') {
        await expect(availabilityElement).toContainText('In Stock');
        await expect(this.getAddToCartButton(productId)).toBeEnabled();
      } else {
        await expect(availabilityElement).toContainText('Out of Stock');
        await expect(this.getAddToCartButton(productId)).toBeDisabled();
      }
    }
  }

  /**
   * Verify filters are applied
   */
  async expectFiltersApplied(activeFilters: {
    price?: { min: number; max: number };
    categories?: string[];
    brands?: string[];
    rating?: number;
    inStockOnly?: boolean;
  }): Promise<void> {
    // Verify active filter indicators
    const activeFilterContainer = this.page.getByTestId('active-filters');

    if (activeFilters.price) {
      const priceFilter = activeFilterContainer.getByTestId('active-price-filter');
      await expect(priceFilter).toBeVisible();
    }

    if (activeFilters.categories?.length) {
      for (const category of activeFilters.categories) {
        const categoryFilter = activeFilterContainer.getByTestId(`active-category-${category}`);
        await expect(categoryFilter).toBeVisible();
      }
    }

    if (activeFilters.brands?.length) {
      for (const brand of activeFilters.brands) {
        const brandFilter = activeFilterContainer.getByTestId(`active-brand-${brand}`);
        await expect(brandFilter).toBeVisible();
      }
    }

    if (activeFilters.rating) {
      const ratingFilter = activeFilterContainer.getByTestId('active-rating-filter');
      await expect(ratingFilter).toBeVisible();
    }

    if (activeFilters.inStockOnly) {
      const availabilityFilter = activeFilterContainer.getByTestId('active-availability-filter');
      await expect(availabilityFilter).toBeVisible();
    }
  }

  /**
   * Verify sort order
   */
  async expectSortOrder(sortType: 'price-asc' | 'price-desc' | 'name-asc' | 'name-desc' | 'rating' | 'newest'): Promise<void> {
    const productPrices: number[] = [];
    const productNames: string[] = [];

    const productCount = await this.productCards.count();

    for (let i = 0; i < Math.min(productCount, 5); i++) {
      const productCard = this.productCards.nth(i);

      if (sortType.includes('price')) {
        const priceText = await productCard.getByTestId('product-price').textContent();
        const price = parseFloat(priceText?.replace(/[^0-9.]/g, '') || '0');
        productPrices.push(price);
      }

      if (sortType.includes('name')) {
        const nameText = await productCard.getByTestId('product-name').textContent();
        productNames.push(nameText || '');
      }
    }

    switch (sortType) {
      case 'price-asc':
        for (let i = 1; i < productPrices.length; i++) {
          expect(productPrices[i]).toBeGreaterThanOrEqual(productPrices[i - 1]);
        }
        break;
      case 'price-desc':
        for (let i = 1; i < productPrices.length; i++) {
          expect(productPrices[i]).toBeLessThanOrEqual(productPrices[i - 1]);
        }
        break;
      case 'name-asc':
        for (let i = 1; i < productNames.length; i++) {
          expect(productNames[i].localeCompare(productNames[i - 1])).toBeGreaterThanOrEqual(0);
        }
        break;
      case 'name-desc':
        for (let i = 1; i < productNames.length; i++) {
          expect(productNames[i].localeCompare(productNames[i - 1])).toBeLessThanOrEqual(0);
        }
        break;
    }
  }

  /**
   * Verify breadcrumbs navigation
   */
  async expectBreadcrumbs(expectedBreadcrumbs: string[]): Promise<void> {
    for (let i = 0; i < expectedBreadcrumbs.length; i++) {
      const breadcrumb = this.breadcrumbs.getByTestId(`breadcrumb-${i}`);
      await expect(breadcrumb).toContainText(expectedBreadcrumbs[i]);
    }
  }

  /**
   * Test infinite scroll
   */
  async testInfiniteScroll(): Promise<void> {
    const initialProductCount = await this.productCards.count();

    // Scroll to bottom of page
    await this.page.evaluate(() => {
      window.scrollTo(0, document.body.scrollHeight);
    });

    // Wait for more products to load
    await this.page.waitForLoadState('networkidle');

    const finalProductCount = await this.productCards.count();
    expect(finalProductCount).toBeGreaterThan(initialProductCount);
  }

  /**
   * Test responsive grid layout
   */
  async testResponsiveGrid(): Promise<void> {
    // Test mobile layout
    await this.page.setViewportSize({ width: 375, height: 667 });
    await this.waitForCategoryLoad();

    const mobileColumns = await this.page.evaluate(() => {
      const grid = document.querySelector('[data-testid="product-grid"]');
      return window.getComputedStyle(grid!).gridTemplateColumns.split(' ').length;
    });
    expect(mobileColumns).toBe(1);

    // Test tablet layout
    await this.page.setViewportSize({ width: 768, height: 1024 });
    await this.waitForCategoryLoad();

    const tabletColumns = await this.page.evaluate(() => {
      const grid = document.querySelector('[data-testid="product-grid"]');
      return window.getComputedStyle(grid!).gridTemplateColumns.split(' ').length;
    });
    expect(tabletColumns).toBeGreaterThanOrEqual(2);

    // Test desktop layout
    await this.page.setViewportSize({ width: 1280, height: 720 });
    await this.waitForCategoryLoad();

    const desktopColumns = await this.page.evaluate(() => {
      const grid = document.querySelector('[data-testid="product-grid"]');
      return window.getComputedStyle(grid!).gridTemplateColumns.split(' ').length;
    });
    expect(desktopColumns).toBeGreaterThanOrEqual(3);
  }
}