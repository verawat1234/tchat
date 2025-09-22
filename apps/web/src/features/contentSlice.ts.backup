import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import type { ContentValue, ContentState } from '../types/content';
import type { RootState } from '../store';

/**
 * Content preferences interface for user content display settings
 */
interface ContentPreferences {
  /** Whether to show draft content */
  showDrafts: boolean;
  /** Use compact view for content lists */
  compactView: boolean;
}

/**
 * Initial state for content slice based on ContentState interface
 */
const initialState: ContentState = {
  selectedLanguage: 'en',
  contentPreferences: {
    showDrafts: false,
    compactView: false,
  },
  lastSyncTime: new Date().toISOString(),
  syncStatus: 'idle',
  fallbackMode: false,
  fallbackContent: {},
};

/**
 * Content slice for managing dynamic content state
 *
 * This slice handles:
 * - Language selection for localized content
 * - User preferences for content display
 * - Synchronization status and fallback mode
 * - Local fallback content cache
 */
const contentSlice = createSlice({
  name: 'content',
  initialState,
  reducers: {
    /**
     * Set the selected language for content localization
     * @param state - Current content state
     * @param action - Action with language code payload
     */
    setSelectedLanguage: (state, action: PayloadAction<string>) => {
      state.selectedLanguage = action.payload;
    },

    /**
     * Update content preferences (partial update)
     * @param state - Current content state
     * @param action - Action with partial content preferences
     */
    updateContentPreferences: (state, action: PayloadAction<Partial<ContentPreferences>>) => {
      state.contentPreferences = {
        ...state.contentPreferences,
        ...action.payload,
      };
    },

    /**
     * Set the synchronization status
     * @param state - Current content state
     * @param action - Action with sync status and optional timestamp
     */
    setSyncStatus: (state, action: PayloadAction<{
      status: 'idle' | 'syncing' | 'error';
      timestamp?: string;
    }>) => {
      state.syncStatus = action.payload.status;
      if (action.payload.timestamp) {
        state.lastSyncTime = action.payload.timestamp;
      } else if (action.payload.status === 'idle') {
        // Update sync time when returning to idle status (successful sync)
        state.lastSyncTime = new Date().toISOString();
      }
    },

    /**
     * Toggle fallback mode on/off
     * @param state - Current content state
     * @param action - Action with fallback mode boolean
     */
    toggleFallbackMode: (state, action: PayloadAction<boolean>) => {
      state.fallbackMode = action.payload;
    },

    /**
     * Update fallback content cache with new content
     * @param state - Current content state
     * @param action - Action with content ID and value to cache
     */
    updateFallbackContent: (state, action: PayloadAction<{
      contentId: string;
      content: ContentValue;
    }>) => {
      state.fallbackContent[action.payload.contentId] = action.payload.content;
    },

    /**
     * Clear all fallback content from cache
     * @param state - Current content state
     */
    clearFallbackContent: (state) => {
      state.fallbackContent = {};
    },

    /**
     * Clear specific fallback content by ID
     * @param state - Current content state
     * @param action - Action with content ID to remove from cache
     */
    removeFallbackContent: (state, action: PayloadAction<string>) => {
      delete state.fallbackContent[action.payload];
    },

    /**
     * Batch update fallback content
     * @param state - Current content state
     * @param action - Action with record of content IDs to content values
     */
    batchUpdateFallbackContent: (state, action: PayloadAction<Record<string, ContentValue>>) => {
      state.fallbackContent = {
        ...state.fallbackContent,
        ...action.payload,
      };
    },
  },
  /**
   * extraReducers for API integration with RTK Query endpoints
   *
   * Handles state updates for:
   * - Content fetch operations (get, list, sync)
   * - Content mutation operations (create, update, delete)
   * - Bulk operations and category management
   * - Error handling and fallback mode activation
   */
  extraReducers: (builder) => {
    // Note: These matchers will be activated when contentApi endpoints are implemented
    // The patterns are ready for integration with actual RTK Query endpoints

    builder
      // ===== QUERY OPERATIONS =====

      // Content fetch operations - set syncing status
      .addMatcher(
        (action) => action.type.endsWith('/pending') && action.type.includes('content'),
        (state) => {
          state.syncStatus = 'syncing';
        }
      )

      // Successful content operations - update sync status and cache fallback content
      .addMatcher(
        (action) => action.type.endsWith('/fulfilled') && action.type.includes('getContent'),
        (state, action) => {
          state.syncStatus = 'idle';
          state.lastSyncTime = new Date().toISOString();

          // Cache successful content fetches for fallback
          if (action.payload && action.meta?.arg) {
            const contentId = typeof action.meta.arg === 'string' ? action.meta.arg : action.meta.arg.id;
            if (contentId && action.payload.value) {
              state.fallbackContent[contentId] = action.payload.value;
            }
          }

          // Disable fallback mode on successful fetch
          if (state.fallbackMode) {
            state.fallbackMode = false;
          }
        }
      )

      // Successful content list operations - batch update fallback cache
      .addMatcher(
        (action) => action.type.endsWith('/fulfilled') && (
          action.type.includes('getContentItems') ||
          action.type.includes('getContentByCategory')
        ),
        (state, action) => {
          state.syncStatus = 'idle';
          state.lastSyncTime = new Date().toISOString();

          // Batch update fallback content from list results
          if (action.payload?.items) {
            const fallbackUpdates: Record<string, ContentValue> = {};
            action.payload.items.forEach((item: any) => {
              if (item.id && item.value) {
                fallbackUpdates[item.id] = item.value;
              }
            });

            if (Object.keys(fallbackUpdates).length > 0) {
              state.fallbackContent = {
                ...state.fallbackContent,
                ...fallbackUpdates,
              };
            }
          }

          // Disable fallback mode on successful list fetch
          if (state.fallbackMode) {
            state.fallbackMode = false;
          }
        }
      )

      // Successful sync operations - comprehensive sync state update
      .addMatcher(
        (action) => action.type.endsWith('/fulfilled') && action.type.includes('syncContent'),
        (state, action) => {
          state.syncStatus = 'idle';

          if (action.payload) {
            // Update sync time from server response
            if (action.payload.syncTime) {
              state.lastSyncTime = action.payload.syncTime;
            }

            // Update fallback content with synced items
            if (action.payload.items) {
              const fallbackUpdates: Record<string, ContentValue> = {};
              action.payload.items.forEach((item: any) => {
                if (item.id && item.value) {
                  fallbackUpdates[item.id] = item.value;
                }
              });

              state.fallbackContent = {
                ...state.fallbackContent,
                ...fallbackUpdates,
              };
            }

            // Remove deleted content from fallback cache
            if (action.payload.deletedIds) {
              action.payload.deletedIds.forEach((deletedId: string) => {
                delete state.fallbackContent[deletedId];
              });
            }
          }

          // Successful sync disables fallback mode
          if (state.fallbackMode) {
            state.fallbackMode = false;
          }
        }
      )

      // ===== MUTATION OPERATIONS =====

      // Successful content mutations - update fallback cache
      .addMatcher(
        (action) => action.type.endsWith('/fulfilled') && (
          action.type.includes('updateContent') ||
          action.type.includes('createContent') ||
          action.type.includes('publishContent')
        ),
        (state, action) => {
          state.syncStatus = 'idle';
          state.lastSyncTime = new Date().toISOString();

          // Update fallback cache with mutated content
          if (action.payload?.id && action.payload?.value) {
            state.fallbackContent[action.payload.id] = action.payload.value;
          }
        }
      )

      // Successful bulk operations - batch update fallback cache
      .addMatcher(
        (action) => action.type.endsWith('/fulfilled') && action.type.includes('bulkUpdate'),
        (state, action) => {
          state.syncStatus = 'idle';
          state.lastSyncTime = new Date().toISOString();

          // Update fallback cache with bulk operation results
          if (Array.isArray(action.payload)) {
            const fallbackUpdates: Record<string, ContentValue> = {};
            action.payload.forEach((item: any) => {
              if (item.id && item.value) {
                fallbackUpdates[item.id] = item.value;
              }
            });

            if (Object.keys(fallbackUpdates).length > 0) {
              state.fallbackContent = {
                ...state.fallbackContent,
                ...fallbackUpdates,
              };
            }
          }
        }
      )

      // Successful archive/delete operations - remove from fallback cache
      .addMatcher(
        (action) => action.type.endsWith('/fulfilled') && (
          action.type.includes('archiveContent') ||
          action.type.includes('deleteContent')
        ),
        (state, action) => {
          state.syncStatus = 'idle';
          state.lastSyncTime = new Date().toISOString();

          // Remove archived/deleted content from fallback cache
          if (action.meta?.arg) {
            const contentId = typeof action.meta.arg === 'string' ? action.meta.arg : action.meta.arg.id;
            if (contentId && state.fallbackContent[contentId]) {
              delete state.fallbackContent[contentId];
            }
          }
        }
      )

      // ===== ERROR HANDLING =====

      // Failed content operations - activate fallback mode
      .addMatcher(
        (action) => action.type.endsWith('/rejected') && action.type.includes('content'),
        (state, action) => {
          state.syncStatus = 'error';

          // Activate fallback mode for network/server errors (not client validation errors)
          if (action.error && !action.error.message?.includes('validation')) {
            state.fallbackMode = true;
          }

          // Log error context for debugging (in development)
          if (process.env.NODE_ENV === 'development') {
            console.warn('Content operation failed:', {
              type: action.type,
              error: action.error,
              meta: action.meta,
            });
          }
        }
      )

      // ===== CATEGORY OPERATIONS =====

      // Successful category operations don't affect content cache directly
      // but they reset sync status
      .addMatcher(
        (action) => action.type.endsWith('/fulfilled') && action.type.includes('Category'),
        (state) => {
          state.syncStatus = 'idle';
          state.lastSyncTime = new Date().toISOString();
        }
      )

      // Failed category operations
      .addMatcher(
        (action) => action.type.endsWith('/rejected') && action.type.includes('Category'),
        (state) => {
          state.syncStatus = 'error';
          // Category failures don't necessarily require fallback mode
        }
      );
  },
});

// Export actions
export const {
  setSelectedLanguage,
  updateContentPreferences,
  setSyncStatus,
  toggleFallbackMode,
  updateFallbackContent,
  clearFallbackContent,
  removeFallbackContent,
  batchUpdateFallbackContent,
} = contentSlice.actions;

// Export reducer as default
export default contentSlice.reducer;

// State selectors for use with useSelector
export const selectContentState = (state: RootState) => state.content;
export const selectSelectedLanguage = (state: RootState) => state.content.selectedLanguage;
export const selectContentPreferences = (state: RootState) => state.content.contentPreferences;
export const selectSyncStatus = (state: RootState) => state.content.syncStatus;
export const selectLastSyncTime = (state: RootState) => state.content.lastSyncTime;
export const selectFallbackMode = (state: RootState) => state.content.fallbackMode;
export const selectFallbackContent = (state: RootState) => state.content.fallbackContent;

/**
 * Selector to get fallback content by ID
 * @param state - Root state
 * @param contentId - Content ID to retrieve
 * @returns Content value if exists in fallback cache
 */
export const selectFallbackContentById = (state: RootState, contentId: string): ContentValue | undefined => {
  return state.content.fallbackContent[contentId];
};

/**
 * Selector to check if content is available in fallback cache
 * @param state - Root state
 * @param contentId - Content ID to check
 * @returns True if content exists in fallback cache
 */
export const selectHasFallbackContent = (state: RootState, contentId: string): boolean => {
  return contentId in state.content.fallbackContent;
};

/**
 * Selector to get the count of items in fallback cache
 * @param state - Root state
 * @returns Number of items in fallback content cache
 */
export const selectFallbackContentCount = (state: RootState): number => {
  return Object.keys(state.content.fallbackContent).length;
};