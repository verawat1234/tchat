import React, { useState, useRef, useEffect, useCallback } from 'react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Slider } from './ui/slider';
import { 
  Play, 
  Pause, 
  Volume2, 
  VolumeX, 
  Maximize, 
  Minimize,
  X,
  Heart,
  MessageCircle,
  Share,
  Bookmark,
  UserPlus,
  UserMinus,
  MoreVertical,
  SkipBack,
  SkipForward,
  Settings,
  Download,
  Flag,
  CheckCircle
} from 'lucide-react';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { toast } from "sonner";

interface Video {
  id: string;
  title: string;
  description: string;
  videoUrl: string;
  thumbnail: string;
  duration: string;
  views: number;
  likes: number;
  uploadTime: string;
  channel: {
    id: string;
    name: string;
    avatar?: string;
    verified?: boolean;
    subscribers?: number;
  };
  tags?: string[];
  quality?: string[];
}

interface FullscreenVideoPlayerProps {
  video: Video;
  isPlaying: boolean;
  onClose: () => void;
  onPlay: () => void;
  onPause: () => void;
  onLike: (videoId: string) => void;
  onShare: (videoId: string) => void;
  onSubscribe: (channelId: string) => void;
  isLiked: boolean;
  isSubscribed: boolean;
  user: any;
}

export function FullscreenVideoPlayer({
  video,
  isPlaying,
  onClose,
  onPlay,
  onPause,
  onLike,
  onShare,
  onSubscribe,
  isLiked,
  isSubscribed,
  user
}: FullscreenVideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1);
  const [isMuted, setIsMuted] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [showControls, setShowControls] = useState(true);
  const [isBuffering, setIsBuffering] = useState(false);
  const [quality, setQuality] = useState('auto');
  const [playbackRate, setPlaybackRate] = useState(1);
  const [showInfo, setShowInfo] = useState(false);
  
  // Auto-hide controls
  const controlsTimeoutRef = useRef<NodeJS.Timeout>();
  
  const resetControlsTimeout = useCallback(() => {
    if (controlsTimeoutRef.current) {
      clearTimeout(controlsTimeoutRef.current);
    }
    setShowControls(true);
    controlsTimeoutRef.current = setTimeout(() => {
      if (isPlaying) {
        setShowControls(false);
      }
    }, 3000);
  }, [isPlaying]);

  // Video event handlers
  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    const handleTimeUpdate = () => setCurrentTime(video.currentTime);
    const handleDurationChange = () => setDuration(video.duration);
    const handleVolumeChange = () => {
      setVolume(video.volume);
      setIsMuted(video.muted);
    };
    const handleWaiting = () => setIsBuffering(true);
    const handleCanPlay = () => setIsBuffering(false);
    const handlePlaying = () => setIsBuffering(false);

    video.addEventListener('timeupdate', handleTimeUpdate);
    video.addEventListener('durationchange', handleDurationChange);
    video.addEventListener('volumechange', handleVolumeChange);
    video.addEventListener('waiting', handleWaiting);
    video.addEventListener('canplay', handleCanPlay);
    video.addEventListener('playing', handlePlaying);

    return () => {
      video.removeEventListener('timeupdate', handleTimeUpdate);
      video.removeEventListener('durationchange', handleDurationChange);
      video.removeEventListener('volumechange', handleVolumeChange);
      video.removeEventListener('waiting', handleWaiting);
      video.removeEventListener('canplay', handleCanPlay);
      video.removeEventListener('playing', handlePlaying);
    };
  }, []);

  // Handle play/pause
  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    if (isPlaying) {
      video.play().catch(console.error);
    } else {
      video.pause();
    }
  }, [isPlaying]);

  // Fullscreen handling
  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => document.removeEventListener('fullscreenchange', handleFullscreenChange);
  }, []);

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      switch (e.code) {
        case 'Space':
          e.preventDefault();
          isPlaying ? onPause() : onPlay();
          break;
        case 'ArrowLeft':
          e.preventDefault();
          seek(currentTime - 10);
          break;
        case 'ArrowRight':
          e.preventDefault();
          seek(currentTime + 10);
          break;
        case 'ArrowUp':
          e.preventDefault();
          setVolume(Math.min(1, volume + 0.1));
          break;
        case 'ArrowDown':
          e.preventDefault();
          setVolume(Math.max(0, volume - 0.1));
          break;
        case 'KeyM':
          e.preventDefault();
          toggleMute();
          break;
        case 'KeyF':
          e.preventDefault();
          toggleFullscreen();
          break;
        case 'Escape':
          e.preventDefault();
          onClose();
          break;
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [currentTime, volume, isPlaying, onPause, onPlay, onClose]);

  // Touch gestures for mobile
  const [touchStart, setTouchStart] = useState<{ x: number; y: number } | null>(null);
  
  const handleTouchStart = (e: React.TouchEvent) => {
    const touch = e.touches[0];
    setTouchStart({ x: touch.clientX, y: touch.clientY });
  };

  const handleTouchEnd = (e: React.TouchEvent) => {
    if (!touchStart) return;
    
    const touch = e.changedTouches[0];
    const deltaX = touch.clientX - touchStart.x;
    const deltaY = touch.clientY - touchStart.y;
    
    // Vertical swipe down to close
    if (deltaY > 100 && Math.abs(deltaX) < 50) {
      onClose();
    }
    
    // Horizontal double tap to seek
    if (Math.abs(deltaX) < 10 && Math.abs(deltaY) < 10) {
      const rect = e.currentTarget.getBoundingClientRect();
      const tapX = touch.clientX - rect.left;
      const tapArea = rect.width / 3;
      
      if (tapX < tapArea) {
        seek(currentTime - 10);
        toast.success('⏪ -10s');
      } else if (tapX > tapArea * 2) {
        seek(currentTime + 10);
        toast.success('⏩ +10s');
      } else {
        isPlaying ? onPause() : onPlay();
      }
    }
    
    setTouchStart(null);
  };

  const seek = (time: number) => {
    if (videoRef.current) {
      videoRef.current.currentTime = Math.max(0, Math.min(duration, time));
    }
  };

  const toggleMute = () => {
    if (videoRef.current) {
      videoRef.current.muted = !videoRef.current.muted;
    }
  };

  const toggleFullscreen = async () => {
    if (!containerRef.current) return;
    
    try {
      if (!document.fullscreenElement) {
        await containerRef.current.requestFullscreen();
      } else {
        await document.exitFullscreen();
      }
    } catch (error) {
      console.error('Fullscreen error:', error);
    }
  };

  const formatTime = (time: number) => {
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  const formatViews = (views: number) => {
    if (views >= 1000000) {
      return `${(views / 1000000).toFixed(1)}M`;
    } else if (views >= 1000) {
      return `${(views / 1000).toFixed(1)}K`;
    }
    return views.toString();
  };

  return (
    <div 
      ref={containerRef}
      className="fixed inset-0 z-[9999] bg-black flex flex-col"
      onMouseMove={resetControlsTimeout}
      onTouchStart={handleTouchStart}
      onTouchEnd={handleTouchEnd}
      onClick={resetControlsTimeout}
    >
      {/* Video Element */}
      <div className="relative flex-1 flex items-center justify-center">
        <video
          ref={videoRef}
          src={video.videoUrl}
          className="w-full h-full object-contain"
          poster={video.thumbnail}
          preload="metadata"
          onLoadedMetadata={() => {
            if (videoRef.current) {
              setDuration(videoRef.current.duration);
            }
          }}
        />
        
        {/* Buffering Indicator */}
        {isBuffering && (
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="w-12 h-12 border-4 border-white/20 border-t-white rounded-full animate-spin"></div>
          </div>
        )}
        
        {/* Close Button */}
        <Button
          variant="ghost"
          size="icon"
          className="absolute top-4 right-4 z-50 bg-black/50 hover:bg-black/70 text-white rounded-full w-10 h-10 sm:w-12 sm:h-12 touch-manipulation"
          onClick={onClose}
        >
          <X className="w-5 h-5 sm:w-6 sm:h-6" />
        </Button>

        {/* Video Info Toggle */}
        <Button
          variant="ghost"
          size="icon"
          className="absolute top-4 left-4 z-50 bg-black/50 hover:bg-black/70 text-white rounded-full w-10 h-10 sm:w-12 sm:h-12 touch-manipulation lg:hidden"
          onClick={() => setShowInfo(!showInfo)}
        >
          <MoreVertical className="w-5 h-5 sm:w-6 sm:h-6" />
        </Button>
      </div>

      {/* Controls Overlay */}
      <div className={`absolute inset-0 pointer-events-none transition-opacity duration-300 ${showControls ? 'opacity-100' : 'opacity-0'}`}>
        {/* Top Gradient */}
        <div className="absolute top-0 left-0 right-0 h-32 bg-gradient-to-b from-black/70 to-transparent"></div>
        
        {/* Bottom Controls */}
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/90 via-black/60 to-transparent p-4 sm:p-6 pointer-events-auto">
          {/* Progress Bar */}
          <div className="mb-4">
            <Slider
              value={[currentTime]}
              max={duration || 100}
              step={1}
              onValueChange={(value) => seek(value[0])}
              className="w-full cursor-pointer"
            />
            <div className="flex justify-between text-xs text-white/80 mt-1">
              <span>{formatTime(currentTime)}</span>
              <span>{formatTime(duration)}</span>
            </div>
          </div>
          
          {/* Control Buttons */}
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 sm:gap-4">
              {/* Play/Pause */}
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/20 rounded-full w-12 h-12 sm:w-14 sm:h-14 touch-manipulation"
                onClick={isPlaying ? onPause : onPlay}
              >
                {isPlaying ? <Pause className="w-6 h-6 sm:w-8 sm:h-8" /> : <Play className="w-6 h-6 sm:w-8 sm:h-8" />}
              </Button>
              
              {/* Skip Buttons */}
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/20 rounded-full w-10 h-10 sm:w-12 sm:h-12 touch-manipulation"
                onClick={() => seek(currentTime - 10)}
              >
                <SkipBack className="w-5 h-5 sm:w-6 sm:h-6" />
              </Button>
              
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/20 rounded-full w-10 h-10 sm:w-12 sm:h-12 touch-manipulation"
                onClick={() => seek(currentTime + 10)}
              >
                <SkipForward className="w-5 h-5 sm:w-6 sm:h-6" />
              </Button>
              
              {/* Volume Control - Desktop Only */}
              <div className="hidden lg:flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-white hover:bg-white/20 rounded-full w-10 h-10"
                  onClick={toggleMute}
                >
                  {isMuted || volume === 0 ? <VolumeX className="w-5 h-5" /> : <Volume2 className="w-5 h-5" />}
                </Button>
                <Slider
                  value={[isMuted ? 0 : volume * 100]}
                  max={100}
                  step={1}
                  onValueChange={(value) => {
                    const newVolume = value[0] / 100;
                    setVolume(newVolume);
                    if (videoRef.current) {
                      videoRef.current.volume = newVolume;
                      videoRef.current.muted = newVolume === 0;
                    }
                  }}
                  className="w-20"
                />
              </div>
            </div>
            
            <div className="flex items-center gap-2">
              {/* Settings Menu */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-white hover:bg-white/20 rounded-full w-10 h-10 sm:w-12 sm:h-12 touch-manipulation"
                  >
                    <Settings className="w-5 h-5 sm:w-6 sm:h-6" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-48 z-[10000]">
                  <DropdownMenuItem>
                    <span>Quality: {quality}</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setPlaybackRate(0.5)}>
                    <span>Speed: 0.5x</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setPlaybackRate(1)}>
                    <span>Speed: 1x</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setPlaybackRate(1.5)}>
                    <span>Speed: 1.5x</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setPlaybackRate(2)}>
                    <span>Speed: 2x</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
              
              {/* Fullscreen Toggle */}
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/20 rounded-full w-10 h-10 sm:w-12 sm:h-12 touch-manipulation"
                onClick={toggleFullscreen}
              >
                {isFullscreen ? <Minimize className="w-5 h-5 sm:w-6 sm:h-6" /> : <Maximize className="w-5 h-5 sm:w-6 sm:h-6" />}
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Video Info Sidebar - Desktop Only or Mobile Popup */}
      <div className={`absolute top-0 right-0 bottom-0 w-80 xl:w-96 bg-black/90 backdrop-blur-sm overflow-y-auto p-4 sm:p-6 transition-transform duration-300 ${showInfo || window.innerWidth >= 1024 ? 'translate-x-0' : 'translate-x-full'} lg:translate-x-0`}>
        {/* Video Info */}
        <div className="space-y-4">
          <div>
            <h2 className="text-white text-lg sm:text-xl font-medium mb-2 line-clamp-3">{video.title}</h2>
            <div className="flex items-center gap-2 text-white/70 text-sm">
              <span>{formatViews(video.views)} views</span>
              <span>•</span>
              <span>{video.uploadTime}</span>
            </div>
          </div>

          {/* Channel Info */}
          <div className="flex items-center gap-3">
            <Avatar className="w-12 h-12">
              <AvatarImage src={video.channel.avatar} />
              <AvatarFallback className="bg-white/20 text-white">{video.channel.name.charAt(0)}</AvatarFallback>
            </Avatar>
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <span className="text-white font-medium">{video.channel.name}</span>
                {video.channel.verified && (
                  <CheckCircle className="w-4 h-4 text-blue-500 fill-current" />
                )}
              </div>
              {video.channel.subscribers && (
                <p className="text-white/70 text-sm">{formatViews(video.channel.subscribers)} subscribers</p>
              )}
            </div>
            <Button
              variant={isSubscribed ? "secondary" : "default"}
              size="sm"
              onClick={() => onSubscribe(video.channel.id)}
              className="touch-manipulation"
            >
              {isSubscribed ? (
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

          {/* Action Buttons */}
          <div className="flex flex-wrap gap-2">
            <Button
              variant="ghost"
              size="sm"
              className={`text-white hover:bg-white/20 ${isLiked ? 'text-red-500' : ''} touch-manipulation`}
              onClick={() => onLike(video.id)}
            >
              <Heart className={`w-4 h-4 mr-2 ${isLiked ? 'fill-current' : ''}`} />
              {formatViews(video.likes)}
            </Button>
            
            <Button
              variant="ghost"
              size="sm"
              className="text-white hover:bg-white/20 touch-manipulation"
              onClick={() => onShare(video.id)}
            >
              <Share className="w-4 h-4 mr-2" />
              Share
            </Button>
            
            <Button
              variant="ghost"
              size="sm"
              className="text-white hover:bg-white/20 touch-manipulation"
            >
              <Bookmark className="w-4 h-4 mr-2" />
              Save
            </Button>
            
            <Button
              variant="ghost"
              size="sm"
              className="text-white hover:bg-white/20 touch-manipulation"
            >
              <Download className="w-4 h-4 mr-2" />
              Download
            </Button>
          </div>

          {/* Description */}
          <div className="text-white/80 text-sm">
            <p className="leading-relaxed">{video.description}</p>
          </div>

          {/* Tags */}
          {video.tags && video.tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {video.tags.map((tag) => (
                <Badge key={tag} variant="secondary" className="bg-white/20 text-white hover:bg-white/30 cursor-pointer">
                  #{tag}
                </Badge>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}