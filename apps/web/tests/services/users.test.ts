/**
 * User Service Contract Tests
 * Tests for GET /api/users, GET /api/users/:id, and PATCH /api/users/:id endpoints
 *
 * These tests validate request/response schemas according to the OpenAPI contract
 * Following TDD approach - tests will fail initially until endpoints are implemented
 *
 * EXPECTED FAILURES (TDD Approach):
 * 1. Schema validation failures - MSW handlers return mock data that doesn't match User schema
 * 2. Missing PATCH endpoint handler - No MSW handler for PATCH /api/users/:id
 * 3. Incomplete user objects - Mock users missing required fields like settings, profile, preferences
 * 4. Field type mismatches - Mock data types don't match schema expectations
 *
 * These failures are intentional and demonstrate that:
 * - Tests are written before implementation
 * - Schema validation is strict and comprehensive
 * - All endpoint contracts are thoroughly tested
 * - Error scenarios are properly handled
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { server } from '@/lib/test-utils/msw/server';
import { http, HttpResponse } from 'msw';
import type { User, UserProfile, UserSettings, UserPreferences } from '@/schema/schema';

// API Base URL
const API_BASE_URL = process.env.VITE_API_URL || 'http://localhost:3001';

// Type definitions for request/response validation
interface UsersListResponse {
  users: User[];
  total: number;
  page?: number;
  limit?: number;
  hasMore?: boolean;
}

interface UserResponse {
  user: User;
}

interface UserUpdateRequest {
  name?: string;
  email?: string;
  avatar?: string;
  profile?: Partial<UserProfile>;
  settings?: Partial<UserSettings>;
  preferences?: Partial<UserPreferences>;
}

interface ErrorResponse {
  error: string;
  code?: string;
  details?: Record<string, any>;
}

// Schema validation helpers
const validateUser = (user: any): user is User => {
  return (
    typeof user.id === 'string' &&
    typeof user.name === 'string' &&
    (typeof user.email === 'string' || user.email === undefined) &&
    (typeof user.phone === 'string' || user.phone === undefined) &&
    (typeof user.avatar === 'string' || user.avatar === undefined) &&
    typeof user.country === 'string' &&
    typeof user.locale === 'string' &&
    typeof user.kycTier === 'number' &&
    [1, 2, 3].includes(user.kycTier) &&
    typeof user.status === 'string' &&
    ['online', 'offline', 'away', 'busy'].includes(user.status) &&
    typeof user.isVerified === 'boolean' &&
    typeof user.settings === 'object' &&
    typeof user.profile === 'object' &&
    typeof user.preferences === 'object' &&
    typeof user.createdAt === 'string' &&
    typeof user.updatedAt === 'string'
  );
};

const validateUsersListResponse = (response: any): response is UsersListResponse => {
  return (
    Array.isArray(response.users) &&
    response.users.every(validateUser) &&
    typeof response.total === 'number' &&
    response.total >= 0
  );
};

const validateUserResponse = (response: any): response is UserResponse => {
  return (
    typeof response === 'object' &&
    response.user &&
    validateUser(response.user)
  );
};

const validateErrorResponse = (response: any): response is ErrorResponse => {
  return (
    typeof response === 'object' &&
    typeof response.error === 'string' &&
    response.error.length > 0
  );
};

describe('User Service Contract Tests', () => {
  beforeEach(() => {
    // Reset MSW handlers before each test
    server.resetHandlers();
  });

  describe('GET /api/users', () => {
    it('should return a list of users with valid schema', async () => {
      // This test will initially fail until the endpoint is properly implemented
      const response = await fetch(`${API_BASE_URL}/api/users`);
      const data = await response.json();

      // Contract validation
      expect(response.status).toBe(200);
      expect(response.headers.get('content-type')).toContain('application/json');

      // Schema validation
      expect(validateUsersListResponse(data)).toBe(true);
      expect(data.users).toBeDefined();
      expect(Array.isArray(data.users)).toBe(true);
      expect(typeof data.total).toBe('number');

      // Each user should conform to User schema
      data.users.forEach((user: any) => {
        expect(validateUser(user)).toBe(true);

        // Required fields validation
        expect(user.id).toBeDefined();
        expect(user.name).toBeDefined();
        expect(user.country).toBeDefined();
        expect(user.locale).toBeDefined();
        expect(user.kycTier).toBeDefined();
        expect(user.status).toBeDefined();
        expect(user.isVerified).toBeDefined();
        expect(user.settings).toBeDefined();
        expect(user.profile).toBeDefined();
        expect(user.preferences).toBeDefined();
        expect(user.createdAt).toBeDefined();
        expect(user.updatedAt).toBeDefined();
      });
    });

    it('should handle query parameters for pagination', async () => {
      const response = await fetch(`${API_BASE_URL}/api/users?page=1&limit=10`);
      const data = await response.json();

      expect(response.status).toBe(200);
      expect(validateUsersListResponse(data)).toBe(true);

      // Should respect pagination parameters
      if (data.page !== undefined) {
        expect(typeof data.page).toBe('number');
      }
      if (data.limit !== undefined) {
        expect(typeof data.limit).toBe('number');
      }
      if (data.hasMore !== undefined) {
        expect(typeof data.hasMore).toBe('boolean');
      }
    });

    it('should handle search query parameter', async () => {
      const response = await fetch(`${API_BASE_URL}/api/users?search=test`);
      const data = await response.json();

      expect(response.status).toBe(200);
      expect(validateUsersListResponse(data)).toBe(true);
    });

    it('should return 500 on server error', async () => {
      // Mock server error
      server.use(
        http.get(`${API_BASE_URL}/api/users`, () => {
          return HttpResponse.json(
            { error: 'Internal Server Error' },
            { status: 500 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users`);
      const data = await response.json();

      expect(response.status).toBe(500);
      expect(validateErrorResponse(data)).toBe(true);
    });

    it('should return 401 when unauthorized', async () => {
      // Mock unauthorized error
      server.use(
        http.get(`${API_BASE_URL}/api/users`, () => {
          return HttpResponse.json(
            { error: 'Unauthorized' },
            { status: 401 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users`);
      const data = await response.json();

      expect(response.status).toBe(401);
      expect(validateErrorResponse(data)).toBe(true);
    });
  });

  describe('GET /api/users/:id', () => {
    const userId = '123e4567-e89b-12d3-a456-426614174000';

    it('should return a single user with valid schema', async () => {
      // This test will initially fail until the endpoint is properly implemented
      const response = await fetch(`${API_BASE_URL}/api/users/${userId}`);
      const data = await response.json();

      // Contract validation
      expect(response.status).toBe(200);
      expect(response.headers.get('content-type')).toContain('application/json');

      // Schema validation
      expect(validateUserResponse(data)).toBe(true);
      expect(data.user).toBeDefined();
      expect(validateUser(data.user)).toBe(true);

      // User ID should match requested ID
      expect(data.user.id).toBe(userId);

      // Required fields validation
      expect(data.user.name).toBeDefined();
      expect(data.user.country).toBeDefined();
      expect(data.user.locale).toBeDefined();
      expect(data.user.kycTier).toBeDefined();
      expect(data.user.status).toBeDefined();
      expect(data.user.isVerified).toBeDefined();
      expect(data.user.settings).toBeDefined();
      expect(data.user.profile).toBeDefined();
      expect(data.user.preferences).toBeDefined();
      expect(data.user.createdAt).toBeDefined();
      expect(data.user.updatedAt).toBeDefined();
    });

    it('should return 404 for non-existent user', async () => {
      // Mock not found error
      server.use(
        http.get(`${API_BASE_URL}/api/users/:id`, ({ params }) => {
          if (params.id === 'non-existent-id') {
            return HttpResponse.json(
              { error: 'User not found' },
              { status: 404 }
            );
          }
          return HttpResponse.json({ user: { id: params.id } });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users/non-existent-id`);
      const data = await response.json();

      expect(response.status).toBe(404);
      expect(validateErrorResponse(data)).toBe(true);
      expect(data.error).toBe('User not found');
    });

    it('should return 400 for invalid user ID format', async () => {
      // Mock validation error
      server.use(
        http.get(`${API_BASE_URL}/api/users/:id`, ({ params }) => {
          if (params.id === 'invalid-uuid') {
            return HttpResponse.json(
              { error: 'Invalid user ID format' },
              { status: 400 }
            );
          }
          return HttpResponse.json({ user: { id: params.id } });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users/invalid-uuid`);
      const data = await response.json();

      expect(response.status).toBe(400);
      expect(validateErrorResponse(data)).toBe(true);
    });

    it('should return 401 when unauthorized', async () => {
      // Mock unauthorized error
      server.use(
        http.get(`${API_BASE_URL}/api/users/:id`, () => {
          return HttpResponse.json(
            { error: 'Unauthorized' },
            { status: 401 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users/${userId}`);
      const data = await response.json();

      expect(response.status).toBe(401);
      expect(validateErrorResponse(data)).toBe(true);
    });
  });

  describe('PATCH /api/users/:id', () => {
    const userId = '123e4567-e89b-12d3-a456-426614174000';

    it('should update user with valid data and return updated user', async () => {
      const updateData: UserUpdateRequest = {
        name: 'Updated Name',
        email: 'updated@example.com',
        profile: {
          displayName: 'Updated Display Name',
          bio: 'Updated bio',
        },
        settings: {
          appearance: {
            theme: 'dark',
            language: 'en-US',
            fontSize: 'medium',
          },
        },
      };

      // Mock successful update
      server.use(
        http.patch(`${API_BASE_URL}/api/users/:id`, async ({ params, request }) => {
          const body = await request.json() as UserUpdateRequest;

          return HttpResponse.json({
            user: {
              id: params.id,
              name: body.name || 'Test User',
              email: body.email || 'test@example.com',
              country: 'TH',
              locale: 'th-TH',
              kycTier: 1,
              status: 'online',
              isVerified: false,
              settings: {
                privacy: {
                  profileVisibility: 'public',
                  phoneVisibility: 'friends',
                  lastSeenVisibility: 'everyone',
                  readReceiptsEnabled: true,
                  onlineStatusVisible: true,
                },
                notifications: {
                  pushEnabled: true,
                  emailEnabled: true,
                  messageNotifications: true,
                  postNotifications: true,
                  eventNotifications: true,
                  paymentNotifications: true,
                  mutedChats: [],
                  mutedUsers: [],
                },
                appearance: body.settings?.appearance || {
                  theme: 'light',
                  language: 'en-US',
                  fontSize: 'medium',
                },
                security: {
                  twoFactorEnabled: false,
                  biometricEnabled: false,
                  autoLockTimeout: 30,
                  trustedDevices: [],
                },
              },
              profile: {
                displayName: body.profile?.displayName,
                bio: body.profile?.bio,
                interests: [],
                languages: ['en-US'],
                timezone: 'Asia/Bangkok',
              },
              preferences: {
                defaultCurrency: 'THB',
                eventCategories: [],
                productCategories: [],
                contentLanguages: ['en-US'],
                contentFilters: [],
              },
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users/${userId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updateData),
      });

      const data = await response.json();

      // Contract validation
      expect(response.status).toBe(200);
      expect(response.headers.get('content-type')).toContain('application/json');

      // Schema validation
      expect(validateUserResponse(data)).toBe(true);
      expect(data.user).toBeDefined();
      expect(validateUser(data.user)).toBe(true);

      // Updated fields validation
      expect(data.user.id).toBe(userId);
      expect(data.user.name).toBe(updateData.name);
      expect(data.user.email).toBe(updateData.email);
      if (updateData.profile?.displayName) {
        expect(data.user.profile.displayName).toBe(updateData.profile.displayName);
      }
      if (updateData.profile?.bio) {
        expect(data.user.profile.bio).toBe(updateData.profile.bio);
      }
    });

    it('should validate request content-type', async () => {
      const response = await fetch(`${API_BASE_URL}/api/users/${userId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'text/plain',
        },
        body: 'invalid data',
      });

      // Should return error for invalid content-type
      expect(response.status).toBe(400);
    });

    it('should validate request body schema', async () => {
      // Mock validation error
      server.use(
        http.patch(`${API_BASE_URL}/api/users/:id`, () => {
          return HttpResponse.json(
            {
              error: 'Invalid request data',
              details: {
                name: 'Name must be a string',
                email: 'Invalid email format',
              }
            },
            { status: 400 }
          );
        })
      );

      const invalidData = {
        name: 123, // Invalid type
        email: 'invalid-email', // Invalid format
        kycTier: 5, // Invalid value
      };

      const response = await fetch(`${API_BASE_URL}/api/users/${userId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(invalidData),
      });

      const data = await response.json();

      expect(response.status).toBe(400);
      expect(validateErrorResponse(data)).toBe(true);
      expect(data.details).toBeDefined();
    });

    it('should return 404 for non-existent user', async () => {
      // Mock not found error
      server.use(
        http.patch(`${API_BASE_URL}/api/users/:id`, ({ params }) => {
          if (params.id === 'non-existent-id') {
            return HttpResponse.json(
              { error: 'User not found' },
              { status: 404 }
            );
          }
          return HttpResponse.json({ user: { id: params.id } });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users/non-existent-id`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: 'Test' }),
      });

      const data = await response.json();

      expect(response.status).toBe(404);
      expect(validateErrorResponse(data)).toBe(true);
      expect(data.error).toBe('User not found');
    });

    it('should return 401 when unauthorized', async () => {
      // Mock unauthorized error
      server.use(
        http.patch(`${API_BASE_URL}/api/users/:id`, () => {
          return HttpResponse.json(
            { error: 'Unauthorized' },
            { status: 401 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users/${userId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: 'Test' }),
      });

      const data = await response.json();

      expect(response.status).toBe(401);
      expect(validateErrorResponse(data)).toBe(true);
    });

    it('should return 403 when forbidden (user trying to update another user)', async () => {
      // Mock forbidden error
      server.use(
        http.patch(`${API_BASE_URL}/api/users/:id`, () => {
          return HttpResponse.json(
            { error: 'Forbidden: Cannot update other users' },
            { status: 403 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users/${userId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: 'Test' }),
      });

      const data = await response.json();

      expect(response.status).toBe(403);
      expect(validateErrorResponse(data)).toBe(true);
    });

    it('should handle partial updates correctly', async () => {
      const partialUpdate: UserUpdateRequest = {
        name: 'Partially Updated Name',
      };

      // Mock successful partial update
      server.use(
        http.patch(`${API_BASE_URL}/api/users/:id`, async ({ params, request }) => {
          const body = await request.json() as UserUpdateRequest;

          return HttpResponse.json({
            user: {
              id: params.id,
              name: body.name || 'Original Name',
              email: 'original@example.com', // Should remain unchanged
              country: 'TH',
              locale: 'th-TH',
              kycTier: 1,
              status: 'online',
              isVerified: false,
              settings: {
                privacy: {
                  profileVisibility: 'public',
                  phoneVisibility: 'friends',
                  lastSeenVisibility: 'everyone',
                  readReceiptsEnabled: true,
                  onlineStatusVisible: true,
                },
                notifications: {
                  pushEnabled: true,
                  emailEnabled: true,
                  messageNotifications: true,
                  postNotifications: true,
                  eventNotifications: true,
                  paymentNotifications: true,
                  mutedChats: [],
                  mutedUsers: [],
                },
                appearance: {
                  theme: 'light',
                  language: 'en-US',
                  fontSize: 'medium',
                },
                security: {
                  twoFactorEnabled: false,
                  biometricEnabled: false,
                  autoLockTimeout: 30,
                  trustedDevices: [],
                },
              },
              profile: {
                interests: [],
                languages: ['en-US'],
                timezone: 'Asia/Bangkok',
              },
              preferences: {
                defaultCurrency: 'THB',
                eventCategories: [],
                productCategories: [],
                contentLanguages: ['en-US'],
                contentFilters: [],
              },
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/users/${userId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(partialUpdate),
      });

      const data = await response.json();

      expect(response.status).toBe(200);
      expect(validateUserResponse(data)).toBe(true);
      expect(data.user.name).toBe(partialUpdate.name);
      expect(data.user.email).toBe('original@example.com'); // Should remain unchanged
    });
  });

  describe('Cross-endpoint consistency', () => {
    const userId = '123e4567-e89b-12d3-a456-426614174000';

    it('should maintain data consistency between GET /api/users and GET /api/users/:id', async () => {
      // Get all users
      const usersResponse = await fetch(`${API_BASE_URL}/api/users`);
      const usersData = await usersResponse.json();

      expect(validateUsersListResponse(usersData)).toBe(true);

      // Get first user individually
      const firstUser = usersData.users[0];
      if (firstUser) {
        const userResponse = await fetch(`${API_BASE_URL}/api/users/${firstUser.id}`);
        const userData = await userResponse.json();

        expect(validateUserResponse(userData)).toBe(true);

        // Compare user objects for consistency
        expect(userData.user.id).toBe(firstUser.id);
        expect(userData.user.name).toBe(firstUser.name);
        expect(userData.user.email).toBe(firstUser.email);
        expect(userData.user.status).toBe(firstUser.status);
      }
    });

    it('should maintain data consistency after PATCH operation', async () => {
      const updateData: UserUpdateRequest = {
        name: 'Consistency Test User',
      };

      // Update user
      const updateResponse = await fetch(`${API_BASE_URL}/api/users/${userId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updateData),
      });

      const updateResult = await updateResponse.json();
      expect(validateUserResponse(updateResult)).toBe(true);

      // Fetch user again to verify consistency
      const getResponse = await fetch(`${API_BASE_URL}/api/users/${userId}`);
      const getResult = await getResponse.json();

      expect(validateUserResponse(getResult)).toBe(true);
      expect(getResult.user.name).toBe(updateData.name);
      expect(getResult.user.updatedAt).toBe(updateResult.user.updatedAt);
    });
  });
});