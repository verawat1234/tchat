# Quickstart: Common UI Components

**Feature**: Common UI Component Breakdown
**Date**: 2025-09-21
**Status**: Implementation Ready

## Overview

This guide provides a quick start for implementing and using the common UI component library in the Tchat application. The component library extracts reusable UI patterns from the existing codebase and provides a consistent, maintainable foundation for future development.

## Prerequisites

- Node.js 18+ with npm
- TypeScript 5.3+
- React 18.3+
- Vite 6.3+ (already configured)
- TailwindCSS v4 (already configured)

## Quick Start

### 1. Component Structure

Create the component library structure:

```bash
# Create component directories
mkdir -p apps/web/src/components/common
mkdir -p apps/web/src/components/common/{pagination,tabs,card,chat-message,badge,layout,header,sidebar}
mkdir -p apps/web/tests/components/common
```

### 2. Install Additional Dependencies

```bash
cd apps/web
npm install @testing-library/react @testing-library/jest-dom @testing-library/user-event vitest jsdom
npm install --save-dev storybook @storybook/react @storybook/react-vite
```

### 3. Basic Component Implementation

Start with a simple component (Badge):

```typescript
// apps/web/src/components/common/badge/Badge.tsx
import React from 'react';
import { BadgeProps } from '../../../specs/001-agent-frontend-specialist/contracts';
import { cn } from '../../utils';

export const Badge: React.FC<BadgeProps> = ({
  variant = 'default',
  size = 'md',
  children,
  className,
  testId,
  ...props
}) => {
  return (
    <span
      data-testid={testId}
      className={cn(
        'inline-flex items-center rounded-full font-medium',
        {
          'bg-gray-100 text-gray-800': variant === 'default',
          'bg-green-100 text-green-800': variant === 'success',
          'bg-yellow-100 text-yellow-800': variant === 'warning',
          'bg-red-100 text-red-800': variant === 'danger',
          'bg-blue-100 text-blue-800': variant === 'info',
        },
        {
          'px-2 py-1 text-xs': size === 'sm',
          'px-2.5 py-1.5 text-sm': size === 'md',
          'px-3 py-2 text-base': size === 'lg',
        },
        className
      )}
      {...props}
    >
      {children}
    </span>
  );
};
```

### 4. Component Testing

Create a test file:

```typescript
// apps/web/tests/components/common/Badge.test.tsx
import { render, screen } from '@testing-library/react';
import { Badge } from '../../../src/components/common/badge/Badge';

describe('Badge', () => {
  it('renders with default props', () => {
    render(<Badge>Test Badge</Badge>);
    expect(screen.getByText('Test Badge')).toBeInTheDocument();
  });

  it('applies correct variant styles', () => {
    render(<Badge variant="success" testId="badge">Success</Badge>);
    const badge = screen.getByTestId('badge');
    expect(badge).toHaveClass('bg-green-100', 'text-green-800');
  });
});
```

### 5. Storybook Story

Create a Storybook story:

```typescript
// apps/web/src/components/common/badge/Badge.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import { Badge } from './Badge';

const meta: Meta<typeof Badge> = {
  title: 'Common/Badge',
  component: Badge,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'success', 'warning', 'danger', 'info'],
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
    },
  },
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: 'Badge',
  },
};

export const Success: Story = {
  args: {
    children: 'Success',
    variant: 'success',
  },
};

export const WithCount: Story = {
  args: {
    children: '99+',
    variant: 'danger',
  },
};
```

### 6. Component Export

Create index files for easy importing:

```typescript
// apps/web/src/components/common/badge/index.ts
export { Badge } from './Badge';
export type { BadgeProps } from '../../../specs/001-agent-frontend-specialist/contracts';

// apps/web/src/components/common/index.ts
export * from './badge';
// Add other components as they're implemented
```

## Usage Examples

### Basic Usage

```tsx
import { Badge } from '@/components/common';

function App() {
  return (
    <div>
      <Badge variant="success">Active</Badge>
      <Badge variant="warning" size="sm">Pending</Badge>
      <Badge variant="danger">Error</Badge>
    </div>
  );
}
```

### With Existing Components

```tsx
import { Badge } from '@/components/common';
import { Button } from '@/components/ui/button';

function NotificationButton() {
  return (
    <div className="relative">
      <Button>Notifications</Button>
      <Badge
        variant="danger"
        size="sm"
        className="absolute -top-2 -right-2"
      >
        5
      </Badge>
    </div>
  );
}
```

## Development Workflow

### 1. Component Development Order

Implement components in dependency order:

1. **Foundation**: Badge, Layout (Grid, Flex, Stack)
2. **Basic Containers**: Card
3. **Navigation**: Tabs, Pagination
4. **Complex Components**: Header, Sidebar
5. **Domain-Specific**: ChatMessage

### 2. Testing Strategy

For each component:

1. **Unit Tests**: Props, rendering, interactions
2. **Visual Tests**: Storybook stories with controls
3. **Integration Tests**: Component composition
4. **Accessibility Tests**: Keyboard navigation, screen readers

### 3. Code Quality

```bash
# Type checking
npm run type-check

# Linting
npm run lint

# Testing
npm run test

# Visual testing
npm run storybook
```

## Configuration Files

### Vitest Config

```typescript
// apps/web/vitest.config.ts
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react-swc';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test-setup.ts'],
  },
});
```

### Test Setup

```typescript
// apps/web/src/test-setup.ts
import '@testing-library/jest-dom';
```

### Storybook Config

```typescript
// apps/web/.storybook/main.ts
import type { StorybookConfig } from '@storybook/react-vite';

const config: StorybookConfig = {
  stories: ['../src/**/*.stories.@(js|jsx|ts|tsx|mdx)'],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-interactions',
  ],
  framework: {
    name: '@storybook/react-vite',
    options: {},
  },
};

export default config;
```

## Performance Considerations

### Bundle Size Monitoring

```bash
# Analyze bundle size
npm run build
npm run analyze
```

### Code Splitting

```typescript
// Lazy load complex components
const ChatMessage = lazy(() => import('./chat-message/ChatMessage'));
```

### Optimization Checklist

- [ ] Tree shaking enabled
- [ ] Dynamic imports for heavy components
- [ ] Memoization for expensive calculations
- [ ] Bundle size under 50KB per component
- [ ] No unnecessary re-renders

## Migration Path

### Phase 1: Core Components (Week 1-2)
- Badge ✓
- Layout components
- Card component

### Phase 2: Navigation (Week 3-4)
- Tabs component
- Pagination component

### Phase 3: Complex Components (Week 5-6)
- Header component
- Sidebar component

### Phase 4: Domain-Specific (Week 7-8)
- ChatMessage component
- Integration and optimization

## Support

### Documentation
- Component props: See TypeScript interfaces in `/contracts/`
- Storybook: `npm run storybook`
- Examples: Check `/stories/` directory

### Testing
- Unit tests: `npm run test`
- E2E tests: `npm run test:e2e`
- Visual regression: Storybook visual tests

### Getting Help

1. Check component contracts in `/contracts/`
2. Review Storybook documentation
3. Run existing tests for examples
4. Check implementation in `/common/` directory