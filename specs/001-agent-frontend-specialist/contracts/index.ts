/**
 * Common UI Component Contracts
 * Type definitions and interfaces for all common UI components
 */

// Base types and interfaces
export * from './base';

// Component contracts
export * from './pagination';
export * from './tabs';
export * from './card';
export * from './chat-message';
export * from './badge';
export * from './layout';
export * from './header';
export * from './sidebar';

/**
 * Component registry for dynamic loading
 */
export const COMPONENT_REGISTRY = {
  pagination: 'Pagination',
  tabs: 'Tabs',
  card: 'Card',
  chatMessage: 'ChatMessage',
  badge: 'Badge',
  layout: 'Layout',
  header: 'Header',
  sidebar: 'Sidebar',
} as const;

/**
 * Component categories for organization
 */
export const COMPONENT_CATEGORIES = {
  navigation: ['pagination', 'tabs', 'sidebar'],
  layout: ['layout', 'card'],
  data: ['chatMessage'],
  feedback: ['badge'],
  typography: ['header'],
} as const;

/**
 * Export types for component registry
 */
export type ComponentName = keyof typeof COMPONENT_REGISTRY;
export type ComponentCategory = keyof typeof COMPONENT_CATEGORIES;