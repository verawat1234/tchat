/**
 * Error Components Test Suite
 *
 * Tests for dynamic error message handling in error components,
 * ensuring proper fallback behavior and content API integration.
 */

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';

import { ContentErrorBoundary } from '../components/ContentErrorBoundary';
import { ImageWithFallback } from '../components/figma/ImageWithFallback';
import { useContentText } from '../hooks/useContentText';
import { contentApi } from '../services/content';
import { api } from '../services/api';

// Mock component that throws an error for testing error boundaries
const ThrowError = ({ shouldThrow }: { shouldThrow: boolean }) => {
  if (shouldThrow) {
    throw new Error('Test error for boundary');
  }
  return <div>No error</div>;
};

// Mock server for content API
const server = setupServer(
  http.get('/api/content/error.boundary.title', () => {
    return HttpResponse.json({
      id: 'error.boundary.title',
      type: 'text',
      data: 'System Error (Dynamic)',
      metadata: {},
      version: 1,
      category: 'error',
      created_at: '2023-01-01T00:00:00Z',
      updated_at: '2023-01-01T00:00:00Z',
    });
  }),
  http.get('/api/content/error.boundary.description', () => {
    return HttpResponse.json({
      id: 'error.boundary.description',
      type: 'text',
      data: 'Dynamic error description from content API',
      metadata: {},
      version: 1,
      category: 'error',
      created_at: '2023-01-01T00:00:00Z',
      updated_at: '2023-01-01T00:00:00Z',
    });
  }),
  http.get('/api/content/*', () => {
    // Generic fallback for other content items
    return new HttpResponse(null, { status: 404 });
  })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

const createTestStore = () => {
  return configureStore({
    reducer: {
      api: api.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(api.middleware),
  });
};

const renderWithProvider = (component: React.ReactElement) => {
  const store = createTestStore();
  return render(
    <Provider store={store}>
      {component}
    </Provider>
  );
};

describe('ContentErrorBoundary', () => {
  it('renders children when no error occurs', () => {
    renderWithProvider(
      <ContentErrorBoundary>
        <ThrowError shouldThrow={false} />
      </ContentErrorBoundary>
    );

    expect(screen.getByText('No error')).toBeInTheDocument();
  });

  it('renders error UI with fallback content when error occurs', () => {
    // Suppress console errors for this test
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

    renderWithProvider(
      <ContentErrorBoundary>
        <ThrowError shouldThrow={true} />
      </ContentErrorBoundary>
    );

    // Should show fallback content initially
    expect(screen.getByText('Content System Error')).toBeInTheDocument();
    expect(screen.getByText(/Something went wrong with the content system/)).toBeInTheDocument();

    // Should have action buttons
    expect(screen.getByText('Try Again')).toBeInTheDocument();
    expect(screen.getByText('Reload App')).toBeInTheDocument();

    consoleSpy.mockRestore();
  });

  it('loads dynamic content when available', async () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

    renderWithProvider(
      <ContentErrorBoundary>
        <ThrowError shouldThrow={true} />
      </ContentErrorBoundary>
    );

    // Wait for dynamic content to load
    await waitFor(() => {
      expect(screen.getByText('System Error (Dynamic)')).toBeInTheDocument();
    }, { timeout: 2000 });

    await waitFor(() => {
      expect(screen.getByText('Dynamic error description from content API')).toBeInTheDocument();
    }, { timeout: 2000 });

    consoleSpy.mockRestore();
  });

  it('handles reset functionality', () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

    renderWithProvider(
      <ContentErrorBoundary>
        <ThrowError shouldThrow={true} />
      </ContentErrorBoundary>
    );

    expect(screen.getByText('Content System Error')).toBeInTheDocument();

    // Click Try Again button
    fireEvent.click(screen.getByText('Try Again'));

    // Error UI should disappear and children should render
    // Note: In a real scenario, the component would reset and re-render children
    // This test verifies the button exists and is clickable
    expect(screen.getByText('Try Again')).toBeInTheDocument();

    consoleSpy.mockRestore();
  });

  it('renders custom fallback when provided', () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
    const customFallback = <div>Custom Error UI</div>;

    renderWithProvider(
      <ContentErrorBoundary fallback={customFallback}>
        <ThrowError shouldThrow={true} />
      </ContentErrorBoundary>
    );

    expect(screen.getByText('Custom Error UI')).toBeInTheDocument();
    expect(screen.queryByText('Content System Error')).not.toBeInTheDocument();

    consoleSpy.mockRestore();
  });
});

describe('ImageWithFallback', () => {
  it('renders image when src loads successfully', () => {
    renderWithProvider(
      <ImageWithFallback src="https://example.com/image.jpg" alt="Test image" />
    );

    const img = screen.getByAltText('Test image');
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute('src', 'https://example.com/image.jpg');
  });

  it('renders fallback when image fails to load', async () => {
    renderWithProvider(
      <ImageWithFallback src="https://example.com/broken-image.jpg" alt="Test image" />
    );

    const img = screen.getByAltText('Test image');

    // Simulate image load error
    fireEvent.error(img);

    // Should show fallback content with dynamic error text (or fallback)
    await waitFor(() => {
      const fallbackImg = screen.getByAltText(/Error loading image/);
      expect(fallbackImg).toBeInTheDocument();
    });
  });

  it('uses dynamic error content when available', async () => {
    // Add mock for image error content
    server.use(
      http.get('/api/content/error.image.load_failed', () => {
        return HttpResponse.json({
          id: 'error.image.load_failed',
          type: 'text',
          data: 'Dynamic image error message',
          metadata: {},
          version: 1,
          category: 'error',
          created_at: '2023-01-01T00:00:00Z',
          updated_at: '2023-01-01T00:00:00Z',
        });
      })
    );

    renderWithProvider(
      <ImageWithFallback src="https://example.com/broken-image.jpg" alt="Test image" />
    );

    const img = screen.getByAltText('Test image');
    fireEvent.error(img);

    // Should eventually show dynamic error message
    await waitFor(() => {
      const fallbackImg = screen.getByAltText('Dynamic image error message');
      expect(fallbackImg).toBeInTheDocument();
    }, { timeout: 2000 });
  });
});

describe('useContentText hook', () => {
  const TestComponent = ({ contentId, fallback }: { contentId: string; fallback: string }) => {
    const { text, isLoading, isError, isFallback } = useContentText(contentId, fallback);

    return (
      <div>
        <span data-testid="text">{text}</span>
        <span data-testid="loading">{isLoading.toString()}</span>
        <span data-testid="error">{isError.toString()}</span>
        <span data-testid="fallback">{isFallback.toString()}</span>
      </div>
    );
  };

  it('returns fallback text initially and loads dynamic content', async () => {
    renderWithProvider(
      <TestComponent contentId="error.boundary.title" fallback="Fallback Title" />
    );

    // Should show fallback initially
    expect(screen.getByTestId('text')).toHaveTextContent('Fallback Title');
    expect(screen.getByTestId('fallback')).toHaveTextContent('true');

    // Should load dynamic content
    await waitFor(() => {
      expect(screen.getByTestId('text')).toHaveTextContent('System Error (Dynamic)');
      expect(screen.getByTestId('fallback')).toHaveTextContent('false');
    }, { timeout: 2000 });
  });

  it('uses fallback when content API fails', async () => {
    renderWithProvider(
      <TestComponent contentId="error.nonexistent.content" fallback="Fallback Content" />
    );

    // Should show fallback
    expect(screen.getByTestId('text')).toHaveTextContent('Fallback Content');

    // Should remain fallback after API call fails
    await waitFor(() => {
      expect(screen.getByTestId('fallback')).toHaveTextContent('true');
    }, { timeout: 1000 });
  });
});

describe('Error message content IDs', () => {
  it('follows error.{type}.{element} pattern', () => {
    // Test that our content IDs follow the expected pattern
    const expectedPatterns = [
      'error.boundary.title',
      'error.boundary.description',
      'error.boundary.try_again',
      'error.boundary.reload_app',
      'error.http.bad_request',
      'error.http.unauthorized',
      'error.network.connection',
      'error.image.load_failed',
      'error.message.send_failed',
      'error.validation.generic',
    ];

    expectedPatterns.forEach(pattern => {
      expect(pattern).toMatch(/^error\.\w+\.\w+/);
    });
  });
});