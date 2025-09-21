import React, { useState } from 'react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader } from './ui/card';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Separator } from './ui/separator';
import { 
  Download, Play, Pause, Heart, Share2, MessageSquareReply, Copy, 
  Receipt, CreditCard, MapPin, Gift, Calendar, Star, Phone, Video,
  FileText, Image as ImageIcon, Music, Film, Archive, Package,
  ShoppingCart, Truck, CheckCircle, Clock, AlertCircle, XCircle,
  Zap, Award, Sparkles, Banknote, TrendingUp, Users, Building2,
  BarChart3, ChefHat, MapPinned, Timer, Utensils, Bike, User
} from 'lucide-react';
import { toast } from "sonner";

interface MessageData {
  id: string;
  senderId: string;
  senderName: string;
  content: string;
  timestamp: string;
  type: 'text' | 'markdown' | 'invoice' | 'bill' | 'order' | 'sticker' | 'poll' | 'contact' | 'location' | 'payment' | 'voice' | 'file' | 'image' | 'video' | 'music' | 'system';
  isOwn: boolean;
  
  // Rich content data
  metadata?: {
    // Order data
    order?: {
      number: string;
      items: Array<{
        name: string;
        quantity: number;
        price: number;
        total: number;
        notes?: string;
        image?: string;
      }>;
      subtotal: number;
      deliveryFee?: number;
      tax?: number;
      discount?: number;
      total: number;
      currency: string;
      status: 'pending' | 'confirmed' | 'preparing' | 'ready' | 'out-for-delivery' | 'delivered' | 'cancelled';
      estimatedTime?: string;
      deliveryTime?: string;
      customer: {
        name: string;
        phone: string;
        address?: string;
        email?: string;
      };
      shop: {
        name: string;
        address: string;
        phone: string;
      };
      deliveryType: 'pickup' | 'delivery';
      paymentStatus: 'pending' | 'paid' | 'failed';
      paymentMethod?: string;
      createdAt: string;
      notes?: string;
    };

    // Invoice/Bill data
    invoice?: {
      number: string;
      items: Array<{
        name: string;
        quantity: number;
        price: number;
        total: number;
      }>;
      subtotal: number;
      tax?: number;
      total: number;
      currency: string;
      dueDate?: string;
      status: 'paid' | 'pending' | 'overdue' | 'cancelled';
      from: {
        name: string;
        address?: string;
        phone?: string;
      };
      to: {
        name: string;
        address?: string;
        phone?: string;
      };
    };
    
    // Sticker data
    sticker?: {
      emoji: string;
      animation?: 'bounce' | 'shake' | 'spin' | 'pulse';
      size: 'small' | 'medium' | 'large';
      pack?: string;
    };
    
    // Contact card
    contact?: {
      name: string;
      phone?: string;
      email?: string;
      avatar?: string;
      company?: string;
      title?: string;
    };
    
    // Location data
    location?: {
      name: string;
      address: string;
      latitude: number;
      longitude: number;
      thumbnail?: string;
    };
    
    // Payment data
    payment?: {
      amount: number;
      currency: string;
      method: string;
      status: 'completed' | 'pending' | 'failed';
      recipient?: string;
      reference?: string;
      fees?: number;
    };
    
    // Poll data
    poll?: {
      question: string;
      options: Array<{
        text: string;
        votes: number;
        voters: string[];
      }>;
      totalVotes: number;
      allowMultiple: boolean;
      anonymous: boolean;
      expiresAt?: string;
    };
    
    // Media data
    media?: {
      url: string;
      thumbnail?: string;
      filename?: string;
      size?: string;
      duration?: string;
      width?: number;
      height?: number;
    };
  };
}

interface MarkdownMessageProps {
  message: MessageData;
  onReply?: (message: MessageData) => void;
  onReact?: (messageId: string, emoji: string) => void;
  onShare?: (message: MessageData) => void;
  onCopy?: (content: string) => void;
  onVote?: (messageId: string, optionIndex: number) => void;
}

export function MarkdownMessage({ 
  message, 
  onReply, 
  onReact, 
  onShare, 
  onCopy,
  onVote 
}: MarkdownMessageProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [selectedPollOptions, setSelectedPollOptions] = useState<number[]>([]);

  // Parse markdown-like content
  const parseMarkdown = (text: string) => {
    return text
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.*?)\*/g, '<em>$1</em>')
      .replace(/`(.*?)`/g, '<code class="bg-muted px-1 rounded text-sm">$1</code>')
      .replace(/~~(.*?)~~/g, '<del>$1</del>')
      .replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2" class="text-primary underline" target="_blank">$1</a>')
      .replace(/\n/g, '<br>');
  };

  const formatCurrency = (amount: number, currency: string) => {
    const symbols = { THB: '฿', USD: '$', EUR: '€', IDR: 'Rp', VND: '₫', MYR: 'RM', PHP: '₱', SGD: 'S$' };
    return `${symbols[currency as keyof typeof symbols] || currency} ${amount.toLocaleString()}`;
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'paid':
      case 'completed':
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'pending':
        return <Clock className="w-4 h-4 text-yellow-500" />;
      case 'overdue':
      case 'failed':
        return <AlertCircle className="w-4 h-4 text-red-500" />;
      case 'cancelled':
        return <XCircle className="w-4 h-4 text-gray-500" />;
      default:
        return null;
    }
  };

  const renderOrderMessage = () => {
    const order = message.metadata?.order;
    if (!order) return null;

    const getStatusColor = (status: string) => {
      switch (status) {
        case 'confirmed':
        case 'delivered':
          return 'text-green-600 bg-green-100';
        case 'preparing':
        case 'ready':
          return 'text-blue-600 bg-blue-100';
        case 'out-for-delivery':
          return 'text-orange-600 bg-orange-100';
        case 'pending':
          return 'text-yellow-600 bg-yellow-100';
        case 'cancelled':
          return 'text-red-600 bg-red-100';
        default:
          return 'text-gray-600 bg-gray-100';
      }
    };

    const getStatusIcon = (status: string) => {
      switch (status) {
        case 'confirmed':
          return <CheckCircle className="w-4 h-4" />;
        case 'preparing':
          return <ChefHat className="w-4 h-4" />;
        case 'ready':
          return <Utensils className="w-4 h-4" />;
        case 'out-for-delivery':
          return <Bike className="w-4 h-4" />;
        case 'delivered':
          return <CheckCircle className="w-4 h-4" />;
        case 'pending':
          return <Clock className="w-4 h-4" />;
        case 'cancelled':
          return <XCircle className="w-4 h-4" />;
        default:
          return <Clock className="w-4 h-4" />;
      }
    };

    return (
      <Card className="max-w-md">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <ShoppingCart className="w-5 h-5 text-chart-2" />
              <span className="font-medium">Order #{order.number}</span>
            </div>
            <div className="flex items-center gap-1">
              <div className={`flex items-center gap-1 px-2 py-1 rounded-full text-xs ${getStatusColor(order.status)}`}>
                {getStatusIcon(order.status)}
                <span className="capitalize">{order.status.replace('-', ' ')}</span>
              </div>
            </div>
          </div>
          
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Calendar className="w-4 h-4" />
            <span>{order.createdAt}</span>
            {order.estimatedTime && (
              <>
                <Timer className="w-4 h-4 ml-2" />
                <span>Est: {order.estimatedTime}</span>
              </>
            )}
          </div>
        </CardHeader>
        
        <CardContent className="space-y-4">
          {/* Customer & Delivery Info */}
          <div className="grid grid-cols-1 gap-3 text-sm">
            <div className="flex items-start gap-2">
              <User className="w-4 h-4 mt-0.5 text-muted-foreground" />
              <div>
                <p className="font-medium">{order.customer.name}</p>
                <p className="text-muted-foreground">{order.customer.phone}</p>
                {order.customer.email && (
                  <p className="text-muted-foreground text-xs">{order.customer.email}</p>
                )}
              </div>
            </div>
            
            <div className="flex items-start gap-2">
              <div className="w-4 h-4 mt-0.5 flex items-center justify-center">
                {order.deliveryType === 'delivery' ? 
                  <Truck className="w-4 h-4 text-muted-foreground" /> : 
                  <MapPinned className="w-4 h-4 text-muted-foreground" />
                }
              </div>
              <div>
                <p className="font-medium capitalize">{order.deliveryType}</p>
                {order.deliveryType === 'delivery' && order.customer.address ? (
                  <p className="text-muted-foreground text-xs">{order.customer.address}</p>
                ) : (
                  <p className="text-muted-foreground text-xs">{order.shop.address}</p>
                )}
                {order.deliveryTime && (
                  <p className="text-muted-foreground text-xs">By {order.deliveryTime}</p>
                )}
              </div>
            </div>
          </div>
          
          <Separator />
          
          {/* Order Items */}
          <div className="space-y-2">
            <h4 className="font-medium text-sm">Order Items</h4>
            {order.items.map((item, index) => (
              <div key={index} className="flex items-center gap-3 text-sm">
                {item.image && (
                  <img 
                    src={item.image} 
                    alt={item.name}
                    className="w-8 h-8 rounded object-cover"
                  />
                )}
                <div className="flex-1">
                  <p className="font-medium">{item.name}</p>
                  {item.notes && (
                    <p className="text-muted-foreground text-xs">Note: {item.notes}</p>
                  )}
                  <p className="text-muted-foreground text-xs">
                    {item.quantity}x {formatCurrency(item.price, order.currency)}
                  </p>
                </div>
                <p className="font-medium">{formatCurrency(item.total, order.currency)}</p>
              </div>
            ))}
          </div>
          
          <Separator />
          
          {/* Order Summary */}
          <div className="space-y-1 text-sm">
            <div className="flex justify-between">
              <span>Subtotal:</span>
              <span>{formatCurrency(order.subtotal, order.currency)}</span>
            </div>
            {order.deliveryFee && order.deliveryFee > 0 && (
              <div className="flex justify-between text-muted-foreground">
                <span>Delivery Fee:</span>
                <span>{formatCurrency(order.deliveryFee, order.currency)}</span>
              </div>
            )}
            {order.tax && order.tax > 0 && (
              <div className="flex justify-between text-muted-foreground">
                <span>Tax:</span>
                <span>{formatCurrency(order.tax, order.currency)}</span>
              </div>
            )}
            {order.discount && order.discount > 0 && (
              <div className="flex justify-between text-green-600">
                <span>Discount:</span>
                <span>-{formatCurrency(order.discount, order.currency)}</span>
              </div>
            )}
            <div className="flex justify-between font-medium text-lg pt-2 border-t">
              <span>Total:</span>
              <span className="text-chart-1">{formatCurrency(order.total, order.currency)}</span>
            </div>
          </div>
          
          {/* Payment Status */}
          <div className="flex items-center justify-between text-sm p-3 bg-muted/50 rounded-lg">
            <div className="flex items-center gap-2">
              <CreditCard className="w-4 h-4" />
              <span>Payment Status:</span>
            </div>
            <div className="flex items-center gap-1">
              {getStatusIcon(order.paymentStatus)}
              <span className={`capitalize ${
                order.paymentStatus === 'paid' ? 'text-green-600' :
                order.paymentStatus === 'failed' ? 'text-red-600' : 'text-yellow-600'
              }`}>
                {order.paymentStatus}
              </span>
            </div>
          </div>
          
          {order.notes && (
            <div className="p-3 bg-blue-50 rounded-lg border-l-4 border-blue-200">
              <p className="text-sm text-blue-800">
                <strong>Notes:</strong> {order.notes}
              </p>
            </div>
          )}
          
          <div className="flex gap-2 pt-2">
            <Button size="sm" className="flex-1">
              <Phone className="w-4 h-4 mr-2" />
              Call Customer
            </Button>
            <Button size="sm" variant="outline" onClick={() => onShare?.(message)}>
              <Share2 className="w-4 h-4" />
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  };

  const renderInvoiceMessage = () => {
    const invoice = message.metadata?.invoice;
    if (!invoice) return null;

    return (
      <Card className="w-full max-w-sm sm:max-w-md lg:max-w-lg xl:max-w-xl">
        <CardHeader className="pb-3 px-3 sm:px-6">
          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-2 min-w-0 flex-1">
              <Receipt className="w-4 h-4 sm:w-5 sm:h-5 text-chart-2 flex-shrink-0" />
              <span className="font-medium text-sm sm:text-base truncate">Invoice #{invoice.number}</span>
            </div>
            <div className="flex items-center gap-1 flex-shrink-0">
              {getStatusIcon(invoice.status)}
              <Badge variant={invoice.status === 'paid' ? 'secondary' : 'destructive'} className="text-xs">
                {invoice.status.toUpperCase()}
              </Badge>
            </div>
          </div>
        </CardHeader>
        
        <CardContent className="space-y-4 px-3 sm:px-6">
          {/* From/To */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4 text-sm">
            <div className="space-y-1">
              <p className="font-medium mb-1 text-xs sm:text-sm">From:</p>
              <p className="text-sm sm:text-base">{invoice.from.name}</p>
              {invoice.from.address && (
                <p className="text-muted-foreground text-xs break-words">{invoice.from.address}</p>
              )}
            </div>
            <div className="space-y-1">
              <p className="font-medium mb-1 text-xs sm:text-sm">To:</p>
              <p className="text-sm sm:text-base">{invoice.to.name}</p>
              {invoice.to.address && (
                <p className="text-muted-foreground text-xs break-words">{invoice.to.address}</p>
              )}
            </div>
          </div>
          
          <Separator />
          
          {/* Items */}
          <div className="space-y-3">
            {invoice.items.map((item, index) => (
              <div key={index} className="flex justify-between items-start gap-2 text-sm">
                <div className="flex-1 min-w-0">
                  <p className="text-sm sm:text-base break-words">{item.name}</p>
                  <p className="text-muted-foreground text-xs">
                    {item.quantity}x {formatCurrency(item.price, invoice.currency)}
                  </p>
                </div>
                <p className="font-medium text-sm sm:text-base flex-shrink-0">
                  {formatCurrency(item.total, invoice.currency)}
                </p>
              </div>
            ))}
          </div>
          
          <Separator />
          
          {/* Totals */}
          <div className="space-y-2 text-sm">
            <div className="flex justify-between items-center">
              <span className="text-sm">Subtotal:</span>
              <span className="text-sm">{formatCurrency(invoice.subtotal, invoice.currency)}</span>
            </div>
            {invoice.tax && (
              <div className="flex justify-between items-center text-muted-foreground">
                <span className="text-sm">Tax:</span>
                <span className="text-sm">{formatCurrency(invoice.tax, invoice.currency)}</span>
              </div>
            )}
            <div className="flex justify-between items-center font-medium text-base sm:text-lg pt-2 border-t">
              <span>Total:</span>
              <span className="text-chart-1">{formatCurrency(invoice.total, invoice.currency)}</span>
            </div>
          </div>
          
          {invoice.dueDate && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Calendar className="w-4 h-4 flex-shrink-0" />
              <span className="text-xs sm:text-sm">Due: {invoice.dueDate}</span>
            </div>
          )}
          
          <div className="flex flex-col sm:flex-row gap-2 pt-2">
            <Button size="sm" className="flex-1 touch-manipulation">
              <Download className="w-4 h-4 mr-2" />
              <span className="text-xs sm:text-sm truncate">Download</span>
            </Button>
            <Button 
              size="sm" 
              variant="outline" 
              className="touch-manipulation sm:w-auto"
              onClick={() => onShare?.(message)}
            >
              <Share2 className="w-4 h-4" />
              <span className="ml-2 sm:hidden text-xs">Share</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  };

  const renderStickerMessage = () => {
    const sticker = message.metadata?.sticker;
    if (!sticker) return null;

    const sizeClasses = {
      small: 'text-4xl',
      medium: 'text-6xl',
      large: 'text-8xl'
    };

    const animationClasses = {
      bounce: 'animate-bounce',
      shake: 'animate-pulse',
      spin: 'animate-spin',
      pulse: 'animate-pulse'
    };

    return (
      <div className="flex items-center justify-center p-4">
        <div className={`
          ${sizeClasses[sticker.size]} 
          ${sticker.animation ? animationClasses[sticker.animation] : ''}
          cursor-pointer hover:scale-110 transition-transform
          select-none
        `}>
          {sticker.emoji}
        </div>
      </div>
    );
  };

  const renderContactMessage = () => {
    const contact = message.metadata?.contact;
    if (!contact) return null;

    return (
      <Card className="max-w-sm">
        <CardContent className="p-4">
          <div className="flex items-center gap-3">
            <Avatar className="w-12 h-12">
              <AvatarImage src={contact.avatar} />
              <AvatarFallback>{contact.name.charAt(0)}</AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="font-medium truncate">{contact.name}</p>
              {contact.title && contact.company && (
                <p className="text-sm text-muted-foreground truncate">
                  {contact.title} at {contact.company}
                </p>
              )}
              {contact.phone && (
                <p className="text-sm text-muted-foreground">{contact.phone}</p>
              )}
            </div>
          </div>
          <div className="flex gap-2 mt-3">
            {contact.phone && (
              <Button size="sm" variant="outline" className="flex-1">
                <Phone className="w-4 h-4 mr-2" />
                Call
              </Button>
            )}
            <Button size="sm" variant="outline" className="flex-1">
              <MessageSquareReply className="w-4 h-4 mr-2" />
              Message
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  };

  const renderLocationMessage = () => {
    const location = message.metadata?.location;
    if (!location) return null;

    return (
      <Card className="max-w-sm">
        <CardContent className="p-0">
          <div className="aspect-video bg-muted rounded-t-lg flex items-center justify-center relative overflow-hidden">
            {location.thumbnail ? (
              <img src={location.thumbnail} alt="Location" className="w-full h-full object-cover" />
            ) : (
              <MapPin className="w-12 h-12 text-muted-foreground" />
            )}
            <div className="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent" />
            <div className="absolute bottom-2 left-2 text-white">
              <div className="flex items-center gap-1">
                <MapPin className="w-4 h-4" />
                <span className="text-sm font-medium">{location.name}</span>
              </div>
            </div>
          </div>
          <div className="p-4">
            <p className="text-sm text-muted-foreground">{location.address}</p>
            <div className="flex gap-2 mt-3">
              <Button size="sm" variant="outline" className="flex-1">
                <MapPin className="w-4 h-4 mr-2" />
                Directions
              </Button>
              <Button size="sm" variant="outline" onClick={() => onShare?.(message)}>
                <Share2 className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  };

  const renderPaymentMessage = () => {
    const payment = message.metadata?.payment;
    if (!payment) return null;

    return (
      <Card className="max-w-sm">
        <CardContent className="p-4">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
                <CreditCard className="w-5 h-5 text-green-600" />
              </div>
              <div>
                <p className="font-medium">Payment {payment.status === 'completed' ? 'Sent' : 'Pending'}</p>
                <p className="text-sm text-muted-foreground">{payment.method}</p>
              </div>
            </div>
            {getStatusIcon(payment.status)}
          </div>
          
          <div className="text-center py-4">
            <p className="text-3xl font-bold text-green-600">
              {formatCurrency(payment.amount, payment.currency)}
            </p>
            {payment.recipient && (
              <p className="text-sm text-muted-foreground mt-1">to {payment.recipient}</p>
            )}
          </div>
          
          {payment.reference && (
            <div className="flex items-center justify-between text-sm border-t pt-3">
              <span className="text-muted-foreground">Reference:</span>
              <span className="font-mono">{payment.reference}</span>
            </div>
          )}
          
          {payment.fees && payment.fees > 0 && (
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">Fees:</span>
              <span>{formatCurrency(payment.fees, payment.currency)}</span>
            </div>
          )}
        </CardContent>
      </Card>
    );
  };

  const renderPollMessage = () => {
    const poll = message.metadata?.poll;
    if (!poll) return null;

    const handleVote = (optionIndex: number) => {
      if (poll.allowMultiple) {
        setSelectedPollOptions(prev => 
          prev.includes(optionIndex) 
            ? prev.filter(i => i !== optionIndex)
            : [...prev, optionIndex]
        );
      } else {
        setSelectedPollOptions([optionIndex]);
      }
      onVote?.(message.id, optionIndex);
    };

    return (
      <Card className="max-w-md">
        <CardHeader className="pb-3">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-chart-1/10 rounded-full flex items-center justify-center">
              <BarChart3 className="w-4 h-4 text-chart-1" />
            </div>
            <span className="font-medium">{poll.question}</span>
          </div>
        </CardHeader>
        
        <CardContent className="space-y-2">
          {poll.options.map((option, index) => {
            const percentage = poll.totalVotes > 0 ? (option.votes / poll.totalVotes) * 100 : 0;
            const isSelected = selectedPollOptions.includes(index);
            
            return (
              <button
                key={index}
                onClick={() => handleVote(index)}
                className={`w-full p-3 rounded-lg border text-left relative overflow-hidden transition-colors ${
                  isSelected ? 'border-primary bg-primary/5' : 'border-border hover:bg-accent'
                }`}
              >
                <div 
                  className="absolute inset-0 bg-primary/10 transition-all duration-300"
                  style={{ width: `${percentage}%` }}
                />
                <div className="relative flex items-center justify-between">
                  <span className="font-medium">{option.text}</span>
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-muted-foreground">{option.votes}</span>
                    <span className="text-sm text-muted-foreground">{percentage.toFixed(0)}%</span>
                  </div>
                </div>
              </button>
            );
          })}
          
          <div className="flex items-center justify-between text-sm text-muted-foreground pt-2">
            <span>{poll.totalVotes} vote{poll.totalVotes !== 1 ? 's' : ''}</span>
            {poll.expiresAt && (
              <span>Expires {poll.expiresAt}</span>
            )}
          </div>
        </CardContent>
      </Card>
    );
  };

  const renderVoiceMessage = () => {
    const media = message.metadata?.media;
    
    return (
      <div className="flex items-center gap-3 bg-muted rounded-lg p-3 max-w-xs">
        <Button 
          size="sm" 
          variant="ghost" 
          className="rounded-full w-10 h-10 p-0"
          onClick={() => setIsPlaying(!isPlaying)}
        >
          {isPlaying ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
        </Button>
        
        <div className="flex-1">
          <div className="flex items-center gap-1 mb-1">
            {[...Array(20)].map((_, i) => (
              <div 
                key={i} 
                className={`w-0.5 bg-primary/60 rounded-full transition-all ${
                  isPlaying && i < 8 ? 'h-6 animate-pulse' : 'h-2'
                }`} 
              />
            ))}
          </div>
          <p className="text-xs text-muted-foreground">{media?.duration || '0:12'}</p>
        </div>
      </div>
    );
  };

  const renderMediaMessage = () => {
    const media = message.metadata?.media;
    if (!media) return null;

    if (message.type === 'image') {
      return (
        <div className="max-w-sm">
          <img 
            src={media.url} 
            alt="Shared image" 
            className="w-full rounded-lg max-h-64 object-cover"
          />
          {message.content && (
            <p className="text-sm mt-2 text-muted-foreground">{message.content}</p>
          )}
        </div>
      );
    }

    if (message.type === 'video') {
      return (
        <div className="max-w-sm">
          <div className="relative aspect-video bg-black rounded-lg overflow-hidden">
            <video 
              src={media.url}
              poster={media.thumbnail}
              className="w-full h-full object-cover"
              controls
            />
          </div>
          {message.content && (
            <p className="text-sm mt-2 text-muted-foreground">{message.content}</p>
          )}
        </div>
      );
    }

    return null;
  };

  const renderMessageContent = () => {
    switch (message.type) {
      case 'order':
        return renderOrderMessage();
        
      case 'invoice':
      case 'bill':
        return renderInvoiceMessage();
      
      case 'sticker':
        return renderStickerMessage();
      
      case 'contact':
        return renderContactMessage();
      
      case 'location':
        return renderLocationMessage();
      
      case 'payment':
        return renderPaymentMessage();
      
      case 'poll':
        return renderPollMessage();
      
      case 'voice':
        return renderVoiceMessage();
      
      case 'image':
      case 'video':
        return renderMediaMessage();
      
      case 'markdown':
        return (
          <div 
            className="prose prose-sm max-w-none"
            dangerouslySetInnerHTML={{ __html: parseMarkdown(message.content) }}
          />
        );
      
      default:
        return <p className="text-sm whitespace-pre-line">{message.content}</p>;
    }
  };

  return (
    <div className={`flex gap-3 ${message.isOwn ? 'flex-row-reverse' : ''} group`}>
      {!message.isOwn && (
        <Avatar className="w-8 h-8 flex-shrink-0">
          <AvatarFallback className="text-xs">
            {message.type === 'system' ? <Sparkles className="w-4 h-4" /> : message.senderName.charAt(0)}
          </AvatarFallback>
        </Avatar>
      )}
      
      <div className={`flex flex-col max-w-md ${message.isOwn ? 'items-end' : 'items-start'}`}>
        {!message.isOwn && message.type !== 'sticker' && (
          <span className="text-xs text-muted-foreground mb-1 flex items-center gap-1">
            {message.type === 'system' && <Sparkles className="w-3 h-3" />}
            {message.senderName}
          </span>
        )}
        
        <div className={`${
          message.type === 'sticker' ? '' : 
          message.type === 'system' ? 'bg-chart-1/10 border border-chart-1/20 rounded-lg sm:rounded-xl p-2 sm:p-3 lg:p-4 touch-manipulation' :
          ['invoice', 'bill', 'contact', 'location', 'payment', 'poll'].includes(message.type) ? '' :
          message.isOwn ? 'bg-primary text-primary-foreground rounded-lg sm:rounded-xl px-2 py-1.5 sm:px-3 sm:py-2 lg:px-4 lg:py-3 touch-manipulation max-w-[85%] sm:max-w-[75%] lg:max-w-[65%] text-sm sm:text-base' : 'bg-muted rounded-lg sm:rounded-xl px-2 py-1.5 sm:px-3 sm:py-2 lg:px-4 lg:py-3 touch-manipulation max-w-[85%] sm:max-w-[75%] lg:max-w-[65%] text-sm sm:text-base'
        }`}>
          {renderMessageContent()}
        </div>
        
        <div className="flex items-center gap-2 mt-1">
          <span className="text-xs text-muted-foreground">
            {message.timestamp}
          </span>
          
          {/* Quick reaction buttons */}
          <div className="opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
            <Button 
              size="sm" 
              variant="ghost" 
              className="h-6 w-6 p-0 hover:bg-background/20"
              onClick={() => onReact?.(message.id, '👍')}
            >
              👍
            </Button>
            <Button 
              size="sm" 
              variant="ghost" 
              className="h-6 w-6 p-0 hover:bg-background/20"
              onClick={() => onReact?.(message.id, '❤️')}
            >
              ❤️
            </Button>
            <Button 
              size="sm" 
              variant="ghost" 
              className="h-6 w-6 p-0 hover:bg-background/20"
              onClick={() => onReply?.(message)}
            >
              <MessageSquareReply className="w-3 h-3" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}