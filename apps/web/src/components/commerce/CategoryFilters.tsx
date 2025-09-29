/**
 * CategoryFilters - Advanced filtering component for categories and products
 *
 * Provides comprehensive filtering options for category browsing and product discovery.
 * Supports price ranges, ratings, attributes, and custom filters.
 *
 * Features:
 * - Multiple filter types (price, rating, attributes, etc.)
 * - Real-time filter application
 * - Filter presets and saved filters
 * - Mobile-responsive design
 * - Filter count indicators
 * - Clear all functionality
 * - Advanced filter combinations
 * - Accessibility support
 */

import React, { useState, useEffect, useMemo } from 'react';
import { Button } from '../ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Badge } from '../ui/badge';
import { Label } from '../ui/label';
import { Checkbox } from '../ui/checkbox';
import { Slider } from '../ui/slider';
import { ScrollArea } from '../ui/scroll-area';
import { Separator } from '../ui/separator';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '../ui/collapsible';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../ui/select';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '../ui/sheet';
import {
  Filter,
  X,
  ChevronDown,
  ChevronUp,
  Star,
  Tag,
  Package,
  DollarSign,
  Palette,
  Zap,
  SlidersHorizontal,
  RotateCcw,
  Save,
  Bookmark,
} from 'lucide-react';
import { cn } from '../../lib/utils';
import { useCategory } from './CategoryProvider';
import type { Category, Product } from '../../types/commerce';

// ===== Types =====

export interface CategoryFiltersProps {
  className?: string;
  variant?: 'sidebar' | 'inline' | 'modal';
  showPresets?: boolean;
  showSaveFilter?: boolean;
  onFiltersChange?: (filters: FilterState) => void;
}

export interface FilterState {
  priceRange: [number, number];
  rating: number;
  inStock: boolean;
  onSale: boolean;
  featured: boolean;
  freeShipping: boolean;
  brands: string[];
  colors: string[];
  sizes: string[];
  categories: string[];
  attributes: Record<string, string[]>;
  sortBy: string;
}

interface FilterPreset {
  id: string;
  name: string;
  filters: Partial<FilterState>;
  icon?: React.ReactNode;
}

// ===== Filter Section Component =====

interface FilterSectionProps {
  title: string;
  icon?: React.ReactNode;
  children: React.ReactNode;
  defaultExpanded?: boolean;
  badge?: string | number;
}

function FilterSection({ title, icon, children, defaultExpanded = true, badge }: FilterSectionProps) {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded);

  return (
    <Collapsible open={isExpanded} onOpenChange={setIsExpanded}>
      <CollapsibleTrigger asChild>
        <Button
          variant="ghost"
          className="w-full justify-between p-3 h-auto font-medium"
        >
          <div className="flex items-center gap-2">
            {icon}
            <span>{title}</span>
            {badge && (
              <Badge variant="secondary" className="h-5 px-1 text-xs">
                {badge}
              </Badge>
            )}
          </div>
          {isExpanded ? (
            <ChevronUp className="h-4 w-4" />
          ) : (
            <ChevronDown className="h-4 w-4" />
          )}
        </Button>
      </CollapsibleTrigger>
      <CollapsibleContent className="px-3 pb-3">
        {children}
      </CollapsibleContent>
      <Separator />
    </Collapsible>
  );
}

// ===== Price Range Filter =====

interface PriceRangeFilterProps {
  value: [number, number];
  onChange: (value: [number, number]) => void;
  min?: number;
  max?: number;
  step?: number;
}

function PriceRangeFilter({
  value,
  onChange,
  min = 0,
  max = 10000,
  step = 100
}: PriceRangeFilterProps) {
  const [localValue, setLocalValue] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => onChange(localValue), 300);
    return () => clearTimeout(timer);
  }, [localValue, onChange]);

  return (
    <div className="space-y-4">
      <div className="flex justify-between text-sm">
        <span>฿{localValue[0].toLocaleString()}</span>
        <span>฿{localValue[1].toLocaleString()}</span>
      </div>
      <Slider
        value={localValue}
        onValueChange={setLocalValue}
        min={min}
        max={max}
        step={step}
        className="w-full"
      />
      <div className="flex gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setLocalValue([0, 1000])}
          className="text-xs"
        >
          Under ฿1K
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setLocalValue([1000, 5000])}
          className="text-xs"
        >
          ฿1K-5K
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setLocalValue([5000, max])}
          className="text-xs"
        >
          Over ฿5K
        </Button>
      </div>
    </div>
  );
}

// ===== Rating Filter =====

interface RatingFilterProps {
  value: number;
  onChange: (value: number) => void;
}

function RatingFilter({ value, onChange }: RatingFilterProps) {
  const ratings = [
    { value: 4, label: '4+ stars', count: 1250 },
    { value: 3, label: '3+ stars', count: 2100 },
    { value: 2, label: '2+ stars', count: 2800 },
    { value: 1, label: '1+ stars', count: 3200 },
  ];

  return (
    <div className="space-y-2">
      {ratings.map((rating) => (
        <label
          key={rating.value}
          className="flex items-center gap-3 cursor-pointer hover:bg-muted/50 p-2 rounded"
        >
          <input
            type="radio"
            name="rating"
            checked={value === rating.value}
            onChange={() => onChange(rating.value)}
            className="sr-only"
          />
          <div className="flex items-center gap-1">
            {Array.from({ length: 5 }, (_, i) => (
              <Star
                key={i}
                className={cn(
                  'w-3 h-3',
                  i < rating.value ? 'text-yellow-500 fill-current' : 'text-gray-300'
                )}
              />
            ))}
          </div>
          <span className="text-sm">{rating.label}</span>
          <span className="text-xs text-muted-foreground ml-auto">
            ({rating.count.toLocaleString()})
          </span>
        </label>
      ))}
      {value > 0 && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onChange(0)}
          className="text-xs text-muted-foreground"
        >
          Clear rating filter
        </Button>
      )}
    </div>
  );
}

// ===== Checkbox Group Filter =====

interface CheckboxGroupFilterProps {
  title: string;
  options: Array<{ value: string; label: string; count?: number }>;
  selectedValues: string[];
  onChange: (values: string[]) => void;
  maxVisible?: number;
}

function CheckboxGroupFilter({
  title,
  options,
  selectedValues,
  onChange,
  maxVisible = 5,
}: CheckboxGroupFilterProps) {
  const [showAll, setShowAll] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  const filteredOptions = useMemo(() => {
    return options.filter(option =>
      option.label.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [options, searchQuery]);

  const visibleOptions = showAll
    ? filteredOptions
    : filteredOptions.slice(0, maxVisible);

  const handleToggle = (value: string) => {
    const newValues = selectedValues.includes(value)
      ? selectedValues.filter(v => v !== value)
      : [...selectedValues, value];
    onChange(newValues);
  };

  return (
    <div className="space-y-3">
      {options.length > maxVisible && (
        <div className="relative">
          <input
            type="text"
            placeholder={`Search ${title.toLowerCase()}...`}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full px-3 py-2 text-sm border rounded-md"
          />
        </div>
      )}

      <div className="space-y-1 max-h-48 overflow-y-auto">
        {visibleOptions.map((option) => (
          <label
            key={option.value}
            className="flex items-center gap-3 cursor-pointer hover:bg-muted/50 p-2 rounded"
          >
            <Checkbox
              checked={selectedValues.includes(option.value)}
              onCheckedChange={() => handleToggle(option.value)}
            />
            <span className="text-sm flex-1">{option.label}</span>
            {option.count && (
              <span className="text-xs text-muted-foreground">
                ({option.count.toLocaleString()})
              </span>
            )}
          </label>
        ))}
      </div>

      {filteredOptions.length > maxVisible && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setShowAll(!showAll)}
          className="text-xs"
        >
          {showAll ? 'Show less' : `Show all ${filteredOptions.length}`}
        </Button>
      )}

      {selectedValues.length > 0 && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onChange([])}
          className="text-xs text-muted-foreground"
        >
          Clear {title.toLowerCase()}
        </Button>
      )}
    </div>
  );
}

// ===== Filter Presets =====

const filterPresets: FilterPreset[] = [
  {
    id: 'popular',
    name: 'Popular Items',
    icon: <TrendingUp className="w-4 h-4" />,
    filters: {
      rating: 4,
      featured: true,
    },
  },
  {
    id: 'budget',
    name: 'Budget Friendly',
    icon: <DollarSign className="w-4 h-4" />,
    filters: {
      priceRange: [0, 1000],
      freeShipping: true,
    },
  },
  {
    id: 'premium',
    name: 'Premium Quality',
    icon: <Star className="w-4 h-4" />,
    filters: {
      rating: 4.5,
      priceRange: [2000, 10000],
    },
  },
  {
    id: 'sale',
    name: 'On Sale',
    icon: <Tag className="w-4 h-4" />,
    filters: {
      onSale: true,
    },
  },
];

interface FilterPresetsProps {
  onPresetSelect: (preset: FilterPreset) => void;
}

function FilterPresets({ onPresetSelect }: FilterPresetsProps) {
  return (
    <div className="grid grid-cols-2 gap-2">
      {filterPresets.map((preset) => (
        <Button
          key={preset.id}
          variant="outline"
          size="sm"
          onClick={() => onPresetSelect(preset)}
          className="justify-start gap-2 h-auto p-3"
        >
          {preset.icon}
          <span className="text-xs">{preset.name}</span>
        </Button>
      ))}
    </div>
  );
}

// ===== Active Filters Display =====

interface ActiveFiltersProps {
  filters: FilterState;
  onRemoveFilter: (key: keyof FilterState, value?: string) => void;
  onClearAll: () => void;
}

function ActiveFilters({ filters, onRemoveFilter, onClearAll }: ActiveFiltersProps) {
  const activeFilters = useMemo(() => {
    const active: Array<{ key: keyof FilterState; label: string; value?: string }> = [];

    if (filters.priceRange[0] > 0 || filters.priceRange[1] < 10000) {
      active.push({
        key: 'priceRange',
        label: `฿${filters.priceRange[0].toLocaleString()} - ฿${filters.priceRange[1].toLocaleString()}`,
      });
    }

    if (filters.rating > 0) {
      active.push({
        key: 'rating',
        label: `${filters.rating}+ stars`,
      });
    }

    if (filters.inStock) active.push({ key: 'inStock', label: 'In Stock' });
    if (filters.onSale) active.push({ key: 'onSale', label: 'On Sale' });
    if (filters.featured) active.push({ key: 'featured', label: 'Featured' });
    if (filters.freeShipping) active.push({ key: 'freeShipping', label: 'Free Shipping' });

    filters.brands.forEach(brand => {
      active.push({ key: 'brands', label: brand, value: brand });
    });

    filters.colors.forEach(color => {
      active.push({ key: 'colors', label: color, value: color });
    });

    filters.sizes.forEach(size => {
      active.push({ key: 'sizes', label: size, value: size });
    });

    return active;
  }, [filters]);

  if (activeFilters.length === 0) return null;

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-medium">Active Filters</h4>
        <Button
          variant="ghost"
          size="sm"
          onClick={onClearAll}
          className="text-xs text-muted-foreground"
        >
          Clear All
        </Button>
      </div>
      <div className="flex flex-wrap gap-2">
        {activeFilters.map((filter, index) => (
          <Badge
            key={`${filter.key}-${filter.value || index}`}
            variant="secondary"
            className="gap-1 py-1"
          >
            <span>{filter.label}</span>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => onRemoveFilter(filter.key, filter.value)}
              className="w-3 h-3 p-0 hover:bg-transparent"
            >
              <X className="w-2 h-2" />
            </Button>
          </Badge>
        ))}
      </div>
    </div>
  );
}

// ===== Main CategoryFilters Component =====

export function CategoryFilters({
  className,
  variant = 'sidebar',
  showPresets = true,
  showSaveFilter = false,
  onFiltersChange,
}: CategoryFiltersProps) {
  const { state, updateFilters, clearFilters } = useCategory();
  const [localFilters, setLocalFilters] = useState<FilterState>({
    priceRange: [0, 10000],
    rating: 0,
    inStock: false,
    onSale: false,
    featured: false,
    freeShipping: false,
    brands: [],
    colors: [],
    sizes: [],
    categories: [],
    attributes: {},
    sortBy: 'relevance',
  });

  // Mock data for filter options
  const brandOptions = [
    { value: 'apple', label: 'Apple', count: 245 },
    { value: 'samsung', label: 'Samsung', count: 189 },
    { value: 'sony', label: 'Sony', count: 156 },
    { value: 'lg', label: 'LG', count: 123 },
    { value: 'nike', label: 'Nike', count: 298 },
    { value: 'adidas', label: 'Adidas', count: 267 },
  ];

  const colorOptions = [
    { value: 'black', label: 'Black', count: 567 },
    { value: 'white', label: 'White', count: 423 },
    { value: 'blue', label: 'Blue', count: 312 },
    { value: 'red', label: 'Red', count: 289 },
    { value: 'green', label: 'Green', count: 234 },
    { value: 'gray', label: 'Gray', count: 198 },
  ];

  const sizeOptions = [
    { value: 'xs', label: 'XS', count: 89 },
    { value: 's', label: 'S', count: 156 },
    { value: 'm', label: 'M', count: 234 },
    { value: 'l', label: 'L', count: 198 },
    { value: 'xl', label: 'XL', count: 145 },
    { value: 'xxl', label: 'XXL', count: 87 },
  ];

  const handleFilterChange = <K extends keyof FilterState>(
    key: K,
    value: FilterState[K]
  ) => {
    const newFilters = { ...localFilters, [key]: value };
    setLocalFilters(newFilters);
    onFiltersChange?.(newFilters);
  };

  const handleRemoveFilter = (key: keyof FilterState, value?: string) => {
    let newFilters = { ...localFilters };

    if (value && Array.isArray(newFilters[key])) {
      newFilters[key] = (newFilters[key] as string[]).filter(v => v !== value);
    } else if (key === 'priceRange') {
      newFilters[key] = [0, 10000];
    } else if (typeof newFilters[key] === 'boolean') {
      newFilters[key] = false as any;
    } else if (typeof newFilters[key] === 'number') {
      newFilters[key] = 0 as any;
    }

    setLocalFilters(newFilters);
    onFiltersChange?.(newFilters);
  };

  const handleClearAll = () => {
    const resetFilters: FilterState = {
      priceRange: [0, 10000],
      rating: 0,
      inStock: false,
      onSale: false,
      featured: false,
      freeShipping: false,
      brands: [],
      colors: [],
      sizes: [],
      categories: [],
      attributes: {},
      sortBy: 'relevance',
    };

    setLocalFilters(resetFilters);
    onFiltersChange?.(resetFilters);
  };

  const handlePresetSelect = (preset: FilterPreset) => {
    const newFilters = { ...localFilters, ...preset.filters };
    setLocalFilters(newFilters);
    onFiltersChange?.(newFilters);
  };

  const filterContent = (
    <div className="space-y-1">
      {/* Active Filters */}
      <div className="p-3">
        <ActiveFilters
          filters={localFilters}
          onRemoveFilter={handleRemoveFilter}
          onClearAll={handleClearAll}
        />
      </div>

      {/* Filter Presets */}
      {showPresets && (
        <FilterSection title="Quick Filters" icon={<Zap className="w-4 h-4" />}>
          <FilterPresets onPresetSelect={handlePresetSelect} />
        </FilterSection>
      )}

      {/* Price Range */}
      <FilterSection title="Price Range" icon={<DollarSign className="w-4 h-4" />}>
        <PriceRangeFilter
          value={localFilters.priceRange}
          onChange={(value) => handleFilterChange('priceRange', value)}
        />
      </FilterSection>

      {/* Rating */}
      <FilterSection title="Customer Rating" icon={<Star className="w-4 h-4" />}>
        <RatingFilter
          value={localFilters.rating}
          onChange={(value) => handleFilterChange('rating', value)}
        />
      </FilterSection>

      {/* Availability & Features */}
      <FilterSection title="Availability" icon={<Package className="w-4 h-4" />}>
        <div className="space-y-3">
          {[
            { key: 'inStock', label: 'In Stock' },
            { key: 'onSale', label: 'On Sale' },
            { key: 'featured', label: 'Featured' },
            { key: 'freeShipping', label: 'Free Shipping' },
          ].map((item) => (
            <label key={item.key} className="flex items-center gap-2 cursor-pointer">
              <Checkbox
                checked={localFilters[item.key as keyof FilterState] as boolean}
                onCheckedChange={(checked) =>
                  handleFilterChange(item.key as keyof FilterState, checked as any)
                }
              />
              <span className="text-sm">{item.label}</span>
            </label>
          ))}
        </div>
      </FilterSection>

      {/* Brands */}
      <FilterSection
        title="Brands"
        icon={<Tag className="w-4 h-4" />}
        badge={localFilters.brands.length || undefined}
      >
        <CheckboxGroupFilter
          title="Brands"
          options={brandOptions}
          selectedValues={localFilters.brands}
          onChange={(values) => handleFilterChange('brands', values)}
        />
      </FilterSection>

      {/* Colors */}
      <FilterSection
        title="Colors"
        icon={<Palette className="w-4 h-4" />}
        badge={localFilters.colors.length || undefined}
      >
        <CheckboxGroupFilter
          title="Colors"
          options={colorOptions}
          selectedValues={localFilters.colors}
          onChange={(values) => handleFilterChange('colors', values)}
        />
      </FilterSection>

      {/* Sizes */}
      <FilterSection
        title="Sizes"
        icon={<Package className="w-4 h-4" />}
        badge={localFilters.sizes.length || undefined}
      >
        <CheckboxGroupFilter
          title="Sizes"
          options={sizeOptions}
          selectedValues={localFilters.sizes}
          onChange={(values) => handleFilterChange('sizes', values)}
        />
      </FilterSection>

      {/* Save Filter */}
      {showSaveFilter && (
        <div className="p-3">
          <Button variant="outline" size="sm" className="w-full gap-2">
            <Save className="w-4 h-4" />
            Save Current Filters
          </Button>
        </div>
      )}
    </div>
  );

  // Mobile sheet variant
  if (variant === 'modal') {
    return (
      <Sheet>
        <SheetTrigger asChild>
          <Button variant="outline" className={cn('gap-2', className)}>
            <Filter className="w-4 h-4" />
            Filters
            {Object.values(localFilters).some(value =>
              Array.isArray(value) ? value.length > 0 :
              typeof value === 'boolean' ? value :
              typeof value === 'number' ? value > 0 : false
            ) && (
              <Badge variant="secondary" className="h-5 px-1 text-xs">
                Active
              </Badge>
            )}
          </Button>
        </SheetTrigger>
        <SheetContent side="left" className="w-80">
          <SheetHeader>
            <SheetTitle>Filters</SheetTitle>
            <SheetDescription>
              Refine your search results
            </SheetDescription>
          </SheetHeader>
          <ScrollArea className="h-full pr-4">
            {filterContent}
          </ScrollArea>
        </SheetContent>
      </Sheet>
    );
  }

  // Inline variant
  if (variant === 'inline') {
    return (
      <Card className={className}>
        <CardHeader className="pb-2">
          <CardTitle className="text-base flex items-center gap-2">
            <SlidersHorizontal className="w-4 h-4" />
            Filters
          </CardTitle>
        </CardHeader>
        <CardContent className="pt-0">
          <ScrollArea className="h-96">
            {filterContent}
          </ScrollArea>
        </CardContent>
      </Card>
    );
  }

  // Sidebar variant (default)
  return (
    <div className={cn('w-80 space-y-1', className)}>
      <div className="flex items-center justify-between p-3 border-b">
        <h3 className="font-medium flex items-center gap-2">
          <SlidersHorizontal className="w-4 h-4" />
          Filters
        </h3>
        <Button
          variant="ghost"
          size="sm"
          onClick={handleClearAll}
          className="text-xs"
        >
          <RotateCcw className="w-3 h-3 mr-1" />
          Reset
        </Button>
      </div>
      <ScrollArea className="h-[calc(100vh-12rem)]">
        {filterContent}
      </ScrollArea>
    </div>
  );
}

export default CategoryFilters;