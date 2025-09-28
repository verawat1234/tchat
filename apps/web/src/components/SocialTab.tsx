import React, { useState, useEffect, useMemo } from 'react';
import { Heart, MessageCircle, Share, MoreVertical, Plus, Camera, MapPin, Clock, Users, Zap, Star, Play, Send, Bookmark, UserPlus, UserCheck, Hash, Copy, ChevronDown, ChevronUp, TrendingUp, Search, Eye, Filter, X, Globe, Utensils, Building, Calendar, Bell, ArrowRight, Navigation, Music, Flame, Sparkles, ChevronLeft, ChevronRight, Volume2, PlayCircle, PauseCircle, SkipForward, AlertTriangle } from 'lucide-react';
import { CreatePostSection } from './CreatePostSection';
import { DiscoverTab } from './DiscoverTab';
import { EventsTab } from './EventsTab';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Input } from './ui/input';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './ui/dialog';
import { VisuallyHidden } from './ui/visually-hidden';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { Separator } from './ui/separator';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { Progress } from './ui/progress';
import { toast } from "sonner";
import {
  useGetSocialFeedQuery,
  useGetSocialStoriesQuery,
  useGetUserFriendsQuery,
  useLikeSocialPostMutation,
  useCreateSocialPostMutation,
  useFollowUserMutation
} from '../services/microservicesApi';

interface SocialTabProps {
  user: any;
  onLiveStreamClick?: (streamId: string) => void;
  onPostShare?: (postId: string, postData: any) => void;
}

interface Friend {
  id: string;
  name: string;
  username: string;
  avatar?: string;
  isOnline: boolean;
  lastSeen?: string;
  mutualFriends?: number;
  status?: string;
  isFollowing: boolean;
}

interface Comment {
  id: string;
  user: {
    name: string;
    avatar?: string;
  };
  text: string;
  timestamp: string;
  likes: number;
  isLiked: boolean;
}

interface Post {
  id: string;
  author: {
    name: string;
    avatar?: string;
    verified?: boolean;
    type: 'user' | 'merchant' | 'channel';
  };
  content: string;
  images?: string[];
  timestamp: string;
  likes: number;
  comments: number;
  shares: number;
  isLiked?: boolean;
  location?: string;
  tags?: string[];
  type: 'text' | 'image' | 'live' | 'product';
  product?: {
    name: string;
    price: number;
    currency: string;
  };
  liveData?: {
    viewers: number;
    startTime: string;
    isLive: boolean;
  };
  source?: 'following' | 'trending' | 'interest' | 'sponsored';
}

interface Moment {
  id: string;
  author: {
    name: string;
    avatar?: string;
  };
  preview: string;
  isViewed: boolean;
  isLive?: boolean;
  content?: string;
  timestamp?: string;
  media?: {
    type: 'image' | 'video';
    url: string;
    duration?: number;
  }[];
  expiresAt: string;
}

export function SocialTab({ user, onLiveStreamClick, onPostShare }: SocialTabProps) {
  const [isMounted, setIsMounted] = useState(false);
  const [likedPosts, setLikedPosts] = useState<string[]>([]);
  const [bookmarkedPosts, setBookmarkedPosts] = useState<string[]>([]);
  const [followingUsers, setFollowingUsers] = useState<string[]>(['1', '2', '3', '5']);
  const [commentsOpen, setCommentsOpen] = useState<string | null>(null);
  const [newComment, setNewComment] = useState('');
  const [postComments, setPostComments] = useState<{ [key: string]: Comment[] }>({});
  const [viewingStory, setViewingStory] = useState<Moment | null>(null);
  const [storyIndex, setStoryIndex] = useState(0);
  const [storyProgress, setStoryProgress] = useState(0);
  const [storyMediaIndex, setStoryMediaIndex] = useState(0);
  const [showEventsTab, setShowEventsTab] = useState(false);

  // RTK Query hooks for social data
  const { data: socialFeedData, isLoading: feedLoading, error: feedError } = useGetSocialFeedQuery(
    { type: 'all', limit: 20, page: 1 },
    { skip: !isMounted }
  );

  const { data: socialStoriesData, isLoading: storiesLoading, error: storiesError } = useGetSocialStoriesQuery(
    { active: true, limit: 20 },
    { skip: !isMounted }
  );

  const { data: friendsData, isLoading: friendsLoading, error: friendsError } = useGetUserFriendsQuery(
    { status: 'active', limit: 50 },
    { skip: !isMounted }
  );

  const [likeSocialPost] = useLikeSocialPostMutation();
  const [createSocialPost] = useCreateSocialPostMutation();
  const [followUser] = useFollowUserMutation();

  // Mount effect
  useEffect(() => {
    setIsMounted(true);
  }, []);
  
  // Create Post state
  const [createPostOpen, setCreatePostOpen] = useState(false);
  const [newPostText, setNewPostText] = useState('');
  const [selectedImages, setSelectedImages] = useState<string[]>([]);
  const [postLocation, setPostLocation] = useState('');
  const [postPrivacy, setPostPrivacy] = useState<'public' | 'friends' | 'private'>('public');

  // User's posts state - Real dynamic posts
  const [userPosts, setUserPosts] = useState<Post[]>([]);

  // Moment creation and viewing state
  const [createStoryOpen, setCreateStoryOpen] = useState(false);
  const [storyText, setStoryText] = useState('');
  const [storyBackground, setStoryBackground] = useState('#1a1a1a');
  const [storyType, setStoryType] = useState<'text' | 'image' | 'video'>('text');

  // Transform RTK Query stories data to Moment format
  const stories: Moment[] = useMemo(() => {
    if (!socialStoriesData || storiesLoading || storiesError) {
      // Fallback data while loading or on error
      return [
        {
          id: '0',
          author: { name: 'Your Moment', avatar: user?.avatar || '' },
          preview: '',
          isViewed: false,
          content: 'Add to your moment',
          timestamp: 'now',
          expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
        },
        {
          id: 'fallback-1',
          author: { name: 'Sarah Johnson', avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face' },
          preview: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=400&h=600&fit=crop',
          isViewed: false,
          content: 'Amazing Pad Thai at Chatuchak Market! 🍜✨',
          timestamp: '2h ago',
          expiresAt: new Date(Date.now() + 22 * 60 * 60 * 1000).toISOString(),
          media: [{
            type: 'image',
            url: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=400&h=600&fit=crop',
            duration: 5
          }]
        }
      ];
    }

    // Always include user's own story creation option first
    const userStory: Moment = {
      id: '0',
      author: { name: 'Your Moment', avatar: user?.avatar || '' },
      preview: '',
      isViewed: false,
      content: 'Add to your moment',
      timestamp: 'now',
      expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
    };

    // Transform API data to Moment format
    const apiStories = socialStoriesData.map(story => ({
      id: story.id,
      author: {
        name: story.author?.name || 'Unknown User',
        avatar: story.author?.avatar || ''
      },
      preview: story.preview || story.media?.[0]?.url || '',
      isViewed: story.isViewed || false,
      isLive: story.isLive || false,
      content: story.content || '',
      timestamp: story.timestamp || 'Just now',
      media: story.media || [],
      expiresAt: story.expiresAt || new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
    }));

    return [userStory, ...apiStories];
  }, [socialStoriesData, storiesLoading, storiesError, user?.avatar]);

  // Transform RTK Query feed data to Post format
  const feedPosts: Post[] = useMemo(() => {
    if (!socialFeedData || feedLoading || feedError) {
      // Fallback data while loading or on error
      return [
        {
          id: 'fallback-1',
          author: {
            name: 'Thai Food Explorer',
            avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face',
            verified: true,
            type: 'channel'
          },
          content: 'TRENDING: Secret street food spots only locals know about! 🏮🍜 This hidden gem in Chinatown serves the most authentic Tom Yum I\'ve ever tasted. The owner has been perfecting this recipe for 30 years!',
          images: ['https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=600&h=400&fit=crop'],
          timestamp: '15 min ago',
          likes: 847,
          comments: 123,
          shares: 45,
          location: 'Chinatown, Bangkok',
          tags: ['#HiddenGems', '#TomYum', '#Chinatown', '#StreetFood'],
          type: 'image',
          source: 'trending'
        },
        {
          id: 'fallback-2',
          author: {
            name: 'Sarah Johnson',
            avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face',
            verified: false,
            type: 'user'
          },
          content: 'Just tried the most amazing Pad Thai at Chatuchak Market! 🍜✨ The vendor taught me his secret ingredient - tamarind paste mixed with palm sugar. Mind blown! 🤯',
          images: ['https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=600&h=400&fit=crop'],
          timestamp: '25 min ago',
          likes: 47,
          comments: 12,
          shares: 3,
          location: 'Chatuchak Weekend Market, Bangkok',
          tags: ['#PadThai', '#StreetFood', '#Bangkok'],
          type: 'image',
          source: 'following'
        }
      ];
    }

    // Transform API data to Post format
    return socialFeedData.map(post => ({
      id: post.id,
      author: {
        name: post.author?.name || 'Unknown User',
        avatar: post.author?.avatar || '',
        verified: post.author?.verified || false,
        type: post.author?.type || 'user'
      },
      content: post.content || '',
      images: post.images || [],
      timestamp: post.timestamp || 'Just now',
      likes: post.likes || 0,
      comments: post.comments || 0,
      shares: post.shares || 0,
      isLiked: post.isLiked || false,
      location: post.location || '',
      tags: post.tags || [],
      type: post.type || 'text',
      source: post.source || 'interest',
      product: post.product,
      liveData: post.liveData
    }));
  }, [socialFeedData, feedLoading, feedError]);

  // Transform RTK Query friends data to Friend format
  const friends: Friend[] = useMemo(() => {
    if (!friendsData || friendsLoading || friendsError) {
      // Fallback data while loading or on error
      return [
        {
          id: 'fallback-1',
          name: 'Sarah Johnson',
          username: '@sarah_foodie',
          avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face',
          isOnline: true,
          mutualFriends: 12,
          status: 'Exploring Bangkok street food! 🍜',
          isFollowing: true
        },
        {
          id: 'fallback-2',
          name: 'Mike Chen',
          username: '@mike_travels',
          avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face',
          isOnline: true,
          lastSeen: '2 min ago',
          mutualFriends: 8,
          status: 'Live streaming from floating market!',
          isFollowing: true
        }
      ];
    }

    // Transform API data to Friend format
    return friendsData.map(friend => ({
      id: friend.id,
      name: friend.name || 'Unknown User',
      username: friend.username || '@user',
      avatar: friend.avatar || '',
      isOnline: friend.isOnline || false,
      lastSeen: friend.lastSeen || '',
      mutualFriends: friend.mutualFriends || 0,
      status: friend.status || '',
      isFollowing: friend.isFollowing || false
    }));
  }, [friendsData, friendsLoading, friendsError]);

  // Moment progress timer
  useEffect(() => {
    if (viewingStory && viewingStory.media) {
      const currentMedia = viewingStory.media[storyMediaIndex];
      if (currentMedia) {
        const duration = currentMedia.duration || 5;
        const interval = setInterval(() => {
          setStoryProgress(prev => {
            if (prev >= 100) {
              // Move to next media or close story
              if (storyMediaIndex < viewingStory.media!.length - 1) {
                setStoryMediaIndex(prev => prev + 1);
                return 0;
              } else {
                setViewingStory(null);
                setStoryMediaIndex(0);
                return 0;
              }
            }
            return prev + (100 / (duration * 10)); // Update every 100ms
          });
        }, 100);

        return () => clearInterval(interval);
      }
    }
  }, [viewingStory, storyMediaIndex]);

  // Event handlers
  const handleLike = async (postId: string) => {
    const isCurrentlyLiked = likedPosts.includes(postId);

    try {
      // Optimistic update
      if (isCurrentlyLiked) {
        setLikedPosts(likedPosts.filter(id => id !== postId));
      } else {
        setLikedPosts([...likedPosts, postId]);
        toast.success('Post liked!');
      }

      // RTK Query mutation
      await likeSocialPost({ postId, isLiked: !isCurrentlyLiked }).unwrap();
    } catch (error) {
      // Revert optimistic update on error
      if (isCurrentlyLiked) {
        setLikedPosts([...likedPosts, postId]);
      } else {
        setLikedPosts(likedPosts.filter(id => id !== postId));
      }
      toast.error('Failed to like post');
      console.error('Error liking post:', error);
    }
  };

  const handleJoinLive = (postId: string) => {
    if (onLiveStreamClick) {
      onLiveStreamClick('featured-stream');
    } else {
      toast.success('Joining live stream...');
    }
  };

  const handleBookmark = (postId: string) => {
    if (bookmarkedPosts.includes(postId)) {
      setBookmarkedPosts(bookmarkedPosts.filter(id => id !== postId));
      toast.success('Removed from bookmarks');
    } else {
      setBookmarkedPosts([...bookmarkedPosts, postId]);
      toast.success('Added to bookmarks');
    }
  };

  const handleFollow = (userId: string) => {
    if (followingUsers.includes(userId)) {
      setFollowingUsers(followingUsers.filter(id => id !== userId));
      toast.success('Unfollowed user');
    } else {
      setFollowingUsers([...followingUsers, userId]);
      toast.success('Following user');
    }
  };

  const handleHashtagClick = (hashtag: string) => {
    toast.success(`Searching for ${hashtag}`);
  };

  const handleComment = (postId: string) => {
    setCommentsOpen(postId);
  };

  const handleShare = (postId: string) => {
    // Search in both user posts and feed posts
    const post = [...userPosts, ...feedPosts].find(p => p.id === postId);
    if (post && onPostShare) {
      onPostShare(postId, post);
    } else {
      navigator.clipboard.writeText(`https://telegram-sea.app/post/${postId}`);
      toast.success('Post link copied to clipboard!');
    }
  };

  const handleAddComment = (postId: string) => {
    if (newComment.trim()) {
      const comment: Comment = {
        id: Date.now().toString(),
        user: {
          name: user?.name || 'You',
          avatar: user?.avatar
        },
        text: newComment,
        timestamp: 'just now',
        likes: 0,
        isLiked: false
      };

      setPostComments(prev => ({
        ...prev,
        [postId]: [...(prev[postId] || []), comment]
      }));

      setNewComment('');
      toast.success('Comment added!');
    }
  };

  const handleStoryClick = (story: Moment, index: number) => {
    if (story.author.name === 'Your Moment') {
      setCreateStoryOpen(true);
      return;
    }
    
    setViewingStory(story);
    setStoryIndex(index);
    setStoryProgress(0);
    setStoryMediaIndex(0);
  };

  const handleCreatePost = () => {
    if (newPostText.trim()) {
      toast.success('Post shared to your feed! 🎉');
      setNewPostText('');
      setSelectedImages([]);
      setPostLocation('');
      setCreatePostOpen(false);
    }
  };

  const handleCreateStory = () => {
    if (storyText.trim() || selectedImages.length > 0) {
      toast.success('Moment created! Visible for 24 hours 📸');
      setStoryText('');
      setSelectedImages([]);
      setCreateStoryOpen(false);
    }
  };

  const handleAddFriend = (userId: string, userName: string) => {
    toast.success(`Friend request sent to ${userName}!`);
  };

  // Real post creation handlers for CreatePostSection
  const handleCreateTextPost = (textContent: string) => {
    if (!textContent.trim()) return;

    const newPost: Post = {
      id: `user_post_${Date.now()}`,
      author: {
        name: user?.name || 'You',
        avatar: user?.avatar,
        verified: false,
        type: 'user'
      },
      content: textContent.trim(),
      timestamp: 'just now',
      likes: 0,
      comments: 0,
      shares: 0,
      type: 'text',
      source: 'following'
    };

    // Add to user's posts at the top of the feed
    setUserPosts(prev => [newPost, ...prev]);
    toast.success('Post created successfully! 🎉');
  };

  const handleCreatePhotoPost = () => {
    // Simulate photo post creation with placeholder
    const newPost: Post = {
      id: `user_photo_${Date.now()}`,
      author: {
        name: user?.name || 'You',
        avatar: user?.avatar,
        verified: false,
        type: 'user'
      },
      content: 'Just captured a beautiful moment in Bangkok! 📸✨',
      images: ['https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=600&h=400&fit=crop'],
      timestamp: 'just now',
      likes: 0,
      comments: 0,
      shares: 0,
      type: 'image',
      location: 'Bangkok, Thailand',
      tags: ['#Bangkok', '#Photography'],
      source: 'following'
    };

    setUserPosts(prev => [newPost, ...prev]);
    toast.success('Photo post created! 📷');
  };

  const handleCreateGalleryPost = () => {
    // Simulate gallery post with multiple images
    const newPost: Post = {
      id: `user_gallery_${Date.now()}`,
      author: {
        name: user?.name || 'You',
        avatar: user?.avatar,
        verified: false,
        type: 'user'
      },
      content: 'Amazing street food adventure today! Multiple shots of the best Bangkok has to offer 🍜🥘🍲',
      images: [
        'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=600&h=400&fit=crop',
        'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?w=600&h=400&fit=crop'
      ],
      timestamp: 'just now',
      likes: 0,
      comments: 0,
      shares: 0,
      type: 'image',
      location: 'Street Food District, Bangkok',
      tags: ['#StreetFood', '#Bangkok', '#Gallery'],
      source: 'following'
    };

    setUserPosts(prev => [newPost, ...prev]);
    toast.success('Gallery post created! 🖼️');
  };

  const handleCreateLocationPost = () => {
    // Simulate location-based post
    const newPost: Post = {
      id: `user_location_${Date.now()}`,
      author: {
        name: user?.name || 'You',
        avatar: user?.avatar,
        verified: false,
        type: 'user'
      },
      content: 'Currently at this amazing floating market! The energy here is incredible 🛶✨',
      images: ['https://images.unsplash.com/photo-1743485753872-3b24372fcd24?w=600&h=400&fit=crop'],
      timestamp: 'just now',
      likes: 0,
      comments: 0,
      shares: 0,
      type: 'image',
      location: 'Damnoen Saduak Floating Market, Thailand',
      tags: ['#FloatingMarket', '#Thailand', '#Travel'],
      source: 'following'
    };

    setUserPosts(prev => [newPost, ...prev]);
    toast.success('Location post created! 📍');
  };

  // Combine user posts with feed posts for display
  const allPosts = [...userPosts, ...feedPosts].sort((a, b) => {
    // Sort by timestamp - newest first
    if (a.timestamp === 'just now') return -1;
    if (b.timestamp === 'just now') return 1;
    return 0; // For demo purposes, keep original order for older posts
  });

  const getSourceBadge = (source?: string) => {
    switch (source) {
      case 'trending':
        return <Badge className="bg-orange-500 text-white text-xs">Trending</Badge>;
      case 'sponsored':
        return <Badge variant="secondary" className="text-xs">Sponsored</Badge>;
      case 'interest':
        return <Badge className="bg-purple-500 text-white text-xs">For You</Badge>;
      case 'following':
        return <Badge className="bg-blue-500 text-white text-xs">Following</Badge>;
      default:
        return null;
    }
  };

  const renderPost = (post: Post) => (
    <Card key={post.id} className="mb-4">
      <CardContent className="p-3 sm:p-4 lg:p-6 space-y-3 sm:space-y-4 max-w-full">
        {/* Post Header - Mobile Optimized */}
        <div className="flex items-start gap-3 relative">
          <Avatar className="w-10 h-10 sm:w-11 sm:h-11 cursor-pointer hover:ring-2 hover:ring-ring hover:ring-offset-1 transition-all flex-shrink-0">
            <AvatarImage src={post.author.avatar} />
            <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
              {post.author.name.charAt(0)}
            </AvatarFallback>
          </Avatar>
          
          <div className="flex-1 min-w-0 space-y-1 sm:space-y-2">
            {/* Author Info - Responsive Layout */}
            <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-2">
              <div className="flex items-center gap-2 min-w-0">
                <span className="font-medium cursor-pointer hover:underline text-sm sm:text-base truncate max-w-32 sm:max-w-none">
                  {post.author.name}
                </span>
                {post.author.verified && (
                  <Star className="w-4 h-4 text-yellow-500 fill-current flex-shrink-0" />
                )}
                {post.type === 'live' && post.liveData?.isLive && (
                  <Badge className="bg-red-500 text-white text-xs px-2 py-0.5 flex-shrink-0">LIVE</Badge>
                )}
                {getSourceBadge(post.source)}
              </div>
              
              {/* Add Friend button - Mobile optimized */}
              {!followingUsers.includes(post.author.name) && post.author.type === 'user' && (
                <Button
                  size="sm" 
                  variant="outline"
                  className="h-6 sm:h-7 text-xs px-2 sm:px-3 flex-shrink-0 touch-manipulation"
                  onClick={() => handleAddFriend(post.author.name, post.author.name)}
                >
                  <UserPlus className="w-3 h-3 mr-1" />
                  <span className="hidden sm:inline">Add Friend</span>
                  <span className="sm:hidden">Add</span>
                </Button>
              )}
            </div>
            
            {/* Post Metadata - Responsive Stack */}
            <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3 text-xs sm:text-sm text-muted-foreground">
              <span className="flex-shrink-0">{post.timestamp}</span>
              
              {post.location && (
                <div className="flex items-center gap-1 min-w-0">
                  <MapPin className="w-3 h-3 flex-shrink-0" />
                  <span className="truncate">{post.location}</span>
                </div>
              )}
            </div>
          </div>
          
          {/* More Options Menu - Fixed Position */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8 sm:h-9 sm:w-9 flex-shrink-0 touch-manipulation">
                <MoreVertical className="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48 z-[1000]" sideOffset={8}>
              <DropdownMenuItem onClick={() => handleBookmark(post.id)}>
                <Bookmark className="w-4 h-4 mr-2" />
                {bookmarkedPosts.includes(post.id) ? 'Remove Bookmark' : 'Bookmark'}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleShare(post.id)}>
                <Share className="w-4 h-4 mr-2" />
                Share Post
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Copy className="w-4 h-4 mr-2" />
                Copy Link
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        {/* Post Content - Enhanced Readability */}
        <div className="space-y-3">
          <div className="prose prose-sm sm:prose-base max-w-none">
            <p className="leading-relaxed text-sm sm:text-base text-foreground m-0">
              {post.content.split(' ').map((word, index) => {
                if (word.startsWith('#')) {
                  return (
                    <span
                      key={index}
                      className="text-primary cursor-pointer hover:underline font-medium hover:text-primary/80 transition-colors"
                      onClick={() => handleHashtagClick(word)}
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

        {/* Post Images - Mobile Optimized */}
        {post.images && post.images.length > 0 && (
          <div className="space-y-2">
            <div className="relative rounded-xl overflow-hidden bg-muted">
              <ImageWithFallback
                src={post.images[0]}
                alt="Post content"
                className="w-full h-auto max-h-80 sm:max-h-96 lg:max-h-[28rem] object-cover cursor-pointer hover:scale-105 transition-transform duration-300"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/10 to-transparent pointer-events-none" />
              
              {post.type === 'live' && post.liveData?.isLive && (
                <div className="absolute inset-0 flex items-center justify-center bg-black/20 rounded-xl">
                  <Button 
                    onClick={() => handleJoinLive(post.id)}
                    className="bg-red-500 hover:bg-red-600 text-white touch-manipulation"
                    size="sm"
                  >
                    <Play className="w-4 h-4 mr-2" />
                    <span className="hidden sm:inline">Join Live • {post.liveData.viewers} watching</span>
                    <span className="sm:hidden">Join Live</span>
                  </Button>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Post Actions - Enhanced Mobile Touch Experience */}
        <div className="pt-3 border-t border-border/50">
          <div className="flex items-center justify-between">
            {/* Main Action Buttons */}
            <div className="flex items-center gap-1 sm:gap-3 lg:gap-4">
              <Button
                variant="ghost"
                size="sm"
                className={`h-9 px-2 sm:px-3 hover:bg-red-50 dark:hover:bg-red-950/20 rounded-full transition-colors touch-manipulation min-w-0 ${
                  likedPosts.includes(post.id) ? 'text-red-500 bg-red-50 dark:bg-red-950/20' : 'text-muted-foreground'
                }`}
                onClick={() => handleLike(post.id)}
              >
                <Heart className={`w-4 h-4 mr-1 sm:mr-2 ${likedPosts.includes(post.id) ? 'fill-current' : ''}`} />
                <span className="text-xs sm:text-sm font-medium">
                  {(post.likes + (likedPosts.includes(post.id) ? 1 : 0)) > 999 
                    ? `${Math.floor((post.likes + (likedPosts.includes(post.id) ? 1 : 0)) / 1000)}k` 
                    : (post.likes + (likedPosts.includes(post.id) ? 1 : 0))}
                </span>
              </Button>
              
              <Button
                variant="ghost"
                size="sm"
                className="h-9 px-2 sm:px-3 text-muted-foreground hover:text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-950/20 rounded-full transition-colors touch-manipulation min-w-0"
                onClick={() => handleComment(post.id)}
              >
                <MessageCircle className="w-4 h-4 mr-1 sm:mr-2" />
                <span className="text-xs sm:text-sm font-medium">
                  {post.comments > 999 ? `${Math.floor(post.comments / 1000)}k` : post.comments}
                </span>
              </Button>
              
              <Button
                variant="ghost"
                size="sm"
                className="h-9 px-2 sm:px-3 text-muted-foreground hover:text-green-500 hover:bg-green-50 dark:hover:bg-green-950/20 rounded-full transition-colors touch-manipulation min-w-0"
                onClick={() => handleShare(post.id)}
              >
                <Share className="w-4 h-4 mr-1 sm:mr-2" />
                <span className="text-xs sm:text-sm font-medium hidden sm:inline">
                  {post.shares > 999 ? `${Math.floor(post.shares / 1000)}k` : post.shares}
                </span>
              </Button>
            </div>
            
            {/* Bookmark Button */}
            <Button
              variant="ghost"
              size="sm"
              className={`h-9 w-9 rounded-full transition-colors touch-manipulation ${
                bookmarkedPosts.includes(post.id) 
                  ? 'text-primary bg-primary/10 hover:bg-primary/20' 
                  : 'text-muted-foreground hover:text-primary hover:bg-primary/10'
              }`}
              onClick={() => handleBookmark(post.id)}
            >
              <Bookmark className={`w-4 h-4 ${bookmarkedPosts.includes(post.id) ? 'fill-current' : ''}`} />
            </Button>
          </div>
        </div>

        {/* Comments Section - Enhanced Mobile Layout */}
        {commentsOpen === post.id && (
          <div className="mt-4 pt-4 border-t border-border/50 space-y-4">
            {/* Comment Input */}
            <div className="flex gap-3">
              <Avatar className="w-8 h-8 sm:w-10 sm:h-10 flex-shrink-0">
                <AvatarImage src={user?.avatar} />
                <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
                  {user?.name?.charAt(0) || 'U'}
                </AvatarFallback>
              </Avatar>
              <div className="flex-1 space-y-2">
                <Input
                  placeholder="Write a comment..."
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  onKeyPress={(e) => {
                    if (e.key === 'Enter' && !e.shiftKey) {
                      e.preventDefault();
                      handleAddComment(post.id);
                    }
                  }}
                  className="border-0 bg-muted/50 focus-visible:ring-0 focus-visible:ring-offset-0 text-sm sm:text-base"
                />
                <Button
                  size="sm"
                  className="touch-manipulation"
                  onClick={() => handleAddComment(post.id)}
                  disabled={!newComment.trim()}
                >
                  <Send className="w-3 h-3 mr-1" />
                  Comment
                </Button>
              </div>
            </div>

            {/* Existing Comments */}
            <div className="space-y-4">
              {postComments[post.id]?.map((comment) => (
                <div key={comment.id} className="flex gap-3">
                  <Avatar className="w-8 h-8 sm:w-10 sm:h-10 flex-shrink-0">
                    <AvatarImage src={comment.user.avatar} />
                    <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
                      {comment.user.name.charAt(0)}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 min-w-0">
                    <div className="bg-muted/50 rounded-lg p-3">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-medium text-sm truncate">{comment.user.name}</span>
                        <span className="text-xs text-muted-foreground flex-shrink-0">{comment.timestamp}</span>
                      </div>
                      <p className="text-sm leading-relaxed">{comment.text}</p>
                    </div>
                    <div className="flex items-center gap-4 mt-2">
                      <Button variant="ghost" size="sm" className="h-6 px-2 text-xs text-muted-foreground touch-manipulation">
                        <Heart className="w-3 h-3 mr-1" />
                        {comment.likes}
                      </Button>
                      <Button variant="ghost" size="sm" className="h-6 px-2 text-xs text-muted-foreground touch-manipulation">
                        Reply
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );

  return (
    <div className="h-full flex flex-col">
      {/* Sticky Header with Moments and Create Post */}
      <div className="sticky top-0 z-10 bg-background border-b">
        {/* Moments Section */}
        <div className="p-4 border-b">
          <ScrollArea orientation="horizontal" className="w-full whitespace-nowrap">
            <div className="flex gap-3 pb-2">
              {stories.map((story, index) => (
                <div
                  key={story.id}
                  className="flex flex-col items-center gap-2 cursor-pointer group"
                  onClick={() => handleStoryClick(story, index)}
                >
                  <div className={`relative w-16 h-16 rounded-full p-0.5 ${
                    story.author.name === 'Your Moment' 
                      ? 'bg-gradient-to-tr from-gray-300 to-gray-500' 
                      : story.isViewed 
                        ? 'bg-gray-300' 
                        : 'bg-gradient-to-tr from-orange-400 to-red-600'
                  }`}>
                    <div className="w-full h-full rounded-full bg-background p-0.5">
                      {story.author.name === 'Your Moment' ? (
                        <div className="w-full h-full rounded-full bg-muted flex items-center justify-center">
                          <Plus className="w-6 h-6 text-muted-foreground" />
                        </div>
                      ) : (
                        <ImageWithFallback
                          src={story.preview || story.author.avatar || ''}
                          alt={story.author.name}
                          className="w-full h-full rounded-full object-cover"
                        />
                      )}
                    </div>
                    {story.isLive && (
                      <div className="absolute -bottom-1 -right-1 bg-red-500 text-white text-xs px-1 rounded-full">
                        LIVE
                      </div>
                    )}
                  </div>
                  <span className="text-xs text-center max-w-16 truncate group-hover:text-primary transition-colors">
                    {story.author.name === 'Your Moment' ? 'Your Moment' : story.author.name}
                  </span>
                </div>
              ))}
            </div>
          </ScrollArea>
        </div>

        {/* Create Post Section */}
        <CreatePostSection
          user={user}
          onCreatePost={handleCreatePost}
          newPostText={newPostText}
          setNewPostText={setNewPostText}
          selectedImages={selectedImages}
          setSelectedImages={setSelectedImages}
          postLocation={postLocation}
          setPostLocation={setPostLocation}
          postPrivacy={postPrivacy}
          setPostPrivacy={setPostPrivacy}
          createPostOpen={createPostOpen}
          setCreatePostOpen={setCreatePostOpen}
        />
      </div>

      {/* Main Tabs */}
      <Tabs defaultValue="friends" className="flex-1 flex flex-col">
        <TabsList className="grid w-full grid-cols-4 px-2 sm:px-4 py-2 mx-0">
          <TabsTrigger value="friends" className="flex items-center gap-2">
            <Users className="w-4 h-4" />
            <span className="hidden sm:inline">Friends</span>
          </TabsTrigger>
          <TabsTrigger value="feed" className="flex items-center gap-2">
            <Star className="w-4 h-4" />
            <span className="hidden sm:inline">Feed</span>
          </TabsTrigger>
          <TabsTrigger value="discover" className="flex items-center gap-2">
            <TrendingUp className="w-4 h-4" />
            <span className="hidden sm:inline">Discover</span>
          </TabsTrigger>
          <TabsTrigger value="events" className="flex items-center gap-2">
            <Calendar className="w-4 h-4" />
            <span className="hidden sm:inline">Events</span>
          </TabsTrigger>
        </TabsList>

        <div className="flex-1 overflow-hidden">
          {/* Friends Tab - Friend Posts and Activities */}
          <TabsContent value="friends" className="h-full mt-0">
            <ScrollArea className="h-full">
              <div className="p-4 pb-12 space-y-4">
                {/* Friends Activity Header */}
                <div className="flex items-center justify-between">
                  <h3 className="font-medium flex items-center gap-2">
                    <Users className="w-5 h-5 text-blue-500" />
                    Friends Activity
                  </h3>
                  <Badge variant="secondary">{friends.length} friends</Badge>
                </div>

                {/* Friend Activity Posts */}
                <div className="space-y-4">
                  {/* Sarah's Food Post */}
                  <Card className="mb-4">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-3 mb-3">
                        <Avatar className="cursor-pointer hover:ring-2 hover:ring-ring hover:ring-offset-2 transition-all">
                          <AvatarImage src="https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face" />
                          <AvatarFallback>S</AvatarFallback>
                        </Avatar>
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <span className="font-medium cursor-pointer hover:underline">Sarah Johnson</span>
                            <Badge className="bg-blue-500 text-white text-xs">Friend</Badge>
                          </div>
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <span>25 min ago</span>
                            <span>•</span>
                            <div className="flex items-center gap-1">
                              <MapPin className="w-3 h-3" />
                              <span>Chatuchak Market</span>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div className="mb-3">
                        <p className="text-base leading-relaxed">
                          Just discovered this amazing Pad Thai vendor! 🍜 The uncle has been making this for 40 years and you can taste the experience in every bite. Best 60 baht I've spent today! #PadThai #StreetFood #ChatuchakFinds
                        </p>
                      </div>
                      <ImageWithFallback
                        src="https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=600&h=400&fit=crop"
                        alt="Pad Thai"
                        className="w-full rounded-lg max-h-96 object-cover mb-3"
                      />
                      <div className="flex items-center justify-between pt-2 border-t">
                        <div className="flex items-center gap-4">
                          <Button variant="ghost" size="sm" className="h-8 px-2 text-muted-foreground">
                            <Heart className="w-4 h-4 mr-1" />
                            24
                          </Button>
                          <Button variant="ghost" size="sm" className="h-8 px-2 text-muted-foreground">
                            <MessageCircle className="w-4 h-4 mr-1" />
                            8
                          </Button>
                          <Button variant="ghost" size="sm" className="h-8 px-2 text-muted-foreground">
                            <Share className="w-4 h-4 mr-1" />
                            3
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>

                  {/* Mike's Check-in Activity */}
                  <Card className="mb-4 bg-green-50 border-green-200">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-3 mb-3">
                        <Avatar>
                          <AvatarImage src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face" />
                          <AvatarFallback>M</AvatarFallback>
                        </Avatar>
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <span className="font-medium">Mike Chen</span>
                            <Badge className="bg-green-500 text-white text-xs">
                              <MapPin className="w-3 h-3 mr-1" />
                              Checked in
                            </Badge>
                          </div>
                          <div className="text-sm text-muted-foreground">2 hours ago</div>
                        </div>
                      </div>
                      <p className="text-base mb-2">
                        📍 <span className="font-medium">Wat Pho Temple</span> - Starting the day with morning meditation. The energy here is incredible! 🙏
                      </p>
                      <div className="flex items-center gap-4 text-sm text-green-600">
                        <div className="flex items-center gap-1">
                          <Users className="w-4 h-4" />
                          <span>3 friends also checked in here</span>
                        </div>
                      </div>
                    </CardContent>
                  </Card>

                  {/* Emma's Life Update */}
                  <Card className="mb-4">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-3 mb-3">
                        <Avatar>
                          <AvatarImage src="https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&h=150&fit=crop&crop=face" />
                          <AvatarFallback>E</AvatarFallback>
                        </Avatar>
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <span className="font-medium">Emma Wilson</span>
                            <Badge variant="secondary" className="text-xs">Friend</Badge>
                            <Badge className="bg-orange-500 text-white text-xs">Life Update</Badge>
                          </div>
                          <div className="text-sm text-muted-foreground">4 hours ago</div>
                        </div>
                      </div>
                      <p className="text-base mb-3">
                        Just finished my first Thai cooking class! 👩‍🍳 Learned to make authentic Green Curry from scratch. Can't wait to cook for my family this weekend! The chef was so patient with us beginners 💚
                      </p>
                      <div className="grid grid-cols-2 gap-2 mb-3">
                        <ImageWithFallback
                          src="https://images.unsplash.com/photo-1455619452474-d2be8b1e70cd?w=300&h=200&fit=crop"
                          alt="Cooking class"
                          className="w-full h-24 object-cover rounded-lg"
                        />
                        <ImageWithFallback
                          src="https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=300&h=200&fit=crop"
                          alt="Green curry"
                          className="w-full h-24 object-cover rounded-lg"
                        />
                      </div>
                      <div className="flex items-center justify-between pt-2 border-t">
                        <div className="flex items-center gap-4">
                          <Button variant="ghost" size="sm" className="h-8 px-2 text-muted-foreground">
                            <Heart className="w-4 h-4 mr-1" />
                            18
                          </Button>
                          <Button variant="ghost" size="sm" className="h-8 px-2 text-muted-foreground">
                            <MessageCircle className="w-4 h-4 mr-1" />
                            12
                          </Button>
                          <Button variant="ghost" size="sm" className="h-8 px-2 text-muted-foreground">
                            <Share className="w-4 h-4 mr-1" />
                            5
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>

                  {/* Friend Suggestions */}
                  <Card className="bg-blue-50 border-blue-200">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-2 mb-3">
                        <Users className="w-5 h-5 text-blue-500" />
                        <h4 className="font-medium">Friend Suggestions</h4>
                      </div>
                      <div className="space-y-3">
                        {friends.filter(f => !f.isFollowing).slice(0, 2).map((friend) => (
                          <div key={friend.id} className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                              <Avatar className="w-10 h-10">
                                <AvatarImage src={friend.avatar} />
                                <AvatarFallback>{friend.name.charAt(0)}</AvatarFallback>
                              </Avatar>
                              <div>
                                <div className="font-medium text-sm">{friend.name}</div>
                                <div className="text-xs text-muted-foreground">
                                  {friend.mutualFriends} mutual friends
                                </div>
                              </div>
                            </div>
                            <Button size="sm" onClick={() => handleAddFriend(friend.id, friend.name)}>
                              <UserPlus className="w-3 h-3 mr-1" />
                              Add
                            </Button>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>

                  {/* Active Friends Status */}
                  <Card>
                    <CardContent className="p-4">
                      <div className="flex items-center gap-2 mb-3">
                        <Zap className="w-5 h-5 text-green-500" />
                        <h4 className="font-medium">Friends Online</h4>
                        <Badge className="bg-green-500 text-white text-xs">
                          {friends.filter(f => f.isOnline).length} online
                        </Badge>
                      </div>
                      <div className="flex flex-wrap gap-2">
                        {friends.filter(f => f.isOnline).map((friend) => (
                          <div key={friend.id} className="flex items-center gap-2 bg-muted rounded-full px-3 py-1">
                            <div className="relative">
                              <Avatar className="w-6 h-6">
                                <AvatarImage src={friend.avatar} />
                                <AvatarFallback className="text-xs">{friend.name.charAt(0)}</AvatarFallback>
                              </Avatar>
                              <div className="absolute -bottom-0.5 -right-0.5 w-2 h-2 bg-green-500 border border-background rounded-full"></div>
                            </div>
                            <span className="text-sm font-medium">{friend.name.split(' ')[0]}</span>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </ScrollArea>
          </TabsContent>

          {/* Feed Tab - Facebook-style All Interests */}
          <TabsContent value="feed" className="h-full mt-0">
            <ScrollArea className="h-full">
              <div className="p-4 space-y-4">
                {/* Feed Header */}
                <div className="flex items-center justify-between">
                  <h3 className="font-medium flex items-center gap-2">
                    <Star className="w-5 h-5 text-yellow-500" />
                    Your Feed
                  </h3>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary" className="text-xs">
                      Latest
                    </Badge>
                    <Button variant="outline" size="sm">
                      <Filter className="w-4 h-4 mr-1" />
                      Customize
                    </Button>
                  </div>
                </div>

                {/* Interests Quick Access */}
                <Card className="bg-gradient-to-r from-blue-50 to-purple-50 border-blue-200">
                  <CardContent className="p-4">
                    <div className="flex items-center gap-2 mb-3">
                      <Star className="w-5 h-5 text-yellow-500" />
                      <h4 className="font-medium">Your Interests</h4>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {['🍜 Thai Food', '🏛️ Culture', '🛶 Markets', '🎵 Music', '📱 Tech', '✈️ Travel'].map((interest) => (
                        <Button key={interest} variant="outline" size="sm" className="text-xs">
                          {interest}
                        </Button>
                      ))}
                    </div>
                  </CardContent>
                </Card>

                {/* Algorithm-driven Feed Posts */}
                {feedPosts.map(renderPost)}

                {/* Suggested Content */}
                <Card className="bg-green-50 border-green-200">
                  <CardContent className="p-4">
                    <div className="flex items-center gap-2 mb-3">
                      <TrendingUp className="w-5 h-5 text-green-500" />
                      <h4 className="font-medium">Suggested for You</h4>
                    </div>
                    <div className="space-y-3">
                      <div className="flex items-center gap-3">
                        <ImageWithFallback
                          src="https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=60&h=60&fit=crop"
                          alt="Thai Cooking Group"
                          className="w-12 h-12 object-cover rounded-lg"
                        />
                        <div className="flex-1">
                          <div className="font-medium text-sm">Thai Cooking Enthusiasts</div>
                          <div className="text-xs text-muted-foreground">142K members • Food Group</div>
                        </div>
                        <Button size="sm" variant="outline">
                          <Plus className="w-3 h-3 mr-1" />
                          Join
                        </Button>
                      </div>
                      <div className="flex items-center gap-3">
                        <ImageWithFallback
                          src="https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=60&h=60&fit=crop"
                          alt="Bangkok Markets"
                          className="w-12 h-12 object-cover rounded-lg"
                        />
                        <div className="flex-1">
                          <div className="font-medium text-sm">Bangkok Hidden Gems</div>
                          <div className="text-xs text-muted-foreground">89K members • Local Community</div>
                        </div>
                        <Button size="sm" variant="outline">
                          <Plus className="w-3 h-3 mr-1" />
                          Join
                        </Button>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                {/* Memories/On This Day */}
                <Card className="bg-purple-50 border-purple-200">
                  <CardContent className="p-4">
                    <div className="flex items-center gap-2 mb-3">
                      <Clock className="w-5 h-5 text-purple-500" />
                      <h4 className="font-medium">On This Day</h4>
                    </div>
                    <div className="flex gap-3">
                      <ImageWithFallback
                        src="https://images.unsplash.com/photo-1523050854058-8df90110c9d1?w=60&h=60&fit=crop"
                        alt="Memory"
                        className="w-12 h-12 object-cover rounded-lg"
                      />
                      <div className="flex-1">
                        <div className="text-sm">Your first visit to Songkran Festival</div>
                        <div className="text-xs text-muted-foreground">1 year ago • Bangkok</div>
                        <Button variant="ghost" size="sm" className="h-6 px-2 mt-1 text-xs">
                          <Heart className="w-3 h-3 mr-1" />
                          Share Memory
                        </Button>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </ScrollArea>
          </TabsContent>

          {/* Discover Tab - Enhanced with DiscoverTab Component */}
          <TabsContent value="discover" className="h-full mt-0">
            <DiscoverTab 
              user={user}
              onPostShare={onPostShare}
            />
          </TabsContent>

          {/* Events Tab */}
          <TabsContent value="events" className="h-full mt-0">
            {showEventsTab ? (
              <EventsTab
                user={user}
                onBack={() => setShowEventsTab(false)}
              />
            ) : (
              <div className="p-4 space-y-4">
                {/* Enhanced Events Header */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Calendar className="w-5 h-5 text-chart-1" />
                    <h3 className="font-medium">Local Events</h3>
                    <Badge className="bg-orange-500 text-white text-xs">
                      <Flame className="w-3 h-3 mr-1" />
                      Hot
                    </Badge>
                  </div>
                  <Button 
                    variant="default" 
                    size="sm"
                    onClick={() => setShowEventsTab(true)}
                    className="flex items-center gap-1"
                  >
                    <Eye className="w-4 h-4" />
                    Explore All
                  </Button>
                </div>

                {/* Featured Events Preview */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {/* Trending Music Festival */}
                  <Card className="overflow-hidden relative group cursor-pointer hover:shadow-lg transition-shadow" onClick={() => setShowEventsTab(true)}>
                    <div className="relative">
                      <ImageWithFallback
                        src="https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=800&h=400&fit=crop"
                        alt="Bangkok Electronic Music Festival"
                        className="w-full h-32 object-cover"
                      />
                      <div className="absolute top-2 left-2">
                        <Badge className="bg-red-500 text-white flex items-center gap-1">
                          <Flame className="w-3 h-3" />
                          Trending #1
                        </Badge>
                      </div>
                      <div className="absolute top-2 right-2">
                        <Badge className="bg-white/90 text-gray-900">
                          <Eye className="w-3 h-3 mr-1" />
                          890K views
                        </Badge>
                      </div>
                      <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
                      <div className="absolute bottom-2 left-2 text-white">
                        <div className="flex items-center gap-2 text-sm">
                          <Music className="w-4 h-4" />
                          <span>Music Festival</span>
                          <Users className="w-4 h-4 ml-2" />
                          <span>18,500 going</span>
                        </div>
                      </div>
                    </div>
                    <CardContent className="p-3">
                      <h4 className="font-medium mb-1">Bangkok Electronic Music Festival 2025</h4>
                      <p className="text-sm text-muted-foreground mb-2">Calvin Harris, Armin van Buuren & more</p>
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">March 15-17, 2025</span>
                        <span className="font-medium text-green-600">From ฿2,500</span>
                      </div>
                    </CardContent>
                  </Card>

                  {/* Thai Street Food Festival */}
                  <Card className="overflow-hidden relative group cursor-pointer hover:shadow-lg transition-shadow" onClick={() => setShowEventsTab(true)}>
                    <div className="relative">
                      <ImageWithFallback
                        src="https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=800&h=400&fit=crop"
                        alt="Thai Street Food Championship"
                        className="w-full h-32 object-cover"
                      />
                      <div className="absolute top-2 left-2">
                        <Badge className="bg-green-500 text-white">
                          <Utensils className="w-3 h-3 mr-1" />
                          Food Festival
                        </Badge>
                      </div>
                      <div className="absolute top-2 right-2">
                        <Badge className="bg-white/90 text-gray-900">
                          <Users className="w-3 h-3 mr-1" />
                          12K going
                        </Badge>
                      </div>
                      <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
                      <div className="absolute bottom-2 left-2 text-white">
                        <div className="flex items-center gap-2 text-sm">
                          <Star className="w-4 h-4 fill-current text-yellow-400" />
                          <span>5 friends interested</span>
                        </div>
                      </div>
                    </div>
                    <CardContent className="p-3">
                      <h4 className="font-medium mb-1">Thai Street Food Championship</h4>
                      <p className="text-sm text-muted-foreground mb-2">Master chefs compete & cooking workshops</p>
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-muted-foreground">Feb 28 - Mar 1, 2025</span>
                        <span className="font-medium text-green-600">From ฿500</span>
                      </div>
                    </CardContent>
                  </Card>
                </div>

                {/* Quick Event Categories */}
                <div className="space-y-3">
                  <h4 className="font-medium">Browse by Category</h4>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                    {[
                      { id: 'music', label: 'Music & Festivals', icon: Music, count: '15 events', color: 'bg-purple-500' },
                      { id: 'food', label: 'Food & Markets', icon: Utensils, count: '23 events', color: 'bg-orange-500' },
                      { id: 'cultural', label: 'Cultural Events', icon: Building, count: '8 events', color: 'bg-blue-500' },
                      { id: 'temple', label: 'Temple & Wellness', icon: Star, count: '12 events', color: 'bg-green-500' }
                    ].map((category) => {
                      const Icon = category.icon;
                      return (
                        <Button
                          key={category.id}
                          variant="outline"
                          className="h-auto p-3 flex-col gap-2 hover:shadow-md transition-shadow"
                          onClick={() => setShowEventsTab(true)}
                        >
                          <div className={`w-8 h-8 rounded-full ${category.color} flex items-center justify-center`}>
                            <Icon className="w-4 h-4 text-white" />
                          </div>
                          <div className="text-center">
                            <div className="text-sm font-medium">{category.label}</div>
                            <div className="text-xs text-muted-foreground">{category.count}</div>
                          </div>
                        </Button>
                      );
                    })}
                  </div>
                </div>

                {/* Historical Events Showcase */}
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <h4 className="font-medium">Relive Amazing Moments</h4>
                    <Button variant="ghost" size="sm" onClick={() => setShowEventsTab(true)}>
                      <Clock className="w-4 h-4 mr-1" />
                      View History
                    </Button>
                  </div>
                  <div className="grid grid-cols-2 gap-3">
                    <Card className="overflow-hidden cursor-pointer hover:shadow-md transition-shadow" onClick={() => setShowEventsTab(true)}>
                      <div className="relative">
                        <ImageWithFallback
                          src="https://images.unsplash.com/photo-1523050854058-8df90110c9d1?w=400&h=200&fit=crop"
                          alt="Songkran Music Festival 2024"
                          className="w-full h-20 object-cover"
                        />
                        <Badge className="absolute top-1 left-1 bg-black/70 text-white text-xs">
                          Past Event
                        </Badge>
                      </div>
                      <CardContent className="p-2">
                        <h5 className="text-sm font-medium">Songkran Music Festival 2024</h5>
                        <div className="flex items-center gap-2 text-xs text-muted-foreground mt-1">
                          <Users className="w-3 h-3" />
                          <span>85K attended</span>
                          <Star className="w-3 h-3 text-yellow-500 fill-current" />
                          <span>4.8</span>
                        </div>
                      </CardContent>
                    </Card>
                    
                    <Card className="overflow-hidden cursor-pointer hover:shadow-md transition-shadow" onClick={() => setShowEventsTab(true)}>
                      <div className="relative">
                        <ImageWithFallback
                          src="https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=400&h=200&fit=crop"
                          alt="Floating Market Food Festival 2024"
                          className="w-full h-20 object-cover"
                        />
                        <Badge className="absolute top-1 left-1 bg-black/70 text-white text-xs">
                          Past Event
                        </Badge>
                      </div>
                      <CardContent className="p-2">
                        <h5 className="text-sm font-medium">Floating Market Food Festival</h5>
                        <div className="flex items-center gap-2 text-xs text-muted-foreground mt-1">
                          <Users className="w-3 h-3" />
                          <span>25K attended</span>
                          <Star className="w-3 h-3 text-yellow-500 fill-current" />
                          <span>4.6</span>
                        </div>
                      </CardContent>
                    </Card>
                  </div>
                </div>

                {/* Call to Action */}
                <Card className="bg-gradient-to-r from-chart-1/10 to-chart-2/10 border-chart-1/20">
                  <CardContent className="p-4 text-center">
                    <Calendar className="w-8 h-8 text-chart-1 mx-auto mb-2" />
                    <h4 className="font-medium mb-1">Discover Amazing Events</h4>
                    <p className="text-sm text-muted-foreground mb-3">
                      Explore upcoming festivals, food markets, and cultural events in Thailand
                    </p>
                    <Button onClick={() => setShowEventsTab(true)} className="w-full">
                      <Sparkles className="w-4 h-4 mr-1" />
                      Explore All Events
                    </Button>
                  </CardContent>
                </Card>
              </div>
            )}
          </TabsContent>
        </div>
      </Tabs>

      {/* Moment Viewer Dialog */}
      {viewingStory && (
        <Dialog open={!!viewingStory} onOpenChange={() => setViewingStory(null)}>
          <DialogContent className="max-w-md p-0 bg-black border-0">
            <div className="relative aspect-[9/16] bg-black">
              {/* Progress Bars */}
              <div className="absolute top-4 left-4 right-4 flex gap-1 z-10">
                {viewingStory.media?.map((_, index) => (
                  <div key={index} className="flex-1 h-0.5 bg-white/30 rounded">
                    <div 
                      className="h-full bg-white rounded transition-all duration-100"
                      style={{ 
                        width: index === storyMediaIndex ? `${storyProgress}%` : index < storyMediaIndex ? '100%' : '0%'
                      }}
                    />
                  </div>
                ))}
              </div>

              {/* Story Content */}
              <div className="absolute inset-0">
                {viewingStory.media && (
                  <ImageWithFallback
                    src={viewingStory.media[storyMediaIndex]?.url || ''}
                    alt="Story content"
                    className="w-full h-full object-cover"
                  />
                )}
                
                {/* Story Text Overlay */}
                <div className="absolute bottom-4 left-4 right-4 text-white">
                  <div className="flex items-center gap-2 mb-2">
                    <Avatar className="w-8 h-8 border-2 border-white">
                      <AvatarImage src={viewingStory.author.avatar} />
                      <AvatarFallback>{viewingStory.author.name.charAt(0)}</AvatarFallback>
                    </Avatar>
                    <div>
                      <div className="font-medium text-sm">{viewingStory.author.name}</div>
                      <div className="text-xs text-white/80">{viewingStory.timestamp}</div>
                    </div>
                  </div>
                  {viewingStory.content && (
                    <p className="text-sm bg-black/50 rounded p-2">{viewingStory.content}</p>
                  )}
                </div>

                {/* Navigation */}
                <button 
                  className="absolute left-0 top-0 w-1/3 h-full"
                  onClick={() => {
                    if (storyMediaIndex > 0) {
                      setStoryMediaIndex(prev => prev - 1);
                      setStoryProgress(0);
                    }
                  }}
                />
                <button 
                  className="absolute right-0 top-0 w-1/3 h-full"
                  onClick={() => {
                    if (viewingStory.media && storyMediaIndex < viewingStory.media.length - 1) {
                      setStoryMediaIndex(prev => prev + 1);
                      setStoryProgress(0);
                    } else {
                      setViewingStory(null);
                      setStoryMediaIndex(0);
                    }
                  }}
                />
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}

      {/* Create Story Dialog */}
      <Dialog open={createStoryOpen} onOpenChange={setCreateStoryOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Create Your Moment</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="aspect-[9/16] bg-muted rounded-lg flex items-center justify-center">
              <Camera className="w-12 h-12 text-muted-foreground" />
            </div>
            <Input
              placeholder="Add text to your moment..."
              value={storyText}
              onChange={(e) => setStoryText(e.target.value)}
            />
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => setCreateStoryOpen(false)}>
                Cancel
              </Button>
              <Button onClick={handleCreateStory} className="flex-1">
                <Camera className="w-4 h-4 mr-2" />
                Share Moment
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}