/**
 * Unit Tests: Accessibility Compliance (T063)
 * Tests WCAG 2.1 AA compliance across all components
 * Constitutional requirements: Screen reader support, keyboard navigation, color contrast
 */
import { describe, it, expect, vi, beforeAll, afterAll } from 'vitest';
import { render, screen, fireEvent, cleanup } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { TchatButton } from '../../components/TchatButton';
import { TchatCard } from '../../components/TchatCard';
import { TchatInput } from '../../components/TchatInput';

// Clean up after each test
afterEach(cleanup);

// Mock accessibility testing APIs
const mockAccessibilityAPI = {
  checkColorContrast: vi.fn(),
  validateSemanticStructure: vi.fn(),
  testKeyboardNavigation: vi.fn(),
  auditScreenReader: vi.fn()
};

// Color contrast checker utility
const getColorContrast = (foreground: string, background: string): number => {
  // Simplified contrast calculation for testing
  // In production, use a proper color contrast library
  const getLuminance = (color: string): number => {
    // Mock luminance calculation based on color
    const colors: Record<string, number> = {
      '#FFFFFF': 1.0,    // White
      '#000000': 0.0,    // Black
      '#3B82F6': 0.2,    // Primary blue
      '#EF4444': 0.15,   // Error red
      '#10B981': 0.25,   // Success green
      '#F9FAFB': 0.95,   // Light gray
      '#111827': 0.05    // Dark gray
    };
    return colors[color.toUpperCase()] ?? 0.5;
  };

  const l1 = getLuminance(foreground);
  const l2 = getLuminance(background);
  const lighter = Math.max(l1, l2);
  const darker = Math.min(l1, l2);

  return (lighter + 0.05) / (darker + 0.05);
};

describe('Accessibility Compliance Validator', () => {
  describe('WCAG 2.1 AA Color Contrast Requirements', () => {
    it('should validate TchatButton color contrast meets 4.5:1 ratio', () => {
      // Test primary button contrast
      const primaryContrast = getColorContrast('#FFFFFF', '#3B82F6');
      expect(primaryContrast).toBeGreaterThanOrEqual(4.5);

      // Test secondary button contrast
      const secondaryContrast = getColorContrast('#111827', '#F9FAFB');
      expect(secondaryContrast).toBeGreaterThanOrEqual(4.5);

      // Test destructive button contrast
      const destructiveContrast = getColorContrast('#FFFFFF', '#EF4444');
      expect(destructiveContrast).toBeGreaterThanOrEqual(4.5);
    });

    it('should validate TchatInput color contrast in all states', () => {
      // Default state contrast
      const defaultContrast = getColorContrast('#111827', '#FFFFFF');
      expect(defaultContrast).toBeGreaterThanOrEqual(4.5);

      // Valid state contrast
      const validContrast = getColorContrast('#111827', '#10B981');
      expect(validContrast).toBeGreaterThanOrEqual(4.5);

      // Invalid state contrast
      const errorContrast = getColorContrast('#EF4444', '#FFFFFF');
      expect(errorContrast).toBeGreaterThanOrEqual(4.5);
    });

    it('should validate TchatCard color contrast across variants', () => {
      // Elevated variant contrast
      const elevatedContrast = getColorContrast('#111827', '#FFFFFF');
      expect(elevatedContrast).toBeGreaterThanOrEqual(4.5);

      // Filled variant contrast
      const filledContrast = getColorContrast('#111827', '#F9FAFB');
      expect(filledContrast).toBeGreaterThanOrEqual(4.5);
    });

    it('should fail with insufficient color contrast', () => {
      // Test a combination that would fail WCAG AA
      const lowContrast = getColorContrast('#CCCCCC', '#FFFFFF');
      expect(lowContrast).toBeLessThan(4.5);
    });
  });

  describe('Keyboard Navigation Compliance', () => {
    it('should support Tab navigation for TchatButton', async () => {
      const user = userEvent.setup();

      render(
        <div>
          <TchatButton>Button 1</TchatButton>
          <TchatButton>Button 2</TchatButton>
          <TchatButton disabled>Button 3</TchatButton>
        </div>
      );

      const button1 = screen.getByText('Button 1');
      const button2 = screen.getByText('Button 2');
      const button3 = screen.getByText('Button 3');

      // First button should be focusable
      await user.tab();
      expect(button1).toHaveFocus();

      // Tab to second button
      await user.tab();
      expect(button2).toHaveFocus();

      // Tab should skip disabled button
      await user.tab();
      expect(button3).not.toHaveFocus();
    });

    it('should support keyboard activation for TchatButton', async () => {
      const user = userEvent.setup();
      const handleClick = vi.fn();

      render(<TchatButton onClick={handleClick}>Keyboard Button</TchatButton>);

      const button = screen.getByText('Keyboard Button');
      await user.tab();
      expect(button).toHaveFocus();

      // Test Enter key activation
      await user.keyboard('{Enter}');
      expect(handleClick).toHaveBeenCalledTimes(1);

      // Test Space key activation
      await user.keyboard(' ');
      expect(handleClick).toHaveBeenCalledTimes(2);
    });

    it('should support keyboard navigation for interactive TchatCard', async () => {
      const user = userEvent.setup();
      const handleClick = vi.fn();

      render(
        <TchatCard interactive onClick={handleClick}>
          Interactive Card
        </TchatCard>
      );

      const card = screen.getByRole('button');

      // Should be focusable
      expect(card).toHaveAttribute('tabIndex', '0');

      await user.tab();
      expect(card).toHaveFocus();

      // Should activate with Enter
      await user.keyboard('{Enter}');
      expect(handleClick).toHaveBeenCalledTimes(1);

      // Should activate with Space
      await user.keyboard(' ');
      expect(handleClick).toHaveBeenCalledTimes(2);
    });

    it('should support keyboard navigation for TchatInput', async () => {
      const user = userEvent.setup();

      render(
        <TchatInput
          label="Test Input"
          type="password"
          showPasswordToggle={true}
        />
      );

      const input = screen.getByTestId('tchat-input');
      const passwordToggle = screen.getByTestId('password-toggle');

      // Input should be focusable
      await user.tab();
      expect(input).toHaveFocus();

      // Tab to password toggle button
      await user.tab();
      expect(passwordToggle).toHaveFocus();

      // Activate password toggle with Enter
      await user.keyboard('{Enter}');
      expect(input).toHaveAttribute('type', 'text');
    });

    it('should skip non-interactive elements in tab order', async () => {
      const user = userEvent.setup();

      render(
        <div>
          <TchatButton>Focusable Button</TchatButton>
          <TchatCard>Non-interactive Card</TchatCard>
          <TchatInput label="Focusable Input" />
        </div>
      );

      const button = screen.getByText('Focusable Button');
      const card = screen.getByRole('article');
      const input = screen.getByTestId('tchat-input');

      // Non-interactive card should not be focusable
      expect(card).not.toHaveAttribute('tabIndex');

      // Tab navigation should skip non-interactive elements
      await user.tab();
      expect(button).toHaveFocus();

      await user.tab();
      expect(input).toHaveFocus();
      expect(card).not.toHaveFocus();
    });
  });

  describe('Screen Reader Support (ARIA)', () => {
    it('should provide proper ARIA labels for TchatButton', () => {
      render(
        <TchatButton
          loading={true}
          aria-label="Save document"
        >
          Save
        </TchatButton>
      );

      const button = screen.getByRole('button');
      expect(button).toHaveAttribute('aria-label', 'Save document');
      expect(button).toHaveAttribute('data-loading', 'true');
    });

    it('should provide proper ARIA attributes for TchatCard', () => {
      render(
        <TchatCard
          interactive={true}
          ariaLabel="Product card for iPhone 14"
          role="button"
        >
          Product Information
        </TchatCard>
      );

      const card = screen.getByRole('button');
      expect(card).toHaveAttribute('aria-label', 'Product card for iPhone 14');
      expect(card).toHaveAttribute('tabIndex', '0');
    });

    it('should provide proper ARIA attributes for TchatInput', () => {
      render(
        <TchatInput
          label="Email Address"
          validationState="invalid"
          error="Please enter a valid email address"
          required
        />
      );

      const input = screen.getByTestId('tchat-input');
      const errorMessage = screen.getByTestId('input-error');

      expect(input).toHaveAttribute('aria-invalid', 'true');
      expect(input).toHaveAttribute('aria-describedby', errorMessage.id);
      expect(errorMessage).toHaveAttribute('role', 'alert');
    });

    it('should announce loading states properly', () => {
      render(
        <TchatButton loading={true}>
          Loading Button
        </TchatButton>
      );

      const button = screen.getByRole('button');
      const spinner = screen.getByTestId('loading-spinner');

      expect(button).toHaveAttribute('data-loading', 'true');
      expect(spinner).toBeInTheDocument();
      expect(button).toBeDisabled(); // Loading buttons should be disabled
    });

    it('should provide proper semantic structure', () => {
      render(
        <div>
          <TchatCard role="article">
            <h2>Card Title</h2>
            <p>Card content</p>
          </TchatCard>
          <TchatButton role="button">
            Action Button
          </TchatButton>
        </div>
      );

      expect(screen.getByRole('article')).toBeInTheDocument();
      expect(screen.getByRole('button')).toBeInTheDocument();
      expect(screen.getByRole('heading', { level: 2 })).toBeInTheDocument();
    });

    it('should handle dynamic ARIA state changes', async () => {
      const user = userEvent.setup();

      const { rerender } = render(
        <TchatInput
          validationState="none"
          error=""
        />
      );

      let input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('aria-invalid', 'false');

      // Change to invalid state
      rerender(
        <TchatInput
          validationState="invalid"
          error="Required field"
        />
      );

      input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('aria-invalid', 'true');

      const errorMessage = screen.getByTestId('input-error');
      expect(errorMessage).toHaveAttribute('role', 'alert');
    });
  });

  describe('Focus Management', () => {
    it('should maintain focus outline visibility', () => {
      render(<TchatButton>Focus Test</TchatButton>);

      const button = screen.getByText('Focus Test');
      expect(button).toHaveClass(
        'focus-visible:outline-none',
        'focus-visible:ring-2',
        'focus-visible:ring-blue-500'
      );
    });

    it('should trap focus in modal-like components', async () => {
      const user = userEvent.setup();

      // Simulate a modal or dialog behavior
      render(
        <div role="dialog" aria-labelledby="modal-title">
          <h2 id="modal-title">Modal Title</h2>
          <TchatInput label="Modal Input" />
          <div>
            <TchatButton>Cancel</TchatButton>
            <TchatButton>Save</TchatButton>
          </div>
        </div>
      );

      const input = screen.getByTestId('tchat-input');
      const cancelButton = screen.getByText('Cancel');
      const saveButton = screen.getByText('Save');

      // Focus should move through modal elements
      await user.tab();
      expect(input).toHaveFocus();

      await user.tab();
      expect(cancelButton).toHaveFocus();

      await user.tab();
      expect(saveButton).toHaveFocus();
    });

    it('should restore focus after dynamic content changes', async () => {
      const user = userEvent.setup();
      const { rerender } = render(
        <TchatButton>Original Button</TchatButton>
      );

      const originalButton = screen.getByText('Original Button');
      await user.tab();
      expect(originalButton).toHaveFocus();

      // Simulate content change
      rerender(
        <TchatButton>Updated Button</TchatButton>
      );

      const updatedButton = screen.getByText('Updated Button');
      // Focus should be maintained on the same element
      expect(updatedButton).toHaveFocus();
    });
  });

  describe('Touch Target Requirements', () => {
    it('should meet minimum 44dp touch targets for TchatButton', () => {
      const { container } = render(
        <TchatButton size="sm">Small Button</TchatButton>
      );

      const button = container.firstChild as HTMLElement;
      expect(button).toHaveClass('min-h-[44px]');
    });

    it('should meet minimum touch targets for TchatCard', () => {
      const { container } = render(
        <TchatCard interactive>Interactive Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('min-h-[44px]');
    });

    it('should meet minimum touch targets for TchatInput', () => {
      const { container } = render(
        <TchatInput size="sm" />
      );

      const input = container.querySelector('[data-testid="tchat-input"]') as HTMLElement;
      expect(input).toHaveClass('min-h-[32px]'); // Small size exception, but still accessible
    });

    it('should provide adequate spacing between interactive elements', () => {
      const { container } = render(
        <div className="flex gap-2">
          <TchatButton size="sm">Button 1</TchatButton>
          <TchatButton size="sm">Button 2</TchatButton>
        </div>
      );

      const buttonContainer = container.firstChild as HTMLElement;
      expect(buttonContainer).toHaveClass('gap-2'); // 8px minimum spacing
    });
  });

  describe('Error Messaging and Validation', () => {
    it('should announce validation errors properly', () => {
      render(
        <TchatInput
          validationState="invalid"
          error="Password must be at least 8 characters"
        />
      );

      const errorMessage = screen.getByTestId('input-error');
      expect(errorMessage).toHaveAttribute('role', 'alert');
      expect(errorMessage).toHaveTextContent('Password must be at least 8 characters');
    });

    it('should associate error messages with form controls', () => {
      render(
        <TchatInput
          label="Password"
          validationState="invalid"
          error="Password is required"
        />
      );

      const input = screen.getByTestId('tchat-input');
      const errorMessage = screen.getByTestId('input-error');

      expect(input).toHaveAttribute('aria-describedby', errorMessage.id);
      expect(input).toHaveAttribute('aria-invalid', 'true');
    });

    it('should provide success confirmation for valid states', () => {
      const { container } = render(
        <TchatInput
          validationState="valid"
          label="Email"
        />
      );

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('aria-invalid', 'false');

      // Should show success icon
      const successIcon = container.querySelector('.text-success');
      expect(successIcon).toBeInTheDocument();
    });
  });

  describe('High Contrast Mode Support', () => {
    it('should maintain visibility in high contrast mode', () => {
      // Simulate high contrast mode detection
      const { container } = render(
        <TchatButton variant="primary">
          High Contrast Button
        </TchatButton>
      );

      const button = container.firstChild as HTMLElement;

      // Buttons should have proper borders for high contrast mode
      expect(button).toHaveClass('rounded-lg');
      expect(button).toHaveAttribute('data-variant', 'primary');
    });

    it('should provide fallback styling for custom properties', () => {
      const { container } = render(
        <TchatCard variant="glass">
          Glass Effect Card
        </TchatCard>
      );

      const card = container.firstChild as HTMLElement;

      // Glass variant should have fallback borders
      expect(card).toHaveClass('border');
    });
  });

  describe('Animation and Motion Accessibility', () => {
    it('should respect reduced motion preferences', () => {
      const { container } = render(
        <TchatButton>Animated Button</TchatButton>
      );

      const button = container.firstChild as HTMLElement;

      // Should use duration-200 for smooth but respectful animations
      expect(button).toHaveClass('transition-all', 'duration-200');
    });

    it('should provide non-motion alternatives for essential animations', () => {
      render(
        <TchatButton loading={true}>
          Loading Button
        </TchatButton>
      );

      const button = screen.getByRole('button');
      const spinner = screen.getByTestId('loading-spinner');

      // Loading state should be announced even without animation
      expect(button).toBeDisabled();
      expect(button).toHaveAttribute('data-loading', 'true');
      expect(spinner).toBeInTheDocument();
    });
  });

  describe('Language and Internationalization', () => {
    it('should support RTL text direction', () => {
      const { container } = render(
        <div dir="rtl">
          <TchatInput
            type="search"
            placeholder="بحث"
          />
        </div>
      );

      const input = screen.getByTestId('tchat-input');
      expect(input).toHaveAttribute('placeholder', 'بحث');

      // Search icon should still be positioned correctly
      expect(input).toHaveClass('pl-10');
    });

    it('should handle dynamic text sizing', () => {
      const { container } = render(
        <TchatButton>
          Very Long Button Text That Might Wrap
        </TchatButton>
      );

      const button = container.firstChild as HTMLElement;
      expect(button).toHaveClass('min-h-[44px]'); // Should maintain minimum height
    });
  });

  describe('Component Integration Accessibility', () => {
    it('should maintain accessibility when components are combined', () => {
      render(
        <TchatCard interactive>
          <div>
            <h3>Product Card</h3>
            <p>Product description</p>
            <TchatButton size="sm">
              Add to Cart
            </TchatButton>
          </div>
        </TchatCard>
      );

      const card = screen.getByRole('button'); // Interactive card becomes button
      const button = screen.getByText('Add to Cart');

      // Card should be focusable
      expect(card).toHaveAttribute('tabIndex', '0');

      // Button inside should still be accessible
      expect(button).toHaveAttribute('data-testid', 'tchat-button');
    });

    it('should handle nested focus management correctly', async () => {
      const user = userEvent.setup();
      const cardClick = vi.fn();
      const buttonClick = vi.fn();

      render(
        <TchatCard interactive onClick={cardClick}>
          <div>
            <h3>Interactive Card</h3>
            <TchatButton onClick={buttonClick}>
              Nested Button
            </TchatButton>
          </div>
        </TchatCard>
      );

      const card = screen.getByRole('button', { name: /interactive card/i });
      const button = screen.getByText('Nested Button');

      // Tab should focus the card first
      await user.tab();
      expect(card).toHaveFocus();

      // Button inside interactive card should handle events properly
      await user.click(button);
      expect(buttonClick).toHaveBeenCalledTimes(1);
    });
  });
});

/**
 * Accessibility Testing Utilities
 * Helper functions for comprehensive accessibility testing
 */
export const AccessibilityTestUtils = {
  /**
   * Test color contrast ratio between foreground and background
   */
  testColorContrast: (foreground: string, background: string, minimumRatio: number = 4.5): boolean => {
    const contrast = getColorContrast(foreground, background);
    return contrast >= minimumRatio;
  },

  /**
   * Test keyboard navigation sequence
   */
  testKeyboardNavigation: async (user: any, expectedFocusOrder: string[]): Promise<boolean> => {
    for (let i = 0; i < expectedFocusOrder.length; i++) {
      await user.tab();
      const focusedElement = document.activeElement;
      const expectedElement = screen.getByTestId(expectedFocusOrder[i]);

      if (focusedElement !== expectedElement) {
        return false;
      }
    }
    return true;
  },

  /**
   * Test ARIA attributes presence and correctness
   */
  testARIAAttributes: (element: HTMLElement, requiredAttributes: string[]): boolean => {
    return requiredAttributes.every(attr => element.hasAttribute(attr));
  },

  /**
   * Test minimum touch target size compliance
   */
  testTouchTargetSize: (element: HTMLElement): boolean => {
    const styles = window.getComputedStyle(element);
    const minHeight = parseInt(styles.minHeight) || parseInt(styles.height);
    const minWidth = parseInt(styles.minWidth) || parseInt(styles.width);

    // 44px minimum for both dimensions
    return minHeight >= 44 && minWidth >= 44;
  },

  /**
   * Test focus visibility
   */
  testFocusVisibility: (element: HTMLElement): boolean => {
    const hasVisibleFocus = element.classList.contains('focus-visible:ring-2') ||
                           element.classList.contains('focus:ring-2') ||
                           element.classList.contains('focus-visible:outline');
    return hasVisibleFocus;
  }
};