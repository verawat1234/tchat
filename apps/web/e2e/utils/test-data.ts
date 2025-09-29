/**
 * Test Data Management for Commerce E2E Tests
 * Provides standardized test data and utilities for commerce testing
 */

import { faker } from '@faker-js/faker';

export interface TestUser {
  id?: string;
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  phone?: string;
  address?: TestAddress;
}

export interface TestAddress {
  street: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
}

export interface TestProduct {
  id?: string;
  name: string;
  description: string;
  price: number;
  category: string;
  sku: string;
  stock: number;
  images: string[];
  variants?: TestProductVariant[];
}

export interface TestProductVariant {
  id?: string;
  name: string;
  value: string;
  price?: number;
  stock?: number;
}

export interface TestCoupon {
  code: string;
  type: 'percentage' | 'fixed';
  value: number;
  minAmount?: number;
  maxDiscount?: number;
  expiresAt?: Date;
}

export interface TestCart {
  id?: string;
  userId?: string;
  items: TestCartItem[];
  coupon?: TestCoupon;
  totalAmount: number;
  currency: string;
}

export interface TestCartItem {
  productId: string;
  variantId?: string;
  quantity: number;
  price: number;
}

export class TestDataGenerator {
  /**
   * Generate a test user with realistic data
   */
  static generateUser(overrides: Partial<TestUser> = {}): TestUser {
    const firstName = faker.person.firstName();
    const lastName = faker.person.lastName();

    return {
      email: faker.internet.email({ firstName, lastName }).toLowerCase(),
      password: 'TestPassword123!',
      firstName,
      lastName,
      phone: faker.phone.number(),
      address: this.generateAddress(),
      ...overrides,
    };
  }

  /**
   * Generate a test address
   */
  static generateAddress(overrides: Partial<TestAddress> = {}): TestAddress {
    return {
      street: faker.location.streetAddress(),
      city: faker.location.city(),
      state: faker.location.state(),
      zipCode: faker.location.zipCode(),
      country: 'United States',
      ...overrides,
    };
  }

  /**
   * Generate a test product
   */
  static generateProduct(overrides: Partial<TestProduct> = {}): TestProduct {
    const productName = faker.commerce.productName();
    const price = parseFloat(faker.commerce.price({ min: 10, max: 500 }));

    return {
      name: productName,
      description: faker.commerce.productDescription(),
      price,
      category: faker.commerce.department(),
      sku: faker.string.alphanumeric({ length: 8 }).toUpperCase(),
      stock: faker.number.int({ min: 0, max: 100 }),
      images: [
        faker.image.url({ width: 400, height: 400 }),
        faker.image.url({ width: 400, height: 400 }),
      ],
      variants: this.generateProductVariants(),
      ...overrides,
    };
  }

  /**
   * Generate product variants (size, color, etc.)
   */
  static generateProductVariants(): TestProductVariant[] {
    const colors = ['Red', 'Blue', 'Green', 'Black', 'White'];
    const sizes = ['S', 'M', 'L', 'XL'];

    const variants: TestProductVariant[] = [];

    // Add color variants
    colors.slice(0, faker.number.int({ min: 2, max: 4 })).forEach(color => {
      variants.push({
        name: 'Color',
        value: color,
        stock: faker.number.int({ min: 0, max: 20 }),
      });
    });

    // Add size variants
    sizes.slice(0, faker.number.int({ min: 2, max: 4 })).forEach(size => {
      variants.push({
        name: 'Size',
        value: size,
        stock: faker.number.int({ min: 0, max: 20 }),
      });
    });

    return variants;
  }

  /**
   * Generate a test coupon
   */
  static generateCoupon(overrides: Partial<TestCoupon> = {}): TestCoupon {
    const type = faker.helpers.arrayElement(['percentage', 'fixed'] as const);

    return {
      code: faker.string.alphanumeric({ length: 8 }).toUpperCase(),
      type,
      value: type === 'percentage'
        ? faker.number.int({ min: 5, max: 50 })
        : parseFloat(faker.commerce.price({ min: 5, max: 100 })),
      minAmount: parseFloat(faker.commerce.price({ min: 50, max: 200 })),
      maxDiscount: type === 'percentage'
        ? parseFloat(faker.commerce.price({ min: 10, max: 100 }))
        : undefined,
      expiresAt: faker.date.future(),
      ...overrides,
    };
  }

  /**
   * Generate a test cart with items
   */
  static generateCart(products: TestProduct[], overrides: Partial<TestCart> = {}): TestCart {
    const itemCount = faker.number.int({ min: 1, max: 5 });
    const selectedProducts = faker.helpers.arrayElements(products, itemCount);

    const items: TestCartItem[] = selectedProducts.map(product => ({
      productId: product.id!,
      variantId: product.variants?.[0]?.id,
      quantity: faker.number.int({ min: 1, max: 3 }),
      price: product.price,
    }));

    const totalAmount = items.reduce((sum, item) => sum + (item.price * item.quantity), 0);

    return {
      items,
      totalAmount,
      currency: 'USD',
      ...overrides,
    };
  }
}

/**
 * Predefined test data sets for consistent testing
 */
export const TestDataSets = {
  /**
   * Standard test users for different scenarios
   */
  users: {
    new: TestDataGenerator.generateUser({
      email: 'new-user@test.com',
      firstName: 'New',
      lastName: 'User',
    }),
    existing: TestDataGenerator.generateUser({
      email: 'existing-user@test.com',
      firstName: 'Existing',
      lastName: 'User',
    }),
    premium: TestDataGenerator.generateUser({
      email: 'premium-user@test.com',
      firstName: 'Premium',
      lastName: 'User',
    }),
  },

  /**
   * Standard test products for different categories
   */
  products: {
    electronics: TestDataGenerator.generateProduct({
      name: 'Test Smartphone',
      category: 'Electronics',
      price: 299.99,
      sku: 'PHONE001',
    }),
    clothing: TestDataGenerator.generateProduct({
      name: 'Test T-Shirt',
      category: 'Clothing',
      price: 29.99,
      sku: 'SHIRT001',
    }),
    books: TestDataGenerator.generateProduct({
      name: 'Test Book',
      category: 'Books',
      price: 19.99,
      sku: 'BOOK001',
    }),
    outOfStock: TestDataGenerator.generateProduct({
      name: 'Out of Stock Item',
      category: 'Test',
      price: 49.99,
      sku: 'STOCK000',
      stock: 0,
    }),
  },

  /**
   * Standard test coupons
   */
  coupons: {
    percentage: TestDataGenerator.generateCoupon({
      code: 'SAVE20',
      type: 'percentage',
      value: 20,
      minAmount: 100,
      maxDiscount: 50,
    }),
    fixed: TestDataGenerator.generateCoupon({
      code: 'FIXED10',
      type: 'fixed',
      value: 10,
      minAmount: 50,
    }),
    expired: TestDataGenerator.generateCoupon({
      code: 'EXPIRED',
      type: 'percentage',
      value: 15,
      expiresAt: new Date('2020-01-01'),
    }),
  },
};

/**
 * Test data cleanup utilities
 */
export class TestDataCleanup {
  private static createdUsers: string[] = [];
  private static createdProducts: string[] = [];
  private static createdCarts: string[] = [];

  /**
   * Track created entities for cleanup
   */
  static trackUser(userId: string) {
    this.createdUsers.push(userId);
  }

  static trackProduct(productId: string) {
    this.createdProducts.push(productId);
  }

  static trackCart(cartId: string) {
    this.createdCarts.push(cartId);
  }

  /**
   * Clean up all tracked entities
   */
  static async cleanupAll() {
    // In a real implementation, these would make API calls to clean up
    console.log('Cleaning up test data:', {
      users: this.createdUsers.length,
      products: this.createdProducts.length,
      carts: this.createdCarts.length,
    });

    // Clear tracking arrays
    this.createdUsers = [];
    this.createdProducts = [];
    this.createdCarts = [];
  }

  /**
   * Reset all test data to initial state
   */
  static async resetTestData() {
    await this.cleanupAll();
    // Additional reset logic would go here
  }
}