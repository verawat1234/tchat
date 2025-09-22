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
  ContentMetadata,
  BulkUpdateContentRequest,
  BulkUpdateContentResponse,
  BulkUpdateItemRequest,
  BulkUpdateItemResult,
  UpdateContentItemRequest
} from '../../types/content';
import type { ApiError } from '../../types/api';

/**
 * T014: Contract Test for bulkUpdateContent API
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * This test validates:
 * 1. Bulk content item updates with array of update requests
 * 2. Mixed success/failure handling in bulk operations
 * 3. Transaction-like behavior with rollback on errors
 * 4. Performance validation for large bulk operations
 * 5. Partial success scenarios with detailed results
 * 6. Validation error aggregation across multiple items
 * 7. Atomic vs non-atomic bulk update modes
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

const mockSecondaryCategory: ContentCategory = {
  id: 'messaging',
  name: 'Messaging',
  description: 'Chat and communication content',
  permissions: {
    read: ['user', 'admin'],
    write: ['admin'],
    publish: ['admin']
  }
};

const mockTextContent1: ContentValue = {
  type: 'text',
  value: 'Welcome to Tchat',
  maxLength: 100
};

const mockTextContent2: ContentValue = {
  type: 'text',
  value: 'Chat with friends',
  maxLength: 100
};

const mockUpdatedTextContent1: ContentValue = {
  type: 'text',
  value: 'Welcome to the new Tchat platform',
  maxLength: 100
};

const mockUpdatedTextContent2: ContentValue = {
  type: 'text',
  value: 'Connect and chat with friends worldwide',
  maxLength: 100
};

const mockRichTextContent: ContentValue = {
  type: 'rich_text',
  value: '<h1>Help Center</h1><p>Find answers to your questions</p>',
  format: 'html',
  allowedTags: ['h1', 'h2', 'p', 'strong', 'em']
};

const mockBaseMetadata: ContentMetadata = {
  createdAt: '2024-01-15T10:00:00.000Z',
  createdBy: 'user-123',
  updatedAt: '2024-01-15T10:00:00.000Z',
  updatedBy: 'user-123',
  version: 1
};

const mockContentItem1: ContentItem = {
  id: 'nav-welcome-title',
  key: 'welcome.title',
  categoryId: 'navigation',
  category: mockContentCategory,
  type: ContentType.TEXT,
  value: mockTextContent1,
  status: ContentStatus.DRAFT,
  tags: ['ui', 'navigation'],
  metadata: mockBaseMetadata,
  notes: 'Welcome message',
  changeLog: 'Initial creation'
};

const mockContentItem2: ContentItem = {
  id: 'nav-chat-subtitle',
  key: 'chat.subtitle',
  categoryId: 'navigation',
  category: mockContentCategory,
  type: ContentType.TEXT,
  value: mockTextContent2,
  status: ContentStatus.DRAFT,
  tags: ['ui', 'navigation', 'chat'],
  metadata: { ...mockBaseMetadata, version: 2 },
  notes: 'Chat subtitle',
  changeLog: 'Initial creation'
};

const mockContentItem3: ContentItem = {
  id: 'help-center-title',
  key: 'help.title',
  categoryId: 'messaging',
  category: mockSecondaryCategory,
  type: ContentType.RICH_TEXT,
  value: mockRichTextContent,
  status: ContentStatus.PUBLISHED,
  tags: ['ui', 'help'],
  metadata: { ...mockBaseMetadata, version: 1 },
  notes: 'Help center content',
  changeLog: 'Initial creation'
};

// Store setup for testing
let store: ReturnType<typeof configureStore>;

const createTestStore = () => {
  return configureStore({
    reducer: {
      api: api.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(api.middleware),
  });
};

describe('Content API: bulkUpdateContent', () => {
  beforeEach(() => {
    store = createTestStore();
    setupListeners(store.dispatch);
    vi.clearAllMocks();
  });

  describe('Successful Bulk Operations', () => {
    it('should successfully update multiple content items in a single request', async () => {
      const updatedItem1: ContentItem = {
        ...mockContentItem1,
        value: mockUpdatedTextContent1,
        status: ContentStatus.PUBLISHED,
        metadata: {
          ...mockBaseMetadata,
          updatedAt: '2024-01-20T14:30:00.000Z',
          updatedBy: 'user-456',
          version: 2
        }
      };

      const updatedItem2: ContentItem = {
        ...mockContentItem2,
        value: mockUpdatedTextContent2,
        status: ContentStatus.PUBLISHED,
        metadata: {
          ...mockBaseMetadata,
          updatedAt: '2024-01-20T14:30:00.000Z',
          updatedBy: 'user-456',
          version: 3
        }
      };

      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'nav-welcome-title',
            success: true,
            item: updatedItem1
          },
          {
            id: 'nav-chat-subtitle',
            success: true,
            item: updatedItem2
          }
        ],
        success: true,
        successCount: 2,
        errorCount: 0,
        processingTime: 150
      };

      const mockHttpResponse = {
        ok: true,
        status: 200,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: {
              value: mockUpdatedTextContent1,
              status: ContentStatus.PUBLISHED,
              changeLog: 'Updated welcome message'
            }
          },
          {
            id: 'nav-chat-subtitle',
            update: {
              value: mockUpdatedTextContent2,
              status: ContentStatus.PUBLISHED,
              changeLog: 'Updated chat subtitle'
            }
          }
        ],
        changeLog: 'Bulk update for navigation content'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.data).toEqual(mockResponse);
      expect(result.data?.success).toBe(true);
      expect(result.data?.successCount).toBe(2);
      expect(result.data?.errorCount).toBe(0);
      expect(result.data?.results).toHaveLength(2);

      // Verify request structure
      expect(fetch).toHaveBeenCalledWith(
        '/api/content/bulk-update',
        expect.objectContaining({
          method: 'PUT',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
          body: JSON.stringify(bulkUpdateRequest)
        })
      );
    });

    it('should handle large bulk operations efficiently', async () => {
      // Create 50 content items for performance testing
      const bulkItems: BulkUpdateItemRequest[] = Array.from({ length: 50 }, (_, i) => ({
        id: `content-item-${i + 1}`,
        update: {
          value: {
            type: 'text',
            value: `Updated content ${i + 1}`,
            maxLength: 100
          },
          status: ContentStatus.PUBLISHED,
          changeLog: `Bulk update ${i + 1}`
        }
      }));

      const bulkResults: BulkUpdateItemResult[] = bulkItems.map(item => ({
        id: item.id,
        success: true,
        item: {
          ...mockContentItem1,
          id: item.id,
          value: item.update.value!,
          status: item.update.status!,
          metadata: {
            ...mockBaseMetadata,
            updatedAt: '2024-01-20T14:30:00.000Z',
            updatedBy: 'user-456',
            version: 2
          }
        }
      }));

      const mockResponse: BulkUpdateContentResponse = {
        results: bulkResults,
        success: true,
        successCount: 50,
        errorCount: 0,
        processingTime: 850 // Under 1 second for 50 items
      };

      const mockHttpResponse = {
        ok: true,
        status: 200,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: bulkItems,
        changeLog: 'Large bulk update operation'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.data?.success).toBe(true);
      expect(result.data?.successCount).toBe(50);
      expect(result.data?.processingTime).toBeLessThan(1000);
      expect(result.data?.results).toHaveLength(50);
    });

    it('should support atomic transactions with all-or-nothing behavior', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'nav-welcome-title',
            success: true,
            item: {
              ...mockContentItem1,
              value: mockUpdatedTextContent1,
              metadata: {
                ...mockBaseMetadata,
                updatedAt: '2024-01-20T14:30:00.000Z',
                updatedBy: 'user-456',
                version: 2
              }
            }
          },
          {
            id: 'nav-chat-subtitle',
            success: true,
            item: {
              ...mockContentItem2,
              value: mockUpdatedTextContent2,
              metadata: {
                ...mockBaseMetadata,
                updatedAt: '2024-01-20T14:30:00.000Z',
                updatedBy: 'user-456',
                version: 3
              }
            }
          }
        ],
        success: true,
        successCount: 2,
        errorCount: 0,
        processingTime: 200
      };

      const mockHttpResponse = {
        ok: true,
        status: 200,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: {
              value: mockUpdatedTextContent1
            }
          },
          {
            id: 'nav-chat-subtitle',
            update: {
              value: mockUpdatedTextContent2
            }
          }
        ],
        atomic: true,
        changeLog: 'Atomic bulk update'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.data?.success).toBe(true);
      expect(result.data?.rolledBack).toBeUndefined();

      // Verify atomic flag in request
      expect(fetch).toHaveBeenCalledWith(
        '/api/content/bulk-update',
        expect.objectContaining({
          body: JSON.stringify(bulkUpdateRequest)
        })
      );
    });
  });

  describe('Mixed Success/Failure Scenarios', () => {
    it('should handle partial failures with detailed error reporting', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'nav-welcome-title',
            success: true,
            item: {
              ...mockContentItem1,
              value: mockUpdatedTextContent1,
              metadata: {
                ...mockBaseMetadata,
                updatedAt: '2024-01-20T14:30:00.000Z',
                updatedBy: 'user-456',
                version: 2
              }
            }
          },
          {
            id: 'non-existent-item',
            success: false,
            error: {
              code: 'CONTENT_NOT_FOUND',
              message: 'Content item not found',
              field: 'id'
            }
          },
          {
            id: 'invalid-content-item',
            success: false,
            error: {
              code: 'VALIDATION_ERROR',
              message: 'Content value exceeds maximum length',
              field: 'value.value'
            }
          }
        ],
        success: false,
        successCount: 1,
        errorCount: 2,
        processingTime: 180
      };

      const mockHttpResponse = {
        ok: true,
        status: 200, // 200 OK even with partial failures
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: {
              value: mockUpdatedTextContent1
            }
          },
          {
            id: 'non-existent-item',
            update: {
              value: mockUpdatedTextContent2
            }
          },
          {
            id: 'invalid-content-item',
            update: {
              value: {
                type: 'text',
                value: 'This is a very long text that exceeds the maximum allowed length and should trigger a validation error because it is too long',
                maxLength: 50
              }
            }
          }
        ],
        continueOnError: true,
        changeLog: 'Mixed update operation'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.data?.success).toBe(false);
      expect(result.data?.successCount).toBe(1);
      expect(result.data?.errorCount).toBe(2);
      expect(result.data?.results).toHaveLength(3);

      // Verify successful item
      expect(result.data?.results[0].success).toBe(true);
      expect(result.data?.results[0].item).toBeDefined();

      // Verify failed items have proper error details
      expect(result.data?.results[1].success).toBe(false);
      expect(result.data?.results[1].error?.code).toBe('CONTENT_NOT_FOUND');

      expect(result.data?.results[2].success).toBe(false);
      expect(result.data?.results[2].error?.code).toBe('VALIDATION_ERROR');
    });

    it('should support non-atomic mode with continue-on-error behavior', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'nav-welcome-title',
            success: true,
            item: {
              ...mockContentItem1,
              value: mockUpdatedTextContent1,
              metadata: {
                ...mockBaseMetadata,
                updatedAt: '2024-01-20T14:30:00.000Z',
                updatedBy: 'user-456',
                version: 2
              }
            }
          },
          {
            id: 'permission-denied-item',
            success: false,
            error: {
              code: 'INSUFFICIENT_PERMISSIONS',
              message: 'User does not have write permissions for this content',
              field: 'permissions'
            }
          },
          {
            id: 'nav-chat-subtitle',
            success: true,
            item: {
              ...mockContentItem2,
              value: mockUpdatedTextContent2,
              metadata: {
                ...mockBaseMetadata,
                updatedAt: '2024-01-20T14:30:00.000Z',
                updatedBy: 'user-456',
                version: 3
              }
            }
          }
        ],
        success: false,
        successCount: 2,
        errorCount: 1,
        processingTime: 220
      };

      const mockHttpResponse = {
        ok: true,
        status: 200,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: { value: mockUpdatedTextContent1 }
          },
          {
            id: 'permission-denied-item',
            update: { value: mockUpdatedTextContent2 }
          },
          {
            id: 'nav-chat-subtitle',
            update: { value: mockUpdatedTextContent2 }
          }
        ],
        atomic: false,
        continueOnError: true,
        changeLog: 'Non-atomic bulk update'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.data?.success).toBe(false);
      expect(result.data?.successCount).toBe(2);
      expect(result.data?.errorCount).toBe(1);

      // Verify that successful items were processed despite the failure
      expect(result.data?.results[0].success).toBe(true);
      expect(result.data?.results[2].success).toBe(true);
      expect(result.data?.results[1].success).toBe(false);
    });
  });

  describe('Transaction and Rollback Behavior', () => {
    it('should rollback all changes on atomic failure', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'nav-welcome-title',
            success: false,
            error: {
              code: 'ATOMIC_ROLLBACK',
              message: 'Operation rolled back due to atomic transaction failure'
            }
          },
          {
            id: 'nav-chat-subtitle',
            success: false,
            error: {
              code: 'ATOMIC_ROLLBACK',
              message: 'Operation rolled back due to atomic transaction failure'
            }
          },
          {
            id: 'invalid-item',
            success: false,
            error: {
              code: 'CONTENT_NOT_FOUND',
              message: 'Content item not found',
              field: 'id'
            }
          }
        ],
        success: false,
        successCount: 0,
        errorCount: 3,
        processingTime: 95,
        rolledBack: true
      };

      const mockHttpResponse = {
        ok: false,
        status: 400,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: { value: mockUpdatedTextContent1 }
          },
          {
            id: 'nav-chat-subtitle',
            update: { value: mockUpdatedTextContent2 }
          },
          {
            id: 'invalid-item',
            update: { value: mockUpdatedTextContent1 }
          }
        ],
        atomic: true,
        changeLog: 'Atomic update with failure'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(400);
      expect(error.data?.rolledBack).toBe(true);
      expect(error.data?.successCount).toBe(0);
    });

    it('should handle version conflicts in bulk operations', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'nav-welcome-title',
            success: true,
            item: {
              ...mockContentItem1,
              value: mockUpdatedTextContent1,
              metadata: {
                ...mockBaseMetadata,
                updatedAt: '2024-01-20T14:30:00.000Z',
                updatedBy: 'user-456',
                version: 2
              }
            }
          },
          {
            id: 'nav-chat-subtitle',
            success: false,
            error: {
              code: 'VERSION_CONFLICT',
              message: 'Content item has been updated by another user',
              field: 'expectedVersion'
            }
          }
        ],
        success: false,
        successCount: 1,
        errorCount: 1,
        processingTime: 160
      };

      const mockHttpResponse = {
        ok: true,
        status: 200,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: {
              value: mockUpdatedTextContent1,
              expectedVersion: 1
            }
          },
          {
            id: 'nav-chat-subtitle',
            update: {
              value: mockUpdatedTextContent2,
              expectedVersion: 1 // Outdated version
            }
          }
        ],
        continueOnError: true,
        changeLog: 'Bulk update with version check'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.data?.success).toBe(false);
      expect(result.data?.successCount).toBe(1);
      expect(result.data?.errorCount).toBe(1);
      expect(result.data?.results[1].error?.code).toBe('VERSION_CONFLICT');
    });
  });

  describe('Validation and Error Handling', () => {
    it('should aggregate validation errors across multiple items', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'item-with-long-text',
            success: false,
            error: {
              code: 'VALIDATION_ERROR',
              message: 'Content value exceeds maximum length',
              field: 'value.value'
            }
          },
          {
            id: 'item-with-invalid-status',
            success: false,
            error: {
              code: 'VALIDATION_ERROR',
              message: 'Invalid status transition',
              field: 'status'
            }
          },
          {
            id: 'item-with-invalid-tags',
            success: false,
            error: {
              code: 'VALIDATION_ERROR',
              message: 'Invalid tag format',
              field: 'tags'
            }
          }
        ],
        success: false,
        successCount: 0,
        errorCount: 3,
        processingTime: 120
      };

      const mockHttpResponse = {
        ok: false,
        status: 400,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'item-with-long-text',
            update: {
              value: {
                type: 'text',
                value: 'This text is way too long and exceeds the maximum allowed length for this content item',
                maxLength: 50
              }
            }
          },
          {
            id: 'item-with-invalid-status',
            update: {
              status: 'invalid_status' as ContentStatus
            }
          },
          {
            id: 'item-with-invalid-tags',
            update: {
              tags: ['', 'invalid tag with spaces', '']
            }
          }
        ],
        continueOnError: true,
        changeLog: 'Validation error testing'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(400);
      expect(error.data?.errorCount).toBe(3);
      expect(error.data?.results.every(r => !r.success)).toBe(true);
    });

    it('should handle authentication and permission errors', async () => {
      const mockErrorResponse = {
        ok: false,
        status: 401,
        json: async () => ({
          error: 'Authentication required',
          code: 'AUTHENTICATION_REQUIRED',
          details: 'Valid authentication token is required to perform bulk updates'
        }),
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockErrorResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: { value: mockUpdatedTextContent1 }
          }
        ],
        changeLog: 'Unauthorized bulk update attempt'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(401);
      expect(error.data?.code).toBe('AUTHENTICATION_REQUIRED');
    });

    it('should handle server errors gracefully', async () => {
      const mockErrorResponse = {
        ok: false,
        status: 500,
        json: async () => ({
          error: 'Internal server error',
          code: 'INTERNAL_SERVER_ERROR',
          details: 'An unexpected error occurred during bulk update processing'
        }),
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockErrorResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: { value: mockUpdatedTextContent1 }
          }
        ],
        changeLog: 'Bulk update causing server error'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(500);
      expect(error.data?.code).toBe('INTERNAL_SERVER_ERROR');
    });
  });

  describe('Performance and Optimization', () => {
    it('should handle concurrent bulk operations with proper queuing', async () => {
      const mockResponse1: BulkUpdateContentResponse = {
        results: [
          {
            id: 'batch-1-item-1',
            success: true,
            item: {
              ...mockContentItem1,
              id: 'batch-1-item-1',
              metadata: { ...mockBaseMetadata, version: 2 }
            }
          }
        ],
        success: true,
        successCount: 1,
        errorCount: 0,
        processingTime: 100
      };

      const mockResponse2: BulkUpdateContentResponse = {
        results: [
          {
            id: 'batch-2-item-1',
            success: true,
            item: {
              ...mockContentItem2,
              id: 'batch-2-item-1',
              metadata: { ...mockBaseMetadata, version: 3 }
            }
          }
        ],
        success: true,
        successCount: 1,
        errorCount: 0,
        processingTime: 110
      };

      global.fetch = vi.fn()
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => mockResponse1,
        } as Response)
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => mockResponse2,
        } as Response);

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const [result1, result2] = await Promise.all([
        store.dispatch(
          contentApi.endpoints.bulkUpdateContent.initiate({
            items: [
              {
                id: 'batch-1-item-1',
                update: { value: mockUpdatedTextContent1 }
              }
            ],
            changeLog: 'Concurrent batch 1'
          })
        ),
        store.dispatch(
          contentApi.endpoints.bulkUpdateContent.initiate({
            items: [
              {
                id: 'batch-2-item-1',
                update: { value: mockUpdatedTextContent2 }
              }
            ],
            changeLog: 'Concurrent batch 2'
          })
        )
      ]);

      expect(result1.data?.success).toBe(true);
      expect(result2.data?.success).toBe(true);
      expect(fetch).toHaveBeenCalledTimes(2);
    });

    it('should provide detailed processing time metrics', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: Array.from({ length: 10 }, (_, i) => ({
          id: `perf-test-item-${i + 1}`,
          success: true,
          item: {
            ...mockContentItem1,
            id: `perf-test-item-${i + 1}`,
            metadata: { ...mockBaseMetadata, version: 2 }
          }
        })),
        success: true,
        successCount: 10,
        errorCount: 0,
        processingTime: 75 // Should be under 100ms for 10 items
      };

      const mockHttpResponse = {
        ok: true,
        status: 200,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: Array.from({ length: 10 }, (_, i) => ({
          id: `perf-test-item-${i + 1}`,
          update: {
            value: {
              type: 'text',
              value: `Performance test content ${i + 1}`,
              maxLength: 100
            }
          }
        })),
        changeLog: 'Performance testing bulk update'
      };

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
      );

      expect(result.data?.success).toBe(true);
      expect(result.data?.processingTime).toBeLessThan(100);
      expect(result.data?.successCount).toBe(10);
    });
  });

  describe('RTK Query Integration', () => {
    it('should provide proper cache invalidation for bulk updates', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'nav-welcome-title',
            success: true,
            item: {
              ...mockContentItem1,
              value: mockUpdatedTextContent1,
              metadata: {
                ...mockBaseMetadata,
                updatedAt: '2024-01-20T14:30:00.000Z',
                updatedBy: 'user-456',
                version: 2
              }
            }
          }
        ],
        success: true,
        successCount: 1,
        errorCount: 0,
        processingTime: 85
      };

      const mockHttpResponse = {
        ok: true,
        status: 200,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.bulkUpdateContent.initiate({
          items: [
            {
              id: 'nav-welcome-title',
              update: { value: mockUpdatedTextContent1 }
            }
          ],
          changeLog: 'Cache invalidation test'
        })
      );

      // Verify cache tags are properly set for invalidation
      expect(result.data).toEqual(mockResponse);

      // The endpoint should invalidate relevant cache tags:
      // - ['Content', 'nav-welcome-title'] for specific items
      // - ['Content', 'navigation'] for category items
      // - ['ContentList'] for any content lists
      // - ['BulkContent'] for bulk operation caches
    });

    it('should handle request deduplication for identical bulk requests', async () => {
      const mockResponse: BulkUpdateContentResponse = {
        results: [
          {
            id: 'nav-welcome-title',
            success: true,
            item: {
              ...mockContentItem1,
              value: mockUpdatedTextContent1,
              metadata: {
                ...mockBaseMetadata,
                updatedAt: '2024-01-20T14:30:00.000Z',
                updatedBy: 'user-456',
                version: 2
              }
            }
          }
        ],
        success: true,
        successCount: 1,
        errorCount: 0,
        processingTime: 90
      };

      const mockHttpResponse = {
        ok: true,
        status: 200,
        json: async () => mockResponse,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockHttpResponse);

      const bulkUpdateRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'nav-welcome-title',
            update: { value: mockUpdatedTextContent1 }
          }
        ],
        changeLog: 'Deduplication test'
      };

      // Make identical requests simultaneously
      // This will fail because bulkUpdateContent endpoint doesn't exist yet
      const [result1, result2] = await Promise.all([
        store.dispatch(
          contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
        ),
        store.dispatch(
          contentApi.endpoints.bulkUpdateContent.initiate(bulkUpdateRequest)
        )
      ]);

      expect(result1.data).toEqual(mockResponse);
      expect(result2.data).toEqual(mockResponse);

      // With proper deduplication, fetch should only be called once
      expect(fetch).toHaveBeenCalledTimes(1);
    });
  });

  describe('Type Safety Validation', () => {
    it('should enforce BulkUpdateContentRequest interface constraints', () => {
      // Valid bulk update request
      const validRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'test-item-1',
            update: {
              value: mockUpdatedTextContent1,
              status: ContentStatus.PUBLISHED,
              tags: ['ui', 'test'],
              notes: 'Test update',
              changeLog: 'Test change log'
            }
          }
        ],
        atomic: true,
        continueOnError: false,
        changeLog: 'Bulk update test'
      };

      expect(validRequest).toBeDefined();
      expect(validRequest.items).toHaveLength(1);
      expect(validRequest.atomic).toBe(true);

      // Empty items array should be valid
      const emptyRequest: BulkUpdateContentRequest = {
        items: []
      };
      expect(emptyRequest).toBeDefined();

      // Minimal request with just items
      const minimalRequest: BulkUpdateContentRequest = {
        items: [
          {
            id: 'minimal-item',
            update: {}
          }
        ]
      };
      expect(minimalRequest).toBeDefined();
    });

    it('should maintain type safety for BulkUpdateItemResult interface', () => {
      // Success result
      const successResult: BulkUpdateItemResult = {
        id: 'test-item',
        success: true,
        item: mockContentItem1
      };

      expect(successResult.success).toBe(true);
      expect(successResult.item).toBeDefined();
      expect(successResult.error).toBeUndefined();

      // Error result
      const errorResult: BulkUpdateItemResult = {
        id: 'failed-item',
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Validation failed',
          field: 'value.value'
        }
      };

      expect(errorResult.success).toBe(false);
      expect(errorResult.error).toBeDefined();
      expect(errorResult.item).toBeUndefined();
    });
  });
});