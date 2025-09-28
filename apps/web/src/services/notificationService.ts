/**
 * Notification Service for Real-time Content Updates
 *
 * Provides a comprehensive notification system for content management operations,
 * including real-time updates, offline queue management, and user notifications.
 *
 * Features:
 * - Real-time WebSocket notifications for content changes
 * - Toast notifications for user feedback
 * - Notification preferences and settings management
 * - Offline notification queueing and sync
 * - Content-specific notification filtering
 * - Cross-tab notification synchronization
 */

import { getRealTimeService, type RealTimeMessage } from './realTimeConnectionService';

// =============================================================================
// Type Definitions
// =============================================================================

export type NotificationType =
  | 'content_updated'
  | 'content_created'
  | 'content_published'
  | 'content_archived'
  | 'content_deleted'
  | 'content_conflict'
  | 'content_error'
  | 'sync_complete'
  | 'offline_queue_processed';

export type NotificationPriority = 'low' | 'medium' | 'high' | 'critical';

export type NotificationStatus = 'pending' | 'delivered' | 'read' | 'dismissed' | 'expired';

export interface ContentNotification {
  id: string;
  type: NotificationType;
  priority: NotificationPriority;
  status: NotificationStatus;
  title: string;
  message: string;
  contentId?: string;
  category?: string;
  userId?: string;
  timestamp: string;
  expiresAt?: string;
  actions?: NotificationAction[];
  metadata?: Record<string, any>;
}

export interface NotificationAction {
  id: string;
  label: string;
  type: 'primary' | 'secondary' | 'destructive';
  action: () => void | Promise<void>;
}

export interface NotificationPreferences {
  enabled: boolean;
  categories: {
    [category: string]: boolean;
  };
  types: {
    [type in NotificationType]: boolean;
  };
  delivery: {
    toast: boolean;
    persistent: boolean;
    sound: boolean;
  };
  priority: {
    showLow: boolean;
    showMedium: boolean;
    showHigh: boolean;
    showCritical: boolean;
  };
}

export interface NotificationQueue {
  notifications: ContentNotification[];
  lastProcessed: string;
  retryCount: number;
  maxRetries: number;
}

// =============================================================================
// Notification Service Class
// =============================================================================

export class NotificationService {
  private notifications: Map<string, ContentNotification> = new Map();
  private subscribers: Map<NotificationType, Array<(notification: ContentNotification) => void>> = new Map();
  private globalSubscribers: Array<(notification: ContentNotification) => void> = [];
  private offlineQueue: NotificationQueue = {
    notifications: [],
    lastProcessed: new Date().toISOString(),
    retryCount: 0,
    maxRetries: 3,
  };
  private preferences: NotificationPreferences = this.getDefaultPreferences();
  private realTimeService = getRealTimeService();

  constructor() {
    this.initializeRealTimeSubscriptions();
    this.loadPreferences();
    this.loadOfflineQueue();
    this.setupCrossTabSync();
  }

  // =========================================================================
  // Public API
  // =========================================================================

  /**
   * Subscribe to notifications of a specific type
   */
  subscribe(type: NotificationType, callback: (notification: ContentNotification) => void): () => void {
    if (!this.subscribers.has(type)) {
      this.subscribers.set(type, []);
    }

    const callbacks = this.subscribers.get(type)!;
    callbacks.push(callback);

    return () => {
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    };
  }

  /**
   * Subscribe to all notifications
   */
  subscribeAll(callback: (notification: ContentNotification) => void): () => void {
    this.globalSubscribers.push(callback);

    return () => {
      const index = this.globalSubscribers.indexOf(callback);
      if (index > -1) {
        this.globalSubscribers.splice(index, 1);
      }
    };
  }

  /**
   * Send a notification
   */
  async notify(notification: Omit<ContentNotification, 'id' | 'timestamp' | 'status'>): Promise<string> {
    const fullNotification: ContentNotification = {
      ...notification,
      id: this.generateNotificationId(),
      timestamp: new Date().toISOString(),
      status: 'pending',
    };

    // Check if notification should be shown based on preferences
    if (!this.shouldShowNotification(fullNotification)) {
      return fullNotification.id;
    }

    // Store notification
    this.notifications.set(fullNotification.id, fullNotification);

    // Send via real-time service if connected
    if (this.realTimeService?.isConnected()) {
      try {
        const sent = this.realTimeService.send({
          type: 'notification',
          data: {
            notification: fullNotification,
            timestamp: new Date().toISOString(),
          },
        });

        if (sent) {
          fullNotification.status = 'delivered';
        }
      } catch (error) {
        console.error('Failed to send real-time notification:', error);
        this.queueForOfflineDelivery(fullNotification);
      }
    } else {
      // Queue for offline delivery
      this.queueForOfflineDelivery(fullNotification);
    }

    // Dispatch to subscribers
    this.dispatchNotification(fullNotification);

    // Store in localStorage for persistence
    this.persistNotification(fullNotification);

    return fullNotification.id;
  }

  /**
   * Mark notification as read
   */
  markAsRead(notificationId: string): void {
    const notification = this.notifications.get(notificationId);
    if (notification) {
      notification.status = 'read';
      this.notifications.set(notificationId, notification);
      this.persistNotification(notification);
    }
  }

  /**
   * Dismiss notification
   */
  dismiss(notificationId: string): void {
    const notification = this.notifications.get(notificationId);
    if (notification) {
      notification.status = 'dismissed';
      this.notifications.set(notificationId, notification);
      this.persistNotification(notification);
    }
  }

  /**
   * Clear all notifications
   */
  clearAll(): void {
    this.notifications.clear();
    localStorage.removeItem('tchat_notifications');
  }

  /**
   * Get all notifications
   */
  getAllNotifications(): ContentNotification[] {
    return Array.from(this.notifications.values()).sort(
      (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    );
  }

  /**
   * Get unread notifications
   */
  getUnreadNotifications(): ContentNotification[] {
    return this.getAllNotifications().filter(n => n.status === 'pending' || n.status === 'delivered');
  }

  /**
   * Update notification preferences
   */
  updatePreferences(preferences: Partial<NotificationPreferences>): void {
    this.preferences = { ...this.preferences, ...preferences };
    localStorage.setItem('tchat_notification_preferences', JSON.stringify(this.preferences));
  }

  /**
   * Get current preferences
   */
  getPreferences(): NotificationPreferences {
    return { ...this.preferences };
  }

  /**
   * Process offline notification queue
   */
  async processOfflineQueue(): Promise<void> {
    if (!this.realTimeService?.isConnected() || this.offlineQueue.notifications.length === 0) {
      return;
    }

    const notifications = [...this.offlineQueue.notifications];
    this.offlineQueue.notifications = [];
    this.offlineQueue.retryCount++;

    try {
      for (const notification of notifications) {
        const sent = this.realTimeService.send({
          type: 'notification',
          data: {
            notification,
            timestamp: new Date().toISOString(),
            fromQueue: true,
          },
        });

        if (sent) {
          notification.status = 'delivered';
          this.notifications.set(notification.id, notification);
          this.persistNotification(notification);
        } else {
          // Re-queue if failed
          if (this.offlineQueue.retryCount < this.offlineQueue.maxRetries) {
            this.offlineQueue.notifications.push(notification);
          }
        }
      }

      this.offlineQueue.lastProcessed = new Date().toISOString();
      this.persistOfflineQueue();

      // Notify about successful queue processing
      if (notifications.length > 0) {
        this.notify({
          type: 'offline_queue_processed',
          priority: 'low',
          title: 'Notifications Synced',
          message: `${notifications.length} notification(s) delivered`,
        });
      }
    } catch (error) {
      console.error('Failed to process offline queue:', error);
    }
  }

  // =========================================================================
  // Content-Specific Notification Helpers
  // =========================================================================

  /**
   * Notify about content updates
   */
  async notifyContentUpdate(contentId: string, category: string, updatedBy?: string): Promise<string> {
    return this.notify({
      type: 'content_updated',
      priority: 'medium',
      title: 'Content Updated',
      message: `Content in ${category} has been updated${updatedBy ? ` by ${updatedBy}` : ''}`,
      contentId,
      category,
      actions: [
        {
          id: 'view',
          label: 'View Changes',
          type: 'primary',
          action: () => {
            // Navigate to content
            window.location.href = `/content/${contentId}`;
          },
        },
      ],
    });
  }

  /**
   * Notify about content conflicts
   */
  async notifyContentConflict(contentId: string, category: string): Promise<string> {
    return this.notify({
      type: 'content_conflict',
      priority: 'high',
      title: 'Content Conflict',
      message: `Multiple users are editing content in ${category}. Resolution required.`,
      contentId,
      category,
      actions: [
        {
          id: 'resolve',
          label: 'Resolve Conflict',
          type: 'primary',
          action: () => {
            // Navigate to conflict resolution
            window.location.href = `/content/${contentId}/conflicts`;
          },
        },
      ],
    });
  }

  /**
   * Notify about content errors
   */
  async notifyContentError(error: string, contentId?: string, category?: string): Promise<string> {
    return this.notify({
      type: 'content_error',
      priority: 'high',
      title: 'Content Operation Failed',
      message: error,
      contentId,
      category,
      actions: [
        {
          id: 'retry',
          label: 'Try Again',
          type: 'primary',
          action: () => {
            // Trigger retry logic
            window.location.reload();
          },
        },
      ],
    });
  }

  // =========================================================================
  // Private Methods
  // =========================================================================

  private initializeRealTimeSubscriptions(): void {
    if (!this.realTimeService) {
      return;
    }

    // Subscribe to content update notifications
    this.realTimeService.subscribe('content_updated', (message: RealTimeMessage) => {
      const { contentId, category, updatedBy } = message.data;
      this.notifyContentUpdate(contentId, category, updatedBy);
    });

    // Subscribe to connection status for queue processing
    this.realTimeService.onStatusChange((status) => {
      if (status.status === 'connected') {
        this.processOfflineQueue();
      }
    });
  }

  private shouldShowNotification(notification: ContentNotification): boolean {
    // Check if notifications are enabled
    if (!this.preferences.enabled) {
      return false;
    }

    // Check type preferences
    if (!this.preferences.types[notification.type]) {
      return false;
    }

    // Check category preferences
    if (notification.category && !this.preferences.categories[notification.category]) {
      return false;
    }

    // Check priority preferences
    switch (notification.priority) {
      case 'low':
        return this.preferences.priority.showLow;
      case 'medium':
        return this.preferences.priority.showMedium;
      case 'high':
        return this.preferences.priority.showHigh;
      case 'critical':
        return this.preferences.priority.showCritical;
      default:
        return true;
    }
  }

  private dispatchNotification(notification: ContentNotification): void {
    // Dispatch to type-specific subscribers
    const typeSubscribers = this.subscribers.get(notification.type) || [];
    typeSubscribers.forEach(callback => {
      try {
        callback(notification);
      } catch (error) {
        console.error('Error in notification subscriber:', error);
      }
    });

    // Dispatch to global subscribers
    this.globalSubscribers.forEach(callback => {
      try {
        callback(notification);
      } catch (error) {
        console.error('Error in global notification subscriber:', error);
      }
    });

    // Broadcast to other tabs
    this.broadcastToOtherTabs(notification);
  }

  private queueForOfflineDelivery(notification: ContentNotification): void {
    this.offlineQueue.notifications.push(notification);
    this.persistOfflineQueue();
  }

  private persistNotification(notification: ContentNotification): void {
    const existingNotifications = this.loadNotificationsFromStorage();
    const updatedNotifications = existingNotifications.filter(n => n.id !== notification.id);
    updatedNotifications.push(notification);

    // Keep only recent notifications (last 100)
    const recentNotifications = updatedNotifications
      .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
      .slice(0, 100);

    localStorage.setItem('tchat_notifications', JSON.stringify(recentNotifications));
  }

  private persistOfflineQueue(): void {
    localStorage.setItem('tchat_notification_queue', JSON.stringify(this.offlineQueue));
  }

  private loadNotificationsFromStorage(): ContentNotification[] {
    try {
      const stored = localStorage.getItem('tchat_notifications');
      return stored ? JSON.parse(stored) : [];
    } catch (error) {
      console.error('Failed to load notifications from storage:', error);
      return [];
    }
  }

  private loadPreferences(): void {
    try {
      const stored = localStorage.getItem('tchat_notification_preferences');
      if (stored) {
        this.preferences = { ...this.getDefaultPreferences(), ...JSON.parse(stored) };
      }
    } catch (error) {
      console.error('Failed to load notification preferences:', error);
    }
  }

  private loadOfflineQueue(): void {
    try {
      const stored = localStorage.getItem('tchat_notification_queue');
      if (stored) {
        this.offlineQueue = { ...this.offlineQueue, ...JSON.parse(stored) };
      }
    } catch (error) {
      console.error('Failed to load offline queue:', error);
    }
  }

  private setupCrossTabSync(): void {
    window.addEventListener('storage', (event) => {
      if (event.key === 'tchat_notifications' && event.newValue) {
        try {
          const notifications: ContentNotification[] = JSON.parse(event.newValue);
          notifications.forEach(notification => {
            if (!this.notifications.has(notification.id)) {
              this.notifications.set(notification.id, notification);
              this.dispatchNotification(notification);
            }
          });
        } catch (error) {
          console.error('Failed to sync notifications across tabs:', error);
        }
      }
    });
  }

  private broadcastToOtherTabs(notification: ContentNotification): void {
    // Use localStorage to broadcast to other tabs
    const event = {
      type: 'notification',
      notification,
      timestamp: new Date().toISOString(),
    };

    localStorage.setItem('tchat_notification_broadcast', JSON.stringify(event));

    // Clear after a brief moment to prevent accumulation
    setTimeout(() => {
      localStorage.removeItem('tchat_notification_broadcast');
    }, 100);
  }

  private generateNotificationId(): string {
    return `notification_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  private getDefaultPreferences(): NotificationPreferences {
    return {
      enabled: true,
      categories: {},
      types: {
        content_updated: true,
        content_created: true,
        content_published: true,
        content_archived: true,
        content_deleted: true,
        content_conflict: true,
        content_error: true,
        sync_complete: true,
        offline_queue_processed: false,
      },
      delivery: {
        toast: true,
        persistent: true,
        sound: false,
      },
      priority: {
        showLow: false,
        showMedium: true,
        showHigh: true,
        showCritical: true,
      },
    };
  }
}

// =============================================================================
// Singleton Instance
// =============================================================================

export const notificationService = new NotificationService();

// =============================================================================
// React Hook
// =============================================================================

import { useState, useEffect } from 'react';

export function useNotifications(type?: NotificationType) {
  const [notifications, setNotifications] = useState<ContentNotification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    const updateNotifications = () => {
      const allNotifications = notificationService.getAllNotifications();
      const filteredNotifications = type
        ? allNotifications.filter(n => n.type === type)
        : allNotifications;

      setNotifications(filteredNotifications);
      setUnreadCount(notificationService.getUnreadNotifications().length);
    };

    updateNotifications();

    const unsubscribe = type
      ? notificationService.subscribe(type, updateNotifications)
      : notificationService.subscribeAll(updateNotifications);

    return unsubscribe;
  }, [type]);

  return {
    notifications,
    unreadCount,
    markAsRead: notificationService.markAsRead.bind(notificationService),
    dismiss: notificationService.dismiss.bind(notificationService),
    clearAll: notificationService.clearAll.bind(notificationService),
  };
}