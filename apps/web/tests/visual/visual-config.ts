/**
 * Visual Regression Testing Configuration
 * Cross-platform UI consistency validation for Web vs iOS
 */

export interface VisualTestConfig {
  threshold: number;
  clip?: { x: number; y: number; width: number; height: number };
  fullPage?: boolean;
  animations?: 'disabled' | 'allow';
  mask?: string[];
}

export interface ComponentTestCase {
  name: string;
  component: string;
  variants: string[];
  sizes: string[];
  states?: string[];
  props?: Record<string, any>;
}

// Cross-platform visual consistency thresholds
export const VISUAL_THRESHOLDS = {
  // 95% visual consistency requirement
  DEFAULT: 0.05, // 5% difference threshold
  STRICT: 0.02,  // 2% for critical components
  RELAXED: 0.1,  // 10% for complex animations
} as const;

// iOS-specific viewport configurations to match native screens
export const IOS_VIEWPORTS = {
  'iPhone 12': { width: 390, height: 844 },
  'iPhone 12 Pro Max': { width: 428, height: 926 },
  'iPhone SE': { width: 375, height: 667 },
  'iPad Air': { width: 820, height: 1180 },
  'iPad Pro 12.9': { width: 1024, height: 1366 },
} as const;

// Component categories for systematic testing
export const COMPONENT_CATEGORIES = {
  CORE_INTERACTIVE: [
    'TchatButton',
    'TchatInput',
    'TchatCard',
    'TchatTabs',
    'TchatAlert',
    'TchatButton',
    'TchatToast',
    'TchatTooltip'
  ],
  MISSING_HIGH_PRIORITY: [
    'TchatDialog',
    'TchatDrawer',
    'TchatPopover',
    'TchatDropdownMenu',
    'TchatCommand'
  ],
  DATA_DISPLAY: [
    'TchatCalendar',
    'TchatChart',
    'TchatCarousel',
    'TchatTable',
    'TchatProgress'
  ],
  LAYOUT: [
    'TchatAccordion',
    'TchatCollapsible',
    'TchatMenubar',
    'TchatNavigationMenu',
    'TchatAvatar'
  ],
  FEEDBACK: [
    'TchatBadge',
    'TchatSkeleton',
    'TchatSlider',
    'TchatSeparator'
  ],
  SPECIALIZED: [
    'TchatForm',
    'TchatHoverCard',
    'TchatInputOtp',
    'TchatResizable',
    'TchatScrollArea',
    'TchatSheet',
    'TchatToggleGroup',
    'TchatToggle',
    'TchatRadioGroup'
  ]
} as const;

// Standard variants across all components
export const STANDARD_VARIANTS = ['primary', 'secondary', 'outline', 'ghost', 'destructive'] as const;
export const STANDARD_SIZES = ['small', 'medium', 'large'] as const;
export const STANDARD_STATES = ['default', 'hover', 'active', 'disabled', 'loading'] as const;

// Test configuration for systematic component testing
export const TEST_CASES: ComponentTestCase[] = [
  {
    name: 'TchatButton',
    component: 'Button',
    variants: [...STANDARD_VARIANTS],
    sizes: [...STANDARD_SIZES],
    states: [...STANDARD_STATES],
    props: { text: 'Test Button' }
  },
  {
    name: 'TchatInput',
    component: 'Input',
    variants: ['default', 'error', 'success'],
    sizes: [...STANDARD_SIZES],
    states: ['default', 'focus', 'disabled'],
    props: { placeholder: 'Enter text...' }
  },
  {
    name: 'TchatCard',
    component: 'Card',
    variants: ['elevated', 'outlined', 'filled', 'glass'],
    sizes: ['small', 'medium', 'large'],
    props: { title: 'Test Card', content: 'Card content for testing' }
  }
];

// Visual test utilities for cross-platform comparison
export const VISUAL_TEST_CONFIG: VisualTestConfig = {
  threshold: VISUAL_THRESHOLDS.DEFAULT,
  animations: 'disabled', // Disable animations for consistent screenshots
  fullPage: false,
  clip: { x: 0, y: 0, width: 800, height: 600 } // Standard component clip area
};