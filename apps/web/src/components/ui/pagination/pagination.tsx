/**
 * Pagination Components
 * Navigation controls for paginated content
 */

import React from 'react';
import { cn } from '@/utils/cn';
import type {
  PaginationProps,
  PaginationButtonProps,
  PageSizeSelectorProps,
  PaginationState
} from '../../../../specs/001-agent-frontend-specialist/contracts/pagination';

/**
 * Calculate visible pages for pagination
 */
const calculatePaginationState = (
  currentPage: number,
  totalPages: number,
  maxVisiblePages: number = 7
): PaginationState => {
  const pages: number[] = [];
  let showFirstEllipsis = false;
  let showLastEllipsis = false;

  if (totalPages <= maxVisiblePages) {
    // Show all pages if total is less than max
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i);
    }
  } else {
    const sidePages = Math.floor((maxVisiblePages - 1) / 2);
    let startPage = Math.max(1, currentPage - sidePages);
    let endPage = Math.min(totalPages, currentPage + sidePages);

    // Adjust if we're near the beginning or end
    if (currentPage <= sidePages) {
      endPage = Math.min(totalPages, maxVisiblePages);
    } else if (currentPage >= totalPages - sidePages) {
      startPage = Math.max(1, totalPages - maxVisiblePages + 1);
    }

    // Add pages
    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }

    // Check for ellipsis
    showFirstEllipsis = startPage > 1;
    showLastEllipsis = endPage < totalPages;
  }

  return {
    visiblePages: pages,
    canGoPrevious: currentPage > 1,
    canGoNext: currentPage < totalPages,
    showFirstEllipsis,
    showLastEllipsis
  };
};

/**
 * Individual pagination button component
 */
export const PaginationButton = React.forwardRef<HTMLButtonElement, PaginationButtonProps>(
  ({
    className,
    testId,
    page,
    active = false,
    disabled = false,
    onClick,
    ...props
  }, ref) => {
    const getLabel = () => {
      if (typeof page === 'number') return page.toString();

      const labelMap = {
        first: 'First',
        prev: 'Previous',
        next: 'Next',
        last: 'Last',
        ellipsis: '...'
      };
      return labelMap[page] || page;
    };

    const isEllipsis = page === 'ellipsis';

    if (isEllipsis) {
      return (
        <span
          data-testid={testId}
          className={cn(
            'px-3 py-2 text-sm text-gray-500',
            className
          )}
        >
          ...
        </span>
      );
    }

    return (
      <button
        ref={ref}
        type="button"
        data-testid={testId}
        className={cn(
          // Base styles
          'px-3 py-2 text-sm font-medium rounded-md transition-colors duration-200',
          'border border-gray-300 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500',

          // Active state
          active && 'bg-blue-600 text-white border-blue-600 hover:bg-blue-700',

          // Disabled state
          disabled && 'opacity-50 cursor-not-allowed hover:bg-transparent',

          // Custom className
          className
        )}
        disabled={disabled}
        onClick={onClick}
        aria-current={active ? 'page' : undefined}
        aria-label={getLabel()}
        {...props}
      >
        {getLabel()}
      </button>
    );
  }
);

PaginationButton.displayName = 'PaginationButton';

/**
 * Page size selector component
 */
export const PageSizeSelector = React.forwardRef<HTMLSelectElement, PageSizeSelectorProps>(
  ({
    className,
    testId,
    sizes,
    value,
    onChange,
    disabled = false,
    ...props
  }, ref) => {
    const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
      onChange(parseInt(event.target.value, 10));
    };

    return (
      <select
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'px-3 py-2 text-sm border border-gray-300 rounded-md',
          'bg-white focus:outline-none focus:ring-2 focus:ring-blue-500',

          // Disabled state
          disabled && 'opacity-50 cursor-not-allowed',

          // Custom className
          className
        )}
        value={value}
        onChange={handleChange}
        disabled={disabled}
        {...props}
      >
        {sizes.map(size => (
          <option key={size} value={size}>
            {size}
          </option>
        ))}
      </select>
    );
  }
);

PageSizeSelector.displayName = 'PageSizeSelector';

/**
 * Main pagination component
 */
export const Pagination = React.forwardRef<HTMLNavElement, PaginationProps>(
  ({
    className,
    testId,
    currentPage,
    totalPages,
    onPageChange,
    showFirstLast = false,
    showPrevNext = true,
    maxVisiblePages = 7,
    disabled = false,
    showPageSize = false,
    pageSizes,
    pageSize,
    onPageSizeChange,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedBy,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const paginationState = calculatePaginationState(currentPage, totalPages, maxVisiblePages);

    const handlePageChange = (page: number) => {
      if (!disabled && page !== currentPage && page >= 1 && page <= totalPages) {
        onPageChange(page);
      }
    };

    const handlePrevious = () => {
      handlePageChange(currentPage - 1);
    };

    const handleNext = () => {
      handlePageChange(currentPage + 1);
    };

    const handleFirst = () => {
      handlePageChange(1);
    };

    const handleLast = () => {
      handlePageChange(totalPages);
    };

    return (
      <nav
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex items-center justify-between',

          // Custom className
          className
        )}
        aria-label={ariaLabel || 'Pagination navigation'}
        aria-describedby={ariaDescribedBy}
        aria-expanded={ariaExpanded}
        aria-disabled={ariaDisabled}
        role={role || 'navigation'}
        tabIndex={tabIndex}
        {...props}
      >
        <div className="flex items-center space-x-2">
          {/* First button */}
          {showFirstLast && (
            <PaginationButton
              page="first"
              disabled={disabled || !paginationState.canGoPrevious}
              onClick={handleFirst}
            />
          )}

          {/* Previous button */}
          {showPrevNext && (
            <PaginationButton
              page="prev"
              disabled={disabled || !paginationState.canGoPrevious}
              onClick={handlePrevious}
            />
          )}

          {/* First ellipsis */}
          {paginationState.showFirstEllipsis && (
            <PaginationButton page="ellipsis" />
          )}

          {/* Page buttons */}
          {paginationState.visiblePages.map(page => (
            <PaginationButton
              key={page}
              page={page}
              active={page === currentPage}
              disabled={disabled}
              onClick={() => handlePageChange(page)}
            />
          ))}

          {/* Last ellipsis */}
          {paginationState.showLastEllipsis && (
            <PaginationButton page="ellipsis" />
          )}

          {/* Next button */}
          {showPrevNext && (
            <PaginationButton
              page="next"
              disabled={disabled || !paginationState.canGoNext}
              onClick={handleNext}
            />
          )}

          {/* Last button */}
          {showFirstLast && (
            <PaginationButton
              page="last"
              disabled={disabled || !paginationState.canGoNext}
              onClick={handleLast}
            />
          )}
        </div>

        {/* Page size selector */}
        {showPageSize && pageSizes && pageSize && onPageSizeChange && (
          <div className="flex items-center space-x-2">
            <span className="text-sm text-gray-700">Show:</span>
            <PageSizeSelector
              sizes={pageSizes}
              value={pageSize}
              onChange={onPageSizeChange}
              disabled={disabled}
            />
            <span className="text-sm text-gray-700">per page</span>
          </div>
        )}
      </nav>
    );
  }
);

Pagination.displayName = 'Pagination';

export type {
  PaginationProps,
  PaginationButtonProps,
  PageSizeSelectorProps,
  PaginationState
};