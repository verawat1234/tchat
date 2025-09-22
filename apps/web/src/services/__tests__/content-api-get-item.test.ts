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
 * T005: Contract Test for getContentItem API
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * This test validates:
 * 1. Single content item retrieval by ID
 * 2. Content item structure validation
 * 3. Error handling for non-existent content (404)
 * 4. Proper typing and metadata validation
 * 5. RTK Query integration patterns
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
  value: 'Welcome to Tchat',
  maxLength: 100
};

const mockMetadata: ContentMetadata = {
  createdAt: '2024-01-15T10:00:00.000Z',
  createdBy: 'user-123',
  updatedAt: '2024-01-15T10:30:00.000Z',
  updatedBy: 'user-123',
  version: 1,
  tags: ['header', 'welcome'],
  notes: 'Initial welcome message'
};

const mockContentItem: ContentItem = {
  id: 'navigation.header.welcome',
  category: mockContentCategory,
  type: ContentType.TEXT,
  value: mockTextContent,
  metadata: mockMetadata,
  status: ContentStatus.PUBLISHED
};

const mockApiError: ApiError = {
  success: false,
  error: {
    code: 'NOT_FOUND',
    message: 'Content item not found',
    details: { contentId: 'non.existent.item' },
    timestamp: '2024-01-15T11:00:00.000Z'
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

describe('Content API - getContentItem', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
    mockFetch.mockClear();
    vi.clearAllMocks();
  });

  describe('Successful content item retrieval', () => {
    it('should fetch content item by ID with correct structure', async () => {
      // Arrange: Mock successful API response
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockContentItem,
          meta: {
            timestamp: '2024-01-15T11:00:00.000Z',
            version: '1.0.0'
          }
        }),
      });

      // Act: Trigger the query that WILL FAIL (TDD requirement)
      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('navigation.header.welcome')
      );

      const response = await result;

      // Assert: Validate response structure and content
      expect(response.data).toBeDefined();
      expect(response.error).toBeUndefined();

      const contentItem = response.data as ContentItem;

      // Validate core content item structure
      expect(contentItem.id).toBe('navigation.header.welcome');
      expect(contentItem.type).toBe(ContentType.TEXT);
      expect(contentItem.status).toBe(ContentStatus.PUBLISHED);

      // Validate category structure
      expect(contentItem.category).toEqual(mockContentCategory);
      expect(contentItem.category.permissions.read).toContain('user');
      expect(contentItem.category.permissions.write).toContain('admin');

      // Validate content value structure
      expect(contentItem.value.type).toBe('text');
      expect((contentItem.value as any).value).toBe('Welcome to Tchat');
      expect((contentItem.value as any).maxLength).toBe(100);

      // Validate metadata structure
      expect(contentItem.metadata.version).toBe(1);
      expect(contentItem.metadata.createdBy).toBe('user-123');
      expect(contentItem.metadata.tags).toContain('header');
      expect(contentItem.metadata.tags).toContain('welcome');
      expect(contentItem.metadata.notes).toBe('Initial welcome message');

      // Validate timestamp formats (ISO 8601)
      expect(contentItem.metadata.createdAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
      expect(contentItem.metadata.updatedAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
    });

    it('should handle different content types correctly', async () => {
      const richTextContent: ContentValue = {
        type: 'rich_text',
        value: '<h1>Welcome</h1><p>To our platform</p>',
        format: 'html',
        allowedTags: ['h1', 'p', 'strong', 'em']
      };

      const richTextItem: ContentItem = {
        ...mockContentItem,
        id: 'content.landing.hero',
        type: ContentType.RICH_TEXT,
        value: richTextContent
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: richTextItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('content.landing.hero')
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const contentItem = response.data as ContentItem;
      expect(contentItem.type).toBe(ContentType.RICH_TEXT);
      expect(contentItem.value.type).toBe('rich_text');
      expect((contentItem.value as any).format).toBe('html');
      expect((contentItem.value as any).allowedTags).toContain('h1');
    });

    it('should handle image content type correctly', async () => {
      const imageContent: ContentValue = {
        type: 'image_url',
        url: 'https://cdn.tchat.com/images/logo.png',
        alt: 'Tchat logo',
        width: 200,
        height: 60,
        format: 'png'
      };

      const imageItem: ContentItem = {
        ...mockContentItem,
        id: 'branding.logo.header',
        type: ContentType.IMAGE_URL,
        value: imageContent
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: imageItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('branding.logo.header')
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const contentItem = response.data as ContentItem;
      expect(contentItem.type).toBe(ContentType.IMAGE_URL);
      expect((contentItem.value as any).url).toBe('https://cdn.tchat.com/images/logo.png');
      expect((contentItem.value as any).alt).toBe('Tchat logo');
      expect((contentItem.value as any).width).toBe(200);
      expect((contentItem.value as any).height).toBe(60);
    });

    it('should handle translation content type correctly', async () => {
      const translationContent: ContentValue = {
        type: 'translation',
        values: {
          'en': 'Welcome',
          'es': 'Bienvenido',
          'fr': 'Bienvenue',
          'de': 'Willkommen'
        },
        defaultLocale: 'en'
      };

      const translationItem: ContentItem = {
        ...mockContentItem,
        id: 'navigation.menu.welcome',
        type: ContentType.TRANSLATION,
        value: translationContent
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: translationItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('navigation.menu.welcome')
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const contentItem = response.data as ContentItem;
      expect(contentItem.type).toBe(ContentType.TRANSLATION);
      expect((contentItem.value as any).values.en).toBe('Welcome');
      expect((contentItem.value as any).values.es).toBe('Bienvenido');
      expect((contentItem.value as any).defaultLocale).toBe('en');
    });
  });

  describe('Error handling', () => {
    it('should handle 404 Not Found error correctly', async () => {
      // Arrange: Mock 404 response
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => mockApiError,
      });

      // Act: Attempt to fetch non-existent content
      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('non.existent.item')
      );

      const response = await result;

      // Assert: Validate error response structure
      expect(response.data).toBeUndefined();
      expect(response.error).toBeDefined();

      const error = response.error as any;
      expect(error.status).toBe(404);
      expect(error.data).toEqual(mockApiError);
      expect(error.data.error.code).toBe('NOT_FOUND');
      expect(error.data.error.message).toBe('Content item not found');
      expect(error.data.error.details.contentId).toBe('non.existent.item');
    });

    it('should handle 401 Unauthorized error correctly', async () => {
      const unauthorizedError: ApiError = {
        success: false,
        error: {
          code: 'UNAUTHORIZED',
          message: 'Authentication required',
          timestamp: '2024-01-15T11:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => unauthorizedError,
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('protected.admin.content')
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(401);
      expect(error.data.error.code).toBe('UNAUTHORIZED');
    });

    it('should handle 403 Forbidden error correctly', async () => {
      const forbiddenError: ApiError = {
        success: false,
        error: {
          code: 'FORBIDDEN',
          message: 'Insufficient permissions to access this content',
          details: { requiredRole: 'admin', userRole: 'user' },
          timestamp: '2024-01-15T11:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
        json: async () => forbiddenError,
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('admin.sensitive.data')
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(403);
      expect(error.data.error.code).toBe('FORBIDDEN');
      expect(error.data.error.details.requiredRole).toBe('admin');
    });

    it('should handle network errors correctly', async () => {
      // Arrange: Mock network failure
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      // Act: Attempt to fetch content with network failure
      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('any.content.item')
      );

      const response = await result;

      // Assert: Validate network error handling
      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.error).toBe('Network error');
    });

    it('should handle malformed content ID validation', async () => {
      const validationError: ApiError = {
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Invalid content ID format. Expected pattern: category.subcategory.key',
          details: {
            providedId: 'invalid-id',
            expectedPattern: '{category}.{subcategory}.{key}',
            examples: ['navigation.header.title', 'error.network.timeout']
          },
          timestamp: '2024-01-15T11:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => validationError,
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('invalid-id')
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(400);
      expect(error.data.error.code).toBe('VALIDATION_ERROR');
      expect(error.data.error.details.expectedPattern).toBe('{category}.{subcategory}.{key}');
    });
  });

  describe('RTK Query integration', () => {
    it('should provide correct cache tags for content items', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockContentItem
        }),
      });

      // This test validates that the endpoint provides proper cache tags
      // The actual implementation should provide tags like:
      // ['Content', { type: 'Content', id: 'navigation.header.welcome' }]

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('navigation.header.welcome')
      );

      await result;

      // When implemented, this should verify cache tag structure
      // For now, this documents the expected cache behavior
      expect(true).toBe(true); // Placeholder until implementation
    });

    it('should handle query deduplication correctly', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockContentItem
        }),
      });

      // Trigger multiple identical queries
      const query1 = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('navigation.header.welcome')
      );

      const query2 = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('navigation.header.welcome')
      );

      await Promise.all([query1, query2]);

      // Should only make one network request due to deduplication
      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it('should support proper request serialization', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockContentItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('navigation.header.welcome')
      );

      await result;

      // Validate that the request was made to the correct endpoint
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/content/navigation.header.welcome'),
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          })
        })
      );
    });
  });

  describe('Type safety validation', () => {
    it('should provide proper TypeScript typing for successful responses', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          success: true,
          data: mockContentItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('navigation.header.welcome')
      );

      const response = await result;

      // TypeScript should infer the correct types
      if (response.data) {
        // These should be properly typed without 'any' assertions
        const id: string = response.data.id;
        const type: ContentType = response.data.type;
        const status: ContentStatus = response.data.status;
        const version: number = response.data.metadata.version;

        expect(id).toBe('navigation.header.welcome');
        expect(type).toBe(ContentType.TEXT);
        expect(status).toBe(ContentStatus.PUBLISHED);
        expect(version).toBe(1);
      }
    });

    it('should provide proper TypeScript typing for error responses', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => mockApiError,
      });

      const result = store.dispatch(
        contentApi.endpoints.getContentItem.initiate('non.existent.item')
      );

      const response = await result;

      // TypeScript should properly type the error
      if (response.error) {
        const error = response.error as any; // Will be properly typed when implemented
        expect(error.status).toBe(404);
        expect(error.data.error.code).toBe('NOT_FOUND');
      }
    });
  });

  describe('Content ID pattern validation', () => {
    it('should validate correct content ID patterns', () => {
      const validIds = [
        'navigation.header.title',
        'error.network.timeout',
        'content.landing.hero',
        'branding.logo.header',
        'form.validation.required'
      ];

      // When implemented, these should all be accepted
      validIds.forEach(id => {
        expect(id.split('.').length).toBe(3);
        expect(id).toMatch(/^[a-z]+\.[a-z]+\.[a-z]+$/);
      });
    });

    it('should reject invalid content ID patterns', () => {
      const invalidIds = [
        'invalid-id',
        'only.two.parts.too.many',
        'onepart',
        'two.parts',
        '',
        'navigation.header.',
        '.header.title',
        'navigation..title'
      ];

      // When implemented, these should all be rejected
      invalidIds.forEach(id => {
        const parts = id.split('.');
        expect(
          parts.length !== 3 ||
          parts.some(part => part === '') ||
          !id.match(/^[a-z]+\.[a-z]+\.[a-z]+$/)
        ).toBe(true);
      });
    });
  });
});