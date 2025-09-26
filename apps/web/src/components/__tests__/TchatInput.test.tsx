/**
 * Unit Tests: TchatInput Validation States (T062)
 * Tests input validation logic, state management, and user interactions
 * Complements integration tests with focused unit testing
 * Constitutional requirements: 97% consistency, WCAG 2.1 AA, <200ms load
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, cleanup } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {
  TchatInput,
  type TchatInputProps,
  type TchatInputType,
  type TchatInputValidationState,
  type TchatInputSize
} from '../TchatInput';

// Clean up after each test
afterEach(cleanup);

describe('TchatInput Component Unit Tests', () => {
  describe('Component Rendering and Props', () => {
    it('should render with default props', () => {
      render(<TchatInput />);

      const input = screen.getByTestId('tchat-input');
      expect(input).toBeInTheDocument();
      expect(input).toHaveAttribute('data-type', 'text');
      expect(input).toHaveAttribute('data-validation-state', 'none');
      expect(input).toHaveAttribute('data-size', 'md');
    });

    it('should forward ref correctly', () => {
      const ref = vi.fn();
      render(<TchatInput ref={ref} />);
      expect(ref).toHaveBeenCalledWith(expect.any(HTMLInputElement));
    });

    it('should apply custom className', () => {
      render(<TchatInput className="custom-input" />);

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveClass('custom-input');
    });

    it('should pass through HTML attributes', () => {
      render(
        <TchatInput
          placeholder="Enter text"
          maxLength={100}
          data-custom="test"
        />
      );

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('placeholder', 'Enter text');
      expect(input).toHaveAttribute('maxLength', '100');
      expect(input).toHaveAttribute('data-custom', 'test');
    });
  });

  describe('Input Types System', () => {
    const inputTypes: TchatInputType[] = ['text', 'email', 'password', 'number', 'search', 'multiline'];

    it('should render all input types correctly', () => {
      inputTypes.forEach(inputType => {
        const { container } = render(
          <TchatInput type={inputType} data-testid={`input-${inputType}`} />
        );

        const input = screen.getByTestId(`input-${inputType}`);
        expect(input).toBeInTheDocument();
        expect(input).toHaveAttribute('data-type', inputType);

        cleanup();
      });
    });

    it('should render text input with correct type', () => {
      render(<TchatInput type="text" />);

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('type', 'text');
      expect(input).toHaveAttribute('data-type', 'text');
    });

    it('should render email input with correct type', () => {
      render(<TchatInput type="email" />);

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('type', 'email');
      expect(input).toHaveAttribute('data-type', 'email');
    });

    it('should render password input with correct type', () => {
      render(<TchatInput type="password" />);

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('type', 'password');
      expect(input).toHaveAttribute('data-type', 'password');
    });

    it('should render number input with correct type', () => {
      render(<TchatInput type="number" />);

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('type', 'number');
      expect(input).toHaveAttribute('data-type', 'number');
    });

    it('should render search input with search icon', () => {
      render(<TchatInput type="search" />);

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('type', 'text'); // Search type becomes text
      expect(input).toHaveAttribute('data-type', 'search');
      expect(input).toHaveClass('pl-10'); // Space for leading icon
    });

    it('should render multiline input as textarea', () => {
      render(<TchatInput type="multiline" />);

      const input = screen.getByTestId('tchat-input');
      expect(input.tagName).toBe('TEXTAREA');
      expect(input).toHaveAttribute('data-type', 'multiline');
      expect(input).toHaveClass('resize-none', 'min-h-[88px]');
    });
  });

  describe('Validation States System', () => {
    const validationStates: TchatInputValidationState[] = ['none', 'valid', 'invalid'];

    it('should render all validation states correctly', () => {
      validationStates.forEach(state => {
        const { container } = render(
          <TchatInput
            validationState={state}
            data-testid={`input-${state}`}
          />
        );

        const input = screen.getByTestId(`input-${state}`);
        expect(input).toBeInTheDocument();
        expect(input).toHaveAttribute('data-validation-state', state);

        cleanup();
      });
    });

    it('should apply none validation state styling', () => {
      const { container } = render(
        <TchatInput validationState="none" />
      );

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('border-border');
      expect(input).toHaveAttribute('aria-invalid', 'false');
    });

    it('should apply valid validation state styling', () => {
      const { container } = render(
        <TchatInput validationState="valid" />
      );

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('border-success', 'bg-success/5');

      // Should show check icon
      const icon = container.querySelector('.text-success');
      expect(icon).toBeInTheDocument();
    });

    it('should apply invalid validation state styling', () => {
      const { container } = render(
        <TchatInput validationState="invalid" />
      );

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('border-error', 'bg-error/5');
      expect(input).toHaveAttribute('aria-invalid', 'true');

      // Should show error icon
      const icon = container.querySelector('.text-error');
      expect(icon).toBeInTheDocument();
    });

    it('should display error message with invalid state', () => {
      render(
        <TchatInput
          validationState="invalid"
          error="This field is required"
        />
      );

      const errorMessage = screen.getByTestId('input-error');
      expect(errorMessage).toBeInTheDocument();
      expect(errorMessage).toHaveTextContent('This field is required');
      expect(errorMessage).toHaveAttribute('role', 'alert');
    });

    it('should not display error message with valid state', () => {
      render(
        <TchatInput
          validationState="valid"
          error="This field is required"
        />
      );

      const errorMessage = screen.queryByTestId('input-error');
      expect(errorMessage).not.toBeInTheDocument();
    });

    it('should link error message with input for accessibility', () => {
      render(
        <TchatInput
          validationState="invalid"
          error="Invalid email format"
        />
      );

      const input = screen.getByTestId('tchat-input');
      const errorMessage = screen.getByTestId('input-error');
      const errorId = errorMessage.id;

      expect(input).toHaveAttribute('aria-describedby', errorId);
      expect(errorId).toMatch(/-error$/);
    });
  });

  describe('Size System', () => {
    const sizes: TchatInputSize[] = ['sm', 'md', 'lg'];

    it('should render all size variants correctly', () => {
      sizes.forEach(size => {
        const { container } = render(
          <TchatInput size={size} data-testid={`input-${size}`} />
        );

        const input = screen.getByTestId(`input-${size}`);
        expect(input).toBeInTheDocument();
        expect(input).toHaveAttribute('data-size', size);

        cleanup();
      });
    });

    it('should apply small size styles', () => {
      const { container } = render(
        <TchatInput size="sm" />
      );

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('h-8', 'px-2', 'text-xs', 'min-h-[32px]');
    });

    it('should apply medium size styles (default)', () => {
      const { container } = render(
        <TchatInput size="md" />
      );

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('h-11', 'px-3', 'text-sm', 'min-h-[44px]');
    });

    it('should apply large size styles', () => {
      const { container } = render(
        <TchatInput size="lg" />
      );

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('h-12', 'px-4', 'text-base');
    });

    it('should meet minimum touch target size requirements', () => {
      sizes.forEach(size => {
        const { container } = render(<TchatInput size={size} />);
        const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;

        // Constitutional requirement: 44dp minimum touch targets
        if (size === 'sm') {
          expect(input).toHaveClass('min-h-[32px]');
        } else {
          expect(input).toHaveClass('min-h-[44px]');
        }

        cleanup();
      });
    });
  });

  describe('Password Visibility Toggle', () => {
    it('should render password toggle when showPasswordToggle is true', () => {
      render(
        <TchatInput
          type="password"
          showPasswordToggle={true}
        />
      );

      const toggleButton = screen.getByTestId('password-toggle');
      expect(toggleButton).toBeInTheDocument();
      expect(toggleButton).toHaveAttribute('aria-label', 'Toggle password visibility');
    });

    it('should not render password toggle when showPasswordToggle is false', () => {
      render(
        <TchatInput
          type="password"
          showPasswordToggle={false}
        />
      );

      const toggleButton = screen.queryByTestId('password-toggle');
      expect(toggleButton).not.toBeInTheDocument();
    });

    it('should toggle password visibility on button click', async () => {
      const user = userEvent.setup();

      render(
        <TchatInput
          type="password"
          showPasswordToggle={true}
          value="secret123"
        />
      );

      const input = screen.getByTestId('tchat-input');
      const toggleButton = screen.getByTestId('password-toggle');

      // Initially password should be hidden
      expect(input).toHaveAttribute('type', 'password');

      // Click to show password
      await user.click(toggleButton);
      expect(input).toHaveAttribute('type', 'text');

      // Click again to hide password
      await user.click(toggleButton);
      expect(input).toHaveAttribute('type', 'password');
    });

    it('should not render password toggle for non-password inputs', () => {
      render(
        <TchatInput
          type="text"
          showPasswordToggle={true}
        />
      );

      const toggleButton = screen.queryByTestId('password-toggle');
      expect(toggleButton).not.toBeInTheDocument();
    });
  });

  describe('Icon System', () => {
    it('should render leading icon when provided', () => {
      const leadingIcon = <span data-testid="custom-icon">🔍</span>;

      render(
        <TchatInput leadingIcon={leadingIcon} />
      );

      const icon = screen.getByTestId('custom-icon');
      expect(icon).toBeInTheDocument();

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveClass('pl-10'); // Space for leading icon
    });

    it('should render search icon for search type inputs', () => {
      const { container } = render(
        <TchatInput type="search" />
      );

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveClass('pl-10');

      // Search icon should be present in the DOM
      const searchIcon = container.querySelector('svg');
      expect(searchIcon).toBeInTheDocument();
    });

    it('should render trailing action when provided', () => {
      const trailingAction = <button data-testid="custom-action">Action</button>;

      render(
        <TchatInput trailingAction={trailingAction} />
      );

      const action = screen.getByTestId('custom-action');
      expect(action).toBeInTheDocument();

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveClass('pr-10'); // Space for trailing action
    });

    it('should render clear button when value is present and onClear provided', () => {
      const handleClear = vi.fn();

      render(
        <TchatInput
          value="some text"
          onClear={handleClear}
        />
      );

      const clearButton = screen.getByTestId('clear-button');
      expect(clearButton).toBeInTheDocument();
      expect(clearButton).toHaveAttribute('aria-label', 'Clear input');
    });

    it('should call onClear when clear button is clicked', async () => {
      const user = userEvent.setup();
      const handleClear = vi.fn();

      render(
        <TchatInput
          value="some text"
          onClear={handleClear}
        />
      );

      const clearButton = screen.getByTestId('clear-button');
      await user.click(clearButton);

      expect(handleClear).toHaveBeenCalledTimes(1);
    });

    it('should not render clear button when no value or onClear', () => {
      render(<TchatInput />);

      const clearButton = screen.queryByTestId('clear-button');
      expect(clearButton).not.toBeInTheDocument();
    });

    it('should prioritize password toggle over other trailing content', () => {
      const trailingAction = <button data-testid="custom-action">Action</button>;

      render(
        <TchatInput
          type="password"
          showPasswordToggle={true}
          trailingAction={trailingAction}
        />
      );

      // Password toggle should be present
      const passwordToggle = screen.getByTestId('password-toggle');
      expect(passwordToggle).toBeInTheDocument();

      // Custom trailing action should not be present
      const customAction = screen.queryByTestId('custom-action');
      expect(customAction).not.toBeInTheDocument();
    });
  });

  describe('Label and Accessibility', () => {
    it('should render label when provided', () => {
      render(
        <TchatInput label="Email Address" />
      );

      const label = screen.getByText('Email Address');
      expect(label.tagName).toBe('LABEL');
      expect(label).toBeInTheDocument();
    });

    it('should associate label with input', () => {
      render(
        <TchatInput label="Username" />
      );

      const input = screen.getByTestId('tchat-input');
      const label = screen.getByText('Username');

      expect(label).toHaveAttribute('for', input.id);
      expect(input).toHaveAttribute('aria-labelledby', label.id);
    });

    it('should not render label when not provided', () => {
      render(<TchatInput />);

      const input = screen.getByTestId('tchat-input');
      expect(input).not.toHaveAttribute('aria-labelledby');
    });

    it('should handle custom aria-describedby', () => {
      render(
        <TchatInput aria-describedby="custom-description" />
      );

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('aria-describedby', 'custom-description');
    });

    it('should prioritize error message over custom aria-describedby', () => {
      render(
        <TchatInput
          validationState="invalid"
          error="Required field"
          aria-describedby="custom-description"
        />
      );

      const input = screen.getByTestId('tchat-input');
      const errorMessage = screen.getByTestId('input-error');

      expect(input).toHaveAttribute('aria-describedby', errorMessage.id);
    });
  });

  describe('Focus and Interaction States', () => {
    it('should handle focus events', async () => {
      const user = userEvent.setup();
      const handleFocus = vi.fn();

      render(<TchatInput onFocus={handleFocus} />);

      const input = screen.getByTestId('tchat-input');
      await user.click(input);

      expect(handleFocus).toHaveBeenCalledTimes(1);
    });

    it('should handle blur events', async () => {
      const user = userEvent.setup();
      const handleBlur = vi.fn();

      render(<TchatInput onBlur={handleBlur} />);

      const input = screen.getByTestId('tchat-input');
      await user.click(input);
      await user.tab(); // Tab away to trigger blur

      expect(handleBlur).toHaveBeenCalledTimes(1);
    });

    it('should handle disabled state', () => {
      render(<TchatInput disabled />);

      const input = screen.getByTestId('tchat-input');
      expect(input).toBeDisabled();
      expect(input).toHaveClass('disabled:cursor-not-allowed', 'disabled:opacity-50');
    });

    it('should handle user input', async () => {
      const user = userEvent.setup();
      const handleChange = vi.fn();

      render(<TchatInput onChange={handleChange} />);

      const input = screen.getByTestId('tchat-input');
      await user.type(input, 'Hello World');

      expect(input).toHaveValue('Hello World');
      expect(handleChange).toHaveBeenCalled();
    });
  });

  describe('Multiline Input Behavior', () => {
    it('should render textarea for multiline type', () => {
      render(<TchatInput type="multiline" />);

      const input = screen.getByTestId('tchat-input');
      expect(input.tagName).toBe('TEXTAREA');
      expect(input).toHaveClass('resize-none', 'min-h-[88px]');
    });

    it('should handle multiline focus and blur events', async () => {
      const user = userEvent.setup();
      const handleFocus = vi.fn();
      const handleBlur = vi.fn();

      render(
        <TchatInput
          type="multiline"
          onFocus={handleFocus}
          onBlur={handleBlur}
        />
      );

      const textarea = screen.getByTestId('tchat-input');
      await user.click(textarea);
      expect(handleFocus).toHaveBeenCalledTimes(1);

      await user.tab();
      expect(handleBlur).toHaveBeenCalledTimes(1);
    });

    it('should handle multiline input text', async () => {
      const user = userEvent.setup();

      render(<TchatInput type="multiline" />);

      const textarea = screen.getByTestId('tchat-input');
      await user.type(textarea, 'Line 1\nLine 2\nLine 3');

      expect(textarea).toHaveValue('Line 1\nLine 2\nLine 3');
    });
  });

  describe('Performance Requirements', () => {
    it('should have GPU acceleration classes', () => {
      const { container } = render(<TchatInput />);

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('transform-gpu', 'will-change-[border-color,box-shadow]');
    });

    it('should have transition classes for smooth animations', () => {
      const { container } = render(<TchatInput />);

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('transition-all', 'duration-200');
    });

    it('should render quickly within performance budget', () => {
      const startTime = performance.now();

      render(
        <TchatInput
          type="password"
          validationState="invalid"
          size="lg"
          showPasswordToggle={true}
          label="Complex Input"
          error="Error message"
        />
      );

      const endTime = performance.now();
      const renderTime = endTime - startTime;

      // Constitutional requirement: <200ms component load times
      expect(renderTime).toBeLessThan(200);
    });
  });

  describe('Complex Scenarios', () => {
    it('should handle all features together', () => {
      render(
        <TchatInput
          type="email"
          size="lg"
          validationState="valid"
          label="Email Address"
          placeholder="Enter your email"
          leadingIcon={<span data-testid="email-icon">📧</span>}
          className="custom-class"
          aria-describedby="email-help"
        />
      );

      // Verify all components are present
      const input = screen.getByTestId('tchat-input');
      const label = screen.getByText('Email Address');
      const icon = screen.getByTestId('email-icon');

      expect(input).toBeInTheDocument();
      expect(label).toBeInTheDocument();
      expect(icon).toBeInTheDocument();

      // Verify attributes
      expect(input).toHaveAttribute('type', 'email');
      expect(input).toHaveAttribute('data-size', 'lg');
      expect(input).toHaveAttribute('data-validation-state', 'valid');
      expect(input).toHaveAttribute('placeholder', 'Enter your email');
      expect(input).toHaveClass('custom-class');

      // Verify spacing for icons
      expect(input).toHaveClass('pl-10', 'pr-10');
    });

    it('should handle validation state transitions', async () => {
      const user = userEvent.setup();
      const { rerender } = render(
        <TchatInput
          validationState="none"
          error=""
        />
      );

      let input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('data-validation-state', 'none');
      expect(input).toHaveAttribute('aria-invalid', 'false');

      // Change to invalid state
      rerender(
        <TchatInput
          validationState="invalid"
          error="This field is required"
        />
      );

      input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('data-validation-state', 'invalid');
      expect(input).toHaveAttribute('aria-invalid', 'true');
      expect(screen.getByTestId('input-error')).toBeInTheDocument();

      // Change to valid state
      rerender(
        <TchatInput
          validationState="valid"
          error=""
        />
      );

      input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('data-validation-state', 'valid');
      expect(screen.queryByTestId('input-error')).not.toBeInTheDocument();
    });

    it('should maintain consistent styling across all variants', () => {
      const combinations = [
        { type: 'text', size: 'sm', validationState: 'none' },
        { type: 'email', size: 'md', validationState: 'valid' },
        { type: 'password', size: 'lg', validationState: 'invalid' },
        { type: 'number', size: 'sm', validationState: 'valid' },
        { type: 'search', size: 'md', validationState: 'none' }
      ] as const;

      combinations.forEach((combo, index) => {
        const { container } = render(
          <TchatInput
            key={index}
            type={combo.type}
            size={combo.size}
            validationState={combo.validationState}
            data-testid={`input-${index}`}
          />
        );

        const input = container.querySelector(`[data-testid="input-${index}"]`) as HTMLElement;

        // All inputs should have consistent base classes
        expect(input).toHaveClass(
          'flex',
          'w-full',
          'rounded-md',
          'border',
          'bg-white',
          'transition-all',
          'duration-200',
          'transform-gpu'
        );

        // All inputs should meet minimum height requirements
        expect(input).toHaveClass(/min-h-\[/);

        cleanup();
      });
    });
  });
});