/**
 * SocialFeed Component
 *
 * Main social feed component that displays personalized posts with real-time updates.
 * Features infinite scroll, optimistic updates, and comprehensive social interactions.
 */

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { Heart, MessageCircle, Share, MoreVertical, Bookmark, UserPlus, MapPin, Star, Calendar, TrendingUp, Filter, AlertTriangle } from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Input } from './ui/input';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './ui/dialog';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { toast } from "sonner";
import {
  useGetSocialFeedQuery,
  useAddReactionMutation,
  useRemoveReactionMutation,
  useCreateCommentMutation,
  useShareContentMutation,
  transformPostToLegacy,
  getSocialDisplayName,
  canEditPost,
  getReactionEmoji,
} from '../services/socialApi';
import type { Post, Comment, SocialFeedRequest } from '../types/social';

interface SocialFeedProps {
  user: any;
  algorithm?: 'chronological' | 'personalized' | 'trending';
  region?: 'TH' | 'SG' | 'ID' | 'MY' | 'PH' | 'VN';
  onPostClick?: (postId: string) => void;
  onUserClick?: (userId: string) => void;
  className?: string;
}

interface PostCardProps {
  post: Post;
  currentUser: any;
  onLike: (postId: string) => Promise<void>;
  onComment: (postId: string, content: string) => Promise<void>;
  onShare: (postId: string) => Promise<void>;
  onBookmark: (postId: string) => void;
  onPostClick?: (postId: string) => void;
  onUserClick?: (userId: string) => void;
}

const PostCard: React.FC<PostCardProps> = ({
  post,
  currentUser,
  onLike,
  onComment,
  onShare,
  onBookmark,
  onPostClick,
  onUserClick,
}) => {
  const [showComments, setShowComments] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [isLiking, setIsLiking] = useState(false);
  const [isSharing, setIsSharing] = useState(false);
  const [bookmarkedPosts, setBookmarkedPosts] = useState<string[]>([]);

  const handleLike = async () => {
    if (isLiking) return;
    setIsLiking(true);
    try {
      await onLike(post.id);
    } finally {
      setIsLiking(false);
    }
  };

  const handleAddComment = async () => {
    if (!newComment.trim()) return;
    try {
      await onComment(post.id, newComment);
      setNewComment('');
      toast.success('Comment added!');
    } catch (error) {
      toast.error('Failed to add comment');
    }
  };

  const handleShare = async () => {
    if (isSharing) return;
    setIsSharing(true);
    try {
      await onShare(post.id);
    } finally {
      setIsSharing(false);
    }
  };

  const handleBookmark = () => {
    const isBookmarked = bookmarkedPosts.includes(post.id);
    if (isBookmarked) {
      setBookmarkedPosts(prev => prev.filter(id => id !== post.id));
      toast.success('Removed from bookmarks');
    } else {
      setBookmarkedPosts(prev => [...prev, post.id]);
      toast.success('Added to bookmarks');
    }
    onBookmark(post.id);
  };

  const formatTimestamp = (timestamp: string) => {
    const now = new Date();
    const date = new Date(timestamp);
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / (1000 * 60));
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffMins < 1) return 'just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;

    return date.toLocaleDateString();
  };

  const getSourceBadge = () => {
    if (post.isTrending) {
      return <Badge className="bg-orange-500 text-white text-xs">Trending</Badge>;
    }
    if (post.isPinned) {
      return <Badge className="bg-blue-500 text-white text-xs">Pinned</Badge>;
    }
    return null;
  };

  return (
    <Card className="mb-4 hover:shadow-md transition-shadow">
      <CardContent className="p-3 sm:p-4 lg:p-6 space-y-3 sm:space-y-4">
        {/* Post Header */}
        <div className="flex items-start gap-3 relative">
          <Avatar
            className="w-10 h-10 sm:w-11 sm:h-11 cursor-pointer hover:ring-2 hover:ring-ring hover:ring-offset-1 transition-all flex-shrink-0"
            onClick={() => onUserClick?.(post.authorId)}
          >
            <AvatarImage src={post.author?.avatar} />
            <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
              {getSocialDisplayName(post.author || {} as any).charAt(0)}
            </AvatarFallback>
          </Avatar>

          <div className="flex-1 min-w-0 space-y-1 sm:space-y-2">
            <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-2">
              <div className="flex items-center gap-2 min-w-0">
                <span
                  className="font-medium cursor-pointer hover:underline text-sm sm:text-base truncate max-w-32 sm:max-w-none"
                  onClick={() => onUserClick?.(post.authorId)}
                >
                  {getSocialDisplayName(post.author || {} as any)}
                </span>
                {post.author?.isSocialVerified && (
                  <Star className="w-4 h-4 text-yellow-500 fill-current flex-shrink-0" />
                )}
                {getSourceBadge()}
              </div>
            </div>

            <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3 text-xs sm:text-sm text-muted-foreground">
              <span className="flex-shrink-0">{formatTimestamp(post.createdAt)}</span>

              {post.metadata?.location && (
                <div className="flex items-center gap-1 min-w-0">
                  <MapPin className="w-3 h-3 flex-shrink-0" />
                  <span className="truncate">{post.metadata.location}</span>
                </div>
              )}

              {post.isEdited && (
                <span className="text-xs text-muted-foreground">edited</span>
              )}
            </div>
          </div>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8 sm:h-9 sm:w-9 flex-shrink-0 touch-manipulation">
                <MoreVertical className="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48 z-[1000]" sideOffset={8}>
              <DropdownMenuItem onClick={handleBookmark}>
                <Bookmark className="w-4 h-4 mr-2" />
                {bookmarkedPosts.includes(post.id) ? 'Remove Bookmark' : 'Bookmark'}
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleShare}>
                <Share className="w-4 h-4 mr-2" />
                Share Post
              </DropdownMenuItem>
              {canEditPost(post, currentUser?.id) && (
                <DropdownMenuItem onClick={() => onPostClick?.(post.id)}>
                  <Star className="w-4 h-4 mr-2" />
                  Edit Post
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        {/* Post Content */}
        <div className="space-y-3">
          <div className="prose prose-sm sm:prose-base max-w-none">
            <p className="leading-relaxed text-sm sm:text-base text-foreground m-0">
              {post.content.split(' ').map((word, index) => {
                if (word.startsWith('#')) {
                  return (
                    <span
                      key={index}
                      className="text-primary cursor-pointer hover:underline font-medium hover:text-primary/80 transition-colors"
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

        {/* Post Images */}
        {post.mediaUrls && post.mediaUrls.length > 0 && (
          <div className="space-y-2">
            <div className="relative rounded-xl overflow-hidden bg-muted">
              <ImageWithFallback
                src={post.mediaUrls[0]}
                alt="Post content"
                className="w-full h-auto max-h-80 sm:max-h-96 lg:max-h-[28rem] object-cover cursor-pointer hover:scale-105 transition-transform duration-300"
                onClick={() => onPostClick?.(post.id)}
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/10 to-transparent pointer-events-none" />
            </div>
          </div>
        )}

        {/* Post Tags */}
        {post.tags && post.tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {post.tags.map((tag, index) => (
              <Badge key={index} variant="secondary" className="text-xs cursor-pointer hover:bg-primary/20">
                #{tag}
              </Badge>
            ))}
          </div>
        )}

        {/* Post Actions */}
        <div className="pt-3 border-t border-border/50">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1 sm:gap-3 lg:gap-4">
              <Button
                variant="ghost"
                size="sm"
                className={`h-9 px-2 sm:px-3 hover:bg-red-50 dark:hover:bg-red-950/20 rounded-full transition-colors touch-manipulation min-w-0 ${
                  post.isLiked ? 'text-red-500 bg-red-50 dark:bg-red-950/20' : 'text-muted-foreground'
                }`}
                onClick={handleLike}
                disabled={isLiking}
              >
                <Heart className={`w-4 h-4 mr-1 sm:mr-2 ${post.isLiked ? 'fill-current' : ''}`} />
                <span className="text-xs sm:text-sm font-medium">
                  {post.likesCount > 999 ? `${Math.floor(post.likesCount / 1000)}k` : post.likesCount}
                </span>
              </Button>

              <Button
                variant="ghost"
                size="sm"
                className="h-9 px-2 sm:px-3 text-muted-foreground hover:text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-950/20 rounded-full transition-colors touch-manipulation min-w-0"
                onClick={() => setShowComments(!showComments)}
              >
                <MessageCircle className="w-4 h-4 mr-1 sm:mr-2" />
                <span className="text-xs sm:text-sm font-medium">
                  {post.commentsCount > 999 ? `${Math.floor(post.commentsCount / 1000)}k` : post.commentsCount}
                </span>
              </Button>

              <Button
                variant="ghost"
                size="sm"
                className="h-9 px-2 sm:px-3 text-muted-foreground hover:text-green-500 hover:bg-green-50 dark:hover:bg-green-950/20 rounded-full transition-colors touch-manipulation min-w-0"
                onClick={handleShare}
                disabled={isSharing}
              >
                <Share className="w-4 h-4 mr-1 sm:mr-2" />
                <span className="text-xs sm:text-sm font-medium hidden sm:inline">
                  {post.sharesCount > 999 ? `${Math.floor(post.sharesCount / 1000)}k` : post.sharesCount}
                </span>
              </Button>
            </div>

            <Button
              variant="ghost"
              size="sm"
              className={`h-9 w-9 rounded-full transition-colors touch-manipulation ${
                bookmarkedPosts.includes(post.id)
                  ? 'text-primary bg-primary/10 hover:bg-primary/20'
                  : 'text-muted-foreground hover:text-primary hover:bg-primary/10'
              }`}
              onClick={handleBookmark}
            >
              <Bookmark className={`w-4 h-4 ${bookmarkedPosts.includes(post.id) ? 'fill-current' : ''}`} />
            </Button>
          </div>
        </div>

        {/* Comments Section */}
        {showComments && (
          <div className="mt-4 pt-4 border-t border-border/50 space-y-4">
            <div className="flex gap-3">
              <Avatar className="w-8 h-8 sm:w-10 sm:h-10 flex-shrink-0">
                <AvatarImage src={currentUser?.avatar} />
                <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
                  {currentUser?.name?.charAt(0) || 'U'}
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
                      handleAddComment();
                    }
                  }}
                  className="border-0 bg-muted/50 focus-visible:ring-0 focus-visible:ring-offset-0 text-sm sm:text-base"
                />
                <Button
                  size="sm"
                  className="touch-manipulation"
                  onClick={handleAddComment}
                  disabled={!newComment.trim()}
                >
                  <MessageCircle className="w-3 h-3 mr-1" />
                  Comment
                </Button>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export const SocialFeed: React.FC<SocialFeedProps> = ({
  user,
  algorithm = 'personalized',
  region,
  onPostClick,
  onUserClick,
  className = '',
}) => {
  const [feedParams, setFeedParams] = useState<SocialFeedRequest>({
    algorithm,
    limit: 20,
    region,
  });

  // RTK Query for social feed
  const {
    data: feedData,
    isLoading: feedLoading,
    error: feedError,
    refetch: refetchFeed,
  } = useGetSocialFeedQuery(feedParams, {
    skip: !user?.id,
    pollingInterval: 60000, // Poll every minute for fresh content
  });

  // Mutations
  const [addReaction] = useAddReactionMutation();
  const [removeReaction] = useRemoveReactionMutation();
  const [createComment] = useCreateCommentMutation();
  const [shareContent] = useShareContentMutation();

  // Handlers
  const handleLike = useCallback(async (postId: string) => {
    try {
      const post = feedData?.posts.find(p => p.id === postId);
      if (!post) return;

      if (post.isLiked) {
        await removeReaction({ targetId: postId, targetType: 'post' }).unwrap();
        toast.success('Reaction removed');
      } else {
        await addReaction({ targetId: postId, targetType: 'post', type: 'like' }).unwrap();
        toast.success('Post liked!');
      }
    } catch (error) {
      toast.error('Failed to update reaction');
      console.error('Error updating reaction:', error);
    }
  }, [feedData?.posts, addReaction, removeReaction]);

  const handleComment = useCallback(async (postId: string, content: string) => {
    try {
      await createComment({ postId, content }).unwrap();
    } catch (error) {
      throw new Error('Failed to add comment');
    }
  }, [createComment]);

  const handleShare = useCallback(async (postId: string) => {
    try {
      await shareContent({
        contentId: postId,
        contentType: 'post',
        platform: 'internal',
      }).unwrap();

      // Copy to clipboard as well
      await navigator.clipboard.writeText(`${window.location.origin}/post/${postId}`);
      toast.success('Post shared and link copied!');
    } catch (error) {
      // Fallback to clipboard only
      try {
        await navigator.clipboard.writeText(`${window.location.origin}/post/${postId}`);
        toast.success('Post link copied to clipboard!');
      } catch {
        toast.error('Failed to share post');
      }
    }
  }, [shareContent]);

  const handleBookmark = useCallback((postId: string) => {
    // This could be implemented with a bookmark API endpoint
    // For now, we'll just show the UI feedback
  }, []);

  const handleAlgorithmChange = (newAlgorithm: typeof algorithm) => {
    setFeedParams(prev => ({ ...prev, algorithm: newAlgorithm }));
  };

  // Loading state
  if (feedLoading && !feedData) {
    return (
      <div className={`space-y-4 ${className}`}>
        {[...Array(3)].map((_, i) => (
          <Card key={i} className="mb-4">
            <CardContent className="p-4 space-y-4">
              <div className="flex items-start gap-3">
                <div className="w-10 h-10 bg-muted rounded-full animate-pulse" />
                <div className="flex-1 space-y-2">
                  <div className="h-4 bg-muted rounded animate-pulse w-1/3" />
                  <div className="h-3 bg-muted rounded animate-pulse w-1/4" />
                </div>
              </div>
              <div className="space-y-2">
                <div className="h-4 bg-muted rounded animate-pulse" />
                <div className="h-4 bg-muted rounded animate-pulse w-3/4" />
              </div>
              <div className="h-32 bg-muted rounded animate-pulse" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  // Error state
  if (feedError) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <AlertTriangle className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
        <h3 className="font-medium mb-2">Unable to load feed</h3>
        <p className="text-sm text-muted-foreground mb-4">
          There was an issue loading your social feed. Please try again.
        </p>
        <Button onClick={() => refetchFeed()} variant="outline">
          <TrendingUp className="w-4 h-4 mr-2" />
          Retry
        </Button>
      </div>
    );
  }

  // Empty state
  if (!feedData?.posts || feedData.posts.length === 0) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <Calendar className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
        <h3 className="font-medium mb-2">No posts yet</h3>
        <p className="text-sm text-muted-foreground mb-4">
          Follow some users or join communities to see posts in your feed.
        </p>
        <Button onClick={() => refetchFeed()} variant="outline">
          <TrendingUp className="w-4 h-4 mr-2" />
          Refresh Feed
        </Button>
      </div>
    );
  }

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Feed Controls */}
      <div className="flex items-center justify-between p-4 bg-muted/50 rounded-lg">
        <div className="flex items-center gap-2">
          <TrendingUp className="w-5 h-5 text-primary" />
          <h3 className="font-medium">Your Feed</h3>
          <Badge variant="secondary" className="text-xs">
            {feedData.posts.length} posts
          </Badge>
        </div>

        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm">
                <Filter className="w-4 h-4 mr-1" />
                {algorithm.charAt(0).toUpperCase() + algorithm.slice(1)}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => handleAlgorithmChange('personalized')}>
                <Star className="w-4 h-4 mr-2" />
                Personalized
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleAlgorithmChange('chronological')}>
                <Calendar className="w-4 h-4 mr-2" />
                Latest
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleAlgorithmChange('trending')}>
                <TrendingUp className="w-4 h-4 mr-2" />
                Trending
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {/* Posts */}
      <div className="space-y-4">
        {feedData.posts.map((post) => (
          <PostCard
            key={post.id}
            post={post}
            currentUser={user}
            onLike={handleLike}
            onComment={handleComment}
            onShare={handleShare}
            onBookmark={handleBookmark}
            onPostClick={onPostClick}
            onUserClick={onUserClick}
          />
        ))}
      </div>

      {/* Load More */}
      {feedData.hasMore && (
        <div className="text-center py-4">
          <Button
            variant="outline"
            onClick={() => {
              // Implement pagination with cursor
              setFeedParams(prev => ({ ...prev, cursor: feedData.cursor }));
            }}
          >
            Load More Posts
          </Button>
        </div>
      )}
    </div>
  );
};

export default SocialFeed;