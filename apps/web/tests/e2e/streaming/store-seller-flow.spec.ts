import { test, expect } from '@playwright/test';

/**
 * E2E Tests for Store Seller Live Stream Workflow
 *
 * Task: T064
 * Description: Complete store seller flow from stream creation to end
 * Reference: /specs/029-implement-live-on/quickstart.md lines 11-102
 *
 * Test Flow:
 * 1. Login as store seller (KYC Tier 1+)
 * 2. Create scheduled stream
 * 3. Start live stream with WebRTC
 * 4. Feature products during stream
 * 5. Verify player UI and product overlay
 * 6. Send and moderate chat messages
 * 7. End stream
 * 8. Verify recording availability
 *
 * Performance Targets:
 * - Stream creation: <500ms
 * - Stream start: <2s
 * - Product overlay: <2s
 * - Chat moderation: <1s
 * - Recording availability: <30s
 */

// Test data
const STORE_SELLER = {
  email: 'seller@tchat.com',
  password: 'StorePass123!',
  kyc_tier: 1,
  user_id: '550e8400-e29b-41d4-a716-446655440000',
};

const BUYER = {
  email: 'buyer@tchat.com',
  password: 'BuyerPass123!',
  user_id: '660e8400-e29b-41d4-a716-446655440001',
};

const TEST_STREAM = {
  title: 'New iPhone 15 Pro Unboxing & Demo',
  description: 'Live demonstration of iPhone 15 Pro features and accessories',
  stream_type: 'store',
  privacy_setting: 'public',
};

const TEST_PRODUCT = {
  product_id: 'prod-iphone15-pro-256gb',
  display_position: 'overlay',
  display_priority: 1,
};

// Mock WebRTC SDP offer/answer
const MOCK_WEBRTC_OFFER = 'v=0\r\no=- 123 0 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=msid-semantic: WMS\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:mock\r\na=ice-pwd:mockpassword\r\na=fingerprint:sha-256 00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00\r\na=setup:actpass\r\na=mid:0\r\na=sendrecv\r\na=rtcp-mux\r\na=rtpmap:96 VP8/90000';

const MOCK_WEBRTC_ANSWER = 'v=0\r\no=- 456 0 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=msid-semantic: WMS\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:server\r\na=ice-pwd:serverpassword\r\na=fingerprint:sha-256 11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11\r\na=setup:active\r\na=mid:0\r\na=sendrecv\r\na=rtcp-mux\r\na=rtpmap:96 VP8/90000';

// Helper: Generate future timestamp for scheduled stream
function getFutureTimestamp(minutesFromNow: number = 30): string {
  const now = new Date();
  now.setMinutes(now.getMinutes() + minutesFromNow);
  return now.toISOString();
}

// Helper: Check if timestamp is approximately 30 days in the future
function isApproximately30DaysLater(timestamp: string, referenceTime: string): boolean {
  const target = new Date(timestamp).getTime();
  const reference = new Date(referenceTime).getTime();
  const diff = target - reference;
  const thirtyDays = 30 * 24 * 60 * 60 * 1000;
  // Allow 1 hour tolerance
  return Math.abs(diff - thirtyDays) < 60 * 60 * 1000;
}

test.describe('Store Seller Live Stream E2E Flow', () => {
  let authToken: string;
  let streamId: string;
  let buyerAuthToken: string;
  let chatMessageId: string;

  test.beforeAll(async ({ request }) => {
    // Setup: Authenticate store seller
    const loginResponse = await request.post('http://localhost:8080/api/v1/auth/login', {
      data: {
        email: STORE_SELLER.email,
        password: STORE_SELLER.password,
      },
    });

    expect(loginResponse.ok()).toBeTruthy();
    const loginData = await loginResponse.json();
    authToken = loginData.access_token;

    // Setup: Authenticate buyer for chat testing
    const buyerLoginResponse = await request.post('http://localhost:8080/api/v1/auth/login', {
      data: {
        email: BUYER.email,
        password: BUYER.password,
      },
    });

    expect(buyerLoginResponse.ok()).toBeTruthy();
    const buyerLoginData = await buyerLoginResponse.json();
    buyerAuthToken = buyerLoginData.access_token;
  });

  test.beforeEach(async ({ page }) => {
    // Mock WebRTC APIs for browser testing
    await page.addInitScript(() => {
      // Mock MediaDevices.getUserMedia
      Object.defineProperty(navigator.mediaDevices, 'getUserMedia', {
        value: async (constraints: MediaStreamConstraints) => {
          const mockStream = new MediaStream();

          if (constraints.audio) {
            const audioTrack = {
              kind: 'audio',
              enabled: true,
              id: 'mock-audio-track',
              label: 'Mock Audio Track',
              muted: false,
              readyState: 'live',
              stop: () => {},
              addEventListener: () => {},
              removeEventListener: () => {},
              contentHint: '',
              onended: null,
              onmute: null,
              onunmute: null,
              applyConstraints: async () => {},
              clone: () => ({} as MediaStreamTrack),
              getCapabilities: () => ({}),
              getConstraints: () => ({}),
              getSettings: () => ({}),
            } as unknown as MediaStreamTrack;
            mockStream.addTrack(audioTrack);
          }

          if (constraints.video) {
            const videoTrack = {
              kind: 'video',
              enabled: true,
              id: 'mock-video-track',
              label: 'Mock Video Track',
              muted: false,
              readyState: 'live',
              stop: () => {},
              addEventListener: () => {},
              removeEventListener: () => {},
              contentHint: '',
              onended: null,
              onmute: null,
              onunmute: null,
              applyConstraints: async () => {},
              clone: () => ({} as MediaStreamTrack),
              getCapabilities: () => ({}),
              getConstraints: () => ({}),
              getSettings: () => ({}),
            } as unknown as MediaStreamTrack;
            mockStream.addTrack(videoTrack);
          }

          return mockStream;
        },
        writable: true,
      });

      // Mock RTCPeerConnection
      (window as any).RTCPeerConnection = class MockRTCPeerConnection {
        localDescription: RTCSessionDescription | null = null;
        remoteDescription: RTCSessionDescription | null = null;
        connectionState: RTCPeerConnectionState = 'new';
        iceConnectionState: RTCIceConnectionState = 'new';
        signalingState: RTCSignalingState = 'stable';
        onicecandidate: ((event: RTCPeerConnectionIceEvent) => void) | null = null;
        onconnectionstatechange: (() => void) | null = null;
        ontrack: ((event: RTCTrackEvent) => void) | null = null;

        async createOffer(): Promise<RTCSessionDescriptionInit> {
          return {
            type: 'offer',
            sdp: (window as any).__MOCK_WEBRTC_OFFER || 'mock-offer-sdp',
          };
        }

        async createAnswer(): Promise<RTCSessionDescriptionInit> {
          return {
            type: 'answer',
            sdp: (window as any).__MOCK_WEBRTC_ANSWER || 'mock-answer-sdp',
          };
        }

        async setLocalDescription(description: RTCSessionDescriptionInit): Promise<void> {
          this.localDescription = description as RTCSessionDescription;
          // Simulate ICE candidate generation
          setTimeout(() => {
            if (this.onicecandidate) {
              this.onicecandidate({ candidate: null } as RTCPeerConnectionIceEvent);
            }
          }, 100);
        }

        async setRemoteDescription(description: RTCSessionDescriptionInit): Promise<void> {
          this.remoteDescription = description as RTCSessionDescription;
          this.connectionState = 'connected';
          this.iceConnectionState = 'connected';
          if (this.onconnectionstatechange) {
            this.onconnectionstatechange();
          }
        }

        addStream(): void {}
        addTrack(): RTCRtpSender {
          return {} as RTCRtpSender;
        }
        close(): void {
          this.connectionState = 'closed';
          this.iceConnectionState = 'closed';
        }
        addEventListener(): void {}
        removeEventListener(): void {}
      };
    });

    // Set mock WebRTC offer for use in tests
    await page.evaluate((offer) => {
      (window as any).__MOCK_WEBRTC_OFFER = offer;
    }, MOCK_WEBRTC_OFFER);
  });

  test.afterEach(async ({ request }) => {
    // Cleanup: End stream if it exists
    if (streamId && authToken) {
      try {
        await request.post(`http://localhost:8080/api/v1/streams/${streamId}/end`, {
          headers: {
            'Authorization': `Bearer ${authToken}`,
          },
        });
      } catch (error) {
        // Ignore cleanup errors
        console.log('Cleanup: Stream may already be ended');
      }
    }
  });

  test('should complete full store seller live stream workflow', async ({ page, request }) => {
    /**
     * STEP 1: Verify KYC Status
     */
    const userResponse = await request.get('http://localhost:8080/api/v1/users/me', {
      headers: {
        'Authorization': `Bearer ${authToken}`,
      },
    });

    expect(userResponse.ok()).toBeTruthy();
    const userData = await userResponse.json();
    expect(userData.kyc_tier).toBeGreaterThanOrEqual(1);
    expect(userData.email_verified).toBe(true);
    console.log('✓ Step 1: KYC status verified (Tier 1+)');

    /**
     * STEP 2: Create Scheduled Stream
     */
    const scheduledStartTime = getFutureTimestamp(30);
    const createStreamStart = Date.now();

    const createStreamResponse = await request.post('http://localhost:8080/api/v1/streams', {
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
      data: {
        ...TEST_STREAM,
        scheduled_start_time: scheduledStartTime,
      },
    });

    const createStreamDuration = Date.now() - createStreamStart;
    expect(createStreamResponse.status()).toBe(201);
    expect(createStreamDuration).toBeLessThan(500); // <500ms target

    const streamData = await createStreamResponse.json();
    streamId = streamData.id;

    // Validate stream creation response
    expect(streamData.broadcaster_id).toBeTruthy();
    expect(streamData.stream_type).toBe('store');
    expect(streamData.status).toBe('scheduled');
    expect(streamData.stream_key).toMatch(/^rtmp:\/\//);
    expect(streamData.max_capacity).toBe(50000); // Mega-scale support
    expect(streamData.title).toBe(TEST_STREAM.title);

    console.log(`✓ Step 2: Stream created (${createStreamDuration}ms) - ID: ${streamId}`);

    /**
     * STEP 3: Start Live Stream with WebRTC
     */
    const startStreamStart = Date.now();

    const startStreamResponse = await request.post(
      `http://localhost:8080/api/v1/streams/${streamId}/start`,
      {
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json',
        },
        data: {
          webrtc_offer: MOCK_WEBRTC_OFFER,
        },
      }
    );

    const startStreamDuration = Date.now() - startStreamStart;
    expect(startStreamResponse.ok()).toBeTruthy();
    expect(startStreamDuration).toBeLessThan(2000); // <2s target

    const startData = await startStreamResponse.json();

    // Validate stream start response
    expect(startData.webrtc_answer).toBeTruthy();
    expect(startData.webrtc_session_id).toBeTruthy();
    expect(startData.primary_server_id).toBeTruthy();
    expect(startData.quality_layers).toEqual(['360p', '720p', '1080p']);

    console.log(`✓ Step 3: Stream started (${startStreamDuration}ms) - Session: ${startData.webrtc_session_id}`);

    // Verify stream status transition
    const streamStatusResponse = await request.get(
      `http://localhost:8080/api/v1/streams/${streamId}`,
      {
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      }
    );

    expect(streamStatusResponse.ok()).toBeTruthy();
    const streamStatus = await streamStatusResponse.json();
    expect(streamStatus.status).toBe('live');
    expect(streamStatus.actual_start_time).toBeTruthy();

    console.log('✓ Step 3a: Stream status transitioned to "live"');

    /**
     * STEP 4: Feature Products During Stream
     */
    const featureProductStart = Date.now();

    const featureProductResponse = await request.post(
      `http://localhost:8080/api/v1/streams/${streamId}/products`,
      {
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json',
        },
        data: TEST_PRODUCT,
      }
    );

    const featureProductDuration = Date.now() - featureProductStart;
    expect(featureProductResponse.status()).toBe(201);
    expect(featureProductDuration).toBeLessThan(2000); // <2s target

    const featuredProduct = await featureProductResponse.json();

    // Validate featured product response
    expect(featuredProduct.id).toBeTruthy();
    expect(featuredProduct.stream_id).toBe(streamId);
    expect(featuredProduct.product_id).toBe(TEST_PRODUCT.product_id);
    expect(featuredProduct.display_position).toBe('overlay');
    expect(featuredProduct.view_count).toBe(0);
    expect(featuredProduct.click_count).toBe(0);

    console.log(`✓ Step 4: Product featured (${featureProductDuration}ms) - ID: ${featuredProduct.id}`);

    /**
     * STEP 5: Verify Player UI and Product Overlay
     */
    await page.goto(`http://localhost:3000/streams/${streamId}`);

    // Wait for video player to load
    const videoElement = page.locator('[data-testid="stream-video-player"]');
    await expect(videoElement).toBeVisible({ timeout: 5000 });

    // Verify quality selector
    const qualitySelector = page.locator('[data-testid="quality-selector"]');
    await expect(qualitySelector).toBeVisible();

    // Check available quality options
    const qualityOptions = page.locator('[data-testid="quality-option"]');
    await expect(qualityOptions).toHaveCount(3);
    await expect(qualityOptions.nth(0)).toContainText('360p');
    await expect(qualityOptions.nth(1)).toContainText('720p');
    await expect(qualityOptions.nth(2)).toContainText('1080p');

    // Verify product overlay appears within 2 seconds
    const productOverlay = page.locator('[data-testid="product-overlay"]');
    await expect(productOverlay).toBeVisible({ timeout: 2000 });
    await expect(productOverlay).toContainText(TEST_PRODUCT.product_id);

    console.log('✓ Step 5: Player UI and product overlay verified');

    /**
     * STEP 6: Send and Moderate Chat Messages
     */
    // Buyer sends chat message
    const sendChatStart = Date.now();

    const sendChatResponse = await request.post(
      `http://localhost:8080/api/v1/streams/${streamId}/chat`,
      {
        headers: {
          'Authorization': `Bearer ${buyerAuthToken}`,
          'Content-Type': 'application/json',
        },
        data: {
          message_text: 'Does it come with a charger?',
        },
      }
    );

    const sendChatDuration = Date.now() - sendChatStart;
    expect(sendChatResponse.status()).toBe(201);

    const chatMessage = await sendChatResponse.json();
    chatMessageId = chatMessage.id;

    expect(chatMessage.message_text).toBe('Does it come with a charger?');
    expect(chatMessage.user_id).toBe(BUYER.user_id);

    console.log(`✓ Step 6a: Chat message sent (${sendChatDuration}ms) - ID: ${chatMessageId}`);

    // Seller moderates chat (removes message)
    const moderateChatStart = Date.now();

    const moderateChatResponse = await request.delete(
      `http://localhost:8080/api/v1/streams/${streamId}/chat/${chatMessageId}`,
      {
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      }
    );

    const moderateChatDuration = Date.now() - moderateChatStart;
    expect(moderateChatResponse.status()).toBe(204);
    expect(moderateChatDuration).toBeLessThan(1000); // <1s target

    // Verify message is removed
    const chatHistoryResponse = await request.get(
      `http://localhost:8080/api/v1/streams/${streamId}/chat`,
      {
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      }
    );

    expect(chatHistoryResponse.ok()).toBeTruthy();
    const chatHistory = await chatHistoryResponse.json();
    const removedMessage = chatHistory.messages?.find((msg: any) => msg.id === chatMessageId);

    if (removedMessage) {
      expect(removedMessage.moderation_status).toBe('removed');
    }

    console.log(`✓ Step 6b: Chat moderated (${moderateChatDuration}ms) - Message removed`);

    /**
     * STEP 7: End Stream
     */
    const endStreamStart = Date.now();

    const endStreamResponse = await request.post(
      `http://localhost:8080/api/v1/streams/${streamId}/end`,
      {
        headers: {
          'Authorization': `Bearer ${authToken}`,
        },
      }
    );

    const endStreamDuration = Date.now() - endStreamStart;
    expect(endStreamResponse.ok()).toBeTruthy();

    const endData = await endStreamResponse.json();

    // Validate stream end response
    expect(endData.stream_id).toBe(streamId);
    expect(endData.status).toBe('ended');
    expect(endData.end_time).toBeTruthy();
    expect(endData.duration_seconds).toBeGreaterThan(0);
    expect(endData.peak_viewer_count).toBeGreaterThanOrEqual(0);

    console.log(`✓ Step 7: Stream ended (${endStreamDuration}ms) - Duration: ${endData.duration_seconds}s`);

    /**
     * STEP 8: Verify Recording Availability
     */
    // Wait for recording to be available (up to 30 seconds)
    let recordingAvailable = false;
    let recordingData: any = null;
    const recordingCheckStart = Date.now();

    for (let attempt = 0; attempt < 6; attempt++) {
      await page.waitForTimeout(5000); // Wait 5 seconds between checks

      const recordingResponse = await request.get(
        `http://localhost:8080/api/v1/streams/${streamId}`,
        {
          headers: {
            'Authorization': `Bearer ${authToken}`,
          },
        }
      );

      if (recordingResponse.ok()) {
        const data = await recordingResponse.json();
        if (data.recording_url) {
          recordingAvailable = true;
          recordingData = data;
          break;
        }
      }
    }

    const recordingCheckDuration = Date.now() - recordingCheckStart;

    expect(recordingAvailable).toBe(true);
    expect(recordingCheckDuration).toBeLessThan(30000); // <30s target
    expect(recordingData.recording_url).toMatch(/^https?:\/\//);
    expect(recordingData.recording_url).toContain('.m3u8');

    // Verify recording expiry date (30 days from end_time)
    expect(recordingData.recording_expiry_date).toBeTruthy();
    expect(
      isApproximately30DaysLater(
        recordingData.recording_expiry_date,
        recordingData.end_time || endData.end_time
      )
    ).toBe(true);

    console.log(`✓ Step 8: Recording available (${recordingCheckDuration}ms) - URL: ${recordingData.recording_url}`);
    console.log(`✓ Recording expiry: ${recordingData.recording_expiry_date} (30 days)`);

    console.log('\n✅ All steps completed successfully!');
  });

  test('should handle stream creation validation errors', async ({ request }) => {
    // Test: Missing required fields
    const invalidStreamResponse = await request.post('http://localhost:8080/api/v1/streams', {
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
      data: {
        title: '', // Empty title
        stream_type: 'store',
      },
    });

    expect(invalidStreamResponse.status()).toBe(400);
    const errorData = await invalidStreamResponse.json();
    expect(errorData.error).toBeTruthy();

    console.log('✓ Validation: Empty title rejected');

    // Test: Invalid stream type
    const invalidTypeResponse = await request.post('http://localhost:8080/api/v1/streams', {
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
      data: {
        title: 'Test Stream',
        stream_type: 'invalid_type',
      },
    });

    expect(invalidTypeResponse.status()).toBe(400);

    console.log('✓ Validation: Invalid stream type rejected');
  });

  test('should prevent unauthorized stream operations', async ({ request }) => {
    // Create a stream as authorized user
    const createResponse = await request.post('http://localhost:8080/api/v1/streams', {
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
      data: {
        ...TEST_STREAM,
        scheduled_start_time: getFutureTimestamp(30),
      },
    });

    expect(createResponse.status()).toBe(201);
    const streamData = await createResponse.json();
    streamId = streamData.id;

    // Try to start stream with buyer token (unauthorized)
    const unauthorizedStartResponse = await request.post(
      `http://localhost:8080/api/v1/streams/${streamId}/start`,
      {
        headers: {
          'Authorization': `Bearer ${buyerAuthToken}`,
          'Content-Type': 'application/json',
        },
        data: {
          webrtc_offer: MOCK_WEBRTC_OFFER,
        },
      }
    );

    expect(unauthorizedStartResponse.status()).toBe(403);

    console.log('✓ Authorization: Unauthorized start prevented');

    // Try to moderate chat with buyer token (unauthorized)
    const unauthorizedModerateResponse = await request.delete(
      `http://localhost:8080/api/v1/streams/${streamId}/chat/fake-message-id`,
      {
        headers: {
          'Authorization': `Bearer ${buyerAuthToken}`,
        },
      }
    );

    expect(unauthorizedModerateResponse.status()).toBe(403);

    console.log('✓ Authorization: Unauthorized moderation prevented');
  });

  test('should handle product featuring errors gracefully', async ({ request }) => {
    // Create and start a stream
    const createResponse = await request.post('http://localhost:8080/api/v1/streams', {
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
      data: {
        ...TEST_STREAM,
        scheduled_start_time: getFutureTimestamp(30),
      },
    });

    expect(createResponse.status()).toBe(201);
    const streamData = await createResponse.json();
    streamId = streamData.id;

    // Try to feature product on non-live stream
    const featureBeforeStartResponse = await request.post(
      `http://localhost:8080/api/v1/streams/${streamId}/products`,
      {
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json',
        },
        data: TEST_PRODUCT,
      }
    );

    expect(featureBeforeStartResponse.status()).toBe(400);
    const errorData = await featureBeforeStartResponse.json();
    expect(errorData.error).toContain('not live');

    console.log('✓ Product featuring: Cannot feature on non-live stream');

    // Try to feature non-existent product
    await request.post(`http://localhost:8080/api/v1/streams/${streamId}/start`, {
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
      data: {
        webrtc_offer: MOCK_WEBRTC_OFFER,
      },
    });

    const invalidProductResponse = await request.post(
      `http://localhost:8080/api/v1/streams/${streamId}/products`,
      {
        headers: {
          'Authorization': `Bearer ${authToken}`,
          'Content-Type': 'application/json',
        },
        data: {
          product_id: 'non-existent-product-id',
          display_position: 'overlay',
        },
      }
    );

    expect(invalidProductResponse.status()).toBe(404);

    console.log('✓ Product featuring: Non-existent product rejected');
  });
});