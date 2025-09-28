/**
 * ContentLoader Component Examples
 *
 * This file demonstrates various usage patterns for the ContentLoader component
 * across different scenarios and content types.
 */

import React from 'react';
import ContentLoader from './ContentLoader';
import { useSelector } from 'react-redux';
import { selectOperationLoading } from '../features/loadingSlice';
import { api } from '../services/api';

// ============================================================================
// Basic Usage Examples
// ============================================================================

export const BasicLoadingExample: React.FC = () => {
  const [isLoading, setIsLoading] = React.useState(true);

  React.useEffect(() => {
    const timer = setTimeout(() => setIsLoading(false), 3000);
    return () => clearTimeout(timer);
  }, []);

  return (
    <ContentLoader isLoading={isLoading} contentType="text" skeletonCount={3}>
      <div className="space-y-4">
        <h2 className="text-2xl font-bold">Article Title</h2>
        <p>This is the main content that will appear after loading completes.</p>
        <p>The ContentLoader shows skeleton placeholders while loading.</p>
      </div>
    </ContentLoader>
  );
};

// ============================================================================
// RTK Query Integration Examples
// ============================================================================

// Example API endpoint (would be in your actual API service)
const exampleApi = api.injectEndpoints({
  endpoints: (builder) => ({
    fetchUserProfile: builder.query<any, string>({
      query: (userId) => `/users/${userId}`,
      providesTags: ['User'],
    }),
    fetchPosts: builder.query<any[], void>({
      query: () => '/posts',
      providesTags: ['Content'],
    }),
  }),
});

export const RTKQueryExample: React.FC<{ userId: string }> = ({ userId }) => {
  const { data: user, isLoading, error } = exampleApi.useFetchUserProfileQuery(userId);

  return (
    <ContentLoader
      isLoading={isLoading}
      error={error}
      contentType="card"
      rtkQueryKey="fetchUserProfile"
      onRetry={() => {
        // RTK Query automatically handles retries, but you can manually refetch
        // or implement custom retry logic here
      }}
    >
      {user && (
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center space-x-4">
            <img
              src={user.avatar}
              alt={user.name}
              className="w-16 h-16 rounded-full"
            />
            <div>
              <h3 className="text-xl font-semibold">{user.name}</h3>
              <p className="text-gray-600">{user.email}</p>
            </div>
          </div>
          <p className="mt-4">{user.bio}</p>
        </div>
      )}
    </ContentLoader>
  );
};

// ============================================================================
// Progress Loading Examples
// ============================================================================

export const ProgressLoadingExample: React.FC = () => {
  const [progress, setProgress] = React.useState(0);
  const [isLoading, setIsLoading] = React.useState(true);

  React.useEffect(() => {
    const interval = setInterval(() => {
      setProgress((prev) => {
        if (prev >= 100) {
          setIsLoading(false);
          clearInterval(interval);
          return 100;
        }
        return prev + 10;
      });
    }, 500);

    return () => clearInterval(interval);
  }, []);

  return (
    <ContentLoader
      isLoading={isLoading}
      loadingType="progress"
      progress={progress}
      progressMessage="Uploading files..."
    >
      <div className="text-center p-8">
        <h3 className="text-lg font-semibold text-green-600">Upload Complete!</h3>
        <p>Your files have been successfully uploaded.</p>
      </div>
    </ContentLoader>
  );
};

// ============================================================================
// Error Handling Examples
// ============================================================================

export const ErrorHandlingExample: React.FC = () => {
  const [error, setError] = React.useState<string | null>("Network connection failed");
  const [retryCount, setRetryCount] = React.useState(0);

  const handleRetry = () => {
    setRetryCount(prev => prev + 1);
    setError(null);

    // Simulate retry attempt
    setTimeout(() => {
      if (retryCount < 2) {
        setError("Still failing... try again");
      } else {
        setError(null); // Success after 3 attempts
      }
    }, 1000);
  };

  return (
    <ContentLoader
      error={error}
      onRetry={handleRetry}
      maxRetries={3}
      retryCount={retryCount}
    >
      <div className="text-center p-8">
        <h3 className="text-lg font-semibold text-green-600">Data Loaded Successfully!</h3>
        <p>This content is now available.</p>
      </div>
    </ContentLoader>
  );
};

// ============================================================================
// Fallback Content Examples
// ============================================================================

export const FallbackContentExample: React.FC = () => {
  const [isOnline, setIsOnline] = React.useState(navigator.onLine);
  const [hasCache, setHasCache] = React.useState(true);

  React.useEffect(() => {
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  return (
    <ContentLoader
      isLoading={false}
      isOffline={!isOnline}
      isFallback={hasCache && !isOnline}
      fallbackMessage={isOnline ? "Using cached content" : "Showing offline content"}
    >
      <div className="space-y-4">
        <h2 className="text-2xl font-bold">Cached Article</h2>
        <p>This content is available offline or from cache.</p>
        <div className="text-sm text-gray-500">
          Status: {isOnline ? "Online" : "Offline"} |
          Cache: {hasCache ? "Available" : "Not available"}
        </div>
      </div>
    </ContentLoader>
  );
};

// ============================================================================
// Different Content Type Examples
// ============================================================================

export const ContentTypeExamples: React.FC = () => {
  const [activeType, setActiveType] = React.useState<string>('text');
  const [isLoading, setIsLoading] = React.useState(false);

  const contentTypes = [
    'text', 'image', 'card', 'list', 'table', 'chat', 'media'
  ];

  const toggleLoading = (type: string) => {
    setActiveType(type);
    setIsLoading(true);
    setTimeout(() => setIsLoading(false), 2000);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap gap-2">
        {contentTypes.map((type) => (
          <button
            key={type}
            onClick={() => toggleLoading(type)}
            className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
          >
            {type}
          </button>
        ))}
      </div>

      <ContentLoader
        isLoading={isLoading}
        contentType={activeType as any}
        skeletonCount={3}
        size="lg"
      >
        <div className="p-6 bg-gray-50 rounded-lg">
          <h3 className="text-lg font-semibold mb-4">
            {activeType.toUpperCase()} Content Loaded
          </h3>
          <p>This represents the actual content for {activeType} type.</p>
        </div>
      </ContentLoader>
    </div>
  );
};

// ============================================================================
// Custom Skeleton Example
// ============================================================================

export const CustomSkeletonExample: React.FC = () => {
  const [isLoading, setIsLoading] = React.useState(true);

  const customSkeleton = () => (
    <div className="space-y-4">
      <div className="flex items-center space-x-4">
        <div className="w-20 h-20 bg-gray-200 rounded-full animate-pulse" />
        <div className="flex-1 space-y-2">
          <div className="h-6 bg-gray-200 rounded animate-pulse" />
          <div className="h-4 bg-gray-200 rounded w-3/4 animate-pulse" />
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4">
        {[1, 2, 3].map((i) => (
          <div key={i} className="space-y-2">
            <div className="h-32 bg-gray-200 rounded animate-pulse" />
            <div className="h-4 bg-gray-200 rounded animate-pulse" />
          </div>
        ))}
      </div>

      <div className="space-y-2">
        <div className="h-4 bg-gray-200 rounded animate-pulse" />
        <div className="h-4 bg-gray-200 rounded w-5/6 animate-pulse" />
        <div className="h-4 bg-gray-200 rounded w-4/6 animate-pulse" />
      </div>
    </div>
  );

  React.useEffect(() => {
    const timer = setTimeout(() => setIsLoading(false), 3000);
    return () => clearTimeout(timer);
  }, []);

  return (
    <ContentLoader
      isLoading={isLoading}
      renderSkeleton={customSkeleton}
    >
      <div className="space-y-4">
        <div className="flex items-center space-x-4">
          <img
            src="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=80&h=80&fit=crop&auto=format"
            alt="Profile"
            className="w-20 h-20 rounded-full"
          />
          <div>
            <h3 className="text-xl font-semibold">John Doe</h3>
            <p className="text-gray-600">Software Developer</p>
          </div>
        </div>

        <div className="grid grid-cols-3 gap-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="space-y-2">
              <img
                src={`https://images.unsplash.com/photo-${[
                  '1506905925661-52d59fd7e096',
                  '1516117172878-fd2c41f4a759',
                  '1506755594592-349d12a7c9c3'
                ][i-1]}?w=200&h=128&fit=crop&auto=format`}
                alt={`Image ${i}`}
                className="w-full h-32 object-cover rounded"
              />
              <p className="text-sm">Caption {i}</p>
            </div>
          ))}
        </div>

        <div>
          <p>This is the actual content that appears after the custom skeleton loading is complete.</p>
          <p>The custom skeleton mimics the exact layout of this content.</p>
          <p>This provides a seamless loading experience for users.</p>
        </div>
      </div>
    </ContentLoader>
  );
};

// ============================================================================
// Accessibility Example
// ============================================================================

export const AccessibilityExample: React.FC = () => {
  const [isLoading, setIsLoading] = React.useState(false);
  const [announceChanges, setAnnounceChanges] = React.useState(true);
  const [reduceMotion, setReduceMotion] = React.useState(false);

  const startLoading = () => {
    setIsLoading(true);
    setTimeout(() => setIsLoading(false), 3000);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap gap-4">
        <button
          onClick={startLoading}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          disabled={isLoading}
        >
          {isLoading ? 'Loading...' : 'Start Loading'}
        </button>

        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={announceChanges}
            onChange={(e) => setAnnounceChanges(e.target.checked)}
          />
          <span>Announce changes to screen readers</span>
        </label>

        <label className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={reduceMotion}
            onChange={(e) => setReduceMotion(e.target.checked)}
          />
          <span>Reduce motion (accessibility)</span>
        </label>
      </div>

      <ContentLoader
        isLoading={isLoading}
        loadingLabel="Loading important user data, please wait"
        announceChanges={announceChanges}
        reduceMotion={reduceMotion}
        contentType="card"
      >
        <div className="bg-white border rounded-lg p-6">
          <h3 className="text-xl font-semibold mb-4">Accessible Content</h3>
          <p>This content demonstrates proper accessibility features:</p>
          <ul className="list-disc list-inside mt-2 space-y-1">
            <li>Screen reader announcements for loading states</li>
            <li>Proper ARIA labels and roles</li>
            <li>Reduced motion support</li>
            <li>Keyboard navigation support</li>
            <li>High contrast loading indicators</li>
          </ul>
        </div>
      </ContentLoader>
    </div>
  );
};

// ============================================================================
// Performance Example
// ============================================================================

export const PerformanceExample: React.FC = () => {
  const [showExpensive, setShowExpensive] = React.useState(false);
  const [isLoading, setIsLoading] = React.useState(false);

  const loadExpensiveContent = () => {
    setIsLoading(true);
    // Simulate expensive operation
    setTimeout(() => {
      setShowExpensive(true);
      setIsLoading(false);
    }, 2000);
  };

  return (
    <div className="space-y-4">
      <button
        onClick={loadExpensiveContent}
        className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
        disabled={isLoading}
      >
        Load Heavy Content
      </button>

      <ContentLoader
        isLoading={isLoading}
        deferRender={true}  // Don't render skeleton immediately
        size="lg"           // Consistent sizing to prevent layout shift
        contentType="table"
      >
        {showExpensive && (
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Heavy Data Table</h3>
            <div className="overflow-x-auto">
              <table className="min-w-full border-collapse border border-gray-300">
                <thead>
                  <tr className="bg-gray-50">
                    {['ID', 'Name', 'Email', 'Role', 'Created', 'Status'].map((header) => (
                      <th key={header} className="border border-gray-300 px-4 py-2 text-left">
                        {header}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {Array.from({ length: 10 }).map((_, i) => (
                    <tr key={i} className={i % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>
                      <td className="border border-gray-300 px-4 py-2">{i + 1}</td>
                      <td className="border border-gray-300 px-4 py-2">User {i + 1}</td>
                      <td className="border border-gray-300 px-4 py-2">user{i + 1}@example.com</td>
                      <td className="border border-gray-300 px-4 py-2">Developer</td>
                      <td className="border border-gray-300 px-4 py-2">2024-01-{String(i + 1).padStart(2, '0')}</td>
                      <td className="border border-gray-300 px-4 py-2">
                        <span className="px-2 py-1 text-xs bg-green-100 text-green-800 rounded">
                          Active
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </ContentLoader>
    </div>
  );
};

// ============================================================================
// Complete Demo Component
// ============================================================================

export const ContentLoaderDemo: React.FC = () => {
  const [activeExample, setActiveExample] = React.useState('basic');

  const examples = {
    basic: { component: BasicLoadingExample, title: 'Basic Loading' },
    rtk: { component: RTKQueryExample, title: 'RTK Query Integration', props: { userId: '123' } },
    progress: { component: ProgressLoadingExample, title: 'Progress Loading' },
    error: { component: ErrorHandlingExample, title: 'Error Handling' },
    fallback: { component: FallbackContentExample, title: 'Fallback Content' },
    types: { component: ContentTypeExamples, title: 'Content Types' },
    custom: { component: CustomSkeletonExample, title: 'Custom Skeleton' },
    accessibility: { component: AccessibilityExample, title: 'Accessibility' },
    performance: { component: PerformanceExample, title: 'Performance' },
  };

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      <div>
        <h1 className="text-3xl font-bold mb-4">ContentLoader Examples</h1>
        <p className="text-gray-600 mb-6">
          Comprehensive examples of the ContentLoader component in various scenarios.
        </p>
      </div>

      <div className="flex flex-wrap gap-2 mb-6">
        {Object.entries(examples).map(([key, { title }]) => (
          <button
            key={key}
            onClick={() => setActiveExample(key)}
            className={`px-3 py-2 rounded transition-colors ${
              activeExample === key
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200 hover:bg-gray-300'
            }`}
          >
            {title}
          </button>
        ))}
      </div>

      <div className="border rounded-lg p-6 min-h-96">
        <h2 className="text-xl font-semibold mb-4">
          {examples[activeExample as keyof typeof examples].title}
        </h2>

        {React.createElement(
          examples[activeExample as keyof typeof examples].component,
          examples[activeExample as keyof typeof examples].props || {}
        )}
      </div>
    </div>
  );
};

export default ContentLoaderDemo;