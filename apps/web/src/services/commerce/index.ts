/**
 * Commerce Services Export
 *
 * Centralized export for all commerce-related services, hooks, and utilities.
 */

// Re-export commerce API
export {
  commerceApi,
  commerceEndpoints,
  // Business hooks
  useGetBusinessesQuery,
  useGetBusinessQuery,
  useCreateBusinessMutation,
  useUpdateBusinessMutation,
  useDeleteBusinessMutation,
  // Product hooks
  useGetProductsQuery,
  useGetProductQuery,
  useCreateProductMutation,
  useUpdateProductMutation,
  useDeleteProductMutation,
  // Cart hooks
  useGetCartQuery,
  useGetCartByIdQuery,
  useAddToCartMutation,
  useUpdateCartItemMutation,
  useRemoveFromCartMutation,
  useApplyCouponMutation,
  useRemoveCouponMutation,
  useValidateCartQuery,
  useMergeCartsMutation,
  useGetAbandonedCartsQuery,
  useCreateAbandonmentTrackingMutation,
  // Category hooks
  useGetCategoriesQuery,
  useGetRootCategoriesQuery,
  useGetCategoryQuery,
  useCreateCategoryMutation,
  useUpdateCategoryMutation,
  useDeleteCategoryMutation,
  useAddProductToCategoryMutation,
  useRemoveProductFromCategoryMutation,
  useTrackCategoryViewMutation,
  useGetCategoryAnalyticsQuery,
  // Review hooks
  useGetReviewsQuery,
  useGetReviewQuery,
  useCreateReviewMutation,
  useUpdateReviewMutation,
  useDeleteReviewMutation,
  // Wishlist hooks
  useGetWishlistsQuery,
  useGetWishlistQuery,
  useCreateWishlistMutation,
  useUpdateWishlistMutation,
  useDeleteWishlistMutation,
  useAddToWishlistMutation,
  useRemoveFromWishlistMutation,
} from '../commerceApi';

// Re-export commerce types
export type {
  // Core Types
  Business,
  Product,
  Cart,
  Category,
  Review,
  Wishlist,
  CartValidation,
  CartAbandonmentTracking,
  CategoryAnalytics,
  // Enums
  BusinessVerificationStatus,
  ProductStatus,
  ProductType,
  CartStatus,
  CategoryStatus,
  CategoryType,
  ReviewType,
  ReviewStatus,
  WishlistType,
  WishlistPrivacy,
  // Request Types
  CreateBusinessRequest,
  UpdateBusinessRequest,
  BusinessFilters,
  CreateProductRequest,
  UpdateProductRequest,
  ProductFilters,
  AddToCartRequest,
  UpdateCartItemRequest,
  ApplyCouponRequest,
  MergeCartRequest,
  CreateAbandonmentTrackingRequest,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  AddProductToCategoryRequest,
  TrackCategoryViewRequest,
  CreateReviewRequest,
  UpdateReviewRequest,
  CreateWishlistRequest,
  UpdateWishlistRequest,
  AddToWishlistRequest,
  // Response Types
  BusinessResponse,
  ProductResponse,
  CartResponse,
  AbandonmentTrackingResponse,
  CategoryResponse,
  ReviewResponse,
  WishlistResponse,
  // Utility Types
  Pagination,
  SortOptions,
} from '../../types/commerce';

// Commerce-specific utilities and helpers
export const commerceUtils = {
  /**
   * Format price for display
   */
  formatPrice: (price: string, currency: string = 'USD'): string => {
    const numPrice = parseFloat(price);
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(numPrice);
  },

  /**
   * Calculate cart total
   */
  calculateCartTotal: (items: any[]): string => {
    const total = items.reduce((sum, item) => {
      return sum + parseFloat(item.totalPrice || '0');
    }, 0);
    return total.toFixed(2);
  },

  /**
   * Generate cart item key for deduplication
   */
  getCartItemKey: (productId: string, variantId?: string): string => {
    return variantId ? `${productId}-${variantId}` : productId;
  },

  /**
   * Check if product is available
   */
  isProductAvailable: (product: any): boolean => {
    return product.status === 'active' &&
           product.isActive &&
           (!product.inventory?.trackQuantity || product.inventory?.quantity > 0);
  },

  /**
   * Get product availability text
   */
  getAvailabilityText: (product: any): string => {
    if (!product.isActive || product.status !== 'active') {
      return 'Unavailable';
    }

    if (product.inventory?.trackQuantity) {
      const qty = product.inventory.quantity;
      if (qty <= 0) {
        return product.inventory.allowBackorder ? 'Backorder' : 'Out of Stock';
      }
      if (qty <= product.inventory.lowStockThreshold) {
        return `Low Stock (${qty} left)`;
      }
      return 'In Stock';
    }

    return 'Available';
  },

  /**
   * Generate product slug
   */
  generateSlug: (name: string): string => {
    return name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/(^-|-$)/g, '');
  },

  /**
   * Validate SKU format
   */
  isValidSku: (sku: string): boolean => {
    // Basic SKU validation - alphanumeric with dashes and underscores
    return /^[a-zA-Z0-9-_]+$/.test(sku) && sku.length >= 3 && sku.length <= 50;
  },

  /**
   * Calculate discount percentage
   */
  calculateDiscountPercentage: (originalPrice: string, salePrice: string): number => {
    const original = parseFloat(originalPrice);
    const sale = parseFloat(salePrice);
    if (original <= 0 || sale >= original) return 0;
    return Math.round(((original - sale) / original) * 100);
  },

  /**
   * Get star rating display
   */
  getStarRating: (rating: string): { full: number; half: boolean; empty: number } => {
    const numRating = parseFloat(rating);
    const full = Math.floor(numRating);
    const half = numRating % 1 >= 0.5;
    const empty = 5 - full - (half ? 1 : 0);

    return { full, half, empty };
  },

  /**
   * Validate email for reviews
   */
  isValidEmail: (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  },

  /**
   * Get business category display name
   */
  getCategoryDisplayName: (category: string): string => {
    return category
      .split('_')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ');
  },

  /**
   * Calculate estimated shipping days
   */
  getEstimatedShippingDays: (handlingTime: string, shippingMethod: string = 'standard'): number => {
    const handling = parseInt(handlingTime) || 1;
    const shipping = shippingMethod === 'express' ? 2 : 5;
    return handling + shipping;
  },
};

// Commerce constants
export const COMMERCE_CONSTANTS = {
  // Pagination defaults
  DEFAULT_PAGE_SIZE: 20,
  MAX_PAGE_SIZE: 100,

  // Cart limits
  MAX_CART_ITEMS: 100,
  MAX_ITEM_QUANTITY: 99,
  CART_EXPIRY_DAYS: 30,

  // Product limits
  MAX_PRODUCT_IMAGES: 10,
  MAX_PRODUCT_TAGS: 20,
  MAX_SKU_LENGTH: 50,

  // Review limits
  MIN_RATING: 1,
  MAX_RATING: 5,
  MAX_REVIEW_LENGTH: 2000,
  MAX_REVIEW_IMAGES: 5,

  // Wishlist limits
  MAX_WISHLISTS_PER_USER: 10,
  MAX_WISHLIST_ITEMS: 200,

  // Business limits
  MAX_BUSINESS_NAME_LENGTH: 100,
  MAX_BUSINESS_DESCRIPTION_LENGTH: 1000,

  // Currency codes supported
  SUPPORTED_CURRENCIES: ['USD', 'EUR', 'GBP', 'THB', 'SGD', 'MYR', 'IDR', 'VND', 'PHP'],

  // Product statuses
  PRODUCT_STATUSES: ['draft', 'active', 'inactive', 'archived'] as const,

  // Cart statuses
  CART_STATUSES: ['active', 'abandoned', 'converted', 'expired'] as const,

  // Review types
  REVIEW_TYPES: ['product', 'business', 'order'] as const,
} as const;