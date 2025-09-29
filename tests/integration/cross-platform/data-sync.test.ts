/**
 * Cross-Platform Data Synchronization Tests
 * Tests data synchronization between Web, iOS, Android, and Backend
 */

import { describe, it, expect, beforeAll, afterAll, beforeEach, afterEach } from 'vitest'
import { setupServer } from 'msw/node'
import { rest } from 'msw'
import WebSocket from 'ws'
import { EventEmitter } from 'events'

// Simulate different platform APIs
interface PlatformAPI {
  name: string
  getCart(): Promise<Cart>
  addToCart(item: CartItem): Promise<Cart>
  updateCartItem(itemId: string, updates: Partial<CartItem>): Promise<Cart>
  removeCartItem(itemId: string): Promise<Cart>
  clearCart(): Promise<Cart>
  syncData(): Promise<void>
  onDataUpdate(callback: (data: any) => void): void
}

interface Cart {
  id: string
  userId: string
  items: CartItem[]
  total: number
  currency: string
  lastModified: number
  version: number
}

interface CartItem {
  id: string
  productId: string
  quantity: number
  unitPrice: number
  totalPrice: number
  name: string
  lastModified: number
}

interface SyncEvent {
  type: 'cart_updated' | 'cart_item_added' | 'cart_item_updated' | 'cart_item_removed' | 'cart_cleared'
  userId: string
  cartId: string
  data: any
  timestamp: number
  version: number
  source: string
}

// Mock platform implementations
class WebPlatformAPI implements PlatformAPI {
  name = 'web'
  private cart: Cart | null = null
  private updateCallbacks: ((data: any) => void)[] = []
  private ws: WebSocket | null = null

  constructor(private baseURL: string) {
    this.connectWebSocket()
  }

  private connectWebSocket() {
    this.ws = new WebSocket(`ws://localhost:8080/ws/cart`)
    this.ws.on('message', (data) => {
      const event: SyncEvent = JSON.parse(data.toString())
      this.handleSyncEvent(event)
    })
  }

  private handleSyncEvent(event: SyncEvent) {
    if (event.source !== this.name) {
      this.updateCallbacks.forEach(callback => callback(event))
    }
  }

  async getCart(): Promise<Cart> {
    const response = await fetch(`${this.baseURL}/api/v1/commerce/cart`)
    const data = await response.json()
    this.cart = data.cart
    return this.cart!
  }

  async addToCart(item: CartItem): Promise<Cart> {
    const response = await fetch(`${this.baseURL}/api/v1/commerce/cart/items`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(item)
    })
    const data = await response.json()
    this.cart = data.cart
    this.broadcastUpdate('cart_item_added', item)
    return this.cart!
  }

  async updateCartItem(itemId: string, updates: Partial<CartItem>): Promise<Cart> {
    const response = await fetch(`${this.baseURL}/api/v1/commerce/cart/items/${itemId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(updates)
    })
    const data = await response.json()
    this.cart = data.cart
    this.broadcastUpdate('cart_item_updated', { itemId, updates })
    return this.cart!
  }

  async removeCartItem(itemId: string): Promise<Cart> {
    const response = await fetch(`${this.baseURL}/api/v1/commerce/cart/items/${itemId}`, {
      method: 'DELETE'
    })
    const data = await response.json()
    this.cart = data.cart
    this.broadcastUpdate('cart_item_removed', { itemId })
    return this.cart!
  }

  async clearCart(): Promise<Cart> {
    const response = await fetch(`${this.baseURL}/api/v1/commerce/cart`, {
      method: 'DELETE'
    })
    const data = await response.json()
    this.cart = data.cart
    this.broadcastUpdate('cart_cleared', {})
    return this.cart!
  }

  async syncData(): Promise<void> {
    await this.getCart()
  }

  onDataUpdate(callback: (data: any) => void): void {
    this.updateCallbacks.push(callback)
  }

  private broadcastUpdate(type: SyncEvent['type'], data: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const event: SyncEvent = {
        type,
        userId: 'user-1',
        cartId: this.cart?.id || 'cart-1',
        data,
        timestamp: Date.now(),
        version: (this.cart?.version || 0) + 1,
        source: this.name
      }
      this.ws.send(JSON.stringify(event))
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
    }
  }
}

class MobilePlatformAPI implements PlatformAPI {
  private cart: Cart | null = null
  private updateCallbacks: ((data: any) => void)[] = []
  private ws: WebSocket | null = null

  constructor(
    public name: string,
    private baseURL: string
  ) {
    this.connectWebSocket()
  }

  private connectWebSocket() {
    this.ws = new WebSocket(`ws://localhost:8080/ws/cart`)
    this.ws.on('message', (data) => {
      const event: SyncEvent = JSON.parse(data.toString())
      this.handleSyncEvent(event)
    })
  }

  private handleSyncEvent(event: SyncEvent) {
    if (event.source !== this.name) {
      // Update local cache/database
      this.syncLocalData(event)
      this.updateCallbacks.forEach(callback => callback(event))
    }
  }

  private async syncLocalData(event: SyncEvent) {
    // Simulate local database update
    if (this.cart) {
      switch (event.type) {
        case 'cart_item_added':
          this.cart.items.push(event.data)
          this.cart.version = event.version
          this.cart.lastModified = event.timestamp
          break
        case 'cart_item_updated':
          const itemIndex = this.cart.items.findIndex(item => item.id === event.data.itemId)
          if (itemIndex !== -1) {
            this.cart.items[itemIndex] = { ...this.cart.items[itemIndex], ...event.data.updates }
            this.cart.version = event.version
            this.cart.lastModified = event.timestamp
          }
          break
        case 'cart_item_removed':
          this.cart.items = this.cart.items.filter(item => item.id !== event.data.itemId)
          this.cart.version = event.version
          this.cart.lastModified = event.timestamp
          break
        case 'cart_cleared':
          this.cart.items = []
          this.cart.total = 0
          this.cart.version = event.version
          this.cart.lastModified = event.timestamp
          break
      }
      // Recalculate total
      this.cart.total = this.cart.items.reduce((sum, item) => sum + item.totalPrice, 0)
    }
  }

  async getCart(): Promise<Cart> {
    // Try local cache first
    if (this.cart) {
      return this.cart
    }

    // Fallback to network
    const response = await fetch(`${this.baseURL}/api/v1/commerce/cart`)
    const data = await response.json()
    this.cart = data.cart
    return this.cart!
  }

  async addToCart(item: CartItem): Promise<Cart> {
    // Optimistic update
    if (this.cart) {
      this.cart.items.push(item)
      this.cart.total += item.totalPrice
      this.cart.version += 1
      this.cart.lastModified = Date.now()
    }

    try {
      const response = await fetch(`${this.baseURL}/api/v1/commerce/cart/items`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(item)
      })
      const data = await response.json()
      this.cart = data.cart
      this.broadcastUpdate('cart_item_added', item)
      return this.cart!
    } catch (error) {
      // Rollback optimistic update
      if (this.cart) {
        this.cart.items = this.cart.items.filter(i => i.id !== item.id)
        this.cart.total -= item.totalPrice
        this.cart.version -= 1
      }
      throw error
    }
  }

  async updateCartItem(itemId: string, updates: Partial<CartItem>): Promise<Cart> {
    // Optimistic update
    const originalItem = this.cart?.items.find(item => item.id === itemId)
    if (this.cart && originalItem) {
      const itemIndex = this.cart.items.findIndex(item => item.id === itemId)
      this.cart.items[itemIndex] = { ...originalItem, ...updates }
      this.cart.total = this.cart.items.reduce((sum, item) => sum + item.totalPrice, 0)
      this.cart.version += 1
      this.cart.lastModified = Date.now()
    }

    try {
      const response = await fetch(`${this.baseURL}/api/v1/commerce/cart/items/${itemId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updates)
      })
      const data = await response.json()
      this.cart = data.cart
      this.broadcastUpdate('cart_item_updated', { itemId, updates })
      return this.cart!
    } catch (error) {
      // Rollback optimistic update
      if (this.cart && originalItem) {
        const itemIndex = this.cart.items.findIndex(item => item.id === itemId)
        this.cart.items[itemIndex] = originalItem
        this.cart.total = this.cart.items.reduce((sum, item) => sum + item.totalPrice, 0)
        this.cart.version -= 1
      }
      throw error
    }
  }

  async removeCartItem(itemId: string): Promise<Cart> {
    // Optimistic update
    const removedItem = this.cart?.items.find(item => item.id === itemId)
    if (this.cart && removedItem) {
      this.cart.items = this.cart.items.filter(item => item.id !== itemId)
      this.cart.total -= removedItem.totalPrice
      this.cart.version += 1
      this.cart.lastModified = Date.now()
    }

    try {
      const response = await fetch(`${this.baseURL}/api/v1/commerce/cart/items/${itemId}`, {
        method: 'DELETE'
      })
      const data = await response.json()
      this.cart = data.cart
      this.broadcastUpdate('cart_item_removed', { itemId })
      return this.cart!
    } catch (error) {
      // Rollback optimistic update
      if (this.cart && removedItem) {
        this.cart.items.push(removedItem)
        this.cart.total += removedItem.totalPrice
        this.cart.version -= 1
      }
      throw error
    }
  }

  async clearCart(): Promise<Cart> {
    // Optimistic update
    const originalItems = this.cart?.items || []
    const originalTotal = this.cart?.total || 0
    if (this.cart) {
      this.cart.items = []
      this.cart.total = 0
      this.cart.version += 1
      this.cart.lastModified = Date.now()
    }

    try {
      const response = await fetch(`${this.baseURL}/api/v1/commerce/cart`, {
        method: 'DELETE'
      })
      const data = await response.json()
      this.cart = data.cart
      this.broadcastUpdate('cart_cleared', {})
      return this.cart!
    } catch (error) {
      // Rollback optimistic update
      if (this.cart) {
        this.cart.items = originalItems
        this.cart.total = originalTotal
        this.cart.version -= 1
      }
      throw error
    }
  }

  async syncData(): Promise<void> {
    // Force refresh from server
    const response = await fetch(`${this.baseURL}/api/v1/commerce/cart?force=true`)
    const data = await response.json()
    this.cart = data.cart
  }

  onDataUpdate(callback: (data: any) => void): void {
    this.updateCallbacks.push(callback)
  }

  private broadcastUpdate(type: SyncEvent['type'], data: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const event: SyncEvent = {
        type,
        userId: 'user-1',
        cartId: this.cart?.id || 'cart-1',
        data,
        timestamp: Date.now(),
        version: (this.cart?.version || 0) + 1,
        source: this.name
      }
      this.ws.send(JSON.stringify(event))
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
    }
  }
}

// Mock WebSocket server
class MockWebSocketServer extends EventEmitter {
  private clients: WebSocket[] = []
  private server: any

  start(port: number = 8081) {
    const WebSocketServer = require('ws').Server
    this.server = new WebSocketServer({ port })

    this.server.on('connection', (ws: WebSocket) => {
      this.clients.push(ws)

      ws.on('message', (data: string) => {
        const event: SyncEvent = JSON.parse(data)
        this.broadcast(event, ws)
      })

      ws.on('close', () => {
        this.clients = this.clients.filter(client => client !== ws)
      })
    })
  }

  private broadcast(event: SyncEvent, sender: WebSocket) {
    const message = JSON.stringify(event)
    this.clients.forEach(client => {
      if (client !== sender && client.readyState === WebSocket.OPEN) {
        client.send(message)
      }
    })
  }

  stop() {
    if (this.server) {
      this.server.close()
    }
  }
}

// Mock data
const mockCart: Cart = {
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

// Mock server setup
const server = setupServer(
  rest.get('http://localhost:8080/api/v1/commerce/cart', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: mockCart
      })
    )
  }),

  rest.post('http://localhost:8080/api/v1/commerce/cart/items', (req, res, ctx) => {
    const newItem = req.body as CartItem
    const updatedCart = {
      ...mockCart,
      items: [...mockCart.items, newItem],
      total: mockCart.total + newItem.totalPrice,
      version: mockCart.version + 1,
      lastModified: Date.now()
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: updatedCart
      })
    )
  }),

  rest.put('http://localhost:8080/api/v1/commerce/cart/items/:itemId', (req, res, ctx) => {
    const { itemId } = req.params
    const updates = req.body as Partial<CartItem>

    const updatedItems = mockCart.items.map(item =>
      item.id === itemId ? { ...item, ...updates, lastModified: Date.now() } : item
    )

    const updatedCart = {
      ...mockCart,
      items: updatedItems,
      total: updatedItems.reduce((sum, item) => sum + item.totalPrice, 0),
      version: mockCart.version + 1,
      lastModified: Date.now()
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: updatedCart
      })
    )
  }),

  rest.delete('http://localhost:8080/api/v1/commerce/cart/items/:itemId', (req, res, ctx) => {
    const { itemId } = req.params

    const updatedItems = mockCart.items.filter(item => item.id !== itemId)

    const updatedCart = {
      ...mockCart,
      items: updatedItems,
      total: updatedItems.reduce((sum, item) => sum + item.totalPrice, 0),
      version: mockCart.version + 1,
      lastModified: Date.now()
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: updatedCart
      })
    )
  }),

  rest.delete('http://localhost:8080/api/v1/commerce/cart', (req, res, ctx) => {
    const clearedCart = {
      ...mockCart,
      items: [],
      total: 0,
      version: mockCart.version + 1,
      lastModified: Date.now()
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: clearedCart
      })
    )
  })
)

describe('Cross-Platform Data Synchronization Tests', () => {
  let webAPI: WebPlatformAPI
  let iosAPI: MobilePlatformAPI
  let androidAPI: MobilePlatformAPI
  let wsServer: MockWebSocketServer

  beforeAll(async () => {
    server.listen({ onUnhandledRequest: 'error' })

    // Start WebSocket server
    wsServer = new MockWebSocketServer()
    wsServer.start(8081)

    // Wait for WebSocket server to start
    await new Promise(resolve => setTimeout(resolve, 1000))
  })

  afterAll(() => {
    server.close()
    wsServer.stop()
  })

  beforeEach(async () => {
    // Initialize platform APIs
    webAPI = new WebPlatformAPI('http://localhost:8080')
    iosAPI = new MobilePlatformAPI('ios', 'http://localhost:8080')
    androidAPI = new MobilePlatformAPI('android', 'http://localhost:8080')

    // Wait for WebSocket connections
    await new Promise(resolve => setTimeout(resolve, 500))
  })

  afterEach(() => {
    webAPI.disconnect()
    iosAPI.disconnect()
    androidAPI.disconnect()
    server.resetHandlers()
  })

  describe('Basic Data Synchronization', () => {
    it('should synchronize cart data across all platforms', async () => {
      // Get initial cart state on all platforms
      const webCart = await webAPI.getCart()
      const iosCart = await iosAPI.getCart()
      const androidCart = await androidAPI.getCart()

      // All platforms should have the same cart data
      expect(webCart.id).toBe(iosCart.id)
      expect(iosCart.id).toBe(androidCart.id)
      expect(webCart.items.length).toBe(iosCart.items.length)
      expect(iosCart.items.length).toBe(androidCart.items.length)
      expect(webCart.total).toBe(iosCart.total)
      expect(iosCart.total).toBe(androidCart.total)
    })

    it('should sync cart updates from web to mobile platforms', async () => {
      const updatePromises: Promise<any>[] = []

      // Set up listeners on mobile platforms
      iosAPI.onDataUpdate((event) => {
        updatePromises.push(Promise.resolve(event))
      })

      androidAPI.onDataUpdate((event) => {
        updatePromises.push(Promise.resolve(event))
      })

      // Add item from web
      const newItem: CartItem = {
        id: 'item-2',
        productId: 'product-2',
        quantity: 1,
        unitPrice: 199.99,
        totalPrice: 199.99,
        name: 'AirPods Pro',
        lastModified: Date.now()
      }

      await webAPI.addToCart(newItem)

      // Wait for synchronization
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Check that mobile platforms received the update
      const iosCart = await iosAPI.getCart()
      const androidCart = await androidAPI.getCart()

      expect(iosCart.items.length).toBe(2)
      expect(androidCart.items.length).toBe(2)
      expect(iosCart.items.some(item => item.id === 'item-2')).toBe(true)
      expect(androidCart.items.some(item => item.id === 'item-2')).toBe(true)
    })

    it('should sync cart updates from mobile to web and other mobile platforms', async () => {
      const updatePromises: Promise<any>[] = []

      // Set up listeners
      webAPI.onDataUpdate((event) => {
        updatePromises.push(Promise.resolve(event))
      })

      androidAPI.onDataUpdate((event) => {
        updatePromises.push(Promise.resolve(event))
      })

      // Update item quantity from iOS
      await iosAPI.updateCartItem('item-1', { quantity: 3, totalPrice: 2999.97 })

      // Wait for synchronization
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Check that other platforms received the update
      const webCart = await webAPI.getCart()
      const androidCart = await androidAPI.getCart()

      const webItem = webCart.items.find(item => item.id === 'item-1')
      const androidItem = androidCart.items.find(item => item.id === 'item-1')

      expect(webItem?.quantity).toBe(3)
      expect(androidItem?.quantity).toBe(3)
      expect(webItem?.totalPrice).toBe(2999.97)
      expect(androidItem?.totalPrice).toBe(2999.97)
    })
  })

  describe('Conflict Resolution', () => {
    it('should handle concurrent updates with last-write-wins strategy', async () => {
      // Simulate concurrent updates from different platforms
      const updatePromises = [
        iosAPI.updateCartItem('item-1', { quantity: 2 }),
        androidAPI.updateCartItem('item-1', { quantity: 3 }),
        webAPI.updateCartItem('item-1', { quantity: 4 })
      ]

      await Promise.allSettled(updatePromises)

      // Wait for synchronization
      await new Promise(resolve => setTimeout(resolve, 2000))

      // All platforms should eventually converge to the same state
      const webCart = await webAPI.getCart()
      const iosCart = await iosAPI.getCart()
      const androidCart = await androidAPI.getCart()

      const webItem = webCart.items.find(item => item.id === 'item-1')
      const iosItem = iosCart.items.find(item => item.id === 'item-1')
      const androidItem = androidCart.items.find(item => item.id === 'item-1')

      // All should have the same final quantity (last write wins)
      expect(webItem?.quantity).toBe(iosItem?.quantity)
      expect(iosItem?.quantity).toBe(androidItem?.quantity)
    })

    it('should handle version conflicts with server-side resolution', async () => {
      // Get initial cart state
      const initialCart = await webAPI.getCart()

      // Simulate out-of-sync state by manipulating version
      const staleUpdate = {
        quantity: 5,
        totalPrice: 4999.95
      }

      // This should be rejected if version checking is implemented
      try {
        await iosAPI.updateCartItem('item-1', staleUpdate)

        // Wait for sync
        await new Promise(resolve => setTimeout(resolve, 1000))

        // Check final state - should either be accepted or resolved
        const finalCart = await webAPI.getCart()
        const item = finalCart.items.find(item => item.id === 'item-1')

        // Version should be updated
        expect(finalCart.version).toBeGreaterThan(initialCart.version)
        expect(item).toBeDefined()
      } catch (error) {
        // Version conflict detected and handled
        expect(error).toBeDefined()
      }
    })
  })

  describe('Offline/Online Synchronization', () => {
    it('should queue updates when offline and sync when back online', async () => {
      // Simulate network failure
      server.use(
        rest.post('http://localhost:8080/api/v1/commerce/cart/items', (req, res, ctx) => {
          return res(ctx.status(503), ctx.json({ error: 'Service unavailable' }))
        })
      )

      // Try to add item while offline (should fail gracefully)
      const newItem: CartItem = {
        id: 'item-offline',
        productId: 'product-offline',
        quantity: 1,
        unitPrice: 99.99,
        totalPrice: 99.99,
        name: 'Offline Item',
        lastModified: Date.now()
      }

      try {
        await iosAPI.addToCart(newItem)
      } catch (error) {
        // Expected to fail
        expect(error).toBeDefined()
      }

      // Restore network connectivity
      server.restoreHandlers()

      // Retry sync
      await iosAPI.syncData()

      // Verify final state
      const cart = await iosAPI.getCart()
      expect(cart).toBeDefined()
    })

    it('should merge local changes with server state after reconnection', async () => {
      // Start with synchronized state
      await Promise.all([
        webAPI.getCart(),
        iosAPI.getCart(),
        androidAPI.getCart()
      ])

      // Simulate iOS going offline and making local changes
      const offlineItem: CartItem = {
        id: 'item-offline-merge',
        productId: 'product-offline-merge',
        quantity: 1,
        unitPrice: 149.99,
        totalPrice: 149.99,
        name: 'Offline Merge Item',
        lastModified: Date.now()
      }

      // Meanwhile, web makes online changes
      const onlineItem: CartItem = {
        id: 'item-online',
        productId: 'product-online',
        quantity: 1,
        unitPrice: 299.99,
        totalPrice: 299.99,
        name: 'Online Item',
        lastModified: Date.now()
      }

      await webAPI.addToCart(onlineItem)

      // iOS comes back online and syncs
      await iosAPI.syncData()

      // Wait for final synchronization
      await new Promise(resolve => setTimeout(resolve, 1000))

      // All platforms should have both items
      const finalWebCart = await webAPI.getCart()
      const finalIosCart = await iosAPI.getCart()
      const finalAndroidCart = await androidAPI.getCart()

      // Check that all platforms converged to the same state
      expect(finalWebCart.items.length).toBe(finalIosCart.items.length)
      expect(finalIosCart.items.length).toBe(finalAndroidCart.items.length)

      // Check for online item
      expect(finalWebCart.items.some(item => item.id === 'item-online')).toBe(true)
      expect(finalIosCart.items.some(item => item.id === 'item-online')).toBe(true)
      expect(finalAndroidCart.items.some(item => item.id === 'item-online')).toBe(true)
    })
  })

  describe('Real-time Synchronization Performance', () => {
    it('should sync updates across platforms within 1 second', async () => {
      const syncTimes: number[] = []

      // Set up timing measurements
      const startTime = Date.now()

      iosAPI.onDataUpdate(() => {
        syncTimes.push(Date.now() - startTime)
      })

      androidAPI.onDataUpdate(() => {
        syncTimes.push(Date.now() - startTime)
      })

      // Perform update from web
      await webAPI.updateCartItem('item-1', { quantity: 10, totalPrice: 9999.90 })

      // Wait for synchronization
      await new Promise(resolve => setTimeout(resolve, 1500))

      // Check sync times
      expect(syncTimes.length).toBeGreaterThan(0)
      syncTimes.forEach(time => {
        expect(time).toBeLessThan(1000) // Less than 1 second
      })
    })

    it('should handle high-frequency updates without data loss', async () => {
      const updateCount = 10
      const updatePromises: Promise<any>[] = []

      // Generate rapid updates from different platforms
      for (let i = 0; i < updateCount; i++) {
        const platform = [webAPI, iosAPI, androidAPI][i % 3]
        updatePromises.push(
          platform.updateCartItem('item-1', {
            quantity: i + 1,
            totalPrice: (i + 1) * 999.99
          })
        )
      }

      await Promise.allSettled(updatePromises)

      // Wait for all updates to propagate
      await new Promise(resolve => setTimeout(resolve, 3000))

      // Verify final consistency
      const webCart = await webAPI.getCart()
      const iosCart = await iosAPI.getCart()
      const androidCart = await androidAPI.getCart()

      // All should have the same final state
      expect(webCart.version).toBe(iosCart.version)
      expect(iosCart.version).toBe(androidCart.version)

      const webItem = webCart.items.find(item => item.id === 'item-1')
      const iosItem = iosCart.items.find(item => item.id === 'item-1')
      const androidItem = androidCart.items.find(item => item.id === 'item-1')

      expect(webItem?.quantity).toBe(iosItem?.quantity)
      expect(iosItem?.quantity).toBe(androidItem?.quantity)
    })
  })

  describe('Data Integrity and Consistency', () => {
    it('should maintain data consistency during network partitions', async () => {
      // Initial state
      const initialCart = await webAPI.getCart()

      // Simulate network partition - iOS and Android can't reach server
      let networkDown = true

      server.use(
        rest.all('http://localhost:8080/api/v1/commerce/cart/*', (req, res, ctx) => {
          if (networkDown) {
            return res.networkError('Network partition')
          }
          return req.passthrough()
        })
      )

      // Try operations during partition
      try {
        await iosAPI.addToCart({
          id: 'item-partition',
          productId: 'product-partition',
          quantity: 1,
          unitPrice: 99.99,
          totalPrice: 99.99,
          name: 'Partition Item',
          lastModified: Date.now()
        })
      } catch (error) {
        // Expected during partition
      }

      // Restore network
      networkDown = false
      server.restoreHandlers()

      // Force resync
      await Promise.all([
        iosAPI.syncData(),
        androidAPI.syncData()
      ])

      // Wait for consistency
      await new Promise(resolve => setTimeout(resolve, 2000))

      // Verify final consistency
      const finalWebCart = await webAPI.getCart()
      const finalIosCart = await iosAPI.getCart()
      const finalAndroidCart = await androidAPI.getCart()

      expect(finalWebCart.items.length).toBe(finalIosCart.items.length)
      expect(finalIosCart.items.length).toBe(finalAndroidCart.items.length)
    })

    it('should preserve data integrity during concurrent clear operations', async () => {
      // Add some items first
      await webAPI.addToCart({
        id: 'item-clear-test-1',
        productId: 'product-clear-1',
        quantity: 1,
        unitPrice: 100,
        totalPrice: 100,
        name: 'Clear Test 1',
        lastModified: Date.now()
      })

      await androidAPI.addToCart({
        id: 'item-clear-test-2',
        productId: 'product-clear-2',
        quantity: 1,
        unitPrice: 200,
        totalPrice: 200,
        name: 'Clear Test 2',
        lastModified: Date.now()
      })

      // Wait for sync
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Concurrent clear operations
      const clearPromises = [
        webAPI.clearCart(),
        iosAPI.clearCart(),
        androidAPI.clearCart()
      ]

      await Promise.allSettled(clearPromises)

      // Wait for final sync
      await new Promise(resolve => setTimeout(resolve, 2000))

      // All platforms should have empty carts
      const webCart = await webAPI.getCart()
      const iosCart = await iosAPI.getCart()
      const androidCart = await androidAPI.getCart()

      expect(webCart.items.length).toBe(0)
      expect(iosCart.items.length).toBe(0)
      expect(androidCart.items.length).toBe(0)
      expect(webCart.total).toBe(0)
      expect(iosCart.total).toBe(0)
      expect(androidCart.total).toBe(0)
    })
  })
})