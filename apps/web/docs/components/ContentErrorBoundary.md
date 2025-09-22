# ContentErrorBoundary

Enhanced error boundary component for handling content management failures with comprehensive recovery mechanisms, intelligent error classification, and RTK Query integration.

## Features

- **Content-Specific Error Handling**: Intelligent classification of content-related errors
- **Automatic Retry Logic**: Exponential backoff retry for recoverable errors
- **Network Awareness**: Auto-retry when network connection is restored
- **RTK Query Integration**: Seamless integration with RTK Query error states
- **Performance Optimized**: Fast mode for performance-critical areas
- **Error Reporting**: Comprehensive error logging and reporting
- **Context Management**: Centralized error state management
- **Recovery Mechanisms**: Multiple recovery strategies and fallback options

## Basic Usage

```tsx
import { ContentErrorBoundary } from '../components/ContentErrorBoundary';

function App() {
  return (
    <ContentErrorBoundary maxRetries={3} retryDelay={2000}>
      <YourContentComponents />
    </ContentErrorBoundary>
  );
}
```

## Advanced Usage

### With Custom Error Handling

```tsx
import { ContentErrorBoundary, ContentError } from '../components/ContentErrorBoundary';

function App() {
  const handleError = (error: Error, errorInfo: ErrorInfo) => {
    // Custom error handling logic
    console.error('Content error occurred:', error);
  };

  const reportError = (error: ContentError) => {
    // Send to error tracking service
    analytics.captureException(error);
  };

  return (
    <ContentErrorBoundary
      maxRetries={5}
      retryDelay={1000}
      onError={handleError}
      reportError={reportError}
    >
      <YourContentComponents />
    </ContentErrorBoundary>
  );
}
```

### With Custom Fallback UI

```tsx
const CustomErrorFallback = () => (
  <div className="error-fallback">
    <h2>Content Unavailable</h2>
    <p>We're working to restore this content.</p>
  </div>
);

function App() {
  return (
    <ContentErrorBoundary fallback={<CustomErrorFallback />}>
      <YourContentComponents />
    </ContentErrorBoundary>
  );
}
```

## Error Types

The component automatically classifies errors into specific types:

- **NETWORK**: Network connectivity issues, fetch failures
- **LOADING**: Content loading timeouts, cancelled requests
- **PARSING**: JSON parsing errors, malformed content
- **RTK_QUERY**: RTK Query specific errors
- **CONTENT_MISSING**: 404 errors, missing content
- **PERMISSION**: Authorization failures, access denied
- **VALIDATION**: Content validation failures
- **UNKNOWN**: Unclassified errors

## RTK Query Integration

### Using the RTK Query Error Hook

```tsx
import { useRTKQueryErrorIntegration } from '../components/ContentErrorBoundary';
import { useGetContentQuery } from '../store/api';

function ContentComponent() {
  const { data, error, isLoading } = useGetContentQuery('content-id');
  const { handleRTKError } = useRTKQueryErrorIntegration();

  // This will trigger the error boundary if there's an RTK Query error
  handleRTKError(error, 'fetchContent');

  if (isLoading) return <div>Loading...</div>;
  return <div>{data?.content}</div>;
}
```

## Error Recovery Hook

For components that need custom recovery logic:

```tsx
import { useContentErrorRecovery } from '../components/ContentErrorBoundary';

function ContentComponent() {
  const { recoveryAttempts, isRecovering, attemptRecovery, resetRecovery } = useContentErrorRecovery();

  const handleRetry = () => {
    attemptRecovery(async () => {
      // Custom recovery logic
      await refetchContent();
      await validateContent();
    }, 3); // max 3 attempts
  };

  return (
    <div>
      {recoveryAttempts > 0 && (
        <p>Recovery attempts: {recoveryAttempts}</p>
      )}
      <button onClick={handleRetry} disabled={isRecovering}>
        {isRecovering ? 'Recovering...' : 'Retry'}
      </button>
    </div>
  );
}
```

## Performance Optimization

### Fast Content Error Boundary

For performance-critical content areas where you need minimal overhead:

```tsx
import { FastContentErrorBoundary } from '../components/ContentErrorBoundary';

function PerformanceCriticalArea() {
  return (
    <FastContentErrorBoundary>
      <HighFrequencyUpdatingComponent />
    </FastContentErrorBoundary>
  );
}
```

### Higher-Order Component

Wrap components automatically with error boundaries:

```tsx
import { withContentErrorRecovery } from '../components/ContentErrorBoundary';

const SafeContentComponent = withContentErrorRecovery(YourComponent, {
  maxRetries: 3,
  retryDelay: 1000,
  fastMode: false, // Use full error boundary
});
```

## Centralized Error Management

### Using the Error Context

```tsx
import { ContentErrorProvider, useContentErrorContext } from '../components/ContentErrorBoundary';

// Wrap your app with the provider
function App() {
  return (
    <ContentErrorProvider>
      <YourApp />
    </ContentErrorProvider>
  );
}

// Use the context in components
function ErrorSummary() {
  const { errorHistory, hasRecentErrors, clearErrors } = useContentErrorContext();

  return (
    <div>
      {hasRecentErrors && (
        <div className="error-alert">
          Recent errors detected ({errorHistory.length})
          <button onClick={clearErrors}>Clear</button>
        </div>
      )}
    </div>
  );
}
```

## Configuration Options

### ContentErrorBoundary Props

```tsx
interface Props {
  children: ReactNode;
  fallback?: ReactNode;           // Custom fallback UI
  maxRetries?: number;            // Maximum retry attempts (default: 3)
  retryDelay?: number;            // Base retry delay in ms (default: 2000)
  onError?: (error: Error, errorInfo: ErrorInfo) => void;  // Error callback
  reportError?: (error: ContentError) => void;             // Error reporting
}
```

### Error Recovery Options

```tsx
const options = {
  maxRetries: 5,        // Maximum retry attempts
  retryDelay: 1000,     // Base delay between retries
  fallback: <CustomUI />, // Custom fallback component
  fastMode: true,       // Use lightweight error boundary
};
```

## Error Reporting

The component automatically reports errors in development and production:

### Development Mode
- Console logging with detailed error information
- Component stack traces
- Error classification details

### Production Mode
- Local storage logging for debugging
- Integration points for error tracking services
- Sanitized error reporting

### Custom Error Reporting

```tsx
const reportError = (error: ContentError) => {
  // Send to your error tracking service
  if (process.env.NODE_ENV === 'production') {
    errorTrackingService.captureException({
      message: error.message,
      type: error.type,
      isRecoverable: error.isRecoverable,
      retryCount: error.retryCount,
      timestamp: error.timestamp,
      url: window.location.href,
      userAgent: navigator.userAgent,
    });
  }
};
```

## Network Integration

The component automatically monitors network status and:

- Shows network status indicators for network-related errors
- Auto-retries when connection is restored
- Adjusts retry strategies based on connectivity
- Provides offline/online status information

## Best Practices

1. **Place at Strategic Boundaries**: Wrap content-heavy components that might fail
2. **Use Appropriate Retry Counts**: Balance user experience with server load
3. **Implement Custom Recovery**: Provide domain-specific recovery strategies
4. **Monitor Error Patterns**: Use error reporting to identify common failure modes
5. **Test Error Scenarios**: Ensure graceful degradation in all failure modes
6. **Progressive Enhancement**: Start with basic boundaries, add features as needed

## Testing

The component includes comprehensive tests covering:

- Basic error boundary functionality
- Error classification accuracy
- Retry mechanism behavior
- Network status integration
- Context management
- RTK Query integration
- Performance optimization features

Run tests with:
```bash
npm test ContentErrorBoundary.test.tsx
```

## Migration from Basic Error Boundaries

If you're upgrading from a basic error boundary:

```tsx
// Before
<ErrorBoundary>
  <Content />
</ErrorBoundary>

// After
<ContentErrorBoundary maxRetries={3} retryDelay={2000}>
  <Content />
</ContentErrorBoundary>
```

The enhanced boundary is backward compatible and provides additional features without breaking existing implementations.