/**
 * Pagination Component Contract Tests
 * CRITICAL: These tests MUST FAIL until Pagination component is implemented
 */

import { render, screen, fireEvent } from '@testing-library/react';
import { Pagination, PaginationButton, PageSizeSelector } from './pagination';
import type { PaginationProps, PaginationButtonProps, PageSizeSelectorProps } from '../../../../specs/001-agent-frontend-specialist/contracts/pagination';

describe('Pagination Contract Tests', () => {
  const defaultProps: PaginationProps = {
    currentPage: 1,
    totalPages: 10,
    onPageChange: vi.fn(),
    testId: 'pagination-test'
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Basic Rendering', () => {
    test('renders pagination with required props', () => {
      render(<Pagination {...defaultProps} />);
      expect(screen.getByTestId('pagination-test')).toBeInTheDocument();
    });

    test('displays current page correctly', () => {
      render(<Pagination {...defaultProps} currentPage={5} />);
      const currentPageButton = screen.getByRole('button', { name: '5' });
      expect(currentPageButton).toHaveAttribute('aria-current', 'page');
    });

    test('applies custom className', () => {
      render(<Pagination {...defaultProps} className="custom-pagination" />);
      const pagination = screen.getByTestId('pagination-test');
      expect(pagination).toHaveClass('custom-pagination');
    });
  });

  describe('Page Navigation', () => {
    test('calls onPageChange when page button is clicked', () => {
      const onPageChange = vi.fn();
      render(<Pagination {...defaultProps} onPageChange={onPageChange} currentPage={1} />);

      const page2Button = screen.getByRole('button', { name: '2' });
      fireEvent.click(page2Button);
      expect(onPageChange).toHaveBeenCalledWith(2);
    });

    test('does not call onPageChange for current page', () => {
      const onPageChange = vi.fn();
      render(<Pagination {...defaultProps} onPageChange={onPageChange} currentPage={3} />);

      const currentPageButton = screen.getByRole('button', { name: '3' });
      fireEvent.click(currentPageButton);
      expect(onPageChange).not.toHaveBeenCalled();
    });
  });

  describe('Previous/Next Navigation', () => {
    test('shows previous/next buttons when showPrevNext is true', () => {
      render(<Pagination {...defaultProps} showPrevNext />);

      expect(screen.getByRole('button', { name: /previous/i })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /next/i })).toBeInTheDocument();
    });

    test('disables previous button on first page', () => {
      render(<Pagination {...defaultProps} currentPage={1} showPrevNext />);

      const prevButton = screen.getByRole('button', { name: /previous/i });
      expect(prevButton).toBeDisabled();
    });

    test('disables next button on last page', () => {
      render(<Pagination {...defaultProps} currentPage={10} totalPages={10} showPrevNext />);

      const nextButton = screen.getByRole('button', { name: /next/i });
      expect(nextButton).toBeDisabled();
    });

    test('calls onPageChange with correct page for previous/next', () => {
      const onPageChange = vi.fn();
      render(<Pagination {...defaultProps} currentPage={5} onPageChange={onPageChange} showPrevNext />);

      const prevButton = screen.getByRole('button', { name: /previous/i });
      const nextButton = screen.getByRole('button', { name: /next/i });

      fireEvent.click(prevButton);
      expect(onPageChange).toHaveBeenCalledWith(4);

      fireEvent.click(nextButton);
      expect(onPageChange).toHaveBeenCalledWith(6);
    });
  });

  describe('First/Last Navigation', () => {
    test('shows first/last buttons when showFirstLast is true', () => {
      render(<Pagination {...defaultProps} showFirstLast />);

      expect(screen.getByRole('button', { name: /first/i })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /last/i })).toBeInTheDocument();
    });

    test('disables first button on first page', () => {
      render(<Pagination {...defaultProps} currentPage={1} showFirstLast />);

      const firstButton = screen.getByRole('button', { name: /first/i });
      expect(firstButton).toBeDisabled();
    });

    test('disables last button on last page', () => {
      render(<Pagination {...defaultProps} currentPage={10} totalPages={10} showFirstLast />);

      const lastButton = screen.getByRole('button', { name: /last/i });
      expect(lastButton).toBeDisabled();
    });

    test('calls onPageChange with correct page for first/last', () => {
      const onPageChange = vi.fn();
      render(<Pagination {...defaultProps} currentPage={5} totalPages={10} onPageChange={onPageChange} showFirstLast />);

      const firstButton = screen.getByRole('button', { name: /first/i });
      const lastButton = screen.getByRole('button', { name: /last/i });

      fireEvent.click(firstButton);
      expect(onPageChange).toHaveBeenCalledWith(1);

      fireEvent.click(lastButton);
      expect(onPageChange).toHaveBeenCalledWith(10);
    });
  });

  describe('Visible Pages', () => {
    test('limits visible pages based on maxVisiblePages', () => {
      render(<Pagination {...defaultProps} totalPages={20} maxVisiblePages={5} currentPage={10} />);

      // Should show approximately 5 page buttons (may vary based on algorithm)
      const pageButtons = screen.getAllByRole('button').filter(button =>
        /^\d+$/.test(button.textContent || '')
      );
      expect(pageButtons.length).toBeLessThanOrEqual(5);
    });

    test('shows ellipsis when there are more pages', () => {
      render(<Pagination {...defaultProps} totalPages={20} maxVisiblePages={5} currentPage={10} />);

      const ellipsis = screen.getAllByText('...');
      expect(ellipsis.length).toBeGreaterThan(0);
    });
  });

  describe('Disabled State', () => {
    test('disables all buttons when disabled is true', () => {
      render(
        <Pagination
          {...defaultProps}
          disabled
          showPrevNext
          showFirstLast
        />
      );

      const buttons = screen.getAllByRole('button');
      buttons.forEach(button => {
        expect(button).toBeDisabled();
      });
    });
  });

  describe('Page Size Selector', () => {
    test('shows page size selector when showPageSize is true', () => {
      render(
        <Pagination
          {...defaultProps}
          showPageSize
          pageSizes={[10, 20, 50]}
          pageSize={10}
          onPageSizeChange={vi.fn()}
        />
      );

      expect(screen.getByRole('combobox')).toBeInTheDocument();
    });

    test('calls onPageSizeChange when page size changes', () => {
      const onPageSizeChange = vi.fn();
      render(
        <Pagination
          {...defaultProps}
          showPageSize
          pageSizes={[10, 20, 50]}
          pageSize={10}
          onPageSizeChange={onPageSizeChange}
        />
      );

      const select = screen.getByRole('combobox');
      fireEvent.change(select, { target: { value: '20' } });
      expect(onPageSizeChange).toHaveBeenCalledWith(20);
    });
  });

  describe('Accessibility', () => {
    test('supports ARIA label', () => {
      render(<Pagination {...defaultProps} aria-label="Search results pagination" />);
      const pagination = screen.getByTestId('pagination-test');
      expect(pagination).toHaveAttribute('aria-label', 'Search results pagination');
    });

    test('has proper navigation role', () => {
      render(<Pagination {...defaultProps} role="navigation" />);
      const pagination = screen.getByTestId('pagination-test');
      expect(pagination).toHaveAttribute('role', 'navigation');
    });

    test('supports keyboard navigation', () => {
      render(<Pagination {...defaultProps} tabIndex={0} />);
      const pagination = screen.getByTestId('pagination-test');
      expect(pagination).toHaveAttribute('tabIndex', '0');
    });
  });
});

describe('PaginationButton Contract Tests', () => {
  const defaultButtonProps: PaginationButtonProps = {
    page: 1,
    onClick: vi.fn(),
    testId: 'pagination-button-test'
  };

  test('renders button with page number', () => {
    render(<PaginationButton {...defaultButtonProps} page={5} />);
    expect(screen.getByRole('button', { name: '5' })).toBeInTheDocument();
  });

  test('renders button with action type', () => {
    render(<PaginationButton {...defaultButtonProps} page="next" />);
    expect(screen.getByRole('button', { name: /next/i })).toBeInTheDocument();
  });

  test('applies active state correctly', () => {
    render(<PaginationButton {...defaultButtonProps} active />);
    const button = screen.getByRole('button');
    expect(button).toHaveAttribute('aria-current', 'page');
  });

  test('applies disabled state correctly', () => {
    render(<PaginationButton {...defaultButtonProps} disabled />);
    const button = screen.getByRole('button');
    expect(button).toBeDisabled();
  });

  test('calls onClick when clicked', () => {
    const onClick = vi.fn();
    render(<PaginationButton {...defaultButtonProps} onClick={onClick} />);

    const button = screen.getByRole('button');
    fireEvent.click(button);
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});

describe('PageSizeSelector Contract Tests', () => {
  const defaultSelectorProps: PageSizeSelectorProps = {
    sizes: [10, 20, 50, 100],
    value: 20,
    onChange: vi.fn(),
    testId: 'page-size-selector-test'
  };

  test('renders select with all size options', () => {
    render(<PageSizeSelector {...defaultSelectorProps} />);

    const select = screen.getByRole('combobox');
    expect(select).toBeInTheDocument();

    defaultSelectorProps.sizes.forEach(size => {
      expect(screen.getByRole('option', { name: size.toString() })).toBeInTheDocument();
    });
  });

  test('shows current value as selected', () => {
    render(<PageSizeSelector {...defaultSelectorProps} value={50} />);

    const select = screen.getByRole('combobox') as HTMLSelectElement;
    expect(select.value).toBe('50');
  });

  test('calls onChange when selection changes', () => {
    const onChange = vi.fn();
    render(<PageSizeSelector {...defaultSelectorProps} onChange={onChange} />);

    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '100' } });
    expect(onChange).toHaveBeenCalledWith(100);
  });

  test('applies disabled state correctly', () => {
    render(<PageSizeSelector {...defaultSelectorProps} disabled />);

    const select = screen.getByRole('combobox');
    expect(select).toBeDisabled();
  });
});