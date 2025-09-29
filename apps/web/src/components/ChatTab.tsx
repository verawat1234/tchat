import React, { useState, useRef, useCallback } from 'react';
import { Search, Plus, Phone, Video, MoreVertical, Paperclip, Mic, Send, Bot, Star, Users, Lock, Speaker, File, Image, Camera, Play, Pause, Download, Tag, Target, Briefcase, Crown, Zap, HardDrive, Upload, AlertTriangle, CheckCircle, Filter, Hash, Building2, ShoppingBag, Trophy, Gauge, BarChart3, FolderOpen, Settings, MessageCircle, TrendingUp, Brain, UserPlus, Sparkles, Wrench, DollarSign, Clock, Activity, Coffee, ArrowDown, Store, Reply, Forward, Copy, Edit3, Trash2, Pin, Archive, VolumeX, Volume2, UserX, RotateCcw, Smile, MapPin, Gift, CreditCard, Bookmark } from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Separator } from './ui/separator';
import { Card, CardContent, CardHeader } from './ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Progress } from './ui/progress';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from './ui/dropdown-menu';
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuLabel, ContextMenuSeparator, ContextMenuTrigger } from './ui/context-menu';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog';
import { Textarea } from './ui/textarea';
import { ChatMessages, ChatHeaderActions } from './ChatActions';
import { AdvancedChatMessages } from './EnhancedChatMessages';
import { ChatInput } from './ChatInput';
import { WorkspaceSwitcher } from './WorkspaceSwitcher';
import { toast } from "sonner";
import { useGetUserChatsQuery } from '../services/microservicesApi';
import { callService } from '../services/webrtc/CallService';

interface ChatTabProps {
  user: any;
  onVideoCall?: (callee: any, isIncoming?: boolean) => void;
  onVoiceCall?: (callee: any, isIncoming?: boolean) => void;
  onNewChat?: () => void;
  selectedWorkspace?: string;
  onWorkspaceChange?: (workspaceId: string) => void;
  userWorkspaces?: any[];
}

interface Dialog {
  id: string;
  name: string;
  lastMessage: string;
  timestamp: string;
  unread: number;
  type: 'user' | 'group' | 'channel' | 'bot' | 'business';
  avatar?: string;
  isOnline?: boolean;
  isPinned?: boolean;
  isE2E?: boolean;
  leadScore?: number;
  tags?: string[];
  businessType?: 'shop' | 'customer' | 'partner' | 'supplier';
  revenue?: number;
  lastActivity?: string;
}

interface Message {
  id: string;
  senderId: string;
  senderName: string;
  content: string;
  timestamp: string;
  type: 'text' | 'payment' | 'system' | 'voice' | 'file' | 'image';
  isOwn: boolean;
  fileUrl?: string;
  fileName?: string;
  fileSize?: string;
  duration?: string;
}

interface WorkBot {
  id: string;
  name: string;
  description: string;
  type: 'lead-qualifier' | 'customer-support' | 'order-manager' | 'analytics' | 'scheduler';
  isActive: boolean;
  conversations: number;
  successRate: number;
  createdDate: string;
}

export function ChatTab({ 
  user, 
  onVideoCall, 
  onVoiceCall, 
  onNewChat, 
  selectedWorkspace = 'golden-mango-shop',
  onWorkspaceChange,
  userWorkspaces = []
}: ChatTabProps) {
  const [selectedDialog, setSelectedDialog] = useState<string | null>('golden-mango');
  const [messageInput, setMessageInput] = useState('');

  // Fetch user chats from API
  const {
    data: userChats,
    isLoading: chatsLoading,
    error: chatsError
  } = useGetUserChatsQuery();
  const [isRecording, setIsRecording] = useState(false);
  const [isSpeakerOn, setIsSpeakerOn] = useState(false);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [showTagFilter, setShowTagFilter] = useState(false);
  const [replyToMessage, setReplyToMessage] = useState<Message | null>(null);
  const [editingMessage, setEditingMessage] = useState<Message | null>(null);
  const [showAttachmentMenu, setShowAttachmentMenu] = useState(false);
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const [selectedMessages, setSelectedMessages] = useState<string[]>([]);
  const [isSelectionMode, setIsSelectionMode] = useState(false);
  const [chatActions, setChatActions] = useState({
    isMuted: false,
    isPinned: false,
    isArchived: false,
    isBlocked: false
  });
  
  const fileInputRef = useRef<HTMLInputElement>(null);
  const imageInputRef = useRef<HTMLInputElement>(null);
  const voiceRecorderRef = useRef<MediaRecorder | null>(null);
  const [recordingDuration, setRecordingDuration] = useState(0);
  const recordingTimer = useRef<NodeJS.Timeout | null>(null);


  // Personal Chat Data (All conversations: friends, family, AND chats with shops as customer)
  const personalDialogs: Dialog[] = [
    {
      id: 'family',
      name: 'Family Group',
      lastMessage: 'Mom: Dinner at 7pm! 🍽️',
      timestamp: '5 min',
      unread: 3,
      type: 'group',
      isOnline: true,
      tags: ['Family', 'Personal'],
      avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face'
    },
    {
      id: 'sarah',
      name: 'Sarah',
      lastMessage: 'See you at the coffee shop! ☕',
      timestamp: '15 min',
      unread: 1,
      type: 'user',
      isOnline: true,
      tags: ['Friends'],
      avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face'
    },
    {
      id: 'pad-thai-shop',
      name: 'Pad Thai Corner',
      lastMessage: 'Your order #247 is ready for pickup! 🍜',
      timestamp: '30 min',
      unread: 1,
      type: 'business',
      isOnline: true,
      tags: ['Restaurant', 'Orders'],
      avatar: 'https://images.unsplash.com/photo-1559847844-5315695b6a77?w=150&h=150&fit=crop'
    },
    {
      id: 'mike',
      name: 'Mike Chen',
      lastMessage: 'Thanks for the pad thai recommendation!',
      timestamp: '1 hour',
      unread: 0,
      type: 'user',
      isOnline: false,
      tags: ['Friends'],
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face'
    },
    {
      id: 'grocery-mart',
      name: 'FreshMart Grocery',
      lastMessage: 'Hi! Your weekly groceries are 20% off today 🛒',
      timestamp: '2 hours',
      unread: 0,
      type: 'business',
      isOnline: false,
      tags: ['Shopping', 'Groceries'],
      avatar: 'https://images.unsplash.com/photo-1542838132-92c53300491e?w=150&h=150&fit=crop'
    },
    {
      id: 'thai-news',
      name: 'Thailand News',
      lastMessage: 'Breaking: New digital payment regulations announced',
      timestamp: '6 hours',
      unread: 0,
      type: 'channel',
      tags: ['News', 'Updates']
    }
  ];

  // Business Chat Data by Workspace (ONLY for shop owners managing customer support)
  const businessDialogsByWorkspace: Record<string, Dialog[]> = {
    'golden-mango-shop': [
      {
        id: 'customer-jane',
        name: 'Jane Wilson',
        lastMessage: 'Do you deliver to Sukhumvit area? 📍',
        timestamp: '5 min',
        unread: 2,
        type: 'business',
        leadScore: 82,
        businessType: 'customer',
        revenue: 1200,
        tags: ['New Customer', 'Delivery Inquiry'],
        avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face',
        lastActivity: 'First-time customer'
      },
      {
        id: 'customer-david',
        name: 'David Kim',
        lastMessage: 'Great service! Will order again next week 👍',
        timestamp: '15 min',
        unread: 0,
        type: 'business',
        leadScore: 95,
        businessType: 'customer',
        revenue: 8500,
        tags: ['VIP Customer', 'Repeat Orders', 'High Value'],
        avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
        lastActivity: 'Ordered 12 times this month'
      },
      {
        id: 'support-bot-golden',
        name: 'Restaurant Support Bot',
        lastMessage: 'Handling 3 customer inquiries automatically',
        timestamp: '45 min',
        unread: 0,
        type: 'bot',
        leadScore: 95,
        tags: ['AI', 'Support', 'Restaurant'],
        lastActivity: 'Active now'
      }
    ],
    'thai-coffee-chain': [
      {
        id: 'coffee-customer-lisa',
        name: 'Lisa Chen',
        lastMessage: 'Can I order 50 coffees for office meeting? ☕',
        timestamp: '10 min',
        unread: 1,
        type: 'business',
        leadScore: 75,
        businessType: 'customer',
        revenue: 2800,
        tags: ['Bulk Orders', 'Corporate'],
        avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop&crop=face',
        lastActivity: 'Bulk order inquiry'
      },
      {
        id: 'coffee-customer-tom',
        name: 'Tom Anderson',
        lastMessage: 'Love the new seasonal blend! 🍂',
        timestamp: '25 min',
        unread: 0,
        type: 'business',
        leadScore: 88,
        businessType: 'customer',
        revenue: 450,
        tags: ['Regular Customer', 'Seasonal'],
        avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
        lastActivity: 'Weekly regular'
      },
      {
        id: 'coffee-bot',
        name: 'Coffee Support Bot',
        lastMessage: 'Handling loyalty program inquiries',
        timestamp: '1 hour',
        unread: 0,
        type: 'bot',
        leadScore: 90,
        tags: ['AI', 'Support', 'Coffee'],
        lastActivity: 'Active now'
      }
    ],
    'sea-tech-company': [
      {
        id: 'tech-client-innovate',
        name: 'InnovateCorp',
        lastMessage: 'Need help with API integration 🔧',
        timestamp: '3 min',
        unread: 3,
        type: 'business',
        leadScore: 92,
        businessType: 'customer',
        revenue: 0, // Not applicable for tech support
        tags: ['Enterprise', 'API Support', 'High Priority'],
        avatar: 'https://images.unsplash.com/photo-1486312338219-ce68d2c6f44d?w=150&h=150&fit=crop',
        lastActivity: 'Technical issue reported'
      },
      {
        id: 'tech-client-startup',
        name: 'StartupXYZ',
        lastMessage: 'Thanks for the quick resolution! 🚀',
        timestamp: '20 min',
        unread: 0,
        type: 'business',
        leadScore: 85,
        businessType: 'customer',
        revenue: 0,
        tags: ['Startup', 'Resolved', 'Satisfied'],
        avatar: 'https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=150&h=150&fit=crop',
        lastActivity: 'Issue resolved'
      },
      {
        id: 'tech-support-bot',
        name: 'Tech Support Assistant',
        lastMessage: 'Automated 15 ticket responses today',
        timestamp: '30 min',
        unread: 0,
        type: 'bot',
        leadScore: 95,
        tags: ['AI', 'Tech Support', 'Automation'],
        lastActivity: 'Active now'
      }
    ]
  };

  const businessDialogs = businessDialogsByWorkspace[selectedWorkspace] || [];

  // Work Bots Data
  const workBots: WorkBot[] = [
    {
      id: 'lead-qualifier',
      name: 'Lead Qualifier',
      description: 'Automatically qualifies incoming leads and assigns scores',
      type: 'lead-qualifier',
      isActive: true,
      conversations: 156,
      successRate: 87,
      createdDate: '2024-01-15'
    },
    {
      id: 'customer-support',
      name: 'Customer Support',
      description: 'Handles common customer inquiries and escalates complex issues',
      type: 'customer-support',
      isActive: true,
      conversations: 342,
      successRate: 92,
      createdDate: '2024-01-10'
    },
    {
      id: 'order-manager',
      name: 'Order Manager',
      description: 'Processes orders and updates customers on status',
      type: 'order-manager',
      isActive: false,
      conversations: 89,
      successRate: 94,
      createdDate: '2024-01-20'
    }
  ];

  const availableTags = [
    'High Priority', 'VIP', 'New Customer', 'Hot Lead', 'Cold Lead',
    'Food Business', 'Tech Support', 'Family', 'Friends', 'Work',
    'AI', 'Support', 'Sales', 'Automation', 'Partner'
  ];

  // Voice and Video Call Handlers
  const handleVoiceCall = useCallback(async () => {
    if (!selectedDialog) return;

    // Find the selected dialog from all available dialogs
    const allDialogs = [...personalDialogs, ...businessDialogs];
    const currentDialog = allDialogs.find(d => d.id === selectedDialog);

    if (!currentDialog) {
      toast.error('Unable to find contact information for call');
      return;
    }

    try {
      // Initialize the call service if not already done
      await callService.initialize();

      // Initiate voice call
      const callId = await callService.initiateCall(currentDialog.id, currentDialog.name, {
        type: 'voice',
        enableVideo: false,
        enableAudio: true,
        quality: 'auto'
      });

      // Call the provided onVoiceCall handler if available
      onVoiceCall?.(currentDialog.name, false);

      // Show success feedback
      toast.success(`Starting voice call with ${currentDialog.name}...`);

    } catch (error) {
      console.error('Failed to start voice call:', error);
      toast.error('Failed to start voice call. Please try again.');
    }
  }, [selectedDialog, personalDialogs, businessDialogs, onVoiceCall]);

  const handleVideoCall = useCallback(async () => {
    if (!selectedDialog) return;

    // Find the selected dialog from all available dialogs
    const allDialogs = [...personalDialogs, ...businessDialogs];
    const currentDialog = allDialogs.find(d => d.id === selectedDialog);

    if (!currentDialog) {
      toast.error('Unable to find contact information for call');
      return;
    }

    try {
      // Initialize the call service if not already done
      await callService.initialize();

      // Initiate video call
      const callId = await callService.initiateCall(currentDialog.id, currentDialog.name, {
        type: 'video',
        enableVideo: true,
        enableAudio: true,
        quality: 'auto'
      });

      // Call the provided onVideoCall handler if available
      onVideoCall?.(currentDialog.name, true);

      // Show success feedback
      toast.success(`Starting video call with ${currentDialog.name}...`);

    } catch (error) {
      console.error('Failed to start video call:', error);
      toast.error('Failed to start video call. Please try again.');
    }
  }, [selectedDialog, personalDialogs, businessDialogs, onVideoCall]);

  // Sample messages for family chat
  const familyMessages: Message[] = [
    {
      id: '1',
      senderId: 'mom',
      senderName: 'Mom',
      content: 'ลูกกินข้าวหรือยัง? (Have you eaten?)',
      timestamp: '10:25 AM',
      type: 'text',
      isOwn: false
    },
    {
      id: '2',
      senderId: user?.id || 'user',
      senderName: 'You',
      content: 'กินแล้วครับ แม่ (Already ate, Mom)',
      timestamp: '10:27 AM',
      type: 'text',
      isOwn: true
    },
    {
      id: '3',
      senderId: 'mom',
      senderName: 'Mom',
      content: 'Dinner at 7pm! 🍽️',
      timestamp: '10:30 AM',
      type: 'text',
      isOwn: false
    },
    {
      id: '4',
      senderId: user?.id || 'user',
      senderName: 'You',
      content: 'Voice message',
      timestamp: '10:31 AM',
      type: 'voice',
      isOwn: true,
      duration: '0:12'
    }
  ];

  // Sample messages for business chat - Golden Mango
  const businessMessages: Message[] = [
    {
      id: '1',
      senderId: 'ai-bot',
      senderName: 'AI Assistant',
      content: '🤖 Lead qualified: Golden Mango Restaurant (Score: 92/100)\n• High revenue potential: ฿25K/month\n• Regular customer (3 orders this week)\n• VIP status confirmed',
      timestamp: '9:00 AM',
      type: 'system',
      isOwn: false
    },
    {
      id: '2',
      senderId: 'golden-mango',
      senderName: 'Golden Mango Restaurant',
      content: 'สวัสดีครับ! We have a new catering menu available for corporate events. 🥭',
      timestamp: '9:15 AM',
      type: 'text',
      isOwn: false
    },
    {
      id: '3',
      senderId: 'golden-mango',
      senderName: 'Golden Mango Restaurant',
      content: 'Perfect for office parties and meetings. 50-500 people capacity.',
      timestamp: '9:16 AM',
      type: 'text',
      isOwn: false
    },
    {
      id: '4',
      senderId: user?.id || 'user',
      senderName: 'You',
      content: 'That sounds great! Can you send me the menu and pricing details?',
      timestamp: '9:20 AM',
      type: 'text',
      isOwn: true
    },
    {
      id: '5',
      senderId: 'golden-mango',
      senderName: 'Golden Mango Restaurant',
      content: 'Catering_Menu_2024.pdf',
      timestamp: '9:22 AM',
      type: 'file',
      isOwn: false,
      fileName: 'Catering_Menu_2024.pdf',
      fileSize: '2.1 MB'
    },
    {
      id: '6',
      senderId: 'ai-bot',
      senderName: 'AI Assistant',
      content: '💡 Suggestion: This customer has 92% lead score. Consider offering:\n• 10% discount for bulk orders\n• Priority delivery\n• Dedicated account manager',
      timestamp: '9:25 AM',
      type: 'system',
      isOwn: false
    },
    {
      id: '7',
      senderId: user?.id || 'user',
      senderName: 'You',
      content: 'We can offer 10% discount for orders over ฿5,000 and free delivery within 5km. Interested?',
      timestamp: '9:30 AM',
      type: 'text',
      isOwn: true
    },
    {
      id: '8',
      senderId: 'golden-mango',
      senderName: 'Golden Mango Restaurant',
      content: 'Excellent! We typically order ฿8,000-฿12,000 per event. When can we schedule a tasting session?',
      timestamp: '9:35 AM',
      type: 'text',
      isOwn: false
    }
  ];

  const handleSendMessage = () => {
    if (!messageInput.trim()) return;
    
    if (editingMessage) {
      toast.success('Message updated!');
      setEditingMessage(null);
    } else if (replyToMessage) {
      toast.success('Reply sent!');
      setReplyToMessage(null);
    } else {
      toast.success('Message sent!');
    }
    
    setMessageInput('');
  };

  // Message Actions
  const handleReplyToMessage = (message: Message) => {
    setReplyToMessage(message);
    setIsSelectionMode(false);
  };

  const handleEditMessage = (message: Message) => {
    if (message.isOwn) {
      setEditingMessage(message);
      setMessageInput(message.content);
    }
  };

  const handleDeleteMessage = (messageId: string) => {
    toast.success('Message deleted');
  };

  const handleForwardMessage = (messageId: string) => {
    setIsSelectionMode(false);
    toast.success('Forward message dialog would open');
  };

  const handleCopyMessage = (content: string) => {
    navigator.clipboard.writeText(content);
    toast.success('Message copied to clipboard');
  };

  const handlePinMessage = (messageId: string) => {
    toast.success('Message pinned');
  };

  const handleSelectMessage = (messageId: string) => {
    if (isSelectionMode) {
      setSelectedMessages(prev => 
        prev.includes(messageId) 
          ? prev.filter(id => id !== messageId)
          : [...prev, messageId]
      );
    }
  };

  const toggleSelectionMode = () => {
    setIsSelectionMode(!isSelectionMode);
    setSelectedMessages([]);
  };

  const clearSelection = () => {
    setSelectedMessages([]);
    setIsSelectionMode(false);
  };

  // File and Media Actions
  const handleFileUpload = () => {
    fileInputRef.current?.click();
  };

  const handleImageUpload = () => {
    imageInputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      toast.success(`Uploading ${file.name}...`);
      setShowAttachmentMenu(false);
    }
  };

  const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      toast.success(`Uploading image...`);
      setShowAttachmentMenu(false);
    }
  };

  const handleLocationShare = () => {
    toast.success('Sharing current location...');
    setShowAttachmentMenu(false);
  };

  const handlePaymentSend = () => {
    toast.success('Payment dialog would open');
    setShowAttachmentMenu(false);
  };

  // Voice Recording
  const startVoiceRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const mediaRecorder = new MediaRecorder(stream);
      voiceRecorderRef.current = mediaRecorder;
      
      mediaRecorder.start();
      setIsRecording(true);
      setRecordingDuration(0);
      
      recordingTimer.current = setInterval(() => {
        setRecordingDuration(prev => prev + 1);
      }, 1000);
      
      toast.success('Recording started');
    } catch (error) {
      toast.error('Microphone access denied');
    }
  };

  const stopVoiceRecording = () => {
    if (voiceRecorderRef.current && isRecording) {
      voiceRecorderRef.current.stop();
      voiceRecorderRef.current.stream.getTracks().forEach(track => track.stop());
      
      setIsRecording(false);
      setRecordingDuration(0);
      
      if (recordingTimer.current) {
        clearInterval(recordingTimer.current);
      }
      
      toast.success('Voice message recorded');
    }
  };

  const cancelVoiceRecording = () => {
    if (voiceRecorderRef.current && isRecording) {
      voiceRecorderRef.current.stop();
      voiceRecorderRef.current.stream.getTracks().forEach(track => track.stop());
      
      setIsRecording(false);
      setRecordingDuration(0);
      
      if (recordingTimer.current) {
        clearInterval(recordingTimer.current);
      }
      
      toast.success('Recording cancelled');
    }
  };

  // Chat Management Actions
  const handleMuteChat = () => {
    setChatActions(prev => ({ ...prev, isMuted: !prev.isMuted }));
    toast.success(chatActions.isMuted ? 'Chat unmuted' : 'Chat muted');
  };

  const handlePinChat = () => {
    setChatActions(prev => ({ ...prev, isPinned: !prev.isPinned }));
    toast.success(chatActions.isPinned ? 'Chat unpinned' : 'Chat pinned');
  };

  const handleArchiveChat = () => {
    setChatActions(prev => ({ ...prev, isArchived: !prev.isArchived }));
    toast.success('Chat archived');
  };

  const handleBlockContact = () => {
    setChatActions(prev => ({ ...prev, isBlocked: !prev.isBlocked }));
    toast.success(chatActions.isBlocked ? 'Contact unblocked' : 'Contact blocked');
  };

  const handleClearHistory = () => {
    toast.success('Chat history cleared');
  };

  const formatRecordingTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const handleCreateBot = (type: string) => {
    toast.success(`Creating ${type} bot...`);
  };

  const handleToggleBot = (botId: string) => {
    toast.success('Bot status updated!');
  };

  const getLeadScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600 bg-green-100';
    if (score >= 60) return 'text-yellow-600 bg-yellow-100';
    if (score >= 40) return 'text-orange-600 bg-orange-100';
    return 'text-red-600 bg-red-100';
  };

  const getLeadScoreIcon = (score: number) => {
    if (score >= 80) return <Zap className="w-3 h-3" />;
    if (score >= 60) return <Target className="w-3 h-3" />;
    if (score >= 40) return <BarChart3 className="w-3 h-3" />;
    return <Gauge className="w-3 h-3" />;
  };

  const getBotTypeIcon = (type: string) => {
    switch (type) {
      case 'lead-qualifier': return <Target className="w-4 h-4 text-green-500" />;
      case 'customer-support': return <Users className="w-4 h-4 text-blue-500" />;
      case 'order-manager': return <ShoppingBag className="w-4 h-4 text-purple-500" />;
      case 'analytics': return <BarChart3 className="w-4 h-4 text-orange-500" />;
      case 'scheduler': return <Clock className="w-4 h-4 text-red-500" />;
      default: return <Bot className="w-4 h-4 text-gray-500" />;
    }
  };

  const formatRevenue = (amount: number) => {
    if (amount >= 1000000) return `฿${(amount / 1000000).toFixed(1)}M`;
    if (amount >= 1000) return `฿${(amount / 1000).toFixed(0)}K`;
    return `฿${amount}`;
  };

  const renderDialogList = (dialogs: Dialog[], isWork = false) => (
    <div className="p-2">
      {/* Loading State */}
      {chatsLoading && (
        <div className="flex items-center justify-center py-8">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-sm text-muted-foreground">Loading chats...</p>
          </div>
        </div>
      )}

      {/* Error State */}
      {chatsError && (
        <div className="flex items-center justify-center py-8">
          <div className="text-center">
            <AlertTriangle className="w-8 h-8 text-destructive mx-auto mb-4" />
            <p className="text-sm text-muted-foreground">Failed to load chats</p>
          </div>
        </div>
      )}

      {/* Chat List */}
      {!chatsLoading && !chatsError && dialogs.map((dialog) => (
        <button
          key={dialog.id}
          onClick={() => setSelectedDialog(dialog.id)}
          className={`w-full p-3 rounded-lg text-left hover:bg-accent/50 transition-colors mb-2 ${
            selectedDialog === dialog.id ? 'bg-accent' : ''
          }`}
        >
          <div className="flex items-start gap-3">
            <div className="relative">
              <Avatar className="w-12 h-12">
                <AvatarImage src={dialog.avatar} />
                <AvatarFallback>
                  {dialog.type === 'bot' ? <Bot className="w-6 h-6" /> :
                   dialog.type === 'group' ? <Users className="w-6 h-6" /> :
                   dialog.type === 'business' ? <Building2 className="w-6 h-6" /> :
                   dialog.name.charAt(0)}
                </AvatarFallback>
              </Avatar>
              {dialog.isOnline && (
                <div className="absolute -bottom-1 -right-1 w-4 h-4 bg-green-500 rounded-full border-2 border-background"></div>
              )}
            </div>
            
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <span className="font-medium truncate">{dialog.name}</span>
                {dialog.isPinned && <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />}
                {dialog.isE2E && <Lock className="w-4 h-4 text-green-500" />}
                {dialog.type === 'bot' && <Bot className="w-4 h-4 text-chart-1" />}
                {dialog.type === 'business' && <Building2 className="w-4 h-4 text-chart-2" />}
              </div>
              
              <p className="text-sm text-muted-foreground truncate mb-1">
                {dialog.lastMessage}
              </p>
              
              {/* Lead Score & Revenue - Only for work tab */}
              {isWork && dialog.leadScore !== undefined && (
                <div className="flex items-center gap-2 mb-1">
                  <div className={`flex items-center gap-1 px-2 py-1 rounded-full text-xs ${getLeadScoreColor(dialog.leadScore)}`}>
                    {getLeadScoreIcon(dialog.leadScore)}
                    <span>{dialog.leadScore}</span>
                  </div>
                  {dialog.revenue !== undefined && dialog.revenue > 0 && (
                    <Badge variant="secondary" className="text-xs">
                      {formatRevenue(dialog.revenue)}
                    </Badge>
                  )}
                </div>
              )}
              
              {/* Tags */}
              {dialog.tags && dialog.tags.length > 0 && (
                <div className="flex flex-wrap gap-1 mt-1">
                  {dialog.tags.slice(0, 2).map(tag => (
                    <Badge key={tag} variant="outline" className="text-xs px-1 py-0">
                      <Hash className="w-2 h-2 mr-1" />
                      {tag}
                    </Badge>
                  ))}
                  {dialog.tags.length > 2 && (
                    <Badge variant="outline" className="text-xs px-1 py-0">
                      +{dialog.tags.length - 2}
                    </Badge>
                  )}
                </div>
              )}
            </div>
            
            <div className="flex flex-col items-end gap-1">
              <span className="text-xs text-muted-foreground">
                {dialog.timestamp}
              </span>
              {dialog.unread > 0 && (
                <Badge className="w-5 h-5 text-xs p-0 flex items-center justify-center">
                  {dialog.unread}
                </Badge>
              )}
            </div>
          </div>
        </button>
      ))}
    </div>
  );

  const renderWorkDashboard = () => (
    <div className="p-4 space-y-6">
      {/* Quick Stats */}
      <div className="grid grid-cols-2 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-green-500" />
              <span className="text-sm font-medium">Today's Leads</span>
            </div>
            <p className="text-2xl font-bold">23</p>
            <p className="text-xs text-muted-foreground">+15% from yesterday</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <DollarSign className="w-4 h-4 text-blue-500" />
              <span className="text-sm font-medium">Revenue</span>
            </div>
            <p className="text-2xl font-bold">฿45K</p>
            <p className="text-xs text-muted-foreground">This week</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Bot className="w-4 h-4 text-purple-500" />
              <span className="text-sm font-medium">Active Bots</span>
            </div>
            <p className="text-2xl font-bold">{workBots.filter(bot => bot.isActive).length}</p>
            <p className="text-xs text-muted-foreground">Running now</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Activity className="w-4 h-4 text-orange-500" />
              <span className="text-sm font-medium">Conversations</span>
            </div>
            <p className="text-2xl font-bold">587</p>
            <p className="text-xs text-muted-foreground">Total this month</p>
          </CardContent>
        </Card>
      </div>
      
      {/* Active Bots */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-medium flex items-center gap-2">
            <Bot className="w-5 h-5 text-chart-1" />
            AI Bots
          </h3>
          <Button size="sm" onClick={() => handleCreateBot('new')}>
            <Plus className="w-4 h-4 mr-1" />
            Create Bot
          </Button>
        </div>
        
        <div className="space-y-3">
          {workBots.map((bot) => (
            <Card key={bot.id}>
              <CardContent className="p-4">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    {getBotTypeIcon(bot.type)}
                    <div>
                      <h4 className="font-medium">{bot.name}</h4>
                      <p className="text-sm text-muted-foreground">{bot.description}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant={bot.isActive ? "default" : "secondary"}>
                      {bot.isActive ? 'Active' : 'Inactive'}
                    </Badge>
                    <Button 
                      variant="ghost" 
                      size="sm"
                      onClick={() => handleToggleBot(bot.id)}
                    >
                      <Settings className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
                
                <div className="grid grid-cols-3 gap-4 text-sm">
                  <div>
                    <p className="text-muted-foreground">Conversations</p>
                    <p className="font-medium">{bot.conversations}</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Success Rate</p>
                    <p className="font-medium">{bot.successRate}%</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">Created</p>
                    <p className="font-medium">{new Date(bot.createdDate).toLocaleDateString()}</p>
                  </div>
                </div>
                
                <div className="mt-3">
                  <div className="flex justify-between text-sm mb-1">
                    <span>Performance</span>
                    <span>{bot.successRate}%</span>
                  </div>
                  <Progress value={bot.successRate} className="h-2" />
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
      
      {/* Quick Actions */}
      <div>
        <h3 className="font-medium mb-4">Quick Actions</h3>
        <div className="grid grid-cols-2 gap-3">
          <Button variant="outline" onClick={() => handleCreateBot('lead-qualifier')}>
            <Target className="w-4 h-4 mr-2" />
            Lead Qualifier Bot
          </Button>
          <Button variant="outline" onClick={() => handleCreateBot('customer-support')}>
            <Users className="w-4 h-4 mr-2" />
            Support Bot
          </Button>
          <Button variant="outline" onClick={() => handleCreateBot('order-manager')}>
            <ShoppingBag className="w-4 h-4 mr-2" />
            Order Manager Bot
          </Button>
          <Button variant="outline" onClick={() => handleCreateBot('analytics')}>
            <BarChart3 className="w-4 h-4 mr-2" />
            Analytics Bot
          </Button>
        </div>
      </div>
    </div>
  );

  return (
    <div className="flex h-full">
      {/* Left Sidebar with Tabs */}
      <div className="w-80 border-r border-border flex flex-col">
        <Tabs defaultValue="chat" className="flex-1 flex flex-col">
          <div className="p-4 border-b border-border">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="chat" className="flex items-center gap-2">
                <MessageCircle className="w-4 h-4" />
                <div className="flex flex-col items-start">
                  <span>Chat</span>
                  <span className="text-xs text-muted-foreground">All conversations</span>
                </div>
              </TabsTrigger>
              <TabsTrigger value="work" className="flex items-center gap-2">
                <Briefcase className="w-4 h-4" />
                <div className="flex flex-col items-start">
                  <span>Work</span>
                  <span className="text-xs text-muted-foreground">Business support</span>
                </div>
              </TabsTrigger>
            </TabsList>
          </div>
          
          {/* Chat Tab Content */}
          <TabsContent value="chat" className="flex-1 flex flex-col mt-0">
            <div className="p-4 border-b border-border">
              <div className="relative mb-3">
                <Search className="absolute left-3 top-3 w-4 h-4 text-muted-foreground" />
                <Input placeholder="Search personal chats..." className="pl-10" />
              </div>
            </div>
            
            <ScrollArea className="flex-1">
              {renderDialogList(personalDialogs, false)}
            </ScrollArea>
            
            <div className="p-4 border-t border-border">
              <Button className="w-full" variant="outline" onClick={onNewChat}>
                <Plus className="w-4 h-4 mr-2" />
                New Chat
              </Button>
            </div>
          </TabsContent>
          
          {/* Work Tab Content */}
          <TabsContent value="work" className="flex-1 flex flex-col mt-0">
            <div className="p-4 border-b border-border">
              {/* Workspace Switcher */}
              {userWorkspaces.length > 1 && (
                <div className="mb-3">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="outline" className="w-full justify-between">
                        <div className="flex items-center gap-2">
                          <Avatar className="w-6 h-6">
                            <AvatarImage src={userWorkspaces.find(w => w.id === selectedWorkspace)?.avatar} />
                            <AvatarFallback>
                              {userWorkspaces.find(w => w.id === selectedWorkspace)?.type === 'shop' ? 
                                <Store className="w-3 h-3" /> : 
                                <Building2 className="w-3 h-3" />
                              }
                            </AvatarFallback>
                          </Avatar>
                          <div className="text-left">
                            <p className="text-sm font-medium">
                              {userWorkspaces.find(w => w.id === selectedWorkspace)?.name}
                            </p>
                            <p className="text-xs text-muted-foreground">
                              {userWorkspaces.find(w => w.id === selectedWorkspace)?.role}
                            </p>
                          </div>
                        </div>
                        <ArrowDown className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent className="w-full">
                      <DropdownMenuLabel>Switch Workspace</DropdownMenuLabel>
                      <DropdownMenuSeparator />
                      {userWorkspaces.map((workspace) => (
                        <DropdownMenuItem
                          key={workspace.id}
                          onClick={() => onWorkspaceChange?.(workspace.id)}
                          className={selectedWorkspace === workspace.id ? 'bg-accent' : ''}
                        >
                          <div className="flex items-center gap-3 w-full">
                            <Avatar className="w-8 h-8">
                              <AvatarImage src={workspace.avatar} />
                              <AvatarFallback>
                                {workspace.type === 'shop' ? 
                                  <Store className="w-4 h-4" /> : 
                                  <Building2 className="w-4 h-4" />
                                }
                              </AvatarFallback>
                            </Avatar>
                            <div className="flex-1">
                              <p className="font-medium">{workspace.name}</p>
                              <p className="text-xs text-muted-foreground">{workspace.role}</p>
                              <div className="flex items-center gap-2 mt-1">
                                <Badge variant="secondary" className="text-xs">
                                  {workspace.customerCount} customers
                                </Badge>
                                {workspace.revenue && (
                                  <Badge variant="outline" className="text-xs">
                                    {workspace.revenue}
                                  </Badge>
                                )}
                              </div>
                            </div>
                          </div>
                        </DropdownMenuItem>
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              )}
              
              <div className="relative mb-3">
                <Search className="absolute left-3 top-3 w-4 h-4 text-muted-foreground" />
                <Input placeholder="Search business chats..." className="pl-10" />
              </div>
              
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowTagFilter(!showTagFilter)}
                  className="flex-1"
                >
                  <Filter className="w-4 h-4 mr-2" />
                  Filter by Tags
                  {selectedTags.length > 0 && (
                    <Badge className="ml-2">{selectedTags.length}</Badge>
                  )}
                </Button>
              </div>
              
              {showTagFilter && (
                <div className="mt-3 p-3 bg-muted rounded-lg">
                  <div className="grid grid-cols-2 gap-2">
                    {availableTags.slice(0, 8).map(tag => (
                      <Button
                        key={tag}
                        variant={selectedTags.includes(tag) ? 'default' : 'outline'}
                        size="sm"
                        className="text-xs h-8"
                        onClick={() => {
                          setSelectedTags(prev => 
                            prev.includes(tag)
                              ? prev.filter(t => t !== tag)
                              : [...prev, tag]
                          );
                        }}
                      >
                        <Hash className="w-3 h-3 mr-1" />
                        {tag}
                      </Button>
                    ))}
                  </div>
                </div>
              )}
            </div>
            
            <ScrollArea className="flex-1">
              {renderDialogList(businessDialogs, true)}
            </ScrollArea>
            
            <div className="p-4 border-t border-border">
              <Button className="w-full" variant="outline" onClick={() => handleCreateBot('new')}>
                <Sparkles className="w-4 h-4 mr-2" />
                Create AI Bot
              </Button>
            </div>
          </TabsContent>
        </Tabs>
      </div>

      {/* Chat Window */}
      <div className="flex-1 flex flex-col">
        {selectedDialog ? (
          <>
            {/* Chat Header */}
            <div className="border-b border-border p-4 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Avatar className="w-10 h-10">
                  <AvatarImage src={personalDialogs.find(d => d.id === selectedDialog)?.avatar || businessDialogs.find(d => d.id === selectedDialog)?.avatar} />
                  <AvatarFallback>
                    {selectedDialog === 'family' ? <Users className="w-5 h-5" /> : 'F'}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <div className="flex items-center gap-2">
                    <h3 className="font-medium">
                      {personalDialogs.find(d => d.id === selectedDialog)?.name || businessDialogs.find(d => d.id === selectedDialog)?.name}
                    </h3>
                    {(() => {
                      const dialog = businessDialogs.find(d => d.id === selectedDialog);
                      if (dialog?.leadScore !== undefined) {
                        return (
                          <div className={`flex items-center gap-1 px-2 py-1 rounded-full text-xs ${getLeadScoreColor(dialog.leadScore)}`}>
                            {getLeadScoreIcon(dialog.leadScore)}
                            <span>{dialog.leadScore}</span>
                          </div>
                        );
                      }
                      return null;
                    })()}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {selectedDialog === 'family' ? 'Family Group • 4 members' : 'Online'}
                  </p>
                </div>
              </div>
              
              <div className="flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleVoiceCall()}
                  aria-label="Start voice call"
                  data-testid="voice-call-button"
                >
                  <Phone className="w-5 h-5" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleVideoCall()}
                  aria-label="Start video call"
                  data-testid="video-call-button"
                >
                  <Video className="w-5 h-5" />
                </Button>
                <Button variant="ghost" size="icon">
                  <MoreVertical className="w-5 h-5" />
                </Button>
              </div>
            </div>

            {/* Messages */}
            <ScrollArea className="flex-1">
              <AdvancedChatMessages
                messages={selectedDialog === 'family' ? familyMessages :
                         selectedDialog === 'golden-mango' ? businessMessages :
                         []}
                selectedMessages={[]}
                isSelectionMode={false}
                onReplyToMessage={(message) => {
                  // Handle reply functionality
                  console.log('Reply to:', message);
                }}
                onEditMessage={(message) => {
                  // Handle edit functionality
                  console.log('Edit:', message);
                }}
                onDeleteMessage={(messageId) => {
                  // Handle delete functionality
                  console.log('Delete:', messageId);
                }}
                onForwardMessage={(messageId) => {
                  // Handle forward functionality
                  console.log('Forward:', messageId);
                }}
                onCopyMessage={(content) => {
                  // Handle copy functionality
                  navigator.clipboard.writeText(content);
                  toast("Message copied to clipboard");
                }}
                onPinMessage={(messageId) => {
                  // Handle pin functionality
                  console.log('Pin:', messageId);
                }}
                onSelectMessage={(messageId) => {
                  // Handle selection functionality
                  console.log('Select:', messageId);
                }}
                chatId={selectedDialog || 'default-chat'}
                currentUserId="current-user"
              />
            </ScrollArea>

            {/* Message Input */}
            <div className="border-t border-border p-4">
              <div className="flex items-center gap-2">
                <Button variant="ghost" size="icon">
                  <Paperclip className="w-5 h-5" />
                </Button>
                <Input
                  placeholder="Type a message..."
                  value={messageInput}
                  onChange={(e) => setMessageInput(e.target.value)}
                  onKeyPress={(e) => {
                    if (e.key === 'Enter' && !e.shiftKey) {
                      e.preventDefault();
                      handleSendMessage();
                    }
                  }}
                  className="flex-1"
                />
                <Button
                  variant="ghost"
                  size="icon"
                  onMouseDown={() => setIsRecording(true)}
                  onMouseUp={() => setIsRecording(false)}
                  onMouseLeave={() => setIsRecording(false)}
                >
                  <Mic className={`w-5 h-5 ${isRecording ? 'text-red-500' : ''}`} />
                </Button>
                <Button 
                  size="icon" 
                  onClick={handleSendMessage}
                  disabled={!messageInput.trim()}
                >
                  <Send className="w-5 h-5" />
                </Button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <Coffee className="w-16 h-16 text-muted-foreground mx-auto mb-4" />
              <h3 className="font-medium mb-2">Select a chat to start messaging</h3>
              <p className="text-sm text-muted-foreground">
                Choose from your personal chats or business conversations
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}