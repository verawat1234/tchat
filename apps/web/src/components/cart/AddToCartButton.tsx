/**
 * AddToCartButton - Product Integration Component
 *
 * A comprehensive add-to-cart button with quantity selection, variant support,
 * and optimistic updates. Designed for seamless product integration.
 *
 * Features:
 * - Quantity selection with validation
 * - Product variant support
 * - Stock checking and validation
 * - Optimistic updates with fallback
 * - Loading states and animations
 * - Error handling with retry
 * - Accessibility compliance
 * - Multiple display modes
 */

import React, { useState, useCallback, useMemo, useEffect } from 'react';
import { cn } from '../../lib/utils';
import { useCart } from './CartProvider';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Badge } from '../ui/badge/badge';
import { Popover, PopoverContent, PopoverTrigger } from '../ui/popover';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Separator } from '../ui/separator';
import {
  AlertCircle,
  Check,
  ChevronDown,
  Minus,
  Package,
  Plus,
  ShoppingCart,
  Truck,
} from 'lucide-react';
import type { AddToCartRequest, Product } from '../../types/commerce';

// ===== Types =====

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

export interface AddToCartButtonProps {
  /** Product information */
  product: Product;
  /** Available product variants */
  variants?: ProductVariant[];
  /** Default selected variant ID */
  defaultVariantId?: string;
  /** Default quantity */
  defaultQuantity?: number;
  /** Maximum quantity allowed */
  maxQuantity?: number;
  /** Button display mode */
  mode?: 'button' | 'inline' | 'dropdown' | 'detailed';
  /** Button size */
  size?: 'sm' | 'default' | 'lg';
  /** Button variant */
  variant?: 'default' | 'outline' | 'secondary' | 'ghost';
  /** Whether to show quantity selector */
  showQuantity?: boolean;
  /** Whether to show variant selector */
  showVariants?: boolean;
  /** Whether to show stock information */
  showStock?: boolean;
  /** Whether to show price */
  showPrice?: boolean;
  /** Custom button text */
  buttonText?: string;
  /** Custom added to cart text */
  addedText?: string;
  /** Custom class name */
  className?: string;
  /** Success callback */
  onSuccess?: (item: AddToCartRequest) => void;
  /** Error callback */
  onError?: (error: unknown) => void;
  /** Disabled state */
  disabled?: boolean;
  /** Loading state */
  loading?: boolean;
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
const calculateSavings = (price: string, comparePrice?: string): number => {
  if (!comparePrice) return 0;
  const current = parseFloat(price);
  const original = parseFloat(comparePrice);
  return original > current ? original - current : 0;
};

// ===== Component =====

/**
 * AddToCartButton - Comprehensive add-to-cart functionality
 *
 * @param product - Product data
 * @param variants - Available variants
 * @param defaultVariantId - Default selected variant
 * @param defaultQuantity - Default quantity
 * @param maxQuantity - Maximum allowed quantity
 * @param mode - Display mode
 * @param size - Button size
 * @param variant - Button variant
 * @param showQuantity - Show quantity selector
 * @param showVariants - Show variant selector
 * @param showStock - Show stock information
 * @param showPrice - Show price display
 * @param buttonText - Custom button text
 * @param addedText - Custom success text
 * @param className - Custom styling
 * @param onSuccess - Success callback
 * @param onError - Error callback
 * @param disabled - Disabled state
 * @param loading - Loading state
 */
export const AddToCartButton: React.FC<AddToCartButtonProps> = ({
  product,
  variants = [],
  defaultVariantId,
  defaultQuantity = 1,
  maxQuantity,
  mode = 'button',
  size = 'default',
  variant = 'default',
  showQuantity = true,
  showVariants = true,
  showStock = false,
  showPrice = false,
  buttonText = 'Add to Cart',
  addedText = 'Added to Cart',
  className,
  onSuccess,
  onError,
  disabled = false,
  loading = false,
}) => {
  const { addToCart, isLoading, getCartItemQuantity, isInCart } = useCart();

  // ===== Local State =====

  const [selectedVariantId, setSelectedVariantId] = useState<string | undefined>(
    defaultVariantId || variants[0]?.id
  );
  const [quantity, setQuantity] = useState(defaultQuantity);
  const [isAdding, setIsAdding] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  // ===== Computed Values =====

  const selectedVariant = useMemo(
    () => variants.find(v => v.id === selectedVariantId),
    [variants, selectedVariantId]
  );

  const currentStock = selectedVariant?.stock ?? product.inventory.quantity;
  const currentMaxQuantity = Math.min(
    maxQuantity ?? Infinity,
    selectedVariant?.maxQuantity ?? product.inventory.maxPerOrder,
    currentStock
  );

  const isOutOfStock = currentStock === 0 || !product.isActive;
  const isLowStock = currentStock > 0 && currentStock <= 5;
  const exceedsStock = quantity > currentStock;
  const exceedsMax = quantity > currentMaxQuantity;

  const currentPrice = selectedVariant?.price ?? product.price;
  const comparePrice = selectedVariant?.compareAtPrice ?? product.compareAtPrice;
  const savings = calculateSavings(currentPrice, comparePrice);

  const cartQuantity = getCartItemQuantity(product.id, selectedVariantId);
  const inCart = isInCart(product.id, selectedVariantId);

  const canAddToCart = !isOutOfStock && !exceedsStock && !exceedsMax && !disabled && !loading;

  // ===== Effects =====

  useEffect(() => {
    if (showSuccess) {
      const timer = setTimeout(() => setShowSuccess(false), 2000);
      return () => clearTimeout(timer);
    }
  }, [showSuccess]);

  useEffect(() => {
    setError(null);
  }, [selectedVariantId, quantity]);

  // ===== Handlers =====

  /**
   * Handle add to cart action
   */
  const handleAddToCart = useCallback(async () => {
    if (!canAddToCart) return;

    setIsAdding(true);
    setError(null);

    try {
      const item: AddToCartRequest = {
        productId: product.id,
        variantId: selectedVariantId,
        quantity,
      };

      await addToCart(item);
      setShowSuccess(true);
      onSuccess?.(item);

      // Reset quantity in button mode
      if (mode === 'button') {
        setQuantity(defaultQuantity);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to add item to cart';
      setError(errorMessage);
      onError?.(err);
    } finally {
      setIsAdding(false);
    }
  }, [canAddToCart, product.id, selectedVariantId, quantity, addToCart, onSuccess, onError, mode, defaultQuantity]);

  /**
   * Handle quantity change with validation
   */
  const handleQuantityChange = useCallback((newQuantity: number) => {
    if (newQuantity < 1) return;
    if (newQuantity > currentMaxQuantity) return;
    setQuantity(newQuantity);
  }, [currentMaxQuantity]);

  // ===== Render Helpers =====

  const renderQuantitySelector = () => {
    if (!showQuantity || mode === 'button') return null;

    return (
      <div className="flex items-center gap-1">
        <Button
          variant="outline"
          size="icon"
          className="h-8 w-8"
          onClick={() => handleQuantityChange(quantity - 1)}
          disabled={quantity <= 1 || isAdding}
          aria-label="Decrease quantity"
        >
          <Minus className="h-3 w-3" />
        </Button>

        <Input
          type="number"
          value={quantity}
          onChange={(e) => {
            const value = parseInt(e.target.value, 10);
            if (!isNaN(value) && value > 0) {
              handleQuantityChange(value);
            }
          }}
          className="w-16 h-8 text-center"
          min={1}
          max={currentMaxQuantity}
          disabled={isAdding}
          aria-label="Quantity"
        />

        <Button
          variant="outline"
          size="icon"
          className="h-8 w-8"
          onClick={() => handleQuantityChange(quantity + 1)}
          disabled={quantity >= currentMaxQuantity || isAdding}
          aria-label="Increase quantity"
        >
          <Plus className="h-3 w-3" />
        </Button>
      </div>
    );
  };

  const renderVariantSelector = () => {
    if (!showVariants || variants.length === 0) return null;

    return (
      <Select
        value={selectedVariantId}
        onValueChange={setSelectedVariantId}
        disabled={isAdding}
      >
        <SelectTrigger className="w-full">
          <SelectValue placeholder="Select variant" />
        </SelectTrigger>
        <SelectContent>
          {variants.map((variant) => (
            <SelectItem
              key={variant.id}
              value={variant.id}
              disabled={!variant.isAvailable}
            >
              <div className="flex items-center justify-between w-full">
                <span>{variant.name}</span>
                <div className="flex items-center gap-2 ml-2">
                  <span className="font-medium">
                    {formatCurrency(variant.price)}
                  </span>
                  {variant.compareAtPrice && (
                    <span className="text-xs text-muted-foreground line-through">
                      {formatCurrency(variant.compareAtPrice)}
                    </span>
                  )}
                  {variant.stock <= 5 && variant.stock > 0 && (
                    <Badge variant="secondary" className="text-xs">
                      {variant.stock} left
                    </Badge>
                  )}
                </div>
              </div>
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    );
  };

  const renderPriceInfo = () => {
    if (!showPrice) return null;

    return (
      <div className="flex items-center justify-between">
        <div>
          <span className="font-semibold text-lg">
            {formatCurrency(currentPrice)}
          </span>
          {comparePrice && (
            <span className="text-sm text-muted-foreground line-through ml-2">
              {formatCurrency(comparePrice)}
            </span>
          )}
        </div>
        {savings > 0 && (
          <Badge variant="default" className="bg-success text-success-foreground">
            Save {formatCurrency(savings.toString())}
          </Badge>
        )}
      </div>
    );
  };

  const renderStockInfo = () => {
    if (!showStock) return null;

    if (isOutOfStock) {
      return (
        <div className="flex items-center gap-2 text-destructive">
          <AlertCircle className="w-4 h-4" />
          <span className="text-sm">Out of stock</span>
        </div>
      );
    }

    if (isLowStock) {
      return (
        <div className="flex items-center gap-2 text-warning-foreground">
          <Package className="w-4 h-4" />
          <span className="text-sm">Only {currentStock} left</span>
        </div>
      );
    }

    return (
      <div className="flex items-center gap-2 text-success-foreground">
        <Truck className="w-4 h-4" />
        <span className="text-sm">In stock</span>
      </div>
    );
  };

  const renderError = () => {
    if (!error) return null;

    return (
      <div className="flex items-center gap-2 text-destructive text-sm">
        <AlertCircle className="w-4 h-4" />
        <span>{error}</span>
      </div>
    );
  };

  const renderMainButton = () => {
    const isButtonLoading = isAdding || loading || isLoading;

    const buttonContent = () => {
      if (showSuccess) {
        return (
          <>
            <Check className="w-4 h-4 mr-2" />
            {addedText}
          </>
        );
      }

      if (isButtonLoading) {
        return (
          <>
            <div className="w-4 h-4 mr-2 animate-spin rounded-full border-2 border-current border-t-transparent" />
            Adding...
          </>
        );
      }

      if (inCart && mode === 'button') {
        return (
          <>
            <Check className="w-4 h-4 mr-2" />
            In Cart ({cartQuantity})
          </>
        );
      }

      return (
        <>
          <ShoppingCart className="w-4 h-4 mr-2" />
          {buttonText}
        </>
      );
    };

    return (
      <Button
        variant={showSuccess ? 'default' : variant}
        size={size}
        onClick={handleAddToCart}
        disabled={!canAddToCart || isButtonLoading}
        className={cn(
          showSuccess && "bg-success text-success-foreground hover:bg-success/90",
          className
        )}
        aria-label={`${buttonText} - ${formatCurrency(currentPrice)}`}
      >
        {buttonContent()}
      </Button>
    );
  };

  // ===== Mode-specific Renders =====

  if (mode === 'button') {
    return (
      <div className="space-y-2">
        {renderMainButton()}
        {renderError()}
      </div>
    );
  }

  if (mode === 'inline') {
    return (
      <div className={cn("flex items-center gap-3", className)}>
        {renderQuantitySelector()}
        {renderMainButton()}
        {renderError()}
      </div>
    );
  }

  if (mode === 'dropdown') {
    return (
      <Popover open={isDropdownOpen} onOpenChange={setIsDropdownOpen}>
        <PopoverTrigger asChild>
          <Button variant={variant} size={size} className={className}>
            <ShoppingCart className="w-4 h-4 mr-2" />
            {buttonText}
            <ChevronDown className="w-4 h-4 ml-2" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-80 p-4" side="bottom" align="start">
          <div className="space-y-4">
            {renderPriceInfo()}
            {renderVariantSelector()}
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Quantity:</span>
              {renderQuantitySelector()}
            </div>
            {renderStockInfo()}
            <Separator />
            {renderMainButton()}
            {renderError()}
          </div>
        </PopoverContent>
      </Popover>
    );
  }

  if (mode === 'detailed') {
    return (
      <div className={cn("space-y-4 p-4 border rounded-lg", className)}>
        {renderPriceInfo()}
        {showVariants && variants.length > 0 && (
          <div>
            <label className="text-sm font-medium mb-2 block">Variant:</label>
            {renderVariantSelector()}
          </div>
        )}
        <div className="flex items-center justify-between">
          <label className="text-sm font-medium">Quantity:</label>
          {renderQuantitySelector()}
        </div>
        {renderStockInfo()}
        <Separator />
        {renderMainButton()}
        {renderError()}
      </div>
    );
  }

  return renderMainButton();
};

// ===== Export Component =====

export default AddToCartButton;