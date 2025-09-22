/**
 * Card Component Contracts
 */

import { BaseComponent, ClickHandler, AccessibilityProps } from './base';

/**
 * Card root container properties
 */
export interface CardProps extends BaseComponent, AccessibilityProps {
  /** Visual variant */
  variant?: 'default' | 'outline' | 'ghost' | 'filled';
  /** Padding size */
  padding?: 'none' | 'sm' | 'md' | 'lg';
  /** Shadow level */
  shadow?: 'none' | 'sm' | 'md' | 'lg';
  /** Whether card is interactive (clickable) */
  interactive?: boolean;
  /** Click handler for interactive cards */
  onClick?: ClickHandler;
  /** Whether card is loading */
  loading?: boolean;
  /** Whether card is selected/active */
  selected?: boolean;
}

/**
 * Card header properties
 */
export interface CardHeaderProps extends BaseComponent {
  /** Header title */
  title?: string;
  /** Header subtitle */
  subtitle?: string;
  /** Action elements (buttons, menu, etc.) */
  actions?: React.ReactNode;
  /** Avatar or icon */
  avatar?: React.ReactNode;
  /** Whether header has border */
  border?: boolean;
}

/**
 * Card content properties
 */
export interface CardContentProps extends BaseComponent {
  /** Padding override (inherits from Card by default) */
  padding?: 'inherit' | 'none' | 'sm' | 'md' | 'lg';
  /** Whether content is scrollable */
  scrollable?: boolean;
  /** Maximum height for scrollable content */
  maxHeight?: string | number;
}

/**
 * Card footer properties
 */
export interface CardFooterProps extends BaseComponent {
  /** Content justification */
  justify?: 'start' | 'center' | 'end' | 'between' | 'around';
  /** Whether footer has border */
  border?: boolean;
  /** Padding override */
  padding?: 'inherit' | 'none' | 'sm' | 'md' | 'lg';
}

/**
 * Card media properties (for image/video content)
 */
export interface CardMediaProps extends BaseComponent {
  /** Media source URL */
  src: string;
  /** Alternative text */
  alt?: string;
  /** Media type */
  type?: 'image' | 'video';
  /** Aspect ratio */
  aspectRatio?: 'square' | 'video' | 'auto';
  /** Object fit */
  objectFit?: 'cover' | 'contain' | 'fill' | 'none';
}