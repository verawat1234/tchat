/**
 * Header Components
 * Page headers, breadcrumbs, and navigation elements
 */

import React from 'react';
import { cn } from '@/utils/cn';
import { useGetContentItemQuery } from '../../../services/contentApi';
import { useSelector } from 'react-redux';
import {
  selectFallbackContentById,
  selectContentSyncMetadata,
  selectContentLoadingState,
  selectSelectedLanguage
} from '../../../features/contentSelectors';
import type {
  HeaderProps,
  BreadcrumbProps,
  PageHeaderProps,
  BreadcrumbItem
} from '../../../../specs/001-agent-frontend-specialist/contracts/header';

/**
 * Custom hook for getting content with fallback support and language awareness
 */
const useContentWithFallback = (contentId: string, defaultText: string = '') => {
  const selectedLanguage = useSelector(selectSelectedLanguage);

  // Create language-specific content ID
  const languageContentId = `${contentId}.${selectedLanguage}`;

  const { data: contentData, isLoading, error } = useGetContentItemQuery(languageContentId);
  const fallbackSelector = selectFallbackContentById(languageContentId);
  const fallbackContent = useSelector(fallbackSelector);

  // Fallback to base content ID if language-specific doesn't exist
  const { data: baseContentData } = useGetContentItemQuery(contentId, {
    skip: !!contentData || !!fallbackContent
  });
  const baseFallbackSelector = selectFallbackContentById(contentId);
  const baseFallbackContent = useSelector(baseFallbackSelector);

  // Return content value with fallback hierarchy
  const getContentText = (content: any) => {
    if (!content) return null;
    if (content.type === 'text') return content.value;
    if (content.type === 'rich_text') return content.value;
    if (content.type === 'translation') {
      return content.values[selectedLanguage] || content.values[content.defaultLocale];
    }
    return null;
  };

  const content = getContentText(contentData?.value) ||
                  getContentText(fallbackContent) ||
                  getContentText(baseContentData?.value) ||
                  getContentText(baseFallbackContent) ||
                  defaultText;

  return {
    content,
    isLoading,
    hasError: !!error,
    hasFallback: !!fallbackContent || !!baseFallbackContent,
    hasRemoteContent: !!contentData || !!baseContentData,
    currentLanguage: selectedLanguage
  };
};

/**
 * Header component for page titles and navigation
 */
export const Header = React.forwardRef<HTMLElement, HeaderProps>(
  ({
    className,
    testId,
    title,
    subtitle,
    level = 1,
    actions,
    breadcrumbs,
    sticky = false,
    border = false,
    size = 'md',
    centered = false,
    background = 'default',
    icon,
    showBack = false,
    onBack,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedby,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const HeadingTag = `h${level}` as keyof JSX.IntrinsicElements;

    // Get dynamic content with fallback
    const backButtonText = useContentWithFallback('header.navigation.back', 'Go back');

    // Get content sync status
    const syncMetadata = useSelector(selectContentSyncMetadata);
    const contentLoadingState = useSelector(selectContentLoadingState);

    // Check if any content is still loading
    const isContentLoading = backButtonText.isLoading || contentLoadingState.isLoading;

    const getSizeClass = () => {
      const sizeMap = {
        sm: 'text-lg',
        md: 'text-xl',
        lg: 'text-2xl',
        xl: 'text-3xl'
      };
      return sizeMap[size];
    };

    const getBackgroundClass = () => {
      const backgroundMap = {
        transparent: 'bg-transparent',
        default: 'bg-white',
        muted: 'bg-gray-50'
      };
      return backgroundMap[background];
    };

    const handleBackClick = () => {
      if (onBack) {
        onBack();
      }
    };

    return (
      <header
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'w-full',

          // Background
          getBackgroundClass(),
          `header-bg-${background}`,

          // Sticky
          sticky && 'sticky top-0 z-10 header-sticky',

          // Border
          border && 'border-b border-gray-200 header-border',

          // Centered
          centered && 'text-center header-centered',

          // Size
          `header-${size}`,

          // Custom className
          className
        )}
        aria-label={ariaLabel}
        aria-describedby={ariaDescribedby}
        aria-expanded={ariaExpanded}
        aria-disabled={ariaDisabled}
        role={role || 'banner'}
        tabIndex={tabIndex}
        {...props}
      >
        {/* Content loading indicator */}
        {isContentLoading && (
          <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-blue-200 via-blue-500 to-blue-200 animate-pulse" />
        )}

        <div className="px-4 py-6">
          {/* Breadcrumbs */}
          {breadcrumbs && breadcrumbs.length > 0 && (
            <Breadcrumb items={breadcrumbs} className="mb-4" />
          )}

          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              {/* Back button */}
              {showBack && (
                <button
                  type="button"
                  onClick={handleBackClick}
                  className="p-2 text-gray-500 hover:text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md"
                  aria-label={backButtonText.content}
                  disabled={backButtonText.isLoading}
                >
                  {backButtonText.isLoading ? (
                    <div className="w-5 h-5 animate-spin rounded-full border-2 border-gray-300 border-t-gray-700" />
                  ) : (
                    <svg
                      className="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M15 19l-7-7 7-7"
                      />
                    </svg>
                  )}
                </button>
              )}

              {/* Icon */}
              {icon && (
                <div className="flex items-center">
                  {icon}
                </div>
              )}

              {/* Title and subtitle */}
              <div className="min-w-0 flex-1">
                <HeadingTag
                  className={cn(
                    'font-bold text-gray-900 truncate',
                    getSizeClass()
                  )}
                >
                  {title}
                </HeadingTag>
                {subtitle && (
                  <p className="mt-1 text-sm text-gray-500 truncate">
                    {subtitle}
                  </p>
                )}
              </div>
            </div>

            {/* Content sync status indicator */}
            <div className="flex items-center space-x-2">
              {/* Show content status */}
              {(syncMetadata.fallbackMode || contentLoadingState.hasError || backButtonText.hasFallback) && (
                <div
                  className={cn(
                    'flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium',
                    syncMetadata.fallbackMode
                      ? 'bg-yellow-100 text-yellow-800'
                      : contentLoadingState.hasError
                      ? 'bg-red-100 text-red-800'
                      : 'bg-blue-100 text-blue-800'
                  )}
                  title={syncMetadata.displayStatus}
                >
                  <div
                    className={cn(
                      'w-2 h-2 rounded-full',
                      syncMetadata.fallbackMode
                        ? 'bg-yellow-400'
                        : contentLoadingState.hasError
                        ? 'bg-red-400'
                        : 'bg-blue-400'
                    )}
                  />
                  <span>
                    {syncMetadata.fallbackMode
                      ? 'Offline'
                      : contentLoadingState.hasError
                      ? 'Error'
                      : 'Cached'}
                  </span>
                </div>
              )}

              {/* Actions */}
              {actions && (
                <>
                  {actions}
                </>
              )}
            </div>
          </div>
        </div>
      </header>
    );
  }
);

Header.displayName = 'Header';

/**
 * Breadcrumb component for navigation hierarchy
 */
export const Breadcrumb = React.forwardRef<HTMLNavElement, BreadcrumbProps>(
  ({
    className,
    testId,
    items,
    separator,
    maxItems,
    showHome = false,
    size = 'md',
    ...props
  }, ref) => {
    // Get dynamic content for breadcrumb navigation
    const breadcrumbAriaLabel = useContentWithFallback('header.breadcrumb.navigation', 'Breadcrumb');
    const getSizeClass = () => {
      const sizeMap = {
        sm: 'text-xs',
        md: 'text-sm',
        lg: 'text-base'
      };
      return sizeMap[size];
    };

    // Truncate items if maxItems is specified
    let displayItems = items;
    let showEllipsis = false;

    if (maxItems && items.length > maxItems) {
      showEllipsis = true;
      // Show first item and last (maxItems - 1) items
      displayItems = [
        items[0], // First item
        ...items.slice(-(maxItems - 1)) // Last (maxItems - 1) items
      ];
    }

    const defaultSeparator = separator || (
      <svg
        className="w-4 h-4 text-gray-400"
        fill="currentColor"
        viewBox="0 0 20 20"
      >
        <path
          fillRule="evenodd"
          d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
          clipRule="evenodd"
        />
      </svg>
    );

    return (
      <nav
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex items-center space-x-2',

          // Size
          getSizeClass(),
          `breadcrumb-${size}`,

          // Custom className
          className
        )}
        aria-label={breadcrumbAriaLabel.content}
        {...props}
      >
        {/* Home icon */}
        {showHome && (
          <>
            <div data-icon="home" className="w-4 h-4 text-gray-500">
              <svg fill="currentColor" viewBox="0 0 20 20">
                <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" />
              </svg>
            </div>
            {items.length > 0 && defaultSeparator}
          </>
        )}

        {displayItems.map((item, index) => {
          const isLast = index === displayItems.length - 1;
          const shouldShowEllipsis = showEllipsis && index === 1;

          return (
            <React.Fragment key={`${item.label}-${index}`}>
              {/* Show ellipsis after first item if truncated */}
              {shouldShowEllipsis && (
                <>
                  <span className="text-gray-500">...</span>
                  {defaultSeparator}
                </>
              )}

              {/* Breadcrumb item */}
              <div className="flex items-center">
                {/* Icon */}
                {item.icon && (
                  <span className="mr-1">
                    {item.icon}
                  </span>
                )}

                {/* Item content */}
                {item.href && !item.disabled ? (
                  <a
                    href={item.href}
                    className={cn(
                      'hover:text-gray-700 transition-colors duration-200',
                      isLast ? 'text-gray-900 font-medium' : 'text-gray-500'
                    )}
                    aria-current={isLast ? 'page' : undefined}
                  >
                    {item.label}
                  </a>
                ) : item.onClick && !item.disabled ? (
                  <button
                    type="button"
                    onClick={item.onClick}
                    className={cn(
                      'hover:text-gray-700 transition-colors duration-200 text-left',
                      isLast ? 'text-gray-900 font-medium' : 'text-gray-500'
                    )}
                    aria-current={isLast ? 'page' : undefined}
                  >
                    {item.label}
                  </button>
                ) : (
                  <span
                    className={cn(
                      isLast ? 'text-gray-900 font-medium' : 'text-gray-500',
                      item.disabled && 'opacity-50'
                    )}
                    aria-current={isLast ? 'page' : undefined}
                    aria-disabled={item.disabled}
                  >
                    {item.label}
                  </span>
                )}
              </div>

              {/* Separator */}
              {!isLast && (
                <span className="text-gray-400">
                  {defaultSeparator}
                </span>
              )}
            </React.Fragment>
          );
        })}
      </nav>
    );
  }
);

Breadcrumb.displayName = 'Breadcrumb';

/**
 * Page header component for complex page layouts
 */
export const PageHeader = React.forwardRef<HTMLElement, PageHeaderProps>(
  ({
    className,
    testId,
    title,
    description,
    breadcrumbs,
    actions,
    secondaryActions,
    tabs,
    fullWidth = false,
    background,
    avatar,
    status,
    metadata,
    ...props
  }, ref) => {
    // Content loading state indicator
    const [isContentReady, setIsContentReady] = React.useState(false);

    React.useEffect(() => {
      // Simulate content readiness check
      const timer = setTimeout(() => setIsContentReady(true), 100);
      return () => clearTimeout(timer);
    }, []);
    return (
      <header
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'w-full bg-white border-b border-gray-200',

          // Full width
          fullWidth && 'page-header-full-width',

          // Custom className
          className
        )}
        {...props}
      >
        {/* Background */}
        {background && (
          <div className="absolute inset-0 -z-10">
            {background}
          </div>
        )}

        <div className={cn(
          'px-4 py-6',
          !fullWidth && 'max-w-7xl mx-auto'
        )}>
          {/* Breadcrumbs */}
          {breadcrumbs && breadcrumbs.length > 0 && (
            <Breadcrumb items={breadcrumbs} className="mb-4" />
          )}

          {/* Main header content */}
          <div className="flex items-start justify-between">
            <div className="flex items-start space-x-4 min-w-0 flex-1">
              {/* Avatar */}
              {avatar && (
                <div className="flex-shrink-0">
                  {avatar}
                </div>
              )}

              <div className="min-w-0 flex-1">
                <div className="flex items-center space-x-3">
                  <h1 className="text-2xl font-bold text-gray-900 truncate">
                    {title}
                  </h1>

                  {/* Status */}
                  {status && (
                    <div className="flex-shrink-0">
                      {status}
                    </div>
                  )}
                </div>

                {/* Description */}
                {description && (
                  <p className="mt-2 text-sm text-gray-600">
                    {description}
                  </p>
                )}

                {/* Metadata */}
                {metadata && (
                  <div className="mt-3">
                    {metadata}
                  </div>
                )}
              </div>
            </div>

            {/* Actions */}
            <div className="flex items-center space-x-3 ml-4">
              {/* Secondary actions */}
              {secondaryActions && (
                <div className="flex items-center space-x-2">
                  {secondaryActions}
                </div>
              )}

              {/* Primary actions */}
              {actions && (
                <div className="flex items-center space-x-2">
                  {actions}
                </div>
              )}
            </div>
          </div>

          {/* Tabs */}
          {tabs && (
            <div className="mt-6">
              {tabs}
            </div>
          )}
        </div>
      </header>
    );
  }
);

PageHeader.displayName = 'PageHeader';

export type {
  HeaderProps,
  BreadcrumbProps,
  PageHeaderProps,
  BreadcrumbItem
};