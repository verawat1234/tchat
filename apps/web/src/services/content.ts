import { api } from './api';
import type { ContentItem } from '../types/content';

/**
 * Content API Service
 *
 * RTK Query service for managing dynamic content items.
 * This service provides endpoints for CRUD operations on content items,
 * supporting the content management system with proper caching and type safety.
 *
 * IMPLEMENTATION STATUS: STUB - Endpoints not yet implemented
 * This file exists to support T005 contract test, which MUST FAIL until implementation is complete.
 */

export const contentApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getContentItem: builder.query<ContentItem, string>({
      query: (contentId) => ({
        url: `/content/${contentId}`,
        method: 'GET',
      }),
      providesTags: (result, error, contentId) => [
        { type: 'Content', id: contentId },
        { type: 'Content', id: 'LIST' },
      ],
    }),
  }),
  overrideExisting: false,
});

// Export hooks for the content API endpoints
export const {
  useGetContentItemQuery,
} = contentApi;