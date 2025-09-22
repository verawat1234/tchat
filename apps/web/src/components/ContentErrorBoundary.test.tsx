import React from 'react';
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import {
  ContentErrorBoundary,
  ContentErrorType,
  ContentError,
  FastContentErrorBoundary,
  ContentErrorProvider,
  useContentErrorContext,
  useRTKQueryErrorIntegration,
  useContentErrorRecovery,
  withContentErrorRecovery
} from './ContentErrorBoundary';

// Mock the useContentTexts hook
vi.mock('../hooks/useContentText', () => ({
  useContentTexts: vi.fn(() => ({
    'error.boundary.title': { text: 'Content System Error' },
    'error.boundary.description': { text: 'Something went wrong with the content system.' },
    'error.boundary.details_summary': { text: 'Error Details (Development)' },
    'error.boundary.try_again': { text: 'Try Again' },
    'error.boundary.reload_app': { text: 'Reload App' },
    'error.boundary.retry_automatic': { text: 'Auto-retry in progress...' },
    'error.boundary.retry_count': { text: 'Retry 1/3' }
  }))
}));

// Test component that throws an error
const ThrowError: React.FC<{ shouldThrow?: boolean; errorType?: string }> = ({
  shouldThrow = true,
  errorType = 'general'
}) => {
  if (shouldThrow) {
    const error = new Error(`Test ${errorType} error`);
    error.name = errorType === 'network' ? 'NetworkError' : 'Error';
    throw error;
  }
  return <div>No error</div>;
};

// Test component for testing RTK Query integration
const RTKQueryTestComponent: React.FC<{ error?: any; query?: string }> = ({
  error,
  query = 'testQuery'
}) => {
  const { handleRTKError } = useRTKQueryErrorIntegration();

  // This would typically be called when RTK Query returns an error
  handleRTKError(error, query);

  return <div>RTK Query component</div>;
};

describe('ContentErrorBoundary', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Mock console methods to avoid noise in tests
    vi.spyOn(console, 'error').mockImplementation(() => {});
    vi.spyOn(console, 'warn').mockImplementation(() => {});
    vi.spyOn(console, 'group').mockImplementation(() => {});
    vi.spyOn(console, 'groupEnd').mockImplementation(() => {});
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('Basic Error Boundary Functionality', () => {
    it('renders children when no error occurs', () => {
      render(
        <ContentErrorBoundary>
          <ThrowError shouldThrow={false} />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('No error')).toBeInTheDocument();
    });

    it('catches and displays error when child component throws', () => {
      render(
        <ContentErrorBoundary>
          <ThrowError errorType="general" />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('Content System Error')).toBeInTheDocument();
      expect(screen.getByText('Something went wrong with the content system.')).toBeInTheDocument();
    });

    it('renders custom fallback when provided', () => {
      const customFallback = <div>Custom error message</div>;

      render(
        <ContentErrorBoundary fallback={customFallback}>
          <ThrowError />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('Custom error message')).toBeInTheDocument();
      expect(screen.queryByText('Content System Error')).not.toBeInTheDocument();
    });
  });

  describe('Error Classification', () => {
    it('classifies network errors correctly', () => {
      render(
        <ContentErrorBoundary>
          <ThrowError errorType="network" />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('Network Connection Error')).toBeInTheDocument();
    });

    it('shows appropriate error messages for different error types', async () => {
      const { rerender } = render(
        <ContentErrorBoundary>
          <ThrowError shouldThrow={false} />
        </ContentErrorBoundary>
      );

      // Test network error
      rerender(
        <ContentErrorBoundary>
          <ThrowError errorType="network" />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('Network Connection Error')).toBeInTheDocument();
    });
  });

  describe('Retry Mechanism', () => {
    it('shows retry button for recoverable errors', () => {
      render(
        <ContentErrorBoundary maxRetries={3}>
          <ThrowError errorType="network" />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('Try Again')).toBeInTheDocument();
    });

    it('handles manual retry attempts', async () => {
      const { rerender } = render(
        <ContentErrorBoundary maxRetries={3}>
          <ThrowError errorType="network" />
        </ContentErrorBoundary>
      );

      const retryButton = screen.getByText('Try Again');
      fireEvent.click(retryButton);

      // Wait for retry logic
      await waitFor(() => {
        expect(screen.getByText('Retry 1/3')).toBeInTheDocument();
      });
    });

    it('disables retry when max retries reached', async () => {
      render(
        <ContentErrorBoundary maxRetries={1}>
          <ThrowError errorType="network" />
        </ContentErrorBoundary>
      );

      const retryButton = screen.getByText('Try Again');
      fireEvent.click(retryButton);

      await waitFor(() => {
        expect(screen.queryByText('Try Again')).not.toBeInTheDocument();
      });
    });
  });

  describe('Reset Functionality', () => {
    it('resets error state when reset button is clicked', async () => {
      const { rerender } = render(
        <ContentErrorBoundary>
          <ThrowError />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('Content System Error')).toBeInTheDocument();

      const resetButton = screen.getByText('Reset Component');
      fireEvent.click(resetButton);

      // After reset, render a non-throwing component
      rerender(
        <ContentErrorBoundary>
          <ThrowError shouldThrow={false} />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('No error')).toBeInTheDocument();
    });
  });

  describe('Error Reporting', () => {
    it('calls onError callback when error occurs', () => {
      const onError = vi.fn();

      render(
        <ContentErrorBoundary onError={onError}>
          <ThrowError />
        </ContentErrorBoundary>
      );

      expect(onError).toHaveBeenCalledWith(
        expect.any(Error),
        expect.objectContaining({
          componentStack: expect.any(String)
        })
      );
    });

    it('calls reportError callback with ContentError object', () => {
      const reportError = vi.fn();

      render(
        <ContentErrorBoundary reportError={reportError}>
          <ThrowError />
        </ContentErrorBoundary>
      );

      expect(reportError).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.any(String),
          message: expect.any(String),
          originalError: expect.any(Error),
          timestamp: expect.any(Number),
          isRecoverable: expect.any(Boolean)
        })
      );
    });
  });

  describe('Network Status Integration', () => {
    it('shows network status for network errors', () => {
      // Mock navigator.onLine
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: false
      });

      render(
        <ContentErrorBoundary>
          <ThrowError errorType="network" />
        </ContentErrorBoundary>
      );

      expect(screen.getByText('Network: Offline')).toBeInTheDocument();
    });

    it('auto-retries when network comes back online', async () => {
      // Mock navigator.onLine as false initially
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: false
      });

      render(
        <ContentErrorBoundary maxRetries={3}>
          <ThrowError errorType="network" />
        </ContentErrorBoundary>
      );

      // Simulate network coming back online
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: true
      });

      // Trigger online event
      act(() => {
        window.dispatchEvent(new Event('online'));
      });

      await waitFor(() => {
        expect(screen.getByText('Auto-retry in progress...')).toBeInTheDocument();
      }, { timeout: 2000 });
    });
  });
});

describe('FastContentErrorBoundary', () => {
  it('renders minimal error UI', () => {
    render(
      <FastContentErrorBoundary>
        <ThrowError />
      </FastContentErrorBoundary>
    );

    expect(screen.getByText('Content temporarily unavailable')).toBeInTheDocument();
  });

  it('uses custom fallback when provided', () => {
    const customFallback = <div>Fast error fallback</div>;

    render(
      <FastContentErrorBoundary fallback={customFallback}>
        <ThrowError />
      </FastContentErrorBoundary>
    );

    expect(screen.getByText('Fast error fallback')).toBeInTheDocument();
  });

  it('calls onError callback', () => {
    const onError = vi.fn();

    render(
      <FastContentErrorBoundary onError={onError}>
        <ThrowError />
      </FastContentErrorBoundary>
    );

    expect(onError).toHaveBeenCalledWith(expect.any(Error));
  });
});

describe('ContentErrorProvider and Context', () => {
  const TestConsumer: React.FC = () => {
    const { reportError, errorHistory, hasRecentErrors, clearErrors } = useContentErrorContext();

    const handleReportError = () => {
      reportError({
        type: ContentErrorType.NETWORK,
        message: 'Test error',
        originalError: new Error('Test'),
        timestamp: Date.now(),
        retryCount: 0,
        isRecoverable: true
      });
    };

    return (
      <div>
        <button onClick={handleReportError}>Report Error</button>
        <button onClick={clearErrors}>Clear Errors</button>
        <div>Error Count: {errorHistory.length}</div>
        <div>Has Recent Errors: {hasRecentErrors.toString()}</div>
      </div>
    );
  };

  it('provides error context to children', () => {
    render(
      <ContentErrorProvider>
        <TestConsumer />
      </ContentErrorProvider>
    );

    expect(screen.getByText('Error Count: 0')).toBeInTheDocument();
    expect(screen.getByText('Has Recent Errors: false')).toBeInTheDocument();
  });

  it('manages error history correctly', () => {
    render(
      <ContentErrorProvider>
        <TestConsumer />
      </ContentErrorProvider>
    );

    const reportButton = screen.getByText('Report Error');
    fireEvent.click(reportButton);

    expect(screen.getByText('Error Count: 1')).toBeInTheDocument();
    expect(screen.getByText('Has Recent Errors: true')).toBeInTheDocument();
  });

  it('clears error history', () => {
    render(
      <ContentErrorProvider>
        <TestConsumer />
      </ContentErrorProvider>
    );

    const reportButton = screen.getByText('Report Error');
    const clearButton = screen.getByText('Clear Errors');

    fireEvent.click(reportButton);
    expect(screen.getByText('Error Count: 1')).toBeInTheDocument();

    fireEvent.click(clearButton);
    expect(screen.getByText('Error Count: 0')).toBeInTheDocument();
  });
});

describe('RTK Query Integration', () => {
  it('throws error when RTK Query error is provided', () => {
    const rtkError = { message: 'Network request failed' };

    expect(() => {
      render(
        <ContentErrorBoundary>
          <RTKQueryTestComponent error={rtkError} query="fetchUsers" />
        </ContentErrorBoundary>
      );
    }).not.toThrow(); // Error boundary should catch it

    expect(screen.getByText('Data Service Error')).toBeInTheDocument();
  });

  it('does not throw when no RTK Query error', () => {
    render(
      <ContentErrorBoundary>
        <RTKQueryTestComponent query="fetchUsers" />
      </ContentErrorBoundary>
    );

    expect(screen.getByText('RTK Query component')).toBeInTheDocument();
  });
});

describe('Content Error Recovery Hook', () => {
  const TestRecoveryComponent: React.FC = () => {
    const { recoveryAttempts, isRecovering, attemptRecovery, resetRecovery } = useContentErrorRecovery();

    const handleRecovery = () => {
      attemptRecovery(async () => {
        // Simulate async recovery operation
        await new Promise(resolve => setTimeout(resolve, 100));
      });
    };

    return (
      <div>
        <div>Recovery Attempts: {recoveryAttempts}</div>
        <div>Is Recovering: {isRecovering.toString()}</div>
        <button onClick={handleRecovery}>Attempt Recovery</button>
        <button onClick={resetRecovery}>Reset Recovery</button>
      </div>
    );
  };

  it('manages recovery state correctly', async () => {
    render(<TestRecoveryComponent />);

    expect(screen.getByText('Recovery Attempts: 0')).toBeInTheDocument();
    expect(screen.getByText('Is Recovering: false')).toBeInTheDocument();

    const recoveryButton = screen.getByText('Attempt Recovery');
    fireEvent.click(recoveryButton);

    expect(screen.getByText('Is Recovering: true')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText('Is Recovering: false')).toBeInTheDocument();
    });
  });

  it('resets recovery state', () => {
    render(<TestRecoveryComponent />);

    const resetButton = screen.getByText('Reset Recovery');
    fireEvent.click(resetButton);

    expect(screen.getByText('Recovery Attempts: 0')).toBeInTheDocument();
    expect(screen.getByText('Is Recovering: false')).toBeInTheDocument();
  });
});

describe('withContentErrorRecovery HOC', () => {
  const TestComponent: React.FC<{ shouldThrow?: boolean }> = ({ shouldThrow = false }) => {
    if (shouldThrow) throw new Error('Test error');
    return <div>Test Component</div>;
  };

  it('wraps component with error boundary', () => {
    const WrappedComponent = withContentErrorRecovery(TestComponent);

    render(<WrappedComponent shouldThrow={false} />);
    expect(screen.getByText('Test Component')).toBeInTheDocument();
  });

  it('catches errors in wrapped component', () => {
    const WrappedComponent = withContentErrorRecovery(TestComponent);

    render(<WrappedComponent shouldThrow={true} />);
    expect(screen.getByText('Content System Error')).toBeInTheDocument();
  });

  it('uses fast mode when specified', () => {
    const WrappedComponent = withContentErrorRecovery(TestComponent, { fastMode: true });

    render(<WrappedComponent shouldThrow={true} />);
    expect(screen.getByText('Content temporarily unavailable')).toBeInTheDocument();
  });

  it('applies custom options', () => {
    const customFallback = <div>Custom HOC Fallback</div>;
    const WrappedComponent = withContentErrorRecovery(TestComponent, {
      fallback: customFallback,
      maxRetries: 5
    });

    render(<WrappedComponent shouldThrow={true} />);
    expect(screen.getByText('Custom HOC Fallback')).toBeInTheDocument();
  });
});