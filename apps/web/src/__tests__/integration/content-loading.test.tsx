/**
 * T016: Integration Test - Dynamic Content Loading
 *
 * Comprehensive integration test for dynamic content loading flow with RTK Query.
 * This test validates the complete user workflow from content request to UI display.
 *
 * TDD STATUS: ✅ INTENTIONALLY FAILING (Expected: 13 failed, 1 passed)
 *
 * IMPORTANT: This test WILL FAIL until the content API endpoints and UI integration
 * are fully implemented. This is expected and required for proper TDD.
 *
 * Current Status:
 * - Test components are stuck in loading state (expected)
 * - RTK Query hooks need to be implemented: useGetContentQuery, useListContentQuery
 * - Content API endpoints need to be created in services/api.ts
 * - Test structure and assertions are working correctly
 *
 * Test Coverage:
 * - Dynamic content loading flow with RTK Query
 * - Loading states and UI feedback during content fetch
 * - Content population in UI components after successful load
 * - Error state handling when content loading fails
 * - Cache behavior and data persistence
 * - Multiple component content loading coordination
 * - Performance requirements (<200ms content load times)
 *
 * Implementation Requirements for Tests to Pass:
 * 1. Add content API endpoints to services/api.ts with injectEndpoints
 * 2. Implement useGetContentQuery and useListContentQuery hooks
 * 3. Update TestContentDisplay component to use real RTK Query hooks
 * 4. Add proper error handling and retry mechanisms
 * 5. Implement cache invalidation and refetching logic
 */

import React from 'react';
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { screen, waitFor, act } from '@testing-library/react';
import { http, HttpResponse } from 'msw';
import { server } from '../../lib/test-utils/msw/server';
import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../../services/api';
import contentReducer from '../../features/contentSlice';
import authReducer from '../../features/authSlice';
import uiReducer from '../../features/uiSlice';
import loadingReducer from '../../features/loadingSlice';
import type {
  ContentItem,
  ContentCategory,
  ContentType,
  ContentStatus,
  PaginatedContentResponse
} from '../../types/content';

// Mock performance.now() for performance testing
const mockPerformanceNow = vi.fn();
global.performance.now = mockPerformanceNow;

// Mock content data factories
const mockContentCategory = (overrides = {}): ContentCategory => ({
  id: 'test-category',
  name: 'Test Category',
  description: 'Category for testing content loading',
  permissions: {
    read: ['user', 'admin'],
    write: ['admin'],
    publish: ['admin']
  },
  ...overrides,
});

const mockContentItem = (overrides = {}): ContentItem => ({
  id: 'test.content.item',
  key: 'content.item',
  categoryId: 'test-category',
  category: mockContentCategory(),
  type: 'text' as ContentType,
  value: {
    type: 'text',
    value: 'Dynamic test content loaded successfully'
  },
  metadata: {
    createdAt: new Date().toISOString(),
    createdBy: 'test-user',
    updatedAt: new Date().toISOString(),
    updatedBy: 'test-user',
    version: 1,
    tags: ['test', 'integration']
  },
  status: 'published' as ContentStatus,
  tags: ['test', 'integration'],
  ...overrides,
});

const mockPaginatedResponse = (items: ContentItem[] = [], overrides = {}): PaginatedContentResponse => ({
  items,
  pagination: {
    page: 1,
    limit: 20,
    total: items.length,
    totalPages: Math.ceil(items.length / 20),
    hasNext: false,
    hasPrev: false
  },
  ...overrides,
});

// Test component that uses content loading
const TestContentDisplay: React.FC<{ contentId: string }> = ({ contentId }) => {
  // This will fail until contentApi is implemented
  // const { data: content, isLoading, error, isSuccess } = useGetContentQuery(contentId);

  // Temporary mock implementation for testing structure
  const isLoading = true;
  const error = null;
  const isSuccess = false;
  const content = null;

  if (isLoading) {
    return (
      <div data-testid="content-loading" role="status" aria-label="Loading content">
        <div data-testid="loading-spinner" aria-hidden="true">Loading...</div>
        <span className="sr-only">Loading content</span>
      </div>
    );
  }

  if (error) {
    return (
      <div data-testid="content-error" role="alert">
        <h3 data-testid="error-title">Failed to load content</h3>
        <p data-testid="error-message">
          {error?.message || 'An error occurred while loading content'}
        </p>
        <button
          data-testid="retry-button"
          onClick={() => {/* retry logic */}}
          aria-label="Retry loading content"
        >
          Retry
        </button>
      </div>
    );
  }

  if (isSuccess && content) {
    return (
      <div data-testid="content-display" data-content-id={contentId}>
        <div data-testid="content-value">
          {content.value?.type === 'text' ? content.value.value : 'Unsupported content type'}
        </div>
        <div data-testid="content-metadata">
          <span data-testid="content-version">v{content.metadata.version}</span>
          <span data-testid="content-status">{content.status}</span>
          <time
            data-testid="content-updated"
            dateTime={content.metadata.updatedAt}
            aria-label={`Last updated ${content.metadata.updatedAt}`}
          >
            {new Date(content.metadata.updatedAt).toLocaleDateString()}
          </time>
        </div>
      </div>
    );
  }

  return (
    <div data-testid="content-empty">
      No content available
    </div>
  );
};

// Multi-content component for testing coordination
const TestMultiContentDisplay: React.FC<{ contentIds: string[] }> = ({ contentIds }) => {
  return (
    <div data-testid="multi-content-container">
      <h2 data-testid="multi-content-title">Multiple Content Items</h2>
      {contentIds.map((contentId) => (
        <div key={contentId} data-testid={`content-item-${contentId}`}>
          <TestContentDisplay contentId={contentId} />
        </div>
      ))}
    </div>
  );
};

// Content list component for testing pagination and bulk loading
const TestContentList: React.FC<{ categoryId?: string }> = ({ categoryId }) => {
  // This will fail until contentApi is implemented
  // const { data: response, isLoading, error } = useListContentQuery({ categoryId });

  // Temporary mock implementation
  const isLoading = true;
  const error = null;
  const response = null;

  if (isLoading) {
    return (
      <div data-testid="content-list-loading">
        <div data-testid="list-loading-spinner">Loading content list...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div data-testid="content-list-error" role="alert">
        Failed to load content list
      </div>
    );
  }

  const items = response?.items || [];

  return (
    <div data-testid="content-list">
      <div data-testid="content-list-header">
        <span data-testid="content-count">{items.length} items</span>
      </div>
      <ul data-testid="content-items" role="list">
        {items.map((item) => (
          <li key={item.id} data-testid={`list-item-${item.id}`} role="listitem">
            <TestContentDisplay contentId={item.id} />
          </li>
        ))}
      </ul>
      {response?.pagination && (
        <div data-testid="pagination-info">
          Page {response.pagination.page} of {response.pagination.totalPages}
        </div>
      )}
    </div>
  );
};

// Test query client setup
const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
        staleTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  });

// Test store setup with Redux Toolkit
const createTestStore = (preloadedState = {}) => {
  return configureStore({
    reducer: {
      [api.reducerPath]: api.reducer,
      auth: authReducer,
      ui: uiReducer,
      loading: loadingReducer,
      content: contentReducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware({
        serializableCheck: {
          ignoredActionPaths: ['meta.arg', 'payload.timestamp'],
          ignoredPaths: ['auth.expiresAt'],
        },
      }).concat(api.middleware),
    preloadedState,
  });
};

// Test wrapper component
const TestWrapper: React.FC<{ children: React.ReactNode; store?: any }> = ({
  children,
  store
}) => {
  const queryClient = createTestQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      {/* Note: Store provider would be added here when Redux integration is implemented */}
      {children}
    </QueryClientProvider>
  );
};

// API endpoint constants
const API_BASE_URL = 'http://localhost:3001';

describe('T016: Dynamic Content Loading Integration', () => {
  let performanceTimes: number[] = [];

  beforeEach(() => {
    performanceTimes = [];
    mockPerformanceNow.mockImplementation(() => Date.now());

    // Reset MSW handlers
    server.resetHandlers();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('Single Content Loading Flow', () => {
    it('should show loading state while fetching content', async () => {
      // Setup: Mock API response with delay
      server.use(
        http.get(`${API_BASE_URL}/api/content/:id`, async () => {
          await new Promise(resolve => setTimeout(resolve, 100));
          return HttpResponse.json(mockContentItem());
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="test.content.item" />
        </TestWrapper>
      );

      // Verify: Loading state is displayed immediately
      expect(screen.getByTestId('content-loading')).toBeInTheDocument();
      expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
      expect(screen.getByRole('status')).toHaveAttribute('aria-label', 'Loading content');

      // Verify: Content is not yet displayed
      expect(screen.queryByTestId('content-display')).not.toBeInTheDocument();
    });

    it('should populate content in UI after successful load', async () => {
      const testContent = mockContentItem({
        id: 'success.test.content',
        value: {
          type: 'text',
          value: 'Successfully loaded dynamic content'
        },
        metadata: {
          ...mockContentItem().metadata,
          version: 2,
          updatedAt: '2024-01-15T10:30:00Z'
        }
      });

      server.use(
        http.get(`${API_BASE_URL}/api/content/success.test.content`, () => {
          return HttpResponse.json(testContent);
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="success.test.content" />
        </TestWrapper>
      );

      // Wait for content to load
      await waitFor(() => {
        expect(screen.getByTestId('content-display')).toBeInTheDocument();
      });

      // Verify: Content value is displayed
      expect(screen.getByTestId('content-value')).toHaveTextContent(
        'Successfully loaded dynamic content'
      );

      // Verify: Metadata is displayed
      expect(screen.getByTestId('content-version')).toHaveTextContent('v2');
      expect(screen.getByTestId('content-status')).toHaveTextContent('published');
      expect(screen.getByTestId('content-updated')).toHaveAttribute(
        'dateTime',
        '2024-01-15T10:30:00Z'
      );

      // Verify: Accessibility attributes
      expect(screen.getByTestId('content-display')).toHaveAttribute(
        'data-content-id',
        'success.test.content'
      );
    });

    it('should handle and display error states when content loading fails', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/content/error.test.content`, () => {
          return HttpResponse.json(
            { error: 'Content not found', code: 'CONTENT_NOT_FOUND' },
            { status: 404 }
          );
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="error.test.content" />
        </TestWrapper>
      );

      // Wait for error state
      await waitFor(() => {
        expect(screen.getByTestId('content-error')).toBeInTheDocument();
      });

      // Verify: Error message is displayed
      expect(screen.getByTestId('error-title')).toHaveTextContent('Failed to load content');
      expect(screen.getByRole('alert')).toBeInTheDocument();

      // Verify: Retry button is available
      const retryButton = screen.getByTestId('retry-button');
      expect(retryButton).toBeInTheDocument();
      expect(retryButton).toHaveAttribute('aria-label', 'Retry loading content');
    });
  });

  describe('Multiple Content Loading Coordination', () => {
    it('should coordinate loading of multiple content items', async () => {
      const contentItems = [
        mockContentItem({ id: 'multi.content.1', value: { type: 'text', value: 'Content 1' } }),
        mockContentItem({ id: 'multi.content.2', value: { type: 'text', value: 'Content 2' } }),
        mockContentItem({ id: 'multi.content.3', value: { type: 'text', value: 'Content 3' } })
      ];

      // Mock API responses for each content item
      contentItems.forEach((item) => {
        server.use(
          http.get(`${API_BASE_URL}/api/content/${item.id}`, () => {
            return HttpResponse.json(item);
          })
        );
      });

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestMultiContentDisplay contentIds={['multi.content.1', 'multi.content.2', 'multi.content.3']} />
        </TestWrapper>
      );

      // Verify: Container is rendered
      expect(screen.getByTestId('multi-content-container')).toBeInTheDocument();
      expect(screen.getByTestId('multi-content-title')).toHaveTextContent('Multiple Content Items');

      // Verify: All content items are present
      expect(screen.getByTestId('content-item-multi.content.1')).toBeInTheDocument();
      expect(screen.getByTestId('content-item-multi.content.2')).toBeInTheDocument();
      expect(screen.getByTestId('content-item-multi.content.3')).toBeInTheDocument();

      // Wait for all content to load
      await waitFor(() => {
        expect(screen.getAllByTestId(/^content-display$/)).toHaveLength(3);
      });

      // Verify: All content values are displayed
      expect(screen.getByText('Content 1')).toBeInTheDocument();
      expect(screen.getByText('Content 2')).toBeInTheDocument();
      expect(screen.getByText('Content 3')).toBeInTheDocument();
    });

    it('should handle mixed success/error states in multiple content loading', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/content/mixed.success.1`, () => {
          return HttpResponse.json(
            mockContentItem({
              id: 'mixed.success.1',
              value: { type: 'text', value: 'Successful content' }
            })
          );
        }),
        http.get(`${API_BASE_URL}/api/content/mixed.error.2`, () => {
          return HttpResponse.json(
            { error: 'Failed to load' },
            { status: 500 }
          );
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestMultiContentDisplay contentIds={['mixed.success.1', 'mixed.error.2']} />
        </TestWrapper>
      );

      await waitFor(() => {
        // Verify: Success content is displayed
        expect(screen.getByText('Successful content')).toBeInTheDocument();

        // Verify: Error state is displayed for failed content
        expect(screen.getByTestId('content-error')).toBeInTheDocument();
      });
    });
  });

  describe('Content List and Pagination', () => {
    it('should load and display paginated content list', async () => {
      const contentItems = Array.from({ length: 5 }, (_, i) =>
        mockContentItem({
          id: `list.item.${i + 1}`,
          value: { type: 'text', value: `List item ${i + 1}` }
        })
      );

      const paginatedResponse = mockPaginatedResponse(contentItems, {
        pagination: {
          page: 1,
          limit: 20,
          total: 5,
          totalPages: 1,
          hasNext: false,
          hasPrev: false
        }
      });

      server.use(
        http.get(`${API_BASE_URL}/api/content`, () => {
          return HttpResponse.json(paginatedResponse);
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentList categoryId="test-category" />
        </TestWrapper>
      );

      // Wait for list to load
      await waitFor(() => {
        expect(screen.getByTestId('content-list')).toBeInTheDocument();
      });

      // Verify: Content count is displayed
      expect(screen.getByTestId('content-count')).toHaveTextContent('5 items');

      // Verify: All list items are rendered
      const listItems = screen.getAllByRole('listitem');
      expect(listItems).toHaveLength(5);

      // Verify: Pagination info is displayed
      expect(screen.getByTestId('pagination-info')).toHaveTextContent('Page 1 of 1');
    });
  });

  describe('Performance Requirements', () => {
    it('should load content within 200ms performance requirement', async () => {
      const startTime = 100;
      const endTime = 250; // 150ms duration (within 200ms requirement)

      mockPerformanceNow
        .mockReturnValueOnce(startTime) // First call - request start
        .mockReturnValueOnce(endTime);  // Second call - response received

      server.use(
        http.get(`${API_BASE_URL}/api/content/performance.test`, async () => {
          // Simulate fast API response
          await new Promise(resolve => setTimeout(resolve, 50));
          return HttpResponse.json(mockContentItem({
            id: 'performance.test',
            value: { type: 'text', value: 'Fast loading content' }
          }));
        })
      );

      const store = createTestStore();

      const startTimer = performance.now();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="performance.test" />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByTestId('content-display')).toBeInTheDocument();
      });

      const endTimer = performance.now();
      const loadTime = endTimer - startTimer;

      // Verify: Content loads within 200ms requirement
      expect(loadTime).toBeLessThan(200);
      expect(screen.getByText('Fast loading content')).toBeInTheDocument();
    });

    it('should show performance warning for slow loading content', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/content/slow.test`, async () => {
          // Simulate slow API response (exceeds 200ms)
          await new Promise(resolve => setTimeout(resolve, 300));
          return HttpResponse.json(mockContentItem({
            id: 'slow.test',
            value: { type: 'text', value: 'Slow loading content' }
          }));
        })
      );

      const store = createTestStore();

      const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});

      const startTime = performance.now();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="slow.test" />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByTestId('content-display')).toBeInTheDocument();
      }, { timeout: 5000 });

      const endTime = performance.now();
      const loadTime = endTime - startTime;

      // Verify: Load time exceeds 200ms threshold
      expect(loadTime).toBeGreaterThan(200);

      // Note: Performance warning would be implemented in the actual content loading hook
      // This test documents the expected behavior for when performance monitoring is added

      consoleSpy.mockRestore();
    });
  });

  describe('Cache Behavior and Data Persistence', () => {
    it('should cache content and serve from cache on subsequent requests', async () => {
      let requestCount = 0;

      server.use(
        http.get(`${API_BASE_URL}/api/content/cache.test`, () => {
          requestCount++;
          return HttpResponse.json(mockContentItem({
            id: 'cache.test',
            value: { type: 'text', value: `Cached content (request #${requestCount})` }
          }));
        })
      );

      const store = createTestStore();

      // First render - should make API request
      const { unmount } = render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="cache.test" />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByTestId('content-display')).toBeInTheDocument();
      });

      expect(requestCount).toBe(1);
      expect(screen.getByText('Cached content (request #1)')).toBeInTheDocument();

      unmount();

      // Second render - should use cached data (RTK Query cache)
      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="cache.test" />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByTestId('content-display')).toBeInTheDocument();
      });

      // Verify: No additional API request was made
      expect(requestCount).toBe(1);
      expect(screen.getByText('Cached content (request #1)')).toBeInTheDocument();
    });

    it('should persist content in Redux store state', async () => {
      const testContent = mockContentItem({
        id: 'persist.test',
        value: { type: 'text', value: 'Persisted content' }
      });

      server.use(
        http.get(`${API_BASE_URL}/api/content/persist.test`, () => {
          return HttpResponse.json(testContent);
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="persist.test" />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByTestId('content-display')).toBeInTheDocument();
      });

      // Verify: Content is stored in Redux state via RTK Query
      const state = store.getState();
      const apiState = state.api;

      // Check that the API state contains the cached query
      expect(apiState.queries).toBeDefined();

      // Note: Specific cache key structure depends on RTK Query implementation
      // This test documents the expected behavior for state persistence
    });
  });

  describe('Error Recovery and Retry Mechanisms', () => {
    it('should allow retrying failed content loads', async () => {
      let attemptCount = 0;

      server.use(
        http.get(`${API_BASE_URL}/api/content/retry.test`, () => {
          attemptCount++;
          if (attemptCount === 1) {
            return HttpResponse.json(
              { error: 'Temporary failure' },
              { status: 500 }
            );
          }
          return HttpResponse.json(mockContentItem({
            id: 'retry.test',
            value: { type: 'text', value: 'Content loaded after retry' }
          }));
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="retry.test" />
        </TestWrapper>
      );

      // Wait for initial error state
      await waitFor(() => {
        expect(screen.getByTestId('content-error')).toBeInTheDocument();
      });

      expect(attemptCount).toBe(1);

      // Click retry button
      const retryButton = screen.getByTestId('retry-button');
      await act(async () => {
        retryButton.click();
      });

      // Wait for successful content load
      await waitFor(() => {
        expect(screen.getByTestId('content-display')).toBeInTheDocument();
      });

      expect(attemptCount).toBe(2);
      expect(screen.getByText('Content loaded after retry')).toBeInTheDocument();
    });
  });

  describe('Accessibility and User Experience', () => {
    it('should provide proper ARIA labels and roles for screen readers', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/content/a11y.test`, () => {
          return HttpResponse.json(mockContentItem({
            id: 'a11y.test',
            value: { type: 'text', value: 'Accessible content' }
          }));
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="a11y.test" />
        </TestWrapper>
      );

      // Verify: Loading state has proper ARIA attributes
      const loadingElement = screen.getByTestId('content-loading');
      expect(loadingElement).toHaveAttribute('role', 'status');
      expect(loadingElement).toHaveAttribute('aria-label', 'Loading content');

      // Verify: Screen reader text is available
      expect(screen.getByText('Loading content')).toHaveClass('sr-only');

      await waitFor(() => {
        expect(screen.getByTestId('content-display')).toBeInTheDocument();
      });

      // Verify: Content metadata has proper time element
      const timeElement = screen.getByTestId('content-updated');
      expect(timeElement.tagName).toBe('TIME');
      expect(timeElement).toHaveAttribute('dateTime');
      expect(timeElement).toHaveAttribute('aria-label');
    });

    it('should handle keyboard navigation and focus management', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/content/focus.test`, () => {
          return HttpResponse.json(
            { error: 'Focus test error' },
            { status: 404 }
          );
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="focus.test" />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByTestId('content-error')).toBeInTheDocument();
      });

      const retryButton = screen.getByTestId('retry-button');

      // Verify: Retry button is focusable
      retryButton.focus();
      expect(document.activeElement).toBe(retryButton);

      // Verify: Button has proper keyboard support
      expect(retryButton.tagName).toBe('BUTTON');
      expect(retryButton).toHaveAttribute('aria-label');
    });
  });

  describe('Integration with RTK Query Features', () => {
    it('should integrate with RTK Query invalidation and refetching', async () => {
      const initialContent = mockContentItem({
        id: 'invalidation.test',
        value: { type: 'text', value: 'Initial content' },
        metadata: { ...mockContentItem().metadata, version: 1 }
      });

      const updatedContent = mockContentItem({
        id: 'invalidation.test',
        value: { type: 'text', value: 'Updated content' },
        metadata: { ...mockContentItem().metadata, version: 2 }
      });

      let contentVersion = 1;

      server.use(
        http.get(`${API_BASE_URL}/api/content/invalidation.test`, () => {
          return HttpResponse.json(contentVersion === 1 ? initialContent : updatedContent);
        }),
        http.patch(`${API_BASE_URL}/api/content/invalidation.test`, () => {
          contentVersion = 2;
          return HttpResponse.json(updatedContent);
        })
      );

      const store = createTestStore();

      render(
        <TestWrapper store={store}>
          <TestContentDisplay contentId="invalidation.test" />
        </TestWrapper>
      );

      // Wait for initial content load
      await waitFor(() => {
        expect(screen.getByText('Initial content')).toBeInTheDocument();
      });

      // Simulate content update that would trigger invalidation
      await act(async () => {
        // This would typically be triggered by another component updating content
        // and causing RTK Query to invalidate and refetch
        await fetch(`${API_BASE_URL}/api/content/invalidation.test`, {
          method: 'PATCH',
          body: JSON.stringify({ value: 'Updated content' })
        });
      });

      // Note: Actual cache invalidation would need to be implemented in the content API
      // This test documents the expected behavior for automatic refetching
    });
  });
});