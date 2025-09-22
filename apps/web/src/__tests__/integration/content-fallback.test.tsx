/**
 * T018: Content Fallback System Integration Test
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL until implementation is complete.
 * These tests are designed to fail initially as they test functionality that doesn't exist yet.
 *
 * EXPECTED FAILURES:
 * - All tests timeout or fail assertions because mock components don't implement:
 *   1. Content fetching hooks/logic
 *   2. Fallback mode activation
 *   3. localStorage fallback mechanism
 *   4. Error boundary integration
 *   5. User notification system
 *   6. Retry/recovery functionality
 *
 * IMPLEMENTATION REQUIREMENTS (components that need to be built):
 *
 * 1. MockContentDisplay component should:
 *    - Use a content hook (e.g., useContent) to fetch content by key
 *    - Detect API failures and activate fallback mode
 *    - Load content from localStorage when API fails
 *    - Update UI state (content-value, content-status, fallback-indicator)
 *    - Handle loading, error, and success states properly
 *
 * 2. MockContentManager component should:
 *    - Connect to Redux store content slice
 *    - Display real sync-status and fallback-status from store
 *    - Implement retry-sync functionality (trigger content refetch)
 *    - Implement clear-cache functionality (clear localStorage)
 *    - Dispatch appropriate Redux actions
 *
 * 3. MockNotificationDisplay component should:
 *    - React to fallback mode changes in Redux store
 *    - Show/hide notification-message based on state
 *    - Display appropriate messages for different error states
 *
 * 4. Content Management System Integration:
 *    - API endpoints for content fetching (MSW handlers already configured)
 *    - Content hooks for React components
 *    - Redux middleware for localStorage persistence
 *    - Error boundary components for graceful failures
 *    - Auto-retry mechanisms and recovery logic
 *
 * 5. localStorage Fallback System:
 *    - Automatic content caching on successful API calls
 *    - Fallback content loading when API fails
 *    - Data corruption handling (malformed JSON)
 *    - Cache invalidation and clearing
 *
 * TEST SCENARIOS COVERED:
 * - API failure with localStorage fallback activation
 * - Graceful degradation from dynamic to static content
 * - localStorage mechanism with corruption handling
 * - Error boundary integration for component failures
 * - User notification system for fallback mode
 * - Recovery and retry functionality
 * - Content consistency during mixed online/offline states
 * - Concurrent updates and race condition handling
 *
 * NEXT STEPS FOR IMPLEMENTATION:
 * 1. Create content fetching hooks (useContent, useContentManager)
 * 2. Implement Redux middleware for localStorage persistence
 * 3. Build error boundaries for content components
 * 4. Create notification system components
 * 5. Implement retry and recovery mechanisms
 * 6. Add real API integration (replace mock handlers)
 *
 * Once implementation is complete, these tests should pass and provide
 * comprehensive coverage of the entire content fallback system.
 */

import React from 'react';
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { screen, within, waitFor, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { server } from '../../lib/test-utils/msw/server';
import { http, HttpResponse } from 'msw';
import { render } from '@testing-library/react';
import { store as defaultStore } from '../../store';
import contentReducer, {
  toggleFallbackMode,
  updateFallbackContent,
  setSyncStatus,
  selectFallbackMode,
  selectFallbackContent,
  selectSyncStatus,
} from '../../features/contentSlice';
import type { ContentValue } from '../../types/content';

// Mock content component that will use the fallback system
const MockContentDisplay = ({ contentKey }: { contentKey: string }) => {
  // This component doesn't exist yet - will be implemented later
  // It should use hooks to fetch content and handle fallbacks
  return (
    <div data-testid={`content-display-${contentKey}`}>
      <div data-testid="content-value">Loading...</div>
      <div data-testid="content-status">idle</div>
      <div data-testid="fallback-indicator" style={{ display: 'none' }}>
        Using fallback content
      </div>
    </div>
  );
};

// Mock content management component for testing recovery
const MockContentManager = () => {
  // This component doesn't exist yet - will be implemented later
  return (
    <div data-testid="content-manager">
      <button data-testid="retry-sync">Retry Sync</button>
      <button data-testid="clear-cache">Clear Cache</button>
      <div data-testid="sync-status">idle</div>
      <div data-testid="fallback-status">false</div>
    </div>
  );
};

// Mock notification component for user feedback
const MockNotificationDisplay = () => {
  // This component doesn't exist yet - will be implemented later
  return (
    <div data-testid="notification-display">
      <div data-testid="notification-message" style={{ display: 'none' }}>
        Using offline content. Some data may be outdated.
      </div>
    </div>
  );
};

// Integration test app component
const TestApp = () => (
  <div data-testid="test-app">
    <MockNotificationDisplay />
    <MockContentManager />
    <MockContentDisplay contentKey="navigation.header.title" />
    <MockContentDisplay contentKey="error.network.timeout" />
    <MockContentDisplay contentKey="button.submit.text" />
  </div>
);

// Test store factory with content slice
const createTestStore = (initialState = {}) => {
  return configureStore({
    reducer: {
      content: contentReducer,
    },
    preloadedState: {
      content: {
        selectedLanguage: 'en',
        contentPreferences: {
          showDrafts: false,
          compactView: false,
        },
        lastSyncTime: new Date().toISOString(),
        syncStatus: 'idle' as const,
        fallbackMode: false,
        fallbackContent: {},
        ...initialState,
      },
    },
  });
};

// Test content data
const mockContentItems = {
  'navigation.header.title': {
    type: 'text' as const,
    value: 'Welcome to Tchat',
  },
  'error.network.timeout': {
    type: 'text' as const,
    value: 'Network request timed out. Please try again.',
  },
  'button.submit.text': {
    type: 'text' as const,
    value: 'Submit',
  },
};

const fallbackContentItems = {
  'navigation.header.title': {
    type: 'text' as const,
    value: 'Welcome (Offline)',
  },
  'error.network.timeout': {
    type: 'text' as const,
    value: 'Connection error (Cached)',
  },
  'button.submit.text': {
    type: 'text' as const,
    value: 'Submit (Offline)',
  },
};

// localStorage mock
const localStorageMock = (() => {
  let store: Record<string, string> = {};

  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value;
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key];
    }),
    clear: vi.fn(() => {
      store = {};
    }),
    length: 0,
    key: vi.fn(),
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

describe('T018: Content Fallback System Integration', () => {
  let testStore: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    testStore = createTestStore();
    localStorageMock.clear();
    vi.clearAllMocks();
    server.resetHandlers();
  });

  afterEach(() => {
    localStorageMock.clear();
  });

  describe('API Failure and Fallback Activation', () => {
    it('should activate fallback mode when API calls fail', async () => {
      // CRITICAL: This test MUST FAIL initially - no implementation exists yet

      // Setup: Mock API failures
      server.use(
        http.get('*/api/content/*', () => {
          return HttpResponse.error();
        })
      );

      // Pre-populate localStorage with fallback content
      const fallbackData = JSON.stringify(fallbackContentItems);
      localStorageMock.setItem('tchat-content-fallback', fallbackData);
      const user = userEvent.setup();

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for initial load and API failure detection
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('error');
      }, { timeout: 5000 });

      // Verify fallback mode is activated
      await waitFor(() => {
        expect(screen.getByTestId('fallback-status')).toHaveTextContent('true');
      });

      // Verify fallback content is displayed
      const titleDisplay = screen.getByTestId('content-display-navigation.header.title');
      expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome (Offline)');

      const errorDisplay = screen.getByTestId('content-display-error.network.timeout');
      expect(within(errorDisplay).getByTestId('content-value')).toHaveTextContent('Connection error (Cached)');

      // Verify fallback indicators are visible
      expect(within(titleDisplay).getByTestId('fallback-indicator')).toBeVisible();
      expect(within(errorDisplay).getByTestId('fallback-indicator')).toBeVisible();

      // Verify user notification about fallback mode
      const notification = screen.getByTestId('notification-message');
      expect(notification).toBeVisible();
      expect(notification).toHaveTextContent('Using offline content');
    });

    it('should gracefully degrade from dynamic to static content', async () => {
      // Setup: Partial API failures (some endpoints work, others fail)
      let callCount = 0;
      server.use(
        http.get('*/api/content/navigation.header.title', () => {
          return HttpResponse.json(mockContentItems['navigation.header.title']);
        }),
        http.get('*/api/content/*', () => {
          callCount++;
          if (callCount > 2) {
            return HttpResponse.error();
          }
          return HttpResponse.json({ type: 'text', value: 'Loading...' });
        })
      );

      // Pre-populate localStorage with fallback content
      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for mixed content loading
      await waitFor(() => {
        const titleDisplay = screen.getByTestId('content-display-navigation.header.title');
        expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome to Tchat');
      });

      // Verify partial fallback activation
      await waitFor(() => {
        const errorDisplay = screen.getByTestId('content-display-error.network.timeout');
        expect(within(errorDisplay).getByTestId('content-value')).toHaveTextContent('Connection error (Cached)');
        expect(within(errorDisplay).getByTestId('fallback-indicator')).toBeVisible();
      });

      // Verify sync status shows mixed state
      expect(screen.getByTestId('sync-status')).toHaveTextContent('error');
    });
  });

  describe('LocalStorage Fallback Mechanism', () => {
    it('should load content from localStorage when API is unavailable', async () => {
      // Setup: Complete API failure
      server.use(
        http.get('*/api/content/*', () => {
          return HttpResponse.error();
        })
      );

      // Pre-populate localStorage with extensive fallback content
      const extensiveFallback = {
        ...fallbackContentItems,
        'additional.content.key': {
          type: 'rich_text' as const,
          value: '<p>Rich text fallback</p>',
          format: 'html' as const,
        },
      };
      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(extensiveFallback));

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Verify localStorage is accessed
      await waitFor(() => {
        expect(localStorageMock.getItem).toHaveBeenCalledWith('tchat-content-fallback');
      });

      // Verify fallback content is loaded from localStorage
      await waitFor(() => {
        expect(screen.getByTestId('fallback-status')).toHaveTextContent('true');
      });

      // Verify content displays fallback values
      const displays = screen.getAllByTestId(/^content-display-/);
      displays.forEach(display => {
        expect(within(display).getByTestId('fallback-indicator')).toBeVisible();
      });
    });

    it('should handle corrupted localStorage data gracefully', async () => {
      // Setup: API failure with corrupted localStorage
      server.use(
        http.get('*/api/content/*', () => {
          return HttpResponse.error();
        })
      );

      // Corrupt localStorage data
      localStorageMock.setItem('tchat-content-fallback', 'invalid-json-data');

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Verify graceful handling of corrupted data
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('error');
      });

      // Verify fallback mode is not activated due to corrupt data
      expect(screen.getByTestId('fallback-status')).toHaveTextContent('false');

      // Verify error notification is shown
      const notification = screen.getByTestId('notification-message');
      expect(notification).toBeVisible();
      expect(notification).toHaveTextContent(/error|failed/i);
    });

    it('should persist successful content to localStorage for future fallback', async () => {
      // Setup: Successful API calls
      server.use(
        http.get('*/api/content/navigation.header.title', () => {
          return HttpResponse.json(mockContentItems['navigation.header.title']);
        }),
        http.get('*/api/content/error.network.timeout', () => {
          return HttpResponse.json(mockContentItems['error.network.timeout']);
        }),
        http.get('*/api/content/button.submit.text', () => {
          return HttpResponse.json(mockContentItems['button.submit.text']);
        })
      );

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for successful content loading
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('idle');
      });

      // Verify content is persisted to localStorage
      await waitFor(() => {
        expect(localStorageMock.setItem).toHaveBeenCalledWith(
          'tchat-content-fallback',
          expect.stringContaining('Welcome to Tchat')
        );
      });

      // Verify localStorage contains the fetched content
      const storedData = localStorageMock.setItem.mock.calls.find(
        call => call[0] === 'tchat-content-fallback'
      )?.[1];

      if (storedData) {
        const parsedData = JSON.parse(storedData);
        expect(parsedData['navigation.header.title']).toEqual(mockContentItems['navigation.header.title']);
      }
    });
  });

  describe('Error Boundary Handling', () => {
    it('should handle component errors during fallback rendering', async () => {
      // Setup: API failure and localStorage with valid fallback
      server.use(
        http.get('*/api/content/*', () => {
          return HttpResponse.error();
        })
      );

      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      // Mock console.error to catch error boundary logs
      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      const ErrorProneComponent = () => {
        // Simulate a component that throws during fallback rendering
        throw new Error('Fallback rendering error');
      };

      const TestAppWithError = () => (
        <div data-testid="test-app">
          <MockNotificationDisplay />
          <MockContentManager />
          <ErrorProneComponent />
          <MockContentDisplay contentKey="navigation.header.title" />
        </div>
      );

      render(
        <Provider store={testStore}>
          <TestAppWithError />
        </Provider>
      );

      // Verify error boundary catches the error
      await waitFor(() => {
        expect(consoleErrorSpy).toHaveBeenCalled();
      });

      // Verify other content still loads despite component error
      await waitFor(() => {
        expect(screen.getByTestId('content-display-navigation.header.title')).toBeInTheDocument();
      });

      consoleErrorSpy.mockRestore();
    });

    it('should display user-friendly error messages when content loading fails', async () => {
      // Setup: Complete failure scenario
      server.use(
        http.get('*/api/content/*', () => {
          return HttpResponse.json({ error: 'Service temporarily unavailable' }, { status: 503 });
        })
      );

      // No localStorage fallback available
      localStorageMock.clear();

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for error state
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('error');
      });

      // Verify user-friendly error notification
      const notification = screen.getByTestId('notification-message');
      expect(notification).toBeVisible();
      expect(notification).toHaveTextContent(/temporarily unavailable|try again/i);

      // Verify content displays show error state
      const displays = screen.getAllByTestId(/^content-display-/);
      displays.forEach(display => {
        const status = within(display).getByTestId('content-status');
        expect(status).toHaveTextContent('error');
      });
    });
  });

  describe('User Notification System', () => {
    it('should notify users when entering fallback mode', async () => {
      // Setup: API failure with localStorage fallback
      server.use(
        http.get('*/api/content/*', () => {
          return HttpResponse.error();
        })
      );

      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for fallback activation
      await waitFor(() => {
        expect(screen.getByTestId('fallback-status')).toHaveTextContent('true');
      });

      // Verify notification is displayed
      const notification = screen.getByTestId('notification-message');
      expect(notification).toBeVisible();
      expect(notification).toHaveTextContent(/offline|cached|outdated/i);

      // Verify notification persists during fallback mode
      await new Promise(resolve => setTimeout(resolve, 1000));
      expect(notification).toBeVisible();
    });

    it('should provide retry options to users', async () => {
      // Setup: Initial API failure
      server.use(
        http.get('*/api/content/*', () => {
          return HttpResponse.error();
        })
      );

      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      const user = userEvent.setup();

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for fallback mode
      await waitFor(() => {
        expect(screen.getByTestId('fallback-status')).toHaveTextContent('true');
      });

      // Verify retry button is available
      const retryButton = screen.getByTestId('retry-sync');
      expect(retryButton).toBeVisible();
      expect(retryButton).toBeEnabled();

      // Test retry functionality
      await user.click(retryButton);

      // Verify retry attempt
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('syncing');
      });
    });
  });

  describe('Recovery from Fallback Mode', () => {
    it('should recover automatically when API becomes available', async () => {
      // Setup: Initial API failure, then recovery
      let apiAvailable = false;

      server.use(
        http.get('*/api/content/*', ({ request }) => {
          if (!apiAvailable) {
            return HttpResponse.error();
          }

          const url = new URL(request.url);
          const contentKey = url.pathname.split('/').pop();
          const content = mockContentItems[contentKey as keyof typeof mockContentItems];

          if (content) {
            return HttpResponse.json(content);
          }

          return HttpResponse.json({ error: 'Content not found' }, { status: 404 });
        })
      );

      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      const user = userEvent.setup();

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for initial fallback mode
      await waitFor(() => {
        expect(screen.getByTestId('fallback-status')).toHaveTextContent('true');
      });

      // Verify fallback content is displayed
      const titleDisplay = screen.getByTestId('content-display-navigation.header.title');
      expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome (Offline)');

      // Simulate API recovery
      apiAvailable = true;

      // Trigger retry
      const retryButton = screen.getByTestId('retry-sync');
      await user.click(retryButton);

      // Wait for recovery
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('idle');
      }, { timeout: 5000 });

      // Verify fallback mode is disabled
      expect(screen.getByTestId('fallback-status')).toHaveTextContent('false');

      // Verify fresh content is displayed
      await waitFor(() => {
        expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome to Tchat');
      });

      // Verify fallback indicators are hidden
      expect(within(titleDisplay).getByTestId('fallback-indicator')).not.toBeVisible();

      // Verify notification is hidden
      const notification = screen.getByTestId('notification-message');
      expect(notification).not.toBeVisible();
    });

    it('should maintain data consistency during fallback-to-live transition', async () => {
      // Setup: Initial API failure with localStorage fallback
      let apiCallCount = 0;

      server.use(
        http.get('*/api/content/*', ({ request }) => {
          apiCallCount++;

          if (apiCallCount <= 3) {
            return HttpResponse.error();
          }

          const url = new URL(request.url);
          const contentKey = url.pathname.split('/').pop();
          const content = mockContentItems[contentKey as keyof typeof mockContentItems];

          if (content) {
            return HttpResponse.json({
              ...content,
              value: content.value + ' (Updated)',
            });
          }

          return HttpResponse.json({ error: 'Content not found' }, { status: 404 });
        })
      );

      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      const user = userEvent.setup();

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for fallback mode
      await waitFor(() => {
        expect(screen.getByTestId('fallback-status')).toHaveTextContent('true');
      });

      const titleDisplay = screen.getByTestId('content-display-navigation.header.title');

      // Verify initial fallback content
      expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome (Offline)');

      // Trigger recovery
      const retryButton = screen.getByTestId('retry-sync');
      await user.click(retryButton);

      // Wait for live content
      await waitFor(() => {
        expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome to Tchat (Updated)');
      }, { timeout: 5000 });

      // Verify no content flickering or inconsistent states
      expect(screen.getByTestId('fallback-status')).toHaveTextContent('false');
      expect(screen.getByTestId('sync-status')).toHaveTextContent('idle');

      // Verify updated content is persisted to localStorage
      await waitFor(() => {
        expect(localStorageMock.setItem).toHaveBeenCalledWith(
          'tchat-content-fallback',
          expect.stringContaining('Welcome to Tchat (Updated)')
        );
      });
    });

    it('should handle cache clearing and fresh content loading', async () => {
      // Setup: API available with fresh content
      server.use(
        http.get('*/api/content/*', ({ request }) => {
          const url = new URL(request.url);
          const contentKey = url.pathname.split('/').pop();
          const content = mockContentItems[contentKey as keyof typeof mockContentItems];

          if (content) {
            return HttpResponse.json({
              ...content,
              value: content.value + ' (Fresh)',
            });
          }

          return HttpResponse.json({ error: 'Content not found' }, { status: 404 });
        })
      );

      // Pre-populate localStorage with stale content
      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      const user = userEvent.setup();

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for initial load (should use fresh content)
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('idle');
      });

      const titleDisplay = screen.getByTestId('content-display-navigation.header.title');
      expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome to Tchat (Fresh)');

      // Test cache clearing
      const clearCacheButton = screen.getByTestId('clear-cache');
      await user.click(clearCacheButton);

      // Verify localStorage is cleared
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('tchat-content-fallback');

      // Verify content reloads fresh
      await waitFor(() => {
        expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome to Tchat (Fresh)');
      });
    });
  });

  describe('Content Consistency During Fallback', () => {
    it('should maintain content state consistency across components', async () => {
      // Setup: Mixed API success/failure scenario
      server.use(
        http.get('*/api/content/navigation.header.title', () => {
          return HttpResponse.json(mockContentItems['navigation.header.title']);
        }),
        http.get('*/api/content/error.network.timeout', () => {
          return HttpResponse.error();
        }),
        http.get('*/api/content/button.submit.text', () => {
          return HttpResponse.json(mockContentItems['button.submit.text']);
        })
      );

      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for mixed loading to complete
      await waitFor(() => {
        const titleDisplay = screen.getByTestId('content-display-navigation.header.title');
        expect(within(titleDisplay).getByTestId('content-value')).toHaveTextContent('Welcome to Tchat');
      });

      await waitFor(() => {
        const errorDisplay = screen.getByTestId('content-display-error.network.timeout');
        expect(within(errorDisplay).getByTestId('content-value')).toHaveTextContent('Connection error (Cached)');
      });

      // Verify fallback indicators show correct state per component
      const titleDisplay = screen.getByTestId('content-display-navigation.header.title');
      const errorDisplay = screen.getByTestId('content-display-error.network.timeout');
      const buttonDisplay = screen.getByTestId('content-display-button.submit.text');

      expect(within(titleDisplay).getByTestId('fallback-indicator')).not.toBeVisible();
      expect(within(errorDisplay).getByTestId('fallback-indicator')).toBeVisible();
      expect(within(buttonDisplay).getByTestId('fallback-indicator')).not.toBeVisible();

      // Verify global state reflects mixed content state
      expect(screen.getByTestId('sync-status')).toHaveTextContent('error');
    });

    it('should handle concurrent content updates and fallback scenarios', async () => {
      // Setup: Delayed and mixed API responses
      let responseCount = 0;

      server.use(
        http.get('*/api/content/*', async ({ request }) => {
          responseCount++;
          const url = new URL(request.url);
          const contentKey = url.pathname.split('/').pop();

          // Simulate different response times and failures
          if (contentKey === 'navigation.header.title') {
            await new Promise(resolve => setTimeout(resolve, 100));
            return HttpResponse.json(mockContentItems[contentKey as keyof typeof mockContentItems]);
          }

          if (contentKey === 'error.network.timeout') {
            await new Promise(resolve => setTimeout(resolve, 200));
            if (responseCount % 2 === 0) {
              return HttpResponse.error();
            }
            return HttpResponse.json(mockContentItems[contentKey as keyof typeof mockContentItems]);
          }

          return HttpResponse.json(mockContentItems[contentKey as keyof typeof mockContentItems]);
        })
      );

      localStorageMock.setItem('tchat-content-fallback', JSON.stringify(fallbackContentItems));

      render(
        <Provider store={testStore}>
          <TestApp />
        </Provider>
      );

      // Wait for various loading states to resolve
      await waitFor(() => {
        const displays = screen.getAllByTestId(/^content-display-/);
        const allLoaded = displays.every(display => {
          const value = within(display).getByTestId('content-value');
          return !value.textContent?.includes('Loading');
        });
        return allLoaded;
      }, { timeout: 10000 });

      // Verify final state consistency
      const syncStatus = screen.getByTestId('sync-status');
      expect(['idle', 'error']).toContain(syncStatus.textContent);

      // Verify no content shows loading state
      const displays = screen.getAllByTestId(/^content-display-/);
      displays.forEach(display => {
        const value = within(display).getByTestId('content-value');
        expect(value.textContent).not.toContain('Loading');
      });
    });
  });
});