/**
 * CallService - WebRTC Call Management Service
 *
 * Comprehensive WebRTC service for managing voice and video calls
 * Handles peer connections, media streams, and call state management
 * Integrates with SignalingClient for WebRTC coordination
 */

import { SignalingClient } from './SignalingClient';

export type CallType = 'voice' | 'video';
export type CallState = 'idle' | 'calling' | 'ringing' | 'connected' | 'ended' | 'error';
export type MediaDevice = 'camera' | 'microphone' | 'speaker';

export interface CallParticipant {
  id: string;
  name: string;
  avatar?: string;
  isLocal: boolean;
  audioEnabled: boolean;
  videoEnabled: boolean;
  isScreenSharing: boolean;
}

export interface CallSession {
  id: string;
  type: CallType;
  state: CallState;
  participants: CallParticipant[];
  startTime?: Date;
  endTime?: Date;
  duration?: number;
  quality: {
    video: 'low' | 'medium' | 'high' | 'auto';
    audio: 'low' | 'medium' | 'high' | 'auto';
  };
  networkStats: {
    latency: number;
    bandwidth: number;
    packetLoss: number;
  };
}

export interface MediaConstraints {
  video: boolean | MediaTrackConstraints;
  audio: boolean | MediaTrackConstraints;
}

export interface CallOptions {
  type: CallType;
  enableVideo: boolean;
  enableAudio: boolean;
  quality: 'low' | 'medium' | 'high' | 'auto';
}

export interface CallEventHandlers {
  onStateChange?: (state: CallState) => void;
  onParticipantJoined?: (participant: CallParticipant) => void;
  onParticipantLeft?: (participantId: string) => void;
  onMediaToggle?: (participantId: string, type: 'audio' | 'video', enabled: boolean) => void;
  onScreenShareStart?: (participantId: string) => void;
  onScreenShareEnd?: (participantId: string) => void;
  onNetworkQualityChange?: (quality: 'poor' | 'fair' | 'good' | 'excellent') => void;
  onError?: (error: Error) => void;
}

export class CallService {
  private peerConnection: RTCPeerConnection | null = null;
  private localStream: MediaStream | null = null;
  private remoteStream: MediaStream | null = null;
  private screenShareStream: MediaStream | null = null;
  private signalingClient: SignalingClient;
  private currentCall: CallSession | null = null;
  private eventHandlers: CallEventHandlers = {};
  private isInitialized = false;

  // WebRTC Configuration
  private readonly rtcConfig: RTCConfiguration = {
    iceServers: [
      { urls: 'stun:stun.l.google.com:19302' },
      { urls: 'stun:stun1.l.google.com:19302' },
      { urls: 'stun:stun2.l.google.com:19302' }
    ],
    iceCandidatePoolSize: 10
  };

  constructor() {
    this.signalingClient = new SignalingClient();
    this.setupSignalingHandlers();
  }

  /**
   * Initialize the call service
   */
  async initialize(): Promise<void> {
    if (this.isInitialized) return;

    try {
      await this.signalingClient.connect();
      this.isInitialized = true;
    } catch (error) {
      console.error('Failed to initialize CallService:', error);
      throw error;
    }
  }

  /**
   * Initiate a new call
   */
  async initiateCall(
    recipientId: string,
    recipientName: string,
    options: CallOptions
  ): Promise<string> {
    if (!this.isInitialized) {
      await this.initialize();
    }

    if (this.currentCall && this.currentCall.state !== 'ended') {
      throw new Error('Another call is already in progress');
    }

    const callId = this.generateCallId();

    this.currentCall = {
      id: callId,
      type: options.type,
      state: 'calling',
      participants: [
        {
          id: 'local',
          name: 'You',
          isLocal: true,
          audioEnabled: options.enableAudio,
          videoEnabled: options.enableVideo,
          isScreenSharing: false
        },
        {
          id: recipientId,
          name: recipientName,
          isLocal: false,
          audioEnabled: true,
          videoEnabled: options.type === 'video',
          isScreenSharing: false
        }
      ],
      startTime: new Date(),
      quality: {
        video: options.quality,
        audio: options.quality
      },
      networkStats: {
        latency: 0,
        bandwidth: 0,
        packetLoss: 0
      }
    };

    try {
      // Setup peer connection
      await this.setupPeerConnection();

      // Get user media
      await this.getUserMedia({
        video: options.enableVideo,
        audio: options.enableAudio
      });

      // Create offer
      const offer = await this.peerConnection!.createOffer();
      await this.peerConnection!.setLocalDescription(offer);

      // Send call initiation through signaling
      await this.signalingClient.sendMessage({
        type: 'call-initiate',
        callId,
        recipientId,
        callType: options.type,
        offer: offer.sdp,
        hasVideo: options.enableVideo
      });

      this.notifyStateChange('calling');

      return callId;
    } catch (error) {
      await this.endCall();
      throw error;
    }
  }

  /**
   * Answer an incoming call
   */
  async answerCall(callId: string, options: CallOptions): Promise<void> {
    if (!this.isInitialized) {
      await this.initialize();
    }

    try {
      // Setup peer connection
      await this.setupPeerConnection();

      // Get user media
      await this.getUserMedia({
        video: options.enableVideo,
        audio: options.enableAudio
      });

      // Update call state
      if (this.currentCall) {
        this.currentCall.state = 'connected';
        this.notifyStateChange('connected');
      }

      // Send answer through signaling
      await this.signalingClient.sendMessage({
        type: 'call-answer',
        callId,
        hasVideo: options.enableVideo
      });

    } catch (error) {
      await this.endCall();
      throw error;
    }
  }

  /**
   * Decline an incoming call
   */
  async declineCall(callId: string): Promise<void> {
    await this.signalingClient.sendMessage({
      type: 'call-decline',
      callId
    });
  }

  /**
   * End the current call
   */
  async endCall(): Promise<void> {
    if (this.currentCall) {
      this.currentCall.state = 'ended';
      this.currentCall.endTime = new Date();
      if (this.currentCall.startTime) {
        this.currentCall.duration = this.currentCall.endTime.getTime() - this.currentCall.startTime.getTime();
      }

      // Send end call signal
      await this.signalingClient.sendMessage({
        type: 'call-end',
        callId: this.currentCall.id
      });

      this.notifyStateChange('ended');
    }

    await this.cleanup();
  }

  /**
   * Toggle audio mute
   */
  toggleAudio(): boolean {
    if (!this.localStream) return false;

    const audioTrack = this.localStream.getAudioTracks()[0];
    if (audioTrack) {
      audioTrack.enabled = !audioTrack.enabled;

      if (this.currentCall) {
        const localParticipant = this.currentCall.participants.find(p => p.isLocal);
        if (localParticipant) {
          localParticipant.audioEnabled = audioTrack.enabled;
          this.eventHandlers.onMediaToggle?.('local', 'audio', audioTrack.enabled);
        }
      }

      // Notify remote peer
      this.signalingClient.sendMessage({
        type: 'media-toggle',
        callId: this.currentCall?.id,
        mediaType: 'audio',
        enabled: audioTrack.enabled
      });

      return audioTrack.enabled;
    }
    return false;
  }

  /**
   * Toggle video mute
   */
  toggleVideo(): boolean {
    if (!this.localStream) return false;

    const videoTrack = this.localStream.getVideoTracks()[0];
    if (videoTrack) {
      videoTrack.enabled = !videoTrack.enabled;

      if (this.currentCall) {
        const localParticipant = this.currentCall.participants.find(p => p.isLocal);
        if (localParticipant) {
          localParticipant.videoEnabled = videoTrack.enabled;
          this.eventHandlers.onMediaToggle?.('local', 'video', videoTrack.enabled);
        }
      }

      // Notify remote peer
      this.signalingClient.sendMessage({
        type: 'media-toggle',
        callId: this.currentCall?.id,
        mediaType: 'video',
        enabled: videoTrack.enabled
      });

      return videoTrack.enabled;
    }
    return false;
  }

  /**
   * Switch camera (front/back)
   */
  async switchCamera(): Promise<void> {
    if (!this.localStream || this.currentCall?.type !== 'video') return;

    const videoTrack = this.localStream.getVideoTracks()[0];
    if (!videoTrack) return;

    try {
      const devices = await navigator.mediaDevices.enumerateDevices();
      const videoDevices = devices.filter(device => device.kind === 'videoinput');

      if (videoDevices.length < 2) return; // No camera switching available

      // Get current device ID
      const currentDeviceId = videoTrack.getSettings().deviceId;
      const nextDevice = videoDevices.find(device => device.deviceId !== currentDeviceId) || videoDevices[0];

      // Replace video track
      const newStream = await navigator.mediaDevices.getUserMedia({
        video: { deviceId: nextDevice.deviceId },
        audio: false
      });

      const newVideoTrack = newStream.getVideoTracks()[0];
      const sender = this.peerConnection?.getSenders().find(s => s.track === videoTrack);

      if (sender && newVideoTrack) {
        await sender.replaceTrack(newVideoTrack);

        // Replace track in local stream
        this.localStream.removeTrack(videoTrack);
        this.localStream.addTrack(newVideoTrack);

        videoTrack.stop();
      }
    } catch (error) {
      console.error('Failed to switch camera:', error);
      this.eventHandlers.onError?.(error as Error);
    }
  }

  /**
   * Start screen sharing
   */
  async startScreenShare(): Promise<void> {
    if (this.screenShareStream) return; // Already sharing

    try {
      const screenStream = await navigator.mediaDevices.getDisplayMedia({
        video: true,
        audio: true
      });

      this.screenShareStream = screenStream;

      // Replace video track with screen share
      const videoTrack = screenStream.getVideoTracks()[0];
      const sender = this.peerConnection?.getSenders().find(
        s => s.track && s.track.kind === 'video'
      );

      if (sender && videoTrack) {
        await sender.replaceTrack(videoTrack);
      }

      // Handle screen share end
      videoTrack.onended = () => {
        this.stopScreenShare();
      };

      if (this.currentCall) {
        const localParticipant = this.currentCall.participants.find(p => p.isLocal);
        if (localParticipant) {
          localParticipant.isScreenSharing = true;
          this.eventHandlers.onScreenShareStart?.('local');
        }
      }

      // Notify remote peer
      this.signalingClient.sendMessage({
        type: 'screen-share-start',
        callId: this.currentCall?.id
      });

    } catch (error) {
      console.error('Failed to start screen sharing:', error);
      this.eventHandlers.onError?.(error as Error);
    }
  }

  /**
   * Stop screen sharing
   */
  async stopScreenShare(): Promise<void> {
    if (!this.screenShareStream) return;

    try {
      // Stop screen share tracks
      this.screenShareStream.getTracks().forEach(track => track.stop());
      this.screenShareStream = null;

      // Restore camera video
      if (this.localStream && this.currentCall?.type === 'video') {
        const cameraVideoTrack = this.localStream.getVideoTracks()[0];
        const sender = this.peerConnection?.getSenders().find(
          s => s.track && s.track.kind === 'video'
        );

        if (sender && cameraVideoTrack) {
          await sender.replaceTrack(cameraVideoTrack);
        }
      }

      if (this.currentCall) {
        const localParticipant = this.currentCall.participants.find(p => p.isLocal);
        if (localParticipant) {
          localParticipant.isScreenSharing = false;
          this.eventHandlers.onScreenShareEnd?.('local');
        }
      }

      // Notify remote peer
      this.signalingClient.sendMessage({
        type: 'screen-share-end',
        callId: this.currentCall?.id
      });

    } catch (error) {
      console.error('Failed to stop screen sharing:', error);
      this.eventHandlers.onError?.(error as Error);
    }
  }

  /**
   * Get available media devices
   */
  async getAvailableDevices(): Promise<{
    cameras: MediaDeviceInfo[];
    microphones: MediaDeviceInfo[];
    speakers: MediaDeviceInfo[];
  }> {
    const devices = await navigator.mediaDevices.enumerateDevices();

    return {
      cameras: devices.filter(device => device.kind === 'videoinput'),
      microphones: devices.filter(device => device.kind === 'audioinput'),
      speakers: devices.filter(device => device.kind === 'audiooutput')
    };
  }

  /**
   * Set event handlers
   */
  setEventHandlers(handlers: CallEventHandlers): void {
    this.eventHandlers = { ...this.eventHandlers, ...handlers };
  }

  /**
   * Get current call session
   */
  getCurrentCall(): CallSession | null {
    return this.currentCall;
  }

  /**
   * Get local media stream
   */
  getLocalStream(): MediaStream | null {
    return this.localStream;
  }

  /**
   * Get remote media stream
   */
  getRemoteStream(): MediaStream | null {
    return this.remoteStream;
  }

  /**
   * Get call statistics
   */
  async getCallStats(): Promise<RTCStatsReport | null> {
    if (!this.peerConnection) return null;
    return await this.peerConnection.getStats();
  }

  // Private methods

  private async setupPeerConnection(): Promise<void> {
    if (this.peerConnection) {
      this.peerConnection.close();
    }

    this.peerConnection = new RTCPeerConnection(this.rtcConfig);

    // Handle ICE candidates
    this.peerConnection.onicecandidate = (event) => {
      if (event.candidate) {
        this.signalingClient.sendMessage({
          type: 'ice-candidate',
          callId: this.currentCall?.id,
          candidate: event.candidate
        });
      }
    };

    // Handle remote stream
    this.peerConnection.ontrack = (event) => {
      this.remoteStream = event.streams[0];
    };

    // Handle connection state changes
    this.peerConnection.onconnectionstatechange = () => {
      const state = this.peerConnection?.connectionState;
      console.log('Peer connection state:', state);

      if (state === 'connected' && this.currentCall) {
        this.currentCall.state = 'connected';
        this.notifyStateChange('connected');
      } else if (state === 'failed' || state === 'disconnected') {
        this.handleConnectionFailure();
      }
    };

    // Monitor network quality
    this.startNetworkQualityMonitoring();
  }

  private async getUserMedia(constraints: MediaConstraints): Promise<void> {
    try {
      this.localStream = await navigator.mediaDevices.getUserMedia(constraints);

      // Add tracks to peer connection
      this.localStream.getTracks().forEach(track => {
        this.peerConnection?.addTrack(track, this.localStream!);
      });

    } catch (error) {
      console.error('Failed to get user media:', error);
      throw new Error('Failed to access camera/microphone');
    }
  }

  private setupSignalingHandlers(): void {
    this.signalingClient.setEventHandlers({
      onMessage: this.handleSignalingMessage.bind(this),
      onError: (error) => this.eventHandlers.onError?.(error),
      onConnectionChange: (connected) => {
        if (!connected && this.currentCall?.state === 'connected') {
          this.handleConnectionFailure();
        }
      }
    });
  }

  private async handleSignalingMessage(message: any): Promise<void> {
    try {
      switch (message.type) {
        case 'call-offer':
          await this.handleCallOffer(message);
          break;
        case 'call-answer':
          await this.handleCallAnswer(message);
          break;
        case 'ice-candidate':
          await this.handleIceCandidate(message);
          break;
        case 'call-end':
          await this.endCall();
          break;
        case 'media-toggle':
          this.handleMediaToggle(message);
          break;
        case 'screen-share-start':
          this.handleScreenShareStart(message);
          break;
        case 'screen-share-end':
          this.handleScreenShareEnd(message);
          break;
      }
    } catch (error) {
      console.error('Error handling signaling message:', error);
      this.eventHandlers.onError?.(error as Error);
    }
  }

  private async handleCallOffer(message: any): Promise<void> {
    // This would be called when receiving an incoming call
    // Implementation depends on UI integration
  }

  private async handleCallAnswer(message: any): Promise<void> {
    if (message.answer && this.peerConnection) {
      await this.peerConnection.setRemoteDescription({
        type: 'answer',
        sdp: message.answer
      });
    }
  }

  private async handleIceCandidate(message: any): Promise<void> {
    if (message.candidate && this.peerConnection) {
      await this.peerConnection.addIceCandidate(message.candidate);
    }
  }

  private handleMediaToggle(message: any): void {
    if (this.currentCall) {
      const participant = this.currentCall.participants.find(p => !p.isLocal);
      if (participant) {
        if (message.mediaType === 'audio') {
          participant.audioEnabled = message.enabled;
        } else if (message.mediaType === 'video') {
          participant.videoEnabled = message.enabled;
        }
        this.eventHandlers.onMediaToggle?.(participant.id, message.mediaType, message.enabled);
      }
    }
  }

  private handleScreenShareStart(message: any): void {
    if (this.currentCall) {
      const participant = this.currentCall.participants.find(p => !p.isLocal);
      if (participant) {
        participant.isScreenSharing = true;
        this.eventHandlers.onScreenShareStart?.(participant.id);
      }
    }
  }

  private handleScreenShareEnd(message: any): void {
    if (this.currentCall) {
      const participant = this.currentCall.participants.find(p => !p.isLocal);
      if (participant) {
        participant.isScreenSharing = false;
        this.eventHandlers.onScreenShareEnd?.(participant.id);
      }
    }
  }

  private handleConnectionFailure(): void {
    if (this.currentCall) {
      this.currentCall.state = 'error';
      this.notifyStateChange('error');
      this.eventHandlers.onError?.(new Error('Connection failed'));
    }
  }

  private startNetworkQualityMonitoring(): void {
    // Monitor connection quality every 5 seconds
    const interval = setInterval(async () => {
      if (!this.peerConnection || this.currentCall?.state !== 'connected') {
        clearInterval(interval);
        return;
      }

      try {
        const stats = await this.peerConnection.getStats();
        const networkQuality = this.analyzeNetworkQuality(stats);

        if (this.currentCall) {
          this.currentCall.networkStats = networkQuality;

          let quality: 'poor' | 'fair' | 'good' | 'excellent' = 'good';
          if (networkQuality.packetLoss > 5) quality = 'poor';
          else if (networkQuality.latency > 300) quality = 'fair';
          else if (networkQuality.latency < 100 && networkQuality.packetLoss < 1) quality = 'excellent';

          this.eventHandlers.onNetworkQualityChange?.(quality);
        }
      } catch (error) {
        console.error('Failed to get network stats:', error);
      }
    }, 5000);
  }

  private analyzeNetworkQuality(stats: RTCStatsReport): { latency: number; bandwidth: number; packetLoss: number } {
    let latency = 0;
    let bandwidth = 0;
    let packetLoss = 0;

    stats.forEach((report) => {
      if (report.type === 'candidate-pair' && report.state === 'succeeded') {
        latency = report.currentRoundTripTime * 1000 || 0;
      }
      if (report.type === 'outbound-rtp') {
        bandwidth += report.bytesSent || 0;
      }
      if (report.type === 'inbound-rtp') {
        const packetsReceived = report.packetsReceived || 0;
        const packetsLost = report.packetsLost || 0;
        if (packetsReceived > 0) {
          packetLoss = (packetsLost / (packetsReceived + packetsLost)) * 100;
        }
      }
    });

    return { latency, bandwidth, packetLoss };
  }

  private notifyStateChange(state: CallState): void {
    if (this.currentCall) {
      this.currentCall.state = state;
    }
    this.eventHandlers.onStateChange?.(state);
  }

  private async cleanup(): Promise<void> {
    // Stop all tracks
    this.localStream?.getTracks().forEach(track => track.stop());
    this.screenShareStream?.getTracks().forEach(track => track.stop());

    // Close peer connection
    this.peerConnection?.close();

    // Clear streams
    this.localStream = null;
    this.remoteStream = null;
    this.screenShareStream = null;
    this.peerConnection = null;
  }

  private generateCallId(): string {
    return `call_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
}

// Singleton instance
export const callService = new CallService();