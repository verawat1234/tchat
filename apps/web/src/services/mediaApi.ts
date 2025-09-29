// Media Store RTK Query API
// Generated for Media Store Tabs feature implementation

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type {
  MediaCategory,
  MediaSubtab,
  MediaContentItem,
  MediaProduct,
  MediaCartItem,
  MediaOrder,
  MediaCategoriesResponse,
  MediaContentResponse,
  MediaFeaturedResponse,
  MediaSubtabsResponse,
  MediaSearchResponse,
  AddMediaToCartRequest,
  AddMediaToCartResponse,
  UnifiedCartResponse,
  MediaCheckoutValidationRequest,
  MediaCheckoutValidationResponse,
  MediaOrdersResponse,
  MediaApiError
} from '../types/media';

// Base query configuration
const baseQuery = fetchBaseQuery({
  baseUrl: '/api/v1',
  prepareHeaders: (headers, { getState }) => {
    // Add authentication token if available
    const token = (getState() as any)?.auth?.token;
    if (token) {
      headers.set('authorization', `Bearer ${token}`);
    }
    headers.set('content-type', 'application/json');
    return headers;
  },
});

// Media API slice
export const mediaApi = createApi({
  reducerPath: 'mediaApi',
  baseQuery,
  tagTypes: [
    'MediaCategories',
    'MediaContent',
    'MediaFeatured',
    'MediaSubtabs',
    'MediaProducts',
    'MediaCart',
    'MediaOrders'
  ],
  endpoints: (builder) => ({
    // Media Categories
    getMediaCategories: builder.query<MediaCategoriesResponse, void>({
      query: () => '/media/categories',
      providesTags: ['MediaCategories'],
    }),

    getMediaCategory: builder.query<MediaCategory, string>({
      query: (categoryId) => `/media/categories/${categoryId}`,
      providesTags: (result, error, categoryId) => [
        { type: 'MediaCategories', id: categoryId }
      ],
    }),

    // Media Subtabs
    getMovieSubtabs: builder.query<MediaSubtabsResponse, void>({
      query: () => '/media/movies/subtabs',
      providesTags: ['MediaSubtabs'],
    }),

    // Media Content
    getContentByCategory: builder.query<
      MediaContentResponse,
      { categoryId: string; page?: number; limit?: number; subtab?: string }
    >({
      query: ({ categoryId, page = 1, limit = 20, subtab }) => {
        const params = new URLSearchParams({
          page: page.toString(),
          limit: limit.toString(),
        });
        if (subtab) params.append('subtab', subtab);

        return `/media/category/${categoryId}/content?${params}`;
      },
      providesTags: (result, error, { categoryId }) => [
        { type: 'MediaContent', id: categoryId }
      ],
    }),

    getFeaturedContent: builder.query<
      MediaFeaturedResponse,
      { limit?: number; categoryId?: string }
    >({
      query: ({ limit = 10, categoryId }) => {
        const params = new URLSearchParams();
        if (limit) params.append('limit', limit.toString());
        if (categoryId) params.append('categoryId', categoryId);

        return `/media/featured?${params}`;
      },
      providesTags: ['MediaFeatured'],
    }),

    searchMediaContent: builder.query<
      MediaSearchResponse,
      { query: string; categoryId?: string; page?: number; limit?: number }
    >({
      query: ({ query, categoryId, page = 1, limit = 20 }) => {
        const params = new URLSearchParams({
          q: query,
          page: page.toString(),
          limit: limit.toString(),
        });
        if (categoryId) params.append('categoryId', categoryId);

        return `/media/search?${params}`;
      },
      providesTags: ['MediaContent'],
    }),

    // Store Integration - Products
    getMediaProducts: builder.query<
      { products: MediaProduct[]; pagination: any },
      { categoryId?: string; page?: number; limit?: number }
    >({
      query: ({ categoryId, page = 1, limit = 20 }) => {
        const params = new URLSearchParams({
          page: page.toString(),
          limit: limit.toString(),
        });
        if (categoryId) params.append('categoryId', categoryId);

        return `/store/products/media?${params}`;
      },
      providesTags: ['MediaProducts'],
    }),

    // Store Integration - Cart
    addMediaToCart: builder.mutation<AddMediaToCartResponse, AddMediaToCartRequest>({
      query: (body) => ({
        url: '/store/cart/add-media',
        method: 'POST',
        body,
      }),
      invalidatesTags: ['MediaCart'],
    }),

    getUnifiedCart: builder.query<UnifiedCartResponse, void>({
      query: () => '/store/cart',
      providesTags: ['MediaCart'],
    }),

    removeMediaFromCart: builder.mutation<void, { cartItemId: string }>({
      query: ({ cartItemId }) => ({
        url: `/store/cart/items/${cartItemId}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['MediaCart'],
    }),

    updateMediaCartItem: builder.mutation<
      MediaCartItem,
      { cartItemId: string; quantity: number }
    >({
      query: ({ cartItemId, quantity }) => ({
        url: `/store/cart/items/${cartItemId}`,
        method: 'PATCH',
        body: { quantity },
      }),
      invalidatesTags: ['MediaCart'],
    }),

    // Store Integration - Checkout
    validateMediaCheckout: builder.mutation<
      MediaCheckoutValidationResponse,
      MediaCheckoutValidationRequest
    >({
      query: (body) => ({
        url: '/store/checkout/media-validation',
        method: 'POST',
        body,
      }),
    }),

    processMediaCheckout: builder.mutation<
      MediaOrder,
      {
        cartId: string;
        mediaItems: MediaCartItem[];
        paymentMethod: string;
        billingAddress?: string;
      }
    >({
      query: (body) => ({
        url: '/store/checkout/media',
        method: 'POST',
        body,
      }),
      invalidatesTags: ['MediaCart', 'MediaOrders'],
    }),

    // Store Integration - Orders
    getMediaOrders: builder.query<
      MediaOrdersResponse,
      { page?: number; limit?: number; status?: string }
    >({
      query: ({ page = 1, limit = 20, status }) => {
        const params = new URLSearchParams({
          page: page.toString(),
          limit: limit.toString(),
        });
        if (status) params.append('status', status);

        return `/store/orders/media?${params}`;
      },
      providesTags: ['MediaOrders'],
    }),

    getMediaOrder: builder.query<MediaOrder, string>({
      query: (orderId) => `/store/orders/media/${orderId}`,
      providesTags: (result, error, orderId) => [
        { type: 'MediaOrders', id: orderId }
      ],
    }),

    downloadMediaContent: builder.query<
      { downloadUrl: string; expiresAt: string },
      { orderItemId: string }
    >({
      query: ({ orderItemId }) => `/store/orders/media/items/${orderItemId}/download`,
    }),
  }),
});

// Export hooks for components
export const {
  // Categories
  useGetMediaCategoriesQuery,
  useGetMediaCategoryQuery,

  // Subtabs
  useGetMovieSubtabsQuery,

  // Content
  useGetContentByCategoryQuery,
  useGetFeaturedContentQuery,
  useSearchMediaContentQuery,

  // Products
  useGetMediaProductsQuery,

  // Cart
  useAddMediaToCartMutation,
  useGetUnifiedCartQuery,
  useRemoveMediaFromCartMutation,
  useUpdateMediaCartItemMutation,

  // Checkout
  useValidateMediaCheckoutMutation,
  useProcessMediaCheckoutMutation,

  // Orders
  useGetMediaOrdersQuery,
  useGetMediaOrderQuery,
  useDownloadMediaContentQuery,
} = mediaApi;

// Export reducer
export const mediaApiReducer = mediaApi.reducer;

// Export middleware
export const mediaApiMiddleware = mediaApi.middleware;

// Utility types for API states
export type MediaApiState = ReturnType<typeof mediaApi.reducer>;

// Error handling utility
export const isMediaApiError = (error: any): error is MediaApiError => {
  return error && typeof error === 'object' && 'error' in error && 'message' in error;
};

// Cache tag utilities
export const invalidateMediaCache = (dispatch: any) => {
  dispatch(mediaApi.util.invalidateTags(['MediaContent', 'MediaFeatured', 'MediaProducts']));
};

export const resetMediaCache = (dispatch: any) => {
  dispatch(mediaApi.util.resetApiState());
};

// Prefetch utilities
export const prefetchMediaCategories = (dispatch: any) => {
  dispatch(mediaApi.util.prefetch('getMediaCategories', undefined, { force: false }));
};

export const prefetchFeaturedContent = (dispatch: any, categoryId?: string) => {
  dispatch(mediaApi.util.prefetch('getFeaturedContent', { categoryId }, { force: false }));
};

// Optimistic update utilities
export const optimisticCartUpdate = (
  dispatch: any,
  cartData: UnifiedCartResponse,
  newItem: MediaCartItem
) => {
  dispatch(
    mediaApi.util.updateQueryData('getUnifiedCart', undefined, (draft) => {
      draft.mediaItems.push(newItem);
      draft.totalMediaAmount += newItem.totalPrice;
      draft.totalAmount += newItem.totalPrice;
      draft.itemsCount += newItem.quantity;
    })
  );
};