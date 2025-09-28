import { configureStore, combineReducers } from '@reduxjs/toolkit';
import { api } from '../services/api';
import authReducer from '../features/authSlice';
import uiReducer from '../features/uiSlice';
import loadingReducer from '../features/loadingSlice';
import contentReducer from '../features/contentSlice';
import { authMiddleware } from './middleware/authMiddleware';
import { errorMiddleware } from './middleware/errorMiddleware';
import { contentFallbackMiddleware } from './middleware/contentFallbackMiddleware';
import { socialMiddleware } from './middleware/socialMiddleware';

// Root reducer without persistence (temporary)
const rootReducer = combineReducers({
  [api.reducerPath]: api.reducer,
  auth: authReducer,
  ui: uiReducer,
  loading: loadingReducer,
  content: contentReducer,
});

export const store = configureStore({
  reducer: rootReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // Ignore these field paths in all actions
        ignoredActionPaths: ['meta.arg', 'payload.timestamp'],
        // Ignore these paths in the state
        ignoredPaths: ['auth.expiresAt'],
      },
    })
    .concat(api.middleware)
    .concat(contentFallbackMiddleware.middleware)
    .concat(socialMiddleware.middleware)
    .prepend(authMiddleware.middleware)
    .prepend(errorMiddleware.middleware),
  devTools: process.env.NODE_ENV !== 'production' && {
    name: 'Tchat Redux Store',
    trace: true,
    traceLimit: 25,
    maxAge: 50,

    // Enhanced action sanitization for security
    actionSanitizer: (action) => {
      const sensitiveActions = ['Token', 'auth', 'login', 'refresh', 'setTokens'];
      const isSensitive = sensitiveActions.some(keyword =>
        action.type.toLowerCase().includes(keyword.toLowerCase())
      );

      if (isSensitive && action.payload) {
        return {
          ...action,
          payload: {
            ...action.payload,
            accessToken: action.payload.accessToken ? '[REDACTED]' : undefined,
            refreshToken: action.payload.refreshToken ? '[REDACTED]' : undefined,
            password: action.payload.password ? '[REDACTED]' : undefined,
            tokens: action.payload.tokens ? {
              ...action.payload.tokens,
              accessToken: '[REDACTED]',
              refreshToken: '[REDACTED]',
            } : undefined,
          },
        };
      }

      return action;
    },

    // Enhanced state sanitization
    stateSanitizer: (state) => ({
      ...state,
      auth: {
        ...state.auth,
        accessToken: state.auth?.accessToken ? '[REDACTED]' : null,
        refreshToken: state.auth?.refreshToken ? '[REDACTED]' : null,
      },
      // Hide any cached sensitive API data
      api: {
        ...state.api,
        queries: Object.keys(state.api?.queries || {}).reduce((acc, key) => {
          const query = state.api.queries[key];
          if (key.includes('auth') || key.includes('login')) {
            acc[key] = {
              ...query,
              data: query.data ? '[REDACTED]' : query.data,
            };
          } else {
            acc[key] = query;
          }
          return acc;
        }, {}),
      },
    }),

    // Filter out noisy actions in development
    predicate: (state, action) => {
      const noisyActions = [
        'listenerMiddleware',
        'api/internal',
      ];

      // Hide noisy actions but allow important ones
      return !noisyActions.some(noisy => action.type.includes(noisy));
    },

    // Advanced features for better debugging
    serialize: {
      options: {
        undefined: true,
        function: true,
        symbol: true,
      },
    },

    // Action type filtering for cleaner debugging
    actionsBlacklist: [
      'api/internal',
      'listenerMiddleware',
    ],

    // Features configuration
    features: {
      pause: true,
      lock: true,
      persist: true,
      export: true,
      import: 'custom',
      jump: true,
      skip: true,
      reorder: true,
      dispatch: true,
      test: true,
    },
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;