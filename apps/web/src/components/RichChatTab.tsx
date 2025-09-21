import React, { useState, useRef, useCallback } from 'react';
import { Search, Plus, Phone, Video, MoreVertical, MessageCircle, Briefcase, Hash, ArrowLeft } from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Card, CardContent } from './ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { WorkspaceSwitcher } from './WorkspaceSwitcher';
import { MarkdownMessage } from './MarkdownMessage';
import { RichChatInput } from './RichChatInput';
import { ChatHeaderActions } from './ChatActions';
import { toast } from "sonner";

interface ChatTabProps {
  user: any;
  onVideoCall?: (callee: any, isIncoming?: boolean) => void;
  onVoiceCall?: (callee: any, isIncoming?: boolean) => void;
  onNewChat?: () => void;
  selectedWorkspace?: string;
  onWorkspaceChange?: (workspaceId: string) => void;
  userWorkspaces?: any[];
  onShopChatClick?: (shopId: string) => void;
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
  leadScore?: number;
  tags?: string[];
  businessType?: 'shop' | 'customer' | 'partner' | 'supplier';
  revenue?: number;
  lastActivity?: string;
}

interface MessageData {
  id: string;
  senderId: string;
  senderName: string;
  content: string;
  timestamp: string;
  type: 'text' | 'markdown' | 'invoice' | 'bill' | 'order' | 'sticker' | 'poll' | 'contact' | 'location' | 'payment' | 'voice' | 'file' | 'image' | 'video' | 'music' | 'system';
  isOwn: boolean;
  metadata?: any;
}

export function RichChatTab({ 
  user, 
  onVideoCall, 
  onVoiceCall, 
  onNewChat, 
  selectedWorkspace = 'golden-mango-shop',
  onWorkspaceChange,
  userWorkspaces = [],
  onShopChatClick
}: ChatTabProps) {
  const [selectedDialog, setSelectedDialog] = useState<string | null>('family');
  const [messageInput, setMessageInput] = useState('');
  const [isRecording, setIsRecording] = useState(false);
  const [recordingDuration, setRecordingDuration] = useState(0);
  const [replyToMessage, setReplyToMessage] = useState<MessageData | null>(null);
  const [editingMessage, setEditingMessage] = useState<MessageData | null>(null);
  const [messages, setMessages] = useState<Record<string, MessageData[]>>({});
  const [chatActions, setChatActions] = useState({
    isMuted: false,
    isPinned: false,
    isArchived: false,
    isBlocked: false
  });

  const voiceRecorderRef = useRef<MediaRecorder | null>(null);
  const recordingTimer = useRef<NodeJS.Timeout | null>(null);

  // Personal Chat Data
  const personalDialogs: Dialog[] = [
    {
      id: 'family',
      name: 'Family Group',
      lastMessage: 'Payment sent: ฿500 🎉',
      timestamp: '2 min',
      unread: 2,
      type: 'group',
      isOnline: true,
      tags: ['Family', 'Personal'],
      avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face'
    },
    {
      id: 'thai-restaurant',
      name: 'Golden Mango Restaurant',
      lastMessage: 'Order #ORD-2024-001 created for customer',
      timestamp: '15 min',
      unread: 1,
      type: 'business',
      isOnline: true,
      tags: ['Restaurant', 'Business'],
      avatar: 'https://images.unsplash.com/photo-1559847844-5315695b6a77?w=150&h=150&fit=crop'
    },
    {
      id: 'customer-order-chat',
      name: 'Thai Coffee Orders',
      lastMessage: 'Order ready for pickup!',
      timestamp: '45 min',
      unread: 0,
      type: 'business',
      isOnline: false,
      tags: ['Coffee', 'Orders'],
      avatar: 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=150&h=150&fit=crop'
    },
    {
      id: 'poll-group',
      name: 'Weekend Plans',
      lastMessage: 'Poll: Where should we go this weekend?',
      timestamp: '1 hour',
      unread: 0,
      type: 'group',
      isOnline: false,
      tags: ['Friends', 'Social'],
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face'
    }
  ];

  // Rich sample messages with different types
  const sampleMessages: Record<string, MessageData[]> = {
    'family': [
      {
        id: '1',
        senderId: 'mom',
        senderName: 'Mom',
        content: 'ลูกกินข้าวหรือยัง? **Have you eaten?**',
        timestamp: '10:25 AM',
        type: 'markdown',
        isOwn: false
      },
      {
        id: '2',
        senderId: user?.id || 'user',
        senderName: 'You',
        content: 'กินแล้วครับ แม่ *Already ate, Mom* 😊',
        timestamp: '10:27 AM',
        type: 'markdown',
        isOwn: true
      },
      {
        id: '3',
        senderId: 'dad',
        senderName: 'Dad',
        content: 'Here\'s money for groceries',
        timestamp: '10:30 AM',
        type: 'payment',
        isOwn: false,
        metadata: {
          payment: {
            amount: 500,
            currency: 'THB',
            method: 'PromptPay',
            status: 'completed',
            recipient: 'You',
            reference: 'REF123456789'
          }
        }
      },
      {
        id: '4',
        senderId: 'sister',
        senderName: 'Ploy',
        content: '🎉',
        timestamp: '10:32 AM',
        type: 'sticker',
        isOwn: false,
        metadata: {
          sticker: {
            emoji: '🎉',
            size: 'large',
            animation: 'bounce'
          }
        }
      },
      {
        id: '5',
        senderId: user?.id || 'user',
        senderName: 'You',
        content: 'Voice message',
        timestamp: '10:35 AM',
        type: 'voice',
        isOwn: true,
        metadata: {
          media: {
            duration: '0:15',
            url: ''
          }
        }
      }
    ],
    'thai-restaurant': [
      {
        id: '1',
        senderId: 'restaurant',
        senderName: 'Golden Mango Restaurant',
        content: 'New order created for customer delivery',
        timestamp: '1:30 PM',
        type: 'order',
        isOwn: false,
        metadata: {
          order: {
            number: 'ORD-2024-001',
            items: [
              { 
                name: 'Pad Thai with Shrimp', 
                quantity: 2, 
                price: 180, 
                total: 360,
                notes: 'Extra spicy, no bean sprouts',
                image: 'https://images.unsplash.com/photo-1559847844-5315695b6a77?w=100&h=100&fit=crop'
              },
              { 
                name: 'Tom Yum Soup (Large)', 
                quantity: 1, 
                price: 120, 
                total: 120,
                notes: 'Medium spice level',
                image: 'https://images.unsplash.com/photo-1542795110-4eea0ce69c8c?w=100&h=100&fit=crop'
              },
              { 
                name: 'Thai Iced Tea', 
                quantity: 2, 
                price: 50, 
                total: 100,
                notes: 'Less sweet',
                image: 'https://images.unsplash.com/photo-1570968915860-54d5c301fa9f?w=100&h=100&fit=crop'
              }
            ],
            subtotal: 580,
            deliveryFee: 30,
            tax: 40.6,
            total: 650.6,
            currency: 'THB',
            status: 'confirmed',
            estimatedTime: '30-45 min',
            deliveryTime: '2:30 PM',
            customer: {
              name: 'Somchai Pattana',
              phone: '+66 81 234 5678',
              address: '456 Silom Road, Bangkok 10500, Thailand',
              email: 'somchai@email.com'
            },
            shop: {
              name: 'Golden Mango Restaurant',
              address: '123 Sukhumvit Road, Bangkok 10110',
              phone: '+66 2 123 4567'
            },
            deliveryType: 'delivery',
            paymentStatus: 'paid',
            paymentMethod: 'PromptPay',
            createdAt: 'January 25, 2024 1:30 PM',
            notes: 'Customer is allergic to peanuts - please ensure no cross-contamination'
          }
        }
      },
      {
        id: '2',
        senderId: 'restaurant',
        senderName: 'Golden Mango Restaurant',
        content: 'Thank you for your business! Here\'s your invoice for the catering order.',
        timestamp: '2:00 PM',
        type: 'invoice',
        isOwn: false,
        metadata: {
          invoice: {
            number: 'INV-2024-001',
            items: [
              { name: 'Pad Thai Catering (50 servings)', quantity: 1, price: 2500, total: 2500 },
              { name: 'Tom Yum Soup (5 liters)', quantity: 1, price: 800, total: 800 },
              { name: 'Green Curry (5 liters)', quantity: 1, price: 900, total: 900 },
              { name: 'Jasmine Rice (10kg)', quantity: 1, price: 300, total: 300 }
            ],
            subtotal: 4500,
            tax: 315,
            total: 4815,
            currency: 'THB',
            status: 'paid',
            from: {
              name: 'Golden Mango Restaurant',
              address: '123 Sukhumvit Road, Bangkok 10110',
              phone: '+66 2 123 4567'
            },
            to: {
              name: 'Corporate Event Services',
              address: '456 Silom Road, Bangkok 10500',
              phone: '+66 2 987 6543'
            },
            dueDate: 'January 30, 2024'
          }
        }
      },
      {
        id: '3',
        senderId: user?.id || 'user',
        senderName: 'You',
        content: 'Perfect! The food was **amazing** and everyone loved it. Will definitely order again! 🙏',
        timestamp: '2:15 PM',
        type: 'markdown',
        isOwn: true
      },
      {
        id: '4',
        senderId: user?.id || 'user',
        senderName: 'You',
        content: 'Order confirmed! Customer has been notified. **Estimated delivery: 2:30 PM** 🛵',
        timestamp: '1:35 PM',
        type: 'markdown',
        isOwn: true
      },
      {
        id: '5',
        senderId: 'restaurant',
        senderName: 'Golden Mango Restaurant',
        content: 'Dr. Somchai - Head Chef',
        timestamp: '2:20 PM',
        type: 'contact',
        isOwn: false,
        metadata: {
          contact: {
            name: 'Dr. Somchai Kitcharoen',
            phone: '+66 81 234 5678',
            email: 'chef@goldenmango.co.th',
            company: 'Golden Mango Restaurant',
            title: 'Head Chef & Owner',
            avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face'
          }
        }
      }
    ],
    'poll-group': [
      {
        id: '1',
        senderId: 'sarah',
        senderName: 'Sarah',
        content: 'Weekend Plans Poll',
        timestamp: '11:00 AM',
        type: 'poll',
        isOwn: false,
        metadata: {
          poll: {
            question: 'Where should we go this weekend?',
            options: [
              { text: 'Chatuchak Weekend Market 🛍️', votes: 5, voters: ['sarah', 'mike', 'jane', 'tom', 'lisa'] },
              { text: 'Ayutthaya Historical Park 🏛️', votes: 3, voters: ['david', 'anna', 'john'] },
              { text: 'Floating Market 🛶', votes: 7, voters: ['ploy', 'som', 'niran', 'kate', 'alex', 'nina', 'bob'] },
              { text: 'Stay home and chill 🏠', votes: 2, voters: ['lazy1', 'lazy2'] }
            ],
            totalVotes: 17,
            allowMultiple: false,
            anonymous: false
          }
        }
      },
      {
        id: '2',
        senderId: user?.id || 'user',
        senderName: 'You',
        content: 'I vote for the **floating market**! 🛶 Haven\'t been there in ages.',
        timestamp: '11:05 AM',
        type: 'markdown',
        isOwn: true
      },
      {
        id: '3',
        senderId: 'mike',
        senderName: 'Mike',
        content: 'Damnoen Saduak Floating Market',
        timestamp: '11:10 AM',
        type: 'location',
        isOwn: false,
        metadata: {
          location: {
            name: 'Damnoen Saduak Floating Market',
            address: 'Damnoen Saduak District, Ratchaburi 70130, Thailand',
            latitude: 13.5186,
            longitude: 99.9550,
            thumbnail: 'https://images.unsplash.com/photo-1570717191584-95b6fa2f3bb6?w=400&h=200&fit=crop'
          }
        }
      }
    ],
    'customer-order-chat': [
      {
        id: '1',
        senderId: 'coffee-shop',
        senderName: 'Thai Coffee Chain',
        content: 'Order created for office delivery',
        timestamp: '12:15 PM',
        type: 'order',
        isOwn: false,
        metadata: {
          order: {
            number: 'ORD-COF-456',
            items: [
              { 
                name: 'Thai Iced Coffee (Large)', 
                quantity: 8, 
                price: 65, 
                total: 520,
                notes: 'Office delivery - mixed sweetness levels',
                image: 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=100&h=100&fit=crop'
              },
              { 
                name: 'Americano (Medium)', 
                quantity: 4, 
                price: 55, 
                total: 220,
                notes: 'Hot, no sugar',
                image: 'https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=100&h=100&fit=crop'
              },
              { 
                name: 'Croissant', 
                quantity: 6, 
                price: 45, 
                total: 270,
                notes: 'Warm if possible',
                image: 'https://images.unsplash.com/photo-1555507036-ab794f77c9d2?w=100&h=100&fit=crop'
              }
            ],
            subtotal: 1010,
            deliveryFee: 50,
            tax: 70.7,
            total: 1130.7,
            currency: 'THB',
            status: 'ready',
            estimatedTime: '20-30 min',
            deliveryTime: '1:00 PM',
            customer: {
              name: 'Tech Startup Office',
              phone: '+66 2 987 6543',
              address: '789 Innovation Tower, Bangkok 10400',
              email: 'office@techstartup.co.th'
            },
            shop: {
              name: 'Thai Coffee Chain - Silom Branch',
              address: '321 Silom Road, Bangkok 10500',
              phone: '+66 2 555 0123'
            },
            deliveryType: 'delivery',
            paymentStatus: 'paid',
            paymentMethod: 'Corporate Account',
            createdAt: 'January 25, 2024 12:15 PM',
            notes: 'Bulk office order - please call reception when arriving'
          }
        }
      },
      {
        id: '2',
        senderId: user?.id || 'user',
        senderName: 'You',
        content: 'Order is **ready for pickup**! Delivery driver can collect anytime. ☕',
        timestamp: '12:45 PM',
        type: 'markdown',
        isOwn: true
      },
      {
        id: '3',
        senderId: 'coffee-shop',
        senderName: 'Thai Coffee Chain',
        content: 'Perfect! Driver is on the way. Customer will receive delivery notification.',
        timestamp: '12:50 PM',
        type: 'text',
        isOwn: false
      }
    ]
  };

  // Initialize messages
  React.useEffect(() => {
    setMessages(sampleMessages);
  }, []);

  const currentMessages = messages[selectedDialog || ''] || [];
  const currentDialog = personalDialogs.find(d => d.id === selectedDialog);

  const handleSendMessage = (messageData: Partial<MessageData>) => {
    if (!selectedDialog) return;

    const newMessage: MessageData = {
      id: Date.now().toString(),
      senderId: user?.id || 'user',
      senderName: 'You',
      timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      isOwn: true,
      content: messageInput,
      type: 'text',
      ...messageData
    };

    setMessages(prev => ({
      ...prev,
      [selectedDialog]: [...(prev[selectedDialog] || []), newMessage]
    }));

    setMessageInput('');
    setReplyToMessage(null);
    setEditingMessage(null);
    
    toast.success('Message sent!');
  };

  const handleReplyToMessage = (message: MessageData) => {
    setReplyToMessage(message);
  };

  const handleShareMessage = (message: MessageData) => {
    navigator.clipboard.writeText(message.content);
    toast.success('Message copied to clipboard');
  };

  const handleCopyMessage = (content: string) => {
    navigator.clipboard.writeText(content);
    toast.success('Copied to clipboard');
  };

  const handleReactToMessage = (messageId: string, emoji: string) => {
    toast.success(`Reacted with ${emoji}`);
  };

  const handleVoteInPoll = (messageId: string, optionIndex: number) => {
    if (!selectedDialog) return;

    setMessages(prev => ({
      ...prev,
      [selectedDialog]: prev[selectedDialog]?.map(msg => {
        if (msg.id === messageId && msg.metadata?.poll) {
          const poll = { ...msg.metadata.poll };
          const userId = user?.id || 'user';
          
          // Remove user from all options first (for single choice polls)
          if (!poll.allowMultiple) {
            poll.options.forEach((opt: any) => {
              opt.voters = opt.voters.filter((voter: string) => voter !== userId);
            });
          }

          // Add user to selected option
          if (!poll.options[optionIndex].voters.includes(userId)) {
            poll.options[optionIndex].voters.push(userId);
          } else {
            // Remove if already voted for this option
            poll.options[optionIndex].voters = poll.options[optionIndex].voters.filter((voter: string) => voter !== userId);
          }

          // Recalculate votes and total
          poll.options.forEach((opt: any) => {
            opt.votes = opt.voters.length;
          });
          poll.totalVotes = poll.options.reduce((sum: number, opt: any) => sum + opt.votes, 0);

          return {
            ...msg,
            metadata: {
              ...msg.metadata,
              poll
            }
          };
        }
        return msg;
      }) || []
    }));

    toast.success('Vote recorded!');
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

      // Send voice message
      handleSendMessage({
        content: 'Voice message',
        type: 'voice',
        metadata: {
          media: {
            duration: `0:${recordingDuration.toString().padStart(2, '0')}`,
            url: ''
          }
        }
      });
      
      toast.success('Voice message sent!');
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

  const toggleSelectionMode = () => {
    toast.success('Selection mode toggled');
  };

  const renderDialogList = (dialogs: Dialog[]) => (
    <div className="p-1 sm:p-2">
      {dialogs.map((dialog) => (
        <button
          key={dialog.id}
          onClick={() => setSelectedDialog(dialog.id)}
          className={`w-full p-2 sm:p-3 lg:p-4 rounded-lg text-left hover:bg-accent/50 transition-colors mb-1 sm:mb-2 touch-manipulation active:scale-[0.98] ${
            selectedDialog === dialog.id ? 'bg-accent' : ''
          }`}
        >
          <div className="flex items-start gap-2 sm:gap-3 lg:gap-4">
            <div className="relative flex-shrink-0">
              <Avatar className="w-10 h-10 sm:w-12 sm:h-12 lg:w-14 lg:h-14">
                <AvatarImage src={dialog.avatar} />
                <AvatarFallback className="text-sm sm:text-base">
                  {dialog.name.charAt(0)}
                </AvatarFallback>
              </Avatar>
              {dialog.isOnline && (
                <div className="absolute -bottom-0.5 -right-0.5 sm:-bottom-1 sm:-right-1 w-3 h-3 sm:w-4 sm:h-4 lg:w-5 lg:h-5 bg-green-500 rounded-full border-1 sm:border-2 border-background"></div>
              )}
            </div>
            
            <div className="flex-1 min-w-0 overflow-hidden">
              <div className="flex items-center gap-1 sm:gap-2 mb-0.5 sm:mb-1">
                <span className="text-sm sm:text-base lg:text-lg font-medium truncate">{dialog.name}</span>
              </div>
              
              <p className="text-xs sm:text-sm lg:text-base text-muted-foreground truncate mb-1 leading-tight">
                {dialog.lastMessage}
              </p>
              
              {/* Tags */}
              {dialog.tags && dialog.tags.length > 0 && (
                <div className="flex flex-wrap gap-0.5 sm:gap-1 mt-1 sm:mt-1.5 overflow-hidden">
                  {dialog.tags.slice(0, window.innerWidth < 640 ? 1 : 2).map(tag => (
                    <Badge key={tag} variant="outline" className="text-[10px] sm:text-xs px-1 py-0 flex-shrink-0 max-w-20 sm:max-w-none">
                      <Hash className="w-2 h-2 mr-0.5 sm:mr-1 flex-shrink-0" />
                      <span className="truncate">{tag}</span>
                    </Badge>
                  ))}
                  {dialog.tags.length > (window.innerWidth < 640 ? 1 : 2) && (
                    <Badge variant="outline" className="text-[10px] sm:text-xs px-1 py-0 flex-shrink-0">
                      +{dialog.tags.length - (window.innerWidth < 640 ? 1 : 2)}
                    </Badge>
                  )}
                </div>
              )}
            </div>
            
            <div className="flex flex-col items-end gap-0.5 sm:gap-1 flex-shrink-0">
              <span className="text-[10px] sm:text-xs lg:text-sm text-muted-foreground whitespace-nowrap">
                {dialog.timestamp}
              </span>
              {dialog.unread > 0 && (
                <Badge className="w-4 h-4 sm:w-5 sm:h-5 lg:w-6 lg:h-6 text-[10px] sm:text-xs p-0 flex items-center justify-center min-w-4 sm:min-w-5 lg:min-w-6">
                  {dialog.unread > 99 ? '99+' : dialog.unread}
                </Badge>
              )}
            </div>
          </div>
        </button>
      ))}
    </div>
  );

  return (
    <div className="h-full flex flex-col">
      <Tabs defaultValue="chat" className="flex flex-col w-full h-full">
        {/* Mobile-optimized tabs - hidden on mobile when chat is open */}
        <TabsList className={`grid w-full grid-cols-2 md:grid ${selectedDialog ? 'hidden md:grid' : 'grid'}`}>
          <TabsTrigger value="chat" className="flex items-center gap-1.5 sm:gap-2 touch-manipulation text-xs sm:text-sm">
            <MessageCircle className="w-3 h-3 sm:w-4 sm:h-4" />
            <span className="font-medium">Chat</span>
          </TabsTrigger>
          <TabsTrigger value="work" className="flex items-center gap-1.5 sm:gap-2 touch-manipulation text-xs sm:text-sm">
            <Briefcase className="w-3 h-3 sm:w-4 sm:h-4" />
            <span className="font-medium">Work</span>
          </TabsTrigger>
        </TabsList>

        {/* Personal Chat Tab */}
        <TabsContent value="chat" className="flex-1 flex flex-col md:flex-row">
          {/* Chat List - Mobile: Full width when no chat selected, Desktop: Fixed sidebar */}
          <div className={`${selectedDialog ? 'hidden md:flex' : 'flex'} w-full md:w-80 border-r-0 md:border-r border-border flex-col`}>
            {/* Mobile back button when chat is selected */}
            {selectedDialog && (
              <div className="flex md:hidden items-center gap-2 p-3 border-b border-border">
                <Button 
                  variant="ghost" 
                  size="icon" 
                  className="w-8 h-8"
                  onClick={() => setSelectedDialog(null)}
                >
                  <ArrowLeft className="w-4 h-4" />
                </Button>
                <span className="font-medium">Back to Chats</span>
              </div>
            )}
            
            <div className="p-3 md:p-4 border-b border-border">
              <div className="flex items-center gap-2 mb-3">
                <div className="flex-1">
                  <Input
                    placeholder="Search chats..."
                    className="w-full h-9 md:h-10"
                    startIcon={<Search className="w-4 h-4" />}
                  />
                </div>
                <Button size="icon" variant="ghost" onClick={onNewChat} className="w-9 h-9 md:w-10 md:h-10 touch-manipulation">
                  <Plus className="w-4 h-4" />
                </Button>
              </div>
            </div>

            <ScrollArea className="flex-1 mobile-scroll">
              {renderDialogList(personalDialogs)}
            </ScrollArea>
          </div>

          {/* Chat Content - Mobile: Full width when chat selected, Desktop: Flex remainder */}
          <div className={`${selectedDialog ? 'flex' : 'hidden md:flex'} flex-1 flex-col`}>
            {selectedDialog && currentDialog ? (
              <>
                {/* Chat Header - Mobile optimized */}
                <div className="sticky top-0 z-30 border-b border-border p-3 md:p-4 flex items-center justify-between min-h-[60px] md:min-h-auto bg-card/95 backdrop-blur-sm">
                  {/* Mobile back button */}
                  <div className="flex items-center gap-2 md:gap-3">
                    <Button 
                      variant="ghost" 
                      size="icon" 
                      className="w-8 h-8 md:hidden touch-manipulation"
                      onClick={() => setSelectedDialog(null)}
                    >
                      <ArrowLeft className="w-4 h-4" />
                    </Button>
                    <Avatar className="w-8 h-8 md:w-10 md:h-10">
                      <AvatarImage src={currentDialog.avatar} />
                      <AvatarFallback>
                        {currentDialog.name.charAt(0)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="min-w-0">
                      <h3 className="font-medium text-sm md:text-base truncate">{currentDialog.name}</h3>
                      <p className="text-xs md:text-sm text-muted-foreground">
                        {currentDialog.isOnline ? 'Online' : 'Last seen recently'}
                      </p>
                    </div>
                  </div>

                  <ChatHeaderActions
                    selectedDialog={currentDialog}
                    chatActions={chatActions}
                    onVideoCall={() => onVideoCall?.(currentDialog)}
                    onVoiceCall={() => onVoiceCall?.(currentDialog)}
                    onMuteChat={handleMuteChat}
                    onPinChat={handlePinChat}
                    onArchiveChat={handleArchiveChat}
                    onBlockContact={handleBlockContact}
                    onClearHistory={handleClearHistory}
                    toggleSelectionMode={toggleSelectionMode}
                    isSelectionMode={false}
                  />
                </div>

                {/* Messages - Mobile optimized spacing */}
                <ScrollArea className="flex-1 p-3 md:p-4 mobile-scroll">
                  <div className="space-y-4 md:space-y-6 w-full min-w-0 overflow-hidden">
                    {currentMessages.map((message) => (
                      <div key={message.id} className="w-full min-w-0 overflow-hidden">
                        <MarkdownMessage
                          message={message}
                          onReply={handleReplyToMessage}
                          onReact={handleReactToMessage}
                          onShare={handleShareMessage}
                          onCopy={handleCopyMessage}
                          onVote={handleVoteInPoll}
                        />
                      </div>
                    ))}
                  </div>
                </ScrollArea>

                {/* Rich Message Input - Mobile optimized */}
                <div className="border-t border-border">
                  <RichChatInput
                    messageInput={messageInput}
                    setMessageInput={setMessageInput}
                    onSendMessage={handleSendMessage}
                    replyToMessage={replyToMessage}
                    setReplyToMessage={setReplyToMessage}
                    editingMessage={editingMessage}
                    setEditingMessage={setEditingMessage}
                    isRecording={isRecording}
                    recordingDuration={recordingDuration}
                    onStartRecording={startVoiceRecording}
                    onStopRecording={stopVoiceRecording}
                    onCancelRecording={cancelVoiceRecording}
                  />
                </div>
              </>
            ) : (
              <div className="flex-1 flex items-center justify-center p-4">
                <div className="text-center max-w-sm">
                  <MessageCircle className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <p className="text-muted-foreground mb-2">Select a chat to start messaging</p>
                  <p className="text-sm text-muted-foreground">
                    Try rich messages like **markdown**, invoices, polls, and more!
                  </p>
                </div>
              </div>
            )}
          </div>
        </TabsContent>

        {/* Work Tab - Mobile optimized */}
        <TabsContent value="work" className="flex-1 flex flex-col">
          <div className="border-b border-border">
            <WorkspaceSwitcher
              selectedWorkspace={selectedWorkspace || ''}
              userWorkspaces={userWorkspaces}
              onWorkspaceChange={onWorkspaceChange || (() => {})}
              variant="prominent"
              showMetrics={true}
            />
          </div>

          <div className="flex-1 flex flex-col h-full">
            {/* Shop Chat Header */}
            <div className="flex items-center justify-between p-4 border-b border-border">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-gradient-to-br from-chart-2 to-chart-3 rounded-2xl flex items-center justify-center">
                  <Briefcase className="w-5 h-5 text-white" />
                </div>
                <div>
                  <h3 className="font-semibold">Shop Chats</h3>
                  <p className="text-sm text-muted-foreground">Business conversations</p>
                </div>
              </div>
              <Button variant="ghost" size="sm" className="text-chart-2">
                <Plus className="w-4 h-4 mr-1" />
                New
              </Button>
            </div>

            {/* Active Shop Conversations */}
            <ScrollArea className="flex-1">
              <div className="p-4 space-y-4">
                {/* Quick Actions for Shop Chat */}
                <div className="grid grid-cols-2 gap-3 mb-6">
                  <Button variant="outline" className="h-auto p-3 flex flex-col items-center gap-2 group"
                          onClick={() => onShopChatClick?.('golden-mango')}>
                    <div className="w-8 h-8 rounded-full bg-chart-1/10 group-hover:bg-chart-1/20 flex items-center justify-center transition-colors">
                      <MessageCircle className="w-4 h-4 text-chart-1" />
                    </div>
                    <span className="text-xs font-medium">Quick Inquiry</span>
                  </Button>
                  <Button variant="outline" className="h-auto p-3 flex flex-col items-center gap-2 group">
                    <div className="w-8 h-8 rounded-full bg-chart-2/10 group-hover:bg-chart-2/20 flex items-center justify-center transition-colors">
                      <Search className="w-4 h-4 text-chart-2" />
                    </div>
                    <span className="text-xs font-medium">Find Vendor</span>
                  </Button>
                </div>

                {/* Recent Conversations */}
                <div className="space-y-3">
                  <h4 className="text-sm font-medium text-muted-foreground mb-3">Recent Conversations</h4>
                  
                  {/* Golden Mango Restaurant */}
                  <Card className="p-3 cursor-pointer hover:shadow-sm transition-shadow border-l-4 border-l-chart-2"
                        onClick={() => onShopChatClick?.('golden-mango')}>
                    <div className="flex items-start gap-3">
                      <Avatar className="w-10 h-10 flex-shrink-0">
                        <AvatarImage src="https://images.unsplash.com/photo-1743485753872-3b24372fcd24?w=150&h=150&fit=crop" />
                        <AvatarFallback>GM</AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between mb-1">
                          <p className="font-medium text-sm truncate">Golden Mango Restaurant</p>
                          <div className="flex items-center gap-1">
                            <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                            <span className="text-xs text-muted-foreground">Online</span>
                          </div>
                        </div>
                        <p className="text-xs text-muted-foreground mb-2 line-clamp-2">
                          "Your Pad Thai order is ready for pickup! 🍜 Order #PTH-2024-001"
                        </p>
                        <div className="flex items-center justify-between">
                          <Badge variant="secondary" className="text-xs">
                            Order Update
                          </Badge>
                          <span className="text-xs text-muted-foreground">2m ago</span>
                        </div>
                      </div>
                    </div>
                  </Card>

                  {/* Thai Coffee Chain */}
                  <Card className="p-3 cursor-pointer hover:shadow-sm transition-shadow"
                        onClick={() => onShopChatClick?.('thai-coffee')}>
                    <div className="flex items-start gap-3">
                      <Avatar className="w-10 h-10 flex-shrink-0">
                        <AvatarImage src="https://images.unsplash.com/photo-1559847844-5315695b6a77?w=150&h=150&fit=crop" />
                        <AvatarFallback>TC</AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between mb-1">
                          <p className="font-medium text-sm truncate">Thai Coffee Chain</p>
                          <span className="text-xs text-muted-foreground">15m</span>
                        </div>
                        <p className="text-xs text-muted-foreground mb-2 line-clamp-2">
                          "Thanks for the 5-star review! ⭐ Enjoy 10% off your next coffee order"
                        </p>
                        <div className="flex items-center gap-2">
                          <Badge variant="outline" className="text-xs">
                            Review Response
                          </Badge>
                          <Badge className="text-xs bg-green-100 text-green-800 border-green-200">
                            10% Discount
                          </Badge>
                        </div>
                      </div>
                    </div>
                  </Card>

                  {/* Street Food Vendor */}
                  <Card className="p-3 cursor-pointer hover:shadow-sm transition-shadow"
                        onClick={() => onShopChatClick?.('street-food')}>
                    <div className="flex items-start gap-3">
                      <Avatar className="w-10 h-10 flex-shrink-0">
                        <AvatarImage src="https://images.unsplash.com/photo-1590947132387-155cc02f3212?w=150&h=150&fit=crop" />
                        <AvatarFallback>SF</AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between mb-1">
                          <p className="font-medium text-sm truncate">Somtam Street Cart</p>
                          <span className="text-xs text-muted-foreground">1h</span>
                        </div>
                        <p className="text-xs text-muted-foreground mb-2 line-clamp-2">
                          "Fresh papaya salad available! 🥗 Made with organic ingredients"
                        </p>
                        <div className="flex items-center gap-2">
                          <Badge variant="outline" className="text-xs">
                            Daily Special
                          </Badge>
                          <span className="text-xs font-medium text-chart-2">฿45</span>
                        </div>
                      </div>
                    </div>
                  </Card>

                  {/* Electronics Shop */}
                  <Card className="p-3 cursor-pointer hover:shadow-sm transition-shadow">
                    <div className="flex items-start gap-3">
                      <Avatar className="w-10 h-10 flex-shrink-0 bg-chart-4/10">
                        <AvatarFallback className="text-chart-4">ES</AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between mb-1">
                          <p className="font-medium text-sm truncate">Bangkok Electronics</p>
                          <span className="text-xs text-muted-foreground">3h</span>
                        </div>
                        <p className="text-xs text-muted-foreground mb-2 line-clamp-2">
                          "iPhone 15 Pro Max in stock! Payment via PromptPay available 📱"
                        </p>
                        <div className="flex items-center gap-2">
                          <Badge variant="outline" className="text-xs">
                            Product Inquiry
                          </Badge>
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <span>🇹🇭</span>
                            <span>PromptPay</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </Card>
                </div>

                {/* Suggested Shops to Chat */}
                <div className="pt-4 border-t border-border">
                  <h4 className="text-sm font-medium text-muted-foreground mb-3">Discover Local Vendors</h4>
                  <div className="grid grid-cols-2 gap-3">
                    <Button variant="outline" className="h-auto p-3 flex flex-col items-center gap-2 group">
                      <Avatar className="w-8 h-8">
                        <AvatarImage src="https://images.unsplash.com/photo-1565299624946-b28f40a0ca4b?w=150&h=150&fit=crop" />
                        <AvatarFallback>PS</AvatarFallback>
                      </Avatar>
                      <div className="text-center">
                        <p className="text-xs font-medium">Pizza Station</p>
                        <p className="text-xs text-muted-foreground">Italian • 0.5km</p>
                      </div>
                    </Button>
                    <Button variant="outline" className="h-auto p-3 flex flex-col items-center gap-2 group">
                      <Avatar className="w-8 h-8">
                        <AvatarImage src="https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=150&h=150&fit=crop" />
                        <AvatarFallback>SM</AvatarFallback>
                      </Avatar>
                      <div className="text-center">
                        <p className="text-xs font-medium">Siam Market</p>
                        <p className="text-xs text-muted-foreground">Grocery • 0.8km</p>
                      </div>
                    </Button>
                  </div>
                </div>

                {/* Chat Features Info */}
                <div className="mt-6 p-4 bg-gradient-to-r from-chart-1/10 to-chart-2/10 rounded-xl border border-chart-1/20">
                  <div className="flex items-start gap-3">
                    <div className="w-10 h-10 bg-gradient-to-br from-chart-1 to-chart-2 rounded-xl flex items-center justify-center flex-shrink-0">
                      <MessageCircle className="w-5 h-5 text-white" />
                    </div>
                    <div className="flex-1">
                      <h4 className="font-medium text-sm mb-1">SEA Business Chat Features</h4>
                      <ul className="text-xs text-muted-foreground space-y-1">
                        <li>• Real-time order tracking & updates</li>
                        <li>• Multi-language support (TH, ID, VN, MY, PH)</li>
                        <li>• QR payment integration (PromptPay, QRIS, VietQR)</li>
                        <li>• Product catalogs & quick ordering</li>
                        <li>• Location-based vendor discovery</li>
                      </ul>
                    </div>
                  </div>
                </div>

                {/* Bottom Padding for Mobile */}
                <div className="h-20"></div>
              </div>
            </ScrollArea>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}