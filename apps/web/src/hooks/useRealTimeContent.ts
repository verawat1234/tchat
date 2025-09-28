/**
 * React hooks for real-time content management
 */

import { useState, useEffect, useCallback, useRef } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import type { RootState } from '../store';
import {
  getRealTimeService,
  type RealTimeMessage,
  type ConnectionStatus,
} from '../services/realTimeConnectionService';
import { useGetContentItemQuery } from '../services/content';

export interface RealTimeContentState {
  value: string;
  lastUpdated: string | null;
  syncStatus: 'idle' | 'syncing' | 'error' | 'queued';
  connectionStatus: 'connecting' | 'connected' | 'disconnected' | 'error' | 'reconnecting';
  version: number;
}

export interface RealTimeContentOptions {
  throttleMs?: number;
  preserveLocalChanges?: boolean;
}

/**
 * Hook for managing real-time content updates
 */
export const useRealTimeContent = (
  contentId: string,
  options: RealTimeContentOptions = {}
): RealTimeContentState & {
  updateContent: (value: string) => Promise<void>;
  isUserEditing: boolean;
  setIsUserEditing: (editing: boolean) => void;
} => {
  const {
    throttleMs = 500,
    preserveLocalChanges = true,
  } = options;

  // Local state
  const [value, setValue] = useState<string>('Loading...');
  const [lastUpdated, setLastUpdated] = useState<string | null>(null);
  const [syncStatus, setSyncStatus] = useState<'idle' | 'syncing' | 'error' | 'queued'>('idle');
  const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus['status']>('disconnected');
  const [version, setVersion] = useState<number>(0);
  const [isUserEditing, setIsUserEditing] = useState<boolean>(false);

  // Refs for throttling and state management
  const throttleTimer = useRef<NodeJS.Timeout | null>(null);
  const lastUpdateRef = useRef<string | null>(null);

  // RTK Query for initial content load
  const {
    data: contentItem,
    isLoading,
    error: loadError,
  } = useGetContentItemQuery(contentId);

  // Get real-time service
  const realTimeService = getRealTimeService();

  // Initialize content from API
  useEffect(() => {
    if (contentItem?.data) {
      const newValue = typeof contentItem.data === 'string' ? contentItem.data : 'No content';
      setValue(newValue);
      setLastUpdated(contentItem.metadata?.updatedAt || null);
      setVersion(contentItem.metadata?.version || 0);
    }
  }, [contentItem]);

  // Handle real-time content updates
  const handleContentUpdate = useCallback((message: RealTimeMessage) => {
    if (message.data.contentId !== contentId) {
      return; // Not for this content
    }

    // Throttle updates to prevent UI flooding
    if (throttleTimer.current) {
      clearTimeout(throttleTimer.current);
    }

    throttleTimer.current = setTimeout(() => {
      // Don't update if user is actively editing and preservation is enabled
      if (preserveLocalChanges && isUserEditing) {
        console.log('Skipping real-time update - user is editing');
        return;
      }

      // Prevent duplicate updates
      if (lastUpdateRef.current === message.data.updatedAt) {
        return;
      }

      lastUpdateRef.current = message.data.updatedAt;

      setValue(message.data.value || 'No content');
      setLastUpdated(message.data.updatedAt || new Date().toISOString());
      setVersion(message.data.version || version + 1);
      setSyncStatus('idle');
    }, throttleMs);
  }, [contentId, throttleMs, preserveLocalChanges, isUserEditing, version]);

  // Handle connection status updates
  const handleStatusUpdate = useCallback((status: ConnectionStatus) => {
    setConnectionStatus(status.status);
  }, []);

  // Subscribe to real-time updates
  useEffect(() => {
    if (!realTimeService) {
      return;
    }

    const unsubscribeContent = realTimeService.subscribe('content_updated', handleContentUpdate);
    const unsubscribeStatus = realTimeService.onStatusChange(handleStatusUpdate);

    // Set initial status
    setConnectionStatus(realTimeService.getStatus().status);

    return () => {
      unsubscribeContent();
      unsubscribeStatus();
      if (throttleTimer.current) {
        clearTimeout(throttleTimer.current);
      }
    };
  }, [realTimeService, handleContentUpdate, handleStatusUpdate]);

  // Function to update content
  const updateContent = useCallback(async (newValue: string): Promise<void> => {
    setSyncStatus('syncing');

    try {
      // Send real-time update if connected
      if (realTimeService?.isConnected()) {
        const updateMessage: RealTimeMessage = {
          type: 'content_update',
          data: {
            contentId,
            value: newValue,
            timestamp: new Date().toISOString(),
          },
        };

        const sent = realTimeService.send(updateMessage);
        if (!sent) {
          throw new Error('Failed to send real-time update');
        }
      } else {
        // Queue for later if offline
        setSyncStatus('queued');
        console.log('Update queued - service offline');
        return;
      }

      // Update local state optimistically
      setValue(newValue);
      setLastUpdated(new Date().toISOString());
      setVersion(prev => prev + 1);
      setSyncStatus('idle');

    } catch (error) {
      console.error('Failed to update content:', error);
      setSyncStatus('error');
      throw error;
    }
  }, [contentId, realTimeService]);

  return {
    value,
    lastUpdated,
    syncStatus,
    connectionStatus,
    version,
    updateContent,
    isUserEditing,
    setIsUserEditing,
  };
};

/**
 * Hook for managing real-time connection status
 */
export const useRealTimeConnectionStatus = () => {
  const [status, setStatus] = useState<ConnectionStatus>({ status: 'disconnected' });
  const realTimeService = getRealTimeService();

  useEffect(() => {
    if (!realTimeService) {
      return;
    }

    const unsubscribe = realTimeService.onStatusChange(setStatus);

    // Set initial status
    setStatus(realTimeService.getStatus());

    return unsubscribe;
  }, [realTimeService]);

  return {
    status: status.status,
    lastConnectedAt: status.lastConnectedAt,
    reconnectAttempts: status.reconnectAttempts,
    isConnected: status.status === 'connected',
    isConnecting: status.status === 'connecting',
    isReconnecting: status.status === 'reconnecting',
    hasError: status.status === 'error',
  };
};

/**
 * Hook for cross-tab content synchronization
 */
export const useCrossTabSync = (contentId: string) => {
  const [updates, setUpdates] = useState(0);

  useEffect(() => {
    const handleStorageEvent = (event: StorageEvent) => {
      if (event.key === 'content_update' && event.newValue) {
        try {
          const updateData = JSON.parse(event.newValue);
          if (updateData.contentId === contentId) {
            setUpdates(prev => prev + 1);

            // Dispatch custom event for components to handle
            window.dispatchEvent(new CustomEvent('cross-tab-content-update', {
              detail: updateData
            }));
          }
        } catch (error) {
          console.error('Failed to parse cross-tab update:', error);
        }
      }
    };

    window.addEventListener('storage', handleStorageEvent);

    return () => {
      window.removeEventListener('storage', handleStorageEvent);
    };
  }, [contentId]);

  const broadcastUpdate = useCallback((value: string) => {
    const updateData = {
      contentId,
      value,
      timestamp: new Date().toISOString(),
    };

    localStorage.setItem('content_update', JSON.stringify(updateData));

    // Clear after a brief moment to prevent accumulation
    setTimeout(() => {
      localStorage.removeItem('content_update');
    }, 100);
  }, [contentId]);

  return {
    updates,
    broadcastUpdate,
  };
};