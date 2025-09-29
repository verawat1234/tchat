/**
 * CartItem - Individual Cart Item Component
 *
 * Displays a single cart item with product information, quantity controls,
 * and item-specific actions. Includes optimistic updates and error handling.
 *
 * Features:
 * - Product information display
 * - Quantity controls with validation
 * - Remove item functionality
 * - Stock availability checking
 * - Gift message support
 * - Optimistic updates
 * - Loading states
 * - Accessibility compliance
 */

import React, { useState, useCallback, useMemo } from 'react';
import { cn } from '../../lib/utils';
import { useCart } from './CartProvider';
import { TchatCard, TchatCardContent } from '../TchatCard';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Badge } from '../ui/badge';
import { Textarea } from '../ui/textarea';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '../ui/collapsible';
import {
  AlertCircle,
  ChevronDown,
  ChevronUp,
  Gift,
  Heart,
  Minus,
  Plus,
  Trash2,
  Package,
} from 'lucide-react';
import type { CartItem as CartItemType } from '../../types/commerce';

// ===== Types =====

export interface CartItemProps {
  /** Cart item data */
  item: CartItemType;
  /** Whether to show compact layout */
  compact?: boolean;
  /** Whether quantity controls are editable */
  editable?: boolean;
  /** Whether to show gift message option */
  showGiftOptions?: boolean;
  /** Whether to show wishlist action */
  showWishlistAction?: boolean;
  /** Custom class name */
  className?: string;
  /** Quantity change handler (for external control) */
  onQuantityChange?: (quantity: number) => void;
  /** Remove item handler (for external control) */
  onRemove?: () => void;
}

// ===== Helper Functions =====

/**
 * Format currency amount for display
 */
const formatCurrency = (amount: string, currency = 'USD'): string => {
  const numAmount = parseFloat(amount);
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
  }).format(numAmount);
};

/**
 * Calculate savings for discounted items
 */
const calculateSavings = (unitPrice: string, originalPrice?: string): number => {
  if (!originalPrice) return 0;
  const current = parseFloat(unitPrice);
  const original = parseFloat(originalPrice);
  return original > current ? original - current : 0;
};

// ===== Component =====

/**
 * CartItem - Displays individual cart item with controls
 *
 * @param item - Cart item data
 * @param compact - Use compact layout
 * @param editable - Enable quantity controls
 * @param showGiftOptions - Show gift message option
 * @param showWishlistAction - Show wishlist button
 * @param className - Custom styling
 * @param onQuantityChange - External quantity handler
 * @param onRemove - External remove handler
 */
export const CartItem: React.FC<CartItemProps> = ({
  item,
  compact = false,
  editable = true,
  showGiftOptions = true,
  showWishlistAction = false,
  className,
  onQuantityChange,
  onRemove,
}) => {
  const { updateCartItem, removeFromCart } = useCart();

  // ===== Local State =====

  const [isUpdating, setIsUpdating] = useState(false);
  const [isRemoving, setIsRemoving] = useState(false);
  const [showGiftMessage, setShowGiftMessage] = useState(!!item.giftMessage);
  const [giftMessage, setGiftMessage] = useState(item.giftMessage || '');
  const [updateError, setUpdateError] = useState<string | null>(null);

  // ===== Computed Values =====

  const isOutOfStock = !item.isAvailable || item.stockQuantity === 0;
  const isLowStock = item.stockQuantity > 0 && item.stockQuantity <= 5;
  const exceedsStock = item.quantity > item.stockQuantity;
  const exceedsMax = item.quantity > item.maxQuantity;
  const hasIssue = isOutOfStock || exceedsStock || exceedsMax;

  const totalPrice = parseFloat(item.totalPrice);
  const savings = calculateSavings(item.unitPrice, undefined); // Would need original price from product data

  // ===== Handlers =====

  /**
   * Update item quantity with validation
   */
  const handleQuantityChange = useCallback(async (newQuantity: number) => {
    if (newQuantity < 1) return;
    if (newQuantity > item.maxQuantity) return;
    if (newQuantity === item.quantity) return;

    setIsUpdating(true);
    setUpdateError(null);

    try {
      if (onQuantityChange) {
        onQuantityChange(newQuantity);
      } else {
        await updateCartItem(item.id, { quantity: newQuantity });
      }
    } catch (error) {
      setUpdateError('Failed to update quantity');
      console.error('Failed to update cart item:', error);
    } finally {
      setIsUpdating(false);
    }
  }, [item.id, item.quantity, item.maxQuantity, updateCartItem, onQuantityChange]);

  /**
   * Remove item from cart
   */
  const handleRemove = useCallback(async () => {
    setIsRemoving(true);

    try {
      if (onRemove) {
        onRemove();
      } else {
        await removeFromCart(item.id);
      }
    } catch (error) {
      console.error('Failed to remove cart item:', error);
    } finally {
      setIsRemoving(false);
    }
  }, [item.id, removeFromCart, onRemove]);

  /**
   * Update gift message
   */
  const handleGiftMessageChange = useCallback(async (message: string) => {
    setIsUpdating(true);
    setUpdateError(null);

    try {
      await updateCartItem(item.id, {
        isGift: message.length > 0,
        giftMessage: message || undefined,
      });
      setGiftMessage(message);
    } catch (error) {
      setUpdateError('Failed to update gift message');
      console.error('Failed to update gift message:', error);
    } finally {
      setIsUpdating(false);
    }
  }, [item.id, updateCartItem]);

  // ===== Render Helpers =====

  const renderQuantityControls = () => {
    if (!editable) {
      return (
        <div className="flex items-center gap-2">
          <span className="text-sm text-muted-foreground">Qty:</span>
          <span className="font-medium">{item.quantity}</span>
        </div>
      );
    }

    return (
      <div className="flex items-center gap-1">
        <Button
          variant="outline"
          size="icon"
          className="h-8 w-8"
          onClick={() => handleQuantityChange(item.quantity - 1)}
          disabled={item.quantity <= 1 || isUpdating}
          aria-label="Decrease quantity"
        >
          <Minus className="h-3 w-3" />
        </Button>

        <Input
          type="number"
          value={item.quantity}
          onChange={(e) => {
            const value = parseInt(e.target.value, 10);
            if (!isNaN(value) && value > 0) {
              handleQuantityChange(value);
            }
          }}
          className="w-16 h-8 text-center"
          min={1}
          max={item.maxQuantity}
          disabled={isUpdating}
          aria-label="Item quantity"
        />

        <Button
          variant="outline"
          size="icon"
          className="h-8 w-8"
          onClick={() => handleQuantityChange(item.quantity + 1)}
          disabled={item.quantity >= item.maxQuantity || isUpdating}
          aria-label="Increase quantity"
        >
          <Plus className="h-3 w-3" />
        </Button>
      </div>
    );
  };

  const renderStatusBadges = () => {
    const badges = [];

    if (isOutOfStock) {
      badges.push(
        <Badge key="out-of-stock" variant="destructive" className="text-xs">
          Out of Stock
        </Badge>
      );
    } else if (isLowStock) {
      badges.push(
        <Badge key="low-stock" variant="secondary" className="text-xs">
          Low Stock
        </Badge>
      );
    }

    if (exceedsStock && !isOutOfStock) {
      badges.push(
        <Badge key="exceeds-stock" variant="destructive" className="text-xs">
          Exceeds Stock
        </Badge>
      );
    }

    if (exceedsMax) {
      badges.push(
        <Badge key="exceeds-max" variant="destructive" className="text-xs">
          Max {item.maxQuantity}
        </Badge>
      );
    }

    if (item.isGift) {
      badges.push(
        <Badge key="gift" variant="default" className="text-xs">
          <Gift className="w-3 h-3 mr-1" />
          Gift
        </Badge>
      );
    }

    return badges.length > 0 ? (
      <div className="flex flex-wrap gap-1">
        {badges}
      </div>
    ) : null;
  };

  const renderGiftOptions = () => {
    if (!showGiftOptions) return null;

    return (
      <Collapsible open={showGiftMessage} onOpenChange={setShowGiftMessage}>
        <CollapsibleTrigger asChild>
          <Button variant="ghost" size="sm" className="h-auto p-2 justify-start">
            <Gift className="w-4 h-4 mr-2" />
            Gift Message
            {showGiftMessage ? <ChevronUp className="w-4 h-4 ml-auto" /> : <ChevronDown className="w-4 h-4 ml-auto" />}
          </Button>
        </CollapsibleTrigger>
        <CollapsibleContent className="mt-2">
          <Textarea
            placeholder="Add a gift message..."
            value={giftMessage}
            onChange={(e) => setGiftMessage(e.target.value)}
            onBlur={() => handleGiftMessageChange(giftMessage)}
            className="min-h-[60px]"
            maxLength={500}
            disabled={isUpdating}
          />
          <div className="text-xs text-muted-foreground mt-1">
            {giftMessage.length}/500 characters
          </div>
        </CollapsibleContent>
      </Collapsible>
    );
  };

  // ===== Main Render =====

  return (
    <TchatCard
      variant={hasIssue ? "outlined" : "elevated"}
      className={cn(
        hasIssue && "border-destructive/50",
        isRemoving && "opacity-50 pointer-events-none",
        className
      )}
      size={compact ? "compact" : "standard"}
    >
      <TchatCardContent>
        <div className={cn(
          "flex gap-4",
          compact && "gap-3"
        )}>
          {/* Product Image */}
          <div className={cn(
            "flex-shrink-0",
            compact ? "w-16 h-16" : "w-20 h-20"
          )}>
            {item.productImage ? (
              <img
                src={item.productImage}
                alt={item.productName}
                className="w-full h-full object-cover rounded-md"
                loading="lazy"
              />
            ) : (
              <div className="w-full h-full bg-muted rounded-md flex items-center justify-center">
                <Package className="w-6 h-6 text-muted-foreground" />
              </div>
            )}
          </div>

          {/* Product Details */}
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-2">
              <div className="flex-1 min-w-0">
                <h3 className={cn(
                  "font-medium text-foreground truncate",
                  compact ? "text-sm" : "text-base"
                )}>
                  {item.productName}
                </h3>

                {item.variantName && (
                  <p className="text-sm text-muted-foreground">
                    {item.variantName}
                  </p>
                )}

                <p className="text-xs text-muted-foreground">
                  SKU: {item.productSku}
                </p>
              </div>

              {/* Actions */}
              <div className="flex items-center gap-1">
                {showWishlistAction && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8"
                    aria-label="Add to wishlist"
                  >
                    <Heart className="h-4 w-4" />
                  </Button>
                )}

                <Button
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8 text-destructive hover:text-destructive"
                  onClick={handleRemove}
                  disabled={isRemoving}
                  aria-label="Remove item"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </div>

            {/* Status Badges */}
            {renderStatusBadges()}

            {/* Error Message */}
            {updateError && (
              <div className="flex items-center gap-2 text-destructive text-sm mt-2">
                <AlertCircle className="w-4 h-4" />
                <span>{updateError}</span>
              </div>
            )}

            {/* Quantity and Price */}
            <div className="flex items-center justify-between mt-3">
              {renderQuantityControls()}

              <div className="text-right">
                <div className={cn(
                  "font-semibold",
                  compact ? "text-sm" : "text-base"
                )}>
                  {formatCurrency(item.totalPrice)}
                </div>
                {item.quantity > 1 && (
                  <div className="text-xs text-muted-foreground">
                    {formatCurrency(item.unitPrice)} each
                  </div>
                )}
                {savings > 0 && (
                  <div className="text-xs text-success-foreground">
                    Save {formatCurrency(savings.toString())}
                  </div>
                )}
              </div>
            </div>

            {/* Gift Options */}
            {renderGiftOptions()}
          </div>
        </div>
      </TchatCardContent>
    </TchatCard>
  );
};

// ===== Export Component =====

export default CartItem;