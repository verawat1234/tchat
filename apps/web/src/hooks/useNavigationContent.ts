/**
 * Navigation Content Hook
 *
 * Custom hook for managing dynamic navigation content with fallbacks.
 * Provides type-safe access to navigation text with loading states and error handling.
 */

import { useGetContentItemQuery, useGetContentByCategoryQuery } from '../services/contentApi';
import type { ContentItem } from '../types/content';

/**
 * Navigation content IDs following the pattern: navigation.{section}.{element}
 */
export const NAVIGATION_CONTENT_IDS = {
  // Tab navigation
  TABS: {
    CHAT: 'navigation.tabs.chat',
    STORE: 'navigation.tabs.store',
    SOCIAL: 'navigation.tabs.social',
    VIDEO: 'navigation.tabs.video',
    MORE: 'navigation.tabs.more',
    WORK: 'navigation.tabs.work',
  },

  // Header actions
  HEADER: {
    APP_TITLE: 'navigation.header.app_title',
    SEARCH: 'navigation.header.search',
    QR_SCANNER: 'navigation.header.qr_scanner',
    CART: 'navigation.header.cart',
    NOTIFICATIONS: 'navigation.header.notifications',
    SETTINGS: 'navigation.header.settings',
    WALLET: 'navigation.header.wallet',
  },

  // Settings and preferences
  SETTINGS: {
    LANGUAGE: 'navigation.settings.language',
    CURRENCY: 'navigation.settings.currency',
    VERSION: 'navigation.settings.version',
    DATA_USAGE: 'navigation.settings.data_usage',
    HELP_SUPPORT: 'navigation.settings.help_support',
    SIGN_OUT: 'navigation.settings.sign_out',
  },

  // Quick actions
  ACTIONS: {
    WALLET: 'navigation.actions.wallet',
    SETTINGS: 'navigation.actions.settings',
    WORK: 'navigation.actions.work',
    QR_SCANNER: 'navigation.actions.qr_scanner',
  },

  // SEA Features
  FEATURES: {
    MINI_APPS: 'navigation.features.mini_apps',
    DATA_SAVER: 'navigation.features.data_saver',
    QR_PAYMENTS: 'navigation.features.qr_payments',
  },

  // Notifications
  NOTIFICATIONS: {
    TITLE: 'navigation.notifications.title',
    MARK_ALL_READ: 'navigation.notifications.mark_all_read',
    VIEW_ALL: 'navigation.notifications.view_all',
    NO_NOTIFICATIONS: 'navigation.notifications.no_notifications',
  },
} as const;

/**
 * Fallback content for navigation elements
 */
const NAVIGATION_FALLBACKS = {
  [NAVIGATION_CONTENT_IDS.TABS.CHAT]: 'Chat',
  [NAVIGATION_CONTENT_IDS.TABS.STORE]: 'Store',
  [NAVIGATION_CONTENT_IDS.TABS.SOCIAL]: 'Social',
  [NAVIGATION_CONTENT_IDS.TABS.VIDEO]: 'Video',
  [NAVIGATION_CONTENT_IDS.TABS.MORE]: 'More',
  [NAVIGATION_CONTENT_IDS.TABS.WORK]: 'Work',

  [NAVIGATION_CONTENT_IDS.HEADER.APP_TITLE]: 'Telegram SEA',
  [NAVIGATION_CONTENT_IDS.HEADER.SEARCH]: 'Search',
  [NAVIGATION_CONTENT_IDS.HEADER.QR_SCANNER]: 'QR Scanner',
  [NAVIGATION_CONTENT_IDS.HEADER.CART]: 'Cart',
  [NAVIGATION_CONTENT_IDS.HEADER.NOTIFICATIONS]: 'Notifications',
  [NAVIGATION_CONTENT_IDS.HEADER.SETTINGS]: 'Settings',
  [NAVIGATION_CONTENT_IDS.HEADER.WALLET]: 'Wallet',

  [NAVIGATION_CONTENT_IDS.SETTINGS.LANGUAGE]: 'Language',
  [NAVIGATION_CONTENT_IDS.SETTINGS.CURRENCY]: 'Currency',
  [NAVIGATION_CONTENT_IDS.SETTINGS.VERSION]: 'Version',
  [NAVIGATION_CONTENT_IDS.SETTINGS.DATA_USAGE]: 'Data Usage',
  [NAVIGATION_CONTENT_IDS.SETTINGS.HELP_SUPPORT]: 'Help & Support',
  [NAVIGATION_CONTENT_IDS.SETTINGS.SIGN_OUT]: 'Sign Out',

  [NAVIGATION_CONTENT_IDS.ACTIONS.WALLET]: 'Wallet',
  [NAVIGATION_CONTENT_IDS.ACTIONS.SETTINGS]: 'Settings',
  [NAVIGATION_CONTENT_IDS.ACTIONS.WORK]: 'Work',
  [NAVIGATION_CONTENT_IDS.ACTIONS.QR_SCANNER]: 'QR Scanner',

  [NAVIGATION_CONTENT_IDS.FEATURES.MINI_APPS]: 'Mini-Apps',
  [NAVIGATION_CONTENT_IDS.FEATURES.DATA_SAVER]: 'Data Saver',
  [NAVIGATION_CONTENT_IDS.FEATURES.QR_PAYMENTS]: 'QR Payments',

  [NAVIGATION_CONTENT_IDS.NOTIFICATIONS.TITLE]: 'Notifications',
  [NAVIGATION_CONTENT_IDS.NOTIFICATIONS.MARK_ALL_READ]: 'Mark all read',
  [NAVIGATION_CONTENT_IDS.NOTIFICATIONS.VIEW_ALL]: 'View All Notifications',
  [NAVIGATION_CONTENT_IDS.NOTIFICATIONS.NO_NOTIFICATIONS]: 'No notifications',
} as const;

/**
 * Result type for navigation content
 */
export interface NavigationContentResult {
  /** The resolved content text */
  text: string;
  /** Whether content is currently loading */
  isLoading: boolean;
  /** Whether an error occurred */
  isError: boolean;
  /** Whether fallback content is being used */
  isFallback: boolean;
  /** The raw content item (if available) */
  contentItem?: ContentItem;
}

/**
 * Get content text from a ContentItem
 */
const getContentText = (contentItem: ContentItem | undefined): string | undefined => {
  if (!contentItem || contentItem.status !== 'published') {
    return undefined;
  }

  // Handle different content value types
  if (typeof contentItem.value === 'string') {
    return contentItem.value;
  }

  if (typeof contentItem.value === 'object' && contentItem.value !== null) {
    // Handle localized content
    if ('text' in contentItem.value && typeof contentItem.value.text === 'string') {
      return contentItem.value.text;
    }

    // Handle multi-language content (this would be enhanced based on actual content structure)
    if ('en' in contentItem.value && typeof contentItem.value.en === 'string') {
      return contentItem.value.en;
    }
  }

  return undefined;
};

/**
 * Hook for getting individual navigation content item
 */
export const useNavigationContent = (contentId: string): NavigationContentResult => {
  const {
    data: contentItem,
    isLoading,
    isError
  } = useGetContentItemQuery(contentId, {
    // Reduce polling frequency for navigation content
    pollingInterval: 300000, // 5 minutes
    // Skip if contentId is not valid
    skip: !contentId,
  });

  const text = getContentText(contentItem);
  const fallbackText = NAVIGATION_FALLBACKS[contentId as keyof typeof NAVIGATION_FALLBACKS];
  const isFallback = !text && !!fallbackText;

  return {
    text: text || fallbackText || contentId,
    isLoading,
    isError,
    isFallback,
    contentItem,
  };
};

/**
 * Hook for getting multiple navigation content items efficiently
 */
export const useNavigationContentBatch = (contentIds: string[]): Record<string, NavigationContentResult> => {
  // Get navigation category content (more efficient than individual queries)
  const {
    data: categoryResponse,
    isLoading: categoryLoading,
    isError: categoryError
  } = useGetContentByCategoryQuery('navigation', {
    pollingInterval: 300000, // 5 minutes
  });

  // Extract items from paginated response structure
  const categoryContent = categoryResponse?.items || [];

  // Create content map for O(1) lookup
  const contentMap = new Map<string, ContentItem>();
  categoryContent.forEach(item => {
    contentMap.set(item.id, item);
  });

  // Build result for each requested content ID
  const result: Record<string, NavigationContentResult> = {};

  contentIds.forEach(contentId => {
    const contentItem = contentMap.get(contentId);
    const text = getContentText(contentItem);
    const fallbackText = NAVIGATION_FALLBACKS[contentId as keyof typeof NAVIGATION_FALLBACKS];
    const isFallback = !text && !!fallbackText;

    result[contentId] = {
      text: text || fallbackText || contentId,
      isLoading: categoryLoading,
      isError: categoryError,
      isFallback,
      contentItem,
    };
  });

  return result;
};

/**
 * Specialized hooks for common navigation sections
 */

export const useTabNavigationContent = () => {
  return useNavigationContentBatch([
    NAVIGATION_CONTENT_IDS.TABS.CHAT,
    NAVIGATION_CONTENT_IDS.TABS.STORE,
    NAVIGATION_CONTENT_IDS.TABS.SOCIAL,
    NAVIGATION_CONTENT_IDS.TABS.VIDEO,
    NAVIGATION_CONTENT_IDS.TABS.MORE,
    NAVIGATION_CONTENT_IDS.TABS.WORK,
  ]);
};

export const useHeaderNavigationContent = () => {
  return useNavigationContentBatch([
    NAVIGATION_CONTENT_IDS.HEADER.APP_TITLE,
    NAVIGATION_CONTENT_IDS.HEADER.SEARCH,
    NAVIGATION_CONTENT_IDS.HEADER.QR_SCANNER,
    NAVIGATION_CONTENT_IDS.HEADER.CART,
    NAVIGATION_CONTENT_IDS.HEADER.NOTIFICATIONS,
    NAVIGATION_CONTENT_IDS.HEADER.SETTINGS,
    NAVIGATION_CONTENT_IDS.HEADER.WALLET,
  ]);
};

export const useSettingsNavigationContent = () => {
  return useNavigationContentBatch([
    NAVIGATION_CONTENT_IDS.SETTINGS.LANGUAGE,
    NAVIGATION_CONTENT_IDS.SETTINGS.CURRENCY,
    NAVIGATION_CONTENT_IDS.SETTINGS.VERSION,
    NAVIGATION_CONTENT_IDS.SETTINGS.DATA_USAGE,
    NAVIGATION_CONTENT_IDS.SETTINGS.HELP_SUPPORT,
    NAVIGATION_CONTENT_IDS.SETTINGS.SIGN_OUT,
  ]);
};

export const useNotificationsNavigationContent = () => {
  return useNavigationContentBatch([
    NAVIGATION_CONTENT_IDS.NOTIFICATIONS.TITLE,
    NAVIGATION_CONTENT_IDS.NOTIFICATIONS.MARK_ALL_READ,
    NAVIGATION_CONTENT_IDS.NOTIFICATIONS.VIEW_ALL,
    NAVIGATION_CONTENT_IDS.NOTIFICATIONS.NO_NOTIFICATIONS,
  ]);
};

export const useQuickActionsContent = () => {
  return useNavigationContentBatch([
    NAVIGATION_CONTENT_IDS.ACTIONS.WALLET,
    NAVIGATION_CONTENT_IDS.ACTIONS.SETTINGS,
    NAVIGATION_CONTENT_IDS.ACTIONS.WORK,
    NAVIGATION_CONTENT_IDS.ACTIONS.QR_SCANNER,
  ]);
};

export const useSEAFeaturesContent = () => {
  return useNavigationContentBatch([
    NAVIGATION_CONTENT_IDS.FEATURES.MINI_APPS,
    NAVIGATION_CONTENT_IDS.FEATURES.DATA_SAVER,
    NAVIGATION_CONTENT_IDS.FEATURES.QR_PAYMENTS,
  ]);
};