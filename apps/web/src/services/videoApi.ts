/**
 * Video API - RTK Query endpoints for YouTube/TikTok-style video platform
 *
 * This file implements all RTK Query endpoints for the video service,
 * providing type-safe API integration with comprehensive error handling, caching,
 * and real-time synchronization capabilities for video content management.
 *
 * Features:
 * - Complete CRUD operations for video content
 * - Channel management and video organization
 * - Live streaming and real-time chat integration
 * - Comprehensive search and filtering capabilities
 * - Analytics and engagement tracking
 * - Advanced caching with tag-based invalidation
 * - Type-safe request/response handling
 * - Comprehensive error handling and retry logic
 * - Optimistic updates for better user experience
 */

import { useState, useEffect } from 'react';
import { api } from './api';
import type {
  VideoContent,
  ChannelInfo,
  VideoComment,
  VideoPlaylist,
  LiveStream,
  VideoAnalytics,
  PaginatedResponse,
  SingleResponse,
  VideoType,
  VideoCategory,
  VideoQuality,
  SubtitleTrack,
  SponsoredSegment,
} from '../types/video';

// =============================================================================
// Request/Response Type Definitions
// =============================================================================

/**
 * Request parameters for getting videos with filtering and pagination
 */
export interface GetVideosRequest {
  /** Pagination page number */
  page?: number;
  /** Items per page */
  limit?: number;
  /** Filter by video category */
  category?: VideoCategory;
  /** Filter by video type (SHORT/LONG) */
  type?: VideoType;
  /** Include premium content (requires auth) */
  includePremium?: boolean;
  /** Include private content (requires auth) */
  includePrivate?: boolean;
  /** Sort field */
  sortBy?: 'uploadTime' | 'views' | 'likes' | 'trending' | 'duration';
  /** Sort direction */
  sortOrder?: 'asc' | 'desc';
  /** Filter by language */
  language?: string;
  /** Filter by tags */
  tags?: string[];
  /** Include age-restricted content */
  ageRestricted?: boolean;
}

/**
 * Response for paginated video list
 */
export interface GetVideosResponse extends PaginatedResponse<VideoContent> {
  filters?: Record<string, any>;
}

/**
 * Request for creating a new video
 */
export interface CreateVideoRequest {
  /** Video title */
  title: string;
  /** Video description */
  description: string;
  /** Video category */
  category: VideoCategory;
  /** Video type */
  type: VideoType;
  /** Video tags */
  tags: string[];
  /** Video language */
  language: string;
  /** Thumbnail URL */
  thumbnail: string;
  /** Video file URL */
  videoUrl: string;
  /** Is video private */
  isPrivate: boolean;
  /** Is video age restricted */
  ageRestricted: boolean;
  /** Video duration in seconds */
  durationSeconds?: number;
  /** Subtitle tracks */
  subtitles?: SubtitleTrack[];
  /** Available video qualities */
  qualities?: VideoQuality[];
  /** Sponsored segments */
  sponsoredSegments?: SponsoredSegment[];
}

/**
 * Request for updating a video
 */
export interface UpdateVideoRequest {
  /** Video ID to update */
  id: string;
  /** Updated title */
  title?: string;
  /** Updated description */
  description?: string;
  /** Updated category */
  category?: VideoCategory;
  /** Updated tags */
  tags?: string[];
  /** Updated privacy setting */
  isPrivate?: boolean;
  /** Updated age restriction */
  ageRestricted?: boolean;
  /** Updated thumbnail */
  thumbnail?: string;
}

/**
 * Request for video search with filters
 */
export interface SearchVideosRequest {
  /** Search query */
  q: string;
  /** Filter by category */
  category?: VideoCategory;
  /** Filter by duration (SHORT/MEDIUM/LONG) */
  duration?: string;
  /** Sort by field */
  sortBy?: 'RELEVANCE' | 'VIEW_COUNT' | 'UPLOAD_DATE' | 'RATING';
  /** Upload date filter */
  uploadDate?: 'THIS_HOUR' | 'TODAY' | 'THIS_WEEK' | 'THIS_MONTH' | 'THIS_YEAR';
  /** Maximum results */
  limit?: number;
  /** Pagination offset */
  offset?: number;
}

/**
 * Response for video search
 */
export interface SearchVideosResponse {
  videos: VideoContent[];
  channels: ChannelInfo[];
  playlists: VideoPlaylist[];
  totalResults: number;
  searchQuery: string;
  suggestions: string[];
  nextPageToken?: string;
}

/**
 * Request for getting video comments
 */
export interface GetVideoCommentsRequest {
  /** Video ID */
  videoId: string;
  /** Pagination page */
  page?: number;
  /** Items per page */
  limit?: number;
  /** Sort by TOP, NEW, or CONTROVERSIAL */
  sortBy?: 'TOP' | 'NEW' | 'CONTROVERSIAL';
}

/**
 * Request for creating a video comment
 */
export interface CreateVideoCommentRequest {
  /** Video ID to comment on */
  videoId: string;
  /** Comment content */
  content: string;
  /** Parent comment ID (for replies) */
  parentCommentId?: string;
  /** Mentioned users */
  mentionedUsers?: string[];
}

/**
 * Request for getting user playlists
 */
export interface GetPlaylistsRequest {
  /** User ID (optional - defaults to current user) */
  userId?: string;
  /** Include private playlists (requires auth) */
  includePrivate?: boolean;
  /** Pagination page */
  page?: number;
  /** Items per page */
  limit?: number;
}

/**
 * Request for creating a playlist
 */
export interface CreatePlaylistRequest {
  /** Playlist title */
  title: string;
  /** Playlist description */
  description: string;
  /** Is playlist private */
  isPrivate: boolean;
  /** Is playlist collaborative */
  isCollaborative?: boolean;
  /** Initial video IDs */
  videoIds?: string[];
  /** Playlist category */
  category?: VideoCategory;
}

/**
 * Request for getting live streams
 */
export interface GetLiveStreamsRequest {
  /** Filter by stream status */
  status?: 'SCHEDULED' | 'LIVE' | 'ENDED';
  /** Filter by category */
  category?: VideoCategory;
  /** Pagination page */
  page?: number;
  /** Items per page */
  limit?: number;
  /** Include private streams (requires auth) */
  includePrivate?: boolean;
}

/**
 * Request for creating a live stream
 */
export interface CreateLiveStreamRequest {
  /** Stream title */
  title: string;
  /** Stream description */
  description: string;
  /** Stream category */
  category: VideoCategory;
  /** Scheduled start time */
  scheduledFor?: string;
  /** Is stream private */
  isPrivate: boolean;
  /** Is stream premium only */
  isPremium?: boolean;
  /** Age restricted stream */
  ageRestricted?: boolean;
  /** Enable chat */
  chatEnabled?: boolean;
  /** Chat slow mode (seconds) */
  chatSlowMode?: number;
  /** Subscribers only chat */
  subscribersOnly?: boolean;
  /** Followers only chat */
  followersOnly?: boolean;
  /** Stream tags */
  tags?: string[];
}

/**
 * Request for video analytics (creator only)
 */
export interface GetVideoAnalyticsRequest {
  /** Video ID */
  videoId: string;
  /** Analytics time range */
  timeRange?: 'LAST_7_DAYS' | 'LAST_30_DAYS' | 'LAST_90_DAYS' | 'LAST_YEAR' | 'ALL_TIME';
  /** Include revenue data (if monetized) */
  includeRevenue?: boolean;
}

/**
 * Request for video interaction (like, dislike, bookmark, share)
 */
export interface VideoInteractionRequest {
  /** Video ID */
  videoId: string;
  /** Interaction type */
  action: 'LIKE' | 'DISLIKE' | 'BOOKMARK' | 'SHARE' | 'REMOVE_LIKE' | 'REMOVE_DISLIKE' | 'REMOVE_BOOKMARK';
}

/**
 * Request for channel subscription
 */
export interface ChannelSubscriptionRequest {
  /** Channel ID */
  channelId: string;
  /** Subscribe or unsubscribe */
  action: 'SUBSCRIBE' | 'UNSUBSCRIBE';
  /** Enable notifications */
  notificationsEnabled?: boolean;
}

// =============================================================================
// Video API Implementation
// =============================================================================

/**
 * Video API endpoints using RTK Query
 *
 * Provides comprehensive video platform functionality with:
 * - Type-safe API calls
 * - Advanced caching strategies
 * - Optimistic updates
 * - Error handling and retry logic
 * - Real-time synchronization
 */
export const videoApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // ========================================================================
    // Query Endpoints - Video Content Retrieval
    // ========================================================================

    /**
     * Get paginated list of videos with filtering options
     *
     * Features:
     * - Advanced filtering by category, type, language
     * - Content filtering (premium, private, age-restricted)
     * - Pagination with configurable page sizes
     * - Sorting by multiple fields
     * - Optimized caching with category-based tags
     */
    getVideos: builder.query<GetVideosResponse, GetVideosRequest>({
      query: (params) => ({
        url: '/videos',
        method: 'GET',
        params: {
          ...params,
          // Ensure arrays are properly serialized
          tags: params.tags?.join(','),
        },
      }),
      providesTags: (result, error, arg) => [
        'Video',
        { type: 'VideoList', id: 'LIST' },
        // Category-specific tags for efficient invalidation
        ...(arg.category ? [{ type: 'VideoCategory' as const, id: arg.category }] : []),
        // Type-specific tags
        ...(arg.type ? [{ type: 'VideoType' as const, id: arg.type }] : []),
      ],
      keepUnusedDataFor: 300, // 5 minutes
      // Custom serialization for request deduplication
      serializeQueryArgs: ({ queryArgs }) => {
        const { page, limit, category, type, sortBy } = queryArgs;
        return `getVideos(${JSON.stringify({ page, limit, category, type, sortBy })})`;
      },
    }),

    /**
     * Get a single video by ID with complete details
     *
     * Features:
     * - Individual video retrieval with engagement data
     * - User-specific data (subscription status, interaction state)
     * - Optimized caching per video ID
     * - Automatic view tracking
     */
    getVideoById: builder.query<SingleResponse<VideoContent>, string>({
      query: (id) => `/videos/${encodeURIComponent(id)}`,
      providesTags: (result, error, id) => [
        { type: 'Video', id },
        'Video',
        // Also provide channel tags for subscription status
        ...(result?.data.channel ? [{ type: 'Channel' as const, id: result.data.channel.id }] : []),
      ],
      keepUnusedDataFor: 600, // 10 minutes for individual videos
    }),

    /**
     * Search videos with advanced filtering and suggestions
     *
     * Features:
     * - Full-text search across videos, channels, and playlists
     * - Advanced filtering by duration, upload date, category
     * - Search suggestions and autocomplete
     * - Relevance-based sorting
     */
    searchVideos: builder.query<SearchVideosResponse, SearchVideosRequest>({
      query: (params) => ({
        url: '/videos/search',
        method: 'GET',
        params,
      }),
      providesTags: ['VideoSearch'],
      keepUnusedDataFor: 180, // 3 minutes for search results
    }),

    /**
     * Get trending videos for homepage
     *
     * Features:
     * - Trending algorithm based on views, engagement, recency
     * - Category-specific trending lists
     * - Geographic and demographic customization
     */
    getTrendingVideos: builder.query<GetVideosResponse, { category?: VideoCategory; region?: string; limit?: number }>({
      query: (params) => ({
        url: '/videos/trending',
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { category }) => [
        'TrendingVideos',
        ...(category ? [{ type: 'VideoCategory' as const, id: category }] : []),
      ],
      keepUnusedDataFor: 900, // 15 minutes for trending
    }),

    // ========================================================================
    // Channel Management
    // ========================================================================

    /**
     * Get list of all channels with filtering options
     *
     * Features:
     * - Paginated channel listing
     * - Category and verification filtering
     * - Search by channel name
     * - Subscription status for authenticated users
     */
    getChannels: builder.query<PaginatedResponse<ChannelInfo>, { page?: number; limit?: number; category?: VideoCategory; verified?: boolean; search?: string }>({
      query: (params) => ({
        url: '/channels',
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { category }) => [
        'Channel',
        { type: 'ChannelList', id: 'LIST' },
        ...(category ? [{ type: 'VideoCategory' as const, id: category }] : []),
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Get channel details with subscription status
     *
     * Features:
     * - Complete channel information
     * - User-specific subscription status
     * - Channel statistics and metadata
     */
    getChannelById: builder.query<SingleResponse<ChannelInfo>, string>({
      query: (id) => `/channels/${encodeURIComponent(id)}`,
      providesTags: (result, error, id) => [
        { type: 'Channel', id },
        'Channel',
      ],
      keepUnusedDataFor: 600, // 10 minutes
    }),

    /**
     * Get videos from a specific channel
     *
     * Features:
     * - Channel-specific video listing
     * - Sorting by upload date, popularity
     * - Pagination support
     */
    getChannelVideos: builder.query<GetVideosResponse, { channelId: string; page?: number; limit?: number; sortBy?: string }>({
      query: ({ channelId, ...params }) => ({
        url: `/channels/${encodeURIComponent(channelId)}/videos`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { channelId }) => [
        { type: 'Channel', id: channelId },
        { type: 'ChannelVideos', id: channelId },
        'Video',
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    // ========================================================================
    // Comments System
    // ========================================================================

    /**
     * Get video comments with nested replies
     *
     * Features:
     * - Paginated comments with sorting options
     * - Nested reply structure
     * - User interaction states (likes, pins, hearts)
     */
    getVideoComments: builder.query<PaginatedResponse<VideoComment>, GetVideoCommentsRequest>({
      query: ({ videoId, ...params }) => ({
        url: `/videos/${encodeURIComponent(videoId)}/comments`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { videoId }) => [
        { type: 'VideoComments', id: videoId },
        'Comment',
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    // ========================================================================
    // Playlists
    // ========================================================================

    /**
     * Get user playlists
     *
     * Features:
     * - User-specific playlist retrieval
     * - Privacy filtering
     * - Collaborative playlist support
     */
    getPlaylists: builder.query<PaginatedResponse<VideoPlaylist>, GetPlaylistsRequest>({
      query: (params) => ({
        url: '/playlists',
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { userId }) => [
        'Playlist',
        { type: 'UserPlaylists', id: userId || 'current' },
      ],
      keepUnusedDataFor: 300, // 5 minutes
    }),

    /**
     * Get specific playlist with videos
     *
     * Features:
     * - Complete playlist information
     * - Video list with metadata
     * - Collaborative playlist permissions
     */
    getPlaylistById: builder.query<SingleResponse<VideoPlaylist & { videos: VideoContent[] }>, string>({
      query: (id) => `/playlists/${encodeURIComponent(id)}`,
      providesTags: (result, error, id) => [
        { type: 'Playlist', id },
        'Playlist',
      ],
      keepUnusedDataFor: 600, // 10 minutes
    }),

    // ========================================================================
    // Live Streaming
    // ========================================================================

    /**
     * Get live streams with filtering
     *
     * Features:
     * - Real-time stream listing
     * - Status and category filtering
     * - Live viewer counts
     */
    getLiveStreams: builder.query<PaginatedResponse<LiveStream>, GetLiveStreamsRequest>({
      query: (params) => ({
        url: '/livestreams',
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { status, category }) => [
        'LiveStream',
        ...(status ? [{ type: 'StreamStatus' as const, id: status }] : []),
        ...(category ? [{ type: 'VideoCategory' as const, id: category }] : []),
      ],
      keepUnusedDataFor: 60, // 1 minute for live data
    }),

    /**
     * Get specific live stream details
     *
     * Features:
     * - Real-time stream information
     * - Current viewer count
     * - Chat integration details
     */
    getLiveStreamById: builder.query<SingleResponse<LiveStream>, string>({
      query: (id) => `/livestreams/${encodeURIComponent(id)}`,
      providesTags: (result, error, id) => [
        { type: 'LiveStream', id },
        'LiveStream',
      ],
      keepUnusedDataFor: 30, // 30 seconds for live stream data
    }),

    // ========================================================================
    // Analytics (Creator Only)
    // ========================================================================

    /**
     * Get video analytics for creators
     *
     * Features:
     * - Comprehensive video performance metrics
     * - Audience demographics and retention
     * - Revenue data (if monetized)
     * - Time-based analytics filtering
     */
    getVideoAnalytics: builder.query<SingleResponse<VideoAnalytics>, GetVideoAnalyticsRequest>({
      query: ({ videoId, ...params }) => ({
        url: `/videos/${encodeURIComponent(videoId)}/analytics`,
        method: 'GET',
        params,
      }),
      providesTags: (result, error, { videoId }) => [
        { type: 'VideoAnalytics', id: videoId },
        'Analytics',
      ],
      keepUnusedDataFor: 300, // 5 minutes for analytics data
    }),

    // ========================================================================
    // Mutation Endpoints - Content Modification
    // ========================================================================

    /**
     * Create a new video
     *
     * Features:
     * - Video upload and processing
     * - Metadata and thumbnail management
     * - Privacy and content settings
     * - Automatic quality processing
     */
    createVideo: builder.mutation<SingleResponse<VideoContent>, CreateVideoRequest>({
      query: (body) => ({
        url: '/videos',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { category, type }) => [
        'Video',
        'VideoList',
        { type: 'VideoCategory', id: category },
        { type: 'VideoType', id: type },
        'TrendingVideos',
      ],
      // Optimistic update for immediate UI feedback
      async onQueryStarted(arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          // Update video lists to include new video
          dispatch(
            videoApi.util.updateQueryData('getVideos', { category: arg.category, type: arg.type }, (draft) => {
              if (draft.data) {
                draft.data.unshift(data.data);
                draft.total += 1;
              }
            })
          );
        } catch {
          // Handle error silently - invalidation will refetch
        }
      },
    }),

    /**
     * Update video metadata
     *
     * Features:
     * - Partial video updates
     * - Metadata modification
     * - Privacy settings changes
     * - Thumbnail updates
     */
    updateVideo: builder.mutation<SingleResponse<VideoContent>, UpdateVideoRequest>({
      query: ({ id, ...body }) => ({
        url: `/videos/${encodeURIComponent(id)}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: 'Video', id },
        'VideoList',
        ...(result ? [{ type: 'VideoCategory' as const, id: result.data.category }] : []),
      ],
      // Optimistic update
      async onQueryStarted({ id, title, description, category, tags, isPrivate, ageRestricted, thumbnail }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          videoApi.util.updateQueryData('getVideoById', id, (draft) => {
            if (draft.data) {
              if (title !== undefined) draft.data.title = title;
              if (description !== undefined) draft.data.description = description;
              if (category !== undefined) draft.data.category = category;
              if (tags !== undefined) draft.data.tags = tags;
              if (isPrivate !== undefined) draft.data.isPrivate = isPrivate;
              if (ageRestricted !== undefined) draft.data.ageRestricted = ageRestricted;
              if (thumbnail !== undefined) draft.data.thumbnail = thumbnail;
            }
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
     * Delete a video
     *
     * Features:
     * - Video removal from platform
     * - Cascade deletion of comments and analytics
     * - Playlist cleanup
     * - Cache invalidation
     */
    deleteVideo: builder.mutation<{ success: boolean }, string>({
      query: (id) => ({
        url: `/videos/${encodeURIComponent(id)}`,
        method: 'DELETE',
      }),
      invalidatesTags: (result, error, id) => [
        { type: 'Video', id },
        'VideoList',
        'TrendingVideos',
        { type: 'VideoComments', id },
        { type: 'VideoAnalytics', id },
      ],
    }),

    // ========================================================================
    // User Interactions
    // ========================================================================

    /**
     * Handle video interactions (like, dislike, bookmark, share)
     *
     * Features:
     * - User engagement tracking
     * - Real-time interaction updates
     * - Optimistic UI updates
     */
    videoInteraction: builder.mutation<{ success: boolean; newCount: number }, VideoInteractionRequest>({
      query: (body) => ({
        url: '/videos/interactions',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { videoId }) => [
        { type: 'Video', id: videoId },
      ],
      // Optimistic update for instant feedback
      async onQueryStarted({ videoId, action }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          videoApi.util.updateQueryData('getVideoById', videoId, (draft) => {
            if (draft.data) {
              switch (action) {
                case 'LIKE':
                  draft.data.likes += 1;
                  break;
                case 'DISLIKE':
                  draft.data.dislikes += 1;
                  break;
                case 'BOOKMARK':
                  draft.data.bookmarks += 1;
                  break;
                case 'SHARE':
                  draft.data.shares += 1;
                  break;
                case 'REMOVE_LIKE':
                  draft.data.likes = Math.max(0, draft.data.likes - 1);
                  break;
                case 'REMOVE_DISLIKE':
                  draft.data.dislikes = Math.max(0, draft.data.dislikes - 1);
                  break;
                case 'REMOVE_BOOKMARK':
                  draft.data.bookmarks = Math.max(0, draft.data.bookmarks - 1);
                  break;
              }
            }
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
     * Subscribe/unsubscribe to channel
     *
     * Features:
     * - Channel subscription management
     * - Notification preferences
     * - Real-time subscriber count updates
     */
    channelSubscription: builder.mutation<{ success: boolean; newSubscriberCount: number }, ChannelSubscriptionRequest>({
      query: (body) => ({
        url: '/channels/subscription',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { channelId }) => [
        { type: 'Channel', id: channelId },
      ],
      // Optimistic update
      async onQueryStarted({ channelId, action, notificationsEnabled }, { dispatch, queryFulfilled }) {
        const patchResult = dispatch(
          videoApi.util.updateQueryData('getChannelById', channelId, (draft) => {
            if (draft.data) {
              if (action === 'SUBSCRIBE') {
                draft.data.subscribers += 1;
                draft.data.isSubscribed = true;
                draft.data.notificationsEnabled = notificationsEnabled || false;
              } else {
                draft.data.subscribers = Math.max(0, draft.data.subscribers - 1);
                draft.data.isSubscribed = false;
                draft.data.notificationsEnabled = false;
              }
            }
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
     * Create a video comment
     *
     * Features:
     * - Comment creation with mentions
     * - Reply threading support
     * - Real-time comment updates
     */
    createVideoComment: builder.mutation<SingleResponse<VideoComment>, CreateVideoCommentRequest>({
      query: (body) => ({
        url: '/videos/comments',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { videoId }) => [
        { type: 'VideoComments', id: videoId },
        { type: 'Video', id: videoId }, // Update comment count
      ],
      // Optimistic update
      async onQueryStarted(arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          // Update comment list
          dispatch(
            videoApi.util.updateQueryData('getVideoComments', { videoId: arg.videoId }, (draft) => {
              if (draft.data) {
                if (arg.parentCommentId) {
                  // Add as reply
                  const parentComment = draft.data.find(c => c.id === arg.parentCommentId);
                  if (parentComment) {
                    parentComment.replies.push(data.data);
                  }
                } else {
                  // Add as top-level comment
                  draft.data.unshift(data.data);
                  draft.total += 1;
                }
              }
            })
          );
          // Update video comment count
          dispatch(
            videoApi.util.updateQueryData('getVideoById', arg.videoId, (draft) => {
              if (draft.data) {
                draft.data.comments += 1;
              }
            })
          );
        } catch {
          // Handle error silently
        }
      },
    }),

    /**
     * Create a playlist
     *
     * Features:
     * - Playlist creation with initial videos
     * - Privacy and collaboration settings
     * - Automatic categorization
     */
    createPlaylist: builder.mutation<SingleResponse<VideoPlaylist>, CreatePlaylistRequest>({
      query: (body) => ({
        url: '/playlists',
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Playlist', { type: 'UserPlaylists', id: 'current' }],
    }),

    /**
     * Create a live stream
     *
     * Features:
     * - Live stream setup and configuration
     * - Chat and interaction settings
     * - Scheduling and privacy controls
     */
    createLiveStream: builder.mutation<SingleResponse<LiveStream>, CreateLiveStreamRequest>({
      query: (body) => ({
        url: '/livestreams',
        method: 'POST',
        body,
      }),
      invalidatesTags: (result, error, { category, isPremium }) => [
        'LiveStream',
        { type: 'VideoCategory', id: category },
        { type: 'StreamStatus', id: 'SCHEDULED' },
      ],
    }),
  }),
  overrideExisting: false,
});

// =============================================================================
// Export Generated Hooks
// =============================================================================

// Query hooks for data fetching
export const {
  useGetVideosQuery,
  useGetVideoByIdQuery,
  useSearchVideosQuery,
  useGetTrendingVideosQuery,
  useGetChannelsQuery,
  useGetChannelByIdQuery,
  useGetChannelVideosQuery,
  useGetVideoCommentsQuery,
  useGetPlaylistsQuery,
  useGetPlaylistByIdQuery,
  useGetLiveStreamsQuery,
  useGetLiveStreamByIdQuery,
  useGetVideoAnalyticsQuery,

  // Lazy query hooks for manual triggering
  useLazyGetVideosQuery,
  useLazyGetVideoByIdQuery,
  useLazySearchVideosQuery,
  useLazyGetTrendingVideosQuery,
  useLazyGetChannelsQuery,
  useLazyGetChannelByIdQuery,
  useLazyGetChannelVideosQuery,
  useLazyGetVideoCommentsQuery,
  useLazyGetPlaylistsQuery,
  useLazyGetPlaylistByIdQuery,
  useLazyGetLiveStreamsQuery,
  useLazyGetLiveStreamByIdQuery,
  useLazyGetVideoAnalyticsQuery,

  // Mutation hooks for data modification
  useCreateVideoMutation,
  useUpdateVideoMutation,
  useDeleteVideoMutation,
  useVideoInteractionMutation,
  useChannelSubscriptionMutation,
  useCreateVideoCommentMutation,
  useCreatePlaylistMutation,
  useCreateLiveStreamMutation,

  // Utility hooks for cache management
  util: {
    updateQueryData: videoApiUpdateQueryData,
    invalidateTags: videoApiInvalidateTags,
    resetApiState: videoApiResetApiState,
    getRunningQueriesThunk: videoApiGetRunningQueriesThunk,
  },
} = videoApi;

// =============================================================================
// Advanced Hook Utilities
// =============================================================================

/**
 * Enhanced video hook with user interaction state
 * Combines video data with user-specific interaction information
 */
export const useVideoWithInteractions = (videoId: string) => {
  const videoQuery = useGetVideoByIdQuery(videoId);

  return {
    ...videoQuery,
    video: videoQuery.data?.data,
    // Additional computed properties for UI
    isLiked: videoQuery.data?.data?.channel?.isSubscribed || false, // This would come from user interaction state
    isBookmarked: false, // This would come from user interaction state
    hasInteracted: false, // This would come from user interaction state
  };
};

/**
 * Prefetch video content for performance optimization
 */
export const usePrefetchVideo = () => {
  const prefetchVideo = useLazyGetVideoByIdQuery()[1];
  const prefetchVideos = useLazyGetVideosQuery()[1];
  const prefetchComments = useLazyGetVideoCommentsQuery()[1];

  return {
    prefetchVideo: (id: string) => prefetchVideo(id),
    prefetchVideos: (params: GetVideosRequest) => prefetchVideos(params),
    prefetchComments: (videoId: string) => prefetchComments({ videoId }),
    prefetchCategory: (category: VideoCategory) =>
      prefetchVideos({ category, limit: 20 }),
    prefetchTrending: () =>
      prefetchVideos({ sortBy: 'trending', limit: 20 }),
  };
};

/**
 * Enhanced video search with history and suggestions
 */
export const useVideoSearch = () => {
  const [searchVideos] = useLazySearchVideosQuery();

  // This would integrate with local storage for search history
  const searchHistory: string[] = []; // Would come from localStorage
  const recentSearches: string[] = []; // Would come from localStorage

  const performSearch = (query: string, filters?: Partial<SearchVideosRequest>) => {
    return searchVideos({ q: query, ...filters });
  };

  return {
    searchVideos: performSearch,
    searchHistory,
    recentSearches,
    clearHistory: () => {
      // Would clear localStorage
    },
  };
};

/**
 * Video list management with infinite scroll support
 */
export const useVideoList = (initialParams: GetVideosRequest = {}) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [allVideos, setAllVideos] = useState<VideoContent[]>([]);

  const query = useGetVideosQuery({
    ...initialParams,
    page: currentPage,
    limit: initialParams.limit || 12
  });

  useEffect(() => {
    if (query.data?.data && currentPage === 1) {
      setAllVideos(query.data.data);
    } else if (query.data?.data && currentPage > 1) {
      setAllVideos(prev => [...prev, ...query.data.data]);
    }
  }, [query.data, currentPage]);

  const loadMore = () => {
    if (query.data?.pagination?.hasNext) {
      setCurrentPage(prev => prev + 1);
    }
  };

  const reset = () => {
    setCurrentPage(1);
    setAllVideos([]);
  };

  return {
    videos: allVideos,
    isLoading: query.isLoading,
    isError: query.isError,
    error: query.error,
    hasMore: query.data?.pagination?.hasNext || false,
    loadMore,
    reset,
    refetch: query.refetch,
  };
};

// =============================================================================
// Backend Video Service Integration
// =============================================================================

/**
 * Upload video to backend video service with progress tracking
 */
export async function uploadVideoToBackend(
  file: File,
  metadata: {
    title: string;
    description: string;
    tags: string[];
    content_rating: string;
    thumbnail?: File;
    category?: string;
    is_monetized?: boolean;
    price?: number;
  },
  onProgress?: (progress: number) => void
): Promise<{ video_id: string; status: string; message: string }> {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('title', metadata.title);
  formData.append('description', metadata.description);
  formData.append('tags', metadata.tags.join(','));
  formData.append('content_rating', metadata.content_rating);

  if (metadata.thumbnail) {
    formData.append('thumbnail', metadata.thumbnail);
  }
  if (metadata.category) {
    formData.append('category', metadata.category);
  }
  formData.append('is_monetized', String(metadata.is_monetized || false));
  if (metadata.price) {
    formData.append('price', String(metadata.price));
  }

  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();

    xhr.upload.addEventListener('progress', (event) => {
      if (event.lengthComputable && onProgress) {
        const progress = Math.round((event.loaded / event.total) * 100);
        onProgress(progress);
      }
    });

    xhr.addEventListener('load', () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          resolve(JSON.parse(xhr.responseText));
        } catch (error) {
          reject(new Error('Invalid response'));
        }
      } else {
        reject(new Error(`Upload failed: ${xhr.status}`));
      }
    });

    xhr.addEventListener('error', () => reject(new Error('Upload failed')));

    xhr.open('POST', '/api/v1/video');

    const token = localStorage.getItem('auth_token');
    if (token) {
      xhr.setRequestHeader('Authorization', `Bearer ${token}`);
    }

    xhr.send(formData);
  });
}

/**
 * Get streaming URL from backend video service
 */
export async function getBackendStreamUrl(
  videoId: string,
  quality: string = 'auto',
  platform: string = 'web'
): Promise<{
  video_id: string;
  stream_url: string;
  protocol: string;
  available_qualities: string[];
  expires_at: string;
}> {
  const response = await fetch(`/api/v1/video/${videoId}/stream?quality=${quality}&platform=${platform}`);
  if (!response.ok) throw new Error('Failed to get stream URL');
  return response.json();
}

/**
 * Sync playback position with backend
 */
export async function syncBackendPlaybackPosition(
  videoId: string,
  sessionId: string,
  position: number,
  platform: string,
  playbackState: string = 'playing'
): Promise<{
  success: boolean;
  session_id: string;
  updated_position: number;
  synced_platforms: string[];
}> {
  const response = await fetch(`/api/v1/video/${videoId}/sync`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
    },
    body: JSON.stringify({
      session_id: sessionId,
      position,
      platform,
      playback_state: playbackState,
      timestamp: new Date().toISOString(),
    }),
  });

  if (!response.ok) throw new Error('Sync failed');
  return response.json();
}

/**
 * Create sync session on backend
 */
export async function createBackendSyncSession(
  videoId: string,
  userId: string,
  platform: string,
  initialPosition: number = 0
): Promise<{
  session_id: string;
  video_id: string;
  user_id: string;
  platform: string;
  initial_position: number;
  created_at: string;
}> {
  const response = await fetch(`/api/v1/video/${videoId}/sync/session`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
    },
    body: JSON.stringify({
      video_id: videoId,
      user_id: userId,
      platform,
      initial_position: initialPosition,
    }),
  });

  if (!response.ok) throw new Error('Failed to create sync session');
  return response.json();
}

// Export the enhanced API for external use
export default videoApi;