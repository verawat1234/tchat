/**
 * FeaturedContentCarouselOptimized Component
 *
 * High-performance featured content carousel for Stream categories.
 * Implements virtual scrolling, smart image loading, and gesture optimization
 * to achieve smooth 60fps animations and <500ms image load times.
 *
 * Performance Features:
 * - Virtual scrolling for large content lists
 * - Progressive image loading with WebP/AVIF support
 * - Intersection Observer for lazy loading
 * - Memory-efficient viewport management
 * - Hardware-accelerated animations
 * - Prefetch strategies for smooth navigation
 */

import React, { memo, useCallback, useMemo, useRef, useEffect, useState } from 'react';
import { PerformanceMonitor } from './utils/PerformanceMonitor';

export interface FeaturedContentItem {
  id: string;
  title: string;
  subtitle?: string;
  image_url: string;
  price?: number;
  discount?: number;
  rating?: number;
  category: string;
  priority: number;
  metadata?: {
    duration?: string;
    author?: string;
    release_date?: string;
    content_type: 'video' | 'audio' | 'text' | 'interactive';
  };
}

export interface FeaturedContentCarouselProps {
  items: FeaturedContentItem[];
  category: string;
  onItemClick: (item: FeaturedContentItem) => void;
  onAddToCart?: (item: FeaturedContentItem) => void;
  performanceMonitor?: PerformanceMonitor;
  className?: string;
  autoPlay?: boolean;
  autoPlayInterval?: number;
}

// Performance budgets and constants
const IMAGE_LOAD_BUDGET = 500; // 500ms budget for image loading
const ANIMATION_DURATION = 300; // 300ms for smooth animations
const VIEWPORT_BUFFER = 2; // Render 2 items beyond visible area
const INTERSECTION_THRESHOLD = 0.1; // 10% visibility threshold
const PREFETCH_DISTANCE = 3; // Prefetch 3 items ahead

// Image format preferences for modern browsers
const IMAGE_FORMATS = {
  avif: 'image/avif',
  webp: 'image/webp',
  jpeg: 'image/jpeg'
};

export const FeaturedContentCarouselOptimized: React.FC<FeaturedContentCarouselProps> = memo(({
  items,
  category,
  onItemClick,
  onAddToCart,
  performanceMonitor,
  className = '',
  autoPlay = false,
  autoPlayInterval = 5000
}) => {
  const carouselRef = useRef<HTMLDivElement>(null);
  const observerRef = useRef<IntersectionObserver | null>(null);
  const autoPlayRef = useRef<NodeJS.Timeout | null>(null);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [visibleRange, setVisibleRange] = useState({ start: 0, end: 6 });
  const [loadedImages, setLoadedImages] = useState<Set<string>>(new Set());
  const [imageErrors, setImageErrors] = useState<Set<string>>(new Set());

  // Memoized sorted items by priority
  const sortedItems = useMemo(() => {
    return [...items].sort((a, b) => a.priority - b.priority);
  }, [items]);

  // Calculate visible items with virtual scrolling
  const visibleItems = useMemo(() => {
    const start = Math.max(0, visibleRange.start - VIEWPORT_BUFFER);
    const end = Math.min(sortedItems.length, visibleRange.end + VIEWPORT_BUFFER);
    return sortedItems.slice(start, end).map((item, index) => ({
      ...item,
      virtualIndex: start + index
    }));
  }, [sortedItems, visibleRange]);

  // Generate optimized image sources with format support
  const generateImageSources = useCallback((imageUrl: string) => {
    const baseUrl = imageUrl.replace(/\.(jpe?g|png)$/i, '');
    return [
      { srcSet: `${baseUrl}.avif`, type: IMAGE_FORMATS.avif },
      { srcSet: `${baseUrl}.webp`, type: IMAGE_FORMATS.webp },
      { srcSet: imageUrl, type: IMAGE_FORMATS.jpeg }
    ];
  }, []);

  // Optimized image loading with performance monitoring
  const handleImageLoad = useCallback((itemId: string, loadStartTime: number) => {
    const loadTime = performance.now() - loadStartTime;

    setLoadedImages(prev => new Set(prev).add(itemId));

    performanceMonitor?.recordOperation(
      'featured_image_load',
      loadTime,
      'load',
      {
        item_id: itemId,
        category,
        budget: IMAGE_LOAD_BUDGET,
        success: true
      }
    );
  }, [performanceMonitor, category]);

  // Handle image load errors with fallback
  const handleImageError = useCallback((itemId: string, error: Event) => {
    setImageErrors(prev => new Set(prev).add(itemId));

    performanceMonitor?.recordOperation(
      'featured_image_error',
      0,
      'error',
      {
        item_id: itemId,
        category,
        error: 'Image load failed'
      }
    );
  }, [performanceMonitor, category]);

  // Optimized scroll handler with viewport calculation
  const handleScroll = useCallback(() => {
    if (!carouselRef.current) return;

    const { scrollLeft, clientWidth } = carouselRef.current;
    const itemWidth = 280; // Approximate item width
    const itemsPerView = Math.ceil(clientWidth / itemWidth);

    const start = Math.floor(scrollLeft / itemWidth);
    const end = start + itemsPerView;

    setVisibleRange({ start, end });

    // Update current index for auto-play
    const centerIndex = Math.round((scrollLeft + clientWidth / 2) / itemWidth);
    setCurrentIndex(Math.max(0, Math.min(centerIndex, sortedItems.length - 1)));
  }, [sortedItems.length]);

  // Navigation handlers with smooth scrolling
  const scrollToItem = useCallback((index: number) => {
    if (!carouselRef.current) return;

    const itemWidth = 280;
    const targetScroll = index * itemWidth;

    const startTime = performance.now();

    carouselRef.current.scrollTo({
      left: targetScroll,
      behavior: 'smooth'
    });

    // Monitor scroll performance
    requestAnimationFrame(() => {
      const duration = performance.now() - startTime;
      performanceMonitor?.recordOperation(
        'carousel_navigation',
        duration,
        'interaction',
        { target_index: index, category }
      );
    });

    setCurrentIndex(index);
  }, [performanceMonitor, category]);

  const handlePrevious = useCallback(() => {
    const prevIndex = Math.max(0, currentIndex - 1);
    scrollToItem(prevIndex);
  }, [currentIndex, scrollToItem]);

  const handleNext = useCallback(() => {
    const nextIndex = Math.min(sortedItems.length - 1, currentIndex + 1);
    scrollToItem(nextIndex);
  }, [currentIndex, sortedItems.length, scrollToItem]);

  // Item interaction handlers with performance tracking
  const handleItemClick = useCallback((item: FeaturedContentItem) => {
    const startTime = performance.now();

    onItemClick(item);

    performanceMonitor?.recordOperation(
      'featured_item_click',
      performance.now() - startTime,
      'interaction',
      { item_id: item.id, category: item.category }
    );
  }, [onItemClick, performanceMonitor]);

  const handleAddToCart = useCallback((item: FeaturedContentItem, e: React.MouseEvent) => {
    e.stopPropagation();

    const startTime = performance.now();

    onAddToCart?.(item);

    performanceMonitor?.recordOperation(
      'add_to_cart_featured',
      performance.now() - startTime,
      'interaction',
      { item_id: item.id, category: item.category, price: item.price }
    );
  }, [onAddToCart, performanceMonitor]);

  // Auto-play functionality
  useEffect(() => {
    if (!autoPlay) return;

    autoPlayRef.current = setInterval(() => {
      setCurrentIndex(prev => {
        const nextIndex = (prev + 1) % sortedItems.length;
        scrollToItem(nextIndex);
        return nextIndex;
      });
    }, autoPlayInterval);

    return () => {
      if (autoPlayRef.current) {
        clearInterval(autoPlayRef.current);
      }
    };
  }, [autoPlay, autoPlayInterval, sortedItems.length, scrollToItem]);

  // Intersection Observer for lazy loading
  useEffect(() => {
    observerRef.current = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const img = entry.target as HTMLImageElement;
            const itemId = img.dataset.itemId;

            if (itemId && !loadedImages.has(itemId)) {
              const loadStartTime = performance.now();
              img.dataset.loadStart = loadStartTime.toString();
            }
          }
        });
      },
      { threshold: INTERSECTION_THRESHOLD }
    );

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [loadedImages]);

  // Prefetch upcoming images
  useEffect(() => {
    const prefetchIndex = currentIndex + PREFETCH_DISTANCE;
    if (prefetchIndex < sortedItems.length) {
      const prefetchItem = sortedItems[prefetchIndex];
      if (!loadedImages.has(prefetchItem.id)) {
        const link = document.createElement('link');
        link.rel = 'prefetch';
        link.href = prefetchItem.image_url;
        document.head.appendChild(link);
      }
    }
  }, [currentIndex, sortedItems, loadedImages]);

  // Render individual content item
  const renderContentItem = useCallback((item: FeaturedContentItem & { virtualIndex: number }) => {
    const isLoaded = loadedImages.has(item.id);
    const hasError = imageErrors.has(item.id);
    const sources = generateImageSources(item.image_url);

    return (
      <div
        key={item.id}
        className="featured-item"
        data-testid="FeaturedItem"
        onClick={() => handleItemClick(item)}
        style={{ transform: `translateX(${item.virtualIndex * 280}px)` }}
      >
        <div className="item-image-container">
          {!hasError ? (
            <picture>
              {sources.map((source, index) => (
                <source key={index} srcSet={source.srcSet} type={source.type} />
              ))}
              <img
                src={item.image_url}
                alt={item.title}
                className={`item-image ${isLoaded ? 'loaded' : 'loading'}`}
                data-item-id={item.id}
                loading="lazy"
                onLoad={(e) => {
                  const loadStart = parseFloat((e.target as HTMLImageElement).dataset.loadStart || '0');
                  if (loadStart) {
                    handleImageLoad(item.id, loadStart);
                  }
                }}
                onError={(e) => handleImageError(item.id, e)}
                ref={(img) => {
                  if (img && observerRef.current) {
                    observerRef.current.observe(img);
                  }
                }}
              />
            </picture>
          ) : (
            <div className="image-placeholder" aria-label="Image unavailable">
              <span className="placeholder-icon">📷</span>
            </div>
          )}

          {!isLoaded && !hasError && (
            <div className="loading-skeleton" aria-hidden="true" />
          )}
        </div>

        <div className="item-content">
          <h3 className="item-title">{item.title}</h3>
          {item.subtitle && <p className="item-subtitle">{item.subtitle}</p>}

          <div className="item-metadata">
            {item.metadata?.author && (
              <span className="item-author">by {item.metadata.author}</span>
            )}
            {item.metadata?.duration && (
              <span className="item-duration">{item.metadata.duration}</span>
            )}
            {item.rating && (
              <span className="item-rating" aria-label={`Rating: ${item.rating} out of 5`}>
                ⭐ {item.rating.toFixed(1)}
              </span>
            )}
          </div>

          <div className="item-actions">
            {item.price && (
              <div className="item-pricing">
                {item.discount && (
                  <span className="original-price">${item.price.toFixed(2)}</span>
                )}
                <span className="current-price">
                  ${((item.price * (1 - (item.discount || 0) / 100))).toFixed(2)}
                </span>
                {item.discount && (
                  <span className="discount-badge">{item.discount}% off</span>
                )}
              </div>
            )}

            {onAddToCart && item.price && (
              <button
                className="add-to-cart-btn"
                onClick={(e) => handleAddToCart(item, e)}
                data-testid="AddToCartButton"
                aria-label={`Add ${item.title} to cart`}
              >
                Add to Cart
              </button>
            )}
          </div>
        </div>
      </div>
    );
  }, [
    loadedImages,
    imageErrors,
    generateImageSources,
    handleItemClick,
    handleImageLoad,
    handleImageError,
    handleAddToCart,
    onAddToCart
  ]);

  if (sortedItems.length === 0) {
    return (
      <div className={`featured-carousel-empty ${className}`}>
        <p>No featured content available</p>
      </div>
    );
  }

  return (
    <div className={`featured-carousel-optimized ${className}`}>
      <div className="carousel-header">
        <h2 className="carousel-title">Featured in {category}</h2>
        <div className="carousel-controls">
          <button
            className="nav-btn prev"
            onClick={handlePrevious}
            disabled={currentIndex === 0}
            aria-label="Previous items"
          >
            ←
          </button>
          <span className="carousel-indicator">
            {currentIndex + 1} / {sortedItems.length}
          </span>
          <button
            className="nav-btn next"
            onClick={handleNext}
            disabled={currentIndex === sortedItems.length - 1}
            aria-label="Next items"
          >
            →
          </button>
        </div>
      </div>

      <div
        ref={carouselRef}
        className="carousel-container"
        data-testid="FeaturedCarousel"
        onScroll={handleScroll}
        role="region"
        aria-label="Featured content carousel"
      >
        <div className="carousel-track">
          {visibleItems.map(renderContentItem)}
        </div>
      </div>

      {/* Progress indicators */}
      <div className="carousel-dots" role="tablist" aria-label="Carousel navigation">
        {sortedItems.map((_, index) => (
          <button
            key={index}
            className={`carousel-dot ${index === currentIndex ? 'active' : ''}`}
            onClick={() => scrollToItem(index)}
            role="tab"
            aria-selected={index === currentIndex}
            aria-label={`Go to item ${index + 1}`}
          />
        ))}
      </div>
    </div>
  );
});

FeaturedContentCarouselOptimized.displayName = 'FeaturedContentCarouselOptimized';

export default FeaturedContentCarouselOptimized;