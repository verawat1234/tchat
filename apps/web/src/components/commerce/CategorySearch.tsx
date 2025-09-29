/**
 * CategorySearch - Search component for categories and products
 *
 * Provides advanced search functionality with autocomplete, filters, and suggestions.
 * Supports searching across categories, products, and global search.
 *
 * Features:
 * - Real-time search with debouncing
 * - Autocomplete with suggestions
 * - Recent searches history
 * - Search filters and scoping
 * - Keyboard navigation support
 * - Voice search support
 * - Search analytics tracking
 * - Mobile-optimized interface
 */

import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Badge } from '../ui/badge';
import { ScrollArea } from '../ui/scroll-area';
import { Skeleton } from '../ui/skeleton';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from '../ui/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '../ui/popover';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../ui/dialog';
import {
  Search,
  X,
  Clock,
  TrendingUp,
  Filter,
  Mic,
  ArrowRight,
  Package,
  Folder,
  History,
  Star,
} from 'lucide-react';
import { cn } from '../../lib/utils';
import { useCategory } from './CategoryProvider';
import { useGetProductsQuery } from '../../services/commerceApi';
import type { Category, Product } from '../../types/commerce';

// ===== Types =====

export interface CategorySearchProps {
  className?: string;
  placeholder?: string;
  variant?: 'default' | 'compact' | 'expanded';
  showFilters?: boolean;
  showVoiceSearch?: boolean;
  showRecentSearches?: boolean;
  autoFocus?: boolean;
  onSearchSubmit?: (query: string) => void;
  onCategorySelect?: (category: Category) => void;
  onProductSelect?: (product: Product) => void;
}

interface SearchSuggestion {
  type: 'category' | 'product' | 'recent' | 'trending';
  id: string;
  title: string;
  subtitle?: string;
  icon?: React.ReactNode;
  badge?: string;
}

// ===== Search Filters Component =====

interface SearchFiltersProps {
  onFiltersChange: (filters: SearchFilters) => void;
}

interface SearchFilters {
  scope: 'all' | 'categories' | 'products';
  priceRange: [number, number];
  rating: number;
  inStock: boolean;
  featured: boolean;
}

function SearchFilters({ onFiltersChange }: SearchFiltersProps) {
  const [filters, setFilters] = useState<SearchFilters>({
    scope: 'all',
    priceRange: [0, 10000],
    rating: 0,
    inStock: false,
    featured: false,
  });

  const handleFilterChange = <K extends keyof SearchFilters>(
    key: K,
    value: SearchFilters[K]
  ) => {
    const newFilters = { ...filters, [key]: value };
    setFilters(newFilters);
    onFiltersChange(newFilters);
  };

  const activeFilterCount = Object.entries(filters).filter(([key, value]) => {
    if (key === 'scope') return value !== 'all';
    if (key === 'priceRange') return value[0] > 0 || value[1] < 10000;
    if (key === 'rating') return value > 0;
    return value === true;
  }).length;

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2">
          <Filter className="w-4 h-4" />
          Filters
          {activeFilterCount > 0 && (
            <Badge variant="secondary" className="h-5 px-1 text-xs">
              {activeFilterCount}
            </Badge>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-80 p-4" align="start">
        <div className="space-y-4">
          <div>
            <label className="text-sm font-medium mb-2 block">Search in</label>
            <div className="flex gap-2">
              <Button
                variant={filters.scope === 'all' ? 'default' : 'outline'}
                size="sm"
                onClick={() => handleFilterChange('scope', 'all')}
              >
                All
              </Button>
              <Button
                variant={filters.scope === 'categories' ? 'default' : 'outline'}
                size="sm"
                onClick={() => handleFilterChange('scope', 'categories')}
              >
                Categories
              </Button>
              <Button
                variant={filters.scope === 'products' ? 'default' : 'outline'}
                size="sm"
                onClick={() => handleFilterChange('scope', 'products')}
              >
                Products
              </Button>
            </div>
          </div>

          {filters.scope !== 'categories' && (
            <>
              <div>
                <label className="text-sm font-medium mb-2 block">
                  Price Range: ฿{filters.priceRange[0]} - ฿{filters.priceRange[1]}
                </label>
                <div className="space-y-2">
                  <input
                    type="range"
                    min="0"
                    max="10000"
                    step="100"
                    value={filters.priceRange[0]}
                    onChange={(e) => handleFilterChange('priceRange', [
                      Number(e.target.value),
                      filters.priceRange[1]
                    ])}
                    className="w-full"
                  />
                  <input
                    type="range"
                    min="0"
                    max="10000"
                    step="100"
                    value={filters.priceRange[1]}
                    onChange={(e) => handleFilterChange('priceRange', [
                      filters.priceRange[0],
                      Number(e.target.value)
                    ])}
                    className="w-full"
                  />
                </div>
              </div>

              <div>
                <label className="text-sm font-medium mb-2 block">
                  Minimum Rating: {filters.rating} stars
                </label>
                <input
                  type="range"
                  min="0"
                  max="5"
                  step="0.5"
                  value={filters.rating}
                  onChange={(e) => handleFilterChange('rating', Number(e.target.value))}
                  className="w-full"
                />
              </div>

              <div className="space-y-2">
                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={filters.inStock}
                    onChange={(e) => handleFilterChange('inStock', e.target.checked)}
                  />
                  <span className="text-sm">In stock only</span>
                </label>

                <label className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={filters.featured}
                    onChange={(e) => handleFilterChange('featured', e.target.checked)}
                  />
                  <span className="text-sm">Featured items</span>
                </label>
              </div>
            </>
          )}

          {activeFilterCount > 0 && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                const resetFilters: SearchFilters = {
                  scope: 'all',
                  priceRange: [0, 10000],
                  rating: 0,
                  inStock: false,
                  featured: false,
                };
                setFilters(resetFilters);
                onFiltersChange(resetFilters);
              }}
              className="w-full"
            >
              Clear all filters
            </Button>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}

// ===== Voice Search Component =====

interface VoiceSearchProps {
  onResult: (text: string) => void;
}

function VoiceSearch({ onResult }: VoiceSearchProps) {
  const [isListening, setIsListening] = useState(false);
  const [isSupported, setIsSupported] = useState(false);
  const recognitionRef = useRef<any>(null);

  useEffect(() => {
    // Check if speech recognition is supported
    const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
    setIsSupported(!!SpeechRecognition);

    if (SpeechRecognition) {
      recognitionRef.current = new SpeechRecognition();
      recognitionRef.current.continuous = false;
      recognitionRef.current.interimResults = false;
      recognitionRef.current.lang = 'en-US';

      recognitionRef.current.onresult = (event: any) => {
        const transcript = event.results[0][0].transcript;
        onResult(transcript);
        setIsListening(false);
      };

      recognitionRef.current.onerror = () => {
        setIsListening(false);
      };

      recognitionRef.current.onend = () => {
        setIsListening(false);
      };
    }

    return () => {
      if (recognitionRef.current) {
        recognitionRef.current.stop();
      }
    };
  }, [onResult]);

  const handleVoiceSearch = () => {
    if (!isSupported || !recognitionRef.current) return;

    if (isListening) {
      recognitionRef.current.stop();
    } else {
      recognitionRef.current.start();
      setIsListening(true);
    }
  };

  if (!isSupported) return null;

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={handleVoiceSearch}
      className={cn(
        'absolute right-2 top-1/2 -translate-y-1/2',
        isListening && 'text-red-500 animate-pulse'
      )}
      title={isListening ? 'Stop listening' : 'Voice search'}
    >
      <Mic className="w-4 h-4" />
    </Button>
  );
}

// ===== Search Suggestions Component =====

interface SearchSuggestionsProps {
  query: string;
  suggestions: SearchSuggestion[];
  recentSearches: string[];
  onSuggestionSelect: (suggestion: SearchSuggestion | string) => void;
  onClearRecentSearches: () => void;
}

function SearchSuggestions({
  query,
  suggestions,
  recentSearches,
  onSuggestionSelect,
  onClearRecentSearches,
}: SearchSuggestionsProps) {
  if (!query && recentSearches.length === 0 && suggestions.length === 0) {
    return (
      <div className="p-4 text-center text-muted-foreground">
        Start typing to see suggestions
      </div>
    );
  }

  return (
    <Command>
      <CommandList>
        {!query && recentSearches.length > 0 && (
          <CommandGroup heading="Recent Searches">
            {recentSearches.slice(0, 5).map((search, index) => (
              <CommandItem
                key={index}
                onSelect={() => onSuggestionSelect(search)}
                className="flex items-center gap-2"
              >
                <Clock className="w-4 h-4 text-muted-foreground" />
                <span>{search}</span>
              </CommandItem>
            ))}
            <CommandSeparator />
            <CommandItem onSelect={onClearRecentSearches} className="text-muted-foreground">
              <History className="w-4 h-4 mr-2" />
              Clear recent searches
            </CommandItem>
          </CommandGroup>
        )}

        {query && suggestions.length === 0 && (
          <CommandEmpty>No results found for "{query}"</CommandEmpty>
        )}

        {suggestions.length > 0 && (
          <>
            {/* Category suggestions */}
            {suggestions.filter(s => s.type === 'category').length > 0 && (
              <CommandGroup heading="Categories">
                {suggestions
                  .filter(s => s.type === 'category')
                  .slice(0, 3)
                  .map((suggestion) => (
                    <CommandItem
                      key={suggestion.id}
                      onSelect={() => onSuggestionSelect(suggestion)}
                      className="flex items-center gap-2"
                    >
                      <Folder className="w-4 h-4 text-blue-500" />
                      <div className="flex-1">
                        <div className="font-medium">{suggestion.title}</div>
                        {suggestion.subtitle && (
                          <div className="text-xs text-muted-foreground">
                            {suggestion.subtitle}
                          </div>
                        )}
                      </div>
                      {suggestion.badge && (
                        <Badge variant="secondary" className="text-xs">
                          {suggestion.badge}
                        </Badge>
                      )}
                    </CommandItem>
                  ))}
              </CommandGroup>
            )}

            {/* Product suggestions */}
            {suggestions.filter(s => s.type === 'product').length > 0 && (
              <CommandGroup heading="Products">
                {suggestions
                  .filter(s => s.type === 'product')
                  .slice(0, 5)
                  .map((suggestion) => (
                    <CommandItem
                      key={suggestion.id}
                      onSelect={() => onSuggestionSelect(suggestion)}
                      className="flex items-center gap-2"
                    >
                      <Package className="w-4 h-4 text-green-500" />
                      <div className="flex-1">
                        <div className="font-medium">{suggestion.title}</div>
                        {suggestion.subtitle && (
                          <div className="text-xs text-muted-foreground">
                            {suggestion.subtitle}
                          </div>
                        )}
                      </div>
                      {suggestion.badge && (
                        <Badge variant="secondary" className="text-xs">
                          {suggestion.badge}
                        </Badge>
                      )}
                    </CommandItem>
                  ))}
              </CommandGroup>
            )}

            {/* Trending suggestions */}
            {suggestions.filter(s => s.type === 'trending').length > 0 && (
              <CommandGroup heading="Trending">
                {suggestions
                  .filter(s => s.type === 'trending')
                  .slice(0, 3)
                  .map((suggestion) => (
                    <CommandItem
                      key={suggestion.id}
                      onSelect={() => onSuggestionSelect(suggestion)}
                      className="flex items-center gap-2"
                    >
                      <TrendingUp className="w-4 h-4 text-orange-500" />
                      <span>{suggestion.title}</span>
                    </CommandItem>
                  ))}
              </CommandGroup>
            )}
          </>
        )}
      </CommandList>
    </Command>
  );
}

// ===== Main CategorySearch Component =====

export function CategorySearch({
  className,
  placeholder = 'Search categories and products...',
  variant = 'default',
  showFilters = true,
  showVoiceSearch = true,
  showRecentSearches = true,
  autoFocus = false,
  onSearchSubmit,
  onCategorySelect,
  onProductSelect,
}: CategorySearchProps) {
  const {
    state,
    setSearchQuery,
    addRecentSearch,
    clearRecentSearches,
    categoriesQuery,
    searchResults,
  } = useCategory();

  const [isOpen, setIsOpen] = useState(false);
  const [searchFilters, setSearchFilters] = useState<SearchFilters>({
    scope: 'all',
    priceRange: [0, 10000],
    rating: 0,
    inStock: false,
    featured: false,
  });

  const inputRef = useRef<HTMLInputElement>(null);

  // Debounced search query
  const [debouncedQuery, setDebouncedQuery] = useState(state.search.query);

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(state.search.query);
    }, 300);

    return () => clearTimeout(timer);
  }, [state.search.query]);

  // Search for products when query changes
  const { data: productsData } = useGetProductsQuery({
    filters: { search: debouncedQuery },
    pagination: { page: 1, pageSize: 10 },
  }, {
    skip: !debouncedQuery || searchFilters.scope === 'categories',
  });

  // Build suggestions
  const suggestions = useMemo((): SearchSuggestion[] => {
    const results: SearchSuggestion[] = [];

    // Add category suggestions
    if (searchFilters.scope !== 'products') {
      searchResults.slice(0, 3).forEach(category => {
        results.push({
          type: 'category',
          id: category.id,
          title: category.name,
          subtitle: category.shortDescription,
          badge: `${category.activeProductCount} products`,
        });
      });
    }

    // Add product suggestions
    if (searchFilters.scope !== 'categories' && productsData?.products) {
      productsData.products.slice(0, 5).forEach(product => {
        results.push({
          type: 'product',
          id: product.id,
          title: product.name,
          subtitle: product.category,
          badge: `฿${parseFloat(product.price).toFixed(0)}`,
        });
      });
    }

    // Add trending suggestions (mock data)
    if (!debouncedQuery) {
      ['smartphones', 'laptops', 'headphones'].forEach((term, index) => {
        results.push({
          type: 'trending',
          id: `trending-${index}`,
          title: term,
        });
      });
    }

    return results;
  }, [searchResults, productsData, searchFilters.scope, debouncedQuery]);

  const handleInputChange = (value: string) => {
    setSearchQuery(value);
    setIsOpen(value.length > 0 || showRecentSearches);
  };

  const handleSuggestionSelect = (suggestion: SearchSuggestion | string) => {
    if (typeof suggestion === 'string') {
      setSearchQuery(suggestion);
      addRecentSearch(suggestion);
      onSearchSubmit?.(suggestion);
    } else {
      const query = suggestion.title;
      setSearchQuery(query);
      addRecentSearch(query);

      if (suggestion.type === 'category') {
        const category = categoriesQuery.data?.categories.find(c => c.id === suggestion.id);
        if (category) {
          onCategorySelect?.(category);
        }
      } else if (suggestion.type === 'product') {
        const product = productsData?.products.find(p => p.id === suggestion.id);
        if (product) {
          onProductSelect?.(product);
        }
      }

      onSearchSubmit?.(query);
    }

    setIsOpen(false);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (state.search.query.trim()) {
      addRecentSearch(state.search.query);
      onSearchSubmit?.(state.search.query);
      setIsOpen(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      setIsOpen(false);
      inputRef.current?.blur();
    }
  };

  const handleVoiceResult = (text: string) => {
    setSearchQuery(text);
    addRecentSearch(text);
    onSearchSubmit?.(text);
  };

  // Compact variant
  if (variant === 'compact') {
    return (
      <div className={cn('relative', className)}>
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            ref={inputRef}
            value={state.search.query}
            onChange={(e) => handleInputChange(e.target.value)}
            onFocus={() => setIsOpen(true)}
            onKeyDown={handleKeyDown}
            placeholder={placeholder}
            className="pl-9 pr-10"
            autoFocus={autoFocus}
          />
          {state.search.query && (
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setSearchQuery('')}
              className="absolute right-1 top-1/2 -translate-y-1/2 w-6 h-6"
            >
              <X className="w-3 h-3" />
            </Button>
          )}
        </div>

        {isOpen && (
          <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-popover border rounded-md shadow-lg">
            <SearchSuggestions
              query={debouncedQuery}
              suggestions={suggestions}
              recentSearches={showRecentSearches ? state.search.recentSearches : []}
              onSuggestionSelect={handleSuggestionSelect}
              onClearRecentSearches={clearRecentSearches}
            />
          </div>
        )}
      </div>
    );
  }

  // Expanded variant (modal)
  if (variant === 'expanded') {
    return (
      <Dialog>
        <DialogTrigger asChild>
          <Button variant="outline" className={cn('justify-start text-muted-foreground', className)}>
            <Search className="w-4 h-4 mr-2" />
            {placeholder}
          </Button>
        </DialogTrigger>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Search</DialogTitle>
          </DialogHeader>
          <CategorySearch
            variant="default"
            autoFocus
            showFilters={showFilters}
            showVoiceSearch={showVoiceSearch}
            showRecentSearches={showRecentSearches}
            onSearchSubmit={onSearchSubmit}
            onCategorySelect={onCategorySelect}
            onProductSelect={onProductSelect}
          />
        </DialogContent>
      </Dialog>
    );
  }

  // Default variant
  return (
    <div className={cn('relative', className)}>
      <form onSubmit={handleSubmit} className="flex gap-2">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            ref={inputRef}
            value={state.search.query}
            onChange={(e) => handleInputChange(e.target.value)}
            onFocus={() => setIsOpen(true)}
            onKeyDown={handleKeyDown}
            placeholder={placeholder}
            className="pl-9 pr-12"
            autoFocus={autoFocus}
          />

          {state.search.query && (
            <Button
              type="button"
              variant="ghost"
              size="icon"
              onClick={() => setSearchQuery('')}
              className="absolute right-8 top-1/2 -translate-y-1/2 w-6 h-6"
            >
              <X className="w-3 h-3" />
            </Button>
          )}

          {showVoiceSearch && (
            <VoiceSearch onResult={handleVoiceResult} />
          )}
        </div>

        {showFilters && (
          <SearchFilters onFiltersChange={setSearchFilters} />
        )}

        <Button type="submit" disabled={!state.search.query.trim()}>
          Search
        </Button>
      </form>

      {isOpen && (
        <div className="absolute top-full left-0 right-0 z-50 mt-2 bg-popover border rounded-md shadow-lg max-h-96 overflow-hidden">
          <SearchSuggestions
            query={debouncedQuery}
            suggestions={suggestions}
            recentSearches={showRecentSearches ? state.search.recentSearches : []}
            onSuggestionSelect={handleSuggestionSelect}
            onClearRecentSearches={clearRecentSearches}
          />
        </div>
      )}
    </div>
  );
}

export default CategorySearch;