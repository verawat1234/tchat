// T052 - MessageComponentFactory with dynamic loading
/**
 * MessageComponentFactory
 * Dynamic factory for loading and rendering message components based on message type
 * Supports lazy loading, error boundaries, and performance optimization
 */

import React, { Suspense, memo, lazy, ComponentType } from 'react';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { Card, CardContent } from '../ui/card';
import { Badge } from '../ui/badge';
import { AlertCircle, Loader2, MessageSquare } from 'lucide-react';
import { cn } from '../../lib/utils';

// Lazy load all message components for code splitting
const ReplyMessage = lazy(() => import('./ReplyMessage'));
const QuizMessage = lazy(() => import('./QuizMessage'));
const EventMessage = lazy(() => import('./EventMessage'));
const SurveyMessage = lazy(() => import('./SurveyMessage'));
const ProductMessage = lazy(() => import('./ProductMessage'));
const FormMessage = lazy(() => import('./FormMessage'));
const RichCardMessage = lazy(() => import('./RichCardMessage'));
const StatusUpdateMessage = lazy(() => import('./StatusUpdateMessage'));
const VenueMessage = lazy(() => import('./VenueMessage'));
const SpreadsheetMessage = lazy(() => import('./SpreadsheetMessage'));
const EmbedMessage = lazy(() => import('./EmbedMessage'));

// Base props interface that all message components share
export interface BaseMessageProps {
  message: MessageData;
  onInteraction?: (interaction: InteractionRequest) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  performanceMode?: boolean;
}

// Extended props for specific component types
export interface MessageComponentProps extends BaseMessageProps {
  [key: string]: any;
}

// Component map for dynamic loading
const MESSAGE_COMPONENT_MAP: Record<MessageType, ComponentType<any>> = {
  [MessageType.TEXT]: () => null, // Text messages handled elsewhere
  [MessageType.REPLY]: ReplyMessage,
  [MessageType.GIF]: () => null, // GIF messages handled elsewhere
  [MessageType.QUIZ]: QuizMessage,
  [MessageType.EVENT]: EventMessage,
  [MessageType.DOCUMENT]: () => null, // Document messages handled elsewhere
  [MessageType.SURVEY]: SurveyMessage,
  [MessageType.PRODUCT]: ProductMessage,
  [MessageType.FORM]: FormMessage,
  [MessageType.RICH_CARD]: RichCardMessage,
  [MessageType.STATUS_UPDATE]: StatusUpdateMessage,
  [MessageType.VENUE]: VenueMessage,
  [MessageType.SPREADSHEET]: SpreadsheetMessage,
  [MessageType.EMBED]: EmbedMessage,
};

// Error boundary for message component failures
interface MessageErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

class MessageErrorBoundary extends React.Component<
  { children: React.ReactNode; messageType: MessageType; messageId: string },
  MessageErrorBoundaryState
> {
  constructor(props: { children: React.ReactNode; messageType: MessageType; messageId: string }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): MessageErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Message component error:', {
      messageType: this.props.messageType,
      messageId: this.props.messageId,
      error: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
    });
  }

  render() {
    if (this.state.hasError) {
      return (
        <Card className="border-destructive/50 bg-destructive/5">
          <CardContent className="p-4">
            <div className="flex items-start gap-3">
              <AlertCircle className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <h4 className="font-medium text-destructive">Message Error</h4>
                  <Badge variant="outline" className="text-xs">
                    {this.props.messageType}
                  </Badge>
                </div>
                <p className="text-sm text-muted-foreground">
                  Unable to display this {this.props.messageType.toLowerCase()} message.
                </p>
                {process.env.NODE_ENV === 'development' && this.state.error && (
                  <details className="mt-2">
                    <summary className="text-xs text-muted-foreground cursor-pointer">
                      Technical details
                    </summary>
                    <pre className="text-xs text-muted-foreground mt-1 p-2 bg-muted rounded overflow-auto">
                      {this.state.error.message}
                    </pre>
                  </details>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
      );
    }

    return this.props.children;
  }
}

// Loading fallback component
const MessageLoadingFallback: React.FC<{
  messageType: MessageType;
  compactMode?: boolean;
}> = memo(({ messageType, compactMode = false }) => (
  <Card className="animate-pulse">
    <CardContent className={cn("p-4", compactMode && "p-3")}>
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 bg-muted rounded-full flex-shrink-0" />
        <div className="space-y-2 flex-1">
          <div className="flex items-center gap-2">
            <div className="h-4 bg-muted rounded w-24" />
            <Badge variant="outline" className="text-xs">
              <Loader2 className="w-3 h-3 mr-1 animate-spin" />
              Loading {messageType.toLowerCase()}...
            </Badge>
          </div>
          <div className="space-y-1">
            <div className="h-3 bg-muted rounded w-full" />
            <div className="h-3 bg-muted rounded w-3/4" />
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
));

MessageLoadingFallback.displayName = 'MessageLoadingFallback';

// Unsupported message type fallback
const UnsupportedMessage: React.FC<{
  message: MessageData;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
}> = memo(({ message, showAvatar = true, showTimestamp = true, compactMode = false }) => (
  <Card className="border-orange-200 bg-orange-50 dark:border-orange-800 dark:bg-orange-950/50">
    <CardContent className={cn("p-4", compactMode && "p-3")}>
      <div className="flex items-start gap-3">
        {showAvatar && (
          <div className="w-10 h-10 bg-orange-100 dark:bg-orange-900/50 rounded-full flex items-center justify-center flex-shrink-0">
            <MessageSquare className="w-5 h-5 text-orange-600 dark:text-orange-400" />
          </div>
        )}
        <div className="space-y-2 flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="font-medium text-foreground">
              {message.senderName}
            </span>
            <Badge variant="outline" className="text-xs border-orange-500 text-orange-600">
              {message.type}
            </Badge>
            {showTimestamp && (
              <span className="text-xs text-muted-foreground">
                {message.timestamp.toLocaleDateString()}
              </span>
            )}
          </div>
          <div className="text-sm text-muted-foreground">
            <p>This message type is not yet supported in this version.</p>
            <p className="text-xs mt-1">
              Message ID: <code className="bg-muted px-1 rounded">{message.id}</code>
            </p>
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
));

UnsupportedMessage.displayName = 'UnsupportedMessage';

// Main factory component
export interface MessageComponentFactoryProps extends BaseMessageProps {
  // Component-specific props can be passed through
  [key: string]: any;
}

export const MessageComponentFactory: React.FC<MessageComponentFactoryProps> = memo((props) => {
  const { message, ...componentProps } = props;
  const { type: messageType, id: messageId } = message;

  // Get the component for this message type
  const Component = MESSAGE_COMPONENT_MAP[messageType];

  // Handle unsupported message types
  if (!Component) {
    return (
      <UnsupportedMessage
        message={message}
        showAvatar={props.showAvatar}
        showTimestamp={props.showTimestamp}
        compactMode={props.compactMode}
      />
    );
  }

  // Handle null components (simple message types handled elsewhere)
  if (Component === null || (typeof Component === 'function' && Component() === null)) {
    return null;
  }

  // Render with error boundary and suspense
  return (
    <MessageErrorBoundary messageType={messageType} messageId={messageId}>
      <Suspense
        fallback={
          <MessageLoadingFallback
            messageType={messageType}
            compactMode={props.compactMode}
          />
        }
      >
        <Component {...componentProps} message={message} />
      </Suspense>
    </MessageErrorBoundary>
  );
});

MessageComponentFactory.displayName = 'MessageComponentFactory';

// Utility function to check if a message type is supported
export const isMessageTypeSupported = (messageType: MessageType): boolean => {
  const Component = MESSAGE_COMPONENT_MAP[messageType];
  return Component !== undefined && Component !== null &&
         !(typeof Component === 'function' && Component() === null);
};

// Utility function to get component loading priority
export const getComponentLoadingPriority = (messageType: MessageType): 'high' | 'normal' | 'low' => {
  // Prioritize commonly used interactive components
  switch (messageType) {
    case MessageType.REPLY:
    case MessageType.QUIZ:
    case MessageType.SURVEY:
      return 'high';

    case MessageType.EVENT:
    case MessageType.FORM:
    case MessageType.PRODUCT:
      return 'normal';

    default:
      return 'low';
  }
};

// Utility function to preload components based on priority
export const preloadMessageComponents = async (
  messageTypes: MessageType[],
  priority: 'high' | 'normal' | 'low' = 'normal'
): Promise<void> => {
  const componentsToLoad = messageTypes
    .filter(type => isMessageTypeSupported(type))
    .filter(type => getComponentLoadingPriority(type) === priority || priority === 'normal');

  const loadPromises = componentsToLoad.map(async (messageType) => {
    try {
      const Component = MESSAGE_COMPONENT_MAP[messageType];
      if (Component && typeof Component === 'object' && 'then' in Component) {
        await Component;
      }
    } catch (error) {
      console.warn(`Failed to preload component for ${messageType}:`, error);
    }
  });

  await Promise.allSettled(loadPromises);
};

// Performance monitoring hook
export const useMessageComponentPerformance = (messageType: MessageType) => {
  const [loadTime, setLoadTime] = React.useState<number | null>(null);
  const [renderTime, setRenderTime] = React.useState<number | null>(null);

  React.useEffect(() => {
    const startTime = performance.now();

    // Measure component load time
    const Component = MESSAGE_COMPONENT_MAP[messageType];
    if (Component && typeof Component === 'object' && 'then' in Component) {
      Component.then(() => {
        setLoadTime(performance.now() - startTime);
      }).catch(() => {
        setLoadTime(-1); // Error indicator
      });
    }

    // Measure render time
    const renderStartTime = performance.now();
    const cleanup = () => {
      setRenderTime(performance.now() - renderStartTime);
    };

    return cleanup;
  }, [messageType]);

  return { loadTime, renderTime };
};

// Export component map for testing and debugging
export { MESSAGE_COMPONENT_MAP };

export default MessageComponentFactory;