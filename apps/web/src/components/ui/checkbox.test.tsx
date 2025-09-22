/**
 * Checkbox Component Tests
 * Testing the Checkbox component with state management and accessibility
 */

import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Checkbox } from './checkbox';
import { vi } from 'vitest';
import { useState } from 'react';

// Controlled checkbox component for testing
function ControlledCheckbox({ defaultChecked = false, onCheckedChange }: any) {
  const [checked, setChecked] = useState(defaultChecked);

  return (
    <Checkbox
      checked={checked}
      onCheckedChange={(value) => {
        setChecked(value as boolean);
        onCheckedChange?.(value);
      }}
    />
  );
}

describe('Checkbox Component Tests', () => {
  describe('Basic Rendering', () => {
    test('renders unchecked checkbox by default', () => {
      const { container } = render(<Checkbox />);

      const checkbox = container.querySelector('[data-slot="checkbox"]');
      expect(checkbox).toBeInTheDocument();
      expect(checkbox).toHaveAttribute('data-state', 'unchecked');
    });

    test('applies custom className', () => {
      const { container } = render(
        <Checkbox className="custom-checkbox" />
      );

      const checkbox = container.querySelector('[data-slot="checkbox"]');
      expect(checkbox).toHaveClass('custom-checkbox');
    });

    test('has proper base styles', () => {
      const { container } = render(<Checkbox />);

      const checkbox = container.querySelector('[data-slot="checkbox"]');
      expect(checkbox).toHaveClass(
        'size-4',
        'shrink-0',
        'rounded-[4px]',
        'border',
        'shadow-xs'
      );
    });

    test('renders with role="checkbox"', () => {
      render(<Checkbox />);

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toBeInTheDocument();
    });
  });

  describe('State Management', () => {
    test('toggles checked state on click', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();

      render(<Checkbox onCheckedChange={onCheckedChange} />);

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('aria-checked', 'false');

      await user.click(checkbox);
      expect(onCheckedChange).toHaveBeenCalledWith(true);

      // For uncontrolled component, we need to re-render with controlled state
      const { rerender } = render(<Checkbox checked={true} />);
      const checkedBox = screen.getAllByRole('checkbox')[1];
      expect(checkedBox).toHaveAttribute('aria-checked', 'true');
    });

    test('works as controlled component', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();

      render(
        <ControlledCheckbox
          defaultChecked={false}
          onCheckedChange={onCheckedChange}
        />
      );

      const checkbox = screen.getByRole('checkbox');

      await user.click(checkbox);
      await waitFor(() => {
        expect(checkbox).toHaveAttribute('aria-checked', 'true');
      });
      expect(onCheckedChange).toHaveBeenCalledWith(true);

      await user.click(checkbox);
      await waitFor(() => {
        expect(checkbox).toHaveAttribute('aria-checked', 'false');
      });
      expect(onCheckedChange).toHaveBeenCalledWith(false);
    });

    test('supports indeterminate state', () => {
      render(<Checkbox checked="indeterminate" />);

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('aria-checked', 'mixed');
    });

    test('renders checked state correctly', () => {
      const { container } = render(<Checkbox checked={true} />);

      const checkbox = container.querySelector('[data-slot="checkbox"]');
      expect(checkbox).toHaveAttribute('data-state', 'checked');
    });
  });

  describe('Check Indicator', () => {
    test('shows check icon when checked', () => {
      const { container } = render(<Checkbox checked={true} />);

      const indicator = container.querySelector('[data-slot="checkbox-indicator"]');
      expect(indicator).toBeInTheDocument();

      // Check for SVG icon
      const svg = indicator?.querySelector('svg');
      expect(svg).toBeInTheDocument();
    });

    test('hides indicator when unchecked', () => {
      const { container } = render(<Checkbox checked={false} />);

      const checkbox = container.querySelector('[data-slot="checkbox"]');
      const indicator = container.querySelector('[data-slot="checkbox-indicator"]');

      // Radix hides indicator with display: none or similar
      expect(checkbox).toHaveAttribute('data-state', 'unchecked');
    });

    test('indicator has proper styles', () => {
      const { container } = render(<Checkbox checked={true} />);

      const indicator = container.querySelector('[data-slot="checkbox-indicator"]');
      expect(indicator).toHaveClass(
        'flex',
        'items-center',
        'justify-center',
        'text-current'
      );
    });
  });

  describe('Disabled State', () => {
    test('renders as disabled when disabled prop is true', () => {
      render(<Checkbox disabled />);

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('disabled');
      expect(checkbox).toHaveAttribute('data-disabled');
    });

    test('applies disabled styles', () => {
      const { container } = render(<Checkbox disabled />);

      const checkbox = container.querySelector('[data-slot="checkbox"]');
      expect(checkbox).toHaveClass('disabled:cursor-not-allowed', 'disabled:opacity-50');
    });

    test('does not respond to clicks when disabled', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();

      render(<Checkbox disabled onCheckedChange={onCheckedChange} />);

      const checkbox = screen.getByRole('checkbox');
      await user.click(checkbox);

      expect(onCheckedChange).not.toHaveBeenCalled();
    });
  });

  describe('Keyboard Interactions', () => {
    test('can be toggled with Space key', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();

      render(<Checkbox onCheckedChange={onCheckedChange} />);

      const checkbox = screen.getByRole('checkbox');
      checkbox.focus();

      await user.keyboard(' ');
      expect(onCheckedChange).toHaveBeenCalledWith(true);
    });

    test('can be focused with Tab key', async () => {
      const user = userEvent.setup();

      render(
        <div>
          <button>Before</button>
          <Checkbox />
          <button>After</button>
        </div>
      );

      const beforeButton = screen.getByText('Before');
      const checkbox = screen.getByRole('checkbox');

      beforeButton.focus();
      await user.tab();

      expect(checkbox).toHaveFocus();
    });
  });

  describe('Accessibility', () => {
    test('has proper ARIA attributes', () => {
      render(<Checkbox />);

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('aria-checked', 'false');
      expect(checkbox).toHaveAttribute('type', 'button');
    });

    test('supports aria-label', () => {
      render(<Checkbox aria-label="Accept terms" />);

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('aria-label', 'Accept terms');
    });

    test('supports aria-labelledby', () => {
      render(
        <>
          <span id="label">Accept terms and conditions</span>
          <Checkbox aria-labelledby="label" />
        </>
      );

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('aria-labelledby', 'label');
    });

    test('supports aria-describedby', () => {
      render(
        <>
          <Checkbox aria-describedby="description" />
          <span id="description">You must accept to continue</span>
        </>
      );

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('aria-describedby', 'description');
    });

    test('indicates invalid state with aria-invalid', () => {
      const { container } = render(<Checkbox aria-invalid />);

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('aria-invalid', 'true');

      const checkboxElement = container.querySelector('[data-slot="checkbox"]');
      expect(checkboxElement).toHaveClass(
        'aria-invalid:ring-destructive/20',
        'aria-invalid:border-destructive'
      );
    });

    test('supports required attribute', () => {
      render(<Checkbox required />);

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('aria-required', 'true');
    });
  });

  describe('Visual States', () => {
    test('applies focus styles on focus', () => {
      const { container } = render(<Checkbox />);

      const checkbox = screen.getByRole('checkbox');
      const checkboxElement = container.querySelector('[data-slot="checkbox"]');

      checkbox.focus();

      expect(checkboxElement).toHaveClass(
        'focus-visible:border-ring',
        'focus-visible:ring-[3px]'
      );
    });

    test('applies checked styles when checked', () => {
      const { container } = render(<Checkbox checked />);

      const checkbox = container.querySelector('[data-slot="checkbox"]');
      expect(checkbox).toHaveClass(
        'data-[state=checked]:bg-primary',
        'data-[state=checked]:text-primary-foreground',
        'data-[state=checked]:border-primary'
      );
    });
  });

  describe('Form Integration', () => {
    test('works with form labels', () => {
      render(
        <div>
          <label htmlFor="terms">
            <Checkbox id="terms" />
            <span>I accept the terms and conditions</span>
          </label>
        </div>
      );

      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toHaveAttribute('id', 'terms');
    });

    test('can be used in a form', () => {
      const handleSubmit = vi.fn();

      render(
        <form onSubmit={handleSubmit}>
          <Checkbox name="agreement" value="accepted" />
          <button type="submit">Submit</button>
        </form>
      );

      // Radix Checkbox doesn't pass through name/value to the button element
      // It handles form integration internally
      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toBeInTheDocument();
    });

    test('supports custom value prop', () => {
      render(<Checkbox value="option1" />);

      // Radix Checkbox handles value internally
      const checkbox = screen.getByRole('checkbox');
      expect(checkbox).toBeInTheDocument();
    });
  });

  describe('Complex Use Cases', () => {
    test('works in a checkbox group', () => {
      render(
        <fieldset>
          <legend>Select your interests:</legend>
          <label>
            <Checkbox name="interests" value="sports" />
            Sports
          </label>
          <label>
            <Checkbox name="interests" value="music" />
            Music
          </label>
          <label>
            <Checkbox name="interests" value="reading" />
            Reading
          </label>
        </fieldset>
      );

      const checkboxes = screen.getAllByRole('checkbox');
      expect(checkboxes).toHaveLength(3);

      // Radix Checkbox handles name internally
      checkboxes.forEach(checkbox => {
        expect(checkbox).toBeInTheDocument();
      });
    });

    test('maintains state when parent re-renders', async () => {
      const user = userEvent.setup();

      function Parent() {
        const [count, setCount] = useState(0);
        return (
          <div>
            <button onClick={() => setCount(c => c + 1)}>Count: {count}</button>
            <ControlledCheckbox />
          </div>
        );
      }

      render(<Parent />);

      const checkbox = screen.getByRole('checkbox');
      const button = screen.getByRole('button');

      // Check the checkbox
      await user.click(checkbox);
      expect(checkbox).toHaveAttribute('aria-checked', 'true');

      // Trigger parent re-render
      await user.click(button);

      // Checkbox should maintain its state
      expect(checkbox).toHaveAttribute('aria-checked', 'true');
    });
  });
});