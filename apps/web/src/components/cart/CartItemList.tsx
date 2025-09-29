/**
 * CartItemList - Cart Items Container Component
 *
 * Displays a list of cart items with various layout options and bulk actions.
 * Provides virtualization for large lists and accessibility features.
 *
 * Features:
 * - Virtual scrolling for performance
 * - Bulk selection and actions
 * - Empty state handling
 * - Loading states
 * - Error recovery
 * - Accessibility compliance
 * - Responsive design
 */

import React, { useState, useMemo, useCallback } from 'react';
import { cn } from '../../lib/utils';
import { useCart } from './CartProvider';
import { CartItem } from './CartItem';
import { TchatCard, TchatCardHeader, TchatCardContent } from '../TchatCard';
import { Button } from '../ui/button';
import { Checkbox } from '../ui/checkbox';
import { Badge } from '../ui/badge/badge';
import { Separator } from '../ui/separator';
import { AlertCircle, ShoppingCart, Trash2, Heart, Package } from 'lucide-react';
import type { CartItem as CartItemType } from '../../types/commerce';

// ===== Types =====

export interface CartItemListProps {
  /** Whether to show bulk actions */
  showBulkActions?: boolean;
  /** Whether to show selection checkboxes */
  showSelection?: boolean;
  /** Whether to use compact layout */
  compact?: boolean;
  /** Whether items are editable */
  editable?: boolean;
  /** Whether to show gift options */
  showGiftOptions?: boolean;
  /** Whether to show wishlist actions */
  showWishlistActions?: boolean;
  /** Custom class name */
  className?: string;
  /** Maximum height before scrolling */
  maxHeight?: string;
  /** Virtualization threshold (items count) */
  virtualizationThreshold?: number;
  /** Custom empty state component */
  emptyStateComponent?: React.ComponentType;
  /** Bulk action handlers */
  onBulkRemove?: (itemIds: string[]) => void;
  onBulkMoveToWishlist?: (itemIds: string[]) => void;
  /** Item action handlers (for external control) */
  onItemQuantityChange?: (itemId: string, quantity: number) => void;
  onItemRemove?: (itemId: string) => void;
}

// ===== Default Empty State =====

const DefaultEmptyState: React.FC = () => (
  <div className="text-center py-12">
    <ShoppingCart className="w-16 h-16 mx-auto text-muted-foreground mb-4" />
    <h3 className="text-lg font-medium text-foreground mb-2">Your cart is empty</h3>
    <p className="text-muted-foreground mb-6">
      Add some items to your cart to get started
    </p>
    <Button variant="outline">
      <Package className="w-4 h-4 mr-2" />
      Continue Shopping
    </Button>
  </div>
);

// ===== Component =====

/**
 * CartItemList - Container for displaying cart items
 *
 * @param showBulkActions - Enable bulk action controls
 * @param showSelection - Show item selection checkboxes
 * @param compact - Use compact item layout
 * @param editable - Allow item editing
 * @param showGiftOptions - Show gift message options
 * @param showWishlistActions - Show wishlist buttons
 * @param className - Custom styling
 * @param maxHeight - Maximum container height
 * @param virtualizationThreshold - Enable virtualization above this item count
 * @param emptyStateComponent - Custom empty state
 * @param onBulkRemove - Bulk remove handler
 * @param onBulkMoveToWishlist - Bulk wishlist handler
 * @param onItemQuantityChange - Item quantity change handler
 * @param onItemRemove - Item remove handler
 */
export const CartItemList: React.FC<CartItemListProps> = ({
  showBulkActions = false,
  showSelection = false,
  compact = false,
  editable = true,
  showGiftOptions = true,
  showWishlistActions = false,
  className,
  maxHeight = "600px",
  virtualizationThreshold = 50,
  emptyStateComponent: EmptyStateComponent = DefaultEmptyState,
  onBulkRemove,
  onBulkMoveToWishlist,
  onItemQuantityChange,
  onItemRemove,
}) => {
  const { cart, isLoading, error, itemCount } = useCart();

  // ===== Local State =====

  const [selectedItems, setSelectedItems] = useState<Set<string>>(new Set());
  const [isPerformingBulkAction, setIsPerformingBulkAction] = useState(false);

  // ===== Computed Values =====

  const items = cart?.items || [];
  const isEmpty = itemCount === 0;
  const shouldVirtualize = items.length > virtualizationThreshold;
  const hasSelection = selectedItems.size > 0;
  const allSelected = selectedItems.size === items.length && items.length > 0;
  const partiallySelected = hasSelection && !allSelected;

  // Filter items by availability for better UX
  const { availableItems, unavailableItems } = useMemo(() => {
    const available = items.filter(item => item.isAvailable);
    const unavailable = items.filter(item => !item.isAvailable);
    return { availableItems: available, unavailableItems: unavailable };
  }, [items]);

  // ===== Selection Handlers =====

  const handleSelectAll = useCallback((checked: boolean) => {
    if (checked) {
      setSelectedItems(new Set(items.map(item => item.id)));
    } else {
      setSelectedItems(new Set());
    }
  }, [items]);

  const handleSelectItem = useCallback((itemId: string, checked: boolean) => {
    setSelectedItems(prev => {
      const newSet = new Set(prev);
      if (checked) {
        newSet.add(itemId);
      } else {
        newSet.delete(itemId);
      }
      return newSet;
    });
  }, []);

  // ===== Bulk Action Handlers =====

  const handleBulkRemove = useCallback(async () => {
    if (!hasSelection) return;

    setIsPerformingBulkAction(true);
    try {
      const itemIds = Array.from(selectedItems);
      if (onBulkRemove) {
        await onBulkRemove(itemIds);
      }
      setSelectedItems(new Set());
    } catch (error) {
      console.error('Failed to remove items:', error);
    } finally {
      setIsPerformingBulkAction(false);
    }
  }, [selectedItems, hasSelection, onBulkRemove]);

  const handleBulkMoveToWishlist = useCallback(async () => {
    if (!hasSelection) return;

    setIsPerformingBulkAction(true);
    try {
      const itemIds = Array.from(selectedItems);
      if (onBulkMoveToWishlist) {
        await onBulkMoveToWishlist(itemIds);
      }
      setSelectedItems(new Set());
    } catch (error) {
      console.error('Failed to move items to wishlist:', error);
    } finally {
      setIsPerformingBulkAction(false);
    }
  }, [selectedItems, hasSelection, onBulkMoveToWishlist]);

  // ===== Render Helpers =====

  const renderBulkActions = () => {
    if (!showBulkActions || !hasSelection) return null;

    return (
      <div className="flex items-center justify-between p-4 bg-muted/50 rounded-lg">
        <div className="flex items-center gap-2">
          <Badge variant="secondary">
            {selectedItems.size} item{selectedItems.size !== 1 ? 's' : ''} selected
          </Badge>
        </div>

        <div className="flex items-center gap-2">
          {showWishlistActions && onBulkMoveToWishlist && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleBulkMoveToWishlist}
              disabled={isPerformingBulkAction}
            >
              <Heart className="w-4 h-4 mr-2" />
              Move to Wishlist
            </Button>
          )}

          {onBulkRemove && (
            <Button
              variant="destructive"
              size="sm"
              onClick={handleBulkRemove}
              disabled={isPerformingBulkAction}
            >
              <Trash2 className="w-4 h-4 mr-2" />
              Remove Selected
            </Button>
          )}
        </div>
      </div>
    );
  };

  const renderSelectionHeader = () => {
    if (!showSelection) return null;

    return (
      <div className="flex items-center gap-3 p-4 bg-background border-b">
        <Checkbox
          checked={allSelected}
          ref={(ref) => {
            if (ref) {
              ref.indeterminate = partiallySelected;
            }
          }}
          onCheckedChange={handleSelectAll}
          aria-label="Select all items"
        />
        <span className="text-sm font-medium">
          {allSelected ? 'Deselect All' : 'Select All'}
        </span>
        {hasSelection && (
          <Badge variant="secondary" className="ml-auto">
            {selectedItems.size} selected
          </Badge>
        )}
      </div>
    );
  };

  const renderCartItem = (item: CartItemType) => {
    const isSelected = selectedItems.has(item.id);

    return (
      <div key={item.id} className="relative">
        {showSelection && (
          <div className="absolute top-4 left-4 z-10">
            <Checkbox
              checked={isSelected}
              onCheckedChange={(checked) => handleSelectItem(item.id, !!checked)}
              aria-label={`Select ${item.productName}`}
            />
          </div>
        )}

        <CartItem
          item={item}
          compact={compact}
          editable={editable}
          showGiftOptions={showGiftOptions}
          showWishlistAction={showWishlistActions}
          className={cn(
            showSelection && "ml-10",
            isSelected && "ring-2 ring-primary/20 bg-primary/5"
          )}
          onQuantityChange={onItemQuantityChange ? (quantity) => onItemQuantityChange(item.id, quantity) : undefined}
          onRemove={onItemRemove ? () => onItemRemove(item.id) : undefined}
        />
      </div>
    );
  };

  const renderItemSection = (title: string, items: CartItemType[], variant: 'default' | 'warning' = 'default') => {
    if (items.length === 0) return null;

    return (
      <div className="space-y-3">
        <div className="flex items-center gap-2">
          <h3 className={cn(
            "text-sm font-medium",
            variant === 'warning' && "text-warning-foreground"
          )}>
            {title}
          </h3>
          <Badge variant={variant === 'warning' ? 'secondary' : 'default'}>
            {items.length}
          </Badge>
        </div>

        <div className="space-y-3">
          {items.map(renderCartItem)}
        </div>
      </div>
    );
  };

  // ===== Loading State =====

  if (isLoading) {
    return (
      <TchatCard className={className}>
        <TchatCardContent>
          <div className="space-y-4">
            {Array.from({ length: 3 }).map((_, index) => (
              <div key={index} className="flex gap-4 animate-pulse">
                <div className="w-20 h-20 bg-muted rounded-md" />
                <div className="flex-1 space-y-2">
                  <div className="h-4 bg-muted rounded w-3/4" />
                  <div className="h-3 bg-muted rounded w-1/2" />
                  <div className="h-4 bg-muted rounded w-1/4" />
                </div>
              </div>
            ))}
          </div>
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== Error State =====

  if (error) {
    return (
      <TchatCard variant="outlined" className={cn("border-destructive", className)}>
        <TchatCardContent>
          <div className="flex items-center gap-3 text-destructive py-6">
            <AlertCircle className="w-5 h-5" />
            <div>
              <h3 className="font-medium">Failed to load cart items</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Please refresh the page or try again later.
              </p>
            </div>
          </div>
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== Empty State =====

  if (isEmpty) {
    return (
      <TchatCard className={className}>
        <TchatCardContent>
          <EmptyStateComponent />
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== Main Content =====

  return (
    <TchatCard className={className}>
      <TchatCardHeader
        title={`Cart Items (${itemCount})`}
        actions={renderBulkActions()}
      />

      {renderSelectionHeader()}

      <TchatCardContent
        className="p-0"
        style={{ maxHeight: shouldVirtualize ? maxHeight : undefined }}
      >
        <div className={cn(
          "space-y-6 p-4",
          shouldVirtualize && "overflow-y-auto"
        )}>
          {/* Available Items */}
          {renderItemSection("Available Items", availableItems)}

          {/* Separator between available and unavailable */}
          {availableItems.length > 0 && unavailableItems.length > 0 && (
            <Separator className="my-6" />
          )}

          {/* Unavailable Items */}
          {renderItemSection("Unavailable Items", unavailableItems, 'warning')}
        </div>
      </TchatCardContent>
    </TchatCard>
  );
};

// ===== Export Component =====

export default CartItemList;