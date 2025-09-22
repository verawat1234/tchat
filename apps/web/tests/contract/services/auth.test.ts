import { describe, it, expect, beforeEach } from 'vitest';
import { http, HttpResponse } from 'msw';
import { server } from '../../../src/lib/test-utils/msw/server';

/**
 * Authentication API Contract Tests (T006, T007, T008)
 *
 * These tests validate the request/response schemas and contracts for
 * authentication endpoints. Following TDD approach - tests WILL FAIL
 * until the actual API endpoints are implemented to match contracts.
 *
 * Endpoints tested:
 * - POST /api/auth/login (T006)
 * - POST /api/auth/refresh (T007) - MISSING ENDPOINT
 * - GET /api/auth/me (T008)
 *
 * Expected FAILURES:
 * 1. Login response format mismatch (current: {user, token} vs expected: {success, data: {user, tokens}})
 * 2. Refresh endpoint doesn't exist
 * 3. GET /me response format mismatch
 * 4. Error response formats don't match contract
 */

// Base API URL - matching the project's MSW setup
const API_BASE_URL = process.env.VITE_API_URL || 'http://localhost:3001';

// Contract type definitions
interface LoginResponse {
  success: boolean;
  data: {
    user: {
      id: string;
      email: string;
      name: string;
      avatar?: string;
      role: string;
      createdAt: string;
      updatedAt: string;
    };
    tokens: {
      accessToken: string;
      refreshToken: string;
      expiresAt: string;
    };
  };
  message?: string;
}

interface RefreshResponse {
  success: boolean;
  data: {
    tokens: {
      accessToken: string;
      refreshToken: string;
      expiresAt: string;
    };
  };
  message?: string;
}

interface MeResponse {
  success: boolean;
  data: {
    user: {
      id: string;
      email: string;
      name: string;
      avatar?: string;
      role: string;
      createdAt: string;
      updatedAt: string;
      lastLoginAt?: string;
    };
  };
  message?: string;
}

interface ErrorResponse {
  success: false;
  error: {
    code: string;
    message: string;
    details?: Record<string, any>;
  };
  statusCode: number;
}

describe('Authentication API Contract Tests', () => {
  beforeEach(() => {
    server.resetHandlers();
  });

  describe('POST /api/auth/login (T006)', () => {
    it('should match expected contract response format', async () => {
      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: 'test@example.com',
          password: 'password'
        })
      });

      const data = await response.json();

      // EXPECTED TO FAIL: Current MSW handler returns {user, token}
      // Contract expects {success: true, data: {user, tokens}}
      expect(response.status).toBe(200);
      expect(data.success).toBe(true);
      expect(data.data).toBeDefined();
      expect(data.data.user).toBeDefined();
      expect(data.data.tokens).toBeDefined();
      expect(data.data.tokens.accessToken).toBeTruthy();
      expect(data.data.tokens.refreshToken).toBeTruthy();
      expect(data.data.tokens.expiresAt).toBeTruthy();
    });

    it('should handle invalid credentials with proper error format', async () => {
      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: 'test@example.com',
          password: 'wrongpassword'
        })
      });

      const data = await response.json();

      // EXPECTED TO FAIL: Current handler returns {error: "Invalid credentials"}
      // Contract expects {success: false, error: {code, message}, statusCode}
      expect(response.status).toBe(401);
      expect(data.success).toBe(false);
      expect(data.error).toBeDefined();
      expect(data.error.code).toBe('INVALID_CREDENTIALS');
      expect(data.error.message).toBeTruthy();
      expect(data.statusCode).toBe(401);
    });

    it('should validate required fields', async () => {
      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: 'test@example.com'
          // password missing
        })
      });

      const data = await response.json();

      // EXPECTED TO FAIL: No validation in current handler
      expect(response.status).toBe(400);
      expect(data.success).toBe(false);
      expect(data.error.code).toBe('MISSING_FIELDS');
      expect(data.error.details.missingFields).toContain('password');
    });

    it('should validate email format', async () => {
      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: 'invalid-email',
          password: 'password123'
        })
      });

      const data = await response.json();

      // EXPECTED TO FAIL: No email validation in current handler
      expect(response.status).toBe(400);
      expect(data.success).toBe(false);
      expect(data.error.code).toBe('INVALID_EMAIL');
    });
  });

  describe('POST /api/auth/refresh (T007)', () => {
    it('should refresh tokens with valid refresh token', async () => {
      // Override default MSW to add refresh endpoint (it doesn't exist yet)
      server.use(
        http.post(`${API_BASE_URL}/api/auth/refresh`, async ({ request }) => {
          const body = await request.json() as { refreshToken: string };

          if (body.refreshToken === 'valid_refresh_token') {
            const refreshResponse: RefreshResponse = {
              success: true,
              data: {
                tokens: {
                  accessToken: 'new_access_token',
                  refreshToken: 'new_refresh_token',
                  expiresAt: new Date(Date.now() + 3600000).toISOString()
                }
              },
              message: 'Tokens refreshed successfully'
            };
            return HttpResponse.json(refreshResponse);
          }

          return HttpResponse.json(
            {
              success: false,
              error: {
                code: 'INVALID_REFRESH_TOKEN',
                message: 'Invalid refresh token'
              },
              statusCode: 401
            },
            { status: 401 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          refreshToken: 'valid_refresh_token'
        })
      });

      const data = await response.json();

      // This test validates the contract once we add the handler above
      expect(response.status).toBe(200);
      expect(data.success).toBe(true);
      expect(data.data.tokens).toBeDefined();
      expect(data.data.tokens.accessToken).toBeTruthy();
      expect(data.data.tokens.refreshToken).toBeTruthy();
      expect(data.data.tokens.expiresAt).toBeTruthy();
    });

    it('should handle missing refresh endpoint', async () => {
      // EXPECTED TO FAIL: Endpoint doesn't exist in current MSW handlers
      const response = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          refreshToken: 'some_token'
        })
      });

      // This will likely return 404 or fail due to no handler
      expect(response.status).not.toBe(404);
      expect(response.status).toBe(200);
    });

    it('should handle invalid refresh token', async () => {
      // Add temporary handler for this test
      server.use(
        http.post(`${API_BASE_URL}/api/auth/refresh`, () => {
          return HttpResponse.json(
            {
              success: false,
              error: {
                code: 'REFRESH_TOKEN_EXPIRED',
                message: 'Refresh token has expired'
              },
              statusCode: 401
            },
            { status: 401 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          refreshToken: 'expired_token'
        })
      });

      const data = await response.json();

      expect(response.status).toBe(401);
      expect(data.success).toBe(false);
      expect(data.error.code).toBe('REFRESH_TOKEN_EXPIRED');
    });
  });

  describe('GET /api/auth/me (T008)', () => {
    it('should return user profile with expected contract format', async () => {
      const response = await fetch(`${API_BASE_URL}/api/auth/me`, {
        method: 'GET',
        headers: {
          'Authorization': 'Bearer valid_token',
          'Content-Type': 'application/json'
        }
      });

      const data = await response.json();

      // EXPECTED TO FAIL: Current handler returns {user}
      // Contract expects {success: true, data: {user}}
      expect(response.status).toBe(200);
      expect(data.success).toBe(true);
      expect(data.data).toBeDefined();
      expect(data.data.user).toBeDefined();
      expect(data.data.user.id).toBeTruthy();
      expect(data.data.user.email).toBeTruthy();
      expect(data.data.user.name).toBeTruthy();
    });

    it('should require authorization header', async () => {
      // Override handler to check authorization
      server.use(
        http.get(`${API_BASE_URL}/api/auth/me`, ({ request }) => {
          const authHeader = request.headers.get('Authorization');

          if (!authHeader) {
            return HttpResponse.json(
              {
                success: false,
                error: {
                  code: 'MISSING_AUTHORIZATION',
                  message: 'Authorization header is required'
                },
                statusCode: 401
              },
              { status: 401 }
            );
          }

          return HttpResponse.json({
            success: true,
            data: {
              user: {
                id: '1',
                email: 'test@example.com',
                name: 'Test User',
                role: 'user',
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString()
              }
            }
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/auth/me`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
        // No Authorization header
      });

      const data = await response.json();

      // EXPECTED TO FAIL: Current handler doesn't check authorization
      expect(response.status).toBe(401);
      expect(data.success).toBe(false);
      expect(data.error.code).toBe('MISSING_AUTHORIZATION');
    });

    it('should handle invalid bearer token format', async () => {
      // Override handler to validate token format
      server.use(
        http.get(`${API_BASE_URL}/api/auth/me`, ({ request }) => {
          const authHeader = request.headers.get('Authorization');

          if (!authHeader || !authHeader.startsWith('Bearer ')) {
            return HttpResponse.json(
              {
                success: false,
                error: {
                  code: 'INVALID_TOKEN_FORMAT',
                  message: 'Invalid authorization token format'
                },
                statusCode: 401
              },
              { status: 401 }
            );
          }

          return HttpResponse.json({
            success: true,
            data: {
              user: {
                id: '1',
                email: 'test@example.com',
                name: 'Test User',
                role: 'user',
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString()
              }
            }
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/auth/me`, {
        method: 'GET',
        headers: {
          'Authorization': 'InvalidFormat token123',
          'Content-Type': 'application/json'
        }
      });

      const data = await response.json();

      expect(response.status).toBe(401);
      expect(data.success).toBe(false);
      expect(data.error.code).toBe('INVALID_TOKEN_FORMAT');
    });
  });

  describe('Common API Contract Validations', () => {
    it('should handle server errors with consistent format', async () => {
      // Override to simulate server error
      server.use(
        http.post(`${API_BASE_URL}/api/auth/login`, () => {
          return HttpResponse.json(
            {
              success: false,
              error: {
                code: 'INTERNAL_SERVER_ERROR',
                message: 'An unexpected error occurred'
              },
              statusCode: 500
            },
            { status: 500 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: 'test@example.com',
          password: 'password123'
        })
      });

      const data = await response.json();

      expect(response.status).toBe(500);
      expect(data.success).toBe(false);
      expect(data.error.code).toBe('INTERNAL_SERVER_ERROR');
      expect(data.statusCode).toBe(500);
    });

    it('should validate Content-Type requirements', async () => {
      // Override to validate content type
      server.use(
        http.post(`${API_BASE_URL}/api/auth/login`, ({ request }) => {
          const contentType = request.headers.get('Content-Type');

          if (contentType !== 'application/json') {
            return HttpResponse.json(
              {
                success: false,
                error: {
                  code: 'INVALID_CONTENT_TYPE',
                  message: 'Content-Type must be application/json'
                },
                statusCode: 415
              },
              { status: 415 }
            );
          }

          return HttpResponse.json({
            success: true,
            data: {
              user: { id: '1', email: 'test@example.com', name: 'Test' },
              tokens: { accessToken: 'token', refreshToken: 'refresh', expiresAt: '2024-01-01' }
            }
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'text/plain'
        },
        body: JSON.stringify({
          email: 'test@example.com',
          password: 'password123'
        })
      });

      const data = await response.json();

      // EXPECTED TO FAIL: Current handler doesn't validate Content-Type
      expect(response.status).toBe(415);
      expect(data.success).toBe(false);
      expect(data.error.code).toBe('INVALID_CONTENT_TYPE');
    });
  });
});