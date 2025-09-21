import React, { useState } from 'react';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import { Separator } from './ui/separator';
import { Input } from './ui/input';
import { ScrollArea } from './ui/scroll-area';
import { 
  ArrowLeft, 
  Minus, 
  Plus, 
  Trash2, 
  ShoppingBag, 
  CreditCard,
  MapPin,
  Clock,
  Check
} from 'lucide-react';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";

interface CartItem {
  id: string;
  quantity: number;
  timestamp: number;
}

interface CartScreenProps {
  user: any;
  onBack: () => void;
  cartItems: CartItem[];
  onUpdateCart: (items: CartItem[]) => void;
}

// Mock product data
const mockProducts = {
  'prod-1': {
    id: 'prod-1',
    name: 'Authentic Pad Thai',
    price: 85,
    image: 'https://images.unsplash.com/photo-1559847844-5315695b6a77?w=200&h=200&fit=crop',
    vendor: 'Bangkok Street Kitchen',
    deliveryTime: '15-25 min'
  },
  'prod-2': {
    id: 'prod-2',
    name: 'Thai Green Curry',
    price: 120,
    image: 'https://images.unsplash.com/photo-1455619452474-d2be8b1e70cd?w=200&h=200&fit=crop',
    vendor: 'Spice Garden',
    deliveryTime: '20-30 min'
  },
  'prod-3': {
    id: 'prod-3',
    name: 'Mango Sticky Rice',
    price: 65,
    image: 'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=200&h=200&fit=crop',
    vendor: 'Thai Dessert House',
    deliveryTime: '10-15 min'
  }
};

export function CartScreen({ user, onBack, cartItems, onUpdateCart }: CartScreenProps) {
  const [promoCode, setPromoCode] = useState('');
  const [selectedAddress, setSelectedAddress] = useState('home');

  const updateQuantity = (itemId: string, newQuantity: number) => {
    if (newQuantity <= 0) {
      removeItem(itemId);
      return;
    }

    const updatedItems = cartItems.map(item => 
      item.id === itemId ? { ...item, quantity: newQuantity } : item
    );
    onUpdateCart(updatedItems);
  };

  const removeItem = (itemId: string) => {
    const updatedItems = cartItems.filter(item => item.id !== itemId);
    onUpdateCart(updatedItems);
    toast.success('Item removed from cart');
  };

  const clearCart = () => {
    onUpdateCart([]);
    toast.success('Cart cleared');
  };

  const applyPromoCode = () => {
    if (promoCode.toLowerCase() === 'sea2024') {
      toast.success('Promo code applied! 10% discount');
    } else {
      toast.error('Invalid promo code');
    }
  };

  const calculateSubtotal = () => {
    return cartItems.reduce((total, item) => {
      const product = mockProducts[item.id as keyof typeof mockProducts];
      return total + (product?.price || 0) * item.quantity;
    }, 0);
  };

  const subtotal = calculateSubtotal();
  const deliveryFee = subtotal > 200 ? 0 : 25;
  const discount = promoCode.toLowerCase() === 'sea2024' ? subtotal * 0.1 : 0;
  const total = subtotal + deliveryFee - discount;

  const handleCheckout = () => {
    if (cartItems.length === 0) {
      toast.error('Your cart is empty');
      return;
    }
    toast.success('Proceeding to payment...');
    // Navigate to payment flow
  };

  if (cartItems.length === 0) {
    return (
      <div className="h-full flex flex-col bg-background">
        {/* Header */}
        <header className="border-b bg-card px-4 py-3 flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={onBack}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="text-lg font-medium">Shopping Cart</h1>
        </header>

        {/* Empty State */}
        <div className="flex-1 flex items-center justify-center p-8">
          <div className="text-center">
            <ShoppingBag className="w-16 h-16 text-muted-foreground mx-auto mb-4" />
            <h2 className="text-lg font-medium mb-2">Your cart is empty</h2>
            <p className="text-muted-foreground mb-6">Add some delicious items from our store!</p>
            <Button onClick={onBack}>
              Continue Shopping
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col bg-background">
      {/* Header */}
      <header className="border-b bg-card px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={onBack}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="text-lg font-medium">Shopping Cart</h1>
          <Badge variant="secondary">
            {cartItems.length} {cartItems.length === 1 ? 'item' : 'items'}
          </Badge>
        </div>
        
        <Button 
          variant="ghost" 
          size="sm" 
          onClick={clearCart}
          className="text-destructive hover:text-destructive"
        >
          Clear All
        </Button>
      </header>

      <div className="flex-1 flex flex-col">
        {/* Cart Items */}
        <ScrollArea className="flex-1">
          <div className="p-4 space-y-4">
            {cartItems.map((item) => {
              const product = mockProducts[item.id as keyof typeof mockProducts];
              if (!product) return null;

              return (
                <Card key={item.id}>
                  <CardContent className="p-4">
                    <div className="flex gap-4">
                      <div className="w-16 h-16 rounded-lg overflow-hidden flex-shrink-0">
                        <ImageWithFallback
                          src={product.image}
                          alt={product.name}
                          className="w-full h-full object-cover"
                        />
                      </div>
                      
                      <div className="flex-1 min-w-0">
                        <h3 className="font-medium line-clamp-1">{product.name}</h3>
                        <p className="text-sm text-muted-foreground">{product.vendor}</p>
                        <div className="flex items-center gap-2 mt-1">
                          <Clock className="w-3 h-3 text-muted-foreground" />
                          <span className="text-xs text-muted-foreground">{product.deliveryTime}</span>
                        </div>
                        <div className="flex items-center justify-between mt-3">
                          <span className="font-medium">฿{product.price}</span>
                          <div className="flex items-center gap-2">
                            <Button
                              variant="outline"
                              size="icon"
                              className="w-8 h-8"
                              onClick={() => updateQuantity(item.id, item.quantity - 1)}
                            >
                              <Minus className="w-3 h-3" />
                            </Button>
                            <span className="w-8 text-center">{item.quantity}</span>
                            <Button
                              variant="outline"
                              size="icon"
                              className="w-8 h-8"
                              onClick={() => updateQuantity(item.id, item.quantity + 1)}
                            >
                              <Plus className="w-3 h-3" />
                            </Button>
                          </div>
                        </div>
                      </div>
                      
                      <Button
                        variant="ghost"
                        size="icon"
                        className="text-destructive hover:text-destructive"
                        onClick={() => removeItem(item.id)}
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </ScrollArea>

        {/* Order Summary */}
        <div className="border-t bg-card p-4 space-y-4">
          {/* Delivery Address */}
          <div className="flex items-center gap-3 p-3 bg-muted rounded-lg">
            <MapPin className="w-5 h-5 text-primary" />
            <div className="flex-1">
              <p className="text-sm font-medium">Delivery to Home</p>
              <p className="text-xs text-muted-foreground">123 Sukhumvit Road, Bangkok</p>
            </div>
            <Button variant="ghost" size="sm">
              Change
            </Button>
          </div>

          {/* Promo Code */}
          <div className="flex gap-2">
            <Input
              placeholder="Enter promo code"
              value={promoCode}
              onChange={(e) => setPromoCode(e.target.value)}
              className="flex-1"
            />
            <Button variant="outline" onClick={applyPromoCode}>
              Apply
            </Button>
          </div>

          <Separator />

          {/* Pricing Breakdown */}
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span>Subtotal</span>
              <span>฿{subtotal}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span>Delivery Fee</span>
              <span className={deliveryFee === 0 ? 'text-green-600 line-through' : ''}>
                ฿{deliveryFee}
              </span>
            </div>
            {discount > 0 && (
              <div className="flex justify-between text-sm text-green-600">
                <span>Discount (SEA2024)</span>
                <span>-฿{discount.toFixed(0)}</span>
              </div>
            )}
            <Separator />
            <div className="flex justify-between font-medium">
              <span>Total</span>
              <span>฿{total.toFixed(0)}</span>
            </div>
          </div>

          {/* Checkout Button */}
          <Button 
            className="w-full"
            size="lg"
            onClick={handleCheckout}
          >
            <CreditCard className="w-5 h-5 mr-2" />
            Checkout ฿{total.toFixed(0)}
          </Button>

          {/* Free Delivery Message */}
          {subtotal < 200 && (
            <p className="text-xs text-center text-muted-foreground">
              Add ฿{200 - subtotal} more for free delivery
            </p>
          )}
        </div>
      </div>
    </div>
  );
}