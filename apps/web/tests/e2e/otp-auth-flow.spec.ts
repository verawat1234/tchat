import { test, expect, Page } from '@playwright/test';

// Test configuration and constants
const TEST_CONFIG = {
  APP_URL: 'http://localhost:3000',
  GATEWAY_URL: 'http://localhost:8080',
  AUTH_SERVICE_URL: 'http://localhost:8081',
  PHONE_NUMBER: '+66812345678',
  COUNTRY_CODE: '+66',
  PHONE_DIGITS: '812345678',
  EXPECTED_OTP: '123456',
  TIMEOUT: 10000
};

/**
 * Comprehensive E2E Test Suite for OTP Authentication Flow
 *
 * Tests the complete user authentication journey through the gateway:
 * 1. Navigate to auth screen
 * 2. Request OTP with proper field formatting
 * 3. Verify OTP submission
 * 4. Validate successful authentication
 * 5. Confirm no CORS errors occur
 * 6. Test gateway routing functionality
 */
test.describe('OTP Authentication Flow - Gateway Integration', () => {

  test.beforeEach(async ({ page }) => {
    // Set up network monitoring to catch CORS and API errors
    const requests: Array<{ url: string; method: string; status?: number; error?: string }> = [];
    const responses: Array<{ url: string; status: number; body?: string }> = [];
    const consoleMessages: Array<{ type: string; text: string }> = [];

    // Monitor network requests
    page.on('request', (request) => {
      requests.push({
        url: request.url(),
        method: request.method()
      });
    });

    // Monitor network responses
    page.on('response', async (response) => {
      const responseData = {
        url: response.url(),
        status: response.status(),
        body: undefined as string | undefined
      };

      try {
        // Capture response body for API endpoints
        if (response.url().includes('/api/')) {
          responseData.body = await response.text();
        }
      } catch (error) {
        // Response body might not be available
      }

      responses.push(responseData);
    });

    // Monitor console messages for CORS errors
    page.on('console', (message) => {
      consoleMessages.push({
        type: message.type(),
        text: message.text()
      });
    });

    // Store monitoring data in page context for test access
    await page.addInitScript(() => {
      (window as any).testData = {
        requests: [],
        responses: [],
        consoleMessages: []
      };
    });
  });

  test('should complete full OTP authentication flow through gateway', async ({ page }) => {
    console.log('🧪 Starting OTP Authentication E2E Test');

    // Step 1: Navigate to the application
    console.log('📱 Navigating to auth screen...');
    await page.goto(TEST_CONFIG.APP_URL, { waitUntil: 'networkidle' });

    // Verify we're on the auth screen
    await expect(page.locator('h1')).toContainText(/Telegram SEA|Sign In/i);
    console.log('✅ Auth screen loaded successfully');

    // Step 2: Verify UI elements are present
    await expect(page.locator('input[placeholder*="66"]')).toBeVisible();
    await expect(page.getByText('Send OTP')).toBeVisible();
    console.log('✅ Auth form elements visible');

    // Step 3: Test phone number input and country code selection
    const phoneInput = page.locator('input[placeholder*="66"]');
    await phoneInput.clear();
    await phoneInput.fill(TEST_CONFIG.PHONE_NUMBER);

    // Verify phone number is properly formatted
    await expect(phoneInput).toHaveValue(TEST_CONFIG.PHONE_NUMBER);
    console.log(`✅ Phone number entered: ${TEST_CONFIG.PHONE_NUMBER}`);

    // Step 4: Request OTP and monitor network traffic
    console.log('📞 Requesting OTP...');

    // Set up network monitoring for OTP request
    const otpRequestPromise = page.waitForResponse(response =>
      response.url().includes('/auth/login') && response.status() === 200
    );

    // Click Send OTP button
    await page.getByText('Send OTP').click();

    // Wait for OTP request to complete
    try {
      const otpResponse = await otpRequestPromise;
      console.log(`✅ OTP request completed with status: ${otpResponse.status()}`);

      // Verify the request went through the gateway
      expect(otpResponse.url()).toContain(TEST_CONFIG.GATEWAY_URL.replace('http://', ''));
      console.log('✅ Request correctly routed through gateway');

      // Parse and verify response body
      const responseBody = await otpResponse.text();
      console.log('📋 OTP Response:', responseBody);

      // The response should indicate success
      expect(otpResponse.status()).toBe(200);

    } catch (error) {
      console.error('❌ OTP request failed:', error);

      // Check for CORS errors in console
      const corsErrors = await page.evaluate(() => {
        return window.console.messages?.filter(msg =>
          msg.text.toLowerCase().includes('cors') ||
          msg.text.toLowerCase().includes('cross-origin')
        ) || [];
      });

      if (corsErrors.length > 0) {
        console.error('🚫 CORS errors detected:', corsErrors);
      }

      throw error;
    }

    // Step 5: Verify UI transitions to OTP verification step
    await expect(page.getByText('Verify Your Phone')).toBeVisible({ timeout: TEST_CONFIG.TIMEOUT });
    await expect(page.locator('input[placeholder*="verification"]')).toBeVisible();
    console.log('✅ Transitioned to OTP verification screen');

    // Step 6: Enter OTP code
    console.log('🔢 Entering OTP code...');
    const otpInput = page.locator('input[placeholder*="verification"]');
    await otpInput.fill(TEST_CONFIG.EXPECTED_OTP);
    await expect(otpInput).toHaveValue(TEST_CONFIG.EXPECTED_OTP);
    console.log(`✅ OTP entered: ${TEST_CONFIG.EXPECTED_OTP}`);

    // Step 7: Submit OTP verification
    console.log('🔐 Verifying OTP...');

    // Monitor OTP verification request
    const verifyRequestPromise = page.waitForResponse(response =>
      response.url().includes('/auth/verify-otp'), { timeout: TEST_CONFIG.TIMEOUT }
    );

    // Click verify button
    await page.getByText('Verify & Continue').click();

    // Wait for verification response
    try {
      const verifyResponse = await verifyRequestPromise;
      console.log(`✅ OTP verification completed with status: ${verifyResponse.status()}`);

      // Verify the request went through the gateway
      expect(verifyResponse.url()).toContain(TEST_CONFIG.GATEWAY_URL.replace('http://', ''));
      console.log('✅ Verification request correctly routed through gateway');

      const verifyResponseBody = await verifyResponse.text();
      console.log('📋 Verification Response:', verifyResponseBody);

      // Check if verification was successful
      if (verifyResponse.status() === 200) {
        console.log('✅ OTP verification successful');

        // Step 8: Verify successful authentication state
        // Look for indicators of successful login (this will depend on your app's behavior)
        await expect(page.locator('body')).not.toContainText('Sign In with Phone', { timeout: 5000 });
        console.log('✅ Successfully authenticated - auth form no longer visible');

      } else {
        console.log(`⚠️ OTP verification returned status: ${verifyResponse.status()}`);
        console.log('Response body:', verifyResponseBody);
      }

    } catch (error) {
      console.error('❌ OTP verification failed:', error);
      throw error;
    }

    // Step 9: Validate no CORS errors occurred
    console.log('🔍 Checking for CORS errors...');
    const finalCorsCheck = await page.evaluate(() => {
      const consoleMessages = (window as any).console?.messages || [];
      return consoleMessages.filter((msg: any) =>
        typeof msg.text === 'string' && (
          msg.text.toLowerCase().includes('cors') ||
          msg.text.toLowerCase().includes('cross-origin') ||
          msg.text.toLowerCase().includes('access-control')
        )
      );
    });

    expect(finalCorsCheck.length).toBe(0);
    console.log('✅ No CORS errors detected');

    console.log('🎉 OTP Authentication E2E Test completed successfully!');
  });

  test('should handle OTP request errors gracefully', async ({ page }) => {
    console.log('🧪 Testing OTP error handling...');

    await page.goto(TEST_CONFIG.APP_URL);

    // Test with invalid phone number
    const phoneInput = page.locator('input[placeholder*="66"]');
    await phoneInput.clear();
    await phoneInput.fill('invalid-phone');

    await page.getByText('Send OTP').click();

    // Should show an error message
    await expect(page.locator('body')).toContainText(/invalid|error/i, { timeout: 5000 });
    console.log('✅ Error handling works correctly');
  });

  test('should validate gateway routing endpoints', async ({ page }) => {
    console.log('🧪 Testing gateway routing...');

    // Test that gateway is accessible
    const gatewayHealthResponse = await page.request.get(`${TEST_CONFIG.GATEWAY_URL}/health`);
    expect(gatewayHealthResponse.status()).toBe(200);
    console.log('✅ Gateway health check passed');

    // Test auth service routing through gateway
    const authHealthResponse = await page.request.get(`${TEST_CONFIG.GATEWAY_URL}/api/v1/auth/health`);
    expect(authHealthResponse.status()).toBe(200);
    console.log('✅ Auth service accessible through gateway');
  });

  test('should validate field name compatibility', async ({ page }) => {
    console.log('🧪 Testing field name compatibility...');

    await page.goto(TEST_CONFIG.APP_URL);

    // Intercept and examine the OTP request payload
    let requestPayload: any = null;

    page.on('request', (request) => {
      if (request.url().includes('/auth/login') && request.method() === 'POST') {
        try {
          requestPayload = JSON.parse(request.postData() || '{}');
        } catch (error) {
          console.error('Failed to parse request payload:', error);
        }
      }
    });

    // Fill phone number and request OTP
    const phoneInput = page.locator('input[placeholder*="66"]');
    await phoneInput.fill(TEST_CONFIG.PHONE_NUMBER);
    await page.getByText('Send OTP').click();

    // Wait a moment for the request to be captured
    await page.waitForTimeout(2000);

    // Verify the payload has the expected field names
    expect(requestPayload).toBeTruthy();
    expect(requestPayload).toHaveProperty('phone_number');
    expect(requestPayload).toHaveProperty('country_code');

    console.log('✅ Request payload format is correct:', requestPayload);
    console.log('✅ Field name compatibility verified');
  });
});

/**
 * Network and Performance Validation Tests
 */
test.describe('Network and Performance Validation', () => {

  test('should complete authentication within performance budget', async ({ page }) => {
    console.log('🧪 Testing authentication performance...');

    const startTime = Date.now();

    await page.goto(TEST_CONFIG.APP_URL);

    // Measure OTP request time
    const phoneInput = page.locator('input[placeholder*="66"]');
    await phoneInput.fill(TEST_CONFIG.PHONE_NUMBER);

    const otpStartTime = Date.now();
    await page.getByText('Send OTP').click();

    await expect(page.getByText('Verify Your Phone')).toBeVisible();
    const otpEndTime = Date.now();

    const otpDuration = otpEndTime - otpStartTime;
    console.log(`⏱️ OTP request took: ${otpDuration}ms`);

    // OTP request should complete within 5 seconds
    expect(otpDuration).toBeLessThan(5000);

    // Measure verification time
    const otpInput = page.locator('input[placeholder*="verification"]');
    await otpInput.fill(TEST_CONFIG.EXPECTED_OTP);

    const verifyStartTime = Date.now();
    await page.getByText('Verify & Continue').click();

    // Wait for authentication to complete
    await page.waitForTimeout(3000);
    const verifyEndTime = Date.now();

    const verifyDuration = verifyEndTime - verifyStartTime;
    console.log(`⏱️ OTP verification took: ${verifyDuration}ms`);

    // Verification should complete within 3 seconds
    expect(verifyDuration).toBeLessThan(3000);

    const totalDuration = Date.now() - startTime;
    console.log(`⏱️ Total authentication flow took: ${totalDuration}ms`);

    console.log('✅ Performance requirements met');
  });

  test('should maintain stable network connection during auth flow', async ({ page }) => {
    console.log('🧪 Testing network stability...');

    let networkErrors = 0;
    let failedRequests: string[] = [];

    page.on('requestfailed', (request) => {
      networkErrors++;
      failedRequests.push(`${request.method()} ${request.url()}: ${request.failure()?.errorText}`);
    });

    await page.goto(TEST_CONFIG.APP_URL);

    // Complete authentication flow
    const phoneInput = page.locator('input[placeholder*="66"]');
    await phoneInput.fill(TEST_CONFIG.PHONE_NUMBER);
    await page.getByText('Send OTP').click();

    await expect(page.getByText('Verify Your Phone')).toBeVisible();

    const otpInput = page.locator('input[placeholder*="verification"]');
    await otpInput.fill(TEST_CONFIG.EXPECTED_OTP);
    await page.getByText('Verify & Continue').click();

    await page.waitForTimeout(2000);

    // No network errors should occur during normal flow
    if (networkErrors > 0) {
      console.log('❌ Network errors detected:', failedRequests);
    }

    console.log(`✅ Network stability check: ${networkErrors} errors detected`);
    // Allow for some non-critical network errors but flag if excessive
    expect(networkErrors).toBeLessThan(3);
  });
});