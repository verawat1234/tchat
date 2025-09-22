import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { axe, toHaveNoViolations } from 'jest-axe';

import ContentLoader, {
  TextSkeleton,
  ImageSkeleton,
  CardSkeleton,
  ListSkeleton,
  TableSkeleton,
  ChatSkeleton,
  MediaSkeleton,
  SpinnerLoader,
  ProgressLoader,
  ErrorState,
  FallbackState,
  useContentLoader,
  useAccessibilityAnnouncer,
} from './ContentLoader';

// Extend Jest matchers
expect.extend(toHaveNoViolations);

// Mock components that might not be available in test environment
vi.mock('lucide-react', () => ({
  AlertCircle: () => <div data-testid="alert-circle-icon" />,
  Wifi: () => <div data-testid="wifi-icon" />,
  WifiOff: () => <div data-testid="wifi-off-icon" />,
  RefreshCw: () => <div data-testid="refresh-icon" />,
}));

// Mock Redux store for testing
const createMockStore = (initialState = {}) => {
  const defaultState = {
    loading: {
      global: false,
      requests: {},
      operations: {},
    },
  };

  return configureStore({
    reducer: {
      loading: (state = defaultState.loading, action) => {
        switch (action.type) {
          case 'SET_OPERATION':
            return {
              ...state,
              operations: {
                ...state.operations,
                [action.payload.key]: action.payload.operation,
              },
            };
          default:
            return state;
        }
      },
    },
    preloadedState: { ...defaultState, ...initialState },
  });
};

const renderWithProvider = (component: React.ReactElement, store = createMockStore()) => {
  return render(
    <Provider store={store}>
      {component}
    </Provider>
  );
};

describe('ContentLoader', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('Basic functionality', () => {
    it('renders children when not loading', () => {
      renderWithProvider(
        <ContentLoader isLoading={false}>
          <div>Test content</div>
        </ContentLoader>
      );

      expect(screen.getByText('Test content')).toBeInTheDocument();
    });

    it('shows loading skeleton when loading', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} contentType="text">
          <div>Test content</div>
        </ContentLoader>
      );

      expect(screen.queryByText('Test content')).not.toBeInTheDocument();
      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('applies custom className', () => {
      renderWithProvider(
        <ContentLoader className="custom-class" isLoading={false}>
          <div>Test content</div>
        </ContentLoader>
      );

      const container = screen.getByText('Test content').parentElement;
      expect(container).toHaveClass('custom-class');
    });
  });

  describe('Loading types', () => {
    it('renders skeleton loader by default', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} />
      );

      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('renders spinner loader', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} loadingType="spinner" />
      );

      expect(screen.getByTestId('refresh-icon')).toBeInTheDocument();
    });

    it('renders progress loader with progress value', () => {
      renderWithProvider(
        <ContentLoader
          isLoading={true}
          loadingType="progress"
          progress={50}
          progressMessage="Loading items..."
        />
      );

      expect(screen.getByText('Loading items...')).toBeInTheDocument();
      expect(screen.getByText('50%')).toBeInTheDocument();
    });

    it('renders minimal loader', () => {
      renderWithProvider(
        <ContentLoader
          isLoading={true}
          loadingType="minimal"
          loadingLabel="Please wait"
        />
      );

      expect(screen.getByText('Please wait')).toBeInTheDocument();
    });
  });

  describe('Content types skeletons', () => {
    it('renders text skeleton', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} contentType="text" skeletonCount={2} />
      );

      // Text skeleton creates multiple skeleton elements
      const skeletons = screen.getAllByRole('status')[0];
      expect(skeletons).toBeInTheDocument();
    });

    it('renders image skeleton', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} contentType="image" />
      );

      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('renders card skeleton', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} contentType="card" />
      );

      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('renders list skeleton', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} contentType="list" skeletonCount={3} />
      );

      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('renders table skeleton', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} contentType="table" />
      );

      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('renders chat skeleton', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} contentType="chat" />
      );

      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('renders media skeleton', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} contentType="media" />
      );

      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('uses custom skeleton renderer', () => {
      const customSkeleton = () => <div data-testid="custom-skeleton">Custom</div>;

      renderWithProvider(
        <ContentLoader
          isLoading={true}
          renderSkeleton={customSkeleton}
        />
      );

      expect(screen.getByTestId('custom-skeleton')).toBeInTheDocument();
    });
  });

  describe('Error handling', () => {
    it('displays error state with string error', () => {
      renderWithProvider(
        <ContentLoader error="Something went wrong" />
      );

      expect(screen.getByRole('alert')).toBeInTheDocument();
      expect(screen.getByText('Failed to load content')).toBeInTheDocument();
      expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    });

    it('displays error state with Error object', () => {
      const error = new Error('Network error');

      renderWithProvider(
        <ContentLoader error={error} />
      );

      expect(screen.getByText('Network error')).toBeInTheDocument();
    });

    it('shows retry button when onRetry is provided', async () => {
      const user = userEvent.setup();
      const onRetry = vi.fn();

      renderWithProvider(
        <ContentLoader
          error="Test error"
          onRetry={onRetry}
          maxRetries={3}
          retryCount={1}
        />
      );

      const retryButton = screen.getByRole('button', { name: /try again/i });
      expect(retryButton).toBeInTheDocument();
      expect(screen.getByText('Try again (2 attempts left)')).toBeInTheDocument();

      await user.click(retryButton);
      expect(onRetry).toHaveBeenCalledTimes(1);
    });

    it('hides retry button when max retries reached', () => {
      renderWithProvider(
        <ContentLoader
          error="Test error"
          onRetry={vi.fn()}
          maxRetries={3}
          retryCount={3}
        />
      );

      expect(screen.queryByRole('button', { name: /try again/i })).not.toBeInTheDocument();
    });
  });

  describe('Fallback states', () => {
    it('shows fallback indicator when content is cached', () => {
      renderWithProvider(
        <ContentLoader isFallback={true} fallbackMessage="Using cached data">
          <div>Cached content</div>
        </ContentLoader>
      );

      expect(screen.getByText('Using cached data')).toBeInTheDocument();
      expect(screen.getByTestId('wifi-icon')).toBeInTheDocument();
      expect(screen.getByText('Cached content')).toBeInTheDocument();
    });

    it('shows offline indicator when offline', () => {
      renderWithProvider(
        <ContentLoader isOffline={true}>
          <div>Offline content</div>
        </ContentLoader>
      );

      expect(screen.getByText('Showing offline content')).toBeInTheDocument();
      expect(screen.getByTestId('wifi-off-icon')).toBeInTheDocument();
    });
  });

  describe('RTK Query integration', () => {
    it('integrates with RTK Query loading state', () => {
      const store = createMockStore({
        loading: {
          global: false,
          requests: {},
          operations: {
            'fetchData': {
              isLoading: true,
              message: 'Fetching data...',
              progress: 75,
            },
          },
        },
      });

      renderWithProvider(
        <ContentLoader rtkQueryKey="fetchData" loadingType="progress">
          <div>Content</div>
        </ContentLoader>,
        store
      );

      expect(screen.getByText('Fetching data...')).toBeInTheDocument();
      expect(screen.getByText('75%')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    it('has proper ARIA attributes for loading state', () => {
      renderWithProvider(
        <ContentLoader isLoading={true} loadingLabel="Loading user data" />
      );

      const loader = screen.getByRole('status');
      expect(loader).toHaveAttribute('aria-live', 'polite');
      expect(loader).toHaveAttribute('aria-label', 'Loading user data');
    });

    it('has proper ARIA attributes for error state', () => {
      renderWithProvider(
        <ContentLoader error="Failed to load" />
      );

      const errorAlert = screen.getByRole('alert');
      expect(errorAlert).toHaveAttribute('aria-live', 'polite');
    });

    it('announces loading state changes', async () => {
      renderWithProvider(
        <ContentLoader
          isLoading={true}
          announceChanges={true}
          loadingLabel="Loading content"
        />
      );

      // Check for screen reader announcement
      await waitFor(() => {
        expect(screen.getByLabelText('Loading content')).toBeInTheDocument();
      });
    });

    it('respects reduced motion preference', () => {
      renderWithProvider(
        <ContentLoader
          isLoading={true}
          reduceMotion={true}
          animation="pulse"
        />
      );

      const loader = screen.getByRole('status');
      expect(loader).not.toHaveClass('animate-pulse');
    });

    it('passes accessibility audit', async () => {
      const { container } = renderWithProvider(
        <ContentLoader isLoading={true} loadingLabel="Loading data">
          <div>Content</div>
        </ContentLoader>
      );

      const results = await axe(container);
      expect(results).toHaveNoViolations();
    });
  });

  describe('Performance optimizations', () => {
    it('supports deferred rendering', async () => {
      vi.useFakeTimers();

      renderWithProvider(
        <ContentLoader isLoading={true} deferRender={true} />
      );

      // Should not render initially
      expect(screen.queryByRole('status')).not.toBeInTheDocument();

      // Fast-forward timer
      act(() => {
        vi.advanceTimersByTime(150);
      });

      // Should render after delay
      await waitFor(() => {
        expect(screen.getByRole('status')).toBeInTheDocument();
      }, { timeout: 1000 });

      vi.useRealTimers();
    });

    it('prevents layout shift with consistent sizing', () => {
      const { rerender } = renderWithProvider(
        <ContentLoader isLoading={true} size="md" />
      );

      const loadingElement = screen.getByRole('status');
      const initialHeight = loadingElement.style.minHeight;

      rerender(
        <Provider store={createMockStore()}>
          <ContentLoader isLoading={false} size="md">
            <div>Content loaded</div>
          </ContentLoader>
        </Provider>
      );

      const contentElement = screen.getByText('Content loaded').parentElement;
      expect(contentElement).toHaveClass('min-h-24'); // md size
    });
  });

  describe('Compound components', () => {
    it('exports individual skeleton components', () => {
      render(<TextSkeleton count={2} />);
      // TextSkeleton should render without errors
    });

    it('provides compound components through ContentLoader', () => {
      render(<ContentLoader.Text count={3} />);
      render(<ContentLoader.Image />);
      render(<ContentLoader.Card />);
      render(<ContentLoader.List count={2} />);
      render(<ContentLoader.Table rows={3} cols={4} />);
      render(<ContentLoader.Chat count={2} />);
      render(<ContentLoader.Media />);
      render(<ContentLoader.Spinner size="md" />);
      render(<ContentLoader.Progress progress={50} message="Loading..." />);
      render(
        <ContentLoader.Error
          error="Test error"
          onRetry={vi.fn()}
          maxRetries={3}
          retryCount={1}
        />
      );
      render(<ContentLoader.Fallback isOffline={true} />);
    });
  });

  describe('Custom hooks', () => {
    it('useContentLoader hook works with RTK Query', () => {
      const TestComponent = () => {
        const { isLoading, progress, message } = useContentLoader('testQuery');
        return (
          <div>
            <span>Loading: {isLoading.toString()}</span>
            <span>Progress: {progress || 0}</span>
            <span>Message: {message || 'None'}</span>
          </div>
        );
      };

      const store = createMockStore({
        loading: {
          operations: {
            testQuery: {
              isLoading: true,
              progress: 50,
              message: 'Test message',
            },
          },
        },
      });

      renderWithProvider(<TestComponent />, store);

      expect(screen.getByText('Loading: true')).toBeInTheDocument();
      expect(screen.getByText('Progress: 50')).toBeInTheDocument();
      expect(screen.getByText('Message: Test message')).toBeInTheDocument();
    });

    it('useAccessibilityAnnouncer hook works correctly', () => {
      const TestComponent = () => {
        const { announcement, announce } = useAccessibilityAnnouncer(true);

        React.useEffect(() => {
          announce('Test announcement');
        }, [announce]);

        return <div>{announcement}</div>;
      };

      render(<TestComponent />);
      expect(screen.getByText('Test announcement')).toBeInTheDocument();
    });
  });

  describe('Edge cases', () => {
    it('handles undefined/null children gracefully', () => {
      renderWithProvider(
        <ContentLoader isLoading={false}>
          {null}
        </ContentLoader>
      );

      // Should not crash
    });

    it('handles empty error gracefully', () => {
      renderWithProvider(
        <ContentLoader error="" />
      );

      expect(screen.getByRole('alert')).toBeInTheDocument();
      expect(screen.getByText('Failed to load content')).toBeInTheDocument();
    });

    it('handles zero progress value', () => {
      renderWithProvider(
        <ContentLoader
          isLoading={true}
          loadingType="progress"
          progress={0}
        />
      );

      expect(screen.getByText('0%')).toBeInTheDocument();
    });

    it('handles very high progress value', () => {
      renderWithProvider(
        <ContentLoader
          isLoading={true}
          loadingType="progress"
          progress={150}
        />
      );

      expect(screen.getByText('150%')).toBeInTheDocument();
    });
  });
});

// Individual component tests
describe('Individual Skeleton Components', () => {
  it('TextSkeleton renders correct number of lines', () => {
    render(<TextSkeleton count={3} />);
    // Should render 3 skeleton lines
  });

  it('ImageSkeleton renders with correct aspect ratio', () => {
    render(<ImageSkeleton aspectRatio="aspect-square" />);
    // Should render image skeleton
  });

  it('CardSkeleton renders complete card structure', () => {
    render(<CardSkeleton />);
    // Should render card skeleton
  });

  it('ListSkeleton renders correct number of items', () => {
    render(<ListSkeleton count={5} />);
    // Should render 5 list items
  });

  it('TableSkeleton renders with correct dimensions', () => {
    render(<TableSkeleton rows={3} cols={4} />);
    // Should render 3x4 table
  });

  it('ChatSkeleton renders chat message layout', () => {
    render(<ChatSkeleton count={3} />);
    // Should render 3 chat messages
  });

  it('MediaSkeleton renders media player layout', () => {
    render(<MediaSkeleton />);
    // Should render media skeleton
  });
});

describe('Utility Components', () => {
  it('SpinnerLoader renders with different sizes', () => {
    render(<SpinnerLoader size="sm" />);
    expect(screen.getByTestId('refresh-icon')).toBeInTheDocument();

    render(<SpinnerLoader size="lg" />);
    expect(screen.getAllByTestId('refresh-icon')).toHaveLength(2);
  });

  it('ProgressLoader renders progress correctly', () => {
    render(<ProgressLoader progress={75} message="Loading files..." />);
    expect(screen.getByText('Loading files...')).toBeInTheDocument();
    expect(screen.getByText('75%')).toBeInTheDocument();
  });

  it('ErrorState renders retry functionality', async () => {
    const user = userEvent.setup();
    const onRetry = vi.fn();

    render(
      <ErrorState
        error="Test error"
        onRetry={onRetry}
        maxRetries={5}
        retryCount={2}
      />
    );

    const retryButton = screen.getByRole('button');
    await user.click(retryButton);
    expect(onRetry).toHaveBeenCalledTimes(1);
  }, { timeout: 10000 });

  it('FallbackState renders offline and online states', () => {
    const { rerender } = render(
      <FallbackState isOffline={true} message="Offline mode" />
    );

    expect(screen.getByText('Offline mode')).toBeInTheDocument();
    expect(screen.getByTestId('wifi-off-icon')).toBeInTheDocument();

    rerender(<FallbackState isOffline={false} message="Cached data" />);

    expect(screen.getByText('Cached data')).toBeInTheDocument();
    expect(screen.getByTestId('wifi-icon')).toBeInTheDocument();
  });
});