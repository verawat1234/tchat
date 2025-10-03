import { createListenerMiddleware, isAnyOf } from '@reduxjs/toolkit';
import type { RootState } from '../index';
import { logout, setTokens } from '../../features/authSlice';
import { api } from '../../services/api';
import { authApi } from '../../services/auth';

console.log('[Auth Middleware] Module loaded! authApi:', authApi);
console.log('[Auth Middleware] verifyOTP matcher:', authApi.endpoints.verifyOTP.matchFulfilled);

export const authMiddleware = createListenerMiddleware();

// Debug: Log ALL actions to see if middleware is running
authMiddleware.startListening({
  predicate: () => true,
  effect: (action) => {
    if (action.type?.includes('api/') || action.type?.includes('auth')) {
      console.log('[Auth Middleware] ALL ACTIONS:', action.type, action);
    }
  }
});

// Listen for successful OTP verification and save tokens directly to localStorage
authMiddleware.startListening({
  matcher: authApi.endpoints.verifyOTP.matchFulfilled,
  effect: async (action, listenerApi) => {
    console.log('[Auth Middleware] verifyOTP fulfilled - saving tokens to localStorage');
    const { payload } = action;
    const expiresAt = Date.now() + (payload.expiresIn * 1000);

    // Save tokens directly to localStorage
    localStorage.setItem('accessToken', payload.accessToken);
    localStorage.setItem('refreshToken', payload.refreshToken);
    localStorage.setItem('expiresAt', String(expiresAt));

    console.log('[Auth Middleware] Tokens saved:', {
      accessToken: payload.accessToken.substring(0, 30) + '...',
      refreshToken: payload.refreshToken.substring(0, 30) + '...',
      expiresAt: new Date(expiresAt).toISOString(),
    });
  },
});

// Listen for successful token refresh and save tokens to localStorage
authMiddleware.startListening({
  matcher: authApi.endpoints.refreshToken.matchFulfilled,
  effect: async (action, listenerApi) => {
    console.log('[Auth Middleware] refreshToken fulfilled - saving tokens to localStorage');
    const { payload } = action;
    const expiresAt = Date.now() + (payload.expiresIn * 1000);

    // Save tokens directly to localStorage
    localStorage.setItem('accessToken', payload.accessToken);
    localStorage.setItem('refreshToken', payload.refreshToken);
    localStorage.setItem('expiresAt', String(expiresAt));

    console.log('[Auth Middleware] Tokens refreshed and saved:', {
      accessToken: payload.accessToken.substring(0, 30) + '...',
      expiresAt: new Date(expiresAt).toISOString(),
    });
  },
});

// Listen for auth-related failures to handle logout
authMiddleware.startListening({
  matcher: isAnyOf(
    authApi.endpoints.refreshToken.matchRejected,
    authApi.endpoints.getCurrentUser.matchRejected
  ),
  effect: async (action, listenerApi) => {
    const { error } = action;

    // If refresh token is invalid or expired, logout user
    if (error && 'status' in error && (error.status === 401 || error.status === 403)) {
      listenerApi.dispatch(logout());
      listenerApi.dispatch(api.util.resetApiState());
    }
  },
});

// Save tokens to localStorage and setup auto-refresh when setTokens is dispatched
authMiddleware.startListening({
  actionCreator: setTokens,
  effect: async (action, listenerApi) => {
    console.log('[Auth Middleware] setTokens action received!');
    const { accessToken, refreshToken, expiresIn } = action.payload;
    const expiresAt = Date.now() + (expiresIn * 1000);

    // Save tokens to localStorage for persistence
    localStorage.setItem('accessToken', accessToken);
    localStorage.setItem('refreshToken', refreshToken);
    localStorage.setItem('expiresAt', String(expiresAt));

    console.log('[Auth Middleware] Tokens saved to localStorage:', {
      accessToken: accessToken.substring(0, 30) + '...',
      refreshToken: refreshToken.substring(0, 30) + '...',
      expiresAt: new Date(expiresAt).toISOString(),
      savedAccessToken: localStorage.getItem('accessToken')?.substring(0, 30) + '...',
      savedRefreshToken: localStorage.getItem('refreshToken')?.substring(0, 30) + '...'
    });

    // Set a timer to auto-refresh token before expiry
    const refreshBuffer = 5 * 60 * 1000; // 5 minutes before expiry
    const refreshAt = expiresAt - refreshBuffer;
    const timeUntilRefresh = refreshAt - Date.now();

    console.log('[Auth Middleware] Auto-refresh timer:', {
      expiresIn,
      expiresAt: new Date(expiresAt).toISOString(),
      refreshBuffer,
      refreshAt: new Date(refreshAt).toISOString(),
      timeUntilRefresh,
      willSchedule: timeUntilRefresh > 0
    });

    // Only schedule auto-refresh if token expiry is more than 6 minutes away
    // This prevents immediate refresh on fresh tokens
    if (timeUntilRefresh > 60000) { // At least 1 minute in the future
      console.log('[Auth Middleware] Scheduling auto-refresh in', timeUntilRefresh, 'ms');
      setTimeout(async () => {
        const state = listenerApi.getState() as RootState;
        const refreshToken = state.auth.refreshToken;

        if (refreshToken) {
          console.log('[Auth Middleware] Auto-refresh timer fired, refreshing token');
          try {
            await listenerApi.dispatch(
              api.endpoints.refreshToken.initiate({ refreshToken })
            ).unwrap();
          } catch (error) {
            console.error('[Auth Middleware] Auto-refresh failed:', error);
            // If refresh fails, logout user
            listenerApi.dispatch(logout());
          }
        }
      }, timeUntilRefresh);
    } else {
      console.log('[Auth Middleware] Skipping auto-refresh timer - timeUntilRefresh too small or negative');
    }
  },
});

// Clear tokens from localStorage on logout and redirect to home
authMiddleware.startListening({
  actionCreator: logout,
  effect: async (action, listenerApi) => {
    console.log('[Auth Middleware] Logout triggered, clearing tokens and redirecting');

    // Clear persisted auth data
    localStorage.removeItem('persist:auth');
    localStorage.removeItem('authToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('accessToken');
    localStorage.removeItem('expiresAt');

    // Reset API state to clear any cached user data
    listenerApi.dispatch(api.util.resetApiState());

    // Redirect to home page
    window.location.href = '/';
    console.log('[Auth Middleware] Logout complete - redirecting to home');
  },
});

// Auto-login after successful OTP verification or token refresh
authMiddleware.startListening({
  matcher: isAnyOf(
    // api.endpoints.verifyOTP.matchFulfilled,
    // api.endpoints.refreshToken.matchFulfilled
  ),
  effect: async (action, listenerApi) => {
    const state = listenerApi.getState() as RootState;
    const { accessToken, expiresAt } = state.auth;

    console.log('[Auth Middleware] Auto-login check:', {
      hasToken: !!accessToken,
      expiresAt: expiresAt ? new Date(expiresAt).toISOString() : 'none',
      isExpired: expiresAt ? Date.now() >= expiresAt : 'no-expiry',
      currentTime: new Date().toISOString()
    });

    // If we have a token and it's not expired, get current user
    if (accessToken && expiresAt && Date.now() < expiresAt) {
      console.log('[Auth Middleware] Fetching current user with valid token');
      try {
        await listenerApi.dispatch(
          api.endpoints.getCurrentUser.initiate()
        ).unwrap();
        console.log('[Auth Middleware] Current user fetched successfully');
      } catch (error) {
        console.error('[Auth Middleware] Failed to get current user:', error);
        // If getting current user fails, logout
        listenerApi.dispatch(logout());
      }
    } else {
      console.log('[Auth Middleware] Skipping getCurrentUser - invalid or expired token');
    }
  },
});

// Auto-login on app startup if valid tokens exist in localStorage
authMiddleware.startListening({
  predicate: (action, currentState, previousState) => {
    // Only trigger on store initialization (when previousState is undefined)
    return !previousState;
  },
  effect: async (action, listenerApi) => {
    const state = listenerApi.getState() as RootState;
    const { accessToken, expiresAt } = state.auth;

    console.log('[Auth Middleware] App startup - checking for existing session:', {
      hasToken: !!accessToken,
      expiresAt: expiresAt ? new Date(expiresAt).toISOString() : 'none',
      isExpired: expiresAt ? Date.now() >= expiresAt : 'no-expiry',
    });

    // If we have a valid token from localStorage, fetch current user
    if (accessToken && expiresAt && Date.now() < expiresAt) {
      console.log('[Auth Middleware] Restoring session - fetching current user');
      try {
        await listenerApi.dispatch(
          api.endpoints.getCurrentUser.initiate()
        ).unwrap();
        console.log('[Auth Middleware] Session restored successfully');
      } catch (error) {
        console.error('[Auth Middleware] Session restoration failed:', error);
        listenerApi.dispatch(logout());
      }
    }
  },
});