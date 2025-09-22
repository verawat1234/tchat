/**
 * Tooltip Test Utilities
 * Simplified testing setup for Radix UI Tooltip components
 */

import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import { TooltipProvider } from './tooltip';

/**
 * Custom render for Tooltip components with provider
 */
export function renderTooltip(ui: React.ReactElement, options = {}) {
  return render(
    <TooltipProvider delayDuration={0} skipDelayDuration={0}>
      {ui}
    </TooltipProvider>,
    options
  );
}

/**
 * Wait for tooltip content to appear in the DOM
 * Radix portals take time to render
 */
export async function waitForTooltipContent(text: string | RegExp) {
  // First wait for any portal to appear
  await waitFor(() => {
    const portals = document.querySelectorAll('[data-radix-portal]');
    if (portals.length === 0) {
      // If no portal, check if content is in regular DOM
      const content = screen.queryByText(text);
      if (!content) {
        throw new Error('No tooltip content found');
      }
    }
  }, { timeout: 100 });

  // Then wait for the specific content
  return waitFor(() => {
    const content = screen.getByText(text);
    expect(content).toBeInTheDocument();
    return content;
  }, { timeout: 500 });
}

/**
 * Helper to check if tooltip is NOT in the DOM
 */
export async function expectTooltipNotVisible(text: string | RegExp) {
  await waitFor(() => {
    expect(screen.queryByText(text)).not.toBeInTheDocument();
  }, { timeout: 500 });
}

/**
 * Mock Radix Tooltip for simpler testing
 */
export const MockTooltip: React.FC<{
  children: React.ReactNode;
  content: string;
  open?: boolean;
  defaultOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
}> = ({ children, content, open, defaultOpen = false, onOpenChange }) => {
  const [isOpen, setIsOpen] = React.useState(defaultOpen);
  const controlledOpen = open !== undefined ? open : isOpen;

  const handleToggle = () => {
    const newState = !controlledOpen;
    setIsOpen(newState);
    onOpenChange?.(newState);
  };

  return (
    <>
      <div onClick={handleToggle} onMouseEnter={() => handleToggle()} onFocus={() => handleToggle()}>
        {children}
      </div>
      {controlledOpen && (
        <div role="tooltip" data-state="open">
          {content}
        </div>
      )}
    </>
  );
};