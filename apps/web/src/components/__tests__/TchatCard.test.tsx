/**
 * Unit Tests: TchatCard Variants (T059)
 * Tests component logic, variant behavior, and accessibility
 * Complements integration tests with focused unit testing
 * Constitutional requirements: 97% consistency, WCAG 2.1 AA, <200ms load
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, cleanup } from '@testing-library/react';
import {
  TchatCard,
  TchatCardHeader,
  TchatCardContent,
  TchatCardFooter,
  type TchatCardProps,
  type TchatCardVariant,
  type TchatCardSize
} from '../TchatCard';

// Clean up after each test
afterEach(cleanup);

describe('TchatCard Component Unit Tests', () => {
  describe('Component Rendering and Props', () => {
    it('should render with default props', () => {
      render(<TchatCard>Default Card</TchatCard>);

      const card = screen.getByTestId('tchat-card');
      expect(card).toBeInTheDocument();
      expect(card).toHaveTextContent('Default Card');
      expect(card).toHaveAttribute('data-variant', 'elevated');
      expect(card).toHaveAttribute('data-size', 'standard');
      expect(card).toHaveAttribute('data-interactive', 'false');
    });

    it('should forward ref correctly', () => {
      const ref = vi.fn();
      render(<TchatCard ref={ref}>Ref Card</TchatCard>);
      expect(ref).toHaveBeenCalledWith(expect.any(HTMLDivElement));
    });

    it('should apply custom className', () => {
      render(
        <TchatCard className="custom-class">
          Custom Class Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveClass('custom-class');
    });

    it('should pass through HTML attributes', () => {
      render(
        <TchatCard data-custom="test" id="test-card">
          HTML Attributes Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('data-custom', 'test');
      expect(card).toHaveAttribute('id', 'test-card');
    });
  });

  describe('Variant System', () => {
    const variants: TchatCardVariant[] = ['elevated', 'outlined', 'filled', 'glass'];

    it('should render all 4 variants correctly', () => {
      variants.forEach(variant => {
        const { container } = render(
          <TchatCard variant={variant} data-testid={`card-${variant}`}>
            {variant} Variant
          </TchatCard>
        );

        const card = screen.getByTestId(`card-${variant}`);
        expect(card).toBeInTheDocument();
        expect(card).toHaveAttribute('data-variant', variant);
        expect(card).toHaveTextContent(`${variant} Variant`);

        cleanup();
      });
    });

    it('should apply elevated variant styles', () => {
      const { container } = render(
        <TchatCard variant="elevated">Elevated Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('bg-white', 'shadow-sm', 'border');
    });

    it('should apply outlined variant styles', () => {
      const { container } = render(
        <TchatCard variant="outlined">Outlined Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('bg-white', 'border');
    });

    it('should apply filled variant styles', () => {
      const { container } = render(
        <TchatCard variant="filled">Filled Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('bg-surface', 'border');
    });

    it('should apply glass variant with glassmorphism styles', () => {
      const { container } = render(
        <TchatCard variant="glass">Glass Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass(
        'bg-white/80',
        'backdrop-blur-sm',
        'relative',
        'overflow-hidden'
      );
    });
  });

  describe('Size System', () => {
    const sizes: TchatCardSize[] = ['compact', 'standard', 'expanded'];

    it('should render all 3 size variants correctly', () => {
      sizes.forEach(size => {
        const { container } = render(
          <TchatCard size={size} data-testid={`card-${size}`}>
            {size} Size
          </TchatCard>
        );

        const card = screen.getByTestId(`card-${size}`);
        expect(card).toBeInTheDocument();
        expect(card).toHaveAttribute('data-size', size);

        cleanup();
      });
    });

    it('should apply compact size padding', () => {
      const { container } = render(
        <TchatCard size="compact">Compact Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('p-3'); // 12dp padding
    });

    it('should apply standard size padding', () => {
      const { container } = render(
        <TchatCard size="standard">Standard Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('p-4'); // 16dp padding
    });

    it('should apply expanded size padding', () => {
      const { container } = render(
        <TchatCard size="expanded">Expanded Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('p-6'); // 24dp padding
    });
  });

  describe('Interactive Behavior', () => {
    it('should not be interactive by default', () => {
      render(<TchatCard>Non-interactive Card</TchatCard>);

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('data-interactive', 'false');
      expect(card).toHaveAttribute('role', 'article');
      expect(card).not.toHaveAttribute('tabIndex');
      expect(card).not.toHaveClass('cursor-pointer');
    });

    it('should be interactive when interactive=true', () => {
      render(
        <TchatCard interactive>
          Interactive Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('data-interactive', 'true');
      expect(card).toHaveAttribute('role', 'button');
      expect(card).toHaveAttribute('tabIndex', '0');
      expect(card).toHaveClass('cursor-pointer');
    });

    it('should handle click events when interactive', () => {
      const handleClick = vi.fn();
      render(
        <TchatCard interactive onClick={handleClick}>
          Clickable Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      fireEvent.click(card);
      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('should not handle click events when not interactive', () => {
      const handleClick = vi.fn();
      render(
        <TchatCard interactive={false} onClick={handleClick}>
          Non-clickable Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      fireEvent.click(card);
      expect(handleClick).not.toHaveBeenCalled();
    });

    it('should handle Enter key when interactive', () => {
      const handleClick = vi.fn();
      render(
        <TchatCard interactive onClick={handleClick}>
          Keyboard Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      fireEvent.keyDown(card, { key: 'Enter', code: 'Enter' });
      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('should handle Space key when interactive', () => {
      const handleClick = vi.fn();
      render(
        <TchatCard interactive onClick={handleClick}>
          Space Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      fireEvent.keyDown(card, { key: ' ', code: 'Space' });
      expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('should not handle other keys when interactive', () => {
      const handleClick = vi.fn();
      render(
        <TchatCard interactive onClick={handleClick}>
          Other Keys Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      fireEvent.keyDown(card, { key: 'Tab', code: 'Tab' });
      fireEvent.keyDown(card, { key: 'Escape', code: 'Escape' });
      expect(handleClick).not.toHaveBeenCalled();
    });

    it('should handle custom onKeyDown callback', () => {
      const handleKeyDown = vi.fn();
      render(
        <TchatCard onKeyDown={handleKeyDown}>
          Custom KeyDown Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      fireEvent.keyDown(card, { key: 'Tab', code: 'Tab' });
      expect(handleKeyDown).toHaveBeenCalledTimes(1);
    });
  });

  describe('Accessibility Features', () => {
    it('should have proper semantic role for non-interactive cards', () => {
      render(<TchatCard>Article Card</TchatCard>);

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('role', 'article');
    });

    it('should have proper semantic role for interactive cards', () => {
      render(<TchatCard interactive>Button Card</TchatCard>);

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('role', 'button');
    });

    it('should support custom role', () => {
      render(
        <TchatCard role="region">
          Custom Role Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('role', 'region');
    });

    it('should support aria-label', () => {
      render(
        <TchatCard ariaLabel="Product card with pricing info">
          Product Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('aria-label', 'Product card with pricing info');
    });

    it('should support contentDescription as aria-label fallback', () => {
      render(
        <TchatCard contentDescription="Product description for screen readers">
          Product Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('aria-label', 'Product description for screen readers');
    });

    it('should prefer ariaLabel over contentDescription', () => {
      render(
        <TchatCard
          ariaLabel="Primary label"
          contentDescription="Fallback description"
        >
          Priority Card
        </TchatCard>
      );

      const card = screen.getByTestId('tchat-card');
      expect(card).toHaveAttribute('aria-label', 'Primary label');
    });

    it('should have focus ring styles for interactive cards', () => {
      const { container } = render(
        <TchatCard interactive>
          Focus Card
        </TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass(
        'focus:outline-none',
        'focus:ring-2',
        'focus:ring-primary'
      );
    });

    it('should meet minimum touch target size (44dp)', () => {
      const { container } = render(
        <TchatCard>Touch Target Card</TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      // Constitutional requirement: 44dp minimum touch targets
      expect(card).toHaveClass('min-h-[44px]');
    });
  });

  describe('Performance Features', () => {
    it('should have GPU acceleration classes', () => {
      const { container } = render(
        <TchatCard variant="elevated">
          Performance Card
        </TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('transform-gpu', 'will-change-transform');
    });

    it('should have transition classes for smooth animations', () => {
      const { container } = render(
        <TchatCard variant="outlined">
          Animated Card
        </TchatCard>
      );

      const card = container.firstChild as HTMLElement;
      expect(card).toHaveClass('transition-all', 'duration-200');
    });

    it('should render quickly within performance budget', () => {
      const startTime = performance.now();

      render(
        <TchatCard variant="glass" size="expanded" interactive>
          Performance Test Card
        </TchatCard>
      );

      const endTime = performance.now();
      const renderTime = endTime - startTime;

      // Constitutional requirement: <200ms component load times
      expect(renderTime).toBeLessThan(200);
    });
  });

  describe('Variant Combinations', () => {
    it('should handle all variant and size combinations', () => {
      const variants: TchatCardVariant[] = ['elevated', 'outlined', 'filled', 'glass'];
      const sizes: TchatCardSize[] = ['compact', 'standard', 'expanded'];

      variants.forEach(variant => {
        sizes.forEach(size => {
          const { container } = render(
            <TchatCard
              variant={variant}
              size={size}
              data-testid={`${variant}-${size}`}
            >
              {variant} {size}
            </TchatCard>
          );

          const card = screen.getByTestId(`${variant}-${size}`);
          expect(card).toBeInTheDocument();
          expect(card).toHaveAttribute('data-variant', variant);
          expect(card).toHaveAttribute('data-size', size);

          cleanup();
        });
      });
    });

    it('should handle interactive variants with all sizes', () => {
      const sizes: TchatCardSize[] = ['compact', 'standard', 'expanded'];

      sizes.forEach(size => {
        const { container } = render(
          <TchatCard
            variant="elevated"
            size={size}
            interactive
            data-testid={`interactive-${size}`}
          >
            Interactive {size}
          </TchatCard>
        );

        const card = screen.getByTestId(`interactive-${size}`);
        expect(card).toHaveAttribute('data-interactive', 'true');
        expect(card).toHaveAttribute('role', 'button');
        expect(card).toHaveClass('cursor-pointer');

        cleanup();
      });
    });
  });
});

describe('TchatCardHeader Component Unit Tests', () => {
  it('should render with title and subtitle', () => {
    render(
      <TchatCardHeader
        title="Card Title"
        subtitle="Card Subtitle"
      />
    );

    expect(screen.getByText('Card Title')).toBeInTheDocument();
    expect(screen.getByText('Card Subtitle')).toBeInTheDocument();
    expect(screen.getByTestId('card-header')).toBeInTheDocument();
  });

  it('should render with actions', () => {
    render(
      <TchatCardHeader
        title="Card with Actions"
        actions={<button>Action</button>}
      />
    );

    expect(screen.getByText('Card with Actions')).toBeInTheDocument();
    expect(screen.getByText('Action')).toBeInTheDocument();
  });

  it('should render children content', () => {
    render(
      <TchatCardHeader title="Header">
        <p>Custom content</p>
      </TchatCardHeader>
    );

    expect(screen.getByText('Custom content')).toBeInTheDocument();
  });

  it('should apply custom className', () => {
    render(
      <TchatCardHeader
        title="Custom Class Header"
        className="custom-header"
      />
    );

    const header = screen.getByTestId('card-header');
    expect(header).toHaveClass('custom-header');
  });

  it('should forward ref correctly', () => {
    const ref = vi.fn();
    render(<TchatCardHeader ref={ref} title="Ref Header" />);
    expect(ref).toHaveBeenCalledWith(expect.any(HTMLDivElement));
  });
});

describe('TchatCardContent Component Unit Tests', () => {
  it('should render children content', () => {
    render(
      <TchatCardContent>
        <p>Card content text</p>
      </TchatCardContent>
    );

    expect(screen.getByText('Card content text')).toBeInTheDocument();
    expect(screen.getByTestId('card-content')).toBeInTheDocument();
  });

  it('should apply custom className', () => {
    render(
      <TchatCardContent className="custom-content">
        Content
      </TchatCardContent>
    );

    const content = screen.getByTestId('card-content');
    expect(content).toHaveClass('custom-content');
  });

  it('should forward ref correctly', () => {
    const ref = vi.fn();
    render(<TchatCardContent ref={ref}>Content</TchatCardContent>);
    expect(ref).toHaveBeenCalledWith(expect.any(HTMLDivElement));
  });

  it('should have proper text styling', () => {
    const { container } = render(
      <TchatCardContent>Styled content</TchatCardContent>
    );

    const content = container.firstChild as HTMLElement;
    expect(content).toHaveClass('text-text-primary');
  });
});

describe('TchatCardFooter Component Unit Tests', () => {
  it('should render children content', () => {
    render(
      <TchatCardFooter>
        <span>Footer content</span>
      </TchatCardFooter>
    );

    expect(screen.getByText('Footer content')).toBeInTheDocument();
    expect(screen.getByTestId('card-footer')).toBeInTheDocument();
  });

  it('should apply custom className', () => {
    render(
      <TchatCardFooter className="custom-footer">
        Footer
      </TchatCardFooter>
    );

    const footer = screen.getByTestId('card-footer');
    expect(footer).toHaveClass('custom-footer');
  });

  it('should forward ref correctly', () => {
    const ref = vi.fn();
    render(<TchatCardFooter ref={ref}>Footer</TchatCardFooter>);
    expect(ref).toHaveBeenCalledWith(expect.any(HTMLDivElement));
  });

  it('should have proper footer styling', () => {
    const { container } = render(
      <TchatCardFooter>Footer with border</TchatCardFooter>
    );

    const footer = container.firstChild as HTMLElement;
    expect(footer).toHaveClass(
      'flex',
      'items-center',
      'justify-between',
      'border-t'
    );
  });
});

describe('Full Card Composition Unit Tests', () => {
  it('should render complete card composition', () => {
    render(
      <TchatCard variant="elevated" size="standard" interactive>
        <TchatCardHeader
          title="Product Card"
          subtitle="Premium Product"
          actions={<button>View</button>}
        />
        <TchatCardContent>
          <p>Product description and details</p>
        </TchatCardContent>
        <TchatCardFooter>
          <span>$99.99</span>
          <button>Add to Cart</button>
        </TchatCardFooter>
      </TchatCard>
    );

    // Verify all components are present
    expect(screen.getByTestId('tchat-card')).toBeInTheDocument();
    expect(screen.getByTestId('card-header')).toBeInTheDocument();
    expect(screen.getByTestId('card-content')).toBeInTheDocument();
    expect(screen.getByTestId('card-footer')).toBeInTheDocument();

    // Verify content
    expect(screen.getByText('Product Card')).toBeInTheDocument();
    expect(screen.getByText('Premium Product')).toBeInTheDocument();
    expect(screen.getByText('Product description and details')).toBeInTheDocument();
    expect(screen.getByText('$99.99')).toBeInTheDocument();
    expect(screen.getByText('Add to Cart')).toBeInTheDocument();
    expect(screen.getByText('View')).toBeInTheDocument();
  });

  it('should maintain proper component hierarchy', () => {
    const { container } = render(
      <TchatCard>
        <TchatCardHeader title="Header" />
        <TchatCardContent>Content</TchatCardContent>
        <TchatCardFooter>Footer</TchatCardFooter>
      </TchatCard>
    );

    const card = container.firstChild as HTMLElement;
    const header = card.querySelector('[data-testid="card-header"]');
    const content = card.querySelector('[data-testid="card-content"]');
    const footer = card.querySelector('[data-testid="card-footer"]');

    expect(header).toBeInTheDocument();
    expect(content).toBeInTheDocument();
    expect(footer).toBeInTheDocument();

    // Verify DOM order
    expect(header?.nextElementSibling).toBe(content);
    expect(content?.nextElementSibling).toBe(footer);
  });
});