/**
 * Real-time Connection Service
 *
 * Provides WebSocket-based real-time communication for content updates
 * Supports connection management, message handling, and error recovery
 */

export interface RealTimeMessage {
  type: string;
  data: any;
  timestamp?: string;
}

export interface ConnectionStatus {
  status: 'connecting' | 'connected' | 'disconnected' | 'error' | 'reconnecting';
  lastConnectedAt?: string;
  reconnectAttempts?: number;
}

export interface RealTimeConnectionConfig {
  url: string;
  token?: string;
  reconnectDelay?: number;
  maxReconnectAttempts?: number;
  heartbeatInterval?: number;
}

type MessageHandler = (message: RealTimeMessage) => void;
type StatusHandler = (status: ConnectionStatus) => void;

export class RealTimeConnectionService {
  private ws: WebSocket | null = null;
  private config: RealTimeConnectionConfig;
  private messageHandlers: Map<string, MessageHandler[]> = new Map();
  private statusHandlers: StatusHandler[] = [];
  private reconnectTimer: NodeJS.Timeout | null = null;
  private heartbeatTimer: NodeJS.Timeout | null = null;
  private currentStatus: ConnectionStatus = { status: 'disconnected' };
  private reconnectAttempts = 0;

  constructor(config: RealTimeConnectionConfig) {
    this.config = {
      reconnectDelay: 5000,
      maxReconnectAttempts: 5,
      heartbeatInterval: 30000,
      ...config,
    };
  }

  /**
   * Establish WebSocket connection
   */
  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.updateStatus({ status: 'connecting' });

        // Construct WebSocket URL with token if available
        const url = this.config.token
          ? `${this.config.url}?token=${this.config.token}`
          : this.config.url;

        this.ws = new WebSocket(url);

        this.ws.onopen = () => {
          this.updateStatus({
            status: 'connected',
            lastConnectedAt: new Date().toISOString(),
            reconnectAttempts: 0
          });
          this.reconnectAttempts = 0;
          this.startHeartbeat();
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message: RealTimeMessage = JSON.parse(event.data);
            this.handleMessage(message);
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
          }
        };

        this.ws.onclose = (event) => {
          this.cleanup();

          if (event.code === 1000) {
            // Normal closure
            this.updateStatus({ status: 'disconnected' });
          } else {
            // Unexpected closure - attempt reconnection
            this.updateStatus({ status: 'error' });
            this.attemptReconnection();
          }
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.updateStatus({ status: 'error' });
          reject(new Error('WebSocket connection failed'));
        };

      } catch (error) {
        this.updateStatus({ status: 'error' });
        reject(error);
      }
    });
  }

  /**
   * Disconnect WebSocket connection
   */
  disconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.ws) {
      this.ws.close(1000, 'Client disconnect');
    }

    this.cleanup();
    this.updateStatus({ status: 'disconnected' });
  }

  /**
   * Send message through WebSocket
   */
  send(message: RealTimeMessage): boolean {
    if (this.ws?.readyState === WebSocket.OPEN) {
      try {
        this.ws.send(JSON.stringify(message));
        return true;
      } catch (error) {
        console.error('Failed to send WebSocket message:', error);
        return false;
      }
    }
    return false;
  }

  /**
   * Subscribe to messages of a specific type
   */
  subscribe(messageType: string, handler: MessageHandler): () => void {
    if (!this.messageHandlers.has(messageType)) {
      this.messageHandlers.set(messageType, []);
    }

    const handlers = this.messageHandlers.get(messageType)!;
    handlers.push(handler);

    // Return unsubscribe function
    return () => {
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    };
  }

  /**
   * Subscribe to connection status changes
   */
  onStatusChange(handler: StatusHandler): () => void {
    this.statusHandlers.push(handler);

    // Return unsubscribe function
    return () => {
      const index = this.statusHandlers.indexOf(handler);
      if (index > -1) {
        this.statusHandlers.splice(index, 1);
      }
    };
  }

  /**
   * Get current connection status
   */
  getStatus(): ConnectionStatus {
    return { ...this.currentStatus };
  }

  /**
   * Get current connection state
   */
  isConnected(): boolean {
    return this.currentStatus.status === 'connected';
  }

  private handleMessage(message: RealTimeMessage): void {
    // Add timestamp if not present
    if (!message.timestamp) {
      message.timestamp = new Date().toISOString();
    }

    // Handle heartbeat/pong messages
    if (message.type === 'pong') {
      return; // Heartbeat response - no action needed
    }

    // Dispatch to registered handlers
    const handlers = this.messageHandlers.get(message.type) || [];
    handlers.forEach(handler => {
      try {
        handler(message);
      } catch (error) {
        console.error(`Error in message handler for type ${message.type}:`, error);
      }
    });
  }

  private updateStatus(status: Partial<ConnectionStatus>): void {
    this.currentStatus = { ...this.currentStatus, ...status };

    // Notify status handlers
    this.statusHandlers.forEach(handler => {
      try {
        handler(this.currentStatus);
      } catch (error) {
        console.error('Error in status handler:', error);
      }
    });
  }

  private attemptReconnection(): void {
    if (this.reconnectAttempts >= this.config.maxReconnectAttempts!) {
      console.error('Max reconnection attempts reached');
      this.updateStatus({ status: 'error' });
      return;
    }

    this.reconnectAttempts++;
    this.updateStatus({
      status: 'reconnecting',
      reconnectAttempts: this.reconnectAttempts
    });

    this.reconnectTimer = setTimeout(() => {
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.config.maxReconnectAttempts})`);
      this.connect().catch(error => {
        console.error('Reconnection failed:', error);
        this.attemptReconnection();
      });
    }, this.config.reconnectDelay!);
  }

  private startHeartbeat(): void {
    this.heartbeatTimer = setInterval(() => {
      if (this.isConnected()) {
        this.send({
          type: 'ping',
          data: { timestamp: new Date().toISOString() }
        });
      }
    }, this.config.heartbeatInterval!);
  }

  private cleanup(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }

    this.ws = null;
  }
}

/**
 * Singleton instance for global use
 */
export let realTimeService: RealTimeConnectionService | null = null;

/**
 * Initialize the real-time service
 */
export const initializeRealTimeService = (config: RealTimeConnectionConfig): RealTimeConnectionService => {
  realTimeService = new RealTimeConnectionService(config);
  return realTimeService;
};

/**
 * Get the current real-time service instance
 */
export const getRealTimeService = (): RealTimeConnectionService | null => {
  return realTimeService;
};

/**
 * Hook for React components to use real-time connection
 */
export const useRealTimeConnection = () => {
  const service = getRealTimeService();

  if (!service) {
    throw new Error('RealTimeConnectionService not initialized. Call initializeRealTimeService first.');
  }

  return service;
};