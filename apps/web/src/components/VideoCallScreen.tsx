import React, { useState, useEffect, useRef, useCallback } from 'react';
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
  ArrowLeft,
  MonitorSpeaker
} from 'lucide-react';
import { Button } from './ui/button';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { callService, CallState, CallSession } from '../services/webrtc/CallService';

interface VideoCallScreenProps {
  user: any;
  callee: {
    id: string;
    name: string;
    avatar?: string;
    isGroup?: boolean;
    members?: number;
  };
  callId?: string;
  isIncoming?: boolean;
  onEndCall: () => void;
  onBack: () => void;
  onCallAccept?: () => void;
  onCallDecline?: () => void;
}

export function VideoCallScreen({
  user,
  callee,
  callId,
  isIncoming = false,
  onEndCall,
  onBack,
  onCallAccept,
  onCallDecline
}: VideoCallScreenProps) {
  // WebRTC state
  const [callState, setCallState] = useState<CallState>('idle');
  const [currentCall, setCurrentCall] = useState<CallSession | null>(null);
  const [isAudioMuted, setIsAudioMuted] = useState(false);
  const [isVideoMuted, setIsVideoMuted] = useState(false);
  const [isScreenSharing, setIsScreenSharing] = useState(false);
  const [networkQuality, setNetworkQuality] = useState<'poor' | 'fair' | 'good' | 'excellent'>('good');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Legacy UI state for compatibility
  const [isConnected, setIsConnected] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [isVideoOn, setIsVideoOn] = useState(true);
  const [isSpeakerOn, setIsSpeakerOn] = useState(false);
  const [callDuration, setCallDuration] = useState(0);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isControlsVisible, setIsControlsVisible] = useState(true);

  // Video element refs
  const localVideoRef = useRef<HTMLVideoElement>(null);
  const remoteVideoRef = useRef<HTMLVideoElement>(null);
  const durationIntervalRef = useRef<NodeJS.Timeout | null>(null);

  // Initialize WebRTC call service
  useEffect(() => {
    const initializeCallService = async () => {
      try {
        await callService.initialize();

        callService.setEventHandlers({
          onStateChange: (state) => {
            setCallState(state);
            setIsConnected(state === 'connected');

            if (state === 'ended' || state === 'error') {
              stopDurationTimer();
              onEndCall?.();
            } else if (state === 'connected') {
              startDurationTimer();
            }
          },
          onMediaToggle: (participantId, type, enabled) => {
            if (participantId === 'local') {
              if (type === 'audio') {
                setIsAudioMuted(!enabled);
                setIsMuted(!enabled);
              } else if (type === 'video') {
                setIsVideoMuted(!enabled);
                setIsVideoOn(enabled);
              }
            }
          },
          onScreenShareStart: (participantId) => {
            if (participantId === 'local') {
              setIsScreenSharing(true);
            }
          },
          onScreenShareEnd: (participantId) => {
            if (participantId === 'local') {
              setIsScreenSharing(false);
            }
          },
          onNetworkQualityChange: (quality) => {
            setNetworkQuality(quality);
          },
          onError: (error) => {
            setError(error.message);
            setIsLoading(false);
          }
        });

      } catch (error) {
        setError('Failed to initialize call service');
        console.error('Call service initialization failed:', error);
      }
    };

    initializeCallService();

    return () => {
      stopDurationTimer();
    };
  }, [onEndCall]);

  // Update current call state
  useEffect(() => {
    const call = callService.getCurrentCall();
    setCurrentCall(call);
  }, [callState]);

  // Setup video streams when call connects
  useEffect(() => {
    if (callState === 'connected') {
      const localStream = callService.getLocalStream();
      const remoteStream = callService.getRemoteStream();

      if (localStream && localVideoRef.current) {
        localVideoRef.current.srcObject = localStream;
        localVideoRef.current.play().catch(console.error);
      }

      if (remoteStream && remoteVideoRef.current) {
        remoteVideoRef.current.srcObject = remoteStream;
        remoteVideoRef.current.play().catch(console.error);
      }
    }
  }, [callState]);

  // Handle incoming call setup
  useEffect(() => {
    if (isIncoming && callId && !currentCall) {
      setCallState('ringing');
    }
  }, [isIncoming, callId, currentCall]);

  const startDurationTimer = useCallback(() => {
    if (durationIntervalRef.current) return;

    durationIntervalRef.current = setInterval(() => {
      setCallDuration(prev => prev + 1);
    }, 1000);
  }, []);

  const stopDurationTimer = useCallback(() => {
    if (durationIntervalRef.current) {
      clearInterval(durationIntervalRef.current);
      durationIntervalRef.current = null;
    }
  }, []);

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const handleAnswer = useCallback(async () => {
    if (!callId) return;

    setIsLoading(true);
    setError(null);

    try {
      await callService.answerCall(callId, {
        type: 'video',
        enableVideo: true,
        enableAudio: true,
        quality: 'auto'
      });

      onCallAccept?.();
    } catch (error) {
      setError('Failed to accept call');
      console.error('Failed to accept call:', error);
    } finally {
      setIsLoading(false);
    }
  }, [callId, onCallAccept]);

  const handleDecline = useCallback(async () => {
    if (!callId) return;

    try {
      await callService.declineCall(callId);
      onCallDecline?.();
    } catch (error) {
      console.error('Failed to decline call:', error);
      onCallDecline?.();
    }
  }, [callId, onCallDecline]);

  const handleEndCall = useCallback(async () => {
    try {
      await callService.endCall();
    } catch (error) {
      console.error('Failed to end call:', error);
    }
  }, []);

  const handleToggleAudio = useCallback(() => {
    const newState = callService.toggleAudio();
    setIsMuted(!newState);
    setIsAudioMuted(!newState);
  }, []);

  const handleToggleVideo = useCallback(() => {
    const newState = callService.toggleVideo();
    setIsVideoOn(newState);
    setIsVideoMuted(!newState);
  }, []);

  const handleSwitchCamera = useCallback(async () => {
    try {
      await callService.switchCamera();
    } catch (error) {
      console.error('Failed to switch camera:', error);
      setError('Failed to switch camera');
    }
  }, []);

  const handleToggleScreenShare = useCallback(async () => {
    try {
      if (isScreenSharing) {
        await callService.stopScreenShare();
      } else {
        await callService.startScreenShare();
      }
    } catch (error) {
      console.error('Failed to toggle screen share:', error);
      setError('Failed to toggle screen sharing');
    }
  }, [isScreenSharing]);

  const toggleControls = () => {
    setIsControlsVisible(!isControlsVisible);
  };

  if (isIncoming && (callState === 'ringing' || !isConnected)) {
    return (
      <div
        className="h-screen bg-gradient-to-br from-chart-1/20 to-chart-2/20 flex flex-col items-center justify-center relative"
        data-testid="incoming-call-modal"
        role="dialog"
        aria-modal="true"
        aria-labelledby="incoming-call-title"
      >
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
              <h2
                id="incoming-call-title"
                className="text-2xl text-white"
                data-testid="caller-name"
              >
                {callee.name}
              </h2>
              {callee.isGroup && (
                <Badge variant="secondary" className="bg-white/20 text-white">
                  {callee.members} members
                </Badge>
              )}
              <p
                className="text-white/80"
                data-testid="call-type"
              >
                Video Call
              </p>
            </div>
          </div>

          {/* Error message */}
          {error && (
            <div className="mb-4 p-3 bg-red-500 bg-opacity-90 rounded-lg text-white">
              {error}
            </div>
          )}

          {/* Call actions */}
          <div className="flex items-center justify-center gap-8">
            <Button
              size="lg"
              variant="destructive"
              className="w-16 h-16 rounded-full bg-red-500 hover:bg-red-600"
              onClick={handleDecline}
              disabled={isLoading}
              aria-label="Decline call"
              data-testid="decline-call-button"
            >
              <PhoneOff className="w-8 h-8" />
            </Button>

            <Button
              size="lg"
              className="w-16 h-16 rounded-full bg-green-500 hover:bg-green-600"
              onClick={handleAnswer}
              disabled={isLoading}
              aria-label="Accept call"
              data-testid="accept-video-call-button"
            >
              <Phone className="w-8 h-8" />
            </Button>
          </div>

          {/* Loading indicator */}
          {isLoading && (
            <div className="mt-4 text-white">
              Connecting...
            </div>
          )}

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
      data-testid="video-call-screen"
    >
      {/* Network warning */}
      {networkQuality === 'poor' && (
        <div
          className="absolute top-4 left-4 bg-red-500 bg-opacity-90 px-3 py-2 rounded-lg text-sm text-white z-30"
          data-testid="network-warning"
        >
          Poor network connection
        </div>
      )}

      {/* Error message */}
      {error && (
        <div
          className="absolute top-4 right-4 bg-red-500 bg-opacity-90 px-3 py-2 rounded-lg text-sm text-white z-30 max-w-xs"
          data-testid="call-error-message"
        >
          {error}
        </div>
      )}

      {/* Call timeout message */}
      {callState === 'error' && error?.includes('timeout') && (
        <div
          className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-yellow-500 bg-opacity-90 px-4 py-3 rounded-lg text-center z-40"
          data-testid="call-timeout-message"
        >
          Call timed out
        </div>
      )}

      {/* Screen share indicator */}
      {isScreenSharing && (
        <div
          className="absolute top-4 left-1/2 transform -translate-x-1/2 bg-blue-500 bg-opacity-90 px-3 py-2 rounded-lg text-sm text-white z-30"
          data-testid="screen-share-indicator"
        >
          Sharing screen
        </div>
      )}

      {/* Main video area */}
      <div className="absolute inset-0 flex items-center justify-center">
        {/* Remote video */}
        <video
          ref={remoteVideoRef}
          autoPlay
          playsInline
          className="w-full h-full object-cover"
          style={{ display: callState === 'connected' && !isVideoMuted ? 'block' : 'none' }}
          data-testid="remote-video"
          aria-label="Remote video stream"
        />

        {/* Remote video placeholder when video is off */}
        {(!isConnected || isVideoMuted || callState !== 'connected') && (
          <div className="w-full h-full bg-gradient-to-br from-muted to-muted-foreground/20 flex items-center justify-center">
            <div className="text-center">
              <Avatar className="w-32 h-32 mx-auto mb-4">
                <AvatarImage src={callee.avatar} />
                <AvatarFallback className="text-2xl">
                  {callee.isGroup ? <Users className="w-12 h-12" /> : callee.name.charAt(0)}
                </AvatarFallback>
              </Avatar>
              <p className="text-white text-lg">{callee.name}</p>
              <p className="text-white/60">
                {callState !== 'connected' ? 'Connecting...' : 'Camera is off'}
              </p>
            </div>
          </div>
        )}

        {/* Local video (Picture-in-Picture) */}
        <Card className={`absolute top-4 right-4 w-32 h-24 bg-muted border-2 border-white/20 overflow-hidden transition-all duration-300 ${isControlsVisible ? 'opacity-100' : 'opacity-0'}`}>
          <CardContent className="p-0 h-full relative">
            {/* Local video stream */}
            <video
              ref={localVideoRef}
              autoPlay
              playsInline
              muted
              className="w-full h-full object-cover"
              style={{ display: isVideoOn && callState === 'connected' ? 'block' : 'none' }}
              data-testid="local-video"
              aria-label="Local video stream"
            />

            {/* Local video placeholder when video is off */}
            {(!isVideoOn || callState !== 'connected') && (
              <div className="w-full h-full bg-muted flex items-center justify-center">
                {callState !== 'connected' ? (
                  <p className="text-white text-xs">Connecting...</p>
                ) : (
                  <VideoOff className="w-6 h-6 text-muted-foreground" />
                )}
              </div>
            )}

            {/* Camera switch button */}
            <Button
              size="icon"
              variant="ghost"
              className="absolute top-1 right-1 w-6 h-6 hover:bg-black/20"
              onClick={handleSwitchCamera}
              disabled={!isVideoOn || callState !== 'connected'}
              aria-label="Switch camera"
              data-testid="switch-camera-button"
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
            <p
              className="text-xs text-white/60"
              data-testid="call-status"
            >
              {callState === 'connected' ? formatDuration(callDuration) :
               callState === 'calling' ? 'Calling...' :
               callState === 'ringing' ? 'Ringing...' : 'Connecting...'}
            </p>
            {callState === 'connected' && (
              <p
                className="text-xs text-white/40 mt-1"
                data-testid="call-duration"
              >
                Duration: {formatDuration(callDuration)}
              </p>
            )}
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
      <div
        className={`absolute bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-black/70 to-transparent transition-all duration-300 ${isControlsVisible ? 'opacity-100' : 'opacity-0'}`}
        data-testid="video-controls"
      >
        <div className="flex items-center justify-center gap-4">
          {/* Mute button */}
          <Button
            size="lg"
            variant={isMuted ? "destructive" : "secondary"}
            className="w-14 h-14 rounded-full"
            onClick={handleToggleAudio}
            disabled={callState !== 'connected' && callState !== 'calling'}
            aria-label={isMuted ? 'Unmute microphone' : 'Mute microphone'}
            aria-pressed={isMuted}
            data-testid="mute-audio-button"
          >
            {isMuted ? <MicOff className="w-6 h-6" /> : <Mic className="w-6 h-6" />}
          </Button>

          {/* Video toggle */}
          <Button
            size="lg"
            variant={!isVideoOn ? "destructive" : "secondary"}
            className="w-14 h-14 rounded-full"
            onClick={handleToggleVideo}
            disabled={callState !== 'connected' && callState !== 'calling'}
            aria-label={isVideoOn ? 'Turn off camera' : 'Turn on camera'}
            aria-pressed={!isVideoOn}
            data-testid="mute-video-button"
          >
            {isVideoOn ? <Video className="w-6 h-6" /> : <VideoOff className="w-6 h-6" />}
          </Button>

          {/* Screen share button */}
          <Button
            size="lg"
            variant={isScreenSharing ? "default" : "secondary"}
            className="w-14 h-14 rounded-full"
            onClick={handleToggleScreenShare}
            disabled={callState !== 'connected'}
            aria-label={isScreenSharing ? 'Stop screen sharing' : 'Start screen sharing'}
            aria-pressed={isScreenSharing}
            data-testid="screen-share-button"
          >
            <MonitorSpeaker className="w-6 h-6" />
          </Button>

          {/* End call */}
          <Button
            size="lg"
            variant="destructive"
            className="w-16 h-16 rounded-full bg-red-500 hover:bg-red-600"
            onClick={handleEndCall}
            disabled={callState === 'ended'}
            aria-label="End call"
            data-testid="end-call-button"
          >
            <PhoneOff className="w-8 h-8" />
          </Button>

          {/* Speaker toggle */}
          <Button
            size="lg"
            variant={isSpeakerOn ? "default" : "secondary"}
            className="w-14 h-14 rounded-full"
            onClick={() => setIsSpeakerOn(!isSpeakerOn)}
            aria-label={isSpeakerOn ? 'Turn off speaker' : 'Turn on speaker'}
            aria-pressed={isSpeakerOn}
          >
            <Speaker className="w-6 h-6" />
          </Button>

          {/* Chat */}
          <Button
            size="lg"
            variant="secondary"
            className="w-14 h-14 rounded-full"
            aria-label="Open chat"
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