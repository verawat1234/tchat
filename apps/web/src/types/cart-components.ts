/**
 * Cart Components Type Definitions
 *
 * Comprehensive TypeScript type definitions for the cart management system.
 * Extends base commerce types with component-specific interfaces.
 */

import type {
  Cart,
  CartItem,
  CartValidation,
  CartValidationIssue,
  Product,
  AddToCartRequest,
  UpdateCartItemRequest,
  ApplyCouponRequest,
} from './commerce';

// ===== Core Context Types =====

/**
 * Cart context value interface
 */
export interface CartContextValue {
  // State
  cart: Cart | null;
  cartValidation: CartValidation | null;
  isLoading: boolean;
  error: string | null;

  // Operations
  addToCart: (item: AddToCartRequest) => Promise<void>;
  updateCartItem: (itemId: string, updates: UpdateCartItemRequest) => Promise<void>;
  removeFromCart: (itemId: string) => Promise<void>;
  applyCoupon: (couponCode: string) => Promise<void>;
  removeCoupon: () => Promise<void>;

  // Utilities
  getCartItem: (productId: string, variantId?: string) => CartItem | undefined;
  getCartItemQuantity: (productId: string, variantId?: string) => number;
  isInCart: (productId: string, variantId?: string) => boolean;

  // Computed totals
  itemCount: number;
  subtotal: string;
  tax: string;
  shipping: string;
  discount: string;
  total: string;
}

/**
 * Cart provider configuration
 */
export interface CartProviderConfig {
  userId?: string;
  sessionId?: string;
  businessId?: string;
  autoValidate?: boolean;
  pollingInterval?: number;
  validationInterval?: number;
}

// ===== Component State Types =====

/**
 * Loading state configuration
 */
export interface LoadingState {
  isLoading: boolean;
  isUpdating: boolean;
  isAdding: boolean;
  isRemoving: boolean;
  isValidating: boolean;
  operation?: string;
}

/**
 * Error state configuration
 */
export interface ErrorState {
  hasError: boolean;
  errorMessage?: string;
  errorCode?: string;
  errorType?: 'network' | 'validation' | 'server' | 'client';
  retry?: () => void;
}

/**
 * Success state configuration
 */
export interface SuccessState {
  hasSuccess: boolean;
  successMessage?: string;
  duration?: number;
  autoHide?: boolean;
}

// ===== Display Types =====

/**
 * Component size variants
 */
export type ComponentSize = 'sm' | 'default' | 'lg';

/**
 * Component display modes
 */
export type DisplayMode = 'button' | 'inline' | 'dropdown' | 'detailed' | 'compact';

/**
 * Component variants
 */
export type ComponentVariant = 'default' | 'outline' | 'secondary' | 'ghost' | 'destructive';

/**
 * Visual states
 */
export type VisualState = 'default' | 'loading' | 'success' | 'error' | 'warning' | 'disabled';

// ===== Product Integration Types =====

/**
 * Product variant for cart operations
 */
export interface ProductVariant {
  id: string;
  name: string;
  price: string;
  compareAtPrice?: string;
  sku: string;
  stock: number;
  maxQuantity: number;
  isAvailable: boolean;
  attributes: Record<string, string>;
}

/**
 * Add to cart configuration
 */
export interface AddToCartConfig {
  productId: string;
  variantId?: string;
  quantity: number;
  maxQuantity?: number;
  isGift?: boolean;
  giftMessage?: string;
  customAttributes?: Record<string, any>;
}

/**
 * Cart item display configuration
 */
export interface CartItemDisplayConfig {
  showImage: boolean;
  showVariant: boolean;
  showSku: boolean;
  showGiftOptions: boolean;
  showWishlistAction: boolean;
  showQuantityControls: boolean;
  showRemoveAction: boolean;
  showStockInfo: boolean;
  showPricing: boolean;
}

// ===== Validation Types =====

/**
 * Cart validation configuration
 */
export interface CartValidationConfig {
  enableRealTimeValidation: boolean;
  validationInterval: number;
  autoRefresh: boolean;
  showDetails: boolean;
  showSummary: boolean;
  groupBySeverity: boolean;
}

/**
 * Validation issue display
 */
export interface ValidationIssueDisplay {
  issue: CartValidationIssue;
  icon: React.ComponentType<any>;
  variant: ComponentVariant;
  recommendation: string;
  canResolve: boolean;
  resolveAction?: () => void;
}

/**
 * Validation status
 */
export interface ValidationStatus {
  type: 'valid' | 'warning' | 'error';
  message: string;
  icon: React.ComponentType<any>;
  count: number;
  details?: CartValidationIssue[];
}

// ===== Coupon Types =====

/**
 * Coupon suggestion
 */
export interface CouponSuggestion {
  code: string;
  description: string;
  type: 'percentage' | 'fixed' | 'shipping';
  requirements?: string;
  validUntil?: string;
  isPopular?: boolean;
}

/**
 * Applied coupon information
 */
export interface AppliedCouponInfo {
  code: string;
  discount: string;
  discountType: 'percentage' | 'fixed';
  discountAmount: number;
  discountPercentage?: number;
  validUntil?: string;
  description?: string;
}

/**
 * Coupon validation result
 */
export interface CouponValidationResult {
  isValid: boolean;
  errorMessage?: string;
  requirements?: string[];
  discount?: {
    amount: string;
    type: 'percentage' | 'fixed';
  };
}

// ===== Event Handling Types =====

/**
 * Cart event types
 */
export type CartEventType =
  | 'item_added'
  | 'item_updated'
  | 'item_removed'
  | 'coupon_applied'
  | 'coupon_removed'
  | 'cart_validated'
  | 'checkout_started'
  | 'error_occurred';

/**
 * Cart event data
 */
export interface CartEventData {
  type: CartEventType;
  timestamp: string;
  data: any;
  context?: string;
}

/**
 * Event handler configuration
 */
export interface EventHandlerConfig {
  onItemAdded?: (item: CartItem) => void;
  onItemUpdated?: (itemId: string, updates: UpdateCartItemRequest) => void;
  onItemRemoved?: (itemId: string) => void;
  onCouponApplied?: (couponCode: string, discount: string) => void;
  onCouponRemoved?: (couponCode: string) => void;
  onCartValidated?: (validation: CartValidation) => void;
  onCheckoutStarted?: () => void;
  onError?: (error: string, context?: string) => void;
}

// ===== Accessibility Types =====

/**
 * Accessibility configuration
 */
export interface AccessibilityConfig {
  enableKeyboardNavigation: boolean;
  enableScreenReaderSupport: boolean;
  enableHighContrastMode: boolean;
  enableReducedMotion: boolean;
  announceChanges: boolean;
  customAriaLabels?: Record<string, string>;
}

/**
 * ARIA attributes for cart components
 */
export interface CartAriaAttributes {
  'aria-label'?: string;
  'aria-labelledby'?: string;
  'aria-describedby'?: string;
  'aria-expanded'?: boolean;
  'aria-selected'?: boolean;
  'aria-disabled'?: boolean;
  'aria-invalid'?: boolean;
  'aria-live'?: 'polite' | 'assertive' | 'off';
  'aria-atomic'?: boolean;
  'aria-busy'?: boolean;
  role?: string;
}

// ===== Performance Types =====

/**
 * Performance configuration
 */
export interface PerformanceConfig {
  enableOptimisticUpdates: boolean;
  enableCaching: boolean;
  cacheTimeout: number;
  enableVirtualization: boolean;
  virtualizationThreshold: number;
  enableLazyLoading: boolean;
  enablePrefetching: boolean;
}

/**
 * Optimization settings
 */
export interface OptimizationSettings {
  debounceDelay: number;
  throttleDelay: number;
  maxRetries: number;
  retryDelay: number;
  batchSize: number;
  enableBatching: boolean;
}

// ===== Layout Types =====

/**
 * Layout configuration
 */
export interface LayoutConfig {
  direction: 'row' | 'column';
  alignment: 'start' | 'center' | 'end' | 'stretch';
  justification: 'start' | 'center' | 'end' | 'between' | 'around' | 'evenly';
  spacing: 'tight' | 'normal' | 'loose';
  responsive: boolean;
  breakpoints?: Record<string, any>;
}

/**
 * Grid configuration
 */
export interface GridConfig {
  columns: number | 'auto';
  rows?: number | 'auto';
  gap: string | number;
  responsive: boolean;
  minItemWidth?: string;
  maxItemWidth?: string;
}

// ===== Animation Types =====

/**
 * Animation configuration
 */
export interface AnimationConfig {
  enableAnimations: boolean;
  duration: number;
  easing: string;
  respectReducedMotion: boolean;
  staggerDelay?: number;
}

/**
 * Transition types
 */
export type TransitionType = 'fade' | 'slide' | 'scale' | 'bounce' | 'none';

// ===== Theme Types =====

/**
 * Cart theme configuration
 */
export interface CartTheme {
  colors: {
    primary: string;
    secondary: string;
    success: string;
    warning: string;
    error: string;
    muted: string;
    background: string;
    foreground: string;
  };
  spacing: {
    xs: string;
    sm: string;
    md: string;
    lg: string;
    xl: string;
  };
  typography: {
    fontFamily: string;
    fontSize: Record<ComponentSize, string>;
    fontWeight: Record<string, number>;
    lineHeight: Record<ComponentSize, string>;
  };
  borderRadius: {
    sm: string;
    md: string;
    lg: string;
  };
  shadows: {
    sm: string;
    md: string;
    lg: string;
  };
}

// ===== Utility Types =====

/**
 * Component ref types
 */
export type CartComponentRef<T = HTMLElement> = React.RefObject<T>;

/**
 * Component class names
 */
export type ComponentClassName = string | undefined;

/**
 * Component children
 */
export type ComponentChildren = React.ReactNode;

/**
 * Component event handlers
 */
export type ComponentEventHandler<T = Event> = (event: T) => void;

/**
 * Optional component props
 */
export type OptionalProps<T> = Partial<T>;

/**
 * Required component props
 */
export type RequiredProps<T, K extends keyof T> = T & Required<Pick<T, K>>;

// ===== Export all types =====

export type {
  // Re-export commerce types
  Cart,
  CartItem,
  CartValidation,
  CartValidationIssue,
  Product,
  AddToCartRequest,
  UpdateCartItemRequest,
  ApplyCouponRequest,
};