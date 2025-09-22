/**
 * Sidebar Component Contract Tests
 * CRITICAL: These tests MUST FAIL until Sidebar component is implemented
 */

import { render, screen, fireEvent } from '@testing-library/react';
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarFooter,
  SidebarNav,
  SidebarNavGroup,
  SidebarToggle
} from './sidebar';
import type {
  SidebarProps,
  SidebarContentProps,
  SidebarHeaderProps,
  SidebarFooterProps,
  SidebarNavProps,
  SidebarNavGroupProps,
  SidebarToggleProps,
  SidebarNavItem
} from '../../../../specs/001-agent-frontend-specialist/contracts/sidebar';

describe('Sidebar Contract Tests', () => {
  const baseProps: SidebarProps = {
    testId: 'sidebar-test'
  };

  describe('Basic Rendering', () => {
    test('renders sidebar with default props', () => {
      render(
        <Sidebar {...baseProps}>
          <SidebarContent>Sidebar content</SidebarContent>
        </Sidebar>
      );

      expect(screen.getByTestId('sidebar-test')).toBeInTheDocument();
      expect(screen.getByText('Sidebar content')).toBeInTheDocument();
    });

    test('applies custom className', () => {
      render(
        <Sidebar {...baseProps} className="custom-sidebar">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-test');
      expect(sidebar).toHaveClass('custom-sidebar');
    });
  });

  describe('Position and Layout', () => {
    test('applies position correctly', () => {
      const positions: SidebarProps['position'][] = ['left', 'right'];

      positions.forEach(position => {
        const { unmount } = render(
          <Sidebar {...baseProps} position={position} testId={`sidebar-${position}`}>
            <SidebarContent>Content</SidebarContent>
          </Sidebar>
        );

        const sidebar = screen.getByTestId(`sidebar-${position}`);
        expect(sidebar).toHaveClass(`sidebar-${position}`);
        unmount();
      });
    });

    test('applies custom width', () => {
      render(
        <Sidebar {...baseProps} width={300} testId="sidebar-width">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-width');
      expect(sidebar).toHaveStyle({ width: '300px' });
    });

    test('applies string width', () => {
      render(
        <Sidebar {...baseProps} width="25%" testId="sidebar-width-percent">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-width-percent');
      expect(sidebar).toHaveStyle({ width: '25%' });
    });

    test('applies variant styles correctly', () => {
      const variants: SidebarProps['variant'][] = ['default', 'floating', 'bordered'];

      variants.forEach(variant => {
        const { unmount } = render(
          <Sidebar {...baseProps} variant={variant} testId={`sidebar-${variant}`}>
            <SidebarContent>Content</SidebarContent>
          </Sidebar>
        );

        const sidebar = screen.getByTestId(`sidebar-${variant}`);
        expect(sidebar).toHaveClass(`sidebar-${variant}`);
        unmount();
      });
    });
  });

  describe('Collapsible Behavior', () => {
    test('supports collapsible functionality', () => {
      render(
        <Sidebar {...baseProps} collapsible testId="sidebar-collapsible">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-collapsible');
      expect(sidebar).toHaveClass('sidebar-collapsible');
    });

    test('applies collapsed state correctly', () => {
      render(
        <Sidebar {...baseProps} collapsed testId="sidebar-collapsed">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-collapsed');
      expect(sidebar).toHaveClass('sidebar-collapsed');
    });

    test('calls onCollapsedChange when collapsed state changes', () => {
      const onCollapsedChange = vi.fn();
      render(
        <Sidebar {...baseProps} collapsible onCollapsedChange={onCollapsedChange}>
          <SidebarContent>Content</SidebarContent>
          <SidebarToggle />
        </Sidebar>
      );

      const toggleButton = screen.getByRole('button');
      fireEvent.click(toggleButton);
      expect(onCollapsedChange).toHaveBeenCalledWith(true);
    });
  });

  describe('Responsive Features', () => {
    test('shows overlay when overlay is true', () => {
      render(
        <Sidebar {...baseProps} overlay testId="sidebar-overlay">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-overlay');
      expect(sidebar).toHaveClass('sidebar-overlay');
    });

    test('applies persistent styling when persistent is true', () => {
      render(
        <Sidebar {...baseProps} persistent testId="sidebar-persistent">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-persistent');
      expect(sidebar).toHaveClass('sidebar-persistent');
    });
  });

  describe('Resizable Functionality', () => {
    test('supports resizable functionality', () => {
      render(
        <Sidebar {...baseProps} resizable testId="sidebar-resizable">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-resizable');
      expect(sidebar).toHaveClass('sidebar-resizable');
      expect(sidebar.querySelector('[data-testid="resize-handle"]')).toBeInTheDocument();
    });

    test('applies min and max width constraints', () => {
      render(
        <Sidebar {...baseProps} resizable minWidth={200} maxWidth={400} testId="sidebar-constraints">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-constraints');
      expect(sidebar).toHaveStyle({
        minWidth: '200px',
        maxWidth: '400px'
      });
    });

    test('calls onResize when resizing', () => {
      const onResize = vi.fn();
      render(
        <Sidebar {...baseProps} resizable onResize={onResize} testId="sidebar-resize">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const resizeHandle = screen.getByTestId('resize-handle');
      fireEvent.mouseDown(resizeHandle);
      fireEvent.mouseMove(resizeHandle, { clientX: 350 });
      fireEvent.mouseUp(resizeHandle);

      expect(onResize).toHaveBeenCalled();
    });
  });

  describe('Accessibility', () => {
    test('supports ARIA label', () => {
      render(
        <Sidebar {...baseProps} aria-label="Main navigation sidebar">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-test');
      expect(sidebar).toHaveAttribute('aria-label', 'Main navigation sidebar');
    });

    test('supports custom role', () => {
      render(
        <Sidebar {...baseProps} role="navigation">
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-test');
      expect(sidebar).toHaveAttribute('role', 'navigation');
    });

    test('supports keyboard navigation', () => {
      render(
        <Sidebar {...baseProps} tabIndex={0}>
          <SidebarContent>Content</SidebarContent>
        </Sidebar>
      );

      const sidebar = screen.getByTestId('sidebar-test');
      expect(sidebar).toHaveAttribute('tabIndex', '0');
    });
  });
});

describe('SidebarContent Contract Tests', () => {
  test('renders content with default props', () => {
    render(
      <SidebarContent testId="sidebar-content-test">
        <p>Sidebar content</p>
      </SidebarContent>
    );

    expect(screen.getByTestId('sidebar-content-test')).toBeInTheDocument();
    expect(screen.getByText('Sidebar content')).toBeInTheDocument();
  });

  test('applies padding correctly', () => {
    const paddings: SidebarContentProps['padding'][] = ['none', 'sm', 'md', 'lg'];

    paddings.forEach(padding => {
      const { unmount } = render(
        <SidebarContent padding={padding} testId={`content-padding-${padding}`}>
          Content
        </SidebarContent>
      );

      const content = screen.getByTestId(`content-padding-${padding}`);
      expect(content).toHaveClass(`sidebar-content-padding-${padding}`);
      unmount();
    });
  });

  test('supports scrollable content', () => {
    render(
      <SidebarContent scrollable testId="sidebar-content-scroll">
        <div style={{ height: '2000px' }}>Very tall content</div>
      </SidebarContent>
    );

    const content = screen.getByTestId('sidebar-content-scroll');
    expect(content).toHaveClass('sidebar-content-scrollable');
  });
});

describe('SidebarHeader Contract Tests', () => {
  test('renders header with title', () => {
    render(
      <SidebarHeader title="Navigation" testId="sidebar-header-test" />
    );

    expect(screen.getByTestId('sidebar-header-test')).toBeInTheDocument();
    expect(screen.getByText('Navigation')).toBeInTheDocument();
  });

  test('renders actions when provided', () => {
    const actions = <button data-testid="header-action">Settings</button>;

    render(
      <SidebarHeader title="Navigation" actions={actions} testId="sidebar-header-actions" />
    );

    expect(screen.getByTestId('header-action')).toBeInTheDocument();
  });

  test('applies border when border is true', () => {
    render(
      <SidebarHeader title="Navigation" border testId="sidebar-header-border" />
    );

    const header = screen.getByTestId('sidebar-header-border');
    expect(header).toHaveClass('sidebar-header-border');
  });

  test('applies sticky styling when sticky is true', () => {
    render(
      <SidebarHeader title="Navigation" sticky testId="sidebar-header-sticky" />
    );

    const header = screen.getByTestId('sidebar-header-sticky');
    expect(header).toHaveClass('sidebar-header-sticky');
  });
});

describe('SidebarFooter Contract Tests', () => {
  test('renders footer with content', () => {
    render(
      <SidebarFooter testId="sidebar-footer-test">
        <p>Footer content</p>
      </SidebarFooter>
    );

    expect(screen.getByTestId('sidebar-footer-test')).toBeInTheDocument();
    expect(screen.getByText('Footer content')).toBeInTheDocument();
  });

  test('applies border when border is true', () => {
    render(
      <SidebarFooter border testId="sidebar-footer-border">
        Footer
      </SidebarFooter>
    );

    const footer = screen.getByTestId('sidebar-footer-border');
    expect(footer).toHaveClass('sidebar-footer-border');
  });

  test('applies sticky styling when sticky is true', () => {
    render(
      <SidebarFooter sticky testId="sidebar-footer-sticky">
        Footer
      </SidebarFooter>
    );

    const footer = screen.getByTestId('sidebar-footer-sticky');
    expect(footer).toHaveClass('sidebar-footer-sticky');
  });
});

describe('SidebarNav Contract Tests', () => {
  const navItems: SidebarNavItem[] = [
    {
      id: 'dashboard',
      label: 'Dashboard',
      href: '/dashboard',
      icon: <span data-testid="dashboard-icon">📊</span>
    },
    {
      id: 'users',
      label: 'Users',
      href: '/users',
      badge: '5'
    },
    {
      id: 'settings',
      label: 'Settings',
      children: [
        { id: 'general', label: 'General', href: '/settings/general' },
        { id: 'security', label: 'Security', href: '/settings/security' }
      ],
      expandable: true
    }
  ];

  test('renders navigation items', () => {
    render(
      <SidebarNav items={navItems} testId="sidebar-nav-test" />
    );

    expect(screen.getByTestId('sidebar-nav-test')).toBeInTheDocument();
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Users')).toBeInTheDocument();
    expect(screen.getByText('Settings')).toBeInTheDocument();
  });

  test('renders icons when provided', () => {
    render(
      <SidebarNav items={navItems} testId="sidebar-nav-icons" />
    );

    expect(screen.getByTestId('dashboard-icon')).toBeInTheDocument();
  });

  test('renders badges when provided', () => {
    render(
      <SidebarNav items={navItems} testId="sidebar-nav-badges" />
    );

    expect(screen.getByText('5')).toBeInTheDocument();
  });

  test('handles item selection', () => {
    const onValueChange = vi.fn();
    render(
      <SidebarNav
        items={navItems}
        value="dashboard"
        onValueChange={onValueChange}
        testId="sidebar-nav-selection"
      />
    );

    const usersItem = screen.getByText('Users');
    fireEvent.click(usersItem);
    expect(onValueChange).toHaveBeenCalledWith('users');
  });

  test('supports expandable items', () => {
    render(
      <SidebarNav items={navItems} testId="sidebar-nav-expandable" />
    );

    const settingsItem = screen.getByText('Settings');
    fireEvent.click(settingsItem);

    expect(screen.getByText('General')).toBeInTheDocument();
    expect(screen.getByText('Security')).toBeInTheDocument();
  });

  test('applies variant styles correctly', () => {
    const variants: SidebarNavProps['variant'][] = ['default', 'pills', 'tree'];

    variants.forEach(variant => {
      const { unmount } = render(
        <SidebarNav items={navItems} variant={variant} testId={`nav-${variant}`} />
      );

      const nav = screen.getByTestId(`nav-${variant}`);
      expect(nav).toHaveClass(`sidebar-nav-${variant}`);
      unmount();
    });
  });

  test('handles multiple expanded sections', () => {
    const onExpandedChange = vi.fn();
    render(
      <SidebarNav
        items={navItems}
        multiple
        onExpandedChange={onExpandedChange}
        testId="sidebar-nav-multiple"
      />
    );

    const nav = screen.getByTestId('sidebar-nav-multiple');
    expect(nav).toHaveClass('sidebar-nav-multiple');
  });

  test('handles disabled items', () => {
    const disabledItems: SidebarNavItem[] = [
      { id: 'active', label: 'Active Item' },
      { id: 'disabled', label: 'Disabled Item', disabled: true }
    ];

    render(
      <SidebarNav items={disabledItems} testId="sidebar-nav-disabled" />
    );

    const disabledItem = screen.getByText('Disabled Item');
    expect(disabledItem).toHaveAttribute('aria-disabled', 'true');
  });
});

describe('SidebarNavGroup Contract Tests', () => {
  test('renders group with title', () => {
    render(
      <SidebarNavGroup title="Main Navigation" testId="sidebar-group-test">
        <div>Group content</div>
      </SidebarNavGroup>
    );

    expect(screen.getByTestId('sidebar-group-test')).toBeInTheDocument();
    expect(screen.getByText('Main Navigation')).toBeInTheDocument();
    expect(screen.getByText('Group content')).toBeInTheDocument();
  });

  test('supports collapsible functionality', () => {
    render(
      <SidebarNavGroup title="Collapsible Group" collapsible testId="sidebar-group-collapsible">
        <div>Group content</div>
      </SidebarNavGroup>
    );

    const group = screen.getByTestId('sidebar-group-collapsible');
    expect(group).toHaveClass('sidebar-group-collapsible');
  });

  test('applies collapsed state correctly', () => {
    render(
      <SidebarNavGroup title="Group" collapsed testId="sidebar-group-collapsed">
        <div>Group content</div>
      </SidebarNavGroup>
    );

    const group = screen.getByTestId('sidebar-group-collapsed');
    expect(group).toHaveClass('sidebar-group-collapsed');
  });

  test('calls onCollapsedChange when toggled', () => {
    const onCollapsedChange = vi.fn();
    render(
      <SidebarNavGroup
        title="Group"
        collapsible
        onCollapsedChange={onCollapsedChange}
        testId="sidebar-group-toggle"
      >
        <div>Content</div>
      </SidebarNavGroup>
    );

    const toggleButton = screen.getByRole('button');
    fireEvent.click(toggleButton);
    expect(onCollapsedChange).toHaveBeenCalledWith(true);
  });
});

describe('SidebarToggle Contract Tests', () => {
  test('renders toggle button', () => {
    render(<SidebarToggle testId="sidebar-toggle-test" />);

    expect(screen.getByTestId('sidebar-toggle-test')).toBeInTheDocument();
    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  test('applies collapsed state correctly', () => {
    render(<SidebarToggle collapsed testId="sidebar-toggle-collapsed" />);

    const toggle = screen.getByTestId('sidebar-toggle-collapsed');
    expect(toggle).toHaveClass('sidebar-toggle-collapsed');
  });

  test('calls onToggle when clicked', () => {
    const onToggle = vi.fn();
    render(<SidebarToggle onToggle={onToggle} testId="sidebar-toggle-click" />);

    const button = screen.getByRole('button');
    fireEvent.click(button);
    expect(onToggle).toHaveBeenCalledTimes(1);
  });

  test('applies position correctly', () => {
    const positions: SidebarToggleProps['position'][] = ['inside', 'outside'];

    positions.forEach(position => {
      const { unmount } = render(
        <SidebarToggle position={position} testId={`toggle-${position}`} />
      );

      const toggle = screen.getByTestId(`toggle-${position}`);
      expect(toggle).toHaveClass(`sidebar-toggle-${position}`);
      unmount();
    });
  });

  test('applies direction correctly', () => {
    const directions: SidebarToggleProps['direction'][] = ['left', 'right'];

    directions.forEach(direction => {
      const { unmount } = render(
        <SidebarToggle direction={direction} testId={`toggle-dir-${direction}`} />
      );

      const toggle = screen.getByTestId(`toggle-dir-${direction}`);
      expect(toggle).toHaveClass(`sidebar-toggle-${direction}`);
      unmount();
    });
  });
});