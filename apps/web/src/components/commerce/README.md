# Category Browsing System

A comprehensive React component library for e-commerce category navigation and product discovery. Built with TypeScript, RTK Query, and accessibility in mind.

## ✨ Features

### 🧭 Navigation
- **Hierarchical Category Tree** - Unlimited depth category navigation with expand/collapse
- **Breadcrumb Navigation** - Smart breadcrumbs with mobile responsive collapsing
- **Quick Navigation** - Jump between sibling categories and parent categories

### 🔍 Search & Discovery
- **Real-time Search** - Debounced search with autocomplete suggestions
- **Voice Search** - Built-in voice search support (when available)
- **Recent Searches** - Persistent search history with localStorage
- **Smart Suggestions** - Category and product suggestions with analytics

### 🎛️ Filtering & Sorting
- **Advanced Filters** - Price range, ratings, brands, colors, sizes, and custom attributes
- **Filter Presets** - Quick filter combinations for common use cases
- **Real-time Application** - Filters apply immediately with optimistic updates
- **Mobile Responsive** - Collapsible filter panels for mobile devices

### 📱 Responsive Design
- **Mobile First** - Optimized for mobile devices with touch-friendly interfaces
- **Responsive Layouts** - Adapts to screen sizes from 320px to 4K displays
- **Progressive Enhancement** - Works without JavaScript, enhanced with it

### ♿ Accessibility
- **WCAG 2.1 AA Compliant** - Full keyboard navigation and screen reader support
- **Focus Management** - Proper focus handling for dynamic content
- **ARIA Attributes** - Comprehensive ARIA labels and descriptions
- **High Contrast** - Supports high contrast modes and custom themes

### ⚡ Performance
- **Lazy Loading** - Components and data load on demand
- **Virtual Scrolling** - Handle thousands of items efficiently
- **Optimistic Updates** - Immediate UI feedback with rollback on errors
- **Intelligent Caching** - RTK Query powered caching with tag invalidation

### 📊 Analytics
- **View Tracking** - Automatic category and product view tracking
- **Search Analytics** - Track search queries and result interactions
- **Performance Monitoring** - Built-in performance metrics collection

## 🚀 Quick Start

### Basic Setup

```tsx
import { CategoryProvider, CategoryTree, CategoryGrid, CategorySearch } from '@/components/commerce';

function CategoryBrowser() {
  return (
    <CategoryProvider businessId="your-business-id">
      <div className="flex gap-6">
        <aside className="w-80">
          <CategorySearch className="mb-4" />
          <CategoryTree />
        </aside>
        <main className="flex-1">
          <CategoryGrid />
        </main>
      </div>
    </CategoryProvider>
  );
}
```

### Full-Featured Layout

```tsx
import {
  CategoryProvider,
  CategoryBreadcrumbs,
  CategorySearch,
  CategoryFilters,
  CategoryTree,
  CategoryGrid,
  CategoryProductList
} from '@/components/commerce';

function FullCategoryBrowser() {
  const [selectedCategoryId, setSelectedCategoryId] = useState<string | null>(null);

  return (
    <CategoryProvider businessId="your-business-id">
      <div className="min-h-screen bg-background">
        {/* Header */}
        <header className="border-b bg-card px-6 py-4">
          <CategorySearch
            placeholder="Search categories and products..."
            showVoiceSearch
            showRecentSearches
          />
        </header>

        {/* Main Content */}
        <div className="flex">
          {/* Sidebar */}
          <aside className="w-80 border-r bg-card p-6 space-y-6">
            <CategoryFilters
              variant="sidebar"
              showPresets
              showSaveFilter
            />
            <CategoryTree
              showProductCounts
              showIcons
              onCategorySelect={(category) => setSelectedCategoryId(category.id)}
            />
          </aside>

          {/* Content Area */}
          <main className="flex-1 p-6 space-y-6">
            <CategoryBreadcrumbs />

            {selectedCategoryId ? (
              <CategoryProductList
                categoryId={selectedCategoryId}
                enableLazyLoading
                enableInfiniteScroll
                gridColumns={4}
              />
            ) : (
              <CategoryGrid
                viewMode="grid"
                gridColumns={4}
                showFilters={false}
              />
            )}
          </main>
        </div>
      </div>
    </CategoryProvider>
  );
}
```

### Mobile Optimized

```tsx
import {
  CategoryProvider,
  CategorySearch,
  CategoryFilters,
  CategoryGrid,
  CompactCategoryTree
} from '@/components/commerce';

function MobileCategoryBrowser() {
  return (
    <CategoryProvider>
      <div className="min-h-screen bg-background p-4 space-y-4">
        <CategorySearch variant="compact" />

        <div className="flex gap-2">
          <CategoryFilters variant="modal" />
          <div className="flex-1" />
        </div>

        <CompactCategoryTree
          onCategorySelect={(category) => console.log('Selected:', category)}
        />

        <CategoryGrid
          viewMode="list"
          gridColumns={1}
          showFilters={false}
          showViewModeToggle={false}
        />
      </div>
    </CategoryProvider>
  );
}
```

## 📖 Component API

### CategoryProvider

Central state management for the category browsing system.

```tsx
interface CategoryProviderProps {
  children: React.ReactNode;
  businessId?: string;  // Optional business scope
}
```

### CategoryTree

Hierarchical category navigation with expand/collapse functionality.

```tsx
interface CategoryTreeProps {
  className?: string;
  maxHeight?: string;                    // Default: '400px'
  showProductCounts?: boolean;           // Default: true
  showIcons?: boolean;                   // Default: true
  allowMultiSelect?: boolean;            // Default: false
  onCategorySelect?: (category: Category) => void;
  onCategoryToggle?: (categoryId: string, isExpanded: boolean) => void;
}
```

### CategoryGrid

Responsive grid layout for displaying categories with filtering and sorting.

```tsx
interface CategoryGridProps {
  className?: string;
  categories?: Category[];               // Optional override
  viewMode?: 'grid' | 'list' | 'compact';
  gridColumns?: 1 | 2 | 3 | 4 | 5 | 6;  // Default: 3
  showFilters?: boolean;                 // Default: true
  showSorting?: boolean;                 // Default: true
  showPagination?: boolean;              // Default: true
  showViewModeToggle?: boolean;          // Default: true
  infiniteScroll?: boolean;              // Default: false
  virtualScroll?: boolean;               // Default: false
  pageSize?: number;                     // Default: 24
  emptyMessage?: string;
  emptyAction?: React.ReactNode;
  onCategorySelect?: (category: Category) => void;
}
```

### CategorySearch

Advanced search with autocomplete, filters, and voice search.

```tsx
interface CategorySearchProps {
  className?: string;
  placeholder?: string;                  // Default: 'Search categories and products...'
  variant?: 'default' | 'compact' | 'expanded';
  showFilters?: boolean;                 // Default: true
  showVoiceSearch?: boolean;             // Default: true
  showRecentSearches?: boolean;          // Default: true
  autoFocus?: boolean;                   // Default: false
  onSearchSubmit?: (query: string) => void;
  onCategorySelect?: (category: Category) => void;
  onProductSelect?: (product: Product) => void;
}
```

### CategoryFilters

Comprehensive filtering interface with multiple filter types.

```tsx
interface CategoryFiltersProps {
  className?: string;
  variant?: 'sidebar' | 'inline' | 'modal';  // Default: 'sidebar'
  showPresets?: boolean;                      // Default: true
  showSaveFilter?: boolean;                   // Default: false
  onFiltersChange?: (filters: FilterState) => void;
}
```

### CategoryBreadcrumbs

Navigation breadcrumbs with responsive collapsing and quick navigation.

```tsx
interface CategoryBreadcrumbsProps {
  className?: string;
  showHome?: boolean;                    // Default: true
  showIcons?: boolean;                   // Default: true
  maxItems?: number;                     // Default: 5
  separator?: 'chevron' | 'slash' | 'arrow' | 'dot';  // Default: 'chevron'
  size?: 'sm' | 'md' | 'lg';            // Default: 'md'
  variant?: 'default' | 'ghost' | 'outline';  // Default: 'ghost'
  onNavigate?: (category: Category | null) => void;
}
```

### CategoryProductList

Product listing with lazy loading, infinite scroll, and multiple view modes.

```tsx
interface CategoryProductListProps {
  categoryId?: string;
  className?: string;
  viewMode?: 'grid' | 'list' | 'compact';
  gridColumns?: 2 | 3 | 4 | 5 | 6;      // Default: 4
  showFilters?: boolean;                 // Default: true
  showSorting?: boolean;                 // Default: true
  enableLazyLoading?: boolean;           // Default: true
  enableVirtualScroll?: boolean;         // Default: false
  enableInfiniteScroll?: boolean;        // Default: true
  pageSize?: number;                     // Default: 24
  onProductSelect?: (product: Product) => void;
  onAddToCart?: (product: Product) => void;
}
```

## 🎨 Customization

### Themes

The components use Tailwind CSS and support custom themes through CSS variables:

```css
:root {
  --category-primary: #3b82f6;
  --category-secondary: #f1f5f9;
  --category-accent: #10b981;
  --category-muted: #64748b;
}
```

### Custom Styling

All components accept `className` props and use `cn()` utility for class merging:

```tsx
<CategoryTree className="custom-tree border-2 border-red-500" />
<CategoryGrid className="my-custom-grid bg-slate-100" />
```

### Component Variants

Many components support multiple variants for different use cases:

```tsx
{/* Search variants */}
<CategorySearch variant="compact" />    {/* Minimal search bar */}
<CategorySearch variant="expanded" />   {/* Full-featured modal */}

{/* Filter variants */}
<CategoryFilters variant="sidebar" />   {/* Full sidebar */}
<CategoryFilters variant="modal" />     {/* Mobile sheet */}
<CategoryFilters variant="inline" />    {/* Inline card */}
```

## 🔌 Integration

### RTK Query

Components automatically integrate with the existing RTK Query commerce API:

```tsx
// These hooks are used internally
import {
  useGetCategoriesQuery,
  useGetRootCategoriesQuery,
  useGetProductsQuery,
  useTrackCategoryViewMutation,
} from '@/services/commerceApi';
```

### Analytics

Category views are automatically tracked when users navigate:

```tsx
// This happens automatically when categories are selected
const trackCategoryView = useTrackCategoryViewMutation();

trackCategoryView({
  categoryId: 'category-123',
  sessionId: 'session-456',
  ipAddress: '0.0.0.0',
  userAgent: navigator.userAgent,
  referrer: document.referrer,
});
```

### State Management

Access the category state from any child component:

```tsx
import { useCategory } from '@/components/commerce';

function MyComponent() {
  const {
    state,                    // Current category state
    setCurrentCategory,       // Navigate to category
    setSearchQuery,          // Update search
    updateFilters,           // Apply filters
    trackCategoryView,       // Track analytics
    categoriesQuery,         // RTK Query result
    hasActiveFilters,        // Computed state
  } = useCategory();

  return (
    <div>
      Current category: {state.currentCategoryId}
      Search query: {state.search.query}
      Active filters: {state.activeFilters.length}
    </div>
  );
}
```

## 🧪 Testing

### Accessibility Testing

Use these tools to verify accessibility compliance:

```bash
# Install accessibility testing tools
npm install --save-dev @axe-core/react jest-axe

# Run accessibility tests
npm run test:a11y
```

### Performance Testing

Monitor performance with built-in metrics:

```tsx
import { CategoryBrowsingDemo } from '@/components/commerce/CategoryBrowsingDemo';

// The demo component includes performance monitoring
<CategoryBrowsingDemo />
```

### Unit Testing

Test individual components:

```tsx
import { render, screen } from '@testing-library/react';
import { CategoryProvider, CategoryTree } from '@/components/commerce';

test('CategoryTree renders categories', () => {
  render(
    <CategoryProvider>
      <CategoryTree />
    </CategoryProvider>
  );

  expect(screen.getByRole('tree')).toBeInTheDocument();
});
```

## 📚 Examples

### E-commerce Store

```tsx
function ProductCatalog() {
  return (
    <CategoryProvider businessId="store-123">
      <div className="min-h-screen">
        <header className="bg-white shadow">
          <div className="container mx-auto px-4 py-6">
            <CategorySearch
              placeholder="Search products..."
              showVoiceSearch
            />
          </div>
        </header>

        <div className="container mx-auto px-4 py-8">
          <div className="grid grid-cols-12 gap-8">
            <aside className="col-span-3">
              <CategoryFilters showPresets />
              <CategoryTree className="mt-6" />
            </aside>

            <main className="col-span-9">
              <CategoryBreadcrumbs className="mb-6" />
              <CategoryGrid
                gridColumns={3}
                enableInfiniteScroll
              />
            </main>
          </div>
        </div>
      </div>
    </CategoryProvider>
  );
}
```

### Marketplace

```tsx
function Marketplace() {
  const [selectedBusiness, setSelectedBusiness] = useState('');

  return (
    <CategoryProvider businessId={selectedBusiness}>
      <div className="flex h-screen">
        <aside className="w-80 bg-card border-r p-6">
          <Select value={selectedBusiness} onValueChange={setSelectedBusiness}>
            <SelectTrigger className="mb-4">
              <SelectValue placeholder="Select store" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="store-1">Electronics Store</SelectItem>
              <SelectItem value="store-2">Fashion Boutique</SelectItem>
              <SelectItem value="store-3">Home & Garden</SelectItem>
            </SelectContent>
          </Select>

          <CategoryTree showProductCounts />
        </aside>

        <main className="flex-1 p-6">
          <CategorySearch className="mb-6" />
          <CategoryGrid viewMode="grid" gridColumns={4} />
        </main>
      </div>
    </CategoryProvider>
  );
}
```

### Content Management

```tsx
function CategoryManager() {
  const [editMode, setEditMode] = useState(false);

  return (
    <CategoryProvider>
      <div className="p-6">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">Category Management</h1>
          <Button onClick={() => setEditMode(!editMode)}>
            {editMode ? 'View Mode' : 'Edit Mode'}
          </Button>
        </div>

        <div className="grid grid-cols-12 gap-6">
          <div className="col-span-4">
            <CategoryTree
              allowMultiSelect={editMode}
              onCategorySelect={(category) => {
                if (editMode) {
                  console.log('Edit category:', category);
                }
              }}
            />
          </div>

          <div className="col-span-8">
            <CategoryGrid
              showFilters={false}
              onCategorySelect={(category) => {
                if (editMode) {
                  console.log('Edit category:', category);
                }
              }}
            />
          </div>
        </div>
      </div>
    </CategoryProvider>
  );
}
```

## 🛠️ Development

### Project Structure

```
components/commerce/
├── CategoryProvider.tsx        # State management
├── CategoryTree.tsx           # Tree navigation
├── CategoryCard.tsx           # Individual category display
├── CategoryGrid.tsx           # Grid layout
├── CategoryBreadcrumbs.tsx    # Navigation breadcrumbs
├── CategoryProductList.tsx    # Product listing
├── CategorySearch.tsx         # Search functionality
├── CategoryFilters.tsx        # Filter interface
├── CategoryBrowsingDemo.tsx   # Demo component
├── index.ts                   # Exports
└── README.md                  # Documentation
```

### Build Commands

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test

# Build for production
npm run build

# Check types
npm run type-check

# Lint code
npm run lint

# Format code
npm run format
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Ensure all tests pass: `npm test`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Support

- 📖 [Documentation](./README.md)
- 🐛 [Issue Tracker](https://github.com/your-repo/issues)
- 💬 [Discussions](https://github.com/your-repo/discussions)
- 📧 [Email Support](mailto:support@yourcompany.com)

## 🙏 Acknowledgments

- Built with [React](https://reactjs.org/) and [TypeScript](https://www.typescriptlang.org/)
- Styled with [Tailwind CSS](https://tailwindcss.com/) and [Radix UI](https://www.radix-ui.com/)
- State management with [RTK Query](https://redux-toolkit.js.org/rtk-query/overview)
- Icons from [Lucide React](https://lucide.dev/)
- Accessibility guidance from [WAI-ARIA](https://www.w3.org/WAI/ARIA/)

---

**Built with ❤️ for modern e-commerce experiences**