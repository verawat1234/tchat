// T045 - ProductMessage component with commerce integration
/**
 * ProductMessage Component
 * Displays product cards with pricing, availability, cart actions, and purchase workflows
 * Supports product variants, bulk discounts, and wishlist management
 */

import React, { useState, useCallback, useMemo, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { ProductContent, ProductVariant, ProductAvailability, DiscountType } from '../../types/ProductContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Separator } from '../ui/separator';
import {
  ShoppingCart,
  Heart,
  Star,
  Package,
  Truck,
  Shield,
  Plus,
  Minus,
  ExternalLink,
  Share,
  Eye,
  AlertCircle,
  Tag,
  Zap,
  Gift,
  Clock
} from 'lucide-react';

// Component Props Interface
interface ProductMessageProps {
  message: MessageData & { content: ProductContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onAddToCart?: (productId: string, variantId?: string, quantity?: number) => void;
  onAddToWishlist?: (productId: string) => void;
  onViewProduct?: (productId: string) => void;
  onShare?: (productId: string) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  showQuickActions?: boolean;
  performanceMode?: boolean;
}

// Animation Variants
const productVariants = {
  initial: { opacity: 0, scale: 0.95, y: 20 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.95, y: -20 }
};

const imageVariants = {
  initial: { opacity: 0, scale: 1.1 },
  animate: { opacity: 1, scale: 1 },
  exit: { opacity: 0, scale: 0.9 }
};

const priceVariants = {
  initial: { opacity: 0, x: -20 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 20 }
};

const actionVariants = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -20 }
};

export const ProductMessage: React.FC<ProductMessageProps> = ({
  message,
  onInteraction,
  onAddToCart,
  onAddToWishlist,
  onViewProduct,
  onShare,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  showQuickActions = true,
  performanceMode = false
}) => {
  const productRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Product state
  const [selectedVariant, setSelectedVariant] = useState<string>(
    content.variants?.[0]?.id || ''
  );
  const [quantity, setQuantity] = useState(1);
  const [isWishlisted, setIsWishlisted] = useState(false);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  // Get selected variant or default product
  const currentVariant = useMemo(() => {
    if (content.variants && selectedVariant) {
      return content.variants.find(v => v.id === selectedVariant) || content.variants[0];
    }
    return null;
  }, [content.variants, selectedVariant]);

  // Calculate pricing
  const pricing = useMemo(() => {
    const basePrice = currentVariant?.price || content.basePrice;
    const originalPrice = currentVariant?.originalPrice || content.originalPrice;
    const discount = originalPrice && basePrice < originalPrice
      ? ((originalPrice - basePrice) / originalPrice) * 100
      : 0;

    // Apply bulk discount if applicable
    let finalPrice = basePrice;
    let bulkDiscountApplied = false;

    if (content.bulkDiscounts && quantity > 1) {
      const applicableDiscount = content.bulkDiscounts
        .filter(d => quantity >= d.minQuantity)
        .sort((a, b) => b.minQuantity - a.minQuantity)[0];

      if (applicableDiscount) {
        if (applicableDiscount.discountType === DiscountType.PERCENTAGE) {
          finalPrice = basePrice * (1 - applicableDiscount.discountValue / 100);
        } else {
          finalPrice = Math.max(0, basePrice - applicableDiscount.discountValue);
        }
        bulkDiscountApplied = true;
      }
    }

    const total = finalPrice * quantity;
    const savings = (basePrice - finalPrice) * quantity;

    return {
      basePrice,
      originalPrice,
      finalPrice,
      total,
      discount,
      savings,
      bulkDiscountApplied
    };
  }, [currentVariant, content.basePrice, content.originalPrice, content.bulkDiscounts, quantity]);

  // Check availability
  const availability = useMemo(() => {
    const stock = currentVariant?.stock ?? content.stock;
    const isLowStock = typeof stock === 'number' && stock <= 5 && stock > 0;
    const isOutOfStock = stock === 0;
    const isUnlimited = stock === null || stock === undefined;

    return {
      stock,
      isLowStock,
      isOutOfStock,
      isUnlimited,
      canPurchase: !isOutOfStock && !readonly
    };
  }, [currentVariant?.stock, content.stock, readonly]);

  // Format price
  const formatPrice = useCallback((price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: content.currency || 'USD'
    }).format(price);
  }, [content.currency]);

  // Handle variant selection
  const handleVariantChange = useCallback((variantId: string) => {
    setSelectedVariant(variantId);
    setCurrentImageIndex(0);

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'product_variant_select',
        data: { productId: content.id, variantId },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [message.id, content.id, onInteraction]);

  // Handle quantity change
  const handleQuantityChange = useCallback((newQuantity: number) => {
    const maxQuantity = availability.stock || 99;
    const validQuantity = Math.max(1, Math.min(newQuantity, maxQuantity));
    setQuantity(validQuantity);
  }, [availability.stock]);

  // Handle add to cart
  const handleAddToCart = useCallback(() => {
    if (!availability.canPurchase) return;

    if (onAddToCart) {
      onAddToCart(content.id, selectedVariant, quantity);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'add_to_cart',
        data: {
          productId: content.id,
          variantId: selectedVariant,
          quantity,
          price: pricing.finalPrice
        },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [
    availability.canPurchase,
    onAddToCart,
    content.id,
    selectedVariant,
    quantity,
    onInteraction,
    message.id,
    pricing.finalPrice
  ]);

  // Handle wishlist toggle
  const handleWishlistToggle = useCallback(() => {
    if (readonly) return;

    setIsWishlisted(!isWishlisted);

    if (onAddToWishlist) {
      onAddToWishlist(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: isWishlisted ? 'remove_from_wishlist' : 'add_to_wishlist',
        data: { productId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, isWishlisted, onAddToWishlist, content.id, onInteraction, message.id]);

  // Handle product view
  const handleViewProduct = useCallback(() => {
    if (onViewProduct) {
      onViewProduct(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'view_product',
        data: { productId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [onViewProduct, content.id, onInteraction, message.id]);

  // Handle share
  const handleShare = useCallback(() => {
    if (onShare) {
      onShare(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'share_product',
        data: { productId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [onShare, content.id, onInteraction, message.id]);

  // Get availability badge
  const getAvailabilityBadge = useCallback(() => {
    if (availability.isOutOfStock) {
      return (
        <Badge variant="destructive" className="text-xs">
          <Package className="w-3 h-3 mr-1" />
          Out of stock
        </Badge>
      );
    }

    if (availability.isLowStock) {
      return (
        <Badge variant="outline" className="text-xs border-orange-500 text-orange-600">
          <AlertCircle className="w-3 h-3 mr-1" />
          Low stock ({availability.stock} left)
        </Badge>
      );
    }

    return (
      <Badge variant="outline" className="text-xs border-green-500 text-green-600">
        <Package className="w-3 h-3 mr-1" />
        In stock
      </Badge>
    );
  }, [availability]);

  // Performance optimization: skip animation in performance mode
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: productVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={productRef}
        className={cn(
          "product-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`product-message-${message.id}`}
        data-product-id={content.id}
        role="article"
        aria-label={`Product: ${content.name} - ${formatPrice(pricing.finalPrice)}`}
      >
        <Card className="product-card overflow-hidden">
          {/* Header */}
          <CardHeader className="space-y-3">
            <div className="flex items-start justify-between gap-3">
              <div className="flex items-center gap-3 min-w-0 flex-1">
                {showAvatar && (
                  <motion.div
                    initial={performanceMode ? {} : { scale: 0.8, opacity: 0 }}
                    animate={performanceMode ? {} : { scale: 1, opacity: 1 }}
                    transition={{ delay: 0.1 }}
                  >
                    <Avatar className={cn(compactMode ? "w-8 h-8" : "w-10 h-10")}>
                      <AvatarImage src={`/avatars/${message.senderName.toLowerCase()}.png`} />
                      <AvatarFallback>
                        {message.senderName.substring(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                  </motion.div>
                )}

                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="font-semibold text-foreground truncate">
                      {message.senderName}
                    </span>
                    {message.isOwn && (
                      <Badge variant="secondary" className="text-xs">You</Badge>
                    )}
                  </div>
                  {showTimestamp && (
                    <p className="text-xs text-muted-foreground mt-1">
                      Shared a product • {message.timestamp.toLocaleDateString()}
                    </p>
                  )}
                </div>
              </div>

              {showQuickActions && (
                <div className="flex items-center gap-1">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleWishlistToggle}
                        className={cn(
                          "h-8 w-8 p-0",
                          isWishlisted && "text-red-500 hover:text-red-600"
                        )}
                      >
                        <Heart className={cn("w-4 h-4", isWishlisted && "fill-current")} />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      {isWishlisted ? 'Remove from wishlist' : 'Add to wishlist'}
                    </TooltipContent>
                  </Tooltip>

                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleShare}
                        className="h-8 w-8 p-0"
                      >
                        <Share className="w-4 h-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Share product</TooltipContent>
                  </Tooltip>
                </div>
              )}
            </div>
          </CardHeader>

          <CardContent className="space-y-4">
            {/* Product image */}
            <div className="relative aspect-square rounded-lg overflow-hidden bg-muted">
              <motion.img
                key={currentImageIndex}
                variants={performanceMode ? {} : imageVariants}
                initial={performanceMode ? {} : "initial"}
                animate={performanceMode ? {} : "animate"}
                exit={performanceMode ? {} : "exit"}
                src={content.images[currentImageIndex] || '/placeholder-product.png'}
                alt={content.name}
                className="w-full h-full object-cover cursor-pointer hover:scale-105 transition-transform duration-200"
                onClick={handleViewProduct}
              />

              {/* Image navigation */}
              {content.images.length > 1 && (
                <div className="absolute bottom-2 left-1/2 transform -translate-x-1/2 flex gap-1">
                  {content.images.map((_, index) => (
                    <button
                      key={index}
                      onClick={() => setCurrentImageIndex(index)}
                      className={cn(
                        "w-2 h-2 rounded-full transition-colors",
                        index === currentImageIndex
                          ? "bg-white"
                          : "bg-white/50 hover:bg-white/75"
                      )}
                      aria-label={`View image ${index + 1}`}
                    />
                  ))}
                </div>
              )}

              {/* Badges */}
              <div className="absolute top-2 left-2 flex flex-col gap-1">
                {pricing.discount > 0 && (
                  <Badge className="bg-red-500 text-white text-xs">
                    <Tag className="w-3 h-3 mr-1" />
                    {Math.round(pricing.discount)}% OFF
                  </Badge>
                )}

                {content.isFeatured && (
                  <Badge className="bg-yellow-500 text-white text-xs">
                    <Star className="w-3 h-3 mr-1" />
                    Featured
                  </Badge>
                )}

                {pricing.bulkDiscountApplied && (
                  <Badge className="bg-blue-500 text-white text-xs">
                    <Gift className="w-3 h-3 mr-1" />
                    Bulk discount
                  </Badge>
                )}
              </div>

              {/* Availability badge */}
              <div className="absolute top-2 right-2">
                {getAvailabilityBadge()}
              </div>
            </div>

            {/* Product details */}
            <div className="space-y-3">
              {/* Name and rating */}
              <div className="space-y-1">
                <h3 className="font-semibold text-lg text-foreground leading-tight">
                  {content.name}
                </h3>

                {content.rating && (
                  <div className="flex items-center gap-2 text-sm">
                    <div className="flex items-center gap-1">
                      <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                      <span className="font-medium">{content.rating.average.toFixed(1)}</span>
                    </div>
                    <span className="text-muted-foreground">
                      ({content.rating.count} reviews)
                    </span>
                  </div>
                )}
              </div>

              {/* Description */}
              {content.description && (
                <p className="text-sm text-muted-foreground line-clamp-2">
                  {content.description}
                </p>
              )}

              {/* Variants */}
              {content.variants && content.variants.length > 1 && (
                <div className="space-y-2">
                  <Label className="text-sm font-medium">Variant</Label>
                  <Select value={selectedVariant} onValueChange={handleVariantChange}>
                    <SelectTrigger className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {content.variants.map((variant) => (
                        <SelectItem key={variant.id} value={variant.id}>
                          {variant.name} - {formatPrice(variant.price)}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )}

              {/* Pricing */}
              <motion.div
                variants={performanceMode ? {} : priceVariants}
                initial={performanceMode ? {} : "initial"}
                animate={performanceMode ? {} : "animate"}
                className="space-y-2"
              >
                <div className="flex items-center gap-3">
                  <span className="text-2xl font-bold text-foreground">
                    {formatPrice(pricing.finalPrice)}
                  </span>

                  {pricing.originalPrice && pricing.discount > 0 && (
                    <span className="text-lg text-muted-foreground line-through">
                      {formatPrice(pricing.originalPrice)}
                    </span>
                  )}

                  {pricing.savings > 0 && quantity > 1 && (
                    <Badge variant="outline" className="text-xs text-green-600 border-green-500">
                      <Zap className="w-3 h-3 mr-1" />
                      Save {formatPrice(pricing.savings)}
                    </Badge>
                  )}
                </div>

                {quantity > 1 && (
                  <div className="text-sm text-muted-foreground">
                    Total: <span className="font-semibold">{formatPrice(pricing.total)}</span>
                    {pricing.bulkDiscountApplied && (
                      <span className="text-green-600 ml-2">
                        (Bulk discount applied!)
                      </span>
                    )}
                  </div>
                )}
              </motion.div>

              <Separator />

              {/* Actions */}
              <motion.div
                variants={performanceMode ? {} : actionVariants}
                initial={performanceMode ? {} : "initial"}
                animate={performanceMode ? {} : "animate"}
                transition={{ delay: 0.2 }}
                className="space-y-3"
              >
                {/* Quantity selector */}
                {!readonly && availability.canPurchase && (
                  <div className="flex items-center justify-between">
                    <Label className="text-sm font-medium">Quantity</Label>
                    <div className="flex items-center gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleQuantityChange(quantity - 1)}
                        disabled={quantity <= 1}
                        className="h-8 w-8 p-0"
                      >
                        <Minus className="w-3 h-3" />
                      </Button>

                      <span className="w-12 text-center text-sm font-medium">
                        {quantity}
                      </span>

                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleQuantityChange(quantity + 1)}
                        disabled={availability.stock !== null && quantity >= availability.stock}
                        className="h-8 w-8 p-0"
                      >
                        <Plus className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                )}

                {/* Action buttons */}
                <div className="flex gap-2">
                  <Button
                    onClick={handleAddToCart}
                    disabled={!availability.canPurchase}
                    className="flex-1 gap-2"
                    size="sm"
                  >
                    <ShoppingCart className="w-4 h-4" />
                    Add to Cart
                  </Button>

                  <Button
                    variant="outline"
                    onClick={handleViewProduct}
                    className="gap-2"
                    size="sm"
                  >
                    <Eye className="w-4 h-4" />
                    View Details
                  </Button>
                </div>
              </motion.div>

              {/* Features */}
              {content.features && content.features.length > 0 && (
                <div className="space-y-2">
                  <div className="flex flex-wrap gap-2">
                    {content.features.slice(0, 3).map((feature, index) => (
                      <div key={index} className="flex items-center gap-1 text-xs text-muted-foreground">
                        <Shield className="w-3 h-3" />
                        <span>{feature}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Performance Debug Info (Development Only) */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            V: {selectedVariant || 'default'} | Q: {quantity} | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedProductMessage = React.memo(ProductMessage, (prevProps, nextProps) => {
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.showQuickActions === nextProps.showQuickActions &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedProductMessage.displayName = 'MemoizedProductMessage';

export default ProductMessage;