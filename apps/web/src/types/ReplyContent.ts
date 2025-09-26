// T028 - Reply message content interfaces
/**
 * Type definitions for reply message content with thread visualization support
 * Supports nested threading, visual connections, and context preservation
 */

import { MessageType } from './MessageData';

// Core Reply Content Interface
export interface ReplyContent {
  readonly originalMessageId: string;
  readonly replyText: string;
  readonly threadDepth: number;
  readonly isThreadStart: boolean;
  readonly originalPreview?: MessagePreview;
  readonly threadMetadata?: ThreadMetadata;
}

// Message Preview for Original Message Context
export interface MessagePreview {
  readonly messageId: string;
  readonly senderName: string;
  readonly senderAvatar?: string;
  readonly contentPreview: string;
  readonly timestamp: Date;
  readonly messageType: MessageType;
  readonly attachmentCount?: number;
  readonly isEdited?: boolean;
}

// Thread Metadata for Visual Rendering
export interface ThreadMetadata {
  readonly threadId: string;
  readonly rootMessageId: string;
  readonly parentMessageId?: string;
  readonly childMessageIds: string[];
  readonly totalThreadMessages: number;
  readonly threadParticipants: ThreadParticipant[];
  readonly threadCreatedAt: Date;
  readonly threadUpdatedAt: Date;
  readonly threadStatus: ThreadStatus;
}

// Thread Participant Information
export interface ThreadParticipant {
  readonly userId: string;
  readonly userName: string;
  readonly userAvatar?: string;
  readonly messageCount: number;
  readonly lastMessageAt: Date;
  readonly isActive: boolean;
}

// Thread Status
export enum ThreadStatus {
  ACTIVE = 'active',
  RESOLVED = 'resolved',
  ARCHIVED = 'archived',
  LOCKED = 'locked'
}

// Thread Navigation Data
export interface ThreadNavigation {
  readonly hasParent: boolean;
  readonly hasChildren: boolean;
  readonly parentMessageId?: string;
  readonly childMessageIds: string[];
  readonly siblingMessageIds: string[];
  readonly threadPosition: number;
  readonly totalThreadSize: number;
}

// Reply Composer State
export interface ReplyComposerState {
  readonly isVisible: boolean;
  readonly replyingTo: MessagePreview;
  readonly draftText: string;
  readonly attachments: ReplyAttachment[];
  readonly mentions: ReplyMention[];
  readonly isSubmitting: boolean;
  readonly validationErrors: string[];
}

// Reply Attachments
export interface ReplyAttachment {
  readonly id: string;
  readonly fileName: string;
  readonly fileType: string;
  readonly fileSize: number;
  readonly uploadUrl?: string;
  readonly thumbnailUrl?: string;
  readonly uploadProgress?: number;
  readonly uploadStatus: UploadStatus;
}

export enum UploadStatus {
  PENDING = 'pending',
  UPLOADING = 'uploading',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled'
}

// Reply Mentions
export interface ReplyMention {
  readonly userId: string;
  readonly userName: string;
  readonly displayName: string;
  readonly startIndex: number;
  readonly length: number;
  readonly type: MentionType;
}

export enum MentionType {
  USER = 'user',
  CHANNEL = 'channel',
  EVERYONE = 'everyone',
  HERE = 'here'
}

// Reply Rendering Configuration
export interface ReplyRenderConfig {
  readonly maxThreadDepth: number;
  readonly showThreadLines: boolean;
  readonly compactMode: boolean;
  readonly showAvatars: boolean;
  readonly showTimestamps: boolean;
  readonly showParticipantCount: boolean;
  readonly enableInlineReplies: boolean;
  readonly collapseOldThreads: boolean;
}

// Reply Interaction Events
export interface ReplyInteractionEvent {
  readonly type: ReplyInteractionType;
  readonly messageId: string;
  readonly userId: string;
  readonly timestamp: Date;
  readonly data?: Record<string, unknown>;
}

export enum ReplyInteractionType {
  THREAD_EXPAND = 'thread_expand',
  THREAD_COLLAPSE = 'thread_collapse',
  THREAD_NAVIGATE = 'thread_navigate',
  REPLY_START = 'reply_start',
  REPLY_CANCEL = 'reply_cancel',
  REPLY_SUBMIT = 'reply_submit',
  THREAD_RESOLVE = 'thread_resolve',
  THREAD_ARCHIVE = 'thread_archive'
}

// Reply Validation Rules
export interface ReplyValidationRules {
  readonly minLength: number;
  readonly maxLength: number;
  readonly maxThreadDepth: number;
  readonly allowEmptyReplies: boolean;
  readonly allowSelfReplies: boolean;
  readonly allowNestedMedia: boolean;
  readonly requireParentAccess: boolean;
}

// Reply Performance Metrics
export interface ReplyPerformanceMetrics {
  readonly threadRenderTime: number;
  readonly messageLoadTime: number;
  readonly interactionLatency: number;
  readonly memoryUsage: number;
  readonly visibleThreadCount: number;
  readonly cachedMessageCount: number;
}

// Threading Algorithm Configuration
export interface ThreadingAlgorithmConfig {
  readonly maxDepthBeforeFlattening: number;
  readonly threadCollapseThreshold: number;
  readonly oldThreadDefinitionHours: number;
  readonly maxVisibleRepliesPerThread: number;
  readonly threadSortOrder: ThreadSortOrder;
  readonly threadGroupingStrategy: ThreadGroupingStrategy;
}

export enum ThreadSortOrder {
  CHRONOLOGICAL = 'chronological',
  REVERSE_CHRONOLOGICAL = 'reverse_chronological',
  MOST_REPLIES = 'most_replies',
  MOST_RECENT_ACTIVITY = 'most_recent_activity'
}

export enum ThreadGroupingStrategy {
  BY_THREAD_ROOT = 'by_thread_root',
  BY_PARTICIPANT = 'by_participant',
  BY_TIME_WINDOW = 'by_time_window',
  FLAT = 'flat'
}

// Reply Search and Filtering
export interface ReplySearchCriteria {
  readonly query?: string;
  readonly authorIds?: string[];
  readonly dateRange?: DateRange;
  readonly threadIds?: string[];
  readonly hasAttachments?: boolean;
  readonly mentions?: string[];
  readonly sortBy: ReplySortBy;
  readonly limit: number;
  readonly offset: number;
}

export enum ReplySortBy {
  RELEVANCE = 'relevance',
  DATE_ASC = 'date_asc',
  DATE_DESC = 'date_desc',
  THREAD_SIZE = 'thread_size',
  AUTHOR_NAME = 'author_name'
}

export interface DateRange {
  readonly startDate: Date;
  readonly endDate: Date;
}

// Reply Analytics Data
export interface ReplyAnalytics {
  readonly threadEngagement: ThreadEngagement;
  readonly participationMetrics: ParticipationMetrics;
  readonly responseTimeMetrics: ResponseTimeMetrics;
  readonly contentAnalytics: ContentAnalytics;
}

export interface ThreadEngagement {
  readonly averageRepliesPerThread: number;
  readonly averageParticipantsPerThread: number;
  readonly threadCompletionRate: number;
  readonly averageThreadDuration: number;
  readonly mostActiveThreads: string[];
}

export interface ParticipationMetrics {
  readonly totalParticipants: number;
  readonly activeParticipants: number;
  readonly averageRepliesPerUser: number;
  readonly topContributors: ThreadParticipant[];
  readonly participationDistribution: Record<string, number>;
}

export interface ResponseTimeMetrics {
  readonly averageResponseTime: number;
  readonly medianResponseTime: number;
  readonly fastestResponseTime: number;
  readonly slowestResponseTime: number;
  readonly responseTimeDistribution: number[];
}

export interface ContentAnalytics {
  readonly averageReplyLength: number;
  readonly mostCommonWords: WordFrequency[];
  readonly mentionFrequency: Record<string, number>;
  readonly attachmentTypes: Record<string, number>;
  readonly sentimentAnalysis?: SentimentScore;
}

export interface WordFrequency {
  readonly word: string;
  readonly frequency: number;
  readonly percentage: number;
}

export interface SentimentScore {
  readonly positive: number;
  readonly neutral: number;
  readonly negative: number;
  readonly overall: number;
}

// Utility Functions and Type Guards
export const isValidReplyContent = (content: unknown): content is ReplyContent => {
  return (
    typeof content === 'object' &&
    content !== null &&
    'originalMessageId' in content &&
    'replyText' in content &&
    'threadDepth' in content &&
    'isThreadStart' in content &&
    typeof (content as any).originalMessageId === 'string' &&
    typeof (content as any).replyText === 'string' &&
    typeof (content as any).threadDepth === 'number' &&
    typeof (content as any).isThreadStart === 'boolean'
  );
};

export const calculateThreadDepth = (
  originalMessageId: string,
  messageHistory: Array<{ id: string; content: ReplyContent | string }>
): number => {
  let depth = 1;
  let currentParent = originalMessageId;

  while (currentParent) {
    const parentMessage = messageHistory.find(msg => msg.id === currentParent);
    if (parentMessage && typeof parentMessage.content === 'object' && 'originalMessageId' in parentMessage.content) {
      currentParent = parentMessage.content.originalMessageId;
      depth++;
    } else {
      break;
    }
  }

  return depth;
};

export const buildThreadTree = (
  messages: Array<{ id: string; content: ReplyContent | string }>
): ThreadNode[] => {
  const messageMap = new Map<string, ThreadNode>();
  const rootNodes: ThreadNode[] = [];

  // Create nodes for all messages
  messages.forEach(message => {
    const node: ThreadNode = {
      id: message.id,
      content: message.content,
      children: [],
      parent: null,
      depth: 0
    };
    messageMap.set(message.id, node);
  });

  // Build parent-child relationships
  messages.forEach(message => {
    if (typeof message.content === 'object' && 'originalMessageId' in message.content) {
      const childNode = messageMap.get(message.id);
      const parentNode = messageMap.get(message.content.originalMessageId);

      if (childNode && parentNode) {
        childNode.parent = parentNode;
        childNode.depth = parentNode.depth + 1;
        parentNode.children.push(childNode);
      }
    } else {
      // Root message
      const rootNode = messageMap.get(message.id);
      if (rootNode) {
        rootNodes.push(rootNode);
      }
    }
  });

  return rootNodes;
};

export interface ThreadNode {
  readonly id: string;
  readonly content: ReplyContent | string;
  readonly children: ThreadNode[];
  readonly parent: ThreadNode | null;
  readonly depth: number;
}

// Export all types for external use
export type {
  MessagePreview,
  ThreadMetadata,
  ThreadParticipant,
  ThreadNavigation,
  ReplyComposerState,
  ReplyAttachment,
  ReplyMention,
  ReplyRenderConfig,
  ReplyInteractionEvent,
  ReplyValidationRules,
  ReplyPerformanceMetrics,
  ThreadingAlgorithmConfig,
  ReplySearchCriteria,
  ReplyAnalytics
};