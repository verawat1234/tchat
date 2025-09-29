import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
  Phone,
  PhoneOff,
  Mic,
  MicOff,
  Speaker,
  MessageSquare,
  Users,
  Volume2,
  VolumeX,
  ArrowLeft,
  UserPlus
} from 'lucide-react';
import { Button } from './ui/button';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { callService, CallState, CallSession } from '../services/webrtc/CallService';

interface VoiceCallScreenProps {
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

export function VoiceCallScreen({
  user,
  callee,
  callId,
  isIncoming = false,
  onEndCall,
  onBack,
  onCallAccept,
  onCallDecline
}: VoiceCallScreenProps) {
  // WebRTC state
  const [callState, setCallState] = useState<CallState>('idle');
  const [currentCall, setCurrentCall] = useState<CallSession | null>(null);
  const [isAudioMuted, setIsAudioMuted] = useState(false);
  const [networkQuality, setNetworkQuality] = useState<'poor' | 'fair' | 'good' | 'excellent'>('good');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Legacy UI state for compatibility
  const [isConnected, setIsConnected] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [isSpeakerOn, setIsSpeakerOn] = useState(false);
  const [callDuration, setCallDuration] = useState(0);
  const [audioLevel, setAudioLevel] = useState(0);

  const audioRef = useRef<HTMLAudioElement>(null);
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
            if (participantId === 'local' && type === 'audio') {
              setIsAudioMuted(!enabled);
              setIsMuted(!enabled);
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

  // Setup audio output when call connects
  useEffect(() => {
    if (callState === 'connected') {
      const remoteStream = callService.getRemoteStream();
      if (remoteStream && audioRef.current) {
        audioRef.current.srcObject = remoteStream;
        audioRef.current.play().catch(console.error);
      }
    }
  }, [callState]);

  // Handle incoming call setup
  useEffect(() => {
    if (isIncoming && callId && !currentCall) {
      setCallState('ringing');
    }
  }, [isIncoming, callId, currentCall]);

  // Legacy audio level simulation for UI animations
  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (isConnected) {
      interval = setInterval(() => {
        setAudioLevel(Math.random() * 100);
      }, 500);
    }
    return () => clearInterval(interval);
  }, [isConnected]);

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
        type: 'voice',
        enableVideo: false,
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

  if (isIncoming && (callState === 'ringing' || !isConnected)) {
    return (
      <div
        className="h-screen bg-gradient-to-br from-chart-1/20 to-chart-3/20 flex flex-col items-center justify-center relative"
        data-testid="incoming-call-modal"
        role="dialog"
        aria-modal="true"
        aria-labelledby="incoming-call-title"
      >
        {/* Hidden audio element for remote stream */}
        <audio
          ref={audioRef}
          autoPlay
          playsInline
          className="hidden"
        />

        {/* Background pattern */}
        <div className="absolute inset-0 bg-black/10 backdrop-blur-sm"></div>
        
        <div className="relative z-10 text-center space-y-8">
          {/* Caller info */}
          <div className="space-y-6">
            <div className="relative">
              <Avatar className="w-40 h-40 mx-auto ring-4 ring-white/30 animate-pulse">
                <AvatarImage src={callee.avatar} />
                <AvatarFallback className="text-3xl">
                  {callee.isGroup ? <Users className="w-16 h-16" /> : callee.name.charAt(0)}
                </AvatarFallback>
              </Avatar>
              
              {/* Audio wave animation */}
              <div className="absolute -bottom-2 left-1/2 transform -translate-x-1/2 flex items-end gap-1">
                {[...Array(5)].map((_, i) => (
                  <div
                    key={i}
                    className="w-1 bg-chart-1 rounded-full animate-pulse"
                    style={{
                      height: `${Math.random() * 20 + 10}px`,
                      animationDelay: `${i * 0.1}s`
                    }}
                  />
                ))}
              </div>
            </div>
            
            <div className="space-y-2">
              <h2
                id="incoming-call-title"
                className="text-3xl"
                data-testid="caller-name"
              >
                {callee.name}
              </h2>
              {callee.isGroup && (
                <Badge variant="secondary" className="bg-white/20">
                  {callee.members} members
                </Badge>
              )}
              <p
                className="text-muted-foreground text-lg"
                data-testid="call-type"
              >
                Voice Call
              </p>
            </div>
          </div>

          {/* Error message */}
          {error && (
            <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
              {error}
            </div>
          )}

          {/* Call actions */}
          <div className="flex items-center justify-center gap-12">
            <Button
              size="lg"
              variant="destructive"
              className="w-20 h-20 rounded-full bg-red-500 hover:bg-red-600 shadow-lg"
              onClick={handleDecline}
              disabled={isLoading}
              aria-label="Decline call"
              data-testid="decline-call-button"
            >
              <PhoneOff className="w-10 h-10" />
            </Button>

            <Button
              size="lg"
              className="w-20 h-20 rounded-full bg-green-500 hover:bg-green-600 shadow-lg"
              onClick={handleAnswer}
              disabled={isLoading}
              aria-label="Accept call"
              data-testid="accept-voice-call-button"
            >
              <Phone className="w-10 h-10" />
            </Button>
          </div>

          {/* Loading indicator */}
          {isLoading && (
            <div className="mt-4 text-gray-600">
              Connecting...
            </div>
          )}

          {/* Quick message */}
          <div className="space-y-3">
            <p className="text-sm text-muted-foreground">Quick reply</p>
            <div className="flex flex-wrap justify-center gap-2">
              <Button variant="outline" size="sm" className="rounded-full">
                Can't talk now
              </Button>
              <Button variant="outline" size="sm" className="rounded-full">
                Call you back
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div
      className="h-screen bg-gradient-to-br from-background to-muted flex flex-col"
      data-testid="voice-call-screen"
    >
      {/* Hidden audio element for remote stream */}
      <audio
        ref={audioRef}
        autoPlay
        playsInline
        className="hidden"
      />

      {/* Network warning */}
      {networkQuality === 'poor' && (
        <div
          className="absolute top-4 left-4 bg-red-500 bg-opacity-90 px-3 py-2 rounded-lg text-sm text-white z-10"
          data-testid="network-warning"
        >
          Poor network connection
        </div>
      )}

      {/* Error message */}
      {error && (
        <div
          className="absolute top-4 right-4 bg-red-500 bg-opacity-90 px-3 py-2 rounded-lg text-sm text-white z-10 max-w-xs"
          data-testid="call-error-message"
        >
          {error}
        </div>
      )}

      {/* Call timeout message */}
      {callState === 'error' && error?.includes('timeout') && (
        <div
          className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-yellow-500 bg-opacity-90 px-4 py-3 rounded-lg text-center z-20"
          data-testid="call-timeout-message"
        >
          Call timed out
        </div>
      )}

      {/* Header */}
      <div className="p-4 flex items-center justify-between border-b">
        <Button
          variant="ghost"
          size="icon"
          onClick={onBack}
        >
          <ArrowLeft className="w-5 h-5" />
        </Button>
        
        <div className="text-center">
          <p className="text-sm text-muted-foreground">Voice call</p>
          <p
            className="text-sm"
            data-testid="call-status"
          >
            {callState === 'connected' ? formatDuration(callDuration) :
             callState === 'calling' ? 'Calling...' :
             callState === 'ringing' ? 'Ringing...' : 'Connecting...'}
          </p>
          {callState === 'connected' && (
            <p
              className="text-xs text-muted-foreground mt-1"
              data-testid="call-duration"
            >
              Duration: {formatDuration(callDuration)}
            </p>
          )}
        </div>

        <Button
          variant="ghost"
          size="icon"
          disabled={!isConnected}
        >
          <UserPlus className="w-5 h-5" />
        </Button>
      </div>

      {/* Main call area */}
      <div className="flex-1 flex flex-col items-center justify-center space-y-8 p-8">
        {/* Caller avatar with audio visualization */}
        <div className="relative">
          <Avatar className="w-48 h-48 ring-4 ring-primary/20">
            <AvatarImage src={callee.avatar} />
            <AvatarFallback className="text-4xl">
              {callee.isGroup ? <Users className="w-20 h-20" /> : callee.name.charAt(0)}
            </AvatarFallback>
          </Avatar>
          
          {/* Audio level indicator */}
          {isConnected && (
            <div className="absolute inset-0 rounded-full border-4 border-chart-1 animate-ping opacity-75" 
                 style={{ animationDuration: '2s' }} />
          )}
          
          {/* Speaking indicator */}
          {isConnected && audioLevel > 50 && (
            <div className="absolute -bottom-4 left-1/2 transform -translate-x-1/2 flex items-center gap-1">
              <Volume2 className="w-4 h-4 text-chart-1" />
              <div className="flex gap-1">
                {[...Array(3)].map((_, i) => (
                  <div
                    key={i}
                    className="w-1 h-4 bg-chart-1 rounded-full animate-pulse"
                    style={{ animationDelay: `${i * 0.2}s` }}
                  />
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Caller info */}
        <div className="text-center space-y-2">
          <h2 className="text-2xl">{callee.name}</h2>
          {callee.isGroup && (
            <Badge variant="secondary">
              {callee.members} members
            </Badge>
          )}
          <p className="text-muted-foreground">
            {isConnected ? 'Connected' : 'Connecting...'}
          </p>
        </div>

        {/* Connection quality indicator */}
        {isConnected && (
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <div className="flex gap-1">
              {[...Array(4)].map((_, i) => (
                <div
                  key={i}
                  className={`w-1 h-3 rounded-full ${
                    i < 3 ? 'bg-green-500' : 'bg-muted'
                  }`}
                />
              ))}
            </div>
            <span>Good connection</span>
          </div>
        )}
      </div>

      {/* Call controls */}
      <div
        className="p-6 bg-card border-t"
        data-testid="voice-controls"
      >
        <div className="flex items-center justify-center gap-6">
          {/* Mute button */}
          <Button
            size="lg"
            variant={isMuted ? "destructive" : "secondary"}
            className="w-16 h-16 rounded-full"
            onClick={handleToggleAudio}
            disabled={callState !== 'connected' && callState !== 'calling'}
            aria-label={isMuted ? 'Unmute microphone' : 'Mute microphone'}
            aria-pressed={isMuted}
            data-testid="mute-audio-button"
          >
            {isMuted ? <MicOff className="w-7 h-7" /> : <Mic className="w-7 h-7" />}
          </Button>

          {/* End call */}
          <Button
            size="lg"
            variant="destructive"
            className="w-20 h-20 rounded-full bg-red-500 hover:bg-red-600"
            onClick={handleEndCall}
            disabled={callState === 'ended'}
            aria-label="End call"
            data-testid="end-call-button"
          >
            <PhoneOff className="w-10 h-10" />
          </Button>

          {/* Speaker toggle */}
          <Button
            size="lg"
            variant={isSpeakerOn ? "default" : "secondary"}
            className="w-16 h-16 rounded-full"
            onClick={() => setIsSpeakerOn(!isSpeakerOn)}
          >
            {isSpeakerOn ? <Volume2 className="w-7 h-7" /> : <VolumeX className="w-7 h-7" />}
          </Button>
        </div>

        {/* Secondary controls */}
        <div className="flex items-center justify-center gap-4 mt-4">
          <Button variant="ghost" size="sm">
            <MessageSquare className="w-4 h-4 mr-2" />
            Message
          </Button>
          
          {callee.isGroup && (
            <Button variant="ghost" size="sm">
              <Users className="w-4 h-4 mr-2" />
              Members
            </Button>
          )}
        </div>

        {/* Mute status */}
        {isMuted && (
          <div className="text-center mt-4">
            <Badge variant="secondary" className="bg-red-500/20 text-red-600">
              <MicOff className="w-3 h-3 mr-1" />
              You're muted
            </Badge>
          </div>
        )}
      </div>
    </div>
  );
}