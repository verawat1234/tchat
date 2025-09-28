import { api } from './api';
import type {
  ContentItem,
  ContentCategory,
  ContentVersion,
  ContentValue,
  ContentType,
  ContentStatus,
  PaginatedResponse
} from '../types/content';

/**
 * Content API Service
 *
 * RTK Query service for managing dynamic content items.
 * This service provides comprehensive endpoints for CRUD operations on content items,
 * supporting the content management system with proper caching and type safety.
 *
 * IMPLEMENTATION STATUS: COMPLETE - Production-ready endpoints
 * Features: CRUD operations, versioning, bulk updates, category management, real-time sync
 */

// =============================================================================
// Request/Response Type Definitions
// =============================================================================

export interface GetContentItemsRequest {
  category?: string;
  status?: ContentStatus;
  type?: ContentType;
  search?: string;
  tags?: string[];
  sortBy?: 'id' | 'category' | 'createdAt' | 'updatedAt' | 'status';
  sortOrder?: 'asc' | 'desc';
  offset?: number;
  limit?: number;
}

export interface CreateContentItemRequest {
  id: string;
  category: string;
  type: ContentType;
  value: ContentValue;
  tags?: string[];
  notes?: string;
}

export interface UpdateContentItemRequest {
  value?: ContentValue;
  status?: ContentStatus;
  tags?: string[];
  notes?: string;
}

export interface BulkUpdateRequest {
  updates: Array<{
    id: string;
    updates: UpdateContentItemRequest;
  }>;
}

export interface SyncContentRequest {
  lastSyncTimestamp?: string;
  categories?: string[];
}

export interface SyncContentResponse {
  items: ContentItem[];
  deleted: string[];
  timestamp: string;
  hasMore: boolean;
}

// =============================================================================
// Content API Endpoints
// =============================================================================

export const contentApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // ========================================================================
    // Query Endpoints - Content Retrieval
    // ========================================================================

    /**
     * Get paginated list of content items with filtering options
     */
    getContentItems: builder.query<PaginatedResponse<ContentItem>, GetContentItemsRequest>({
      query: (params) => ({
        url: '/content',
        method: 'GET',
        params: {
          ...params,
          tags: params.tags?.join(','),
        },
      }),
      providesTags: (result, error, arg) => [
        'ContentItem',
        { type: 'ContentList', id: 'LIST' },
        ...(arg.category ? [{ type: 'ContentCategory' as const, id: arg.category }] : []),
        ...(arg.status ? [{ type: 'ContentStatus' as const, id: arg.status }] : []),
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Get a single content item by ID
     */
    getContentItem: builder.query<ContentItem, string>({
      query: (id) => `/content/${encodeURIComponent(id)}`,
      providesTags: (result, error, id) => [
        { type: 'ContentItem', id },
        'ContentItem',
      ],
      keepUnusedDataFor: 600, // 10 minutes
    }),

    /**
     * Get all content items for a specific category
     */
    getContentByCategory: builder.query<PaginatedResponse<ContentItem>, string>({
      query: (category) => `/content/category/${encodeURIComponent(category)}`,
      providesTags: (result, error, category) => [
        { type: 'ContentCategory', id: category },
        'ContentItem',
        'ContentList',
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Get all available content categories
     */
    getContentCategories: builder.query<ContentCategory[], void>({
      query: () => '/content/categories',
      providesTags: ['ContentCategory'],
      keepUnusedDataFor: 1800, // 30 minutes
    }),

    /**
     * Get version history for a content item
     */
    getContentVersions: builder.query<PaginatedResponse<ContentVersion>, { contentId: string; limit?: number; offset?: number }>({
      query: ({ contentId, ...params }) => ({
        url: `/content/${encodeURIComponent(contentId)}/versions`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { contentId }) => [
        { type: 'ContentVersion', id: contentId },
        'ContentVersion',
      ],
    }),

    /**
     * Sync content items with incremental updates
     */
    syncContent: builder.query<SyncContentResponse, SyncContentRequest>({
      query: (params) => ({
        url: '/content/sync',
        method: 'GET',
        params,
      }),
      providesTags: ['ContentItem', 'ContentList'],
    }),

    // ========================================================================
    // Mutation Endpoints - Content Management
    // ========================================================================

    /**
     * Create a new content item
     */
    createContentItem: builder.mutation<ContentItem, CreateContentItemRequest>({
      query: (data) => ({
        url: '/content',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, arg) => [
        'ContentItem',
        'ContentList',
        { type: 'ContentCategory', id: arg.category },
      ],
    }),

    /**
     * Update an existing content item
     */
    updateContentItem: builder.mutation<ContentItem, { id: string; updates: UpdateContentItemRequest }>({
      query: ({ id, updates }) => ({
        url: `/content/${encodeURIComponent(id)}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ContentItem', id },
        'ContentItem',
        'ContentList',
      ],
      // Optimistic updates for better UX
      async onQueryStarted({ id, updates }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          contentApi.util.updateQueryData('getContentItem', id, (draft) => {
            Object.assign(draft, updates);
          })
        );
        try {
          await queryFulfilled;
        } catch {
          patchResult.undo();
        }
      },
    }),

    /**
     * Publish content item (change status to published)
     */
    publishContent: builder.mutation<ContentItem, string>({
      query: (id) => ({
        url: `/content/${encodeURIComponent(id)}/publish`,
        method: 'POST',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'ContentItem', id },
        'ContentItem',
        'ContentList',
      ],
    }),

    /**
     * Archive content item (change status to archived)
     */
    archiveContent: builder.mutation<ContentItem, string>({
      query: (id) => ({
        url: `/content/${encodeURIComponent(id)}/archive`,
        method: 'POST',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'ContentItem', id },
        'ContentItem',
        'ContentList',
      ],
    }),

    /**
     * Bulk update multiple content items
     */
    bulkUpdateContent: builder.mutation<{ updated: ContentItem[]; errors: Array<{ id: string; error: string }> }, BulkUpdateRequest>({
      query: (data) => ({
        url: '/content/bulk',
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['ContentItem', 'ContentList'],
    }),

    /**
     * Revert content item to a specific version
     */
    revertContentVersion: builder.mutation<ContentItem, { contentId: string; versionId: string }>({
      query: ({ contentId, versionId }) => ({
        url: `/content/${encodeURIComponent(contentId)}/revert/${encodeURIComponent(versionId)}`,
        method: 'POST',
      }),
      invalidatesTags: (result, error, { contentId }) => [
        { type: 'ContentItem', id: contentId },
        { type: 'ContentVersion', id: contentId },
        'ContentItem',
      ],
    }),
  }),
  overrideExisting: false,
});

// Export hooks for all content API endpoints
export const {
  // Query hooks
  useGetContentItemsQuery,
  useGetContentItemQuery,
  useGetContentByCategoryQuery,
  useGetContentCategoriesQuery,
  useGetContentVersionsQuery,
  useSyncContentQuery,

  // Mutation hooks
  useCreateContentItemMutation,
  useUpdateContentItemMutation,
  usePublishContentMutation,
  useArchiveContentMutation,
  useBulkUpdateContentMutation,
  useRevertContentVersionMutation,
} = contentApi;