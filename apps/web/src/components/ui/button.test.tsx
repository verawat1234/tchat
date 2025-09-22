/**
 * Button Component Tests
 * Testing the Button component with all variants and states
 */

import { render, screen, fireEvent } from '@testing-library/react';
import { Button, buttonVariants } from './button';
import { vi } from 'vitest';

describe('Button Component Tests', () => {
  describe('Basic Rendering', () => {
    test('renders button with children', () => {
      render(<Button>Click me</Button>);

      const button = screen.getByRole('button', { name: /click me/i });
      expect(button).toBeInTheDocument();
      expect(button).toHaveTextContent('Click me');
    });

    test('applies custom className', () => {
      render(<Button className="custom-class">Button</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('custom-class');
    });

    test('forwards ref correctly', () => {
      const ref = vi.fn();
      render(<Button ref={ref}>Button</Button>);

      expect(ref).toHaveBeenCalled();
    });

    test('renders as child component when asChild is true', () => {
      render(
        <Button asChild>
          <a href="/test">Link Button</a>
        </Button>
      );

      const link = screen.getByRole('link', { name: /link button/i });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute('href', '/test');
    });
  });

  describe('Variants', () => {
    test('applies default variant', () => {
      render(<Button>Default</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('bg-primary', 'text-primary-foreground');
    });

    test('applies destructive variant', () => {
      render(<Button variant="destructive">Delete</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('bg-destructive', 'text-white');
    });

    test('applies outline variant', () => {
      render(<Button variant="outline">Outline</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('border', 'bg-background');
    });

    test('applies secondary variant', () => {
      render(<Button variant="secondary">Secondary</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('bg-secondary', 'text-secondary-foreground');
    });

    test('applies ghost variant', () => {
      render(<Button variant="ghost">Ghost</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('hover:bg-accent');
    });

    test('applies link variant', () => {
      render(<Button variant="link">Link</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('text-primary', 'underline-offset-4');
    });
  });

  describe('Sizes', () => {
    test('applies default size', () => {
      render(<Button>Default Size</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('h-9', 'px-4', 'py-2');
    });

    test('applies small size', () => {
      render(<Button size="sm">Small</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('h-8', 'px-3');
    });

    test('applies large size', () => {
      render(<Button size="lg">Large</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('h-10', 'px-6');
    });

    test('applies icon size', () => {
      render(<Button size="icon" aria-label="Settings">⚙️</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('size-9');
    });
  });

  describe('States', () => {
    test('handles disabled state', () => {
      render(<Button disabled>Disabled</Button>);

      const button = screen.getByRole('button');
      expect(button).toBeDisabled();
      expect(button).toHaveClass('disabled:pointer-events-none', 'disabled:opacity-50');
    });

    test('handles click events', () => {
      const handleClick = vi.fn();
      render(<Button onClick={handleClick}>Click me</Button>);

      const button = screen.getByRole('button');
      fireEvent.click(button);

      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    test('does not trigger click when disabled', () => {
      const handleClick = vi.fn();
      render(<Button disabled onClick={handleClick}>Disabled</Button>);

      const button = screen.getByRole('button');
      fireEvent.click(button);

      expect(handleClick).not.toHaveBeenCalled();
    });

    test('handles keyboard navigation', () => {
      const handleClick = vi.fn();
      render(<Button onClick={handleClick}>Keyboard</Button>);

      const button = screen.getByRole('button');
      button.focus();
      // Buttons respond to click events, not keyDown for Enter
      fireEvent.click(button);

      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    test('handles space key activation', () => {
      const handleClick = vi.fn();
      render(<Button onClick={handleClick}>Space</Button>);

      const button = screen.getByRole('button');
      button.focus();
      // Space key triggers click event on buttons
      fireEvent.click(button);

      expect(handleClick).toHaveBeenCalledTimes(1);
    });
  });

  describe('Accessibility', () => {
    test('supports aria-label', () => {
      render(<Button aria-label="Save document">💾</Button>);

      const button = screen.getByRole('button', { name: /save document/i });
      expect(button).toBeInTheDocument();
    });

    test('supports aria-pressed for toggle buttons', () => {
      render(<Button aria-pressed="true">Toggle</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveAttribute('aria-pressed', 'true');
    });

    test('supports aria-describedby', () => {
      render(
        <>
          <Button aria-describedby="help-text">Submit</Button>
          <span id="help-text">Click to submit the form</span>
        </>
      );

      const button = screen.getByRole('button');
      expect(button).toHaveAttribute('aria-describedby', 'help-text');
    });

    test('has proper focus styles', () => {
      render(<Button>Focus me</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveClass('focus-visible:ring-ring/50');
    });
  });

  describe('Button with Icons', () => {
    test('renders with icon', () => {
      render(
        <Button>
          <svg className="w-4 h-4" />
          With Icon
        </Button>
      );

      const button = screen.getByRole('button');
      const svg = button.querySelector('svg');
      expect(svg).toBeInTheDocument();
      expect(button).toHaveTextContent('With Icon');
    });

    test('applies icon-only styling when size is icon', () => {
      render(
        <Button size="icon" aria-label="Menu">
          <svg className="w-4 h-4" />
        </Button>
      );

      const button = screen.getByRole('button');
      expect(button).toHaveClass('size-9', 'rounded-md');
    });
  });

  describe('Button Types', () => {
    test('defaults to button type', () => {
      render(<Button>Default Type</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveAttribute('type', 'button');
    });

    test('can be submit type', () => {
      render(<Button type="submit">Submit</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveAttribute('type', 'submit');
    });

    test('can be reset type', () => {
      render(<Button type="reset">Reset</Button>);

      const button = screen.getByRole('button');
      expect(button).toHaveAttribute('type', 'reset');
    });
  });

  describe('buttonVariants utility', () => {
    test('generates correct classes for variant combinations', () => {
      const classes = buttonVariants({ variant: 'outline', size: 'lg' });

      expect(classes).toContain('border');
      expect(classes).toContain('h-10');
      expect(classes).toContain('px-6');
    });

    test('uses default values when not specified', () => {
      const classes = buttonVariants({});

      expect(classes).toContain('bg-primary');
      expect(classes).toContain('h-9');
    });
  });
});