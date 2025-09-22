import { setupServer } from 'msw/node';
import { handlers } from './handlers';

// Setup MSW server with default handlers
export const server = setupServer(...handlers);

// Server lifecycle methods for tests
export const setupMockServer = () => {
  beforeAll(() => {
    server.listen({
      onUnhandledRequest: 'error',
    });
  });

  afterEach(() => {
    server.resetHandlers();
  });

  afterAll(() => {
    server.close();
  });
};

// Helper to override handlers for specific tests
export const mockApiResponse = (
  method: 'get' | 'post' | 'put' | 'patch' | 'delete',
  path: string,
  response: any,
  status = 200
) => {
  const { http, HttpResponse } = require('msw');
  const httpMethod = http[method];

  server.use(
    httpMethod(path, () => {
      return HttpResponse.json(response, { status });
    })
  );
};

// Helper to simulate network errors
export const mockNetworkError = (path: string) => {
  const { http, HttpResponse } = require('msw');

  server.use(
    http.get(path, () => {
      return HttpResponse.error();
    })
  );
};

// Helper to simulate delays
export const mockDelayedResponse = (
  method: 'get' | 'post' | 'put' | 'patch' | 'delete',
  path: string,
  response: any,
  delay = 1000
) => {
  const { http, HttpResponse } = require('msw');
  const httpMethod = http[method];

  server.use(
    httpMethod(path, async () => {
      await new Promise((resolve) => setTimeout(resolve, delay));
      return HttpResponse.json(response);
    })
  );
};