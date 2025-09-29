# Cart Management Components

A comprehensive set of React components for cart management in the Tchat commerce system. Built with TypeScript, React 18.3.1, and RTK Query integration.

## Features

- **Complete cart state management** with optimistic updates
- **Real-time validation** and error handling
- **Accessibility compliance** (WCAG 2.1 AA)
- **Responsive design** with mobile-first approach
- **TypeScript type safety** throughout
- **Performance optimized** with proper memoization
- **Seamless RTK Query integration**

## Components Overview

### Core Components

#### CartProvider
Central state management provider for cart operations.

```tsx
import { CartProvider } from '@/components/cart';

<CartProvider userId="user123" businessId="business456" autoValidate={true}>
  {/* Your cart components */}
</CartProvider>
```

**Props:**
- `userId?: string` - Current user ID
- `sessionId?: string` - Session ID for guest carts
- `businessId?: string` - Business context
- `autoValidate?: boolean` - Enable real-time validation (default: true)

#### CartSummary
Displays cart totals and checkout interface.

```tsx
import { CartSummary } from '@/components/cart';

<CartSummary
  showCheckout={true}
  showDetails={true}
  onCheckout={() => navigate('/checkout')}
/>
```

**Props:**
- `showCheckout?: boolean` - Show checkout button (default: true)
- `showDetails?: boolean` - Show detailed breakdown (default: true)
- `compact?: boolean` - Use compact layout
- `onCheckout?: () => void` - Checkout handler

#### CartItemList
Container for displaying cart items with bulk actions.

```tsx
import { CartItemList } from '@/components/cart';

<CartItemList
  showBulkActions={true}
  showSelection={true}
  editable={true}
  onBulkRemove={(ids) => console.log('Remove', ids)}
/>
```

**Props:**
- `showBulkActions?: boolean` - Enable bulk action controls
- `showSelection?: boolean` - Show item selection checkboxes
- `editable?: boolean` - Allow item editing (default: true)
- `compact?: boolean` - Use compact item layout

#### CartItem
Individual cart item display with controls.

```tsx
import { CartItem } from '@/components/cart';

<CartItem
  item={cartItem}
  editable={true}
  showGiftOptions={true}
  onQuantityChange={(qty) => console.log('New quantity:', qty)}
/>
```

**Props:**
- `item: CartItem` - Cart item data (required)
- `editable?: boolean` - Enable quantity controls (default: true)
- `showGiftOptions?: boolean` - Show gift message option (default: true)
- `compact?: boolean` - Use compact layout

### Product Integration

#### AddToCartButton
Comprehensive add-to-cart functionality with variants support.

```tsx
import { AddToCartButton } from '@/components/cart';

<AddToCartButton
  product={product}
  variants={productVariants}
  mode="detailed"
  showQuantity={true}
  onSuccess={(item) => console.log('Added:', item)}
/>
```

**Props:**
- `product: Product` - Product information (required)
- `variants?: ProductVariant[]` - Available product variants
- `mode?: 'button' | 'inline' | 'dropdown' | 'detailed'` - Display mode
- `showQuantity?: boolean` - Show quantity selector (default: true)
- `showVariants?: boolean` - Show variant selector (default: true)

**Display Modes:**
- `button`: Simple add to cart button
- `inline`: Button with inline quantity controls
- `dropdown`: Button that opens detailed options in popover
- `detailed`: Full detailed view with all options

### Cart Management

#### CartValidation
Real-time cart validation display.

```tsx
import { CartValidation } from '@/components/cart';

<CartValidation
  showDetails={true}
  showSummary={true}
  autoRefresh={true}
  onResolveIssue={(issue) => console.log('Resolve:', issue)}
/>
```

**Props:**
- `showDetails?: boolean` - Show detailed issues (default: true)
- `showSummary?: boolean` - Show validation summary (default: true)
- `autoRefresh?: boolean` - Enable auto-refresh
- `onResolveIssue?: (issue) => void` - Issue resolution handler

#### CouponInput
Discount code input and management.

```tsx
import { CouponInput } from '@/components/cart';

<CouponInput
  showSuggestions={true}
  collapsible={true}
  onCouponApplied={(code, discount) => console.log('Applied:', code, discount)}
/>
```

**Props:**
- `showSuggestions?: boolean` - Show coupon suggestions
- `collapsible?: boolean` - Enable collapsible behavior
- `onCouponApplied?: (code, discount) => void` - Applied callback
- `onCouponRemoved?: (code) => void` - Removed callback

## Usage Examples

### Basic Cart Implementation

```tsx
import React from 'react';
import {
  CartProvider,
  CartItemList,
  CartSummary,
  CartValidation,
  CouponInput
} from '@/components/cart';

function CartPage() {
  const userId = useAuth().user?.id;
  const sessionId = useSession().id;

  return (
    <CartProvider
      userId={userId}
      sessionId={sessionId}
      businessId="business123"
      autoValidate={true}
    >
      <div className="container mx-auto p-4">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main cart content */}
          <div className="lg:col-span-2 space-y-6">
            <CartItemList
              showBulkActions={true}
              showSelection={true}
              showGiftOptions={true}
            />

            <CartValidation
              showDetails={true}
              autoRefresh={true}
            />
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            <CouponInput
              showSuggestions={true}
              collapsible={true}
            />

            <CartSummary
              showDetails={true}
              onCheckout={() => navigate('/checkout')}
            />
          </div>
        </div>
      </div>
    </CartProvider>
  );
}
```

### Product Page Integration

```tsx
import React from 'react';
import { CartProvider, AddToCartButton } from '@/components/cart';

function ProductPage({ product, variants }) {
  return (
    <CartProvider businessId={product.businessId}>
      <div className="product-details">
        {/* Product info */}
        <div className="product-info">
          <h1>{product.name}</h1>
          <p>{product.description}</p>
        </div>

        {/* Add to cart */}
        <AddToCartButton
          product={product}
          variants={variants}
          mode="detailed"
          showQuantity={true}
          showVariants={true}
          showStock={true}
          showPrice={true}
          onSuccess={(item) => {
            toast.success(`Added ${item.quantity} items to cart`);
          }}
          onError={(error) => {
            toast.error(`Failed to add to cart: ${error}`);
          }}
        />
      </div>
    </CartProvider>
  );
}
```

### Compact Cart Widget

```tsx
import React from 'react';
import { CartProvider, CartSummary } from '@/components/cart';

function CartWidget() {
  return (
    <CartProvider userId={currentUser.id} businessId="business123">
      <CartSummary
        compact={true}
        showDetails={false}
        checkoutText="View Cart"
        onCheckout={() => navigate('/cart')}
      />
    </CartProvider>
  );
}
```

## State Management

The cart system uses RTK Query for state management with the following features:

### Optimistic Updates
All cart operations use optimistic updates for immediate UI feedback:

```tsx
// Adding items shows immediate feedback
await addToCart({ productId: '123', quantity: 2 });
// UI updates immediately, rolls back on error
```

### Real-time Validation
Cart validation runs automatically and can be configured:

```tsx
<CartProvider autoValidate={true}> // Enables real-time validation
  <CartValidation autoRefresh={true} /> // Shows validation status
</CartProvider>
```

### Caching Strategy
- Cart data cached with 30-second polling
- Validation cached with 60-second polling
- Automatic cache invalidation on mutations

## Accessibility Features

All components follow WCAG 2.1 AA standards:

### Keyboard Navigation
- Full keyboard support with logical tab order
- Enter/Space activation for interactive elements
- Escape key support for modals and dropdowns

### Screen Reader Support
- Semantic HTML structure
- Proper ARIA labels and descriptions
- Live regions for dynamic content updates
- Clear focus indicators

### Visual Accessibility
- High contrast color ratios
- Proper color usage (not color-only information)
- Scalable text and touch targets
- Reduced motion support

### Implementation Example

```tsx
<AddToCartButton
  product={product}
  aria-label={`Add ${product.name} to cart for ${formatCurrency(product.price)}`}
  onSuccess={(item) => {
    // Announce to screen readers
    announceToScreenReader(`Added ${item.quantity} ${product.name} to cart`);
  }}
/>
```

## Performance Optimization

### Memoization
Components use React.memo and useMemo for optimal performance:

```tsx
// Automatically memoized for performance
const CartItem = React.memo(({ item, ...props }) => {
  const computedData = useMemo(() => {
    return expensiveCalculation(item);
  }, [item]);

  return <div>{/* Component content */}</div>;
});
```

### Virtualization
Large item lists automatically virtualize:

```tsx
<CartItemList
  virtualizationThreshold={50} // Virtualizes if >50 items
  maxHeight="600px"
/>
```

### Lazy Loading
Images and non-critical content lazy load:

```tsx
<img
  src={item.productImage}
  loading="lazy" // Native lazy loading
  alt={item.productName}
/>
```

## Error Handling

Comprehensive error handling throughout:

### Network Errors
```tsx
<CartProvider>
  {/* Automatic retry with exponential backoff */}
  <CartItemList />
</CartProvider>
```

### Validation Errors
```tsx
<CartValidation
  onResolveIssue={(issue) => {
    // Handle specific validation issues
    if (issue.type === 'STOCK_INSUFFICIENT') {
      updateCartItem(issue.productId, { quantity: availableStock });
    }
  }}
/>
```

### User-Friendly Messages
```tsx
<AddToCartButton
  onError={(error) => {
    // Show user-friendly error messages
    toast.error(getUserFriendlyErrorMessage(error));
  }}
/>
```

## TypeScript Support

Full TypeScript support with comprehensive type definitions:

```tsx
import type {
  CartContextValue,
  CartItemProps,
  AddToCartButtonProps,
  ProductVariant
} from '@/components/cart';

// Type-safe component usage
const handleAddToCart = (item: AddToCartRequest) => {
  // TypeScript ensures type safety
};
```

## Styling and Theming

Components use Tailwind CSS with design system integration:

### Custom Styling
```tsx
<CartSummary
  className="border-2 border-primary rounded-lg shadow-lg"
  compact={true}
/>
```

### Theme Integration
```tsx
// Follows existing design system
<CartItem
  item={item}
  className="bg-card text-card-foreground" // Design system colors
/>
```

## Testing

Components are built with testing in mind:

### Test Utilities
```tsx
import { render, screen } from '@testing-library/react';
import { CartProvider, CartSummary } from '@/components/cart';

// Mock cart provider for testing
const renderWithCart = (component, cartData = mockCart) => {
  return render(
    <CartProvider value={cartData}>
      {component}
    </CartProvider>
  );
};
```

### Example Tests
```tsx
test('CartSummary displays correct totals', () => {
  renderWithCart(<CartSummary />);

  expect(screen.getByText('$29.99')).toBeInTheDocument();
  expect(screen.getByRole('button', { name: /checkout/i })).toBeEnabled();
});
```

## Browser Support

- **Modern browsers**: Chrome 88+, Firefox 85+, Safari 14+, Edge 88+
- **Mobile browsers**: iOS Safari 14+, Android Chrome 88+
- **Progressive enhancement**: Graceful degradation for older browsers

## Contributing

When adding new features:

1. Follow existing TypeScript patterns
2. Include comprehensive prop documentation
3. Add accessibility features
4. Include error handling
5. Write tests for new functionality
6. Update this documentation

## License

Part of the Tchat commerce system. See main project license.