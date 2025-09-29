/**
 * CartSummary - Cart Overview Component
 *
 * Displays cart totals, summary information, and provides quick actions.
 * Integrates with CartProvider for real-time cart state.
 *
 * Features:
 * - Real-time cart totals display
 * - Coupon code application
 * - Cart validation status
 * - Loading and error states
 * - Responsive design
 * - Accessibility compliance
 */

import React, { useMemo } from 'react';
import { cn } from '../../lib/utils';
import { useCart } from './CartProvider';
import { TchatCard, TchatCardHeader, TchatCardContent, TchatCardFooter } from '../TchatCard';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Separator } from '../ui/separator';
import { AlertCircle, CheckCircle, ShoppingCart, Tag, Truck } from 'lucide-react';

// ===== Types =====

export interface CartSummaryProps {
  /** Whether to show the checkout button */
  showCheckout?: boolean;
  /** Checkout button text */
  checkoutText?: string;
  /** Checkout button click handler */
  onCheckout?: () => void;
  /** Whether to show detailed breakdown */
  showDetails?: boolean;
  /** Whether to show cart validation status */
  showValidation?: boolean;
  /** Custom class name */
  className?: string;
  /** Compact mode for smaller spaces */
  compact?: boolean;
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
 * Calculate savings amount
 */
const calculateSavings = (subtotal: string, discount: string): number => {
  const subtotalNum = parseFloat(subtotal);
  const discountNum = parseFloat(discount);
  return discountNum > 0 ? discountNum : 0;
};

// ===== Component =====

/**
 * CartSummary - Displays cart totals and summary information
 *
 * @param showCheckout - Whether to show checkout button
 * @param checkoutText - Text for checkout button
 * @param onCheckout - Checkout handler
 * @param showDetails - Whether to show detailed breakdown
 * @param showValidation - Whether to show validation status
 * @param className - Custom styling
 * @param compact - Use compact layout
 */
export const CartSummary: React.FC<CartSummaryProps> = ({
  showCheckout = true,
  checkoutText = 'Proceed to Checkout',
  onCheckout,
  showDetails = true,
  showValidation = true,
  className,
  compact = false,
}) => {
  const {
    cart,
    cartValidation,
    isLoading,
    error,
    itemCount,
    subtotal,
    tax,
    shipping,
    discount,
    total,
  } = useCart();

  // ===== Computed Values =====

  const isEmpty = itemCount === 0;
  const hasDiscount = parseFloat(discount) > 0;
  const savings = calculateSavings(subtotal, discount);

  const currency = cart?.currency || 'USD';

  const validationStatus = useMemo(() => {
    if (!cartValidation) return null;

    const hasErrors = cartValidation.issues.some(issue => issue.severity === 'error');
    const hasWarnings = cartValidation.issues.some(issue => issue.severity === 'warning');

    if (hasErrors) {
      return { type: 'error', message: 'Cart has validation errors', icon: AlertCircle };
    }
    if (hasWarnings) {
      return { type: 'warning', message: 'Cart has warnings', icon: AlertCircle };
    }
    return { type: 'success', message: 'Cart is valid', icon: CheckCircle };
  }, [cartValidation]);

  // ===== Render Helpers =====

  const renderSummaryLine = (label: string, amount: string, highlight = false, icon?: React.ComponentType<any>) => {
    const IconComponent = icon;

    return (
      <div className={cn(
        'flex items-center justify-between',
        highlight && 'text-lg font-semibold text-primary',
        compact && 'text-sm'
      )}>
        <div className="flex items-center gap-2">
          {IconComponent && <IconComponent className="w-4 h-4" />}
          <span>{label}</span>
        </div>
        <span className={cn(highlight && 'font-bold')}>
          {formatCurrency(amount, currency)}
        </span>
      </div>
    );
  };

  const renderValidationBadge = () => {
    if (!showValidation || !validationStatus) return null;

    const { type, message, icon: Icon } = validationStatus;
    const variant = type === 'error' ? 'destructive' : type === 'warning' ? 'secondary' : 'default';

    return (
      <Badge variant={variant} className="flex items-center gap-1">
        <Icon className="w-3 h-3" />
        <span className="text-xs">{message}</span>
      </Badge>
    );
  };

  // ===== Loading State =====

  if (isLoading) {
    return (
      <TchatCard className={cn('animate-pulse', className)} size={compact ? 'compact' : 'standard'}>
        <TchatCardContent>
          <div className="space-y-3">
            <div className="h-4 bg-muted rounded w-1/2" />
            <div className="h-4 bg-muted rounded w-3/4" />
            <div className="h-4 bg-muted rounded w-2/3" />
            <div className="h-6 bg-muted rounded w-full" />
          </div>
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== Error State =====

  if (error) {
    return (
      <TchatCard variant="outlined" className={cn('border-destructive', className)}>
        <TchatCardContent>
          <div className="flex items-center gap-2 text-destructive">
            <AlertCircle className="w-4 h-4" />
            <span className="text-sm">Failed to load cart summary</span>
          </div>
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== Empty State =====

  if (isEmpty) {
    return (
      <TchatCard variant="outlined" className={cn('text-center', className)}>
        <TchatCardContent className="py-8">
          <ShoppingCart className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
          <h3 className="font-medium text-foreground mb-2">Your cart is empty</h3>
          <p className="text-sm text-muted-foreground">
            Add some items to get started
          </p>
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== Main Content =====

  return (
    <TchatCard
      variant="elevated"
      className={className}
      size={compact ? 'compact' : 'standard'}
    >
      <TchatCardHeader
        title={compact ? undefined : 'Order Summary'}
        subtitle={compact ? undefined : `${itemCount} ${itemCount === 1 ? 'item' : 'items'}`}
        actions={renderValidationBadge()}
      />

      <TchatCardContent>
        <div className={cn('space-y-3', compact && 'space-y-2')}>
          {/* Basic totals - always shown */}
          {renderSummaryLine('Subtotal', subtotal)}

          {/* Detailed breakdown - conditional */}
          {showDetails && (
            <>
              {hasDiscount && (
                <div className="space-y-2">
                  {renderSummaryLine('Discount', `-${discount}`, false, Tag)}
                  <div className="text-xs text-success-foreground bg-success/10 px-2 py-1 rounded">
                    You saved {formatCurrency(savings.toString(), currency)}!
                  </div>
                </div>
              )}

              {parseFloat(tax) > 0 && renderSummaryLine('Tax', tax)}
              {parseFloat(shipping) > 0 && renderSummaryLine('Shipping', shipping, false, Truck)}

              <Separator className="my-3" />
            </>
          )}

          {/* Total - always shown */}
          {renderSummaryLine('Total', total, true)}
        </div>
      </TchatCardContent>

      {showCheckout && onCheckout && (
        <TchatCardFooter>
          <Button
            className="w-full"
            size={compact ? 'sm' : 'default'}
            onClick={onCheckout}
            disabled={isEmpty || (validationStatus?.type === 'error')}
            aria-label={`${checkoutText} - ${formatCurrency(total, currency)}`}
          >
            {checkoutText}
          </Button>
        </TchatCardFooter>
      )}
    </TchatCard>
  );
};

// ===== Export Component =====

export default CartSummary;