import { describe, it, expect, beforeAll, afterAll, beforeEach } from 'vitest';
import { setupServer } from 'msw/node';
import { http, HttpResponse, delay } from 'msw';
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../../src/services/api';
import { authApi } from '../../src/services/auth';
import { usersApi } from '../../src/services/users';
import { messagesApi } from '../../src/services/messages';
import { chatsApi } from '../../src/services/chats';

/**
 * Performance Tests for API Response Times
 *
 * Ensures all API endpoints respond within acceptable time limits:
 * - Critical endpoints: <200ms (auth operations)
 * - Standard endpoints: <500ms (CRUD operations)
 * - Heavy endpoints: <1000ms (list operations with large datasets)
 */

const server = setupServer();

describe('API Performance Tests', () => {
  let store: ReturnType<typeof configureStore>;

  beforeAll(() => {
    server.listen({ onUnhandledRequest: 'error' });
  });

  afterAll(() => {
    server.close();
  });

  beforeEach(() => {
    server.resetHandlers();

    // Create fresh store for each test
    store = configureStore({
      reducer: {
        [api.reducerPath]: api.reducer,
      },
      middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(api.middleware),
    });
  });

  describe('Authentication Endpoints - Critical Performance (<200ms)', () => {
    it('should login within 200ms', async () => {
      server.use(
        http.post('/api/auth/login', async () => {
          // Simulate realistic auth processing time
          await delay(150);
          return HttpResponse.json({
            user: {
              id: '1',
              email: 'test@example.com',
              username: 'testuser',
              displayName: 'Test User',
              role: 'user',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
            tokens: {
              accessToken: 'access-token',
              refreshToken: 'refresh-token',
              expiresIn: 3600,
              tokenType: 'Bearer',
            },
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        authApi.endpoints.login.initiate({
          email: 'test@example.com',
          password: 'password',
        })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(responseTime).toBeLessThan(200);
    });

    it('should refresh token within 200ms', async () => {
      server.use(
        http.post('/api/auth/refresh', async () => {
          await delay(100);
          return HttpResponse.json({
            accessToken: 'new-access-token',
            refreshToken: 'new-refresh-token',
            expiresIn: 3600,
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        authApi.endpoints.refreshToken.initiate({
          refreshToken: 'refresh-token',
        })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(responseTime).toBeLessThan(200);
    });

    it('should get current user within 200ms', async () => {
      server.use(
        http.get('/api/auth/me', async () => {
          await delay(120);
          return HttpResponse.json({
            success: true,
            data: {
              id: '1',
              email: 'test@example.com',
              username: 'testuser',
              displayName: 'Test User',
              role: 'user',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        authApi.endpoints.getCurrentUser.initiate()
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(responseTime).toBeLessThan(200);
    });
  });

  describe('CRUD Operations - Standard Performance (<500ms)', () => {
    it('should get user by ID within 500ms', async () => {
      server.use(
        http.get('/api/users/:id', async () => {
          await delay(300);
          return HttpResponse.json({
            success: true,
            data: {
              id: '1',
              email: 'user@example.com',
              username: 'user1',
              displayName: 'User One',
              role: 'user',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        usersApi.endpoints.getUserById.initiate('1')
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(responseTime).toBeLessThan(500);
    });

    it('should update user within 500ms', async () => {
      server.use(
        http.patch('/api/users/:id', async ({ request }) => {
          await delay(350);
          const body = await request.json() as any;
          return HttpResponse.json({
            success: true,
            data: {
              id: '1',
              displayName: body.displayName,
              email: 'user@example.com',
              username: 'user1',
              role: 'user',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        usersApi.endpoints.updateUser.initiate({
          id: '1',
          data: { displayName: 'Updated Name' },
        })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(responseTime).toBeLessThan(500);
    });

    it('should send message within 500ms', async () => {
      server.use(
        http.post('/api/messages', async ({ request }) => {
          await delay(250);
          const body = await request.json() as any;
          return HttpResponse.json({
            success: true,
            data: {
              id: Date.now().toString(),
              chatId: body.chatId,
              userId: 'user1',
              content: body.content,
              type: body.type || 'text',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        messagesApi.endpoints.sendMessage.initiate({
          chatId: 'chat1',
          content: 'Hello World',
          type: 'text',
        })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(responseTime).toBeLessThan(500);
    });

    it('should create chat within 500ms', async () => {
      server.use(
        http.post('/api/chats', async ({ request }) => {
          await delay(400);
          const body = await request.json() as any;
          return HttpResponse.json({
            success: true,
            data: {
              id: Date.now().toString(),
              name: body.name,
              type: body.type,
              participants: body.participants,
              unreadCount: 0,
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        chatsApi.endpoints.createChat.initiate({
          name: 'Test Chat',
          type: 'group',
          participants: ['user1', 'user2'],
        })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(responseTime).toBeLessThan(500);
    });
  });

  describe('List Operations - Heavy Performance (<1000ms)', () => {
    it('should list users within 1000ms', async () => {
      server.use(
        http.get('/api/users', async () => {
          await delay(600);
          // Simulate large user list
          const users = Array.from({ length: 100 }, (_, i) => ({
            id: (i + 1).toString(),
            email: `user${i + 1}@example.com`,
            username: `user${i + 1}`,
            displayName: `User ${i + 1}`,
            role: 'user',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          }));

          return HttpResponse.json({
            success: true,
            data: users,
            meta: {
              pagination: {
                total: 1000,
                page: 1,
                limit: 100,
                totalPages: 10,
                hasNext: true,
                hasPrev: false,
              },
            },
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        usersApi.endpoints.listUsers.initiate({ page: 1, limit: 100 })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(result.data.data).toHaveLength(100);
      expect(responseTime).toBeLessThan(1000);
    });

    it('should list messages within 1000ms', async () => {
      server.use(
        http.get('/api/messages', async () => {
          await delay(800);
          // Simulate large message list
          const messages = Array.from({ length: 50 }, (_, i) => ({
            id: (i + 1).toString(),
            chatId: 'chat1',
            userId: 'user1',
            content: `Message ${i + 1}`,
            type: 'text',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          }));

          return HttpResponse.json({
            items: messages,
            pagination: {
              cursor: '50',
              nextCursor: '51',
              hasMore: true,
            },
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        messagesApi.endpoints.listMessages.initiate({
          chatId: 'chat1',
          limit: 50,
        })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(result.data.items).toHaveLength(50);
      expect(responseTime).toBeLessThan(1000);
    });

    it('should list chats within 1000ms', async () => {
      server.use(
        http.get('/api/chats', async () => {
          await delay(700);
          // Simulate moderate chat list
          const chats = Array.from({ length: 25 }, (_, i) => ({
            id: (i + 1).toString(),
            name: `Chat ${i + 1}`,
            type: 'group',
            participants: ['user1', 'user2'],
            unreadCount: Math.floor(Math.random() * 10),
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          }));

          return HttpResponse.json({
            success: true,
            data: chats,
          });
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        chatsApi.endpoints.listChats.initiate({ limit: 25 })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.data).toBeDefined();
      expect(result.data).toHaveLength(25);
      expect(responseTime).toBeLessThan(1000);
    });
  });

  describe('Concurrent Operations Performance', () => {
    it('should handle concurrent requests efficiently', async () => {
      server.use(
        http.get('/api/auth/me', async () => {
          await delay(100);
          return HttpResponse.json({
            success: true,
            data: {
              id: '1',
              email: 'test@example.com',
              username: 'testuser',
              displayName: 'Test User',
              role: 'user',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        }),
        http.get('/api/chats', async () => {
          await delay(150);
          return HttpResponse.json({
            success: true,
            data: [],
          });
        }),
        http.get('/api/users', async () => {
          await delay(200);
          return HttpResponse.json({
            success: true,
            data: [],
            meta: { pagination: {} },
          });
        })
      );

      const startTime = performance.now();

      // Execute concurrent requests
      const promises = [
        store.dispatch(authApi.endpoints.getCurrentUser.initiate()),
        store.dispatch(chatsApi.endpoints.listChats.initiate({ limit: 10 })),
        store.dispatch(usersApi.endpoints.listUsers.initiate({ limit: 10 })),
      ];

      const results = await Promise.all(promises);

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      // Should complete in roughly the time of the slowest request (200ms)
      // Plus some overhead for concurrent processing
      expect(responseTime).toBeLessThan(300);
      expect(results).toHaveLength(3);
      results.forEach(result => {
        expect(result.data).toBeDefined();
      });
    });
  });

  describe('Cache Performance', () => {
    it('should serve cached data quickly (<50ms)', async () => {
      server.use(
        http.get('/api/users/1', async () => {
          await delay(200);
          return HttpResponse.json({
            success: true,
            data: {
              id: '1',
              email: 'user@example.com',
              username: 'user1',
              displayName: 'User One',
              role: 'user',
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      // First request - should hit the server
      const firstResult = await store.dispatch(
        usersApi.endpoints.getUserById.initiate('1')
      );
      expect(firstResult.data).toBeDefined();

      // Second request - should hit cache
      const startTime = performance.now();

      const cachedResult = await store.dispatch(
        usersApi.endpoints.getUserById.initiate('1')
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(cachedResult.data).toBeDefined();
      expect(responseTime).toBeLessThan(50);
    });
  });

  describe('Error Response Performance', () => {
    it('should handle 404 errors quickly (<200ms)', async () => {
      server.use(
        http.get('/api/users/999', async () => {
          await delay(100);
          return HttpResponse.json(
            {
              success: false,
              error: {
                code: 'NOT_FOUND',
                message: 'User not found',
                timestamp: new Date().toISOString(),
              },
            },
            { status: 404 }
          );
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        usersApi.endpoints.getUserById.initiate('999')
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.error).toBeDefined();
      expect(responseTime).toBeLessThan(200);
    });

    it('should handle validation errors quickly (<200ms)', async () => {
      server.use(
        http.post('/api/auth/login', async () => {
          await delay(80);
          return HttpResponse.json(
            {
              success: false,
              error: {
                code: 'VALIDATION_ERROR',
                message: 'Invalid credentials',
                timestamp: new Date().toISOString(),
              },
            },
            { status: 422 }
          );
        })
      );

      const startTime = performance.now();

      const result = await store.dispatch(
        authApi.endpoints.login.initiate({
          email: 'invalid',
          password: 'wrong',
        })
      );

      const endTime = performance.now();
      const responseTime = endTime - startTime;

      expect(result.error).toBeDefined();
      expect(responseTime).toBeLessThan(200);
    });
  });
});