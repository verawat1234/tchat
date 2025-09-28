import React, { useState, useRef, useEffect, useCallback, useMemo } from 'react';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { ScrollArea } from './ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Input } from './ui/input';
import { Separator } from './ui/separator';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog';
import { Textarea } from './ui/textarea';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { 
  Play, 
  Pause, 
  Heart, 
  MessageCircle, 
  Share, 
  Bookmark, 
  UserPlus, 
  UserMinus,
  Search, 
  Filter, 
  TrendingUp, 
  Clock, 
  Eye, 
  Volume2, 
  VolumeX, 
  Maximize, 
  MoreVertical,
  ArrowLeft,
  ArrowUp,
  ArrowDown,
  Home,
  Compass,
  Flame,
  Music,
  Gamepad2,
  GraduationCap,
  ChevronDown,
  ChevronRight,
  CheckCircle,
  Plus,
  Users,
  Bell,
  BellOff,
  Settings,
  Upload,
  Camera,
  Video,
  Zap,
  Star,
  Grid3X3,
  List,
  Calendar
} from 'lucide-react';
import { toast } from "sonner";

// RTK Query video hooks
import {
  useGetVideosQuery,
  useSearchVideosQuery,
  useGetVideoByIdQuery,
  useGetChannelsQuery,
  useGetChannelByIdQuery,
  useGetChannelVideosQuery,
  useChannelSubscriptionMutation,
  useVideoInteractionMutation,
  useCreateVideoCommentMutation,
  useGetVideoCommentsQuery,
} from '../services/videoApi';

interface Video {
  id: string;
  title: string;
  description: string;
  thumbnail: string;
  videoUrl?: string;
  duration: string;
  views: number;
  likes: number;
  channel: {
    id: string;
    name: string;
    avatar: string;
    subscribers: number;
    verified: boolean;
  };
  uploadTime: string;
  isLive?: boolean;
  category: string;
  tags: string[];
  type: 'short' | 'long';
}

interface Channel {
  id: string;
  name: string;
  avatar: string;
  banner?: string;
  description: string;
  subscribers: number;
  videos: number;
  verified: boolean;
  category: string;
  joinedDate: string;
  location?: string;
  isSubscribed?: boolean;
  notificationsEnabled?: boolean;
}

interface VideoTabProps {
  user: any;
  onBack: () => void;
  onVideoPlay: (videoId: string, videoData?: any) => void;
  onVideoLike: (videoId: string) => void;
  onVideoShare: (videoId: string, videoData?: any) => void;
  onSubscribe: (channelId: string) => void;
  currentVideoId: string;
  watchedVideos: string[];
  likedVideos: string[];
  subscribedChannels?: string[]; // Optional - now derived from RTK Query
}

export function VideoTab({
  user,
  onBack,
  onVideoPlay,
  onVideoLike,
  onVideoShare,
  onSubscribe,
  currentVideoId,
  watchedVideos,
  likedVideos,
  subscribedChannels: propSubscribedChannels // Not used - RTK Query manages subscriptions
}: VideoTabProps) {
  const [currentTab, setCurrentTab] = useState<'shorts' | 'long' | 'channels' | 'subscriptions'>('shorts');
  const [currentShortIndex, setCurrentShortIndex] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  // Remove hardcoded subscribedChannels state - now derived from RTK Query channels data
  const [selectedChannel, setSelectedChannel] = useState<Channel | null>(null);
  const [isCreateChannelOpen, setIsCreateChannelOpen] = useState(false);
  const [newChannelData, setNewChannelData] = useState({
    name: '',
    description: '',
    category: 'food'
  });
  const [videoLoading, setVideoLoading] = useState(false);
  const [videoError, setVideoError] = useState(false);
  const [isMounted, setIsMounted] = useState(false);
  const videoRef = useRef<HTMLVideoElement>(null);

  // Mount guard to prevent premature API calls
  useEffect(() => {
    setIsMounted(true);
    return () => setIsMounted(false);
  }, []);

  // RTK Query hooks for video data - fixed to prevent infinite loops
  const { data: shortVideosData, isLoading: shortVideosLoading, error: shortVideosError } = useGetVideosQuery({
    type: 'short',
    limit: 50
  }, {
    skip: !isMounted
  });

  const { data: longVideosData, isLoading: longVideosLoading, error: longVideosError } = useGetVideosQuery({
    type: 'long',
    limit: 50
  }, {
    skip: !isMounted
  });

  // RTK Query hook for channels data
  const { data: channelsData, isLoading: channelsLoading, error: channelsError } = useGetChannelsQuery({
    limit: 50
  }, {
    skip: !isMounted
  });

  // Log any errors for debugging
  useEffect(() => {
    if (shortVideosError) {
      console.error('Short videos error:', shortVideosError);
    }
  }, [shortVideosError]);

  useEffect(() => {
    if (longVideosError) {
      console.error('Long videos error:', longVideosError);
    }
  }, [longVideosError]);

  useEffect(() => {
    if (channelsError) {
      console.error('Channels error:', channelsError);
    }
  }, [channelsError]);

  // Log RTK Query results for debugging
  useEffect(() => {
    if (shortVideosData) {
      console.log('RTK Query - Short videos data:', shortVideosData);
    }
    if (shortVideosError) {
      console.log('RTK Query - Short videos error:', shortVideosError);
    }
  }, [shortVideosData, shortVideosError]);

  useEffect(() => {
    if (longVideosData) {
      console.log('RTK Query - Long videos data:', longVideosData);
    }
    if (longVideosError) {
      console.log('RTK Query - Long videos error:', longVideosError);
    }
  }, [longVideosData, longVideosError]);

  // Use real channels data from RTK Query
  const effectiveChannels = channelsData?.data || [];

  // Replace hardcoded subscribedChannels with channels that have isSubscribed = true
  const subscribedChannelIds = useMemo(() => {
    return effectiveChannels
      .filter(channel => channel.isSubscribed)
      .map(channel => channel.id);
  }, [effectiveChannels]);

  // Video mutations - using the unified interaction and subscription mutations
  const [videoInteraction] = useVideoInteractionMutation();
  const [channelSubscription] = useChannelSubscriptionMutation();

  // Enhanced handlers that use RTK mutations
  const handleVideoLike = useCallback(async (videoId: string) => {
    try {
      const action = likedVideos.includes(videoId) ? 'unlike' : 'like';
      await videoInteraction({ videoId, action }).unwrap();
      // Also call the prop callback for existing functionality
      onVideoLike(videoId);
    } catch (error) {
      console.error('Error toggling video like:', error);
      toast.error('Failed to update like status');
    }
  }, [likedVideos, videoInteraction, onVideoLike]);

  const handleVideoShare = useCallback(async (videoId: string, videoData?: any) => {
    try {
      await videoInteraction({ videoId, action: 'share' }).unwrap();
      // Also call the prop callback for existing functionality
      onVideoShare(videoId, videoData);
      toast.success('Video shared successfully!');
    } catch (error) {
      console.error('Error sharing video:', error);
      toast.error('Failed to share video');
    }
  }, [videoInteraction, onVideoShare]);

  const handleChannelSubscribe = useCallback(async (channelId: string) => {
    try {
      const action = subscribedChannelIds.includes(channelId) ? 'UNSUBSCRIBE' : 'SUBSCRIBE';
      await channelSubscription({ channelId, action }).unwrap();
      // RTK Query handles optimistic updates automatically
      onSubscribe(channelId);
      toast.success(action === 'UNSUBSCRIBE' ? 'Unsubscribed successfully!' : 'Subscribed successfully!');
    } catch (error) {
      console.error('Error toggling subscription:', error);
      toast.error('Failed to update subscription');
    }
  }, [subscribedChannelIds, channelSubscription, onSubscribe]);


  const categories = [
    { id: 'all', name: 'All', icon: Home },
    { id: 'trending', name: 'Trending', icon: TrendingUp },
    { id: 'food', name: 'Food', icon: Flame },
    { id: 'music', name: 'Music', icon: Music },
    { id: 'entertainment', name: 'Entertainment', icon: Gamepad2 },
    { id: 'education', name: 'Education', icon: GraduationCap },
    { id: 'travel', name: 'Travel', icon: Compass },
  ];

  // Use real API data only - no mock data fallback
  const effectiveShortVideos = shortVideosData?.videos || [];
  const effectiveLongVideos = longVideosData?.videos || [];
  const effectiveChannels = derivedChannels;

  // Memoize filtered shorts to prevent infinite re-renders
  const filteredShorts = useMemo(() => {
    return selectedCategory === 'all'
      ? effectiveShortVideos
      : effectiveShortVideos.filter(video => video.category === selectedCategory);
  }, [selectedCategory, effectiveShortVideos]);

  const filteredLongs = selectedCategory === 'all'
    ? effectiveLongVideos
    : effectiveLongVideos.filter(video => video.category === selectedCategory);

  const filteredChannels = selectedCategory === 'all'
    ? effectiveChannels
    : effectiveChannels.filter(channel => channel.category === selectedCategory);

  const subscribedChannelsList = effectiveChannels.filter(channel =>
    subscribedChannelIds.includes(channel.id)
  );

  const togglePlayPause = () => {
    if (videoRef.current) {
      if (isPlaying) {
        videoRef.current.pause();
      } else {
        videoRef.current.play();
      }
      setIsPlaying(!isPlaying);
    }
  };

  const handleVideoMaximize = (video: Video) => {
    onVideoPlay(video.id, video);
  };

  const toggleMute = () => {
    if (videoRef.current) {
      videoRef.current.muted = !isMuted;
      setIsMuted(!isMuted);
    }
  };

  const handleShortSwipe = (direction: 'up' | 'down') => {
    if (direction === 'up' && currentShortIndex < filteredShorts.length - 1) {
      setCurrentShortIndex(prev => prev + 1);
    } else if (direction === 'down' && currentShortIndex > 0) {
      setCurrentShortIndex(prev => prev - 1);
    }
  };

  const formatViews = (views: number) => {
    if (views >= 1000000) {
      return `${(views / 1000000).toFixed(1)}M`;
    } else if (views >= 1000) {
      return `${(views / 1000).toFixed(1)}K`;
    }
    return views.toString();
  };

  const formatSubscribers = (subs: number) => {
    if (subs >= 1000000) {
      return `${(subs / 1000000).toFixed(1)}M`;
    } else if (subs >= 1000) {
      return `${(subs / 1000).toFixed(1)}K`;
    }
    return subs.toString();
  };


  const handleCreateChannel = () => {
    if (!newChannelData.name.trim()) {
      toast.error('Channel name is required');
      return;
    }
    
    // In real app, would create channel via API
    toast.success(`Channel "${newChannelData.name}" created successfully!`);
    setIsCreateChannelOpen(false);
    setNewChannelData({ name: '', description: '', category: 'food' });
  };

  const handleChannelClick = (channel: Channel) => {
    setSelectedChannel(channel);
  };

  const currentShort = filteredShorts[currentShortIndex];

  // Auto-play shorts when switching (infinite loop fixed by memoizing filteredShorts)
  useEffect(() => {
    const currentShortVideo = filteredShorts[currentShortIndex];
    if (currentShortVideo && videoRef.current && currentTab === 'shorts') {
      const video = videoRef.current;
      video.src = currentShortVideo.videoUrl || '';
      video.muted = isMuted;

      // Auto-play after a short delay
      const playTimer = setTimeout(() => {
        video.play().then(() => {
          setIsPlaying(true);
        }).catch((error) => {
          console.log('Auto-play failed:', error);
          setIsPlaying(false);
        });
      }, 500);

      return () => clearTimeout(playTimer);
    }
  }, [currentShortIndex, currentTab, isMuted, filteredShorts]);

  // Channel Detail View
  if (selectedChannel) {
    return (
      <div className="h-full flex flex-col bg-background">
        {/* Channel Header */}
        <div className="relative">
          {selectedChannel.banner && (
            <div className="h-32 md:h-48 relative">
              <ImageWithFallback
                src={selectedChannel.banner}
                alt={`${selectedChannel.name} banner`}
                className="w-full h-full object-cover"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent" />
            </div>
          )}
          
          <div className="p-4 bg-card">
            <Button 
              variant="ghost" 
              size="sm" 
              className="mb-4"
              onClick={() => setSelectedChannel(null)}
            >
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back to Channels
            </Button>
            
            <div className="flex items-start gap-4">
              <Avatar className="w-20 h-20 border-4 border-background">
                <AvatarImage src={selectedChannel.avatar} />
                <AvatarFallback>{selectedChannel.name.charAt(0)}</AvatarFallback>
              </Avatar>
              
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-2">
                  <h1 className="text-xl font-semibold">{selectedChannel.name}</h1>
                  {selectedChannel.verified && (
                    <CheckCircle className="w-5 h-5 text-blue-500 fill-current" />
                  )}
                </div>
                
                <div className="flex items-center gap-4 text-sm text-muted-foreground mb-3">
                  <span>{formatSubscribers(selectedChannel.subscribers)} subscribers</span>
                  <span>{selectedChannel.videos} videos</span>
                  <span>Joined {selectedChannel.joinedDate}</span>
                </div>
                
                <p className="text-sm mb-4 line-clamp-2">{selectedChannel.description}</p>
                
                <div className="flex items-center gap-2">
                  <Button
                    variant={selectedChannel.isSubscribed ? "outline" : "default"}
                    onClick={() => handleChannelSubscribe(selectedChannel.id)}
                    className="flex items-center gap-2"
                  >
                    {selectedChannel.isSubscribed ? (
                      <>
                        <UserMinus className="w-4 h-4" />
                        Subscribed
                      </>
                    ) : (
                      <>
                        <UserPlus className="w-4 h-4" />
                        Subscribe
                      </>
                    )}
                  </Button>
                  
                  {selectedChannel.isSubscribed && (
                    <Button variant="outline" size="icon">
                      {selectedChannel.notificationsEnabled ? (
                        <Bell className="w-4 h-4" />
                      ) : (
                        <BellOff className="w-4 h-4" />
                      )}
                    </Button>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
        
        {/* Channel Content */}
        <ScrollArea className="flex-1">
          <div className="p-4">
            <Tabs defaultValue="videos" className="w-full">
              <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="videos">Videos</TabsTrigger>
                <TabsTrigger value="shorts">Shorts</TabsTrigger>
                <TabsTrigger value="about">About</TabsTrigger>
              </TabsList>
              
              <TabsContent value="videos" className="mt-4">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {effectiveLongVideos
                    .filter(video => video.channel.id === selectedChannel.id)
                    .map((video) => (
                      <Card 
                        key={video.id} 
                        className="overflow-hidden cursor-pointer hover:shadow-lg transition-shadow"
                        onClick={() => onVideoPlay(video.id, video)}
                      >
                        <div className="aspect-video relative">
                          <ImageWithFallback
                            src={video.thumbnail}
                            alt={video.title}
                            className="w-full h-full object-cover"
                          />
                          <div className="absolute inset-0 flex items-center justify-center bg-black/20 opacity-0 hover:opacity-100 transition-opacity">
                            <div className="w-16 h-16 bg-white/90 rounded-full flex items-center justify-center">
                              <Play className="w-8 h-8 text-black ml-1" />
                            </div>
                          </div>
                          <Badge className="absolute bottom-2 right-2 bg-black/70 text-white">
                            {video.duration}
                          </Badge>
                        </div>
                        <CardContent className="p-3">
                          <h3 className="font-medium line-clamp-2 mb-2">{video.title}</h3>
                          <div className="flex items-center gap-2 text-xs text-muted-foreground">
                            <span>{formatViews(video.views)} views</span>
                            <span>•</span>
                            <span>{video.uploadTime}</span>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                </div>
              </TabsContent>
              
              <TabsContent value="shorts" className="mt-4">
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
                  {effectiveShortVideos
                    .filter(video => video.channel.id === selectedChannel.id)
                    .map((video) => (
                      <Card 
                        key={video.id} 
                        className="overflow-hidden cursor-pointer hover:shadow-lg transition-shadow aspect-[9/16]"
                        onClick={() => onVideoPlay(video.id, video)}
                      >
                        <div className="relative h-full">
                          <ImageWithFallback
                            src={video.thumbnail}
                            alt={video.title}
                            className="w-full h-full object-cover"
                          />
                          <div className="absolute inset-0 flex items-center justify-center bg-black/20 opacity-0 hover:opacity-100 transition-opacity">
                            <div className="w-12 h-12 bg-white/90 rounded-full flex items-center justify-center">
                              <Play className="w-6 h-6 text-black ml-0.5" />
                            </div>
                          </div>
                          <Badge className="absolute bottom-2 right-2 bg-black/70 text-white text-xs">
                            {video.duration}
                          </Badge>
                          <div className="absolute bottom-2 left-2 text-white text-xs">
                            <div className="font-medium line-clamp-2">{video.title}</div>
                            <div className="text-xs opacity-75">{formatViews(video.views)} views</div>
                          </div>
                        </div>
                      </Card>
                    ))}
                </div>
              </TabsContent>
              
              <TabsContent value="about" className="mt-4">
                <Card>
                  <CardContent className="p-6 space-y-4">
                    <div>
                      <h3 className="font-medium mb-2">Description</h3>
                      <p className="text-sm text-muted-foreground">{selectedChannel.description}</p>
                    </div>
                    
                    <Separator />
                    
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <h4 className="font-medium mb-1">Subscribers</h4>
                        <p className="text-sm text-muted-foreground">{formatSubscribers(selectedChannel.subscribers)}</p>
                      </div>
                      <div>
                        <h4 className="font-medium mb-1">Videos</h4>
                        <p className="text-sm text-muted-foreground">{selectedChannel.videos}</p>
                      </div>
                      <div>
                        <h4 className="font-medium mb-1">Joined</h4>
                        <p className="text-sm text-muted-foreground">{selectedChannel.joinedDate}</p>
                      </div>
                      <div>
                        <h4 className="font-medium mb-1">Location</h4>
                        <p className="text-sm text-muted-foreground">{selectedChannel.location}</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </ScrollArea>
      </div>
    );
  }

  return (
    <div className="w-full flex flex-col bg-background min-h-[100vh]">
      <Tabs value={currentTab} onValueChange={(value) => setCurrentTab(value as 'shorts' | 'long' | 'channels' | 'subscriptions')} className="flex-1 flex flex-col">
        {/* Tab Navigation */}
        <div className="border-b bg-card">
          <div className="space-y-3 sm:space-y-4">
            {/* Tab Navigation Section */}
            <div className="px-2 sm:px-3 pt-3 sm:pt-4">
              <div className="flex items-center justify-center gap-1 sm:gap-2 w-full">
                <TabsList className="inline-flex h-auto p-0.5 sm:p-1 bg-muted rounded-full">
                  <TabsTrigger 
                    value="shorts" 
                    className="h-8 sm:h-9 px-3 sm:px-4 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-1.5 sm:gap-2"
                  >
                    <Zap className="w-3 h-3 sm:w-4 sm:h-4" />
                    <span className="text-xs sm:text-sm font-medium">Shorts</span>
                  </TabsTrigger>
                  <TabsTrigger 
                    value="long" 
                    className="h-8 sm:h-9 px-3 sm:px-4 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-1.5 sm:gap-2"
                  >
                    <Video className="w-3 h-3 sm:w-4 sm:h-4" />
                    <span className="text-xs sm:text-sm font-medium">Videos</span>
                  </TabsTrigger>
                  <TabsTrigger 
                    value="channels" 
                    className="h-8 sm:h-9 px-3 sm:px-4 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-1.5 sm:gap-2"
                  >
                    <Users className="w-3 h-3 sm:w-4 sm:h-4" />
                    <span className="text-xs sm:text-sm font-medium">Channels</span>
                  </TabsTrigger>
                  <TabsTrigger 
                    value="subscriptions" 
                    className="h-8 sm:h-9 px-3 sm:px-4 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-1.5 sm:gap-2"
                  >
                    <Bell className="w-3 h-3 sm:w-4 sm:h-4" />
                    <span className="text-xs sm:text-sm font-medium">Subscribed</span>
                  </TabsTrigger>
                </TabsList>
              </div>
            </div>
            
            {/* Filters and Actions Section */}
            <div className="flex items-center justify-between px-3 sm:px-4 pb-2">
              <div className="flex items-center gap-2">
                <Button 
                  variant="ghost" 
                  size="icon" 
                  className="h-8 w-8 sm:h-9 sm:w-9 touch-manipulation"
                >
                  <Search className="w-4 h-4 sm:w-5 sm:h-5" />
                </Button>
                <Button 
                  variant="ghost" 
                  size="icon" 
                  className="h-8 w-8 sm:h-9 sm:w-9 touch-manipulation"
                >
                  <Filter className="w-4 h-4 sm:w-5 sm:h-5" />
                </Button>
              </div>
              
              {currentTab === 'channels' && (
                <Dialog open={isCreateChannelOpen} onOpenChange={setIsCreateChannelOpen}>
                  <DialogTrigger asChild>
                    <Button size="sm" className="flex items-center gap-1.5 sm:gap-2 text-xs sm:text-sm px-2 sm:px-3 h-8 sm:h-9 touch-manipulation">
                      <Plus className="w-3 h-3 sm:w-4 sm:h-4" />
                      <span className="hidden sm:inline">Create Channel</span>
                      <span className="sm:hidden">Create</span>
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="w-[90vw] max-w-md">
                    <DialogHeader>
                      <DialogTitle>Create New Channel</DialogTitle>
                    </DialogHeader>
                    <div className="space-y-4">
                      <div>
                        <label className="text-sm font-medium">Channel Name</label>
                        <Input
                          placeholder="Enter channel name"
                          value={newChannelData.name}
                          onChange={(e) => setNewChannelData(prev => ({ ...prev, name: e.target.value }))}
                          className="mt-1"
                        />
                      </div>
                      <div>
                        <label className="text-sm font-medium">Description</label>
                        <Textarea
                          placeholder="Describe your channel"
                          value={newChannelData.description}
                          onChange={(e) => setNewChannelData(prev => ({ ...prev, description: e.target.value }))}
                          className="mt-1"
                          rows={3}
                        />
                      </div>
                      <div>
                        <label className="text-sm font-medium">Category</label>
                        <select 
                          className="w-full p-2 mt-1 border border-border rounded-md bg-background text-foreground"
                          value={newChannelData.category}
                          onChange={(e) => setNewChannelData(prev => ({ ...prev, category: e.target.value }))}
                        >
                          <option value="food">Food & Cooking</option>
                          <option value="entertainment">Entertainment</option>
                          <option value="education">Education</option>
                          <option value="travel">Travel</option>
                          <option value="music">Music</option>
                          <option value="business">Business</option>
                        </select>
                      </div>
                      <div className="flex flex-col sm:flex-row gap-2 pt-4">
                        <Button onClick={handleCreateChannel} className="flex-1 touch-manipulation">
                          Create Channel
                        </Button>
                        <Button 
                          variant="outline" 
                          onClick={() => setIsCreateChannelOpen(false)}
                          className="touch-manipulation"
                        >
                          Cancel
                        </Button>
                      </div>
                    </div>
                  </DialogContent>
                </Dialog>
              )}
            </div>
          </div>

          {/* Categories */}
          <ScrollArea className="w-full">
            <div className="w-full overflow-x-auto scrollbar-hide">
              <div className="w-full pb-4">
                <div className="flex gap-2 overflow-x-auto scrollbar-hide pb-1 px-4">
                {categories.map((category) => {
                  const IconComponent = category.icon;
                  return (
                    <Button
                      key={category.id}
                      variant={selectedCategory === category.id ? 'default' : 'outline'}
                      size="sm"
                      className="flex items-center gap-2 whitespace-nowrap flex-shrink-0"
                      onClick={() => setSelectedCategory(category.id)}
                    >
                      <IconComponent className="w-4 h-4" />
                      {category.name}
                    </Button>
                  );
                })}
                </div>
              </div>
            </div>
          </ScrollArea>
        </div>

        {/* Shorts Tab Content */}
        <TabsContent value="shorts" className="flex-1 m-0 p-0">
          <div className="h-full relative bg-black">
            {/* Loading state */}
            {shortVideosLoading && (
              <div className="flex items-center justify-center h-full">
                <div className="text-white text-center">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white mx-auto mb-4"></div>
                  <p>Loading shorts...</p>
                </div>
              </div>
            )}

            {/* Error state */}
            {shortVideosError && !shortVideosLoading && (
              <div className="flex items-center justify-center h-full">
                <div className="text-white text-center p-4">
                  <p className="text-lg mb-2">Failed to load shorts</p>
                  <p className="text-sm text-gray-300 mb-4">Please check your connection and try again</p>
                  <Button variant="outline" onClick={() => window.location.reload()}>
                    Retry
                  </Button>
                </div>
              </div>
            )}

            {/* Empty state */}
            {!shortVideosLoading && !shortVideosError && effectiveShortVideos.length === 0 && (
              <div className="flex items-center justify-center h-full">
                <div className="text-white text-center p-4">
                  <p className="text-lg mb-2">No shorts available</p>
                  <p className="text-sm text-gray-300">Check back later for new content</p>
                </div>
              </div>
            )}

            {/* Main content */}
            {!shortVideosLoading && !shortVideosError && currentShort && (
              <>
                {/* Video Player */}
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="relative w-full h-full max-w-md mx-auto">
                    {currentShort?.videoUrl && !videoError ? (
                      <>
                        <video
                          ref={videoRef}
                          src={currentShort.videoUrl}
                          className="w-full h-full object-cover"
                          poster={currentShort.thumbnail}
                          loop
                          muted={isMuted}
                          playsInline
                          onPlay={() => setIsPlaying(true)}
                          onPause={() => setIsPlaying(false)}
                          onLoadStart={() => {
                            setVideoLoading(true);
                            console.log('Video loading started');
                          }}
                          onLoadedData={() => {
                            setVideoLoading(false);
                            console.log('Video loaded');
                          }}
                          onError={(e) => {
                            console.error('Video error:', e);
                            setVideoError(true);
                            setVideoLoading(false);
                          }}
                        />
                        {videoLoading && (
                          <div className="absolute inset-0 flex items-center justify-center bg-black/30">
                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
                          </div>
                        )}
                      </>
                    ) : (
                      <ImageWithFallback
                        src={currentShort.thumbnail}
                        alt={currentShort.title}
                        className="w-full h-full object-cover"
                      />
                    )}
                    
                    {/* Video Controls Overlay */}
                    <div className="absolute inset-0 flex items-center justify-center">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="w-20 h-20 rounded-full bg-black/60 backdrop-blur-sm text-white hover:bg-black/80 border-2 border-white/20 hover:border-white/40 transition-all duration-200 hover:scale-110 active:scale-95 shadow-lg"
                        onClick={togglePlayPause}
                      >
                        {isPlaying ? (
                          <Pause className="w-10 h-10" />
                        ) : (
                          <Play className="w-10 h-10 ml-1" />
                        )}
                      </Button>
                    </div>

                    {/* Navigation Areas */}
                    <div 
                      className="absolute top-0 bottom-20 left-0 right-1/2 z-10"
                      onClick={() => {
                        handleShortSwipe('down');
                        // Ensure video plays after navigation
                        setTimeout(() => {
                          if (videoRef.current && !videoRef.current.paused) return;
                          videoRef.current?.play().then(() => {
                            setIsPlaying(true);
                          }).catch((error) => {
                            console.log('Play failed:', error);
                            setIsPlaying(false);
                          });
                        }, 100);
                      }}
                    />
                    <div 
                      className="absolute top-0 bottom-20 right-0 left-1/2 z-10"
                      onClick={() => handleShortSwipe('up')}
                    />
                  </div>
                </div>

                {/* Right Side Actions */}
                <div className="absolute right-4 bottom-32 flex flex-col gap-6 z-20">
                  <div className="flex flex-col items-center gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      className={`w-12 h-12 rounded-full text-white hover:bg-white/20 ${
                        likedVideos.includes(currentShort.id) ? 'bg-red-500' : 'bg-black/50'
                      }`}
                      onClick={() => handleVideoLike(currentShort.id)}
                    >
                      <Heart className={`w-6 h-6 ${likedVideos.includes(currentShort.id) ? 'fill-current' : ''}`} />
                    </Button>
                    <span className="text-white text-xs">{formatViews(currentShort.likes)}</span>
                  </div>

                  <div className="flex flex-col items-center gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="w-12 h-12 rounded-full bg-black/50 text-white hover:bg-white/20"
                    >
                      <MessageCircle className="w-6 h-6" />
                    </Button>
                    <span className="text-white text-xs">247</span>
                  </div>

                  <div className="flex flex-col items-center gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="w-12 h-12 rounded-full bg-black/50 text-white hover:bg-white/20"
                      onClick={() => handleVideoShare(currentShort.id)}
                    >
                      <Share className="w-6 h-6" />
                    </Button>
                    <span className="text-white text-xs">Share</span>
                  </div>

                  <div className="flex flex-col items-center gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="w-12 h-12 rounded-full bg-black/50 text-white hover:bg-white/20"
                    >
                      <Bookmark className="w-6 h-6" />
                    </Button>
                  </div>

                  <div className="flex flex-col items-center gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="w-12 h-12 rounded-full bg-black/50 text-white hover:bg-white/20"
                      onClick={toggleMute}
                    >
                      {isMuted ? <VolumeX className="w-6 h-6" /> : <Volume2 className="w-6 h-6" />}
                    </Button>
                  </div>
                </div>

                {/* Bottom Info */}
                <div className="absolute bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-black/90 via-black/50 to-transparent text-white z-20 pointer-events-none">
                  <div className="flex items-end gap-3 pointer-events-auto">
                    <Avatar className="w-12 h-12 border-2 border-white/80 shadow-lg flex-shrink-0">
                      <AvatarImage src={currentShort.channel.avatar} />
                      <AvatarFallback className="bg-white/20 text-white border-white/40">
                        {currentShort.channel.name.charAt(0)}
                      </AvatarFallback>
                    </Avatar>
                    
                    <div className="flex-1 min-w-0 space-y-2">
                      <div className="flex items-center gap-2 flex-wrap">
                        <span className="font-semibold text-white drop-shadow-sm truncate">
                          {currentShort.channel.name}
                        </span>
                        {currentShort.channel.verified && (
                          <CheckCircle className="w-4 h-4 text-blue-400 fill-current drop-shadow-sm flex-shrink-0" />
                        )}
                        <Button
                          size="sm"
                          className="bg-red-600 hover:bg-red-700 text-white h-7 px-3 rounded-full font-medium shadow-lg border-0 flex-shrink-0 transition-all duration-200 hover:scale-105"
                          onClick={() => handleChannelSubscribe(currentShort.channel.id)}
                        >
                          <UserPlus className="w-3 h-3 mr-1.5" />
                          Follow
                        </Button>
                      </div>
                      
                      <div className="space-y-1">
                        <p className="text-sm text-white/95 line-clamp-2 leading-relaxed drop-shadow-sm">
                          {currentShort.description}
                        </p>
                        <div className="flex items-center gap-4 text-xs text-white/75 drop-shadow-sm">
                          <span className="flex items-center gap-1">
                            <Eye className="w-3 h-3" />
                            {formatViews(currentShort.views)}
                          </span>
                          <span className="flex items-center gap-1">
                            <Clock className="w-3 h-3" />
                            {currentShort.uploadTime}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                  
                  {/* Safety gradient for better text readability */}
                  <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent pointer-events-none -z-10" />
                </div>

                {/* Navigation Indicators */}
                <div className="absolute top-1/2 right-2 transform -translate-y-1/2 flex flex-col gap-2 z-10">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="w-8 h-8 rounded-full bg-black/30 text-white hover:bg-black/50"
                    onClick={() => handleShortSwipe('down')}
                    disabled={currentShortIndex === 0}
                  >
                    <ArrowUp className="w-4 h-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="w-8 h-8 rounded-full bg-black/30 text-white hover:bg-black/50"
                    onClick={() => handleShortSwipe('up')}
                    disabled={currentShortIndex === filteredShorts.length - 1}
                  >
                    <ArrowDown className="w-4 h-4" />
                  </Button>
                </div>
              </>
            )}
          </div>
        </TabsContent>

        {/* Long Videos Tab Content */}
        <TabsContent value="long" className="flex-1 m-0">
          <ScrollArea className="h-full">
            <div className="p-4 space-y-4">
              {/* Loading state */}
              {longVideosLoading && (
                <div className="flex items-center justify-center py-12">
                  <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
                    <p>Loading videos...</p>
                  </div>
                </div>
              )}

              {/* Error state */}
              {longVideosError && !longVideosLoading && (
                <div className="flex items-center justify-center py-12">
                  <div className="text-center p-4">
                    <p className="text-lg mb-2">Failed to load videos</p>
                    <p className="text-sm text-muted-foreground mb-4">Please check your connection and try again</p>
                    <Button variant="outline" onClick={() => window.location.reload()}>
                      Retry
                    </Button>
                  </div>
                </div>
              )}

              {/* Empty state */}
              {!longVideosLoading && !longVideosError && filteredLongs.length === 0 && (
                <div className="flex items-center justify-center py-12">
                  <div className="text-center p-4">
                    <p className="text-lg mb-2">No videos available</p>
                    <p className="text-sm text-muted-foreground">Check back later for new content</p>
                  </div>
                </div>
              )}

              {/* Main content */}
              {!longVideosLoading && !longVideosError && filteredLongs.map((video) => (
                <Card 
                  key={video.id} 
                  className="overflow-hidden cursor-pointer hover:shadow-lg transition-shadow"
                  onClick={() => onVideoPlay(video.id)}
                >
                  <div className="relative">
                    <div className="aspect-[3/1] relative">
                      {video.videoUrl ? (
                        <video
                          src={video.videoUrl}
                          poster={video.thumbnail}
                          className="w-full h-full object-cover"
                          muted
                          preload="metadata"
                          onMouseEnter={(e) => {
                            const target = e.target as HTMLVideoElement;
                            target.currentTime = 10; // Preview at 10 seconds
                          }}
                        />
                      ) : (
                        <ImageWithFallback
                          src={video.thumbnail}
                          alt={video.title}
                          className="w-full h-full object-cover"
                        />
                      )}
                      <div className="absolute inset-0 bg-black/20 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="w-16 h-16 rounded-full bg-black/50 text-white hover:bg-black/70"
                        >
                          <Play className="w-8 h-8 ml-1" />
                        </Button>
                      </div>
                      <Badge className="absolute bottom-2 right-2 bg-black/70 text-white">
                        {video.duration}
                      </Badge>
                      {video.isLive && (
                        <Badge className="absolute top-2 left-2 bg-red-600 text-white">
                          LIVE
                        </Badge>
                      )}
                      {watchedVideos.includes(video.id) && (
                        <div className="absolute bottom-0 left-0 right-0 h-1 bg-red-600"></div>
                      )}
                    </div>
                  </div>

                  <CardContent className="p-3 sm:p-4 w-full">
                    <div className="flex gap-2 sm:gap-3 w-full min-w-0">
                      <Avatar 
                        className="w-8 h-8 sm:w-10 sm:h-10 flex-shrink-0 cursor-pointer touch-manipulation"
                        onClick={(e) => {
                          e.stopPropagation();
                          const channel = effectiveChannels.find(c => c.id === video.channel.id);
                          if (channel) handleChannelClick(channel);
                        }}
                      >
                        <AvatarImage src={video.channel.avatar} />
                        <AvatarFallback className="text-xs sm:text-sm">{video.channel.name.charAt(0)}</AvatarFallback>
                      </Avatar>
                      
                      <div className="flex-1 min-w-0 overflow-hidden">
                        <h3 className="text-sm sm:text-base font-medium line-clamp-2 mb-1 leading-tight">{video.title}</h3>
                        <div className="flex items-center gap-1 sm:gap-2 text-xs sm:text-sm text-muted-foreground mb-1 flex-wrap">
                          <span 
                            className="cursor-pointer hover:text-primary truncate max-w-32 sm:max-w-none touch-manipulation"
                            onClick={(e) => {
                              e.stopPropagation();
                              const channel = effectiveChannels.find(c => c.id === video.channel.id);
                              if (channel) handleChannelClick(channel);
                            }}
                          >
                            {video.channel.name}
                          </span>
                          {video.channel.verified && (
                            <CheckCircle className="w-3 h-3 text-blue-500 fill-current flex-shrink-0" />
                          )}
                        </div>
                        <div className="flex items-center gap-1 sm:gap-2 text-xs text-muted-foreground flex-wrap">
                          <span className="whitespace-nowrap">{formatViews(video.views)} views</span>
                          <span className="hidden sm:inline">•</span>
                          <span className="whitespace-nowrap">{video.uploadTime}</span>
                        </div>
                      </div>

                      <div className="flex flex-col gap-1 sm:gap-2 flex-shrink-0">
                        <Button variant="ghost" size="icon" className="w-7 h-7 sm:w-8 sm:h-8 touch-manipulation">
                          <MoreVertical className="w-4 h-4" />
                        </Button>
                      </div>
                    </div>

                    <div className="flex items-center justify-between mt-2 sm:mt-3 w-full min-w-0 gap-2">
                      <div className="flex items-center gap-1 sm:gap-2 lg:gap-4 flex-1 min-w-0 overflow-x-auto scrollbar-hide">
                        <Button
                          variant="ghost"
                          size="sm"
                          className={`gap-1 sm:gap-2 flex-shrink-0 h-8 sm:h-9 px-2 sm:px-3 text-xs sm:text-sm touch-manipulation ${likedVideos.includes(video.id) ? 'text-red-500' : ''}`}
                          onClick={(e) => {
                            e.stopPropagation();
                            onVideoLike(video.id);
                          }}
                        >
                          <Heart className={`w-3 h-3 sm:w-4 sm:h-4 ${likedVideos.includes(video.id) ? 'fill-current' : ''}`} />
                          <span className="hidden sm:inline">{formatViews(video.likes)}</span>
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="sm" 
                          className="gap-1 sm:gap-2 flex-shrink-0 h-8 sm:h-9 px-2 sm:px-3 text-xs sm:text-sm touch-manipulation"
                          onClick={(e) => e.stopPropagation()}
                        >
                          <MessageCircle className="w-3 h-3 sm:w-4 sm:h-4" />
                          <span className="hidden sm:inline">142</span>
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="sm" 
                          className="gap-1 sm:gap-2 flex-shrink-0 h-8 sm:h-9 px-2 sm:px-3 text-xs sm:text-sm touch-manipulation"
                          onClick={(e) => {
                            e.stopPropagation();
                            onVideoShare(video.id);
                          }}
                        >
                          <Share className="w-3 h-3 sm:w-4 sm:h-4" />
                          <span className="hidden sm:inline">Share</span>
                        </Button>
                      </div>
                      
                      <Button
                        variant={subscribedChannelIds.includes(video.channel.id) ? "outline" : "default"}
                        size="sm"
                        className="flex-shrink-0 h-8 sm:h-9 px-2 sm:px-3 text-xs sm:text-sm touch-manipulation"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleChannelSubscribe(video.channel.id);
                        }}
                      >
                        {subscribedChannelIds.includes(video.channel.id) ? (
                          <>
                            <UserMinus className="w-3 h-3 sm:w-4 sm:h-4 mr-1 sm:mr-2" />
                            <span className="hidden sm:inline">Subscribed</span>
                            <span className="sm:hidden">✓</span>
                          </>
                        ) : (
                          <>
                            <UserPlus className="w-3 h-3 sm:w-4 sm:h-4 mr-1 sm:mr-2" />
                            <span className="hidden sm:inline">Subscribe</span>
                            <span className="sm:hidden">+</span>
                          </>
                        )}
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </ScrollArea>
        </TabsContent>

        {/* Channels Tab Content */}
        <TabsContent value="channels" className="flex-1 m-0">
          <ScrollArea className="h-full">
            <div className="p-4 space-y-4">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-semibold">Discover Channels</h2>
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm">
                    <Grid3X3 className="w-4 h-4" />
                  </Button>
                  <Button variant="ghost" size="sm">
                    <List className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              {/* Loading state */}
              {channelsLoading && (
                <div className="flex items-center justify-center py-12">
                  <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
                    <p>Loading channels...</p>
                  </div>
                </div>
              )}

              {/* Error state */}
              {channelsError && !channelsLoading && (
                <div className="flex items-center justify-center py-12">
                  <div className="text-center p-4">
                    <p className="text-lg mb-2">Failed to load channels</p>
                    <p className="text-sm text-muted-foreground mb-4">Please check your connection and try again</p>
                    <Button variant="outline" onClick={() => window.location.reload()}>
                      Retry
                    </Button>
                  </div>
                </div>
              )}

              {/* Empty state */}
              {!channelsLoading && !channelsError && effectiveChannels.length === 0 && (
                <div className="flex items-center justify-center py-12">
                  <div className="text-center p-4">
                    <p className="text-lg mb-2">No channels available</p>
                    <p className="text-sm text-muted-foreground">Channels will appear here once videos are loaded</p>
                  </div>
                </div>
              )}

              {/* Main content */}
              {!channelsLoading && !channelsError && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {filteredChannels.map((channel) => (
                    <Card
                      key={channel.id}
                      className="overflow-hidden cursor-pointer hover:shadow-lg transition-shadow"
                      onClick={() => handleChannelClick(channel)}
                    >
                      {channel.banner && (
                        <div className="h-24 relative">
                          <ImageWithFallback
                            src={channel.banner}
                            alt={`${channel.name} banner`}
                            className="w-full h-full object-cover"
                          />
                          <div className="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent" />
                        </div>
                      )}

                      <CardContent className="p-4">
                        <div className="flex items-start gap-4">
                          <Avatar className="w-16 h-16 border-2 border-background">
                            <AvatarImage src={channel.avatar} />
                            <AvatarFallback>{channel.name.charAt(0)}</AvatarFallback>
                          </Avatar>

                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 mb-1">
                              <h3 className="font-medium line-clamp-1">{channel.name}</h3>
                              {channel.verified && (
                                <CheckCircle className="w-4 h-4 text-blue-500 fill-current" />
                              )}
                            </div>

                            <div className="flex items-center gap-3 text-sm text-muted-foreground mb-2">
                              <span>{formatSubscribers(channel.subscribers)} subscribers</span>
                              <span>•</span>
                              <span>{channel.videos} videos</span>
                            </div>

                            <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
                              {channel.description}
                            </p>

                            <Button
                              variant={channel.isSubscribed ? "outline" : "default"}
                              size="sm"
                              className="w-full"
                              onClick={(e) => {
                                e.stopPropagation();
                                handleChannelSubscribe(channel.id);
                              }}
                            >
                              {channel.isSubscribed ? (
                                <>
                                  <UserMinus className="w-4 h-4 mr-2" />
                                  Subscribed
                                </>
                              ) : (
                                <>
                                  <UserPlus className="w-4 h-4 mr-2" />
                                  Subscribe
                                </>
                              )}
                            </Button>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          </ScrollArea>
        </TabsContent>

        {/* Subscriptions Tab Content */}
        <TabsContent value="subscriptions" className="flex-1 m-0">
          <ScrollArea className="h-full">
            <div className="p-3 sm:p-4 lg:p-6">
              <div className="flex items-center justify-between mb-4 sm:mb-6">
                <h2 className="text-base sm:text-lg font-semibold">My Subscriptions</h2>
                <Badge variant="secondary" className="text-xs sm:text-sm">
                  {subscribedChannelsList.length} channels
                </Badge>
              </div>

              {subscribedChannelsList.length === 0 ? (
                <div className="text-center py-8 sm:py-12">
                  <Users className="w-12 h-12 sm:w-16 sm:h-16 text-muted-foreground mx-auto mb-3 sm:mb-4" />
                  <h3 className="text-base sm:text-lg font-medium mb-2">No subscriptions yet</h3>
                  <p className="text-sm sm:text-base text-muted-foreground mb-4 sm:mb-6 px-4">Subscribe to channels to see them here</p>
                  <Button 
                    onClick={() => setCurrentTab('channels')}
                    className="touch-manipulation"
                    size="sm"
                  >
                    Discover Channels
                  </Button>
                </div>
              ) : (
                <div className="space-y-3 sm:space-y-4">
                  {subscribedChannelsList.map((channel) => (
                    <Card 
                      key={channel.id}
                      className="cursor-pointer hover:shadow-md transition-shadow touch-manipulation"
                      onClick={() => handleChannelClick(channel)}
                    >
                      <CardContent className="p-3 sm:p-4">
                        <div className="flex items-start sm:items-center gap-3 sm:gap-4">
                          <Avatar className="w-10 h-10 sm:w-12 sm:h-12 flex-shrink-0">
                            <AvatarImage src={channel.avatar} />
                            <AvatarFallback className="text-sm">{channel.name.charAt(0)}</AvatarFallback>
                          </Avatar>
                          
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-1 sm:gap-2 mb-1">
                              <h3 className="text-sm sm:text-base font-medium truncate">{channel.name}</h3>
                              {channel.verified && (
                                <CheckCircle className="w-3 h-3 sm:w-4 sm:h-4 text-blue-500 fill-current flex-shrink-0" />
                              )}
                            </div>
                            <div className="flex items-center gap-2 sm:gap-3 text-xs sm:text-sm text-muted-foreground flex-wrap">
                              <span className="whitespace-nowrap">{formatSubscribers(channel.subscribers)} subscribers</span>
                              <span className="hidden sm:inline">•</span>
                              <span className="whitespace-nowrap">{channel.videos} videos</span>
                            </div>
                          </div>
                          
                          <div className="flex flex-col sm:flex-row items-end sm:items-center gap-2">
                            <Button 
                              variant="outline" 
                              size="icon" 
                              className="w-8 h-8 sm:w-9 sm:h-9 touch-manipulation flex-shrink-0"
                              onClick={(e) => {
                                e.stopPropagation();
                                // Toggle notifications logic here
                              }}
                            >
                              {channel.notificationsEnabled ? (
                                <Bell className="w-3 h-3 sm:w-4 sm:h-4" />
                              ) : (
                                <BellOff className="w-3 h-3 sm:w-4 sm:h-4" />
                              )}
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              className="text-xs sm:text-sm h-8 sm:h-9 px-2 sm:px-3 touch-manipulation flex-shrink-0"
                              onClick={(e) => {
                                e.stopPropagation();
                                handleChannelSubscribe(channel.id);
                              }}
                            >
                              <UserMinus className="w-3 h-3 sm:w-4 sm:h-4 mr-1 sm:mr-2" />
                              <span className="hidden sm:inline">Unsubscribe</span>
                              <span className="sm:hidden">Unsub</span>
                            </Button>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}

              {/* Recent Videos from Subscribed Channels */}
              {subscribedChannelsList.length > 0 && (
                <div className="mt-6 sm:mt-8">
                  <h3 className="text-base sm:text-lg font-semibold mb-3 sm:mb-4">Latest from your subscriptions</h3>
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
                    {effectiveLongVideos
                      .filter(video => subscribedChannelIds.includes(video.channel.id))
                      .slice(0, 8)
                      .map((video) => (
                        <Card 
                          key={video.id} 
                          className="overflow-hidden cursor-pointer hover:shadow-lg transition-shadow touch-manipulation"
                          onClick={() => onVideoPlay(video.id, video)}
                        >
                          <div className="aspect-video relative">
                            <ImageWithFallback
                              src={video.thumbnail}
                              alt={video.title}
                              className="w-full h-full object-cover"
                            />
                            <Badge className="absolute bottom-1 right-1 sm:bottom-2 sm:right-2 bg-black/70 text-white text-xs">
                              {video.duration}
                            </Badge>
                            {video.isLive && (
                              <Badge className="absolute top-1 left-1 sm:top-2 sm:left-2 bg-red-600 text-white text-xs">
                                LIVE
                              </Badge>
                            )}
                          </div>
                          <CardContent className="p-2 sm:p-3">
                            <h4 className="text-sm sm:text-base font-medium line-clamp-2 mb-1 sm:mb-2 leading-tight">{video.title}</h4>
                            <div className="flex items-center gap-1 sm:gap-2 text-xs sm:text-sm text-muted-foreground mb-1 flex-wrap">
                              <span className="truncate max-w-24 sm:max-w-none">{video.channel.name}</span>
                              {video.channel.verified && (
                                <CheckCircle className="w-2.5 h-2.5 sm:w-3 sm:h-3 text-blue-500 fill-current flex-shrink-0" />
                              )}
                            </div>
                            <div className="flex items-center gap-1 sm:gap-2 text-xs text-muted-foreground flex-wrap">
                              <span className="whitespace-nowrap">{formatViews(video.views)} views</span>
                              <span className="hidden sm:inline">•</span>
                              <span className="whitespace-nowrap">{video.uploadTime}</span>
                            </div>
                          </CardContent>
                        </Card>
                      ))}
                  </div>
                </div>
              )}
            </div>
          </ScrollArea>
        </TabsContent>
      </Tabs>
    </div>
  );
}