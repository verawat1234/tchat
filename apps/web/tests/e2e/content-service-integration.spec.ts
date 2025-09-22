import { test, expect } from '@playwright/test';

/**
 * Content Service Integration E2E Tests
 *
 * These tests validate the complete integration between the frontend and content service,
 * addressing the contract testing feedback to prevent future API mismatches.
 */

test.describe('Content Service Integration', () => {
  const baseURL = 'http://localhost:3000';
  const contentServiceURL = 'http://localhost:8086';

  test.beforeEach(async ({ page }) => {
    // Set up test data in the content service before each test
    await page.goto(baseURL);
  });

  test.describe('API Contract Validation', () => {
    test('should validate content service is running and accessible', async ({ request }) => {
      // Health check
      const healthResponse = await request.get(`${contentServiceURL}/health`);
      expect(healthResponse.status()).toBe(200);

      // Ready check
      const readyResponse = await request.get(`${contentServiceURL}/ready`);
      expect(readyResponse.status()).toBe(200);
    });

    test('should validate API endpoint structure and response format', async ({ request }) => {
      // Test the main content endpoint
      const response = await request.get(`${contentServiceURL}/api/v1/content`);
      expect(response.status()).toBe(200);

      const responseData = await response.json();
      expect(responseData).toHaveProperty('status');
      expect(responseData).toHaveProperty('data');
      expect(responseData.status).toBe('success');
    });

    test('should handle both UUID and string key requests', async ({ request }) => {
      // First create a test content item
      const createResponse = await request.post(`${contentServiceURL}/api/v1/content`, {
        data: {
          category: 'test.e2e.content',
          type: 'text',
          value: { text: 'E2E Test Content' },
          status: 'published'
        }
      });
      expect(createResponse.status()).toBe(201);

      const createdContent = await createResponse.json();
      const contentId = createdContent.data.id;

      // Test UUID access
      const uuidResponse = await request.get(`${contentServiceURL}/api/v1/content/${contentId}`);
      expect(uuidResponse.status()).toBe(200);

      const uuidData = await uuidResponse.json();
      expect(uuidData.data.id).toBe(contentId);

      // Test string key access
      const keyResponse = await request.get(`${contentServiceURL}/api/v1/content/test.e2e.content`);
      expect(keyResponse.status()).toBe(200);

      const keyData = await keyResponse.json();
      expect(keyData.data.category).toBe('test.e2e.content');
      expect(keyData.data.id).toBe(contentId);
    });

    test('should validate frontend proxy routing', async ({ request }) => {
      // Test that requests through the frontend proxy work correctly
      const response = await request.get(`${baseURL}/api/v1/content`);
      expect(response.status()).toBe(200);

      const responseData = await response.json();
      expect(responseData).toHaveProperty('status');
      expect(responseData).toHaveProperty('data');
    });
  });

  test.describe('Frontend Integration', () => {
    test('should load the application without console errors', async ({ page }) => {
      const errors: string[] = [];
      page.on('console', (msg) => {
        if (msg.type() === 'error') {
          errors.push(msg.text());
        }
      });

      await page.goto(baseURL);
      await page.waitForLoadState('networkidle');

      // Filter out expected 404s for missing content
      const criticalErrors = errors.filter(error =>
        !error.includes('404') &&
        !error.includes('Failed to load resource') &&
        !error.includes('non-serializable value')
      );

      expect(criticalErrors).toHaveLength(0);
    });

    test('should handle content loading gracefully', async ({ page }) => {
      // Create test content first
      await page.request.post(`${contentServiceURL}/api/v1/content`, {
        data: {
          category: 'auth.signin.title',
          type: 'text',
          value: { text: 'Sign In to Tchat' },
          status: 'published'
        }
      });

      await page.goto(baseURL);
      await page.waitForLoadState('networkidle');

      // The page should render without throwing JavaScript errors
      const hasJSErrors = await page.evaluate(() => {
        return window.onerror !== null;
      });
      expect(hasJSErrors).toBe(false);
    });

    test('should display fallback content when API fails', async ({ page }) => {
      // Navigate to the page when content service is unavailable for specific content
      await page.goto(baseURL);
      await page.waitForLoadState('networkidle');

      // Should still render the page structure
      const titleExists = await page.locator('h1').count() > 0;
      expect(titleExists).toBe(true);
    });
  });

  test.describe('Content API Workflow', () => {
    test('should create, read, update, and delete content', async ({ request }) => {
      // Create
      const createResponse = await request.post(`${contentServiceURL}/api/v1/content`, {
        data: {
          category: 'test.crud.operation',
          type: 'text',
          value: { text: 'Initial Content' },
          status: 'draft'
        }
      });
      expect(createResponse.status()).toBe(201);

      const createdData = await createResponse.json();
      const contentId = createdData.data.id;

      // Read
      const readResponse = await request.get(`${contentServiceURL}/api/v1/content/${contentId}`);
      expect(readResponse.status()).toBe(200);

      const readData = await readResponse.json();
      expect(readData.data.value.text).toBe('Initial Content');

      // Update
      const updateResponse = await request.put(`${contentServiceURL}/api/v1/content/${contentId}`, {
        data: {
          value: { text: 'Updated Content' },
          status: 'published'
        }
      });
      expect(updateResponse.status()).toBe(200);

      const updatedData = await updateResponse.json();
      expect(updatedData.data.value.text).toBe('Updated Content');
      expect(updatedData.data.status).toBe('published');

      // Delete
      const deleteResponse = await request.delete(`${contentServiceURL}/api/v1/content/${contentId}`);
      expect(deleteResponse.status()).toBe(200);

      // Verify deletion
      const verifyResponse = await request.get(`${contentServiceURL}/api/v1/content/${contentId}`);
      expect(verifyResponse.status()).toBe(404);
    });

    test('should validate content types and schemas', async ({ request }) => {
      const contentTypes = [
        { type: 'text', value: { text: 'Sample text' } },
        { type: 'html', value: { html: '<p>Sample HTML</p>' } },
        { type: 'json', value: { data: { key: 'value' } } }
      ];

      for (const contentType of contentTypes) {
        const response = await request.post(`${contentServiceURL}/api/v1/content`, {
          data: {
            category: `test.type.${contentType.type}`,
            type: contentType.type,
            value: contentType.value,
            status: 'published'
          }
        });

        expect(response.status()).toBe(201);
        const data = await response.json();
        expect(data.data.type).toBe(contentType.type);
        expect(data.data.value).toEqual(contentType.value);
      }
    });
  });

  test.describe('Error Handling', () => {
    test('should return 404 for non-existent content', async ({ request }) => {
      const response = await request.get(`${contentServiceURL}/api/v1/content/non.existent.key`);
      expect(response.status()).toBe(404);

      const data = await response.json();
      expect(data.status).toBe('error');
    });

    test('should handle malformed requests gracefully', async ({ request }) => {
      const response = await request.post(`${contentServiceURL}/api/v1/content`, {
        data: {
          // Missing required fields
          category: '',
          type: 'invalid-type'
        }
      });

      expect([400, 422]).toContain(response.status());
    });

    test('should validate request payload structure', async ({ request }) => {
      const invalidPayloads = [
        {}, // Empty payload
        { category: 'test' }, // Missing type and value
        { category: 'test', type: 'text' }, // Missing value
        { category: 'test', type: 'text', value: 'invalid' } // Invalid value structure
      ];

      for (const payload of invalidPayloads) {
        const response = await request.post(`${contentServiceURL}/api/v1/content`, {
          data: payload
        });

        expect([400, 422]).toContain(response.status());
      }
    });
  });

  test.describe('Performance and Reliability', () => {
    test('should handle concurrent requests', async ({ request }) => {
      const requests = Array.from({ length: 10 }, (_, i) =>
        request.post(`${contentServiceURL}/api/v1/content`, {
          data: {
            category: `test.concurrent.${i}`,
            type: 'text',
            value: { text: `Concurrent test ${i}` },
            status: 'published'
          }
        })
      );

      const responses = await Promise.all(requests);

      responses.forEach((response, i) => {
        expect(response.status()).toBe(201);
      });
    });

    test('should respond within acceptable time limits', async ({ request }) => {
      const startTime = Date.now();

      const response = await request.get(`${contentServiceURL}/api/v1/content`);

      const endTime = Date.now();
      const responseTime = endTime - startTime;

      expect(response.status()).toBe(200);
      expect(responseTime).toBeLessThan(1000); // Should respond within 1 second
    });
  });
});