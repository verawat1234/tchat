import React, { useMemo } from 'react';
import { Button } from '../ui/button';
import { Loader2 } from 'lucide-react';
import {
  useGetStreamCategoryDetailQuery,
  StreamSubtab,
} from '../../services/streamApi';

interface MoviesSubtabsProps {
  selectedSubtab: string | null;
  onSubtabChange: (subtabId: string | null) => void;
  className?: string;
}

export function MoviesSubtabs({
  selectedSubtab,
  onSubtabChange,
  className = ''
}: MoviesSubtabsProps) {
  // Get movie category details including subtabs
  const {
    data: categoryDetailData,
    isLoading: categoryDetailLoading,
    error: categoryDetailError
  } = useGetStreamCategoryDetailQuery('movies');

  // Extract and filter movie subtabs
  const movieSubtabs = useMemo(() => {
    const subtabs = categoryDetailData?.subtabs || [];
    return subtabs.filter(subtab =>
      subtab.categoryId === 'movies' && subtab.isActive
    ).sort((a, b) => a.displayOrder - b.displayOrder);
  }, [categoryDetailData?.subtabs]);

  // Handle subtab selection
  const handleSubtabClick = (subtabId: string | null) => {
    onSubtabChange(subtabId);
  };

  // Loading state
  if (categoryDetailLoading) {
    return (
      <div className={`flex items-center justify-center p-2 ${className}`}>
        <Loader2 className="w-4 h-4 animate-spin" />
        <span className="ml-2 text-xs text-muted-foreground">Loading movie types...</span>
      </div>
    );
  }

  // Error state
  if (categoryDetailError) {
    return (
      <div className={`text-center p-2 ${className}`}>
        <p className="text-xs text-destructive">Failed to load movie types</p>
      </div>
    );
  }

  // No subtabs available
  if (movieSubtabs.length === 0) {
    return null;
  }

  return (
    <div className={`w-full overflow-x-auto scrollbar-hide ${className}`}>
      <div className="flex gap-2 pb-2 min-w-max">
        {/* All Movies Option */}
        <Button
          variant={selectedSubtab === null ? 'default' : 'outline'}
          size="sm"
          onClick={() => handleSubtabClick(null)}
          className="h-8 px-3 flex-shrink-0 whitespace-nowrap text-xs"
          aria-label="Show all movies"
        >
          All Movies
        </Button>

        {/* Dynamic Subtab Options */}
        {movieSubtabs.map((subtab) => {
          // Get filter criteria information for display
          const filterInfo = getFilterDisplayInfo(subtab);

          return (
            <Button
              key={subtab.id}
              variant={selectedSubtab === subtab.id ? 'default' : 'outline'}
              size="sm"
              onClick={() => handleSubtabClick(subtab.id)}
              className="h-8 px-3 flex-shrink-0 whitespace-nowrap text-xs"
              aria-label={`Filter by ${subtab.name}${filterInfo ? ` (${filterInfo})` : ''}`}
              title={filterInfo ? `${subtab.name} - ${filterInfo}` : subtab.name}
            >
              <span>{subtab.name}</span>
              {filterInfo && (
                <span className="ml-1 text-[10px] opacity-75">
                  ({filterInfo})
                </span>
              )}
            </Button>
          );
        })}
      </div>
    </div>
  );
}

// Helper function to extract and format filter criteria for display
function getFilterDisplayInfo(subtab: StreamSubtab): string | null {
  if (!subtab.filterCriteria) return null;

  const criteria = subtab.filterCriteria;

  // Duration-based filtering (primary use case for movies)
  if (criteria.maxDuration !== undefined) {
    const maxMinutes = Math.floor(criteria.maxDuration / 60);
    return `≤${maxMinutes}min`;
  }

  if (criteria.minDuration !== undefined) {
    const minMinutes = Math.floor(criteria.minDuration / 60);
    return `≥${minMinutes}min`;
  }

  // Duration range
  if (criteria.minDuration !== undefined && criteria.maxDuration !== undefined) {
    const minMinutes = Math.floor(criteria.minDuration / 60);
    const maxMinutes = Math.floor(criteria.maxDuration / 60);
    return `${minMinutes}-${maxMinutes}min`;
  }

  // Other filter types
  if (criteria.genre) {
    return `${criteria.genre}`;
  }

  if (criteria.rating) {
    return `${criteria.rating}+`;
  }

  if (criteria.year) {
    return `${criteria.year}`;
  }

  // If there are criteria but we don't recognize the format
  if (Object.keys(criteria).length > 0) {
    return 'filtered';
  }

  return null;
}

export default MoviesSubtabs;