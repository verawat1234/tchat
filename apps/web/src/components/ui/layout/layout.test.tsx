/**
 * Layout Component Contract Tests
 * CRITICAL: These tests MUST FAIL until Layout component is implemented
 */

import { render, screen } from '@testing-library/react';
import { Container, Grid, Flex, Stack, Spacer, Divider } from './layout';
import type {
  ContainerProps,
  GridProps,
  FlexProps,
  StackProps,
  SpacerProps,
  DividerProps
} from '../../../../specs/001-agent-frontend-specialist/contracts/layout';

describe('Container Contract Tests', () => {
  test('renders container with default props', () => {
    render(
      <Container testId="container-test">
        Container content
      </Container>
    );

    expect(screen.getByTestId('container-test')).toBeInTheDocument();
    expect(screen.getByText('Container content')).toBeInTheDocument();
  });

  test('applies maxWidth classes correctly', () => {
    const maxWidths: ContainerProps['maxWidth'][] = ['xs', 'sm', 'md', 'lg', 'xl', '2xl', 'full'];

    maxWidths.forEach(maxWidth => {
      const { unmount } = render(
        <Container maxWidth={maxWidth} testId={`container-${maxWidth}`}>
          Content
        </Container>
      );

      const container = screen.getByTestId(`container-${maxWidth}`);
      expect(container).toHaveClass(`container-${maxWidth}`);
      unmount();
    });
  });

  test('applies center styling when center is true', () => {
    render(
      <Container center testId="container-center">
        Centered content
      </Container>
    );

    const container = screen.getByTestId('container-center');
    expect(container).toHaveClass('container-center');
  });

  test('applies fluid styling when fluid is true', () => {
    render(
      <Container fluid testId="container-fluid">
        Fluid content
      </Container>
    );

    const container = screen.getByTestId('container-fluid');
    expect(container).toHaveClass('container-fluid');
  });

  test('applies padding correctly', () => {
    const paddings: ContainerProps['padding'][] = ['none', 'xs', 'sm', 'md', 'lg', 'xl'];

    paddings.forEach(padding => {
      const { unmount } = render(
        <Container padding={padding} testId={`container-padding-${padding}`}>
          Content
        </Container>
      );

      const container = screen.getByTestId(`container-padding-${padding}`);
      expect(container).toHaveClass(`container-padding-${padding}`);
      unmount();
    });
  });

  test('supports numeric maxWidth', () => {
    render(
      <Container maxWidth={800} testId="container-numeric">
        Content
      </Container>
    );

    const container = screen.getByTestId('container-numeric');
    expect(container).toHaveStyle({ maxWidth: '800px' });
  });
});

describe('Grid Contract Tests', () => {
  test('renders grid with default props', () => {
    render(
      <Grid testId="grid-test">
        <div>Grid item 1</div>
        <div>Grid item 2</div>
      </Grid>
    );

    expect(screen.getByTestId('grid-test')).toBeInTheDocument();
    expect(screen.getByText('Grid item 1')).toBeInTheDocument();
    expect(screen.getByText('Grid item 2')).toBeInTheDocument();
  });

  test('applies cols correctly', () => {
    render(
      <Grid cols={3} testId="grid-cols">
        <div>Item 1</div>
        <div>Item 2</div>
        <div>Item 3</div>
      </Grid>
    );

    const grid = screen.getByTestId('grid-cols');
    expect(grid).toHaveClass('grid-cols-3');
  });

  test('applies auto cols correctly', () => {
    const autoTypes: GridProps['cols'][] = ['auto', 'fit', 'fill'];

    autoTypes.forEach(type => {
      const { unmount } = render(
        <Grid cols={type} testId={`grid-cols-${type}`}>
          <div>Item</div>
        </Grid>
      );

      const grid = screen.getByTestId(`grid-cols-${type}`);
      expect(grid).toHaveClass(`grid-cols-${type}`);
      unmount();
    });
  });

  test('applies rows correctly', () => {
    render(
      <Grid rows={2} testId="grid-rows">
        <div>Item 1</div>
        <div>Item 2</div>
      </Grid>
    );

    const grid = screen.getByTestId('grid-rows');
    expect(grid).toHaveClass('grid-rows-2');
  });

  test('applies gap correctly', () => {
    const gaps: GridProps['gap'][] = ['none', 'xs', 'sm', 'md', 'lg', 'xl'];

    gaps.forEach(gap => {
      const { unmount } = render(
        <Grid gap={gap} testId={`grid-gap-${gap}`}>
          <div>Item</div>
        </Grid>
      );

      const grid = screen.getByTestId(`grid-gap-${gap}`);
      expect(grid).toHaveClass(`grid-gap-${gap}`);
      unmount();
    });
  });

  test('applies separate X and Y gaps', () => {
    render(
      <Grid gapX="md" gapY="lg" testId="grid-gap-xy">
        <div>Item</div>
      </Grid>
    );

    const grid = screen.getByTestId('grid-gap-xy');
    expect(grid).toHaveClass('grid-gap-x-md');
    expect(grid).toHaveClass('grid-gap-y-lg');
  });

  test('supports numeric gap values', () => {
    render(
      <Grid gap={16} testId="grid-gap-numeric">
        <div>Item</div>
      </Grid>
    );

    const grid = screen.getByTestId('grid-gap-numeric');
    expect(grid).toHaveStyle({ gap: '16px' });
  });

  test('supports responsive grid properties', () => {
    render(
      <Grid responsive={{ sm: { cols: 1 }, md: { cols: 2 }, lg: { cols: 3 } }} testId="grid-responsive">
        <div>Item 1</div>
        <div>Item 2</div>
        <div>Item 3</div>
      </Grid>
    );

    const grid = screen.getByTestId('grid-responsive');
    expect(grid).toHaveClass('grid-sm-cols-1');
    expect(grid).toHaveClass('grid-md-cols-2');
    expect(grid).toHaveClass('grid-lg-cols-3');
  });
});

describe('Flex Contract Tests', () => {
  test('renders flex container with default props', () => {
    render(
      <Flex testId="flex-test">
        <div>Flex item 1</div>
        <div>Flex item 2</div>
      </Flex>
    );

    expect(screen.getByTestId('flex-test')).toBeInTheDocument();
    expect(screen.getByText('Flex item 1')).toBeInTheDocument();
    expect(screen.getByText('Flex item 2')).toBeInTheDocument();
  });

  test('applies direction correctly', () => {
    const directions: FlexProps['direction'][] = ['row', 'column', 'row-reverse', 'column-reverse'];

    directions.forEach(direction => {
      const { unmount } = render(
        <Flex direction={direction} testId={`flex-${direction}`}>
          <div>Item</div>
        </Flex>
      );

      const flex = screen.getByTestId(`flex-${direction}`);
      expect(flex).toHaveClass(`flex-${direction}`);
      unmount();
    });
  });

  test('applies justify correctly', () => {
    const justifications: FlexProps['justify'][] = ['start', 'center', 'end', 'between', 'around', 'evenly'];

    justifications.forEach(justify => {
      const { unmount } = render(
        <Flex justify={justify} testId={`flex-justify-${justify}`}>
          <div>Item</div>
        </Flex>
      );

      const flex = screen.getByTestId(`flex-justify-${justify}`);
      expect(flex).toHaveClass(`flex-justify-${justify}`);
      unmount();
    });
  });

  test('applies align correctly', () => {
    const alignments: FlexProps['align'][] = ['start', 'center', 'end', 'stretch', 'baseline'];

    alignments.forEach(align => {
      const { unmount } = render(
        <Flex align={align} testId={`flex-align-${align}`}>
          <div>Item</div>
        </Flex>
      );

      const flex = screen.getByTestId(`flex-align-${align}`);
      expect(flex).toHaveClass(`flex-align-${align}`);
      unmount();
    });
  });

  test('applies wrap correctly', () => {
    render(
      <Flex wrap testId="flex-wrap">
        <div>Item</div>
      </Flex>
    );

    const flex = screen.getByTestId('flex-wrap');
    expect(flex).toHaveClass('flex-wrap');
  });

  test('applies wrap reverse correctly', () => {
    render(
      <Flex wrap="reverse" testId="flex-wrap-reverse">
        <div>Item</div>
      </Flex>
    );

    const flex = screen.getByTestId('flex-wrap-reverse');
    expect(flex).toHaveClass('flex-wrap-reverse');
  });

  test('applies gap correctly', () => {
    const gaps: FlexProps['gap'][] = ['none', 'xs', 'sm', 'md', 'lg', 'xl'];

    gaps.forEach(gap => {
      const { unmount } = render(
        <Flex gap={gap} testId={`flex-gap-${gap}`}>
          <div>Item</div>
        </Flex>
      );

      const flex = screen.getByTestId(`flex-gap-${gap}`);
      expect(flex).toHaveClass(`flex-gap-${gap}`);
      unmount();
    });
  });

  test('supports numeric gap', () => {
    render(
      <Flex gap={16} testId="flex-gap-numeric">
        <div>Item</div>
      </Flex>
    );

    const flex = screen.getByTestId('flex-gap-numeric');
    expect(flex).toHaveStyle({ gap: '16px' });
  });

  test('applies grow when grow is true', () => {
    render(
      <Flex grow testId="flex-grow">
        <div>Item</div>
      </Flex>
    );

    const flex = screen.getByTestId('flex-grow');
    expect(flex).toHaveClass('flex-grow');
  });

  test('applies shrink when shrink is true', () => {
    render(
      <Flex shrink testId="flex-shrink">
        <div>Item</div>
      </Flex>
    );

    const flex = screen.getByTestId('flex-shrink');
    expect(flex).toHaveClass('flex-shrink');
  });
});

describe('Stack Contract Tests', () => {
  test('renders stack with default props', () => {
    render(
      <Stack testId="stack-test">
        <div>Stack item 1</div>
        <div>Stack item 2</div>
      </Stack>
    );

    expect(screen.getByTestId('stack-test')).toBeInTheDocument();
    expect(screen.getByText('Stack item 1')).toBeInTheDocument();
    expect(screen.getByText('Stack item 2')).toBeInTheDocument();
  });

  test('applies space correctly', () => {
    const spaces: StackProps['space'][] = ['none', 'xs', 'sm', 'md', 'lg', 'xl'];

    spaces.forEach(space => {
      const { unmount } = render(
        <Stack space={space} testId={`stack-space-${space}`}>
          <div>Item 1</div>
          <div>Item 2</div>
        </Stack>
      );

      const stack = screen.getByTestId(`stack-space-${space}`);
      expect(stack).toHaveClass(`stack-space-${space}`);
      unmount();
    });
  });

  test('supports numeric space', () => {
    render(
      <Stack space={24} testId="stack-space-numeric">
        <div>Item 1</div>
        <div>Item 2</div>
      </Stack>
    );

    const stack = screen.getByTestId('stack-space-numeric');
    expect(stack).toHaveStyle({ gap: '24px' });
  });

  test('renders divider between items', () => {
    const divider = <hr data-testid="stack-divider" />;

    render(
      <Stack divider={divider} testId="stack-with-divider">
        <div>Item 1</div>
        <div>Item 2</div>
      </Stack>
    );

    expect(screen.getAllByTestId('stack-divider')).toHaveLength(1);
  });

  test('applies align correctly', () => {
    const alignments: StackProps['align'][] = ['start', 'center', 'end', 'stretch'];

    alignments.forEach(align => {
      const { unmount } = render(
        <Stack align={align} testId={`stack-align-${align}`}>
          <div>Item</div>
        </Stack>
      );

      const stack = screen.getByTestId(`stack-align-${align}`);
      expect(stack).toHaveClass(`stack-align-${align}`);
      unmount();
    });
  });
});

describe('Spacer Contract Tests', () => {
  test('renders spacer with default props', () => {
    render(<Spacer testId="spacer-test" />);

    expect(screen.getByTestId('spacer-test')).toBeInTheDocument();
  });

  test('applies size correctly', () => {
    const sizes: SpacerProps['size'][] = ['xs', 'sm', 'md', 'lg', 'xl'];

    sizes.forEach(size => {
      const { unmount } = render(
        <Spacer size={size} testId={`spacer-${size}`} />
      );

      const spacer = screen.getByTestId(`spacer-${size}`);
      expect(spacer).toHaveClass(`spacer-${size}`);
      unmount();
    });
  });

  test('supports numeric size', () => {
    render(
      <Spacer size={32} testId="spacer-numeric" />
    );

    const spacer = screen.getByTestId('spacer-numeric');
    expect(spacer).toHaveStyle({ height: '32px' });
  });

  test('applies direction correctly', () => {
    const directions: SpacerProps['direction'][] = ['horizontal', 'vertical'];

    directions.forEach(direction => {
      const { unmount } = render(
        <Spacer direction={direction} testId={`spacer-${direction}`} />
      );

      const spacer = screen.getByTestId(`spacer-${direction}`);
      expect(spacer).toHaveClass(`spacer-${direction}`);
      unmount();
    });
  });
});

describe('Divider Contract Tests', () => {
  test('renders divider with default props', () => {
    render(<Divider testId="divider-test" />);

    expect(screen.getByTestId('divider-test')).toBeInTheDocument();
  });

  test('applies orientation correctly', () => {
    const orientations: DividerProps['orientation'][] = ['horizontal', 'vertical'];

    orientations.forEach(orientation => {
      const { unmount } = render(
        <Divider orientation={orientation} testId={`divider-${orientation}`} />
      );

      const divider = screen.getByTestId(`divider-${orientation}`);
      expect(divider).toHaveClass(`divider-${orientation}`);
      unmount();
    });
  });

  test('applies variant correctly', () => {
    const variants: DividerProps['variant'][] = ['solid', 'dashed', 'dotted'];

    variants.forEach(variant => {
      const { unmount } = render(
        <Divider variant={variant} testId={`divider-${variant}`} />
      );

      const divider = screen.getByTestId(`divider-${variant}`);
      expect(divider).toHaveClass(`divider-${variant}`);
      unmount();
    });
  });

  test('applies thickness correctly', () => {
    const thicknesses: DividerProps['thickness'][] = ['thin', 'medium', 'thick'];

    thicknesses.forEach(thickness => {
      const { unmount } = render(
        <Divider thickness={thickness} testId={`divider-${thickness}`} />
      );

      const divider = screen.getByTestId(`divider-${thickness}`);
      expect(divider).toHaveClass(`divider-${thickness}`);
      unmount();
    });
  });

  test('renders label when provided', () => {
    render(
      <Divider label="Section Break" testId="divider-label" />
    );

    expect(screen.getByText('Section Break')).toBeInTheDocument();
  });

  test('applies label position correctly', () => {
    const positions: DividerProps['labelPosition'][] = ['left', 'center', 'right'];

    positions.forEach(position => {
      const { unmount } = render(
        <Divider label="Label" labelPosition={position} testId={`divider-label-${position}`} />
      );

      const divider = screen.getByTestId(`divider-label-${position}`);
      expect(divider).toHaveClass(`divider-label-${position}`);
      unmount();
    });
  });
});