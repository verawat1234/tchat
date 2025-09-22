/**
 * MSW (Mock Service Worker) Configuration
 *
 * This module provides mock API endpoints for development and testing.
 * It includes both Node.js server setup for tests and browser worker setup for development.
 */

// Server setup for Node.js/testing environment
export {
  server,
  startMockServer,
  stopMockServer,
  resetMockServer,
  useMockHandlers as useServerHandlers,
} from './server';

// Browser worker setup for development
export {
  worker,
  startMockWorker,
  stopMockWorker,
  resetMockWorker,
  useMockHandlers as useWorkerHandlers,
} from './browser';

// Request handlers and mock data factories
export {
  handlers,
  errorHandlers,
  withDelay,
  createMockUser,
  createMockMessage,
  createMockChannel,
  createMockConversation,
} from './handlers';

// Type definitions for mock data
export type MockUser = ReturnType<typeof import('./handlers').createMockUser>;
export type MockMessage = ReturnType<typeof import('./handlers').createMockMessage>;
export type MockChannel = ReturnType<typeof import('./handlers').createMockChannel>;
export type MockConversation = ReturnType<typeof import('./handlers').createMockConversation>;

/**
 * Environment detection utilities
 */
export const isTestEnvironment = () =>
  typeof process !== 'undefined' && process.env.NODE_ENV === 'test';

export const isDevelopmentEnvironment = () =>
  typeof process !== 'undefined' && process.env.NODE_ENV === 'development';

/**
 * Auto-setup function for easy initialization
 *
 * @param environment - The environment to setup mocks for ('test' | 'development' | 'auto')
 */
export const setupMocks = async (environment: 'test' | 'development' | 'auto' = 'auto') => {
  if (environment === 'auto') {
    if (isTestEnvironment()) {
      environment = 'test';
    } else if (isDevelopmentEnvironment()) {
      environment = 'development';
    } else {
      console.warn('[MSW] Auto-detection failed, defaulting to development mode');
      environment = 'development';
    }
  }

  if (environment === 'test') {
    // Node.js server for testing
    const { startMockServer } = await import('./server');
    startMockServer({ onUnhandledRequest: 'error' });
  } else if (environment === 'development') {
    // Browser worker for development
    const { startMockWorker } = await import('./browser');
    await startMockWorker({ onUnhandledRequest: 'bypass', quiet: false });
  }
};