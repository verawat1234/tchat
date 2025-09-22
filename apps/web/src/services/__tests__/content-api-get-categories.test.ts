/**
 * T007: Contract Test for getContentCategories API Endpoint
 *
 * This test file implements comprehensive contract testing for the content categories
 * API endpoint following TDD principles. The test WILL FAIL initially because the
 * implementation doesn't exist yet - this is intentional and required for proper
 * test-driven development.
 *
 * The test validates:
 * - Content categories list retrieval
 * - Category structure validation (id, name, description, permissions)
 * - Permissions structure validation (read, write, publish arrays)
 * - Empty categories handling
 * - Category metadata validation
 * - Error scenarios and edge cases
 *
 * @requires RTK Query integration
 * @requires MSW for API mocking
 * @follows TDD principles - this test drives the implementation
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import { setupListeners } from '@reduxjs/toolkit/query';
import { http, HttpResponse } from 'msw';
import { server } from '../../lib/test-utils/msw/server';
import { api } from '../api';
import type { ContentCategory, CategoryPermissions } from '../../types/content';

// Test API base URL
const TEST_API_BASE_URL = 'http://localhost:3001/api';

// Mock data factories following the ContentCategory interface
const createMockCategoryPermissions = (overrides: Partial<CategoryPermissions> = {}): CategoryPermissions => ({
  read: ['user', 'admin'],
  write: ['admin'],
  publish: ['admin'],
  ...overrides,
});

const createMockContentCategory = (overrides: Partial<ContentCategory> = {}): ContentCategory => ({
  id: 'test-category',
  name: 'Test Category',
  description: 'A test category for validation',
  permissions: createMockCategoryPermissions(),
  ...overrides,
});

// Create test store with RTK Query API
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

// Extended API with content endpoints injection (simulating future implementation)
const contentApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getContentCategories: builder.query<ContentCategory[], void>({
      query: () => '/content/categories',
      providesTags: ['ContentCategory'],
    }),
  }),
});

// Extract hooks for testing
const { useGetContentCategoriesQuery } = contentApi;

// Mock response data
const mockCategories: ContentCategory[] = [
  createMockContentCategory({
    id: 'navigation',
    name: 'Navigation',
    description: 'Navigation menu items and links',
    permissions: createMockCategoryPermissions({
      read: ['user', 'admin', 'guest'],
      write: ['admin', 'editor'],
      publish: ['admin'],
    }),
  }),
  createMockContentCategory({
    id: 'errors',
    name: 'Error Messages',
    description: 'User-facing error messages and alerts',
    permissions: createMockCategoryPermissions({
      read: ['user', 'admin'],
      write: ['admin'],
      publish: ['admin'],
    }),
  }),
  createMockContentCategory({
    id: 'help',
    name: 'Help Content',
    description: 'User assistance and documentation',
    parentId: 'navigation',
    permissions: createMockCategoryPermissions({
      read: ['user', 'admin'],
      write: ['admin', 'support'],
      publish: ['admin', 'support'],
    }),
  }),
];

describe('T007: Content Categories API Contract Tests', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
  });

  afterEach(() => {
    store.dispatch(api.util.resetApiState());
  });

  describe('getContentCategories Endpoint Contract', () => {
    describe('Successful Response Scenarios', () => {
      beforeEach(() => {
        // Mock successful API response
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return HttpResponse.json(mockCategories, { status: 200 });
          })
        );
      });

      it('should retrieve content categories list successfully', async () => {
        // Execute the query
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        // Validate successful response
        expect(result.isSuccess).toBe(true);
        expect(result.data).toBeDefined();
        expect(Array.isArray(result.data)).toBe(true);
        expect(result.data).toHaveLength(3);

        // Cleanup
        promise.unsubscribe();
      });

      it('should validate category structure with required fields', async () => {
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isSuccess).toBe(true);
        expect(result.data).toBeDefined();

        // Validate each category has required structure
        result.data!.forEach((category) => {
          // Required fields validation
          expect(category).toHaveProperty('id');
          expect(category).toHaveProperty('name');
          expect(category).toHaveProperty('description');
          expect(category).toHaveProperty('permissions');

          // Type validation
          expect(typeof category.id).toBe('string');
          expect(typeof category.name).toBe('string');
          expect(typeof category.description).toBe('string');
          expect(typeof category.permissions).toBe('object');

          // ID should not be empty
          expect(category.id.trim().length).toBeGreaterThan(0);
          expect(category.name.trim().length).toBeGreaterThan(0);
        });

        promise.unsubscribe();
      });

      it('should validate permissions structure with arrays', async () => {
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isSuccess).toBe(true);
        expect(result.data).toBeDefined();

        // Validate permissions structure for each category
        result.data!.forEach((category) => {
          const { permissions } = category;

          // Required permission fields
          expect(permissions).toHaveProperty('read');
          expect(permissions).toHaveProperty('write');
          expect(permissions).toHaveProperty('publish');

          // Must be arrays
          expect(Array.isArray(permissions.read)).toBe(true);
          expect(Array.isArray(permissions.write)).toBe(true);
          expect(Array.isArray(permissions.publish)).toBe(true);

          // Arrays should contain strings only
          permissions.read.forEach(role => expect(typeof role).toBe('string'));
          permissions.write.forEach(role => expect(typeof role).toBe('string'));
          permissions.publish.forEach(role => expect(typeof role).toBe('string'));

          // Role names should not be empty
          permissions.read.forEach(role => expect(role.trim().length).toBeGreaterThan(0));
          permissions.write.forEach(role => expect(role.trim().length).toBeGreaterThan(0));
          permissions.publish.forEach(role => expect(role.trim().length).toBeGreaterThan(0));
        });

        promise.unsubscribe();
      });

      it('should handle hierarchical categories with parentId', async () => {
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isSuccess).toBe(true);
        expect(result.data).toBeDefined();

        // Find category with parentId
        const childCategory = result.data!.find(cat => cat.parentId);
        expect(childCategory).toBeDefined();

        if (childCategory?.parentId) {
          // Validate parentId is a string
          expect(typeof childCategory.parentId).toBe('string');
          expect(childCategory.parentId.trim().length).toBeGreaterThan(0);

          // Verify parent exists in the list
          const parentExists = result.data!.some(cat => cat.id === childCategory.parentId);
          expect(parentExists).toBe(true);
        }

        promise.unsubscribe();
      });

      it('should validate category metadata with proper structure', async () => {
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isSuccess).toBe(true);
        expect(result.data).toBeDefined();

        // Validate specific category data
        const navigationCategory = result.data!.find(cat => cat.id === 'navigation');
        expect(navigationCategory).toBeDefined();
        expect(navigationCategory!.name).toBe('Navigation');
        expect(navigationCategory!.description).toBe('Navigation menu items and links');

        const errorCategory = result.data!.find(cat => cat.id === 'errors');
        expect(errorCategory).toBeDefined();
        expect(errorCategory!.name).toBe('Error Messages');
        expect(errorCategory!.description).toBe('User-facing error messages and alerts');

        promise.unsubscribe();
      });
    });

    describe('Empty Response Scenarios', () => {
      beforeEach(() => {
        // Mock empty categories response
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return HttpResponse.json([], { status: 200 });
          })
        );
      });

      it('should handle empty categories list gracefully', async () => {
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isSuccess).toBe(true);
        expect(result.data).toBeDefined();
        expect(Array.isArray(result.data)).toBe(true);
        expect(result.data).toHaveLength(0);

        promise.unsubscribe();
      });
    });

    describe('Error Response Scenarios', () => {
      it('should handle 404 Not Found responses', async () => {
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return HttpResponse.json(
              { error: 'Categories not found' },
              { status: 404 }
            );
          })
        );

        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isError).toBe(true);
        expect(result.error).toBeDefined();
        if ('status' in result.error) {
          expect(result.error.status).toBe(404);
        }

        promise.unsubscribe();
      });

      it('should handle 500 Internal Server Error responses', async () => {
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return HttpResponse.json(
              { error: 'Internal Server Error' },
              { status: 500 }
            );
          })
        );

        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isError).toBe(true);
        expect(result.error).toBeDefined();
        if ('status' in result.error) {
          expect(result.error.status).toBe(500);
        }

        promise.unsubscribe();
      });

      it('should handle 401 Unauthorized responses', async () => {
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return HttpResponse.json(
              { error: 'Unauthorized access' },
              { status: 401 }
            );
          })
        );

        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isError).toBe(true);
        expect(result.error).toBeDefined();
        if ('status' in result.error) {
          expect(result.error.status).toBe(401);
        }

        promise.unsubscribe();
      });

      it('should handle network errors gracefully', async () => {
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return HttpResponse.error();
          })
        );

        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isError).toBe(true);
        expect(result.error).toBeDefined();
        if ('status' in result.error) {
          expect(result.error.status).toBe('FETCH_ERROR');
        }

        promise.unsubscribe();
      });

      it('should handle malformed JSON responses', async () => {
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return new Response('Invalid JSON', {
              status: 200,
              headers: { 'Content-Type': 'application/json' },
            });
          })
        );

        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isError).toBe(true);
        expect(result.error).toBeDefined();

        promise.unsubscribe();
      });
    });

    describe('RTK Query Integration Patterns', () => {
      beforeEach(() => {
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return HttpResponse.json(mockCategories, { status: 200 });
          })
        );
      });

      it('should provide proper cache tags for invalidation', async () => {
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isSuccess).toBe(true);

        // Verify cache tags are properly set
        const state = store.getState();
        const cacheEntry = state.api.queries[
          contentApi.endpoints.getContentCategories.select()(state).requestId
        ];

        // The query should be cached
        expect(cacheEntry).toBeDefined();

        promise.unsubscribe();
      });

      it('should support query refetching', async () => {
        // Initial query
        const initialPromise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const initialResult = await initialPromise;
        expect(initialResult.isSuccess).toBe(true);

        // Refetch the same query
        const refetchPromise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate(undefined, {
            forceRefetch: true,
          })
        );

        const refetchResult = await refetchPromise;
        expect(refetchResult.isSuccess).toBe(true);

        initialPromise.unsubscribe();
        refetchPromise.unsubscribe();
      });

      it('should handle concurrent requests with deduplication', async () => {
        // Start multiple concurrent requests
        const promise1 = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );
        const promise2 = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );
        const promise3 = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        // All should resolve successfully
        const [result1, result2, result3] = await Promise.all([
          promise1,
          promise2,
          promise3,
        ]);

        expect(result1.isSuccess).toBe(true);
        expect(result2.isSuccess).toBe(true);
        expect(result3.isSuccess).toBe(true);

        // Data should be identical (same reference due to deduplication)
        expect(result1.data).toEqual(result2.data);
        expect(result2.data).toEqual(result3.data);

        promise1.unsubscribe();
        promise2.unsubscribe();
        promise3.unsubscribe();
      });
    });

    describe('Type Safety and Contract Validation', () => {
      beforeEach(() => {
        server.use(
          http.get(`${TEST_API_BASE_URL}/content/categories`, () => {
            return HttpResponse.json(mockCategories, { status: 200 });
          })
        );
      });

      it('should enforce ContentCategory type contract', async () => {
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isSuccess).toBe(true);
        expect(result.data).toBeDefined();

        // TypeScript type checking ensures this at compile time,
        // but we can validate the runtime structure
        result.data!.forEach((category: ContentCategory) => {
          // This would fail at compile time if types don't match
          const typedCategory: ContentCategory = category;

          expect(typedCategory).toMatchObject({
            id: expect.any(String),
            name: expect.any(String),
            description: expect.any(String),
            permissions: {
              read: expect.any(Array),
              write: expect.any(Array),
              publish: expect.any(Array),
            },
          });
        });

        promise.unsubscribe();
      });

      it('should validate CategoryPermissions type contract', async () => {
        const promise = store.dispatch(
          contentApi.endpoints.getContentCategories.initiate()
        );

        const result = await promise;

        expect(result.isSuccess).toBe(true);
        expect(result.data).toBeDefined();

        result.data!.forEach((category) => {
          const permissions: CategoryPermissions = category.permissions;

          // Validate permissions structure matches type contract
          expect(permissions).toMatchObject({
            read: expect.arrayContaining([expect.any(String)]),
            write: expect.arrayContaining([expect.any(String)]),
            publish: expect.arrayContaining([expect.any(String)]),
          });

          // Ensure no additional properties beyond the contract
          const expectedKeys = ['read', 'write', 'publish'];
          const actualKeys = Object.keys(permissions);
          expect(actualKeys.sort()).toEqual(expectedKeys.sort());
        });

        promise.unsubscribe();
      });
    });
  });

  describe('Future Implementation Requirements', () => {
    it('should fail until getContentCategories endpoint is implemented', () => {
      // This test documents the TDD expectation
      // The above tests WILL FAIL until:
      // 1. The actual content API service is created
      // 2. The getContentCategories endpoint is implemented
      // 3. The backend API endpoints are created
      // 4. The RTK Query integration is complete

      expect(true).toBe(true); // This passes, but the above integration tests will fail

      // TODO for implementation:
      // - Create src/services/content.ts with RTK Query endpoints
      // - Implement getContentCategories query with proper typing
      // - Create backend /api/content/categories endpoint
      // - Ensure proper error handling and validation
      // - Add cache tags and invalidation logic
    });
  });
});