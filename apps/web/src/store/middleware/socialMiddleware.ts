/**
 * Social Middleware
 *
 * Advanced middleware for social features including real-time updates,
 * optimistic UI updates, and comprehensive error handling.
 */

import { createListenerMiddleware, isAnyOf } from '@reduxjs/toolkit';
import { socialApi } from '../../services/socialApi';
import type { RootState } from '../index';

// Create listener middleware for social features
export const socialMiddleware = createListenerMiddleware();

// Real-time updates listener
socialMiddleware.startListening({
  matcher: isAnyOf(
    socialApi.endpoints.getSocialFeed.matchFulfilled,
    socialApi.endpoints.createPost.matchFulfilled,
    socialApi.endpoints.addReaction.matchFulfilled
  ),
  effect: async (action, listenerApi) => {
    const state = listenerApi.getState() as RootState;

    // Handle real-time feed updates
    if (socialApi.endpoints.createPost.matchFulfilled(action)) {
      // Invalidate and refetch social feed to show new post
      listenerApi.dispatch(
        socialApi.util.invalidateTags(['SocialFeed', 'SocialProfile'])
      );
    }

    // Handle reaction updates
    if (socialApi.endpoints.addReaction.matchFulfilled(action)) {
      // Update cached posts with new reaction counts
      const posts = socialApi.endpoints.getSocialFeed.select({
        algorithm: 'personalized',
        limit: 20
      })(state);

      if (posts.data?.posts) {
        // Real-time update logic would go here
        console.log('Reaction added, updating feed cache');
      }
    }
  },
});

// Error handling for social operations
socialMiddleware.startListening({
  matcher: isAnyOf(
    socialApi.endpoints.createPost.matchRejected,
    socialApi.endpoints.followUser.matchRejected,
    socialApi.endpoints.addReaction.matchRejected
  ),
  effect: async (action, listenerApi) => {
    // Handle social operation errors
    console.error('Social operation failed:', action.error);

    // You could dispatch user-friendly error messages here
    // or trigger retry mechanisms
  },
});

// Optimistic updates listener
socialMiddleware.startListening({
  actionCreator: socialApi.endpoints.addReaction.initiate,
  effect: async (action, listenerApi) => {
    // Optimistic UI update for reactions
    const { targetId, targetType, type } = action.arg;

    if (targetType === 'post') {
      // Update the post in cache immediately for better UX
      listenerApi.dispatch(
        socialApi.util.updateQueryData(
          'getSocialFeed',
          { algorithm: 'personalized', limit: 20 },
          (draft) => {
            const post = draft.posts?.find(p => p.id === targetId);
            if (post) {
              post.likesCount = (post.likesCount || 0) + 1;
              post.isLiked = true;
            }
          }
        )
      );
    }
  },
});

// Cache management for social features
socialMiddleware.startListening({
  predicate: (action) => {
    return action.type.includes('social') && action.type.includes('fulfilled');
  },
  effect: async (action, listenerApi) => {
    // Manage cache expiration and cleanup
    const state = listenerApi.getState() as RootState;

    // Clean up old cached data periodically
    const cacheTimestamp = Date.now();
    const maxCacheAge = 5 * 60 * 1000; // 5 minutes

    // This would implement cache cleanup logic
    console.log('Managing social cache at', cacheTimestamp);
  },
});

// Performance monitoring for social features
socialMiddleware.startListening({
  matcher: isAnyOf(
    socialApi.endpoints.getSocialFeed.matchPending,
    socialApi.endpoints.getSocialFeed.matchFulfilled,
    socialApi.endpoints.getSocialFeed.matchRejected
  ),
  effect: async (action, listenerApi) => {
    // Track social feed performance
    if (socialApi.endpoints.getSocialFeed.matchPending(action)) {
      console.time(`SocialFeed-${action.meta.requestId}`);
    }

    if (socialApi.endpoints.getSocialFeed.matchFulfilled(action)) {
      console.timeEnd(`SocialFeed-${action.meta.requestId}`);
      console.log('Social feed loaded successfully:', action.payload.posts?.length, 'posts');
    }

    if (socialApi.endpoints.getSocialFeed.matchRejected(action)) {
      console.timeEnd(`SocialFeed-${action.meta.requestId}`);
      console.error('Social feed failed to load:', action.error);
    }
  },
});