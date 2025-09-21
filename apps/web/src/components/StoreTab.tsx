import React, { useState } from 'react';
import { Search, Filter, ShoppingCart, Star, MapPin, Clock, Zap, Heart, QrCode, Wallet, CreditCard, TrendingUp, Compass, Grid3X3, Store, Users, Eye } from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './ui/dialog';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";

interface StoreTabProps {
  user: any;
  onShopClick?: (shopId: string) => void;
  onProductClick?: (productId: string) => void;
  onAddToCart?: (productId: string, quantity?: number) => void;
  onLiveStreamClick?: (streamId: string) => void;
  onProductShare?: (productId: string, productData: any) => void;
  onShopShare?: (shopId: string, shopData: any) => void;
  cartItems?: string[];
}

interface Product {
  id: string;
  name: string;
  price: number;
  currency: string;
  image: string;
  merchant: string;
  rating: number;
  deliveryTime: string;
  distance: string;
  category: string;
  isLive?: boolean;
  discount?: number;
  isHot?: boolean;
  orders?: number;
}

interface Shop {
  id: string;
  name: string;
  description: string;
  image: string;
  avatar: string;
  rating: number;
  deliveryTime: string;
  distance: string;
  category: string;
  isVerified?: boolean;
  followers: number;
  products: number;
  isHot?: boolean;
  specialOffer?: string;
}

export function StoreTab({ user, onShopClick, onProductClick, onAddToCart, onLiveStreamClick, onProductShare, onShopShare, cartItems = [] }: StoreTabProps) {
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [selectedFilter, setSelectedFilter] = useState('all');
  const [showPayment, setShowPayment] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);

  // Filter options
  const filters = [
    { id: 'all', name: 'All', icon: Grid3X3 },
    { id: 'nearby', name: 'Nearby', icon: Compass },
    { id: 'hot', name: 'Hot', icon: TrendingUp },
    { id: 'live', name: 'Live', icon: Zap }
  ];

  // Categories
  const categories = [
    { id: 'all', name: 'All', icon: '🏪' },
    { id: 'food', name: 'Food', icon: '🍜' },
    { id: 'grocery', name: 'Grocery', icon: '🛒' },
    { id: 'fashion', name: 'Fashion', icon: '👕' },
    { id: 'electronics', name: 'Electronics', icon: '📱' },
    { id: 'beauty', name: 'Beauty', icon: '💄' },
    { id: 'home', name: 'Home', icon: '🏠' },
    { id: 'books', name: 'Books', icon: '📚' }
  ];

  // Mock shops data
  const shops: Shop[] = [
    {
      id: 'shop1',
      name: 'Bangkok Street Food Palace',
      description: 'Authentic Thai street food made fresh daily',
      image: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      avatar: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=150&h=150&fit=crop&crop=face',
      rating: 4.8,
      deliveryTime: '15-20 min',
      distance: '0.5 km',
      category: 'food',
      isVerified: true,
      followers: 2840,
      products: 45,
      isHot: true,
      specialOffer: '20% off first order'
    },
    {
      id: 'shop2',
      name: 'Thai Fashion Boutique',
      description: 'Modern Thai fashion and traditional clothing',
      image: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwZmFzaGlvbiUyMGJvdXRpcXVlfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face',
      rating: 4.6,
      deliveryTime: '1-2 days',
      distance: '1.2 km',
      category: 'fashion',
      isVerified: true,
      followers: 1560,
      products: 128,
      specialOffer: 'Free shipping over ฿500'
    },
    {
      id: 'shop3',
      name: 'Fresh Market 24/7',
      description: 'Fresh groceries and daily essentials',
      image: 'https://images.unsplash.com/photo-1488459716781-31db52582fe9?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwZ3JvY2VyeSUyMG1hcmtldHxlbnwxfHx8fDE3NTgzOTQ1MTV8MA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face',
      rating: 4.4,
      deliveryTime: '30-45 min',
      distance: '0.8 km',
      category: 'grocery',
      followers: 892,
      products: 234,
      isHot: true
    },
    {
      id: 'shop4',
      name: 'Tech Zone Thailand',
      description: 'Latest electronics and gadgets',
      image: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0ZWNoJTIwc3RvcmUlMjB0aGFpbGFuZHxlbnwxfHx8fDE3NTgzOTQ1MTV8MA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
      rating: 4.7,
      deliveryTime: '1-3 days',
      distance: '2.1 km',
      category: 'electronics',
      isVerified: true,
      followers: 3420,
      products: 89
    }
  ];

  // Mock products data
  const products: Product[] = [
    {
      id: '1',
      name: 'Pad Thai Goong',
      price: 45,
      currency: 'THB',
      image: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      merchant: 'Bangkok Street Food Palace',
      rating: 4.8,
      deliveryTime: '15-20 min',
      distance: '0.5 km',
      category: 'food',
      isLive: true,
      isHot: true,
      orders: 1234
    },
    {
      id: '2',
      name: 'Som Tam Thai',
      price: 35,
      currency: 'THB',
      image: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      merchant: 'Bangkok Street Food Palace',
      rating: 4.6,
      deliveryTime: '15-20 min',
      distance: '0.5 km',
      category: 'food',
      discount: 20,
      isHot: true,
      orders: 892
    },
    {
      id: '3',
      name: 'Mango Sticky Rice',
      price: 40,
      currency: 'THB',
      image: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwZGVzc2VydCUyMG1hbmdvfGVufDF8fHx8MTc1ODM5NDUxN3ww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      merchant: 'Bangkok Street Food Palace',
      rating: 4.9,
      deliveryTime: '15-20 min',
      distance: '0.5 km',
      category: 'food',
      orders: 567
    },
    {
      id: '4',
      name: 'Thai Silk Dress',
      price: 890,
      currency: 'THB',
      image: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc2lsayUyMGRyZXNzfGVufDF8fHx8MTc1ODM5NDUxN3ww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      merchant: 'Thai Fashion Boutique',
      rating: 4.7,
      deliveryTime: '1-2 days',
      distance: '1.2 km',
      category: 'fashion',
      discount: 15,
      orders: 234
    }
  ];

  const handleAddToCart = (productId: string) => {
    if (onAddToCart) {
      onAddToCart(productId, 1);
    } else {
      toast.success('Added to cart!');
    }
  };

  const handleBuyNow = (product: Product) => {
    setSelectedProduct(product);
    setShowPayment(true);
  };

  const handlePayment = (method: string) => {
    toast.success(`Payment successful via ${method}!`);
    setShowPayment(false);
    setSelectedProduct(null);
  };

  const filteredShops = shops.filter(shop => {
    if (selectedCategory !== 'all' && shop.category !== selectedCategory) return false;
    if (selectedFilter === 'nearby' && parseFloat(shop.distance) > 1) return false;
    if (selectedFilter === 'hot' && !shop.isHot) return false;
    return true;
  });

  const filteredProducts = products.filter(product => {
    if (selectedCategory !== 'all' && product.category !== selectedCategory) return false;
    if (selectedFilter === 'nearby' && parseFloat(product.distance) > 1) return false;
    if (selectedFilter === 'hot' && !product.isHot) return false;
    if (selectedFilter === 'live' && !product.isLive) return false;
    return true;
  });

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="sticky top-0 z-30 border-b border-border p-3 sm:p-4 bg-card/95 backdrop-blur-sm">
        <div className="flex items-center gap-3 sm:gap-4 mb-3 sm:mb-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 w-4 h-4 text-muted-foreground" />
            <Input placeholder="Search shops, products..." className="pl-10" />
          </div>
          <Button variant="outline" size="icon" className="h-10 w-10 touch-manipulation">
            <Filter className="w-5 h-5" />
          </Button>
          <Button variant="outline" size="icon" className="relative h-10 w-10 touch-manipulation">
            <ShoppingCart className="w-5 h-5" />
            {cartItems.length > 0 && (
              <Badge className="absolute -top-2 -right-2 w-5 h-5 text-xs p-0 flex items-center justify-center">
                {cartItems.length}
              </Badge>
            )}
          </Button>
        </div>

        {/* Filter Tags - Icon Only */}
        <div className="flex gap-1.5 sm:gap-2 mb-3 sm:mb-4 overflow-x-auto scrollbar-hide pb-1">
          {filters.map((filter) => {
            const IconComponent = filter.icon;
            return (
              <Button
                key={filter.id}
                variant={selectedFilter === filter.id ? 'default' : 'outline'}
                size="sm"
                onClick={() => setSelectedFilter(filter.id)}
                className="h-8 sm:h-9 w-8 sm:w-9 p-0 flex-shrink-0 rounded-full touch-manipulation hover:scale-105 transition-transform"
                title={filter.name}
              >
                <IconComponent className="w-4 h-4" />
              </Button>
            );
          })}
        </div>

        {/* Categories - Text Labels with Horizontal Scroll */}
        <div className="w-full overflow-x-auto scrollbar-hide">
          <div className="flex gap-1.5 sm:gap-2 pb-2 min-w-max">
            {categories.map((category) => (
              <Button
                key={category.id}
                variant={selectedCategory === category.id ? 'default' : 'outline'}
                size="sm"
                onClick={() => setSelectedCategory(category.id)}
                className="h-8 sm:h-9 px-3 sm:px-4 flex-shrink-0 whitespace-nowrap touch-manipulation hover:scale-105 transition-transform"
              >
                <span className="text-base sm:text-lg mr-1.5">{category.icon}</span>
                <span className="text-sm">{category.name}</span>
              </Button>
            ))}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-hidden">
        <Tabs defaultValue="shops" className="h-full flex flex-col">
          <div className="mx-3 sm:mx-4 mt-3 sm:mt-4 flex justify-center">
            <TabsList className="inline-flex h-auto p-0.5 sm:p-1 bg-muted rounded-full">
              <TabsTrigger 
                value="shops" 
                className="h-9 sm:h-10 px-3 sm:px-4 py-2 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-2"
              >
                <Store className="w-4 h-4 sm:w-5 sm:h-5" />
                <span className="text-sm">Shops</span>
              </TabsTrigger>
              <TabsTrigger 
                value="products" 
                className="h-9 sm:h-10 px-3 sm:px-4 py-2 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-2"
              >
                <Grid3X3 className="w-4 h-4 sm:w-5 sm:h-5" />
                <span className="text-sm">Products</span>
              </TabsTrigger>
              <TabsTrigger 
                value="live" 
                className="h-9 sm:h-10 px-3 sm:px-4 py-2 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-2"
              >
                <Eye className="w-4 h-4 sm:w-5 sm:h-5" />
                <span className="text-sm">Live</span>
              </TabsTrigger>
            </TabsList>
          </div>

          <TabsContent value="shops" className="flex-1 overflow-hidden mt-4">
            <ScrollArea className="h-full px-4">
              <div className="space-y-4 pb-4">
                {filteredShops.map((shop) => (
                  <Card 
                    key={shop.id} 
                    className="overflow-hidden hover:shadow-md transition-shadow cursor-pointer"
                    onClick={() => onShopClick?.(shop.id)}
                  >
                    <div className="relative">
                      <ImageWithFallback
                        src={shop.image}
                        alt={shop.name}
                        className="w-full h-32 object-cover"
                      />
                      {shop.isHot && (
                        <Badge className="absolute top-2 left-2 bg-red-500">
                          <TrendingUp className="w-3 h-3 mr-1" />
                          HOT
                        </Badge>
                      )}
                      {shop.specialOffer && (
                        <Badge className="absolute top-2 right-2 bg-green-500">
                          {shop.specialOffer}
                        </Badge>
                      )}
                    </div>
                    
                    <CardContent className="p-4">
                      <div className="flex items-start gap-3">
                        <Avatar className="w-12 h-12">
                          <AvatarImage src={shop.avatar} />
                          <AvatarFallback>
                            <Store className="w-6 h-6" />
                          </AvatarFallback>
                        </Avatar>
                        
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <h3 className="font-medium">{shop.name}</h3>
                            {shop.isVerified && (
                              <Star className="w-4 h-4 text-chart-1 fill-chart-1" />
                            )}
                          </div>
                          <p className="text-sm text-muted-foreground mb-2">{shop.description}</p>
                          
                          <div className="flex items-center gap-4 mb-2">
                            <div className="flex items-center gap-1">
                              <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                              <span className="text-sm">{shop.rating}</span>
                            </div>
                            <div className="flex items-center gap-1">
                              <Clock className="w-4 h-4 text-muted-foreground" />
                              <span className="text-sm text-muted-foreground">{shop.deliveryTime}</span>
                            </div>
                            <div className="flex items-center gap-1">
                              <MapPin className="w-4 h-4 text-muted-foreground" />
                              <span className="text-sm text-muted-foreground">{shop.distance}</span>
                            </div>
                          </div>
                          
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-4 text-sm text-muted-foreground">
                              <div className="flex items-center gap-1">
                                <Users className="w-4 h-4" />
                                <span>{shop.followers.toLocaleString()} followers</span>
                              </div>
                              <div className="flex items-center gap-1">
                                <Store className="w-4 h-4" />
                                <span>{shop.products} products</span>
                              </div>
                            </div>
                            
                            <Button 
                              size="sm" 
                              onClick={(e) => {
                                e.stopPropagation();
                                onShopClick?.(shop.id);
                              }}
                            >
                              Visit Shop
                            </Button>
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </ScrollArea>
          </TabsContent>

          <TabsContent value="products" className="flex-1 overflow-hidden mt-4">
            <ScrollArea className="h-full px-4">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pb-4">
                {filteredProducts.map((product) => (
                  <Card key={product.id} className="overflow-hidden">
                    <div className="relative cursor-pointer" onClick={() => onProductClick?.(product.id)}>
                      <ImageWithFallback
                        src={product.image}
                        alt={product.name}
                        className="w-full h-32 object-cover"
                      />
                      {product.isLive && (
                        <Badge className="absolute top-2 left-2 bg-red-500">
                          <Zap className="w-3 h-3 mr-1" />
                          LIVE
                        </Badge>
                      )}
                      {product.isHot && (
                        <Badge className="absolute top-2 right-2 bg-orange-500">
                          <TrendingUp className="w-3 h-3 mr-1" />
                          HOT
                        </Badge>
                      )}
                      {product.discount && (
                        <Badge className="absolute bottom-2 right-2 bg-green-500">
                          -{product.discount}%
                        </Badge>
                      )}
                      <Button
                        variant="ghost"
                        size="icon"
                        className="absolute top-2 right-2 bg-white/80 hover:bg-white"
                      >
                        <Heart className="w-4 h-4" />
                      </Button>
                    </div>
                    
                    <CardContent className="p-4">
                      <h3 
                        className="font-medium mb-1 cursor-pointer hover:text-primary" 
                        onClick={() => onProductClick?.(product.id)}
                      >
                        {product.name}
                      </h3>
                      <p className="text-sm text-muted-foreground mb-2">{product.merchant}</p>
                      
                      <div className="flex items-center gap-2 mb-2">
                        <div className="flex items-center gap-1">
                          <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                          <span className="text-sm">{product.rating}</span>
                        </div>
                        <span className="text-sm text-muted-foreground">•</span>
                        <div className="flex items-center gap-1">
                          <Clock className="w-4 h-4 text-muted-foreground" />
                          <span className="text-sm text-muted-foreground">{product.deliveryTime}</span>
                        </div>
                        {product.orders && (
                          <>
                            <span className="text-sm text-muted-foreground">•</span>
                            <div className="flex items-center gap-1">
                              <Eye className="w-4 h-4 text-muted-foreground" />
                              <span className="text-sm text-muted-foreground">{product.orders} sold</span>
                            </div>
                          </>
                        )}
                      </div>
                      
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="font-bold text-lg">
                            ฿{product.price}
                          </span>
                          {product.discount && (
                            <span className="text-sm text-muted-foreground line-through">
                              ฿{Math.round(product.price / (1 - product.discount / 100))}
                            </span>
                          )}
                        </div>
                        
                        <div className="flex gap-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleAddToCart(product.id)}
                          >
                            Add to Cart
                          </Button>
                          <Button size="sm" onClick={() => handleBuyNow(product)}>
                            Buy Now
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </ScrollArea>
          </TabsContent>

          <TabsContent value="live" className="flex-1 overflow-hidden mt-4">
            <ScrollArea className="h-full px-4">
              {/* Featured Live Stream */}
              <Card 
                className="bg-gradient-to-r from-red-500 to-pink-500 text-white mb-4 cursor-pointer hover:shadow-lg transition-shadow"
                onClick={() => onLiveStreamClick?.('featured-stream')}
              >
                <CardContent className="p-6">
                  <div className="flex items-center gap-2 mb-2">
                    <div className="w-3 h-3 bg-white rounded-full animate-pulse"></div>
                    <span className="font-medium">LIVE NOW</span>
                  </div>
                  <h3 className="text-xl font-bold mb-2">Thai Street Food Festival 🍜</h3>
                  <p className="mb-4">Watch Chef Somsak cook authentic Pad Thai & Tom Yum live with exclusive discounts!</p>
                  <div className="flex items-center gap-4">
                    <span>👥 1,247 watching</span>
                    <span>⏰ Started 23 min ago</span>
                    <span>🔥 25% off live orders</span>
                  </div>
                </CardContent>
              </Card>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pb-4">
                {/* Live Stream 1 */}
                <Card className="overflow-hidden hover:shadow-md transition-shadow">
                  <div className="relative">
                    <ImageWithFallback
                      src="https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral"
                      alt="Bangkok Food Market Live"
                      className="w-full h-32 object-cover"
                    />
                    <Badge className="absolute top-2 left-2 bg-red-500">
                      <div className="w-2 h-2 bg-white rounded-full mr-1 animate-pulse"></div>
                      LIVE
                    </Badge>
                    <div className="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
                      856 watching
                    </div>
                  </div>
                  <CardContent className="p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Avatar className="w-8 h-8">
                        <AvatarImage src="https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=150&h=150&fit=crop&crop=face" />
                        <AvatarFallback>BF</AvatarFallback>
                      </Avatar>
                      <div>
                        <h4 className="font-medium text-sm">Bangkok Food Market</h4>
                        <p className="text-xs text-muted-foreground">12 min ago</p>
                      </div>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3">
                      Live cooking demonstration with Tom Yum soup special offers
                    </p>
                    <Button 
                      className="w-full" 
                      size="sm"
                      onClick={() => onLiveStreamClick?.('food-market-stream')}
                    >
                      Join Live Stream
                    </Button>
                  </CardContent>
                </Card>

                {/* Live Stream 2 */}
                <Card className="overflow-hidden hover:shadow-md transition-shadow">
                  <div className="relative">
                    <ImageWithFallback
                      src="https://images.unsplash.com/photo-1441986300917-64674bd600d8?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwZmFzaGlvbiUyMGJvdXRpcXVlfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral"
                      alt="Fashion Show Live"
                      className="w-full h-32 object-cover"
                    />
                    <Badge className="absolute top-2 left-2 bg-red-500">
                      <div className="w-2 h-2 bg-white rounded-full mr-1 animate-pulse"></div>
                      LIVE
                    </Badge>
                    <div className="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
                      432 watching
                    </div>
                  </div>
                  <CardContent className="p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Avatar className="w-8 h-8">
                        <AvatarImage src="https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face" />
                        <AvatarFallback>TF</AvatarFallback>
                      </Avatar>
                      <div>
                        <h4 className="font-medium text-sm">Thai Fashion Live</h4>
                        <p className="text-xs text-muted-foreground">8 min ago</p>
                      </div>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3">
                      Latest Thai silk collection with live shopping discounts
                    </p>
                    <Button 
                      className="w-full" 
                      size="sm"
                      onClick={() => onLiveStreamClick?.('fashion-stream')}
                    >
                      Join Live Stream
                    </Button>
                  </CardContent>
                </Card>

                {/* Live Stream 3 */}
                <Card className="overflow-hidden hover:shadow-md transition-shadow">
                  <div className="relative">
                    <ImageWithFallback
                      src="https://images.unsplash.com/photo-1488459716781-31db52582fe9?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwZ3JvY2VyeSUyMG1hcmtldHxlbnwxfHx8fDE3NTgzOTQ1MTV8MA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral"
                      alt="Fresh Market Live"
                      className="w-full h-32 object-cover"
                    />
                    <Badge className="absolute top-2 left-2 bg-green-500">
                      <div className="w-2 h-2 bg-white rounded-full mr-1 animate-pulse"></div>
                      STARTING
                    </Badge>
                    <div className="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
                      124 waiting
                    </div>
                  </div>
                  <CardContent className="p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Avatar className="w-8 h-8">
                        <AvatarImage src="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face" />
                        <AvatarFallback>FM</AvatarFallback>
                      </Avatar>
                      <div>
                        <h4 className="font-medium text-sm">Fresh Market 24/7</h4>
                        <p className="text-xs text-muted-foreground">Starting in 5 min</p>
                      </div>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3">
                      Fresh produce showcase with wholesale prices
                    </p>
                    <Button 
                      variant="outline" 
                      className="w-full" 
                      size="sm"
                      onClick={() => onLiveStreamClick?.('grocery-stream')}
                    >
                      Set Reminder
                    </Button>
                  </CardContent>
                </Card>

                {/* Live Stream 4 */}
                <Card className="overflow-hidden hover:shadow-md transition-shadow">
                  <div className="relative">
                    <ImageWithFallback
                      src="https://images.unsplash.com/photo-1441986300917-64674bd600d8?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0ZWNoJTIwc3RvcmUlMjB0aGFpbGFuZHxlbnwxfHx8fDE3NTgzOTQ1MTV8MA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral"
                      alt="Tech Zone Live"
                      className="w-full h-32 object-cover"
                    />
                    <Badge className="absolute top-2 left-2 bg-red-500">
                      <div className="w-2 h-2 bg-white rounded-full mr-1 animate-pulse"></div>
                      LIVE
                    </Badge>
                    <div className="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
                      289 watching
                    </div>
                  </div>
                  <CardContent className="p-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Avatar className="w-8 h-8">
                        <AvatarImage src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face" />
                        <AvatarFallback>TZ</AvatarFallback>
                      </Avatar>
                      <div>
                        <h4 className="font-medium text-sm">Tech Zone Thailand</h4>
                        <p className="text-xs text-muted-foreground">15 min ago</p>
                      </div>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3">
                      Latest smartphone unboxing with flash sale prices
                    </p>
                    <Button 
                      className="w-full" 
                      size="sm"
                      onClick={() => onLiveStreamClick?.('tech-stream')}
                    >
                      Join Live Stream
                    </Button>
                  </CardContent>
                </Card>
              </div>
            </ScrollArea>
          </TabsContent>
        </Tabs>
      </div>

      {/* Payment Dialog */}
      <Dialog open={showPayment} onOpenChange={setShowPayment}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Complete Payment</DialogTitle>
          </DialogHeader>
          
          {selectedProduct && (
            <div className="space-y-4">
              <div className="flex items-center gap-3 p-3 bg-muted rounded-lg">
                <ImageWithFallback
                  src={selectedProduct.image}
                  alt={selectedProduct.name}
                  className="w-16 h-16 object-cover rounded"
                />
                <div>
                  <h4 className="font-medium">{selectedProduct.name}</h4>
                  <p className="text-sm text-muted-foreground">{selectedProduct.merchant}</p>
                  <p className="font-bold">฿{selectedProduct.price}</p>
                </div>
              </div>

              <div className="space-y-3">
                <h4 className="font-medium">Choose Payment Method</h4>
                
                <Button
                  variant="outline"
                  className="w-full justify-start gap-3 h-auto p-4"
                  onClick={() => handlePayment('PromptPay')}
                >
                  <QrCode className="w-6 h-6" />
                  <div className="text-left">
                    <p className="font-medium">PromptPay QR</p>
                    <p className="text-sm text-muted-foreground">Scan QR to pay instantly</p>
                  </div>
                </Button>

                <Button
                  variant="outline"
                  className="w-full justify-start gap-3 h-auto p-4"
                  onClick={() => handlePayment('Wallet')}
                >
                  <Wallet className="w-6 h-6" />
                  <div className="text-left">
                    <p className="font-medium">Telegram Wallet</p>
                    <p className="text-sm text-muted-foreground">Balance: ฿1,250.00</p>
                  </div>
                </Button>

                <Button
                  variant="outline"
                  className="w-full justify-start gap-3 h-auto p-4"
                  onClick={() => handlePayment('Credit Card')}
                >
                  <CreditCard className="w-6 h-6" />
                  <div className="text-left">
                    <p className="font-medium">Credit/Debit Card</p>
                    <p className="text-sm text-muted-foreground">Visa, Mastercard accepted</p>
                  </div>
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}