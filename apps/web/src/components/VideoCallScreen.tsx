import React, { useState, useEffect } from 'react';
import { 
  Phone, 
  PhoneOff, 
  Mic, 
  MicOff, 
  Video, 
  VideoOff, 
  Camera, 
  Speaker,
  MessageSquare,
  Users,
  Settings,
  Maximize2,
  Minimize2,
  ArrowLeft
} from 'lucide-react';
import { Button } from './ui/button';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';

interface VideoCallScreenProps {
  user: any;
  callee: {
    id: string;
    name: string;
    avatar?: string;
    isGroup?: boolean;
    members?: number;
  };
  isIncoming?: boolean;
  onEndCall: () => void;
  onBack: () => void;
}

export function VideoCallScreen({ user, callee, isIncoming = false, onEndCall, onBack }: VideoCallScreenProps) {
  const [isConnected, setIsConnected] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [isVideoOn, setIsVideoOn] = useState(true);
  const [isSpeakerOn, setIsSpeakerOn] = useState(false);
  const [callDuration, setCallDuration] = useState(0);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isControlsVisible, setIsControlsVisible] = useState(true);

  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (isConnected) {
      interval = setInterval(() => {
        setCallDuration(prev => prev + 1);
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [isConnected]);

  useEffect(() => {
    if (!isIncoming) {
      // Auto-connect for outgoing calls after 3 seconds
      const timer = setTimeout(() => {
        setIsConnected(true);
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [isIncoming]);

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const handleAnswer = () => {
    setIsConnected(true);
  };

  const handleDecline = () => {
    onEndCall();
  };

  const handleEndCall = () => {
    onEndCall();
  };

  const toggleControls = () => {
    setIsControlsVisible(!isControlsVisible);
  };

  if (isIncoming && !isConnected) {
    return (
      <div className="h-screen bg-gradient-to-br from-chart-1/20 to-chart-2/20 flex flex-col items-center justify-center relative">
        {/* Background blur effect */}
        <div className="absolute inset-0 bg-black/20 backdrop-blur-sm"></div>
        
        <div className="relative z-10 text-center space-y-8">
          {/* Caller info */}
          <div className="space-y-4">
            <Avatar className="w-32 h-32 mx-auto ring-4 ring-white/20">
              <AvatarImage src={callee.avatar} />
              <AvatarFallback className="text-2xl">
                {callee.isGroup ? <Users className="w-12 h-12" /> : callee.name.charAt(0)}
              </AvatarFallback>
            </Avatar>
            
            <div className="space-y-2">
              <h2 className="text-2xl text-white">{callee.name}</h2>
              {callee.isGroup && (
                <Badge variant="secondary" className="bg-white/20 text-white">
                  {callee.members} members
                </Badge>
              )}
              <p className="text-white/80">Incoming video call...</p>
            </div>
          </div>

          {/* Call actions */}
          <div className="flex items-center justify-center gap-8">
            <Button
              size="lg"
              variant="destructive"
              className="w-16 h-16 rounded-full bg-red-500 hover:bg-red-600"
              onClick={handleDecline}
            >
              <PhoneOff className="w-8 h-8" />
            </Button>
            
            <Button
              size="lg"
              className="w-16 h-16 rounded-full bg-green-500 hover:bg-green-600"
              onClick={handleAnswer}
            >
              <Phone className="w-8 h-8" />
            </Button>
          </div>

          {/* Quick actions */}
          <div className="flex items-center justify-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              className="text-white hover:bg-white/20"
              onClick={() => setIsMuted(!isMuted)}
            >
              {isMuted ? <MicOff className="w-5 h-5" /> : <Mic className="w-5 h-5" />}
            </Button>
            
            <Button
              variant="ghost"
              size="icon"
              className="text-white hover:bg-white/20"
              onClick={() => setIsVideoOn(!isVideoOn)}
            >
              {isVideoOn ? <Video className="w-5 h-5" /> : <VideoOff className="w-5 h-5" />}
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div 
      className={`h-screen bg-black relative overflow-hidden ${isFullscreen ? 'cursor-none' : ''}`}
      onClick={toggleControls}
    >
      {/* Main video area */}
      <div className="absolute inset-0 flex items-center justify-center">
        {/* Remote video placeholder */}
        <div className="w-full h-full bg-gradient-to-br from-muted to-muted-foreground/20 flex items-center justify-center">
          {!isVideoOn ? (
            <div className="text-center">
              <Avatar className="w-32 h-32 mx-auto mb-4">
                <AvatarImage src={callee.avatar} />
                <AvatarFallback className="text-2xl">
                  {callee.isGroup ? <Users className="w-12 h-12" /> : callee.name.charAt(0)}
                </AvatarFallback>
              </Avatar>
              <p className="text-white text-lg">{callee.name}</p>
              <p className="text-white/60">Camera is off</p>
            </div>
          ) : (
            <div className="w-full h-full bg-gradient-to-br from-chart-1/30 to-chart-2/30 flex items-center justify-center">
              <p className="text-white/60">Video feed would appear here</p>
            </div>
          )}
        </div>

        {/* Local video (Picture-in-Picture) */}
        <Card className={`absolute top-4 right-4 w-32 h-24 bg-muted border-2 border-white/20 overflow-hidden transition-all duration-300 ${isControlsVisible ? 'opacity-100' : 'opacity-0'}`}>
          <CardContent className="p-0 h-full relative">
            {isVideoOn ? (
              <div className="w-full h-full bg-gradient-to-br from-primary/30 to-primary/60 flex items-center justify-center">
                <p className="text-white text-xs">You</p>
              </div>
            ) : (
              <div className="w-full h-full bg-muted flex items-center justify-center">
                <VideoOff className="w-6 h-6 text-muted-foreground" />
              </div>
            )}
            <Button
              size="icon"
              variant="ghost"
              className="absolute top-1 right-1 w-6 h-6 hover:bg-black/20"
              onClick={() => setIsVideoOn(!isVideoOn)}
            >
              <Camera className="w-3 h-3 text-white" />
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Call info header */}
      <div className={`absolute top-0 left-0 right-0 p-4 bg-gradient-to-b from-black/50 to-transparent transition-all duration-300 ${isControlsVisible ? 'opacity-100' : 'opacity-0'}`}>
        <div className="flex items-center justify-between text-white">
          <Button
            variant="ghost"
            size="icon"
            className="text-white hover:bg-white/20"
            onClick={onBack}
          >
            <ArrowLeft className="w-5 h-5" />
          </Button>
          
          <div className="text-center">
            <p className="text-sm">{callee.name}</p>
            <p className="text-xs text-white/60">
              {isConnected ? formatDuration(callDuration) : 'Connecting...'}
            </p>
          </div>

          <Button
            variant="ghost"
            size="icon"
            className="text-white hover:bg-white/20"
            onClick={() => setIsFullscreen(!isFullscreen)}
          >
            {isFullscreen ? <Minimize2 className="w-5 h-5" /> : <Maximize2 className="w-5 h-5" />}
          </Button>
        </div>
      </div>

      {/* Call controls */}
      <div className={`absolute bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-black/70 to-transparent transition-all duration-300 ${isControlsVisible ? 'opacity-100' : 'opacity-0'}`}>
        <div className="flex items-center justify-center gap-4">
          {/* Mute button */}
          <Button
            size="lg"
            variant={isMuted ? "destructive" : "secondary"}
            className="w-14 h-14 rounded-full"
            onClick={() => setIsMuted(!isMuted)}
          >
            {isMuted ? <MicOff className="w-6 h-6" /> : <Mic className="w-6 h-6" />}
          </Button>

          {/* Video toggle */}
          <Button
            size="lg"
            variant={!isVideoOn ? "destructive" : "secondary"}
            className="w-14 h-14 rounded-full"
            onClick={() => setIsVideoOn(!isVideoOn)}
          >
            {isVideoOn ? <Video className="w-6 h-6" /> : <VideoOff className="w-6 h-6" />}
          </Button>

          {/* End call */}
          <Button
            size="lg"
            variant="destructive"
            className="w-16 h-16 rounded-full bg-red-500 hover:bg-red-600"
            onClick={handleEndCall}
          >
            <PhoneOff className="w-8 h-8" />
          </Button>

          {/* Speaker toggle */}
          <Button
            size="lg"
            variant={isSpeakerOn ? "default" : "secondary"}
            className="w-14 h-14 rounded-full"
            onClick={() => setIsSpeakerOn(!isSpeakerOn)}
          >
            <Speaker className="w-6 h-6" />
          </Button>

          {/* Chat */}
          <Button
            size="lg"
            variant="secondary"
            className="w-14 h-14 rounded-full"
          >
            <MessageSquare className="w-6 h-6" />
          </Button>
        </div>

        {/* Additional controls */}
        <div className="flex items-center justify-center gap-2 mt-4">
          <Button
            variant="ghost"
            size="sm"
            className="text-white hover:bg-white/20"
          >
            <Settings className="w-4 h-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Connection status */}
      {!isConnected && (
        <div className="absolute bottom-20 left-1/2 transform -translate-x-1/2">
          <Badge variant="secondary" className="bg-black/50 text-white">
            Connecting...
          </Badge>
        </div>
      )}
    </div>
  );
}