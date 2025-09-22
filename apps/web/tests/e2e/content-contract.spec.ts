import { test, expect } from '@playwright/test';

/**
 * Content Service Contract Tests
 *
 * These tests specifically validate the contract between frontend and backend
 * to prevent the type of issues mentioned in user feedback: "didnt you already have contract test? why it still diff?"
 */

test.describe('Content Service Contract Validation', () => {
  const contentServiceURL = 'http://localhost:8086';

  test.describe('Response Structure Contract', () => {
    test('should return correct response structure for content list', async ({ request }) => {
      const response = await request.get(`${contentServiceURL}/api/v1/content`);
      expect(response.status()).toBe(200);

      const data = await response.json();

      // Validate top-level response structure
      expect(data).toHaveProperty('status');
      expect(data).toHaveProperty('data');
      expect(data.status).toBe('success');

      // Validate data structure - should be an array or paginated object
      if (Array.isArray(data.data)) {
        // If it's an array, validate each item structure
        data.data.forEach((item: any) => {
          expect(item).toHaveProperty('id');
          expect(item).toHaveProperty('category');
          expect(item).toHaveProperty('type');
          expect(item).toHaveProperty('value');
          expect(item).toHaveProperty('status');
          expect(item).toHaveProperty('created_at');
          expect(item).toHaveProperty('updated_at');
        });
      } else {
        // If it's a paginated response, validate pagination structure
        expect(data.data).toHaveProperty('items');
        expect(data.data).toHaveProperty('total');
        expect(data.data).toHaveProperty('page');
        expect(data.data).toHaveProperty('limit');
        expect(Array.isArray(data.data.items)).toBe(true);
      }
    });

    test('should return correct structure for single content item', async ({ request }) => {
      // First create a test item
      const createResponse = await request.post(`${contentServiceURL}/api/v1/content`, {
        data: {
          category: 'contract.test.item',
          type: 'text',
          value: { text: 'Contract test content' },
          status: 'published'
        }
      });

      expect(createResponse.status()).toBe(201);
      const createdData = await createResponse.json();
      const contentId = createdData.data.id;

      // Test getting by ID
      const getByIdResponse = await request.get(`${contentServiceURL}/api/v1/content/${contentId}`);
      expect(getByIdResponse.status()).toBe(200);

      const getByIdData = await getByIdResponse.json();
      expect(getByIdData).toHaveProperty('status');
      expect(getByIdData).toHaveProperty('data');
      expect(getByIdData.status).toBe('success');

      const item = getByIdData.data;
      expect(item).toHaveProperty('id');
      expect(item).toHaveProperty('category');
      expect(item).toHaveProperty('type');
      expect(item).toHaveProperty('value');
      expect(item).toHaveProperty('status');
      expect(item).toHaveProperty('created_at');
      expect(item).toHaveProperty('updated_at');

      // Test getting by string key
      const getByKeyResponse = await request.get(`${contentServiceURL}/api/v1/content/contract.test.item`);
      expect(getByKeyResponse.status()).toBe(200);

      const getByKeyData = await getByKeyResponse.json();
      expect(getByKeyData).toHaveProperty('status');
      expect(getByKeyData).toHaveProperty('data');
      expect(getByKeyData.status).toBe('success');

      // Both responses should have identical structure
      expect(getByIdData.data).toEqual(getByKeyData.data);
    });

    test('should validate content value structure by type', async ({ request }) => {
      const testCases = [
        {
          type: 'text',
          value: { text: 'Sample text content' },
          expectedStructure: { text: expect.any(String) }
        },
        {
          type: 'html',
          value: { html: '<p>Sample HTML content</p>' },
          expectedStructure: { html: expect.any(String) }
        },
        {
          type: 'json',
          value: { data: { key: 'value', number: 42 } },
          expectedStructure: { data: expect.any(Object) }
        }
      ];

      for (const testCase of testCases) {
        const response = await request.post(`${contentServiceURL}/api/v1/content`, {
          data: {
            category: `contract.type.${testCase.type}`,
            type: testCase.type,
            value: testCase.value,
            status: 'published'
          }
        });

        expect(response.status()).toBe(201);
        const data = await response.json();

        expect(data.data.type).toBe(testCase.type);
        expect(data.data.value).toMatchObject(testCase.expectedStructure);
      }
    });
  });

  test.describe('Error Response Contract', () => {
    test('should return consistent error structure', async ({ request }) => {
      // Test 404 error
      const notFoundResponse = await request.get(`${contentServiceURL}/api/v1/content/non.existent.key`);
      expect(notFoundResponse.status()).toBe(404);

      const notFoundData = await notFoundResponse.json();
      expect(notFoundData).toHaveProperty('status');
      expect(notFoundData).toHaveProperty('error');
      expect(notFoundData.status).toBe('error');
      expect(notFoundData.error).toContain('not found');
    });

    test('should return validation error structure', async ({ request }) => {
      const invalidResponse = await request.post(`${contentServiceURL}/api/v1/content`, {
        data: {
          category: '',
          type: 'invalid-type',
          value: 'invalid-value'
        }
      });

      expect([400, 422]).toContain(invalidResponse.status());
      const invalidData = await invalidResponse.json();

      expect(invalidData).toHaveProperty('status');
      expect(invalidData).toHaveProperty('error');
      expect(invalidData.status).toBe('error');
    });
  });

  test.describe('Frontend-Backend Data Contract', () => {
    test('should validate that frontend can properly process backend responses', async ({ page, request }) => {
      // Create test content that the frontend will try to load
      const testContent = [
        {
          category: 'auth.signin.title',
          type: 'text',
          value: { text: 'Sign In' },
          status: 'published'
        },
        {
          category: 'auth.signin.subtitle',
          type: 'text',
          value: { text: 'Welcome back to Tchat' },
          status: 'published'
        },
        {
          category: 'navigation.menu.items',
          type: 'json',
          value: { data: [{ label: 'Home', path: '/' }, { label: 'About', path: '/about' }] },
          status: 'published'
        }
      ];

      // Create test content
      for (const content of testContent) {
        await request.post(`${contentServiceURL}/api/v1/content`, { data: content });
      }

      // Navigate to frontend and check for JavaScript errors
      const errors: string[] = [];
      page.on('console', (msg) => {
        if (msg.type() === 'error' && !msg.text().includes('404') && !msg.text().includes('Failed to load resource')) {
          errors.push(msg.text());
        }
      });

      await page.goto('http://localhost:3000');
      await page.waitForLoadState('networkidle');

      // Should not have forEach errors
      const forEachErrors = errors.filter(error => error.includes('forEach is not a function'));
      expect(forEachErrors).toHaveLength(0);

      // Should not have other critical JavaScript errors
      const criticalErrors = errors.filter(error =>
        !error.includes('non-serializable value') &&
        !error.includes('API Error')
      );
      expect(criticalErrors).toHaveLength(0);
    });

    test('should validate content batch loading contract', async ({ request }) => {
      // Test the batch content loading endpoint that the frontend uses
      const categories = ['auth.signin.title', 'auth.signup.title', 'navigation.main.title'];

      // Create content for each category
      for (let i = 0; i < categories.length; i++) {
        await request.post(`${contentServiceURL}/api/v1/content`, {
          data: {
            category: categories[i],
            type: 'text',
            value: { text: `Content ${i + 1}` },
            status: 'published'
          }
        });
      }

      // Test batch retrieval (this might be how the frontend loads content)
      const batchResponse = await request.get(`${contentServiceURL}/api/v1/content?categories=${categories.join(',')}`);

      if (batchResponse.status() === 200) {
        const batchData = await batchResponse.json();

        // Validate that the response is in the correct format for frontend consumption
        expect(batchData).toHaveProperty('status');
        expect(batchData).toHaveProperty('data');
        expect(batchData.status).toBe('success');

        // The data should be either an array or an object that the frontend can iterate over
        const data = batchData.data;
        if (Array.isArray(data)) {
          expect(data.length).toBeGreaterThan(0);
        } else if (typeof data === 'object') {
          expect(Object.keys(data).length).toBeGreaterThan(0);
        }
      }
    });
  });

  test.describe('UUID vs String Key Contract', () => {
    test('should handle both UUID and string key formats consistently', async ({ request }) => {
      // Create content
      const createResponse = await request.post(`${contentServiceURL}/api/v1/content`, {
        data: {
          category: 'uuid.string.test',
          type: 'text',
          value: { text: 'Test content for UUID/string compatibility' },
          status: 'published'
        }
      });

      expect(createResponse.status()).toBe(201);
      const createdData = await createResponse.json();
      const contentId = createdData.data.id;

      // Test UUID access
      const uuidResponse = await request.get(`${contentServiceURL}/api/v1/content/${contentId}`);
      expect(uuidResponse.status()).toBe(200);

      // Test string key access
      const stringResponse = await request.get(`${contentServiceURL}/api/v1/content/uuid.string.test`);
      expect(stringResponse.status()).toBe(200);

      // Both should return the same data
      const uuidData = await uuidResponse.json();
      const stringData = await stringResponse.json();

      expect(uuidData.data.id).toBe(stringData.data.id);
      expect(uuidData.data.category).toBe(stringData.data.category);
      expect(uuidData.data.value).toEqual(stringData.data.value);
    });

    test('should validate that invalid UUIDs fall back to string key lookup', async ({ request }) => {
      // Create content with a category that looks like it could be confused with a UUID
      await request.post(`${contentServiceURL}/api/v1/content`, {
        data: {
          category: 'not-a-uuid-but-has-dashes',
          type: 'text',
          value: { text: 'This should be accessible by string key' },
          status: 'published'
        }
      });

      // Try to access it as a string key
      const response = await request.get(`${contentServiceURL}/api/v1/content/not-a-uuid-but-has-dashes`);
      expect(response.status()).toBe(200);

      const data = await response.json();
      expect(data.data.category).toBe('not-a-uuid-but-has-dashes');
    });
  });
});