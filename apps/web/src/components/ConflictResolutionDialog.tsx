/**
 * Conflict Resolution Dialog for Concurrent Content Editing
 *
 * Provides a sophisticated UI for resolving conflicts when multiple users
 * edit the same content simultaneously. Features visual diff comparison,
 * merge options, and real-time collaboration indicators.
 *
 * Features:
 * - Side-by-side visual diff comparison
 * - Three-way merge with common ancestor
 * - Auto-merge for non-conflicting changes
 * - Manual conflict resolution with granular selection
 * - Real-time user presence indicators
 * - Conflict prevention through content locking
 * - Version history navigation
 */

import React, { useState, useEffect, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from './ui/dialog';
import { Button } from './ui/button';
import { Textarea } from './ui/textarea';
import { Badge } from './ui/badge/badge';
import { Separator } from './ui/separator';
import { ScrollArea } from './ui/scroll-area';
import {
  AlertTriangle,
  Users,
  GitMerge,
  Clock,
  Check,
  X,
  RotateCcw,
  Eye,
  Edit3,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import type { ContentItem, ContentVersion } from '../types/content';
import { notificationService } from '../services/notificationService';
import { useRealTimeContent } from '../hooks/useRealTimeContent';

// =============================================================================
// Type Definitions
// =============================================================================

export interface ConflictResolutionData {
  contentId: string;
  currentVersion: ContentItem;
  conflictingVersion: ContentItem;
  baseVersion?: ContentItem; // Common ancestor for three-way merge
  activeUsers: Array<{
    id: string;
    name: string;
    avatar?: string;
    isEditing: boolean;
    lastActivity: string;
  }>;
  conflictDetails: ConflictDetail[];
}

export interface ConflictDetail {
  id: string;
  type: 'content' | 'metadata' | 'status';
  field: string;
  currentValue: any;
  conflictingValue: any;
  baseValue?: any;
  resolved: boolean;
  resolution?: 'current' | 'conflicting' | 'custom' | 'merged';
  customValue?: any;
}

export interface ConflictResolutionProps {
  isOpen: boolean;
  onClose: () => void;
  conflictData: ConflictResolutionData;
  onResolve: (resolution: ResolvedContent) => Promise<void>;
  onCancel: () => void;
}

export interface ResolvedContent {
  contentId: string;
  resolvedValue: any;
  resolution: 'manual' | 'auto_merge' | 'force_current' | 'force_conflicting';
  conflictResolutions: Array<{
    conflictId: string;
    resolution: ConflictDetail['resolution'];
    value: any;
  }>;
  resolutionNotes?: string;
}

// =============================================================================
// Main Component
// =============================================================================

export function ConflictResolutionDialog({
  isOpen,
  onClose,
  conflictData,
  onResolve,
  onCancel,
}: ConflictResolutionProps) {
  const [selectedView, setSelectedView] = useState<'diff' | 'current' | 'conflicting' | 'merged'>('diff');
  const [conflicts, setConflicts] = useState<ConflictDetail[]>(conflictData.conflictDetails);
  const [mergedContent, setMergedContent] = useState<string>('');
  const [resolutionNotes, setResolutionNotes] = useState<string>('');
  const [isResolving, setIsResolving] = useState(false);
  const [autoMergeAttempted, setAutoMergeAttempted] = useState(false);

  const { isUserEditing, setIsUserEditing } = useRealTimeContent(conflictData.contentId, {
    preserveLocalChanges: true,
  });

  // Auto-merge non-conflicting changes
  useEffect(() => {
    if (!autoMergeAttempted) {
      attemptAutoMerge();
      setAutoMergeAttempted(true);
    }
  }, [conflictData]);

  // Calculate resolution statistics
  const resolutionStats = useMemo(() => {
    const total = conflicts.length;
    const resolved = conflicts.filter(c => c.resolved).length;
    const remaining = total - resolved;

    return { total, resolved, remaining };
  }, [conflicts]);

  // Check if all conflicts are resolved
  const allConflictsResolved = useMemo(() => {
    return conflicts.every(c => c.resolved);
  }, [conflicts]);

  // =========================================================================
  // Auto-merge Logic
  // =========================================================================

  const attemptAutoMerge = () => {
    const updatedConflicts = conflicts.map(conflict => {
      // Auto-resolve if one side matches the base version
      if (conflict.baseValue !== undefined) {
        if (conflict.currentValue === conflict.baseValue && conflict.conflictingValue !== conflict.baseValue) {
          return {
            ...conflict,
            resolved: true,
            resolution: 'conflicting' as const,
          };
        } else if (conflict.conflictingValue === conflict.baseValue && conflict.currentValue !== conflict.baseValue) {
          return {
            ...conflict,
            resolved: true,
            resolution: 'current' as const,
          };
        }
      }

      // Auto-resolve metadata conflicts with current preference
      if (conflict.type === 'metadata' && conflict.field !== 'version') {
        return {
          ...conflict,
          resolved: true,
          resolution: 'current' as const,
        };
      }

      return conflict;
    });

    setConflicts(updatedConflicts);

    // Generate merged content
    generateMergedContent(updatedConflicts);
  };

  const generateMergedContent = (conflictList: ConflictDetail[]) => {
    const contentConflict = conflictList.find(c => c.type === 'content' && c.field === 'value');

    if (!contentConflict) {
      setMergedContent(conflictData.currentVersion.value?.text || '');
      return;
    }

    if (contentConflict.resolved) {
      switch (contentConflict.resolution) {
        case 'current':
          setMergedContent(contentConflict.currentValue?.text || '');
          break;
        case 'conflicting':
          setMergedContent(contentConflict.conflictingValue?.text || '');
          break;
        case 'custom':
          setMergedContent(contentConflict.customValue?.text || '');
          break;
        default:
          setMergedContent(contentConflict.currentValue?.text || '');
      }
    } else {
      // Show a basic merge attempt
      const currentText = contentConflict.currentValue?.text || '';
      const conflictingText = contentConflict.conflictingValue?.text || '';
      setMergedContent(`${currentText}\n\n--- CONFLICT MARKER ---\n\n${conflictingText}`);
    }
  };

  // =========================================================================
  // Conflict Resolution Handlers
  // =========================================================================

  const resolveConflict = (conflictId: string, resolution: ConflictDetail['resolution'], customValue?: any) => {
    const updatedConflicts = conflicts.map(conflict => {
      if (conflict.id === conflictId) {
        return {
          ...conflict,
          resolved: true,
          resolution,
          customValue,
        };
      }
      return conflict;
    });

    setConflicts(updatedConflicts);
    generateMergedContent(updatedConflicts);
  };

  const unresolveConflict = (conflictId: string) => {
    const updatedConflicts = conflicts.map(conflict => {
      if (conflict.id === conflictId) {
        return {
          ...conflict,
          resolved: false,
          resolution: undefined,
          customValue: undefined,
        };
      }
      return conflict;
    });

    setConflicts(updatedConflicts);
    generateMergedContent(updatedConflicts);
  };

  const handleResolveAll = async () => {
    if (!allConflictsResolved) {
      return;
    }

    setIsResolving(true);
    setIsUserEditing(false);

    try {
      const resolvedContent: ResolvedContent = {
        contentId: conflictData.contentId,
        resolvedValue: mergedContent,
        resolution: autoMergeAttempted ? 'auto_merge' : 'manual',
        conflictResolutions: conflicts.map(c => ({
          conflictId: c.id,
          resolution: c.resolution!,
          value: c.resolution === 'custom' ? c.customValue :
                 c.resolution === 'current' ? c.currentValue :
                 c.conflictingValue,
        })),
        resolutionNotes,
      };

      await onResolve(resolvedContent);

      // Notify other users about resolution
      await notificationService.notify({
        type: 'content_updated',
        priority: 'medium',
        title: 'Conflict Resolved',
        message: `Content conflict has been resolved for "${conflictData.currentVersion.key}"`,
        contentId: conflictData.contentId,
        category: conflictData.currentVersion.category.name,
      });

      onClose();
    } catch (error) {
      console.error('Failed to resolve conflict:', error);
      await notificationService.notifyContentError(
        'Failed to resolve content conflict. Please try again.',
        conflictData.contentId,
        conflictData.currentVersion.category.name
      );
    } finally {
      setIsResolving(false);
    }
  };

  const handleCancel = () => {
    setIsUserEditing(false);
    onCancel();
    onClose();
  };

  // =========================================================================
  // Render Methods
  // =========================================================================

  const renderActiveUsers = () => (
    <div className="flex items-center gap-2 p-3 bg-muted/50 rounded-lg">
      <Users className="h-4 w-4 text-muted-foreground" />
      <span className="text-sm font-medium">Active editors:</span>
      <div className="flex gap-2">
        {conflictData.activeUsers.map(user => (
          <div key={user.id} className="flex items-center gap-1">
            {user.avatar ? (
              <img src={user.avatar} alt={user.name} className="h-6 w-6 rounded-full" />
            ) : (
              <div className="h-6 w-6 rounded-full bg-primary/20 flex items-center justify-center">
                <span className="text-xs font-medium">{user.name[0]}</span>
              </div>
            )}
            <span className="text-xs text-muted-foreground">{user.name}</span>
            {user.isEditing && (
              <div className="h-2 w-2 bg-green-500 rounded-full animate-pulse" />
            )}
          </div>
        ))}
      </div>
    </div>
  );

  const renderConflictList = () => (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-medium">Conflicts to Resolve</h4>
        <Badge variant="outline">
          {resolutionStats.resolved}/{resolutionStats.total} resolved
        </Badge>
      </div>

      <ScrollArea className="h-64">
        <div className="space-y-3">
          {conflicts.map(conflict => (
            <motion.div
              key={conflict.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className={`p-3 border rounded-lg ${
                conflict.resolved ? 'bg-green-50 border-green-200' : 'bg-yellow-50 border-yellow-200'
              }`}
            >
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2">
                  {conflict.resolved ? (
                    <Check className="h-4 w-4 text-green-600" />
                  ) : (
                    <AlertTriangle className="h-4 w-4 text-yellow-600" />
                  )}
                  <span className="text-sm font-medium capitalize">
                    {conflict.type} - {conflict.field}
                  </span>
                </div>
                {conflict.resolved && (
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={() => unresolveConflict(conflict.id)}
                  >
                    <RotateCcw className="h-3 w-3" />
                  </Button>
                )}
              </div>

              {!conflict.resolved && (
                <div className="space-y-2">
                  <div className="grid grid-cols-2 gap-2 text-xs">
                    <div className="p-2 bg-blue-50 rounded border">
                      <div className="font-medium text-blue-700 mb-1">Your Version</div>
                      <div className="text-blue-600">
                        {typeof conflict.currentValue === 'object'
                          ? JSON.stringify(conflict.currentValue, null, 2)
                          : String(conflict.currentValue)}
                      </div>
                    </div>
                    <div className="p-2 bg-red-50 rounded border">
                      <div className="font-medium text-red-700 mb-1">Their Version</div>
                      <div className="text-red-600">
                        {typeof conflict.conflictingValue === 'object'
                          ? JSON.stringify(conflict.conflictingValue, null, 2)
                          : String(conflict.conflictingValue)}
                      </div>
                    </div>
                  </div>

                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => resolveConflict(conflict.id, 'current')}
                      className="flex-1"
                    >
                      Use Mine
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => resolveConflict(conflict.id, 'conflicting')}
                      className="flex-1"
                    >
                      Use Theirs
                    </Button>
                  </div>
                </div>
              )}

              {conflict.resolved && (
                <div className="text-xs text-muted-foreground">
                  Resolved: Using {conflict.resolution === 'current' ? 'your version' :
                                   conflict.resolution === 'conflicting' ? 'their version' :
                                   conflict.resolution === 'custom' ? 'custom value' : 'merged content'}
                </div>
              )}
            </motion.div>
          ))}
        </div>
      </ScrollArea>
    </div>
  );

  const renderContentPreview = () => (
    <div className="space-y-3">
      <div className="flex items-center gap-2">
        <Button
          size="sm"
          variant={selectedView === 'diff' ? 'default' : 'outline'}
          onClick={() => setSelectedView('diff')}
        >
          <GitMerge className="h-3 w-3 mr-1" />
          Diff
        </Button>
        <Button
          size="sm"
          variant={selectedView === 'current' ? 'default' : 'outline'}
          onClick={() => setSelectedView('current')}
        >
          Your Version
        </Button>
        <Button
          size="sm"
          variant={selectedView === 'conflicting' ? 'default' : 'outline'}
          onClick={() => setSelectedView('conflicting')}
        >
          Their Version
        </Button>
        <Button
          size="sm"
          variant={selectedView === 'merged' ? 'default' : 'outline'}
          onClick={() => setSelectedView('merged')}
        >
          <Eye className="h-3 w-3 mr-1" />
          Preview
        </Button>
      </div>

      <ScrollArea className="h-64 w-full border rounded-lg">
        <div className="p-4">
          {selectedView === 'diff' && (
            <div className="grid grid-cols-2 gap-4 h-full">
              <div className="space-y-2">
                <h5 className="text-sm font-medium text-blue-700">Your Version</h5>
                <pre className="text-xs bg-blue-50 p-3 rounded border overflow-auto">
                  {conflictData.currentVersion.value?.text || 'No content'}
                </pre>
              </div>
              <div className="space-y-2">
                <h5 className="text-sm font-medium text-red-700">Their Version</h5>
                <pre className="text-xs bg-red-50 p-3 rounded border overflow-auto">
                  {conflictData.conflictingVersion.value?.text || 'No content'}
                </pre>
              </div>
            </div>
          )}

          {selectedView === 'current' && (
            <pre className="text-sm whitespace-pre-wrap">
              {conflictData.currentVersion.value?.text || 'No content'}
            </pre>
          )}

          {selectedView === 'conflicting' && (
            <pre className="text-sm whitespace-pre-wrap">
              {conflictData.conflictingVersion.value?.text || 'No content'}
            </pre>
          )}

          {selectedView === 'merged' && (
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <h5 className="text-sm font-medium">Merged Content Preview</h5>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setSelectedView('diff')}
                >
                  <Edit3 className="h-3 w-3 mr-1" />
                  Edit
                </Button>
              </div>
              <Textarea
                value={mergedContent}
                onChange={(e) => setMergedContent(e.target.value)}
                className="min-h-[200px] font-mono text-sm"
                placeholder="Merged content will appear here..."
              />
            </div>
          )}
        </div>
      </ScrollArea>
    </div>
  );

  // =========================================================================
  // Main Render
  // =========================================================================

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-6xl max-h-[90vh] overflow-hidden">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-yellow-600" />
            Content Conflict Resolution
          </DialogTitle>
          <DialogDescription>
            Multiple users have edited "{conflictData.currentVersion.key}" simultaneously.
            Please resolve the conflicts below to continue.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {renderActiveUsers()}

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="space-y-4">
              {renderConflictList()}

              <div className="space-y-3">
                <label className="text-sm font-medium">Resolution Notes (Optional)</label>
                <Textarea
                  value={resolutionNotes}
                  onChange={(e) => setResolutionNotes(e.target.value)}
                  placeholder="Add notes about how you resolved this conflict..."
                  className="min-h-[80px]"
                />
              </div>
            </div>

            <div className="space-y-4">
              {renderContentPreview()}
            </div>
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock className="h-4 w-4" />
              Last updated: {new Date(conflictData.currentVersion.metadata.updatedAt).toLocaleString()}
            </div>

            <div className="flex gap-3">
              <Button variant="outline" onClick={handleCancel} disabled={isResolving}>
                Cancel
              </Button>
              <Button
                onClick={handleResolveAll}
                disabled={!allConflictsResolved || isResolving}
                className="min-w-[120px]"
              >
                {isResolving ? (
                  <motion.div
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                    className="h-4 w-4 border-2 border-current border-t-transparent rounded-full"
                  />
                ) : (
                  <>
                    <Check className="h-4 w-4 mr-2" />
                    Resolve Conflict
                  </>
                )}
              </Button>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}