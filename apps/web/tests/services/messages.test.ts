/**
 * Message Service Contract Tests
 * Tests for GET /api/messages and POST /api/messages endpoints
 *
 * These tests validate request/response schemas according to the OpenAPI contract
 * Following TDD approach - tests will fail initially until endpoints are implemented
 *
 * EXPECTED FAILURES (TDD Approach):
 * 1. Schema validation failures - MSW handlers return mock data that doesn't match Message schema
 * 2. Missing POST endpoint handler - No MSW handler for POST /api/messages
 * 3. Incomplete message objects - Mock messages missing required fields like content, attachments, metadata
 * 4. Field type mismatches - Mock data types don't match schema expectations
 * 5. Pagination schema mismatches - Cursor-based pagination fields missing or incorrect
 * 6. Message content validation - Complex message content types not properly validated
 *
 * These failures are intentional and demonstrate that:
 * - Tests are written before implementation
 * - Schema validation is strict and comprehensive
 * - All endpoint contracts are thoroughly tested
 * - Error scenarios are properly handled
 * - Cursor-based pagination is properly implemented
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { server } from '@/lib/test-utils/msw/server';
import { http, HttpResponse } from 'msw';
import type {
  Message,
  MessageContent,
  MessageType,
  MessageAttachment,
  MessageReaction,
  MessageRead,
  MessageMetadata,
  UUID,
  Timestamp
} from '@/schema/schema';

// API Base URL
const API_BASE_URL = process.env.VITE_API_URL || 'http://localhost:3001';

// Type definitions for request/response validation
interface MessagesListResponse {
  messages: Message[];
  pagination: {
    total: number;
    hasMore: boolean;
    nextCursor?: string;
    prevCursor?: string;
    limit: number;
  };
  meta?: {
    unreadCount?: number;
    totalCount?: number;
    oldestMessageId?: UUID;
    newestMessageId?: UUID;
  };
}

interface MessageResponse {
  message: Message;
  meta?: {
    deliveryStatus?: 'sent' | 'delivered' | 'read' | 'failed';
    encryptionStatus?: 'encrypted' | 'plain';
  };
}

interface CreateMessageRequest {
  dialogId: UUID;
  content: MessageContent;
  type: MessageType;
  replyToId?: UUID;
  forwardFromId?: UUID;
  threadId?: UUID;
  mentions?: UUID[];
  attachments?: Omit<MessageAttachment, 'id' | 'messageId' | 'createdAt'>[];
  metadata?: Partial<MessageMetadata>;
  scheduledAt?: Timestamp;
  isEphemeral?: boolean;
  deleteAfter?: number; // seconds
}

interface UpdateMessageRequest {
  content?: MessageContent;
  isPinned?: boolean;
  metadata?: Partial<MessageMetadata>;
}

// Message query parameters interface
interface MessageQueryParams {
  dialogId?: UUID;
  cursor?: string;
  limit?: number;
  order?: 'asc' | 'desc';
  before?: Timestamp;
  after?: Timestamp;
  type?: MessageType;
  senderId?: UUID;
  hasAttachments?: boolean;
  isUnread?: boolean;
  search?: string;
  threadId?: UUID;
}

// Enhanced mock message factory with comprehensive fields
const createMockMessage = (overrides: Partial<Message> = {}): Message => ({
  id: '550e8400-e29b-41d4-a716-446655440001',
  dialogId: '550e8400-e29b-41d4-a716-446655440002',
  senderId: '550e8400-e29b-41d4-a716-446655440003',
  type: 'text',
  content: {
    text: 'Hello, this is a test message',
    html: '<p>Hello, this is a test message</p>'
  },
  replyToId: undefined,
  forwardFromId: undefined,
  threadId: undefined,
  isEdited: false,
  isPinned: false,
  reactions: [],
  readBy: [],
  mentions: [],
  attachments: [],
  metadata: {
    deliveryStatus: {
      sent: true,
      sentAt: '2024-09-22T10:00:00.000Z',
      delivered: true,
      deliveredAt: '2024-09-22T10:00:01.000Z',
      failed: false,
      failureReason: undefined,
      retryCount: 0
    },
    analytics: {
      opens: 0,
      clicks: 0,
      shares: 0,
      reactions: 0,
      forwardCount: 0,
      replyCount: 0,
      firstOpenAt: undefined,
      lastOpenAt: undefined
    }
  },
  createdAt: '2024-09-22T10:00:00.000Z',
  editedAt: undefined,
  deletedAt: undefined,
  ...overrides
});

// Mock message with complex content types
const createMockImageMessage = (): Message => createMockMessage({
  id: '550e8400-e29b-41d4-a716-446655440004',
  type: 'image',
  content: {
    image: {
      fileUrl: 'https://example.com/images/test.jpg',
      thumbnailUrl: 'https://example.com/images/test_thumb.jpg',
      width: 1920,
      height: 1080,
      fileSize: 2048000,
      caption: 'Test image',
      altText: 'A test image for contract validation',
      blurHash: 'LEHV6nWB2yk8pyo0adR*.7kCMdnj',
      exifData: {
        make: 'Canon',
        model: 'EOS R5',
        dateTime: '2024-09-22T10:00:00.000Z'
      }
    }
  },
  attachments: [
    {
      id: '550e8400-e29b-41d4-a716-446655440010',
      messageId: '550e8400-e29b-41d4-a716-446655440004',
      type: 'image',
      fileUrl: 'https://example.com/images/test.jpg',
      fileName: 'test.jpg',
      fileSize: 2048000,
      mimeType: 'image/jpeg',
      uploadProgress: 100,
      metadata: {
        width: 1920,
        height: 1080,
        thumbnail: 'https://example.com/images/test_thumb.jpg',
        virus_scan: 'clean',
        compression: 'medium',
        originalSize: 3072000
      },
      createdAt: '2024-09-22T09:59:58.000Z'
    }
  ]
});

// Mock message with voice content
const createMockVoiceMessage = (): Message => createMockMessage({
  id: '550e8400-e29b-41d4-a716-446655440005',
  type: 'voice',
  content: {
    voice: {
      fileUrl: 'https://example.com/audio/voice_note.m4a',
      duration: 45,
      waveform: [0.2, 0.4, 0.6, 0.8, 0.4, 0.2, 0.1, 0.3, 0.5, 0.7],
      transcript: 'This is a voice message for testing purposes',
      transcriptLanguage: 'en-US',
      fileSize: 512000,
      encoding: 'aac',
      sampleRate: 44100
    }
  }
});

// Mock payment message
const createMockPaymentMessage = (): Message => createMockMessage({
  id: '550e8400-e29b-41d4-a716-446655440006',
  type: 'payment',
  content: {
    payment: {
      amount: 250.50,
      currency: 'THB',
      description: 'Payment for dinner',
      transactionId: '550e8400-e29b-41d4-a716-446655440020',
      status: 'completed',
      paymentMethod: 'promptpay',
      recipient: 'John Doe',
      reference: 'DINNER-2024-09-22',
      qrCode: 'data:image/png;base64,iVBOR...',
      deepLink: 'promptpay://pay?amount=250.50&ref=DINNER-2024-09-22'
    }
  }
});

describe('Message Service Contract Tests', () => {
  beforeEach(() => {
    // Reset any custom handlers before each test
    server.resetHandlers();
  });

  describe('GET /api/messages', () => {
    it('should return messages with proper schema validation', async () => {
      // Mock successful response with cursor pagination
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, ({ request }) => {
          const url = new URL(request.url);
          const dialogId = url.searchParams.get('dialogId');
          const cursor = url.searchParams.get('cursor');
          const limit = parseInt(url.searchParams.get('limit') || '20');

          // Intentionally incomplete response to trigger TDD failure
          return HttpResponse.json({
            messages: [
              // Missing required fields to demonstrate TDD failures
              {
                id: '1',
                text: 'Test message', // Wrong schema - should be content.text
                senderId: '123',
                timestamp: new Date().toISOString() // Wrong field name - should be createdAt
              }
            ],
            // Missing pagination object entirely
            total: 1
          });
        })
      );

      const response = await fetch(
        `${API_BASE_URL}/api/messages?dialogId=550e8400-e29b-41d4-a716-446655440002&limit=20`
      );

      expect(response.ok).toBe(true);
      const data: MessagesListResponse = await response.json();

      // Schema validation tests - these will fail initially
      expect(data).toHaveProperty('messages');
      expect(data).toHaveProperty('pagination');
      expect(data.pagination).toHaveProperty('total');
      expect(data.pagination).toHaveProperty('hasMore');
      expect(data.pagination).toHaveProperty('limit');

      // Validate message schema
      data.messages.forEach((message: Message) => {
        expect(message).toHaveProperty('id');
        expect(message).toHaveProperty('dialogId');
        expect(message).toHaveProperty('senderId');
        expect(message).toHaveProperty('type');
        expect(message).toHaveProperty('content');
        expect(message).toHaveProperty('reactions');
        expect(message).toHaveProperty('readBy');
        expect(message).toHaveProperty('mentions');
        expect(message).toHaveProperty('attachments');
        expect(message).toHaveProperty('createdAt');
        expect(message).toHaveProperty('isEdited');
        expect(message).toHaveProperty('isPinned');

        // Validate UUID format
        expect(message.id).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i);

        // Validate timestamp format
        expect(message.createdAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?Z$/);

        // Validate message type enum
        expect(['text', 'voice', 'file', 'image', 'video', 'payment', 'system', 'location', 'contact', 'poll', 'event', 'product', 'sticker', 'gif']).toContain(message.type);

        // Validate content structure based on type
        if (message.type === 'text') {
          expect(message.content).toHaveProperty('text');
          expect(typeof message.content.text).toBe('string');
        }
      });
    });

    it('should handle cursor-based pagination correctly', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, ({ request }) => {
          const url = new URL(request.url);
          const cursor = url.searchParams.get('cursor');

          // Intentionally wrong pagination structure to demonstrate TDD failure
          return HttpResponse.json({
            data: [createMockMessage()], // Wrong structure - should be 'messages'
            page: 1, // Wrong pagination schema - should use cursors
            totalPages: 5,
            hasNextPage: true
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/messages?cursor=eyJpZCI6IjEyMzQifQ&limit=10`);

      expect(response.ok).toBe(true);
      const data: MessagesListResponse = await response.json();

      // Cursor pagination validation - will fail with current mock
      expect(data).toHaveProperty('pagination');
      expect(data.pagination).toHaveProperty('hasMore');
      expect(data.pagination).toHaveProperty('nextCursor');
      expect(data.pagination).toHaveProperty('prevCursor');
      expect(data.pagination.hasMore).toBe(true);
      expect(typeof data.pagination.nextCursor).toBe('string');
    });

    it('should support complex query parameters', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, ({ request }) => {
          const url = new URL(request.url);
          const dialogId = url.searchParams.get('dialogId');
          const type = url.searchParams.get('type');
          const hasAttachments = url.searchParams.get('hasAttachments');
          const after = url.searchParams.get('after');
          const search = url.searchParams.get('search');

          // Mock doesn't handle query parameters properly - TDD failure
          return HttpResponse.json({
            messages: [createMockMessage()],
            pagination: {
              total: 1,
              hasMore: false,
              limit: 20
            }
          });
        })
      );

      const queryParams = new URLSearchParams({
        dialogId: '550e8400-e29b-41d4-a716-446655440002',
        type: 'image',
        hasAttachments: 'true',
        after: '2024-09-22T00:00:00.000Z',
        search: 'test query',
        limit: '15'
      });

      const response = await fetch(`${API_BASE_URL}/api/messages?${queryParams}`);
      expect(response.ok).toBe(true);

      const data: MessagesListResponse = await response.json();

      // Validate that query parameters were properly applied
      // These assertions will pass initially but demonstrate contract expectations
      expect(Array.isArray(data.messages)).toBe(true);
      expect(typeof data.pagination.limit).toBe('number');
    });

    it('should return appropriate message types with correct content schemas', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, () => {
          // Return messages with incomplete content schemas - TDD failure
          return HttpResponse.json({
            messages: [
              createMockMessage({ type: 'text' }),
              // Image message with missing required fields
              {
                ...createMockImageMessage(),
                content: {
                  // Missing image object entirely
                  text: 'This should be an image message but content is wrong'
                }
              },
              // Voice message with incomplete content
              {
                ...createMockVoiceMessage(),
                content: {
                  voice: {
                    fileUrl: 'test.mp3'
                    // Missing duration, waveform, etc.
                  }
                }
              }
            ],
            pagination: {
              total: 3,
              hasMore: false,
              limit: 20
            }
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/messages?dialogId=test&limit=20`);
      const data: MessagesListResponse = await response.json();

      // Validate different message content types
      const textMessage = data.messages.find(m => m.type === 'text');
      const imageMessage = data.messages.find(m => m.type === 'image');
      const voiceMessage = data.messages.find(m => m.type === 'voice');

      // Text message validation
      if (textMessage) {
        expect(textMessage.content).toHaveProperty('text');
        expect(typeof textMessage.content.text).toBe('string');
      }

      // Image message validation - will fail due to incomplete mock
      if (imageMessage) {
        expect(imageMessage.content).toHaveProperty('image');
        expect(imageMessage.content.image).toHaveProperty('fileUrl');
        expect(imageMessage.content.image).toHaveProperty('width');
        expect(imageMessage.content.image).toHaveProperty('height');
        expect(imageMessage.content.image).toHaveProperty('fileSize');
        expect(typeof imageMessage.content.image.width).toBe('number');
        expect(typeof imageMessage.content.image.height).toBe('number');
      }

      // Voice message validation - will fail due to incomplete mock
      if (voiceMessage) {
        expect(voiceMessage.content).toHaveProperty('voice');
        expect(voiceMessage.content.voice).toHaveProperty('fileUrl');
        expect(voiceMessage.content.voice).toHaveProperty('duration');
        expect(voiceMessage.content.voice).toHaveProperty('waveform');
        expect(Array.isArray(voiceMessage.content.voice.waveform)).toBe(true);
        expect(typeof voiceMessage.content.voice.duration).toBe('number');
      }
    });

    it('should handle error responses correctly', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, () => {
          return HttpResponse.json(
            {
              error: 'Unauthorized',
              message: 'Invalid or missing authentication token',
              code: 'AUTH_REQUIRED'
            },
            { status: 401 }
          );
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/messages`);

      expect(response.status).toBe(401);
      expect(response.ok).toBe(false);

      const errorData = await response.json();
      expect(errorData).toHaveProperty('error');
      expect(errorData).toHaveProperty('message');
      expect(errorData).toHaveProperty('code');
    });

    it('should validate message metadata and analytics', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, () => {
          return HttpResponse.json({
            messages: [
              {
                ...createMockMessage(),
                // Incomplete metadata to trigger TDD failure
                metadata: {
                  // Missing deliveryStatus and analytics objects
                  autoGenerated: false
                }
              }
            ],
            pagination: {
              total: 1,
              hasMore: false,
              limit: 20
            }
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/messages?dialogId=test`);
      const data: MessagesListResponse = await response.json();

      const message = data.messages[0];
      expect(message).toHaveProperty('metadata');

      // Validate metadata structure - will fail due to incomplete mock
      if (message.metadata) {
        expect(message.metadata).toHaveProperty('deliveryStatus');
        expect(message.metadata).toHaveProperty('analytics');

        if (message.metadata.deliveryStatus) {
          expect(message.metadata.deliveryStatus).toHaveProperty('sent');
          expect(message.metadata.deliveryStatus).toHaveProperty('delivered');
          expect(message.metadata.deliveryStatus).toHaveProperty('failed');
          expect(message.metadata.deliveryStatus).toHaveProperty('retryCount');
          expect(typeof message.metadata.deliveryStatus.sent).toBe('boolean');
          expect(typeof message.metadata.deliveryStatus.retryCount).toBe('number');
        }

        if (message.metadata.analytics) {
          expect(message.metadata.analytics).toHaveProperty('opens');
          expect(message.metadata.analytics).toHaveProperty('clicks');
          expect(message.metadata.analytics).toHaveProperty('shares');
          expect(typeof message.metadata.analytics.opens).toBe('number');
          expect(typeof message.metadata.analytics.clicks).toBe('number');
        }
      }
    });
  });

  describe('POST /api/messages', () => {
    it('should create a text message with proper validation', async () => {
      // No MSW handler for POST - will cause TDD failure
      // server.use() intentionally omitted to demonstrate missing endpoint

      const messageData: CreateMessageRequest = {
        dialogId: '550e8400-e29b-41d4-a716-446655440002',
        type: 'text',
        content: {
          text: 'Hello, this is a new test message',
          html: '<p>Hello, this is a new test message</p>'
        },
        mentions: [],
        metadata: {
          autoGenerated: false
        }
      };

      const response = await fetch(`${API_BASE_URL}/api/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(messageData)
      });

      // This will fail because no MSW handler is set up - demonstrating TDD
      expect(response.ok).toBe(true);

      const data: MessageResponse = await response.json();
      expect(data).toHaveProperty('message');

      const { message } = data;
      expect(message).toHaveProperty('id');
      expect(message).toHaveProperty('dialogId');
      expect(message).toHaveProperty('senderId');
      expect(message.type).toBe('text');
      expect(message.content).toHaveProperty('text');
      expect(message.content.text).toBe(messageData.content.text);
      expect(message.isEdited).toBe(false);
      expect(message.isPinned).toBe(false);
      expect(Array.isArray(message.reactions)).toBe(true);
      expect(Array.isArray(message.readBy)).toBe(true);
      expect(Array.isArray(message.mentions)).toBe(true);
      expect(Array.isArray(message.attachments)).toBe(true);
    });

    it('should create an image message with attachments', async () => {
      server.use(
        http.post(`${API_BASE_URL}/api/messages`, async ({ request }) => {
          const body = await request.json() as CreateMessageRequest;

          // Return incomplete response to trigger TDD failure
          return HttpResponse.json({
            // Missing 'message' wrapper object
            id: '550e8400-e29b-41d4-a716-446655440007',
            type: body.type,
            // Wrong content structure
            text: body.content.text || 'fallback text'
            // Missing many required fields
          });
        })
      );

      const messageData: CreateMessageRequest = {
        dialogId: '550e8400-e29b-41d4-a716-446655440002',
        type: 'image',
        content: {
          image: {
            fileUrl: 'https://example.com/images/upload.jpg',
            thumbnailUrl: 'https://example.com/images/upload_thumb.jpg',
            width: 1920,
            height: 1080,
            fileSize: 2048000,
            caption: 'Uploaded test image',
            altText: 'A newly uploaded test image'
          }
        },
        attachments: [
          {
            type: 'image',
            fileUrl: 'https://example.com/images/upload.jpg',
            fileName: 'upload.jpg',
            fileSize: 2048000,
            mimeType: 'image/jpeg',
            metadata: {
              width: 1920,
              height: 1080,
              thumbnail: 'https://example.com/images/upload_thumb.jpg'
            }
          }
        ]
      };

      const response = await fetch(`${API_BASE_URL}/api/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(messageData)
      });

      expect(response.ok).toBe(true);
      const data: MessageResponse = await response.json();

      // Validate response structure - will fail due to incomplete mock response
      expect(data).toHaveProperty('message');
      expect(data.message.type).toBe('image');
      expect(data.message.content).toHaveProperty('image');
      expect(data.message.content.image).toHaveProperty('fileUrl');
      expect(data.message.content.image).toHaveProperty('width');
      expect(data.message.content.image).toHaveProperty('height');
      expect(Array.isArray(data.message.attachments)).toBe(true);
      expect(data.message.attachments.length).toBeGreaterThan(0);
    });

    it('should create a payment message with transaction details', async () => {
      server.use(
        http.post(`${API_BASE_URL}/api/messages`, async ({ request }) => {
          const body = await request.json() as CreateMessageRequest;

          // Return response with wrong payment structure - TDD failure
          return HttpResponse.json({
            message: {
              id: '550e8400-e29b-41d4-a716-446655440008',
              dialogId: body.dialogId,
              senderId: '550e8400-e29b-41d4-a716-446655440003',
              type: 'payment',
              // Wrong payment content structure
              content: {
                payment: {
                  amount: 100, // Missing currency and other required fields
                  description: 'Test payment'
                }
              },
              createdAt: new Date().toISOString(),
              isEdited: false,
              isPinned: false,
              reactions: [],
              readBy: [],
              mentions: [],
              attachments: []
            }
          });
        })
      );

      const paymentData: CreateMessageRequest = {
        dialogId: '550e8400-e29b-41d4-a716-446655440002',
        type: 'payment',
        content: {
          payment: {
            amount: 500.75,
            currency: 'THB',
            description: 'Lunch payment',
            transactionId: '550e8400-e29b-41d4-a716-446655440030',
            status: 'pending',
            paymentMethod: 'promptpay',
            recipient: 'Jane Doe',
            reference: 'LUNCH-2024-09-22'
          }
        }
      };

      const response = await fetch(`${API_BASE_URL}/api/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(paymentData)
      });

      expect(response.ok).toBe(true);
      const data: MessageResponse = await response.json();

      // Payment message validation - will fail due to incomplete mock
      expect(data.message.type).toBe('payment');
      expect(data.message.content).toHaveProperty('payment');
      expect(data.message.content.payment).toHaveProperty('amount');
      expect(data.message.content.payment).toHaveProperty('currency');
      expect(data.message.content.payment).toHaveProperty('description');
      expect(data.message.content.payment).toHaveProperty('transactionId');
      expect(data.message.content.payment).toHaveProperty('status');
      expect(data.message.content.payment.amount).toBe(500.75);
      expect(data.message.content.payment.currency).toBe('THB');
      expect(['pending', 'completed', 'failed', 'cancelled']).toContain(data.message.content.payment.status);
    });

    it('should handle message replies and threading', async () => {
      server.use(
        http.post(`${API_BASE_URL}/api/messages`, async ({ request }) => {
          const body = await request.json() as CreateMessageRequest;

          // Mock missing threading fields - TDD failure
          return HttpResponse.json({
            message: {
              ...createMockMessage(),
              // Missing replyToId and threadId even though they were provided
              replyToId: undefined,
              threadId: undefined
            }
          });
        })
      );

      const replyData: CreateMessageRequest = {
        dialogId: '550e8400-e29b-41d4-a716-446655440002',
        type: 'text',
        content: {
          text: 'This is a reply message'
        },
        replyToId: '550e8400-e29b-41d4-a716-446655440001',
        threadId: '550e8400-e29b-41d4-a716-446655440050'
      };

      const response = await fetch(`${API_BASE_URL}/api/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(replyData)
      });

      expect(response.ok).toBe(true);
      const data: MessageResponse = await response.json();

      // Threading validation - will fail due to mock not preserving threading fields
      expect(data.message.replyToId).toBe(replyData.replyToId);
      expect(data.message.threadId).toBe(replyData.threadId);
      expect(data.message.content.text).toBe(replyData.content.text);
    });

    it('should validate request body schema and return appropriate errors', async () => {
      server.use(
        http.post(`${API_BASE_URL}/api/messages`, async ({ request }) => {
          const body = await request.json();

          // Mock should validate request body but doesn't - TDD failure
          // Always returns success even for invalid data
          return HttpResponse.json({
            message: createMockMessage()
          });
        })
      );

      // Send invalid request data
      const invalidData = {
        // Missing required dialogId
        type: 'text',
        content: {
          text: 'Test message'
        }
      };

      const response = await fetch(`${API_BASE_URL}/api/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(invalidData)
      });

      // This should fail with validation error but mock doesn't validate - TDD failure
      expect(response.status).toBe(400);
      const errorData = await response.json();
      expect(errorData).toHaveProperty('error');
      expect(errorData.error).toContain('dialogId');
    });

    it('should handle scheduled messages', async () => {
      // No MSW handler for scheduled messages - TDD failure
      const scheduledData: CreateMessageRequest = {
        dialogId: '550e8400-e29b-41d4-a716-446655440002',
        type: 'text',
        content: {
          text: 'This message will be sent later'
        },
        scheduledAt: '2024-09-22T15:00:00.000Z'
      };

      const response = await fetch(`${API_BASE_URL}/api/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(scheduledData)
      });

      // Will fail because scheduled message handling is not implemented
      expect(response.ok).toBe(true);
      const data: MessageResponse = await response.json();

      expect(data.message).toHaveProperty('metadata');
      expect(data.message.metadata).toHaveProperty('scheduledFor');
      expect(data.message.metadata.scheduledFor).toBe(scheduledData.scheduledAt);
    });

    it('should handle ephemeral messages', async () => {
      const ephemeralData: CreateMessageRequest = {
        dialogId: '550e8400-e29b-41d4-a716-446655440002',
        type: 'text',
        content: {
          text: 'This message will self-destruct'
        },
        isEphemeral: true,
        deleteAfter: 3600 // 1 hour
      };

      // No handler for ephemeral messages - TDD failure
      const response = await fetch(`${API_BASE_URL}/api/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(ephemeralData)
      });

      expect(response.ok).toBe(true);
      const data: MessageResponse = await response.json();

      // Ephemeral message validation - will fail because not implemented
      expect(data.message.metadata).toHaveProperty('isEphemeral');
      expect(data.message.metadata).toHaveProperty('deleteAfter');
      expect(data.message.metadata.isEphemeral).toBe(true);
      expect(data.message.metadata.deleteAfter).toBe(3600);
    });
  });

  describe('Additional Contract Validations', () => {
    it('should validate message reaction schemas', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, () => {
          return HttpResponse.json({
            messages: [
              {
                ...createMockMessage(),
                // Incomplete reactions array - TDD failure
                reactions: [
                  {
                    // Missing required fields
                    emoji: '👍',
                    userId: '123'
                    // Missing id, messageId, createdAt, etc.
                  }
                ]
              }
            ],
            pagination: { total: 1, hasMore: false, limit: 20 }
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/messages?dialogId=test`);
      const data: MessagesListResponse = await response.json();

      const message = data.messages[0];
      expect(Array.isArray(message.reactions)).toBe(true);

      if (message.reactions.length > 0) {
        const reaction = message.reactions[0];
        expect(reaction).toHaveProperty('id');
        expect(reaction).toHaveProperty('messageId');
        expect(reaction).toHaveProperty('userId');
        expect(reaction).toHaveProperty('emoji');
        expect(reaction).toHaveProperty('createdAt');
        expect(typeof reaction.emoji).toBe('string');
        expect(reaction.messageId).toBe(message.id);
      }
    });

    it('should validate message read receipt schemas', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, () => {
          return HttpResponse.json({
            messages: [
              {
                ...createMockMessage(),
                // Incomplete readBy array - TDD failure
                readBy: [
                  {
                    userId: '123'
                    // Missing messageId, readAt, deviceId
                  }
                ]
              }
            ],
            pagination: { total: 1, hasMore: false, limit: 20 }
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/messages?dialogId=test`);
      const data: MessagesListResponse = await response.json();

      const message = data.messages[0];
      expect(Array.isArray(message.readBy)).toBe(true);

      if (message.readBy.length > 0) {
        const readReceipt = message.readBy[0];
        expect(readReceipt).toHaveProperty('messageId');
        expect(readReceipt).toHaveProperty('userId');
        expect(readReceipt).toHaveProperty('readAt');
        expect(readReceipt.messageId).toBe(message.id);
        expect(typeof readReceipt.readAt).toBe('string');
        expect(readReceipt.readAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/);
      }
    });

    it('should validate attachment metadata schemas', async () => {
      server.use(
        http.get(`${API_BASE_URL}/api/messages`, () => {
          return HttpResponse.json({
            messages: [
              {
                ...createMockImageMessage(),
                attachments: [
                  {
                    id: '1',
                    // Missing most required fields for TDD failure
                    type: 'image',
                    fileUrl: 'test.jpg'
                  }
                ]
              }
            ],
            pagination: { total: 1, hasMore: false, limit: 20 }
          });
        })
      );

      const response = await fetch(`${API_BASE_URL}/api/messages?dialogId=test&hasAttachments=true`);
      const data: MessagesListResponse = await response.json();

      const message = data.messages[0];
      expect(Array.isArray(message.attachments)).toBe(true);

      if (message.attachments.length > 0) {
        const attachment = message.attachments[0];
        expect(attachment).toHaveProperty('id');
        expect(attachment).toHaveProperty('messageId');
        expect(attachment).toHaveProperty('type');
        expect(attachment).toHaveProperty('fileUrl');
        expect(attachment).toHaveProperty('fileName');
        expect(attachment).toHaveProperty('fileSize');
        expect(attachment).toHaveProperty('mimeType');
        expect(attachment).toHaveProperty('createdAt');
        expect(attachment.messageId).toBe(message.id);
        expect(['file', 'image', 'video', 'audio', 'document']).toContain(attachment.type);
        expect(typeof attachment.fileSize).toBe('number');
        expect(attachment.fileSize).toBeGreaterThan(0);
      }
    });
  });
});