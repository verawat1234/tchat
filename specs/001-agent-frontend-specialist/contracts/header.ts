/**
 * Header Component Contracts
 */

import { BaseComponent, ClickHandler, AccessibilityProps } from './base';

/**
 * Header component properties
 */
export interface HeaderProps extends BaseComponent, AccessibilityProps {
  /** Header title */
  title: string;
  /** Header subtitle */
  subtitle?: string;
  /** Header level (h1-h6) */
  level?: 1 | 2 | 3 | 4 | 5 | 6;
  /** Action elements */
  actions?: React.ReactNode;
  /** Breadcrumb navigation */
  breadcrumbs?: BreadcrumbItem[];
  /** Whether header is sticky */
  sticky?: boolean;
  /** Whether header has bottom border */
  border?: boolean;
  /** Header size variant */
  size?: 'sm' | 'md' | 'lg' | 'xl';
  /** Whether to center content */
  centered?: boolean;
  /** Background variant */
  background?: 'transparent' | 'default' | 'muted';
  /** Icon or avatar */
  icon?: React.ReactNode;
  /** Back button */
  showBack?: boolean;
  /** Back button handler */
  onBack?: ClickHandler;
}

/**
 * Breadcrumb item
 */
export interface BreadcrumbItem {
  /** Item label */
  label: string;
  /** Link href */
  href?: string;
  /** Click handler */
  onClick?: ClickHandler;
  /** Whether item is disabled */
  disabled?: boolean;
  /** Icon for the item */
  icon?: React.ReactNode;
}

/**
 * Breadcrumb component properties
 */
export interface BreadcrumbProps extends BaseComponent {
  /** Breadcrumb items */
  items: BreadcrumbItem[];
  /** Separator element */
  separator?: React.ReactNode;
  /** Maximum number of items to show */
  maxItems?: number;
  /** Whether to show home icon */
  showHome?: boolean;
  /** Size variant */
  size?: 'sm' | 'md' | 'lg';
}

/**
 * Page header properties (complex header with multiple sections)
 */
export interface PageHeaderProps extends BaseComponent {
  /** Page title */
  title: string;
  /** Page description */
  description?: string;
  /** Breadcrumb navigation */
  breadcrumbs?: BreadcrumbItem[];
  /** Primary actions */
  actions?: React.ReactNode;
  /** Secondary actions (usually in dropdown) */
  secondaryActions?: React.ReactNode;
  /** Tab navigation */
  tabs?: React.ReactNode;
  /** Whether header content spans full width */
  fullWidth?: boolean;
  /** Background image or element */
  background?: React.ReactNode;
  /** Avatar or icon */
  avatar?: React.ReactNode;
  /** Status indicator */
  status?: React.ReactNode;
  /** Metadata (tags, dates, etc.) */
  metadata?: React.ReactNode;
}