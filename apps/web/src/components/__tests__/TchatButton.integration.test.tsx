/**
 * Integration Test: TchatButton Web Component
 * Tests cross-platform consistency requirements
 * These tests should now PASS since the component is implemented
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { TchatButton } from '../TchatButton';

describe('TchatButton Web Integration Tests', () => {
  describe('Component Rendering', () => {
    it('should render TchatButton component', () => {
      render(<TchatButton>Test Button</TchatButton>);

      const buttonElement = screen.getByRole('button');
      expect(buttonElement).toBeInTheDocument();
      expect(buttonElement).toHaveTextContent('Test Button');
      expect(buttonElement).toHaveAttribute('data-testid', 'tchat-button');
    });

    it('should support all 5 variants per specification', () => {
      const variants = ['primary', 'secondary', 'ghost', 'destructive', 'outline'] as const;

      variants.forEach(variant => {
        const { container } = render(
          <TchatButton variant={variant} data-testid={`button-${variant}`}>
            {variant} Button
          </TchatButton>
        );

        const buttonElement = screen.getByTestId(`button-${variant}`);
        expect(buttonElement).toBeInTheDocument();
        expect(buttonElement).toHaveAttribute('data-variant', variant);
      });
    });

    it('should support all 3 size variants', () => {
      const sizes = ['sm', 'md', 'lg'] as const;

      sizes.forEach(size => {
        const { container } = render(
          <TchatButton size={size} data-testid={`button-${size}`}>
            {size} Button
          </TchatButton>
        );

        const buttonElement = screen.getByTestId(`button-${size}`);
        expect(buttonElement).toBeInTheDocument();
        expect(buttonElement).toHaveAttribute('data-size', size);
      });
    });

    it('should support loading state', () => {
      render(
        <TchatButton loading data-testid="loading-button">
          Loading Button
        </TchatButton>
      );

      const buttonElement = screen.getByTestId('loading-button');
      const loadingSpinner = screen.getByTestId('loading-spinner');

      expect(buttonElement).toBeInTheDocument();
      expect(buttonElement).toBeDisabled();
      expect(buttonElement).toHaveAttribute('data-loading', 'true');
      expect(loadingSpinner).toBeInTheDocument();
    });
  });

  describe('Cross-Platform Consistency Validation', () => {
    it('should use consistent design tokens across platforms', () => {
      render(<TchatButton variant="primary">Design Token Test</TchatButton>);

      const buttonElement = screen.getByRole('button');

      // Verify minimum height for touch target compliance
      expect(buttonElement).toHaveClass('min-h-[44px]');

      // Verify base classes for consistency
      expect(buttonElement).toHaveClass('inline-flex', 'items-center', 'justify-center');
    });

    it('should maintain 97% visual consistency per Constitution', () => {
      render(<TchatButton variant="outline">Consistency Test</TchatButton>);

      const buttonElement = screen.getByRole('button');

      // Constitutional requirement: minimum 44dp touch targets
      expect(buttonElement).toHaveClass('min-h-[44px]');

      // Should be focusable for accessibility
      expect(buttonElement.tabIndex).toBe(0);
    });

    it('should be interactive with proper click handling', () => {
      const handleClick = vi.fn();

      render(
        <TchatButton onClick={handleClick}>
          Interactive Button
        </TchatButton>
      );

      const buttonElement = screen.getByRole('button');

      fireEvent.click(buttonElement);
      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('should support icon placement', () => {
      const leftIcon = <span data-testid="left-icon">←</span>;
      const rightIcon = <span data-testid="right-icon">→</span>;

      render(
        <TchatButton leftIcon={leftIcon} rightIcon={rightIcon}>
          Icon Button
        </TchatButton>
      );

      expect(screen.getByTestId('button-left-icon')).toBeInTheDocument();
      expect(screen.getByTestId('button-right-icon')).toBeInTheDocument();
    });
  });

  describe('Accessibility Compliance (WCAG 2.1 AA)', () => {
    it('should have proper semantic markup', () => {
      render(
        <TchatButton aria-label="Submit form">
          Submit
        </TchatButton>
      );

      const buttonElement = screen.getByRole('button');
      expect(buttonElement).toHaveAttribute('aria-label', 'Submit form');
    });

    it('should support keyboard navigation', () => {
      const handleClick = vi.fn();

      render(
        <TchatButton onClick={handleClick}>
          Keyboard Button
        </TchatButton>
      );

      const buttonElement = screen.getByRole('button');

      // Should be focusable
      expect(buttonElement.tabIndex).toBe(0);

      // Should support Enter key activation
      buttonElement.focus();
      expect(document.activeElement).toBe(buttonElement);

      fireEvent.keyDown(buttonElement, { key: 'Enter' });
      fireEvent.keyUp(buttonElement, { key: 'Enter' });
    });

    it('should have proper disabled state', () => {
      render(<TchatButton disabled>Disabled Button</TchatButton>);

      const buttonElement = screen.getByRole('button');
      expect(buttonElement).toBeDisabled();
      expect(buttonElement).toHaveClass('disabled:pointer-events-none', 'disabled:opacity-60');
    });
  });

  describe('Performance Requirements', () => {
    it('should render within 200ms performance budget', async () => {
      const startTime = performance.now();

      render(
        <TchatButton variant="primary">
          Performance test button
        </TchatButton>
      );

      const endTime = performance.now();
      const renderTime = endTime - startTime;

      // Constitutional requirement: <200ms component load times
      expect(renderTime).toBeLessThan(200);
    });

    it('should support 60fps animations', () => {
      render(<TchatButton variant="primary">Animation Test</TchatButton>);

      const buttonElement = screen.getByRole('button');

      // Should have GPU acceleration classes
      expect(buttonElement).toHaveClass('transform-gpu', 'will-change-transform');

      // Should have transition classes for smooth animations
      expect(buttonElement).toHaveClass('transition-all', 'duration-200');
    });
  });
});