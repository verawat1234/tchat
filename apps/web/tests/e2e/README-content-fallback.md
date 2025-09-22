# Content Fallback E2E Tests

This document provides comprehensive information about the content fallback E2E tests implemented for T058.

## Overview

The content fallback E2E tests validate the complete offline/online content management system, ensuring users can continue using the application effectively even when the content management system is unavailable.

## Test Coverage

### 1. Offline Scenarios 🌐
- **Network Simulation**: Tests app behavior when network is completely unavailable
- **Cache Loading**: Validates localStorage fallback content loading and storage
- **Navigation**: Ensures navigation works with cached content while offline
- **State Persistence**: Verifies fallback mode persists across page reloads

### 2. API Failures 🚫
- **Server Errors**: Tests 500 Internal Server Error responses
- **Network Timeouts**: Simulates slow/unresponsive APIs
- **Connection Issues**: Tests connection refused and fetch errors
- **Graceful Degradation**: Validates appropriate error messages when no cache available

### 3. Cache Behavior 💾
- **Response Caching**: Validates successful API responses are cached in localStorage
- **Recovery Workflows**: Tests automatic recovery when connectivity is restored
- **Expiration Handling**: Ensures expired cache content is handled appropriately
- **Corruption Recovery**: Tests graceful handling of corrupted cache data

### 4. User Feedback 📱
- **Fallback Indicators**: Validates user feedback about fallback mode status
- **Retry Functionality**: Tests retry mechanisms after failures
- **Sync Status**: Verifies sync status indicators work correctly
- **Accessibility**: Ensures screen readers announce fallback states

### 5. Performance Under Stress ⚡
- **Concurrent Failures**: Tests multiple simultaneous API failures
- **Storage Limits**: Handles localStorage quota exceeded scenarios
- **UI Responsiveness**: Maintains responsive UI during fallback operations
- **Memory Management**: Efficient handling of large cache datasets

### 6. Content Consistency 🔍
- **Structure Preservation**: Ensures fallback content matches expected structure
- **Formatting**: Validates content formatting is preserved in fallback mode
- **Mixed Content Types**: Handles various content types (text, objects, arrays)
- **Accessibility**: Maintains accessibility attributes in fallback mode

### 7. Cross-Browser Compatibility 🌏
- **Browser Support**: Tests consistency across Chrome, Firefox, Safari
- **localStorage Behavior**: Validates localStorage differences across browsers
- **Feature Detection**: Graceful fallback when features unavailable

## Test Architecture

### Helper Functions

#### `setNetworkCondition(page, condition)`
Simulates network conditions:
- `'online'`: Normal network connectivity
- `'offline'`: Complete network disconnection
- `'slow'`: 5-second delay simulation
- `'error'`: Forces 500 server errors

#### `setupCacheContent(page)`
Pre-populates localStorage with test content:
```typescript
const cachedItem = {
  id: 'test-content-123',
  content: 'Test content value for fallback testing',
  cachedAt: Date.now(),
  expiresAt: Date.now() + (24 * 60 * 60 * 1000),
  size: 256,
  accessCount: 1,
  lastAccessed: Date.now(),
  type: 'text',
  version: 1
};
```

#### `checkFallbackIndicators(page, expectedMode)`
Validates fallback mode UI indicators:
- Searches for offline/cached/fallback text
- Checks `data-testid="fallback-indicator"`
- Validates ARIA labels for accessibility

#### `waitForContentOrFallback(page)`
Waits for content loading or fallback activation:
- Watches for `data-testid="content-loaded"`
- Monitors `data-testid="fallback-content"`
- Detects `data-testid="content-error"`

### Mock API Responses

The tests use comprehensive API mocking:

```typescript
const MOCK_CONTENT_RESPONSE = {
  id: 'test-content-123',
  value: 'Test content value for fallback testing',
  type: 'text',
  category: 'navigation',
  version: 1,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString()
};
```

## Running the Tests

### Prerequisites
- Node.js 18+
- Playwright installed (`npm install @playwright/test`)
- Application server running on `http://localhost:3000`

### Basic Test Execution

```bash
# Run all content fallback tests
npx playwright test content-fallback.spec.ts

# Run with UI mode for debugging
npx playwright test content-fallback.spec.ts --ui

# Run specific test suite
npx playwright test content-fallback.spec.ts --grep "Offline Scenarios"

# Run with headed browser for visual debugging
npx playwright test content-fallback.spec.ts --headed

# Generate test report
npx playwright test content-fallback.spec.ts --reporter=html
```

### Test Configuration Options

```bash
# Test different browsers
npx playwright test content-fallback.spec.ts --project=chromium
npx playwright test content-fallback.spec.ts --project=firefox
npx playwright test content-fallback.spec.ts --project=webkit

# Parallel execution
npx playwright test content-fallback.spec.ts --workers=3

# Retry failed tests
npx playwright test content-fallback.spec.ts --retries=2

# Debug mode
npx playwright test content-fallback.spec.ts --debug
```

## Test Data Selectors

The tests use these data-testid selectors for reliable element identification:

### Required Test IDs
Your components should implement these test IDs:

```html
<!-- Content states -->
<div data-testid="content-loaded">Content successfully loaded</div>
<div data-testid="fallback-content">Cached content being displayed</div>
<div data-testid="content-error">Error loading content</div>

<!-- Fallback indicators -->
<div data-testid="fallback-indicator">🔌 Using cached content</div>
<div data-testid="offline-indicator">📶 You're offline</div>

<!-- Interactive elements -->
<button data-testid="retry-button">Try Again</button>
<div data-testid="sync-status">Sync status: online</div>
```

### Accessibility Labels
The tests also check for accessibility:

```html
<!-- ARIA labels for screen readers -->
<div aria-label="Content loaded successfully">...</div>
<div aria-label="Using cached content while offline">...</div>
<div aria-label="Failed to load content">...</div>

<!-- Live regions for announcements -->
<div aria-live="polite" aria-atomic="true">Content loading...</div>
```

## Expected Behaviors

### Online Mode
1. ✅ API requests succeed
2. ✅ Content loads from server
3. ✅ Responses cached in localStorage
4. ✅ No fallback indicators visible
5. ✅ Sync status shows "online"

### Offline Mode
1. ✅ API requests fail
2. ✅ Content loads from localStorage cache
3. ✅ Fallback indicators visible
4. ✅ Navigation still functional
5. ✅ Sync status shows "offline"

### Recovery Mode
1. ✅ Network restored
2. ✅ API requests succeed again
3. ✅ Fresh content loaded and cached
4. ✅ Fallback indicators disappear
5. ✅ Sync status returns to "online"

## Debugging Failed Tests

### Common Issues and Solutions

#### 1. Test Timeout
```bash
Error: Test timeout of 30000ms exceeded
```
**Solution**: Increase timeout or check for slow API responses
```typescript
test('should handle slow APIs', async ({ page }) => {
  test.setTimeout(60000); // Increase timeout
  // ... test code
});
```

#### 2. Element Not Found
```bash
Error: locator.first() - element not found
```
**Solution**: Ensure test IDs are implemented in components
```typescript
// Check if element exists before interacting
if (await page.locator('[data-testid="retry-button"]').count() > 0) {
  await page.locator('[data-testid="retry-button"]').click();
}
```

#### 3. localStorage Issues
```bash
Error: localStorage is not available
```
**Solution**: Check browser context setup
```typescript
// Ensure localStorage is available
const isAvailable = await page.evaluate(() => {
  try {
    localStorage.setItem('test', 'test');
    localStorage.removeItem('test');
    return true;
  } catch {
    return false;
  }
});
```

#### 4. Network Simulation Not Working
```bash
Error: API calls still succeeding when offline
```
**Solution**: Verify route mocking setup
```typescript
// Ensure routes are properly mocked before navigation
await page.route('**/api/**', route => route.abort());
await page.goto('/'); // Navigate after setting up routes
```

## Performance Expectations

### Response Time Targets
- **Cache Loading**: < 500ms
- **Fallback Activation**: < 1000ms
- **Recovery Time**: < 2000ms
- **UI Interactions**: < 200ms

### Storage Limits
- **Max Cache Size**: 5MB
- **Item Count**: < 1000 items
- **Cleanup Frequency**: Every hour

### Browser Support
- **Chrome**: 90+ ✅
- **Firefox**: 88+ ✅
- **Safari**: 14+ ✅
- **Edge**: 90+ ✅

## Continuous Integration

### GitHub Actions Example

```yaml
name: E2E Content Fallback Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    timeout-minutes: 60
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-node@v3
      with:
        node-version: 18
    - name: Install dependencies
      run: npm ci
    - name: Install Playwright
      run: npx playwright install --with-deps
    - name: Start application
      run: npm run dev &
    - name: Wait for app
      run: npx wait-on http://localhost:3000
    - name: Run content fallback tests
      run: npx playwright test content-fallback.spec.ts
    - uses: actions/upload-artifact@v3
      if: always()
      with:
        name: playwright-report
        path: playwright-report/
        retention-days: 30
```

## Troubleshooting Guide

### Test Environment Setup

1. **Verify Application State**
   ```bash
   curl http://localhost:3000/health
   ```

2. **Check Content API Endpoints**
   ```bash
   curl http://localhost:3000/api/content/items/test
   ```

3. **Validate localStorage Support**
   ```javascript
   // In browser console
   localStorage.setItem('test', 'value');
   console.log(localStorage.getItem('test'));
   ```

### Debug Mode Usage

```bash
# Run single test with debug
npx playwright test content-fallback.spec.ts --grep "should show cached content" --debug

# Enable verbose logging
DEBUG=pw:api npx playwright test content-fallback.spec.ts

# Generate trace files
npx playwright test content-fallback.spec.ts --trace on
```

### Visual Testing

```bash
# Update visual baselines
npx playwright test content-fallback.spec.ts --update-snapshots

# Compare visual differences
npx playwright show-report
```

## Contributing

When adding new fallback scenarios:

1. **Follow Test Structure**: Use the established helper functions
2. **Add Test IDs**: Ensure components have proper data-testid attributes
3. **Document Behavior**: Add comments explaining expected outcomes
4. **Include Edge Cases**: Test error conditions and boundary cases
5. **Verify Accessibility**: Include screen reader and keyboard navigation tests

### Example New Test

```typescript
test('should handle [new scenario]', async ({ page }) => {
  // Setup
  await setupCacheContent(page);

  // Action
  await [perform specific action];

  // Verification
  await waitForContentOrFallback(page);
  await checkFallbackIndicators(page, [expected mode]);

  // Assertions
  await expect([specific element]).toBeVisible();
});
```

## Related Documentation

- [Content Fallback Service](../../src/services/contentFallback.ts)
- [Content Middleware](../../src/store/middleware/contentFallbackMiddleware.ts)
- [Content Loader Component](../../src/components/ContentLoader.tsx)
- [Playwright Configuration](../../playwright.config.ts)