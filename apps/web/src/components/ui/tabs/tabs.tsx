/**
 * Tabs Components
 * Accessible tab navigation with content panels
 */

import React, { createContext, useContext, useState, useEffect, useId } from 'react';
import { cn } from '@/utils/cn';
import type {
  TabsProps,
  TabsListProps,
  TabsTriggerProps,
  TabsContentProps,
  TabsContextValue
} from '../../../../specs/001-agent-frontend-specialist/contracts/tabs';

/**
 * Tabs context for sharing state between components
 */
const TabsContext = createContext<TabsContextValue | undefined>(undefined);

const useTabsContext = () => {
  const context = useContext(TabsContext);
  if (!context) {
    throw new Error('Tabs components must be used within a Tabs provider');
  }
  return context;
};

/**
 * Root tabs container component
 */
export const Tabs = React.forwardRef<HTMLDivElement, TabsProps>(
  ({
    className,
    children,
    testId,
    defaultValue,
    value: controlledValue,
    onValueChange,
    orientation = 'horizontal',
    activationMode = 'automatic',
    loop = false,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedBy,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const [internalValue, setInternalValue] = useState(defaultValue || '');
    const isControlled = controlledValue !== undefined;
    const value = isControlled ? controlledValue : internalValue;

    const handleValueChange = (newValue: string) => {
      if (!isControlled) {
        setInternalValue(newValue);
      }
      onValueChange?.(newValue);
    };

    const contextValue: TabsContextValue = {
      value,
      onValueChange: handleValueChange,
      orientation,
      activationMode,
      loop
    };

    return (
      <TabsContext.Provider value={contextValue}>
        <div
          ref={ref}
          data-testid={testId}
          className={cn(
            // Base styles
            orientation === 'horizontal' ? 'flex flex-col' : 'flex flex-row',

            // Custom className
            className
          )}
          aria-label={ariaLabel}
          aria-describedby={ariaDescribedBy}
          aria-expanded={ariaExpanded}
          aria-disabled={ariaDisabled}
          role={role}
          tabIndex={tabIndex}
          {...props}
        >
          {children}
        </div>
      </TabsContext.Provider>
    );
  }
);

Tabs.displayName = 'Tabs';

/**
 * Tabs list container component
 */
export const TabsList = React.forwardRef<HTMLDivElement, TabsListProps>(
  ({
    className,
    children,
    testId,
    loop: listLoop,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedBy,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const { orientation, activationMode, loop: contextLoop } = useTabsContext();
    const effectiveLoop = listLoop ?? contextLoop ?? false;

    const handleKeyDown = (event: React.KeyboardEvent) => {
      const triggers = Array.from(
        (event.currentTarget as HTMLElement).querySelectorAll('[role="tab"]:not([disabled])')
      ) as HTMLElement[];

      if (triggers.length === 0) return;

      const currentIndex = triggers.findIndex(trigger => trigger === event.target);
      if (currentIndex === -1) return;

      let nextIndex = currentIndex;

      switch (event.key) {
        case 'ArrowRight':
        case 'ArrowDown':
          event.preventDefault();
          nextIndex = currentIndex + 1;
          if (nextIndex >= triggers.length) {
            nextIndex = effectiveLoop ? 0 : triggers.length - 1;
          }
          break;

        case 'ArrowLeft':
        case 'ArrowUp':
          event.preventDefault();
          nextIndex = currentIndex - 1;
          if (nextIndex < 0) {
            nextIndex = effectiveLoop ? triggers.length - 1 : 0;
          }
          break;

        case 'Home':
          event.preventDefault();
          nextIndex = 0;
          break;

        case 'End':
          event.preventDefault();
          nextIndex = triggers.length - 1;
          break;

        default:
          return;
      }

      const nextTrigger = triggers[nextIndex];
      if (nextTrigger) {
        nextTrigger.focus();
        if (activationMode === 'automatic') {
          nextTrigger.click();
        }
      }
    };

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex border-b border-gray-200',

          // Orientation styles
          orientation === 'horizontal' ? 'flex-row' : 'flex-col border-r border-b-0',

          // Custom className
          className
        )}
        role={role || 'tablist'}
        aria-orientation={orientation}
        aria-label={ariaLabel}
        aria-describedby={ariaDescribedBy}
        aria-expanded={ariaExpanded}
        aria-disabled={ariaDisabled}
        tabIndex={tabIndex}
        data-activation-mode={activationMode}
        onKeyDown={handleKeyDown}
        {...props}
      >
        {children}
      </div>
    );
  }
);

TabsList.displayName = 'TabsList';

/**
 * Individual tab trigger component
 */
export const TabsTrigger = React.forwardRef<HTMLButtonElement, TabsTriggerProps>(
  ({
    className,
    children,
    testId,
    value: triggerValue,
    disabled = false,
    icon,
    badge,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedBy,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const { value: activeValue, onValueChange, orientation } = useTabsContext();
    const isSelected = activeValue === triggerValue;
    const triggerId = useId();
    const panelId = `panel-${triggerValue}`;

    const handleClick = () => {
      if (!disabled) {
        onValueChange?.(triggerValue);
      }
    };

    const handleKeyDown = (event: React.KeyboardEvent) => {
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault();
        handleClick();
      }
    };

    return (
      <button
        ref={ref}
        type="button"
        data-testid={testId}
        className={cn(
          // Base styles
          'px-4 py-2 text-sm font-medium transition-colors duration-200',
          'border-b-2 border-transparent hover:text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500',

          // Orientation styles
          orientation === 'horizontal'
            ? 'border-b-2'
            : 'border-r-2 border-b-0',

          // Selected state
          isSelected && [
            'text-blue-600',
            orientation === 'horizontal'
              ? 'border-b-blue-600'
              : 'border-r-blue-600'
          ],

          // Disabled state
          disabled && 'opacity-50 cursor-not-allowed hover:text-current',

          // Custom className
          className
        )}
        role={role || 'tab'}
        id={triggerId}
        aria-selected={isSelected}
        aria-controls={panelId}
        aria-disabled={ariaDisabled || disabled}
        aria-label={ariaLabel}
        aria-describedby={ariaDescribedBy}
        aria-expanded={ariaExpanded}
        tabIndex={isSelected ? (tabIndex ?? 0) : -1}
        disabled={disabled}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
        {...props}
      >
        <div className="flex items-center space-x-2">
          {/* Icon */}
          {icon && (
            <span className="flex items-center">
              {icon}
            </span>
          )}

          {/* Label */}
          <span>{children}</span>

          {/* Badge */}
          {badge && (
            <span className="px-1.5 py-0.5 text-xs bg-gray-100 text-gray-600 rounded-full">
              {badge}
            </span>
          )}
        </div>
      </button>
    );
  }
);

TabsTrigger.displayName = 'TabsTrigger';

/**
 * Tab content panel component
 */
export const TabsContent = React.forwardRef<HTMLDivElement, TabsContentProps>(
  ({
    className,
    children,
    testId,
    value: contentValue,
    forceMount = false,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedBy,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const { value: activeValue } = useTabsContext();
    const isSelected = activeValue === contentValue;
    const triggerId = `trigger-${contentValue}`;
    const panelId = useId();

    // Don't render if not selected and not force mounted
    if (!isSelected && !forceMount) {
      return null;
    }

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'mt-4 focus:outline-none',

          // Hide if not selected but force mounted
          !isSelected && forceMount && 'hidden',

          // Custom className
          className
        )}
        role={role || 'tabpanel'}
        id={panelId}
        aria-labelledby={triggerId}
        aria-label={ariaLabel}
        aria-describedby={ariaDescribedBy}
        aria-expanded={ariaExpanded}
        aria-disabled={ariaDisabled}
        tabIndex={tabIndex ?? 0}
        {...props}
      >
        {children}
      </div>
    );
  }
);

TabsContent.displayName = 'TabsContent';

export type {
  TabsProps,
  TabsListProps,
  TabsTriggerProps,
  TabsContentProps,
  TabsContextValue
};