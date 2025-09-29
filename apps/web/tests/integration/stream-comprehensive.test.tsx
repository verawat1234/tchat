import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react'
import { Provider } from 'react-redux'
import { configureStore } from '@reduxjs/toolkit'
import { setupListeners } from '@reduxjs/toolkit/query'
import { rest } from 'msw'
import { setupServer } from 'msw/node'
import userEvent from '@testing-library/user-event'

import { api } from '../../src/services/api'
import { StreamTab } from '../../src/components/StreamTab'
import { streamApi } from '../../src/services/streamApi'

// Test data
const mockCategories = [
  {
    id: 'books',
    name: 'Books',
    displayOrder: 1,
    iconName: 'book-open',
    isActive: true,
    featuredContentEnabled: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  },
  {
    id: 'podcasts',
    name: 'Podcasts',
    displayOrder: 2,
    iconName: 'headphones',
    isActive: true,
    featuredContentEnabled: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  },
  {
    id: 'movies',
    name: 'Movies',
    displayOrder: 4,
    iconName: 'video',
    isActive: true,
    featuredContentEnabled: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  }
]

const mockMovieSubtabs = [
  {
    id: 'movies_short',
    categoryId: 'movies',
    name: 'Short Films',
    displayOrder: 1,
    filterCriteria: { content_type: 'SHORT_MOVIE', max_duration: 1800 },
    isActive: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  },
  {
    id: 'movies_feature',
    categoryId: 'movies',
    name: 'Feature Films',
    displayOrder: 2,
    filterCriteria: { content_type: 'LONG_MOVIE', min_duration: 1800 },
    isActive: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  }
]

const mockStreamContent = [
  {
    id: 'book-1',
    categoryId: 'books',
    title: 'Test Book 1',
    description: 'A fascinating test book',
    thumbnailUrl: 'https://example.com/book1.jpg',
    contentType: 'book' as const,
    price: 9.99,
    currency: 'USD',
    availabilityStatus: 'available' as const,
    isFeatured: true,
    featuredOrder: 1,
    metadata: { author: 'Test Author', genre: 'fiction' },
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  },
  {
    id: 'movie-short-1',
    categoryId: 'movies',
    title: 'Test Short Film',
    description: 'A compelling short film',
    thumbnailUrl: 'https://example.com/short1.jpg',
    contentType: 'short_movie' as const,
    duration: 1200,
    price: 2.99,
    currency: 'USD',
    availabilityStatus: 'available' as const,
    isFeatured: true,
    featuredOrder: 1,
    metadata: { director: 'Test Director', year: 2024 },
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  },
  {
    id: 'movie-feature-1',
    categoryId: 'movies',
    title: 'Test Feature Film',
    description: 'An epic feature film',
    thumbnailUrl: 'https://example.com/feature1.jpg',
    contentType: 'long_movie' as const,
    duration: 7200,
    price: 12.99,
    currency: 'USD',
    availabilityStatus: 'available' as const,
    isFeatured: false,
    metadata: { director: 'Test Director', year: 2024, rating: 'PG-13' },
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  }
]

const mockFeaturedContent = mockStreamContent.filter(item => item.isFeatured)

// MSW server setup
const server = setupServer(
  // Get categories
  rest.get('/api/v1/stream/categories', (req, res, ctx) => {
    return res(
      ctx.json({
        categories: mockCategories,
        total: mockCategories.length,
        success: true
      })
    )
  }),

  // Get category detail
  rest.get('/api/v1/stream/categories/:categoryId', (req, res, ctx) => {
    const { categoryId } = req.params
    const category = mockCategories.find(c => c.id === categoryId)

    if (!category) {
      return res(ctx.status(404), ctx.json({ success: false, error: 'Category not found' }))
    }

    const subtabs = categoryId === 'movies' ? mockMovieSubtabs : []

    return res(
      ctx.json({
        category,
        subtabs,
        stats: { totalContent: mockStreamContent.filter(c => c.categoryId === categoryId).length },
        success: true
      })
    )
  }),

  // Get content by category
  rest.get('/api/v1/stream/content', (req, res, ctx) => {
    const categoryId = req.url.searchParams.get('categoryId')
    const subtabId = req.url.searchParams.get('subtabId')
    const page = parseInt(req.url.searchParams.get('page') || '1')
    const limit = parseInt(req.url.searchParams.get('limit') || '20')

    if (!categoryId) {
      return res(ctx.status(400), ctx.json({ success: false, error: 'Category ID required' }))
    }

    let filteredContent = mockStreamContent.filter(item => item.categoryId === categoryId)

    // Apply subtab filtering
    if (subtabId === 'movies_short') {
      filteredContent = filteredContent.filter(item =>
        item.contentType === 'short_movie' && (item.duration || 0) <= 1800
      )
    } else if (subtabId === 'movies_feature') {
      filteredContent = filteredContent.filter(item =>
        item.contentType === 'long_movie' && (item.duration || 0) > 1800
      )
    }

    const startIndex = (page - 1) * limit
    const endIndex = startIndex + limit
    const paginatedContent = filteredContent.slice(startIndex, endIndex)

    return res(
      ctx.json({
        items: paginatedContent,
        total: filteredContent.length,
        hasMore: endIndex < filteredContent.length,
        success: true
      })
    )
  }),

  // Get featured content
  rest.get('/api/v1/stream/featured', (req, res, ctx) => {
    const categoryId = req.url.searchParams.get('categoryId')
    const limit = parseInt(req.url.searchParams.get('limit') || '10')

    let featured = mockFeaturedContent
    if (categoryId) {
      featured = featured.filter(item => item.categoryId === categoryId)
    }

    featured = featured.slice(0, limit)

    return res(
      ctx.json({
        items: featured,
        total: featured.length,
        hasMore: false,
        success: true
      })
    )
  }),

  // Search content
  rest.get('/api/v1/stream/search', (req, res, ctx) => {
    const query = req.url.searchParams.get('q')
    const categoryId = req.url.searchParams.get('categoryId')
    const page = parseInt(req.url.searchParams.get('page') || '1')
    const limit = parseInt(req.url.searchParams.get('limit') || '20')

    if (!query) {
      return res(ctx.status(400), ctx.json({ success: false, error: 'Search query required' }))
    }

    let results = mockStreamContent.filter(item =>
      item.title.toLowerCase().includes(query.toLowerCase()) ||
      item.description.toLowerCase().includes(query.toLowerCase())
    )

    if (categoryId) {
      results = results.filter(item => item.categoryId === categoryId)
    }

    const startIndex = (page - 1) * limit
    const endIndex = startIndex + limit
    const paginatedResults = results.slice(startIndex, endIndex)

    return res(
      ctx.json({
        items: paginatedResults,
        total: results.length,
        hasMore: endIndex < results.length,
        success: true
      })
    )
  }),

  // Purchase content
  rest.post('/api/v1/stream/content/purchase', (req, res, ctx) => {
    return res(
      ctx.json({
        orderId: 'test-order-123',
        totalAmount: 9.99,
        currency: 'USD',
        success: true,
        message: 'Purchase successful!'
      })
    )
  })
)

// Store setup helper
function createTestStore() {
  return configureStore({
    reducer: {
      api: api.reducer
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(api.middleware)
  })
}

// Test wrapper component
function TestWrapper({ children }: { children: React.ReactNode }) {
  const store = createTestStore()
  setupListeners(store.dispatch)

  return (
    <Provider store={store}>
      {children}
    </Provider>
  )
}

// Test helper to render StreamTab
function renderStreamTab(props = {}) {
  const defaultProps = {
    user: { id: 'test-user', name: 'Test User' },
    onContentClick: vi.fn(),
    onAddToCart: vi.fn(),
    onContentShare: vi.fn(),
    cartItems: []
  }

  return render(
    <TestWrapper>
      <StreamTab {...defaultProps} {...props} />
    </TestWrapper>
  )
}

describe('Stream Store Tabs - Comprehensive Integration Tests', () => {
  beforeEach(() => {
    server.listen()
  })

  afterEach(() => {
    server.resetHandlers()
    vi.clearAllMocks()
  })

  // FLOW TEST 1: Complete Navigation Flow
  describe('Complete Navigation Flow', () => {
    it('should load categories and navigate between tabs', async () => {
      renderStreamTab()

      // Wait for categories to load
      await waitFor(() => {
        expect(screen.getByText('Books')).toBeInTheDocument()
        expect(screen.getByText('Podcasts')).toBeInTheDocument()
        expect(screen.getByText('Movies')).toBeInTheDocument()
      })

      // Default should be Books tab
      const booksTab = screen.getByRole('tab', { name: /books/i })
      expect(booksTab).toHaveAttribute('aria-selected', 'true')

      // Click Podcasts tab
      const podcastsTab = screen.getByRole('tab', { name: /podcasts/i })
      await userEvent.click(podcastsTab)

      await waitFor(() => {
        expect(podcastsTab).toHaveAttribute('aria-selected', 'true')
        expect(booksTab).toHaveAttribute('aria-selected', 'false')
      })

      // Click Movies tab
      const moviesTab = screen.getByRole('tab', { name: /movies/i })
      await userEvent.click(moviesTab)

      await waitFor(() => {
        expect(moviesTab).toHaveAttribute('aria-selected', 'true')
        expect(podcastsTab).toHaveAttribute('aria-selected', 'false')
      })
    })

    it('should maintain tab state during navigation', async () => {
      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Movies')).toBeInTheDocument()
      })

      // Navigate to Movies and check content loads
      const moviesTab = screen.getByRole('tab', { name: /movies/i })
      await userEvent.click(moviesTab)

      await waitFor(() => {
        expect(screen.getByText('Test Short Film')).toBeInTheDocument()
        expect(screen.getByText('Test Feature Film')).toBeInTheDocument()
      })

      // Navigate away and back
      const booksTab = screen.getByRole('tab', { name: /books/i })
      await userEvent.click(booksTab)

      await waitFor(() => {
        expect(screen.getByText('Test Book 1')).toBeInTheDocument()
      })

      // Return to Movies - should remember state
      await userEvent.click(moviesTab)

      await waitFor(() => {
        expect(moviesTab).toHaveAttribute('aria-selected', 'true')
        expect(screen.getByText('Test Short Film')).toBeInTheDocument()
      })
    })
  })

  // FLOW TEST 2: Featured Content Discovery Flow
  describe('Featured Content Discovery Flow', () => {
    it('should display featured content carousel', async () => {
      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Books')).toBeInTheDocument()
      })

      // Should show featured content for books
      await waitFor(() => {
        expect(screen.getByText('Featured')).toBeInTheDocument()
        expect(screen.getByText('Test Book 1')).toBeInTheDocument()
      })

      // Navigate to movies and check featured content
      const moviesTab = screen.getByRole('tab', { name: /movies/i })
      await userEvent.click(moviesTab)

      await waitFor(() => {
        expect(screen.getByText('Test Short Film')).toBeInTheDocument()
      })
    })

    it('should handle featured content interactions', async () => {
      const onContentClick = vi.fn()
      renderStreamTab({ onContentClick })

      await waitFor(() => {
        expect(screen.getByText('Test Book 1')).toBeInTheDocument()
      })

      // Click on featured content
      const featuredItem = screen.getByText('Test Book 1')
      await userEvent.click(featuredItem)

      expect(onContentClick).toHaveBeenCalledWith('book-1')
    })
  })

  // FLOW TEST 3: Movies Subtab Navigation Flow
  describe('Movies Subtab Navigation Flow', () => {
    it('should display and navigate subtabs in movies category', async () => {
      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Movies')).toBeInTheDocument()
      })

      // Navigate to Movies tab
      const moviesTab = screen.getByRole('tab', { name: /movies/i })
      await userEvent.click(moviesTab)

      // Wait for subtabs to appear
      await waitFor(() => {
        expect(screen.getByText('Short Films')).toBeInTheDocument()
        expect(screen.getByText('Feature Films')).toBeInTheDocument()
      })

      // Should show both short and feature films initially
      expect(screen.getByText('Test Short Film')).toBeInTheDocument()
      expect(screen.getByText('Test Feature Film')).toBeInTheDocument()

      // Click Short Films subtab
      const shortFilmsTab = screen.getByText('Short Films')
      await userEvent.click(shortFilmsTab)

      await waitFor(() => {
        expect(screen.getByText('Test Short Film')).toBeInTheDocument()
        // Feature film should be filtered out
        expect(screen.queryByText('Test Feature Film')).not.toBeInTheDocument()
      })

      // Click Feature Films subtab
      const featureFilmsTab = screen.getByText('Feature Films')
      await userEvent.click(featureFilmsTab)

      await waitFor(() => {
        expect(screen.getByText('Test Feature Film')).toBeInTheDocument()
        // Short film should be filtered out
        expect(screen.queryByText('Test Short Film')).not.toBeInTheDocument()
      })
    })
  })

  // FLOW TEST 4: Content Interaction Flow
  describe('Content Interaction Flow', () => {
    it('should handle add to cart interactions', async () => {
      const onAddToCart = vi.fn()
      renderStreamTab({ onAddToCart })

      await waitFor(() => {
        expect(screen.getByText('Test Book 1')).toBeInTheDocument()
      })

      // Find and click add to cart button
      const addToCartButtons = screen.getAllByText('Add to Cart')
      expect(addToCartButtons.length).toBeGreaterThan(0)

      await userEvent.click(addToCartButtons[0])

      expect(onAddToCart).toHaveBeenCalledWith('book-1', 1)
    })

    it('should display content details correctly', async () => {
      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Test Book 1')).toBeInTheDocument()
      })

      // Check content details are displayed
      expect(screen.getByText('A fascinating test book')).toBeInTheDocument()
      expect(screen.getByText('$9.99')).toBeInTheDocument()

      // Navigate to movies
      const moviesTab = screen.getByRole('tab', { name: /movies/i })
      await userEvent.click(moviesTab)

      await waitFor(() => {
        expect(screen.getByText('Test Short Film')).toBeInTheDocument()
        expect(screen.getByText('20m')).toBeInTheDocument() // Duration formatted
        expect(screen.getByText('$2.99')).toBeInTheDocument()
      })
    })
  })

  // FLOW TEST 5: Search Flow
  describe('Search Flow', () => {
    it('should handle content search', async () => {
      // Note: This would require implementing search UI component
      // For now, we test the API endpoint through direct calls
      const store = createTestStore()

      const result = await store.dispatch(
        streamApi.endpoints.searchStreamContent.initiate({
          q: 'test',
          page: 1,
          limit: 10
        })
      )

      expect(result.data?.success).toBe(true)
      expect(result.data?.items.length).toBeGreaterThan(0)
      expect(result.data?.items[0].title).toContain('Test')
    })
  })

  // FLOW TEST 6: Error Handling Flow
  describe('Error Handling Flow', () => {
    it('should handle API errors gracefully', async () => {
      // Mock server error
      server.use(
        rest.get('/api/v1/stream/categories', (req, res, ctx) => {
          return res(ctx.status(500), ctx.json({ success: false, error: 'Server error' }))
        })
      )

      renderStreamTab()

      // Should show loading state and then error state
      await waitFor(() => {
        expect(screen.getByText(/No content available/i)).toBeInTheDocument()
      })
    })

    it('should handle empty content gracefully', async () => {
      // Mock empty response
      server.use(
        rest.get('/api/v1/stream/content', (req, res, ctx) => {
          return res(
            ctx.json({
              items: [],
              total: 0,
              hasMore: false,
              success: true
            })
          )
        })
      )

      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Books')).toBeInTheDocument()
      })

      // Navigate to books to trigger content load
      const booksTab = screen.getByRole('tab', { name: /books/i })
      await userEvent.click(booksTab)

      await waitFor(() => {
        expect(screen.getByText(/No content available/i)).toBeInTheDocument()
      })
    })
  })

  // FLOW TEST 7: Performance Flow
  describe('Performance Flow', () => {
    it('should load content within performance targets', async () => {
      const startTime = performance.now()

      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Test Book 1')).toBeInTheDocument()
      })

      const loadTime = performance.now() - startTime

      // Should load within 3 seconds (3000ms) as per requirements
      expect(loadTime).toBeLessThan(3000)
    })

    it('should handle rapid tab switching efficiently', async () => {
      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Books')).toBeInTheDocument()
      })

      const tabs = ['podcasts', 'movies', 'books']

      for (const tabName of tabs) {
        const tab = screen.getByRole('tab', { name: new RegExp(tabName, 'i') })
        await userEvent.click(tab)

        // Should respond quickly to tab switches
        await waitFor(() => {
          expect(tab).toHaveAttribute('aria-selected', 'true')
        }, { timeout: 200 }) // 200ms max response time
      }
    })
  })

  // FLOW TEST 8: Cross-Platform Data Consistency
  describe('Cross-Platform Data Consistency', () => {
    it('should maintain consistent data structure across API calls', async () => {
      const store = createTestStore()

      // Test categories consistency
      const categoriesResult = await store.dispatch(
        streamApi.endpoints.getStreamCategories.initiate()
      )

      expect(categoriesResult.data?.categories).toBeDefined()
      expect(categoriesResult.data?.categories[0]).toHaveProperty('id')
      expect(categoriesResult.data?.categories[0]).toHaveProperty('name')
      expect(categoriesResult.data?.categories[0]).toHaveProperty('displayOrder')

      // Test content consistency
      const contentResult = await store.dispatch(
        streamApi.endpoints.getStreamContent.initiate({
          categoryId: 'books',
          page: 1,
          limit: 10
        })
      )

      expect(contentResult.data?.items).toBeDefined()
      expect(contentResult.data?.items[0]).toHaveProperty('id')
      expect(contentResult.data?.items[0]).toHaveProperty('title')
      expect(contentResult.data?.items[0]).toHaveProperty('price')
      expect(contentResult.data?.items[0]).toHaveProperty('categoryId')
    })
  })

  // FLOW TEST 9: Accessibility Flow
  describe('Accessibility Flow', () => {
    it('should support keyboard navigation', async () => {
      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Books')).toBeInTheDocument()
      })

      const booksTab = screen.getByRole('tab', { name: /books/i })

      // Focus the tab
      booksTab.focus()
      expect(booksTab).toHaveFocus()

      // Tab navigation
      await userEvent.keyboard('{ArrowRight}')

      const podcastsTab = screen.getByRole('tab', { name: /podcasts/i })
      expect(podcastsTab).toHaveFocus()

      // Activate tab with Enter
      await userEvent.keyboard('{Enter}')

      await waitFor(() => {
        expect(podcastsTab).toHaveAttribute('aria-selected', 'true')
      })
    })

    it('should have proper ARIA labels and roles', async () => {
      renderStreamTab()

      await waitFor(() => {
        expect(screen.getByText('Books')).toBeInTheDocument()
      })

      // Check tab roles
      const tabs = screen.getAllByRole('tab')
      expect(tabs.length).toBeGreaterThan(0)

      // Check tabpanel exists
      const tabpanel = screen.getByRole('tabpanel')
      expect(tabpanel).toBeInTheDocument()

      // Check for tablist
      const tablist = screen.getByRole('tablist')
      expect(tablist).toBeInTheDocument()
    })
  })

  // FLOW TEST 10: State Management Flow
  describe('State Management Flow', () => {
    it('should maintain state consistency in Redux store', async () => {
      const store = createTestStore()

      // Dispatch category query
      const result = await store.dispatch(
        streamApi.endpoints.getStreamCategories.initiate()
      )

      // Check that data is cached in store
      const state = store.getState()
      const apiState = state.api

      // Should have cached the result
      expect(apiState.queries).toBeDefined()

      // Dispatch same query again - should use cache
      const cachedResult = await store.dispatch(
        streamApi.endpoints.getStreamCategories.initiate()
      )

      expect(cachedResult.data).toEqual(result.data)
    })
  })
})