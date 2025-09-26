# Enterprise API Integration Guide

**Comprehensive API Integration Documentation for Production Deployment**
- **Coverage**: 12 RTK Query endpoints, authentication, state management, error handling
- **Platforms**: Web (React/RTK Query), iOS (Swift/Alamofire), Android (Kotlin/Retrofit)
- **Enterprise Features**: Advanced caching, optimistic updates, error recovery, monitoring integration

---

## 1. API Architecture Overview

### 1.1 Enterprise API Architecture

The Tchat API integration follows enterprise-grade patterns for scalability, reliability, and performance:

```typescript
interface EnterpriseAPIArchitecture {
  coreServices: {
    contentManagement: '12 comprehensive endpoints with versioning';
    authentication: 'JWT-based with automatic refresh';
    userManagement: 'Profile and preferences synchronization';
    messaging: 'Real-time messaging with offline support';
    commerce: 'E-commerce integration with payment processing';
  };
  crossCuttingConcerns: {
    caching: 'Advanced tag-based invalidation with Redis backend';
    errorHandling: 'Circuit breakers with exponential backoff';
    monitoring: 'APM integration with performance tracking';
    security: 'OAuth 2.0 with PKCE, rate limiting, CORS management';
    performance: 'Request deduplication, prefetching, background sync';
  };
  dataFlow: {
    realTimeSync: 'WebSocket-based state synchronization';
    offlineSupport: 'localStorage fallback with conflict resolution';
    optimisticUpdates: 'Immediate UI updates with rollback capability';
    backgroundSync: 'Service worker integration for offline-first experience';
  };
}
```

### 1.2 API Service Layer Structure

```typescript
// Core API configuration
export const api = createApi({
  reducerPath: 'api',
  baseQuery: baseQueryWithReauth,
  tagTypes: [
    'Content',
    'ContentCategory',
    'ContentVersion',
    'User',
    'UserProfile',
    'Message',
    'Chat',
    'Notification',
    'Commerce',
    'Payment'
  ],
  endpoints: (builder) => ({
    // Endpoints defined in individual service files
  }),
  keepUnusedDataFor: 60, // 1 minute default cache
  refetchOnMountOrArgChange: 30, // 30 seconds
  refetchOnFocus: true,
  refetchOnReconnect: true
});
```

---

## 2. Content Management API Integration

### 2.1 Comprehensive Content API Endpoints

The content management system provides 12 sophisticated endpoints for enterprise content operations:

#### Core Content Endpoints

```typescript
export const contentApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // 1. Get Content Items (with pagination, filtering, sorting)
    getContentItems: builder.query<
      PaginatedResponse<ContentItem[]>,
      GetContentItemsRequest
    >({
      query: ({ page = 1, limit = 20, category, status, sortBy = 'updatedAt', sortOrder = 'desc', search }) => ({
        url: '/content/items',
        params: {
          page,
          limit,
          ...(category && { category }),
          ...(status && { status }),
          sortBy,
          sortOrder,
          ...(search && { search })
        }
      }),
      providesTags: (result) => [
        { type: 'Content', id: 'LIST' },
        ...(result?.data || []).map(({ id }) => ({ type: 'Content' as const, id }))
      ],
      transformResponse: (response: ApiResponse<PaginatedResponse<ContentItem[]>>) => response.data,
      transformErrorResponse: (response) => ({
        status: response.status,
        message: response.data?.message || 'Failed to fetch content items',
        code: response.data?.code
      })
    }),

    // 2. Get Single Content Item (with version support)
    getContentItem: builder.query<ContentItem, { id: string; version?: string }>({
      query: ({ id, version }) => ({
        url: `/content/items/${id}`,
        params: version ? { version } : {}
      }),
      providesTags: (_result, _error, { id }) => [{ type: 'Content', id }],
      transformResponse: (response: ApiResponse<ContentItem>) => response.data,
      keepCacheDataFor: 300 // 5 minutes for individual items
    }),

    // 3. Get Content by Category (with advanced filtering)
    getContentByCategory: builder.query<
      PaginatedResponse<ContentItem[]>,
      GetContentByCategoryRequest
    >({
      query: ({
        category,
        page = 1,
        limit = 20,
        includeSubcategories = false,
        status = 'published',
        dateRange
      }) => ({
        url: `/content/categories/${category}/items`,
        params: {
          page,
          limit,
          includeSubcategories,
          status,
          ...(dateRange && {
            startDate: dateRange.startDate,
            endDate: dateRange.endDate
          })
        }
      }),
      providesTags: (result, _error, { category }) => [
        { type: 'Content', id: `CATEGORY_${category}` },
        { type: 'ContentCategory', id: category },
        ...(result?.data || []).map(({ id }) => ({ type: 'Content' as const, id }))
      ],
      transformResponse: (response: ApiResponse<PaginatedResponse<ContentItem[]>>) => response.data
    }),

    // 4. Get Content Categories (hierarchical structure)
    getContentCategories: builder.query<ContentCategory[], { includeEmpty?: boolean }>({
      query: ({ includeEmpty = false }) => ({
        url: '/content/categories',
        params: { includeEmpty }
      }),
      providesTags: (result) => [
        { type: 'ContentCategory', id: 'LIST' },
        ...(result || []).map(({ id }) => ({ type: 'ContentCategory' as const, id }))
      ],
      transformResponse: (response: ApiResponse<ContentCategory[]>) => response.data,
      keepCacheDataFor: 600 // 10 minutes for categories
    }),

    // 5. Get Content Versions (version history)
    getContentVersions: builder.query<ContentVersion[], { contentId: string; limit?: number }>({
      query: ({ contentId, limit = 10 }) => ({
        url: `/content/items/${contentId}/versions`,
        params: { limit }
      }),
      providesTags: (_result, _error, { contentId }) => [
        { type: 'ContentVersion', id: `CONTENT_${contentId}` }
      ],
      transformResponse: (response: ApiResponse<ContentVersion[]>) => response.data
    }),

    // 6. Sync Content (incremental synchronization)
    syncContent: builder.query<ContentSyncResponse, ContentSyncRequest>({
      query: ({ lastSyncTimestamp, includeDeleted = false, platforms = ['web'] }) => ({
        url: '/content/sync',
        params: {
          lastSyncTimestamp,
          includeDeleted,
          platforms: platforms.join(',')
        }
      }),
      providesTags: ['Content'],
      transformResponse: (response: ApiResponse<ContentSyncResponse>) => response.data,
      keepCacheDataFor: 0 // No caching for sync operations
    }),

    // 7. Create Content Item (with draft support)
    createContentItem: builder.mutation<ContentItem, CreateContentItemRequest>({
      query: (contentData) => ({
        url: '/content/items',
        method: 'POST',
        body: contentData
      }),
      invalidatesTags: (result) => [
        { type: 'Content', id: 'LIST' },
        ...(result?.category ? [{ type: 'Content', id: `CATEGORY_${result.category}` }] : [])
      ],
      transformResponse: (response: ApiResponse<ContentItem>) => response.data,
      onQueryStarted: async (contentData, { dispatch, queryFulfilled }) => {
        // Optimistic update
        const patchResult = dispatch(
          contentApi.util.updateQueryData('getContentItems', { page: 1 }, (draft) => {
            const optimisticItem: ContentItem = {
              id: `temp-${Date.now()}`,
              ...contentData,
              status: 'draft',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
              version: 1
            };
            draft.data.unshift(optimisticItem);
          })
        );

        try {
          const { data: createdItem } = await queryFulfilled;
          // Update optimistic item with real data
          dispatch(
            contentApi.util.updateQueryData('getContentItems', { page: 1 }, (draft) => {
              const index = draft.data.findIndex(item => item.id.startsWith('temp-'));
              if (index !== -1) {
                draft.data[index] = createdItem;
              }
            })
          );
        } catch {
          // Rollback optimistic update
          patchResult.undo();
        }
      }
    }),

    // 8. Update Content Item (with conflict resolution)
    updateContentItem: builder.mutation<
      ContentItem,
      UpdateContentItemRequest & { id: string }
    >({
      query: ({ id, ...updateData }) => ({
        url: `/content/items/${id}`,
        method: 'PUT',
        body: updateData
      }),
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'Content', id },
        { type: 'Content', id: 'LIST' }
      ],
      transformResponse: (response: ApiResponse<ContentItem>) => response.data,
      onQueryStarted: async ({ id, ...updateData }, { dispatch, queryFulfilled }) => {
        // Optimistic update for individual item
        const patchResult = dispatch(
          contentApi.util.updateQueryData('getContentItem', { id }, (draft) => {
            Object.assign(draft, updateData, { updatedAt: new Date().toISOString() });
          })
        );

        try {
          await queryFulfilled;
        } catch (error) {
          patchResult.undo();

          // Handle conflict resolution
          if (error.error?.status === 409) {
            dispatch(showNotification({
              type: 'warning',
              title: 'Update Conflict',
              message: 'Content was modified by another user. Please refresh and try again.',
              action: {
                label: 'Refresh',
                handler: () => dispatch(contentApi.util.invalidateTags([{ type: 'Content', id }]))
              }
            }));
          }
        }
      }
    }),

    // 9. Publish Content (workflow management)
    publishContent: builder.mutation<ContentItem, { id: string; publishOptions?: PublishOptions }>({
      query: ({ id, publishOptions }) => ({
        url: `/content/items/${id}/publish`,
        method: 'POST',
        body: publishOptions || {}
      }),
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'Content', id },
        { type: 'Content', id: 'LIST' },
        { type: 'Content', id: 'PUBLISHED' }
      ],
      transformResponse: (response: ApiResponse<ContentItem>) => response.data
    }),

    // 10. Archive Content (soft delete with recovery)
    archiveContent: builder.mutation<{ success: boolean }, { id: string; archiveReason?: string }>({
      query: ({ id, archiveReason }) => ({
        url: `/content/items/${id}/archive`,
        method: 'POST',
        body: archiveReason ? { reason: archiveReason } : {}
      }),
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'Content', id },
        { type: 'Content', id: 'LIST' }
      ]
    }),

    // 11. Bulk Update Content (batch operations)
    bulkUpdateContent: builder.mutation<
      BulkUpdateResponse,
      { operations: BulkUpdateOperation[] }
    >({
      query: ({ operations }) => ({
        url: '/content/items/bulk',
        method: 'POST',
        body: { operations }
      }),
      invalidatesTags: ['Content'],
      transformResponse: (response: ApiResponse<BulkUpdateResponse>) => response.data,
      onQueryStarted: async ({ operations }, { dispatch, queryFulfilled }) => {
        // Show progress notification
        dispatch(showNotification({
          type: 'info',
          title: 'Bulk Update',
          message: `Processing ${operations.length} operations...`,
          persistent: true
        }));

        try {
          const result = await queryFulfilled;
          dispatch(hideNotification());

          if (result.data.failedOperations.length > 0) {
            dispatch(showNotification({
              type: 'warning',
              title: 'Bulk Update Completed with Errors',
              message: `${result.data.successfulOperations} succeeded, ${result.data.failedOperations.length} failed`,
              details: result.data.failedOperations.map(op => op.error).join('\n')
            }));
          } else {
            dispatch(showNotification({
              type: 'success',
              title: 'Bulk Update Completed',
              message: `Successfully processed ${result.data.successfulOperations} operations`
            }));
          }
        } catch (error) {
          dispatch(hideNotification());
          dispatch(showNotification({
            type: 'error',
            title: 'Bulk Update Failed',
            message: 'The bulk update operation encountered an error'
          }));
        }
      }
    }),

    // 12. Revert Content Version (version control)
    revertContentVersion: builder.mutation<
      ContentItem,
      { contentId: string; versionId: string; createBackup?: boolean }
    >({
      query: ({ contentId, versionId, createBackup = true }) => ({
        url: `/content/items/${contentId}/revert`,
        method: 'POST',
        body: { versionId, createBackup }
      }),
      invalidatesTags: (_result, _error, { contentId }) => [
        { type: 'Content', id: contentId },
        { type: 'ContentVersion', id: `CONTENT_${contentId}` }
      ],
      transformResponse: (response: ApiResponse<ContentItem>) => response.data
    })
  })
});

// Export hooks for component usage
export const {
  useGetContentItemsQuery,
  useGetContentItemQuery,
  useGetContentByCategoryQuery,
  useGetContentCategoriesQuery,
  useGetContentVersionsQuery,
  useSyncContentQuery,
  useCreateContentItemMutation,
  useUpdateContentItemMutation,
  usePublishContentMutation,
  useArchiveContentMutation,
  useBulkUpdateContentMutation,
  useRevertContentVersionMutation
} = contentApi;
```

### 2.2 Advanced Caching Strategy

#### Tag-Based Cache Invalidation System

```typescript
// Advanced cache invalidation with relationships
export const cacheInvalidationConfig = {
  // Content-related invalidations
  contentOperations: {
    create: (result: ContentItem) => [
      { type: 'Content', id: 'LIST' },
      { type: 'Content', id: `CATEGORY_${result.category}` },
      { type: 'ContentCategory', id: result.category }
    ],
    update: (result: ContentItem, originalData?: ContentItem) => [
      { type: 'Content', id: result.id },
      { type: 'Content', id: 'LIST' },
      ...(originalData && originalData.category !== result.category ? [
        { type: 'Content', id: `CATEGORY_${originalData.category}` },
        { type: 'Content', id: `CATEGORY_${result.category}` }
      ] : [{ type: 'Content', id: `CATEGORY_${result.category}` }])
    ],
    delete: (contentId: string, category?: string) => [
      { type: 'Content', id: contentId },
      { type: 'Content', id: 'LIST' },
      ...(category ? [{ type: 'Content', id: `CATEGORY_${category}` }] : [])
    ],
    publish: (result: ContentItem) => [
      { type: 'Content', id: result.id },
      { type: 'Content', id: 'LIST' },
      { type: 'Content', id: 'PUBLISHED' },
      { type: 'Content', id: `CATEGORY_${result.category}` }
    ]
  },

  // Cross-cutting invalidations
  userOperations: {
    profileUpdate: () => [
      { type: 'User', id: 'PROFILE' },
      { type: 'Content', id: 'USER_CONTENT' }
    ],
    preferencesUpdate: () => [
      { type: 'User', id: 'PREFERENCES' },
      { type: 'Content', id: 'LIST' } // May affect content filtering
    ]
  }
};

// Cache warming strategy
export const cacheWarmingService = {
  async warmEssentialData(userId: string) {
    const dispatch = store.dispatch;

    // Pre-fetch critical data
    const warmingPromises = [
      // Essential content
      dispatch(contentApi.endpoints.getContentItems.initiate({ page: 1, limit: 20 })),
      dispatch(contentApi.endpoints.getContentCategories.initiate({ includeEmpty: false })),

      // User-specific data
      dispatch(userApi.endpoints.getUserProfile.initiate({ userId })),
      dispatch(userApi.endpoints.getUserPreferences.initiate()),

      // Recent activity
      dispatch(contentApi.endpoints.getContentItems.initiate({
        page: 1,
        limit: 10,
        sortBy: 'updatedAt',
        sortOrder: 'desc'
      }))
    ];

    await Promise.allSettled(warmingPromises);
    console.log('✅ Cache warming completed');
  },

  async warmCategoryData(category: string) {
    const dispatch = store.dispatch;

    await dispatch(contentApi.endpoints.getContentByCategory.initiate({
      category,
      page: 1,
      limit: 20
    }));

    console.log(`✅ Category cache warmed for: ${category}`);
  }
};
```

### 2.3 Error Recovery and Resilience

#### Advanced Error Handling System

```typescript
// Enhanced error recovery system
export const apiErrorRecoveryService = {
  async handleNetworkError(error: NetworkError, originalRequest: ApiRequest): Promise<ApiResponse | null> {
    const maxRetries = 3;
    const baseDelay = 1000; // 1 second

    for (let attempt = 1; attempt <= maxRetries; attempt++) {
      const delay = baseDelay * Math.pow(2, attempt - 1); // Exponential backoff

      console.log(`🔄 Retry attempt ${attempt}/${maxRetries} after ${delay}ms`);
      await new Promise(resolve => setTimeout(resolve, delay));

      try {
        const retryResponse = await this.retryRequest(originalRequest);
        console.log(`✅ Request succeeded on attempt ${attempt}`);
        return retryResponse;
      } catch (retryError) {
        if (attempt === maxRetries) {
          console.error('❌ All retry attempts exhausted');
          throw retryError;
        }
      }
    }

    return null;
  },

  async handleServerError(error: ServerError, fallbackStrategy: 'cache' | 'offline' | 'degraded'): Promise<ApiResponse | null> {
    switch (fallbackStrategy) {
      case 'cache':
        return this.attemptCacheFallback(error.endpoint);

      case 'offline':
        return this.activateOfflineMode(error.endpoint);

      case 'degraded':
        return this.provideDegradedService(error.endpoint);

      default:
        return null;
    }
  },

  async attemptCacheFallback(endpoint: string): Promise<ApiResponse | null> {
    const cacheKey = `fallback_${endpoint}`;
    const cachedData = localStorage.getItem(cacheKey);

    if (cachedData) {
      try {
        const parsedData = JSON.parse(cachedData);
        console.log('📦 Serving cached fallback data');

        // Show user notification about cached data
        store.dispatch(showNotification({
          type: 'warning',
          title: 'Using Cached Data',
          message: 'Showing cached content due to connection issues',
          action: {
            label: 'Retry',
            handler: () => this.retryOriginalRequest(endpoint)
          }
        }));

        return parsedData;
      } catch (parseError) {
        console.error('Failed to parse cached fallback data:', parseError);
      }
    }

    return null;
  },

  async activateOfflineMode(endpoint: string): Promise<ApiResponse | null> {
    // Activate offline mode
    store.dispatch(setOfflineMode(true));

    // Use IndexedDB or localStorage for offline data
    const offlineData = await this.getOfflineData(endpoint);

    if (offlineData) {
      store.dispatch(showNotification({
        type: 'info',
        title: 'Offline Mode',
        message: 'Working offline - changes will sync when connection is restored',
        persistent: true
      }));

      return offlineData;
    }

    return null;
  },

  async provideDegradedService(endpoint: string): Promise<ApiResponse | null> {
    // Provide minimal functionality
    const degradedResponse = this.createDegradedResponse(endpoint);

    store.dispatch(showNotification({
      type: 'warning',
      title: 'Limited Functionality',
      message: 'Some features are temporarily unavailable',
      persistent: true
    }));

    return degradedResponse;
  }
};

// Circuit breaker implementation
export class ApiCircuitBreaker {
  private failures = new Map<string, number>();
  private lastFailureTime = new Map<string, number>();
  private readonly failureThreshold = 5;
  private readonly timeoutDuration = 60000; // 1 minute

  async executeRequest<T>(
    endpoint: string,
    requestFn: () => Promise<T>
  ): Promise<T> {
    if (this.isCircuitOpen(endpoint)) {
      throw new Error(`Circuit breaker is open for endpoint: ${endpoint}`);
    }

    try {
      const result = await requestFn();
      this.onSuccess(endpoint);
      return result;
    } catch (error) {
      this.onFailure(endpoint);
      throw error;
    }
  }

  private isCircuitOpen(endpoint: string): boolean {
    const failures = this.failures.get(endpoint) || 0;
    const lastFailure = this.lastFailureTime.get(endpoint) || 0;
    const timeSinceLastFailure = Date.now() - lastFailure;

    // Circuit is open if failure threshold exceeded and timeout not elapsed
    return failures >= this.failureThreshold && timeSinceLastFailure < this.timeoutDuration;
  }

  private onSuccess(endpoint: string): void {
    this.failures.delete(endpoint);
    this.lastFailureTime.delete(endpoint);
  }

  private onFailure(endpoint: string): void {
    const currentFailures = this.failures.get(endpoint) || 0;
    this.failures.set(endpoint, currentFailures + 1);
    this.lastFailureTime.set(endpoint, Date.now());

    if (currentFailures + 1 >= this.failureThreshold) {
      console.warn(`⚠️ Circuit breaker opened for endpoint: ${endpoint}`);
    }
  }
}
```

---

## 3. Authentication and Security Integration

### 3.1 JWT Token Management System

#### Enterprise Authentication Service

```typescript
// Comprehensive JWT authentication service
export const authService = {
  // Token storage with secure fallbacks
  storage: {
    async setTokens(tokens: TokenPair): Promise<void> {
      try {
        // Primary: Secure storage (mobile) or httpOnly cookies (web)
        if (this.isNativeApp()) {
          await this.setSecureStorageTokens(tokens);
        } else {
          await this.setHttpOnlyCookieTokens(tokens);
        }
      } catch (error) {
        // Fallback: Encrypted localStorage
        console.warn('Primary storage failed, using encrypted fallback:', error);
        await this.setEncryptedStorageTokens(tokens);
      }
    },

    async getTokens(): Promise<TokenPair | null> {
      try {
        if (this.isNativeApp()) {
          return await this.getSecureStorageTokens();
        } else {
          return await this.getHttpOnlyCookieTokens();
        }
      } catch (error) {
        console.warn('Primary storage read failed, trying encrypted fallback:', error);
        return await this.getEncryptedStorageTokens();
      }
    },

    async clearTokens(): Promise<void> {
      await Promise.allSettled([
        this.clearSecureStorageTokens(),
        this.clearHttpOnlyCookieTokens(),
        this.clearEncryptedStorageTokens()
      ]);
    }
  },

  // Automatic token refresh
  tokenRefresh: {
    refreshPromise: null as Promise<TokenPair> | null,
    isRefreshing: false,

    async refreshTokens(): Promise<TokenPair> {
      // Prevent concurrent refresh attempts
      if (this.refreshPromise) {
        return this.refreshPromise;
      }

      this.isRefreshing = true;
      this.refreshPromise = this.performRefresh();

      try {
        const newTokens = await this.refreshPromise;
        await authService.storage.setTokens(newTokens);

        // Update Redux store
        store.dispatch(authSlice.actions.setTokens(newTokens));

        console.log('✅ Tokens refreshed successfully');
        return newTokens;
      } catch (error) {
        console.error('❌ Token refresh failed:', error);
        // Clear invalid tokens and redirect to login
        await authService.storage.clearTokens();
        store.dispatch(authSlice.actions.logout());
        throw error;
      } finally {
        this.isRefreshing = false;
        this.refreshPromise = null;
      }
    },

    async performRefresh(): Promise<TokenPair> {
      const currentTokens = await authService.storage.getTokens();

      if (!currentTokens?.refreshToken) {
        throw new Error('No refresh token available');
      }

      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${currentTokens.refreshToken}`
        },
        credentials: 'include'
      });

      if (!response.ok) {
        throw new Error(`Token refresh failed: ${response.status}`);
      }

      const data = await response.json();
      return data.tokens;
    }
  },

  // Authentication state management
  auth: {
    async login(credentials: LoginCredentials): Promise<AuthResult> {
      try {
        const response = await fetch('/api/auth/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(credentials),
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error(`Login failed: ${response.status}`);
        }

        const authResult: AuthResult = await response.json();

        // Store tokens securely
        await authService.storage.setTokens(authResult.tokens);

        // Update Redux store
        store.dispatch(authSlice.actions.loginSuccess({
          user: authResult.user,
          tokens: authResult.tokens,
          permissions: authResult.permissions
        }));

        // Warm cache with user-specific data
        await cacheWarmingService.warmEssentialData(authResult.user.id);

        console.log('✅ Login successful');
        return authResult;

      } catch (error) {
        console.error('❌ Login failed:', error);
        store.dispatch(authSlice.actions.loginFailure(error.message));
        throw error;
      }
    },

    async logout(): Promise<void> {
      try {
        const currentTokens = await authService.storage.getTokens();

        if (currentTokens?.accessToken) {
          // Notify server about logout
          await fetch('/api/auth/logout', {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${currentTokens.accessToken}`
            },
            credentials: 'include'
          });
        }

        // Clear all stored data
        await authService.storage.clearTokens();

        // Clear Redux store
        store.dispatch(authSlice.actions.logout());

        // Clear all API cache
        store.dispatch(api.util.resetApiState());

        console.log('✅ Logout successful');

      } catch (error) {
        console.error('❌ Logout error (proceeding anyway):', error);

        // Clear local data even if server request failed
        await authService.storage.clearTokens();
        store.dispatch(authSlice.actions.logout());
      }
    },

    async verifyToken(): Promise<boolean> {
      try {
        const tokens = await authService.storage.getTokens();

        if (!tokens?.accessToken) {
          return false;
        }

        // Check if token is expired (client-side check)
        if (this.isTokenExpired(tokens.accessToken)) {
          // Attempt refresh
          try {
            await authService.tokenRefresh.refreshTokens();
            return true;
          } catch {
            return false;
          }
        }

        // Verify token with server
        const response = await fetch('/api/auth/verify', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${tokens.accessToken}`
          },
          credentials: 'include'
        });

        return response.ok;

      } catch {
        return false;
      }
    },

    isTokenExpired(token: string): boolean {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        const currentTime = Math.floor(Date.now() / 1000);
        return payload.exp < currentTime;
      } catch {
        return true; // Assume expired if parsing fails
      }
    }
  }
};

// Authentication Redux slice
export const authSlice = createSlice({
  name: 'auth',
  initialState: {
    user: null as User | null,
    tokens: null as TokenPair | null,
    permissions: [] as string[],
    isAuthenticated: false,
    isLoading: false,
    error: null as string | null,
    lastActivity: null as string | null,
    sessionTimeout: 24 * 60 * 60 * 1000 // 24 hours
  },
  reducers: {
    loginStart: (state) => {
      state.isLoading = true;
      state.error = null;
    },
    loginSuccess: (state, action) => {
      state.isLoading = false;
      state.isAuthenticated = true;
      state.user = action.payload.user;
      state.tokens = action.payload.tokens;
      state.permissions = action.payload.permissions;
      state.lastActivity = new Date().toISOString();
      state.error = null;
    },
    loginFailure: (state, action) => {
      state.isLoading = false;
      state.isAuthenticated = false;
      state.error = action.payload;
      state.user = null;
      state.tokens = null;
      state.permissions = [];
    },
    logout: (state) => {
      state.isAuthenticated = false;
      state.user = null;
      state.tokens = null;
      state.permissions = [];
      state.error = null;
      state.lastActivity = null;
    },
    setTokens: (state, action) => {
      state.tokens = action.payload;
      state.lastActivity = new Date().toISOString();
    },
    updateLastActivity: (state) => {
      state.lastActivity = new Date().toISOString();
    }
  }
});
```

### 3.2 Cross-Platform Authentication Integration

#### iOS Authentication Integration

```swift
// iOS JWT Authentication Service
import Foundation
import Security
import Combine

class iOSAuthenticationService: ObservableObject {
    @Published var isAuthenticated = false
    @Published var user: User?
    @Published var isLoading = false

    private let baseURL = "https://api.tchat.com"
    private let keychainService = "com.tchat.tokens"

    // MARK: - Authentication Methods

    func login(credentials: LoginCredentials) async throws -> AuthResult {
        isLoading = true
        defer { isLoading = false }

        let loginRequest = LoginRequest(
            email: credentials.email,
            password: credentials.password,
            deviceId: UIDevice.current.identifierForVendor?.uuidString ?? UUID().uuidString,
            platform: "ios"
        )

        let url = URL(string: "\(baseURL)/auth/login")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = try JSONEncoder().encode(loginRequest)

        let (data, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw AuthError.loginFailed
        }

        let authResult = try JSONDecoder().decode(AuthResult.self, from: data)

        // Store tokens securely
        try await storeTokensSecurely(authResult.tokens)

        // Update state
        await MainActor.run {
            self.user = authResult.user
            self.isAuthenticated = true
        }

        // Warm cache
        await warmUserCache(userId: authResult.user.id)

        return authResult
    }

    func logout() async {
        // Clear tokens from keychain
        await clearStoredTokens()

        // Clear state
        await MainActor.run {
            self.user = nil
            self.isAuthenticated = false
        }

        // Clear all cached data
        URLCache.shared.removeAllCachedResponses()
    }

    // MARK: - Secure Token Storage

    private func storeTokensSecurely(_ tokens: TokenPair) async throws {
        let tokenData = try JSONEncoder().encode(tokens)

        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: keychainService,
            kSecAttrAccount as String: "auth_tokens",
            kSecValueData as String: tokenData
        ]

        // Delete existing item first
        SecItemDelete(query as CFDictionary)

        // Add new item
        let status = SecItemAdd(query as CFDictionary, nil)

        guard status == errSecSuccess else {
            throw AuthError.tokenStorageFailed
        }
    }

    private func retrieveStoredTokens() async throws -> TokenPair? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: keychainService,
            kSecAttrAccount as String: "auth_tokens",
            kSecReturnData as String: true
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        guard status == errSecSuccess,
              let tokenData = result as? Data else {
            return nil
        }

        return try JSONDecoder().decode(TokenPair.self, from: tokenData)
    }

    private func clearStoredTokens() async {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: keychainService,
            kSecAttrAccount as String: "auth_tokens"
        ]

        SecItemDelete(query as CFDictionary)
    }

    // MARK: - Token Refresh

    func refreshTokensIfNeeded() async throws {
        guard let tokens = try await retrieveStoredTokens() else {
            throw AuthError.noTokensFound
        }

        // Check if access token is expired
        if isTokenExpired(tokens.accessToken) {
            let newTokens = try await refreshTokens(refreshToken: tokens.refreshToken)
            try await storeTokensSecurely(newTokens)
        }
    }

    private func refreshTokens(refreshToken: String) async throws -> TokenPair {
        let url = URL(string: "\(baseURL)/auth/refresh")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(refreshToken)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse,
              httpResponse.statusCode == 200 else {
            throw AuthError.tokenRefreshFailed
        }

        let refreshResult = try JSONDecoder().decode(TokenRefreshResult.self, from: data)
        return refreshResult.tokens
    }

    private func isTokenExpired(_ token: String) -> Bool {
        let parts = token.components(separatedBy: ".")
        guard parts.count == 3,
              let payloadData = Data(base64Encoded: parts[1].padding(to: 4)) else {
            return true
        }

        do {
            let payload = try JSONSerialization.jsonObject(with: payloadData) as? [String: Any]
            let exp = payload?["exp"] as? TimeInterval ?? 0
            return Date().timeIntervalSince1970 > exp
        } catch {
            return true
        }
    }
}
```

#### Android Authentication Integration

```kotlin
// Android JWT Authentication Service
class AndroidAuthenticationService @Inject constructor(
    private val apiService: ApiService,
    private val tokenStorage: SecureTokenStorage,
    private val userPreferences: UserPreferences
) {
    private val _authState = MutableStateFlow(AuthState())
    val authState: StateFlow<AuthState> = _authState.asStateFlow()

    suspend fun login(credentials: LoginCredentials): Result<AuthResult> {
        return try {
            _authState.value = _authState.value.copy(isLoading = true)

            val loginRequest = LoginRequest(
                email = credentials.email,
                password = credentials.password,
                deviceId = getDeviceId(),
                platform = "android"
            )

            val authResult = apiService.login(loginRequest)

            // Store tokens securely
            tokenStorage.storeTokens(authResult.tokens)

            // Update auth state
            _authState.value = AuthState(
                isAuthenticated = true,
                user = authResult.user,
                permissions = authResult.permissions,
                isLoading = false
            )

            // Warm cache
            warmUserCache(authResult.user.id)

            Result.success(authResult)

        } catch (e: Exception) {
            _authState.value = _authState.value.copy(
                isLoading = false,
                error = e.message
            )
            Result.failure(e)
        }
    }

    suspend fun logout() {
        try {
            // Notify server
            val tokens = tokenStorage.getTokens()
            if (tokens != null) {
                apiService.logout(tokens.accessToken)
            }
        } catch (e: Exception) {
            // Continue with logout even if server request fails
            Log.w("AuthService", "Server logout failed", e)
        }

        // Clear local data
        tokenStorage.clearTokens()
        userPreferences.clearUserData()

        // Update state
        _authState.value = AuthState()
    }

    suspend fun refreshTokensIfNeeded(): Boolean {
        val tokens = tokenStorage.getTokens() ?: return false

        if (isTokenExpired(tokens.accessToken)) {
            return try {
                val newTokens = refreshTokens(tokens.refreshToken)
                tokenStorage.storeTokens(newTokens)
                true
            } catch (e: Exception) {
                Log.e("AuthService", "Token refresh failed", e)
                logout() // Clear invalid tokens
                false
            }
        }

        return true
    }

    private suspend fun refreshTokens(refreshToken: String): TokenPair {
        return apiService.refreshTokens(refreshToken)
    }

    private fun isTokenExpired(token: String): Boolean {
        return try {
            val parts = token.split(".")
            if (parts.size != 3) return true

            val payload = Base64.decode(parts[1], Base64.DEFAULT)
            val jsonObject = JSONObject(String(payload))
            val exp = jsonObject.getLong("exp")

            System.currentTimeMillis() / 1000 > exp
        } catch (e: Exception) {
            true
        }
    }

    private fun getDeviceId(): String {
        return Settings.Secure.getString(
            context.contentResolver,
            Settings.Secure.ANDROID_ID
        )
    }
}

// Secure token storage implementation
@Singleton
class SecureTokenStorage @Inject constructor(
    private val context: Context
) {
    private val sharedPreferences = EncryptedSharedPreferences.create(
        "auth_tokens",
        MasterKey.Builder(context)
            .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
            .build(),
        context,
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )

    fun storeTokens(tokens: TokenPair) {
        val tokensJson = Gson().toJson(tokens)
        sharedPreferences.edit()
            .putString("tokens", tokensJson)
            .apply()
    }

    fun getTokens(): TokenPair? {
        val tokensJson = sharedPreferences.getString("tokens", null)
        return if (tokensJson != null) {
            try {
                Gson().fromJson(tokensJson, TokenPair::class.java)
            } catch (e: Exception) {
                null
            }
        } else {
            null
        }
    }

    fun clearTokens() {
        sharedPreferences.edit()
            .remove("tokens")
            .apply()
    }
}
```

---

## 4. Real-Time Features and WebSocket Integration

### 4.1 WebSocket Connection Management

#### Enterprise WebSocket Service

```typescript
// Comprehensive WebSocket service for real-time features
export class EnterpriseWebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;
  private reconnectDelay = 1000; // Start with 1 second
  private heartbeatInterval: NodeJS.Timeout | null = null;
  private messageQueue: WebSocketMessage[] = [];
  private subscribers = new Map<string, Set<(data: any) => void>>();

  constructor(private baseUrl: string) {
    this.setupConnectionMonitoring();
  }

  async connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        const tokens = store.getState().auth.tokens;
        const wsUrl = `${this.baseUrl}?token=${tokens?.accessToken}`;

        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
          console.log('🔗 WebSocket connected');
          this.reconnectAttempts = 0;
          this.reconnectDelay = 1000;
          this.startHeartbeat();
          this.flushMessageQueue();
          resolve();
        };

        this.ws.onmessage = (event) => {
          this.handleMessage(JSON.parse(event.data));
        };

        this.ws.onclose = (event) => {
          console.log('🔌 WebSocket disconnected', event.code);
          this.stopHeartbeat();

          if (!event.wasClean && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.scheduleReconnect();
          }
        };

        this.ws.onerror = (error) => {
          console.error('❌ WebSocket error:', error);
          reject(error);
        };

        // Connection timeout
        setTimeout(() => {
          if (this.ws?.readyState !== WebSocket.OPEN) {
            this.ws?.close();
            reject(new Error('WebSocket connection timeout'));
          }
        }, 5000);

      } catch (error) {
        reject(error);
      }
    });
  }

  disconnect(): void {
    this.stopHeartbeat();
    this.ws?.close(1000, 'Client disconnect');
    this.ws = null;
    this.subscribers.clear();
    this.messageQueue = [];
  }

  send(message: WebSocketMessage): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      // Queue message for later
      this.messageQueue.push(message);
    }
  }

  subscribe(event: string, callback: (data: any) => void): () => void {
    if (!this.subscribers.has(event)) {
      this.subscribers.set(event, new Set());
    }

    this.subscribers.get(event)!.add(callback);

    // Return unsubscribe function
    return () => {
      const eventSubscribers = this.subscribers.get(event);
      if (eventSubscribers) {
        eventSubscribers.delete(callback);
        if (eventSubscribers.size === 0) {
          this.subscribers.delete(event);
        }
      }
    };
  }

  private handleMessage(message: WebSocketMessage): void {
    const { type, data } = message;

    // Handle system messages
    switch (type) {
      case 'heartbeat':
        this.send({ type: 'heartbeat_response', data: { timestamp: Date.now() } });
        break;

      case 'content_updated':
        this.handleContentUpdate(data);
        break;

      case 'user_status_changed':
        this.handleUserStatusChange(data);
        break;

      case 'message_received':
        this.handleMessageReceived(data);
        break;

      case 'sync_required':
        this.handleSyncRequired(data);
        break;
    }

    // Notify subscribers
    const eventSubscribers = this.subscribers.get(type);
    if (eventSubscribers) {
      eventSubscribers.forEach(callback => {
        try {
          callback(data);
        } catch (error) {
          console.error('WebSocket subscriber error:', error);
        }
      });
    }
  }

  private handleContentUpdate(data: ContentUpdateData): void {
    // Invalidate relevant cache entries
    store.dispatch(contentApi.util.invalidateTags([
      { type: 'Content', id: data.contentId },
      { type: 'Content', id: 'LIST' }
    ]));

    // Show notification if needed
    if (data.showNotification) {
      store.dispatch(showNotification({
        type: 'info',
        title: 'Content Updated',
        message: `${data.title} has been updated`,
        action: {
          label: 'View',
          handler: () => {
            // Navigate to content
            window.location.href = `/content/${data.contentId}`;
          }
        }
      }));
    }
  }

  private handleSyncRequired(data: SyncRequiredData): void {
    // Trigger background sync
    store.dispatch(contentApi.endpoints.syncContent.initiate({
      lastSyncTimestamp: data.lastSyncTimestamp,
      includeDeleted: true
    }));
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts++;
    const delay = Math.min(this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1), 30000);

    console.log(`🔄 Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);

    setTimeout(() => {
      this.connect().catch(() => {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          this.scheduleReconnect();
        } else {
          console.error('❌ Max reconnection attempts reached');
          this.notifyConnectionLoss();
        }
      });
    }, delay);
  }

  private startHeartbeat(): void {
    this.heartbeatInterval = setInterval(() => {
      this.send({ type: 'heartbeat', data: { timestamp: Date.now() } });
    }, 30000); // Every 30 seconds
  }

  private stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }

  private flushMessageQueue(): void {
    while (this.messageQueue.length > 0) {
      const message = this.messageQueue.shift()!;
      this.send(message);
    }
  }

  private notifyConnectionLoss(): void {
    store.dispatch(showNotification({
      type: 'error',
      title: 'Connection Lost',
      message: 'Real-time features are unavailable. Please refresh the page.',
      persistent: true,
      action: {
        label: 'Refresh',
        handler: () => window.location.reload()
      }
    }));
  }

  private setupConnectionMonitoring(): void {
    // Monitor online/offline status
    window.addEventListener('online', () => {
      if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
        this.connect().catch(console.error);
      }
    });

    window.addEventListener('offline', () => {
      store.dispatch(showNotification({
        type: 'warning',
        title: 'Offline',
        message: 'You are offline. Changes will sync when connection is restored.',
        persistent: true
      }));
    });
  }
}

// WebSocket integration with Redux
export const webSocketMiddleware: Middleware = (store) => (next) => (action) => {
  // Handle WebSocket-related actions
  if (action.type.includes('webSocket/')) {
    switch (action.type) {
      case 'webSocket/connect':
        webSocketService.connect().catch(console.error);
        break;

      case 'webSocket/disconnect':
        webSocketService.disconnect();
        break;

      case 'webSocket/send':
        webSocketService.send(action.payload);
        break;
    }
  }

  return next(action);
};
```

---

## 5. Production Monitoring and Analytics

### 5.1 Performance Monitoring Integration

#### Enterprise Performance Monitoring

```typescript
// Comprehensive performance monitoring system
export class EnterprisePerformanceMonitor {
  private metricsBuffer: PerformanceMetric[] = [];
  private flushInterval = 30000; // 30 seconds
  private maxBufferSize = 1000;

  constructor(private config: PerformanceMonitorConfig) {
    this.setupPerformanceObservers();
    this.startMetricsCollection();
    this.scheduleMetricsFlush();
  }

  private setupPerformanceObservers(): void {
    // Web Vitals monitoring
    this.observeWebVitals();

    // API performance monitoring
    this.observeAPIPerformance();

    // Component performance monitoring
    this.observeComponentPerformance();

    // Custom business metrics
    this.observeBusinessMetrics();
  }

  private observeWebVitals(): void {
    // Largest Contentful Paint
    new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        this.recordMetric({
          name: 'largest_contentful_paint',
          value: entry.startTime,
          timestamp: Date.now(),
          tags: {
            page: window.location.pathname,
            userAgent: navigator.userAgent.substring(0, 100)
          }
        });
      }
    }).observe({ entryTypes: ['largest-contentful-paint'] });

    // First Input Delay
    new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        this.recordMetric({
          name: 'first_input_delay',
          value: entry.processingStart - entry.startTime,
          timestamp: Date.now(),
          tags: {
            page: window.location.pathname,
            inputType: entry.name
          }
        });
      }
    }).observe({ entryTypes: ['first-input'] });

    // Cumulative Layout Shift
    let clsScore = 0;
    new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        if (!entry.hadRecentInput) {
          clsScore += entry.value;
        }
      }

      this.recordMetric({
        name: 'cumulative_layout_shift',
        value: clsScore,
        timestamp: Date.now(),
        tags: {
          page: window.location.pathname
        }
      });
    }).observe({ entryTypes: ['layout-shift'] });
  }

  private observeAPIPerformance(): void {
    const originalFetch = window.fetch;

    window.fetch = async (...args) => {
      const startTime = performance.now();
      const url = typeof args[0] === 'string' ? args[0] : args[0].url;

      try {
        const response = await originalFetch(...args);
        const endTime = performance.now();
        const duration = endTime - startTime;

        this.recordMetric({
          name: 'api_request_duration',
          value: duration,
          timestamp: Date.now(),
          tags: {
            url: new URL(url, window.location.origin).pathname,
            method: args[1]?.method || 'GET',
            status: response.status.toString(),
            success: response.ok.toString()
          }
        });

        // Constitutional compliance check (200ms requirement)
        if (duration > 200) {
          this.recordMetric({
            name: 'constitutional_violation',
            value: duration,
            timestamp: Date.now(),
            tags: {
              type: 'api_response_time',
              url: new URL(url, window.location.origin).pathname,
              requirement: '200ms',
              actual: `${duration.toFixed(1)}ms`
            }
          });
        }

        return response;
      } catch (error) {
        const endTime = performance.now();
        const duration = endTime - startTime;

        this.recordMetric({
          name: 'api_request_error',
          value: duration,
          timestamp: Date.now(),
          tags: {
            url: new URL(url, window.location.origin).pathname,
            method: args[1]?.method || 'GET',
            error: error.message
          }
        });

        throw error;
      }
    };
  }

  private observeComponentPerformance(): void {
    // React component performance monitoring
    const originalCreateElement = React.createElement;

    React.createElement = function(type, props, ...children) {
      const startTime = performance.now();

      try {
        const element = originalCreateElement.apply(this, arguments);
        const endTime = performance.now();
        const renderTime = endTime - startTime;

        // Only monitor complex components (render time > 1ms)
        if (renderTime > 1) {
          performanceMonitor.recordMetric({
            name: 'component_render_time',
            value: renderTime,
            timestamp: Date.now(),
            tags: {
              component: typeof type === 'string' ? type : type.name || 'Anonymous',
              hasProps: (props && Object.keys(props).length > 0).toString()
            }
          });

          // Constitutional compliance check (16ms for 60fps)
          if (renderTime > 16) {
            performanceMonitor.recordMetric({
              name: 'constitutional_violation',
              value: renderTime,
              timestamp: Date.now(),
              tags: {
                type: 'component_render_time',
                component: typeof type === 'string' ? type : type.name || 'Anonymous',
                requirement: '16ms',
                actual: `${renderTime.toFixed(1)}ms`
              }
            });
          }
        }

        return element;
      } catch (error) {
        const endTime = performance.now();
        const renderTime = endTime - startTime;

        performanceMonitor.recordMetric({
          name: 'component_render_error',
          value: renderTime,
          timestamp: Date.now(),
          tags: {
            component: typeof type === 'string' ? type : type.name || 'Anonymous',
            error: error.message
          }
        });

        throw error;
      }
    };
  }

  private observeBusinessMetrics(): void {
    // Track user interactions
    document.addEventListener('click', (event) => {
      const target = event.target as HTMLElement;
      const componentName = target.getAttribute('data-component') || 'unknown';

      this.recordMetric({
        name: 'user_interaction',
        value: 1,
        timestamp: Date.now(),
        tags: {
          type: 'click',
          component: componentName,
          page: window.location.pathname
        }
      });
    });

    // Track form submissions
    document.addEventListener('submit', (event) => {
      const form = event.target as HTMLFormElement;
      const formId = form.id || form.getAttribute('data-form') || 'unknown';

      this.recordMetric({
        name: 'form_submission',
        value: 1,
        timestamp: Date.now(),
        tags: {
          formId,
          page: window.location.pathname
        }
      });
    });
  }

  recordMetric(metric: PerformanceMetric): void {
    this.metricsBuffer.push(metric);

    // Flush buffer if it's getting too large
    if (this.metricsBuffer.length >= this.maxBufferSize) {
      this.flushMetrics();
    }
  }

  private async flushMetrics(): Promise<void> {
    if (this.metricsBuffer.length === 0) return;

    const metricsToFlush = [...this.metricsBuffer];
    this.metricsBuffer = [];

    try {
      await fetch('/api/metrics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${store.getState().auth.tokens?.accessToken}`
        },
        body: JSON.stringify({
          metrics: metricsToFlush,
          session: this.getSessionInfo()
        })
      });

      console.log(`📊 Flushed ${metricsToFlush.length} performance metrics`);

    } catch (error) {
      console.error('Failed to flush performance metrics:', error);

      // Re-queue metrics for retry (keep only recent metrics)
      const recentMetrics = metricsToFlush.filter(
        metric => Date.now() - metric.timestamp < 300000 // 5 minutes
      );
      this.metricsBuffer.unshift(...recentMetrics);
    }
  }

  private scheduleMetricsFlush(): void {
    setInterval(() => {
      this.flushMetrics();
    }, this.flushInterval);
  }

  private getSessionInfo(): SessionInfo {
    return {
      sessionId: this.generateSessionId(),
      userId: store.getState().auth.user?.id,
      userAgent: navigator.userAgent,
      platform: this.detectPlatform(),
      viewport: {
        width: window.innerWidth,
        height: window.innerHeight
      },
      connection: this.getConnectionInfo()
    };
  }

  private getConnectionInfo(): ConnectionInfo {
    const connection = (navigator as any).connection || (navigator as any).mozConnection || (navigator as any).webkitConnection;

    return {
      effectiveType: connection?.effectiveType || 'unknown',
      downlink: connection?.downlink || 0,
      rtt: connection?.rtt || 0,
      saveData: connection?.saveData || false
    };
  }
}

// Initialize performance monitoring
export const performanceMonitor = new EnterprisePerformanceMonitor({
  apiEndpoint: '/api/metrics',
  flushInterval: 30000,
  maxBufferSize: 1000,
  enableWebVitals: true,
  enableAPIMonitoring: true,
  enableComponentMonitoring: true,
  enableBusinessMetrics: true
});
```

---

This comprehensive API Integration Guide provides enterprise-grade documentation for production deployment, covering all 12 RTK Query endpoints, advanced caching strategies, cross-platform authentication, real-time features, error handling, and performance monitoring integration.