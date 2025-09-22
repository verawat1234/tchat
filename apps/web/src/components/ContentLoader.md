# ContentLoader Component

A comprehensive content loading indicator system for React applications with RTK Query integration, accessibility features, and performance optimizations.

## Features

- **Multiple Loading States**: Skeleton loaders, spinners, progress indicators, and minimal loaders
- **Content-Aware Skeletons**: Specialized skeletons for text, images, cards, lists, tables, chat, and media
- **RTK Query Integration**: Seamless integration with Redux Toolkit Query loading states
- **Error Handling**: Comprehensive error states with retry functionality
- **Fallback Content**: Support for offline and cached content indicators
- **Accessibility**: Full ARIA support and screen reader announcements
- **Performance**: Layout shift prevention and deferred rendering
- **Customization**: Flexible props and custom skeleton renderers

## Installation

The ContentLoader component is already installed as part of your project. It requires the following dependencies:

- `class-variance-authority`: For component variants (already installed)
- `clsx` and `tailwind-merge`: For styling utilities (already installed)
- `lucide-react`: For icons (already installed)
- `react-redux`: For RTK Query integration (already installed)

## Basic Usage

### Simple Loading State

```tsx
import ContentLoader from './components/ContentLoader';

function MyComponent() {
  const [isLoading, setIsLoading] = useState(true);

  return (
    <ContentLoader isLoading={isLoading} contentType="text">
      <div>Your content goes here</div>
    </ContentLoader>
  );
}
```

### RTK Query Integration

```tsx
import { useGetUserQuery } from './api/userApi';
import ContentLoader from './components/ContentLoader';

function UserProfile({ userId }: { userId: string }) {
  const { data: user, isLoading, error } = useGetUserQuery(userId);

  return (
    <ContentLoader
      isLoading={isLoading}
      error={error}
      contentType="card"
      rtkQueryKey="getUser"
    >
      {user && (
        <div className="user-profile">
          <h2>{user.name}</h2>
          <p>{user.email}</p>
        </div>
      )}
    </ContentLoader>
  );
}
```

## Content Types

The component supports different content types with appropriate skeleton loaders:

### Text Content
```tsx
<ContentLoader
  isLoading={true}
  contentType="text"
  skeletonCount={3}
>
  <div>Text content</div>
</ContentLoader>
```

### Image Content
```tsx
<ContentLoader
  isLoading={true}
  contentType="image"
>
  <img src="image.jpg" alt="Description" />
</ContentLoader>
```

### Card Content
```tsx
<ContentLoader
  isLoading={true}
  contentType="card"
>
  <div className="card">
    <img src="..." />
    <h3>Title</h3>
    <p>Description</p>
  </div>
</ContentLoader>
```

### List Content
```tsx
<ContentLoader
  isLoading={true}
  contentType="list"
  skeletonCount={5}
>
  <ul>
    {items.map(item => <li key={item.id}>{item.name}</li>)}
  </ul>
</ContentLoader>
```

### Table Content
```tsx
<ContentLoader
  isLoading={true}
  contentType="table"
>
  <table>
    <thead>
      <tr><th>Name</th><th>Email</th></tr>
    </thead>
    <tbody>
      {users.map(user => (
        <tr key={user.id}>
          <td>{user.name}</td>
          <td>{user.email}</td>
        </tr>
      ))}
    </tbody>
  </table>
</ContentLoader>
```

### Chat Content
```tsx
<ContentLoader
  isLoading={true}
  contentType="chat"
  skeletonCount={4}
>
  <div className="chat-messages">
    {messages.map(msg => (
      <div key={msg.id} className="message">
        {msg.content}
      </div>
    ))}
  </div>
</ContentLoader>
```

### Media Content
```tsx
<ContentLoader
  isLoading={true}
  contentType="media"
>
  <div className="media-player">
    <video controls>
      <source src="video.mp4" type="video/mp4" />
    </video>
  </div>
</ContentLoader>
```

## Loading Types

### Skeleton Loader (Default)
```tsx
<ContentLoader
  isLoading={true}
  loadingType="skeleton"
  contentType="text"
>
  Content
</ContentLoader>
```

### Spinner Loader
```tsx
<ContentLoader
  isLoading={true}
  loadingType="spinner"
>
  Content
</ContentLoader>
```

### Progress Loader
```tsx
<ContentLoader
  isLoading={true}
  loadingType="progress"
  progress={75}
  progressMessage="Uploading files..."
>
  Content
</ContentLoader>
```

### Minimal Loader
```tsx
<ContentLoader
  isLoading={true}
  loadingType="minimal"
  loadingLabel="Please wait..."
>
  Content
</ContentLoader>
```

## Error Handling

### Basic Error Display
```tsx
<ContentLoader
  error="Failed to load data"
>
  Content
</ContentLoader>
```

### Error with Retry
```tsx
<ContentLoader
  error={error}
  onRetry={() => refetch()}
  maxRetries={3}
  retryCount={retryCount}
>
  Content
</ContentLoader>
```

## Fallback States

### Offline Content
```tsx
<ContentLoader
  isOffline={!navigator.onLine}
  fallbackMessage="Showing offline content"
>
  {cachedContent}
</ContentLoader>
```

### Cached Content
```tsx
<ContentLoader
  isFallback={true}
  fallbackMessage="Using cached data"
>
  {cachedData}
</ContentLoader>
```

## Accessibility

### Screen Reader Support
```tsx
<ContentLoader
  isLoading={true}
  loadingLabel="Loading user profile"
  announceChanges={true}
>
  Content
</ContentLoader>
```

### Reduced Motion
```tsx
<ContentLoader
  isLoading={true}
  reduceMotion={true} // or auto-detect from user preferences
  animation="none"
>
  Content
</ContentLoader>
```

## Performance Optimizations

### Deferred Rendering
```tsx
<ContentLoader
  isLoading={true}
  deferRender={true} // Delays skeleton rendering by 100ms
>
  Content
</ContentLoader>
```

### Layout Shift Prevention
```tsx
<ContentLoader
  isLoading={true}
  size="lg" // Consistent sizing to prevent layout shifts
>
  Content
</ContentLoader>
```

## Custom Skeletons

### Custom Skeleton Renderer
```tsx
const customSkeleton = () => (
  <div className="custom-skeleton">
    <div className="h-8 w-32 bg-gray-200 rounded animate-pulse" />
    <div className="h-4 w-48 bg-gray-200 rounded animate-pulse mt-2" />
  </div>
);

<ContentLoader
  isLoading={true}
  renderSkeleton={customSkeleton}
>
  Content
</ContentLoader>
```

## Compound Components

The ContentLoader also provides individual components for direct use:

### Individual Skeleton Components
```tsx
import {
  TextSkeleton,
  ImageSkeleton,
  CardSkeleton,
  ListSkeleton,
  TableSkeleton,
  ChatSkeleton,
  MediaSkeleton
} from './components/ContentLoader';

// Use directly
<TextSkeleton count={3} />
<ImageSkeleton aspectRatio="aspect-square" />
<ListSkeleton count={5} />
```

### Individual Utility Components
```tsx
import {
  SpinnerLoader,
  ProgressLoader,
  ErrorState,
  FallbackState
} from './components/ContentLoader';

// Use directly
<SpinnerLoader size="lg" />
<ProgressLoader progress={50} message="Loading..." />
<ErrorState error="Something went wrong" onRetry={handleRetry} />
<FallbackState isOffline={true} />
```

### Through ContentLoader
```tsx
<ContentLoader.Text count={3} />
<ContentLoader.Image />
<ContentLoader.Spinner size="md" />
<ContentLoader.Progress progress={75} />
<ContentLoader.Error error="Failed" onRetry={retry} />
<ContentLoader.Fallback isOffline={true} />
```

## Custom Hooks

### useContentLoader Hook
```tsx
import { useContentLoader } from './components/ContentLoader';

function MyComponent() {
  const { isLoading, progress, message } = useContentLoader('fetchData');

  return (
    <div>
      {isLoading && <div>Loading: {message} ({progress}%)</div>}
    </div>
  );
}
```

### useAccessibilityAnnouncer Hook
```tsx
import { useAccessibilityAnnouncer } from './components/ContentLoader';

function MyComponent() {
  const { announce } = useAccessibilityAnnouncer(true);

  const handleAction = () => {
    announce('Action completed successfully');
  };

  return <button onClick={handleAction}>Do Action</button>;
}
```

## Complete Props Reference

```tsx
interface ContentLoaderProps {
  // Content and styling
  className?: string;
  children?: React.ReactNode;
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full' | 'auto';
  spacing?: 'tight' | 'normal' | 'loose';
  animation?: 'pulse' | 'wave' | 'fade' | 'none';

  // Loading state
  isLoading?: boolean;
  loadingType?: 'skeleton' | 'spinner' | 'progress' | 'minimal';
  contentType?: 'text' | 'image' | 'card' | 'list' | 'table' | 'chat' | 'media' | 'custom';
  skeletonCount?: number;

  // Progress
  progress?: number;
  progressMessage?: string;

  // Error handling
  error?: Error | string | null;
  onRetry?: () => void;
  maxRetries?: number;
  retryCount?: number;

  // Fallback states
  isFallback?: boolean;
  isOffline?: boolean;
  fallbackMessage?: string;

  // RTK Query integration
  rtkQueryKey?: string;

  // Accessibility
  loadingLabel?: string;
  announceChanges?: boolean;

  // Performance
  reduceMotion?: boolean;
  deferRender?: boolean;

  // Customization
  renderSkeleton?: () => React.ReactNode;
}
```

## Best Practices

1. **Choose the Right Content Type**: Use the appropriate `contentType` to match your actual content layout
2. **Consistent Sizing**: Use the `size` prop to prevent layout shifts
3. **Accessibility First**: Always provide meaningful `loadingLabel` and enable `announceChanges`
4. **Performance**: Use `deferRender` for expensive components and `reduceMotion` for accessibility
5. **Error Handling**: Always provide retry functionality for user-initiated actions
6. **RTK Query**: Use `rtkQueryKey` for automatic loading state integration

## Examples

See `ContentLoader.examples.tsx` for comprehensive usage examples including:
- Basic loading states
- RTK Query integration
- Progress loading
- Error handling
- Fallback content
- Different content types
- Custom skeletons
- Accessibility features
- Performance optimizations

## Testing

The component includes comprehensive test coverage in:
- `ContentLoader.test.tsx`: Full integration tests
- `ContentLoader.basic.test.tsx`: Basic component tests

Run tests with:
```bash
npm test ContentLoader
```