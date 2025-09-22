import { createListenerMiddleware, isRejected, isRejectedWithValue } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { addNotification } from '../../features/uiSlice';
import { api } from '../../services/api';
import type { ApiError } from '../../types/api';
import { getErrorMessageSync, getErrorContentIdForStatus, ERROR_CONTENT_IDS } from '../../utils/errorMessages';

export const errorMiddleware = createListenerMiddleware();

// Enhanced error message mapping with dynamic content
const getErrorMessage = (error: any): string => {
  // If the error has a specific message from the server, use it
  if (error?.data?.error?.message) {
    return error.data.error.message;
  }

  // For HTTP status codes, use dynamic content with fallbacks
  if (error?.status) {
    const contentId = getErrorContentIdForStatus(error.status);
    return getErrorMessageSync(contentId);
  }

  // For other error types, try to get dynamic content
  if (error?.message) {
    // For now, return the original message, but could be enhanced
    // to map common error messages to content IDs
    return error.message;
  }

  // Generic fallback with dynamic content
  return getErrorMessageSync('GENERIC');
};

// Enhanced error type checking
const isNetworkError = (error: any): boolean => {
  return error?.name === 'NetworkError' ||
         error?.message?.includes('fetch') ||
         error?.status === 0;
};

const isTimeoutError = (error: any): boolean => {
  return error?.name === 'TimeoutError' ||
         error?.message?.includes('timeout');
};

// Listen for any RTK Query errors
errorMiddleware.startListening({
  matcher: isRejectedWithValue,
  effect: async (action: PayloadAction<any>, listenerApi) => {
    const { error, meta } = action.payload || {};
    const endpointName = meta?.arg?.endpointName;

    // Skip errors for certain endpoints (like auth failures during login)
    const skipNotificationEndpoints = [
      'login',
      'refreshToken',
      'getCurrentUser'
    ];

    if (skipNotificationEndpoints.includes(endpointName)) {
      return;
    }

    let notificationType: 'error' | 'warning' = 'error';
    let message = getErrorMessage(error);

    // Special handling for different error types with dynamic content
    if (isNetworkError(error)) {
      message = getErrorMessageSync('NETWORK_ERROR');
      notificationType = 'warning';
    } else if (isTimeoutError(error)) {
      message = getErrorMessageSync('TIMEOUT_ERROR');
      notificationType = 'warning';
    } else if (error?.status === 429) {
      notificationType = 'warning';
    }

    // Add notification to UI
    listenerApi.dispatch(addNotification({
      type: notificationType,
      message,
      duration: 5000,
    }));

    // Log error for debugging
    console.error(`API Error [${endpointName}]:`, {
      error,
      action: action.type,
      timestamp: new Date().toISOString(),
    });
  },
});

// Listen for network errors specifically
errorMiddleware.startListening({
  matcher: isRejected,
  effect: async (action, listenerApi) => {
    const { error } = action;

    // Handle network connectivity issues
    if (isNetworkError(error)) {
      // Could implement retry logic here
      console.warn('Network connectivity issue detected');

      // Show persistent warning for network issues with dynamic content
      const retryMessage = getErrorMessageSync('NETWORK_ERROR', 'Connection lost. Retrying...');
      listenerApi.dispatch(addNotification({
        type: 'warning',
        message: retryMessage,
        duration: 10000,
      }));
    }
  },
});

// Handle specific API endpoint errors with custom logic
// Note: This will be activated when sendMessage endpoint is injected
errorMiddleware.startListening({
  predicate: (action) => {
    return action.type.endsWith('/rejected') &&
           action.type.includes('sendMessage');
  },
  effect: async (action, listenerApi) => {
    const { error } = action;

    // For message sending failures, show specific guidance with dynamic content
    if (error && 'status' in error) {
      let message = getErrorMessageSync('MESSAGE_SEND_FAILED');

      if (error.status === 413) {
        message = getErrorMessageSync('MESSAGE_TOO_LARGE');
      } else if (error.status === 429) {
        message = getErrorMessageSync('MESSAGE_RATE_LIMIT');
      }

      listenerApi.dispatch(addNotification({
        type: 'error',
        message,
        duration: 5000,
      }));
    }
  },
});

// Handle validation errors with detailed feedback
errorMiddleware.startListening({
  predicate: (action) => {
    return isRejectedWithValue(action) &&
           action.payload?.error?.status === 422;
  },
  effect: async (action, listenerApi) => {
    const { error } = action.payload;
    const details = error?.data?.error?.details;

    if (details && typeof details === 'object') {
      // Show field-specific validation errors
      Object.entries(details).forEach(([field, message]) => {
        listenerApi.dispatch(addNotification({
          type: 'error',
          message: `${field}: ${message}`,
          duration: 7000,
        }));
      });
    } else {
      // Fallback to generic validation error with dynamic content
      const validationMessage = getErrorMessageSync('VALIDATION_GENERIC');
      listenerApi.dispatch(addNotification({
        type: 'error',
        message: validationMessage,
        duration: 5000,
      }));
    }
  },
});

// Auto-clear old notifications
errorMiddleware.startListening({
  actionCreator: addNotification,
  effect: async (action, listenerApi) => {
    const { id, duration = 5000 } = action.payload;

    // Auto-remove notification after duration
    setTimeout(() => {
      listenerApi.dispatch({
        type: 'ui/removeNotification',
        payload: id,
      });
    }, duration);
  },
});