/**
 * Message Status Indicator Component
 *
 * Visual indicators for message delivery status with animations and tooltips.
 * Shows pending, sent, delivered, read, and failed states with appropriate icons.
 */

import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { DeliveryStatus } from '../types/MessageTypes';

interface MessageStatusIndicatorProps {
  status: DeliveryStatus;
  timestamp: string;
  showTooltip?: boolean;
  size?: 'sm' | 'md' | 'lg';
}

export function MessageStatusIndicator({
  status,
  timestamp,
  showTooltip = true,
  size = 'sm'
}: MessageStatusIndicatorProps) {
  const sizeClasses = {
    sm: 'w-3 h-3',
    md: 'w-4 h-4',
    lg: 'w-5 h-5'
  };

  const getStatusIcon = () => {
    switch (status) {
      case DeliveryStatus.PENDING:
        return (
          <motion.div
            className={`${sizeClasses[size]} border-2 border-gray-300 border-t-blue-500 rounded-full`}
            animate={{ rotate: 360 }}
            transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
          />
        );

      case DeliveryStatus.SENT:
        return (
          <motion.svg
            className={`${sizeClasses[size]} text-gray-400`}
            fill="currentColor"
            viewBox="0 0 20 20"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: "spring", stiffness: 300, damping: 20 }}
          >
            <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
          </motion.svg>
        );

      case DeliveryStatus.DELIVERED:
        return (
          <motion.div
            className={`${sizeClasses[size]} relative`}
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: "spring", stiffness: 300, damping: 20 }}
          >
            <svg
              className="text-blue-500 absolute inset-0"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
            </svg>
            <svg
              className="text-blue-500 absolute inset-0 transform translate-x-1"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
            </svg>
          </motion.div>
        );

      case DeliveryStatus.READ:
        return (
          <motion.div
            className={`${sizeClasses[size]} relative`}
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: "spring", stiffness: 300, damping: 20 }}
          >
            <svg
              className="text-green-500 absolute inset-0"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
            </svg>
            <svg
              className="text-green-500 absolute inset-0 transform translate-x-1"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
            </svg>
          </motion.div>
        );

      case DeliveryStatus.FAILED:
        return (
          <motion.svg
            className={`${sizeClasses[size]} text-red-500`}
            fill="currentColor"
            viewBox="0 0 20 20"
            initial={{ scale: 0, rotate: -90 }}
            animate={{ scale: 1, rotate: 0 }}
            transition={{ type: "spring", stiffness: 300, damping: 20 }}
          >
            <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
          </motion.svg>
        );

      default:
        return null;
    }
  };

  const getStatusText = () => {
    switch (status) {
      case DeliveryStatus.PENDING:
        return 'Sending...';
      case DeliveryStatus.SENT:
        return `Sent at ${new Date(timestamp).toLocaleTimeString()}`;
      case DeliveryStatus.DELIVERED:
        return `Delivered at ${new Date(timestamp).toLocaleTimeString()}`;
      case DeliveryStatus.READ:
        return `Read at ${new Date(timestamp).toLocaleTimeString()}`;
      case DeliveryStatus.FAILED:
        return 'Failed to send. Tap to retry.';
      default:
        return '';
    }
  };

  const getStatusColor = () => {
    switch (status) {
      case DeliveryStatus.PENDING:
        return 'text-gray-400';
      case DeliveryStatus.SENT:
        return 'text-gray-500';
      case DeliveryStatus.DELIVERED:
        return 'text-blue-500';
      case DeliveryStatus.READ:
        return 'text-green-500';
      case DeliveryStatus.FAILED:
        return 'text-red-500';
      default:
        return 'text-gray-400';
    }
  };

  return (
    <div className="relative">
      <AnimatePresence mode="wait">
        <motion.div
          key={status}
          className={`flex items-center ${getStatusColor()}`}
          initial={{ opacity: 0, x: 10 }}
          animate={{ opacity: 1, x: 0 }}
          exit={{ opacity: 0, x: -10 }}
          transition={{ duration: 0.2 }}
        >
          {getStatusIcon()}
        </motion.div>
      </AnimatePresence>

      {showTooltip && (
        <motion.div
          className="absolute bottom-full right-0 mb-2 px-2 py-1 bg-gray-900 text-white text-xs rounded shadow-lg whitespace-nowrap opacity-0 hover:opacity-100 transition-opacity z-10"
          initial={{ opacity: 0, y: 5 }}
          whileHover={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2 }}
        >
          {getStatusText()}
          <div className="absolute top-full right-2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-900"></div>
        </motion.div>
      )}
    </div>
  );
}

/**
 * Bulk Message Status Component
 *
 * Shows delivery status for multiple messages with aggregate counts
 */
interface BulkMessageStatusProps {
  statuses: DeliveryStatus[];
  showCounts?: boolean;
}

export function BulkMessageStatus({ statuses, showCounts = true }: BulkMessageStatusProps) {
  const statusCounts = statuses.reduce((acc, status) => {
    acc[status] = (acc[status] || 0) + 1;
    return acc;
  }, {} as Record<DeliveryStatus, number>);

  const totalMessages = statuses.length;
  const deliveredCount = (statusCounts[DeliveryStatus.DELIVERED] || 0) + (statusCounts[DeliveryStatus.READ] || 0);
  const readCount = statusCounts[DeliveryStatus.READ] || 0;
  const failedCount = statusCounts[DeliveryStatus.FAILED] || 0;

  return (
    <div className="flex items-center space-x-2 text-xs text-gray-500">
      {showCounts && (
        <div className="flex space-x-3">
          <span className="flex items-center space-x-1">
            <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
            <span>{deliveredCount}/{totalMessages} delivered</span>
          </span>

          {readCount > 0 && (
            <span className="flex items-center space-x-1">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span>{readCount}/{totalMessages} read</span>
            </span>
          )}

          {failedCount > 0 && (
            <span className="flex items-center space-x-1">
              <div className="w-2 h-2 bg-red-500 rounded-full"></div>
              <span>{failedCount} failed</span>
            </span>
          )}
        </div>
      )}
    </div>
  );
}

/**
 * Read Receipt Avatars Component
 *
 * Shows avatar stack of users who have read the message
 */
interface ReadReceiptAvatarsProps {
  readByUsers: Array<{
    id: string;
    name: string;
    avatar?: string;
    readAt: string;
  }>;
  maxShow?: number;
}

export function ReadReceiptAvatars({ readByUsers, maxShow = 3 }: ReadReceiptAvatarsProps) {
  const visibleUsers = readByUsers.slice(0, maxShow);
  const remainingCount = readByUsers.length - maxShow;

  if (readByUsers.length === 0) {
    return null;
  }

  return (
    <div className="flex items-center space-x-1">
      <div className="flex -space-x-2">
        {visibleUsers.map((user, index) => (
          <motion.div
            key={user.id}
            className="relative"
            initial={{ scale: 0, x: 20 }}
            animate={{ scale: 1, x: 0 }}
            transition={{ delay: index * 0.1, type: "spring", stiffness: 300 }}
          >
            {user.avatar ? (
              <img
                src={user.avatar}
                alt={user.name}
                className="w-4 h-4 rounded-full border border-white"
                title={`${user.name} read at ${new Date(user.readAt).toLocaleTimeString()}`}
              />
            ) : (
              <div
                className="w-4 h-4 rounded-full border border-white bg-blue-500 flex items-center justify-center text-white text-xs font-medium"
                title={`${user.name} read at ${new Date(user.readAt).toLocaleTimeString()}`}
              >
                {user.name.charAt(0).toUpperCase()}
              </div>
            )}
          </motion.div>
        ))}

        {remainingCount > 0 && (
          <motion.div
            className="w-4 h-4 rounded-full border border-white bg-gray-500 flex items-center justify-center text-white text-xs font-medium"
            initial={{ scale: 0, x: 20 }}
            animate={{ scale: 1, x: 0 }}
            transition={{ delay: visibleUsers.length * 0.1, type: "spring", stiffness: 300 }}
            title={`+${remainingCount} more`}
          >
            +{remainingCount}
          </motion.div>
        )}
      </div>

      <span className="text-xs text-gray-500 ml-2">
        Read by {readByUsers.length}
      </span>
    </div>
  );
}