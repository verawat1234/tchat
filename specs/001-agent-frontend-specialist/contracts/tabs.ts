/**
 * Tabs Component Contracts
 */

import { BaseComponent, ChangeHandler, AccessibilityProps } from './base';

/**
 * Root tabs container properties
 */
export interface TabsProps extends BaseComponent, AccessibilityProps {
  /** Default active tab value (uncontrolled) */
  defaultValue?: string;
  /** Active tab value (controlled) */
  value?: string;
  /** Callback when active tab changes */
  onValueChange?: ChangeHandler<string>;
  /** Tab orientation */
  orientation?: 'horizontal' | 'vertical';
  /** How tabs are activated */
  activationMode?: 'automatic' | 'manual';
  /** Whether tabs can be looped through */
  loop?: boolean;
}

/**
 * Tabs list container properties
 */
export interface TabsListProps extends BaseComponent, AccessibilityProps {
  /** Whether keyboard navigation loops */
  loop?: boolean;
}

/**
 * Individual tab trigger properties
 */
export interface TabsTriggerProps extends BaseComponent, AccessibilityProps {
  /** Unique value for this tab */
  value: string;
  /** Whether tab is disabled */
  disabled?: boolean;
  /** Icon to display */
  icon?: React.ReactNode;
  /** Badge content */
  badge?: string | number;
}

/**
 * Tab content panel properties
 */
export interface TabsContentProps extends BaseComponent, AccessibilityProps {
  /** Tab value this content is associated with */
  value: string;
  /** Force mount content even when not active */
  forceMount?: boolean;
}

/**
 * Tabs context value
 */
export interface TabsContextValue {
  /** Currently active tab */
  value?: string;
  /** Function to change active tab */
  onValueChange?: ChangeHandler<string>;
  /** Tab orientation */
  orientation?: 'horizontal' | 'vertical';
  /** Activation mode */
  activationMode?: 'automatic' | 'manual';
  /** Whether navigation loops */
  loop?: boolean;
}