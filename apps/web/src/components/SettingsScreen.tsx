import React, { useState } from 'react';
import { ArrowLeft, User, Shield, Globe, Bell, Wallet, Smartphone, HelpCircle, LogOut, ChevronRight, Moon, Sun, Monitor, Phone, Mail, Lock, Eye, EyeOff, QrCode, CreditCard, MapPin, Volume2, Languages, Calendar, Database, UserCheck } from 'lucide-react';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Switch } from './ui/switch';
import { Badge } from './ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Separator } from './ui/separator';
import { ScrollArea } from './ui/scroll-area';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import { Slider } from './ui/slider';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { toast } from "sonner";

interface SettingsScreenProps {
  user: any;
  onBack: () => void;
  onLogout: () => void;
}

type SettingsPage = 'main' | 'profile' | 'privacy' | 'notifications' | 'language' | 'payments' | 'data' | 'help';

export function SettingsScreen({ user, onBack, onLogout }: SettingsScreenProps) {
  const [currentPage, setCurrentPage] = useState<SettingsPage>('main');
  const [theme, setTheme] = useState('system');
  const [language, setLanguage] = useState('th');
  const [calendar, setCalendar] = useState('buddhist');
  const [e2eeEnabled, setE2eeEnabled] = useState(true);
  const [biometricAuth, setBiometricAuth] = useState(false);
  const [voiceEnabled, setVoiceEnabled] = useState(true);
  const [lowDataMode, setLowDataMode] = useState(false);
  const [offlineMode, setOfflineMode] = useState(true);
  const [notifications, setNotifications] = useState({
    messages: true,
    payments: true,
    social: true,
    marketing: false
  });

  const languages = [
    { code: 'th', name: 'ไทย (Thai)', flag: '🇹🇭' },
    { code: 'id', name: 'Bahasa Indonesia', flag: '🇮🇩' },
    { code: 'vi', name: 'Tiếng Việt', flag: '🇻🇳' },
    { code: 'tl', name: 'Filipino', flag: '🇵🇭' },
    { code: 'ms', name: 'Bahasa Melayu', flag: '🇲🇾' },
    { code: 'en', name: 'English', flag: '🇺🇸' }
  ];

  const kycTiers = [
    { tier: 0, name: 'Basic', description: 'Chat only', limit: '฿0' },
    { tier: 1, name: 'Verified', description: 'P2P payments', limit: '฿50,000/month' },
    { tier: 2, name: 'Enhanced', description: 'Merchant payments', limit: '฿200,000/month' },
    { tier: 3, name: 'Premium', description: 'Cross-border', limit: 'Unlimited' }
  ];

  const handleNotificationToggle = (type: keyof typeof notifications) => {
    setNotifications(prev => ({
      ...prev,
      [type]: !prev[type]
    }));
    toast.success('Notification settings updated');
  };

  const SettingItem = ({ icon: Icon, title, subtitle, action, onClick }: any) => (
    <Button
      variant="ghost"
      className="w-full justify-start h-auto p-4"
      onClick={onClick}
    >
      <div className="flex items-center gap-3 w-full">
        <Icon className="w-5 h-5 text-muted-foreground" />
        <div className="flex-1 text-left">
          <p className="font-medium">{title}</p>
          {subtitle && <p className="text-sm text-muted-foreground">{subtitle}</p>}
        </div>
        {action || <ChevronRight className="w-4 h-4 text-muted-foreground" />}
      </div>
    </Button>
  );

  if (currentPage === 'profile') {
    return (
      <div className="h-full flex flex-col">
        <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => setCurrentPage('main')}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="font-medium">Profile</h1>
        </header>

        <ScrollArea className="flex-1 p-4">
          <div className="space-y-6">
            {/* Profile Picture */}
            <Card>
              <CardContent className="p-6 text-center">
                <Avatar className="w-24 h-24 mx-auto mb-4">
                  <AvatarImage src={user.avatar} />
                  <AvatarFallback className="text-2xl">
                    {user.name.charAt(0).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <h3 className="font-medium text-lg">{user.name}</h3>
                <p className="text-sm text-muted-foreground">{user.phone || user.email}</p>
                <Badge className="mt-2">
                  KYC Tier {user.kycTier}
                </Badge>
                <Button variant="outline" className="mt-4">
                  Change Photo
                </Button>
              </CardContent>
            </Card>

            {/* KYC Status */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <UserCheck className="w-5 h-5" />
                  Verification Status
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {kycTiers.map((tier) => (
                    <div
                      key={tier.tier}
                      className={`p-3 rounded-lg border ${
                        user.kycTier >= tier.tier 
                          ? 'bg-green-50 border-green-200 dark:bg-green-950/20 dark:border-green-800' 
                          : 'bg-muted'
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="font-medium">{tier.name}</p>
                          <p className="text-sm text-muted-foreground">{tier.description}</p>
                        </div>
                        <div className="text-right">
                          <p className="text-sm font-medium">{tier.limit}</p>
                          {user.kycTier >= tier.tier && (
                            <Badge variant="secondary" className="mt-1">Verified</Badge>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
                {user.kycTier < 3 && (
                  <Button className="w-full mt-4">
                    Upgrade Verification
                  </Button>
                )}
              </CardContent>
            </Card>

            {/* Contact Info */}
            <Card>
              <CardHeader>
                <CardTitle>Contact Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label>Phone Number</Label>
                  <div className="flex gap-2">
                    <Input value={user.phone || '+66 XX XXX XXXX'} disabled />
                    <Button variant="outline" size="sm">Change</Button>
                  </div>
                </div>
                <div className="space-y-2">
                  <Label>Email Address</Label>
                  <div className="flex gap-2">
                    <Input value={user.email || 'user@example.com'} disabled />
                    <Button variant="outline" size="sm">Change</Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </ScrollArea>
      </div>
    );
  }

  if (currentPage === 'privacy') {
    return (
      <div className="h-full flex flex-col">
        <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => setCurrentPage('main')}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="font-medium">Privacy & Security</h1>
        </header>

        <ScrollArea className="flex-1 p-4">
          <div className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="w-5 h-5" />
                  Encryption
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">End-to-End Encryption</p>
                    <p className="text-sm text-muted-foreground">Enable E2EE for secret chats</p>
                  </div>
                  <Switch checked={e2eeEnabled} onCheckedChange={setE2eeEnabled} />
                </div>
                <Separator />
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Biometric Authentication</p>
                    <p className="text-sm text-muted-foreground">Use fingerprint/face unlock</p>
                  </div>
                  <Switch checked={biometricAuth} onCheckedChange={setBiometricAuth} />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Data & Privacy</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <SettingItem
                  icon={Database}
                  title="Data Export"
                  subtitle="Download your data"
                  onClick={() => toast.success('Data export requested')}
                />
                <SettingItem
                  icon={Lock}
                  title="Account Deletion"
                  subtitle="Permanently delete your account"
                  onClick={() => toast.error('Account deletion requested')}
                />
                <SettingItem
                  icon={Eye}
                  title="Privacy Policy"
                  subtitle="PDPA (TH/MY), PDP (ID) compliant"
                />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Sessions</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center gap-3">
                      <Smartphone className="w-5 h-5" />
                      <div>
                        <p className="font-medium">iPhone 15 Pro</p>
                        <p className="text-sm text-muted-foreground">Current session • Bangkok</p>
                      </div>
                    </div>
                    <Badge variant="secondary">Active</Badge>
                  </div>
                  <Button variant="destructive" className="w-full">
                    Terminate All Other Sessions
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </ScrollArea>
      </div>
    );
  }

  if (currentPage === 'language') {
    return (
      <div className="h-full flex flex-col">
        <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => setCurrentPage('main')}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="font-medium">Language & Region</h1>
        </header>

        <ScrollArea className="flex-1 p-4">
          <div className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Languages className="w-5 h-5" />
                  Interface Language
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {languages.map((lang) => (
                    <Button
                      key={lang.code}
                      variant={language === lang.code ? 'default' : 'ghost'}
                      className="w-full justify-start"
                      onClick={() => {
                        setLanguage(lang.code);
                        toast.success(`Language changed to ${lang.name}`);
                      }}
                    >
                      <span className="mr-3">{lang.flag}</span>
                      {lang.name}
                    </Button>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Calendar className="w-5 h-5" />
                  Calendar System
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <Button
                    variant={calendar === 'buddhist' ? 'default' : 'ghost'}
                    className="w-full justify-start"
                    onClick={() => setCalendar('buddhist')}
                  >
                    🇹🇭 Buddhist Calendar (BE 2567)
                  </Button>
                  <Button
                    variant={calendar === 'gregorian' ? 'default' : 'ghost'}
                    className="w-full justify-start"
                    onClick={() => setCalendar('gregorian')}
                  >
                    🌍 Gregorian Calendar (2024)
                  </Button>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Volume2 className="w-5 h-5" />
                  Voice & Audio
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Voice Messages</p>
                    <p className="text-sm text-muted-foreground">Enable voice note transcription</p>
                  </div>
                  <Switch checked={voiceEnabled} onCheckedChange={setVoiceEnabled} />
                </div>
                <div className="space-y-2">
                  <Label>Voice Recognition Language</Label>
                  <Select defaultValue="th">
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="th">🇹🇭 Thai</SelectItem>
                      <SelectItem value="id">🇮🇩 Indonesian</SelectItem>
                      <SelectItem value="vi">🇻🇳 Vietnamese</SelectItem>
                      <SelectItem value="en">🇺🇸 English</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </Card>
          </div>
        </ScrollArea>
      </div>
    );
  }

  if (currentPage === 'notifications') {
    return (
      <div className="h-full flex flex-col">
        <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => setCurrentPage('main')}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="font-medium">Notifications</h1>
        </header>

        <ScrollArea className="flex-1 p-4">
          <div className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Bell className="w-5 h-5" />
                  Push Notifications
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Messages</p>
                    <p className="text-sm text-muted-foreground">Chat messages and replies</p>
                  </div>
                  <Switch 
                    checked={notifications.messages} 
                    onCheckedChange={() => handleNotificationToggle('messages')} 
                  />
                </div>
                <Separator />
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Payments</p>
                    <p className="text-sm text-muted-foreground">Payment confirmations and receipts</p>
                  </div>
                  <Switch 
                    checked={notifications.payments} 
                    onCheckedChange={() => handleNotificationToggle('payments')} 
                  />
                </div>
                <Separator />
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Social</p>
                    <p className="text-sm text-muted-foreground">Likes, comments, and mentions</p>
                  </div>
                  <Switch 
                    checked={notifications.social} 
                    onCheckedChange={() => handleNotificationToggle('social')} 
                  />
                </div>
                <Separator />
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Marketing</p>
                    <p className="text-sm text-muted-foreground">Promotions and updates</p>
                  </div>
                  <Switch 
                    checked={notifications.marketing} 
                    onCheckedChange={() => handleNotificationToggle('marketing')} 
                  />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Quiet Hours</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span>Do not disturb</span>
                    <Switch />
                  </div>
                  <div className="space-y-2">
                    <Label>From</Label>
                    <Select defaultValue="22:00">
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="22:00">10:00 PM</SelectItem>
                        <SelectItem value="23:00">11:00 PM</SelectItem>
                        <SelectItem value="00:00">12:00 AM</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label>To</Label>
                    <Select defaultValue="07:00">
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="06:00">6:00 AM</SelectItem>
                        <SelectItem value="07:00">7:00 AM</SelectItem>
                        <SelectItem value="08:00">8:00 AM</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </ScrollArea>
      </div>
    );
  }

  if (currentPage === 'payments') {
    return (
      <div className="h-full flex flex-col">
        <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => setCurrentPage('main')}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="font-medium">Payments & Wallet</h1>
        </header>

        <ScrollArea className="flex-1 p-4">
          <div className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Wallet className="w-5 h-5" />
                  Wallet Balance
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-center py-4">
                  <p className="text-3xl font-bold">฿1,250.00</p>
                  <p className="text-sm text-muted-foreground">Available Balance</p>
                  <div className="flex gap-2 mt-4">
                    <Button className="flex-1">Top Up</Button>
                    <Button variant="outline" className="flex-1">Withdraw</Button>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <QrCode className="w-5 h-5" />
                  Payment Methods
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-blue-500 rounded-lg flex items-center justify-center">
                      <QrCode className="w-5 h-5 text-white" />
                    </div>
                    <div>
                      <p className="font-medium">PromptPay</p>
                      <p className="text-sm text-muted-foreground">+66 XX XXX XXXX</p>
                    </div>
                  </div>
                  <Badge variant="secondary">Primary</Badge>
                </div>
                
                <div className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-green-500 rounded-lg flex items-center justify-center">
                      <CreditCard className="w-5 h-5 text-white" />
                    </div>
                    <div>
                      <p className="font-medium">Kasikorn Bank</p>
                      <p className="text-sm text-muted-foreground">**** 1234</p>
                    </div>
                  </div>
                  <Button variant="ghost" size="sm">Remove</Button>
                </div>

                <Button variant="outline" className="w-full">
                  <Plus className="w-4 h-4 mr-2" />
                  Add Payment Method
                </Button>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Transaction Limits</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span>Daily limit</span>
                    <span className="font-medium">฿10,000</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Monthly limit</span>
                    <span className="font-medium">฿50,000</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Used this month</span>
                    <span className="text-chart-1">฿12,500</span>
                  </div>
                  <Button variant="outline" className="w-full mt-4">
                    Request Limit Increase
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </ScrollArea>
      </div>
    );
  }

  if (currentPage === 'data') {
    return (
      <div className="h-full flex flex-col">
        <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => setCurrentPage('main')}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="font-medium">Data & Storage</h1>
        </header>

        <ScrollArea className="flex-1 p-4">
          <div className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="w-5 h-5" />
                  Data Usage
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Low Data Mode</p>
                    <p className="text-sm text-muted-foreground">Optimize for 3G networks</p>
                  </div>
                  <Switch checked={lowDataMode} onCheckedChange={setLowDataMode} />
                </div>
                <Separator />
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Offline Mode</p>
                    <p className="text-sm text-muted-foreground">Store messages locally</p>
                  </div>
                  <Switch checked={offlineMode} onCheckedChange={setOfflineMode} />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Storage Usage</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span>Messages</span>
                      <span>45 MB</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div className="bg-primary h-2 rounded-full" style={{ width: '30%' }}></div>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span>Media</span>
                      <span>120 MB</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div className="bg-chart-1 h-2 rounded-full" style={{ width: '80%' }}></div>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span>Cache</span>
                      <span>25 MB</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div className="bg-chart-2 h-2 rounded-full" style={{ width: '17%' }}></div>
                    </div>
                  </div>
                  
                  <Button variant="outline" className="w-full mt-4">
                    Clear Cache
                  </Button>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Auto-Download</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="space-y-2">
                    <Label>Photos</Label>
                    <Select defaultValue="wifi">
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="never">Never</SelectItem>
                        <SelectItem value="wifi">Wi-Fi only</SelectItem>
                        <SelectItem value="always">Always</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  
                  <div className="space-y-2">
                    <Label>Videos</Label>
                    <Select defaultValue="never">
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="never">Never</SelectItem>
                        <SelectItem value="wifi">Wi-Fi only</SelectItem>
                        <SelectItem value="always">Always</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </ScrollArea>
      </div>
    );
  }

  if (currentPage === 'help') {
    return (
      <div className="h-full flex flex-col">
        <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => setCurrentPage('main')}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <h1 className="font-medium">Help & Support</h1>
        </header>

        <ScrollArea className="flex-1 p-4">
          <div className="space-y-4">
            <Card>
              <CardContent className="p-6 text-center">
                <h3 className="font-medium text-lg mb-2">Need Help?</h3>
                <p className="text-sm text-muted-foreground mb-4">
                  Our support team is available 24/7 to assist you with any questions or issues.
                </p>
                <div className="grid grid-cols-2 gap-3">
                  <Button variant="outline">
                    <Phone className="w-4 h-4 mr-2" />
                    Call Support
                  </Button>
                  <Button variant="outline">
                    <Mail className="w-4 h-4 mr-2" />
                    Email Us
                  </Button>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>FAQ</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <SettingItem
                  icon={HelpCircle}
                  title="How to set up PromptPay?"
                  subtitle="Step-by-step guide"
                />
                <SettingItem
                  icon={HelpCircle}
                  title="Understanding KYC verification"
                  subtitle="Verification levels explained"
                />
                <SettingItem
                  icon={HelpCircle}
                  title="Payment security"
                  subtitle="How we protect your money"
                />
                <SettingItem
                  icon={HelpCircle}
                  title="Live commerce features"
                  subtitle="Shopping while watching streams"
                />
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>App Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-between">
                  <span>Version</span>
                  <span className="text-muted-foreground">2.1.0 (SEA Edition)</span>
                </div>
                <div className="flex justify-between">
                  <span>Build</span>
                  <span className="text-muted-foreground">2024.03.15</span>
                </div>
                <div className="flex justify-between">
                  <span>Region</span>
                  <span className="text-muted-foreground">🇹🇭 Thailand</span>
                </div>
                <Button variant="outline" className="w-full mt-4">
                  Check for Updates
                </Button>
              </CardContent>
            </Card>
          </div>
        </ScrollArea>
      </div>
    );
  }

  // Main settings page
  return (
    <div className="h-full flex flex-col">
      <header className="border-b border-border bg-card px-4 py-3 flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={onBack}>
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h1 className="font-medium">Settings</h1>
      </header>

      <ScrollArea className="flex-1">
        <div className="p-4 space-y-4">
          {/* Profile Section */}
          <Card>
            <CardContent className="p-4">
              <Button
                variant="ghost"
                className="w-full justify-start h-auto p-4"
                onClick={() => setCurrentPage('profile')}
              >
                <div className="flex items-center gap-3 w-full">
                  <Avatar className="w-12 h-12">
                    <AvatarImage src={user.avatar} />
                    <AvatarFallback>
                      {user.name.charAt(0).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 text-left">
                    <p className="font-medium">{user.name}</p>
                    <p className="text-sm text-muted-foreground">
                      {user.phone || user.email} • KYC Tier {user.kycTier}
                    </p>
                  </div>
                  <ChevronRight className="w-4 h-4 text-muted-foreground" />
                </div>
              </Button>
            </CardContent>
          </Card>

          {/* Settings Options */}
          <Card>
            <CardContent className="p-0">
              <SettingItem
                icon={Shield}
                title="Privacy & Security"
                subtitle="Encryption, sessions, and data protection"
                onClick={() => setCurrentPage('privacy')}
              />
              <Separator />
              <SettingItem
                icon={Bell}
                title="Notifications"
                subtitle="Manage alerts and quiet hours"
                onClick={() => setCurrentPage('notifications')}
              />
              <Separator />
              <SettingItem
                icon={Globe}
                title="Language & Region"
                subtitle={`${languages.find(l => l.code === language)?.flag} ${languages.find(l => l.code === language)?.name}`}
                onClick={() => setCurrentPage('language')}
              />
              <Separator />
              <SettingItem
                icon={Wallet}
                title="Payments & Wallet"
                subtitle="Payment methods and transaction limits"
                onClick={() => setCurrentPage('payments')}
              />
              <Separator />
              <SettingItem
                icon={Database}
                title="Data & Storage"
                subtitle="Manage app data usage and storage"
                onClick={() => setCurrentPage('data')}
              />
            </CardContent>
          </Card>

          {/* Appearance */}
          <Card>
            <CardHeader>
              <CardTitle>Appearance</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label>Theme</Label>
                  <div className="grid grid-cols-3 gap-2">
                    <Button
                      variant={theme === 'light' ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => setTheme('light')}
                    >
                      <Sun className="w-4 h-4 mr-2" />
                      Light
                    </Button>
                    <Button
                      variant={theme === 'dark' ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => setTheme('dark')}
                    >
                      <Moon className="w-4 h-4 mr-2" />
                      Dark
                    </Button>
                    <Button
                      variant={theme === 'system' ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => setTheme('system')}
                    >
                      <Monitor className="w-4 h-4 mr-2" />
                      Auto
                    </Button>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Support */}
          <Card>
            <CardContent className="p-0">
              <SettingItem
                icon={HelpCircle}
                title="Help & Support"
                subtitle="FAQ, contact support, and app info"
                onClick={() => setCurrentPage('help')}
              />
            </CardContent>
          </Card>

          {/* Logout */}
          <Card>
            <CardContent className="p-0">
              <Button
                variant="ghost"
                className="w-full justify-start h-auto p-4 text-destructive hover:text-destructive"
                onClick={onLogout}
              >
                <div className="flex items-center gap-3 w-full">
                  <LogOut className="w-5 h-5" />
                  <span className="font-medium">Log Out</span>
                </div>
              </Button>
            </CardContent>
          </Card>
        </div>
      </ScrollArea>
    </div>
  );
}