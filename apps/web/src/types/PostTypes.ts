/**
 * Post Types - Web Platform
 *
 * Unified Post Type System aligned with mobile Kotlin implementation
 * Cross-platform consistency with 42 post types and comprehensive data structures
 */

/**
 * Unified Post Type System - 42 Post Types
 * Aligned with mobile platform Post.kt for consistent cross-platform experience
 */
export enum PostType {
  // Core Content Types (8)
  TEXT = 'TEXT',                    // Simple status updates
  IMAGE = 'IMAGE',                  // Single/multiple photos
  VIDEO = 'VIDEO',                  // Video posts
  AUDIO = 'AUDIO',                  // Voice notes, music
  LINK_SHARE = 'LINK_SHARE',        // Shared articles/websites
  POST_MESSAGE = 'POST_MESSAGE',    // Message posted to someone's timeline/wall
  REVIEW = 'REVIEW',                // Reviews of places, products, services
  ALBUM = 'ALBUM',                  // Photo collections

  // Rich Media Types (6)
  STORY = 'STORY',                  // Ephemeral 24h content
  REEL = 'REEL',                    // Short-form vertical video
  LIVE_STREAM = 'LIVE_STREAM',      // Live video broadcasts
  PLAYLIST = 'PLAYLIST',            // Music/video collections
  MOOD_BOARD = 'MOOD_BOARD',        // Visual inspiration
  TUTORIAL = 'TUTORIAL',            // How-to content

  // Interactive Content (6)
  POLL = 'POLL',                    // Voting posts
  QUIZ = 'QUIZ',                    // Trivia, personality tests
  SURVEY = 'SURVEY',                // Feedback collection
  Q_AND_A = 'Q_AND_A',              // Ask me anything
  CHALLENGE = 'CHALLENGE',          // Viral challenges/trends
  PETITION = 'PETITION',            // Social causes

  // Social & Location (8)
  CHECK_IN = 'CHECK_IN',            // Location-based posts
  TRAVEL_LOG = 'TRAVEL_LOG',        // Trip updates/itinerary
  LIFE_EVENT = 'LIFE_EVENT',        // Major life moments
  MILESTONE = 'MILESTONE',          // Personal achievements
  MEMORY = 'MEMORY',                // Throwback/flashback posts
  ANNIVERSARY = 'ANNIVERSARY',      // Yearly memories
  RECOMMENDATION = 'RECOMMENDATION', // Place/product suggestions
  GROUP_ACTIVITY = 'GROUP_ACTIVITY', // Group-specific content

  // Commercial & Business (6)
  PRODUCT_SHOWCASE = 'PRODUCT_SHOWCASE', // Selling items
  SERVICE_PROMOTION = 'SERVICE_PROMOTION', // Business services
  EVENT_PROMOTION = 'EVENT_PROMOTION',   // Events/meetups
  JOB_POSTING = 'JOB_POSTING',          // Hiring/career opportunities
  FUNDRAISER = 'FUNDRAISER',            // Charity/personal causes
  COLLABORATION = 'COLLABORATION',      // Creative projects

  // Specialized Content (8)
  RECIPE = 'RECIPE',                // Cooking/food content
  WORKOUT = 'WORKOUT',              // Fitness routines
  BOOK_REVIEW = 'BOOK_REVIEW',      // Reading updates
  MOOD_UPDATE = 'MOOD_UPDATE',      // Emotional status
  ACHIEVEMENT = 'ACHIEVEMENT',      // Gaming/app achievements
  QUOTE = 'QUOTE',                  // Inspirational quotes
  MUSIC = 'MUSIC',                  // Music sharing/streaming
  VENUE = 'VENUE'                   // Venue information/reviews
}

export enum PostContentType {
  TEXT = 'TEXT',
  IMAGE = 'IMAGE',
  VIDEO = 'VIDEO',
  MIXED = 'MIXED',
  LIVE = 'LIVE',
  POLL = 'POLL'
}

export interface PostContent {
  type: PostContentType;
  text?: string;
  images: PostImage[];
  videos: PostVideo[];
  hashtags: string[];
  mentions: string[];
  location?: string;
  poll?: PostPoll;
}

export interface PostImage {
  id: string;
  url: string;
  caption?: string;
  aspectRatio: number;
  filters: string[];
}

export interface PostVideo {
  id: string;
  url: string;
  thumbnailUrl?: string;
  duration: string;
  caption?: string;
  isAutoPlay: boolean;
}

export interface PostPoll {
  question: string;
  options: string[];
  votes: Record<number, number>; // option index -> vote count
  expiresAt?: string;
}

export enum ReactionType {
  LIKE = 'LIKE',
  LOVE = 'LOVE',
  LAUGH = 'LAUGH',
  WOW = 'WOW',
  SAD = 'SAD',
  ANGRY = 'ANGRY',
  CARE = 'CARE',
  CELEBRATE = 'CELEBRATE',
  SUPPORT = 'SUPPORT',
  CURIOUS = 'CURIOUS',
  DISAGREE = 'DISAGREE',
  SPAM = 'SPAM'
}

export interface PostReaction {
  type: ReactionType;
  userId: string;
  timestamp: string;
  userName?: string;
}

export enum ShareType {
  DIRECT_SHARE = 'DIRECT_SHARE',    // Simple reshare
  QUOTE_SHARE = 'QUOTE_SHARE',      // Share with comment
  STORY_SHARE = 'STORY_SHARE',      // Share to story
  MESSAGE_SHARE = 'MESSAGE_SHARE',  // Share via DM
  EXTERNAL_SHARE = 'EXTERNAL_SHARE' // Share outside platform
}

export interface PostShare {
  id: string;
  userId: string;
  userName: string;
  timestamp: string;
  shareType: ShareType;
  addedComment?: string;
  sharedToGroups: string[];
}

export interface PostComment {
  id: string;
  userId: string;
  userName: string;
  userAvatar?: string;
  content: string;
  timestamp: string;
  replies: PostComment[];
  reactions: PostReaction[];
  isEdited: boolean;
  isDeleted: boolean;
  mentionedUsers: string[];
}

export interface PostInteractions {
  reactions: PostReaction[];
  comments: PostComment[];
  shares: PostShare[];
  saves: string[];              // User IDs who saved
  views: number;
  reach: number;                // Unique users reached
  impressions: number;          // Total views
  clickThroughs: number;        // Link/action clicks
  engagementRate: number;       // Calculated engagement %

  // User-specific interaction states
  isLiked: boolean;
  isBookmarked: boolean;

  // Computed properties for legacy compatibility
  likes: number;
  commentsCount: number;
  sharesCount: number;
  savesCount: number;
}

export interface PostUser {
  id: string;
  username: string;
  displayName?: string;
  avatarUrl?: string;
  isVerified: boolean;
  followerCount: number;
  isFollowing: boolean;
}

export enum MoodType {
  HAPPY = 'HAPPY',
  EXCITED = 'EXCITED',
  GRATEFUL = 'GRATEFUL',
  LOVED = 'LOVED',
  BLESSED = 'BLESSED',
  RELAXED = 'RELAXED',
  CONTENT = 'CONTENT',
  MOTIVATED = 'MOTIVATED',
  PROUD = 'PROUD',
  ACCOMPLISHED = 'ACCOMPLISHED',
  TIRED = 'TIRED',
  STRESSED = 'STRESSED',
  SAD = 'SAD',
  ANXIOUS = 'ANXIOUS',
  CONFUSED = 'CONFUSED',
  FRUSTRATED = 'FRUSTRATED',
  ANGRY = 'ANGRY',
  LONELY = 'LONELY',
  NOSTALGIC = 'NOSTALGIC',
  CONTEMPLATIVE = 'CONTEMPLATIVE'
}

export enum FeelingType {
  AMAZING = 'AMAZING',
  FANTASTIC = 'FANTASTIC',
  GOOD = 'GOOD',
  OKAY = 'OKAY',
  MEH = 'MEH',
  NOT_GREAT = 'NOT_GREAT',
  TERRIBLE = 'TERRIBLE'
}

export enum ActivityType {
  EATING = 'EATING',
  DRINKING = 'DRINKING',
  TRAVELING = 'TRAVELING',
  EXERCISING = 'EXERCISING',
  WORKING = 'WORKING',
  STUDYING = 'STUDYING',
  READING = 'READING',
  WATCHING = 'WATCHING',
  LISTENING = 'LISTENING',
  PLAYING = 'PLAYING',
  COOKING = 'COOKING',
  SHOPPING = 'SHOPPING',
  CELEBRATING = 'CELEBRATING',
  RELAXING = 'RELAXING',
  SLEEPING = 'SLEEPING'
}

export enum ContentWarningType {
  NONE = 'NONE',
  SENSITIVE_CONTENT = 'SENSITIVE_CONTENT',
  GRAPHIC_VIOLENCE = 'GRAPHIC_VIOLENCE',
  ADULT_CONTENT = 'ADULT_CONTENT',
  DISTURBING_CONTENT = 'DISTURBING_CONTENT',
  SPOILER = 'SPOILER',
  FLASHING_LIGHTS = 'FLASHING_LIGHTS',
  POLITICAL_CONTENT = 'POLITICAL_CONTENT'
}

export enum LocationCategory {
  RESTAURANT = 'RESTAURANT',
  HOTEL = 'HOTEL',
  ATTRACTION = 'ATTRACTION',
  SHOPPING = 'SHOPPING',
  ENTERTAINMENT = 'ENTERTAINMENT',
  OUTDOORS = 'OUTDOORS',
  TRANSPORTATION = 'TRANSPORTATION',
  HOME = 'HOME',
  WORK = 'WORK',
  EDUCATION = 'EDUCATION',
  HEALTHCARE = 'HEALTHCARE',
  GOVERNMENT = 'GOVERNMENT',
  RELIGIOUS = 'RELIGIOUS',
  SPORTS = 'SPORTS',
  OTHER = 'OTHER'
}

export interface PostLocation {
  name: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  city?: string;
  country?: string;
  category?: LocationCategory;
}

export interface GeofenceRadius {
  latitude: number;
  longitude: number;
  distance: number; // Distance in kilometers
}

export interface LocationRestriction {
  countries: string[];
  regions: string[];
  cities: string[];
  radius?: GeofenceRadius;
}

export interface AgeRange {
  min?: number;
  max?: number;
}

export interface PostAudience {
  includedUsers: string[];     // Specific users who can see
  excludedUsers: string[];     // Users who cannot see
  includedGroups: string[];    // Specific groups/circles
  excludedGroups: string[];    // Excluded groups
  locationRestriction?: LocationRestriction;
  ageRange?: AgeRange;
  interests: string[];         // Interest-based targeting
}

export interface PostMetadata {
  targetType?: string;         // "product", "shop", "user", etc.
  targetId?: string;           // ID of the target
  targetName?: string;         // Name of the target
  rating?: number;             // For reviews
  price?: string;              // For product posts
  category?: string;           // Content category
  tags: string[];
  isSponsored: boolean;
  sponsorName?: string;
  editedAt?: string;
  originalPostId?: string;     // If this is a shared post
  isPromoted: boolean;         // Sponsored/boosted content
  mentionedUsers: string[];    // @mentions
  location?: PostLocation;
  mood?: MoodType;
  feeling?: FeelingType;
  activity?: ActivityType;
  contentWarning?: ContentWarningType;
  language?: string;
  isArchived: boolean;
  archivedAt?: string;
  expiresAt?: string;          // For stories/temporary content
  allowComments: boolean;
  allowShares: boolean;
  allowSaves: boolean;
  isPinned: boolean;
  isSticky: boolean;           // Stays at top of profile/group
}

/**
 * Enhanced Privacy Controls - 8 Privacy Levels
 * Aligned with mobile platform Post.kt
 */
export enum PostPrivacy {
  PUBLIC = 'PUBLIC',              // Visible to everyone
  FRIENDS = 'FRIENDS',            // Friends only
  CLOSE_FRIENDS = 'CLOSE_FRIENDS', // Close friends list
  FOLLOWERS = 'FOLLOWERS',        // Followers only
  MUTUAL_FRIENDS = 'MUTUAL_FRIENDS', // Mutual connections
  CUSTOM = 'CUSTOM',              // Custom audience
  UNLISTED = 'UNLISTED',          // Hidden from feeds but shareable
  PRIVATE = 'PRIVATE'             // Only author can see
}

// Legacy enum for backward compatibility
export enum PostVisibility {
  PUBLIC = 'PUBLIC',
  FRIENDS = 'FRIENDS',
  PRIVATE = 'PRIVATE',
  UNLISTED = 'UNLISTED'
}

/**
 * Enhanced Post Data Interface - Unified Social Platform Architecture
 * Aligned with mobile platform Post data class for cross-platform consistency
 */
export interface Post {
  id: string;
  type: PostType;
  user: PostUser;
  content: PostContent;
  interactions: PostInteractions;
  privacy: PostPrivacy;
  audience?: PostAudience;
  metadata?: PostMetadata;
  createdAt: string;
  updatedAt?: string;
  isEdited: boolean;

  // Legacy field for backward compatibility
  visibility: PostVisibility;
}

export interface PostHashtag {
  tag: string;
  count: number;
  isFollowing: boolean;
  category?: string;
}

/**
 * Post Type Guards and Validation Utilities
 * Similar to mobile platform's Kotlin type guards
 */
export class PostTypeValidator {

  static isImagePost(post: Post): boolean {
    return post.content.type === PostContentType.IMAGE && post.content.images.length > 0;
  }

  static isVideoPost(post: Post): boolean {
    return (post.content.type === PostContentType.VIDEO && post.content.videos.length > 0) ||
           ([PostType.VIDEO, PostType.REEL, PostType.LIVE_STREAM].includes(post.type));
  }

  static isPollPost(post: Post): boolean {
    return post.content.poll !== undefined || post.type === PostType.POLL;
  }

  static isStoryPost(post: Post): boolean {
    return post.type === PostType.STORY && post.metadata?.expiresAt !== undefined;
  }

  static isReviewPost(post: Post): boolean {
    return post.type === PostType.REVIEW && post.metadata?.rating !== undefined;
  }

  static isLocationPost(post: Post): boolean {
    return [PostType.CHECK_IN, PostType.TRAVEL_LOG].includes(post.type) ||
           post.metadata?.location !== undefined;
  }

  static isInteractivePost(post: Post): boolean {
    return [
      PostType.POLL, PostType.QUIZ, PostType.SURVEY,
      PostType.Q_AND_A, PostType.CHALLENGE, PostType.PETITION
    ].includes(post.type);
  }

  static isCommercialPost(post: Post): boolean {
    return [
      PostType.PRODUCT_SHOWCASE, PostType.SERVICE_PROMOTION,
      PostType.EVENT_PROMOTION, PostType.JOB_POSTING, PostType.FUNDRAISER
    ].includes(post.type) || post.metadata?.isPromoted === true;
  }

  static isEphemeralPost(post: Post): boolean {
    return post.type === PostType.STORY || post.metadata?.expiresAt !== undefined;
  }

  static getPostDominantContent(post: Post): PostContentType {
    if (post.content.images.length > 0) return PostContentType.IMAGE;
    if (post.content.videos.length > 0) return PostContentType.VIDEO;
    if (post.content.poll) return PostContentType.POLL;
    if (post.type === PostType.LIVE_STREAM) return PostContentType.LIVE;
    return PostContentType.TEXT;
  }
}

/**
 * Post Validation Rules
 */
export class PostValidationRules {

  static validatePost(post: Post): string[] {
    const errors: string[] = [];

    // Required field validation
    if (!post.id) {
      errors.push("Post ID is required");
    }

    if (!post.user.id) {
      errors.push("User ID is required");
    }

    // Content validation based on type
    switch (post.type) {
      case PostType.IMAGE:
      case PostType.ALBUM:
        if (post.content.images.length === 0) {
          errors.push("Image posts must contain at least one image");
        }
        break;
      case PostType.VIDEO:
      case PostType.REEL:
        if (post.content.videos.length === 0) {
          errors.push("Video posts must contain video content");
        }
        break;
      case PostType.POLL:
        if (!post.content.poll || post.content.poll.options.length < 2) {
          errors.push("Poll posts must have at least 2 options");
        }
        break;
      case PostType.TEXT:
        if (!post.content.text?.trim()) {
          errors.push("Text posts must contain text content");
        }
        break;
      case PostType.REVIEW:
        if (post.metadata?.rating === undefined) {
          errors.push("Review posts must include a rating");
        }
        break;
      case PostType.CHECK_IN:
        if (!post.metadata?.location) {
          errors.push("Check-in posts must include location");
        }
        break;
    }

    // Privacy validation
    if (post.privacy === PostPrivacy.CUSTOM && !post.audience) {
      errors.push("Custom privacy requires audience specification");
    }

    // Engagement validation
    if (post.interactions.engagementRate < 0 || post.interactions.engagementRate > 1) {
      errors.push("Engagement rate must be between 0 and 1");
    }

    return errors;
  }

  static isValidPostType(type: PostType, content: PostContent): boolean {
    switch (type) {
      case PostType.IMAGE:
      case PostType.ALBUM:
        return content.images.length > 0;
      case PostType.VIDEO:
      case PostType.REEL:
      case PostType.LIVE_STREAM:
        return content.videos.length > 0;
      case PostType.POLL:
        return content.poll !== undefined;
      case PostType.TEXT:
        return !!content.text?.trim();
      default:
        return true; // Other types are flexible
    }
  }
}

/**
 * Post Engagement Calculator
 */
export class PostEngagementCalculator {

  static calculateEngagementRate(post: Post): number {
    const totalEngagements = post.interactions.reactions.length +
                            post.interactions.comments.length +
                            post.interactions.shares.length;
    const views = post.interactions.views;

    return views > 0 ? Math.min(totalEngagements / views, 1) : 0;
  }

  static getTopReaction(post: Post): ReactionType | null {
    const reactionCounts = post.interactions.reactions.reduce((acc, reaction) => {
      acc[reaction.type] = (acc[reaction.type] || 0) + 1;
      return acc;
    }, {} as Record<ReactionType, number>);

    const topReaction = Object.entries(reactionCounts)
      .sort(([,a], [,b]) => b - a)[0];

    return topReaction ? topReaction[0] as ReactionType : null;
  }

  static getEngagementSummary(post: Post): Record<string, number> {
    return {
      reactions: post.interactions.reactions.length,
      comments: post.interactions.comments.length,
      shares: post.interactions.shares.length,
      saves: post.interactions.saves.length,
      views: post.interactions.views
    };
  }
}

// Legacy Post interface for backward compatibility
export interface LegacyPost {
  id: string;
  author: {
    name: string;
    avatar?: string;
    verified?: boolean;
    type: 'user' | 'merchant' | 'channel';
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
  type: 'text' | 'image' | 'live' | 'product';
  product?: {
    name: string;
    price: number;
    currency: string;
  };
  liveData?: {
    viewers: number;
    startTime: string;
    isLive: boolean;
  };
  source?: 'following' | 'trending' | 'interest' | 'sponsored';
}

/**
 * Utility functions to convert between legacy and new post formats
 */
export function convertLegacyToPost(legacyPost: LegacyPost): Post {
  const postType = (() => {
    switch (legacyPost.type) {
      case 'text': return PostType.TEXT;
      case 'image': return PostType.IMAGE;
      case 'live': return PostType.LIVE_STREAM;
      case 'product': return PostType.PRODUCT_SHOWCASE;
      default: return PostType.TEXT;
    }
  })();

  return {
    id: legacyPost.id,
    type: postType,
    user: {
      id: legacyPost.id + '_user',
      username: legacyPost.author.name.toLowerCase().replace(/\s+/g, '_'),
      displayName: legacyPost.author.name,
      avatarUrl: legacyPost.author.avatar,
      isVerified: legacyPost.author.verified || false,
      followerCount: 0,
      isFollowing: false
    },
    content: {
      type: legacyPost.images?.length ? PostContentType.IMAGE : PostContentType.TEXT,
      text: legacyPost.content,
      images: legacyPost.images?.map((url, index) => ({
        id: `${legacyPost.id}_img_${index}`,
        url,
        caption: undefined,
        aspectRatio: 1,
        filters: []
      })) || [],
      videos: [],
      hashtags: legacyPost.tags || [],
      mentions: [],
      location: legacyPost.location
    },
    interactions: {
      reactions: [],
      comments: [],
      shares: [],
      saves: [],
      views: legacyPost.likes * 10, // Estimate views from likes
      reach: 0,
      impressions: 0,
      clickThroughs: 0,
      engagementRate: 0,
      isLiked: legacyPost.isLiked || false,
      isBookmarked: false,
      likes: legacyPost.likes,
      commentsCount: legacyPost.comments,
      sharesCount: legacyPost.shares,
      savesCount: 0
    },
    privacy: PostPrivacy.PUBLIC,
    metadata: {
      tags: legacyPost.tags || [],
      isSponsored: legacyPost.source === 'sponsored',
      isPromoted: false,
      mentionedUsers: [],
      isArchived: false,
      allowComments: true,
      allowShares: true,
      allowSaves: true,
      isPinned: false,
      isSticky: false
    },
    createdAt: legacyPost.timestamp,
    isEdited: false,
    visibility: PostVisibility.PUBLIC
  };
}

export function convertPostToLegacy(post: Post): LegacyPost {
  const legacyType = (() => {
    switch (post.type) {
      case PostType.IMAGE:
      case PostType.ALBUM:
        return 'image';
      case PostType.LIVE_STREAM:
        return 'live';
      case PostType.PRODUCT_SHOWCASE:
        return 'product';
      default:
        return 'text';
    }
  })();

  return {
    id: post.id,
    author: {
      name: post.user.displayName || post.user.username,
      avatar: post.user.avatarUrl,
      verified: post.user.isVerified,
      type: 'user'
    },
    content: post.content.text || '',
    images: post.content.images.map(img => img.url),
    timestamp: post.createdAt,
    likes: post.interactions.likes,
    comments: post.interactions.commentsCount,
    shares: post.interactions.sharesCount,
    isLiked: post.interactions.isLiked,
    location: post.content.location || post.metadata?.location?.name,
    tags: post.content.hashtags,
    type: legacyType,
    source: post.metadata?.isSponsored ? 'sponsored' : 'following'
  };
}