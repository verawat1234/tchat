"use client";

import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { AlertCircle, Wifi, WifiOff, RefreshCw } from "lucide-react";
import { cn } from "./ui/utils";
import { Skeleton } from "./ui/skeleton";
import { Progress } from "./ui/progress";
import { useSelector } from "react-redux";
import { selectOperationLoading, selectAnyLoading } from "../features/loadingSlice";
import type { RootState } from "../store";

// ============================================================================
// Types and Interfaces
// ============================================================================

export interface ContentLoaderProps extends VariantProps<typeof contentLoaderVariants> {
  className?: string;
  children?: React.ReactNode;

  // Loading state configuration
  isLoading?: boolean;
  loadingType?: 'skeleton' | 'spinner' | 'progress' | 'minimal';

  // Progress indicator
  progress?: number;
  progressMessage?: string;

  // Error handling
  error?: Error | string | null;
  onRetry?: () => void;

  // Fallback/offline states
  isFallback?: boolean;
  isOffline?: boolean;
  fallbackMessage?: string;

  // Content type-specific skeletons
  contentType?: 'text' | 'image' | 'card' | 'list' | 'table' | 'chat' | 'media' | 'custom';
  skeletonCount?: number;

  // RTK Query integration
  rtkQueryKey?: string;

  // Accessibility
  loadingLabel?: string;
  announceChanges?: boolean;

  // Performance optimizations
  reduceMotion?: boolean;
  deferRender?: boolean;

  // Custom skeleton renderer
  renderSkeleton?: () => React.ReactNode;

  // Retry configuration
  maxRetries?: number;
  retryCount?: number;
}

// ============================================================================
// Component Variants
// ============================================================================

const contentLoaderVariants = cva(
  "relative transition-all duration-200 ease-in-out",
  {
    variants: {
      size: {
        sm: "min-h-16",
        md: "min-h-24",
        lg: "min-h-32",
        xl: "min-h-48",
        full: "min-h-screen",
        auto: "min-h-0",
      },
      spacing: {
        tight: "space-y-2",
        normal: "space-y-4",
        loose: "space-y-6",
      },
      animation: {
        pulse: "animate-pulse",
        wave: "animate-[wave_2s_ease-in-out_infinite]",
        fade: "animate-[fade_1s_ease-in-out_infinite_alternate]",
        none: "",
      },
    },
    defaultVariants: {
      size: "auto",
      spacing: "normal",
      animation: "pulse",
    },
  }
);

// ============================================================================
// Skeleton Components for Different Content Types
// ============================================================================

const TextSkeleton: React.FC<{ count?: number; className?: string }> = ({
  count = 3,
  className
}) => (
  <div className={cn("space-y-2", className)}>
    {Array.from({ length: count }).map((_, i) => (
      <Skeleton
        key={i}
        className={cn(
          "h-4",
          i === 0 ? "w-3/4" : i === count - 1 ? "w-1/2" : "w-full"
        )}
      />
    ))}
  </div>
);

const ImageSkeleton: React.FC<{ className?: string; aspectRatio?: string }> = ({
  className,
  aspectRatio = "aspect-video"
}) => (
  <Skeleton className={cn("w-full rounded-lg", aspectRatio, className)} />
);

const CardSkeleton: React.FC<{ className?: string }> = ({ className }) => (
  <div className={cn("space-y-3", className)}>
    <Skeleton className="h-48 w-full rounded-lg" />
    <div className="space-y-2">
      <Skeleton className="h-4 w-3/4" />
      <Skeleton className="h-4 w-1/2" />
    </div>
  </div>
);

const ListSkeleton: React.FC<{ count?: number; className?: string }> = ({
  count = 5,
  className
}) => (
  <div className={cn("space-y-3", className)}>
    {Array.from({ length: count }).map((_, i) => (
      <div key={i} className="flex items-center space-x-3">
        <Skeleton className="h-10 w-10 rounded-full" />
        <div className="flex-1 space-y-2">
          <Skeleton className="h-4 w-3/4" />
          <Skeleton className="h-3 w-1/2" />
        </div>
      </div>
    ))}
  </div>
);

const TableSkeleton: React.FC<{ rows?: number; cols?: number; className?: string }> = ({
  rows = 5,
  cols = 4,
  className
}) => (
  <div className={cn("space-y-3", className)}>
    {/* Header */}
    <div className="flex space-x-3">
      {Array.from({ length: cols }).map((_, i) => (
        <Skeleton key={i} className="h-6 flex-1" />
      ))}
    </div>
    {/* Rows */}
    {Array.from({ length: rows }).map((_, i) => (
      <div key={i} className="flex space-x-3">
        {Array.from({ length: cols }).map((_, j) => (
          <Skeleton key={j} className="h-4 flex-1" />
        ))}
      </div>
    ))}
  </div>
);

const ChatSkeleton: React.FC<{ count?: number; className?: string }> = ({
  count = 4,
  className
}) => (
  <div className={cn("space-y-4", className)}>
    {Array.from({ length: count }).map((_, i) => (
      <div key={i} className={cn("flex", i % 2 === 0 ? "justify-start" : "justify-end")}>
        <div className={cn("flex max-w-[70%] space-x-2", i % 2 === 0 ? "flex-row" : "flex-row-reverse space-x-reverse")}>
          <Skeleton className="h-8 w-8 rounded-full flex-shrink-0" />
          <div className="space-y-1">
            <Skeleton className="h-4 w-20" />
            <Skeleton className={cn("h-12 rounded-lg", i % 2 === 0 ? "w-32" : "w-28")} />
          </div>
        </div>
      </div>
    ))}
  </div>
);

const MediaSkeleton: React.FC<{ className?: string }> = ({ className }) => (
  <div className={cn("space-y-3", className)}>
    <Skeleton className="h-64 w-full rounded-lg" />
    <div className="flex items-center space-x-3">
      <Skeleton className="h-12 w-12 rounded-full" />
      <div className="flex-1 space-y-1">
        <Skeleton className="h-4 w-1/2" />
        <Skeleton className="h-3 w-1/4" />
      </div>
      <Skeleton className="h-8 w-20" />
    </div>
  </div>
);

// ============================================================================
// Loading Indicators
// ============================================================================

const SpinnerLoader: React.FC<{ size?: 'sm' | 'md' | 'lg'; className?: string }> = ({
  size = 'md',
  className
}) => {
  const sizeClasses = {
    sm: "h-4 w-4",
    md: "h-6 w-6",
    lg: "h-8 w-8"
  };

  return (
    <div className={cn("flex items-center justify-center py-8", className)}>
      <RefreshCw className={cn("animate-spin text-muted-foreground", sizeClasses[size])} />
    </div>
  );
};

const ProgressLoader: React.FC<{
  progress?: number;
  message?: string;
  className?: string
}> = ({ progress = 0, message, className }) => (
  <div className={cn("space-y-3 py-6", className)}>
    <Progress value={progress} className="w-full" />
    <div className="flex items-center justify-between text-sm text-muted-foreground">
      <span>{message || "Loading..."}</span>
      <span>{Math.round(progress)}%</span>
    </div>
  </div>
);

// ============================================================================
// State Indicators
// ============================================================================

const ErrorState: React.FC<{
  error: Error | string;
  onRetry?: () => void;
  maxRetries?: number;
  retryCount?: number;
  className?: string;
}> = ({ error, onRetry, maxRetries = 3, retryCount = 0, className }) => {
  const errorMessage = typeof error === 'string' ? error : error.message;
  const canRetry = onRetry && retryCount < maxRetries;

  return (
    <div className={cn("flex flex-col items-center justify-center py-8 space-y-4", className)}>
      <div className="flex items-center space-x-2 text-destructive">
        <AlertCircle className="h-5 w-5" />
        <span className="font-medium">Failed to load content</span>
      </div>

      <p className="text-sm text-muted-foreground text-center max-w-md">
        {errorMessage}
      </p>

      {canRetry && (
        <button
          onClick={onRetry}
          className="inline-flex items-center space-x-2 px-3 py-2 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
          aria-label={`Retry loading (${retryCount + 1}/${maxRetries} attempts)`}
        >
          <RefreshCw className="h-4 w-4" />
          <span>Try again ({maxRetries - retryCount} attempts left)</span>
        </button>
      )}
    </div>
  );
};

const FallbackState: React.FC<{
  isOffline?: boolean;
  message?: string;
  className?: string;
}> = ({ isOffline, message, className }) => (
  <div className={cn("flex items-center justify-center py-4", className)}>
    <div className="flex items-center space-x-2 px-3 py-2 bg-muted rounded-md">
      {isOffline ? (
        <WifiOff className="h-4 w-4 text-muted-foreground" />
      ) : (
        <Wifi className="h-4 w-4 text-muted-foreground" />
      )}
      <span className="text-sm text-muted-foreground">
        {message || (isOffline ? "Showing offline content" : "Using cached content")}
      </span>
    </div>
  </div>
);

// ============================================================================
// Custom Hooks
// ============================================================================

const useContentLoader = (rtkQueryKey?: string) => {
  const rtkOperation = useSelector((state: RootState) =>
    rtkQueryKey ? selectOperationLoading(rtkQueryKey)(state) : null
  );
  const anyLoading = useSelector(selectAnyLoading);

  return {
    isLoading: rtkOperation?.isLoading || false,
    progress: rtkOperation?.progress,
    message: rtkOperation?.message,
    hasAnyLoading: anyLoading,
  };
};

const useAccessibilityAnnouncer = (announceChanges: boolean = true) => {
  const [announcement, setAnnouncement] = React.useState<string>("");

  const announce = React.useCallback((message: string) => {
    if (announceChanges) {
      setAnnouncement(message);
      // Clear after announcement to allow re-announcing the same message
      setTimeout(() => setAnnouncement(""), 1000);
    }
  }, [announceChanges]);

  return { announcement, announce };
};

// ============================================================================
// Main ContentLoader Component
// ============================================================================

export const ContentLoader: React.FC<ContentLoaderProps> = ({
  className,
  children,
  isLoading: externalLoading,
  loadingType = 'skeleton',
  progress,
  progressMessage,
  error,
  onRetry,
  isFallback = false,
  isOffline = false,
  fallbackMessage,
  contentType = 'text',
  skeletonCount = 3,
  rtkQueryKey,
  loadingLabel = "Loading content",
  announceChanges = true,
  reduceMotion = false,
  deferRender = false,
  renderSkeleton,
  maxRetries = 3,
  retryCount = 0,
  size,
  spacing,
  animation,
  ...props
}) => {
  // RTK Query integration
  const rtkState = useContentLoader(rtkQueryKey);
  const isLoading = externalLoading ?? rtkState.isLoading;
  const loadingProgress = progress ?? rtkState.progress;
  const loadingMessage = progressMessage ?? rtkState.message;

  // Accessibility
  const { announcement, announce } = useAccessibilityAnnouncer(announceChanges);

  // Announce loading state changes
  React.useEffect(() => {
    if (isLoading) {
      announce(loadingMessage || loadingLabel);
    } else if (error) {
      announce(`Error loading content: ${typeof error === 'string' ? error : error.message}`);
    } else if (isFallback) {
      announce(fallbackMessage || "Showing cached content");
    }
  }, [isLoading, error, isFallback, loadingMessage, loadingLabel, fallbackMessage, announce]);

  // Deferred rendering optimization
  const [shouldRender, setShouldRender] = React.useState(!deferRender);
  React.useEffect(() => {
    if (deferRender && isLoading) {
      const timer = setTimeout(() => setShouldRender(true), 100);
      return () => clearTimeout(timer);
    }
  }, [deferRender, isLoading]);

  // Respect user's motion preferences
  const effectiveAnimation = reduceMotion ? 'none' : animation;

  // Render skeleton based on content type
  const renderContentSkeleton = () => {
    if (renderSkeleton) {
      return renderSkeleton();
    }

    const skeletonProps = { count: skeletonCount, className: "w-full" };

    switch (contentType) {
      case 'text':
        return <TextSkeleton {...skeletonProps} />;
      case 'image':
        return <ImageSkeleton className="w-full" />;
      case 'card':
        return <CardSkeleton className="w-full" />;
      case 'list':
        return <ListSkeleton {...skeletonProps} />;
      case 'table':
        return <TableSkeleton className="w-full" />;
      case 'chat':
        return <ChatSkeleton {...skeletonProps} />;
      case 'media':
        return <MediaSkeleton className="w-full" />;
      default:
        return <TextSkeleton {...skeletonProps} />;
    }
  };

  // Render loading indicator based on type
  const renderLoadingIndicator = () => {
    switch (loadingType) {
      case 'spinner':
        return <SpinnerLoader className="w-full" />;
      case 'progress':
        return (
          <ProgressLoader
            progress={loadingProgress}
            message={loadingMessage}
            className="w-full"
          />
        );
      case 'minimal':
        return (
          <div className="flex items-center justify-center py-4">
            <span className="text-sm text-muted-foreground">{loadingMessage || loadingLabel}</span>
          </div>
        );
      case 'skeleton':
      default:
        return renderContentSkeleton();
    }
  };

  if (!shouldRender) {
    return null;
  }

  // Error state
  if (error) {
    return (
      <div
        className={cn(contentLoaderVariants({ size, spacing }), className)}
        role="alert"
        aria-live="polite"
        {...props}
      >
        <ErrorState
          error={error}
          onRetry={onRetry}
          maxRetries={maxRetries}
          retryCount={retryCount}
        />

        {/* Screen reader announcement */}
        <div className="sr-only" aria-live="polite" aria-atomic="true">
          {announcement}
        </div>
      </div>
    );
  }

  // Loading state
  if (isLoading) {
    return (
      <div
        className={cn(
          contentLoaderVariants({ size, spacing, animation: effectiveAnimation }),
          className
        )}
        role="status"
        aria-live="polite"
        aria-label={loadingMessage || loadingLabel}
        {...props}
      >
        {renderLoadingIndicator()}

        {/* Screen reader announcement */}
        <div className="sr-only" aria-live="polite" aria-atomic="true">
          {announcement}
        </div>
      </div>
    );
  }

  // Content loaded successfully
  return (
    <div className={cn(contentLoaderVariants({ size, spacing }), className)} {...props}>
      {/* Fallback/offline indicator */}
      {(isFallback || isOffline) && (
        <FallbackState
          isOffline={isOffline}
          message={fallbackMessage}
        />
      )}

      {/* Main content */}
      {children}

      {/* Screen reader announcement */}
      <div className="sr-only" aria-live="polite" aria-atomic="true">
        {announcement}
      </div>
    </div>
  );
};

// ============================================================================
// Compound Components
// ============================================================================

ContentLoader.Text = TextSkeleton;
ContentLoader.Image = ImageSkeleton;
ContentLoader.Card = CardSkeleton;
ContentLoader.List = ListSkeleton;
ContentLoader.Table = TableSkeleton;
ContentLoader.Chat = ChatSkeleton;
ContentLoader.Media = MediaSkeleton;
ContentLoader.Spinner = SpinnerLoader;
ContentLoader.Progress = ProgressLoader;
ContentLoader.Error = ErrorState;
ContentLoader.Fallback = FallbackState;

// ============================================================================
// Exports
// ============================================================================

export default ContentLoader;

// Export individual components for direct usage
export {
  TextSkeleton,
  ImageSkeleton,
  CardSkeleton,
  ListSkeleton,
  TableSkeleton,
  ChatSkeleton,
  MediaSkeleton,
  SpinnerLoader,
  ProgressLoader,
  ErrorState,
  FallbackState,
  useContentLoader,
  useAccessibilityAnnouncer,
};

// Export types
export type {
  ContentLoaderProps,
};