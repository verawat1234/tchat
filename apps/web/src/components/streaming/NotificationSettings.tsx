import React, { useEffect, useState } from 'react';

interface NotificationPreferences {
  user_id: string;
  live_start: {
    enabled: boolean;
    channels: string[];
  };
  chat_mention: {
    enabled: boolean;
    channels: string[];
  };
  featured_product: {
    enabled: boolean;
    channels: string[];
  };
  milestone: {
    enabled: boolean;
    channels: string[];
  };
  updated_at: string;
}

export const NotificationSettings: React.FC = () => {
  const [preferences, setPreferences] = useState<NotificationPreferences | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  // Load preferences on mount
  useEffect(() => {
    loadPreferences();
  }, []);

  const loadPreferences = async () => {
    try {
      const response = await fetch('/api/v1/notification-preferences', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('auth_token')}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        setPreferences(data);
      }
    } catch (error) {
      console.error('[Notifications] Failed to load preferences:', error);
    }
  };

  const savePreferences = async () => {
    if (!preferences) return;

    setIsSaving(true);

    try {
      const response = await fetch('/api/v1/notification-preferences', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('auth_token')}`,
        },
        body: JSON.stringify(preferences),
      });

      if (response.ok) {
        const updated = await response.json();
        setPreferences(updated);
        console.log('[Notifications] Preferences saved');
      }
    } catch (error) {
      console.error('[Notifications] Failed to save preferences:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const toggleEvent = (eventType: keyof NotificationPreferences, enabled: boolean) => {
    if (!preferences) return;

    setPreferences({
      ...preferences,
      [eventType]: {
        ...preferences[eventType],
        enabled,
      },
    });
  };

  const toggleChannel = (
    eventType: keyof NotificationPreferences,
    channel: string,
    enabled: boolean
  ) => {
    if (!preferences) return;

    const currentEvent = preferences[eventType] as { enabled: boolean; channels: string[] };
    const currentChannels = currentEvent.channels || [];

    const updatedChannels = enabled
      ? [...currentChannels, channel]
      : currentChannels.filter((c) => c !== channel);

    setPreferences({
      ...preferences,
      [eventType]: {
        ...currentEvent,
        channels: updatedChannels,
      },
    });
  };

  if (!preferences) {
    return (
      <div className="notification-settings p-6">
        <div className="text-center text-gray-500">Loading preferences...</div>
      </div>
    );
  }

  return (
    <div className="notification-settings max-w-2xl mx-auto p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Notification Preferences</h2>

      <div className="space-y-6">
        {/* Live Stream Start */}
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="text-lg font-semibold text-gray-900">Live Stream Start</h3>
              <p className="text-sm text-gray-500 mt-1">
                Get notified when someone you follow goes live
              </p>
            </div>
            <input
              type="checkbox"
              checked={preferences.live_start.enabled}
              onChange={(e) => toggleEvent('live_start', e.target.checked)}
              className="w-5 h-5 text-blue-600 rounded"
            />
          </div>

          {preferences.live_start.enabled && (
            <div className="space-y-2 ml-4">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={preferences.live_start.channels.includes('push')}
                  onChange={(e) => toggleChannel('live_start', 'push', e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">Push notifications</span>
              </label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={preferences.live_start.channels.includes('email')}
                  onChange={(e) => toggleChannel('live_start', 'email', e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">Email</span>
              </label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={preferences.live_start.channels.includes('in_app')}
                  onChange={(e) => toggleChannel('live_start', 'in_app', e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">In-app notifications</span>
              </label>
            </div>
          )}
        </div>

        {/* Chat Mention */}
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="text-lg font-semibold text-gray-900">Chat Mention</h3>
              <p className="text-sm text-gray-500 mt-1">
                Get notified when someone mentions you in chat
              </p>
            </div>
            <input
              type="checkbox"
              checked={preferences.chat_mention.enabled}
              onChange={(e) => toggleEvent('chat_mention', e.target.checked)}
              className="w-5 h-5 text-blue-600 rounded"
            />
          </div>

          {preferences.chat_mention.enabled && (
            <div className="space-y-2 ml-4">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={preferences.chat_mention.channels.includes('push')}
                  onChange={(e) => toggleChannel('chat_mention', 'push', e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">Push notifications</span>
              </label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={preferences.chat_mention.channels.includes('in_app')}
                  onChange={(e) => toggleChannel('chat_mention', 'in_app', e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">In-app notifications</span>
              </label>
            </div>
          )}
        </div>

        {/* Featured Product */}
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="text-lg font-semibold text-gray-900">Featured Product</h3>
              <p className="text-sm text-gray-500 mt-1">
                Get notified about special deals during live streams
              </p>
            </div>
            <input
              type="checkbox"
              checked={preferences.featured_product.enabled}
              onChange={(e) => toggleEvent('featured_product', e.target.checked)}
              className="w-5 h-5 text-blue-600 rounded"
            />
          </div>

          {preferences.featured_product.enabled && (
            <div className="space-y-2 ml-4">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={preferences.featured_product.channels.includes('push')}
                  onChange={(e) => toggleChannel('featured_product', 'push', e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">Push notifications</span>
              </label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={preferences.featured_product.channels.includes('in_app')}
                  onChange={(e) => toggleChannel('featured_product', 'in_app', e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">In-app notifications</span>
              </label>
            </div>
          )}
        </div>

        {/* Milestone */}
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="text-lg font-semibold text-gray-900">Stream Milestones</h3>
              <p className="text-sm text-gray-500 mt-1">
                Get notified about viewer milestones (1K, 10K, etc.)
              </p>
            </div>
            <input
              type="checkbox"
              checked={preferences.milestone.enabled}
              onChange={(e) => toggleEvent('milestone', e.target.checked)}
              className="w-5 h-5 text-blue-600 rounded"
            />
          </div>

          {preferences.milestone.enabled && (
            <div className="space-y-2 ml-4">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={preferences.milestone.channels.includes('in_app')}
                  onChange={(e) => toggleChannel('milestone', 'in_app', e.target.checked)}
                  className="w-4 h-4 text-blue-600 rounded mr-2"
                />
                <span className="text-sm text-gray-700">In-app notifications</span>
              </label>
            </div>
          )}
        </div>
      </div>

      {/* Save Button */}
      <div className="mt-8 flex justify-end">
        <button
          onClick={savePreferences}
          disabled={isSaving}
          className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium"
        >
          {isSaving ? 'Saving...' : 'Save Preferences'}
        </button>
      </div>
    </div>
  );
};