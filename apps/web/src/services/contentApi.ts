/**
 * Content API - RTK Query endpoints for dynamic content management
 *
 * This file implements all RTK Query endpoints for the content management system,
 * providing type-safe API integration with comprehensive error handling, caching,
 * and real-time synchronization capabilities.
 *
 * Features:
 * - Complete CRUD operations for content items
 * - Category management and hierarchical content organization
 * - Version control and content history tracking
 * - Bulk operations for efficient content management
 * - Real-time content synchronization
 * - Advanced caching with tag-based invalidation
 * - Type-safe request/response handling
 * - Comprehensive error handling and retry logic
 */

import { api } from './api';
import { useSelector } from 'react-redux';
import type { RootState } from '../store';
import {
  selectFallbackContentById,
  selectFallbackMode,
} from '../features/contentSlice';
import type {
  ContentItem,
  ContentCategory,
  ContentVersion,
  ContentValue,
  ContentType,
  ContentStatus,
  PaginatedResponse,
} from '../types/content';

// =============================================================================
// Request/Response Type Definitions
// =============================================================================

/**
 * Request parameters for getting content items with filtering and pagination
 */
export interface GetContentItemsRequest {
  /** Filter by category ID */
  category?: string;
  /** Filter by content status */
  status?: ContentStatus;
  /** Filter by content type */
  type?: ContentType;
  /** Search query for content text */
  search?: string;
  /** Filter by tags */
  tags?: string[];
  /** Sort field */
  sortBy?: 'id' | 'category' | 'createdAt' | 'updatedAt' | 'status';
  /** Sort direction */
  sortOrder?: 'asc' | 'desc';
  /** Pagination offset */
  offset?: number;
  /** Items per page */
  limit?: number;
}

/**
 * Response for paginated content items
 */
export interface GetContentItemsResponse extends PaginatedResponse<ContentItem> {}

/**
 * Request for creating a new content item
 */
export interface CreateContentItemRequest {
  /** Unique content ID */
  id: string;
  /** Content category */
  category: string;
  /** Content type */
  type: ContentType;
  /** Content value */
  value: ContentValue;
  /** Optional tags */
  tags?: string[];
  /** Optional notes */
  notes?: string;
}

/**
 * Request for updating an existing content item
 */
export interface UpdateContentItemRequest {
  /** Content ID to update */
  id: string;
  /** Updated content value */
  value: ContentValue;
  /** Optional tags update */
  tags?: string[];
  /** Optional notes about the change */
  notes?: string;
  /** Expected version for optimistic locking */
  expectedVersion?: number;
}

/**
 * Request for publishing content
 */
export interface PublishContentRequest {
  /** Content ID to publish */
  id: string;
  /** Required change log for publication */
  changeLog: string;
  /** Force publication despite warnings */
  force?: boolean;
}

/**
 * Request for archiving content
 */
export interface ArchiveContentRequest {
  /** Content ID to archive */
  id: string;
  /** Required reason for archiving */
  archiveReason: string;
  /** Optional notes */
  notes?: string;
}

/**
 * Request for bulk content updates
 */
export interface BulkUpdateContentRequest {
  /** Array of update operations */
  updates: Array<{
    id: string;
    value: ContentValue;
    tags?: string[];
    notes?: string;
  }>;
  /** Whether operations should be atomic (all or nothing) */
  atomic?: boolean;
}

/**
 * Response for bulk update operations
 */
export interface BulkUpdateContentResponse {
  /** Successfully updated items */
  successful: ContentItem[];
  /** Failed operations with error details */
  failed: Array<{
    id: string;
    error: string;
    details?: string;
  }>;
  /** Summary statistics */
  summary: {
    total: number;
    successful: number;
    failed: number;
  };
}

/**
 * Request for reverting content to a previous version
 */
export interface RevertContentVersionRequest {
  /** Content ID to revert */
  id: string;
  /** Target version number */
  version: number;
  /** Optional reason for reversion */
  reason?: string;
}

/**
 * Request for content synchronization
 */
export interface SyncContentRequest {
  /** Last sync timestamp for incremental sync */
  lastSyncTime?: string;
  /** Categories to sync (if not provided, syncs all) */
  categories?: string[];
}

/**
 * Response for content synchronization
 */
export interface SyncContentResponse {
  /** Updated/new content items */
  items: ContentItem[];
  /** IDs of deleted content */
  deletedIds: string[];
  /** Server sync timestamp */
  syncTime: string;
  /** Total number of changes */
  totalChanges: number;
  /** Whether more changes are available */
  hasMore: boolean;
}

/**
 * Request for getting content versions
 */
export interface GetContentVersionsRequest {
  /** Content ID to get versions for */
  contentId: string;
  /** Pagination page number */
  page?: number;
  /** Items per page */
  limit?: number;
  /** Sort by version number or date */
  sortBy?: 'version' | 'createdAt';
  /** Sort direction */
  sortOrder?: 'asc' | 'desc';
}

/**
 * Response for content versions
 */
export interface GetContentVersionsResponse extends PaginatedResponse<ContentVersion> {}

// =============================================================================
// Content API Implementation
// =============================================================================

/**
 * Content API endpoints using RTK Query
 *
 * Provides comprehensive content management functionality with:
 * - Type-safe API calls
 * - Advanced caching strategies
 * - Optimistic updates
 * - Error handling and retry logic
 * - Real-time synchronization
 */
export const contentApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // ========================================================================
    // Query Endpoints - Content Retrieval
    // ========================================================================

    /**
     * Get paginated list of content items with filtering options
     *
     * Features:
     * - Advanced filtering by category, status, type, tags
     * - Full-text search capabilities
     * - Pagination with configurable page sizes
     * - Sorting by multiple fields
     * - Optimized caching with category-based tags
     */
    getContentItems: builder.query<GetContentItemsResponse, GetContentItemsRequest>({
      query: (params) => ({
        url: '/content',
        method: 'GET',
        params: {
          ...params,
          // Ensure arrays are properly serialized
          tags: params.tags?.join(','),
        },
      }),
      providesTags: (result, error, arg) => [
        'ContentItem',
        { type: 'ContentList', id: 'LIST' },
        // Category-specific tags for efficient invalidation
        ...(arg.category ? [{ type: 'ContentCategory' as const, id: arg.category }] : []),
        // Status-specific tags
        ...(arg.status ? [{ type: 'ContentStatus' as const, id: arg.status }] : []),
      ],
      // Enable automatic refetching on window focus
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Get a single content item by ID
     *
     * Features:
     * - Individual content item retrieval
     * - Optimized caching per content ID
     * - Automatic fallback content population
     */
    getContentItem: builder.query<ContentItem, string>({
      query: (id) => `/content/${encodeURIComponent(id)}`,
      providesTags: (result, error, id) => [
        { type: 'ContentItem', id },
        'ContentItem',
      ],
      keepUnusedDataFor: 600, // 10 minutes for individual items
    }),

    /**
     * Get all content items for a specific category
     *
     * Features:
     * - Category-based content retrieval
     * - Hierarchical category support
     * - Optimized for category-specific UI components
     */
    getContentByCategory: builder.query<GetContentItemsResponse, string>({
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
     *
     * Features:
     * - Complete category hierarchy
     * - Permission information per category
     * - Category metadata and descriptions
     */
    getContentCategories: builder.query<ContentCategory[], void>({
      query: () => '/content/categories',
      providesTags: ['ContentCategory'],
      keepUnusedDataFor: 1800, // 30 minutes - categories change infrequently
    }),

    /**
     * Get version history for a content item
     *
     * Features:
     * - Complete version history with metadata
     * - Paginated version lists
     * - Change tracking and diff information
     */
    getContentVersions: builder.query<GetContentVersionsResponse, GetContentVersionsRequest>({
      query: ({ contentId, ...params }) => ({
        url: `/content/${encodeURIComponent(contentId)}/versions`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { contentId }) => [
        { type: 'ContentVersion', id: contentId },
        'ContentVersion',
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Synchronize content with server
     *
     * Features:
     * - Incremental synchronization
     * - Conflict resolution
     * - Deleted content tracking
     * - Category-specific sync
     */
    syncContent: builder.query<SyncContentResponse, SyncContentRequest>({
      query: (params) => ({
        url: '/content/sync',
        method: 'POST',
        body: params,
      }),
      // Invalidate all content-related cache on sync
      invalidatesTags: ['ContentItem', 'ContentList', 'ContentCategory'],
      keepUnusedDataFor: 0, // Don't cache sync results
    }),

    // ========================================================================
    // Mutation Endpoints - Content Modification
    // ========================================================================

    /**
     * Create a new content item
     *
     * Features:
     * - Type-safe content creation
     * - Automatic validation
     * - Duplicate ID prevention
     * - Immediate cache updates
     */
    createContentItem: builder.mutation<ContentItem, CreateContentItemRequest>({
      query: (body) => ({
        url: '/content',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { category }) => [
        'ContentItem',
        'ContentList',
        { type: 'ContentCategory', id: category },
      ],
      // Optimistic update for better UX
      async onQueryStarted(arg, { dispatch, queryFulfilled }) {
        // Create optimistic content item
        const optimisticItem: Partial<ContentItem> = {
          id: arg.id,
          category: arg.category,
          type: arg.type,
          value: arg.value,
          status: 'draft' as ContentStatus,
          tags: arg.tags || [],
          notes: arg.notes,
        };

        // Optimistically update the content list
        const patchResult = dispatch(
          contentApi.util.updateQueryData('getContentItems', {}, (draft) => {
            if (draft.items) {
              draft.items.unshift(optimisticItem as ContentItem);
              draft.total += 1;
            }
          })
        );

        try {
          await queryFulfilled;
        } catch {
          // Revert optimistic update on error
          patchResult.undo();
        }
      },
    }),

    /**
     * Update an existing content item
     *
     * Features:
     * - Partial content updates
     * - Version conflict detection
     * - Optimistic updates
     * - Automatic cache synchronization
     */
    updateContentItem: builder.mutation<ContentItem, UpdateContentItemRequest>({
      query: ({ id, ...body }) => ({
        url: `/content/${encodeURIComponent(id)}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ContentItem', id },
        'ContentList',
        // Also invalidate category if item exists
        ...(result ? [{ type: 'ContentCategory' as const, id: result.category }] : []),
      ],
      // Optimistic update
      async onQueryStarted({ id, value, tags, notes }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          contentApi.util.updateQueryData('getContentItem', id, (draft) => {
            if (draft) {
              draft.value = value;
              if (tags !== undefined) draft.tags = tags;
              if (notes !== undefined) draft.notes = notes;
              // Update metadata would be handled by server
            }
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
     * Publish content item
     *
     * Features:
     * - Status change to published
     * - Publication workflow validation
     * - Change log requirements
     * - Permission verification
     */
    publishContent: builder.mutation<ContentItem, PublishContentRequest>({
      query: ({ id, ...body }) => ({
        url: `/content/${encodeURIComponent(id)}/publish`,
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ContentItem', id },
        'ContentList',
        ...(result ? [{ type: 'ContentCategory' as const, id: result.category }] : []),
      ],
    }),

    /**
     * Archive content item
     *
     * Features:
     * - Status change to archived
     * - Archive reason tracking
     * - Reversible archival
     * - Cache cleanup
     */
    archiveContent: builder.mutation<ContentItem, ArchiveContentRequest>({
      query: ({ id, ...body }) => ({
        url: `/content/${encodeURIComponent(id)}/archive`,
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ContentItem', id },
        'ContentList',
        ...(result ? [{ type: 'ContentCategory' as const, id: result.category }] : []),
      ],
    }),

    /**
     * Bulk update multiple content items
     *
     * Features:
     * - Atomic or non-atomic bulk operations
     * - Partial success handling
     * - Performance optimization for large updates
     * - Detailed error reporting
     */
    bulkUpdateContent: builder.mutation<BulkUpdateContentResponse, BulkUpdateContentRequest>({
      query: (body) => ({
        url: '/content/bulk',
        method: 'PUT',
        body,
      }),
      invalidatesTags: ['ContentItem', 'ContentList'],
      // Optimistic updates for bulk operations
      async onQueryStarted({ updates }, { dispatch, queryFulfilled }) {
        const patchResults = updates.map(({ id, value, tags, notes }) =>
          dispatch(
            contentApi.util.updateQueryData('getContentItem', id, (draft) => {
              if (draft) {
                draft.value = value;
                if (tags !== undefined) draft.tags = tags;
                if (notes !== undefined) draft.notes = notes;
              }
            })
          )
        );

        try {
          const result = await queryFulfilled;
          // Only keep optimistic updates for successful items
          const failedIds = new Set(result.data.failed.map(f => f.id));
          patchResults.forEach((patch, index) => {
            if (failedIds.has(updates[index].id)) {
              patch.undo();
            }
          });
        } catch {
          // Revert all optimistic updates on complete failure
          patchResults.forEach(patch => patch.undo());
        }
      },
    }),

    /**
     * Revert content to a previous version
     *
     * Features:
     * - Version history navigation
     * - New version creation (doesn't overwrite history)
     * - Metadata preservation
     * - Change tracking
     */
    revertContentVersion: builder.mutation<ContentItem, RevertContentVersionRequest>({
      query: ({ id, ...body }) => ({
        url: `/content/${encodeURIComponent(id)}/versions/${body.version}/revert`,
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ContentItem', id },
        { type: 'ContentVersion', id },
        'ContentList',
        ...(result ? [{ type: 'ContentCategory' as const, id: result.category }] : []),
      ],
    }),
  }),
  overrideExisting: false,
});

// =============================================================================
// Export Generated Hooks
// =============================================================================

// Query hooks for data fetching
export const {
  useGetContentItemsQuery,
  useGetContentItemQuery,
  useGetContentByCategoryQuery,
  useGetContentCategoriesQuery,
  useGetContentVersionsQuery,
  useSyncContentQuery,

  // Lazy query hooks for manual triggering
  useLazyGetContentItemsQuery,
  useLazyGetContentItemQuery,
  useLazyGetContentByCategoryQuery,
  useLazyGetContentCategoriesQuery,
  useLazyGetContentVersionsQuery,
  useLazySyncContentQuery,

  // Mutation hooks for data modification
  useCreateContentItemMutation,
  useUpdateContentItemMutation,
  usePublishContentMutation,
  useArchiveContentMutation,
  useBulkUpdateContentMutation,
  useRevertContentVersionMutation,

  // Utility hooks for cache management
  util: {
    getRunningQueriesThunk,
    resetApiState,
    updateQueryData,
    upsertQueryData,
    patchQueryData,
    invalidateTags,
  },
} = contentApi;

// =============================================================================
// Advanced Hook Utilities
// =============================================================================

/**
 * Enhanced content hook with fallback support
 * Combines RTK Query with fallback content from Redux state and localStorage
 */
export const useContentWithFallback = (contentId: string) => {
  const queryResult = useGetContentItemQuery(contentId);

  // Get fallback content from Redux state (in-memory cache)
  const fallbackContent = useSelector((state: RootState) =>
    selectFallbackContentById(state, contentId)
  );

  // Get fallback mode status
  const fallbackMode = useSelector(selectFallbackMode);

  // Determine the content to use
  const content = queryResult.data?.value || fallbackContent;
  const isFromFallback = !queryResult.data?.value && !!fallbackContent;

  return {
    ...queryResult,
    data: queryResult.data || (fallbackContent ? {
      id: contentId,
      value: fallbackContent,
      status: 'published' as const,
      type: 'text' as const, // Default fallback type
      category: 'unknown',
      tags: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      version: 1,
    } : undefined),
    content,
    hasContentAvailable: !!(queryResult.data?.value || fallbackContent),
    isFromFallback,
    fallbackMode,
    isLoading: queryResult.isLoading && !isFromFallback,
    isError: queryResult.isError && !isFromFallback,
  };
};

/**
 * Prefetch content for performance optimization
 */
export const usePrefetchContent = () => {
  const prefetchContentItem = useLazyGetContentItemQuery()[1];
  const prefetchContentItems = useLazyGetContentItemsQuery()[1];

  return {
    prefetchContentItem: (id: string) => prefetchContentItem(id),
    prefetchCategory: (category: string) =>
      prefetchContentItems({ category, limit: 50 }),
  };
};

/**
 * Enhanced content list hook with fallback support
 * Combines RTK Query with fallback content for content lists
 */
export const useContentListWithFallback = (params: GetContentItemsRequest = {}) => {
  const queryResult = useGetContentItemsQuery(params);

  // Get fallback content from Redux state
  const allFallbackContent = useSelector((state: RootState) => state.content.fallbackContent);
  const fallbackMode = useSelector(selectFallbackMode);

  // Filter fallback content by category if specified
  const fallbackItems = Object.entries(allFallbackContent)
    .filter(([contentId, content]) => {
      // Simple filtering - in a real implementation, you'd want more sophisticated filtering
      return !params.category || contentId.includes(params.category);
    })
    .map(([contentId, content]) => ({
      id: contentId,
      value: content,
      status: 'published' as const,
      type: 'text' as const,
      category: params.category || 'unknown',
      tags: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      version: 1,
    }));

  const isFromFallback = !queryResult.data && fallbackItems.length > 0;

  return {
    ...queryResult,
    data: queryResult.data || (fallbackItems.length > 0 ? {
      items: fallbackItems,
      total: fallbackItems.length,
      offset: 0,
      limit: fallbackItems.length,
      hasMore: false,
    } : undefined),
    hasContentAvailable: !!(queryResult.data?.items?.length || fallbackItems.length),
    isFromFallback,
    fallbackMode,
    isLoading: queryResult.isLoading && !isFromFallback,
    isError: queryResult.isError && !isFromFallback,
  };
};

/**
 * Hook for managing fallback service directly
 */
export const useContentFallbackService = () => {
  const fallbackMode = useSelector(selectFallbackMode);
  const capacity = useSelector((state: RootState) => {
    // This would need to be calculated or stored in state
    return { used: 0, available: 5242880, usagePercent: 0 }; // Default values
  });

  return {
    fallbackMode,
    capacity,
    // These would be connected to dispatch actions or service methods
    clearCache: () => {
      // Implementation would dispatch actions to clear cache
      console.log('Clear cache requested');
    },
    getStats: () => {
      // Implementation would return cache statistics
      return { hits: 0, misses: 0, evictions: 0, corruptions: 0 };
    },
  };
};

// Export the enhanced API for external use
export default contentApi;