/**
 * Offline/Online Integration Tests
 * Tests application behavior during offline periods and online reconnection
 */

import { describe, it, expect, beforeAll, afterAll, beforeEach, afterEach } from 'vitest'
import { setupServer } from 'msw/node'
import { rest } from 'msw'

interface OfflineQueueItem {
  id: string
  operation: 'CREATE' | 'UPDATE' | 'DELETE'
  endpoint: string
  data?: any
  timestamp: number
  retryCount: number
  maxRetries: number
}

interface StorageAdapter {
  getItem(key: string): string | null
  setItem(key: string, value: string): void
  removeItem(key: string): void
  clear(): void
}

interface NetworkAdapter {
  isOnline(): boolean
  onOnline(callback: () => void): void
  onOffline(callback: () => void): void
}

// Mock implementations
class MockStorageAdapter implements StorageAdapter {
  private storage: Map<string, string> = new Map()

  getItem(key: string): string | null {
    return this.storage.get(key) || null
  }

  setItem(key: string, value: string): void {
    this.storage.set(key, value)
  }

  removeItem(key: string): void {
    this.storage.delete(key)
  }

  clear(): void {
    this.storage.clear()
  }
}

class MockNetworkAdapter implements NetworkAdapter {
  private online: boolean = true
  private onlineCallbacks: (() => void)[] = []
  private offlineCallbacks: (() => void)[] = []

  isOnline(): boolean {
    return this.online
  }

  setOnline(online: boolean): void {
    const wasOnline = this.online
    this.online = online

    if (wasOnline && !online) {
      this.offlineCallbacks.forEach(callback => callback())
    } else if (!wasOnline && online) {
      this.onlineCallbacks.forEach(callback => callback())
    }
  }

  onOnline(callback: () => void): void {
    this.onlineCallbacks.push(callback)
  }

  onOffline(callback: () => void): void {
    this.offlineCallbacks.push(callback)
  }
}

// Offline-capable API client
class OfflineCapableAPIClient {
  private offlineQueue: OfflineQueueItem[] = []
  private isProcessingQueue: boolean = false

  constructor(
    private baseURL: string,
    private storage: StorageAdapter,
    private network: NetworkAdapter
  ) {
    this.loadOfflineQueue()
    this.setupNetworkListeners()
  }

  private loadOfflineQueue(): void {
    const queueData = this.storage.getItem('offline_queue')
    if (queueData) {
      this.offlineQueue = JSON.parse(queueData)
    }
  }

  private saveOfflineQueue(): void {
    this.storage.setItem('offline_queue', JSON.stringify(this.offlineQueue))
  }

  private setupNetworkListeners(): void {
    this.network.onOnline(() => {
      this.processOfflineQueue()
    })
  }

  async request<T>(
    method: string,
    endpoint: string,
    data?: any,
    options: { retry?: boolean; offlineSupport?: boolean } = {}
  ): Promise<T> {
    const { retry = true, offlineSupport = true } = options

    // Try online request first
    if (this.network.isOnline()) {
      try {
        const response = await this.makeRequest<T>(method, endpoint, data)
        return response
      } catch (error) {
        if (!retry) throw error

        // If request fails and we're supposedly online, we might be experiencing connectivity issues
        if (offlineSupport && this.shouldQueueRequest(method)) {
          return this.queueRequest<T>(method, endpoint, data)
        }
        throw error
      }
    }

    // Handle offline scenario
    if (offlineSupport) {
      if (method === 'GET') {
        return this.getCachedResponse<T>(endpoint)
      } else if (this.shouldQueueRequest(method)) {
        return this.queueRequest<T>(method, endpoint, data)
      }
    }

    throw new Error('Network unavailable and offline support not enabled')
  }

  private async makeRequest<T>(method: string, endpoint: string, data?: any): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer test-token'
      },
      body: data ? JSON.stringify(data) : undefined
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const result = await response.json()

    // Cache successful GET responses
    if (method === 'GET') {
      this.cacheResponse(endpoint, result)
    }

    return result
  }

  private shouldQueueRequest(method: string): boolean {
    return ['POST', 'PUT', 'DELETE', 'PATCH'].includes(method)
  }

  private async queueRequest<T>(method: string, endpoint: string, data?: any): Promise<T> {
    const queueItem: OfflineQueueItem = {
      id: `${Date.now()}-${Math.random()}`,
      operation: method as any,
      endpoint,
      data,
      timestamp: Date.now(),
      retryCount: 0,
      maxRetries: 3
    }

    this.offlineQueue.push(queueItem)
    this.saveOfflineQueue()

    // Return optimistic response for certain operations
    return this.getOptimisticResponse<T>(method, endpoint, data)
  }

  private getCachedResponse<T>(endpoint: string): T {
    const cached = this.storage.getItem(`cache_${endpoint}`)
    if (cached) {
      return JSON.parse(cached)
    }
    throw new Error(`No cached data available for ${endpoint}`)
  }

  private cacheResponse(endpoint: string, data: any): void {
    this.storage.setItem(`cache_${endpoint}`, JSON.stringify(data))
  }

  private getOptimisticResponse<T>(method: string, endpoint: string, data?: any): T {
    // Generate optimistic responses based on the operation
    if (method === 'POST' && endpoint.includes('/cart/items')) {
      return {
        success: true,
        cart: {
          id: 'cart-1',
          userId: 'user-1',
          items: [{ ...data, id: `temp-${Date.now()}` }],
          total: data.totalPrice,
          currency: 'USD',
          lastModified: Date.now(),
          version: 0,
          isOptimistic: true
        }
      } as T
    }

    if (method === 'PUT' && endpoint.includes('/cart/items/')) {
      return {
        success: true,
        cart: {
          id: 'cart-1',
          userId: 'user-1',
          items: [{ ...data, lastModified: Date.now() }],
          total: data.totalPrice || 0,
          currency: 'USD',
          lastModified: Date.now(),
          version: 0,
          isOptimistic: true
        }
      } as T
    }

    if (method === 'DELETE') {
      return {
        success: true,
        message: 'Item queued for deletion',
        isOptimistic: true
      } as T
    }

    return {
      success: true,
      message: 'Operation queued',
      isOptimistic: true
    } as T
  }

  private async processOfflineQueue(): Promise<void> {
    if (this.isProcessingQueue || this.offlineQueue.length === 0) {
      return
    }

    this.isProcessingQueue = true

    try {
      const queueCopy = [...this.offlineQueue]
      this.offlineQueue = []

      for (const item of queueCopy) {
        try {
          await this.makeRequest(item.operation, item.endpoint, item.data)
          // Successfully processed, remove from queue
        } catch (error) {
          item.retryCount++
          if (item.retryCount <= item.maxRetries) {
            // Re-queue for retry
            this.offlineQueue.push(item)
          } else {
            console.warn(`Failed to process queue item after ${item.maxRetries} retries:`, item)
          }
        }
      }

      this.saveOfflineQueue()
    } finally {
      this.isProcessingQueue = false
    }
  }

  getQueueStatus(): { pending: number; failed: number } {
    const failed = this.offlineQueue.filter(item => item.retryCount > item.maxRetries).length
    const pending = this.offlineQueue.length - failed

    return { pending, failed }
  }

  clearQueue(): void {
    this.offlineQueue = []
    this.saveOfflineQueue()
  }

  clearCache(): void {
    // Clear all cached responses
    const keys = []
    for (let i = 0; i < 1000; i++) {
      const key = `cache_${i}`
      if (this.storage.getItem(key)) {
        keys.push(key)
      }
    }
    keys.forEach(key => this.storage.removeItem(key))
  }
}

// Mock server
const mockData = {
  cart: {
    id: 'cart-1',
    userId: 'user-1',
    items: [
      {
        id: 'item-1',
        productId: 'product-1',
        quantity: 1,
        unitPrice: 999.99,
        totalPrice: 999.99,
        name: 'iPhone 15 Pro',
        lastModified: Date.now()
      }
    ],
    total: 999.99,
    currency: 'USD',
    lastModified: Date.now(),
    version: 1
  },
  products: [
    {
      id: 'product-1',
      name: 'iPhone 15 Pro',
      price: 999.99,
      currency: 'USD',
      category: 'smartphones',
      status: 'active'
    }
  ]
}

const server = setupServer(
  rest.get('http://localhost:8080/api/v1/commerce/cart', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: mockData.cart
      })
    )
  }),

  rest.get('http://localhost:8080/api/v1/commerce/products', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        products: mockData.products,
        total: mockData.products.length
      })
    )
  }),

  rest.post('http://localhost:8080/api/v1/commerce/cart/items', (req, res, ctx) => {
    const newItem = req.body as any
    mockData.cart.items.push(newItem)
    mockData.cart.total += newItem.totalPrice
    mockData.cart.version++

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: mockData.cart
      })
    )
  }),

  rest.put('http://localhost:8080/api/v1/commerce/cart/items/:itemId', (req, res, ctx) => {
    const { itemId } = req.params
    const updates = req.body as any

    const itemIndex = mockData.cart.items.findIndex(item => item.id === itemId)
    if (itemIndex !== -1) {
      mockData.cart.items[itemIndex] = { ...mockData.cart.items[itemIndex], ...updates }
      mockData.cart.total = mockData.cart.items.reduce((sum, item) => sum + item.totalPrice, 0)
      mockData.cart.version++
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: mockData.cart
      })
    )
  }),

  rest.delete('http://localhost:8080/api/v1/commerce/cart/items/:itemId', (req, res, ctx) => {
    const { itemId } = req.params

    mockData.cart.items = mockData.cart.items.filter(item => item.id !== itemId)
    mockData.cart.total = mockData.cart.items.reduce((sum, item) => sum + item.totalPrice, 0)
    mockData.cart.version++

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: mockData.cart
      })
    )
  })
)

describe('Offline/Online Integration Tests', () => {
  let apiClient: OfflineCapableAPIClient
  let storage: MockStorageAdapter
  let network: MockNetworkAdapter

  beforeAll(() => {
    server.listen({ onUnhandledRequest: 'error' })
  })

  afterAll(() => {
    server.close()
  })

  beforeEach(() => {
    storage = new MockStorageAdapter()
    network = new MockNetworkAdapter()
    apiClient = new OfflineCapableAPIClient('http://localhost:8080', storage, network)

    // Reset mock data
    mockData.cart = {
      id: 'cart-1',
      userId: 'user-1',
      items: [
        {
          id: 'item-1',
          productId: 'product-1',
          quantity: 1,
          unitPrice: 999.99,
          totalPrice: 999.99,
          name: 'iPhone 15 Pro',
          lastModified: Date.now()
        }
      ],
      total: 999.99,
      currency: 'USD',
      lastModified: Date.now(),
      version: 1
    }
  })

  afterEach(() => {
    server.resetHandlers()
    network.setOnline(true)
  })

  describe('Online Behavior', () => {
    it('should make normal API requests when online', async () => {
      // Ensure we're online
      network.setOnline(true)

      // Make a GET request
      const cartResponse = await apiClient.request('GET', '/api/v1/commerce/cart')
      expect(cartResponse.success).toBe(true)
      expect(cartResponse.cart).toBeDefined()
      expect(cartResponse.cart.id).toBe('cart-1')

      // Make a POST request
      const newItem = {
        productId: 'product-2',
        quantity: 1,
        unitPrice: 199.99,
        totalPrice: 199.99,
        name: 'AirPods Pro'
      }

      const addResponse = await apiClient.request('POST', '/api/v1/commerce/cart/items', newItem)
      expect(addResponse.success).toBe(true)
      expect(addResponse.cart.items.length).toBe(2)
    })

    it('should cache GET responses for offline access', async () => {
      network.setOnline(true)

      // Make request while online (should cache)
      await apiClient.request('GET', '/api/v1/commerce/cart')

      // Verify data is cached
      const cachedData = storage.getItem('cache_/api/v1/commerce/cart')
      expect(cachedData).toBeDefined()

      const parsed = JSON.parse(cachedData!)
      expect(parsed.cart.id).toBe('cart-1')
    })
  })

  describe('Offline Behavior', () => {
    it('should return cached data for GET requests when offline', async () => {
      // First, make request while online to populate cache
      network.setOnline(true)
      await apiClient.request('GET', '/api/v1/commerce/cart')

      // Go offline
      network.setOnline(false)

      // Should return cached data
      const response = await apiClient.request('GET', '/api/v1/commerce/cart')
      expect(response.success).toBe(true)
      expect(response.cart.id).toBe('cart-1')
    })

    it('should queue POST requests when offline', async () => {
      network.setOnline(false)

      const newItem = {
        productId: 'product-offline',
        quantity: 1,
        unitPrice: 99.99,
        totalPrice: 99.99,
        name: 'Offline Item'
      }

      // Should return optimistic response
      const response = await apiClient.request('POST', '/api/v1/commerce/cart/items', newItem)
      expect(response.success).toBe(true)
      expect(response.cart.isOptimistic).toBe(true)

      // Check queue status
      const queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBe(1)
      expect(queueStatus.failed).toBe(0)
    })

    it('should queue PUT requests when offline', async () => {
      network.setOnline(false)

      const updates = {
        quantity: 3,
        totalPrice: 2999.97
      }

      // Should return optimistic response
      const response = await apiClient.request('PUT', '/api/v1/commerce/cart/items/item-1', updates)
      expect(response.success).toBe(true)
      expect(response.cart.isOptimistic).toBe(true)

      // Check queue
      const queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBe(1)
    })

    it('should queue DELETE requests when offline', async () => {
      network.setOnline(false)

      // Should return optimistic response
      const response = await apiClient.request('DELETE', '/api/v1/commerce/cart/items/item-1')
      expect(response.success).toBe(true)
      expect(response.isOptimistic).toBe(true)

      // Check queue
      const queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBe(1)
    })

    it('should handle multiple offline operations', async () => {
      network.setOnline(false)

      // Queue multiple operations
      await apiClient.request('POST', '/api/v1/commerce/cart/items', {
        productId: 'product-1',
        quantity: 1,
        unitPrice: 100,
        totalPrice: 100,
        name: 'Item 1'
      })

      await apiClient.request('POST', '/api/v1/commerce/cart/items', {
        productId: 'product-2',
        quantity: 2,
        unitPrice: 200,
        totalPrice: 400,
        name: 'Item 2'
      })

      await apiClient.request('PUT', '/api/v1/commerce/cart/items/item-1', {
        quantity: 5
      })

      await apiClient.request('DELETE', '/api/v1/commerce/cart/items/item-2')

      // Check queue has all operations
      const queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBe(4)
    })
  })

  describe('Online/Offline Transitions', () => {
    it('should process queued operations when coming back online', async () => {
      // Start offline and queue operations
      network.setOnline(false)

      await apiClient.request('POST', '/api/v1/commerce/cart/items', {
        id: 'item-queued',
        productId: 'product-queued',
        quantity: 1,
        unitPrice: 299.99,
        totalPrice: 299.99,
        name: 'Queued Item'
      })

      // Verify operation is queued
      let queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBe(1)

      // Come back online
      network.setOnline(true)

      // Wait for queue processing
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Queue should be empty
      queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBe(0)

      // Verify the operation was actually processed
      const cartResponse = await apiClient.request('GET', '/api/v1/commerce/cart')
      expect(cartResponse.cart.items.length).toBe(2) // Original + queued
    })

    it('should handle rapid online/offline transitions', async () => {
      // Rapid transitions
      network.setOnline(false)
      await apiClient.request('POST', '/api/v1/commerce/cart/items', { id: 'item-1' })

      network.setOnline(true)
      await new Promise(resolve => setTimeout(resolve, 100))

      network.setOnline(false)
      await apiClient.request('POST', '/api/v1/commerce/cart/items', { id: 'item-2' })

      network.setOnline(true)
      await new Promise(resolve => setTimeout(resolve, 100))

      network.setOnline(false)
      await apiClient.request('POST', '/api/v1/commerce/cart/items', { id: 'item-3' })

      network.setOnline(true)

      // Wait for all processing
      await new Promise(resolve => setTimeout(resolve, 2000))

      // All operations should eventually be processed
      const queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBe(0)
    })

    it('should handle network errors during queue processing', async () => {
      // Queue operation while offline
      network.setOnline(false)
      await apiClient.request('POST', '/api/v1/commerce/cart/items', {
        id: 'item-error',
        productId: 'product-error',
        quantity: 1,
        unitPrice: 99.99,
        totalPrice: 99.99,
        name: 'Error Item'
      })

      // Mock server error
      server.use(
        rest.post('http://localhost:8080/api/v1/commerce/cart/items', (req, res, ctx) => {
          return res(ctx.status(500), ctx.json({ error: 'Server error' }))
        })
      )

      // Come back online
      network.setOnline(true)

      // Wait for queue processing (should fail and retry)
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Should still have pending items (being retried)
      const queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBeGreaterThan(0)

      // Restore server
      server.restoreHandlers()

      // Wait more time for retry
      await new Promise(resolve => setTimeout(resolve, 2000))

      // Should eventually succeed
      const finalQueueStatus = apiClient.getQueueStatus()
      expect(finalQueueStatus.pending).toBe(0)
    })
  })

  describe('Data Consistency', () => {
    it('should maintain data consistency between optimistic updates and server state', async () => {
      // Start with cached cart
      network.setOnline(true)
      const initialCart = await apiClient.request('GET', '/api/v1/commerce/cart')

      // Go offline and make optimistic update
      network.setOnline(false)
      const optimisticResponse = await apiClient.request('PUT', '/api/v1/commerce/cart/items/item-1', {
        quantity: 5,
        totalPrice: 4999.95
      })

      expect(optimisticResponse.cart.isOptimistic).toBe(true)

      // Come back online
      network.setOnline(true)

      // Wait for sync
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Get fresh data
      const finalCart = await apiClient.request('GET', '/api/v1/commerce/cart')

      // Should have the updated quantity
      const updatedItem = finalCart.cart.items.find((item: any) => item.id === 'item-1')
      expect(updatedItem.quantity).toBe(5)
      expect(updatedItem.totalPrice).toBe(4999.95)
    })

    it('should handle conflicts between optimistic updates and server changes', async () => {
      // Make optimistic update while offline
      network.setOnline(false)
      await apiClient.request('PUT', '/api/v1/commerce/cart/items/item-1', {
        quantity: 3
      })

      // Simulate server-side change (another device)
      mockData.cart.items[0].quantity = 2
      mockData.cart.items[0].totalPrice = 1999.98
      mockData.cart.version++

      // Come back online
      network.setOnline(true)

      // Wait for sync
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Get final state
      const finalCart = await apiClient.request('GET', '/api/v1/commerce/cart')

      // Should have reconciled state (last write wins or merge strategy)
      const item = finalCart.cart.items.find((item: any) => item.id === 'item-1')
      expect(item).toBeDefined()
      expect(item.quantity).toBeGreaterThan(0) // Should have some valid quantity
    })

    it('should preserve operation order during queue processing', async () => {
      network.setOnline(false)

      // Queue operations in specific order
      await apiClient.request('POST', '/api/v1/commerce/cart/items', {
        id: 'item-order-1',
        productId: 'product-1',
        quantity: 1,
        unitPrice: 100,
        totalPrice: 100,
        name: 'Order Item 1'
      })

      await apiClient.request('POST', '/api/v1/commerce/cart/items', {
        id: 'item-order-2',
        productId: 'product-2',
        quantity: 1,
        unitPrice: 200,
        totalPrice: 200,
        name: 'Order Item 2'
      })

      await apiClient.request('PUT', '/api/v1/commerce/cart/items/item-order-1', {
        quantity: 2,
        totalPrice: 200
      })

      // Come back online
      network.setOnline(true)

      // Wait for processing
      await new Promise(resolve => setTimeout(resolve, 2000))

      // Verify final state reflects all operations in correct order
      const finalCart = await apiClient.request('GET', '/api/v1/commerce/cart')
      expect(finalCart.cart.items.length).toBeGreaterThanOrEqual(3) // Original + 2 new

      const item1 = finalCart.cart.items.find((item: any) => item.id === 'item-order-1')
      expect(item1?.quantity).toBe(2) // Should reflect the PUT update
    })
  })

  describe('Cache Management', () => {
    it('should invalidate cache when server data is newer', async () => {
      // Cache initial data
      network.setOnline(true)
      await apiClient.request('GET', '/api/v1/commerce/cart')

      // Simulate server-side change
      mockData.cart.items.push({
        id: 'item-server-added',
        productId: 'product-server',
        quantity: 1,
        unitPrice: 150,
        totalPrice: 150,
        name: 'Server Added Item',
        lastModified: Date.now()
      })
      mockData.cart.total += 150
      mockData.cart.version++

      // Request fresh data
      const freshCart = await apiClient.request('GET', '/api/v1/commerce/cart')

      // Should have the server-side changes
      expect(freshCart.cart.items.length).toBe(2)
      expect(freshCart.cart.items.some((item: any) => item.id === 'item-server-added')).toBe(true)
    })

    it('should handle cache expiration', async () => {
      // This would test TTL-based cache expiration
      // For simplicity, we'll test manual cache clearing

      network.setOnline(true)
      await apiClient.request('GET', '/api/v1/commerce/cart')

      // Verify cached
      let cachedData = storage.getItem('cache_/api/v1/commerce/cart')
      expect(cachedData).toBeDefined()

      // Clear cache
      apiClient.clearCache()

      // Verify cache cleared
      cachedData = storage.getItem('cache_/api/v1/commerce/cart')
      expect(cachedData).toBeNull()
    })

    it('should handle storage quota exceeded', async () => {
      // Simulate storage full scenario
      const originalSetItem = storage.setItem.bind(storage)
      let callCount = 0

      storage.setItem = (key: string, value: string) => {
        callCount++
        if (callCount > 5) {
          throw new Error('QuotaExceededError')
        }
        originalSetItem(key, value)
      }

      network.setOnline(false)

      // Try to queue many operations (should handle storage errors gracefully)
      for (let i = 0; i < 10; i++) {
        try {
          await apiClient.request('POST', '/api/v1/commerce/cart/items', {
            id: `item-${i}`,
            productId: `product-${i}`,
            quantity: 1,
            unitPrice: 100,
            totalPrice: 100,
            name: `Item ${i}`
          })
        } catch (error) {
          // Should handle storage errors gracefully
          expect(error.message).toContain('QuotaExceededError')
        }
      }

      // Should still function despite storage errors
      const queueStatus = apiClient.getQueueStatus()
      expect(queueStatus.pending).toBeLessThanOrEqual(5)
    })
  })
})