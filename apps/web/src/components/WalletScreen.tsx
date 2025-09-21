import React, { useState } from 'react';
import { ArrowLeft, Wallet, Plus, Send, QrCode, History, CreditCard, TrendingUp, Eye, EyeOff, ArrowUpRight, ArrowDownLeft, Copy, Share, MoreVertical } from 'lucide-react';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import { ScrollArea } from './ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Separator } from './ui/separator';
import { toast } from "sonner";

interface WalletScreenProps {
  user: any;
  onBack: () => void;
}

interface Transaction {
  id: string;
  type: 'send' | 'receive' | 'topup' | 'withdraw' | 'purchase';
  amount: number;
  currency: string;
  description: string;
  timestamp: string;
  status: 'completed' | 'pending' | 'failed';
  counterpart?: {
    name: string;
    avatar?: string;
  };
}

export function WalletScreen({ user, onBack }: WalletScreenProps) {
  const [showBalance, setShowBalance] = useState(true);
  const [showQRCode, setShowQRCode] = useState(false);
  const [showSendMoney, setShowSendMoney] = useState(false);
  const [sendAmount, setSendAmount] = useState('');
  const [sendRecipient, setSendRecipient] = useState('');

  const balance = 1250.00;
  const currency = 'THB';

  // Mock transaction data
  const transactions: Transaction[] = [
    {
      id: '1',
      type: 'receive',
      amount: 500,
      currency: 'THB',
      description: 'Payment from Mom',
      timestamp: '2 hours ago',
      status: 'completed',
      counterpart: {
        name: 'Mom',
        avatar: ''
      }
    },
    {
      id: '2',
      type: 'purchase',
      amount: -45,
      currency: 'THB',
      description: 'Pad Thai Goong - Somtam Vendor',
      timestamp: 'Yesterday',
      status: 'completed',
      counterpart: {
        name: 'Somtam Vendor',
        avatar: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'
      }
    },
    {
      id: '3',
      type: 'topup',
      amount: 1000,
      currency: 'THB',
      description: 'Bank transfer from Kasikorn',
      timestamp: '2 days ago',
      status: 'completed'
    },
    {
      id: '4',
      type: 'send',
      amount: -200,
      currency: 'THB',
      description: 'Lunch money for sister',
      timestamp: '3 days ago',
      status: 'completed',
      counterpart: {
        name: 'Sister',
        avatar: ''
      }
    },
    {
      id: '5',
      type: 'purchase',
      amount: -35,
      currency: 'THB',
      description: 'Som Tam Thai - Street Food Palace',
      timestamp: '1 week ago',
      status: 'completed',
      counterpart: {
        name: 'Street Food Palace',
        avatar: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'
      }
    }
  ];

  const handleSendMoney = () => {
    if (!sendAmount || !sendRecipient) {
      toast.error('Please fill in all fields');
      return;
    }
    
    toast.success(`Sent ฿${sendAmount} to ${sendRecipient}`);
    setShowSendMoney(false);
    setSendAmount('');
    setSendRecipient('');
  };

  const handleCopyQR = () => {
    toast.success('QR code copied to clipboard');
  };

  const getTransactionIcon = (type: string) => {
    switch (type) {
      case 'send':
        return <ArrowUpRight className="w-5 h-5 text-red-500" />;
      case 'receive':
        return <ArrowDownLeft className="w-5 h-5 text-green-500" />;
      case 'topup':
        return <Plus className="w-5 h-5 text-blue-500" />;
      case 'withdraw':
        return <ArrowUpRight className="w-5 h-5 text-orange-500" />;
      case 'purchase':
        return <ArrowUpRight className="w-5 h-5 text-purple-500" />;
      default:
        return <ArrowUpRight className="w-5 h-5" />;
    }
  };

  const formatAmount = (amount: number) => {
    const sign = amount >= 0 ? '+' : '';
    return `${sign}฿${Math.abs(amount).toLocaleString()}`;
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={onBack}>
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h1 className="font-medium">Wallet</h1>
        <div className="ml-auto">
          <Button variant="ghost" size="icon">
            <MoreVertical className="w-5 h-5" />
          </Button>
        </div>
      </header>

      <ScrollArea className="flex-1">
        <div className="p-3 sm:p-4 lg:p-6 space-y-4 sm:space-y-6 max-w-4xl mx-auto">
          {/* Balance Card */}
          <Card className="bg-gradient-to-br from-primary to-chart-1 text-primary-foreground">
            <CardContent className="p-4 sm:p-6">
              <div className="flex items-center justify-between mb-4 sm:mb-6">
                <div className="flex items-center gap-2 sm:gap-3">
                  <Wallet className="w-5 h-5 sm:w-6 sm:h-6" />
                  <span className="font-medium text-sm sm:text-base">Telegram Wallet</span>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-primary-foreground hover:bg-primary-foreground/20 w-9 h-9 sm:w-10 sm:h-10 touch-manipulation"
                  onClick={() => setShowBalance(!showBalance)}
                >
                  {showBalance ? <Eye className="w-4 h-4 sm:w-5 sm:h-5" /> : <EyeOff className="w-4 h-4 sm:w-5 sm:h-5" />}
                </Button>
              </div>
              
              <div className="mb-6 sm:mb-8">
                <p className="text-xs sm:text-sm opacity-90 mb-2">Available Balance</p>
                <div className="flex items-baseline gap-2 sm:gap-3">
                  <span className="text-2xl sm:text-3xl lg:text-4xl font-bold">
                    {showBalance ? `฿${balance.toLocaleString()}` : '฿••••••'}
                  </span>
                  <span className="text-base sm:text-lg opacity-75">{currency}</span>
                </div>
              </div>

              <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3">
                <Button
                  variant="secondary"
                  size="sm"
                  className="flex flex-col gap-1 h-auto py-3 sm:py-4 touch-manipulation mobile-button-large"
                  onClick={() => setShowSendMoney(true)}
                >
                  <Send className="w-5 h-5 sm:w-6 sm:h-6" />
                  <span className="text-[10px] sm:text-xs font-medium">Send</span>
                </Button>
                
                <Button
                  variant="secondary" 
                  size="sm"
                  className="flex flex-col gap-1 h-auto py-3 sm:py-4 touch-manipulation mobile-button-large"
                  onClick={() => setShowQRCode(true)}
                >
                  <QrCode className="w-5 h-5 sm:w-6 sm:h-6" />
                  <span className="text-[10px] sm:text-xs font-medium">Receive</span>
                </Button>
                
                <Button
                  variant="secondary"
                  size="sm"
                  className="flex flex-col gap-1 h-auto py-3 sm:py-4 touch-manipulation mobile-button-large"
                >
                  <Plus className="w-5 h-5 sm:w-6 sm:h-6" />
                  <span className="text-[10px] sm:text-xs font-medium">Top Up</span>
                </Button>
                
                <Button
                  variant="secondary"
                  size="sm"
                  className="flex flex-col gap-1 h-auto py-3 sm:py-4 touch-manipulation mobile-button-large"
                >
                  <TrendingUp className="w-5 h-5 sm:w-6 sm:h-6" />
                  <span className="text-[10px] sm:text-xs font-medium">Invest</span>
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Quick Actions */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
            <Card className="touch-manipulation">
              <CardContent className="p-3 sm:p-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 sm:w-12 sm:h-12 bg-blue-500 rounded-lg flex items-center justify-center flex-shrink-0">
                    <QrCode className="w-5 h-5 sm:w-6 sm:h-6 text-white" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="font-medium text-sm sm:text-base">PromptPay</p>
                    <p className="text-xs sm:text-sm text-muted-foreground">Scan QR to pay</p>
                  </div>
                </div>
              </CardContent>
            </Card>
            
            <Card className="touch-manipulation">
              <CardContent className="p-3 sm:p-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 sm:w-12 sm:h-12 bg-green-500 rounded-lg flex items-center justify-center flex-shrink-0">
                    <CreditCard className="w-5 h-5 sm:w-6 sm:h-6 text-white" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="font-medium text-sm sm:text-base">Cards</p>
                    <p className="text-xs sm:text-sm text-muted-foreground">Manage cards</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* KYC Status */}
          <Card>
            <CardHeader className="p-3 sm:p-4 lg:p-6">
              <CardTitle className="flex items-center justify-between text-sm sm:text-base">
                <span>Account Status</span>
                <Badge variant="secondary" className="text-xs">KYC Tier {user.kycTier}</Badge>
              </CardTitle>
            </CardHeader>
            <CardContent className="p-3 sm:p-4 lg:p-6 pt-0">
              <div className="space-y-3 sm:space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm sm:text-base">Daily limit</span>
                  <span className="font-medium text-sm sm:text-base">฿10,000</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm sm:text-base">Monthly limit</span>
                  <span className="font-medium text-sm sm:text-base">฿50,000</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm sm:text-base">Used this month</span>
                  <span className="text-chart-1 font-medium text-sm sm:text-base">฿12,500</span>
                </div>
                {user.kycTier < 3 && (
                  <Button variant="outline" className="w-full mt-4 touch-manipulation mobile-button-large">
                    Upgrade Verification
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Recent Transactions */}
          <Card>
            <CardHeader className="p-3 sm:p-4 lg:p-6">
              <CardTitle className="flex items-center justify-between text-sm sm:text-base">
                <span>Recent Transactions</span>
                <Button variant="ghost" size="sm" className="touch-manipulation">
                  <History className="w-3 h-3 sm:w-4 sm:h-4 mr-1 sm:mr-2" />
                  <span className="text-xs sm:text-sm">View All</span>
                </Button>
              </CardTitle>
            </CardHeader>
            <CardContent className="p-3 sm:p-4 lg:p-6 pt-0">
              <div className="space-y-2 sm:space-y-3">
                {transactions.slice(0, 5).map((transaction) => (
                  <div key={transaction.id} className="flex items-center gap-3 p-2 sm:p-3 rounded-lg hover:bg-muted/50 transition-colors touch-manipulation">
                    <div className="flex-shrink-0">
                      {transaction.counterpart?.avatar ? (
                        <Avatar className="w-8 h-8 sm:w-10 sm:h-10">
                          <AvatarImage src={transaction.counterpart.avatar} />
                          <AvatarFallback className="text-xs sm:text-sm">
                            {transaction.counterpart.name.charAt(0)}
                          </AvatarFallback>
                        </Avatar>
                      ) : (
                        <div className="w-8 h-8 sm:w-10 sm:h-10 bg-muted rounded-full flex items-center justify-center">
                          {getTransactionIcon(transaction.type)}
                        </div>
                      )}
                    </div>
                    
                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate text-sm sm:text-base">{transaction.description}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <p className="text-xs sm:text-sm text-muted-foreground">{transaction.timestamp}</p>
                        <Badge
                          variant={
                            transaction.status === 'completed' ? 'secondary' :
                            transaction.status === 'pending' ? 'secondary' : 'destructive'
                          }
                          className="text-[10px] sm:text-xs"
                        >
                          {transaction.status}
                        </Badge>
                      </div>
                    </div>
                    
                    <div className="text-right flex-shrink-0">
                      <p className={`font-medium text-sm sm:text-base ${
                        transaction.amount >= 0 ? 'text-green-600' : 'text-red-600'
                      }`}>
                        {formatAmount(transaction.amount)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </ScrollArea>

      {/* QR Code Dialog */}
      <Dialog open={showQRCode} onOpenChange={setShowQRCode}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Receive Payment</DialogTitle>
          </DialogHeader>
          
          <div className="text-center space-y-4">
            <div className="w-48 h-48 mx-auto bg-muted rounded-lg flex items-center justify-center">
              <QrCode className="w-24 h-24 text-muted-foreground" />
            </div>
            
            <div>
              <p className="font-medium">{user.name}</p>
              <p className="text-sm text-muted-foreground">{user.phone}</p>
            </div>
            
            <div className="grid grid-cols-2 gap-3">
              <Button variant="outline" onClick={handleCopyQR}>
                <Copy className="w-4 h-4 mr-2" />
                Copy
              </Button>
              <Button variant="outline">
                <Share className="w-4 h-4 mr-2" />
                Share
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      {/* Send Money Dialog */}
      <Dialog open={showSendMoney} onOpenChange={setShowSendMoney}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Send Money</DialogTitle>
          </DialogHeader>
          
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Recipient</Label>
              <Input
                placeholder="Phone number or name"
                value={sendRecipient}
                onChange={(e) => setSendRecipient(e.target.value)}
              />
            </div>
            
            <div className="space-y-2">
              <Label>Amount</Label>
              <div className="relative">
                <span className="absolute left-3 top-3 text-muted-foreground">฿</span>
                <Input
                  placeholder="0.00"
                  value={sendAmount}
                  onChange={(e) => setSendAmount(e.target.value)}
                  className="pl-8"
                  type="number"
                />
              </div>
            </div>
            
            <div className="text-sm text-muted-foreground">
              Available balance: ฿{balance.toLocaleString()}
            </div>
            
            <div className="grid grid-cols-2 gap-3">
              <Button variant="outline" onClick={() => setShowSendMoney(false)}>
                Cancel
              </Button>
              <Button onClick={handleSendMoney}>
                Send Money
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}