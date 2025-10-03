import { api } from './api';
import type {
  RefreshTokenRequest,
  RefreshTokenResponse,
  User,
  ApiResponse
} from '../types/api';

// OTP-specific types for the actual backend API
export interface OTPRequest {
  phone_number: string;
  country_code: string;
}

export interface OTPRequestResponse {
  success: boolean;
  message: string;
  request_id: string;
  expires_in: number;
}

export interface OTPVerifyRequest {
  request_id: string;
  code: string;
  phone_number?: string;
}

export interface OTPVerifyResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

/**
 * Authentication API endpoints using RTK Query
 *
 * Provides endpoints for user authentication including login, logout,
 * token refresh, and current user retrieval. Handles automatic token
 * management and cache invalidation.
 *
 * @example
 * ```typescript
 * // Login user
 * const [login] = useLoginMutation();
 * await login({ email: 'user@example.com', password: 'password' });
 *
 * // Get current user
 * const { data: user } = useGetCurrentUserQuery();
 *
 * // Logout
 * const [logout] = useLogoutMutation();
 * await logout();
 * ```
 */
export const authApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Request OTP for phone number authentication
     *
     * @param request OTP request with phone number
     * @param request.phoneNumber User's phone number in international format
     * @returns Promise<OTPRequestResponse> Success confirmation
     *
     * @example
     * ```typescript
     * const [requestOTP, { isLoading, error }] = useRequestOTPMutation();
     *
     * const handleRequestOTP = async () => {
     *   try {
     *     const result = await requestOTP({
     *       phoneNumber: '+66812345678'
     *     }).unwrap();
     *
     *     console.log('OTP sent:', result.message);
     *   } catch (error) {
     *     console.error('OTP request failed:', error);
     *   }
     * };
     * ```
     *
     * @throws {ApiError} When phone number is invalid or server error occurs
     */
    requestOTP: builder.mutation<OTPRequestResponse, OTPRequest>({
      query: (request) => ({
        url: '/auth/login',
        method: 'POST',
        body: request,
      }),
      transformResponse: (response: any) => response.data || response,
    }),

    /**
     * Verify OTP code and complete authentication
     *
     * @param request OTP verification request
     * @param request.phoneNumber User's phone number
     * @param request.code 6-digit OTP code
     * @returns Promise<OTPVerifyResponse> User data and authentication tokens
     *
     * @example
     * ```typescript
     * const [verifyOTP, { isLoading, error }] = useVerifyOTPMutation();
     *
     * const handleVerifyOTP = async () => {
     *   try {
     *     const result = await verifyOTP({
     *       phoneNumber: '+66812345678',
     *       code: '123456'
     *     }).unwrap();
     *
     *     console.log('Logged in user:', result.user);
     *     // Tokens are automatically stored by middleware
     *   } catch (error) {
     *     console.error('OTP verification failed:', error);
     *   }
     * };
     * ```
     *
     * @throws {ApiError} When OTP is invalid or expired
     */
    verifyOTP: builder.mutation<OTPVerifyResponse, OTPVerifyRequest>({
      query: (request) => ({
        url: '/auth/verify-otp',
        method: 'POST',
        body: request,
      }),
      invalidatesTags: ['User', 'Auth'],
      transformResponse: (response: OTPVerifyResponse) => response,
    }),

    /**
     * Log out the current user and clear session
     *
     * Invalidates all cached data and removes authentication tokens.
     * The auth middleware automatically clears localStorage and resets
     * the API state when this mutation succeeds.
     *
     * @returns Promise<void> No return data on success
     *
     * @example
     * ```typescript
     * const [logout, { isLoading }] = useLogoutMutation();
     *
     * const handleLogout = async () => {
     *   try {
     *     await logout().unwrap();
     *     // User is automatically logged out
     *     // All cache data is cleared
     *   } catch (error) {
     *     console.error('Logout failed:', error);
     *   }
     * };
     * ```
     */
    logout: builder.mutation<void, void>({
      query: () => ({
        url: '/auth/logout',
        method: 'POST',
      }),
      invalidatesTags: ['User', 'Message', 'Chat', 'Auth'],
      onQueryStarted: async (_, { dispatch, queryFulfilled }) => {
        try {
          await queryFulfilled;
          // Clear all cached data on logout
          dispatch(api.util.resetApiState());
        } catch {}
      },
    }),

    /**
     * Refresh authentication tokens
     *
     * Uses the refresh token to obtain new access and refresh tokens.
     * This is typically called automatically by the auth middleware
     * before token expiration.
     *
     * @param request Refresh token request
     * @param request.refreshToken Current valid refresh token
     * @returns Promise<RefreshTokenResponse> New tokens with expiration info
     *
     * @example
     * ```typescript
     * const [refreshToken] = useRefreshTokenMutation();
     *
     * const handleRefresh = async () => {
     *   try {
     *     const result = await refreshToken({
     *       refreshToken: 'current-refresh-token'
     *     }).unwrap();
     *
     *     // New tokens are automatically stored by middleware
     *     console.log('Tokens refreshed');
     *   } catch (error) {
     *     // Refresh failed - user will be logged out automatically
     *     console.error('Token refresh failed:', error);
     *   }
     * };
     * ```
     *
     * @throws {ApiError} When refresh token is invalid or expired
     */
    refreshToken: builder.mutation<RefreshTokenResponse, RefreshTokenRequest>({
      query: (request) => ({
        url: '/auth/refresh',
        method: 'POST',
        body: request,
      }),
    }),

    /**
     * Get current authenticated user information
     *
     * Retrieves the current user's profile data. Requires valid
     * authentication token. Data is cached and automatically
     * refreshed on login/logout.
     *
     * @returns Promise<User> Current user profile data
     *
     * @example
     * ```typescript
     * const {
     *   data: user,
     *   isLoading,
     *   error,
     *   refetch
     * } = useGetCurrentUserQuery();
     *
     * if (isLoading) return <div>Loading...</div>;
     * if (error) return <div>Please log in</div>;
     *
     * return (
     *   <div>
     *     <h1>Welcome, {user.displayName}!</h1>
     *     <p>Email: {user.email}</p>
     *     <button onClick={refetch}>Refresh Profile</button>
     *   </div>
     * );
     * ```
     *
     * @throws {ApiError} When user is not authenticated or token is invalid
     */
    getCurrentUser: builder.query<User, void>({
      query: () => '/auth/me',
      providesTags: ['User', 'Auth'],
      transformResponse: (response: ApiResponse<User>) => response.data,
    }),
  }),
});

/**
 * Generated RTK Query hooks for authentication endpoints
 *
 * These hooks provide React integration for all authentication operations
 * with automatic loading states, error handling, and cache management.
 */
export const {
  /**
   * Hook for requesting OTP mutation
   * @returns [mutationTrigger, mutationResult] OTP request mutation trigger and result
   */
  useRequestOTPMutation,

  /**
   * Hook for verifying OTP mutation
   * @returns [mutationTrigger, mutationResult] OTP verification mutation trigger and result
   */
  useVerifyOTPMutation,

  /**
   * Hook for user logout mutation
   * @returns [mutationTrigger, mutationResult] Logout mutation trigger and result
   */
  useLogoutMutation,

  /**
   * Hook for token refresh mutation
   * @returns [mutationTrigger, mutationResult] Refresh token mutation trigger and result
   */
  useRefreshTokenMutation,

  /**
   * Hook for fetching current user data
   * @returns QueryResult Current user query result with data, loading, and error states
   */
  useGetCurrentUserQuery,

  /**
   * Hook for lazy fetching of current user data
   * @returns [trigger, QueryResult] Lazy query trigger and result
   */
  useLazyGetCurrentUserQuery,
} = authApi;