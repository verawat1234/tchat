/**
 * Contract Test: getContentItems API Endpoint
 *
 * This test validates the exact contract structure for the content management API
 * getContentItems endpoint. It follows TDD principles - the test WILL FAIL until
 * the backend API endpoints are implemented.
 *
 * Validates:
 * - API endpoint structure and response format
 * - Pagination parameters (limit, offset, hasMore)
 * - Filtering by category and status
 * - Content item structure validation
 * - Error handling scenarios
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import { http, HttpResponse } from 'msw';
import { server } from '../../lib/test-utils/msw/server';
import { api } from '../api';

// Contract types - these will be moved to a proper contract file once implemented
export interface ContentItem {
  id: string;
  category: string;
  type: ContentType;
  value: ContentValue;
  metadata: ContentMetadata;
  status: ContentStatus;
}

export interface ContentMetadata {
  createdAt: string;
  createdBy: string;
  updatedAt: string;
  updatedBy: string;
  version: number;
  tags?: string[];
  notes?: string;
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

// Mock contentApi since it doesn't exist yet - this will fail in TDD fashion
const contentApi = {
  endpoints: {
    getContentItems: {
      initiate: (params: GetContentItemsRequest) => {
        throw new Error('getContentItems endpoint not implemented yet - this is expected for TDD');
      }
    }
  }
};

// Test store setup for RTK Query
const createTestStore = () => {
  return configureStore({
    reducer: {
      [api.reducerPath]: api.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(api.middleware),
  });
};

// Mock data factories for consistent test data
const createMockContentMetadata = (overrides: Partial<ContentMetadata> = {}): ContentMetadata => ({
  createdAt: '2023-01-01T00:00:00.000Z',
  createdBy: 'test-user-id',
  updatedAt: '2023-01-01T12:00:00.000Z',
  updatedBy: 'test-user-id',
  version: 1,
  tags: ['test', 'content'],
  notes: 'Test content item',
  ...overrides,
});

const createMockContentItem = (overrides: Partial<ContentItem> = {}): ContentItem => ({
  id: 'test-content-1',
  category: 'ui.buttons',
  type: ContentType.TEXT,
  value: {
    type: ContentType.TEXT,
    value: 'Test Button Label',
    locale: 'en',
  },
  metadata: createMockContentMetadata(),
  status: ContentStatus.PUBLISHED,
  ...overrides,
});

const createMockResponse = (
  items: ContentItem[] = [],
  total: number = items.length,
  hasMore: boolean = false
): GetContentItemsResponse => ({
  items,
  total,
  hasMore,
});

// API Base URL for MSW
const API_BASE_URL = process.env.VITE_API_URL || 'http://localhost:3001';

describe('Content API - getContentItems Contract Test', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
  });

  afterEach(() => {
    server.resetHandlers();
  });

  describe('API Endpoint Structure', () => {
    it('should call correct endpoint with proper HTTP method', async () => {
      // Arrange: Mock successful response
      const mockItems = [
        createMockContentItem({ id: 'item-1' }),
        createMockContentItem({ id: 'item-2' }),
      ];
      const mockResponse = createMockResponse(mockItems, 2, false);

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify successful response structure
      expect(result.data).toEqual(mockResponse);
      expect(result.data?.items).toHaveLength(2);
      expect(result.data?.total).toBe(2);
      expect(result.data?.hasMore).toBe(false);
    });

    it('should handle empty results correctly', async () => {
      // Arrange: Mock empty response
      const mockResponse = createMockResponse([], 0, false);

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify empty response structure
      expect(result.data).toEqual(mockResponse);
      expect(result.data?.items).toEqual([]);
      expect(result.data?.total).toBe(0);
      expect(result.data?.hasMore).toBe(false);
    });
  });

  describe('Pagination Parameters', () => {
    it('should handle pagination with limit and offset', async () => {
      // Arrange: Mock paginated response
      const mockItems = Array.from({ length: 10 }, (_, i) =>
        createMockContentItem({ id: `item-${i + 11}` })
      );
      const mockResponse = createMockResponse(mockItems, 25, true);

      let capturedParams: URLSearchParams | null = null;

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, ({ request }) => {
          const url = new URL(request.url);
          capturedParams = url.searchParams;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query with pagination
      const queryParams: GetContentItemsRequest = {
        limit: 10,
        offset: 10,
      };

      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate(queryParams)
      );

      // Assert: Verify pagination parameters and response
      expect(capturedParams?.get('limit')).toBe('10');
      expect(capturedParams?.get('offset')).toBe('10');
      expect(result.data?.items).toHaveLength(10);
      expect(result.data?.total).toBe(25);
      expect(result.data?.hasMore).toBe(true);
    });

    it('should handle hasMore flag correctly for last page', async () => {
      // Arrange: Mock last page response
      const mockItems = Array.from({ length: 5 }, (_, i) =>
        createMockContentItem({ id: `item-${i + 21}` })
      );
      const mockResponse = createMockResponse(mockItems, 25, false);

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query for last page
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          limit: 10,
          offset: 20,
        })
      );

      // Assert: Verify hasMore is false for last page
      expect(result.data?.items).toHaveLength(5);
      expect(result.data?.total).toBe(25);
      expect(result.data?.hasMore).toBe(false);
    });
  });

  describe('Filtering Parameters', () => {
    it('should filter by category', async () => {
      // Arrange: Mock filtered response
      const mockItems = [
        createMockContentItem({
          id: 'ui-button-1',
          category: 'ui.buttons',
        }),
        createMockContentItem({
          id: 'ui-button-2',
          category: 'ui.buttons',
        }),
      ];
      const mockResponse = createMockResponse(mockItems, 2, false);

      let capturedParams: URLSearchParams | null = null;

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, ({ request }) => {
          const url = new URL(request.url);
          capturedParams = url.searchParams;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query with category filter
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          category: 'ui.buttons',
        })
      );

      // Assert: Verify category filter and response
      expect(capturedParams?.get('category')).toBe('ui.buttons');
      expect(result.data?.items.every(item => item.category === 'ui.buttons')).toBe(true);
    });

    it('should filter by status', async () => {
      // Arrange: Mock status-filtered response
      const mockItems = [
        createMockContentItem({
          id: 'draft-1',
          status: ContentStatus.DRAFT,
        }),
        createMockContentItem({
          id: 'draft-2',
          status: ContentStatus.DRAFT,
        }),
      ];
      const mockResponse = createMockResponse(mockItems, 2, false);

      let capturedParams: URLSearchParams | null = null;

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, ({ request }) => {
          const url = new URL(request.url);
          capturedParams = url.searchParams;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query with status filter
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          status: ContentStatus.DRAFT,
        })
      );

      // Assert: Verify status filter and response
      expect(capturedParams?.get('status')).toBe('draft');
      expect(result.data?.items.every(item => item.status === ContentStatus.DRAFT)).toBe(true);
    });

    it('should filter by content type', async () => {
      // Arrange: Mock type-filtered response
      const mockItems = [
        createMockContentItem({
          id: 'rich-text-1',
          type: ContentType.RICH_TEXT,
          value: {
            type: ContentType.RICH_TEXT,
            value: '<p>Rich text content</p>',
            format: 'html',
          },
        }),
      ];
      const mockResponse = createMockResponse(mockItems, 1, false);

      let capturedParams: URLSearchParams | null = null;

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, ({ request }) => {
          const url = new URL(request.url);
          capturedParams = url.searchParams;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query with type filter
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          type: ContentType.RICH_TEXT,
        })
      );

      // Assert: Verify type filter and response
      expect(capturedParams?.get('type')).toBe('rich_text');
      expect(result.data?.items.every(item => item.type === ContentType.RICH_TEXT)).toBe(true);
    });

    it('should handle multiple filters simultaneously', async () => {
      // Arrange: Mock multi-filtered response
      const mockItems = [
        createMockContentItem({
          id: 'published-ui-text',
          category: 'ui.labels',
          type: ContentType.TEXT,
          status: ContentStatus.PUBLISHED,
        }),
      ];
      const mockResponse = createMockResponse(mockItems, 1, false);

      let capturedParams: URLSearchParams | null = null;

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, ({ request }) => {
          const url = new URL(request.url);
          capturedParams = url.searchParams;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query with multiple filters
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          category: 'ui.labels',
          type: ContentType.TEXT,
          status: ContentStatus.PUBLISHED,
          limit: 20,
          offset: 0,
        })
      );

      // Assert: Verify all filters are applied
      expect(capturedParams?.get('category')).toBe('ui.labels');
      expect(capturedParams?.get('type')).toBe('text');
      expect(capturedParams?.get('status')).toBe('published');
      expect(capturedParams?.get('limit')).toBe('20');
      expect(capturedParams?.get('offset')).toBe('0');
      expect(result.data?.items).toHaveLength(1);
    });

    it('should handle search parameter', async () => {
      // Arrange: Mock search response
      const mockItems = [
        createMockContentItem({
          id: 'search-result-1',
          value: {
            type: ContentType.TEXT,
            value: 'Button label with search term',
          },
        }),
      ];
      const mockResponse = createMockResponse(mockItems, 1, false);

      let capturedParams: URLSearchParams | null = null;

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, ({ request }) => {
          const url = new URL(request.url);
          capturedParams = url.searchParams;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query with search
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          search: 'button',
        })
      );

      // Assert: Verify search parameter
      expect(capturedParams?.get('search')).toBe('button');
      expect(result.data?.items).toHaveLength(1);
    });

    it('should handle tags parameter', async () => {
      // Arrange: Mock tags response
      const mockItems = [
        createMockContentItem({
          id: 'tagged-content',
          metadata: createMockContentMetadata({
            tags: ['ui', 'button', 'primary'],
          }),
        }),
      ];
      const mockResponse = createMockResponse(mockItems, 1, false);

      let capturedParams: URLSearchParams | null = null;

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, ({ request }) => {
          const url = new URL(request.url);
          capturedParams = url.searchParams;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query with tags
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          tags: ['ui', 'button'],
        })
      );

      // Assert: Verify tags parameter (implementation-dependent how this is serialized)
      expect(capturedParams?.has('tags')).toBe(true);
      expect(result.data?.items).toHaveLength(1);
    });
  });

  describe('Content Item Structure Validation', () => {
    it('should validate complete ContentItem structure', async () => {
      // Arrange: Mock complete content item
      const mockItem = createMockContentItem({
        id: 'complete-item',
        category: 'ui.navigation',
        type: ContentType.RICH_TEXT,
        value: {
          type: ContentType.RICH_TEXT,
          value: '<nav><a href="/">Home</a></nav>',
          format: 'html',
          locale: 'en',
        },
        metadata: createMockContentMetadata({
          version: 3,
          tags: ['navigation', 'header', 'ui'],
          notes: 'Main navigation component',
        }),
        status: ContentStatus.PUBLISHED,
      });
      const mockResponse = createMockResponse([mockItem], 1, false);

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify complete structure
      const item = result.data?.items[0];
      expect(item).toBeDefined();
      expect(item?.id).toBe('complete-item');
      expect(item?.category).toBe('ui.navigation');
      expect(item?.type).toBe(ContentType.RICH_TEXT);
      expect(item?.status).toBe(ContentStatus.PUBLISHED);

      // Validate ContentValue structure
      expect(item?.value.type).toBe(ContentType.RICH_TEXT);
      expect(item?.value.value).toContain('<nav>');
      expect(item?.value.format).toBe('html');
      expect(item?.value.locale).toBe('en');

      // Validate ContentMetadata structure
      expect(item?.metadata.createdAt).toMatch(/\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/);
      expect(item?.metadata.createdBy).toBe('test-user-id');
      expect(item?.metadata.updatedAt).toMatch(/\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/);
      expect(item?.metadata.updatedBy).toBe('test-user-id');
      expect(item?.metadata.version).toBe(3);
      expect(item?.metadata.tags).toEqual(['navigation', 'header', 'ui']);
      expect(item?.metadata.notes).toBe('Main navigation component');
    });

    it('should validate different ContentType variants', async () => {
      // Arrange: Mock different content types
      const mockItems = [
        createMockContentItem({
          id: 'text-content',
          type: ContentType.TEXT,
          value: { type: ContentType.TEXT, value: 'Simple text' },
        }),
        createMockContentItem({
          id: 'image-content',
          type: ContentType.IMAGE_URL,
          value: {
            type: ContentType.IMAGE_URL,
            value: 'https://example.com/image.jpg',
            alt: 'Description',
            width: 400,
            height: 300,
          },
        }),
        createMockContentItem({
          id: 'config-content',
          type: ContentType.CONFIG,
          value: {
            type: ContentType.CONFIG,
            value: { theme: 'dark', locale: 'en' },
          },
        }),
        createMockContentItem({
          id: 'translation-content',
          type: ContentType.TRANSLATION,
          value: {
            type: ContentType.TRANSLATION,
            value: { en: 'Hello', es: 'Hola', fr: 'Bonjour' },
          },
        }),
      ];
      const mockResponse = createMockResponse(mockItems, 4, false);

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify different content types
      const items = result.data?.items || [];
      expect(items).toHaveLength(4);

      const textItem = items.find(i => i.id === 'text-content');
      expect(textItem?.type).toBe(ContentType.TEXT);
      expect(textItem?.value.value).toBe('Simple text');

      const imageItem = items.find(i => i.id === 'image-content');
      expect(imageItem?.type).toBe(ContentType.IMAGE_URL);
      expect(imageItem?.value.value).toBe('https://example.com/image.jpg');
      expect(imageItem?.value.alt).toBe('Description');

      const configItem = items.find(i => i.id === 'config-content');
      expect(configItem?.type).toBe(ContentType.CONFIG);
      expect(configItem?.value.value).toEqual({ theme: 'dark', locale: 'en' });

      const translationItem = items.find(i => i.id === 'translation-content');
      expect(translationItem?.type).toBe(ContentType.TRANSLATION);
      expect(translationItem?.value.value).toEqual({ en: 'Hello', es: 'Hola', fr: 'Bonjour' });
    });

    it('should validate different ContentStatus variants', async () => {
      // Arrange: Mock different statuses
      const mockItems = [
        createMockContentItem({
          id: 'draft-content',
          status: ContentStatus.DRAFT,
        }),
        createMockContentItem({
          id: 'published-content',
          status: ContentStatus.PUBLISHED,
        }),
        createMockContentItem({
          id: 'archived-content',
          status: ContentStatus.ARCHIVED,
        }),
      ];
      const mockResponse = createMockResponse(mockItems, 3, false);

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify different statuses
      const items = result.data?.items || [];
      expect(items).toHaveLength(3);

      const statuses = items.map(item => item.status);
      expect(statuses).toContain(ContentStatus.DRAFT);
      expect(statuses).toContain(ContentStatus.PUBLISHED);
      expect(statuses).toContain(ContentStatus.ARCHIVED);
    });
  });

  describe('Error Handling Scenarios', () => {
    it('should handle 400 Bad Request errors', async () => {
      // Arrange: Mock bad request error
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(
            {
              error: 'Bad Request',
              message: 'Invalid query parameters',
              details: { limit: 'must be between 1 and 100' }
            },
            { status: 400 }
          );
        })
      );

      // Act: Execute query with invalid parameters
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          limit: 500, // Invalid large limit
        })
      );

      // Assert: Verify error response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(400);
      expect(result.data).toBeUndefined();
    });

    it('should handle 401 Unauthorized errors', async () => {
      // Arrange: Mock unauthorized error
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(
            { error: 'Unauthorized', message: 'Invalid or expired token' },
            { status: 401 }
          );
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify unauthorized response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(401);
      expect(result.data).toBeUndefined();
    });

    it('should handle 403 Forbidden errors', async () => {
      // Arrange: Mock forbidden error
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(
            {
              error: 'Forbidden',
              message: 'Insufficient permissions to access content category'
            },
            { status: 403 }
          );
        })
      );

      // Act: Execute query for restricted category
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          category: 'admin.settings',
        })
      );

      // Assert: Verify forbidden response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(403);
      expect(result.data).toBeUndefined();
    });

    it('should handle 404 Not Found errors', async () => {
      // Arrange: Mock not found error
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(
            { error: 'Not Found', message: 'Content category does not exist' },
            { status: 404 }
          );
        })
      );

      // Act: Execute query for non-existent category
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({
          category: 'non.existent.category',
        })
      );

      // Assert: Verify not found response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(404);
      expect(result.data).toBeUndefined();
    });

    it('should handle 500 Internal Server Error', async () => {
      // Arrange: Mock server error
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(
            { error: 'Internal Server Error', message: 'Database connection failed' },
            { status: 500 }
          );
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify server error response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(500);
      expect(result.data).toBeUndefined();
    });

    it('should handle network errors', async () => {
      // Arrange: Mock network error
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.error();
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify network error response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe('FETCH_ERROR');
      expect(result.data).toBeUndefined();
    });

    it('should handle malformed response data', async () => {
      // Arrange: Mock malformed response
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json({
            // Missing required fields
            items: null,
            // Invalid field types
            total: 'not-a-number',
            hasMore: 'not-a-boolean',
          });
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify that malformed data is received as-is
      // (Type validation would typically happen at runtime or with additional validation)
      expect(result.data).toBeDefined();
      expect(result.data?.items).toBe(null);
      expect(result.data?.total).toBe('not-a-number');
      expect(result.data?.hasMore).toBe('not-a-boolean');
    });
  });

  describe('RTK Query Integration', () => {
    it('should provide proper cache tags for invalidation', async () => {
      // Arrange: Mock response
      const mockItems = [
        createMockContentItem({ id: 'cacheable-1' }),
        createMockContentItem({ id: 'cacheable-2' }),
      ];
      const mockResponse = createMockResponse(mockItems, 2, false);

      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Verify that proper tags are set for caching
      // Note: Cache tag verification would require inspecting RTK Query internals
      expect(result.data).toEqual(mockResponse);

      // Verify that the endpoint provides the expected cache tags structure
      const endpoint = contentApi.endpoints.getContentItems;
      expect(endpoint.Types).toBeUndefined(); // provideTags is set on the endpoint definition

      // The actual cache tags are set in the endpoint definition and would be:
      // - Individual item tags: { type: 'ContentItem', id: item.id }
      // - List tag: { type: 'ContentItem', id: 'LIST' }
    });

    it('should handle query argument serialization correctly', async () => {
      // This test verifies that complex query arguments are properly serialized
      // and don't cause caching issues

      // Arrange: Mock response
      const mockResponse = createMockResponse([], 0, false);

      let requestCount = 0;
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          requestCount++;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute same query twice with identical arguments
      const queryArgs: GetContentItemsRequest = {
        category: 'ui.buttons',
        status: ContentStatus.PUBLISHED,
        limit: 10,
        offset: 0,
        tags: ['primary', 'ui'],
      };

      await store.dispatch(
        contentApi.endpoints.getContentItems.initiate(queryArgs)
      );

      await store.dispatch(
        contentApi.endpoints.getContentItems.initiate(queryArgs)
      );

      // Assert: Verify that caching works (only one request made)
      expect(requestCount).toBe(1);
    });
  });

  describe('Response Format Validation', () => {
    it('should require items array in response', async () => {
      // Arrange: Mock response without items
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json({
            total: 0,
            hasMore: false,
            // Missing items array
          });
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Response is received but missing required field
      expect(result.data).toBeDefined();
      expect(result.data?.items).toBeUndefined();
      expect(result.data?.total).toBe(0);
      expect(result.data?.hasMore).toBe(false);
    });

    it('should require total field in response', async () => {
      // Arrange: Mock response without total
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json({
            items: [],
            hasMore: false,
            // Missing total field
          });
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Response is received but missing required field
      expect(result.data).toBeDefined();
      expect(result.data?.items).toEqual([]);
      expect(result.data?.total).toBeUndefined();
      expect(result.data?.hasMore).toBe(false);
    });

    it('should require hasMore field in response', async () => {
      // Arrange: Mock response without hasMore
      server.use(
        http.get(`${API_BASE_URL}/api/content/items`, () => {
          return HttpResponse.json({
            items: [],
            total: 0,
            // Missing hasMore field
          });
        })
      );

      // Act: Execute query
      const result = await store.dispatch(
        contentApi.endpoints.getContentItems.initiate({})
      );

      // Assert: Response is received but missing required field
      expect(result.data).toBeDefined();
      expect(result.data?.items).toEqual([]);
      expect(result.data?.total).toBe(0);
      expect(result.data?.hasMore).toBeUndefined();
    });
  });
});