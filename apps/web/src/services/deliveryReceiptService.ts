/**
 * Delivery Receipt Service
 *
 * Handles message delivery receipts, read status tracking, and real-time
 * status updates across chat participants.
 *
 * Features:
 * - Automatic delivery receipt generation
 * - Read status tracking and broadcasting
 * - Bulk status updates for chat optimization
 * - Real-time status synchronization
 * - Privacy controls for read receipts
 */

import { getRealTimeService, type RealTimeMessage } from './realTimeConnectionService';
import { DeliveryStatus } from '../types/MessageTypes';

// =============================================================================
// Type Definitions
// =============================================================================

export interface DeliveryReceipt {
  id: string;
  messageId: string;
  chatId: string;
  userId: string;
  userName: string;
  userAvatar?: string;
  status: DeliveryStatus;
  timestamp: string;
  deviceInfo?: {
    platform: string;
    version: string;
  };
}

export interface ReadReceiptSettings {
  sendReadReceipts: boolean;
  showReadReceipts: boolean;
  showOnlineStatus: boolean;
  showTypingIndicators: boolean;
}

export interface MessageStatusUpdate {
  messageId: string;
  chatId: string;
  status: DeliveryStatus;
  timestamp: string;
  userId?: string;
  userName?: string;
  userAvatar?: string;
}

export interface BulkStatusUpdate {
  chatId: string;
  messageIds: string[];
  status: DeliveryStatus;
  timestamp: string;
}

// =============================================================================
// Delivery Receipt Service Class
// =============================================================================

export class DeliveryReceiptService {
  private realTimeService = getRealTimeService();
  private receiptCache: Map<string, DeliveryReceipt[]> = new Map();
  private statusSubscribers: Map<string, Array<(update: MessageStatusUpdate) => void>> = new Map();
  private bulkStatusSubscribers: Map<string, Array<(updates: MessageStatusUpdate[]) => void>> = new Map();
  private settings: ReadReceiptSettings;
  private currentUserId: string;
  private currentUserName: string;
  private pendingReceipts: Map<string, DeliveryReceipt> = new Map();

  constructor() {
    this.settings = this.loadSettings();
    this.currentUserId = 'current-user'; // Should come from auth context
    this.currentUserName = 'Current User'; // Should come from auth context
    this.initializeRealTimeSubscriptions();
    this.loadCachedReceipts();
    this.startPeriodicSync();
  }

  // =========================================================================
  // Public API - Receipt Management
  // =========================================================================

  /**
   * Mark a message as delivered
   */
  async markAsDelivered(messageId: string, chatId: string): Promise<void> {
    if (!this.settings.sendReadReceipts) {
      return;
    }

    const receipt: DeliveryReceipt = {
      id: this.generateReceiptId(),
      messageId,
      chatId,
      userId: this.currentUserId,
      userName: this.currentUserName,
      status: DeliveryStatus.DELIVERED,
      timestamp: new Date().toISOString(),
      deviceInfo: {
        platform: navigator.platform,
        version: navigator.userAgent,
      },
    };

    await this.sendReceipt(receipt);
  }

  /**
   * Mark a message as read
   */
  async markAsRead(messageId: string, chatId: string): Promise<void> {
    if (!this.settings.sendReadReceipts) {
      return;
    }

    const receipt: DeliveryReceipt = {
      id: this.generateReceiptId(),
      messageId,
      chatId,
      userId: this.currentUserId,
      userName: this.currentUserName,
      status: DeliveryStatus.READ,
      timestamp: new Date().toISOString(),
      deviceInfo: {
        platform: navigator.platform,
        version: navigator.userAgent,
      },
    };

    await this.sendReceipt(receipt);
  }

  /**
   * Mark multiple messages as read (bulk operation)
   */
  async markMultipleAsRead(messageIds: string[], chatId: string): Promise<void> {
    if (!this.settings.sendReadReceipts || messageIds.length === 0) {
      return;
    }

    const timestamp = new Date().toISOString();
    const receipts: DeliveryReceipt[] = messageIds.map(messageId => ({
      id: this.generateReceiptId(),
      messageId,
      chatId,
      userId: this.currentUserId,
      userName: this.currentUserName,
      status: DeliveryStatus.READ,
      timestamp,
      deviceInfo: {
        platform: navigator.platform,
        version: navigator.userAgent,
      },
    }));

    await this.sendBulkReceipts(receipts);
  }

  /**
   * Get delivery receipts for a message
   */
  getMessageReceipts(messageId: string): DeliveryReceipt[] {
    return Array.from(this.receiptCache.values())
      .flat()
      .filter(receipt => receipt.messageId === messageId)
      .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());
  }

  /**
   * Get read receipts for a message (users who have read it)
   */
  getMessageReadReceipts(messageId: string): DeliveryReceipt[] {
    return this.getMessageReceipts(messageId)
      .filter(receipt => receipt.status === DeliveryStatus.READ);
  }

  /**
   * Get delivery status for a message
   */
  getMessageDeliveryStatus(messageId: string): DeliveryStatus {
    const receipts = this.getMessageReceipts(messageId);

    if (receipts.length === 0) {
      return DeliveryStatus.SENT;
    }

    // If any user has read the message
    if (receipts.some(r => r.status === DeliveryStatus.READ)) {
      return DeliveryStatus.READ;
    }

    // If any user has received the message
    if (receipts.some(r => r.status === DeliveryStatus.DELIVERED)) {
      return DeliveryStatus.DELIVERED;
    }

    return DeliveryStatus.SENT;
  }

  /**
   * Get read status summary for multiple messages
   */
  getBulkReadStatus(messageIds: string[]): Record<string, DeliveryStatus> {
    const result: Record<string, DeliveryStatus> = {};

    messageIds.forEach(messageId => {
      result[messageId] = this.getMessageDeliveryStatus(messageId);
    });

    return result;
  }

  // =========================================================================
  // Settings Management
  // =========================================================================

  /**
   * Update read receipt settings
   */
  updateSettings(newSettings: Partial<ReadReceiptSettings>): void {
    this.settings = { ...this.settings, ...newSettings };
    this.saveSettings();
  }

  /**
   * Get current settings
   */
  getSettings(): ReadReceiptSettings {
    return { ...this.settings };
  }

  // =========================================================================
  // Subscriptions
  // =========================================================================

  /**
   * Subscribe to status updates for a specific message
   */
  subscribeToMessageStatus(messageId: string, callback: (update: MessageStatusUpdate) => void): () => void {
    if (!this.statusSubscribers.has(messageId)) {
      this.statusSubscribers.set(messageId, []);
    }

    const subscribers = this.statusSubscribers.get(messageId)!;
    subscribers.push(callback);

    return () => {
      const index = subscribers.indexOf(callback);
      if (index > -1) {
        subscribers.splice(index, 1);
      }
    };
  }

  /**
   * Subscribe to bulk status updates for a chat
   */
  subscribeToBulkStatus(chatId: string, callback: (updates: MessageStatusUpdate[]) => void): () => void {
    if (!this.bulkStatusSubscribers.has(chatId)) {
      this.bulkStatusSubscribers.set(chatId, []);
    }

    const subscribers = this.bulkStatusSubscribers.get(chatId)!;
    subscribers.push(callback);

    return () => {
      const index = subscribers.indexOf(callback);
      if (index > -1) {
        subscribers.splice(index, 1);
      }
    };
  }

  // =========================================================================
  // Visibility and Privacy
  // =========================================================================

  /**
   * Check if read receipts should be shown for a message
   */
  shouldShowReadReceipts(chatId: string): boolean {
    return this.settings.showReadReceipts;
  }

  /**
   * Mark messages as viewed (when chat is opened)
   */
  async markChatAsViewed(chatId: string, messageIds: string[]): Promise<void> {
    if (!this.settings.sendReadReceipts) {
      return;
    }

    // Mark unread messages as read
    const unreadMessageIds = messageIds.filter(messageId => {
      const status = this.getMessageDeliveryStatus(messageId);
      return status !== DeliveryStatus.READ;
    });

    if (unreadMessageIds.length > 0) {
      await this.markMultipleAsRead(unreadMessageIds, chatId);
    }
  }

  // =========================================================================
  // Private Methods
  // =========================================================================

  private async sendReceipt(receipt: DeliveryReceipt): Promise<void> {
    // Store receipt immediately
    this.storeReceipt(receipt);

    try {
      if (this.realTimeService?.isConnected()) {
        const sent = this.realTimeService.send({
          type: 'delivery_receipt',
          data: receipt,
        });

        if (!sent) {
          // Store for retry later
          this.pendingReceipts.set(receipt.id, receipt);
        }
      } else {
        // Store for retry when connection is restored
        this.pendingReceipts.set(receipt.id, receipt);
      }

      // Notify subscribers
      this.notifyStatusUpdate({
        messageId: receipt.messageId,
        chatId: receipt.chatId,
        status: receipt.status,
        timestamp: receipt.timestamp,
        userId: receipt.userId,
        userName: receipt.userName,
      });
    } catch (error) {
      console.error('Failed to send delivery receipt:', error);
      this.pendingReceipts.set(receipt.id, receipt);
    }
  }

  private async sendBulkReceipts(receipts: DeliveryReceipt[]): Promise<void> {
    // Store receipts immediately
    receipts.forEach(receipt => this.storeReceipt(receipt));

    try {
      if (this.realTimeService?.isConnected()) {
        const sent = this.realTimeService.send({
          type: 'bulk_delivery_receipts',
          data: receipts,
        });

        if (!sent) {
          receipts.forEach(receipt => {
            this.pendingReceipts.set(receipt.id, receipt);
          });
        }
      } else {
        receipts.forEach(receipt => {
          this.pendingReceipts.set(receipt.id, receipt);
        });
      }

      // Notify bulk subscribers
      const updates: MessageStatusUpdate[] = receipts.map(receipt => ({
        messageId: receipt.messageId,
        chatId: receipt.chatId,
        status: receipt.status,
        timestamp: receipt.timestamp,
        userId: receipt.userId,
        userName: receipt.userName,
      }));

      this.notifyBulkStatusUpdate(receipts[0].chatId, updates);
    } catch (error) {
      console.error('Failed to send bulk delivery receipts:', error);
      receipts.forEach(receipt => {
        this.pendingReceipts.set(receipt.id, receipt);
      });
    }
  }

  private initializeRealTimeSubscriptions(): void {
    if (!this.realTimeService) {
      return;
    }

    // Delivery receipt received
    this.realTimeService.subscribe('delivery_receipt', (message: RealTimeMessage) => {
      const receipt = message.data as DeliveryReceipt;
      this.handleReceivedReceipt(receipt);
    });

    // Bulk delivery receipts received
    this.realTimeService.subscribe('bulk_delivery_receipts', (message: RealTimeMessage) => {
      const receipts = message.data as DeliveryReceipt[];
      receipts.forEach(receipt => this.handleReceivedReceipt(receipt));
    });

    // Connection restored - send pending receipts
    this.realTimeService.subscribe('connection_restored', () => {
      this.sendPendingReceipts();
    });
  }

  private handleReceivedReceipt(receipt: DeliveryReceipt): void {
    // Don't process our own receipts
    if (receipt.userId === this.currentUserId) {
      return;
    }

    this.storeReceipt(receipt);

    // Notify subscribers
    this.notifyStatusUpdate({
      messageId: receipt.messageId,
      chatId: receipt.chatId,
      status: receipt.status,
      timestamp: receipt.timestamp,
      userId: receipt.userId,
      userName: receipt.userName,
    });
  }

  private storeReceipt(receipt: DeliveryReceipt): void {
    const chatReceipts = this.receiptCache.get(receipt.chatId) || [];

    // Remove existing receipt from same user for same message
    const filteredReceipts = chatReceipts.filter(
      r => !(r.messageId === receipt.messageId && r.userId === receipt.userId)
    );

    filteredReceipts.push(receipt);
    this.receiptCache.set(receipt.chatId, filteredReceipts);
    this.persistReceipts();
  }

  private notifyStatusUpdate(update: MessageStatusUpdate): void {
    const subscribers = this.statusSubscribers.get(update.messageId) || [];
    subscribers.forEach(callback => {
      try {
        callback(update);
      } catch (error) {
        console.error('Error in status subscriber:', error);
      }
    });
  }

  private notifyBulkStatusUpdate(chatId: string, updates: MessageStatusUpdate[]): void {
    const subscribers = this.bulkStatusSubscribers.get(chatId) || [];
    subscribers.forEach(callback => {
      try {
        callback(updates);
      } catch (error) {
        console.error('Error in bulk status subscriber:', error);
      }
    });
  }

  private async sendPendingReceipts(): Promise<void> {
    const pending = Array.from(this.pendingReceipts.values());
    this.pendingReceipts.clear();

    for (const receipt of pending) {
      await this.sendReceipt(receipt);
    }
  }

  private startPeriodicSync(): void {
    // Sync pending receipts every 30 seconds
    setInterval(() => {
      if (this.pendingReceipts.size > 0 && this.realTimeService?.isConnected()) {
        this.sendPendingReceipts();
      }
    }, 30000);
  }

  private persistReceipts(): void {
    try {
      const allReceipts: Record<string, DeliveryReceipt[]> = {};
      this.receiptCache.forEach((receipts, chatId) => {
        // Keep only recent receipts (last 1000 per chat)
        allReceipts[chatId] = receipts
          .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
          .slice(0, 1000);
      });

      localStorage.setItem('tchat_delivery_receipts', JSON.stringify(allReceipts));
    } catch (error) {
      console.error('Failed to persist delivery receipts:', error);
    }
  }

  private loadCachedReceipts(): void {
    try {
      const stored = localStorage.getItem('tchat_delivery_receipts');
      if (stored) {
        const allReceipts: Record<string, DeliveryReceipt[]> = JSON.parse(stored);
        Object.entries(allReceipts).forEach(([chatId, receipts]) => {
          this.receiptCache.set(chatId, receipts);
        });
      }
    } catch (error) {
      console.error('Failed to load cached delivery receipts:', error);
    }
  }

  private loadSettings(): ReadReceiptSettings {
    try {
      const stored = localStorage.getItem('tchat_receipt_settings');
      if (stored) {
        return JSON.parse(stored);
      }
    } catch (error) {
      console.error('Failed to load receipt settings:', error);
    }

    // Default settings
    return {
      sendReadReceipts: true,
      showReadReceipts: true,
      showOnlineStatus: true,
      showTypingIndicators: true,
    };
  }

  private saveSettings(): void {
    try {
      localStorage.setItem('tchat_receipt_settings', JSON.stringify(this.settings));
    } catch (error) {
      console.error('Failed to save receipt settings:', error);
    }
  }

  private generateReceiptId(): string {
    return `receipt_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
}

// =============================================================================
// Singleton Instance and React Hooks
// =============================================================================

export const deliveryReceiptService = new DeliveryReceiptService();

// React hook for delivery receipts
import { useState, useEffect, useCallback } from 'react';

export function useDeliveryReceipts(messageId?: string, chatId?: string) {
  const [receipts, setReceipts] = useState<DeliveryReceipt[]>([]);
  const [status, setStatus] = useState<DeliveryStatus>(DeliveryStatus.SENT);

  useEffect(() => {
    if (messageId) {
      const messageReceipts = deliveryReceiptService.getMessageReceipts(messageId);
      setReceipts(messageReceipts);
      setStatus(deliveryReceiptService.getMessageDeliveryStatus(messageId));

      // Subscribe to status updates
      const unsubscribe = deliveryReceiptService.subscribeToMessageStatus(messageId, (update) => {
        setStatus(update.status);
        const updatedReceipts = deliveryReceiptService.getMessageReceipts(messageId);
        setReceipts(updatedReceipts);
      });

      return unsubscribe;
    }
  }, [messageId]);

  const markAsRead = useCallback(async () => {
    if (messageId && chatId) {
      await deliveryReceiptService.markAsRead(messageId, chatId);
    }
  }, [messageId, chatId]);

  const markAsDelivered = useCallback(async () => {
    if (messageId && chatId) {
      await deliveryReceiptService.markAsDelivered(messageId, chatId);
    }
  }, [messageId, chatId]);

  return {
    receipts,
    status,
    readReceipts: receipts.filter(r => r.status === DeliveryStatus.READ),
    markAsRead,
    markAsDelivered,
  };
}