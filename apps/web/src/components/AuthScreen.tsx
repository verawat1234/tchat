import React, { useState } from 'react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Badge } from './ui/badge';
import { MessageCircle, Phone, Mail, Shield, Zap, Globe, Wallet } from 'lucide-react';
import { toast } from "sonner";
import { useGetContentItemQuery } from '../services/contentApi';
import { useSelector, useDispatch } from 'react-redux';
import { selectFallbackContentById } from '../features/contentSelectors';
import { useRequestOTPMutation, useVerifyOTPMutation } from '../services/auth';
import { setTokens } from '../features/authSlice';

interface AuthScreenProps {
  onAuth: (userData: any) => void;
}

/**
 * Custom hook for getting authentication content with fallback support
 * Skip API calls on initial load to improve performance
 */
const useAuthContentWithFallback = (contentId: string, defaultText: string = '') => {
  const { data: contentData, isLoading, error } = useGetContentItemQuery(contentId, {
    // Skip API calls completely for auth screen - use static content for better performance
    skip: true,
  });
  const fallbackSelector = selectFallbackContentById(contentId);
  const fallbackContent = useSelector(fallbackSelector);

  // Use static default text instead of making API calls
  const content = defaultText;

  return {
    content,
    isLoading: false,
    hasError: false,
    hasFallback: false,
    hasRemoteContent: false
  };
};

export function AuthScreen({ onAuth }: AuthScreenProps) {
  const [phoneNumber, setPhoneNumber] = useState('+66812345678');
  const [otpCode, setOtpCode] = useState('');
  const [step, setStep] = useState<'input' | 'verify'>('input');
  const [loading, setLoading] = useState(false);

  // Redux hooks
  const dispatch = useDispatch();
  const [requestOTP, { isLoading: otpRequestLoading }] = useRequestOTPMutation();
  const [verifyOTP, { isLoading: otpVerifyLoading }] = useVerifyOTPMutation();

  // Dynamic content with fallbacks
  const appTitle = useAuthContentWithFallback('auth.app.title', 'Telegram SEA Edition');
  const appDescription = useAuthContentWithFallback('auth.app.description', 'Cloud messaging, payments, and social commerce built for Southeast Asia');
  const signInTitle = useAuthContentWithFallback('auth.signin.title', 'Sign In with Phone');
  const signInDescription = useAuthContentWithFallback('auth.signin.description', 'Enter your phone number to receive an OTP');
  const phonePlaceholder = useAuthContentWithFallback('auth.signin.phone.placeholder', '+66 XX XXX XXXX');
  const phoneHelperText = useAuthContentWithFallback('auth.signin.phone.helper', "We'll send you a 6-digit OTP via SMS");
  const sendOtpButton = useAuthContentWithFallback('auth.signin.button.sendotp', 'Send OTP');
  const sendingButton = useAuthContentWithFallback('auth.signin.button.sending', 'Sending...');
  const verifyTitle = useAuthContentWithFallback('auth.verify.title', 'Verify Your Phone');
  const verifyPhoneDescription = useAuthContentWithFallback('auth.verify.phone.description', 'We sent a code to {phoneNumber}');
  const verifyPlaceholder = useAuthContentWithFallback('auth.verify.placeholder', 'Enter verification code');
  const verifyButton = useAuthContentWithFallback('auth.verify.button.verify', 'Verify & Continue');
  const verifyingButton = useAuthContentWithFallback('auth.verify.button.verifying', 'Verifying...');
  const backButton = useAuthContentWithFallback('auth.verify.button.back', 'Back to Phone');
  const privacyText = useAuthContentWithFallback('auth.privacy.text', 'By continuing, you agree to our Terms of Service and Privacy Policy. Built for PDPA (TH/MY), PDP (ID) compliance.');

  // Feature labels
  const featureEncrypted = useAuthContentWithFallback('auth.features.encrypted', 'End-to-End Encrypted');
  const featureLowData = useAuthContentWithFallback('auth.features.lowdata', 'Ultra Low Data');
  const featurePayments = useAuthContentWithFallback('auth.features.payments', 'QR Payments');
  const featureLanguages = useAuthContentWithFallback('auth.features.languages', 'SEA Languages');

  // Error messages
  const errorPhoneRequired = useAuthContentWithFallback('auth.errors.phone.required', 'Please enter your phone number');
  const errorCodeRequired = useAuthContentWithFallback('auth.errors.code.required', 'Please enter the verification code');

  // Success messages
  const successOtpSent = useAuthContentWithFallback('auth.success.otp.sent', 'OTP sent to {phoneNumber}');
  const successWelcome = useAuthContentWithFallback('auth.success.welcome', 'Welcome to Telegram SEA!');

  const handleRequestOTP = async () => {
    if (!phoneNumber) {
      toast.error(errorPhoneRequired.content);
      return;
    }

    try {
      setLoading(true);

      // Parse phone number to separate country code and phone number
      const phoneNumberMatch = phoneNumber.match(/^(\+\d{1,3})(\d+)$/);
      if (!phoneNumberMatch) {
        toast.error('Invalid phone number format');
        return;
      }

      const [, countryCodePrefix, phone] = phoneNumberMatch;

      // Convert country code prefix to country code
      const countryCodeMap: { [key: string]: string } = {
        '+66': 'TH',
        '+62': 'ID',
        '+63': 'PH',
        '+84': 'VN',
        '+60': 'MY',
        '+65': 'SG'
      };

      const countryCode = countryCodeMap[countryCodePrefix] || 'TH';

      // Call the real OTP request API with correct field names
      const result = await requestOTP({
        phone_number: phoneNumber,
        country_code: countryCode
      }).unwrap();

      setStep('verify');
      toast.success(successOtpSent.content.replace('{phoneNumber}', phoneNumber));
    } catch (error: any) {
      console.error('OTP request failed:', error);
      toast.error(error?.data?.message || 'Failed to send OTP. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleVerify = async () => {
    if (!otpCode) {
      toast.error(errorCodeRequired.content);
      return;
    }

    try {
      setLoading(true);

      // Call the real OTP verification API
      const result = await verifyOTP({
        phoneNumber: phoneNumber,
        code: otpCode
      }).unwrap();

      // Store tokens in Redux
      dispatch(setTokens({
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
        expiresIn: result.expiresIn || 3600
      }));

      // Call the parent callback with user data
      onAuth(result.user);

      toast.success(successWelcome.content);
    } catch (error: any) {
      console.error('OTP verification failed:', error);
      toast.error(error?.data?.message || 'Invalid OTP code. Please try again.');
    } finally {
      setLoading(false);
    }
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
                verifyTitle.content
              )}
            </CardTitle>
            <CardDescription>
              {isContentLoading ? (
                <span className="block h-4 bg-gray-200 rounded animate-pulse mt-2" />
              ) : (
                verifyPhoneDescription.content.replace('{phoneNumber}', phoneNumber)
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
              disabled={loading || otpVerifyLoading || otpCode.length < 4 || isContentLoading}
              aria-label={loading || otpVerifyLoading ? verifyingButton.content : verifyButton.content}
            >
              {isContentLoading ? (
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 bg-gray-200 rounded animate-pulse" />
                  <div className="w-16 h-4 bg-gray-200 rounded animate-pulse" />
                </div>
              ) : (loading || otpVerifyLoading) ? (
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
              aria-label={backButton.content}
            >
              {isContentLoading ? (
                <div className="w-20 h-4 bg-gray-200 rounded animate-pulse" />
              ) : (
                backButton.content
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
              <span className="block h-6 w-96 bg-gray-200 rounded animate-pulse mx-auto" />
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
                  <span className="block h-4 w-20 bg-gray-200 rounded animate-pulse mx-auto" />
                ) : (
                  featureEncrypted.content
                )}
              </p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Zap className="w-6 h-6 text-chart-2 mb-2 mx-auto" />
              <p className="text-sm">
                {featureLowData.isLoading ? (
                  <span className="block h-4 w-20 bg-gray-200 rounded animate-pulse mx-auto" />
                ) : (
                  featureLowData.content
                )}
              </p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Wallet className="w-6 h-6 text-chart-3 mb-2 mx-auto" />
              <p className="text-sm">
                {featurePayments.isLoading ? (
                  <span className="block h-4 w-20 bg-gray-200 rounded animate-pulse mx-auto" />
                ) : (
                  featurePayments.content
                )}
              </p>
            </div>
            <div className="bg-card rounded-lg p-4 border">
              <Globe className="w-6 h-6 text-chart-4 mb-2 mx-auto" />
              <p className="text-sm">
                {featureLanguages.isLoading ? (
                  <span className="block h-4 w-20 bg-gray-200 rounded animate-pulse mx-auto" />
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
                <span className="block h-4 w-48 bg-gray-200 rounded animate-pulse" />
              ) : (
                signInDescription.content
              )}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
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
                <p className="text-xs text-muted-foreground mt-2">
                  Demo phone: +66812345678 → OTP: 123456
                </p>
              </div>
            </div>

            <Button
              onClick={handleRequestOTP}
              className="w-full mt-4"
              disabled={loading || otpRequestLoading || isMainContentLoading}
              aria-label={
                loading || otpRequestLoading
                  ? sendingButton.content
                  : sendOtpButton.content
              }
            >
              {isMainContentLoading ? (
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 bg-gray-200 rounded animate-pulse" />
                  <div className="w-20 h-4 bg-gray-200 rounded animate-pulse" />
                </div>
              ) : (loading || otpRequestLoading) ? (
                sendingButton.content
              ) : (
                sendOtpButton.content
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
