import { chromium } from 'playwright';

(async () => {
  console.log('🚀 Starting browser automation...\n');
  
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();

  // Track all requests to auth/login
  const authRequests = [];
  page.on('request', request => {
    if (request.url().includes('/auth/login')) {
      const requestData = {
        url: request.url(),
        method: request.method(),
        headers: request.headers(),
        postData: request.postData()
      };
      authRequests.push(requestData);
      
      console.log('📤 REQUEST TO /auth/login');
      console.log('URL:', request.url());
      console.log('Method:', request.method());
      console.log('Body:', request.postData());
      console.log('---\n');
    }
  });

  // Track responses
  page.on('response', async response => {
    if (response.url().includes('/auth/login')) {
      console.log('📥 RESPONSE FROM /auth/login');
      console.log('Status:', response.status(), response.statusText());
      try {
        const body = await response.text();
        console.log('Body:', body);
      } catch (e) {
        console.log('Could not read body:', e.message);
      }
      console.log('---\n');
    }
  });

  console.log('📱 Navigating to http://localhost:3000...');
  await page.goto('http://localhost:3000', { waitUntil: 'networkidle' });
  
  console.log('📸 Taking screenshot of page...');
  await page.screenshot({ path: '/tmp/auth-page.png' });
  
  console.log('🔍 Looking for phone input field...');
  const phoneInput = await page.locator('input').first();
  
  if (phoneInput) {
    console.log('✏️  Found input, filling phone number: +66812345678');
    await phoneInput.fill('+66812345678');
    
    console.log('🔘 Looking for Send OTP button...');
    const sendButton = await page.locator('button').filter({ hasText: /send.*otp/i }).first();
    
    if (sendButton) {
      console.log('🔘 Clicking Send OTP button...');
      await sendButton.click();
      
      console.log('⏳ Waiting for response...\n');
      await page.waitForTimeout(3000);
    }
  }
  
  console.log('\n📊 VERIFICATION SUMMARY:');
  console.log('='.repeat(50));
  
  if (authRequests.length > 0) {
    const lastRequest = authRequests[authRequests.length - 1];
    const postData = JSON.parse(lastRequest.postData || '{}');
    
    console.log('✅ Request was sent to Railway gateway');
    console.log('📍 URL:', lastRequest.url);
    console.log('📦 Payload:', JSON.stringify(postData, null, 2));
    
    if (postData.country_code && postData.country_code.length === 2) {
      console.log('✅ Country code format is CORRECT:', postData.country_code);
      console.log('   (2-character ISO code as required by backend)');
    } else {
      console.log('❌ Country code format is INCORRECT:', postData.country_code);
    }
    
    if (postData.phone_number) {
      console.log('✅ Phone number sent:', postData.phone_number);
    }
  } else {
    console.log('ℹ️  No requests captured - check /tmp/auth-page.png for page state');
  }
  
  console.log('='.repeat(50));
  console.log('\n✨ Verification complete! Closing browser in 3 seconds...');
  
  await page.waitForTimeout(3000);
  await browser.close();
})();
