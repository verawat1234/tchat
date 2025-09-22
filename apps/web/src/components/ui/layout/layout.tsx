/**
 * Layout Components
 * Flexible layout primitives for building responsive interfaces
 */

import React from 'react';
import { cn } from '@/utils/cn';
import type {
  ContainerProps,
  GridProps,
  FlexProps,
  StackProps,
  SpacerProps,
  DividerProps
} from '../../../../specs/001-agent-frontend-specialist/contracts/layout';

/**
 * Container component for responsive content containment
 */
export const Container = React.forwardRef<HTMLDivElement, ContainerProps>(
  ({
    className,
    children,
    testId,
    maxWidth = 'full',
    center = false,
    padding,
    fluid = false,
    ...props
  }, ref) => {
    const getMaxWidthClass = () => {
      if (typeof maxWidth === 'number') return '';
      const widthMap = {
        xs: 'max-w-xs',
        sm: 'max-w-sm',
        md: 'max-w-md',
        lg: 'max-w-lg',
        xl: 'max-w-xl',
        '2xl': 'max-w-2xl',
        full: 'max-w-full'
      };
      return widthMap[maxWidth];
    };

    const getPaddingClass = () => {
      if (!padding) return '';
      const paddingMap = {
        none: 'p-0',
        xs: 'p-2',
        sm: 'p-4',
        md: 'p-6',
        lg: 'p-8',
        xl: 'p-12'
      };
      return paddingMap[padding];
    };

    const style = typeof maxWidth === 'number' ? { maxWidth: `${maxWidth}px` } : {};

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'w-full',

          // Max width
          getMaxWidthClass(),
          typeof maxWidth !== 'number' && `container-${maxWidth}`,

          // Center alignment
          center && 'mx-auto container-center',

          // Fluid styling
          fluid && 'container-fluid',

          // Padding
          getPaddingClass(),
          padding && `container-padding-${padding}`,

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

Container.displayName = 'Container';

/**
 * Grid component for CSS Grid layouts
 */
export const Grid = React.forwardRef<HTMLDivElement, GridProps>(
  ({
    className,
    children,
    testId,
    cols,
    rows,
    gap,
    gapX,
    gapY,
    responsive,
    ...props
  }, ref) => {
    const getColsClass = () => {
      if (!cols) return '';
      if (typeof cols === 'number') return `grid-cols-${cols}`;
      return `grid-cols-${cols}`;
    };

    const getRowsClass = () => {
      if (!rows) return '';
      if (typeof rows === 'number') return `grid-rows-${rows}`;
      return `grid-rows-${rows}`;
    };

    const getGapClass = () => {
      if (!gap) return '';
      if (typeof gap === 'number') return '';
      const gapMap = {
        none: 'gap-0',
        xs: 'gap-1',
        sm: 'gap-2',
        md: 'gap-4',
        lg: 'gap-6',
        xl: 'gap-8'
      };
      return gapMap[gap];
    };

    const getGapXClass = () => {
      if (!gapX) return '';
      if (typeof gapX === 'number') return '';
      const gapMap = {
        none: 'gap-x-0',
        xs: 'gap-x-1',
        sm: 'gap-x-2',
        md: 'gap-x-4',
        lg: 'gap-x-6',
        xl: 'gap-x-8'
      };
      return gapMap[gapX];
    };

    const getGapYClass = () => {
      if (!gapY) return '';
      if (typeof gapY === 'number') return '';
      const gapMap = {
        none: 'gap-y-0',
        xs: 'gap-y-1',
        sm: 'gap-y-2',
        md: 'gap-y-4',
        lg: 'gap-y-6',
        xl: 'gap-y-8'
      };
      return gapMap[gapY];
    };

    const getResponsiveClasses = () => {
      if (!responsive) return '';
      const classes: string[] = [];

      Object.entries(responsive).forEach(([breakpoint, config]) => {
        if (config.cols) {
          classes.push(`grid-${breakpoint}-cols-${config.cols}`);
        }
      });

      return classes.join(' ');
    };

    const style: React.CSSProperties = {};
    if (typeof gap === 'number') style.gap = `${gap}px`;

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'grid',

          // Columns and rows
          getColsClass(),
          getRowsClass(),

          // Gap
          getGapClass(),
          getGapXClass(),
          getGapYClass(),
          gap && `grid-gap-${gap}`,
          gapX && `grid-gap-x-${gapX}`,
          gapY && `grid-gap-y-${gapY}`,

          // Responsive
          getResponsiveClasses(),

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

Grid.displayName = 'Grid';

/**
 * Flex component for flexbox layouts
 */
export const Flex = React.forwardRef<HTMLDivElement, FlexProps>(
  ({
    className,
    children,
    testId,
    direction = 'row',
    gap,
    align,
    justify,
    wrap,
    grow = false,
    shrink = false,
    ...props
  }, ref) => {
    const getDirectionClass = () => {
      const directionMap = {
        row: 'flex-row',
        column: 'flex-col',
        'row-reverse': 'flex-row-reverse',
        'column-reverse': 'flex-col-reverse'
      };
      return directionMap[direction];
    };

    const getGapClass = () => {
      if (!gap) return '';
      if (typeof gap === 'number') return '';
      const gapMap = {
        none: 'gap-0',
        xs: 'gap-1',
        sm: 'gap-2',
        md: 'gap-4',
        lg: 'gap-6',
        xl: 'gap-8'
      };
      return gapMap[gap];
    };

    const getAlignClass = () => {
      if (!align) return '';
      const alignMap = {
        start: 'items-start',
        center: 'items-center',
        end: 'items-end',
        stretch: 'items-stretch',
        baseline: 'items-baseline'
      };
      return alignMap[align];
    };

    const getJustifyClass = () => {
      if (!justify) return '';
      const justifyMap = {
        start: 'justify-start',
        center: 'justify-center',
        end: 'justify-end',
        between: 'justify-between',
        around: 'justify-around',
        evenly: 'justify-evenly'
      };
      return justifyMap[justify];
    };

    const getWrapClass = () => {
      if (wrap === true) return 'flex-wrap';
      if (wrap === 'reverse') return 'flex-wrap-reverse';
      return '';
    };

    const style: React.CSSProperties = {};
    if (typeof gap === 'number') style.gap = `${gap}px`;

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex',

          // Direction
          getDirectionClass(),
          `flex-${direction}`,

          // Alignment
          getAlignClass(),
          align && `flex-align-${align}`,

          // Justify
          getJustifyClass(),
          justify && `flex-justify-${justify}`,

          // Wrap
          getWrapClass(),
          wrap === true && 'flex-wrap',
          wrap === 'reverse' && 'flex-wrap-reverse',

          // Gap
          getGapClass(),
          gap && `flex-gap-${gap}`,

          // Grow and shrink
          grow && 'flex-grow flex-grow',
          shrink && 'flex-shrink flex-shrink',

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

Flex.displayName = 'Flex';

/**
 * Stack component for vertical layouts with consistent spacing
 */
export const Stack = React.forwardRef<HTMLDivElement, StackProps>(
  ({
    className,
    children,
    testId,
    space,
    divider,
    align,
    ...props
  }, ref) => {
    const getSpaceClass = () => {
      if (!space) return '';
      if (typeof space === 'number') return '';
      const spaceMap = {
        none: 'space-y-0',
        xs: 'space-y-1',
        sm: 'space-y-2',
        md: 'space-y-4',
        lg: 'space-y-6',
        xl: 'space-y-8'
      };
      return spaceMap[space];
    };

    const getAlignClass = () => {
      if (!align) return '';
      const alignMap = {
        start: 'items-start',
        center: 'items-center',
        end: 'items-end',
        stretch: 'items-stretch'
      };
      return alignMap[align];
    };

    const childrenArray = React.Children.toArray(children);
    const childrenWithDividers = divider
      ? childrenArray.reduce((acc: React.ReactNode[], child, index) => {
          acc.push(child);
          if (index < childrenArray.length - 1) {
            acc.push(React.cloneElement(divider as React.ReactElement, { key: `divider-${index}` }));
          }
          return acc;
        }, [])
      : childrenArray;

    const style: React.CSSProperties = {};
    if (typeof space === 'number') style.gap = `${space}px`;

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex flex-col',

          // Space
          getSpaceClass(),
          space && `stack-space-${space}`,

          // Align
          getAlignClass(),
          align && `stack-align-${align}`,

          // Custom className
          className
        )}
        style={style}
        {...props}
      >
        {childrenWithDividers}
      </div>
    );
  }
);

Stack.displayName = 'Stack';

/**
 * Spacer component for adding space between elements
 */
export const Spacer = React.forwardRef<HTMLDivElement, SpacerProps>(
  ({
    className,
    testId,
    size = 'md',
    direction = 'vertical',
    ...props
  }, ref) => {
    const getSizeClass = () => {
      if (typeof size === 'number') return '';
      const sizeMap = {
        xs: direction === 'horizontal' ? 'w-1' : 'h-1',
        sm: direction === 'horizontal' ? 'w-2' : 'h-2',
        md: direction === 'horizontal' ? 'w-4' : 'h-4',
        lg: direction === 'horizontal' ? 'w-6' : 'h-6',
        xl: direction === 'horizontal' ? 'w-8' : 'h-8'
      };
      return sizeMap[size];
    };

    const style: React.CSSProperties = {};
    if (typeof size === 'number') {
      if (direction === 'horizontal') {
        style.width = `${size}px`;
      } else {
        style.height = `${size}px`;
      }
    }

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex-shrink-0',

          // Size
          getSizeClass(),
          `spacer-${size}`,

          // Direction
          `spacer-${direction}`,

          // Custom className
          className
        )}
        style={style}
        {...props}
      />
    );
  }
);

Spacer.displayName = 'Spacer';

/**
 * Divider component for visual separation
 */
export const Divider = React.forwardRef<HTMLDivElement, DividerProps>(
  ({
    className,
    testId,
    orientation = 'horizontal',
    variant = 'solid',
    thickness = 'medium',
    label,
    labelPosition = 'center',
    ...props
  }, ref) => {
    const getOrientationClass = () => {
      return orientation === 'horizontal' ? 'w-full h-px' : 'h-full w-px';
    };

    const getVariantClass = () => {
      const variantMap = {
        solid: 'border-solid',
        dashed: 'border-dashed',
        dotted: 'border-dotted'
      };
      return variantMap[variant];
    };

    const getThicknessClass = () => {
      const thicknessMap = {
        thin: 'border-t',
        medium: 'border-t-2',
        thick: 'border-t-4'
      };
      return orientation === 'horizontal' ? thicknessMap[thickness] : thicknessMap[thickness].replace('t', 'l');
    };

    if (label) {
      return (
        <div
          ref={ref}
          data-testid={testId}
          className={cn(
            // Base styles
            'relative flex items-center',

            // Orientation
            `divider-${orientation}`,

            // Label position
            `divider-label-${labelPosition}`,

            // Custom className
            className
          )}
          {...props}
        >
          <div
            className={cn(
              'border-gray-300',
              getVariantClass(),
              getThicknessClass(),
              labelPosition === 'left' ? 'w-4' : 'flex-1'
            )}
          />
          <span className={cn(
            'px-3 text-sm text-gray-500',
            labelPosition === 'center' && 'mx-3'
          )}>
            {label}
          </span>
          <div
            className={cn(
              'border-gray-300',
              getVariantClass(),
              getThicknessClass(),
              labelPosition === 'right' ? 'w-4' : 'flex-1'
            )}
          />
        </div>
      );
    }

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'border-gray-300',

          // Orientation
          getOrientationClass(),
          `divider-${orientation}`,

          // Variant
          getVariantClass(),
          `divider-${variant}`,

          // Thickness
          getThicknessClass(),
          `divider-${thickness}`,

          // Custom className
          className
        )}
        {...props}
      />
    );
  }
);

Divider.displayName = 'Divider';

export type {
  ContainerProps,
  GridProps,
  FlexProps,
  StackProps,
  SpacerProps,
  DividerProps
};