/**
 * Badge Component
 * Displays small status descriptors for UI elements
 */

import React from 'react';
import { cn } from '@/utils/cn';
import type { BadgeProps, BadgeGroupProps } from '../../../../specs/001-agent-frontend-specialist/contracts/badge';

const badgeVariants = {
  default: 'bg-gray-100 text-gray-900 border-gray-200',
  success: 'bg-green-100 text-green-800 border-green-200',
  warning: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  danger: 'bg-red-100 text-red-800 border-red-200',
  info: 'bg-blue-100 text-blue-800 border-blue-200',
  secondary: 'bg-gray-100 text-gray-600 border-gray-200'
};

const badgeSizes = {
  sm: 'text-xs px-1.5 py-0.5 min-h-[16px]',
  md: 'text-sm px-2 py-0.5 min-h-[20px]',
  lg: 'text-base px-2.5 py-1 min-h-[24px]'
};

/**
 * Badge component for displaying status, counts, and labels
 */
export const Badge = React.forwardRef<HTMLSpanElement, BadgeProps>(
  ({
    className,
    children,
    testId,
    variant = 'default',
    size = 'md',
    dot = false,
    pulse = false,
    count,
    max = 99,
    showZero = false,
    content,
    icon,
    outline = false,
    removable = false,
    onRemove,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedby,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    // Handle count display logic
    const getDisplayContent = () => {
      if (content !== undefined) return content;
      if (count !== undefined) {
        if (count === 0 && !showZero) return null;
        return count > max ? `${max}+` : count;
      }
      return children;
    };

    const displayContent = getDisplayContent();

    // Don't render if no content and count is 0 without showZero
    if (displayContent === null && !dot) return null;

    return (
      <span
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'inline-flex items-center font-medium rounded-full border',

          // Size styles
          badgeSizes[size],
          `badge-${size}`,

          // Variant styles
          badgeVariants[variant],
          `badge-${variant}`,

          // Outline styles
          outline && 'bg-transparent badge-outline',

          // Dot mode
          dot && 'w-2 h-2 p-0 min-h-0 badge-dot',

          // Pulse animation
          pulse && 'animate-pulse badge-pulse',

          // Interactive styles for removable badges
          removable && 'pr-1',

          // Custom className
          className
        )}
        aria-label={ariaLabel}
        aria-describedby={ariaDescribedby}
        aria-expanded={ariaExpanded}
        aria-disabled={ariaDisabled}
        role={role}
        tabIndex={tabIndex}
        {...props}
      >
        {/* Icon */}
        {icon && !dot && (
          <span className="mr-1 flex items-center">
            {icon}
          </span>
        )}

        {/* Content */}
        {!dot && displayContent}

        {/* Remove button */}
        {removable && onRemove && (
          <button
            type="button"
            className={cn(
              'ml-1 flex items-center justify-center w-4 h-4 rounded-full',
              'hover:bg-black/10 focus:outline-none focus:ring-1 focus:ring-black/20',
              'transition-colors duration-150'
            )}
            onClick={onRemove}
            aria-label="Remove badge"
            tabIndex={-1}
          >
            <svg
              className="w-3 h-3"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fillRule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        )}
      </span>
    );
  }
);

Badge.displayName = 'Badge';

/**
 * BadgeGroup component for displaying multiple badges
 */
export const BadgeGroup = React.forwardRef<HTMLDivElement, BadgeGroupProps>(
  ({
    className,
    children,
    testId,
    max,
    spacing = 'normal',
    wrap = true,
    ...props
  }, ref) => {
    const childrenArray = React.Children.toArray(children);
    const visibleChildren = max ? childrenArray.slice(0, max) : childrenArray;

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'inline-flex items-center',

          // Spacing styles
          spacing === 'tight' && 'gap-1 badge-group-tight',
          spacing === 'normal' && 'gap-2 badge-group-normal',
          spacing === 'loose' && 'gap-3 badge-group-loose',

          // Wrap styles
          wrap && 'flex-wrap badge-group-wrap',

          // Custom className
          className
        )}
        {...props}
      >
        {visibleChildren}

        {/* Show count of hidden items */}
        {max && childrenArray.length > max && (
          <Badge variant="secondary" size="sm">
            +{childrenArray.length - max}
          </Badge>
        )}
      </div>
    );
  }
);

BadgeGroup.displayName = 'BadgeGroup';

export type { BadgeProps, BadgeGroupProps };