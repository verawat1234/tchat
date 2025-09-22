/**
 * Tabs Component Contract Tests
 * CRITICAL: These tests MUST FAIL until Tabs component is implemented
 */

import { render, screen, fireEvent } from '@testing-library/react';
import { Tabs, TabsList, TabsTrigger, TabsContent } from './tabs';
import type { TabsProps, TabsListProps, TabsTriggerProps, TabsContentProps } from '../../../../specs/001-agent-frontend-specialist/contracts/tabs';

describe('Tabs Contract Tests', () => {
  const defaultTabsStructure = (
    <Tabs defaultValue="tab1" testId="tabs-test">
      <TabsList testId="tabs-list">
        <TabsTrigger value="tab1" testId="tab-trigger-1">Tab 1</TabsTrigger>
        <TabsTrigger value="tab2" testId="tab-trigger-2">Tab 2</TabsTrigger>
        <TabsTrigger value="tab3" testId="tab-trigger-3">Tab 3</TabsTrigger>
      </TabsList>
      <TabsContent value="tab1" testId="tab-content-1">Content 1</TabsContent>
      <TabsContent value="tab2" testId="tab-content-2">Content 2</TabsContent>
      <TabsContent value="tab3" testId="tab-content-3">Content 3</TabsContent>
    </Tabs>
  );

  describe('Basic Rendering', () => {
    test('renders tabs container with default value', () => {
      render(defaultTabsStructure);

      expect(screen.getByTestId('tabs-test')).toBeInTheDocument();
      expect(screen.getByTestId('tabs-list')).toBeInTheDocument();
      expect(screen.getByText('Content 1')).toBeInTheDocument();
    });

    test('applies custom className to tabs container', () => {
      render(
        <Tabs defaultValue="tab1" className="custom-tabs" testId="tabs-custom">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>
      );

      const tabs = screen.getByTestId('tabs-custom');
      expect(tabs).toHaveClass('custom-tabs');
    });
  });

  describe('Controlled vs Uncontrolled', () => {
    test('works in uncontrolled mode with defaultValue', () => {
      render(defaultTabsStructure);

      expect(screen.getByText('Content 1')).toBeInTheDocument();
      expect(screen.queryByText('Content 2')).not.toBeInTheDocument();
    });

    test('works in controlled mode with value and onValueChange', () => {
      const onValueChange = vi.fn();
      render(
        <Tabs value="tab2" onValueChange={onValueChange} testId="tabs-controlled">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>
      );

      expect(screen.getByText('Content 2')).toBeInTheDocument();
      expect(screen.queryByText('Content 1')).not.toBeInTheDocument();
    });

    test('calls onValueChange when tab is clicked in controlled mode', () => {
      const onValueChange = vi.fn();
      render(
        <Tabs value="tab1" onValueChange={onValueChange}>
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>
      );

      fireEvent.click(screen.getByText('Tab 2'));
      expect(onValueChange).toHaveBeenCalledWith('tab2');
    });
  });

  describe('Tab Navigation', () => {
    test('switches content when tab trigger is clicked', () => {
      render(defaultTabsStructure);

      expect(screen.getByText('Content 1')).toBeInTheDocument();

      fireEvent.click(screen.getByTestId('tab-trigger-2'));
      expect(screen.getByText('Content 2')).toBeInTheDocument();
      expect(screen.queryByText('Content 1')).not.toBeInTheDocument();
    });

    test('sets active tab trigger correctly', () => {
      render(defaultTabsStructure);

      const tab1Trigger = screen.getByTestId('tab-trigger-1');
      const tab2Trigger = screen.getByTestId('tab-trigger-2');

      expect(tab1Trigger).toHaveAttribute('aria-selected', 'true');
      expect(tab2Trigger).toHaveAttribute('aria-selected', 'false');

      fireEvent.click(tab2Trigger);

      expect(tab1Trigger).toHaveAttribute('aria-selected', 'false');
      expect(tab2Trigger).toHaveAttribute('aria-selected', 'true');
    });
  });

  describe('Orientation', () => {
    test('applies horizontal orientation by default', () => {
      render(defaultTabsStructure);

      const tabsList = screen.getByTestId('tabs-list');
      expect(tabsList).toHaveAttribute('aria-orientation', 'horizontal');
    });

    test('applies vertical orientation when specified', () => {
      render(
        <Tabs defaultValue="tab1" orientation="vertical">
          <TabsList testId="tabs-list-vertical">
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>
      );

      const tabsList = screen.getByTestId('tabs-list-vertical');
      expect(tabsList).toHaveAttribute('aria-orientation', 'vertical');
    });
  });

  describe('Activation Mode', () => {
    test('supports automatic activation mode', () => {
      render(
        <Tabs defaultValue="tab1" activationMode="automatic">
          <TabsList testId="tabs-list-auto">
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>
      );

      const tabsList = screen.getByTestId('tabs-list-auto');
      expect(tabsList).toHaveAttribute('data-activation-mode', 'automatic');
    });

    test('supports manual activation mode', () => {
      render(
        <Tabs defaultValue="tab1" activationMode="manual">
          <TabsList testId="tabs-list-manual">
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>
      );

      const tabsList = screen.getByTestId('tabs-list-manual');
      expect(tabsList).toHaveAttribute('data-activation-mode', 'manual');
    });
  });

  describe('Keyboard Navigation', () => {
    test('supports keyboard navigation with arrow keys', () => {
      render(defaultTabsStructure);

      const tab1Trigger = screen.getByTestId('tab-trigger-1');
      const tab2Trigger = screen.getByTestId('tab-trigger-2');

      tab1Trigger.focus();
      fireEvent.keyDown(tab1Trigger, { key: 'ArrowRight' });

      expect(tab2Trigger).toHaveFocus();
    });

    test('supports loop navigation when enabled', () => {
      render(
        <Tabs defaultValue="tab1" loop>
          <TabsList>
            <TabsTrigger value="tab1" testId="tab-trigger-1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2" testId="tab-trigger-2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>
      );

      const tab2Trigger = screen.getByTestId('tab-trigger-2');
      const tab1Trigger = screen.getByTestId('tab-trigger-1');

      tab2Trigger.focus();
      fireEvent.keyDown(tab2Trigger, { key: 'ArrowRight' });

      expect(tab1Trigger).toHaveFocus();
    });
  });

  describe('Disabled State', () => {
    test('supports disabled tab triggers', () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2" disabled testId="tab-trigger-disabled">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>
      );

      const disabledTab = screen.getByTestId('tab-trigger-disabled');
      expect(disabledTab).toHaveAttribute('aria-disabled', 'true');
      expect(disabledTab).toBeDisabled();
    });

    test('skips disabled tabs in keyboard navigation', () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1" testId="tab-trigger-1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2" disabled>Tab 2</TabsTrigger>
            <TabsTrigger value="tab3" testId="tab-trigger-3">Tab 3</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
          <TabsContent value="tab3">Content 3</TabsContent>
        </Tabs>
      );

      const tab1Trigger = screen.getByTestId('tab-trigger-1');
      const tab3Trigger = screen.getByTestId('tab-trigger-3');

      tab1Trigger.focus();
      fireEvent.keyDown(tab1Trigger, { key: 'ArrowRight' });

      expect(tab3Trigger).toHaveFocus();
    });
  });

  describe('Tab Trigger Features', () => {
    test('renders icon when provided', () => {
      const icon = <span data-testid="tab-icon">★</span>;
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1" icon={icon}>Tab with Icon</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>
      );

      expect(screen.getByTestId('tab-icon')).toBeInTheDocument();
    });

    test('renders badge when provided', () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1" badge="3">Tab with Badge</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>
      );

      expect(screen.getByText('3')).toBeInTheDocument();
    });
  });

  describe('Tab Content Features', () => {
    test('supports forceMount to keep content mounted', () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2" forceMount testId="tab-content-force-mount">
            Content 2 (Force Mounted)
          </TabsContent>
        </Tabs>
      );

      // Content should be in DOM even though tab is not active
      expect(screen.getByTestId('tab-content-force-mount')).toBeInTheDocument();
    });
  });

  describe('Accessibility', () => {
    test('sets proper ARIA attributes on tabs container', () => {
      render(
        <Tabs defaultValue="tab1" aria-label="Main navigation" testId="tabs-aria">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>
      );

      const tabs = screen.getByTestId('tabs-aria');
      expect(tabs).toHaveAttribute('aria-label', 'Main navigation');
    });

    test('sets proper ARIA attributes on tabs list', () => {
      render(defaultTabsStructure);

      const tabsList = screen.getByTestId('tabs-list');
      expect(tabsList).toHaveAttribute('role', 'tablist');
    });

    test('sets proper ARIA attributes on tab triggers', () => {
      render(defaultTabsStructure);

      const tab1Trigger = screen.getByTestId('tab-trigger-1');
      expect(tab1Trigger).toHaveAttribute('role', 'tab');
      expect(tab1Trigger).toHaveAttribute('aria-selected', 'true');
      expect(tab1Trigger).toHaveAttribute('aria-controls');
    });

    test('sets proper ARIA attributes on tab content', () => {
      render(defaultTabsStructure);

      const tabContent = screen.getByTestId('tab-content-1');
      expect(tabContent).toHaveAttribute('role', 'tabpanel');
      expect(tabContent).toHaveAttribute('aria-labelledby');
    });

    test('supports custom tabIndex', () => {
      render(
        <Tabs defaultValue="tab1" tabIndex={0} testId="tabs-tabindex">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>
      );

      const tabs = screen.getByTestId('tabs-tabindex');
      expect(tabs).toHaveAttribute('tabIndex', '0');
    });
  });
});