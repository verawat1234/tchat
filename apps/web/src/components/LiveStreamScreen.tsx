import React, { useState, useEffect } from 'react';
import { 
  ArrowLeft, 
  Heart, 
  MessageCircle, 
  Share, 
  ShoppingCart,
  Users,
  Eye,
  Gift,
  Star,
  Volume2,
  VolumeX,
  Maximize,
  MoreVertical,
  Send,
  Smile,
  ThumbsUp,
  Zap,
  ShoppingBag,
  Crown,
  Flame
} from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './ui/dialog';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";

interface LiveStreamScreenProps {
  user: any;
  streamId: string;
  onBack: () => void;
  onProductClick?: (productId: string) => void;
  onAddToCart?: (productId: string) => void;
}

interface StreamData {
  id: string;
  title: string;
  streamer: {
    id: string;
    name: string;
    avatar: string;
    verified: boolean;
    followers: number;
  };
  viewerCount: number;
  likes: number;
  duration: string;
  category: string;
  products: Product[];
  isLive: boolean;
}

interface Product {
  id: string;
  name: string;
  price: number;
  originalPrice?: number;
  currency: string;
  image: string;
  inStock: boolean;
  discount?: number;
}

interface ChatMessage {
  id: string;
  user: {
    name: string;
    avatar?: string;
    isVip?: boolean;
    isModerator?: boolean;
  };
  message: string;
  timestamp: string;
  type: 'message' | 'gift' | 'product' | 'system';
  giftType?: string;
  productId?: string;
}

interface Gift {
  id: string;
  name: string;
  icon: string;
  price: number;
  animation?: string;
}

export function LiveStreamScreen({ user, streamId, onBack, onProductClick, onAddToCart }: LiveStreamScreenProps) {
  const [isFollowing, setIsFollowing] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [showGifts, setShowGifts] = useState(false);
  const [showProducts, setShowProducts] = useState(false);
  const [newMessage, setNewMessage] = useState('');
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [currentViewers, setCurrentViewers] = useState(0);
  const [currentLikes, setCurrentLikes] = useState(0);

  // Mock stream data
  const streamData: StreamData = {
    id: streamId,
    title: 'Bangkok Street Food Tour 🍜 Trying Authentic Pad Thai & Tom Yum!',
    streamer: {
      id: 'streamer1',
      name: 'Chef Somsak',
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
      verified: true,
      followers: 45200
    },
    viewerCount: 1247,
    likes: 3456,
    duration: '1:23:45',
    category: 'Food & Cooking',
    isLive: true,
    products: [
      {
        id: '1',
        name: 'Premium Pad Thai Kit',
        price: 89,
        originalPrice: 120,
        currency: 'THB',
        image: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
        inStock: true,
        discount: 26
      },
      {
        id: '2',
        name: 'Tom Yum Paste (Authentic)',
        price: 45,
        currency: 'THB',
        image: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0b20lMjB5dW0lMjBzb3VwfGVufDF8fHx8MTc1ODM5NDUxN3ww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
        inStock: true
      },
      {
        id: '3',
        name: 'Thai Cooking Utensil Set',
        price: 156,
        originalPrice: 200,
        currency: 'THB',
        image: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
        inStock: false,
        discount: 22
      }
    ]
  };

  // Mock gifts data
  const gifts: Gift[] = [
    { id: '1', name: 'Heart', icon: '❤️', price: 1 },
    { id: '2', name: 'Rose', icon: '🌹', price: 5 },
    { id: '3', name: 'Thai Tea', icon: '🧋', price: 10 },
    { id: '4', name: 'Gold Crown', icon: '👑', price: 50 },
    { id: '5', name: 'Dragon', icon: '🐉', price: 100 },
    { id: '6', name: 'Fireworks', icon: '🎆', price: 200 }
  ];

  // Mock chat messages
  const initialMessages: ChatMessage[] = [
    {
      id: '1',
      user: { name: 'FoodieThailand', avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face', isVip: true },
      message: 'This looks so delicious! 😍',
      timestamp: '2 min ago',
      type: 'message'
    },
    {
      id: '2',
      user: { name: 'BangkokEats' },
      message: 'Where is this place located?',
      timestamp: '1 min ago',
      type: 'message'
    },
    {
      id: '3',
      user: { name: 'SpicyLover', isModerator: true },
      message: 'Can you make it extra spicy? 🌶️',
      timestamp: '30s ago',
      type: 'message'
    },
    {
      id: '4',
      user: { name: 'ThaiFoodFan' },
      message: 'Just ordered the Pad Thai kit! 🛒',
      timestamp: '10s ago',
      type: 'system'
    }
  ];

  useEffect(() => {
    setMessages(initialMessages);
    setCurrentViewers(streamData.viewerCount);
    setCurrentLikes(streamData.likes);

    // Simulate live updates
    const interval = setInterval(() => {
      // Add random messages
      if (Math.random() > 0.7) {
        const randomMessages = [
          'This is amazing! 👏',
          'How long does delivery take?',
          'Love your cooking style! ❤️',
          'Can you ship to Phuket?',
          'Best Thai food stream! 🔥',
          'Just followed you! 😊'
        ];
        const randomUsers = [
          { name: 'Thai_Foodie', avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face' },
          { name: 'Bangkok_Lover' },
          { name: 'Spice_Master', isVip: true },
          { name: 'Street_Food_Explorer' }
        ];
        
        const newMsg: ChatMessage = {
          id: Date.now().toString(),
          user: randomUsers[Math.floor(Math.random() * randomUsers.length)],
          message: randomMessages[Math.floor(Math.random() * randomMessages.length)],
          timestamp: 'now',
          type: 'message'
        };
        
        setMessages(prev => [...prev, newMsg].slice(-50)); // Keep last 50 messages
      }

      // Update viewer count
      setCurrentViewers(prev => prev + Math.floor(Math.random() * 10) - 4);
      
      // Update likes occasionally
      if (Math.random() > 0.8) {
        setCurrentLikes(prev => prev + Math.floor(Math.random() * 5) + 1);
      }
    }, 3000);

    return () => clearInterval(interval);
  }, []);

  const handleSendMessage = () => {
    if (newMessage.trim()) {
      const message: ChatMessage = {
        id: Date.now().toString(),
        user: { 
          name: user?.name || 'You', 
          avatar: user?.avatar 
        },
        message: newMessage,
        timestamp: 'now',
        type: 'message'
      };
      
      setMessages(prev => [...prev, message]);
      setNewMessage('');
    }
  };

  const handleSendGift = (gift: Gift) => {
    const message: ChatMessage = {
      id: Date.now().toString(),
      user: { 
        name: user?.name || 'You', 
        avatar: user?.avatar 
      },
      message: `Sent ${gift.name} ${gift.icon}`,
      timestamp: 'now',
      type: 'gift',
      giftType: gift.name
    };
    
    setMessages(prev => [...prev, message]);
    setShowGifts(false);
    toast.success(`Sent ${gift.name} for ฿${gift.price}!`);
  };

  const handleLike = () => {
    // Add heart animation effect
    setCurrentLikes(prev => prev + 1);
    toast.success('❤️');
  };

  const handleFollow = () => {
    setIsFollowing(!isFollowing);
    toast.success(isFollowing ? 'Unfollowed' : 'Following Chef Somsak');
  };

  const handleAddToCart = (productId: string) => {
    if (onAddToCart) {
      onAddToCart(productId);
    }
    toast.success('Added to cart!');
  };

  const handleProductClick = (product: Product) => {
    setSelectedProduct(product);
  };

  const handleShare = () => {
    toast.success('Stream link copied to clipboard!');
  };

  return (
    <div className="h-screen bg-black flex flex-col relative">
      {/* Header Overlay */}
      <div className="absolute top-0 left-0 right-0 z-20 bg-gradient-to-b from-black/80 to-transparent p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="icon" onClick={onBack} className="text-white">
              <ArrowLeft className="w-5 h-5" />
            </Button>
            
            <div className="flex items-center gap-2">
              <Avatar className="w-8 h-8">
                <AvatarImage src={streamData.streamer.avatar} />
                <AvatarFallback>{streamData.streamer.name.charAt(0)}</AvatarFallback>
              </Avatar>
              <div>
                <div className="flex items-center gap-1">
                  <span className="text-white font-medium text-sm">{streamData.streamer.name}</span>
                  {streamData.streamer.verified && (
                    <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                  )}
                </div>
                <div className="flex items-center gap-2 text-xs text-white/80">
                  <span>{streamData.streamer.followers.toLocaleString()} followers</span>
                </div>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant={isFollowing ? "secondary" : "default"}
              size="sm"
              onClick={handleFollow}
              className="text-xs"
            >
              {isFollowing ? 'Following' : 'Follow'}
            </Button>
            <Button variant="ghost" size="icon" className="text-white" onClick={handleShare}>
              <Share className="w-4 h-4" />
            </Button>
            <Button variant="ghost" size="icon" className="text-white" onClick={() => setIsMuted(!isMuted)}>
              {isMuted ? <VolumeX className="w-4 h-4" /> : <Volume2 className="w-4 h-4" />}
            </Button>
          </div>
        </div>

        <div className="mt-3">
          <h2 className="text-white font-medium text-sm mb-2 line-clamp-2">
            {streamData.title}
          </h2>
          <div className="flex items-center gap-4 text-xs text-white/80">
            <Badge className="bg-red-500">
              <div className="w-2 h-2 bg-white rounded-full mr-1 animate-pulse"></div>
              LIVE
            </Badge>
            <div className="flex items-center gap-1">
              <Eye className="w-3 h-3" />
              <span>{currentViewers.toLocaleString()}</span>
            </div>
            <div className="flex items-center gap-1">
              <Heart className="w-3 h-3" />
              <span>{currentLikes.toLocaleString()}</span>
            </div>
            <span>{streamData.duration}</span>
          </div>
        </div>
      </div>

      {/* Video Player Area */}
      <div className="flex-1 relative bg-gray-900 flex items-center justify-center">
        <ImageWithFallback
          src="https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral"
          alt="Live Stream"
          className="w-full h-full object-cover"
        />
        
        {/* Live Indicators */}
        <div className="absolute top-20 right-4 flex flex-col gap-3">
          <Button
            variant="ghost"
            size="icon"
            className="bg-black/50 text-white rounded-full w-12 h-12"
            onClick={handleLike}
          >
            <Heart className="w-6 h-6" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="bg-black/50 text-white rounded-full w-12 h-12"
            onClick={() => setShowGifts(true)}
          >
            <Gift className="w-6 h-6" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="bg-black/50 text-white rounded-full w-12 h-12"
            onClick={() => setShowProducts(true)}
          >
            <ShoppingBag className="w-6 h-6" />
          </Button>
        </div>

        {/* Floating Hearts Animation */}
        <div className="absolute bottom-20 right-4 pointer-events-none">
          {/* Add floating heart animations here */}
        </div>
      </div>

      {/* Bottom Section */}
      <div className="bg-black/90 backdrop-blur-sm border-t border-white/20">
        {/* Products Bar */}
        <div className="p-3 border-b border-white/10">
          <ScrollArea className="w-full">
            <div className="flex gap-3 pb-2">
              {streamData.products.map((product) => (
                <Card 
                  key={product.id} 
                  className="flex-shrink-0 w-32 bg-white/10 border-white/20 cursor-pointer hover:bg-white/20 transition-colors"
                  onClick={() => handleProductClick(product)}
                >
                  <CardContent className="p-2">
                    <ImageWithFallback
                      src={product.image}
                      alt={product.name}
                      className="w-full h-16 object-cover rounded mb-2"
                    />
                    <p className="text-white text-xs font-medium truncate">{product.name}</p>
                    <div className="flex items-center gap-1 mt-1">
                      <span className="text-yellow-400 text-xs font-bold">
                        ฿{product.price}
                      </span>
                      {product.originalPrice && (
                        <span className="text-white/60 text-xs line-through">
                          ฿{product.originalPrice}
                        </span>
                      )}
                    </div>
                    {product.discount && (
                      <Badge className="bg-red-500 text-xs mt-1">
                        -{product.discount}%
                      </Badge>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>
          </ScrollArea>
        </div>

        {/* Chat Section */}
        <div className="flex h-40">
          {/* Messages */}
          <div className="flex-1 flex flex-col">
            <ScrollArea className="flex-1 px-3 py-2">
              <div className="space-y-2">
                {messages.map((message) => (
                  <div key={message.id} className="flex items-start gap-2">
                    {message.user.avatar ? (
                      <Avatar className="w-6 h-6">
                        <AvatarImage src={message.user.avatar} />
                        <AvatarFallback className="text-xs">
                          {message.user.name.charAt(0)}
                        </AvatarFallback>
                      </Avatar>
                    ) : (
                      <div className="w-6 h-6 bg-gradient-to-r from-chart-1 to-chart-2 rounded-full flex items-center justify-center">
                        <span className="text-white text-xs font-bold">
                          {message.user.name.charAt(0)}
                        </span>
                      </div>
                    )}
                    
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-1 mb-1">
                        <span className={`text-xs font-medium ${
                          message.user.isVip ? 'text-yellow-400' : 
                          message.user.isModerator ? 'text-green-400' : 
                          'text-white'
                        }`}>
                          {message.user.name}
                        </span>
                        {message.user.isVip && <Crown className="w-3 h-3 text-yellow-400" />}
                        {message.user.isModerator && <Star className="w-3 h-3 text-green-400" />}
                        {message.type === 'gift' && <Gift className="w-3 h-3 text-pink-400" />}
                      </div>
                      <p className={`text-xs ${
                        message.type === 'gift' ? 'text-pink-300' :
                        message.type === 'system' ? 'text-blue-300' :
                        'text-white/90'
                      }`}>
                        {message.message}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </ScrollArea>

            {/* Message Input */}
            <div className="p-3 border-t border-white/10">
              <div className="flex gap-2">
                <Input
                  placeholder="Say something..."
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
                  className="flex-1 bg-white/10 border-white/20 text-white placeholder:text-white/60"
                />
                <Button size="icon" onClick={handleSendMessage} disabled={!newMessage.trim()}>
                  <Send className="w-4 h-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Gifts Modal */}
      <Dialog open={showGifts} onOpenChange={setShowGifts}>
        <DialogContent className="bg-black/90 border-white/20 text-white">
          <DialogHeader>
            <DialogTitle>Send a Gift</DialogTitle>
          </DialogHeader>
          <div className="grid grid-cols-3 gap-4">
            {gifts.map((gift) => (
              <Button
                key={gift.id}
                variant="outline"
                className="flex flex-col gap-2 h-20 bg-white/10 border-white/20 hover:bg-white/20"
                onClick={() => handleSendGift(gift)}
              >
                <span className="text-2xl">{gift.icon}</span>
                <div className="text-center">
                  <p className="text-xs">{gift.name}</p>
                  <p className="text-xs text-yellow-400">฿{gift.price}</p>
                </div>
              </Button>
            ))}
          </div>
        </DialogContent>
      </Dialog>

      {/* Products Modal */}
      <Dialog open={showProducts} onOpenChange={setShowProducts}>
        <DialogContent className="bg-black/90 border-white/20 text-white max-w-md">
          <DialogHeader>
            <DialogTitle>Featured Products</DialogTitle>
          </DialogHeader>
          <div className="space-y-3">
            {streamData.products.map((product) => (
              <Card key={product.id} className="bg-white/10 border-white/20">
                <CardContent className="p-3">
                  <div className="flex gap-3">
                    <ImageWithFallback
                      src={product.image}
                      alt={product.name}
                      className="w-16 h-16 object-cover rounded"
                    />
                    <div className="flex-1">
                      <h4 className="font-medium text-sm mb-1">{product.name}</h4>
                      <div className="flex items-center gap-2 mb-2">
                        <span className="text-yellow-400 font-bold">฿{product.price}</span>
                        {product.originalPrice && (
                          <span className="text-white/60 text-sm line-through">
                            ฿{product.originalPrice}
                          </span>
                        )}
                      </div>
                      <Button
                        size="sm"
                        onClick={() => handleAddToCart(product.id)}
                        disabled={!product.inStock}
                        className="w-full"
                      >
                        {product.inStock ? 'Add to Cart' : 'Out of Stock'}
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </DialogContent>
      </Dialog>

      {/* Product Detail Modal */}
      <Dialog open={!!selectedProduct} onOpenChange={() => setSelectedProduct(null)}>
        <DialogContent className="bg-black/90 border-white/20 text-white">
          {selectedProduct && (
            <>
              <DialogHeader>
                <DialogTitle>{selectedProduct.name}</DialogTitle>
              </DialogHeader>
              <div className="space-y-4">
                <ImageWithFallback
                  src={selectedProduct.image}
                  alt={selectedProduct.name}
                  className="w-full h-48 object-cover rounded"
                />
                <div className="flex items-center gap-3">
                  <span className="text-2xl font-bold text-yellow-400">
                    ฿{selectedProduct.price}
                  </span>
                  {selectedProduct.originalPrice && (
                    <>
                      <span className="text-lg text-white/60 line-through">
                        ฿{selectedProduct.originalPrice}
                      </span>
                      <Badge className="bg-red-500">
                        -{selectedProduct.discount}% OFF
                      </Badge>
                    </>
                  )}
                </div>
                <div className="flex gap-2">
                  <Button
                    className="flex-1"
                    onClick={() => {
                      handleAddToCart(selectedProduct.id);
                      setSelectedProduct(null);
                    }}
                    disabled={!selectedProduct.inStock}
                  >
                    {selectedProduct.inStock ? 'Add to Cart' : 'Out of Stock'}
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => {
                      onProductClick?.(selectedProduct.id);
                      setSelectedProduct(null);
                    }}
                  >
                    View Details
                  </Button>
                </div>
              </div>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}