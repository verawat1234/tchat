import { describe, it, expect, beforeEach, vi } from 'vitest';
import uiReducer, {
  setTheme,
  toggleSidebar,
  setSidebarOpen,
  openModal,
  closeModal,
  addNotification,
  removeNotification,
  clearNotifications,
  setLoading,
} from '../../src/features/uiSlice';

describe('UI Slice', () => {
  const initialState = {
    theme: 'system' as const,
    sidebarOpen: true,
    activeModal: null,
    notifications: [],
    isLoading: false,
    loadingMessage: null,
  };

  beforeEach(() => {
    // Reset system time for consistent testing
    vi.setSystemTime(new Date('2023-01-01T00:00:00.000Z'));
  });

  describe('Initial State', () => {
    it('should return the initial state', () => {
      expect(uiReducer(undefined, { type: 'unknown' })).toEqual(initialState);
    });
  });

  describe('Theme Actions', () => {
    it('should set theme to light', () => {
      const action = setTheme('light');
      const state = uiReducer(initialState, action);

      expect(state.theme).toBe('light');
    });

    it('should set theme to dark', () => {
      const action = setTheme('dark');
      const state = uiReducer(initialState, action);

      expect(state.theme).toBe('dark');
    });

    it('should set theme to system', () => {
      const stateWithDarkTheme = { ...initialState, theme: 'dark' as const };
      const action = setTheme('system');
      const state = uiReducer(stateWithDarkTheme, action);

      expect(state.theme).toBe('system');
    });
  });

  describe('Sidebar Actions', () => {
    it('should toggle sidebar from open to closed', () => {
      const action = toggleSidebar();
      const state = uiReducer(initialState, action);

      expect(state.sidebarOpen).toBe(false);
    });

    it('should toggle sidebar from closed to open', () => {
      const closedSidebarState = { ...initialState, sidebarOpen: false };
      const action = toggleSidebar();
      const state = uiReducer(closedSidebarState, action);

      expect(state.sidebarOpen).toBe(true);
    });

    it('should set sidebar open explicitly', () => {
      const closedSidebarState = { ...initialState, sidebarOpen: false };
      const action = setSidebarOpen(true);
      const state = uiReducer(closedSidebarState, action);

      expect(state.sidebarOpen).toBe(true);
    });

    it('should set sidebar closed explicitly', () => {
      const action = setSidebarOpen(false);
      const state = uiReducer(initialState, action);

      expect(state.sidebarOpen).toBe(false);
    });
  });

  describe('Modal Actions', () => {
    it('should open a modal', () => {
      const action = openModal('settings');
      const state = uiReducer(initialState, action);

      expect(state.activeModal).toBe('settings');
    });

    it('should replace active modal', () => {
      const stateWithModal = { ...initialState, activeModal: 'settings' };
      const action = openModal('profile');
      const state = uiReducer(stateWithModal, action);

      expect(state.activeModal).toBe('profile');
    });

    it('should close modal', () => {
      const stateWithModal = { ...initialState, activeModal: 'settings' };
      const action = closeModal();
      const state = uiReducer(stateWithModal, action);

      expect(state.activeModal).toBe(null);
    });

    it('should handle close modal when no modal is open', () => {
      const action = closeModal();
      const state = uiReducer(initialState, action);

      expect(state.activeModal).toBe(null);
    });
  });

  describe('Notification Actions', () => {
    beforeEach(() => {
      // Mock Math.random for consistent notification IDs
      vi.spyOn(Math, 'random').mockReturnValue(0.123456789);
    });

    it('should add a notification', () => {
      const notification = {
        type: 'success' as const,
        message: 'Operation completed',
        duration: 5000,
      };

      const action = addNotification(notification);
      const state = uiReducer(initialState, action);

      expect(state.notifications).toHaveLength(1);
      expect(state.notifications[0]).toEqual({
        ...notification,
        id: 'notification-0-0.123456789',
        timestamp: Date.now(),
      });
    });

    it('should add multiple notifications', () => {
      let state = uiReducer(initialState, addNotification({
        type: 'info',
        message: 'First notification',
      }));

      // Advance time and change random for second notification
      vi.setSystemTime(new Date('2023-01-01T00:01:00.000Z'));
      vi.spyOn(Math, 'random').mockReturnValue(0.987654321);

      state = uiReducer(state, addNotification({
        type: 'warning',
        message: 'Second notification',
      }));

      expect(state.notifications).toHaveLength(2);
      expect(state.notifications[0].message).toBe('First notification');
      expect(state.notifications[1].message).toBe('Second notification');
      expect(state.notifications[1].id).toBe('notification-60000-0.987654321');
    });

    it('should remove a notification by id', () => {
      const stateWithNotifications = {
        ...initialState,
        notifications: [
          {
            id: 'notification-1',
            type: 'info' as const,
            message: 'First',
            timestamp: Date.now(),
          },
          {
            id: 'notification-2',
            type: 'success' as const,
            message: 'Second',
            timestamp: Date.now(),
          },
        ],
      };

      const action = removeNotification('notification-1');
      const state = uiReducer(stateWithNotifications, action);

      expect(state.notifications).toHaveLength(1);
      expect(state.notifications[0].id).toBe('notification-2');
    });

    it('should handle remove non-existent notification', () => {
      const stateWithNotifications = {
        ...initialState,
        notifications: [
          {
            id: 'notification-1',
            type: 'info' as const,
            message: 'First',
            timestamp: Date.now(),
          },
        ],
      };

      const action = removeNotification('non-existent');
      const state = uiReducer(stateWithNotifications, action);

      expect(state.notifications).toHaveLength(1);
      expect(state.notifications[0].id).toBe('notification-1');
    });

    it('should clear all notifications', () => {
      const stateWithNotifications = {
        ...initialState,
        notifications: [
          {
            id: 'notification-1',
            type: 'info' as const,
            message: 'First',
            timestamp: Date.now(),
          },
          {
            id: 'notification-2',
            type: 'error' as const,
            message: 'Second',
            timestamp: Date.now(),
          },
        ],
      };

      const action = clearNotifications();
      const state = uiReducer(stateWithNotifications, action);

      expect(state.notifications).toHaveLength(0);
    });

    it('should add notification with default duration', () => {
      const notification = {
        type: 'error' as const,
        message: 'Error occurred',
      };

      const action = addNotification(notification);
      const state = uiReducer(initialState, action);

      expect(state.notifications[0]).toEqual({
        ...notification,
        id: 'notification-0-0.123456789',
        timestamp: Date.now(),
      });
      expect(state.notifications[0].duration).toBeUndefined();
    });
  });

  describe('Loading Actions', () => {
    it('should set loading state', () => {
      const action = setLoading({ isLoading: true, message: 'Processing...' });
      const state = uiReducer(initialState, action);

      expect(state.isLoading).toBe(true);
      expect(state.loadingMessage).toBe('Processing...');
    });

    it('should clear loading state', () => {
      const loadingState = {
        ...initialState,
        isLoading: true,
        loadingMessage: 'Loading...',
      };

      const action = setLoading({ isLoading: false });
      const state = uiReducer(loadingState, action);

      expect(state.isLoading).toBe(false);
      expect(state.loadingMessage).toBe(null);
    });

    it('should set loading without message', () => {
      const action = setLoading({ isLoading: true });
      const state = uiReducer(initialState, action);

      expect(state.isLoading).toBe(true);
      expect(state.loadingMessage).toBe(null);
    });

    it('should update loading message while loading', () => {
      const loadingState = {
        ...initialState,
        isLoading: true,
        loadingMessage: 'Initial message',
      };

      const action = setLoading({ isLoading: true, message: 'Updated message' });
      const state = uiReducer(loadingState, action);

      expect(state.isLoading).toBe(true);
      expect(state.loadingMessage).toBe('Updated message');
    });
  });

  describe('Complex State Interactions', () => {
    it('should handle multiple simultaneous UI changes', () => {
      let state = initialState;

      // Set theme to dark
      state = uiReducer(state, setTheme('dark'));

      // Close sidebar
      state = uiReducer(state, setSidebarOpen(false));

      // Open modal
      state = uiReducer(state, openModal('settings'));

      // Add notification
      state = uiReducer(state, addNotification({
        type: 'info',
        message: 'Settings opened',
      }));

      // Set loading
      state = uiReducer(state, setLoading({ isLoading: true, message: 'Saving...' }));

      expect(state).toEqual({
        theme: 'dark',
        sidebarOpen: false,
        activeModal: 'settings',
        notifications: [
          {
            id: 'notification-0-0.123456789',
            type: 'info',
            message: 'Settings opened',
            timestamp: Date.now(),
          },
        ],
        isLoading: true,
        loadingMessage: 'Saving...',
      });
    });

    it('should maintain independent state properties', () => {
      let state = uiReducer(initialState, setTheme('dark'));
      state = uiReducer(state, openModal('profile'));

      // Changing theme shouldn't affect modal
      state = uiReducer(state, setTheme('light'));
      expect(state.activeModal).toBe('profile');

      // Closing modal shouldn't affect theme
      state = uiReducer(state, closeModal());
      expect(state.theme).toBe('light');
    });
  });

  describe('Edge Cases', () => {
    it('should handle notification with very long message', () => {
      const longMessage = 'A'.repeat(1000);
      const action = addNotification({
        type: 'info',
        message: longMessage,
      });

      const state = uiReducer(initialState, action);
      expect(state.notifications[0].message).toBe(longMessage);
    });

    it('should handle notification with zero duration', () => {
      const action = addNotification({
        type: 'info',
        message: 'Test',
        duration: 0,
      });

      const state = uiReducer(initialState, action);
      expect(state.notifications[0].duration).toBe(0);
    });

    it('should handle empty notification message', () => {
      const action = addNotification({
        type: 'info',
        message: '',
      });

      const state = uiReducer(initialState, action);
      expect(state.notifications[0].message).toBe('');
    });
  });
});