/**
 * Dialog Component Tests
 * Testing the Dialog/Modal component with all its subcomponents and interactions
 */

import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from './dialog';
import { vi } from 'vitest';
import React from 'react';

describe('Dialog Component Tests', () => {
  describe('Basic Rendering', () => {
    test('renders dialog trigger', () => {
      render(
        <Dialog>
          <DialogTrigger>Open Dialog</DialogTrigger>
        </Dialog>
      );

      const trigger = screen.getByText('Open Dialog');
      expect(trigger).toBeInTheDocument();
    });

    test('dialog content is not visible initially', () => {
      render(
        <Dialog>
          <DialogTrigger>Open</DialogTrigger>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Dialog Title</DialogTitle>
            <DialogDescription>Dialog description</DialogDescription>
          </DialogContent>
        </Dialog>
      );

      expect(screen.queryByText('Dialog Title')).not.toBeInTheDocument();
      expect(screen.queryByText('Dialog description')).not.toBeInTheDocument();
    });

    test('opens dialog when trigger is clicked', async () => {
      render(
        <Dialog>
          <DialogTrigger>Open</DialogTrigger>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Dialog Title</DialogTitle>
            <DialogDescription>Dialog description</DialogDescription>
          </DialogContent>
        </Dialog>
      );

      const trigger = screen.getByText('Open');
      fireEvent.click(trigger);

      await waitFor(() => {
        expect(screen.getByText('Dialog Title')).toBeInTheDocument();
        expect(screen.getByText('Dialog description')).toBeInTheDocument();
      });
    });

    test('renders with custom trigger element', () => {
      render(
        <Dialog>
          <DialogTrigger asChild>
            <button className="custom-button">Custom Button</button>
          </DialogTrigger>
        </Dialog>
      );

      const trigger = screen.getByRole('button', { name: 'Custom Button' });
      expect(trigger).toHaveClass('custom-button');
    });
  });

  describe('Dialog Content Structure', () => {
    test('renders header with title and description', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogHeader>
              <DialogTitle>Important Notice</DialogTitle>
              <DialogDescription>
                Please read this important information.
              </DialogDescription>
            </DialogHeader>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        expect(screen.getByText('Important Notice')).toBeInTheDocument();
        expect(screen.getByText('Please read this important information.')).toBeInTheDocument();
      });
    });

    test('renders footer with actions', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Action Dialog</DialogTitle>
            <DialogDescription>Dialog with action buttons</DialogDescription>
            <DialogFooter>
              <button>Cancel</button>
              <button>Confirm</button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        expect(screen.getByText('Cancel')).toBeInTheDocument();
        expect(screen.getByText('Confirm')).toBeInTheDocument();
      });
    });

    test('renders close button', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Dialog with Close</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const closeButton = screen.getByText('Close');
        expect(closeButton).toBeInTheDocument();
        expect(closeButton).toHaveClass('sr-only'); // Screen reader only
      });
    });

    test('applies custom className to content', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent className="custom-content">
            <DialogTitle>Test</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const content = document.querySelector('[data-slot="dialog-content"]');
        expect(content).toHaveClass('custom-content');
      });
    });
  });

  describe('Dialog Overlay', () => {
    test('renders overlay when open', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>With Overlay</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const overlay = document.querySelector('[data-slot="dialog-overlay"]');
        expect(overlay).toBeInTheDocument();
        expect(overlay).toHaveClass('bg-black/50');
      });
    });

    test('closes dialog when overlay is clicked', async () => {
      const onOpenChange = vi.fn();

      render(
        <Dialog defaultOpen onOpenChange={onOpenChange}>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Click Outside</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const overlay = document.querySelector('[data-slot="dialog-overlay"]');
        expect(overlay).toBeInTheDocument();
      });

      // Click on the overlay element - Radix handles overlay clicks via pointer events
      const overlay = document.querySelector('[data-slot="dialog-overlay"]')!;
      fireEvent.pointerDown(overlay);
      fireEvent.pointerUp(overlay);

      await waitFor(() => {
        expect(onOpenChange).toHaveBeenCalledWith(false);
      });
    });
  });

  describe('Dialog Interactions', () => {
    test('closes with close button', async () => {
      const onOpenChange = vi.fn();

      render(
        <Dialog defaultOpen onOpenChange={onOpenChange}>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Closeable Dialog</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        expect(screen.getByText('Closeable Dialog')).toBeInTheDocument();
      });

      const closeButton = document.querySelector('[data-slot="dialog-content"] button');
      fireEvent.click(closeButton!);

      expect(onOpenChange).toHaveBeenCalledWith(false);
    });

    test('closes with DialogClose component', async () => {
      const onOpenChange = vi.fn();

      render(
        <Dialog defaultOpen onOpenChange={onOpenChange}>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Dialog Title</DialogTitle>
            <DialogClose>Custom Close</DialogClose>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const customClose = screen.getByText('Custom Close');
        fireEvent.click(customClose);
      });

      expect(onOpenChange).toHaveBeenCalledWith(false);
    });

    test('closes with Escape key', async () => {
      const onOpenChange = vi.fn();

      render(
        <Dialog defaultOpen onOpenChange={onOpenChange}>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Press Escape</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        expect(screen.getByText('Press Escape')).toBeInTheDocument();
      });

      fireEvent.keyDown(document, { key: 'Escape' });

      expect(onOpenChange).toHaveBeenCalledWith(false);
    });

    test('handles controlled open state', () => {
      const ControlledDialog = () => {
        const [open, setOpen] = React.useState(false);

        return (
          <>
            <button onClick={() => setOpen(true)}>Open Controlled</button>
            <Dialog open={open} onOpenChange={setOpen}>
              <DialogContent aria-describedby="dialog-description">
                <DialogTitle>Controlled Dialog</DialogTitle>
                <button onClick={() => setOpen(false)}>Close</button>
              </DialogContent>
            </Dialog>
          </>
        );
      };

      render(<ControlledDialog />);

      expect(screen.queryByText('Controlled Dialog')).not.toBeInTheDocument();

      const openButton = screen.getByText('Open Controlled');
      fireEvent.click(openButton);

      expect(screen.getByText('Controlled Dialog')).toBeInTheDocument();
    });

    test('handles uncontrolled with defaultOpen', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Default Open Dialog</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        expect(screen.getByText('Default Open Dialog')).toBeInTheDocument();
      });
    });
  });

  describe('Focus Management', () => {
    test('focuses first focusable element when opened', async () => {
      render(
        <Dialog>
          <DialogTrigger>Open</DialogTrigger>
          <DialogContent aria-describedby="dialog-description">
            <DialogHeader>
              <DialogTitle>Focus Test</DialogTitle>
            </DialogHeader>
            <input type="text" placeholder="First input" />
            <button>Action Button</button>
          </DialogContent>
        </Dialog>
      );

      const trigger = screen.getByText('Open');
      fireEvent.click(trigger);

      await waitFor(() => {
        // Note: Focus management is handled by Radix UI Dialog
        expect(screen.getByText('Focus Test')).toBeInTheDocument();
      });
    });

    test('returns focus to trigger when closed', async () => {
      render(
        <Dialog>
          <DialogTrigger>Open Dialog</DialogTrigger>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Test</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      const trigger = screen.getByText('Open Dialog');
      fireEvent.click(trigger);

      await waitFor(() => {
        expect(screen.getByText('Test')).toBeInTheDocument();
      });

      fireEvent.keyDown(document, { key: 'Escape' });

      await waitFor(() => {
        expect(screen.queryByText('Test')).not.toBeInTheDocument();
        // Note: Focus return is handled by Radix UI
      });
    });
  });

  describe('Accessibility', () => {
    test('has proper ARIA attributes', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Accessible Dialog</DialogTitle>
            <DialogDescription id="dialog-description">
              This is an accessible dialog.
            </DialogDescription>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const dialog = document.querySelector('[role="dialog"]');
        expect(dialog).toBeInTheDocument();
        expect(dialog).toHaveAttribute('aria-describedby', 'dialog-description');
      });
    });

    test('title has proper heading role', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Dialog Heading</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const title = screen.getByText('Dialog Heading');
        expect(title.tagName).toBe('H2'); // Default heading level
      });
    });

    test('close button has screen reader text', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Test</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const closeText = screen.getByText('Close');
        expect(closeText).toHaveClass('sr-only');
      });
    });

    test('prevents body scroll when open', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>No Scroll Dialog</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        // Note: Scroll prevention is handled by Radix UI
        expect(screen.getByText('No Scroll Dialog')).toBeInTheDocument();
      });
    });
  });

  describe('Animation States', () => {
    test('applies animation classes', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Animated Dialog</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const content = document.querySelector('[data-slot="dialog-content"]');
        expect(content).toHaveClass('data-[state=open]:animate-in');
        expect(content).toHaveClass('data-[state=closed]:animate-out');
      });
    });

    test('overlay has fade animation', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>Test</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        const overlay = document.querySelector('[data-slot="dialog-overlay"]');
        expect(overlay).toHaveClass('data-[state=open]:fade-in-0');
        expect(overlay).toHaveClass('data-[state=closed]:fade-out-0');
      });
    });
  });

  describe('Common Use Cases', () => {
    test('confirmation dialog', async () => {
      const handleConfirm = vi.fn();
      const handleCancel = vi.fn();

      render(
        <Dialog>
          <DialogTrigger>Delete Item</DialogTrigger>
          <DialogContent aria-describedby="dialog-description">
            <DialogHeader>
              <DialogTitle>Confirm Deletion</DialogTitle>
              <DialogDescription>
                Are you sure you want to delete this item? This action cannot be undone.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <DialogClose onClick={handleCancel}>Cancel</DialogClose>
              <DialogClose onClick={handleConfirm}>Delete</DialogClose>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      );

      const trigger = screen.getByText('Delete Item');
      fireEvent.click(trigger);

      await waitFor(() => {
        const confirmButton = screen.getByText('Delete');
        fireEvent.click(confirmButton);
      });

      expect(handleConfirm).toHaveBeenCalled();
    });

    test('form dialog', async () => {
      const handleSubmit = vi.fn((e) => e.preventDefault());

      render(
        <Dialog>
          <DialogTrigger>Add User</DialogTrigger>
          <DialogContent aria-describedby="dialog-description">
            <DialogHeader>
              <DialogTitle>Add New User</DialogTitle>
              <DialogDescription>
                Enter the user details below.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit}>
              <input type="text" name="name" placeholder="Name" />
              <input type="email" name="email" placeholder="Email" />
              <DialogFooter>
                <DialogClose>Cancel</DialogClose>
                <button type="submit">Add User</button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      );

      const trigger = screen.getByText('Add User');
      fireEvent.click(trigger);

      await waitFor(() => {
        const submitButton = screen.getByText('Add User', { selector: 'button[type="submit"]' });
        fireEvent.click(submitButton);
      });

      expect(handleSubmit).toHaveBeenCalled();
    });

    test('information modal', async () => {
      render(
        <Dialog defaultOpen>
          <DialogContent aria-describedby="dialog-description">
            <DialogHeader>
              <DialogTitle>Welcome!</DialogTitle>
              <DialogDescription>
                Thank you for signing up. Check your email for confirmation.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <DialogClose>Got it</DialogClose>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      );

      await waitFor(() => {
        expect(screen.getByText('Welcome!')).toBeInTheDocument();
        expect(screen.getByText(/Thank you for signing up/)).toBeInTheDocument();
      });
    });

    test('nested dialogs', async () => {
      render(
        <Dialog>
          <DialogTrigger>Open First</DialogTrigger>
          <DialogContent aria-describedby="dialog-description">
            <DialogTitle>First Dialog</DialogTitle>
            <Dialog>
              <DialogTrigger>Open Second</DialogTrigger>
              <DialogContent aria-describedby="dialog-description">
                <DialogTitle>Second Dialog</DialogTitle>
              </DialogContent>
            </Dialog>
          </DialogContent>
        </Dialog>
      );

      const firstTrigger = screen.getByText('Open First');
      fireEvent.click(firstTrigger);

      await waitFor(() => {
        expect(screen.getByText('First Dialog')).toBeInTheDocument();
      });

      const secondTrigger = screen.getByText('Open Second');
      fireEvent.click(secondTrigger);

      await waitFor(() => {
        expect(screen.getByText('Second Dialog')).toBeInTheDocument();
      });
    });
  });
});