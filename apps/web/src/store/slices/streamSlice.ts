import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { StreamCategory, StreamContentItem, StreamSubtab } from '../../services/streamApi';

// Stream content state management
interface StreamState {
  // Current navigation state
  currentCategory: string;
  currentSubtab: string | null;

  // Content cache
  contentCache: {
    [key: string]: StreamContentItem[];
  };
  featuredCache: {
    [key: string]: StreamContentItem[];
  };
  categoriesCache: StreamCategory[];
  subtabsCache: {
    [categoryId: string]: StreamSubtab[];
  };

  // User preferences
  preferences: {
    autoplayFeatured: boolean;
    defaultCategory: string;
    viewMode: 'grid' | 'list';
    sortOrder: 'featured' | 'newest' | 'price_asc' | 'price_desc' | 'rating';
    showSubtabs: boolean;
  };

  // UI state
  ui: {
    isCarouselPlaying: boolean;
    selectedContentId: string | null;
    showPurchaseDialog: boolean;
    searchQuery: string;
    filters: {
      priceRange: [number, number] | null;
      categories: string[];
      availability: ('available' | 'coming_soon' | 'unavailable')[];
      contentTypes: string[];
    };
  };

  // Performance tracking
  performance: {
    lastLoadTime: number | null;
    cacheHitRate: number;
    totalApiCalls: number;
    errorCount: number;
  };
}

// Initial state
const initialState: StreamState = {
  currentCategory: 'books',
  currentSubtab: null,

  contentCache: {},
  featuredCache: {},
  categoriesCache: [],
  subtabsCache: {},

  preferences: {
    autoplayFeatured: true,
    defaultCategory: 'books',
    viewMode: 'grid',
    sortOrder: 'featured',
    showSubtabs: true,
  },

  ui: {
    isCarouselPlaying: false,
    selectedContentId: null,
    showPurchaseDialog: false,
    searchQuery: '',
    filters: {
      priceRange: null,
      categories: [],
      availability: ['available'],
      contentTypes: [],
    },
  },

  performance: {
    lastLoadTime: null,
    cacheHitRate: 0,
    totalApiCalls: 0,
    errorCount: 0,
  },
};

// Stream slice
const streamSlice = createSlice({
  name: 'stream',
  initialState,
  reducers: {
    // Navigation actions
    setCurrentCategory: (state, action: PayloadAction<string>) => {
      state.currentCategory = action.payload;
      state.currentSubtab = null; // Reset subtab when changing category
    },

    setCurrentSubtab: (state, action: PayloadAction<string | null>) => {
      state.currentSubtab = action.payload;
    },

    // Content cache actions
    cacheContent: (state, action: PayloadAction<{ key: string; content: StreamContentItem[] }>) => {
      const { key, content } = action.payload;
      state.contentCache[key] = content;
      state.performance.cacheHitRate = Object.keys(state.contentCache).length / (state.performance.totalApiCalls || 1);
    },

    cacheFeaturedContent: (state, action: PayloadAction<{ categoryId: string; content: StreamContentItem[] }>) => {
      const { categoryId, content } = action.payload;
      state.featuredCache[categoryId] = content;
    },

    cacheCategories: (state, action: PayloadAction<StreamCategory[]>) => {
      state.categoriesCache = action.payload;
    },

    cacheSubtabs: (state, action: PayloadAction<{ categoryId: string; subtabs: StreamSubtab[] }>) => {
      const { categoryId, subtabs } = action.payload;
      state.subtabsCache[categoryId] = subtabs;
    },

    // User preferences actions
    updatePreferences: (state, action: PayloadAction<Partial<StreamState['preferences']>>) => {
      state.preferences = { ...state.preferences, ...action.payload };
    },

    setDefaultCategory: (state, action: PayloadAction<string>) => {
      state.preferences.defaultCategory = action.payload;
    },

    setViewMode: (state, action: PayloadAction<'grid' | 'list'>) => {
      state.preferences.viewMode = action.payload;
    },

    setSortOrder: (state, action: PayloadAction<StreamState['preferences']['sortOrder']>) => {
      state.preferences.sortOrder = action.payload;
    },

    toggleAutoplay: (state) => {
      state.preferences.autoplayFeatured = !state.preferences.autoplayFeatured;
    },

    toggleSubtabs: (state) => {
      state.preferences.showSubtabs = !state.preferences.showSubtabs;
    },

    // UI state actions
    setCarouselPlaying: (state, action: PayloadAction<boolean>) => {
      state.ui.isCarouselPlaying = action.payload;
    },

    setSelectedContent: (state, action: PayloadAction<string | null>) => {
      state.ui.selectedContentId = action.payload;
    },

    togglePurchaseDialog: (state, action: PayloadAction<boolean>) => {
      state.ui.showPurchaseDialog = action.payload;
    },

    setSearchQuery: (state, action: PayloadAction<string>) => {
      state.ui.searchQuery = action.payload;
    },

    updateFilters: (state, action: PayloadAction<Partial<StreamState['ui']['filters']>>) => {
      state.ui.filters = { ...state.ui.filters, ...action.payload };
    },

    clearFilters: (state) => {
      state.ui.filters = {
        priceRange: null,
        categories: [],
        availability: ['available'],
        contentTypes: [],
      };
    },

    // Performance tracking actions
    incrementApiCalls: (state) => {
      state.performance.totalApiCalls += 1;
    },

    incrementErrors: (state) => {
      state.performance.errorCount += 1;
    },

    setLastLoadTime: (state, action: PayloadAction<number>) => {
      state.performance.lastLoadTime = action.payload;
    },

    updateCacheHitRate: (state) => {
      state.performance.cacheHitRate = Object.keys(state.contentCache).length / (state.performance.totalApiCalls || 1);
    },

    // Utility actions
    clearCache: (state) => {
      state.contentCache = {};
      state.featuredCache = {};
      state.subtabsCache = {};
    },

    resetToDefaults: (state) => {
      state.currentCategory = state.preferences.defaultCategory;
      state.currentSubtab = null;
      state.ui.searchQuery = '';
      state.ui.selectedContentId = null;
      state.ui.filters = initialState.ui.filters;
    },

    // Session management
    initializeFromSession: (state, action: PayloadAction<Partial<StreamState>>) => {
      const sessionData = action.payload;
      if (sessionData.currentCategory) state.currentCategory = sessionData.currentCategory;
      if (sessionData.currentSubtab !== undefined) state.currentSubtab = sessionData.currentSubtab;
      if (sessionData.preferences) {
        state.preferences = { ...state.preferences, ...sessionData.preferences };
      }
      if (sessionData.ui) {
        state.ui = { ...state.ui, ...sessionData.ui };
      }
    },
  },
});

// Export actions
export const {
  // Navigation
  setCurrentCategory,
  setCurrentSubtab,

  // Cache management
  cacheContent,
  cacheFeaturedContent,
  cacheCategories,
  cacheSubtabs,
  clearCache,

  // User preferences
  updatePreferences,
  setDefaultCategory,
  setViewMode,
  setSortOrder,
  toggleAutoplay,
  toggleSubtabs,

  // UI state
  setCarouselPlaying,
  setSelectedContent,
  togglePurchaseDialog,
  setSearchQuery,
  updateFilters,
  clearFilters,

  // Performance
  incrementApiCalls,
  incrementErrors,
  setLastLoadTime,
  updateCacheHitRate,

  // Utility
  resetToDefaults,
  initializeFromSession,
} = streamSlice.actions;

// Selectors
export const selectCurrentCategory = (state: { stream: StreamState }) => state.stream.currentCategory;
export const selectCurrentSubtab = (state: { stream: StreamState }) => state.stream.currentSubtab;
export const selectStreamPreferences = (state: { stream: StreamState }) => state.stream.preferences;
export const selectStreamUI = (state: { stream: StreamState }) => state.stream.ui;
export const selectStreamPerformance = (state: { stream: StreamState }) => state.stream.performance;

export const selectCachedContent = (state: { stream: StreamState }, key: string) =>
  state.stream.contentCache[key] || [];

export const selectCachedFeatured = (state: { stream: StreamState }, categoryId: string) =>
  state.stream.featuredCache[categoryId] || [];

export const selectCachedCategories = (state: { stream: StreamState }) =>
  state.stream.categoriesCache;

export const selectCachedSubtabs = (state: { stream: StreamState }, categoryId: string) =>
  state.stream.subtabsCache[categoryId] || [];

export const selectIsContentCached = (state: { stream: StreamState }, key: string) =>
  key in state.stream.contentCache;

export const selectCurrentCategorySubtabs = (state: { stream: StreamState }) =>
  state.stream.subtabsCache[state.stream.currentCategory] || [];

// Complex selectors
export const selectFilteredContent = (state: { stream: StreamState }, content: StreamContentItem[]) => {
  const { searchQuery, filters } = state.stream.ui;
  const { sortOrder } = state.stream.preferences;

  let filtered = content;

  // Apply search filter
  if (searchQuery) {
    filtered = filtered.filter(item =>
      item.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.description.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }

  // Apply availability filter
  if (filters.availability.length > 0) {
    filtered = filtered.filter(item =>
      filters.availability.includes(item.availabilityStatus)
    );
  }

  // Apply price range filter
  if (filters.priceRange) {
    const [min, max] = filters.priceRange;
    filtered = filtered.filter(item => item.price >= min && item.price <= max);
  }

  // Apply content type filter
  if (filters.contentTypes.length > 0) {
    filtered = filtered.filter(item =>
      filters.contentTypes.includes(item.contentType)
    );
  }

  // Apply sorting
  switch (sortOrder) {
    case 'featured':
      filtered.sort((a, b) => {
        if (a.isFeatured && !b.isFeatured) return -1;
        if (!a.isFeatured && b.isFeatured) return 1;
        return (a.featuredOrder || 999) - (b.featuredOrder || 999);
      });
      break;
    case 'newest':
      filtered.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
      break;
    case 'price_asc':
      filtered.sort((a, b) => a.price - b.price);
      break;
    case 'price_desc':
      filtered.sort((a, b) => b.price - a.price);
      break;
    case 'rating':
      // Sort by metadata rating if available
      filtered.sort((a, b) => {
        const ratingA = a.metadata?.rating || 0;
        const ratingB = b.metadata?.rating || 0;
        return ratingB - ratingA;
      });
      break;
  }

  return filtered;
};

export default streamSlice.reducer;