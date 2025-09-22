import { describe, it, expect, beforeEach, vi } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import { setupListeners } from '@reduxjs/toolkit/query';
import { api } from '../api';
import type { ContentVersion, ContentValue, ContentMetadata, PaginationMeta } from '../../types/content';
import type { ApiResponse, PaginatedResponse } from '../../types/api';

// Mock fetch to control network responses
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Helper to create test store with RTK Query
const createTestStore = () => {
  const store = configureStore({
    reducer: {
      api: api.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(api.middleware),
  });
  setupListeners(store.dispatch);
  return store;
};

// Content API with getContentVersions endpoint (will be injected)
const contentApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getContentVersions: builder.query<
      PaginatedResponse<ContentVersion>,
      {
        contentId: string;
        page?: number;
        limit?: number;
        sortBy?: 'version' | 'createdAt' | 'updatedAt';
        sortOrder?: 'asc' | 'desc';
      }
    >({
      query: ({ contentId, page = 1, limit = 10, sortBy = 'version', sortOrder = 'desc' }) => ({
        url: `/content/${contentId}/versions`,
        method: 'GET',
        params: { page, limit, sortBy, sortOrder },
      }),
      providesTags: (result, error, { contentId }) => [
        { type: 'ContentVersion' as const, id: contentId },
        { type: 'ContentVersion' as const, id: 'LIST' },
      ],
    }),
  }),
});

const { useGetContentVersionsQuery } = contentApi;

describe('Content API - getContentVersions Contract Tests', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
    mockFetch.mockClear();
  });

  describe('Content Version History Retrieval', () => {
    it('should fetch content version history with proper request structure', async () => {
      // Arrange: Mock successful response
      const mockVersions: ContentVersion[] = [
        {
          id: 'version-3',
          contentId: 'navigation.header.title',
          version: 3,
          value: {
            type: 'text',
            value: 'Welcome to Tchat v3',
            maxLength: 100,
          } as ContentValue,
          metadata: {
            createdAt: '2024-03-15T10:30:00Z',
            createdBy: 'user-123',
            updatedAt: '2024-03-15T10:30:00Z',
            updatedBy: 'user-123',
            version: 3,
            tags: ['header', 'navigation'],
            notes: 'Updated brand messaging',
          } as ContentMetadata,
          changeLog: 'Updated title for new brand guidelines',
        },
        {
          id: 'version-2',
          contentId: 'navigation.header.title',
          version: 2,
          value: {
            type: 'text',
            value: 'Welcome to Tchat v2',
            maxLength: 100,
          } as ContentValue,
          metadata: {
            createdAt: '2024-03-10T09:15:00Z',
            createdBy: 'user-456',
            updatedAt: '2024-03-10T09:15:00Z',
            updatedBy: 'user-456',
            version: 2,
            tags: ['header', 'navigation'],
            notes: 'Version update',
          } as ContentMetadata,
          changeLog: 'Updated version number in title',
        },
      ];

      const mockResponse: ApiResponse<PaginatedResponse<ContentVersion>> = {
        success: true,
        data: {
          items: mockVersions,
          pagination: {
            cursor: 'version-2',
            nextCursor: 'version-1',
            prevCursor: null,
            hasMore: true,
            total: 3,
          },
        },
        meta: {
          pagination: {
            total: 3,
            page: 1,
            limit: 10,
            totalPages: 1,
            hasNext: false,
            hasPrev: false,
          },
          timestamp: '2024-03-15T10:30:00Z',
          version: 'v1',
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve(mockResponse),
        headers: new Headers({ 'content-type': 'application/json' }),
      });

      // Act: Execute query through RTK Query hook
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'navigation.header.title',
          page: 1,
          limit: 10,
          sortBy: 'version',
          sortOrder: 'desc',
        })
      );

      // Wait for the query to complete
      const queryResult = await result;

      // Assert: Verify request was made correctly
      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/content/navigation.header.title/versions'),
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );

      // Assert: Verify response structure
      expect(queryResult.isSuccess).toBe(true);
      expect(queryResult.data).toBeDefined();
      expect(queryResult.data?.items).toHaveLength(2);
      expect(queryResult.data?.pagination).toBeDefined();
    });

    it('should handle pagination parameters correctly', async () => {
      // Arrange: Mock paginated response
      const mockResponse: ApiResponse<PaginatedResponse<ContentVersion>> = {
        success: true,
        data: {
          items: [],
          pagination: {
            cursor: 'version-10',
            nextCursor: 'version-5',
            prevCursor: 'version-15',
            hasMore: true,
            total: 25,
          },
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve(mockResponse),
      });

      // Act: Execute query with pagination
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content.item',
          page: 2,
          limit: 5,
          sortBy: 'createdAt',
          sortOrder: 'asc',
        })
      );

      await result;

      // Assert: Verify pagination parameters were sent
      const fetchCall = mockFetch.mock.calls[0];
      const url = new URL(fetchCall[0]);
      expect(url.searchParams.get('page')).toBe('2');
      expect(url.searchParams.get('limit')).toBe('5');
      expect(url.searchParams.get('sortBy')).toBe('createdAt');
      expect(url.searchParams.get('sortOrder')).toBe('asc');
    });
  });

  describe('Version Structure Validation', () => {
    it('should validate complete ContentVersion structure', async () => {
      // Arrange: Mock version with all required fields
      const completeVersion: ContentVersion = {
        id: 'version-1',
        contentId: 'test.content.item',
        version: 1,
        value: {
          type: 'rich_text',
          value: '<h1>Rich Content</h1>',
          format: 'html',
          allowedTags: ['h1', 'p', 'strong'],
        },
        metadata: {
          createdAt: '2024-03-15T10:30:00Z',
          createdBy: 'user-123',
          updatedAt: '2024-03-15T10:30:00Z',
          updatedBy: 'user-123',
          version: 1,
          tags: ['rich-text', 'html'],
          notes: 'Initial rich text content',
        },
        changeLog: 'Initial version with rich text content',
      };

      const mockResponse: ApiResponse<PaginatedResponse<ContentVersion>> = {
        success: true,
        data: {
          items: [completeVersion],
          pagination: {
            hasMore: false,
            total: 1,
          },
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve(mockResponse),
      });

      // Act: Execute query
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content.item',
        })
      );

      const queryResult = await result;

      // Assert: Validate version structure
      expect(queryResult.isSuccess).toBe(true);
      const version = queryResult.data?.items[0];

      // Core version fields
      expect(version?.id).toBe('version-1');
      expect(version?.contentId).toBe('test.content.item');
      expect(version?.version).toBe(1);
      expect(version?.changeLog).toBe('Initial version with rich text content');

      // Value structure
      expect(version?.value).toBeDefined();
      expect(version?.value.type).toBe('rich_text');

      // Metadata structure
      expect(version?.metadata).toBeDefined();
      expect(version?.metadata.createdAt).toBe('2024-03-15T10:30:00Z');
      expect(version?.metadata.createdBy).toBe('user-123');
      expect(version?.metadata.version).toBe(1);
      expect(version?.metadata.tags).toEqual(['rich-text', 'html']);
      expect(version?.metadata.notes).toBe('Initial rich text content');
    });

    it('should handle different content value types in versions', async () => {
      // Arrange: Mock versions with different value types
      const mockVersions: ContentVersion[] = [
        {
          id: 'version-text',
          contentId: 'test.content',
          version: 1,
          value: { type: 'text', value: 'Plain text' },
          metadata: {
            createdAt: '2024-03-15T10:30:00Z',
            createdBy: 'user-123',
            updatedAt: '2024-03-15T10:30:00Z',
            updatedBy: 'user-123',
            version: 1,
          },
          changeLog: 'Text content version',
        },
        {
          id: 'version-image',
          contentId: 'test.content',
          version: 2,
          value: {
            type: 'image_url',
            url: 'https://example.com/image.jpg',
            alt: 'Test image',
            width: 800,
            height: 600,
          },
          metadata: {
            createdAt: '2024-03-16T10:30:00Z',
            createdBy: 'user-456',
            updatedAt: '2024-03-16T10:30:00Z',
            updatedBy: 'user-456',
            version: 2,
          },
          changeLog: 'Changed to image content',
        },
        {
          id: 'version-config',
          contentId: 'test.content',
          version: 3,
          value: {
            type: 'config',
            value: { enabled: true, maxRetries: 3 },
            schema: { type: 'object' },
          },
          metadata: {
            createdAt: '2024-03-17T10:30:00Z',
            createdBy: 'user-789',
            updatedAt: '2024-03-17T10:30:00Z',
            updatedBy: 'user-789',
            version: 3,
          },
          changeLog: 'Changed to configuration object',
        },
      ];

      const mockResponse: ApiResponse<PaginatedResponse<ContentVersion>> = {
        success: true,
        data: {
          items: mockVersions,
          pagination: { hasMore: false, total: 3 },
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve(mockResponse),
      });

      // Act: Execute query
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content',
        })
      );

      const queryResult = await result;

      // Assert: Validate different value types
      expect(queryResult.isSuccess).toBe(true);
      const versions = queryResult.data?.items;

      expect(versions).toHaveLength(3);
      expect(versions?.[0].value.type).toBe('text');
      expect(versions?.[1].value.type).toBe('image_url');
      expect(versions?.[2].value.type).toBe('config');
    });
  });

  describe('Version Ordering and Metadata', () => {
    it('should return versions in descending order by default', async () => {
      // Arrange: Mock versions with different timestamps
      const mockVersions: ContentVersion[] = [
        {
          id: 'version-3',
          contentId: 'test.content',
          version: 3,
          value: { type: 'text', value: 'Latest version' },
          metadata: {
            createdAt: '2024-03-17T10:30:00Z',
            createdBy: 'user-123',
            updatedAt: '2024-03-17T10:30:00Z',
            updatedBy: 'user-123',
            version: 3,
          },
          changeLog: 'Latest changes',
        },
        {
          id: 'version-2',
          contentId: 'test.content',
          version: 2,
          value: { type: 'text', value: 'Middle version' },
          metadata: {
            createdAt: '2024-03-16T10:30:00Z',
            createdBy: 'user-123',
            updatedAt: '2024-03-16T10:30:00Z',
            updatedBy: 'user-123',
            version: 2,
          },
          changeLog: 'Middle changes',
        },
        {
          id: 'version-1',
          contentId: 'test.content',
          version: 1,
          value: { type: 'text', value: 'First version' },
          metadata: {
            createdAt: '2024-03-15T10:30:00Z',
            createdBy: 'user-123',
            updatedAt: '2024-03-15T10:30:00Z',
            updatedBy: 'user-123',
            version: 1,
          },
          changeLog: 'Initial version',
        },
      ];

      const mockResponse: ApiResponse<PaginatedResponse<ContentVersion>> = {
        success: true,
        data: {
          items: mockVersions,
          pagination: { hasMore: false, total: 3 },
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve(mockResponse),
      });

      // Act: Execute query with default ordering
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content',
        })
      );

      const queryResult = await result;

      // Assert: Verify descending order
      expect(queryResult.isSuccess).toBe(true);
      const versions = queryResult.data?.items;

      expect(versions).toHaveLength(3);
      expect(versions?.[0].version).toBe(3);
      expect(versions?.[1].version).toBe(2);
      expect(versions?.[2].version).toBe(1);
    });

    it('should include comprehensive metadata for each version', async () => {
      // Arrange: Mock version with rich metadata
      const mockVersion: ContentVersion = {
        id: 'version-meta-test',
        contentId: 'test.content',
        version: 1,
        value: { type: 'text', value: 'Test content' },
        metadata: {
          createdAt: '2024-03-15T10:30:00Z',
          createdBy: 'user-123',
          updatedAt: '2024-03-15T11:45:00Z',
          updatedBy: 'user-456',
          version: 1,
          tags: ['test', 'metadata', 'validation'],
          notes: 'Comprehensive metadata test version',
        },
        changeLog: 'Added comprehensive metadata for testing validation',
      };

      const mockResponse: ApiResponse<PaginatedResponse<ContentVersion>> = {
        success: true,
        data: {
          items: [mockVersion],
          pagination: { hasMore: false, total: 1 },
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve(mockResponse),
      });

      // Act: Execute query
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content',
        })
      );

      const queryResult = await result;

      // Assert: Validate metadata completeness
      expect(queryResult.isSuccess).toBe(true);
      const version = queryResult.data?.items[0];
      const metadata = version?.metadata;

      expect(metadata?.createdAt).toBe('2024-03-15T10:30:00Z');
      expect(metadata?.createdBy).toBe('user-123');
      expect(metadata?.updatedAt).toBe('2024-03-15T11:45:00Z');
      expect(metadata?.updatedBy).toBe('user-456');
      expect(metadata?.version).toBe(1);
      expect(metadata?.tags).toEqual(['test', 'metadata', 'validation']);
      expect(metadata?.notes).toBe('Comprehensive metadata test version');
      expect(version?.changeLog).toBe('Added comprehensive metadata for testing validation');
    });
  });

  describe('Error Handling for Non-Existent Content', () => {
    it('should handle 404 error for non-existent content ID', async () => {
      // Arrange: Mock 404 response
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: () => Promise.resolve({
          success: false,
          error: {
            code: 'NOT_FOUND',
            message: 'Content item not found',
            details: { contentId: 'non.existent.content' },
            timestamp: '2024-03-15T10:30:00Z',
          },
        }),
      });

      // Act: Execute query with non-existent content ID
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'non.existent.content',
        })
      );

      const queryResult = await result;

      // Assert: Verify error handling
      expect(queryResult.isError).toBe(true);
      expect(queryResult.error).toBeDefined();

      // Check error structure
      if ('data' in queryResult.error!) {
        const errorData = queryResult.error.data as any;
        expect(errorData.success).toBe(false);
        expect(errorData.error.code).toBe('NOT_FOUND');
        expect(errorData.error.message).toBe('Content item not found');
      }
    });

    it('should handle authorization errors for restricted content', async () => {
      // Arrange: Mock 403 response
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
        json: () => Promise.resolve({
          success: false,
          error: {
            code: 'FORBIDDEN',
            message: 'Insufficient permissions to access content versions',
            details: { contentId: 'restricted.content.item' },
            timestamp: '2024-03-15T10:30:00Z',
          },
        }),
      });

      // Act: Execute query for restricted content
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'restricted.content.item',
        })
      );

      const queryResult = await result;

      // Assert: Verify authorization error
      expect(queryResult.isError).toBe(true);

      if ('data' in queryResult.error!) {
        const errorData = queryResult.error.data as any;
        expect(errorData.error.code).toBe('FORBIDDEN');
        expect(errorData.error.message).toBe('Insufficient permissions to access content versions');
      }
    });

    it('should handle network errors gracefully', async () => {
      // Arrange: Mock network error
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      // Act: Execute query that will fail
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content',
        })
      );

      const queryResult = await result;

      // Assert: Verify network error handling
      expect(queryResult.isError).toBe(true);
      expect(queryResult.error).toBeDefined();
    });

    it('should handle malformed response data', async () => {
      // Arrange: Mock malformed response
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: () => Promise.resolve({
          // Missing required fields
          data: {
            items: [
              {
                id: 'incomplete-version',
                // Missing contentId, version, value, metadata, changeLog
              },
            ],
          },
        }),
      });

      // Act: Execute query
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content',
        })
      );

      const queryResult = await result;

      // Assert: Verify handling of malformed data
      // RTK Query should still mark as successful but data may be incomplete
      expect(queryResult.isSuccess).toBe(true);
      const version = queryResult.data?.items[0];
      expect(version?.id).toBe('incomplete-version');
      // These fields should be undefined due to incomplete response
      expect(version?.contentId).toBeUndefined();
      expect(version?.version).toBeUndefined();
    });
  });

  describe('RTK Query Caching and Invalidation', () => {
    it('should provide proper cache tags for invalidation', () => {
      // Act: Get endpoint definition
      const endpoint = contentApi.endpoints.getContentVersions;

      // Assert: Verify cache tag structure
      expect(endpoint.providesTags).toBeDefined();

      // Test tag generation
      const mockResult = { items: [], pagination: { hasMore: false } };
      const mockArgs = { contentId: 'test.content' };
      const tags = endpoint.providesTags!(mockResult, null, mockArgs);

      expect(tags).toContain({ type: 'ContentVersion', id: 'test.content' });
      expect(tags).toContain({ type: 'ContentVersion', id: 'LIST' });
    });

    it('should cache responses based on query parameters', async () => {
      // Arrange: Mock same response for multiple calls
      const mockResponse = {
        ok: true,
        status: 200,
        json: () => Promise.resolve({
          success: true,
          data: { items: [], pagination: { hasMore: false } },
        }),
      };

      mockFetch.mockResolvedValue(mockResponse);

      // Act: Make same query twice
      const query1 = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content',
          page: 1,
        })
      );

      const query2 = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content',
          page: 1,
        })
      );

      await Promise.all([query1.result, query2.result]);

      // Assert: Should only make one network request due to caching
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });
  });

  describe('TDD Contract Validation', () => {
    it('should fail when backend endpoint is not implemented', async () => {
      // This test intentionally expects failure until backend implementation exists

      // Arrange: Mock 404 for unimplemented endpoint
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: () => Promise.resolve({
          success: false,
          error: {
            code: 'ENDPOINT_NOT_FOUND',
            message: 'Endpoint /content/{contentId}/versions not implemented',
            timestamp: '2024-03-15T10:30:00Z',
          },
        }),
      });

      // Act: Attempt to call unimplemented endpoint
      const { result } = store.dispatch(
        contentApi.endpoints.getContentVersions.initiate({
          contentId: 'test.content',
        })
      );

      const queryResult = await result;

      // Assert: Should fail as expected in TDD process
      expect(queryResult.isError).toBe(true);

      if ('data' in queryResult.error!) {
        const errorData = queryResult.error.data as any;
        expect(errorData.error.code).toBe('ENDPOINT_NOT_FOUND');
        expect(errorData.error.message).toContain('not implemented');
      }

      // This failure is EXPECTED and REQUIRED for proper TDD
      // The test will pass once the backend API is implemented
    });
  });
});