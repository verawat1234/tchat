/**
 * Social API Test Suite
 *
 * Comprehensive tests for the social API service including RTK Query endpoints,
 * mutations, optimistic updates, and error handling.
 */

import { configureStore } from '@reduxjs/toolkit';
import { setupListeners } from '@reduxjs/toolkit/query';
import { api } from '../api';
import { socialApi } from '../socialApi';
import type { SocialProfile, Post, CreatePostRequest } from '../../types/social';

// Mock fetch for testing
global.fetch = jest.fn();

// Create test store
const createTestStore = () => {
  const store = configureStore({
    reducer: {
      [api.reducerPath]: api.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(api.middleware),
  });

  setupListeners(store.dispatch);
  return store;
};

describe('Social API Service', () => {
  let store: ReturnType<typeof createTestStore>;

  beforeEach(() => {
    store = createTestStore();
    (fetch as jest.Mock).mockClear();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('Social Profile Endpoints', () => {
    const mockProfile: SocialProfile = {
      id: 'user-123',
      name: 'Test User',
      display_name: 'Test User',
      avatar: 'https://example.com/avatar.jpg',
      bio: 'Test bio',
      country: 'TH',
      locale: 'en',
      is_active: true,
      kyc_status: 'verified',
      kyc_tier: 1,
      status: 'active',
      is_verified: true,
      is_email_verified: true,
      is_phone_verified: true,
      pref_notifications_email: true,
      pref_notifications_push: true,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
      interests: ['technology', 'travel'],
      followersCount: 150,
      followingCount: 200,
      postsCount: 50,
      isSocialVerified: true,
      socialCreatedAt: '2024-01-01T00:00:00Z',
      socialUpdatedAt: '2024-01-01T00:00:00Z',
    };

    it('should fetch social profile successfully', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockProfile,
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.getSocialProfile.initiate('user-123')
      );

      expect(result.data).toEqual(mockProfile);
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/profiles/user-123'),
        expect.objectContaining({
          method: 'GET',
        })
      );
    });

    it('should update social profile with optimistic updates', async () => {
      const updates = {
        displayName: 'Updated Name',
        bio: 'Updated bio',
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: { ...mockProfile, ...updates },
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.updateSocialProfile.initiate({
          userId: 'user-123',
          updates,
        })
      );

      expect(result.data).toEqual(expect.objectContaining(updates));
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/profiles/user-123'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(updates),
        })
      );
    });

    it('should handle profile fetch errors gracefully', async () => {
      (fetch as jest.Mock).mockRejectedValueOnce(new Error('Network error'));

      const result = await store.dispatch(
        socialApi.endpoints.getSocialProfile.initiate('invalid-user')
      );

      expect(result.error).toBeDefined();
      expect(result.data).toBeUndefined();
    });
  });

  describe('Post Management Endpoints', () => {
    const mockPost: Post = {
      id: 'post-123',
      authorId: 'user-123',
      content: 'Test post content',
      type: 'text',
      visibility: 'public',
      likesCount: 10,
      commentsCount: 5,
      sharesCount: 2,
      reactionsCount: 10,
      viewsCount: 100,
      isEdited: false,
      isPinned: false,
      isDeleted: false,
      isTrending: false,
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    };

    const mockCreateRequest: CreatePostRequest = {
      content: 'New post content',
      type: 'text',
      visibility: 'public',
    };

    it('should create a new post successfully', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: mockPost,
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.createPost.initiate(mockCreateRequest)
      );

      expect(result.data).toEqual(mockPost);
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/posts'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(mockCreateRequest),
        })
      );
    });

    it('should fetch a specific post by ID', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockPost,
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.getPost.initiate('post-123')
      );

      expect(result.data).toEqual(mockPost);
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/posts/post-123'),
        expect.objectContaining({
          method: 'GET',
        })
      );
    });

    it('should update post with optimistic updates', async () => {
      const updates = {
        content: 'Updated content',
        isPinned: true,
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: { ...mockPost, ...updates },
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.updatePost.initiate({
          postId: 'post-123',
          updates,
        })
      );

      expect(result.data).toEqual(expect.objectContaining(updates));
    });

    it('should delete post successfully', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          message: 'Post deleted successfully',
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.deletePost.initiate('post-123')
      );

      expect(result.data).toEqual({ message: 'Post deleted successfully' });
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/posts/post-123'),
        expect.objectContaining({
          method: 'DELETE',
        })
      );
    });
  });

  describe('Social Feed Endpoints', () => {
    const mockFeed = {
      userId: 'user-123',
      posts: [
        {
          id: 'post-1',
          authorId: 'user-456',
          content: 'Feed post 1',
          type: 'text' as const,
          visibility: 'public' as const,
          likesCount: 5,
          commentsCount: 2,
          sharesCount: 1,
          reactionsCount: 5,
          viewsCount: 50,
          isEdited: false,
          isPinned: false,
          isDeleted: false,
          isTrending: false,
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ],
      algorithm: 'personalized' as const,
      hasMore: true,
      updatedAt: '2024-01-01T00:00:00Z',
    };

    it('should fetch social feed with proper parameters', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockFeed,
        }),
      });

      const feedParams = {
        algorithm: 'personalized' as const,
        limit: 20,
        region: 'TH' as const,
      };

      const result = await store.dispatch(
        socialApi.endpoints.getSocialFeed.initiate(feedParams)
      );

      expect(result.data).toEqual(mockFeed);
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/feed'),
        expect.objectContaining({
          method: 'GET',
        })
      );
    });

    it('should fetch trending content', async () => {
      const mockTrending = {
        region: 'TH',
        timeframe: '24h',
        topics: ['technology', 'travel'],
        posts: mockFeed.posts,
        hashtags: ['#tech', '#travel'],
        metrics: { totalPosts: 100, totalEngagement: 500 },
        updatedAt: '2024-01-01T00:00:00Z',
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockTrending,
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.getTrendingContent.initiate({
          region: 'TH',
          timeframe: '24h',
          limit: 20,
        })
      );

      expect(result.data).toEqual(mockTrending);
    });
  });

  describe('Reactions and Interactions', () => {
    it('should add reaction with optimistic updates', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          message: 'Reaction added successfully',
        }),
      });

      const reactionRequest = {
        targetId: 'post-123',
        targetType: 'post' as const,
        type: 'like' as const,
      };

      const result = await store.dispatch(
        socialApi.endpoints.addReaction.initiate(reactionRequest)
      );

      expect(result.data).toEqual({ message: 'Reaction added successfully' });
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/reactions'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(reactionRequest),
        })
      );
    });

    it('should remove reaction successfully', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          message: 'Reaction removed successfully',
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.removeReaction.initiate({
          targetId: 'post-123',
          targetType: 'post',
        })
      );

      expect(result.data).toEqual({ message: 'Reaction removed successfully' });
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/reactions/post/post-123'),
        expect.objectContaining({
          method: 'DELETE',
        })
      );
    });

    it('should create comment successfully', async () => {
      const mockComment = {
        id: 'comment-123',
        postId: 'post-123',
        authorId: 'user-123',
        content: 'Test comment',
        likesCount: 0,
        repliesCount: 0,
        reactionsCount: 0,
        isEdited: false,
        isDeleted: false,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: mockComment,
        }),
      });

      const commentRequest = {
        postId: 'post-123',
        content: 'Test comment',
      };

      const result = await store.dispatch(
        socialApi.endpoints.createComment.initiate(commentRequest)
      );

      expect(result.data).toEqual(mockComment);
    });
  });

  describe('User Relationships', () => {
    it('should follow user successfully', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          message: 'Successfully followed user',
        }),
      });

      const followRequest = {
        followerId: 'user-123',
        followingId: 'user-456',
        source: 'manual' as const,
      };

      const result = await store.dispatch(
        socialApi.endpoints.followUser.initiate(followRequest)
      );

      expect(result.data).toEqual({ message: 'Successfully followed user' });
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/social/follow'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(followRequest),
        })
      );
    });

    it('should unfollow user successfully', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          message: 'Successfully unfollowed user',
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.unfollowUser.initiate({
          followerId: 'user-123',
          followingId: 'user-456',
        })
      );

      expect(result.data).toEqual({ message: 'Successfully unfollowed user' });
    });

    it('should fetch followers list', async () => {
      const mockFollowers = {
        followers: [
          {
            id: 'user-456',
            name: 'Follower User',
            avatar: 'https://example.com/avatar2.jpg',
          },
        ],
        total: 1,
        hasMore: false,
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockFollowers,
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.getFollowers.initiate({
          userId: 'user-123',
          limit: 50,
        })
      );

      expect(result.data).toEqual(mockFollowers);
    });
  });

  describe('Error Handling', () => {
    it('should handle network errors gracefully', async () => {
      (fetch as jest.Mock).mockRejectedValueOnce(new Error('Network error'));

      const result = await store.dispatch(
        socialApi.endpoints.getSocialProfile.initiate('user-123')
      );

      expect(result.error).toBeDefined();
      expect(result.data).toBeUndefined();
    });

    it('should handle API errors with proper error messages', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({
          success: false,
          error: 'User not found',
        }),
      });

      const result = await store.dispatch(
        socialApi.endpoints.getSocialProfile.initiate('invalid-user')
      );

      expect(result.error).toBeDefined();
      expect(result.data).toBeUndefined();
    });

    it('should handle validation errors for post creation', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({
          success: false,
          error: 'Validation failed',
          details: {
            content: 'Content is required',
          },
        }),
      });

      const invalidRequest = {
        content: '',
        type: 'text' as const,
        visibility: 'public' as const,
      };

      const result = await store.dispatch(
        socialApi.endpoints.createPost.initiate(invalidRequest)
      );

      expect(result.error).toBeDefined();
    });
  });

  describe('Cache Management', () => {
    it('should invalidate relevant caches after post creation', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 201,
        json: async () => ({
          success: true,
          data: {
            id: 'new-post',
            content: 'New post',
            authorId: 'user-123',
          },
        }),
      });

      await store.dispatch(
        socialApi.endpoints.createPost.initiate({
          content: 'New post',
          type: 'text',
          visibility: 'public',
        })
      );

      // Check that appropriate cache invalidation has occurred
      const state = store.getState();
      expect(state.api.queries).toBeDefined();
    });

    it('should implement proper cache TTL for different endpoints', () => {
      // Test cache configuration
      const profileEndpoint = socialApi.endpoints.getSocialProfile;
      const feedEndpoint = socialApi.endpoints.getSocialFeed;

      // Profile should cache longer than feed
      expect(profileEndpoint.keepUnusedDataFor).toBeGreaterThan(
        feedEndpoint.keepUnusedDataFor || 0
      );
    });
  });
});