/**
 * Commerce Components - Category Browsing System
 *
 * Complete category browsing and product discovery system for e-commerce.
 * Provides hierarchical navigation, search, filtering, and product display.
 *
 * Features:
 * - Comprehensive category management
 * - Advanced search and filtering
 * - Multiple view modes and layouts
 * - Analytics integration
 * - Accessibility compliant
 * - Mobile responsive
 * - Performance optimized
 */

// ===== Core Provider =====
export { CategoryProvider, useCategory } from './CategoryProvider';
export type {
  CategorySearchState,
  CategoryFilterState,
  CategoryBreadcrumb,
  CategoryViewState,
  CategoryState,
} from './CategoryProvider';

// ===== Navigation Components =====
export { CategoryTree, CompactCategoryTree } from './CategoryTree';
export { default as CategoryBreadcrumbs, CompactCategoryBreadcrumbs } from './CategoryBreadcrumbs';

// ===== Display Components =====
export { CategoryCard, CategoryCardSkeleton } from './CategoryCard';
export type { CategoryCardProps } from './CategoryCard';

export { CategoryGrid } from './CategoryGrid';
export type { CategoryGridProps } from './CategoryGrid';

export { CategoryProductList } from './CategoryProductList';
export type { CategoryProductListProps } from './CategoryProductList';

// ===== Search & Filter Components =====
export { CategorySearch } from './CategorySearch';
export type { CategorySearchProps } from './CategorySearch';

export { CategoryFilters } from './CategoryFilters';
export type { CategoryFiltersProps, FilterState } from './CategoryFilters';

// ===== Helper Functions =====
export {
  buildCategoryPath,
  getCategoryChildren,
  getCategoryDepth,
} from './CategoryProvider';

// ===== Component Combinations =====

/**
 * Complete category browsing layout with sidebar navigation
 */
export const CategoryBrowsingLayout = {
  Provider: CategoryProvider,
  Sidebar: CategoryTree,
  Search: CategorySearch,
  Filters: CategoryFilters,
  Breadcrumbs: CategoryBreadcrumbs,
  Grid: CategoryGrid,
  ProductList: CategoryProductList,
};

/**
 * Compact category navigation for mobile
 */
export const CompactCategoryLayout = {
  Provider: CategoryProvider,
  Navigation: CompactCategoryTree,
  Search: CategorySearch,
  Breadcrumbs: CompactCategoryBreadcrumbs,
  Grid: CategoryGrid,
};

/**
 * Category showcase for homepage
 */
export const CategoryShowcase = {
  Provider: CategoryProvider,
  Cards: CategoryCard,
  Grid: CategoryGrid,
};

// ===== Usage Examples =====

/**
 * Basic category browsing setup:
 *
 * ```tsx
 * import { CategoryProvider, CategoryTree, CategoryGrid, CategorySearch } from '@/components/commerce';
 *
 * function CategoryBrowser() {
 *   return (
 *     <CategoryProvider businessId="your-business-id">
 *       <div className="flex gap-6">
 *         <aside className="w-80">
 *           <CategorySearch className="mb-4" />
 *           <CategoryTree />
 *         </aside>
 *         <main className="flex-1">
 *           <CategoryGrid />
 *         </main>
 *       </div>
 *     </CategoryProvider>
 *   );
 * }
 * ```
 *
 * Product browsing with filters:
 *
 * ```tsx
 * import { CategoryProvider, CategoryBreadcrumbs, CategoryFilters, CategoryProductList } from '@/components/commerce';
 *
 * function ProductBrowser({ categoryId }: { categoryId: string }) {
 *   return (
 *     <CategoryProvider>
 *       <div className="space-y-4">
 *         <CategoryBreadcrumbs />
 *         <div className="flex gap-6">
 *           <aside className="w-80">
 *             <CategoryFilters />
 *           </aside>
 *           <main className="flex-1">
 *             <CategoryProductList categoryId={categoryId} />
 *           </main>
 *         </div>
 *       </div>
 *     </CategoryProvider>
 *   );
 * }
 * ```
 *
 * Mobile-optimized layout:
 *
 * ```tsx
 * import { CategoryProvider, CategorySearch, CategoryGrid, CategoryFilters } from '@/components/commerce';
 *
 * function MobileCategoryBrowser() {
 *   return (
 *     <CategoryProvider>
 *       <div className="space-y-4 p-4">
 *         <CategorySearch variant="compact" />
 *         <div className="flex gap-2">
 *           <CategoryFilters variant="modal" />
 *           <div className="flex-1" />
 *         </div>
 *         <CategoryGrid viewMode="list" gridColumns={1} />
 *       </div>
 *     </CategoryProvider>
 *   );
 * }
 * ```
 */