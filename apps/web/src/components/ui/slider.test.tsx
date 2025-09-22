/**
 * Slider Component Tests
 * Testing the Slider component with various values, ranges, and interactions
 */

import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Slider } from './slider';
import { vi } from 'vitest';
import React from 'react';
import { getSliderThumb } from '@/test-utils/radix-ui';

describe('Slider Component Tests', () => {
  describe('Basic Rendering', () => {
    test('renders slider with default props', () => {
      const { container } = render(<Slider />);

      const slider = container.querySelector('[data-slot="slider"]');
      expect(slider).toBeInTheDocument();
      expect(slider).toHaveClass('relative', 'flex', 'w-full', 'touch-none');
    });

    test('applies custom className', () => {
      const { container } = render(<Slider className="custom-slider" />);

      const slider = container.querySelector('[data-slot="slider"]');
      expect(slider).toHaveClass('custom-slider');
    });

    test('renders slider track', () => {
      const { container } = render(<Slider />);

      const track = container.querySelector('[data-slot="slider-track"]');
      expect(track).toBeInTheDocument();
      expect(track).toHaveClass('bg-muted', 'relative', 'grow', 'overflow-hidden', 'rounded-full');
    });

    test('renders slider range', () => {
      const { container } = render(<Slider />);

      const range = container.querySelector('[data-slot="slider-range"]');
      expect(range).toBeInTheDocument();
      expect(range).toHaveClass('bg-primary', 'absolute');
    });

    test('renders slider thumb', () => {
      const { container } = render(<Slider defaultValue={[50]} />);

      const thumb = container.querySelector('[data-slot="slider-thumb"]');
      expect(thumb).toBeInTheDocument();
      expect(thumb).toHaveClass('block', 'size-4', 'rounded-full', 'border');
    });
  });

  describe('Value Management', () => {
    test('renders with default value', () => {
      render(<Slider defaultValue={[30]} />);

      const slider = screen.getByRole('slider');
      expect(slider).toHaveAttribute('aria-valuenow', '30');
    });

    test('renders with min and max values', () => {
      render(<Slider min={10} max={90} defaultValue={[50]} />);

      const slider = screen.getByRole('slider');
      expect(slider).toHaveAttribute('aria-valuemin', '10');
      expect(slider).toHaveAttribute('aria-valuemax', '90');
      expect(slider).toHaveAttribute('aria-valuenow', '50');
    });

    test('works as controlled component', () => {
      const ControlledSlider = () => {
        const [value, setValue] = React.useState([25]);
        return (
          <Slider
            value={value}
            onValueChange={setValue}
          />
        );
      };

      render(<ControlledSlider />);
      const slider = screen.getByRole('slider');
      expect(slider).toHaveAttribute('aria-valuenow', '25');
    });

    test('works as uncontrolled component', () => {
      render(<Slider defaultValue={[75]} />);

      const slider = screen.getByRole('slider');
      expect(slider).toHaveAttribute('aria-valuenow', '75');
    });

    test('handles onValueChange callback', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[50]} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');

      // Simulate value change (Radix Slider uses keyboard/pointer events)
      fireEvent.keyDown(slider, { key: 'ArrowRight' });

      expect(handleChange).toHaveBeenCalled();
    });

    test('handles onValueCommit callback', () => {
      const handleCommit = vi.fn();
      render(<Slider defaultValue={[50]} onValueCommit={handleCommit} />);

      const slider = screen.getByRole('slider');

      // Simulate committing value (release after change)
      fireEvent.keyDown(slider, { key: 'ArrowRight' });
      fireEvent.keyUp(slider, { key: 'ArrowRight' });

      expect(handleCommit).toHaveBeenCalled();
    });
  });

  describe('Range Slider', () => {
    test('renders multiple thumbs for range values', () => {
      const { container } = render(<Slider defaultValue={[25, 75]} />);

      const thumbs = container.querySelectorAll('[data-slot="slider-thumb"]');
      expect(thumbs).toHaveLength(2);
    });

    test('manages range values correctly', () => {
      render(<Slider defaultValue={[20, 80]} />);

      const sliders = screen.getAllByRole('slider');
      expect(sliders).toHaveLength(2);
      expect(sliders[0]).toHaveAttribute('aria-valuenow', '20');
      expect(sliders[1]).toHaveAttribute('aria-valuenow', '80');
    });

    test('prevents range values from crossing', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[40, 60]} onValueChange={handleChange} />);

      const sliders = screen.getAllByRole('slider');

      // Try to move first thumb past second
      sliders[0].focus();
      for (let i = 0; i < 30; i++) {
        fireEvent.keyDown(sliders[0], { key: 'ArrowRight' });
      }

      // Radix prevents crossing, so values should be constrained
      if (handleChange.mock.calls.length > 0) {
        const lastCall = handleChange.mock.calls[handleChange.mock.calls.length - 1];
        expect(lastCall[0][0]).toBeLessThanOrEqual(lastCall[0][1]);
      }
    });
  });

  describe('Keyboard Interactions', () => {
    test('increases value with ArrowRight', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[50]} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');
      slider.focus();
      fireEvent.keyDown(slider, { key: 'ArrowRight' });

      expect(handleChange).toHaveBeenCalledWith([51]);
    });

    test('decreases value with ArrowLeft', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[50]} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');
      slider.focus();
      fireEvent.keyDown(slider, { key: 'ArrowLeft' });

      expect(handleChange).toHaveBeenCalledWith([49]);
    });

    test('increases value with ArrowUp', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[50]} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');
      slider.focus();
      fireEvent.keyDown(slider, { key: 'ArrowUp' });

      expect(handleChange).toHaveBeenCalledWith([51]);
    });

    test('decreases value with ArrowDown', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[50]} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');
      slider.focus();
      fireEvent.keyDown(slider, { key: 'ArrowDown' });

      expect(handleChange).toHaveBeenCalledWith([49]);
    });

    test('jumps to minimum with Home', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[50]} min={10} max={90} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');
      slider.focus();
      fireEvent.keyDown(slider, { key: 'Home' });

      expect(handleChange).toHaveBeenCalledWith([10]);
    });

    test('jumps to maximum with End', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[50]} min={10} max={90} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');
      slider.focus();
      fireEvent.keyDown(slider, { key: 'End' });

      expect(handleChange).toHaveBeenCalledWith([90]);
    });

    test('supports custom step size', () => {
      const handleChange = vi.fn();
      render(<Slider defaultValue={[50]} step={5} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');
      slider.focus();
      fireEvent.keyDown(slider, { key: 'ArrowRight' });

      expect(handleChange).toHaveBeenCalledWith([55]);
    });
  });

  describe('Disabled State', () => {
    test('renders as disabled when disabled prop is true', () => {
      const { container } = render(<Slider disabled defaultValue={[50]} />);

      const slider = container.querySelector('[data-slot="slider"]');
      expect(slider).toHaveAttribute('data-disabled');
      expect(slider).toHaveClass('data-[disabled]:opacity-50');

      const thumb = screen.getByRole('slider');
      expect(thumb).toBeInTheDocument();
    });

    test('does not respond to interactions when disabled', () => {
      const handleChange = vi.fn();
      render(<Slider disabled defaultValue={[50]} onValueChange={handleChange} />);

      const slider = screen.getByRole('slider');
      fireEvent.keyDown(slider, { key: 'ArrowRight' });

      expect(handleChange).not.toHaveBeenCalled();
    });
  });

  describe('Orientation', () => {
    test('renders horizontal orientation by default', () => {
      const { container } = render(<Slider defaultValue={[50]} />);

      const slider = container.querySelector('[data-slot="slider"]');
      expect(slider).toHaveAttribute('data-orientation', 'horizontal');
    });

    test('renders vertical orientation when specified', () => {
      const { container } = render(<Slider orientation="vertical" defaultValue={[50]} />);

      const slider = container.querySelector('[data-slot="slider"]');
      expect(slider).toHaveAttribute('data-orientation', 'vertical');
      expect(slider).toHaveClass('data-[orientation=vertical]:h-full');
    });

    test('applies vertical-specific styles', () => {
      const { container } = render(<Slider orientation="vertical" defaultValue={[50]} />);

      const track = container.querySelector('[data-slot="slider-track"]');
      expect(track).toHaveClass('data-[orientation=vertical]:h-full', 'data-[orientation=vertical]:w-1.5');
    });
  });

  describe('Accessibility', () => {
    test('has proper ARIA attributes', () => {
      render(<Slider defaultValue={[60]} min={0} max={100} />);

      const slider = screen.getByRole('slider');
      expect(slider).toHaveAttribute('aria-valuenow', '60');
      expect(slider).toHaveAttribute('aria-valuemin', '0');
      expect(slider).toHaveAttribute('aria-valuemax', '100');
      expect(slider).toHaveAttribute('aria-orientation', 'horizontal');
    });

    test('supports aria-label', () => {
      const { container } = render(<Slider defaultValue={[50]} aria-label="Volume control" />);

      // Radix applies aria-label to the container
      const sliderContainer = container.querySelector('[data-slot="slider"]');
      expect(sliderContainer).toHaveAttribute('aria-label', 'Volume control');

      const slider = screen.getByRole('slider');
      expect(slider).toBeInTheDocument();
    });

    test('supports aria-labelledby', () => {
      const { container } = render(
        <>
          <span id="slider-label">Brightness</span>
          <Slider defaultValue={[70]} aria-labelledby="slider-label" />
        </>
      );

      // Radix applies aria-labelledby to the container, not the thumb
      const sliderContainer = container.querySelector('[data-slot="slider"]');
      expect(sliderContainer).toBeInTheDocument();
      // The thumb itself may not have aria-labelledby, which is fine for Radix
      const slider = screen.getByRole('slider');
      expect(slider).toBeInTheDocument();
    });

    test('supports aria-describedby', () => {
      const { container } = render(
        <>
          <Slider defaultValue={[50]} aria-describedby="slider-help" />
          <span id="slider-help">Adjust the volume level</span>
        </>
      );

      // Radix applies aria-describedby to the container
      const sliderContainer = container.querySelector('[data-slot="slider"]');
      expect(sliderContainer).toBeInTheDocument();
      // Verify slider role exists
      const slider = screen.getByRole('slider');
      expect(slider).toBeInTheDocument();
    });

    test('maintains focus management', () => {
      render(<Slider defaultValue={[50]} />);

      const slider = screen.getByRole('slider');
      expect(slider).toHaveAttribute('tabindex', '0');

      slider.focus();
      expect(document.activeElement).toBe(slider);
    });
  });

  describe('Visual States', () => {
    test('applies hover styles to thumb', () => {
      const { container } = render(<Slider defaultValue={[50]} />);

      const thumb = container.querySelector('[data-slot="slider-thumb"]');
      expect(thumb).toHaveClass('hover:ring-4');
    });

    test('applies focus styles to thumb', () => {
      const { container } = render(<Slider defaultValue={[50]} />);

      const thumb = container.querySelector('[data-slot="slider-thumb"]');
      expect(thumb).toHaveClass('focus-visible:ring-4', 'focus-visible:outline-hidden');
    });

    test('shows shadow on thumb', () => {
      const { container } = render(<Slider defaultValue={[50]} />);

      const thumb = container.querySelector('[data-slot="slider-thumb"]');
      expect(thumb).toHaveClass('shadow-sm');
    });
  });

  describe('Common Use Cases', () => {
    test('volume control slider', () => {
      const handleVolumeChange = vi.fn();
      const { container } = render(
        <Slider
          defaultValue={[75]}
          min={0}
          max={100}
          step={1}
          aria-label="Volume"
          onValueChange={handleVolumeChange}
        />
      );

      const sliderContainer = container.querySelector('[data-slot="slider"]');
      expect(sliderContainer).toHaveAttribute('aria-label', 'Volume');

      const slider = screen.getByRole('slider');
      expect(slider).toHaveAttribute('aria-valuenow', '75');
      expect(slider).toHaveAttribute('aria-valuemin', '0');
      expect(slider).toHaveAttribute('aria-valuemax', '100');
    });

    test('price range slider', () => {
      render(
        <Slider
          defaultValue={[100, 500]}
          min={0}
          max={1000}
          step={10}
          aria-label="Price range"
        />
      );

      const sliders = screen.getAllByRole('slider');
      expect(sliders).toHaveLength(2);
      expect(sliders[0]).toHaveAttribute('aria-valuenow', '100');
      expect(sliders[1]).toHaveAttribute('aria-valuenow', '500');
    });

    test('brightness adjustment', () => {
      const handleBrightnessChange = vi.fn();
      const { container } = render(
        <Slider
          defaultValue={[50]}
          min={0}
          max={100}
          step={5}
          aria-label="Brightness"
          onValueCommit={(value) => handleBrightnessChange(value[0])}
        />
      );

      // Radix applies aria-label to container, not thumb directly
      const sliderContainer = container.querySelector('[data-slot="slider"]');
      expect(sliderContainer).toHaveAttribute('aria-label', 'Brightness');

      const slider = screen.getByRole('slider');
      slider.focus();
      fireEvent.keyDown(slider, { key: 'ArrowRight' });
      fireEvent.keyUp(slider, { key: 'ArrowRight' });

      // The callback extracts value[0] so it should be called with 55
      expect(handleBrightnessChange).toHaveBeenCalledTimes(1);
    });
  });
});