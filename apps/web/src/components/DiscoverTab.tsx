import React, { useState, useEffect, useMemo } from 'react';
import { Heart, MessageCircle, Share, MoreVertical, Plus, Camera, MapPin, Clock, Users, Zap, Star, Play, Send, Bookmark, UserPlus, UserCheck, Hash, Copy, ChevronDown, ChevronUp, TrendingUp, Search, Eye, Filter, X, Globe, Utensils, Building, Calendar, Bell, ArrowRight, Navigation, Music, Flame, Sparkles, Repeat2, BarChart3, AlertTriangle, Loader2 } from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Input } from './ui/input';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";
import {
  useGetTrendingPostsQuery,
  useGetTrendingTopicsQuery,
  useGetSuggestedUsersQuery,
  useSearchDiscoverQuery
} from '../services/microservicesApi';

interface TrendingPost {
  id: string;
  author: {
    name: string;
    username: string;
    avatar?: string;
    verified?: boolean;
    type: 'user' | 'merchant' | 'channel';
  };
  content: string;
  images?: string[];
  timestamp: string;
  likes: number;
  reposts: number;
  comments: number;
  shares: number;
  isLiked?: boolean;
  isReposted?: boolean;
  location?: string;
  tags?: string[];
  type: 'text' | 'image' | 'live' | 'product';
  engagementRate: number;
  trendingRank?: number;
}

interface TrendingTopic {
  id: string;
  title: string;
  category: 'hashtag' | 'topic' | 'person' | 'place';
  posts: number;
  trending: boolean;
  trendingRank?: number;
  description?: string;
  image?: string;
  engagementRate: number;
  timeframe: '24h' | '7d' | '30d';
}

interface SuggestedUser {
  id: string;
  name: string;
  username: string;
  avatar?: string;
  verified?: boolean;
  bio?: string;
  followers: number;
  mutualFollows: number;
  category: 'food' | 'culture' | 'business' | 'influencer' | 'local';
}

interface DiscoverTabProps {
  user: any;
  onPostShare?: (postId: string, postData: any) => void;
}

export function DiscoverTab({ user, onPostShare }: DiscoverTabProps) {
  const [discoverCategory, setDiscoverCategory] = useState<string>('trending');
  const [discoverSearch, setDiscoverSearch] = useState('');
  const [showFilters, setShowFilters] = useState(false);
  const [followedTopics, setFollowedTopics] = useState<string[]>(['1', '2']);
  const [followedUsers, setFollowedUsers] = useState<string[]>(['1', '3']);
  const [likedPosts, setLikedPosts] = useState<string[]>([]);
  const [repostedPosts, setRepostedPosts] = useState<string[]>([]);
  const [bookmarkedPosts, setBookmarkedPosts] = useState<string[]>([]);
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  // RTK Query hooks for discover data
  const { data: postsData, isLoading: postsLoading, error: postsError } = useGetTrendingPostsQuery({
    category: discoverCategory,
    limit: 20
  }, {
    skip: !isMounted
  });

  const { data: topicsData, isLoading: topicsLoading, error: topicsError } = useGetTrendingTopicsQuery({
    timeframe: '24h',
    limit: 10
  }, {
    skip: !isMounted
  });

  const { data: usersData, isLoading: usersLoading, error: usersError } = useGetSuggestedUsersQuery({
    category: 'all',
    limit: 8
  }, {
    skip: !isMounted
  });

  // Fallback trending posts data
  const fallbackPostsData: TrendingPost[] = [
    {
      id: '1',
      author: {
        name: 'Thai Food Explorer',
        username: '@thaifoodexplorer',
        avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face',
        verified: true,
        type: 'channel'
      },
      content: 'This hidden gem in Chinatown serves the most authentic Tom Yum I\'ve ever tasted! 🍜 The secret? 30 years of perfecting the recipe. The broth is so rich and aromatic, you can smell it from the street. Must try! #BangkokFood #HiddenGems #TomYum',
      images: ['https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=600&h=400&fit=crop'],
      timestamp: '2h ago',
      likes: 2847,
      reposts: 892,
      comments: 456,
      shares: 234,
      location: 'Chinatown, Bangkok',
      tags: ['#BangkokFood', '#HiddenGems', '#TomYum'],
      type: 'image',
      engagementRate: 8.2,
      trendingRank: 1
    },
    {
      id: '2',
      author: {
        name: 'Bangkok Street Life',
        username: '@bangkokstreets',
        avatar: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=150&h=150&fit=crop&crop=face',
        verified: true,
        type: 'channel'
      },
      content: 'Early morning at Chatuchak Market is pure magic ✨ The energy, the colors, the smells... This is Thailand! Watch vendors setting up for day while sipping fresh coconut water 🥥',
      images: ['https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=600&h=400&fit=crop'],
      timestamp: '4h ago',
      likes: 1923,
      reposts: 567,
      comments: 289,
      shares: 145,
      location: 'Chatuchak Weekend Market',
      tags: ['#ChatuchakMarket', '#BangkokLife', '#Thailand'],
      type: 'image',
      engagementRate: 6.8,
      trendingRank: 2
    }
  ];

  // Transform RTK Query posts data to local format
  const trendingPosts = useMemo(() => {
    if (!postsData || postsLoading || postsError) {
      return fallbackPostsData;
    }
    return postsData.map(post => ({
      id: post.id,
      author: {
        name: post.author?.name || 'Unknown User',
        username: post.author?.username || '@unknown',
        avatar: post.author?.avatar,
        verified: post.author?.verified || false,
        type: post.author?.type || 'user'
      },
      content: post.content || 'No content available',
      images: post.images || [],
      timestamp: post.timestamp || 'Unknown time',
      likes: post.likes || 0,
      reposts: post.reposts || 0,
      comments: post.comments || 0,
      shares: post.shares || 0,
      isLiked: post.isLiked || false,
      isReposted: post.isReposted || false,
      location: post.location,
      tags: post.tags || [],
      type: post.type || 'text',
      engagementRate: post.engagementRate || 0,
      trendingRank: post.trendingRank
    }));
  }, [postsData, postsLoading, postsError]);

  // Fallback trending topics data
  const fallbackTopicsData: TrendingTopic[] = [
    {
      id: '1',
      title: '#SongkranFestival',
      category: 'hashtag',
      posts: 127456,
      trending: true,
      trendingRank: 1,
      description: 'Thailand\'s water festival celebrations',
      engagementRate: 15.2,
      timeframe: '24h'
    },
    {
      id: '2',
      title: '#BangkokFood',
      category: 'hashtag',
      posts: 89234,
      trending: true,
      trendingRank: 2,
      description: 'Street food and restaurant discoveries',
      engagementRate: 12.8,
      timeframe: '24h'
    },
    {
      id: '3',
      title: '#ThaiCulture',
      category: 'hashtag',
      posts: 56789,
      trending: true,
      trendingRank: 3,
      description: 'Traditional culture and modern Thailand',
      engagementRate: 9.4,
      timeframe: '7d'
    }
  ];

  // Transform RTK Query topics data to local format
  const trendingTopics = useMemo(() => {
    if (!topicsData || topicsLoading || topicsError) {
      return fallbackTopicsData;
    }
    return topicsData.map(topic => ({
      id: topic.id,
      title: topic.title || 'Unknown Topic',
      category: topic.category || 'hashtag',
      posts: topic.posts || 0,
      trending: topic.trending || false,
      trendingRank: topic.trendingRank,
      description: topic.description,
      image: topic.image,
      engagementRate: topic.engagementRate || 0,
      timeframe: topic.timeframe || '24h'
    }));
  }, [topicsData, topicsLoading, topicsError]);

  // Fallback suggested users data
  const fallbackUsersData: SuggestedUser[] = [
    {
      id: '1',
      name: 'Chef Siriporn',
      username: '@chefsiriporn',
      avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face',
      verified: true,
      bio: 'Michelin-starred Thai chef sharing authentic recipes',
      followers: 245000,
      mutualFollows: 12,
      category: 'food'
    },
    {
      id: '2',
      name: 'Bangkok Explorer',
      username: '@bangkokexplorer',
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
      verified: true,
      bio: 'Discovering hidden gems in Bangkok daily',
      followers: 89000,
      mutualFollows: 8,
      category: 'local'
    }
  ];

  // Transform RTK Query users data to local format
  const suggestedUsers = useMemo(() => {
    if (!usersData || usersLoading || usersError) {
      return fallbackUsersData;
    }
    return usersData.map(user => ({
      id: user.id,
      name: user.name || 'Unknown User',
      username: user.username || '@unknown',
      avatar: user.avatar,
      verified: user.verified || false,
      bio: user.bio,
      followers: user.followers || 0,
      mutualFollows: user.mutualFollows || 0,
      category: user.category || 'influencer'
    }));
  }, [usersData, usersLoading, usersError]);

  // Event handlers
  const handleFollowTopic = (topicId: string, topicTitle: string) => {
    if (followedTopics.includes(topicId)) {
      setFollowedTopics(prev => prev.filter(id => id !== topicId));
      toast.success(`Unfollowed ${topicTitle}`);
    } else {
      setFollowedTopics(prev => [...prev, topicId]);
      toast.success(`Following ${topicTitle}`);
    }
  };

  const handleFollowUser = (userId: string, userName: string) => {
    if (followedUsers.includes(userId)) {
      setFollowedUsers(prev => prev.filter(id => id !== userId));
      toast.success(`Unfollowed @${userName}`);
    } else {
      setFollowedUsers(prev => [...prev, userId]);
      toast.success(`Following @${userName}`);
    }
  };

  const handleLikePost = (postId: string) => {
    if (likedPosts.includes(postId)) {
      setLikedPosts(prev => prev.filter(id => id !== postId));
    } else {
      setLikedPosts(prev => [...prev, postId]);
      toast.success('Post liked!');
    }
  };

  const handleRepost = (postId: string) => {
    if (repostedPosts.includes(postId)) {
      setRepostedPosts(prev => prev.filter(id => id !== postId));
      toast.success('Repost removed');
    } else {
      setRepostedPosts(prev => [...prev, postId]);
      toast.success('Reposted!');
    }
  };

  const handleBookmarkPost = (postId: string) => {
    if (bookmarkedPosts.includes(postId)) {
      setBookmarkedPosts(prev => prev.filter(id => id !== postId));
      toast.success('Removed from bookmarks');
    } else {
      setBookmarkedPosts(prev => [...prev, postId]);
      toast.success('Post bookmarked');
    }
  };

  const handleSharePost = (postId: string, post: TrendingPost) => {
    if (onPostShare) {
      onPostShare(postId, post);
    } else {
      navigator.clipboard.writeText(`Check out this trending post: ${post.content.slice(0, 100)}...`);
      toast.success('Post link copied to clipboard!');
    }
  };

  const handleTopicClick = (topic: TrendingTopic) => {
    toast.success(`Exploring ${topic.title}`);
  };

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'trending':
        return <TrendingUp className="w-4 h-4" />;
      case 'hashtags':
        return <Hash className="w-4 h-4" />;
      case 'people':
        return <Users className="w-4 h-4" />;
      case 'posts':
        return <MessageCircle className="w-4 h-4" />;
      default:
        return <Globe className="w-4 h-4" />;
    }
  };

  const filteredTopics = trendingTopics.filter(topic => {
    if (discoverCategory === 'hashtags') return topic.category === 'hashtag';
    if (discoverCategory === 'people') return topic.category === 'person';
    return true;
  }).filter(topic => 
    discoverSearch === '' || 
    topic.title.toLowerCase().includes(discoverSearch.toLowerCase()) ||
    topic.description?.toLowerCase().includes(discoverSearch.toLowerCase())
  );

  const filteredPosts = trendingPosts.filter(post =>
    discoverSearch === '' ||
    post.content.toLowerCase().includes(discoverSearch.toLowerCase()) ||
    post.author.name.toLowerCase().includes(discoverSearch.toLowerCase()) ||
    post.tags?.some(tag => tag.toLowerCase().includes(discoverSearch.toLowerCase()))
  );

  const renderPost = (post: TrendingPost) => (
    <Card key={post.id} className="mb-4 hover:shadow-md transition-shadow">
      <CardContent className="p-2 sm:p-3 lg:p-4 space-y-2 sm:space-y-3 w-full overflow-hidden">
        {/* Post Header - Ultra Mobile Constrained */}
        <div className="flex items-start gap-2 w-full min-w-0">
          <Avatar className="w-8 h-8 sm:w-9 sm:h-9 lg:w-10 lg:h-10 cursor-pointer hover:ring-2 hover:ring-ring hover:ring-offset-1 transition-all flex-shrink-0">
            <AvatarImage src={post.author.avatar} />
            <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white text-xs">
              {post.author.name.charAt(0)}
            </AvatarFallback>
          </Avatar>
          
          <div className="flex-1 min-w-0 overflow-hidden">
            {/* Author Info - Maximum Width Constraint */}
            <div className="w-full min-w-0 space-y-0.5">
              <div className="flex items-center justify-between w-full min-w-0 gap-1">
                <div className="flex items-center gap-1 min-w-0 flex-1 overflow-hidden">
                  <span className="font-medium cursor-pointer hover:underline text-xs sm:text-sm truncate max-w-16 sm:max-w-20">
                    {post.author.name}
                  </span>
                  <span className="text-[10px] sm:text-xs text-muted-foreground truncate max-w-12 sm:max-w-16">
                    {post.author.username}
                  </span>
                  {post.author.verified && (
                    <Star className="w-2.5 h-2.5 sm:w-3 sm:h-3 text-yellow-500 fill-current flex-shrink-0" />
                  )}
                </div>
                
                {post.trendingRank && (
                  <Badge className="bg-gradient-to-r from-orange-500 to-red-500 text-white text-[8px] sm:text-[10px] px-1 sm:px-1.5 py-0.5 flex-shrink-0">
                    <TrendingUp className="w-2 h-2 sm:w-2.5 sm:h-2.5 mr-0.5" />
                    #{post.trendingRank}
                  </Badge>
                )}
              </div>
              
              {/* Post Metadata - Single Line Overflow */}
              <div className="flex items-center gap-1 text-[9px] sm:text-[10px] text-muted-foreground w-full min-w-0 overflow-hidden">
                <span className="flex-shrink-0 truncate max-w-12 sm:max-w-16">{post.timestamp}</span>
                
                {post.location && (
                  <>
                    <span className="flex-shrink-0">•</span>
                    <div className="flex items-center gap-0.5 min-w-0">
                      <MapPin className="w-2 h-2 sm:w-2.5 sm:h-2.5 flex-shrink-0" />
                      <span className="truncate max-w-10 sm:max-w-12">{post.location}</span>
                    </div>
                  </>
                )}
                
                <span className="flex-shrink-0">•</span>
                <div className="flex items-center gap-0.5 flex-shrink-0">
                  <BarChart3 className="w-2 h-2 sm:w-2.5 sm:h-2.5" />
                  <span className="whitespace-nowrap">{post.engagementRate}%</span>
                </div>
              </div>
            </div>
          </div>
          
          {/* More Options Menu - Ultra Compact */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-6 w-6 sm:h-7 sm:w-7 flex-shrink-0 touch-manipulation">
                <MoreVertical className="w-3 h-3 sm:w-3.5 sm:h-3.5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              <DropdownMenuItem onClick={() => handleBookmarkPost(post.id)}>
                <Bookmark className={`w-4 h-4 mr-2 ${bookmarkedPosts.includes(post.id) ? 'fill-current' : ''}`} />
                {bookmarkedPosts.includes(post.id) ? 'Remove Bookmark' : 'Bookmark'}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleSharePost(post.id, post)}>
                <Share className="w-4 h-4 mr-2" />
                Share Post
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Bell className="w-4 h-4 mr-2" />
                Get Notifications
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        {/* Post Content - Width Constrained */}
        <div className="w-full min-w-0 overflow-hidden">
          <div className="prose prose-sm max-w-none w-full">
            <p className="leading-relaxed text-sm sm:text-base text-foreground m-0 break-words overflow-wrap-anywhere w-full">
              {post.content.split(' ').map((word, index) => {
                if (word.startsWith('#')) {
                  return (
                    <span
                      key={index}
                      className="text-primary cursor-pointer hover:underline font-medium hover:text-primary/80 transition-colors break-all"
                      onClick={() => toast.success(`Exploring ${word}`)}
                    >
                      {word}{' '}
                    </span>
                  );
                }
                return word + ' ';
              })}
            </p>
          </div>
        </div>

        {/* Post Images - Contained Width */}
        {post.images && post.images.length > 0 && (
          <div className="w-full min-w-0 space-y-2">
            <div className="relative rounded-lg sm:rounded-xl overflow-hidden bg-muted w-full">
              <ImageWithFallback
                src={post.images[0]}
                alt="Post content"
                className="w-full h-auto max-h-60 sm:max-h-72 lg:max-h-80 object-cover cursor-pointer hover:scale-105 transition-transform duration-300"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/5 to-transparent pointer-events-none" />
            </div>
            
            {post.images.length > 1 && (
              <div className="flex gap-1 sm:gap-1.5 overflow-x-auto pb-1 scrollbar-hide w-full">
                {post.images.slice(1, 4).map((image, index) => (
                  <div key={index} className="relative flex-shrink-0">
                    <ImageWithFallback
                      src={image}
                      alt={`Post image ${index + 2}`}
                      className="w-10 h-10 sm:w-12 sm:h-12 lg:w-16 lg:h-16 object-cover rounded-sm sm:rounded-md cursor-pointer hover:opacity-80 transition-opacity"
                    />
                    {index === 2 && post.images.length > 4 && (
                      <div className="absolute inset-0 bg-black/60 rounded-sm sm:rounded-md flex items-center justify-center">
                        <span className="text-white text-[8px] sm:text-[10px] font-medium">
                          +{post.images.length - 4}
                        </span>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Post Actions - Maximum Width Constraint */}
        <div className="pt-2 border-t border-border/50 w-full min-w-0">
          <div className="flex items-center w-full min-w-0 gap-1">
            {/* Main Action Buttons - No Overflow */}
            <div className="flex items-center gap-0.5 sm:gap-1 flex-1 min-w-0 overflow-hidden">
              <Button
                variant="ghost"
                size="sm"
                className="h-7 sm:h-8 px-1 sm:px-1.5 text-muted-foreground hover:text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-950/20 rounded-full transition-colors touch-manipulation min-w-0 flex-shrink-0"
              >
                <MessageCircle className="w-3 h-3 sm:w-3.5 sm:h-3.5 mr-0.5" />
                <span className="text-[9px] sm:text-[10px] font-medium">
                  {post.comments > 999 ? `${Math.floor(post.comments / 1000)}k` : post.comments}
                </span>
              </Button>
              
              <Button
                variant="ghost"
                size="sm"
                className={`h-7 sm:h-8 px-1 sm:px-1.5 hover:text-green-500 hover:bg-green-50 dark:hover:bg-green-950/20 rounded-full transition-colors touch-manipulation min-w-0 flex-shrink-0 ${
                  repostedPosts.includes(post.id) ? 'text-green-500 bg-green-50 dark:bg-green-950/20' : 'text-muted-foreground'
                }`}
                onClick={() => handleRepost(post.id)}
              >
                <Repeat2 className="w-3 h-3 sm:w-3.5 sm:h-3.5 mr-0.5" />
                <span className="text-[9px] sm:text-[10px] font-medium">
                  {(post.reposts + (repostedPosts.includes(post.id) ? 1 : 0)) > 999 
                    ? `${Math.floor((post.reposts + (repostedPosts.includes(post.id) ? 1 : 0)) / 1000)}k` 
                    : (post.reposts + (repostedPosts.includes(post.id) ? 1 : 0))}
                </span>
              </Button>
              
              <Button
                variant="ghost"
                size="sm"
                className={`h-7 sm:h-8 px-1 sm:px-1.5 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-950/20 rounded-full transition-colors touch-manipulation min-w-0 flex-shrink-0 ${
                  likedPosts.includes(post.id) ? 'text-red-500 bg-red-50 dark:bg-red-950/20' : 'text-muted-foreground'
                }`}
                onClick={() => handleLikePost(post.id)}
              >
                <Heart className={`w-3 h-3 sm:w-3.5 sm:h-3.5 mr-0.5 ${likedPosts.includes(post.id) ? 'fill-current' : ''}`} />
                <span className="text-[9px] sm:text-[10px] font-medium">
                  {(post.likes + (likedPosts.includes(post.id) ? 1 : 0)) > 999 
                    ? `${Math.floor((post.likes + (likedPosts.includes(post.id) ? 1 : 0)) / 1000)}k` 
                    : (post.likes + (likedPosts.includes(post.id) ? 1 : 0))}
                </span>
              </Button>
              
              <Button
                variant="ghost"
                size="sm"
                className="h-7 sm:h-8 px-1 sm:px-1.5 text-muted-foreground hover:text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-950/20 rounded-full transition-colors touch-manipulation min-w-0 flex-shrink-0"
                onClick={() => handleSharePost(post.id, post)}
              >
                <Share className="w-3 h-3 sm:w-3.5 sm:h-3.5" />
                <span className="text-[9px] sm:text-[10px] font-medium hidden sm:inline ml-0.5">
                  {post.shares > 999 ? `${Math.floor(post.shares / 1000)}k` : post.shares}
                </span>
              </Button>
            </div>
            
            {/* Bookmark Button - Fixed Size */}
            <Button
              variant="ghost"
              size="sm"
              className={`h-7 sm:h-8 w-7 sm:w-8 rounded-full transition-colors touch-manipulation flex-shrink-0 ${
                bookmarkedPosts.includes(post.id) 
                  ? 'text-primary bg-primary/10 hover:bg-primary/20' 
                  : 'text-muted-foreground hover:text-primary hover:bg-primary/10'
              }`}
              onClick={() => handleBookmarkPost(post.id)}
            >
              <Bookmark className={`w-3 h-3 sm:w-3.5 sm:h-3.5 ${bookmarkedPosts.includes(post.id) ? 'fill-current' : ''}`} />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );

  return (
    <ScrollArea className="h-full">
      <div className="space-y-4">
        {/* Twitter-Style Search and Filters */}
        <div className="px-4 space-y-3 pt-4">
          <div className="flex items-center gap-2">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  placeholder="Search posts, hashtags, people..."
                  value={discoverSearch}
                  onChange={(e) => setDiscoverSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
              {discoverSearch && (
                <Button
                  variant="ghost"
                  size="sm"
                  className="absolute right-1 top-1/2 transform -translate-y-1/2 h-8 w-8 p-0"
                  onClick={() => setDiscoverSearch('')}
                >
                  <X className="w-4 h-4" />
                </Button>
              )}
            </div>
          </div>

          {/* Category Filters - Twitter Style */}
          <div className="w-full">
            <div className="flex gap-1.5 sm:gap-2 overflow-x-auto pb-2 scrollbar-hide w-full">
              {['trending', 'posts', 'hashtags', 'people'].map((category) => (
                <Button
                  key={category}
                  variant={discoverCategory === category ? "default" : "outline"}
                  size="sm"
                  className="flex items-center gap-1.5 sm:gap-2 whitespace-nowrap flex-shrink-0 min-w-0 touch-manipulation px-3 sm:px-4 h-8 sm:h-9"
                  onClick={() => setDiscoverCategory(category)}
                >
                  {getCategoryIcon(category)}
                  <span className="text-xs sm:text-sm">{category.charAt(0).toUpperCase() + category.slice(1)}</span>
                </Button>
              ))}
            </div>
          </div>
        </div>

        <div className="px-2 sm:px-3 lg:px-4 space-y-2 sm:space-y-3 lg:space-y-4 max-w-full min-w-0 overflow-hidden">
          {/* Trending Posts */}
          {(discoverCategory === 'trending' || discoverCategory === 'posts') && (
            <div className="space-y-3 w-full -mx-4">
              <div className="flex items-center justify-between px-4">
                <h3 className="font-medium flex items-center gap-2">
                  <TrendingUp className="w-5 h-5 text-orange-500" />
                  Trending Posts
                </h3>
                <Badge variant="secondary" className="text-xs">
                  {filteredPosts.length} posts
                </Badge>
              </div>

              {postsLoading ? (
                <Card>
                  <CardContent className="p-8 text-center">
                    <Loader2 className="w-8 h-8 text-primary mx-auto mb-4 animate-spin" />
                    <p className="text-muted-foreground">Loading trending posts...</p>
                  </CardContent>
                </Card>
              ) : postsError ? (
                <Card>
                  <CardContent className="p-8 text-center">
                    <AlertTriangle className="w-8 h-8 text-destructive mx-auto mb-4" />
                    <p className="text-destructive mb-2">Failed to load trending posts</p>
                    <p className="text-muted-foreground text-sm">Using cached content</p>
                  </CardContent>
                </Card>
              ) : filteredPosts.length === 0 ? (
                <Card>
                  <CardContent className="p-8 text-center">
                    <MessageCircle className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                    <p className="text-muted-foreground">No posts found matching your search</p>
                  </CardContent>
                </Card>
              ) : (
                <div className="space-y-0">
                  {filteredPosts.map(renderPost)}
                </div>
              )}
            </div>
          )}

          {/* Trending Hashtags & Topics */}
          {discoverCategory !== 'posts' && (
            <div className="space-y-2 sm:space-y-3 w-full -mx-4 px-0">
              <div className="flex items-center justify-between">
                <h3 className="font-medium flex items-center gap-2">
                  {discoverCategory === 'hashtags' ? (
                    <>
                      <Hash className="w-5 h-5 text-blue-500" />
                      Trending Hashtags
                    </>
                  ) : discoverCategory === 'people' ? (
                    <>
                      <Users className="w-5 h-5 text-purple-500" />
                      Who to Follow
                    </>
                  ) : (
                    <>
                      <TrendingUp className="w-5 h-5 text-orange-500" />
                      Trending Now
                    </>
                  )}
                </h3>
                <Badge variant="secondary" className="text-xs">
                  {discoverCategory === 'people' ? suggestedUsers.length : filteredTopics.length} {discoverCategory === 'people' ? 'suggestions' : 'topics'}
                </Badge>
              </div>

              {discoverCategory === 'people' ? (
                <div className="space-y-3">
                  {usersLoading ? (
                    <Card>
                      <CardContent className="p-8 text-center">
                        <Loader2 className="w-8 h-8 text-primary mx-auto mb-4 animate-spin" />
                        <p className="text-muted-foreground">Loading suggested users...</p>
                      </CardContent>
                    </Card>
                  ) : usersError ? (
                    <Card>
                      <CardContent className="p-8 text-center">
                        <AlertTriangle className="w-8 h-8 text-destructive mx-auto mb-4" />
                        <p className="text-destructive mb-2">Failed to load suggested users</p>
                        <p className="text-muted-foreground text-sm">Using cached content</p>
                      </CardContent>
                    </Card>
                  ) : suggestedUsers.length === 0 ? (
                    <Card>
                      <CardContent className="p-8 text-center">
                        <Users className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                        <p className="text-muted-foreground">No users found to suggest</p>
                      </CardContent>
                    </Card>
                  ) : (
                    suggestedUsers.map((suggestedUser) => (
                    <Card key={suggestedUser.id} className="hover:shadow-md transition-shadow">
                      <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <Avatar className="w-12 h-12">
                              <AvatarImage src={suggestedUser.avatar} />
                              <AvatarFallback>{suggestedUser.name.charAt(0)}</AvatarFallback>
                            </Avatar>
                            <div className="flex-1">
                              <div className="flex items-center gap-2">
                                <span className="font-medium">{suggestedUser.name}</span>
                                <span className="text-sm text-muted-foreground">{suggestedUser.username}</span>
                                {suggestedUser.verified && (
                                  <Star className="w-4 h-4 text-yellow-500 fill-current" />
                                )}
                              </div>
                              <p className="text-sm text-muted-foreground line-clamp-1">{suggestedUser.bio}</p>
                              <div className="flex items-center gap-4 text-xs text-muted-foreground mt-1">
                                <span>{suggestedUser.followers.toLocaleString()} followers</span>
                                <span>{suggestedUser.mutualFollows} mutual follows</span>
                              </div>
                            </div>
                          </div>
                          <Button
                            variant={followedUsers.includes(suggestedUser.id) ? "secondary" : "default"}
                            size="sm"
                            onClick={() => handleFollowUser(suggestedUser.id, suggestedUser.username.replace('@', ''))}
                          >
                            {followedUsers.includes(suggestedUser.id) ? (
                              <>
                                <UserCheck className="w-3 h-3 mr-1" />
                                Following
                              </>
                            ) : (
                              <>
                                <UserPlus className="w-3 h-3 mr-1" />
                                Follow
                              </>
                            )}
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                    ))
                  )}
                </div>
              ) : (
                <div className="space-y-3">
                  {topicsLoading ? (
                    <Card>
                      <CardContent className="p-8 text-center">
                        <Loader2 className="w-8 h-8 text-primary mx-auto mb-4 animate-spin" />
                        <p className="text-muted-foreground">Loading trending topics...</p>
                      </CardContent>
                    </Card>
                  ) : topicsError ? (
                    <Card>
                      <CardContent className="p-8 text-center">
                        <AlertTriangle className="w-8 h-8 text-destructive mx-auto mb-4" />
                        <p className="text-destructive mb-2">Failed to load trending topics</p>
                        <p className="text-muted-foreground text-sm">Using cached content</p>
                      </CardContent>
                    </Card>
                  ) : filteredTopics.length === 0 ? (
                    <Card>
                      <CardContent className="p-8 text-center">
                        <Hash className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                        <p className="text-muted-foreground">No topics found matching your search</p>
                      </CardContent>
                    </Card>
                  ) : (
                    filteredTopics.map((topic) => (
                    <Card key={topic.id} className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => handleTopicClick(topic)}>
                      <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                          <div className="flex-1">
                            <div className="flex items-center gap-2 mb-1">
                              <span className="font-medium text-lg">{topic.title}</span>
                              {topic.trending && topic.trendingRank && (
                                <Badge className="bg-orange-500 text-white text-xs">
                                  #{topic.trendingRank} Trending
                                </Badge>
                              )}
                              {topic.category === 'hashtag' && <Hash className="w-3 h-3 text-blue-500" />}
                              {topic.category === 'person' && <Users className="w-3 h-3 text-purple-500" />}
                              {topic.category === 'place' && <MapPin className="w-3 h-3 text-green-500" />}
                            </div>
                            <p className="text-sm text-muted-foreground mb-2">{topic.description}</p>
                            <div className="flex items-center gap-4 text-sm text-muted-foreground">
                              <span>{topic.posts.toLocaleString()} posts</span>
                              <span>{topic.engagementRate}% engagement</span>
                              <span>Trending {topic.timeframe}</span>
                            </div>
                          </div>
                          <div className="flex items-center gap-2 ml-4">
                            <Button
                              variant={followedTopics.includes(topic.id) ? "secondary" : "outline"}
                              size="sm"
                              onClick={(e) => {
                                e.stopPropagation();
                                handleFollowTopic(topic.id, topic.title);
                              }}
                            >
                              {followedTopics.includes(topic.id) ? (
                                <>
                                  <UserCheck className="w-3 h-3 mr-1" />
                                  Following
                                </>
                              ) : (
                                <>
                                  <UserPlus className="w-3 h-3 mr-1" />
                                  Follow
                                </>
                              )}
                            </Button>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                    ))
                  )}
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </ScrollArea>
  );
}