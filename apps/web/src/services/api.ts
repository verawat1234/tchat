import { createApi, fetchBaseQuery, BaseQueryFn, FetchArgs, FetchBaseQueryError, retry } from '@reduxjs/toolkit/query/react';
import type { RootState } from '../store';
import { RefreshTokenResponse } from '../types/api';

const baseQuery = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_API_URL || 'http://localhost:3001/api',
  prepareHeaders: (headers, { getState }) => {
    const token = (getState() as RootState).auth?.accessToken;
    if (token) {
      headers.set('authorization', `Bearer ${token}`);
    }
    headers.set('Content-Type', 'application/json');
    return headers;
  },
  credentials: 'include',
  timeout: 10000, // 10 second timeout
});

// Enhanced retry logic with exponential backoff
const baseQueryWithRetry = retry(
  async (args, api, extraOptions) => {
    const result = await baseQuery(args, api, extraOptions);

    // Retry on network errors, timeouts, and 5xx errors
    if (result.error) {
      const { status } = result.error as FetchBaseQueryError;
      if (
        status === 'FETCH_ERROR' ||
        status === 'TIMEOUT_ERROR' ||
        (typeof status === 'number' && status >= 500)
      ) {
        throw result.error;
      }
    }

    return result;
  },
  {
    maxRetries: 3,
    backoff: (attempt) => {
      // Exponential backoff: 1s, 2s, 4s
      return Math.min(1000 * Math.pow(2, attempt), 10000);
    },
  }
);

const baseQueryWithReauth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  let result = await baseQueryWithRetry(args, api, extraOptions);

  if (result.error && result.error.status === 401) {
    // Try to get a new token
    const refreshToken = (api.getState() as RootState).auth?.refreshToken;

    if (refreshToken) {
      const refreshResult = await baseQuery(
        {
          url: '/auth/refresh',
          method: 'POST',
          body: { refreshToken },
        },
        api,
        extraOptions
      );

      if (refreshResult.data) {
        const data = refreshResult.data as RefreshTokenResponse;
        // Store the new tokens
        api.dispatch({
          type: 'auth/setTokens',
          payload: {
            accessToken: data.accessToken,
            refreshToken: data.refreshToken,
            expiresIn: data.expiresIn,
          }
        });

        // Retry the original query with new token
        result = await baseQuery(args, api, extraOptions);
      } else {
        // Refresh failed - logout user
        api.dispatch({ type: 'auth/logout' });
      }
    } else {
      // No refresh token available - logout user
      api.dispatch({ type: 'auth/logout' });
    }
  }

  return result;
};

/**
 * Main RTK Query API service configuration
 *
 * Provides the base API configuration with automatic authentication,
 * retry logic, caching, and request deduplication. All specific API
 * endpoints are injected into this base configuration.
 *
 * Features:
 * - Automatic JWT token attachment and refresh
 * - Exponential backoff retry logic for network/server errors
 * - 60-second cache TTL with tag-based invalidation
 * - Request deduplication to prevent duplicate API calls
 * - Redux Persist integration for offline support
 * - Comprehensive error handling and user feedback
 *
 * @example
 * ```typescript
 * // Injecting endpoints into the API
 * export const myApi = api.injectEndpoints({
 *   endpoints: (builder) => ({
 *     getItems: builder.query<Item[], void>({
 *       query: () => '/items',
 *       providesTags: ['Item'],
 *     }),
 *   }),
 * });
 * ```
 */
export const api = createApi({
  reducerPath: 'api',
  baseQuery: baseQueryWithReauth,
  tagTypes: [
    'User',
    'Message',
    'Chat',
    'Auth',
    'UserProfile',
    'ChatList',
    'MessageList',
    'Notification',
    'Settings',
    'Content',
    'ContentItem',
    'ContentCategory',
    'ContentVersion'
  ],
  endpoints: () => ({}),
  refetchOnMountOrArgChange: 30,
  refetchOnFocus: true,
  refetchOnReconnect: true,
  keepUnusedDataFor: 60, // 60 seconds cache
  extractRehydrationInfo(action, { reducerPath }) {
    if (action.type === 'persist/REHYDRATE') {
      return action.payload?.[reducerPath];
    }
  },
  // Request deduplication configuration
  serializeQueryArgs: ({ queryArgs, endpointDefinition, endpointName }) => {
    // Custom serialization for request deduplication
    if (endpointName === 'listMessages') {
      // For messages, deduplicate based on chatId only (ignore cursor for deduplication)
      return `${endpointName}(${JSON.stringify({ chatId: queryArgs.chatId })})`;
    }

    if (endpointName === 'listUsers') {
      // For users, deduplicate based on search and limit
      const { search, limit } = queryArgs;
      return `${endpointName}(${JSON.stringify({ search, limit })})`;
    }

    // Default serialization for other endpoints
    return `${endpointName}(${JSON.stringify(queryArgs)})`;
  },
});

// Export hooks for usage in functional components
export const {
  util: { getRunningQueriesThunk, getRunningMutationsThunk },
} = api;