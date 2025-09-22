/**
 * Input Component Tests
 * Testing the Input component with all states and interactions
 */

import { render, screen, fireEvent } from '@testing-library/react';
import { Input } from './input';
import { vi } from 'vitest';

describe('Input Component Tests', () => {
  describe('Basic Rendering', () => {
    test('renders input element', () => {
      render(<Input placeholder="Enter text" />);

      const input = screen.getByPlaceholderText(/enter text/i);
      expect(input).toBeInTheDocument();
    });

    test('applies custom className', () => {
      render(<Input className="custom-input" />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveClass('custom-input');
    });

    test('forwards ref correctly', () => {
      const ref = vi.fn();
      render(<Input ref={ref} />);

      expect(ref).toHaveBeenCalled();
    });

    test('has proper base styles', () => {
      render(<Input />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveClass('flex', 'h-9', 'w-full', 'rounded-md');
    });
  });

  describe('Input Types', () => {
    test('defaults to text type', () => {
      render(<Input />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('type', 'text');
    });

    test('supports email type', () => {
      render(<Input type="email" placeholder="Email" />);

      const input = screen.getByPlaceholderText(/email/i);
      expect(input).toHaveAttribute('type', 'email');
    });

    test('supports password type', () => {
      render(<Input type="password" placeholder="Password" />);

      const input = screen.getByPlaceholderText(/password/i);
      expect(input).toHaveAttribute('type', 'password');
    });

    test('supports number type', () => {
      render(<Input type="number" placeholder="Age" />);

      const input = screen.getByPlaceholderText(/age/i);
      expect(input).toHaveAttribute('type', 'number');
    });

    test('supports search type', () => {
      render(<Input type="search" placeholder="Search" />);

      const input = screen.getByRole('searchbox');
      expect(input).toHaveAttribute('type', 'search');
    });

    test('supports tel type', () => {
      render(<Input type="tel" placeholder="Phone" />);

      const input = screen.getByPlaceholderText(/phone/i);
      expect(input).toHaveAttribute('type', 'tel');
    });

    test('supports url type', () => {
      render(<Input type="url" placeholder="Website" />);

      const input = screen.getByPlaceholderText(/website/i);
      expect(input).toHaveAttribute('type', 'url');
    });

    test('supports date type', () => {
      render(<Input type="date" />);

      const input = screen.getByDisplayValue('');
      expect(input).toHaveAttribute('type', 'date');
    });
  });

  describe('Value and Change Handling', () => {
    test('displays initial value', () => {
      render(<Input defaultValue="Initial value" />);

      const input = screen.getByDisplayValue('Initial value');
      expect(input).toBeInTheDocument();
    });

    test('handles controlled value', () => {
      const { rerender } = render(<Input value="Controlled" readOnly />);

      const input = screen.getByRole('textbox') as HTMLInputElement;
      expect(input.value).toBe('Controlled');

      rerender(<Input value="Updated" readOnly />);
      expect(input.value).toBe('Updated');
    });

    test('handles onChange event', () => {
      const handleChange = vi.fn();
      render(<Input onChange={handleChange} />);

      const input = screen.getByRole('textbox');
      fireEvent.change(input, { target: { value: 'New text' } });

      expect(handleChange).toHaveBeenCalledTimes(1);
      expect(handleChange).toHaveBeenCalledWith(
        expect.objectContaining({
          target: expect.objectContaining({ value: 'New text' })
        })
      );
    });

    test('handles onFocus event', () => {
      const handleFocus = vi.fn();
      render(<Input onFocus={handleFocus} />);

      const input = screen.getByRole('textbox');
      fireEvent.focus(input);

      expect(handleFocus).toHaveBeenCalledTimes(1);
    });

    test('handles onBlur event', () => {
      const handleBlur = vi.fn();
      render(<Input onBlur={handleBlur} />);

      const input = screen.getByRole('textbox');
      fireEvent.blur(input);

      expect(handleBlur).toHaveBeenCalledTimes(1);
    });

    test('handles onKeyDown event', () => {
      const handleKeyDown = vi.fn();
      render(<Input onKeyDown={handleKeyDown} />);

      const input = screen.getByRole('textbox');
      fireEvent.keyDown(input, { key: 'Enter' });

      expect(handleKeyDown).toHaveBeenCalledTimes(1);
      expect(handleKeyDown).toHaveBeenCalledWith(
        expect.objectContaining({
          key: 'Enter'
        })
      );
    });
  });

  describe('States', () => {
    test('handles disabled state', () => {
      render(<Input disabled placeholder="Disabled input" />);

      const input = screen.getByPlaceholderText(/disabled input/i);
      expect(input).toBeDisabled();
      expect(input).toHaveClass('disabled:cursor-not-allowed', 'disabled:opacity-50');
    });

    test('handles readOnly state', () => {
      render(<Input readOnly value="Read only" />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('readonly');
    });

    test('does not accept input when disabled', () => {
      const { container } = render(<Input disabled defaultValue="Initial" />);

      const input = container.querySelector('input[disabled]') as HTMLInputElement;
      expect(input).toBeDisabled();
      expect(input.value).toBe('Initial');

      // Verify the input has disabled styling
      expect(input).toHaveClass('disabled:cursor-not-allowed');
      expect(input).toHaveClass('disabled:opacity-50');
    });

    test('handles required state', () => {
      render(<Input required />);

      const input = screen.getByRole('textbox');
      expect(input).toBeRequired();
    });
  });

  describe('Attributes', () => {
    test('supports placeholder', () => {
      render(<Input placeholder="Enter your name" />);

      const input = screen.getByPlaceholderText(/enter your name/i);
      expect(input).toBeInTheDocument();
    });

    test('supports name attribute', () => {
      render(<Input name="username" />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('name', 'username');
    });

    test('supports id attribute', () => {
      render(<Input id="email-input" />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('id', 'email-input');
    });

    test('supports maxLength', () => {
      render(<Input maxLength={10} />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('maxLength', '10');
    });

    test('supports minLength', () => {
      render(<Input minLength={3} />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('minLength', '3');
    });

    test('supports pattern', () => {
      render(<Input pattern="[0-9]{3}" />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('pattern', '[0-9]{3}');
    });

    test('supports autoComplete', () => {
      render(<Input autoComplete="email" />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('autoComplete', 'email');
    });

    test('supports autoFocus', () => {
      render(<Input autoFocus />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveFocus();
    });
  });

  describe('Accessibility', () => {
    test('supports aria-label', () => {
      render(<Input aria-label="Email address" />);

      const input = screen.getByRole('textbox', { name: /email address/i });
      expect(input).toBeInTheDocument();
    });

    test('supports aria-describedby', () => {
      render(
        <>
          <Input aria-describedby="email-help" />
          <span id="email-help">Enter a valid email address</span>
        </>
      );

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('aria-describedby', 'email-help');
    });

    test('supports aria-invalid', () => {
      render(<Input aria-invalid="true" />);

      const input = screen.getByRole('textbox');
      expect(input).toHaveAttribute('aria-invalid', 'true');
      expect(input).toHaveClass('aria-invalid:border-destructive');
    });

    test('has proper focus styles', () => {
      render(<Input />);

      const input = screen.getByRole('textbox');
      // Check for the actual focus styles in the component
      expect(input).toHaveClass('focus-visible:border-ring');
      expect(input).toHaveClass('focus-visible:ring-ring/50');
      expect(input).toHaveClass('focus-visible:ring-[3px]');
    });
  });

  describe('File Input', () => {
    test('supports file type', () => {
      const { container } = render(<Input type="file" />);

      const input = container.querySelector('input[type="file"]');
      expect(input).toBeInTheDocument();
    });

    test('supports accept attribute for file input', () => {
      const { container } = render(<Input type="file" accept="image/*" />);

      const input = container.querySelector('input[type="file"]');
      expect(input).toHaveAttribute('accept', 'image/*');
    });

    test('supports multiple file selection', () => {
      const { container } = render(<Input type="file" multiple />);

      const input = container.querySelector('input[type="file"]');
      expect(input).toHaveAttribute('multiple');
    });
  });

  describe('Form Integration', () => {
    test('participates in form submission', () => {
      const handleSubmit = vi.fn((e) => e.preventDefault());

      render(
        <form onSubmit={handleSubmit}>
          <Input name="testInput" defaultValue="test value" />
          <button type="submit">Submit</button>
        </form>
      );

      const button = screen.getByRole('button', { name: /submit/i });
      fireEvent.click(button);

      expect(handleSubmit).toHaveBeenCalledTimes(1);
    });

    test('respects form validation', () => {
      render(
        <form>
          <Input required />
          <button type="submit">Submit</button>
        </form>
      );

      const input = screen.getByRole('textbox') as HTMLInputElement;
      expect(input.validity.valid).toBe(false);
      expect(input.validity.valueMissing).toBe(true);
    });
  });
});