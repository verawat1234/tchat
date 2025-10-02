import React, { useEffect, useRef, useState } from 'react';
import { SignalingClient } from '../../services/streaming/signalingClient';

interface ChatMessage {
  message_id: string;
  stream_id: string;
  user_id: string;
  message: string;
  timestamp: string;
  moderation_status: string;
}

interface ChatPanelProps {
  streamId: string;
  signalingClient: SignalingClient;
  isBroadcaster: boolean;
  currentUserId: string;
}

export const ChatPanel: React.FC<ChatPanelProps> = ({
  streamId,
  signalingClient,
  isBroadcaster,
  currentUserId,
}) => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputMessage, setInputMessage] = useState('');
  const [isRateLimited, setIsRateLimited] = useState(false);
  const [rateLimitCountdown, setRateLimitCountdown] = useState(0);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Listen for incoming chat messages
  useEffect(() => {
    const unsubscribe = signalingClient.on('CHAT', (payload: unknown) => {
      const chatMessage = payload as ChatMessage;
      setMessages((prev) => [...prev, chatMessage]);
    });

    return () => {
      unsubscribe();
    };
  }, [signalingClient]);

  // Handle send message
  const handleSendMessage = async () => {
    if (!inputMessage.trim() || isRateLimited) return;

    // Send via signaling client
    signalingClient.sendChat(inputMessage);

    // Clear input
    setInputMessage('');

    // Apply rate limiting (5 messages per second = 200ms cooldown)
    setIsRateLimited(true);
    setRateLimitCountdown(200);

    const interval = setInterval(() => {
      setRateLimitCountdown((prev) => {
        if (prev <= 10) {
          setIsRateLimited(false);
          clearInterval(interval);
          return 0;
        }
        return prev - 10;
      });
    }, 10);
  };

  // Handle delete message (broadcaster only)
  const handleDeleteMessage = async (messageId: string) => {
    try {
      const response = await fetch(`/api/v1/streams/${streamId}/chat/${messageId}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('auth_token')}`,
        },
      });

      if (response.ok) {
        // Remove message from UI
        setMessages((prev) => prev.filter((msg) => msg.message_id !== messageId));
      }
    } catch (error) {
      console.error('[Chat] Failed to delete message:', error);
    }
  };

  // Handle enter key
  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <div className="chat-panel flex flex-col h-full bg-gray-50 rounded-lg border border-gray-200">
      {/* Chat Header */}
      <div className="flex-shrink-0 px-4 py-3 border-b border-gray-200 bg-white rounded-t-lg">
        <h3 className="text-lg font-semibold text-gray-900">Live Chat</h3>
        <p className="text-xs text-gray-500 mt-1">{messages.length} messages</p>
      </div>

      {/* Messages List */}
      <div className="flex-1 overflow-y-auto px-4 py-3 space-y-3">
        {messages.length === 0 ? (
          <div className="text-center text-gray-500 text-sm mt-8">
            No messages yet. Be the first to chat!
          </div>
        ) : (
          messages.map((message) => (
            <div
              key={message.message_id}
              className="flex items-start space-x-2 group"
            >
              {/* Avatar Placeholder */}
              <div className="flex-shrink-0 w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
                U
              </div>

              {/* Message Content */}
              <div className="flex-1 min-w-0">
                <div className="flex items-baseline space-x-2">
                  <span className="text-sm font-medium text-gray-900">
                    User {message.user_id.slice(0, 8)}
                  </span>
                  <span className="text-xs text-gray-500">
                    {new Date(message.timestamp).toLocaleTimeString()}
                  </span>
                </div>
                <p className="text-sm text-gray-700 mt-0.5 break-words">
                  {message.message}
                </p>
              </div>

              {/* Delete Button (broadcaster only) */}
              {isBroadcaster && (
                <button
                  onClick={() => handleDeleteMessage(message.message_id)}
                  className="flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity text-red-500 hover:text-red-700"
                  title="Delete message"
                >
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              )}
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Message Input */}
      <div className="flex-shrink-0 px-4 py-3 border-t border-gray-200 bg-white rounded-b-lg">
        <div className="flex items-end space-x-2">
          <input
            type="text"
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Type a message..."
            disabled={isRateLimited}
            maxLength={500}
            className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100 disabled:cursor-not-allowed"
          />
          <button
            onClick={handleSendMessage}
            disabled={!inputMessage.trim() || isRateLimited}
            className="flex-shrink-0 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
          >
            Send
          </button>
        </div>

        {/* Rate Limit Indicator */}
        {isRateLimited && (
          <div className="mt-2 text-xs text-amber-600 flex items-center">
            <svg
              className="w-4 h-4 mr-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            Please wait {Math.ceil(rateLimitCountdown / 100)}ms before sending another message
          </div>
        )}

        {/* Character Count */}
        <div className="mt-1 text-xs text-gray-500 text-right">
          {inputMessage.length}/500
        </div>
      </div>
    </div>
  );
};