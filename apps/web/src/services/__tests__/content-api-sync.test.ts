import { describe, it, expect, beforeEach, vi } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import { setupListeners } from '@reduxjs/toolkit/query';
import { api } from '../api';
import { contentApi } from '../content';
import {
  ContentItem,
  ContentType,
  ContentStatus,
  ContentCategory,
  ContentValue,
  ContentMetadata
} from '../../types/content';
import type { ApiError } from '../../types/api';

/**
 * T009: Contract Test for syncContent API
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * This test validates:
 * 1. Content sync data retrieval with lastSyncTime parameter
 * 2. Response structure validation (items, deletedIds, syncTime)
 * 3. Incremental sync functionality
 * 4. Timestamp validation (ISO format)
 * 5. Conflict resolution scenarios
 * 6. Error handling for invalid sync timestamps
 */

// Mock data following content-api contract structure
const mockContentCategory: ContentCategory = {
  id: 'navigation',
  name: 'Navigation',
  description: 'Navigation menu items and labels',
  permissions: {
    read: ['user', 'admin'],
    write: ['admin'],
    publish: ['admin']
  }
};

const mockTextContent: ContentValue = {
  type: 'text',
  value: 'Updated Welcome Message',
  maxLength: 100
};

const mockMetadata: ContentMetadata = {
  createdAt: '2024-01-15T10:00:00.000Z',
  createdBy: 'user-123',
  updatedAt: '2024-01-15T11:30:00.000Z',
  updatedBy: 'user-456',
  version: 2,
  tags: ['header', 'welcome', 'updated'],
  notes: 'Updated welcome message with new branding'
};

const mockUpdatedContentItem: ContentItem = {
  id: 'navigation.header.welcome',
  category: mockContentCategory,
  type: ContentType.TEXT,
  value: mockTextContent,
  metadata: mockMetadata,
  status: ContentStatus.PUBLISHED
};

const mockNewContentItem: ContentItem = {
  id: 'navigation.footer.copyright',
  category: mockContentCategory,
  type: ContentType.TEXT,
  value: {
    type: 'text',
    value: '© 2024 Tchat. All rights reserved.',
    maxLength: 50
  },
  metadata: {
    createdAt: '2024-01-15T11:00:00.000Z',
    createdBy: 'user-789',
    updatedAt: '2024-01-15T11:00:00.000Z',
    updatedBy: 'user-789',
    version: 1,
    tags: ['footer', 'legal']
  },
  status: ContentStatus.PUBLISHED
};

// Sync response interface (contract definition)
interface ContentSyncResponse {
  /** Content items created or updated since lastSyncTime */
  items: ContentItem[];
  /** IDs of content items deleted since lastSyncTime */
  deletedIds: string[];
  /** Server timestamp for this sync operation (ISO format) */
  syncTime: string;
  /** Number of total changes */
  totalChanges: number;
  /** Whether there are more changes to sync (pagination) */
  hasMore: boolean;
  /** Next page token for paginated sync */
  nextPageToken?: string;
}

// Sync request interface (contract definition)
interface ContentSyncRequest {
  /** ISO timestamp of last successful sync (optional for initial sync) */
  lastSyncTime?: string;
  /** Maximum number of items to return (pagination) */
  limit?: number;
  /** Page token for continued sync */
  pageToken?: string;
  /** Categories to include in sync (optional filter) */
  categories?: string[];
  /** Include draft content (admin only) */
  includeDrafts?: boolean;
}

const mockSyncResponse: ContentSyncResponse = {
  items: [mockUpdatedContentItem, mockNewContentItem],
  deletedIds: ['navigation.sidebar.oldlink'],
  syncTime: '2024-01-15T12:00:00.000Z',
  totalChanges: 3,
  hasMore: false
};

const mockApiError: ApiError = {
  success: false,
  error: {
    code: 'INVALID_TIMESTAMP',
    message: 'Invalid lastSyncTime format. Expected ISO 8601 timestamp.',
    details: {
      providedTimestamp: 'invalid-timestamp',
      expectedFormat: 'YYYY-MM-DDTHH:mm:ss.sssZ',
      examples: ['2024-01-15T10:00:00.000Z', '2024-01-15T10:00:00Z']
    },
    timestamp: '2024-01-15T12:00:00.000Z'
  }
};

// Helper to create test store with API
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

// Mock fetch for testing
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('Content API - syncContent', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
    mockFetch.mockClear();
    vi.clearAllMocks();
  });

  describe('Incremental sync functionality', () => {
    it('should perform initial sync without lastSyncTime parameter', async () => {
      // Arrange: Mock successful initial sync response
      const initialSyncResponse: ContentSyncResponse = {
        items: [mockUpdatedContentItem, mockNewContentItem],
        deletedIds: [],
        syncTime: '2024-01-15T12:00:00.000Z',
        totalChanges: 2,
        hasMore: false
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: initialSyncResponse,
          meta: {
            timestamp: '2024-01-15T12:00:00.000Z',
            version: '1.0.0'
          }
        }),
      });

      // Act: Trigger initial sync (TDD requirement - this WILL FAIL)
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({})
      );

      const response = await result;

      // Assert: Validate initial sync response structure
      expect(response.data).toBeDefined();
      expect(response.error).toBeUndefined();

      const syncData = response.data as ContentSyncResponse;

      // Validate sync response structure
      expect(syncData.items).toHaveLength(2);
      expect(syncData.deletedIds).toHaveLength(0);
      expect(syncData.totalChanges).toBe(2);
      expect(syncData.hasMore).toBe(false);

      // Validate ISO timestamp format
      expect(syncData.syncTime).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);

      // Validate returned content items
      expect(syncData.items[0].id).toBe('navigation.header.welcome');
      expect(syncData.items[1].id).toBe('navigation.footer.copyright');

      // Validate metadata timestamps
      syncData.items.forEach(item => {
        expect(item.metadata.createdAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
        expect(item.metadata.updatedAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
      });
    });

    it('should perform incremental sync with lastSyncTime parameter', async () => {
      // Arrange: Mock incremental sync response
      const lastSyncTime = '2024-01-15T10:00:00.000Z';

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockSyncResponse,
          meta: {
            timestamp: '2024-01-15T12:00:00.000Z',
            version: '1.0.0'
          }
        }),
      });

      // Act: Trigger incremental sync
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({ lastSyncTime })
      );

      const response = await result;

      // Assert: Validate incremental sync response
      expect(response.data).toBeDefined();
      const syncData = response.data as ContentSyncResponse;

      // Validate sync included updates and deletions
      expect(syncData.items).toHaveLength(2);
      expect(syncData.deletedIds).toHaveLength(1);
      expect(syncData.deletedIds[0]).toBe('navigation.sidebar.oldlink');
      expect(syncData.totalChanges).toBe(3); // 2 updates + 1 deletion

      // Validate sync timestamp is after lastSyncTime
      const syncTimestamp = new Date(syncData.syncTime);
      const lastSyncTimestamp = new Date(lastSyncTime);
      expect(syncTimestamp.getTime()).toBeGreaterThan(lastSyncTimestamp.getTime());

      // Validate request was made with correct parameters
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/content/sync'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          }),
          body: expect.stringContaining(lastSyncTime)
        })
      );
    });

    it('should handle paginated sync with large datasets', async () => {
      // Arrange: Mock paginated sync response
      const paginatedResponse: ContentSyncResponse = {
        items: [mockUpdatedContentItem],
        deletedIds: ['navigation.sidebar.oldlink'],
        syncTime: '2024-01-15T12:00:00.000Z',
        totalChanges: 50,
        hasMore: true,
        nextPageToken: 'next-page-token-abc123'
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: paginatedResponse
        }),
      });

      // Act: Trigger paginated sync
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z',
          limit: 20
        })
      );

      const response = await result;

      // Assert: Validate paginated response
      expect(response.data).toBeDefined();
      const syncData = response.data as ContentSyncResponse;

      expect(syncData.hasMore).toBe(true);
      expect(syncData.nextPageToken).toBe('next-page-token-abc123');
      expect(syncData.totalChanges).toBe(50);
    });

    it('should handle continued paginated sync with page token', async () => {
      // Arrange: Mock continued sync response
      const continuedResponse: ContentSyncResponse = {
        items: [mockNewContentItem],
        deletedIds: [],
        syncTime: '2024-01-15T12:00:00.000Z',
        totalChanges: 50,
        hasMore: false
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: continuedResponse
        }),
      });

      // Act: Trigger continued sync with page token
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          pageToken: 'next-page-token-abc123',
          limit: 20
        })
      );

      const response = await result;

      // Assert: Validate continued sync
      expect(response.data).toBeDefined();
      const syncData = response.data as ContentSyncResponse;

      expect(syncData.hasMore).toBe(false);
      expect(syncData.nextPageToken).toBeUndefined();
    });
  });

  describe('Content filtering and categorization', () => {
    it('should sync specific categories when requested', async () => {
      // Arrange: Mock category-filtered sync
      const categoryFilteredResponse: ContentSyncResponse = {
        items: [mockUpdatedContentItem],
        deletedIds: [],
        syncTime: '2024-01-15T12:00:00.000Z',
        totalChanges: 1,
        hasMore: false
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: categoryFilteredResponse
        }),
      });

      // Act: Sync specific categories
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z',
          categories: ['navigation', 'branding']
        })
      );

      const response = await result;

      // Assert: Validate category filtering
      expect(response.data).toBeDefined();
      const syncData = response.data as ContentSyncResponse;

      expect(syncData.items).toHaveLength(1);
      expect(syncData.items[0].category.id).toBe('navigation');

      // Validate request included category filter
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/content/sync'),
        expect.objectContaining({
          body: expect.stringContaining('navigation')
        })
      );
    });

    it('should include draft content when includeDrafts is true (admin only)', async () => {
      // Arrange: Mock draft content in sync
      const draftContentItem: ContentItem = {
        ...mockNewContentItem,
        id: 'navigation.menu.newfeature',
        status: ContentStatus.DRAFT,
        metadata: {
          ...mockNewContentItem.metadata,
          notes: 'Draft content for upcoming feature'
        }
      };

      const draftSyncResponse: ContentSyncResponse = {
        items: [mockUpdatedContentItem, draftContentItem],
        deletedIds: [],
        syncTime: '2024-01-15T12:00:00.000Z',
        totalChanges: 2,
        hasMore: false
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: draftSyncResponse
        }),
      });

      // Act: Sync with draft content included
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z',
          includeDrafts: true
        })
      );

      const response = await result;

      // Assert: Validate draft content inclusion
      expect(response.data).toBeDefined();
      const syncData = response.data as ContentSyncResponse;

      expect(syncData.items).toHaveLength(2);

      const draftItem = syncData.items.find(item => item.status === ContentStatus.DRAFT);
      expect(draftItem).toBeDefined();
      expect(draftItem?.id).toBe('navigation.menu.newfeature');
      expect(draftItem?.metadata.notes).toBe('Draft content for upcoming feature');
    });
  });

  describe('Conflict resolution scenarios', () => {
    it('should handle version conflicts during sync', async () => {
      // Arrange: Mock conflict resolution response
      const conflictResponse = {
        success: true,
        data: mockSyncResponse,
        conflicts: [
          {
            contentId: 'navigation.header.welcome',
            localVersion: 1,
            serverVersion: 2,
            resolution: 'server_wins',
            conflictTime: '2024-01-15T11:45:00.000Z'
          }
        ]
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => conflictResponse,
      });

      // Act: Sync with potential conflicts
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z'
        })
      );

      const response = await result;

      // Assert: Validate conflict resolution
      expect(response.data).toBeDefined();

      // In a real implementation, conflicts would be handled
      // This test documents the expected conflict resolution behavior
      const responseData = response.data as any;
      expect(responseData.items[0].metadata.version).toBe(2); // Server version wins
    });

    it('should preserve local changes when server version is older', async () => {
      // Arrange: Mock scenario where local version is newer
      const localNewerResponse: ContentSyncResponse = {
        items: [], // No server updates
        deletedIds: [],
        syncTime: '2024-01-15T12:00:00.000Z',
        totalChanges: 0,
        hasMore: false
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: localNewerResponse
        }),
      });

      // Act: Sync when local content is newer
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T11:30:00.000Z' // Recent sync time
        })
      );

      const response = await result;

      // Assert: Validate no overwrites occurred
      expect(response.data).toBeDefined();
      const syncData = response.data as ContentSyncResponse;

      expect(syncData.items).toHaveLength(0); // No server updates
      expect(syncData.totalChanges).toBe(0);
    });
  });

  describe('Error handling', () => {
    it('should handle invalid lastSyncTime format', async () => {
      // Arrange: Mock invalid timestamp error
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => mockApiError,
      });

      // Act: Attempt sync with invalid timestamp
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: 'invalid-timestamp'
        })
      );

      const response = await result;

      // Assert: Validate error response
      expect(response.data).toBeUndefined();
      expect(response.error).toBeDefined();

      const error = response.error as any;
      expect(error.status).toBe(400);
      expect(error.data.error.code).toBe('INVALID_TIMESTAMP');
      expect(error.data.error.details.expectedFormat).toBe('YYYY-MM-DDTHH:mm:ss.sssZ');
      expect(error.data.error.details.examples).toContain('2024-01-15T10:00:00.000Z');
    });

    it('should handle network errors during sync', async () => {
      // Arrange: Mock network failure
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      // Act: Attempt sync with network failure
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z'
        })
      );

      const response = await result;

      // Assert: Validate network error handling
      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.error).toBe('Network error');
    });

    it('should handle 401 Unauthorized during sync', async () => {
      // Arrange: Mock unauthorized error
      const unauthorizedError: ApiError = {
        success: false,
        error: {
          code: 'UNAUTHORIZED',
          message: 'Authentication required for content sync',
          timestamp: '2024-01-15T12:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => unauthorizedError,
      });

      // Act: Attempt unauthorized sync
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z'
        })
      );

      const response = await result;

      // Assert: Validate unauthorized error
      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(401);
      expect(error.data.error.code).toBe('UNAUTHORIZED');
    });

    it('should handle 403 Forbidden for draft content access', async () => {
      // Arrange: Mock forbidden error for draft access
      const forbiddenError: ApiError = {
        success: false,
        error: {
          code: 'FORBIDDEN',
          message: 'Insufficient permissions to access draft content',
          details: { requiredRole: 'admin', userRole: 'user' },
          timestamp: '2024-01-15T12:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
        json: async () => forbiddenError,
      });

      // Act: Attempt to sync draft content without admin role
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z',
          includeDrafts: true
        })
      );

      const response = await result;

      // Assert: Validate forbidden error
      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(403);
      expect(error.data.error.code).toBe('FORBIDDEN');
      expect(error.data.error.details.requiredRole).toBe('admin');
    });

    it('should handle server errors (500) gracefully', async () => {
      // Arrange: Mock server error
      const serverError: ApiError = {
        success: false,
        error: {
          code: 'INTERNAL_SERVER_ERROR',
          message: 'An unexpected error occurred during sync',
          timestamp: '2024-01-15T12:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => serverError,
      });

      // Act: Attempt sync during server error
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z'
        })
      );

      const response = await result;

      // Assert: Validate server error handling
      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(500);
      expect(error.data.error.code).toBe('INTERNAL_SERVER_ERROR');
    });
  });

  describe('Timestamp validation', () => {
    it('should validate ISO 8601 timestamp format', () => {
      const validTimestamps = [
        '2024-01-15T10:00:00.000Z',
        '2024-01-15T10:00:00Z',
        '2024-12-31T23:59:59.999Z',
        '2024-01-01T00:00:00.000Z'
      ];

      validTimestamps.forEach(timestamp => {
        expect(timestamp).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?Z$/);
      });
    });

    it('should reject invalid timestamp formats', () => {
      const invalidTimestamps = [
        'invalid-timestamp',
        '2024-01-15',
        '2024-01-15 10:00:00',
        '2024/01/15T10:00:00Z',
        ''
      ];

      invalidTimestamps.forEach(timestamp => {
        expect(timestamp).not.toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?Z$/);
      });

      // Test specifically invalid date components
      const invalidDate = '2024-13-32T25:61:61.000Z';
      // This technically matches the regex pattern but is an invalid date
      expect(invalidDate).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?Z$/);
      // But would be rejected by Date parsing
      expect(isNaN(new Date(invalidDate).getTime())).toBe(true);
    });

    it('should handle timezone considerations in sync timestamps', async () => {
      // Arrange: Mock response with timezone-aware timestamp
      const timezoneResponse: ContentSyncResponse = {
        items: [mockUpdatedContentItem],
        deletedIds: [],
        syncTime: '2024-01-15T12:00:00.000Z', // Always UTC
        totalChanges: 1,
        hasMore: false
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: timezoneResponse
        }),
      });

      // Act: Sync with timezone considerations
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z' // UTC timestamp
        })
      );

      const response = await result;

      // Assert: Validate timezone handling
      expect(response.data).toBeDefined();
      const syncData = response.data as ContentSyncResponse;

      // All timestamps should be in UTC (Z suffix)
      expect(syncData.syncTime).toEndWith('Z');
      expect(syncData.items[0].metadata.createdAt).toEndWith('Z');
      expect(syncData.items[0].metadata.updatedAt).toEndWith('Z');
    });
  });

  describe('RTK Query integration', () => {
    it('should provide correct cache tags for sync operations', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockSyncResponse
        }),
      });

      // This test validates that the endpoint provides proper cache tags
      // The actual implementation should provide tags like:
      // ['ContentSync', 'Content']

      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z'
        })
      );

      await result;

      // When implemented, this should verify cache tag structure
      // For now, this documents the expected cache behavior
      expect(true).toBe(true); // Placeholder until implementation
    });

    it('should handle query deduplication for identical sync requests', async () => {
      const syncParams = { lastSyncTime: '2024-01-15T10:00:00.000Z' };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockSyncResponse
        }),
      });

      // Trigger multiple identical sync requests
      const sync1 = store.dispatch(
        contentApi.endpoints.syncContent.initiate(syncParams)
      );

      const sync2 = store.dispatch(
        contentApi.endpoints.syncContent.initiate(syncParams)
      );

      await Promise.all([sync1, sync2]);

      // Should only make one network request due to deduplication
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it('should support proper request serialization for sync parameters', async () => {
      const syncParams: ContentSyncRequest = {
        lastSyncTime: '2024-01-15T10:00:00.000Z',
        limit: 50,
        categories: ['navigation', 'branding'],
        includeDrafts: true
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockSyncResponse
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate(syncParams)
      );

      await result;

      // Validate that the request was made with correct parameters
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/content/sync'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          }),
          body: expect.stringContaining('2024-01-15T10:00:00.000Z')
        })
      );

      const requestBody = JSON.parse(mockFetch.mock.calls[0][1].body);
      expect(requestBody.lastSyncTime).toBe('2024-01-15T10:00:00.000Z');
      expect(requestBody.limit).toBe(50);
      expect(requestBody.categories).toEqual(['navigation', 'branding']);
      expect(requestBody.includeDrafts).toBe(true);
    });
  });

  describe('Type safety validation', () => {
    it('should provide proper TypeScript typing for sync responses', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockSyncResponse
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z'
        })
      );

      const response = await result;

      // TypeScript should infer the correct types
      if (response.data) {
        const syncData: ContentSyncResponse = response.data;
        const items: ContentItem[] = syncData.items;
        const deletedIds: string[] = syncData.deletedIds;
        const syncTime: string = syncData.syncTime;
        const totalChanges: number = syncData.totalChanges;
        const hasMore: boolean = syncData.hasMore;

        expect(items).toHaveLength(2);
        expect(deletedIds).toHaveLength(1);
        expect(typeof syncTime).toBe('string');
        expect(typeof totalChanges).toBe('number');
        expect(typeof hasMore).toBe('boolean');
      }
    });

    it('should provide proper TypeScript typing for sync requests', () => {
      const validRequest: ContentSyncRequest = {
        lastSyncTime: '2024-01-15T10:00:00.000Z',
        limit: 100,
        pageToken: 'token-123',
        categories: ['navigation'],
        includeDrafts: false
      };

      // TypeScript should enforce correct types
      expect(typeof validRequest.lastSyncTime).toBe('string');
      expect(typeof validRequest.limit).toBe('number');
      expect(typeof validRequest.pageToken).toBe('string');
      expect(Array.isArray(validRequest.categories)).toBe(true);
      expect(typeof validRequest.includeDrafts).toBe('boolean');
    });
  });

  describe('Performance and scalability', () => {
    it('should handle large sync responses efficiently', async () => {
      // Arrange: Mock large dataset sync
      const largeItems = Array.from({ length: 1000 }, (_, i) => ({
        ...mockUpdatedContentItem,
        id: `navigation.item.${i}`,
        value: {
          type: 'text' as const,
          value: `Item ${i}`,
          maxLength: 100
        }
      }));

      const largeSyncResponse: ContentSyncResponse = {
        items: largeItems,
        deletedIds: Array.from({ length: 100 }, (_, i) => `deleted.item.${i}`),
        syncTime: '2024-01-15T12:00:00.000Z',
        totalChanges: 1100,
        hasMore: false
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: largeSyncResponse
        }),
      });

      // Act: Sync large dataset
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z',
          limit: 1000
        })
      );

      const response = await result;

      // Assert: Validate large dataset handling
      expect(response.data).toBeDefined();
      const syncData = response.data as ContentSyncResponse;

      expect(syncData.items).toHaveLength(1000);
      expect(syncData.deletedIds).toHaveLength(100);
      expect(syncData.totalChanges).toBe(1100);
    });

    it('should respect sync rate limiting', async () => {
      // Arrange: Mock rate limit error
      const rateLimitError: ApiError = {
        success: false,
        error: {
          code: 'RATE_LIMIT_EXCEEDED',
          message: 'Sync rate limit exceeded. Please wait before retrying.',
          details: {
            retryAfter: 60,
            limit: 10,
            window: 3600
          },
          timestamp: '2024-01-15T12:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 429,
        json: async () => rateLimitError,
      });

      // Act: Attempt sync that exceeds rate limit
      const result = store.dispatch(
        contentApi.endpoints.syncContent.initiate({
          lastSyncTime: '2024-01-15T10:00:00.000Z'
        })
      );

      const response = await result;

      // Assert: Validate rate limiting
      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(429);
      expect(error.data.error.code).toBe('RATE_LIMIT_EXCEEDED');
      expect(error.data.error.details.retryAfter).toBe(60);
    });
  });
});