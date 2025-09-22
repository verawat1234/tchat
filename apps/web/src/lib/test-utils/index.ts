import React, { ReactElement } from 'react';
import { render, RenderOptions, RenderResult, queries } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';

// Custom queries for enhanced element selection
const customQueries = {
  getByTestId: (container: HTMLElement, id: string) => {
    return container.querySelector(`[data-testid="${id}"]`);
  },
  getAllByTestId: (container: HTMLElement, id: string) => {
    return Array.from(container.querySelectorAll(`[data-testid="${id}"]`));
  },
  queryByTestId: (container: HTMLElement, id: string) => {
    return container.querySelector(`[data-testid="${id}"]`) || null;
  },
  findByTestId: async (container: HTMLElement, id: string) => {
    const element = container.querySelector(`[data-testid="${id}"]`);
    if (!element) {
      throw new Error(`Element with test ID "${id}" not found`);
    }
    return element;
  },
};

// Create a custom query client for testing
const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
        staleTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  });

// Provider wrapper for tests
interface TestProviderProps {
  children: React.ReactNode;
  queryClient?: QueryClient;
}

export function TestProvider({ children, queryClient }: TestProviderProps) {
  const testQueryClient = queryClient || createTestQueryClient();

  return (
    <QueryClientProvider client={testQueryClient}>
      {children}
    </QueryClientProvider>
  );
}

// Extended render function with providers
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  queryClient?: QueryClient;
}

export function customRender(
  ui: ReactElement,
  options?: CustomRenderOptions
): RenderResult & { user: ReturnType<typeof userEvent.setup> } {
  const { queryClient, ...renderOptions } = options || {};

  const Wrapper = ({ children }: { children: React.ReactNode }) => (
    <TestProvider queryClient={queryClient}>{children}</TestProvider>
  );

  const renderResult = render(ui, {
    wrapper: Wrapper,
    queries: {
      ...queries,
      ...customQueries,
    },
    ...renderOptions,
  });

  return {
    ...renderResult,
    user: userEvent.setup(),
  };
}

// Re-export everything from React Testing Library
export * from '@testing-library/react';
export { customRender as render, createTestQueryClient };

// Utility functions for common testing patterns
export const waitForLoadingToFinish = () => {
  return new Promise((resolve) => setTimeout(resolve, 0));
};

export const mockConsole = () => {
  const originalConsole = {
    error: console.error,
    warn: console.warn,
    log: console.log,
  };

  beforeEach(() => {
    console.error = vi.fn();
    console.warn = vi.fn();
    console.log = vi.fn();
  });

  afterEach(() => {
    console.error = originalConsole.error;
    console.warn = originalConsole.warn;
    console.log = originalConsole.log;
  });

  return {
    expectNoErrors: () => {
      expect(console.error).not.toHaveBeenCalled();
    },
    expectNoWarnings: () => {
      expect(console.warn).not.toHaveBeenCalled();
    },
  };
};

// Accessibility testing helpers
export const checkAccessibility = async (container: HTMLElement) => {
  const requiredAttributes = {
    buttons: ['aria-label', 'aria-describedby'],
    inputs: ['aria-label', 'aria-describedby', 'aria-invalid'],
    images: ['alt'],
    links: ['aria-label'],
  };

  const violations: string[] = [];

  // Check buttons
  container.querySelectorAll('button').forEach((button) => {
    if (!button.textContent && !button.getAttribute('aria-label')) {
      violations.push(`Button missing accessible label: ${button.outerHTML}`);
    }
  });

  // Check inputs
  container.querySelectorAll('input').forEach((input) => {
    const id = input.id;
    const label = id ? container.querySelector(`label[for="${id}"]`) : null;
    if (!label && !input.getAttribute('aria-label')) {
      violations.push(`Input missing accessible label: ${input.outerHTML}`);
    }
  });

  // Check images
  container.querySelectorAll('img').forEach((img) => {
    if (!img.getAttribute('alt')) {
      violations.push(`Image missing alt text: ${img.outerHTML}`);
    }
  });

  return {
    violations,
    passes: violations.length === 0,
  };
};

// Component state testing helpers
export const expectToBeDisabled = (element: HTMLElement) => {
  expect(element).toHaveAttribute('disabled');
  expect(element).toHaveAttribute('aria-disabled', 'true');
};

export const expectToBeEnabled = (element: HTMLElement) => {
  expect(element).not.toHaveAttribute('disabled');
  expect(element).not.toHaveAttribute('aria-disabled', 'true');
};

export const expectToHaveFocus = (element: HTMLElement) => {
  expect(document.activeElement).toBe(element);
};

// Form testing helpers
export const fillForm = async (
  user: ReturnType<typeof userEvent.setup>,
  formData: Record<string, string>
) => {
  for (const [fieldName, value] of Object.entries(formData)) {
    const field = document.querySelector(`[name="${fieldName}"]`) as HTMLElement;
    if (field) {
      await user.clear(field);
      await user.type(field, value);
    }
  }
};

// Component visibility helpers
export const expectToBeVisible = (element: HTMLElement) => {
  expect(element).toBeVisible();
  expect(element).not.toHaveStyle({ display: 'none' });
  expect(element).not.toHaveClass('hidden');
};

export const expectToBeHidden = (element: HTMLElement) => {
  expect(element).not.toBeVisible();
};