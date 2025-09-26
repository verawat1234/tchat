// T050 - SpreadsheetMessage component with data tables
/**
 * SpreadsheetMessage Component
 * Displays spreadsheet data with sorting, filtering, and editing capabilities
 * Supports large datasets, real-time updates, and collaborative editing
 */

import React, { useState, useCallback, useMemo, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { SpreadsheetContent, SpreadsheetCell, SpreadsheetColumn, SortDirection } from '../../types/SpreadsheetContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Input } from '../ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Separator } from '../ui/separator';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../ui/table';
import {
  Download,
  Upload,
  Plus,
  Minus,
  Edit3,
  Save,
  X,
  Check,
  ArrowUpDown,
  ArrowUp,
  ArrowDown,
  Filter,
  Search,
  MoreHorizontal,
  FileSpreadsheet,
  Share,
  Eye,
  Lock,
  Unlock,
  RefreshCw,
  Users,
  Calculator,
  TrendingUp,
  BarChart3,
  Grid3X3,
  Maximize,
  Minimize
} from 'lucide-react';

// Component Props Interface
interface SpreadsheetMessageProps {
  message: MessageData & { content: SpreadsheetContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onCellUpdate?: (rowId: string, columnId: string, value: string) => void;
  onRowAdd?: (sheetId: string) => void;
  onRowDelete?: (rowId: string) => void;
  onShare?: (sheetId: string) => void;
  onDownload?: (sheetId: string, format: 'csv' | 'xlsx') => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  allowEditing?: boolean;
  maxDisplayRows?: number;
  performanceMode?: boolean;
}

// Animation Variants
const spreadsheetVariants = {
  initial: { opacity: 0, scale: 0.98, y: 20 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.98, y: -20 }
};

const tableVariants = {
  initial: { opacity: 0, y: 10 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -10 }
};

const rowVariants = {
  initial: { opacity: 0, x: -10 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 10 }
};

export const SpreadsheetMessage: React.FC<SpreadsheetMessageProps> = ({
  message,
  onInteraction,
  onCellUpdate,
  onRowAdd,
  onRowDelete,
  onShare,
  onDownload,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  allowEditing = true,
  maxDisplayRows = 10,
  performanceMode = false
}) => {
  const spreadsheetRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Spreadsheet state
  const [sortColumn, setSortColumn] = useState<string | null>(null);
  const [sortDirection, setSortDirection] = useState<SortDirection>(SortDirection.ASC);
  const [filterColumn, setFilterColumn] = useState<string>('');
  const [filterValue, setFilterValue] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [editingCell, setEditingCell] = useState<{ rowId: string; columnId: string } | null>(null);
  const [editValue, setEditValue] = useState<string>('');
  const [isExpanded, setIsExpanded] = useState(false);

  // Process and filter data
  const processedData = useMemo(() => {
    let filteredRows = [...content.rows];

    // Apply search filter
    if (searchQuery.trim()) {
      filteredRows = filteredRows.filter(row =>
        content.columns.some(column => {
          const cell = row.cells.find(c => c.columnId === column.id);
          return cell?.value.toLowerCase().includes(searchQuery.toLowerCase());
        })
      );
    }

    // Apply column filter
    if (filterColumn && filterValue.trim()) {
      filteredRows = filteredRows.filter(row => {
        const cell = row.cells.find(c => c.columnId === filterColumn);
        return cell?.value.toLowerCase().includes(filterValue.toLowerCase());
      });
    }

    // Apply sorting
    if (sortColumn) {
      const column = content.columns.find(c => c.id === sortColumn);
      if (column) {
        filteredRows.sort((a, b) => {
          const cellA = a.cells.find(c => c.columnId === sortColumn);
          const cellB = b.cells.find(c => c.columnId === sortColumn);

          const valueA = cellA?.value || '';
          const valueB = cellB?.value || '';

          // Numeric sorting for number columns
          if (column.type === 'number') {
            const numA = parseFloat(valueA) || 0;
            const numB = parseFloat(valueB) || 0;
            return sortDirection === SortDirection.ASC ? numA - numB : numB - numA;
          }

          // Date sorting for date columns
          if (column.type === 'date') {
            const dateA = new Date(valueA).getTime() || 0;
            const dateB = new Date(valueB).getTime() || 0;
            return sortDirection === SortDirection.ASC ? dateA - dateB : dateB - dateA;
          }

          // String sorting for text columns
          const comparison = valueA.localeCompare(valueB);
          return sortDirection === SortDirection.ASC ? comparison : -comparison;
        });
      }
    }

    return filteredRows;
  }, [content.rows, content.columns, searchQuery, filterColumn, filterValue, sortColumn, sortDirection]);

  // Get visible rows based on expansion and max display
  const visibleRows = useMemo(() => {
    return isExpanded ? processedData : processedData.slice(0, maxDisplayRows);
  }, [processedData, isExpanded, maxDisplayRows]);

  // Handle column sorting
  const handleSort = useCallback((columnId: string) => {
    if (sortColumn === columnId) {
      // Toggle sort direction or remove sorting
      if (sortDirection === SortDirection.ASC) {
        setSortDirection(SortDirection.DESC);
      } else {
        setSortColumn(null);
        setSortDirection(SortDirection.ASC);
      }
    } else {
      setSortColumn(columnId);
      setSortDirection(SortDirection.ASC);
    }
  }, [sortColumn, sortDirection]);

  // Handle cell editing
  const handleCellClick = useCallback((rowId: string, columnId: string, currentValue: string) => {
    if (readonly || !allowEditing) return;

    const column = content.columns.find(c => c.id === columnId);
    if (column?.readonly) return;

    setEditingCell({ rowId, columnId });
    setEditValue(currentValue);
  }, [readonly, allowEditing, content.columns]);

  // Handle cell save
  const handleCellSave = useCallback(() => {
    if (!editingCell) return;

    if (onCellUpdate) {
      onCellUpdate(editingCell.rowId, editingCell.columnId, editValue);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'cell_update',
        data: {
          sheetId: content.id,
          rowId: editingCell.rowId,
          columnId: editingCell.columnId,
          value: editValue
        },
        userId: 'current-user',
        timestamp: new Date()
      });
    }

    setEditingCell(null);
    setEditValue('');
  }, [editingCell, editValue, onCellUpdate, onInteraction, message.id, content.id]);

  // Handle cell cancel
  const handleCellCancel = useCallback(() => {
    setEditingCell(null);
    setEditValue('');
  }, []);

  // Handle add row
  const handleAddRow = useCallback(() => {
    if (readonly || !allowEditing) return;

    if (onRowAdd) {
      onRowAdd(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'row_add',
        data: { sheetId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, allowEditing, onRowAdd, content.id, onInteraction, message.id]);

  // Handle share
  const handleShare = useCallback(() => {
    if (readonly) return;

    if (onShare) {
      onShare(content.id);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'spreadsheet_share',
        data: { sheetId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, onShare, content.id, onInteraction, message.id]);

  // Handle download
  const handleDownload = useCallback((format: 'csv' | 'xlsx') => {
    if (readonly) return;

    if (onDownload) {
      onDownload(content.id, format);
    }

    if (onInteraction) {
      onInteraction({
        messageId: message.id,
        interactionType: 'spreadsheet_download',
        data: { sheetId: content.id, format },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, onDownload, content.id, onInteraction, message.id]);

  // Get cell display value
  const getCellDisplayValue = useCallback((cell: SpreadsheetCell, column: SpreadsheetColumn) => {
    if (!cell.value) return '';

    switch (column.type) {
      case 'currency':
        const amount = parseFloat(cell.value);
        return isNaN(amount) ? cell.value : new Intl.NumberFormat('en-US', {
          style: 'currency',
          currency: 'USD'
        }).format(amount);

      case 'percentage':
        const percent = parseFloat(cell.value);
        return isNaN(percent) ? cell.value : `${percent}%`;

      case 'date':
        const date = new Date(cell.value);
        return isNaN(date.getTime()) ? cell.value : date.toLocaleDateString();

      default:
        return cell.value;
    }
  }, []);

  // Get sort icon
  const getSortIcon = useCallback((columnId: string) => {
    if (sortColumn !== columnId) return <ArrowUpDown className="w-3 h-3" />;
    return sortDirection === SortDirection.ASC ?
      <ArrowUp className="w-3 h-3" /> :
      <ArrowDown className="w-3 h-3" />;
  }, [sortColumn, sortDirection]);

  // Performance optimization
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: spreadsheetVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={spreadsheetRef}
        className={cn(
          "spreadsheet-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`spreadsheet-message-${message.id}`}
        data-sheet-id={content.id}
        role="article"
        aria-label={`Spreadsheet: ${content.title} with ${content.rows.length} rows`}
      >
        <Card className="spreadsheet-card">
          {/* Header */}
          <CardHeader className="space-y-3">
            <div className="flex items-start justify-between gap-3">
              <div className="flex items-center gap-3 min-w-0 flex-1">
                {showAvatar && (
                  <motion.div
                    initial={performanceMode ? {} : { scale: 0.8, opacity: 0 }}
                    animate={performanceMode ? {} : { scale: 1, opacity: 1 }}
                    transition={{ delay: 0.1 }}
                  >
                    <Avatar className={cn(compactMode ? "w-8 h-8" : "w-10 h-10")}>
                      <AvatarImage src={`/avatars/${message.senderName.toLowerCase()}.png`} />
                      <AvatarFallback>
                        {message.senderName.substring(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                  </motion.div>
                )}

                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="font-semibold text-foreground truncate">
                      {message.senderName}
                    </span>
                    {message.isOwn && (
                      <Badge variant="secondary" className="text-xs">You</Badge>
                    )}
                  </div>
                  {showTimestamp && (
                    <p className="text-xs text-muted-foreground mt-1">
                      Shared a spreadsheet • {message.timestamp.toLocaleDateString()}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-2">
                <Badge variant="outline" className="text-xs">
                  <Grid3X3 className="w-3 h-3 mr-1" />
                  {content.rows.length} rows
                </Badge>

                {content.isLocked && (
                  <Badge variant="secondary" className="text-xs">
                    <Lock className="w-3 h-3 mr-1" />
                    Locked
                  </Badge>
                )}

                {content.allowsCollaboration && (
                  <Badge variant="outline" className="text-xs">
                    <Users className="w-3 h-3 mr-1" />
                    Shared
                  </Badge>
                )}
              </div>
            </div>

            {/* Title and description */}
            <div className="space-y-2">
              <h3 className="font-semibold text-lg text-foreground leading-tight">
                {content.title}
              </h3>

              {content.description && (
                <p className="text-sm text-muted-foreground">
                  {content.description}
                </p>
              )}
            </div>

            {/* Controls */}
            <div className="flex items-center justify-between gap-4 flex-wrap">
              <div className="flex items-center gap-2 flex-wrap">
                {/* Search */}
                <div className="relative">
                  <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 w-3 h-3 text-muted-foreground" />
                  <Input
                    placeholder="Search..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-7 h-8 w-32 text-xs"
                  />
                </div>

                {/* Filter */}
                <div className="flex items-center gap-1">
                  <Select value={filterColumn} onValueChange={setFilterColumn}>
                    <SelectTrigger className="h-8 w-24 text-xs">
                      <SelectValue placeholder="Filter" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">No filter</SelectItem>
                      {content.columns.map(column => (
                        <SelectItem key={column.id} value={column.id}>
                          {column.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>

                  {filterColumn && (
                    <Input
                      placeholder="Value..."
                      value={filterValue}
                      onChange={(e) => setFilterValue(e.target.value)}
                      className="h-8 w-20 text-xs"
                    />
                  )}
                </div>
              </div>

              <div className="flex items-center gap-1">
                {allowEditing && !readonly && (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleAddRow}
                        className="h-8 w-8 p-0"
                      >
                        <Plus className="w-4 h-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Add row</TooltipContent>
                  </Tooltip>
                )}

                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setIsExpanded(!isExpanded)}
                      className="h-8 w-8 p-0"
                    >
                      {isExpanded ? <Minimize className="w-4 h-4" /> : <Maximize className="w-4 h-4" />}
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    {isExpanded ? 'Show less' : 'Show all'}
                  </TooltipContent>
                </Tooltip>

                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={handleShare}
                      disabled={readonly}
                      className="h-8 w-8 p-0"
                    >
                      <Share className="w-4 h-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Share spreadsheet</TooltipContent>
                </Tooltip>

                <Select onValueChange={(value) => handleDownload(value as 'csv' | 'xlsx')}>
                  <SelectTrigger asChild>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          variant="ghost"
                          size="sm"
                          disabled={readonly}
                          className="h-8 w-8 p-0"
                        >
                          <Download className="w-4 h-4" />
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>Download</TooltipContent>
                    </Tooltip>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="csv">Download CSV</SelectItem>
                    <SelectItem value="xlsx">Download Excel</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </CardHeader>

          {/* Spreadsheet table */}
          <CardContent>
            <motion.div
              variants={performanceMode ? {} : tableVariants}
              initial={performanceMode ? {} : "initial"}
              animate={performanceMode ? {} : "animate"}
              className="relative border rounded-lg overflow-hidden"
            >
              <div className="overflow-x-auto max-h-96">
                <Table>
                  <TableHeader>
                    <TableRow className="bg-muted/50">
                      {content.columns.map((column) => (
                        <TableHead
                          key={column.id}
                          className={cn(
                            "relative px-3 py-2 text-xs font-medium",
                            !column.readonly && !readonly && "cursor-pointer hover:bg-muted"
                          )}
                          onClick={() => handleSort(column.id)}
                        >
                          <div className="flex items-center gap-2">
                            <span className="truncate">{column.name}</span>
                            {getSortIcon(column.id)}
                            {column.required && (
                              <span className="text-destructive">*</span>
                            )}
                          </div>
                        </TableHead>
                      ))}
                      {allowEditing && !readonly && (
                        <TableHead className="w-12"></TableHead>
                      )}
                    </TableRow>
                  </TableHeader>

                  <TableBody>
                    <AnimatePresence mode="popLayout">
                      {visibleRows.map((row, rowIndex) => (
                        <motion.tr
                          key={row.id}
                          variants={performanceMode ? {} : rowVariants}
                          initial={performanceMode ? {} : "initial"}
                          animate={performanceMode ? {} : "animate"}
                          exit={performanceMode ? {} : "exit"}
                          transition={{ delay: rowIndex * 0.02 }}
                          className="border-b transition-colors hover:bg-muted/50"
                        >
                          {content.columns.map((column) => {
                            const cell = row.cells.find(c => c.columnId === column.id);
                            const isEditing = editingCell?.rowId === row.id && editingCell?.columnId === column.id;

                            return (
                              <TableCell
                                key={`${row.id}-${column.id}`}
                                className={cn(
                                  "px-3 py-2 text-xs",
                                  !column.readonly && !readonly && allowEditing && "cursor-pointer hover:bg-muted/25",
                                  cell?.hasError && "bg-red-50 text-red-900"
                                )}
                                onClick={() => handleCellClick(row.id, column.id, cell?.value || '')}
                              >
                                {isEditing ? (
                                  <div className="flex items-center gap-1">
                                    <Input
                                      value={editValue}
                                      onChange={(e) => setEditValue(e.target.value)}
                                      onKeyDown={(e) => {
                                        if (e.key === 'Enter') handleCellSave();
                                        if (e.key === 'Escape') handleCellCancel();
                                      }}
                                      className="h-6 text-xs"
                                      autoFocus
                                    />
                                    <Button
                                      variant="ghost"
                                      size="sm"
                                      onClick={handleCellSave}
                                      className="h-6 w-6 p-0"
                                    >
                                      <Check className="w-3 h-3" />
                                    </Button>
                                    <Button
                                      variant="ghost"
                                      size="sm"
                                      onClick={handleCellCancel}
                                      className="h-6 w-6 p-0"
                                    >
                                      <X className="w-3 h-3" />
                                    </Button>
                                  </div>
                                ) : (
                                  <div className="flex items-center justify-between">
                                    <span className="truncate">
                                      {cell ? getCellDisplayValue(cell, column) : ''}
                                    </span>
                                    {cell?.hasFormula && (
                                      <Calculator className="w-3 h-3 text-muted-foreground flex-shrink-0 ml-1" />
                                    )}
                                  </div>
                                )}
                              </TableCell>
                            );
                          })}

                          {allowEditing && !readonly && (
                            <TableCell className="px-2 py-2">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => onRowDelete?.(row.id)}
                                className="h-6 w-6 p-0 text-destructive hover:text-destructive"
                              >
                                <Minus className="w-3 h-3" />
                              </Button>
                            </TableCell>
                          )}
                        </motion.tr>
                      ))}
                    </AnimatePresence>
                  </TableBody>
                </Table>
              </div>

              {/* Show more indicator */}
              {processedData.length > maxDisplayRows && !isExpanded && (
                <div className="absolute bottom-0 left-0 right-0 h-12 bg-gradient-to-t from-background to-transparent flex items-end justify-center pb-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setIsExpanded(true)}
                    className="gap-2"
                  >
                    <Eye className="w-3 h-3" />
                    Show {processedData.length - maxDisplayRows} more rows
                  </Button>
                </div>
              )}
            </motion.div>

            {/* Statistics */}
            {content.statistics && (
              <div className="mt-4 pt-4 border-t">
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                  {content.statistics.totalRows && (
                    <div>
                      <div className="text-2xl font-bold text-foreground">
                        {content.statistics.totalRows.toLocaleString()}
                      </div>
                      <div className="text-xs text-muted-foreground">Total Rows</div>
                    </div>
                  )}

                  {content.statistics.averageValue !== undefined && (
                    <div>
                      <div className="text-2xl font-bold text-foreground">
                        {content.statistics.averageValue.toFixed(2)}
                      </div>
                      <div className="text-xs text-muted-foreground">Average</div>
                    </div>
                  )}

                  {content.statistics.sum !== undefined && (
                    <div>
                      <div className="text-2xl font-bold text-foreground">
                        {content.statistics.sum.toLocaleString()}
                      </div>
                      <div className="text-xs text-muted-foreground">Sum</div>
                    </div>
                  )}

                  {content.statistics.lastUpdated && (
                    <div>
                      <div className="text-sm font-medium text-foreground">
                        {new Date(content.statistics.lastUpdated).toLocaleDateString()}
                      </div>
                      <div className="text-xs text-muted-foreground">Last Updated</div>
                    </div>
                  )}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Performance Debug Info */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            {visibleRows.length}/{processedData.length} rows | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedSpreadsheetMessage = React.memo(SpreadsheetMessage, (prevProps, nextProps) => {
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.allowEditing === nextProps.allowEditing &&
    prevProps.maxDisplayRows === nextProps.maxDisplayRows &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedSpreadsheetMessage.displayName = 'MemoizedSpreadsheetMessage';

export default SpreadsheetMessage;