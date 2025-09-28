import React, { useState } from 'react';
import { toast } from "sonner";
import { Plus, Camera, Image, MapPin } from 'lucide-react';
import { Button } from './ui/button';
import { Card, CardContent } from './ui/card';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';

interface CreatePostSectionProps {
  user: any;
  onCreatePhotoPost?: () => void;
  onCreateGalleryPost?: () => void;
  onCreateLocationPost?: () => void;
  onCreateTextPost?: (text: string) => void;
}

export function CreatePostSection({
  user,
  onCreatePhotoPost,
  onCreateGalleryPost,
  onCreateLocationPost,
  onCreateTextPost
}: CreatePostSectionProps) {
  const [textContent, setTextContent] = useState('');

  const handleCreatePost = () => {
    const text = textContent.trim();
    if (!text) return;
    
    // Call parent callback if available, otherwise handle internally
    if (onCreateTextPost) {
      onCreateTextPost(text);
    } else {
      // Internal post creation with toast feedback
      toast.success(`Post created: "${text.slice(0, 50)}${text.length > 50 ? '...' : ''}"`);
    }
    
    // Clear the text input
    setTextContent('');
    
    // Reset textarea height
    const textarea = document.querySelector('textarea');
    if (textarea) {
      textarea.style.height = 'auto';
    }
  };

  const handleCreatePhotoPost = () => {
    if (onCreatePhotoPost) {
      onCreatePhotoPost();
    } else {
      toast.success('📷 Opening camera for photo post...');
    }
  };

  const handleCreateGalleryPost = () => {
    if (onCreateGalleryPost) {
      onCreateGalleryPost();
    } else {
      toast.success('🖼️ Opening gallery for photo selection...');
    }
  };

  const handleCreateLocationPost = () => {
    if (onCreateLocationPost) {
      onCreateLocationPost();
    } else {
      toast.success('📍 Adding current location to post...');
    }
  };
  return (
    <div className="px-4 py-3 bg-card border-b border-border" data-testid="create-post-section">
      <Card className="hover:bg-accent/30 transition-colors" data-testid="create-post-card">
        <CardContent className="p-4 space-y-4" data-testid="create-post-content">
          <div className="flex items-center gap-4" data-testid="create-post-input-area">
            <Avatar className="w-12 h-12 ring-2 ring-primary/10" data-testid="create-post-user-avatar">
              <AvatarImage src={user?.avatar || ''} />
              <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
                {(user?.name || 'You').charAt(0)}
              </AvatarFallback>
            </Avatar>

            <div className="flex-1 relative" data-testid="create-post-textarea-container">
              <textarea
                placeholder="What's happening in Bangkok today? 🍜"
                className="w-full bg-muted/60 hover:bg-muted/80 focus:bg-muted/80 rounded-2xl px-5 py-3 text-sm border border-border/50 focus:border-primary/30 focus:outline-none focus:ring-2 focus:ring-primary/10 transition-all resize-none min-h-[44px] max-h-32"
                rows={1}
                value={textContent}
                onChange={(e) => setTextContent(e.target.value)}
                onInput={(e) => {
                  const target = e.target as HTMLTextAreaElement;
                  target.style.height = 'auto';
                  target.style.height = Math.min(target.scrollHeight, 128) + 'px';
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    handleCreatePost();
                  }
                }}
                data-testid="create-post-textarea-input"
              />
            </div>
          </div>
          
          <div className="flex items-center justify-between" data-testid="create-post-actions">
            <div className="flex items-center gap-1" data-testid="create-post-media-buttons">
              <Button
                variant="ghost"
                size="sm"
                className="text-chart-1 hover:bg-chart-1/10 rounded-xl px-3 h-8 gap-2 touch-manipulation"
                onClick={handleCreatePhotoPost}
                data-testid="create-post-photo-button"
              >
                <Camera className="w-4 h-4" />
                <span className="text-xs hidden sm:inline">Photo</span>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="text-chart-2 hover:bg-chart-2/10 rounded-xl px-3 h-8 gap-2 touch-manipulation"
                onClick={handleCreateGalleryPost}
                data-testid="create-post-gallery-button"
              >
                <Image className="w-4 h-4" />
                <span className="text-xs hidden sm:inline">Gallery</span>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="text-chart-3 hover:bg-chart-3/10 rounded-xl px-3 h-8 gap-2 touch-manipulation"
                onClick={handleCreateLocationPost}
                data-testid="create-post-location-button"
              >
                <MapPin className="w-4 h-4" />
                <span className="text-xs hidden sm:inline">Location</span>
              </Button>
            </div>

            <Button
              size="sm"
              className="gap-2 rounded-xl bg-gradient-to-r from-primary to-primary/90 hover:from-primary/90 hover:to-primary/80 shadow-sm touch-manipulation"
              onClick={handleCreatePost}
              disabled={!textContent.trim()}
              data-testid="create-post-submit-button"
              data-disabled={!textContent.trim() ? 'true' : 'false'}
            >
              <Plus className="w-4 h-4" />
              <span className="hidden sm:inline">Create Post</span>
              <span className="sm:hidden">Post</span>
            </Button>
          </div>
          
          <div className="text-xs text-muted-foreground/80 bg-muted/30 rounded-lg px-3 py-2 border-l-2 border-chart-1/30" data-testid="create-post-hint">
            💡 Share your Bangkok adventures, food discoveries, or local insights with the community
          </div>
        </CardContent>
      </Card>
    </div>
  );
}