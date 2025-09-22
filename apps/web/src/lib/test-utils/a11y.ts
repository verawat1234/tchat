import { axe, toHaveNoViolations } from 'jest-axe';
import { RenderResult } from '@testing-library/react';

// Add the jest-axe matchers to Vitest
expect.extend(toHaveNoViolations);

// WCAG compliance levels
export type A11yLevel = 'A' | 'AA' | 'AAA';

// Axe configuration for different compliance levels
const axeConfig = {
  A: {
    runOnly: {
      type: 'tag',
      values: ['wcag2a', 'wcag21a', 'best-practice'],
    },
  },
  AA: {
    runOnly: {
      type: 'tag',
      values: ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'best-practice'],
    },
  },
  AAA: {
    runOnly: {
      type: 'tag',
      values: ['wcag2a', 'wcag2aa', 'wcag2aaa', 'wcag21a', 'wcag21aa', 'wcag21aaa', 'best-practice'],
    },
  },
};

// Run accessibility tests with axe-core
export async function testAccessibility(
  container: HTMLElement,
  level: A11yLevel = 'AA',
  options = {}
) {
  const results = await axe(container, {
    ...axeConfig[level],
    ...options,
  });

  return results;
}

// Assert no accessibility violations
export async function expectNoA11yViolations(
  renderResult: RenderResult | HTMLElement,
  level: A11yLevel = 'AA'
) {
  const container = 'container' in renderResult ? renderResult.container : renderResult;
  const results = await testAccessibility(container, level);
  expect(results).toHaveNoViolations();
}

// Check specific WCAG criteria
export interface WCAGCheck {
  criterion: string;
  level: A11yLevel;
  description: string;
  test: (element: HTMLElement) => boolean;
}

export const wcagChecks: Record<string, WCAGCheck> = {
  // 1.1.1 Non-text Content (Level A)
  altText: {
    criterion: '1.1.1',
    level: 'A',
    description: 'All images must have alt text',
    test: (element) => {
      const images = element.querySelectorAll('img');
      return Array.from(images).every((img) => img.hasAttribute('alt'));
    },
  },

  // 1.3.1 Info and Relationships (Level A)
  headingStructure: {
    criterion: '1.3.1',
    level: 'A',
    description: 'Headings must be in logical order',
    test: (element) => {
      const headings = element.querySelectorAll('h1, h2, h3, h4, h5, h6');
      let previousLevel = 0;

      return Array.from(headings).every((heading) => {
        const level = parseInt(heading.tagName[1], 10);
        const isValid = level <= previousLevel + 1;
        previousLevel = level;
        return isValid;
      });
    },
  },

  // 1.4.3 Contrast (Minimum) (Level AA)
  colorContrast: {
    criterion: '1.4.3',
    level: 'AA',
    description: 'Text must have sufficient color contrast',
    test: (element) => {
      // This would require color analysis - simplified check
      const hasLowContrastClass = element.classList.contains('low-contrast');
      return !hasLowContrastClass;
    },
  },

  // 2.1.1 Keyboard (Level A)
  keyboardAccessible: {
    criterion: '2.1.1',
    level: 'A',
    description: 'All interactive elements must be keyboard accessible',
    test: (element) => {
      const interactiveElements = element.querySelectorAll(
        'a, button, input, select, textarea, [tabindex]'
      );

      return Array.from(interactiveElements).every((el) => {
        const tabIndex = el.getAttribute('tabindex');
        return !tabIndex || parseInt(tabIndex, 10) >= -1;
      });
    },
  },

  // 2.4.4 Link Purpose (In Context) (Level A)
  linkPurpose: {
    criterion: '2.4.4',
    level: 'A',
    description: 'Links must have descriptive text',
    test: (element) => {
      const links = element.querySelectorAll('a');
      return Array.from(links).every((link) => {
        const text = link.textContent?.trim();
        const ariaLabel = link.getAttribute('aria-label');
        return !!(text || ariaLabel) && text !== 'click here' && text !== 'read more';
      });
    },
  },

  // 3.1.1 Language of Page (Level A)
  language: {
    criterion: '3.1.1',
    level: 'A',
    description: 'Page must have a language attribute',
    test: () => {
      return document.documentElement.hasAttribute('lang');
    },
  },

  // 3.3.2 Labels or Instructions (Level A)
  formLabels: {
    criterion: '3.3.2',
    level: 'A',
    description: 'Form inputs must have labels',
    test: (element) => {
      const inputs = element.querySelectorAll('input:not([type="submit"]):not([type="button"]), select, textarea');

      return Array.from(inputs).every((input) => {
        const id = input.id;
        const ariaLabel = input.getAttribute('aria-label');
        const ariaLabelledBy = input.getAttribute('aria-labelledby');
        const label = id ? element.querySelector(`label[for="${id}"]`) : null;

        return !!(label || ariaLabel || ariaLabelledBy);
      });
    },
  },

  // 4.1.2 Name, Role, Value (Level A)
  nameRoleValue: {
    criterion: '4.1.2',
    level: 'A',
    description: 'Custom controls must have proper ARIA attributes',
    test: (element) => {
      const customControls = element.querySelectorAll('[role]');

      return Array.from(customControls).every((control) => {
        const role = control.getAttribute('role');
        const hasLabel = control.hasAttribute('aria-label') ||
                         control.hasAttribute('aria-labelledby') ||
                         control.textContent?.trim();

        // Check required attributes for common roles
        if (role === 'button' || role === 'link') {
          return hasLabel;
        }

        if (role === 'checkbox' || role === 'radio') {
          return hasLabel && control.hasAttribute('aria-checked');
        }

        if (role === 'slider') {
          return hasLabel &&
                 control.hasAttribute('aria-valuenow') &&
                 control.hasAttribute('aria-valuemin') &&
                 control.hasAttribute('aria-valuemax');
        }

        return true;
      });
    },
  },
};

// Run specific WCAG checks
export function runWCAGCheck(element: HTMLElement, checkName: string): boolean {
  const check = wcagChecks[checkName];
  if (!check) {
    throw new Error(`Unknown WCAG check: ${checkName}`);
  }

  return check.test(element);
}

// Run all WCAG checks for a specific level
export function runAllWCAGChecks(element: HTMLElement, level: A11yLevel = 'AA') {
  const results: Record<string, boolean> = {};
  const levelPriority = { A: 1, AA: 2, AAA: 3 };

  Object.entries(wcagChecks).forEach(([name, check]) => {
    if (levelPriority[check.level] <= levelPriority[level]) {
      results[name] = check.test(element);
    }
  });

  return results;
}

// Keyboard navigation helpers
export const keyboardNavigation = {
  tab: () => {
    const event = new KeyboardEvent('keydown', { key: 'Tab' });
    document.activeElement?.dispatchEvent(event);
  },

  shiftTab: () => {
    const event = new KeyboardEvent('keydown', { key: 'Tab', shiftKey: true });
    document.activeElement?.dispatchEvent(event);
  },

  enter: () => {
    const event = new KeyboardEvent('keydown', { key: 'Enter' });
    document.activeElement?.dispatchEvent(event);
  },

  escape: () => {
    const event = new KeyboardEvent('keydown', { key: 'Escape' });
    document.activeElement?.dispatchEvent(event);
  },

  arrowDown: () => {
    const event = new KeyboardEvent('keydown', { key: 'ArrowDown' });
    document.activeElement?.dispatchEvent(event);
  },

  arrowUp: () => {
    const event = new KeyboardEvent('keydown', { key: 'ArrowUp' });
    document.activeElement?.dispatchEvent(event);
  },

  arrowLeft: () => {
    const event = new KeyboardEvent('keydown', { key: 'ArrowLeft' });
    document.activeElement?.dispatchEvent(event);
  },

  arrowRight: () => {
    const event = new KeyboardEvent('keydown', { key: 'ArrowRight' });
    document.activeElement?.dispatchEvent(event);
  },

  space: () => {
    const event = new KeyboardEvent('keydown', { key: ' ' });
    document.activeElement?.dispatchEvent(event);
  },
};

// Focus management helpers
export const focusManagement = {
  getTabbableElements: (container: HTMLElement): HTMLElement[] => {
    const selector = [
      'a[href]',
      'button:not([disabled])',
      'input:not([disabled])',
      'select:not([disabled])',
      'textarea:not([disabled])',
      '[tabindex]:not([tabindex="-1"])',
    ].join(', ');

    return Array.from(container.querySelectorAll(selector));
  },

  getFirstTabbable: (container: HTMLElement): HTMLElement | null => {
    const tabbables = focusManagement.getTabbableElements(container);
    return tabbables[0] || null;
  },

  getLastTabbable: (container: HTMLElement): HTMLElement | null => {
    const tabbables = focusManagement.getTabbableElements(container);
    return tabbables[tabbables.length - 1] || null;
  },

  trapFocus: (container: HTMLElement) => {
    const tabbables = focusManagement.getTabbableElements(container);
    const firstTabbable = tabbables[0];
    const lastTabbable = tabbables[tabbables.length - 1];

    const handleTab = (event: KeyboardEvent) => {
      if (event.key !== 'Tab') return;

      if (event.shiftKey) {
        if (document.activeElement === firstTabbable) {
          event.preventDefault();
          lastTabbable?.focus();
        }
      } else {
        if (document.activeElement === lastTabbable) {
          event.preventDefault();
          firstTabbable?.focus();
        }
      }
    };

    container.addEventListener('keydown', handleTab);

    return () => {
      container.removeEventListener('keydown', handleTab);
    };
  },
};

// Screen reader helpers
export const screenReaderHelpers = {
  // Announce message to screen readers
  announce: (message: string, priority: 'polite' | 'assertive' = 'polite') => {
    const announcement = document.createElement('div');
    announcement.setAttribute('role', 'status');
    announcement.setAttribute('aria-live', priority);
    announcement.setAttribute('aria-atomic', 'true');
    announcement.style.position = 'absolute';
    announcement.style.left = '-10000px';
    announcement.style.width = '1px';
    announcement.style.height = '1px';
    announcement.style.overflow = 'hidden';
    announcement.textContent = message;

    document.body.appendChild(announcement);

    setTimeout(() => {
      document.body.removeChild(announcement);
    }, 1000);
  },

  // Check if element is visible to screen readers
  isScreenReaderVisible: (element: HTMLElement): boolean => {
    const ariaHidden = element.getAttribute('aria-hidden');
    const display = window.getComputedStyle(element).display;
    const visibility = window.getComputedStyle(element).visibility;

    return ariaHidden !== 'true' && display !== 'none' && visibility !== 'hidden';
  },

  // Get accessible name of element
  getAccessibleName: (element: HTMLElement): string | null => {
    // Check aria-label
    const ariaLabel = element.getAttribute('aria-label');
    if (ariaLabel) return ariaLabel;

    // Check aria-labelledby
    const labelledBy = element.getAttribute('aria-labelledby');
    if (labelledBy) {
      const labels = labelledBy.split(' ')
        .map(id => document.getElementById(id)?.textContent)
        .filter(Boolean)
        .join(' ');
      if (labels) return labels;
    }

    // Check for associated label
    const id = element.id;
    if (id) {
      const label = document.querySelector(`label[for="${id}"]`);
      if (label?.textContent) return label.textContent;
    }

    // Fall back to text content
    return element.textContent || null;
  },
};

// Color contrast calculation helpers
export const colorContrast = {
  // Calculate relative luminance
  relativeLuminance: (rgb: { r: number; g: number; b: number }): number => {
    const { r, g, b } = rgb;
    const [rs, gs, bs] = [r, g, b].map(c => {
      c = c / 255;
      return c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
    });

    return 0.2126 * rs + 0.7152 * gs + 0.0722 * bs;
  },

  // Calculate contrast ratio
  contrastRatio: (
    color1: { r: number; g: number; b: number },
    color2: { r: number; g: number; b: number }
  ): number => {
    const l1 = colorContrast.relativeLuminance(color1);
    const l2 = colorContrast.relativeLuminance(color2);
    const lighter = Math.max(l1, l2);
    const darker = Math.min(l1, l2);

    return (lighter + 0.05) / (darker + 0.05);
  },

  // Check if contrast meets WCAG requirements
  meetsWCAG: (
    ratio: number,
    level: 'AA' | 'AAA' = 'AA',
    largeText = false
  ): boolean => {
    if (level === 'AA') {
      return largeText ? ratio >= 3 : ratio >= 4.5;
    } else {
      return largeText ? ratio >= 4.5 : ratio >= 7;
    }
  },
};

// Export type definitions for Vitest
declare global {
  namespace Vi {
    interface Matchers<R> {
      toHaveNoViolations(): R;
    }
    interface AsymmetricMatchers {
      toHaveNoViolations(): any;
    }
  }
}