import { createListenerMiddleware, isAnyOf } from '@reduxjs/toolkit';
import type { RootState } from '../index';
import { logout, setTokens } from '../../features/authSlice';
import { api } from '../../services/api';

export const authMiddleware = createListenerMiddleware();

// Listen for token refresh success to update local state
authMiddleware.startListening({
  matcher: isAnyOf(api.endpoints.refreshToken.matchFulfilled),
  effect: async (action, listenerApi) => {
    const { payload } = action;

    listenerApi.dispatch(setTokens({
      accessToken: payload.accessToken,
      refreshToken: payload.refreshToken,
      expiresIn: payload.expiresIn,
    }));
  },
});

// Listen for auth-related failures to handle logout
authMiddleware.startListening({
  matcher: isAnyOf(
    api.endpoints.refreshToken.matchRejected,
    api.endpoints.getCurrentUser.matchRejected
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

// Auto-logout when token expires
authMiddleware.startListening({
  actionCreator: setTokens,
  effect: async (action, listenerApi) => {
    const { expiresIn } = action.payload;
    const expiresAt = Date.now() + (expiresIn * 1000);

    // Set a timer to auto-refresh token before expiry
    const refreshBuffer = 5 * 60 * 1000; // 5 minutes before expiry
    const refreshAt = expiresAt - refreshBuffer;
    const timeUntilRefresh = refreshAt - Date.now();

    if (timeUntilRefresh > 0) {
      setTimeout(async () => {
        const state = listenerApi.getState() as RootState;
        const refreshToken = state.auth.refreshToken;

        if (refreshToken) {
          try {
            await listenerApi.dispatch(
              api.endpoints.refreshToken.initiate({ refreshToken })
            ).unwrap();
          } catch (error) {
            // If refresh fails, logout user
            listenerApi.dispatch(logout());
          }
        }
      }, timeUntilRefresh);
    }
  },
});

// Clear tokens from localStorage on logout
authMiddleware.startListening({
  actionCreator: logout,
  effect: async (action, listenerApi) => {
    // Clear persisted auth data
    localStorage.removeItem('persist:auth');
    localStorage.removeItem('authToken');
    localStorage.removeItem('refreshToken');

    // Reset API state to clear any cached user data
    listenerApi.dispatch(api.util.resetApiState());
  },
});

// Auto-login on app startup if valid tokens exist
authMiddleware.startListening({
  predicate: (action, currentState, previousState) => {
    // Trigger on store initialization or auth state changes
    return !previousState ||
           (currentState as RootState).auth.accessToken !==
           (previousState as RootState).auth?.accessToken;
  },
  effect: async (action, listenerApi) => {
    const state = listenerApi.getState() as RootState;
    const { accessToken, expiresAt } = state.auth;

    // If we have a token and it's not expired, get current user
    if (accessToken && expiresAt && Date.now() < expiresAt) {
      try {
        await listenerApi.dispatch(
          api.endpoints.getCurrentUser.initiate()
        ).unwrap();
      } catch (error) {
        // If getting current user fails, logout
        listenerApi.dispatch(logout());
      }
    }
  },
});