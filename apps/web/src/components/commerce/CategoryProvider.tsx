/**
 * CategoryProvider - Context provider for category browsing state management
 *
 * Manages hierarchical category navigation, search state, filters, and analytics tracking.
 * Provides centralized state management for all category-related components.
 *
 * Features:
 * - Hierarchical category navigation with breadcrumbs
 * - Search and filter state management
 * - Analytics tracking integration
 * - Optimistic updates and caching
 * - Cross-component state synchronization
 */

import React, { createContext, useContext, useReducer, useEffect, useMemo } from 'react';
import { useGetCategoriesQuery, useGetRootCategoriesQuery, useTrackCategoryViewMutation } from '../../services/commerceApi';
import type { Category, CategoryResponse, ProductFilters, SortOptions } from '../../types/commerce';

// ===== State Types =====

export interface CategorySearchState {
  query: string;
  suggestions: Category[];
  isSearching: boolean;
  recentSearches: string[];
}

export interface CategoryFilterState {
  priceRange: {
    min: number;
    max: number;
  };
  rating: number;
  inStock: boolean;
  onSale: boolean;
  featured: boolean;
  brands: string[];
  attributes: Record<string, string[]>;
}

export interface CategoryBreadcrumb {
  id: string;
  name: string;
  slug: string;
}

export interface CategoryViewState {
  viewMode: 'grid' | 'list';
  itemsPerPage: number;
  sortBy: SortOptions;
  showFilters: boolean;
  selectedCategory: Category | null;
  parentCategories: Category[];
  childCategories: Category[];
  breadcrumbs: CategoryBreadcrumb[];
}

export interface CategoryState {
  // Navigation state
  currentCategoryId: string | null;
  viewState: CategoryViewState;

  // Search state
  search: CategorySearchState;

  // Filter state
  filters: CategoryFilterState;
  activeFilters: string[];

  // UI state
  isLoading: boolean;
  error: string | null;

  // Analytics
  viewHistory: string[];
  sessionStartTime: number;
}

// ===== Action Types =====

type CategoryAction =
  | { type: 'SET_CURRENT_CATEGORY'; payload: string | null }
  | { type: 'SET_VIEW_MODE'; payload: 'grid' | 'list' }
  | { type: 'SET_SORT_BY'; payload: SortOptions }
  | { type: 'SET_ITEMS_PER_PAGE'; payload: number }
  | { type: 'TOGGLE_FILTERS' }
  | { type: 'SET_SEARCH_QUERY'; payload: string }
  | { type: 'SET_SEARCH_SUGGESTIONS'; payload: Category[] }
  | { type: 'SET_SEARCHING'; payload: boolean }
  | { type: 'ADD_RECENT_SEARCH'; payload: string }
  | { type: 'CLEAR_RECENT_SEARCHES' }
  | { type: 'UPDATE_FILTERS'; payload: Partial<CategoryFilterState> }
  | { type: 'CLEAR_FILTERS' }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'TRACK_CATEGORY_VIEW'; payload: string }
  | { type: 'UPDATE_BREADCRUMBS'; payload: CategoryBreadcrumb[] }
  | { type: 'SET_PARENT_CATEGORIES'; payload: Category[] }
  | { type: 'SET_CHILD_CATEGORIES'; payload: Category[] };

// ===== Initial State =====

const initialState: CategoryState = {
  currentCategoryId: null,
  viewState: {
    viewMode: 'grid',
    itemsPerPage: 20,
    sortBy: { field: 'name', order: 'asc' },
    showFilters: false,
    selectedCategory: null,
    parentCategories: [],
    childCategories: [],
    breadcrumbs: [],
  },
  search: {
    query: '',
    suggestions: [],
    isSearching: false,
    recentSearches: JSON.parse(localStorage.getItem('tchat-category-recent-searches') || '[]'),
  },
  filters: {
    priceRange: { min: 0, max: 10000 },
    rating: 0,
    inStock: false,
    onSale: false,
    featured: false,
    brands: [],
    attributes: {},
  },
  activeFilters: [],
  isLoading: false,
  error: null,
  viewHistory: [],
  sessionStartTime: Date.now(),
};

// ===== Reducer =====

function categoryReducer(state: CategoryState, action: CategoryAction): CategoryState {
  switch (action.type) {
    case 'SET_CURRENT_CATEGORY':
      return {
        ...state,
        currentCategoryId: action.payload,
      };

    case 'SET_VIEW_MODE':
      return {
        ...state,
        viewState: {
          ...state.viewState,
          viewMode: action.payload,
        },
      };

    case 'SET_SORT_BY':
      return {
        ...state,
        viewState: {
          ...state.viewState,
          sortBy: action.payload,
        },
      };

    case 'SET_ITEMS_PER_PAGE':
      return {
        ...state,
        viewState: {
          ...state.viewState,
          itemsPerPage: action.payload,
        },
      };

    case 'TOGGLE_FILTERS':
      return {
        ...state,
        viewState: {
          ...state.viewState,
          showFilters: !state.viewState.showFilters,
        },
      };

    case 'SET_SEARCH_QUERY':
      return {
        ...state,
        search: {
          ...state.search,
          query: action.payload,
        },
      };

    case 'SET_SEARCH_SUGGESTIONS':
      return {
        ...state,
        search: {
          ...state.search,
          suggestions: action.payload,
        },
      };

    case 'SET_SEARCHING':
      return {
        ...state,
        search: {
          ...state.search,
          isSearching: action.payload,
        },
      };

    case 'ADD_RECENT_SEARCH':
      const newRecentSearches = [
        action.payload,
        ...state.search.recentSearches.filter(s => s !== action.payload),
      ].slice(0, 10);

      // Save to localStorage
      localStorage.setItem('tchat-category-recent-searches', JSON.stringify(newRecentSearches));

      return {
        ...state,
        search: {
          ...state.search,
          recentSearches: newRecentSearches,
        },
      };

    case 'CLEAR_RECENT_SEARCHES':
      localStorage.removeItem('tchat-category-recent-searches');
      return {
        ...state,
        search: {
          ...state.search,
          recentSearches: [],
        },
      };

    case 'UPDATE_FILTERS':
      const updatedFilters = { ...state.filters, ...action.payload };
      const activeFilters = Object.entries(updatedFilters)
        .filter(([key, value]) => {
          if (key === 'priceRange') {
            const range = value as { min: number; max: number };
            return range.min > 0 || range.max < 10000;
          }
          if (key === 'rating') return (value as number) > 0;
          if (typeof value === 'boolean') return value;
          if (Array.isArray(value)) return value.length > 0;
          if (typeof value === 'object') return Object.keys(value).length > 0;
          return false;
        })
        .map(([key]) => key);

      return {
        ...state,
        filters: updatedFilters,
        activeFilters,
      };

    case 'CLEAR_FILTERS':
      return {
        ...state,
        filters: initialState.filters,
        activeFilters: [],
      };

    case 'SET_LOADING':
      return {
        ...state,
        isLoading: action.payload,
      };

    case 'SET_ERROR':
      return {
        ...state,
        error: action.payload,
        isLoading: false,
      };

    case 'TRACK_CATEGORY_VIEW':
      return {
        ...state,
        viewHistory: [action.payload, ...state.viewHistory.slice(0, 19)],
      };

    case 'UPDATE_BREADCRUMBS':
      return {
        ...state,
        viewState: {
          ...state.viewState,
          breadcrumbs: action.payload,
        },
      };

    case 'SET_PARENT_CATEGORIES':
      return {
        ...state,
        viewState: {
          ...state.viewState,
          parentCategories: action.payload,
        },
      };

    case 'SET_CHILD_CATEGORIES':
      return {
        ...state,
        viewState: {
          ...state.viewState,
          childCategories: action.payload,
        },
      };

    default:
      return state;
  }
}

// ===== Context =====

interface CategoryContextValue {
  // State
  state: CategoryState;

  // Navigation actions
  setCurrentCategory: (id: string | null) => void;
  setViewMode: (mode: 'grid' | 'list') => void;
  setSortBy: (sort: SortOptions) => void;
  setItemsPerPage: (count: number) => void;
  toggleFilters: () => void;

  // Search actions
  setSearchQuery: (query: string) => void;
  addRecentSearch: (query: string) => void;
  clearRecentSearches: () => void;

  // Filter actions
  updateFilters: (filters: Partial<CategoryFilterState>) => void;
  clearFilters: () => void;

  // Analytics actions
  trackCategoryView: (categoryId: string) => void;

  // Data queries
  categoriesQuery: ReturnType<typeof useGetCategoriesQuery>;
  rootCategoriesQuery: ReturnType<typeof useGetRootCategoriesQuery>;

  // Computed values
  filteredProducts: ProductFilters;
  hasActiveFilters: boolean;
  searchResults: Category[];
}

const CategoryContext = createContext<CategoryContextValue | null>(null);

// ===== Provider Component =====

interface CategoryProviderProps {
  children: React.ReactNode;
  businessId?: string;
}

export function CategoryProvider({ children, businessId }: CategoryProviderProps) {
  const [state, dispatch] = useReducer(categoryReducer, initialState);
  const [trackCategoryViewMutation] = useTrackCategoryViewMutation();

  // API queries
  const categoriesQuery = useGetCategoriesQuery({
    businessId,
    pagination: { page: 1, pageSize: 100 },
    sort: { field: 'sortOrder', order: 'asc' },
  });

  const rootCategoriesQuery = useGetRootCategoriesQuery({
    businessId,
  });

  // Action handlers
  const setCurrentCategory = (id: string | null) => {
    dispatch({ type: 'SET_CURRENT_CATEGORY', payload: id });
  };

  const setViewMode = (mode: 'grid' | 'list') => {
    dispatch({ type: 'SET_VIEW_MODE', payload: mode });
  };

  const setSortBy = (sort: SortOptions) => {
    dispatch({ type: 'SET_SORT_BY', payload: sort });
  };

  const setItemsPerPage = (count: number) => {
    dispatch({ type: 'SET_ITEMS_PER_PAGE', payload: count });
  };

  const toggleFilters = () => {
    dispatch({ type: 'TOGGLE_FILTERS' });
  };

  const setSearchQuery = (query: string) => {
    dispatch({ type: 'SET_SEARCH_QUERY', payload: query });

    if (query.length > 0) {
      dispatch({ type: 'SET_SEARCHING', payload: true });

      // Simulate search suggestions (in real app, this would be an API call)
      setTimeout(() => {
        const suggestions = categoriesQuery.data?.categories.filter(cat =>
          cat.name.toLowerCase().includes(query.toLowerCase())
        ) || [];
        dispatch({ type: 'SET_SEARCH_SUGGESTIONS', payload: suggestions });
        dispatch({ type: 'SET_SEARCHING', payload: false });
      }, 300);
    } else {
      dispatch({ type: 'SET_SEARCH_SUGGESTIONS', payload: [] });
      dispatch({ type: 'SET_SEARCHING', payload: false });
    }
  };

  const addRecentSearch = (query: string) => {
    if (query.trim()) {
      dispatch({ type: 'ADD_RECENT_SEARCH', payload: query.trim() });
    }
  };

  const clearRecentSearches = () => {
    dispatch({ type: 'CLEAR_RECENT_SEARCHES' });
  };

  const updateFilters = (filters: Partial<CategoryFilterState>) => {
    dispatch({ type: 'UPDATE_FILTERS', payload: filters });
  };

  const clearFilters = () => {
    dispatch({ type: 'CLEAR_FILTERS' });
  };

  const trackCategoryView = async (categoryId: string) => {
    dispatch({ type: 'TRACK_CATEGORY_VIEW', payload: categoryId });

    try {
      await trackCategoryViewMutation({
        categoryId,
        sessionId: `session-${state.sessionStartTime}`,
        ipAddress: '0.0.0.0', // Would be determined by backend
        userAgent: navigator.userAgent,
        referrer: document.referrer,
      }).unwrap();
    } catch (error) {
      console.warn('Failed to track category view:', error);
    }
  };

  // Update breadcrumbs when current category changes
  useEffect(() => {
    if (!state.currentCategoryId || !categoriesQuery.data) return;

    const buildBreadcrumbs = (categoryId: string): CategoryBreadcrumb[] => {
      const category = categoriesQuery.data.categories.find(c => c.id === categoryId);
      if (!category) return [];

      const breadcrumbs: CategoryBreadcrumb[] = [{
        id: category.id,
        name: category.name,
        slug: category.seo.slug,
      }];

      if (category.parentId) {
        const parentBreadcrumbs = buildBreadcrumbs(category.parentId);
        return [...parentBreadcrumbs, ...breadcrumbs];
      }

      return breadcrumbs;
    };

    const breadcrumbs = buildBreadcrumbs(state.currentCategoryId);
    dispatch({ type: 'UPDATE_BREADCRUMBS', payload: breadcrumbs });

    // Set parent and child categories
    const currentCategory = categoriesQuery.data.categories.find(c => c.id === state.currentCategoryId);
    if (currentCategory) {
      const parentCategories = categoriesQuery.data.categories.filter(c =>
        c.id === currentCategory.parentId
      );
      const childCategories = categoriesQuery.data.categories.filter(c =>
        c.parentId === currentCategory.id
      );

      dispatch({ type: 'SET_PARENT_CATEGORIES', payload: parentCategories });
      dispatch({ type: 'SET_CHILD_CATEGORIES', payload: childCategories });
    }
  }, [state.currentCategoryId, categoriesQuery.data]);

  // Computed values
  const filteredProducts = useMemo((): ProductFilters => {
    const filters: ProductFilters = {};

    if (state.currentCategoryId) {
      filters.category = state.currentCategoryId;
    }

    if (state.search.query) {
      filters.search = state.search.query;
    }

    return filters;
  }, [state.currentCategoryId, state.search.query]);

  const hasActiveFilters = state.activeFilters.length > 0 || !!state.search.query;

  const searchResults = useMemo(() => {
    if (!state.search.query || !categoriesQuery.data) return [];

    return categoriesQuery.data.categories.filter(category =>
      category.name.toLowerCase().includes(state.search.query.toLowerCase()) ||
      category.description.toLowerCase().includes(state.search.query.toLowerCase())
    );
  }, [state.search.query, categoriesQuery.data]);

  const contextValue: CategoryContextValue = {
    state,
    setCurrentCategory,
    setViewMode,
    setSortBy,
    setItemsPerPage,
    toggleFilters,
    setSearchQuery,
    addRecentSearch,
    clearRecentSearches,
    updateFilters,
    clearFilters,
    trackCategoryView,
    categoriesQuery,
    rootCategoriesQuery,
    filteredProducts,
    hasActiveFilters,
    searchResults,
  };

  return (
    <CategoryContext.Provider value={contextValue}>
      {children}
    </CategoryContext.Provider>
  );
}

// ===== Hook =====

export function useCategory() {
  const context = useContext(CategoryContext);
  if (!context) {
    throw new Error('useCategory must be used within a CategoryProvider');
  }
  return context;
}

// ===== Helper Functions =====

export function buildCategoryPath(category: Category, allCategories: Category[]): Category[] {
  const path: Category[] = [category];

  let currentCategory = category;
  while (currentCategory.parentId) {
    const parent = allCategories.find(c => c.id === currentCategory.parentId);
    if (!parent) break;
    path.unshift(parent);
    currentCategory = parent;
  }

  return path;
}

export function getCategoryChildren(categoryId: string, allCategories: Category[]): Category[] {
  return allCategories.filter(category => category.parentId === categoryId);
}

export function getCategoryDepth(category: Category, allCategories: Category[]): number {
  let depth = 0;
  let currentCategory = category;

  while (currentCategory.parentId) {
    const parent = allCategories.find(c => c.id === currentCategory.parentId);
    if (!parent) break;
    depth++;
    currentCategory = parent;
  }

  return depth;
}