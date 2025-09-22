import React, { useState } from 'react';
import {
  ContentErrorBoundary,
  ContentErrorProvider,
  FastContentErrorBoundary,
  useContentErrorContext,
  useRTKQueryErrorIntegration,
  withContentErrorRecovery,
  ContentErrorType
} from '../ContentErrorBoundary';
import { Button } from '../ui/button';

// Example component that might fail
const UnstableContentComponent: React.FC<{ shouldFail?: boolean; errorType?: string }> = ({
  shouldFail = false,
  errorType = 'general'
}) => {
  if (shouldFail) {
    let error: Error;
    switch (errorType) {
      case 'network':
        error = new Error('Network request failed');
        error.name = 'NetworkError';
        break;
      case 'parsing':
        error = new Error('JSON parsing failed');
        break;
      case 'rtk':
        error = new Error('RTK Query endpoint failed');
        error.name = 'RTKQueryError';
        break;
      default:
        error = new Error('Something went wrong');
    }
    throw error;
  }

  return (
    <div className="p-4 border rounded-lg">
      <h3 className="font-semibold mb-2">Content Loaded Successfully</h3>
      <p>This content loaded without any errors.</p>
    </div>
  );
};

// Example of RTK Query integration
const RTKQueryExample: React.FC<{ hasError?: boolean }> = ({ hasError = false }) => {
  const { handleRTKError } = useRTKQueryErrorIntegration();

  // Simulate RTK Query response
  const mockError = hasError ? { message: 'Failed to fetch content' } : null;

  // This would trigger the error boundary if there's an error
  handleRTKError(mockError, 'fetchContent');

  return (
    <div className="p-4 border rounded-lg">
      <h3 className="font-semibold mb-2">RTK Query Content</h3>
      <p>This content was loaded via RTK Query.</p>
    </div>
  );
};

// Example error summary component
const ErrorSummaryComponent: React.FC = () => {
  const { errorHistory, hasRecentErrors, clearErrors } = useContentErrorContext();

  if (!hasRecentErrors) {
    return null;
  }

  return (
    <div className="p-4 mb-4 bg-yellow-50 border border-yellow-200 rounded-lg">
      <div className="flex justify-between items-center">
        <div>
          <h4 className="font-semibold text-yellow-800">Recent Errors Detected</h4>
          <p className="text-sm text-yellow-700">
            {errorHistory.length} error(s) in the last 5 minutes
          </p>
        </div>
        <Button onClick={clearErrors} variant="outline" size="sm">
          Clear
        </Button>
      </div>
    </div>
  );
};

// Enhanced component with error recovery
const EnhancedContentComponent = withContentErrorRecovery(UnstableContentComponent, {
  maxRetries: 2,
  retryDelay: 1000,
});

// Main example component
export const ContentErrorBoundaryExample: React.FC = () => {
  const [errorScenario, setErrorScenario] = useState<string>('none');
  const [componentType, setComponentType] = useState<string>('standard');

  const renderContentComponent = () => {
    const shouldFail = errorScenario !== 'none';
    const errorType = errorScenario;

    switch (componentType) {
      case 'fast':
        return (
          <FastContentErrorBoundary>
            <UnstableContentComponent shouldFail={shouldFail} errorType={errorType} />
          </FastContentErrorBoundary>
        );

      case 'enhanced':
        return <EnhancedContentComponent shouldFail={shouldFail} errorType={errorType} />;

      case 'rtk':
        return (
          <ContentErrorBoundary maxRetries={3}>
            <RTKQueryExample hasError={shouldFail} />
          </ContentErrorBoundary>
        );

      default:
        return (
          <ContentErrorBoundary
            maxRetries={3}
            retryDelay={2000}
            onError={(error, errorInfo) => {
              console.log('Custom error handler:', error.message);
            }}
          >
            <UnstableContentComponent shouldFail={shouldFail} errorType={errorType} />
          </ContentErrorBoundary>
        );
    }
  };

  return (
    <ContentErrorProvider>
      <div className="space-y-6 p-6">
        <div>
          <h2 className="text-2xl font-bold mb-4">Content Error Boundary Examples</h2>
          <p className="text-gray-600 mb-6">
            Demonstration of various error boundary configurations and error scenarios.
          </p>
        </div>

        {/* Error Summary */}
        <ErrorSummaryComponent />

        {/* Controls */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4 bg-gray-50 rounded-lg">
          <div>
            <label className="block text-sm font-medium mb-2">Error Scenario:</label>
            <select
              value={errorScenario}
              onChange={(e) => setErrorScenario(e.target.value)}
              className="w-full p-2 border rounded"
            >
              <option value="none">No Error</option>
              <option value="network">Network Error</option>
              <option value="parsing">Parsing Error</option>
              <option value="rtk">RTK Query Error</option>
              <option value="general">General Error</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Component Type:</label>
            <select
              value={componentType}
              onChange={(e) => setComponentType(e.target.value)}
              className="w-full p-2 border rounded"
            >
              <option value="standard">Standard Error Boundary</option>
              <option value="fast">Fast Error Boundary</option>
              <option value="enhanced">HOC Enhanced</option>
              <option value="rtk">RTK Query Integration</option>
            </select>
          </div>
        </div>

        {/* Content Area */}
        <div className="border-2 border-dashed border-gray-300 rounded-lg p-4">
          <h3 className="font-semibold mb-4">Content Area ({componentType})</h3>
          {renderContentComponent()}
        </div>

        {/* Feature Showcase */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div className="p-4 border rounded-lg">
            <h4 className="font-semibold mb-2">🔄 Auto-Retry</h4>
            <p className="text-sm text-gray-600">
              Automatically retries failed operations with exponential backoff.
            </p>
          </div>

          <div className="p-4 border rounded-lg">
            <h4 className="font-semibold mb-2">🌐 Network Aware</h4>
            <p className="text-sm text-gray-600">
              Monitors network status and auto-retries when connection restored.
            </p>
          </div>

          <div className="p-4 border rounded-lg">
            <h4 className="font-semibold mb-2">🎯 Error Classification</h4>
            <p className="text-sm text-gray-600">
              Intelligently classifies errors for appropriate handling strategies.
            </p>
          </div>

          <div className="p-4 border rounded-lg">
            <h4 className="font-semibold mb-2">⚡ Performance Optimized</h4>
            <p className="text-sm text-gray-600">
              Fast mode available for performance-critical components.
            </p>
          </div>

          <div className="p-4 border rounded-lg">
            <h4 className="font-semibold mb-2">🔗 RTK Query Integration</h4>
            <p className="text-sm text-gray-600">
              Seamless integration with RTK Query error states.
            </p>
          </div>

          <div className="p-4 border rounded-lg">
            <h4 className="font-semibold mb-2">📊 Error Reporting</h4>
            <p className="text-sm text-gray-600">
              Comprehensive error logging and reporting capabilities.
            </p>
          </div>
        </div>

        {/* Usage Tips */}
        <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <h4 className="font-semibold text-blue-800 mb-2">💡 Usage Tips</h4>
          <ul className="text-sm text-blue-700 space-y-1">
            <li>• Use standard boundary for most content areas</li>
            <li>• Use fast boundary for performance-critical areas</li>
            <li>• Implement custom error reporting for production monitoring</li>
            <li>• Test different error scenarios during development</li>
            <li>• Monitor error patterns to improve content reliability</li>
          </ul>
        </div>
      </div>
    </ContentErrorProvider>
  );
};

export default ContentErrorBoundaryExample;