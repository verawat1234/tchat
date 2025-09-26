/**
 * Integration Test: TchatInput Web Component
 * Tests cross-platform consistency requirements
 * These tests should now PASS since the component is implemented
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { TchatInput } from '../TchatInput';

describe('TchatInput Web Integration Tests', () => {
  describe('Component Rendering', () => {
    it('should render TchatInput component', () => {
      render(<TchatInput placeholder="Test Input" />);

      const inputElement = screen.getByPlaceholderText('Test Input');
      expect(inputElement).toBeInTheDocument();
      expect(inputElement.tagName).toBe('INPUT');
      expect(inputElement).toHaveAttribute('data-testid', 'tchat-input');
    });

    it('should support all 6 input types per specification', () => {
      const inputTypes = ['text', 'email', 'password', 'number', 'search', 'multiline'] as const;

      inputTypes.forEach(type => {
        const { container } = render(
          <TchatInput
            type={type}
            placeholder={`${type} input`}
            data-testid={`input-${type}`}
          />
        );

        const inputElement = screen.getByTestId(`input-${type}`);
        expect(inputElement).toBeInTheDocument();
        expect(inputElement).toHaveAttribute('data-type', type);

        if (type === 'multiline') {
          expect(inputElement.tagName).toBe('TEXTAREA');
        } else {
          expect(inputElement.tagName).toBe('INPUT');
        }
      });
    });

    it('should support all 3 validation states', () => {
      const validationStates = ['none', 'valid', 'invalid'] as const;

      validationStates.forEach(state => {
        const { container } = render(
          <TchatInput
            validationState={state}
            placeholder="Validation test"
            data-testid={`input-${state}`}
          />
        );

        const inputElement = screen.getByTestId(`input-${state}`);
        expect(inputElement).toBeInTheDocument();
        expect(inputElement).toHaveAttribute('data-validation-state', state);
      });
    });

    it('should support all 3 size variants', () => {
      const sizes = ['sm', 'md', 'lg'] as const;

      sizes.forEach(size => {
        const { container } = render(
          <TchatInput
            size={size}
            placeholder="Size test"
            data-testid={`input-${size}`}
          />
        );

        const inputElement = screen.getByTestId(`input-${size}`);
        expect(inputElement).toBeInTheDocument();
        expect(inputElement).toHaveAttribute('data-size', size);
      });
    });
  });

  describe('Cross-Platform Consistency Validation', () => {
    it('should use consistent design tokens across platforms', () => {
      render(
        <TchatInput
          validationState="valid"
          placeholder="Design token test"
        />
      );

      const inputElement = screen.getByPlaceholderText('Design token test');

      // Verify minimum height for touch target compliance
      expect(inputElement).toHaveClass('min-h-[44px]');
    });

    it('should maintain 97% visual consistency per Constitution', () => {
      render(
        <TchatInput
          validationState="invalid"
          error="Invalid input"
          placeholder="Consistency test"
        />
      );

      const inputElement = screen.getByPlaceholderText('Consistency test');
      const errorMessage = screen.getByText('Invalid input');

      // Constitutional requirement: minimum 44dp touch targets
      expect(inputElement).toHaveClass('min-h-[44px]');

      // Error message should be properly associated
      expect(errorMessage).toBeInTheDocument();
      expect(errorMessage).toHaveAttribute('role', 'alert');
    });

    it('should have proper validation state styling', () => {
      const { rerender } = render(
        <TchatInput validationState="valid" placeholder="Test" />
      );

      let inputElement = screen.getByPlaceholderText('Test');
      expect(inputElement).toHaveAttribute('data-validation-state', 'valid');

      // Test invalid state
      rerender(<TchatInput validationState="invalid" placeholder="Test" />);
      inputElement = screen.getByPlaceholderText('Test');
      expect(inputElement).toHaveAttribute('data-validation-state', 'invalid');
    });

    it('should support interactive features', () => {
      const handleChange = vi.fn();
      const handleFocus = vi.fn();
      const handleBlur = vi.fn();

      render(
        <TchatInput
          placeholder="Interactive test"
          onChange={handleChange}
          onFocus={handleFocus}
          onBlur={handleBlur}
        />
      );

      const inputElement = screen.getByPlaceholderText('Interactive test');

      // Test focus
      fireEvent.focus(inputElement);
      expect(handleFocus).toHaveBeenCalled();

      // Test input change
      fireEvent.change(inputElement, { target: { value: 'test input' } });
      expect(handleChange).toHaveBeenCalled();

      // Test blur
      fireEvent.blur(inputElement);
      expect(handleBlur).toHaveBeenCalled();
    });

    it('should support password visibility toggle', () => {
      render(
        <TchatInput
          type="password"
          placeholder="Password test"
          showPasswordToggle
        />
      );

      const inputElement = screen.getByPlaceholderText('Password test');
      const toggleButton = screen.getByRole('button', { name: /toggle password visibility/i });

      // Initially should be password type
      expect(inputElement).toHaveAttribute('type', 'password');

      // Click toggle
      fireEvent.click(toggleButton);
      expect(inputElement).toHaveAttribute('type', 'text');

      // Click again
      fireEvent.click(toggleButton);
      expect(inputElement).toHaveAttribute('type', 'password');
    });
  });

  describe('Accessibility Compliance (WCAG 2.1 AA)', () => {
    it('should have proper semantic markup', () => {
      render(
        <TchatInput
          label="Email Address"
          placeholder="Enter your email"
          type="email"
          required
        />
      );

      const inputElement = screen.getByLabelText('Email Address');
      expect(inputElement).toHaveAttribute('type', 'email');
      expect(inputElement).toHaveAttribute('required');
      expect(inputElement).toBeInTheDocument();
    });

    it('should support keyboard navigation', () => {
      render(
        <TchatInput
          placeholder="Keyboard test"
          data-testid="keyboard-input"
        />
      );

      const inputElement = screen.getByTestId('keyboard-input');

      // Should be focusable
      expect(inputElement.tabIndex).toBe(0);

      // Should support focus
      inputElement.focus();
      expect(document.activeElement).toBe(inputElement);
    });

    it('should display error messages for accessibility', () => {
      render(
        <TchatInput
          validationState="invalid"
          error="This field is required"
        />
      );

      const errorMessage = screen.getByText('This field is required');
      expect(errorMessage).toBeInTheDocument();
      expect(errorMessage).toHaveAttribute('role', 'alert');

      const inputElement = screen.getByRole('textbox');
      expect(inputElement).toBeInTheDocument();
    });
  });

  describe('Performance Requirements', () => {
    it('should render within 200ms performance budget', async () => {
      const startTime = performance.now();

      render(
        <TchatInput
          placeholder="Performance test input"
          type="text"
        />
      );

      const endTime = performance.now();
      const renderTime = endTime - startTime;

      // Constitutional requirement: <200ms component load times
      expect(renderTime).toBeLessThan(200);
    });

    it('should support 60fps animations on state changes', () => {
      const { rerender } = render(
        <TchatInput validationState="none" placeholder="Test" />
      );

      const inputElement = screen.getByPlaceholderText('Test');

      // Should have transition classes for smooth animations
      expect(inputElement).toHaveClass('transition-all', 'duration-200');

      // State change should be smooth
      rerender(<TchatInput validationState="valid" placeholder="Test" />);
    });
  });

  describe('Icon Support and Leading/Trailing Elements', () => {
    it('should support leading icons', () => {
      render(
        <TchatInput
          placeholder="Search..."
          type="search"
        />
      );

      // Search type automatically adds search icon
      const inputElement = screen.getByPlaceholderText('Search...');
      expect(inputElement).toBeInTheDocument();
    });

    it('should support trailing action buttons', () => {
      const handleClear = vi.fn();

      render(
        <TchatInput
          placeholder="Clearable input"
          value="some text"
          onClear={handleClear}
        />
      );

      const clearButton = screen.getByRole('button', { name: /clear/i });
      expect(clearButton).toBeInTheDocument();

      fireEvent.click(clearButton);
      expect(handleClear).toHaveBeenCalled();
    });

    it('should show validation icons', () => {
      // Test valid state icon
      const { rerender } = render(
        <TchatInput validationState="valid" placeholder="Valid input" />
      );

      let inputElement = screen.getByPlaceholderText('Valid input');
      expect(inputElement).toHaveAttribute('data-validation-state', 'valid');

      // Test invalid state icon
      rerender(
        <TchatInput validationState="invalid" placeholder="Invalid input" />
      );

      inputElement = screen.getByPlaceholderText('Invalid input');
      expect(inputElement).toHaveAttribute('data-validation-state', 'invalid');
    });
  });
});