/**
 * Pagination Component Contracts
 */

import { BaseComponent, ChangeHandler, AccessibilityProps } from './base';

/**
 * Pagination component properties
 */
export interface PaginationProps extends BaseComponent, AccessibilityProps {
  /** Current active page (1-based) */
  currentPage: number;
  /** Total number of pages */
  totalPages: number;
  /** Callback when page changes */
  onPageChange: ChangeHandler<number>;
  /** Show first/last page buttons */
  showFirstLast?: boolean;
  /** Show previous/next page buttons */
  showPrevNext?: boolean;
  /** Maximum number of visible page buttons */
  maxVisiblePages?: number;
  /** Disable all pagination controls */
  disabled?: boolean;
  /** Show page size selector */
  showPageSize?: boolean;
  /** Available page sizes */
  pageSizes?: number[];
  /** Current page size */
  pageSize?: number;
  /** Callback when page size changes */
  onPageSizeChange?: ChangeHandler<number>;
}

/**
 * Internal pagination state
 */
export interface PaginationState {
  /** Array of visible page numbers */
  visiblePages: number[];
  /** Whether previous navigation is possible */
  canGoPrevious: boolean;
  /** Whether next navigation is possible */
  canGoNext: boolean;
  /** Whether first page button should be shown */
  showFirstEllipsis: boolean;
  /** Whether last page button should be shown */
  showLastEllipsis: boolean;
}

/**
 * Pagination button properties
 */
export interface PaginationButtonProps extends BaseComponent {
  /** Page number or action type */
  page: number | 'first' | 'prev' | 'next' | 'last' | 'ellipsis';
  /** Whether button is active */
  active?: boolean;
  /** Whether button is disabled */
  disabled?: boolean;
  /** Click handler */
  onClick?: () => void;
}

/**
 * Page size selector properties
 */
export interface PageSizeSelectorProps extends BaseComponent {
  /** Available page sizes */
  sizes: number[];
  /** Current page size */
  value: number;
  /** Change handler */
  onChange: ChangeHandler<number>;
  /** Whether selector is disabled */
  disabled?: boolean;
}