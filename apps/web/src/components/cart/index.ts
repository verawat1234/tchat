/**
 * Cart Components - Comprehensive Cart Management System
 *
 * A complete set of reusable cart management components for the commerce system.
 * Built with TypeScript, React 18.3.1, and RTK Query integration.
 *
 * Features:
 * - Complete cart state management
 * - Optimistic updates with rollback
 * - Real-time validation
 * - Accessibility compliance (WCAG 2.1 AA)
 * - Responsive design
 * - Error handling and recovery
 * - Loading states and animations
 * - TypeScript type safety
 *
 * Usage:
 * ```tsx
 * import { CartProvider, CartSummary, CartItemList, AddToCartButton } from '@/components/cart';
 *
 * function App() {
 *   return (
 *     <CartProvider userId="user123" businessId="business456">
 *       <CartItemList />
 *       <CartSummary onCheckout={() => navigate('/checkout')} />
 *     </CartProvider>
 *   );
 * }
 * ```
 */

// ===== Core Components =====

export { CartProvider, useCart } from './CartProvider';
export type { CartProviderProps, CartContextValue } from './CartProvider';

export { CartSummary } from './CartSummary';
export type { CartSummaryProps } from './CartSummary';

export { CartItem } from './CartItem';
export type { CartItemProps } from './CartItem';

export { CartItemList } from './CartItemList';
export type { CartItemListProps } from './CartItemList';

export { AddToCartButton } from './AddToCartButton';
export type { AddToCartButtonProps, ProductVariant } from './AddToCartButton';

export { CartValidation } from './CartValidation';
export type { CartValidationProps } from './CartValidation';

export { CouponInput } from './CouponInput';
export type { CouponInputProps } from './CouponInput';

// ===== Default Exports =====

export { default as CartSummaryDefault } from './CartSummary';
export { default as CartItemDefault } from './CartItem';
export { default as CartItemListDefault } from './CartItemList';
export { default as AddToCartButtonDefault } from './AddToCartButton';
export { default as CartValidationDefault } from './CartValidation';
export { default as CouponInputDefault } from './CouponInput';

// ===== Component Groups =====

/**
 * Core cart components for basic functionality
 */
export const CoreCartComponents = {
  Provider: CartProvider,
  Summary: CartSummary,
  ItemList: CartItemList,
  Item: CartItem,
} as const;

/**
 * Product integration components
 */
export const ProductComponents = {
  AddToCartButton,
} as const;

/**
 * Cart management components
 */
export const CartManagementComponents = {
  Validation: CartValidation,
  CouponInput,
} as const;

/**
 * All cart components
 */
export const CartComponents = {
  ...CoreCartComponents,
  ...ProductComponents,
  ...CartManagementComponents,
} as const;

// ===== Utility Types =====

/**
 * Union type of all cart component names
 */
export type CartComponentName = keyof typeof CartComponents;

/**
 * Cart system configuration
 */
export interface CartSystemConfig {
  /** Business context */
  businessId?: string;
  /** User context */
  userId?: string;
  /** Session context for guest users */
  sessionId?: string;
  /** Auto-validation settings */
  autoValidate?: boolean;
  /** Polling intervals */
  polling?: {
    cart?: number;
    validation?: number;
  };
  /** Feature flags */
  features?: {
    wishlist?: boolean;
    giftMessages?: boolean;
    bulkActions?: boolean;
    suggestions?: boolean;
  };
}

/**
 * Cart event handlers
 */
export interface CartEventHandlers {
  /** Cart state change events */
  onCartChange?: (cart: any) => void;
  onItemAdded?: (item: any) => void;
  onItemRemoved?: (itemId: string) => void;
  onItemUpdated?: (itemId: string, updates: any) => void;

  /** Coupon events */
  onCouponApplied?: (couponCode: string, discount: string) => void;
  onCouponRemoved?: (couponCode: string) => void;

  /** Checkout events */
  onCheckoutStart?: () => void;
  onCheckoutComplete?: (orderId: string) => void;

  /** Error events */
  onError?: (error: string, context?: string) => void;
}

// ===== Helper Functions =====

/**
 * Cart utility functions
 */
export const CartUtils = {
  /**
   * Format currency amount
   */
  formatCurrency: (amount: string, currency = 'USD'): string => {
    const numAmount = parseFloat(amount);
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
      minimumFractionDigits: 2,
    }).format(numAmount);
  },

  /**
   * Calculate cart totals
   */
  calculateTotals: (cart: any) => {
    if (!cart) return { itemCount: 0, total: '0.00' };

    return {
      itemCount: cart.itemCount || 0,
      subtotal: cart.subtotalAmount || '0.00',
      tax: cart.taxAmount || '0.00',
      shipping: cart.shippingAmount || '0.00',
      discount: cart.discountAmount || '0.00',
      total: cart.totalAmount || '0.00',
    };
  },

  /**
   * Check if cart is valid for checkout
   */
  isValidForCheckout: (cart: any, validation: any): boolean => {
    if (!cart || cart.itemCount === 0) return false;
    if (!validation) return true;

    const hasErrors = validation.issues?.some((issue: any) => issue.severity === 'error');
    return !hasErrors;
  },

  /**
   * Get cart validation status
   */
  getValidationStatus: (validation: any): 'valid' | 'warning' | 'error' | null => {
    if (!validation) return null;

    const { issues } = validation;
    const hasErrors = issues.some((issue: any) => issue.severity === 'error');
    const hasWarnings = issues.some((issue: any) => issue.severity === 'warning');

    if (hasErrors) return 'error';
    if (hasWarnings) return 'warning';
    return 'valid';
  },
} as const;

// ===== Constants =====

/**
 * Default cart configuration
 */
export const DEFAULT_CART_CONFIG: CartSystemConfig = {
  autoValidate: true,
  polling: {
    cart: 30000, // 30 seconds
    validation: 60000, // 60 seconds
  },
  features: {
    wishlist: true,
    giftMessages: true,
    bulkActions: true,
    suggestions: true,
  },
} as const;

/**
 * Cart component sizes
 */
export const CART_SIZES = {
  sm: 'small',
  default: 'default',
  lg: 'large',
} as const;

/**
 * Cart display modes
 */
export const CART_MODES = {
  button: 'button',
  inline: 'inline',
  dropdown: 'dropdown',
  detailed: 'detailed',
} as const;

// ===== Version Info =====

/**
 * Cart components version information
 */
export const CART_COMPONENTS_VERSION = {
  major: 1,
  minor: 0,
  patch: 0,
  version: '1.0.0',
  build: new Date().toISOString(),
} as const;