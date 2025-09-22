/**
 * Chat Message Components
 * Chat message display with reactions, status, and typing indicators
 */

import React from 'react';
import { cn } from '@/utils/cn';
import type {
  ChatMessageProps,
  TypingIndicatorProps,
  MessageGroupProps,
  MessageBubbleProps,
  MediaContent,
  SystemContent,
  MessageStatus,
  Reaction,
  User
} from '../../../../specs/001-agent-frontend-specialist/contracts/chat-message';

/**
 * Chat message component for displaying messages
 */
export const ChatMessage = React.forwardRef<HTMLDivElement, ChatMessageProps>(
  ({
    className,
    testId,
    type,
    content,
    timestamp,
    sender,
    isOwn = false,
    status,
    reactions,
    reply,
    editing = false,
    selected = false,
    showTimestamp = true,
    showAvatar = true,
    showSender = false,
    onClick,
    onReaction,
    onReply,
    onEdit,
    onDelete,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedBy,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const getMessageContent = () => {
      if (type === 'system') {
        const systemContent = content as SystemContent;
        return (
          <div className="text-center text-sm text-gray-500">
            {systemContent.message}
          </div>
        );
      }

      if (type === 'typing') {
        return <TypingIndicator users={sender ? [sender] : []} />;
      }

      if (type === 'text') {
        return <span className="message-text">{content as string}</span>;
      }

      if (['image', 'video', 'audio', 'file'].includes(type)) {
        const mediaContent = content as MediaContent;

        if (type === 'image') {
          return (
            <div className="message-media message-image">
              <img
                src={mediaContent.url}
                alt={mediaContent.filename || 'Image'}
                className="max-w-full rounded"
              />
            </div>
          );
        }

        if (type === 'video') {
          return (
            <div className="message-media message-video">
              <video
                src={mediaContent.url}
                controls
                className="max-w-full rounded"
              >
                Your browser does not support the video tag.
              </video>
            </div>
          );
        }

        if (type === 'audio') {
          return (
            <div className="message-media message-audio">
              <audio
                src={mediaContent.url}
                controls
                className="w-full"
              >
                Your browser does not support the audio tag.
              </audio>
            </div>
          );
        }

        if (type === 'file') {
          return (
            <div className="message-media message-file flex items-center space-x-2 p-2 bg-gray-100 rounded">
              <svg className="w-6 h-6 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clipRule="evenodd" />
              </svg>
              <div className="flex-1">
                <div className="text-sm font-medium">{mediaContent.filename || 'File'}</div>
                {mediaContent.size && (
                  <div className="text-xs text-gray-500">
                    {formatFileSize(mediaContent.size)}
                  </div>
                )}
              </div>
            </div>
          );
        }
      }

      return null;
    };

    const formatFileSize = (bytes: number): string => {
      if (bytes === 0) return '0 Bytes';
      const k = 1024;
      const sizes = ['Bytes', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
    };

    const formatTimestamp = (date: Date): string => {
      const hours = date.getHours().toString().padStart(2, '0');
      const minutes = date.getMinutes().toString().padStart(2, '0');
      return `${hours}:${minutes}`;
    };

    const getStatusIcon = (status: MessageStatus) => {
      switch (status) {
        case 'sending':
          return <span className="text-gray-400">⏱</span>;
        case 'sent':
          return <span className="text-gray-400">✓</span>;
        case 'delivered':
          return <span className="text-gray-400">✓✓</span>;
        case 'read':
          return <span className="text-blue-500">✓✓</span>;
        case 'failed':
          return <span className="text-red-500">✗</span>;
        default:
          return null;
      }
    };

    const handleClick = () => {
      if (onClick) {
        onClick();
      }
    };

    const handleKeyDown = (event: React.KeyboardEvent) => {
      if (onClick && (event.key === 'Enter' || event.key === ' ')) {
        event.preventDefault();
        onClick();
      }
    };

    const handleReaction = (emoji: string) => {
      if (onReaction) {
        onReaction(emoji);
      }
    };

    // System messages
    if (type === 'system') {
      return (
        <div
          ref={ref}
          data-testid={testId}
          className={cn(
            'py-2 px-4 text-center',
            'chat-message',
            'chat-message-system',
            'message-system',
            className
          )}
          {...props}
        >
          {getMessageContent()}
        </div>
      );
    }

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex',
          isOwn ? 'justify-end' : 'justify-start',
          'mb-2',

          // Message type
          `message-${type}`,

          // State classes
          selected && 'bg-blue-50 message-selected chat-message-selected',
          editing && 'opacity-75 message-editing chat-message-editing',
          isOwn && 'chat-message-own',

          // Custom classes
          'chat-message',
          className
        )}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
        aria-label={ariaLabel}
        aria-describedby={ariaDescribedBy}
        aria-expanded={ariaExpanded}
        aria-disabled={ariaDisabled}
        role={role || 'article'}
        tabIndex={onClick ? (tabIndex ?? 0) : tabIndex}
        data-status={status}
        {...props}
      >
        <div className={cn(
          'flex',
          isOwn ? 'flex-row-reverse' : 'flex-row',
          'items-end space-x-2'
        )}>
          {/* Avatar */}
          {showAvatar && sender && !isOwn && (
            <div className="flex-shrink-0 message-avatar">
              {sender.avatar ? (
                <img
                  src={sender.avatar}
                  alt={sender.name}
                  className="w-8 h-8 rounded-full"
                />
              ) : (
                <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center">
                  <span className="text-sm font-medium text-gray-600">
                    {sender.name.charAt(0).toUpperCase()}
                  </span>
                </div>
              )}
            </div>
          )}

          <div className={cn(
            'flex flex-col',
            isOwn && 'items-end'
          )}>
            {/* Sender name */}
            {showSender && sender && !isOwn && (
              <div className="text-xs text-gray-500 mb-1 message-sender">
                {sender.name}
              </div>
            )}

            {/* Reply */}
            {reply && (
              <div className="message-reply mb-1 px-2 py-1 bg-gray-100 rounded text-sm text-gray-600 border-l-2 border-gray-400">
                {typeof reply.content === 'string' ? reply.content : 'Media'}
              </div>
            )}

            {/* Message bubble */}
            <div
              className={cn(
                'px-3 py-2 rounded-lg max-w-xs lg:max-w-md',
                isOwn
                  ? 'bg-blue-500 text-white message-own chat-message-own'
                  : 'bg-gray-200 text-gray-900 message-other',
                type === 'typing' && 'message-typing'
              )}
            >
              {getMessageContent()}
            </div>

            {/* Reactions */}
            {reactions && reactions.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-1 message-reactions">
                {reactions.map((reaction, index) => (
                  <button
                    key={`${reaction.emoji}-${index}`}
                    className={cn(
                      'px-1.5 py-0.5 rounded-full text-xs flex items-center space-x-1',
                      reaction.hasUserReacted
                        ? 'bg-blue-100 text-blue-700 reaction-active'
                        : 'bg-gray-100 text-gray-700 reaction-inactive'
                    )}
                    onClick={(e) => {
                      e.stopPropagation();
                      handleReaction(reaction.emoji);
                    }}
                  >
                    <span>{reaction.emoji}</span>
                    <span>{reaction.count}</span>
                  </button>
                ))}
              </div>
            )}

            {/* Timestamp and status */}
            {(showTimestamp || status) && (
              <div className="flex items-center space-x-1 mt-1 text-xs text-gray-500">
                {showTimestamp && (
                  <span className="message-timestamp">
                    {formatTimestamp(timestamp)}
                  </span>
                )}
                {status && isOwn && (
                  <span className="message-status">
                    {getStatusIcon(status)}
                  </span>
                )}
              </div>
            )}

            {/* Action buttons */}
            {(onReply || onEdit || onDelete) && (
              <div className="flex space-x-2 mt-1 message-actions">
                {onReply && (
                  <button
                    className="text-xs text-gray-500 hover:text-gray-700"
                    onClick={(e) => {
                      e.stopPropagation();
                      onReply();
                    }}
                  >
                    Reply
                  </button>
                )}
                {onEdit && isOwn && (
                  <button
                    className="text-xs text-gray-500 hover:text-gray-700"
                    onClick={(e) => {
                      e.stopPropagation();
                      onEdit();
                    }}
                    aria-label="Edit"
                  >
                    Edit
                  </button>
                )}
                {onDelete && isOwn && (
                  <button
                    className="text-xs text-red-500 hover:text-red-700"
                    onClick={(e) => {
                      e.stopPropagation();
                      onDelete();
                    }}
                    aria-label="Delete"
                  >
                    Delete
                  </button>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }
);

ChatMessage.displayName = 'ChatMessage';

/**
 * Typing indicator component
 */
export const TypingIndicator = React.forwardRef<HTMLDivElement, TypingIndicatorProps>(
  ({ className, testId, users, maxUsers = 3, ...props }, ref) => {
    const displayUsers = users.slice(0, maxUsers);
    const remainingCount = Math.max(0, users.length - maxUsers);

    const getTypingText = () => {
      if (displayUsers.length === 0) return '';

      if (displayUsers.length === 1) {
        if (remainingCount > 0) {
          if (remainingCount === 1) {
            return `${displayUsers[0].name} and 1 other are typing`;
          }
          return `${displayUsers[0].name} and ${remainingCount} others are typing`;
        }
        return `${displayUsers[0].name} is typing`;
      }

      if (displayUsers.length === 2) {
        if (remainingCount > 0) {
          if (remainingCount === 1) {
            return `${displayUsers[0].name}, ${displayUsers[1].name} and 1 other are typing`;
          }
          return `${displayUsers[0].name}, ${displayUsers[1].name} and ${remainingCount} others are typing`;
        }
        return `${displayUsers[0].name} and ${displayUsers[1].name} are typing`;
      }

      // 3+ users displayed
      const names = displayUsers.map(u => u.name).join(', ');
      if (remainingCount > 0) {
        if (remainingCount === 1) {
          return `${names} and 1 other are typing`;
        }
        return `${names} and ${remainingCount} others are typing`;
      }
      return `${names} are typing`;
    };

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          'flex items-center space-x-2',
          'typing-indicator',
          className
        )}
        {...props}
      >
        <div className="flex space-x-1">
          <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce typing-dot" style={{ animationDelay: '0ms' }}></span>
          <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce typing-dot" style={{ animationDelay: '150ms' }}></span>
          <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce typing-dot" style={{ animationDelay: '300ms' }}></span>
        </div>
        {users.length > 0 && (
          <span className="text-xs text-gray-500 typing-text">
            {getTypingText()}
          </span>
        )}
      </div>
    );
  }
);

TypingIndicator.displayName = 'TypingIndicator';

/**
 * Message group component for grouping consecutive messages
 */
export const MessageGroup = React.forwardRef<HTMLDivElement, MessageGroupProps>(
  ({
    className,
    testId,
    messages,
    sender,
    isOwn = false,
    showAvatar = true,
    showTimestamp = true,
    ...props
  }, ref) => {
    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          'mb-4',
          'message-group',
          isOwn ? 'message-group-own' : 'message-group-other',
          className
        )}
        {...props}
      >
        {/* Group header with avatar */}
        <div className={cn(
          'flex items-end',
          isOwn ? 'justify-end' : 'justify-start',
          'mb-1'
        )}>
          {showAvatar && !isOwn && (
            <div className="flex-shrink-0 mr-2 group-avatar">
              {sender.avatar ? (
                <img
                  src={sender.avatar}
                  alt={sender.name}
                  className="w-8 h-8 rounded-full"
                />
              ) : (
                <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center">
                  <span className="text-sm font-medium text-gray-600">
                    {sender.name.charAt(0).toUpperCase()}
                  </span>
                </div>
              )}
            </div>
          )}

          <div className={cn(
            'flex flex-col',
            isOwn && 'items-end'
          )}>
            {/* Sender name */}
            {!isOwn && (
              <div className="text-xs text-gray-500 mb-1 group-sender">
                {sender.name}
              </div>
            )}

            {/* Messages */}
            <div className="space-y-1 group-messages">
              {messages.map((message, index) => (
                <ChatMessage
                  key={index}
                  {...message}
                  sender={sender}
                  isOwn={isOwn}
                  showAvatar={false}
                  showSender={false}
                  showTimestamp={index === messages.length - 1 && showTimestamp}
                />
              ))}
            </div>
          </div>
        </div>
      </div>
    );
  }
);

MessageGroup.displayName = 'MessageGroup';

/**
 * Message bubble component
 */
export const MessageBubble = React.forwardRef<HTMLDivElement, MessageBubbleProps>(
  ({
    className,
    testId,
    children,
    variant = 'default',
    grouped = false,
    position = 'single',
    tail = false,
    ...props
  }, ref) => {
    const getBubbleStyles = () => {
      const baseStyles = 'px-3 py-2';

      let variantStyles = '';
      switch (variant) {
        case 'own':
          variantStyles = 'bg-blue-500 text-white';
          break;
        case 'system':
          variantStyles = 'bg-gray-100 text-gray-600 text-center';
          break;
        default:
          variantStyles = 'bg-gray-200 text-gray-900';
      }

      let radiusStyles = 'rounded-lg';
      if (grouped) {
        switch (position) {
          case 'first':
            radiusStyles = variant === 'own'
              ? 'rounded-t-lg rounded-bl-lg rounded-br'
              : 'rounded-t-lg rounded-br-lg rounded-bl';
            break;
          case 'middle':
            radiusStyles = variant === 'own'
              ? 'rounded-l-lg rounded-br'
              : 'rounded-r-lg rounded-bl';
            break;
          case 'last':
            radiusStyles = variant === 'own'
              ? 'rounded-b-lg rounded-tl-lg rounded-tr'
              : 'rounded-b-lg rounded-tr-lg rounded-tl';
            break;
          default:
            radiusStyles = 'rounded-lg';
        }
      }

      return cn(baseStyles, variantStyles, radiusStyles);
    };

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          getBubbleStyles(),
          'message-bubble',
          `message-bubble-${variant}`,
          `bubble-${variant}`,
          grouped && 'bubble-grouped',
          position && `message-bubble-${position} bubble-${position}`,
          tail && 'message-bubble-tail bubble-tail',
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);

MessageBubble.displayName = 'MessageBubble';

export type {
  ChatMessageProps,
  TypingIndicatorProps,
  MessageGroupProps,
  MessageBubbleProps,
  MediaContent,
  SystemContent,
  MessageStatus,
  Reaction,
  User
};