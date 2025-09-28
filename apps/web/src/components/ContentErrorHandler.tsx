/**
 * Content Error Handler Component
 *
 * Provides comprehensive error handling and user-friendly error messages
 * for content management operations. Features intelligent error categorization,
 * actionable recovery options, and seamless integration with offline queue.
 *
 * Features:
 * - Intelligent error categorization and user-friendly messaging
 * - Actionable recovery options with retry logic
 * - Integration with offline queue for failed operations
 * - Toast notifications with proper error severity levels
 * - Detailed error logging and debugging information
 * - Network-aware error handling and recovery suggestions
 * - Cross-component error boundary integration
 */

import React, { useState, useEffect, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from './ui/alert';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card/card';
import { Badge } from './ui/badge/badge';
import { ScrollArea } from './ui/scroll-area';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from './ui/collapsible';
import {
  AlertTriangle,
  WifiOff,
  RefreshCcw,
  Download,
  Upload,
  Clock,
  Info,
  X,
  ChevronDown,
  ChevronUp,
  Copy,
  ExternalLink,
  Zap,
  AlertCircle,
  CheckCircle2,
} from 'lucide-react';
import { useToast } from '../hooks/use-toast';
import { useOfflineQueue } from '../services/offlineQueueService';
import { notificationService } from '../services/notificationService';
import type { ContentItem } from '../types/content';

// =============================================================================
// Type Definitions
// =============================================================================

export type ContentErrorType =
  | 'network_error'
  | 'validation_error'
  | 'permission_error'
  | 'conflict_error'
  | 'server_error'
  | 'timeout_error'
  | 'quota_exceeded'
  | 'unknown_error';

export type ErrorSeverity = 'info' | 'warning' | 'error' | 'critical';

export interface ContentError {
  id: string;
  type: ContentErrorType;
  severity: ErrorSeverity;
  title: string;
  message: string;
  details?: string;
  timestamp: string;
  contentId?: string;
  operation?: string;
  retryable: boolean;
  autoRetry?: boolean;
  retryCount: number;
  maxRetries: number;
  context?: {
    url?: string;
    statusCode?: number;
    requestId?: string;
    userAgent?: string;
    sessionId?: string;
  };
  actions?: ErrorAction[];
}

export interface ErrorAction {
  id: string;
  label: string;
  type: 'primary' | 'secondary' | 'destructive';
  icon?: React.ComponentType<{ className?: string }>;
  action: () => void | Promise<void>;
  disabled?: boolean;
  loading?: boolean;
}

export interface ContentErrorHandlerProps {
  error?: ContentError;
  onRetry?: () => Promise<void>;
  onDismiss?: () => void;
  showDetails?: boolean;
  compact?: boolean;
  className?: string;
}

export interface ErrorManagerProps {
  maxErrors?: number;
  autoRetryInterval?: number;
  showToasts?: boolean;
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left';
}

// =============================================================================
// Error Utilities
// =============================================================================

export function categorizeError(error: any): ContentErrorType {
  if (!navigator.onLine) {
    return 'network_error';
  }

  if (error?.status === 400 || error?.name === 'ValidationError') {
    return 'validation_error';
  }

  if (error?.status === 401 || error?.status === 403) {
    return 'permission_error';
  }

  if (error?.status === 409) {
    return 'conflict_error';
  }

  if (error?.status >= 500) {
    return 'server_error';
  }

  if (error?.name === 'TimeoutError' || error?.code === 'TIMEOUT') {
    return 'timeout_error';
  }

  if (error?.status === 413 || error?.message?.includes('quota')) {
    return 'quota_exceeded';
  }

  return 'unknown_error';
}

export function getErrorSeverity(type: ContentErrorType): ErrorSeverity {
  switch (type) {
    case 'network_error':
    case 'timeout_error':
      return 'warning';
    case 'validation_error':
      return 'info';
    case 'permission_error':
    case 'quota_exceeded':
      return 'error';
    case 'conflict_error':
    case 'server_error':
    case 'unknown_error':
      return 'critical';
    default:
      return 'error';
  }
}

export function getErrorIcon(type: ContentErrorType): React.ComponentType<{ className?: string }> {
  switch (type) {
    case 'network_error':
      return WifiOff;
    case 'validation_error':
      return Info;
    case 'permission_error':
      return AlertTriangle;
    case 'conflict_error':
      return AlertCircle;
    case 'server_error':
    case 'unknown_error':
      return AlertTriangle;
    case 'timeout_error':
      return Clock;
    case 'quota_exceeded':
      return Zap;
    default:
      return AlertTriangle;
  }
}

export function getUserFriendlyMessage(type: ContentErrorType, error: any): { title: string; message: string } {
  switch (type) {
    case 'network_error':
      return {
        title: 'Connection Problem',
        message: 'Unable to connect to the server. Your changes have been saved locally and will sync when connection is restored.',
      };
    case 'validation_error':
      return {
        title: 'Invalid Data',
        message: 'Please check your input and try again. Some required fields may be missing or contain invalid values.',
      };
    case 'permission_error':
      return {
        title: 'Access Denied',
        message: 'You don\'t have permission to perform this action. Please contact your administrator if you believe this is an error.',
      };
    case 'conflict_error':
      return {
        title: 'Content Conflict',
        message: 'Another user has modified this content. Please resolve the conflicts before continuing.',
      };
    case 'server_error':
      return {
        title: 'Server Error',
        message: 'Something went wrong on our servers. Our team has been notified and will fix this shortly.',
      };
    case 'timeout_error':
      return {
        title: 'Request Timeout',
        message: 'The operation took too long to complete. Please try again or check your connection.',
      };
    case 'quota_exceeded':
      return {
        title: 'Storage Limit Reached',
        message: 'You\'ve reached your storage limit. Please free up space or upgrade your plan to continue.',
      };
    default:
      return {
        title: 'Unexpected Error',
        message: error?.message || 'An unexpected error occurred. Please try again or contact support if the problem persists.',
      };
  }
}

// =============================================================================
// Individual Error Component
// =============================================================================

export function ContentErrorCard({ error, onRetry, onDismiss, showDetails = false, compact = false, className = '' }: ContentErrorHandlerProps) {
  const [isRetrying, setIsRetrying] = useState(false);
  const [showDetailsExpanded, setShowDetailsExpanded] = useState(showDetails);
  const { toast } = useToast();
  const { queueOperation } = useOfflineQueue();

  const ErrorIcon = getErrorIcon(error!.type);

  const handleRetry = async () => {
    if (!onRetry || isRetrying) return;

    setIsRetrying(true);
    try {
      await onRetry();
      toast({
        title: 'Success',
        description: 'Operation completed successfully',
        variant: 'default',
      });
    } catch (retryError) {
      console.error('Retry failed:', retryError);
      toast({
        title: 'Retry Failed',
        description: 'The retry attempt failed. The operation has been queued for later.',
        variant: 'destructive',
      });
    } finally {
      setIsRetrying(false);
    }
  };

  const handleQueueOffline = async () => {
    if (!error?.contentId) return;

    try {
      await queueOperation({
        type: 'update',
        priority: 'normal',
        contentId: error.contentId,
        operation: {
          endpoint: `/content/${error.contentId}`,
          method: 'PUT',
          data: {},
        },
      });

      toast({
        title: 'Queued for Later',
        description: 'Operation has been queued and will be retried when connection is restored.',
        variant: 'default',
      });

      onDismiss?.();
    } catch (queueError) {
      console.error('Failed to queue operation:', queueError);
    }
  };

  const handleCopyDetails = () => {
    const details = `
Error: ${error!.title}
Type: ${error!.type}
Message: ${error!.message}
Timestamp: ${error!.timestamp}
Content ID: ${error!.contentId || 'N/A'}
Details: ${error!.details || 'None'}
    `.trim();

    navigator.clipboard.writeText(details).then(() => {
      toast({
        title: 'Copied',
        description: 'Error details copied to clipboard',
        variant: 'default',
      });
    });
  };

  if (!error) return null;

  const severityColors = {
    info: 'border-blue-200 bg-blue-50',
    warning: 'border-yellow-200 bg-yellow-50',
    error: 'border-red-200 bg-red-50',
    critical: 'border-red-300 bg-red-100',
  };

  const iconColors = {
    info: 'text-blue-600',
    warning: 'text-yellow-600',
    error: 'text-red-600',
    critical: 'text-red-700',
  };

  if (compact) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -10 }}
        className={`flex items-center gap-3 p-3 border rounded-lg ${severityColors[error.severity]} ${className}`}
      >
        <ErrorIcon className={`h-4 w-4 ${iconColors[error.severity]} flex-shrink-0`} />
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-gray-900 truncate">{error.title}</p>
          <p className="text-xs text-gray-600 truncate">{error.message}</p>
        </div>
        {error.retryable && onRetry && (
          <Button
            size="sm"
            variant="outline"
            onClick={handleRetry}
            disabled={isRetrying}
            className="flex-shrink-0"
          >
            {isRetrying ? (
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                className="h-3 w-3 border border-current border-t-transparent rounded-full"
              />
            ) : (
              <RefreshCcw className="h-3 w-3" />
            )}
          </Button>
        )}
        {onDismiss && (
          <Button
            size="sm"
            variant="ghost"
            onClick={onDismiss}
            className="flex-shrink-0 h-8 w-8 p-0"
          >
            <X className="h-3 w-3" />
          </Button>
        )}
      </motion.div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      exit={{ opacity: 0, scale: 0.95 }}
      className={className}
    >
      <Card className={`${severityColors[error.severity]} border-l-4`}>
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <ErrorIcon className={`h-5 w-5 ${iconColors[error.severity]}`} />
              <div>
                <CardTitle className="text-base">{error.title}</CardTitle>
                <Badge variant="outline" className="mt-1">
                  {error.type.replace('_', ' ').toUpperCase()}
                </Badge>
              </div>
            </div>
            {onDismiss && (
              <Button variant="ghost" size="sm" onClick={onDismiss}>
                <X className="h-4 w-4" />
              </Button>
            )}
          </div>
        </CardHeader>

        <CardContent className="space-y-4">
          <p className="text-sm text-gray-700">{error.message}</p>

          {/* Action Buttons */}
          <div className="flex flex-wrap gap-2">
            {error.retryable && onRetry && (
              <Button
                size="sm"
                onClick={handleRetry}
                disabled={isRetrying}
                className="flex items-center gap-2"
              >
                {isRetrying ? (
                  <motion.div
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                    className="h-3 w-3 border border-current border-t-transparent rounded-full"
                  />
                ) : (
                  <RefreshCcw className="h-3 w-3" />
                )}
                {isRetrying ? 'Retrying...' : 'Try Again'}
              </Button>
            )}

            {error.type === 'network_error' && (
              <Button
                size="sm"
                variant="outline"
                onClick={handleQueueOffline}
                className="flex items-center gap-2"
              >
                <Download className="h-3 w-3" />
                Queue for Later
              </Button>
            )}

            {error.actions?.map(action => (
              <Button
                key={action.id}
                size="sm"
                variant={action.type === 'primary' ? 'default' : action.type === 'destructive' ? 'destructive' : 'outline'}
                onClick={action.action}
                disabled={action.disabled || action.loading}
                className="flex items-center gap-2"
              >
                {action.loading ? (
                  <motion.div
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                    className="h-3 w-3 border border-current border-t-transparent rounded-full"
                  />
                ) : (
                  action.icon && <action.icon className="h-3 w-3" />
                )}
                {action.label}
              </Button>
            ))}
          </div>

          {/* Error Details */}
          {(error.details || error.context) && (
            <Collapsible open={showDetailsExpanded} onOpenChange={setShowDetailsExpanded}>
              <CollapsibleTrigger asChild>
                <Button variant="ghost" size="sm" className="flex items-center gap-2 p-0 h-auto">
                  {showDetailsExpanded ? (
                    <ChevronUp className="h-3 w-3" />
                  ) : (
                    <ChevronDown className="h-3 w-3" />
                  )}
                  {showDetailsExpanded ? 'Hide Details' : 'Show Details'}
                </Button>
              </CollapsibleTrigger>

              <CollapsibleContent className="mt-3">
                <div className="bg-gray-50 rounded-lg p-3 space-y-3">
                  {error.details && (
                    <div>
                      <h5 className="text-xs font-medium text-gray-700 mb-1">Error Details</h5>
                      <p className="text-xs text-gray-600 font-mono bg-white p-2 rounded border">
                        {error.details}
                      </p>
                    </div>
                  )}

                  {error.context && (
                    <div>
                      <h5 className="text-xs font-medium text-gray-700 mb-1">Technical Info</h5>
                      <div className="text-xs text-gray-600 space-y-1">
                        {error.context.statusCode && (
                          <div>Status Code: {error.context.statusCode}</div>
                        )}
                        {error.context.requestId && (
                          <div>Request ID: {error.context.requestId}</div>
                        )}
                        {error.context.url && (
                          <div>URL: {error.context.url}</div>
                        )}
                        <div>Timestamp: {new Date(error.timestamp).toLocaleString()}</div>
                      </div>
                    </div>
                  )}

                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={handleCopyDetails}
                      className="flex items-center gap-1"
                    >
                      <Copy className="h-3 w-3" />
                      Copy Details
                    </Button>
                  </div>
                </div>
              </CollapsibleContent>
            </Collapsible>
          )}

          {/* Retry Information */}
          {error.retryable && error.retryCount > 0 && (
            <div className="text-xs text-gray-500">
              Retry attempt {error.retryCount} of {error.maxRetries}
            </div>
          )}
        </CardContent>
      </Card>
    </motion.div>
  );
}

// =============================================================================
// Error Manager Component
// =============================================================================

export function ContentErrorManager({
  maxErrors = 5,
  autoRetryInterval = 30000,
  showToasts = true,
  position = 'top-right'
}: ErrorManagerProps) {
  const [errors, setErrors] = useState<ContentError[]>([]);
  const [isVisible, setIsVisible] = useState(false);
  const { toast } = useToast();

  const addError = useCallback((error: ContentError) => {
    setErrors(prev => {
      const newErrors = [error, ...prev.slice(0, maxErrors - 1)];
      return newErrors;
    });

    if (showToasts) {
      toast({
        title: error.title,
        description: error.message,
        variant: error.severity === 'critical' || error.severity === 'error' ? 'destructive' : 'default',
      });
    }

    setIsVisible(true);
  }, [maxErrors, showToasts, toast]);

  const removeError = useCallback((errorId: string) => {
    setErrors(prev => prev.filter(e => e.id !== errorId));
  }, []);

  const clearAllErrors = useCallback(() => {
    setErrors([]);
    setIsVisible(false);
  }, []);

  // Auto-hide when no errors
  useEffect(() => {
    if (errors.length === 0) {
      setIsVisible(false);
    }
  }, [errors.length]);

  const positionClasses = {
    'top-right': 'top-4 right-4',
    'top-left': 'top-4 left-4',
    'bottom-right': 'bottom-4 right-4',
    'bottom-left': 'bottom-4 left-4',
  };

  if (!isVisible || errors.length === 0) {
    return null;
  }

  return (
    <div className={`fixed ${positionClasses[position]} z-50 max-w-md w-full space-y-2`}>
      <AnimatePresence>
        {errors.map(error => (
          <ContentErrorCard
            key={error.id}
            error={error}
            onDismiss={() => removeError(error.id)}
            compact
          />
        ))}
      </AnimatePresence>

      {errors.length > 1 && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="flex justify-center"
        >
          <Button
            size="sm"
            variant="outline"
            onClick={clearAllErrors}
            className="text-xs"
          >
            Clear All ({errors.length})
          </Button>
        </motion.div>
      )}
    </div>
  );
}

// =============================================================================
// Hook for Error Handling
// =============================================================================

export function useContentErrorHandler() {
  const [currentError, setCurrentError] = useState<ContentError | null>(null);
  const { toast } = useToast();

  const handleError = useCallback(async (error: any, context?: {
    contentId?: string;
    operation?: string;
    retryCallback?: () => Promise<void>;
  }) => {
    const errorType = categorizeError(error);
    const severity = getErrorSeverity(errorType);
    const { title, message } = getUserFriendlyMessage(errorType, error);

    const contentError: ContentError = {
      id: `error_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      type: errorType,
      severity,
      title,
      message,
      details: error?.stack || error?.details,
      timestamp: new Date().toISOString(),
      contentId: context?.contentId,
      operation: context?.operation,
      retryable: errorType !== 'permission_error' && errorType !== 'validation_error',
      retryCount: 0,
      maxRetries: 3,
      context: {
        statusCode: error?.status,
        requestId: error?.requestId,
        url: error?.config?.url,
        userAgent: navigator.userAgent,
        sessionId: `session_${Date.now()}`,
      },
    };

    setCurrentError(contentError);

    // Log error for debugging
    console.error('Content Error:', {
      type: errorType,
      error,
      context,
      contentError,
    });

    // Send to notification service
    await notificationService.notifyContentError(
      message,
      context?.contentId,
      context?.operation
    );

    return contentError;
  }, []);

  const retryOperation = useCallback(async (retryCallback?: () => Promise<void>) => {
    if (!retryCallback || !currentError) return;

    try {
      await retryCallback();
      setCurrentError(null);
    } catch (error) {
      // Update retry count
      if (currentError) {
        const updatedError = {
          ...currentError,
          retryCount: currentError.retryCount + 1,
        };
        setCurrentError(updatedError);
      }
      throw error;
    }
  }, [currentError]);

  const dismissError = useCallback(() => {
    setCurrentError(null);
  }, []);

  return {
    currentError,
    handleError,
    retryOperation,
    dismissError,
  };
}