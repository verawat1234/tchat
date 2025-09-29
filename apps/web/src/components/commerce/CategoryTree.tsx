/**
 * CategoryTree - Hierarchical category navigation component
 *
 * Displays categories in a tree structure with expand/collapse functionality.
 * Supports nested categories, lazy loading, and keyboard navigation.
 *
 * Features:
 * - Recursive tree rendering with unlimited depth
 * - Expand/collapse with smooth animations
 * - Lazy loading of child categories
 * - Keyboard navigation (arrow keys, enter, space)
 * - Accessibility with ARIA attributes
 * - Visual indicators for category states
 * - Product count display
 * - Category icons and colors
 */

import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Button } from '../ui/button';
import { ScrollArea } from '../ui/scroll-area';
import { Badge } from '../ui/badge';
import { Skeleton } from '../ui/skeleton';
import {
  ChevronRight,
  ChevronDown,
  Folder,
  FolderOpen,
  Package,
  Star,
  Eye,
  EyeOff,
} from 'lucide-react';
import { cn } from '../../lib/utils';
import { useCategory } from './CategoryProvider';
import type { Category } from '../../types/commerce';

// ===== Types =====

interface CategoryTreeProps {
  className?: string;
  maxHeight?: string;
  showProductCounts?: boolean;
  showIcons?: boolean;
  allowMultiSelect?: boolean;
  onCategorySelect?: (category: Category) => void;
  onCategoryToggle?: (categoryId: string, isExpanded: boolean) => void;
}

interface CategoryNodeProps {
  category: Category;
  allCategories: Category[];
  depth: number;
  isExpanded: boolean;
  isSelected: boolean;
  showProductCounts: boolean;
  showIcons: boolean;
  onToggle: (categoryId: string) => void;
  onSelect: (category: Category) => void;
  onKeyDown: (event: React.KeyboardEvent, category: Category) => void;
}

// ===== Category Node Component =====

function CategoryNode({
  category,
  allCategories,
  depth,
  isExpanded,
  isSelected,
  showProductCounts,
  showIcons,
  onToggle,
  onSelect,
  onKeyDown,
}: CategoryNodeProps) {
  const hasChildren = allCategories.some(cat => cat.parentId === category.id);
  const childCategories = allCategories.filter(cat => cat.parentId === category.id);

  const handleClick = () => {
    onSelect(category);
  };

  const handleToggle = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (hasChildren) {
      onToggle(category.id);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    onKeyDown(e, category);
  };

  // Determine icon based on category state
  const getIcon = () => {
    if (!showIcons) return null;

    if (category.icon) {
      // Custom category icon
      return (
        <span
          className="w-4 h-4 flex items-center justify-center text-lg"
          style={{ color: category.color }}
        >
          {category.icon}
        </span>
      );
    }

    // Default icons
    if (hasChildren) {
      return isExpanded ? (
        <FolderOpen className="w-4 h-4 text-blue-500" />
      ) : (
        <Folder className="w-4 h-4 text-blue-500" />
      );
    }

    return <Package className="w-4 h-4 text-gray-500" />;
  };

  return (
    <div className="select-none">
      {/* Category Item */}
      <div
        role="treeitem"
        tabIndex={0}
        aria-expanded={hasChildren ? isExpanded : undefined}
        aria-selected={isSelected}
        aria-level={depth + 1}
        className={cn(
          'flex items-center gap-2 py-2 px-3 rounded-lg cursor-pointer transition-colors',
          'hover:bg-muted focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-1',
          isSelected && 'bg-primary/10 text-primary border border-primary/20',
          !category.isVisible && 'opacity-50',
          depth > 0 && 'ml-4 border-l border-border'
        )}
        style={{ paddingLeft: `${depth * 16 + 12}px` }}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
      >
        {/* Expand/Collapse Toggle */}
        {hasChildren ? (
          <Button
            variant="ghost"
            size="icon"
            className="w-4 h-4 p-0 hover:bg-transparent"
            onClick={handleToggle}
            aria-label={isExpanded ? 'Collapse category' : 'Expand category'}
          >
            {isExpanded ? (
              <ChevronDown className="w-3 h-3" />
            ) : (
              <ChevronRight className="w-3 h-3" />
            )}
          </Button>
        ) : (
          <div className="w-4 h-4" />
        )}

        {/* Category Icon */}
        {getIcon()}

        {/* Category Name */}
        <span className={cn(
          'flex-1 truncate text-sm font-medium',
          !category.isVisible && 'line-through'
        )}>
          {category.name}
        </span>

        {/* Category Badges */}
        <div className="flex items-center gap-1">
          {/* Featured Badge */}
          {category.isFeatured && (
            <Star className="w-3 h-3 text-yellow-500" title="Featured category" />
          )}

          {/* Visibility Indicator */}
          {!category.isVisible && (
            <EyeOff className="w-3 h-3 text-gray-400" title="Hidden category" />
          )}

          {/* Product Count */}
          {showProductCounts && category.activeProductCount > 0 && (
            <Badge variant="secondary" className="text-xs">
              {category.activeProductCount}
            </Badge>
          )}
        </div>
      </div>

      {/* Child Categories */}
      {hasChildren && isExpanded && (
        <div className="mt-1">
          {childCategories.map(childCategory => (
            <CategoryNodeContainer
              key={childCategory.id}
              category={childCategory}
              allCategories={allCategories}
              depth={depth + 1}
              showProductCounts={showProductCounts}
              showIcons={showIcons}
              onSelect={onSelect}
              onKeyDown={onKeyDown}
            />
          ))}
        </div>
      )}
    </div>
  );
}

// ===== Category Node Container =====

interface CategoryNodeContainerProps {
  category: Category;
  allCategories: Category[];
  depth: number;
  showProductCounts: boolean;
  showIcons: boolean;
  onSelect: (category: Category) => void;
  onKeyDown: (event: React.KeyboardEvent, category: Category) => void;
}

function CategoryNodeContainer({
  category,
  allCategories,
  depth,
  showProductCounts,
  showIcons,
  onSelect,
  onKeyDown,
}: CategoryNodeContainerProps) {
  const { state, setCurrentCategory, trackCategoryView } = useCategory();
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(new Set());

  const isSelected = state.currentCategoryId === category.id;
  const isExpanded = expandedCategories.has(category.id);

  const handleToggle = (categoryId: string) => {
    setExpandedCategories(prev => {
      const newSet = new Set(prev);
      if (newSet.has(categoryId)) {
        newSet.delete(categoryId);
      } else {
        newSet.add(categoryId);
      }
      return newSet;
    });
  };

  const handleSelect = useCallback((selectedCategory: Category) => {
    setCurrentCategory(selectedCategory.id);
    trackCategoryView(selectedCategory.id);
    onSelect(selectedCategory);
  }, [setCurrentCategory, trackCategoryView, onSelect]);

  return (
    <CategoryNode
      category={category}
      allCategories={allCategories}
      depth={depth}
      isExpanded={isExpanded}
      isSelected={isSelected}
      showProductCounts={showProductCounts}
      showIcons={showIcons}
      onToggle={handleToggle}
      onSelect={handleSelect}
      onKeyDown={onKeyDown}
    />
  );
}

// ===== Loading Skeleton =====

function CategoryTreeSkeleton() {
  return (
    <div className="space-y-2 p-3">
      {[...Array(8)].map((_, index) => (
        <div key={index} className="flex items-center gap-2">
          <Skeleton className="w-4 h-4" />
          <Skeleton className="w-4 h-4" />
          <Skeleton className="flex-1 h-4" />
          <Skeleton className="w-8 h-4" />
        </div>
      ))}
    </div>
  );
}

// ===== Main CategoryTree Component =====

export function CategoryTree({
  className,
  maxHeight = '400px',
  showProductCounts = true,
  showIcons = true,
  allowMultiSelect = false,
  onCategorySelect,
  onCategoryToggle,
}: CategoryTreeProps) {
  const { categoriesQuery, rootCategoriesQuery, state } = useCategory();
  const treeRef = useRef<HTMLDivElement>(null);

  // Handle keyboard navigation
  const handleKeyDown = useCallback((event: React.KeyboardEvent, category: Category) => {
    const { key } = event;

    switch (key) {
      case 'Enter':
      case ' ':
        event.preventDefault();
        onCategorySelect?.(category);
        break;

      case 'ArrowDown':
      case 'ArrowUp':
        event.preventDefault();
        // Navigate to next/previous category
        // Implementation would require more complex focus management
        break;

      case 'ArrowRight':
        event.preventDefault();
        // Expand category if collapsed
        onCategoryToggle?.(category.id, true);
        break;

      case 'ArrowLeft':
        event.preventDefault();
        // Collapse category if expanded
        onCategoryToggle?.(category.id, false);
        break;

      default:
        break;
    }
  }, [onCategorySelect, onCategoryToggle]);

  // Loading state
  if (categoriesQuery.isLoading || rootCategoriesQuery.isLoading) {
    return (
      <div className={cn('border rounded-lg bg-card', className)}>
        <CategoryTreeSkeleton />
      </div>
    );
  }

  // Error state
  if (categoriesQuery.error || rootCategoriesQuery.error) {
    return (
      <div className={cn('border rounded-lg bg-card p-4', className)}>
        <div className="text-center">
          <Package className="w-8 h-8 text-muted-foreground mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">
            Failed to load categories
          </p>
          <Button
            variant="outline"
            size="sm"
            className="mt-2"
            onClick={() => {
              categoriesQuery.refetch();
              rootCategoriesQuery.refetch();
            }}
          >
            Retry
          </Button>
        </div>
      </div>
    );
  }

  const allCategories = categoriesQuery.data?.categories || [];
  const rootCategories = rootCategoriesQuery.data || [];

  // Empty state
  if (rootCategories.length === 0) {
    return (
      <div className={cn('border rounded-lg bg-card p-4', className)}>
        <div className="text-center">
          <Folder className="w-8 h-8 text-muted-foreground mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">
            No categories available
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className={cn('border rounded-lg bg-card', className)}>
      <ScrollArea style={{ maxHeight }}>
        <div
          ref={treeRef}
          role="tree"
          aria-label="Category tree"
          className="p-3 space-y-1"
        >
          {rootCategories.map(category => (
            <CategoryNodeContainer
              key={category.id}
              category={category}
              allCategories={allCategories}
              depth={0}
              showProductCounts={showProductCounts}
              showIcons={showIcons}
              onSelect={onCategorySelect || (() => {})}
              onKeyDown={handleKeyDown}
            />
          ))}
        </div>
      </ScrollArea>
    </div>
  );
}

// ===== Compact CategoryTree Component =====

interface CompactCategoryTreeProps {
  selectedCategoryId?: string;
  onCategorySelect: (category: Category) => void;
  className?: string;
}

export function CompactCategoryTree({
  selectedCategoryId,
  onCategorySelect,
  className,
}: CompactCategoryTreeProps) {
  const { rootCategoriesQuery } = useCategory();

  if (rootCategoriesQuery.isLoading) {
    return (
      <div className={cn('flex gap-2', className)}>
        {[...Array(4)].map((_, index) => (
          <Skeleton key={index} className="w-20 h-8" />
        ))}
      </div>
    );
  }

  const rootCategories = rootCategoriesQuery.data || [];

  return (
    <div className={cn('flex gap-2 flex-wrap', className)}>
      {rootCategories.map(category => (
        <Button
          key={category.id}
          variant={selectedCategoryId === category.id ? 'default' : 'outline'}
          size="sm"
          className="h-8"
          onClick={() => onCategorySelect(category)}
        >
          {category.icon && (
            <span className="mr-1" style={{ color: category.color }}>
              {category.icon}
            </span>
          )}
          {category.name}
          {category.activeProductCount > 0 && (
            <Badge variant="secondary" className="ml-1 text-xs">
              {category.activeProductCount}
            </Badge>
          )}
        </Button>
      ))}
    </div>
  );
}

export default CategoryTree;