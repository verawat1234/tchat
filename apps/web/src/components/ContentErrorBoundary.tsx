import React, { Component, ErrorInfo, ReactNode, useCallback, useEffect, useState, useContext, useMemo } from 'react';
import { AlertCircle, RefreshCw, Home, Clock, Wifi, Database } from 'lucide-react';
import { Button } from './ui/button';
import { useContentTexts } from '../hooks/useContentText';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  maxRetries?: number;
  retryDelay?: number;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
  reportError?: (error: ContentError) => void;
}

interface State {
  hasError: boolean;
  error?: Error;
  errorInfo?: ErrorInfo;
  retryCount: number;
  errorType: ContentErrorType;
  lastErrorTime: number;
}

/**
 * Content-specific error types for better error handling
 */
export enum ContentErrorType {
  NETWORK = 'NETWORK',
  PARSING = 'PARSING',
  LOADING = 'LOADING',
  VALIDATION = 'VALIDATION',
  RTK_QUERY = 'RTK_QUERY',
  CONTENT_MISSING = 'CONTENT_MISSING',
  PERMISSION = 'PERMISSION',
  UNKNOWN = 'UNKNOWN'
}

/**
 * Enhanced error interface for content failures
 */
export interface ContentError {
  type: ContentErrorType;
  message: string;
  originalError: Error;
  timestamp: number;
  retryCount: number;
  componentStack?: string;
  isRecoverable: boolean;
}

/**
 * Error classification utilities
 */
const classifyContentError = (error: Error): { type: ContentErrorType; isRecoverable: boolean } => {
  const errorMessage = error.message.toLowerCase();
  const errorName = error.name.toLowerCase();

  // Network-related errors
  if (
    errorMessage.includes('network') ||
    errorMessage.includes('fetch') ||
    errorMessage.includes('connection') ||
    errorName.includes('networkerror') ||
    errorMessage.includes('offline')
  ) {
    return { type: ContentErrorType.NETWORK, isRecoverable: true };
  }

  // RTK Query specific errors
  if (
    errorMessage.includes('rtk') ||
    errorMessage.includes('query') ||
    errorMessage.includes('endpoint') ||
    errorName.includes('rtkqueryerror')
  ) {
    return { type: ContentErrorType.RTK_QUERY, isRecoverable: true };
  }

  // Content parsing errors
  if (
    errorMessage.includes('parse') ||
    errorMessage.includes('json') ||
    errorMessage.includes('syntax') ||
    errorMessage.includes('malformed')
  ) {
    return { type: ContentErrorType.PARSING, isRecoverable: false };
  }

  // Content loading errors
  if (
    errorMessage.includes('load') ||
    errorMessage.includes('timeout') ||
    errorMessage.includes('abort') ||
    errorMessage.includes('cancelled')
  ) {
    return { type: ContentErrorType.LOADING, isRecoverable: true };
  }

  // Content validation errors
  if (
    errorMessage.includes('validation') ||
    errorMessage.includes('invalid') ||
    errorMessage.includes('schema') ||
    errorMessage.includes('format')
  ) {
    return { type: ContentErrorType.VALIDATION, isRecoverable: false };
  }

  // Content missing errors
  if (
    errorMessage.includes('not found') ||
    errorMessage.includes('missing') ||
    errorMessage.includes('404') ||
    errorMessage.includes('unavailable')
  ) {
    return { type: ContentErrorType.CONTENT_MISSING, isRecoverable: true };
  }

  // Permission errors
  if (
    errorMessage.includes('permission') ||
    errorMessage.includes('unauthorized') ||
    errorMessage.includes('forbidden') ||
    errorMessage.includes('401') ||
    errorMessage.includes('403')
  ) {
    return { type: ContentErrorType.PERMISSION, isRecoverable: false };
  }

  return { type: ContentErrorType.UNKNOWN, isRecoverable: true };
};

/**
 * Error reporting utility
 */
const reportContentError = (error: ContentError) => {
  // Log to console in development
  if (process.env.NODE_ENV === 'development') {
    console.group('🚨 Content Error Boundary - Error Report');
    console.error('Type:', error.type);
    console.error('Message:', error.message);
    console.error('Recoverable:', error.isRecoverable);
    console.error('Retry Count:', error.retryCount);
    console.error('Original Error:', error.originalError);
    if (error.componentStack) {
      console.error('Component Stack:', error.componentStack);
    }
    console.groupEnd();
  }

  // In production, send to error reporting service
  if (process.env.NODE_ENV === 'production') {
    try {
      // Example: Send to error tracking service
      // errorTrackingService.captureException(error);

      // Store locally for debugging
      const errorLog = {
        ...error,
        userAgent: navigator.userAgent,
        url: window.location.href,
        timestamp: new Date().toISOString()
      };

      localStorage.setItem(
        `content-error-${error.timestamp}`,
        JSON.stringify(errorLog)
      );
    } catch (reportingError) {
      console.warn('Failed to report content error:', reportingError);
    }
  }
};

/**
 * Enhanced Error Display Component with Recovery Features
 */
interface ErrorDisplayProps {
  error?: Error;
  errorInfo?: ErrorInfo;
  errorType: ContentErrorType;
  retryCount: number;
  maxRetries: number;
  onReset: () => void;
  onReload: () => void;
  onRetry?: () => void;
  isRetrying?: boolean;
}

function ErrorDisplay({
  error,
  errorInfo,
  errorType,
  retryCount,
  maxRetries,
  onReset,
  onReload,
  onRetry,
  isRetrying = false
}: ErrorDisplayProps) {
  const [isNetworkOnline, setIsNetworkOnline] = useState(navigator.onLine);

  // Monitor network status for network-related errors
  useEffect(() => {
    const handleOnline = () => setIsNetworkOnline(true);
    const handleOffline = () => setIsNetworkOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  // Auto-retry for network errors when connection is restored
  useEffect(() => {
    if (
      errorType === ContentErrorType.NETWORK &&
      isNetworkOnline &&
      retryCount < maxRetries &&
      onRetry &&
      !isRetrying
    ) {
      const retryTimer = setTimeout(() => {
        onRetry();
      }, 1000);

      return () => clearTimeout(retryTimer);
    }
  }, [isNetworkOnline, errorType, retryCount, maxRetries, onRetry, isRetrying]);

  const content = useContentTexts([
    { id: 'error.boundary.title', fallback: getErrorTitle(errorType) },
    { id: 'error.boundary.description', fallback: getErrorDescription(errorType, isNetworkOnline) },
    { id: 'error.boundary.details_summary', fallback: 'Error Details (Development)' },
    { id: 'error.boundary.try_again', fallback: 'Try Again' },
    { id: 'error.boundary.reload_app', fallback: 'Reload App' },
    { id: 'error.boundary.retry_automatic', fallback: 'Auto-retry in progress...' },
    { id: 'error.boundary.retry_count', fallback: `Retry ${retryCount}/${maxRetries}` },
  ]);

  const getErrorIcon = () => {
    switch (errorType) {
      case ContentErrorType.NETWORK:
        return <Wifi className={`w-12 h-12 mx-auto ${isNetworkOnline ? 'text-warning' : 'text-destructive'}`} />;
      case ContentErrorType.LOADING:
        return <Clock className="w-12 h-12 text-warning mx-auto" />;
      case ContentErrorType.RTK_QUERY:
      case ContentErrorType.CONTENT_MISSING:
        return <Database className="w-12 h-12 text-warning mx-auto" />;
      default:
        return <AlertCircle className="w-12 h-12 text-destructive mx-auto" />;
    }
  };

  const canRetry = retryCount < maxRetries && onRetry && (
    errorType === ContentErrorType.NETWORK ||
    errorType === ContentErrorType.LOADING ||
    errorType === ContentErrorType.RTK_QUERY ||
    errorType === ContentErrorType.CONTENT_MISSING
  );

  return (
    <div className="min-h-[200px] flex items-center justify-center p-4">
      <div className="text-center max-w-md">
        <div className="mb-4">
          {getErrorIcon()}
        </div>

        <h3 className="text-lg font-semibold mb-2">
          {content['error.boundary.title'].text}
        </h3>

        <p className="text-sm text-muted-foreground mb-4">
          {content['error.boundary.description'].text}
        </p>

        {/* Network status indicator */}
        {errorType === ContentErrorType.NETWORK && (
          <div className={`flex items-center justify-center gap-2 mb-4 text-sm ${
            isNetworkOnline ? 'text-success' : 'text-destructive'
          }`}>
            <Wifi className="w-4 h-4" />
            {isNetworkOnline ? 'Network: Online' : 'Network: Offline'}
          </div>
        )}

        {/* Retry information */}
        {retryCount > 0 && (
          <div className="text-sm text-muted-foreground mb-4">
            {content['error.boundary.retry_count'].text}
          </div>
        )}

        {/* Auto-retry indicator */}
        {isRetrying && (
          <div className="flex items-center justify-center gap-2 mb-4 text-sm text-muted-foreground">
            <RefreshCw className="w-4 h-4 animate-spin" />
            {content['error.boundary.retry_automatic'].text}
          </div>
        )}

        {/* Development error details */}
        {process.env.NODE_ENV === 'development' && error && (
          <details className="mb-4 text-left">
            <summary className="cursor-pointer text-sm font-medium mb-2">
              {content['error.boundary.details_summary'].text}
            </summary>
            <div className="bg-muted p-3 rounded text-xs overflow-auto max-h-32">
              <div className="mb-2">
                <strong>Type:</strong> {errorType}
              </div>
              <pre>{error.toString()}</pre>
              {errorInfo && (
                <pre className="mt-2">{errorInfo.componentStack}</pre>
              )}
            </div>
          </details>
        )}

        <div className="flex gap-2 justify-center flex-wrap">
          {canRetry && !isRetrying && (
            <Button
              onClick={onRetry}
              variant="outline"
              size="sm"
              className="flex items-center gap-2"
              disabled={isRetrying}
            >
              <RefreshCw className="w-4 h-4" />
              {content['error.boundary.try_again'].text}
            </Button>
          )}

          <Button
            onClick={onReset}
            variant="outline"
            size="sm"
            className="flex items-center gap-2"
            disabled={isRetrying}
          >
            <AlertCircle className="w-4 h-4" />
            Reset Component
          </Button>

          <Button
            onClick={onReload}
            variant="default"
            size="sm"
            className="flex items-center gap-2"
            disabled={isRetrying}
          >
            <Home className="w-4 h-4" />
            {content['error.boundary.reload_app'].text}
          </Button>
        </div>
      </div>
    </div>
  );
}

/**
 * Get appropriate error title based on error type
 */
function getErrorTitle(errorType: ContentErrorType): string {
  switch (errorType) {
    case ContentErrorType.NETWORK:
      return 'Network Connection Error';
    case ContentErrorType.LOADING:
      return 'Content Loading Error';
    case ContentErrorType.PARSING:
      return 'Content Format Error';
    case ContentErrorType.RTK_QUERY:
      return 'Data Service Error';
    case ContentErrorType.CONTENT_MISSING:
      return 'Content Not Available';
    case ContentErrorType.PERMISSION:
      return 'Access Permission Error';
    case ContentErrorType.VALIDATION:
      return 'Content Validation Error';
    default:
      return 'Content System Error';
  }
}

/**
 * Get appropriate error description based on error type
 */
function getErrorDescription(errorType: ContentErrorType, isOnline: boolean): string {
  switch (errorType) {
    case ContentErrorType.NETWORK:
      return isOnline
        ? 'Unable to connect to the content service. This might be a temporary server issue.'
        : 'No internet connection detected. Please check your network connection.';
    case ContentErrorType.LOADING:
      return 'Content is taking longer than expected to load. This might be due to network conditions.';
    case ContentErrorType.PARSING:
      return 'The content format is invalid or corrupted. Please contact support if this persists.';
    case ContentErrorType.RTK_QUERY:
      return 'There was an issue with the data service. We\'re working to resolve this automatically.';
    case ContentErrorType.CONTENT_MISSING:
      return 'The requested content is currently unavailable. It might have been moved or deleted.';
    case ContentErrorType.PERMISSION:
      return 'You don\'t have permission to access this content. Please check your account status.';
    case ContentErrorType.VALIDATION:
      return 'The content failed validation checks. This indicates a data integrity issue.';
    default:
      return 'Something went wrong with the content system. This might be a temporary issue.';
  }
}

/**
 * Enhanced Error Boundary for Content System Operations
 *
 * Provides comprehensive error handling for content-related failures
 * with intelligent recovery, retry mechanisms, and detailed reporting.
 */
export class ContentErrorBoundary extends Component<Props, State> {
  private retryTimeoutId: NodeJS.Timeout | null = null;

  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      retryCount: 0,
      errorType: ContentErrorType.UNKNOWN,
      lastErrorTime: 0
    };
  }

  static getDerivedStateFromError(error: Error): Partial<State> {
    const { type, isRecoverable } = classifyContentError(error);

    return {
      hasError: true,
      error,
      errorType: type,
      lastErrorTime: Date.now()
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Update state with error information
    this.setState({
      error,
      errorInfo,
    });

    // Create comprehensive error report
    const contentError: ContentError = {
      type: this.state.errorType,
      message: error.message,
      originalError: error,
      timestamp: Date.now(),
      retryCount: this.state.retryCount,
      componentStack: errorInfo.componentStack,
      isRecoverable: classifyContentError(error).isRecoverable
    };

    // Report the error
    reportContentError(contentError);

    // Call user-provided error handler
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }

    // Call user-provided error reporter
    if (this.props.reportError) {
      this.props.reportError(contentError);
    }

    // Attempt automatic recovery for certain error types
    this.attemptAutoRecovery(contentError);
  }

  componentWillUnmount() {
    if (this.retryTimeoutId) {
      clearTimeout(this.retryTimeoutId);
    }
  }

  private attemptAutoRecovery = (contentError: ContentError) => {
    const maxRetries = this.props.maxRetries || 3;
    const retryDelay = this.props.retryDelay || 2000;

    // Only auto-retry for recoverable errors
    if (
      contentError.isRecoverable &&
      this.state.retryCount < maxRetries &&
      (contentError.type === ContentErrorType.NETWORK ||
       contentError.type === ContentErrorType.LOADING ||
       contentError.type === ContentErrorType.RTK_QUERY)
    ) {
      this.retryTimeoutId = setTimeout(() => {
        this.handleRetry();
      }, retryDelay * Math.pow(2, this.state.retryCount)); // Exponential backoff
    }
  };

  private handleReset = () => {
    if (this.retryTimeoutId) {
      clearTimeout(this.retryTimeoutId);
      this.retryTimeoutId = null;
    }

    this.setState({
      hasError: false,
      error: undefined,
      errorInfo: undefined,
      retryCount: 0,
      errorType: ContentErrorType.UNKNOWN,
      lastErrorTime: 0
    });
  };

  private handleReload = () => {
    window.location.reload();
  };

  private handleRetry = () => {
    const newRetryCount = this.state.retryCount + 1;
    const maxRetries = this.props.maxRetries || 3;

    if (newRetryCount <= maxRetries) {
      this.setState({
        retryCount: newRetryCount,
        lastErrorTime: Date.now()
      });

      // Try to recover by resetting the component
      setTimeout(() => {
        this.handleReset();
      }, 100);
    }
  };

  render() {
    if (this.state.hasError) {
      // Custom fallback UI
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // Enhanced error UI with recovery features
      return (
        <ErrorDisplay
          error={this.state.error}
          errorInfo={this.state.errorInfo}
          errorType={this.state.errorType}
          retryCount={this.state.retryCount}
          maxRetries={this.props.maxRetries || 3}
          onReset={this.handleReset}
          onReload={this.handleReload}
          onRetry={this.handleRetry}
          isRetrying={this.retryTimeoutId !== null}
        />
      );
    }

    return this.props.children;
  }
}

/**
 * Hook version of error boundary for functional components
 * Uses React Error Boundary library pattern
 */
export const withContentErrorBoundary = <P extends object>(
  Component: React.ComponentType<P>,
  fallback?: ReactNode
) => {
  const WrappedComponent = (props: P) => (
    <ContentErrorBoundary fallback={fallback}>
      <Component {...props} />
    </ContentErrorBoundary>
  );

  WrappedComponent.displayName = `withContentErrorBoundary(${Component.displayName || Component.name})`;

  return WrappedComponent;
};

/**
 * RTK Query Error Integration Hook
 *
 * Integrates content error boundary with RTK Query error states
 */
export const useRTKQueryErrorIntegration = () => {
  const handleRTKError = useCallback((error: any, query: string) => {
    if (error) {
      // Create a standardized error for the boundary
      const rtkError = new Error(`RTK Query Error in ${query}: ${error.message || 'Unknown error'}`);
      rtkError.name = 'RTKQueryError';

      // Trigger error boundary by throwing during render
      throw rtkError;
    }
  }, []);

  return { handleRTKError };
};

/**
 * Content Error Recovery Hook
 *
 * Provides recovery utilities for content-related operations
 */
export const useContentErrorRecovery = () => {
  const [recoveryAttempts, setRecoveryAttempts] = useState(0);
  const [isRecovering, setIsRecovering] = useState(false);

  const attemptRecovery = useCallback(async (
    recoveryFn: () => Promise<void>,
    maxAttempts: number = 3
  ) => {
    if (recoveryAttempts >= maxAttempts) {
      console.warn('Max recovery attempts reached');
      return false;
    }

    setIsRecovering(true);
    try {
      await recoveryFn();
      setRecoveryAttempts(0); // Reset on success
      setIsRecovering(false);
      return true;
    } catch (error) {
      setRecoveryAttempts(prev => prev + 1);
      setIsRecovering(false);
      console.error('Recovery attempt failed:', error);
      return false;
    }
  }, [recoveryAttempts]);

  const resetRecovery = useCallback(() => {
    setRecoveryAttempts(0);
    setIsRecovering(false);
  }, []);

  return {
    recoveryAttempts,
    isRecovering,
    attemptRecovery,
    resetRecovery
  };
};

/**
 * Performance-Optimized Content Error Boundary
 *
 * A lighter version for performance-critical content areas
 */
interface FastContentErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error) => void;
}

export class FastContentErrorBoundary extends Component<
  FastContentErrorBoundaryProps,
  { hasError: boolean }
> {
  constructor(props: FastContentErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(): { hasError: boolean } {
    return { hasError: true };
  }

  componentDidCatch(error: Error) {
    // Minimal logging for performance
    if (process.env.NODE_ENV === 'development') {
      console.error('Fast Content Error Boundary:', error.message);
    }

    if (this.props.onError) {
      this.props.onError(error);
    }
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="p-4 text-center text-sm text-muted-foreground">
          Content temporarily unavailable
        </div>
      );
    }

    return this.props.children;
  }
}

/**
 * Content Error Context for sharing error state across components
 */
export interface ContentErrorContextValue {
  reportError: (error: ContentError) => void;
  clearErrors: () => void;
  errorHistory: ContentError[];
  hasRecentErrors: boolean;
}

export const ContentErrorContext = React.createContext<ContentErrorContextValue | null>(null);

/**
 * Content Error Provider for centralized error management
 */
export const ContentErrorProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [errorHistory, setErrorHistory] = useState<ContentError[]>([]);

  const reportError = useCallback((error: ContentError) => {
    setErrorHistory(prev => {
      const newHistory = [...prev, error];
      // Keep only last 10 errors for memory management
      return newHistory.slice(-10);
    });
  }, []);

  const clearErrors = useCallback(() => {
    setErrorHistory([]);
  }, []);

  const hasRecentErrors = useMemo(() => {
    const recentThreshold = Date.now() - 5 * 60 * 1000; // 5 minutes
    return errorHistory.some(error => error.timestamp > recentThreshold);
  }, [errorHistory]);

  const contextValue: ContentErrorContextValue = {
    reportError,
    clearErrors,
    errorHistory,
    hasRecentErrors
  };

  return (
    <ContentErrorContext.Provider value={contextValue}>
      {children}
    </ContentErrorContext.Provider>
  );
};

/**
 * Hook to use content error context
 */
export const useContentErrorContext = () => {
  const context = useContext(ContentErrorContext);
  if (!context) {
    throw new Error('useContentErrorContext must be used within ContentErrorProvider');
  }
  return context;
};

/**
 * Higher-order component for automatic error boundary wrapping
 */
export const withContentErrorRecovery = <P extends object>(
  Component: React.ComponentType<P>,
  options?: {
    maxRetries?: number;
    retryDelay?: number;
    fallback?: ReactNode;
    fastMode?: boolean;
  }
) => {
  const WrappedComponent = (props: P) => {
    const ErrorBoundary = options?.fastMode ? FastContentErrorBoundary : ContentErrorBoundary;

    return (
      <ErrorBoundary
        maxRetries={options?.maxRetries}
        retryDelay={options?.retryDelay}
        fallback={options?.fallback}
      >
        <Component {...props} />
      </ErrorBoundary>
    );
  };

  WrappedComponent.displayName = `withContentErrorRecovery(${Component.displayName || Component.name})`;
  return WrappedComponent;
};

export default ContentErrorBoundary;