import { describe, it, expect, beforeEach } from 'vitest';
import authReducer, { setTokens, setUser, logout } from '../../src/features/authSlice';
import { authApi } from '../../src/services/auth';
import type { User } from '../../src/types/api';

describe('Auth Slice', () => {
  const initialState = {
    isAuthenticated: false,
    user: null,
    accessToken: null,
    refreshToken: null,
    expiresAt: null,
  };

  const mockUser: User = {
    id: '1',
    email: 'test@example.com',
    username: 'testuser',
    displayName: 'Test User',
    role: 'user',
    createdAt: '2023-01-01T00:00:00.000Z',
    updatedAt: '2023-01-01T00:00:00.000Z',
  };

  const mockTokens = {
    accessToken: 'mock-access-token',
    refreshToken: 'mock-refresh-token',
    expiresIn: 3600,
  };

  beforeEach(() => {
    // Reset Date.now for consistent testing
    vi.setSystemTime(new Date('2023-01-01T00:00:00.000Z'));
  });

  describe('Initial State', () => {
    it('should return the initial state', () => {
      expect(authReducer(undefined, { type: 'unknown' })).toEqual(initialState);
    });
  });

  describe('Actions', () => {
    describe('setTokens', () => {
      it('should set tokens and mark as authenticated', () => {
        const action = setTokens(mockTokens);
        const state = authReducer(initialState, action);

        expect(state).toEqual({
          ...initialState,
          isAuthenticated: true,
          accessToken: 'mock-access-token',
          refreshToken: 'mock-refresh-token',
          expiresAt: Date.now() + (3600 * 1000),
        });
      });

      it('should calculate correct expiration time', () => {
        const currentTime = Date.now();
        const action = setTokens({ ...mockTokens, expiresIn: 7200 });
        const state = authReducer(initialState, action);

        expect(state.expiresAt).toBe(currentTime + (7200 * 1000));
      });
    });

    describe('setUser', () => {
      it('should set user data', () => {
        const action = setUser(mockUser);
        const state = authReducer(initialState, action);

        expect(state).toEqual({
          ...initialState,
          user: mockUser,
        });
      });

      it('should update user data if already set', () => {
        const existingState = {
          ...initialState,
          user: mockUser,
        };

        const updatedUser = {
          ...mockUser,
          displayName: 'Updated Name',
        };

        const action = setUser(updatedUser);
        const state = authReducer(existingState, action);

        expect(state.user?.displayName).toBe('Updated Name');
      });
    });

    describe('logout', () => {
      it('should reset all auth state', () => {
        const authenticatedState = {
          isAuthenticated: true,
          user: mockUser,
          accessToken: 'token',
          refreshToken: 'refresh',
          expiresAt: Date.now() + 3600000,
        };

        const action = logout();
        const state = authReducer(authenticatedState, action);

        expect(state).toEqual(initialState);
      });
    });
  });

  describe('Extra Reducers - API Responses', () => {
    describe('login endpoint', () => {
      it('should handle login fulfilled', () => {
        const loginResponse = {
          user: mockUser,
          tokens: {
            accessToken: 'new-access-token',
            refreshToken: 'new-refresh-token',
            expiresIn: 3600,
            tokenType: 'Bearer' as const,
          },
        };

        const action = {
          type: authApi.endpoints.login.matchFulfilled.type,
          payload: loginResponse,
        };

        const state = authReducer(initialState, action);

        expect(state).toEqual({
          isAuthenticated: true,
          user: mockUser,
          accessToken: 'new-access-token',
          refreshToken: 'new-refresh-token',
          expiresAt: Date.now() + (3600 * 1000),
        });
      });
    });

    describe('logout endpoint', () => {
      it('should handle logout fulfilled', () => {
        const authenticatedState = {
          isAuthenticated: true,
          user: mockUser,
          accessToken: 'token',
          refreshToken: 'refresh',
          expiresAt: Date.now() + 3600000,
        };

        const action = {
          type: authApi.endpoints.logout.matchFulfilled.type,
          payload: undefined,
        };

        const state = authReducer(authenticatedState, action);

        expect(state).toEqual(initialState);
      });
    });

    describe('getCurrentUser endpoint', () => {
      it('should handle getCurrentUser fulfilled', () => {
        const stateWithTokens = {
          ...initialState,
          isAuthenticated: false,
          accessToken: 'token',
          refreshToken: 'refresh',
          expiresAt: Date.now() + 3600000,
        };

        const action = {
          type: authApi.endpoints.getCurrentUser.matchFulfilled.type,
          payload: mockUser,
        };

        const state = authReducer(stateWithTokens, action);

        expect(state).toEqual({
          ...stateWithTokens,
          isAuthenticated: true,
          user: mockUser,
        });
      });
    });

    describe('refreshToken endpoint', () => {
      it('should handle refreshToken fulfilled', () => {
        const stateWithOldTokens = {
          ...initialState,
          isAuthenticated: true,
          user: mockUser,
          accessToken: 'old-token',
          refreshToken: 'old-refresh',
          expiresAt: Date.now() - 1000, // Expired
        };

        const refreshResponse = {
          accessToken: 'new-access-token',
          refreshToken: 'new-refresh-token',
          expiresIn: 3600,
        };

        const action = {
          type: authApi.endpoints.refreshToken.matchFulfilled.type,
          payload: refreshResponse,
        };

        const state = authReducer(stateWithOldTokens, action);

        expect(state).toEqual({
          ...stateWithOldTokens,
          accessToken: 'new-access-token',
          refreshToken: 'new-refresh-token',
          expiresAt: Date.now() + (3600 * 1000),
        });
      });
    });
  });

  describe('State Transitions', () => {
    it('should handle complete auth flow', () => {
      let state = authReducer(undefined, { type: 'unknown' });
      expect(state.isAuthenticated).toBe(false);

      // Login
      const loginAction = {
        type: authApi.endpoints.login.matchFulfilled.type,
        payload: {
          user: mockUser,
          tokens: {
            accessToken: 'access-token',
            refreshToken: 'refresh-token',
            expiresIn: 3600,
            tokenType: 'Bearer' as const,
          },
        },
      };
      state = authReducer(state, loginAction);
      expect(state.isAuthenticated).toBe(true);
      expect(state.user).toEqual(mockUser);

      // Token refresh
      const refreshAction = {
        type: authApi.endpoints.refreshToken.matchFulfilled.type,
        payload: {
          accessToken: 'new-access-token',
          refreshToken: 'new-refresh-token',
          expiresIn: 3600,
        },
      };
      state = authReducer(state, refreshAction);
      expect(state.accessToken).toBe('new-access-token');
      expect(state.isAuthenticated).toBe(true);

      // Logout
      const logoutAction = logout();
      state = authReducer(state, logoutAction);
      expect(state.isAuthenticated).toBe(false);
      expect(state.user).toBe(null);
      expect(state.accessToken).toBe(null);
    });
  });

  describe('Edge Cases', () => {
    it('should handle setTokens with zero expiration', () => {
      const action = setTokens({ ...mockTokens, expiresIn: 0 });
      const state = authReducer(initialState, action);

      expect(state.expiresAt).toBe(Date.now());
      expect(state.isAuthenticated).toBe(true);
    });

    it('should handle setUser with null', () => {
      const stateWithUser = {
        ...initialState,
        user: mockUser,
      };

      const action = setUser(null as any);
      const state = authReducer(stateWithUser, action);

      expect(state.user).toBe(null);
    });

    it('should maintain other state when setting tokens', () => {
      const stateWithUser = {
        ...initialState,
        user: mockUser,
      };

      const action = setTokens(mockTokens);
      const state = authReducer(stateWithUser, action);

      expect(state.user).toEqual(mockUser);
      expect(state.isAuthenticated).toBe(true);
    });
  });
});