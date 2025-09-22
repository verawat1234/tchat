/**
 * Layout Component Contracts
 */

import { BaseComponent, AccessibilityProps } from './base';

/**
 * Flex layout properties
 */
export interface FlexProps extends BaseComponent {
  /** Flex direction */
  direction?: 'row' | 'column' | 'row-reverse' | 'column-reverse';
  /** Gap between items */
  gap?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl' | number;
  /** Align items */
  align?: 'start' | 'center' | 'end' | 'stretch' | 'baseline';
  /** Justify content */
  justify?: 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly';
  /** Whether to wrap items */
  wrap?: boolean | 'reverse';
  /** Whether to grow to fill container */
  grow?: boolean;
  /** Whether to shrink */
  shrink?: boolean;
}

/**
 * Grid layout properties
 */
export interface GridProps extends BaseComponent {
  /** Number of columns or auto */
  cols?: number | 'auto' | 'fit' | 'fill';
  /** Number of rows or auto */
  rows?: number | 'auto' | 'fit' | 'fill';
  /** Gap between grid items */
  gap?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl' | number;
  /** Column gap specifically */
  gapX?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl' | number;
  /** Row gap specifically */
  gapY?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl' | number;
  /** Responsive breakpoints */
  responsive?: {
    sm?: Partial<GridProps>;
    md?: Partial<GridProps>;
    lg?: Partial<GridProps>;
    xl?: Partial<GridProps>;
  };
}

/**
 * Stack layout properties (vertical spacing)
 */
export interface StackProps extends BaseComponent {
  /** Space between items */
  space?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl' | number;
  /** Divider element between items */
  divider?: React.ReactNode;
  /** Alignment of items */
  align?: 'start' | 'center' | 'end' | 'stretch';
}

/**
 * Container properties
 */
export interface ContainerProps extends BaseComponent {
  /** Maximum width */
  maxWidth?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full' | number;
  /** Whether to center the container */
  center?: boolean;
  /** Padding */
  padding?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  /** Whether container is fluid (full width) */
  fluid?: boolean;
}

/**
 * Spacer properties (for adding space between elements)
 */
export interface SpacerProps extends BaseComponent {
  /** Size of the spacer */
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | number;
  /** Direction of the spacer */
  direction?: 'horizontal' | 'vertical';
}

/**
 * Divider properties
 */
export interface DividerProps extends BaseComponent {
  /** Divider orientation */
  orientation?: 'horizontal' | 'vertical';
  /** Divider style */
  variant?: 'solid' | 'dashed' | 'dotted';
  /** Thickness of the divider */
  thickness?: 'thin' | 'medium' | 'thick';
  /** Label to show on the divider */
  label?: string;
  /** Position of the label */
  labelPosition?: 'left' | 'center' | 'right';
}