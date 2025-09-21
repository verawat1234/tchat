import React, { useState, useEffect, useRef, useCallback } from 'react';
import { MessageCircle, Store, Users, Bell, Settings, Search, Plus, QrCode, Wallet, Globe, CreditCard, Play, Heart, User, TrendingUp, Briefcase, ShoppingCart, Languages, Banknote, ArrowDown, Zap } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { Button } from './components/ui/button';
import { Badge } from './components/ui/badge';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from './components/ui/dropdown-menu';
import { Avatar, AvatarFallback, AvatarImage } from './components/ui/avatar';
import { ScrollArea } from './components/ui/scroll-area';
import { AuthScreen } from './components/AuthScreen';
import { RichChatTab } from './components/RichChatTab';
import { StoreTab } from './components/StoreTab';
import { WorkspaceSwitcher } from './components/WorkspaceSwitcher';
import { WorkspaceTab } from './components/WorkspaceTab';
import { SocialTab } from './components/SocialTab';
import { SettingsScreen } from './components/SettingsScreen';
import { WalletScreen } from './components/WalletScreen';
import { QRScannerScreen } from './components/QRScannerScreen';
import { SearchScreen } from './components/SearchScreen';
import { VideoCallScreen } from './components/VideoCallScreen';
import { VoiceCallScreen } from './components/VoiceCallScreen';
import { NewChatScreen } from './components/NewChatScreen';
import { ShopPage } from './components/ShopPage';
import { ProductPage } from './components/ProductPage';
import { LiveStreamScreen } from './components/LiveStreamScreen';
import { ShopChatScreen } from './components/ShopChatScreen';
import { VideoTab } from './components/VideoTab';
import { CartScreen } from './components/CartScreen';
import { ShareDialog } from './components/ShareDialog';
import { FullscreenVideoPlayer } from './components/FullscreenVideoPlayer';
import { NotificationsScreen } from './components/NotificationsScreen';
import { Toaster } from './components/ui/sonner';
import { toast } from "sonner";

type Screen = 'chat' | 'store' | 'social' | 'video' | 'work' | 'more' | 'settings' | 'wallet' | 'cart' | 'qr-scanner' | 'search' | 'video-call' | 'voice-call' | 'new-chat' | 'shop' | 'product' | 'live-stream' | 'notifications' | 'shop-chat';

interface Notification {
  id: string;
  type: 'message' | 'payment' | 'social' | 'merchant' | 'system';
  title: string;
  description: string;
  timestamp: string;
  read: boolean;
  avatar?: string;
  amount?: number;
}

export default function App() {
  const [currentScreen, setCurrentScreen] = useState<Screen>('chat');
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<any>(null);
  const [currentLanguage, setCurrentLanguage] = useState('th');
  const [currentCurrency, setCurrentCurrency] = useState('THB');
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState<string>('golden-mango-shop');
  const [isKeyboardOpen, setIsKeyboardOpen] = useState(false);
  const [viewportHeight, setViewportHeight] = useState<number>(0);
  
  // Touch/Swipe handling
  const touchStartX = useRef<number>(0);
  const touchStartY = useRef<number>(0);
  const mainContentRef = useRef<HTMLDivElement>(null);
  
  const [selectedShopId, setSelectedShopId] = useState<string>('');
  const [selectedProductId, setSelectedProductId] = useState<string>('');
  const [selectedStreamId, setSelectedStreamId] = useState<string>('');
  const [watchedVideos, setWatchedVideos] = useState<string[]>([]);
  const [likedVideos, setLikedVideos] = useState<string[]>([]);
  const [currentVideoId, setCurrentVideoId] = useState<string>('');
  const [cartItems, setCartItems] = useState<any[]>([]);
  const [subscribedChannels, setSubscribedChannels] = useState<string[]>(['thai-chef', 'thai-culture']);
  const [fullscreenVideo, setFullscreenVideo] = useState<{
    video: any;
    isPlaying: boolean;
  } | null>(null);
  const [callData, setCallData] = useState<{
    callee: {
      id: string;
      name: string;
      avatar?: string;
      isGroup?: boolean;
      members?: number;
    };
    isIncoming: boolean;
  } | null>(null);
  const [shareDialog, setShareDialog] = useState<{
    open: boolean;
    content: {
      type: 'post' | 'video' | 'product' | 'live-stream' | 'shop' | 'workspace-file';
      id: string;
      title: string;
      description?: string;
      image?: string;
      author?: {
        name: string;
        avatar?: string;
      };
      price?: {
        amount: number;
        currency: string;
      };
      metadata?: any;
    } | null;
  }>({
    open: false,
    content: null
  });
  const [notifications, setNotifications] = useState<Notification[]>([
    {
      id: '1',
      type: 'message',
      title: 'Family Group',
      description: 'Mom: Dinner at 7pm! 🍽️',
      timestamp: '2 min ago',
      read: false,
      avatar: ''
    },
    {
      id: '2',
      type: 'payment',
      title: 'Payment Received',
      description: 'Received ฿500 from Mom via PromptPay',
      timestamp: '1 hour ago',
      read: false,
      amount: 500
    },
    {
      id: '3',
      type: 'social',
      title: 'New Like',
      description: 'John liked your street food photo',
      timestamp: '3 hours ago',
      read: false,
      avatar: ''
    },
    {
      id: '4',
      type: 'merchant',
      title: 'Order Update',
      description: 'Your Pad Thai is ready for pickup!',
      timestamp: '5 hours ago',
      read: true,
      avatar: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'
    },
    {
      id: '5',
      type: 'system',
      title: 'Songkran Festival Sale!',
      description: 'Up to 50% off on Thai street food vendors',
      timestamp: '1 day ago',
      read: true
    }
  ]);

  // Language and Currency Data
  const languages = [
    { code: 'th', name: 'ไทย', flag: '🇹🇭' },
    { code: 'id', name: 'Indonesia', flag: '🇮🇩' },
    { code: 'vi', name: 'Việt Nam', flag: '🇻🇳' },
    { code: 'ms', name: 'Malaysia', flag: '🇲🇾' },
    { code: 'tl', name: 'Philippines', flag: '🇵🇭' },
    { code: 'en', name: 'English', flag: '🇺🇸' }
  ];

  const currencies = [
    { code: 'THB', symbol: '฿', name: 'Thai Baht' },
    { code: 'IDR', symbol: 'Rp', name: 'Indonesian Rupiah' },
    { code: 'VND', symbol: '₫', name: 'Vietnamese Dong' },
    { code: 'MYR', symbol: 'RM', name: 'Malaysian Ringgit' },
    { code: 'PHP', symbol: '₱', name: 'Philippine Peso' },
    { code: 'SGD', symbol: 'S$', name: 'Singapore Dollar' }
  ];

  const qrPaymentMethods = [
    { name: 'PromptPay', country: 'Thailand', icon: '🇹🇭', action: () => toast.success('Opening PromptPay...') },
    { name: 'QRIS', country: 'Indonesia', icon: '🇮🇩', action: () => toast.success('Opening QRIS...') },
    { name: 'VietQR', country: 'Vietnam', icon: '🇻🇳', action: () => toast.success('Opening VietQR...') },
    { name: 'DuitNow', country: 'Malaysia', icon: '🇲🇾', action: () => toast.success('Opening DuitNow QR...') },
    { name: 'InstaPay', country: 'Philippines', icon: '🇵🇭', action: () => toast.success('Opening InstaPay QR...') }
  ];

  // User's shops and workplaces
  const userWorkspaces = [
    {
      id: 'golden-mango-shop',
      name: 'Golden Mango Restaurant',
      type: 'shop',
      role: 'Owner',
      avatar: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?w=150&h=150&fit=crop',
      customerCount: 245,
      revenue: '฿125K'
    },
    {
      id: 'thai-coffee-chain',
      name: 'Thai Coffee Chain',
      type: 'shop',
      role: 'Manager',
      avatar: 'https://images.unsplash.com/photo-1559847844-5315695b6a77?w=150&h=150&fit=crop',
      customerCount: 89,
      revenue: '฿67K'
    },
    {
      id: 'sea-tech-company',
      name: 'SEA Tech Solutions',
      type: 'company',
      role: 'Customer Support Lead',
      avatar: 'https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d?w=150&h=150&fit=crop',
      customerCount: 156,
      revenue: null // Not revenue-focused for company work
    }
  ];

  // Tab management for swipe navigation
  const tabs: Screen[] = ['chat', 'store', 'social', 'video', 'more'];
  const currentTabIndex = tabs.indexOf(currentScreen);

  // Handle swipe gestures
  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    touchStartX.current = e.touches[0].clientX;
    touchStartY.current = e.touches[0].clientY;
  }, []);

  const handleTouchEnd = useCallback((e: React.TouchEvent) => {
    if (!touchStartX.current || !touchStartY.current) return;

    const touchEndX = e.changedTouches[0].clientX;
    const touchEndY = e.changedTouches[0].clientY;
    const deltaX = touchStartX.current - touchEndX;
    const deltaY = touchStartY.current - touchEndY;

    // Only handle horizontal swipes (ignore vertical scrolling)
    if (Math.abs(deltaX) > Math.abs(deltaY) && Math.abs(deltaX) > 50) {
      if (deltaX > 0 && currentTabIndex < tabs.length - 1) {
        // Swipe left - next tab
        setCurrentScreen(tabs[currentTabIndex + 1]);
        // Add haptic feedback if available
        if (navigator.vibrate) navigator.vibrate(50);
      } else if (deltaX < 0 && currentTabIndex > 0) {
        // Swipe right - previous tab
        setCurrentScreen(tabs[currentTabIndex - 1]);
        if (navigator.vibrate) navigator.vibrate(50);
      }
    }

    // Pull to refresh (swipe down at top)
    if (deltaY < -100 && touchStartY.current < 100 && !isRefreshing) {
      handlePullToRefresh();
    }

    touchStartX.current = 0;
    touchStartY.current = 0;
  }, [currentTabIndex, isRefreshing]);

  // Pull to refresh functionality
  const handlePullToRefresh = useCallback(async () => {
    if (isRefreshing) return;
    
    setIsRefreshing(true);
    if (navigator.vibrate) navigator.vibrate(100);
    
    // Simulate refresh delay
    await new Promise(resolve => setTimeout(resolve, 1500));
    
    toast.success('Content refreshed!');
    setIsRefreshing(false);
  }, [isRefreshing]);

  // Animation variants for screen transitions
  const screenVariants = {
    initial: { 
      opacity: 0, 
      x: 20,
      scale: 0.98
    },
    animate: { 
      opacity: 1, 
      x: 0,
      scale: 1,
      transition: {
        duration: 0.3,
        ease: [0.22, 1, 0.36, 1]
      }
    },
    exit: { 
      opacity: 0, 
      x: -20,
      scale: 0.98,
      transition: {
        duration: 0.2,
        ease: [0.22, 1, 0.36, 1]
      }
    }
  };

  const slideVariants = {
    initial: (direction: number) => ({
      x: direction > 0 ? 300 : -300,
      opacity: 0
    }),
    animate: {
      x: 0,
      opacity: 1,
      transition: {
        type: "spring",
        stiffness: 300,
        damping: 30
      }
    },
    exit: (direction: number) => ({
      x: direction < 0 ? 300 : -300,
      opacity: 0,
      transition: {
        type: "spring",
        stiffness: 300,
        damping: 30
      }
    })
  };

  const modalVariants = {
    initial: {
      opacity: 0,
      scale: 0.95,
      y: 10
    },
    animate: {
      opacity: 1,
      scale: 1,
      y: 0,
      transition: {
        duration: 0.2,
        ease: [0.22, 1, 0.36, 1]
      }
    },
    exit: {
      opacity: 0,
      scale: 0.95,
      y: 10,
      transition: {
        duration: 0.15,
        ease: [0.22, 1, 0.36, 1]
      }
    }
  };

  const staggerChildrenVariants = {
    animate: {
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const childVariants = {
    initial: { opacity: 0, y: 20 },
    animate: { 
      opacity: 1, 
      y: 0,
      transition: {
        duration: 0.3,
        ease: [0.22, 1, 0.36, 1]
      }
    }
  };

  // Language switching
  const handleLanguageSwitch = useCallback((langCode: string) => {
    setCurrentLanguage(langCode);
    const lang = languages.find(l => l.code === langCode);
    toast.success(`Language switched to ${lang?.name}`);
  }, []);

  // Currency switching  
  const handleCurrencySwitch = useCallback((currencyCode: string) => {
    setCurrentCurrency(currencyCode);
    const currency = currencies.find(c => c.code === currencyCode);
    toast.success(`Currency changed to ${currency?.name}`);
  }, []);

  // Simulate checking auth state
  useEffect(() => {
    const savedAuth = localStorage.getItem('telegram-sea-auth');
    if (savedAuth) {
      const authData = JSON.parse(savedAuth);
      setIsAuthenticated(true);
      setUser(authData.user);
    }
  }, []);

  // Enhanced keyboard visibility for mobile layout optimization
  useEffect(() => {
    if (typeof window === 'undefined') return;

    let initialViewportHeight = window.visualViewport?.height || window.innerHeight;
    setViewportHeight(initialViewportHeight);
    
    const handleViewportChange = () => {
      const currentHeight = window.visualViewport?.height || window.innerHeight;
      const heightDifference = initialViewportHeight - currentHeight;
      
      setViewportHeight(currentHeight);
      // If viewport shrunk by more than 150px, likely keyboard is open
      setIsKeyboardOpen(heightDifference > 150);
    };

    // Use visualViewport API if available (iOS Safari)
    if (window.visualViewport) {
      window.visualViewport.addEventListener('resize', handleViewportChange);
      return () => {
        window.visualViewport?.removeEventListener('resize', handleViewportChange);
      };
    } else {
      // Fallback for other browsers
      window.addEventListener('resize', handleViewportChange);
      return () => {
        window.removeEventListener('resize', handleViewportChange);
      };
    }
  }, []);

  // Simulate incoming calls
  useEffect(() => {
    if (isAuthenticated) {
      // Random incoming calls for demo
      const callTimer = setTimeout(() => {
        if (Math.random() > 0.7 && currentScreen === 'chat') {
          const isVideo = Math.random() > 0.5;
          const callee = {
            id: 'demo-caller',
            name: 'Mom',
            avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face',
            isGroup: false
          };
          
          if (isVideo) {
            startVideoCall(callee, true);
          } else {
            startVoiceCall(callee, true);
          }
        }
      }, 10000); // 10 seconds after auth

      return () => clearTimeout(callTimer);
    }
  }, [isAuthenticated, currentScreen]);

  const handleAuth = (userData: any) => {
    setIsAuthenticated(true);
    setUser(userData);
    localStorage.setItem('telegram-sea-auth', JSON.stringify({ user: userData }));
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    setUser(null);
    setCurrentScreen('chat');
    localStorage.removeItem('telegram-sea-auth');
  };

  const handleBackToMain = () => {
    setCurrentScreen('chat');
  };

  const startVideoCall = (callee: any, isIncoming = false) => {
    setCallData({ callee, isIncoming });
    setCurrentScreen('video-call');
  };

  const startVoiceCall = (callee: any, isIncoming = false) => {
    setCallData({ callee, isIncoming });
    setCurrentScreen('voice-call');
  };

  const endCall = () => {
    setCallData(null);
    setCurrentScreen('chat');
  };

  const handleCreateChat = (chatData: any) => {
    // In real app, would create new chat via API
    console.log('Creating new chat:', chatData);
    setCurrentScreen('chat');
    toast.success(`Started chat with ${chatData.name}`);
  };

  const handleShopClick = (shopId: string) => {
    setSelectedShopId(shopId);
    setCurrentScreen('shop');
  };

  const handleShopChatClick = (shopId: string) => {
    setSelectedShopId(shopId);
    setCurrentScreen('shop-chat');
  };

  const handleProductClick = (productId: string) => {
    setSelectedProductId(productId);
    setCurrentScreen('product');
  };

  const handleVideoPlay = (videoId: string, videoData?: any) => {
    setCurrentVideoId(videoId);
    setWatchedVideos(prev => [...new Set([...prev, videoId])]);
    
    // Open fullscreen video player if video data is provided
    if (videoData) {
      setFullscreenVideo({
        video: videoData,
        isPlaying: true
      });
    }
  };

  const handleFullscreenVideoClose = () => {
    setFullscreenVideo(null);
  };

  const handleFullscreenVideoPlay = () => {
    if (fullscreenVideo) {
      setFullscreenVideo(prev => prev ? { ...prev, isPlaying: true } : null);
    }
  };

  const handleFullscreenVideoPause = () => {
    if (fullscreenVideo) {
      setFullscreenVideo(prev => prev ? { ...prev, isPlaying: false } : null);
    }
  };

  const handleVideoLike = (videoId: string) => {
    setLikedVideos(prev => 
      prev.includes(videoId) 
        ? prev.filter(id => id !== videoId)
        : [...prev, videoId]
    );
    toast.success(likedVideos.includes(videoId) ? 'Video unliked' : 'Video liked!');
  };

  const handleVideoShare = (videoId: string, videoData?: any) => {
    // Open comprehensive share dialog for video
    setShareDialog({
      open: true,
      content: {
        type: 'video',
        id: videoId,
        title: videoData?.title || `Video ${videoId}`,
        description: videoData?.description,
        image: videoData?.thumbnail,
        author: videoData?.author,
        metadata: videoData
      }
    });
  };

  const handlePostShare = (postId: string, postData: any) => {
    setShareDialog({
      open: true,
      content: {
        type: 'post',
        id: postId,
        title: postData.content?.slice(0, 50) + (postData.content?.length > 50 ? '...' : ''),
        description: postData.content,
        image: postData.images?.[0],
        author: postData.author,
        metadata: postData
      }
    });
  };

  const handleProductShare = (productId: string, productData: any) => {
    setShareDialog({
      open: true,
      content: {
        type: 'product',
        id: productId,
        title: productData.name,
        description: productData.description,
        image: productData.image,
        price: {
          amount: productData.price,
          currency: productData.currency || 'THB'
        },
        author: productData.vendor,
        metadata: productData
      }
    });
  };

  const handleShopShare = (shopId: string, shopData: any) => {
    setShareDialog({
      open: true,
      content: {
        type: 'shop',
        id: shopId,
        title: shopData.name,
        description: shopData.description,
        image: shopData.image,
        author: {
          name: shopData.name,
          avatar: shopData.avatar
        },
        metadata: shopData
      }
    });
  };

  const handleLiveStreamShare = (streamId: string, streamData: any) => {
    setShareDialog({
      open: true,
      content: {
        type: 'live-stream',
        id: streamId,
        title: streamData.title,
        description: `Live stream by ${streamData.author?.name}`,
        image: streamData.thumbnail,
        author: streamData.author,
        metadata: streamData
      }
    });
  };

  const handleWorkspaceFileShare = (fileId: string, fileData: any) => {
    setShareDialog({
      open: true,
      content: {
        type: 'workspace-file',
        id: fileId,
        title: fileData.name,
        description: `${fileData.type} file - ${fileData.size || 'Unknown size'}`,
        author: {
          name: fileData.owner,
          avatar: fileData.ownerAvatar
        },
        metadata: fileData
      }
    });
  };

  const handleVideoSubscribe = (channelId: string) => {
    if (subscribedChannels.includes(channelId)) {
      setSubscribedChannels(prev => prev.filter(id => id !== channelId));
      toast.success('Unsubscribed from channel');
    } else {
      setSubscribedChannels(prev => [...prev, channelId]);
      toast.success('Subscribed to channel!');
    }
  };

  const handleLiveStreamClick = (streamId: string) => {
    setSelectedStreamId(streamId);
    setCurrentScreen('live-stream');
  };

  const handleNotificationClick = (notificationId: string) => {
    setNotifications(prev => 
      prev.map(notification => 
        notification.id === notificationId 
          ? { ...notification, read: true }
          : notification
      )
    );
  };

  const markAllAsRead = () => {
    setNotifications(prev => 
      prev.map(notification => ({ ...notification, read: true }))
    );
    toast.success('All notifications marked as read');
  };

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'message':
        return <MessageCircle className="w-5 h-5 text-chart-1" />;
      case 'payment':
        return <CreditCard className="w-5 h-5 text-green-500" />;
      case 'social':
        return <Heart className="w-5 h-5 text-red-500" />;
      case 'merchant':
        return <Store className="w-5 h-5 text-chart-2" />;
      case 'system':
        return <TrendingUp className="w-5 h-5 text-chart-4" />;
      default:
        return <Bell className="w-5 h-5 text-muted-foreground" />;
    }
  };

  const unreadCount = notifications.filter(n => !n.read).length;

  if (!isAuthenticated) {
    return <AuthScreen onAuth={handleAuth} />;
  }

  // Render different screens with animations
  if (currentScreen === 'notifications') {
    return (
      <motion.div
        key="notifications"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <NotificationsScreen 
          user={user} 
          onBack={handleBackToMain}
          notifications={notifications}
          onNotificationUpdate={setNotifications}
          currentCurrency={currentCurrency}
          currencies={currencies}
        />
      </motion.div>
    );
  }

  if (currentScreen === 'settings') {
    return (
      <motion.div
        key="settings"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <SettingsScreen user={user} onBack={handleBackToMain} onLogout={handleLogout} />
      </motion.div>
    );
  }

  if (currentScreen === 'wallet') {
    return (
      <motion.div
        key="wallet"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <WalletScreen user={user} onBack={handleBackToMain} />
      </motion.div>
    );
  }

  if (currentScreen === 'cart') {
    return (
      <motion.div
        key="cart"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <CartScreen user={user} onBack={handleBackToMain} cartItems={cartItems} onUpdateCart={setCartItems} />
      </motion.div>
    );
  }

  if (currentScreen === 'qr-scanner') {
    return (
      <motion.div
        key="qr-scanner"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <QRScannerScreen user={user} onBack={handleBackToMain} />
      </motion.div>
    );
  }

  if (currentScreen === 'search') {
    return (
      <motion.div
        key="search"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <SearchScreen user={user} onBack={handleBackToMain} />
      </motion.div>
    );
  }

  if (currentScreen === 'new-chat') {
    return (
      <motion.div
        key="new-chat"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <NewChatScreen user={user} onBack={handleBackToMain} onCreateChat={handleCreateChat} />
      </motion.div>
    );
  }

  if (currentScreen === 'shop') {
    return (
      <motion.div
        key="shop"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <ShopPage 
          user={user} 
          shopId={selectedShopId} 
          onBack={handleBackToMain} 
          onProductClick={handleProductClick}
          onAddToCart={(productId: string, quantity: number = 1) => {
            const newItem = { id: productId, quantity, timestamp: Date.now() };
            setCartItems(prev => [...prev, newItem]);
            toast.success(`Added ${quantity}x item to cart!`);
          }}
        />
      </motion.div>
    );
  }

  if (currentScreen === 'product') {
    return (
      <motion.div
        key="product"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <ProductPage 
          user={user} 
          productId={selectedProductId} 
          onBack={handleBackToMain} 
          onShopClick={handleShopClick}
          onAddToCart={(productId: string, quantity: number = 1) => {
            const newItem = { id: productId, quantity, timestamp: Date.now() };
            setCartItems(prev => [...prev, newItem]);
            toast.success(`Added ${quantity}x item to cart!`);
          }}
          onBuyNow={(productId: string, quantity: number = 1) => {
            toast.success(`Proceeding to checkout for ${quantity}x item(s)`);
          }}
        />
      </motion.div>
    );
  }

  if (currentScreen === 'live-stream') {
    return (
      <motion.div
        key="live-stream"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <LiveStreamScreen 
          user={user} 
          streamId={selectedStreamId} 
          onBack={handleBackToMain} 
          onProductClick={handleProductClick}
          onAddToCart={(productId: string, quantity: number = 1) => {
            const newItem = { id: productId, quantity, timestamp: Date.now() };
            setCartItems(prev => [...prev, newItem]);
            toast.success(`Added ${quantity}x item to cart!`);
          }}
        />
      </motion.div>
    );
  }

  if (currentScreen === 'shop-chat') {
    return (
      <motion.div
        key="shop-chat"
        initial="initial"
        animate="animate"
        exit="exit"
        variants={screenVariants}
      >
        <ShopChatScreen 
          user={user} 
          shopId={selectedShopId} 
          onBack={handleBackToMain}
          onCall={(shopId) => {
            toast.success(`Calling ${shopId}...`);
            // Could implement actual call functionality here
          }}
          onVideoCall={(shopId) => {
            toast.success(`Starting video call with ${shopId}...`);
            // Could implement actual video call functionality here
          }}
          onOrderNow={(productId) => {
            handleProductClick(productId);
          }}
        />
      </motion.div>
    );
  }

  if (currentScreen === 'video-call' && callData) {
    return (
      <motion.div
        key="video-call"
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.9 }}
        transition={{ duration: 0.3, ease: [0.22, 1, 0.36, 1] }}
      >
        <VideoCallScreen
          user={user}
          callee={callData.callee}
          isIncoming={callData.isIncoming}
          onEndCall={endCall}
          onBack={handleBackToMain}
        />
      </motion.div>
    );
  }

  if (currentScreen === 'voice-call' && callData) {
    return (
      <motion.div
        key="voice-call"
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.9 }}
        transition={{ duration: 0.3, ease: [0.22, 1, 0.36, 1] }}
      >
        <VoiceCallScreen
          user={user}
          callee={callData.callee}
          isIncoming={callData.isIncoming}
          onEndCall={endCall}
          onBack={handleBackToMain}
        />
      </motion.div>
    );
  }

  return (
    <div 
      className="h-screen mobile-full-height bg-background flex flex-col mobile-text overflow-hidden min-w-[320px] max-w-full"
      onTouchStart={handleTouchStart}
      onTouchEnd={handleTouchEnd}
      style={{
        height: isKeyboardOpen && viewportHeight ? `${viewportHeight}px` : '100vh',
        minHeight: isKeyboardOpen && viewportHeight ? `${viewportHeight}px` : '100vh'
      }}
    >
      <Toaster />
      
      {/* Fullscreen Video Player */}
      {fullscreenVideo && (
        <FullscreenVideoPlayer
          video={fullscreenVideo.video}
          isPlaying={fullscreenVideo.isPlaying}
          onClose={handleFullscreenVideoClose}
          onPlay={handleFullscreenVideoPlay}
          onPause={handleFullscreenVideoPause}
          onLike={handleVideoLike}
          onShare={handleVideoShare}
          onSubscribe={handleVideoSubscribe}
          isLiked={likedVideos.includes(fullscreenVideo.video.id)}
          isSubscribed={subscribedChannels.includes(fullscreenVideo.video.channel.id)}
          user={user}
        />
      )}
      
      {/* Share Dialog */}
      {shareDialog.content && (
        <ShareDialog
          open={shareDialog.open}
          onOpenChange={(open) => setShareDialog(prev => ({ ...prev, open }))}
          content={shareDialog.content}
          user={user}
        />
      )}
      
      {/* Enhanced Responsive Header */}
      <motion.header 
        className="sticky top-0 z-50 border-b border-border bg-card/95 backdrop-blur-sm px-2 sm:px-3 lg:px-4 py-2 sm:py-2.5 flex items-center justify-between relative min-h-[52px] sm:min-h-[56px] lg:min-h-[60px] dropdown-container" 
        style={{ overflow: 'visible' }}
        initial={{ y: -60, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.3, ease: [0.22, 1, 0.36, 1] }}
      >
        {/* Pull to refresh indicator */}
        <AnimatePresence>
          {isRefreshing && (
            <motion.div 
              className="absolute top-0 left-0 right-0 h-1 bg-primary/20 z-50"
              initial={{ scaleX: 0 }}
              animate={{ scaleX: 1 }}
              exit={{ scaleX: 0, transition: { duration: 0.2 } }}
              style={{ originX: 0 }}
            >
              <motion.div 
                className="h-full bg-primary"
                animate={{ 
                  x: ['-100%', '100%'],
                  transition: { 
                    repeat: Infinity, 
                    duration: 1, 
                    ease: 'linear' 
                  }
                }}
              />
            </motion.div>
          )}
        </AnimatePresence>
        
        {/* Left Section - App Brand & Context */}
        <div className="flex items-center gap-1 sm:gap-2 lg:gap-3 min-w-0 flex-1">
          {/* App Logo & Title */}
          <div className="flex items-center gap-1.5 sm:gap-2 min-w-0">
            <div className="w-7 h-7 sm:w-8 sm:h-8 lg:w-9 lg:h-9 bg-primary rounded-full flex items-center justify-center flex-shrink-0">
              <MessageCircle className="w-4 h-4 sm:w-5 sm:h-5 lg:w-5 lg:h-5 text-primary-foreground" />
            </div>
            <span className="font-medium text-sm sm:text-base lg:text-lg truncate">Telegram SEA</span>
          </div>
          
          {/* Workspace Switcher - Show on Chat and Work screens */}
          {(currentScreen === 'chat' || currentScreen === 'work') && userWorkspaces.length > 0 && (
            <div className="hidden sm:block">
              <WorkspaceSwitcher
                selectedWorkspace={selectedWorkspace}
                userWorkspaces={userWorkspaces}
                onWorkspaceChange={setSelectedWorkspace}
                variant="compact"
                showMetrics={false}
              />
            </div>
          )}

          {/* Language & Currency Switcher - Desktop Only */}
          <div className="hidden lg:flex items-center gap-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm" className="h-7 px-2 text-xs">
                  <span className="flex items-center gap-1">
                    {languages.find(l => l.code === currentLanguage)?.flag} 
                    <span className="hidden xl:inline">
                      {languages.find(l => l.code === currentLanguage)?.name}
                    </span>
                  </span>
                  <ArrowDown className="w-3 h-3 ml-1" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start" className="w-48 z-[1000]" sideOffset={8}>
                <DropdownMenuLabel className="flex items-center gap-2">
                  <Languages className="w-4 h-4" />
                  Language
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                {languages.map((lang) => (
                  <DropdownMenuItem
                    key={lang.code}
                    onClick={() => handleLanguageSwitch(lang.code)}
                    className={currentLanguage === lang.code ? 'bg-accent' : ''}
                  >
                    <span className="mr-2">{lang.flag}</span>
                    {lang.name}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm" className="h-7 px-2 text-xs">
                  <span className="flex items-center gap-1">
                    {currencies.find(c => c.code === currentCurrency)?.symbol}
                    <span className="hidden xl:inline">{currentCurrency}</span>
                  </span>
                  <ArrowDown className="w-3 h-3 ml-1" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start" className="w-48 z-[1000]" sideOffset={8}>
                <DropdownMenuLabel className="flex items-center gap-2">
                  <Banknote className="w-4 h-4" />
                  Currency
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                {currencies.map((currency) => (
                  <DropdownMenuItem
                    key={currency.code}
                    onClick={() => handleCurrencySwitch(currency.code)}
                    className={currentCurrency === currency.code ? 'bg-accent' : ''}
                  >
                    <span className="mr-2">{currency.symbol}</span>
                    {currency.name}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>

        {/* Right Section - Actions */}
        <div className="flex items-center gap-1 sm:gap-2 flex-shrink-0">
          {/* Search Button */}
          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button 
              variant="ghost" 
              size="icon" 
              className="w-8 h-8 sm:w-9 sm:h-9 touch-manipulation transition-all duration-200"
              onClick={() => setCurrentScreen('search')}
            >
              <Search className="w-4 h-4 sm:w-5 sm:h-5" />
            </Button>
          </motion.div>
          

          
          {/* QR Scanner Button */}
          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button 
              variant="ghost" 
              size="icon" 
              className="relative w-8 h-8 sm:w-9 sm:h-9 touch-manipulation transition-all duration-200"
              onClick={() => setCurrentScreen('qr-scanner')}
              title="QR Scanner"
            >
              <QrCode className="w-4 h-4 sm:w-5 sm:h-5" />
              <motion.div
                animate={{ rotate: [0, 10, -10, 0] }}
                transition={{ repeat: Infinity, duration: 2, ease: "easeInOut" }}
              >
                <Zap className="absolute -bottom-0.5 -right-0.5 w-2.5 h-2.5 sm:w-3 sm:h-3 text-green-500" />
              </motion.div>
            </Button>
          </motion.div>

          {/* Cart Button */}
          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button 
              variant="ghost" 
              size="icon" 
              className="relative w-8 h-8 sm:w-9 sm:h-9 touch-manipulation transition-all duration-200"
              onClick={() => setCurrentScreen('cart')}
            >
              <ShoppingCart className="w-4 h-4 sm:w-5 sm:h-5" />
              <AnimatePresence>
                {cartItems.length > 0 && (
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    exit={{ scale: 0 }}
                    transition={{ type: "spring", stiffness: 500, damping: 30 }}
                  >
                    <Badge className="absolute -top-1 -right-1 w-4 h-4 sm:w-5 sm:h-5 text-[10px] sm:text-xs p-0 flex items-center justify-center bg-chart-2">
                      {cartItems.length}
                    </Badge>
                  </motion.div>
                )}
              </AnimatePresence>
            </Button>
          </motion.div>
          
          {/* Notifications */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="relative w-8 h-8 sm:w-9 sm:h-9 touch-manipulation">
                <Bell className="w-4 h-4 sm:w-5 sm:h-5" />
                {unreadCount > 0 && (
                  <Badge className="absolute -top-1 -right-1 w-4 h-4 sm:w-5 sm:h-5 text-[10px] sm:text-xs p-0 flex items-center justify-center bg-destructive">
                    {unreadCount}
                  </Badge>
                )}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-80 sm:w-96 z-[1000] max-h-[80vh]" sideOffset={8}>
              <div className="flex items-center justify-between p-3">
                <DropdownMenuLabel className="p-0">Notifications</DropdownMenuLabel>
                {unreadCount > 0 && (
                  <Button variant="ghost" size="sm" onClick={markAllAsRead} className="text-xs h-auto p-1">
                    Mark all read
                  </Button>
                )}
              </div>
              <DropdownMenuSeparator />
              
              <ScrollArea className="h-80 sm:h-96">
                {notifications.length > 0 ? (
                  <div className="space-y-1">
                    {notifications.map((notification) => (
                      <DropdownMenuItem
                        key={notification.id}
                        className="p-3 cursor-pointer focus:bg-accent touch-manipulation"
                        onClick={() => handleNotificationClick(notification.id)}
                      >
                        <div className="flex items-start gap-3 w-full">
                          <div className="flex-shrink-0 mt-1">
                            {notification.avatar ? (
                              <Avatar className="w-8 h-8 sm:w-10 sm:h-10">
                                <AvatarImage src={notification.avatar} />
                                <AvatarFallback>
                                  {notification.title.charAt(0)}
                                </AvatarFallback>
                              </Avatar>
                            ) : (
                              <div className="w-8 h-8 sm:w-10 sm:h-10 bg-muted rounded-full flex items-center justify-center">
                                {getNotificationIcon(notification.type)}
                              </div>
                            )}
                          </div>
                          
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 mb-1">
                              <p className={`text-sm ${!notification.read ? 'font-medium' : ''} truncate`}>
                                {notification.title}
                              </p>
                              {!notification.read && (
                                <div className="w-2 h-2 bg-primary rounded-full flex-shrink-0"></div>
                              )}
                            </div>
                            <p className="text-xs text-muted-foreground line-clamp-2 mb-1">
                              {notification.description}
                            </p>
                            <div className="flex items-center justify-between">
                              <p className="text-xs text-muted-foreground">
                                {notification.timestamp}
                              </p>
                              {notification.amount && (
                                <Badge variant="secondary" className="text-xs">
                                  {currencies.find(c => c.code === currentCurrency)?.symbol}{notification.amount}
                                </Badge>
                              )}
                            </div>
                          </div>
                        </div>
                      </DropdownMenuItem>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <Bell className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                    <p className="text-sm text-muted-foreground">No notifications</p>
                  </div>
                )}
              </ScrollArea>
              
              {notifications.length > 0 && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem 
                    className="justify-center p-3 touch-manipulation"
                    onClick={() => setCurrentScreen('notifications')}
                  >
                    <Button variant="ghost" size="sm" className="w-full">
                      View All Notifications
                    </Button>
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </motion.header>

      {/* Main Content Area */}
      <div className="flex-1 flex min-h-0">
        {/* Desktop Sidebar Navigation */}
        <motion.nav 
          className="hidden lg:flex w-16 xl:w-20 bg-sidebar border-r border-sidebar-border flex-col items-center py-3 gap-2 xl:gap-3"
          initial={{ x: -80, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          transition={{ duration: 0.3, ease: [0.22, 1, 0.36, 1] }}
        >
          <motion.div
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <Button
              variant={currentScreen === 'chat' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 relative transition-all duration-200"
              onClick={() => setCurrentScreen('chat')}
            >
              <MessageCircle className="w-5 h-5 xl:w-6 xl:h-6" />
              <motion.div
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                transition={{ delay: 0.2, type: "spring", stiffness: 500, damping: 30 }}
              >
                <Badge className="absolute -top-1 -right-1 w-4 h-4 xl:w-5 xl:h-5 text-[10px] xl:text-xs p-0 flex items-center justify-center">
                  3
                </Badge>
              </motion.div>
            </Button>
          </motion.div>
          
          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'store' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 transition-all duration-200"
              onClick={() => setCurrentScreen('store')}
            >
              <Store className="w-5 h-5 xl:w-6 xl:h-6" />
            </Button>
          </motion.div>
          
          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'social' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 relative transition-all duration-200"
              onClick={() => setCurrentScreen('social')}
            >
              <Users className="w-5 h-5 xl:w-6 xl:h-6" />
              <motion.div
                className="absolute -top-1 -right-1 w-2 h-2 bg-chart-1 rounded-full"
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ repeat: Infinity, duration: 2 }}
              />
            </Button>
          </motion.div>

          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'video' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 relative transition-all duration-200"
              onClick={() => setCurrentScreen('video')}
            >
              <Play className="w-5 h-5 xl:w-6 xl:h-6" />
              <AnimatePresence>
                {watchedVideos.length > 0 && (
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    exit={{ scale: 0 }}
                    transition={{ type: "spring", stiffness: 500, damping: 30 }}
                  >
                    <Badge className="absolute -top-1 -right-1 w-4 h-4 xl:w-5 xl:h-5 text-[10px] xl:text-xs p-0 flex items-center justify-center bg-chart-1">
                      {watchedVideos.length}
                    </Badge>
                  </motion.div>
                )}
              </AnimatePresence>
            </Button>
          </motion.div>

          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'work' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 relative transition-all duration-200"
              onClick={() => setCurrentScreen('work')}
            >
              <Briefcase className="w-5 h-5 xl:w-6 xl:h-6" />
              <motion.div
                className="absolute -top-1 -right-1 w-2 h-2 bg-chart-2 rounded-full"
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ repeat: Infinity, duration: 2, delay: 0.5 }}
              />
            </Button>
          </motion.div>

          <motion.div 
            className="mt-auto space-y-2 xl:space-y-3"
            variants={staggerChildrenVariants}
            initial="initial"
            animate="animate"
          >
            <motion.div variants={childVariants} whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button 
                variant={currentScreen === 'settings' ? 'default' : 'ghost'}
                size="icon" 
                className="w-12 h-12 xl:w-14 xl:h-14 transition-all duration-200"
                onClick={() => setCurrentScreen('settings')}
              >
                <Settings className="w-5 h-5 xl:w-6 xl:h-6" />
              </Button>
            </motion.div>
            <motion.div variants={childVariants} whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button 
                variant={currentScreen === 'wallet' ? 'default' : 'ghost'}
                size="icon" 
                className="w-12 h-12 xl:w-14 xl:h-14 transition-all duration-200"
                onClick={() => setCurrentScreen('wallet')}
              >
                <Wallet className="w-5 h-5 xl:w-6 xl:h-6" />
              </Button>
            </motion.div>
          </motion.div>
        </motion.nav>

        {/* Tab Content */}
        <main 
          className="flex-1 relative overflow-hidden min-h-0 bg-background w-full"
          ref={mainContentRef}
        >
          {/* Scrollable Content Container */}
          <div 
            className="h-full w-full overflow-y-auto overflow-x-hidden mobile-scroll scrollbar-hide"
            style={{
              paddingBottom: isKeyboardOpen 
                ? '0.5rem' 
                : 'max(5rem, calc(5rem + env(safe-area-inset-bottom, 0)))',
              height: '100%',
              position: 'relative'
            }}
          >
            <div className="min-h-full w-full">
              <AnimatePresence mode="wait">
                {currentScreen === 'chat' && (
                  <motion.div
                    key="chat-tab"
                    initial="initial"
                    animate="animate"
                    exit="exit"
                    variants={slideVariants}
                    custom={currentTabIndex}
                    className="w-full h-full min-h-0"
                  >
                    <RichChatTab 
                      user={user} 
                      onVideoCall={startVideoCall} 
                      onVoiceCall={startVoiceCall} 
                      onNewChat={() => setCurrentScreen('new-chat')}
                      selectedWorkspace={selectedWorkspace}
                      onWorkspaceChange={setSelectedWorkspace}
                      userWorkspaces={userWorkspaces}
                      onShopChatClick={handleShopChatClick}
                    />
                  </motion.div>
                )}
                {currentScreen === 'store' && (
                  <motion.div
                    key="store-tab"
                    initial="initial"
                    animate="animate"
                    exit="exit"
                    variants={slideVariants}
                    custom={currentTabIndex}
                    className="w-full h-full min-h-0"
                  >
                    <StoreTab 
                      user={user} 
                      onShopClick={handleShopClick}
                      onProductClick={handleProductClick}
                      onAddToCart={(productId: string, quantity: number = 1) => {
                        const newItem = { id: productId, quantity, timestamp: Date.now() };
                        setCartItems(prev => [...prev, newItem]);
                        toast.success(`Added ${quantity}x item to cart!`);
                      }}
                      onLiveStreamClick={handleLiveStreamClick}
                      onProductShare={handleProductShare}
                      onShopShare={handleShopShare}
                      cartItems={cartItems}
                    />
                  </motion.div>
                )}
                {currentScreen === 'social' && (
                  <motion.div
                    key="social-tab"
                    initial="initial"
                    animate="animate"
                    exit="exit"
                    variants={slideVariants}
                    custom={currentTabIndex}
                    className="w-full h-full min-h-0"
                  >
                    <SocialTab 
                      user={user} 
                      onLiveStreamClick={handleLiveStreamClick}
                      onPostShare={handlePostShare}
                    />
                  </motion.div>
                )}
                {currentScreen === 'video' && (
                  <motion.div
                    key="video-tab"
                    initial="initial"
                    animate="animate"
                    exit="exit"
                    variants={slideVariants}
                    custom={currentTabIndex}
                    className="w-full h-full min-h-0"
                  >
                    <VideoTab 
                      user={user} 
                      onBack={handleBackToMain}
                      onVideoPlay={handleVideoPlay}
                      onVideoLike={handleVideoLike}
                      onVideoShare={handleVideoShare}
                      onSubscribe={handleVideoSubscribe}
                      currentVideoId={currentVideoId}
                      watchedVideos={watchedVideos}
                      likedVideos={likedVideos}
                      subscribedChannels={subscribedChannels}
                    />
                  </motion.div>
                )}
                {currentScreen === 'work' && (
                  <motion.div
                    key="work-tab"
                    initial="initial"
                    animate="animate"
                    exit="exit"
                    variants={slideVariants}
                    custom={currentTabIndex}
                    className="w-full h-full min-h-0"
                  >
                    <WorkspaceTab 
                      user={user} 
                      onBack={handleBackToMain}
                    />
                  </motion.div>
                )}
                {currentScreen === 'more' && (
                  <motion.div
                    key="more-tab"
                    initial="initial"
                    animate="animate"
                    exit="exit"
                    variants={slideVariants}
                    custom={currentTabIndex}
                    className="w-full h-full min-h-0"
                  >
                    {/* Enhanced More Tab Content */}
                    <div className="w-full h-full overflow-y-auto mobile-scroll scrollbar-hide">
                      <div className="w-full max-w-2xl mx-auto px-4 py-6 space-y-8">
                      {/* User Profile Section */}
                      <motion.section 
                        className="flex items-center gap-4 p-6 bg-gradient-to-r from-card to-card/50 rounded-2xl border shadow-sm"
                        variants={childVariants}
                        whileHover={{ scale: 1.02, transition: { duration: 0.2 } }}
                        whileTap={{ scale: 0.98 }}
                      >
                        <Avatar className="w-20 h-20 ring-2 ring-primary/10 flex-shrink-0">
                          <AvatarImage src={user?.avatar} />
                          <AvatarFallback className="text-lg font-medium">
                            {user?.name?.charAt(0) || 'U'}
                          </AvatarFallback>
                        </Avatar>
                        <div className="flex-1 space-y-1 min-w-0">
                          <h2 className="text-xl font-semibold truncate">{user?.name || 'User'}</h2>
                          <p className="text-sm text-muted-foreground truncate">
                            {user?.email || user?.phone || 'Telegram SEA User'}
                          </p>
                          <div className="flex items-center gap-2 pt-1 flex-wrap">
                            <Badge variant="secondary" className="text-xs">
                              SEA Edition
                            </Badge>
                            <Badge variant="outline" className="text-xs">
                              {currentLanguage.toUpperCase()}
                            </Badge>
                            <Badge variant="outline" className="text-xs bg-green-50 text-green-700 border-green-200">
                              ✓ Verified
                            </Badge>
                          </div>
                        </div>
                      </motion.section>

                      {/* Quick Actions Grid */}
                      <motion.section 
                        className="space-y-4"
                        variants={staggerChildrenVariants}
                        initial="initial"
                        animate="animate"
                      >
                        <h3 className="text-lg font-semibold">Quick Actions</h3>
                        <div className="grid grid-cols-2 gap-4">
                          <motion.div variants={childVariants} whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                            <Button
                              variant="outline"
                              className="w-full aspect-square min-h-[120px] p-4 flex flex-col gap-3 justify-center items-center touch-manipulation transition-all duration-200 hover:shadow-md rounded-xl group"
                              onClick={() => setCurrentScreen('wallet')}
                            >
                              <div className="w-12 h-12 rounded-full bg-chart-1/10 flex items-center justify-center group-hover:bg-chart-1/20 transition-colors flex-shrink-0">
                                <Wallet className="w-6 h-6 text-chart-1" />
                              </div>
                              <span className="text-sm font-medium text-center">Wallet</span>
                            </Button>
                          </motion.div>
                          
                          <motion.div variants={childVariants} whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                            <Button
                              variant="outline"
                              className="w-full aspect-square min-h-[120px] p-4 flex flex-col gap-3 justify-center items-center touch-manipulation transition-all duration-200 hover:shadow-md rounded-xl group"
                              onClick={() => setCurrentScreen('settings')}
                            >
                              <div className="w-12 h-12 rounded-full bg-chart-2/10 flex items-center justify-center group-hover:bg-chart-2/20 transition-colors flex-shrink-0">
                                <Settings className="w-6 h-6 text-chart-2" />
                              </div>
                              <span className="text-sm font-medium text-center">Settings</span>
                            </Button>
                          </motion.div>
                          
                          <motion.div variants={childVariants} whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                            <Button
                              variant="outline"
                              className="w-full aspect-square min-h-[120px] p-4 flex flex-col gap-3 justify-center items-center touch-manipulation transition-all duration-200 hover:shadow-md rounded-xl group"
                              onClick={() => setCurrentScreen('work')}
                            >
                              <div className="w-12 h-12 rounded-full bg-chart-3/10 flex items-center justify-center group-hover:bg-chart-3/20 transition-colors flex-shrink-0">
                                <Briefcase className="w-6 h-6 text-chart-3" />
                              </div>
                              <span className="text-sm font-medium text-center">Work</span>
                            </Button>
                          </motion.div>
                          
                          <motion.div variants={childVariants} whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                            <Button
                              variant="outline"
                              className="w-full aspect-square min-h-[120px] p-4 flex flex-col gap-3 justify-center items-center touch-manipulation transition-all duration-200 hover:shadow-md rounded-xl group"
                              onClick={() => setCurrentScreen('qr-scanner')}
                            >
                              <div className="w-12 h-12 rounded-full bg-chart-4/10 flex items-center justify-center group-hover:bg-chart-4/20 transition-colors flex-shrink-0">
                                <QrCode className="w-6 h-6 text-chart-4" />
                              </div>
                              <span className="text-sm font-medium text-center">QR Scanner</span>
                            </Button>
                          </motion.div>
                        </div>
                      </motion.section>

                      {/* SEA Super-App Features */}
                      <motion.section 
                        className="space-y-4"
                        variants={staggerChildrenVariants}
                        initial="initial"
                        animate="animate"
                      >
                        <h3 className="text-lg font-semibold">SEA Features</h3>
                        <div className="grid grid-cols-1 gap-3">
                          <motion.div 
                            variants={childVariants}
                            className="flex items-center justify-between p-4 bg-card rounded-xl border hover:shadow-sm transition-shadow cursor-pointer"
                            onClick={() => toast.success('Mini-apps coming soon!')}
                          >
                            <div className="flex items-center gap-3 flex-1 min-w-0">
                              <div className="w-10 h-10 rounded-full bg-chart-1/10 flex items-center justify-center flex-shrink-0">
                                <Globe className="w-5 h-5 text-chart-1" />
                              </div>
                              <div className="min-w-0 flex-1">
                                <span className="text-sm font-medium block">Mini-Apps</span>
                                <p className="text-xs text-muted-foreground">3rd party apps & services</p>
                              </div>
                            </div>
                            <Badge variant="secondary" className="text-xs">
                              Coming Soon
                            </Badge>
                          </motion.div>

                          <motion.div 
                            variants={childVariants}
                            className="flex items-center justify-between p-4 bg-card rounded-xl border hover:shadow-sm transition-shadow cursor-pointer"
                            onClick={() => toast.success('Data saver mode enabled')}
                          >
                            <div className="flex items-center gap-3 flex-1 min-w-0">
                              <div className="w-10 h-10 rounded-full bg-chart-2/10 flex items-center justify-center flex-shrink-0">
                                <Zap className="w-5 h-5 text-chart-2" />
                              </div>
                              <div className="min-w-0 flex-1">
                                <span className="text-sm font-medium block">Data Saver</span>
                                <p className="text-xs text-muted-foreground">Ultra-low data usage mode</p>
                              </div>
                            </div>
                            <div className="w-8 h-5 bg-chart-2/20 rounded-full relative">
                              <div className="absolute left-1 top-1 w-3 h-3 bg-chart-2 rounded-full"></div>
                            </div>
                          </motion.div>

                          <motion.div 
                            variants={childVariants}
                            className="flex items-center justify-between p-4 bg-card rounded-xl border hover:shadow-sm transition-shadow cursor-pointer"
                            onClick={() => {
                              const methods = qrPaymentMethods.map(m => m.name).join(', ');
                              toast.success(`QR Payments: ${methods}`);
                            }}
                          >
                            <div className="flex items-center gap-3 flex-1 min-w-0">
                              <div className="w-10 h-10 rounded-full bg-chart-3/10 flex items-center justify-center flex-shrink-0">
                                <CreditCard className="w-5 h-5 text-chart-3" />
                              </div>
                              <div className="min-w-0 flex-1">
                                <span className="text-sm font-medium block">QR Payments</span>
                                <p className="text-xs text-muted-foreground">PromptPay, QRIS, VietQR & more</p>
                              </div>
                            </div>
                            <div className="flex items-center gap-1">
                              {qrPaymentMethods.slice(0, 3).map((method) => (
                                <span key={method.country} className="text-xs">
                                  {method.icon}
                                </span>
                              ))}
                            </div>
                          </motion.div>
                        </div>
                      </motion.section>

                      {/* Preferences Section */}
                      <motion.section 
                        className="space-y-4"
                        variants={staggerChildrenVariants}
                        initial="initial"
                        animate="animate"
                      >
                        <h3 className="text-lg font-semibold">Preferences</h3>
                        
                        <div className="space-y-3">
                          <motion.div 
                            variants={childVariants}
                            className="flex items-center justify-between p-4 bg-card rounded-xl border hover:shadow-sm transition-shadow"
                          >
                            <div className="flex items-center gap-3 flex-1 min-w-0">
                              <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center flex-shrink-0">
                                <Languages className="w-5 h-5 text-muted-foreground" />
                              </div>
                              <div className="min-w-0 flex-1">
                                <span className="text-sm font-medium block">Language</span>
                                <p className="text-xs text-muted-foreground">App display language</p>
                              </div>
                            </div>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="sm" className="h-auto p-2 hover:bg-accent rounded-lg flex-shrink-0">
                                  <span className="flex items-center gap-2">
                                    {languages.find(l => l.code === currentLanguage)?.flag}
                                    <span className="text-sm font-medium">
                                      {languages.find(l => l.code === currentLanguage)?.name}
                                    </span>
                                    <ArrowDown className="w-3 h-3 ml-1" />
                                  </span>
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end" className="w-52 z-[1000]">
                                {languages.map((lang) => (
                                  <DropdownMenuItem
                                    key={lang.code}
                                    onClick={() => handleLanguageSwitch(lang.code)}
                                    className={`${currentLanguage === lang.code ? 'bg-accent' : ''} cursor-pointer`}
                                  >
                                    <span className="mr-3">{lang.flag}</span>
                                    <span>{lang.name}</span>
                                  </DropdownMenuItem>
                                ))}
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </motion.div>

                          <motion.div 
                            variants={childVariants}
                            className="flex items-center justify-between p-4 bg-card rounded-xl border hover:shadow-sm transition-shadow"
                          >
                            <div className="flex items-center gap-3 flex-1 min-w-0">
                              <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center flex-shrink-0">
                                <Banknote className="w-5 h-5 text-muted-foreground" />
                              </div>
                              <div className="min-w-0 flex-1">
                                <span className="text-sm font-medium block">Currency</span>
                                <p className="text-xs text-muted-foreground">Payment currency</p>
                              </div>
                            </div>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="sm" className="h-auto p-2 hover:bg-accent rounded-lg flex-shrink-0">
                                  <span className="flex items-center gap-2">
                                    <span className="text-sm font-medium">
                                      {currencies.find(c => c.code === currentCurrency)?.symbol} {currentCurrency}
                                    </span>
                                    <ArrowDown className="w-3 h-3 ml-1" />
                                  </span>
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end" className="w-52 z-[1000]">
                                {currencies.map((currency) => (
                                  <DropdownMenuItem
                                    key={currency.code}
                                    onClick={() => handleCurrencySwitch(currency.code)}
                                    className={`${currentCurrency === currency.code ? 'bg-accent' : ''} cursor-pointer`}
                                  >
                                    <span className="mr-3">{currency.symbol}</span>
                                    <span>{currency.name}</span>
                                  </DropdownMenuItem>
                                ))}
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </motion.div>
                        </div>
                      </motion.section>

                      {/* App Info Section */}
                      <motion.section 
                        className="space-y-4"
                        variants={staggerChildrenVariants}
                        initial="initial"
                        animate="animate"
                      >
                        <h3 className="text-lg font-semibold">App Information</h3>
                        <div className="space-y-3">
                          <motion.div 
                            variants={childVariants}
                            className="flex items-center justify-between p-4 bg-card rounded-xl border hover:shadow-sm transition-shadow cursor-pointer"
                            onClick={() => toast.success('Telegram SEA Edition v2.0.1')}
                          >
                            <div className="flex items-center gap-3 flex-1 min-w-0">
                              <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
                                <MessageCircle className="w-5 h-5 text-primary" />
                              </div>
                              <div className="min-w-0 flex-1">
                                <span className="text-sm font-medium block">Version</span>
                                <p className="text-xs text-muted-foreground">Telegram SEA Edition</p>
                              </div>
                            </div>
                            <span className="text-sm text-muted-foreground">v2.0.1</span>
                          </motion.div>

                          <motion.div 
                            variants={childVariants}
                            className="flex items-center justify-between p-4 bg-card rounded-xl border hover:shadow-sm transition-shadow cursor-pointer"
                            onClick={() => toast.success('Data usage: 2.3 MB today')}
                          >
                            <div className="flex items-center gap-3 flex-1 min-w-0">
                              <div className="w-10 h-10 rounded-full bg-chart-2/10 flex items-center justify-center flex-shrink-0">
                                <TrendingUp className="w-5 h-5 text-chart-2" />
                              </div>
                              <div className="min-w-0 flex-1">
                                <span className="text-sm font-medium block">Data Usage</span>
                                <p className="text-xs text-muted-foreground">Ultra-low consumption</p>
                              </div>
                            </div>
                            <span className="text-sm text-green-600">2.3 MB</span>
                          </motion.div>

                          <motion.div 
                            variants={childVariants}
                            className="flex items-center justify-between p-4 bg-card rounded-xl border hover:shadow-sm transition-shadow cursor-pointer"
                            onClick={() => toast.success('Help & Support coming soon')}
                          >
                            <div className="flex items-center gap-3 flex-1 min-w-0">
                              <div className="w-10 h-10 rounded-full bg-chart-4/10 flex items-center justify-center flex-shrink-0">
                                <Bell className="w-5 h-5 text-chart-4" />
                              </div>
                              <div className="min-w-0 flex-1">
                                <span className="text-sm font-medium block">Help & Support</span>
                                <p className="text-xs text-muted-foreground">Get help and report issues</p>
                              </div>
                            </div>
                            <ArrowDown className="w-4 h-4 text-muted-foreground rotate-[-90deg]" />
                          </motion.div>
                        </div>
                      </motion.section>

                      {/* Account Section */}
                      <motion.section 
                        className="space-y-4 pb-6"
                        variants={staggerChildrenVariants}
                        initial="initial"
                        animate="animate"
                      >
                        <h3 className="text-lg font-semibold">Account</h3>
                        <motion.div variants={childVariants}>
                          <Button
                            variant="outline"
                            className="w-full justify-start p-4 text-destructive hover:text-destructive hover:bg-destructive/5 border-destructive/20 hover:border-destructive/30 touch-manipulation rounded-xl"
                            onClick={handleLogout}
                          >
                            <div className="w-10 h-10 rounded-full bg-destructive/10 flex items-center justify-center mr-3 flex-shrink-0">
                              <User className="w-5 h-5 text-destructive" />
                            </div>
                            <div className="text-left min-w-0 flex-1">
                              <div className="font-medium">Sign Out</div>
                              <div className="text-xs text-muted-foreground">Logout from your account</div>
                            </div>
                          </Button>
                        </motion.div>
                      </motion.section>
                      </div>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          </div>
        </main>
      </div>

      {/* Enhanced Mobile Bottom Navigation */}
      <motion.nav 
        className={`lg:hidden fixed bottom-0 left-0 right-0 z-50 border-t border-border bg-card/95 backdrop-blur-sm px-1 py-1.5 flex items-center justify-around ${
          isKeyboardOpen ? 'translate-y-full opacity-0' : 'translate-y-0 opacity-100'
        }`}
        style={{
          paddingBottom: `max(0.375rem, env(safe-area-inset-bottom, 0.375rem))`
        }}
        initial={{ y: 100, opacity: 0 }}
        animate={{ 
          y: isKeyboardOpen ? 100 : 0, 
          opacity: isKeyboardOpen ? 0 : 1,
          transition: { duration: 0.3, ease: [0.22, 1, 0.36, 1] }
        }}
      >
        {/* Active Tab Indicator */}
        <div className="absolute top-0 left-0 right-0 h-0.5 bg-muted">
          <motion.div 
            className="h-full bg-primary"
            animate={{ 
              x: `${currentTabIndex * 100}%`
            }}
            transition={{ 
              type: "spring", 
              stiffness: 300, 
              damping: 30 
            }}
            style={{ width: '20%' }}
          />
        </div>
        
        <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
          <Button
            variant={currentScreen === 'chat' ? 'default' : 'ghost'}
            size="sm"
            className="flex flex-col gap-0.5 h-auto py-2 relative min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200 mobile-button-large"
            onClick={() => setCurrentScreen('chat')}
          >
            <MessageCircle className="w-5 h-5 mb-0.5" />
            <span className="text-[10px] leading-none font-medium">Chat</span>
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ delay: 0.1, type: "spring", stiffness: 500, damping: 30 }}
            >
              <Badge className="absolute -top-0.5 -right-0.5 w-4 h-4 text-[8px] p-0 flex items-center justify-center">
                3
              </Badge>
            </motion.div>
          </Button>
        </motion.div>
        
        <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
          <Button
            variant={currentScreen === 'store' ? 'default' : 'ghost'}
            size="sm"
            className="flex flex-col gap-0.5 h-auto py-2 min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200 mobile-button-large"
            onClick={() => setCurrentScreen('store')}
          >
            <Store className="w-5 h-5 mb-0.5" />
            <span className="text-[10px] leading-none font-medium">Store</span>
          </Button>
        </motion.div>
        
        <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
          <Button
            variant={currentScreen === 'social' ? 'default' : 'ghost'}
            size="sm"
            className="flex flex-col gap-0.5 h-auto py-2 relative min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200 mobile-button-large"
            onClick={() => setCurrentScreen('social')}
          >
            <Users className="w-5 h-5 mb-0.5" />
            <span className="text-[10px] leading-none font-medium">Social</span>
            <motion.div
              className="absolute -top-0.5 -right-0.5 w-2.5 h-2.5 bg-chart-1 rounded-full border border-card"
              animate={{ scale: [1, 1.2, 1] }}
              transition={{ repeat: Infinity, duration: 2 }}
            />
          </Button>
        </motion.div>
        
        <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
          <Button
            variant={currentScreen === 'video' ? 'default' : 'ghost'}
            size="sm"
            className="flex flex-col gap-0.5 h-auto py-2 relative min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200 mobile-button-large"
            onClick={() => setCurrentScreen('video')}
          >
            <Play className="w-5 h-5 mb-0.5" />
            <span className="text-[10px] leading-none font-medium">Video</span>
            <AnimatePresence>
              {watchedVideos.length > 0 && (
                <motion.div
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  exit={{ scale: 0 }}
                  transition={{ type: "spring", stiffness: 500, damping: 30 }}
                >
                  <Badge className="absolute -top-0.5 -right-0.5 w-4 h-4 text-[8px] p-0 flex items-center justify-center bg-chart-1">
                    {watchedVideos.length}
                  </Badge>
                </motion.div>
              )}
            </AnimatePresence>
          </Button>
        </motion.div>
        
        <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
          <Button
            variant={currentScreen === 'more' ? 'default' : 'ghost'}
            size="sm"
            className="flex flex-col gap-0.5 h-auto py-2 min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200 mobile-button-large"
            onClick={() => setCurrentScreen('more')}
          >
            <User className="w-5 h-5 mb-0.5" />
            <span className="text-[10px] leading-none font-medium">More</span>
          </Button>
        </motion.div>

      </motion.nav>
    </div>
  );
}