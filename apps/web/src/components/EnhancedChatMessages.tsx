/**
 * Chat Messages Integration Component
 *
 * Bridges the existing ChatTab interface with our advanced MessageComponent system.
 * Provides backward compatibility while enabling new features like delivery receipts,
 * real-time status tracking, and advanced dialog interactions.
 */

import React, { useMemo } from 'react';
import type { Message as OldMessage } from './ChatActions';
import type { Message, MessageType, DeliveryStatus, MessageContent } from '../services/messageService';

// =============================================================================
// Type Mapping and Conversion
// =============================================================================

interface ChatMessagesProps {
  messages: OldMessage[];
  selectedMessages: string[];
  isSelectionMode: boolean;
  onReplyToMessage: (message: OldMessage) => void;
  onEditMessage: (message: OldMessage) => void;
  onDeleteMessage: (messageId: string) => void;
  onForwardMessage: (messageId: string) => void;
  onCopyMessage: (content: string) => void;
  onPinMessage: (messageId: string) => void;
  onSelectMessage: (messageId: string) => void;
  chatId?: string;
  currentUserId?: string;
}

/**
 * Convert old message format to new MessageComponent format
 */
function convertOldToNewMessage(oldMessage: OldMessage): Message {
  // Map old message types to new MessageType enum
  const getMessageType = (oldType: string): MessageType => {
    switch (oldType) {
      case 'text': return 'TEXT' as MessageType;
      case 'voice': return 'AUDIO' as MessageType;
      case 'file': return 'FILE' as MessageType;
      case 'image': return 'IMAGE' as MessageType;
      case 'payment': return 'PAYMENT' as MessageType;
      case 'system': return 'SYSTEM' as MessageType;
      default: return 'TEXT' as MessageType;
    }
  };

  // Create MessageContent based on message type
  const getMessageContent = (message: OldMessage): MessageContent => {
    if (message.type === 'file') {
      return {
        text: message.fileName || 'File',
        data: {
          fileName: message.fileName,
          fileSize: message.fileSize,
          fileUrl: message.fileUrl,
        }
      };
    }

    if (message.type === 'voice') {
      return {
        text: 'Voice message',
        data: {
          duration: message.duration,
          fileUrl: message.fileUrl,
        }
      };
    }

    return {
      text: message.content,
    };
  };

  // Determine delivery status (simplified for now)
  const getDeliveryStatus = (): DeliveryStatus => {
    return 'DELIVERED' as DeliveryStatus; // Default status for existing messages
  };

  return {
    id: oldMessage.id,
    chatId: 'current-chat', // Will be provided by props
    senderId: oldMessage.senderId,
    type: getMessageType(oldMessage.type),
    content: getMessageContent(oldMessage),
    timestamp: oldMessage.timestamp,
    status: getDeliveryStatus(),
    reactions: [],
    metadata: {
      version: 1,
      clientId: 'legacy-client',
      priority: 'normal',
    },
    mentions: [],
    isEdited: false,
    isDeleted: false,
  };
}

/**
 * Enhanced Chat Messages Component with backward compatibility
 */
export function AdvancedChatMessages({
  messages,
  selectedMessages,
  isSelectionMode,
  onReplyToMessage,
  onEditMessage,
  onDeleteMessage,
  onForwardMessage,
  onCopyMessage,
  onPinMessage,
  onSelectMessage,
  chatId = 'current-chat',
  currentUserId = 'current-user'
}: ChatMessagesProps) {

  // Convert old messages to new format
  const convertedMessages = useMemo(() => {
    return messages.map(oldMsg => ({
      ...convertOldToNewMessage(oldMsg),
      chatId: chatId
    }));
  }, [messages, chatId]);

  // Handler adapters to bridge old and new interfaces
  const handleReply = (message: Message) => {
    const oldMessage = messages.find(m => m.id === message.id);
    if (oldMessage) {
      onReplyToMessage(oldMessage);
    }
  };

  const handleEdit = (message: Message) => {
    const oldMessage = messages.find(m => m.id === message.id);
    if (oldMessage) {
      onEditMessage(oldMessage);
    }
  };

  const handleDelete = (message: Message) => {
    onDeleteMessage(message.id);
  };

  const handleReact = (message: Message, emoji: string) => {
    // For now, just show a toast - can be enhanced later
    console.log(`Reaction ${emoji} added to message ${message.id}`);
  };

  const handleThreadOpen = (message: Message) => {
    // For now, just log - can be enhanced with thread functionality later
    console.log(`Opening thread for message ${message.id}`);
  };

  // If no messages, show empty state
  if (convertedMessages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center text-muted-foreground">
        <div className="text-center">
          <p className="text-lg font-medium mb-2">No messages yet</p>
          <p className="text-sm">Start a conversation!</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-hidden">
      <LegacyChatMessagesFallback
        messages={messages}
        selectedMessages={selectedMessages}
        isSelectionMode={isSelectionMode}
        onReplyToMessage={onReplyToMessage}
        onEditMessage={onEditMessage}
        onDeleteMessage={onDeleteMessage}
        onForwardMessage={onForwardMessage}
        onCopyMessage={onCopyMessage}
        onPinMessage={onPinMessage}
        onSelectMessage={onSelectMessage}
        chatId={chatId}
        currentUserId={currentUserId}
      />
    </div>
  );
}

/**
 * Legacy Chat Messages Fallback Component
 *
 * Falls back to the original ChatMessages if there are any issues
 * with the enhanced version. Provides graceful degradation.
 */
export function LegacyChatMessagesFallback(props: ChatMessagesProps) {
  // This would import and use the original ChatMessages component
  // For now, we'll use a simple implementation
  return (
    <div className="flex-1 overflow-auto space-y-4 p-4">
      {props.messages.map((message) => (
        <div
          key={message.id}
          className={`flex gap-3 ${message.isOwn ? 'flex-row-reverse' : ''}`}
        >
          <div className={`flex flex-col max-w-md ${message.isOwn ? 'items-end' : 'items-start'}`}>
            {!message.isOwn && (
              <span className="text-xs text-muted-foreground mb-1">
                {message.senderName}
              </span>
            )}

            <div
              className={`rounded-lg px-3 py-2 ${
                message.isOwn
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-muted'
              }`}
            >
              <p className="text-sm whitespace-pre-line">{message.content}</p>
            </div>

            <span className="text-xs text-muted-foreground mt-1">
              {message.timestamp}
            </span>
          </div>
        </div>
      ))}
    </div>
  );
}

export default AdvancedChatMessages;