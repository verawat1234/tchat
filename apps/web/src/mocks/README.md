# Mock Service Worker (MSW) Configuration

This directory contains the Mock Service Worker (MSW) configuration for the Tchat application. MSW allows us to intercept network requests and return mock responses during development and testing.

## Files Structure

```
src/mocks/
├── README.md          # This documentation
├── index.ts           # Main exports and auto-setup
├── server.ts          # Node.js server setup for testing
├── browser.ts         # Browser worker setup for development
├── handlers.ts        # Request handlers and mock data
└── utils.ts           # Testing utilities and helpers
```

## Quick Start

### For Testing (Vitest/Jest)

The MSW server is automatically configured in the test setup. No additional setup is required.

```typescript
// Already configured in src/test-setup.ts
import { server } from '@/mocks/server';
```

### For Development (Browser)

To enable mocking in development mode:

```typescript
import { startMockWorker } from '@/mocks';

// Start the service worker
if (process.env.NODE_ENV === 'development') {
  startMockWorker();
}
```

### Auto Setup

Use the auto-setup function for easy initialization:

```typescript
import { setupMocks } from '@/mocks';

// Auto-detect environment and setup accordingly
setupMocks('auto');

// Or specify the environment
setupMocks('development'); // or 'test'
```

## Available Mock Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `POST /api/auth/logout` - User logout
- `GET /api/auth/me` - Get current user

### Users
- `GET /api/users` - List users (with pagination)
- `GET /api/users/:id` - Get user by ID

### Chat/Messaging
- `GET /api/channels` - List channels
- `GET /api/channels/:id/messages` - Get channel messages
- `POST /api/channels/:id/messages` - Send message to channel
- `GET /api/conversations` - List conversations
- `GET /api/conversations/:id/messages` - Get conversation messages

### File Upload
- `POST /api/upload` - Upload file

### Health Check
- `GET /api/health` - API health status

## Testing Utilities

### Mock API Responses

```typescript
import { mockApiResponse } from '@/mocks/utils';

// Mock a successful response
mockApiResponse('get', '/users', { users: [] });

// Mock an error response
mockApiError('post', '/users', { error: 'Validation failed' }, 400);
```

### Mock Authentication States

```typescript
import { mockAuthStates } from '@/mocks/utils';

// Mock authenticated user
mockAuthStates.authenticated({ id: '1', name: 'Test User' });

// Mock unauthenticated state
mockAuthStates.unauthenticated();

// Mock expired token
mockAuthStates.tokenExpired();
```

### Mock Network Conditions

```typescript
import { mockNetworkError, mockDelayedResponse } from '@/mocks/utils';

// Simulate network error
mockNetworkError('get', '/users');

// Simulate slow response
mockDelayedResponse('get', '/users', { users: [] }, 2000);
```

### Mock Pagination

```typescript
import { mockPaginatedResponse } from '@/mocks/utils';

const users = [/* array of user objects */];
mockPaginatedResponse('/users', users, 1, 10, 100);
```

## Mock Data Factories

Use these functions to create consistent mock data:

```typescript
import {
  createMockUser,
  createMockMessage,
  createMockChannel,
  createMockConversation
} from '@/mocks';

const user = createMockUser({ name: 'Custom Name' });
const message = createMockMessage({ content: 'Hello World' });
const channel = createMockChannel({ name: 'general' });
const conversation = createMockConversation({ unreadCount: 5 });
```

## Error Handlers

Test different error scenarios:

```typescript
import { enableGlobalErrorHandlers } from '@/mocks/utils';

// Enable server errors for all requests
enableGlobalErrorHandlers('serverError');

// Other available types:
// 'networkError', 'unauthorized', 'forbidden', 'notFound', 'rateLimit'
```

## CRUD Operations

Quickly mock CRUD endpoints:

```typescript
import { createCrudMocks } from '@/mocks/utils';
import { server } from '@/mocks/server';

const mockUsers = [/* array of users */];
const userCrudHandlers = createCrudMocks('users', mockUsers);

// Use in tests
server.use(...userCrudHandlers);
```

## Configuration

### Environment Variables

- `VITE_API_URL` - Base URL for API requests (default: `http://localhost:3001`)

### Customization

To add new endpoints, edit `src/mocks/handlers.ts`:

```typescript
import { http, HttpResponse } from 'msw';

export const customHandlers = [
  http.get('/api/custom-endpoint', () => {
    return HttpResponse.json({ data: 'custom response' });
  }),
];

// Add to main handlers array
export const handlers = [
  ...authHandlers,
  ...userHandlers,
  ...customHandlers, // Add your custom handlers
];
```

## Best Practices

1. **Use Mock Data Factories**: Create consistent mock data using the provided factory functions
2. **Reset Between Tests**: The server automatically resets handlers between tests
3. **Override Handlers**: Use `server.use()` to override specific endpoints in individual tests
4. **Test Error States**: Always test both success and error scenarios
5. **Realistic Data**: Use realistic mock data that matches your actual API responses
6. **Document Changes**: Update this README when adding new endpoints or utilities

## Troubleshooting

### MSW Not Working in Tests

Ensure the server is properly set up in your test configuration:

```typescript
// In test-setup.ts
import { server } from '@/mocks/server';

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());
```

### MSW Not Working in Development

1. Ensure the service worker is registered
2. Check browser console for MSW messages
3. Verify the `public/mockServiceWorker.js` file exists (run `npx msw init public/`)

### Unhandled Requests

Configure how MSW handles unmatched requests:

```typescript
// In server.ts
server.listen({
  onUnhandledRequest: 'error', // 'bypass' | 'warn' | 'error'
});
```

## Resources

- [MSW Documentation](https://mswjs.io/docs/)
- [MSW API Reference](https://mswjs.io/docs/api/)
- [Testing with MSW](https://mswjs.io/docs/getting-started/mocks/rest-api)