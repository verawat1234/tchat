/**
 * Card Components
 * Flexible content containers with various layouts and styling options
 */

import React from 'react';
import { cn } from '@/utils/cn';
import type {
  CardProps,
  CardHeaderProps,
  CardContentProps,
  CardFooterProps,
  CardMediaProps
} from '../../../../specs/001-agent-frontend-specialist/contracts/card';

/**
 * Card root container component
 */
export const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({
    className,
    children,
    testId,
    variant = 'default',
    padding = 'md',
    shadow = 'sm',
    interactive = false,
    onClick,
    loading = false,
    selected = false,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedby,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const getVariantClass = () => {
      const variantMap = {
        default: 'bg-white border border-gray-200',
        outline: 'bg-transparent border border-gray-300',
        ghost: 'bg-transparent border-0',
        filled: 'bg-gray-50 border border-gray-100'
      };
      return variantMap[variant];
    };

    const getPaddingClass = () => {
      const paddingMap = {
        none: 'p-0',
        sm: 'p-3',
        md: 'p-6',
        lg: 'p-8'
      };
      return paddingMap[padding];
    };

    const getShadowClass = () => {
      const shadowMap = {
        none: 'shadow-none',
        sm: 'shadow-sm',
        md: 'shadow-md',
        lg: 'shadow-lg'
      };
      return shadowMap[shadow];
    };

    const handleClick = () => {
      if (interactive && onClick) {
        onClick();
      }
    };

    const handleKeyDown = (event: React.KeyboardEvent) => {
      if (interactive && onClick && (event.key === 'Enter' || event.key === ' ')) {
        event.preventDefault();
        onClick();
      }
    };

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'rounded-lg overflow-hidden',

          // Variant styles
          getVariantClass(),
          `card-${variant}`,

          // Padding styles
          getPaddingClass(),
          `card-padding-${padding}`,

          // Shadow styles
          getShadowClass(),
          `card-shadow-${shadow}`,

          // Interactive styles
          interactive && 'cursor-pointer transition-all duration-200 hover:shadow-md card-interactive',

          // Loading state
          loading && 'opacity-60 pointer-events-none card-loading',

          // Selected state
          selected && 'ring-2 ring-blue-500 card-selected',

          // Custom className
          className
        )}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
        tabIndex={interactive ? (tabIndex ?? 0) : tabIndex}
        aria-label={ariaLabel}
        aria-describedby={ariaDescribedby}
        aria-expanded={ariaExpanded}
        aria-disabled={ariaDisabled}
        aria-busy={loading}
        aria-selected={selected}
        role={role}
        {...props}
      >
        {children}
      </div>
    );
  }
);

Card.displayName = 'Card';

/**
 * Card header component
 */
export const CardHeader = React.forwardRef<HTMLDivElement, CardHeaderProps>(
  ({
    className,
    children,
    testId,
    title,
    subtitle,
    actions,
    avatar,
    border = false,
    ...props
  }, ref) => {
    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex items-center justify-between',

          // Border
          border && 'border-b border-gray-200 card-header-border',

          // Custom className
          className
        )}
        {...props}
      >
        <div className="flex items-center space-x-3">
          {/* Avatar */}
          {avatar && (
            <div className="flex-shrink-0">
              {avatar}
            </div>
          )}

          {/* Title and subtitle */}
          {(title || subtitle) ? (
            <div className="min-w-0 flex-1">
              {title && (
                <h3 className="text-lg font-semibold text-gray-900 truncate">
                  {title}
                </h3>
              )}
              {subtitle && (
                <p className="text-sm text-gray-500 truncate">
                  {subtitle}
                </p>
              )}
            </div>
          ) : children}
        </div>

        {/* Actions */}
        {actions && (
          <div className="flex-shrink-0 ml-4">
            {actions}
          </div>
        )}
      </div>
    );
  }
);

CardHeader.displayName = 'CardHeader';

/**
 * Card content component
 */
export const CardContent = React.forwardRef<HTMLDivElement, CardContentProps>(
  ({
    className,
    children,
    testId,
    padding = 'inherit',
    scrollable = false,
    maxHeight,
    ...props
  }, ref) => {
    const getPaddingClass = () => {
      if (padding === 'inherit') return '';
      const paddingMap = {
        none: 'p-0',
        sm: 'p-3',
        md: 'p-6',
        lg: 'p-8'
      };
      return paddingMap[padding];
    };

    const style: React.CSSProperties = {};
    if (maxHeight) {
      style.maxHeight = typeof maxHeight === 'string' ? maxHeight : `${maxHeight}px`;
    }

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex-1',

          // Padding override
          padding !== 'inherit' && getPaddingClass(),
          padding !== 'inherit' && `card-content-padding-${padding}`,

          // Scrollable
          scrollable && 'overflow-auto card-content-scrollable',

          // Custom className
          className
        )}
        style={style}
        {...props}
      >
        {children}
      </div>
    );
  }
);

CardContent.displayName = 'CardContent';

/**
 * Card footer component
 */
export const CardFooter = React.forwardRef<HTMLDivElement, CardFooterProps>(
  ({
    className,
    children,
    testId,
    justify = 'end',
    border = false,
    padding = 'inherit',
    ...props
  }, ref) => {
    const getJustifyClass = () => {
      const justifyMap = {
        start: 'justify-start',
        center: 'justify-center',
        end: 'justify-end',
        between: 'justify-between',
        around: 'justify-around'
      };
      return justifyMap[justify];
    };

    const getPaddingClass = () => {
      if (padding === 'inherit') return '';
      const paddingMap = {
        none: 'p-0',
        sm: 'p-3',
        md: 'p-6',
        lg: 'p-8'
      };
      return paddingMap[padding];
    };

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex items-center',

          // Justify content
          getJustifyClass(),
          `card-footer-justify-${justify}`,

          // Border
          border && 'border-t border-gray-200 card-footer-border',

          // Padding override
          padding !== 'inherit' && getPaddingClass(),
          padding !== 'inherit' && `card-footer-padding-${padding}`,

          // Custom className
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);

CardFooter.displayName = 'CardFooter';

/**
 * Card title component
 */
export const CardTitle = React.forwardRef<HTMLHeadingElement, React.HTMLAttributes<HTMLHeadingElement>>(
  ({ className, children, ...props }, ref) => (
    <h3
      ref={ref}
      className={cn("text-lg font-medium leading-6 text-gray-900", className)}
      {...props}
    >
      {children}
    </h3>
  )
);

CardTitle.displayName = 'CardTitle';

/**
 * Card description component
 */
export const CardDescription = React.forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  ({ className, children, ...props }, ref) => (
    <p
      ref={ref}
      className={cn("text-sm text-gray-600", className)}
      {...props}
    >
      {children}
    </p>
  )
);

CardDescription.displayName = 'CardDescription';

/**
 * Card media component for images and videos
 */
export const CardMedia = React.forwardRef<HTMLDivElement, CardMediaProps>(
  ({
    className,
    testId,
    src,
    alt,
    type = 'image',
    aspectRatio = 'auto',
    objectFit = 'cover',
    ...props
  }, ref) => {
    const getAspectRatioClass = () => {
      const aspectMap = {
        square: 'aspect-square',
        video: 'aspect-video',
        auto: 'aspect-auto'
      };
      return aspectMap[aspectRatio];
    };

    const getObjectFitClass = () => {
      const fitMap = {
        cover: 'object-cover',
        contain: 'object-contain',
        fill: 'object-fill',
        none: 'object-none'
      };
      return fitMap[objectFit];
    };

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'relative overflow-hidden',

          // Aspect ratio
          getAspectRatioClass(),
          `card-media-aspect-${aspectRatio}`,

          // Object fit
          getObjectFitClass(),
          `card-media-fit-${objectFit}`,

          // Custom className
          className
        )}
        {...props}
      >
        {type === 'image' && (
          <img
            src={src}
            alt={alt}
            className={cn(
              'w-full h-full',
              getObjectFitClass()
            )}
          />
        )}

        {type === 'video' && (
          <video
            src={src}
            className={cn(
              'w-full h-full',
              getObjectFitClass()
            )}
            controls
          />
        )}
      </div>
    );
  }
);

CardMedia.displayName = 'CardMedia';

export type {
  CardProps,
  CardHeaderProps,
  CardContentProps,
  CardFooterProps,
  CardMediaProps
};