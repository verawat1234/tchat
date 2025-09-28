import React, { useState, useRef, useEffect, useMemo } from 'react';
import { ArrowLeft, Phone, Video, MoreVertical, Send, Plus, Camera, Mic, MapPin, Star, Clock, CreditCard, Package, CheckCircle2, Info } from 'lucide-react';
import { useGetChatMessagesQuery } from '../services/microservicesApi';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { toast } from "sonner";
import { motion, AnimatePresence } from 'framer-motion';

interface ShopChatScreenProps {
  user: any;
  shopId: string;
  onBack: () => void;
  onCall?: (shopId: string) => void;
  onVideoCall?: (shopId: string) => void;
  onOrderNow?: (productId: string) => void;
}

interface Message {
  id: string;
  type: 'text' | 'order' | 'payment' | 'location' | 'product' | 'system';
  content: string;
  sender: 'user' | 'shop' | 'system';
  timestamp: Date;
  data?: any;
  status?: 'sent' | 'delivered' | 'read';
}

// Mock shop data - in real app, this would come from API
const getShopData = (shopId: string) => {
  const shops = {
    'golden-mango': {
      id: 'golden-mango',
      name: 'Golden Mango Restaurant',
      avatar: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?w=150&h=150&fit=crop',
      description: 'Authentic Thai street food & traditional dishes',
      rating: 4.8,
      isOnline: true,
      lastSeen: 'Online',
      location: 'Sukhumvit Rd, Bangkok',
      specialties: ['Pad Thai', 'Green Curry', 'Mango Sticky Rice'],
      paymentMethods: ['PromptPay', 'Cash', 'Card']
    },
    'thai-coffee': {
      id: 'thai-coffee',
      name: 'Thai Coffee Chain',
      avatar: 'https://images.unsplash.com/photo-1559847844-5315695b6a77?w=150&h=150&fit=crop',
      description: 'Premium Thai coffee & traditional beverages',
      rating: 4.6,
      isOnline: true,
      lastSeen: 'Online',
      location: 'Siam Square, Bangkok',
      specialties: ['Thai Iced Coffee', 'Matcha Latte', 'Thai Tea'],
      paymentMethods: ['PromptPay', 'True Wallet', 'Card']
    },
    'street-food': {
      id: 'street-food',
      name: 'Somtam Street Cart',
      avatar: 'https://images.unsplash.com/photo-1590947132387-155cc02f3212?w=150&h=150&fit=crop',
      description: 'Fresh papaya salad & Isaan specialties',
      rating: 4.9,
      isOnline: false,
      lastSeen: '2 minutes ago',
      location: 'Chatuchak Market, Bangkok',
      specialties: ['Som Tam', 'Grilled Chicken', 'Sticky Rice'],
      paymentMethods: ['PromptPay', 'Cash']
    }
  };
  return shops[shopId] || shops['golden-mango'];
};

// Mock messages - in real app, these would come from API
const getInitialMessages = (shopId: string): Message[] => {
  const baseMessages = [
    {
      id: '1',
      type: 'system' as const,
      content: 'Chat encryption enabled. Messages are secure.',
      sender: 'system' as const,
      timestamp: new Date(Date.now() - 86400000),
      status: 'read' as const
    },
    {
      id: '2',
      type: 'text' as const,
      content: 'Hello! Welcome to our restaurant. How can I help you today? 🍜',
      sender: 'shop' as const,
      timestamp: new Date(Date.now() - 7200000),
      status: 'read' as const
    },
    {
      id: '3',
      type: 'text' as const,
      content: "Hi! I'd like to know about your Pad Thai. Is it available today?",
      sender: 'user' as const,
      timestamp: new Date(Date.now() - 7100000),
      status: 'read' as const
    },
    {
      id: '4',
      type: 'text' as const,
      content: "Yes! Our signature Pad Thai is available. It comes with fresh prawns, tofu, bean sprouts, and our special tamarind sauce. ฿120",
      sender: 'shop' as const,
      timestamp: new Date(Date.now() - 7000000),
      status: 'read' as const
    },
    {
      id: '5',
      type: 'product' as const,
      content: 'Signature Pad Thai',
      sender: 'shop' as const,
      timestamp: new Date(Date.now() - 6900000),
      data: {
        name: 'Signature Pad Thai',
        price: 120,
        currency: 'THB',
        image: 'https://images.unsplash.com/photo-1559314809-0f31657403cb?w=300&h=200&fit=crop',
        description: 'Traditional Thai stir-fried noodles with prawns, tofu, and special sauce'
      },
      status: 'read' as const
    }
  ];

  // Add shop-specific messages
  if (shopId === 'thai-coffee') {
    return [
      ...baseMessages.slice(0, 2),
      {
        id: '3',
        type: 'text' as const,
        content: "Hi! What's your most popular coffee today?",
        sender: 'user' as const,
        timestamp: new Date(Date.now() - 7100000),
        status: 'read' as const
      },
      {
        id: '4',
        type: 'text' as const,
        content: "Our Thai Iced Coffee is very popular! Made with locally sourced beans and traditional brewing method. ฿65",
        sender: 'shop' as const,
        timestamp: new Date(Date.now() - 7000000),
        status: 'read' as const
      },
      {
        id: '5',
        type: 'product' as const,
        content: 'Thai Iced Coffee',
        sender: 'shop' as const,
        timestamp: new Date(Date.now() - 6900000),
        data: {
          name: 'Thai Iced Coffee',
          price: 65,
          currency: 'THB',
          image: 'https://images.unsplash.com/photo-1517701604599-bb29b565090c?w=300&h=200&fit=crop',
          description: 'Traditional Thai coffee with condensed milk and ice'
        },
        status: 'read' as const
      }
    ];
  }

  return baseMessages;
};

export function ShopChatScreen({ user, shopId, onBack, onCall, onVideoCall, onOrderNow }: ShopChatScreenProps) {
  const [inputValue, setInputValue] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // RTK Query for chat messages
  const {
    data: messagesData,
    isLoading: messagesLoading,
    error: messagesError
  } = useGetChatMessagesQuery({
    chatId: shopId,
    page: 1,
    limit: 50
  });

  const messages: Message[] = useMemo(() => {
    if (messagesLoading || !messagesData?.data) {
      // Fallback data while loading
      return getInitialMessages(shopId);
    }

    return messagesData.data.map((msg: any) => ({
      id: msg.id || msg.message_id || `msg-${Math.random()}`,
      type: msg.type || msg.message_type || 'text',
      content: msg.content || msg.message || msg.text || '',
      sender: msg.sender || msg.sender_type || (msg.from_user ? 'user' : 'shop'),
      timestamp: new Date(msg.timestamp || msg.created_at || Date.now()),
      data: msg.data || msg.metadata || undefined,
      status: msg.status || msg.message_status || 'sent'
    }));
  }, [messagesData, messagesLoading, shopId]);

  const [localMessages, setLocalMessages] = useState<Message[]>([]);
  const allMessages = [...messages, ...localMessages];

  const shopData = getShopData(shopId);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (scrollAreaRef.current) {
      const scrollContainer = scrollAreaRef.current.querySelector('[data-radix-scroll-area-viewport]');
      if (scrollContainer) {
        scrollContainer.scrollTop = scrollContainer.scrollHeight;
      }
    }
  }, [allMessages]);

  const sendMessage = () => {
    if (!inputValue.trim()) return;

    const newMessage: Message = {
      id: Date.now().toString(),
      type: 'text',
      content: inputValue,
      sender: 'user',
      timestamp: new Date(),
      status: 'sent'
    };

    setLocalMessages(prev => [...prev, newMessage]);
    setInputValue('');

    // Simulate shop typing and response
    setTimeout(() => {
      setIsTyping(true);
    }, 500);

    setTimeout(() => {
      setIsTyping(false);
      const responses = [
        "Thanks for your message! I'll check that for you.",
        "That sounds great! Let me help you with that.",
        "Sure thing! Our chef will prepare that fresh for you.",
        "Perfect choice! Would you like to add anything else?",
        "Coming right up! Estimated time: 15-20 minutes."
      ];
      
      const shopResponse: Message = {
        id: (Date.now() + 1).toString(),
        type: 'text',
        content: responses[Math.floor(Math.random() * responses.length)],
        sender: 'shop',
        timestamp: new Date(),
        status: 'sent'
      };
      
      setLocalMessages(prev => [...prev, shopResponse]);
    }, 2000);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const sendQuickAction = (type: string, content: string, data?: any) => {
    const newMessage: Message = {
      id: Date.now().toString(),
      type: type as any,
      content,
      sender: 'user',
      timestamp: new Date(),
      data,
      status: 'sent'
    };

    setLocalMessages(prev => [...prev, newMessage]);
    
    // Quick responses for different actions
    setTimeout(() => {
      let response = '';
      switch (type) {
        case 'location':
          response = "Thanks! I can see your location. Delivery time will be about 25-30 minutes.";
          break;
        case 'payment':
          response = "Payment confirmed! Your order is being prepared. 🍽️";
          break;
        default:
          response = "Thanks for that!";
      }
      
      const shopResponse: Message = {
        id: (Date.now() + 1).toString(),
        type: 'text',
        content: response,
        sender: 'shop',
        timestamp: new Date(),
        status: 'sent'
      };
      
      setLocalMessages(prev => [...prev, shopResponse]);
    }, 1000);
  };

  const renderMessage = (message: Message) => {
    const isUser = message.sender === 'user';
    const isSystem = message.sender === 'system';

    if (isSystem) {
      return (
        <div key={message.id} className="flex justify-center my-4">
          <div className="bg-muted/50 text-muted-foreground text-xs px-3 py-1 rounded-full flex items-center gap-1">
            <Info className="w-3 h-3" />
            {message.content}
          </div>
        </div>
      );
    }

    return (
      <motion.div
        key={message.id}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className={`flex gap-2 mb-4 ${isUser ? 'flex-row-reverse' : 'flex-row'}`}
      >
        {!isUser && (
          <Avatar className="w-8 h-8 flex-shrink-0">
            <AvatarImage src={shopData.avatar} />
            <AvatarFallback>{shopData.name.charAt(0)}</AvatarFallback>
          </Avatar>
        )}
        
        <div className={`flex flex-col max-w-[80%] ${isUser ? 'items-end' : 'items-start'}`}>
          {message.type === 'product' && message.data ? (
            <Card className={`${isUser ? 'bg-primary text-primary-foreground' : 'bg-muted'} border-0 shadow-sm`}>
              <CardContent className="p-3">
                <div className="flex gap-3">
                  <img 
                    src={message.data.image} 
                    alt={message.data.name}
                    className="w-16 h-16 rounded-lg object-cover flex-shrink-0"
                  />
                  <div className="flex-1 min-w-0">
                    <h4 className="font-medium text-sm mb-1">{message.data.name}</h4>
                    <p className="text-xs opacity-80 mb-2 line-clamp-2">{message.data.description}</p>
                    <div className="flex items-center justify-between">
                      <span className="font-semibold">฿{message.data.price}</span>
                      <Button 
                        size="sm" 
                        variant={isUser ? "secondary" : "default"}
                        className="h-6 px-2 text-xs"
                        onClick={() => onOrderNow?.(message.data.id)}
                      >
                        Order Now
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ) : message.type === 'order' && message.data ? (
            <Card className={`${isUser ? 'bg-primary text-primary-foreground' : 'bg-muted'} border-0 shadow-sm`}>
              <CardContent className="p-3">
                <div className="flex items-center gap-2 mb-2">
                  <Package className="w-4 h-4" />
                  <span className="font-medium text-sm">Order #{message.data.orderId}</span>
                  <Badge variant="secondary" className="text-xs">
                    {message.data.status}
                  </Badge>
                </div>
                <p className="text-xs opacity-80">{message.content}</p>
              </CardContent>
            </Card>
          ) : (
            <div className={`px-3 py-2 rounded-2xl ${
              isUser 
                ? 'bg-primary text-primary-foreground rounded-br-md' 
                : 'bg-muted rounded-bl-md'
            }`}>
              <p className="text-sm whitespace-pre-wrap">{message.content}</p>
            </div>
          )}
          
          <div className="flex items-center gap-1 mt-1">
            <span className="text-xs text-muted-foreground">
              {message.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
            </span>
            {isUser && message.status && (
              <div className="flex items-center">
                {message.status === 'sent' && <div className="w-1 h-1 bg-muted-foreground rounded-full" />}
                {message.status === 'delivered' && <CheckCircle2 className="w-3 h-3 text-muted-foreground" />}
                {message.status === 'read' && <CheckCircle2 className="w-3 h-3 text-chart-1" />}
              </div>
            )}
          </div>
        </div>
      </motion.div>
    );
  };

  return (
    <div className="h-full flex flex-col bg-background">
      {/* Chat Header */}
      <motion.header 
        className="flex items-center gap-3 p-4 border-b border-border bg-card/95 backdrop-blur-sm"
        initial={{ y: -20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Button variant="ghost" size="icon" onClick={onBack} className="flex-shrink-0">
          <ArrowLeft className="w-5 h-5" />
        </Button>
        
        <Avatar className="w-10 h-10 flex-shrink-0">
          <AvatarImage src={shopData.avatar} />
          <AvatarFallback>{shopData.name.charAt(0)}</AvatarFallback>
        </Avatar>
        
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="font-medium truncate">{shopData.name}</h3>
            <div className="flex items-center gap-1">
              <Star className="w-3 h-3 fill-yellow-500 text-yellow-500" />
              <span className="text-xs text-muted-foreground">{shopData.rating}</span>
            </div>
          </div>
          <div className="flex items-center gap-2">
            {shopData.isOnline ? (
              <div className="flex items-center gap-1">
                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                <span className="text-xs text-muted-foreground">Online</span>
              </div>
            ) : (
              <div className="flex items-center gap-1">
                <Clock className="w-3 h-3 text-muted-foreground" />
                <span className="text-xs text-muted-foreground">{shopData.lastSeen}</span>
              </div>
            )}
            <div className="flex items-center gap-1">
              <MapPin className="w-3 h-3 text-muted-foreground" />
              <span className="text-xs text-muted-foreground truncate">{shopData.location}</span>
            </div>
          </div>
        </div>
        
        <div className="flex items-center gap-1 flex-shrink-0">
          <Button 
            variant="ghost" 
            size="icon" 
            className="w-9 h-9"
            onClick={() => onCall?.(shopId)}
          >
            <Phone className="w-4 h-4" />
          </Button>
          <Button 
            variant="ghost" 
            size="icon" 
            className="w-9 h-9"
            onClick={() => onVideoCall?.(shopId)}
          >
            <Video className="w-4 h-4" />
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="w-9 h-9">
                <MoreVertical className="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => toast.success("Shop info opened")}>
                Shop Information
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => toast.success("Menu opened")}>
                View Menu
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => toast.success("Order history opened")}>
                Order History
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => toast.success("Report shop")}>
                Report Shop
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </motion.header>

      {/* Quick Actions Bar */}
      <motion.div 
        className="p-3 border-b border-border bg-muted/30"
        initial={{ y: -10, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.1 }}
      >
        <div className="flex gap-2 overflow-x-auto scrollbar-hide">
          <Button 
            variant="outline" 
            size="sm" 
            className="flex-shrink-0 gap-1"
            onClick={() => sendQuickAction('text', "Can I see your menu please?")}
          >
            <Package className="w-3 h-3" />
            Menu
          </Button>
          <Button 
            variant="outline" 
            size="sm" 
            className="flex-shrink-0 gap-1"
            onClick={() => sendQuickAction('location', "Here is my location for delivery", { lat: 13.7563, lng: 100.5018 })}
          >
            <MapPin className="w-3 h-3" />
            Share Location
          </Button>
          <Button 
            variant="outline" 
            size="sm" 
            className="flex-shrink-0 gap-1"
            onClick={() => toast.success("Opening PromptPay QR...")}
          >
            <CreditCard className="w-3 h-3" />
            PromptPay
          </Button>
          <Button 
            variant="outline" 
            size="sm" 
            className="flex-shrink-0 gap-1"
            onClick={() => sendQuickAction('text', "What are your most popular dishes?")}
          >
            <Star className="w-3 h-3" />
            Popular
          </Button>
        </div>
      </motion.div>

      {/* Messages Area */}
      <ScrollArea className="flex-1 px-4" ref={scrollAreaRef}>
        <div className="py-4">
          <AnimatePresence>
            {allMessages.map(renderMessage)}
            
            {/* Typing Indicator */}
            {isTyping && (
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                className="flex gap-2 mb-4"
              >
                <Avatar className="w-8 h-8 flex-shrink-0">
                  <AvatarImage src={shopData.avatar} />
                  <AvatarFallback>{shopData.name.charAt(0)}</AvatarFallback>
                </Avatar>
                <div className="bg-muted px-3 py-2 rounded-2xl rounded-bl-md">
                  <div className="flex gap-1">
                    <div className="w-2 h-2 bg-muted-foreground/50 rounded-full animate-bounce"></div>
                    <div className="w-2 h-2 bg-muted-foreground/50 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                    <div className="w-2 h-2 bg-muted-foreground/50 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                  </div>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </ScrollArea>

      {/* Chat Input */}
      <motion.div 
        className="p-4 border-t border-border bg-card/95 backdrop-blur-sm"
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.2 }}
      >
        <div className="flex items-end gap-2">
          <Button variant="ghost" size="icon" className="flex-shrink-0 mb-1">
            <Plus className="w-5 h-5" />
          </Button>
          
          <div className="flex-1 relative">
            <Input
              ref={inputRef}
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder={`Message ${shopData.name}...`}
              className="pr-20 resize-none rounded-full bg-input-background"
              rows={1}
            />
            <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1">
              <Button variant="ghost" size="icon" className="w-8 h-8">
                <Camera className="w-4 h-4" />
              </Button>
              <Button variant="ghost" size="icon" className="w-8 h-8">
                <Mic className="w-4 h-4" />
              </Button>
            </div>
          </div>
          
          <Button 
            onClick={sendMessage}
            disabled={!inputValue.trim()}
            className="flex-shrink-0 rounded-full w-10 h-10 p-0"
          >
            <Send className="w-4 h-4" />
          </Button>
        </div>
        
        {/* Payment Methods */}
        <div className="flex items-center gap-2 mt-2 pt-2 border-t border-border/50">
          <span className="text-xs text-muted-foreground">Accepts:</span>
          {shopData.paymentMethods.map((method) => (
            <Badge key={method} variant="outline" className="text-xs">
              {method === 'PromptPay' && '🇹🇭'} {method}
            </Badge>
          ))}
        </div>
      </motion.div>
    </div>
  );
}