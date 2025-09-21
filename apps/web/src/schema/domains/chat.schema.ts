/**
 * Chat & Messaging Domain Schema
 *
 * Handles dialogs, messages, attachments, reactions, and real-time messaging features
 * Supports business chats, rich media, voice messages, and Southeast Asian payment integrations
 */

import { UUID, Timestamp, Currency, Locale } from '../schema';

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
