import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {
  TextSkeleton,
  ImageSkeleton,
  CardSkeleton,
  ListSkeleton,
  TableSkeleton,
  ChatSkeleton,
  MediaSkeleton,
  SpinnerLoader,
  ProgressLoader,
  ErrorState,
  FallbackState,
} from './ContentLoader';

// Mock lucide-react icons
vi.mock('lucide-react', () => ({
  AlertCircle: () => <div data-testid="alert-circle-icon" />,
  Wifi: () => <div data-testid="wifi-icon" />,
  WifiOff: () => <div data-testid="wifi-off-icon" />,
  RefreshCw: () => <div data-testid="refresh-icon" />,
}));

describe('ContentLoader Basic Components', () => {
  describe('Skeleton Components', () => {
    it('TextSkeleton renders with default count', () => {
      render(<TextSkeleton />);
      // Should render without errors
    });

    it('TextSkeleton renders with custom count', () => {
      render(<TextSkeleton count={5} />);
      // Should render without errors
    });

    it('ImageSkeleton renders with default aspect ratio', () => {
      render(<ImageSkeleton />);
      // Should render without errors
    });

    it('ImageSkeleton renders with custom aspect ratio', () => {
      render(<ImageSkeleton aspectRatio="aspect-square" />);
      // Should render without errors
    });

    it('CardSkeleton renders', () => {
      render(<CardSkeleton />);
      // Should render without errors
    });

    it('ListSkeleton renders with default count', () => {
      render(<ListSkeleton />);
      // Should render without errors
    });

    it('ListSkeleton renders with custom count', () => {
      render(<ListSkeleton count={3} />);
      // Should render without errors
    });

    it('TableSkeleton renders with default dimensions', () => {
      render(<TableSkeleton />);
      // Should render without errors
    });

    it('TableSkeleton renders with custom dimensions', () => {
      render(<TableSkeleton rows={3} cols={4} />);
      // Should render without errors
    });

    it('ChatSkeleton renders', () => {
      render(<ChatSkeleton />);
      // Should render without errors
    });

    it('MediaSkeleton renders', () => {
      render(<MediaSkeleton />);
      // Should render without errors
    });
  });

  describe('Loading Indicators', () => {
    it('SpinnerLoader renders with default size', () => {
      render(<SpinnerLoader />);
      expect(screen.getByTestId('refresh-icon')).toBeInTheDocument();
    });

    it('SpinnerLoader renders with small size', () => {
      render(<SpinnerLoader size="sm" />);
      expect(screen.getByTestId('refresh-icon')).toBeInTheDocument();
    });

    it('SpinnerLoader renders with large size', () => {
      render(<SpinnerLoader size="lg" />);
      expect(screen.getByTestId('refresh-icon')).toBeInTheDocument();
    });

    it('ProgressLoader renders with progress value', () => {
      render(<ProgressLoader progress={50} message="Loading..." />);
      expect(screen.getByText('Loading...')).toBeInTheDocument();
      expect(screen.getByText('50%')).toBeInTheDocument();
    });

    it('ProgressLoader renders without message', () => {
      render(<ProgressLoader progress={75} />);
      expect(screen.getByText('Loading...')).toBeInTheDocument();
      expect(screen.getByText('75%')).toBeInTheDocument();
    });
  });

  describe('Error State', () => {
    it('renders error message', () => {
      render(<ErrorState error="Test error message" />);
      expect(screen.getByText('Failed to load content')).toBeInTheDocument();
      expect(screen.getByText('Test error message')).toBeInTheDocument();
      expect(screen.getByTestId('alert-circle-icon')).toBeInTheDocument();
    });

    it('renders Error object message', () => {
      const error = new Error('Network connection failed');
      render(<ErrorState error={error} />);
      expect(screen.getByText('Network connection failed')).toBeInTheDocument();
    });

    it('shows retry button when onRetry provided', () => {
      const onRetry = vi.fn();
      render(
        <ErrorState
          error="Test error"
          onRetry={onRetry}
          maxRetries={3}
          retryCount={1}
        />
      );

      const retryButton = screen.getByRole('button');
      expect(retryButton).toBeInTheDocument();
      expect(screen.getByText('Try again (2 attempts left)')).toBeInTheDocument();
    });

    it('hides retry button when max retries reached', () => {
      const onRetry = vi.fn();
      render(
        <ErrorState
          error="Test error"
          onRetry={onRetry}
          maxRetries={3}
          retryCount={3}
        />
      );

      expect(screen.queryByRole('button')).not.toBeInTheDocument();
    });

    it('calls onRetry when retry button clicked', async () => {
      const user = userEvent.setup();
      const onRetry = vi.fn();

      render(
        <ErrorState
          error="Test error"
          onRetry={onRetry}
          maxRetries={3}
          retryCount={0}
        />
      );

      const retryButton = screen.getByRole('button');
      await user.click(retryButton);
      expect(onRetry).toHaveBeenCalledTimes(1);
    });
  });

  describe('Fallback State', () => {
    it('renders online fallback state', () => {
      render(<FallbackState isOffline={false} message="Using cached data" />);
      expect(screen.getByText('Using cached data')).toBeInTheDocument();
      expect(screen.getByTestId('wifi-icon')).toBeInTheDocument();
    });

    it('renders offline fallback state', () => {
      render(<FallbackState isOffline={true} message="Offline mode active" />);
      expect(screen.getByText('Offline mode active')).toBeInTheDocument();
      expect(screen.getByTestId('wifi-off-icon')).toBeInTheDocument();
    });

    it('renders default offline message', () => {
      render(<FallbackState isOffline={true} />);
      expect(screen.getByText('Showing offline content')).toBeInTheDocument();
    });

    it('renders default online message', () => {
      render(<FallbackState isOffline={false} />);
      expect(screen.getByText('Using cached content')).toBeInTheDocument();
    });
  });

  describe('Component Props', () => {
    it('TextSkeleton accepts custom className', () => {
      const { container } = render(<TextSkeleton className="custom-class" />);
      expect(container.firstChild).toHaveClass('custom-class');
    });

    it('ImageSkeleton accepts custom className', () => {
      const { container } = render(<ImageSkeleton className="custom-class" />);
      expect(container.firstChild).toHaveClass('custom-class');
    });

    it('ProgressLoader handles undefined progress', () => {
      render(<ProgressLoader />);
      expect(screen.getByText('0%')).toBeInTheDocument();
    });

    it('ErrorState with no retry function', () => {
      render(<ErrorState error="Test error" />);
      expect(screen.queryByRole('button')).not.toBeInTheDocument();
    });
  });

  describe('Edge Cases', () => {
    it('handles empty error string', () => {
      render(<ErrorState error="" />);
      expect(screen.getByText('Failed to load content')).toBeInTheDocument();
    });

    it('handles zero progress', () => {
      render(<ProgressLoader progress={0} />);
      expect(screen.getByText('0%')).toBeInTheDocument();
    });

    it('handles progress over 100', () => {
      render(<ProgressLoader progress={150} />);
      expect(screen.getByText('150%')).toBeInTheDocument();
    });

    it('handles negative progress', () => {
      render(<ProgressLoader progress={-10} />);
      expect(screen.getByText('-10%')).toBeInTheDocument();
    });
  });
});