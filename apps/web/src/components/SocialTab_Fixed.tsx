import React, { useState } from 'react';
import { CreatePostSection } from './CreatePostSection';
import { Card, CardContent } from './ui/card';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { Heart, MessageCircle, Share, Bookmark, Star, MapPin } from 'lucide-react';
import { toast } from "sonner";

interface SocialTabProps {
  user: any;
  onLiveStreamClick?: (streamId: string) => void;
  onPostShare?: (postId: string, postData: any) => void;
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
  location?: string;
  tags?: string[];
  type: 'text' | 'image' | 'live' | 'product';
  source?: 'following' | 'trending' | 'interest' | 'sponsored';
}

export function SocialTab({ user, onLiveStreamClick, onPostShare }: SocialTabProps) {
  const [userPosts, setUserPosts] = useState<Post[]>([]);
  const [likedPosts, setLikedPosts] = useState<string[]>([]);
  const [bookmarkedPosts, setBookmarkedPosts] = useState<string[]>([]);

  // Sample feed posts
  const feedPosts: Post[] = [
    {
      id: '1',
      author: {
        name: 'Thai Food Explorer',
        avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&h=150&fit=crop&crop=face',
        verified: true,
        type: 'channel'
      },
      content: 'TRENDING: Secret street food spots only locals know about! 🏮🍜 This hidden gem in Chinatown serves the most authentic Tom Yum I\'ve ever tasted.',
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
      id: '2',
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

  // Real post creation handlers
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

    setUserPosts(prev => [newPost, ...prev]);
    toast.success('Post created successfully! 🎉');
  };

  const handleCreatePhotoPost = () => {
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

  // Combine and sort posts
  const allPosts = [...userPosts, ...feedPosts].sort((a, b) => {
    if (a.timestamp === 'just now') return -1;
    if (b.timestamp === 'just now') return 1;
    return 0;
  });

  const handleLike = (postId: string) => {
    if (likedPosts.includes(postId)) {
      setLikedPosts(likedPosts.filter(id => id !== postId));
    } else {
      setLikedPosts([...likedPosts, postId]);
      toast.success('Post liked!');
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

  const handleShare = (postId: string) => {
    const post = allPosts.find(p => p.id === postId);
    if (post && onPostShare) {
      onPostShare(postId, post);
    } else {
      navigator.clipboard.writeText(`https://telegram-sea.app/post/${postId}`);
      toast.success('Post link copied to clipboard!');
    }
  };

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
      <CardContent className="p-4 space-y-4">
        {/* Post Header */}
        <div className="flex items-start gap-3">
          <Avatar className="w-10 h-10 cursor-pointer">
            <AvatarImage src={post.author.avatar} />
            <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
              {post.author.name.charAt(0)}
            </AvatarFallback>
          </Avatar>
          
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <span className="font-medium cursor-pointer hover:underline text-sm">
                {post.author.name}
              </span>
              {post.author.verified && (
                <Star className="w-4 h-4 text-yellow-500 fill-current" />
              )}
              {getSourceBadge(post.source)}
            </div>
            
            <div className="flex items-center gap-3 text-xs text-muted-foreground">
              <span>{post.timestamp}</span>
              {post.location && (
                <div className="flex items-center gap-1">
                  <MapPin className="w-3 h-3" />
                  <span className="truncate">{post.location}</span>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Post Content */}
        <div className="space-y-3">
          <p className="leading-relaxed text-sm text-foreground">
            {post.content}
          </p>
        </div>

        {/* Post Images */}
        {post.images && post.images.length > 0 && (
          <div className="space-y-2">
            <div className="relative rounded-xl overflow-hidden bg-muted">
              <ImageWithFallback
                src={post.images[0]}
                alt="Post content"
                className="w-full h-auto max-h-80 object-cover"
              />
            </div>
          </div>
        )}

        {/* Post Actions */}
        <div className="pt-3 border-t border-border/50">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Button
                variant="ghost"
                size="sm"
                className={`h-9 px-3 rounded-full transition-colors ${
                  likedPosts.includes(post.id) ? 'text-red-500 bg-red-50' : 'text-muted-foreground'
                }`}
                onClick={() => handleLike(post.id)}
              >
                <Heart className={`w-4 h-4 mr-2 ${likedPosts.includes(post.id) ? 'fill-current' : ''}`} />
                <span className="text-sm font-medium">
                  {post.likes + (likedPosts.includes(post.id) ? 1 : 0)}
                </span>
              </Button>
              
              <Button
                variant="ghost"
                size="sm"
                className="h-9 px-3 text-muted-foreground rounded-full"
              >
                <MessageCircle className="w-4 h-4 mr-2" />
                <span className="text-sm font-medium">{post.comments}</span>
              </Button>
              
              <Button
                variant="ghost"
                size="sm"
                className="h-9 px-3 text-muted-foreground rounded-full"
                onClick={() => handleShare(post.id)}
              >
                <Share className="w-4 h-4 mr-2" />
                <span className="text-sm font-medium">{post.shares}</span>
              </Button>
            </div>
            
            <Button
              variant="ghost"
              size="sm"
              className={`h-9 w-9 rounded-full ${
                bookmarkedPosts.includes(post.id) 
                  ? 'text-primary bg-primary/10' 
                  : 'text-muted-foreground'
              }`}
              onClick={() => handleBookmark(post.id)}
            >
              <Bookmark className={`w-4 h-4 ${bookmarkedPosts.includes(post.id) ? 'fill-current' : ''}`} />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );

  return (
    <div className="w-full h-full flex flex-col min-h-0 min-w-[320px] max-w-full">
      {/* Create Post Section - Fixed at top */}
      <div className="flex-shrink-0">
        <CreatePostSection
          user={user}
          onCreateTextPost={handleCreateTextPost}
          onCreatePhotoPost={handleCreatePhotoPost}
          onCreateGalleryPost={handleCreateGalleryPost}
          onCreateLocationPost={handleCreateLocationPost}
        />
      </div>

      {/* Posts Feed - Scrollable */}
      <div className="flex-1 overflow-y-auto mobile-scroll scrollbar-hide">
        <div className="px-4 pb-4 space-y-4">
          {allPosts.map(renderPost)}
        </div>
      </div>
    </div>
  );
}