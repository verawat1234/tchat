/**
 * Message Reading Hook
 *
 * Custom hook for managing message reading status and visibility tracking.
 * Automatically marks messages as read when they become visible and handles
 * batch read operations for performance.
 */

import { useEffect, useRef, useCallback } from 'react';
import { deliveryReceiptService } from '../services/deliveryReceiptService';

interface UseMessageReadingOptions {
  enabled?: boolean;
  readDelay?: number; // Delay before marking as read (ms)
  threshold?: number; // Intersection threshold (0-1)
}

export function useMessageReading(
  messageId: string,
  chatId: string,
  isOwnMessage: boolean,
  options: UseMessageReadingOptions = {}
) {
  const {
    enabled = true,
    readDelay = 2000,
    threshold = 0.5
  } = options;

  const elementRef = useRef<HTMLElement>(null);
  const timeoutRef = useRef<NodeJS.Timeout>();
  const hasBeenRead = useRef(false);

  const markAsRead = useCallback(async () => {
    if (hasBeenRead.current || isOwnMessage || !enabled) {
      return;
    }

    try {
      await deliveryReceiptService.markAsRead(messageId, chatId);
      hasBeenRead.current = true;
    } catch (error) {
      console.error('Failed to mark message as read:', error);
    }
  }, [messageId, chatId, isOwnMessage, enabled]);

  const handleIntersection = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      const entry = entries[0];

      if (entry.isIntersecting && !hasBeenRead.current) {
        // Start countdown to mark as read
        timeoutRef.current = setTimeout(() => {
          markAsRead();
        }, readDelay);
      } else if (!entry.isIntersecting && timeoutRef.current) {
        // Cancel if message goes out of view before delay
        clearTimeout(timeoutRef.current);
        timeoutRef.current = undefined;
      }
    },
    [markAsRead, readDelay]
  );

  useEffect(() => {
    if (!enabled || isOwnMessage || !elementRef.current) {
      return;
    }

    const observer = new IntersectionObserver(handleIntersection, {
      threshold,
      rootMargin: '0px'
    });

    const element = elementRef.current;
    observer.observe(element);

    return () => {
      observer.unobserve(element);
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, [enabled, isOwnMessage, handleIntersection, threshold]);

  return {
    elementRef,
    markAsRead,
    hasBeenRead: hasBeenRead.current
  };
}

/**
 * Batch Message Reading Hook
 *
 * Manages reading status for multiple messages in a chat.
 * Optimizes by batching read receipts for better performance.
 */
export function useBatchMessageReading(chatId: string, isEnabled = true) {
  const pendingReads = useRef<Set<string>>(new Set());
  const batchTimeoutRef = useRef<NodeJS.Timeout>();

  const addToReadQueue = useCallback((messageId: string) => {
    if (!isEnabled) return;

    pendingReads.current.add(messageId);

    // Clear existing timeout
    if (batchTimeoutRef.current) {
      clearTimeout(batchTimeoutRef.current);
    }

    // Batch read operations for better performance
    batchTimeoutRef.current = setTimeout(async () => {
      const messageIds = Array.from(pendingReads.current);
      pendingReads.current.clear();

      if (messageIds.length > 0) {
        try {
          await deliveryReceiptService.markMultipleAsRead(messageIds, chatId);
        } catch (error) {
          console.error('Failed to batch mark messages as read:', error);
          // Re-add failed messages to queue for retry
          messageIds.forEach(id => pendingReads.current.add(id));
        }
      }
    }, 1000); // Batch after 1 second
  }, [chatId, isEnabled]);

  const markChatAsViewed = useCallback(async (messageIds: string[]) => {
    if (!isEnabled) return;

    try {
      await deliveryReceiptService.markChatAsViewed(chatId, messageIds);
    } catch (error) {
      console.error('Failed to mark chat as viewed:', error);
    }
  }, [chatId, isEnabled]);

  useEffect(() => {
    return () => {
      if (batchTimeoutRef.current) {
        clearTimeout(batchTimeoutRef.current);
      }
    };
  }, []);

  return {
    addToReadQueue,
    markChatAsViewed,
    pendingCount: pendingReads.current.size
  };
}

/**
 * Message Visibility Hook
 *
 * Tracks when messages become visible for analytics and read receipts.
 */
export function useMessageVisibility(
  messageId: string,
  onVisible?: (messageId: string) => void,
  options: { threshold?: number; once?: boolean } = {}
) {
  const { threshold = 0.5, once = false } = options;
  const elementRef = useRef<HTMLElement>(null);
  const hasBeenVisible = useRef(false);

  const handleIntersection = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      const entry = entries[0];

      if (entry.isIntersecting) {
        if (!hasBeenVisible.current || !once) {
          onVisible?.(messageId);
          hasBeenVisible.current = true;
        }
      }
    },
    [messageId, onVisible, once]
  );

  useEffect(() => {
    if (!elementRef.current) return;

    const observer = new IntersectionObserver(handleIntersection, {
      threshold,
      rootMargin: '0px'
    });

    const element = elementRef.current;
    observer.observe(element);

    return () => {
      observer.unobserve(element);
    };
  }, [handleIntersection, threshold]);

  return {
    elementRef,
    isVisible: hasBeenVisible.current
  };
}