import { test, expect } from '@playwright/test';

/**
 * Focused CORS and Gateway Validation Test
 *
 * This test specifically validates that:
 * 1. CORS is properly configured
 * 2. Gateway routing works correctly
 * 3. Field names are compatible between frontend and backend
 * 4. No network errors occur during authentication flow
 */
test.describe('CORS and Gateway Validation', () => {
  test('should validate CORS fix and gateway routing for OTP flow', async ({ page }) => {
    console.log('🧪 Testing CORS fix and gateway routing');

    // Monitor for CORS errors
    const corsErrors: string[] = [];
    const networkErrors: string[] = [];
    let otpRequestSucceeded = false;
    let requestPayload: any = null;

    // Monitor console for CORS errors
    page.on('console', (message) => {
      const text = message.text();
      if (text.toLowerCase().includes('cors') ||
          text.toLowerCase().includes('cross-origin') ||
          text.toLowerCase().includes('access-control')) {
        corsErrors.push(text);
      }
    });

    // Monitor network requests
    page.on('request', (request) => {
      // Capture OTP request payload
      if (request.url().includes('/auth/login') && request.method() === 'POST') {
        try {
          requestPayload = JSON.parse(request.postData() || '{}');
        } catch (error) {
          console.log('Could not parse request payload');
        }
      }
    });

    // Monitor responses for successful OTP request
    page.on('response', (response) => {
      if (response.url().includes('/auth/login') && response.status() === 200) {
        otpRequestSucceeded = true;
        console.log('✅ OTP request succeeded through gateway');
      }
    });

    page.on('requestfailed', (request) => {
      // Only track API-related failures, ignore health checks to other services
      if (request.url().includes('/api/')) {
        networkErrors.push(`${request.method()} ${request.url()}: ${request.failure()?.errorText}`);
      }
    });

    // Navigate to auth screen
    await page.goto('http://localhost:3000');
    console.log('📱 Navigated to auth screen');

    // Verify auth screen loads correctly
    await expect(page.locator('h1')).toContainText(/Telegram SEA|Sign In/i);
    await expect(page.locator('input[placeholder*="66"]')).toBeVisible();
    console.log('✅ Auth screen elements visible');

    // Fill phone number
    const phoneInput = page.locator('input[placeholder*="66"]');
    await phoneInput.clear();
    await phoneInput.fill('+66812345678');
    console.log('✅ Phone number entered');

    // Request OTP
    await page.getByText('Send OTP').click();
    console.log('📞 OTP request initiated');

    // Wait for request to complete
    await page.waitForTimeout(3000);

    // Verify gateway routing worked
    expect(otpRequestSucceeded).toBe(true);
    console.log('✅ OTP request successfully routed through gateway');

    // Verify field names are correct
    expect(requestPayload).toBeTruthy();
    expect(requestPayload).toHaveProperty('phone_number');
    expect(requestPayload).toHaveProperty('country_code');
    expect(requestPayload.phone_number).toBe('+66812345678');
    expect(requestPayload.country_code).toBe('TH');
    console.log('✅ Field names are correctly formatted:', requestPayload);

    // Verify no CORS errors occurred
    expect(corsErrors).toHaveLength(0);
    console.log('✅ No CORS errors detected');

    // Verify no critical API network errors
    const criticalApiErrors = networkErrors.filter(error =>
      error.includes('/auth/') || error.includes('cors')
    );
    expect(criticalApiErrors).toHaveLength(0);
    console.log('✅ No critical API network errors');

    console.log('🎉 CORS fix and gateway routing validation completed successfully!');
  });

  test('should validate gateway health and auth service routing', async ({ page }) => {
    console.log('🧪 Testing gateway health and service routing');

    // Test gateway health
    const gatewayHealthResponse = await page.request.get('http://localhost:8080/health');
    expect(gatewayHealthResponse.status()).toBe(200);

    const gatewayHealthData = await gatewayHealthResponse.json();
    expect(gatewayHealthData).toHaveProperty('status', 'healthy');
    console.log('✅ Gateway health check passed');

    // Test auth service routing through gateway
    const authHealthResponse = await page.request.get('http://localhost:8080/api/v1/auth/health');
    expect(authHealthResponse.status()).toBe(200);

    const authHealthData = await authHealthResponse.json();
    expect(authHealthData.success).toBe(true);
    console.log('✅ Auth service accessible through gateway');

    // Test direct OTP request to verify field compatibility
    const otpResponse = await page.request.post('http://localhost:8080/api/v1/auth/login', {
      headers: {
        'Content-Type': 'application/json',
        'Origin': 'http://localhost:3000'
      },
      data: {
        phone_number: '+66812345678',
        country_code: 'TH'
      }
    });

    expect(otpResponse.status()).toBe(200);
    const otpData = await otpResponse.json();
    expect(otpData.success).toBe(true);
    expect(otpData.data.success).toBe(true);
    expect(otpData.data.message).toBe('OTP sent successfully');
    console.log('✅ OTP request with correct field names successful');

    console.log('🎉 Gateway health and service routing validation completed!');
  });
});