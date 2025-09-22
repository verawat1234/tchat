/**
 * Content Fallback Middleware - RTK Query integration for localStorage fallback
 *
 * This middleware intercepts RTK Query actions to provide seamless integration
 * with the localStorage-based content fallback system. It handles caching
 * successful API responses and providing fallback content when API calls fail.
 *
 * Features:
 * - Automatic caching of successful content API responses
 * - Fallback content retrieval for failed API calls
 * - Integration with content slice state management
 * - Error handling and fallback mode activation
 * - Performance optimization through selective caching
 */

import { createListenerMiddleware, isAnyOf } from '@reduxjs/toolkit';
import type { RootState } from '../index';
import {
  contentFallbackService,
  cacheRTKQueryResponse,
  getFallbackContent,
} from '../../services/contentFallback';
import {
  toggleFallbackMode,
  updateFallbackContent,
  batchUpdateFallbackContent,
  setSyncStatus,
} from '../../features/contentSlice';

// =============================================================================
// Middleware Configuration
// =============================================================================

/** Content fallback listener middleware */
export const contentFallbackMiddleware = createListenerMiddleware();

// =============================================================================
// Helper Functions
// =============================================================================

/**
 * Check if an action is a content-related RTK Query action
 */
function isContentAction(action: any): boolean {
  return (
    action.type?.includes('content') ||
    action.type?.includes('Content') ||
    action.meta?.arg?.endpointName?.toLowerCase().includes('content')
  );
}

/**
 * Extract endpoint information from RTK Query action
 */
function extractEndpointInfo(action: any): {
  endpointName: string;
  args: any;
  result?: any;
  error?: any;
} {
  const endpointName = action.meta?.arg?.endpointName || '';
  const args = action.meta?.arg?.originalArgs;
  const result = action.payload;
  const error = action.error;

  return { endpointName, args, result, error };
}

/**
 * Check if error should trigger fallback mode
 */
function shouldActivateFallback(error: any): boolean {
  if (!error) return false;

  // Network errors
  if (error.status === 'FETCH_ERROR' || error.status === 'TIMEOUT_ERROR') {
    return true;
  }

  // Server errors (5xx)
  if (typeof error.status === 'number' && error.status >= 500) {
    return true;
  }

  // Connection refused, network unreachable
  if (error.message?.includes('NetworkError') || error.message?.includes('fetch')) {
    return true;
  }

  return false;
}

// =============================================================================
// Success Handlers - Cache API Responses
// =============================================================================

/**
 * Handle successful content API responses by caching them
 */
contentFallbackMiddleware.startListening({
  matcher: isAnyOf(
    (action) => action.type.endsWith('/fulfilled') && isContentAction(action)
  ),
  effect: async (action, listenerApi) => {
    try {
      const { endpointName, args, result } = extractEndpointInfo(action);

      if (!endpointName || !result) return;

      // Cache the successful response
      await cacheRTKQueryResponse(endpointName, args, result, action.meta);

      // Update Redux state with cached content for immediate access
      const state = listenerApi.getState() as RootState;

      if (endpointName === 'getContentItem' && result.id && result.value) {
        listenerApi.dispatch(updateFallbackContent({
          contentId: result.id,
          content: result.value,
        }));
      } else if ((endpointName === 'getContentItems' || endpointName === 'getContentByCategory') && Array.isArray(result.items || result)) {
        const items = result.items || result;
        const fallbackUpdates: Record<string, any> = {};

        items.forEach((item: any) => {
          if (item.id && item.value) {
            fallbackUpdates[item.id] = item.value;
          }
        });

        if (Object.keys(fallbackUpdates).length > 0) {
          listenerApi.dispatch(batchUpdateFallbackContent(fallbackUpdates));
        }
      }

      // If we were in fallback mode, disable it on successful fetch
      if (state.content.fallbackMode) {
        listenerApi.dispatch(toggleFallbackMode(false));
      }

      // Update sync status
      listenerApi.dispatch(setSyncStatus({
        status: 'idle',
        timestamp: new Date().toISOString()
      }));

    } catch (error) {
      console.warn('Failed to cache content response:', error);
    }
  },
});

// =============================================================================
// Error Handlers - Provide Fallback Content
// =============================================================================

/**
 * Handle failed content API requests by providing fallback content
 */
contentFallbackMiddleware.startListening({
  matcher: isAnyOf(
    (action) => action.type.endsWith('/rejected') && isContentAction(action)
  ),
  effect: async (action, listenerApi) => {
    try {
      const { endpointName, args, error } = extractEndpointInfo(action);
      const state = listenerApi.getState() as RootState;

      // Check if we should activate fallback mode
      if (shouldActivateFallback(error)) {
        // Activate fallback mode
        if (!state.content.fallbackMode) {
          listenerApi.dispatch(toggleFallbackMode(true));
        }

        // Try to get fallback content
        const fallbackContent = await getFallbackContent(endpointName, args);

        if (fallbackContent) {
          // We have fallback content available
          console.info(`Providing fallback content for ${endpointName}`, {
            args,
            fallbackAvailable: true,
          });

          // For single content items, update the fallback content in state
          if (endpointName === 'getContentItem' && typeof args === 'string') {
            listenerApi.dispatch(updateFallbackContent({
              contentId: args,
              content: fallbackContent,
            }));
          }
        } else {
          console.warn(`No fallback content available for ${endpointName}`, { args });
        }
      }

      // Update sync status to error
      listenerApi.dispatch(setSyncStatus({ status: 'error' }));

    } catch (fallbackError) {
      console.error('Error handling content fallback:', fallbackError);
    }
  },
});

// =============================================================================
// Pending Handlers - Update Sync Status
// =============================================================================

/**
 * Handle pending content API requests by updating sync status
 */
contentFallbackMiddleware.startListening({
  matcher: isAnyOf(
    (action) => action.type.endsWith('/pending') && isContentAction(action)
  ),
  effect: async (action, listenerApi) => {
    listenerApi.dispatch(setSyncStatus({ status: 'syncing' }));
  },
});

// =============================================================================
// Mutation Handlers - Update Cache on Content Changes
// =============================================================================

/**
 * Handle successful content mutations by updating cache
 */
contentFallbackMiddleware.startListening({
  matcher: isAnyOf(
    (action) =>
      action.type.endsWith('/fulfilled') &&
      isContentAction(action) &&
      (action.type.includes('create') ||
       action.type.includes('update') ||
       action.type.includes('publish'))
  ),
  effect: async (action, listenerApi) => {
    try {
      const { result } = extractEndpointInfo(action);

      // Update cache with mutated content
      if (result?.id && result?.value && result?.type) {
        await contentFallbackService.cacheContent(
          result.id,
          result.value,
          result.type,
          {
            category: result.category,
            version: result.version,
          }
        );

        // Update Redux state
        listenerApi.dispatch(updateFallbackContent({
          contentId: result.id,
          content: result.value,
        }));
      }

      // Update sync status
      listenerApi.dispatch(setSyncStatus({
        status: 'idle',
        timestamp: new Date().toISOString()
      }));

    } catch (error) {
      console.warn('Failed to update cache after mutation:', error);
    }
  },
});

// =============================================================================
// Archive/Delete Handlers - Remove from Cache
// =============================================================================

/**
 * Handle content archive/delete by removing from cache
 */
contentFallbackMiddleware.startListening({
  matcher: isAnyOf(
    (action) =>
      action.type.endsWith('/fulfilled') &&
      isContentAction(action) &&
      (action.type.includes('archive') || action.type.includes('delete'))
  ),
  effect: async (action, listenerApi) => {
    try {
      const { args } = extractEndpointInfo(action);
      const contentId = typeof args === 'string' ? args : args?.id;

      if (contentId) {
        // Remove from cache
        await contentFallbackService.removeContent(contentId);

        // Remove from Redux state
        // Note: This would require adding a removeFromFallbackContent action
        // to the content slice if not already present
      }

      // Update sync status
      listenerApi.dispatch(setSyncStatus({
        status: 'idle',
        timestamp: new Date().toISOString()
      }));

    } catch (error) {
      console.warn('Failed to remove content from cache:', error);
    }
  },
});

// =============================================================================
// Bulk Operations Handler
// =============================================================================

/**
 * Handle bulk content operations
 */
contentFallbackMiddleware.startListening({
  matcher: isAnyOf(
    (action) =>
      action.type.endsWith('/fulfilled') &&
      isContentAction(action) &&
      action.type.includes('bulk')
  ),
  effect: async (action, listenerApi) => {
    try {
      const { result } = extractEndpointInfo(action);

      if (result?.successful && Array.isArray(result.successful)) {
        // Cache successful bulk updates
        const items = result.successful.map((item: any) => ({
          contentId: item.id,
          content: item.value,
          type: item.type,
          category: item.category,
          version: item.version,
        }));

        await contentFallbackService.batchCacheContent(items);

        // Update Redux state
        const fallbackUpdates: Record<string, any> = {};
        result.successful.forEach((item: any) => {
          if (item.id && item.value) {
            fallbackUpdates[item.id] = item.value;
          }
        });

        if (Object.keys(fallbackUpdates).length > 0) {
          listenerApi.dispatch(batchUpdateFallbackContent(fallbackUpdates));
        }
      }

      // Update sync status
      listenerApi.dispatch(setSyncStatus({
        status: 'idle',
        timestamp: new Date().toISOString()
      }));

    } catch (error) {
      console.warn('Failed to handle bulk operation cache update:', error);
    }
  },
});

// =============================================================================
// Service Initialization
// =============================================================================

/**
 * Initialize fallback service on application start
 */
contentFallbackMiddleware.startListening({
  actionCreator: 'persist/REHYDRATE' as any,
  effect: async () => {
    try {
      await contentFallbackService.initialize();
      console.info('Content fallback service initialized successfully');
    } catch (error) {
      console.error('Failed to initialize content fallback service:', error);
    }
  },
});

// =============================================================================
// Maintenance Tasks
// =============================================================================

/**
 * Perform periodic maintenance
 */
let maintenanceInterval: NodeJS.Timeout;

contentFallbackMiddleware.startListening({
  predicate: () => true,
  effect: async (action, listenerApi) => {
    // Set up maintenance interval once
    if (!maintenanceInterval) {
      maintenanceInterval = setInterval(async () => {
        try {
          await contentFallbackService.performMaintenance();
        } catch (error) {
          console.warn('Fallback service maintenance failed:', error);
        }
      }, 60 * 60 * 1000); // Every hour
    }
  },
});

// =============================================================================
// Export
// =============================================================================

export default contentFallbackMiddleware;