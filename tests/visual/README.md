# Cross-Platform Visual Regression Testing Infrastructure

Comprehensive visual testing system for ensuring >95% consistency between web and iOS UI components.

## Architecture

### Core Components

- **VisualTestRunner**: Main testing orchestrator with iOS-specific viewport and styling
- **Visual Config**: Centralized configuration for thresholds, viewports, and test cases
- **Test Utils**: Cross-platform screenshot comparison and consistency validation
- **Component Tests**: Individual test files for each UI component

### iOS Integration Strategy

1. **Viewport Matching**: Tests use iOS device viewports (iPhone 12, iPad Air, etc.)
2. **Font Rendering**: Applies iOS-specific font smoothing and system fonts
3. **Design Token Validation**: Compares computed CSS values against iOS specifications
4. **Cross-Platform Screenshots**: Captures web components styled to match iOS appearance

## Test Categories

### Priority 1: Core Interactive Components ✅
- TchatButton, TchatInput, TchatCard, TchatTabs, TchatAlert, TchatToast, TchatTooltip
- **Status**: Testing infrastructure ready

### Priority 2: Missing High Priority Components 🚧
- TchatDialog, TchatDrawer, TchatPopover, TchatDropdownMenu, TchatCommand
- **Status**: Tests written, will pass once components are implemented

### Priority 3: Data Display & Layout 📋
- Calendar, Chart, Carousel, Table, Progress, Accordion, etc.
- **Status**: Tests planned in tasks.md (T041-T068)

## Usage

### Running Visual Tests

```bash
# Run all visual regression tests
npm run test:e2e tests/visual/

# Run specific component tests
npm run test:e2e tests/visual/components/dialog.spec.ts

# Run with UI mode for debugging
npm run test:e2e:ui tests/visual/

# Generate visual report
npm run test:e2e tests/visual/cross-platform-runner.spec.ts
```

### Component Test Development

1. **Create Component Test File**:
   ```typescript
   // tests/visual/components/new-component.spec.ts
   import { test, expect } from '@playwright/test';
   import { VisualTestRunner } from '../visual-test-utils';
   ```

2. **Configure Test Case**:
   ```typescript
   const results = await visualTester.testComponentVariants(
     'TchatNewComponent',
     ['primary', 'secondary'],
     ['small', 'medium'],
     { threshold: VISUAL_THRESHOLDS.STRICT }
   );
   ```

3. **Validate Results**:
   ```typescript
   expect(result.consistencyScore).toBeGreaterThan(0.95);
   ```

## Quality Gates

### Consistency Thresholds
- **Default**: 5% visual difference allowed
- **Strict**: 2% for critical components (dialogs, forms)
- **Relaxed**: 10% for complex animations

### Performance Targets
- **Load Time**: <200ms per component
- **Memory Usage**: <100MB mobile, <500MB desktop
- **Render Time**: <3s on 3G networks

### Accessibility Compliance
- **WCAG 2.1 AA**: Minimum compliance level
- **Keyboard Navigation**: Full keyboard support
- **Screen Reader**: Proper ARIA attributes
- **Touch Targets**: Minimum 44px iOS compliance

## Integration with Development Workflow

### TDD Approach (Required)
1. **Write Visual Test First**: Create failing test for new component
2. **Implement Component**: Build component to pass visual consistency
3. **Validate Cross-Platform**: Ensure >95% iOS consistency score
4. **Update Test Cases**: Add variants, sizes, states as needed

### CI/CD Integration
- Tests run automatically on component changes
- Visual diff reports generated for PR reviews
- Consistency threshold violations block deployment
- Cross-platform screenshots archived for comparison

## File Structure

```
tests/visual/
├── README.md                          # This documentation
├── visual-config.ts                   # Configuration and constants
├── visual-test-utils.ts               # Testing utilities and runner
├── cross-platform-runner.spec.ts     # Comprehensive test suite
└── components/                        # Individual component tests
    ├── dialog.spec.ts                 # TchatDialog tests (example)
    ├── drawer.spec.ts                 # TchatDrawer tests
    ├── popover.spec.ts                # TchatPopover tests
    └── ...                            # Additional component tests
```

## Implementation Status

- ✅ **Infrastructure Setup**: Core testing framework operational
- ✅ **Configuration System**: Thresholds, viewports, test cases defined
- ✅ **iOS Integration**: Device viewports and styling configured
- ✅ **Example Tests**: Dialog component test demonstrates patterns
- ✅ **Performance Monitoring**: Load time and memory validation
- 🚧 **Component Coverage**: Expanding to all 39+ missing iOS components
- 📋 **CI/CD Integration**: Planned for automated quality gates

## Next Steps

1. **Implement Missing Components**: Create iOS components per tasks.md T036-T068
2. **Expand Test Coverage**: Add visual tests for all component variants
3. **iOS Screenshot Capture**: Integrate actual iOS app screenshot comparison
4. **Automated Reporting**: Generate consistency reports for stakeholders
5. **Performance Optimization**: Ensure <200ms render time targets

## Related Tasks

- **T002** ✅: Set up visual regression testing infrastructure
- **T031-T035**: Visual regression tests for Priority 1 components
- **T041-T064**: Visual tests for remaining component categories
- **T079**: Comprehensive visual regression testing execution
- **T080**: Validate 95% visual consistency score achievement