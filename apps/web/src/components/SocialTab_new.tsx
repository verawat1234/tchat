import React, { useState } from 'react';
import { Heart, MessageCircle, Share, MoreVertical, Plus, Camera, MapPin, Clock, Users, Zap, Star, Play, Send, Bookmark, UserPlus, UserCheck, Hash, Copy, ChevronDown, ChevronUp, TrendingUp, Search, Eye, Filter, X, Globe, Utensils, Building, Calendar } from 'lucide-react';
import { CreatePostSection } from './CreatePostSection';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Input } from './ui/input';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './ui/dialog';
import { Separator } from './ui/separator';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";

interface SocialTabProps {
  user: any;
  onLiveStreamClick?: (streamId: string) => void;
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
}

interface Story {
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
  }[];
}

export function SocialTab({ user, onLiveStreamClick }: SocialTabProps) {
  const [likedPosts, setLikedPosts] = useState<string[]>([]);
  const [bookmarkedPosts, setBookmarkedPosts] = useState<string[]>([]);
  const [followingUsers, setFollowingUsers] = useState<string[]>([]);
  const [commentsOpen, setCommentsOpen] = useState<string | null>(null);
  const [shareOpen, setShareOpen] = useState<string | null>(null);
  const [newComment, setNewComment] = useState('');
  const [postComments, setPostComments] = useState<{ [key: string]: Comment[] }>({});
  const [expandedPosts, setExpandedPosts] = useState<string[]>([]);
  const [viewingStory, setViewingStory] = useState<Story | null>(null);
  const [storyIndex, setStoryIndex] = useState(0);
  const [storyProgress, setStoryProgress] = useState(0);
  const [discoverCategory, setDiscoverCategory] = useState<string>('trending');
  const [discoverSearch, setDiscoverSearch] = useState('');
  
  // Create Post state
  const [createPostOpen, setCreatePostOpen] = useState(false);
  const [newPostText, setNewPostText] = useState('');
  const [selectedImages, setSelectedImages] = useState<string[]>([]);
  const [postLocation, setPostLocation] = useState('');
  const [postPrivacy, setPostPrivacy] = useState<'public' | 'friends' | 'private'>('public');

  // Mock data
  const stories: Story[] = [
    {
      id: '1',
      author: { name: 'Your Story', avatar: user?.avatar || '' },
      preview: '',
      isViewed: false,
      content: 'Add to your story',
      timestamp: 'now'
    },
    {
      id: '2',
      author: { name: 'Sarah Johnson', avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face' },
      preview: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      isViewed: false,
      content: 'Amazing Pad Thai at Chatuchak Market! 🍜✨',
      timestamp: '2h ago',
      media: [{
        type: 'image',
        url: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'
      }]
    },
    // ... other stories remain the same
  ];

  // Create Post handlers
  const handleCreatePost = () => {
    if (newPostText.trim()) {
      // In real app, would send post to backend
      console.log('Creating post:', {
        text: newPostText,
        images: selectedImages,
        location: postLocation,
        privacy: postPrivacy
      });
      
      // Reset form
      setNewPostText('');
      setSelectedImages([]);
      setPostLocation('');
      setPostPrivacy('public');
      setCreatePostOpen(false);
      
      toast.success('Post shared successfully! 🎉');
    }
  };

  const handleImageSelect = () => {
    // In real app, would open image picker
    const demoImages = [
      'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGFpJTIwc3RyZWV0JTIwZm9vZCUyMGNvb2tpbmd8ZW58MXx8fHwxNzU4Mzk0NTE3fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral',
      'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzb3V0aGVhc3QlMjBhc2lhJTIwbWFya2V0JTIwdmVuZG9yfGVufDF8fHx8MTc1ODM5NDUxNXww&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'
    ];
    const randomImage = demoImages[Math.floor(Math.random() * demoImages.length)];
    setSelectedImages(prev => [...prev, randomImage]);
    toast.success('Image added!');
  };

  const removeImage = (index: number) => {
    setSelectedImages(prev => prev.filter((_, i) => i !== index));
    toast.success('Image removed');
  };

  // Other handlers remain the same...
  const handleLike = (postId: string) => {
    if (likedPosts.includes(postId)) {
      setLikedPosts(likedPosts.filter(id => id !== postId));
    } else {
      setLikedPosts([...likedPosts, postId]);
      toast.success('Post liked!');
    }
  };

  return (
    <div className="h-full flex flex-col">
      {/* Stories Section */}
      <div className="border-b border-border bg-card">
        <div className="p-4">
          <ScrollArea className="w-full whitespace-nowrap">
            <div className="flex space-x-4">
              {stories.map((story) => (
                <div key={story.id} className="flex flex-col items-center space-y-2 flex-shrink-0">
                  <div className={`relative ${story.isLive ? 'bg-gradient-to-r from-red-500 to-pink-500' : story.isViewed ? 'bg-gray-300' : 'bg-gradient-to-r from-purple-500 to-pink-500'} p-1 rounded-full`}>
                    <Avatar className="w-16 h-16 border-2 border-white">
                      <AvatarImage src={story.preview || story.author.avatar} />
                      <AvatarFallback>
                        {story.author.name === 'Your Story' ? (
                          <Plus className="w-6 h-6" />
                        ) : (
                          story.author.name.charAt(0)
                        )}
                      </AvatarFallback>
                    </Avatar>
                    {story.isLive && (
                      <Badge className="absolute -bottom-1 left-1/2 transform -translate-x-1/2 bg-red-500 text-white text-xs">
                        LIVE
                      </Badge>
                    )}
                  </div>
                  <span className="text-xs text-center w-20 truncate">
                    {story.author.name}
                  </span>
                </div>
              ))}
            </div>
          </ScrollArea>
        </div>
      </div>

      {/* Create Post Section - Placed here between Stories and Tabs */}
      <CreatePostSection
        user={user}
        createPostOpen={createPostOpen}
        setCreatePostOpen={setCreatePostOpen}
        newPostText={newPostText}
        setNewPostText={setNewPostText}
        selectedImages={selectedImages}
        setSelectedImages={setSelectedImages}
        postLocation={postLocation}
        setPostLocation={setPostLocation}
        postPrivacy={postPrivacy}
        setPostPrivacy={setPostPrivacy}
        handleCreatePost={handleCreatePost}
        handleImageSelect={handleImageSelect}
        removeImage={removeImage}
      />

      {/* Rest of the Social Tab content with tabs */}
      <Tabs defaultValue="feed" className="flex-1 flex flex-col">
        <TabsList className="grid w-full grid-cols-4 bg-muted m-2">
          <TabsTrigger value="feed">Feed</TabsTrigger>
          <TabsTrigger value="friends">Friends</TabsTrigger>
          <TabsTrigger value="discover">Discover</TabsTrigger>
          <TabsTrigger value="events">Events</TabsTrigger>
        </TabsList>

        <TabsContent value="feed" className="flex-1 overflow-hidden">
          <ScrollArea className="h-full">
            <div className="space-y-4 p-4">
              {/* Sample posts content */}
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-3 mb-3">
                    <Avatar>
                      <AvatarImage src="https://images.unsplash.com/photo-1494790108755-2616b612b820?w=150&h=150&fit=crop&crop=face" />
                      <AvatarFallback>U</AvatarFallback>
                    </Avatar>
                    <div>
                      <p className="font-medium">Sample User</p>
                      <p className="text-sm text-muted-foreground">2 hours ago</p>
                    </div>
                  </div>
                  <p className="mb-3">Just tried the amazing Pad Thai from the local vendor! 🍜</p>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <Button variant="ghost" size="sm" className="gap-2">
                        <Heart className="w-4 h-4" />
                        <span>24</span>
                      </Button>
                      <Button variant="ghost" size="sm" className="gap-2">
                        <MessageCircle className="w-4 h-4" />
                        <span>5</span>
                      </Button>
                      <Button variant="ghost" size="sm" className="gap-2">
                        <Share className="w-4 h-4" />
                        <span>2</span>
                      </Button>
                    </div>
                    <Button variant="ghost" size="sm">
                      <Bookmark className="w-4 h-4" />
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </div>
          </ScrollArea>
        </TabsContent>

        <TabsContent value="friends" className="flex-1 overflow-hidden">
          <ScrollArea className="h-full">
            <div className="p-4">
              <p className="text-muted-foreground">Friends tab content...</p>
            </div>
          </ScrollArea>
        </TabsContent>

        <TabsContent value="discover" className="flex-1 overflow-hidden">
          <ScrollArea className="h-full">
            <div className="p-4">
              <p className="text-muted-foreground">Discover tab content...</p>
            </div>
          </ScrollArea>
        </TabsContent>

        <TabsContent value="events" className="flex-1 overflow-hidden">
          <ScrollArea className="h-full">
            <div className="p-4">
              <p className="text-muted-foreground">Events tab content...</p>
            </div>
          </ScrollArea>
        </TabsContent>
      </Tabs>
    </div>
  );
}