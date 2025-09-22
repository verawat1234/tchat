/**
 * Sidebar Components
 * Navigation sidebar with collapsible panels and nested navigation
 */

import React, { createContext, useContext, useState, useEffect, useRef } from 'react';
import { cn } from '@/utils/cn';
import type {
  SidebarProps,
  SidebarContentProps,
  SidebarHeaderProps,
  SidebarFooterProps,
  SidebarNavProps,
  SidebarNavGroupProps,
  SidebarToggleProps,
  SidebarNavItem
} from '../../../../specs/001-agent-frontend-specialist/contracts/sidebar';

/**
 * Sidebar context for sharing state
 */
interface SidebarContextValue {
  collapsed: boolean;
  onCollapsedChange?: (collapsed: boolean) => void;
  position: 'left' | 'right';
}

const SidebarContext = createContext<SidebarContextValue | undefined>(undefined);

const useSidebarContext = () => {
  const context = useContext(SidebarContext);
  return context;
};

/**
 * Main sidebar container component
 */
export const Sidebar = React.forwardRef<HTMLDivElement, SidebarProps>(
  ({
    className,
    children,
    testId,
    position = 'left',
    width = 250,
    collapsible = false,
    collapsed: controlledCollapsed,
    onCollapsedChange,
    overlay = false,
    persistent = false,
    variant = 'default',
    resizable = false,
    minWidth = 200,
    maxWidth = 400,
    onResize,
    'aria-label': ariaLabel,
    'aria-describedby': ariaDescribedBy,
    'aria-expanded': ariaExpanded,
    'aria-disabled': ariaDisabled,
    role,
    tabIndex,
    ...props
  }, ref) => {
    const [internalCollapsed, setInternalCollapsed] = useState(false);
    const [currentWidth, setCurrentWidth] = useState<number | string>(
      typeof width === 'number' ? width :
      typeof width === 'string' ? width :
      250
    );
    const resizeRef = useRef<HTMLDivElement>(null);
    const isResizing = useRef(false);

    const isControlled = controlledCollapsed !== undefined;
    const collapsed = isControlled ? controlledCollapsed : internalCollapsed;

    const handleCollapsedChange = (newCollapsed: boolean) => {
      if (!isControlled) {
        setInternalCollapsed(newCollapsed);
      }
      onCollapsedChange?.(newCollapsed);
    };

    const contextValue: SidebarContextValue = {
      collapsed,
      onCollapsedChange: handleCollapsedChange,
      position
    };

    // Handle resizing
    useEffect(() => {
      if (!resizable || !resizeRef.current) return;

      const handleMouseMove = (e: MouseEvent) => {
        if (!isResizing.current) return;

        const newWidth = position === 'left'
          ? e.clientX
          : window.innerWidth - e.clientX;

        const clampedWidth = Math.min(Math.max(newWidth, minWidth), maxWidth);
        setCurrentWidth(clampedWidth);
        onResize?.(clampedWidth);
      };

      const handleMouseUp = () => {
        isResizing.current = false;
        document.body.style.userSelect = '';
        document.body.style.cursor = '';
      };

      const handleMouseDown = (e: MouseEvent) => {
        const target = e.target as HTMLElement;
        if (target.classList.contains('sidebar-resize-handle')) {
          isResizing.current = true;
          document.body.style.userSelect = 'none';
          document.body.style.cursor = 'ew-resize';
        }
      };

      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
      resizeRef.current.addEventListener('mousedown', handleMouseDown);

      return () => {
        document.removeEventListener('mousemove', handleMouseMove);
        document.removeEventListener('mouseup', handleMouseUp);
        resizeRef.current?.removeEventListener('mousedown', handleMouseDown);
      };
    }, [resizable, position, minWidth, maxWidth, onResize]);

    const getVariantClasses = () => {
      switch (variant) {
        case 'floating':
          return 'shadow-lg rounded-lg m-2';
        case 'bordered':
          return 'border-r border-gray-200';
        default:
          return 'bg-gray-50';
      }
    };

    const getSidebarWidth = () => {
      if (collapsed) return '60px';
      if (typeof currentWidth === 'string') return currentWidth;
      return `${currentWidth}px`;
    };

    return (
      <SidebarContext.Provider value={contextValue}>
        {/* Overlay for mobile */}
        {overlay && !collapsed && (
          <div
            className="fixed inset-0 bg-black/50 z-40 lg:hidden sidebar-overlay"
            onClick={() => handleCollapsedChange(true)}
          />
        )}

        <aside
          ref={ref}
          data-testid={testId}
          className={cn(
            // Base styles
            'flex flex-col h-full transition-all duration-200',

            // Position
            position === 'left' ? 'sidebar-left' : 'sidebar-right',

            // Width and collapsed state
            collapsed && 'sidebar-collapsed',
            !collapsed && 'sidebar-expanded',
            collapsible && 'sidebar-collapsible',

            // Variant
            `sidebar-${variant}`,
            getVariantClasses(),

            // Persistent
            persistent && 'sidebar-persistent',
            resizable && 'sidebar-resizable',

            // Overlay mode (mobile)
            overlay && 'fixed z-50 lg:relative lg:z-auto sidebar-overlay',

            // Custom classes
            'sidebar',
            className
          )}
          style={{
            width: getSidebarWidth(),
            [position]: position === 'left' ? 0 : undefined,
            ...(resizable && { minWidth: `${minWidth}px`, maxWidth: `${maxWidth}px` })
          }}
          aria-label={ariaLabel}
          aria-describedby={ariaDescribedBy}
          aria-expanded={ariaExpanded || !collapsed}
          aria-disabled={ariaDisabled}
          role={role || 'navigation'}
          tabIndex={tabIndex}
          {...props}
        >
          {children}

          {/* Resize handle */}
          {resizable && !collapsed && (
            <div
              ref={resizeRef}
              data-testid="resize-handle"
              className={cn(
                'absolute top-0 bottom-0 w-1 bg-transparent hover:bg-blue-500 cursor-ew-resize sidebar-resize-handle',
                position === 'left' ? '-right-1' : '-left-1'
              )}
            />
          )}
        </aside>
      </SidebarContext.Provider>
    );
  }
);

Sidebar.displayName = 'Sidebar';

/**
 * Sidebar content container
 */
export const SidebarContent = React.forwardRef<HTMLDivElement, SidebarContentProps>(
  ({
    className,
    children,
    testId,
    padding = 'md',
    scrollable = true,
    ...props
  }, ref) => {
    const getPaddingClass = () => {
      const paddingMap = {
        none: 'p-0',
        sm: 'p-2',
        md: 'p-4',
        lg: 'p-6'
      };
      return paddingMap[padding];
    };

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex-1',

          // Padding
          getPaddingClass(),
          `sidebar-content-padding-${padding}`,

          // Scrolling
          scrollable && 'overflow-y-auto sidebar-content-scrollable',

          // Custom classes
          'sidebar-content',
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);

SidebarContent.displayName = 'SidebarContent';

/**
 * Sidebar header component
 */
export const SidebarHeader = React.forwardRef<HTMLDivElement, SidebarHeaderProps>(
  ({
    className,
    children,
    testId,
    title,
    actions,
    border = true,
    sticky = false,
    ...props
  }, ref) => {
    const context = useSidebarContext();
  const collapsed = context?.collapsed ?? false;

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'flex items-center justify-between px-4 py-3',

          // Border
          border && 'border-b border-gray-200 sidebar-header-border',

          // Sticky
          sticky && 'sticky top-0 z-10 bg-white sidebar-header-sticky',

          // Collapsed
          collapsed && 'justify-center',

          // Custom classes
          'sidebar-header',
          className
        )}
        {...props}
      >
        {collapsed ? (
          // Show only first letter or icon when collapsed
          title && (
            <span className="text-lg font-semibold">
              {title.charAt(0).toUpperCase()}
            </span>
          )
        ) : (
          <>
            {title && (
              <h2 className="text-lg font-semibold truncate">{title}</h2>
            )}
            {actions && (
              <div className="flex items-center space-x-2">{actions}</div>
            )}
            {children}
          </>
        )}
      </div>
    );
  }
);

SidebarHeader.displayName = 'SidebarHeader';

/**
 * Sidebar footer component
 */
export const SidebarFooter = React.forwardRef<HTMLDivElement, SidebarFooterProps>(
  ({
    className,
    children,
    testId,
    border = true,
    sticky = false,
    ...props
  }, ref) => {
    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'px-4 py-3',

          // Border
          border && 'border-t border-gray-200 sidebar-footer-border',

          // Sticky
          sticky && 'sticky bottom-0 z-10 bg-white sidebar-footer-sticky',

          // Custom classes
          'sidebar-footer',
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);

SidebarFooter.displayName = 'SidebarFooter';

/**
 * Sidebar navigation component
 */
export const SidebarNav = React.forwardRef<HTMLDivElement, SidebarNavProps>(
  ({
    className,
    testId,
    items,
    value: controlledValue,
    onValueChange,
    multiple = false,
    defaultExpanded = [],
    expanded: controlledExpanded,
    onExpandedChange,
    variant = 'default',
    ...props
  }, ref) => {
    const [internalValue, setInternalValue] = useState<string>('');
    const [internalExpanded, setInternalExpanded] = useState<string[]>(defaultExpanded);
    const context = useSidebarContext();
  const collapsed = context?.collapsed ?? false;

    const isValueControlled = controlledValue !== undefined;
    const value = isValueControlled ? controlledValue : internalValue;

    const isExpandedControlled = controlledExpanded !== undefined;
    const expanded = isExpandedControlled ? controlledExpanded : internalExpanded;

    const handleValueChange = (newValue: string) => {
      if (!isValueControlled) {
        setInternalValue(newValue);
      }
      onValueChange?.(newValue);
    };

    const handleExpandedChange = (itemId: string) => {
      let newExpanded: string[];

      if (expanded.includes(itemId)) {
        newExpanded = expanded.filter(id => id !== itemId);
      } else {
        newExpanded = multiple ? [...expanded, itemId] : [itemId];
      }

      if (!isExpandedControlled) {
        setInternalExpanded(newExpanded);
      }
      onExpandedChange?.(newExpanded);
    };

    const renderNavItem = (item: SidebarNavItem, level = 0): JSX.Element => {
      const isActive = value === item.id;
      const isExpanded = expanded.includes(item.id);
      const hasChildren = item.children && item.children.length > 0;

      const handleClick = () => {
        if (item.disabled) return;

        if (hasChildren && item.expandable !== false) {
          handleExpandedChange(item.id);
        }

        if (!hasChildren || item.href || item.onClick) {
          handleValueChange(item.id);
          item.onClick?.();
        }
      };

      const getVariantClasses = () => {
        switch (variant) {
          case 'pills':
            return isActive
              ? 'bg-blue-500 text-white'
              : 'hover:bg-gray-100';
          case 'tree':
            return isActive
              ? 'bg-blue-50 border-l-2 border-blue-500'
              : 'hover:bg-gray-50';
          default:
            return isActive
              ? 'bg-gray-100 text-gray-900'
              : 'hover:bg-gray-50';
        }
      };

      const ItemContent = item.href ? 'a' : 'button';

      return (
        <div key={item.id} className="nav-item">
          <ItemContent
            href={item.href}
            onClick={handleClick}
            disabled={item.disabled}
            aria-disabled={item.disabled ? 'true' : undefined}
            className={cn(
              // Base styles
              'flex items-center w-full px-3 py-2 text-sm transition-colors duration-150',

              // Level indentation
              level > 0 && `pl-${(level + 1) * 3}`,

              // Variant styles
              getVariantClasses(),

              // State classes
              item.disabled && 'opacity-50 cursor-not-allowed nav-item-disabled',
              isActive && 'nav-item-active',

              // Collapsed mode
              collapsed && 'justify-center'
            )}
            aria-expanded={hasChildren ? isExpanded : undefined}
            aria-current={isActive ? 'page' : undefined}
          >
            {/* Icon */}
            {item.icon && (
              <span className={cn(
                'flex items-center',
                !collapsed && 'mr-2'
              )}>
                {item.icon}
              </span>
            )}

            {/* Label - hidden when collapsed */}
            {!collapsed && (
              <>
                <span
                  className="flex-1 truncate"
                  aria-disabled={item.disabled ? 'true' : undefined}
                >
                  {item.label}
                </span>

                {/* Badge */}
                {item.badge && (
                  <span className="px-1.5 py-0.5 text-xs bg-gray-100 text-gray-600 rounded-full">
                    {item.badge}
                  </span>
                )}

                {/* Expand icon */}
                {hasChildren && (
                  <svg
                    className={cn(
                      'w-4 h-4 transition-transform duration-150',
                      isExpanded && 'rotate-90'
                    )}
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                )}
              </>
            )}
          </ItemContent>

          {/* Children */}
          {hasChildren && isExpanded && !collapsed && (
            <div className="nav-item-children">
              {item.children!.map(child => renderNavItem(child, level + 1))}
            </div>
          )}
        </div>
      );
    };

    return (
      <nav
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'space-y-1',

          // Variant
          `sidebar-nav sidebar-nav-${variant}`,

          // Multiple navigation
          multiple && 'sidebar-nav-multiple',

          // Custom classes
          className
        )}
        {...props}
      >
        {items.map(item => renderNavItem(item))}
      </nav>
    );
  }
);

SidebarNav.displayName = 'SidebarNav';

/**
 * Sidebar navigation group component
 */
export const SidebarNavGroup = React.forwardRef<HTMLDivElement, SidebarNavGroupProps>(
  ({
    className,
    children,
    testId,
    title,
    collapsible = false,
    defaultCollapsed = false,
    collapsed: controlledCollapsed,
    onCollapsedChange,
    ...props
  }, ref) => {
    const [internalCollapsed, setInternalCollapsed] = useState(defaultCollapsed);
    const context = useSidebarContext();
  const sidebarCollapsed = context?.collapsed ?? false;

    const isControlled = controlledCollapsed !== undefined;
    const collapsed = isControlled ? controlledCollapsed : internalCollapsed;

    const handleToggle = () => {
      const newCollapsed = !collapsed;
      if (!isControlled) {
        setInternalCollapsed(newCollapsed);
      }
      onCollapsedChange?.(newCollapsed);
    };

    return (
      <div
        ref={ref}
        data-testid={testId}
        className={cn(
          // Base styles
          'mb-4',

          // Custom classes
          'sidebar-nav-group',
          collapsed && 'sidebar-group-collapsed',
          collapsible && 'sidebar-group-collapsible',
          className
        )}
        {...props}
      >
        {title && !sidebarCollapsed && (
          <button
            type="button"
            onClick={collapsible ? handleToggle : undefined}
            className={cn(
              'flex items-center w-full px-3 py-1 text-xs font-semibold text-gray-500 uppercase tracking-wider',
              collapsible && 'cursor-pointer hover:text-gray-700'
            )}
            aria-expanded={collapsible ? !collapsed : undefined}
          >
            {collapsible && (
              <svg
                className={cn(
                  'w-3 h-3 mr-1 transition-transform duration-150',
                  !collapsed && 'rotate-90'
                )}
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                  clipRule="evenodd"
                />
              </svg>
            )}
            {title}
          </button>
        )}

        {(!collapsible || !collapsed) && !sidebarCollapsed && (
          <div className="nav-group-content">
            {children}
          </div>
        )}
      </div>
    );
  }
);

SidebarNavGroup.displayName = 'SidebarNavGroup';

/**
 * Sidebar toggle button component
 */
export const SidebarToggle = React.forwardRef<HTMLButtonElement, SidebarToggleProps>(
  ({
    className,
    testId,
    collapsed: propCollapsed,
    onToggle,
    position = 'inside',
    direction = 'left',
    ...props
  }, ref) => {
    const context = useSidebarContext();
    const collapsed = propCollapsed ?? context?.collapsed ?? false;

    const handleClick = () => {
      onToggle?.();
      context?.onCollapsedChange?.(!collapsed);
    };

    return (
      <button
        ref={ref}
        type="button"
        data-testid={testId}
        onClick={handleClick}
        className={cn(
          // Base styles
          'p-2 rounded-md transition-colors duration-150',
          'hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500',

          // Position
          position === 'inside' ? 'sidebar-toggle-inside' : 'sidebar-toggle-outside',
          direction === 'left' ? 'sidebar-toggle-left' : 'sidebar-toggle-right',
          collapsed && 'sidebar-toggle-collapsed',

          // Custom classes
          'sidebar-toggle',
          className
        )}
        aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        aria-expanded={!collapsed}
        {...props}
      >
        <svg
          className={cn(
            'w-5 h-5 transition-transform duration-150',
            collapsed && direction === 'left' && 'rotate-180',
            collapsed && direction === 'right' && '-rotate-180'
          )}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d={direction === 'left'
              ? 'M11 19l-7-7 7-7M4 12h16'
              : 'M13 5l7 7-7 7M20 12H4'
            }
          />
        </svg>
      </button>
    );
  }
);

SidebarToggle.displayName = 'SidebarToggle';

export type {
  SidebarProps,
  SidebarContentProps,
  SidebarHeaderProps,
  SidebarFooterProps,
  SidebarNavProps,
  SidebarNavGroupProps,
  SidebarToggleProps,
  SidebarNavItem
};