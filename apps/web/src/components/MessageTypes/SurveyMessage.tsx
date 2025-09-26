// T044 - SurveyMessage component with interactive polling
/**
 * SurveyMessage Component
 * Displays interactive surveys and polls with real-time result visualization
 * Supports multiple question types, anonymous voting, and result analytics
 */

import React, { useState, useCallback, useMemo, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { SurveyContent, SurveyQuestion, SurveyQuestionType, SurveyStatus } from '../../types/SurveyContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Progress } from '../ui/progress';
import { Separator } from '../ui/separator';
import { RadioGroup, RadioGroupItem } from '../ui/radio-group';
import { Checkbox } from '../ui/checkbox';
import { Textarea } from '../ui/textarea';
import { Label } from '../ui/label';
import {
  BarChart3,
  Clock,
  Users,
  CheckCircle,
  Circle,
  Square,
  Type,
  Hash,
  Eye,
  EyeOff,
  Send,
  RotateCcw,
  TrendingUp,
  AlertCircle
} from 'lucide-react';

// Component Props Interface
interface SurveyMessageProps {
  message: MessageData & { content: SurveyContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onSurveySubmit?: (surveyId: string, answers: Record<string, any>) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  showResults?: boolean;
  showLiveResults?: boolean;
  performanceMode?: boolean;
}

// Survey Response Interface
interface SurveyResponse {
  questionId: string;
  answer: string | string[] | number | boolean;
  confidence?: number;
}

// Animation Variants
const surveyVariants = {
  initial: { opacity: 0, scale: 0.98, y: 15 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.98, y: -15 }
};

const questionVariants = {
  initial: { opacity: 0, x: -10 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 10 }
};

const resultsVariants = {
  initial: { opacity: 0, height: 0 },
  animate: { opacity: 1, height: 'auto' },
  exit: { opacity: 0, height: 0 }
};

const progressVariants = {
  initial: { width: 0 },
  animate: { width: '100%' }
};

export const SurveyMessage: React.FC<SurveyMessageProps> = ({
  message,
  onInteraction,
  onSurveySubmit,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  showResults = false,
  showLiveResults = true,
  performanceMode = false
}) => {
  const surveyRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Survey state
  const [responses, setResponses] = useState<Record<string, SurveyResponse>>({});
  const [currentPage, setCurrentPage] = useState(0);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [showResultsView, setShowResultsView] = useState(showResults);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Calculate survey progress
  const surveyProgress = useMemo(() => {
    if (!content.questions.length) return { completed: 0, total: 0, percentage: 0 };

    const requiredQuestions = content.questions.filter(q => q.required);
    const completedRequired = requiredQuestions.filter(q => responses[q.id]).length;
    const totalResponses = Object.keys(responses).length;

    return {
      completed: totalResponses,
      total: content.questions.length,
      requiredCompleted: completedRequired,
      requiredTotal: requiredQuestions.length,
      percentage: content.questions.length > 0 ? (totalResponses / content.questions.length) * 100 : 0,
      canSubmit: completedRequired === requiredQuestions.length
    };
  }, [responses, content.questions]);

  // Get survey timing info
  const surveyTiming = useMemo(() => {
    const now = new Date();
    const endTime = content.endTime ? new Date(content.endTime) : null;
    const isExpired = endTime ? endTime <= now : false;
    const timeRemaining = endTime ? Math.max(0, endTime.getTime() - now.getTime()) : null;

    return {
      isExpired,
      timeRemaining,
      hasDeadline: !!endTime
    };
  }, [content.endTime]);

  // Format time remaining
  const formatTimeRemaining = useCallback((milliseconds: number) => {
    const days = Math.floor(milliseconds / (1000 * 60 * 60 * 24));
    const hours = Math.floor((milliseconds % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const minutes = Math.floor((milliseconds % (1000 * 60 * 60)) / (1000 * 60));

    if (days > 0) return `${days}d ${hours}h`;
    if (hours > 0) return `${hours}h ${minutes}m`;
    return `${minutes}m`;
  }, []);

  // Handle question response
  const handleQuestionResponse = useCallback((questionId: string, answer: any) => {
    if (readonly || isSubmitted || surveyTiming.isExpired) return;

    setResponses(prev => ({
      ...prev,
      [questionId]: {
        questionId,
        answer,
        confidence: 1
      }
    }));

    // Send interaction for live updates
    if (onInteraction && showLiveResults) {
      onInteraction({
        messageId: message.id,
        interactionType: 'survey_response',
        data: { questionId, answer, surveyId: content.id },
        userId: 'current-user',
        timestamp: new Date()
      });
    }
  }, [readonly, isSubmitted, surveyTiming.isExpired, onInteraction, showLiveResults, message.id, content.id]);

  // Handle survey submission
  const handleSubmit = useCallback(async () => {
    if (!surveyProgress.canSubmit || isSubmitting) return;

    setIsSubmitting(true);

    try {
      if (onSurveySubmit) {
        const answerMap = Object.fromEntries(
          Object.values(responses).map(response => [response.questionId, response.answer])
        );
        await onSurveySubmit(content.id, answerMap);
      }

      if (onInteraction) {
        onInteraction({
          messageId: message.id,
          interactionType: 'survey_submit',
          data: { surveyId: content.id, responses: Object.values(responses) },
          userId: 'current-user',
          timestamp: new Date()
        });
      }

      setIsSubmitted(true);
      setShowResultsView(true);
    } catch (error) {
      console.error('Survey submission failed:', error);
    } finally {
      setIsSubmitting(false);
    }
  }, [surveyProgress.canSubmit, isSubmitting, responses, content.id, onSurveySubmit, onInteraction, message.id]);

  // Reset survey
  const handleReset = useCallback(() => {
    if (readonly) return;

    setResponses({});
    setCurrentPage(0);
    setIsSubmitted(false);
    setShowResultsView(false);
  }, [readonly]);

  // Get question icon
  const getQuestionIcon = useCallback((type: SurveyQuestionType) => {
    switch (type) {
      case SurveyQuestionType.MULTIPLE_CHOICE: return <Circle className="w-4 h-4" />;
      case SurveyQuestionType.MULTIPLE_SELECT: return <Square className="w-4 h-4" />;
      case SurveyQuestionType.TEXT: return <Type className="w-4 h-4" />;
      case SurveyQuestionType.TEXTAREA: return <Type className="w-4 h-4" />;
      case SurveyQuestionType.NUMBER: return <Hash className="w-4 h-4" />;
      case SurveyQuestionType.RATING: return <TrendingUp className="w-4 h-4" />;
      default: return <Circle className="w-4 h-4" />;
    }
  }, []);

  // Render question based on type
  const renderQuestion = useCallback((question: SurveyQuestion, index: number) => {
    const response = responses[question.id];
    const isAnswered = !!response;

    const QuestionWrapper = ({ children }: { children: React.ReactNode }) => (
      <motion.div
        key={question.id}
        variants={performanceMode ? {} : questionVariants}
        initial={performanceMode ? {} : "initial"}
        animate={performanceMode ? {} : "animate"}
        exit={performanceMode ? {} : "exit"}
        transition={{ delay: index * 0.05 }}
        className={cn(
          "question-container space-y-3 p-4 rounded-lg border",
          isAnswered && "bg-green-50/50 border-green-200 dark:bg-green-950/20 dark:border-green-800",
          !isAnswered && "bg-card"
        )}
      >
        <div className="flex items-start gap-2">
          {getQuestionIcon(question.type)}
          <div className="flex-1 space-y-3">
            <div className="flex items-start justify-between gap-2">
              <div>
                <h4 className="font-medium text-sm text-foreground leading-tight">
                  {question.question}
                  {question.required && (
                    <span className="text-destructive ml-1" aria-label="Required">*</span>
                  )}
                </h4>
                {question.description && (
                  <p className="text-xs text-muted-foreground mt-1">
                    {question.description}
                  </p>
                )}
              </div>
              {isAnswered && <CheckCircle className="w-4 h-4 text-green-600 flex-shrink-0" />}
            </div>
            {children}
          </div>
        </div>
      </motion.div>
    );

    switch (question.type) {
      case SurveyQuestionType.MULTIPLE_CHOICE:
        return (
          <QuestionWrapper key={question.id}>
            <RadioGroup
              value={response?.answer as string || ''}
              onValueChange={(value) => handleQuestionResponse(question.id, value)}
              disabled={readonly || isSubmitted}
            >
              {question.options?.map((option, optionIndex) => (
                <div key={optionIndex} className="flex items-center space-x-2">
                  <RadioGroupItem value={option} id={`${question.id}-${optionIndex}`} />
                  <Label
                    htmlFor={`${question.id}-${optionIndex}`}
                    className="text-sm font-normal cursor-pointer flex-1"
                  >
                    {option}
                  </Label>
                </div>
              ))}
            </RadioGroup>
          </QuestionWrapper>
        );

      case SurveyQuestionType.MULTIPLE_SELECT:
        return (
          <QuestionWrapper key={question.id}>
            <div className="space-y-2">
              {question.options?.map((option, optionIndex) => {
                const currentAnswers = (response?.answer as string[]) || [];
                const isChecked = currentAnswers.includes(option);

                return (
                  <div key={optionIndex} className="flex items-center space-x-2">
                    <Checkbox
                      id={`${question.id}-${optionIndex}`}
                      checked={isChecked}
                      onCheckedChange={(checked) => {
                        const newAnswers = checked
                          ? [...currentAnswers, option]
                          : currentAnswers.filter(a => a !== option);
                        handleQuestionResponse(question.id, newAnswers);
                      }}
                      disabled={readonly || isSubmitted}
                    />
                    <Label
                      htmlFor={`${question.id}-${optionIndex}`}
                      className="text-sm font-normal cursor-pointer flex-1"
                    >
                      {option}
                    </Label>
                  </div>
                );
              })}
            </div>
          </QuestionWrapper>
        );

      case SurveyQuestionType.TEXT:
        return (
          <QuestionWrapper key={question.id}>
            <input
              type="text"
              value={(response?.answer as string) || ''}
              onChange={(e) => handleQuestionResponse(question.id, e.target.value)}
              disabled={readonly || isSubmitted}
              placeholder="Type your answer..."
              className="w-full px-3 py-2 text-sm border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-primary/20 disabled:opacity-50"
            />
          </QuestionWrapper>
        );

      case SurveyQuestionType.TEXTAREA:
        return (
          <QuestionWrapper key={question.id}>
            <Textarea
              value={(response?.answer as string) || ''}
              onChange={(e) => handleQuestionResponse(question.id, e.target.value)}
              disabled={readonly || isSubmitted}
              placeholder="Type your detailed answer..."
              rows={3}
              className="resize-none"
            />
          </QuestionWrapper>
        );

      case SurveyQuestionType.NUMBER:
        return (
          <QuestionWrapper key={question.id}>
            <input
              type="number"
              value={(response?.answer as number) || ''}
              onChange={(e) => handleQuestionResponse(question.id, parseInt(e.target.value) || 0)}
              disabled={readonly || isSubmitted}
              placeholder="Enter a number..."
              min={question.validation?.min}
              max={question.validation?.max}
              className="w-full px-3 py-2 text-sm border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-primary/20 disabled:opacity-50"
            />
          </QuestionWrapper>
        );

      case SurveyQuestionType.RATING:
        const maxRating = question.validation?.max || 5;
        const currentRating = (response?.answer as number) || 0;

        return (
          <QuestionWrapper key={question.id}>
            <div className="flex items-center gap-2">
              {Array.from({ length: maxRating }, (_, i) => i + 1).map((rating) => (
                <Button
                  key={rating}
                  variant={currentRating >= rating ? "default" : "outline"}
                  size="sm"
                  onClick={() => handleQuestionResponse(question.id, rating)}
                  disabled={readonly || isSubmitted}
                  className="w-8 h-8 p-0"
                >
                  {rating}
                </Button>
              ))}
              <span className="text-sm text-muted-foreground ml-2">
                {currentRating > 0 ? `${currentRating}/${maxRating}` : `Rate 1-${maxRating}`}
              </span>
            </div>
          </QuestionWrapper>
        );

      default:
        return null;
    }
  }, [responses, readonly, isSubmitted, handleQuestionResponse, performanceMode, getQuestionIcon]);

  // Render results visualization
  const renderResults = useCallback(() => {
    if (!showResultsView || !content.analytics) return null;

    return (
      <motion.div
        variants={performanceMode ? {} : resultsVariants}
        initial={performanceMode ? {} : "initial"}
        animate={performanceMode ? {} : "animate"}
        exit={performanceMode ? {} : "exit"}
        className="space-y-4"
      >
        <Separator />
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <h4 className="font-medium text-sm text-foreground flex items-center gap-2">
              <BarChart3 className="w-4 h-4" />
              Survey Results
            </h4>
            <Badge variant="secondary" className="text-xs">
              {content.analytics.totalResponses} responses
            </Badge>
          </div>

          {content.questions.map((question, index) => {
            const questionStats = content.analytics?.questionStats?.[question.id];
            if (!questionStats) return null;

            return (
              <div key={question.id} className="space-y-2">
                <p className="text-sm font-medium">{question.question}</p>

                {question.type === SurveyQuestionType.MULTIPLE_CHOICE && questionStats.optionCounts && (
                  <div className="space-y-1">
                    {Object.entries(questionStats.optionCounts).map(([option, count]) => {
                      const percentage = content.analytics?.totalResponses
                        ? (count / content.analytics.totalResponses) * 100
                        : 0;

                      return (
                        <div key={option} className="space-y-1">
                          <div className="flex justify-between text-xs">
                            <span className="text-muted-foreground truncate">{option}</span>
                            <span className="font-medium">{count} ({percentage.toFixed(0)}%)</span>
                          </div>
                          <motion.div
                            variants={performanceMode ? {} : progressVariants}
                            initial={performanceMode ? {} : "initial"}
                            animate={performanceMode ? {} : "animate"}
                            transition={{ delay: index * 0.1 }}
                          >
                            <Progress value={percentage} className="h-2" />
                          </motion.div>
                        </div>
                      );
                    })}
                  </div>
                )}

                {question.type === SurveyQuestionType.RATING && (
                  <div className="text-sm text-muted-foreground">
                    Average: {questionStats.averageRating?.toFixed(1)}/5
                    ({questionStats.responseCount} responses)
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </motion.div>
    );
  }, [showResultsView, content.analytics, content.questions, performanceMode]);

  // Performance optimization: skip animation in performance mode
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: surveyVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={surveyRef}
        className={cn(
          "survey-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`survey-message-${message.id}`}
        data-survey-status={content.status}
        role="article"
        aria-label={`Survey: ${content.title}`}
      >
        <Card className="survey-card">
          <CardHeader className="space-y-3">
            {/* Header with sender info and survey status */}
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
                      Survey • {message.timestamp.toLocaleDateString()}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-2">
                <Badge variant={content.isAnonymous ? "secondary" : "outline"} className="text-xs">
                  {content.isAnonymous ? (
                    <>
                      <EyeOff className="w-3 h-3 mr-1" />
                      Anonymous
                    </>
                  ) : (
                    <>
                      <Eye className="w-3 h-3 mr-1" />
                      Public
                    </>
                  )}
                </Badge>

                {surveyTiming.hasDeadline && surveyTiming.timeRemaining && surveyTiming.timeRemaining > 0 && (
                  <Badge variant="outline" className="text-xs">
                    <Clock className="w-3 h-3 mr-1" />
                    {formatTimeRemaining(surveyTiming.timeRemaining)}
                  </Badge>
                )}

                {surveyTiming.isExpired && (
                  <Badge variant="destructive" className="text-xs">
                    <AlertCircle className="w-3 h-3 mr-1" />
                    Expired
                  </Badge>
                )}
              </div>
            </div>

            {/* Survey title and description */}
            <div className="space-y-2">
              <CardTitle className={cn(
                "text-lg font-bold text-foreground leading-tight",
                compactMode && "text-base"
              )}>
                {content.title}
              </CardTitle>
              {content.description && (
                <p className="text-muted-foreground text-sm leading-relaxed">
                  {content.description}
                </p>
              )}
            </div>

            {/* Progress bar */}
            {!readonly && !surveyTiming.isExpired && !isSubmitted && (
              <div className="space-y-2">
                <div className="flex justify-between text-xs">
                  <span className="text-muted-foreground">
                    Progress: {surveyProgress.completed}/{surveyProgress.total} questions
                  </span>
                  <span className="font-medium">
                    {Math.round(surveyProgress.percentage)}%
                  </span>
                </div>
                <Progress value={surveyProgress.percentage} className="h-2" />
              </div>
            )}
          </CardHeader>

          <CardContent className="space-y-4">
            {/* Questions */}
            {!readonly && !surveyTiming.isExpired && !isSubmitted ? (
              <AnimatePresence mode="wait">
                <div className="space-y-4">
                  {content.questions.map((question, index) => renderQuestion(question, index))}
                </div>
              </AnimatePresence>
            ) : (
              <div className="text-center py-8">
                <div className="text-muted-foreground">
                  {isSubmitted && (
                    <>
                      <CheckCircle className="w-12 h-12 mx-auto mb-3 text-green-600" />
                      <p className="font-medium">Survey submitted successfully!</p>
                      <p className="text-sm">Thank you for your responses.</p>
                    </>
                  )}
                  {surveyTiming.isExpired && !isSubmitted && (
                    <>
                      <AlertCircle className="w-12 h-12 mx-auto mb-3 text-orange-600" />
                      <p className="font-medium">Survey expired</p>
                      <p className="text-sm">This survey is no longer accepting responses.</p>
                    </>
                  )}
                  {readonly && (
                    <>
                      <Eye className="w-12 h-12 mx-auto mb-3 text-muted-foreground" />
                      <p className="font-medium">Survey preview</p>
                      <p className="text-sm">This is a read-only view of the survey.</p>
                    </>
                  )}
                </div>
              </div>
            )}

            {/* Action buttons */}
            {!readonly && !surveyTiming.isExpired && !isSubmitted && (
              <motion.div
                initial={performanceMode ? {} : { opacity: 0, y: 10 }}
                animate={performanceMode ? {} : { opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
                className="flex items-center justify-between pt-2"
              >
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <Users className="w-4 h-4" />
                  <span>{content.analytics?.totalResponses || 0} responses</span>
                </div>

                <div className="flex items-center gap-2">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        onClick={handleReset}
                        variant="ghost"
                        size="sm"
                        disabled={Object.keys(responses).length === 0}
                        className="gap-2"
                      >
                        <RotateCcw className="w-4 h-4" />
                        Reset
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Clear all responses</TooltipContent>
                  </Tooltip>

                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        onClick={handleSubmit}
                        disabled={!surveyProgress.canSubmit || isSubmitting}
                        size="sm"
                        className="gap-2"
                      >
                        <Send className="w-4 h-4" />
                        {isSubmitting ? 'Submitting...' : 'Submit Survey'}
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      {surveyProgress.canSubmit
                        ? 'Submit your survey responses'
                        : `Please answer ${surveyProgress.requiredTotal - surveyProgress.requiredCompleted} more required questions`
                      }
                    </TooltipContent>
                  </Tooltip>
                </div>
              </motion.div>
            )}

            {/* Results section */}
            <AnimatePresence mode="wait">
              {renderResults()}
            </AnimatePresence>

            {/* Results toggle */}
            {(isSubmitted || showResults) && content.analytics && (
              <div className="flex justify-center pt-2">
                <Button
                  onClick={() => setShowResultsView(!showResultsView)}
                  variant="ghost"
                  size="sm"
                  className="gap-2"
                >
                  <BarChart3 className="w-4 h-4" />
                  {showResultsView ? 'Hide Results' : 'View Results'}
                </Button>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Performance Debug Info (Development Only) */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            Q: {surveyProgress.completed}/{surveyProgress.total} | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedSurveyMessage = React.memo(SurveyMessage, (prevProps, nextProps) => {
  // Custom comparison function for performance
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.showResults === nextProps.showResults &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedSurveyMessage.displayName = 'MemoizedSurveyMessage';

// Export both versions
export default SurveyMessage;