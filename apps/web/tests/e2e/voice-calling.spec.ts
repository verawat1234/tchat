import { test, expect } from '@playwright/test';

/**
 * E2E Tests for Voice Calling Feature
 *
 * These tests validate the complete voice calling workflow from the web interface,
 * including WebRTC integration, real-time signaling, and user interactions.
 */

test.describe('Voice Calling E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the chat application
    await page.goto('/');

    // Mock WebRTC getUserMedia for testing
    await page.addInitScript(() => {
      // Mock MediaDevices.getUserMedia
      Object.defineProperty(navigator.mediaDevices, 'getUserMedia', {
        value: async (constraints: MediaStreamConstraints) => {
          // Create mock audio stream
          const mockStream = new MediaStream();
          const mockTrack = {
            kind: 'audio',
            enabled: true,
            id: 'mock-audio-track',
            label: 'Mock Audio Track',
            muted: false,
            readyState: 'live',
            stop: () => {},
            addEventListener: () => {},
            removeEventListener: () => {},
          } as MediaStreamTrack;
          mockStream.addTrack(mockTrack);
          return mockStream;
        },
        writable: true,
      });

      // Mock RTCPeerConnection
      (window as any).RTCPeerConnection = class MockRTCPeerConnection {
        localDescription = null;
        remoteDescription = null;
        connectionState = 'new';
        iceConnectionState = 'new';
        signalingState = 'stable';
        onicecandidate = null;
        onconnectionstatechange = null;

        async createOffer() {
          return {
            type: 'offer',
            sdp: 'mock-offer-sdp'
          };
        }

        async createAnswer() {
          return {
            type: 'answer',
            sdp: 'mock-answer-sdp'
          };
        }

        async setLocalDescription(description: any) {
          this.localDescription = description;
        }

        async setRemoteDescription(description: any) {
          this.remoteDescription = description;
        }

        addStream() {}
        addTrack() {}
        close() {}
        addEventListener() {}
        removeEventListener() {}
      };
    });
  });

  test('should initiate a voice call successfully', async ({ page }) => {
    // Arrange: Login and navigate to chat
    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    // Wait for chat interface to load
    await expect(page.locator('[data-testid="chat-container"]')).toBeVisible();

    // Select a contact to call
    await page.click('[data-testid="contact-item"]');
    await expect(page.locator('[data-testid="chat-header"]')).toBeVisible();

    // Act: Initiate voice call
    await page.click('[data-testid="voice-call-button"]');

    // Assert: Call UI should appear
    await expect(page.locator('[data-testid="call-interface"]')).toBeVisible();
    await expect(page.locator('[data-testid="call-status"]')).toContainText('Calling...');
    await expect(page.locator('[data-testid="call-type-indicator"]')).toContainText('Voice Call');

    // Verify audio controls are visible
    await expect(page.locator('[data-testid="mute-button"]')).toBeVisible();
    await expect(page.locator('[data-testid="end-call-button"]')).toBeVisible();

    // Verify video controls are not visible for voice call
    await expect(page.locator('[data-testid="video-toggle-button"]')).not.toBeVisible();
    await expect(page.locator('[data-testid="camera-switch-button"]')).not.toBeVisible();
  });

  test('should handle call acceptance flow', async ({ page, context }) => {
    // Create two browser contexts to simulate caller and callee
    const calleeContext = await context.browser()?.newContext();
    const calleePage = await calleeContext?.newPage();

    if (!calleePage) {
      throw new Error('Failed to create callee page');
    }

    // Setup callee
    await calleePage.goto('/');
    await calleePage.fill('[data-testid="email-input"]', 'callee@example.com');
    await calleePage.fill('[data-testid="password-input"]', 'password123');
    await calleePage.click('[data-testid="login-button"]');

    // Setup caller
    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    // Wait for both to be ready
    await expect(page.locator('[data-testid="chat-container"]')).toBeVisible();
    await expect(calleePage.locator('[data-testid="chat-container"]')).toBeVisible();

    // Caller initiates call
    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');

    // Callee should receive incoming call notification
    await expect(calleePage.locator('[data-testid="incoming-call-notification"]')).toBeVisible();
    await expect(calleePage.locator('[data-testid="caller-name"]')).toContainText('caller@example.com');
    await expect(calleePage.locator('[data-testid="call-type-incoming"]')).toContainText('Voice Call');

    // Callee accepts the call
    await calleePage.click('[data-testid="accept-call-button"]');

    // Both should now be in active call state
    await expect(page.locator('[data-testid="call-status"]')).toContainText('Connected');
    await expect(calleePage.locator('[data-testid="call-status"]')).toContainText('Connected');

    // Verify call duration timer is running
    await expect(page.locator('[data-testid="call-timer"]')).toBeVisible();
    await expect(calleePage.locator('[data-testid="call-timer"]')).toBeVisible();

    // Clean up
    await calleeContext?.close();
  });

  test('should handle call decline flow', async ({ page, context }) => {
    const calleeContext = await context.browser()?.newContext();
    const calleePage = await calleeContext?.newPage();

    if (!calleePage) {
      throw new Error('Failed to create callee page');
    }

    // Setup both users
    await calleePage.goto('/');
    await calleePage.fill('[data-testid="email-input"]', 'callee@example.com');
    await calleePage.fill('[data-testid="password-input"]', 'password123');
    await calleePage.click('[data-testid="login-button"]');

    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    // Caller initiates call
    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');

    // Callee declines the call
    await expect(calleePage.locator('[data-testid="incoming-call-notification"]')).toBeVisible();
    await calleePage.click('[data-testid="decline-call-button"]');

    // Caller should see call declined message
    await expect(page.locator('[data-testid="call-status"]')).toContainText('Call Declined');
    await expect(page.locator('[data-testid="call-interface"]')).not.toBeVisible({ timeout: 5000 });

    // Callee should return to normal chat
    await expect(calleePage.locator('[data-testid="incoming-call-notification"]')).not.toBeVisible();

    await calleeContext?.close();
  });

  test('should handle audio controls during call', async ({ page, context }) => {
    const calleeContext = await context.browser()?.newContext();
    const calleePage = await calleeContext?.newPage();

    if (!calleePage) {
      throw new Error('Failed to create callee page');
    }

    // Setup and establish call
    await calleePage.goto('/');
    await calleePage.fill('[data-testid="email-input"]', 'callee@example.com');
    await calleePage.fill('[data-testid="password-input"]', 'password123');
    await calleePage.click('[data-testid="login-button"]');

    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');
    await calleePage.click('[data-testid="accept-call-button"]');

    // Wait for call to be established
    await expect(page.locator('[data-testid="call-status"]')).toContainText('Connected');

    // Test mute functionality
    await expect(page.locator('[data-testid="mute-button"]')).not.toHaveClass(/muted/);
    await page.click('[data-testid="mute-button"]');
    await expect(page.locator('[data-testid="mute-button"]')).toHaveClass(/muted/);
    await expect(page.locator('[data-testid="mute-indicator"]')).toBeVisible();

    // Test unmute
    await page.click('[data-testid="mute-button"]');
    await expect(page.locator('[data-testid="mute-button"]')).not.toHaveClass(/muted/);
    await expect(page.locator('[data-testid="mute-indicator"]')).not.toBeVisible();

    // Test speaker toggle (if available)
    if (await page.locator('[data-testid="speaker-button"]').isVisible()) {
      await page.click('[data-testid="speaker-button"]');
      await expect(page.locator('[data-testid="speaker-button"]')).toHaveClass(/active/);
    }

    await calleeContext?.close();
  });

  test('should end call properly', async ({ page, context }) => {
    const calleeContext = await context.browser()?.newContext();
    const calleePage = await calleeContext?.newPage();

    if (!calleePage) {
      throw new Error('Failed to create callee page');
    }

    // Setup and establish call
    await calleePage.goto('/');
    await calleePage.fill('[data-testid="email-input"]', 'callee@example.com');
    await calleePage.fill('[data-testid="password-input"]', 'password123');
    await calleePage.click('[data-testid="login-button"]');

    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');
    await calleePage.click('[data-testid="accept-call-button"]');

    // Wait for call to be established
    await expect(page.locator('[data-testid="call-status"]')).toContainText('Connected');

    // Record call duration before ending
    const callTimer = page.locator('[data-testid="call-timer"]');
    await expect(callTimer).toBeVisible();

    // End call from caller side
    await page.click('[data-testid="end-call-button"]');

    // Both sides should show call ended
    await expect(page.locator('[data-testid="call-interface"]')).not.toBeVisible({ timeout: 5000 });
    await expect(calleePage.locator('[data-testid="call-interface"]')).not.toBeVisible({ timeout: 5000 });

    // Should return to normal chat interface
    await expect(page.locator('[data-testid="chat-container"]')).toBeVisible();
    await expect(calleePage.locator('[data-testid="chat-container"]')).toBeVisible();

    await calleeContext?.close();
  });

  test('should handle call timeout', async ({ page }) => {
    // Setup caller
    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');

    // Verify call is in progress
    await expect(page.locator('[data-testid="call-status"]')).toContainText('Calling...');

    // Mock no answer timeout (wait for timeout or simulate it)
    await page.waitForTimeout(30000); // Wait for call timeout

    // Should show timeout message and return to chat
    await expect(page.locator('[data-testid="call-status"]')).toContainText('No Answer');
    await expect(page.locator('[data-testid="call-interface"]')).not.toBeVisible({ timeout: 5000 });
  });

  test('should handle network issues during call', async ({ page, context }) => {
    const calleeContext = await context.browser()?.newContext();
    const calleePage = await calleeContext?.newPage();

    if (!calleePage) {
      throw new Error('Failed to create callee page');
    }

    // Setup and establish call
    await calleePage.goto('/');
    await calleePage.fill('[data-testid="email-input"]', 'callee@example.com');
    await calleePage.fill('[data-testid="password-input"]', 'password123');
    await calleePage.click('[data-testid="login-button"]');

    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');
    await calleePage.click('[data-testid="accept-call-button"]');

    // Wait for call to be established
    await expect(page.locator('[data-testid="call-status"]')).toContainText('Connected');

    // Simulate network issues by going offline
    await page.context().setOffline(true);

    // Should show connection issues
    await expect(page.locator('[data-testid="call-status"]')).toContainText(/Reconnecting|Connection Lost/, { timeout: 10000 });
    await expect(page.locator('[data-testid="connection-indicator"]')).toBeVisible();

    // Restore connection
    await page.context().setOffline(false);

    // Should attempt to reconnect
    await expect(page.locator('[data-testid="call-status"]')).toContainText(/Connected|Reconnected/, { timeout: 15000 });

    await calleeContext?.close();
  });

  test('should display call history after call ends', async ({ page, context }) => {
    const calleeContext = await context.browser()?.newContext();
    const calleePage = await calleeContext?.newPage();

    if (!calleePage) {
      throw new Error('Failed to create callee page');
    }

    // Setup and complete a call
    await calleePage.goto('/');
    await calleePage.fill('[data-testid="email-input"]', 'callee@example.com');
    await calleePage.fill('[data-testid="password-input"]', 'password123');
    await calleePage.click('[data-testid="login-button"]');

    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');
    await calleePage.click('[data-testid="accept-call-button"]');

    // Wait a bit then end call
    await page.waitForTimeout(2000);
    await page.click('[data-testid="end-call-button"]');

    // Navigate to call history
    await page.click('[data-testid="call-history-tab"]');

    // Verify call appears in history
    await expect(page.locator('[data-testid="call-history-list"]')).toBeVisible();

    const historyItems = page.locator('[data-testid="call-history-item"]');
    await expect(historyItems.first()).toBeVisible();

    // Verify call details
    const firstCall = historyItems.first();
    await expect(firstCall.locator('[data-testid="call-type"]')).toContainText('Voice');
    await expect(firstCall.locator('[data-testid="call-direction"]')).toContainText('Outgoing');
    await expect(firstCall.locator('[data-testid="call-duration"]')).toBeVisible();
    await expect(firstCall.locator('[data-testid="call-timestamp"]')).toBeVisible();

    await calleeContext?.close();
  });

  test('should handle multiple simultaneous call attempts', async ({ page }) => {
    // Setup caller
    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    // Try to initiate multiple calls
    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');

    // Attempt second call while first is active
    await page.click('[data-testid="contact-item-2"]');

    // Should show error or disable call button
    const secondCallButton = page.locator('[data-testid="voice-call-button"]');
    await expect(secondCallButton).toBeDisabled();

    // Or should show error message
    await secondCallButton.click();
    await expect(page.locator('[data-testid="error-message"]')).toContainText(/already in call|call in progress/);
  });

  test('should test accessibility features', async ({ page }) => {
    // Setup
    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');

    // Test keyboard navigation
    await page.keyboard.press('Tab');
    await expect(page.locator('[data-testid="mute-button"]:focus')).toBeVisible();

    await page.keyboard.press('Tab');
    await expect(page.locator('[data-testid="end-call-button"]:focus')).toBeVisible();

    // Test ARIA labels
    await expect(page.locator('[data-testid="mute-button"]')).toHaveAttribute('aria-label', /mute|unmute/i);
    await expect(page.locator('[data-testid="end-call-button"]')).toHaveAttribute('aria-label', /end call/i);

    // Test screen reader announcements
    await expect(page.locator('[data-testid="call-status"]')).toHaveAttribute('aria-live', 'polite');

    // Test high contrast mode compatibility
    await page.emulateMedia({ prefers_color_scheme: 'dark' });
    await expect(page.locator('[data-testid="call-interface"]')).toBeVisible();
  });

  test('should handle poor network conditions', async ({ page, context }) => {
    // Setup slow network conditions
    await page.route('**/*', route => {
      // Add delay to simulate slow network
      setTimeout(() => route.continue(), 1000);
    });

    const calleeContext = await context.browser()?.newContext();
    const calleePage = await calleeContext?.newPage();

    if (!calleePage) {
      throw new Error('Failed to create callee page');
    }

    // Setup and attempt call under poor conditions
    await calleePage.goto('/');
    await calleePage.fill('[data-testid="email-input"]', 'callee@example.com');
    await calleePage.fill('[data-testid="password-input"]', 'password123');
    await calleePage.click('[data-testid="login-button"]');

    await page.fill('[data-testid="email-input"]', 'caller@example.com');
    await page.fill('[data-testid="password-input"]', 'password123');
    await page.click('[data-testid="login-button"]');

    await page.click('[data-testid="contact-item"]');
    await page.click('[data-testid="voice-call-button"]');

    // Should show connection quality indicators
    await expect(page.locator('[data-testid="connection-quality"]')).toBeVisible();
    await expect(page.locator('[data-testid="connection-quality"]')).toContainText(/poor|fair/i);

    await calleeContext?.close();
  });
});