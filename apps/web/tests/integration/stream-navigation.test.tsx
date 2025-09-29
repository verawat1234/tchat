import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import StreamTabs from '../../src/components/store/StreamTabs';
import { streamSlice } from '../../src/store/slices/streamSlice';

// Mock the Stream API
vi.mock('../../src/services/streamApi', () => ({
  streamApi: {
    useGetStreamCategoriesQuery: vi.fn(),
    useGetStreamContentQuery: vi.fn(),
    useGetStreamFeaturedQuery: vi.fn(),
  },
}));

describe('Stream Tab Navigation Integration Tests', () => {
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
   * T009: Frontend Integration Test - Stream Tab Navigation
   * This test MUST FAIL until Stream tab navigation is implemented
   */
  it('should render all 6 Stream category tabs', async () => {
    // Mock API response for categories
    const mockCategories = [
      { id: 'books', name: 'Books', displayOrder: 1, iconName: 'book-open', isActive: true },
      { id: 'podcasts', name: 'Podcasts', displayOrder: 2, iconName: 'microphone', isActive: true },
      { id: 'cartoons', name: 'Cartoons', displayOrder: 3, iconName: 'film', isActive: true },
      { id: 'movies', name: 'Movies', displayOrder: 4, iconName: 'video', isActive: true },
      { id: 'music', name: 'Music', displayOrder: 5, iconName: 'music', isActive: true },
      { id: 'art', name: 'Art', displayOrder: 6, iconName: 'palette', isActive: true },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamCategoriesQuery as any).mockReturnValue({
      data: { categories: mockCategories, total: 6 },
      isLoading: false,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamTabs />
      </Provider>
    );

    // This should FAIL until StreamTabs component is implemented
    expect(screen.getByRole('tablist')).toBeInTheDocument();

    // All 6 category tabs should be present
    expect(screen.getByRole('tab', { name: /books/i })).toBeInTheDocument();
    expect(screen.getByRole('tab', { name: /podcasts/i })).toBeInTheDocument();
    expect(screen.getByRole('tab', { name: /cartoons/i })).toBeInTheDocument();
    expect(screen.getByRole('tab', { name: /movies/i })).toBeInTheDocument();
    expect(screen.getByRole('tab', { name: /music/i })).toBeInTheDocument();
    expect(screen.getByRole('tab', { name: /art/i })).toBeInTheDocument();
  });

  it('should switch between tabs with smooth transitions', async () => {
    // This test MUST FAIL until tab switching is implemented
    const mockCategories = [
      { id: 'books', name: 'Books', displayOrder: 1, iconName: 'book-open', isActive: true },
      { id: 'podcasts', name: 'Podcasts', displayOrder: 2, iconName: 'microphone', isActive: true },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamCategoriesQuery as any).mockReturnValue({
      data: { categories: mockCategories, total: 2 },
      isLoading: false,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamTabs />
      </Provider>
    );

    const booksTab = screen.getByRole('tab', { name: /books/i });
    const podcastsTab = screen.getByRole('tab', { name: /podcasts/i });

    // Initial state - books should be selected
    expect(booksTab).toHaveAttribute('aria-selected', 'true');

    // Click podcasts tab
    fireEvent.click(podcastsTab);

    // Should transition to podcasts tab
    await waitFor(() => {
      expect(podcastsTab).toHaveAttribute('aria-selected', 'true');
      expect(booksTab).toHaveAttribute('aria-selected', 'false');
    });
  });

  it('should load content when tab is selected', async () => {
    // This test MUST FAIL until content loading is implemented
    const mockCategories = [
      { id: 'books', name: 'Books', displayOrder: 1, iconName: 'book-open', isActive: true },
    ];

    const mockContent = [
      {
        id: 'book-1',
        categoryId: 'books',
        title: 'Test Book',
        description: 'A test book',
        thumbnailUrl: 'test.jpg',
        contentType: 'book',
        price: 9.99,
        currency: 'USD',
        availabilityStatus: 'available',
        isFeatured: false,
      },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamCategoriesQuery as any).mockReturnValue({
      data: { categories: mockCategories, total: 1 },
      isLoading: false,
      error: null,
    });

    (streamApi.useGetStreamContentQuery as any).mockReturnValue({
      data: { items: mockContent, total: 1, hasMore: false },
      isLoading: false,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamTabs />
      </Provider>
    );

    // Content should load for the active tab
    await waitFor(() => {
      expect(screen.getByText('Test Book')).toBeInTheDocument();
    });
  });

  it('should persist tab selection state', async () => {
    // This test MUST FAIL until state persistence is implemented
    const mockCategories = [
      { id: 'books', name: 'Books', displayOrder: 1, iconName: 'book-open', isActive: true },
      { id: 'podcasts', name: 'Podcasts', displayOrder: 2, iconName: 'microphone', isActive: true },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamCategoriesQuery as any).mockReturnValue({
      data: { categories: mockCategories, total: 2 },
      isLoading: false,
      error: null,
    });

    const { rerender } = render(
      <Provider store={store}>
        <StreamTabs />
      </Provider>
    );

    // Select podcasts tab
    const podcastsTab = screen.getByRole('tab', { name: /podcasts/i });
    fireEvent.click(podcastsTab);

    // Rerender component
    rerender(
      <Provider store={store}>
        <StreamTabs />
      </Provider>
    );

    // Tab selection should persist
    await waitFor(() => {
      expect(podcastsTab).toHaveAttribute('aria-selected', 'true');
    });
  });

  it('should meet performance requirements (<200ms transitions)', async () => {
    // This test MUST FAIL until performance requirements are met
    const mockCategories = [
      { id: 'books', name: 'Books', displayOrder: 1, iconName: 'book-open', isActive: true },
      { id: 'podcasts', name: 'Podcasts', displayOrder: 2, iconName: 'microphone', isActive: true },
    ];

    const { streamApi } = await import('../../src/services/streamApi');
    (streamApi.useGetStreamCategoriesQuery as any).mockReturnValue({
      data: { categories: mockCategories, total: 2 },
      isLoading: false,
      error: null,
    });

    render(
      <Provider store={store}>
        <StreamTabs />
      </Provider>
    );

    const podcastsTab = screen.getByRole('tab', { name: /podcasts/i });

    // Measure transition time
    const startTime = performance.now();
    fireEvent.click(podcastsTab);

    await waitFor(() => {
      expect(podcastsTab).toHaveAttribute('aria-selected', 'true');
    });

    const endTime = performance.now();
    const transitionTime = endTime - startTime;

    // Should transition in less than 200ms
    expect(transitionTime).toBeLessThan(200);
  });
});