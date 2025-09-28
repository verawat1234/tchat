import React, { useState, useMemo } from 'react';
import { useGetShopProductsQuery } from '../services/microservicesApi';
import { 
  ArrowLeft, 
  Star, 
  MapPin, 
  Clock, 
  Users, 
  Store, 
  Heart, 
  Share, 
  MessageCircle,
  Phone,
  Search,
  Filter,
  ShoppingCart,
  Zap,
  TrendingUp,
  Eye
} from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";

interface ShopPageProps {
  user: any;
  shopId: string;
  onBack: () => void;
  onProductClick: (productId: string) => void;
  onAddToCart: (productId: string) => void;
}

interface Product {
  id: string;
  name: string;
  price: number;
  currency: string;
  image: string;
  rating: number;
  category: string;
  isLive?: boolean;
  discount?: number;
  isHot?: boolean;
  orders?: number;
  stock?: number;
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
  phone?: string;
  address?: string;
  openHours?: string;
  joinedDate?: string;
}

export function ShopPage({ user, shopId, onBack, onProductClick, onAddToCart }: ShopPageProps) {
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [isFollowing, setIsFollowing] = useState(false);

  // Mock shop data - in real app would fetch based on shopId
  const shop: Shop = {
    id: shopId,
    name: 'Bangkok Street Food Palace',
    description: 'Authentic Thai street food made fresh daily with traditional recipes passed down through generations. We pride ourselves on using the freshest ingredients and maintaining the authentic taste of Bangkok street food.',
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
    specialOffer: '20% off first order',
    phone: '+66 89 123 4567',
    address: 'Chatuchak Market, Phahon Yothin Rd, Bangkok 10900',
    openHours: '6:00 AM - 11:00 PM',
    joinedDate: 'March 2020'
  };

  // RTK Query for shop products
  const {
    data: productsData,
    isLoading: productsLoading,
    error: productsError
  } = useGetShopProductsQuery({
    shopId: 'thai-street-food-paradise',
    page: 1,
    limit: 20,
    category: activeCategory === 'all' ? undefined : activeCategory
  });

  const products: Product[] = useMemo(() => {
    if (productsLoading || !productsData) {
      // Fallback data while loading
      return [
        {
          id: '1',
          name: 'Pad Thai Goong',
          price: 45,
          currency: 'THB',
          image: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
          rating: 4.8,
          category: 'main',
          isLive: true,
          isHot: true,
          orders: 1234,
          stock: 15
        },
        {
          id: '2',
          name: 'Som Tam Thai',
          price: 35,
          currency: 'THB',
          image: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
          rating: 4.6,
          category: 'salad',
          discount: 20,
          isHot: true,
          orders: 892,
          stock: 8
        }
      ];
    }

    return productsData.map((product: any) => ({
      id: product.id || product.product_id || `product-${Math.random()}`,
      name: product.name || product.title || 'Product',
      price: product.price || 0,
      currency: product.currency || 'THB',
      image: product.image || product.image_url || product.thumbnail || 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=1080',
      rating: product.rating || product.average_rating || 4.5,
      category: product.category || product.product_category || 'main',
      isLive: product.isLive || product.is_live || false,
      isHot: product.isHot || product.is_hot || product.trending || false,
      discount: product.discount || product.discount_percentage || undefined,
      orders: product.orders || product.total_orders || 0,
      stock: product.stock || product.available_quantity || 0
    }));
  }, [productsData, productsLoading, activeCategory]);

  const categories = [
    { id: 'all', name: 'All', count: products.length },
    { id: 'main', name: 'Main Dishes', count: products.filter(p => p.category === 'main').length },
    { id: 'salad', name: 'Salads', count: products.filter(p => p.category === 'salad').length },
    { id: 'soup', name: 'Soups', count: products.filter(p => p.category === 'soup').length },
    { id: 'dessert', name: 'Desserts', count: products.filter(p => p.category === 'dessert').length },
    { id: 'drink', name: 'Drinks', count: products.filter(p => p.category === 'drink').length }
  ];

  const filteredProducts = products.filter(product => {
    const matchesCategory = selectedCategory === 'all' || product.category === selectedCategory;
    const matchesSearch = product.name.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesCategory && matchesSearch;
  });

  const handleFollow = () => {
    setIsFollowing(!isFollowing);
    toast.success(isFollowing ? 'Unfollowed shop' : 'Following shop');
  };

  const handleAddToCart = (productId: string) => {
    onAddToCart(productId);
    toast.success('Added to cart!');
  };

  const handleShare = () => {
    toast.success('Shop link copied to clipboard!');
  };

  const handleMessage = () => {
    toast.success('Opening chat with shop owner...');
  };

  const handleCall = () => {
    toast.success('Calling shop...');
  };

  return (
    <div className="h-screen bg-background flex flex-col">
      {/* Header */}
      <div className="border-b border-border p-4 flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={onBack}>
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h1 className="text-lg font-medium flex-1 truncate">{shop.name}</h1>
        <Button variant="ghost" size="icon" onClick={handleShare}>
          <Share className="w-5 h-5" />
        </Button>
      </div>

      <ScrollArea className="flex-1">
        <div className="space-y-3 sm:space-y-4 lg:space-y-6 pb-3 sm:pb-4 lg:pb-6">
          {/* Shop Header */}
          <div className="relative">
            <ImageWithFallback
              src={shop.image}
              alt={shop.name}
              className="w-full h-40 sm:h-48 lg:h-56 object-cover"
            />
            {shop.isHot && (
              <Badge className="absolute top-2 left-2 sm:top-4 sm:left-4 bg-red-500 text-xs sm:text-sm">
                <TrendingUp className="w-2.5 h-2.5 sm:w-3 sm:h-3 mr-0.5 sm:mr-1" />
                HOT
              </Badge>
            )}
            {shop.specialOffer && (
              <Badge className="absolute top-2 right-2 sm:top-4 sm:right-4 bg-green-500 text-xs sm:text-sm">
                {shop.specialOffer}
              </Badge>
            )}
          </div>

          {/* Shop Info */}
          <div className="px-3 sm:px-4 lg:px-6 space-y-3 sm:space-y-4">
            <div className="flex items-start gap-3 sm:gap-4">
              <Avatar className="w-12 h-12 sm:w-16 sm:h-16 lg:w-20 lg:h-20 flex-shrink-0">
                <AvatarImage src={shop.avatar} />
                <AvatarFallback className="text-sm sm:text-base lg:text-lg">
                  <Store className="w-5 h-5 sm:w-8 sm:h-8 lg:w-10 lg:h-10" />
                </AvatarFallback>
              </Avatar>
              
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-1 sm:gap-2 mb-1 sm:mb-2">
                  <h2 className="text-base sm:text-xl lg:text-2xl font-medium truncate">{shop.name}</h2>
                  {shop.isVerified && (
                    <Star className="w-4 h-4 sm:w-5 sm:h-5 lg:w-6 lg:h-6 text-chart-1 fill-chart-1 flex-shrink-0" />
                  )}
                </div>
                
                <div className="flex items-center gap-2 sm:gap-4 mb-1 sm:mb-2 flex-wrap">
                  <div className="flex items-center gap-0.5 sm:gap-1">
                    <Star className="w-3 h-3 sm:w-4 sm:h-4 text-yellow-500 fill-yellow-500" />
                    <span className="text-xs sm:text-sm font-medium">{shop.rating}</span>
                  </div>
                  <div className="flex items-center gap-0.5 sm:gap-1">
                    <Clock className="w-3 h-3 sm:w-4 sm:h-4 text-muted-foreground" />
                    <span className="text-xs sm:text-sm text-muted-foreground whitespace-nowrap">{shop.deliveryTime}</span>
                  </div>
                  <div className="flex items-center gap-0.5 sm:gap-1">
                    <MapPin className="w-3 h-3 sm:w-4 sm:h-4 text-muted-foreground" />
                    <span className="text-xs sm:text-sm text-muted-foreground whitespace-nowrap">{shop.distance}</span>
                  </div>
                </div>

                <div className="flex items-center gap-2 sm:gap-4 text-xs sm:text-sm text-muted-foreground flex-wrap">
                  <div className="flex items-center gap-0.5 sm:gap-1">
                    <Users className="w-3 h-3 sm:w-4 sm:h-4" />
                    <span className="whitespace-nowrap">{shop.followers.toLocaleString()} followers</span>
                  </div>
                  <div className="flex items-center gap-0.5 sm:gap-1">
                    <Store className="w-3 h-3 sm:w-4 sm:h-4" />
                    <span className="whitespace-nowrap">{shop.products} products</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex gap-2 sm:gap-3">
              <Button 
                className="flex-1 h-9 sm:h-10 lg:h-11 text-xs sm:text-sm touch-manipulation" 
                variant={isFollowing ? "outline" : "default"}
                onClick={handleFollow}
              >
                <Heart className={`w-3 h-3 sm:w-4 sm:h-4 mr-1 sm:mr-2 ${isFollowing ? 'fill-red-500 text-red-500' : ''}`} />
                <span className="hidden sm:inline">{isFollowing ? 'Following' : 'Follow'}</span>
                <span className="sm:hidden">{isFollowing ? '✓' : '+'}</span>
              </Button>
              <Button 
                variant="outline" 
                onClick={handleMessage}
                className="h-9 sm:h-10 lg:h-11 px-2 sm:px-4 text-xs sm:text-sm touch-manipulation"
              >
                <MessageCircle className="w-3 h-3 sm:w-4 sm:h-4 mr-1 sm:mr-2" />
                <span className="hidden sm:inline">Message</span>
                <span className="sm:hidden">Chat</span>
              </Button>
              <Button 
                variant="outline" 
                size="icon" 
                onClick={handleCall}
                className="w-9 h-9 sm:w-10 sm:h-10 lg:w-11 lg:h-11 flex-shrink-0 touch-manipulation"
              >
                <Phone className="w-3 h-3 sm:w-4 sm:h-4" />
              </Button>
            </div>

            <p className="text-xs sm:text-sm lg:text-base text-muted-foreground leading-relaxed">
              {shop.description}
            </p>
          </div>

          {/* Shop Details Tabs */}
          <div className="px-3 sm:px-4 lg:px-6">
            <Tabs defaultValue="products" className="w-full">
              <TabsList className="grid w-full grid-cols-3 h-9 sm:h-10 lg:h-11">
                <TabsTrigger value="products" className="text-xs sm:text-sm touch-manipulation">Products</TabsTrigger>
                <TabsTrigger value="reviews" className="text-xs sm:text-sm touch-manipulation">Reviews</TabsTrigger>
                <TabsTrigger value="info" className="text-xs sm:text-sm touch-manipulation">Info</TabsTrigger>
              </TabsList>

              <TabsContent value="products" className="space-y-3 sm:space-y-4 mt-3 sm:mt-4">
                {/* Search and Filter */}
                <div className="flex gap-2">
                  <div className="relative flex-1">
                    <Search className="absolute left-2 sm:left-3 top-2.5 sm:top-3 w-3 h-3 sm:w-4 sm:h-4 text-muted-foreground" />
                    <Input
                      placeholder="Search products..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-8 sm:pl-10 h-9 sm:h-10 text-xs sm:text-sm"
                    />
                  </div>
                  <Button variant="outline" size="icon" className="w-9 h-9 sm:w-10 sm:h-10 flex-shrink-0 touch-manipulation">
                    <Filter className="w-3 h-3 sm:w-4 sm:h-4" />
                  </Button>
                </div>

                {/* Category Filters */}
                <ScrollArea className="w-full">
                  <div className="flex gap-1.5 sm:gap-2 pb-1 sm:pb-2">
                    {categories.map((category) => (
                      <Button
                        key={category.id}
                        variant={selectedCategory === category.id ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => setSelectedCategory(category.id)}
                        className="whitespace-nowrap text-xs sm:text-sm h-8 sm:h-9 px-2 sm:px-3 touch-manipulation"
                      >
                        {category.name} ({category.count})
                      </Button>
                    ))}
                  </div>
                </ScrollArea>

                {/* Products Grid */}
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
                  {filteredProducts.map((product) => (
                    <Card key={product.id} className="overflow-hidden cursor-pointer hover:shadow-md transition-shadow touch-manipulation">
                      <div className="relative" onClick={() => onProductClick(product.id)}>
                        <ImageWithFallback
                          src={product.image}
                          alt={product.name}
                          className="w-full h-28 sm:h-32 lg:h-36 object-cover"
                        />
                        {product.isLive && (
                          <Badge className="absolute top-1 left-1 sm:top-2 sm:left-2 bg-red-500 text-xs">
                            <Zap className="w-2 h-2 sm:w-3 sm:h-3 mr-0.5 sm:mr-1" />
                            LIVE
                          </Badge>
                        )}
                        {product.isHot && (
                          <Badge className="absolute top-1 right-1 sm:top-2 sm:right-2 bg-orange-500 text-xs">
                            <TrendingUp className="w-2 h-2 sm:w-3 sm:h-3 mr-0.5 sm:mr-1" />
                            HOT
                          </Badge>
                        )}
                        {product.discount && (
                          <Badge className="absolute bottom-1 right-1 sm:bottom-2 sm:right-2 bg-green-500 text-xs">
                            -{product.discount}%
                          </Badge>
                        )}
                      </div>
                      
                      <CardContent className="p-2 sm:p-3">
                        <h3 className="text-sm sm:text-base font-medium mb-1 sm:mb-2 cursor-pointer line-clamp-2 leading-tight" onClick={() => onProductClick(product.id)}>
                          {product.name}
                        </h3>
                        
                        <div className="flex items-center gap-1 sm:gap-2 mb-1 sm:mb-2 text-xs sm:text-sm flex-wrap">
                          <div className="flex items-center gap-0.5 sm:gap-1">
                            <Star className="w-2.5 h-2.5 sm:w-3 sm:h-3 text-yellow-500 fill-yellow-500" />
                            <span className="font-medium">{product.rating}</span>
                          </div>
                          {product.orders && (
                            <>
                              <span className="text-muted-foreground hidden sm:inline">•</span>
                              <div className="flex items-center gap-0.5 sm:gap-1">
                                <Eye className="w-2.5 h-2.5 sm:w-3 sm:h-3 text-muted-foreground" />
                                <span className="text-muted-foreground whitespace-nowrap">{product.orders} sold</span>
                              </div>
                            </>
                          )}
                          {product.stock && (
                            <>
                              <span className="text-muted-foreground hidden sm:inline">•</span>
                              <span className="text-muted-foreground whitespace-nowrap">{product.stock} left</span>
                            </>
                          )}
                        </div>
                        
                        <div className="flex items-center justify-between gap-2">
                          <div className="flex items-center gap-1 sm:gap-2 min-w-0">
                            <span className="text-sm sm:text-base font-bold truncate">
                              ฿{product.price}
                            </span>
                            {product.discount && (
                              <span className="text-xs sm:text-sm text-muted-foreground line-through truncate">
                                ฿{Math.round(product.price / (1 - product.discount / 100))}
                              </span>
                            )}
                          </div>
                          
                          <Button 
                            size="sm" 
                            onClick={(e) => {
                              e.stopPropagation();
                              handleAddToCart(product.id);
                            }}
                            className="h-7 sm:h-8 px-2 sm:px-3 text-xs sm:text-sm flex-shrink-0 touch-manipulation"
                          >
                            <ShoppingCart className="w-2.5 h-2.5 sm:w-3 sm:h-3 mr-1" />
                            <span className="hidden sm:inline">Add</span>
                            <span className="sm:hidden">+</span>
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </TabsContent>

              <TabsContent value="reviews" className="space-y-3 sm:space-y-4 mt-3 sm:mt-4">
                <div className="space-y-3 sm:space-y-4">
                  <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3 sm:gap-4">
                    <div className="text-center flex-shrink-0">
                      <div className="text-xl sm:text-2xl lg:text-3xl font-bold">{shop.rating}</div>
                      <div className="flex items-center gap-0.5 sm:gap-1 justify-center">
                        {Array.from({ length: 5 }).map((_, i) => (
                          <Star 
                            key={i} 
                            className={`w-3 h-3 sm:w-4 sm:h-4 ${i < Math.floor(shop.rating) ? 'text-yellow-500 fill-yellow-500' : 'text-muted-foreground'}`} 
                          />
                        ))}
                      </div>
                      <div className="text-xs sm:text-sm text-muted-foreground">1,234 reviews</div>
                    </div>
                    <div className="flex-1 space-y-1.5 sm:space-y-2 w-full">
                      {[5, 4, 3, 2, 1].map(stars => (
                        <div key={stars} className="flex items-center gap-1.5 sm:gap-2">
                          <span className="text-xs sm:text-sm w-2 flex-shrink-0">{stars}</span>
                          <Star className="w-2.5 h-2.5 sm:w-3 sm:h-3 text-yellow-500 fill-yellow-500 flex-shrink-0" />
                          <div className="flex-1 h-1.5 sm:h-2 bg-muted rounded-full">
                            <div 
                              className="h-full bg-yellow-500 rounded-full transition-all duration-300" 
                              style={{ width: `${stars === 5 ? 70 : stars === 4 ? 20 : stars === 3 ? 7 : stars === 2 ? 2 : 1}%` }}
                            />
                          </div>
                          <span className="text-xs sm:text-sm text-muted-foreground w-6 sm:w-8 text-right flex-shrink-0">
                            {stars === 5 ? '70%' : stars === 4 ? '20%' : stars === 3 ? '7%' : stars === 2 ? '2%' : '1%'}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </TabsContent>

              <TabsContent value="info" className="space-y-3 sm:space-y-4 mt-3 sm:mt-4">
                <div className="space-y-3 sm:space-y-4">
                  <div>
                    <h4 className="text-sm sm:text-base font-medium mb-1.5 sm:mb-2">Contact Information</h4>
                    <div className="space-y-1.5 sm:space-y-2 text-xs sm:text-sm">
                      <div className="flex items-center gap-1.5 sm:gap-2">
                        <Phone className="w-3 h-3 sm:w-4 sm:h-4 text-muted-foreground flex-shrink-0" />
                        <span>{shop.phone}</span>
                      </div>
                      <div className="flex items-start gap-1.5 sm:gap-2">
                        <MapPin className="w-3 h-3 sm:w-4 sm:h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                        <span className="leading-relaxed">{shop.address}</span>
                      </div>
                    </div>
                  </div>

                  <div>
                    <h4 className="text-sm sm:text-base font-medium mb-1.5 sm:mb-2">Business Hours</h4>
                    <p className="text-xs sm:text-sm text-muted-foreground leading-relaxed">{shop.openHours}</p>
                  </div>

                  <div>
                    <h4 className="text-sm sm:text-base font-medium mb-1.5 sm:mb-2">About</h4>
                    <p className="text-xs sm:text-sm text-muted-foreground leading-relaxed">
                      Joined Telegram SEA Marketplace in {shop.joinedDate}
                    </p>
                  </div>
                </div>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      </ScrollArea>
    </div>
  );
}