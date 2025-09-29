/**
 * CouponInput - Discount Code Input Component
 *
 * Provides coupon code application and management functionality.
 * Integrates with CartProvider for cart-specific coupon operations.
 *
 * Features:
 * - Coupon code input and validation
 * - Real-time coupon verification
 * - Applied coupon display and removal
 * - Discount amount calculation
 * - Error handling and feedback
 * - Loading states
 * - Accessibility compliance
 */

import React, { useState, useCallback, useMemo, useEffect } from 'react';
import { cn } from '../../lib/utils';
import { useCart } from './CartProvider';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Badge } from '../ui/badge/badge';
import { Label } from '../ui/label';
import { Alert, AlertDescription } from '../ui/alert';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '../ui/collapsible';
import {
  AlertCircle,
  Check,
  ChevronDown,
  ChevronUp,
  Loader2,
  Percent,
  Tag,
  X,
  Gift,
  Zap,
} from 'lucide-react';

// ===== Types =====

export interface CouponInputProps {
  /** Placeholder text for input */
  placeholder?: string;
  /** Whether to show applied coupon info */
  showAppliedCoupon?: boolean;
  /** Whether to show discount breakdown */
  showDiscountDetails?: boolean;
  /** Whether to show coupon suggestions */
  showSuggestions?: boolean;
  /** Whether to collapse when no coupon applied */
  collapsible?: boolean;
  /** Default collapsed state */
  defaultCollapsed?: boolean;
  /** Input size */
  size?: 'sm' | 'default' | 'lg';
  /** Custom class name */
  className?: string;
  /** Compact display mode */
  compact?: boolean;
  /** Success callback */
  onCouponApplied?: (couponCode: string, discount: string) => void;
  /** Remove callback */
  onCouponRemoved?: (couponCode: string) => void;
  /** Error callback */
  onError?: (error: string) => void;
}

// ===== Helper Functions =====

/**
 * Format currency amount
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
 * Calculate discount percentage
 */
const calculateDiscountPercentage = (discount: string, subtotal: string): number => {
  const discountAmount = parseFloat(discount);
  const subtotalAmount = parseFloat(subtotal);

  if (subtotalAmount === 0) return 0;
  return Math.round((discountAmount / subtotalAmount) * 100);
};

/**
 * Validate coupon code format
 */
const validateCouponCode = (code: string): string | null => {
  if (!code.trim()) {
    return 'Please enter a coupon code';
  }

  if (code.length < 3) {
    return 'Coupon code must be at least 3 characters';
  }

  if (code.length > 50) {
    return 'Coupon code is too long';
  }

  if (!/^[A-Za-z0-9-_]+$/.test(code)) {
    return 'Coupon code can only contain letters, numbers, hyphens, and underscores';
  }

  return null;
};

// ===== Mock Coupon Suggestions =====
const COUPON_SUGGESTIONS = [
  { code: 'WELCOME10', description: '10% off your first order', type: 'percentage' },
  { code: 'FREESHIP', description: 'Free shipping on orders over $50', type: 'shipping' },
  { code: 'SAVE20', description: '$20 off orders over $100', type: 'fixed' },
];

// ===== Component =====

/**
 * CouponInput - Discount code input and management
 *
 * @param placeholder - Input placeholder text
 * @param showAppliedCoupon - Show applied coupon details
 * @param showDiscountDetails - Show discount breakdown
 * @param showSuggestions - Show coupon suggestions
 * @param collapsible - Enable collapsible behavior
 * @param defaultCollapsed - Default collapsed state
 * @param size - Input size
 * @param className - Custom styling
 * @param compact - Use compact layout
 * @param onCouponApplied - Applied callback
 * @param onCouponRemoved - Removed callback
 * @param onError - Error callback
 */
export const CouponInput: React.FC<CouponInputProps> = ({
  placeholder = 'Enter coupon code',
  showAppliedCoupon = true,
  showDiscountDetails = true,
  showSuggestions = false,
  collapsible = false,
  defaultCollapsed = false,
  size = 'default',
  className,
  compact = false,
  onCouponApplied,
  onCouponRemoved,
  onError,
}) => {
  const {
    cart,
    applyCoupon,
    removeCoupon,
    isLoading,
    subtotal,
    discount,
  } = useCart();

  // ===== Local State =====

  const [couponCode, setCouponCode] = useState('');
  const [isApplying, setIsApplying] = useState(false);
  const [isRemoving, setIsRemoving] = useState(false);
  const [validationError, setValidationError] = useState<string | null>(null);
  const [apiError, setApiError] = useState<string | null>(null);
  const [isCollapsed, setIsCollapsed] = useState(defaultCollapsed);
  const [showSuggestionsExpanded, setShowSuggestionsExpanded] = useState(false);

  // ===== Computed Values =====

  const appliedCoupon = cart?.couponCode;
  const appliedDiscount = cart?.couponDiscount || '0';
  const hasCoupon = !!appliedCoupon;
  const hasDiscount = parseFloat(appliedDiscount) > 0;

  const discountPercentage = useMemo(() => {
    if (!hasDiscount) return 0;
    return calculateDiscountPercentage(appliedDiscount, subtotal);
  }, [hasDiscount, appliedDiscount, subtotal]);

  const currency = cart?.currency || 'USD';

  const isProcessing = isApplying || isRemoving || isLoading;

  // ===== Effects =====

  useEffect(() => {
    // Clear validation error when input changes
    if (validationError && couponCode) {
      setValidationError(null);
    }
  }, [couponCode, validationError]);

  useEffect(() => {
    // Clear API error after a delay
    if (apiError) {
      const timer = setTimeout(() => setApiError(null), 5000);
      return () => clearTimeout(timer);
    }
  }, [apiError]);

  // ===== Handlers =====

  /**
   * Apply coupon code
   */
  const handleApplyCoupon = useCallback(async () => {
    const trimmedCode = couponCode.trim().toUpperCase();

    // Validate coupon code format
    const error = validateCouponCode(trimmedCode);
    if (error) {
      setValidationError(error);
      return;
    }

    setIsApplying(true);
    setApiError(null);
    setValidationError(null);

    try {
      await applyCoupon(trimmedCode);
      setCouponCode(''); // Clear input after successful application
      onCouponApplied?.(trimmedCode, appliedDiscount);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to apply coupon';
      setApiError(errorMessage);
      onError?.(errorMessage);
    } finally {
      setIsApplying(false);
    }
  }, [couponCode, applyCoupon, appliedDiscount, onCouponApplied, onError]);

  /**
   * Remove applied coupon
   */
  const handleRemoveCoupon = useCallback(async () => {
    if (!appliedCoupon) return;

    setIsRemoving(true);
    setApiError(null);

    try {
      await removeCoupon();
      onCouponRemoved?.(appliedCoupon);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to remove coupon';
      setApiError(errorMessage);
      onError?.(errorMessage);
    } finally {
      setIsRemoving(false);
    }
  }, [appliedCoupon, removeCoupon, onCouponRemoved, onError]);

  /**
   * Apply suggested coupon
   */
  const handleApplySuggestion = useCallback(async (suggestionCode: string) => {
    setCouponCode(suggestionCode);
    // Auto-apply the suggestion
    setIsApplying(true);
    setApiError(null);
    setValidationError(null);

    try {
      await applyCoupon(suggestionCode);
      setCouponCode('');
      setShowSuggestionsExpanded(false);
      onCouponApplied?.(suggestionCode, appliedDiscount);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to apply coupon';
      setApiError(errorMessage);
      onError?.(errorMessage);
    } finally {
      setIsApplying(false);
    }
  }, [applyCoupon, appliedDiscount, onCouponApplied, onError]);

  /**
   * Handle Enter key
   */
  const handleKeyPress = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && couponCode.trim() && !isProcessing) {
      handleApplyCoupon();
    }
  }, [couponCode, isProcessing, handleApplyCoupon]);

  // ===== Render Helpers =====

  const renderAppliedCoupon = () => {
    if (!showAppliedCoupon || !hasCoupon) return null;

    return (
      <div className="flex items-center justify-between p-3 bg-success/10 border border-success/20 rounded-lg">
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-8 h-8 bg-success/20 rounded-full">
            <Tag className="w-4 h-4 text-success-foreground" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <span className="font-medium text-success-foreground">
                {appliedCoupon}
              </span>
              <Badge variant="default" className="bg-success text-success-foreground">
                Applied
              </Badge>
            </div>
            {hasDiscount && showDiscountDetails && (
              <div className="text-sm text-success-foreground/80">
                Save {formatCurrency(appliedDiscount, currency)}
                {discountPercentage > 0 && ` (${discountPercentage}% off)`}
              </div>
            )}
          </div>
        </div>
        <Button
          variant="ghost"
          size="sm"
          onClick={handleRemoveCoupon}
          disabled={isRemoving}
          className="text-success-foreground hover:text-success-foreground/80"
          aria-label={`Remove coupon ${appliedCoupon}`}
        >
          {isRemoving ? (
            <Loader2 className="w-4 h-4 animate-spin" />
          ) : (
            <X className="w-4 h-4" />
          )}
        </Button>
      </div>
    );
  };

  const renderCouponInput = () => {
    if (hasCoupon && !showSuggestions) return null;

    return (
      <div className="space-y-3">
        {!compact && (
          <Label htmlFor="coupon-input" className="text-sm font-medium">
            Discount Code
          </Label>
        )}

        <div className="flex gap-2">
          <div className="flex-1">
            <Input
              id="coupon-input"
              type="text"
              placeholder={placeholder}
              value={couponCode}
              onChange={(e) => setCouponCode(e.target.value.toUpperCase())}
              onKeyPress={handleKeyPress}
              disabled={isProcessing}
              className={cn(
                validationError && "border-destructive focus-visible:ring-destructive",
                size === 'sm' && "h-8 text-sm",
                size === 'lg' && "h-12 text-lg"
              )}
              aria-describedby={validationError ? "coupon-error" : undefined}
            />
          </div>

          <Button
            variant="outline"
            size={size}
            onClick={handleApplyCoupon}
            disabled={!couponCode.trim() || isProcessing}
            className="flex items-center gap-2"
          >
            {isApplying ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Tag className="w-4 h-4" />
            )}
            {compact ? 'Apply' : 'Apply Coupon'}
          </Button>
        </div>
      </div>
    );
  };

  const renderSuggestions = () => {
    if (!showSuggestions || hasCoupon) return null;

    return (
      <Collapsible open={showSuggestionsExpanded} onOpenChange={setShowSuggestionsExpanded}>
        <CollapsibleTrigger asChild>
          <Button variant="ghost" className="w-full justify-between p-0 h-auto text-sm">
            <span className="flex items-center gap-2">
              <Gift className="w-4 h-4" />
              Available Offers
            </span>
            {showSuggestionsExpanded ? (
              <ChevronUp className="w-4 h-4" />
            ) : (
              <ChevronDown className="w-4 h-4" />
            )}
          </Button>
        </CollapsibleTrigger>
        <CollapsibleContent className="mt-3">
          <div className="space-y-2">
            {COUPON_SUGGESTIONS.map((suggestion) => (
              <div
                key={suggestion.code}
                className="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="flex items-center justify-center w-8 h-8 bg-primary/10 rounded-full">
                    {suggestion.type === 'percentage' && <Percent className="w-4 h-4 text-primary" />}
                    {suggestion.type === 'shipping' && <Truck className="w-4 h-4 text-primary" />}
                    {suggestion.type === 'fixed' && <Zap className="w-4 h-4 text-primary" />}
                  </div>
                  <div>
                    <div className="font-medium text-sm">{suggestion.code}</div>
                    <div className="text-xs text-muted-foreground">{suggestion.description}</div>
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleApplySuggestion(suggestion.code)}
                  disabled={isProcessing}
                >
                  Apply
                </Button>
              </div>
            ))}
          </div>
        </CollapsibleContent>
      </Collapsible>
    );
  };

  const renderErrors = () => {
    const error = validationError || apiError;
    if (!error) return null;

    return (
      <Alert variant="destructive" className="mt-3">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription id="coupon-error">
          {error}
        </AlertDescription>
      </Alert>
    );
  };

  const renderCollapsibleContent = () => {
    return (
      <div className="space-y-4">
        {renderCouponInput()}
        {renderSuggestions()}
        {renderErrors()}
      </div>
    );
  };

  // ===== Main Render =====

  if (collapsible && !hasCoupon) {
    return (
      <div className={className}>
        {renderAppliedCoupon()}

        <Collapsible open={!isCollapsed} onOpenChange={(open) => setIsCollapsed(!open)}>
          <CollapsibleTrigger asChild>
            <Button variant="ghost" className="w-full justify-between p-0 h-auto">
              <span className="flex items-center gap-2">
                <Tag className="w-4 h-4" />
                Add Discount Code
              </span>
              {isCollapsed ? (
                <ChevronDown className="w-4 h-4" />
              ) : (
                <ChevronUp className="w-4 h-4" />
              )}
            </Button>
          </CollapsibleTrigger>
          <CollapsibleContent className="mt-3">
            {renderCollapsibleContent()}
          </CollapsibleContent>
        </Collapsible>
      </div>
    );
  }

  return (
    <div className={className}>
      {renderAppliedCoupon()}
      {!hasCoupon && renderCollapsibleContent()}
    </div>
  );
};

// ===== Export Component =====

export default CouponInput;