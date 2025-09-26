// T068 - Redux store integration for message types
/**
 * MessageTypes Redux Slice
 * Manages local state for message types UI and interactions
 */

import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { MessageType, MessageData } from '../types/MessageData';

// State interfaces
export interface MessageTypesState {
  // UI State
  selectedMessages: Set<string>;
  activeThread: string | null;
  searchQuery: string;
  searchFilters: MessageSearchFilters;
  composing: ComposingState | null;

  // Interaction State
  interactions: Record<string, MessageInteractionState>;

  // Cache State
  messageCache: Record<string, MessageData>;
  threadCache: Record<string, ThreadCacheEntry>;

  // Performance State
  performanceMetrics: PerformanceMetrics;

  // Error State
  errors: MessageError[];
}

export interface MessageSearchFilters {
  messageTypes: MessageType[];
  dateRange?: {
    start: string;
    end: string;
  };
  senderIds: string[];
  tags: string[];
  hasAttachments?: boolean;
}

export interface ComposingState {
  chatId: string;
  messageType: MessageType;
  content: any;
  replyToId?: string;
  threadId?: string;
  isDraft: boolean;
  lastSaved?: string;
}

export interface MessageInteractionState {
  messageId: string;
  interactionType: string;
  isLoading: boolean;
  error?: string;
  optimisticUpdate?: any;
}

export interface ThreadCacheEntry {
  threadId: string;
  messages: MessageData[];
  lastUpdated: string;
  totalCount: number;
}

export interface PerformanceMetrics {
  averageRenderTime: Record<MessageType, number>;
  renderCount: Record<MessageType, number>;
  errorCount: Record<MessageType, number>;
  lastUpdated: string;
}

export interface MessageError {
  id: string;
  type: 'validation' | 'network' | 'component' | 'interaction';
  messageId?: string;
  messageType?: MessageType;
  message: string;
  timestamp: string;
  resolved: boolean;
}

// Initial state
const initialState: MessageTypesState = {
  selectedMessages: new Set<string>(),
  activeThread: null,
  searchQuery: '',
  searchFilters: {
    messageTypes: [],
    senderIds: [],
    tags: [],
  },
  composing: null,
  interactions: {},
  messageCache: {},
  threadCache: {},
  performanceMetrics: {
    averageRenderTime: {} as Record<MessageType, number>,
    renderCount: {} as Record<MessageType, number>,
    errorCount: {} as Record<MessageType, number>,
    lastUpdated: new Date().toISOString(),
  },
  errors: [],
};

// Utility functions for Set serialization
const setToArray = (set: Set<string>): string[] => Array.from(set);
const arrayToSet = (array: string[]): Set<string> => new Set(array);

// Redux slice
const messageTypesSlice = createSlice({
  name: 'messageTypes',
  initialState,
  reducers: {
    // Selection actions
    toggleMessageSelection: (state, action: PayloadAction<string>) => {
      const messageId = action.payload;
      const selectedArray = setToArray(state.selectedMessages);

      if (selectedArray.includes(messageId)) {
        state.selectedMessages = arrayToSet(
          selectedArray.filter(id => id !== messageId)
        );
      } else {
        state.selectedMessages = arrayToSet([...selectedArray, messageId]);
      }
    },

    selectAllMessages: (state, action: PayloadAction<string[]>) => {
      state.selectedMessages = arrayToSet(action.payload);
    },

    clearMessageSelection: (state) => {
      state.selectedMessages = new Set<string>();
    },

    // Thread actions
    setActiveThread: (state, action: PayloadAction<string | null>) => {
      state.activeThread = action.payload;
    },

    // Search actions
    setSearchQuery: (state, action: PayloadAction<string>) => {
      state.searchQuery = action.payload;
    },

    updateSearchFilters: (state, action: PayloadAction<Partial<MessageSearchFilters>>) => {
      state.searchFilters = { ...state.searchFilters, ...action.payload };
    },

    clearSearchFilters: (state) => {
      state.searchFilters = {
        messageTypes: [],
        senderIds: [],
        tags: [],
      };
      state.searchQuery = '';
    },

    // Composing actions
    startComposing: (state, action: PayloadAction<{
      chatId: string;
      messageType: MessageType;
      replyToId?: string;
      threadId?: string;
    }>) => {
      state.composing = {
        ...action.payload,
        content: {},
        isDraft: false,
      };
    },

    updateComposingContent: (state, action: PayloadAction<any>) => {
      if (state.composing) {
        state.composing.content = action.payload;
        state.composing.isDraft = true;
        state.composing.lastSaved = new Date().toISOString();
      }
    },

    saveDraft: (state) => {
      if (state.composing) {
        state.composing.isDraft = true;
        state.composing.lastSaved = new Date().toISOString();
      }
    },

    clearComposing: (state) => {
      state.composing = null;
    },

    // Interaction actions
    startInteraction: (state, action: PayloadAction<{
      messageId: string;
      interactionType: string;
      optimisticUpdate?: any;
    }>) => {
      const { messageId, interactionType, optimisticUpdate } = action.payload;
      state.interactions[messageId] = {
        messageId,
        interactionType,
        isLoading: true,
        optimisticUpdate,
      };

      // Apply optimistic update to message cache
      if (optimisticUpdate && state.messageCache[messageId]) {
        const message = state.messageCache[messageId];
        Object.assign(message.content, optimisticUpdate);
      }
    },

    completeInteraction: (state, action: PayloadAction<{
      messageId: string;
      result?: any;
      error?: string;
    }>) => {
      const { messageId, result, error } = action.payload;
      const interaction = state.interactions[messageId];

      if (interaction) {
        interaction.isLoading = false;
        interaction.error = error;

        // Update message cache with real result
        if (result && state.messageCache[messageId]) {
          Object.assign(state.messageCache[messageId], result);
        }

        // Remove optimistic update if there was an error
        if (error && interaction.optimisticUpdate && state.messageCache[messageId]) {
          // Revert optimistic update - simplified approach
          delete state.messageCache[messageId].content.hasInteracted;
        }
      }
    },

    clearInteraction: (state, action: PayloadAction<string>) => {
      delete state.interactions[action.payload];
    },

    // Cache actions
    updateMessageCache: (state, action: PayloadAction<MessageData[]>) => {
      action.payload.forEach(message => {
        state.messageCache[message.id] = message;
      });
    },

    removeFromMessageCache: (state, action: PayloadAction<string[]>) => {
      action.payload.forEach(messageId => {
        delete state.messageCache[messageId];
      });
    },

    updateThreadCache: (state, action: PayloadAction<{
      threadId: string;
      messages: MessageData[];
      totalCount: number;
    }>) => {
      const { threadId, messages, totalCount } = action.payload;
      state.threadCache[threadId] = {
        threadId,
        messages,
        totalCount,
        lastUpdated: new Date().toISOString(),
      };
    },

    clearThreadCache: (state, action: PayloadAction<string>) => {
      delete state.threadCache[action.payload];
    },

    // Performance actions
    recordRenderTime: (state, action: PayloadAction<{
      messageType: MessageType;
      renderTime: number;
      hadError?: boolean;
    }>) => {
      const { messageType, renderTime, hadError } = action.payload;
      const metrics = state.performanceMetrics;

      // Update render count
      metrics.renderCount[messageType] = (metrics.renderCount[messageType] || 0) + 1;

      // Update average render time
      const currentAverage = metrics.averageRenderTime[messageType] || 0;
      const currentCount = metrics.renderCount[messageType];
      metrics.averageRenderTime[messageType] =
        (currentAverage * (currentCount - 1) + renderTime) / currentCount;

      // Update error count
      if (hadError) {
        metrics.errorCount[messageType] = (metrics.errorCount[messageType] || 0) + 1;
      }

      metrics.lastUpdated = new Date().toISOString();
    },

    resetPerformanceMetrics: (state) => {
      state.performanceMetrics = {
        averageRenderTime: {} as Record<MessageType, number>,
        renderCount: {} as Record<MessageType, number>,
        errorCount: {} as Record<MessageType, number>,
        lastUpdated: new Date().toISOString(),
      };
    },

    // Error actions
    addError: (state, action: PayloadAction<Omit<MessageError, 'id' | 'timestamp' | 'resolved'>>) => {
      const error: MessageError = {
        ...action.payload,
        id: `error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: new Date().toISOString(),
        resolved: false,
      };
      state.errors.push(error);

      // Keep only last 50 errors to prevent memory issues
      if (state.errors.length > 50) {
        state.errors = state.errors.slice(-50);
      }
    },

    resolveError: (state, action: PayloadAction<string>) => {
      const error = state.errors.find(e => e.id === action.payload);
      if (error) {
        error.resolved = true;
      }
    },

    clearErrors: (state, action: PayloadAction<{
      type?: MessageError['type'];
      messageType?: MessageType;
      resolved?: boolean;
    }> = { payload: {} }) => {
      const { type, messageType, resolved } = action.payload;

      state.errors = state.errors.filter(error => {
        if (type && error.type !== type) return true;
        if (messageType && error.messageType !== messageType) return true;
        if (resolved !== undefined && error.resolved !== resolved) return true;
        return false;
      });
    },
  },
});

// Export actions
export const {
  // Selection
  toggleMessageSelection,
  selectAllMessages,
  clearMessageSelection,

  // Thread
  setActiveThread,

  // Search
  setSearchQuery,
  updateSearchFilters,
  clearSearchFilters,

  // Composing
  startComposing,
  updateComposingContent,
  saveDraft,
  clearComposing,

  // Interaction
  startInteraction,
  completeInteraction,
  clearInteraction,

  // Cache
  updateMessageCache,
  removeFromMessageCache,
  updateThreadCache,
  clearThreadCache,

  // Performance
  recordRenderTime,
  resetPerformanceMetrics,

  // Error
  addError,
  resolveError,
  clearErrors,
} = messageTypesSlice.actions;

// Selectors
export const selectSelectedMessages = (state: { messageTypes: MessageTypesState }) =>
  Array.from(state.messageTypes.selectedMessages);

export const selectActiveThread = (state: { messageTypes: MessageTypesState }) =>
  state.messageTypes.activeThread;

export const selectSearchState = (state: { messageTypes: MessageTypesState }) => ({
  query: state.messageTypes.searchQuery,
  filters: state.messageTypes.searchFilters,
});

export const selectComposingState = (state: { messageTypes: MessageTypesState }) =>
  state.messageTypes.composing;

export const selectMessageFromCache = (messageId: string) =>
  (state: { messageTypes: MessageTypesState }) =>
    state.messageTypes.messageCache[messageId];

export const selectThreadFromCache = (threadId: string) =>
  (state: { messageTypes: MessageTypesState }) =>
    state.messageTypes.threadCache[threadId];

export const selectInteractionState = (messageId: string) =>
  (state: { messageTypes: MessageTypesState }) =>
    state.messageTypes.interactions[messageId];

export const selectPerformanceMetrics = (state: { messageTypes: MessageTypesState }) =>
  state.messageTypes.performanceMetrics;

export const selectUnresolvedErrors = (state: { messageTypes: MessageTypesState }) =>
  state.messageTypes.errors.filter(error => !error.resolved);

export const selectErrorsByType = (errorType: MessageError['type']) =>
  (state: { messageTypes: MessageTypesState }) =>
    state.messageTypes.errors.filter(error => error.type === errorType && !error.resolved);

// Export reducer
export default messageTypesSlice.reducer;