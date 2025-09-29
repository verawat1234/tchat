/**
 * CategoryBreadcrumbs - Navigation breadcrumb component for category hierarchy
 *
 * Displays the current category path as interactive breadcrumbs with navigation.
 * Supports collapsing on mobile, custom separators, and accessibility features.
 *
 * Features:
 * - Hierarchical breadcrumb navigation
 * - Mobile-responsive collapsing
 * - Custom separators and styling
 * - Keyboard navigation support
 * - ARIA accessibility attributes
 * - Category icons and colors
 * - Loading and error states
 * - Quick navigation dropdown
 */

import React, { useState, useRef, useEffect } from 'react';
import { Button } from '../ui/button';
import { ScrollArea } from '../ui/scroll-area';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu';
import { Skeleton } from '../ui/skeleton';
import {
  ChevronRight,
  Home,
  MoreHorizontal,
  ChevronDown,
  Store,
  Folder,
  Package,
} from 'lucide-react';
import { cn } from '../../lib/utils';
import { useCategory } from './CategoryProvider';
import type { Category } from '../../types/commerce';

// ===== Types =====

export interface CategoryBreadcrumbsProps {
  className?: string;
  showHome?: boolean;
  showIcons?: boolean;
  maxItems?: number;
  separator?: 'chevron' | 'slash' | 'arrow' | 'dot';
  size?: 'sm' | 'md' | 'lg';
  variant?: 'default' | 'ghost' | 'outline';
  onNavigate?: (category: Category | null) => void;
}

interface BreadcrumbItem {
  id: string | null;
  name: string;
  category?: Category;
  icon?: React.ReactNode;
  isHome?: boolean;
}

// ===== Separator Component =====

function BreadcrumbSeparator({ type }: { type: 'chevron' | 'slash' | 'arrow' | 'dot' }) {
  const separators = {
    chevron: <ChevronRight className="w-4 h-4 text-muted-foreground" />,
    slash: <span className="text-muted-foreground mx-2">/</span>,
    arrow: <span className="text-muted-foreground mx-2">→</span>,
    dot: <span className="text-muted-foreground mx-2">•</span>,
  };

  return separators[type];
}

// ===== Breadcrumb Item Component =====

interface BreadcrumbItemProps {
  item: BreadcrumbItem;
  isLast: boolean;
  showIcons: boolean;
  size: 'sm' | 'md' | 'lg';
  variant: 'default' | 'ghost' | 'outline';
  onClick: (item: BreadcrumbItem) => void;
}

function BreadcrumbItemComponent({
  item,
  isLast,
  showIcons,
  size,
  variant,
  onClick,
}: BreadcrumbItemProps) {
  const buttonSize = size === 'sm' ? 'sm' : size === 'lg' ? 'lg' : 'default';

  const handleClick = () => {
    if (!isLast) {
      onClick(item);
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if ((event.key === 'Enter' || event.key === ' ') && !isLast) {
      event.preventDefault();
      onClick(item);
    }
  };

  return (
    <Button
      variant={isLast ? 'default' : variant}
      size={buttonSize}
      className={cn(
        'h-auto py-1 px-2 font-normal transition-colors',
        isLast
          ? 'cursor-default bg-primary text-primary-foreground'
          : 'hover:bg-muted cursor-pointer',
        size === 'sm' && 'text-xs',
        size === 'lg' && 'text-base'
      )}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      disabled={isLast}
      aria-current={isLast ? 'page' : undefined}
    >
      <div className="flex items-center gap-1.5">
        {showIcons && item.icon && (
          <span className="flex-shrink-0">
            {item.icon}
          </span>
        )}
        <span className="truncate max-w-32 sm:max-w-48">
          {item.name}
        </span>
      </div>
    </Button>
  );
}

// ===== Collapsed Breadcrumbs Component =====

interface CollapsedBreadcrumbsProps {
  items: BreadcrumbItem[];
  currentIndex: number;
  onClick: (item: BreadcrumbItem) => void;
}

function CollapsedBreadcrumbs({ items, currentIndex, onClick }: CollapsedBreadcrumbsProps) {
  const hiddenItems = items.slice(1, currentIndex);

  if (hiddenItems.length === 0) return null;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="h-auto py-1 px-2"
          aria-label={`Show ${hiddenItems.length} hidden breadcrumb items`}
        >
          <MoreHorizontal className="w-4 h-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" className="w-56">
        {hiddenItems.map((item) => (
          <DropdownMenuItem
            key={item.id || 'home'}
            onClick={() => onClick(item)}
            className="flex items-center gap-2"
          >
            {item.icon && (
              <span className="flex-shrink-0 w-4 h-4">
                {item.icon}
              </span>
            )}
            <span className="truncate">{item.name}</span>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

// ===== Quick Navigation Component =====

interface QuickNavigationProps {
  currentCategory: Category | null;
  onNavigate: (category: Category | null) => void;
}

function QuickNavigation({ currentCategory, onNavigate }: QuickNavigationProps) {
  const { categoriesQuery } = useCategory();

  const siblingCategories = currentCategory && categoriesQuery.data
    ? categoriesQuery.data.categories.filter(cat =>
        cat.parentId === currentCategory.parentId &&
        cat.id !== currentCategory.id &&
        cat.isVisible
      )
    : [];

  if (!currentCategory || siblingCategories.length === 0) return null;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="h-auto py-1 px-1 ml-1"
          aria-label="Navigate to sibling categories"
        >
          <ChevronDown className="w-3 h-3" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-48">
        {siblingCategories.map((category) => (
          <DropdownMenuItem
            key={category.id}
            onClick={() => onNavigate(category)}
            className="flex items-center gap-2"
          >
            {category.icon ? (
              <span style={{ color: category.color }}>
                {category.icon}
              </span>
            ) : (
              <Package className="w-4 h-4 text-muted-foreground" />
            )}
            <span className="truncate">{category.name}</span>
            {category.activeProductCount > 0 && (
              <span className="text-xs text-muted-foreground ml-auto">
                {category.activeProductCount}
              </span>
            )}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

// ===== Loading Skeleton =====

function BreadcrumbsSkeleton() {
  return (
    <div className="flex items-center gap-2">
      <Skeleton className="w-16 h-6" />
      <ChevronRight className="w-4 h-4 text-muted-foreground" />
      <Skeleton className="w-24 h-6" />
      <ChevronRight className="w-4 h-4 text-muted-foreground" />
      <Skeleton className="w-20 h-6" />
    </div>
  );
}

// ===== Main CategoryBreadcrumbs Component =====

export function CategoryBreadcrumbs({
  className,
  showHome = true,
  showIcons = true,
  maxItems = 5,
  separator = 'chevron',
  size = 'md',
  variant = 'ghost',
  onNavigate,
}: CategoryBreadcrumbsProps) {
  const { state, categoriesQuery, setCurrentCategory, trackCategoryView } = useCategory();
  const [isCollapsed, setIsCollapsed] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  const currentCategory = state.currentCategoryId
    ? categoriesQuery.data?.categories.find(cat => cat.id === state.currentCategoryId)
    : null;

  // Build breadcrumb items from current category path
  const breadcrumbItems = React.useMemo((): BreadcrumbItem[] => {
    const items: BreadcrumbItem[] = [];

    // Add home item
    if (showHome) {
      items.push({
        id: null,
        name: 'Store',
        icon: showIcons ? <Store className="w-4 h-4" /> : undefined,
        isHome: true,
      });
    }

    // Build path from current category
    if (currentCategory && categoriesQuery.data) {
      const path: Category[] = [];
      let current = currentCategory;

      // Build path from current to root
      while (current) {
        path.unshift(current);
        const parent = categoriesQuery.data.categories.find(cat => cat.id === current.parentId);
        current = parent || null;
      }

      // Add each category in the path
      path.forEach(category => {
        const icon = showIcons ? (
          category.icon ? (
            <span style={{ color: category.color }}>
              {category.icon}
            </span>
          ) : (
            <Folder className="w-4 h-4" style={{ color: category.color }} />
          )
        ) : undefined;

        items.push({
          id: category.id,
          name: category.name,
          category,
          icon,
        });
      });
    }

    return items;
  }, [currentCategory, categoriesQuery.data, showHome, showIcons]);

  // Handle navigation
  const handleNavigate = React.useCallback((item: BreadcrumbItem) => {
    if (item.isHome) {
      setCurrentCategory(null);
      onNavigate?.(null);
    } else if (item.category) {
      setCurrentCategory(item.category.id);
      trackCategoryView(item.category.id);
      onNavigate?.(item.category);
    }
  }, [setCurrentCategory, trackCategoryView, onNavigate]);

  // Handle responsive collapsing
  useEffect(() => {
    const handleResize = () => {
      if (containerRef.current) {
        const containerWidth = containerRef.current.offsetWidth;
        const shouldCollapse = containerWidth < 400 && breadcrumbItems.length > 3;
        setIsCollapsed(shouldCollapse);
      }
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [breadcrumbItems.length]);

  // Loading state
  if (categoriesQuery.isLoading) {
    return (
      <div className={cn('flex items-center', className)}>
        <BreadcrumbsSkeleton />
      </div>
    );
  }

  // Empty state (just home)
  if (breadcrumbItems.length <= 1) {
    return (
      <div className={cn('flex items-center', className)} ref={containerRef}>
        {breadcrumbItems.map((item) => (
          <BreadcrumbItemComponent
            key={item.id || 'home'}
            item={item}
            isLast={true}
            showIcons={showIcons}
            size={size}
            variant={variant}
            onClick={handleNavigate}
          />
        ))}
      </div>
    );
  }

  // Determine which items to show
  const visibleItems = React.useMemo(() => {
    if (!isCollapsed || breadcrumbItems.length <= maxItems) {
      return breadcrumbItems;
    }

    // Show first item, collapsed indicator, and last few items
    const lastItems = breadcrumbItems.slice(-2);
    return [breadcrumbItems[0], ...lastItems];
  }, [breadcrumbItems, isCollapsed, maxItems]);

  const shouldShowCollapsed = isCollapsed && breadcrumbItems.length > maxItems;

  return (
    <div className={cn('flex items-center flex-wrap gap-1', className)} ref={containerRef}>
      {visibleItems.map((item, index) => (
        <React.Fragment key={item.id || 'home'}>
          {/* Show collapsed items dropdown */}
          {shouldShowCollapsed && index === 1 && (
            <>
              <BreadcrumbSeparator type={separator} />
              <CollapsedBreadcrumbs
                items={breadcrumbItems}
                currentIndex={breadcrumbItems.length - 2}
                onClick={handleNavigate}
              />
            </>
          )}

          {/* Show separator */}
          {index > 0 && (!shouldShowCollapsed || index !== 1) && (
            <BreadcrumbSeparator type={separator} />
          )}

          {/* Show breadcrumb item */}
          <div className="flex items-center">
            <BreadcrumbItemComponent
              item={item}
              isLast={index === visibleItems.length - 1}
              showIcons={showIcons}
              size={size}
              variant={variant}
              onClick={handleNavigate}
            />

            {/* Show quick navigation for last item */}
            {index === visibleItems.length - 1 && item.category && (
              <QuickNavigation
                currentCategory={item.category}
                onNavigate={handleNavigate}
              />
            )}
          </div>
        </React.Fragment>
      ))}
    </div>
  );
}

// ===== Compact CategoryBreadcrumbs Component =====

interface CompactCategoryBreadcrumbsProps {
  className?: string;
  onNavigate?: (category: Category | null) => void;
}

export function CompactCategoryBreadcrumbs({
  className,
  onNavigate,
}: CompactCategoryBreadcrumbsProps) {
  return (
    <CategoryBreadcrumbs
      className={className}
      showHome={false}
      showIcons={false}
      maxItems={3}
      separator="slash"
      size="sm"
      variant="ghost"
      onNavigate={onNavigate}
    />
  );
}

export default CategoryBreadcrumbs;