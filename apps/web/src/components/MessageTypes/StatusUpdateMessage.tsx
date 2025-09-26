// T048 - StatusUpdateMessage component with activity tracking
/**
 * StatusUpdateMessage Component
 * Displays status updates with activity indicators, reactions, and presence information
 * Supports rich status content, location sharing, and social interactions
 */

import React, { useState, useCallback, useMemo, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { StatusUpdateContent, StatusType, ActivityType, StatusReaction } from '../../types/StatusUpdateContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Separator } from '../ui/separator';
import {
  Circle,
  MapPin,
  Clock,
  Users,
  Heart,
  ThumbsUp,
  Smile,
  Zap,
  Coffee,
  Plane,
  Car,
  Home,
  Building,
  Music,
  Camera,
  Gamepad2,
  BookOpen,
  Utensils,
  Dumbbell,
  Moon,
  Sun,
  CloudRain,
  Snowflake,
  Thermometer,
  Wind,
  Eye,
  Share,
  MessageCircle,
  MoreHorizontal,
  Calendar,
  Star,
  AlertCircle
} from 'lucide-react';

// Component Props Interface
interface StatusUpdateMessageProps {
  message: MessageData & { content: StatusUpdateContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onReaction?: (statusId: string, reaction: StatusReaction) => void;
  onShare?: (statusId: string) => void;
  onComment?: (statusId: string) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  showReactions?: boolean;
  showLocation?: boolean;
  performanceMode?: boolean;
}

// Animation Variants
const statusVariants = {
  initial: { opacity: 0, scale: 0.95, y: 20 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.95, y: -20 }
};

const activityVariants = {
  initial: { opacity: 0, x: -20 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 20 }
};

const reactionVariants = {
  initial: { opacity: 0, scale: 0.5 },
  animate: { opacity: 1, scale: 1 },
  exit: { opacity: 0, scale: 0.5 }
};

const presenceVariants = {
  initial: { scale: 0 },
  animate: { scale: 1 },
  exit: { scale: 0 }
};

export const StatusUpdateMessage: React.FC<StatusUpdateMessageProps> = ({
  message,
  onInteraction,
  onReaction,
  onShare,
  onComment,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  showReactions = true,
  showLocation = true,
  performanceMode = false
}) => {
  const statusRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Status state
  const [userReaction, setUserReaction] = useState<StatusReaction | null>(
    content.reactions?.find(r => r.userId === 'current-user')?.type || null
  );
  const [showAllReactions, setShowAllReactions] = useState(false);
  const [isLiked, setIsLiked] = useState(false);

  // Calculate time since status
  const timeSince = useMemo(() => {
    const now = new Date();
    const statusTime = new Date(content.timestamp);
    const diffInMinutes = Math.floor((now.getTime() - statusTime.getTime()) / (1000 * 60));

    if (diffInMinutes < 1) return 'now';
    if (diffInMinutes < 60) return `${diffInMinutes}m`;
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h`;
    return `${Math.floor(diffInMinutes / 1440)}d`;
  }, [content.timestamp]);

  // Get status type styling
  const getStatusTypeColor = useCallback((type: StatusType) => {
    switch (type) {
      case StatusType.ONLINE: return 'bg-green-500';
      case StatusType.AWAY: return 'bg-yellow-500';
      case StatusType.BUSY: return 'bg-red-500';
      case StatusType.INVISIBLE: return 'bg-gray-400';
      case StatusType.CUSTOM: return 'bg-blue-500';
      default: return 'bg-gray-400';
    }
  }, []);

  // Get activity icon
  const getActivityIcon = useCallback((activity: ActivityType) => {
    switch (activity) {
      case ActivityType.WORKING: return <Building className="w-4 h-4" />;
      case ActivityType.TRAVELING: return <Plane className="w-4 h-4" />;
      case ActivityType.EATING: return <Utensils className="w-4 h-4" />;
      case ActivityType.EXERCISING: return <Dumbbell className="w-4 h-4" />;
      case ActivityType.SLEEPING: return <Moon className="w-4 h-4" />;
      case ActivityType.STUDYING: return <BookOpen className="w-4 h-4" />;
      case ActivityType.GAMING: return <Gamepad2 className="w-4 h-4" />;
      case ActivityType.LISTENING_MUSIC: return <Music className="w-4 h-4" />;
      case ActivityType.WATCHING: return <Eye className="w-4 h-4" />;
      case ActivityType.DRIVING: return <Car className="w-4 h-4" />;
      case ActivityType.AT_HOME: return <Home className="w-4 h-4" />;
      case ActivityType.COFFEE: return <Coffee className="w-4 h-4" />;
      case ActivityType.PHOTOGRAPHY: return <Camera className="w-4 h-4" />;
      default: return <Circle className="w-4 h-4" />;
    }
  }, []);

  // Get weather icon
  const getWeatherIcon = useCallback((condition: string) => {
    switch (condition.toLowerCase()) {
      case 'sunny': case 'clear': return <Sun className="w-4 h-4 text-yellow-500" />;
      case 'rainy': case 'rain': return <CloudRain className="w-4 h-4 text-blue-500" />;
      case 'snowy': case 'snow': return <Snowflake className="w-4 h-4 text-blue-200" />;
      case 'windy': case 'wind': return <Wind className="w-4 h-4 text-gray-500" />;
      default: return <Thermometer className="w-4 h-4 text-gray-500" />;
    }
  }, []);

  // Handle reaction
  const handleReaction = useCallback((reactionType: StatusReaction) => {
    if (readonly) return;

    const newReaction = userReaction === reactionType ? null : reactionType;
    setUserReaction(newReaction);

    if (onReaction) {
      onReaction(content.id, reactionType);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'status_reaction',
        data: {
          statusId: content.id,
          reaction: reactionType,
          action: newReaction ? 'add' : 'remove'
        },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, userReaction, onReaction, content.id, onInteraction, message.id]);

  // Handle like toggle
  const handleLike = useCallback(() => {
    if (readonly) return;

    setIsLiked(!isLiked);
    handleReaction(StatusReaction.LIKE);
  }, [readonly, isLiked, handleReaction]);

  // Handle share
  const handleShare = useCallback(() => {
    if (readonly) return;

    if (onShare) {
      onShare(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'status_share',
        data: { statusId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, onShare, content.id, onInteraction, message.id]);

  // Handle comment
  const handleComment = useCallback(() => {
    if (readonly) return;

    if (onComment) {
      onComment(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'status_comment',
        data: { statusId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, onComment, content.id, onInteraction, message.id]);

  // Get reaction counts
  const reactionCounts = useMemo(() => {
    if (!content.reactions) return {};

    return content.reactions.reduce((counts, reaction) => {
      counts[reaction.type] = (counts[reaction.type] || 0) + 1;
      return counts;
    }, {} as Record<StatusReaction, number>);
  }, [content.reactions]);

  // Get visible reactions
  const visibleReactions = useMemo(() => {
    const reactions = Object.entries(reactionCounts);
    return showAllReactions ? reactions : reactions.slice(0, 3);
  }, [reactionCounts, showAllReactions]);

  // Get reaction icon
  const getReactionIcon = useCallback((reaction: StatusReaction) => {
    switch (reaction) {
      case StatusReaction.LIKE: return <ThumbsUp className="w-3 h-3" />;
      case StatusReaction.LOVE: return <Heart className="w-3 h-3" />;
      case StatusReaction.LAUGH: return <Smile className="w-3 h-3" />;
      case StatusReaction.WOW: return <Zap className="w-3 h-3" />;
      case StatusReaction.CONGRATS: return <Star className="w-3 h-3" />;
      case StatusReaction.SUPPORT: return <Users className="w-3 h-3" />;
      default: return <Circle className="w-3 h-3" />;
    }
  }, []);

  // Performance optimization
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: statusVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={statusRef}
        className={cn(
          "status-update-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`status-update-message-${message.id}`}
        data-status-type={content.statusType}
        role="article"
        aria-label={`Status update from ${message.senderName}: ${content.text || content.activity}`}
      >
        <Card className="status-card">
          {/* Header */}
          <CardHeader className="space-y-3">
            <div className="flex items-start justify-between gap-3">
              <div className="flex items-center gap-3 min-w-0 flex-1">
                {showAvatar && (
                  <motion.div
                    initial={performanceMode ? {} : { scale: 0.8, opacity: 0 }}
                    animate={performanceMode ? {} : { scale: 1, opacity: 1 }}
                    transition={{ delay: 0.1 }}
                    className="relative"
                  >
                    <Avatar className={cn(compactMode ? "w-10 h-10" : "w-12 h-12")}>
                      <AvatarImage src={`/avatars/${message.senderName.toLowerCase()}.png`} />
                      <AvatarFallback>
                        {message.senderName.substring(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>

                    {/* Presence indicator */}
                    <motion.div
                      variants={performanceMode ? {} : presenceVariants}
                      initial={performanceMode ? {} : "initial"}
                      animate={performanceMode ? {} : "animate"}
                      className={cn(
                        "absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-background",
                        getStatusTypeColor(content.statusType)
                      )}
                      title={content.statusType}
                    />
                  </motion.div>
                )}

                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="font-semibold text-foreground truncate">
                      {message.senderName}
                    </span>
                    {message.isOwn && (
                      <Badge variant="secondary" className="text-xs">You</Badge>
                    )}
                    {content.isHighlighted && (
                      <Badge variant="default" className="text-xs">
                        <Star className="w-3 h-3 mr-1" />
                        Featured
                      </Badge>
                    )}
                  </div>

                  {/* Activity and timing */}
                  <div className="flex items-center gap-2 flex-wrap text-sm text-muted-foreground mt-1">
                    {content.activity && (
                      <motion.div
                        variants={performanceMode ? {} : activityVariants}
                        initial={performanceMode ? {} : "initial"}
                        animate={performanceMode ? {} : "animate"}
                        className="flex items-center gap-1"
                      >
                        {getActivityIcon(content.activity)}
                        <span>{content.activity.replace('_', ' ')}</span>
                      </motion.div>
                    )}

                    {showTimestamp && (
                      <div className="flex items-center gap-1">
                        <Clock className="w-3 h-3" />
                        <span>{timeSince}</span>
                      </div>
                    )}

                    {content.expiresAt && (
                      <div className="flex items-center gap-1 text-orange-600">
                        <AlertCircle className="w-3 h-3" />
                        <span>Expires in {Math.floor((new Date(content.expiresAt).getTime() - Date.now()) / (1000 * 60 * 60))}h</span>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-1">
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleShare}
                      disabled={readonly}
                      className="h-8 w-8 p-0"
                    >
                      <Share className="w-4 h-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Share status</TooltipContent>
                </Tooltip>
              </div>
            </div>
          </CardHeader>

          {/* Status content */}
          <CardContent className="space-y-4">
            {/* Status text */}
            {content.text && (
              <motion.div
                initial={performanceMode ? {} : { opacity: 0, y: 10 }}
                animate={performanceMode ? {} : { opacity: 1, y: 0 }}
                transition={{ delay: 0.1 }}
                className="space-y-2"
              >
                <p className="text-foreground leading-relaxed">
                  {content.text}
                </p>
              </motion.div>
            )}

            {/* Status media */}
            {content.media && (
              <motion.div
                initial={performanceMode ? {} : { opacity: 0, scale: 0.95 }}
                animate={performanceMode ? {} : { opacity: 1, scale: 1 }}
                transition={{ delay: 0.2 }}
                className="aspect-video rounded-lg overflow-hidden bg-muted"
              >
                {content.media.type === 'image' ? (
                  <img
                    src={content.media.url}
                    alt={content.media.alt || 'Status image'}
                    className="w-full h-full object-cover"
                  />
                ) : content.media.type === 'video' ? (
                  <video
                    src={content.media.url}
                    poster={content.media.thumbnailUrl}
                    className="w-full h-full object-cover"
                    controls
                  />
                ) : null}
              </motion.div>
            )}

            {/* Location and context */}
            <div className="flex flex-wrap gap-4 text-sm">
              {content.location && showLocation && (
                <motion.div
                  variants={performanceMode ? {} : activityVariants}
                  initial={performanceMode ? {} : "initial"}
                  animate={performanceMode ? {} : "animate"}
                  transition={{ delay: 0.3 }}
                  className="flex items-center gap-2 text-muted-foreground"
                >
                  <MapPin className="w-4 h-4" />
                  <span>{content.location.name}</span>
                  {content.location.coordinates && (
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-6 px-2 text-xs"
                      onClick={() => {/* Handle location view */}}
                    >
                      View map
                    </Button>
                  )}
                </motion.div>
              )}

              {content.weather && (
                <motion.div
                  variants={performanceMode ? {} : activityVariants}
                  initial={performanceMode ? {} : "initial"}
                  animate={performanceMode ? {} : "animate"}
                  transition={{ delay: 0.35 }}
                  className="flex items-center gap-2 text-muted-foreground"
                >
                  {getWeatherIcon(content.weather.condition)}
                  <span>{content.weather.temperature}°{content.weather.unit}</span>
                  <span className="capitalize">{content.weather.condition}</span>
                </motion.div>
              )}

              {content.mood && (
                <motion.div
                  variants={performanceMode ? {} : activityVariants}
                  initial={performanceMode ? {} : "initial"}
                  animate={performanceMode ? {} : "animate"}
                  transition={{ delay: 0.4 }}
                  className="flex items-center gap-2 text-muted-foreground"
                >
                  <Circle className={cn(
                    "w-3 h-3 rounded-full",
                    content.mood === 'happy' && "bg-green-500",
                    content.mood === 'sad' && "bg-blue-500",
                    content.mood === 'excited' && "bg-yellow-500",
                    content.mood === 'tired' && "bg-purple-500",
                    content.mood === 'relaxed' && "bg-cyan-500"
                  )} />
                  <span className="capitalize">{content.mood}</span>
                </motion.div>
              )}
            </div>

            {/* Reactions and engagement */}
            {showReactions && (content.reactions && content.reactions.length > 0 || !readonly) && (
              <>
                <Separator />

                <div className="space-y-3">
                  {/* Reaction summary */}
                  {content.reactions && content.reactions.length > 0 && (
                    <motion.div
                      initial={performanceMode ? {} : { opacity: 0 }}
                      animate={performanceMode ? {} : { opacity: 1 }}
                      transition={{ delay: 0.4 }}
                      className="flex items-center justify-between"
                    >
                      <div className="flex items-center gap-2">
                        <AnimatePresence mode="popLayout">
                          {visibleReactions.map(([reaction, count]) => (
                            <motion.div
                              key={reaction}
                              variants={performanceMode ? {} : reactionVariants}
                              initial={performanceMode ? {} : "initial"}
                              animate={performanceMode ? {} : "animate"}
                              exit={performanceMode ? {} : "exit"}
                              className="flex items-center gap-1 bg-muted px-2 py-1 rounded-full text-xs"
                            >
                              {getReactionIcon(reaction as StatusReaction)}
                              <span>{count}</span>
                            </motion.div>
                          ))}
                        </AnimatePresence>

                        {Object.keys(reactionCounts).length > 3 && (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setShowAllReactions(!showAllReactions)}
                            className="h-6 px-2 text-xs"
                          >
                            {showAllReactions ? 'Show less' : `+${Object.keys(reactionCounts).length - 3} more`}
                          </Button>
                        )}
                      </div>

                      {content.commentCount && content.commentCount > 0 && (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={handleComment}
                          disabled={readonly}
                          className="text-xs text-muted-foreground hover:text-foreground"
                        >
                          <MessageCircle className="w-3 h-3 mr-1" />
                          {content.commentCount} comment{content.commentCount !== 1 ? 's' : ''}
                        </Button>
                      )}
                    </motion.div>
                  )}

                  {/* Quick reactions */}
                  {!readonly && (
                    <motion.div
                      initial={performanceMode ? {} : { opacity: 0, y: 10 }}
                      animate={performanceMode ? {} : { opacity: 1, y: 0 }}
                      transition={{ delay: 0.5 }}
                      className="flex items-center gap-2"
                    >
                      <Button
                        variant={isLiked ? "default" : "outline"}
                        size="sm"
                        onClick={handleLike}
                        className="gap-2"
                      >
                        <ThumbsUp className="w-4 h-4" />
                        Like
                      </Button>

                      <Button
                        variant="outline"
                        size="sm"
                        onClick={handleComment}
                        className="gap-2"
                      >
                        <MessageCircle className="w-4 h-4" />
                        Comment
                      </Button>

                      {/* Reaction picker trigger */}
                      <div className="flex items-center">
                        {[StatusReaction.LOVE, StatusReaction.LAUGH, StatusReaction.WOW].map((reaction) => (
                          <Tooltip key={reaction}>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleReaction(reaction)}
                                className={cn(
                                  "h-8 w-8 p-0",
                                  userReaction === reaction && "bg-primary/10 text-primary"
                                )}
                              >
                                {getReactionIcon(reaction)}
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent className="capitalize">
                              {reaction.replace('_', ' ')}
                            </TooltipContent>
                          </Tooltip>
                        ))}
                      </div>
                    </motion.div>
                  )}
                </div>
              </>
            )}
          </CardContent>
        </Card>

        {/* Performance Debug Info */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            {content.statusType} | {content.activity || 'no activity'} | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedStatusUpdateMessage = React.memo(StatusUpdateMessage, (prevProps, nextProps) => {
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.showReactions === nextProps.showReactions &&
    prevProps.showLocation === nextProps.showLocation &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedStatusUpdateMessage.displayName = 'MemoizedStatusUpdateMessage';

export default StatusUpdateMessage;