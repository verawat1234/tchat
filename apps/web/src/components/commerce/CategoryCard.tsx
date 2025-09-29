/**
 * CategoryCard - Individual category display component
 *
 * Displays a single category with image, name, description, and product count.
 * Supports multiple display variants and interactive states.
 *
 * Features:
 * - Multiple card variants (default, compact, featured, minimal)
 * - Hover effects and animations
 * - Accessibility with ARIA attributes
 * - Image fallback and optimization
 * - Product count badges
 * - Featured category indicators
 * - Category color theming
 * - Loading and error states
 */

import React, { useState, useRef, useEffect } from 'react';
import { Card, CardContent, CardHeader } from '../ui/card';
import { Badge } from '../ui/badge';
import { Button } from '../ui/button';
import { Skeleton } from '../ui/skeleton';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import {
  Star,
  Package,
  ChevronRight,
  Eye,
  TrendingUp,
  Folder,
  Image as ImageIcon,
} from 'lucide-react';
import { cn } from '../../lib/utils';
import { useCategory } from './CategoryProvider';
import { ImageWithFallback } from '../figma/ImageWithFallback';
import type { Category } from '../../types/commerce';

// ===== Types =====

export interface CategoryCardProps {
  category: Category;
  variant?: 'default' | 'compact' | 'featured' | 'minimal' | 'grid';
  size?: 'sm' | 'md' | 'lg';
  showDescription?: boolean;
  showProductCount?: boolean;
  showAnalytics?: boolean;
  showSubcategories?: boolean;
  maxSubcategories?: number;
  className?: string;
  onClick?: (category: Category) => void;
  onSubcategoryClick?: (subcategory: Category) => void;
}

// ===== Subcategory List Component =====

interface SubcategoryListProps {
  parentCategory: Category;
  maxCount: number;
  onSubcategoryClick?: (subcategory: Category) => void;
}

function SubcategoryList({ parentCategory, maxCount, onSubcategoryClick }: SubcategoryListProps) {
  const { categoriesQuery } = useCategory();

  const subcategories = categoriesQuery.data?.categories.filter(
    cat => cat.parentId === parentCategory.id && cat.isVisible
  ).slice(0, maxCount) || [];

  if (subcategories.length === 0) return null;

  return (
    <div className="mt-3 pt-3 border-t border-border">
      <h4 className="text-xs font-medium text-muted-foreground mb-2">Subcategories</h4>
      <div className="flex flex-wrap gap-1">
        {subcategories.map(subcategory => (
          <Button
            key={subcategory.id}
            variant="secondary"
            size="sm"
            className="h-6 px-2 text-xs"
            onClick={(e) => {
              e.stopPropagation();
              onSubcategoryClick?.(subcategory);
            }}
          >
            {subcategory.name}
          </Button>
        ))}
        {parentCategory.childrenCount > maxCount && (
          <Badge variant="outline" className="h-6 px-2 text-xs">
            +{parentCategory.childrenCount - maxCount} more
          </Badge>
        )}
      </div>
    </div>
  );
}

// ===== Analytics Display Component =====

interface AnalyticsDisplayProps {
  category: Category;
  compact?: boolean;
}

function AnalyticsDisplay({ category, compact = false }: AnalyticsDisplayProps) {
  if (compact) {
    return (
      <div className="flex items-center gap-2 text-xs text-muted-foreground">
        <div className="flex items-center gap-1">
          <Eye className="w-3 h-3" />
          <span>Popular</span>
        </div>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 gap-2 mt-3 pt-3 border-t border-border">
      <div className="text-center">
        <div className="text-lg font-semibold">2.1k</div>
        <div className="text-xs text-muted-foreground">Views</div>
      </div>
      <div className="text-center">
        <div className="text-lg font-semibold flex items-center justify-center gap-1">
          <TrendingUp className="w-3 h-3 text-green-500" />
          <span>12%</span>
        </div>
        <div className="text-xs text-muted-foreground">Growth</div>
      </div>
    </div>
  );
}

// ===== Loading Skeleton =====

export function CategoryCardSkeleton({ variant = 'default' }: { variant?: CategoryCardProps['variant'] }) {
  const isCompact = variant === 'compact' || variant === 'minimal';

  return (
    <Card className={cn(
      'overflow-hidden',
      isCompact ? 'h-20' : 'h-48'
    )}>
      <CardContent className={cn(
        'p-4',
        isCompact && 'flex items-center gap-3'
      )}>
        {!isCompact && (
          <Skeleton className="w-full h-24 mb-4 rounded" />
        )}

        {isCompact && (
          <Skeleton className="w-12 h-12 rounded" />
        )}

        <div className="flex-1 space-y-2">
          <Skeleton className="h-4 w-3/4" />
          {!isCompact && (
            <>
              <Skeleton className="h-3 w-full" />
              <Skeleton className="h-3 w-1/2" />
            </>
          )}
        </div>

        {!isCompact && (
          <div className="flex justify-between items-center mt-4">
            <Skeleton className="h-6 w-12" />
            <Skeleton className="h-8 w-8 rounded" />
          </div>
        )}
      </CardContent>
    </Card>
  );
}

// ===== Main CategoryCard Component =====

export function CategoryCard({
  category,
  variant = 'default',
  size = 'md',
  showDescription = true,
  showProductCount = true,
  showAnalytics = false,
  showSubcategories = false,
  maxSubcategories = 3,
  className,
  onClick,
  onSubcategoryClick,
}: CategoryCardProps) {
  const { setCurrentCategory, trackCategoryView } = useCategory();
  const [imageLoaded, setImageLoaded] = useState(false);
  const [imageError, setImageError] = useState(false);
  const cardRef = useRef<HTMLDivElement>(null);

  // Handle click
  const handleClick = () => {
    setCurrentCategory(category.id);
    trackCategoryView(category.id);
    onClick?.(category);
  };

  // Handle keyboard interaction
  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleClick();
    }
  };

  // Determine card dimensions based on variant and size
  const getCardClasses = () => {
    const baseClasses = 'group cursor-pointer overflow-hidden transition-all duration-200 hover:shadow-lg';

    switch (variant) {
      case 'compact':
        return cn(baseClasses, 'h-20');
      case 'minimal':
        return cn(baseClasses, 'h-16 border-none shadow-none');
      case 'featured':
        return cn(baseClasses, 'h-64 bg-gradient-to-br from-primary/5 to-primary/10 border-primary/20');
      case 'grid':
        return cn(baseClasses, 'aspect-square');
      default:
        return cn(baseClasses, size === 'sm' ? 'h-40' : size === 'lg' ? 'h-56' : 'h-48');
    }
  };

  // Generate category icon fallback
  const getCategoryIcon = () => {
    if (category.icon) {
      return category.icon;
    }
    return category.allowProducts ? '📦' : '📁';
  };

  // Render minimal variant
  if (variant === 'minimal') {
    return (
      <div
        ref={cardRef}
        role="button"
        tabIndex={0}
        aria-label={`Browse ${category.name} category`}
        className={cn(
          'flex items-center gap-3 p-3 rounded-lg hover:bg-muted/50 transition-colors',
          className
        )}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
      >
        <div
          className="w-8 h-8 rounded-full flex items-center justify-center text-sm"
          style={{ backgroundColor: `${category.color}20`, color: category.color }}
        >
          {getCategoryIcon()}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="font-medium truncate">{category.name}</h3>
          {showProductCount && category.activeProductCount > 0 && (
            <p className="text-xs text-muted-foreground">
              {category.activeProductCount} products
            </p>
          )}
        </div>
        <ChevronRight className="w-4 h-4 text-muted-foreground group-hover:text-foreground transition-colors" />
      </div>
    );
  }

  // Render compact variant
  if (variant === 'compact') {
    return (
      <Card
        ref={cardRef}
        role="button"
        tabIndex={0}
        aria-label={`Browse ${category.name} category`}
        className={cn(getCardClasses(), className)}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
      >
        <CardContent className="flex items-center gap-3 p-4 h-full">
          {/* Category Image/Icon */}
          <div className="relative flex-shrink-0">
            {category.image ? (
              <Avatar className="w-12 h-12">
                <AvatarImage
                  src={category.image.url}
                  alt={category.image.alt || category.name}
                  onLoad={() => setImageLoaded(true)}
                  onError={() => setImageError(true)}
                />
                <AvatarFallback style={{ backgroundColor: `${category.color}20`, color: category.color }}>
                  {getCategoryIcon()}
                </AvatarFallback>
              </Avatar>
            ) : (
              <div
                className="w-12 h-12 rounded-full flex items-center justify-center text-lg"
                style={{ backgroundColor: `${category.color}20`, color: category.color }}
              >
                {getCategoryIcon()}
              </div>
            )}

            {category.isFeatured && (
              <Star className="absolute -top-1 -right-1 w-4 h-4 text-yellow-500 fill-current" />
            )}
          </div>

          {/* Category Info */}
          <div className="flex-1 min-w-0">
            <h3 className="font-medium truncate group-hover:text-primary transition-colors">
              {category.name}
            </h3>
            {showProductCount && category.activeProductCount > 0 && (
              <p className="text-sm text-muted-foreground">
                {category.activeProductCount} products
              </p>
            )}
          </div>

          {/* Action Arrow */}
          <ChevronRight className="w-5 h-5 text-muted-foreground group-hover:text-primary group-hover:translate-x-1 transition-all" />
        </CardContent>
      </Card>
    );
  }

  // Render default, featured, or grid variant
  return (
    <Card
      ref={cardRef}
      role="button"
      tabIndex={0}
      aria-label={`Browse ${category.name} category`}
      className={cn(getCardClasses(), className)}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
    >
      <CardContent className="p-0 h-full flex flex-col">
        {/* Category Image */}
        <div className="relative flex-1 min-h-0">
          {category.image ? (
            <ImageWithFallback
              src={category.image.url}
              alt={category.image.alt || category.name}
              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              onLoad={() => setImageLoaded(true)}
              onError={() => setImageError(true)}
            />
          ) : (
            <div
              className="w-full h-full flex items-center justify-center text-4xl group-hover:scale-105 transition-transform duration-300"
              style={{ backgroundColor: `${category.color}10`, color: category.color }}
            >
              {getCategoryIcon()}
            </div>
          )}

          {/* Overlay Icons */}
          <div className="absolute top-3 right-3 flex gap-2">
            {category.isFeatured && (
              <div className="bg-yellow-500 text-white p-1 rounded-full">
                <Star className="w-3 h-3 fill-current" />
              </div>
            )}
          </div>

          {/* Gradient Overlay for better text readability */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>

        {/* Category Content */}
        <div className="p-4 space-y-3">
          {/* Category Header */}
          <div>
            <h3 className="font-semibold text-lg leading-tight group-hover:text-primary transition-colors line-clamp-2">
              {category.name}
            </h3>

            {showDescription && category.shortDescription && (
              <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                {category.shortDescription}
              </p>
            )}
          </div>

          {/* Category Metrics */}
          <div className="flex items-center justify-between">
            {showProductCount && (
              <div className="flex items-center gap-2">
                <Package className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm font-medium">
                  {category.activeProductCount} products
                </span>
              </div>
            )}

            {variant === 'featured' && showAnalytics && (
              <AnalyticsDisplay category={category} compact />
            )}
          </div>

          {/* Subcategories */}
          {showSubcategories && (
            <SubcategoryList
              parentCategory={category}
              maxCount={maxSubcategories}
              onSubcategoryClick={onSubcategoryClick}
            />
          )}

          {/* Analytics (for featured variant) */}
          {variant === 'featured' && showAnalytics && (
            <AnalyticsDisplay category={category} />
          )}
        </div>
      </CardContent>
    </Card>
  );
}

export default CategoryCard;