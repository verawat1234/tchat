import { test, expect, Page, BrowserContext } from '@playwright/test';

/**
 * Viewer Experience E2E Tests
 *
 * Comprehensive end-to-end testing for viewer experience in live streaming
 * Tests complete viewer journey: discovery → join → chat → commerce → leave
 * Covers store streams with product overlay and cart integration
 *
 * Test Flow:
 * 1. Setup: Login as viewer, create active store stream (as broadcaster)
 * 2. Discover streams: GET /api/v1/streams?status=live&stream_type=store
 * 3. Join stream: Navigate to stream page, verify player loads
 * 4. Verify player: Check video element, quality selector, controls
 * 5. Send chat: Type and submit chat message
 * 6. Verify chat: Message appears within 1 second
 * 7. Test rate limiting: Send 6 messages rapidly, verify 6th blocked
 * 8. View product: Check product overlay displays featured product
 * 9. Add to cart: Click "Add to Cart" on product overlay
 * 10. Verify cart: Check product added via GET /api/v1/cart/items
 * 11. Leave stream: Close player, verify session ended
 */

// API Configuration
const API_BASE_URL = process.env.VITE_API_URL || 'http://localhost:8080/api/v1';
const CHAT_RATE_LIMIT_MS = 200;
const CHAT_RATE_LIMIT_COUNT = 5;

// Test Data
interface StreamData {
  id: string;
  broadcaster_id: string;
  title: string;
  stream_type: 'store';
  status: 'live';
  featured_product_id?: string;
}

interface ProductData {
  id: string;
  name: string;
  price: number;
  image_url: string;
  description: string;
}

interface ChatMessage {
  id: string;
  user_id: string;
  username: string;
  message: string;
  timestamp: string;
}

// Mock WebRTC APIs for viewer experience
const mockViewerWebRTC = `
  // Mock getUserMedia for video/audio capture (viewer typically doesn't need this)
  window.navigator.mediaDevices = {
    getUserMedia: async (constraints) => {
      console.log('Mock viewer getUserMedia called:', constraints);
      const stream = new MediaStream();

      if (constraints.video) {
        const videoTrack = {
          kind: 'video',
          enabled: true,
          muted: false,
          readyState: 'live',
          stop: () => console.log('Video track stopped'),
          addEventListener: () => {},
          removeEventListener: () => {}
        };
        stream.addTrack(videoTrack);
      }

      if (constraints.audio) {
        const audioTrack = {
          kind: 'audio',
          enabled: true,
          muted: false,
          readyState: 'live',
          stop: () => console.log('Audio track stopped'),
          addEventListener: () => {},
          removeEventListener: () => {}
        };
        stream.addTrack(audioTrack);
      }

      return stream;
    },
    enumerateDevices: async () => []
  };

  // Mock RTCPeerConnection for viewer-side WebRTC
  window.RTCPeerConnection = class MockViewerRTCPeerConnection {
    constructor(config) {
      console.log('Mock viewer RTCPeerConnection created:', config);
      this.localDescription = null;
      this.remoteDescription = null;
      this.connectionState = 'connected';
      this.iceConnectionState = 'connected';
      this.signalingState = 'stable';
      this.onicecandidate = null;
      this.ontrack = null;
      this.onconnectionstatechange = null;
      this.oniceconnectionstatechange = null;
      this.onsignalingstatechange = null;
      this._streams = new Set();

      // Simulate connection established
      setTimeout(() => {
        if (this.onconnectionstatechange) {
          this.connectionState = 'connected';
          this.onconnectionstatechange({ target: this });
        }
        if (this.oniceconnectionstatechange) {
          this.iceConnectionState = 'connected';
          this.oniceconnectionstatechange({ target: this });
        }
      }, 100);
    }

    async setRemoteDescription(description) {
      console.log('Mock setRemoteDescription called:', description.type);
      this.remoteDescription = description;
      this.signalingState = 'have-remote-offer';

      // Trigger ontrack event to simulate receiving broadcaster's stream
      if (this.ontrack) {
        const stream = new MediaStream();
        const videoTrack = {
          kind: 'video',
          enabled: true,
          muted: false,
          readyState: 'live',
          stop: () => {},
          addEventListener: () => {},
          removeEventListener: () => {}
        };
        stream.addTrack(videoTrack);

        this.ontrack({
          track: videoTrack,
          streams: [stream]
        });
      }
    }

    async createAnswer(options) {
      console.log('Mock createAnswer called');
      return {
        type: 'answer',
        sdp: 'mock-viewer-answer-sdp-' + Date.now()
      };
    }

    async setLocalDescription(description) {
      console.log('Mock setLocalDescription called:', description.type);
      this.localDescription = description;
      this.signalingState = 'stable';
    }

    async addIceCandidate(candidate) {
      console.log('Mock addIceCandidate called');
      return Promise.resolve();
    }

    addTrack(track, stream) {
      console.log('Mock addTrack called:', track.kind);
      this._streams.add(stream);
      return {
        track,
        sender: { track },
        stop: () => {}
      };
    }

    close() {
      console.log('Mock peer connection closed');
      this.connectionState = 'closed';
      this.iceConnectionState = 'closed';
    }
  };

  // Mock WebSocket for chat and signaling
  window.MockWebSocket = class MockWebSocket {
    constructor(url) {
      console.log('Mock WebSocket created:', url);
      this.url = url;
      this.readyState = 1; // OPEN
      this.onopen = null;
      this.onmessage = null;
      this.onerror = null;
      this.onclose = null;

      setTimeout(() => {
        if (this.onopen) this.onopen({ type: 'open' });
      }, 50);
    }

    send(data) {
      console.log('Mock WebSocket send:', data);
      // Simulate echo back for chat messages
      setTimeout(() => {
        if (this.onmessage) {
          const parsed = JSON.parse(data);
          if (parsed.type === 'chat_message') {
            this.onmessage({
              data: JSON.stringify({
                type: 'chat_message',
                data: {
                  id: 'msg-' + Date.now(),
                  user_id: 'viewer-123',
                  username: 'TestViewer',
                  message: parsed.message,
                  timestamp: new Date().toISOString()
                }
              })
            });
          }
        }
      }, 100);
    }

    close() {
      console.log('Mock WebSocket closed');
      this.readyState = 3; // CLOSED
      if (this.onclose) this.onclose({ type: 'close' });
    }
  };
`;

// Helper Functions
async function loginAsViewer(page: Page): Promise<string> {
  // Navigate to login page
  await page.goto('/login');

  // Fill in viewer credentials
  await page.fill('[data-testid="email-input"]', 'viewer@test.com');
  await page.fill('[data-testid="password-input"]', 'viewer123');

  // Click login button
  await page.click('[data-testid="login-button"]');

  // Wait for navigation to complete
  await page.waitForURL('/dashboard', { timeout: 10000 });

  // Extract auth token from localStorage
  const token = await page.evaluate(() => {
    return localStorage.getItem('auth_token') || 'mock-viewer-token';
  });

  return token;
}

async function createLiveStoreStream(
  page: Page,
  broadcasterToken: string
): Promise<StreamData> {
  // Create stream via API
  const response = await page.request.post(`${API_BASE_URL}/streams`, {
    headers: {
      'Authorization': `Bearer ${broadcasterToken}`,
      'Content-Type': 'application/json'
    },
    data: {
      title: 'Live Store Test Stream',
      description: 'Test store stream with featured products',
      stream_type: 'store',
      privacy_setting: 'public',
      featured_product_id: 'product-test-001'
    }
  });

  expect(response.ok()).toBeTruthy();
  const stream = await response.json();

  // Start broadcasting
  const startResponse = await page.request.post(
    `${API_BASE_URL}/streams/${stream.id}/start`,
    {
      headers: {
        'Authorization': `Bearer ${broadcasterToken}`,
        'Content-Type': 'application/json'
      },
      data: {
        webrtc_offer: 'mock-offer-sdp'
      }
    }
  );

  expect(startResponse.ok()).toBeTruthy();

  return {
    ...stream,
    status: 'live' as const
  };
}

async function cleanupStream(
  page: Page,
  streamId: string,
  broadcasterToken: string
): Promise<void> {
  // End stream
  await page.request.post(`${API_BASE_URL}/streams/${streamId}/end`, {
    headers: {
      'Authorization': `Bearer ${broadcasterToken}`,
      'Content-Type': 'application/json'
    }
  });

  // Delete stream
  await page.request.delete(`${API_BASE_URL}/streams/${streamId}`, {
    headers: {
      'Authorization': `Bearer ${broadcasterToken}`
    }
  });
}

async function cleanupCart(page: Page, viewerToken: string): Promise<void> {
  // Clear cart items
  await page.request.delete(`${API_BASE_URL}/cart/items`, {
    headers: {
      'Authorization': `Bearer ${viewerToken}`
    }
  });
}

// Test Suite
test.describe('Viewer Experience E2E', () => {
  let viewerToken: string;
  let broadcasterToken: string;
  let testStream: StreamData;

  test.beforeEach(async ({ page, context }) => {
    // Inject WebRTC mocks
    await context.addInitScript(mockViewerWebRTC);

    // Login as viewer
    viewerToken = await loginAsViewer(page);

    // Mock broadcaster token (in real scenario, would create broadcaster account)
    broadcasterToken = 'mock-broadcaster-token';

    // Create live store stream
    testStream = await createLiveStoreStream(page, broadcasterToken);
  });

  test.afterEach(async ({ page }) => {
    // Cleanup stream and cart
    await cleanupStream(page, testStream.id, broadcasterToken);
    await cleanupCart(page, viewerToken);
  });

  test('Complete viewer journey: discover → join → chat → commerce → leave', async ({
    page
  }) => {
    // Step 1: Discover live streams
    await test.step('Discover live store streams', async () => {
      // Navigate to streams page
      await page.goto('/streams');

      // Wait for streams list to load
      await page.waitForSelector('[data-testid="streams-list"]', {
        timeout: 5000
      });

      // Apply filter for live store streams
      await page.click('[data-testid="filter-type-store"]');
      await page.click('[data-testid="filter-status-live"]');

      // Verify test stream appears in list
      const streamCard = page.locator(
        `[data-testid="stream-card-${testStream.id}"]`
      );
      await expect(streamCard).toBeVisible();

      // Verify stream metadata
      await expect(streamCard.locator('[data-testid="stream-title"]')).toHaveText(
        testStream.title
      );
      await expect(
        streamCard.locator('[data-testid="stream-status"]')
      ).toHaveText('LIVE');
    });

    // Step 2: Join stream
    await test.step('Join stream and verify player loads', async () => {
      // Click on stream to join
      await page.click(`[data-testid="stream-card-${testStream.id}"]`);

      // Wait for stream page navigation
      await page.waitForURL(`/streams/${testStream.id}`, { timeout: 5000 });

      // Verify video player renders
      const videoPlayer = page.locator('[data-testid="video-player"]');
      await expect(videoPlayer).toBeVisible();

      // Verify video element is present
      const videoElement = page.locator('video[data-testid="video-element"]');
      await expect(videoElement).toBeVisible();

      // Verify quality selector is available
      const qualitySelector = page.locator(
        '[data-testid="quality-selector"]'
      );
      await expect(qualitySelector).toBeVisible();

      // Verify controls are present
      const playPauseBtn = page.locator('[data-testid="play-pause-button"]');
      const volumeControl = page.locator('[data-testid="volume-control"]');
      const fullscreenBtn = page.locator('[data-testid="fullscreen-button"]');

      await expect(playPauseBtn).toBeVisible();
      await expect(volumeControl).toBeVisible();
      await expect(fullscreenBtn).toBeVisible();
    });

    // Step 3: Send chat message
    await test.step('Send chat message', async () => {
      // Wait for chat panel to load
      const chatPanel = page.locator('[data-testid="chat-panel"]');
      await expect(chatPanel).toBeVisible();

      // Type chat message
      const chatInput = page.locator('[data-testid="chat-input"]');
      await chatInput.fill('Hello from the viewer!');

      // Submit message
      await page.click('[data-testid="chat-send-button"]');

      // Verify message appears in chat list within 1 second
      const chatMessage = page.locator(
        '[data-testid="chat-message"]:has-text("Hello from the viewer!")'
      );
      await expect(chatMessage).toBeVisible({ timeout: 1000 });

      // Verify message metadata
      await expect(
        chatMessage.locator('[data-testid="message-username"]')
      ).toHaveText('TestViewer');
    });

    // Step 4: Test chat rate limiting
    await test.step('Test chat rate limiting', async () => {
      const chatInput = page.locator('[data-testid="chat-input"]');
      const sendButton = page.locator('[data-testid="chat-send-button"]');

      // Send messages rapidly (should trigger rate limit on 6th)
      for (let i = 1; i <= 6; i++) {
        await chatInput.fill(`Rapid message ${i}`);
        await sendButton.click();

        if (i === 6) {
          // 6th message should show rate limit UI
          const rateLimitNotice = page.locator(
            '[data-testid="rate-limit-notice"]'
          );
          await expect(rateLimitNotice).toBeVisible({ timeout: 500 });
          await expect(rateLimitNotice).toContainText('Too many messages');

          // Send button should be disabled
          await expect(sendButton).toBeDisabled();

          // Wait for rate limit to clear
          await page.waitForTimeout(CHAT_RATE_LIMIT_MS + 100);

          // Verify button re-enabled
          await expect(sendButton).toBeEnabled();
          await expect(rateLimitNotice).not.toBeVisible();
        } else {
          // First 5 messages should send successfully
          const message = page.locator(
            `[data-testid="chat-message"]:has-text("Rapid message ${i}")`
          );
          await expect(message).toBeVisible({ timeout: 1000 });
        }
      }
    });

    // Step 5: View product overlay
    await test.step('View product overlay', async () => {
      // Wait for product overlay to appear
      const productOverlay = page.locator('[data-testid="product-overlay"]');
      await expect(productOverlay).toBeVisible({ timeout: 3000 });

      // Verify featured product details
      const productName = productOverlay.locator(
        '[data-testid="product-name"]'
      );
      const productPrice = productOverlay.locator(
        '[data-testid="product-price"]'
      );
      const productImage = productOverlay.locator(
        '[data-testid="product-image"]'
      );
      const addToCartBtn = productOverlay.locator(
        '[data-testid="add-to-cart-button"]'
      );

      await expect(productName).toBeVisible();
      await expect(productPrice).toBeVisible();
      await expect(productImage).toBeVisible();
      await expect(addToCartBtn).toBeVisible();
      await expect(addToCartBtn).toBeEnabled();
    });

    // Step 6: Add product to cart
    await test.step('Add product to cart', async () => {
      // Click "Add to Cart" button
      const addToCartBtn = page.locator('[data-testid="add-to-cart-button"]');
      await addToCartBtn.click();

      // Verify success notification
      const successNotification = page.locator(
        '[data-testid="cart-success-notification"]'
      );
      await expect(successNotification).toBeVisible({ timeout: 2000 });
      await expect(successNotification).toContainText('Added to cart');

      // Verify cart badge updates
      const cartBadge = page.locator('[data-testid="cart-badge"]');
      await expect(cartBadge).toHaveText('1');
    });

    // Step 7: Verify cart contents
    await test.step('Verify product added to cart via API', async () => {
      // Fetch cart items via API
      const response = await page.request.get(`${API_BASE_URL}/cart/items`, {
        headers: {
          'Authorization': `Bearer ${viewerToken}`
        }
      });

      expect(response.ok()).toBeTruthy();
      const cartData = await response.json();

      // Verify cart contains the product
      expect(cartData.items).toHaveLength(1);
      expect(cartData.items[0].product_id).toBe('product-test-001');
      expect(cartData.items[0].quantity).toBe(1);
    });

    // Step 8: Leave stream
    await test.step('Leave stream and verify session ended', async () => {
      // Get initial viewer count
      const viewerCount = page.locator('[data-testid="viewer-count"]');
      const initialCount = await viewerCount.textContent();

      // Close video player
      const closeButton = page.locator('[data-testid="close-player-button"]');
      await closeButton.click();

      // Wait for navigation away from stream page
      await page.waitForURL((url) => !url.pathname.includes('/streams/'), {
        timeout: 5000
      });

      // Verify viewer session ended (check via API)
      const streamResponse = await page.request.get(
        `${API_BASE_URL}/streams/${testStream.id}`,
        {
          headers: {
            'Authorization': `Bearer ${viewerToken}`
          }
        }
      );

      expect(streamResponse.ok()).toBeTruthy();
      const streamData = await streamResponse.json();

      // Viewer count should have decremented
      expect(streamData.current_viewer_count).toBeLessThan(
        parseInt(initialCount || '0')
      );
    });
  });

  test('Viewer can navigate between streams without cart loss', async ({
    page
  }) => {
    // Add product to cart from first stream
    await page.goto(`/streams/${testStream.id}`);
    await page.waitForSelector('[data-testid="video-player"]');

    const addToCartBtn = page.locator('[data-testid="add-to-cart-button"]');
    await addToCartBtn.click();

    // Verify cart badge shows 1 item
    const cartBadge = page.locator('[data-testid="cart-badge"]');
    await expect(cartBadge).toHaveText('1');

    // Navigate back to streams list
    await page.goto('/streams');
    await page.waitForSelector('[data-testid="streams-list"]');

    // Verify cart persists
    await expect(cartBadge).toHaveText('1');

    // Navigate to different stream
    await page.goto(`/streams/${testStream.id}`);

    // Verify cart still shows 1 item
    await expect(cartBadge).toHaveText('1');
  });

  test('Viewer receives real-time chat messages from other viewers', async ({
    page,
    context
  }) => {
    // Open stream in first tab (viewer 1)
    await page.goto(`/streams/${testStream.id}`);
    await page.waitForSelector('[data-testid="video-player"]');

    // Open stream in second tab (viewer 2)
    const page2 = await context.newPage();
    await page2.addInitScript(mockViewerWebRTC);
    await page2.goto(`/streams/${testStream.id}`);
    await page2.waitForSelector('[data-testid="video-player"]');

    // Send message from viewer 2
    const chatInput2 = page2.locator('[data-testid="chat-input"]');
    await chatInput2.fill('Hello from viewer 2!');
    await page2.click('[data-testid="chat-send-button"]');

    // Verify message appears in viewer 1's chat
    const chatMessage = page.locator(
      '[data-testid="chat-message"]:has-text("Hello from viewer 2!")'
    );
    await expect(chatMessage).toBeVisible({ timeout: 2000 });

    // Close second tab
    await page2.close();
  });

  test('Viewer can recover from connection loss', async ({ page }) => {
    // Join stream
    await page.goto(`/streams/${testStream.id}`);
    await page.waitForSelector('[data-testid="video-player"]');

    // Simulate connection loss
    await page.evaluate(() => {
      // @ts-ignore
      if (window.mockPeerConnection) {
        // @ts-ignore
        window.mockPeerConnection.connectionState = 'disconnected';
        // @ts-ignore
        if (window.mockPeerConnection.onconnectionstatechange) {
          // @ts-ignore
          window.mockPeerConnection.onconnectionstatechange({
            // @ts-ignore
            target: window.mockPeerConnection
          });
        }
      }
    });

    // Verify reconnection notice appears
    const reconnectNotice = page.locator('[data-testid="reconnect-notice"]');
    await expect(reconnectNotice).toBeVisible({ timeout: 3000 });
    await expect(reconnectNotice).toContainText('Reconnecting');

    // Simulate reconnection
    await page.evaluate(() => {
      // @ts-ignore
      if (window.mockPeerConnection) {
        // @ts-ignore
        window.mockPeerConnection.connectionState = 'connected';
        // @ts-ignore
        if (window.mockPeerConnection.onconnectionstatechange) {
          // @ts-ignore
          window.mockPeerConnection.onconnectionstatechange({
            // @ts-ignore
            target: window.mockPeerConnection
          });
        }
      }
    });

    // Verify reconnection successful
    await expect(reconnectNotice).not.toBeVisible({ timeout: 3000 });

    // Verify video player still functional
    const videoPlayer = page.locator('[data-testid="video-player"]');
    await expect(videoPlayer).toBeVisible();
  });

  test('Accessibility: keyboard navigation and screen reader support', async ({
    page
  }) => {
    // Join stream
    await page.goto(`/streams/${testStream.id}`);
    await page.waitForSelector('[data-testid="video-player"]');

    // Test keyboard navigation to video controls
    await page.keyboard.press('Tab'); // Focus play/pause
    const playPauseBtn = page.locator('[data-testid="play-pause-button"]');
    await expect(playPauseBtn).toBeFocused();

    await page.keyboard.press('Tab'); // Focus volume control
    const volumeControl = page.locator('[data-testid="volume-control"]');
    await expect(volumeControl).toBeFocused();

    // Test keyboard navigation to chat
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    const chatInput = page.locator('[data-testid="chat-input"]');
    await expect(chatInput).toBeFocused();

    // Test keyboard message send
    await chatInput.type('Keyboard navigation test');
    await page.keyboard.press('Enter');

    const chatMessage = page.locator(
      '[data-testid="chat-message"]:has-text("Keyboard navigation test")'
    );
    await expect(chatMessage).toBeVisible({ timeout: 1000 });

    // Verify ARIA labels
    await expect(playPauseBtn).toHaveAttribute('aria-label', 'Play or pause');
    await expect(volumeControl).toHaveAttribute('aria-label', 'Volume control');
    await expect(chatInput).toHaveAttribute('aria-label', 'Chat message input');
  });

  test('Performance: stream loads within 3 seconds', async ({ page }) => {
    const startTime = Date.now();

    // Navigate to stream
    await page.goto(`/streams/${testStream.id}`);

    // Wait for video player to load
    await page.waitForSelector('[data-testid="video-player"]');
    await page.waitForSelector('video[data-testid="video-element"]');

    const loadTime = Date.now() - startTime;

    // Verify load time under 3 seconds
    expect(loadTime).toBeLessThan(3000);

    // Verify no console errors
    const errors: string[] = [];
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      }
    });

    await page.waitForTimeout(1000);
    expect(errors).toHaveLength(0);
  });

  test('Mobile: responsive viewer experience', async ({ page, context }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });

    // Join stream
    await page.goto(`/streams/${testStream.id}`);
    await page.waitForSelector('[data-testid="video-player"]');

    // Verify mobile-optimized layout
    const videoPlayer = page.locator('[data-testid="video-player"]');
    const chatPanel = page.locator('[data-testid="chat-panel"]');
    const productOverlay = page.locator('[data-testid="product-overlay"]');

    // Video player should fill viewport width
    const playerBox = await videoPlayer.boundingBox();
    expect(playerBox?.width).toBeCloseTo(375, 10);

    // Chat panel should be collapsible on mobile
    const chatToggle = page.locator('[data-testid="chat-toggle-mobile"]');
    await expect(chatToggle).toBeVisible();

    // Toggle chat visibility
    await chatToggle.click();
    await expect(chatPanel).toHaveClass(/collapsed/);

    await chatToggle.click();
    await expect(chatPanel).not.toHaveClass(/collapsed/);

    // Product overlay should be responsive
    const overlayBox = await productOverlay.boundingBox();
    expect(overlayBox?.width).toBeLessThanOrEqual(375);

    // Touch gestures for video controls
    await page.tap('[data-testid="video-player"]');
    const controlsOverlay = page.locator('[data-testid="controls-overlay"]');
    await expect(controlsOverlay).toBeVisible();

    // Controls should auto-hide after 3 seconds
    await page.waitForTimeout(3000);
    await expect(controlsOverlay).not.toBeVisible();
  });
});

/**
 * Test Summary:
 *
 * 1. Complete Viewer Journey Test (Primary)
 *    - Discover live store streams
 *    - Join stream and verify player
 *    - Send chat messages with rate limiting
 *    - View and interact with product overlay
 *    - Add product to cart
 *    - Verify cart via API
 *    - Leave stream gracefully
 *
 * 2. Cart Persistence Test
 *    - Verify cart persists across navigation
 *    - Cart state maintained between streams
 *
 * 3. Real-time Chat Test
 *    - Multi-viewer chat synchronization
 *    - Real-time message broadcasting
 *
 * 4. Connection Recovery Test
 *    - Handle disconnection gracefully
 *    - Automatic reconnection logic
 *    - User feedback during reconnection
 *
 * 5. Accessibility Test
 *    - Keyboard navigation support
 *    - ARIA labels for screen readers
 *    - Focus management
 *
 * 6. Performance Test
 *    - Stream loads within 3 seconds
 *    - No console errors during playback
 *
 * 7. Mobile Responsiveness Test
 *    - Responsive layout on mobile devices
 *    - Touch gesture support
 *    - Collapsible UI elements
 *    - Auto-hiding controls
 *
 * Coverage:
 * - User discovery and joining flows
 * - WebRTC viewer connection
 * - Real-time chat with rate limiting
 * - Commerce integration (products, cart)
 * - Session management
 * - Error recovery
 * - Accessibility compliance
 * - Performance validation
 * - Mobile user experience
 */