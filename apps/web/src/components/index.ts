/**
 * Tchat Design System - Component Exports
 * Cross-platform component library for Web, iOS, and Android
 * Constitutional compliance: 97% visual consistency, WCAG 2.1 AA, <200ms load time
 */

// Core Components
export {
  TchatButton,
  type TchatButtonProps,
  type TchatButtonVariant,
  type TchatButtonSize,
} from './TchatButton';

export {
  TchatInput,
  type TchatInputProps,
  type TchatInputType,
  type TchatInputValidationState,
  type TchatInputSize,
} from './TchatInput';

export {
  TchatCard,
  TchatCardHeader,
  TchatCardContent,
  TchatCardFooter,
  type TchatCardProps,
  type TchatCardVariant,
  type TchatCardSize,
  type TchatCardHeaderProps,
  type TchatCardContentProps,
  type TchatCardFooterProps,
} from './TchatCard';

// Component metadata for cross-platform validation
export const TCHAT_COMPONENTS = {
  TchatButton: {
    variants: ['primary', 'secondary', 'ghost', 'destructive', 'outline'] as const,
    sizes: ['sm', 'md', 'lg'] as const,
    features: ['loading', 'icons', 'accessibility', 'animations'],
    status: 'implemented',
  },
  TchatInput: {
    types: ['text', 'email', 'password', 'number', 'search', 'multiline'] as const,
    validationStates: ['none', 'valid', 'invalid'] as const,
    sizes: ['sm', 'md', 'lg'] as const,
    features: ['validation', 'icons', 'passwordToggle', 'accessibility'],
    status: 'implemented',
  },
  TchatCard: {
    variants: ['elevated', 'outlined', 'filled', 'glass'] as const,
    sizes: ['compact', 'standard', 'expanded'] as const,
    features: ['interactive', 'glassmorphism', 'accessibility', 'animations'],
    status: 'implemented',
  },
} as const;

// Design system version and metadata
export const DESIGN_SYSTEM_VERSION = '1.0.0';
export const CROSS_PLATFORM_CONSISTENCY_TARGET = 0.97; // 97% Constitutional requirement
export const SUPPORTED_PLATFORMS = ['web', 'ios', 'android'] as const;