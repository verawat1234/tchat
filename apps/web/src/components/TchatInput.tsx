/**
 * TchatInput Component - Web Implementation
 * Cross-platform design system input component with validation states
 * Constitutional requirements: 97% visual consistency, WCAG 2.1 AA, <200ms load time
 */
import React, { forwardRef, useState } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '../lib/utils';

// Input variants using class-variance-authority for type safety
const inputVariants = cva(
  // Base styles - shared across all variants
  [
    'flex w-full rounded-md border bg-white px-3 py-2',
    'text-sm placeholder:text-text-muted',
    'transition-all duration-200',
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2',
    'disabled:cursor-not-allowed disabled:opacity-50',
    // Touch target compliance (Constitutional requirement: 44dp minimum)
    'min-h-[44px]',
    // Animation support for 60fps performance
    'transform-gpu will-change-[border-color,box-shadow]',
  ],
  {
    variants: {
      // 3 Validation states with distinct visual feedback
      validationState: {
        none: [
          'border-border',
          'hover:border-border-hover',
          'focus:border-primary',
        ],
        valid: [
          'border-success bg-success/5',
          'hover:border-success/80',
          'focus:border-success focus:ring-success/30',
        ],
        invalid: [
          'border-error bg-error/5',
          'hover:border-error/80',
          'focus:border-error focus:ring-error/30',
        ],
      },
      // 3 Size variants for different use cases
      size: {
        sm: ['h-8 px-2 text-xs min-h-[32px]'], // Small: 32dp height
        md: ['h-11 px-3 text-sm min-h-[44px]'], // Medium: 44dp height (default)
        lg: ['h-12 px-4 text-base min-h-[48dp]'], // Large: 48dp height
      },
    },
    defaultVariants: {
      validationState: 'none',
      size: 'md',
    },
  }
);

// Eye icon for password visibility toggle
const EyeIcon = ({ open }: { open: boolean }) => (
  <svg
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className="text-text-muted"
  >
    {open ? (
      // Eye open icon
      <>
        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
        <circle cx="12" cy="12" r="3" />
      </>
    ) : (
      // Eye closed icon
      <>
        <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
        <line x1="1" y1="1" x2="23" y2="23" />
      </>
    )}
  </svg>
);

// Success checkmark icon
const CheckIcon = () => (
  <svg
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className="text-success"
  >
    <polyline points="20,6 9,17 4,12" />
  </svg>
);

// Error X icon
const XIcon = () => (
  <svg
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className="text-error"
  >
    <line x1="18" y1="6" x2="6" y2="18" />
    <line x1="6" y1="6" x2="18" y2="18" />
  </svg>
);

// Search icon
const SearchIcon = () => (
  <svg
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className="text-text-muted"
  >
    <circle cx="11" cy="11" r="8" />
    <path d="m21 21-4.35-4.35" />
  </svg>
);

export interface TchatInputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'>,
    VariantProps<typeof inputVariants> {
  /** Input type - affects keyboard and behavior */
  type?: 'text' | 'email' | 'password' | 'number' | 'search' | 'multiline';
  /** Validation state with visual feedback */
  validationState?: 'none' | 'valid' | 'invalid';
  /** Error message displayed below input */
  error?: string;
  /** Label text displayed above input */
  label?: string;
  /** Show password visibility toggle for password inputs */
  showPasswordToggle?: boolean;
  /** Leading icon element */
  leadingIcon?: React.ReactNode;
  /** Trailing action element */
  trailingAction?: React.ReactNode;
  /** Custom class name */
  className?: string;
  /** Callback for clear action */
  onClear?: () => void;
  /** ARIA describedby for error messages */
  'aria-describedby'?: string;
  /** Content description for accessibility */
  contentDescription?: string;
}

/**
 * TchatInput - Cross-platform design system input component
 *
 * Features:
 * - 6 input types (text, email, password, number, search, multiline)
 * - 3 validation states (none, valid, invalid) with visual feedback
 * - 3 size variants (sm, md, lg) with proper touch targets
 * - Password visibility toggle
 * - Icon support (leading/trailing)
 * - WCAG 2.1 AA accessibility compliance
 * - Animated borders and focus states
 * - Cross-platform visual consistency (97% target)
 */
export const TchatInput = forwardRef<HTMLInputElement, TchatInputProps>(
  (
    {
      type = 'text',
      validationState = 'none',
      size = 'md',
      error,
      label,
      showPasswordToggle = false,
      leadingIcon,
      trailingAction,
      className,
      disabled,
      onClear,
      'aria-describedby': ariaDescribedBy,
      contentDescription,
      ...props
    },
    ref
  ) => {
    const [showPassword, setShowPassword] = useState(false);
    const [isFocused, setIsFocused] = useState(false);

    // Handle password visibility toggle
    const togglePasswordVisibility = () => {
      setShowPassword(!showPassword);
    };

    // Determine actual input type
    const inputType = type === 'password' && showPassword ? 'text' :
                     type === 'search' ? 'text' :
                     type === 'multiline' ? 'text' : type;

    // Generate unique IDs for accessibility
    const inputId = React.useId();
    const errorId = `${inputId}-error`;
    const labelId = `${inputId}-label`;

    // Determine trailing content
    const getTrailingContent = () => {
      if (type === 'password' && showPasswordToggle) {
        return (
          <button
            type="button"
            className="p-1 hover:bg-surface rounded-sm transition-colors"
            onClick={togglePasswordVisibility}
            aria-label="Toggle password visibility"
            data-testid="password-toggle"
          >
            <EyeIcon open={showPassword} />
          </button>
        );
      }

      if (validationState === 'valid') {
        return <CheckIcon />;
      }

      if (validationState === 'invalid') {
        return <XIcon />;
      }

      if (trailingAction) {
        return trailingAction;
      }

      if (props.value && onClear) {
        return (
          <button
            type="button"
            className="p-1 hover:bg-surface rounded-sm transition-colors"
            onClick={onClear}
            aria-label="Clear input"
            data-testid="clear-button"
          >
            <XIcon />
          </button>
        );
      }

      return null;
    };

    // Determine leading content
    const getLeadingContent = () => {
      if (type === 'search') {
        return <SearchIcon />;
      }

      if (leadingIcon) {
        return leadingIcon;
      }

      return null;
    };

    const trailingContent = getTrailingContent();
    const leadingContent = getLeadingContent();

    return (
      <div className="w-full">
        {/* Label */}
        {label && (
          <label
            htmlFor={inputId}
            id={labelId}
            className="block text-sm font-medium text-text-primary mb-1"
          >
            {label}
          </label>
        )}

        {/* Input container */}
        <div className="relative">
          {/* Leading icon */}
          {leadingContent && (
            <div className="absolute left-3 top-1/2 transform -translate-y-1/2 flex items-center">
              {leadingContent}
            </div>
          )}

          {/* Input element */}
          {type === 'multiline' ? (
            <textarea
              ref={ref as any}
              id={inputId}
              className={cn(
                inputVariants({ validationState, size }),
                leadingContent && 'pl-10',
                trailingContent && 'pr-10',
                'resize-none min-h-[88px]', // Multiline specific styles
                className
              )}
              disabled={disabled}
              aria-describedby={error ? errorId : ariaDescribedBy}
              aria-invalid={validationState === 'invalid'}
              aria-labelledby={label ? labelId : undefined}
              data-testid="tchat-input"
              data-type={type}
              data-validation-state={validationState}
              data-size={size}
              onFocus={(e) => {
                setIsFocused(true);
                props.onFocus?.(e as any);
              }}
              onBlur={(e) => {
                setIsFocused(false);
                props.onBlur?.(e as any);
              }}
              {...(props as any)}
            />
          ) : (
            <input
              ref={ref}
              type={inputType}
              id={inputId}
              className={cn(
                inputVariants({ validationState, size }),
                leadingContent && 'pl-10',
                trailingContent && 'pr-10',
                className
              )}
              disabled={disabled}
              aria-describedby={error ? errorId : ariaDescribedBy}
              aria-invalid={validationState === 'invalid'}
              aria-labelledby={label ? labelId : undefined}
              data-testid="tchat-input"
              data-type={type}
              data-validation-state={validationState}
              data-size={size}
              onFocus={(e) => {
                setIsFocused(true);
                props.onFocus?.(e);
              }}
              onBlur={(e) => {
                setIsFocused(false);
                props.onBlur?.(e);
              }}
              {...props}
            />
          )}

          {/* Trailing content */}
          {trailingContent && (
            <div className="absolute right-3 top-1/2 transform -translate-y-1/2 flex items-center">
              {trailingContent}
            </div>
          )}
        </div>

        {/* Error message */}
        {error && (
          <p
            id={errorId}
            className="mt-1 text-sm text-error"
            role="alert"
            data-testid="input-error"
          >
            {error}
          </p>
        )}
      </div>
    );
  }
);

TchatInput.displayName = 'TchatInput';

// Export types for external use
export type TchatInputType = TchatInputProps['type'];
export type TchatInputValidationState = TchatInputProps['validationState'];
export type TchatInputSize = VariantProps<typeof inputVariants>['size'];