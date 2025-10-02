/**
 * E2E Test: Video Creator Live Stream Workflow
 *
 * Comprehensive end-to-end testing for content creator live streaming workflow:
 * - Authentication and KYC verification
 * - Instant stream creation and initialization
 * - WebRTC peer connection establishment
 * - Real-time viewer reactions and interactions
 * - Analytics tracking and validation
 * - Stream termination and recording management
 *
 * Test Scenario: Content Creator Instant Video Stream
 * 1. Login as verified content creator
 * 2. Create instant stream (no scheduled_time)
 * 3. Start stream with WebRTC SDP offer
 * 4. Simulate viewer reactions (heart, fire, clap)
 * 5. Verify reaction animations and real-time updates
 * 6. View analytics (peak viewers, reactions, duration)
 * 7. End stream and verify recording expiry (30 days)
 *
 * Reference: /specs/029-implement-live-on/quickstart.md
 */

import { test, expect, Page } from '@playwright/test';

// Test constants
const TEST_CREATOR = {
  email: 'creator@tchat.com',
  password: 'Test123!@#',
  userId: '550e8400-e29b-41d4-a716-446655440000',
  emailVerified: true,
  phoneVerified: true,
  kycTier: 0, // Content creators need email/phone verification (KYC Tier 0+)
};

const TEST_STREAM = {
  title: 'Live Gaming Session - Valorant Ranked',
  description: 'Join me for ranked gameplay and tips!',
  context: 'video' as const,
  category: 'gaming',
  tags: ['gaming', 'valorant', 'fps', 'competitive'],
  thumbnailUrl: 'https://cdn.tchat.com/thumbnails/gaming-session.jpg',
  qualityLayers: ['360p', '720p', '1080p'],
};

const VIEWER_REACTIONS = [
  { type: 'heart', emoji: '❤️', count: 15 },
  { type: 'fire', emoji: '🔥', count: 8 },
  { type: 'clap', emoji: '👏', count: 12 },
];

// Helper functions
async function mockWebRTC(page: Page) {
  await page.addInitScript(() => {
    // Mock MediaDevices.getUserMedia for video/audio capture
    Object.defineProperty(navigator.mediaDevices, 'getUserMedia', {
      value: async (constraints: MediaStreamConstraints) => {
        const mockStream = new MediaStream();

        // Add video track if requested
        if (constraints.video) {
          const mockVideoTrack = {
            kind: 'video',
            enabled: true,
            id: 'mock-video-track',
            label: 'Mock Video Track',
            muted: false,
            readyState: 'live',
            stop: () => {},
            addEventListener: () => {},
            removeEventListener: () => {},
            getSettings: () => ({ width: 1920, height: 1080, frameRate: 30 }),
          } as unknown as MediaStreamTrack;
          mockStream.addTrack(mockVideoTrack);
        }

        // Add audio track if requested
        if (constraints.audio) {
          const mockAudioTrack = {
            kind: 'audio',
            enabled: true,
            id: 'mock-audio-track',
            label: 'Mock Audio Track',
            muted: false,
            readyState: 'live',
            stop: () => {},
            addEventListener: () => {},
            removeEventListener: () => {},
          } as unknown as MediaStreamTrack;
          mockStream.addTrack(mockAudioTrack);
        }

        return mockStream;
      },
      writable: true,
    });

    // Mock RTCPeerConnection for WebRTC streaming
    (window as any).RTCPeerConnection = class MockRTCPeerConnection {
      localDescription: RTCSessionDescription | null = null;
      remoteDescription: RTCSessionDescription | null = null;
      connectionState: RTCPeerConnectionState = 'new';
      iceConnectionState: RTCIceConnectionState = 'new';
      signalingState: RTCSignalingState = 'stable';
      onicecandidate: ((event: RTCPeerConnectionIceEvent) => void) | null = null;
      onconnectionstatechange: ((event: Event) => void) | null = null;
      ontrack: ((event: RTCTrackEvent) => void) | null = null;

      async createOffer(): Promise<RTCSessionDescriptionInit> {
        return {
          type: 'offer',
          sdp: 'v=0\r\no=- 123456789 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0 1\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nm=audio 9 UDP/TLS/RTP/SAVPF 111',
        };
      }

      async createAnswer(): Promise<RTCSessionDescriptionInit> {
        return {
          type: 'answer',
          sdp: 'v=0\r\no=- 987654321 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0 1\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nm=audio 9 UDP/TLS/RTP/SAVPF 111',
        };
      }

      async setLocalDescription(description: RTCSessionDescriptionInit): Promise<void> {
        this.localDescription = description as RTCSessionDescription;
        this.signalingState = 'have-local-offer';
      }

      async setRemoteDescription(description: RTCSessionDescriptionInit): Promise<void> {
        this.remoteDescription = description as RTCSessionDescription;
        this.signalingState = 'stable';
        this.connectionState = 'connected';
        this.iceConnectionState = 'connected';

        // Simulate connection established
        if (this.onconnectionstatechange) {
          this.onconnectionstatechange(new Event('connectionstatechange'));
        }
      }

      addTrack(track: MediaStreamTrack, stream: MediaStream): RTCRtpSender {
        return {} as RTCRtpSender;
      }

      addTransceiver(trackOrKind: MediaStreamTrack | string, init?: RTCRtpTransceiverInit): RTCRtpTransceiver {
        return {} as RTCRtpTransceiver;
      }

      close(): void {
        this.connectionState = 'closed';
      }

      addEventListener(): void {}
      removeEventListener(): void {}
    };
  });
}

async function loginAsCreator(page: Page) {
  await page.goto('/login');

  // Fill login credentials
  await page.fill('[data-testid="email-input"]', TEST_CREATOR.email);
  await page.fill('[data-testid="password-input"]', TEST_CREATOR.password);
  await page.click('[data-testid="login-button"]');

  // Wait for successful login and dashboard navigation
  await expect(page).toHaveURL(/\/(dashboard|home)/);
  await page.waitForSelector('[data-testid="user-profile"]');
}

async function verifyCreatorEligibility(page: Page) {
  // Mock user verification status API
  await page.route('**/api/v1/users/me', (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        id: TEST_CREATOR.userId,
        email: TEST_CREATOR.email,
        email_verified: TEST_CREATOR.emailVerified,
        phone_verified: TEST_CREATOR.phoneVerified,
        kyc_tier: TEST_CREATOR.kycTier,
        can_stream: true,
        stream_contexts: ['video'],
      }),
    });
  });

  // Verify user status
  const response = await page.request.get('http://localhost:8080/api/v1/users/me', {
    headers: {
      'Authorization': `Bearer mock-jwt-token-${TEST_CREATOR.userId}`,
    },
  });

  expect(response.ok()).toBeTruthy();
  const userData = await response.json();
  expect(userData.email_verified).toBe(true);
  expect(userData.can_stream).toBe(true);
}

async function createInstantStream(page: Page): Promise<string> {
  let streamId = '';

  // Mock stream creation API
  await page.route('**/api/v1/streams', (route) => {
    if (route.request().method() === 'POST') {
      streamId = `stream-${Date.now()}`;
      const requestBody = JSON.parse(route.request().postData() || '{}');

      route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          id: streamId,
          user_id: TEST_CREATOR.userId,
          title: requestBody.title,
          description: requestBody.description,
          context: requestBody.context,
          category: requestBody.category,
          tags: requestBody.tags,
          thumbnail_url: requestBody.thumbnail_url,
          scheduled_time: null, // Instant stream (no scheduled time)
          status: 'scheduled',
          created_at: new Date().toISOString(),
          visibility: 'public',
          max_viewers: 10000,
          chat_enabled: true,
          reactions_enabled: true,
          products_enabled: false,
        }),
      });
    }
  });

  // Navigate to stream creation page
  await page.goto('/streams/create');
  await expect(page.locator('[data-testid="stream-create-form"]')).toBeVisible();

  // Fill stream details
  await page.fill('[data-testid="stream-title-input"]', TEST_STREAM.title);
  await page.fill('[data-testid="stream-description-input"]', TEST_STREAM.description);
  await page.selectOption('[data-testid="stream-context-select"]', TEST_STREAM.context);
  await page.selectOption('[data-testid="stream-category-select"]', TEST_STREAM.category);

  // Add tags
  for (const tag of TEST_STREAM.tags) {
    await page.fill('[data-testid="stream-tags-input"]', tag);
    await page.keyboard.press('Enter');
  }

  // Select instant stream (no scheduling)
  await page.click('[data-testid="instant-stream-checkbox"]');

  // Submit stream creation
  await page.click('[data-testid="create-stream-button"]');

  // Wait for success notification
  await expect(page.locator('.toast')).toContainText('Stream created successfully');

  // Verify navigation to stream setup page
  await page.waitForURL(/\/streams\/.*\/setup/);

  // Extract stream ID from URL
  const url = page.url();
  streamId = url.match(/\/streams\/([^\/]+)\/setup/)?.[1] || streamId;

  return streamId;
}

async function startStream(page: Page, streamId: string) {
  // Mock start stream API
  await page.route(`**/api/v1/streams/${streamId}/start`, (route) => {
    if (route.request().method() === 'POST') {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: streamId,
          status: 'live',
          actual_start_time: new Date().toISOString(),
          rtc_session_id: `rtc-${streamId}`,
          ice_servers: [
            { urls: 'stun:stun.tchat.com:3478' },
            { urls: 'turn:turn.tchat.com:3478', username: 'tchat', credential: 'secret' },
          ],
          quality_layers: TEST_STREAM.qualityLayers,
        }),
      });
    }
  });

  // Mock WebRTC offer/answer exchange
  await page.route('**/api/v1/streams/*/webrtc/offer', (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        type: 'answer',
        sdp: 'v=0\r\no=- 987654321 2 IN IP4 10.0.0.1\r\ns=WebRTC Stream\r\nt=0 0\r\na=group:BUNDLE 0 1\r\nm=video 9 UDP/TLS/RTP/SAVPF 96\r\nm=audio 9 UDP/TLS/RTP/SAVPF 111',
      }),
    });
  });

  // Start streaming
  await page.click('[data-testid="start-stream-button"]');

  // Wait for stream to go live
  await expect(page.locator('[data-testid="stream-status"]')).toContainText('Live');
  await expect(page.locator('[data-testid="live-indicator"]')).toBeVisible();

  // Verify video player is active
  await expect(page.locator('[data-testid="stream-video-player"]')).toBeVisible();
  await expect(page.locator('[data-testid="quality-selector"]')).toBeVisible();
}

async function simulateViewerReactions(page: Page, streamId: string) {
  // Mock reaction API
  await page.route(`**/api/v1/streams/${streamId}/react`, (route) => {
    if (route.request().method() === 'POST') {
      const requestBody = JSON.parse(route.request().postData() || '{}');
      route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          reaction_type: requestBody.reaction_type,
          timestamp: new Date().toISOString(),
          animation_id: `anim-${Date.now()}`,
        }),
      });
    }
  });

  // Simulate reactions from multiple viewers
  for (const reaction of VIEWER_REACTIONS) {
    for (let i = 0; i < reaction.count; i++) {
      // Send reaction API call
      await page.request.post(`http://localhost:8080/api/v1/streams/${streamId}/react`, {
        headers: {
          'Authorization': `Bearer viewer-token-${i}`,
          'Content-Type': 'application/json',
        },
        data: {
          reaction_type: reaction.type,
        },
      });

      // Wait briefly to simulate realistic timing
      await page.waitForTimeout(100);
    }
  }
}

async function verifyReactionAnimations(page: Page) {
  // Wait for reaction overlay container to be visible
  await expect(page.locator('[data-testid="reaction-overlay"]')).toBeVisible();

  // Verify reaction animations appear
  for (const reaction of VIEWER_REACTIONS) {
    const reactionElements = page.locator(`[data-testid="reaction-${reaction.type}"]`);
    const count = await reactionElements.count();

    // At least some reactions should be visible (they animate and disappear)
    expect(count).toBeGreaterThan(0);
  }

  // Verify reaction animations have proper styling (floating effect)
  const firstReaction = page.locator('[data-testid^="reaction-"]').first();
  await expect(firstReaction).toHaveCSS('animation-name', /float|rise|bounce/);
}

async function viewAnalytics(page: Page, streamId: string): Promise<any> {
  const mockAnalytics = {
    stream_id: streamId,
    peak_viewers: 245,
    total_views: 312,
    avg_watch_time: 423, // seconds (7 minutes)
    reactions_count: VIEWER_REACTIONS.reduce((sum, r) => sum + r.count, 0),
    reactions_breakdown: VIEWER_REACTIONS.reduce((acc, r) => {
      acc[r.type] = r.count;
      return acc;
    }, {} as Record<string, number>),
    chat_messages: 89,
    duration: 1847, // seconds (30 minutes 47 seconds)
    revenue: 0, // Video streams don't have product sales
  };

  // Mock analytics API
  await page.route(`**/api/v1/streams/${streamId}/analytics`, (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(mockAnalytics),
    });
  });

  // Open analytics panel
  await page.click('[data-testid="analytics-button"]');
  await expect(page.locator('[data-testid="analytics-panel"]')).toBeVisible();

  // Wait for analytics to load
  await page.waitForSelector('[data-testid="analytics-peak-viewers"]');

  return mockAnalytics;
}

async function verifyAnalyticsMetrics(page: Page, expectedAnalytics: any) {
  // Verify peak viewers
  const peakViewersElement = page.locator('[data-testid="analytics-peak-viewers"]');
  await expect(peakViewersElement).toContainText(expectedAnalytics.peak_viewers.toString());

  // Verify total views
  const totalViewsElement = page.locator('[data-testid="analytics-total-views"]');
  await expect(totalViewsElement).toContainText(expectedAnalytics.total_views.toString());

  // Verify average watch time
  const avgWatchTimeElement = page.locator('[data-testid="analytics-avg-watch-time"]');
  await expect(avgWatchTimeElement).toBeVisible();

  // Verify reactions count
  const reactionsCountElement = page.locator('[data-testid="analytics-reactions-count"]');
  await expect(reactionsCountElement).toContainText(expectedAnalytics.reactions_count.toString());

  // Verify reactions breakdown
  for (const [type, count] of Object.entries(expectedAnalytics.reactions_breakdown)) {
    const reactionBreakdown = page.locator(`[data-testid="analytics-reaction-${type}"]`);
    await expect(reactionBreakdown).toContainText(count.toString());
  }

  // Verify duration is calculated correctly
  const durationElement = page.locator('[data-testid="analytics-duration"]');
  await expect(durationElement).toBeVisible();
  const durationText = await durationElement.textContent();
  expect(durationText).toMatch(/\d+:\d{2}/); // Format: MM:SS or HH:MM:SS
}

async function endStream(page: Page, streamId: string) {
  const endTime = new Date();
  const expiryDate = new Date(endTime.getTime() + 30 * 24 * 60 * 60 * 1000); // +30 days

  // Mock end stream API
  await page.route(`**/api/v1/streams/${streamId}/end`, (route) => {
    if (route.request().method() === 'POST') {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: streamId,
          status: 'ended',
          actual_end_time: endTime.toISOString(),
          recording_url: `https://cdn.tchat.com/recordings/${streamId}.mp4`,
          recording_duration: 1847, // seconds
          recording_expires_at: expiryDate.toISOString(),
        }),
      });
    }
  });

  // End the stream
  await page.click('[data-testid="end-stream-button"]');

  // Confirm end stream dialog
  await page.click('[data-testid="confirm-end-stream-button"]');

  // Wait for stream to end
  await expect(page.locator('[data-testid="stream-status"]')).toContainText('Ended');
  await expect(page.locator('[data-testid="live-indicator"]')).not.toBeVisible();

  // Verify success notification
  await expect(page.locator('.toast')).toContainText('Stream ended successfully');
}

async function verifyRecording(page: Page, streamId: string) {
  const expiryDate = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000);

  // Mock recording API
  await page.route(`**/api/v1/streams/${streamId}/recording`, (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        stream_id: streamId,
        recording_url: `https://cdn.tchat.com/recordings/${streamId}.mp4`,
        thumbnail_url: `https://cdn.tchat.com/thumbnails/${streamId}.jpg`,
        duration: 1847,
        file_size: 245678901, // bytes (~245 MB)
        quality: '1080p',
        format: 'mp4',
        created_at: new Date().toISOString(),
        expires_at: expiryDate.toISOString(),
        is_public: true,
        download_count: 0,
        view_count: 0,
      }),
    });
  });

  // Navigate to recording page
  await page.click('[data-testid="view-recording-button"]');
  await page.waitForURL(/\/streams\/.*\/recording/);

  // Verify recording is available
  await expect(page.locator('[data-testid="recording-player"]')).toBeVisible();
  await expect(page.locator('[data-testid="recording-download-button"]')).toBeVisible();

  // Verify expiry date is displayed (30 days from now)
  const expiryElement = page.locator('[data-testid="recording-expiry-date"]');
  await expect(expiryElement).toBeVisible();
  const expiryText = await expiryElement.textContent();

  // Check that expiry is approximately 30 days from now (allow 1 day margin)
  const displayedExpiry = new Date(expiryText || '');
  const daysDifference = Math.abs((displayedExpiry.getTime() - expiryDate.getTime()) / (1000 * 60 * 60 * 24));
  expect(daysDifference).toBeLessThan(1);
}

// Main test suite
test.describe('Video Creator Stream Workflow', () => {
  test.beforeEach(async ({ page }) => {
    // Setup WebRTC mocking before each test
    await mockWebRTC(page);
  });

  test('Complete video creator instant stream flow', async ({ page }) => {
    // Step 1: Login as verified content creator
    await loginAsCreator(page);

    // Step 2: Verify creator eligibility (email/phone verified)
    await verifyCreatorEligibility(page);

    // Step 3: Create instant stream (no scheduled_time)
    const streamId = await createInstantStream(page);
    expect(streamId).toBeTruthy();

    // Step 4: Start stream with WebRTC setup
    await startStream(page, streamId);

    // Step 5: Simulate viewer reactions (heart, fire, clap)
    await simulateViewerReactions(page, streamId);

    // Step 6: Verify reaction animations appear in UI
    await verifyReactionAnimations(page);

    // Step 7: View analytics
    const analytics = await viewAnalytics(page, streamId);

    // Step 8: Verify analytics accuracy
    await verifyAnalyticsMetrics(page, analytics);
    expect(analytics.peak_viewers).toBeGreaterThanOrEqual(1);
    expect(analytics.reactions_count).toBe(VIEWER_REACTIONS.reduce((sum, r) => sum + r.count, 0));

    // Step 9: End stream
    await endStream(page, streamId);

    // Step 10: Verify recording with 30-day expiry
    await verifyRecording(page, streamId);
  });

  test('Stream status transitions correctly', async ({ page }) => {
    await loginAsCreator(page);
    await verifyCreatorEligibility(page);

    const streamId = await createInstantStream(page);

    // Verify initial status: scheduled
    await expect(page.locator('[data-testid="stream-status"]')).toContainText('Scheduled');

    // Start stream → status: live
    await startStream(page, streamId);
    await expect(page.locator('[data-testid="stream-status"]')).toContainText('Live');

    // End stream → status: ended
    await endStream(page, streamId);
    await expect(page.locator('[data-testid="stream-status"]')).toContainText('Ended');
  });

  test('Reactions appear within 500ms of sending', async ({ page }) => {
    await loginAsCreator(page);
    const streamId = await createInstantStream(page);
    await startStream(page, streamId);

    // Mock single reaction API
    await page.route(`**/api/v1/streams/${streamId}/react`, (route) => {
      route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          reaction_type: 'heart',
          timestamp: new Date().toISOString(),
          animation_id: `anim-test`,
        }),
      });
    });

    // Send reaction and measure time
    const startTime = Date.now();

    await page.request.post(`http://localhost:8080/api/v1/streams/${streamId}/react`, {
      headers: {
        'Authorization': 'Bearer test-viewer-token',
        'Content-Type': 'application/json',
      },
      data: { reaction_type: 'heart' },
    });

    // Wait for reaction to appear
    await page.waitForSelector('[data-testid="reaction-heart"]', { timeout: 1000 });
    const endTime = Date.now();

    // Verify reaction appeared within 500ms
    const responseTime = endTime - startTime;
    expect(responseTime).toBeLessThan(500);
  });

  test('Analytics calculation is accurate', async ({ page }) => {
    await loginAsCreator(page);
    const streamId = await createInstantStream(page);
    await startStream(page, streamId);

    // Simulate specific viewer actions
    const expectedReactions = 25; // 10 hearts + 8 fires + 7 claps
    const specificReactions = [
      { type: 'heart', count: 10 },
      { type: 'fire', count: 8 },
      { type: 'clap', count: 7 },
    ];

    for (const reaction of specificReactions) {
      for (let i = 0; i < reaction.count; i++) {
        await page.request.post(`http://localhost:8080/api/v1/streams/${streamId}/react`, {
          headers: { 'Authorization': `Bearer viewer-${i}`, 'Content-Type': 'application/json' },
          data: { reaction_type: reaction.type },
        });
      }
    }

    // View analytics
    await page.route(`**/api/v1/streams/${streamId}/analytics`, (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          stream_id: streamId,
          peak_viewers: 15,
          total_views: 23,
          reactions_count: expectedReactions,
          reactions_breakdown: {
            heart: 10,
            fire: 8,
            clap: 7,
          },
        }),
      });
    });

    await page.click('[data-testid="analytics-button"]');

    // Verify exact reaction counts
    const reactionsCount = page.locator('[data-testid="analytics-reactions-count"]');
    await expect(reactionsCount).toContainText(expectedReactions.toString());

    // Verify breakdown accuracy
    await expect(page.locator('[data-testid="analytics-reaction-heart"]')).toContainText('10');
    await expect(page.locator('[data-testid="analytics-reaction-fire"]')).toContainText('8');
    await expect(page.locator('[data-testid="analytics-reaction-clap"]')).toContainText('7');
  });

  test('Recording expires after 30 days', async ({ page }) => {
    await loginAsCreator(page);
    const streamId = await createInstantStream(page);
    await startStream(page, streamId);
    await endStream(page, streamId);

    const now = new Date();
    const expiry = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000);

    // Mock recording with exact 30-day expiry
    await page.route(`**/api/v1/streams/${streamId}/recording`, (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          stream_id: streamId,
          recording_url: `https://cdn.tchat.com/recordings/${streamId}.mp4`,
          created_at: now.toISOString(),
          expires_at: expiry.toISOString(),
        }),
      });
    });

    await page.click('[data-testid="view-recording-button"]');

    // Verify expiry date is exactly 30 days
    const expiryElement = page.locator('[data-testid="recording-expiry-date"]');
    const expiryText = await expiryElement.textContent();
    const displayedExpiry = new Date(expiryText || '');

    const daysDifference = Math.round((displayedExpiry.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
    expect(daysDifference).toBe(30);
  });

  test('Quality selector shows all available layers', async ({ page }) => {
    await loginAsCreator(page);
    const streamId = await createInstantStream(page);
    await startStream(page, streamId);

    // Click quality selector
    await page.click('[data-testid="quality-selector"]');

    // Verify all quality layers are available
    for (const quality of TEST_STREAM.qualityLayers) {
      const qualityOption = page.locator(`[data-testid="quality-option-${quality}"]`);
      await expect(qualityOption).toBeVisible();
    }

    // Verify Auto quality option is present
    await expect(page.locator('[data-testid="quality-option-auto"]')).toBeVisible();
  });

  test('WebRTC connection establishes successfully', async ({ page }) => {
    await loginAsCreator(page);
    const streamId = await createInstantStream(page);

    // Monitor WebRTC connection state
    let connectionState = 'new';
    page.on('console', (msg) => {
      if (msg.text().includes('RTCPeerConnection state:')) {
        connectionState = msg.text().split(':')[1].trim();
      }
    });

    await startStream(page, streamId);

    // Wait for connection to establish
    await page.waitForTimeout(1000);

    // Verify connection state is connected
    expect(['connected', 'completed']).toContain(connectionState);

    // Verify stream is live
    await expect(page.locator('[data-testid="live-indicator"]')).toBeVisible();
  });

  test('Test data cleanup after stream end', async ({ page }) => {
    await loginAsCreator(page);
    const streamId = await createInstantStream(page);
    await startStream(page, streamId);
    await endStream(page, streamId);

    // Verify cleanup notification
    await expect(page.locator('.toast')).toContainText('Stream ended successfully');

    // Verify WebRTC connection is closed
    const videoPlayer = page.locator('[data-testid="stream-video-player"]');
    await expect(videoPlayer).not.toBeVisible();

    // Verify live indicator is removed
    const liveIndicator = page.locator('[data-testid="live-indicator"]');
    await expect(liveIndicator).not.toBeVisible();
  });
});

test.describe('Video Creator Stream Error Handling', () => {
  test('Handles stream creation failure gracefully', async ({ page }) => {
    await mockWebRTC(page);
    await loginAsCreator(page);

    // Mock stream creation failure
    await page.route('**/api/v1/streams', (route) => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({
          error: 'Internal server error',
          message: 'Failed to create stream',
        }),
      });
    });

    await page.goto('/streams/create');
    await page.fill('[data-testid="stream-title-input"]', TEST_STREAM.title);
    await page.click('[data-testid="create-stream-button"]');

    // Verify error message
    await expect(page.locator('.toast-error')).toContainText('Failed to create stream');
  });

  test('Handles WebRTC connection failure', async ({ page }) => {
    await mockWebRTC(page);
    await loginAsCreator(page);
    const streamId = await createInstantStream(page);

    // Mock WebRTC offer failure
    await page.route('**/api/v1/streams/*/webrtc/offer', (route) => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({
          error: 'WebRTC connection failed',
        }),
      });
    });

    await page.click('[data-testid="start-stream-button"]');

    // Verify error notification
    await expect(page.locator('.toast-error')).toContainText(/WebRTC|connection/i);
  });

  test('Handles reaction API failure without breaking stream', async ({ page }) => {
    await mockWebRTC(page);
    await loginAsCreator(page);
    const streamId = await createInstantStream(page);
    await startStream(page, streamId);

    // Mock reaction API failure
    await page.route(`**/api/v1/streams/${streamId}/react`, (route) => {
      route.fulfill({ status: 500 });
    });

    // Attempt to send reaction
    await page.request.post(`http://localhost:8080/api/v1/streams/${streamId}/react`, {
      headers: { 'Authorization': 'Bearer test-token', 'Content-Type': 'application/json' },
      data: { reaction_type: 'heart' },
    });

    // Verify stream continues to work
    await expect(page.locator('[data-testid="stream-status"]')).toContainText('Live');
    await expect(page.locator('[data-testid="live-indicator"]')).toBeVisible();
  });
});