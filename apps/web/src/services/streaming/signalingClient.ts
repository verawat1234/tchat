type MessageType =
  | 'OFFER'
  | 'ANSWER'
  | 'ICE_CANDIDATE'
  | 'CHAT'
  | 'REACTION'
  | 'VIEWER_JOIN'
  | 'VIEWER_LEAVE';

interface SignalingMessage {
  type: MessageType;
  payload: unknown;
  timestamp: string;
}

interface ChatMessage {
  message_id: string;
  stream_id: string;
  user_id: string;
  message: string;
  timestamp: string;
  moderation_status: string;
}

interface Reaction {
  reaction_id: string;
  stream_id: string;
  user_id: string;
  reaction_type: string;
  timestamp: string;
}

export class SignalingClient {
  private ws: WebSocket | null = null;
  private streamId: string | null = null;
  private token: string | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private messageHandlers: Map<MessageType, ((payload: unknown) => void)[]> = new Map();
  private isConnected = false;
  private heartbeatInterval: number | null = null;

  constructor() {
    // Initialize message handler map
    const messageTypes: MessageType[] = [
      'OFFER',
      'ANSWER',
      'ICE_CANDIDATE',
      'CHAT',
      'REACTION',
      'VIEWER_JOIN',
      'VIEWER_LEAVE',
    ];
    messageTypes.forEach((type) => {
      this.messageHandlers.set(type, []);
    });
  }

  /**
   * Connect to signaling server
   */
  connect(streamId: string, token: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.streamId = streamId;
      this.token = token;

      const wsUrl = `ws://localhost:8080/ws/signaling?token=${token}&stream_id=${streamId}`;
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        console.log('[Signaling] Connected to signaling server');
        this.isConnected = true;
        this.reconnectAttempts = 0;
        this.startHeartbeat();
        resolve();
      };

      this.ws.onmessage = (event) => {
        this.handleMessage(event.data);
      };

      this.ws.onerror = (error) => {
        console.error('[Signaling] WebSocket error:', error);
        reject(error);
      };

      this.ws.onclose = () => {
        console.log('[Signaling] Disconnected from signaling server');
        this.isConnected = false;
        this.stopHeartbeat();
        this.attemptReconnect();
      };
    });
  }

  /**
   * Disconnect from signaling server
   */
  disconnect(): void {
    this.reconnectAttempts = this.maxReconnectAttempts; // Prevent reconnection
    this.stopHeartbeat();

    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }

    this.isConnected = false;
    console.log('[Signaling] Disconnected');
  }

  /**
   * Send signaling message
   */
  send(type: MessageType, payload: unknown): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn('[Signaling] WebSocket not connected');
      return;
    }

    const message: SignalingMessage = {
      type,
      payload,
      timestamp: new Date().toISOString(),
    };

    this.ws.send(JSON.stringify(message));
  }

  /**
   * Register message handler
   */
  on(type: MessageType, handler: (payload: unknown) => void): () => void {
    const handlers = this.messageHandlers.get(type);
    if (handlers) {
      handlers.push(handler);
    }

    // Return unsubscribe function
    return () => {
      const handlers = this.messageHandlers.get(type);
      if (handlers) {
        const index = handlers.indexOf(handler);
        if (index > -1) {
          handlers.splice(index, 1);
        }
      }
    };
  }

  /**
   * Send chat message
   */
  sendChat(message: string): void {
    this.send('CHAT', { message });
  }

  /**
   * Send reaction
   */
  sendReaction(reactionType: string): void {
    this.send('REACTION', { reaction_type: reactionType });
  }

  /**
   * Check if connected
   */
  isSignalingConnected(): boolean {
    return this.isConnected;
  }

  /**
   * Handle incoming message
   */
  private handleMessage(data: string): void {
    try {
      const message: SignalingMessage = JSON.parse(data);

      // Dispatch to registered handlers
      const handlers = this.messageHandlers.get(message.type);
      if (handlers) {
        handlers.forEach((handler) => {
          handler(message.payload);
        });
      }
    } catch (error) {
      console.error('[Signaling] Failed to parse message:', error);
    }
  }

  /**
   * Attempt to reconnect
   */
  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('[Signaling] Max reconnect attempts reached');
      return;
    }

    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1); // Exponential backoff

    console.log(`[Signaling] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

    setTimeout(() => {
      if (this.streamId && this.token) {
        this.connect(this.streamId, this.token).catch((error) => {
          console.error('[Signaling] Reconnect failed:', error);
        });
      }
    }, delay);
  }

  /**
   * Start heartbeat ping/pong
   */
  private startHeartbeat(): void {
    this.heartbeatInterval = window.setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: 'PING' }));
      }
    }, 30000); // 30 seconds
  }

  /**
   * Stop heartbeat
   */
  private stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }
}