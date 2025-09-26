// T049 - VenueMessage component with location and business info
/**
 * VenueMessage Component
 * Displays venue information with location details, reviews, and booking capabilities
 * Supports business hours, contact information, and direction services
 */

import React, { useState, useCallback, useMemo, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { VenueContent, BusinessHours, VenueFeature } from '../../types/VenueContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Separator } from '../ui/separator';
import {
  MapPin,
  Phone,
  Globe,
  Clock,
  Star,
  Navigation,
  Share,
  Bookmark,
  Calendar,
  Users,
  Car,
  Wifi,
  CreditCard,
  Accessibility,
  Coffee,
  Utensils,
  ShoppingCart,
  Music,
  Camera,
  Shield,
  Heart,
  ExternalLink,
  ChevronRight,
  ChevronLeft,
  Eye,
  AlertCircle,
  CheckCircle,
  XCircle
} from 'lucide-react';

// Component Props Interface
interface VenueMessageProps {
  message: MessageData & { content: VenueContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onGetDirections?: (venueId: string, coordinates: { lat: number; lng: number }) => void;
  onBookVenue?: (venueId: string) => void;
  onCall?: (phoneNumber: string) => void;
  onVisitWebsite?: (url: string) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  showPhotos?: boolean;
  showReviews?: boolean;
  performanceMode?: boolean;
}

// Animation Variants
const venueVariants = {
  initial: { opacity: 0, scale: 0.95, y: 20 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.95, y: -20 }
};

const photoVariants = {
  initial: { opacity: 0, scale: 1.1 },
  animate: { opacity: 1, scale: 1 },
  exit: { opacity: 0, scale: 0.9 }
};

const actionVariants = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -20 }
};

const hoursVariants = {
  initial: { opacity: 0, height: 0 },
  animate: { opacity: 1, height: 'auto' },
  exit: { opacity: 0, height: 0 }
};

export const VenueMessage: React.FC<VenueMessageProps> = ({
  message,
  onInteraction,
  onGetDirections,
  onBookVenue,
  onCall,
  onVisitWebsite,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  showPhotos = true,
  showReviews = true,
  performanceMode = false
}) => {
  const venueRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Venue state
  const [currentPhotoIndex, setCurrentPhotoIndex] = useState(0);
  const [showHours, setShowHours] = useState(false);
  const [isBookmarked, setIsBookmarked] = useState(false);

  // Calculate current status
  const venueStatus = useMemo(() => {
    if (!content.businessHours) return { isOpen: null, status: 'Unknown' };

    const now = new Date();
    const currentDay = now.getDay(); // 0 = Sunday, 1 = Monday, etc.
    const currentTime = now.getHours() * 60 + now.getMinutes();

    const todayHours = content.businessHours.find(h => h.dayOfWeek === currentDay);

    if (!todayHours) {
      return { isOpen: false, status: 'Closed today' };
    }

    if (!todayHours.isOpen) {
      return { isOpen: false, status: 'Closed today' };
    }

    const openTime = todayHours.openTime.hours * 60 + todayHours.openTime.minutes;
    const closeTime = todayHours.closeTime.hours * 60 + todayHours.closeTime.minutes;

    if (currentTime >= openTime && currentTime <= closeTime) {
      // Calculate closing time
      const timeToClose = closeTime - currentTime;
      const hoursToClose = Math.floor(timeToClose / 60);
      const minutesToClose = timeToClose % 60;

      if (timeToClose <= 30) {
        return { isOpen: true, status: `Closing soon (${minutesToClose}m)` };
      }

      return {
        isOpen: true,
        status: `Open until ${todayHours.closeTime.hours.toString().padStart(2, '0')}:${todayHours.closeTime.minutes.toString().padStart(2, '0')}`
      };
    }

    if (currentTime < openTime) {
      const timeToOpen = openTime - currentTime;
      const hoursToOpen = Math.floor(timeToOpen / 60);
      const minutesToOpen = timeToOpen % 60;

      if (hoursToOpen === 0) {
        return { isOpen: false, status: `Opens in ${minutesToOpen}m` };
      }

      return {
        isOpen: false,
        status: `Opens at ${todayHours.openTime.hours.toString().padStart(2, '0')}:${todayHours.openTime.minutes.toString().padStart(2, '0')}`
      };
    }

    return { isOpen: false, status: 'Closed' };
  }, [content.businessHours]);

  // Format rating
  const formatRating = useCallback((rating: number) => {
    return rating.toFixed(1);
  }, []);

  // Handle photo navigation
  const handlePrevPhoto = useCallback(() => {
    if (!content.photos?.length) return;
    setCurrentPhotoIndex((prev) =>
      prev > 0 ? prev - 1 : content.photos!.length - 1
    );
  }, [content.photos]);

  const handleNextPhoto = useCallback(() => {
    if (!content.photos?.length) return;
    setCurrentPhotoIndex((prev) =>
      prev < content.photos!.length - 1 ? prev + 1 : 0
    );
  }, [content.photos]);

  // Handle directions
  const handleGetDirections = useCallback(() => {
    if (readonly || !content.coordinates) return;

    if (onGetDirections) {
      onGetDirections(content.id, content.coordinates);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'get_directions',
        data: { venueId: content.id, coordinates: content.coordinates },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, content.coordinates, content.id, onGetDirections, onInteraction, message.id]);

  // Handle booking
  const handleBookVenue = useCallback(() => {
    if (readonly) return;

    if (onBookVenue) {
      onBookVenue(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'book_venue',
        data: { venueId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, content.id, onBookVenue, onInteraction, message.id]);

  // Handle call
  const handleCall = useCallback(() => {
    if (readonly || !content.phone) return;

    if (onCall) {
      onCall(content.phone);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'call_venue',
        data: { venueId: content.id, phone: content.phone },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, content.phone, content.id, onCall, onInteraction, message.id]);

  // Handle website visit
  const handleVisitWebsite = useCallback(() => {
    if (readonly || !content.website) return;

    if (onVisitWebsite) {
      onVisitWebsite(content.website);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'visit_website',
        data: { venueId: content.id, website: content.website },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, content.website, content.id, onVisitWebsite, onInteraction, message.id]);

  // Handle bookmark toggle
  const handleBookmark = useCallback(() => {
    if (readonly) return;

    setIsBookmarked(!isBookmarked);

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: isBookmarked ? 'remove_bookmark' : 'add_bookmark',
        data: { venueId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, isBookmarked, onInteraction, message.id, content.id]);

  // Get feature icon
  const getFeatureIcon = useCallback((feature: VenueFeature) => {
    switch (feature) {
      case VenueFeature.WIFI: return <Wifi className="w-4 h-4" />;
      case VenueFeature.PARKING: return <Car className="w-4 h-4" />;
      case VenueFeature.CREDIT_CARDS: return <CreditCard className="w-4 h-4" />;
      case VenueFeature.WHEELCHAIR_ACCESSIBLE: return <Accessibility className="w-4 h-4" />;
      case VenueFeature.RESERVATIONS: return <Calendar className="w-4 h-4" />;
      case VenueFeature.OUTDOOR_SEATING: return <Users className="w-4 h-4" />;
      case VenueFeature.LIVE_MUSIC: return <Music className="w-4 h-4" />;
      case VenueFeature.TAKEOUT: return <ShoppingCart className="w-4 h-4" />;
      case VenueFeature.DELIVERY: return <Navigation className="w-4 h-4" />;
      case VenueFeature.KID_FRIENDLY: return <Heart className="w-4 h-4" />;
      default: return <CheckCircle className="w-4 h-4" />;
    }
  }, []);

  // Get status color
  const getStatusColor = useCallback((isOpen: boolean | null) => {
    if (isOpen === null) return 'text-muted-foreground';
    return isOpen ? 'text-green-600' : 'text-red-600';
  }, []);

  // Format day name
  const formatDayName = useCallback((dayOfWeek: number) => {
    const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
    return days[dayOfWeek];
  }, []);

  // Format time
  const formatTime = useCallback((time: { hours: number; minutes: number }) => {
    return `${time.hours.toString().padStart(2, '0')}:${time.minutes.toString().padStart(2, '0')}`;
  }, []);

  // Performance optimization
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: venueVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={venueRef}
        className={cn(
          "venue-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`venue-message-${message.id}`}
        data-venue-id={content.id}
        role="article"
        aria-label={`Venue: ${content.name} - ${content.category || 'Business'}`}
      >
        <Card className="venue-card overflow-hidden">
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
                      Shared a venue • {message.timestamp.toLocaleDateString()}
                    </p>
                  )}
                </div>
              </div>

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
                    {isBookmarked ? 'Remove bookmark' : 'Bookmark venue'}
                  </TooltipContent>
                </Tooltip>
              </div>
            </div>
          </CardHeader>

          <CardContent className="space-y-4">
            {/* Venue photo */}
            {showPhotos && content.photos && content.photos.length > 0 && (
              <div className="relative aspect-video rounded-lg overflow-hidden bg-muted">
                <motion.img
                  key={currentPhotoIndex}
                  variants={performanceMode ? {} : photoVariants}
                  initial={performanceMode ? {} : "initial"}
                  animate={performanceMode ? {} : "animate"}
                  src={content.photos[currentPhotoIndex]}
                  alt={`${content.name} - Photo ${currentPhotoIndex + 1}`}
                  className="w-full h-full object-cover"
                />

                {/* Photo navigation */}
                {content.photos.length > 1 && (
                  <>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handlePrevPhoto}
                      className="absolute left-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 bg-background/80 backdrop-blur-sm"
                      disabled={readonly}
                    >
                      <ChevronLeft className="w-4 h-4" />
                    </Button>

                    <Button
                      variant="outline"
                      size="sm"
                      onClick={handleNextPhoto}
                      className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 bg-background/80 backdrop-blur-sm"
                      disabled={readonly}
                    >
                      <ChevronRight className="w-4 h-4" />
                    </Button>

                    {/* Photo indicators */}
                    <div className="absolute bottom-2 left-1/2 transform -translate-x-1/2 flex gap-1">
                      {content.photos.map((_, index) => (
                        <button
                          key={index}
                          onClick={() => setCurrentPhotoIndex(index)}
                          className={cn(
                            "w-2 h-2 rounded-full transition-colors",
                            index === currentPhotoIndex
                              ? "bg-white"
                              : "bg-white/50 hover:bg-white/75"
                          )}
                          aria-label={`View photo ${index + 1}`}
                        />
                      ))}
                    </div>
                  </>
                )}

                {/* Status overlay */}
                <div className="absolute top-2 right-2">
                  <Badge
                    variant={venueStatus.isOpen ? "default" : "secondary"}
                    className={cn(
                      "text-xs",
                      venueStatus.isOpen === true && "bg-green-600 text-white",
                      venueStatus.isOpen === false && "bg-red-600 text-white"
                    )}
                  >
                    {venueStatus.status}
                  </Badge>
                </div>
              </div>
            )}

            {/* Venue details */}
            <div className="space-y-3">
              {/* Name and rating */}
              <div className="space-y-2">
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0 flex-1">
                    <h3 className="font-bold text-lg text-foreground leading-tight">
                      {content.name}
                    </h3>

                    {content.category && (
                      <p className="text-sm text-muted-foreground">
                        {content.category}
                      </p>
                    )}
                  </div>

                  {content.rating && (
                    <div className="flex items-center gap-1 bg-muted px-2 py-1 rounded-full">
                      <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                      <span className="text-sm font-medium">
                        {formatRating(content.rating.average)}
                      </span>
                      {content.rating.count && (
                        <span className="text-xs text-muted-foreground">
                          ({content.rating.count})
                        </span>
                      )}
                    </div>
                  )}
                </div>

                {/* Price range */}
                {content.priceRange && (
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-muted-foreground">Price:</span>
                    <Badge variant="outline" className="text-xs">
                      {content.priceRange}
                    </Badge>
                  </div>
                )}
              </div>

              {/* Description */}
              {content.description && (
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {content.description}
                </p>
              )}

              {/* Address and contact */}
              <div className="space-y-2">
                <div className="flex items-start gap-2">
                  <MapPin className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
                  <div className="min-w-0 flex-1">
                    <p className="text-sm text-foreground">{content.address}</p>
                    {content.distance && (
                      <p className="text-xs text-muted-foreground">
                        {content.distance.toFixed(1)} km away
                      </p>
                    )}
                  </div>
                </div>

                {content.phone && (
                  <div className="flex items-center gap-2">
                    <Phone className="w-4 h-4 text-primary" />
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleCall}
                      disabled={readonly}
                      className="h-auto p-0 text-sm text-foreground hover:text-primary"
                    >
                      {content.phone}
                    </Button>
                  </div>
                )}

                {content.website && (
                  <div className="flex items-center gap-2">
                    <Globe className="w-4 h-4 text-primary" />
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleVisitWebsite}
                      disabled={readonly}
                      className="h-auto p-0 text-sm text-foreground hover:text-primary"
                    >
                      <span className="truncate">{content.website}</span>
                      <ExternalLink className="w-3 h-3 ml-1 flex-shrink-0" />
                    </Button>
                  </div>
                )}
              </div>

              {/* Business hours */}
              {content.businessHours && (
                <div className="space-y-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setShowHours(!showHours)}
                    className="h-auto p-0 text-sm font-medium text-foreground flex items-center gap-2"
                  >
                    <Clock className="w-4 h-4" />
                    <span className={getStatusColor(venueStatus.isOpen)}>
                      {venueStatus.status}
                    </span>
                    <ChevronRight className={cn(
                      "w-3 h-3 transition-transform",
                      showHours && "rotate-90"
                    )} />
                  </Button>

                  <AnimatePresence>
                    {showHours && (
                      <motion.div
                        variants={performanceMode ? {} : hoursVariants}
                        initial={performanceMode ? {} : "initial"}
                        animate={performanceMode ? {} : "animate"}
                        exit={performanceMode ? {} : "exit"}
                        className="ml-6 space-y-1"
                      >
                        {content.businessHours.map((hours) => (
                          <div key={hours.dayOfWeek} className="flex justify-between text-xs">
                            <span className="text-muted-foreground">
                              {formatDayName(hours.dayOfWeek)}
                            </span>
                            <span className="font-medium">
                              {hours.isOpen
                                ? `${formatTime(hours.openTime)} - ${formatTime(hours.closeTime)}`
                                : 'Closed'
                              }
                            </span>
                          </div>
                        ))}
                      </motion.div>
                    )}
                  </AnimatePresence>
                </div>
              )}

              {/* Features */}
              {content.features && content.features.length > 0 && (
                <>
                  <Separator />
                  <div className="space-y-2">
                    <h4 className="text-sm font-medium text-foreground">Features</h4>
                    <div className="flex flex-wrap gap-2">
                      {content.features.map((feature, index) => (
                        <div
                          key={index}
                          className="flex items-center gap-1 bg-muted px-2 py-1 rounded-full text-xs"
                        >
                          {getFeatureIcon(feature)}
                          <span className="capitalize">{feature.replace('_', ' ')}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                </>
              )}

              {/* Reviews preview */}
              {showReviews && content.reviews && content.reviews.length > 0 && (
                <>
                  <Separator />
                  <div className="space-y-2">
                    <h4 className="text-sm font-medium text-foreground">Recent Reviews</h4>
                    <div className="space-y-2">
                      {content.reviews.slice(0, 2).map((review, index) => (
                        <div key={index} className="bg-muted p-3 rounded-lg space-y-1">
                          <div className="flex items-center justify-between">
                            <span className="text-sm font-medium">{review.author}</span>
                            <div className="flex items-center gap-1">
                              <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
                              <span className="text-xs">{review.rating}</span>
                            </div>
                          </div>
                          <p className="text-xs text-muted-foreground line-clamp-2">
                            {review.text}
                          </p>
                          <span className="text-xs text-muted-foreground">
                            {new Date(review.date).toLocaleDateString()}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                </>
              )}

              <Separator />

              {/* Action buttons */}
              <motion.div
                variants={performanceMode ? {} : actionVariants}
                initial={performanceMode ? {} : "initial"}
                animate={performanceMode ? {} : "animate"}
                transition={{ delay: 0.2 }}
                className="flex flex-wrap gap-2"
              >
                <Button
                  onClick={handleGetDirections}
                  disabled={readonly || !content.coordinates}
                  className="gap-2"
                  size="sm"
                >
                  <Navigation className="w-4 h-4" />
                  Directions
                </Button>

                {content.acceptsReservations && (
                  <Button
                    variant="outline"
                    onClick={handleBookVenue}
                    disabled={readonly}
                    className="gap-2"
                    size="sm"
                  >
                    <Calendar className="w-4 h-4" />
                    Book
                  </Button>
                )}

                {content.phone && (
                  <Button
                    variant="outline"
                    onClick={handleCall}
                    disabled={readonly}
                    className="gap-2"
                    size="sm"
                  >
                    <Phone className="w-4 h-4" />
                    Call
                  </Button>
                )}

                <Button
                  variant="outline"
                  onClick={() => {/* Handle share */}}
                  disabled={readonly}
                  className="gap-2"
                  size="sm"
                >
                  <Share className="w-4 h-4" />
                  Share
                </Button>
              </motion.div>
            </div>
          </CardContent>
        </Card>

        {/* Performance Debug Info */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            {content.category || 'venue'} | {venueStatus.isOpen ? 'open' : 'closed'} | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedVenueMessage = React.memo(VenueMessage, (prevProps, nextProps) => {
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.showPhotos === nextProps.showPhotos &&
    prevProps.showReviews === nextProps.showReviews &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedVenueMessage.displayName = 'MemoizedVenueMessage';

export default VenueMessage;