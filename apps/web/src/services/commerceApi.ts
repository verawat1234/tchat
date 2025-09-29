/**
 * Commerce API Service
 *
 * RTK Query API slice for the commerce microservice.
 * Provides comprehensive CRUD operations for businesses, products, carts, categories, reviews, and wishlists.
 *
 * Features:
 * - Automatic caching with tag-based invalidation
 * - Optimistic updates for cart operations
 * - Error handling with fallback support
 * - Request deduplication
 * - TypeScript type safety
 */

import { api } from './api';
import type {
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
} from '../types/commerce';

/**
 * Commerce API endpoints
 */
export const commerceApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // ===== Business Endpoints =====

    /**
     * Get all businesses with filtering and pagination
     */
    getBusinesses: builder.query<BusinessResponse, {
      filters?: BusinessFilters;
      pagination?: Pagination;
      sort?: SortOptions;
    }>({
      query: ({ filters = {}, pagination = { page: 1, pageSize: 20 }, sort }) => ({
        url: '/commerce/businesses',
        method: 'GET',
        params: {
          ...filters,
          page: pagination.page,
          pageSize: pagination.pageSize,
          ...(sort && { sortField: sort.field, sortOrder: sort.order }),
        },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.businesses.map(({ id }) => ({ type: 'Business' as const, id })),
              { type: 'Business', id: 'LIST' },
            ]
          : [{ type: 'Business', id: 'LIST' }],
    }),

    /**
     * Get a business by ID
     */
    getBusiness: builder.query<Business, string>({
      query: (id) => `/commerce/businesses/${id}`,
      providesTags: (result, error, id) => [{ type: 'Business', id }],
    }),

    /**
     * Create a new business
     */
    createBusiness: builder.mutation<Business, CreateBusinessRequest>({
      query: (business) => ({
        url: '/commerce/businesses',
        method: 'POST',
        body: business,
      }),
      invalidatesTags: [{ type: 'Business', id: 'LIST' }],
    }),

    /**
     * Update a business
     */
    updateBusiness: builder.mutation<Business, { id: string; updates: UpdateBusinessRequest }>({
      query: ({ id, updates }) => ({
        url: `/commerce/businesses/${id}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Business', id },
        { type: 'Business', id: 'LIST' },
      ],
    }),

    /**
     * Delete a business
     */
    deleteBusiness: builder.mutation<void, string>({
      query: (id) => ({
        url: `/commerce/businesses/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Business', id },
        { type: 'Business', id: 'LIST' },
      ],
    }),

    // ===== Product Endpoints =====

    /**
     * Get all products with filtering and pagination
     */
    getProducts: builder.query<ProductResponse, {
      filters?: ProductFilters;
      pagination?: Pagination;
      sort?: SortOptions;
    }>({
      query: ({ filters = {}, pagination = { page: 1, pageSize: 20 }, sort }) => ({
        url: '/commerce/products',
        method: 'GET',
        params: {
          ...filters,
          page: pagination.page,
          pageSize: pagination.pageSize,
          ...(sort && { sortField: sort.field, sortOrder: sort.order }),
        },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.products.map(({ id }) => ({ type: 'Product' as const, id })),
              { type: 'Product', id: 'LIST' },
            ]
          : [{ type: 'Product', id: 'LIST' }],
    }),

    /**
     * Get a product by ID
     */
    getProduct: builder.query<Product, string>({
      query: (id) => `/commerce/products/${id}`,
      providesTags: (result, error, id) => [{ type: 'Product', id }],
    }),

    /**
     * Create a new product
     */
    createProduct: builder.mutation<Product, CreateProductRequest>({
      query: (product) => ({
        url: '/commerce/products',
        method: 'POST',
        body: product,
      }),
      invalidatesTags: [{ type: 'Product', id: 'LIST' }],
    }),

    /**
     * Update a product
     */
    updateProduct: builder.mutation<Product, { id: string; updates: UpdateProductRequest }>({
      query: ({ id, updates }) => ({
        url: `/commerce/products/${id}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Product', id },
        { type: 'Product', id: 'LIST' },
      ],
    }),

    /**
     * Delete a product
     */
    deleteProduct: builder.mutation<void, string>({
      query: (id) => ({
        url: `/commerce/products/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Product', id },
        { type: 'Product', id: 'LIST' },
      ],
    }),

    // ===== Cart Endpoints =====

    /**
     * Get current user's cart
     */
    getCart: builder.query<Cart, { userId?: string; sessionId?: string }>({
      query: ({ userId, sessionId }) => ({
        url: '/commerce/carts/current',
        method: 'GET',
        params: { userId, sessionId },
      }),
      providesTags: [{ type: 'Cart', id: 'CURRENT' }],
    }),

    /**
     * Get cart by ID
     */
    getCartById: builder.query<Cart, string>({
      query: (id) => `/commerce/carts/${id}`,
      providesTags: (result, error, id) => [{ type: 'Cart', id }],
    }),

    /**
     * Add item to cart with optimistic update
     */
    addToCart: builder.mutation<Cart, { cartId?: string; item: AddToCartRequest }>({
      query: ({ cartId, item }) => ({
        url: cartId ? `/commerce/carts/${cartId}/items` : '/commerce/carts/items',
        method: 'POST',
        body: item,
      }),
      // Optimistic update
      async onQueryStarted({ cartId }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          commerceApi.util.updateQueryData('getCart', { userId: undefined, sessionId: undefined }, (draft) => {
            // Add optimistic item to cart
            if (draft) {
              draft.itemCount += 1;
              // Note: Real implementation would add the actual item
            }
          })
        );
        try {
          await queryFulfilled;
        } catch {
          patchResult.undo();
        }
      },
      invalidatesTags: [
        { type: 'Cart', id: 'CURRENT' },
        { type: 'Cart', id: 'LIST' },
      ],
    }),

    /**
     * Update cart item
     */
    updateCartItem: builder.mutation<Cart, {
      cartId: string;
      itemId: string;
      updates: UpdateCartItemRequest
    }>({
      query: ({ cartId, itemId, updates }) => ({
        url: `/commerce/carts/${cartId}/items/${itemId}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { cartId }) => [
        { type: 'Cart', id: cartId },
        { type: 'Cart', id: 'CURRENT' },
      ],
    }),

    /**
     * Remove item from cart
     */
    removeFromCart: builder.mutation<Cart, { cartId: string; itemId: string }>({
      query: ({ cartId, itemId }) => ({
        url: `/commerce/carts/${cartId}/items/${itemId}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, { cartId }) => [
        { type: 'Cart', id: cartId },
        { type: 'Cart', id: 'CURRENT' },
      ],
    }),

    /**
     * Apply coupon to cart
     */
    applyCoupon: builder.mutation<Cart, { cartId: string; coupon: ApplyCouponRequest }>({
      query: ({ cartId, coupon }) => ({
        url: `/commerce/carts/${cartId}/coupons`,
        method: 'POST',
        body: coupon,
      }),
      invalidatesTags: (result, error, { cartId }) => [
        { type: 'Cart', id: cartId },
        { type: 'Cart', id: 'CURRENT' },
      ],
    }),

    /**
     * Remove coupon from cart
     */
    removeCoupon: builder.mutation<Cart, string>({
      query: (cartId) => ({
        url: `/commerce/carts/${cartId}/coupons`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, cartId) => [
        { type: 'Cart', id: cartId },
        { type: 'Cart', id: 'CURRENT' },
      ],
    }),

    /**
     * Validate cart contents
     */
    validateCart: builder.query<CartValidation, string>({
      query: (cartId) => `/commerce/carts/${cartId}/validate`,
      providesTags: (result, error, cartId) => [{ type: 'Cart', id: `${cartId}-validation` }],
    }),

    /**
     * Merge carts (guest to user)
     */
    mergeCarts: builder.mutation<Cart, MergeCartRequest>({
      query: (mergeRequest) => ({
        url: '/commerce/carts/merge',
        method: 'POST',
        body: mergeRequest,
      }),
      invalidatesTags: [
        { type: 'Cart', id: 'CURRENT' },
        { type: 'Cart', id: 'LIST' },
      ],
    }),

    /**
     * Get abandoned carts
     */
    getAbandonedCarts: builder.query<CartResponse, {
      filters?: Record<string, any>;
      pagination?: Pagination;
    }>({
      query: ({ filters = {}, pagination = { page: 1, pageSize: 20 } }) => ({
        url: '/commerce/carts/abandoned',
        method: 'GET',
        params: {
          ...filters,
          page: pagination.page,
          pageSize: pagination.pageSize,
        },
      }),
      providesTags: [{ type: 'Cart', id: 'ABANDONED' }],
    }),

    /**
     * Create abandonment tracking
     */
    createAbandonmentTracking: builder.mutation<CartAbandonmentTracking, CreateAbandonmentTrackingRequest>({
      query: (tracking) => ({
        url: '/commerce/carts/abandonment',
        method: 'POST',
        body: tracking,
      }),
      invalidatesTags: [{ type: 'Cart', id: 'ABANDONED' }],
    }),

    // ===== Category Endpoints =====

    /**
     * Get all categories
     */
    getCategories: builder.query<CategoryResponse, {
      businessId?: string;
      pagination?: Pagination;
      sort?: SortOptions;
    }>({
      query: ({ businessId, pagination = { page: 1, pageSize: 50 }, sort }) => ({
        url: '/commerce/categories',
        method: 'GET',
        params: {
          ...(businessId && { businessId }),
          page: pagination.page,
          pageSize: pagination.pageSize,
          ...(sort && { sortField: sort.field, sortOrder: sort.order }),
        },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.categories.map(({ id }) => ({ type: 'Category' as const, id })),
              { type: 'Category', id: 'LIST' },
            ]
          : [{ type: 'Category', id: 'LIST' }],
    }),

    /**
     * Get root categories
     */
    getRootCategories: builder.query<Category[], { businessId?: string }>({
      query: ({ businessId }) => ({
        url: '/commerce/categories/root',
        method: 'GET',
        params: { ...(businessId && { businessId }) },
      }),
      providesTags: [{ type: 'Category', id: 'ROOT' }],
    }),

    /**
     * Get category by ID
     */
    getCategory: builder.query<Category, string>({
      query: (id) => `/commerce/categories/${id}`,
      providesTags: (result, error, id) => [{ type: 'Category', id }],
    }),

    /**
     * Create a new category
     */
    createCategory: builder.mutation<Category, CreateCategoryRequest>({
      query: (category) => ({
        url: '/commerce/categories',
        method: 'POST',
        body: category,
      }),
      invalidatesTags: [
        { type: 'Category', id: 'LIST' },
        { type: 'Category', id: 'ROOT' },
      ],
    }),

    /**
     * Update a category
     */
    updateCategory: builder.mutation<Category, { id: string; updates: UpdateCategoryRequest }>({
      query: ({ id, updates }) => ({
        url: `/commerce/categories/${id}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Category', id },
        { type: 'Category', id: 'LIST' },
        { type: 'Category', id: 'ROOT' },
      ],
    }),

    /**
     * Delete a category
     */
    deleteCategory: builder.mutation<void, string>({
      query: (id) => ({
        url: `/commerce/categories/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Category', id },
        { type: 'Category', id: 'LIST' },
        { type: 'Category', id: 'ROOT' },
      ],
    }),

    /**
     * Add product to category
     */
    addProductToCategory: builder.mutation<void, AddProductToCategoryRequest>({
      query: (request) => ({
        url: '/commerce/categories/products',
        method: 'POST',
        body: request,
      }),
      invalidatesTags: (result, error, { categoryId, productId }) => [
        { type: 'Category', id: categoryId },
        { type: 'Product', id: productId },
      ],
    }),

    /**
     * Remove product from category
     */
    removeProductFromCategory: builder.mutation<void, { productId: string; categoryId: string }>({
      query: ({ productId, categoryId }) => ({
        url: `/commerce/categories/products`,
        method: 'DELETE',
        params: { productId, categoryId },
      }),
      invalidatesTags: (result, error, { categoryId, productId }) => [
        { type: 'Category', id: categoryId },
        { type: 'Product', id: productId },
      ],
    }),

    /**
     * Track category view
     */
    trackCategoryView: builder.mutation<void, TrackCategoryViewRequest>({
      query: (viewData) => ({
        url: '/commerce/categories/views',
        method: 'POST',
        body: viewData,
      }),
      // No cache invalidation for analytics tracking
    }),

    /**
     * Get category analytics
     */
    getCategoryAnalytics: builder.query<CategoryAnalytics, {
      categoryId: string;
      dateFrom: string;
      dateTo: string;
    }>({
      query: ({ categoryId, dateFrom, dateTo }) => ({
        url: `/commerce/categories/${categoryId}/analytics`,
        method: 'GET',
        params: { dateFrom, dateTo },
      }),
      providesTags: (result, error, { categoryId }) => [
        { type: 'Category', id: `${categoryId}-analytics` },
      ],
    }),

    // ===== Review Endpoints =====

    /**
     * Get reviews with filtering
     */
    getReviews: builder.query<ReviewResponse, {
      productId?: string;
      businessId?: string;
      pagination?: Pagination;
      sort?: SortOptions;
    }>({
      query: ({ productId, businessId, pagination = { page: 1, pageSize: 20 }, sort }) => ({
        url: '/commerce/reviews',
        method: 'GET',
        params: {
          ...(productId && { productId }),
          ...(businessId && { businessId }),
          page: pagination.page,
          pageSize: pagination.pageSize,
          ...(sort && { sortField: sort.field, sortOrder: sort.order }),
        },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.reviews.map(({ id }) => ({ type: 'Review' as const, id })),
              { type: 'Review', id: 'LIST' },
            ]
          : [{ type: 'Review', id: 'LIST' }],
    }),

    /**
     * Get review by ID
     */
    getReview: builder.query<Review, string>({
      query: (id) => `/commerce/reviews/${id}`,
      providesTags: (result, error, id) => [{ type: 'Review', id }],
    }),

    /**
     * Create a new review
     */
    createReview: builder.mutation<Review, CreateReviewRequest>({
      query: (review) => ({
        url: '/commerce/reviews',
        method: 'POST',
        body: review,
      }),
      invalidatesTags: [{ type: 'Review', id: 'LIST' }],
    }),

    /**
     * Update a review
     */
    updateReview: builder.mutation<Review, { id: string; updates: UpdateReviewRequest }>({
      query: ({ id, updates }) => ({
        url: `/commerce/reviews/${id}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Review', id },
        { type: 'Review', id: 'LIST' },
      ],
    }),

    /**
     * Delete a review
     */
    deleteReview: builder.mutation<void, string>({
      query: (id) => ({
        url: `/commerce/reviews/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Review', id },
        { type: 'Review', id: 'LIST' },
      ],
    }),

    // ===== Wishlist Endpoints =====

    /**
     * Get user's wishlists
     */
    getWishlists: builder.query<WishlistResponse, {
      userId: string;
      pagination?: Pagination;
    }>({
      query: ({ userId, pagination = { page: 1, pageSize: 20 } }) => ({
        url: '/commerce/wishlists',
        method: 'GET',
        params: {
          userId,
          page: pagination.page,
          pageSize: pagination.pageSize,
        },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.wishlists.map(({ id }) => ({ type: 'Wishlist' as const, id })),
              { type: 'Wishlist', id: 'LIST' },
            ]
          : [{ type: 'Wishlist', id: 'LIST' }],
    }),

    /**
     * Get wishlist by ID
     */
    getWishlist: builder.query<Wishlist, string>({
      query: (id) => `/commerce/wishlists/${id}`,
      providesTags: (result, error, id) => [{ type: 'Wishlist', id }],
    }),

    /**
     * Create a new wishlist
     */
    createWishlist: builder.mutation<Wishlist, CreateWishlistRequest>({
      query: (wishlist) => ({
        url: '/commerce/wishlists',
        method: 'POST',
        body: wishlist,
      }),
      invalidatesTags: [{ type: 'Wishlist', id: 'LIST' }],
    }),

    /**
     * Update a wishlist
     */
    updateWishlist: builder.mutation<Wishlist, { id: string; updates: UpdateWishlistRequest }>({
      query: ({ id, updates }) => ({
        url: `/commerce/wishlists/${id}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Wishlist', id },
        { type: 'Wishlist', id: 'LIST' },
      ],
    }),

    /**
     * Delete a wishlist
     */
    deleteWishlist: builder.mutation<void, string>({
      query: (id) => ({
        url: `/commerce/wishlists/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Wishlist', id },
        { type: 'Wishlist', id: 'LIST' },
      ],
    }),

    /**
     * Add product to wishlist
     */
    addToWishlist: builder.mutation<Wishlist, {
      wishlistId: string;
      item: AddToWishlistRequest
    }>({
      query: ({ wishlistId, item }) => ({
        url: `/commerce/wishlists/${wishlistId}/items`,
        method: 'POST',
        body: item,
      }),
      invalidatesTags: (result, error, { wishlistId }) => [
        { type: 'Wishlist', id: wishlistId },
        { type: 'Wishlist', id: 'LIST' },
      ],
    }),

    /**
     * Remove product from wishlist
     */
    removeFromWishlist: builder.mutation<Wishlist, {
      wishlistId: string;
      itemId: string
    }>({
      query: ({ wishlistId, itemId }) => ({
        url: `/commerce/wishlists/${wishlistId}/items/${itemId}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, { wishlistId }) => [
        { type: 'Wishlist', id: wishlistId },
        { type: 'Wishlist', id: 'LIST' },
      ],
    }),
  }),
});

// Export hooks for use in React components
export const {
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
} = commerceApi;

// Export the endpoints for use in other files
export const commerceEndpoints = commerceApi.endpoints;