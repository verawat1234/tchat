/**
 * Chat Message Component Contracts
 */

import { BaseComponent, ClickHandler, AccessibilityProps } from './base';

/**
 * User information
 */
export interface User {
  /** Unique user ID */
  id: string;
  /** Display name */
  name: string;
  /** Avatar URL */
  avatar?: string;
  /** Online status */
  status?: 'online' | 'offline' | 'away' | 'busy';
}

/**
 * Media content structure
 */
export interface MediaContent {
  /** Media URL */
  url: string;
  /** Media type */
  type: 'image' | 'video' | 'audio' | 'file';
  /** Original filename */
  filename?: string;
  /** File size in bytes */
  size?: number;
  /** Thumbnail URL */
  thumbnail?: string;
  /** Media dimensions */
  dimensions?: {
    width: number;
    height: number;
  };
  /** Duration for audio/video */
  duration?: number;
}

/**
 * System message content
 */
export interface SystemContent {
  /** System message type */
  type: 'join' | 'leave' | 'rename' | 'notification' | 'error' | 'info';
  /** Message text */
  message: string;
  /** Additional metadata */
  metadata?: Record<string, any>;
}

/**
 * Message reaction
 */
export interface Reaction {
  /** Emoji character */
  emoji: string;
  /** Number of users who reacted */
  count: number;
  /** User IDs who reacted */
  users: string[];
  /** Whether current user has reacted */
  hasUserReacted: boolean;
}

/**
 * Message status
 */
export type MessageStatus = 'sending' | 'sent' | 'delivered' | 'read' | 'failed';

/**
 * Message content types
 */
export type MessageContent = string | MediaContent | SystemContent;

/**
 * Chat message properties
 */
export interface ChatMessageProps extends BaseComponent, AccessibilityProps {
  /** Message type */
  type: 'text' | 'image' | 'video' | 'audio' | 'file' | 'system' | 'typing';
  /** Message content */
  content: MessageContent;
  /** Message timestamp */
  timestamp: Date;
  /** Message sender */
  sender?: User;
  /** Whether message is from current user */
  isOwn?: boolean;
  /** Delivery/read status */
  status?: MessageStatus;
  /** Message reactions */
  reactions?: Reaction[];
  /** Replied-to message */
  reply?: Partial<ChatMessageProps>;
  /** Whether message is being edited */
  editing?: boolean;
  /** Whether message is selected */
  selected?: boolean;
  /** Whether to show timestamp */
  showTimestamp?: boolean;
  /** Whether to show avatar */
  showAvatar?: boolean;
  /** Whether to show sender name */
  showSender?: boolean;
  /** Click handler for message */
  onClick?: ClickHandler;
  /** Reaction handler */
  onReaction?: (emoji: string) => void;
  /** Reply handler */
  onReply?: () => void;
  /** Edit handler */
  onEdit?: () => void;
  /** Delete handler */
  onDelete?: () => void;
}

/**
 * Typing indicator properties
 */
export interface TypingIndicatorProps extends BaseComponent {
  /** Users who are typing */
  users: User[];
  /** Maximum users to show */
  maxUsers?: number;
}

/**
 * Message group properties (for grouping consecutive messages)
 */
export interface MessageGroupProps extends BaseComponent {
  /** Messages in the group */
  messages: ChatMessageProps[];
  /** Group sender */
  sender: User;
  /** Whether group is from current user */
  isOwn?: boolean;
  /** Whether to show avatar for group */
  showAvatar?: boolean;
  /** Whether to show timestamp for group */
  showTimestamp?: boolean;
}

/**
 * Message bubble properties
 */
export interface MessageBubbleProps extends BaseComponent {
  /** Bubble variant */
  variant?: 'default' | 'own' | 'system';
  /** Whether bubble is part of a group */
  grouped?: boolean;
  /** Bubble position in group */
  position?: 'first' | 'middle' | 'last' | 'single';
  /** Whether bubble has tail */
  tail?: boolean;
}