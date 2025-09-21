import React, { useRef, useState } from 'react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader } from './ui/card';
import { Textarea } from './ui/textarea';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from './ui/dropdown-menu';
import { 
  Send, Paperclip, Mic, Smile, File, Image, Camera, MapPin, CreditCard, Gift,
  X, Pause, Play, Receipt, Users, BarChart3, Contact, MessageSquareText,
  Bold, Italic, Code, Link2, List, Hash, AtSign, Calendar, Star,
  Zap, Sparkles, DollarSign, Package, ShoppingCart, Clock, Truck, ChefHat
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
  metadata?: any;
}

interface RichChatInputProps {
  messageInput: string;
  setMessageInput: (value: string) => void;
  onSendMessage: (messageData: Partial<MessageData>) => void;
  replyToMessage: MessageData | null;
  setReplyToMessage: (message: MessageData | null) => void;
  editingMessage: MessageData | null;
  setEditingMessage: (message: MessageData | null) => void;
  isRecording: boolean;
  recordingDuration: number;
  onStartRecording: () => void;
  onStopRecording: () => void;
  onCancelRecording: () => void;
}

export function RichChatInput({
  messageInput,
  setMessageInput,
  onSendMessage,
  replyToMessage,
  setReplyToMessage,
  editingMessage,
  setEditingMessage,
  isRecording,
  recordingDuration,
  onStartRecording,
  onStopRecording,
  onCancelRecording
}: RichChatInputProps) {
  const [showAttachmentMenu, setShowAttachmentMenu] = useState(false);
  const [showRichContentDialog, setShowRichContentDialog] = useState<string | null>(null);
  const [orderForm, setOrderForm] = useState({
    number: '',
    items: [{ name: '', quantity: 1, price: 0, notes: '', image: '' }],
    customer: { name: '', phone: '', address: '', email: '' },
    shop: { name: 'Golden Mango Restaurant', address: '123 Sukhumvit Road, Bangkok', phone: '+66 2 123 4567' },
    deliveryType: 'delivery' as 'pickup' | 'delivery',
    deliveryFee: 30,
    currency: 'THB',
    estimatedTime: '30-45 min',
    deliveryTime: '',
    notes: '',
    paymentMethod: 'PromptPay'
  });
  const [invoiceForm, setInvoiceForm] = useState({
    number: '',
    items: [{ name: '', quantity: 1, price: 0 }],
    from: { name: '', address: '', phone: '' },
    to: { name: '', address: '', phone: '' },
    currency: 'THB'
  });
  const [pollForm, setPollForm] = useState({
    question: '',
    options: ['', ''],
    allowMultiple: false,
    anonymous: true
  });
  const [contactForm, setContactForm] = useState({
    name: '',
    phone: '',
    email: '',
    company: '',
    title: ''
  });
  const [paymentForm, setPaymentForm] = useState({
    amount: 0,
    currency: 'THB',
    method: 'PromptPay',
    recipient: ''
  });

  const fileInputRef = useRef<HTMLInputElement>(null);
  const imageInputRef = useRef<HTMLInputElement>(null);

  const formatRecordingTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const sendOrder = () => {
    const items = orderForm.items.filter(item => item.name.trim());
    if (items.length === 0) {
      toast.error('Please add at least one item');
      return;
    }
    
    if (!orderForm.customer.name.trim() || !orderForm.customer.phone.trim()) {
      toast.error('Customer name and phone are required');
      return;
    }

    const subtotal = items.reduce((sum, item) => sum + (item.quantity * item.price), 0);
    const tax = subtotal * 0.07; // 7% tax
    const deliveryFee = orderForm.deliveryType === 'delivery' ? orderForm.deliveryFee : 0;
    const total = subtotal + tax + deliveryFee;

    // Generate order number if not provided
    const orderNumber = orderForm.number || `ORD-${Date.now()}`;

    onSendMessage({
      content: `New order created: ${orderNumber}`,
      type: 'order',
      metadata: {
        order: {
          number: orderNumber,
          items: items.map(item => ({
            ...item,
            total: item.quantity * item.price
          })),
          subtotal,
          deliveryFee,
          tax,
          total,
          currency: orderForm.currency,
          status: 'pending',
          estimatedTime: orderForm.estimatedTime,
          deliveryTime: orderForm.deliveryTime,
          customer: orderForm.customer,
          shop: orderForm.shop,
          deliveryType: orderForm.deliveryType,
          paymentStatus: 'pending',
          paymentMethod: orderForm.paymentMethod,
          createdAt: new Date().toLocaleString(),
          notes: orderForm.notes
        }
      }
    });

    // Reset form
    setOrderForm({
      number: '',
      items: [{ name: '', quantity: 1, price: 0, notes: '', image: '' }],
      customer: { name: '', phone: '', address: '', email: '' },
      shop: { name: 'Golden Mango Restaurant', address: '123 Sukhumvit Road, Bangkok', phone: '+66 2 123 4567' },
      deliveryType: 'delivery',
      deliveryFee: 30,
      currency: 'THB',
      estimatedTime: '30-45 min',
      deliveryTime: '',
      notes: '',
      paymentMethod: 'PromptPay'
    });
    setShowRichContentDialog(null);
  };

  const insertMarkdownFormat = (format: string) => {
    const textarea = document.querySelector('input[type="text"]') as HTMLInputElement;
    if (!textarea) return;

    const start = textarea.selectionStart || 0;
    const end = textarea.selectionEnd || 0;
    const selectedText = messageInput.substring(start, end);
    
    let formatted = '';
    switch (format) {
      case 'bold':
        formatted = selectedText ? `**${selectedText}**` : '****';
        break;
      case 'italic':
        formatted = selectedText ? `*${selectedText}*` : '**';
        break;
      case 'code':
        formatted = selectedText ? `\`${selectedText}\`` : '``';
        break;
      case 'link':
        formatted = selectedText ? `[${selectedText}](url)` : '[text](url)';
        break;
    }

    const newText = messageInput.substring(0, start) + formatted + messageInput.substring(end);
    setMessageInput(newText);
  };

  const sendInvoice = () => {
    const items = invoiceForm.items.filter(item => item.name.trim());
    const subtotal = items.reduce((sum, item) => sum + (item.quantity * item.price), 0);
    const tax = subtotal * 0.07; // 7% tax
    const total = subtotal + tax;

    onSendMessage({
      content: `Invoice #${invoiceForm.number}`,
      type: 'invoice',
      metadata: {
        invoice: {
          number: invoiceForm.number,
          items: items.map(item => ({
            ...item,
            total: item.quantity * item.price
          })),
          subtotal,
          tax,
          total,
          currency: invoiceForm.currency,
          status: 'pending',
          from: invoiceForm.from,
          to: invoiceForm.to,
          dueDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toLocaleDateString()
        }
      }
    });

    // Reset form
    setInvoiceForm({
      number: '',
      items: [{ name: '', quantity: 1, price: 0 }],
      from: { name: '', address: '', phone: '' },
      to: { name: '', address: '', phone: '' },
      currency: 'THB'
    });
    setShowRichContentDialog(null);
  };

  const sendPoll = () => {
    const validOptions = pollForm.options.filter(opt => opt.trim());
    if (validOptions.length < 2) {
      toast.error('Poll needs at least 2 options');
      return;
    }

    onSendMessage({
      content: pollForm.question,
      type: 'poll',
      metadata: {
        poll: {
          question: pollForm.question,
          options: validOptions.map(text => ({
            text,
            votes: 0,
            voters: []
          })),
          totalVotes: 0,
          allowMultiple: pollForm.allowMultiple,
          anonymous: pollForm.anonymous
        }
      }
    });

    // Reset form
    setPollForm({
      question: '',
      options: ['', ''],
      allowMultiple: false,
      anonymous: true
    });
    setShowRichContentDialog(null);
  };

  const sendContact = () => {
    if (!contactForm.name.trim()) {
      toast.error('Contact name is required');
      return;
    }

    onSendMessage({
      content: `Contact: ${contactForm.name}`,
      type: 'contact',
      metadata: {
        contact: contactForm
      }
    });

    // Reset form
    setContactForm({
      name: '',
      phone: '',
      email: '',
      company: '',
      title: ''
    });
    setShowRichContentDialog(null);
  };

  const sendPayment = () => {
    if (paymentForm.amount <= 0) {
      toast.error('Payment amount must be greater than 0');
      return;
    }

    const symbols = { THB: '฿', USD: '$', EUR: '€', IDR: 'Rp', VND: '₫', MYR: 'RM', PHP: '₱', SGD: 'S$' };
    const symbol = symbols[paymentForm.currency as keyof typeof symbols] || paymentForm.currency;

    onSendMessage({
      content: `Payment of ${symbol}${paymentForm.amount.toLocaleString()}`,
      type: 'payment',
      metadata: {
        payment: {
          ...paymentForm,
          status: 'pending',
          reference: `REF${Date.now()}`
        }
      }
    });

    // Reset form
    setPaymentForm({
      amount: 0,
      currency: 'THB',
      method: 'PromptPay',
      recipient: ''
    });
    setShowRichContentDialog(null);
  };

  const sendSticker = (emoji: string, size: 'small' | 'medium' | 'large' = 'medium') => {
    onSendMessage({
      content: emoji,
      type: 'sticker',
      metadata: {
        sticker: {
          emoji,
          size,
          animation: 'bounce'
        }
      }
    });
  };

  const handleSendMessage = () => {
    if (!messageInput.trim()) return;

    // Check if message contains markdown
    const hasMarkdown = /[*_`\[\]()~]/.test(messageInput);
    
    onSendMessage({
      content: messageInput,
      type: hasMarkdown ? 'markdown' : 'text'
    });
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const cancelReply = () => {
    setReplyToMessage(null);
  };

  const cancelEdit = () => {
    setEditingMessage(null);
    setMessageInput('');
  };

  const addOrderItem = () => {
    setOrderForm(prev => ({
      ...prev,
      items: [...prev.items, { name: '', quantity: 1, price: 0, notes: '', image: '' }]
    }));
  };

  const removeOrderItem = (index: number) => {
    setOrderForm(prev => ({
      ...prev,
      items: prev.items.filter((_, i) => i !== index)
    }));
  };

  const addInvoiceItem = () => {
    setInvoiceForm(prev => ({
      ...prev,
      items: [...prev.items, { name: '', quantity: 1, price: 0 }]
    }));
  };

  const removeInvoiceItem = (index: number) => {
    setInvoiceForm(prev => ({
      ...prev,
      items: prev.items.filter((_, i) => i !== index)
    }));
  };

  const addPollOption = () => {
    setPollForm(prev => ({
      ...prev,
      options: [...prev.options, '']
    }));
  };

  const removePollOption = (index: number) => {
    setPollForm(prev => ({
      ...prev,
      options: prev.options.filter((_, i) => i !== index)
    }));
  };

  const renderOrderDialog = () => (
    <Dialog open={showRichContentDialog === 'order'} onOpenChange={() => setShowRichContentDialog(null)}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <ShoppingCart className="w-5 h-5" />
            Create Customer Order
          </DialogTitle>
          <DialogDescription>
            Create a new order for a customer at your shop
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Order Details */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium">Order Number (Optional)</label>
              <Input
                value={orderForm.number}
                onChange={(e) => setOrderForm(prev => ({ ...prev, number: e.target.value }))}
                placeholder="Auto-generated if empty"
              />
            </div>
            <div>
              <label className="text-sm font-medium">Currency</label>
              <select 
                className="w-full p-2 border rounded-md"
                value={orderForm.currency}
                onChange={(e) => setOrderForm(prev => ({ ...prev, currency: e.target.value }))}
              >
                <option value="THB">Thai Baht (฿)</option>
                <option value="USD">US Dollar ($)</option>
                <option value="IDR">Indonesian Rupiah (Rp)</option>
                <option value="VND">Vietnamese Dong (₫)</option>
                <option value="MYR">Malaysian Ringgit (RM)</option>
                <option value="PHP">Philippine Peso (₱)</option>
                <option value="SGD">Singapore Dollar (S$)</option>
              </select>
            </div>
          </div>

          {/* Customer Information */}
          <div>
            <h4 className="font-medium mb-3">Customer Information</h4>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="text-sm font-medium">Customer Name *</label>
                <Input
                  value={orderForm.customer.name}
                  onChange={(e) => setOrderForm(prev => ({ 
                    ...prev, 
                    customer: { ...prev.customer, name: e.target.value }
                  }))}
                  placeholder="John Smith"
                />
              </div>
              <div>
                <label className="text-sm font-medium">Phone Number *</label>
                <Input
                  value={orderForm.customer.phone}
                  onChange={(e) => setOrderForm(prev => ({ 
                    ...prev, 
                    customer: { ...prev.customer, phone: e.target.value }
                  }))}
                  placeholder="+66 81 234 5678"
                />
              </div>
              <div className="col-span-2">
                <label className="text-sm font-medium">Email (Optional)</label>
                <Input
                  value={orderForm.customer.email}
                  onChange={(e) => setOrderForm(prev => ({ 
                    ...prev, 
                    customer: { ...prev.customer, email: e.target.value }
                  }))}
                  placeholder="john@example.com"
                />
              </div>
            </div>
          </div>

          {/* Delivery Information */}
          <div>
            <h4 className="font-medium mb-3">Delivery Information</h4>
            <div className="space-y-4">
              <div className="flex gap-4">
                <label className="flex items-center gap-2">
                  <input
                    type="radio"
                    name="deliveryType"
                    checked={orderForm.deliveryType === 'pickup'}
                    onChange={() => setOrderForm(prev => ({ ...prev, deliveryType: 'pickup' }))}
                  />
                  <span className="text-sm">Pickup at Store</span>
                </label>
                <label className="flex items-center gap-2">
                  <input
                    type="radio"
                    name="deliveryType"
                    checked={orderForm.deliveryType === 'delivery'}
                    onChange={() => setOrderForm(prev => ({ ...prev, deliveryType: 'delivery' }))}
                  />
                  <span className="text-sm">Home Delivery</span>
                </label>
              </div>
              
              {orderForm.deliveryType === 'delivery' && (
                <div className="grid grid-cols-2 gap-4">
                  <div className="col-span-2">
                    <label className="text-sm font-medium">Delivery Address</label>
                    <Textarea
                      value={orderForm.customer.address}
                      onChange={(e) => setOrderForm(prev => ({ 
                        ...prev, 
                        customer: { ...prev.customer, address: e.target.value }
                      }))}
                      placeholder="123 Main Street, Bangkok 10110"
                      rows={2}
                    />
                  </div>
                  <div>
                    <label className="text-sm font-medium">Delivery Fee</label>
                    <Input
                      type="number"
                      value={orderForm.deliveryFee}
                      onChange={(e) => setOrderForm(prev => ({ ...prev, deliveryFee: parseFloat(e.target.value) || 0 }))}
                      placeholder="30"
                    />
                  </div>
                  <div>
                    <label className="text-sm font-medium">Delivery Time</label>
                    <Input
                      value={orderForm.deliveryTime}
                      onChange={(e) => setOrderForm(prev => ({ ...prev, deliveryTime: e.target.value }))}
                      placeholder="e.g., 2:00 PM"
                    />
                  </div>
                </div>
              )}
              
              <div>
                <label className="text-sm font-medium">Estimated Preparation Time</label>
                <select 
                  className="w-full p-2 border rounded-md"
                  value={orderForm.estimatedTime}
                  onChange={(e) => setOrderForm(prev => ({ ...prev, estimatedTime: e.target.value }))}
                >
                  <option value="15-20 min">15-20 minutes</option>
                  <option value="20-30 min">20-30 minutes</option>
                  <option value="30-45 min">30-45 minutes</option>
                  <option value="45-60 min">45-60 minutes</option>
                  <option value="1-1.5 hours">1-1.5 hours</option>
                </select>
              </div>
            </div>
          </div>

          {/* Order Items */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <h4 className="font-medium">Order Items</h4>
              <Button size="sm" onClick={addOrderItem}>
                <Package className="w-4 h-4 mr-2" />
                Add Item
              </Button>
            </div>
            <div className="space-y-3">
              {orderForm.items.map((item, index) => (
                <div key={index} className="p-3 border rounded-lg">
                  <div className="grid grid-cols-3 gap-2 mb-2">
                    <div className="col-span-2">
                      <Input
                        placeholder="Item name (e.g., Pad Thai)"
                        value={item.name}
                        onChange={(e) => {
                          const newItems = [...orderForm.items];
                          newItems[index].name = e.target.value;
                          setOrderForm(prev => ({ ...prev, items: newItems }));
                        }}
                      />
                    </div>
                    <div className="flex gap-1">
                      <Input
                        type="number"
                        placeholder="Qty"
                        value={item.quantity}
                        onChange={(e) => {
                          const newItems = [...orderForm.items];
                          newItems[index].quantity = parseInt(e.target.value) || 0;
                          setOrderForm(prev => ({ ...prev, items: newItems }));
                        }}
                      />
                      <Input
                        type="number"
                        placeholder="Price"
                        step="0.01"
                        value={item.price}
                        onChange={(e) => {
                          const newItems = [...orderForm.items];
                          newItems[index].price = parseFloat(e.target.value) || 0;
                          setOrderForm(prev => ({ ...prev, items: newItems }));
                        }}
                      />
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Input
                      placeholder="Special notes (e.g., extra spicy, no onions)"
                      value={item.notes}
                      onChange={(e) => {
                        const newItems = [...orderForm.items];
                        newItems[index].notes = e.target.value;
                        setOrderForm(prev => ({ ...prev, items: newItems }));
                      }}
                      className="flex-1"
                    />
                    {orderForm.items.length > 1 && (
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => removeOrderItem(index)}
                      >
                        <X className="w-4 h-4" />
                      </Button>
                    )}
                  </div>
                  <div className="text-right text-sm text-muted-foreground mt-1">
                    Total: {((item.quantity || 0) * (item.price || 0)).toFixed(2)}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Payment Method */}
          <div>
            <label className="text-sm font-medium">Payment Method</label>
            <select 
              className="w-full p-2 border rounded-md"
              value={orderForm.paymentMethod}
              onChange={(e) => setOrderForm(prev => ({ ...prev, paymentMethod: e.target.value }))}
            >
              <option value="PromptPay">PromptPay (Thailand)</option>
              <option value="Cash">Cash on Delivery/Pickup</option>
              <option value="Card">Credit/Debit Card</option>
              <option value="QRIS">QRIS (Indonesia)</option>
              <option value="VietQR">VietQR (Vietnam)</option>
              <option value="DuitNow">DuitNow QR (Malaysia)</option>
              <option value="InstaPay">InstaPay QR (Philippines)</option>
            </select>
          </div>

          {/* Order Notes */}
          <div>
            <label className="text-sm font-medium">Order Notes (Optional)</label>
            <Textarea
              value={orderForm.notes}
              onChange={(e) => setOrderForm(prev => ({ ...prev, notes: e.target.value }))}
              placeholder="Any special instructions or notes for this order..."
              rows={2}
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => setShowRichContentDialog(null)}>
            Cancel
          </Button>
          <Button onClick={sendOrder}>
            <ShoppingCart className="w-4 h-4 mr-2" />
            Create Order
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );

  const renderInvoiceDialog = () => (
    <Dialog open={showRichContentDialog === 'invoice'} onOpenChange={() => setShowRichContentDialog(null)}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Receipt className="w-5 h-5" />
            Create Invoice
          </DialogTitle>
          <DialogDescription>
            Generate a professional invoice to send in your chat
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Invoice Details */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium">Invoice Number</label>
              <Input
                value={invoiceForm.number}
                onChange={(e) => setInvoiceForm(prev => ({ ...prev, number: e.target.value }))}
                placeholder="INV-001"
              />
            </div>
            <div>
              <label className="text-sm font-medium">Currency</label>
              <select 
                className="w-full p-2 border rounded-md"
                value={invoiceForm.currency}
                onChange={(e) => setInvoiceForm(prev => ({ ...prev, currency: e.target.value }))}
              >
                <option value="THB">Thai Baht (฿)</option>
                <option value="USD">US Dollar ($)</option>
                <option value="EUR">Euro (€)</option>
                <option value="IDR">Indonesian Rupiah (Rp)</option>
                <option value="VND">Vietnamese Dong (₫)</option>
                <option value="MYR">Malaysian Ringgit (RM)</option>
                <option value="PHP">Philippine Peso (₱)</option>
                <option value="SGD">Singapore Dollar (S$)</option>
              </select>
            </div>
          </div>

          {/* From/To */}
          <div className="grid grid-cols-2 gap-6">
            <div>
              <h4 className="font-medium mb-2">From:</h4>
              <div className="space-y-2">
                <Input
                  placeholder="Company/Name"
                  value={invoiceForm.from.name}
                  onChange={(e) => setInvoiceForm(prev => ({ 
                    ...prev, 
                    from: { ...prev.from, name: e.target.value }
                  }))}
                />
                <Input
                  placeholder="Address"
                  value={invoiceForm.from.address}
                  onChange={(e) => setInvoiceForm(prev => ({ 
                    ...prev, 
                    from: { ...prev.from, address: e.target.value }
                  }))}
                />
                <Input
                  placeholder="Phone"
                  value={invoiceForm.from.phone}
                  onChange={(e) => setInvoiceForm(prev => ({ 
                    ...prev, 
                    from: { ...prev.from, phone: e.target.value }
                  }))}
                />
              </div>
            </div>
            <div>
              <h4 className="font-medium mb-2">To:</h4>
              <div className="space-y-2">
                <Input
                  placeholder="Company/Name"
                  value={invoiceForm.to.name}
                  onChange={(e) => setInvoiceForm(prev => ({ 
                    ...prev, 
                    to: { ...prev.to, name: e.target.value }
                  }))}
                />
                <Input
                  placeholder="Address"
                  value={invoiceForm.to.address}
                  onChange={(e) => setInvoiceForm(prev => ({ 
                    ...prev, 
                    to: { ...prev.to, address: e.target.value }
                  }))}
                />
                <Input
                  placeholder="Phone"
                  value={invoiceForm.to.phone}
                  onChange={(e) => setInvoiceForm(prev => ({ 
                    ...prev, 
                    to: { ...prev.to, phone: e.target.value }
                  }))}
                />
              </div>
            </div>
          </div>

          {/* Items */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <h4 className="font-medium">Items</h4>
              <Button size="sm" onClick={addInvoiceItem}>
                <Package className="w-4 h-4 mr-2" />
                Add Item
              </Button>
            </div>
            <div className="space-y-2">
              {invoiceForm.items.map((item, index) => (
                <div key={index} className="flex gap-2 items-end">
                  <div className="flex-1">
                    <Input
                      placeholder="Item name"
                      value={item.name}
                      onChange={(e) => {
                        const newItems = [...invoiceForm.items];
                        newItems[index].name = e.target.value;
                        setInvoiceForm(prev => ({ ...prev, items: newItems }));
                      }}
                    />
                  </div>
                  <div className="w-20">
                    <Input
                      type="number"
                      placeholder="Qty"
                      value={item.quantity}
                      onChange={(e) => {
                        const newItems = [...invoiceForm.items];
                        newItems[index].quantity = parseInt(e.target.value) || 0;
                        setInvoiceForm(prev => ({ ...prev, items: newItems }));
                      }}
                    />
                  </div>
                  <div className="w-24">
                    <Input
                      type="number"
                      placeholder="Price"
                      step="0.01"
                      value={item.price}
                      onChange={(e) => {
                        const newItems = [...invoiceForm.items];
                        newItems[index].price = parseFloat(e.target.value) || 0;
                        setInvoiceForm(prev => ({ ...prev, items: newItems }));
                      }}
                    />
                  </div>
                  {invoiceForm.items.length > 1 && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => removeInvoiceItem(index)}
                    >
                      <X className="w-4 h-4" />
                    </Button>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => setShowRichContentDialog(null)}>
            Cancel
          </Button>
          <Button onClick={sendInvoice}>
            <Send className="w-4 h-4 mr-2" />
            Send Invoice
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );

  const renderPollDialog = () => (
    <Dialog open={showRichContentDialog === 'poll'} onOpenChange={() => setShowRichContentDialog(null)}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <BarChart3 className="w-5 h-5" />
            Create Poll
          </DialogTitle>
          <DialogDescription>
            Ask a question and get responses from chat participants
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div>
            <label className="text-sm font-medium">Question</label>
            <Input
              placeholder="What's your question?"
              value={pollForm.question}
              onChange={(e) => setPollForm(prev => ({ ...prev, question: e.target.value }))}
            />
          </div>

          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="text-sm font-medium">Options</label>
              <Button size="sm" onClick={addPollOption}>
                Add Option
              </Button>
            </div>
            <div className="space-y-2">
              {pollForm.options.map((option, index) => (
                <div key={index} className="flex gap-2">
                  <Input
                    placeholder={`Option ${index + 1}`}
                    value={option}
                    onChange={(e) => {
                      const newOptions = [...pollForm.options];
                      newOptions[index] = e.target.value;
                      setPollForm(prev => ({ ...prev, options: newOptions }));
                    }}
                  />
                  {pollForm.options.length > 2 && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => removePollOption(index)}
                    >
                      <X className="w-4 h-4" />
                    </Button>
                  )}
                </div>
              ))}
            </div>
          </div>

          <div className="flex items-center gap-4">
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={pollForm.allowMultiple}
                onChange={(e) => setPollForm(prev => ({ ...prev, allowMultiple: e.target.checked }))}
              />
              <span className="text-sm">Allow multiple answers</span>
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={pollForm.anonymous}
                onChange={(e) => setPollForm(prev => ({ ...prev, anonymous: e.target.checked }))}
              />
              <span className="text-sm">Anonymous voting</span>
            </label>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => setShowRichContentDialog(null)}>
            Cancel
          </Button>
          <Button onClick={sendPoll}>
            Create Poll
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );

  const renderContactDialog = () => (
    <Dialog open={showRichContentDialog === 'contact'} onOpenChange={() => setShowRichContentDialog(null)}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Contact className="w-5 h-5" />
            Share Contact
          </DialogTitle>
          <DialogDescription>
            Share contact information in the chat
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <Input
            placeholder="Full Name *"
            value={contactForm.name}
            onChange={(e) => setContactForm(prev => ({ ...prev, name: e.target.value }))}
          />
          <Input
            placeholder="Phone Number"
            value={contactForm.phone}
            onChange={(e) => setContactForm(prev => ({ ...prev, phone: e.target.value }))}
          />
          <Input
            placeholder="Email Address"
            value={contactForm.email}
            onChange={(e) => setContactForm(prev => ({ ...prev, email: e.target.value }))}
          />
          <Input
            placeholder="Company"
            value={contactForm.company}
            onChange={(e) => setContactForm(prev => ({ ...prev, company: e.target.value }))}
          />
          <Input
            placeholder="Job Title"
            value={contactForm.title}
            onChange={(e) => setContactForm(prev => ({ ...prev, title: e.target.value }))}
          />
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => setShowRichContentDialog(null)}>
            Cancel
          </Button>
          <Button onClick={sendContact}>
            Share Contact
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );

  const renderPaymentDialog = () => (
    <Dialog open={showRichContentDialog === 'payment'} onOpenChange={() => setShowRichContentDialog(null)}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <CreditCard className="w-5 h-5" />
            Send Payment
          </DialogTitle>
          <DialogDescription>
            Send a payment to someone in the chat
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium">Amount</label>
              <Input
                type="number"
                step="0.01"
                placeholder="0.00"
                value={paymentForm.amount}
                onChange={(e) => setPaymentForm(prev => ({ ...prev, amount: parseFloat(e.target.value) || 0 }))}
              />
            </div>
            <div>
              <label className="text-sm font-medium">Currency</label>
              <select 
                className="w-full p-2 border rounded-md"
                value={paymentForm.currency}
                onChange={(e) => setPaymentForm(prev => ({ ...prev, currency: e.target.value }))}
              >
                <option value="THB">Thai Baht (฿)</option>
                <option value="USD">US Dollar ($)</option>
                <option value="IDR">Indonesian Rupiah (Rp)</option>
                <option value="VND">Vietnamese Dong (₫)</option>
                <option value="MYR">Malaysian Ringgit (RM)</option>
                <option value="PHP">Philippine Peso (₱)</option>
                <option value="SGD">Singapore Dollar (S$)</option>
              </select>
            </div>
          </div>

          <div>
            <label className="text-sm font-medium">Payment Method</label>
            <select 
              className="w-full p-2 border rounded-md"
              value={paymentForm.method}
              onChange={(e) => setPaymentForm(prev => ({ ...prev, method: e.target.value }))}
            >
              <option value="PromptPay">PromptPay (Thailand)</option>
              <option value="QRIS">QRIS (Indonesia)</option>
              <option value="VietQR">VietQR (Vietnam)</option>
              <option value="DuitNow">DuitNow QR (Malaysia)</option>
              <option value="InstaPay">InstaPay QR (Philippines)</option>
            </select>
          </div>

          <Input
            placeholder="Recipient (optional)"
            value={paymentForm.recipient}
            onChange={(e) => setPaymentForm(prev => ({ ...prev, recipient: e.target.value }))}
          />
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => setShowRichContentDialog(null)}>
            Cancel
          </Button>
          <Button onClick={sendPayment}>
            Send Payment
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );

  return (
    <div className="border-t border-border bg-card">
      {/* Render Dialogs */}
      {renderOrderDialog()}
      {renderInvoiceDialog()}
      {renderPollDialog()}
      {renderContactDialog()}  
      {renderPaymentDialog()}

      {/* Reply/Edit Preview */}
      {(replyToMessage || editingMessage) && (
        <div className="px-4 py-2 border-b border-border bg-muted/50">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <p className="text-xs text-muted-foreground mb-1">
                {editingMessage ? 'Editing message' : `Replying to ${replyToMessage?.senderName}`}
              </p>
              <p className="text-sm truncate">
                {editingMessage?.content || replyToMessage?.content}
              </p>
            </div>
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={editingMessage ? cancelEdit : cancelReply}
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
        </div>
      )}

      {/* Voice Recording UI */}
      {isRecording && (
        <div className="px-4 py-3 bg-red-50 border-b border-border">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-3 h-3 bg-red-500 rounded-full animate-pulse"></div>
              <span className="text-sm font-medium text-red-700">
                Recording... {formatRecordingTime(recordingDuration)}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={onCancelRecording}>
                Cancel
              </Button>
              <Button size="sm" onClick={onStopRecording}>
                <Pause className="w-4 h-4 mr-1" />
                Send
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Main Input Area */}
      <div className="fixed bottom-16 left-0 right-0 z-40 md:relative md:bottom-auto md:left-auto md:right-auto md:z-auto bg-card/95 backdrop-blur-sm border-t border-border md:border-t-0 p-3 md:p-4"
        style={{
          paddingBottom: `max(0.75rem, env(safe-area-inset-bottom) + 0.75rem)`
        }}
      >
        <div className="flex items-end gap-2">
          {/* Attachment Menu */}
          <DropdownMenu open={showAttachmentMenu} onOpenChange={setShowAttachmentMenu}>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="flex-shrink-0">
                <Paperclip className="w-5 h-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-64">
              <DropdownMenuLabel>Rich Content</DropdownMenuLabel>
              <DropdownMenuSeparator />
              
              <DropdownMenuItem onClick={() => setShowRichContentDialog('order')}>
                <ShoppingCart className="w-4 h-4 mr-2" />
                Create Order
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={() => setShowRichContentDialog('invoice')}>
                <Receipt className="w-4 h-4 mr-2" />
                Invoice & Bill
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={() => setShowRichContentDialog('poll')}>
                <BarChart3 className="w-4 h-4 mr-2" />
                Create Poll
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={() => setShowRichContentDialog('contact')}>
                <Contact className="w-4 h-4 mr-2" />
                Share Contact
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={() => setShowRichContentDialog('payment')}>
                <CreditCard className="w-4 h-4 mr-2" />
                Send Payment
              </DropdownMenuItem>
              
              <DropdownMenuSeparator />
              
              <DropdownMenuItem onClick={() => toast.success('Photo picker would open')}>
                <Image className="w-4 h-4 mr-2" />
                Photo & Video
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={() => toast.success('File picker would open')}>
                <File className="w-4 h-4 mr-2" />
                Document
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={() => toast.success('Location picker would open')}>
                <MapPin className="w-4 h-4 mr-2" />
                Location
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* Message Input */}
          <div className="flex-1 relative">
            <Input
              placeholder={
                editingMessage ? "Edit message..." : 
                replyToMessage ? "Reply..." : 
                "Type a message... Use **bold**, *italic*, `code`"
              }
              value={messageInput}
              onChange={(e) => setMessageInput(e.target.value)}
              onKeyPress={handleKeyPress}
              className="pr-10"
              disabled={isRecording}
            />
            
            {/* Emoji Button */}
            <Button 
              variant="ghost" 
              size="icon" 
              className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8"
              onClick={() => toast.success('Emoji picker would open')}
            >
              <Smile className="w-4 h-4" />
            </Button>
          </div>

          {/* Voice/Send Button */}
          {messageInput.trim() ? (
            <Button onClick={handleSendMessage} size="icon" className="flex-shrink-0">
              <Send className="w-5 h-5" />
            </Button>
          ) : (
            <Button 
              variant={isRecording ? "destructive" : "ghost"}
              size="icon" 
              className="flex-shrink-0"
              onClick={isRecording ? onStopRecording : onStartRecording}
            >
              <Mic className="w-5 h-5" />
            </Button>
          )}
        </div>

        {/* Rich Content Quick Actions */}
        <div className="flex items-center gap-2 mt-2 flex-wrap">
          {/* Markdown formatting */}
          <div className="flex items-center gap-1">
            <Button 
              size="sm" 
              variant="outline" 
              className="h-6 px-2 text-xs"
              onClick={() => insertMarkdownFormat('bold')}
            >
              <Bold className="w-3 h-3" />
            </Button>
            <Button 
              size="sm" 
              variant="outline" 
              className="h-6 px-2 text-xs"
              onClick={() => insertMarkdownFormat('italic')}
            >
              <Italic className="w-3 h-3" />
            </Button>
            <Button 
              size="sm" 
              variant="outline" 
              className="h-6 px-2 text-xs"
              onClick={() => insertMarkdownFormat('code')}
            >
              <Code className="w-3 h-3" />
            </Button>
          </div>

          {/* Quick stickers */}
          <div className="flex items-center gap-1">
            <Button 
              size="sm" 
              variant="outline" 
              className="h-6 px-2 text-xs"
              onClick={() => sendSticker('🎉', 'large')}
            >
              🎉
            </Button>
            <Button 
              size="sm" 
              variant="outline" 
              className="h-6 px-2 text-xs"
              onClick={() => sendSticker('👏', 'large')}
            >
              👏
            </Button>
            <Button 
              size="sm" 
              variant="outline" 
              className="h-6 px-2 text-xs"
              onClick={() => sendSticker('💯', 'large')}
            >
              💯
            </Button>
          </div>

          {/* Quick reactions */}
          <div className="flex items-center gap-1">
            <Badge 
              variant="outline" 
              className="text-xs cursor-pointer hover:bg-accent h-6" 
              onClick={() => setMessageInput(messageInput + "👍")}
            >
              👍
            </Badge>
            <Badge 
              variant="outline" 
              className="text-xs cursor-pointer hover:bg-accent h-6" 
              onClick={() => setMessageInput(messageInput + "❤️")}
            >
              ❤️
            </Badge>
            <Badge 
              variant="outline" 
              className="text-xs cursor-pointer hover:bg-accent h-6" 
              onClick={() => setMessageInput(messageInput + "😊")}
            >
              😊
            </Badge>
          </div>
        </div>
      </div>

      {/* Hidden File Inputs */}
      <input
        ref={fileInputRef}
        type="file"
        hidden
        accept=".pdf,.doc,.docx,.txt,.xls,.xlsx,.ppt,.pptx"
      />
      <input
        ref={imageInputRef}
        type="file"
        hidden
        accept="image/*,video/*"
        multiple
      />
    </div>
  );
}