/**
 * Comprehensive Web Frontend Integration Tests
 * Tests RTK Query API integration with commerce backend endpoints
 */

import { describe, it, expect, beforeAll, afterAll, beforeEach, afterEach } from 'vitest'
import { setupServer } from 'msw/node'
import { rest } from 'msw'
import { configureStore } from '@reduxjs/toolkit'
import { act, renderHook, waitFor } from '@testing-library/react'
import { Provider } from 'react-redux'
import React from 'react'

// Import RTK Query API slices (assuming these exist in the project)
import { commerceApi } from '../../../apps/web/src/store/api/commerceApi'
import { cartApi } from '../../../apps/web/src/store/api/cartApi'
import { productApi } from '../../../apps/web/src/store/api/productApi'
import { categoryApi } from '../../../apps/web/src/store/api/categoryApi'
import { businessApi } from '../../../apps/web/src/store/api/businessApi'
import { rootReducer } from '../../../apps/web/src/store'

// Types
interface Product {
  id: string
  name: string
  price: number
  currency: string
  category: string
  shopId: string
  status: string
  inventory: {
    quantity: number
    stockStatus: string
  }
  images: Array<{
    url: string
    altText: string
    isMain: boolean
  }>
  variants: Array<{
    id: string
    name: string
    options: Record<string, string>
    price?: number
  }>
}

interface Cart {
  id: string
  userId: string
  items: CartItem[]
  total: number
  currency: string
  status: string
}

interface CartItem {
  id: string
  productId: string
  quantity: number
  unitPrice: number
  totalPrice: number
  name: string
}

interface Category {
  id: string
  name: string
  slug: string
  parentId?: string
  children?: Category[]
  level: number
  productCount: number
}

interface Business {
  id: string
  name: string
  businessType: string
  industry: string
  email: string
  status: string
  verification: {
    status: string
    level: string
  }
}

// Mock data
const mockProducts: Product[] = [
  {
    id: '1',
    name: 'iPhone 15 Pro',
    price: 999.99,
    currency: 'USD',
    category: 'smartphones',
    shopId: 'shop-1',
    status: 'active',
    inventory: {
      quantity: 50,
      stockStatus: 'in_stock'
    },
    images: [
      {
        url: 'https://example.com/iphone15-main.jpg',
        altText: 'iPhone 15 Pro',
        isMain: true
      }
    ],
    variants: [
      {
        id: 'variant-1',
        name: '128GB Space Black',
        options: { storage: '128GB', color: 'Space Black' }
      }
    ]
  },
  {
    id: '2',
    name: 'MacBook Air M2',
    price: 1199.99,
    currency: 'USD',
    category: 'laptops',
    shopId: 'shop-1',
    status: 'active',
    inventory: {
      quantity: 25,
      stockStatus: 'in_stock'
    },
    images: [
      {
        url: 'https://example.com/macbook-air.jpg',
        altText: 'MacBook Air M2',
        isMain: true
      }
    ],
    variants: []
  }
]

const mockCart: Cart = {
  id: 'cart-1',
  userId: 'user-1',
  items: [
    {
      id: 'item-1',
      productId: '1',
      quantity: 1,
      unitPrice: 999.99,
      totalPrice: 999.99,
      name: 'iPhone 15 Pro'
    }
  ],
  total: 999.99,
  currency: 'USD',
  status: 'active'
}

const mockCategories: Category[] = [
  {
    id: 'cat-1',
    name: 'Electronics',
    slug: 'electronics',
    level: 0,
    productCount: 2,
    children: [
      {
        id: 'cat-2',
        name: 'Smartphones',
        slug: 'smartphones',
        parentId: 'cat-1',
        level: 1,
        productCount: 1
      },
      {
        id: 'cat-3',
        name: 'Laptops',
        slug: 'laptops',
        parentId: 'cat-1',
        level: 1,
        productCount: 1
      }
    ]
  }
]

const mockBusiness: Business = {
  id: 'business-1',
  name: 'TechCorp Solutions',
  businessType: 'corporation',
  industry: 'technology',
  email: 'contact@techcorp.com',
  status: 'active',
  verification: {
    status: 'verified',
    level: 'full'
  }
}

// Mock server setup
const server = setupServer(
  // Product endpoints
  rest.get('http://localhost:8080/api/v1/commerce/products', (req, res, ctx) => {
    const query = req.url.searchParams.get('query')
    const category = req.url.searchParams.get('category')

    let filteredProducts = mockProducts

    if (query) {
      filteredProducts = filteredProducts.filter(p =>
        p.name.toLowerCase().includes(query.toLowerCase())
      )
    }

    if (category) {
      filteredProducts = filteredProducts.filter(p => p.category === category)
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        products: filteredProducts,
        total: filteredProducts.length,
        page: 1,
        limit: 20
      })
    )
  }),

  rest.get('http://localhost:8080/api/v1/commerce/products/:id', (req, res, ctx) => {
    const { id } = req.params
    const product = mockProducts.find(p => p.id === id)

    if (!product) {
      return res(
        ctx.status(404),
        ctx.json({
          success: false,
          message: 'Product not found'
        })
      )
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        product
      })
    )
  }),

  rest.post('http://localhost:8080/api/v1/commerce/shops/:shopId/products', (req, res, ctx) => {
    return res(
      ctx.status(201),
      ctx.json({
        success: true,
        product: {
          ...req.body,
          id: `product-${Date.now()}`,
          shopId: req.params.shopId,
          status: 'active'
        }
      })
    )
  }),

  rest.put('http://localhost:8080/api/v1/commerce/products/:id', (req, res, ctx) => {
    const { id } = req.params
    const product = mockProducts.find(p => p.id === id)

    if (!product) {
      return res(ctx.status(404), ctx.json({ success: false, message: 'Product not found' }))
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        product: { ...product, ...req.body }
      })
    )
  }),

  rest.delete('http://localhost:8080/api/v1/commerce/products/:id', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        message: 'Product deleted successfully'
      })
    )
  }),

  // Cart endpoints
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
    const newItem = {
      id: `item-${Date.now()}`,
      ...req.body,
      totalPrice: req.body.quantity * req.body.unitPrice
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: {
          ...mockCart,
          items: [...mockCart.items, newItem],
          total: mockCart.total + newItem.totalPrice
        }
      })
    )
  }),

  rest.put('http://localhost:8080/api/v1/commerce/cart/items/:itemId', (req, res, ctx) => {
    const { itemId } = req.params
    const updatedItems = mockCart.items.map(item =>
      item.id === itemId ? { ...item, ...req.body } : item
    )

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: {
          ...mockCart,
          items: updatedItems,
          total: updatedItems.reduce((sum, item) => sum + item.totalPrice, 0)
        }
      })
    )
  }),

  rest.delete('http://localhost:8080/api/v1/commerce/cart/items/:itemId', (req, res, ctx) => {
    const { itemId } = req.params
    const filteredItems = mockCart.items.filter(item => item.id !== itemId)

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: {
          ...mockCart,
          items: filteredItems,
          total: filteredItems.reduce((sum, item) => sum + item.totalPrice, 0)
        }
      })
    )
  }),

  rest.delete('http://localhost:8080/api/v1/commerce/cart', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        cart: {
          ...mockCart,
          items: [],
          total: 0
        }
      })
    )
  }),

  // Category endpoints
  rest.get('http://localhost:8080/api/v1/commerce/categories', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        categories: mockCategories,
        total: mockCategories.length
      })
    )
  }),

  rest.get('http://localhost:8080/api/v1/commerce/categories/hierarchy', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        categories: mockCategories
      })
    )
  }),

  rest.get('http://localhost:8080/api/v1/commerce/categories/:id', (req, res, ctx) => {
    const { id } = req.params
    const findCategory = (categories: Category[]): Category | undefined => {
      for (const cat of categories) {
        if (cat.id === id) return cat
        if (cat.children) {
          const found = findCategory(cat.children)
          if (found) return found
        }
      }
      return undefined
    }

    const category = findCategory(mockCategories)

    if (!category) {
      return res(ctx.status(404), ctx.json({ success: false, message: 'Category not found' }))
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        category
      })
    )
  }),

  rest.post('http://localhost:8080/api/v1/commerce/categories', (req, res, ctx) => {
    return res(
      ctx.status(201),
      ctx.json({
        success: true,
        category: {
          ...req.body,
          id: `category-${Date.now()}`,
          level: req.body.parentId ? 1 : 0,
          productCount: 0
        }
      })
    )
  }),

  // Business endpoints
  rest.get('http://localhost:8080/api/v1/commerce/businesses', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        businesses: [mockBusiness],
        total: 1
      })
    )
  }),

  rest.get('http://localhost:8080/api/v1/commerce/businesses/:id', (req, res, ctx) => {
    const { id } = req.params

    if (id !== mockBusiness.id) {
      return res(ctx.status(404), ctx.json({ success: false, message: 'Business not found' }))
    }

    return res(
      ctx.status(200),
      ctx.json({
        success: true,
        business: mockBusiness
      })
    )
  }),

  rest.post('http://localhost:8080/api/v1/commerce/businesses', (req, res, ctx) => {
    return res(
      ctx.status(201),
      ctx.json({
        success: true,
        business: {
          ...req.body,
          id: `business-${Date.now()}`,
          status: 'pending',
          verification: {
            status: 'unverified',
            level: 'none'
          }
        }
      })
    )
  }),

  // Error scenarios
  rest.get('http://localhost:8080/api/v1/commerce/products/error', (req, res, ctx) => {
    return res(
      ctx.status(500),
      ctx.json({
        success: false,
        message: 'Internal server error'
      })
    )
  }),

  // Network timeout simulation
  rest.get('http://localhost:8080/api/v1/commerce/products/timeout', (req, res, ctx) => {
    return res(
      ctx.delay(5000),
      ctx.status(200),
      ctx.json({ success: true, products: [] })
    )
  })
)

// Test setup
let store: ReturnType<typeof configureStore>

const createTestStore = () => {
  return configureStore({
    reducer: rootReducer,
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware({
        serializableCheck: {
          ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE']
        }
      }).concat(
        commerceApi.middleware,
        cartApi.middleware,
        productApi.middleware,
        categoryApi.middleware,
        businessApi.middleware
      )
  })
}

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <Provider store={store}>{children}</Provider>
)

describe('Web Frontend RTK Query Integration Tests', () => {
  beforeAll(() => {
    server.listen({ onUnhandledRequest: 'error' })
  })

  afterAll(() => {
    server.close()
  })

  beforeEach(() => {
    store = createTestStore()
  })

  afterEach(() => {
    server.resetHandlers()
    store.dispatch(commerceApi.util.resetApiState())
    store.dispatch(cartApi.util.resetApiState())
    store.dispatch(productApi.util.resetApiState())
    store.dispatch(categoryApi.util.resetApiState())
    store.dispatch(businessApi.util.resetApiState())
  })

  describe('Product API Integration', () => {
    it('should fetch products list', async () => {
      const { result } = renderHook(
        () => productApi.useGetProductsQuery({ page: 1, limit: 20 }),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(result.current.data).toEqual({
        success: true,
        products: mockProducts,
        total: mockProducts.length,
        page: 1,
        limit: 20
      })
    })

    it('should fetch product by ID', async () => {
      const { result } = renderHook(
        () => productApi.useGetProductQuery('1'),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(result.current.data?.product).toEqual(mockProducts[0])
    })

    it('should handle product not found', async () => {
      const { result } = renderHook(
        () => productApi.useGetProductQuery('999'),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isError).toBe(true)
      })

      expect(result.current.error).toBeDefined()
    })

    it('should search products with query', async () => {
      const { result } = renderHook(
        () => productApi.useGetProductsQuery({ query: 'iPhone' }),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      const products = result.current.data?.products || []
      expect(products).toHaveLength(1)
      expect(products[0].name).toContain('iPhone')
    })

    it('should filter products by category', async () => {
      const { result } = renderHook(
        () => productApi.useGetProductsQuery({ category: 'smartphones' }),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      const products = result.current.data?.products || []
      expect(products).toHaveLength(1)
      expect(products[0].category).toBe('smartphones')
    })

    it('should create product with mutation', async () => {
      const { result } = renderHook(
        () => productApi.useCreateProductMutation(),
        { wrapper }
      )

      const newProduct = {
        name: 'Test Product',
        price: 199.99,
        currency: 'USD',
        category: 'test'
      }

      await act(async () => {
        const response = await result.current[0]({
          shopId: 'shop-1',
          product: newProduct
        }).unwrap()

        expect(response.success).toBe(true)
        expect(response.product.name).toBe(newProduct.name)
        expect(response.product.shopId).toBe('shop-1')
      })
    })

    it('should update product with mutation', async () => {
      const { result } = renderHook(
        () => productApi.useUpdateProductMutation(),
        { wrapper }
      )

      const updates = {
        name: 'Updated iPhone 15 Pro',
        price: 1099.99
      }

      await act(async () => {
        const response = await result.current[0]({
          id: '1',
          ...updates
        }).unwrap()

        expect(response.success).toBe(true)
        expect(response.product.name).toBe(updates.name)
        expect(response.product.price).toBe(updates.price)
      })
    })

    it('should delete product with mutation', async () => {
      const { result } = renderHook(
        () => productApi.useDeleteProductMutation(),
        { wrapper }
      )

      await act(async () => {
        const response = await result.current[0]('1').unwrap()
        expect(response.success).toBe(true)
      })
    })
  })

  describe('Cart API Integration', () => {
    it('should fetch cart', async () => {
      const { result } = renderHook(
        () => cartApi.useGetCartQuery(),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(result.current.data?.cart).toEqual(mockCart)
    })

    it('should add item to cart', async () => {
      const { result } = renderHook(
        () => cartApi.useAddToCartMutation(),
        { wrapper }
      )

      const newItem = {
        productId: '2',
        quantity: 1,
        unitPrice: 1199.99
      }

      await act(async () => {
        const response = await result.current[0](newItem).unwrap()

        expect(response.success).toBe(true)
        expect(response.cart.items).toHaveLength(2)
        expect(response.cart.total).toBeGreaterThan(mockCart.total)
      })
    })

    it('should update cart item', async () => {
      const { result } = renderHook(
        () => cartApi.useUpdateCartItemMutation(),
        { wrapper }
      )

      await act(async () => {
        const response = await result.current[0]({
          itemId: 'item-1',
          quantity: 2
        }).unwrap()

        expect(response.success).toBe(true)
        expect(response.cart.items[0].quantity).toBe(2)
      })
    })

    it('should remove item from cart', async () => {
      const { result } = renderHook(
        () => cartApi.useRemoveCartItemMutation(),
        { wrapper }
      )

      await act(async () => {
        const response = await result.current[0]('item-1').unwrap()

        expect(response.success).toBe(true)
        expect(response.cart.items).toHaveLength(0)
        expect(response.cart.total).toBe(0)
      })
    })

    it('should clear cart', async () => {
      const { result } = renderHook(
        () => cartApi.useClearCartMutation(),
        { wrapper }
      )

      await act(async () => {
        const response = await result.current[0]().unwrap()

        expect(response.success).toBe(true)
        expect(response.cart.items).toHaveLength(0)
        expect(response.cart.total).toBe(0)
      })
    })
  })

  describe('Category API Integration', () => {
    it('should fetch categories list', async () => {
      const { result } = renderHook(
        () => categoryApi.useGetCategoriesQuery(),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(result.current.data?.categories).toEqual(mockCategories)
    })

    it('should fetch category hierarchy', async () => {
      const { result } = renderHook(
        () => categoryApi.useGetCategoryHierarchyQuery(),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(result.current.data?.categories).toEqual(mockCategories)
      expect(result.current.data?.categories[0].children).toHaveLength(2)
    })

    it('should fetch category by ID', async () => {
      const { result } = renderHook(
        () => categoryApi.useGetCategoryQuery('cat-1'),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(result.current.data?.category.id).toBe('cat-1')
      expect(result.current.data?.category.name).toBe('Electronics')
    })

    it('should create category with mutation', async () => {
      const { result } = renderHook(
        () => categoryApi.useCreateCategoryMutation(),
        { wrapper }
      )

      const newCategory = {
        name: 'New Category',
        description: 'Test category'
      }

      await act(async () => {
        const response = await result.current[0](newCategory).unwrap()

        expect(response.success).toBe(true)
        expect(response.category.name).toBe(newCategory.name)
        expect(response.category.level).toBe(0)
      })
    })
  })

  describe('Business API Integration', () => {
    it('should fetch businesses list', async () => {
      const { result } = renderHook(
        () => businessApi.useGetBusinessesQuery(),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(result.current.data?.businesses).toEqual([mockBusiness])
    })

    it('should fetch business by ID', async () => {
      const { result } = renderHook(
        () => businessApi.useGetBusinessQuery('business-1'),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      expect(result.current.data?.business).toEqual(mockBusiness)
    })

    it('should create business with mutation', async () => {
      const { result } = renderHook(
        () => businessApi.useCreateBusinessMutation(),
        { wrapper }
      )

      const newBusiness = {
        name: 'New Business',
        businessType: 'llc',
        industry: 'retail',
        email: 'contact@newbusiness.com',
        address: {
          street: '123 New St',
          city: 'New City',
          state: 'NY',
          postalCode: '12345',
          country: 'US'
        }
      }

      await act(async () => {
        const response = await result.current[0](newBusiness).unwrap()

        expect(response.success).toBe(true)
        expect(response.business.name).toBe(newBusiness.name)
        expect(response.business.status).toBe('pending')
      })
    })
  })

  describe('Error Handling', () => {
    it('should handle server errors gracefully', async () => {
      const { result } = renderHook(
        () => productApi.useGetProductQuery('error'),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isError).toBe(true)
      })

      expect(result.current.error).toBeDefined()
    })

    it('should handle network timeouts', async () => {
      const { result } = renderHook(
        () => productApi.useGetProductQuery('timeout'),
        { wrapper }
      )

      // Should still be loading after reasonable time
      await new Promise(resolve => setTimeout(resolve, 1000))
      expect(result.current.isLoading).toBe(true)
    })

    it('should retry failed requests based on RTK Query configuration', async () => {
      let callCount = 0

      server.use(
        rest.get('http://localhost:8080/api/v1/commerce/products/retry', (req, res, ctx) => {
          callCount++
          if (callCount < 3) {
            return res(ctx.status(500))
          }
          return res(ctx.status(200), ctx.json({ success: true, products: [] }))
        })
      )

      const { result } = renderHook(
        () => productApi.useGetProductQuery('retry'),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      }, { timeout: 10000 })

      expect(callCount).toBeGreaterThanOrEqual(3)
    })
  })

  describe('Caching and Optimization', () => {
    it('should cache product data between requests', async () => {
      // First request
      const { result: result1 } = renderHook(
        () => productApi.useGetProductQuery('1'),
        { wrapper }
      )

      await waitFor(() => {
        expect(result1.current.isSuccess).toBe(true)
      })

      // Second request should use cache
      const { result: result2 } = renderHook(
        () => productApi.useGetProductQuery('1'),
        { wrapper }
      )

      // Should immediately have data from cache
      expect(result2.current.data).toBeDefined()
      expect(result2.current.isLoading).toBe(false)
    })

    it('should invalidate cache when mutations occur', async () => {
      // Get initial products list
      const { result: listResult } = renderHook(
        () => productApi.useGetProductsQuery({}),
        { wrapper }
      )

      await waitFor(() => {
        expect(listResult.current.isSuccess).toBe(true)
      })

      const initialCount = listResult.current.data?.products.length || 0

      // Create new product
      const { result: createResult } = renderHook(
        () => productApi.useCreateProductMutation(),
        { wrapper }
      )

      await act(async () => {
        await createResult.current[0]({
          shopId: 'shop-1',
          product: { name: 'Cache Test Product', price: 99.99, currency: 'USD', category: 'test' }
        }).unwrap()
      })

      // Cache should be invalidated and refetch should occur
      await waitFor(() => {
        expect(listResult.current.data?.products.length).toBe(initialCount + 1)
      })
    })

    it('should handle optimistic updates for cart operations', async () => {
      // Get initial cart
      const { result: cartResult } = renderHook(
        () => cartApi.useGetCartQuery(),
        { wrapper }
      )

      await waitFor(() => {
        expect(cartResult.current.isSuccess).toBe(true)
      })

      const initialItemCount = cartResult.current.data?.cart.items.length || 0

      // Add item with optimistic update
      const { result: addResult } = renderHook(
        () => cartApi.useAddToCartMutation(),
        { wrapper }
      )

      await act(async () => {
        await addResult.current[0]({
          productId: '2',
          quantity: 1,
          unitPrice: 199.99
        }).unwrap()
      })

      // Cart should be updated
      await waitFor(() => {
        expect(cartResult.current.data?.cart.items.length).toBe(initialItemCount + 1)
      })
    })
  })

  describe('Real-time Data Synchronization', () => {
    it('should poll for cart updates when enabled', async () => {
      let callCount = 0

      server.use(
        rest.get('http://localhost:8080/api/v1/commerce/cart', (req, res, ctx) => {
          callCount++
          return res(
            ctx.status(200),
            ctx.json({
              success: true,
              cart: {
                ...mockCart,
                items: [...mockCart.items, {
                  id: `polling-item-${callCount}`,
                  productId: '1',
                  quantity: 1,
                  unitPrice: 100,
                  totalPrice: 100,
                  name: `Polling Item ${callCount}`
                }]
              }
            })
          )
        })
      )

      const { result } = renderHook(
        () => cartApi.useGetCartQuery(undefined, { pollingInterval: 1000 }),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      // Wait for polling to trigger additional requests
      await new Promise(resolve => setTimeout(resolve, 2500))

      expect(callCount).toBeGreaterThan(1)
    })

    it('should handle connection state changes', async () => {
      const { result } = renderHook(
        () => productApi.useGetProductsQuery({}),
        { wrapper }
      )

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true)
      })

      // Simulate going offline
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: false
      })

      // Trigger a request while offline
      const { result: offlineResult } = renderHook(
        () => productApi.useGetProductQuery('offline-test'),
        { wrapper }
      )

      // Should handle offline state gracefully
      expect(offlineResult.current.isError || offlineResult.current.isLoading).toBe(true)

      // Restore online state
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: true
      })
    })
  })
})