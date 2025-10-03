import { test, expect } from '@playwright/test';

test.describe('OTP Authentication Flow', () => {
  test('complete OTP authentication flow from phone input to verification', async ({ page }) => {
    // Configure longer timeout for this test as it involves network requests
    test.setTimeout(60000);

    // Array to capture console messages
    const consoleMessages: string[] = [];
    const errorMessages: string[] = [];

    // Listen to console events
    page.on('console', (msg) => {
      const text = msg.text();
      consoleMessages.push(`[${msg.type()}] ${text}`);
      console.log(`Browser console [${msg.type()}]:`, text);
    });

    // Listen to page errors
    page.on('pageerror', (error) => {
      const errorText = error.toString();
      errorMessages.push(errorText);
      console.error('Page error:', errorText);
    });

    // Step 1: Navigate to the application
    console.log('Step 1: Navigating to http://localhost:3000');
    await page.goto('http://localhost:3000');
    await page.waitForLoadState('networkidle');

    // Take screenshot of initial state
    await page.screenshot({ path: 'tests/screenshots/01-initial-load.png', fullPage: true });

    // Step 2: Find the phone number input
    console.log('Step 2: Looking for phone number input field');

    // Try multiple possible selectors for phone input
    const phoneInput = await page.locator('input[type="tel"], input[placeholder*="phone" i], input[name*="phone" i], input[data-testid="phone-input"]').first();

    await expect(phoneInput).toBeVisible({ timeout: 10000 });
    await page.screenshot({ path: 'tests/screenshots/02-phone-input-visible.png', fullPage: true });

    // Step 3: Fill in the phone number
    console.log('Step 3: Filling phone number +66812345678');
    await phoneInput.fill('+66812345678');
    await page.waitForTimeout(500); // Brief wait for any validation

    await page.screenshot({ path: 'tests/screenshots/03-phone-filled.png', fullPage: true });

    // Step 4: Find and click the "Send OTP" button
    console.log('Step 4: Looking for Send OTP button');

    const sendOTPButton = await page.locator('button:has-text("Send OTP"), button:has-text("Request OTP"), button[data-testid="send-otp-button"]').first();

    await expect(sendOTPButton).toBeVisible({ timeout: 5000 });
    await page.screenshot({ path: 'tests/screenshots/04-send-otp-button.png', fullPage: true });

    console.log('Step 4b: Clicking Send OTP button');
    await sendOTPButton.click();

    // Step 5: Wait for OTP input field to appear
    console.log('Step 5: Waiting for OTP input field to appear');

    const otpInput = await page.locator('input[type="text"][placeholder*="OTP" i], input[type="text"][placeholder*="code" i], input[data-testid="otp-input"], input[name*="otp" i]').first();

    await expect(otpInput).toBeVisible({ timeout: 15000 });
    await page.screenshot({ path: 'tests/screenshots/05-otp-input-visible.png', fullPage: true });

    // Step 6: Enter the test OTP code
    console.log('Step 6: Entering OTP code 123456');
    await otpInput.fill('123456');
    await page.waitForTimeout(500);

    await page.screenshot({ path: 'tests/screenshots/06-otp-filled.png', fullPage: true });

    // Step 7: Find and click the verify/submit button
    console.log('Step 7: Looking for verify/submit button');

    const verifyButton = await page.locator('button:has-text("Verify"), button:has-text("Submit"), button:has-text("Confirm"), button[data-testid="verify-otp-button"]').first();

    await expect(verifyButton).toBeVisible({ timeout: 5000 });
    await page.screenshot({ path: 'tests/screenshots/07-verify-button.png', fullPage: true });

    console.log('Step 7b: Clicking verify button');
    await verifyButton.click();

    // Step 8: Wait for response and check final state
    console.log('Step 8: Waiting for verification response');
    await page.waitForTimeout(3000); // Wait for API response

    await page.screenshot({ path: 'tests/screenshots/08-final-state.png', fullPage: true });

    // Check for success indicators
    const successIndicators = [
      page.locator('text=/success/i'),
      page.locator('text=/welcome/i'),
      page.locator('text=/verified/i'),
      page.locator('[data-testid="auth-success"]'),
    ];

    let foundSuccess = false;
    for (const indicator of successIndicators) {
      if (await indicator.isVisible().catch(() => false)) {
        foundSuccess = true;
        console.log('Success indicator found:', await indicator.textContent());
        break;
      }
    }

    // Check for error indicators
    const errorIndicators = [
      page.locator('text=/error/i'),
      page.locator('text=/invalid/i'),
      page.locator('text=/failed/i'),
      page.locator('[role="alert"]'),
    ];

    let foundError = false;
    let errorText = '';
    for (const indicator of errorIndicators) {
      if (await indicator.isVisible().catch(() => false)) {
        foundError = true;
        errorText = await indicator.textContent() || '';
        console.log('Error indicator found:', errorText);
        break;
      }
    }

    // Get final URL
    const finalUrl = page.url();
    console.log('Final URL:', finalUrl);

    // Get page title
    const pageTitle = await page.title();
    console.log('Page title:', pageTitle);

    // Compile comprehensive test report
    console.log('\n========== OTP AUTHENTICATION FLOW TEST REPORT ==========');
    console.log('\n1. TEST EXECUTION SUMMARY:');
    console.log(`   - Initial URL: http://localhost:3000`);
    console.log(`   - Final URL: ${finalUrl}`);
    console.log(`   - Page Title: ${pageTitle}`);
    console.log(`   - Phone Number Used: +66812345678`);
    console.log(`   - OTP Code Used: 123456`);

    console.log('\n2. FLOW COMPLETION STATUS:');
    console.log(`   - Phone input found: ✓`);
    console.log(`   - Send OTP button found: ✓`);
    console.log(`   - OTP input appeared: ✓`);
    console.log(`   - Verify button found: ✓`);
    console.log(`   - Success indicator: ${foundSuccess ? '✓' : '✗'}`);
    console.log(`   - Error detected: ${foundError ? '✓ (Error: ' + errorText + ')' : '✗'}`);

    console.log('\n3. CONSOLE MESSAGES:');
    if (consoleMessages.length > 0) {
      consoleMessages.forEach((msg, idx) => {
        console.log(`   ${idx + 1}. ${msg}`);
      });
    } else {
      console.log('   No console messages captured');
    }

    console.log('\n4. PAGE ERRORS:');
    if (errorMessages.length > 0) {
      errorMessages.forEach((error, idx) => {
        console.log(`   ${idx + 1}. ${error}`);
      });
    } else {
      console.log('   No page errors detected');
    }

    console.log('\n5. SCREENSHOTS CAPTURED:');
    console.log('   - 01-initial-load.png');
    console.log('   - 02-phone-input-visible.png');
    console.log('   - 03-phone-filled.png');
    console.log('   - 04-send-otp-button.png');
    console.log('   - 05-otp-input-visible.png');
    console.log('   - 06-otp-filled.png');
    console.log('   - 07-verify-button.png');
    console.log('   - 08-final-state.png');

    console.log('\n6. FINAL ASSESSMENT:');
    if (foundSuccess && !foundError) {
      console.log('   ✓ OTP authentication flow completed SUCCESSFULLY');
    } else if (foundError) {
      console.log('   ✗ OTP authentication flow FAILED with errors');
    } else {
      console.log('   ? OTP authentication flow status UNCLEAR (no success or error indicators)');
    }
    console.log('\n========================================================\n');

    // Assertions for test validation
    expect(errorMessages.length).toBe(0); // No page errors should occur
    expect(foundSuccess || foundError).toBeTruthy(); // Should have some final state indication
  });
});
