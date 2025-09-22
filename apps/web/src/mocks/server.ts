import { setupServer } from 'msw/node';
import { handlers } from './handlers';

/**
 * MSW Server Configuration for Testing
 *
 * This server is used for mocking API endpoints during development and testing.
 * It intercepts HTTP requests and returns mock responses based on the defined handlers.
 */
export const server = setupServer(...handlers);

/**
 * Start the mock server with custom configuration
 *
 * @param options - Configuration options for the server
 */
export const startMockServer = (options?: {
  onUnhandledRequest?: 'bypass' | 'warn' | 'error';
}) => {
  server.listen({
    onUnhandledRequest: options?.onUnhandledRequest || 'warn',
  });
};

/**
 * Stop the mock server and clean up
 */
export const stopMockServer = () => {
  server.close();
};

/**
 * Reset all handlers to their initial state
 * Use this between tests to ensure clean state
 */
export const resetMockServer = () => {
  server.resetHandlers();
};

/**
 * Use additional handlers for specific test scenarios
 *
 * @param additionalHandlers - Array of MSW handlers to use
 */
export const useMockHandlers = (...additionalHandlers: Parameters<typeof server.use>) => {
  server.use(...additionalHandlers);
};

// Export the server instance as default for convenience
export default server;