/**
 * Product Page Object Model
 * Encapsulates product detail page interactions and selectors
 */

import { Page, Locator, expect } from '@playwright/test';
import { TestProduct } from '../../utils/test-data';

export class ProductPage {
  readonly page: Page;
  readonly productContainer: Locator;
  readonly productTitle: Locator;
  readonly productDescription: Locator;
  readonly productPrice: Locator;
  readonly originalPrice: Locator;
  readonly discountBadge: Locator;
  readonly productImages: Locator;
  readonly mainImage: Locator;
  readonly thumbnailImages: Locator;
  readonly imageGallery: Locator;
  readonly variantSelector: Locator;
  readonly quantitySelector: Locator;
  readonly addToCartButton: Locator;
  readonly buyNowButton: Locator;
  readonly addToWishlistButton: Locator;
  readonly shareButton: Locator;
  readonly productRating: Locator;
  readonly reviewCount: Locator;
  readonly productSku: Locator;
  readonly productAvailability: Locator;
  readonly stockCount: Locator;
  readonly shippingInfo: Locator;
  readonly returnPolicy: Locator;
  readonly productTabs: Locator;
  readonly relatedProducts: Locator;
  readonly recentlyViewed: Locator;
  readonly breadcrumbs: Locator;

  // Review elements
  readonly reviewsSection: Locator;
  readonly writeReviewButton: Locator;
  readonly reviewForm: Locator;
  readonly reviewsList: Locator;

  // Product specification elements
  readonly specificationsTab: Locator;
  readonly specificationsTable: Locator;

  constructor(page: Page) {
    this.page = page;
    this.productContainer = page.getByTestId('product-container');
    this.productTitle = page.getByTestId('product-title');
    this.productDescription = page.getByTestId('product-description');
    this.productPrice = page.getByTestId('product-price');
    this.originalPrice = page.getByTestId('original-price');
    this.discountBadge = page.getByTestId('discount-badge');
    this.productImages = page.getByTestId('product-images');
    this.mainImage = page.getByTestId('main-product-image');
    this.thumbnailImages = page.getByTestId('thumbnail-images');
    this.imageGallery = page.getByTestId('image-gallery');
    this.variantSelector = page.getByTestId('variant-selector');
    this.quantitySelector = page.getByTestId('quantity-selector');
    this.addToCartButton = page.getByTestId('add-to-cart-button');
    this.buyNowButton = page.getByTestId('buy-now-button');
    this.addToWishlistButton = page.getByTestId('add-to-wishlist-button');
    this.shareButton = page.getByTestId('share-button');
    this.productRating = page.getByTestId('product-rating');
    this.reviewCount = page.getByTestId('review-count');
    this.productSku = page.getByTestId('product-sku');
    this.productAvailability = page.getByTestId('product-availability');
    this.stockCount = page.getByTestId('stock-count');
    this.shippingInfo = page.getByTestId('shipping-info');
    this.returnPolicy = page.getByTestId('return-policy');
    this.productTabs = page.getByTestId('product-tabs');
    this.relatedProducts = page.getByTestId('related-products');
    this.recentlyViewed = page.getByTestId('recently-viewed');
    this.breadcrumbs = page.getByTestId('breadcrumbs');

    // Review elements
    this.reviewsSection = page.getByTestId('reviews-section');
    this.writeReviewButton = page.getByTestId('write-review-button');
    this.reviewForm = page.getByTestId('review-form');
    this.reviewsList = page.getByTestId('reviews-list');

    // Product specification elements
    this.specificationsTab = page.getByTestId('specifications-tab');
    this.specificationsTable = page.getByTestId('specifications-table');
  }

  /**
   * Navigate to product page
   */
  async goto(productId: string): Promise<void> {
    await this.page.goto(`/products/${productId}`);
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Wait for product page to load
   */
  async waitForProductLoad(): Promise<void> {
    await this.productContainer.waitFor({ state: 'visible' });
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Get variant option
   */
  getVariantOption(variantType: string, value: string): Locator {
    return this.variantSelector.locator(`[data-variant-type="${variantType}"]`).getByTestId(`variant-${value}`);
  }

  /**
   * Select product variant
   */
  async selectVariant(variantType: string, value: string): Promise<void> {
    const variantOption = this.getVariantOption(variantType, value);
    await variantOption.click();

    // Wait for price and availability to update
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Set product quantity
   */
  async setQuantity(quantity: number): Promise<void> {
    // Try quantity input first
    const quantityInput = this.quantitySelector.locator('input[type="number"]');
    if (await quantityInput.isVisible({ timeout: 1000 })) {
      await quantityInput.fill(quantity.toString());
      return;
    }

    // Try quantity dropdown
    const quantityDropdown = this.quantitySelector.locator('select');
    if (await quantityDropdown.isVisible({ timeout: 1000 })) {
      await quantityDropdown.selectOption(quantity.toString());
      return;
    }

    // Try quantity buttons
    const currentQuantity = await this.getCurrentQuantity();
    const difference = quantity - currentQuantity;

    if (difference > 0) {
      const increaseButton = this.quantitySelector.getByTestId('increase-quantity');
      for (let i = 0; i < difference; i++) {
        await increaseButton.click();
      }
    } else if (difference < 0) {
      const decreaseButton = this.quantitySelector.getByTestId('decrease-quantity');
      for (let i = 0; i < Math.abs(difference); i++) {
        await decreaseButton.click();
      }
    }
  }

  /**
   * Get current quantity
   */
  async getCurrentQuantity(): Promise<number> {
    const quantityInput = this.quantitySelector.locator('input[type="number"]');
    if (await quantityInput.isVisible({ timeout: 1000 })) {
      const value = await quantityInput.inputValue();
      return parseInt(value) || 1;
    }

    const quantityDisplay = this.quantitySelector.getByTestId('quantity-display');
    if (await quantityDisplay.isVisible({ timeout: 1000 })) {
      const text = await quantityDisplay.textContent();
      return parseInt(text || '1') || 1;
    }

    return 1;
  }

  /**
   * Add product to cart
   */
  async addToCart(quantity: number = 1): Promise<void> {
    await this.setQuantity(quantity);
    await this.addToCartButton.click();

    // Wait for add to cart confirmation
    const confirmation = this.page.getByTestId('add-to-cart-confirmation');
    await expect(confirmation).toBeVisible({ timeout: 5000 });
  }

  /**
   * Buy product now (direct checkout)
   */
  async buyNow(quantity: number = 1): Promise<void> {
    await this.setQuantity(quantity);
    await this.buyNowButton.click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Add product to wishlist
   */
  async addToWishlist(): Promise<void> {
    await this.addToWishlistButton.click();

    // Wait for wishlist confirmation
    const confirmation = this.page.getByTestId('wishlist-confirmation');
    await expect(confirmation).toBeVisible({ timeout: 5000 });
  }

  /**
   * Share product
   */
  async shareProduct(method: 'copy-link' | 'email' | 'facebook' | 'twitter'): Promise<void> {
    await this.shareButton.click();

    const shareModal = this.page.getByTestId('share-modal');
    await expect(shareModal).toBeVisible();

    const shareOption = shareModal.getByTestId(`share-${method}`);
    await shareOption.click();

    if (method === 'copy-link') {
      const copyConfirmation = this.page.getByTestId('copy-confirmation');
      await expect(copyConfirmation).toBeVisible({ timeout: 3000 });
    }
  }

  /**
   * View product image gallery
   */
  async viewImageGallery(): Promise<void> {
    await this.mainImage.click();

    const galleryModal = this.page.getByTestId('image-gallery-modal');
    await expect(galleryModal).toBeVisible();
  }

  /**
   * Navigate through product images
   */
  async selectProductImage(index: number): Promise<void> {
    const thumbnail = this.thumbnailImages.locator('.thumbnail').nth(index);
    await thumbnail.click();

    // Wait for main image to update
    await this.page.waitForTimeout(500);
  }

  /**
   * Switch product tabs
   */
  async switchToTab(tabName: 'description' | 'specifications' | 'reviews' | 'shipping'): Promise<void> {
    const tab = this.productTabs.getByTestId(`${tabName}-tab`);
    await tab.click();

    const tabContent = this.page.getByTestId(`${tabName}-content`);
    await expect(tabContent).toBeVisible();
  }

  /**
   * Write a product review
   */
  async writeReview(review: {
    rating: number;
    title: string;
    comment: string;
    name?: string;
    email?: string;
  }): Promise<void> {
    await this.writeReviewButton.click();

    const reviewModal = this.page.getByTestId('review-modal');
    await expect(reviewModal).toBeVisible();

    // Set rating
    const ratingStars = reviewModal.getByTestId('rating-stars');
    const star = ratingStars.locator(`.star-${review.rating}`);
    await star.click();

    // Fill review form
    await reviewModal.getByTestId('review-title').fill(review.title);
    await reviewModal.getByTestId('review-comment').fill(review.comment);

    if (review.name) {
      await reviewModal.getByTestId('reviewer-name').fill(review.name);
    }

    if (review.email) {
      await reviewModal.getByTestId('reviewer-email').fill(review.email);
    }

    // Submit review
    const submitButton = reviewModal.getByTestId('submit-review');
    await submitButton.click();

    // Wait for review submission confirmation
    const confirmation = this.page.getByTestId('review-submitted');
    await expect(confirmation).toBeVisible({ timeout: 5000 });
  }

  /**
   * Verify product details
   */
  async expectProductDetails(product: {
    name?: string;
    price?: string;
    originalPrice?: string;
    discount?: string;
    sku?: string;
    availability?: 'in-stock' | 'out-of-stock' | 'limited-stock';
    rating?: string;
    reviewCount?: string;
  }): Promise<void> {
    if (product.name) {
      await expect(this.productTitle).toContainText(product.name);
    }

    if (product.price) {
      await expect(this.productPrice).toContainText(product.price);
    }

    if (product.originalPrice) {
      await expect(this.originalPrice).toContainText(product.originalPrice);
    }

    if (product.discount) {
      await expect(this.discountBadge).toContainText(product.discount);
    }

    if (product.sku) {
      await expect(this.productSku).toContainText(product.sku);
    }

    if (product.availability) {
      switch (product.availability) {
        case 'in-stock':
          await expect(this.productAvailability).toContainText('In Stock');
          await expect(this.addToCartButton).toBeEnabled();
          break;
        case 'out-of-stock':
          await expect(this.productAvailability).toContainText('Out of Stock');
          await expect(this.addToCartButton).toBeDisabled();
          break;
        case 'limited-stock':
          await expect(this.productAvailability).toContainText('Limited Stock');
          await expect(this.addToCartButton).toBeEnabled();
          break;
      }
    }

    if (product.rating) {
      await expect(this.productRating).toContainText(product.rating);
    }

    if (product.reviewCount) {
      await expect(this.reviewCount).toContainText(product.reviewCount);
    }
  }

  /**
   * Verify product variants
   */
  async expectVariantsAvailable(variants: {
    type: string;
    options: string[];
  }[]): Promise<void> {
    for (const variant of variants) {
      const variantContainer = this.variantSelector.locator(`[data-variant-type="${variant.type}"]`);
      await expect(variantContainer).toBeVisible();

      for (const option of variant.options) {
        const variantOption = this.getVariantOption(variant.type, option);
        await expect(variantOption).toBeVisible();
      }
    }
  }

  /**
   * Verify selected variant
   */
  async expectVariantSelected(variantType: string, value: string): Promise<void> {
    const variantOption = this.getVariantOption(variantType, value);
    await expect(variantOption).toHaveAttribute('aria-selected', 'true');
  }

  /**
   * Verify related products
   */
  async expectRelatedProducts(minCount: number = 1): Promise<void> {
    await expect(this.relatedProducts).toBeVisible();

    const relatedProductCards = this.relatedProducts.locator('[data-testid="product-card"]');
    const count = await relatedProductCards.count();
    expect(count).toBeGreaterThanOrEqual(minCount);
  }

  /**
   * Verify breadcrumbs
   */
  async expectBreadcrumbs(expectedBreadcrumbs: string[]): Promise<void> {
    for (let i = 0; i < expectedBreadcrumbs.length; i++) {
      const breadcrumb = this.breadcrumbs.getByTestId(`breadcrumb-${i}`);
      await expect(breadcrumb).toContainText(expectedBreadcrumbs[i]);
    }
  }

  /**
   * Verify product images
   */
  async expectProductImages(minCount: number = 1): Promise<void> {
    await expect(this.mainImage).toBeVisible();

    if (minCount > 1) {
      const thumbnails = this.thumbnailImages.locator('.thumbnail');
      const count = await thumbnails.count();
      expect(count).toBeGreaterThanOrEqual(minCount);
    }
  }

  /**
   * Verify quantity constraints
   */
  async expectQuantityConstraints(min: number = 1, max?: number): Promise<void> {
    const quantityInput = this.quantitySelector.locator('input[type="number"]');

    if (await quantityInput.isVisible({ timeout: 1000 })) {
      await expect(quantityInput).toHaveAttribute('min', min.toString());

      if (max !== undefined) {
        await expect(quantityInput).toHaveAttribute('max', max.toString());
      }
    }
  }

  /**
   * Test product page responsiveness
   */
  async testResponsiveLayout(): Promise<void> {
    // Test mobile layout
    await this.page.setViewportSize({ width: 375, height: 667 });
    await this.waitForProductLoad();
    await expect(this.productContainer).toBeVisible();

    // Verify mobile-specific layout changes
    const productInfo = this.page.getByTestId('product-info');
    const mobileLayout = await productInfo.evaluate(el => {
      const style = window.getComputedStyle(el);
      return style.flexDirection === 'column';
    });
    expect(mobileLayout).toBe(true);

    // Test tablet layout
    await this.page.setViewportSize({ width: 768, height: 1024 });
    await this.waitForProductLoad();
    await expect(this.productContainer).toBeVisible();

    // Test desktop layout
    await this.page.setViewportSize({ width: 1280, height: 720 });
    await this.waitForProductLoad();
    await expect(this.productContainer).toBeVisible();

    const desktopLayout = await productInfo.evaluate(el => {
      const style = window.getComputedStyle(el);
      return style.flexDirection === 'row';
    });
    expect(desktopLayout).toBe(true);
  }

  /**
   * Test product zoom functionality
   */
  async testImageZoom(): Promise<void> {
    // Hover over main image to trigger zoom
    await this.mainImage.hover();

    const zoomContainer = this.page.getByTestId('image-zoom');
    await expect(zoomContainer).toBeVisible({ timeout: 2000 });

    // Move mouse to test zoom tracking
    await this.mainImage.hover({ position: { x: 100, y: 100 } });
    await this.page.waitForTimeout(500);

    await this.mainImage.hover({ position: { x: 200, y: 200 } });
    await this.page.waitForTimeout(500);
  }

  /**
   * Test product recommendations
   */
  async testProductRecommendations(): Promise<void> {
    await expect(this.relatedProducts).toBeVisible();

    const relatedProductCards = this.relatedProducts.locator('[data-testid="product-card"]');
    const count = await relatedProductCards.count();

    if (count > 0) {
      // Click on first related product
      const firstRelatedProduct = relatedProductCards.first();
      const productLink = firstRelatedProduct.getByTestId('product-link');
      await productLink.click();

      // Verify navigation to related product
      await this.waitForProductLoad();
      await expect(this.productContainer).toBeVisible();
    }
  }
}