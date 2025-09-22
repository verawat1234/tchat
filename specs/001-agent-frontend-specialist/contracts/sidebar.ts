/**
 * Sidebar Component Contracts
 */

import { BaseComponent, ChangeHandler, ClickHandler, AccessibilityProps } from './base';

/**
 * Sidebar component properties
 */
export interface SidebarProps extends BaseComponent, AccessibilityProps {
  /** Sidebar position */
  position?: 'left' | 'right';
  /** Sidebar width */
  width?: number | string;
  /** Whether sidebar is collapsible */
  collapsible?: boolean;
  /** Controlled collapsed state */
  collapsed?: boolean;
  /** Collapsed state change handler */
  onCollapsedChange?: ChangeHandler<boolean>;
  /** Whether to show overlay on mobile */
  overlay?: boolean;
  /** Whether sidebar persists across breakpoints */
  persistent?: boolean;
  /** Sidebar variant */
  variant?: 'default' | 'floating' | 'bordered';
  /** Whether to show resize handle */
  resizable?: boolean;
  /** Minimum width when resizing */
  minWidth?: number;
  /** Maximum width when resizing */
  maxWidth?: number;
  /** Resize handler */
  onResize?: ChangeHandler<number>;
}

/**
 * Sidebar content properties
 */
export interface SidebarContentProps extends BaseComponent {
  /** Content padding */
  padding?: 'none' | 'sm' | 'md' | 'lg';
  /** Whether content is scrollable */
  scrollable?: boolean;
}

/**
 * Sidebar header properties
 */
export interface SidebarHeaderProps extends BaseComponent {
  /** Header title */
  title?: string;
  /** Header actions */
  actions?: React.ReactNode;
  /** Whether header has border */
  border?: boolean;
  /** Whether header is sticky */
  sticky?: boolean;
}

/**
 * Sidebar footer properties
 */
export interface SidebarFooterProps extends BaseComponent {
  /** Whether footer has border */
  border?: boolean;
  /** Whether footer is sticky */
  sticky?: boolean;
}

/**
 * Sidebar navigation item
 */
export interface SidebarNavItem {
  /** Unique item ID */
  id: string;
  /** Item label */
  label: string;
  /** Item icon */
  icon?: React.ReactNode;
  /** Link href */
  href?: string;
  /** Click handler */
  onClick?: ClickHandler;
  /** Whether item is disabled */
  disabled?: boolean;
  /** Badge content */
  badge?: string | number;
  /** Child items for nested navigation */
  children?: SidebarNavItem[];
  /** Whether item is expandable */
  expandable?: boolean;
  /** Default expanded state */
  defaultExpanded?: boolean;
}

/**
 * Sidebar navigation properties
 */
export interface SidebarNavProps extends BaseComponent {
  /** Navigation items */
  items: SidebarNavItem[];
  /** Currently active item */
  value?: string;
  /** Active item change handler */
  onValueChange?: ChangeHandler<string>;
  /** Whether to allow multiple expanded sections */
  multiple?: boolean;
  /** Default expanded items */
  defaultExpanded?: string[];
  /** Expanded items (controlled) */
  expanded?: string[];
  /** Expanded change handler */
  onExpandedChange?: ChangeHandler<string[]>;
  /** Navigation variant */
  variant?: 'default' | 'pills' | 'tree';
}

/**
 * Sidebar navigation group properties
 */
export interface SidebarNavGroupProps extends BaseComponent {
  /** Group title */
  title?: string;
  /** Whether group is collapsible */
  collapsible?: boolean;
  /** Default collapsed state */
  defaultCollapsed?: boolean;
  /** Controlled collapsed state */
  collapsed?: boolean;
  /** Collapsed change handler */
  onCollapsedChange?: ChangeHandler<boolean>;
}

/**
 * Sidebar toggle button properties
 */
export interface SidebarToggleProps extends BaseComponent {
  /** Current collapsed state */
  collapsed?: boolean;
  /** Toggle handler */
  onToggle?: ClickHandler;
  /** Button position */
  position?: 'inside' | 'outside';
  /** Toggle direction */
  direction?: 'left' | 'right';
}