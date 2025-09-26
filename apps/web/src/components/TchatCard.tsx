/**
 * TchatCard Component - Web Implementation
 * Cross-platform design system card component with 4 sophisticated variants
 * Constitutional requirements: 97% visual consistency, WCAG 2.1 AA, <200ms load time
 */
import React, { forwardRef } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '../lib/utils';

// Card variants using class-variance-authority for type safety
const cardVariants = cva(
  // Base styles - shared across all variants
  [
    'rounded-lg transition-all duration-200',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2',
    // Touch target compliance when interactive (Constitutional requirement: 44dp minimum)
    'min-h-[44px]',
    // Animation support for 60fps performance
    'transform-gpu will-change-transform',
  ],
  {
    variants: {
      // 4 Sophisticated Variants per specification
      variant: {
        elevated: [
          'bg-white shadow-sm border border-border/50',
          'hover:shadow-md hover:border-border',
          'active:shadow-lg active:scale-[0.98]',
        ],
        outlined: [
          'bg-white border border-border',
          'hover:border-border-hover hover:shadow-sm',
          'active:border-border-hover active:shadow-md active:scale-[0.98]',
        ],
        filled: [
          'bg-surface border border-border/30',
          'hover:bg-surface/80 hover:border-border',
          'active:bg-surface/90 active:scale-[0.98]',
        ],
        glass: [
          'bg-white/80 backdrop-blur-sm border border-white/20',
          'shadow-lg shadow-black/5',
          'hover:bg-white/90 hover:shadow-xl hover:shadow-black/10',
          'active:bg-white/95 active:scale-[0.98]',
          // Glassmorphism specific styles
          'before:absolute before:inset-0 before:rounded-lg',
          'before:bg-gradient-to-br before:from-white/20 before:to-transparent',
          'before:pointer-events-none',
          'relative overflow-hidden',
        ],
      },
      // 3 Size variants for different content densities
      size: {
        compact: ['p-3'], // 12dp padding
        standard: ['p-4'], // 16dp padding (default)
        expanded: ['p-6'], // 24dp padding
      },
      // Interactive state for clickable cards
      interactive: {
        true: [
          'cursor-pointer',
          'focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2',
        ],
        false: [],
      },
    },
    defaultVariants: {
      variant: 'elevated',
      size: 'standard',
      interactive: false,
    },
  }
);

export interface TchatCardProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof cardVariants> {
  /** Card variant affecting visual appearance */
  variant?: 'elevated' | 'outlined' | 'filled' | 'glass';
  /** Size variant affecting padding */
  size?: 'compact' | 'standard' | 'expanded';
  /** Whether the card is interactive (clickable) */
  interactive?: boolean;
  /** Click handler for interactive cards */
  onClick?: React.MouseEventHandler<HTMLDivElement>;
  /** Keyboard handler for interactive cards */
  onKeyDown?: React.KeyboardEventHandler<HTMLDivElement>;
  /** Custom class name */
  className?: string;
  /** Children content */
  children?: React.ReactNode;
  /** ARIA label for accessibility */
  ariaLabel?: string;
  /** Content description for accessibility */
  contentDescription?: string;
  /** Role for semantic markup */
  role?: string;
}

/**
 * TchatCard - Cross-platform design system card component
 *
 * Features:
 * - 4 sophisticated variants (elevated, outlined, filled, glass)
 * - 3 size variants (compact, standard, expanded) with consistent padding
 * - Interactive state with proper focus management
 * - Glassmorphism effect for modern UI aesthetics
 * - WCAG 2.1 AA accessibility compliance
 * - 60fps animations with GPU acceleration
 * - Cross-platform visual consistency (97% target)
 */
export const TchatCard = forwardRef<HTMLDivElement, TchatCardProps>(
  (
    {
      variant = 'elevated',
      size = 'standard',
      interactive = false,
      onClick,
      onKeyDown,
      className,
      children,
      ariaLabel,
      contentDescription,
      role,
      ...props
    },
    ref
  ) => {
    // Handle keyboard interactions for accessibility
    const handleKeyDown = (event: React.KeyboardEvent<HTMLDivElement>) => {
      if (interactive && (event.key === 'Enter' || event.key === ' ')) {
        event.preventDefault();
        onClick?.(event as any);
      }
      onKeyDown?.(event);
    };

    // Determine semantic role
    const cardRole = role || (interactive ? 'button' : 'article');

    // Determine tabIndex for keyboard navigation
    const tabIndex = interactive ? 0 : undefined;

    return (
      <div
        ref={ref}
        className={cn(cardVariants({ variant, size, interactive }), className)}
        onClick={interactive ? onClick : undefined}
        onKeyDown={interactive ? handleKeyDown : onKeyDown}
        role={cardRole}
        tabIndex={tabIndex}
        aria-label={ariaLabel || contentDescription}
        data-testid="tchat-card"
        data-variant={variant}
        data-size={size}
        data-interactive={interactive}
        {...props}
      >
        {children}
      </div>
    );
  }
);

TchatCard.displayName = 'TchatCard';

// Export types for external use
export type TchatCardVariant = VariantProps<typeof cardVariants>['variant'];
export type TchatCardSize = VariantProps<typeof cardVariants>['size'];

/**
 * TchatCardHeader - Optional card header component
 */
export interface TchatCardHeaderProps extends React.HTMLAttributes<HTMLDivElement> {
  title?: string;
  subtitle?: string;
  actions?: React.ReactNode;
  className?: string;
}

export const TchatCardHeader = forwardRef<HTMLDivElement, TchatCardHeaderProps>(
  ({ title, subtitle, actions, className, children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={cn('flex items-start justify-between mb-3', className)}
        data-testid="card-header"
        {...props}
      >
        <div className="flex-1 min-w-0">
          {title && (
            <h3 className="text-lg font-medium text-text-primary truncate mb-1">
              {title}
            </h3>
          )}
          {subtitle && (
            <p className="text-sm text-text-muted truncate">
              {subtitle}
            </p>
          )}
          {children}
        </div>
        {actions && (
          <div className="flex items-center gap-2 ml-3">
            {actions}
          </div>
        )}
      </div>
    );
  }
);

TchatCardHeader.displayName = 'TchatCardHeader';

/**
 * TchatCardContent - Optional card content wrapper
 */
export interface TchatCardContentProps extends React.HTMLAttributes<HTMLDivElement> {
  className?: string;
}

export const TchatCardContent = forwardRef<HTMLDivElement, TchatCardContentProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={cn('text-text-primary', className)}
        data-testid="card-content"
        {...props}
      >
        {children}
      </div>
    );
  }
);

TchatCardContent.displayName = 'TchatCardContent';

/**
 * TchatCardFooter - Optional card footer component
 */
export interface TchatCardFooterProps extends React.HTMLAttributes<HTMLDivElement> {
  className?: string;
}

export const TchatCardFooter = forwardRef<HTMLDivElement, TchatCardFooterProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={cn('flex items-center justify-between mt-4 pt-3 border-t border-border/30', className)}
        data-testid="card-footer"
        {...props}
      >
        {children}
      </div>
    );
  }
);

TchatCardFooter.displayName = 'TchatCardFooter';