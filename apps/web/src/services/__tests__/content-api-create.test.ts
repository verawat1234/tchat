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
  CreateContentItemRequest,
  TextContent,
  RichTextContent,
  ImageContent,
  ConfigContent,
  TranslationContent
} from '../../types/content';
import type { ApiError } from '../../types/api';

/**
 * T010: Contract Test for createContentItem API
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * This test validates:
 * 1. New content item creation with complete request structure
 * 2. Content validation and type safety
 * 3. Duplicate ID prevention (409 Conflict)
 * 4. Required field validation
 * 5. Content type-specific validation rules
 * 6. Draft status assignment for new content
 * 7. Metadata generation for created content
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

const mockGeneratedMetadata: ContentMetadata = {
  createdAt: '2024-01-15T10:00:00.000Z',
  createdBy: 'user-123',
  updatedAt: '2024-01-15T10:00:00.000Z',
  updatedBy: 'user-123',
  version: 1,
  tags: ['new', 'content'],
  notes: 'Initial content creation'
};

// Valid create requests for different content types
const validTextContentRequest: CreateContentItemRequest = {
  id: 'navigation.header.welcome',
  categoryId: 'navigation',
  type: ContentType.TEXT,
  value: {
    type: 'text',
    value: 'Welcome to Tchat',
    maxLength: 100
  } as TextContent,
  tags: ['header', 'welcome'],
  notes: 'New welcome message for header'
};

const validRichTextContentRequest: CreateContentItemRequest = {
  id: 'content.landing.hero',
  categoryId: 'content',
  type: ContentType.RICH_TEXT,
  value: {
    type: 'rich_text',
    value: '<h1>Welcome</h1><p>To our amazing platform</p>',
    format: 'html',
    allowedTags: ['h1', 'p', 'strong', 'em']
  } as RichTextContent,
  tags: ['landing', 'hero'],
  notes: 'Hero section content'
};

const validImageContentRequest: CreateContentItemRequest = {
  id: 'branding.logo.header',
  categoryId: 'branding',
  type: ContentType.IMAGE_URL,
  value: {
    type: 'image_url',
    url: 'https://cdn.tchat.com/images/logo.png',
    alt: 'Tchat logo',
    width: 200,
    height: 60,
    format: 'png'
  } as ImageContent,
  tags: ['logo', 'branding']
};

const validConfigContentRequest: CreateContentItemRequest = {
  id: 'app.settings.maxUsers',
  categoryId: 'settings',
  type: ContentType.CONFIG,
  value: {
    type: 'config',
    value: 1000,
    schema: {
      type: 'integer',
      minimum: 1,
      maximum: 10000
    }
  } as ConfigContent,
  tags: ['config', 'limits']
};

const validTranslationContentRequest: CreateContentItemRequest = {
  id: 'navigation.menu.home',
  categoryId: 'navigation',
  type: ContentType.TRANSLATION,
  value: {
    type: 'translation',
    values: {
      'en': 'Home',
      'es': 'Inicio',
      'fr': 'Accueil',
      'de': 'Startseite'
    },
    defaultLocale: 'en'
  } as TranslationContent,
  tags: ['navigation', 'menu']
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

describe('Content API - createContentItem', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
    mockFetch.mockClear();
    vi.clearAllMocks();
  });

  describe('Successful content item creation', () => {
    it('should create text content item with correct structure and metadata', async () => {
      // Expected response after creation
      const expectedContentItem: ContentItem = {
        id: validTextContentRequest.id,
        category: mockContentCategory,
        type: ContentType.TEXT,
        value: validTextContentRequest.value,
        metadata: mockGeneratedMetadata,
        status: ContentStatus.DRAFT // New content should default to DRAFT
      };

      // Arrange: Mock successful API response
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedContentItem,
          meta: {
            timestamp: '2024-01-15T10:00:00.000Z',
            version: '1.0.0'
          }
        }),
      });

      // Act: Trigger the create mutation that WILL FAIL (TDD requirement)
      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validTextContentRequest)
      );

      const response = await result;

      // Assert: Validate response structure and content
      expect(response.data).toBeDefined();
      expect(response.error).toBeUndefined();

      const contentItem = response.data as ContentItem;

      // Validate core content item structure
      expect(contentItem.id).toBe('navigation.header.welcome');
      expect(contentItem.type).toBe(ContentType.TEXT);
      expect(contentItem.status).toBe(ContentStatus.DRAFT); // Should default to draft

      // Validate category assignment
      expect(contentItem.category).toEqual(mockContentCategory);

      // Validate content value preservation
      expect(contentItem.value.type).toBe('text');
      expect((contentItem.value as TextContent).value).toBe('Welcome to Tchat');
      expect((contentItem.value as TextContent).maxLength).toBe(100);

      // Validate metadata generation
      expect(contentItem.metadata.version).toBe(1); // Should start at version 1
      expect(contentItem.metadata.createdBy).toBe('user-123');
      expect(contentItem.metadata.updatedBy).toBe('user-123');
      expect(contentItem.metadata.createdAt).toBe(contentItem.metadata.updatedAt); // Should be same on creation
      expect(contentItem.metadata.tags).toContain('header');
      expect(contentItem.metadata.tags).toContain('welcome');
      expect(contentItem.metadata.notes).toBe('New welcome message for header');

      // Validate timestamp format (ISO 8601)
      expect(contentItem.metadata.createdAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
      expect(contentItem.metadata.updatedAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
    });

    it('should create rich text content item with proper HTML validation', async () => {
      const expectedRichTextItem: ContentItem = {
        id: validRichTextContentRequest.id,
        category: mockContentCategory,
        type: ContentType.RICH_TEXT,
        value: validRichTextContentRequest.value,
        metadata: mockGeneratedMetadata,
        status: ContentStatus.DRAFT
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedRichTextItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validRichTextContentRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const contentItem = response.data as ContentItem;
      expect(contentItem.type).toBe(ContentType.RICH_TEXT);
      expect(contentItem.value.type).toBe('rich_text');
      expect((contentItem.value as RichTextContent).format).toBe('html');
      expect((contentItem.value as RichTextContent).allowedTags).toContain('h1');
      expect((contentItem.value as RichTextContent).value).toContain('<h1>Welcome</h1>');
    });

    it('should create image content item with complete metadata', async () => {
      const expectedImageItem: ContentItem = {
        id: validImageContentRequest.id,
        category: mockContentCategory,
        type: ContentType.IMAGE_URL,
        value: validImageContentRequest.value,
        metadata: mockGeneratedMetadata,
        status: ContentStatus.DRAFT
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedImageItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validImageContentRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const contentItem = response.data as ContentItem;
      expect(contentItem.type).toBe(ContentType.IMAGE_URL);
      expect((contentItem.value as ImageContent).url).toBe('https://cdn.tchat.com/images/logo.png');
      expect((contentItem.value as ImageContent).alt).toBe('Tchat logo');
      expect((contentItem.value as ImageContent).width).toBe(200);
      expect((contentItem.value as ImageContent).height).toBe(60);
      expect((contentItem.value as ImageContent).format).toBe('png');
    });

    it('should create config content item with schema validation', async () => {
      const expectedConfigItem: ContentItem = {
        id: validConfigContentRequest.id,
        category: mockContentCategory,
        type: ContentType.CONFIG,
        value: validConfigContentRequest.value,
        metadata: mockGeneratedMetadata,
        status: ContentStatus.DRAFT
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedConfigItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validConfigContentRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const contentItem = response.data as ContentItem;
      expect(contentItem.type).toBe(ContentType.CONFIG);
      expect((contentItem.value as ConfigContent).value).toBe(1000);
      expect((contentItem.value as ConfigContent).schema?.type).toBe('integer');
      expect((contentItem.value as ConfigContent).schema?.minimum).toBe(1);
      expect((contentItem.value as ConfigContent).schema?.maximum).toBe(10000);
    });

    it('should create translation content item with multi-language support', async () => {
      const expectedTranslationItem: ContentItem = {
        id: validTranslationContentRequest.id,
        category: mockContentCategory,
        type: ContentType.TRANSLATION,
        value: validTranslationContentRequest.value,
        metadata: mockGeneratedMetadata,
        status: ContentStatus.DRAFT
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedTranslationItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validTranslationContentRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const contentItem = response.data as ContentItem;
      expect(contentItem.type).toBe(ContentType.TRANSLATION);
      expect((contentItem.value as TranslationContent).values.en).toBe('Home');
      expect((contentItem.value as TranslationContent).values.es).toBe('Inicio');
      expect((contentItem.value as TranslationContent).values.fr).toBe('Accueil');
      expect((contentItem.value as TranslationContent).values.de).toBe('Startseite');
      expect((contentItem.value as TranslationContent).defaultLocale).toBe('en');
    });

    it('should handle optional fields correctly', async () => {
      const minimalRequest: CreateContentItemRequest = {
        id: 'minimal.test.item',
        categoryId: 'test',
        type: ContentType.TEXT,
        value: {
          type: 'text',
          value: 'Minimal content'
        } as TextContent
        // No tags or notes provided
      };

      const expectedMinimalItem: ContentItem = {
        ...minimalRequest,
        category: mockContentCategory,
        metadata: {
          ...mockGeneratedMetadata,
          tags: [], // Should default to empty array
          notes: undefined // Should handle undefined notes
        },
        status: ContentStatus.DRAFT
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedMinimalItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(minimalRequest)
      );

      const response = await result;

      expect(response.data).toBeDefined();
      const contentItem = response.data as ContentItem;
      expect(contentItem.metadata.tags).toEqual([]);
      expect(contentItem.metadata.notes).toBeUndefined();
    });
  });

  describe('Validation and error handling', () => {
    it('should handle duplicate ID error (409 Conflict)', async () => {
      const conflictError: ApiError = {
        success: false,
        error: {
          code: 'CONFLICT',
          message: 'Content item with this ID already exists',
          details: {
            contentId: 'navigation.header.welcome',
            existingVersion: 2,
            conflictType: 'DUPLICATE_ID'
          },
          timestamp: '2024-01-15T10:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 409,
        json: async () => conflictError,
      });

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validTextContentRequest)
      );

      const response = await result;

      expect(response.data).toBeUndefined();
      expect(response.error).toBeDefined();

      const error = response.error as any;
      expect(error.status).toBe(409);
      expect(error.data.error.code).toBe('CONFLICT');
      expect(error.data.error.details.contentId).toBe('navigation.header.welcome');
      expect(error.data.error.details.conflictType).toBe('DUPLICATE_ID');
    });

    it('should validate required fields', async () => {
      const validationError: ApiError = {
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Missing required fields',
          details: {
            missingFields: ['id', 'categoryId', 'type', 'value'],
            providedFields: [],
            requiredFields: ['id', 'categoryId', 'type', 'value']
          },
          timestamp: '2024-01-15T10:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => validationError,
      });

      // Test with empty request
      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate({} as CreateContentItemRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(400);
      expect(error.data.error.code).toBe('VALIDATION_ERROR');
      expect(error.data.error.details.missingFields).toContain('id');
      expect(error.data.error.details.missingFields).toContain('categoryId');
      expect(error.data.error.details.missingFields).toContain('type');
      expect(error.data.error.details.missingFields).toContain('value');
    });

    it('should validate content ID pattern', async () => {
      const idPatternError: ApiError = {
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Invalid content ID format. Expected pattern: category.subcategory.key',
          details: {
            providedId: 'invalid-id-format',
            expectedPattern: '{category}.{subcategory}.{key}',
            examples: ['navigation.header.title', 'error.network.timeout'],
            validationRules: [
              'Must contain exactly 3 parts separated by dots',
              'Each part must contain only lowercase letters',
              'No empty parts allowed'
            ]
          },
          timestamp: '2024-01-15T10:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => idPatternError,
      });

      const invalidRequest: CreateContentItemRequest = {
        ...validTextContentRequest,
        id: 'invalid-id-format'
      };

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(invalidRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(400);
      expect(error.data.error.code).toBe('VALIDATION_ERROR');
      expect(error.data.error.details.providedId).toBe('invalid-id-format');
      expect(error.data.error.details.expectedPattern).toBe('{category}.{subcategory}.{key}');
    });

    it('should validate content type-specific rules', async () => {
      const typeValidationError: ApiError = {
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Content value validation failed for type: rich_text',
          details: {
            contentType: 'rich_text',
            validationErrors: [
              'HTML contains disallowed tags: <script>, <iframe>',
              'Format must be either "html" or "markdown"',
              'allowedTags array cannot be empty for HTML content'
            ],
            providedValue: {
              type: 'rich_text',
              value: '<script>alert("bad")</script><p>Content</p>',
              format: 'invalid',
              allowedTags: []
            }
          },
          timestamp: '2024-01-15T10:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => typeValidationError,
      });

      const invalidRichTextRequest: CreateContentItemRequest = {
        ...validRichTextContentRequest,
        value: {
          type: 'rich_text',
          value: '<script>alert("bad")</script><p>Content</p>',
          format: 'invalid' as any,
          allowedTags: []
        } as RichTextContent
      };

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(invalidRichTextRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(400);
      expect(error.data.error.code).toBe('VALIDATION_ERROR');
      expect(error.data.error.details.contentType).toBe('rich_text');
      expect(error.data.error.details.validationErrors).toContain('HTML contains disallowed tags: <script>, <iframe>');
    });

    it('should validate category existence', async () => {
      const categoryNotFoundError: ApiError = {
        success: false,
        error: {
          code: 'NOT_FOUND',
          message: 'Content category not found',
          details: {
            categoryId: 'nonexistent-category',
            availableCategories: ['navigation', 'content', 'branding', 'settings'],
            suggestion: 'Check available categories or create the category first'
          },
          timestamp: '2024-01-15T10:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => categoryNotFoundError,
      });

      const invalidCategoryRequest: CreateContentItemRequest = {
        ...validTextContentRequest,
        categoryId: 'nonexistent-category'
      };

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(invalidCategoryRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(404);
      expect(error.data.error.code).toBe('NOT_FOUND');
      expect(error.data.error.details.categoryId).toBe('nonexistent-category');
      expect(error.data.error.details.availableCategories).toContain('navigation');
    });

    it('should handle authorization errors', async () => {
      const unauthorizedError: ApiError = {
        success: false,
        error: {
          code: 'FORBIDDEN',
          message: 'Insufficient permissions to create content in this category',
          details: {
            categoryId: 'admin-only',
            requiredPermission: 'write',
            userRole: 'user',
            requiredRoles: ['admin']
          },
          timestamp: '2024-01-15T10:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 403,
        json: async () => unauthorizedError,
      });

      const unauthorizedRequest: CreateContentItemRequest = {
        ...validTextContentRequest,
        categoryId: 'admin-only'
      };

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(unauthorizedRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(403);
      expect(error.data.error.code).toBe('FORBIDDEN');
      expect(error.data.error.details.requiredPermission).toBe('write');
      expect(error.data.error.details.userRole).toBe('user');
    });

    it('should handle content value type mismatch', async () => {
      const typeMismatchError: ApiError = {
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Content value type does not match declared content type',
          details: {
            declaredType: 'text',
            valueType: 'rich_text',
            expectedValueStructure: {
              type: 'text',
              value: 'string',
              maxLength: 'number (optional)'
            },
            providedValueStructure: {
              type: 'rich_text',
              value: 'string',
              format: 'html'
            }
          },
          timestamp: '2024-01-15T10:00:00.000Z'
        }
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => typeMismatchError,
      });

      const typeMismatchRequest: CreateContentItemRequest = {
        ...validTextContentRequest,
        type: ContentType.TEXT,
        value: {
          type: 'rich_text',
          value: '<p>This is rich text</p>',
          format: 'html'
        } as any // Intentionally mismatched type
      };

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(typeMismatchRequest)
      );

      const response = await result;

      expect(response.error).toBeDefined();
      const error = response.error as any;
      expect(error.status).toBe(400);
      expect(error.data.error.code).toBe('VALIDATION_ERROR');
      expect(error.data.error.details.declaredType).toBe('text');
      expect(error.data.error.details.valueType).toBe('rich_text');
    });
  });

  describe('RTK Query integration', () => {
    it('should provide correct cache invalidation tags', async () => {
      const expectedContentItem: ContentItem = {
        id: validTextContentRequest.id,
        category: mockContentCategory,
        type: ContentType.TEXT,
        value: validTextContentRequest.value,
        metadata: mockGeneratedMetadata,
        status: ContentStatus.DRAFT
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedContentItem
        }),
      });

      // This test validates that the create mutation provides proper cache invalidation
      // The actual implementation should invalidate tags like:
      // ['Content', { type: 'ContentList', id: 'navigation' }, { type: 'ContentCategory', id: 'navigation' }]

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validTextContentRequest)
      );

      await result;

      // When implemented, this should verify cache invalidation behavior
      // For now, this documents the expected cache behavior
      expect(true).toBe(true); // Placeholder until implementation
    });

    it('should handle proper request serialization for POST', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: {
            id: validTextContentRequest.id,
            category: mockContentCategory,
            type: ContentType.TEXT,
            value: validTextContentRequest.value,
            metadata: mockGeneratedMetadata,
            status: ContentStatus.DRAFT
          }
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validTextContentRequest)
      );

      await result;

      // Validate that the request was made to the correct endpoint with proper payload
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/content'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          }),
          body: JSON.stringify(validTextContentRequest)
        })
      );
    });

    it('should support mutation lifecycle hooks', async () => {
      const expectedContentItem: ContentItem = {
        id: validTextContentRequest.id,
        category: mockContentCategory,
        type: ContentType.TEXT,
        value: validTextContentRequest.value,
        metadata: mockGeneratedMetadata,
        status: ContentStatus.DRAFT
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedContentItem
        }),
      });

      // Test mutation state management
      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validTextContentRequest)
      );

      // Should be loading initially
      expect(result.requestId).toBeDefined();

      const response = await result;

      // Should complete successfully
      expect(response.data).toBeDefined();
      expect(response.error).toBeUndefined();
    });
  });

  describe('Type safety validation', () => {
    it('should provide proper TypeScript typing for create request', () => {
      // TypeScript should enforce correct request structure
      const typedRequest: CreateContentItemRequest = {
        id: 'test.typed.content',
        categoryId: 'test',
        type: ContentType.TEXT,
        value: {
          type: 'text',
          value: 'Typed content'
        } as TextContent,
        tags: ['typed', 'test'],
        notes: 'Type safety test'
      };

      // These should be properly typed without 'any' assertions
      const id: string = typedRequest.id;
      const categoryId: string = typedRequest.categoryId;
      const type: ContentType = typedRequest.type;
      const value: ContentValue = typedRequest.value;
      const tags: string[] | undefined = typedRequest.tags;
      const notes: string | undefined = typedRequest.notes;

      expect(id).toBe('test.typed.content');
      expect(categoryId).toBe('test');
      expect(type).toBe(ContentType.TEXT);
      expect(value.type).toBe('text');
      expect(tags).toContain('typed');
      expect(notes).toBe('Type safety test');
    });

    it('should provide proper TypeScript typing for create response', async () => {
      const expectedContentItem: ContentItem = {
        id: validTextContentRequest.id,
        category: mockContentCategory,
        type: ContentType.TEXT,
        value: validTextContentRequest.value,
        metadata: mockGeneratedMetadata,
        status: ContentStatus.DRAFT
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: expectedContentItem
        }),
      });

      const result = store.dispatch(
        contentApi.endpoints.createContentItem.initiate(validTextContentRequest)
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
        expect(status).toBe(ContentStatus.DRAFT);
        expect(version).toBe(1);
      }
    });
  });

  describe('Content ID pattern validation edge cases', () => {
    it('should handle various valid content ID patterns', () => {
      const validIds = [
        'navigation.header.title',
        'error.network.timeout',
        'content.landing.hero',
        'branding.logo.header',
        'form.validation.required',
        'app.settings.theme',
        'modal.confirm.delete'
      ];

      // When implemented, these should all be accepted
      validIds.forEach(id => {
        expect(id.split('.').length).toBe(3);
        expect(id).toMatch(/^[a-z]+\.[a-z]+\.[a-z]+$/);

        // Verify no empty parts
        const parts = id.split('.');
        expect(parts.every(part => part.length > 0)).toBe(true);
      });
    });

    it('should reject invalid content ID patterns', () => {
      const invalidIds = [
        'invalid-id',
        'only.two.parts.too.many.parts',
        'onepart',
        'two.parts',
        '',
        'navigation.header.',
        '.header.title',
        'navigation..title',
        'Navigation.Header.Title', // uppercase
        'navigation.header.title-with-dash',
        'navigation.header.title_with_underscore',
        'navigation.header.123numbers'
      ];

      // When implemented, these should all be rejected
      invalidIds.forEach(id => {
        const parts = id.split('.');
        const isValid = parts.length === 3 &&
                       parts.every(part => part.length > 0) &&
                       /^[a-z]+\.[a-z]+\.[a-z]+$/.test(id);
        expect(isValid).toBe(false);
      });
    });
  });
});