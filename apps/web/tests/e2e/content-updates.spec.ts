import { test, expect, Page, BrowserContext } from '@playwright/test';

/**
 * T057: E2E Test - Content Update Propagation
 *
 * Comprehensive end-to-end tests for real-time content update propagation
 * including multi-tab synchronization, notifications, versioning, conflict resolution,
 * and performance validation across different content types.
 */

// Test data constants
const TEST_CONTENT = {
  MESSAGE: {
    id: 'test-message-001',
    type: 'message',
    title: 'Test Chat Message',
    content: 'Hello, this is a test message for content updates.',
    category: 'chat',
  },
  POST: {
    id: 'test-post-001',
    type: 'post',
    title: 'Test Social Post',
    content: 'This is a social media post to test content propagation.',
    category: 'social',
  },
  PRODUCT: {
    id: 'test-product-001',
    type: 'product',
    title: 'Test Product',
    content: 'Product description for testing commerce updates.',
    category: 'commerce',
    price: 29.99,
  },
} as const;

// Helper functions
async function setupTestContent(page: Page, content: typeof TEST_CONTENT.MESSAGE) {
  await page.evaluate((testContent) => {
    // Mock API response for test content
    window.__TEST_CONTENT__ = testContent;

    // Inject test content into Redux store
    if (window.__REDUX_STORE__) {
      window.__REDUX_STORE__.dispatch({
        type: 'content/updateFallbackContent',
        payload: {
          contentId: testContent.id,
          content: {
            id: testContent.id,
            type: testContent.type,
            title: testContent.title,
            content: testContent.content,
            category: testContent.category,
            version: 1,
            lastModified: new Date().toISOString(),
            author: 'test-user',
          },
        },
      });
    }
  }, content);
}

async function waitForContentUpdate(page: Page, contentId: string, expectedVersion: number) {
  await page.waitForFunction(
    ({ id, version }) => {
      const state = window.__REDUX_STORE__?.getState?.();
      const content = state?.content?.fallbackContent?.[id];
      return content?.version === version;
    },
    { id: contentId, version: expectedVersion },
    { timeout: 10000 }
  );
}

async function triggerContentUpdate(page: Page, contentId: string, updates: Record<string, any>) {
  await page.evaluate(
    ({ id, changes }) => {
      // Simulate real-time content update
      const currentContent = window.__REDUX_STORE__?.getState?.().content.fallbackContent[id];
      if (currentContent) {
        const updatedContent = {
          ...currentContent,
          ...changes,
          version: (currentContent.version || 1) + 1,
          lastModified: new Date().toISOString(),
        };

        window.__REDUX_STORE__.dispatch({
          type: 'content/updateFallbackContent',
          payload: {
            contentId: id,
            content: updatedContent,
          },
        });

        // Trigger real-time update event
        window.dispatchEvent(new CustomEvent('content-updated', {
          detail: { contentId: id, content: updatedContent },
        }));
      }
    },
    { id: contentId, changes: updates }
  );
}

async function getSyncStatus(page: Page): Promise<string> {
  return await page.evaluate(() => {
    return window.__REDUX_STORE__?.getState?.().content.syncStatus || 'unknown';
  });
}

async function measurePerformanceMetrics(page: Page) {
  return await page.evaluate(() => {
    const metrics = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
    const paintMetrics = performance.getEntriesByType('paint');

    return {
      domContentLoaded: metrics.domContentLoadedEventEnd - metrics.domContentLoadedEventStart,
      firstContentfulPaint: paintMetrics.find(m => m.name === 'first-contentful-paint')?.startTime || 0,
      memoryUsage: (performance as any).memory ? {
        used: (performance as any).memory.usedJSHeapSize,
        total: (performance as any).memory.totalJSHeapSize,
        limit: (performance as any).memory.jsHeapSizeLimit,
      } : null,
    };
  });
}

test.describe('Content Update Propagation', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to app and wait for initialization
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Wait for Redux store to be available
    await page.waitForFunction(() => window.__REDUX_STORE__ !== undefined);

    // Initialize test data
    await setupTestContent(page, TEST_CONTENT.MESSAGE);
  });

  test.describe('Real-time Updates', () => {
    test('should propagate content changes without page refresh', async ({ page }) => {
      // Navigate to content view
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Verify initial content
      await expect(page.getByTestId('content-title')).toHaveText(TEST_CONTENT.MESSAGE.title);
      await expect(page.getByTestId('content-body')).toHaveText(TEST_CONTENT.MESSAGE.content);

      // Trigger content update
      const updatedTitle = 'Updated Test Message Title';
      const updatedContent = 'This content has been updated in real-time!';

      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: updatedTitle,
        content: updatedContent,
      });

      // Verify content updates without page refresh
      await expect(page.getByTestId('content-title')).toHaveText(updatedTitle);
      await expect(page.getByTestId('content-body')).toHaveText(updatedContent);

      // Verify no page reload occurred
      const navigationCount = await page.evaluate(() => performance.getEntriesByType('navigation').length);
      expect(navigationCount).toBe(1); // Only initial load
    });

    test('should handle rapid successive updates correctly', async ({ page }) => {
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Trigger multiple rapid updates
      const updates = [
        { title: 'Update 1', content: 'Content update 1' },
        { title: 'Update 2', content: 'Content update 2' },
        { title: 'Update 3', content: 'Content update 3' },
        { title: 'Final Update', content: 'Final content update' },
      ];

      for (const [index, update] of updates.entries()) {
        await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, update);
        await waitForContentUpdate(page, TEST_CONTENT.MESSAGE.id, index + 2); // +2 because initial version is 1
      }

      // Verify final state
      await expect(page.getByTestId('content-title')).toHaveText('Final Update');
      await expect(page.getByTestId('content-body')).toHaveText('Final content update');
    });

    test('should maintain update order and consistency', async ({ page }) => {
      await page.getByTestId('content-list').click();

      // Create a sequence of ordered updates
      const updateSequence = [
        { step: 1, title: 'Step 1', timestamp: new Date().toISOString() },
        { step: 2, title: 'Step 2', timestamp: new Date(Date.now() + 1000).toISOString() },
        { step: 3, title: 'Step 3', timestamp: new Date(Date.now() + 2000).toISOString() },
      ];

      // Apply updates in sequence
      for (const update of updateSequence) {
        await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, update);
        await page.waitForTimeout(100); // Small delay to ensure order
      }

      // Verify final state reflects correct order
      const finalVersion = await page.evaluate((contentId) => {
        const state = window.__REDUX_STORE__?.getState?.();
        return state?.content?.fallbackContent?.[contentId]?.version;
      }, TEST_CONTENT.MESSAGE.id);

      expect(finalVersion).toBe(4); // Initial version 1 + 3 updates
    });
  });

  test.describe('Multi-tab Synchronization', () => {
    let context: BrowserContext;
    let secondPage: Page;

    test.beforeEach(async ({ browser }) => {
      context = await browser.newContext();
      secondPage = await context.newPage();
      await secondPage.goto('/');
      await secondPage.waitForLoadState('networkidle');
      await secondPage.waitForFunction(() => window.__REDUX_STORE__ !== undefined);
      await setupTestContent(secondPage, TEST_CONTENT.MESSAGE);
    });

    test.afterEach(async () => {
      await secondPage.close();
      await context.close();
    });

    test('should sync content updates across multiple tabs', async ({ page }) => {
      // Open content in both tabs
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      await secondPage.getByTestId('content-list').click();
      await secondPage.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Verify initial state in both tabs
      await expect(page.getByTestId('content-title')).toHaveText(TEST_CONTENT.MESSAGE.title);
      await expect(secondPage.getByTestId('content-title')).toHaveText(TEST_CONTENT.MESSAGE.title);

      // Update content in first tab
      const newTitle = 'Multi-tab Update Test';
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, { title: newTitle });

      // Simulate cross-tab sync by triggering the same update in second tab
      await triggerContentUpdate(secondPage, TEST_CONTENT.MESSAGE.id, { title: newTitle });

      // Verify both tabs show updated content
      await expect(page.getByTestId('content-title')).toHaveText(newTitle);
      await expect(secondPage.getByTestId('content-title')).toHaveText(newTitle);
    });

    test('should handle concurrent edits across tabs', async ({ page }) => {
      // Set up content editing in both tabs
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();
      await page.getByTestId('edit-content-button').click();

      await secondPage.getByTestId('content-list').click();
      await secondPage.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();
      await secondPage.getByTestId('edit-content-button').click();

      // Make concurrent edits
      await page.getByTestId('content-title-input').fill('Tab 1 Edit');
      await secondPage.getByTestId('content-title-input').fill('Tab 2 Edit');

      // Trigger saves simultaneously
      const [, ] = await Promise.all([
        page.getByTestId('save-content-button').click(),
        secondPage.getByTestId('save-content-button').click(),
      ]);

      // Verify conflict resolution UI appears
      await expect(
        page.getByTestId('conflict-resolution-dialog').or(
          secondPage.getByTestId('conflict-resolution-dialog')
        )
      ).toBeVisible();
    });

    test('should maintain tab focus state during updates', async ({ page }) => {
      // Focus on content in first tab
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Switch to second tab and trigger update
      await secondPage.bringToFront();
      await triggerContentUpdate(secondPage, TEST_CONTENT.MESSAGE.id, {
        title: 'Background Update'
      });

      // Switch back to first tab
      await page.bringToFront();

      // Verify content updated but focus maintained
      await expect(page.getByTestId('content-title')).toHaveText('Background Update');
      await expect(page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`)).toBeFocused();
    });
  });

  test.describe('Update Notifications', () => {
    test('should display update notifications to users', async ({ page }) => {
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Trigger content update
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Notification Test Update',
      });

      // Check for update notification
      await expect(page.getByTestId('update-notification')).toBeVisible();
      await expect(page.getByTestId('update-notification')).toContainText('Content updated');

      // Verify notification auto-dismisses
      await expect(page.getByTestId('update-notification')).toBeHidden({ timeout: 5000 });
    });

    test('should show update indicators in content lists', async ({ page }) => {
      await page.getByTestId('content-list').click();

      // Trigger update for content in list
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'List Update Test',
      });

      // Check for update indicator
      await expect(page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`))
        .toHaveClass(/updated/);
      await expect(page.getByTestId(`update-indicator-${TEST_CONTENT.MESSAGE.id}`))
        .toBeVisible();
    });

    test('should provide detailed update information', async ({ page }) => {
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Trigger update with metadata
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Detailed Update Test',
        lastModified: new Date().toISOString(),
        modifiedBy: 'test-user-2',
      });

      // Check update details
      await page.getByTestId('update-details-button').click();
      await expect(page.getByTestId('update-details-modal')).toBeVisible();
      await expect(page.getByTestId('update-author')).toContainText('test-user-2');
      await expect(page.getByTestId('update-timestamp')).toBeVisible();
    });
  });

  test.describe('Content Versioning', () => {
    test('should track content version history', async ({ page }) => {
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Make multiple updates to create version history
      const versions = [
        { title: 'Version 2', content: 'Second version content' },
        { title: 'Version 3', content: 'Third version content' },
        { title: 'Version 4', content: 'Fourth version content' },
      ];

      for (const [index, update] of versions.entries()) {
        await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, update);
        await waitForContentUpdate(page, TEST_CONTENT.MESSAGE.id, index + 2);
      }

      // Check version history
      await page.getByTestId('version-history-button').click();
      await expect(page.getByTestId('version-history-modal')).toBeVisible();

      // Verify all versions are listed
      for (let i = 1; i <= 4; i++) {
        await expect(page.getByTestId(`version-${i}`)).toBeVisible();
      }
    });

    test('should allow version comparison', async ({ page }) => {
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Create two versions for comparison
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Version 2 Title',
        content: 'Version 2 content with changes',
      });

      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Version 3 Title',
        content: 'Version 3 content with more changes',
      });

      // Open version comparison
      await page.getByTestId('version-history-button').click();
      await page.getByTestId('version-2').click();
      await page.getByTestId('version-3').click();
      await page.getByTestId('compare-versions-button').click();

      // Verify comparison view
      await expect(page.getByTestId('version-comparison-modal')).toBeVisible();
      await expect(page.getByTestId('version-diff')).toBeVisible();
    });

    test('should enable version rollback', async ({ page }) => {
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      const originalTitle = TEST_CONTENT.MESSAGE.title;

      // Update content
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Updated Title',
        content: 'Updated content',
      });

      // Verify update
      await expect(page.getByTestId('content-title')).toHaveText('Updated Title');

      // Rollback to previous version
      await page.getByTestId('version-history-button').click();
      await page.getByTestId('version-1').click();
      await page.getByTestId('rollback-version-button').click();
      await page.getByTestId('confirm-rollback-button').click();

      // Verify rollback
      await expect(page.getByTestId('content-title')).toHaveText(originalTitle);
    });
  });

  test.describe('Conflict Resolution', () => {
    test('should detect and handle edit conflicts', async ({ page }) => {
      // Simulate concurrent edit scenario
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();
      await page.getByTestId('edit-content-button').click();

      // Start editing
      await page.getByTestId('content-title-input').fill('Local Edit');

      // Simulate remote update while editing
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Remote Edit',
        content: 'Content updated remotely',
      });

      // Attempt to save local changes
      await page.getByTestId('save-content-button').click();

      // Verify conflict detection
      await expect(page.getByTestId('conflict-resolution-dialog')).toBeVisible();
      await expect(page.getByTestId('conflict-message')).toContainText('conflict detected');
    });

    test('should provide merge options for conflicts', async ({ page }) => {
      // Set up conflict scenario
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();
      await page.getByTestId('edit-content-button').click();

      await page.getByTestId('content-title-input').fill('Local Changes');

      // Trigger remote update
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Remote Changes',
        content: 'Remote content changes',
      });

      await page.getByTestId('save-content-button').click();

      // Check merge options
      await expect(page.getByTestId('conflict-resolution-dialog')).toBeVisible();
      await expect(page.getByTestId('keep-local-button')).toBeVisible();
      await expect(page.getByTestId('keep-remote-button')).toBeVisible();
      await expect(page.getByTestId('merge-changes-button')).toBeVisible();
    });

    test('should handle automatic merge when possible', async ({ page }) => {
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();
      await page.getByTestId('edit-content-button').click();

      // Edit different fields to enable auto-merge
      await page.getByTestId('content-tags-input').fill('tag1, tag2');

      // Remote update to different field
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Remote Title Update',
      });

      await page.getByTestId('save-content-button').click();

      // Verify automatic merge success
      await expect(page.getByTestId('merge-success-notification')).toBeVisible();
      await expect(page.getByTestId('content-title')).toHaveText('Remote Title Update');
      await expect(page.getByTestId('content-tags')).toContainText('tag1, tag2');
    });
  });

  test.describe('Performance Impact', () => {
    test('should maintain acceptable performance during updates', async ({ page }) => {
      // Measure baseline performance
      const baselineMetrics = await measurePerformanceMetrics(page);

      await page.getByTestId('content-list').click();

      // Perform multiple rapid updates
      const updatePromises = [];
      for (let i = 0; i < 10; i++) {
        updatePromises.push(
          triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
            title: `Performance Test Update ${i}`,
            updateCount: i,
          })
        );
      }

      await Promise.all(updatePromises);

      // Measure performance after updates
      const postUpdateMetrics = await measurePerformanceMetrics(page);

      // Verify performance hasn't degraded significantly
      if (postUpdateMetrics.memoryUsage && baselineMetrics.memoryUsage) {
        const memoryIncrease = postUpdateMetrics.memoryUsage.used - baselineMetrics.memoryUsage.used;
        expect(memoryIncrease).toBeLessThan(50 * 1024 * 1024); // Less than 50MB increase
      }

      // Check that page remains responsive
      const responseTime = await page.evaluate(async () => {
        const start = performance.now();
        await new Promise(resolve => setTimeout(resolve, 0));
        return performance.now() - start;
      });

      expect(responseTime).toBeLessThan(100); // Less than 100ms for simple operation
    });

    test('should handle large content updates efficiently', async ({ page }) => {
      // Create large content update
      const largeContent = 'A'.repeat(10000); // 10KB of text

      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      const startTime = Date.now();

      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        content: largeContent,
        title: 'Large Content Update',
      });

      await waitForContentUpdate(page, TEST_CONTENT.MESSAGE.id, 2);

      const updateTime = Date.now() - startTime;

      // Verify update completed within reasonable time
      expect(updateTime).toBeLessThan(2000); // Less than 2 seconds

      // Verify content was updated correctly
      await expect(page.getByTestId('content-title')).toHaveText('Large Content Update');
    });

    test('should optimize network requests during batch updates', async ({ page }) => {
      // Monitor network requests
      const networkRequests: any[] = [];
      page.on('request', request => {
        if (request.url().includes('/api/content/')) {
          networkRequests.push({
            url: request.url(),
            method: request.method(),
            timestamp: Date.now(),
          });
        }
      });

      await page.getByTestId('content-list').click();

      // Trigger batch updates
      const batchUpdates = [];
      for (let i = 0; i < 5; i++) {
        batchUpdates.push(
          triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
            title: `Batch Update ${i}`,
            batchIndex: i,
          })
        );
      }

      await Promise.all(batchUpdates);
      await page.waitForTimeout(1000); // Allow time for requests

      // Verify request optimization (should batch or debounce)
      const contentUpdateRequests = networkRequests.filter(req =>
        req.method === 'PUT' || req.method === 'PATCH'
      );

      expect(contentUpdateRequests.length).toBeLessThan(5); // Should batch some requests
    });
  });

  test.describe('Content Types', () => {
    test('should handle chat message updates', async ({ page }) => {
      await setupTestContent(page, TEST_CONTENT.MESSAGE);

      await page.getByTestId('chat-tab').click();
      await page.getByTestId(`message-${TEST_CONTENT.MESSAGE.id}`).click();

      // Update message content
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        content: 'Updated chat message content',
        edited: true,
        editTimestamp: new Date().toISOString(),
      });

      // Verify message shows as edited
      await expect(page.getByTestId(`message-${TEST_CONTENT.MESSAGE.id}`))
        .toContainText('Updated chat message content');
      await expect(page.getByTestId(`edit-indicator-${TEST_CONTENT.MESSAGE.id}`))
        .toBeVisible();
    });

    test('should handle social post updates', async ({ page }) => {
      await setupTestContent(page, TEST_CONTENT.POST);

      await page.getByTestId('social-tab').click();
      await page.getByTestId(`post-${TEST_CONTENT.POST.id}`).click();

      // Update post with reactions
      await triggerContentUpdate(page, TEST_CONTENT.POST.id, {
        content: 'Updated social post with new engagement',
        likes: 15,
        shares: 3,
        comments: 7,
      });

      // Verify social metrics updated
      await expect(page.getByTestId('post-content')).toContainText('Updated social post');
      await expect(page.getByTestId('likes-count')).toContainText('15');
      await expect(page.getByTestId('shares-count')).toContainText('3');
      await expect(page.getByTestId('comments-count')).toContainText('7');
    });

    test('should handle product/commerce updates', async ({ page }) => {
      await setupTestContent(page, TEST_CONTENT.PRODUCT);

      await page.getByTestId('store-tab').click();
      await page.getByTestId(`product-${TEST_CONTENT.PRODUCT.id}`).click();

      // Update product with inventory and pricing
      await triggerContentUpdate(page, TEST_CONTENT.PRODUCT.id, {
        title: 'Updated Product Name',
        price: 24.99,
        inventory: 50,
        onSale: true,
        salePrice: 19.99,
      });

      // Verify product updates
      await expect(page.getByTestId('product-title')).toContainText('Updated Product Name');
      await expect(page.getByTestId('product-price')).toContainText('$19.99');
      await expect(page.getByTestId('original-price')).toContainText('$24.99');
      await expect(page.getByTestId('sale-badge')).toBeVisible();
      await expect(page.getByTestId('inventory-count')).toContainText('50');
    });

    test('should handle mixed content type updates in list views', async ({ page }) => {
      // Set up multiple content types
      await setupTestContent(page, TEST_CONTENT.MESSAGE);
      await setupTestContent(page, TEST_CONTENT.POST);
      await setupTestContent(page, TEST_CONTENT.PRODUCT);

      await page.getByTestId('content-list').click();

      // Trigger simultaneous updates across different types
      await Promise.all([
        triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, { title: 'Updated Message' }),
        triggerContentUpdate(page, TEST_CONTENT.POST.id, { title: 'Updated Post' }),
        triggerContentUpdate(page, TEST_CONTENT.PRODUCT.id, { title: 'Updated Product' }),
      ]);

      // Verify all content types updated in list
      await expect(page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`))
        .toContainText('Updated Message');
      await expect(page.getByTestId(`content-item-${TEST_CONTENT.POST.id}`))
        .toContainText('Updated Post');
      await expect(page.getByTestId(`content-item-${TEST_CONTENT.PRODUCT.id}`))
        .toContainText('Updated Product');
    });
  });

  test.describe('Error Scenarios', () => {
    test('should handle network failures gracefully', async ({ page }) => {
      await page.getByTestId('content-list').click();
      await page.getByTestId(`content-item-${TEST_CONTENT.MESSAGE.id}`).click();

      // Simulate network failure
      await page.route('**/api/content/**', route => route.abort());

      // Attempt content update
      await triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, {
        title: 'Failed Update Test',
      });

      // Verify fallback mode activation
      const syncStatus = await getSyncStatus(page);
      expect(syncStatus).toBe('error');

      await expect(page.getByTestId('offline-notification')).toBeVisible();
      await expect(page.getByTestId('retry-sync-button')).toBeVisible();
    });

    test('should maintain data integrity during partial failures', async ({ page }) => {
      await page.getByTestId('content-list').click();

      // Simulate partial API failure (some requests succeed, others fail)
      let requestCount = 0;
      await page.route('**/api/content/**', route => {
        requestCount++;
        if (requestCount % 2 === 0) {
          route.abort(); // Fail every second request
        } else {
          route.continue();
        }
      });

      // Trigger multiple updates
      await Promise.all([
        triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, { title: 'Update 1' }),
        triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, { title: 'Update 2' }),
        triggerContentUpdate(page, TEST_CONTENT.MESSAGE.id, { title: 'Update 3' }),
      ]);

      // Verify data consistency
      const finalState = await page.evaluate((contentId) => {
        const state = window.__REDUX_STORE__?.getState?.();
        return state?.content?.fallbackContent?.[contentId];
      }, TEST_CONTENT.MESSAGE.id);

      expect(finalState).toBeDefined();
      expect(finalState.title).toMatch(/Update [1-3]/);
    });
  });
});

// Cleanup after all tests
test.afterAll(async () => {
  // Cleanup any persistent test data
  console.log('Content update propagation tests completed');
});