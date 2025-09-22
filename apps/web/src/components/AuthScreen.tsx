import React, { useState } from 'react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Badge } from './ui/badge';
import { MessageCircle, Phone, Mail, Shield, Zap, Globe, Wallet } from 'lucide-react';
import { toast } from "sonner";
import { useGetContentItemQuery } from '../services/contentApi';
import { useSelector } from 'react-redux';
import { selectFallbackContentById } from '../features/contentSelectors';

interface AuthScreenProps {
  onAuth: (userData: any) => void;
}

/**
 * Custom hook for getting authentication content with fallback support
 */
const useAuthContentWithFallback = (contentId: string, defaultText: string = '') => {
  const { data: contentData, isLoading, error } = useGetContentItemQuery(contentId);
  const fallbackSelector = selectFallbackContentById(contentId);
  const fallbackContent = useSelector(fallbackSelector);

  // Return content value with fallback hierarchy
  const content = contentData?.value?.text || fallbackContent?.text || defaultText;

  return {
    content,
    isLoading,
    hasError: !!error,
    hasFallback: !!fallbackContent,
    hasRemoteContent: !!contentData
  };
};

export function AuthScreen({ onAuth }: AuthScreenProps) {
  const [phoneNumber, setPhoneNumber] = useState('');
  const [email, setEmail] = useState('');
  const [otpCode, setOtpCode] = useState('');
  const [step, setStep] = useState<'input' | 'verify'>('input');
  const [authMethod, setAuthMethod] = useState<'phone' | 'email'>('phone');
  const [loading, setLoading] = useState(false);

  // Dynamic content with fallbacks
  const appTitle = useAuthContentWithFallback('auth.app.title', 'Telegram SEA Edition');
  const appDescription = useAuthContentWithFallback('auth.app.description', 'Cloud messaging, payments, and social commerce built for Southeast Asia');
  const signInTitle = useAuthContentWithFallback('auth.signin.title', 'Sign In');
  const signInDescription = useAuthContentWithFallback('auth.signin.description', 'Choose your preferred sign-in method');
  const phoneTabLabel = useAuthContentWithFallback('auth.signin.phone.tab', 'Phone OTP');
  const emailTabLabel = useAuthContentWithFallback('auth.signin.email.tab', 'Magic Link');
  const phonePlaceholder = useAuthContentWithFallback('auth.signin.phone.placeholder', '+66 XX XXX XXXX');
  const phoneHelperText = useAuthContentWithFallback('auth.signin.phone.helper', "We'll send you a 6-digit OTP via SMS");
  const emailPlaceholder = useAuthContentWithFallback('auth.signin.email.placeholder', 'your.email@example.com');
  const emailHelperText = useAuthContentWithFallback('auth.signin.email.helper', "We'll send you a secure magic link");
  const sendOtpButton = useAuthContentWithFallback('auth.signin.button.sendotp', 'Send OTP');
  const sendMagicLinkButton = useAuthContentWithFallback('auth.signin.button.magiclink', 'Send Magic Link');
  const sendingButton = useAuthContentWithFallback('auth.signin.button.sending', 'Sending...');
  const verifyTitle = useAuthContentWithFallback('auth.verify.title', 'Verify Your {method}');
  const verifyPhoneDescription = useAuthContentWithFallback('auth.verify.phone.description', 'We sent a code to {phoneNumber}');
  const verifyEmailDescription = useAuthContentWithFallback('auth.verify.email.description', 'Check your email at {email}');
  const verifyPlaceholder = useAuthContentWithFallback('auth.verify.placeholder', 'Enter verification code');
  const verifyButton = useAuthContentWithFallback('auth.verify.button.verify', 'Verify & Continue');
  const verifyingButton = useAuthContentWithFallback('auth.verify.button.verifying', 'Verifying...');
  const backButton = useAuthContentWithFallback('auth.verify.button.back', 'Back to {method}');
  const privacyText = useAuthContentWithFallback('auth.privacy.text', 'By continuing, you agree to our Terms of Service and Privacy Policy. Built for PDPA (TH/MY), PDP (ID) compliance.');

  // Feature labels
  const featureEncrypted = useAuthContentWithFallback('auth.features.encrypted', 'End-to-End Encrypted');
  const featureLowData = useAuthContentWithFallback('auth.features.lowdata', 'Ultra Low Data');
  const featurePayments = useAuthContentWithFallback('auth.features.payments', 'QR Payments');
  const featureLanguages = useAuthContentWithFallback('auth.features.languages', 'SEA Languages');

  // Error messages
  const errorPhoneRequired = useAuthContentWithFallback('auth.errors.phone.required', 'Please enter your phone number');
  const errorEmailRequired = useAuthContentWithFallback('auth.errors.email.required', 'Please enter your email address');
  const errorCodeRequired = useAuthContentWithFallback('auth.errors.code.required', 'Please enter the verification code');

  // Success messages
  const successOtpSent = useAuthContentWithFallback('auth.success.otp.sent', 'OTP sent to {phoneNumber}');
  const successMagicLinkSent = useAuthContentWithFallback('auth.success.magiclink.sent', 'Magic link sent to {email}');
  const successWelcome = useAuthContentWithFallback('auth.success.welcome', 'Welcome to Telegram SEA!');

  const handleSendCode = async () => {
    if (authMethod === 'phone' && !phoneNumber) {
      toast.error(errorPhoneRequired.content);
      return;
    }
    if (authMethod === 'email' && !email) {
      toast.error(errorEmailRequired.content);
      return;
    }

    setLoading(true);

    // Simulate API call
    setTimeout(() => {
      setLoading(false);
      setStep('verify');
      toast.success(
        authMethod === 'phone'
          ? successOtpSent.content.replace('{phoneNumber}', phoneNumber)
          : successMagicLinkSent.content.replace('{email}', email)
      );
    }, 1500);
  };

  const handleVerify = async () => {
    if (!otpCode) {
      toast.error(errorCodeRequired.content);
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
      toast.success(successWelcome.content);
    }, 1000);
  };

  // Check if any critical content is still loading
  const isContentLoading = verifyTitle.isLoading || verifyButton.isLoading || backButton.isLoading;

  if (step === 'verify') {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary/5 to-chart-1/5 flex items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center mx-auto mb-4">
              <MessageCircle className="w-8 h-8 text-primary-foreground" />
            </div>
            <CardTitle>
              {isContentLoading ? (
                <div className="h-6 bg-gray-200 rounded animate-pulse" />
              ) : (
                verifyTitle.content.replace('{method}', authMethod === 'phone' ? 'Phone' : 'Email')
              )}
            </CardTitle>
            <CardDescription>
              {isContentLoading ? (
                <div className="h-4 bg-gray-200 rounded animate-pulse mt-2" />
              ) : authMethod === 'phone' ? (
                verifyPhoneDescription.content.replace('{phoneNumber}', phoneNumber)
              ) : (
                verifyEmailDescription.content.replace('{email}', email)
              )}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Input
                placeholder={isContentLoading ? 'Loading...' : verifyPlaceholder.content}
                value={otpCode}
                onChange={(e) => setOtpCode(e.target.value)}
                maxLength={6}
                className="text-center text-lg tracking-wider"
                disabled={isContentLoading}
                aria-label={verifyPlaceholder.content}
              />
            </div>

            <Button
              onClick={handleVerify}
              className="w-full"
              disabled={loading || otpCode.length < 4 || isContentLoading}
              aria-label={loading ? verifyingButton.content : verifyButton.content}
            >
              {isContentLoading ? (
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 bg-gray-200 rounded animate-pulse" />
                  <div className="w-16 h-4 bg-gray-200 rounded animate-pulse" />
                </div>
              ) : loading ? (
                verifyingButton.content
              ) : (
                verifyButton.content
              )}
            </Button>

            <Button
              variant="ghost"
              onClick={() => setStep('input')}
              className="w-full"
              disabled={isContentLoading}
              aria-label={backButton.content.replace('{method}', authMethod === 'phone' ? 'Phone' : 'Email')}
            >
              {isContentLoading ? (
                <div className="w-20 h-4 bg-gray-200 rounded animate-pulse" />
              ) : (
                backButton.content.replace('{method}', authMethod === 'phone' ? 'Phone' : 'Email')
              )}
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Check if main content is loading
  const isMainContentLoading = appTitle.isLoading || appDescription.isLoading || signInTitle.isLoading;

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary/5 to-chart-1/5 flex items-center justify-center p-4">
      <div className="w-full max-w-4xl">
        {/* Hero Section */}
        <div className="text-center mb-8">
          <div className="flex items-center justify-center gap-3 mb-4">
            <div className="w-12 h-12 bg-primary rounded-full flex items-center justify-center">
              <MessageCircle className="w-6 h-6 text-primary-foreground" />
            </div>
            <h1 className="text-3xl font-bold">
              {isMainContentLoading ? (
                <div className="h-8 w-64 bg-gray-200 rounded animate-pulse" />
              ) : (
                appTitle.content
              )}
            </h1>
          </div>
          <p className="text-muted-foreground text-lg mb-6">
            {isMainContentLoading ? (
              <div className="h-6 w-96 bg-gray-200 rounded animate-pulse mx-auto" />
            ) : (
              appDescription.content
            )}
          </p>
          
          {/* Feature Highlights */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
            <div className="bg-card rounded-lg p-4 border">
              <Shield className="w-6 h-6 text-chart-1 mb-2 mx-auto" />
              <p className="text-sm">
                {featureEncrypted.isLoading ? (
                  <div className="h-4 w-20 bg-gray-200 rounded animate-pulse mx-auto" />
                ) : (
                  featureEncrypted.content
                )}
              </p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Zap className="w-6 h-6 text-chart-2 mb-2 mx-auto" />
              <p className="text-sm">
                {featureLowData.isLoading ? (
                  <div className="h-4 w-20 bg-gray-200 rounded animate-pulse mx-auto" />
                ) : (
                  featureLowData.content
                )}
              </p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Wallet className="w-6 h-6 text-chart-3 mb-2 mx-auto" />
              <p className="text-sm">
                {featurePayments.isLoading ? (
                  <div className="h-4 w-20 bg-gray-200 rounded animate-pulse mx-auto" />
                ) : (
                  featurePayments.content
                )}
              </p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Globe className="w-6 h-6 text-chart-4 mb-2 mx-auto" />
              <p className="text-sm">
                {featureLanguages.isLoading ? (
                  <div className="h-4 w-20 bg-gray-200 rounded animate-pulse mx-auto" />
                ) : (
                  featureLanguages.content
                )}
              </p>
            </div>
          </div>
        </div>

        {/* Auth Form */}
        <Card className="max-w-md mx-auto">
          <CardHeader>
            <CardTitle>
              {isMainContentLoading ? (
                <div className="h-6 w-24 bg-gray-200 rounded animate-pulse" />
              ) : (
                signInTitle.content
              )}
            </CardTitle>
            <CardDescription>
              {isMainContentLoading ? (
                <div className="h-4 w-48 bg-gray-200 rounded animate-pulse" />
              ) : (
                signInDescription.content
              )}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Tabs value={authMethod} onValueChange={(v) => setAuthMethod(v as 'phone' | 'email')}>
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger
                  value="phone"
                  className="flex items-center gap-2"
                  disabled={isMainContentLoading}
                  aria-label={phoneTabLabel.content}
                >
                  <Phone className="w-4 h-4" />
                  {phoneTabLabel.isLoading ? (
                    <div className="h-4 w-16 bg-gray-200 rounded animate-pulse" />
                  ) : (
                    phoneTabLabel.content
                  )}
                </TabsTrigger>
                <TabsTrigger
                  value="email"
                  className="flex items-center gap-2"
                  disabled={isMainContentLoading}
                  aria-label={emailTabLabel.content}
                >
                  <Mail className="w-4 h-4" />
                  {emailTabLabel.isLoading ? (
                    <div className="h-4 w-20 bg-gray-200 rounded animate-pulse" />
                  ) : (
                    emailTabLabel.content
                  )}
                </TabsTrigger>
              </TabsList>
              
              <TabsContent value="phone" className="space-y-4 mt-4">
                <div>
                  <div className="flex gap-2 mb-2">
                    <Badge
                      variant="outline"
                      className="cursor-pointer hover:bg-accent transition-colors touch-manipulation"
                      onClick={() => !isMainContentLoading && setPhoneNumber('+66 ')}
                      aria-label="Thailand +66"
                    >
                      🇹🇭 +66
                    </Badge>
                    <Badge
                      variant="outline"
                      className="cursor-pointer hover:bg-accent transition-colors touch-manipulation"
                      onClick={() => !isMainContentLoading && setPhoneNumber('+62 ')}
                      aria-label="Indonesia +62"
                    >
                      🇮🇩 +62
                    </Badge>
                    <Badge
                      variant="outline"
                      className="cursor-pointer hover:bg-accent transition-colors touch-manipulation"
                      onClick={() => !isMainContentLoading && setPhoneNumber('+63 ')}
                      aria-label="Philippines +63"
                    >
                      🇵🇭 +63
                    </Badge>
                    <Badge
                      variant="outline"
                      className="cursor-pointer hover:bg-accent transition-colors touch-manipulation"
                      onClick={() => !isMainContentLoading && setPhoneNumber('+84 ')}
                      aria-label="Vietnam +84"
                    >
                      🇻🇳 +84
                    </Badge>
                  </div>
                  <Input
                    placeholder={phonePlaceholder.isLoading ? 'Loading...' : phonePlaceholder.content}
                    value={phoneNumber}
                    onChange={(e) => setPhoneNumber(e.target.value)}
                    disabled={isMainContentLoading}
                    aria-label={phonePlaceholder.content}
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    {phoneHelperText.isLoading ? (
                      <div className="h-3 w-48 bg-gray-200 rounded animate-pulse" />
                    ) : (
                      phoneHelperText.content
                    )}
                  </p>
                </div>
              </TabsContent>
              
              <TabsContent value="email" className="space-y-4 mt-4">
                <div>
                  <Input
                    type="email"
                    placeholder={emailPlaceholder.isLoading ? 'Loading...' : emailPlaceholder.content}
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    disabled={isMainContentLoading}
                    aria-label={emailPlaceholder.content}
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    {emailHelperText.isLoading ? (
                      <div className="h-3 w-44 bg-gray-200 rounded animate-pulse" />
                    ) : (
                      emailHelperText.content
                    )}
                  </p>
                </div>
              </TabsContent>
            </Tabs>

            <Button
              onClick={handleSendCode}
              className="w-full mt-4"
              disabled={loading || isMainContentLoading}
              aria-label={
                loading
                  ? sendingButton.content
                  : authMethod === 'phone'
                  ? sendOtpButton.content
                  : sendMagicLinkButton.content
              }
            >
              {isMainContentLoading ? (
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 bg-gray-200 rounded animate-pulse" />
                  <div className="w-20 h-4 bg-gray-200 rounded animate-pulse" />
                </div>
              ) : loading ? (
                sendingButton.content
              ) : authMethod === 'phone' ? (
                sendOtpButton.content
              ) : (
                sendMagicLinkButton.content
              )}
            </Button>

            <p className="text-xs text-muted-foreground text-center mt-4">
              {privacyText.isLoading ? (
                <div className="space-y-1">
                  <div className="h-3 w-full bg-gray-200 rounded animate-pulse" />
                  <div className="h-3 w-3/4 bg-gray-200 rounded animate-pulse mx-auto" />
                </div>
              ) : (
                privacyText.content
              )}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}