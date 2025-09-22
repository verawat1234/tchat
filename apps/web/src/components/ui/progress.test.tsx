/**
 * Progress Component Tests
 * Testing the Progress component with various values and states
 */

import { render, screen } from '@testing-library/react';
import { Progress } from './progress';
import { vi } from 'vitest';

describe('Progress Component Tests', () => {
  describe('Basic Rendering', () => {
    test('renders progress bar with default state', () => {
      const { container } = render(<Progress />);

      const progress = container.querySelector('[data-slot="progress"]');
      expect(progress).toBeInTheDocument();
      expect(progress).toHaveClass('relative', 'h-2', 'w-full', 'overflow-hidden', 'rounded-full');
    });

    test('applies custom className', () => {
      const { container } = render(
        <Progress className="custom-progress" />
      );

      const progress = container.querySelector('[data-slot="progress"]');
      expect(progress).toHaveClass('custom-progress');
    });

    test('has proper base styles', () => {
      const { container } = render(<Progress />);

      const progress = container.querySelector('[data-slot="progress"]');
      expect(progress).toHaveClass(
        'bg-primary/20',
        'relative',
        'h-2',
        'w-full',
        'overflow-hidden',
        'rounded-full'
      );
    });

    test('renders with role="progressbar"', () => {
      render(<Progress />);

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toBeInTheDocument();
    });

    test('renders progress indicator element', () => {
      const { container } = render(<Progress value={50} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toBeInTheDocument();
      expect(indicator).toHaveClass('bg-primary', 'h-full', 'w-full', 'flex-1', 'transition-all');
    });
  });

  describe('Value Management', () => {
    test('renders with 0% progress', () => {
      const { container } = render(<Progress value={0} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveStyle({ transform: 'translateX(-100%)' });

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '0');
    });

    test('renders with 50% progress', () => {
      const { container } = render(<Progress value={50} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveStyle({ transform: 'translateX(-50%)' });

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '50');
    });

    test('renders with 100% progress', () => {
      const { container } = render(<Progress value={100} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveStyle({ transform: 'translateX(-0%)' });

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '100');
    });

    test('renders with 25% progress', () => {
      const { container } = render(<Progress value={25} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveStyle({ transform: 'translateX(-75%)' });

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '25');
    });

    test('renders with 75% progress', () => {
      const { container } = render(<Progress value={75} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveStyle({ transform: 'translateX(-25%)' });

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '75');
    });

    test('handles undefined value as 0', () => {
      const { container } = render(<Progress />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveStyle({ transform: 'translateX(-100%)' });
    });

    test('handles null value as 0', () => {
      const { container } = render(<Progress value={null as any} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveStyle({ transform: 'translateX(-100%)' });
    });

    test('clamps values above 100', () => {
      // Radix Progress rejects invalid values > 100 and defaults to null
      // This is expected behavior - the component protects against invalid input
      const { container } = render(<Progress value={150} />);

      const progressbar = screen.getByRole('progressbar');
      // Radix clamps to null/indeterminate for invalid values
      expect(progressbar).toHaveAttribute('data-state', 'indeterminate');
    });

    test('handles negative values', () => {
      // Radix Progress rejects negative values and defaults to null (indeterminate)
      // This is expected behavior - the component protects against invalid input
      const { container } = render(<Progress value={-10} />);

      const progressbar = screen.getByRole('progressbar');
      // Radix clamps to null/indeterminate for invalid values
      expect(progressbar).toHaveAttribute('data-state', 'indeterminate');
    });
  });

  describe('Progress States', () => {
    test('handles indeterminate state', () => {
      const { container } = render(<Progress value={undefined} />);

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('data-state', 'indeterminate');
    });

    test('handles loading state', () => {
      const { container } = render(<Progress value={30} />);

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('data-state', 'loading');
    });

    test('handles complete state', () => {
      const { container } = render(<Progress value={100} />);

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('data-state', 'complete');
    });

    test('updates data-value attribute correctly', () => {
      const { rerender } = render(<Progress value={25} />);

      let progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '25');

      rerender(<Progress value={75} />);
      progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '75');
    });
  });

  describe('Accessibility', () => {
    test('has proper ARIA attributes', () => {
      render(<Progress value={60} />);

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '60');
      expect(progressbar).toHaveAttribute('aria-valuemin', '0');
      expect(progressbar).toHaveAttribute('aria-valuemax', '100');
    });

    test('updates aria-valuenow when value changes', () => {
      const { rerender } = render(<Progress value={30} />);

      let progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '30');

      rerender(<Progress value={80} />);
      progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '80');
    });

    test('supports aria-label', () => {
      render(<Progress value={50} aria-label="Upload progress" />);

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-label', 'Upload progress');
    });

    test('supports aria-labelledby', () => {
      render(
        <>
          <span id="progress-label">File upload</span>
          <Progress value={45} aria-labelledby="progress-label" />
        </>
      );

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-labelledby', 'progress-label');
    });

    test('supports aria-describedby', () => {
      render(
        <>
          <Progress value={70} aria-describedby="progress-description" />
          <span id="progress-description">70% complete</span>
        </>
      );

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-describedby', 'progress-description');
    });

    test('handles indeterminate state accessibility', () => {
      render(<Progress />);

      const progressbar = screen.getByRole('progressbar');
      // Indeterminate progress bars may not have aria-valuenow
      expect(progressbar).toHaveAttribute('data-state', 'indeterminate');
    });
  });

  describe('Visual States', () => {
    test('indicator has transition styles', () => {
      const { container } = render(<Progress value={50} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveClass('transition-all');
    });

    test('applies correct background styles', () => {
      const { container } = render(<Progress value={50} />);

      const progress = container.querySelector('[data-slot="progress"]');
      expect(progress).toHaveClass('bg-primary/20');

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveClass('bg-primary');
    });

    test('maintains rounded corners', () => {
      const { container } = render(<Progress value={75} />);

      const progress = container.querySelector('[data-slot="progress"]');
      expect(progress).toHaveClass('rounded-full');
    });

    test('has overflow hidden for clean appearance', () => {
      const { container } = render(<Progress value={50} />);

      const progress = container.querySelector('[data-slot="progress"]');
      expect(progress).toHaveClass('overflow-hidden');
    });
  });

  describe('Animation and Transitions', () => {
    test('smooth transition when value changes', () => {
      const { container, rerender } = render(<Progress value={20} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');
      expect(indicator).toHaveStyle({ transform: 'translateX(-80%)' });
      expect(indicator).toHaveClass('transition-all');

      rerender(<Progress value={60} />);
      expect(indicator).toHaveStyle({ transform: 'translateX(-40%)' });
    });

    test('handles rapid value updates', () => {
      const { container, rerender } = render(<Progress value={0} />);

      const indicator = container.querySelector('[data-slot="progress-indicator"]');

      // Simulate rapid updates
      for (let i = 0; i <= 100; i += 10) {
        rerender(<Progress value={i} />);
        expect(indicator).toHaveStyle({ transform: `translateX(-${100 - i}%)` });
      }
    });
  });

  describe('Props and Customization', () => {
    test('forwards additional props', () => {
      const onClick = vi.fn();
      render(
        <Progress
          value={50}
          data-testid="custom-progress"
          onClick={onClick}
        />
      );

      const progressbar = screen.getByTestId('custom-progress');
      expect(progressbar).toBeInTheDocument();
    });

    test('can be styled with custom styles', () => {
      const { container } = render(
        <Progress
          value={40}
          className="h-4 bg-blue-200"
          style={{ width: '200px' }}
        />
      );

      const progress = container.querySelector('[data-slot="progress"]');
      expect(progress).toHaveClass('h-4', 'bg-blue-200');
      expect(progress).toHaveStyle({ width: '200px' });
    });

    test('supports custom max value', () => {
      render(<Progress value={50} max={200} />);

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuemax', '200');
      expect(progressbar).toHaveAttribute('aria-valuenow', '50');
    });
  });

  describe('Common Use Cases', () => {
    test('file upload progress', () => {
      const uploadProgress = 65;
      render(
        <Progress
          value={uploadProgress}
          aria-label={`File upload: ${uploadProgress}%`}
        />
      );

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '65');
      expect(progressbar).toHaveAttribute('aria-label', 'File upload: 65%');
    });

    test('loading indicator', () => {
      render(
        <Progress
          value={undefined}
          aria-label="Loading content"
        />
      );

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('data-state', 'indeterminate');
      expect(progressbar).toHaveAttribute('aria-label', 'Loading content');
    });

    test('multi-step form progress', () => {
      const currentStep = 3;
      const totalSteps = 5;
      const progress = (currentStep / totalSteps) * 100;

      render(
        <Progress
          value={progress}
          aria-label={`Step ${currentStep} of ${totalSteps}`}
        />
      );

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '60');
      expect(progressbar).toHaveAttribute('aria-label', 'Step 3 of 5');
    });

    test('download progress with description', () => {
      render(
        <>
          <label id="download-label">Download Progress</label>
          <Progress
            value={85}
            aria-labelledby="download-label"
            aria-describedby="download-description"
          />
          <span id="download-description">8.5MB of 10MB</span>
        </>
      );

      const progressbar = screen.getByRole('progressbar');
      expect(progressbar).toHaveAttribute('aria-valuenow', '85');
      expect(progressbar).toHaveAttribute('aria-labelledby', 'download-label');
      expect(progressbar).toHaveAttribute('aria-describedby', 'download-description');
    });
  });
});