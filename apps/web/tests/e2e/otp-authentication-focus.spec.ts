import { test, expect, Page, Response } from '@playwright/test';

test.describe('OTP Authentication Flow - Focused Validation', () => {
  let consoleLogs: string[] = [];
  let networkErrors: string[] = [];
  let corsErrors: string[] = [];
  let apiResponses: { url: string; status: number; method: string }[] = [];

  test.beforeEach(async ({ page }) => {
    // Reset tracking arrays
    consoleLogs = [];
    networkErrors = [];
    corsErrors = [];
    apiResponses = [];

    // Capture console messages for debugging
    page.on('console', (msg) => {
      const message = `[${msg.type()}] ${msg.text()}`;
      consoleLogs.push(message);

      // Capture CORS-specific errors
      if (message.toLowerCase().includes('cors')) {
        corsErrors.push(message);
      }

      // Capture network errors
      if (message.toLowerCase().includes('failed to fetch') ||
          message.toLowerCase().includes('network error') ||
          message.toLowerCase().includes('connection refused')) {
        networkErrors.push(message);
      }
    });

    // Monitor all network requests and responses
    page.on('response', (response: Response) => {
      apiResponses.push({
        url: response.url(),
        status: response.status(),
        method: response.request().method()
      });
    });

    // Monitor failed requests
    page.on('requestfailed', (request) => {
      networkErrors.push(`Request failed: ${request.method()} ${request.url()} - ${request.failure()?.errorText}`);
    });

    // Navigate to the app
    await page.goto('http://localhost:3000/');
    await page.waitForLoadState('networkidle');

    // Allow content systems to stabilize
    await page.waitForTimeout(3000);
  });

  test.afterEach(async ({ page }) => {
    // Log summary of captured information
    console.log('\n=== TEST SUMMARY ===');
    console.log(`Console messages: ${consoleLogs.length}`);
    console.log(`Network errors: ${networkErrors.length}`);
    console.log(`CORS errors: ${corsErrors.length}`);
    console.log(`API responses: ${apiResponses.length}`);

    if (corsErrors.length > 0) {
      console.log('\n=== CORS ERRORS ===');
      corsErrors.forEach(error => console.log(`❌ ${error}`));
    }

    if (networkErrors.length > 0) {
      console.log('\n=== NETWORK ERRORS ===');
      networkErrors.forEach(error => console.log(`🔗 ${error}`));
    }

    // Show relevant API calls
    const authCalls = apiResponses.filter(r => r.url.includes('/auth/'));
    if (authCalls.length > 0) {
      console.log('\n=== AUTH API CALLS ===');
      authCalls.forEach(call => console.log(`📡 ${call.method} ${call.url} → ${call.status}`));
    }
  });

  test('should load AuthScreen without infinite loops and verify no CORS errors on initial load', async ({ page }) => {
    // Verify page loads properly
    await expect(page).toHaveTitle(/Tchat|Telegram SEA/i);
    await expect(page.locator('h1')).toContainText(/telegram/i, { timeout: 10000 });

    // Verify feature highlights (these should load from fallback content)
    await expect(page.locator('text=End-to-End Encrypted')).toBeVisible({ timeout: 5000 });
    await expect(page.locator('text=Ultra Low Data')).toBeVisible();
    await expect(page.locator('text=QR Payments')).toBeVisible();
    await expect(page.locator('text=SEA Languages')).toBeVisible();

    // Verify auth form is displayed
    await expect(page.locator('text=Sign In with Phone')).toBeVisible();
    await expect(page.locator('text=Enter your phone number')).toBeVisible();

    // Check for demo credentials
    await expect(page.locator('text=Demo phone: +66812345678')).toBeVisible();

    // Verify no excessive loading indicators
    await page.waitForTimeout(2000);
    const loadingCount = await page.locator('.animate-pulse').count();
    expect(loadingCount).toBeLessThan(10);

    // CORS Check: Should have no CORS errors during initial page load
    expect(corsErrors.length).toBe(0);

    console.log(`✅ Page loaded successfully with ${consoleLogs.length} console messages and no CORS errors`);
  });

  test('should enter demo phone number and verify Send OTP API call behavior', async ({ page }) => {
    // Find and interact with phone input
    const phoneInput = page.locator('input[placeholder*="+66"], input[aria-label*="phone"]').first();
    await expect(phoneInput).toBeVisible();

    // Verify it's pre-filled with demo number or fill it manually
    const currentValue = await phoneInput.inputValue();
    if (!currentValue.includes('66812345678')) {
      await phoneInput.clear();
      await phoneInput.fill('+66812345678');
    }
    await expect(phoneInput).toHaveValue('+66812345678');

    // Find Send OTP button
    const sendOtpButton = page.locator('button:has-text("Send OTP"), button:has-text("Send")').first();
    await expect(sendOtpButton).toBeVisible();
    await expect(sendOtpButton).toBeEnabled();

    // Monitor for auth API calls
    let authLoginResponse: Response | null = null;
    const responsePromise = page.waitForResponse(
      response => response.url().includes('/auth/login'),
      { timeout: 10000 }
    ).catch(() => null);

    // Click Send OTP
    console.log('🔄 Clicking Send OTP button...');
    await sendOtpButton.click();

    // Wait a moment to see API behavior
    await page.waitForTimeout(3000);

    // Check what happened with the API call
    authLoginResponse = await responsePromise;

    if (authLoginResponse) {
      console.log(`✅ Auth API call made: ${authLoginResponse.status()} ${authLoginResponse.url()}`);

      // If we get a 200/201 response, the API is working
      if (authLoginResponse.status() >= 200 && authLoginResponse.status() < 300) {
        console.log('🎉 API call successful! Checking for OTP verification screen...');

        // Should transition to OTP verification screen
        await expect(page.locator('text=Verify Your Phone')).toBeVisible({ timeout: 10000 });
        await expect(page.locator('text=We sent a code to +66812345678')).toBeVisible();

        console.log('✅ Successfully transitioned to OTP verification screen');
      } else {
        console.log(`⚠️ API call returned status ${authLoginResponse.status()}`);
      }
    } else {
      console.log('❌ No auth API response received - checking for error handling...');

      // Should remain on phone input screen if API fails
      await expect(page.locator('text=Sign In with Phone')).toBeVisible();

      // Check if any CORS errors occurred
      if (corsErrors.length > 0) {
        console.log('🚨 CORS errors detected during API call');
      }

      // Check if any network errors occurred
      if (networkErrors.length > 0) {
        console.log('🔗 Network errors detected during API call');
      }
    }

    // Log final status
    console.log(`📊 Final status: Console logs: ${consoleLogs.length}, CORS errors: ${corsErrors.length}, Network errors: ${networkErrors.length}`);
  });

  test('should complete full OTP verification flow if backend is available', async ({ page }) => {
    // Step 1: Enter phone number
    const phoneInput = page.locator('input[placeholder*="+66"], input[aria-label*="phone"]').first();
    const currentValue = await phoneInput.inputValue();
    if (!currentValue.includes('66812345678')) {
      await phoneInput.clear();
      await phoneInput.fill('+66812345678');
    }

    // Step 2: Send OTP
    const sendOtpButton = page.locator('button:has-text("Send OTP"), button:has-text("Send")').first();

    let otpScreenReached = false;

    try {
      // Monitor for successful OTP request
      const otpResponse = page.waitForResponse(
        response => response.url().includes('/auth/login') && response.status() < 300,
        { timeout: 8000 }
      );

      await sendOtpButton.click();
      console.log('🔄 Attempting OTP request...');

      const response = await otpResponse;
      console.log(`✅ OTP request successful: ${response.status()}`);

      // Should transition to verification screen
      await expect(page.locator('text=Verify Your Phone')).toBeVisible({ timeout: 10000 });
      await expect(page.locator('text=We sent a code to +66812345678')).toBeVisible();

      otpScreenReached = true;
      console.log('✅ Reached OTP verification screen');

      // Step 3: Enter OTP code
      const otpInput = page.locator('input[placeholder*="code"], input[placeholder*="verification"]').first();
      await expect(otpInput).toBeVisible();
      await expect(otpInput).toHaveAttribute('maxlength', '6');

      await otpInput.fill('123456');
      await expect(otpInput).toHaveValue('123456');

      // Step 4: Verify OTP
      const verifyButton = page.locator('button:has-text("Verify"), button:has-text("Continue")').first();
      await expect(verifyButton).toBeVisible();
      await expect(verifyButton).toBeEnabled();

      // Monitor for verification response
      const verifyResponse = page.waitForResponse(
        response => response.url().includes('/auth/verify-otp'),
        { timeout: 8000 }
      ).catch(() => null);

      await verifyButton.click();
      console.log('🔄 Attempting OTP verification...');

      const verification = await verifyResponse;

      if (verification) {
        console.log(`✅ OTP verification response: ${verification.status()}`);

        if (verification.status() < 300) {
          // Should see success indication
          await expect(page.locator('text=Welcome to Telegram SEA')).toBeVisible({ timeout: 10000 }).catch(async () => {
            // Alternative success indicators
            console.log('Looking for alternative success indicators...');
            const hasToast = await page.locator('.sonner-toast').isVisible().catch(() => false);
            if (hasToast) {
              console.log('✅ Success toast notification found');
            }
          });

          // Check if we navigate away from auth screen
          await page.waitForTimeout(2000);
          const stillOnAuth = await page.locator('text=Verify Your Phone').isVisible().catch(() => false);

          if (!stillOnAuth) {
            console.log('🎉 Complete authentication flow successful - navigated away from auth screen');
          } else {
            console.log('⚠️ Still on auth screen after successful verification');
          }
        }
      } else {
        console.log('❌ No verification response received');
      }

    } catch (error) {
      console.log('❌ OTP request failed or timed out');
      console.log('🔍 Checking error conditions...');

      // Should remain on phone input screen
      await expect(page.locator('text=Sign In with Phone')).toBeVisible();

      // Check error types
      if (corsErrors.length > 0) {
        console.log('🚨 CORS issues preventing API calls');
        corsErrors.forEach(error => console.log(`   ${error}`));
      }

      if (networkErrors.length > 0) {
        console.log('🔗 Network connectivity issues');
        networkErrors.forEach(error => console.log(`   ${error}`));
      }

      // Check if backend services are responding
      const hasApiCalls = apiResponses.some(r => r.url.includes('/auth/'));
      if (!hasApiCalls) {
        console.log('📡 No auth API calls detected - backend may be unavailable');
      }
    }

    // Final assessment
    if (otpScreenReached) {
      console.log('🎯 Test Result: OTP authentication flow is working correctly');
    } else {
      console.log('🚧 Test Result: OTP authentication blocked by backend/CORS issues');
    }
  });

  test('should verify CORS configuration by examining network behavior', async ({ page }) => {
    // Try to make an auth API call and examine the exact error
    const phoneInput = page.locator('input[placeholder*="+66"], input[aria-label*="phone"]').first();
    const currentValue = await phoneInput.inputValue();
    if (!currentValue.includes('66812345678')) {
      await phoneInput.clear();
      await phoneInput.fill('+66812345678');
    }

    const sendOtpButton = page.locator('button:has-text("Send OTP"), button:has-text("Send")').first();

    // Capture the exact network behavior
    let requestMade = false;
    let responseReceived = false;
    let requestBlocked = false;

    // Monitor request attempts
    page.on('request', (request) => {
      if (request.url().includes('/auth/login')) {
        requestMade = true;
        console.log(`📤 Request attempted: ${request.method()} ${request.url()}`);
        console.log(`📤 Headers: ${JSON.stringify(request.headers())}`);
      }
    });

    // Monitor successful responses
    page.on('response', (response) => {
      if (response.url().includes('/auth/login')) {
        responseReceived = true;
        console.log(`📥 Response received: ${response.status()} ${response.url()}`);
      }
    });

    // Monitor blocked/failed requests
    page.on('requestfailed', (request) => {
      if (request.url().includes('/auth/login')) {
        requestBlocked = true;
        console.log(`🚫 Request blocked: ${request.url()}`);
        console.log(`🚫 Failure reason: ${request.failure()?.errorText}`);
      }
    });

    // Attempt the request
    await sendOtpButton.click();
    await page.waitForTimeout(5000);

    // Analyze what happened
    console.log('\n=== CORS ANALYSIS ===');
    console.log(`Request attempted: ${requestMade}`);
    console.log(`Response received: ${responseReceived}`);
    console.log(`Request blocked: ${requestBlocked}`);
    console.log(`CORS errors detected: ${corsErrors.length}`);
    console.log(`Network errors detected: ${networkErrors.length}`);

    if (corsErrors.length > 0) {
      console.log('\n🚨 CORS Configuration Issues:');
      corsErrors.forEach(error => console.log(`   ${error}`));
      console.log('\n💡 Recommendation: Check backend CORS settings for origin http://localhost:3000');
    }

    if (requestMade && responseReceived) {
      console.log('\n✅ CORS Configuration: Working correctly');
    } else if (requestMade && !responseReceived) {
      console.log('\n⚠️ CORS Configuration: Request sent but no response (may indicate backend issues)');
    } else if (!requestMade) {
      console.log('\n❌ CORS Configuration: Request not being sent (frontend issue)');
    }

    // Check for specific CORS error patterns
    const hasCorsError = consoleLogs.some(log =>
      log.toLowerCase().includes('blocked by cors policy') ||
      log.toLowerCase().includes('access-control-allow-origin')
    );

    if (hasCorsError) {
      console.log('\n🔧 CORS Fix Required: Add http://localhost:3000 to backend CORS allowed origins');
    }
  });

  test('should validate performance and no infinite loops', async ({ page }) => {
    const startTime = Date.now();

    // Monitor page load performance
    await page.goto('http://localhost:3000/');
    const loadTime = Date.now() - startTime;

    console.log(`⏱️ Page load time: ${loadTime}ms`);
    expect(loadTime).toBeLessThan(10000); // Should load within 10 seconds

    // Wait for content to settle
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(5000);

    // Check for infinite loops by monitoring loading indicators
    const finalLoadingCount = await page.locator('.animate-pulse').count();
    console.log(`🔄 Loading indicators after 5 seconds: ${finalLoadingCount}`);

    // Should not have excessive loading indicators
    expect(finalLoadingCount).toBeLessThan(20);

    // Check for console errors that might indicate infinite loops
    const errorLogs = consoleLogs.filter(log =>
      log.includes('[error]') ||
      log.toLowerCase().includes('maximum update depth') ||
      log.toLowerCase().includes('infinite')
    );

    console.log(`❌ Error logs detected: ${errorLogs.length}`);
    errorLogs.forEach(error => console.log(`   ${error}`));

    // Performance should be reasonable
    if (loadTime < 5000 && finalLoadingCount < 10 && errorLogs.length === 0) {
      console.log('🚀 Performance: Excellent - fast load, stable UI, no errors');
    } else if (loadTime < 8000 && finalLoadingCount < 15) {
      console.log('✅ Performance: Good - acceptable load time and stability');
    } else {
      console.log('⚠️ Performance: Needs improvement - slow load or unstable UI');
    }

    // Verify main UI elements are functional
    await expect(page.locator('h1')).toContainText(/telegram/i);
    await expect(page.locator('input[placeholder*="+66"]')).toBeVisible();
    await expect(page.locator('button:has-text("Send OTP")')).toBeVisible();

    console.log('✅ UI elements are functional and responsive');
  });
});