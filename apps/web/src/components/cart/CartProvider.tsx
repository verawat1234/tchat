/**
 * CartProvider - Cart Context Provider Component
 *
 * Provides cart state management and operations across the application.
 * Integrates with RTK Query commerce API for cart operations.
 *
 * Features:
 * - Centralized cart state management
 * - Optimistic updates for better UX
 * - Session and user cart management
 * - Cart validation and error handling
 * - Automatic cart merging on login
 */

import React, { createContext, useContext, useCallback, useEffect, useMemo } from 'react';
import {
  useGetCartQuery,
  useAddToCartMutation,
  useUpdateCartItemMutation,
  useRemoveFromCartMutation,
  useApplyCouponMutation,
  useRemoveCouponMutation,
  useValidateCartQuery,
  useMergeCartsMutation,
} from '../../services/commerceApi';
import type {
  Cart,
  CartItem,
  CartValidation,
  AddToCartRequest,
  UpdateCartItemRequest,
  ApplyCouponRequest,
} from '../../types/commerce';

// ===== Types =====

export interface CartContextValue {
  // Cart state
  cart: Cart | null;
  cartValidation: CartValidation | null;
  isLoading: boolean;
  error: string | null;

  // Cart operations
  addToCart: (item: AddToCartRequest) => Promise<void>;
  updateCartItem: (itemId: string, updates: UpdateCartItemRequest) => Promise<void>;
  removeFromCart: (itemId: string) => Promise<void>;
  applyCoupon: (couponCode: string) => Promise<void>;
  removeCoupon: () => Promise<void>;

  // Cart utilities
  getCartItem: (productId: string, variantId?: string) => CartItem | undefined;
  getCartItemQuantity: (productId: string, variantId?: string) => number;
  isInCart: (productId: string, variantId?: string) => boolean;

  // Cart totals (computed)
  itemCount: number;
  subtotal: string;
  tax: string;
  shipping: string;
  discount: string;
  total: string;
}

export interface CartProviderProps {
  children: React.ReactNode;
  userId?: string;
  sessionId?: string;
  businessId?: string;
  autoValidate?: boolean;
}

// ===== Context =====

const CartContext = createContext<CartContextValue | null>(null);

/**
 * Hook to access cart context
 * @throws Error if used outside CartProvider
 */
export const useCart = (): CartContextValue => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};

// ===== Provider Component =====

/**
 * CartProvider - Provides cart state and operations
 *
 * @param children - Child components
 * @param userId - Current user ID (optional for guest carts)
 * @param sessionId - Session ID for guest carts
 * @param businessId - Business context for cart
 * @param autoValidate - Whether to automatically validate cart (default: true)
 */
export const CartProvider: React.FC<CartProviderProps> = ({
  children,
  userId,
  sessionId,
  businessId,
  autoValidate = true,
}) => {
  // ===== RTK Query Hooks =====

  const {
    data: cart,
    isLoading: isCartLoading,
    error: cartError,
    refetch: refetchCart,
  } = useGetCartQuery(
    { userId, sessionId },
    {
      skip: !userId && !sessionId,
      pollingInterval: 30000, // Poll every 30s for cart changes
    }
  );

  const {
    data: cartValidation,
    isLoading: isValidationLoading,
    error: validationError,
  } = useValidateCartQuery(
    cart?.id || '',
    {
      skip: !cart?.id || !autoValidate,
      pollingInterval: 60000, // Poll every 60s for validation
    }
  );

  const [addToCartMutation, { isLoading: isAddingToCart }] = useAddToCartMutation();
  const [updateCartItemMutation, { isLoading: isUpdatingItem }] = useUpdateCartItemMutation();
  const [removeFromCartMutation, { isLoading: isRemovingItem }] = useRemoveFromCartMutation();
  const [applyCouponMutation, { isLoading: isApplyingCoupon }] = useApplyCouponMutation();
  const [removeCouponMutation, { isLoading: isRemovingCoupon }] = useRemoveCouponMutation();
  const [mergeCartsMutation] = useMergeCartsMutation();

  // ===== State =====

  const isLoading = isCartLoading || isValidationLoading || isAddingToCart ||
                   isUpdatingItem || isRemovingItem || isApplyingCoupon || isRemovingCoupon;

  const error = cartError || validationError;

  // ===== Cart Operations =====

  /**
   * Add item to cart with optimistic update
   */
  const addToCart = useCallback(async (item: AddToCartRequest) => {
    try {
      await addToCartMutation({
        cartId: cart?.id,
        item: {
          ...item,
          // Ensure business context if available
          ...(businessId && { businessId }),
        },
      }).unwrap();
    } catch (error) {
      console.error('Failed to add item to cart:', error);
      throw error;
    }
  }, [addToCartMutation, cart?.id, businessId]);

  /**
   * Update cart item quantity or properties
   */
  const updateCartItem = useCallback(async (itemId: string, updates: UpdateCartItemRequest) => {
    if (!cart?.id) {
      throw new Error('No active cart found');
    }

    try {
      await updateCartItemMutation({
        cartId: cart.id,
        itemId,
        updates,
      }).unwrap();
    } catch (error) {
      console.error('Failed to update cart item:', error);
      throw error;
    }
  }, [updateCartItemMutation, cart?.id]);

  /**
   * Remove item from cart
   */
  const removeFromCart = useCallback(async (itemId: string) => {
    if (!cart?.id) {
      throw new Error('No active cart found');
    }

    try {
      await removeFromCartMutation({
        cartId: cart.id,
        itemId,
      }).unwrap();
    } catch (error) {
      console.error('Failed to remove item from cart:', error);
      throw error;
    }
  }, [removeFromCartMutation, cart?.id]);

  /**
   * Apply coupon code to cart
   */
  const applyCoupon = useCallback(async (couponCode: string) => {
    if (!cart?.id) {
      throw new Error('No active cart found');
    }

    try {
      await applyCouponMutation({
        cartId: cart.id,
        coupon: { couponCode },
      }).unwrap();
    } catch (error) {
      console.error('Failed to apply coupon:', error);
      throw error;
    }
  }, [applyCouponMutation, cart?.id]);

  /**
   * Remove coupon from cart
   */
  const removeCoupon = useCallback(async () => {
    if (!cart?.id) {
      throw new Error('No active cart found');
    }

    try {
      await removeCouponMutation(cart.id).unwrap();
    } catch (error) {
      console.error('Failed to remove coupon:', error);
      throw error;
    }
  }, [removeCouponMutation, cart?.id]);

  // ===== Cart Utilities =====

  /**
   * Get cart item by product and variant ID
   */
  const getCartItem = useCallback((productId: string, variantId?: string): CartItem | undefined => {
    if (!cart?.items) return undefined;

    return cart.items.find(item =>
      item.productId === productId &&
      (variantId ? item.variantId === variantId : !item.variantId)
    );
  }, [cart?.items]);

  /**
   * Get quantity of specific product in cart
   */
  const getCartItemQuantity = useCallback((productId: string, variantId?: string): number => {
    const item = getCartItem(productId, variantId);
    return item?.quantity || 0;
  }, [getCartItem]);

  /**
   * Check if product is in cart
   */
  const isInCart = useCallback((productId: string, variantId?: string): boolean => {
    return getCartItemQuantity(productId, variantId) > 0;
  }, [getCartItemQuantity]);

  // ===== Computed Values =====

  const cartTotals = useMemo(() => {
    if (!cart) {
      return {
        itemCount: 0,
        subtotal: '0.00',
        tax: '0.00',
        shipping: '0.00',
        discount: '0.00',
        total: '0.00',
      };
    }

    return {
      itemCount: cart.itemCount,
      subtotal: cart.subtotalAmount,
      tax: cart.taxAmount,
      shipping: cart.shippingAmount,
      discount: cart.discountAmount,
      total: cart.totalAmount,
    };
  }, [cart]);

  // ===== Effects =====

  /**
   * Auto-merge carts when user logs in
   */
  useEffect(() => {
    if (userId && sessionId && cart) {
      // Check if we have a guest cart that needs merging
      const shouldMerge = !cart.userId && cart.sessionId === sessionId;

      if (shouldMerge) {
        mergeCartsMutation({ userId, sessionId })
          .unwrap()
          .then(() => {
            refetchCart();
          })
          .catch((error) => {
            console.error('Failed to merge carts:', error);
          });
      }
    }
  }, [userId, sessionId, cart, mergeCartsMutation, refetchCart]);

  // ===== Context Value =====

  const contextValue: CartContextValue = {
    // State
    cart: cart || null,
    cartValidation: cartValidation || null,
    isLoading,
    error: error ? String(error) : null,

    // Operations
    addToCart,
    updateCartItem,
    removeFromCart,
    applyCoupon,
    removeCoupon,

    // Utilities
    getCartItem,
    getCartItemQuantity,
    isInCart,

    // Totals
    ...cartTotals,
  };

  return (
    <CartContext.Provider value={contextValue}>
      {children}
    </CartContext.Provider>
  );
};

// ===== Export Types =====

export type { CartContextValue, CartProviderProps };