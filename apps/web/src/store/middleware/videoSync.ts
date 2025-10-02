/**
 * Video Sync Middleware
 *
 * Real-time cross-platform synchronization middleware for video playback:
 * - Automatic position sync every 5 seconds during playback
 * - Debounced sync to prevent excessive API calls
 * - Conflict detection and resolution
 * - Network resilience with retry logic
 * - Performance monitoring and metrics
 */

import { Middleware } from '@reduxjs/toolkit';
import type { RootState } from '../index';
import {
  updatePlaybackPosition,
  setPlaybackState,
  updateSyncStatus,
  syncPlayback,
} from '../slices/videoSlice';
import { syncBackendPlaybackPosition } from '../../services/videoApi';

// =============================================================================
// Configuration
// =============================================================================

const SYNC_INTERVAL = 5000; // Sync every 5 seconds during playback
const SYNC_DEBOUNCE = 1000; // Debounce rapid position changes
const MAX_RETRY_ATTEMPTS = 3;
const RETRY_DELAY = 2000;

// =============================================================================
// Sync State Management
// =============================================================================

interface SyncState {
  lastSyncTime: number;
  pendingSync: boolean;
  syncTimer: NodeJS.Timeout | null;
  debounceTimer: NodeJS.Timeout | null;
  retryCount: number;
}

const syncStates: Map<string, SyncState> = new Map();

/**
 * Get or create sync state for a video
 */
function getSyncState(videoId: string): SyncState {
  if (!syncStates.has(videoId)) {
    syncStates.set(videoId, {
      lastSyncTime: 0,
      pendingSync: false,
      syncTimer: null,
      debounceTimer: null,
      retryCount: 0,
    });
  }
  return syncStates.get(videoId)!;
}

/**
 * Clear sync timers for a video
 */
function clearSyncTimers(videoId: string) {
  const state = syncStates.get(videoId);
  if (state) {
    if (state.syncTimer) {
      clearInterval(state.syncTimer);
      state.syncTimer = null;
    }
    if (state.debounceTimer) {
      clearTimeout(state.debounceTimer);
      state.debounceTimer = null;
    }
  }
}

// =============================================================================
// Sync Functions
// =============================================================================

/**
 * Perform sync with retry logic
 */
async function performSync(
  videoId: string,
  sessionId: string,
  position: number,
  platform: string,
  playbackState: string,
  dispatch: any
): Promise<boolean> {
  const state = getSyncState(videoId);

  try {
    const result = await syncBackendPlaybackPosition(
      videoId,
      sessionId,
      position,
      platform,
      playbackState
    );

    // Update sync status
    dispatch(
      updateSyncStatus({
        videoId,
        status: {
          last_sync_time: new Date().toISOString(),
          synced_platforms: result.synced_platforms,
          conflict_detected: result.conflict_detected || false,
        },
      })
    );

    // Reset retry count on success
    state.retryCount = 0;
    state.lastSyncTime = Date.now();
    state.pendingSync = false;

    return true;
  } catch (error) {
    console.error('Sync failed:', error);

    // Retry logic
    if (state.retryCount < MAX_RETRY_ATTEMPTS) {
      state.retryCount++;
      console.log(`Retrying sync (attempt ${state.retryCount}/${MAX_RETRY_ATTEMPTS})`);

      await new Promise((resolve) => setTimeout(resolve, RETRY_DELAY));
      return performSync(videoId, sessionId, position, platform, playbackState, dispatch);
    }

    state.pendingSync = false;
    return false;
  }
}

/**
 * Debounced sync function
 */
function debouncedSync(
  videoId: string,
  sessionId: string,
  position: number,
  platform: string,
  playbackState: string,
  dispatch: any
) {
  const state = getSyncState(videoId);

  // Clear existing debounce timer
  if (state.debounceTimer) {
    clearTimeout(state.debounceTimer);
  }

  // Set new debounce timer
  state.debounceTimer = setTimeout(() => {
    if (!state.pendingSync) {
      state.pendingSync = true;
      performSync(videoId, sessionId, position, platform, playbackState, dispatch);
    }
  }, SYNC_DEBOUNCE);
}

/**
 * Start automatic sync interval
 */
function startSyncInterval(
  videoId: string,
  sessionId: string,
  platform: string,
  getState: () => RootState,
  dispatch: any
) {
  const state = getSyncState(videoId);

  // Clear existing timer
  clearSyncTimers(videoId);

  // Start new interval
  state.syncTimer = setInterval(() => {
    const currentState = getState();
    const playbackState = currentState.video.playbackState[videoId];

    if (!playbackState) return;

    // Only sync if playing
    if (playbackState.state === 'playing') {
      const position = playbackState.position;
      const currentPlaybackState = playbackState.state;

      // Perform sync
      performSync(
        videoId,
        sessionId,
        position,
        platform,
        currentPlaybackState,
        dispatch
      );
    }
  }, SYNC_INTERVAL);
}

/**
 * Stop automatic sync interval
 */
function stopSyncInterval(videoId: string) {
  clearSyncTimers(videoId);
  syncStates.delete(videoId);
}

// =============================================================================
// Middleware
// =============================================================================

export const videoSyncMiddleware: Middleware<{}, RootState> = (store) => (next) => (action) => {
  const result = next(action);

  // Get current state after action
  const state = store.getState();

  // Handle playback position updates
  if (updatePlaybackPosition.match(action)) {
    const { videoId, position, sessionId } = action.payload;

    if (sessionId) {
      const playbackState = state.video.playbackState[videoId];
      if (playbackState && playbackState.state === 'playing') {
        // Debounce sync during playback
        debouncedSync(
          videoId,
          sessionId,
          position,
          'web', // Platform - could be dynamic based on user agent
          playbackState.state,
          store.dispatch
        );
      }
    }
  }

  // Handle playback state changes
  if (setPlaybackState.match(action)) {
    const { videoId, state: newState } = action.payload;
    const playbackState = state.video.playbackState[videoId];

    if (playbackState && playbackState.sessionId) {
      const sessionId = playbackState.sessionId;
      const position = playbackState.position;

      // Start/stop sync interval based on playback state
      if (newState === 'playing') {
        // Start automatic sync
        startSyncInterval(
          videoId,
          sessionId,
          'web',
          store.getState,
          store.dispatch
        );

        // Immediate sync on play
        performSync(
          videoId,
          sessionId,
          position,
          'web',
          'playing',
          store.dispatch
        );
      } else if (newState === 'paused' || newState === 'ended') {
        // Stop automatic sync
        stopSyncInterval(videoId);

        // Final sync on pause/end
        performSync(
          videoId,
          sessionId,
          position,
          'web',
          newState,
          store.dispatch
        );
      }
    }
  }

  // Handle sync playback fulfilled
  if (syncPlayback.fulfilled.match(action)) {
    const { videoId } = action.payload;
    const syncState = getSyncState(videoId);
    syncState.lastSyncTime = Date.now();
    syncState.pendingSync = false;
  }

  // Handle sync playback rejected
  if (syncPlayback.rejected.match(action)) {
    const videoId = action.meta.arg.videoId;
    const syncState = getSyncState(videoId);

    // Retry logic
    if (syncState.retryCount < MAX_RETRY_ATTEMPTS) {
      syncState.retryCount++;
      console.log(`Retrying sync (attempt ${syncState.retryCount}/${MAX_RETRY_ATTEMPTS})`);

      setTimeout(() => {
        store.dispatch(syncPlayback(action.meta.arg));
      }, RETRY_DELAY);
    } else {
      console.error('Max retry attempts reached for sync');
      syncState.pendingSync = false;
    }
  }

  return result;
};

// =============================================================================
// Cleanup
// =============================================================================

/**
 * Clean up all sync timers (call on app unmount)
 */
export function cleanupVideoSync() {
  for (const [videoId] of syncStates) {
    clearSyncTimers(videoId);
  }
  syncStates.clear();
}

export default videoSyncMiddleware;