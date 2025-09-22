import React from 'react';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger, DialogClose } from '@/components/ui/dialog';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger, DropdownMenuSub, DropdownMenuSubContent, DropdownMenuSubTrigger } from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

/**
 * Component Showcase Page for E2E Testing
 * This page displays all portal-based components for testing
 */
export default function ComponentsPage() {
  return (
    <div className="container mx-auto p-8 space-y-8">
      <h1 className="text-4xl font-bold mb-8">Component Showcase</h1>

      {/* Tooltip Section */}
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">Tooltips</h2>
        <div className="flex gap-4 flex-wrap">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button data-testid="tooltip-trigger-1" variant="outline">
                  Hover for tooltip
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>This is a simple tooltip</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button data-testid="tooltip-trigger-2" variant="outline">
                  Another tooltip
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>This is another tooltip with different content</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <TooltipProvider delayDuration={800}>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button data-testid="tooltip-with-delay" variant="outline">
                  Tooltip with delay
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>This tooltip appears after a delay</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      </section>

      {/* Dialog Section */}
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">Dialogs</h2>
        <div className="flex gap-4 flex-wrap">
          <Dialog>
            <DialogTrigger asChild>
              <Button data-testid="dialog-trigger">Open Dialog</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Dialog Title</DialogTitle>
                <DialogDescription id="dialog-description">
                  This is a dialog description that explains the purpose of this dialog.
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div>
                  <Label htmlFor="dialog-input">Input Field</Label>
                  <Input id="dialog-input" placeholder="Type something..." />
                </div>

                {/* Nested dropdown in dialog */}
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button data-testid="dropdown-trigger" variant="outline">
                      Options in Dialog
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent>
                    <DropdownMenuItem>Option 1</DropdownMenuItem>
                    <DropdownMenuItem>Option 2</DropdownMenuItem>
                    <DropdownMenuItem>Option 3</DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
              <DialogClose asChild>
                <Button data-testid="dialog-close" variant="ghost">
                  Close
                </Button>
              </DialogClose>
            </DialogContent>
          </Dialog>

          <Dialog>
            <DialogTrigger asChild>
              <Button data-testid="dialog-trigger-2" variant="outline">
                Another Dialog
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Second Dialog</DialogTitle>
                <DialogDescription>
                  This is a different dialog with different content.
                </DialogDescription>
              </DialogHeader>
              <p>Some content goes here...</p>
            </DialogContent>
          </Dialog>
        </div>
      </section>

      {/* Dropdown Section */}
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">Dropdowns</h2>
        <div className="flex gap-4 flex-wrap">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button data-testid="dropdown-trigger" variant="outline">
                Open Dropdown
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem>Profile</DropdownMenuItem>
              <DropdownMenuItem>Settings</DropdownMenuItem>
              <DropdownMenuItem>Billing</DropdownMenuItem>
              <DropdownMenuItem>Sign Out</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button data-testid="dropdown-trigger-submenu" variant="outline">
                Dropdown with Submenu
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem>New File</DropdownMenuItem>
              <DropdownMenuSub>
                <DropdownMenuSubTrigger>More Options</DropdownMenuSubTrigger>
                <DropdownMenuSubContent>
                  <DropdownMenuItem>Sub Option 1</DropdownMenuItem>
                  <DropdownMenuItem>Sub Option 2</DropdownMenuItem>
                  <DropdownMenuItem>Sub Option 3</DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuSub>
              <DropdownMenuItem>Save</DropdownMenuItem>
              <DropdownMenuItem>Exit</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </section>

      {/* Complex Interactions */}
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">Complex Interactions</h2>
        <div className="space-y-4">
          <div className="p-4 border rounded">
            <p className="mb-2">Multiple tooltips in a row:</p>
            <div className="flex gap-2">
              <TooltipProvider>
                {[1, 2, 3, 4, 5].map((i) => (
                  <Tooltip key={i}>
                    <TooltipTrigger asChild>
                      <Button
                        data-testid={`tooltip-trigger-row-${i}`}
                        variant="outline"
                        size="sm"
                      >
                        {i}
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Tooltip content for button {i}</p>
                    </TooltipContent>
                  </Tooltip>
                ))}
              </TooltipProvider>
            </div>
          </div>

          <div className="p-4 border rounded">
            <p className="mb-2">Scrollable area with tooltips:</p>
            <div className="h-32 overflow-y-auto border p-2 space-y-2">
              <TooltipProvider>
                {Array.from({ length: 10 }, (_, i) => (
                  <Tooltip key={i}>
                    <TooltipTrigger asChild>
                      <Button
                        data-testid={`tooltip-scroll-${i}`}
                        variant="ghost"
                        className="w-full justify-start"
                      >
                        Item {i + 1}
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Details for item {i + 1}</p>
                    </TooltipContent>
                  </Tooltip>
                ))}
              </TooltipProvider>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}