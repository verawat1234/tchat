/**
 * Microservice Configuration for Direct Service Access
 *
 * Configuration for connecting to backend services running on specific ports
 * instead of using the API gateway. Each service runs on its own port.
 */

export const SERVICE_CONFIG = {
  // Gateway for general API access
  gateway: {
    baseUrl: 'http://localhost:8080',
    apiPrefix: '/api/v1'
  },

  // Individual microservices
  services: {
    auth: {
      port: 8081,
      baseUrl: 'http://localhost:8081',
      healthEndpoint: '/health',
      apiPrefix: '/api/v1'
    },
    commerce: {
      port: 8082,
      baseUrl: 'http://localhost:8082',
      healthEndpoint: '/health',
      apiPrefix: '/api/v1'
    },
    content: {
      port: 8083,
      baseUrl: 'http://localhost:8083',
      healthEndpoint: '/health',
      apiPrefix: '/api/v1'
    },
    messaging: {
      port: 8084,
      baseUrl: 'http://localhost:8084',
      healthEndpoint: '/health',
      apiPrefix: '/api/v1'
    },
    notification: {
      port: 8085,
      baseUrl: 'http://localhost:8085',
      healthEndpoint: '/health',
      apiPrefix: '/api/v1'
    },
    payment: {
      port: 8086,
      baseUrl: 'http://localhost:8086',
      healthEndpoint: '/health',
      apiPrefix: '/api/v1'
    },
    video: {
      port: 8091,
      baseUrl: 'http://localhost:8091',
      healthEndpoint: '/health',
      apiPrefix: '/api/v1'
    }
  }
} as const;

/**
 * Service endpoint mappings for API routing
 */
export const SERVICE_ROUTES = {
  // Authentication & User Management
  '/auth': 'auth',
  '/users': 'auth',
  '/profiles': 'auth',

  // Commerce & Products
  '/products': 'commerce',
  '/orders': 'commerce',
  '/cart': 'commerce',
  '/shop': 'commerce',

  // Content Management
  '/content': 'content',
  '/categories': 'content',
  '/versions': 'content',

  // Messaging & Chat
  '/chats': 'messaging',
  '/messages': 'messaging',
  '/dialogs': 'messaging',

  // Notifications
  '/notifications': 'notification',
  '/alerts': 'notification',

  // Payments & Wallet
  '/payments': 'payment',
  '/wallet': 'payment',
  '/transactions': 'payment',

  // Video & Media
  '/video': 'video',
  '/videos': 'video',
  '/channels': 'video',
  '/playlists': 'video',
  '/comments': 'video',
  '/livestreams': 'video'
} as const;

/**
 * Get the appropriate service configuration for a given endpoint
 */
export function getServiceForEndpoint(endpoint: string): keyof typeof SERVICE_CONFIG.services | 'gateway' {
  // Remove query parameters and normalize endpoint
  const normalizedEndpoint = endpoint.split('?')[0];

  // Find matching service route
  for (const [route, service] of Object.entries(SERVICE_ROUTES)) {
    if (normalizedEndpoint.startsWith(route)) {
      return service as keyof typeof SERVICE_CONFIG.services;
    }
  }

  // Default to gateway for unknown routes
  return 'gateway';
}

/**
 * Build the full URL for a service endpoint
 */
export function buildServiceUrl(endpoint: string): string {
  const service = getServiceForEndpoint(endpoint);

  if (service === 'gateway') {
    return `${SERVICE_CONFIG.gateway.baseUrl}${SERVICE_CONFIG.gateway.apiPrefix}${endpoint}`;
  }

  const serviceConfig = SERVICE_CONFIG.services[service];
  return `${serviceConfig.baseUrl}${serviceConfig.apiPrefix}${endpoint}`;
}

/**
 * Health check configuration for service discovery
 */
export const HEALTH_CHECK_CONFIG = {
  interval: 30000, // 30 seconds
  timeout: 5000,   // 5 seconds
  retries: 3
};

/**
 * Service availability tracking
 */
export interface ServiceHealth {
  service: string;
  available: boolean;
  lastChecked: number;
  responseTime?: number;
  error?: string;
}

/**
 * Environment-based service configuration
 */
export function getServiceConfig() {
  const isDevelopment = import.meta.env.DEV;
  const useDirectServices = import.meta.env.VITE_USE_DIRECT_SERVICES === 'true';

  return {
    useDirect: isDevelopment && useDirectServices,
    fallbackToGateway: true,
    healthCheckEnabled: isDevelopment
  };
}