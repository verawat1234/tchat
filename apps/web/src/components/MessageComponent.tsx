/**
 * Comprehensive Message Component
 *
 * Advanced message display and interaction component with threading,
 * replies, reactions, file attachments, and rich content support.
 *
 * Features:
 * - Message threading and reply functionality
 * - Rich message types (text, media, files, etc.)
 * - Real-time reactions and emoji support
 * - Message editing and deletion
 * - File attachments with preview
 * - Typing indicators and delivery receipts
 * - Message search and filtering
 * - Accessibility and keyboard navigation
 * - Mobile-responsive design
 */

import React, { useState, useEffect, useCallback, useRef, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { formatDistanceToNow, format } from 'date-fns';
import {
  MessageSquare,
  MoreHorizontal,
  Reply,
  Edit,
  Trash2,
  Copy,
  Forward,
  Pin,
  Star,
  Download,
  ExternalLink,
  Eye,
  EyeOff,
  Volume2,
  VolumeX,
  Play,
  Pause,
  FileText,
  Image as ImageIcon,
  Video,
  Music,
  Paperclip,
  MapPin,
  Phone,
  Calendar,
  User,
  Users,
  Hash,
  AtSign,
  Smile,
  ThumbsUp,
  Heart,
  Laugh,
  Angry,
  Cry,
  Check,
  CheckCheck,
  Clock,
  AlertCircle,
} from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Textarea } from './ui/textarea';
import { Badge } from './ui/badge/badge';
import { Avatar } from './ui/avatar';
import { ScrollArea } from './ui/scroll-area';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { useDialog, useDialogHelpers } from './DialogSystem';
import { useMessages } from '../services/messageService';
import { useDeliveryReceipts } from '../services/deliveryReceiptService';
import { MessageStatusIndicator, ReadReceiptAvatars } from './MessageStatusIndicator';
import { useMessageReading } from '../hooks/useMessageReading';
import type {
  Message,
  MessageContent,
  MessageAttachment,
  MessageReaction,
  TypingIndicator,
  MessageType,
} from '../services/messageService';
import type { DeliveryStatus } from '../types/MessageTypes';

// =============================================================================
// Type Definitions
// =============================================================================

export interface MessageComponentProps {
  chatId: string;
  currentUserId: string;
  showHeader?: boolean;
  showTypingIndicators?: boolean;
  allowThreading?: boolean;
  allowReactions?: boolean;
  allowEditing?: boolean;
  compact?: boolean;
  className?: string;
}

export interface MessageItemProps {
  message: Message;
  currentUserId: string;
  isGrouped?: boolean;
  showTimestamp?: boolean;
  allowThreading?: boolean;
  allowReactions?: boolean;
  allowEditing?: boolean;
  onReply?: (message: Message) => void;
  onEdit?: (message: Message) => void;
  onDelete?: (message: Message) => void;
  onReact?: (message: Message, emoji: string) => void;
  onThreadOpen?: (message: Message) => void;
}

export interface MessageInputProps {
  chatId: string;
  onSend: (content: MessageContent, type: MessageType) => Promise<void>;
  onTyping?: (isTyping: boolean) => void;
  replyTo?: Message;
  onCancelReply?: () => void;
  placeholder?: string;
  disabled?: boolean;
}

export interface MessageThreadProps {
  parentMessage: Message;
  chatId: string;
  currentUserId: string;
  onClose: () => void;
}

// =============================================================================
// Main Message Component
// =============================================================================

export function MessageComponent({
  chatId,
  currentUserId,
  showHeader = true,
  showTypingIndicators = true,
  allowThreading = true,
  allowReactions = true,
  allowEditing = true,
  compact = false,
  className = '',
}: MessageComponentProps) {
  const {
    messages,
    loading,
    typingUsers,
    sendMessage,
    editMessage,
    deleteMessage,
    addReaction,
    sendTyping,
  } = useMessages(chatId);

  const [replyingTo, setReplyingTo] = useState<Message | null>(null);
  const [editingMessage, setEditingMessage] = useState<Message | null>(null);
  const [threadMessage, setThreadMessage] = useState<Message | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [filteredMessages, setFilteredMessages] = useState<Message[]>([]);

  const { showConfirmation, showForm, openDialog } = useDialogHelpers();
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Filter messages based on search
  useEffect(() => {
    if (!searchQuery.trim()) {
      setFilteredMessages(messages);
      return;
    }

    const filtered = messages.filter(message =>
      message.content.text?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      message.content.html?.toLowerCase().includes(searchQuery.toLowerCase())
    );
    setFilteredMessages(filtered);
  }, [messages, searchQuery]);

  // Scroll to bottom when new messages arrive
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages.length]);

  // Group consecutive messages from same user
  const groupedMessages = useMemo(() => {
    const groups: Array<{ user: string; messages: Message[] }> = [];
    let currentGroup: { user: string; messages: Message[] } | null = null;

    filteredMessages.forEach(message => {
      if (!currentGroup || currentGroup.user !== message.senderId) {
        currentGroup = { user: message.senderId, messages: [message] };
        groups.push(currentGroup);
      } else {
        currentGroup.messages.push(message);
      }
    });

    return groups;
  }, [filteredMessages]);

  // Handle message actions
  const handleReply = useCallback((message: Message) => {
    setReplyingTo(message);
    inputRef.current?.focus();
  }, []);

  const handleEdit = useCallback((message: Message) => {
    setEditingMessage(message);
    showForm(
      'Edit Message',
      <MessageEditForm
        message={message}
        onSave={async (content) => {
          await editMessage(message.id, content);
          setEditingMessage(null);
        }}
        onCancel={() => setEditingMessage(null)}
      />,
      [],
      { size: 'lg' }
    );
  }, [editMessage, showForm]);

  const handleDelete = useCallback((message: Message) => {
    showConfirmation(
      'Delete Message',
      'Are you sure you want to delete this message? This action cannot be undone.',
      async () => {
        await deleteMessage(message.id);
      }
    );
  }, [deleteMessage, showConfirmation]);

  const handleReact = useCallback(async (message: Message, emoji: string) => {
    await addReaction(message.id, emoji);
  }, [addReaction]);

  const handleThreadOpen = useCallback((message: Message) => {
    setThreadMessage(message);
    openDialog({
      type: 'sheet',
      title: 'Thread',
      content: (
        <MessageThread
          parentMessage={message}
          chatId={chatId}
          currentUserId={currentUserId}
          onClose={() => setThreadMessage(null)}
        />
      ),
      size: 'lg',
      closable: true,
    });
  }, [openDialog, chatId, currentUserId]);

  const handleSendMessage = useCallback(async (content: MessageContent, type: MessageType) => {
    await sendMessage({
      type,
      content,
      replyToId: replyingTo?.id,
    });
    setReplyingTo(null);
  }, [sendMessage, replyingTo]);

  const handleTyping = useCallback((isTyping: boolean) => {
    sendTyping(isTyping);
  }, [sendTyping]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
          className="h-8 w-8 border-2 border-primary border-t-transparent rounded-full"
        />
      </div>
    );
  }

  return (
    <div className={`flex flex-col h-full bg-white ${className}`}>
      {/* Header */}
      {showHeader && (
        <div className="flex items-center justify-between p-4 border-b bg-gray-50">
          <div className="flex items-center gap-3">
            <h2 className="text-lg font-semibold text-gray-900">Messages</h2>
            <Badge variant="outline">
              {messages.length} message{messages.length !== 1 ? 's' : ''}
            </Badge>
          </div>
          <div className="flex items-center gap-2">
            <Input
              placeholder="Search messages..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-48"
            />
          </div>
        </div>
      )}

      {/* Messages List */}
      <ScrollArea className="flex-1 p-4">
        <div className="space-y-4">
          <AnimatePresence mode="popLayout">
            {groupedMessages.map((group, groupIndex) => (
              <div key={`group-${groupIndex}`} className="space-y-1">
                {group.messages.map((message, messageIndex) => (
                  <MessageItem
                    key={message.id}
                    message={message}
                    currentUserId={currentUserId}
                    isGrouped={messageIndex > 0}
                    showTimestamp={messageIndex === group.messages.length - 1}
                    allowThreading={allowThreading}
                    allowReactions={allowReactions}
                    allowEditing={allowEditing}
                    onReply={handleReply}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                    onReact={handleReact}
                    onThreadOpen={handleThreadOpen}
                  />
                ))}
              </div>
            ))}
          </AnimatePresence>

          {/* Typing Indicators */}
          {showTypingIndicators && typingUsers.length > 0 && (
            <TypingIndicators users={typingUsers} />
          )}

          <div ref={messagesEndRef} />
        </div>
      </ScrollArea>

      {/* Reply Preview */}
      {replyingTo && (
        <ReplyPreview
          message={replyingTo}
          onCancel={() => setReplyingTo(null)}
        />
      )}

      {/* Message Input */}
      <MessageInput
        chatId={chatId}
        onSend={handleSendMessage}
        onTyping={handleTyping}
        replyTo={replyingTo}
        onCancelReply={() => setReplyingTo(null)}
        disabled={false}
      />
    </div>
  );
}

// =============================================================================
// Message Item Component
// =============================================================================

function MessageItem({
  message,
  currentUserId,
  isGrouped = false,
  showTimestamp = true,
  allowThreading = true,
  allowReactions = true,
  allowEditing = true,
  onReply,
  onEdit,
  onDelete,
  onReact,
  onThreadOpen,
}: MessageItemProps) {
  const [showActions, setShowActions] = useState(false);
  const [showReactions, setShowReactions] = useState(false);
  const isOwnMessage = message.senderId === currentUserId;

  // Use delivery receipt hook for read receipts and status
  const { receipts, status, readReceipts } = useDeliveryReceipts(
    message.id,
    message.chatId
  );

  // Use message reading hook for automatic read tracking
  const { elementRef } = useMessageReading(
    message.id,
    message.chatId,
    isOwnMessage,
    {
      enabled: true,
      readDelay: 2000, // Mark as read after 2 seconds of viewing
      threshold: 0.5   // 50% of message must be visible
    }
  );

  const handleReactionSelect = (emoji: string) => {
    onReact?.(message, emoji);
    setShowReactions(false);
  };


  return (
    <motion.div
      ref={elementRef}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -10 }}
      className={`group flex gap-3 ${isOwnMessage ? 'flex-row-reverse' : ''} ${
        isGrouped ? 'mt-1' : 'mt-4'
      }`}
      onMouseEnter={() => setShowActions(true)}
      onMouseLeave={() => setShowActions(false)}
    >
      {/* Avatar */}
      {!isGrouped && (
        <Avatar className="h-8 w-8 flex-shrink-0">
          <div className="h-full w-full bg-primary/20 flex items-center justify-center">
            <span className="text-xs font-medium">
              {message.senderId.charAt(0).toUpperCase()}
            </span>
          </div>
        </Avatar>
      )}

      {/* Message Content */}
      <div className={`flex-1 min-w-0 ${isGrouped && !isOwnMessage ? 'ml-11' : ''}`}>
        {/* Sender Name & Timestamp */}
        {!isGrouped && (
          <div className={`flex items-center gap-2 mb-1 ${isOwnMessage ? 'justify-end' : ''}`}>
            <span className="text-sm font-medium text-gray-900">
              {isOwnMessage ? 'You' : `User ${message.senderId}`}
            </span>
            <span className="text-xs text-gray-500">
              {format(new Date(message.timestamp), 'HH:mm')}
            </span>
          </div>
        )}

        {/* Message Bubble */}
        <div
          className={`relative max-w-lg ${
            isOwnMessage
              ? 'ml-auto bg-blue-500 text-white rounded-l-2xl rounded-br-sm'
              : 'bg-gray-100 text-gray-900 rounded-r-2xl rounded-bl-sm'
          } rounded-t-2xl px-4 py-2`}
        >
          {/* Reply Context */}
          {message.replyToId && (
            <div className="mb-2 p-2 bg-black/10 rounded border-l-2 border-current opacity-75">
              <p className="text-xs">Replying to a message</p>
            </div>
          )}

          {/* Message Content */}
          <MessageContentRenderer content={message.content} type={message.type} />

          {/* Attachments */}
          {message.attachments && message.attachments.length > 0 && (
            <MessageAttachments attachments={message.attachments} />
          )}

          {/* Message Status & Read Receipts */}
          {isOwnMessage && (
            <div className="flex items-center justify-end mt-1 space-x-2 opacity-75">
              <MessageStatusIndicator
                status={status}
                timestamp={message.timestamp}
                size="sm"
              />
              {readReceipts.length > 0 && (
                <ReadReceiptAvatars
                  readByUsers={readReceipts.map(receipt => ({
                    id: receipt.userId,
                    name: receipt.userName,
                    avatar: receipt.userAvatar,
                    readAt: receipt.timestamp,
                  }))}
                  maxShow={3}
                />
              )}
            </div>
          )}
        </div>

        {/* Reactions */}
        {message.reactions.length > 0 && (
          <MessageReactions
            reactions={message.reactions}
            onReact={handleReactionSelect}
            className="mt-1"
          />
        )}

        {/* Thread Info */}
        {message.threadId && allowThreading && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onThreadOpen?.(message)}
            className="mt-1 text-xs text-blue-600 hover:text-blue-700"
          >
            <MessageSquare className="h-3 w-3 mr-1" />
            View thread
          </Button>
        )}

        {/* Timestamp for grouped messages */}
        {showTimestamp && isGrouped && (
          <p className={`text-xs text-gray-500 mt-1 ${isOwnMessage ? 'text-right' : ''}`}>
            {format(new Date(message.timestamp), 'HH:mm')}
          </p>
        )}
      </div>

      {/* Message Actions */}
      <AnimatePresence>
        {showActions && (
          <motion.div
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.8 }}
            className={`flex items-center gap-1 ${isOwnMessage ? 'flex-row-reverse' : ''}`}
          >
            {allowReactions && (
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setShowReactions(!showReactions)}
                      className="h-6 w-6 p-0"
                    >
                      <Smile className="h-3 w-3" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Add reaction</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            )}

            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onReply?.(message)}
                    className="h-6 w-6 p-0"
                  >
                    <Reply className="h-3 w-3" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Reply</TooltipContent>
              </Tooltip>
            </TooltipProvider>

            {allowThreading && (
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => onThreadOpen?.(message)}
                      className="h-6 w-6 p-0"
                    >
                      <MessageSquare className="h-3 w-3" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Start thread</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            )}

            {isOwnMessage && allowEditing && (
              <>
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => onEdit?.(message)}
                        className="h-6 w-6 p-0"
                      >
                        <Edit className="h-3 w-3" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Edit</TooltipContent>
                  </Tooltip>
                </TooltipProvider>

                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => onDelete?.(message)}
                        className="h-6 w-6 p-0 text-red-600 hover:text-red-700"
                      >
                        <Trash2 className="h-3 w-3" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Delete</TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </>
            )}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Reaction Picker */}
      {showReactions && (
        <ReactionPicker
          onSelect={handleReactionSelect}
          onClose={() => setShowReactions(false)}
          position={isOwnMessage ? 'left' : 'right'}
        />
      )}
    </motion.div>
  );
}

// =============================================================================
// Supporting Components
// =============================================================================

function MessageContentRenderer({ content, type }: { content: MessageContent; type: MessageType }) {
  switch (type) {
    case 'text':
      return (
        <div className="text-sm">
          {content.html ? (
            <div dangerouslySetInnerHTML={{ __html: content.html }} />
          ) : (
            <p className="whitespace-pre-wrap">{content.text}</p>
          )}
        </div>
      );

    case 'image':
      return (
        <div className="space-y-2">
          {content.text && <p className="text-sm">{content.text}</p>}
          <img
            src={content.data?.url}
            alt={content.data?.alt || 'Image'}
            className="max-w-full rounded-lg"
          />
        </div>
      );

    case 'video':
      return (
        <div className="space-y-2">
          {content.text && <p className="text-sm">{content.text}</p>}
          <video
            src={content.data?.url}
            controls
            className="max-w-full rounded-lg"
          />
        </div>
      );

    case 'audio':
      return (
        <div className="space-y-2">
          {content.text && <p className="text-sm">{content.text}</p>}
          <audio src={content.data?.url} controls className="w-full" />
        </div>
      );

    case 'location':
      return (
        <div className="flex items-center gap-2 text-sm">
          <MapPin className="h-4 w-4" />
          <span>{content.data?.address || 'Location shared'}</span>
        </div>
      );

    case 'contact':
      return (
        <div className="flex items-center gap-2 text-sm">
          <User className="h-4 w-4" />
          <span>{content.data?.name || 'Contact shared'}</span>
        </div>
      );

    default:
      return <p className="text-sm">{content.text}</p>;
  }
}

function MessageAttachments({ attachments }: { attachments: MessageAttachment[] }) {
  return (
    <div className="mt-2 space-y-2">
      {attachments.map(attachment => (
        <div
          key={attachment.id}
          className="flex items-center gap-2 p-2 bg-black/10 rounded border"
        >
          <FileIcon type={attachment.type} />
          <div className="flex-1 min-w-0">
            <p className="text-xs font-medium truncate">{attachment.name}</p>
            <p className="text-xs opacity-75">
              {formatFileSize(attachment.size)}
            </p>
          </div>
          <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
            <Download className="h-3 w-3" />
          </Button>
        </div>
      ))}
    </div>
  );
}

function FileIcon({ type }: { type: MessageAttachment['type'] }) {
  switch (type) {
    case 'image':
      return <ImageIcon className="h-4 w-4" />;
    case 'video':
      return <Video className="h-4 w-4" />;
    case 'audio':
      return <Music className="h-4 w-4" />;
    default:
      return <FileText className="h-4 w-4" />;
  }
}

function MessageReactions({
  reactions,
  onReact,
  className = '',
}: {
  reactions: MessageReaction[];
  onReact: (emoji: string) => void;
  className?: string;
}) {
  // Group reactions by emoji
  const groupedReactions = reactions.reduce((acc, reaction) => {
    if (!acc[reaction.emoji]) {
      acc[reaction.emoji] = [];
    }
    acc[reaction.emoji].push(reaction);
    return acc;
  }, {} as Record<string, MessageReaction[]>);

  return (
    <div className={`flex flex-wrap gap-1 ${className}`}>
      {Object.entries(groupedReactions).map(([emoji, reactionList]) => (
        <Button
          key={emoji}
          variant="outline"
          size="sm"
          onClick={() => onReact(emoji)}
          className="h-6 px-2 text-xs"
        >
          <span className="mr-1">{emoji}</span>
          {reactionList.length}
        </Button>
      ))}
    </div>
  );
}

function ReactionPicker({
  onSelect,
  onClose,
  position = 'right',
}: {
  onSelect: (emoji: string) => void;
  onClose: () => void;
  position?: 'left' | 'right';
}) {
  const commonEmojis = ['👍', '❤️', '😂', '😮', '😢', '😡', '👏', '🎉'];

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.8 }}
      animate={{ opacity: 1, scale: 1 }}
      exit={{ opacity: 0, scale: 0.8 }}
      className={`absolute z-10 bg-white border rounded-lg shadow-lg p-2 ${
        position === 'left' ? 'right-0' : 'left-0'
      }`}
    >
      <div className="flex gap-1">
        {commonEmojis.map(emoji => (
          <Button
            key={emoji}
            variant="ghost"
            size="sm"
            onClick={() => onSelect(emoji)}
            className="h-8 w-8 p-0 text-lg hover:bg-gray-100"
          >
            {emoji}
          </Button>
        ))}
      </div>
    </motion.div>
  );
}

function TypingIndicators({ users }: { users: TypingIndicator[] }) {
  const typingNames = users
    .filter(user => user.isTyping)
    .map(user => user.userName)
    .slice(0, 3);

  if (typingNames.length === 0) return null;

  const displayText = typingNames.length === 1
    ? `${typingNames[0]} is typing...`
    : typingNames.length === 2
    ? `${typingNames[0]} and ${typingNames[1]} are typing...`
    : `${typingNames[0]}, ${typingNames[1]} and ${typingNames.length - 2} others are typing...`;

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -10 }}
      className="flex items-center gap-2 p-3 text-sm text-gray-600 italic"
    >
      <div className="flex gap-1">
        {[1, 2, 3].map(i => (
          <motion.div
            key={i}
            animate={{ y: [0, -4, 0] }}
            transition={{
              repeat: Infinity,
              duration: 1,
              delay: i * 0.2,
            }}
            className="h-1 w-1 bg-gray-400 rounded-full"
          />
        ))}
      </div>
      {displayText}
    </motion.div>
  );
}

function ReplyPreview({ message, onCancel }: { message: Message; onCancel: () => void }) {
  return (
    <div className="flex items-center gap-3 p-3 bg-gray-50 border-t">
      <div className="flex-1">
        <div className="flex items-center gap-2">
          <Reply className="h-4 w-4 text-gray-500" />
          <span className="text-sm font-medium text-gray-700">
            Replying to {message.senderId === 'current-user' ? 'yourself' : `User ${message.senderId}`}
          </span>
        </div>
        <p className="text-sm text-gray-600 truncate mt-1">
          {message.content.text || 'Media message'}
        </p>
      </div>
      <Button variant="ghost" size="sm" onClick={onCancel}>
        <X className="h-4 w-4" />
      </Button>
    </div>
  );
}

function MessageInput({
  chatId,
  onSend,
  onTyping,
  replyTo,
  onCancelReply,
  placeholder = 'Type a message...',
  disabled = false,
}: MessageInputProps) {
  const [message, setMessage] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const typingTimeoutRef = useRef<NodeJS.Timeout>();

  const handleInputChange = useCallback((value: string) => {
    setMessage(value);

    // Handle typing indicators
    if (value.trim() && !isTyping) {
      setIsTyping(true);
      onTyping?.(true);
    }

    // Clear previous timeout
    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    // Set new timeout to stop typing
    typingTimeoutRef.current = setTimeout(() => {
      setIsTyping(false);
      onTyping?.(false);
    }, 1000);
  }, [isTyping, onTyping]);

  const handleSend = useCallback(async () => {
    if (!message.trim() || disabled) return;

    try {
      await onSend({ text: message.trim() }, 'text');
      setMessage('');
      setIsTyping(false);
      onTyping?.(false);
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  }, [message, disabled, onSend, onTyping]);

  const handleKeyPress = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  }, [handleSend]);

  return (
    <div className="border-t bg-white p-4">
      <div className="flex items-end gap-3">
        <div className="flex-1">
          <Textarea
            value={message}
            onChange={(e) => handleInputChange(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder={placeholder}
            disabled={disabled}
            className="min-h-[40px] max-h-32 resize-none"
            rows={1}
          />
        </div>
        <Button
          onClick={handleSend}
          disabled={!message.trim() || disabled}
          className="flex-shrink-0"
        >
          Send
        </Button>
      </div>
    </div>
  );
}

function MessageEditForm({
  message,
  onSave,
  onCancel,
}: {
  message: Message;
  onSave: (content: MessageContent) => Promise<void>;
  onCancel: () => void;
}) {
  const [content, setContent] = useState(message.content.text || '');
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    setSaving(true);
    try {
      await onSave({ text: content });
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-4">
      <Textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        className="min-h-[100px]"
        placeholder="Edit your message..."
      />
      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onCancel} disabled={saving}>
          Cancel
        </Button>
        <Button onClick={handleSave} disabled={saving || !content.trim()}>
          {saving ? 'Saving...' : 'Save Changes'}
        </Button>
      </div>
    </div>
  );
}

function MessageThread({
  parentMessage,
  chatId,
  currentUserId,
  onClose,
}: MessageThreadProps) {
  // This would typically load thread messages from the API
  return (
    <div className="h-full flex flex-col">
      <div className="flex items-center justify-between p-4 border-b">
        <h3 className="font-semibold">Thread</h3>
        <Button variant="ghost" size="sm" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>

      <div className="flex-1 p-4">
        <MessageItem
          message={parentMessage}
          currentUserId={currentUserId}
          allowThreading={false}
          allowReactions={true}
          allowEditing={false}
        />
        {/* Thread messages would be loaded and displayed here */}
      </div>

      <MessageInput
        chatId={chatId}
        onSend={async () => {
          // Handle thread message sending
        }}
        placeholder="Reply to thread..."
      />
    </div>
  );
}

// =============================================================================
// Utility Functions
// =============================================================================

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Export main component
export default MessageComponent;