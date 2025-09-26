// T066 - Integration tests for MessageTypes API
/**
 * MessageTypes API Integration Tests
 * Comprehensive testing for RTK Query endpoints
 * Tests CRUD operations, error handling, and caching behavior
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { setupApiStore } from '../utils/test-store';
import { server } from '../mocks/server';
import { rest } from 'msw';
import {
  messageTypesApi,
  CreateMessageRequest,
  MessageInteractionRequest,
  MessageSearchRequest,
  MessageAnalyticsRequest,
} from '../../src/services/api/messageTypes';
import { MessageType, MessageData } from '../../src/types/MessageData';

// Mock data for testing
const mockMessageData: MessageData = {
  id: 'msg-123',
  senderId: 'user-123',
  senderName: 'John Doe',
  timestamp: new Date('2024-01-15T10:30:00Z'),
  type: MessageType.REPLY,
  isOwn: false,
  content: {
    originalMessageId: 'msg-456',
    replyText: 'This is a test reply',
    threadDepth: 1,
    isThreadStart: false,
  },
};

const mockQuizMessage: MessageData = {
  id: 'msg-quiz-123',
  senderId: 'user-456',
  senderName: 'Jane Smith',
  timestamp: new Date('2024-01-15T11:00:00Z'),
  type: MessageType.QUIZ,
  isOwn: true,
  content: {
    title: 'JavaScript Quiz',
    description: 'Test your JavaScript knowledge',
    questions: [
      {
        id: 'q1',
        type: 'multiple_choice',
        question: 'What is the output of typeof null?',
        options: ['null', 'object', 'undefined', 'string'],
        correctAnswer: 1,
        explanation: 'typeof null returns "object" due to a historical bug in JavaScript',
      },
    ],
    timeLimit: 300,
    showResults: true,
    allowRetake: false,
  },
};

const mockThread = {
  id: 'thread-123',
  rootMessageId: 'msg-456',
  messageCount: 5,
  participants: ['user-123', 'user-456', 'user-789'],
  lastActivity: '2024-01-15T12:00:00Z',
};

describe('MessageTypes API Integration Tests', () => {
  let storeRef: ReturnType<typeof setupApiStore>;

  beforeEach(() => {
    storeRef = setupApiStore(messageTypesApi);
  });

  afterEach(() => {
    // Clean up any ongoing queries
    storeRef.store.dispatch(messageTypesApi.util.resetApiState());
  });

  describe('createMessage', () => {
    it('should create a reply message successfully', async () => {
      // Arrange
      const createRequest: CreateMessageRequest = {
        chatId: 'chat-123',
        type: MessageType.REPLY,
        content: {
          originalMessageId: 'msg-456',
          replyText: 'This is a test reply',
          threadDepth: 1,
          isThreadStart: false,
        },
        replyToId: 'msg-456',
      };

      server.use(
        rest.post('/api/v1/messages', (req, res, ctx) => {
          return res(
            ctx.status(201),
            ctx.json({
              message: mockMessageData,
              thread: mockThread,
              notifications: [
                {
                  type: 'reply',
                  userId: 'user-456',
                  messageId: mockMessageData.id,
                  content: 'Someone replied to your message',
                },
              ],
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.createMessage.initiate(createRequest)
      );

      // Assert
      expect(result.data).toBeDefined();
      expect(result.data?.message.id).toBe('msg-123');
      expect(result.data?.message.type).toBe(MessageType.REPLY);
      expect(result.data?.thread).toBeDefined();
      expect(result.data?.notifications).toHaveLength(1);
    });

    it('should create a quiz message with complex content', async () => {
      // Arrange
      const createRequest: CreateMessageRequest = {
        chatId: 'chat-123',
        type: MessageType.QUIZ,
        content: mockQuizMessage.content,
      };

      server.use(
        rest.post('/api/v1/messages', (req, res, ctx) => {
          return res(
            ctx.status(201),
            ctx.json({
              message: mockQuizMessage,
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.createMessage.initiate(createRequest)
      );

      // Assert
      expect(result.data).toBeDefined();
      expect(result.data?.message.type).toBe(MessageType.QUIZ);
      expect(result.data?.message.content).toMatchObject({
        title: 'JavaScript Quiz',
        questions: expect.arrayContaining([
          expect.objectContaining({
            type: 'multiple_choice',
            options: expect.any(Array),
          }),
        ]),
      });
    });

    it('should handle validation errors', async () => {
      // Arrange
      const invalidRequest: CreateMessageRequest = {
        chatId: '',
        type: MessageType.REPLY,
        content: {},
      };

      server.use(
        rest.post('/api/v1/messages', (req, res, ctx) => {
          return res(
            ctx.status(400),
            ctx.json({
              error: 'Validation failed',
              details: [
                { field: 'chatId', message: 'Chat ID is required' },
                { field: 'content.replyText', message: 'Reply text is required' },
              ],
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.createMessage.initiate(invalidRequest)
      );

      // Assert
      expect(result.error).toBeDefined();
      expect(result.error).toMatchObject({
        status: 400,
        data: expect.objectContaining({
          error: 'Validation failed',
        }),
      });
    });

    it('should retry on transient failures', async () => {
      // Arrange
      let attempts = 0;
      const createRequest: CreateMessageRequest = {
        chatId: 'chat-123',
        type: MessageType.REPLY,
        content: {
          originalMessageId: 'msg-456',
          replyText: 'Retry test',
          threadDepth: 1,
          isThreadStart: false,
        },
      };

      server.use(
        rest.post('/api/v1/messages', (req, res, ctx) => {
          attempts++;
          if (attempts === 1) {
            return res(ctx.status(500), ctx.json({ error: 'Internal server error' }));
          }
          return res(ctx.status(201), ctx.json({ message: mockMessageData }));
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.createMessage.initiate(createRequest)
      );

      // Assert
      expect(attempts).toBe(2);
      expect(result.data).toBeDefined();
    });
  });

  describe('interactWithMessage', () => {
    it('should handle quiz answer submission', async () => {
      // Arrange
      const interactionRequest: MessageInteractionRequest = {
        messageId: 'msg-quiz-123',
        interactionType: 'quiz_answer',
        data: {
          questionId: 'q1',
          answer: 1,
          timeSpent: 15,
        },
      };

      server.use(
        rest.post('/api/v1/messages/msg-quiz-123/interactions', (req, res, ctx) => {
          return res(
            ctx.status(200),
            ctx.json({
              success: true,
              message: {
                ...mockQuizMessage,
                content: {
                  ...mockQuizMessage.content,
                  userAnswers: [
                    {
                      questionId: 'q1',
                      answer: 1,
                      isCorrect: true,
                      timeSpent: 15,
                    },
                  ],
                },
              },
              result: {
                score: 100,
                correctAnswers: 1,
                totalQuestions: 1,
                completedAt: new Date().toISOString(),
              },
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.interactWithMessage.initiate(interactionRequest)
      );

      // Assert
      expect(result.data).toBeDefined();
      expect(result.data?.success).toBe(true);
      expect(result.data?.result.score).toBe(100);
    });

    it('should handle poll voting', async () => {
      // Arrange
      const interactionRequest: MessageInteractionRequest = {
        messageId: 'msg-poll-123',
        interactionType: 'vote',
        data: {
          optionId: 'option-1',
        },
      };

      server.use(
        rest.post('/api/v1/messages/msg-poll-123/interactions', (req, res, ctx) => {
          return res(
            ctx.status(200),
            ctx.json({
              success: true,
              message: {
                id: 'msg-poll-123',
                type: MessageType.SURVEY,
                content: {
                  title: 'Team Lunch Vote',
                  options: [
                    { id: 'option-1', text: 'Pizza', votes: 5 },
                    { id: 'option-2', text: 'Sushi', votes: 3 },
                  ],
                  userVote: 'option-1',
                },
              },
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.interactWithMessage.initiate(interactionRequest)
      );

      // Assert
      expect(result.data?.success).toBe(true);
      expect(result.data?.message.content).toMatchObject({
        userVote: 'option-1',
        options: expect.arrayContaining([
          expect.objectContaining({ id: 'option-1', votes: 5 }),
        ]),
      });
    });
  });

  describe('searchMessages', () => {
    it('should search messages by type and content', async () => {
      // Arrange
      const searchRequest: MessageSearchRequest = {
        chatId: 'chat-123',
        query: 'JavaScript',
        messageTypes: [MessageType.QUIZ, MessageType.FORM],
        limit: 20,
      };

      server.use(
        rest.get('/api/v1/messages/search', (req, res, ctx) => {
          const url = req.url;
          expect(url.searchParams.get('query')).toBe('JavaScript');
          expect(url.searchParams.get('messageTypes')).toBe('quiz,form');

          return res(
            ctx.status(200),
            ctx.json({
              data: [mockQuizMessage],
              totalCount: 1,
              facets: {
                messageTypes: { quiz: 1, form: 0 },
                senders: { 'user-456': 1 },
                dateHistogram: { '2024-01-15': 1 },
                tags: { javascript: 1, programming: 1 },
              },
              suggestions: ['JavaScript basics', 'JavaScript advanced'],
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.searchMessages.initiate(searchRequest)
      );

      // Assert
      expect(result.data).toBeDefined();
      expect(result.data?.messages).toHaveLength(1);
      expect(result.data?.messages[0].type).toBe(MessageType.QUIZ);
      expect(result.data?.facets.messageTypes.quiz).toBe(1);
      expect(result.data?.suggestions).toContain('JavaScript basics');
    });

    it('should handle empty search results', async () => {
      // Arrange
      const searchRequest: MessageSearchRequest = {
        chatId: 'chat-123',
        query: 'nonexistent',
      };

      server.use(
        rest.get('/api/v1/messages/search', (req, res, ctx) => {
          return res(
            ctx.status(200),
            ctx.json({
              data: [],
              totalCount: 0,
              facets: {
                messageTypes: {},
                senders: {},
                dateHistogram: {},
                tags: {},
              },
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.searchMessages.initiate(searchRequest)
      );

      // Assert
      expect(result.data?.messages).toHaveLength(0);
      expect(result.data?.totalCount).toBe(0);
    });
  });

  describe('getMessageAnalytics', () => {
    it('should retrieve analytics for message types', async () => {
      // Arrange
      const analyticsRequest: MessageAnalyticsRequest = {
        chatId: 'chat-123',
        messageTypes: [MessageType.QUIZ, MessageType.SURVEY],
        timeRange: {
          start: '2024-01-01T00:00:00Z',
          end: '2024-01-31T23:59:59Z',
        },
        groupBy: 'day',
      };

      server.use(
        rest.get('/api/v1/messages/analytics', (req, res, ctx) => {
          return res(
            ctx.status(200),
            ctx.json({
              messageStats: [
                {
                  type: MessageType.QUIZ,
                  count: 15,
                  averageEngagement: 0.85,
                  completionRate: 0.92,
                  errorRate: 0.03,
                },
                {
                  type: MessageType.SURVEY,
                  count: 8,
                  averageEngagement: 0.76,
                  completionRate: 0.88,
                  errorRate: 0.01,
                },
              ],
              engagementMetrics: {
                totalInteractions: 450,
                uniqueUsers: 23,
                averageResponseTime: 125.5,
                retentionRate: 0.78,
              },
              trends: [
                {
                  date: '2024-01-15',
                  messageCount: 3,
                  engagementScore: 0.89,
                  completionRate: 0.94,
                },
              ],
              topPerformers: [
                {
                  messageId: 'msg-quiz-123',
                  type: MessageType.QUIZ,
                  engagementScore: 0.95,
                  interactionCount: 87,
                  completionRate: 0.98,
                },
              ],
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.getMessageAnalytics.initiate(analyticsRequest)
      );

      // Assert
      expect(result.data).toBeDefined();
      expect(result.data?.messageStats).toHaveLength(2);
      expect(result.data?.messageStats[0].type).toBe(MessageType.QUIZ);
      expect(result.data?.messageStats[0].completionRate).toBe(0.92);
      expect(result.data?.engagementMetrics.uniqueUsers).toBe(23);
      expect(result.data?.topPerformers).toHaveLength(1);
    });
  });

  describe('validateMessage', () => {
    it('should validate quiz message content', async () => {
      // Arrange
      const validationRequest = {
        type: MessageType.QUIZ,
        content: mockQuizMessage.content,
        chatId: 'chat-123',
      };

      server.use(
        rest.post('/api/v1/messages/validate', (req, res, ctx) => {
          return res(
            ctx.status(200),
            ctx.json({
              isValid: true,
              errors: [],
              warnings: [
                {
                  field: 'timeLimit',
                  code: 'RECOMMENDED_VALUE',
                  message: 'Consider extending time limit for better user experience',
                  suggestion: 'Increase to 600 seconds for complex quizzes',
                },
              ],
              sanitizedContent: validationRequest.content,
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.validateMessage.initiate(validationRequest)
      );

      // Assert
      expect(result.data?.isValid).toBe(true);
      expect(result.data?.errors).toHaveLength(0);
      expect(result.data?.warnings).toHaveLength(1);
      expect(result.data?.warnings[0].suggestion).toContain('600 seconds');
    });

    it('should return validation errors for invalid content', async () => {
      // Arrange
      const invalidRequest = {
        type: MessageType.QUIZ,
        content: { title: '' }, // Invalid: missing required fields
        chatId: 'chat-123',
      };

      server.use(
        rest.post('/api/v1/messages/validate', (req, res, ctx) => {
          return res(
            ctx.status(200),
            ctx.json({
              isValid: false,
              errors: [
                {
                  field: 'title',
                  code: 'REQUIRED_FIELD',
                  message: 'Title cannot be empty',
                  severity: 'error',
                },
                {
                  field: 'questions',
                  code: 'REQUIRED_FIELD',
                  message: 'At least one question is required',
                  severity: 'critical',
                },
              ],
              warnings: [],
            })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.validateMessage.initiate(invalidRequest)
      );

      // Assert
      expect(result.data?.isValid).toBe(false);
      expect(result.data?.errors).toHaveLength(2);
      expect(result.data?.errors.find(e => e.severity === 'critical')).toBeTruthy();
    });
  });

  describe('caching behavior', () => {
    it('should cache search results appropriately', async () => {
      // Arrange
      const searchRequest: MessageSearchRequest = {
        chatId: 'chat-123',
        query: 'test',
      };

      let requestCount = 0;
      server.use(
        rest.get('/api/v1/messages/search', (req, res, ctx) => {
          requestCount++;
          return res(
            ctx.status(200),
            ctx.json({
              data: [mockMessageData],
              totalCount: 1,
              facets: { messageTypes: {}, senders: {}, dateHistogram: {}, tags: {} },
            })
          );
        })
      );

      // Act
      await storeRef.store.dispatch(
        messageTypesApi.endpoints.searchMessages.initiate(searchRequest)
      );

      await storeRef.store.dispatch(
        messageTypesApi.endpoints.searchMessages.initiate(searchRequest)
      );

      // Assert - Second request should use cache
      expect(requestCount).toBe(1);
    });

    it('should invalidate cache on message creation', async () => {
      // Arrange
      const chatId = 'chat-123';

      // First, populate cache with messages
      server.use(
        rest.get('/api/v1/messages', (req, res, ctx) => {
          return res(
            ctx.status(200),
            ctx.json({
              data: [mockMessageData],
              totalCount: 1,
            })
          );
        }),
        rest.post('/api/v1/messages', (req, res, ctx) => {
          return res(
            ctx.status(201),
            ctx.json({
              message: {
                ...mockMessageData,
                id: 'new-msg-123',
              },
            })
          );
        })
      );

      // Get initial messages (populates cache)
      await storeRef.store.dispatch(
        messageTypesApi.endpoints.getMessages.initiate({ chatId })
      );

      // Act - Create new message (should invalidate cache)
      await storeRef.store.dispatch(
        messageTypesApi.endpoints.createMessage.initiate({
          chatId,
          type: MessageType.REPLY,
          content: mockMessageData.content,
        })
      );

      // Assert - Cache should be invalidated for the chat
      const state = storeRef.store.getState();
      const queryState = state.messageTypesApi.queries[`getMessages({"chatId":"${chatId}"})`];

      // The query should be invalidated (status changed)
      expect(queryState?.status).not.toBe('fulfilled');
    });
  });

  describe('error handling', () => {
    it('should handle network errors gracefully', async () => {
      // Arrange
      server.use(
        rest.post('/api/v1/messages', (req, res, ctx) => {
          return res.networkError('Network connection failed');
        })
      );

      const createRequest: CreateMessageRequest = {
        chatId: 'chat-123',
        type: MessageType.REPLY,
        content: mockMessageData.content,
      };

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.createMessage.initiate(createRequest)
      );

      // Assert
      expect(result.error).toBeDefined();
      expect(result.error).toMatchObject({
        name: 'FetchError',
      });
    });

    it('should handle rate limiting with retry', async () => {
      // Arrange
      let attempts = 0;
      const createRequest: CreateMessageRequest = {
        chatId: 'chat-123',
        type: MessageType.REPLY,
        content: mockMessageData.content,
      };

      server.use(
        rest.post('/api/v1/messages', (req, res, ctx) => {
          attempts++;
          if (attempts === 1) {
            return res(
              ctx.status(429),
              ctx.json({ error: 'Rate limit exceeded' })
            );
          }
          return res(
            ctx.status(201),
            ctx.json({ message: mockMessageData })
          );
        })
      );

      // Act
      const result = await storeRef.store.dispatch(
        messageTypesApi.endpoints.createMessage.initiate(createRequest)
      );

      // Assert - Should succeed after retry
      expect(attempts).toBe(2);
      expect(result.data).toBeDefined();
    });
  });
});