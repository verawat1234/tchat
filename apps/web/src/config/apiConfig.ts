/**
 * API Configuration
 *
 * Centralized API configuration that builds URLs from environment variables.
 * Supports flexible configuration for different environments (development, staging, production).
 *
 * Environment Variables:
 * - VITE_GATEWAY_BASE_URL: Base URL of the API gateway (e.g., https://gateway-service-production-d78d.up.railway.app)
 * - VITE_API_SERVICE_PATH: Service path segment (default: 'api')
 * - VITE_API_VERSION: API version (default: 'v1')
 * - VITE_API_URL: Complete API URL (legacy fallback, computed from above if not set)
 */

/**
 * Get the gateway base URL from environment
 */
export const getGatewayBaseUrl = (): string => {
  return import.meta.env.VITE_GATEWAY_BASE_URL || 'http://localhost:8080';
};

/**
 * Get the API service path from environment
 */
export const getApiServicePath = (): string => {
  return import.meta.env.VITE_API_SERVICE_PATH || 'api';
};

/**
 * Get the API version from environment
 */
export const getApiVersion = (): string => {
  return import.meta.env.VITE_API_VERSION || 'v1';
};

/**
 * Build the complete API URL from environment variables
 *
 * Format: {GATEWAY_BASE_URL}/{API_SERVICE_PATH}/{API_VERSION}
 * Example: https://gateway-service-production-d78d.up.railway.app/api/v1
 *
 * Falls back to VITE_API_URL if set directly in environment.
 */
export const getApiUrl = (): string => {
  // Use direct VITE_API_URL if provided (legacy support)
  if (import.meta.env.VITE_API_URL && !import.meta.env.VITE_API_URL.includes('${')) {
    return import.meta.env.VITE_API_URL;
  }

  // Build from components
  const baseUrl = getGatewayBaseUrl();
  const servicePath = getApiServicePath();
  const version = getApiVersion();

  return `${baseUrl}/${servicePath}/${version}`;
};

/**
 * Get WebSocket URL from environment
 */
export const getWebSocketUrl = (): string => {
  return import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
};

/**
 * Check if using direct service access (bypassing gateway)
 */
export const useDirectServices = (): boolean => {
  return import.meta.env.VITE_USE_DIRECT_SERVICES === 'true';
};

/**
 * API configuration object
 */
export const apiConfig = {
  gatewayBaseUrl: getGatewayBaseUrl(),
  apiServicePath: getApiServicePath(),
  apiVersion: getApiVersion(),
  apiUrl: getApiUrl(),
  wsUrl: getWebSocketUrl(),
  useDirectServices: useDirectServices(),
  isDevelopment: import.meta.env.DEV,
  isProduction: import.meta.env.PROD,
  debug: import.meta.env.VITE_DEBUG === 'true',
} as const;

// Log configuration in development
if (import.meta.env.DEV && import.meta.env.VITE_DEBUG === 'true') {
  console.log('[API Config]', {
    gatewayBaseUrl: apiConfig.gatewayBaseUrl,
    apiUrl: apiConfig.apiUrl,
    wsUrl: apiConfig.wsUrl,
    useDirectServices: apiConfig.useDirectServices,
  });
}
