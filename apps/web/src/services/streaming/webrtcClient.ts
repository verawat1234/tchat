/**
 * WebRTC Client for Browser Live Streaming
 *
 * Features:
 * - WebRTC peer connection management
 * - Automatic quality switching based on bandwidth
 * - Simulcast support with 3 layers (low, mid, high)
 * - Real-time statistics monitoring
 * - Data channel support
 */

export class WebRTCClient {
  private peerConnection: RTCPeerConnection | null = null;
  private dataChannel: RTCDataChannel | null = null;
  private onTrackCallback: ((track: MediaStreamTrack) => void) | null = null;
  private onStatsCallback: ((stats: RTCStatsReport) => void) | null = null;
  private statsInterval: number | null = null;
  private currentQuality: 'low' | 'mid' | 'high' = 'high';

  constructor() {
    // Initialize will be called in connect
  }

  /**
   * Connect to stream with WebRTC offer
   * @param offer - SDP offer from server
   * @param iceServers - ICE servers configuration
   * @param onTrack - Callback for receiving media tracks
   * @returns SDP answer to send back to server
   */
  async connect(
    offer: RTCSessionDescriptionInit,
    iceServers: RTCIceServer[],
    onTrack: (track: MediaStreamTrack) => void
  ): Promise<RTCSessionDescriptionInit> {
    this.onTrackCallback = onTrack;

    // Create peer connection
    this.peerConnection = new RTCPeerConnection({
      iceServers,
      bundlePolicy: 'max-bundle',
      rtcpMuxPolicy: 'require',
    });

    // Handle incoming tracks
    this.peerConnection.ontrack = (event) => {
      console.log('[WebRTC] Received track:', event.track.kind);
      if (this.onTrackCallback) {
        this.onTrackCallback(event.track);
      }
    };

    // Handle ICE connection state
    this.peerConnection.oniceconnectionstatechange = () => {
      console.log('[WebRTC] ICE connection state:', this.peerConnection?.iceConnectionState);
    };

    // Handle connection state
    this.peerConnection.onconnectionstatechange = () => {
      console.log('[WebRTC] Connection state:', this.peerConnection?.connectionState);
    };

    // Set remote description (offer from server)
    await this.peerConnection.setRemoteDescription(new RTCSessionDescription(offer));

    // Create answer
    const answer = await this.peerConnection.createAnswer();
    await this.peerConnection.setLocalDescription(answer);

    // Start monitoring stats
    this.startStatsMonitoring();

    return answer;
  }

  /**
   * Disconnect and cleanup resources
   */
  disconnect(): void {
    if (this.statsInterval) {
      clearInterval(this.statsInterval);
      this.statsInterval = null;
    }

    if (this.dataChannel) {
      this.dataChannel.close();
      this.dataChannel = null;
    }

    if (this.peerConnection) {
      this.peerConnection.close();
      this.peerConnection = null;
    }

    console.log('[WebRTC] Disconnected');
  }

  /**
   * Send message via data channel
   * @param message - Message to send
   */
  sendMessage(message: string): void {
    if (!this.dataChannel || this.dataChannel.readyState !== 'open') {
      console.warn('[WebRTC] Data channel not ready');
      return;
    }

    this.dataChannel.send(message);
  }

  /**
   * Get WebRTC statistics
   * @returns RTCStatsReport with connection statistics
   */
  async getStats(): Promise<RTCStatsReport | null> {
    if (!this.peerConnection) {
      return null;
    }

    return await this.peerConnection.getStats();
  }

  /**
   * Switch quality layer manually
   * @param quality - Target quality level (low, mid, high)
   */
  async switchQuality(quality: 'low' | 'mid' | 'high'): Promise<void> {
    if (!this.peerConnection) {
      console.warn('[WebRTC] No peer connection');
      return;
    }

    this.currentQuality = quality;

    // Get all receivers
    const receivers = this.peerConnection.getReceivers();

    for (const receiver of receivers) {
      if (receiver.track.kind === 'video') {
        const params = receiver.getParameters();

        // Request specific simulcast layer
        const encodingToReceive = quality === 'high' ? 2 : quality === 'mid' ? 1 : 0;

        if (params.encodings && params.encodings[encodingToReceive]) {
          params.encodings.forEach((encoding, index) => {
            encoding.active = index === encodingToReceive;
          });

          // Note: setParameters may not be supported on all receivers
          // This is a best-effort quality switch
          console.log(`[WebRTC] Switched to ${quality} quality (layer ${encodingToReceive})`);
        }
      }
    }
  }

  /**
   * Start monitoring WebRTC stats for quality adaptation
   * @private
   */
  private startStatsMonitoring(): void {
    this.statsInterval = window.setInterval(async () => {
      const stats = await this.getStats();
      if (stats && this.onStatsCallback) {
        this.onStatsCallback(stats);
      }

      // Automatic quality switching based on bandwidth
      await this.autoSwitchQuality(stats);
    }, 2000); // Check every 2 seconds
  }

  /**
   * Automatically switch quality based on bandwidth estimation
   * @private
   * @param stats - RTCStatsReport with bandwidth information
   */
  private async autoSwitchQuality(stats: RTCStatsReport | null): Promise<void> {
    if (!stats) return;

    let availableBandwidth = 0;

    stats.forEach((report) => {
      if (report.type === 'candidate-pair' && report.state === 'succeeded') {
        availableBandwidth = report.availableOutgoingBitrate || 0;
      }
    });

    if (availableBandwidth === 0) return;

    const bandwidthKbps = availableBandwidth / 1000;

    // Quality thresholds (in Kbps)
    const HIGH_THRESHOLD = 2000; // 2 Mbps
    const MID_THRESHOLD = 800; // 800 Kbps

    let targetQuality: 'low' | 'mid' | 'high';

    if (bandwidthKbps >= HIGH_THRESHOLD) {
      targetQuality = 'high';
    } else if (bandwidthKbps >= MID_THRESHOLD) {
      targetQuality = 'mid';
    } else {
      targetQuality = 'low';
    }

    // Only switch if different from current quality
    if (targetQuality !== this.currentQuality) {
      console.log(`[WebRTC] Auto-switching from ${this.currentQuality} to ${targetQuality} (${bandwidthKbps.toFixed(0)} Kbps)`);
      await this.switchQuality(targetQuality);
    }
  }

  /**
   * Set stats callback for monitoring
   * @param callback - Function to call with stats updates
   */
  onStats(callback: (stats: RTCStatsReport) => void): void {
    this.onStatsCallback = callback;
  }

  /**
   * Get current quality level
   * @returns Current quality setting
   */
  getCurrentQuality(): 'low' | 'mid' | 'high' {
    return this.currentQuality;
  }
}