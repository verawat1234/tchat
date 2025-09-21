import React, { useState } from 'react';
import { 
  ArrowLeft, 
  Search, 
  Users, 
  UserPlus, 
  MessageSquare, 
  Hash, 
  Bot,
  Star,
  Phone,
  MapPin,
  QrCode,
  Plus
} from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { Separator } from './ui/separator';

interface NewChatScreenProps {
  user: any;
  onBack: () => void;
  onCreateChat: (chatData: any) => void;
}

interface Contact {
  id: string;
  name: string;
  username?: string;
  avatar?: string;
  status: 'online' | 'offline' | 'recently' | 'lastweek';
  isVerified?: boolean;
  isMutual?: boolean;
  phone?: string;
  location?: string;
}

type ChatType = 'all' | 'contacts' | 'groups' | 'channels' | 'bots';

export function NewChatScreen({ user, onBack, onCreateChat }: NewChatScreenProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedType, setSelectedType] = useState<ChatType>('all');

  // Mock contacts data
  const contacts: Contact[] = [
    {
      id: 'contact1',
      name: 'Mom',
      username: 'mom_thailand',
      avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face',
      status: 'online',
      isVerified: true,
      phone: '+66 89 123 4567',
      location: 'Bangkok, Thailand'
    },
    {
      id: 'contact2',
      name: 'Dad',
      username: 'dad_bangkok',
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
      status: 'recently',
      phone: '+66 89 765 4321'
    },
    {
      id: 'contact3',
      name: 'Sister',
      username: 'sis_cute',
      avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop&crop=face',
      status: 'online',
      isMutual: true
    },
    {
      id: 'contact4',
      name: 'Best Friend',
      username: 'bestie_forever',
      avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face',
      status: 'offline',
      isMutual: true
    },
    {
      id: 'contact5',
      name: 'Work Colleague',
      username: 'work_mate',
      status: 'lastweek'
    }
  ];

  const groups = [
    {
      id: 'group1',
      name: 'Family Group',
      members: 5,
      type: 'group',
      description: 'Our lovely family chat'
    },
    {
      id: 'group2',
      name: 'Work Team',
      members: 12,
      type: 'group',
      description: 'Daily work discussions'
    },
    {
      id: 'group3',
      name: 'Thai Food Lovers',
      members: 156,
      type: 'group',
      description: 'Share your favorite Thai recipes!'
    }
  ];

  const channels = [
    {
      id: 'channel1',
      name: 'Thailand News',
      subscribers: 15420,
      type: 'channel',
      description: 'Latest news from Thailand',
      isVerified: true
    },
    {
      id: 'channel2',
      name: 'Bangkok Traffic',
      subscribers: 8934,
      type: 'channel',
      description: 'Real-time traffic updates'
    },
    {
      id: 'channel3',
      name: 'Street Food Guide',
      subscribers: 23156,
      type: 'channel',
      description: 'Best street food spots in SEA'
    }
  ];

  const bots = [
    {
      id: 'bot1',
      name: 'PromptPay Bot',
      username: 'promptpay_bot',
      type: 'bot',
      description: 'Send money via PromptPay easily',
      isVerified: true
    },
    {
      id: 'bot2',
      name: 'Weather Bot',
      username: 'weather_th_bot',
      type: 'bot',
      description: 'Get weather updates for Thailand'
    },
    {
      id: 'bot3',
      name: 'Translation Bot',
      username: 'translate_sea_bot',
      type: 'bot',
      description: 'Translate between Thai, English, and more'
    }
  ];

  const filteredContacts = contacts.filter(contact =>
    contact.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    contact.username?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const filteredGroups = groups.filter(group =>
    group.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const filteredChannels = channels.filter(channel =>
    channel.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const filteredBots = bots.filter(bot =>
    bot.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    bot.username?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleCreateChat = (contact: Contact) => {
    onCreateChat({
      id: contact.id,
      name: contact.name,
      type: 'user',
      avatar: contact.avatar,
      username: contact.username
    });
    onBack();
  };

  const handleCreateGroup = () => {
    // In real app, would open group creation flow
    onCreateChat({
      id: 'new-group',
      name: 'New Group',
      type: 'group',
      members: 1
    });
    onBack();
  };

  const handleCreateChannel = () => {
    // In real app, would open channel creation flow
    onCreateChat({
      id: 'new-channel',
      name: 'New Channel',
      type: 'channel',
      subscribers: 0
    });
    onBack();
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online':
        return 'text-green-500';
      case 'recently':
        return 'text-yellow-500';
      default:
        return 'text-muted-foreground';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'online':
        return 'online';
      case 'recently':
        return 'last seen recently';
      case 'lastweek':
        return 'last seen within a week';
      default:
        return 'offline';
    }
  };

  return (
    <div className="h-screen bg-background flex flex-col">
      {/* Header */}
      <div className="border-b border-border p-4 flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={onBack}>
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h1 className="text-lg font-medium">New Chat</h1>
      </div>

      {/* Search */}
      <div className="p-4 border-b border-border">
        <div className="relative">
          <Search className="absolute left-3 top-3 w-4 h-4 text-muted-foreground" />
          <Input
            placeholder="Search for people and groups"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
      </div>

      {/* Quick Actions */}
      <div className="p-4 space-y-2">
        <Button 
          variant="ghost" 
          className="w-full justify-start gap-3 h-12"
          onClick={handleCreateGroup}
        >
          <div className="w-10 h-10 bg-chart-1 rounded-full flex items-center justify-center">
            <Users className="w-5 h-5 text-white" />
          </div>
          <div className="text-left">
            <p className="font-medium">New Group</p>
            <p className="text-sm text-muted-foreground">Create a group chat</p>
          </div>
        </Button>

        <Button 
          variant="ghost" 
          className="w-full justify-start gap-3 h-12"
          onClick={handleCreateChannel}
        >
          <div className="w-10 h-10 bg-chart-2 rounded-full flex items-center justify-center">
            <Hash className="w-5 h-5 text-white" />
          </div>
          <div className="text-left">
            <p className="font-medium">New Channel</p>
            <p className="text-sm text-muted-foreground">Broadcast to many people</p>
          </div>
        </Button>

        <Button variant="ghost" className="w-full justify-start gap-3 h-12">
          <div className="w-10 h-10 bg-chart-3 rounded-full flex items-center justify-center">
            <QrCode className="w-5 h-5 text-white" />
          </div>
          <div className="text-left">
            <p className="font-medium">Scan QR Code</p>
            <p className="text-sm text-muted-foreground">Add contact via QR</p>
          </div>
        </Button>

        <Button variant="ghost" className="w-full justify-start gap-3 h-12">
          <div className="w-10 h-10 bg-chart-4 rounded-full flex items-center justify-center">
            <UserPlus className="w-5 h-5 text-white" />
          </div>
          <div className="text-left">
            <p className="font-medium">Invite Friends</p>
            <p className="text-sm text-muted-foreground">Share invite link</p>
          </div>
        </Button>
      </div>

      <Separator />

      {/* Filter Tabs */}
      <div className="flex border-b border-border">
        {[
          { key: 'all', label: 'All', icon: MessageSquare },
          { key: 'contacts', label: 'Contacts', icon: UserPlus },
          { key: 'groups', label: 'Groups', icon: Users },
          { key: 'channels', label: 'Channels', icon: Hash },
          { key: 'bots', label: 'Bots', icon: Bot }
        ].map(({ key, label, icon: Icon }) => (
          <Button
            key={key}
            variant={selectedType === key ? 'default' : 'ghost'}
            size="sm"
            className="flex-1 rounded-none"
            onClick={() => setSelectedType(key as ChatType)}
          >
            <Icon className="w-4 h-4 mr-1" />
            {label}
          </Button>
        ))}
      </div>

      {/* Results */}
      <ScrollArea className="flex-1">
        <div className="p-4 space-y-4">
          {/* Contacts */}
          {(selectedType === 'all' || selectedType === 'contacts') && (
            <div>
              {selectedType === 'all' && filteredContacts.length > 0 && (
                <h3 className="text-sm font-medium text-muted-foreground mb-3">Contacts</h3>
              )}
              <div className="space-y-2">
                {filteredContacts.map((contact) => (
                  <Button
                    key={contact.id}
                    variant="ghost"
                    className="w-full justify-start gap-3 h-16 p-3"
                    onClick={() => handleCreateChat(contact)}
                  >
                    <div className="relative">
                      <Avatar className="w-12 h-12">
                        <AvatarImage src={contact.avatar} />
                        <AvatarFallback>{contact.name.charAt(0)}</AvatarFallback>
                      </Avatar>
                      {contact.status === 'online' && (
                        <div className="absolute -bottom-1 -right-1 w-4 h-4 bg-green-500 rounded-full border-2 border-background"></div>
                      )}
                    </div>
                    
                    <div className="flex-1 text-left">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">{contact.name}</p>
                        {contact.isVerified && (
                          <Star className="w-4 h-4 text-chart-1 fill-chart-1" />
                        )}
                        {contact.isMutual && (
                          <Badge variant="secondary" className="text-xs">Mutual</Badge>
                        )}
                      </div>
                      <div className="flex items-center gap-1">
                        {contact.username && (
                          <p className="text-sm text-muted-foreground">@{contact.username}</p>
                        )}
                        <span className="text-sm text-muted-foreground">•</span>
                        <p className={`text-sm ${getStatusColor(contact.status)}`}>
                          {getStatusText(contact.status)}
                        </p>
                      </div>
                      {contact.location && (
                        <div className="flex items-center gap-1 mt-1">
                          <MapPin className="w-3 h-3 text-muted-foreground" />
                          <p className="text-xs text-muted-foreground">{contact.location}</p>
                        </div>
                      )}
                    </div>

                    {contact.phone && (
                      <Button size="icon" variant="ghost" className="w-8 h-8">
                        <Phone className="w-4 h-4" />
                      </Button>
                    )}
                  </Button>
                ))}
              </div>
            </div>
          )}

          {/* Groups */}
          {(selectedType === 'all' || selectedType === 'groups') && filteredGroups.length > 0 && (
            <div>
              {selectedType === 'all' && (
                <h3 className="text-sm font-medium text-muted-foreground mb-3">Groups</h3>
              )}
              <div className="space-y-2">
                {filteredGroups.map((group) => (
                  <Button
                    key={group.id}
                    variant="ghost"
                    className="w-full justify-start gap-3 h-16 p-3"
                    onClick={() => onCreateChat(group)}
                  >
                    <Avatar className="w-12 h-12">
                      <AvatarFallback>
                        <Users className="w-6 h-6" />
                      </AvatarFallback>
                    </Avatar>
                    
                    <div className="flex-1 text-left">
                      <p className="font-medium">{group.name}</p>
                      <p className="text-sm text-muted-foreground">
                        {group.members} members
                      </p>
                      {group.description && (
                        <p className="text-xs text-muted-foreground truncate">
                          {group.description}
                        </p>
                      )}
                    </div>
                  </Button>
                ))}
              </div>
            </div>
          )}

          {/* Channels */}
          {(selectedType === 'all' || selectedType === 'channels') && filteredChannels.length > 0 && (
            <div>
              {selectedType === 'all' && (
                <h3 className="text-sm font-medium text-muted-foreground mb-3">Channels</h3>
              )}
              <div className="space-y-2">
                {filteredChannels.map((channel) => (
                  <Button
                    key={channel.id}
                    variant="ghost"
                    className="w-full justify-start gap-3 h-16 p-3"
                    onClick={() => onCreateChat(channel)}
                  >
                    <Avatar className="w-12 h-12">
                      <AvatarFallback>
                        <Hash className="w-6 h-6" />
                      </AvatarFallback>
                    </Avatar>
                    
                    <div className="flex-1 text-left">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">{channel.name}</p>
                        {channel.isVerified && (
                          <Star className="w-4 h-4 text-chart-1 fill-chart-1" />
                        )}
                      </div>
                      <p className="text-sm text-muted-foreground">
                        {channel.subscribers.toLocaleString()} subscribers
                      </p>
                      {channel.description && (
                        <p className="text-xs text-muted-foreground truncate">
                          {channel.description}
                        </p>
                      )}
                    </div>
                  </Button>
                ))}
              </div>
            </div>
          )}

          {/* Bots */}
          {(selectedType === 'all' || selectedType === 'bots') && filteredBots.length > 0 && (
            <div>
              {selectedType === 'all' && (
                <h3 className="text-sm font-medium text-muted-foreground mb-3">Bots</h3>
              )}
              <div className="space-y-2">
                {filteredBots.map((bot) => (
                  <Button
                    key={bot.id}
                    variant="ghost"
                    className="w-full justify-start gap-3 h-16 p-3"
                    onClick={() => onCreateChat(bot)}
                  >
                    <Avatar className="w-12 h-12">
                      <AvatarFallback>
                        <Bot className="w-6 h-6" />
                      </AvatarFallback>
                    </Avatar>
                    
                    <div className="flex-1 text-left">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">{bot.name}</p>
                        {bot.isVerified && (
                          <Star className="w-4 h-4 text-chart-1 fill-chart-1" />
                        )}
                        <Badge variant="secondary" className="text-xs">BOT</Badge>
                      </div>
                      <p className="text-sm text-muted-foreground">@{bot.username}</p>
                      {bot.description && (
                        <p className="text-xs text-muted-foreground truncate">
                          {bot.description}
                        </p>
                      )}
                    </div>
                  </Button>
                ))}
              </div>
            </div>
          )}

          {/* No Results */}
          {searchQuery && 
           filteredContacts.length === 0 && 
           filteredGroups.length === 0 && 
           filteredChannels.length === 0 && 
           filteredBots.length === 0 && (
            <div className="text-center py-8">
              <MessageSquare className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="font-medium mb-2">No results found</h3>
              <p className="text-sm text-muted-foreground">
                Try searching with a different term or create a new group
              </p>
            </div>
          )}
        </div>
      </ScrollArea>
    </div>
  );
}