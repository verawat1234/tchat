import React, { useState, useMemo } from 'react';
import { BookOpen, Mic, Film, Video, Music, Palette, ChevronRight, Loader2 } from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";
import {
  useGetStreamCategoriesQuery,
  useGetStreamContentQuery,
  useGetStreamFeaturedQuery,
  usePurchaseStreamContentMutation,
  StreamCategory,
  StreamContentItem,
  StreamSubtab,
} from '../services/streamApi';
import StreamTabs from './store/StreamTabs';
import FeaturedContentCarousel from './store/FeaturedContentCarousel';

// Interface definitions now imported from streamApi

interface StreamTabProps {
  user: any;
  onContentClick?: (contentId: string) => void;
  onAddToCart?: (contentId: string, quantity?: number) => void;
  onContentShare?: (contentId: string, contentData: any) => void;
  cartItems?: string[];
}

export function StreamTab({
  user,
  onContentClick,
  onAddToCart,
  onContentShare,
  cartItems = []
}: StreamTabProps) {
  const [selectedCategory, setSelectedCategory] = useState('books');
  const [selectedSubtab, setSelectedSubtab] = useState<string | null>(null);

  // Categories query for content display
  const {
    data: categoriesData
  } = useGetStreamCategoriesQuery();

  const {
    data: contentData,
    isLoading: contentLoading
  } = useGetStreamContentQuery({
    categoryId: selectedCategory,
    page: 1,
    limit: 20,
    ...(selectedSubtab && { subtabId: selectedSubtab })
  }, {
    skip: !selectedCategory
  });

  const {
    data: featuredData,
    isLoading: featuredLoading
  } = useGetStreamFeaturedQuery({
    categoryId: selectedCategory,
    limit: 10
  }, {
    skip: !selectedCategory
  });

  const [purchaseContent] = usePurchaseStreamContentMutation();

  // Icon mapping for empty state display
  const categoryIcons = {
    'book-open': BookOpen,
    'microphone': Mic,
    'film': Film,
    'video': Video,
    'music': Music,
    'palette': Palette
  } as const;

  // Extract data from API responses
  const categories = categoriesData?.categories || [];
  const contentItems = contentData?.items || [];
  const featuredItems = featuredData?.items || [];

  // Regular content comes directly from API (already filtered by category and subtab)
  const regularContent = useMemo(() => {
    return contentItems.filter(item => !item.isFeatured);
  }, [contentItems]);

  // Featured content comes from separate API call
  const featuredContent = useMemo(() => {
    return featuredItems.sort((a, b) => (a.featuredOrder || 999) - (b.featuredOrder || 999));
  }, [featuredItems]);

  const handleAddToCart = async (contentId: string) => {
    if (onAddToCart) {
      onAddToCart(contentId, 1);
    } else {
      // Use the purchase mutation for direct purchase
      try {
        const result = await purchaseContent({
          mediaContentId: contentId,
          quantity: 1,
          mediaLicense: 'personal',
        }).unwrap();

        if (result.success) {
          toast.success(result.message || 'Purchase successful!');
        } else {
          toast.error(result.message || 'Purchase failed');
        }
      } catch (error) {
        console.error('Purchase error:', error);
        toast.error('Failed to process purchase');
      }
    }
  };

  const formatDuration = (seconds?: number) => {
    if (!seconds) return null;
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    }
    return `${minutes}m`;
  };

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency
    }).format(price);
  };

  return (
    <div className="flex flex-col h-full">
      {/* Stream Navigation Tabs */}
      <StreamTabs
        selectedCategory={selectedCategory}
        selectedSubtab={selectedSubtab}
        onCategoryChange={setSelectedCategory}
        onSubtabChange={setSelectedSubtab}
      />

      {/* Content */}
      <div className="flex-1 overflow-hidden">
        <ScrollArea className="h-full px-4">
          <div className="space-y-6 py-4">
            {/* Featured Content Carousel */}
            <FeaturedContentCarousel
              content={featuredContent}
              isLoading={featuredLoading}
              onContentClick={onContentClick}
              onAddToCart={onAddToCart}
              onSeeAllClick={() => {
                // TODO: Implement see all featured content
                console.log('See all featured content');
              }}
            />

            {/* Regular Content Grid */}
            {contentLoading ? (
              <div className="space-y-4">
                <h2 className="text-lg font-semibold mb-4">Loading content...</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {[1, 2, 3, 4, 5, 6].map((i) => (
                    <div key={i} className="h-64 bg-muted animate-pulse rounded-lg" />
                  ))}
                </div>
              </div>
            ) : regularContent.length > 0 && (
              <div>
                <h2 className="text-lg font-semibold mb-4">
                  Browse {categories.find(c => c.id === selectedCategory)?.name}
                </h2>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {regularContent.map((content) => (
                    <Card
                      key={content.id}
                      className="overflow-hidden hover:shadow-md transition-shadow cursor-pointer"
                      onClick={() => onContentClick?.(content.id)}
                    >
                      <div className="relative">
                        <ImageWithFallback
                          src={content.thumbnailUrl}
                          alt={content.title}
                          className="w-full h-32 object-cover"
                        />
                        {content.availabilityStatus !== 'available' && (
                          <Badge className="absolute top-2 left-2 bg-orange-500">
                            {content.availabilityStatus === 'coming_soon' ? 'Coming Soon' : 'Unavailable'}
                          </Badge>
                        )}
                        {content.duration && (
                          <Badge className="absolute bottom-2 right-2 bg-black/70 text-white">
                            {formatDuration(content.duration)}
                          </Badge>
                        )}
                      </div>

                      <CardContent className="p-4">
                        <h3 className="font-medium mb-1 line-clamp-2">{content.title}</h3>
                        <p className="text-sm text-muted-foreground mb-3 line-clamp-2">{content.description}</p>

                        <div className="flex items-center justify-between">
                          <span className="font-bold">
                            {formatPrice(content.price, content.currency)}
                          </span>

                          <Button
                            size="sm"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleAddToCart(content.id);
                            }}
                            disabled={content.availabilityStatus !== 'available'}
                          >
                            Add to Cart
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>
            )}

            {/* Empty State */}
            {!contentLoading && !featuredLoading && featuredContent.length === 0 && regularContent.length === 0 && (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center mb-4">
                  {React.createElement(categoryIcons[categories.find(c => c.id === selectedCategory)?.iconName as keyof typeof categoryIcons] || BookOpen, { className: "w-8 h-8 text-muted-foreground" })}
                </div>
                <h3 className="text-lg font-medium mb-2">No content available</h3>
                <p className="text-muted-foreground">
                  Check back later for new {categories.find(c => c.id === selectedCategory)?.name.toLowerCase()} content.
                </p>
              </div>
            )}
          </div>
        </ScrollArea>
      </div>
    </div>
  );
}