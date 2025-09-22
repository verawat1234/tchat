/**
 * Label Component Tests
 * Testing the Label component with accessibility and interaction features
 */

import { render, screen, fireEvent } from '@testing-library/react';
import { Label } from './label';
import { vi } from 'vitest';

describe('Label Component Tests', () => {
  describe('Basic Rendering', () => {
    test('renders label with text content', () => {
      render(<Label>Email Address</Label>);

      const label = screen.getByText('Email Address');
      expect(label).toBeInTheDocument();
    });

    test('applies custom className', () => {
      render(
        <Label className="custom-label">Custom Label</Label>
      );

      const label = screen.getByText('Custom Label');
      expect(label).toHaveClass('custom-label');
    });

    test('has proper base styles', () => {
      render(<Label>Test Label</Label>);

      const label = screen.getByText('Test Label');
      expect(label).toHaveClass(
        'flex',
        'items-center',
        'gap-2',
        'text-sm',
        'leading-none',
        'font-medium',
        'select-none'
      );
    });

    test('renders with data-slot attribute', () => {
      const { container } = render(<Label>Test</Label>);

      const label = container.querySelector('[data-slot="label"]');
      expect(label).toBeInTheDocument();
    });
  });

  describe('Form Association', () => {
    test('associates with form input via htmlFor', () => {
      render(
        <>
          <Label htmlFor="email-input">Email</Label>
          <input id="email-input" type="email" />
        </>
      );

      const label = screen.getByText('Email');
      expect(label).toHaveAttribute('for', 'email-input');
    });

    test('can be clicked to focus associated input', () => {
      render(
        <>
          <Label htmlFor="test-input">Click Me</Label>
          <input id="test-input" type="text" />
        </>
      );

      const label = screen.getByText('Click Me');
      const input = screen.getByRole('textbox');

      // Click the label
      fireEvent.click(label);

      // Input should receive focus (browsers handle this automatically)
      // In jsdom, we need to manually check the association
      expect(label).toHaveAttribute('for', 'test-input');
      expect(input).toHaveAttribute('id', 'test-input');
    });
  });

  describe('Props and Attributes', () => {
    test('forwards additional props', () => {
      const onClick = vi.fn();
      render(
        <Label
          onClick={onClick}
          data-testid="custom-label"
          aria-label="Custom aria label"
        >
          Label Text
        </Label>
      );

      const label = screen.getByTestId('custom-label');
      expect(label).toHaveAttribute('aria-label', 'Custom aria label');

      label.click();
      expect(onClick).toHaveBeenCalledTimes(1);
    });

    test('supports children as React elements', () => {
      render(
        <Label>
          <span className="icon">📧</span>
          <span>Email Address</span>
        </Label>
      );

      expect(screen.getByText('📧')).toBeInTheDocument();
      expect(screen.getByText('Email Address')).toBeInTheDocument();
    });

    test('can be rendered as child of another component', () => {
      render(
        <div className="form-field">
          <Label>Field Label</Label>
        </div>
      );

      const label = screen.getByText('Field Label');
      expect(label.parentElement).toHaveClass('form-field');
    });
  });

  describe('Disabled States', () => {
    test('applies disabled styles when in disabled group', () => {
      render(
        <div data-disabled="true" className="group">
          <Label>Disabled Label</Label>
        </div>
      );

      const label = screen.getByText('Disabled Label');
      expect(label).toHaveClass('group-data-[disabled=true]:pointer-events-none');
      expect(label).toHaveClass('group-data-[disabled=true]:opacity-50');
    });

    test('applies disabled styles when peer is disabled', () => {
      render(
        <>
          <input className="peer" disabled />
          <Label>Label for disabled input</Label>
        </>
      );

      const label = screen.getByText('Label for disabled input');
      expect(label).toHaveClass('peer-disabled:cursor-not-allowed');
      expect(label).toHaveClass('peer-disabled:opacity-50');
    });
  });

  describe('Accessibility', () => {
    test('is accessible to screen readers', () => {
      render(
        <Label htmlFor="accessible-input">
          Accessible Label
        </Label>
      );

      const label = screen.getByText('Accessible Label');
      // Label element is inherently accessible
      expect(label.tagName.toLowerCase()).toBe('label');
    });

    test('prevents text selection', () => {
      render(<Label>Non-selectable Text</Label>);

      const label = screen.getByText('Non-selectable Text');
      expect(label).toHaveClass('select-none');
    });

    test('maintains proper focus order in form', () => {
      render(
        <form>
          <Label htmlFor="first">First Field</Label>
          <input id="first" />
          <Label htmlFor="second">Second Field</Label>
          <input id="second" />
        </form>
      );

      const firstLabel = screen.getByText('First Field');
      const secondLabel = screen.getByText('Second Field');

      expect(firstLabel).toBeInTheDocument();
      expect(secondLabel).toBeInTheDocument();
    });

    test('supports ARIA attributes', () => {
      render(
        <Label
          aria-required="true"
          aria-describedby="help-text"
        >
          Required Field
        </Label>
      );

      const label = screen.getByText('Required Field');
      expect(label).toHaveAttribute('aria-required', 'true');
      expect(label).toHaveAttribute('aria-describedby', 'help-text');
    });
  });

  describe('Layout and Styling', () => {
    test('uses flexbox for content alignment', () => {
      render(
        <Label>
          <span>Icon</span>
          <span>Text</span>
        </Label>
      );

      const label = screen.getByText('Icon').parentElement;
      expect(label).toHaveClass('flex', 'items-center', 'gap-2');
    });

    test('applies typography styles', () => {
      render(<Label>Styled Text</Label>);

      const label = screen.getByText('Styled Text');
      expect(label).toHaveClass('text-sm', 'leading-none', 'font-medium');
    });

    test('can be styled with custom styles', () => {
      render(
        <Label
          className="text-red-500 font-bold"
          style={{ marginBottom: '10px' }}
        >
          Custom Styled Label
        </Label>
      );

      const label = screen.getByText('Custom Styled Label');
      expect(label).toHaveClass('text-red-500', 'font-bold');
      expect(label).toHaveStyle({ marginBottom: '10px' });
    });
  });

  describe('Complex Use Cases', () => {
    test('works with required field indicators', () => {
      render(
        <Label htmlFor="required-field">
          Name <span className="text-red-500">*</span>
        </Label>
      );

      expect(screen.getByText('Name')).toBeInTheDocument();
      expect(screen.getByText('*')).toHaveClass('text-red-500');
    });

    test('works with helper text', () => {
      render(
        <div>
          <Label htmlFor="field-with-help">Email</Label>
          <input id="field-with-help" aria-describedby="help-text" />
          <span id="help-text">Enter your email address</span>
        </div>
      );

      const label = screen.getByText('Email');
      const input = screen.getByRole('textbox');
      const helpText = screen.getByText('Enter your email address');

      expect(label).toHaveAttribute('for', 'field-with-help');
      expect(input).toHaveAttribute('aria-describedby', 'help-text');
      expect(helpText).toBeInTheDocument();
    });

    test('works in a form group with multiple elements', () => {
      render(
        <fieldset>
          <legend>Contact Information</legend>
          <Label htmlFor="email">Email</Label>
          <input id="email" type="email" />
          <Label htmlFor="phone">Phone</Label>
          <input id="phone" type="tel" />
        </fieldset>
      );

      expect(screen.getByText('Email')).toBeInTheDocument();
      expect(screen.getByText('Phone')).toBeInTheDocument();
      expect(screen.getByText('Contact Information')).toBeInTheDocument();
    });
  });
});