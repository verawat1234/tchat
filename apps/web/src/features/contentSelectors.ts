import { createSelector } from '@reduxjs/toolkit';
import type { RootState } from '../store';
import type { ContentValue, ContentState } from '../types/content';

/**
 * Content Selectors
 *
 * Comprehensive selectors for efficient content state access with memoization
 * and performance optimizations. These selectors provide:
 * - Basic state access
 * - Fallback content management
 * - Computed values with memoization
 * - Parameterized selectors for filtering
 */

// =============================================================================
// Basic State Selectors
// =============================================================================

/**
 * Select the entire content state
 * @param state - Root Redux state
 * @returns Complete content state
 */
export const selectContentState = (state: RootState): ContentState => state.content;

/**
 * Select the currently selected language
 * @param state - Root Redux state
 * @returns Current language code (e.g., 'en', 'es', 'fr')
 */
export const selectSelectedLanguage = (state: RootState): string =>
  state.content.selectedLanguage;

/**
 * Select user content preferences
 * @param state - Root Redux state
 * @returns Content preferences object
 */
export const selectContentPreferences = (state: RootState) =>
  state.content.contentPreferences;

/**
 * Select current synchronization status
 * @param state - Root Redux state
 * @returns Sync status: 'idle', 'syncing', or 'error'
 */
export const selectSyncStatus = (state: RootState) =>
  state.content.syncStatus;

/**
 * Select last successful synchronization timestamp
 * @param state - Root Redux state
 * @returns ISO timestamp string of last sync
 */
export const selectLastSyncTime = (state: RootState): string =>
  state.content.lastSyncTime;

/**
 * Select fallback mode status
 * @param state - Root Redux state
 * @returns True if using local fallback content
 */
export const selectFallbackMode = (state: RootState): boolean =>
  state.content.fallbackMode;

// =============================================================================
// Fallback Content Selectors
// =============================================================================

/**
 * Select all fallback content from cache
 * @param state - Root Redux state
 * @returns Record of content IDs to content values
 */
export const selectFallbackContent = (state: RootState): Record<string, ContentValue> =>
  state.content.fallbackContent;

/**
 * Factory function to create a selector for specific fallback content by ID
 * @param contentId - Content ID to retrieve
 * @returns Selector function that returns content value or undefined
 */
export const selectFallbackContentById = (contentId: string) =>
  createSelector(
    [selectFallbackContent],
    (fallbackContent): ContentValue | undefined => fallbackContent[contentId]
  );

/**
 * Select whether any fallback content exists in cache
 * @param state - Root Redux state
 * @returns True if any fallback content is available
 */
export const selectHasFallbackContent = createSelector(
  [selectFallbackContent],
  (fallbackContent): boolean => Object.keys(fallbackContent).length > 0
);

/**
 * Select the count of cached fallback content items
 * @param state - Root Redux state
 * @returns Number of items in fallback content cache
 */
export const selectFallbackContentCount = createSelector(
  [selectFallbackContent],
  (fallbackContent): number => Object.keys(fallbackContent).length
);

// =============================================================================
// Computed Selectors with Memoization
// =============================================================================

/**
 * Select content preferences with default values applied
 * Useful for ensuring consistent behavior when preferences are undefined
 */
export const selectContentPreferencesWithDefaults = createSelector(
  [selectContentPreferences],
  (preferences) => ({
    showDrafts: preferences?.showDrafts ?? false,
    compactView: preferences?.compactView ?? false,
  })
);

/**
 * Select user-friendly display status for synchronization
 * Transforms technical sync status into user-readable format
 */
export const selectSyncStatusDisplay = createSelector(
  [selectSyncStatus, selectFallbackMode],
  (syncStatus, fallbackMode): string => {
    if (fallbackMode) {
      return 'Offline Mode';
    }

    switch (syncStatus) {
      case 'idle':
        return 'Synchronized';
      case 'syncing':
        return 'Synchronizing...';
      case 'error':
        return 'Sync Error';
      default:
        return 'Unknown';
    }
  }
);

/**
 * Select whether content is stale and needs refresh
 * Content is considered stale if it hasn't been synced in the last 5 minutes
 */
export const selectIsContentStale = createSelector(
  [selectLastSyncTime, selectSyncStatus],
  (lastSyncTime, syncStatus): boolean => {
    if (syncStatus === 'syncing') {
      return false; // Currently syncing, not stale
    }

    const lastSync = new Date(lastSyncTime);
    const now = new Date();
    const fiveMinutesAgo = new Date(now.getTime() - 5 * 60 * 1000);

    return lastSync < fiveMinutesAgo;
  }
);

/**
 * Factory function to create a selector for fallback content filtered by category
 * Uses content ID pattern matching to filter by category prefix
 * @param category - Category to filter by (e.g., 'navigation', 'error')
 * @returns Selector function that returns filtered fallback content
 */
export const selectFallbackContentByCategory = (category: string) =>
  createSelector(
    [selectFallbackContent],
    (fallbackContent): Record<string, ContentValue> => {
      const categoryPrefix = `${category}.`;
      return Object.entries(fallbackContent)
        .filter(([contentId]) => contentId.startsWith(categoryPrefix))
        .reduce((filtered, [contentId, content]) => {
          filtered[contentId] = content;
          return filtered;
        }, {} as Record<string, ContentValue>);
    }
  );

// =============================================================================
// Performance-Optimized Parameterized Selectors
// =============================================================================

/**
 * High-performance selector factory for checking if specific content exists in fallback
 * Uses memoization to avoid repeated lookups for the same content ID
 * @param contentId - Content ID to check
 * @returns Selector function that returns boolean
 */
export const selectHasSpecificFallbackContent = (contentId: string) =>
  createSelector(
    [selectFallbackContent],
    (fallbackContent): boolean => contentId in fallbackContent
  );

/**
 * Selector for content synchronization metadata
 * Provides comprehensive sync state information in a single object
 */
export const selectContentSyncMetadata = createSelector(
  [selectSyncStatus, selectLastSyncTime, selectFallbackMode, selectIsContentStale],
  (syncStatus, lastSyncTime, fallbackMode, isStale) => ({
    status: syncStatus,
    lastSyncTime,
    fallbackMode,
    isStale,
    canSync: !fallbackMode && syncStatus !== 'syncing',
    displayStatus: fallbackMode ? 'Offline Mode' :
                  syncStatus === 'idle' ? 'Synchronized' :
                  syncStatus === 'syncing' ? 'Synchronizing...' : 'Sync Error',
  })
);

/**
 * Selector for content cache statistics
 * Provides analytics about the fallback content cache
 */
export const selectContentCacheStatistics = createSelector(
  [selectFallbackContent],
  (fallbackContent) => {
    const entries = Object.entries(fallbackContent);
    const categories = new Set(
      entries.map(([contentId]) => contentId.split('.')[0])
    );

    const typeStats = entries.reduce((stats, [, content]) => {
      const type = content.type;
      stats[type] = (stats[type] || 0) + 1;
      return stats;
    }, {} as Record<string, number>);

    return {
      totalItems: entries.length,
      uniqueCategories: categories.size,
      categories: Array.from(categories),
      typeDistribution: typeStats,
      isEmpty: entries.length === 0,
    };
  }
);

// =============================================================================
// Advanced Computed Selectors
// =============================================================================

/**
 * Selector for determining content loading state
 * Combines multiple state indicators to provide comprehensive loading status
 */
export const selectContentLoadingState = createSelector(
  [selectSyncStatus, selectFallbackMode, selectHasFallbackContent],
  (syncStatus, fallbackMode, hasFallbackContent) => ({
    isLoading: syncStatus === 'syncing',
    hasError: syncStatus === 'error',
    hasContent: hasFallbackContent || (!fallbackMode && syncStatus === 'idle'),
    showFallbackContent: fallbackMode && hasFallbackContent,
    shouldShowLoader: syncStatus === 'syncing' && !hasFallbackContent,
    shouldShowError: syncStatus === 'error' && !hasFallbackContent,
  })
);

/**
 * Selector for content preferences with validation
 * Ensures preferences are valid and provides corrected values
 */
export const selectValidatedContentPreferences = createSelector(
  [selectContentPreferences],
  (preferences) => {
    // Validate and correct preferences
    const validated = {
      showDrafts: Boolean(preferences?.showDrafts),
      compactView: Boolean(preferences?.compactView),
    };

    return {
      ...validated,
      isValid: preferences !== null && typeof preferences === 'object',
      needsCorrection: preferences?.showDrafts !== validated.showDrafts ||
                      preferences?.compactView !== validated.compactView,
    };
  }
);

// =============================================================================
// Export All Selectors
// =============================================================================

export {
  // Re-export basic selectors for consistency with existing imports
  selectContentState as selectContentStateBase,
  selectSelectedLanguage as selectSelectedLanguageBase,
  selectContentPreferences as selectContentPreferencesBase,
  selectSyncStatus as selectSyncStatusBase,
  selectLastSyncTime as selectLastSyncTimeBase,
  selectFallbackMode as selectFallbackModeBase,
  selectFallbackContent as selectFallbackContentBase,
};