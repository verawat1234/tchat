/**
 * Contract Test: archiveContent API Endpoint
 *
 * This test validates the exact contract structure for the content management API
 * archiveContent mutation endpoint. It follows TDD principles - the test WILL FAIL until
 * the backend API endpoints are implemented.
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * Validates:
 * - Content archival workflow and state transitions
 * - Status change to archived from published/draft
 * - Archive metadata (archivedAt, archivedBy, archiveReason)
 * - Archive reason requirement and validation
 * - Reversible archive operations
 * - Permission-based archive control
 * - Error handling for already archived content
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import { http, HttpResponse } from 'msw';
import { server } from '../../lib/test-utils/msw/server';
import { api } from '../api';
import { ContentStatus, ContentType } from '../../types/content';

// Contract interfaces for archiveContent functionality
export interface ArchiveContentRequest {
  /** Content item ID to archive */
  contentId: string;
  /** Reason for archiving (required) */
  archiveReason: string;
  /** Optional notes about the archival */
  notes?: string;
}

export interface ArchiveMetadata {
  /** ISO timestamp when content was archived */
  archivedAt: string;
  /** User ID who performed the archival */
  archivedBy: string;
  /** Reason provided for archiving */
  archiveReason: string;
  /** Optional notes about the archival */
  notes?: string;
}

export interface ContentItemWithArchive {
  id: string;
  category: string;
  type: ContentType;
  value: any;
  metadata: {
    createdAt: string;
    createdBy: string;
    updatedAt: string;
    updatedBy: string;
    version: number;
    tags?: string[];
    notes?: string;
  };
  status: ContentStatus;
  /** Archive metadata (only present if status is ARCHIVED) */
  archiveMetadata?: ArchiveMetadata;
}

export interface ArchiveContentResponse {
  /** The archived content item with archive metadata */
  item: ContentItemWithArchive;
  /** Success message */
  message: string;
}

export interface RestoreContentRequest {
  /** Content item ID to restore from archive */
  contentId: string;
  /** Target status to restore to (DRAFT or PUBLISHED) */
  targetStatus: ContentStatus.DRAFT | ContentStatus.PUBLISHED;
  /** Optional notes about the restoration */
  notes?: string;
}

export interface RestoreContentResponse {
  /** The restored content item */
  item: ContentItemWithArchive;
  /** Success message */
  message: string;
}

// Mock contentApi since archiveContent endpoint doesn't exist yet - this will fail in TDD fashion
const contentApi = {
  endpoints: {
    archiveContent: {
      initiate: (params: ArchiveContentRequest) => {
        throw new Error('archiveContent endpoint not implemented yet - this is expected for TDD');
      }
    },
    restoreContent: {
      initiate: (params: RestoreContentRequest) => {
        throw new Error('restoreContent endpoint not implemented yet - this is expected for TDD');
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
const createMockContentItem = (overrides: Partial<ContentItemWithArchive> = {}): ContentItemWithArchive => ({
  id: 'test-content-1',
  category: 'ui.buttons',
  type: ContentType.TEXT,
  value: {
    type: ContentType.TEXT,
    value: 'Test Button Label',
    locale: 'en',
  },
  metadata: {
    createdAt: '2023-01-01T00:00:00.000Z',
    createdBy: 'user-1',
    updatedAt: '2023-01-01T12:00:00.000Z',
    updatedBy: 'user-1',
    version: 1,
    tags: ['ui', 'button'],
    notes: 'Test content item',
  },
  status: ContentStatus.PUBLISHED,
  ...overrides,
});

const createMockArchiveMetadata = (overrides: Partial<ArchiveMetadata> = {}): ArchiveMetadata => ({
  archivedAt: '2023-01-02T10:00:00.000Z',
  archivedBy: 'user-1',
  archiveReason: 'Content no longer needed',
  notes: 'Archived during content cleanup',
  ...overrides,
});

// API Base URL for MSW
const API_BASE_URL = process.env.VITE_API_URL || 'http://localhost:3001';

describe('Content API - archiveContent Contract Test', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
  });

  afterEach(() => {
    server.resetHandlers();
  });

  describe('Archive Content Endpoint Structure', () => {
    it('should call correct endpoint with proper HTTP method and payload', async () => {
      // Arrange: Mock successful archive response
      const mockItem = createMockContentItem({
        id: 'content-to-archive',
        status: ContentStatus.ARCHIVED,
        archiveMetadata: createMockArchiveMetadata({
          archiveReason: 'Content redesign',
          notes: 'Replacing with new design system components',
        }),
      });

      const mockResponse: ArchiveContentResponse = {
        item: mockItem,
        message: 'Content archived successfully',
      };

      let capturedRequestBody: any = null;
      let capturedMethod: string = '';

      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, async ({ request }) => {
          capturedMethod = request.method;
          capturedRequestBody = await request.json();
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute archive mutation
      const archiveRequest: ArchiveContentRequest = {
        contentId: 'content-to-archive',
        archiveReason: 'Content redesign',
        notes: 'Replacing with new design system components',
      };

      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate(archiveRequest)
      );

      // Assert: Verify HTTP method and payload structure
      expect(capturedMethod).toBe('PATCH');
      expect(capturedRequestBody).toEqual(archiveRequest);
      expect(result.data).toEqual(mockResponse);
      expect(result.data?.item.status).toBe(ContentStatus.ARCHIVED);
      expect(result.data?.item.archiveMetadata).toBeDefined();
    });

    it('should return archived content with complete archive metadata', async () => {
      // Arrange: Mock complete archive response
      const archiveMetadata = createMockArchiveMetadata({
        archivedAt: '2023-01-02T15:30:00.000Z',
        archivedBy: 'admin-user',
        archiveReason: 'Policy violation',
        notes: 'Content violates community guidelines',
      });

      const mockItem = createMockContentItem({
        id: 'policy-violation-content',
        status: ContentStatus.ARCHIVED,
        archiveMetadata,
      });

      const mockResponse: ArchiveContentResponse = {
        item: mockItem,
        message: 'Content archived due to policy violation',
      };

      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute archive mutation
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'policy-violation-content',
          archiveReason: 'Policy violation',
          notes: 'Content violates community guidelines',
        })
      );

      // Assert: Verify complete archive metadata structure
      const archivedItem = result.data?.item;
      expect(archivedItem?.status).toBe(ContentStatus.ARCHIVED);
      expect(archivedItem?.archiveMetadata).toBeDefined();
      expect(archivedItem?.archiveMetadata?.archivedAt).toMatch(/\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/);
      expect(archivedItem?.archiveMetadata?.archivedBy).toBe('admin-user');
      expect(archivedItem?.archiveMetadata?.archiveReason).toBe('Policy violation');
      expect(archivedItem?.archiveMetadata?.notes).toBe('Content violates community guidelines');
    });
  });

  describe('Archive Reason Validation', () => {
    it('should require archiveReason in request', async () => {
      // Arrange: Mock validation error for missing archive reason
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json(
            {
              error: 'Validation Error',
              message: 'Archive reason is required',
              details: { archiveReason: 'This field is required' }
            },
            { status: 400 }
          );
        })
      );

      // Act: Attempt to archive without reason
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'test-content',
          archiveReason: '', // Empty reason
        })
      );

      // Assert: Verify validation error
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(400);
      expect(result.data).toBeUndefined();
    });

    it('should validate minimum length for archive reason', async () => {
      // Arrange: Mock validation error for short archive reason
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json(
            {
              error: 'Validation Error',
              message: 'Archive reason must be at least 10 characters',
              details: { archiveReason: 'Must be at least 10 characters long' }
            },
            { status: 400 }
          );
        })
      );

      // Act: Attempt to archive with short reason
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'test-content',
          archiveReason: 'Short', // Too short
        })
      );

      // Assert: Verify validation error
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(400);
      expect(result.data).toBeUndefined();
    });

    it('should accept valid archive reasons', async () => {
      // Arrange: Mock successful archive with valid reason
      const mockItem = createMockContentItem({
        status: ContentStatus.ARCHIVED,
        archiveMetadata: createMockArchiveMetadata({
          archiveReason: 'Content is outdated and no longer relevant to users',
        }),
      });

      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json({
            item: mockItem,
            message: 'Content archived successfully',
          });
        })
      );

      // Act: Archive with valid reason
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'test-content',
          archiveReason: 'Content is outdated and no longer relevant to users',
        })
      );

      // Assert: Verify successful archive
      expect(result.data).toBeDefined();
      expect(result.data?.item.archiveMetadata?.archiveReason).toBe(
        'Content is outdated and no longer relevant to users'
      );
    });
  });

  describe('Status Transition Validation', () => {
    it('should archive content from PUBLISHED status', async () => {
      // Arrange: Mock transitioning from published to archived
      const originalItem = createMockContentItem({
        id: 'published-content',
        status: ContentStatus.PUBLISHED,
      });

      const archivedItem = createMockContentItem({
        id: 'published-content',
        status: ContentStatus.ARCHIVED,
        archiveMetadata: createMockArchiveMetadata(),
      });

      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json({
            item: archivedItem,
            message: 'Published content archived successfully',
          });
        })
      );

      // Act: Archive published content
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'published-content',
          archiveReason: 'End of campaign period',
        })
      );

      // Assert: Verify status transition
      expect(result.data?.item.status).toBe(ContentStatus.ARCHIVED);
      expect(result.data?.item.archiveMetadata).toBeDefined();
    });

    it('should archive content from DRAFT status', async () => {
      // Arrange: Mock transitioning from draft to archived
      const archivedItem = createMockContentItem({
        id: 'draft-content',
        status: ContentStatus.ARCHIVED,
        archiveMetadata: createMockArchiveMetadata({
          archiveReason: 'Project cancelled before publication',
        }),
      });

      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json({
            item: archivedItem,
            message: 'Draft content archived successfully',
          });
        })
      );

      // Act: Archive draft content
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'draft-content',
          archiveReason: 'Project cancelled before publication',
        })
      );

      // Assert: Verify status transition
      expect(result.data?.item.status).toBe(ContentStatus.ARCHIVED);
      expect(result.data?.item.archiveMetadata?.archiveReason).toBe(
        'Project cancelled before publication'
      );
    });

    it('should handle already archived content appropriately', async () => {
      // Arrange: Mock error for already archived content
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json(
            {
              error: 'Invalid State Transition',
              message: 'Content is already archived',
              details: { status: 'Content cannot be archived from ARCHIVED status' }
            },
            { status: 400 }
          );
        })
      );

      // Act: Attempt to archive already archived content
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'already-archived-content',
          archiveReason: 'Duplicate archive attempt',
        })
      );

      // Assert: Verify appropriate error handling
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(400);
      expect(result.data).toBeUndefined();
    });
  });

  describe('Permission-Based Archive Control', () => {
    it('should handle unauthorized archive attempts', async () => {
      // Arrange: Mock unauthorized error
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json(
            {
              error: 'Unauthorized',
              message: 'Invalid or expired authentication token'
            },
            { status: 401 }
          );
        })
      );

      // Act: Attempt archive without proper authentication
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'protected-content',
          archiveReason: 'Unauthorized archive attempt',
        })
      );

      // Assert: Verify unauthorized response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(401);
      expect(result.data).toBeUndefined();
    });

    it('should handle insufficient permissions for archive', async () => {
      // Arrange: Mock forbidden error
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json(
            {
              error: 'Forbidden',
              message: 'Insufficient permissions to archive content in this category',
              details: { requiredRole: 'content-admin', userRole: 'content-editor' }
            },
            { status: 403 }
          );
        })
      );

      // Act: Attempt archive with insufficient permissions
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'admin-only-content',
          archiveReason: 'Insufficient permissions test',
        })
      );

      // Assert: Verify forbidden response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(403);
      expect(result.data).toBeUndefined();
    });

    it('should allow archive with proper permissions', async () => {
      // Arrange: Mock successful archive with proper permissions
      const mockItem = createMockContentItem({
        id: 'admin-content',
        status: ContentStatus.ARCHIVED,
        archiveMetadata: createMockArchiveMetadata({
          archivedBy: 'admin-user',
          archiveReason: 'Administrative content cleanup',
        }),
      });

      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json({
            item: mockItem,
            message: 'Content archived successfully by admin',
          });
        })
      );

      // Act: Archive with proper admin permissions
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'admin-content',
          archiveReason: 'Administrative content cleanup',
        })
      );

      // Assert: Verify successful archive
      expect(result.data).toBeDefined();
      expect(result.data?.item.archiveMetadata?.archivedBy).toBe('admin-user');
    });
  });

  describe('Reversible Archive Operations (Restore)', () => {
    it('should restore archived content to DRAFT status', async () => {
      // Arrange: Mock successful restore to draft
      const restoredItem = createMockContentItem({
        id: 'restored-content',
        status: ContentStatus.DRAFT,
        // archiveMetadata should be removed or nullified after restore
        archiveMetadata: undefined,
      });

      const mockResponse: RestoreContentResponse = {
        item: restoredItem,
        message: 'Content restored to draft status',
      };

      server.use(
        http.patch(`${API_BASE_URL}/api/content/restore`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Restore content to draft
      const result = await store.dispatch(
        contentApi.endpoints.restoreContent.initiate({
          contentId: 'restored-content',
          targetStatus: ContentStatus.DRAFT,
          notes: 'Restoring for further editing',
        })
      );

      // Assert: Verify successful restore
      expect(result.data?.item.status).toBe(ContentStatus.DRAFT);
      expect(result.data?.item.archiveMetadata).toBeUndefined();
      expect(result.data?.message).toBe('Content restored to draft status');
    });

    it('should restore archived content to PUBLISHED status', async () => {
      // Arrange: Mock successful restore to published
      const restoredItem = createMockContentItem({
        id: 'restored-published-content',
        status: ContentStatus.PUBLISHED,
        archiveMetadata: undefined,
      });

      const mockResponse: RestoreContentResponse = {
        item: restoredItem,
        message: 'Content restored to published status',
      };

      server.use(
        http.patch(`${API_BASE_URL}/api/content/restore`, () => {
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Restore content to published
      const result = await store.dispatch(
        contentApi.endpoints.restoreContent.initiate({
          contentId: 'restored-published-content',
          targetStatus: ContentStatus.PUBLISHED,
          notes: 'Content is relevant again',
        })
      );

      // Assert: Verify successful restore
      expect(result.data?.item.status).toBe(ContentStatus.PUBLISHED);
      expect(result.data?.item.archiveMetadata).toBeUndefined();
    });

    it('should validate target status for restore operations', async () => {
      // Arrange: Mock validation error for invalid target status
      server.use(
        http.patch(`${API_BASE_URL}/api/content/restore`, () => {
          return HttpResponse.json(
            {
              error: 'Validation Error',
              message: 'Invalid target status for restore operation',
              details: { targetStatus: 'Must be either DRAFT or PUBLISHED' }
            },
            { status: 400 }
          );
        })
      );

      // Act: Attempt restore with invalid target status
      const result = await store.dispatch(
        contentApi.endpoints.restoreContent.initiate({
          contentId: 'test-content',
          targetStatus: ContentStatus.ARCHIVED as any, // Invalid target
          notes: 'Invalid restore attempt',
        })
      );

      // Assert: Verify validation error
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(400);
    });

    it('should handle restore of non-archived content', async () => {
      // Arrange: Mock error for restoring non-archived content
      server.use(
        http.patch(`${API_BASE_URL}/api/content/restore`, () => {
          return HttpResponse.json(
            {
              error: 'Invalid State Transition',
              message: 'Only archived content can be restored',
              details: { currentStatus: 'PUBLISHED' }
            },
            { status: 400 }
          );
        })
      );

      // Act: Attempt to restore published content
      const result = await store.dispatch(
        contentApi.endpoints.restoreContent.initiate({
          contentId: 'published-content',
          targetStatus: ContentStatus.DRAFT,
        })
      );

      // Assert: Verify appropriate error
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(400);
    });
  });

  describe('Error Handling Scenarios', () => {
    it('should handle content not found errors', async () => {
      // Arrange: Mock not found error
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json(
            {
              error: 'Not Found',
              message: 'Content item not found',
              details: { contentId: 'non-existent-content' }
            },
            { status: 404 }
          );
        })
      );

      // Act: Attempt to archive non-existent content
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'non-existent-content',
          archiveReason: 'Testing not found scenario',
        })
      );

      // Assert: Verify not found response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(404);
      expect(result.data).toBeUndefined();
    });

    it('should handle server errors during archive operation', async () => {
      // Arrange: Mock internal server error
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json(
            {
              error: 'Internal Server Error',
              message: 'Database transaction failed during archive operation'
            },
            { status: 500 }
          );
        })
      );

      // Act: Attempt archive that triggers server error
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'problematic-content',
          archiveReason: 'Server error test case',
        })
      );

      // Assert: Verify server error response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe(500);
      expect(result.data).toBeUndefined();
    });

    it('should handle network errors during archive operation', async () => {
      // Arrange: Mock network error
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.error();
        })
      );

      // Act: Attempt archive with network failure
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'test-content',
          archiveReason: 'Network error test case',
        })
      );

      // Assert: Verify network error response
      expect(result.error).toBeDefined();
      expect(result.error?.status).toBe('FETCH_ERROR');
      expect(result.data).toBeUndefined();
    });
  });

  describe('RTK Query Integration and Caching', () => {
    it('should provide proper cache tags for archive mutations', async () => {
      // Arrange: Mock successful archive response
      const mockItem = createMockContentItem({
        status: ContentStatus.ARCHIVED,
        archiveMetadata: createMockArchiveMetadata(),
      });

      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          return HttpResponse.json({
            item: mockItem,
            message: 'Content archived successfully',
          });
        })
      );

      // Act: Execute archive mutation
      const result = await store.dispatch(
        contentApi.endpoints.archiveContent.initiate({
          contentId: 'test-content',
          archiveReason: 'Cache invalidation test',
        })
      );

      // Assert: Verify that the mutation provides proper cache invalidation
      // The actual cache tags are set in the endpoint definition and would be:
      // - Individual item tag: { type: 'ContentItem', id: contentId }
      // - List tag: { type: 'ContentItem', id: 'LIST' }
      expect(result.data).toBeDefined();
      expect(result.data?.item.status).toBe(ContentStatus.ARCHIVED);

      // Verify that the endpoint exists and can be called
      const endpoint = contentApi.endpoints.archiveContent;
      expect(endpoint).toBeDefined();
    });

    it('should handle mutation argument serialization correctly', async () => {
      // Arrange: Mock response for serialization test
      const mockResponse = {
        item: createMockContentItem({ status: ContentStatus.ARCHIVED }),
        message: 'Archived successfully',
      };

      let requestCount = 0;
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          requestCount++;
          return HttpResponse.json(mockResponse);
        })
      );

      // Act: Execute same mutation twice (should not be deduplicated like queries)
      const mutationArgs = {
        contentId: 'serialization-test',
        archiveReason: 'Testing mutation serialization',
      };

      await store.dispatch(
        contentApi.endpoints.archiveContent.initiate(mutationArgs)
      );

      await store.dispatch(
        contentApi.endpoints.archiveContent.initiate(mutationArgs)
      );

      // Assert: Verify that mutations are not deduplicated (both requests made)
      expect(requestCount).toBe(2);
    });
  });

  describe('Archive Workflow Integration', () => {
    it('should support bulk archive operations tracking', async () => {
      // This test validates that the archive endpoint can be used as part of
      // bulk operations even though it's a single-item endpoint

      // Arrange: Mock multiple successful archive operations
      const contentIds = ['bulk-1', 'bulk-2', 'bulk-3'];
      const mockResponses = contentIds.map(id => ({
        item: createMockContentItem({
          id,
          status: ContentStatus.ARCHIVED,
          archiveMetadata: createMockArchiveMetadata({
            archiveReason: 'Bulk archive operation',
          }),
        }),
        message: `Content ${id} archived successfully`,
      }));

      let responseIndex = 0;
      server.use(
        http.patch(`${API_BASE_URL}/api/content/archive`, () => {
          const response = mockResponses[responseIndex];
          responseIndex++;
          return HttpResponse.json(response);
        })
      );

      // Act: Execute multiple archive operations
      const results = await Promise.all(
        contentIds.map(id =>
          store.dispatch(
            contentApi.endpoints.archiveContent.initiate({
              contentId: id,
              archiveReason: 'Bulk archive operation',
            })
          )
        )
      );

      // Assert: Verify all operations completed successfully
      results.forEach((result, index) => {
        expect(result.data?.item.id).toBe(contentIds[index]);
        expect(result.data?.item.status).toBe(ContentStatus.ARCHIVED);
        expect(result.data?.item.archiveMetadata?.archiveReason).toBe('Bulk archive operation');
      });
    });
  });
});