import React, { useState, useMemo, useEffect } from 'react';
import { Share, Copy, MessageCircle, Users, Facebook, Twitter, Instagram, Link, QrCode, Mail, Download, Star, Heart, Bookmark, Send, CheckCircle, Video, ShoppingCart, Radio, Store, FileText, Phone, MessageSquare } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { ScrollArea } from './ui/scroll-area';
import { Separator } from './ui/separator';
import { toast } from "sonner";
import { useGetUserFriendsQuery } from '../services/microservicesApi';

interface ShareDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  content: {
    type: 'post' | 'video' | 'product' | 'live-stream' | 'shop' | 'workspace-file';
    id: string;
    title: string;
    description?: string;
    image?: string;
    author?: {
      name: string;
      avatar?: string;
    };
    url?: string;
    price?: {
      amount: number;
      currency: string;
    };
    metadata?: any;
  };
  user: any;
}

interface Contact {
  id: string;
  name: string;
  avatar?: string;
  status: 'online' | 'offline';
  lastSeen?: string;
  type: 'individual' | 'group';
  isFrequent?: boolean;
}

export function ShareDialog({ open, onOpenChange, content, user }: ShareDialogProps) {
  const [copiedLink, setCopiedLink] = useState(false);
  const [selectedContacts, setSelectedContacts] = useState<string[]>([]);
  const [shareMessage, setShareMessage] = useState('');
  const [currentTab, setCurrentTab] = useState<'contacts' | 'social' | 'link'>('contacts');

  // RTK Query for user friends/contacts
  const {
    data: friendsData,
    isLoading: friendsLoading,
    error: friendsError
  } = useGetUserFriendsQuery({
    status: 'all',
    limit: 50
  });

  const contacts: Contact[] = useMemo(() => {
    if (friendsLoading || !friendsData) {
      // Fallback data while loading
      return [
        {
          id: '1',
          name: 'Family Group',
          avatar: '',
          status: 'online',
          type: 'group',
          isFrequent: true
        },
        {
          id: '2',
          name: 'Mom',
          avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=40&h=40&fit=crop',
          status: 'online',
          type: 'individual',
          isFrequent: true
        }
      ];
    }

    return friendsData.map((friend: any) => ({
      id: friend.id || friend.user_id || friend.friend_id || `contact-${Math.random()}`,
      name: friend.name || friend.display_name || friend.username || 'Unknown Contact',
      avatar: friend.avatar || friend.profile_picture || friend.image || undefined,
      status: friend.status || friend.online_status || (friend.is_online ? 'online' : 'offline'),
      lastSeen: friend.lastSeen || friend.last_seen || friend.last_active || undefined,
      type: friend.type || friend.chat_type || (friend.is_group ? 'group' : 'individual'),
      isFrequent: friend.isFrequent || friend.is_frequent || friend.interaction_score > 50 || false
    }));
  }, [friendsData, friendsLoading]);

  const frequentContacts = contacts.filter(c => c.isFrequent);
  const allContacts = contacts.filter(c => !c.isFrequent);

  const generateShareUrl = () => {
    const baseUrl = 'https://telegram-sea.app';
    switch (content.type) {
      case 'post':
        return `${baseUrl}/post/${content.id}`;
      case 'video':
        return `${baseUrl}/video/${content.id}`;
      case 'product':
        return `${baseUrl}/product/${content.id}`;
      case 'live-stream':
        return `${baseUrl}/live/${content.id}`;
      case 'shop':
        return `${baseUrl}/shop/${content.id}`;
      case 'workspace-file':
        return `${baseUrl}/workspace/file/${content.id}`;
      default:
        return `${baseUrl}/${content.type}/${content.id}`;
    }
  };

  const getShareText = () => {
    const url = generateShareUrl();
    let text = '';
    
    switch (content.type) {
      case 'post':
        text = `Check out this post: "${content.title}"`;
        break;
      case 'video':
        text = `Watch this video: "${content.title}"`;
        break;
      case 'product':
        text = `Check out this product: "${content.title}"${content.price ? ` - ${content.price.currency} ${content.price.amount}` : ''}`;
        break;
      case 'live-stream':
        text = `Join this live stream: "${content.title}"`;
        break;
      case 'shop':
        text = `Visit this shop: "${content.title}"`;
        break;
      case 'workspace-file':
        text = `Check out this file: "${content.title}"`;
        break;
      default:
        text = `Check this out: "${content.title}"`;
    }
    
    return `${text}\n\n${url}`;
  };

  const handleCopyLink = async () => {
    const shareText = getShareText();
    try {
      await navigator.clipboard.writeText(shareText);
      setCopiedLink(true);
      toast.success('Link copied to clipboard!');
      setTimeout(() => setCopiedLink(false), 2000);
    } catch (err) {
      toast.error('Failed to copy link');
    }
  };

  const handleContactToggle = (contactId: string) => {
    setSelectedContacts(prev => 
      prev.includes(contactId) 
        ? prev.filter(id => id !== contactId)
        : [...prev, contactId]
    );
  };

  const handleSendToContacts = () => {
    if (selectedContacts.length === 0) {
      toast.error('Please select at least one contact');
      return;
    }

    const contactNames = contacts
      .filter(c => selectedContacts.includes(c.id))
      .map(c => c.name)
      .join(', ');

    toast.success(`Shared with ${contactNames}!`);
    setSelectedContacts([]);
    setShareMessage('');
    onOpenChange(false);
  };

  const handleSocialShare = (platform: string) => {
    const shareText = getShareText();
    const url = generateShareUrl();
    
    let shareUrl = '';
    switch (platform) {
      case 'whatsapp':
        shareUrl = `https://wa.me/?text=${encodeURIComponent(shareText)}`;
        break;
      case 'telegram':
        shareUrl = `https://t.me/share/url?url=${encodeURIComponent(url)}&text=${encodeURIComponent(content.title)}`;
        break;
      case 'facebook':
        shareUrl = `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(url)}`;
        break;
      case 'twitter':
        shareUrl = `https://twitter.com/intent/tweet?text=${encodeURIComponent(shareText)}`;
        break;
      case 'instagram':
        // Instagram doesn't support direct URL sharing, so copy to clipboard
        handleCopyLink();
        toast.success('Link copied! Paste it in your Instagram story or bio');
        return;
      case 'email':
        shareUrl = `mailto:?subject=${encodeURIComponent(content.title)}&body=${encodeURIComponent(shareText)}`;
        break;
      default:
        return;
    }
    
    window.open(shareUrl, '_blank');
    toast.success(`Opening ${platform} to share!`);
  };

  const getContentIcon = () => {
    switch (content.type) {
      case 'post':
        return <MessageCircle className="w-5 h-5 text-chart-1" />;
      case 'video':
        return <Video className="w-5 h-5 text-chart-2" />;
      case 'product':
        return <ShoppingCart className="w-5 h-5 text-chart-3" />;
      case 'live-stream':
        return <Radio className="w-5 h-5 text-red-500" />;
      case 'shop':
        return <Store className="w-5 h-5 text-chart-4" />;
      case 'workspace-file':
        return <FileText className="w-5 h-5 text-chart-5" />;
      default:
        return <Share className="w-5 h-5 text-muted-foreground" />;
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md max-h-[80vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Share className="w-5 h-5" />
            Share {content.type.charAt(0).toUpperCase() + content.type.slice(1).replace('-', ' ')}
          </DialogTitle>
        </DialogHeader>

        {/* Content Preview */}
        <Card className="mb-4">
          <CardContent className="p-3">
            <div className="flex items-start gap-3">
              {content.image ? (
                <img src={content.image} alt={content.title} className="w-12 h-12 rounded object-cover" />
              ) : (
                <div className="w-12 h-12 bg-muted rounded flex items-center justify-center">
                  {getContentIcon()}
                </div>
              )}
              <div className="flex-1 min-w-0">
                <h4 className="font-medium text-sm line-clamp-2">{content.title}</h4>
                {content.description && (
                  <p className="text-xs text-muted-foreground line-clamp-2 mt-1">{content.description}</p>
                )}
                {content.author && (
                  <div className="flex items-center gap-2 mt-2">
                    <Avatar className="w-4 h-4">
                      <AvatarImage src={content.author.avatar} />
                      <AvatarFallback className="text-xs">{content.author.name.charAt(0)}</AvatarFallback>
                    </Avatar>
                    <span className="text-xs text-muted-foreground">by {content.author.name}</span>
                  </div>
                )}
                {content.price && (
                  <Badge variant="secondary" className="text-xs mt-2">
                    {content.price.currency} {content.price.amount.toLocaleString()}
                  </Badge>
                )}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Tab Navigation */}
        <div className="flex gap-1 bg-muted p-1 rounded-lg mb-4">
          <Button
            variant={currentTab === 'contacts' ? 'default' : 'ghost'}
            size="sm"
            className="flex-1"
            onClick={() => setCurrentTab('contacts')}
          >
            <MessageCircle className="w-4 h-4 mr-2" />
            Contacts
          </Button>
          <Button
            variant={currentTab === 'social' ? 'default' : 'ghost'}
            size="sm"
            className="flex-1"
            onClick={() => setCurrentTab('social')}
          >
            <Users className="w-4 h-4 mr-2" />
            Social
          </Button>
          <Button
            variant={currentTab === 'link' ? 'default' : 'ghost'}
            size="sm"
            className="flex-1"
            onClick={() => setCurrentTab('link')}
          >
            <Link className="w-4 h-4 mr-2" />
            Link
          </Button>
        </div>

        <div className="flex-1 overflow-hidden">
          {currentTab === 'contacts' && (
            <div className="space-y-4">
              {/* Message Input */}
              <div>
                <label className="text-sm font-medium mb-2 block">Add a message (optional)</label>
                <Input
                  placeholder="Say something about this..."
                  value={shareMessage}
                  onChange={(e) => setShareMessage(e.target.value)}
                  className="text-sm"
                />
              </div>

              {/* Frequent Contacts */}
              {frequentContacts.length > 0 && (
                <div>
                  <h4 className="text-sm font-medium mb-2">Frequent</h4>
                  <div className="space-y-2">
                    {frequentContacts.map((contact) => (
                      <div
                        key={contact.id}
                        className={`flex items-center gap-3 p-2 rounded-lg cursor-pointer transition-colors ${
                          selectedContacts.includes(contact.id) 
                            ? 'bg-primary/10 border border-primary/20' 
                            : 'hover:bg-muted'
                        }`}
                        onClick={() => handleContactToggle(contact.id)}
                      >
                        <div className="relative">
                          <Avatar className="w-10 h-10">
                            {contact.avatar ? (
                              <AvatarImage src={contact.avatar} />
                            ) : (
                              <div className="w-full h-full bg-chart-2 rounded-full flex items-center justify-center">
                                {contact.type === 'group' ? (
                                  <Users className="w-5 h-5 text-white" />
                                ) : (
                                  <span className="text-white font-medium">{contact.name.charAt(0)}</span>
                                )}
                              </div>
                            )}
                            <AvatarFallback>{contact.name.charAt(0)}</AvatarFallback>
                          </Avatar>
                          {contact.status === 'online' && contact.type === 'individual' && (
                            <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-green-500 border-2 border-white rounded-full"></div>
                          )}
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="font-medium text-sm">{contact.name}</p>
                          <div className="flex items-center gap-2">
                            {contact.type === 'group' && <Badge variant="outline" className="text-xs">Group</Badge>}
                            {contact.status === 'offline' && contact.lastSeen && (
                              <span className="text-xs text-muted-foreground">Last seen {contact.lastSeen}</span>
                            )}
                            {contact.status === 'online' && (
                              <span className="text-xs text-green-600">Online</span>
                            )}
                          </div>
                        </div>
                        {selectedContacts.includes(contact.id) && (
                          <CheckCircle className="w-5 h-5 text-primary" />
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* All Contacts */}
              <div>
                <h4 className="text-sm font-medium mb-2">All Contacts</h4>
                <ScrollArea className="h-40">
                  <div className="space-y-2">
                    {allContacts.map((contact) => (
                      <div
                        key={contact.id}
                        className={`flex items-center gap-3 p-2 rounded-lg cursor-pointer transition-colors ${
                          selectedContacts.includes(contact.id) 
                            ? 'bg-primary/10 border border-primary/20' 
                            : 'hover:bg-muted'
                        }`}
                        onClick={() => handleContactToggle(contact.id)}
                      >
                        <div className="relative">
                          <Avatar className="w-8 h-8">
                            {contact.avatar ? (
                              <AvatarImage src={contact.avatar} />
                            ) : (
                              <div className="w-full h-full bg-chart-3 rounded-full flex items-center justify-center">
                                {contact.type === 'group' ? (
                                  <Users className="w-4 h-4 text-white" />
                                ) : (
                                  <span className="text-white text-sm">{contact.name.charAt(0)}</span>
                                )}
                              </div>
                            )}
                            <AvatarFallback className="text-xs">{contact.name.charAt(0)}</AvatarFallback>
                          </Avatar>
                          {contact.status === 'online' && contact.type === 'individual' && (
                            <div className="absolute -bottom-1 -right-1 w-2 h-2 bg-green-500 border border-white rounded-full"></div>
                          )}
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm">{contact.name}</p>
                          {contact.status === 'offline' && contact.lastSeen && (
                            <span className="text-xs text-muted-foreground">Last seen {contact.lastSeen}</span>
                          )}
                        </div>
                        {selectedContacts.includes(contact.id) && (
                          <CheckCircle className="w-4 h-4 text-primary" />
                        )}
                      </div>
                    ))}
                  </div>
                </ScrollArea>
              </div>

              {/* Send Button */}
              <Button 
                onClick={handleSendToContacts} 
                disabled={selectedContacts.length === 0}
                className="w-full"
              >
                <Send className="w-4 h-4 mr-2" />
                Send to {selectedContacts.length} contact{selectedContacts.length !== 1 ? 's' : ''}
              </Button>
            </div>
          )}

          {currentTab === 'social' && (
            <div className="space-y-4">
              <div className="text-center">
                <h4 className="text-lg font-semibold bg-gradient-to-r from-pink-500 to-purple-500 bg-clip-text text-transparent mb-1">
                  Share the vibe ✨
                </h4>
                <p className="text-sm text-gray-500">Let your friends discover this amazing content</p>
              </div>

              {/* Featured Platforms - Instagram/TikTok style grid */}
              <div className="grid grid-cols-3 gap-4">
                <div
                  className="flex flex-col items-center p-4 rounded-2xl bg-gradient-to-b from-green-50 to-green-100 border border-green-200 cursor-pointer hover:scale-105 transition-all duration-200"
                  onClick={() => handleSocialShare('whatsapp')}
                >
                  <div className="w-14 h-14 bg-gradient-to-r from-green-400 to-green-500 rounded-2xl flex items-center justify-center mb-2 shadow-lg">
                    <span className="text-2xl">💬</span>
                  </div>
                  <span className="text-sm font-medium text-green-700">WhatsApp</span>
                </div>

                <div
                  className="flex flex-col items-center p-4 rounded-2xl bg-gradient-to-b from-pink-50 to-purple-100 border border-pink-200 cursor-pointer hover:scale-105 transition-all duration-200"
                  onClick={() => handleSocialShare('instagram')}
                >
                  <div className="w-14 h-14 bg-gradient-to-r from-purple-500 via-pink-500 to-orange-400 rounded-2xl flex items-center justify-center mb-2 shadow-lg">
                    <span className="text-2xl">📷</span>
                  </div>
                  <span className="text-sm font-medium text-pink-700">Instagram</span>
                </div>

                <div
                  className="flex flex-col items-center p-4 rounded-2xl bg-gradient-to-b from-blue-50 to-blue-100 border border-blue-200 cursor-pointer hover:scale-105 transition-all duration-200"
                  onClick={() => handleSocialShare('facebook')}
                >
                  <div className="w-14 h-14 bg-gradient-to-r from-blue-500 to-blue-600 rounded-2xl flex items-center justify-center mb-2 shadow-lg">
                    <span className="text-2xl">📘</span>
                  </div>
                  <span className="text-sm font-medium text-blue-700">Facebook</span>
                </div>

                <div
                  className="flex flex-col items-center p-4 rounded-2xl bg-gradient-to-b from-gray-50 to-gray-100 border border-gray-200 cursor-pointer hover:scale-105 transition-all duration-200"
                  onClick={() => handleSocialShare('twitter')}
                >
                  <div className="w-14 h-14 bg-gradient-to-r from-gray-800 to-black rounded-2xl flex items-center justify-center mb-2 shadow-lg">
                    <span className="text-2xl">🐦</span>
                  </div>
                  <span className="text-sm font-medium text-gray-700">Twitter</span>
                </div>

                <div
                  className="flex flex-col items-center p-4 rounded-2xl bg-gradient-to-b from-cyan-50 to-blue-100 border border-cyan-200 cursor-pointer hover:scale-105 transition-all duration-200"
                  onClick={() => handleSocialShare('telegram')}
                >
                  <div className="w-14 h-14 bg-gradient-to-r from-cyan-400 to-blue-500 rounded-2xl flex items-center justify-center mb-2 shadow-lg">
                    <span className="text-2xl">✈️</span>
                  </div>
                  <span className="text-sm font-medium text-cyan-700">Telegram</span>
                </div>

                <div
                  className="flex flex-col items-center p-4 rounded-2xl bg-gradient-to-b from-red-50 to-red-100 border border-red-200 cursor-pointer hover:scale-105 transition-all duration-200"
                  onClick={() => handleSocialShare('tiktok')}
                >
                  <div className="w-14 h-14 bg-gradient-to-r from-black via-red-500 to-cyan-400 rounded-2xl flex items-center justify-center mb-2 shadow-lg">
                    <span className="text-2xl">🎵</span>
                  </div>
                  <span className="text-sm font-medium text-red-700">TikTok</span>
                </div>
              </div>

              {/* More Options */}
              <div className="pt-4 border-t border-gray-100">
                <div className="grid grid-cols-2 gap-3">
                  <Button
                    variant="outline"
                    className="flex items-center gap-2 h-11 rounded-xl hover:bg-gray-50"
                    onClick={() => handleSocialShare('line')}
                  >
                    <div className="w-6 h-6 bg-green-500 rounded-lg flex items-center justify-center">
                      <span className="text-white text-xs font-bold">L</span>
                    </div>
                    <span className="font-medium">LINE</span>
                  </Button>

                  <Button
                    variant="outline"
                    className="flex items-center gap-2 h-11 rounded-xl hover:bg-gray-50"
                    onClick={() => handleSocialShare('email')}
                  >
                    <div className="w-6 h-6 bg-gray-500 rounded-lg flex items-center justify-center">
                      <Mail className="w-4 h-4 text-white" />
                    </div>
                    <span className="font-medium">Email</span>
                  </Button>

                  <Button
                    variant="outline"
                    className="flex items-center gap-2 h-11 rounded-xl hover:bg-gray-50 col-span-2"
                    onClick={() => handleSocialShare('more')}
                  >
                    <div className="w-6 h-6 bg-gradient-to-r from-purple-400 to-pink-400 rounded-lg flex items-center justify-center">
                      <span className="text-white text-sm">⋯</span>
                    </div>
                    <span className="font-medium">More options</span>
                  </Button>
                </div>
              </div>

              {/* Fun Call to Action */}
              <div className="text-center pt-2">
                <p className="text-xs text-gray-400">
                  Spread the love! 💕 Your friends will thank you later
                </p>
              </div>
            </div>
          )}

          {currentTab === 'link' && (
            <div className="space-y-4">
              <div>
                <label className="text-sm font-medium mb-2 block">Share Link</label>
                <div className="flex gap-2">
                  <Input
                    readOnly
                    value={generateShareUrl()}
                    className="text-sm"
                  />
                  <Button 
                    variant="outline" 
                    size="icon"
                    onClick={handleCopyLink}
                    className={copiedLink ? 'bg-green-100 text-green-700' : ''}
                  >
                    {copiedLink ? <CheckCircle className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
                  </Button>
                </div>
              </div>

              <Separator />

              <div>
                <h4 className="text-sm font-medium mb-3">Quick Actions</h4>
                <div className="grid grid-cols-2 gap-2">
                  <Button variant="outline" size="sm" onClick={handleCopyLink}>
                    <Copy className="w-4 h-4 mr-2" />
                    Copy Link
                  </Button>
                  <Button variant="outline" size="sm" onClick={() => toast.success('QR code generation coming soon!')}>
                    <QrCode className="w-4 h-4 mr-2" />
                    QR Code
                  </Button>
                  <Button variant="outline" size="sm" onClick={() => toast.success('Download feature coming soon!')}>
                    <Download className="w-4 h-4 mr-2" />
                    Download
                  </Button>
                  <Button variant="outline" size="sm" onClick={() => toast.success('Added to bookmarks!')}>
                    <Bookmark className="w-4 h-4 mr-2" />
                    Bookmark
                  </Button>
                </div>
              </div>

              <Card className="bg-muted/50">
                <CardContent className="p-3">
                  <div className="flex items-start gap-2">
                    <Star className="w-4 h-4 text-yellow-500 mt-0.5" />
                    <div>
                      <p className="text-sm font-medium">Share and Earn</p>
                      <p className="text-xs text-muted-foreground">
                        Get rewards when friends engage with your shared content!
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}