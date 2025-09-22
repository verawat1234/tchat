import { screen, waitFor, within, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react';

/**
 * Radix UI Testing Utilities
 * Helper functions for testing Radix UI components with their async behaviors
 */

// Helper to wait for Radix UI portal content to be rendered
export async function waitForRadixPortal(testId?: string) {
  await waitFor(() => {
    const portals = document.querySelectorAll('[data-radix-portal]');
    expect(portals.length).toBeGreaterThan(0);

    if (testId) {
      const portal = screen.getByTestId(testId);
      expect(portal).toBeInTheDocument();
    }
  });
}

// Helper for testing Radix UI Select components
export const radixSelect = {
  async open(triggerElement: HTMLElement) {
    await userEvent.click(triggerElement);
    // Wait for the portal and content to render
    await waitFor(() => {
      const content = document.querySelector('[role="listbox"]');
      expect(content).toBeInTheDocument();
    });
  },

  async selectOption(optionText: string) {
    const option = await screen.findByRole('option', { name: optionText });
    await userEvent.click(option);
    // Wait for the select to close
    await waitFor(() => {
      const content = document.querySelector('[role="listbox"]');
      expect(content).not.toBeInTheDocument();
    });
  },

  async openAndSelect(triggerElement: HTMLElement, optionText: string) {
    await this.open(triggerElement);
    await this.selectOption(optionText);
  },
};

// Helper for testing Radix UI Dialog components
export const radixDialog = {
  async open(triggerElement: HTMLElement) {
    await userEvent.click(triggerElement);
    // Wait for dialog to be in the document
    await waitFor(() => {
      const dialog = screen.getByRole('dialog');
      expect(dialog).toBeInTheDocument();
    });
  },

  async close(method: 'escape' | 'overlay' | 'close-button' = 'escape') {
    switch (method) {
      case 'escape':
        await userEvent.keyboard('{Escape}');
        break;
      case 'overlay':
        const overlay = document.querySelector('[data-radix-dialog-overlay]');
        if (overlay) {
          await userEvent.click(overlay as HTMLElement);
        }
        break;
      case 'close-button':
        const closeButton = screen.getByRole('button', { name: /close/i });
        await userEvent.click(closeButton);
        break;
    }

    // Wait for dialog to be removed
    await waitFor(() => {
      const dialog = screen.queryByRole('dialog');
      expect(dialog).not.toBeInTheDocument();
    });
  },

  async expectOpen() {
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  },

  async expectClosed() {
    await waitFor(() => {
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });
  },
};

// Helper for testing Radix UI Dropdown Menu components
export const radixDropdown = {
  async open(triggerElement: HTMLElement) {
    await userEvent.click(triggerElement);
    await waitFor(() => {
      const menu = screen.getByRole('menu');
      expect(menu).toBeInTheDocument();
    });
  },

  async selectItem(itemText: string) {
    const item = await screen.findByRole('menuitem', { name: itemText });
    await userEvent.click(item);
  },

  async close() {
    await userEvent.keyboard('{Escape}');
    await waitFor(() => {
      const menu = screen.queryByRole('menu');
      expect(menu).not.toBeInTheDocument();
    });
  },
};

// Helper for testing Radix UI Accordion components
export const radixAccordion = {
  async expandItem(triggerElement: HTMLElement) {
    const isExpanded = triggerElement.getAttribute('aria-expanded') === 'true';

    if (!isExpanded) {
      await userEvent.click(triggerElement);
      await waitFor(() => {
        expect(triggerElement).toHaveAttribute('aria-expanded', 'true');
      });
    }
  },

  async collapseItem(triggerElement: HTMLElement) {
    const isExpanded = triggerElement.getAttribute('aria-expanded') === 'true';

    if (isExpanded) {
      await userEvent.click(triggerElement);
      await waitFor(() => {
        expect(triggerElement).toHaveAttribute('aria-expanded', 'false');
      });
    }
  },

  async expectContentVisible(content: string) {
    await waitFor(() => {
      expect(screen.getByText(content)).toBeVisible();
    });
  },
};

// Helper for testing Radix UI Tooltip components
export const radixTooltip = {
  async show(triggerElement: HTMLElement) {
    await userEvent.hover(triggerElement);
    // Tooltips usually have a delay
    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 100));
    });
  },

  async hide(triggerElement: HTMLElement) {
    await userEvent.unhover(triggerElement);
    await act(async () => {
      await new Promise(resolve => setTimeout(resolve, 100));
    });
  },

  async expectVisible(content: string) {
    await waitFor(() => {
      expect(screen.getByRole('tooltip', { name: content })).toBeInTheDocument();
    });
  },

  async expectHidden(content: string) {
    await waitFor(() => {
      expect(screen.queryByRole('tooltip', { name: content })).not.toBeInTheDocument();
    });
  },
};

// Helper for testing Radix UI Popover components
export const radixPopover = {
  async open(triggerElement: HTMLElement) {
    await userEvent.click(triggerElement);
    await waitFor(() => {
      const popover = document.querySelector('[data-radix-popover-content]');
      expect(popover).toBeInTheDocument();
    });
  },

  async close() {
    await userEvent.keyboard('{Escape}');
    await waitFor(() => {
      const popover = document.querySelector('[data-radix-popover-content]');
      expect(popover).not.toBeInTheDocument();
    });
  },
};

// Helper for testing Radix UI Avatar components with async image loading
export const radixAvatar = {
  async waitForImageLoad(container: HTMLElement) {
    const img = container.querySelector('img');
    if (img) {
      fireEvent.load(img);
      await waitFor(() => {
        // Check if image is displayed (Radix hides it until loaded)
        expect(img).toBeVisible();
      });
    }
  },

  async triggerImageError(container: HTMLElement) {
    const img = container.querySelector('img');
    if (img) {
      fireEvent.error(img);
      await waitFor(() => {
        // Fallback should be visible after error
        const fallback = container.querySelector('[data-slot="avatar-fallback"]');
        expect(fallback).toBeVisible();
      });
    }
  },

  expectFallback(container: HTMLElement, text: string) {
    const fallback = within(container).getByText(text);
    expect(fallback).toBeInTheDocument();
  },
};

// Helper for testing Radix UI Switch components
export const radixSwitch = {
  async toggle(switchElement: HTMLElement) {
    const isChecked = switchElement.getAttribute('aria-checked') === 'true';
    await userEvent.click(switchElement);

    await waitFor(() => {
      const newState = switchElement.getAttribute('aria-checked') === 'true';
      expect(newState).toBe(!isChecked);
    });
  },

  async check(switchElement: HTMLElement) {
    const isChecked = switchElement.getAttribute('aria-checked') === 'true';
    if (!isChecked) {
      await this.toggle(switchElement);
    }
  },

  async uncheck(switchElement: HTMLElement) {
    const isChecked = switchElement.getAttribute('aria-checked') === 'true';
    if (isChecked) {
      await this.toggle(switchElement);
    }
  },
};

// Helper for testing Radix UI Tabs components
export const radixTabs = {
  async selectTab(tabText: string) {
    const tab = screen.getByRole('tab', { name: tabText });
    await userEvent.click(tab);

    await waitFor(() => {
      expect(tab).toHaveAttribute('aria-selected', 'true');
    });
  },

  async expectPanelContent(content: string) {
    await waitFor(() => {
      const panel = screen.getByRole('tabpanel');
      expect(within(panel).getByText(content)).toBeInTheDocument();
    });
  },
};

// Helper for testing async loading states
export async function waitForLoadingToComplete(
  loadingTestId = 'loading',
  options = { timeout: 3000 }
) {
  // Wait for loading indicator to appear
  const loading = await screen.findByTestId(loadingTestId);
  expect(loading).toBeInTheDocument();

  // Wait for loading indicator to disappear
  await waitFor(
    () => {
      expect(screen.queryByTestId(loadingTestId)).not.toBeInTheDocument();
    },
    options
  );
}

// Helper for testing components with delays
export async function waitForDelay(ms: number) {
  await act(async () => {
    await new Promise(resolve => setTimeout(resolve, ms));
  });
}

// Helper for testing focus management in Radix UI components
export const radixFocus = {
  async expectFocusTrapped(container: HTMLElement) {
    const focusableElements = container.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );

    const firstElement = focusableElements[0] as HTMLElement;
    const lastElement = focusableElements[focusableElements.length - 1] as HTMLElement;

    // Tab from last element should go to first
    lastElement.focus();
    await userEvent.tab();
    expect(document.activeElement).toBe(firstElement);

    // Shift+Tab from first element should go to last
    firstElement.focus();
    await userEvent.tab({ shift: true });
    expect(document.activeElement).toBe(lastElement);
  },

  async expectFocusReturned(triggerElement: HTMLElement) {
    triggerElement.focus();
    expect(document.activeElement).toBe(triggerElement);
  },
};

// Helper to handle Radix UI's portal rendering
export function withinRadixPortal() {
  const portalRoot = document.querySelector('[data-radix-portal]');
  if (!portalRoot) {
    throw new Error('No Radix portal found in document');
  }
  return within(portalRoot as HTMLElement);
}

// Export all helpers as a single object for convenience
export const radixUI = {
  select: radixSelect,
  dialog: radixDialog,
  dropdown: radixDropdown,
  accordion: radixAccordion,
  tooltip: radixTooltip,
  popover: radixPopover,
  avatar: radixAvatar,
  switch: radixSwitch,
  tabs: radixTabs,
  focus: radixFocus,
  waitForPortal: waitForRadixPortal,
  waitForDelay,
  waitForLoadingToComplete,
  withinPortal: withinRadixPortal,
};