import React, { useMemo } from 'react';
import { BookOpen, Mic, Film, Video, Music, Palette, Loader2 } from 'lucide-react';
import { Button } from '../ui/button';
import {
  useGetStreamCategoriesQuery,
  useGetStreamCategoryDetailQuery,
  StreamCategory,
  StreamSubtab,
} from '../../services/streamApi';
import MoviesSubtabs from './MoviesSubtabs';

interface StreamTabsProps {
  selectedCategory: string;
  selectedSubtab: string | null;
  onCategoryChange: (categoryId: string) => void;
  onSubtabChange: (subtabId: string | null) => void;
  className?: string;
}

// Icon mapping for categories
const categoryIcons = {
  'book-open': BookOpen,
  'microphone': Mic,
  'film': Film,
  'video': Video,
  'music': Music,
  'palette': Palette
} as const;

export function StreamTabs({
  selectedCategory,
  selectedSubtab,
  onCategoryChange,
  onSubtabChange,
  className = ''
}: StreamTabsProps) {
  // RTK Query hooks for data fetching
  const {
    data: categoriesData,
    isLoading: categoriesLoading,
    error: categoriesError
  } = useGetStreamCategoriesQuery();

  const {
    data: categoryDetailData,
    isLoading: categoryDetailLoading
  } = useGetStreamCategoryDetailQuery(selectedCategory, {
    skip: !selectedCategory
  });

  // Extract data from API responses
  const categories = categoriesData?.categories || [];
  const subtabs = categoryDetailData?.subtabs || [];

  // Filter subtabs for current category
  const currentCategorySubtabs = useMemo(() => {
    return subtabs.filter(subtab => subtab.categoryId === selectedCategory && subtab.isActive);
  }, [subtabs, selectedCategory]);

  // Handle category change and reset subtab
  const handleCategoryChange = (categoryId: string) => {
    onCategoryChange(categoryId);
    onSubtabChange(null); // Reset subtab when changing category
  };

  // Loading state
  if (categoriesLoading) {
    return (
      <div className={`flex items-center justify-center p-4 ${className}`}>
        <Loader2 className="w-6 h-6 animate-spin" />
        <span className="ml-2 text-sm">Loading categories...</span>
      </div>
    );
  }

  // Error state
  if (categoriesError) {
    return (
      <div className={`flex items-center justify-center p-4 ${className}`}>
        <p className="text-destructive text-sm">Failed to load categories</p>
      </div>
    );
  }

  return (
    <div className={`sticky top-0 z-30 border-b border-border p-3 bg-card/95 backdrop-blur-sm ${className}`}>
      {/* Main Category Navigation */}
      <div className="w-full overflow-x-auto scrollbar-hide">
        <div className="flex gap-2 pb-2 min-w-max">
          {categories.map((category) => {
            const IconComponent = categoryIcons[category.iconName as keyof typeof categoryIcons] || BookOpen;
            return (
              <Button
                key={category.id}
                variant={selectedCategory === category.id ? 'default' : 'outline'}
                size="sm"
                onClick={() => handleCategoryChange(category.id)}
                className="h-9 px-4 flex-shrink-0 whitespace-nowrap touch-manipulation hover:scale-105 transition-transform"
                aria-label={`Switch to ${category.name} category`}
              >
                <IconComponent className="w-4 h-4 mr-2" aria-hidden="true" />
                <span className="text-sm">{category.name}</span>
              </Button>
            );
          })}
        </div>
      </div>

      {/* Movies Subtabs - Specialized component for movies category */}
      {selectedCategory === 'movies' && (
        <div className="mt-3">
          <MoviesSubtabs
            selectedSubtab={selectedSubtab}
            onSubtabChange={onSubtabChange}
          />
        </div>
      )}

      {/* Generic Subtabs Navigation - For other categories with subtabs */}
      {selectedCategory !== 'movies' && currentCategorySubtabs.length > 0 && (
        <div className="mt-3 w-full overflow-x-auto scrollbar-hide">
          <div className="flex gap-2 pb-2 min-w-max">
            {/* All Content Option */}
            <Button
              variant={selectedSubtab === null ? 'default' : 'outline'}
              size="sm"
              onClick={() => onSubtabChange(null)}
              className="h-8 px-3 flex-shrink-0 whitespace-nowrap text-xs"
              aria-label={`Show all ${categories.find(c => c.id === selectedCategory)?.name} content`}
            >
              All
            </Button>

            {/* Subtab Options */}
            {currentCategorySubtabs.map((subtab) => (
              <Button
                key={subtab.id}
                variant={selectedSubtab === subtab.id ? 'default' : 'outline'}
                size="sm"
                onClick={() => onSubtabChange(subtab.id)}
                className="h-8 px-3 flex-shrink-0 whitespace-nowrap text-xs"
                aria-label={`Filter by ${subtab.name}`}
              >
                {subtab.name}
              </Button>
            ))}
          </div>

          {/* Loading indicator for subtabs */}
          {categoryDetailLoading && (
            <div className="flex items-center justify-center py-2">
              <Loader2 className="w-4 h-4 animate-spin" />
              <span className="ml-2 text-xs text-muted-foreground">Loading subtabs...</span>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default StreamTabs;