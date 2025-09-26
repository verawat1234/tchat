// T039 - ReplyMessage component with thread visualization
/**
 * ReplyMessage Component
 * Displays reply messages with visual thread connections and context preservation
 * Supports nested threading, original message preview, and keyboard navigation
 */

import React, { useMemo, useCallback, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { ReplyContent, MessagePreview, ThreadMetadata } from '../../types/ReplyContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent } from '../ui/card';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import {
  MessageCircle,
  CornerDownRight,
  Clock,
  MoreHorizontal,
  Eye,
  Flag,
  Share,
  Copy
} from 'lucide-react';

// Component Props Interface
interface ReplyMessageProps {
  message: MessageData & { content: ReplyContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onReplyToReply?: (message: MessageData) => void;
  onThreadExpand?: (threadId: string) => void;
  onThreadCollapse?: (threadId: string) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  showThreadLine?: boolean;
  compactMode?: boolean;
  isInThread?: boolean;
  threadPosition?: { current: number; total: number };
  maxDepth?: number;
  performanceMode?: boolean;
}

// Animation Variants
const replyVariants = {
  initial: { opacity: 0, x: -20, y: 10 },
  animate: { opacity: 1, x: 0, y: 0 },
  exit: { opacity: 0, x: 20, y: -10 }
};

const threadLineVariants = {
  initial: { scaleY: 0 },
  animate: { scaleY: 1 },
  exit: { scaleY: 0 }
};

const previewVariants = {
  initial: { opacity: 0, height: 0 },
  animate: { opacity: 1, height: 'auto' },
  exit: { opacity: 0, height: 0 }
};

export const ReplyMessage: React.FC<ReplyMessageProps> = ({
  message,
  onInteraction,
  onReplyToReply,
  onThreadExpand,
  onThreadCollapse,
  className,
  showAvatar = true,
  showTimestamp = true,
  showThreadLine = true,
  compactMode = false,
  isInThread = false,
  threadPosition,
  maxDepth = 5,
  performanceMode = false
}) => {
  const replyRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Calculate visual properties based on thread depth
  const threadDepth = Math.min(content.threadDepth, maxDepth);
  const indentLevel = threadDepth * (compactMode ? 16 : 24);
  const threadLineColor = useMemo(() => {
    const colors = [
      'border-blue-300',
      'border-green-300',
      'border-purple-300',
      'border-orange-300',
      'border-pink-300'
    ];
    return colors[threadDepth % colors.length];
  }, [threadDepth]);

  // Format timestamp
  const formatTimestamp = useCallback((date: Date) => {
    const now = new Date();
    const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60));

    if (diffInMinutes < 1) return 'now';
    if (diffInMinutes < 60) return `${diffInMinutes}m`;
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h`;
    return date.toLocaleDateString();
  }, []);

  // Handle interactions
  const handleInteraction = useCallback((type: string, data: Record<string, unknown> = {}) => {
    if (!onInteraction) return;

    onInteraction({
      messageId: message.id,
      interactionType: type as any,
      data,
      userId: 'current-user', // Would come from auth context
      timestamp: new Date()
    });
  }, [message.id, onInteraction]);

  // Handle reply action
  const handleReplyAction = useCallback(() => {
    if (onReplyToReply) {
      onReplyToReply(message);
    }
    handleInteraction('reply_start', { originalMessageId: message.id });
  }, [message, onReplyToReply, handleInteraction]);

  // Handle copy action
  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(content.replyText);
      handleInteraction('copy_message', { success: true });
      // Show toast notification (would integrate with toast system)
    } catch (error) {
      handleInteraction('copy_message', { success: false, error: error.message });
    }
  }, [content.replyText, handleInteraction]);

  // Handle keyboard navigation
  const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
    switch (event.key) {
      case 'r':
        if (event.ctrlKey || event.metaKey) {
          event.preventDefault();
          handleReplyAction();
        }
        break;
      case 'c':
        if (event.ctrlKey || event.metaKey) {
          event.preventDefault();
          handleCopy();
        }
        break;
      case 'Enter':
        if (event.altKey) {
          event.preventDefault();
          handleReplyAction();
        }
        break;
    }
  }, [handleReplyAction, handleCopy]);

  // Performance optimization: memoize original preview
  const originalPreview = useMemo(() => {
    if (!content.originalPreview) return null;

    return (
      <motion.div
        variants={previewVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        className="mt-2 mb-3"
      >
        <Card className="bg-muted/50 border-l-4 border-l-primary/50">
          <CardContent className="p-3">
            <div className="flex items-center gap-2 mb-1">
              <Avatar className="w-4 h-4">
                <AvatarImage src={`/avatars/${content.originalPreview.senderName.toLowerCase()}.png`} />
                <AvatarFallback className="text-xs">
                  {content.originalPreview.senderName.substring(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <span className="text-xs font-medium text-muted-foreground">
                {content.originalPreview.senderName}
              </span>
              {showTimestamp && (
                <span className="text-xs text-muted-foreground/60">
                  {formatTimestamp(content.originalPreview.timestamp)}
                </span>
              )}
            </div>
            <p className="text-sm text-muted-foreground line-clamp-2">
              {content.originalPreview.contentPreview}
            </p>
          </CardContent>
        </Card>
      </motion.div>
    );
  }, [content.originalPreview, showTimestamp, formatTimestamp]);

  // Performance optimization: skip animation in performance mode
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: replyVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.2, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={replyRef}
        className={cn(
          "reply-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        style={{
          paddingLeft: `${indentLevel}px`,
          marginTop: compactMode ? '4px' : '8px'
        }}
        data-testid={`reply-message-${message.id}`}
        data-thread-depth={threadDepth}
        tabIndex={0}
        onKeyDown={handleKeyDown}
        role="article"
        aria-label={`Reply from ${message.senderName}: ${content.replyText}`}
      >
        {/* Thread Connection Line */}
        {showThreadLine && threadDepth > 0 && (
          <motion.div
            variants={performanceMode ? {} : threadLineVariants}
            initial={performanceMode ? {} : "initial"}
            animate={performanceMode ? {} : "animate"}
            className={cn(
              "reply-thread-line absolute left-4 top-0 w-0.5 h-full",
              "before:absolute before:top-6 before:left-0 before:w-4 before:h-0.5",
              "before:border-t-2 before:border-l-2 before:border-r-2",
              "before:border-b-0 before:rounded-tl-md",
              threadLineColor,
              `before:${threadLineColor}`
            )}
            data-testid="thread-connector"
            style={{ left: `${(threadDepth - 1) * (compactMode ? 16 : 24) + 12}px` }}
          />
        )}

        {/* Main Reply Content */}
        <div className="flex gap-3 min-w-0 flex-1">
          {/* Avatar */}
          {showAvatar && (
            <motion.div
              initial={performanceMode ? {} : { scale: 0.8, opacity: 0 }}
              animate={performanceMode ? {} : { scale: 1, opacity: 1 }}
              transition={{ delay: 0.1 }}
              className="flex-shrink-0"
            >
              <Avatar className={cn(
                compactMode ? "w-6 h-6" : "w-8 h-8",
                "ring-2 ring-background"
              )}>
                <AvatarImage src={`/avatars/${message.senderName.toLowerCase()}.png`} />
                <AvatarFallback className={compactMode ? "text-xs" : "text-sm"}>
                  {message.senderName.substring(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
            </motion.div>
          )}

          {/* Content Container */}
          <div className="flex-1 min-w-0 space-y-1">
            {/* Header with sender info and timestamp */}
            <div className="flex items-center gap-2 flex-wrap">
              <span className="font-medium text-sm text-foreground truncate">
                {message.senderName}
              </span>

              {message.isOwn && (
                <span className="text-xs text-muted-foreground bg-primary/10 px-1.5 py-0.5 rounded-full">
                  You
                </span>
              )}

              {content.isThreadStart && (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <MessageCircle className="w-3 h-3 text-primary" />
                  </TooltipTrigger>
                  <TooltipContent>Thread starter</TooltipContent>
                </Tooltip>
              )}

              {showTimestamp && (
                <div className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Clock className="w-3 h-3" />
                  <time dateTime={message.timestamp.toISOString()}>
                    {formatTimestamp(message.timestamp)}
                  </time>
                </div>
              )}

              {threadPosition && (
                <span className="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded">
                  {threadPosition.current} of {threadPosition.total}
                </span>
              )}
            </div>

            {/* Original Message Preview */}
            <AnimatePresence mode="wait">
              {originalPreview}
            </AnimatePresence>

            {/* Reply Text Content */}
            <motion.div
              initial={performanceMode ? {} : { opacity: 0, y: 10 }}
              animate={performanceMode ? {} : { opacity: 1, y: 0 }}
              transition={{ delay: 0.15 }}
              className="reply-content bg-card border border-border rounded-lg p-3 shadow-sm"
            >
              <div className="flex items-start gap-2">
                <CornerDownRight className="w-4 h-4 text-primary/60 flex-shrink-0 mt-0.5" />
                <p className="text-sm text-foreground leading-relaxed whitespace-pre-wrap break-words">
                  {content.replyText}
                </p>
              </div>
            </motion.div>

            {/* Action Buttons */}
            <motion.div
              initial={performanceMode ? {} : { opacity: 0 }}
              animate={performanceMode ? {} : { opacity: 1 }}
              transition={{ delay: 0.2 }}
              className={cn(
                "flex items-center gap-1 pt-1",
                "opacity-0 group-hover:opacity-100 group-focus-within:opacity-100",
                "transition-opacity duration-200"
              )}
            >
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={handleReplyAction}
                    className="h-7 px-2 text-xs"
                    aria-label="Reply to this message"
                  >
                    <MessageCircle className="w-3 h-3 mr-1" />
                    Reply
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Reply to this message (Ctrl+R)</TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={handleCopy}
                    className="h-7 px-2 text-xs"
                    aria-label="Copy message"
                  >
                    <Copy className="w-3 h-3" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Copy message (Ctrl+C)</TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleInteraction('share_message')}
                    className="h-7 px-2 text-xs"
                    aria-label="Share message"
                  >
                    <Share className="w-3 h-3" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Share message</TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleInteraction('flag_message')}
                    className="h-7 px-2 text-xs"
                    aria-label="Flag message"
                  >
                    <Flag className="w-3 h-3" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Flag message</TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleInteraction('view_thread')}
                    className="h-7 px-2 text-xs"
                    aria-label="More actions"
                  >
                    <MoreHorizontal className="w-3 h-3" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>More actions</TooltipContent>
              </Tooltip>
            </motion.div>
          </div>
        </div>

        {/* Thread Expansion Indicator */}
        {content.threadMetadata && (
          <motion.div
            initial={performanceMode ? {} : { opacity: 0, scale: 0.9 }}
            animate={performanceMode ? {} : { opacity: 1, scale: 1 }}
            transition={{ delay: 0.3 }}
            className="mt-2 ml-12 text-xs text-muted-foreground flex items-center gap-2"
          >
            <MessageCircle className="w-3 h-3" />
            <span>
              {content.threadMetadata.totalThreadMessages} messages in thread
            </span>
            <span>•</span>
            <span>
              {content.threadMetadata.threadParticipants.length} participants
            </span>
          </motion.div>
        )}

        {/* Performance Debug Info (Development Only) */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            D:{threadDepth} P:{performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedReplyMessage = React.memo(ReplyMessage, (prevProps, nextProps) => {
  // Custom comparison function for performance
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.performanceMode === nextProps.performanceMode &&
    JSON.stringify(prevProps.threadPosition) === JSON.stringify(nextProps.threadPosition)
  );
});

MemoizedReplyMessage.displayName = 'MemoizedReplyMessage';

// Export both versions
export default ReplyMessage;