/**
 * SignalingClient - WebSocket Signaling for WebRTC
 *
 * Handles WebSocket communication for WebRTC signaling
 * Manages connection, message routing, and error handling
 * Supports call initiation, ICE candidate exchange, and media coordination
 */

export interface SignalingMessage {
  type: string;
  callId?: string;
  recipientId?: string;
  callType?: 'voice' | 'video';
  offer?: string;
  answer?: string;
  candidate?: RTCIceCandidate;
  hasVideo?: boolean;
  mediaType?: 'audio' | 'video';
  enabled?: boolean;
  [key: string]: any;
}

export interface SignalingEventHandlers {
  onMessage?: (message: SignalingMessage) => void;
  onError?: (error: Error) => void;
  onConnectionChange?: (connected: boolean) => void;
  onReconnecting?: () => void;
  onReconnected?: () => void;
}

export class SignalingClient {
  private ws: WebSocket | null = null;
  private eventHandlers: SignalingEventHandlers = {};
  private isConnected = false;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000; // Start with 1 second
  private maxReconnectDelay = 30000; // Max 30 seconds
  private heartbeatInterval: NodeJS.Timeout | null = null;
  private heartbeatDelay = 30000; // 30 seconds
  private messageQueue: SignalingMessage[] = [];
  private connectionId: string | null = null;

  private readonly signalingUrl: string;

  constructor() {
    // Use environment variable or default to localhost
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = import.meta.env.VITE_SIGNALING_HOST || 'localhost:8080';
    this.signalingUrl = `${protocol}//${host}/ws/calling`;
  }

  /**
   * Connect to the signaling server
   */
  async connect(): Promise<void> {
    if (this.isConnected) return;

    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(this.signalingUrl);

        const connectTimeout = setTimeout(() => {
          this.ws?.close();
          reject(new Error('Connection timeout'));
        }, 10000); // 10 second timeout

        this.ws.onopen = () => {
          clearTimeout(connectTimeout);
          this.isConnected = true;
          this.reconnectAttempts = 0;
          this.reconnectDelay = 1000;

          console.log('Signaling client connected');
          this.eventHandlers.onConnectionChange?.(true);

          // Start heartbeat
          this.startHeartbeat();

          // Send queued messages
          this.flushMessageQueue();

          resolve();
        };

        this.ws.onmessage = (event) => {
          this.handleMessage(event);
        };

        this.ws.onclose = (event) => {
          clearTimeout(connectTimeout);
          this.handleDisconnection(event);
        };

        this.ws.onerror = (error) => {
          clearTimeout(connectTimeout);
          console.error('Signaling WebSocket error:', error);
          this.eventHandlers.onError?.(new Error('WebSocket connection error'));
          reject(error);
        };

      } catch (error) {
        reject(error);
      }
    });
  }

  /**
   * Disconnect from the signaling server
   */
  disconnect(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }

    if (this.ws) {
      this.ws.close(1000, 'Client disconnect');
      this.ws = null;
    }

    this.isConnected = false;
    this.eventHandlers.onConnectionChange?.(false);
  }

  /**
   * Send a message through the signaling channel
   */
  async sendMessage(message: SignalingMessage): Promise<void> {
    if (!this.isConnected || !this.ws || this.ws.readyState !== WebSocket.OPEN) {
      // Queue message for when connection is restored
      this.messageQueue.push(message);

      // Attempt to reconnect if not currently connected
      if (!this.isConnected) {
        this.attemptReconnect();
      }
      return;
    }

    try {
      const messageWithMetadata = {
        ...message,
        timestamp: Date.now(),
        connectionId: this.connectionId
      };

      this.ws.send(JSON.stringify(messageWithMetadata));
    } catch (error) {
      console.error('Failed to send message:', error);
      this.eventHandlers.onError?.(error as Error);

      // Queue message for retry
      this.messageQueue.push(message);
    }
  }

  /**
   * Set event handlers
   */
  setEventHandlers(handlers: SignalingEventHandlers): void {
    this.eventHandlers = { ...this.eventHandlers, ...handlers };
  }

  /**
   * Get connection status
   */
  isConnectedToServer(): boolean {
    return this.isConnected && this.ws?.readyState === WebSocket.OPEN;
  }

  /**
   * Get connection ID
   */
  getConnectionId(): string | null {
    return this.connectionId;
  }

  // Private methods

  private handleMessage(event: MessageEvent): void {
    try {
      const message: SignalingMessage = JSON.parse(event.data);

      // Handle special control messages
      switch (message.type) {
        case 'connection-established':
          this.connectionId = message.connectionId;
          console.log('Signaling connection established:', this.connectionId);
          break;

        case 'pong':
          // Heartbeat response - connection is alive
          break;

        case 'error':
          console.error('Signaling server error:', message.error);
          this.eventHandlers.onError?.(new Error(message.error));
          break;

        default:
          // Forward message to call service
          this.eventHandlers.onMessage?.(message);
          break;
      }

    } catch (error) {
      console.error('Failed to parse signaling message:', error);
      this.eventHandlers.onError?.(new Error('Invalid signaling message format'));
    }
  }

  private handleDisconnection(event: CloseEvent): void {
    console.log('Signaling client disconnected:', event.code, event.reason);

    this.isConnected = false;
    this.connectionId = null;

    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }

    this.eventHandlers.onConnectionChange?.(false);

    // Attempt to reconnect unless this was a deliberate close
    if (event.code !== 1000) {
      this.attemptReconnect();
    }
  }

  private async attemptReconnect(): Promise<void> {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      this.eventHandlers.onError?.(new Error('Failed to reconnect to signaling server'));
      return;
    }

    this.reconnectAttempts++;
    this.eventHandlers.onReconnecting?.();

    console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts}) in ${this.reconnectDelay}ms`);

    setTimeout(async () => {
      try {
        await this.connect();
        console.log('Reconnected to signaling server');
        this.eventHandlers.onReconnected?.();
      } catch (error) {
        console.error('Reconnection failed:', error);

        // Exponential backoff with jitter
        this.reconnectDelay = Math.min(
          this.reconnectDelay * 2 + Math.random() * 1000,
          this.maxReconnectDelay
        );

        // Try again
        this.attemptReconnect();
      }
    }, this.reconnectDelay);
  }

  private startHeartbeat(): void {
    if (this.heartbeatInterval) return;

    this.heartbeatInterval = setInterval(() => {
      if (this.isConnectedToServer()) {
        this.sendMessage({ type: 'ping' });
      }
    }, this.heartbeatDelay);
  }

  private flushMessageQueue(): void {
    while (this.messageQueue.length > 0 && this.isConnectedToServer()) {
      const message = this.messageQueue.shift();
      if (message) {
        this.sendMessage(message);
      }
    }
  }
}

// Message type definitions for type safety
export const SignalingMessageTypes = {
  // Call lifecycle
  CALL_INITIATE: 'call-initiate',
  CALL_ANSWER: 'call-answer',
  CALL_DECLINE: 'call-decline',
  CALL_END: 'call-end',
  CALL_TIMEOUT: 'call-timeout',

  // WebRTC signaling
  OFFER: 'offer',
  ANSWER: 'answer',
  ICE_CANDIDATE: 'ice-candidate',

  // Media control
  MEDIA_TOGGLE: 'media-toggle',
  SCREEN_SHARE_START: 'screen-share-start',
  SCREEN_SHARE_END: 'screen-share-end',

  // Room management
  JOIN_ROOM: 'join-room',
  LEAVE_ROOM: 'leave-room',
  ROOM_PARTICIPANTS: 'room-participants',

  // Connection management
  PING: 'ping',
  PONG: 'pong',
  CONNECTION_ESTABLISHED: 'connection-established',
  ERROR: 'error'
} as const;

export type SignalingMessageType = typeof SignalingMessageTypes[keyof typeof SignalingMessageTypes];