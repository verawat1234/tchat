/**
 * Badge Component Contract Tests
 * CRITICAL: These tests MUST FAIL until Badge component is implemented
 */

import { render, screen, fireEvent } from '@testing-library/react';
import { Badge, BadgeGroup } from './badge';
import type { BadgeProps, BadgeGroupProps } from '../../../../specs/001-agent-frontend-specialist/contracts/badge';

describe('Badge Contract Tests', () => {
  describe('Basic Rendering', () => {
    test('renders badge with default props', () => {
      render(<Badge testId="badge-test">Default Badge</Badge>);
      expect(screen.getByTestId('badge-test')).toBeInTheDocument();
      expect(screen.getByText('Default Badge')).toBeInTheDocument();
    });

    test('applies custom className', () => {
      render(<Badge className="custom-class" testId="badge-test">Badge</Badge>);
      const badge = screen.getByTestId('badge-test');
      expect(badge).toHaveClass('custom-class');
    });
  });

  describe('Variant Handling', () => {
    const variants: BadgeProps['variant'][] = ['default', 'success', 'warning', 'danger', 'info', 'secondary'];

    variants.forEach(variant => {
      test(`renders ${variant} variant correctly`, () => {
        render(<Badge variant={variant} testId={`badge-${variant}`}>Badge</Badge>);
        const badge = screen.getByTestId(`badge-${variant}`);
        expect(badge).toBeInTheDocument();
        expect(badge).toHaveClass(`badge-${variant}`);
      });
    });
  });

  describe('Size Variations', () => {
    const sizes: BadgeProps['size'][] = ['sm', 'md', 'lg'];

    sizes.forEach(size => {
      test(`renders ${size} size correctly`, () => {
        render(<Badge size={size} testId={`badge-${size}`}>Badge</Badge>);
        const badge = screen.getByTestId(`badge-${size}`);
        expect(badge).toHaveClass(`badge-${size}`);
      });
    });
  });

  describe('Count and Content', () => {
    test('displays count when provided', () => {
      render(<Badge count={5} testId="badge-count" />);
      expect(screen.getByText('5')).toBeInTheDocument();
    });

    test('shows "99+" when count exceeds max', () => {
      render(<Badge count={150} max={99} testId="badge-max" />);
      expect(screen.getByText('99+')).toBeInTheDocument();
    });

    test('displays custom content', () => {
      render(<Badge content="NEW" testId="badge-content" />);
      expect(screen.getByText('NEW')).toBeInTheDocument();
    });

    test('shows count 0 when showZero is true', () => {
      render(<Badge count={0} showZero testId="badge-zero" />);
      expect(screen.getByText('0')).toBeInTheDocument();
    });

    test('hides count 0 when showZero is false', () => {
      const { container } = render(<Badge count={0} showZero={false} testId="badge-hidden" />);
      expect(container.firstChild).toBeNull();
    });
  });

  describe('Dot Mode', () => {
    test('renders as dot when dot prop is true', () => {
      render(<Badge dot testId="badge-dot" />);
      const badge = screen.getByTestId('badge-dot');
      expect(badge).toHaveClass('badge-dot');
    });
  });

  describe('Animation and Styling', () => {
    test('applies pulse animation when pulse is true', () => {
      render(<Badge pulse testId="badge-pulse">Badge</Badge>);
      const badge = screen.getByTestId('badge-pulse');
      expect(badge).toHaveClass('badge-pulse');
    });

    test('applies outline styling when outline is true', () => {
      render(<Badge outline testId="badge-outline">Badge</Badge>);
      const badge = screen.getByTestId('badge-outline');
      expect(badge).toHaveClass('badge-outline');
    });
  });

  describe('Removable Badge', () => {
    test('shows remove button when removable is true', () => {
      const onRemove = vi.fn();
      render(<Badge removable onRemove={onRemove} testId="badge-removable">Badge</Badge>);

      const removeButton = screen.getByRole('button');
      expect(removeButton).toBeInTheDocument();
    });

    test('calls onRemove when remove button is clicked', () => {
      const onRemove = vi.fn();
      render(<Badge removable onRemove={onRemove} testId="badge-removable">Badge</Badge>);

      const removeButton = screen.getByRole('button');
      fireEvent.click(removeButton);
      expect(onRemove).toHaveBeenCalledTimes(1);
    });
  });

  describe('Icon Support', () => {
    test('renders icon when provided', () => {
      const icon = <span data-testid="badge-icon">★</span>;
      render(<Badge icon={icon} testId="badge-with-icon">Badge</Badge>);

      expect(screen.getByTestId('badge-icon')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    test('supports ARIA label', () => {
      render(<Badge aria-label="Notification badge" testId="badge-aria">Badge</Badge>);
      const badge = screen.getByTestId('badge-aria');
      expect(badge).toHaveAttribute('aria-label', 'Notification badge');
    });

    test('supports role attribute', () => {
      render(<Badge role="status" testId="badge-role">Badge</Badge>);
      const badge = screen.getByTestId('badge-role');
      expect(badge).toHaveAttribute('role', 'status');
    });

    test('supports tabIndex for keyboard navigation', () => {
      render(<Badge tabIndex={0} testId="badge-tab">Badge</Badge>);
      const badge = screen.getByTestId('badge-tab');
      expect(badge).toHaveAttribute('tabIndex', '0');
    });
  });
});

describe('BadgeGroup Contract Tests', () => {
  test('renders multiple badges in group', () => {
    render(
      <BadgeGroup testId="badge-group">
        <Badge testId="badge-1">Badge 1</Badge>
        <Badge testId="badge-2">Badge 2</Badge>
        <Badge testId="badge-3">Badge 3</Badge>
      </BadgeGroup>
    );

    expect(screen.getByTestId('badge-group')).toBeInTheDocument();
    expect(screen.getByTestId('badge-1')).toBeInTheDocument();
    expect(screen.getByTestId('badge-2')).toBeInTheDocument();
    expect(screen.getByTestId('badge-3')).toBeInTheDocument();
  });

  test('applies spacing classes', () => {
    const spacings: BadgeGroupProps['spacing'][] = ['tight', 'normal', 'loose'];

    spacings.forEach(spacing => {
      const { unmount } = render(
        <BadgeGroup spacing={spacing} testId={`group-${spacing}`}>
          <Badge>Badge</Badge>
        </BadgeGroup>
      );

      const group = screen.getByTestId(`group-${spacing}`);
      expect(group).toHaveClass(`badge-group-${spacing}`);
      unmount();
    });
  });

  test('limits visible badges when max is set', () => {
    render(
      <BadgeGroup max={2} testId="badge-group-max">
        <Badge testId="badge-1">Badge 1</Badge>
        <Badge testId="badge-2">Badge 2</Badge>
        <Badge testId="badge-3">Badge 3</Badge>
      </BadgeGroup>
    );

    expect(screen.getByTestId('badge-1')).toBeInTheDocument();
    expect(screen.getByTestId('badge-2')).toBeInTheDocument();
    expect(screen.queryByTestId('badge-3')).not.toBeInTheDocument();
  });

  test('applies wrap classes when wrap is enabled', () => {
    render(
      <BadgeGroup wrap testId="badge-group-wrap">
        <Badge>Badge</Badge>
      </BadgeGroup>
    );

    const group = screen.getByTestId('badge-group-wrap');
    expect(group).toHaveClass('badge-group-wrap');
  });
});