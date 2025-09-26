// T051 - EmbedMessage component with rich embeds
/**
 * EmbedMessage Component
 * Displays rich embeds from external sources with link previews, media, and metadata
 * Supports various embed types including websites, social media, and media content
 */

import React, { useState, useCallback, useMemo, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { EmbedContent, EmbedType, EmbedProvider } from '../../types/EmbedContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Separator } from '../ui/separator';
import {
  ExternalLink,
  Play,
  Pause,
  Volume2,
  VolumeX,
  Maximize,
  Share,
  Bookmark,
  Eye,
  Heart,
  MessageSquare,
  Repeat2,
  Download,
  AlertTriangle,
  Globe,
  Calendar,
  Clock,
  User,
  Hash,
  Star,
  TrendingUp,
  Music,
  Video,
  Image,
  FileText,
  Code,
  MoreHorizontal,
  RefreshCw,
  ChevronUp,
  ChevronDown
} from 'lucide-react';

// Component Props Interface
interface EmbedMessageProps {
  message: MessageData & { content: EmbedContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onVisitLink?: (url: string) => void;
  onShare?: (embedId: string, url: string) => void;
  onBookmark?: (embedId: string) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  autoPlayMedia?: boolean;
  showMetadata?: boolean;
  performanceMode?: boolean;
}

// Animation Variants
const embedVariants = {
  initial: { opacity: 0, scale: 0.98, y: 20 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.98, y: -20 }
};

const mediaVariants = {
  initial: { opacity: 0, scale: 1.05 },
  animate: { opacity: 1, scale: 1 },
  exit: { opacity: 0, scale: 0.95 }
};

const contentVariants = {
  initial: { opacity: 0, y: 10 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -10 }
};

const metadataVariants = {
  initial: { opacity: 0, height: 0 },
  animate: { opacity: 1, height: 'auto' },
  exit: { opacity: 0, height: 0 }
};

export const EmbedMessage: React.FC<EmbedMessageProps> = ({
  message,
  onInteraction,
  onVisitLink,
  onShare,
  onBookmark,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  autoPlayMedia = false,
  showMetadata = true,
  performanceMode = false
}) => {
  const embedRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Embed state
  const [isBookmarked, setIsBookmarked] = useState(false);
  const [showFullDescription, setShowFullDescription] = useState(false);
  const [mediaError, setMediaError] = useState(false);
  const [isMediaLoading, setIsMediaLoading] = useState(true);
  const [showAdvancedMetadata, setShowAdvancedMetadata] = useState(false);

  // Media controls state
  const [isPlaying, setIsPlaying] = useState(autoPlayMedia);
  const [isMuted, setIsMuted] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);

  // Get provider styling
  const getProviderStyling = useCallback((provider: EmbedProvider) => {
    switch (provider) {
      case EmbedProvider.YOUTUBE:
        return { color: 'bg-red-500', icon: Video };
      case EmbedProvider.TWITTER:
        return { color: 'bg-blue-500', icon: MessageSquare };
      case EmbedProvider.GITHUB:
        return { color: 'bg-gray-800', icon: Code };
      case EmbedProvider.SPOTIFY:
        return { color: 'bg-green-500', icon: Music };
      case EmbedProvider.INSTAGRAM:
        return { color: 'bg-gradient-to-r from-purple-500 to-pink-500', icon: Image };
      case EmbedProvider.LINKEDIN:
        return { color: 'bg-blue-600', icon: User };
      case EmbedProvider.FIGMA:
        return { color: 'bg-purple-500', icon: FileText };
      default:
        return { color: 'bg-muted', icon: Globe };
    }
  }, []);

  // Get embed type icon
  const getEmbedTypeIcon = useCallback((type: EmbedType) => {
    switch (type) {
      case EmbedType.VIDEO: return Video;
      case EmbedType.AUDIO: return Music;
      case EmbedType.IMAGE: return Image;
      case EmbedType.DOCUMENT: return FileText;
      case EmbedType.SOCIAL: return MessageSquare;
      case EmbedType.CODE: return Code;
      case EmbedType.ARTICLE: return FileText;
      default: return Globe;
    }
  }, []);

  // Handle link visit
  const handleVisitLink = useCallback(() => {
    if (readonly) return;

    if (onVisitLink) {
      onVisitLink(content.url);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'embed_visit',
        data: { embedId: content.id, url: content.url },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, onVisitLink, content.url, content.id, onInteraction, message.id]);

  // Handle share
  const handleShare = useCallback(() => {
    if (readonly) return;

    if (onShare) {
      onShare(content.id, content.url);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'embed_share',
        data: { embedId: content.id, url: content.url },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, onShare, content.id, content.url, onInteraction, message.id]);

  // Handle bookmark
  const handleBookmark = useCallback(() => {
    if (readonly) return;

    setIsBookmarked(!isBookmarked);

    if (onBookmark) {
      onBookmark(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: isBookmarked ? 'remove_bookmark' : 'add_bookmark',
        data: { embedId: content.id, url: content.url },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, isBookmarked, onBookmark, content.id, content.url, onInteraction, message.id]);

  // Handle media controls
  const handlePlayPause = useCallback(() => {
    setIsPlaying(!isPlaying);
  }, [isPlaying]);

  const handleMuteToggle = useCallback(() => {
    setIsMuted(!isMuted);
  }, [isMuted]);

  const handleFullscreen = useCallback(() => {
    setIsFullscreen(!isFullscreen);
  }, [isFullscreen]);

  // Format file size
  const formatFileSize = useCallback((bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }, []);

  // Format duration
  const formatDuration = useCallback((seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);

    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    }
    return `${minutes}:${secs.toString().padStart(2, '0')}`;
  }, []);

  // Render embed media
  const renderEmbedMedia = useCallback(() => {
    if (!content.media) return null;

    const handleMediaLoad = () => setIsMediaLoading(false);
    const handleMediaError = () => {
      setIsMediaLoading(false);
      setMediaError(true);
    };

    return (
      <motion.div
        variants={performanceMode ? {} : mediaVariants}
        initial={performanceMode ? {} : "initial"}
        animate={performanceMode ? {} : "animate"}
        className="relative aspect-video bg-muted rounded-lg overflow-hidden"
      >
        {mediaError ? (
          <div className="w-full h-full flex flex-col items-center justify-center text-muted-foreground">
            <AlertTriangle className="w-8 h-8 mb-2" />
            <p className="text-sm">Failed to load media</p>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                setMediaError(false);
                setIsMediaLoading(true);
              }}
              className="mt-2"
            >
              <RefreshCw className="w-3 h-3 mr-1" />
              Retry
            </Button>
          </div>
        ) : (
          <>
            {content.type === EmbedType.VIDEO ? (
              <div className="relative w-full h-full">
                <video
                  src={content.media.url}
                  poster={content.media.thumbnailUrl}
                  className="w-full h-full object-cover"
                  autoPlay={autoPlayMedia}
                  muted={isMuted}
                  onLoadedData={handleMediaLoad}
                  onError={handleMediaError}
                  controls={false}
                />

                {/* Custom video controls */}
                <div className="absolute inset-0 bg-black/20 opacity-0 hover:opacity-100 transition-opacity duration-200">
                  <div className="absolute bottom-4 left-4 right-4 flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handlePlayPause}
                        className="text-white hover:text-white hover:bg-white/20"
                      >
                        {isPlaying ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
                      </Button>

                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleMuteToggle}
                        className="text-white hover:text-white hover:bg-white/20"
                      >
                        {isMuted ? <VolumeX className="w-4 h-4" /> : <Volume2 className="w-4 h-4" />}
                      </Button>

                      {content.duration && (
                        <span className="text-white text-xs">
                          {formatDuration(content.duration)}
                        </span>
                      )}
                    </div>

                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleFullscreen}
                      className="text-white hover:text-white hover:bg-white/20"
                    >
                      <Maximize className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </div>
            ) : content.type === EmbedType.AUDIO ? (
              <div className="w-full h-32 bg-gradient-to-r from-purple-500 to-blue-500 flex items-center justify-center">
                <div className="text-center text-white">
                  <Music className="w-12 h-12 mx-auto mb-2" />
                  <p className="text-sm font-medium">{content.title}</p>
                  {content.duration && (
                    <p className="text-xs opacity-75">
                      {formatDuration(content.duration)}
                    </p>
                  )}
                </div>
              </div>
            ) : (
              <img
                src={content.media.url}
                alt={content.media.alt || content.title}
                className="w-full h-full object-cover cursor-pointer hover:scale-105 transition-transform duration-200"
                onLoad={handleMediaLoad}
                onError={handleMediaError}
                onClick={handleVisitLink}
              />
            )}

            {/* Loading overlay */}
            {isMediaLoading && (
              <div className="absolute inset-0 bg-muted flex items-center justify-center">
                <RefreshCw className="w-6 h-6 animate-spin text-muted-foreground" />
              </div>
            )}
          </>
        )}
      </motion.div>
    );
  }, [
    content.media,
    content.type,
    content.title,
    content.duration,
    mediaError,
    isMediaLoading,
    autoPlayMedia,
    isMuted,
    isPlaying,
    handlePlayPause,
    handleMuteToggle,
    handleFullscreen,
    handleVisitLink,
    formatDuration,
    performanceMode
  ]);

  // Get provider info
  const providerStyling = getProviderStyling(content.provider);
  const EmbedTypeIcon = getEmbedTypeIcon(content.type);

  // Performance optimization
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: embedVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={embedRef}
        className={cn(
          "embed-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`embed-message-${message.id}`}
        data-embed-type={content.type}
        data-provider={content.provider}
        role="article"
        aria-label={`Embed: ${content.title} from ${content.provider}`}
      >
        <Card className="embed-card overflow-hidden">
          {/* Header */}
          <CardHeader className="space-y-3">
            <div className="flex items-start justify-between gap-3">
              <div className="flex items-center gap-3 min-w-0 flex-1">
                {showAvatar && (
                  <motion.div
                    initial={performanceMode ? {} : { scale: 0.8, opacity: 0 }}
                    animate={performanceMode ? {} : { scale: 1, opacity: 1 }}
                    transition={{ delay: 0.1 }}
                  >
                    <Avatar className={cn(compactMode ? "w-8 h-8" : "w-10 h-10")}>
                      <AvatarImage src={`/avatars/${message.senderName.toLowerCase()}.png`} />
                      <AvatarFallback>
                        {message.senderName.substring(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
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
                  </div>
                  {showTimestamp && (
                    <p className="text-xs text-muted-foreground mt-1">
                      Shared a link • {message.timestamp.toLocaleDateString()}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-2">
                {/* Provider badge */}
                <Badge variant="outline" className="text-xs">
                  <providerStyling.icon className="w-3 h-3 mr-1" />
                  {content.provider}
                </Badge>

                {/* Type badge */}
                <Badge variant="secondary" className="text-xs">
                  <EmbedTypeIcon className="w-3 h-3 mr-1" />
                  {content.type}
                </Badge>

                {/* Action buttons */}
                <div className="flex items-center gap-1">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleBookmark}
                        disabled={readonly}
                        className={cn(
                          "h-8 w-8 p-0",
                          isBookmarked && "text-yellow-500 hover:text-yellow-600"
                        )}
                      >
                        <Bookmark className={cn("w-4 h-4", isBookmarked && "fill-current")} />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      {isBookmarked ? 'Remove bookmark' : 'Bookmark'}
                    </TooltipContent>
                  </Tooltip>

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
                    <TooltipContent>Share</TooltipContent>
                  </Tooltip>
                </div>
              </div>
            </div>
          </CardHeader>

          <CardContent className="space-y-4">
            {/* Embed media */}
            {content.media && renderEmbedMedia()}

            {/* Content details */}
            <motion.div
              variants={performanceMode ? {} : contentVariants}
              initial={performanceMode ? {} : "initial"}
              animate={performanceMode ? {} : "animate"}
              transition={{ delay: 0.1 }}
              className="space-y-3"
            >
              {/* Title and URL */}
              <div className="space-y-2">
                <h3 className="font-semibold text-lg text-foreground leading-tight line-clamp-2">
                  {content.title}
                </h3>

                <Button
                  variant="ghost"
                  onClick={handleVisitLink}
                  disabled={readonly}
                  className="h-auto p-0 text-sm text-primary hover:text-primary/80 font-normal justify-start"
                >
                  <Globe className="w-3 h-3 mr-2 flex-shrink-0" />
                  <span className="truncate">{content.url}</span>
                  <ExternalLink className="w-3 h-3 ml-1 flex-shrink-0" />
                </Button>
              </div>

              {/* Description */}
              {content.description && (
                <div className="space-y-2">
                  <p className={cn(
                    "text-sm text-muted-foreground leading-relaxed",
                    !showFullDescription && "line-clamp-3"
                  )}>
                    {content.description}
                  </p>

                  {content.description.length > 200 && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setShowFullDescription(!showFullDescription)}
                      className="h-auto p-0 text-xs text-muted-foreground hover:text-foreground"
                    >
                      {showFullDescription ? (
                        <>
                          <ChevronUp className="w-3 h-3 mr-1" />
                          Show less
                        </>
                      ) : (
                        <>
                          <ChevronDown className="w-3 h-3 mr-1" />
                          Show more
                        </>
                      )}
                    </Button>
                  )}
                </div>
              )}

              {/* Basic metadata */}
              {showMetadata && (
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                  {content.author && (
                    <div className="flex items-center gap-2">
                      <User className="w-4 h-4 text-muted-foreground" />
                      <span className="truncate text-muted-foreground">{content.author}</span>
                    </div>
                  )}

                  {content.publishedAt && (
                    <div className="flex items-center gap-2">
                      <Calendar className="w-4 h-4 text-muted-foreground" />
                      <span className="text-muted-foreground">
                        {new Date(content.publishedAt).toLocaleDateString()}
                      </span>
                    </div>
                  )}

                  {content.duration && (
                    <div className="flex items-center gap-2">
                      <Clock className="w-4 h-4 text-muted-foreground" />
                      <span className="text-muted-foreground">
                        {formatDuration(content.duration)}
                      </span>
                    </div>
                  )}

                  {content.fileSize && (
                    <div className="flex items-center gap-2">
                      <Download className="w-4 h-4 text-muted-foreground" />
                      <span className="text-muted-foreground">
                        {formatFileSize(content.fileSize)}
                      </span>
                    </div>
                  )}
                </div>
              )}

              {/* Advanced metadata */}
              {showMetadata && content.metadata && Object.keys(content.metadata).length > 0 && (
                <>
                  <Separator />
                  <div className="space-y-2">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setShowAdvancedMetadata(!showAdvancedMetadata)}
                      className="h-auto p-0 text-xs font-medium text-foreground flex items-center gap-1"
                    >
                      <MoreHorizontal className="w-3 h-3" />
                      Advanced Details
                      <ChevronDown className={cn(
                        "w-3 h-3 transition-transform",
                        showAdvancedMetadata && "rotate-180"
                      )} />
                    </Button>

                    <AnimatePresence>
                      {showAdvancedMetadata && (
                        <motion.div
                          variants={performanceMode ? {} : metadataVariants}
                          initial={performanceMode ? {} : "initial"}
                          animate={performanceMode ? {} : "animate"}
                          exit={performanceMode ? {} : "exit"}
                          className="grid grid-cols-1 md:grid-cols-2 gap-2 text-xs"
                        >
                          {Object.entries(content.metadata).map(([key, value]) => (
                            <div key={key} className="flex justify-between gap-2">
                              <span className="text-muted-foreground capitalize">
                                {key.replace(/([A-Z])/g, ' $1').trim()}:
                              </span>
                              <span className="font-medium text-right truncate">
                                {String(value)}
                              </span>
                            </div>
                          ))}
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </div>
                </>
              )}

              {/* Social metrics */}
              {(content.likes || content.shares || content.comments || content.views) && (
                <>
                  <Separator />
                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center gap-4">
                      {content.views && (
                        <div className="flex items-center gap-1 text-muted-foreground">
                          <Eye className="w-4 h-4" />
                          <span>{content.views.toLocaleString()}</span>
                        </div>
                      )}

                      {content.likes && (
                        <div className="flex items-center gap-1 text-muted-foreground">
                          <Heart className="w-4 h-4" />
                          <span>{content.likes.toLocaleString()}</span>
                        </div>
                      )}

                      {content.comments && (
                        <div className="flex items-center gap-1 text-muted-foreground">
                          <MessageSquare className="w-4 h-4" />
                          <span>{content.comments.toLocaleString()}</span>
                        </div>
                      )}

                      {content.shares && (
                        <div className="flex items-center gap-1 text-muted-foreground">
                          <Repeat2 className="w-4 h-4" />
                          <span>{content.shares.toLocaleString()}</span>
                        </div>
                      )}
                    </div>

                    {content.rating && (
                      <div className="flex items-center gap-1">
                        <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                        <span className="font-medium">{content.rating.toFixed(1)}</span>
                      </div>
                    )}
                  </div>
                </>
              )}

              {/* Action button */}
              <div className="pt-2">
                <Button
                  onClick={handleVisitLink}
                  disabled={readonly}
                  className="w-full gap-2"
                  size="sm"
                >
                  <ExternalLink className="w-4 h-4" />
                  Visit {content.provider}
                </Button>
              </div>
            </motion.div>
          </CardContent>
        </Card>

        {/* Performance Debug Info */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            {content.type} | {content.provider} | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedEmbedMessage = React.memo(EmbedMessage, (prevProps, nextProps) => {
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.autoPlayMedia === nextProps.autoPlayMedia &&
    prevProps.showMetadata === nextProps.showMetadata &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedEmbedMessage.displayName = 'MemoizedEmbedMessage';

export default EmbedMessage;