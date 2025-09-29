import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import StreamCategoryContent from '../../src/components/store/StreamCategoryContent';
import StreamFeaturedCarousel from '../../src/components/store/StreamFeaturedCarousel';
import { streamSlice } from '../../src/store/slices/streamSlice';

// Mock the Stream API
vi.mock('../../src/services/streamApi', () => ({
  streamApi: {
    useGetStreamContentQuery: vi.fn(),
    useGetStreamFeaturedQuery: vi.fn(),
  },
}));

describe('Stream Content Display Integration Tests', () => {
  let store: any;

  beforeEach(() => {
    // Create a fresh store for each test
    store = configureStore({
      reducer: {
        stream: streamSlice.reducer,
      },
    });
  });

  /**
   * T010: Frontend Integration Test - Stream Content Display
   * This test MUST FAIL until Stream content display is implemented
   */
  it('should display featured content carousel', async () => {
    // Mock featured content API response
    const mockFeaturedContent = [
      {
        id: 'featured-1',
        categoryId: 'books',
        title: 'Featured Book 1',
        description: 'A featured book',
        thumbnailUrl: 'featured1.jpg',
        contentType: 'book',
        price: 19.99,
        currency: 'USD',
        availabilityStatus: 'available',
        isFeatured: true,
        featuredOrder: 1,
      },
      {
        id: 'featured-2',
        categoryId: 'books',
        title: 'Featured Book 2',
        description: 'Another featured book',
        thumbnailUrl: 'featured2.jpg',
        contentType: 'book',
        price: 24.99,
        currency: 'USD',
        availabilityStatus: 'available',
        isFeatured: true,
        featuredOrder: 2,
      },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamFeaturedQuery as any).mockReturnValue({
      data: { items: mockFeaturedContent, total: 2, hasMore: false },
      isLoading: false,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamFeaturedCarousel categoryId="books" />
      </Provider>
    );

    // This should FAIL until StreamFeaturedCarousel component is implemented
    expect(screen.getByRole('region', { name: /featured content/i })).toBeInTheDocument();

    // Featured items should be displayed
    expect(screen.getByText('Featured Book 1')).toBeInTheDocument();
    expect(screen.getByText('Featured Book 2')).toBeInTheDocument();

    // Should display prices
    expect(screen.getByText('$19.99')).toBeInTheDocument();
    expect(screen.getByText('$24.99')).toBeInTheDocument();
  });

  it('should display browseable content in grid layout', async () => {
    // Mock content API response
    const mockContent = [
      {
        id: 'content-1',
        categoryId: 'books',
        title: 'Regular Book 1',
        description: 'A regular book',
        thumbnailUrl: 'book1.jpg',
        contentType: 'book',
        price: 9.99,
        currency: 'USD',
        availabilityStatus: 'available',
        isFeatured: false,
      },
      {
        id: 'content-2',
        categoryId: 'books',
        title: 'Regular Book 2',
        description: 'Another regular book',
        thumbnailUrl: 'book2.jpg',
        contentType: 'book',
        price: 14.99,
        currency: 'USD',
        availabilityStatus: 'available',
        isFeatured: false,
      },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamContentQuery as any).mockReturnValue({
      data: { items: mockContent, total: 2, hasMore: false },
      isLoading: false,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamCategoryContent categoryId="books" />
      </Provider>
    );

    // This should FAIL until StreamCategoryContent component is implemented
    expect(screen.getByRole('region', { name: /content grid/i })).toBeInTheDocument();

    // Content items should be displayed
    expect(screen.getByText('Regular Book 1')).toBeInTheDocument();
    expect(screen.getByText('Regular Book 2')).toBeInTheDocument();
  });

  it('should support horizontal scrolling in featured carousel', async () => {
    // Mock multiple featured items for scrolling
    const mockFeaturedContent = Array.from({ length: 10 }, (_, i) => ({
      id: `featured-${i + 1}`,
      categoryId: 'books',
      title: `Featured Book ${i + 1}`,
      description: `Featured book ${i + 1}`,
      thumbnailUrl: `featured${i + 1}.jpg`,
      contentType: 'book',
      price: 19.99 + i,
      currency: 'USD',
      availabilityStatus: 'available',
      isFeatured: true,
      featuredOrder: i + 1,
    }));

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamFeaturedQuery as any).mockReturnValue({
      data: { items: mockFeaturedContent, total: 10, hasMore: false },
      isLoading: false,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamFeaturedCarousel categoryId="books" />
      </Provider>
    );

    // Should have horizontal scroll functionality
    const carousel = screen.getByRole('region', { name: /featured content/i });

    // Check for scroll indicators or controls
    expect(carousel).toHaveStyle({ overflowX: 'auto' });

    // Should display first few items
    expect(screen.getByText('Featured Book 1')).toBeInTheDocument();
    expect(screen.getByText('Featured Book 2')).toBeInTheDocument();
  });

  it('should handle responsive design for different screen sizes', async () => {
    const mockContent = [
      {
        id: 'content-1',
        categoryId: 'books',
        title: 'Responsive Book',
        description: 'A book for responsive testing',
        thumbnailUrl: 'responsive.jpg',
        contentType: 'book',
        price: 12.99,
        currency: 'USD',
        availabilityStatus: 'available',
        isFeatured: false,
      },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamContentQuery as any).mockReturnValue({
      data: { items: mockContent, total: 1, hasMore: false },
      isLoading: false,
      error: null,
    });

    // Test mobile viewport
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 375,
    });

    render(
      <Provider store={store}>
        <StreamCategoryContent categoryId="books" />
      </Provider>
    );

    // Should adapt to mobile layout
    const contentGrid = screen.getByRole('region', { name: /content grid/i });

    // Should have responsive grid classes
    expect(contentGrid).toHaveClass(/grid/);
    expect(contentGrid).toHaveClass(/responsive/);
  });

  it('should display loading states while content is fetching', async () => {
    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamContentQuery as any).mockReturnValue({
      data: null,
      isLoading: true,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamCategoryContent categoryId="books" />
      </Provider>
    );

    // Should show loading indicator
    expect(screen.getByRole('status', { name: /loading/i })).toBeInTheDocument();

    // Or loading skeleton
    expect(screen.getByTestId('content-loading-skeleton')).toBeInTheDocument();
  });

  it('should handle error states gracefully', async () => {
    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamContentQuery as any).mockReturnValue({
      data: null,
      isLoading: false,
      error: { message: 'Failed to load content' },
    });

    render(
      <Provider store={store}>
        <StreamCategoryContent categoryId="books" />
      </Provider>
    );

    // Should show error message
    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText(/failed to load content/i)).toBeInTheDocument();

    // Should provide retry option
    expect(screen.getByRole('button', { name: /retry/i })).toBeInTheDocument();
  });

  it('should meet Core Web Vitals requirements', async () => {
    // This test MUST FAIL until performance requirements are met
    const mockContent = [
      {
        id: 'content-1',
        categoryId: 'books',
        title: 'Performance Book',
        description: 'A book for performance testing',
        thumbnailUrl: 'performance.jpg',
        contentType: 'book',
        price: 15.99,
        currency: 'USD',
        availabilityStatus: 'available',
        isFeatured: false,
      },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamContentQuery as any).mockReturnValue({
      data: { items: mockContent, total: 1, hasMore: false },
      isLoading: false,
      error: null,
    });

    const startTime = performance.now();

    render(
      <Provider store={store}>
        <StreamCategoryContent categoryId="books" />
      </Provider>
    );

    await waitFor(() => {
      expect(screen.getByText('Performance Book')).toBeInTheDocument();
    });

    const endTime = performance.now();
    const renderTime = endTime - startTime;

    // Should render within Core Web Vitals targets
    expect(renderTime).toBeLessThan(1000); // LCP < 1s for this component
  });

  it('should support visual distinction between featured and regular content', async () => {
    const mockFeaturedContent = [
      {
        id: 'featured-1',
        categoryId: 'books',
        title: 'Featured Book',
        description: 'A featured book',
        thumbnailUrl: 'featured.jpg',
        contentType: 'book',
        price: 29.99,
        currency: 'USD',
        availabilityStatus: 'available',
        isFeatured: true,
        featuredOrder: 1,
      },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamFeaturedQuery as any).mockReturnValue({
      data: { items: mockFeaturedContent, total: 1, hasMore: false },
      isLoading: false,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamFeaturedCarousel categoryId="books" />
      </Provider>
    );

    // Featured content should have distinctive styling
    const featuredItem = screen.getByText('Featured Book').closest('[data-testid="featured-item"]');
    expect(featuredItem).toHaveClass(/featured/);
    expect(featuredItem).toHaveClass(/highlight/);
  });
});