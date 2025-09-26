/**
 * TchatButton Component - Web Implementation
 * Cross-platform design system component with 5 sophisticated variants
 * Constitutional requirements: 97% visual consistency, WCAG 2.1 AA, <200ms load time
 */
import React, { forwardRef } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '../lib/utils';

// Button variants using class-variance-authority for type safety
const buttonVariants = cva(
  // Base styles - shared across all variants
  [
    'inline-flex items-center justify-center gap-2',
    'rounded-lg font-medium transition-all duration-200',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2',
    'disabled:pointer-events-none disabled:opacity-60',
    'relative overflow-hidden',
    // Touch target compliance (Constitutional requirement: 44dp minimum)
    'min-h-[44px]',
    // Animation support for 60fps performance
    'transform-gpu will-change-transform',
  ],
  {
    variants: {
      // 5 Sophisticated Variants per specification
      variant: {
        primary: [
          'bg-primary text-white shadow-sm',
          'hover:bg-primary/90 hover:shadow-md',
          'active:bg-primary/95 active:scale-[0.98]',
          'disabled:bg-primary/50',
        ],
        secondary: [
          'bg-surface border border-border text-text-primary shadow-sm',
          'hover:bg-surface/80 hover:shadow-md hover:border-border-hover',
          'active:bg-surface/90 active:scale-[0.98]',
          'disabled:bg-surface/50 disabled:border-border/50',
        ],
        ghost: [
          'text-primary hover:bg-primary/10 hover:text-primary',
          'active:bg-primary/20 active:scale-[0.98]',
          'disabled:text-primary/50',
        ],
        destructive: [
          'bg-error text-white shadow-sm',
          'hover:bg-error/90 hover:shadow-md',
          'active:bg-error/95 active:scale-[0.98]',
          'disabled:bg-error/50',
        ],
        outline: [
          'border border-border text-text-primary bg-transparent',
          'hover:bg-surface hover:border-border-hover',
          'active:bg-surface/80 active:scale-[0.98]',
          'disabled:border-border/50 disabled:text-text-primary/50',
        ],
      },
      // 3 Size variants for different use cases
      size: {
        sm: ['h-8 px-3 text-sm min-w-[64px]'], // Small: 32dp height
        md: ['h-11 px-4 text-base min-w-[88px]'], // Medium: 44dp height (default)
        lg: ['h-12 px-6 text-lg min-w-[112px]'], // Large: 48dp height
      },
    },
    defaultVariants: {
      variant: 'primary',
      size: 'md',
    },
  }
);

// Loading spinner component for loading state
const LoadingSpinner = ({ size }: { size: 'sm' | 'md' | 'lg' }) => {
  const spinnerSize = {
    sm: 'w-3 h-3',
    md: 'w-4 h-4',
    lg: 'w-5 h-5',
  }[size];

  return (
    <svg
      className={cn('animate-spin text-current', spinnerSize)}
      fill="none"
      viewBox="0 0 24 24"
      data-testid="loading-spinner"
    >
      <circle
        className="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth="4"
      />
      <path
        className="opacity-75"
        fill="currentColor"
        d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  );
};

export interface TchatButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  /** Loading state - shows spinner and disables interaction */
  loading?: boolean;
  /** Left icon element */
  leftIcon?: React.ReactNode;
  /** Right icon element */
  rightIcon?: React.ReactNode;
  /** Custom class name */
  className?: string;
  /** Children content */
  children?: React.ReactNode;
}

/**
 * TchatButton - Cross-platform design system button component
 *
 * Features:
 * - 5 sophisticated variants (primary, secondary, ghost, destructive, outline)
 * - 3 size variants (sm, md, lg) with proper touch targets
 * - Loading state with animated spinner
 * - Icon support (left/right positioning)
 * - WCAG 2.1 AA accessibility compliance
 * - 60fps animations with GPU acceleration
 * - Cross-platform visual consistency (97% target)
 */
export const TchatButton = forwardRef<HTMLButtonElement, TchatButtonProps>(
  (
    {
      variant,
      size = 'md',
      loading = false,
      leftIcon,
      rightIcon,
      className,
      disabled,
      children,
      ...props
    },
    ref
  ) => {
    // Disable button when loading
    const isDisabled = disabled || loading;

    return (
      <button
        ref={ref}
        className={cn(buttonVariants({ variant, size }), className)}
        disabled={isDisabled}
        data-testid="tchat-button"
        data-variant={variant}
        data-size={size}
        data-loading={loading}
        {...props}
      >
        {/* Left icon or loading spinner */}
        {loading ? (
          <LoadingSpinner size={size} />
        ) : (
          leftIcon && (
            <span className="flex-shrink-0" data-testid="button-left-icon">
              {leftIcon}
            </span>
          )
        )}

        {/* Button text content */}
        {children && (
          <span
            className={cn(
              'flex-1 truncate',
              loading && 'opacity-70'
            )}
            data-testid="button-content"
          >
            {children}
          </span>
        )}

        {/* Right icon (hidden during loading) */}
        {!loading && rightIcon && (
          <span className="flex-shrink-0" data-testid="button-right-icon">
            {rightIcon}
          </span>
        )}
      </button>
    );
  }
);

TchatButton.displayName = 'TchatButton';

// Export variants type for external use
export type TchatButtonVariant = VariantProps<typeof buttonVariants>['variant'];
export type TchatButtonSize = VariantProps<typeof buttonVariants>['size'];