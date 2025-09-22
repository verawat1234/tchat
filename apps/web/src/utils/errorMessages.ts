import { store } from '../store';
import { contentApi } from '../services/content';

/**
 * Error message content ID mapping
 *
 * Maps error types and HTTP status codes to content IDs following
 * the pattern: error.{type}.{element}
 */
export const ERROR_CONTENT_IDS = {
  // HTTP Status Code Errors
  HTTP_400: 'error.http.bad_request',
  HTTP_401: 'error.http.unauthorized',
  HTTP_403: 'error.http.forbidden',
  HTTP_404: 'error.http.not_found',
  HTTP_409: 'error.http.conflict',
  HTTP_422: 'error.http.validation',
  HTTP_429: 'error.http.rate_limit',
  HTTP_500: 'error.http.server_error',
  HTTP_503: 'error.http.service_unavailable',

  // Network Errors
  NETWORK_ERROR: 'error.network.connection',
  TIMEOUT_ERROR: 'error.network.timeout',

  // Message-specific Errors
  MESSAGE_TOO_LARGE: 'error.message.too_large',
  MESSAGE_RATE_LIMIT: 'error.message.rate_limit',
  MESSAGE_SEND_FAILED: 'error.message.send_failed',

  // Validation Errors
  VALIDATION_GENERIC: 'error.validation.generic',

  // Image Errors
  IMAGE_LOAD_FAILED: 'error.image.load_failed',

  // Generic Fallback
  GENERIC: 'error.generic.unexpected',
} as const;

/**
 * Fallback error messages for critical scenarios
 *
 * These are used when the content API is unavailable or fails,
 * ensuring users always see helpful error messages.
 */
export const FALLBACK_ERROR_MESSAGES = {
  [ERROR_CONTENT_IDS.HTTP_400]: 'Invalid request. Please check your input.',
  [ERROR_CONTENT_IDS.HTTP_401]: 'Please log in to continue.',
  [ERROR_CONTENT_IDS.HTTP_403]: 'You don\'t have permission to perform this action.',
  [ERROR_CONTENT_IDS.HTTP_404]: 'The requested resource was not found.',
  [ERROR_CONTENT_IDS.HTTP_409]: 'This action conflicts with existing data.',
  [ERROR_CONTENT_IDS.HTTP_422]: 'Please check your input and try again.',
  [ERROR_CONTENT_IDS.HTTP_429]: 'Too many requests. Please wait a moment.',
  [ERROR_CONTENT_IDS.HTTP_500]: 'Server error. Please try again later.',
  [ERROR_CONTENT_IDS.HTTP_503]: 'Service unavailable. Please try again later.',
  [ERROR_CONTENT_IDS.NETWORK_ERROR]: 'Network error. Please check your connection.',
  [ERROR_CONTENT_IDS.TIMEOUT_ERROR]: 'Request timed out. Please try again.',
  [ERROR_CONTENT_IDS.MESSAGE_TOO_LARGE]: 'Message too large. Please reduce the size and try again.',
  [ERROR_CONTENT_IDS.MESSAGE_RATE_LIMIT]: 'Sending messages too quickly. Please wait a moment.',
  [ERROR_CONTENT_IDS.MESSAGE_SEND_FAILED]: 'Failed to send message.',
  [ERROR_CONTENT_IDS.VALIDATION_GENERIC]: 'Please check your input and try again.',
  [ERROR_CONTENT_IDS.IMAGE_LOAD_FAILED]: 'Error loading image',
  [ERROR_CONTENT_IDS.GENERIC]: 'An unexpected error occurred.',
} as const;

/**
 * Get dynamic error message with fallback
 *
 * This function attempts to fetch a dynamic error message from the content API.
 * If the content is not available, it returns the appropriate fallback message.
 * For critical error scenarios, it ensures robust fallback behavior.
 */
export async function getDynamicErrorMessage(
  contentId: keyof typeof ERROR_CONTENT_IDS,
  fallbackMessage?: string
): Promise<string> {
  try {
    // Try to get from cache first
    const state = store.getState();
    const cached = contentApi.endpoints.getContentItem.select(ERROR_CONTENT_IDS[contentId])(state);

    if (cached.data?.data && typeof cached.data.data === 'string') {
      return cached.data.data;
    }

    // If not cached, initiate fetch but don't wait for critical errors
    if (!cached.isLoading) {
      store.dispatch(contentApi.endpoints.getContentItem.initiate(ERROR_CONTENT_IDS[contentId]));
    }

    // Return fallback immediately for error scenarios
    return fallbackMessage || FALLBACK_ERROR_MESSAGES[ERROR_CONTENT_IDS[contentId]] || FALLBACK_ERROR_MESSAGES[ERROR_CONTENT_IDS.GENERIC];
  } catch (error) {
    // If anything goes wrong, use fallback
    console.warn('Failed to get dynamic error message:', error);
    return fallbackMessage || FALLBACK_ERROR_MESSAGES[ERROR_CONTENT_IDS[contentId]] || FALLBACK_ERROR_MESSAGES[ERROR_CONTENT_IDS.GENERIC];
  }
}

/**
 * Synchronous error message getter with fallback
 *
 * For scenarios where async error message fetching is not practical,
 * this function provides immediate fallback messages while potentially
 * triggering background content fetches.
 */
export function getErrorMessageSync(
  contentId: keyof typeof ERROR_CONTENT_IDS,
  fallbackMessage?: string
): string {
  try {
    // Check if we have cached content
    const state = store.getState();
    const cached = contentApi.endpoints.getContentItem.select(ERROR_CONTENT_IDS[contentId])(state);

    if (cached.data?.data && typeof cached.data.data === 'string') {
      return cached.data.data;
    }

    // If not cached and not loading, initiate background fetch
    if (!cached.isLoading) {
      store.dispatch(contentApi.endpoints.getContentItem.initiate(ERROR_CONTENT_IDS[contentId]));
    }

    // Return fallback immediately
    return fallbackMessage || FALLBACK_ERROR_MESSAGES[ERROR_CONTENT_IDS[contentId]] || FALLBACK_ERROR_MESSAGES[ERROR_CONTENT_IDS.GENERIC];
  } catch (error) {
    // If anything goes wrong, use fallback
    console.warn('Failed to get error message:', error);
    return fallbackMessage || FALLBACK_ERROR_MESSAGES[ERROR_CONTENT_IDS[contentId]] || FALLBACK_ERROR_MESSAGES[ERROR_CONTENT_IDS.GENERIC];
  }
}

/**
 * Map HTTP status codes to content IDs
 */
export function getErrorContentIdForStatus(status: number): keyof typeof ERROR_CONTENT_IDS {
  switch (status) {
    case 400: return 'HTTP_400';
    case 401: return 'HTTP_401';
    case 403: return 'HTTP_403';
    case 404: return 'HTTP_404';
    case 409: return 'HTTP_409';
    case 422: return 'HTTP_422';
    case 429: return 'HTTP_429';
    case 500: return 'HTTP_500';
    case 503: return 'HTTP_503';
    default: return 'GENERIC';
  }
}