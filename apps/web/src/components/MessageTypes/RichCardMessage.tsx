// T047 - RichCardMessage component with interactive cards
/**
 * RichCardMessage Component
 * Displays rich cards with images, actions, and interactive elements
 * Supports carousel layouts, quick actions, and card expansion
 */

import React, { useState, useCallback, useMemo, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { RichCardContent, Card as CardData, CardAction, CardActionType } from '../../types/RichCardContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Separator } from '../ui/separator';
import {
  ChevronLeft,
  ChevronRight,
  ExternalLink,
  Download,
  Share,
  Phone,
  Mail,
  MapPin,
  Calendar,
  Clock,
  Star,
  Heart,
  ThumbsUp,
  MessageCircle,
  Play,
  Pause,
  Volume2,
  VolumeX,
  Maximize,
  MoreHorizontal,
  Eye,
  Bookmark
} from 'lucide-react';

// Component Props Interface
interface RichCardMessageProps {
  message: MessageData & { content: RichCardContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onCardAction?: (cardId: string, action: CardAction) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  autoPlayMedia?: boolean;
  performanceMode?: boolean;
}

// Animation Variants
const cardVariants = {
  initial: { opacity: 0, scale: 0.95, y: 20 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.95, y: -20 }
};

const carouselVariants = {
  initial: { opacity: 0, x: 50 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: -50 }
};

const actionButtonVariants = {
  initial: { scale: 1 },
  hover: { scale: 1.05 },
  tap: { scale: 0.95 }
};

const mediaVariants = {
  initial: { opacity: 0, scale: 1.1 },
  animate: { opacity: 1, scale: 1 },
  exit: { opacity: 0, scale: 0.9 }
};

export const RichCardMessage: React.FC<RichCardMessageProps> = ({
  message,
  onInteraction,
  onCardAction,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  autoPlayMedia = false,
  performanceMode = false
}) => {
  const cardRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Carousel state
  const [currentCardIndex, setCurrentCardIndex] = useState(0);
  const [expandedCard, setExpandedCard] = useState<string | null>(null);

  // Media state
  const [mediaStates, setMediaStates] = useState<Record<string, {
    playing?: boolean;
    muted?: boolean;
    currentTime?: number;
  }>>({});

  // Get current card
  const currentCard = useMemo(() => {
    return content.cards[currentCardIndex] || content.cards[0];
  }, [content.cards, currentCardIndex]);

  // Handle carousel navigation
  const handlePrevCard = useCallback(() => {
    setCurrentCardIndex((prev) =>
      prev > 0 ? prev - 1 : content.cards.length - 1
    );
  }, [content.cards.length]);

  const handleNextCard = useCallback(() => {
    setCurrentCardIndex((prev) =>
      prev < content.cards.length - 1 ? prev + 1 : 0
    );
  }, [content.cards.length]);

  // Handle card action
  const handleCardAction = useCallback((card: CardData, action: CardAction) => {
    if (readonly) return;

    switch (action.type) {
      case CardActionType.LINK:
        if (action.url) {
          window.open(action.url, action.openInNewTab ? '_blank' : '_self');
        }
        break;

      case CardActionType.PHONE:
        if (action.phoneNumber) {
          window.location.href = `tel:${action.phoneNumber}`;
        }
        break;

      case CardActionType.EMAIL:
        if (action.email) {
          window.location.href = `mailto:${action.email}`;
        }
        break;

      case CardActionType.SHARE:
        if (navigator.share && action.shareData) {
          navigator.share(action.shareData).catch(console.error);
        }
        break;

      case CardActionType.DOWNLOAD:
        if (action.downloadUrl) {
          const link = document.createElement('a');
          link.href = action.downloadUrl;
          link.download = action.filename || 'download';
          link.click();
        }
        break;

      case CardActionType.CUSTOM:
        if (onCardAction) {
          onCardAction(card.id, action);
        }
        break;
    }

    // Send interaction event
    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'card_action',
        data: {
          cardId: card.id,
          actionType: action.type,
          actionLabel: action.label
        },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, onCardAction, onInteraction, message.id]);

  // Handle card expansion
  const handleCardExpansion = useCallback((cardId: string) => {
    setExpandedCard(expandedCard === cardId ? null : cardId);
  }, [expandedCard]);

  // Handle media controls
  const handleMediaControl = useCallback((cardId: string, control: 'play' | 'pause' | 'mute' | 'unmute') => {
    setMediaStates(prev => ({
      ...prev,
      [cardId]: {
        ...prev[cardId],
        playing: control === 'play' ? true : control === 'pause' ? false : prev[cardId]?.playing,
        muted: control === 'mute' ? true : control === 'unmute' ? false : prev[cardId]?.muted
      }
    }));
  }, []);

  // Get action icon
  const getActionIcon = useCallback((type: CardActionType) => {
    switch (type) {
      case CardActionType.LINK: return <ExternalLink className="w-4 h-4" />;
      case CardActionType.DOWNLOAD: return <Download className="w-4 h-4" />;
      case CardActionType.SHARE: return <Share className="w-4 h-4" />;
      case CardActionType.PHONE: return <Phone className="w-4 h-4" />;
      case CardActionType.EMAIL: return <Mail className="w-4 h-4" />;
      case CardActionType.LOCATION: return <MapPin className="w-4 h-4" />;
      case CardActionType.CALENDAR: return <Calendar className="w-4 h-4" />;
      case CardActionType.LIKE: return <ThumbsUp className="w-4 h-4" />;
      case CardActionType.BOOKMARK: return <Bookmark className="w-4 h-4" />;
      default: return <MoreHorizontal className="w-4 h-4" />;
    }
  }, []);

  // Render card media
  const renderCardMedia = useCallback((card: CardData) => {
    if (!card.media) return null;

    const mediaState = mediaStates[card.id] || {};

    switch (card.media.type) {
      case 'image':
        return (
          <motion.div
            variants={performanceMode ? {} : mediaVariants}
            initial={performanceMode ? {} : "initial"}
            animate={performanceMode ? {} : "animate"}
            className="relative aspect-video rounded-lg overflow-hidden bg-muted"
          >
            <img
              src={card.media.url}
              alt={card.media.alt || card.title}
              className="w-full h-full object-cover cursor-pointer hover:scale-105 transition-transform duration-200"
              onClick={() => handleCardExpansion(card.id)}
            />
            {card.media.overlayText && (
              <div className="absolute inset-0 bg-black/40 flex items-center justify-center">
                <p className="text-white font-semibold text-center px-4">
                  {card.media.overlayText}
                </p>
              </div>
            )}
          </motion.div>
        );

      case 'video':
        return (
          <div className="relative aspect-video rounded-lg overflow-hidden bg-black">
            <video
              src={card.media.url}
              poster={card.media.thumbnailUrl}
              className="w-full h-full object-cover"
              muted={mediaState.muted}
              autoPlay={autoPlayMedia}
              controls={false}
            />

            {/* Custom video controls */}
            <div className="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent opacity-0 hover:opacity-100 transition-opacity">
              <div className="absolute bottom-4 left-4 right-4 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleMediaControl(card.id, mediaState.playing ? 'pause' : 'play')}
                    className="text-white hover:text-white hover:bg-white/20"
                  >
                    {mediaState.playing ?
                      <Pause className="w-4 h-4" /> :
                      <Play className="w-4 h-4" />
                    }
                  </Button>

                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleMediaControl(card.id, mediaState.muted ? 'unmute' : 'mute')}
                    className="text-white hover:text-white hover:bg-white/20"
                  >
                    {mediaState.muted ?
                      <VolumeX className="w-4 h-4" /> :
                      <Volume2 className="w-4 h-4" />
                    }
                  </Button>
                </div>

                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleCardExpansion(card.id)}
                  className="text-white hover:text-white hover:bg-white/20"
                >
                  <Maximize className="w-4 h-4" />
                </Button>
              </div>
            </div>
          </div>
        );

      default:
        return null;
    }
  }, [mediaStates, autoPlayMedia, handleCardExpansion, handleMediaControl, performanceMode]);

  // Render individual card
  const renderCard = useCallback((card: CardData, index: number) => {
    const isExpanded = expandedCard === card.id;

    return (
      <motion.div
        key={card.id}
        variants={performanceMode ? {} : carouselVariants}
        initial={performanceMode ? {} : "initial"}
        animate={performanceMode ? {} : "animate"}
        exit={performanceMode ? {} : "exit"}
        transition={{ delay: index * 0.1 }}
        className={cn(
          "card-item",
          content.layout === 'carousel' && "min-w-full",
          content.layout === 'grid' && "flex-1 min-w-0"
        )}
      >
        <Card className={cn(
          "h-full",
          isExpanded && "ring-2 ring-primary/50"
        )}>
          {/* Card media */}
          {renderCardMedia(card)}

          {/* Card content */}
          <CardContent className="p-4 space-y-3">
            {/* Header with title and subtitle */}
            <div className="space-y-1">
              <h3 className="font-semibold text-foreground leading-tight">
                {card.title}
              </h3>

              {card.subtitle && (
                <p className="text-sm text-muted-foreground">
                  {card.subtitle}
                </p>
              )}
            </div>

            {/* Description */}
            {card.description && (
              <p className={cn(
                "text-sm text-muted-foreground",
                isExpanded ? "line-clamp-none" : "line-clamp-3"
              )}>
                {card.description}
              </p>
            )}

            {/* Metadata */}
            {card.metadata && (
              <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
                {card.metadata.author && (
                  <div className="flex items-center gap-1">
                    <Avatar className="w-4 h-4">
                      <AvatarImage src={`/avatars/${card.metadata.author.toLowerCase()}.png`} />
                      <AvatarFallback className="text-xs">
                        {card.metadata.author.substring(0, 1).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    <span>{card.metadata.author}</span>
                  </div>
                )}

                {card.metadata.publishedAt && (
                  <div className="flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    <span>{new Date(card.metadata.publishedAt).toLocaleDateString()}</span>
                  </div>
                )}

                {card.metadata.rating && (
                  <div className="flex items-center gap-1">
                    <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
                    <span>{card.metadata.rating}/5</span>
                  </div>
                )}

                {card.metadata.price && (
                  <Badge variant="secondary" className="text-xs">
                    {card.metadata.price}
                  </Badge>
                )}

                {card.metadata.category && (
                  <Badge variant="outline" className="text-xs">
                    {card.metadata.category}
                  </Badge>
                )}
              </div>
            )}

            {/* Actions */}
            {card.actions && card.actions.length > 0 && (
              <>
                <Separator />
                <div className="flex flex-wrap gap-2">
                  {card.actions.map((action, actionIndex) => (
                    <motion.div
                      key={actionIndex}
                      variants={performanceMode ? {} : actionButtonVariants}
                      whileHover={performanceMode ? {} : "hover"}
                      whileTap={performanceMode ? {} : "tap"}
                    >
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button
                            variant={action.primary ? "default" : "outline"}
                            size="sm"
                            onClick={() => handleCardAction(card, action)}
                            disabled={readonly}
                            className="gap-2"
                          >
                            {getActionIcon(action.type)}
                            {action.label}
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                          {action.tooltip || action.label}
                        </TooltipContent>
                      </Tooltip>
                    </motion.div>
                  ))}
                </div>
              </>
            )}

            {/* Expansion toggle */}
            {card.description && card.description.length > 150 && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleCardExpansion(card.id)}
                className="w-full gap-2 text-xs"
              >
                <Eye className="w-3 h-3" />
                {isExpanded ? 'Show less' : 'Show more'}
              </Button>
            )}
          </CardContent>
        </Card>
      </motion.div>
    );
  }, [
    expandedCard,
    content.layout,
    renderCardMedia,
    handleCardAction,
    readonly,
    getActionIcon,
    handleCardExpansion,
    performanceMode
  ]);

  // Performance optimization
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: cardVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={cardRef}
        className={cn(
          "rich-card-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`rich-card-message-${message.id}`}
        data-card-count={content.cards.length}
        role="article"
        aria-label={`Rich card: ${content.title || 'Interactive cards'}`}
      >
        <Card className="rich-card-container">
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
                      Shared {content.cards.length} card{content.cards.length !== 1 ? 's' : ''} • {message.timestamp.toLocaleDateString()}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-2">
                <Badge variant="outline" className="text-xs">
                  {content.layout}
                </Badge>

                {content.cards.length > 1 && content.layout === 'carousel' && (
                  <Badge variant="secondary" className="text-xs">
                    {currentCardIndex + 1}/{content.cards.length}
                  </Badge>
                )}
              </div>
            </div>

            {/* Title and description */}
            {(content.title || content.description) && (
              <div className="space-y-1">
                {content.title && (
                  <h2 className="font-semibold text-foreground">
                    {content.title}
                  </h2>
                )}

                {content.description && (
                  <p className="text-sm text-muted-foreground">
                    {content.description}
                  </p>
                )}
              </div>
            )}
          </CardHeader>

          {/* Cards container */}
          <CardContent className="space-y-4">
            {content.layout === 'carousel' ? (
              <>
                {/* Carousel view */}
                <div className="relative">
                  <AnimatePresence mode="wait">
                    {renderCard(currentCard, currentCardIndex)}
                  </AnimatePresence>

                  {/* Carousel navigation */}
                  {content.cards.length > 1 && (
                    <>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={handlePrevCard}
                        className="absolute left-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 bg-background/80 backdrop-blur-sm"
                        disabled={readonly}
                      >
                        <ChevronLeft className="w-4 h-4" />
                      </Button>

                      <Button
                        variant="outline"
                        size="sm"
                        onClick={handleNextCard}
                        className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 bg-background/80 backdrop-blur-sm"
                        disabled={readonly}
                      >
                        <ChevronRight className="w-4 h-4" />
                      </Button>

                      {/* Carousel indicators */}
                      <div className="flex justify-center gap-1 mt-4">
                        {content.cards.map((_, index) => (
                          <button
                            key={index}
                            onClick={() => setCurrentCardIndex(index)}
                            className={cn(
                              "w-2 h-2 rounded-full transition-colors",
                              index === currentCardIndex
                                ? "bg-primary"
                                : "bg-muted-foreground/30 hover:bg-muted-foreground/50"
                            )}
                            aria-label={`Go to card ${index + 1}`}
                          />
                        ))}
                      </div>
                    </>
                  )}
                </div>
              </>
            ) : (
              /* Grid view */
              <div className={cn(
                "grid gap-4",
                content.cards.length === 1 && "grid-cols-1",
                content.cards.length === 2 && "grid-cols-1 md:grid-cols-2",
                content.cards.length >= 3 && "grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
              )}>
                {content.cards.map((card, index) => renderCard(card, index))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Performance Debug Info */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            {content.layout}: {content.cards.length} cards | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedRichCardMessage = React.memo(RichCardMessage, (prevProps, nextProps) => {
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.autoPlayMedia === nextProps.autoPlayMedia &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedRichCardMessage.displayName = 'MemoizedRichCardMessage';

export default RichCardMessage;