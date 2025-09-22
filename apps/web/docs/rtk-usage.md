# RTK Redux Toolkit Usage Guide

Complete guide for using Redux Toolkit with RTK Query in the Tchat application.

## Table of Contents

- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [API Services](#api-services)
- [State Management](#state-management)
- [Caching & Performance](#caching--performance)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Best Practices](#best-practices)

## Quick Start

### Basic Setup

The Redux store is already configured in your application. To use it in components:

```typescript
import { useAppSelector, useAppDispatch } from '../store/hooks';
import { useLoginMutation, useGetCurrentUserQuery } from '../services/auth';

export function MyComponent() {
  const isAuthenticated = useAppSelector(state => state.auth.isAuthenticated);
  const dispatch = useAppDispatch();

  return (
    <div>
      {isAuthenticated ? 'Logged in' : 'Not logged in'}
    </div>
  );
}
```

### Store Structure

```typescript
interface RootState {
  auth: AuthState;           // Authentication state
  ui: UIState;              // UI preferences and notifications
  loading: LoadingState;    // Loading states for operations
  api: ApiState;            // RTK Query cache and metadata
}
```

## Authentication

### Login Flow

```typescript
import { useLoginMutation } from '../services/auth';

export function LoginForm() {
  const [login, { isLoading, error }] = useLoginMutation();

  const handleLogin = async (credentials: LoginRequest) => {
    try {
      const result = await login(credentials).unwrap();
      // User automatically logged in via middleware
      console.log('Login successful:', result.user);
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  return (
    <form onSubmit={handleLogin}>
      {/* Form implementation */}
    </form>
  );
}
```

### Authentication State

```typescript
// Check authentication status
const isAuthenticated = useAppSelector(state => state.auth.isAuthenticated);
const user = useAppSelector(state => state.auth.user);
const tokenExpiresAt = useAppSelector(state => state.auth.expiresAt);

// Manual logout
const dispatch = useAppDispatch();
const handleLogout = () => {
  dispatch(logout());
};
```

### Protected Routes

```typescript
import { useGetCurrentUserQuery } from '../services/auth';

export function ProtectedComponent() {
  const { data: user, isLoading, error } = useGetCurrentUserQuery();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Please log in</div>;

  return <div>Welcome, {user.displayName}!</div>;
}
```

### Token Refresh

Token refresh is handled automatically by middleware:

- Tokens are refreshed 5 minutes before expiry
- Failed refresh attempts trigger automatic logout
- Manual refresh available via `useRefreshTokenMutation()`

## API Services

### Users API

```typescript
import {
  useListUsersQuery,
  useGetUserByIdQuery,
  useUpdateUserMutation,
} from '../services/users';

export function UsersList() {
  const { data, isLoading, error } = useListUsersQuery({
    page: 1,
    limit: 20,
    search: 'john',
  });

  if (isLoading) return <div>Loading users...</div>;
  if (error) return <div>Error loading users</div>;

  return (
    <ul>
      {data?.data.map(user => (
        <li key={user.id}>{user.displayName}</li>
      ))}
    </ul>
  );
}

export function UserProfile({ userId }: { userId: string }) {
  const { data: user } = useGetUserByIdQuery(userId);
  const [updateUser] = useUpdateUserMutation();

  const handleUpdate = async () => {
    await updateUser({
      id: userId,
      data: { displayName: 'New Name' }
    });
  };

  return (
    <div>
      <h1>{user?.displayName}</h1>
      <button onClick={handleUpdate}>Update</button>
    </div>
  );
}
```

### Messages API

```typescript
import {
  useListMessagesQuery,
  useSendMessageMutation,
  useInfiniteMessages,
} from '../services/messages';

export function ChatMessages({ chatId }: { chatId: string }) {
  const { data, isLoading } = useListMessagesQuery({ chatId });
  const [sendMessage] = useSendMessageMutation();

  const handleSend = async (content: string) => {
    // Optimistic update - message appears immediately
    await sendMessage({
      chatId,
      content,
      type: 'text',
    });
  };

  return (
    <div>
      {data?.items.map(message => (
        <div key={message.id}>{message.content}</div>
      ))}
      <button onClick={() => handleSend('Hello!')}>
        Send Message
      </button>
    </div>
  );
}

// Infinite scroll implementation
export function InfiniteMessagesList({ chatId }: { chatId: string }) {
  const { loadMore, data, isLoading } = useInfiniteMessages(chatId);

  return (
    <div>
      {data?.items.map(message => (
        <div key={message.id}>{message.content}</div>
      ))}
      <button onClick={() => loadMore(data?.pagination.nextCursor)}>
        Load More
      </button>
    </div>
  );
}
```

### Chats API

```typescript
import {
  useListChatsQuery,
  useCreateChatMutation,
  useMarkChatAsReadMutation,
} from '../services/chats';

export function ChatsList() {
  const { data: chats } = useListChatsQuery({ limit: 20 });
  const [createChat] = useCreateChatMutation();
  const [markAsRead] = useMarkChatAsReadMutation();

  const handleCreateChat = async () => {
    await createChat({
      name: 'New Chat',
      type: 'group',
      participants: ['user1', 'user2'],
    });
  };

  return (
    <div>
      {chats?.map(chat => (
        <div key={chat.id} onClick={() => markAsRead(chat.id)}>
          {chat.name}
          {chat.unreadCount > 0 && (
            <span className="badge">{chat.unreadCount}</span>
          )}
        </div>
      ))}
      <button onClick={handleCreateChat}>Create Chat</button>
    </div>
  );
}
```

## State Management

### UI State

```typescript
import {
  setTheme,
  toggleSidebar,
  addNotification,
  setLoading,
} from '../features/uiSlice';

export function ThemeToggle() {
  const theme = useAppSelector(state => state.ui.theme);
  const dispatch = useAppDispatch();

  return (
    <button onClick={() => dispatch(setTheme('dark'))}>
      Current theme: {theme}
    </button>
  );
}

export function NotificationSystem() {
  const notifications = useAppSelector(state => state.ui.notifications);
  const dispatch = useAppDispatch();

  const showNotification = () => {
    dispatch(addNotification({
      type: 'success',
      message: 'Operation completed!',
      duration: 5000,
    }));
  };

  return (
    <div>
      {notifications.map(notification => (
        <div key={notification.id} className={`alert alert-${notification.type}`}>
          {notification.message}
        </div>
      ))}
      <button onClick={showNotification}>Show Notification</button>
    </div>
  );
}
```

### Loading States

```typescript
import {
  selectGlobalLoading,
  selectOperationLoading,
  selectAnyLoading,
} from '../features/loadingSlice';

export function LoadingIndicator() {
  const isGlobalLoading = useAppSelector(selectGlobalLoading);
  const loginLoading = useAppSelector(selectOperationLoading('login'));
  const anyLoading = useAppSelector(selectAnyLoading);

  if (isGlobalLoading) {
    return <div className="global-spinner">Loading...</div>;
  }

  if (loginLoading?.isLoading) {
    return <div>{loginLoading.message}</div>;
  }

  return null;
}
```

## Caching & Performance

### Cache Configuration

```typescript
// Cache TTL: 60 seconds
// Refetch on focus: enabled
// Refetch on reconnect: enabled
// Refetch on mount: 30 seconds

// Manual cache invalidation
import { api } from '../services/api';

const dispatch = useAppDispatch();

// Invalidate specific tags
dispatch(api.util.invalidateTags(['User', 'Message']));

// Reset entire API cache
dispatch(api.util.resetApiState());
```

### Prefetching

```typescript
import { usePrefetch } from '../services/prefetch';

export function ChatPreview({ chatId }: { chatId: string }) {
  const { prefetchChatMessages } = usePrefetch();

  return (
    <div
      onMouseEnter={() => prefetchChatMessages(chatId)}
      onFocus={() => prefetchChatMessages(chatId)}
    >
      Chat Preview
    </div>
  );
}
```

### Optimistic Updates

Optimistic updates are built into mutations:

```typescript
// Messages appear immediately, rollback on failure
const [sendMessage] = useSendMessageMutation();

await sendMessage({ chatId, content: 'Hello!' });
// Message appears in UI immediately
// If API call fails, message is automatically removed
```

### Request Deduplication

RTK Query automatically deduplicates identical requests:

```typescript
// These will result in only one API call
const query1 = useListUsersQuery({ page: 1 });
const query2 = useListUsersQuery({ page: 1 });
const query3 = useListUsersQuery({ page: 1 });
```

## Error Handling

### Automatic Error Handling

Errors are automatically handled by middleware:

- Network errors show "Connection lost" notifications
- Validation errors (422) show field-specific messages
- Rate limiting (429) shows retry guidance
- Server errors (5xx) show generic error messages

### Manual Error Handling

```typescript
export function UserForm() {
  const [updateUser, { error }] = useUpdateUserMutation();

  const handleSubmit = async (data) => {
    try {
      await updateUser({ id: 'user1', data }).unwrap();
    } catch (error) {
      if (isApiError(error)) {
        console.log('Error code:', error.data.error.code);
        console.log('Error message:', error.data.error.message);

        if (error.status === 422) {
          // Handle validation errors
          const details = error.data.error.details;
          Object.entries(details).forEach(([field, message]) => {
            console.log(`${field}: ${message}`);
          });
        }
      }
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {error && <div className="error">{error.message}</div>}
      {/* Form fields */}
    </form>
  );
}
```

### Retry Logic

Automatic retry with exponential backoff:

- Network errors: 3 retries with 1s, 2s, 4s delays
- Server errors (5xx): 3 retries with backoff
- 401/403 errors: No retry (handled by auth middleware)

## Testing

### Testing Components with RTK

```typescript
import { renderWithProviders } from '../test-utils/renderWithProviders';

describe('UserProfile', () => {
  it('should display user information', async () => {
    const { getByText } = renderWithProviders(
      <UserProfile userId="1" />,
      {
        preloadedState: {
          auth: { isAuthenticated: true, user: mockUser },
        },
      }
    );

    await waitFor(() => {
      expect(getByText('John Doe')).toBeInTheDocument();
    });
  });
});
```

### Mocking API Calls

```typescript
import { setupServer } from 'msw/node';
import { handlers } from '../mocks/handlers';

const server = setupServer(...handlers);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

test('handles user fetch error', async () => {
  server.use(
    http.get('/api/users/1', () => {
      return HttpResponse.json(
        { error: 'User not found' },
        { status: 404 }
      );
    })
  );

  // Test error handling
});
```

## Best Practices

### 1. Use Typed Hooks

```typescript
// ✅ Good
import { useAppSelector, useAppDispatch } from '../store/hooks';

// ❌ Avoid
import { useSelector, useDispatch } from 'react-redux';
```

### 2. Handle Loading States

```typescript
// ✅ Good
const { data, isLoading, error } = useListUsersQuery();

if (isLoading) return <Spinner />;
if (error) return <ErrorMessage error={error} />;

// ❌ Avoid
const { data } = useListUsersQuery();
return <div>{data?.map(...)}</div>; // Can crash if data is undefined
```

### 3. Use Specific Selectors

```typescript
// ✅ Good
const username = useAppSelector(state => state.auth.user?.username);

// ❌ Avoid
const auth = useAppSelector(state => state.auth);
const username = auth.user?.username;
```

### 4. Destructure Mutations Properly

```typescript
// ✅ Good
const [createUser, { isLoading, error }] = useCreateUserMutation();

// ❌ Avoid
const mutation = useCreateUserMutation();
const isLoading = mutation[1].isLoading;
```

### 5. Handle Optimistic Updates

```typescript
// ✅ Good - Let RTK Query handle optimistic updates
const [sendMessage] = useSendMessageMutation();
await sendMessage(messageData);

// ❌ Avoid - Manual optimistic updates
const [sendMessage] = useSendMessageMutation();
dispatch(addMessageOptimistically(messageData));
await sendMessage(messageData);
```

### 6. Use Prefetching Wisely

```typescript
// ✅ Good - Prefetch on hover/focus
<div onMouseEnter={() => prefetchUserProfile(userId)}>

// ❌ Avoid - Aggressive prefetching
useEffect(() => {
  allUserIds.forEach(id => prefetchUserProfile(id));
}, []);
```

### 7. Cache Management

```typescript
// ✅ Good - Use tags for related data
providesTags: (result, error, id) => [{ type: 'User', id }]
invalidatesTags: [{ type: 'User', id }]

// ❌ Avoid - Invalidating everything
invalidatesTags: ['User', 'Message', 'Chat']
```

### 8. Error Boundaries

```typescript
// ✅ Good - Wrap your app in error boundaries
<ErrorBoundary>
  <Provider store={store}>
    <PersistGate persistor={persistor}>
      <App />
    </PersistGate>
  </Provider>
</ErrorBoundary>
```

## Debugging

### Redux DevTools

The application includes enhanced Redux DevTools configuration:

- **Time travel debugging**: Step through actions
- **Action filtering**: Hide noisy actions
- **State sanitization**: Sensitive data is redacted
- **Trace**: See where actions were dispatched from

### Common Issues

1. **"Cannot read property of undefined"**
   - Always check loading states before accessing data
   - Use optional chaining: `user?.displayName`

2. **"Too many re-renders"**
   - Avoid creating new objects in selectors
   - Use `createSelector` for derived data

3. **"Network request failed"**
   - Check network connectivity
   - Verify API endpoints are correct
   - Check CORS configuration

4. **"Token expired"**
   - Automatic refresh should handle this
   - Check token expiration times
   - Verify refresh token is valid

## Performance Tips

1. **Use selective re-rendering**
   - Split large components into smaller ones
   - Use `React.memo` for expensive components

2. **Optimize selectors**
   - Use `createSelector` for computed data
   - Keep selectors simple and focused

3. **Manage cache size**
   - Set appropriate `keepUnusedDataFor` values
   - Invalidate unused cache entries

4. **Batch mutations**
   - Use `useBatchMutation` for bulk operations
   - Combine related updates

5. **Lazy load data**
   - Use pagination for large datasets
   - Implement infinite scroll with cursors

This guide covers the core patterns for using RTK in the Tchat application. For more advanced usage, refer to the [RTK Query documentation](https://redux-toolkit.js.org/rtk-query/overview) and examine the implementation in the services directory.