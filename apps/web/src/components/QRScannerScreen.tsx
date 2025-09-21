import React, { useState } from 'react';
import { ArrowLeft, QrCode, Camera, Image, Flashlight, FlashlightOff, Upload, Zap, CreditCard, User, ShoppingCart } from 'lucide-react';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from './ui/dialog';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Separator } from './ui/separator';
import { toast } from "sonner";

interface QRScannerScreenProps {
  user: any;
  onBack: () => void;
}

interface QRResult {
  type: 'payment' | 'contact' | 'merchant' | 'product' | 'url';
  data: any;
}

export function QRScannerScreen({ user, onBack }: QRScannerScreenProps) {
  const [flashEnabled, setFlashEnabled] = useState(false);
  const [scanResult, setScanResult] = useState<QRResult | null>(null);
  const [showPayment, setShowPayment] = useState(false);
  const [paymentAmount, setPaymentAmount] = useState('');

  // Mock QR scan results for demo
  const mockScanResults = {
    promptpay: {
      type: 'payment' as const,
      data: {
        method: 'PromptPay',
        recipient: 'Somtam Vendor',
        phone: '+66 XX XXX XXXX',
        merchantId: 'merchant_123',
        amount: null, // Open amount
        avatar: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'
      }
    },
    product: {
      type: 'product' as const,
      data: {
        id: 'prod_123',
        name: 'Pad Thai Goong',
        price: 45,
        merchant: 'Bangkok Street Food',
        image: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'
      }
    },
    contact: {
      type: 'contact' as const,
      data: {
        name: 'John Doe',
        phone: '+66 XX XXX XXXX',
        avatar: null
      }
    }
  };

  const handleScanDemo = (type: keyof typeof mockScanResults) => {
    const result = mockScanResults[type];
    setScanResult(result);
    
    if (result.type === 'payment') {
      setShowPayment(true);
    } else if (result.type === 'product') {
      toast.success(`Found product: ${result.data.name}`);
    } else if (result.type === 'contact') {
      toast.success(`Contact found: ${result.data.name}`);
    }
  };

  const handlePayment = () => {
    if (!paymentAmount) {
      toast.error('Please enter payment amount');
      return;
    }
    
    toast.success(`Payment of ฿${paymentAmount} sent to ${scanResult?.data.recipient}`);
    setShowPayment(false);
    setScanResult(null);
    setPaymentAmount('');
  };

  const handleAddContact = () => {
    if (scanResult?.type === 'contact') {
      toast.success(`Added ${scanResult.data.name} to contacts`);
      setScanResult(null);
    }
  };

  const handleBuyProduct = () => {
    if (scanResult?.type === 'product') {
      toast.success(`Added ${scanResult.data.name} to cart`);
      setScanResult(null);
    }
  };

  return (
    <div className="h-screen mobile-full-height flex flex-col bg-black text-white min-h-0 overflow-hidden">
      {/* Header */}
      <header className="bg-black/80 backdrop-blur-sm px-4 py-3 flex items-center justify-between relative z-10 flex-shrink-0" style={{paddingTop: 'max(0.75rem, env(safe-area-inset-top, 0.75rem))' }}>
        <Button variant="ghost" size="icon" onClick={onBack} className="text-white hover:bg-white/20">
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h1 className="font-medium">QR Scanner</h1>
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setFlashEnabled(!flashEnabled)}
          className="text-white hover:bg-white/20"
        >
          {flashEnabled ? <FlashlightOff className="w-5 h-5" /> : <Flashlight className="w-5 h-5" />}
        </Button>
      </header>

      {/* Camera View */}
      <div className="flex-1 relative bg-gray-900 flex items-center justify-center min-h-0">
        {/* Mock camera view */}
        <div className="absolute inset-0 w-full h-full bg-gradient-to-b from-gray-800 to-gray-900 flex items-center justify-center">
          {/* Scan frame */}
          <div className="relative w-64 h-64 border-2 border-white rounded-lg">
            <div className="absolute top-0 left-0 w-6 h-6 border-t-4 border-l-4 border-primary rounded-tl-lg"></div>
            <div className="absolute top-0 right-0 w-6 h-6 border-t-4 border-r-4 border-primary rounded-tr-lg"></div>
            <div className="absolute bottom-0 left-0 w-6 h-6 border-b-4 border-l-4 border-primary rounded-bl-lg"></div>
            <div className="absolute bottom-0 right-0 w-6 h-6 border-b-4 border-r-4 border-primary rounded-br-lg"></div>
            
            {/* Scanning line animation */}
            <div className="absolute top-0 left-0 w-full h-1 bg-primary animate-pulse"></div>
          </div>
          
          {/* Instructions */}
          <div className="absolute bottom-24 sm:bottom-32 left-0 right-0 text-center px-4">
            <p className="text-white/80 mb-2">Point your camera at a QR code</p>
            <p className="text-sm text-white/60">
              Supports PromptPay, products, contacts, and more
            </p>
          </div>
        </div>
      </div>

      {/* Bottom Controls */}
      <div className="bg-black/80 backdrop-blur-sm p-4 space-y-4 flex-shrink-0" style={{paddingBottom: 'max(1rem, env(safe-area-inset-bottom, 1rem))' }}>
        {/* Quick Actions */}
        <div className="grid grid-cols-3 gap-4">
          <Button
            variant="outline"
            className="flex flex-col gap-2 h-auto py-4 bg-white/10 border-white/20 text-white hover:bg-white/20"
            onClick={() => handleScanDemo('promptpay')}
          >
            <CreditCard className="w-6 h-6" />
            <span className="text-xs">PromptPay</span>
          </Button>
          
          <Button
            variant="outline"
            className="flex flex-col gap-2 h-auto py-4 bg-white/10 border-white/20 text-white hover:bg-white/20"
            onClick={() => handleScanDemo('product')}
          >
            <ShoppingCart className="w-6 h-6" />
            <span className="text-xs">Product</span>
          </Button>
          
          <Button
            variant="outline"
            className="flex flex-col gap-2 h-auto py-4 bg-white/10 border-white/20 text-white hover:bg-white/20"
            onClick={() => handleScanDemo('contact')}
          >
            <User className="w-6 h-6" />
            <span className="text-xs">Contact</span>
          </Button>
        </div>

        {/* Upload from gallery */}
        <Button
          variant="outline"
          className="w-full bg-white/10 border-white/20 text-white hover:bg-white/20"
        >
          <Image className="w-5 h-5 mr-2" />
          Scan from Gallery
        </Button>
      </div>

      {/* Payment Dialog */}
      <Dialog open={showPayment} onOpenChange={setShowPayment}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>PromptPay Payment</DialogTitle>
            <DialogDescription>
              Complete your payment using PromptPay to the selected recipient.
            </DialogDescription>
          </DialogHeader>
          
          {scanResult?.type === 'payment' && (
            <div className="space-y-4">
              <div className="flex items-center gap-3 p-3 bg-muted rounded-lg">
                <Avatar className="w-12 h-12">
                  <AvatarImage src={scanResult.data.avatar} />
                  <AvatarFallback>
                    {scanResult.data.recipient.charAt(0)}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <p className="font-medium">{scanResult.data.recipient}</p>
                  <p className="text-sm text-muted-foreground">{scanResult.data.phone}</p>
                </div>
              </div>

              <div className="space-y-2">
                <Label>Amount</Label>
                <div className="relative">
                  <span className="absolute left-3 top-3 text-muted-foreground">฿</span>
                  <Input
                    placeholder="0.00"
                    value={paymentAmount}
                    onChange={(e) => setPaymentAmount(e.target.value)}
                    className="pl-8"
                    type="number"
                  />
                </div>
              </div>

              <div className="text-sm text-muted-foreground">
                Available balance: ฿1,250.00
              </div>

              <div className="grid grid-cols-2 gap-3">
                <Button variant="outline" onClick={() => setShowPayment(false)}>
                  Cancel
                </Button>
                <Button onClick={handlePayment}>
                  Pay Now
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Contact Result Dialog */}
      <Dialog open={scanResult?.type === 'contact'} onOpenChange={() => setScanResult(null)}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Contact Found</DialogTitle>
            <DialogDescription>
              Add this contact to your address book or view their details.
            </DialogDescription>
          </DialogHeader>
          
          {scanResult?.type === 'contact' && (
            <div className="space-y-4">
              <div className="text-center">
                <Avatar className="w-20 h-20 mx-auto mb-3">
                  <AvatarImage src={scanResult.data.avatar} />
                  <AvatarFallback className="text-2xl">
                    {scanResult.data.name.charAt(0)}
                  </AvatarFallback>
                </Avatar>
                <h3 className="font-medium text-lg">{scanResult.data.name}</h3>
                <p className="text-sm text-muted-foreground">{scanResult.data.phone}</p>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <Button variant="outline" onClick={() => setScanResult(null)}>
                  Cancel
                </Button>
                <Button onClick={handleAddContact}>
                  Add Contact
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Product Result Dialog */}
      <Dialog open={scanResult?.type === 'product'} onOpenChange={() => setScanResult(null)}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Product Found</DialogTitle>
            <DialogDescription>
              View product details and add to your cart or proceed with purchase.
            </DialogDescription>
          </DialogHeader>
          
          {scanResult?.type === 'product' && (
            <div className="space-y-4">
              <div className="text-center">
                <img
                  src={scanResult.data.image}
                  alt={scanResult.data.name}
                  className="w-32 h-32 object-cover rounded-lg mx-auto mb-3"
                />
                <h3 className="font-medium text-lg">{scanResult.data.name}</h3>
                <p className="text-sm text-muted-foreground">{scanResult.data.merchant}</p>
                <p className="text-xl font-bold text-primary mt-2">
                  ฿{scanResult.data.price}
                </p>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <Button variant="outline" onClick={() => setScanResult(null)}>
                  Cancel
                </Button>
                <Button onClick={handleBuyProduct}>
                  Add to Cart
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}