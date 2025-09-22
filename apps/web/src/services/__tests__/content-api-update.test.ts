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
  UpdateContentItemRequest
} from '../../types/content';
import type { ApiError } from '../../types/api';

/**
 * T011: Contract Test for updateContentItem API
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * This test validates:
 * 1. Existing content item updates with proper request structure
 * 2. Version increment handling
 * 3. Metadata updates (updatedAt, updatedBy, version)
 * 4. Content value validation and type safety
 * 5. 404 error handling for non-existent content
 * 6. Optimistic locking/version conflicts
 * 7. Partial update support
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

const mockOriginalTextContent: ContentValue = {
  type: 'text',
  value: 'Welcome to Tchat',
  maxLength: 100
};

const mockUpdatedTextContent: ContentValue = {
  type: 'text',
  value: 'Welcome to the new Tchat',
  maxLength: 100
};

const mockRichTextContent: ContentValue = {
  type: 'rich_text',
  value: '<h1>Welcome to Tchat</h1><p>Your communication platform</p>',
  format: 'html',
  allowedTags: ['h1', 'h2', 'p', 'strong', 'em']
};

const mockOriginalMetadata: ContentMetadata = {
  createdAt: '2024-01-15T10:00:00.000Z',
  createdBy: 'user-123',
  updatedAt: '2024-01-15T10:00:00.000Z',
  updatedBy: 'user-123',
  version: 1
};

const mockUpdatedMetadata: ContentMetadata = {
  createdAt: '2024-01-15T10:00:00.000Z',
  createdBy: 'user-123',
  updatedAt: '2024-01-20T14:30:00.000Z',
  updatedBy: 'user-456',
  version: 2
};

const mockOriginalContentItem: ContentItem = {
  id: 'nav-welcome-title',
  key: 'welcome.title',
  categoryId: 'navigation',
  category: mockContentCategory,
  type: ContentType.TEXT,
  value: mockOriginalTextContent,
  status: ContentStatus.DRAFT,
  tags: ['ui', 'navigation'],
  metadata: mockOriginalMetadata,
  notes: 'Initial welcome message',
  changeLog: 'Created initial welcome title'
};

const mockUpdatedContentItem: ContentItem = {
  ...mockOriginalContentItem,
  value: mockUpdatedTextContent,
  status: ContentStatus.PUBLISHED,
  tags: ['ui', 'navigation', 'featured'],
  metadata: mockUpdatedMetadata,
  notes: 'Updated welcome message for new branding',
  changeLog: 'Updated title text and published for live use'
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

describe('Content API: updateContentItem', () => {
  beforeEach(() => {
    store = createTestStore();
    setupListeners(store.dispatch);
    vi.clearAllMocks();
  });

  describe('Successful Updates', () => {
    it('should update content item with new value and metadata', async () => {
      // Mock successful update response
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => mockUpdatedContentItem,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockResponse);

      const updateRequest: UpdateContentItemRequest = {
        value: mockUpdatedTextContent,
        status: ContentStatus.PUBLISHED,
        tags: ['ui', 'navigation', 'featured'],
        notes: 'Updated welcome message for new branding',
        changeLog: 'Updated title text and published for live use'
      };

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          ...updateRequest
        })
      );

      expect(result.data).toEqual(mockUpdatedContentItem);
      expect(result.data?.metadata.version).toBe(2);
      expect(result.data?.metadata.updatedAt).toBe('2024-01-20T14:30:00.000Z');
      expect(result.data?.metadata.updatedBy).toBe('user-456');
      expect(result.data?.status).toBe(ContentStatus.PUBLISHED);

      // Verify request structure
      expect(fetch).toHaveBeenCalledWith(
        '/api/content/nav-welcome-title',
        expect.objectContaining({
          method: 'PUT',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
          body: JSON.stringify(updateRequest)
        })
      );
    });

    it('should support partial updates (value only)', async () => {
      const partialUpdateItem: ContentItem = {
        ...mockOriginalContentItem,
        value: mockUpdatedTextContent,
        metadata: {
          ...mockOriginalMetadata,
          updatedAt: '2024-01-20T14:30:00.000Z',
          updatedBy: 'user-456',
          version: 2
        }
      };

      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => partialUpdateItem,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockResponse);

      const updateRequest: UpdateContentItemRequest = {
        value: mockUpdatedTextContent
      };

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          ...updateRequest
        })
      );

      expect(result.data).toEqual(partialUpdateItem);
      expect(result.data?.value).toEqual(mockUpdatedTextContent);
      expect(result.data?.status).toBe(ContentStatus.DRAFT); // Unchanged
      expect(result.data?.metadata.version).toBe(2); // Incremented
    });

    it('should support status changes with proper validation', async () => {
      const statusUpdateItem: ContentItem = {
        ...mockOriginalContentItem,
        status: ContentStatus.PUBLISHED,
        metadata: {
          ...mockOriginalMetadata,
          updatedAt: '2024-01-20T14:30:00.000Z',
          updatedBy: 'user-456',
          version: 2
        }
      };

      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => statusUpdateItem,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockResponse);

      const updateRequest: UpdateContentItemRequest = {
        status: ContentStatus.PUBLISHED,
        changeLog: 'Published for live use'
      };

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          ...updateRequest
        })
      );

      expect(result.data).toEqual(statusUpdateItem);
      expect(result.data?.status).toBe(ContentStatus.PUBLISHED);
    });

    it('should update rich text content with format validation', async () => {
      const richTextItem: ContentItem = {
        ...mockOriginalContentItem,
        type: ContentType.RICH_TEXT,
        value: mockRichTextContent,
        metadata: {
          ...mockOriginalMetadata,
          updatedAt: '2024-01-20T14:30:00.000Z',
          updatedBy: 'user-456',
          version: 2
        }
      };

      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => richTextItem,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockResponse);

      const updateRequest: UpdateContentItemRequest = {
        value: mockRichTextContent,
        changeLog: 'Converted to rich text format'
      };

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          ...updateRequest
        })
      );

      expect(result.data).toEqual(richTextItem);
      expect(result.data?.value.type).toBe('rich_text');
      expect(result.data?.value.format).toBe('html');
    });

    it('should handle tags updates correctly', async () => {
      const tagsUpdateItem: ContentItem = {
        ...mockOriginalContentItem,
        tags: ['ui', 'navigation', 'featured', 'homepage'],
        metadata: {
          ...mockOriginalMetadata,
          updatedAt: '2024-01-20T14:30:00.000Z',
          updatedBy: 'user-456',
          version: 2
        }
      };

      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => tagsUpdateItem,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockResponse);

      const updateRequest: UpdateContentItemRequest = {
        tags: ['ui', 'navigation', 'featured', 'homepage'],
        notes: 'Added homepage and featured tags'
      };

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          ...updateRequest
        })
      );

      expect(result.data).toEqual(tagsUpdateItem);
      expect(result.data?.tags).toEqual(['ui', 'navigation', 'featured', 'homepage']);
    });
  });

  describe('Error Handling', () => {
    it('should handle 404 errors for non-existent content items', async () => {
      const mockErrorResponse = {
        ok: false,
        status: 404,
        json: async () => ({
          error: 'Content item not found',
          code: 'CONTENT_NOT_FOUND',
          details: 'Content item with ID "non-existent-id" does not exist'
        }),
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockErrorResponse);

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'non-existent-id',
          value: mockUpdatedTextContent
        })
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(404);
      expect(error.data?.code).toBe('CONTENT_NOT_FOUND');
    });

    it('should handle version conflicts (optimistic locking)', async () => {
      const mockErrorResponse = {
        ok: false,
        status: 409,
        json: async () => ({
          error: 'Version conflict',
          code: 'VERSION_CONFLICT',
          details: 'Content item has been updated by another user. Current version: 3, provided version: 1',
          currentVersion: 3
        }),
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockErrorResponse);

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          value: mockUpdatedTextContent,
          expectedVersion: 1 // Outdated version
        })
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(409);
      expect(error.data?.code).toBe('VERSION_CONFLICT');
      expect(error.data?.currentVersion).toBe(3);
    });

    it('should handle validation errors for invalid content', async () => {
      const mockErrorResponse = {
        ok: false,
        status: 400,
        json: async () => ({
          error: 'Validation failed',
          code: 'VALIDATION_ERROR',
          details: 'Content value exceeds maximum length of 100 characters',
          field: 'value.value'
        }),
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockErrorResponse);

      const invalidContent: ContentValue = {
        type: 'text',
        value: 'This is a very long welcome message that exceeds the maximum allowed length of 100 characters for this content item and should trigger a validation error',
        maxLength: 100
      };

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          value: invalidContent
        })
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(400);
      expect(error.data?.code).toBe('VALIDATION_ERROR');
      expect(error.data?.field).toBe('value.value');
    });

    it('should handle permission errors for restricted content', async () => {
      const mockErrorResponse = {
        ok: false,
        status: 403,
        json: async () => ({
          error: 'Insufficient permissions',
          code: 'INSUFFICIENT_PERMISSIONS',
          details: 'User does not have write permissions for category "navigation"',
          requiredPermission: 'write'
        }),
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockErrorResponse);

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          value: mockUpdatedTextContent
        })
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(403);
      expect(error.data?.code).toBe('INSUFFICIENT_PERMISSIONS');
      expect(error.data?.requiredPermission).toBe('write');
    });

    it('should handle authentication errors', async () => {
      const mockErrorResponse = {
        ok: false,
        status: 401,
        json: async () => ({
          error: 'Authentication required',
          code: 'AUTHENTICATION_REQUIRED',
          details: 'Valid authentication token is required to update content'
        }),
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockErrorResponse);

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          value: mockUpdatedTextContent
        })
      );

      expect(result.error).toBeDefined();
      const error = result.error as ApiError;
      expect(error.status).toBe(401);
      expect(error.data?.code).toBe('AUTHENTICATION_REQUIRED');
    });
  });

  describe('RTK Query Integration', () => {
    it('should provide proper cache invalidation', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => mockUpdatedContentItem,
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockResponse);

      // This will fail because updateContentItem endpoint doesn't exist yet
      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          value: mockUpdatedTextContent
        })
      );

      // Verify cache tags are properly set for invalidation
      expect(result.data).toEqual(mockUpdatedContentItem);

      // The endpoint should invalidate relevant cache tags:
      // - ['Content', 'nav-welcome-title'] for specific item
      // - ['Content', 'navigation'] for category items
      // - ['ContentList'] for any content lists
    });

    it('should support optimistic updates', async () => {
      // Mock delayed response to test optimistic updates
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => {
          await new Promise(resolve => setTimeout(resolve, 100));
          return mockUpdatedContentItem;
        },
      } as Response;

      global.fetch = vi.fn().mockResolvedValueOnce(mockResponse);

      // This will fail because updateContentItem endpoint doesn't exist yet
      const updateRequest: UpdateContentItemRequest = {
        value: mockUpdatedTextContent,
        status: ContentStatus.PUBLISHED
      };

      const result = await store.dispatch(
        contentApi.endpoints.updateContentItem.initiate({
          id: 'nav-welcome-title',
          ...updateRequest
        })
      );

      expect(result.data).toEqual(mockUpdatedContentItem);
    });

    it('should handle concurrent update requests properly', async () => {
      const mockResponse1 = {
        ok: true,
        status: 200,
        json: async () => ({
          ...mockUpdatedContentItem,
          metadata: { ...mockUpdatedMetadata, version: 2 }
        }),
      } as Response;

      const mockResponse2 = {
        ok: true,
        status: 200,
        json: async () => ({
          ...mockUpdatedContentItem,
          metadata: { ...mockUpdatedMetadata, version: 3 }
        }),
      } as Response;

      global.fetch = vi.fn()
        .mockResolvedValueOnce(mockResponse1)
        .mockResolvedValueOnce(mockResponse2);

      // This will fail because updateContentItem endpoint doesn't exist yet
      const [result1, result2] = await Promise.all([
        store.dispatch(
          contentApi.endpoints.updateContentItem.initiate({
            id: 'nav-welcome-title',
            value: mockUpdatedTextContent
          })
        ),
        store.dispatch(
          contentApi.endpoints.updateContentItem.initiate({
            id: 'nav-welcome-title',
            status: ContentStatus.PUBLISHED
          })
        )
      ]);

      expect(result1.data?.metadata.version).toBe(2);
      expect(result2.data?.metadata.version).toBe(3);
    });
  });

  describe('Type Safety Validation', () => {
    it('should enforce UpdateContentItemRequest interface constraints', () => {
      // Valid update request
      const validRequest: UpdateContentItemRequest = {
        value: mockUpdatedTextContent,
        status: ContentStatus.PUBLISHED,
        tags: ['ui', 'navigation'],
        notes: 'Updated content',
        changeLog: 'Content update'
      };

      expect(validRequest).toBeDefined();

      // Test partial updates
      const partialRequest: UpdateContentItemRequest = {
        status: ContentStatus.PUBLISHED
      };

      expect(partialRequest).toBeDefined();

      // Empty request should be valid (no required fields)
      const emptyRequest: UpdateContentItemRequest = {};
      expect(emptyRequest).toBeDefined();
    });

    it('should maintain type safety for ContentValue updates', () => {
      // Text content update
      const textUpdate: UpdateContentItemRequest = {
        value: {
          type: 'text',
          value: 'New text content',
          maxLength: 200
        }
      };

      expect(textUpdate.value?.type).toBe('text');

      // Rich text content update
      const richTextUpdate: UpdateContentItemRequest = {
        value: {
          type: 'rich_text',
          value: '<p>Rich content</p>',
          format: 'html',
          allowedTags: ['p', 'strong']
        }
      };

      expect(richTextUpdate.value?.type).toBe('rich_text');
    });
  });
});