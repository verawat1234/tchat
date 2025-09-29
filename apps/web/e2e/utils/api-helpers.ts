/**
 * API Helper Utilities for E2E Testing
 * Provides utilities for interacting with backend services during tests
 */

import { APIRequestContext, expect } from '@playwright/test';
import { TestUser, TestProduct, TestCart, TestCoupon } from './test-data';

export interface ApiResponse<T = any> {
  data?: T;
  error?: string;
  status: number;
}

export class ApiHelpers {
  constructor(private request: APIRequestContext) {}

  /**
   * Authentication helpers
   */
  async login(email: string, password: string): Promise<{ token: string; user: any }> {
    const response = await this.request.post('/api/v1/auth/login', {
      data: { email, password }
    });

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result;
  }

  async logout(token: string): Promise<void> {
    const response = await this.request.post('/api/v1/auth/logout', {
      headers: { Authorization: `Bearer ${token}` }
    });

    expect(response.status()).toBe(200);
  }

  async createUser(user: TestUser): Promise<{ id: string; token: string }> {
    const response = await this.request.post('/api/v1/auth/register', {
      data: user
    });

    expect(response.status()).toBe(201);
    const result = await response.json();
    return result;
  }

  /**
   * Product management helpers
   */
  async createProduct(product: TestProduct, token: string): Promise<{ id: string }> {
    const response = await this.request.post('/api/v1/commerce/products', {
      headers: { Authorization: `Bearer ${token}` },
      data: product
    });

    expect(response.status()).toBe(201);
    const result = await response.json();
    return result;
  }

  async getProduct(productId: string): Promise<TestProduct> {
    const response = await this.request.get(`/api/v1/commerce/products/${productId}`);

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  async getProducts(params?: {
    category?: string;
    search?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ products: TestProduct[]; total: number }> {
    let url = '/api/v1/commerce/products';

    if (params) {
      const searchParams = new URLSearchParams();
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          searchParams.append(key, value.toString());
        }
      });
      url += `?${searchParams.toString()}`;
    }

    const response = await this.request.get(url);
    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  async updateProductStock(productId: string, stock: number, token: string): Promise<void> {
    const response = await this.request.patch(`/api/v1/commerce/products/${productId}`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { stock }
    });

    expect(response.status()).toBe(200);
  }

  /**
   * Cart management helpers
   */
  async getCart(token: string): Promise<TestCart | null> {
    const response = await this.request.get('/api/v1/commerce/cart', {
      headers: { Authorization: `Bearer ${token}` }
    });

    if (response.status() === 404) {
      return null;
    }

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  async addToCart(productId: string, quantity: number, variantId?: string, token?: string): Promise<TestCart> {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};

    const response = await this.request.post('/api/v1/commerce/cart/items', {
      headers,
      data: { productId, quantity, variantId }
    });

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  async updateCartItem(itemId: string, quantity: number, token?: string): Promise<TestCart> {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};

    const response = await this.request.patch(`/api/v1/commerce/cart/items/${itemId}`, {
      headers,
      data: { quantity }
    });

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  async removeFromCart(itemId: string, token?: string): Promise<TestCart> {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};

    const response = await this.request.delete(`/api/v1/commerce/cart/items/${itemId}`, {
      headers
    });

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  async clearCart(token?: string): Promise<void> {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};

    const response = await this.request.delete('/api/v1/commerce/cart', {
      headers
    });

    expect(response.status()).toBe(200);
  }

  async applyCoupon(couponCode: string, token?: string): Promise<TestCart> {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};

    const response = await this.request.post('/api/v1/commerce/cart/coupon', {
      headers,
      data: { code: couponCode }
    });

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  async removeCoupon(token?: string): Promise<TestCart> {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};

    const response = await this.request.delete('/api/v1/commerce/cart/coupon', {
      headers
    });

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  /**
   * Coupon management helpers
   */
  async createCoupon(coupon: TestCoupon, token: string): Promise<{ id: string }> {
    const response = await this.request.post('/api/v1/commerce/coupons', {
      headers: { Authorization: `Bearer ${token}` },
      data: coupon
    });

    expect(response.status()).toBe(201);
    const result = await response.json();
    return result;
  }

  async validateCoupon(couponCode: string, cartTotal: number): Promise<{ valid: boolean; discount: number }> {
    const response = await this.request.post('/api/v1/commerce/coupons/validate', {
      data: { code: couponCode, cartTotal }
    });

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  /**
   * Order management helpers
   */
  async createOrder(token: string, paymentMethod?: string): Promise<{ id: string; status: string }> {
    const response = await this.request.post('/api/v1/commerce/orders', {
      headers: { Authorization: `Bearer ${token}` },
      data: { paymentMethod: paymentMethod || 'stripe' }
    });

    expect(response.status()).toBe(201);
    const result = await response.json();
    return result;
  }

  async getOrder(orderId: string, token: string): Promise<any> {
    const response = await this.request.get(`/api/v1/commerce/orders/${orderId}`, {
      headers: { Authorization: `Bearer ${token}` }
    });

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  /**
   * Category management helpers
   */
  async getCategories(): Promise<string[]> {
    const response = await this.request.get('/api/v1/commerce/categories');

    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  /**
   * Search helpers
   */
  async searchProducts(query: string, filters?: {
    category?: string;
    minPrice?: number;
    maxPrice?: number;
    inStock?: boolean;
  }): Promise<{ products: TestProduct[]; total: number }> {
    const params = new URLSearchParams({ q: query });

    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined) {
          params.append(key, value.toString());
        }
      });
    }

    const response = await this.request.get(`/api/v1/commerce/search?${params.toString()}`);
    expect(response.status()).toBe(200);
    const result = await response.json();
    return result.data;
  }

  /**
   * Test data setup helpers
   */
  async setupTestProducts(count: number = 10, token: string): Promise<TestProduct[]> {
    const { TestDataGenerator } = await import('./test-data');
    const products: TestProduct[] = [];

    for (let i = 0; i < count; i++) {
      const product = TestDataGenerator.generateProduct();
      const created = await this.createProduct(product, token);
      products.push({ ...product, id: created.id });
    }

    return products;
  }

  async setupTestCategories(token: string): Promise<string[]> {
    // In a real implementation, this would create test categories
    return ['Electronics', 'Clothing', 'Books', 'Home & Garden', 'Sports'];
  }

  async setupTestCoupons(token: string): Promise<TestCoupon[]> {
    const { TestDataSets } = await import('./test-data');
    const coupons: TestCoupon[] = [];

    for (const coupon of Object.values(TestDataSets.coupons)) {
      const created = await this.createCoupon(coupon, token);
      coupons.push({ ...coupon, id: created.id });
    }

    return coupons;
  }

  /**
   * Performance monitoring helpers
   */
  async measureApiPerformance<T>(
    operation: () => Promise<T>,
    threshold: number = 1000
  ): Promise<{ result: T; duration: number }> {
    const startTime = Date.now();
    const result = await operation();
    const duration = Date.now() - startTime;

    expect(duration).toBeLessThan(threshold);

    return { result, duration };
  }

  /**
   * Cleanup helpers
   */
  async cleanupTestData(token: string): Promise<void> {
    // Clean up test products, coupons, etc.
    // This would typically involve calling admin endpoints to clean up test data
    console.log('Cleaning up test data via API');
  }
}

/**
 * Factory function to create API helpers with request context
 */
export function createApiHelpers(request: APIRequestContext): ApiHelpers {
  return new ApiHelpers(request);
}