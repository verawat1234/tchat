/**
 * Content Management API Contracts
 *
 * RTK Query endpoint definitions for dynamic content management.
 * These contracts define the API interface between frontend and backend.
 */

import { api } from '../../../apps/web/src/services/api';

// Types from data model
export interface ContentItem {
  id: string;
  category: string;
  type: ContentType;
  value: ContentValue;
  metadata: ContentMetadata;
  status: ContentStatus;
}

export interface ContentCategory {
  id: string;
  name: string;
  description: string;
  parentId?: string;
  permissions: CategoryPermissions;
}

export interface ContentVersion {
  id: string;
  contentId: string;
  version: number;
  value: ContentValue;
  metadata: ContentMetadata;
  changeLog: string;
}

export enum ContentType {
  TEXT = 'text',
  RICH_TEXT = 'rich_text',
  IMAGE_URL = 'image_url',
  CONFIG = 'config',
  TRANSLATION = 'translation'
}

export enum ContentStatus {
  DRAFT = 'draft',
  PUBLISHED = 'published',
  ARCHIVED = 'archived'
}

export type ContentValue = {
  type: ContentType;
  value: any;
  [key: string]: any;
};

export interface ContentMetadata {
  createdAt: string;
  createdBy: string;
  updatedAt: string;
  updatedBy: string;
  version: number;
  tags?: string[];
  notes?: string;
}

export interface CategoryPermissions {
  read: string[];
  write: string[];
  publish: string[];
}

// Request/Response types
export interface GetContentItemsRequest {
  category?: string;
  status?: ContentStatus;
  type?: ContentType;
  limit?: number;
  offset?: number;
  search?: string;
  tags?: string[];
}

export interface GetContentItemsResponse {
  items: ContentItem[];
  total: number;
  hasMore: boolean;
}

export interface UpdateContentItemRequest {
  id: string;
  value: ContentValue;
  notes?: string;
}

export interface CreateContentItemRequest {
  id: string;
  category: string;
  type: ContentType;
  value: ContentValue;
  tags?: string[];
  notes?: string;
}

export interface PublishContentRequest {
  id: string;
  changeLog: string;
}

export interface GetContentVersionsRequest {
  contentId: string;
  limit?: number;
  offset?: number;
}

export interface GetContentVersionsResponse {
  versions: ContentVersion[];
  total: number;
}

export interface BulkUpdateContentRequest {
  updates: Array<{
    id: string;
    value: ContentValue;
    notes?: string;
  }>;
}

export interface ContentSyncRequest {
  lastSyncTime?: string;
  categories?: string[];
}

export interface ContentSyncResponse {
  items: ContentItem[];
  deletedIds: string[];
  syncTime: string;
}

// RTK Query API definition
export const contentApi = api.injectEndpoints({
  endpoints: (builder) => ({

    // Query endpoints
    getContentItems: builder.query<GetContentItemsResponse, GetContentItemsRequest>({
      query: (params) => ({
        url: '/content/items',
        params,
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.items.map((item) => ({ type: 'ContentItem' as const, id: item.id })),
              { type: 'ContentItem', id: 'LIST' },
            ]
          : [{ type: 'ContentItem', id: 'LIST' }],
    }),

    getContentItem: builder.query<ContentItem, string>({
      query: (id) => `/content/items/${id}`,
      providesTags: (result, error, id) => [{ type: 'ContentItem', id }],
    }),

    getContentByCategory: builder.query<ContentItem[], string>({
      query: (category) => `/content/categories/${category}/items`,
      providesTags: (result, error, category) => [
        { type: 'ContentItem', id: `CATEGORY_${category}` },
      ],
    }),

    getContentCategories: builder.query<ContentCategory[], void>({
      query: () => '/content/categories',
      providesTags: [{ type: 'ContentCategory', id: 'LIST' }],
    }),

    getContentVersions: builder.query<GetContentVersionsResponse, GetContentVersionsRequest>({
      query: ({ contentId, ...params }) => ({
        url: `/content/items/${contentId}/versions`,
        params,
      }),
      providesTags: (result, error, { contentId }) => [
        { type: 'ContentVersion', id: contentId },
      ],
    }),

    syncContent: builder.query<ContentSyncResponse, ContentSyncRequest>({
      query: (params) => ({
        url: '/content/sync',
        params,
      }),
      providesTags: [{ type: 'ContentItem', id: 'SYNC' }],
    }),

    // Mutation endpoints
    createContentItem: builder.mutation<ContentItem, CreateContentItemRequest>({
      query: (data) => ({
        url: '/content/items',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: [
        { type: 'ContentItem', id: 'LIST' },
        (arg) => ({ type: 'ContentItem', id: `CATEGORY_${arg.category}` }),
      ],
    }),

    updateContentItem: builder.mutation<ContentItem, UpdateContentItemRequest>({
      query: ({ id, ...data }) => ({
        url: `/content/items/${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ContentItem', id },
        { type: 'ContentItem', id: 'LIST' },
      ],
    }),

    publishContent: builder.mutation<ContentItem, PublishContentRequest>({
      query: ({ id, ...data }) => ({
        url: `/content/items/${id}/publish`,
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ContentItem', id },
        { type: 'ContentItem', id: 'LIST' },
        { type: 'ContentVersion', id },
      ],
    }),

    archiveContent: builder.mutation<ContentItem, string>({
      query: (id) => ({
        url: `/content/items/${id}/archive`,
        method: 'POST',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'ContentItem', id },
        { type: 'ContentItem', id: 'LIST' },
      ],
    }),

    bulkUpdateContent: builder.mutation<ContentItem[], BulkUpdateContentRequest>({
      query: (data) => ({
        url: '/content/items/bulk',
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: [{ type: 'ContentItem', id: 'LIST' }],
    }),

    revertContentVersion: builder.mutation<ContentItem, { id: string; version: number }>({
      query: ({ id, version }) => ({
        url: `/content/items/${id}/revert/${version}`,
        method: 'POST',
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'ContentItem', id },
        { type: 'ContentVersion', id },
      ],
    }),

  }),
  overrideExisting: false,
});

// Export hooks
export const {
  useGetContentItemsQuery,
  useGetContentItemQuery,
  useGetContentByCategoryQuery,
  useGetContentCategoriesQuery,
  useGetContentVersionsQuery,
  useSyncContentQuery,
  useCreateContentItemMutation,
  useUpdateContentItemMutation,
  usePublishContentMutation,
  useArchiveContentMutation,
  useBulkUpdateContentMutation,
  useRevertContentVersionMutation,
} = contentApi;