/**
 * Social Service TypeScript Models
 *
 * Complete type definitions for social features aligned with backend Go models.
 * Supports the comprehensive social service API with proper typing and validation.
 */

export interface SocialProfile {
  // Core User Fields (from shared.User)
  id: string;
  username?: string;
  phone?: string;
  phone_number?: string;
  email?: string;
  name: string;
  display_name?: string;
  first_name?: string;
  last_name?: string;
  avatar?: string;
  country: string;
  country_code?: string;
  locale: string;
  language?: string;
  timezone?: string;
  timezone_alias?: string;
  bio?: string;
  date_of_birth?: string;
  gender?: string;
  is_active: boolean;
  kyc_status: string;
  kyc_tier: number;
  status: string;
  last_seen?: string;
  last_active_at?: string;
  is_verified: boolean;
  is_email_verified: boolean;
  is_phone_verified: boolean;
  pref_theme?: string;
  pref_language?: string;
  pref_notifications_email: boolean;
  pref_notifications_push: boolean;
  pref_privacy_level?: string;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;

  // Social-Specific Fields
  interests?: string[];
  socialLinks?: Record<string, any>;
  socialPreferences?: Record<string, any>;
  followersCount: number;
  followingCount: number;
  postsCount: number;
  isSocialVerified: boolean;
  socialCreatedAt: string;
  socialUpdatedAt: string;
}

export interface Post {
  id: string;
  authorId: string;
  communityId?: string;
  content: string;
  type: 'text' | 'image' | 'video' | 'link' | 'poll';
  metadata?: Record<string, any>;
  tags?: string[];
  visibility: 'public' | 'members' | 'private' | 'followers';
  mediaUrls?: string[];
  linkPreview?: Record<string, any>;

  // Interaction counts
  likesCount: number;
  commentsCount: number;
  sharesCount: number;
  reactionsCount: number;
  viewsCount: number;

  // Status flags
  isEdited: boolean;
  isPinned: boolean;
  isDeleted: boolean;
  isTrending: boolean;

  // Timestamps
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;

  // Enhanced fields for UI (populated by API)
  author?: SocialProfile;
  userReaction?: Reaction;
  isLiked?: boolean;
  isBookmarked?: boolean;
  comments?: Comment[];
}

export interface Comment {
  id: string;
  postId: string;
  authorId: string;
  parentId?: string;
  content: string;
  metadata?: Record<string, any>;

  // Interaction counts
  likesCount: number;
  repliesCount: number;
  reactionsCount: number;

  // Status flags
  isEdited: boolean;
  isDeleted: boolean;

  // Timestamps
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;

  // Enhanced fields for UI
  author?: SocialProfile;
  userReaction?: Reaction;
  replies?: Comment[];
}

export interface Reaction {
  id: string;
  userId: string;
  targetId: string;
  targetType: 'post' | 'comment';
  type: 'like' | 'love' | 'laugh' | 'angry' | 'sad' | 'wow';
  createdAt: string;
  updatedAt: string;
}

export interface UserRelationship {
  id: string;
  followerId: string;
  followingId: string;
  status: 'pending' | 'active' | 'blocked';
  source: 'discovery' | 'suggestion' | 'search' | 'follow_back' | 'manual';
  isMutual: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface SocialFeed {
  userId: string;
  posts: Post[];
  algorithm: 'chronological' | 'personalized' | 'trending';
  region?: string;
  cursor?: string;
  hasMore: boolean;
  updatedAt: string;
}

export interface TrendingContent {
  region: string;
  timeframe: '1h' | '24h' | '7d';
  topics: string[];
  posts: Post[];
  hashtags: string[];
  metrics: Record<string, any>;
  updatedAt: string;
}

export interface UserActivity {
  id: string;
  userId: string;
  activityType: 'post' | 'comment' | 'reaction' | 'follow' | 'share';
  targetId?: string;
  targetType?: 'post' | 'comment' | 'user' | 'community';
  metadata?: Record<string, any>;
  region: string;
  platform: 'web' | 'mobile' | 'api' | 'kmp_android' | 'kmp_ios';
  createdAt: string;
}

export interface UserAnalytics {
  userId: string;
  period: '1h' | '24h' | '7d' | '30d';
  followers: Record<string, any>;
  following: Record<string, any>;
  engagement: Record<string, any>;
  reach: Record<string, any>;
  growth: Record<string, any>;
  demographics: Record<string, any>;
  updatedAt: string;
  isRealTime: boolean;
  cacheExpiry?: string;
}

// Request Types
export interface UpdateSocialProfileRequest {
  displayName?: string;
  bio?: string;
  avatar?: string;
  interests?: string[];
  socialLinks?: Record<string, any>;
  socialPreferences?: Record<string, any>;
}

export interface CreatePostRequest {
  communityId?: string;
  content: string;
  type: 'text' | 'image' | 'video' | 'link' | 'poll';
  metadata?: Record<string, any>;
  tags?: string[];
  visibility: 'public' | 'members' | 'private' | 'followers';
  mediaUrls?: string[];
}

export interface UpdatePostRequest {
  content?: string;
  tags?: string[];
  metadata?: Record<string, any>;
  isPinned?: boolean;
}

export interface CreateCommentRequest {
  postId: string;
  content: string;
  parentId?: string;
  metadata?: Record<string, any>;
}

export interface CreateReactionRequest {
  targetId: string;
  targetType: 'post' | 'comment';
  type: 'like' | 'love' | 'laugh' | 'angry' | 'sad' | 'wow';
}

export interface FollowRequest {
  followerId: string;
  followingId: string;
  source: 'discovery' | 'suggestion' | 'search' | 'follow_back' | 'manual';
}

export interface UserDiscoveryRequest {
  region?: 'TH' | 'SG' | 'ID' | 'MY' | 'PH' | 'VN';
  interests?: string[];
  limit?: number;
  offset?: number;
}

export interface SocialFeedRequest {
  algorithm?: 'chronological' | 'personalized' | 'trending';
  limit?: number;
  cursor?: string;
  region?: 'TH' | 'SG' | 'ID' | 'MY' | 'PH' | 'VN';
}

export interface TrendingRequest {
  region?: 'TH' | 'SG' | 'ID' | 'MY' | 'PH' | 'VN';
  timeframe?: '1h' | '24h' | '7d';
  category?: string;
  limit?: number;
}

export interface ShareRequest {
  contentId: string;
  contentType: 'post' | 'comment' | 'community';
  platform: 'internal' | 'external';
  message?: string;
  privacy?: string;
  metadata?: Record<string, any>;
}

// Response Types with Pagination
export interface FollowersResponse {
  followers: SocialProfile[];
  total: number;
  hasMore: boolean;
  cursor?: string;
}

export interface FollowingResponse {
  following: SocialProfile[];
  total: number;
  hasMore: boolean;
  cursor?: string;
}

export interface PaginatedPostsResponse {
  posts: Post[];
  total: number;
  hasMore: boolean;
  cursor?: string;
}

// Legacy support for existing components
export interface LegacyPost {
  id: string;
  author: {
    name: string;
    avatar?: string;
    verified?: boolean;
    type?: 'user' | 'channel';
  };
  content: string;
  images?: string[];
  timestamp: string;
  likes: number;
  comments: number;
  shares: number;
  isLiked?: boolean;
  location?: string;
  tags?: string[];
  type?: 'text' | 'image' | 'video' | 'live';
  source?: 'trending' | 'sponsored' | 'interest' | 'following';
  product?: any;
  liveData?: any;
}

// Utility types for social features
export type SocialNotificationType =
  | 'follow'
  | 'unfollow'
  | 'like'
  | 'comment'
  | 'mention'
  | 'share'
  | 'post_trending';

export interface SocialNotification {
  id: string;
  type: SocialNotificationType;
  actorId: string;
  targetId: string;
  targetType: 'post' | 'comment' | 'user';
  message: string;
  isRead: boolean;
  createdAt: string;
  actor?: SocialProfile;
}

// Real-time updates
export interface SocialRealtimeEvent {
  type: 'post_created' | 'post_updated' | 'post_deleted' | 'reaction_added' | 'comment_added' | 'follow_event';
  data: Post | Comment | Reaction | UserRelationship;
  userId: string;
  timestamp: string;
}