import { http, HttpResponse } from 'msw';

// Base API URL - can be overridden in tests
export const API_BASE_URL = process.env.VITE_API_URL || 'http://localhost:3001';

// Mock data factories
export const mockUser = (overrides = {}) => ({
  id: '1',
  name: 'Test User',
  email: 'test@example.com',
  avatar: '/avatars/default.png',
  role: 'user',
  createdAt: new Date().toISOString(),
  ...overrides,
});

export const mockPost = (overrides = {}) => ({
  id: '1',
  title: 'Test Post',
  content: 'This is a test post content',
  authorId: '1',
  author: mockUser(),
  likes: 0,
  comments: [],
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
  ...overrides,
});

export const mockProduct = (overrides = {}) => ({
  id: '1',
  name: 'Test Product',
  description: 'This is a test product',
  price: 99.99,
  image: '/products/default.png',
  inStock: true,
  category: 'electronics',
  rating: 4.5,
  reviews: 42,
  ...overrides,
});

export const mockMessage = (overrides = {}) => ({
  id: '1',
  text: 'Test message',
  senderId: '1',
  sender: mockUser(),
  conversationId: '1',
  timestamp: new Date().toISOString(),
  read: false,
  ...overrides,
});

export const mockNotification = (overrides = {}) => ({
  id: '1',
  type: 'info',
  title: 'Test Notification',
  message: 'This is a test notification',
  read: false,
  createdAt: new Date().toISOString(),
  ...overrides,
});

// Default handlers for common endpoints
export const handlers = [
  // Authentication
  http.post(`${API_BASE_URL}/api/auth/login`, async ({ request }) => {
    const body = await request.json() as { email: string; password: string };

    if (body.email === 'test@example.com' && body.password === 'password') {
      return HttpResponse.json({
        user: mockUser(),
        token: 'mock-jwt-token',
      });
    }

    return HttpResponse.json(
      { error: 'Invalid credentials' },
      { status: 401 }
    );
  }),

  http.post(`${API_BASE_URL}/api/auth/logout`, () => {
    return HttpResponse.json({ success: true });
  }),

  http.get(`${API_BASE_URL}/api/auth/me`, () => {
    return HttpResponse.json({ user: mockUser() });
  }),

  // Users
  http.get(`${API_BASE_URL}/api/users`, () => {
    return HttpResponse.json({
      users: [
        mockUser({ id: '1', name: 'User 1' }),
        mockUser({ id: '2', name: 'User 2' }),
        mockUser({ id: '3', name: 'User 3' }),
      ],
      total: 3,
    });
  }),

  http.get(`${API_BASE_URL}/api/users/:id`, ({ params }) => {
    return HttpResponse.json({
      user: mockUser({ id: params.id as string }),
    });
  }),

  // Posts
  http.get(`${API_BASE_URL}/api/posts`, () => {
    return HttpResponse.json({
      posts: [
        mockPost({ id: '1', title: 'Post 1' }),
        mockPost({ id: '2', title: 'Post 2' }),
        mockPost({ id: '3', title: 'Post 3' }),
      ],
      total: 3,
    });
  }),

  http.post(`${API_BASE_URL}/api/posts`, async ({ request }) => {
    const body = await request.json() as { title: string; content: string };
    return HttpResponse.json({
      post: mockPost({ ...body, id: Date.now().toString() }),
    });
  }),

  // Products
  http.get(`${API_BASE_URL}/api/products`, () => {
    return HttpResponse.json({
      products: [
        mockProduct({ id: '1', name: 'Product 1' }),
        mockProduct({ id: '2', name: 'Product 2' }),
        mockProduct({ id: '3', name: 'Product 3' }),
      ],
      total: 3,
    });
  }),

  http.get(`${API_BASE_URL}/api/products/:id`, ({ params }) => {
    return HttpResponse.json({
      product: mockProduct({ id: params.id as string }),
    });
  }),

  // Cart
  http.get(`${API_BASE_URL}/api/cart`, () => {
    return HttpResponse.json({
      items: [
        {
          product: mockProduct({ id: '1' }),
          quantity: 2,
        },
      ],
      total: 199.98,
    });
  }),

  http.post(`${API_BASE_URL}/api/cart/add`, async ({ request }) => {
    const body = await request.json() as { productId: string; quantity: number };
    return HttpResponse.json({
      success: true,
      item: {
        product: mockProduct({ id: body.productId }),
        quantity: body.quantity,
      },
    });
  }),

  // Messages/Chat
  http.get(`${API_BASE_URL}/api/messages`, () => {
    return HttpResponse.json({
      messages: [
        mockMessage({ id: '1', text: 'Hello' }),
        mockMessage({ id: '2', text: 'Hi there!' }),
        mockMessage({ id: '3', text: 'How are you?' }),
      ],
      total: 3,
    });
  }),

  http.post(`${API_BASE_URL}/api/messages`, async ({ request }) => {
    const body = await request.json() as { text: string; conversationId: string };
    return HttpResponse.json({
      message: mockMessage({
        ...body,
        id: Date.now().toString(),
      }),
    });
  }),

  // Notifications
  http.get(`${API_BASE_URL}/api/notifications`, () => {
    return HttpResponse.json({
      notifications: [
        mockNotification({ id: '1', type: 'info' }),
        mockNotification({ id: '2', type: 'warning' }),
        mockNotification({ id: '3', type: 'success' }),
      ],
      unreadCount: 2,
    });
  }),

  http.patch(`${API_BASE_URL}/api/notifications/:id/read`, ({ params }) => {
    return HttpResponse.json({
      notification: mockNotification({
        id: params.id as string,
        read: true,
      }),
    });
  }),

  // File upload
  http.post(`${API_BASE_URL}/api/upload`, async ({ request }) => {
    const formData = await request.formData();
    const file = formData.get('file') as File;

    if (file) {
      return HttpResponse.json({
        url: `/uploads/${file.name}`,
        filename: file.name,
        size: file.size,
      });
    }

    return HttpResponse.json(
      { error: 'No file provided' },
      { status: 400 }
    );
  }),

  // Search
  http.get(`${API_BASE_URL}/api/search`, ({ request }) => {
    const url = new URL(request.url);
    const query = url.searchParams.get('q');

    return HttpResponse.json({
      results: [
        { type: 'user', item: mockUser({ name: `User matching ${query}` }) },
        { type: 'post', item: mockPost({ title: `Post about ${query}` }) },
        { type: 'product', item: mockProduct({ name: `Product: ${query}` }) },
      ],
      query,
    });
  }),

  // Health check
  http.get(`${API_BASE_URL}/api/health`, () => {
    return HttpResponse.json({
      status: 'ok',
      timestamp: new Date().toISOString(),
    });
  }),
];

// Error response handlers for testing error states
export const errorHandlers = {
  serverError: http.get(`${API_BASE_URL}/*`, () => {
    return HttpResponse.json(
      { error: 'Internal Server Error' },
      { status: 500 }
    );
  }),

  networkError: http.get(`${API_BASE_URL}/*`, () => {
    return HttpResponse.error();
  }),

  unauthorized: http.get(`${API_BASE_URL}/*`, () => {
    return HttpResponse.json(
      { error: 'Unauthorized' },
      { status: 401 }
    );
  }),

  notFound: http.get(`${API_BASE_URL}/*`, () => {
    return HttpResponse.json(
      { error: 'Not Found' },
      { status: 404 }
    );
  }),

  rateLimit: http.get(`${API_BASE_URL}/*`, () => {
    return HttpResponse.json(
      { error: 'Too Many Requests' },
      { status: 429, headers: { 'Retry-After': '60' } }
    );
  }),
};

// Delay utilities for testing loading states
export const withDelay = (handler: typeof http.get, delay = 1000) => {
  return http.get(handler.info.path, async (info) => {
    await new Promise((resolve) => setTimeout(resolve, delay));
    return handler(info);
  });
};

// Custom handler factory for specific test scenarios
export const createCustomHandler = (
  method: 'get' | 'post' | 'put' | 'patch' | 'delete',
  path: string,
  response: any,
  status = 200
) => {
  const httpMethod = http[method];
  return httpMethod(`${API_BASE_URL}${path}`, () => {
    return HttpResponse.json(response, { status });
  });
};