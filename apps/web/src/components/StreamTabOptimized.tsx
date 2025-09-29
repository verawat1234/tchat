/**
 * Optimized StreamTab Component
 *
 * Performance-optimized version of the StreamTab component targeting <1s content load times.
 * Implements advanced caching, prefetching, virtual scrolling, image optimization,
 * and memory management for superior user experience.
 *
 * Performance Targets:
 * - Initial load: <1s
 * - Tab switching: <200ms
 * - Content rendering: <16ms (60fps)
 * - Memory usage: <100MB mobile, <500MB desktop
 * - Cache efficiency: >95% hit rate
 *
 * Optimization Techniques:
 * - Smart prefetching with prediction algorithms
 * - Virtual scrolling for large content lists
 * - Image lazy loading with progressive enhancement
 * - Memoized components and selectors
 * - Background cache warming
 * - Intersection Observer for viewport optimization
 * - Service Worker integration for offline caching
 */

import React, { useState, useMemo, useCallback, useEffect, useRef } from 'react';
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
  streamApi,
  StreamCategory,
  StreamContentItem,
  StreamSubtab,
} from '../services/streamApi';
import { useAppDispatch } from '../store/hooks';
import StreamTabsOptimized from './store/StreamTabsOptimized';
import FeaturedContentCarouselOptimized from './store/FeaturedContentCarouselOptimized';
import { VirtualizedContentGrid } from './store/VirtualizedContentGrid';
import { ContentPreloader } from './store/ContentPreloader';
import { PerformanceMonitor } from './utils/PerformanceMonitor';

// Performance configuration
const PERFORMANCE_CONFIG = {
  CONTENT_PREFETCH_DELAY: 150, // ms
  CACHE_WARMING_INTERVAL: 30000, // 30s
  VIRTUAL_SCROLL_THRESHOLD: 50, // items
  IMAGE_LAZY_LOAD_MARGIN: '200px',
  DEBOUNCE_DELAY: 100, // ms
  MAX_CACHE_SIZE: 1000, // items
  PRELOAD_ADJACENT_CATEGORIES: 2,
} as const;

interface StreamTabOptimizedProps {
  user: any;
  onContentClick?: (contentId: string) => void;
  onAddToCart?: (contentId: string, quantity?: number) => void;
  onContentShare?: (contentId: string, contentData: any) => void;
  cartItems?: string[];
  performanceMode?: 'standard' | 'high' | 'battery';
}

export function StreamTabOptimized({
  user,
  onContentClick,
  onAddToCart,
  onContentShare,
  cartItems = [],
  performanceMode = 'standard'
}: StreamTabOptimizedProps) {
  const dispatch = useAppDispatch();
  const performanceMonitor = useRef(new PerformanceMonitor()).current;

  // State management with performance tracking
  const [selectedCategory, setSelectedCategory] = useState('books');
  const [selectedSubtab, setSelectedSubtab] = useState<string | null>(null);
  const [isTransitioning, setIsTransitioning] = useState(false);

  // Refs for performance optimization
  const contentContainerRef = useRef<HTMLDivElement>(null);
  const prefetchTimerRef = useRef<NodeJS.Timeout | null>(null);
  const transitionStartRef = useRef<number>(0);

  // Memoized performance configuration based on mode
  const config = useMemo(() => {
    const baseConfig = PERFORMANCE_CONFIG;
    switch (performanceMode) {
      case 'high':
        return {
          ...baseConfig,
          CONTENT_PREFETCH_DELAY: 50,
          PRELOAD_ADJACENT_CATEGORIES: 3,
          VIRTUAL_SCROLL_THRESHOLD: 30,
        };
      case 'battery':
        return {
          ...baseConfig,
          CONTENT_PREFETCH_DELAY: 300,
          PRELOAD_ADJACENT_CATEGORIES: 1,
          CACHE_WARMING_INTERVAL: 60000,
        };
      default:
        return baseConfig;
    }
  }, [performanceMode]);

  // RTK Query hooks with performance optimizations
  const {
    data: categoriesData,
    isLoading: categoriesLoading,
    error: categoriesError
  } = useGetStreamCategoriesQuery(undefined, {
    // Aggressive caching for categories (rarely change)
    pollingInterval: 300000, // 5 minutes
    skipPollingIfUnfocused: true,
  });

  const {
    data: contentData,
    isLoading: contentLoading,
    isFetching: contentFetching,
    error: contentError
  } = useGetStreamContentQuery({
    categoryId: selectedCategory,
    page: 1,
    limit: 20,
    ...(selectedSubtab && { subtabId: selectedSubtab })
  }, {
    skip: !selectedCategory,
    // Smart refetching strategy
    refetchOnMountOrArgChange: 30, // 30 seconds
    refetchOnFocus: false,
    refetchOnReconnect: true,
  });

  const {
    data: featuredData,
    isLoading: featuredLoading,
    isFetching: featuredFetching
  } = useGetStreamFeaturedQuery({
    categoryId: selectedCategory,
    limit: 10
  }, {
    skip: !selectedCategory,
    // Prefetch featured content aggressively
    refetchOnMountOrArgChange: 60, // 1 minute
  });

  const [purchaseContent] = usePurchaseStreamContentMutation();

  // Memoized data extraction with performance monitoring
  const categories = useMemo(() => {
    const start = performance.now();
    const result = categoriesData?.categories || [];
    performanceMonitor.recordOperation('categories_memoization', performance.now() - start);
    return result;
  }, [categoriesData?.categories, performanceMonitor]);

  const contentItems = useMemo(() => {
    const start = performance.now();
    const result = contentData?.items || [];
    performanceMonitor.recordOperation('content_memoization', performance.now() - start);
    return result;
  }, [contentData?.items, performanceMonitor]);

  const featuredItems = useMemo(() => {
    const start = performance.now();
    const result = (featuredData?.items || [])
      .sort((a, b) => (a.featuredOrder || 999) - (b.featuredOrder || 999));
    performanceMonitor.recordOperation('featured_memoization', performance.now() - start);
    return result;
  }, [featuredData?.items, performanceMonitor]);

  // Performance-optimized content filtering
  const { regularContent, featuredContent } = useMemo(() => {
    const start = performance.now();

    const regular = contentItems.filter(item => !item.isFeatured);
    const featured = featuredItems;

    performanceMonitor.recordOperation('content_filtering', performance.now() - start);

    return {
      regularContent: regular,
      featuredContent: featured
    };
  }, [contentItems, featuredItems, performanceMonitor]);

  // Smart prefetching strategy
  const prefetchAdjacentCategories = useCallback((currentCategoryId: string) => {
    if (prefetchTimerRef.current) {
      clearTimeout(prefetchTimerRef.current);
    }

    prefetchTimerRef.current = setTimeout(() => {
      const currentIndex = categories.findIndex(cat => cat.id === currentCategoryId);
      if (currentIndex === -1) return;

      const adjacentCategories = [];

      // Prefetch previous and next categories
      for (let i = 1; i <= config.PRELOAD_ADJACENT_CATEGORIES; i++) {
        if (currentIndex - i >= 0) {
          adjacentCategories.push(categories[currentIndex - i]);
        }
        if (currentIndex + i < categories.length) {
          adjacentCategories.push(categories[currentIndex + i]);
        }
      }

      // Prefetch content for adjacent categories
      adjacentCategories.forEach(category => {
        dispatch(streamApi.util.prefetch('getStreamContent', {
          categoryId: category.id,
          page: 1,
          limit: 20
        }, { force: false }));

        dispatch(streamApi.util.prefetch('getStreamFeatured', {
          categoryId: category.id,
          limit: 10
        }, { force: false }));
      });

      performanceMonitor.recordOperation('prefetch_adjacent', adjacentCategories.length);
    }, config.CONTENT_PREFETCH_DELAY);
  }, [categories, dispatch, config, performanceMonitor]);

  // Optimized category change handler with performance tracking
  const handleCategoryChange = useCallback((categoryId: string) => {
    transitionStartRef.current = performance.now();
    setIsTransitioning(true);

    // Batch state updates
    setSelectedCategory(categoryId);
    setSelectedSubtab(null);

    // Prefetch adjacent categories
    prefetchAdjacentCategories(categoryId);

    // Track transition performance
    requestAnimationFrame(() => {
      const transitionTime = performance.now() - transitionStartRef.current;
      performanceMonitor.recordOperation('category_transition', transitionTime);
      setIsTransitioning(false);

      // Warn if transition is slow
      if (transitionTime > 200) {
        console.warn(`Slow category transition: ${transitionTime}ms for ${categoryId}`);
      }
    });
  }, [prefetchAdjacentCategories, performanceMonitor]);

  // Optimized subtab change handler
  const handleSubtabChange = useCallback((subtabId: string | null) => {
    const start = performance.now();
    setSelectedSubtab(subtabId);

    const transitionTime = performance.now() - start;
    performanceMonitor.recordOperation('subtab_transition', transitionTime);
  }, [performanceMonitor]);

  // Optimized add to cart handler with error recovery
  const handleAddToCart = useCallback(async (contentId: string) => {
    const start = performance.now();

    try {
      if (onAddToCart) {
        onAddToCart(contentId, 1);
      } else {
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
      }

      const operationTime = performance.now() - start;
      performanceMonitor.recordOperation('add_to_cart', operationTime);
    } catch (error) {
      console.error('Purchase error:', error);
      toast.error('Failed to process purchase');

      const errorTime = performance.now() - start;
      performanceMonitor.recordError('add_to_cart_error', errorTime);
    }
  }, [onAddToCart, purchaseContent, performanceMonitor]);

  // Performance monitoring and cache warming
  useEffect(() => {
    // Start performance monitoring
    performanceMonitor.startSession();

    // Warm up cache with popular categories
    const warmupTimer = setTimeout(() => {
      if (categories.length > 0) {
        const popularCategories = categories.slice(0, 3); // Books, Podcasts, Cartoons

        popularCategories.forEach(category => {
          if (category.id !== selectedCategory) {
            dispatch(streamApi.util.prefetch('getStreamContent', {
              categoryId: category.id,
              page: 1,
              limit: 20
            }, { force: false }));
          }
        });

        performanceMonitor.recordOperation('cache_warmup', popularCategories.length);
      }
    }, config.CACHE_WARMING_INTERVAL);

    return () => {
      clearTimeout(warmupTimer);
      if (prefetchTimerRef.current) {
        clearTimeout(prefetchTimerRef.current);
      }
      performanceMonitor.endSession();
    };
  }, [categories, dispatch, selectedCategory, config, performanceMonitor]);

  // Trigger prefetching on category selection
  useEffect(() => {
    if (selectedCategory && categories.length > 0) {
      prefetchAdjacentCategories(selectedCategory);
    }
  }, [selectedCategory, categories, prefetchAdjacentCategories]);

  // Icon mapping for performance (memoized)
  const categoryIcons = useMemo(() => ({
    'book-open': BookOpen,
    'microphone': Mic,
    'film': Film,
    'video': Video,
    'music': Music,
    'palette': Palette
  }), []);

  // Utility functions (memoized)
  const formatDuration = useCallback((seconds?: number) => {
    if (!seconds) return null;
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    }
    return `${minutes}m`;
  }, []);

  const formatPrice = useCallback((price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency
    }).format(price);
  }, []);

  // Error boundary fallback
  if (categoriesError || contentError) {
    return (
      <div className="flex flex-col items-center justify-center h-full p-8">
        <div className="text-center">
          <h3 className="text-lg font-semibold mb-2">Unable to load content</h3>
          <p className="text-muted-foreground mb-4">
            {categoriesError ? 'Failed to load categories' : 'Failed to load content'}
          </p>
          <Button
            onClick={() => window.location.reload()}
            variant="outline"
          >
            Retry
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full" ref={contentContainerRef}>
      {/* Performance-optimized Stream Navigation */}
      <StreamTabsOptimized
        selectedCategory={selectedCategory}
        selectedSubtab={selectedSubtab}
        onCategoryChange={handleCategoryChange}
        onSubtabChange={handleSubtabChange}
        isTransitioning={isTransitioning}
        performanceMode={performanceMode}
      />

      {/* Optimized Content Area */}
      <div className="flex-1 overflow-hidden">
        <ScrollArea className="h-full px-4">
          <div className="space-y-6 py-4">
            {/* Performance-optimized Featured Content */}
            <FeaturedContentCarouselOptimized
              content={featuredContent}
              isLoading={featuredLoading || featuredFetching}
              onContentClick={onContentClick}
              onAddToCart={handleAddToCart}
              onSeeAllClick={() => {
                console.log('See all featured content');
              }}
              performanceMode={performanceMode}
            />

            {/* Virtualized Content Grid for Performance */}
            {contentLoading && !contentFetching ? (
              <div className="space-y-4">
                <h2 className="text-lg font-semibold mb-4">Loading content...</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {Array.from({ length: 6 }, (_, i) => (
                    <div key={i} className="h-64 bg-muted animate-pulse rounded-lg" />
                  ))}
                </div>
              </div>
            ) : regularContent.length > 0 ? (
              <div>
                <h2 className="text-lg font-semibold mb-4">
                  Browse {categories.find(c => c.id === selectedCategory)?.name}
                </h2>

                {/* Use virtualized grid for large content lists */}
                {regularContent.length > config.VIRTUAL_SCROLL_THRESHOLD ? (
                  <VirtualizedContentGrid
                    content={regularContent}
                    onContentClick={onContentClick}
                    onAddToCart={handleAddToCart}
                    formatDuration={formatDuration}
                    formatPrice={formatPrice}
                    cartItems={cartItems}
                    performanceMode={performanceMode}
                  />
                ) : (
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
                            loading="lazy"
                            decoding="async"
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
                          <h3 className="font-medium truncate">{content.title}</h3>
                          <p className="text-sm text-muted-foreground truncate mt-1">
                            {content.description}
                          </p>
                          <div className="flex items-center justify-between mt-3">
                            <span className="font-semibold">
                              {formatPrice(content.price, content.currency)}
                            </span>
                            <Button
                              size="sm"
                              variant={cartItems.includes(content.id) ? "secondary" : "default"}
                              onClick={(e) => {
                                e.stopPropagation();
                                handleAddToCart(content.id);
                              }}
                              disabled={content.availabilityStatus !== 'available'}
                            >
                              {cartItems.includes(content.id) ? 'In Cart' : 'Add to Cart'}
                            </Button>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                )}
              </div>
            ) : !contentLoading && (
              <div className="text-center py-12">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-4">
                  {React.createElement(categoryIcons['book-open'], {
                    className: "w-8 h-8 text-muted-foreground"
                  })}
                </div>
                <h3 className="text-lg font-semibold mb-2">No content available</h3>
                <p className="text-muted-foreground">
                  Check back later for new content in this category.
                </p>
              </div>
            )}
          </div>
        </ScrollArea>
      </div>

      {/* Background Content Preloader */}
      <ContentPreloader
        selectedCategory={selectedCategory}
        categories={categories}
        performanceMode={performanceMode}
      />
    </div>
  );
}