/**
 * SocialExample Component
 *
 * Comprehensive example demonstrating the integration of all social components
 * with the real social API service. Shows how to use SocialFeed, SocialProfile,
 * and PostCreation together in a complete social platform interface.
 */

import React, { useState } from 'react';
import { User, Plus, TrendingUp, Users, Star, MessageCircle, Settings } from 'lucide-react';
import { Button } from './ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import SocialFeed from './SocialFeed';
import SocialProfile from './SocialProfile';
import PostCreation from './PostCreation';
import { toast } from "sonner";

interface SocialExampleProps {
  user: any;
  className?: string;
}

export const SocialExample: React.FC<SocialExampleProps> = ({
  user,
  className = '',
}) => {
  const [activeTab, setActiveTab] = useState('feed');
  const [showCreatePost, setShowCreatePost] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
  const [showProfile, setShowProfile] = useState(false);

  const handlePostCreated = (post: any) => {
    setShowCreatePost(false);
    toast.success('Post created successfully!');
    // The RTK Query cache will automatically update the feed
  };

  const handleUserClick = (userId: string) => {
    setSelectedUserId(userId);
    setShowProfile(true);
  };

  const handlePostClick = (postId: string) => {
    // Navigate to post detail or open post modal
    toast.info(`Opening post ${postId}`);
  };

  const handleMessage = (userId: string) => {
    // Navigate to messages or open chat
    toast.info(`Opening chat with user ${userId}`);
  };

  return (
    <div className={`h-full flex flex-col ${className}`}>
      {/* Header */}
      <div className="sticky top-0 z-10 bg-background border-b">
        <div className="flex items-center justify-between p-4">
          <div className="flex items-center gap-2">
            <Users className="w-6 h-6 text-primary" />
            <h1 className="text-xl font-semibold">Social</h1>
            <Badge variant="secondary" className="text-xs">Live</Badge>
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowCreatePost(true)}
            >
              <Plus className="w-4 h-4 mr-2" />
              Create Post
            </Button>

            <Button variant="ghost" size="icon">
              <Settings className="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-hidden">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="h-full flex flex-col">
          <TabsList className="grid w-full grid-cols-4 px-4 py-2">
            <TabsTrigger value="feed" className="flex items-center gap-2">
              <TrendingUp className="w-4 h-4" />
              <span className="hidden sm:inline">Feed</span>
            </TabsTrigger>
            <TabsTrigger value="profile" className="flex items-center gap-2">
              <User className="w-4 h-4" />
              <span className="hidden sm:inline">Profile</span>
            </TabsTrigger>
            <TabsTrigger value="trending" className="flex items-center gap-2">
              <Star className="w-4 h-4" />
              <span className="hidden sm:inline">Trending</span>
            </TabsTrigger>
            <TabsTrigger value="messages" className="flex items-center gap-2">
              <MessageCircle className="w-4 h-4" />
              <span className="hidden sm:inline">Messages</span>
            </TabsTrigger>
          </TabsList>

          <div className="flex-1 overflow-hidden">
            {/* Social Feed Tab */}
            <TabsContent value="feed" className="h-full mt-0">
              <div className="h-full overflow-auto p-4">
                <SocialFeed
                  user={user}
                  algorithm="personalized"
                  region="TH"
                  onPostClick={handlePostClick}
                  onUserClick={handleUserClick}
                />
              </div>
            </TabsContent>

            {/* User Profile Tab */}
            <TabsContent value="profile" className="h-full mt-0">
              <div className="h-full overflow-auto p-4">
                <SocialProfile
                  userId={user?.id || 'current-user'}
                  currentUser={user}
                  onMessage={handleMessage}
                  onEditProfile={() => toast.info('Opening profile settings')}
                />
              </div>
            </TabsContent>

            {/* Trending Tab */}
            <TabsContent value="trending" className="h-full mt-0">
              <div className="h-full overflow-auto p-4">
                <SocialFeed
                  user={user}
                  algorithm="trending"
                  region="TH"
                  onPostClick={handlePostClick}
                  onUserClick={handleUserClick}
                />
              </div>
            </TabsContent>

            {/* Messages Tab */}
            <TabsContent value="messages" className="h-full mt-0">
              <div className="h-full flex items-center justify-center">
                <Card className="w-full max-w-md">
                  <CardHeader>
                    <CardTitle className="text-center">Messages</CardTitle>
                  </CardHeader>
                  <CardContent className="text-center">
                    <MessageCircle className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                    <p className="text-muted-foreground mb-4">
                      Direct messaging functionality would be integrated here.
                    </p>
                    <Button onClick={() => toast.info('Messages feature coming soon!')}>
                      Explore Messages
                    </Button>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>
          </div>
        </Tabs>
      </div>

      {/* Create Post Dialog */}
      <Dialog open={showCreatePost} onOpenChange={setShowCreatePost}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-hidden">
          <DialogHeader>
            <DialogTitle>Create New Post</DialogTitle>
          </DialogHeader>
          <div className="overflow-auto max-h-[calc(90vh-100px)]">
            <PostCreation
              user={user}
              onPostCreated={handlePostCreated}
              onCancel={() => setShowCreatePost(false)}
            />
          </div>
        </DialogContent>
      </Dialog>

      {/* User Profile Dialog */}
      <Dialog open={showProfile} onOpenChange={setShowProfile}>
        <DialogContent className="max-w-4xl max-h-[90vh] overflow-hidden">
          <DialogHeader>
            <DialogTitle>User Profile</DialogTitle>
          </DialogHeader>
          <div className="overflow-auto max-h-[calc(90vh-100px)]">
            {selectedUserId && (
              <SocialProfile
                userId={selectedUserId}
                currentUser={user}
                onMessage={handleMessage}
                onEditProfile={() => {
                  setShowProfile(false);
                  toast.info('Opening profile settings');
                }}
              />
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default SocialExample;