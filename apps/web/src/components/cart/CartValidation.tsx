/**
 * CartValidation - Cart Validation Display Component
 *
 * Displays cart validation status, issues, and recommendations.
 * Integrates with CartProvider for real-time validation.
 *
 * Features:
 * - Real-time validation status
 * - Issue categorization and severity
 * - Actionable recommendations
 * - Auto-refresh validation
 * - Accessibility compliance
 * - Responsive design
 */

import React, { useMemo } from 'react';
import { cn } from '../../lib/utils';
import { useCart } from './CartProvider';
import { TchatCard, TchatCardHeader, TchatCardContent, TchatCardFooter } from '../TchatCard';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge/badge';
import { Alert, AlertDescription, AlertTitle } from '../ui/alert';
import { Separator } from '../ui/separator';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '../ui/collapsible';
import {
  AlertCircle,
  AlertTriangle,
  CheckCircle,
  ChevronDown,
  ChevronUp,
  Info,
  RefreshCw,
  ShieldCheck,
  Package,
  CreditCard,
  Truck,
} from 'lucide-react';
import type { CartValidation as CartValidationType, CartValidationIssue } from '../../types/commerce';

// ===== Types =====

export interface CartValidationProps {
  /** Whether to show detailed issues */
  showDetails?: boolean;
  /** Whether to show validation summary */
  showSummary?: boolean;
  /** Whether to auto-refresh validation */
  autoRefresh?: boolean;
  /** Refresh interval in milliseconds */
  refreshInterval?: number;
  /** Whether to show actions */
  showActions?: boolean;
  /** Custom class name */
  className?: string;
  /** Compact display mode */
  compact?: boolean;
  /** Validation refresh handler */
  onRefresh?: () => void;
  /** Issue action handler */
  onResolveIssue?: (issue: CartValidationIssue) => void;
}

// ===== Helper Functions =====

/**
 * Get icon for validation issue severity
 */
const getIssueIcon = (severity: CartValidationIssue['severity']) => {
  switch (severity) {
    case 'error':
      return AlertCircle;
    case 'warning':
      return AlertTriangle;
    case 'info':
      return Info;
    default:
      return Info;
  }
};

/**
 * Get color variant for issue severity
 */
const getIssueVariant = (severity: CartValidationIssue['severity']) => {
  switch (severity) {
    case 'error':
      return 'destructive';
    case 'warning':
      return 'secondary';
    case 'info':
      return 'default';
    default:
      return 'default';
  }
};

/**
 * Format currency amount
 */
const formatCurrency = (amount: string, currency = 'USD'): string => {
  const numAmount = parseFloat(amount);
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
  }).format(numAmount);
};

/**
 * Get recommendations based on issue type
 */
const getIssueRecommendation = (issue: CartValidationIssue): string => {
  const issueType = issue.type.toLowerCase();

  if (issueType.includes('stock')) {
    return 'Update quantities or remove out-of-stock items';
  }
  if (issueType.includes('shipping')) {
    return 'Update shipping address or select different shipping method';
  }
  if (issueType.includes('payment')) {
    return 'Update payment method or billing information';
  }
  if (issueType.includes('coupon') || issueType.includes('discount')) {
    return 'Remove invalid coupon or check coupon requirements';
  }
  if (issueType.includes('minimum')) {
    return 'Add more items to meet minimum order requirements';
  }

  return 'Review and update cart items';
};

// ===== Component =====

/**
 * CartValidation - Displays cart validation status and issues
 *
 * @param showDetails - Show detailed issue breakdown
 * @param showSummary - Show validation summary
 * @param autoRefresh - Enable auto-refresh
 * @param refreshInterval - Refresh interval
 * @param showActions - Show action buttons
 * @param className - Custom styling
 * @param compact - Use compact layout
 * @param onRefresh - Refresh handler
 * @param onResolveIssue - Issue resolution handler
 */
export const CartValidation: React.FC<CartValidationProps> = ({
  showDetails = true,
  showSummary = true,
  autoRefresh = false,
  refreshInterval = 30000,
  showActions = true,
  className,
  compact = false,
  onRefresh,
  onResolveIssue,
}) => {
  const { cartValidation, isLoading, error } = useCart();

  // ===== Local State =====

  const [isDetailsExpanded, setIsDetailsExpanded] = React.useState(false);
  const [isRefreshing, setIsRefreshing] = React.useState(false);

  // ===== Computed Values =====

  const validationData = useMemo(() => {
    if (!cartValidation) return null;

    const { issues } = cartValidation;
    const errors = issues.filter(issue => issue.severity === 'error');
    const warnings = issues.filter(issue => issue.severity === 'warning');
    const infos = issues.filter(issue => issue.severity === 'info');

    const overallStatus = errors.length > 0 ? 'error' : warnings.length > 0 ? 'warning' : 'valid';

    return {
      ...cartValidation,
      errors,
      warnings,
      infos,
      overallStatus,
      hasIssues: issues.length > 0,
    };
  }, [cartValidation]);

  // ===== Effects =====

  React.useEffect(() => {
    if (!autoRefresh || !refreshInterval) return;

    const interval = setInterval(() => {
      // Validation will auto-refresh via RTK Query polling
    }, refreshInterval);

    return () => clearInterval(interval);
  }, [autoRefresh, refreshInterval]);

  // ===== Handlers =====

  const handleRefresh = React.useCallback(async () => {
    setIsRefreshing(true);
    try {
      onRefresh?.();
      // Wait for animation
      await new Promise(resolve => setTimeout(resolve, 1000));
    } finally {
      setIsRefreshing(false);
    }
  }, [onRefresh]);

  // ===== Render Helpers =====

  const renderValidationStatus = () => {
    if (!validationData) return null;

    const { overallStatus, errors, warnings, infos } = validationData;

    const statusConfig = {
      valid: {
        icon: CheckCircle,
        color: 'text-success-foreground',
        bgColor: 'bg-success/10',
        message: 'Cart is valid and ready for checkout',
      },
      warning: {
        icon: AlertTriangle,
        color: 'text-warning-foreground',
        bgColor: 'bg-warning/10',
        message: `Cart has ${warnings.length} warning${warnings.length !== 1 ? 's' : ''}`,
      },
      error: {
        icon: AlertCircle,
        color: 'text-destructive',
        bgColor: 'bg-destructive/10',
        message: `Cart has ${errors.length} error${errors.length !== 1 ? 's' : ''} that must be resolved`,
      },
    };

    const config = statusConfig[overallStatus];
    const Icon = config.icon;

    return (
      <div className={cn('flex items-center gap-3 p-3 rounded-lg', config.bgColor)}>
        <Icon className={cn('w-5 h-5', config.color)} />
        <div className="flex-1">
          <p className={cn('font-medium', config.color)}>
            {config.message}
          </p>
        </div>
        <div className="flex items-center gap-2">
          {errors.length > 0 && (
            <Badge variant="destructive" className="text-xs">
              {errors.length} error{errors.length !== 1 ? 's' : ''}
            </Badge>
          )}
          {warnings.length > 0 && (
            <Badge variant="secondary" className="text-xs">
              {warnings.length} warning{warnings.length !== 1 ? 's' : ''}
            </Badge>
          )}
          {infos.length > 0 && (
            <Badge variant="default" className="text-xs">
              {infos.length} info
            </Badge>
          )}
        </div>
      </div>
    );
  };

  const renderValidationSummary = () => {
    if (!showSummary || !validationData) return null;

    const { totalItems, totalValue, currency, estimatedShipping, estimatedTax, estimatedTotal } = validationData;

    return (
      <div className="space-y-3">
        <h4 className="font-medium text-sm">Validation Summary</h4>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div className="flex justify-between">
            <span className="text-muted-foreground">Items:</span>
            <span>{totalItems}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Subtotal:</span>
            <span>{formatCurrency(totalValue, currency)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Est. Shipping:</span>
            <span>{formatCurrency(estimatedShipping, currency)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Est. Tax:</span>
            <span>{formatCurrency(estimatedTax, currency)}</span>
          </div>
        </div>
        <Separator />
        <div className="flex justify-between font-medium">
          <span>Estimated Total:</span>
          <span>{formatCurrency(estimatedTotal, currency)}</span>
        </div>
      </div>
    );
  };

  const renderIssueItem = (issue: CartValidationIssue) => {
    const Icon = getIssueIcon(issue.severity);
    const variant = getIssueVariant(issue.severity) as any;
    const recommendation = getIssueRecommendation(issue);

    return (
      <Alert key={`${issue.type}-${issue.productId}`} variant={variant} className="p-4">
        <Icon className="h-4 w-4" />
        <AlertTitle className="text-sm font-medium">
          {issue.type.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}
        </AlertTitle>
        <AlertDescription className="text-sm mt-1">
          <div className="space-y-2">
            <p>{issue.message}</p>
            {recommendation && (
              <p className="text-xs text-muted-foreground italic">
                Recommendation: {recommendation}
              </p>
            )}
            {showActions && onResolveIssue && issue.severity === 'error' && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => onResolveIssue(issue)}
                className="mt-2"
              >
                Resolve Issue
              </Button>
            )}
          </div>
        </AlertDescription>
      </Alert>
    );
  };

  const renderIssuesList = () => {
    if (!showDetails || !validationData?.hasIssues) return null;

    const { errors, warnings, infos } = validationData;
    const allIssues = [...errors, ...warnings, ...infos];

    return (
      <Collapsible open={isDetailsExpanded} onOpenChange={setIsDetailsExpanded}>
        <CollapsibleTrigger asChild>
          <Button variant="ghost" className="w-full justify-between p-0 h-auto">
            <span className="text-sm font-medium">
              View {allIssues.length} Issue{allIssues.length !== 1 ? 's' : ''}
            </span>
            {isDetailsExpanded ? (
              <ChevronUp className="w-4 h-4" />
            ) : (
              <ChevronDown className="w-4 h-4" />
            )}
          </Button>
        </CollapsibleTrigger>
        <CollapsibleContent className="mt-3">
          <div className="space-y-3">
            {allIssues.map(renderIssueItem)}
          </div>
        </CollapsibleContent>
      </Collapsible>
    );
  };

  const renderActions = () => {
    if (!showActions) return null;

    return (
      <div className="flex items-center justify-between">
        <div className="text-xs text-muted-foreground">
          Last validated: {validationData ? new Date(validationData.validatedAt).toLocaleTimeString() : 'Never'}
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={handleRefresh}
          disabled={isRefreshing || isLoading}
        >
          <RefreshCw className={cn('w-4 h-4 mr-2', isRefreshing && 'animate-spin')} />
          Refresh
        </Button>
      </div>
    );
  };

  // ===== Loading State =====

  if (isLoading) {
    return (
      <TchatCard className={cn('animate-pulse', className)} size={compact ? 'compact' : 'standard'}>
        <TchatCardContent>
          <div className="space-y-3">
            <div className="h-12 bg-muted rounded-lg" />
            <div className="h-4 bg-muted rounded w-3/4" />
            <div className="h-4 bg-muted rounded w-1/2" />
          </div>
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== Error State =====

  if (error) {
    return (
      <TchatCard variant="outlined" className={cn('border-destructive', className)}>
        <TchatCardContent>
          <div className="flex items-center gap-3 text-destructive py-4">
            <AlertCircle className="w-5 h-5" />
            <div>
              <h3 className="font-medium">Validation Error</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Unable to validate cart. Please try again.
              </p>
            </div>
          </div>
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== No Validation Data =====

  if (!validationData) {
    return (
      <TchatCard variant="outlined" className={className}>
        <TchatCardContent>
          <div className="flex items-center gap-3 py-4">
            <ShieldCheck className="w-5 h-5 text-muted-foreground" />
            <div>
              <h3 className="font-medium">No Validation Data</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Cart validation will appear here when items are added.
              </p>
            </div>
          </div>
        </TchatCardContent>
      </TchatCard>
    );
  }

  // ===== Main Content =====

  return (
    <TchatCard className={className} size={compact ? 'compact' : 'standard'}>
      <TchatCardHeader
        title={compact ? undefined : 'Cart Validation'}
        subtitle={compact ? undefined : 'Real-time cart status and issues'}
      />

      <TchatCardContent>
        <div className="space-y-4">
          {renderValidationStatus()}
          {renderValidationSummary()}
          {renderIssuesList()}
        </div>
      </TchatCardContent>

      {(showActions || validationData.hasIssues) && (
        <TchatCardFooter>
          {renderActions()}
        </TchatCardFooter>
      )}
    </TchatCard>
  );
};

// ===== Export Component =====

export default CartValidation;