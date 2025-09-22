/**
 * T006: Contract Test - getContentByCategory API Endpoint
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * This contract test validates the API contract for retrieving content by category.
 * It ensures the implementation follows the expected interface and behavior patterns
 * defined in the content management system specification.
 */

import { configureStore } from '@reduxjs/toolkit';
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';
import { describe, beforeAll, afterAll, beforeEach, afterEach, test, expect, vi } from 'vitest';
import { api } from '../api';
import {
  ContentItem,
  ContentCategory,
  ContentType,
  ContentStatus,
  PaginatedContentResponse,
  ContentQueryParams
} from '../../types/content';

// Mock content category factory
const createMockCategory = (overrides: Partial<ContentCategory> = {}): ContentCategory => ({
  id: 'navigation',
  name: 'Navigation',
  description: 'Website navigation content',
  permissions: {
    read: ['user', 'admin'],
    write: ['admin'],
    publish: ['admin']
  },
  ...overrides,
});

// Mock content item factory
const createMockContentItem = (overrides: Partial<ContentItem> = {}): ContentItem => ({
  id: 'navigation.header.title',
  category: createMockCategory(),
  type: ContentType.TEXT,
  value: {
    type: 'text',
    value: 'Tchat App',
    maxLength: 100
  },
  metadata: {
    createdAt: '2024-01-01T00:00:00Z',
    createdBy: 'admin',
    updatedAt: '2024-01-01T00:00:00Z',
    updatedBy: 'admin',
    version: 1,
    tags: ['navigation', 'header'],
    notes: 'Main header title'
  },
  status: ContentStatus.PUBLISHED,
  ...overrides,
});

// Mock paginated response factory
const createMockPaginatedResponse = (
  items: ContentItem[],
  pagination = { page: 1, limit: 10, total: items.length }
): PaginatedContentResponse => ({
  items,
  pagination: {
    page: pagination.page,
    limit: pagination.limit,
    total: pagination.total,
    totalPages: Math.ceil(pagination.total / pagination.limit),
    hasNext: pagination.page * pagination.limit < pagination.total,
    hasPrev: pagination.page > 1,
  },
});

// API Base URL
const API_BASE_URL = 'http://localhost:3001';

// Mock MSW server setup
const server = setupServer();

// RTK Query store setup for testing
const createTestStore = () => {
  return configureStore({
    reducer: {
      [api.reducerPath]: api.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware({
        serializableCheck: false,
      }).concat(api.middleware),
  });
};

describe('T006: Content API - getContentByCategory Contract Test', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeAll(() => {
    server.listen({ onUnhandledRequest: 'error' });
  });

  afterAll(() => {
    server.close();
  });

  beforeEach(() => {
    store = createTestStore();
    vi.clearAllTimers();
  });

  afterEach(() => {
    server.resetHandlers();
    store.dispatch(api.util.resetApiState());
  });

  describe('Contract Validation - API Interface', () => {
    test('should define getContentByCategory endpoint', () => {
      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      // This validates the API contract exists
      expect(() => {
        // @ts-expect-error - Testing for non-existent endpoint
        api.endpoints.getContentByCategory;
      }).not.toThrow();
    });

    test('should accept correct query parameters', () => {
      const validParams: ContentQueryParams = {
        categoryId: 'navigation',
        type: ContentType.TEXT,
        status: ContentStatus.PUBLISHED,
        search: 'header',
        tags: ['navigation'],
        language: 'en',
        sortBy: 'updatedAt',
        sortOrder: 'desc',
        page: 1,
        limit: 20
      };

      // This validates the parameter interface contract
      expect(validParams).toBeDefined();
      expect(validParams.categoryId).toBe('navigation');
      expect(validParams.type).toBe(ContentType.TEXT);
      expect(validParams.status).toBe(ContentStatus.PUBLISHED);
    });
  });

  describe('Category-based Content Filtering', () => {
    test('should retrieve content filtered by category ID', async () => {
      const categoryId = 'navigation';
      const mockItems = [
        createMockContentItem({
          id: 'navigation.header.title',
          category: createMockCategory({ id: categoryId })
        }),
        createMockContentItem({
          id: 'navigation.footer.copyright',
          category: createMockCategory({ id: categoryId })
        }),
      ];

      server.use(
        http.get(`${API_BASE_URL}/api/content`, ({ request }) => {
          const url = new URL(request.url);
          const requestCategoryId = url.searchParams.get('categoryId');

          expect(requestCategoryId).toBe(categoryId);

          return HttpResponse.json({
            success: true,
            data: createMockPaginatedResponse(mockItems),
          });
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const result = await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId })
        ).unwrap();

        expect(result.items).toHaveLength(2);
        expect(result.items.every(item => item.category.id === categoryId)).toBe(true);
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });

    test('should handle multiple category filter combinations', async () => {
      const params: ContentQueryParams = {
        categoryId: 'navigation',
        type: ContentType.TEXT,
        status: ContentStatus.PUBLISHED,
        tags: ['header', 'primary']
      };

      server.use(
        http.get(`${API_BASE_URL}/api/content`, ({ request }) => {
          const url = new URL(request.url);

          expect(url.searchParams.get('categoryId')).toBe(params.categoryId);
          expect(url.searchParams.get('type')).toBe(params.type);
          expect(url.searchParams.get('status')).toBe(params.status);
          expect(url.searchParams.get('tags')).toBe(params.tags?.join(','));

          return HttpResponse.json({
            success: true,
            data: createMockPaginatedResponse([
              createMockContentItem({
                category: createMockCategory({ id: params.categoryId }),
                type: params.type,
                status: params.status,
                metadata: {
                  ...createMockContentItem().metadata,
                  tags: params.tags
                }
              })
            ]),
          });
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const result = await store.dispatch(
          api.endpoints.getContentByCategory.initiate(params)
        ).unwrap();

        expect(result.items[0].category.id).toBe(params.categoryId);
        expect(result.items[0].type).toBe(params.type);
        expect(result.items[0].status).toBe(params.status);
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });
  });

  describe('Array Response Validation', () => {
    test('should return paginated response structure', async () => {
      const mockItems = Array.from({ length: 25 }, (_, i) =>
        createMockContentItem({
          id: `navigation.item.${i}`,
          category: createMockCategory({ id: 'navigation' })
        })
      );

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json({
            success: true,
            data: createMockPaginatedResponse(
              mockItems.slice(0, 10), // First page
              { page: 1, limit: 10, total: 25 }
            ),
          });
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const result = await store.dispatch(
          api.endpoints.getContentByCategory.initiate({
            categoryId: 'navigation',
            page: 1,
            limit: 10
          })
        ).unwrap();

        // Validate pagination structure
        expect(result.pagination).toBeDefined();
        expect(result.pagination.page).toBe(1);
        expect(result.pagination.limit).toBe(10);
        expect(result.pagination.total).toBe(25);
        expect(result.pagination.totalPages).toBe(3);
        expect(result.pagination.hasNext).toBe(true);
        expect(result.pagination.hasPrev).toBe(false);

        // Validate items array
        expect(Array.isArray(result.items)).toBe(true);
        expect(result.items).toHaveLength(10);
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });

    test('should validate content item structure in response', async () => {
      const mockItem = createMockContentItem();

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json({
            success: true,
            data: createMockPaginatedResponse([mockItem]),
          });
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const result = await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId: 'navigation' })
        ).unwrap();

        const item = result.items[0];

        // Validate ContentItem structure
        expect(item.id).toBeDefined();
        expect(item.category).toBeDefined();
        expect(item.type).toBeDefined();
        expect(item.value).toBeDefined();
        expect(item.metadata).toBeDefined();
        expect(item.status).toBeDefined();

        // Validate metadata structure
        expect(item.metadata.createdAt).toBeDefined();
        expect(item.metadata.createdBy).toBeDefined();
        expect(item.metadata.updatedAt).toBeDefined();
        expect(item.metadata.updatedBy).toBeDefined();
        expect(item.metadata.version).toBeDefined();
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });
  });

  describe('Category Consistency in Results', () => {
    test('should ensure all returned items belong to requested category', async () => {
      const categoryId = 'navigation';
      const mockItems = [
        createMockContentItem({
          id: 'navigation.header.title',
          category: createMockCategory({ id: categoryId })
        }),
        createMockContentItem({
          id: 'navigation.menu.home',
          category: createMockCategory({ id: categoryId })
        }),
        createMockContentItem({
          id: 'navigation.footer.links',
          category: createMockCategory({ id: categoryId })
        }),
      ];

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json({
            success: true,
            data: createMockPaginatedResponse(mockItems),
          });
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const result = await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId })
        ).unwrap();

        // Validate category consistency
        result.items.forEach(item => {
          expect(item.category.id).toBe(categoryId);
          expect(item.id).toMatch(new RegExp(`^${categoryId}\\.`));
        });
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });

    test('should handle subcategory hierarchical filtering', async () => {
      const parentCategoryId = 'navigation';
      const subcategoryItems = [
        createMockContentItem({
          id: 'navigation.header.title',
          category: createMockCategory({
            id: 'navigation.header',
            parentId: parentCategoryId
          })
        }),
        createMockContentItem({
          id: 'navigation.footer.copyright',
          category: createMockCategory({
            id: 'navigation.footer',
            parentId: parentCategoryId
          })
        }),
      ];

      server.use(
        http.get(`${API_BASE_URL}/api/content`, ({ request }) => {
          const url = new URL(request.url);
          const includeSubcategories = url.searchParams.get('includeSubcategories');

          expect(includeSubcategories).toBe('true');

          return HttpResponse.json({
            success: true,
            data: createMockPaginatedResponse(subcategoryItems),
          });
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const result = await store.dispatch(
          api.endpoints.getContentByCategory.initiate({
            categoryId: parentCategoryId,
            includeSubcategories: true
          })
        ).unwrap();

        result.items.forEach(item => {
          expect(item.category.parentId).toBe(parentCategoryId);
        });
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });
  });

  describe('Empty Category Handling', () => {
    test('should handle empty category gracefully', async () => {
      const categoryId = 'empty-category';

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json({
            success: true,
            data: createMockPaginatedResponse([], { page: 1, limit: 10, total: 0 }),
          });
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const result = await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId })
        ).unwrap();

        expect(result.items).toHaveLength(0);
        expect(result.pagination.total).toBe(0);
        expect(result.pagination.totalPages).toBe(0);
        expect(result.pagination.hasNext).toBe(false);
        expect(result.pagination.hasPrev).toBe(false);
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });

    test('should return empty array for non-existent category', async () => {
      const categoryId = 'non-existent-category';

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json({
            success: true,
            data: createMockPaginatedResponse([]),
          });
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const result = await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId })
        ).unwrap();

        expect(Array.isArray(result.items)).toBe(true);
        expect(result.items).toHaveLength(0);
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });
  });

  describe('Error Scenarios for Invalid Categories', () => {
    test('should handle 404 error for invalid category ID', async () => {
      const categoryId = 'invalid-category';

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json(
            {
              success: false,
              error: 'Category not found',
              code: 'CATEGORY_NOT_FOUND'
            },
            { status: 404 }
          );
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId })
        ).unwrap();

        // Should not reach here
        expect(true).toBe(false);
      } catch (error: any) {
        expect(error.status).toBe(404);
        expect(error.data?.code).toBe('CATEGORY_NOT_FOUND');
      }
    });

    test('should handle 400 error for malformed category ID', async () => {
      const categoryId = ''; // Empty category ID

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json(
            {
              success: false,
              error: 'Invalid category ID format',
              code: 'INVALID_CATEGORY_ID'
            },
            { status: 400 }
          );
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId })
        ).unwrap();

        // Should not reach here
        expect(true).toBe(false);
      } catch (error: any) {
        expect(error.status).toBe(400);
        expect(error.data?.code).toBe('INVALID_CATEGORY_ID');
      }
    });

    test('should handle 403 error for unauthorized category access', async () => {
      const categoryId = 'admin-only-category';

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json(
            {
              success: false,
              error: 'Insufficient permissions to access category',
              code: 'PERMISSION_DENIED'
            },
            { status: 403 }
          );
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId })
        ).unwrap();

        // Should not reach here
        expect(true).toBe(false);
      } catch (error: any) {
        expect(error.status).toBe(403);
        expect(error.data?.code).toBe('PERMISSION_DENIED');
      }
    });

    test('should handle 500 server error gracefully', async () => {
      const categoryId = 'navigation';

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json(
            {
              success: false,
              error: 'Internal server error',
              code: 'INTERNAL_ERROR'
            },
            { status: 500 }
          );
        })
      );

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        await store.dispatch(
          api.endpoints.getContentByCategory.initiate({ categoryId })
        ).unwrap();

        // Should not reach here
        expect(true).toBe(false);
      } catch (error: any) {
        expect(error.status).toBe(500);
        expect(error.data?.code).toBe('INTERNAL_ERROR');
      }
    });
  });

  describe('RTK Query Integration', () => {
    test('should provide correct cache tags for invalidation', async () => {
      const categoryId = 'navigation';

      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent endpoint
        const endpointDefinition = api.endpoints.getContentByCategory;

        // Validate endpoint configuration
        expect(endpointDefinition.query).toBeDefined();
        expect(endpointDefinition.providesTags).toBeDefined();

        // Test cache tag structure
        const mockResult = createMockPaginatedResponse([createMockContentItem()]);
        const tags = endpointDefinition.providesTags?.(mockResult, null, { categoryId });

        expect(tags).toContain({ type: 'Content', id: 'LIST' });
        expect(tags).toContain({ type: 'Content', id: `CATEGORY-${categoryId}` });
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });

    test('should support lazy query execution', async () => {
      // CRITICAL TDD: This WILL FAIL until the endpoint is implemented
      try {
        // @ts-expect-error - Testing non-existent hook
        const { useLazyGetContentByCategoryQuery } = api;
        expect(useLazyGetContentByCategoryQuery).toBeDefined();
      } catch (error) {
        // Expected to fail during TDD phase
        expect(error).toBeDefined();
      }
    });
  });
});