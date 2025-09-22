/**
 * Base Component Contracts
 * Core interfaces and types for all common UI components
 */

import { ReactNode } from 'react';

/**
 * Base properties shared by all components
 */
export interface BaseComponent {
  /** Additional CSS class names */
  className?: string;
  /** Child elements */
  children?: ReactNode;
  /** Test ID for automated testing */
  testId?: string;
  /** Component variant for styling */
  variant?: string;
  /** Component size */
  size?: 'sm' | 'md' | 'lg' | 'xl';
}

/**
 * Theme token structure for consistent styling
 */
export interface ThemeTokens {
  colors: {
    primary: string;
    secondary: string;
    accent: string;
    muted: string;
    background: string;
    foreground: string;
    border: string;
    destructive: string;
    success: string;
    warning: string;
  };
  spacing: {
    xs: string;
    sm: string;
    md: string;
    lg: string;
    xl: string;
  };
  typography: {
    fontSize: Record<string, string>;
    fontWeight: Record<string, number>;
    lineHeight: Record<string, string>;
  };
  borderRadius: {
    none: string;
    sm: string;
    md: string;
    lg: string;
    full: string;
  };
  shadow: {
    none: string;
    sm: string;
    md: string;
    lg: string;
    xl: string;
  };
}

/**
 * Common event handler types
 */
export type ClickHandler = () => void;
export type ChangeHandler<T> = (value: T) => void;
export type KeyboardHandler = (event: KeyboardEvent) => void;

/**
 * Accessibility properties
 */
export interface AccessibilityProps {
  /** ARIA label */
  'aria-label'?: string;
  /** ARIA description */
  'aria-describedby'?: string;
  /** ARIA expanded state */
  'aria-expanded'?: boolean;
  /** ARIA disabled state */
  'aria-disabled'?: boolean;
  /** Role for screen readers */
  role?: string;
  /** Tab index for keyboard navigation */
  tabIndex?: number;
}