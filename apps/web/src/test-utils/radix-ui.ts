/**
 * Radix UI Testing Utilities
 * Helper functions for testing Radix UI components
 */

import { waitFor, within } from '@testing-library/react';

/**
 * Wait for Radix UI portal content to render
 * Radix components often render in portals which appear asynchronously
 */
export async function waitForPortal(timeout = 1000) {
  await waitFor(() => {
    const portals = document.querySelectorAll('[data-radix-portal]');
    if (portals.length === 0) {
      throw new Error('Portal not yet rendered');
    }
  }, { timeout });
}

/**
 * Get the actual input element for Radix slider thumbs
 * Radix Slider renders the actual slider role on nested elements
 */
export function getSliderThumb(container: HTMLElement, index = 0) {
  const thumbs = container.querySelectorAll('[role="slider"]');
  return thumbs[index] as HTMLElement;
}

/**
 * Get the progress bar element with ARIA attributes
 * Radix Progress sets ARIA attributes on the root element
 */
export function getProgressBar(container: HTMLElement) {
  const progress = container.querySelector('[role="progressbar"]') ||
                   container.querySelector('[data-slot="progress"]');
  return progress as HTMLElement;
}

/**
 * Wait for tooltip to appear
 * Tooltips have delays and animations
 */
export async function waitForTooltip(text: string | RegExp, timeout = 2000) {
  return waitFor(() => {
    const tooltips = Array.from(document.querySelectorAll('[role="tooltip"]'));
    const found = tooltips.find(el => {
      const content = el.textContent || '';
      return typeof text === 'string' ? content.includes(text) : text.test(content);
    });
    if (!found) {
      throw new Error(`Tooltip with text "${text}" not found`);
    }
    return found;
  }, { timeout });
}

/**
 * Get dialog content including portal-rendered content
 */
export function getDialogContent() {
  return document.querySelector('[role="dialog"]') as HTMLElement;
}

/**
 * Check if element has focus within (for composite widgets)
 */
export function hasFocusWithin(element: HTMLElement): boolean {
  return element.contains(document.activeElement);
}

/**
 * Wait for animation to complete
 * Useful for Radix components with enter/exit animations
 */
export async function waitForAnimation(element: HTMLElement, timeout = 500) {
  return new Promise(resolve => {
    const handleAnimationEnd = () => {
      element.removeEventListener('animationend', handleAnimationEnd);
      element.removeEventListener('transitionend', handleAnimationEnd);
      resolve(element);
    };

    element.addEventListener('animationend', handleAnimationEnd);
    element.addEventListener('transitionend', handleAnimationEnd);

    // Fallback timeout
    setTimeout(() => resolve(element), timeout);
  });
}

/**
 * Get all focusable elements within a container
 * Useful for testing focus management
 */
export function getFocusableElements(container: HTMLElement): HTMLElement[] {
  const selector = 'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])';
  return Array.from(container.querySelectorAll(selector)) as HTMLElement[];
}

/**
 * Trigger hover state properly for Radix components
 */
export async function triggerHover(element: HTMLElement) {
  const mouseEnter = new MouseEvent('mouseenter', {
    bubbles: true,
    cancelable: true,
    view: window,
  });

  const mouseOver = new MouseEvent('mouseover', {
    bubbles: true,
    cancelable: true,
    view: window,
  });

  element.dispatchEvent(mouseEnter);
  element.dispatchEvent(mouseOver);

  // Wait a tick for state updates
  await new Promise(resolve => setTimeout(resolve, 0));
}

/**
 * Trigger unhover state properly for Radix components
 */
export async function triggerUnhover(element: HTMLElement) {
  const mouseLeave = new MouseEvent('mouseleave', {
    bubbles: true,
    cancelable: true,
    view: window,
  });

  const mouseOut = new MouseEvent('mouseout', {
    bubbles: true,
    cancelable: true,
    view: window,
  });

  element.dispatchEvent(mouseLeave);
  element.dispatchEvent(mouseOut);

  // Wait a tick for state updates
  await new Promise(resolve => setTimeout(resolve, 0));
}

/**
 * Helper to test accessibility attributes
 */
export function expectAccessibleName(element: HTMLElement, name: string) {
  const ariaLabel = element.getAttribute('aria-label');
  const ariaLabelledBy = element.getAttribute('aria-labelledby');

  if (ariaLabel) {
    expect(ariaLabel).toBe(name);
  } else if (ariaLabelledBy) {
    const labelElement = document.getElementById(ariaLabelledBy);
    expect(labelElement?.textContent).toBe(name);
  } else {
    // Check for accessible name computed from content
    expect(element.textContent).toContain(name);
  }
}

/**
 * Create a test wrapper with Radix provider setup
 */
export function createRadixWrapper(children: React.ReactNode) {
  // This can be extended with any required Radix providers
  return <>{children}</>;
}