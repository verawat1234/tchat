import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { api } from '../services/api';

interface LoadingState {
  global: boolean;
  requests: Record<string, boolean>;
  operations: Record<string, {
    isLoading: boolean;
    progress?: number;
    message?: string;
  }>;
}

const initialState: LoadingState = {
  global: false,
  requests: {},
  operations: {},
};

const loadingSlice = createSlice({
  name: 'loading',
  initialState,
  reducers: {
    setGlobalLoading: (state, action: PayloadAction<boolean>) => {
      state.global = action.payload;
    },
    setRequestLoading: (state, action: PayloadAction<{ key: string; isLoading: boolean }>) => {
      const { key, isLoading } = action.payload;
      if (isLoading) {
        state.requests[key] = true;
      } else {
        delete state.requests[key];
      }
    },
    setOperationLoading: (state, action: PayloadAction<{
      key: string;
      isLoading: boolean;
      progress?: number;
      message?: string;
    }>) => {
      const { key, isLoading, progress, message } = action.payload;
      if (isLoading) {
        state.operations[key] = { isLoading, progress, message };
      } else {
        delete state.operations[key];
      }
    },
    clearAllLoading: (state) => {
      state.global = false;
      state.requests = {};
      state.operations = {};
    },
  },
  extraReducers: (builder) => {
    // Auto-track RTK Query loading states
    // Note: These matchers will be added when specific endpoints are injected

    // Generic matchers for any API pending/fulfilled/rejected actions
    builder
      .addMatcher(
        (action) => action.type.endsWith('/pending') && action.type.includes('api/'),
        (state, action) => {
          const endpointName = action.type.split('/')[1];
          if (endpointName) {
            state.operations[endpointName] = {
              isLoading: true,
              message: `Loading ${endpointName}...`,
            };
          }
        }
      )
      .addMatcher(
        (action) => (action.type.endsWith('/fulfilled') || action.type.endsWith('/rejected')) && action.type.includes('api/'),
        (state, action) => {
          const endpointName = action.type.split('/')[1];
          if (endpointName && state.operations[endpointName]) {
            delete state.operations[endpointName];
          }
        }
      );
  },
});

export const {
  setGlobalLoading,
  setRequestLoading,
  setOperationLoading,
  clearAllLoading,
} = loadingSlice.actions;

export default loadingSlice.reducer;

// Selectors
export const selectGlobalLoading = (state: { loading: LoadingState }) => state.loading.global;
export const selectRequestLoading = (key: string) => (state: { loading: LoadingState }) =>
  state.loading.requests[key] || false;
export const selectOperationLoading = (key: string) => (state: { loading: LoadingState }) =>
  state.loading.operations[key];
export const selectAnyLoading = (state: { loading: LoadingState }) =>
  state.loading.global ||
  Object.keys(state.loading.requests).length > 0 ||
  Object.keys(state.loading.operations).length > 0;