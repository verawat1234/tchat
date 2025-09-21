/**
 * Comprehensive Schema for Telegram SEA Edition
 * Generated from UI component analysis
 * Covers all domains: User, Chat, Social, Commerce, Events, Wallet, Workspace, Video, Discovery, Activities
 */

// =============================================================================
// UTILITY TYPES
// =============================================================================

export type UUID = string;
export type Timestamp = string;
export type Currency = 'THB' | 'SGD' | 'IDR' | 'MYR' | 'PHP' | 'VND' | 'USD';
export type Locale = 'th-TH' | 'id-ID' | 'ms-MY' | 'vi-VN' | 'en-US';
export type CountryCode = 'TH' | 'ID' | 'MY' | 'VN' | 'SG' | 'PH';

// =============================================================================
// USER & AUTHENTICATION DOMAIN
// =============================================================================

export interface User {
  id: UUID;
  phone?: string;
  email?: string;
  name: string;
  avatar?: string;
  country: CountryCode;
  locale: Locale;
  kycTier: 1 | 2 | 3;
  status: UserStatus;
  lastSeen?: Timestamp;
  isVerified: boolean;
  settings: UserSettings;
  profile: UserProfile;
  preferences: UserPreferences;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type UserStatus = 'online' | 'offline' | 'away' | 'busy';

export interface UserProfile {
  displayName?: string;
  bio?: string;
  birthday?: string;
  gender?: 'male' | 'female' | 'other' | 'prefer_not_to_say';
  location?: string;
  website?: string;
  occupation?: string;
  interests: string[];
  languages: Locale[];
  timezone: string;
}

export interface UserSettings {
  privacy: PrivacySettings;
  notifications: NotificationSettings;
  appearance: AppearanceSettings;
  security: SecuritySettings;
}

export interface PrivacySettings {
  profileVisibility: 'public' | 'friends' | 'private';
  phoneVisibility: 'public' | 'friends' | 'private';
  lastSeenVisibility: 'everyone' | 'friends' | 'nobody';
  readReceiptsEnabled: boolean;
  onlineStatusVisible: boolean;
}

export interface NotificationSettings {
  pushEnabled: boolean;
  emailEnabled: boolean;
  messageNotifications: boolean;
  postNotifications: boolean;
  eventNotifications: boolean;
  paymentNotifications: boolean;
  mutedChats: UUID[];
  mutedUsers: UUID[];
}

export interface AppearanceSettings {
  theme: 'light' | 'dark' | 'auto';
  language: Locale;
  fontSize: 'small' | 'medium' | 'large';
  chatWallpaper?: string;
}

export interface SecuritySettings {
  twoFactorEnabled: boolean;
  biometricEnabled: boolean;
  autoLockTimeout: number; // minutes
  trustedDevices: TrustedDevice[];
}

export interface TrustedDevice {
  id: UUID;
  deviceName: string;
  deviceType: 'mobile' | 'desktop' | 'tablet';
  lastUsed: Timestamp;
  ipAddress?: string;
  location?: string;
}

export interface UserPreferences {
  defaultCurrency: Currency;
  defaultPaymentMethod?: UUID;
  eventCategories: string[];
  productCategories: string[];
  contentLanguages: Locale[];
  contentFilters: string[];
}

export interface KYCInfo {
  tier: 1 | 2 | 3;
  status: 'pending' | 'approved' | 'rejected' | 'incomplete';
  documents: KYCDocument[];
  verifiedAt?: Timestamp;
  expiresAt?: Timestamp;
  dailyLimit: number;
  monthlyLimit: number;
  usedThisMonth: number;
}

export interface KYCDocument {
  id: UUID;
  type: 'id_card' | 'passport' | 'driving_license' | 'utility_bill' | 'bank_statement';
  fileUrl: string;
  status: 'pending' | 'approved' | 'rejected';
  uploadedAt: Timestamp;
  verifiedAt?: Timestamp;
  rejectionReason?: string;
  metadata?: Record<string, any>;
}

export interface Friend {
  id: UUID;
  userId: UUID;
  friendId: UUID;
  status: FriendStatus;
  mutualFriends: number;
  commonInterests: string[];
  createdAt: Timestamp;
  acceptedAt?: Timestamp;
}

export type FriendStatus = 'pending' | 'accepted' | 'blocked';

export interface AuthSession {
  id: UUID;
  userId: UUID;
  deviceId: string;
  deviceInfo: DeviceInfo;
  accessTokenHash: string;
  refreshTokenHash: string;
  expiresAt: Timestamp;
  isActive: boolean;
  ipAddress?: string;
  location?: SessionLocation;
  userAgent?: string;
  createdAt: Timestamp;
  lastUsed: Timestamp;
}

export interface DeviceInfo {
  type: 'mobile' | 'desktop' | 'tablet' | 'web';
  os: string;
  browser?: string;
  appVersion?: string;
  pushToken?: string;
  notificationEnabled: boolean;
}

export interface SessionLocation {
  country?: CountryCode;
  city?: string;
  region?: string;
  latitude?: number;
  longitude?: number;
  timezone?: string;
  ipLookupProvider?: string;
}

export interface OTPVerification {
  id: UUID;
  userId?: UUID;
  phone?: string;
  email?: string;
  code: string;
  type: 'registration' | 'login' | 'password_reset' | 'phone_verification';
  attempts: number;
  maxAttempts: number;
  isUsed: boolean;
  expiresAt: Timestamp;
  createdAt: Timestamp;
  usedAt?: Timestamp;
}

export interface UserActivity {
  id: UUID;
  userId: UUID;
  action: UserActionType;
  targetType?: string;
  targetId?: UUID;
  metadata?: Record<string, any>;
  ipAddress?: string;
  userAgent?: string;
  location?: string;
  createdAt: Timestamp;
}

export type UserActionType =
  | 'login'
  | 'logout'
  | 'profile_update'
  | 'settings_change'
  | 'kyc_submission'
  | 'password_change'
  | 'device_added'
  | 'friend_request'
  | 'friend_accept'
  | 'privacy_change';

export interface UserStats {
  userId: UUID;
  totalMessages: number;
  totalPosts: number;
  totalFriends: number;
  totalOrders: number;
  totalSpent: number;
  eventsAttended: number;
  videosWatched: number;
  achievementsEarned: number;
  joinedAt: Timestamp;
  lastActive: Timestamp;
  updatedAt: Timestamp;
}

export const KYC_LIMITS = {
  1: { daily: 5000, monthly: 20000 },
  2: { daily: 50000, monthly: 200000 },
  3: { daily: 500000, monthly: 2000000 }
} as const;

export const SEA_COUNTRIES: CountryCode[] = ['TH', 'ID', 'MY', 'VN', 'SG', 'PH'];

export const LOCALE_COUNTRY_MAP: Record<Locale, CountryCode> = {
  'th-TH': 'TH',
  'id-ID': 'ID',
  'ms-MY': 'MY',
  'vi-VN': 'VN',
  'en-US': 'SG'
};

export const DEFAULT_PREFERENCES: Record<CountryCode, Partial<UserPreferences>> = {
  TH: {
    defaultCurrency: 'THB',
    contentLanguages: ['th-TH', 'en-US'],
    eventCategories: ['music', 'food', 'cultural', 'temple']
  },
  ID: {
    defaultCurrency: 'IDR',
    contentLanguages: ['id-ID', 'en-US'],
    eventCategories: ['music', 'food', 'cultural']
  },
  MY: {
    defaultCurrency: 'MYR',
    contentLanguages: ['ms-MY', 'en-US'],
    eventCategories: ['music', 'food', 'cultural']
  },
  VN: {
    defaultCurrency: 'VND',
    contentLanguages: ['vi-VN', 'en-US'],
    eventCategories: ['music', 'food', 'cultural']
  },
  SG: {
    defaultCurrency: 'SGD',
    contentLanguages: ['en-US'],
    eventCategories: ['music', 'food', 'business']
  },
  PH: {
    defaultCurrency: 'PHP',
    contentLanguages: ['en-US'],
    eventCategories: ['music', 'food', 'cultural']
  }
};

// =============================================================================
// CHAT & MESSAGING DOMAIN
// =============================================================================

// =============================================================================
// DIALOG MANAGEMENT
// =============================================================================

export interface Dialog {
  id: UUID;
  type: DialogType;
  name?: string;
  description?: string;
  avatar?: string;
  participants: UUID[];
  admins: UUID[];
  owners: UUID[];
  lastMessageId?: UUID;
  lastMessage?: Message;
  unreadCount: number;
  mutedUntil?: Timestamp;
  isPinned: boolean;
  isArchived: boolean;
  permissions: DialogPermissions;
  settings: DialogSettings;
  metadata?: DialogMetadata;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type DialogType = 'user' | 'group' | 'channel' | 'bot' | 'business';

export interface DialogPermissions {
  canSendMessages: boolean;
  canSendMedia: boolean;
  canSendVoice: boolean;
  canSendFiles: boolean;
  canAddMembers: boolean;
  canChangeInfo: boolean;
  canPinMessages: boolean;
  canDeleteMessages: boolean;
}

export interface DialogSettings {
  allowedMessageTypes: MessageType[];
  autoDeleteTimeout?: number; // seconds
  slowModeDelay?: number; // seconds
  maxMembers?: number;
  requireApproval: boolean;
  allowInviteLinks: boolean;
  businessMode?: BusinessChatSettings;
}

export interface BusinessChatSettings {
  isBusinessChat: boolean;
  leadScore?: number;
  customerTier?: 'bronze' | 'silver' | 'gold' | 'platinum';
  tags: string[];
  assignedAgent?: UUID;
  automatedResponses: boolean;
  businessHours?: BusinessHours[];
}

export interface BusinessHours {
  dayOfWeek: number; // 0-6, Sunday = 0
  openTime: string; // HH:mm
  closeTime: string; // HH:mm
  isClosed: boolean;
  timezone: string;
}

export interface DialogMetadata {
  customerInfo?: CustomerInfo;
  productContext?: ProductContext[];
  orderHistory?: UUID[];
  supportTickets?: UUID[];
  notes?: string;
}

export interface CustomerInfo {
  preferredLanguage?: Locale;
  location?: string;
  purchaseHistory: number;
  totalSpent: number;
  currency: Currency;
  loyaltyLevel?: string;
  specialRequests?: string[];
}

export interface ProductContext {
  productId: UUID;
  inquiryType: 'question' | 'complaint' | 'support' | 'purchase_intent';
  status: 'open' | 'resolved' | 'escalated';
  priority: 'low' | 'medium' | 'high' | 'urgent';
}

// =============================================================================
// MESSAGE SYSTEM
// =============================================================================

export interface Message {
  id: UUID;
  dialogId: UUID;
  senderId: UUID;
  type: MessageType;
  content: MessageContent;
  replyToId?: UUID;
  forwardFromId?: UUID;
  threadId?: UUID;
  isEdited: boolean;
  isPinned: boolean;
  reactions: MessageReaction[];
  readBy: MessageRead[];
  mentions: UUID[];
  attachments: MessageAttachment[];
  metadata?: MessageMetadata;
  createdAt: Timestamp;
  editedAt?: Timestamp;
  deletedAt?: Timestamp;
}

export type MessageType =
  | 'text'
  | 'voice'
  | 'file'
  | 'image'
  | 'video'
  | 'payment'
  | 'system'
  | 'location'
  | 'contact'
  | 'poll'
  | 'event'
  | 'product'
  | 'sticker'
  | 'gif';

export interface MessageContent {
  text?: string;
  html?: string;
  voice?: VoiceMessage;
  file?: FileMessage;
  image?: ImageMessage;
  video?: VideoMessage;
  payment?: PaymentMessage;
  system?: SystemMessage;
  location?: LocationMessage;
  contact?: ContactMessage;
  poll?: PollMessage;
  event?: EventMessage;
  product?: ProductMessage;
  sticker?: StickerMessage;
  gif?: GifMessage;
}

// =============================================================================
// MESSAGE CONTENT TYPES
// =============================================================================

export interface VoiceMessage {
  fileUrl: string;
  duration: number; // seconds
  waveform: number[];
  transcript?: string;
  transcriptLanguage?: Locale;
  fileSize: number;
  encoding: string;
  sampleRate: number;
}

export interface FileMessage {
  fileUrl: string;
  fileName: string;
  fileSize: number;
  mimeType: string;
  thumbnail?: string;
  virus_scan_status?: 'pending' | 'clean' | 'infected' | 'error';
  downloadCount?: number;
}

export interface ImageMessage {
  fileUrl: string;
  thumbnailUrl?: string;
  width: number;
  height: number;
  fileSize: number;
  caption?: string;
  altText?: string;
  blurHash?: string;
  exifData?: Record<string, any>;
}

export interface VideoMessage {
  fileUrl: string;
  thumbnailUrl?: string;
  duration: number; // seconds
  width: number;
  height: number;
  fileSize: number;
  caption?: string;
  playbackCount?: number;
  quality?: '720p' | '1080p' | '480p' | '360p';
}

export interface PaymentMessage {
  amount: number;
  currency: Currency;
  description: string;
  transactionId: UUID;
  status: 'pending' | 'completed' | 'failed' | 'cancelled';
  paymentMethod?: string;
  recipient?: string;
  reference?: string;
  qrCode?: string;
  deepLink?: string;
}

export interface SystemMessage {
  type: SystemMessageType;
  data: Record<string, any>;
  isAutomated: boolean;
  templateId?: string;
}

export type SystemMessageType =
  | 'user_joined'
  | 'user_left'
  | 'title_changed'
  | 'avatar_changed'
  | 'settings_changed'
  | 'member_promoted'
  | 'member_restricted'
  | 'message_pinned'
  | 'business_hours_changed'
  | 'auto_response';

export interface LocationMessage {
  latitude: number;
  longitude: number;
  address?: string;
  placeName?: string;
  venueId?: UUID;
  livePeriod?: number; // seconds
  accuracy?: number; // meters
  heading?: number; // degrees
}

export interface ContactMessage {
  phoneNumber: string;
  firstName: string;
  lastName?: string;
  userId?: UUID;
  vCard?: string;
  avatar?: string;
}

export interface PollMessage {
  question: string;
  options: PollOption[];
  allowMultiple: boolean;
  isAnonymous: boolean;
  correctOptionId?: number;
  explanation?: string;
  closePeriod?: number; // seconds
  closeDate?: Timestamp;
  totalVoters: number;
}

export interface PollOption {
  id: number;
  text: string;
  voterCount: number;
  voters?: UUID[];
  percentage: number;
}

export interface EventMessage {
  eventId: UUID;
  title: string;
  description?: string;
  imageUrl?: string;
  startDate: Timestamp;
  endDate?: Timestamp;
  location?: string;
  ticketsAvailable: boolean;
  priceRange?: { min: number; max: number; currency: Currency };
}

export interface ProductMessage {
  productId: UUID;
  title: string;
  description?: string;
  imageUrl?: string;
  price: number;
  compareAtPrice?: number;
  currency: Currency;
  shopId: UUID;
  shopName: string;
  isInStock: boolean;
  rating?: number;
  reviewCount?: number;
}

export interface StickerMessage {
  stickerId: UUID;
  stickerUrl: string;
  packId?: UUID;
  packName?: string;
  emoji?: string;
  width: number;
  height: number;
  isAnimated: boolean;
}

export interface GifMessage {
  gifUrl: string;
  thumbnailUrl?: string;
  width: number;
  height: number;
  duration?: number;
  fileSize: number;
  source?: string;
  searchTerm?: string;
}

// =============================================================================
// MESSAGE INTERACTIONS
// =============================================================================

export interface MessageReaction {
  id: UUID;
  messageId: UUID;
  userId: UUID;
  emoji: string;
  skinTone?: string;
  createdAt: Timestamp;
}

export interface MessageRead {
  messageId: UUID;
  userId: UUID;
  readAt: Timestamp;
  deviceId?: string;
}

export interface MessageAttachment {
  id: UUID;
  messageId: UUID;
  type: AttachmentType;
  fileUrl: string;
  fileName: string;
  fileSize: number;
  mimeType: string;
  uploadProgress?: number;
  downloadProgress?: number;
  metadata?: AttachmentMetadata;
  createdAt: Timestamp;
}

export type AttachmentType = 'file' | 'image' | 'video' | 'audio' | 'document';

export interface AttachmentMetadata {
  width?: number;
  height?: number;
  duration?: number;
  thumbnail?: string;
  virus_scan?: 'pending' | 'clean' | 'infected';
  compression?: 'none' | 'low' | 'medium' | 'high';
  originalSize?: number;
}

export interface MessageMetadata {
  editHistory?: MessageEdit[];
  deliveryStatus?: MessageDeliveryStatus;
  businessContext?: BusinessMessageContext;
  analytics?: MessageAnalytics;
  autoGenerated?: boolean;
  templateId?: string;
  campaignId?: UUID;
}

export interface MessageEdit {
  editedAt: Timestamp;
  previousContent: string;
  reason?: string;
}

export interface MessageDeliveryStatus {
  sent: boolean;
  sentAt?: Timestamp;
  delivered: boolean;
  deliveredAt?: Timestamp;
  failed: boolean;
  failureReason?: string;
  retryCount: number;
}

export interface BusinessMessageContext {
  isCustomerSupport: boolean;
  ticketId?: UUID;
  priority?: 'low' | 'medium' | 'high' | 'urgent';
  category?: string;
  tags?: string[];
  assignedAgent?: UUID;
  escalated?: boolean;
  satisfactionRating?: number;
}

export interface MessageAnalytics {
  opens: number;
  clicks: number;
  shares: number;
  reactions: number;
  forwardCount: number;
  replyCount: number;
  firstOpenAt?: Timestamp;
  lastOpenAt?: Timestamp;
}

// =============================================================================
// RICH TEXT & FORMATTING
// =============================================================================

export interface RichTextMessage {
  content: string;
  entities: MessageEntity[];
  formatting: MessageFormatting;
}

export interface MessageEntity {
  type: EntityType;
  offset: number;
  length: number;
  url?: string;
  userId?: UUID;
  language?: string;
  customData?: Record<string, any>;
}

export type EntityType =
  | 'mention'
  | 'hashtag'
  | 'url'
  | 'email'
  | 'phone'
  | 'bold'
  | 'italic'
  | 'underline'
  | 'strikethrough'
  | 'code'
  | 'pre'
  | 'spoiler'
  | 'custom_emoji';

export interface MessageFormatting {
  allowMarkdown: boolean;
  allowHtml: boolean;
  maxLength: number;
  allowedEntities: EntityType[];
  customStyles?: Record<string, string>;
}

// =============================================================================
// CHAT INPUT & COMPOSITION
// =============================================================================

export interface ChatInput {
  dialogId: UUID;
  content: string;
  draftId?: UUID;
  replyToId?: UUID;
  attachments: PendingAttachment[];
  mentions: UUID[];
  scheduling?: MessageScheduling;
  formatting: MessageFormatting;
  voiceRecording?: VoiceRecording;
}

export interface PendingAttachment {
  id: UUID;
  file: File | Blob;
  type: AttachmentType;
  uploadProgress: number;
  thumbnail?: string;
  metadata?: Partial<AttachmentMetadata>;
}

export interface MessageScheduling {
  scheduledFor: Timestamp;
  timezone: string;
  recurring?: RecurringSchedule;
}

export interface RecurringSchedule {
  pattern: 'daily' | 'weekly' | 'monthly';
  interval: number;
  endDate?: Timestamp;
  maxOccurrences?: number;
}

export interface VoiceRecording {
  isRecording: boolean;
  duration: number;
  audioBlob?: Blob;
  waveform: number[];
  pausedAt?: number;
  segments: VoiceSegment[];
}

export interface VoiceSegment {
  startTime: number;
  endTime: number;
  volume: number;
  pauseDuration?: number;
}

// =============================================================================
// DRAFT MANAGEMENT
// =============================================================================

export interface MessageDraft {
  id: UUID;
  dialogId: UUID;
  userId: UUID;
  content: string;
  attachments: UUID[];
  replyToId?: UUID;
  mentions: UUID[];
  lastModified: Timestamp;
  autoSaveEnabled: boolean;
  expiresAt?: Timestamp;
}

// =============================================================================
// SEARCH & FILTERING
// =============================================================================

export interface MessageSearchQuery {
  query: string;
  dialogId?: UUID;
  senderId?: UUID;
  messageType?: MessageType;
  dateRange?: { start: Timestamp; end: Timestamp };
  hasAttachments?: boolean;
  isForwarded?: boolean;
  isEdited?: boolean;
  limit: number;
  offset: number;
}

export interface MessageSearchResult {
  messageId: UUID;
  dialogId: UUID;
  snippet: string;
  highlightedContent: string;
  matchScore: number;
  context: { before: Message[]; after: Message[] };
  matchedEntities: MessageEntity[];
}

// =============================================================================
// BUSINESS LOGIC CONSTANTS
// =============================================================================

/**
 * Message limits by type and user tier
 */
export const MESSAGE_LIMITS = {
  text: { free: 4096, premium: 8192 },
  voice: { free: 300, premium: 600 }, // seconds
  file: { free: 50, premium: 500 }, // MB
  video: { free: 100, premium: 1000 } // MB
} as const;

/**
 * Auto-delete timeouts
 */
export const AUTO_DELETE_OPTIONS = [
  { label: '1 minute', value: 60 },
  { label: '1 hour', value: 3600 },
  { label: '1 day', value: 86400 },
  { label: '1 week', value: 604800 },
  { label: '1 month', value: 2592000 }
] as const;

/**
 * Supported file types for attachments
 */
export const SUPPORTED_FILE_TYPES = {
  image: ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'],
  video: ['mp4', 'mov', 'avi', 'mkv', 'webm'],
  audio: ['mp3', 'wav', 'ogg', 'm4a', 'aac'],
  document: ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt']
} as const;

/**
 * Maximum participants by dialog type
 */
export const MAX_PARTICIPANTS = {
  user: 2,
  group: 200,
  channel: 100000,
  business: 10
} as const;
// =============================================================================
// SOCIAL DOMAIN
// =============================================================================

export interface Post {
  id: UUID;
  authorId: UUID;
  type: 'text' | 'image' | 'video' | 'link' | 'poll' | 'moment';
  content: PostContent;
  privacy: 'public' | 'friends' | 'private';
  location?: Location;
  tags: string[];
  mentions: UUID[];
  engagement: PostEngagement;
  isPromoted: boolean;
  scheduledAt?: Timestamp;
  createdAt: Timestamp;
  updatedAt: Timestamp;
  deletedAt?: Timestamp;
}

export interface PostContent {
  text?: string;
  images?: PostImage[];
  videos?: PostVideo[];
  link?: PostLink;
  poll?: PostPoll;
  moment?: PostMoment;
}

export interface PostImage {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  width: number;
  height: number;
  caption?: string;
  order: number;
}

export interface PostVideo {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  duration: number;
  width: number;
  height: number;
  caption?: string;
  order: number;
}

export interface PostLink {
  url: string;
  title?: string;
  description?: string;
  imageUrl?: string;
  siteName?: string;
}

export interface PostPoll {
  question: string;
  options: string[];
  allowMultiple: boolean;
  expiresAt?: Timestamp;
  votes: PollVote[];
}

export interface PollVote {
  userId: UUID;
  optionIndex: number;
  createdAt: Timestamp;
}

export interface PostMoment {
  id: UUID;
  type: 'story' | 'highlight';
  mediaUrl: string;
  thumbnailUrl?: string;
  duration?: number;
  expiresAt?: Timestamp;
  viewCount: number;
  viewers: MomentView[];
}

export interface MomentView {
  userId: UUID;
  viewedAt: Timestamp;
  watchDuration?: number;
}

export interface PostEngagement {
  likeCount: number;
  commentCount: number;
  shareCount: number;
  viewCount: number;
  likes: PostLike[];
  comments: Comment[];
  shares: PostShare[];
}

export interface PostLike {
  id: UUID;
  postId: UUID;
  userId: UUID;
  createdAt: Timestamp;
}

export interface Comment {
  id: UUID;
  postId: UUID;
  authorId: UUID;
  parentId?: UUID; // for threaded comments
  content: string;
  images?: string[];
  likeCount: number;
  replyCount: number;
  likes: CommentLike[];
  replies: Comment[];
  createdAt: Timestamp;
  editedAt?: Timestamp;
  deletedAt?: Timestamp;
}

export interface CommentLike {
  id: UUID;
  commentId: UUID;
  userId: UUID;
  createdAt: Timestamp;
}

export interface PostShare {
  id: UUID;
  postId: UUID;
  userId: UUID;
  type: 'story' | 'message' | 'external';
  targetId?: UUID; // dialog id for message shares
  createdAt: Timestamp;
}

export interface FriendActivity {
  id: UUID;
  userId: UUID;
  friendId: UUID;
  type: 'post_liked' | 'post_commented' | 'post_shared' | 'friend_added' | 'event_joined';
  targetId: UUID;
  createdAt: Timestamp;
}

// =============================================================================
// E-COMMERCE DOMAIN
// =============================================================================

// =============================================================================
// PRODUCT MANAGEMENT
// =============================================================================

export interface Product {
  id: UUID;
  shopId: UUID;
  title: string;
  description: string;
  shortDescription?: string;
  images: ProductImage[];
  videos?: ProductVideo[];
  price: number;
  compareAtPrice?: number;
  currency: Currency;
  cost?: number;
  sku?: string;
  barcode?: string;
  inventory: ProductInventory;
  variants: ProductVariant[];
  category: ProductCategory;
  tags: string[];
  attributes: ProductAttribute[];
  seo: ProductSEO;
  status: ProductStatus;
  isDigital: boolean;
  weight?: number;
  dimensions?: ProductDimensions;
  shipping: ProductShipping;
  ratings: ProductRating;
  reviews: ProductReview[];
  localization: ProductLocalization;
  compliance: ProductCompliance;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type ProductStatus = 'draft' | 'active' | 'out_of_stock' | 'discontinued' | 'archived';

export interface ProductImage {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  alt?: string;
  order: number;
  isPrimary: boolean;
  variants?: UUID[]; // Which variants this image applies to
}

export interface ProductVideo {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  duration: number;
  title?: string;
  order: number;
  type: 'demo' | 'review' | 'unboxing' | 'tutorial';
}

export interface ProductInventory {
  trackQuantity: boolean;
  quantity?: number;
  lowStockThreshold?: number;
  isInStock: boolean;
  allowBackorders: boolean;
  location?: string;
  reservedQuantity?: number;
  damagedQuantity?: number;
}

export interface ProductVariant {
  id: UUID;
  title: string;
  sku?: string;
  price?: number;
  compareAtPrice?: number;
  cost?: number;
  inventory: ProductInventory;
  options: VariantOption[];
  image?: string;
  weight?: number;
  dimensions?: ProductDimensions;
  isDefault: boolean;
}

export interface VariantOption {
  name: string; // e.g., "Color", "Size", "Material"
  value: string; // e.g., "Red", "Large", "Cotton"
  displayName?: string; // Localized display name
  colorCode?: string; // For color variants
  imageUrl?: string; // For image-based variants
}

export interface ProductDimensions {
  length: number;
  width: number;
  height: number;
  unit: 'cm' | 'in' | 'mm';
}

export interface ProductShipping {
  isShippingRequired: boolean;
  weight?: number;
  dimensions?: ProductDimensions;
  shippingClass?: string;
  processingTime?: number; // days
  shippingMethods: ShippingMethod[];
  restrictions?: ShippingRestriction[];
}

export interface ShippingMethod {
  id: UUID;
  name: string;
  description?: string;
  price: number;
  currency: Currency;
  estimatedDays: { min: number; max: number };
  regions: CountryCode[];
  carrier?: string;
  trackingEnabled: boolean;
}

export interface ShippingRestriction {
  countries: CountryCode[];
  reason: string;
  type: 'prohibited' | 'restricted' | 'requires_permit';
}

export interface ProductCategory {
  id: UUID;
  name: string;
  slug: string;
  parentId?: UUID;
  description?: string;
  image?: string;
  icon?: string;
  isActive: boolean;
  sortOrder: number;
  seoTitle?: string;
  seoDescription?: string;
  localization: CategoryLocalization;
}

export interface CategoryLocalization {
  [locale: string]: {
    name: string;
    description?: string;
    seoTitle?: string;
    seoDescription?: string;
  };
}

export interface ProductAttribute {
  name: string;
  value: string;
  displayName?: string;
  isVariant: boolean;
  isFilterable: boolean;
  isRequired: boolean;
  type: 'text' | 'number' | 'boolean' | 'date' | 'color' | 'image';
  unit?: string;
}

export interface ProductSEO {
  title?: string;
  description?: string;
  keywords: string[];
  slug: string;
  metaTags?: Record<string, string>;
  canonicalUrl?: string;
  openGraphImage?: string;
}

export interface ProductRating {
  averageRating: number;
  totalReviews: number;
  distribution: { [key: number]: number }; // rating -> count
  verifiedPurchaseRating?: number;
  recentRating?: number; // Last 30 days
}

export interface ProductReview {
  id: UUID;
  productId: UUID;
  userId: UUID;
  variantId?: UUID;
  orderId?: UUID;
  rating: number;
  title?: string;
  comment: string;
  pros?: string[];
  cons?: string[];
  images?: ReviewImage[];
  videos?: ReviewVideo[];
  isVerifiedPurchase: boolean;
  helpfulCount: number;
  reportCount: number;
  moderationStatus: 'pending' | 'approved' | 'rejected' | 'hidden';
  language: Locale;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface ReviewImage {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  caption?: string;
  order: number;
}

export interface ReviewVideo {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  duration: number;
  caption?: string;
  order: number;
}

export interface ProductLocalization {
  [locale: string]: {
    title: string;
    description: string;
    shortDescription?: string;
    tags: string[];
    attributes: { [key: string]: string };
  };
}

export interface ProductCompliance {
  certifications: ProductCertification[];
  warnings: string[];
  ageRestriction?: number;
  countryRestrictions: CountryCode[];
  requires_id_verification?: boolean;
}

export interface ProductCertification {
  type: string; // e.g., "CE", "FDA", "HALAL", "Organic"
  number?: string;
  issuer: string;
  issuedAt: Timestamp;
  expiresAt?: Timestamp;
  documentUrl?: string;
}

// =============================================================================
// SHOP MANAGEMENT
// =============================================================================

export interface Shop {
  id: UUID;
  ownerId: UUID;
  name: string;
  description?: string;
  avatar?: string;
  coverImage?: string;
  isVerified: boolean;
  verificationLevel: 'none' | 'basic' | 'premium' | 'enterprise';
  status: ShopStatus;
  settings: ShopSettings;
  contact: ShopContact;
  location?: ShopLocation;
  stats: ShopStats;
  policies: ShopPolicies;
  categories: string[];
  tags: string[];
  subscription: ShopSubscription;
  compliance: ShopCompliance;
  localization: ShopLocalization;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type ShopStatus = 'active' | 'suspended' | 'under_review' | 'closed' | 'maintenance';

export interface ShopSettings {
  isPublic: boolean;
  allowReviews: boolean;
  autoApproveOrders: boolean;
  currency: Currency;
  timezone: string;
  businessHours: ShopBusinessHours[];
  minimumOrder?: number;
  freeShippingThreshold?: number;
  taxSettings: TaxSettings;
  returnSettings: ReturnSettings;
}

export interface ShopBusinessHours {
  dayOfWeek: number; // 0-6, Sunday = 0
  openTime: string; // HH:mm
  closeTime: string; // HH:mm
  isClosed: boolean;
  isDeliveryAvailable: boolean;
  specialNotes?: string;
}

export interface TaxSettings {
  includeTax: boolean;
  taxRate: number;
  taxLabel: string;
  exemptProducts: UUID[];
  countrySpecific: { [country: string]: number };
}

export interface ReturnSettings {
  allowReturns: boolean;
  returnPeriod: number; // days
  returnShipping: 'customer_pays' | 'shop_pays' | 'free';
  conditions: string[];
  restockingFee?: number;
}

export interface ShopContact {
  email?: string;
  phone?: string;
  whatsapp?: string;
  website?: string;
  socialMedia: { [platform: string]: string };
  supportHours?: ShopBusinessHours[];
  responseTime?: number; // hours
}

export interface ShopLocation {
  address: string;
  city: string;
  state: string;
  country: CountryCode;
  postalCode: string;
  coordinates?: { latitude: number; longitude: number };
  isPhysicalStore: boolean;
  storeHours?: ShopBusinessHours[];
}

export interface ShopStats {
  totalProducts: number;
  activeProducts: number;
  totalOrders: number;
  totalRevenue: number;
  averageOrderValue: number;
  averageRating: number;
  totalReviews: number;
  responseTime: number; // hours
  responseRate: number; // percentage
  returnRate: number; // percentage
  customerRetentionRate: number; // percentage
  monthlyStats: MonthlyShopStats[];
}

export interface MonthlyShopStats {
  month: string; // YYYY-MM
  orders: number;
  revenue: number;
  newCustomers: number;
  returningCustomers: number;
  averageOrderValue: number;
  topProducts: { productId: UUID; sales: number }[];
}

export interface ShopPolicies {
  returnPolicy?: string;
  shippingPolicy?: string;
  privacyPolicy?: string;
  termsOfService?: string;
  refundPolicy?: string;
  warrantyPolicy?: string;
  lastUpdated: Timestamp;
}

export interface ShopSubscription {
  plan: 'free' | 'basic' | 'professional' | 'enterprise';
  status: 'active' | 'cancelled' | 'expired' | 'suspended';
  billingCycle: 'monthly' | 'yearly';
  price: number;
  currency: Currency;
  features: ShopFeature[];
  limits: ShopLimits;
  startDate: Timestamp;
  endDate?: Timestamp;
  autoRenew: boolean;
}

export interface ShopFeature {
  id: string;
  name: string;
  isEnabled: boolean;
  limit?: number;
  usage?: number;
}

export interface ShopLimits {
  maxProducts: number;
  maxCategories: number;
  maxImages: number;
  maxBandwidth: number; // GB
  maxOrders: number; // per month
  customDomain: boolean;
  advancedAnalytics: boolean;
  prioritySupport: boolean;
}

export interface ShopCompliance {
  businessLicense?: BusinessLicense;
  taxRegistration?: TaxRegistration;
  certifications: ShopCertification[];
  insurancePolicies: InsurancePolicy[];
  complianceChecks: ComplianceCheck[];
}

export interface BusinessLicense {
  number: string;
  issuer: string;
  issuedAt: Timestamp;
  expiresAt: Timestamp;
  documentUrl?: string;
  status: 'valid' | 'expired' | 'suspended';
}

export interface TaxRegistration {
  number: string;
  country: CountryCode;
  registeredAt: Timestamp;
  isActive: boolean;
}

export interface ShopCertification {
  type: string;
  issuer: string;
  validFrom: Timestamp;
  validTo: Timestamp;
  documentUrl?: string;
}

export interface InsurancePolicy {
  type: 'liability' | 'product' | 'cyber' | 'general';
  provider: string;
  policyNumber: string;
  coverage: number;
  currency: Currency;
  validFrom: Timestamp;
  validTo: Timestamp;
}

export interface ComplianceCheck {
  id: UUID;
  type: string;
  status: 'pending' | 'passed' | 'failed' | 'review_required';
  checkedAt: Timestamp;
  nextCheckAt?: Timestamp;
  notes?: string;
  documentUrls: string[];
}

export interface ShopLocalization {
  [locale: string]: {
    name: string;
    description?: string;
    policies?: Partial<ShopPolicies>;
    categories: string[];
  };
}

// =============================================================================
// CART & ORDERING
// =============================================================================

export interface Cart {
  id: UUID;
  userId: UUID;
  sessionId?: string; // For guest users
  items: CartItem[];
  subtotal: number;
  discount: number;
  tax: number;
  shipping: number;
  total: number;
  currency: Currency;
  couponCode?: string;
  shippingAddress?: Address;
  billingAddress?: Address;
  paymentMethod?: PaymentMethodInfo;
  notes?: string;
  metadata?: CartMetadata;
  expiresAt?: Timestamp;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface CartItem {
  id: UUID;
  cartId: UUID;
  productId: UUID;
  variantId?: UUID;
  quantity: number;
  price: number;
  compareAtPrice?: number;
  title: string;
  image?: string;
  variant?: string;
  shopId: UUID;
  shopName: string;
  isAvailable: boolean;
  isDigital: boolean;
  shippingRequired: boolean;
  weight?: number;
  dimensions?: ProductDimensions;
  customizations?: ProductCustomization[];
  addedAt: Timestamp;
  updatedAt: Timestamp;
}

export interface ProductCustomization {
  type: 'text' | 'image' | 'color' | 'option';
  label: string;
  value: string;
  price?: number; // Additional cost
  displayOrder: number;
}

export interface CartMetadata {
  promoCode?: string;
  referralCode?: string;
  affiliateId?: UUID;
  utmSource?: string;
  utmMedium?: string;
  utmCampaign?: string;
  deviceInfo?: string;
  ipAddress?: string;
  geoLocation?: string;
}

export interface Order {
  id: UUID;
  orderNumber: string;
  userId: UUID;
  customerInfo: CustomerInfo;
  status: OrderStatus;
  items: OrderItem[];
  subtotal: number;
  discount: number;
  tax: number;
  shipping: number;
  total: number;
  currency: Currency;
  couponCode?: string;
  shippingAddress: Address;
  billingAddress?: Address;
  payment: OrderPayment;
  fulfillment: OrderFulfillment;
  communications: OrderCommunication[];
  timeline: OrderTimeline[];
  refunds: OrderRefund[];
  notes?: string;
  metadata?: OrderMetadata;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'processing'
  | 'shipped'
  | 'delivered'
  | 'cancelled'
  | 'refunded'
  | 'disputed';

export interface CustomerInfo {
  name: string;
  email?: string;
  phone?: string;
  preferredLanguage: Locale;
  isGuest: boolean;
  customerNotes?: string;
}

export interface OrderItem {
  id: UUID;
  orderId: UUID;
  productId: UUID;
  variantId?: UUID;
  quantity: number;
  price: number;
  total: number;
  title: string;
  image?: string;
  variant?: string;
  shopId: UUID;
  shopName: string;
  sku?: string;
  weight?: number;
  customizations?: ProductCustomization[];
  fulfillmentStatus: 'pending' | 'processing' | 'shipped' | 'delivered' | 'cancelled' | 'returned';
  trackingNumber?: string;
  returnStatus?: 'none' | 'requested' | 'approved' | 'returned' | 'refunded';
}

export interface OrderPayment {
  method: PaymentMethodInfo;
  status: PaymentStatus;
  transactionId?: string;
  amount: number;
  currency: Currency;
  fees: PaymentFee[];
  processedAt?: Timestamp;
  failureReason?: string;
  authorization?: PaymentAuthorization;
  refunds: PaymentRefund[];
}

export type PaymentStatus = 'pending' | 'authorized' | 'captured' | 'failed' | 'refunded' | 'disputed';

export interface PaymentMethodInfo {
  type: 'wallet' | 'bank_transfer' | 'credit_card' | 'qr_payment' | 'cod' | 'installment';
  details: PaymentMethodDetails;
  isStored: boolean;
  metadata?: Record<string, any>;
}

export interface InstallmentPlan {
  provider: string;
  months: number;
  monthlyAmount: number;
  totalAmount: number;
  interestRate: number;
  firstPayment: Timestamp;
}

export interface PaymentFee {
  type: 'processing' | 'gateway' | 'currency_conversion' | 'installment';
  amount: number;
  currency: Currency;
  description: string;
}

export interface PaymentAuthorization {
  authorizationId: string;
  amount: number;
  currency: Currency;
  expiresAt: Timestamp;
  capturedAmount?: number;
  capturedAt?: Timestamp;
}

export interface PaymentRefund {
  id: UUID;
  amount: number;
  currency: Currency;
  reason: string;
  status: 'pending' | 'completed' | 'failed';
  refundId?: string;
  processedAt?: Timestamp;
  failureReason?: string;
}

export interface OrderFulfillment {
  status: 'pending' | 'processing' | 'partially_shipped' | 'shipped' | 'delivered' | 'cancelled';
  trackingNumbers: TrackingInfo[];
  estimatedDelivery?: Timestamp;
  actualDelivery?: Timestamp;
  shippingCarrier?: string;
  shippingMethod?: string;
  shippingCost: number;
  packaging?: PackagingInfo;
  deliveryInstructions?: string;
  deliveryAttempts: DeliveryAttempt[];
}

export interface TrackingInfo {
  trackingNumber: string;
  carrier: string;
  url?: string;
  status: 'created' | 'picked_up' | 'in_transit' | 'out_for_delivery' | 'delivered' | 'failed';
  lastUpdate: Timestamp;
  estimatedDelivery?: Timestamp;
  items: UUID[]; // OrderItem IDs
}

export interface PackagingInfo {
  type: 'envelope' | 'box' | 'tube' | 'custom';
  dimensions?: ProductDimensions;
  weight: number;
  materials: string[];
  isEcoFriendly: boolean;
}

export interface DeliveryAttempt {
  attemptNumber: number;
  attemptedAt: Timestamp;
  status: 'failed' | 'delivered' | 'rescheduled';
  reason?: string;
  nextAttempt?: Timestamp;
  signature?: string;
  photo?: string;
}

export interface OrderCommunication {
  id: UUID;
  orderId: UUID;
  type: 'email' | 'sms' | 'push' | 'in_app';
  direction: 'outbound' | 'inbound';
  subject?: string;
  content: string;
  sentAt: Timestamp;
  deliveredAt?: Timestamp;
  readAt?: Timestamp;
  responseRequired: boolean;
  templateId?: string;
}

export interface OrderTimeline {
  id: UUID;
  orderId: UUID;
  status: string;
  description: string;
  timestamp: Timestamp;
  userId?: UUID;
  metadata?: Record<string, any>;
  isPublic: boolean;
  notification?: {
    sent: boolean;
    channels: string[];
    sentAt?: Timestamp;
  };
}

export interface OrderRefund {
  id: UUID;
  orderId: UUID;
  amount: number;
  currency: Currency;
  reason: RefundReason;
  status: 'requested' | 'approved' | 'processing' | 'completed' | 'rejected';
  items: RefundItem[];
  refundMethod: 'original_payment' | 'store_credit' | 'bank_transfer';
  processedAt?: Timestamp;
  notes?: string;
  attachments: string[];
}

export type RefundReason =
  | 'defective_product'
  | 'wrong_item'
  | 'not_as_described'
  | 'arrived_late'
  | 'customer_changed_mind'
  | 'duplicate_order'
  | 'fraud'
  | 'other';

export interface RefundItem {
  orderItemId: UUID;
  quantity: number;
  amount: number;
  reason?: string;
  condition?: 'new' | 'opened' | 'used' | 'damaged';
  photos?: string[];
}

export interface OrderMetadata {
  source: 'web' | 'mobile' | 'api' | 'admin';
  referrer?: string;
  utmParams?: Record<string, string>;
  deviceInfo?: string;
  ipAddress?: string;
  fraudScore?: number;
  riskLevel?: 'low' | 'medium' | 'high';
  affiliateId?: UUID;
  promotionIds: UUID[];
}

export interface Address {
  id?: UUID;
  firstName: string;
  lastName: string;
  company?: string;
  address1: string;
  address2?: string;
  city: string;
  province: string;
  country: CountryCode;
  postalCode: string;
  phone?: string;
  email?: string;
  isDefault: boolean;
  type: 'shipping' | 'billing' | 'both';
  coordinates?: { latitude: number; longitude: number };
  deliveryInstructions?: string;
  accessCodes?: string;
  validatedAt?: Timestamp;
  validationStatus?: 'valid' | 'invalid' | 'unverified';
}

// =============================================================================
// BUSINESS LOGIC CONSTANTS
// =============================================================================

/**
 * Product status transitions
 */
export const PRODUCT_STATUS_TRANSITIONS: Record<ProductStatus, ProductStatus[]> = {
  draft: ['active', 'archived'],
  active: ['out_of_stock', 'discontinued', 'archived'],
  out_of_stock: ['active', 'discontinued', 'archived'],
  discontinued: ['active', 'archived'],
  archived: ['draft']
};

/**
 * Order status transitions
 */
export const ORDER_STATUS_TRANSITIONS: Record<OrderStatus, OrderStatus[]> = {
  pending: ['confirmed', 'cancelled'],
  confirmed: ['processing', 'cancelled'],
  processing: ['shipped', 'cancelled'],
  shipped: ['delivered', 'cancelled'],
  delivered: ['refunded', 'disputed'],
  cancelled: [],
  refunded: ['disputed'],
  disputed: []
};

/**
 * Currency configurations for SEA markets
 */
export const CURRENCY_CONFIG: Record<Currency, {
  symbol: string;
  decimals: number;
  placement: 'before' | 'after';
  countries: CountryCode[];
}> = {
  THB: { symbol: '฿', decimals: 2, placement: 'before', countries: ['TH'] },
  IDR: { symbol: 'Rp', decimals: 0, placement: 'before', countries: ['ID'] },
  MYR: { symbol: 'RM', decimals: 2, placement: 'before', countries: ['MY'] },
  SGD: { symbol: 'S$', decimals: 2, placement: 'before', countries: ['SG'] },
  PHP: { symbol: '₱', decimals: 2, placement: 'before', countries: ['PH'] },
  VND: { symbol: '₫', decimals: 0, placement: 'after', countries: ['VN'] },
  USD: { symbol: '$', decimals: 2, placement: 'before', countries: [] }
};

/**
 * Default shipping methods by country
 */
export const DEFAULT_SHIPPING_METHODS: Record<CountryCode, string[]> = {
  TH: ['Thailand Post', 'Kerry Express', 'J&T Express', 'Flash Express'],
  ID: ['JNE', 'TIKI', 'Pos Indonesia', 'J&T Express'],
  MY: ['Pos Malaysia', 'City-Link', 'Ninja Van', 'J&T Express'],
  SG: ['SingPost', 'Ninja Van', 'Qxpress', 'J&T Express'],
  PH: ['LBC', 'J&T Express', '2GO Express', 'Ninja Van'],
  VN: ['Vietnam Post', 'Giao Hang Nhanh', 'Viettel Post', 'J&T Express']
};

// =============================================================================
// WALLET & PAYMENT DOMAIN
// =============================================================================

export interface Wallet {
  id: UUID;
  userId: UUID;
  balance: number;
  currency: Currency;
  frozenBalance: number;
  availableBalance: number;
  dailyLimit: number;
  monthlyLimit: number;
  usedThisMonth: number;
  status: 'active' | 'suspended' | 'closed';
  isPrimary: boolean;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface Transaction {
  id: UUID;
  walletId: UUID;
  type: TransactionType;
  amount: number;
  currency: Currency;
  fee: number;
  netAmount: number;
  status: TransactionStatus;
  description: string;
  reference?: string;
  metadata: TransactionMetadata;
  counterpart?: TransactionCounterpart;
  category?: TransactionCategory;
  tags: string[];
  balanceBefore?: number;
  balanceAfter?: number;
  createdAt: Timestamp;
  processedAt?: Timestamp;
  completedAt?: Timestamp;
}

export type TransactionType =
  | 'send'
  | 'receive'
  | 'topup'
  | 'withdraw'
  | 'purchase'
  | 'refund'
  | 'fee'
  | 'reward'
  | 'cashback';

export type TransactionStatus =
  | 'pending'
  | 'processing'
  | 'completed'
  | 'failed'
  | 'cancelled'
  | 'expired';

export type TransactionCategory =
  | 'food'
  | 'transport'
  | 'shopping'
  | 'entertainment'
  | 'bills'
  | 'transfer'
  | 'topup'
  | 'other';

export interface TransactionMetadata {
  orderId?: UUID;
  productId?: UUID;
  eventId?: UUID;
  dialogId?: UUID;
  messageId?: UUID;
  location?: Location;
  merchant?: string;
  paymentMethod?: string;
  promotionId?: UUID;
}

export interface TransactionCounterpart {
  id?: UUID;
  name: string;
  avatar?: string;
  type: 'user' | 'merchant' | 'system';
  identifier?: string; // phone, email, etc.
}

export interface PaymentMethod {
  id: UUID;
  userId: UUID;
  type: PaymentMethodType;
  name: string;
  isDefault: boolean;
  isActive: boolean;
  details: PaymentMethodDetails;
  verificationStatus: 'pending' | 'verified' | 'failed';
  expiresAt?: Timestamp;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type PaymentMethodType =
  | 'wallet'
  | 'bank_account'
  | 'credit_card'
  | 'debit_card'
  | 'promptpay'
  | 'qr_code'
  | 'bank_transfer'
  | 'e_wallet';

export interface PaymentMethodDetails {
  // Credit/Debit Card
  cardLast4?: string;
  cardBrand?: string;
  cardExpiry?: string;
  last4?: string;
  brand?: string;
  expiryMonth?: number;
  expiryYear?: number;

  // Bank Account
  bankName?: string;
  accountNumber?: string;
  accountName?: string;
  routingNumber?: string;
  swiftCode?: string;
  accountLast4?: string;

  // PromptPay / QR Payments
  promptpayId?: string;
  promptpayType?: 'mobile' | 'id_card' | 'e_wallet';
  qrProvider?: string;
  qrReference?: string;
  qrData?: string;
  qrImageUrl?: string;

  // Cash on Delivery
  codFee?: number;

  // E-Wallet
  walletProvider?: string;
  walletId?: string;
  phoneNumber?: string;
  eWalletId?: string;

  // Installment
  installmentPlan?: InstallmentPlan;

  // Verification
  verifiedAt?: Timestamp;
  verificationMethod?: 'otp' | 'document' | 'manual';
  referenceId?: string;
}

export interface QRPayment {
  id: UUID;
  userId: UUID;
  amount: number;
  currency: Currency;
  description?: string;
  qrCode: string;
  qrImageUrl: string;
  expiresAt: Timestamp;
  status: 'active' | 'used' | 'expired' | 'cancelled';
  usedBy?: UUID;
  usedAt?: Timestamp;
  createdAt: Timestamp;
}

// =============================================================================
// EVENT DOMAIN
// =============================================================================

export interface Event {
  id: UUID;
  title: string;
  description: string;
  shortDescription?: string;
  category: EventCategory;
  type: EventType;
  status: EventStatus;
  organizer: EventOrganizer;
  venue: EventVenue;
  schedule: EventSchedule;
  ticketing: EventTicketing;
  media: EventMedia;
  lineup?: EventLineup;
  amenities: string[];
  ageRestriction?: string;
  tags: string[];
  popularity: EventPopularity;
  socialProof: EventSocialProof;
  reviews: EventReview[];
  weather?: EventWeather;
  isPromoted: boolean;
  isFeatured: boolean;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type EventCategory =
  | 'music'
  | 'food'
  | 'cultural'
  | 'festival'
  | 'temple'
  | 'market'
  | 'sports'
  | 'technology'
  | 'business'
  | 'art'
  | 'education';

export type EventType =
  | 'concert'
  | 'festival'
  | 'conference'
  | 'workshop'
  | 'exhibition'
  | 'competition'
  | 'celebration'
  | 'ceremony';

export type EventStatus =
  | 'draft'
  | 'published'
  | 'cancelled'
  | 'postponed'
  | 'completed';

export interface EventOrganizer {
  id: UUID;
  name: string;
  avatar?: string;
  isVerified: boolean;
  pastEvents: number;
  rating: number;
  contact: {
    email?: string;
    phone?: string;
    website?: string;
  };
}

export interface EventVenue {
  id?: UUID;
  name: string;
  address: string;
  location: Location;
  capacity: number;
  facilities: string[];
  accessibility: string[];
  parking: boolean;
  publicTransport: string[];
}

export interface EventSchedule {
  startDate: Timestamp;
  endDate: Timestamp;
  timezone: string;
  sessions?: EventSession[];
  doors?: Timestamp;
  curfew?: Timestamp;
}

export interface EventSession {
  id: UUID;
  title: string;
  description?: string;
  startTime: Timestamp;
  endTime: Timestamp;
  stage?: string;
  speakers?: EventSpeaker[];
}

export interface EventSpeaker {
  id: UUID;
  name: string;
  title?: string;
  company?: string;
  avatar?: string;
  bio?: string;
}

export interface EventTicketing {
  available: boolean;
  ticketTypes: TicketType[];
  salesStart?: Timestamp;
  salesEnd: Timestamp;
  terms?: string;
  refundPolicy?: string;
}

export interface TicketType {
  id: UUID;
  name: string;
  description?: string;
  price: number;
  currency: Currency;
  quantity: number;
  sold: number;
  available: number;
  perks: string[];
  transferable: boolean;
  refundable: boolean;
  salesStart?: Timestamp;
  salesEnd?: Timestamp;
  accessLevel: 'general' | 'vip' | 'backstage' | 'premium';
}

export interface EventMedia {
  coverImage: string;
  gallery: EventImage[];
  videos: EventVideo[];
  livestreamUrl?: string;
}

export interface EventImage {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  caption?: string;
  order: number;
  type: 'cover' | 'gallery' | 'venue' | 'lineup';
}

export interface EventVideo {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  title?: string;
  duration: number;
  type: 'promo' | 'recap' | 'livestream';
}

export interface EventLineup {
  headliners: Artist[];
  supporting: Artist[];
  schedule?: LineupSchedule[];
}

export interface Artist {
  id: UUID;
  name: string;
  image?: string;
  genre: string;
  popularity: number;
  socialMedia?: { [platform: string]: string };
  isHeadliner: boolean;
}

export interface LineupSchedule {
  day: string;
  stages: LineupStage[];
}

export interface LineupStage {
  name: string;
  location?: string;
  acts: LineupAct[];
}

export interface LineupAct {
  time: string;
  artist: string;
  duration: string;
  type?: 'performance' | 'dj_set' | 'talk' | 'workshop';
}

export interface EventPopularity {
  trending: boolean;
  rank?: number;
  interest: number;
  attending: number;
  views: number;
  shares: number;
  saves: number;
}

export interface EventSocialProof {
  friendsGoing: EventAttendee[];
  influencersGoing: EventInfluencer[];
  totalFriendsGoing: number;
  recommendedBy: UUID[];
}

export interface EventAttendee {
  userId: UUID;
  name: string;
  avatar?: string;
  status: 'interested' | 'going' | 'maybe';
  ticketType?: string;
  checkedIn: boolean;
  checkedInAt?: Timestamp;
}

export interface EventInfluencer {
  userId: UUID;
  name: string;
  avatar?: string;
  isVerified: boolean;
  followers: number;
  category: string;
}

export interface EventReview {
  id: UUID;
  eventId: UUID;
  userId: UUID;
  rating: number;
  title?: string;
  comment: string;
  images?: string[];
  attendedDate?: Timestamp;
  isVerifiedAttendee: boolean;
  helpfulCount: number;
  createdAt: Timestamp;
}

export interface EventWeather {
  forecast: string;
  temperature: string;
  humidity?: number;
  precipitation?: number;
  recommendation: string;
  updatedAt: Timestamp;
}

export interface Ticket {
  id: UUID;
  eventId: UUID;
  userId: UUID;
  ticketTypeId: UUID;
  ticketNumber: string;
  qrCode: string;
  status: TicketStatus;
  price: number;
  currency: Currency;
  purchasedAt: Timestamp;
  transferredTo?: UUID;
  transferredAt?: Timestamp;
  checkedIn: boolean;
  checkedInAt?: Timestamp;
  refunded: boolean;
  refundedAt?: Timestamp;
  refundAmount?: number;
}

export type TicketStatus =
  | 'active'
  | 'transferred'
  | 'refunded'
  | 'cancelled'
  | 'expired'
  | 'used';

// =============================================================================
// WORKSPACE & BUSINESS DOMAIN
// =============================================================================

export interface Workspace {
  id: UUID;
  name: string;
  description?: string;
  avatar?: string;
  type: 'personal' | 'team' | 'enterprise';
  status: 'active' | 'suspended' | 'archived';
  ownerId: UUID;
  settings: WorkspaceSettings;
  subscription: WorkspaceSubscription;
  features: WorkspaceFeature[];
  stats: WorkspaceStats;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface WorkspaceSettings {
  isPublic: boolean;
  allowInvites: boolean;
  requireApproval: boolean;
  timezone: string;
  locale: Locale;
  currency: Currency;
  branding: WorkspaceBranding;
}

export interface WorkspaceBranding {
  logo?: string;
  primaryColor?: string;
  secondaryColor?: string;
  customDomain?: string;
}

export interface WorkspaceSubscription {
  plan: 'free' | 'starter' | 'professional' | 'enterprise';
  status: 'active' | 'cancelled' | 'expired' | 'trial';
  billingCycle: 'monthly' | 'yearly';
  amount: number;
  currency: Currency;
  trialEndsAt?: Timestamp;
  renewsAt?: Timestamp;
  cancelledAt?: Timestamp;
}

export interface WorkspaceFeature {
  id: string;
  name: string;
  isEnabled: boolean;
  config?: Record<string, any>;
  usageLimit?: number;
  usageCount?: number;
}

export interface WorkspaceStats {
  memberCount: number;
  projectCount: number;
  messageCount: number;
  storageUsed: number; // bytes
  storageLimit: number; // bytes
  apiCalls: number;
  lastActivity: Timestamp;
}

export interface WorkspaceMember {
  id: UUID;
  workspaceId: UUID;
  userId: UUID;
  role: WorkspaceRole;
  permissions: WorkspacePermission[];
  status: 'active' | 'invited' | 'suspended';
  joinedAt: Timestamp;
  invitedBy?: UUID;
  lastSeen?: Timestamp;
}

export type WorkspaceRole = 'owner' | 'admin' | 'member' | 'guest' | 'viewer';

export interface WorkspacePermission {
  resource: string;
  actions: string[];
}

export interface Project {
  id: UUID;
  workspaceId: UUID;
  name: string;
  description?: string;
  avatar?: string;
  status: ProjectStatus;
  priority: 'low' | 'medium' | 'high' | 'urgent';
  ownerId: UUID;
  members: UUID[];
  settings: ProjectSettings;
  progress: ProjectProgress;
  budget?: ProjectBudget;
  timeline: ProjectTimeline;
  tags: string[];
  attachments: ProjectAttachment[];
  createdAt: Timestamp;
  updatedAt: Timestamp;
  completedAt?: Timestamp;
}

export type ProjectStatus =
  | 'planning'
  | 'active'
  | 'on_hold'
  | 'completed'
  | 'cancelled'
  | 'archived';

export interface ProjectSettings {
  isPublic: boolean;
  allowComments: boolean;
  requireApproval: boolean;
  template?: string;
  methodology: 'agile' | 'waterfall' | 'kanban' | 'scrum' | 'custom';
}

export interface ProjectProgress {
  percentage: number;
  tasksTotal: number;
  tasksCompleted: number;
  milestonesTotal: number;
  milestonesCompleted: number;
  hoursEstimated?: number;
  hoursLogged?: number;
}

export interface ProjectBudget {
  total: number;
  spent: number;
  currency: Currency;
  breakdown: BudgetCategory[];
}

export interface BudgetCategory {
  name: string;
  allocated: number;
  spent: number;
  percentage: number;
}

export interface ProjectTimeline {
  startDate: Timestamp;
  endDate: Timestamp;
  milestones: ProjectMilestone[];
  phases: ProjectPhase[];
}

export interface ProjectMilestone {
  id: UUID;
  name: string;
  description?: string;
  dueDate: Timestamp;
  status: 'pending' | 'completed' | 'overdue';
  completedAt?: Timestamp;
  dependencies: UUID[];
}

export interface ProjectPhase {
  id: UUID;
  name: string;
  description?: string;
  startDate: Timestamp;
  endDate: Timestamp;
  status: 'upcoming' | 'active' | 'completed' | 'delayed';
  tasks: UUID[];
}

export interface ProjectAttachment {
  id: UUID;
  projectId: UUID;
  name: string;
  fileUrl: string;
  fileSize: number;
  mimeType: string;
  uploadedBy: UUID;
  uploadedAt: Timestamp;
}

export interface Task {
  id: UUID;
  projectId: UUID;
  title: string;
  description?: string;
  status: TaskStatus;
  priority: 'low' | 'medium' | 'high' | 'urgent';
  assigneeId?: UUID;
  creatorId: UUID;
  reviewerId?: UUID;
  labels: string[];
  estimatedHours?: number;
  loggedHours?: number;
  startDate?: Timestamp;
  dueDate?: Timestamp;
  completedAt?: Timestamp;
  dependencies: UUID[];
  subtasks: UUID[];
  attachments: TaskAttachment[];
  comments: TaskComment[];
  history: TaskHistory[];
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type TaskStatus =
  | 'backlog'
  | 'todo'
  | 'in_progress'
  | 'in_review'
  | 'testing'
  | 'done'
  | 'cancelled';

export interface TaskAttachment {
  id: UUID;
  taskId: UUID;
  name: string;
  fileUrl: string;
  fileSize: number;
  mimeType: string;
  uploadedBy: UUID;
  uploadedAt: Timestamp;
}

export interface TaskComment {
  id: UUID;
  taskId: UUID;
  authorId: UUID;
  content: string;
  attachments?: string[];
  createdAt: Timestamp;
  editedAt?: Timestamp;
}

export interface TaskHistory {
  id: UUID;
  taskId: UUID;
  userId: UUID;
  action: string;
  field?: string;
  oldValue?: string;
  newValue?: string;
  createdAt: Timestamp;
}

export interface Invoice {
  id: UUID;
  workspaceId: UUID;
  invoiceNumber: string;
  customerId?: UUID;
  customerInfo: InvoiceCustomer;
  status: InvoiceStatus;
  items: InvoiceItem[];
  subtotal: number;
  tax: number;
  discount: number;
  total: number;
  currency: Currency;
  dueDate: Timestamp;
  issuedDate: Timestamp;
  paidDate?: Timestamp;
  notes?: string;
  terms?: string;
  paymentMethod?: PaymentMethod;
  attachments: InvoiceAttachment[];
  reminder: InvoiceReminder;
  createdBy: UUID;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type InvoiceStatus =
  | 'draft'
  | 'sent'
  | 'viewed'
  | 'overdue'
  | 'paid'
  | 'cancelled'
  | 'refunded';

export interface InvoiceCustomer {
  name: string;
  email?: string;
  phone?: string;
  company?: string;
  address?: Address;
  taxId?: string;
}

export interface InvoiceItem {
  id: UUID;
  description: string;
  quantity: number;
  rate: number;
  amount: number;
  taxable: boolean;
  category?: string;
}

export interface InvoiceAttachment {
  id: UUID;
  invoiceId: UUID;
  name: string;
  fileUrl: string;
  fileSize: number;
  mimeType: string;
  uploadedAt: Timestamp;
}

export interface InvoiceReminder {
  isEnabled: boolean;
  daysBefore: number[];
  lastSent?: Timestamp;
  nextSend?: Timestamp;
}

// =============================================================================
// VIDEO & STREAMING DOMAIN
// =============================================================================

export interface Video {
  id: UUID;
  channelId: UUID;
  title: string;
  description?: string;
  thumbnail: string;
  url: string;
  duration: number; // seconds
  views: number;
  likes: number;
  dislikes: number;
  comments: number;
  shares: number;
  category: VideoCategory;
  tags: string[];
  language: Locale;
  privacy: 'public' | 'unlisted' | 'private';
  status: 'processing' | 'published' | 'deleted' | 'blocked';
  monetized: boolean;
  ageRestricted: boolean;
  location?: Location;
  recordedAt?: Timestamp;
  publishedAt: Timestamp;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type VideoCategory =
  | 'music'
  | 'gaming'
  | 'education'
  | 'entertainment'
  | 'news'
  | 'sports'
  | 'technology'
  | 'food'
  | 'travel'
  | 'lifestyle'
  | 'business'
  | 'health';

export interface Channel {
  id: UUID;
  ownerId: UUID;
  name: string;
  description?: string;
  avatar?: string;
  banner?: string;
  isVerified: boolean;
  subscriberCount: number;
  videoCount: number;
  totalViews: number;
  totalWatchTime: number; // minutes
  category: ChannelCategory;
  country: CountryCode;
  language: Locale;
  settings: ChannelSettings;
  monetization: ChannelMonetization;
  stats: ChannelStats;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type ChannelCategory =
  | 'personal'
  | 'business'
  | 'entertainment'
  | 'education'
  | 'music'
  | 'gaming'
  | 'news'
  | 'sports';

export interface ChannelSettings {
  allowComments: boolean;
  allowSubscriptions: boolean;
  showStats: boolean;
  moderateComments: boolean;
  defaultVideoPrivacy: 'public' | 'unlisted' | 'private';
  brandingWatermark?: string;
}

export interface ChannelMonetization {
  isEnabled: boolean;
  adRevenue: number;
  currency: Currency;
  threshold: number;
  lastPayout?: Timestamp;
  totalEarnings: number;
}

export interface ChannelStats {
  dailyViews: number;
  weeklyViews: number;
  monthlyViews: number;
  dailySubscribers: number;
  weeklySubscribers: number;
  monthlySubscribers: number;
  averageViewDuration: number; // seconds
  engagementRate: number; // percentage
}

export interface LiveStream {
  id: UUID;
  channelId: UUID;
  title: string;
  description?: string;
  thumbnail?: string;
  streamUrl: string;
  chatUrl?: string;
  status: LiveStreamStatus;
  viewerCount: number;
  maxViewers: number;
  likes: number;
  dislikes: number;
  chatEnabled: boolean;
  category: VideoCategory;
  tags: string[];
  language: Locale;
  scheduledAt?: Timestamp;
  startedAt?: Timestamp;
  endedAt?: Timestamp;
  duration?: number; // seconds
  recording: StreamRecording;
  quality: StreamQuality;
  monetization: StreamMonetization;
  moderation: StreamModeration;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type LiveStreamStatus =
  | 'scheduled'
  | 'live'
  | 'ended'
  | 'cancelled'
  | 'error';

export interface StreamRecording {
  isEnabled: boolean;
  autoPublish: boolean;
  videoId?: UUID;
  fileUrl?: string;
  duration?: number;
}

export interface StreamQuality {
  resolution: '720p' | '1080p' | '1440p' | '2160p';
  bitrate: number;
  framerate: number;
  codec: string;
  adaptiveBitrate: boolean;
}

export interface StreamMonetization {
  isEnabled: boolean;
  superChatEnabled: boolean;
  membershipsEnabled: boolean;
  adsEnabled: boolean;
  donationsEnabled: boolean;
}

export interface StreamModeration {
  slowMode: boolean;
  slowModeDelay: number; // seconds
  subscribersOnly: boolean;
  wordsBlacklist: string[];
  moderators: UUID[];
  bannedUsers: UUID[];
}

export interface Subscription {
  id: UUID;
  subscriberId: UUID;
  channelId: UUID;
  notificationsEnabled: boolean;
  tier: 'free' | 'premium' | 'member';
  isPublic: boolean;
  subscribedAt: Timestamp;
  unsubscribedAt?: Timestamp;
}

export interface VideoView {
  id: UUID;
  videoId: UUID;
  viewerId?: UUID;
  sessionId: string;
  watchTime: number; // seconds
  percentage: number;
  quality: string;
  device: string;
  location?: Location;
  referrer?: string;
  createdAt: Timestamp;
}

export interface VideoComment {
  id: UUID;
  videoId: UUID;
  authorId: UUID;
  parentId?: UUID; // for replies
  content: string;
  likes: number;
  dislikes: number;
  replies: number;
  isPinned: boolean;
  isHeartedByCreator: boolean;
  createdAt: Timestamp;
  editedAt?: Timestamp;
  deletedAt?: Timestamp;
}

// =============================================================================
// DISCOVERY & LOCATION DOMAIN
// =============================================================================

export interface Location {
  id?: UUID;
  latitude: number;
  longitude: number;
  address?: string;
  city?: string;
  state?: string;
  country?: CountryCode;
  postalCode?: string;
  formattedAddress?: string;
  placeId?: string;
  placeName?: string;
  category?: LocationCategory;
}

export type LocationCategory =
  | 'restaurant'
  | 'shopping'
  | 'entertainment'
  | 'transport'
  | 'hospital'
  | 'school'
  | 'hotel'
  | 'tourist_attraction'
  | 'park'
  | 'temple'
  | 'market';

export interface Place {
  id: UUID;
  name: string;
  description?: string;
  category: LocationCategory;
  subcategory?: string;
  location: Location;
  contact: PlaceContact;
  hours: BusinessHours[];
  rating: PlaceRating;
  photos: PlacePhoto[];
  amenities: string[];
  priceLevel: 1 | 2 | 3 | 4; // $ to $$$$
  tags: string[];
  isVerified: boolean;
  claimedBy?: UUID;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface PlaceContact {
  phone?: string;
  website?: string;
  email?: string;
  socialMedia?: { [platform: string]: string };
}

export interface PlaceRating {
  averageRating: number;
  totalReviews: number;
  distribution: { [key: number]: number };
}

export interface PlacePhoto {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  caption?: string;
  uploadedBy?: UUID;
  uploadedAt: Timestamp;
  isVerified: boolean;
}

export interface PlaceReview {
  id: UUID;
  placeId: UUID;
  userId: UUID;
  rating: number;
  title?: string;
  content: string;
  photos?: string[];
  visitDate?: Timestamp;
  isVerified: boolean;
  helpfulCount: number;
  reportCount: number;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface Discovery {
  id: UUID;
  type: DiscoveryType;
  title: string;
  description?: string;
  image?: string;
  targetId: UUID; // refers to event, place, product, etc.
  targetType: 'event' | 'place' | 'product' | 'channel' | 'post';
  category: string;
  tags: string[];
  location?: Location;
  priority: number;
  isSponsored: boolean;
  audience: DiscoveryAudience;
  schedule: DiscoverySchedule;
  metrics: DiscoveryMetrics;
  status: 'active' | 'paused' | 'completed' | 'cancelled';
  createdBy: UUID;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type DiscoveryType =
  | 'trending'
  | 'recommended'
  | 'nearby'
  | 'popular'
  | 'new'
  | 'sponsored'
  | 'featured';

export interface DiscoveryAudience {
  targetCountries: CountryCode[];
  targetCities?: string[];
  targetLanguages: Locale[];
  targetAge?: { min: number; max: number };
  targetGender?: 'male' | 'female' | 'all';
  targetInterests: string[];
  excludeUsers?: UUID[];
}

export interface DiscoverySchedule {
  startDate: Timestamp;
  endDate?: Timestamp;
  timezone: string;
  isActive: boolean;
}

export interface DiscoveryMetrics {
  impressions: number;
  clicks: number;
  conversions: number;
  clickThroughRate: number;
  conversionRate: number;
  engagement: number;
  reach: number;
  frequency: number;
}

// =============================================================================
// ACTIVITY & GAMIFICATION DOMAIN
// =============================================================================

export interface Activity {
  id: UUID;
  userId: UUID;
  type: ActivityType;
  targetId?: UUID;
  targetType?: string;
  description: string;
  points: number;
  metadata: ActivityMetadata;
  isPublic: boolean;
  location?: Location;
  createdAt: Timestamp;
}

export type ActivityType =
  | 'message_sent'
  | 'post_created'
  | 'post_liked'
  | 'comment_added'
  | 'event_joined'
  | 'product_purchased'
  | 'friend_added'
  | 'place_visited'
  | 'video_watched'
  | 'achievement_unlocked'
  | 'challenge_completed'
  | 'streak_maintained';

export interface ActivityMetadata {
  amount?: number;
  currency?: Currency;
  duration?: number;
  distance?: number;
  participants?: UUID[];
  tags?: string[];
  extra?: Record<string, any>;
}

export interface Achievement {
  id: UUID;
  name: string;
  description: string;
  icon: string;
  category: AchievementCategory;
  type: AchievementType;
  rarity: 'common' | 'rare' | 'epic' | 'legendary';
  points: number;
  requirements: AchievementRequirement[];
  rewards: AchievementReward[];
  isActive: boolean;
  createdAt: Timestamp;
}

export type AchievementCategory =
  | 'social'
  | 'shopping'
  | 'events'
  | 'exploration'
  | 'content'
  | 'engagement'
  | 'loyalty'
  | 'special';

export type AchievementType =
  | 'milestone'
  | 'streak'
  | 'collection'
  | 'challenge'
  | 'seasonal'
  | 'hidden';

export interface AchievementRequirement {
  type: string;
  target: number;
  current?: number;
  timeframe?: number; // days
  conditions?: Record<string, any>;
}

export interface AchievementReward {
  type: 'points' | 'badge' | 'discount' | 'feature' | 'item';
  value: number | string;
  description: string;
  expiresAt?: Timestamp;
}

export interface UserAchievement {
  id: UUID;
  userId: UUID;
  achievementId: UUID;
  progress: AchievementProgress[];
  isCompleted: boolean;
  completedAt?: Timestamp;
  claimedAt?: Timestamp;
  notifiedAt?: Timestamp;
}

export interface AchievementProgress {
  requirementIndex: number;
  current: number;
  target: number;
  percentage: number;
  lastUpdated: Timestamp;
}

export interface Challenge {
  id: UUID;
  name: string;
  description: string;
  image?: string;
  category: string;
  type: ChallengeType;
  difficulty: 'easy' | 'medium' | 'hard' | 'expert';
  duration: number; // days
  participants: number;
  maxParticipants?: number;
  requirements: ChallengeRequirement[];
  rewards: ChallengeReward[];
  leaderboard: ChallengeLeaderboard[];
  status: 'upcoming' | 'active' | 'completed' | 'cancelled';
  startDate: Timestamp;
  endDate: Timestamp;
  createdBy: UUID;
  createdAt: Timestamp;
}

export type ChallengeType =
  | 'individual'
  | 'team'
  | 'community'
  | 'daily'
  | 'weekly'
  | 'seasonal';

export interface ChallengeRequirement {
  description: string;
  target: number;
  unit: string;
  timeframe?: number; // days
}

export interface ChallengeReward {
  rank: number;
  type: 'points' | 'badge' | 'prize' | 'discount';
  value: number | string;
  description: string;
}

export interface ChallengeLeaderboard {
  rank: number;
  userId: UUID;
  score: number;
  progress: number; // percentage
  lastUpdated: Timestamp;
}

export interface UserChallenge {
  id: UUID;
  userId: UUID;
  challengeId: UUID;
  status: 'registered' | 'active' | 'completed' | 'abandoned';
  progress: ChallengeProgress[];
  score: number;
  rank?: number;
  joinedAt: Timestamp;
  completedAt?: Timestamp;
}

export interface ChallengeProgress {
  requirementIndex: number;
  current: number;
  target: number;
  percentage: number;
  lastUpdated: Timestamp;
}

export interface UserStats {
  userId: UUID;
  totalPoints: number;
  level: number;
  currentLevelPoints: number;
  nextLevelPoints: number;
  achievements: number;
  challengesCompleted: number;
  streakDays: number;
  longestStreak: number;
  lastActivity: Timestamp;
  stats: { [category: string]: number };
}

// =============================================================================
// NOTIFICATION DOMAIN
// =============================================================================

export interface Notification {
  id: UUID;
  userId: UUID;
  type: NotificationType;
  title: string;
  message: string;
  icon?: string;
  image?: string;
  actionUrl?: string;
  actionText?: string;
  category: NotificationCategory;
  priority: 'low' | 'normal' | 'high' | 'urgent';
  data: NotificationData;
  channels: NotificationChannel[];
  status: NotificationStatus;
  scheduledAt?: Timestamp;
  sentAt?: Timestamp;
  readAt?: Timestamp;
  clickedAt?: Timestamp;
  expiresAt?: Timestamp;
  createdAt: Timestamp;
}

export type NotificationType =
  | 'message'
  | 'friend_request'
  | 'post_like'
  | 'post_comment'
  | 'event_reminder'
  | 'payment_received'
  | 'order_update'
  | 'achievement'
  | 'system'
  | 'promotion';

export type NotificationCategory =
  | 'social'
  | 'transactional'
  | 'promotional'
  | 'system'
  | 'reminder'
  | 'achievement';

export type NotificationChannel =
  | 'push'
  | 'email'
  | 'sms'
  | 'in_app'
  | 'webhook';

export type NotificationStatus =
  | 'pending'
  | 'sent'
  | 'delivered'
  | 'failed'
  | 'cancelled';

export interface NotificationData {
  targetId?: UUID;
  targetType?: string;
  senderId?: UUID;
  amount?: number;
  currency?: Currency;
  metadata?: Record<string, any>;
}

export interface NotificationPreference {
  id: UUID;
  userId: UUID;
  category: NotificationCategory;
  channels: NotificationChannel[];
  isEnabled: boolean;
  schedule?: NotificationSchedule;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface NotificationSchedule {
  timezone: string;
  quietHours: {
    start: string; // HH:mm
    end: string; // HH:mm
  };
  daysOfWeek: number[]; // 0-6, Sunday = 0
  frequency: 'immediate' | 'batched' | 'daily' | 'weekly';
}

// =============================================================================
// SEARCH DOMAIN
// =============================================================================

export interface SearchQuery {
  id: UUID;
  userId?: UUID;
  query: string;
  type: SearchType;
  filters: SearchFilters;
  results: SearchResult[];
  resultCount: number;
  executionTime: number; // ms
  location?: Location;
  createdAt: Timestamp;
}

export type SearchType =
  | 'general'
  | 'user'
  | 'product'
  | 'event'
  | 'place'
  | 'video'
  | 'post'
  | 'message';

export interface SearchFilters {
  category?: string;
  location?: Location;
  priceRange?: { min: number; max: number };
  dateRange?: { start: Timestamp; end: Timestamp };
  rating?: number;
  verified?: boolean;
  distance?: number; // km
  sortBy?: 'relevance' | 'date' | 'price' | 'rating' | 'distance';
  sortOrder?: 'asc' | 'desc';
}

export interface SearchResult {
  id: UUID;
  type: string;
  title: string;
  description?: string;
  image?: string;
  url?: string;
  score: number;
  ranking: number;
  metadata: Record<string, any>;
  createdAt: Timestamp;
}

export interface SearchSuggestion {
  id: UUID;
  query: string;
  type: SearchType;
  popularity: number;
  isPromoted: boolean;
  location?: CountryCode;
  language: Locale;
  createdAt: Timestamp;
}
