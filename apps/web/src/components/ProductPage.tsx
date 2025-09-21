import React, { useState } from 'react';
import { 
  ArrowLeft, 
  Star, 
  Heart, 
  Share, 
  MessageCircle,
  ShoppingCart,
  Plus,
  Minus,
  Truck,
  Shield,
  Clock,
  Store,
  Eye,
  ThumbsUp,
  ThumbsDown,
  MoreHorizontal
} from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Separator } from './ui/separator';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";

interface ProductPageProps {
  user: any;
  productId: string;
  onBack: () => void;
  onShopClick: (shopId: string) => void;
  onAddToCart: (productId: string, quantity: number) => void;
  onBuyNow: (productId: string, quantity: number) => void;
}

interface Product {
  id: string;
  name: string;
  price: number;
  originalPrice?: number;
  currency: string;
  images: string[];
  description: string;
  merchant: {
    id: string;
    name: string;
    avatar: string;
    rating: number;
    isVerified: boolean;
  };
  rating: number;
  reviewCount: number;
  category: string;
  isLive?: boolean;
  discount?: number;
  isHot?: boolean;
  orders?: number;
  stock?: number;
  specifications?: { [key: string]: string };
  variants?: { name: string; options: string[] }[];
}

interface Review {
  id: string;
  user: {
    name: string;
    avatar?: string;
  };
  rating: number;
  comment: string;
  date: string;
  helpful: number;
  images?: string[];
}

export function ProductPage({ user, productId, onBack, onShopClick, onAddToCart, onBuyNow }: ProductPageProps) {
  const [quantity, setQuantity] = useState(1);
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);
  const [isFavorite, setIsFavorite] = useState(false);
  const [selectedVariants, setSelectedVariants] = useState<{ [key: string]: string }>({});

  // Mock product data - in real app would fetch based on productId
  const product: Product = {
    id: productId,
    name: 'Pad Thai Goong (Authentic Bangkok Style)',
    price: 45,
    originalPrice: 60,
    currency: 'THB',
    images: [
      'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwbm9vZGxlc3xlbnwxfHx8fDE3NTgzOTQ1MTd8MA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'
    ],
    description: 'Authentic Pad Thai made with fresh rice noodles, tiger prawns, bean sprouts, eggs, and our signature tamarind sauce. Prepared using traditional Bangkok street food techniques passed down through generations. Each serving includes fresh lime, crushed peanuts, and chili flakes on the side.',
    merchant: {
      id: 'shop1',
      name: 'Bangkok Street Food Palace',
      avatar: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=150&h=150&fit=crop&crop=face',
      rating: 4.8,
      isVerified: true
    },
    rating: 4.8,
    reviewCount: 1234,
    category: 'food',
    isLive: true,
    discount: 25,
    isHot: true,
    orders: 1234,
    stock: 15,
    specifications: {
      'Serving Size': '1 plate (350g)',
      'Spice Level': 'Medium (customizable)',
      'Allergens': 'Shellfish, Eggs, Nuts',
      'Preparation Time': '8-12 minutes',
      'Calories': '~450 per serving'
    },
    variants: [
      {
        name: 'Spice Level',
        options: ['Mild', 'Medium', 'Spicy', 'Extra Spicy']
      },
      {
        name: 'Add-ons',
        options: ['Extra Prawns (+฿15)', 'Extra Vegetables (+฿10)', 'Extra Sauce (+฿5)']
      }
    ]
  };

  // Mock reviews data
  const reviews: Review[] = [
    {
      id: '1',
      user: {
        name: 'Somchai K.',
        avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face'
      },
      rating: 5,
      comment: 'Amazing authentic taste! Reminds me of the street food from my childhood in Bangkok. The prawns were fresh and the sauce was perfectly balanced.',
      date: '2 days ago',
      helpful: 12,
      images: ['https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=400&h=300&fit=crop']
    },
    {
      id: '2',
      user: {
        name: 'Sarah M.',
        avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face'
      },
      rating: 4,
      comment: 'Very good Pad Thai! Fast delivery and still hot when arrived. Would definitely order again.',
      date: '1 week ago',
      helpful: 8
    },
    {
      id: '3',
      user: {
        name: 'Michael T.'
      },
      rating: 5,
      comment: 'Best Pad Thai in Bangkok! The portion size is generous and the flavor is incredible.',
      date: '2 weeks ago',
      helpful: 15
    }
  ];

  const handleQuantityChange = (delta: number) => {
    const newQuantity = Math.max(1, Math.min(product.stock || 99, quantity + delta));
    setQuantity(newQuantity);
  };

  const handleAddToCart = () => {
    onAddToCart(product.id, quantity);
    toast.success(`Added ${quantity}x ${product.name} to cart!`);
  };

  const handleBuyNow = () => {
    onBuyNow(product.id, quantity);
  };

  const handleFavorite = () => {
    setIsFavorite(!isFavorite);
    toast.success(isFavorite ? 'Removed from favorites' : 'Added to favorites');
  };

  const handleShare = () => {
    toast.success('Product link copied to clipboard!');
  };

  const handleMessage = () => {
    toast.success('Opening chat with merchant...');
  };

  const handleVariantSelect = (variantName: string, option: string) => {
    setSelectedVariants(prev => ({
      ...prev,
      [variantName]: option
    }));
  };

  return (
    <div className="h-screen bg-background flex flex-col">
      {/* Header */}
      <div className="border-b border-border p-4 flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={onBack}>
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h1 className="text-lg font-medium flex-1 truncate">Product Details</h1>
        <Button variant="ghost" size="icon" onClick={handleShare}>
          <Share className="w-5 h-5" />
        </Button>
        <Button variant="ghost" size="icon" onClick={handleFavorite}>
          <Heart className={`w-5 h-5 ${isFavorite ? 'fill-red-500 text-red-500' : ''}`} />
        </Button>
      </div>

      <ScrollArea className="flex-1">
        <div className="space-y-6">
          {/* Product Images */}
          <div className="relative">
            <ImageWithFallback
              src={product.images[selectedImageIndex]}
              alt={product.name}
              className="w-full h-80 object-cover"
            />
            
            {/* Badges */}
            <div className="absolute top-4 left-4 flex flex-col gap-2">
              {product.isLive && (
                <Badge className="bg-red-500">
                  🔴 LIVE
                </Badge>
              )}
              {product.isHot && (
                <Badge className="bg-orange-500">
                  🔥 HOT
                </Badge>
              )}
              {product.discount && (
                <Badge className="bg-green-500">
                  -{product.discount}% OFF
                </Badge>
              )}
            </div>

            {/* Image Indicators */}
            {product.images.length > 1 && (
              <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex gap-2">
                {product.images.map((_, index) => (
                  <button
                    key={index}
                    onClick={() => setSelectedImageIndex(index)}
                    className={`w-2 h-2 rounded-full ${
                      index === selectedImageIndex ? 'bg-white' : 'bg-white/50'
                    }`}
                  />
                ))}
              </div>
            )}
          </div>

          {/* Product Info */}
          <div className="px-4 space-y-4">
            <div>
              <h1 className="text-xl font-medium mb-2">{product.name}</h1>
              <div className="flex items-center gap-2 mb-3">
                <div className="flex items-center gap-1">
                  <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                  <span className="font-medium">{product.rating}</span>
                </div>
                <span className="text-sm text-muted-foreground">
                  ({product.reviewCount.toLocaleString()} reviews)
                </span>
                <span className="text-sm text-muted-foreground">•</span>
                <div className="flex items-center gap-1">
                  <Eye className="w-4 h-4 text-muted-foreground" />
                  <span className="text-sm text-muted-foreground">{product.orders} sold</span>
                </div>
              </div>

              {/* Price */}
              <div className="flex items-center gap-3 mb-4">
                <span className="text-2xl font-bold text-red-500">
                  ฿{product.price}
                </span>
                {product.originalPrice && (
                  <span className="text-lg text-muted-foreground line-through">
                    ฿{product.originalPrice}
                  </span>
                )}
                {product.discount && (
                  <Badge variant="destructive">
                    Save ฿{product.originalPrice! - product.price}
                  </Badge>
                )}
              </div>

              {/* Stock */}
              {product.stock && (
                <div className="flex items-center gap-2 mb-4">
                  <span className="text-sm text-muted-foreground">Stock:</span>
                  <Badge variant={product.stock > 10 ? "secondary" : "destructive"}>
                    {product.stock} left
                  </Badge>
                </div>
              )}
            </div>

            {/* Merchant Info */}
            <Card className="cursor-pointer" onClick={() => onShopClick(product.merchant.id)}>
              <CardContent className="p-4">
                <div className="flex items-center gap-3">
                  <Avatar className="w-12 h-12">
                    <AvatarImage src={product.merchant.avatar} />
                    <AvatarFallback>
                      <Store className="w-6 h-6" />
                    </AvatarFallback>
                  </Avatar>
                  
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <h3 className="font-medium">{product.merchant.name}</h3>
                      {product.merchant.isVerified && (
                        <Star className="w-4 h-4 text-chart-1 fill-chart-1" />
                      )}
                    </div>
                    <div className="flex items-center gap-1">
                      <Star className="w-3 h-3 text-yellow-500 fill-yellow-500" />
                      <span className="text-sm">{product.merchant.rating}</span>
                      <span className="text-sm text-muted-foreground ml-2">Shop rating</span>
                    </div>
                  </div>
                  
                  <Button variant="outline" size="sm" onClick={handleMessage}>
                    <MessageCircle className="w-4 h-4 mr-2" />
                    Chat
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* Variants */}
            {product.variants && (
              <div className="space-y-4">
                {product.variants.map((variant) => (
                  <div key={variant.name}>
                    <h4 className="font-medium mb-2">{variant.name}</h4>
                    <div className="flex flex-wrap gap-2">
                      {variant.options.map((option) => (
                        <Button
                          key={option}
                          variant={selectedVariants[variant.name] === option ? "default" : "outline"}
                          size="sm"
                          onClick={() => handleVariantSelect(variant.name, option)}
                        >
                          {option}
                        </Button>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}

            {/* Quantity and Actions */}
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <span className="font-medium">Quantity:</span>
                <div className="flex items-center gap-2">
                  <Button 
                    variant="outline" 
                    size="icon"
                    onClick={() => handleQuantityChange(-1)}
                    disabled={quantity <= 1}
                  >
                    <Minus className="w-4 h-4" />
                  </Button>
                  <span className="w-12 text-center font-medium">{quantity}</span>
                  <Button 
                    variant="outline" 
                    size="icon"
                    onClick={() => handleQuantityChange(1)}
                    disabled={quantity >= (product.stock || 99)}
                  >
                    <Plus className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              <div className="flex gap-2">
                <Button variant="outline" className="flex-1" onClick={handleAddToCart}>
                  <ShoppingCart className="w-4 h-4 mr-2" />
                  Add to Cart
                </Button>
                <Button className="flex-1" onClick={handleBuyNow}>
                  Buy Now
                </Button>
              </div>
            </div>

            {/* Delivery Info */}
            <Card>
              <CardContent className="p-4 space-y-3">
                <div className="flex items-center gap-3">
                  <Truck className="w-5 h-5 text-chart-2" />
                  <div>
                    <p className="font-medium">Fast Delivery</p>
                    <p className="text-sm text-muted-foreground">15-20 minutes</p>
                  </div>
                </div>
                <Separator />
                <div className="flex items-center gap-3">
                  <Shield className="w-5 h-5 text-green-500" />
                  <div>
                    <p className="font-medium">Quality Guarantee</p>
                    <p className="text-sm text-muted-foreground">Fresh ingredients only</p>
                  </div>
                </div>
                <Separator />
                <div className="flex items-center gap-3">
                  <Clock className="w-5 h-5 text-chart-4" />
                  <div>
                    <p className="font-medium">Preparation Time</p>
                    <p className="text-sm text-muted-foreground">8-12 minutes</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Product Details Tabs */}
          <Tabs defaultValue="description" className="px-4">
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="description">Description</TabsTrigger>
              <TabsTrigger value="specifications">Specs</TabsTrigger>
              <TabsTrigger value="reviews">Reviews</TabsTrigger>
            </TabsList>

            <TabsContent value="description" className="mt-4">
              <div className="space-y-4">
                <p className="text-sm leading-relaxed">{product.description}</p>
              </div>
            </TabsContent>

            <TabsContent value="specifications" className="mt-4">
              {product.specifications && (
                <div className="space-y-3">
                  {Object.entries(product.specifications).map(([key, value]) => (
                    <div key={key} className="flex justify-between items-center py-2 border-b border-border last:border-0">
                      <span className="text-sm font-medium">{key}</span>
                      <span className="text-sm text-muted-foreground">{value}</span>
                    </div>
                  ))}
                </div>
              )}
            </TabsContent>

            <TabsContent value="reviews" className="mt-4">
              <div className="space-y-4">
                {/* Review Summary */}
                <Card>
                  <CardContent className="p-4">
                    <div className="flex items-center gap-4 mb-4">
                      <div className="text-center">
                        <div className="text-2xl font-bold">{product.rating}</div>
                        <div className="flex items-center gap-1 justify-center">
                          {Array.from({ length: 5 }).map((_, i) => (
                            <Star 
                              key={i} 
                              className={`w-4 h-4 ${i < Math.floor(product.rating) ? 'text-yellow-500 fill-yellow-500' : 'text-muted-foreground'}`} 
                            />
                          ))}
                        </div>
                        <div className="text-sm text-muted-foreground">{product.reviewCount} reviews</div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                {/* Individual Reviews */}
                <div className="space-y-4">
                  {reviews.map((review) => (
                    <Card key={review.id}>
                      <CardContent className="p-4">
                        <div className="flex items-start gap-3">
                          <Avatar className="w-10 h-10">
                            <AvatarImage src={review.user.avatar} />
                            <AvatarFallback>{review.user.name.charAt(0)}</AvatarFallback>
                          </Avatar>
                          
                          <div className="flex-1">
                            <div className="flex items-center gap-2 mb-2">
                              <span className="font-medium">{review.user.name}</span>
                              <div className="flex items-center gap-1">
                                {Array.from({ length: 5 }).map((_, i) => (
                                  <Star 
                                    key={i} 
                                    className={`w-3 h-3 ${i < review.rating ? 'text-yellow-500 fill-yellow-500' : 'text-muted-foreground'}`} 
                                  />
                                ))}
                              </div>
                              <span className="text-sm text-muted-foreground">{review.date}</span>
                            </div>
                            
                            <p className="text-sm mb-3">{review.comment}</p>
                            
                            {review.images && (
                              <div className="flex gap-2 mb-3">
                                {review.images.map((image, index) => (
                                  <ImageWithFallback
                                    key={index}
                                    src={image}
                                    alt="Review image"
                                    className="w-16 h-16 object-cover rounded"
                                  />
                                ))}
                              </div>
                            )}
                            
                            <div className="flex items-center gap-4">
                              <Button variant="ghost" size="sm">
                                <ThumbsUp className="w-4 h-4 mr-1" />
                                Helpful ({review.helpful})
                              </Button>
                              <Button variant="ghost" size="sm">
                                <ThumbsDown className="w-4 h-4 mr-1" />
                              </Button>
                              <Button variant="ghost" size="sm">
                                <MoreHorizontal className="w-4 h-4" />
                              </Button>
                            </div>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </ScrollArea>
    </div>
  );
}