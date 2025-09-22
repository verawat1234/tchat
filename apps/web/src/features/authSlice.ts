import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import type { User, AuthTokens } from '../types/api';
import { authApi } from '../services/auth';

interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  expiresAt: number | null;
}

const initialState: AuthState = {
  isAuthenticated: false,
  user: null,
  accessToken: null,
  refreshToken: null,
  expiresAt: null,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setTokens: (state, action: PayloadAction<{
      accessToken: string;
      refreshToken: string;
      expiresIn: number;
    }>) => {
      state.accessToken = action.payload.accessToken;
      state.refreshToken = action.payload.refreshToken;
      state.expiresAt = Date.now() + (action.payload.expiresIn * 1000);
      state.isAuthenticated = true;
    },
    setUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload;
    },
    logout: (state) => {
      state.isAuthenticated = false;
      state.user = null;
      state.accessToken = null;
      state.refreshToken = null;
      state.expiresAt = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addMatcher(
        authApi.endpoints.login.matchFulfilled,
        (state, { payload }) => {
          state.user = payload.user;
          state.accessToken = payload.tokens.accessToken;
          state.refreshToken = payload.tokens.refreshToken;
          state.expiresAt = Date.now() + (payload.tokens.expiresIn * 1000);
          state.isAuthenticated = true;
        }
      )
      .addMatcher(
        authApi.endpoints.logout.matchFulfilled,
        (state) => {
          state.isAuthenticated = false;
          state.user = null;
          state.accessToken = null;
          state.refreshToken = null;
          state.expiresAt = null;
        }
      )
      .addMatcher(
        authApi.endpoints.getCurrentUser.matchFulfilled,
        (state, { payload }) => {
          state.user = payload;
          state.isAuthenticated = true;
        }
      )
      .addMatcher(
        authApi.endpoints.refreshToken.matchFulfilled,
        (state, { payload }) => {
          state.accessToken = payload.accessToken;
          state.refreshToken = payload.refreshToken;
          state.expiresAt = Date.now() + (payload.expiresIn * 1000);
        }
      );
  },
});

export const { setTokens, setUser, logout } = authSlice.actions;
export default authSlice.reducer;