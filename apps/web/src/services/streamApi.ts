import { api } from './api';

// Stream content types matching backend models
export interface StreamCategory {
  id: string;
  name: string;
  displayOrder: number;
  iconName: string;
  isActive: boolean;
  featuredContentEnabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface StreamSubtab {
  id: string;
  categoryId: string;
  name: string;
  displayOrder: number;
  filterCriteria: Record<string, any>;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface StreamContentItem {
  id: string;
  categoryId: string;
  title: string;
  description: string;
  thumbnailUrl: string;
  contentType: 'book' | 'podcast' | 'cartoon' | 'short_movie' | 'long_movie' | 'music' | 'art';
  duration?: number;
  price: number;
  currency: string;
  availabilityStatus: 'available' | 'coming_soon' | 'unavailable';
  isFeatured: boolean;
  featuredOrder?: number;
  metadata?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface ContentResponse {
  items: StreamContentItem[];
  total: number;
  hasMore: boolean;
}

export interface FeaturedResponse {
  items: StreamContentItem[];
  total: number;
  hasMore: boolean;
}

export interface TabNavigationState {
  id: number;
  userId: string;
  currentCategoryId: string;
  currentSubtabId?: string;
  lastVisitedAt: string;
  sessionId: string;
  devicePlatform: string;
  autoplayEnabled: boolean;
  showSubtabs: boolean;
  preferredViewMode: string;
  createdAt: string;
  updatedAt: string;
}

export interface StreamUserPreference {
  id: number;
  userId: string;
  preferredCategories: string;
  blockedCategories: string;
  autoplayEnabled: boolean;
  highQualityPreferred: boolean;
  offlineDownloadEnabled: boolean;
  notificationsEnabled: boolean;
  languagePreference: string;
  regionPreference: string;
  createdAt: string;
  updatedAt: string;
}

export interface StreamPurchaseRequest {
  mediaContentId: string;
  quantity: number;
  mediaLicense: string;
  downloadFormat?: string;
  cartId?: string;
}

export interface StreamPurchaseResponse {
  orderId: string;
  totalAmount: number;
  currency: string;
  success: boolean;
  message: string;
}

// Query parameters interfaces
export interface GetContentParams {
  categoryId: string;
  page?: number;
  limit?: number;
  subtabId?: string;
}

export interface GetFeaturedParams {
  categoryId: string;
  limit?: number;
}

export interface SearchContentParams {
  q: string;
  categoryId?: string;
  page?: number;
  limit?: number;
}

export interface UpdateNavigationParams {
  userId: string;
  sessionId: string;
  categoryId: string;
  subtabId?: string;
}

export interface UpdateProgressParams {
  userId: string;
  contentId: string;
  sessionId: string;
  progress: number;
}

/**
 * Stream API service with RTK Query endpoints
 *
 * Provides endpoints for:
 * - Stream categories and subtabs management
 * - Content browsing and search functionality
 * - Featured content display
 * - User navigation state tracking
 * - Content purchase functionality
 * - User preferences management
 */
export const streamApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // Category Management Endpoints
    getStreamCategories: builder.query<{ categories: StreamCategory[]; total: number; success: boolean }, void>({
      query: () => '/stream/categories',
      providesTags: ['StreamCategory'],
    }),

    getStreamCategoryDetail: builder.query<{
      category: StreamCategory;
      subtabs: StreamSubtab[];
      stats: Record<string, any>;
      success: boolean;
    }, string>({
      query: (categoryId) => `/stream/categories/${categoryId}`,
      providesTags: (result, error, categoryId) => [
        { type: 'StreamCategory', id: categoryId },
        'StreamSubtab',
      ],
    }),

    // Content Browsing Endpoints
    getStreamContent: builder.query<ContentResponse, GetContentParams>({
      query: ({ categoryId, page = 1, limit = 20, subtabId }) => ({
        url: '/stream/content',
        method: 'GET',
        params: {
          categoryId,
          page,
          limit,
          ...(subtabId && { subtabId }),
        },
      }),
      providesTags: (result, error, { categoryId, subtabId }) => [
        { type: 'StreamContent', id: `${categoryId}-${subtabId || 'all'}` },
        'StreamContent',
      ],
      serializeQueryArgs: ({ queryArgs }) => {
        // Custom serialization for content queries to handle pagination
        const { categoryId, subtabId, limit } = queryArgs;
        return `getStreamContent(${JSON.stringify({ categoryId, subtabId, limit })})`;
      },
    }),

    getStreamContentDetail: builder.query<StreamContentItem, string>({
      query: (contentId) => `/stream/content/${contentId}`,
      providesTags: (result, error, contentId) => [
        { type: 'StreamContent', id: contentId },
      ],
    }),

    getStreamFeatured: builder.query<FeaturedResponse, GetFeaturedParams>({
      query: ({ categoryId, limit = 10 }) => ({
        url: '/stream/featured',
        method: 'GET',
        params: {
          categoryId,
          limit,
        },
      }),
      providesTags: (result, error, { categoryId }) => [
        { type: 'StreamFeatured', id: categoryId },
        'StreamFeatured',
      ],
    }),

    searchStreamContent: builder.query<ContentResponse, SearchContentParams>({
      query: ({ q, categoryId, page = 1, limit = 20 }) => ({
        url: '/stream/search',
        method: 'GET',
        params: {
          q,
          ...(categoryId && { categoryId }),
          page,
          limit,
        },
      }),
      providesTags: ['StreamContent'],
    }),

    // Content Purchase Endpoint
    purchaseStreamContent: builder.mutation<StreamPurchaseResponse, StreamPurchaseRequest>({
      query: (purchaseData) => ({
        url: '/stream/content/purchase',
        method: 'POST',
        body: purchaseData,
      }),
      invalidatesTags: ['StreamPurchase'],
    }),

    // User Navigation and Session Endpoints
    getUserNavigationState: builder.query<TabNavigationState, string>({
      query: (userId) => `/stream/navigation?userId=${userId}`,
      providesTags: (result, error, userId) => [
        { type: 'StreamNavigation', id: userId },
      ],
    }),

    updateUserNavigationState: builder.mutation<TabNavigationState, UpdateNavigationParams>({
      query: ({ userId, sessionId, categoryId, subtabId }) => ({
        url: '/stream/navigation',
        method: 'PUT',
        body: {
          userId,
          sessionId,
          categoryId,
          ...(subtabId && { subtabId }),
        },
      }),
      invalidatesTags: (result, error, { userId }) => [
        { type: 'StreamNavigation', id: userId },
      ],
    }),

    updateContentViewProgress: builder.mutation<{ success: boolean }, UpdateProgressParams>({
      query: ({ userId, contentId, sessionId, progress }) => ({
        url: `/stream/content/${contentId}/progress`,
        method: 'PUT',
        body: {
          userId,
          sessionId,
          progress,
        },
      }),
      invalidatesTags: (result, error, { contentId }) => [
        { type: 'StreamContent', id: contentId },
      ],
    }),

    // User Preferences Endpoints
    getUserPreferences: builder.query<StreamUserPreference, string>({
      query: (userId) => `/stream/preferences?userId=${userId}`,
      providesTags: (result, error, userId) => [
        { type: 'StreamPreferences', id: userId },
      ],
    }),

    updateUserPreferences: builder.mutation<{ success: boolean }, {
      userId: string;
      preferences: Partial<Omit<StreamUserPreference, 'id' | 'userId' | 'createdAt' | 'updatedAt'>>;
    }>({
      query: ({ userId, preferences }) => ({
        url: '/stream/preferences',
        method: 'PUT',
        body: {
          userId,
          ...preferences,
        },
      }),
      invalidatesTags: (result, error, { userId }) => [
        { type: 'StreamPreferences', id: userId },
      ],
    }),
  }),
});

// Export hooks for use in components
export const {
  // Category hooks
  useGetStreamCategoriesQuery,
  useGetStreamCategoryDetailQuery,

  // Content hooks
  useGetStreamContentQuery,
  useGetStreamContentDetailQuery,
  useGetStreamFeaturedQuery,
  useSearchStreamContentQuery,

  // Purchase hooks
  usePurchaseStreamContentMutation,

  // Navigation hooks
  useGetUserNavigationStateQuery,
  useUpdateUserNavigationStateMutation,
  useUpdateContentViewProgressMutation,

  // Preferences hooks
  useGetUserPreferencesQuery,
  useUpdateUserPreferencesMutation,

  // Utility hooks
  useLazyGetStreamContentQuery,
  useLazySearchStreamContentQuery,
} = streamApi;

// Export types for external use
export type {
  StreamCategory,
  StreamSubtab,
  StreamContentItem,
  ContentResponse,
  FeaturedResponse,
  TabNavigationState,
  StreamUserPreference,
  StreamPurchaseRequest,
  StreamPurchaseResponse,
  GetContentParams,
  GetFeaturedParams,
  SearchContentParams,
  UpdateNavigationParams,
  UpdateProgressParams,
};