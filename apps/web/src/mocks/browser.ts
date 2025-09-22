import { setupWorker } from 'msw/browser';
import { handlers } from './handlers';

/**
 * MSW Browser Worker Configuration
 *
 * This worker is used for mocking API endpoints in the browser during development.
 * It intercepts network requests and returns mock responses.
 */
export const worker = setupWorker(...handlers);

/**
 * Start the mock service worker in the browser
 *
 * @param options - Configuration options for the worker
 */
export const startMockWorker = async (options?: {
  onUnhandledRequest?: 'bypass' | 'warn' | 'error';
  quiet?: boolean;
}) => {
  try {
    await worker.start({
      onUnhandledRequest: options?.onUnhandledRequest || 'bypass',
      quiet: options?.quiet || false,
    });

    if (!options?.quiet) {
      console.log('[MSW] Service Worker started successfully');
    }
  } catch (error) {
    console.error('[MSW] Failed to start Service Worker:', error);
  }
};

/**
 * Stop the mock service worker
 */
export const stopMockWorker = () => {
  worker.stop();
  console.log('[MSW] Service Worker stopped');
};

/**
 * Use additional handlers for specific scenarios
 *
 * @param additionalHandlers - Array of MSW handlers to use
 */
export const useMockHandlers = (...additionalHandlers: Parameters<typeof worker.use>) => {
  worker.use(...additionalHandlers);
};

/**
 * Reset all handlers to their initial state
 */
export const resetMockWorker = () => {
  worker.resetHandlers();
};

// Export the worker instance as default for convenience
export default worker;