import { http, HttpResponse } from 'msw';
import { server } from './server';
import { errorHandlers } from './handlers';

/**
 * Utility functions for MSW testing scenarios
 */

/**
 * Mock a successful API response for a specific endpoint
 *
 * @param method - HTTP method
 * @param path - API endpoint path (relative to base URL)
 * @param response - Mock response data
 * @param status - HTTP status code (default: 200)
 */
export const mockApiResponse = (
  method: 'get' | 'post' | 'put' | 'patch' | 'delete',
  path: string,
  response: any,
  status: number = 200
) => {
  const baseUrl = process.env.VITE_API_URL || 'http://localhost:3001/api';
  const fullPath = path.startsWith('/') ? `${baseUrl}${path}` : `${baseUrl}/${path}`;

  const httpMethod = http[method];

  server.use(
    httpMethod(fullPath, () => {
      return HttpResponse.json(response, { status });
    })
  );
};

/**
 * Mock an API error response
 *
 * @param method - HTTP method
 * @param path - API endpoint path
 * @param error - Error response data
 * @param status - HTTP status code
 */
export const mockApiError = (
  method: 'get' | 'post' | 'put' | 'patch' | 'delete',
  path: string,
  error: any,
  status: number = 500
) => {
  mockApiResponse(method, path, error, status);
};

/**
 * Mock a network error for a specific endpoint
 *
 * @param method - HTTP method
 * @param path - API endpoint path
 */
export const mockNetworkError = (
  method: 'get' | 'post' | 'put' | 'patch' | 'delete',
  path: string
) => {
  const baseUrl = process.env.VITE_API_URL || 'http://localhost:3001/api';
  const fullPath = path.startsWith('/') ? `${baseUrl}${path}` : `${baseUrl}/${path}`;

  const httpMethod = http[method];

  server.use(
    httpMethod(fullPath, () => {
      return HttpResponse.error();
    })
  );
};

/**
 * Mock a delayed response for testing loading states
 *
 * @param method - HTTP method
 * @param path - API endpoint path
 * @param response - Mock response data
 * @param delay - Delay in milliseconds
 * @param status - HTTP status code
 */
export const mockDelayedResponse = (
  method: 'get' | 'post' | 'put' | 'patch' | 'delete',
  path: string,
  response: any,
  delay: number = 1000,
  status: number = 200
) => {
  const baseUrl = process.env.VITE_API_URL || 'http://localhost:3001/api';
  const fullPath = path.startsWith('/') ? `${baseUrl}${path}` : `${baseUrl}/${path}`;

  const httpMethod = http[method];

  server.use(
    httpMethod(fullPath, async () => {
      await new Promise(resolve => setTimeout(resolve, delay));
      return HttpResponse.json(response, { status });
    })
  );
};

/**
 * Enable global error handlers for testing error scenarios
 *
 * @param type - Type of error to simulate
 */
export const enableGlobalErrorHandlers = (
  type: 'serverError' | 'networkError' | 'unauthorized' | 'forbidden' | 'notFound' | 'rateLimit'
) => {
  server.use(errorHandlers[type]);
};

/**
 * Mock authentication states
 */
export const mockAuthStates = {
  /**
   * Mock authenticated user state
   */
  authenticated: (user?: any) => {
    const baseUrl = process.env.VITE_API_URL || 'http://localhost:3001/api';
    const mockUser = user || {
      id: '1',
      name: 'Test User',
      email: 'test@example.com',
      role: 'user',
    };

    server.use(
      http.get(`${baseUrl}/auth/me`, () => {
        return HttpResponse.json({ user: mockUser });
      })
    );
  },

  /**
   * Mock unauthenticated state
   */
  unauthenticated: () => {
    const baseUrl = process.env.VITE_API_URL || 'http://localhost:3001/api';

    server.use(
      http.get(`${baseUrl}/auth/me`, () => {
        return HttpResponse.json(
          { error: 'Unauthorized' },
          { status: 401 }
        );
      })
    );
  },

  /**
   * Mock expired token state
   */
  tokenExpired: () => {
    const baseUrl = process.env.VITE_API_URL || 'http://localhost:3001/api';

    server.use(
      http.get(`${baseUrl}/auth/me`, () => {
        return HttpResponse.json(
          { error: 'Token expired' },
          { status: 401 }
        );
      })
    );
  },
};

/**
 * Mock pagination responses
 *
 * @param endpoint - API endpoint
 * @param data - Array of data items
 * @param page - Current page number
 * @param limit - Items per page
 * @param total - Total number of items
 */
export const mockPaginatedResponse = (
  endpoint: string,
  data: any[],
  page: number = 1,
  limit: number = 10,
  total?: number
) => {
  const baseUrl = process.env.VITE_API_URL || 'http://localhost:3001/api';
  const fullPath = endpoint.startsWith('/') ? `${baseUrl}${endpoint}` : `${baseUrl}/${endpoint}`;

  const startIndex = (page - 1) * limit;
  const endIndex = startIndex + limit;
  const paginatedData = data.slice(startIndex, endIndex);
  const totalItems = total || data.length;
  const totalPages = Math.ceil(totalItems / limit);

  server.use(
    http.get(fullPath, ({ request }) => {
      const url = new URL(request.url);
      const requestedPage = parseInt(url.searchParams.get('page') || '1');
      const requestedLimit = parseInt(url.searchParams.get('limit') || '10');

      if (requestedPage !== page || requestedLimit !== limit) {
        const requestedStartIndex = (requestedPage - 1) * requestedLimit;
        const requestedEndIndex = requestedStartIndex + requestedLimit;
        const requestedData = data.slice(requestedStartIndex, requestedEndIndex);

        return HttpResponse.json({
          data: requestedData,
          pagination: {
            page: requestedPage,
            limit: requestedLimit,
            total: totalItems,
            totalPages: Math.ceil(totalItems / requestedLimit),
            hasNext: requestedPage < Math.ceil(totalItems / requestedLimit),
            hasPrev: requestedPage > 1,
          },
        });
      }

      return HttpResponse.json({
        data: paginatedData,
        pagination: {
          page,
          limit,
          total: totalItems,
          totalPages,
          hasNext: page < totalPages,
          hasPrev: page > 1,
        },
      });
    })
  );
};

/**
 * Reset all mock overrides and restore default handlers
 */
export const resetAllMocks = () => {
  server.resetHandlers();
};

/**
 * Helper to create mock handlers for CRUD operations
 *
 * @param resource - Resource name (e.g., 'users', 'posts')
 * @param mockData - Array of mock data items
 */
export const createCrudMocks = (resource: string, mockData: any[]) => {
  const baseUrl = process.env.VITE_API_URL || 'http://localhost:3001/api';

  return [
    // GET /resource
    http.get(`${baseUrl}/${resource}`, () => {
      return HttpResponse.json({
        data: mockData,
        total: mockData.length,
      });
    }),

    // GET /resource/:id
    http.get(`${baseUrl}/${resource}/:id`, ({ params }) => {
      const item = mockData.find(item => item.id === params.id);
      if (!item) {
        return HttpResponse.json(
          { error: `${resource} not found` },
          { status: 404 }
        );
      }
      return HttpResponse.json({ data: item });
    }),

    // POST /resource
    http.post(`${baseUrl}/${resource}`, async ({ request }) => {
      const body = await request.json();
      const newItem = {
        id: (mockData.length + 1).toString(),
        ...body,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      return HttpResponse.json({ data: newItem }, { status: 201 });
    }),

    // PUT /resource/:id
    http.put(`${baseUrl}/${resource}/:id`, async ({ params, request }) => {
      const body = await request.json();
      const item = mockData.find(item => item.id === params.id);
      if (!item) {
        return HttpResponse.json(
          { error: `${resource} not found` },
          { status: 404 }
        );
      }
      const updatedItem = {
        ...item,
        ...body,
        updatedAt: new Date().toISOString(),
      };
      return HttpResponse.json({ data: updatedItem });
    }),

    // DELETE /resource/:id
    http.delete(`${baseUrl}/${resource}/:id`, ({ params }) => {
      const item = mockData.find(item => item.id === params.id);
      if (!item) {
        return HttpResponse.json(
          { error: `${resource} not found` },
          { status: 404 }
        );
      }
      return HttpResponse.json({ success: true });
    }),
  ];
};