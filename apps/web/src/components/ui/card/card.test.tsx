/**
 * Card Component Contract Tests
 * CRITICAL: These tests MUST FAIL until Card component is implemented
 */

import { render, screen, fireEvent } from '@testing-library/react';
import { Card, CardHeader, CardContent, CardFooter, CardMedia } from './card';
import type { CardProps, CardHeaderProps, CardContentProps, CardFooterProps, CardMediaProps } from '../../../../specs/001-agent-frontend-specialist/contracts/card';

describe('Card Contract Tests', () => {
  describe('Basic Rendering', () => {
    test('renders card with default props', () => {
      render(
        <Card testId="card-test">
          <CardContent>Card content</CardContent>
        </Card>
      );

      expect(screen.getByTestId('card-test')).toBeInTheDocument();
      expect(screen.getByText('Card content')).toBeInTheDocument();
    });

    test('applies custom className', () => {
      render(
        <Card className="custom-card" testId="card-custom">
          <CardContent>Content</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-custom');
      expect(card).toHaveClass('custom-card');
    });
  });

  describe('Variants', () => {
    const variants: CardProps['variant'][] = ['default', 'outline', 'ghost', 'filled'];

    variants.forEach(variant => {
      test(`renders ${variant} variant correctly`, () => {
        render(
          <Card variant={variant} testId={`card-${variant}`}>
            <CardContent>Content</CardContent>
          </Card>
        );

        const card = screen.getByTestId(`card-${variant}`);
        expect(card).toHaveClass(`card-${variant}`);
      });
    });
  });

  describe('Padding and Spacing', () => {
    const paddings: CardProps['padding'][] = ['none', 'sm', 'md', 'lg'];

    paddings.forEach(padding => {
      test(`applies ${padding} padding correctly`, () => {
        render(
          <Card padding={padding} testId={`card-padding-${padding}`}>
            <CardContent>Content</CardContent>
          </Card>
        );

        const card = screen.getByTestId(`card-padding-${padding}`);
        expect(card).toHaveClass(`card-padding-${padding}`);
      });
    });
  });

  describe('Shadow Levels', () => {
    const shadows: CardProps['shadow'][] = ['none', 'sm', 'md', 'lg'];

    shadows.forEach(shadow => {
      test(`applies ${shadow} shadow correctly`, () => {
        render(
          <Card shadow={shadow} testId={`card-shadow-${shadow}`}>
            <CardContent>Content</CardContent>
          </Card>
        );

        const card = screen.getByTestId(`card-shadow-${shadow}`);
        expect(card).toHaveClass(`card-shadow-${shadow}`);
      });
    });
  });

  describe('Interactive Card', () => {
    test('applies interactive styles when interactive is true', () => {
      render(
        <Card interactive testId="card-interactive">
          <CardContent>Clickable content</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-interactive');
      expect(card).toHaveClass('card-interactive');
    });

    test('calls onClick when interactive card is clicked', () => {
      const onClick = vi.fn();
      render(
        <Card interactive onClick={onClick} testId="card-clickable">
          <CardContent>Clickable content</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-clickable');
      fireEvent.click(card);
      expect(onClick).toHaveBeenCalledTimes(1);
    });

    test('supports keyboard interaction for interactive cards', () => {
      const onClick = vi.fn();
      render(
        <Card interactive onClick={onClick} testId="card-keyboard">
          <CardContent>Keyboard accessible</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-keyboard');
      fireEvent.keyDown(card, { key: 'Enter' });
      expect(onClick).toHaveBeenCalledTimes(1);

      fireEvent.keyDown(card, { key: ' ' });
      expect(onClick).toHaveBeenCalledTimes(2);
    });
  });

  describe('Loading State', () => {
    test('shows loading state correctly', () => {
      render(
        <Card loading testId="card-loading">
          <CardContent>Loading content</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-loading');
      expect(card).toHaveClass('card-loading');
      expect(card).toHaveAttribute('aria-busy', 'true');
    });
  });

  describe('Selected State', () => {
    test('shows selected state correctly', () => {
      render(
        <Card selected testId="card-selected">
          <CardContent>Selected content</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-selected');
      expect(card).toHaveClass('card-selected');
      expect(card).toHaveAttribute('aria-selected', 'true');
    });
  });

  describe('Accessibility', () => {
    test('supports ARIA label', () => {
      render(
        <Card aria-label="User profile card" testId="card-aria">
          <CardContent>Profile content</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-aria');
      expect(card).toHaveAttribute('aria-label', 'User profile card');
    });

    test('supports custom role', () => {
      render(
        <Card role="article" testId="card-role">
          <CardContent>Article content</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-role');
      expect(card).toHaveAttribute('role', 'article');
    });

    test('supports tabIndex for keyboard navigation', () => {
      render(
        <Card tabIndex={0} testId="card-tab">
          <CardContent>Focusable content</CardContent>
        </Card>
      );

      const card = screen.getByTestId('card-tab');
      expect(card).toHaveAttribute('tabIndex', '0');
    });
  });
});

describe('CardHeader Contract Tests', () => {
  test('renders header with title and subtitle', () => {
    render(
      <Card>
        <CardHeader
          title="Card Title"
          subtitle="Card Subtitle"
          testId="card-header-test"
        />
      </Card>
    );

    expect(screen.getByText('Card Title')).toBeInTheDocument();
    expect(screen.getByText('Card Subtitle')).toBeInTheDocument();
  });

  test('renders header with actions', () => {
    const actions = (
      <button data-testid="header-action">Action</button>
    );

    render(
      <Card>
        <CardHeader
          title="Title"
          actions={actions}
          testId="card-header-actions"
        />
      </Card>
    );

    expect(screen.getByTestId('header-action')).toBeInTheDocument();
  });

  test('renders header with avatar', () => {
    const avatar = <img data-testid="header-avatar" src="/avatar.jpg" alt="Avatar" />;

    render(
      <Card>
        <CardHeader
          title="Title"
          avatar={avatar}
          testId="card-header-avatar"
        />
      </Card>
    );

    expect(screen.getByTestId('header-avatar')).toBeInTheDocument();
  });

  test('applies border when border prop is true', () => {
    render(
      <Card>
        <CardHeader
          title="Title"
          border
          testId="card-header-border"
        />
      </Card>
    );

    const header = screen.getByTestId('card-header-border');
    expect(header).toHaveClass('card-header-border');
  });

  test('renders children when no title prop provided', () => {
    render(
      <Card>
        <CardHeader testId="card-header-children">
          <h3>Custom Header Content</h3>
        </CardHeader>
      </Card>
    );

    expect(screen.getByText('Custom Header Content')).toBeInTheDocument();
  });
});

describe('CardContent Contract Tests', () => {
  test('renders content with default props', () => {
    render(
      <Card>
        <CardContent testId="card-content-test">
          <p>Content text</p>
        </CardContent>
      </Card>
    );

    expect(screen.getByText('Content text')).toBeInTheDocument();
  });

  test('applies custom padding override', () => {
    const paddings: CardContentProps['padding'][] = ['inherit', 'none', 'sm', 'md', 'lg'];

    paddings.forEach(padding => {
      const { unmount } = render(
        <Card>
          <CardContent padding={padding} testId={`content-padding-${padding}`}>
            Content
          </CardContent>
        </Card>
      );

      if (padding !== 'inherit') {
        const content = screen.getByTestId(`content-padding-${padding}`);
        expect(content).toHaveClass(`card-content-padding-${padding}`);
      }
      unmount();
    });
  });

  test('supports scrollable content', () => {
    render(
      <Card>
        <CardContent scrollable testId="card-content-scroll">
          <div style={{ height: '200px' }}>Scrollable content</div>
        </CardContent>
      </Card>
    );

    const content = screen.getByTestId('card-content-scroll');
    expect(content).toHaveClass('card-content-scrollable');
  });

  test('applies maxHeight for scrollable content', () => {
    render(
      <Card>
        <CardContent scrollable maxHeight="300px" testId="card-content-max-height">
          Long content
        </CardContent>
      </Card>
    );

    const content = screen.getByTestId('card-content-max-height');
    expect(content).toHaveStyle({ maxHeight: '300px' });
  });
});

describe('CardFooter Contract Tests', () => {
  test('renders footer with children', () => {
    render(
      <Card>
        <CardFooter testId="card-footer-test">
          <button>Footer Button</button>
        </CardFooter>
      </Card>
    );

    expect(screen.getByText('Footer Button')).toBeInTheDocument();
  });

  test('applies justification classes', () => {
    const justifications: CardFooterProps['justify'][] = ['start', 'center', 'end', 'between', 'around'];

    justifications.forEach(justify => {
      const { unmount } = render(
        <Card>
          <CardFooter justify={justify} testId={`footer-justify-${justify}`}>
            Content
          </CardFooter>
        </Card>
      );

      const footer = screen.getByTestId(`footer-justify-${justify}`);
      expect(footer).toHaveClass(`card-footer-justify-${justify}`);
      unmount();
    });
  });

  test('applies border when border prop is true', () => {
    render(
      <Card>
        <CardFooter border testId="card-footer-border">
          Footer content
        </CardFooter>
      </Card>
    );

    const footer = screen.getByTestId('card-footer-border');
    expect(footer).toHaveClass('card-footer-border');
  });

  test('applies custom padding override', () => {
    const paddings: CardFooterProps['padding'][] = ['inherit', 'none', 'sm', 'md', 'lg'];

    paddings.forEach(padding => {
      const { unmount } = render(
        <Card>
          <CardFooter padding={padding} testId={`footer-padding-${padding}`}>
            Content
          </CardFooter>
        </Card>
      );

      if (padding !== 'inherit') {
        const footer = screen.getByTestId(`footer-padding-${padding}`);
        expect(footer).toHaveClass(`card-footer-padding-${padding}`);
      }
      unmount();
    });
  });
});

describe('CardMedia Contract Tests', () => {
  test('renders image media correctly', () => {
    render(
      <Card>
        <CardMedia
          src="/test-image.jpg"
          alt="Test image"
          type="image"
          testId="card-media-image"
        />
      </Card>
    );

    const media = screen.getByTestId('card-media-image');
    const image = screen.getByAltText('Test image');

    expect(media).toBeInTheDocument();
    expect(image).toHaveAttribute('src', '/test-image.jpg');
  });

  test('renders video media correctly', () => {
    render(
      <Card>
        <CardMedia
          src="/test-video.mp4"
          type="video"
          testId="card-media-video"
        />
      </Card>
    );

    const media = screen.getByTestId('card-media-video');
    expect(media).toBeInTheDocument();

    const video = media.querySelector('video');
    expect(video).toHaveAttribute('src', '/test-video.mp4');
  });

  test('applies aspect ratio classes', () => {
    const aspectRatios: CardMediaProps['aspectRatio'][] = ['square', 'video', 'auto'];

    aspectRatios.forEach(aspectRatio => {
      const { unmount } = render(
        <Card>
          <CardMedia
            src="/test.jpg"
            aspectRatio={aspectRatio}
            testId={`media-aspect-${aspectRatio}`}
          />
        </Card>
      );

      const media = screen.getByTestId(`media-aspect-${aspectRatio}`);
      expect(media).toHaveClass(`card-media-aspect-${aspectRatio}`);
      unmount();
    });
  });

  test('applies object fit classes', () => {
    const objectFits: CardMediaProps['objectFit'][] = ['cover', 'contain', 'fill', 'none'];

    objectFits.forEach(objectFit => {
      const { unmount } = render(
        <Card>
          <CardMedia
            src="/test.jpg"
            objectFit={objectFit}
            testId={`media-fit-${objectFit}`}
          />
        </Card>
      );

      const media = screen.getByTestId(`media-fit-${objectFit}`);
      expect(media).toHaveClass(`card-media-fit-${objectFit}`);
      unmount();
    });
  });
});