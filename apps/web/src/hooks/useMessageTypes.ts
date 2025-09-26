// T067 - Custom React hooks for message types
/**
 * useMessageTypes - Custom hooks for message type operations
 * Provides optimized hooks for common message operations with caching and error handling
 */

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  useCreateMessageMutation,
  useUpdateMessageMutation,
  useDeleteMessageMutation,
  useInteractWithMessageMutation,
  useGetMessagesQuery,
  useGetMessageQuery,
  useValidateMessageMutation,
  useSearchMessagesQuery,
  useLazySearchMessagesQuery,
  useGetMessageAnalyticsQuery,
  useGetThreadQuery,
  useBulkUpdateMessagesMutation,
  messageTypesApi,
  CreateMessageRequest,
  MessageInteractionRequest,
  MessageSearchRequest,
} from '../services/api/messageTypes';
import { MessageData, MessageType, MessageContent } from '../types/MessageData';
import { useMessageRegistry } from '../components/MessageTypes/MessageRegistry';
import { RootState } from '../store';

// Hook for creating messages with validation and optimistic updates
export const useCreateMessage = () => {
  const [createMessage, { isLoading, error }] = useCreateMessageMutation();
  const [validateMessage] = useValidateMessageMutation();
  const registry = useMessageRegistry();

  const createMessageWithValidation = useCallback(
    async (request: CreateMessageRequest) => {
      try {
        // Pre-validate using registry
        const registryValidation = registry.validateMessage({
          id: 'temp',
          senderId: 'current-user',
          senderName: 'Current User',
          timestamp: new Date(),
          type: request.type,
          isOwn: true,
          content: request.content,
        });

        if (!registryValidation.isValid) {
          throw new Error(
            `Validation failed: ${registryValidation.errors.map(e => e.message).join(', ')}`
          );
        }

        // Server-side validation
        const serverValidation = await validateMessage({
          type: request.type,
          content: request.content,
          chatId: request.chatId,
        }).unwrap();

        if (!serverValidation.isValid) {
          throw new Error(
            `Server validation failed: ${serverValidation.errors.map(e => e.message).join(', ')}`
          );
        }

        // Create message with validated content
        const result = await createMessage({
          ...request,
          content: serverValidation.sanitizedContent || request.content,
        }).unwrap();

        // Record usage analytics
        registry.recordUsage(request.type, performance.now(), false);

        return result;
      } catch (err) {
        // Record error in analytics
        registry.recordUsage(request.type, performance.now(), true);
        throw err;
      }
    },
    [createMessage, validateMessage, registry]
  );

  return {
    createMessage: createMessageWithValidation,
    isLoading,
    error,
  };
};

// Hook for message interactions with optimistic updates
export const useMessageInteraction = () => {
  const [interactWithMessage, { isLoading, error }] = useInteractWithMessageMutation();
  const dispatch = useDispatch();

  const interact = useCallback(
    async (request: MessageInteractionRequest) => {
      try {
        // Optimistic update for better UX
        const optimisticUpdate = dispatch(
          messageTypesApi.util.updateQueryData(
            'getMessage',
            { messageId: request.messageId },
            (draft) => {
              // Apply optimistic changes based on interaction type
              switch (request.interactionType) {
                case 'react':
                  // Add reaction optimistically
                  if (!draft.content.reactions) {
                    draft.content.reactions = [];
                  }
                  draft.content.reactions.push({
                    emoji: request.data.emoji,
                    userId: 'current-user',
                    timestamp: new Date().toISOString(),
                  });
                  break;
                case 'vote':
                  // Update vote count optimistically
                  if (draft.content.options) {
                    const option = draft.content.options.find(
                      opt => opt.id === request.data.optionId
                    );
                    if (option) {
                      option.votes = (option.votes || 0) + 1;
                      draft.content.userVote = request.data.optionId;
                    }
                  }
                  break;
                default:
                  // For other interactions, just mark as interacted
                  draft.content.hasInteracted = true;
              }
            }
          )
        );

        // Perform actual interaction
        const result = await interactWithMessage(request).unwrap();

        return result;
      } catch (err) {
        // Revert optimistic update on error
        dispatch(
          messageTypesApi.util.updateQueryData(
            'getMessage',
            { messageId: request.messageId },
            (draft) => {
              // Revert optimistic changes
              // This is a simplified revert - in production, you'd want more sophisticated rollback
              draft.content.hasInteracted = false;
            }
          )
        );
        throw err;
      }
    },
    [interactWithMessage, dispatch]
  );

  return {
    interact,
    isLoading,
    error,
  };
};

// Hook for paginated message loading with infinite scroll support
export const useMessages = (chatId: string, options: {
  messageTypes?: MessageType[];
  includeThreads?: boolean;
  pageSize?: number;
}) => {
  const { messageTypes, includeThreads = false, pageSize = 50 } = options;
  const [cursor, setCursor] = useState<string | undefined>();
  const [allMessages, setAllMessages] = useState<MessageData[]>([]);

  const {
    data,
    isLoading,
    isFetching,
    error,
    refetch
  } = useGetMessagesQuery({
    chatId,
    messageTypes,
    includeThreads,
    limit: pageSize,
    cursor,
  });

  // Accumulate messages for infinite scroll
  useEffect(() => {
    if (data?.messages) {
      if (!cursor) {
        // First page - replace all messages
        setAllMessages(data.messages);
      } else {
        // Subsequent pages - append messages
        setAllMessages(prev => [...prev, ...data.messages]);
      }
    }
  }, [data?.messages, cursor]);

  const loadMore = useCallback(() => {
    if (data?.nextCursor && !isFetching) {
      setCursor(data.nextCursor);
    }
  }, [data?.nextCursor, isFetching]);

  const hasMore = useMemo(() => !!data?.nextCursor, [data?.nextCursor]);

  const reset = useCallback(() => {
    setCursor(undefined);
    setAllMessages([]);
    refetch();
  }, [refetch]);

  return {
    messages: allMessages,
    isLoading: isLoading && !cursor, // Only show loading for initial load
    isFetching,
    error,
    loadMore,
    hasMore,
    reset,
    totalCount: data?.totalCount || 0,
  };
};

// Hook for message search with debouncing and caching
export const useMessageSearch = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchFilters, setSearchFilters] = useState<Partial<MessageSearchRequest>>({});
  const debouncedQuery = useDebounce(searchQuery, 300);
  const [triggerSearch, { data, isLoading, error, isFetching }] = useLazySearchMessagesQuery();

  // Trigger search when query or filters change
  useEffect(() => {
    if (debouncedQuery || Object.keys(searchFilters).length > 0) {
      triggerSearch({
        query: debouncedQuery,
        ...searchFilters,
      });
    }
  }, [debouncedQuery, searchFilters, triggerSearch]);

  const updateFilters = useCallback((filters: Partial<MessageSearchRequest>) => {
    setSearchFilters(prev => ({ ...prev, ...filters }));
  }, []);

  const clearFilters = useCallback(() => {
    setSearchFilters({});
    setSearchQuery('');
  }, []);

  return {
    query: searchQuery,
    setQuery: setSearchQuery,
    filters: searchFilters,
    updateFilters,
    clearFilters,
    results: data?.messages || [],
    facets: data?.facets,
    suggestions: data?.suggestions || [],
    isLoading: isLoading || isFetching,
    error,
    totalCount: data?.totalCount || 0,
  };
};

// Hook for message analytics with caching and refresh capabilities
export const useMessageAnalytics = (
  chatId: string,
  timeRange: { start: string; end: string },
  options?: {
    messageTypes?: MessageType[];
    groupBy?: 'day' | 'week' | 'month';
    autoRefresh?: boolean;
    refreshInterval?: number;
  }
) => {
  const {
    messageTypes,
    groupBy = 'day',
    autoRefresh = false,
    refreshInterval = 300000, // 5 minutes
  } = options || {};

  const {
    data,
    isLoading,
    error,
    refetch
  } = useGetMessageAnalyticsQuery(
    {
      chatId,
      messageTypes,
      timeRange,
      groupBy,
    },
    {
      // Cache for 5 minutes
      pollingInterval: autoRefresh ? refreshInterval : undefined,
    }
  );

  // Computed analytics
  const analytics = useMemo(() => {
    if (!data) return null;

    const totalMessages = data.messageStats.reduce((sum, stat) => sum + stat.count, 0);
    const averageEngagement = data.messageStats.reduce(
      (sum, stat) => sum + stat.averageEngagement * stat.count,
      0
    ) / totalMessages;

    const topMessageTypes = [...data.messageStats]
      .sort((a, b) => b.averageEngagement - a.averageEngagement)
      .slice(0, 3);

    return {
      ...data,
      totalMessages,
      averageEngagement,
      topMessageTypes,
    };
  }, [data]);

  return {
    analytics,
    isLoading,
    error,
    refetch,
  };
};

// Hook for thread management
export const useThread = (threadId: string | null, options?: {
  autoLoad?: boolean;
  pageSize?: number;
}) => {
  const { autoLoad = true, pageSize = 20 } = options || {};
  const [offset, setOffset] = useState(0);

  const {
    data,
    isLoading,
    error,
    refetch
  } = useGetThreadQuery(
    {
      threadId: threadId!,
      limit: pageSize,
      offset,
    },
    {
      skip: !threadId || !autoLoad,
    }
  );

  const loadMore = useCallback(() => {
    if (data?.messages && data.messages.length === pageSize) {
      setOffset(prev => prev + pageSize);
    }
  }, [data?.messages, pageSize]);

  const reset = useCallback(() => {
    setOffset(0);
    refetch();
  }, [refetch]);

  return {
    thread: data?.thread,
    messages: data?.messages || [],
    isLoading,
    error,
    loadMore,
    reset,
  };
};

// Hook for bulk operations
export const useBulkMessageOperations = () => {
  const [bulkUpdate, { isLoading, error }] = useBulkUpdateMessagesMutation();
  const [selectedMessages, setSelectedMessages] = useState<Set<string>>(new Set());

  const toggleSelection = useCallback((messageId: string) => {
    setSelectedMessages(prev => {
      const newSet = new Set(prev);
      if (newSet.has(messageId)) {
        newSet.delete(messageId);
      } else {
        newSet.add(messageId);
      }
      return newSet;
    });
  }, []);

  const selectAll = useCallback((messageIds: string[]) => {
    setSelectedMessages(new Set(messageIds));
  }, []);

  const clearSelection = useCallback(() => {
    setSelectedMessages(new Set());
  }, []);

  const bulkOperation = useCallback(
    async (operation: 'update' | 'delete' | 'archive', updates?: Partial<MessageContent>) => {
      if (selectedMessages.size === 0) return;

      const result = await bulkUpdate({
        messageIds: Array.from(selectedMessages),
        updates: updates || {},
        operation,
      }).unwrap();

      // Clear selection after successful operation
      clearSelection();

      return result;
    },
    [bulkUpdate, selectedMessages, clearSelection]
  );

  return {
    selectedMessages,
    toggleSelection,
    selectAll,
    clearSelection,
    bulkOperation,
    isLoading,
    error,
    selectedCount: selectedMessages.size,
  };
};

// Helper hooks
const useDebounce = (value: string, delay: number) => {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
};

// Hook for message component performance monitoring
export const useMessagePerformance = (messageType: MessageType) => {
  const startTime = useRef<number>();
  const registry = useMessageRegistry();

  useEffect(() => {
    startTime.current = performance.now();

    return () => {
      if (startTime.current) {
        const renderTime = performance.now() - startTime.current;
        registry.recordUsage(messageType, renderTime, false);
      }
    };
  }, [messageType, registry]);
};

// Hook for real-time message updates via WebSocket
export const useMessageUpdates = (chatId: string) => {
  const dispatch = useDispatch();
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    // WebSocket connection for real-time updates
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/ws/messages/${chatId}`;

    wsRef.current = new WebSocket(wsUrl);

    wsRef.current.onmessage = (event) => {
      try {
        const update = JSON.parse(event.data);

        switch (update.type) {
          case 'MESSAGE_CREATED':
            // Add new message to cache
            dispatch(
              messageTypesApi.util.updateQueryData(
                'getMessages',
                { chatId },
                (draft) => {
                  draft.messages.unshift(update.message);
                  draft.totalCount++;
                }
              )
            );
            break;

          case 'MESSAGE_UPDATED':
            // Update existing message in cache
            dispatch(
              messageTypesApi.util.updateQueryData(
                'getMessage',
                { messageId: update.message.id },
                (draft) => {
                  Object.assign(draft, update.message);
                }
              )
            );
            break;

          case 'MESSAGE_DELETED':
            // Remove message from cache
            dispatch(
              messageTypesApi.util.updateQueryData(
                'getMessages',
                { chatId },
                (draft) => {
                  draft.messages = draft.messages.filter(
                    msg => msg.id !== update.messageId
                  );
                  draft.totalCount--;
                }
              )
            );
            break;
        }
      } catch (error) {
        console.error('Failed to process WebSocket message:', error);
      }
    };

    wsRef.current.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    return () => {
      wsRef.current?.close();
    };
  }, [chatId, dispatch]);

  const sendMessage = useCallback((message: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    }
  }, []);

  return { sendMessage };
};