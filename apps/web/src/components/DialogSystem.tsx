/**
 * Comprehensive Dialog System
 *
 * Advanced dialog component system with multiple dialog types, animations,
 * and comprehensive functionality for the Tchat application.
 *
 * Features:
 * - Multiple dialog types (modal, sheet, popover, toast)
 * - Rich content support (forms, media, actions)
 * - Animation and transition system
 * - Keyboard and accessibility support
 * - Mobile-responsive design
 * - Dialog stacking and z-index management
 * - Custom dialog templates for common use cases
 */

import React, { useState, useEffect, useCallback, createContext, useContext } from 'react';
import { motion, AnimatePresence, PanInfo } from 'framer-motion';
import { createPortal } from 'react-dom';
import {
  X,
  ChevronDown,
  AlertTriangle,
  CheckCircle,
  Info,
  AlertCircle,
  Trash2,
  Share,
  Edit,
  Settings,
  Users,
  MessageSquare,
  Phone,
  Video,
  Download,
  Upload,
  Calendar,
  Image,
  File,
  ExternalLink,
} from 'lucide-react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Textarea } from './ui/textarea';
import { Badge } from './ui/badge/badge';
import { Avatar } from './ui/avatar';
import { ScrollArea } from './ui/scroll-area';

// =============================================================================
// Type Definitions
// =============================================================================

export type DialogType =
  | 'modal'       // Full modal dialog
  | 'sheet'       // Side sheet (mobile drawer)
  | 'popover'     // Small popover dialog
  | 'fullscreen'  // Fullscreen modal
  | 'bottomSheet' // Bottom sheet (mobile)
  | 'toast'       // Toast notification dialog
  | 'confirmation' // Confirmation dialog
  | 'form'        // Form dialog
  | 'media'       // Media viewer dialog
  | 'picker';     // Selection picker dialog

export type DialogSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | 'full';

export type DialogPosition =
  | 'center'
  | 'top'
  | 'bottom'
  | 'left'
  | 'right'
  | 'topLeft'
  | 'topRight'
  | 'bottomLeft'
  | 'bottomRight';

export interface DialogAction {
  id: string;
  label: string;
  variant: 'primary' | 'secondary' | 'destructive' | 'ghost';
  icon?: React.ComponentType<{ className?: string }>;
  onClick: () => void | Promise<void>;
  disabled?: boolean;
  loading?: boolean;
  shortcut?: string;
}

export interface DialogConfig {
  id: string;
  type: DialogType;
  title?: string;
  description?: string;
  content?: React.ReactNode;
  size?: DialogSize;
  position?: DialogPosition;
  closable?: boolean;
  dismissable?: boolean; // Click outside to close
  persistent?: boolean; // Don't close on escape
  showOverlay?: boolean;
  overlayBlur?: boolean;
  animate?: boolean;
  duration?: number; // Auto-close after duration (ms)
  actions?: DialogAction[];
  onClose?: () => void;
  onOpen?: () => void;
  className?: string;
  zIndex?: number;
}

export interface DialogContextType {
  dialogs: DialogConfig[];
  openDialog: (config: Omit<DialogConfig, 'id'>) => string;
  closeDialog: (id: string) => void;
  closeAllDialogs: () => void;
  updateDialog: (id: string, updates: Partial<DialogConfig>) => void;
  isDialogOpen: (id: string) => boolean;
}

// =============================================================================
// Dialog Context
// =============================================================================

const DialogContext = createContext<DialogContextType | null>(null);

export function useDialog() {
  const context = useContext(DialogContext);
  if (!context) {
    throw new Error('useDialog must be used within a DialogProvider');
  }
  return context;
}

// =============================================================================
// Dialog Provider
// =============================================================================

export function DialogProvider({ children }: { children: React.ReactNode }) {
  const [dialogs, setDialogs] = useState<DialogConfig[]>([]);

  const openDialog = useCallback((config: Omit<DialogConfig, 'id'>) => {
    const id = `dialog_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    const dialogConfig: DialogConfig = {
      ...config,
      id,
      zIndex: 1000 + dialogs.length,
    };

    setDialogs(prev => [...prev, dialogConfig]);

    // Auto-close if duration is set
    if (config.duration) {
      setTimeout(() => {
        closeDialog(id);
      }, config.duration);
    }

    // Call onOpen callback
    config.onOpen?.();

    return id;
  }, [dialogs.length]);

  const closeDialog = useCallback((id: string) => {
    setDialogs(prev => {
      const dialog = prev.find(d => d.id === id);
      if (dialog?.onClose) {
        dialog.onClose();
      }
      return prev.filter(d => d.id !== id);
    });
  }, []);

  const closeAllDialogs = useCallback(() => {
    setDialogs(prev => {
      prev.forEach(dialog => dialog.onClose?.());
      return [];
    });
  }, []);

  const updateDialog = useCallback((id: string, updates: Partial<DialogConfig>) => {
    setDialogs(prev => prev.map(dialog =>
      dialog.id === id ? { ...dialog, ...updates } : dialog
    ));
  }, []);

  const isDialogOpen = useCallback((id: string) => {
    return dialogs.some(d => d.id === id);
  }, [dialogs]);

  // Handle escape key
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        const topDialog = dialogs[dialogs.length - 1];
        if (topDialog && !topDialog.persistent && topDialog.closable !== false) {
          closeDialog(topDialog.id);
        }
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [dialogs, closeDialog]);

  const contextValue: DialogContextType = {
    dialogs,
    openDialog,
    closeDialog,
    closeAllDialogs,
    updateDialog,
    isDialogOpen,
  };

  return (
    <DialogContext.Provider value={contextValue}>
      {children}
      <DialogRenderer />
    </DialogContext.Provider>
  );
}

// =============================================================================
// Dialog Renderer
// =============================================================================

function DialogRenderer() {
  const { dialogs } = useDialog();

  return createPortal(
    <AnimatePresence mode="multiple">
      {dialogs.map(dialog => (
        <DialogComponent key={dialog.id} config={dialog} />
      ))}
    </AnimatePresence>,
    document.body
  );
}

// =============================================================================
// Individual Dialog Component
// =============================================================================

function DialogComponent({ config }: { config: DialogConfig }) {
  const { closeDialog } = useDialog();
  const [isDragging, setIsDragging] = useState(false);

  const handleClose = useCallback(() => {
    if (config.closable !== false) {
      closeDialog(config.id);
    }
  }, [config.id, config.closable, closeDialog]);

  const handleOverlayClick = useCallback(() => {
    if (config.dismissable !== false) {
      handleClose();
    }
  }, [config.dismissable, handleClose]);

  const handleDragEnd = useCallback((event: any, info: PanInfo) => {
    setIsDragging(false);

    // Close dialog if dragged down significantly (for mobile)
    if (config.type === 'bottomSheet' && info.offset.y > 200) {
      handleClose();
    }
  }, [config.type, handleClose]);

  // Get dialog styles based on type and size
  const getDialogStyles = () => {
    const baseClasses = "bg-white rounded-lg shadow-xl border";

    const sizeClasses = {
      xs: "max-w-xs",
      sm: "max-w-sm",
      md: "max-w-md",
      lg: "max-w-lg",
      xl: "max-w-xl",
      full: "w-full h-full rounded-none",
    };

    const size = config.size || 'md';

    switch (config.type) {
      case 'fullscreen':
        return `${baseClasses} w-full h-full rounded-none`;
      case 'sheet':
        return `${baseClasses} h-full w-96 rounded-l-lg rounded-r-none`;
      case 'bottomSheet':
        return `${baseClasses} w-full rounded-t-lg rounded-b-none`;
      case 'popover':
        return `${baseClasses} ${sizeClasses.sm}`;
      case 'toast':
        return `${baseClasses} ${sizeClasses.sm} border-l-4`;
      default:
        return `${baseClasses} ${sizeClasses[size]} w-full max-h-[90vh]`;
    }
  };

  // Get container styles and animations
  const getContainerProps = () => {
    const baseClasses = "fixed inset-0 flex items-center justify-center p-4";

    switch (config.type) {
      case 'sheet':
        return {
          className: "fixed inset-0 flex items-center justify-end",
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          exit: { opacity: 0 },
        };
      case 'bottomSheet':
        return {
          className: "fixed inset-0 flex items-end justify-center",
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          exit: { opacity: 0 },
        };
      case 'toast':
        return {
          className: "fixed top-4 right-4 z-50",
          initial: { opacity: 0, x: 100 },
          animate: { opacity: 1, x: 0 },
          exit: { opacity: 0, x: 100 },
        };
      default:
        return {
          className: baseClasses,
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          exit: { opacity: 0 },
        };
    }
  };

  const getDialogProps = () => {
    const baseProps = {
      className: getDialogStyles(),
      onClick: (e: React.MouseEvent) => e.stopPropagation(),
    };

    switch (config.type) {
      case 'sheet':
        return {
          ...baseProps,
          initial: { x: '100%' },
          animate: { x: 0 },
          exit: { x: '100%' },
          transition: { type: 'spring', damping: 25, stiffness: 200 },
        };
      case 'bottomSheet':
        return {
          ...baseProps,
          initial: { y: '100%' },
          animate: { y: 0 },
          exit: { y: '100%' },
          transition: { type: 'spring', damping: 25, stiffness: 200 },
          drag: 'y',
          dragConstraints: { top: 0 },
          onDragStart: () => setIsDragging(true),
          onDragEnd: handleDragEnd,
        };
      case 'modal':
      case 'fullscreen':
        return {
          ...baseProps,
          initial: { scale: 0.95, opacity: 0 },
          animate: { scale: 1, opacity: 1 },
          exit: { scale: 0.95, opacity: 0 },
          transition: { type: 'spring', damping: 25, stiffness: 200 },
        };
      case 'toast':
        return {
          ...baseProps,
          initial: { scale: 0.8, opacity: 0 },
          animate: { scale: 1, opacity: 1 },
          exit: { scale: 0.8, opacity: 0 },
        };
      default:
        return {
          ...baseProps,
          initial: { scale: 0.95, opacity: 0 },
          animate: { scale: 1, opacity: 1 },
          exit: { scale: 0.95, opacity: 0 },
        };
    }
  };

  const containerProps = getContainerProps();
  const dialogProps = getDialogProps();

  // Render different dialog types
  const renderDialogContent = () => {
    switch (config.type) {
      case 'toast':
        return <ToastDialog config={config} onClose={handleClose} />;
      case 'confirmation':
        return <ConfirmationDialog config={config} onClose={handleClose} />;
      case 'form':
        return <FormDialog config={config} onClose={handleClose} />;
      case 'media':
        return <MediaDialog config={config} onClose={handleClose} />;
      case 'picker':
        return <PickerDialog config={config} onClose={handleClose} />;
      default:
        return <DefaultDialog config={config} onClose={handleClose} />;
    }
  };

  return (
    <motion.div
      {...containerProps}
      style={{ zIndex: config.zIndex }}
    >
      {/* Overlay */}
      {config.showOverlay !== false && config.type !== 'toast' && (
        <motion.div
          className={`absolute inset-0 bg-black/50 ${config.overlayBlur ? 'backdrop-blur-sm' : ''}`}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          onClick={handleOverlayClick}
        />
      )}

      {/* Dialog */}
      <motion.div {...dialogProps}>
        {config.type === 'bottomSheet' && (
          <div className="flex justify-center py-2">
            <div className="w-8 h-1 bg-gray-300 rounded-full" />
          </div>
        )}
        {renderDialogContent()}
      </motion.div>
    </motion.div>
  );
}

// =============================================================================
// Dialog Type Components
// =============================================================================

function DefaultDialog({ config, onClose }: { config: DialogConfig; onClose: () => void }) {
  return (
    <div className="flex flex-col max-h-full">
      {/* Header */}
      {(config.title || config.closable !== false) && (
        <div className="flex items-center justify-between p-6 pb-4 border-b">
          <div className="flex-1">
            {config.title && (
              <h2 className="text-lg font-semibold text-gray-900">{config.title}</h2>
            )}
            {config.description && (
              <p className="text-sm text-gray-600 mt-1">{config.description}</p>
            )}
          </div>
          {config.closable !== false && (
            <Button variant="ghost" size="sm" onClick={onClose} className="ml-4">
              <X className="h-4 w-4" />
            </Button>
          )}
        </div>
      )}

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        {config.content && (
          <div className="p-6">
            {config.content}
          </div>
        )}
      </div>

      {/* Actions */}
      {config.actions && config.actions.length > 0 && (
        <div className="flex items-center justify-end gap-3 p-6 pt-4 border-t bg-gray-50">
          {config.actions.map(action => (
            <Button
              key={action.id}
              variant={action.variant}
              onClick={action.onClick}
              disabled={action.disabled || action.loading}
              className="flex items-center gap-2"
            >
              {action.loading ? (
                <motion.div
                  animate={{ rotate: 360 }}
                  transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                  className="h-4 w-4 border-2 border-current border-t-transparent rounded-full"
                />
              ) : (
                action.icon && <action.icon className="h-4 w-4" />
              )}
              {action.label}
              {action.shortcut && (
                <Badge variant="outline" className="ml-1 text-xs">
                  {action.shortcut}
                </Badge>
              )}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}

function ToastDialog({ config, onClose }: { config: DialogConfig; onClose: () => void }) {
  const getToastIcon = () => {
    // Determine icon based on dialog context or variant
    if (config.title?.toLowerCase().includes('error') || config.description?.toLowerCase().includes('error')) {
      return <AlertTriangle className="h-5 w-5 text-red-500" />;
    }
    if (config.title?.toLowerCase().includes('success') || config.description?.toLowerCase().includes('success')) {
      return <CheckCircle className="h-5 w-5 text-green-500" />;
    }
    if (config.title?.toLowerCase().includes('warning')) {
      return <AlertCircle className="h-5 w-5 text-yellow-500" />;
    }
    return <Info className="h-5 w-5 text-blue-500" />;
  };

  return (
    <div className="flex items-start gap-3 p-4">
      {getToastIcon()}
      <div className="flex-1 min-w-0">
        {config.title && (
          <h3 className="text-sm font-medium text-gray-900">{config.title}</h3>
        )}
        {config.description && (
          <p className="text-sm text-gray-600 mt-1">{config.description}</p>
        )}
        {config.content}
      </div>
      {config.closable !== false && (
        <Button variant="ghost" size="sm" onClick={onClose} className="flex-shrink-0">
          <X className="h-4 w-4" />
        </Button>
      )}
    </div>
  );
}

function ConfirmationDialog({ config, onClose }: { config: DialogConfig; onClose: () => void }) {
  return (
    <div className="p-6">
      <div className="flex items-center gap-4 mb-4">
        <div className="flex-shrink-0">
          <AlertTriangle className="h-8 w-8 text-red-500" />
        </div>
        <div className="flex-1">
          <h3 className="text-lg font-medium text-gray-900">
            {config.title || 'Confirm Action'}
          </h3>
          {config.description && (
            <p className="text-sm text-gray-600 mt-1">{config.description}</p>
          )}
        </div>
      </div>

      {config.content && (
        <div className="mb-6">{config.content}</div>
      )}

      <div className="flex items-center justify-end gap-3">
        <Button variant="outline" onClick={onClose}>
          Cancel
        </Button>
        {config.actions?.map(action => (
          <Button
            key={action.id}
            variant={action.variant}
            onClick={action.onClick}
            disabled={action.disabled || action.loading}
          >
            {action.loading ? (
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                className="h-4 w-4 border-2 border-current border-t-transparent rounded-full mr-2"
              />
            ) : (
              action.icon && <action.icon className="h-4 w-4 mr-2" />
            )}
            {action.label}
          </Button>
        ))}
      </div>
    </div>
  );
}

function FormDialog({ config, onClose }: { config: DialogConfig; onClose: () => void }) {
  return (
    <div className="flex flex-col max-h-full">
      {/* Header */}
      <div className="flex items-center justify-between p-6 pb-4 border-b">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">
            {config.title || 'Form'}
          </h2>
          {config.description && (
            <p className="text-sm text-gray-600 mt-1">{config.description}</p>
          )}
        </div>
        {config.closable !== false && (
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>

      {/* Form Content */}
      <ScrollArea className="flex-1">
        <div className="p-6">
          {config.content}
        </div>
      </ScrollArea>

      {/* Form Actions */}
      {config.actions && (
        <div className="flex items-center justify-end gap-3 p-6 pt-4 border-t bg-gray-50">
          {config.actions.map(action => (
            <Button
              key={action.id}
              variant={action.variant}
              onClick={action.onClick}
              disabled={action.disabled || action.loading}
            >
              {action.loading ? (
                <motion.div
                  animate={{ rotate: 360 }}
                  transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                  className="h-4 w-4 border-2 border-current border-t-transparent rounded-full mr-2"
                />
              ) : (
                action.icon && <action.icon className="h-4 w-4 mr-2" />
              )}
              {action.label}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}

function MediaDialog({ config, onClose }: { config: DialogConfig; onClose: () => void }) {
  return (
    <div className="relative w-full h-full flex flex-col bg-black">
      {/* Header */}
      <div className="absolute top-0 left-0 right-0 z-10 flex items-center justify-between p-4 bg-gradient-to-b from-black/50 to-transparent">
        <div className="text-white">
          {config.title && (
            <h2 className="text-lg font-medium">{config.title}</h2>
          )}
          {config.description && (
            <p className="text-sm opacity-75">{config.description}</p>
          )}
        </div>
        <Button variant="ghost" size="sm" onClick={onClose} className="text-white hover:bg-white/20">
          <X className="h-5 w-5" />
        </Button>
      </div>

      {/* Media Content */}
      <div className="flex-1 flex items-center justify-center p-4">
        {config.content}
      </div>

      {/* Actions */}
      {config.actions && (
        <div className="absolute bottom-0 left-0 right-0 z-10 flex items-center justify-center gap-3 p-4 bg-gradient-to-t from-black/50 to-transparent">
          {config.actions.map(action => (
            <Button
              key={action.id}
              variant="outline"
              onClick={action.onClick}
              disabled={action.disabled || action.loading}
              className="bg-white/10 border-white/20 text-white hover:bg-white/20"
            >
              {action.icon && <action.icon className="h-4 w-4 mr-2" />}
              {action.label}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}

function PickerDialog({ config, onClose }: { config: DialogConfig; onClose: () => void }) {
  return (
    <div className="flex flex-col max-h-full">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b">
        <h3 className="text-lg font-medium text-gray-900">
          {config.title || 'Select Item'}
        </h3>
        {config.closable !== false && (
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>

      {/* Search/Filter */}
      <div className="p-4 border-b bg-gray-50">
        <Input
          placeholder="Search..."
          className="w-full"
        />
      </div>

      {/* Options */}
      <ScrollArea className="flex-1 max-h-96">
        <div className="p-2">
          {config.content}
        </div>
      </ScrollArea>

      {/* Actions */}
      {config.actions && (
        <div className="flex items-center justify-end gap-3 p-4 border-t">
          {config.actions.map(action => (
            <Button
              key={action.id}
              variant={action.variant}
              onClick={action.onClick}
              disabled={action.disabled || action.loading}
            >
              {action.icon && <action.icon className="h-4 w-4 mr-2" />}
              {action.label}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}

// =============================================================================
// Pre-built Dialog Functions
// =============================================================================

export function useDialogHelpers() {
  const { openDialog, closeDialog } = useDialog();

  const showConfirmation = useCallback((
    title: string,
    message: string,
    onConfirm: () => void | Promise<void>,
    options: Partial<DialogConfig> = {}
  ) => {
    return openDialog({
      type: 'confirmation',
      title,
      description: message,
      actions: [
        {
          id: 'confirm',
          label: 'Confirm',
          variant: 'destructive',
          onClick: onConfirm,
        },
      ],
      ...options,
    });
  }, [openDialog]);

  const showToast = useCallback((
    title: string,
    description?: string,
    duration: number = 5000
  ) => {
    return openDialog({
      type: 'toast',
      title,
      description,
      duration,
      closable: true,
      dismissable: true,
    });
  }, [openDialog]);

  const showForm = useCallback((
    title: string,
    content: React.ReactNode,
    actions: DialogAction[],
    options: Partial<DialogConfig> = {}
  ) => {
    return openDialog({
      type: 'form',
      title,
      content,
      actions,
      size: 'lg',
      ...options,
    });
  }, [openDialog]);

  const showMediaViewer = useCallback((
    content: React.ReactNode,
    title?: string,
    actions?: DialogAction[]
  ) => {
    return openDialog({
      type: 'media',
      title,
      content,
      actions,
      size: 'full',
      showOverlay: true,
      overlayBlur: true,
    });
  }, [openDialog]);

  return {
    showConfirmation,
    showToast,
    showForm,
    showMediaViewer,
    openDialog,
    closeDialog,
  };
}

// Export types only (functions are already exported above)
export type {
  DialogConfig,
  DialogAction,
  DialogType,
  DialogSize,
  DialogPosition,
};