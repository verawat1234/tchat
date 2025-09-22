# Quickstart: RTK Backend API Integration

## Prerequisites

- Node.js 18+ and npm 9+
- Running backend API at `http://localhost:3001`
- React development environment

## Installation

```bash
# Install RTK and dependencies
npm install @reduxjs/toolkit react-redux

# Install dev dependencies for testing
npm install --save-dev msw @mswjs/data
```

## Basic Setup

### 1. Configure the Redux Store

Create `apps/web/src/store/index.ts`:
```typescript
import { configureStore } from '@reduxjs/toolkit';
import { api } from '../services/api';

export const store = configureStore({
  reducer: {
    [api.reducerPath]: api.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(api.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
```

### 2. Create Typed Hooks

Create `apps/web/src/store/hooks.ts`:
```typescript
import { useDispatch, useSelector } from 'react-redux';
import type { RootState, AppDispatch } from './index';

export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
```

### 3. Setup Base API Service

Create `apps/web/src/services/api.ts`:
```typescript
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

export const api = createApi({
  baseQuery: fetchBaseQuery({
    baseUrl: 'http://localhost:3001/api',
    prepareHeaders: (headers, { getState }) => {
      const token = (getState() as any).auth?.accessToken;
      if (token) {
        headers.set('authorization', `Bearer ${token}`);
      }
      return headers;
    },
  }),
  tagTypes: ['User', 'Message', 'Chat'],
  endpoints: () => ({}),
});
```

### 4. Provide Store to React

Update `apps/web/src/main.tsx`:
```typescript
import { Provider } from 'react-redux';
import { store } from './store';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Provider store={store}>
      <App />
    </Provider>
  </React.StrictMode>
);
```

## Creating API Endpoints

### Example: Authentication Service

Create `apps/web/src/services/auth.ts`:
```typescript
import { api } from './api';

export const authApi = api.injectEndpoints({
  endpoints: (builder) => ({
    login: builder.mutation({
      query: (credentials) => ({
        url: '/auth/login',
        method: 'POST',
        body: credentials,
      }),
      invalidatesTags: ['User'],
    }),
    getCurrentUser: builder.query({
      query: () => '/auth/me',
      providesTags: ['User'],
    }),
  }),
});

export const { useLoginMutation, useGetCurrentUserQuery } = authApi;
```

## Using RTK Query in Components

### Example: Login Component

```tsx
import { useLoginMutation } from '../services/auth';

function LoginForm() {
  const [login, { isLoading, error }] = useLoginMutation();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      const result = await login({
        email: 'user@example.com',
        password: 'password'
      }).unwrap();
      console.log('Login successful:', result);
    } catch (err) {
      console.error('Login failed:', err);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* form fields */}
      <button disabled={isLoading}>
        {isLoading ? 'Logging in...' : 'Login'}
      </button>
      {error && <div>Error: {error.message}</div>}
    </form>
  );
}
```

### Example: Fetching Data

```tsx
import { useGetCurrentUserQuery } from '../services/auth';

function Profile() {
  const { data, isLoading, error, refetch } = useGetCurrentUserQuery();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading profile</div>;

  return (
    <div>
      <h1>{data?.user.displayName}</h1>
      <button onClick={refetch}>Refresh</button>
    </div>
  );
}
```

## Testing with MSW

### Setup MSW Handlers

Create `apps/web/src/mocks/handlers.ts`:
```typescript
import { http, HttpResponse } from 'msw';

export const handlers = [
  http.post('/api/auth/login', () => {
    return HttpResponse.json({
      user: { id: '1', email: 'test@example.com' },
      tokens: { accessToken: 'test-token', refreshToken: 'refresh' }
    });
  }),
  http.get('/api/auth/me', () => {
    return HttpResponse.json({
      success: true,
      data: { id: '1', email: 'test@example.com' }
    });
  }),
];
```

### Run Tests

```bash
# Run unit tests
npm test

# Run integration tests
npm run test:integration
```

## Optimistic Updates

Example of optimistic update for creating a message:

```typescript
const [sendMessage] = useSendMessageMutation();

const handleSend = async (content: string) => {
  const tempId = Date.now().toString();

  try {
    await sendMessage({
      content,
      chatId: currentChatId,
      // Optimistic update
      onQueryStarted(arg, { dispatch, queryFulfilled }) {
        // Update cache optimistically
        const patchResult = dispatch(
          api.util.updateQueryData('getMessages', chatId, (draft) => {
            draft.items.push({
              id: tempId,
              content,
              status: 'pending'
            });
          })
        );

        // Rollback on error
        queryFulfilled.catch(patchResult.undo);
      },
    }).unwrap();
  } catch (error) {
    console.error('Failed to send message:', error);
  }
};
```

## Cache Management

### Manual Cache Invalidation

```typescript
import { api } from '../services/api';
import { useAppDispatch } from '../store/hooks';

function Settings() {
  const dispatch = useAppDispatch();

  const handleClearCache = () => {
    // Invalidate all cached data
    dispatch(api.util.invalidateTags(['User', 'Message', 'Chat']));
  };

  const handleRefreshUser = () => {
    // Invalidate specific tag
    dispatch(api.util.invalidateTags(['User']));
  };

  return (
    <div>
      <button onClick={handleClearCache}>Clear All Cache</button>
      <button onClick={handleRefreshUser}>Refresh User Data</button>
    </div>
  );
}
```

## Error Handling

### Global Error Handler

```typescript
import { isRejectedWithValue } from '@reduxjs/toolkit';
import type { Middleware } from '@reduxjs/toolkit';

export const rtkQueryErrorLogger: Middleware = () => (next) => (action) => {
  if (isRejectedWithValue(action)) {
    console.error('API call failed:', action.payload);

    // Handle specific error codes
    if (action.payload.status === 401) {
      // Trigger logout or token refresh
    }
  }

  return next(action);
};
```

## Performance Monitoring

### Debug with Redux DevTools

1. Install Redux DevTools Extension
2. Open DevTools in browser
3. Navigate to Redux tab
4. Monitor:
   - API call timing
   - Cache hits/misses
   - State changes
   - Action dispatches

## Verification Steps

1. **Store Configuration**
   - [ ] Redux store created with RTK Query middleware
   - [ ] TypeScript types properly configured
   - [ ] Provider wrapping app root

2. **API Setup**
   - [ ] Base API service created
   - [ ] Auth headers automatically attached
   - [ ] Tag types defined

3. **Endpoints**
   - [ ] Login endpoint working
   - [ ] Data fetching working
   - [ ] Cache invalidation working

4. **Testing**
   - [ ] MSW handlers created
   - [ ] Tests can mock API calls
   - [ ] Integration tests passing

5. **Developer Experience**
   - [ ] Redux DevTools showing API calls
   - [ ] TypeScript autocomplete working
   - [ ] Error messages helpful

## Common Issues

### CORS Errors
- Ensure backend allows frontend origin
- Check headers configuration

### Token Expiration
- Implement token refresh middleware
- Handle 401 responses globally

### Cache Not Updating
- Verify tag configuration
- Check invalidation logic

### TypeScript Errors
- Run `npm run type-check`
- Ensure all types imported correctly

## Next Steps

1. Add more API endpoints as needed
2. Implement token refresh logic
3. Add request retry logic
4. Setup error boundaries
5. Add loading skeletons
6. Implement offline support