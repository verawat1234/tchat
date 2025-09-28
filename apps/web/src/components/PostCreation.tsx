/**
 * PostCreation Component
 *
 * Comprehensive post creation interface with media upload, tagging, and privacy controls.
 * Supports text, image, video, and link posts with rich editing capabilities.
 */

import React, { useState, useRef, useCallback } from 'react';
import {
  Image,
  Video,
  Link,
  MapPin,
  Users,
  Globe,
  Lock,
  UserCheck,
  Hash,
  Smile,
  X,
  Upload,
  FileText,
  Calendar,
  TrendingUp,
  AlertCircle,
  CheckCircle,
  Loader2
} from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Textarea } from './ui/textarea';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from './ui/dialog';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger, DropdownMenuSeparator } from './ui/dropdown-menu';
import { Progress } from './ui/progress';
import { toast } from "sonner";
import { useCreatePostMutation, getSocialDisplayName } from '../services/socialApi';
import type { CreatePostRequest } from '../types/social';

interface PostCreationProps {
  user: any;
  onPostCreated?: (post: any) => void;
  onCancel?: () => void;
  className?: string;
  defaultVisibility?: 'public' | 'members' | 'private' | 'followers';
  communityId?: string;
}

interface MediaFile {
  id: string;
  file: File;
  url: string;
  type: 'image' | 'video';
  size: number;
  progress?: number;
}

const PostCreation: React.FC<PostCreationProps> = ({
  user,
  onPostCreated,
  onCancel,
  className = '',
  defaultVisibility = 'public',
  communityId,
}) => {
  const [content, setContent] = useState('');
  const [postType, setPostType] = useState<'text' | 'image' | 'video' | 'link' | 'poll'>('text');
  const [visibility, setVisibility] = useState<'public' | 'members' | 'private' | 'followers'>(defaultVisibility);
  const [tags, setTags] = useState<string[]>([]);
  const [currentTag, setCurrentTag] = useState('');
  const [location, setLocation] = useState('');
  const [mediaFiles, setMediaFiles] = useState<MediaFile[]>([]);
  const [linkUrl, setLinkUrl] = useState('');
  const [linkPreview, setLinkPreview] = useState<any>(null);
  const [isLoadingPreview, setIsLoadingPreview] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const fileInputRef = useRef<HTMLInputElement>(null);
  const [createPost, { isLoading: isCreating }] = useCreatePostMutation();

  // Character limits
  const MAX_CONTENT_LENGTH = 2000;
  const MAX_TAGS = 10;
  const MAX_MEDIA_FILES = 10;
  const MAX_FILE_SIZE = 50 * 1024 * 1024; // 50MB

  // Handle content change
  const handleContentChange = (value: string) => {
    if (value.length <= MAX_CONTENT_LENGTH) {
      setContent(value);
    }
  };

  // Handle tag addition
  const addTag = useCallback(() => {
    const tag = currentTag.trim().toLowerCase();
    if (tag && !tags.includes(tag) && tags.length < MAX_TAGS) {
      setTags(prev => [...prev, tag]);
      setCurrentTag('');
    }
  }, [currentTag, tags]);

  // Handle tag removal
  const removeTag = useCallback((tagToRemove: string) => {
    setTags(prev => prev.filter(tag => tag !== tagToRemove));
  }, []);

  // Handle media file selection
  const handleFileSelect = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || []);

    files.forEach(file => {
      if (file.size > MAX_FILE_SIZE) {
        toast.error(`File ${file.name} is too large (max 50MB)`);
        return;
      }

      if (mediaFiles.length >= MAX_MEDIA_FILES) {
        toast.error(`Maximum ${MAX_MEDIA_FILES} files allowed`);
        return;
      }

      const mediaFile: MediaFile = {
        id: `${Date.now()}-${Math.random()}`,
        file,
        url: URL.createObjectURL(file),
        type: file.type.startsWith('video/') ? 'video' : 'image',
        size: file.size,
      };

      setMediaFiles(prev => [...prev, mediaFile]);

      // Set post type based on first file
      if (mediaFiles.length === 0) {
        setPostType(mediaFile.type);
      }
    });

    // Reset file input
    if (event.target) {
      event.target.value = '';
    }
  }, [mediaFiles]);

  // Remove media file
  const removeMediaFile = useCallback((fileId: string) => {
    setMediaFiles(prev => {
      const updated = prev.filter(f => f.id !== fileId);
      // Revoke object URL to prevent memory leaks
      const fileToRemove = prev.find(f => f.id === fileId);
      if (fileToRemove) {
        URL.revokeObjectURL(fileToRemove.url);
      }
      // Reset post type if no media files left
      if (updated.length === 0) {
        setPostType('text');
      }
      return updated;
    });
  }, []);

  // Handle link preview
  const handleLinkPreview = useCallback(async () => {
    if (!linkUrl.trim()) return;

    setIsLoadingPreview(true);
    try {
      // In a real app, you'd call a link preview service
      // For now, we'll create a mock preview
      const mockPreview = {
        title: 'Link Preview',
        description: 'This is a preview of the shared link',
        image: null,
        url: linkUrl,
      };
      setLinkPreview(mockPreview);
      setPostType('link');
    } catch (error) {
      toast.error('Failed to load link preview');
    } finally {
      setIsLoadingPreview(false);
    }
  }, [linkUrl]);

  // Handle post submission
  const handleSubmit = async () => {
    if (!content.trim() && mediaFiles.length === 0 && !linkUrl.trim()) {
      toast.error('Please add some content to your post');
      return;
    }

    setIsUploading(true);

    try {
      // In a real app, you'd upload media files first and get URLs
      const mediaUrls: string[] = [];

      if (mediaFiles.length > 0) {
        // Mock upload process
        for (const mediaFile of mediaFiles) {
          // Simulate upload progress
          await new Promise(resolve => setTimeout(resolve, 500));
          mediaUrls.push(mediaFile.url);
        }
      }

      const postData: CreatePostRequest = {
        content: content.trim(),
        type: postType,
        visibility,
        tags: tags.length > 0 ? tags : undefined,
        mediaUrls: mediaUrls.length > 0 ? mediaUrls : undefined,
        communityId: communityId || undefined,
        metadata: {
          ...(location && { location }),
          ...(linkPreview && { linkPreview }),
        },
      };

      const result = await createPost(postData).unwrap();

      toast.success('Post created successfully!');

      // Clean up
      mediaFiles.forEach(file => URL.revokeObjectURL(file.url));
      setContent('');
      setTags([]);
      setMediaFiles([]);
      setLinkUrl('');
      setLinkPreview(null);
      setLocation('');
      setPostType('text');
      setCurrentTag('');

      if (onPostCreated) {
        onPostCreated(result);
      }

    } catch (error) {
      toast.error('Failed to create post');
      console.error('Error creating post:', error);
    } finally {
      setIsUploading(false);
    }
  };

  // Get visibility icon and label
  const getVisibilityInfo = () => {
    switch (visibility) {
      case 'public':
        return { icon: Globe, label: 'Public', description: 'Anyone can see this post' };
      case 'followers':
        return { icon: UserCheck, label: 'Followers', description: 'Only your followers can see this' };
      case 'members':
        return { icon: Users, label: 'Members', description: 'Only community members can see this' };
      case 'private':
        return { icon: Lock, label: 'Private', description: 'Only you can see this post' };
      default:
        return { icon: Globe, label: 'Public', description: 'Anyone can see this post' };
    }
  };

  const visibilityInfo = getVisibilityInfo();
  const remainingChars = MAX_CONTENT_LENGTH - content.length;
  const isOverLimit = remainingChars < 0;

  return (
    <Card className={`w-full ${className}`}>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FileText className="w-5 h-5" />
          Create Post
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* User Info */}
        <div className="flex items-center gap-3">
          <Avatar>
            <AvatarImage src={user?.avatar} />
            <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
              {getSocialDisplayName(user || {}).charAt(0)}
            </AvatarFallback>
          </Avatar>
          <div>
            <div className="font-medium">{getSocialDisplayName(user || {})}</div>
            <div className="text-sm text-muted-foreground flex items-center gap-1">
              <visibilityInfo.icon className="w-3 h-3" />
              {visibilityInfo.label}
            </div>
          </div>
        </div>

        {/* Content Input */}
        <div className="space-y-2">
          <Textarea
            placeholder="What's on your mind?"
            value={content}
            onChange={(e) => handleContentChange(e.target.value)}
            className="min-h-[120px] border-0 focus-visible:ring-0 focus-visible:ring-offset-0 resize-none text-base"
            disabled={isCreating || isUploading}
          />

          <div className="flex items-center justify-between text-sm">
            <div className="flex items-center gap-2">
              {isOverLimit && (
                <AlertCircle className="w-4 h-4 text-destructive" />
              )}
              <span className={isOverLimit ? 'text-destructive' : 'text-muted-foreground'}>
                {remainingChars} characters remaining
              </span>
            </div>

            {content.length > 0 && (
              <Progress
                value={(content.length / MAX_CONTENT_LENGTH) * 100}
                className="w-16 h-2"
              />
            )}
          </div>
        </div>

        {/* Link Input (for link posts) */}
        {postType === 'link' && (
          <div className="space-y-2">
            <div className="flex gap-2">
              <Input
                placeholder="Paste a link..."
                value={linkUrl}
                onChange={(e) => setLinkUrl(e.target.value)}
                className="flex-1"
              />
              <Button
                variant="outline"
                onClick={handleLinkPreview}
                disabled={!linkUrl.trim() || isLoadingPreview}
              >
                {isLoadingPreview ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  'Preview'
                )}
              </Button>
            </div>

            {linkPreview && (
              <Card className="border-l-4 border-l-primary">
                <CardContent className="p-3">
                  <div className="font-medium">{linkPreview.title}</div>
                  <div className="text-sm text-muted-foreground">{linkPreview.description}</div>
                  <div className="text-xs text-muted-foreground mt-1">{linkPreview.url}</div>
                </CardContent>
              </Card>
            )}
          </div>
        )}

        {/* Media Files */}
        {mediaFiles.length > 0 && (
          <div className="space-y-2">
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
              {mediaFiles.map((file) => (
                <div key={file.id} className="relative group">
                  <div className="aspect-square rounded-lg overflow-hidden bg-muted">
                    {file.type === 'image' ? (
                      <img
                        src={file.url}
                        alt="Upload preview"
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <video
                        src={file.url}
                        className="w-full h-full object-cover"
                        controls={false}
                      />
                    )}
                  </div>
                  <Button
                    variant="destructive"
                    size="icon"
                    className="absolute top-1 right-1 w-6 h-6 opacity-0 group-hover:opacity-100 transition-opacity"
                    onClick={() => removeMediaFile(file.id)}
                  >
                    <X className="w-3 h-3" />
                  </Button>
                  {file.type === 'video' && (
                    <Video className="absolute bottom-1 left-1 w-4 h-4 text-white" />
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Tags */}
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {tags.map((tag) => (
              <Badge
                key={tag}
                variant="secondary"
                className="cursor-pointer hover:bg-destructive hover:text-destructive-foreground"
                onClick={() => removeTag(tag)}
              >
                #{tag} ×
              </Badge>
            ))}
          </div>
        )}

        {/* Advanced Options */}
        {showAdvanced && (
          <div className="space-y-4 p-4 bg-muted/50 rounded-lg">
            {/* Tag Input */}
            <div className="space-y-2">
              <Label>Tags</Label>
              <div className="flex gap-2">
                <Input
                  placeholder="Add a tag..."
                  value={currentTag}
                  onChange={(e) => setCurrentTag(e.target.value)}
                  onKeyPress={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      addTag();
                    }
                  }}
                  className="flex-1"
                  disabled={tags.length >= MAX_TAGS}
                />
                <Button
                  variant="outline"
                  onClick={addTag}
                  disabled={!currentTag.trim() || tags.length >= MAX_TAGS}
                >
                  Add
                </Button>
              </div>
              <div className="text-xs text-muted-foreground">
                {tags.length}/{MAX_TAGS} tags
              </div>
            </div>

            {/* Location */}
            <div className="space-y-2">
              <Label>Location</Label>
              <Input
                placeholder="Add a location..."
                value={location}
                onChange={(e) => setLocation(e.target.value)}
              />
            </div>

            {/* Visibility */}
            <div className="space-y-2">
              <Label>Who can see this post?</Label>
              <Select value={visibility} onValueChange={(value: any) => setVisibility(value)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="public">
                    <div className="flex items-center gap-2">
                      <Globe className="w-4 h-4" />
                      Public
                    </div>
                  </SelectItem>
                  <SelectItem value="followers">
                    <div className="flex items-center gap-2">
                      <UserCheck className="w-4 h-4" />
                      Followers
                    </div>
                  </SelectItem>
                  {communityId && (
                    <SelectItem value="members">
                      <div className="flex items-center gap-2">
                        <Users className="w-4 h-4" />
                        Members
                      </div>
                    </SelectItem>
                  )}
                  <SelectItem value="private">
                    <div className="flex items-center gap-2">
                      <Lock className="w-4 h-4" />
                      Private
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
              <div className="text-xs text-muted-foreground">
                {visibilityInfo.description}
              </div>
            </div>
          </div>
        )}

        {/* Toolbar */}
        <div className="flex items-center justify-between pt-4 border-t">
          <div className="flex items-center gap-2">
            {/* Media Upload */}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm">
                  <Image className="w-4 h-4 mr-2" />
                  Media
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start">
                <DropdownMenuItem onClick={() => fileInputRef.current?.click()}>
                  <Image className="w-4 h-4 mr-2" />
                  Photo/Video
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => setPostType('link')}>
                  <Link className="w-4 h-4 mr-2" />
                  Link
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            {/* Add Tags */}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setShowAdvanced(!showAdvanced)}
            >
              <Hash className="w-4 h-4 mr-2" />
              {showAdvanced ? 'Less' : 'More'}
            </Button>

            {/* Location */}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setShowAdvanced(true)}
            >
              <MapPin className="w-4 h-4 mr-2" />
              Location
            </Button>
          </div>

          {/* Submit Actions */}
          <div className="flex items-center gap-2">
            {onCancel && (
              <Button
                variant="outline"
                onClick={onCancel}
                disabled={isCreating || isUploading}
              >
                Cancel
              </Button>
            )}

            <Button
              onClick={handleSubmit}
              disabled={
                isCreating ||
                isUploading ||
                isOverLimit ||
                (!content.trim() && mediaFiles.length === 0 && !linkUrl.trim())
              }
            >
              {isCreating || isUploading ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  {isUploading ? 'Uploading...' : 'Creating...'}
                </>
              ) : (
                <>
                  <CheckCircle className="w-4 h-4 mr-2" />
                  Post
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Hidden File Input */}
        <input
          ref={fileInputRef}
          type="file"
          multiple
          accept="image/*,video/*"
          onChange={handleFileSelect}
          className="hidden"
        />
      </CardContent>
    </Card>
  );
};

export default PostCreation;