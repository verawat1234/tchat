/**
 * Simple Toast Hook for Error Messages
 */

import { useState, useCallback } from 'react';

export interface ToastProps {
  title: string;
  description?: string;
  variant?: 'default' | 'destructive';
}

export function useToast() {
  const [toasts, setToasts] = useState<ToastProps[]>([]);

  const toast = useCallback((props: ToastProps) => {
    console.log('Toast:', props);

    // Simple implementation - just log for now
    // In a real app, this would trigger a toast notification
    const toastMessage = `${props.title}${props.description ? ': ' + props.description : ''}`;

    if (props.variant === 'destructive') {
      console.error(toastMessage);
    } else {
      console.info(toastMessage);
    }

    setToasts(prev => [...prev, props]);

    // Auto-remove after 5 seconds
    setTimeout(() => {
      setToasts(prev => prev.slice(1));
    }, 5000);
  }, []);

  return { toast, toasts };
}