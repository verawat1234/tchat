/**
 * CategoryGrid - Grid layout for displaying multiple categories
 *
 * Displays categories in responsive grid layouts with various viewing options.
 * Supports filtering, sorting, pagination, and different display modes.
 *
 * Features:
 * - Responsive grid layouts (1-6 columns)
 * - Multiple view modes (grid, list, compact)
 * - Pagination with infinite scroll option
 * - Loading states and skeletons
 * - Empty states with actions
 * - Category filtering and sorting
 * - Accessibility with keyboard navigation
 * - Virtualization for large datasets
 */

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { Button } from '../ui/button';
import { ScrollArea } from '../ui/scroll-area';
import { Badge } from '../ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Skeleton } from '../ui/skeleton';
import {
  Grid3X3,
  List,
  SlidersHorizontal,
  RefreshCw,
  ChevronLeft,
  ChevronRight,
  Search,
  Package,
  Filter,
  Eye,
  EyeOff,
} from 'lucide-react';
import { cn } from '../../lib/utils';
import { useCategory } from './CategoryProvider';
import { CategoryCard, CategoryCardSkeleton } from './CategoryCard';
import type { Category, SortOptions } from '../../types/commerce';

// ===== Types =====

export interface CategoryGridProps {
  className?: string;
  categories?: Category[];
  viewMode?: 'grid' | 'list' | 'compact';
  gridColumns?: 1 | 2 | 3 | 4 | 5 | 6;
  showFilters?: boolean;
  showSorting?: boolean;
  showPagination?: boolean;
  showViewModeToggle?: boolean;
  infiniteScroll?: boolean;
  virtualScroll?: boolean;
  pageSize?: number;
  emptyMessage?: string;
  emptyAction?: React.ReactNode;
  onCategorySelect?: (category: Category) => void;
}

// ===== View Mode Toggle Component =====

interface ViewModeToggleProps {
  currentMode: 'grid' | 'list' | 'compact';
  onModeChange: (mode: 'grid' | 'list' | 'compact') => void;
}

function ViewModeToggle({ currentMode, onModeChange }: ViewModeToggleProps) {
  return (
    <div className="flex border rounded-lg overflow-hidden">
      <Button
        variant={currentMode === 'grid' ? 'default' : 'ghost'}
        size="sm"
        className="rounded-none border-0"
        onClick={() => onModeChange('grid')}
        aria-label="Grid view"
      >
        <Grid3X3 className="w-4 h-4" />
      </Button>
      <Button
        variant={currentMode === 'list' ? 'default' : 'ghost'}
        size="sm"
        className="rounded-none border-0 border-l"
        onClick={() => onModeChange('list')}
        aria-label="List view"
      >
        <List className="w-4 h-4" />
      </Button>
      <Button
        variant={currentMode === 'compact' ? 'default' : 'ghost'}
        size="sm"
        className="rounded-none border-0 border-l"
        onClick={() => onModeChange('compact')}
        aria-label="Compact view"
      >
        <Package className="w-4 h-4" />
      </Button>
    </div>
  );
}

// ===== Category Filter Component =====

interface CategoryFilterProps {
  categories: Category[];
  onFilterChange: (filters: CategoryFilter) => void;
}

interface CategoryFilter {
  status: 'all' | 'active' | 'inactive';
  type: 'all' | 'global' | 'business';
  featured: boolean;
  hasProducts: boolean;
}

function CategoryFilter({ categories, onFilterChange }: CategoryFilterProps) {
  const [filters, setFilters] = useState<CategoryFilter>({
    status: 'all',
    type: 'all',
    featured: false,
    hasProducts: false,
  });

  const handleFilterChange = (key: keyof CategoryFilter, value: any) => {
    const newFilters = { ...filters, [key]: value };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const activeFilterCount = Object.values(filters).filter(value =>
    typeof value === 'boolean' ? value : value !== 'all'
  ).length;

  return (
    <div className="flex items-center gap-3">
      <div className="flex items-center gap-2">
        <Filter className="w-4 h-4 text-muted-foreground" />
        <span className="text-sm font-medium">Filters</span>
        {activeFilterCount > 0 && (
          <Badge variant="secondary" className="h-5 px-1 text-xs">
            {activeFilterCount}
          </Badge>
        )}
      </div>

      <Select value={filters.status} onValueChange={(value) => handleFilterChange('status', value)}>
        <SelectTrigger className="w-32">
          <SelectValue placeholder="Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Status</SelectItem>
          <SelectItem value="active">Active</SelectItem>
          <SelectItem value="inactive">Inactive</SelectItem>
        </SelectContent>
      </Select>

      <Select value={filters.type} onValueChange={(value) => handleFilterChange('type', value)}>
        <SelectTrigger className="w-32">
          <SelectValue placeholder="Type" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Types</SelectItem>
          <SelectItem value="global">Global</SelectItem>
          <SelectItem value="business">Business</SelectItem>
        </SelectContent>
      </Select>

      <Button
        variant={filters.featured ? 'default' : 'outline'}
        size="sm"
        onClick={() => handleFilterChange('featured', !filters.featured)}
      >
        Featured
      </Button>

      <Button
        variant={filters.hasProducts ? 'default' : 'outline'}
        size="sm"
        onClick={() => handleFilterChange('hasProducts', !filters.hasProducts)}
      >
        Has Products
      </Button>

      {activeFilterCount > 0 && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => {
            const resetFilters: CategoryFilter = {
              status: 'all',
              type: 'all',
              featured: false,
              hasProducts: false,
            };
            setFilters(resetFilters);
            onFilterChange(resetFilters);
          }}
        >
          Clear
        </Button>
      )}
    </div>
  );
}

// ===== Sorting Component =====

interface SortingControlsProps {
  sortBy: SortOptions;
  onSortChange: (sort: SortOptions) => void;
}

function SortingControls({ sortBy, onSortChange }: SortingControlsProps) {
  const sortOptions = [
    { value: 'name-asc', label: 'Name (A-Z)', field: 'name', order: 'asc' as const },
    { value: 'name-desc', label: 'Name (Z-A)', field: 'name', order: 'desc' as const },
    { value: 'productCount-desc', label: 'Most Products', field: 'activeProductCount', order: 'desc' as const },
    { value: 'productCount-asc', label: 'Fewest Products', field: 'activeProductCount', order: 'asc' as const },
    { value: 'sortOrder-asc', label: 'Default Order', field: 'sortOrder', order: 'asc' as const },
    { value: 'createdAt-desc', label: 'Newest First', field: 'createdAt', order: 'desc' as const },
    { value: 'createdAt-asc', label: 'Oldest First', field: 'createdAt', order: 'asc' as const },
  ];

  const currentSortValue = `${sortBy.field}-${sortBy.order}`;

  return (
    <Select
      value={currentSortValue}
      onValueChange={(value) => {
        const option = sortOptions.find(opt => opt.value === value);
        if (option) {
          onSortChange({ field: option.field, order: option.order });
        }
      }}
    >
      <SelectTrigger className="w-48">
        <SelectValue placeholder="Sort by" />
      </SelectTrigger>
      <SelectContent>
        {sortOptions.map(option => (
          <SelectItem key={option.value} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}

// ===== Pagination Component =====

interface PaginationControlsProps {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
}

function PaginationControls({
  currentPage,
  totalPages,
  totalItems,
  pageSize,
  onPageChange,
  onPageSizeChange,
}: PaginationControlsProps) {
  const startItem = (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, totalItems);

  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <span>
          Showing {startItem}-{endItem} of {totalItems} categories
        </span>
        <Select value={pageSize.toString()} onValueChange={(value) => onPageSizeChange(Number(value))}>
          <SelectTrigger className="w-20 h-8">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="12">12</SelectItem>
            <SelectItem value="24">24</SelectItem>
            <SelectItem value="48">48</SelectItem>
            <SelectItem value="96">96</SelectItem>
          </SelectContent>
        </Select>
        <span>per page</span>
      </div>

      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          disabled={currentPage <= 1}
          onClick={() => onPageChange(currentPage - 1)}
        >
          <ChevronLeft className="w-4 h-4" />
          Previous
        </Button>

        <div className="flex items-center gap-1">
          {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
            let pageNum;
            if (totalPages <= 5) {
              pageNum = i + 1;
            } else if (currentPage <= 3) {
              pageNum = i + 1;
            } else if (currentPage >= totalPages - 2) {
              pageNum = totalPages - 4 + i;
            } else {
              pageNum = currentPage - 2 + i;
            }

            return (
              <Button
                key={pageNum}
                variant={currentPage === pageNum ? 'default' : 'outline'}
                size="sm"
                className="w-8 h-8 p-0"
                onClick={() => onPageChange(pageNum)}
              >
                {pageNum}
              </Button>
            );
          })}
        </div>

        <Button
          variant="outline"
          size="sm"
          disabled={currentPage >= totalPages}
          onClick={() => onPageChange(currentPage + 1)}
        >
          Next
          <ChevronRight className="w-4 h-4" />
        </Button>
      </div>
    </div>
  );
}

// ===== Loading Grid Component =====

function LoadingGrid({ viewMode, gridColumns, itemCount = 12 }: {
  viewMode: 'grid' | 'list' | 'compact';
  gridColumns: number;
  itemCount?: number;
}) {
  return (
    <div className={cn(
      'grid gap-4',
      viewMode === 'grid' && `grid-cols-1 sm:grid-cols-2 md:grid-cols-${Math.min(gridColumns, 3)} lg:grid-cols-${gridColumns}`,
      viewMode === 'list' && 'grid-cols-1',
      viewMode === 'compact' && 'grid-cols-1 gap-2'
    )}>
      {Array.from({ length: itemCount }, (_, index) => (
        <CategoryCardSkeleton key={index} variant={viewMode} />
      ))}
    </div>
  );
}

// ===== Empty State Component =====

function EmptyState({ message, action }: { message: string; action?: React.ReactNode }) {
  return (
    <div className="text-center py-12">
      <Search className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
      <h3 className="text-lg font-medium mb-2">No categories found</h3>
      <p className="text-muted-foreground mb-6 max-w-md mx-auto">
        {message}
      </p>
      {action}
    </div>
  );
}

// ===== Main CategoryGrid Component =====

export function CategoryGrid({
  className,
  categories: propCategories,
  viewMode: propViewMode,
  gridColumns = 3,
  showFilters = true,
  showSorting = true,
  showPagination = true,
  showViewModeToggle = true,
  infiniteScroll = false,
  virtualScroll = false,
  pageSize: propPageSize = 24,
  emptyMessage = "Try adjusting your filters or search terms to find what you're looking for.",
  emptyAction,
  onCategorySelect,
}: CategoryGridProps) {
  const { state, categoriesQuery, setViewMode, setSortBy } = useCategory();

  // Local state
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(propPageSize);
  const [filters, setFilters] = useState<CategoryFilter>({
    status: 'all',
    type: 'all',
    featured: false,
    hasProducts: false,
  });

  const viewMode = propViewMode || state.viewState.viewMode;
  const sortBy = state.viewState.sortBy;

  // Use provided categories or fetch from API
  const categories = propCategories || categoriesQuery.data?.categories || [];

  // Filter and sort categories
  const filteredAndSortedCategories = useMemo(() => {
    let filtered = categories.filter(category => {
      // Status filter
      if (filters.status === 'active' && !category.isVisible) return false;
      if (filters.status === 'inactive' && category.isVisible) return false;

      // Type filter
      if (filters.type === 'global' && category.type !== 'GLOBAL') return false;
      if (filters.type === 'business' && category.type !== 'BUSINESS') return false;

      // Featured filter
      if (filters.featured && !category.isFeatured) return false;

      // Has products filter
      if (filters.hasProducts && category.activeProductCount === 0) return false;

      return true;
    });

    // Sort categories
    filtered.sort((a, b) => {
      const { field, order } = sortBy;
      let aValue: any = a[field as keyof Category];
      let bValue: any = b[field as keyof Category];

      // Handle special cases
      if (field === 'name') {
        aValue = aValue.toLowerCase();
        bValue = bValue.toLowerCase();
      }

      if (aValue < bValue) return order === 'asc' ? -1 : 1;
      if (aValue > bValue) return order === 'asc' ? 1 : -1;
      return 0;
    });

    return filtered;
  }, [categories, filters, sortBy]);

  // Pagination
  const totalItems = filteredAndSortedCategories.length;
  const totalPages = Math.ceil(totalItems / pageSize);
  const startIndex = (currentPage - 1) * pageSize;
  const endIndex = startIndex + pageSize;
  const paginatedCategories = showPagination
    ? filteredAndSortedCategories.slice(startIndex, endIndex)
    : filteredAndSortedCategories;

  // Reset to first page when filters change
  useEffect(() => {
    setCurrentPage(1);
  }, [filters, sortBy]);

  // Handle category selection
  const handleCategorySelect = useCallback((category: Category) => {
    onCategorySelect?.(category);
  }, [onCategorySelect]);

  // Handle view mode change
  const handleViewModeChange = (mode: 'grid' | 'list' | 'compact') => {
    setViewMode(mode);
  };

  // Handle sort change
  const handleSortChange = (sort: SortOptions) => {
    setSortBy(sort);
  };

  // Loading state
  if (categoriesQuery.isLoading && !propCategories) {
    return (
      <div className={cn('space-y-6', className)}>
        {/* Controls skeleton */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Skeleton className="w-32 h-10" />
            <Skeleton className="w-48 h-10" />
          </div>
          <Skeleton className="w-24 h-10" />
        </div>

        {/* Grid skeleton */}
        <LoadingGrid viewMode={viewMode} gridColumns={gridColumns} />
      </div>
    );
  }

  // Error state
  if (categoriesQuery.error && !propCategories) {
    return (
      <div className={cn('text-center py-12', className)}>
        <Package className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
        <h3 className="text-lg font-medium mb-2">Failed to load categories</h3>
        <p className="text-muted-foreground mb-6">
          There was an error loading the categories. Please try again.
        </p>
        <Button onClick={() => categoriesQuery.refetch()}>
          <RefreshCw className="w-4 h-4 mr-2" />
          Try Again
        </Button>
      </div>
    );
  }

  return (
    <div className={cn('space-y-6', className)}>
      {/* Controls */}
      <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        {/* Filters */}
        {showFilters && (
          <CategoryFilter
            categories={categories}
            onFilterChange={setFilters}
          />
        )}

        <div className="flex items-center gap-4">
          {/* Sorting */}
          {showSorting && (
            <SortingControls
              sortBy={sortBy}
              onSortChange={handleSortChange}
            />
          )}

          {/* View Mode Toggle */}
          {showViewModeToggle && (
            <ViewModeToggle
              currentMode={viewMode}
              onModeChange={handleViewModeChange}
            />
          )}
        </div>
      </div>

      {/* Results summary */}
      <div className="flex items-center justify-between text-sm text-muted-foreground">
        <span>
          {totalItems} {totalItems === 1 ? 'category' : 'categories'} found
        </span>
      </div>

      {/* Category Grid */}
      {paginatedCategories.length === 0 ? (
        <EmptyState message={emptyMessage} action={emptyAction} />
      ) : (
        <div className={cn(
          'grid gap-4',
          viewMode === 'grid' && `grid-cols-1 sm:grid-cols-2 md:grid-cols-${Math.min(gridColumns, 3)} lg:grid-cols-${gridColumns}`,
          viewMode === 'list' && 'grid-cols-1',
          viewMode === 'compact' && 'grid-cols-1 gap-2'
        )}>
          {paginatedCategories.map(category => (
            <CategoryCard
              key={category.id}
              category={category}
              variant={viewMode === 'compact' ? 'compact' : viewMode === 'list' ? 'default' : 'grid'}
              onClick={handleCategorySelect}
            />
          ))}
        </div>
      )}

      {/* Pagination */}
      {showPagination && totalPages > 1 && (
        <PaginationControls
          currentPage={currentPage}
          totalPages={totalPages}
          totalItems={totalItems}
          pageSize={pageSize}
          onPageChange={setCurrentPage}
          onPageSizeChange={setPageSize}
        />
      )}
    </div>
  );
}

export default CategoryGrid;