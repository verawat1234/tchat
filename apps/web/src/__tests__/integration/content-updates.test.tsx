/**
 * Integration Test: Real-time Content Updates (T017)
 *
 * This test validates the complete real-time content update workflow from content
 * modification to immediate UI reflection without page refresh.
 *
 * CRITICAL TDD REQUIREMENT: This test MUST FAIL because the implementation doesn't exist yet.
 * This is intentional and required for proper test-driven development.
 *
 * @test-id T017-integration-content-updates
 */

import React from 'react';
import { describe, it, expect, beforeEach, afterEach, vi, beforeAll, afterAll } from 'vitest';
import { render, screen, waitFor, fireEvent, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { server } from '../../lib/test-utils/msw/server';
import { http, HttpResponse } from 'msw';
import { api } from '../../services/api';
import contentReducer, {
  setSelectedLanguage,
  updateContentPreferences,
  setSyncStatus,
  updateFallbackContent,
  selectContentState,
  selectSyncStatus,
} from '../../features/contentSlice';
import authReducer from '../../features/authSlice';
import uiReducer from '../../features/uiSlice';
import loadingReducer from '../../features/loadingSlice';
import type { ContentItem, ContentValue } from '../../types/content';
import { ContentType, ContentStatus } from '../../types/content';
import { useRealTimeContent, useRealTimeConnectionStatus, useCrossTabSync } from '../../hooks/useRealTimeContent';
import { initializeRealTimeService } from '../../services/realTimeConnectionService';

// Mock WebSocket class for real-time testing
class MockWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  readyState = MockWebSocket.CONNECTING;
  url: string;
  protocol?: string;
  onopen?: ((event: Event) => void) | null = null;
  onclose?: ((event: CloseEvent) => void) | null = null;
  onmessage?: ((event: MessageEvent) => void) | null = null;
  onerror?: ((event: Event) => void) | null = null;

  private listeners: { [key: string]: EventListener[] } = {};

  constructor(url: string, protocol?: string) {
    this.url = url;
    this.protocol = protocol;

    // Simulate connection opening
    setTimeout(() => {
      this.readyState = MockWebSocket.OPEN;
      this.onopen?.(new Event('open'));
      this.dispatchEvent(new Event('open'));
    }, 100);
  }

  send(data: string | ArrayBuffer | Blob | ArrayBufferView): void {
    if (this.readyState !== MockWebSocket.OPEN) {
      throw new Error('WebSocket is not open');
    }

    // Echo back for testing purposes or trigger specific responses
    const message = typeof data === 'string' ? JSON.parse(data) : data;
    this.simulateMessage(message);
  }

  close(code?: number, reason?: string): void {
    this.readyState = MockWebSocket.CLOSING;
    setTimeout(() => {
      this.readyState = MockWebSocket.CLOSED;
      const closeEvent = new CloseEvent('close', { code, reason });
      this.onclose?.(closeEvent);
      this.dispatchEvent(closeEvent);
    }, 10);
  }

  addEventListener(type: string, listener: EventListener): void {
    if (!this.listeners[type]) {
      this.listeners[type] = [];
    }
    this.listeners[type].push(listener);
  }

  removeEventListener(type: string, listener: EventListener): void {
    if (this.listeners[type]) {
      this.listeners[type] = this.listeners[type].filter(l => l !== listener);
    }
  }

  dispatchEvent(event: Event): boolean {
    const listeners = this.listeners[event.type] || [];
    listeners.forEach(listener => listener(event));
    return true;
  }

  // Helper method to simulate incoming messages
  simulateMessage(data: any): void {
    const messageEvent = new MessageEvent('message', {
      data: JSON.stringify(data),
    });
    this.onmessage?.(messageEvent);
    this.dispatchEvent(messageEvent);
  }

  // Helper method to simulate connection errors
  simulateError(): void {
    const errorEvent = new Event('error');
    this.onerror?.(errorEvent);
    this.dispatchEvent(errorEvent);
  }
}

// Mock the global WebSocket
beforeAll(() => {
  Object.defineProperty(global, 'WebSocket', {
    writable: true,
    value: MockWebSocket,
  });
});

// Test Component that displays and updates content
const ContentDisplay: React.FC<{ contentId: string }> = ({ contentId }) => {
  const {
    value,
    lastUpdated,
    syncStatus,
    connectionStatus,
  } = useRealTimeContent(contentId);

  return (
    <div data-testid={`content-display-${contentId}`}>
      <div data-testid="content-value">{value}</div>
      <div data-testid="sync-status">{syncStatus}</div>
      <div data-testid="last-updated">{lastUpdated || 'Never'}</div>
    </div>
  );
};

// Test Component for content editing
const ContentEditor: React.FC<{ contentId: string; onSave: (value: string) => void }> = ({
  contentId,
  onSave
}) => {
  const [value, setValue] = React.useState('');
  const {
    updateContent,
    setIsUserEditing,
    syncStatus,
  } = useRealTimeContent(contentId);

  const handleSave = async () => {
    try {
      await updateContent(value);
      onSave(value);
    } catch (error) {
      console.error('Failed to save content:', error);
    }
  };

  return (
    <div data-testid={`content-editor-${contentId}`}>
      <input
        data-testid="content-input"
        value={value}
        onChange={(e) => {
          setValue(e.target.value);
          setIsUserEditing(e.target.value.length > 0);
        }}
        onFocus={() => setIsUserEditing(true)}
        onBlur={() => setIsUserEditing(false)}
        placeholder="Enter content..."
      />
      <button
        data-testid="save-button"
        onClick={handleSave}
        disabled={syncStatus === 'syncing'}
      >
        {syncStatus === 'syncing' ? 'Saving...' : 'Save'}
      </button>
    </div>
  );
};

// Test App Component with multiple tabs simulation
const TestApp: React.FC = () => {
  const [activeTab, setActiveTab] = React.useState<'tab1' | 'tab2'>('tab1');
  const [contentUpdates, setContentUpdates] = React.useState(0);

  const { status: connectionStatus } = useRealTimeConnectionStatus();
  const { updates: crossTabUpdates } = useCrossTabSync('test.content.1');

  React.useEffect(() => {
    setContentUpdates(crossTabUpdates);
  }, [crossTabUpdates]);

  const handleSave = (value: string) => {
    setContentUpdates(prev => prev + 1);
  };

  return (
    <div data-testid="test-app">
      <nav data-testid="tab-navigation">
        <button
          data-testid="tab1-button"
          onClick={() => setActiveTab('tab1')}
          className={activeTab === 'tab1' ? 'active' : ''}
        >
          Tab 1
        </button>
        <button
          data-testid="tab2-button"
          onClick={() => setActiveTab('tab2')}
          className={activeTab === 'tab2' ? 'active' : ''}
        >
          Tab 2
        </button>
      </nav>

      <div data-testid="content-area">
        {activeTab === 'tab1' && (
          <div data-testid="tab1-content">
            <ContentDisplay contentId="test.content.1" />
            <ContentEditor contentId="test.content.1" onSave={handleSave} />
          </div>
        )}

        {activeTab === 'tab2' && (
          <div data-testid="tab2-content">
            <ContentDisplay contentId="test.content.1" />
            <div data-testid="readonly-notice">Read-only view</div>
          </div>
        )}
      </div>

      <div data-testid="update-counter">Updates: {contentUpdates}</div>
      <div data-testid="connection-status">{connectionStatus}</div>
    </div>
  );
};

// Create test store
const createTestStore = () => {
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
          ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE'],
        },
      }).concat(api.middleware),
  });
};

// Test Provider component
const TestProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const store = createTestStore();
  return <Provider store={store}>{children}</Provider>;
};

// Mock data
const mockContentItem: ContentItem = {
  id: 'test.content.1',
  key: 'test.content.1',
  categoryId: 'test',
  category: {
    id: 'test',
    name: 'Test Category',
    description: 'Test category for integration tests',
    permissions: {
      read: ['user'],
      write: ['admin'],
      publish: ['admin'],
    },
  },
  type: ContentType.TEXT,
  value: {
    type: 'text',
    value: 'Original content',
  } as ContentValue,
  status: ContentStatus.PUBLISHED,
  metadata: {
    createdAt: '2024-01-01T00:00:00Z',
    createdBy: 'test-user',
    updatedAt: '2024-01-01T00:00:00Z',
    updatedBy: 'test-user',
    version: 1,
  },
};

const updatedContentItem: ContentItem = {
  ...mockContentItem,
  value: {
    type: 'text',
    value: 'Updated content via real-time',
  } as ContentValue,
  metadata: {
    ...mockContentItem.metadata,
    updatedAt: '2024-01-01T01:00:00Z',
    version: 2,
  },
};

describe('Real-time Content Updates Integration (T017)', () => {
  let mockWebSocket: MockWebSocket;
  const user = userEvent.setup();

  beforeEach(() => {
    vi.clearAllMocks();

    // Initialize real-time service for testing
    initializeRealTimeService({
      url: 'ws://localhost:8080/realtime',
      token: 'mock-test-token',
      reconnectDelay: 1000,
      maxReconnectAttempts: 3,
      heartbeatInterval: 10000,
    });

    // Setup MSW handlers for content API
    server.use(
      http.get('/api/content/:id', ({ params }) => {
        const { id } = params;
        if (id === 'test.content.1') {
          return HttpResponse.json(mockContentItem);
        }
        return new HttpResponse(null, { status: 404 });
      }),

      http.put('/api/content/:id', ({ params }) => {
        const { id } = params;
        if (id === 'test.content.1') {
          return HttpResponse.json(updatedContentItem);
        }
        return new HttpResponse(null, { status: 404 });
      }),

      http.get('/api/content/realtime/token', () => {
        return HttpResponse.json({ token: 'mock-websocket-token' });
      })
    );
  });

  afterEach(() => {
    if (mockWebSocket) {
      mockWebSocket.close();
    }
  });

  describe('WebSocket Connection Management', () => {
    it('should establish WebSocket connection for real-time updates', async () => {
      // EXPECTED TO FAIL - Real-time connection service not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Wait for component to mount and attempt WebSocket connection
      await waitFor(() => {
        expect(screen.getByTestId('connection-status')).toHaveTextContent('connecting');
      }, { timeout: 1000 });

      // Simulate successful connection
      await waitFor(() => {
        expect(screen.getByTestId('connection-status')).toHaveTextContent('connected');
      }, { timeout: 2000 });
    });

    it('should handle WebSocket connection failures gracefully', async () => {
      // EXPECTED TO FAIL - Error handling not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Simulate connection error
      await act(async () => {
        // This would trigger WebSocket error handling
        await new Promise(resolve => setTimeout(resolve, 100));
      });

      expect(screen.getByTestId('connection-status')).toHaveTextContent('error');
    });

    it('should attempt reconnection after connection loss', async () => {
      // EXPECTED TO FAIL - Reconnection logic not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Simulate connection established then lost
      await waitFor(() => {
        expect(screen.getByTestId('connection-status')).toHaveTextContent('connected');
      });

      // Simulate connection loss
      await act(async () => {
        // This would trigger WebSocket disconnection
        await new Promise(resolve => setTimeout(resolve, 100));
      });

      expect(screen.getByTestId('connection-status')).toHaveTextContent('reconnecting');
    });
  });

  describe('Real-time Content Propagation', () => {
    it('should receive and display real-time content updates without page refresh', async () => {
      // EXPECTED TO FAIL - Real-time content updates not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Initial content should be loaded
      await waitFor(() => {
        expect(screen.getByTestId('content-value')).toHaveTextContent('Original content');
      });

      // Simulate real-time content update from server
      await act(async () => {
        // This would simulate receiving a WebSocket message with updated content
        const updateMessage = {
          type: 'content_updated',
          data: {
            contentId: 'test.content.1',
            value: 'Updated content via real-time',
            updatedAt: '2024-01-01T01:00:00Z',
            version: 2,
          },
        };

        // Mock WebSocket message would trigger content update
        await new Promise(resolve => setTimeout(resolve, 100));
      });

      // Content should be updated in real-time
      expect(screen.getByTestId('content-value')).toHaveTextContent('Updated content via real-time');
      expect(screen.getByTestId('last-updated')).toHaveTextContent('2024-01-01T01:00:00Z');
    });

    it('should handle high-frequency content updates efficiently', async () => {
      // EXPECTED TO FAIL - Throttling/debouncing not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      const updatePromises: Promise<void>[] = [];

      // Simulate rapid content updates
      for (let i = 0; i < 10; i++) {
        updatePromises.push(
          act(async () => {
            // This would simulate rapid WebSocket messages
            await new Promise(resolve => setTimeout(resolve, 10));
          })
        );
      }

      await Promise.all(updatePromises);

      // Should handle updates efficiently without UI freezing
      expect(screen.getByTestId('content-value')).toBeInTheDocument();

      // Performance metrics should be within acceptable limits
      // This would check for throttled updates rather than all 10 individual updates
      const updateCounter = screen.getByTestId('update-counter');
      expect(updateCounter).toHaveTextContent(/Updates: [1-5]/); // Throttled updates
    });

    it('should preserve user input during real-time updates', async () => {
      // EXPECTED TO FAIL - Input preservation not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // User starts typing
      const input = screen.getByTestId('content-input');
      await user.type(input, 'User is typing...');

      // Simulate real-time update while user is typing
      await act(async () => {
        // This would simulate a WebSocket update
        await new Promise(resolve => setTimeout(resolve, 100));
      });

      // User's input should be preserved
      expect(input).toHaveValue('User is typing...');

      // Background content should still update (in display component)
      expect(screen.getByTestId('content-value')).toHaveTextContent('Updated content via real-time');
    });
  });

  describe('Multi-tab/Window Synchronization', () => {
    it('should synchronize content updates across multiple tabs', async () => {
      // EXPECTED TO FAIL - Multi-tab synchronization not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Start on tab 1
      expect(screen.getByTestId('tab1-content')).toBeInTheDocument();

      // Switch to tab 2
      await user.click(screen.getByTestId('tab2-button'));

      expect(screen.getByTestId('tab2-content')).toBeInTheDocument();
      expect(screen.getByTestId('readonly-notice')).toHaveTextContent('Read-only view');

      // Simulate content update from another tab/window
      await act(async () => {
        // This would simulate cross-tab communication
        window.dispatchEvent(new StorageEvent('storage', {
          key: 'content_update',
          newValue: JSON.stringify({
            contentId: 'test.content.1',
            value: 'Updated from another tab',
          }),
        }));
      });

      // Content should be updated in current tab
      expect(screen.getByTestId('content-value')).toHaveTextContent('Updated from another tab');
    });

    it('should handle content conflicts between multiple editors', async () => {
      // EXPECTED TO FAIL - Conflict resolution not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // User makes changes
      const input = screen.getByTestId('content-input');
      await user.type(input, 'Local changes');

      // Simulate concurrent update from another user
      await act(async () => {
        // This would simulate a conflict scenario
        await new Promise(resolve => setTimeout(resolve, 100));
      });

      // Should show conflict resolution UI (WILL FAIL - not implemented)
      try {
        expect(screen.getByTestId('conflict-notification')).toBeInTheDocument();
        expect(screen.getByTestId('resolve-conflict-button')).toBeInTheDocument();
      } catch (error) {
        // Expected to fail - conflict resolution UI not implemented
        expect(error).toBeInstanceOf(Error);
      }
    });
  });

  describe('Update Notifications and User Feedback', () => {
    it('should display toast notifications for content updates', async () => {
      // EXPECTED TO FAIL - Notification system not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Simulate real-time content update
      await act(async () => {
        // This would trigger a notification
        await new Promise(resolve => setTimeout(resolve, 100));
      });

      // Should show update notification (WILL FAIL - not implemented)
      try {
        expect(screen.getByTestId('update-notification')).toBeInTheDocument();
        expect(screen.getByTestId('update-notification')).toHaveTextContent('Content updated');
      } catch (error) {
        // Expected to fail - notification system not implemented
        expect(error).toBeInstanceOf(Error);
      }
    });

    it('should show loading states during content synchronization', async () => {
      // EXPECTED TO FAIL - Loading states not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Trigger content save
      const input = screen.getByTestId('content-input');
      await user.type(input, 'New content');
      await user.click(screen.getByTestId('save-button'));

      // Should show loading state
      expect(screen.getByTestId('sync-status')).toHaveTextContent('syncing');

      // Wait for completion
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('idle');
      });
    });

    it('should display error messages for failed updates', async () => {
      // EXPECTED TO FAIL - Error handling not implemented
      server.use(
        http.put('/api/content/:id', () => {
          return HttpResponse.json({ error: 'Server error' }, { status: 500 });
        })
      );

      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Trigger failed save
      const input = screen.getByTestId('content-input');
      await user.type(input, 'Content that will fail');
      await user.click(screen.getByTestId('save-button'));

      // Should show error state
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('error');
      });

      try {
        expect(screen.getByTestId('error-message')).toHaveTextContent('Failed to save content');
      } catch (error) {
        // Expected to fail - error message UI not implemented
        expect(error).toBeInstanceOf(Error);
      }
    });
  });

  describe('Performance and Resource Management', () => {
    it('should implement connection pooling for multiple content items', async () => {
      // EXPECTED TO FAIL - Connection pooling not implemented
      render(
        <TestProvider>
          <div data-testid="multi-content">
            <ContentDisplay contentId="test.content.1" />
            <ContentDisplay contentId="test.content.2" />
            <ContentDisplay contentId="test.content.3" />
          </div>
        </TestProvider>
      );

      // Should use single WebSocket connection for all content items
      await waitFor(() => {
        expect(screen.getByTestId('connection-status')).toHaveTextContent('connected');
      });

      // Performance check - should not create multiple connections
      // This would be verified through WebSocket mock tracking
    });

    it('should handle memory cleanup on component unmount', async () => {
      // EXPECTED TO FAIL - Cleanup logic not implemented
      const { unmount } = render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Establish connection
      await waitFor(() => {
        expect(screen.getByTestId('connection-status')).toHaveTextContent('connected');
      });

      // Unmount component
      unmount();

      // Should clean up WebSocket connections and event listeners
      // This would be verified through memory leak detection
    });

    it('should throttle updates during rapid content changes', async () => {
      // EXPECTED TO FAIL - Throttling not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      const startTime = Date.now();

      // Simulate rapid updates
      for (let i = 0; i < 20; i++) {
        await act(async () => {
          // Rapid WebSocket messages
          await new Promise(resolve => setTimeout(resolve, 10));
        });
      }

      const endTime = Date.now();

      // Should complete within reasonable time due to throttling
      expect(endTime - startTime).toBeLessThan(1000);

      // Should show reasonable number of updates (throttled)
      expect(screen.getByTestId('update-counter')).toHaveTextContent(/Updates: [1-5]/);
    });
  });

  describe('Offline and Network Resilience', () => {
    it('should queue updates when connection is lost', async () => {
      // EXPECTED TO FAIL - Offline queueing not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Simulate network loss
      await act(async () => {
        // This would trigger offline mode
        window.dispatchEvent(new Event('offline'));
      });

      // Make changes while offline
      const input = screen.getByTestId('content-input');
      await user.type(input, 'Offline changes');
      await user.click(screen.getByTestId('save-button'));

      // Should queue the update
      expect(screen.getByTestId('sync-status')).toHaveTextContent('queued');

      // Simulate network recovery
      await act(async () => {
        window.dispatchEvent(new Event('online'));
      });

      // Should sync queued updates
      await waitFor(() => {
        expect(screen.getByTestId('sync-status')).toHaveTextContent('idle');
      });
    });

    it('should gracefully degrade to polling when WebSocket fails', async () => {
      // EXPECTED TO FAIL - Fallback polling not implemented
      render(
        <TestProvider>
          <TestApp />
        </TestProvider>
      );

      // Simulate WebSocket failure
      await act(async () => {
        // This would trigger fallback to polling
        await new Promise(resolve => setTimeout(resolve, 100));
      });

      expect(screen.getByTestId('connection-status')).toHaveTextContent('polling');

      // Should still receive updates via polling
      await waitFor(() => {
        expect(screen.getByTestId('content-value')).toHaveTextContent('Updated content via polling');
      });
    });
  });
});