import { createApi, fetchBaseQuery, BaseQueryFn, FetchArgs, FetchBaseQueryError, retry } from '@reduxjs/toolkit/query/react';
import type { RootState } from '../store';
import { RefreshTokenResponse } from '../types/api';
import { buildServiceUrl, getServiceConfig } from './serviceConfig';
import { fallbackDataService } from './fallbackData';

// Service-aware base query that routes to appropriate microservices
const createServiceAwareBaseQuery = () => {
  const serviceConfig = getServiceConfig();

  return fetchBaseQuery({
    baseUrl: '', // Will be set dynamically per request
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
};

// Enhanced base query with service routing
const baseQuery: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions
) => {
  const serviceConfig = getServiceConfig();
  const baseQueryFn = createServiceAwareBaseQuery();

  // Handle URL routing for service-aware requests
  if (typeof args === 'string') {
    const fullUrl = serviceConfig.useDirect ? buildServiceUrl(args) : `${import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'}${args}`;
    return baseQueryFn(fullUrl, api, extraOptions);
  } else {
    // For FetchArgs objects
    const endpoint = args.url;
    const fullUrl = serviceConfig.useDirect ? buildServiceUrl(endpoint) : `${import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'}${endpoint}`;

    return baseQueryFn({
      ...args,
      url: fullUrl
    }, api, extraOptions);
  }
};

// Enhanced retry logic with exponential backoff
const baseQueryWithRetry = retry(
  async (args, api, extraOptions) => {
    const result = await baseQuery(args, api, extraOptions);

    // Only retry on network errors, timeouts, and 5xx server errors
    // Don't retry on CORS, 4xx client errors, or parsing errors
    if (result.error) {
      const { status } = result.error as FetchBaseQueryError;
      const errorDetails = result.error as FetchBaseQueryError & { error?: string };
      const errorMessage = typeof errorDetails.error === 'string' ? errorDetails.error : '';
      const isLikelyCorsError = status === 'FETCH_ERROR' && /cors|failed to fetch/i.test(errorMessage);

      // Log CORS and client errors but don't retry
      if (
        status === 'PARSING_ERROR' ||
        status === 'CUSTOM_ERROR' ||
        (typeof status === 'number' && status >= 400 && status < 500)
      ) {
        console.warn(`API error (${status}) - not retrying:`, result.error);
        return result; // Return error result without retrying
      }

      if (isLikelyCorsError) {
        console.warn('CORS or fetch error detected - not retrying:', result.error);
        return result;
      }

      // Only retry on network/server errors
      if (
        status === 'FETCH_ERROR' ||
        status === 'TIMEOUT_ERROR' ||
        (typeof status === 'number' && status >= 500)
      ) {
        throw result.error; // This triggers retry
      }
    }

    return result;
  },
  {
    maxRetries: 2, // Reduced retries
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

  // Handle fallback data for development when endpoints are not available
  if (result.error && import.meta.env.DEV) {
    const endpoint = typeof args === 'string' ? args : args.url;
    const method = typeof args === 'string' ? 'GET' : (args.method || 'GET');

    // Only use fallback for GET requests to avoid unintended side effects
    if (method === 'GET' && fallbackDataService.shouldUseFallback(endpoint, result.error)) {
      console.log(`🔄 API endpoint ${endpoint} not available, using fallback data`);

      const fallbackResponse = fallbackDataService.createFallbackResponse(endpoint);

      // Return fallback data with success status
      return {
        data: fallbackResponse,
        meta: {
          request: args,
          response: { status: 200 },
          fallback: true
        }
      };
    }
  }

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
    'ContentVersion',
    'Video',
    'Channel',
    'Comment',
    'Playlist',
    'LiveStream',
    'Search',
    'Analytics',
    // Social service tags
    'SocialProfile',
    'SocialPost',
    'SocialComment',
    'SocialReaction',
    'SocialFeed',
    'SocialTrending',
    'SocialStories',
    'SocialFriends',
    'SocialFollowers',
    'SocialFollowing',
    'SocialAnalytics',
    'UserRelationship'
  ],
  endpoints: () => ({}),
  refetchOnMountOrArgChange: false, // Disable automatic refetch
  refetchOnFocus: false, // Disable refetch on focus
  refetchOnReconnect: false, // Disable refetch on reconnect
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
