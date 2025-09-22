import { http, HttpResponse } from 'msw';
import type { RequestHandler } from 'msw';

// Base API URL - can be configured via environment variables
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:3001/api';

/**
 * Mock Data Factories
 *
 * These functions generate consistent mock data for testing
 */
export const createMockUser = (overrides: Partial<any> = {}) => ({
  id: '1',
  name: 'Test User',
  email: 'test@example.com',
  avatar: null,
  role: 'user',
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
  ...overrides,
});

export const createMockMessage = (overrides: Partial<any> = {}) => ({
  id: '1',
  content: 'Hello, this is a test message',
  authorId: '1',
  channelId: '1',
  timestamp: new Date().toISOString(),
  edited: false,
  reactions: [],
  ...overrides,
});

export const createMockChannel = (overrides: Partial<any> = {}) => ({
  id: '1',
  name: 'general',
  description: 'General discussion channel',
  type: 'public',
  createdAt: new Date().toISOString(),
  memberCount: 1,
  ...overrides,
});

export const createMockConversation = (overrides: Partial<any> = {}) => ({
  id: '1',
  participants: [createMockUser()],
  lastMessage: createMockMessage(),
  lastActivity: new Date().toISOString(),
  unreadCount: 0,
  ...overrides,
});

/**
 * Authentication Handlers
 */
const authHandlers: RequestHandler[] = [
  // POST /api/auth/login
  http.post(`${API_BASE}/auth/login`, async ({ request }) => {
    const body = await request.json() as { email: string; password: string };

    // Mock successful login
    if (body.email && body.password) {
      return HttpResponse.json({
        user: createMockUser({ email: body.email }),
        token: 'mock-jwt-token',
        refreshToken: 'mock-refresh-token',
      });
    }

    // Mock failed login
    return HttpResponse.json(
      { error: 'Invalid credentials' },
      { status: 401 }
    );
  }),

  // POST /api/auth/register
  http.post(`${API_BASE}/auth/register`, async ({ request }) => {
    const body = await request.json() as { email: string; password: string; name: string };

    return HttpResponse.json({
      user: createMockUser({
        email: body.email,
        name: body.name,
        id: Date.now().toString(),
      }),
      token: 'mock-jwt-token',
    });
  }),

  // POST /api/auth/logout
  http.post(`${API_BASE}/auth/logout`, () => {
    return HttpResponse.json({ success: true });
  }),

  // GET /api/auth/me
  http.get(`${API_BASE}/auth/me`, () => {
    return HttpResponse.json({ user: createMockUser() });
  }),
];

/**
 * User Management Handlers
 */
const userHandlers: RequestHandler[] = [
  // GET /api/users
  http.get(`${API_BASE}/users`, ({ request }) => {
    const url = new URL(request.url);
    const page = parseInt(url.searchParams.get('page') || '1');
    const limit = parseInt(url.searchParams.get('limit') || '10');

    const users = Array.from({ length: limit }, (_, i) =>
      createMockUser({
        id: (page * limit - limit + i + 1).toString(),
        name: `User ${page * limit - limit + i + 1}`,
        email: `user${page * limit - limit + i + 1}@example.com`,
      })
    );

    return HttpResponse.json({
      users,
      pagination: {
        page,
        limit,
        total: 100,
        totalPages: Math.ceil(100 / limit),
      },
    });
  }),

  // GET /api/users/:id
  http.get(`${API_BASE}/users/:id`, ({ params }) => {
    const { id } = params;

    return HttpResponse.json({
      user: createMockUser({
        id: id as string,
        name: `User ${id}`,
        email: `user${id}@example.com`,
      }),
    });
  }),
];

/**
 * Chat/Messaging Handlers
 */
const chatHandlers: RequestHandler[] = [
  // GET /api/channels
  http.get(`${API_BASE}/channels`, () => {
    const channels = [
      createMockChannel({ id: '1', name: 'general' }),
      createMockChannel({ id: '2', name: 'random' }),
      createMockChannel({ id: '3', name: 'tech-talk' }),
    ];

    return HttpResponse.json({ channels });
  }),

  // GET /api/channels/:id/messages
  http.get(`${API_BASE}/channels/:channelId/messages`, ({ params, request }) => {
    const { channelId } = params;
    const url = new URL(request.url);
    const before = url.searchParams.get('before');
    const limit = parseInt(url.searchParams.get('limit') || '50');

    const messages = Array.from({ length: limit }, (_, i) =>
      createMockMessage({
        id: (i + 1).toString(),
        content: `Message ${i + 1} in channel ${channelId}`,
        channelId: channelId as string,
        authorId: Math.random() > 0.5 ? '1' : '2',
        timestamp: new Date(Date.now() - i * 60000).toISOString(),
      })
    );

    return HttpResponse.json({
      messages: messages.reverse(),
      hasMore: true,
    });
  }),

  // POST /api/channels/:id/messages
  http.post(`${API_BASE}/channels/:channelId/messages`, async ({ params, request }) => {
    const { channelId } = params;
    const body = await request.json() as { content: string };

    const newMessage = createMockMessage({
      id: Date.now().toString(),
      content: body.content,
      channelId: channelId as string,
      authorId: '1',
      timestamp: new Date().toISOString(),
    });

    return HttpResponse.json({ message: newMessage }, { status: 201 });
  }),

  // GET /api/conversations
  http.get(`${API_BASE}/conversations`, () => {
    const conversations = [
      createMockConversation({ id: '1' }),
      createMockConversation({ id: '2', unreadCount: 3 }),
    ];

    return HttpResponse.json({ conversations });
  }),

  // GET /api/conversations/:id/messages
  http.get(`${API_BASE}/conversations/:conversationId/messages`, ({ params }) => {
    const { conversationId } = params;

    const messages = Array.from({ length: 20 }, (_, i) =>
      createMockMessage({
        id: (i + 1).toString(),
        content: `Direct message ${i + 1}`,
        authorId: i % 2 === 0 ? '1' : '2',
        timestamp: new Date(Date.now() - i * 60000).toISOString(),
      })
    );

    return HttpResponse.json({
      messages: messages.reverse(),
      hasMore: false,
    });
  }),
];

/**
 * File Upload Handlers
 */
const uploadHandlers: RequestHandler[] = [
  // POST /api/upload
  http.post(`${API_BASE}/upload`, async ({ request }) => {
    const formData = await request.formData();
    const file = formData.get('file') as File;

    if (!file) {
      return HttpResponse.json(
        { error: 'No file provided' },
        { status: 400 }
      );
    }

    // Simulate file upload
    return HttpResponse.json({
      id: Date.now().toString(),
      filename: file.name,
      size: file.size,
      url: `/uploads/${Date.now()}-${file.name}`,
      mimetype: file.type,
      uploadedAt: new Date().toISOString(),
    });
  }),
];

/**
 * Health Check Handler
 */
const healthHandlers: RequestHandler[] = [
  // GET /api/health
  http.get(`${API_BASE}/health`, () => {
    return HttpResponse.json({
      status: 'ok',
      timestamp: new Date().toISOString(),
      version: '1.0.0',
    });
  }),
];

/**
 * All request handlers combined
 */
export const handlers: RequestHandler[] = [
  ...authHandlers,
  ...userHandlers,
  ...chatHandlers,
  ...uploadHandlers,
  ...healthHandlers,
];

/**
 * Error response handlers for testing error scenarios
 */
export const errorHandlers = {
  // Simulate server errors (500)
  serverError: http.all(`${API_BASE}/*`, () => {
    return HttpResponse.json(
      { error: 'Internal Server Error' },
      { status: 500 }
    );
  }),

  // Simulate network errors
  networkError: http.all(`${API_BASE}/*`, () => {
    return HttpResponse.error();
  }),

  // Simulate unauthorized errors (401)
  unauthorized: http.all(`${API_BASE}/*`, () => {
    return HttpResponse.json(
      { error: 'Unauthorized' },
      { status: 401 }
    );
  }),

  // Simulate forbidden errors (403)
  forbidden: http.all(`${API_BASE}/*`, () => {
    return HttpResponse.json(
      { error: 'Forbidden' },
      { status: 403 }
    );
  }),

  // Simulate not found errors (404)
  notFound: http.all(`${API_BASE}/*`, () => {
    return HttpResponse.json(
      { error: 'Not Found' },
      { status: 404 }
    );
  }),

  // Simulate rate limiting (429)
  rateLimit: http.all(`${API_BASE}/*`, () => {
    return HttpResponse.json(
      { error: 'Too Many Requests' },
      { status: 429, headers: { 'Retry-After': '60' } }
    );
  }),
};

/**
 * Utility function to create delayed responses
 */
export const withDelay = (handler: RequestHandler, delay: number = 1000): RequestHandler => {
  return async (info) => {
    await new Promise(resolve => setTimeout(resolve, delay));
    return handler(info);
  };
};

/**
 * Export default handlers
 */
export default handlers;