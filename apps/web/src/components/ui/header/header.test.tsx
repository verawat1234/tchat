/**
 * Header Component Contract Tests
 * CRITICAL: These tests MUST FAIL until Header component is implemented
 */

import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { Header, Breadcrumb, PageHeader } from './header';
import { api } from '../../../services/api';
import type {
  HeaderProps,
  BreadcrumbProps,
  PageHeaderProps,
  BreadcrumbItem
} from '../../../../specs/001-agent-frontend-specialist/contracts/header';

// Mock store setup for dynamic content tests
const createMockStore = (contentState = {}) => {
  return configureStore({
    reducer: {
      api: api.reducer,
      content: (state = {
        selectedLanguage: 'en',
        fallbackContent: {
          'header.navigation.back': { type: 'text', value: 'Go back' },
          'header.breadcrumb.navigation': { type: 'text', value: 'Breadcrumb' }
        },
        syncStatus: 'idle',
        fallbackMode: false,
        lastSyncTime: new Date().toISOString(),
        contentPreferences: { showDrafts: false, compactView: false },
        ...contentState
      }, action) => state,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(api.middleware),
  });
};

const renderWithProvider = (component: React.ReactElement, contentState = {}) => {
  const store = createMockStore(contentState);
  return render(
    <Provider store={store}>
      {component}
    </Provider>
  );
};

describe('Header Contract Tests', () => {
  const baseProps: HeaderProps = {
    title: 'Page Title',
    testId: 'header-test'
  };

  describe('Basic Rendering', () => {
    test('renders header with required title', () => {
      render(<Header {...baseProps} />);

      expect(screen.getByTestId('header-test')).toBeInTheDocument();
      expect(screen.getByText('Page Title')).toBeInTheDocument();
    });

    test('renders with subtitle when provided', () => {
      render(<Header {...baseProps} subtitle="Page Subtitle" />);

      expect(screen.getByText('Page Title')).toBeInTheDocument();
      expect(screen.getByText('Page Subtitle')).toBeInTheDocument();
    });

    test('applies custom className', () => {
      render(<Header {...baseProps} className="custom-header" />);

      const header = screen.getByTestId('header-test');
      expect(header).toHaveClass('custom-header');
    });
  });

  describe('Header Levels', () => {
    test('renders correct heading levels', () => {
      const levels: HeaderProps['level'][] = [1, 2, 3, 4, 5, 6];

      levels.forEach(level => {
        const { unmount } = render(
          <Header {...baseProps} level={level} testId={`header-h${level}`} />
        );

        const heading = screen.getByRole('heading', { level });
        expect(heading).toBeInTheDocument();
        expect(heading).toHaveTextContent('Page Title');
        unmount();
      });
    });
  });

  describe('Size Variants', () => {
    test('applies size classes correctly', () => {
      const sizes: HeaderProps['size'][] = ['sm', 'md', 'lg', 'xl'];

      sizes.forEach(size => {
        const { unmount } = render(
          <Header {...baseProps} size={size} testId={`header-${size}`} />
        );

        const header = screen.getByTestId(`header-${size}`);
        expect(header).toHaveClass(`header-${size}`);
        unmount();
      });
    });
  });

  describe('Actions and Navigation', () => {
    test('renders actions when provided', () => {
      const actions = (
        <button data-testid="header-action">Action Button</button>
      );

      render(<Header {...baseProps} actions={actions} />);

      expect(screen.getByTestId('header-action')).toBeInTheDocument();
    });

    test('renders breadcrumbs when provided', () => {
      const breadcrumbs: BreadcrumbItem[] = [
        { label: 'Home', href: '/' },
        { label: 'Products', href: '/products' },
        { label: 'Current Page' }
      ];

      render(<Header {...baseProps} breadcrumbs={breadcrumbs} />);

      expect(screen.getByText('Home')).toBeInTheDocument();
      expect(screen.getByText('Products')).toBeInTheDocument();
      expect(screen.getByText('Current Page')).toBeInTheDocument();
    });
  });

  describe('Visual Features', () => {
    test('applies sticky styling when sticky is true', () => {
      render(<Header {...baseProps} sticky testId="sticky-header" />);

      const header = screen.getByTestId('sticky-header');
      expect(header).toHaveClass('header-sticky');
    });

    test('applies border when border is true', () => {
      render(<Header {...baseProps} border testId="border-header" />);

      const header = screen.getByTestId('border-header');
      expect(header).toHaveClass('header-border');
    });

    test('applies centered styling when centered is true', () => {
      render(<Header {...baseProps} centered testId="centered-header" />);

      const header = screen.getByTestId('centered-header');
      expect(header).toHaveClass('header-centered');
    });

    test('applies background variants correctly', () => {
      const backgrounds: HeaderProps['background'][] = ['transparent', 'default', 'muted'];

      backgrounds.forEach(background => {
        const { unmount } = render(
          <Header {...baseProps} background={background} testId={`header-bg-${background}`} />
        );

        const header = screen.getByTestId(`header-bg-${background}`);
        expect(header).toHaveClass(`header-bg-${background}`);
        unmount();
      });
    });

    test('renders icon when provided', () => {
      const icon = <span data-testid="header-icon">📄</span>;

      render(<Header {...baseProps} icon={icon} />);

      expect(screen.getByTestId('header-icon')).toBeInTheDocument();
    });
  });

  describe('Back Navigation', () => {
    test('shows back button when showBack is true', () => {
      render(<Header {...baseProps} showBack />);

      const backButton = screen.getByRole('button', { name: /back/i });
      expect(backButton).toBeInTheDocument();
    });

    test('calls onBack when back button is clicked', () => {
      const onBack = vi.fn();
      render(<Header {...baseProps} showBack onBack={onBack} />);

      const backButton = screen.getByRole('button', { name: /back/i });
      fireEvent.click(backButton);
      expect(onBack).toHaveBeenCalledTimes(1);
    });
  });

  describe('Accessibility', () => {
    test('supports ARIA label', () => {
      render(<Header {...baseProps} aria-label="Main page header" />);

      const header = screen.getByTestId('header-test');
      expect(header).toHaveAttribute('aria-label', 'Main page header');
    });

    test('supports custom role', () => {
      render(<Header {...baseProps} role="banner" />);

      const header = screen.getByTestId('header-test');
      expect(header).toHaveAttribute('role', 'banner');
    });

    test('supports keyboard navigation', () => {
      render(<Header {...baseProps} tabIndex={0} />);

      const header = screen.getByTestId('header-test');
      expect(header).toHaveAttribute('tabIndex', '0');
    });
  });
});

describe('Breadcrumb Contract Tests', () => {
  const breadcrumbItems: BreadcrumbItem[] = [
    { label: 'Home', href: '/' },
    { label: 'Products', href: '/products' },
    { label: 'Category', href: '/products/category' },
    { label: 'Current Item' }
  ];

  test('renders all breadcrumb items', () => {
    render(<Breadcrumb items={breadcrumbItems} testId="breadcrumb-test" />);

    expect(screen.getByTestId('breadcrumb-test')).toBeInTheDocument();
    breadcrumbItems.forEach(item => {
      expect(screen.getByText(item.label)).toBeInTheDocument();
    });
  });

  test('renders custom separator when provided', () => {
    const separator = <span data-testid="custom-separator">→</span>;

    render(<Breadcrumb items={breadcrumbItems} separator={separator} testId="breadcrumb-separator" />);

    const separators = screen.getAllByTestId('custom-separator');
    expect(separators.length).toBe(breadcrumbItems.length - 1);
  });

  test('limits items when maxItems is set', () => {
    render(<Breadcrumb items={breadcrumbItems} maxItems={2} testId="breadcrumb-limited" />);

    expect(screen.getByText('Home')).toBeInTheDocument();
    expect(screen.getByText('Current Item')).toBeInTheDocument();
    expect(screen.getByText('...')).toBeInTheDocument();
  });

  test('shows home icon when showHome is true', () => {
    render(<Breadcrumb items={breadcrumbItems} showHome testId="breadcrumb-home" />);

    const homeIcon = screen.getByTestId('breadcrumb-home').querySelector('[data-icon="home"]');
    expect(homeIcon).toBeInTheDocument();
  });

  test('applies size variants correctly', () => {
    const sizes: BreadcrumbProps['size'][] = ['sm', 'md', 'lg'];

    sizes.forEach(size => {
      const { unmount } = render(
        <Breadcrumb items={breadcrumbItems} size={size} testId={`breadcrumb-${size}`} />
      );

      const breadcrumb = screen.getByTestId(`breadcrumb-${size}`);
      expect(breadcrumb).toHaveClass(`breadcrumb-${size}`);
      unmount();
    });
  });

  test('handles item clicks', () => {
    const onClick = vi.fn();
    const itemsWithClick = [
      { label: 'Home', onClick },
      { label: 'Current' }
    ];

    render(<Breadcrumb items={itemsWithClick} testId="breadcrumb-click" />);

    fireEvent.click(screen.getByText('Home'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  test('renders disabled items correctly', () => {
    const itemsWithDisabled = [
      { label: 'Home', href: '/' },
      { label: 'Disabled', disabled: true },
      { label: 'Current' }
    ];

    render(<Breadcrumb items={itemsWithDisabled} testId="breadcrumb-disabled" />);

    const disabledItem = screen.getByText('Disabled');
    expect(disabledItem).toHaveAttribute('aria-disabled', 'true');
  });

  test('renders item icons when provided', () => {
    const itemsWithIcons = [
      { label: 'Home', icon: <span data-testid="home-icon">🏠</span> },
      { label: 'Current' }
    ];

    render(<Breadcrumb items={itemsWithIcons} testId="breadcrumb-icons" />);

    expect(screen.getByTestId('home-icon')).toBeInTheDocument();
  });
});

describe('PageHeader Contract Tests', () => {
  const basePageHeaderProps: PageHeaderProps = {
    title: 'Page Title',
    testId: 'page-header-test'
  };

  test('renders page header with title', () => {
    render(<PageHeader {...basePageHeaderProps} />);

    expect(screen.getByTestId('page-header-test')).toBeInTheDocument();
    expect(screen.getByText('Page Title')).toBeInTheDocument();
  });

  test('renders description when provided', () => {
    render(
      <PageHeader
        {...basePageHeaderProps}
        description="This is a page description"
      />
    );

    expect(screen.getByText('This is a page description')).toBeInTheDocument();
  });

  test('renders primary and secondary actions', () => {
    const actions = <button data-testid="primary-action">Primary</button>;
    const secondaryActions = <button data-testid="secondary-action">Secondary</button>;

    render(
      <PageHeader
        {...basePageHeaderProps}
        actions={actions}
        secondaryActions={secondaryActions}
      />
    );

    expect(screen.getByTestId('primary-action')).toBeInTheDocument();
    expect(screen.getByTestId('secondary-action')).toBeInTheDocument();
  });

  test('renders tabs when provided', () => {
    const tabs = <div data-testid="page-tabs">Tab Navigation</div>;

    render(<PageHeader {...basePageHeaderProps} tabs={tabs} />);

    expect(screen.getByTestId('page-tabs')).toBeInTheDocument();
  });

  test('applies fullWidth styling when fullWidth is true', () => {
    render(<PageHeader {...basePageHeaderProps} fullWidth testId="full-width-header" />);

    const header = screen.getByTestId('full-width-header');
    expect(header).toHaveClass('page-header-full-width');
  });

  test('renders background when provided', () => {
    const background = <div data-testid="header-background">Background</div>;

    render(<PageHeader {...basePageHeaderProps} background={background} />);

    expect(screen.getByTestId('header-background')).toBeInTheDocument();
  });

  test('renders avatar when provided', () => {
    const avatar = <img data-testid="page-avatar" src="/avatar.jpg" alt="Avatar" />;

    render(<PageHeader {...basePageHeaderProps} avatar={avatar} />);

    expect(screen.getByTestId('page-avatar')).toBeInTheDocument();
  });

  test('renders status indicator when provided', () => {
    const status = <span data-testid="page-status">Online</span>;

    render(<PageHeader {...basePageHeaderProps} status={status} />);

    expect(screen.getByTestId('page-status')).toBeInTheDocument();
  });

  test('renders metadata when provided', () => {
    const metadata = <div data-testid="page-metadata">Created: 2023-01-01</div>;

    render(<PageHeader {...basePageHeaderProps} metadata={metadata} />);

    expect(screen.getByTestId('page-metadata')).toBeInTheDocument();
  });

  test('renders breadcrumbs in page header', () => {
    const breadcrumbs: BreadcrumbItem[] = [
      { label: 'Home', href: '/' },
      { label: 'Current Page' }
    ];

    render(<PageHeader {...basePageHeaderProps} breadcrumbs={breadcrumbs} />);

    expect(screen.getByText('Home')).toBeInTheDocument();
    expect(screen.getByText('Current Page')).toBeInTheDocument();
  });
});

describe('Header Dynamic Content Integration Tests', () => {
  describe('Content API Integration', () => {
    test('uses dynamic back button text from content API', async () => {
      renderWithProvider(
        <Header title="Test Page" showBack onBack={() => {}} />
      );

      await waitFor(() => {
        const backButton = screen.getByRole('button');
        expect(backButton).toHaveAttribute('aria-label', 'Go back');
      });
    });

    test('shows loading state while content loads', () => {
      renderWithProvider(
        <Header title="Test Page" showBack onBack={() => {}} />
      );

      // Check for loading indicator
      const header = screen.getByRole('banner');
      const loadingIndicator = header.querySelector('.animate-pulse');
      expect(loadingIndicator).toBeInTheDocument();
    });

    test('falls back to default text when content fails', async () => {
      renderWithProvider(
        <Header title="Test Page" showBack onBack={() => {}} />,
        { fallbackContent: {} } // Empty fallback content
      );

      await waitFor(() => {
        const backButton = screen.getByRole('button');
        expect(backButton).toHaveAttribute('aria-label', 'Go back');
      });
    });
  });

  describe('Language Support', () => {
    test('loads language-specific content', async () => {
      renderWithProvider(
        <Header title="Test Page" showBack onBack={() => {}} />,
        {
          selectedLanguage: 'es',
          fallbackContent: {
            'header.navigation.back.es': { type: 'text', value: 'Regresar' }
          }
        }
      );

      await waitFor(() => {
        const backButton = screen.getByRole('button');
        expect(backButton).toHaveAttribute('aria-label', 'Regresar');
      });
    });

    test('falls back to base language when specific language not available', async () => {
      renderWithProvider(
        <Header title="Test Page" showBack onBack={() => {}} />,
        {
          selectedLanguage: 'fr',
          fallbackContent: {
            'header.navigation.back': { type: 'text', value: 'Go back' }
          }
        }
      );

      await waitFor(() => {
        const backButton = screen.getByRole('button');
        expect(backButton).toHaveAttribute('aria-label', 'Go back');
      });
    });
  });

  describe('Sync Status Indicators', () => {
    test('shows offline indicator when in fallback mode', () => {
      renderWithProvider(
        <Header title="Test Page" showBack onBack={() => {}} />,
        { fallbackMode: true }
      );

      expect(screen.getByText('Offline')).toBeInTheDocument();
    });

    test('shows error indicator when sync fails', () => {
      renderWithProvider(
        <Header title="Test Page" showBack onBack={() => {}} />,
        { syncStatus: 'error' }
      );

      expect(screen.getByText('Error')).toBeInTheDocument();
    });

    test('shows cached indicator when using fallback content', () => {
      renderWithProvider(
        <Header title="Test Page" showBack onBack={() => {}} />
      );

      expect(screen.getByText('Cached')).toBeInTheDocument();
    });
  });

  describe('Breadcrumb Dynamic Content', () => {
    test('uses dynamic aria-label for breadcrumbs', async () => {
      const breadcrumbItems: BreadcrumbItem[] = [
        { label: 'Home', href: '/' },
        { label: 'Page', href: '/page' }
      ];

      renderWithProvider(
        <Breadcrumb items={breadcrumbItems} testId="dynamic-breadcrumb" />
      );

      await waitFor(() => {
        const breadcrumb = screen.getByTestId('dynamic-breadcrumb');
        expect(breadcrumb).toHaveAttribute('aria-label', 'Breadcrumb');
      });
    });
  });
});