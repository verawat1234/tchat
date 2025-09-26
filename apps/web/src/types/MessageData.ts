// T027 - MessageData and MessageType interfaces
/**
 * Core message entity with type discrimination and React-specific state management
 * Supports all 13 new web message component types for cross-platform parity
 */

import { ReactNode } from 'react';

// Core Message Entity
export interface MessageData {
  readonly id: string;
  readonly senderId: string;
  readonly senderName: string;
  readonly timestamp: Date;
  readonly type: MessageType;
  readonly isOwn: boolean;
  readonly content: MessageContent;
  readonly metadata?: MessageMetadata;
  readonly renderState?: MessageRenderState;
}

// Message Type Enum - All supported message types including 13 new web implementations
export enum MessageType {
  // Existing types (for compatibility)
  TEXT = 'text',
  MARKDOWN = 'markdown',
  INVOICE = 'invoice',
  BILL = 'bill',
  ORDER = 'order',
  STICKER = 'sticker',
  POLL = 'poll',
  CONTACT = 'contact',
  LOCATION = 'location',
  PAYMENT = 'payment',
  VOICE = 'voice',
  FILE = 'file',
  IMAGE = 'image',
  VIDEO = 'video',
  MUSIC = 'music',
  SYSTEM = 'system',

  // New types being implemented for web
  REPLY = 'reply',
  GIF = 'gif',
  DOCUMENT = 'document',
  SPREADSHEET = 'spreadsheet',
  VENUE = 'venue',
  QUIZ = 'quiz',
  SURVEY = 'survey',
  PRODUCT = 'product',
  EVENT = 'event',
  REMINDER = 'reminder',
  CALENDAR = 'calendar',
  STATUS = 'status',
  ANNOUNCEMENT = 'announcement',
  FORM = 'form',
  CARD = 'card',
  EMBED = 'embed'
}

// Discriminated union for message content based on type
export type MessageContent =
  | string // For simple text messages
  | ReplyContent
  | AnimatedGifContent
  | DocumentContent
  | SpreadsheetContent
  | VenueContent
  | QuizContent
  | SurveyContent
  | ProductContent
  | EventContent
  | ReminderContent
  | StatusUpdateContent
  | AnnouncementContent
  | FormContent
  | RichCardContent
  | EmbedContent;

// Message Metadata
export interface MessageMetadata {
  readonly replyToId?: string;
  readonly forwardedFrom?: string;
  readonly editedAt?: Date;
  readonly reactions?: MessageReaction[];
  readonly readBy?: string[];
  readonly priority?: MessagePriority;
  readonly expiresAt?: Date;
  readonly tags?: string[];
  readonly deliveryStatus?: DeliveryStatus;
  readonly errorInfo?: MessageErrorInfo;
}

// Message Reactions
export interface MessageReaction {
  readonly emoji: string;
  readonly userIds: string[];
  readonly count: number;
  readonly timestamp: Date;
}

// Message Priority
export enum MessagePriority {
  LOW = 'low',
  NORMAL = 'normal',
  HIGH = 'high',
  URGENT = 'urgent'
}

// Delivery Status
export enum DeliveryStatus {
  PENDING = 'pending',
  SENT = 'sent',
  DELIVERED = 'delivered',
  READ = 'read',
  FAILED = 'failed'
}

// Message Error Information
export interface MessageErrorInfo {
  readonly code: string;
  readonly message: string;
  readonly timestamp: Date;
  readonly retryable: boolean;
  readonly details?: Record<string, unknown>;
}

// Message Render State (for UI components)
export interface MessageRenderState {
  readonly isVisible: boolean;
  readonly isExpanded: boolean;
  readonly loadingState: LoadingState;
  readonly interactionState?: InteractionState;
  readonly validationErrors?: ValidationError[];
  readonly componentProps?: Record<string, unknown>;
}

// Loading States
export enum LoadingState {
  IDLE = 'idle',
  LOADING = 'loading',
  SUCCESS = 'success',
  ERROR = 'error'
}

// Interaction States for interactive messages
export interface InteractionState {
  readonly hasInteracted: boolean;
  readonly interactionCount: number;
  readonly lastInteraction?: Date;
  readonly pendingInteractions?: string[];
  readonly interactionResults?: InteractionResult[];
}

// Interaction Results
export interface InteractionResult {
  readonly interactionId: string;
  readonly type: InteractionType;
  readonly result: unknown;
  readonly timestamp: Date;
  readonly success: boolean;
  readonly errorMessage?: string;
}

// Interaction Types
export enum InteractionType {
  QUIZ_ANSWER = 'quiz_answer',
  SURVEY_RESPONSE = 'survey_response',
  EVENT_RSVP = 'event_rsvp',
  FORM_SUBMIT = 'form_submit',
  REACTION_ADD = 'reaction_add',
  REACTION_REMOVE = 'reaction_remove',
  STATUS_VIEW = 'status_view',
  CARD_ACTION = 'card_action',
  BUTTON_CLICK = 'button_click',
  LINK_CLICK = 'link_click'
}

// Validation Errors
export interface ValidationError {
  readonly field: string;
  readonly message: string;
  readonly code: string;
  readonly severity: ValidationSeverity;
}

export enum ValidationSeverity {
  INFO = 'info',
  WARNING = 'warning',
  ERROR = 'error'
}

// Forward declarations for content types (implemented in separate files)
export interface ReplyContent {
  readonly originalMessageId: string;
  readonly replyText: string;
  readonly threadDepth: number;
  readonly isThreadStart: boolean;
  readonly originalPreview?: MessagePreview;
}

export interface AnimatedGifContent {
  readonly url: string;
  readonly thumbnailUrl?: string;
  readonly width: number;
  readonly height: number;
  readonly duration?: number;
  readonly fileSize: number;
  readonly format: 'gif' | 'webp';
  readonly autoPlay: boolean;
  readonly loopCount: number;
}

export interface DocumentContent {
  readonly url: string;
  readonly fileName: string;
  readonly fileSize: number;
  readonly mimeType: string;
  readonly pageCount?: number;
  readonly thumbnailUrl?: string;
  readonly previewSupported: boolean;
  readonly textSearchable: boolean;
}

export interface SpreadsheetContent {
  readonly url: string;
  readonly fileName: string;
  readonly fileSize: number;
  readonly sheetNames: string[];
  readonly rowCount: number;
  readonly columnCount: number;
  readonly hasFormulas: boolean;
  readonly thumbnailUrl?: string;
}

export interface VenueContent {
  readonly id: string;
  readonly name: string;
  readonly address: string;
  readonly coordinates?: Coordinates;
  readonly category: string;
  readonly rating?: number;
  readonly priceRange?: PriceRange;
  readonly hours?: BusinessHours[];
  readonly contact?: ContactInfo;
  readonly images?: string[];
}

export interface QuizContent {
  readonly id: string;
  readonly title: string;
  readonly description?: string;
  readonly questions: QuizQuestion[];
  readonly timeLimit?: number;
  readonly passingScore?: number;
  readonly allowRetakes: boolean;
  readonly showAnswersAfter: 'never' | 'submission' | 'completion';
}

export interface SurveyContent {
  readonly id: string;
  readonly title: string;
  readonly description?: string;
  readonly questions: SurveyQuestion[];
  readonly allowAnonymous: boolean;
  readonly showProgress: boolean;
  readonly autoSave: boolean;
}

export interface ProductContent {
  readonly id: string;
  readonly name: string;
  readonly description: string;
  readonly price: number;
  readonly currency: string;
  readonly availability: ProductAvailability;
  readonly images: string[];
  readonly specifications?: Record<string, string>;
  readonly vendor: VendorInfo;
  readonly rating?: ProductRating;
}

export interface EventContent {
  readonly id: string;
  readonly title: string;
  readonly description?: string;
  readonly startTime: Date;
  readonly endTime: Date;
  readonly timezone: string;
  readonly location?: EventLocation;
  readonly isAllDay: boolean;
  readonly organizer: EventOrganizer;
  readonly attendees: EventAttendee[];
  readonly rsvpDeadline?: Date;
  readonly maxAttendees?: number;
}

export interface ReminderContent {
  readonly id: string;
  readonly title: string;
  readonly message: string;
  readonly remindAt: Date;
  readonly recurring?: RecurrenceRule;
  readonly priority: MessagePriority;
  readonly actionRequired: boolean;
  readonly completedAt?: Date;
}

export interface StatusUpdateContent {
  readonly id: string;
  readonly message: string;
  readonly type: StatusType;
  readonly visibility: StatusVisibility;
  readonly expiresAt: Date;
  readonly allowReactions: boolean;
  readonly trackViews: boolean;
  readonly priority: MessagePriority;
}

export interface AnnouncementContent {
  readonly id: string;
  readonly title: string;
  readonly message: string;
  readonly type: AnnouncementType;
  readonly priority: MessagePriority;
  readonly audience: string[];
  readonly expiresAt?: Date;
  readonly acknowledgeRequired: boolean;
  readonly acknowledgedBy?: string[];
}

export interface FormContent {
  readonly id: string;
  readonly title: string;
  readonly description?: string;
  readonly fields: FormField[];
  readonly submitButtonText: string;
  readonly allowMultipleSubmissions: boolean;
  readonly enableSaveProgress: boolean;
}

export interface RichCardContent {
  readonly id: string;
  readonly title: string;
  readonly description?: string;
  readonly imageUrl?: string;
  readonly actions: CardAction[];
  readonly metadata?: Record<string, unknown>;
  readonly template: CardTemplate;
}

export interface EmbedContent {
  readonly url: string;
  readonly title?: string;
  readonly description?: string;
  readonly imageUrl?: string;
  readonly siteName?: string;
  readonly embedType: EmbedType;
  readonly embedHtml?: string;
  readonly metadata?: EmbedMetadata;
}

// Supporting types and enums
export interface MessagePreview {
  readonly messageId: string;
  readonly senderName: string;
  readonly contentPreview: string;
  readonly timestamp: Date;
  readonly messageType: MessageType;
}

export interface Coordinates {
  readonly latitude: number;
  readonly longitude: number;
}

export enum PriceRange {
  BUDGET = 'budget',
  MODERATE = 'moderate',
  EXPENSIVE = 'expensive',
  LUXURY = 'luxury'
}

export interface BusinessHours {
  readonly dayOfWeek: number; // 0-6, Sunday-Saturday
  readonly openTime: string;   // HH:mm format
  readonly closeTime: string;  // HH:mm format
  readonly closed: boolean;
}

export interface ContactInfo {
  readonly phone?: string;
  readonly email?: string;
  readonly website?: string;
  readonly socialMedia?: Record<string, string>;
}

export interface QuizQuestion {
  readonly id: string;
  readonly question: string;
  readonly type: QuestionType;
  readonly options?: string[];
  readonly correctAnswer: string | string[];
  readonly explanation?: string;
  readonly points: number;
}

export interface SurveyQuestion {
  readonly id: string;
  readonly question: string;
  readonly type: QuestionType;
  readonly options?: string[];
  readonly isRequired: boolean;
  readonly branching?: BranchingRule[];
}

export enum QuestionType {
  MULTIPLE_CHOICE = 'multiple_choice',
  MULTIPLE_SELECT = 'multiple_select',
  TRUE_FALSE = 'true_false',
  SHORT_ANSWER = 'short_answer',
  LONG_ANSWER = 'long_answer',
  SCALE = 'scale',
  RATING = 'rating'
}

export interface BranchingRule {
  readonly condition: BranchingCondition;
  readonly targetQuestion: string;
}

export interface BranchingCondition {
  readonly operator: 'equals' | 'not_equals' | 'greater_than' | 'less_than' | 'contains';
  readonly value: unknown;
}

export enum ProductAvailability {
  IN_STOCK = 'in_stock',
  OUT_OF_STOCK = 'out_of_stock',
  PREORDER = 'preorder',
  DISCONTINUED = 'discontinued'
}

export interface VendorInfo {
  readonly name: string;
  readonly verificationLevel: VerificationLevel;
  readonly rating?: number;
  readonly reviewCount?: number;
}

export enum VerificationLevel {
  UNVERIFIED = 'unverified',
  EMAIL_VERIFIED = 'email_verified',
  PHONE_VERIFIED = 'phone_verified',
  BUSINESS_VERIFIED = 'business_verified',
  PREMIUM_VERIFIED = 'premium_verified'
}

export interface ProductRating {
  readonly average: number;
  readonly count: number;
  readonly distribution: number[]; // [1-star, 2-star, 3-star, 4-star, 5-star] counts
}

export interface EventLocation {
  readonly name: string;
  readonly address?: string;
  readonly coordinates?: Coordinates;
  readonly type: LocationType;
  readonly url?: string;
}

export enum LocationType {
  PHYSICAL = 'physical',
  VIRTUAL = 'virtual',
  HYBRID = 'hybrid'
}

export interface EventOrganizer {
  readonly userId: string;
  readonly name: string;
  readonly email?: string;
  readonly avatar?: string;
}

export interface EventAttendee {
  readonly userId: string;
  readonly name: string;
  readonly status: AttendeeStatus;
  readonly rsvpAt: Date;
  readonly notes?: string;
}

export enum AttendeeStatus {
  ATTENDING = 'attending',
  NOT_ATTENDING = 'not_attending',
  MAYBE = 'maybe',
  PENDING = 'pending'
}

export interface RecurrenceRule {
  readonly frequency: RecurrenceFrequency;
  readonly interval: number;
  readonly endDate?: Date;
  readonly occurrences?: number;
}

export enum RecurrenceFrequency {
  DAILY = 'daily',
  WEEKLY = 'weekly',
  MONTHLY = 'monthly',
  YEARLY = 'yearly'
}

export enum StatusType {
  AVAILABILITY = 'availability',
  MOOD = 'mood',
  ACTIVITY = 'activity',
  LOCATION = 'location',
  CUSTOM = 'custom'
}

export enum StatusVisibility {
  PUBLIC = 'public',
  CONTACTS = 'contacts',
  CLOSE_FRIENDS = 'close_friends',
  PRIVATE = 'private'
}

export enum AnnouncementType {
  GENERAL = 'general',
  URGENT = 'urgent',
  MAINTENANCE = 'maintenance',
  FEATURE = 'feature',
  POLICY = 'policy'
}

export interface FormField {
  readonly id: string;
  readonly label: string;
  readonly type: FormFieldType;
  readonly isRequired: boolean;
  readonly placeholder?: string;
  readonly validation?: FieldValidation;
  readonly options?: string[];
  readonly defaultValue?: unknown;
}

export enum FormFieldType {
  TEXT = 'text',
  EMAIL = 'email',
  NUMBER = 'number',
  DATE = 'date',
  TIME = 'time',
  DATETIME = 'datetime',
  SELECT = 'select',
  MULTI_SELECT = 'multi_select',
  RADIO = 'radio',
  CHECKBOX = 'checkbox',
  TEXTAREA = 'textarea',
  FILE = 'file',
  URL = 'url',
  PHONE = 'phone'
}

export interface FieldValidation {
  readonly minLength?: number;
  readonly maxLength?: number;
  readonly pattern?: string;
  readonly min?: number;
  readonly max?: number;
  readonly customValidation?: string;
}

export interface CardAction {
  readonly id: string;
  readonly type: CardActionType;
  readonly label: string;
  readonly style: CardActionStyle;
  readonly payload?: Record<string, unknown>;
  readonly url?: string;
  readonly disabled?: boolean;
}

export enum CardActionType {
  BUTTON = 'button',
  LINK = 'link',
  SUBMIT = 'submit',
  SHARE = 'share',
  DOWNLOAD = 'download'
}

export enum CardActionStyle {
  PRIMARY = 'primary',
  SECONDARY = 'secondary',
  DESTRUCTIVE = 'destructive',
  GHOST = 'ghost',
  OUTLINE = 'outline'
}

export enum CardTemplate {
  BASIC = 'basic',
  MEDIA = 'media',
  PRODUCT = 'product',
  ARTICLE = 'article',
  EVENT = 'event',
  PROFILE = 'profile',
  CUSTOM = 'custom'
}

export enum EmbedType {
  WEBSITE = 'website',
  VIDEO = 'video',
  AUDIO = 'audio',
  IMAGE = 'image',
  DOCUMENT = 'document',
  SOCIAL = 'social',
  RICH = 'rich'
}

export interface EmbedMetadata {
  readonly author?: string;
  readonly publishedAt?: Date;
  readonly duration?: number;
  readonly dimensions?: { width: number; height: number };
  readonly tags?: string[];
}

// Utility types
export type MessageComponentProps<T extends MessageContent = MessageContent> = {
  message: MessageData & { content: T };
  onInteraction?: (interaction: InteractionRequest) => void;
  className?: string;
  children?: ReactNode;
};

export interface InteractionRequest {
  readonly messageId: string;
  readonly interactionType: InteractionType;
  readonly data: Record<string, unknown>;
  readonly userId: string;
  readonly timestamp?: Date;
}

// Type guards
export const isReplyMessage = (content: MessageContent): content is ReplyContent => {
  return typeof content === 'object' && content !== null && 'originalMessageId' in content;
};

export const isQuizMessage = (content: MessageContent): content is QuizContent => {
  return typeof content === 'object' && content !== null && 'questions' in content;
};

export const isEventMessage = (content: MessageContent): content is EventContent => {
  return typeof content === 'object' && content !== null && 'startTime' in content && 'endTime' in content;
};

// Message validation utilities
export const validateMessageData = (message: Partial<MessageData>): ValidationError[] => {
  const errors: ValidationError[] = [];

  if (!message.id) {
    errors.push({
      field: 'id',
      message: 'Message ID is required',
      code: 'REQUIRED_FIELD',
      severity: ValidationSeverity.ERROR
    });
  }

  if (!message.senderId) {
    errors.push({
      field: 'senderId',
      message: 'Sender ID is required',
      code: 'REQUIRED_FIELD',
      severity: ValidationSeverity.ERROR
    });
  }

  if (!message.type || !Object.values(MessageType).includes(message.type)) {
    errors.push({
      field: 'type',
      message: 'Valid message type is required',
      code: 'INVALID_TYPE',
      severity: ValidationSeverity.ERROR
    });
  }

  return errors;
};