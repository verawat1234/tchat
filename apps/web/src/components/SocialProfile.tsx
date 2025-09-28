/**
 * SocialProfile Component
 *
 * Comprehensive user profile component for social features.
 * Displays user information, statistics, and provides interaction options.
 */

import React, { useState, useMemo } from 'react';
import { User, UserPlus, UserCheck, Settings, MoreVertical, MapPin, Calendar, Star, Heart, MessageCircle, Share, Users, TrendingUp, Edit, Link, Mail, Phone, Globe, Shield, Eye, EyeOff } from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from './ui/dialog';
import { Input } from './ui/input';
import { Textarea } from './ui/textarea';
import { Label } from './ui/label';
import { Switch } from './ui/switch';
import { Separator } from './ui/separator';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { toast } from "sonner";
import {
  useGetSocialProfileQuery,
  useUpdateSocialProfileMutation,
  useFollowUserMutation,
  useUnfollowUserMutation,
  useGetFollowersQuery,
  useGetFollowingQuery,
  useGetUserAnalyticsQuery,
  getSocialDisplayName,
} from '../services/socialApi';
import type { SocialProfile, UpdateSocialProfileRequest } from '../types/social';

interface SocialProfileProps {
  userId: string;
  currentUser?: any;
  onMessage?: (userId: string) => void;
  onEditProfile?: () => void;
  className?: string;
}

interface ProfileStatsProps {
  profile: SocialProfile;
  onFollowersClick: () => void;
  onFollowingClick: () => void;
}

const ProfileStats: React.FC<ProfileStatsProps> = ({
  profile,
  onFollowersClick,
  onFollowingClick,
}) => {
  return (
    <div className="grid grid-cols-3 gap-4 text-center">
      <div>
        <div className="text-2xl font-bold text-foreground">{profile.postsCount}</div>
        <div className="text-sm text-muted-foreground">Posts</div>
      </div>
      <button
        onClick={onFollowersClick}
        className="hover:bg-muted/50 rounded-lg p-2 transition-colors"
      >
        <div className="text-2xl font-bold text-foreground">{profile.followersCount}</div>
        <div className="text-sm text-muted-foreground">Followers</div>
      </button>
      <button
        onClick={onFollowingClick}
        className="hover:bg-muted/50 rounded-lg p-2 transition-colors"
      >
        <div className="text-2xl font-bold text-foreground">{profile.followingCount}</div>
        <div className="text-sm text-muted-foreground">Following</div>
      </button>
    </div>
  );
};

interface EditProfileDialogProps {
  profile: SocialProfile;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSave: (updates: UpdateSocialProfileRequest) => void;
}

const EditProfileDialog: React.FC<EditProfileDialogProps> = ({
  profile,
  open,
  onOpenChange,
  onSave,
}) => {
  const [formData, setFormData] = useState<UpdateSocialProfileRequest>({
    displayName: profile.display_name || '',
    bio: profile.bio || '',
    avatar: profile.avatar || '',
    interests: profile.interests || [],
    socialLinks: profile.socialLinks || {},
    socialPreferences: profile.socialPreferences || {},
  });

  const [newInterest, setNewInterest] = useState('');
  const [newSocialLink, setNewSocialLink] = useState({ platform: '', url: '' });

  const handleSave = () => {
    onSave(formData);
    onOpenChange(false);
  };

  const addInterest = () => {
    if (newInterest.trim() && !formData.interests?.includes(newInterest.trim())) {
      setFormData(prev => ({
        ...prev,
        interests: [...(prev.interests || []), newInterest.trim()]
      }));
      setNewInterest('');
    }
  };

  const removeInterest = (interest: string) => {
    setFormData(prev => ({
      ...prev,
      interests: prev.interests?.filter(i => i !== interest) || []
    }));
  };

  const addSocialLink = () => {
    if (newSocialLink.platform.trim() && newSocialLink.url.trim()) {
      setFormData(prev => ({
        ...prev,
        socialLinks: {
          ...prev.socialLinks,
          [newSocialLink.platform]: newSocialLink.url
        }
      }));
      setNewSocialLink({ platform: '', url: '' });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Edit Profile</DialogTitle>
          <DialogDescription>
            Update your profile information and social preferences
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Basic Information */}
          <div className="space-y-4">
            <div>
              <Label htmlFor="displayName">Display Name</Label>
              <Input
                id="displayName"
                value={formData.displayName || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, displayName: e.target.value }))}
                placeholder="Your display name"
              />
            </div>

            <div>
              <Label htmlFor="bio">Bio</Label>
              <Textarea
                id="bio"
                value={formData.bio || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, bio: e.target.value }))}
                placeholder="Tell us about yourself..."
                rows={3}
              />
            </div>

            <div>
              <Label htmlFor="avatar">Avatar URL</Label>
              <Input
                id="avatar"
                value={formData.avatar || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, avatar: e.target.value }))}
                placeholder="https://example.com/avatar.jpg"
              />
            </div>
          </div>

          <Separator />

          {/* Interests */}
          <div className="space-y-4">
            <Label>Interests</Label>
            <div className="flex gap-2">
              <Input
                value={newInterest}
                onChange={(e) => setNewInterest(e.target.value)}
                placeholder="Add an interest"
                onKeyPress={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault();
                    addInterest();
                  }
                }}
              />
              <Button onClick={addInterest} size="sm">Add</Button>
            </div>
            <div className="flex flex-wrap gap-2">
              {formData.interests?.map((interest, index) => (
                <Badge
                  key={index}
                  variant="secondary"
                  className="cursor-pointer hover:bg-destructive hover:text-destructive-foreground"
                  onClick={() => removeInterest(interest)}
                >
                  {interest} ×
                </Badge>
              ))}
            </div>
          </div>

          <Separator />

          {/* Social Links */}
          <div className="space-y-4">
            <Label>Social Links</Label>
            <div className="flex gap-2">
              <Input
                value={newSocialLink.platform}
                onChange={(e) => setNewSocialLink(prev => ({ ...prev, platform: e.target.value }))}
                placeholder="Platform (e.g., twitter)"
                className="flex-1"
              />
              <Input
                value={newSocialLink.url}
                onChange={(e) => setNewSocialLink(prev => ({ ...prev, url: e.target.value }))}
                placeholder="URL"
                className="flex-1"
              />
              <Button onClick={addSocialLink} size="sm">Add</Button>
            </div>
            <div className="space-y-2">
              {Object.entries(formData.socialLinks || {}).map(([platform, url]) => (
                <div key={platform} className="flex items-center justify-between p-2 bg-muted rounded">
                  <div className="flex items-center gap-2">
                    <Link className="w-4 h-4" />
                    <span className="text-sm font-medium">{platform}</span>
                    <span className="text-sm text-muted-foreground truncate">{url as string}</span>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => {
                      const newLinks = { ...formData.socialLinks };
                      delete newLinks[platform];
                      setFormData(prev => ({ ...prev, socialLinks: newLinks }));
                    }}
                  >
                    ×
                  </Button>
                </div>
              ))}
            </div>
          </div>

          <div className="flex gap-2 pt-4">
            <Button onClick={handleSave} className="flex-1">
              Save Changes
            </Button>
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export const SocialProfile: React.FC<SocialProfileProps> = ({
  userId,
  currentUser,
  onMessage,
  onEditProfile,
  className = '',
}) => {
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [followersDialogOpen, setFollowersDialogOpen] = useState(false);
  const [followingDialogOpen, setFollowingDialogOpen] = useState(false);
  const [isFollowing, setIsFollowing] = useState(false);

  // RTK Query hooks
  const { data: profile, isLoading: profileLoading, error: profileError } = useGetSocialProfileQuery(userId);
  const { data: followers } = useGetFollowersQuery({ userId, limit: 50 });
  const { data: following } = useGetFollowingQuery({ userId, limit: 50 });
  const { data: analytics } = useGetUserAnalyticsQuery({ userId, period: '30d' });

  // Mutations
  const [updateProfile] = useUpdateSocialProfileMutation();
  const [followUser] = useFollowUserMutation();
  const [unfollowUser] = useUnfollowUserMutation();

  const isOwnProfile = currentUser?.id === userId;

  // Handlers
  const handleFollowToggle = async () => {
    if (!currentUser?.id) {
      toast.error('Please log in to follow users');
      return;
    }

    try {
      if (isFollowing) {
        await unfollowUser({
          followerId: currentUser.id,
          followingId: userId,
        }).unwrap();
        setIsFollowing(false);
        toast.success('User unfollowed');
      } else {
        await followUser({
          followerId: currentUser.id,
          followingId: userId,
          source: 'manual',
        }).unwrap();
        setIsFollowing(true);
        toast.success('User followed!');
      }
    } catch (error) {
      toast.error('Failed to update follow status');
    }
  };

  const handleUpdateProfile = async (updates: UpdateSocialProfileRequest) => {
    try {
      await updateProfile({ userId, updates }).unwrap();
      toast.success('Profile updated successfully!');
    } catch (error) {
      toast.error('Failed to update profile');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
    });
  };

  // Loading state
  if (profileLoading) {
    return (
      <div className={`space-y-6 ${className}`}>
        <Card>
          <CardContent className="p-6 space-y-4">
            <div className="flex items-start gap-4">
              <div className="w-20 h-20 bg-muted rounded-full animate-pulse" />
              <div className="flex-1 space-y-2">
                <div className="h-6 bg-muted rounded animate-pulse w-1/3" />
                <div className="h-4 bg-muted rounded animate-pulse w-1/2" />
                <div className="h-4 bg-muted rounded animate-pulse w-2/3" />
              </div>
            </div>
            <div className="grid grid-cols-3 gap-4">
              {[...Array(3)].map((_, i) => (
                <div key={i} className="text-center space-y-2">
                  <div className="h-8 bg-muted rounded animate-pulse" />
                  <div className="h-4 bg-muted rounded animate-pulse" />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Error state
  if (profileError || !profile) {
    return (
      <div className={`text-center py-8 ${className}`}>
        <User className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
        <h3 className="font-medium mb-2">Profile not found</h3>
        <p className="text-sm text-muted-foreground">
          This user profile could not be loaded.
        </p>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Profile Header */}
      <Card>
        <CardContent className="p-6">
          <div className="flex flex-col sm:flex-row items-start gap-4">
            {/* Avatar */}
            <Avatar className="w-20 h-20 sm:w-24 sm:h-24">
              <AvatarImage src={profile.avatar} />
              <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white text-xl">
                {getSocialDisplayName(profile).charAt(0)}
              </AvatarFallback>
            </Avatar>

            {/* Profile Info */}
            <div className="flex-1 w-full space-y-4">
              <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <h1 className="text-2xl font-bold">{getSocialDisplayName(profile)}</h1>
                    {profile.isSocialVerified && (
                      <Star className="w-5 h-5 text-yellow-500 fill-current" />
                    )}
                  </div>

                  {profile.username && (
                    <p className="text-muted-foreground">@{profile.username}</p>
                  )}

                  {profile.bio && (
                    <p className="text-foreground leading-relaxed">{profile.bio}</p>
                  )}

                  <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
                    {profile.country && (
                      <div className="flex items-center gap-1">
                        <MapPin className="w-4 h-4" />
                        <span>{profile.country}</span>
                      </div>
                    )}

                    <div className="flex items-center gap-1">
                      <Calendar className="w-4 h-4" />
                      <span>Joined {formatDate(profile.socialCreatedAt)}</span>
                    </div>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="flex gap-2">
                  {isOwnProfile ? (
                    <>
                      <Button
                        variant="outline"
                        onClick={() => setEditDialogOpen(true)}
                      >
                        <Edit className="w-4 h-4 mr-2" />
                        Edit Profile
                      </Button>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon">
                            <MoreVertical className="w-4 h-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={onEditProfile}>
                            <Settings className="w-4 h-4 mr-2" />
                            Settings
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </>
                  ) : (
                    <>
                      <Button
                        variant={isFollowing ? "outline" : "default"}
                        onClick={handleFollowToggle}
                      >
                        {isFollowing ? (
                          <>
                            <UserCheck className="w-4 h-4 mr-2" />
                            Following
                          </>
                        ) : (
                          <>
                            <UserPlus className="w-4 h-4 mr-2" />
                            Follow
                          </>
                        )}
                      </Button>

                      {onMessage && (
                        <Button
                          variant="outline"
                          onClick={() => onMessage(userId)}
                        >
                          <MessageCircle className="w-4 h-4 mr-2" />
                          Message
                        </Button>
                      )}
                    </>
                  )}
                </div>
              </div>

              {/* Profile Stats */}
              <ProfileStats
                profile={profile}
                onFollowersClick={() => setFollowersDialogOpen(true)}
                onFollowingClick={() => setFollowingDialogOpen(true)}
              />
            </div>
          </div>

          {/* Interests */}
          {profile.interests && profile.interests.length > 0 && (
            <div className="mt-6 pt-6 border-t">
              <h3 className="font-medium mb-3">Interests</h3>
              <div className="flex flex-wrap gap-2">
                {profile.interests.map((interest, index) => (
                  <Badge key={index} variant="secondary">
                    {interest}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          {/* Social Links */}
          {profile.socialLinks && Object.keys(profile.socialLinks).length > 0 && (
            <div className="mt-6 pt-6 border-t">
              <h3 className="font-medium mb-3">Social Links</h3>
              <div className="flex flex-wrap gap-2">
                {Object.entries(profile.socialLinks).map(([platform, url]) => (
                  <Button
                    key={platform}
                    variant="outline"
                    size="sm"
                    onClick={() => window.open(url as string, '_blank')}
                  >
                    <Link className="w-4 h-4 mr-2" />
                    {platform}
                  </Button>
                ))}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Analytics (for own profile) */}
      {isOwnProfile && analytics && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="w-5 h-5" />
              Profile Analytics
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-chart-1">
                  {(analytics.engagement as any)?.totalLikes || 0}
                </div>
                <div className="text-sm text-muted-foreground">Total Likes</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-chart-2">
                  {(analytics.engagement as any)?.totalComments || 0}
                </div>
                <div className="text-sm text-muted-foreground">Comments</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-chart-3">
                  {(analytics.reach as any)?.profileViews || 0}
                </div>
                <div className="text-sm text-muted-foreground">Profile Views</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-chart-4">
                  {(analytics.growth as any)?.followerGrowth || 0}
                </div>
                <div className="text-sm text-muted-foreground">Growth This Month</div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Edit Profile Dialog */}
      <EditProfileDialog
        profile={profile}
        open={editDialogOpen}
        onOpenChange={setEditDialogOpen}
        onSave={handleUpdateProfile}
      />

      {/* Followers Dialog */}
      <Dialog open={followersDialogOpen} onOpenChange={setFollowersDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Followers</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 max-h-96 overflow-y-auto">
            {followers?.followers.map((follower) => (
              <div key={follower.id} className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Avatar>
                    <AvatarImage src={follower.avatar} />
                    <AvatarFallback>{getSocialDisplayName(follower).charAt(0)}</AvatarFallback>
                  </Avatar>
                  <div>
                    <div className="font-medium">{getSocialDisplayName(follower)}</div>
                    {follower.username && (
                      <div className="text-sm text-muted-foreground">@{follower.username}</div>
                    )}
                  </div>
                </div>
                <Button size="sm" variant="outline">
                  View Profile
                </Button>
              </div>
            ))}
          </div>
        </DialogContent>
      </Dialog>

      {/* Following Dialog */}
      <Dialog open={followingDialogOpen} onOpenChange={setFollowingDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Following</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 max-h-96 overflow-y-auto">
            {following?.following.map((followedUser) => (
              <div key={followedUser.id} className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <Avatar>
                    <AvatarImage src={followedUser.avatar} />
                    <AvatarFallback>{getSocialDisplayName(followedUser).charAt(0)}</AvatarFallback>
                  </Avatar>
                  <div>
                    <div className="font-medium">{getSocialDisplayName(followedUser)}</div>
                    {followedUser.username && (
                      <div className="text-sm text-muted-foreground">@{followedUser.username}</div>
                    )}
                  </div>
                </div>
                <Button size="sm" variant="outline">
                  View Profile
                </Button>
              </div>
            ))}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default SocialProfile;