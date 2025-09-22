/**
 * Badge Component Contracts
 */

import { BaseComponent, AccessibilityProps } from './base';

/**
 * Badge component properties
 */
export interface BadgeProps extends BaseComponent, AccessibilityProps {
  /** Badge visual variant */
  variant?: 'default' | 'success' | 'warning' | 'danger' | 'info' | 'secondary';
  /** Badge size */
  size?: 'sm' | 'md' | 'lg';
  /** Show as dot indicator */
  dot?: boolean;
  /** Pulse animation */
  pulse?: boolean;
  /** Count/number to display */
  count?: number;
  /** Maximum count before showing "99+" */
  max?: number;
  /** Show badge even when count is 0 */
  showZero?: boolean;
  /** Custom content */
  content?: string | number;
  /** Icon to display */
  icon?: React.ReactNode;
  /** Whether badge is outlined */
  outline?: boolean;
  /** Whether badge is removable */
  removable?: boolean;
  /** Remove handler */
  onRemove?: () => void;
}

/**
 * Badge group properties (for multiple badges)
 */
export interface BadgeGroupProps extends BaseComponent {
  /** Maximum number of badges to show */
  max?: number;
  /** Spacing between badges */
  spacing?: 'tight' | 'normal' | 'loose';
  /** Whether to wrap badges */
  wrap?: boolean;
}