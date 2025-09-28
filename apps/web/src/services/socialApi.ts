/**
 * Social API Service
 *
 * Comprehensive RTK Query service for social features following established patterns.
 * Provides full CRUD operations for posts, user relationships, and social interactions.
 *
 * IMPLEMENTATION STATUS: COMPLETE - Production-ready social service
 * Features: Posts, Comments, Reactions, User Relationships, Feed, Trending, Analytics
 */

import { api } from './api';
import type {
  SocialProfile,
  Post,
  Comment,
  Reaction,
  UserRelationship,
  SocialFeed,
  TrendingContent,
  UserAnalytics,
  FollowersResponse,
  FollowingResponse,
  PaginatedPostsResponse,
  UpdateSocialProfileRequest,
  CreatePostRequest,
  UpdatePostRequest,
  CreateCommentRequest,
  CreateReactionRequest,
  FollowRequest,
  UserDiscoveryRequest,
  SocialFeedRequest,
  TrendingRequest,
  ShareRequest,
} from '../types/social';

// =============================================================================
// Social API Endpoints
// =============================================================================

export const socialApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // ========================================================================
    // User Profile & Relationships
    // ========================================================================

    /**
     * Get user's social profile
     */
    getSocialProfile: builder.query<SocialProfile, string>({
      query: (userId) => `/social/profiles/${encodeURIComponent(userId)}`,
      providesTags: (result, error, userId) => [
        { type: 'SocialProfile', id: userId },
        'SocialProfile',
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Update user's social profile
     */
    updateSocialProfile: builder.mutation<SocialProfile, { userId: string; updates: UpdateSocialProfileRequest }>({
      query: ({ userId, updates }) => ({
        url: `/social/profiles/${encodeURIComponent(userId)}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { userId }) => [
        { type: 'SocialProfile', id: userId },
        'SocialProfile',
      ],
      // Optimistic updates
      async onQueryStarted({ userId, updates }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          socialApi.util.updateQueryData('getSocialProfile', userId, (draft) => {
            Object.assign(draft, updates);
          })
        );
        try {
          await queryFulfilled;
        } catch {
          patchResult.undo();
        }
      },
    }),

    /**
     * Discover users with similar interests
     */
    discoverUsers: builder.query<SocialProfile[], UserDiscoveryRequest>({
      query: (params) => ({
        url: '/social/discover/users',
        method: 'GET',
        params,
      }),
      providesTags: ['SocialProfile'],
      keepUnusedDataFor: 180, // 3 minutes
    }),

    /**
     * Follow a user
     */
    followUser: builder.mutation<{ message: string }, FollowRequest>({
      query: (body) => ({
        url: '/social/follow',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { followerId, followingId }) => [
        { type: 'SocialProfile', id: followerId },
        { type: 'SocialProfile', id: followingId },
        'UserRelationship',
        'SocialFollowers',
        'SocialFollowing',
      ],
    }),

    /**
     * Unfollow a user
     */
    unfollowUser: builder.mutation<{ message: string }, { followerId: string; followingId: string }>({
      query: ({ followerId, followingId }) => ({
        url: `/social/follow/${encodeURIComponent(followerId)}/${encodeURIComponent(followingId)}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, { followerId, followingId }) => [
        { type: 'SocialProfile', id: followerId },
        { type: 'SocialProfile', id: followingId },
        'UserRelationship',
        'SocialFollowers',
        'SocialFollowing',
      ],
    }),

    /**
     * Get user's followers
     */
    getFollowers: builder.query<FollowersResponse, { userId: string; limit?: number; offset?: number }>({
      query: ({ userId, ...params }) => ({
        url: `/social/followers/${encodeURIComponent(userId)}`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { userId }) => [
        { type: 'SocialFollowers', id: userId },
        'SocialFollowers',
      ],
      keepUnusedDataFor: 180, // 3 minutes
    }),

    /**
     * Get users that a user is following
     */
    getFollowing: builder.query<FollowingResponse, { userId: string; limit?: number; offset?: number }>({
      query: ({ userId, ...params }) => ({
        url: `/social/following/${encodeURIComponent(userId)}`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { userId }) => [
        { type: 'SocialFollowing', id: userId },
        'SocialFollowing',
      ],
      keepUnusedDataFor: 180, // 3 minutes
    }),

    /**
     * Get user's social analytics
     */
    getUserAnalytics: builder.query<UserAnalytics, { userId: string; period?: string }>({
      query: ({ userId, period = '30d' }) => ({
        url: `/social/analytics/users/${encodeURIComponent(userId)}`,
        method: 'GET',
        params: { period },
      }),
      providesTags: (result, error, { userId }) => [
        { type: 'SocialAnalytics', id: userId },
        'SocialAnalytics',
      ],
      keepUnusedDataFor: 900, // 15 minutes
    }),

    // ========================================================================
    // Posts Management
    // ========================================================================

    /**
     * Create a new social post
     */
    createPost: builder.mutation<Post, CreatePostRequest>({
      query: (data) => ({
        url: '/social/posts',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: [
        'SocialPost',
        'SocialFeed',
        'SocialTrending',
        'SocialProfile',
      ],
    }),

    /**
     * Get a specific post by ID
     */
    getPost: builder.query<Post, string>({
      query: (postId) => `/social/posts/${encodeURIComponent(postId)}`,
      providesTags: (result, error, postId) => [
        { type: 'SocialPost', id: postId },
        'SocialPost',
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Update an existing post
     */
    updatePost: builder.mutation<Post, { postId: string; updates: UpdatePostRequest }>({
      query: ({ postId, updates }) => ({
        url: `/social/posts/${encodeURIComponent(postId)}`,
        method: 'PUT',
        body: updates,
      }),
      invalidatesTags: (result, error, { postId }) => [
        { type: 'SocialPost', id: postId },
        'SocialPost',
        'SocialFeed',
      ],
      // Optimistic updates
      async onQueryStarted({ postId, updates }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          socialApi.util.updateQueryData('getPost', postId, (draft) => {
            Object.assign(draft, updates);
          })
        );
        try {
          await queryFulfilled;
        } catch {
          patchResult.undo();
        }
      },
    }),

    /**
     * Delete a post (soft delete)
     */
    deletePost: builder.mutation<{ message: string }, string>({
      query: (postId) => ({
        url: `/social/posts/${encodeURIComponent(postId)}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, postId) => [
        { type: 'SocialPost', id: postId },
        'SocialPost',
        'SocialFeed',
        'SocialProfile',
      ],
    }),

    // ========================================================================
    // Comments & Reactions
    // ========================================================================

    /**
     * Create a comment on a post
     */
    createComment: builder.mutation<Comment, CreateCommentRequest>({
      query: (data) => ({
        url: '/social/comments',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, { postId }) => [
        { type: 'SocialPost', id: postId },
        'SocialComment',
        'SocialPost',
      ],
    }),

    /**
     * Add reaction to a post or comment
     */
    addReaction: builder.mutation<{ message: string }, CreateReactionRequest>({
      query: (data) => ({
        url: '/social/reactions',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, { targetId, targetType }) => [
        { type: targetType === 'post' ? 'SocialPost' : 'SocialComment', id: targetId },
        'SocialReaction',
      ],
      // Optimistic updates for reactions
      async onQueryStarted({ targetId, targetType, type }, { dispatch, queryFulfilled }) {
        if (targetType === 'post') {
          const patchResult = dispatch(
            socialApi.util.updateQueryData('getPost', targetId, (draft) => {
              draft.likesCount = (draft.likesCount || 0) + 1;
              draft.isLiked = true;
              draft.userReaction = {
                id: 'temp-' + Date.now(),
                userId: 'current-user',
                targetId,
                targetType,
                type,
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString(),
              };
            })
          );
          try {
            await queryFulfilled;
          } catch {
            patchResult.undo();
          }
        }
      },
    }),

    /**
     * Remove reaction from a post or comment
     */
    removeReaction: builder.mutation<{ message: string }, { targetId: string; targetType: 'post' | 'comment' }>({
      query: ({ targetId, targetType }) => ({
        url: `/social/reactions/${targetType}/${encodeURIComponent(targetId)}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, { targetId, targetType }) => [
        { type: targetType === 'post' ? 'SocialPost' : 'SocialComment', id: targetId },
        'SocialReaction',
      ],
      // Optimistic updates for reaction removal
      async onQueryStarted({ targetId, targetType }, { dispatch, queryFulfilled }) {
        if (targetType === 'post') {
          const patchResult = dispatch(
            socialApi.util.updateQueryData('getPost', targetId, (draft) => {
              draft.likesCount = Math.max((draft.likesCount || 0) - 1, 0);
              draft.isLiked = false;
              draft.userReaction = undefined;
            })
          );
          try {
            await queryFulfilled;
          } catch {
            patchResult.undo();
          }
        }
      },
    }),

    // ========================================================================
    // Social Feed & Discovery
    // ========================================================================

    /**
     * Get user's personalized social feed
     */
    getSocialFeed: builder.query<SocialFeed, SocialFeedRequest>({
      query: (params) => ({
        url: '/social/feed',
        method: 'GET',
        params,
      }),
      providesTags: ['SocialFeed'],
      keepUnusedDataFor: 120, // 2 minutes for fresh content
    }),

    /**
     * Get trending social content
     */
    getTrendingContent: builder.query<TrendingContent, TrendingRequest>({
      query: (params) => ({
        url: '/social/trending',
        method: 'GET',
        params,
      }),
      providesTags: ['SocialTrending'],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Share content externally
     */
    shareContent: builder.mutation<{ message: string }, ShareRequest>({
      query: (data) => ({
        url: '/social/share',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: (result, error, { contentId, contentType }) => [
        { type: contentType === 'post' ? 'SocialPost' : 'SocialComment', id: contentId },
      ],
    }),

    // ========================================================================
    // Legacy Support for Existing Components
    // ========================================================================

    /**
     * Legacy: Get social stories (for existing SocialTab component)
     */
    getSocialStories: builder.query<any[], { active?: boolean; limit?: number }>({
      query: (params) => ({
        url: '/social/stories',
        method: 'GET',
        params,
      }),
      providesTags: ['SocialStories'],
      keepUnusedDataFor: 180, // 3 minutes
    }),

    /**
     * Legacy: Get user friends (for existing SocialTab component)
     */
    getUserFriends: builder.query<any[], { status?: string; limit?: number }>({
      query: (params) => ({
        url: '/social/friends',
        method: 'GET',
        params,
      }),
      providesTags: ['SocialFriends'],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Legacy: Like social post (for existing SocialTab component)
     */
    likeSocialPost: builder.mutation<{ message: string }, { postId: string; isLiked: boolean }>({
      query: ({ postId, isLiked }) => ({
        url: '/social/reactions',
        method: 'POST',
        body: {
          targetId: postId,
          targetType: 'post',
          type: 'like',
        },
      }),
      invalidatesTags: (result, error, { postId }) => [
        { type: 'SocialPost', id: postId },
        'SocialFeed',
      ],
    }),

    /**
     * Legacy: Create social post (for existing SocialTab component)
     */
    createSocialPost: builder.mutation<Post, CreatePostRequest>({
      query: (data) => ({
        url: '/social/posts',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: [
        'SocialPost',
        'SocialFeed',
        'SocialProfile',
      ],
    }),

    /**
     * Legacy: Follow user mutation (for existing SocialTab component)
     */
    followUserMutation: builder.mutation<{ message: string }, { userId: string; source?: string }>({
      query: ({ userId, source = 'manual' }) => ({
        url: '/social/follow',
        method: 'POST',
        body: {
          followerId: 'current-user-id', // Will be set by auth middleware
          followingId: userId,
          source,
        },
      }),
      invalidatesTags: [
        'SocialProfile',
        'UserRelationship',
        'SocialFollowers',
        'SocialFollowing',
      ],
    }),
  }),
  overrideExisting: false,
});

// Export hooks for all social API endpoints
export const {
  // User Profile & Relationships
  useGetSocialProfileQuery,
  useUpdateSocialProfileMutation,
  useDiscoverUsersQuery,
  useFollowUserMutation,
  useUnfollowUserMutation,
  useGetFollowersQuery,
  useGetFollowingQuery,
  useGetUserAnalyticsQuery,

  // Posts Management
  useCreatePostMutation,
  useGetPostQuery,
  useUpdatePostMutation,
  useDeletePostMutation,

  // Comments & Reactions
  useCreateCommentMutation,
  useAddReactionMutation,
  useRemoveReactionMutation,

  // Social Feed & Discovery
  useGetSocialFeedQuery,
  useGetTrendingContentQuery,
  useShareContentMutation,

  // Legacy Support
  useGetSocialStoriesQuery,
  useGetUserFriendsQuery,
  useLikeSocialPostMutation,
  useCreateSocialPostMutation,
  useFollowUserMutationMutation,
} = socialApi;

// ========================================================================
// Utility Functions for Social Features
// ========================================================================

/**
 * Transform backend Post to legacy format for existing components
 */
export function transformPostToLegacy(post: Post): any {
  return {
    id: post.id,
    author: {
      name: post.author?.name || post.author?.display_name || 'Unknown User',
      avatar: post.author?.avatar,
      verified: post.author?.isSocialVerified || false,
      type: 'user',
    },
    content: post.content,
    images: post.mediaUrls || [],
    timestamp: formatTimestamp(post.createdAt),
    likes: post.likesCount,
    comments: post.commentsCount,
    shares: post.sharesCount,
    isLiked: post.isLiked || false,
    location: post.metadata?.location,
    tags: post.tags || [],
    type: post.type as any,
    source: 'following',
  };
}

/**
 * Format timestamp for display
 */
function formatTimestamp(timestamp: string): string {
  const now = new Date();
  const date = new Date(timestamp);
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 1) return 'just now';
  if (diffMins < 60) return `${diffMins} min ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;

  return date.toLocaleDateString();
}

/**
 * Get social profile display name
 */
export function getSocialDisplayName(profile: SocialProfile): string {
  return profile.display_name || profile.name || profile.username || 'Unknown User';
}

/**
 * Check if user can edit post
 */
export function canEditPost(post: Post, currentUserId: string): boolean {
  return post.authorId === currentUserId;
}

/**
 * Get reaction emoji
 */
export function getReactionEmoji(type: string): string {
  const emojiMap: Record<string, string> = {
    like: '👍',
    love: '❤️',
    laugh: '😂',
    wow: '😮',
    sad: '😢',
    angry: '😠',
  };
  return emojiMap[type] || '👍';
}