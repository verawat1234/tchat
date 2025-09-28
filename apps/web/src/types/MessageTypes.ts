/**
 * Message Types - Web Platform
 *
 * Comprehensive message type system aligned with mobile Kotlin implementation
 * Cross-platform consistency with enhanced message structure and functionality
 */

/**
 * Message Type Enum - Aligned with mobile CommunicationEnums.kt
 * 22 comprehensive message types supporting rich communication features
 */
export enum MessageType {
  // Core message types
  TEXT = 'text',
  IMAGE = 'image',
  VIDEO = 'video',
  AUDIO = 'audio',
  FILE = 'file',
  LOCATION = 'location',
  CONTACT = 'contact',
  STICKER = 'sticker',
  GIF = 'gif',
  POLL = 'poll',
  EVENT = 'event',
  SYSTEM = 'system',
  DELETED = 'deleted',

  // Extended message types (Phase D: T036-T041)
  EMBED = 'embed',
  EVENT_MESSAGE = 'event_message',
  FORM = 'form',
  LOCATION_MESSAGE = 'location_message',
  PAYMENT = 'payment',
  FILE_MESSAGE = 'file_message',

  // Commerce types
  PRODUCT = 'product',
  INVOICE = 'invoice',
  ORDER = 'order',
  ORDER_STATUS_UPDATE = 'order_status_update',
  PAYMENT_REQUEST = 'payment_request',
  QUOTATION = 'quotation'
}

/**
 * Delivery Status Enum - Aligned with mobile DeliveryStatus
 */
export enum DeliveryStatus {
  PENDING = 'pending',
  SENT = 'sent',
  DELIVERED = 'delivered',
  READ = 'read',
  FAILED = 'failed'
}

/**
 * Attachment Type Enum - Aligned with mobile AttachmentType
 */
export enum AttachmentType {
  IMAGE = 'IMAGE',
  VIDEO = 'VIDEO',
  AUDIO = 'AUDIO',
  FILE = 'FILE',
  LOCATION = 'LOCATION'
}

/**
 * Message Attachment Interface - Aligned with mobile MessageAttachment
 */
export interface MessageAttachment {
  id: string;
  type: AttachmentType;
  url: string;
  thumbnail?: string;
  filename?: string;
  fileSize?: number;
  mimeType?: string;
  width?: number;
  height?: number;
  duration?: number; // in milliseconds for audio/video
  caption?: string;
  metadata: Record<string, string>;
}

/**
 * Message Reaction Interface - Aligned with mobile MessageReaction
 */
export interface MessageReaction {
  emoji: string;
  userId: string;
  userName: string;
  timestamp: string;
}

/**
 * Message Reply Interface - Aligned with mobile MessageReply
 */
export interface MessageReply {
  messageId: string;
  senderId: string;
  senderName: string;
  content: string;
  type: MessageType;
  timestamp: string;
}

/**
 * Enhanced Message Interface - Aligned with mobile Message data class
 * Comprehensive message structure supporting rich content and real-time features
 */
export interface Message {
  id: string;
  chatId: string;
  senderId: string;
  senderName: string;
  senderAvatar?: string;
  type: MessageType;
  content: string;
  isEdited: boolean;
  isPinned: boolean;
  isDeleted: boolean;
  replyToId?: string;
  reactions: MessageReaction[];
  attachments: MessageAttachment[];
  createdAt: string;
  editedAt?: string;
  deletedAt?: string;
  deliveryStatus: DeliveryStatus;
  readBy: string[]; // User IDs who have read this message
}

/**
 * Legacy Message Interface for backward compatibility
 */
export interface LegacyMessage {
  id: string;
  senderId: string;
  senderName: string;
  content: string;
  timestamp: string;
  type: 'text' | 'image' | 'video' | 'audio' | 'file' | 'location' | 'payment' | 'system';
  isOwn: boolean;
  fileUrl?: string;
  fileName?: string;
  fileSize?: string;
  metadata?: any;
}

/**
 * MessageData interface used in chat components
 */
export interface MessageData {
  id: string;
  senderId: string;
  senderName: string;
  content: string;
  timestamp: string;
  type: MessageType;
  isOwn: boolean;
  metadata?: any;
  attachments?: MessageAttachment[];
  reactions?: MessageReaction[];
  replyTo?: MessageReply;
  deliveryStatus?: DeliveryStatus;
}

/**
 * Message Type Utilities - Similar to mobile platform functions
 */
export class MessageTypeUtils {

  /**
   * Get media message types aligned with mobile getMediaTypes()
   */
  static getMediaTypes(): MessageType[] {
    return [
      MessageType.IMAGE,
      MessageType.VIDEO,
      MessageType.AUDIO,
      MessageType.FILE,
      MessageType.GIF,
      MessageType.FILE_MESSAGE
    ];
  }

  /**
   * Get interactive message types aligned with mobile getInteractiveTypes()
   */
  static getInteractiveTypes(): MessageType[] {
    return [
      MessageType.POLL,
      MessageType.FORM,
      MessageType.EVENT,
      MessageType.EVENT_MESSAGE,
      MessageType.EMBED
    ];
  }

  /**
   * Get commerce message types
   */
  static getCommerceTypes(): MessageType[] {
    return [
      MessageType.PRODUCT,
      MessageType.INVOICE,
      MessageType.ORDER,
      MessageType.ORDER_STATUS_UPDATE,
      MessageType.PAYMENT_REQUEST,
      MessageType.QUOTATION,
      MessageType.PAYMENT
    ];
  }

  /**
   * Check if message type requires content
   */
  static requiresContent(type: MessageType): boolean {
    return type !== MessageType.SYSTEM && type !== MessageType.DELETED;
  }

  /**
   * Get display name for message type
   */
  static getDisplayName(type: MessageType): string {
    const displayNames: Record<MessageType, string> = {
      [MessageType.TEXT]: 'Text Message',
      [MessageType.IMAGE]: 'Image',
      [MessageType.VIDEO]: 'Video',
      [MessageType.AUDIO]: 'Audio',
      [MessageType.FILE]: 'File',
      [MessageType.LOCATION]: 'Location',
      [MessageType.CONTACT]: 'Contact',
      [MessageType.STICKER]: 'Sticker',
      [MessageType.GIF]: 'GIF',
      [MessageType.POLL]: 'Poll',
      [MessageType.EVENT]: 'Event',
      [MessageType.SYSTEM]: 'System Message',
      [MessageType.DELETED]: 'Deleted Message',
      [MessageType.EMBED]: 'Rich Embed',
      [MessageType.EVENT_MESSAGE]: 'Calendar Event',
      [MessageType.FORM]: 'Interactive Form',
      [MessageType.LOCATION_MESSAGE]: 'Enhanced Location',
      [MessageType.PAYMENT]: 'Payment',
      [MessageType.FILE_MESSAGE]: 'Advanced File',
      [MessageType.PRODUCT]: 'Product',
      [MessageType.INVOICE]: 'Invoice',
      [MessageType.ORDER]: 'Order',
      [MessageType.ORDER_STATUS_UPDATE]: 'Order Status Update',
      [MessageType.PAYMENT_REQUEST]: 'Payment Request',
      [MessageType.QUOTATION]: 'Quotation'
    };

    return displayNames[type] || 'Unknown Message';
  }

  /**
   * Check if message type is media
   */
  static isMediaMessage(type: MessageType): boolean {
    return this.getMediaTypes().includes(type);
  }

  /**
   * Check if message type is interactive
   */
  static isInteractiveMessage(type: MessageType): boolean {
    return this.getInteractiveTypes().includes(type);
  }

  /**
   * Check if message type is commerce-related
   */
  static isCommerceMessage(type: MessageType): boolean {
    return this.getCommerceTypes().includes(type);
  }

  /**
   * From value string to MessageType enum
   */
  static fromValue(value: string): MessageType | null {
    const values = Object.values(MessageType);
    return values.find(type => type === value) || null;
  }
}

/**
 * Message Utility Functions - Aligned with mobile extension functions
 */
export class MessageUtils {

  /**
   * Check if message is from current user
   */
  static isFromCurrentUser(message: Message, currentUserId: string): boolean {
    return message.senderId === currentUserId;
  }

  /**
   * Check if message has attachments
   */
  static hasAttachments(message: Message): boolean {
    return message.attachments.length > 0;
  }

  /**
   * Check if message has reactions
   */
  static hasReactions(message: Message): boolean {
    return message.reactions.length > 0;
  }

  /**
   * Get reaction count
   */
  static getReactionCount(message: Message): number {
    return message.reactions.length;
  }

  /**
   * Check if user has reacted to message
   */
  static hasUserReacted(message: Message, userId: string): boolean {
    return message.reactions.some(reaction => reaction.userId === userId);
  }

  /**
   * Get user's reaction to message
   */
  static getUserReaction(message: Message, userId: string): MessageReaction | null {
    return message.reactions.find(reaction => reaction.userId === userId) || null;
  }

  /**
   * Get display content for message
   */
  static getDisplayContent(message: Message): string {
    if (message.isDeleted) return "This message was deleted";

    switch (message.type) {
      case MessageType.IMAGE:
        return message.attachments.length > 0
          ? (message.attachments[0].caption || "📷 Image")
          : "📷 Image";
      case MessageType.VIDEO:
        return message.attachments.length > 0
          ? (message.attachments[0].caption || "🎥 Video")
          : "🎥 Video";
      case MessageType.AUDIO:
        return "🎵 Audio";
      case MessageType.FILE:
      case MessageType.FILE_MESSAGE:
        return message.attachments.length > 0
          ? `📄 ${message.attachments[0].filename || "File"}`
          : "📄 File";
      case MessageType.LOCATION:
      case MessageType.LOCATION_MESSAGE:
        return "📍 Location";
      case MessageType.STICKER:
        return "Sticker";
      case MessageType.GIF:
        return "🎭 GIF";
      case MessageType.POLL:
        return "📊 Poll";
      case MessageType.EVENT:
      case MessageType.EVENT_MESSAGE:
        return "📅 Event";
      case MessageType.FORM:
        return "📝 Form";
      case MessageType.EMBED:
        return "🔗 Rich Content";
      case MessageType.PRODUCT:
        return "🛍️ Product";
      case MessageType.INVOICE:
        return "🧾 Invoice";
      case MessageType.ORDER:
        return "📦 Order";
      case MessageType.PAYMENT:
      case MessageType.PAYMENT_REQUEST:
        return "💳 Payment";
      case MessageType.QUOTATION:
        return "💰 Quotation";
      case MessageType.SYSTEM:
        return message.content;
      default:
        return message.content;
    }
  }

  /**
   * Check if message is a media message
   */
  static isMediaMessage(message: Message): boolean {
    return MessageTypeUtils.isMediaMessage(message.type);
  }

  /**
   * Check if message can be edited
   */
  static canBeEdited(message: Message, currentUserId: string): boolean {
    return message.senderId === currentUserId &&
           !message.isDeleted &&
           message.type === MessageType.TEXT;
  }

  /**
   * Check if message can be deleted
   */
  static canBeDeleted(message: Message, currentUserId: string): boolean {
    return message.senderId === currentUserId && !message.isDeleted;
  }

  /**
   * Check if message can be replied to
   */
  static canBeRepliedTo(message: Message): boolean {
    return !message.isDeleted && message.type !== MessageType.SYSTEM;
  }
}

/**
 * Message List Utilities - Aligned with mobile extension functions
 */
export class MessageListUtils {

  /**
   * Group messages by date
   */
  static groupByDate(messages: Message[]): Record<string, Message[]> {
    return messages.reduce((groups, message) => {
      // Extract date from timestamp (assuming ISO format)
      const date = message.createdAt.split('T')[0] || message.createdAt;
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(message);
      return groups;
    }, {} as Record<string, Message[]>);
  }

  /**
   * Group consecutive messages from same sender
   */
  static groupConsecutiveMessages(messages: Message[]): Message[][] {
    if (messages.length === 0) return [];

    const groups: Message[][] = [];
    let currentGroup: Message[] = [messages[0]];

    for (let i = 1; i < messages.length; i++) {
      const current = messages[i];
      const previous = messages[i - 1];

      // Group if same sender and within reasonable time
      if (current.senderId === previous.senderId &&
          this.shouldGroupMessages(previous, current)) {
        currentGroup.push(current);
      } else {
        groups.push(currentGroup);
        currentGroup = [current];
      }
    }

    groups.push(currentGroup);
    return groups;
  }

  /**
   * Check if messages should be grouped together
   */
  private static shouldGroupMessages(previous: Message, current: Message): boolean {
    // Simple time-based grouping (in a real app, you'd parse timestamps properly)
    // For now, always group consecutive messages from same sender
    return true;
  }

  /**
   * Search messages by content
   */
  static searchByContent(messages: Message[], query: string): Message[] {
    if (!query.trim()) return messages;

    const lowerQuery = query.toLowerCase();
    return messages.filter(message =>
      MessageUtils.getDisplayContent(message).toLowerCase().includes(lowerQuery) ||
      message.senderName.toLowerCase().includes(lowerQuery)
    );
  }

  /**
   * Filter messages by type
   */
  static filterByType(messages: Message[], type: MessageType): Message[] {
    return messages.filter(message => message.type === type);
  }

  /**
   * Filter media messages
   */
  static filterMediaMessages(messages: Message[]): Message[] {
    return messages.filter(message => MessageUtils.isMediaMessage(message));
  }

  /**
   * Filter unread messages
   */
  static filterUnread(messages: Message[], lastReadMessageId?: string): Message[] {
    if (!lastReadMessageId) return messages;

    const lastReadIndex = messages.findIndex(message => message.id === lastReadMessageId);
    return lastReadIndex >= 0 ? messages.slice(lastReadIndex + 1) : messages;
  }
}

/**
 * Conversion utilities between legacy and new message formats
 */
export class MessageConverter {

  /**
   * Convert legacy message to new format
   */
  static convertLegacyToMessage(legacyMessage: LegacyMessage): Message {
    const messageType = MessageTypeUtils.fromValue(legacyMessage.type) || MessageType.TEXT;

    return {
      id: legacyMessage.id,
      chatId: '', // Not available in legacy format
      senderId: legacyMessage.senderId,
      senderName: legacyMessage.senderName,
      senderAvatar: undefined,
      type: messageType,
      content: legacyMessage.content,
      isEdited: false,
      isPinned: false,
      isDeleted: false,
      replyToId: undefined,
      reactions: [],
      attachments: legacyMessage.fileUrl ? [{
        id: `${legacyMessage.id}_attachment`,
        type: this.getAttachmentTypeFromMessageType(messageType),
        url: legacyMessage.fileUrl,
        filename: legacyMessage.fileName,
        fileSize: legacyMessage.fileSize ? parseInt(legacyMessage.fileSize) : undefined,
        caption: undefined,
        metadata: {}
      }] : [],
      createdAt: legacyMessage.timestamp,
      editedAt: undefined,
      deletedAt: undefined,
      deliveryStatus: DeliveryStatus.DELIVERED,
      readBy: []
    };
  }

  /**
   * Convert new message to legacy format
   */
  static convertMessageToLegacy(message: Message): LegacyMessage {
    const legacyType = this.getLegacyTypeFromMessageType(message.type);
    const firstAttachment = message.attachments[0];

    return {
      id: message.id,
      senderId: message.senderId,
      senderName: message.senderName,
      content: message.content,
      timestamp: message.createdAt,
      type: legacyType,
      isOwn: false, // This needs to be set based on current user context
      fileUrl: firstAttachment?.url,
      fileName: firstAttachment?.filename,
      fileSize: firstAttachment?.fileSize?.toString(),
      metadata: {
        reactions: message.reactions,
        isEdited: message.isEdited,
        isPinned: message.isPinned,
        deliveryStatus: message.deliveryStatus
      }
    };
  }

  /**
   * Convert MessageData to Message
   */
  static convertMessageDataToMessage(messageData: MessageData): Message {
    return {
      id: messageData.id,
      chatId: '', // Not available in MessageData
      senderId: messageData.senderId,
      senderName: messageData.senderName,
      senderAvatar: undefined,
      type: messageData.type,
      content: messageData.content,
      isEdited: false,
      isPinned: false,
      isDeleted: false,
      replyToId: messageData.replyTo?.messageId,
      reactions: messageData.reactions || [],
      attachments: messageData.attachments || [],
      createdAt: messageData.timestamp,
      editedAt: undefined,
      deletedAt: undefined,
      deliveryStatus: messageData.deliveryStatus || DeliveryStatus.DELIVERED,
      readBy: []
    };
  }

  /**
   * Get attachment type from message type
   */
  private static getAttachmentTypeFromMessageType(messageType: MessageType): AttachmentType {
    switch (messageType) {
      case MessageType.IMAGE: return AttachmentType.IMAGE;
      case MessageType.VIDEO: return AttachmentType.VIDEO;
      case MessageType.AUDIO: return AttachmentType.AUDIO;
      case MessageType.LOCATION:
      case MessageType.LOCATION_MESSAGE: return AttachmentType.LOCATION;
      default: return AttachmentType.FILE;
    }
  }

  /**
   * Get legacy type from message type
   */
  private static getLegacyTypeFromMessageType(messageType: MessageType): LegacyMessage['type'] {
    switch (messageType) {
      case MessageType.IMAGE: return 'image';
      case MessageType.VIDEO: return 'video';
      case MessageType.AUDIO: return 'audio';
      case MessageType.FILE:
      case MessageType.FILE_MESSAGE: return 'file';
      case MessageType.LOCATION:
      case MessageType.LOCATION_MESSAGE: return 'location';
      case MessageType.PAYMENT:
      case MessageType.PAYMENT_REQUEST: return 'payment';
      case MessageType.SYSTEM: return 'system';
      default: return 'text';
    }
  }
}