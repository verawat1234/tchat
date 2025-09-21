import React, { useState } from 'react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Badge } from './ui/badge';
import { MessageCircle, Phone, Mail, Shield, Zap, Globe, Wallet } from 'lucide-react';
import { toast } from "sonner";

interface AuthScreenProps {
  onAuth: (userData: any) => void;
}

export function AuthScreen({ onAuth }: AuthScreenProps) {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [email, setEmail] = useState('');
  const [otpCode, setOtpCode] = useState('');
  const [step, setStep] = useState<'input' | 'verify'>('input');
  const [authMethod, setAuthMethod] = useState<'phone' | 'email'>('phone');
  const [loading, setLoading] = useState(false);

  const handleSendCode = async () => {
    if (authMethod === 'phone' && !phoneNumber) {
      toast.error('Please enter your phone number');
      return;
    }
    if (authMethod === 'email' && !email) {
      toast.error('Please enter your email address');
      return;
    }

    setLoading(true);
    
    // Simulate API call
    setTimeout(() => {
      setLoading(false);
      setStep('verify');
      toast.success(
        authMethod === 'phone' 
          ? `OTP sent to ${phoneNumber}` 
          : `Magic link sent to ${email}`
      );
    }, 1500);
  };

  const handleVerify = async () => {
    if (!otpCode) {
      toast.error('Please enter the verification code');
      return;
    }

    setLoading(true);

    // Simulate verification
    setTimeout(() => {
      setLoading(false);
      const userData = {
        id: 'user_' + Date.now(),
        phone: phoneNumber,
        email: email,
        country: 'TH',
        locale: 'th-TH',
        kycTier: 1,
        name: authMethod === 'phone' ? phoneNumber : email.split('@')[0],
        avatar: null
      };
      onAuth(userData);
      toast.success('Welcome to Telegram SEA!');
    }, 1000);
  };

  if (step === 'verify') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary/5 to-chart-1/5 flex items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center mx-auto mb-4">
              <MessageCircle className="w-8 h-8 text-primary-foreground" />
            </div>
            <CardTitle>Verify Your {authMethod === 'phone' ? 'Phone' : 'Email'}</CardTitle>
            <CardDescription>
              {authMethod === 'phone' 
                ? `We sent a code to ${phoneNumber}` 
                : `Check your email at ${email}`}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Input
                placeholder={authMethod === 'phone' ? 'Enter 6-digit code' : 'Enter verification code'}
                value={otpCode}
                onChange={(e) => setOtpCode(e.target.value)}
                maxLength={6}
                className="text-center text-lg tracking-wider"
              />
            </div>
            
            <Button 
              onClick={handleVerify} 
              className="w-full" 
              disabled={loading || otpCode.length < 4}
            >
              {loading ? 'Verifying...' : 'Verify & Continue'}
            </Button>

            <Button 
              variant="ghost" 
              onClick={() => setStep('input')} 
              className="w-full"
            >
              Back to {authMethod === 'phone' ? 'Phone' : 'Email'}
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary/5 to-chart-1/5 flex items-center justify-center p-4">
      <div className="w-full max-w-4xl">
        {/* Hero Section */}
        <div className="text-center mb-8">
          <div className="flex items-center justify-center gap-3 mb-4">
            <div className="w-12 h-12 bg-primary rounded-full flex items-center justify-center">
              <MessageCircle className="w-6 h-6 text-primary-foreground" />
            </div>
            <h1 className="text-3xl font-bold">Telegram SEA Edition</h1>
          </div>
          <p className="text-muted-foreground text-lg mb-6">
            Cloud messaging, payments, and social commerce built for Southeast Asia
          </p>
          
          {/* Feature Highlights */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
            <div className="bg-card rounded-lg p-4 border">
              <Shield className="w-6 h-6 text-chart-1 mb-2 mx-auto" />
              <p className="text-sm">End-to-End Encrypted</p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Zap className="w-6 h-6 text-chart-2 mb-2 mx-auto" />
              <p className="text-sm">Ultra Low Data</p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Wallet className="w-6 h-6 text-chart-3 mb-2 mx-auto" />
              <p className="text-sm">QR Payments</p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Globe className="w-6 h-6 text-chart-4 mb-2 mx-auto" />
              <p className="text-sm">SEA Languages</p>
            </div>
          </div>
        </div>

        {/* Auth Form */}
        <Card className="max-w-md mx-auto">
          <CardHeader>
            <CardTitle>Sign In</CardTitle>
            <CardDescription>
              Choose your preferred sign-in method
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Tabs value={authMethod} onValueChange={(v) => setAuthMethod(v as 'phone' | 'email')}>
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="phone" className="flex items-center gap-2">
                  <Phone className="w-4 h-4" />
                  Phone OTP
                </TabsTrigger>
                <TabsTrigger value="email" className="flex items-center gap-2">
                  <Mail className="w-4 h-4" />
                  Magic Link
                </TabsTrigger>
              </TabsList>
              
              <TabsContent value="phone" className="space-y-4 mt-4">
                <div>
                  <div className="flex gap-2 mb-2">
                    <Badge 
                      variant="outline" 
                      className="cursor-pointer hover:bg-accent transition-colors touch-manipulation"
                      onClick={() => setPhoneNumber('+66 ')}
                    >
                      🇹🇭 +66
                    </Badge>
                    <Badge 
                      variant="outline" 
                      className="cursor-pointer hover:bg-accent transition-colors touch-manipulation"
                      onClick={() => setPhoneNumber('+62 ')}
                    >
                      🇮🇩 +62
                    </Badge>
                    <Badge 
                      variant="outline" 
                      className="cursor-pointer hover:bg-accent transition-colors touch-manipulation"
                      onClick={() => setPhoneNumber('+63 ')}
                    >
                      🇵🇭 +63
                    </Badge>
                    <Badge 
                      variant="outline" 
                      className="cursor-pointer hover:bg-accent transition-colors touch-manipulation"
                      onClick={() => setPhoneNumber('+84 ')}
                    >
                      🇻🇳 +84
                    </Badge>
                  </div>
                  <Input
                    placeholder="+66 XX XXX XXXX"
                    value={phoneNumber}
                    onChange={(e) => setPhoneNumber(e.target.value)}
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    We'll send you a 6-digit OTP via SMS
                  </p>
                </div>
              </TabsContent>
              
              <TabsContent value="email" className="space-y-4 mt-4">
                <div>
                  <Input
                    type="email"
                    placeholder="your.email@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    We'll send you a secure magic link
                  </p>
                </div>
              </TabsContent>
            </Tabs>

            <Button 
              onClick={handleSendCode} 
              className="w-full mt-4" 
              disabled={loading}
            >
              {loading ? 'Sending...' : `Send ${authMethod === 'phone' ? 'OTP' : 'Magic Link'}`}
            </Button>

            <p className="text-xs text-muted-foreground text-center mt-4">
              By continuing, you agree to our Terms of Service and Privacy Policy. 
              Built for PDPA (TH/MY), PDP (ID) compliance.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}