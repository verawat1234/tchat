import { useGetContentItemQuery } from '../services/content';
import type { ContentItem } from '../types/content';
import { useSelector } from 'react-redux';
import type { RootState } from '../store';

/**
 * Custom hook to get current locale from various sources
 *
 * Attempts to determine the user's preferred locale from:
 * 1. User preferences from Redux store (authenticated user's locale/contentLanguages)
 * 2. Browser language
 * 3. Navigator language
 * 4. Default fallback
 */
function useCurrentLocale(): string {
  const currentUser = useSelector((state: RootState) => state.auth.user);

  // 1. Try authenticated user's locale preference
  if (currentUser?.locale) {
    // Extract primary language code (e.g., 'en' from 'en-US' or 'th-TH')
    return currentUser.locale.split('-')[0].toLowerCase();
  }

  // 2. Try user's content language preferences (if available)
  if (currentUser?.preferences?.contentLanguages?.length > 0) {
    const primaryContentLang = currentUser.preferences.contentLanguages[0];
    return primaryContentLang.split('-')[0].toLowerCase();
  }

  // 3. Try browser language
  if (typeof window !== 'undefined') {
    // Get primary language from browser
    const browserLang = window.navigator.language;
    if (browserLang) {
      // Extract primary language code (e.g., 'en' from 'en-US')
      return browserLang.split('-')[0].toLowerCase();
    }

    // Try navigator languages array
    const languages = window.navigator.languages;
    if (languages && languages.length > 0) {
      return languages[0].split('-')[0].toLowerCase();
    }
  }

  // 4. Default fallback
  return 'en';
}

/**
 * Hook for fetching dynamic content text with fallback support
 *
 * This hook provides a convenient way to fetch text content from the content API
 * with automatic fallback to provided default text in case of loading states or errors.
 * Ideal for error messages, UI text, and other dynamic content.
 */
export function useContentText(
  contentId: string,
  fallbackText: string,
  options?: {
    /** Skip the API call if true */
    skip?: boolean;
    /** Polling interval in milliseconds */
    pollingInterval?: number;
  }
): {
  /** The resolved text content (either from API or fallback) */
  text: string;
  /** Whether the content is currently loading */
  isLoading: boolean;
  /** Whether there was an error fetching content */
  isError: boolean;
  /** Whether the returned text is from fallback (not API) */
  isFallback: boolean;
} {
  // Get current locale using the hook
  const currentLocale = useCurrentLocale();

  const {
    data: contentItem,
    isLoading,
    isError,
  } = useGetContentItemQuery(contentId, {
    skip: options?.skip,
    pollingInterval: options?.pollingInterval,
  });

  // Extract text from content item
  const getTextFromContentItem = (item: ContentItem, locale: string): string => {
    if (!item?.data) {
      return fallbackText;
    }

    // Handle different content types
    switch (item.type) {
      case 'text':
      case 'rich_text':
        return typeof item.data === 'string' ? item.data : fallbackText;

      case 'translation':
        // For translation type, try to get the current locale or default
        if (typeof item.data === 'object' && item.data !== null) {
          const translations = item.data as Record<string, string>;
          // Use the passed locale parameter (from user preferences)
          return translations[locale] ||
                 translations['en'] ||
                 translations['en-US'] ||
                 Object.values(translations)[0] ||
                 fallbackText;
        }
        return fallbackText;

      case 'config':
        // For config type, try to convert to string
        if (typeof item.data === 'string') {
          return item.data;
        } else if (typeof item.data === 'number' || typeof item.data === 'boolean') {
          return String(item.data);
        }
        return fallbackText;

      default:
        return fallbackText;
    }
  };

  const text = contentItem ? getTextFromContentItem(contentItem, currentLocale) : fallbackText;
  const isFallback = !contentItem || isError || text === fallbackText;

  return {
    text,
    isLoading,
    isError,
    isFallback,
  };
}

/**
 * Hook for fetching multiple content items efficiently
 *
 * Useful when you need to fetch several related content items at once,
 * such as all error messages for a component.
 */
export function useContentTexts(
  contentIds: Array<{ id: string; fallback: string }>,
  options?: {
    skip?: boolean;
    pollingInterval?: number;
  }
): Record<string, {
  text: string;
  isLoading: boolean;
  isError: boolean;
  isFallback: boolean;
}> {
  const results: Record<string, ReturnType<typeof useContentText>> = {};

  contentIds.forEach(({ id, fallback }) => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    results[id] = useContentText(id, fallback, options);
  });

  return results;
}