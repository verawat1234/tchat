import { describe, it, expect } from 'vitest';
import { configureStore } from '@reduxjs/toolkit';
import contentReducer, {
  setSelectedLanguage,
  updateContentPreferences,
  setSyncStatus,
  toggleFallbackMode,
  updateFallbackContent,
  clearFallbackContent,
  removeFallbackContent,
  batchUpdateFallbackContent,
  selectSelectedLanguage,
  selectContentPreferences,
  selectSyncStatus,
  selectFallbackMode,
  selectFallbackContent,
  selectFallbackContentById,
  selectHasFallbackContent,
  selectFallbackContentCount,
} from './contentSlice';
import type { ContentValue } from '../types/content';

// Helper to create a test store
const createTestStore = (initialState = {}) => {
  return configureStore({
    reducer: {
      content: contentReducer,
    },
    preloadedState: {
      content: {
        selectedLanguage: 'en',
        contentPreferences: {
          showDrafts: false,
          compactView: false,
        },
        lastSyncTime: new Date().toISOString(),
        syncStatus: 'idle' as const,
        fallbackMode: false,
        fallbackContent: {},
        ...initialState,
      },
    },
  });
};

describe('contentSlice', () => {
  describe('initial state', () => {
    it('should have correct initial state', () => {
      const store = createTestStore();
      const state = store.getState().content;

      expect(state.selectedLanguage).toBe('en');
      expect(state.contentPreferences.showDrafts).toBe(false);
      expect(state.contentPreferences.compactView).toBe(false);
      expect(state.syncStatus).toBe('idle');
      expect(state.fallbackMode).toBe(false);
      expect(state.fallbackContent).toEqual({});
    });
  });

  describe('setSelectedLanguage', () => {
    it('should update selected language', () => {
      const store = createTestStore();

      store.dispatch(setSelectedLanguage('es'));

      const state = store.getState().content;
      expect(state.selectedLanguage).toBe('es');
    });
  });

  describe('updateContentPreferences', () => {
    it('should update content preferences partially', () => {
      const store = createTestStore();

      store.dispatch(updateContentPreferences({ showDrafts: true }));

      const state = store.getState().content;
      expect(state.contentPreferences.showDrafts).toBe(true);
      expect(state.contentPreferences.compactView).toBe(false); // Should remain unchanged
    });

    it('should update multiple preferences at once', () => {
      const store = createTestStore();

      store.dispatch(updateContentPreferences({
        showDrafts: true,
        compactView: true
      }));

      const state = store.getState().content;
      expect(state.contentPreferences.showDrafts).toBe(true);
      expect(state.contentPreferences.compactView).toBe(true);
    });
  });

  describe('setSyncStatus', () => {
    it('should update sync status without timestamp', () => {
      const store = createTestStore();

      store.dispatch(setSyncStatus({ status: 'syncing' }));

      const state = store.getState().content;
      expect(state.syncStatus).toBe('syncing');
    });

    it('should update sync status with custom timestamp', () => {
      const store = createTestStore();
      const customTime = '2023-01-01T00:00:00.000Z';

      store.dispatch(setSyncStatus({
        status: 'error',
        timestamp: customTime
      }));

      const state = store.getState().content;
      expect(state.syncStatus).toBe('error');
      expect(state.lastSyncTime).toBe(customTime);
    });

    it('should update lastSyncTime when status becomes idle', () => {
      const store = createTestStore();
      const initialTime = store.getState().content.lastSyncTime;

      // Small delay to ensure different timestamp
      setTimeout(() => {
        store.dispatch(setSyncStatus({ status: 'idle' }));

        const state = store.getState().content;
        expect(state.syncStatus).toBe('idle');
        expect(state.lastSyncTime).not.toBe(initialTime);
      }, 1);
    });
  });

  describe('toggleFallbackMode', () => {
    it('should toggle fallback mode on', () => {
      const store = createTestStore();

      store.dispatch(toggleFallbackMode(true));

      const state = store.getState().content;
      expect(state.fallbackMode).toBe(true);
    });

    it('should toggle fallback mode off', () => {
      const store = createTestStore({ fallbackMode: true });

      store.dispatch(toggleFallbackMode(false));

      const state = store.getState().content;
      expect(state.fallbackMode).toBe(false);
    });
  });

  describe('fallback content management', () => {
    const mockContent: ContentValue = {
      type: 'text',
      value: 'Test content',
    };

    it('should update fallback content', () => {
      const store = createTestStore();

      store.dispatch(updateFallbackContent({
        contentId: 'test.content.key',
        content: mockContent,
      }));

      const state = store.getState().content;
      expect(state.fallbackContent['test.content.key']).toEqual(mockContent);
    });

    it('should clear all fallback content', () => {
      const store = createTestStore({
        fallbackContent: {
          'test.key.1': mockContent,
          'test.key.2': mockContent,
        },
      });

      store.dispatch(clearFallbackContent());

      const state = store.getState().content;
      expect(state.fallbackContent).toEqual({});
    });

    it('should remove specific fallback content', () => {
      const store = createTestStore({
        fallbackContent: {
          'test.key.1': mockContent,
          'test.key.2': mockContent,
        },
      });

      store.dispatch(removeFallbackContent('test.key.1'));

      const state = store.getState().content;
      expect(state.fallbackContent['test.key.1']).toBeUndefined();
      expect(state.fallbackContent['test.key.2']).toEqual(mockContent);
    });

    it('should batch update fallback content', () => {
      const store = createTestStore();
      const batchContent = {
        'test.key.1': mockContent,
        'test.key.2': { type: 'text' as const, value: 'Second content' },
      };

      store.dispatch(batchUpdateFallbackContent(batchContent));

      const state = store.getState().content;
      expect(state.fallbackContent).toEqual(batchContent);
    });
  });

  describe('selectors', () => {
    const store = createTestStore({
      selectedLanguage: 'fr',
      contentPreferences: { showDrafts: true, compactView: true },
      syncStatus: 'syncing' as const,
      fallbackMode: true,
      fallbackContent: {
        'test.key.1': { type: 'text' as const, value: 'Test content' },
        'test.key.2': { type: 'text' as const, value: 'Another content' },
      },
    });

    it('should select selected language', () => {
      const language = selectSelectedLanguage(store.getState());
      expect(language).toBe('fr');
    });

    it('should select content preferences', () => {
      const preferences = selectContentPreferences(store.getState());
      expect(preferences).toEqual({ showDrafts: true, compactView: true });
    });

    it('should select sync status', () => {
      const status = selectSyncStatus(store.getState());
      expect(status).toBe('syncing');
    });

    it('should select fallback mode', () => {
      const fallbackMode = selectFallbackMode(store.getState());
      expect(fallbackMode).toBe(true);
    });

    it('should select fallback content', () => {
      const fallbackContent = selectFallbackContent(store.getState());
      expect(Object.keys(fallbackContent)).toHaveLength(2);
    });

    it('should select fallback content by ID', () => {
      const content = selectFallbackContentById(store.getState(), 'test.key.1');
      expect(content).toEqual({ type: 'text', value: 'Test content' });
    });

    it('should return undefined for non-existent content ID', () => {
      const content = selectFallbackContentById(store.getState(), 'non.existent.key');
      expect(content).toBeUndefined();
    });

    it('should check if fallback content exists', () => {
      const hasContent = selectHasFallbackContent(store.getState(), 'test.key.1');
      const hasNoContent = selectHasFallbackContent(store.getState(), 'non.existent.key');

      expect(hasContent).toBe(true);
      expect(hasNoContent).toBe(false);
    });

    it('should count fallback content items', () => {
      const count = selectFallbackContentCount(store.getState());
      expect(count).toBe(2);
    });
  });
});