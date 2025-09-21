import React, { useState } from 'react';
import { ArrowLeft, Bell, MessageCircle, CreditCard, Heart, Store, TrendingUp, Check, Trash2, Archive, Filter, Search, MoreVertical, Settings, BellOff } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { Button } from './ui/button';
import { Card, CardContent } from './ui/card';
import { Badge } from './ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { ScrollArea } from './ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Input } from './ui/input';
import { Separator } from './ui/separator';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from './ui/dropdown-menu';
import { Checkbox } from './ui/checkbox';
import { toast } from "sonner";

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

interface NotificationsScreenProps {
  user: any;
  onBack: () => void;
  notifications: Notification[];
  onNotificationUpdate: (notifications: Notification[]) => void;
  currentCurrency: string;
  currencies: Array<{ code: string; symbol: string; name: string }>;
}

export function NotificationsScreen({ 
  user, 
  onBack, 
  notifications: initialNotifications,
  onNotificationUpdate,
  currentCurrency,
  currencies
}: NotificationsScreenProps) {
  const [notifications, setNotifications] = useState<Notification[]>(initialNotifications);
  const [selectedTab, setSelectedTab] = useState<'all' | 'unread' | 'messages' | 'payments' | 'social'>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedNotifications, setSelectedNotifications] = useState<string[]>([]);
  const [isSelectMode, setIsSelectMode] = useState(false);

  // Animation variants
  const containerVariants = {
    initial: { opacity: 0, y: 20 },
    animate: { 
      opacity: 1, 
      y: 0,
      transition: {
        duration: 0.3,
        ease: [0.22, 1, 0.36, 1],
        staggerChildren: 0.1
      }
    },
    exit: { 
      opacity: 0, 
      y: -20,
      transition: {
        duration: 0.2,
        ease: [0.22, 1, 0.36, 1]
      }
    }
  };

  const itemVariants = {
    initial: { opacity: 0, x: 20 },
    animate: { 
      opacity: 1, 
      x: 0,
      transition: {
        duration: 0.3,
        ease: [0.22, 1, 0.36, 1]
      }
    },
    exit: { 
      opacity: 0, 
      x: -20,
      transition: {
        duration: 0.2,
        ease: [0.22, 1, 0.36, 1]
      }
    }
  };

  const handleNotificationClick = (notificationId: string) => {
    if (isSelectMode) {
      toggleNotificationSelection(notificationId);
      return;
    }

    const updatedNotifications = notifications.map(notification => 
      notification.id === notificationId 
        ? { ...notification, read: true }
        : notification
    );
    setNotifications(updatedNotifications);
    onNotificationUpdate(updatedNotifications);
    
    // Simulate navigation to relevant screen based on notification type
    const notification = notifications.find(n => n.id === notificationId);
    if (notification) {
      switch (notification.type) {
        case 'message':
          toast.success('Opening chat...');
          break;
        case 'payment':
          toast.success('Opening wallet...');
          break;
        case 'social':
          toast.success('Opening social feed...');
          break;
        case 'merchant':
          toast.success('Opening store...');
          break;
        case 'system':
          toast.success('Opening notification details...');
          break;
      }
    }
  };

  const markAsRead = (notificationIds: string[]) => {
    const updatedNotifications = notifications.map(notification => 
      notificationIds.includes(notification.id) 
        ? { ...notification, read: true }
        : notification
    );
    setNotifications(updatedNotifications);
    onNotificationUpdate(updatedNotifications);
    toast.success(`Marked ${notificationIds.length} notification(s) as read`);
  };

  const markAsUnread = (notificationIds: string[]) => {
    const updatedNotifications = notifications.map(notification => 
      notificationIds.includes(notification.id) 
        ? { ...notification, read: false }
        : notification
    );
    setNotifications(updatedNotifications);
    onNotificationUpdate(updatedNotifications);
    toast.success(`Marked ${notificationIds.length} notification(s) as unread`);
  };

  const deleteNotifications = (notificationIds: string[]) => {
    const updatedNotifications = notifications.filter(notification => 
      !notificationIds.includes(notification.id)
    );
    setNotifications(updatedNotifications);
    onNotificationUpdate(updatedNotifications);
    toast.success(`Deleted ${notificationIds.length} notification(s)`);
    setSelectedNotifications([]);
    setIsSelectMode(false);
  };

  const markAllAsRead = () => {
    const filteredNotifications = getFilteredNotifications();
    const unreadIds = filteredNotifications.filter(n => !n.read).map(n => n.id);
    if (unreadIds.length > 0) {
      markAsRead(unreadIds);
    }
  };

  const toggleNotificationSelection = (notificationId: string) => {
    setSelectedNotifications(prev => 
      prev.includes(notificationId)
        ? prev.filter(id => id !== notificationId)
        : [...prev, notificationId]
    );
  };

  const selectAllNotifications = () => {
    const filteredNotifications = getFilteredNotifications();
    setSelectedNotifications(filteredNotifications.map(n => n.id));
  };

  const clearSelection = () => {
    setSelectedNotifications([]);
    setIsSelectMode(false);
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

  const getFilteredNotifications = () => {
    let filtered = notifications;

    // Filter by tab
    switch (selectedTab) {
      case 'unread':
        filtered = filtered.filter(n => !n.read);
        break;
      case 'messages':
        filtered = filtered.filter(n => n.type === 'message');
        break;
      case 'payments':
        filtered = filtered.filter(n => n.type === 'payment');
        break;
      case 'social':
        filtered = filtered.filter(n => n.type === 'social');
        break;
    }

    // Filter by search
    if (searchQuery) {
      filtered = filtered.filter(notification =>
        notification.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        notification.description.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    // Sort by timestamp (newest first)
    return filtered.sort((a, b) => {
      if (!a.read && b.read) return -1;
      if (a.read && !b.read) return 1;
      return 0;
    });
  };

  const filteredNotifications = getFilteredNotifications();
  const unreadCount = notifications.filter(n => !n.read).length;
  const currencySymbol = currencies.find(c => c.code === currentCurrency)?.symbol || '฿';

  return (
    <motion.div
      className="h-screen mobile-full-height bg-background flex flex-col min-w-[320px] max-w-full"
      variants={containerVariants}
      initial="initial"
      animate="animate"
      exit="exit"
    >
      {/* Header */}
      <motion.div 
        className="sticky top-0 z-50 bg-card border-b px-4 py-3 flex items-center gap-3"
        variants={itemVariants}
      >
        <Button
          variant="ghost"
          size="icon"
          className="flex-shrink-0 touch-manipulation"
          onClick={onBack}
        >
          <ArrowLeft className="w-5 h-5" />
        </Button>

        <div className="flex-1 min-w-0">
          <h1 className="text-lg font-semibold truncate">Notifications</h1>
          {unreadCount > 0 && (
            <p className="text-sm text-muted-foreground">
              {unreadCount} unread
            </p>
          )}
        </div>

        {/* Header Actions */}
        <div className="flex items-center gap-2">
          {isSelectMode ? (
            <>
              <Button
                variant="ghost"
                size="sm"
                onClick={clearSelection}
                className="text-xs"
              >
                Cancel
              </Button>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon" className="touch-manipulation">
                    <MoreVertical className="w-4 h-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-48">
                  <DropdownMenuItem onClick={selectAllNotifications}>
                    <Check className="w-4 h-4 mr-2" />
                    Select All
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem 
                    onClick={() => markAsRead(selectedNotifications)}
                    disabled={selectedNotifications.length === 0}
                  >
                    <Bell className="w-4 h-4 mr-2" />
                    Mark as Read
                  </DropdownMenuItem>
                  <DropdownMenuItem 
                    onClick={() => markAsUnread(selectedNotifications)}
                    disabled={selectedNotifications.length === 0}
                  >
                    <BellOff className="w-4 h-4 mr-2" />
                    Mark as Unread
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem 
                    onClick={() => deleteNotifications(selectedNotifications)}
                    disabled={selectedNotifications.length === 0}
                    className="text-destructive"
                  >
                    <Trash2 className="w-4 h-4 mr-2" />
                    Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </>
          ) : (
            <>
              <Button
                variant="ghost"
                size="sm"
                onClick={markAllAsRead}
                disabled={unreadCount === 0}
                className="text-xs touch-manipulation"
              >
                Mark all read
              </Button>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon" className="touch-manipulation">
                    <MoreVertical className="w-4 h-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-48">
                  <DropdownMenuItem onClick={() => setIsSelectMode(true)}>
                    <Check className="w-4 h-4 mr-2" />
                    Select Multiple
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem>
                    <Settings className="w-4 h-4 mr-2" />
                    Notification Settings
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </>
          )}
        </div>
      </motion.div>

      {/* Search Bar */}
      <motion.div 
        className="px-4 py-3 border-b"
        variants={itemVariants}
      >
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
          <Input
            placeholder="Search notifications..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
      </motion.div>

      {/* Filter Tabs */}
      <motion.div variants={itemVariants}>
        <Tabs value={selectedTab} onValueChange={(value) => setSelectedTab(value as typeof selectedTab)}>
          <div className="px-4 py-2 border-b">
            <TabsList className="w-full justify-start h-auto p-1 bg-muted">
              <TabsTrigger value="all" className="text-xs px-3 py-2">
                All ({notifications.length})
              </TabsTrigger>
              <TabsTrigger value="unread" className="text-xs px-3 py-2">
                Unread ({unreadCount})
              </TabsTrigger>
              <TabsTrigger value="messages" className="text-xs px-3 py-2">
                Messages ({notifications.filter(n => n.type === 'message').length})
              </TabsTrigger>
              <TabsTrigger value="payments" className="text-xs px-3 py-2">
                Payments ({notifications.filter(n => n.type === 'payment').length})
              </TabsTrigger>
              <TabsTrigger value="social" className="text-xs px-3 py-2">
                Social ({notifications.filter(n => n.type === 'social').length})
              </TabsTrigger>
            </TabsList>
          </div>

          {/* Notifications List */}
          <TabsContent value={selectedTab} className="flex-1 m-0">
            <ScrollArea className="h-full">
              <AnimatePresence mode="popLayout">
                {filteredNotifications.length > 0 ? (
                  <div className="divide-y">
                    {filteredNotifications.map((notification, index) => (
                      <motion.div
                        key={notification.id}
                        variants={itemVariants}
                        initial="initial"
                        animate="animate"
                        exit="exit"
                        layout
                        transition={{
                          layout: { duration: 0.2, ease: [0.22, 1, 0.36, 1] }
                        }}
                        className={`relative ${!notification.read ? 'bg-accent/30' : ''}`}
                      >
                        <div
                          className={`p-4 cursor-pointer hover:bg-accent/50 transition-colors touch-manipulation ${
                            selectedNotifications.includes(notification.id) ? 'bg-primary/10' : ''
                          }`}
                          onClick={() => handleNotificationClick(notification.id)}
                        >
                          <div className="flex items-start gap-3">
                            {/* Selection Checkbox */}
                            {isSelectMode && (
                              <Checkbox
                                checked={selectedNotifications.includes(notification.id)}
                                onChange={() => toggleNotificationSelection(notification.id)}
                                className="mt-1"
                              />
                            )}

                            {/* Avatar/Icon */}
                            <div className="flex-shrink-0">
                              {notification.avatar ? (
                                <Avatar className="w-10 h-10">
                                  <AvatarImage src={notification.avatar} />
                                  <AvatarFallback>
                                    {notification.title.charAt(0)}
                                  </AvatarFallback>
                                </Avatar>
                              ) : (
                                <div className="w-10 h-10 bg-muted rounded-full flex items-center justify-center">
                                  {getNotificationIcon(notification.type)}
                                </div>
                              )}
                            </div>

                            {/* Content */}
                            <div className="flex-1 min-w-0">
                              <div className="flex items-start justify-between gap-2 mb-1">
                                <h3 className={`text-sm ${!notification.read ? 'font-semibold' : 'font-medium'} truncate`}>
                                  {notification.title}
                                </h3>
                                <div className="flex items-center gap-2 flex-shrink-0">
                                  {!notification.read && (
                                    <div className="w-2 h-2 bg-primary rounded-full"></div>
                                  )}
                                  <span className="text-xs text-muted-foreground whitespace-nowrap">
                                    {notification.timestamp}
                                  </span>
                                </div>
                              </div>
                              
                              <p className="text-sm text-muted-foreground line-clamp-2 mb-2">
                                {notification.description}
                              </p>

                              {/* Amount Badge */}
                              {notification.amount && (
                                <Badge variant="secondary" className="text-xs">
                                  {currencySymbol}{notification.amount.toLocaleString()}
                                </Badge>
                              )}
                            </div>
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                ) : (
                  <motion.div 
                    className="text-center py-12"
                    variants={itemVariants}
                  >
                    <Bell className="w-16 h-16 text-muted-foreground mx-auto mb-4" />
                    <h3 className="text-lg font-medium mb-2">No notifications</h3>
                    <p className="text-sm text-muted-foreground px-4">
                      {searchQuery 
                        ? `No notifications found for "${searchQuery}"`
                        : selectedTab === 'unread' 
                          ? "You're all caught up! No unread notifications."
                          : "No notifications in this category yet."
                      }
                    </p>
                  </motion.div>
                )}
              </AnimatePresence>
            </ScrollArea>
          </TabsContent>
        </Tabs>
      </motion.div>

      {/* Selection Footer */}
      <AnimatePresence>
        {isSelectMode && selectedNotifications.length > 0 && (
          <motion.div
            initial={{ y: 100, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            exit={{ y: 100, opacity: 0 }}
            transition={{ duration: 0.3, ease: [0.22, 1, 0.36, 1] }}
            className="sticky bottom-0 bg-card border-t px-4 py-3"
          >
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">
                {selectedNotifications.length} selected
              </span>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => markAsRead(selectedNotifications)}
                  className="touch-manipulation"
                >
                  Mark Read
                </Button>
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={() => deleteNotifications(selectedNotifications)}
                  className="touch-manipulation"
                >
                  Delete
                </Button>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  );
}