import React from 'react';
import { ChevronRight, Loader2 } from 'lucide-react';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Card, CardContent } from '../ui/card';
import { ImageWithFallback } from '../figma/ImageWithFallback';
import { StreamContentItem } from '../../services/streamApi';

interface FeaturedContentCarouselProps {
  content: StreamContentItem[];
  isLoading: boolean;
  onContentClick?: (contentId: string) => void;
  onAddToCart?: (contentId: string, quantity?: number) => void;
  onSeeAllClick?: () => void;
  className?: string;
}

export function FeaturedContentCarousel({
  content,
  isLoading,
  onContentClick,
  onAddToCart,
  onSeeAllClick,
  className = ''
}: FeaturedContentCarouselProps) {
  // Format duration helper
  const formatDuration = (seconds?: number) => {
    if (!seconds) return null;
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    }
    return `${minutes}m`;
  };

  // Format price helper
  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency
    }).format(price);
  };

  // Handle add to cart
  const handleAddToCart = (contentId: string, event: React.MouseEvent) => {
    event.stopPropagation();
    onAddToCart?.(contentId, 1);
  };

  if (isLoading) {
    return (
      <div className={`space-y-4 ${className}`}>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">Featured</h2>
        </div>
        <div className="flex gap-4 overflow-x-auto scrollbar-hide pb-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex-shrink-0 w-64 h-48 bg-muted animate-pulse rounded-lg" />
          ))}
        </div>
      </div>
    );
  }

  if (content.length === 0) {
    return null;
  }

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">Featured</h2>
        {onSeeAllClick && (
          <Button
            variant="ghost"
            size="sm"
            onClick={onSeeAllClick}
            className="text-sm hover:text-primary"
          >
            See All
            <ChevronRight className="w-4 h-4 ml-1" />
          </Button>
        )}
      </div>

      {/* Carousel Container */}
      <div className="relative">
        <div className="flex gap-4 overflow-x-auto scrollbar-hide pb-2">
          {content.map((item) => (
            <Card
              key={item.id}
              className="flex-shrink-0 w-64 overflow-hidden hover:shadow-md transition-shadow cursor-pointer"
              onClick={() => onContentClick?.(item.id)}
            >
              {/* Thumbnail with Overlays */}
              <div className="relative">
                <ImageWithFallback
                  src={item.thumbnailUrl}
                  alt={item.title}
                  className="w-full h-32 object-cover"
                />

                {/* Availability Status Badge */}
                {item.availabilityStatus !== 'available' && (
                  <Badge className="absolute top-2 left-2 bg-orange-500">
                    {item.availabilityStatus === 'coming_soon' ? 'Coming Soon' : 'Unavailable'}
                  </Badge>
                )}

                {/* Duration Badge */}
                {item.duration && (
                  <Badge className="absolute bottom-2 right-2 bg-black/70 text-white">
                    {formatDuration(item.duration)}
                  </Badge>
                )}

                {/* Featured Badge */}
                {item.isFeatured && (
                  <Badge className="absolute top-2 right-2 bg-gradient-to-r from-yellow-400 to-orange-500 text-black">
                    Featured
                  </Badge>
                )}
              </div>

              {/* Content */}
              <CardContent className="p-4">
                <h3 className="font-medium mb-1 line-clamp-2">{item.title}</h3>
                <p className="text-sm text-muted-foreground mb-2 line-clamp-2">{item.description}</p>

                {/* Price and Action */}
                <div className="flex items-center justify-between">
                  <span className="font-bold text-lg">
                    {formatPrice(item.price, item.currency)}
                  </span>

                  {onAddToCart && (
                    <Button
                      size="sm"
                      onClick={(e) => handleAddToCart(item.id, e)}
                      disabled={item.availabilityStatus !== 'available'}
                      className="flex-shrink-0"
                    >
                      Add to Cart
                    </Button>
                  )}
                </div>

                {/* Metadata */}
                {item.metadata && (
                  <div className="mt-2 flex flex-wrap gap-1">
                    {item.metadata.tags?.slice(0, 2).map((tag: string, index: number) => (
                      <Badge key={index} variant="secondary" className="text-xs">
                        {tag}
                      </Badge>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Scroll Indicator (Optional) */}
        {content.length > 3 && (
          <div className="flex justify-center mt-2 space-x-1">
            {Array.from({ length: Math.ceil(content.length / 3) }).map((_, index) => (
              <div
                key={index}
                className="w-2 h-2 rounded-full bg-muted transition-colors"
                aria-hidden="true"
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default FeaturedContentCarousel;