import React, { useState, useCallback, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { motion, AnimatePresence } from 'framer-motion';
import { Button } from './ui/button';
import {
  MessageCircle,
  ShoppingBag,
  Users,
  Video,
  MoreHorizontal,
  Settings,
  Wallet,
  Bell
} from 'lucide-react';
import { AuthScreen } from './AuthScreen';
import { selectFallbackMode } from '../features/contentSlice';

// Define the Screen type
type Screen = 'chat' | 'store' | 'social' | 'video' | 'more' | 'notifications' | 'settings' | 'wallet' | 'cart' | 'qr-scanner' | 'search' | 'new-chat' | 'shop' | 'product' | 'live-stream' | 'shop-chat' | 'video-call' | 'voice-call';

export function Layout() {
  const navigate = useNavigate();
  const location = useLocation();
  const fallbackMode = useSelector(selectFallbackMode);

  // Authentication state
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<any>(null);
  const [currentLanguage, setCurrentLanguage] = useState('th');

  // Get current route for navigation highlighting
  const getCurrentRoute = () => {
    const path = location.pathname;
    if (path === '/' || path === '/chat') return 'chat';
    if (path.startsWith('/store')) return 'store';
    if (path.startsWith('/social')) return 'social';
    if (path.startsWith('/video')) return 'video';
    if (path.startsWith('/more')) return 'more';
    return 'chat';
  };

  const currentScreen = getCurrentRoute();

  // Handle authentication
  const handleAuth = (userData: any) => {
    setIsAuthenticated(true);
    setUser(userData);
  };

  const handleSignOut = () => {
    setIsAuthenticated(false);
    setUser(null);
    navigate('/chat');
  };

  // Navigation handlers
  const navigateToScreen = useCallback((screen: string) => {
    navigate(`/${screen}`);
  }, [navigate]);

  // Touch handlers for mobile swipe navigation
  const tabs: Screen[] = ['chat', 'store', 'social', 'video', 'more'];
  const currentTabIndex = tabs.indexOf(currentScreen as Screen);

  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    const touch = e.touches[0];
    window.touchStartX = touch.clientX;
  }, []);

  const handleTouchEnd = useCallback((e: React.TouchEvent) => {
    if (!window.touchStartX) return;

    const touch = e.changedTouches[0];
    const touchEndX = touch.clientX;
    const diff = window.touchStartX - touchEndX;
    const threshold = 100;

    if (Math.abs(diff) > threshold) {
      if (diff > 0 && currentTabIndex < tabs.length - 1) {
        // Swipe left - next tab
        navigateToScreen(tabs[currentTabIndex + 1]);
      } else if (diff < 0 && currentTabIndex > 0) {
        // Swipe right - previous tab
        navigateToScreen(tabs[currentTabIndex - 1]);
      }
    }

    window.touchStartX = null;
  }, [currentTabIndex, navigateToScreen, tabs]);

  // Show auth screen if not authenticated
  if (!isAuthenticated) {
    return (
      <AuthScreen
        onAuth={handleAuth}
        currentLanguage={currentLanguage}
        onLanguageChange={setCurrentLanguage}
        fallbackMode={fallbackMode}
      />
    );
  }

  return (
    <div className="h-screen flex flex-col bg-background overflow-hidden">
      {/* Main Content Area */}
      <div
        className="flex-1 flex overflow-hidden"
        onTouchStart={handleTouchStart}
        onTouchEnd={handleTouchEnd}
      >
        {/* Desktop Sidebar */}
        <motion.div
          className="hidden lg:flex flex-col w-20 xl:w-24 bg-muted/30 border-r border-border relative z-20"
          initial={{ x: -80, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          transition={{ duration: 0.3, ease: [0.22, 1, 0.36, 1] }}
          aria-label="Main navigation"
        >
          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'chat' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 relative transition-all duration-200"
              onClick={() => navigateToScreen('chat')}
              aria-label="Chat"
            >
              <MessageCircle className="w-5 h-5 xl:w-6 xl:h-6" />
            </Button>
          </motion.div>

          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'store' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 transition-all duration-200"
              onClick={() => navigateToScreen('store')}
              aria-label="Store"
            >
              <ShoppingBag className="w-5 h-5 xl:w-6 xl:h-6" />
            </Button>
          </motion.div>

          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'social' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 relative transition-all duration-200"
              onClick={() => navigateToScreen('social')}
              aria-label="Social"
            >
              <Users className="w-5 h-5 xl:w-6 xl:h-6" />
            </Button>
          </motion.div>

          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'video' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 relative transition-all duration-200"
              onClick={() => navigateToScreen('video')}
              aria-label="Video"
            >
              <Video className="w-5 h-5 xl:w-6 xl:h-6" />
            </Button>
          </motion.div>

          <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
            <Button
              variant={currentScreen === 'more' ? 'default' : 'ghost'}
              size="icon"
              className="w-12 h-12 xl:w-14 xl:h-14 relative transition-all duration-200"
              onClick={() => navigateToScreen('more')}
              aria-label="More"
            >
              <MoreHorizontal className="w-5 h-5 xl:w-6 xl:h-6" />
            </Button>
          </motion.div>

          {/* Settings and other actions */}
          <div className="mt-auto mb-4 space-y-2">
            <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button
                variant="ghost"
                size="icon"
                className="w-12 h-12 xl:w-14 xl:h-14 transition-all duration-200"
                onClick={() => navigateToScreen('settings')}
                aria-label="Settings"
              >
                <Settings className="w-5 h-5 xl:w-6 xl:h-6" />
              </Button>
            </motion.div>
            <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button
                variant="ghost"
                size="icon"
                className="w-12 h-12 xl:w-14 xl:h-14 transition-all duration-200"
                onClick={() => navigateToScreen('wallet')}
                aria-label="Wallet"
              >
                <Wallet className="w-5 h-5 xl:w-6 xl:h-6" />
              </Button>
            </motion.div>
          </div>
        </motion.div>

        {/* Main Content */}
        <div className="flex-1 flex flex-col overflow-hidden">
          <div className="min-h-full w-full">
            <AnimatePresence mode="wait">
              <motion.div
                key={location.pathname}
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.2 }}
                className="h-full"
              >
                <Outlet />
              </motion.div>
            </AnimatePresence>
          </div>
        </div>
      </div>

      {/* Mobile Bottom Navigation */}
      <div className="lg:hidden bg-background/95 backdrop-blur-sm border-t border-border">
        <div className="flex items-center justify-around px-2 py-2 max-w-md mx-auto">
          <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
            <Button
              variant={currentScreen === 'chat' ? 'default' : 'ghost'}
              size="sm"
              className="flex flex-col gap-0.5 h-auto py-2 relative min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200"
              onClick={() => navigateToScreen('chat')}
              aria-label="Chat"
            >
              <MessageCircle className="w-5 h-5" />
              <span className="text-xs font-medium">Chat</span>
            </Button>
          </motion.div>

          <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
            <Button
              variant={currentScreen === 'store' ? 'default' : 'ghost'}
              size="sm"
              className="flex flex-col gap-0.5 h-auto py-2 min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200"
              onClick={() => navigateToScreen('store')}
              aria-label="Store"
            >
              <ShoppingBag className="w-5 h-5" />
              <span className="text-xs font-medium">Store</span>
            </Button>
          </motion.div>

          <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
            <Button
              variant={currentScreen === 'social' ? 'default' : 'ghost'}
              size="sm"
              className="flex flex-col gap-0.5 h-auto py-2 relative min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200"
              onClick={() => navigateToScreen('social')}
              aria-label="Social"
            >
              <Users className="w-5 h-5" />
              <span className="text-xs font-medium">Social</span>
            </Button>
          </motion.div>

          <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
            <Button
              variant={currentScreen === 'video' ? 'default' : 'ghost'}
              size="sm"
              className="flex flex-col gap-0.5 h-auto py-2 relative min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200"
              onClick={() => navigateToScreen('video')}
              aria-label="Video"
            >
              <Video className="w-5 h-5" />
              <span className="text-xs font-medium">Video</span>
            </Button>
          </motion.div>

          <motion.div whileTap={{ scale: 0.9 }} transition={{ type: "spring", stiffness: 400, damping: 17 }}>
            <Button
              variant={currentScreen === 'more' ? 'default' : 'ghost'}
              size="sm"
              className="flex flex-col gap-0.5 h-auto py-2 min-w-0 px-1.5 rounded-lg touch-manipulation transition-all duration-200"
              onClick={() => navigateToScreen('more')}
              aria-label="More"
            >
              <MoreHorizontal className="w-5 h-5" />
              <span className="text-xs font-medium">More</span>
            </Button>
          </motion.div>
        </div>
      </div>
    </div>
  );
}

export default Layout;