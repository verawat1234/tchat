/**
 * Switch Component Tests
 * Testing the Switch component with toggle states and accessibility
 */

import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Switch } from './switch';
import { vi } from 'vitest';
import React from 'react';

describe('Switch Component Tests', () => {
  describe('Basic Rendering', () => {
    test('renders switch with default state', () => {
      const { container } = render(<Switch />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toBeInTheDocument();
      expect(switchElement).toHaveClass('inline-flex', 'h-[1.15rem]', 'w-8', 'rounded-full');
    });

    test('applies custom className', () => {
      const { container } = render(<Switch className="custom-switch" />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toHaveClass('custom-switch');
    });

    test('renders switch thumb', () => {
      const { container } = render(<Switch />);

      const thumb = container.querySelector('[data-slot="switch-thumb"]');
      expect(thumb).toBeInTheDocument();
      expect(thumb).toHaveClass('block', 'size-4', 'rounded-full');
    });

    test('has role="switch" for accessibility', () => {
      render(<Switch />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toBeInTheDocument();
    });
  });

  describe('State Management', () => {
    test('renders unchecked by default', () => {
      const { container } = render(<Switch />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toHaveAttribute('data-state', 'unchecked');

      const switchRole = screen.getByRole('switch');
      expect(switchRole).toHaveAttribute('aria-checked', 'false');
    });

    test('toggles state on click', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();
      render(<Switch onCheckedChange={onCheckedChange} />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-checked', 'false');

      await user.click(switchElement);
      expect(onCheckedChange).toHaveBeenCalledWith(true);
    });

    test('works as controlled component', () => {
      const ControlledSwitch = () => {
        const [checked, setChecked] = React.useState(false);
        return (
          <Switch
            checked={checked}
            onCheckedChange={setChecked}
          />
        );
      };

      const { rerender } = render(<ControlledSwitch />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-checked', 'false');

      fireEvent.click(switchElement);

      // Re-render to see updated state
      rerender(<ControlledSwitch />);
      expect(switchElement).toHaveAttribute('aria-checked', 'true');
    });

    test('works as uncontrolled component', async () => {
      const user = userEvent.setup();
      render(<Switch defaultChecked={false} />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-checked', 'false');

      await user.click(switchElement);

      await waitFor(() => {
        expect(switchElement).toHaveAttribute('aria-checked', 'true');
      });
    });

    test('accepts defaultChecked prop', () => {
      render(<Switch defaultChecked={true} />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-checked', 'true');
    });

    test('handles checked prop', () => {
      const { container } = render(<Switch checked={true} />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toHaveAttribute('data-state', 'checked');

      const switchRole = screen.getByRole('switch');
      expect(switchRole).toHaveAttribute('aria-checked', 'true');
    });
  });

  describe('Thumb Animation', () => {
    test('thumb position changes with state', () => {
      const { container, rerender } = render(<Switch checked={false} />);

      const thumb = container.querySelector('[data-slot="switch-thumb"]');
      expect(thumb).toHaveClass('data-[state=unchecked]:translate-x-0');

      rerender(<Switch checked={true} />);
      expect(thumb).toHaveClass('data-[state=checked]:translate-x-[calc(100%-2px)]');
    });

    test('thumb has transition styles', () => {
      const { container } = render(<Switch />);

      const thumb = container.querySelector('[data-slot="switch-thumb"]');
      expect(thumb).toHaveClass('transition-transform');
    });
  });

  describe('Keyboard Interactions', () => {
    test('toggles with Space key', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();
      render(<Switch onCheckedChange={onCheckedChange} />);

      const switchElement = screen.getByRole('switch');
      switchElement.focus();

      await user.keyboard(' ');
      expect(onCheckedChange).toHaveBeenCalledWith(true);
    });

    test('toggles with Enter key', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();
      render(<Switch onCheckedChange={onCheckedChange} />);

      const switchElement = screen.getByRole('switch');
      switchElement.focus();

      await user.keyboard('{Enter}');
      expect(onCheckedChange).toHaveBeenCalledWith(true);
    });

    test('maintains focus after toggle', async () => {
      const user = userEvent.setup();
      render(<Switch />);

      const switchElement = screen.getByRole('switch');
      switchElement.focus();

      expect(document.activeElement).toBe(switchElement);

      await user.keyboard(' ');

      expect(document.activeElement).toBe(switchElement);
    });

    test('can be focused with Tab key', async () => {
      const user = userEvent.setup();
      render(
        <>
          <button>Before</button>
          <Switch />
          <button>After</button>
        </>
      );

      const beforeButton = screen.getByText('Before');
      const switchElement = screen.getByRole('switch');
      const afterButton = screen.getByText('After');

      beforeButton.focus();
      await user.tab();

      expect(document.activeElement).toBe(switchElement);

      await user.tab();
      expect(document.activeElement).toBe(afterButton);
    });
  });

  describe('Disabled State', () => {
    test('renders as disabled when disabled prop is true', () => {
      const { container } = render(<Switch disabled />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toHaveClass('disabled:cursor-not-allowed', 'disabled:opacity-50');
      expect(switchElement).toHaveAttribute('data-disabled');

      const switchRole = screen.getByRole('switch');
      expect(switchRole).toBeInTheDocument();
    });

    test('does not respond to clicks when disabled', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();
      render(<Switch disabled onCheckedChange={onCheckedChange} />);

      const switchElement = screen.getByRole('switch');
      await user.click(switchElement);

      expect(onCheckedChange).not.toHaveBeenCalled();
    });

    test('does not respond to keyboard when disabled', async () => {
      const user = userEvent.setup();
      const onCheckedChange = vi.fn();
      render(<Switch disabled onCheckedChange={onCheckedChange} />);

      const switchElement = screen.getByRole('switch');
      switchElement.focus();

      await user.keyboard(' ');
      expect(onCheckedChange).not.toHaveBeenCalled();
    });
  });

  describe('Visual States', () => {
    test('applies checked background color', () => {
      const { container } = render(<Switch checked />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toHaveClass('data-[state=checked]:bg-primary');
    });

    test('applies unchecked background color', () => {
      const { container } = render(<Switch checked={false} />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toHaveClass('data-[state=unchecked]:bg-switch-background');
    });

    test('applies focus styles', () => {
      const { container } = render(<Switch />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toHaveClass(
        'focus-visible:border-ring',
        'focus-visible:ring-ring/50',
        'focus-visible:ring-[3px]'
      );
    });

    test('has transition styles', () => {
      const { container } = render(<Switch />);

      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toHaveClass('transition-all');
    });
  });

  describe('Accessibility', () => {
    test('has proper ARIA attributes', () => {
      render(<Switch />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-checked', 'false');
      expect(switchElement).toHaveAttribute('type', 'button');
    });

    test('supports aria-label', () => {
      render(<Switch aria-label="Dark mode" />);

      const switchElement = screen.getByRole('switch', { name: 'Dark mode' });
      expect(switchElement).toBeInTheDocument();
    });

    test('supports aria-labelledby', () => {
      render(
        <>
          <label id="theme-label">Theme</label>
          <Switch aria-labelledby="theme-label" />
        </>
      );

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-labelledby', 'theme-label');
    });

    test('supports aria-describedby', () => {
      render(
        <>
          <Switch aria-describedby="theme-help" />
          <span id="theme-help">Toggle between light and dark theme</span>
        </>
      );

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-describedby', 'theme-help');
    });

    test('indicates invalid state with aria-invalid', () => {
      const { container } = render(<Switch aria-invalid="true" />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-invalid', 'true');
    });

    test('supports required attribute', () => {
      render(<Switch required />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-required', 'true');
    });
  });

  describe('Form Integration', () => {
    test('works with form labels', () => {
      render(
        <div>
          <label htmlFor="notifications">
            Enable notifications
            <Switch id="notifications" />
          </label>
        </div>
      );

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('id', 'notifications');
    });

    test('can be used in a form', () => {
      const handleSubmit = vi.fn((e) => e.preventDefault());
      const { container } = render(
        <form onSubmit={handleSubmit}>
          <Switch name="newsletter" />
          <button type="submit">Submit</button>
        </form>
      );

      // Radix Switch might use a hidden input for form integration
      const switchElement = container.querySelector('[data-slot="switch"]');
      expect(switchElement).toBeInTheDocument();
      // Verify the switch can be part of a form
      const form = container.querySelector('form');
      expect(form).toContainElement(switchElement);
    });

    test('supports value prop', () => {
      render(<Switch value="on" />);

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('value', 'on');
    });
  });

  describe('Common Use Cases', () => {
    test('dark mode toggle', async () => {
      const user = userEvent.setup();
      const handleThemeChange = vi.fn();

      render(
        <Switch
          aria-label="Dark mode"
          defaultChecked={false}
          onCheckedChange={handleThemeChange}
        />
      );

      const switchElement = screen.getByRole('switch', { name: 'Dark mode' });
      expect(switchElement).toHaveAttribute('aria-checked', 'false');

      await user.click(switchElement);
      expect(handleThemeChange).toHaveBeenCalledWith(true);
    });

    test('notification settings', () => {
      render(
        <>
          <label id="notifications-label">Email notifications</label>
          <Switch
            defaultChecked={true}
            aria-labelledby="notifications-label"
            aria-describedby="notifications-desc"
          />
          <span id="notifications-desc">Receive updates via email</span>
        </>
      );

      const switchElement = screen.getByRole('switch');
      expect(switchElement).toHaveAttribute('aria-checked', 'true');
      expect(switchElement).toHaveAttribute('aria-labelledby', 'notifications-label');
      expect(switchElement).toHaveAttribute('aria-describedby', 'notifications-desc');
    });

    test('feature toggle', () => {
      const features = {
        analytics: true,
        newsletter: false,
        betaFeatures: false
      };

      const { container } = render(
        <>
          {Object.entries(features).map(([key, enabled]) => (
            <Switch
              key={key}
              defaultChecked={enabled}
              aria-label={key}
            />
          ))}
        </>
      );

      const switches = screen.getAllByRole('switch');
      expect(switches).toHaveLength(3);
      expect(switches[0]).toHaveAttribute('aria-checked', 'true');
      expect(switches[1]).toHaveAttribute('aria-checked', 'false');
      expect(switches[2]).toHaveAttribute('aria-checked', 'false');
    });
  });
});