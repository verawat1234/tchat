// T031 - Quiz content and interaction interfaces
/**
 * Type definitions for interactive quiz message content
 * Supports multiple question types, scoring, analytics, and real-time interaction
 */

// Core Quiz Content Interface
export interface QuizContent {
  readonly id: string;
  readonly title: string;
  readonly description?: string;
  readonly questions: QuizQuestion[];
  readonly timeLimit?: number; // in seconds
  readonly passingScore?: number; // percentage (0-100)
  readonly allowRetakes: boolean;
  readonly showAnswersAfter: ShowAnswersPolicy;
  readonly quizSettings: QuizSettings;
  readonly analytics?: QuizAnalytics;
}

// Quiz Settings Configuration
export interface QuizSettings {
  readonly randomizeQuestions: boolean;
  readonly randomizeAnswers: boolean;
  readonly showProgressBar: boolean;
  readonly allowSkipQuestions: boolean;
  readonly showQuestionNumbers: boolean;
  readonly requireAllQuestions: boolean;
  readonly enableHints: boolean;
  readonly enableReviewMode: boolean;
  readonly maxAttempts?: number;
  readonly lockAfterSubmission: boolean;
}

// When to Show Answers
export enum ShowAnswersPolicy {
  NEVER = 'never',
  AFTER_SUBMISSION = 'submission',
  AFTER_COMPLETION = 'completion',
  AFTER_TIME_LIMIT = 'time_limit',
  IMMEDIATELY = 'immediately'
}

// Quiz Question Definition
export interface QuizQuestion {
  readonly id: string;
  readonly question: string;
  readonly type: QuizQuestionType;
  readonly options?: string[];
  readonly correctAnswer: string | string[] | number | boolean;
  readonly explanation?: string;
  readonly points: number;
  readonly hints?: QuizHint[];
  readonly media?: QuizMedia;
  readonly timeLimit?: number;
  readonly tags?: string[];
}

// Question Types
export enum QuizQuestionType {
  MULTIPLE_CHOICE = 'multiple_choice',
  MULTIPLE_SELECT = 'multiple_select',
  TRUE_FALSE = 'true_false',
  SHORT_ANSWER = 'short_answer',
  LONG_ANSWER = 'long_answer',
  NUMBER = 'number',
  SCALE = 'scale',
  ORDERING = 'ordering',
  MATCHING = 'matching',
  FILL_BLANK = 'fill_blank',
  DRAG_DROP = 'drag_drop'
}

// Quiz Hints
export interface QuizHint {
  readonly id: string;
  readonly text: string;
  readonly costInPoints?: number;
  readonly availableAfter?: number; // seconds into question
  readonly maxUses?: number;
}

// Quiz Media Attachments
export interface QuizMedia {
  readonly type: QuizMediaType;
  readonly url: string;
  readonly thumbnailUrl?: string;
  readonly alt?: string;
  readonly caption?: string;
  readonly dimensions?: MediaDimensions;
}

export enum QuizMediaType {
  IMAGE = 'image',
  VIDEO = 'video',
  AUDIO = 'audio',
  DOCUMENT = 'document'
}

export interface MediaDimensions {
  readonly width: number;
  readonly height: number;
}

// Quiz Submission Data
export interface QuizSubmission {
  readonly quizId: string;
  readonly userId: string;
  readonly submissionId: string;
  readonly answers: QuizAnswerMap;
  readonly startTime: Date;
  readonly endTime: Date;
  readonly timeSpent: number; // in seconds
  readonly hintsUsed: HintUsage[];
  readonly questionOrder: string[];
  readonly isCompleted: boolean;
  readonly submissionData: QuizSubmissionMetadata;
}

export type QuizAnswerMap = Record<string, QuizAnswerValue>;

export type QuizAnswerValue = string | string[] | number | boolean | null;

export interface HintUsage {
  readonly questionId: string;
  readonly hintId: string;
  readonly usedAt: Date;
  readonly pointsCost: number;
}

export interface QuizSubmissionMetadata {
  readonly userAgent: string;
  readonly ipAddress?: string;
  readonly screenResolution?: string;
  readonly timezone: string;
  readonly language: string;
  readonly sessionId: string;
}

// Quiz Results and Scoring
export interface QuizResult {
  readonly submissionId: string;
  readonly quizId: string;
  readonly userId: string;
  readonly totalPoints: number;
  readonly earnedPoints: number;
  readonly percentage: number;
  readonly passed: boolean;
  readonly grade?: QuizGrade;
  readonly questionResults: QuestionResult[];
  readonly timingAnalysis: TimingAnalysis;
  readonly feedback?: QuizFeedback;
}

export interface QuestionResult {
  readonly questionId: string;
  readonly userAnswer: QuizAnswerValue;
  readonly correctAnswer: string | string[] | number | boolean;
  readonly isCorrect: boolean;
  readonly pointsEarned: number;
  readonly timeSpent: number;
  readonly hintsUsed: number;
  readonly partialCredit?: number;
  readonly feedback?: string;
}

export interface TimingAnalysis {
  readonly totalTime: number;
  readonly averageTimePerQuestion: number;
  readonly fastestQuestion: { questionId: string; time: number };
  readonly slowestQuestion: { questionId: string; time: number };
  readonly timeDistribution: number[];
}

export interface QuizFeedback {
  readonly overall: string;
  readonly strengths: string[];
  readonly improvements: string[];
  readonly recommendedResources: Resource[];
  readonly nextSteps: string[];
}

export interface Resource {
  readonly title: string;
  readonly type: ResourceType;
  readonly url: string;
  readonly description?: string;
}

export enum ResourceType {
  ARTICLE = 'article',
  VIDEO = 'video',
  COURSE = 'course',
  BOOK = 'book',
  PRACTICE = 'practice'
}

// Quiz Grading System
export interface QuizGrade {
  readonly letter: string; // A, B, C, D, F
  readonly gpa: number; // 4.0 scale
  readonly description: string;
  readonly color: string; // hex color for display
}

// Real-time Quiz State
export interface QuizState {
  readonly quizId: string;
  readonly status: QuizStatus;
  readonly currentQuestion: number;
  readonly totalQuestions: number;
  readonly timeRemaining?: number;
  readonly answers: QuizAnswerMap;
  readonly flaggedQuestions: string[];
  readonly reviewMode: boolean;
  readonly canSubmit: boolean;
  readonly validationErrors: QuizValidationError[];
}

export enum QuizStatus {
  NOT_STARTED = 'not_started',
  IN_PROGRESS = 'in_progress',
  PAUSED = 'paused',
  COMPLETED = 'completed',
  TIME_EXPIRED = 'time_expired',
  SUBMITTED = 'submitted'
}

export interface QuizValidationError {
  readonly questionId?: string;
  readonly field: string;
  readonly message: string;
  readonly severity: ValidationSeverity;
}

export enum ValidationSeverity {
  INFO = 'info',
  WARNING = 'warning',
  ERROR = 'error'
}

// Quiz Analytics and Reporting
export interface QuizAnalytics {
  readonly participation: ParticipationStats;
  readonly performance: PerformanceStats;
  readonly questionAnalysis: QuestionAnalysisStats[];
  readonly timing: TimingStats;
  readonly completion: CompletionStats;
}

export interface ParticipationStats {
  readonly totalAttempts: number;
  readonly uniqueParticipants: number;
  readonly completionRate: number;
  readonly averageAttempts: number;
  readonly peakParticipationTime: Date;
  readonly participantDemographics?: Record<string, number>;
}

export interface PerformanceStats {
  readonly averageScore: number;
  readonly medianScore: number;
  readonly highestScore: number;
  readonly lowestScore: number;
  readonly passRate: number;
  readonly scoreDistribution: ScoreDistribution;
  readonly improvementOverTime: number[];
}

export interface ScoreDistribution {
  readonly ranges: ScoreRange[];
  readonly standardDeviation: number;
  readonly percentiles: Percentile[];
}

export interface ScoreRange {
  readonly min: number;
  readonly max: number;
  readonly count: number;
  readonly percentage: number;
}

export interface Percentile {
  readonly percentile: number;
  readonly score: number;
}

export interface QuestionAnalysisStats {
  readonly questionId: string;
  readonly difficulty: DifficultyLevel;
  readonly discrimination: number;
  readonly correctAnswerRate: number;
  readonly averageTime: number;
  readonly skipRate: number;
  readonly hintUsageRate: number;
  readonly answerDistribution: AnswerDistribution;
  readonly commonMistakes: string[];
}

export enum DifficultyLevel {
  VERY_EASY = 'very_easy',
  EASY = 'easy',
  MODERATE = 'moderate',
  HARD = 'hard',
  VERY_HARD = 'very_hard'
}

export interface AnswerDistribution {
  readonly answers: Record<string, AnswerStats>;
  readonly mostCommonWrong: string;
  readonly leastCommonCorrect: string;
}

export interface AnswerStats {
  readonly count: number;
  readonly percentage: number;
  readonly isCorrect: boolean;
  readonly averageConfidence?: number;
}

export interface TimingStats {
  readonly averageCompletionTime: number;
  readonly medianCompletionTime: number;
  readonly fastestCompletion: number;
  readonly slowestCompletion: number;
  readonly timeoutRate: number;
  readonly optimalTimeRange: { min: number; max: number };
}

export interface CompletionStats {
  readonly totalStarted: number;
  readonly totalCompleted: number;
  readonly completionRate: number;
  readonly averageProgressWhenAbandoned: number;
  readonly dropoffPoints: DropoffPoint[];
}

export interface DropoffPoint {
  readonly questionNumber: number;
  readonly dropoffRate: number;
  readonly reasons?: string[];
}

// Quiz Configuration Templates
export interface QuizTemplate {
  readonly id: string;
  readonly name: string;
  readonly description: string;
  readonly category: QuizCategory;
  readonly settings: QuizSettings;
  readonly questionTemplates: QuestionTemplate[];
  readonly isPublic: boolean;
  readonly createdBy: string;
  readonly usage: number;
}

export enum QuizCategory {
  EDUCATION = 'education',
  TRAINING = 'training',
  ASSESSMENT = 'assessment',
  SURVEY = 'survey',
  CERTIFICATION = 'certification',
  ENTERTAINMENT = 'entertainment',
  MARKET_RESEARCH = 'market_research'
}

export interface QuestionTemplate {
  readonly type: QuizQuestionType;
  readonly title: string;
  readonly structure: QuestionStructure;
  readonly validationRules: ValidationRule[];
  readonly scoringMethod: ScoringMethod;
}

export interface QuestionStructure {
  readonly requiredFields: string[];
  readonly optionalFields: string[];
  readonly constraints: Record<string, unknown>;
  readonly defaultValues: Record<string, unknown>;
}

export interface ValidationRule {
  readonly field: string;
  readonly rule: string;
  readonly message: string;
  readonly parameters?: Record<string, unknown>;
}

export enum ScoringMethod {
  ALL_OR_NOTHING = 'all_or_nothing',
  PARTIAL_CREDIT = 'partial_credit',
  WEIGHTED = 'weighted',
  NEGATIVE_MARKING = 'negative_marking',
  CUSTOM = 'custom'
}

// Quiz Interaction Events
export interface QuizInteractionEvent {
  readonly type: QuizInteractionType;
  readonly quizId: string;
  readonly questionId?: string;
  readonly userId: string;
  readonly timestamp: Date;
  readonly data?: Record<string, unknown>;
}

export enum QuizInteractionType {
  QUIZ_STARTED = 'quiz_started',
  QUESTION_VIEWED = 'question_viewed',
  ANSWER_CHANGED = 'answer_changed',
  QUESTION_FLAGGED = 'question_flagged',
  HINT_REQUESTED = 'hint_requested',
  QUESTION_SKIPPED = 'question_skipped',
  QUIZ_PAUSED = 'quiz_paused',
  QUIZ_RESUMED = 'quiz_resumed',
  QUIZ_SUBMITTED = 'quiz_submitted',
  REVIEW_MODE_ENTERED = 'review_mode_entered'
}

// Utility Functions and Type Guards
export const isValidQuizContent = (content: unknown): content is QuizContent => {
  return (
    typeof content === 'object' &&
    content !== null &&
    'id' in content &&
    'title' in content &&
    'questions' in content &&
    'allowRetakes' in content &&
    'showAnswersAfter' in content &&
    Array.isArray((content as any).questions) &&
    (content as any).questions.length > 0
  );
};

export const calculateQuizScore = (
  questions: QuizQuestion[],
  answers: QuizAnswerMap,
  scoringMethod: ScoringMethod = ScoringMethod.ALL_OR_NOTHING
): QuizResult => {
  const questionResults: QuestionResult[] = [];
  let totalPoints = 0;
  let earnedPoints = 0;

  questions.forEach(question => {
    totalPoints += question.points;
    const userAnswer = answers[question.id];
    const isCorrect = compareAnswers(question.correctAnswer, userAnswer, question.type);

    let pointsEarned = 0;
    if (isCorrect) {
      pointsEarned = question.points;
    } else if (scoringMethod === ScoringMethod.PARTIAL_CREDIT) {
      pointsEarned = calculatePartialCredit(question, userAnswer);
    }

    earnedPoints += pointsEarned;

    questionResults.push({
      questionId: question.id,
      userAnswer,
      correctAnswer: question.correctAnswer,
      isCorrect,
      pointsEarned,
      timeSpent: 0, // Would need to be tracked separately
      hintsUsed: 0, // Would need to be tracked separately
    });
  });

  const percentage = totalPoints > 0 ? (earnedPoints / totalPoints) * 100 : 0;

  return {
    submissionId: '', // Would be generated
    quizId: '', // Would be provided
    userId: '', // Would be provided
    totalPoints,
    earnedPoints,
    percentage,
    passed: percentage >= 70, // Default passing score
    questionResults,
    timingAnalysis: {
      totalTime: 0,
      averageTimePerQuestion: 0,
      fastestQuestion: { questionId: '', time: 0 },
      slowestQuestion: { questionId: '', time: 0 },
      timeDistribution: []
    }
  };
};

const compareAnswers = (
  correct: string | string[] | number | boolean,
  user: QuizAnswerValue,
  type: QuizQuestionType
): boolean => {
  if (user === null || user === undefined) return false;

  switch (type) {
    case QuizQuestionType.MULTIPLE_SELECT:
      return Array.isArray(correct) && Array.isArray(user) &&
        correct.length === user.length &&
        correct.every(answer => user.includes(answer));

    case QuizQuestionType.TRUE_FALSE:
      return correct === user;

    case QuizQuestionType.NUMBER:
      return typeof correct === 'number' && typeof user === 'number' &&
        Math.abs(correct - user) < 0.001; // Account for floating point precision

    case QuizQuestionType.SHORT_ANSWER:
    case QuizQuestionType.LONG_ANSWER:
      return typeof correct === 'string' && typeof user === 'string' &&
        correct.toLowerCase().trim() === user.toLowerCase().trim();

    default:
      return correct === user;
  }
};

const calculatePartialCredit = (
  question: QuizQuestion,
  userAnswer: QuizAnswerValue
): number => {
  // Partial credit calculation logic would depend on question type and requirements
  // This is a simplified implementation
  if (question.type === QuizQuestionType.MULTIPLE_SELECT &&
      Array.isArray(question.correctAnswer) &&
      Array.isArray(userAnswer)) {

    const correctAnswers = question.correctAnswer as string[];
    const userAnswers = userAnswer as string[];

    const correctSelected = userAnswers.filter(answer => correctAnswers.includes(answer)).length;
    const incorrectSelected = userAnswers.filter(answer => !correctAnswers.includes(answer)).length;

    const maxPoints = question.points;
    const partialCredit = (correctSelected / correctAnswers.length) * maxPoints;
    const penalty = (incorrectSelected / correctAnswers.length) * maxPoints;

    return Math.max(0, partialCredit - penalty);
  }

  return 0;
};

// Export all types for external use
export type {
  QuizSettings,
  QuizQuestion,
  QuizHint,
  QuizMedia,
  QuizSubmission,
  QuizResult,
  QuestionResult,
  TimingAnalysis,
  QuizFeedback,
  QuizGrade,
  QuizState,
  QuizValidationError,
  QuizAnalytics,
  QuizTemplate,
  QuizInteractionEvent
};