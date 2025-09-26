// T042 - QuizMessage component with interactive questions
/**
 * QuizMessage Component
 * Interactive quiz component with real-time validation, scoring, and analytics
 * Supports multiple question types, hints, time limits, and accessibility
 */

import React, { useState, useEffect, useCallback, useMemo, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, InteractionRequest } from '../../types/MessageData';
import {
  QuizContent,
  QuizState,
  QuizStatus,
  QuizQuestion,
  QuizQuestionType,
  QuizAnswerValue,
  QuizValidationError
} from '../../types/QuizContent';
import { Button } from '../ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Progress } from '../ui/progress';
import { Badge } from '../ui/badge';
import { Separator } from '../ui/separator';
import { Alert, AlertDescription } from '../ui/alert';
import { Checkbox } from '../ui/checkbox';
import { RadioGroup, RadioGroupItem } from '../ui/radio-group';
import { Label } from '../ui/label';
import { Input } from '../ui/input';
import { Textarea } from '../ui/textarea';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import {
  Play,
  Pause,
  RotateCcw,
  CheckCircle,
  XCircle,
  Clock,
  Users,
  Award,
  HelpCircle,
  Flag,
  ChevronLeft,
  ChevronRight,
  Send,
  AlertTriangle
} from 'lucide-react';

// Component Props Interface
interface QuizMessageProps {
  message: MessageData & { content: QuizContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  className?: string;
  readonly?: boolean;
  showResults?: boolean;
  allowRetake?: boolean;
  currentUserId?: string;
}

// Question Component Props
interface QuestionProps {
  question: QuizQuestion;
  answer: QuizAnswerValue;
  onAnswerChange: (questionId: string, answer: QuizAnswerValue) => void;
  readonly?: boolean;
  showResult?: boolean;
  isCorrect?: boolean;
  explanation?: string;
  onHintRequest?: (questionId: string, hintId: string) => void;
}

// Animation Variants
const quizVariants = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -20 }
};

const questionVariants = {
  initial: { opacity: 0, x: 20 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: -20 }
};

// Question Component
const QuestionComponent: React.FC<QuestionProps> = ({
  question,
  answer,
  onAnswerChange,
  readonly = false,
  showResult = false,
  isCorrect = false,
  explanation,
  onHintRequest
}) => {
  const [showHints, setShowHints] = useState(false);

  const handleAnswerChange = useCallback((newAnswer: QuizAnswerValue) => {
    if (!readonly) {
      onAnswerChange(question.id, newAnswer);
    }
  }, [question.id, onAnswerChange, readonly]);

  const renderQuestionContent = () => {
    switch (question.type) {
      case QuizQuestionType.MULTIPLE_CHOICE:
        return (
          <RadioGroup
            value={answer as string || ''}
            onValueChange={handleAnswerChange}
            disabled={readonly}
            className="space-y-2"
          >
            {question.options?.map((option, index) => (
              <div
                key={index}
                className={cn(
                  "flex items-center space-x-2 p-2 rounded-md transition-colors",
                  showResult && question.correctAnswer === option && "bg-green-50 border border-green-200",
                  showResult && answer === option && question.correctAnswer !== option && "bg-red-50 border border-red-200"
                )}
              >
                <RadioGroupItem value={option} id={`${question.id}-${index}`} />
                <Label
                  htmlFor={`${question.id}-${index}`}
                  className={cn(
                    "flex-1 cursor-pointer",
                    showResult && question.correctAnswer === option && "font-semibold text-green-700",
                    showResult && answer === option && question.correctAnswer !== option && "text-red-700"
                  )}
                >
                  {option}
                  {showResult && question.correctAnswer === option && (
                    <CheckCircle className="inline-block w-4 h-4 ml-2 text-green-600" />
                  )}
                  {showResult && answer === option && question.correctAnswer !== option && (
                    <XCircle className="inline-block w-4 h-4 ml-2 text-red-600" />
                  )}
                </Label>
              </div>
            ))}
          </RadioGroup>
        );

      case QuizQuestionType.MULTIPLE_SELECT:
        const selectedAnswers = Array.isArray(answer) ? answer : [];
        const correctAnswers = Array.isArray(question.correctAnswer) ? question.correctAnswer : [];

        return (
          <div className="space-y-2">
            {question.options?.map((option, index) => {
              const isSelected = selectedAnswers.includes(option);
              const isCorrectOption = correctAnswers.includes(option);
              const shouldBeSelected = showResult && isCorrectOption;
              const wrongSelection = showResult && isSelected && !isCorrectOption;

              return (
                <div
                  key={index}
                  className={cn(
                    "flex items-center space-x-2 p-2 rounded-md transition-colors",
                    shouldBeSelected && "bg-green-50 border border-green-200",
                    wrongSelection && "bg-red-50 border border-red-200"
                  )}
                >
                  <Checkbox
                    id={`${question.id}-${index}`}
                    checked={isSelected}
                    disabled={readonly}
                    onCheckedChange={(checked) => {
                      if (checked) {
                        handleAnswerChange([...selectedAnswers, option]);
                      } else {
                        handleAnswerChange(selectedAnswers.filter(a => a !== option));
                      }
                    }}
                  />
                  <Label
                    htmlFor={`${question.id}-${index}`}
                    className={cn(
                      "flex-1 cursor-pointer",
                      shouldBeSelected && "font-semibold text-green-700",
                      wrongSelection && "text-red-700"
                    )}
                  >
                    {option}
                    {shouldBeSelected && (
                      <CheckCircle className="inline-block w-4 h-4 ml-2 text-green-600" />
                    )}
                    {wrongSelection && (
                      <XCircle className="inline-block w-4 h-4 ml-2 text-red-600" />
                    )}
                  </Label>
                </div>
              );
            })}
          </div>
        );

      case QuizQuestionType.TRUE_FALSE:
        return (
          <RadioGroup
            value={answer?.toString() || ''}
            onValueChange={(value) => handleAnswerChange(value === 'true')}
            disabled={readonly}
            className="flex space-x-6"
          >
            <div className={cn(
              "flex items-center space-x-2 p-3 rounded-md border-2 transition-all",
              showResult && question.correctAnswer === true && "border-green-300 bg-green-50",
              showResult && answer === true && question.correctAnswer !== true && "border-red-300 bg-red-50"
            )}>
              <RadioGroupItem value="true" id={`${question.id}-true`} />
              <Label htmlFor={`${question.id}-true`} className="cursor-pointer font-medium">
                True
              </Label>
            </div>
            <div className={cn(
              "flex items-center space-x-2 p-3 rounded-md border-2 transition-all",
              showResult && question.correctAnswer === false && "border-green-300 bg-green-50",
              showResult && answer === false && question.correctAnswer !== false && "border-red-300 bg-red-50"
            )}>
              <RadioGroupItem value="false" id={`${question.id}-false`} />
              <Label htmlFor={`${question.id}-false`} className="cursor-pointer font-medium">
                False
              </Label>
            </div>
          </RadioGroup>
        );

      case QuizQuestionType.SHORT_ANSWER:
        return (
          <Input
            value={answer as string || ''}
            onChange={(e) => handleAnswerChange(e.target.value)}
            placeholder="Enter your answer..."
            disabled={readonly}
            className={cn(
              showResult && isCorrect && "border-green-300 bg-green-50",
              showResult && !isCorrect && answer && "border-red-300 bg-red-50"
            )}
          />
        );

      case QuizQuestionType.LONG_ANSWER:
        return (
          <Textarea
            value={answer as string || ''}
            onChange={(e) => handleAnswerChange(e.target.value)}
            placeholder="Enter your detailed answer..."
            rows={4}
            disabled={readonly}
            className={cn(
              showResult && isCorrect && "border-green-300 bg-green-50",
              showResult && !isCorrect && answer && "border-red-300 bg-red-50"
            )}
          />
        );

      case QuizQuestionType.NUMBER:
        return (
          <Input
            type="number"
            value={answer as number || ''}
            onChange={(e) => handleAnswerChange(parseFloat(e.target.value) || 0)}
            placeholder="Enter a number..."
            disabled={readonly}
            className={cn(
              showResult && isCorrect && "border-green-300 bg-green-50",
              showResult && !isCorrect && answer !== null && "border-red-300 bg-red-50"
            )}
          />
        );

      default:
        return <div className="text-muted-foreground">Unsupported question type</div>;
    }
  };

  return (
    <motion.div
      variants={questionVariants}
      initial="initial"
      animate="animate"
      exit="exit"
      className="space-y-4"
    >
      {/* Question Header */}
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <h4 className="text-lg font-medium leading-tight mb-2">
            {question.question}
          </h4>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Badge variant="secondary">{question.points} points</Badge>
            {question.hints && question.hints.length > 0 && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowHints(!showHints)}
                className="h-6 px-2 text-xs"
              >
                <HelpCircle className="w-3 h-3 mr-1" />
                {question.hints.length} hints
              </Button>
            )}
          </div>
        </div>
      </div>

      {/* Question Media */}
      {question.media && (
        <div className="my-4">
          {question.media.type === 'image' && (
            <img
              src={question.media.url}
              alt={question.media.alt || 'Question image'}
              className="max-w-full h-auto rounded-md"
            />
          )}
          {question.media.type === 'video' && (
            <video
              controls
              className="max-w-full h-auto rounded-md"
              poster={question.media.thumbnailUrl}
            >
              <source src={question.media.url} type="video/mp4" />
            </video>
          )}
        </div>
      )}

      {/* Answer Input */}
      <div className="space-y-3">
        {renderQuestionContent()}
      </div>

      {/* Hints */}
      <AnimatePresence>
        {showHints && question.hints && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="space-y-2"
          >
            {question.hints.map((hint, index) => (
              <Alert key={hint.id}>
                <HelpCircle className="h-4 w-4" />
                <AlertDescription>
                  <strong>Hint {index + 1}:</strong> {hint.text}
                  {hint.costInPoints && (
                    <Badge variant="outline" className="ml-2">
                      -{hint.costInPoints} points
                    </Badge>
                  )}
                </AlertDescription>
              </Alert>
            ))}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Result and Explanation */}
      {showResult && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className={cn(
            "p-3 rounded-md border",
            isCorrect ? "bg-green-50 border-green-200" : "bg-red-50 border-red-200"
          )}
        >
          <div className="flex items-center gap-2 mb-2">
            {isCorrect ? (
              <CheckCircle className="w-5 h-5 text-green-600" />
            ) : (
              <XCircle className="w-5 h-5 text-red-600" />
            )}
            <span className={cn(
              "font-medium",
              isCorrect ? "text-green-700" : "text-red-700"
            )}>
              {isCorrect ? "Correct!" : "Incorrect"}
            </span>
          </div>
          {explanation && (
            <p className="text-sm text-muted-foreground">{explanation}</p>
          )}
        </motion.div>
      )}
    </motion.div>
  );
};

// Main QuizMessage Component
export const QuizMessage: React.FC<QuizMessageProps> = ({
  message,
  onInteraction,
  className,
  readonly = false,
  showResults = false,
  allowRetake = false,
  currentUserId = 'current-user'
}) => {
  const { content } = message;
  const timerRef = useRef<NodeJS.Timeout>();

  // Quiz State Management
  const [quizState, setQuizState] = useState<QuizState>({
    quizId: content.id,
    status: QuizStatus.NOT_STARTED,
    currentQuestion: 0,
    totalQuestions: content.questions.length,
    answers: {},
    flaggedQuestions: [],
    reviewMode: false,
    canSubmit: false,
    validationErrors: []
  });

  const [timeRemaining, setTimeRemaining] = useState<number | null>(
    content.timeLimit || null
  );

  // Timer Effect
  useEffect(() => {
    if (
      quizState.status === QuizStatus.IN_PROGRESS &&
      timeRemaining !== null &&
      timeRemaining > 0
    ) {
      timerRef.current = setInterval(() => {
        setTimeRemaining(prev => {
          if (prev !== null && prev <= 1) {
            setQuizState(prev => ({ ...prev, status: QuizStatus.TIME_EXPIRED }));
            return 0;
          }
          return prev !== null ? prev - 1 : null;
        });
      }, 1000);
    }

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
    };
  }, [quizState.status, timeRemaining]);

  // Format time display
  const formatTime = useCallback((seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }, []);

  // Calculate progress
  const progress = useMemo(() => {
    const answeredQuestions = Object.keys(quizState.answers).length;
    return (answeredQuestions / content.questions.length) * 100;
  }, [quizState.answers, content.questions.length]);

  // Handle interaction
  const handleInteraction = useCallback((type: string, data: Record<string, unknown> = {}) => {
    if (!onInteraction) return;

    onInteraction({
      messageId: message.id,
      interactionType: type as any,
      data: { ...data, quizId: content.id },
      userId: currentUserId,
      timestamp: new Date()
    });
  }, [message.id, content.id, currentUserId, onInteraction]);

  // Start Quiz
  const startQuiz = useCallback(() => {
    setQuizState(prev => ({
      ...prev,
      status: QuizStatus.IN_PROGRESS,
      currentQuestion: 0
    }));
    handleInteraction('quiz_started');
  }, [handleInteraction]);

  // Answer Change Handler
  const handleAnswerChange = useCallback((questionId: string, answer: QuizAnswerValue) => {
    setQuizState(prev => ({
      ...prev,
      answers: { ...prev.answers, [questionId]: answer },
      canSubmit: Object.keys({ ...prev.answers, [questionId]: answer }).length > 0
    }));

    handleInteraction('answer_changed', {
      questionId,
      answer,
      timestamp: new Date()
    });
  }, [handleInteraction]);

  // Submit Quiz
  const submitQuiz = useCallback(() => {
    setQuizState(prev => ({ ...prev, status: QuizStatus.SUBMITTED }));
    handleInteraction('quiz_submitted', {
      answers: quizState.answers,
      timeSpent: content.timeLimit ? (content.timeLimit - (timeRemaining || 0)) : 0,
      completedAt: new Date()
    });
  }, [quizState.answers, content.timeLimit, timeRemaining, handleInteraction]);

  // Navigation
  const goToQuestion = useCallback((questionIndex: number) => {
    setQuizState(prev => ({ ...prev, currentQuestion: questionIndex }));
  }, []);

  const nextQuestion = useCallback(() => {
    if (quizState.currentQuestion < content.questions.length - 1) {
      goToQuestion(quizState.currentQuestion + 1);
    }
  }, [quizState.currentQuestion, content.questions.length, goToQuestion]);

  const prevQuestion = useCallback(() => {
    if (quizState.currentQuestion > 0) {
      goToQuestion(quizState.currentQuestion - 1);
    }
  }, [quizState.currentQuestion, goToQuestion]);

  const currentQuestion = content.questions[quizState.currentQuestion];

  return (
    <TooltipProvider>
      <motion.div
        variants={quizVariants}
        initial="initial"
        animate="animate"
        exit="exit"
        className={cn(
          "quiz-message w-full max-w-4xl mx-auto",
          className
        )}
        data-testid={`quiz-message-${message.id}`}
      >
        <Card className="overflow-hidden">
          {/* Quiz Header */}
          <CardHeader className="pb-4">
            <div className="flex items-start justify-between">
              <div className="space-y-2">
                <CardTitle className="text-xl">{content.title}</CardTitle>
                {content.description && (
                  <p className="text-muted-foreground">{content.description}</p>
                )}
              </div>
              <div className="flex items-center gap-2">
                {timeRemaining !== null && quizState.status === QuizStatus.IN_PROGRESS && (
                  <Badge variant={timeRemaining < 60 ? "destructive" : "secondary"}>
                    <Clock className="w-3 h-3 mr-1" />
                    {formatTime(timeRemaining)}
                  </Badge>
                )}
                <Badge variant="outline">
                  <Users className="w-3 h-3 mr-1" />
                  {content.questions.length} questions
                </Badge>
                <Badge variant="outline">
                  <Award className="w-3 h-3 mr-1" />
                  {content.questions.reduce((sum, q) => sum + q.points, 0)} points
                </Badge>
              </div>
            </div>

            {/* Progress Bar */}
            {quizState.status === QuizStatus.IN_PROGRESS && (
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Progress</span>
                  <span>{Math.round(progress)}% complete</span>
                </div>
                <Progress value={progress} className="h-2" />
              </div>
            )}
          </CardHeader>

          <Separator />

          <CardContent className="pt-6">
            {/* Not Started State */}
            {quizState.status === QuizStatus.NOT_STARTED && !readonly && (
              <div className="text-center space-y-6 py-8">
                <div className="space-y-2">
                  <h3 className="text-lg font-semibold">Ready to start the quiz?</h3>
                  <p className="text-muted-foreground max-w-md mx-auto">
                    {content.description || `This quiz contains ${content.questions.length} questions.`}
                    {content.timeLimit && (
                      <> You'll have {Math.floor(content.timeLimit / 60)} minutes to complete it.</>
                    )}
                  </p>
                </div>

                <div className="flex justify-center gap-3">
                  <Button onClick={startQuiz} size="lg">
                    <Play className="w-4 h-4 mr-2" />
                    Start Quiz
                  </Button>
                </div>

                {/* Quiz Settings Display */}
                <div className="grid grid-cols-2 gap-4 max-w-md mx-auto text-sm">
                  {content.timeLimit && (
                    <div className="text-center">
                      <Clock className="w-5 h-5 mx-auto mb-1 text-muted-foreground" />
                      <p className="font-medium">{Math.floor(content.timeLimit / 60)} minutes</p>
                      <p className="text-muted-foreground">Time limit</p>
                    </div>
                  )}
                  {content.passingScore && (
                    <div className="text-center">
                      <Award className="w-5 h-5 mx-auto mb-1 text-muted-foreground" />
                      <p className="font-medium">{content.passingScore}%</p>
                      <p className="text-muted-foreground">Passing score</p>
                    </div>
                  )}
                  {content.allowRetakes && (
                    <div className="text-center col-span-2">
                      <RotateCcw className="w-5 h-5 mx-auto mb-1 text-muted-foreground" />
                      <p className="text-muted-foreground">Retakes allowed</p>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* In Progress State */}
            {quizState.status === QuizStatus.IN_PROGRESS && currentQuestion && (
              <div className="space-y-6">
                {/* Question Counter */}
                <div className="flex justify-between items-center">
                  <div className="text-sm text-muted-foreground">
                    Question {quizState.currentQuestion + 1} of {content.questions.length}
                  </div>
                  <div className="flex gap-1">
                    {content.questions.map((_, index) => (
                      <button
                        key={index}
                        onClick={() => goToQuestion(index)}
                        className={cn(
                          "w-8 h-8 rounded-full text-xs font-medium transition-colors",
                          index === quizState.currentQuestion && "bg-primary text-primary-foreground",
                          index < quizState.currentQuestion && quizState.answers[content.questions[index].id] && "bg-green-100 text-green-700",
                          index !== quizState.currentQuestion && !quizState.answers[content.questions[index].id] && "bg-muted text-muted-foreground hover:bg-muted/80"
                        )}
                      >
                        {index + 1}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Current Question */}
                <QuestionComponent
                  question={currentQuestion}
                  answer={quizState.answers[currentQuestion.id] || null}
                  onAnswerChange={handleAnswerChange}
                  readonly={readonly}
                />

                {/* Navigation */}
                <div className="flex justify-between items-center pt-4">
                  <Button
                    variant="outline"
                    onClick={prevQuestion}
                    disabled={quizState.currentQuestion === 0}
                  >
                    <ChevronLeft className="w-4 h-4 mr-1" />
                    Previous
                  </Button>

                  <div className="flex gap-2">
                    {quizState.currentQuestion < content.questions.length - 1 ? (
                      <Button onClick={nextQuestion}>
                        Next
                        <ChevronRight className="w-4 h-4 ml-1" />
                      </Button>
                    ) : (
                      <Button
                        onClick={submitQuiz}
                        disabled={!quizState.canSubmit}
                        className="bg-green-600 hover:bg-green-700"
                      >
                        <Send className="w-4 h-4 mr-2" />
                        Submit Quiz
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            )}

            {/* Completed/Time Expired State */}
            {(quizState.status === QuizStatus.SUBMITTED || quizState.status === QuizStatus.TIME_EXPIRED) && (
              <div className="text-center space-y-6 py-8">
                <div className="space-y-2">
                  <h3 className="text-lg font-semibold">
                    {quizState.status === QuizStatus.TIME_EXPIRED ? "Time's up!" : "Quiz submitted!"}
                  </h3>
                  <p className="text-muted-foreground">
                    {quizState.status === QuizStatus.TIME_EXPIRED
                      ? "Your time has expired. Your answers have been automatically submitted."
                      : "Thank you for completing the quiz. Results will be available shortly."
                    }
                  </p>
                </div>

                {allowRetake && (
                  <Button
                    variant="outline"
                    onClick={() => setQuizState(prev => ({ ...prev, status: QuizStatus.NOT_STARTED, answers: {} }))}
                  >
                    <RotateCcw className="w-4 h-4 mr-2" />
                    Retake Quiz
                  </Button>
                )}
              </div>
            )}

            {/* Results Display */}
            {showResults && (
              <div className="space-y-6">
                <Separator />
                <div className="space-y-4">
                  <h3 className="text-lg font-semibold">Quiz Results</h3>
                  {content.questions.map((question, index) => (
                    <QuestionComponent
                      key={question.id}
                      question={question}
                      answer={quizState.answers[question.id] || null}
                      onAnswerChange={() => {}} // Read-only
                      readonly={true}
                      showResult={true}
                      isCorrect={quizState.answers[question.id] === question.correctAnswer}
                      explanation={question.explanation}
                    />
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </motion.div>
    </TooltipProvider>
  );
};

// Memoized version for performance
export const MemoizedQuizMessage = React.memo(QuizMessage);
MemoizedQuizMessage.displayName = 'MemoizedQuizMessage';

export default QuizMessage;