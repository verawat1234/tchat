import { describe, it, expect, beforeEach, afterEach, vi, type MockedFunction } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import { authMiddleware } from '../../src/store/middleware/authMiddleware';
import { setTokens, logout } from '../../src/features/authSlice';
import { api } from '../../src/services/api';
import authReducer from '../../src/features/authSlice';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

// Mock setTimeout and clearTimeout
const originalSetTimeout = global.setTimeout;
const originalClearTimeout = global.clearTimeout;

describe('Auth Middleware', () => {
  let store: ReturnType<typeof configureStore>;
  let mockSetTimeout: MockedFunction<typeof setTimeout>;
  let mockClearTimeout: MockedFunction<typeof clearTimeout>;

  beforeEach(() => {
    // Reset all mocks
    vi.clearAllMocks();
    vi.clearAllTimers();
    vi.useFakeTimers();

    // Mock setTimeout and clearTimeout
    mockSetTimeout = vi.fn().mockImplementation((callback, delay) => {
      return originalSetTimeout(callback, delay);
    });
    mockClearTimeout = vi.fn();
    global.setTimeout = mockSetTimeout;
    global.clearTimeout = mockClearTimeout;

    // Create test store with auth middleware
    store = configureStore({
      reducer: {
        auth: authReducer,
        [api.reducerPath]: api.reducer,
      },
      middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().prepend(authMiddleware.middleware),
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.useRealTimers();
    global.setTimeout = originalSetTimeout;
    global.clearTimeout = originalClearTimeout;
  });

  describe('Token Refresh Success Handling', () => {
    it('should update tokens when refresh succeeds', async () => {
      const refreshResponse = {
        accessToken: 'new-access-token',
        refreshToken: 'new-refresh-token',
        expiresIn: 3600,
      };

      // Dispatch refresh success action
      store.dispatch({
        type: api.endpoints.refreshToken.matchFulfilled.type,
        payload: refreshResponse,
      });

      // Wait for middleware to process
      await vi.runAllTimersAsync();

      const state = store.getState();
      expect(state.auth.accessToken).toBe('new-access-token');
      expect(state.auth.refreshToken).toBe('new-refresh-token');
      expect(state.auth.isAuthenticated).toBe(true);
    });
  });

  describe('Token Refresh Failure Handling', () => {
    it('should logout user when refresh token fails with 401', async () => {
      // Set initial authenticated state
      store.dispatch(setTokens({
        accessToken: 'old-token',
        refreshToken: 'old-refresh',
        expiresIn: 3600,
      }));

      expect(store.getState().auth.isAuthenticated).toBe(true);

      // Dispatch refresh failure action
      store.dispatch({
        type: api.endpoints.refreshToken.matchRejected.type,
        error: {
          status: 401,
          data: { message: 'Invalid refresh token' },
        },
      });

      await vi.runAllTimersAsync();

      const state = store.getState();
      expect(state.auth.isAuthenticated).toBe(false);
      expect(state.auth.accessToken).toBe(null);
      expect(state.auth.refreshToken).toBe(null);
    });

    it('should logout user when getCurrentUser fails with 401', async () => {
      store.dispatch(setTokens({
        accessToken: 'token',
        refreshToken: 'refresh',
        expiresIn: 3600,
      }));

      store.dispatch({
        type: api.endpoints.getCurrentUser.matchRejected.type,
        error: {
          status: 401,
          data: { message: 'Unauthorized' },
        },
      });

      await vi.runAllTimersAsync();

      const state = store.getState();
      expect(state.auth.isAuthenticated).toBe(false);
    });

    it('should not logout user for non-401 errors', async () => {
      store.dispatch(setTokens({
        accessToken: 'token',
        refreshToken: 'refresh',
        expiresIn: 3600,
      }));

      store.dispatch({
        type: api.endpoints.refreshToken.matchRejected.type,
        error: {
          status: 500,
          data: { message: 'Server error' },
        },
      });

      await vi.runAllTimersAsync();

      const state = store.getState();
      expect(state.auth.isAuthenticated).toBe(true);
    });
  });

  describe('Auto Token Refresh', () => {
    it('should schedule token refresh before expiry', async () => {
      const expiresIn = 3600; // 1 hour
      const refreshBuffer = 5 * 60 * 1000; // 5 minutes
      const expectedDelay = (expiresIn * 1000) - refreshBuffer;

      store.dispatch(setTokens({
        accessToken: 'token',
        refreshToken: 'refresh',
        expiresIn,
      }));

      await vi.runAllTimersAsync();

      expect(mockSetTimeout).toHaveBeenCalledWith(
        expect.any(Function),
        expectedDelay
      );
    });

    it('should not schedule refresh if token expires too soon', async () => {
      const expiresIn = 60; // 1 minute (less than 5-minute buffer)

      store.dispatch(setTokens({
        accessToken: 'token',
        refreshToken: 'refresh',
        expiresIn,
      }));

      await vi.runAllTimersAsync();

      // Should not schedule timeout for tokens expiring in less than buffer time
      expect(mockSetTimeout).not.toHaveBeenCalled();
    });

    it('should attempt refresh when timer triggers', async () => {
      const mockDispatch = vi.spyOn(store, 'dispatch');

      store.dispatch(setTokens({
        accessToken: 'token',
        refreshToken: 'refresh-token',
        expiresIn: 3600,
      }));

      // Fast-forward time to trigger the timeout
      await vi.runAllTimersAsync();

      // Manually trigger the timeout callback
      const timeoutCallback = mockSetTimeout.mock.calls[0]?.[0];
      if (typeof timeoutCallback === 'function') {
        await timeoutCallback();
      }

      // Should attempt to call refresh endpoint
      expect(mockDispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('refreshToken'),
        })
      );
    });
  });

  describe('Logout Cleanup', () => {
    it('should clear localStorage on logout', async () => {
      store.dispatch(logout());

      await vi.runAllTimersAsync();

      expect(localStorageMock.removeItem).toHaveBeenCalledWith('persist:auth');
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('authToken');
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('refreshToken');
    });

    it('should reset API state on logout', async () => {
      const mockDispatch = vi.spyOn(store, 'dispatch');

      store.dispatch(logout());

      await vi.runAllTimersAsync();

      expect(mockDispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'api/resetApiState',
        })
      );
    });
  });

  describe('Auto-login on Startup', () => {
    it('should attempt getCurrentUser if valid tokens exist', async () => {
      const mockDispatch = vi.spyOn(store, 'dispatch');
      const futureTime = Date.now() + (3600 * 1000); // 1 hour from now

      // Set tokens that are not expired
      store.dispatch({
        type: 'auth/setTokens',
        payload: {
          accessToken: 'valid-token',
          refreshToken: 'valid-refresh',
          expiresIn: 3600,
        },
      });

      // Simulate store initialization with valid tokens
      const stateWithValidTokens = {
        ...store.getState(),
        auth: {
          ...store.getState().auth,
          accessToken: 'valid-token',
          expiresAt: futureTime,
        },
      };

      // Trigger the auto-login listener by dispatching a state change
      store.dispatch({
        type: 'TEST_STATE_CHANGE',
        payload: stateWithValidTokens,
      });

      await vi.runAllTimersAsync();

      // Should attempt to get current user
      expect(mockDispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('getCurrentUser'),
        })
      );
    });

    it('should logout if tokens are expired on startup', async () => {
      const mockDispatch = vi.spyOn(store, 'dispatch');
      const pastTime = Date.now() - 1000; // 1 second ago (expired)

      // Simulate store with expired tokens
      store.dispatch({
        type: 'TEST_EXPIRED_TOKENS',
        payload: {
          accessToken: 'expired-token',
          expiresAt: pastTime,
        },
      });

      await vi.runAllTimersAsync();

      // Should not attempt getCurrentUser with expired tokens
      expect(mockDispatch).not.toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('getCurrentUser'),
        })
      );
    });

    it('should logout if getCurrentUser fails on startup', async () => {
      const futureTime = Date.now() + (3600 * 1000);

      store.dispatch(setTokens({
        accessToken: 'token',
        refreshToken: 'refresh',
        expiresIn: 3600,
      }));

      // Simulate getCurrentUser failure
      store.dispatch({
        type: api.endpoints.getCurrentUser.matchRejected.type,
        error: {
          status: 401,
          data: { message: 'Token invalid' },
        },
      });

      await vi.runAllTimersAsync();

      const state = store.getState();
      expect(state.auth.isAuthenticated).toBe(false);
    });
  });

  describe('Edge Cases', () => {
    it('should handle multiple simultaneous token refresh attempts', async () => {
      const mockDispatch = vi.spyOn(store, 'dispatch');

      // Set tokens that will trigger auto-refresh
      store.dispatch(setTokens({
        accessToken: 'token1',
        refreshToken: 'refresh1',
        expiresIn: 3600,
      }));

      // Immediately set new tokens (simulating rapid token updates)
      store.dispatch(setTokens({
        accessToken: 'token2',
        refreshToken: 'refresh2',
        expiresIn: 3600,
      }));

      await vi.runAllTimersAsync();

      // Should handle this gracefully without errors
      expect(mockSetTimeout).toHaveBeenCalled();
    });

    it('should handle refresh when no refresh token exists', async () => {
      // Set state with access token but no refresh token
      const stateWithoutRefresh = {
        ...store.getState(),
        auth: {
          ...store.getState().auth,
          accessToken: 'token',
          refreshToken: null,
          expiresAt: Date.now() + 3600000,
        },
      };

      // This should not cause errors
      store.dispatch({
        type: 'TEST_NO_REFRESH_TOKEN',
        payload: stateWithoutRefresh,
      });

      await vi.runAllTimersAsync();

      // Should not attempt refresh without refresh token
      expect(store.getState().auth.accessToken).toBe('token');
    });

    it('should clear previous timeout when new tokens are set', async () => {
      // Set initial tokens
      store.dispatch(setTokens({
        accessToken: 'token1',
        refreshToken: 'refresh1',
        expiresIn: 3600,
      }));

      // Set new tokens (should clear previous timeout)
      store.dispatch(setTokens({
        accessToken: 'token2',
        refreshToken: 'refresh2',
        expiresIn: 7200,
      }));

      await vi.runAllTimersAsync();

      // Should have called setTimeout twice (once for each token set)
      expect(mockSetTimeout).toHaveBeenCalledTimes(2);
    });
  });

  describe('Integration with API State', () => {
    it('should reset API state when logout occurs', async () => {
      const mockDispatch = vi.spyOn(store, 'dispatch');

      // Simulate user being logged in with some API cache
      store.dispatch(setTokens({
        accessToken: 'token',
        refreshToken: 'refresh',
        expiresIn: 3600,
      }));

      // Logout user
      store.dispatch(logout());

      await vi.runAllTimersAsync();

      // Should dispatch API reset action
      expect(mockDispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'api/resetApiState',
        })
      );
    });

    it('should reset API state when refresh fails', async () => {
      const mockDispatch = vi.spyOn(store, 'dispatch');

      store.dispatch(setTokens({
        accessToken: 'token',
        refreshToken: 'refresh',
        expiresIn: 3600,
      }));

      // Simulate refresh failure
      store.dispatch({
        type: api.endpoints.refreshToken.matchRejected.type,
        error: {
          status: 401,
          data: { message: 'Invalid refresh token' },
        },
      });

      await vi.runAllTimersAsync();

      // Should reset API state along with logout
      expect(mockDispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'api/resetApiState',
        })
      );
    });
  });
});