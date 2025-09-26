// T043 - EventMessage component with RSVP and calendar integration
/**
 * EventMessage Component
 * Displays event invitations with RSVP functionality, calendar integration, and attendee management
 * Supports recurring events, reminders, and rich event details with location mapping
 */

import React, { useState, useCallback, useMemo, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { EventContent, RSVPStatus, EventType } from '../../types/EventContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Separator } from '../ui/separator';
import {
  Calendar,
  Clock,
  MapPin,
  Users,
  Plus,
  Minus,
  Check,
  X,
  AlertCircle,
  ExternalLink,
  Bell,
  RefreshCcw,
  Globe
} from 'lucide-react';

// Component Props Interface
interface EventMessageProps {
  message: MessageData & { content: EventContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onRSVP?: (eventId: string, status: RSVPStatus) => void;
  onAddToCalendar?: (event: EventContent) => void;
  onViewLocation?: (location: { lat: number; lng: number; address: string }) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  showAttendeeLimit?: number;
  performanceMode?: boolean;
}

// Animation Variants
const eventVariants = {
  initial: { opacity: 0, scale: 0.95, y: 20 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.95, y: -10 }
};

const rsvpButtonVariants = {
  initial: { scale: 1 },
  hover: { scale: 1.05 },
  tap: { scale: 0.95 },
  selected: { scale: 1.02 }
};

const attendeeVariants = {
  initial: { opacity: 0, x: -20 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 20 }
};

export const EventMessage: React.FC<EventMessageProps> = ({
  message,
  onInteraction,
  onRSVP,
  onAddToCalendar,
  onViewLocation,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  showAttendeeLimit = 5,
  performanceMode = false
}) => {
  const eventRef = useRef<HTMLDivElement>(null);
  const { content } = message;
  const [userRSVP, setUserRSVP] = useState<RSVPStatus>(
    content.attendees?.find(a => a.userId === 'current-user')?.status || RSVPStatus.NO_RESPONSE
  );
  const [showAllAttendees, setShowAllAttendees] = useState(false);

  // Calculate event timing
  const eventTiming = useMemo(() => {
    const now = new Date();
    const startTime = new Date(content.startTime);
    const endTime = new Date(content.endTime);

    const isUpcoming = startTime > now;
    const isOngoing = startTime <= now && endTime >= now;
    const isPast = endTime < now;

    const timeUntilStart = Math.max(0, startTime.getTime() - now.getTime());
    const duration = endTime.getTime() - startTime.getTime();

    return {
      isUpcoming,
      isOngoing,
      isPast,
      timeUntilStart,
      duration,
      status: isPast ? 'past' : isOngoing ? 'ongoing' : 'upcoming'
    };
  }, [content.startTime, content.endTime]);

  // Format time display
  const formatEventTime = useCallback((date: Date) => {
    const now = new Date();
    const eventDate = new Date(date);
    const diffInDays = Math.floor((eventDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));

    const timeString = eventDate.toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true
    });

    if (diffInDays === 0) return `Today at ${timeString}`;
    if (diffInDays === 1) return `Tomorrow at ${timeString}`;
    if (diffInDays === -1) return `Yesterday at ${timeString}`;
    if (diffInDays > 1 && diffInDays <= 7) {
      return `${eventDate.toLocaleDateString([], { weekday: 'long' })} at ${timeString}`;
    }

    return `${eventDate.toLocaleDateString([], {
      month: 'short',
      day: 'numeric',
      year: eventDate.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
    })} at ${timeString}`;
  }, []);

  // Format duration
  const formatDuration = useCallback((milliseconds: number) => {
    const hours = Math.floor(milliseconds / (1000 * 60 * 60));
    const minutes = Math.floor((milliseconds % (1000 * 60 * 60)) / (1000 * 60));

    if (hours === 0) return `${minutes} min`;
    if (minutes === 0) return `${hours}h`;
    return `${hours}h ${minutes}m`;
  }, []);

  // Handle RSVP action
  const handleRSVP = useCallback((status: RSVPStatus) => {
    if (readonly) return;

    setUserRSVP(status);

    if (onRSVP) {
      onRSVP(content.id, status);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'event_rsvp',
        data: { eventId: content.id, status },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [content.id, message.id, readonly, onRSVP, onInteraction]);

  // Handle add to calendar
  const handleAddToCalendar = useCallback(() => {
    if (onAddToCalendar) {
      onAddToCalendar(content);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'calendar_add',
        data: { eventId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [content, message.id, onAddToCalendar, onInteraction]);

  // Handle location view
  const handleViewLocation = useCallback(() => {
    if (content.location && onViewLocation) {
      onViewLocation({
        lat: content.location.latitude || 0,
        lng: content.location.longitude || 0,
        address: content.location.address
      });
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'location_view',
        data: { eventId: content.id, location: content.location },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [content.location, content.id, message.id, onViewLocation, onInteraction]);

  // Get attendee counts
  const attendeeCounts = useMemo(() => {
    if (!content.attendees) return { going: 0, maybe: 0, notGoing: 0, noResponse: 0 };

    return content.attendees.reduce((counts, attendee) => {
      counts[attendee.status]++;
      return counts;
    }, { going: 0, maybe: 0, notGoing: 0, noResponse: 0 });
  }, [content.attendees]);

  // Get visible attendees
  const visibleAttendees = useMemo(() => {
    if (!content.attendees) return [];

    const goingAttendees = content.attendees
      .filter(a => a.status === RSVPStatus.GOING)
      .sort((a, b) => new Date(b.respondedAt).getTime() - new Date(a.respondedAt).getTime());

    return showAllAttendees ? goingAttendees : goingAttendees.slice(0, showAttendeeLimit);
  }, [content.attendees, showAllAttendees, showAttendeeLimit]);

  // Get event type styling
  const getEventTypeColor = useCallback((type: EventType) => {
    switch (type) {
      case EventType.MEETING: return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case EventType.SOCIAL: return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      case EventType.CONFERENCE: return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200';
      case EventType.WORKSHOP: return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200';
      case EventType.WEBINAR: return 'bg-teal-100 text-teal-800 dark:bg-teal-900 dark:text-teal-200';
      default: return 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200';
    }
  }, []);

  // Get status color
  const getStatusColor = useCallback((status: string) => {
    switch (status) {
      case 'ongoing': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      case 'upcoming': return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case 'past': return 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200';
      default: return 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200';
    }
  }, []);

  // Performance optimization: skip animation in performance mode
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: eventVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={eventRef}
        className={cn(
          "event-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`event-message-${message.id}`}
        data-event-status={eventTiming.status}
        role="article"
        aria-label={`Event: ${content.title} - ${formatEventTime(new Date(content.startTime))}`}
      >
        <Card className={cn(
          "event-card border-2",
          eventTiming.isOngoing && "border-green-200 bg-green-50/50 dark:border-green-800 dark:bg-green-950/50",
          eventTiming.isPast && "opacity-75",
          compactMode ? "p-3" : "p-4"
        )}>
          <CardHeader className={cn("space-y-3", compactMode && "space-y-2")}>
            {/* Header with sender info and event status */}
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
                      {formatEventTime(message.timestamp)} • Event invitation
                    </p>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-2">
                <Badge className={getEventTypeColor(content.type)}>
                  {content.type}
                </Badge>
                <Badge className={getStatusColor(eventTiming.status)}>
                  {eventTiming.status}
                </Badge>
              </div>
            </div>

            {/* Event title and description */}
            <div className="space-y-2">
              <CardTitle className={cn(
                "text-xl font-bold text-foreground leading-tight",
                compactMode && "text-lg"
              )}>
                {content.title}
              </CardTitle>
              {content.description && (
                <p className="text-muted-foreground text-sm leading-relaxed">
                  {content.description}
                </p>
              )}
            </div>
          </CardHeader>

          <CardContent className={cn("space-y-4", compactMode && "space-y-3")}>
            {/* Event details */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Time and duration */}
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <Calendar className="w-4 h-4 text-primary" />
                  <div>
                    <p className="font-medium text-sm">
                      {formatEventTime(new Date(content.startTime))}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      Duration: {formatDuration(eventTiming.duration)}
                    </p>
                  </div>
                </div>

                {content.isRecurring && (
                  <div className="flex items-center gap-2">
                    <RefreshCcw className="w-4 h-4 text-primary" />
                    <span className="text-sm text-muted-foreground">
                      Recurring event
                    </span>
                  </div>
                )}

                {eventTiming.isUpcoming && eventTiming.timeUntilStart < 24 * 60 * 60 * 1000 && (
                  <div className="flex items-center gap-2">
                    <Bell className="w-4 h-4 text-orange-500" />
                    <span className="text-sm text-orange-600 dark:text-orange-400">
                      Starts in {formatDuration(eventTiming.timeUntilStart)}
                    </span>
                  </div>
                )}
              </div>

              {/* Location and attendees */}
              <div className="space-y-3">
                {content.location && (
                  <div className="flex items-start gap-2">
                    <MapPin className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
                    <div className="min-w-0 flex-1">
                      <p className="text-sm font-medium">{content.location.venue || 'Location'}</p>
                      <p className="text-xs text-muted-foreground line-clamp-2">
                        {content.location.address}
                      </p>
                      {content.location.latitude && content.location.longitude && (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={handleViewLocation}
                          className="h-6 px-2 text-xs mt-1"
                        >
                          <ExternalLink className="w-3 h-3 mr-1" />
                          View on map
                        </Button>
                      )}
                    </div>
                  </div>
                )}

                {content.isOnline && (
                  <div className="flex items-center gap-2">
                    <Globe className="w-4 h-4 text-primary" />
                    <span className="text-sm text-muted-foreground">
                      Online event
                    </span>
                  </div>
                )}

                <div className="flex items-center gap-2">
                  <Users className="w-4 h-4 text-primary" />
                  <div className="text-sm">
                    <span className="font-medium">{attendeeCounts.going}</span>
                    <span className="text-muted-foreground"> going</span>
                    {attendeeCounts.maybe > 0 && (
                      <>
                        <span className="text-muted-foreground"> • </span>
                        <span className="font-medium">{attendeeCounts.maybe}</span>
                        <span className="text-muted-foreground"> maybe</span>
                      </>
                    )}
                  </div>
                </div>
              </div>
            </div>

            <Separator />

            {/* RSVP Actions */}
            {!readonly && !eventTiming.isPast && (
              <motion.div
                initial={performanceMode ? {} : { opacity: 0, y: 10 }}
                animate={performanceMode ? {} : { opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
                className="flex flex-wrap gap-2"
              >
                <motion.div
                  variants={performanceMode ? {} : rsvpButtonVariants}
                  whileHover={performanceMode ? {} : "hover"}
                  whileTap={performanceMode ? {} : "tap"}
                  animate={userRSVP === RSVPStatus.GOING ? "selected" : "initial"}
                >
                  <Button
                    onClick={() => handleRSVP(RSVPStatus.GOING)}
                    variant={userRSVP === RSVPStatus.GOING ? "default" : "outline"}
                    size="sm"
                    className="gap-2"
                  >
                    <Check className="w-4 h-4" />
                    Going
                  </Button>
                </motion.div>

                <motion.div
                  variants={performanceMode ? {} : rsvpButtonVariants}
                  whileHover={performanceMode ? {} : "hover"}
                  whileTap={performanceMode ? {} : "tap"}
                  animate={userRSVP === RSVPStatus.MAYBE ? "selected" : "initial"}
                >
                  <Button
                    onClick={() => handleRSVP(RSVPStatus.MAYBE)}
                    variant={userRSVP === RSVPStatus.MAYBE ? "default" : "outline"}
                    size="sm"
                    className="gap-2"
                  >
                    <AlertCircle className="w-4 h-4" />
                    Maybe
                  </Button>
                </motion.div>

                <motion.div
                  variants={performanceMode ? {} : rsvpButtonVariants}
                  whileHover={performanceMode ? {} : "hover"}
                  whileTap={performanceMode ? {} : "tap"}
                  animate={userRSVP === RSVPStatus.NOT_GOING ? "selected" : "initial"}
                >
                  <Button
                    onClick={() => handleRSVP(RSVPStatus.NOT_GOING)}
                    variant={userRSVP === RSVPStatus.NOT_GOING ? "default" : "outline"}
                    size="sm"
                    className="gap-2"
                  >
                    <X className="w-4 h-4" />
                    Can't go
                  </Button>
                </motion.div>

                <Separator orientation="vertical" className="h-6" />

                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      onClick={handleAddToCalendar}
                      variant="outline"
                      size="sm"
                      className="gap-2"
                    >
                      <Plus className="w-4 h-4" />
                      Add to calendar
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Add event to your calendar</TooltipContent>
                </Tooltip>
              </motion.div>
            )}

            {/* Attendee List */}
            {visibleAttendees.length > 0 && (
              <motion.div
                initial={performanceMode ? {} : { opacity: 0 }}
                animate={performanceMode ? {} : { opacity: 1 }}
                transition={{ delay: 0.3 }}
                className="space-y-3"
              >
                <Separator />

                <div className="space-y-2">
                  <h4 className="text-sm font-medium text-foreground">
                    Going ({attendeeCounts.going})
                  </h4>

                  <div className="flex flex-wrap gap-2">
                    <AnimatePresence mode="popLayout">
                      {visibleAttendees.map((attendee) => (
                        <motion.div
                          key={attendee.userId}
                          variants={performanceMode ? {} : attendeeVariants}
                          initial={performanceMode ? {} : "initial"}
                          animate={performanceMode ? {} : "animate"}
                          exit={performanceMode ? {} : "exit"}
                          transition={{ duration: 0.2 }}
                          className="flex items-center gap-2 bg-green-50 dark:bg-green-950/50 px-2 py-1 rounded-full"
                        >
                          <Avatar className="w-6 h-6">
                            <AvatarImage src={`/avatars/${attendee.userName.toLowerCase()}.png`} />
                            <AvatarFallback className="text-xs">
                              {attendee.userName.substring(0, 2).toUpperCase()}
                            </AvatarFallback>
                          </Avatar>
                          <span className="text-sm text-green-800 dark:text-green-200">
                            {attendee.userName}
                          </span>
                        </motion.div>
                      ))}
                    </AnimatePresence>
                  </div>

                  {content.attendees && content.attendees.filter(a => a.status === RSVPStatus.GOING).length > showAttendeeLimit && (
                    <Button
                      onClick={() => setShowAllAttendees(!showAllAttendees)}
                      variant="ghost"
                      size="sm"
                      className="h-8 px-2 text-xs"
                    >
                      {showAllAttendees ? (
                        <>
                          <Minus className="w-3 h-3 mr-1" />
                          Show less
                        </>
                      ) : (
                        <>
                          <Plus className="w-3 h-3 mr-1" />
                          Show {content.attendees.filter(a => a.status === RSVPStatus.GOING).length - showAttendeeLimit} more
                        </>
                      )}
                    </Button>
                  )}
                </div>
              </motion.div>
            )}
          </CardContent>
        </Card>

        {/* Performance Debug Info (Development Only) */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            Status: {eventTiming.status} | RSVP: {userRSVP} | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedEventMessage = React.memo(EventMessage, (prevProps, nextProps) => {
  // Custom comparison function for performance
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedEventMessage.displayName = 'MemoizedEventMessage';

// Export both versions
export default EventMessage;