/**
 * ChatMessage Component Contract Tests
 * CRITICAL: These tests MUST FAIL until ChatMessage component is implemented
 */

import { render, screen, fireEvent } from '@testing-library/react';
import { ChatMessage, TypingIndicator, MessageGroup, MessageBubble } from './chat-message';
import type {
  ChatMessageProps,
  TypingIndicatorProps,
  MessageGroupProps,
  MessageBubbleProps,
  User,
  MediaContent,
  SystemContent,
  Reaction,
  MessageStatus
} from '../../../../specs/001-agent-frontend-specialist/contracts/chat-message';

// Mock data for tests
const mockUser: User = {
  id: '1',
  name: 'John Doe',
  avatar: '/avatar.jpg',
  status: 'online'
};

const mockOtherUser: User = {
  id: '2',
  name: 'Jane Smith',
  avatar: '/avatar2.jpg',
  status: 'offline'
};

describe('ChatMessage Contract Tests', () => {
  const baseProps: ChatMessageProps = {
    type: 'text',
    content: 'Hello, this is a test message',
    timestamp: new Date('2023-01-01T12:00:00Z'),
    sender: mockUser,
    testId: 'chat-message-test'
  };

  describe('Basic Rendering', () => {
    test('renders text message correctly', () => {
      render(<ChatMessage {...baseProps} showSender={true} />);

      expect(screen.getByTestId('chat-message-test')).toBeInTheDocument();
      expect(screen.getByText('Hello, this is a test message')).toBeInTheDocument();
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    test('applies custom className', () => {
      render(<ChatMessage {...baseProps} className="custom-message" />);

      const message = screen.getByTestId('chat-message-test');
      expect(message).toHaveClass('custom-message');
    });
  });

  describe('Message Types', () => {
    test('renders image message correctly', () => {
      const imageContent: MediaContent = {
        url: '/test-image.jpg',
        type: 'image',
        filename: 'test.jpg',
        size: 1024,
        dimensions: { width: 800, height: 600 }
      };

      render(
        <ChatMessage
          {...baseProps}
          type="image"
          content={imageContent}
          testId="image-message"
        />
      );

      expect(screen.getByTestId('image-message')).toBeInTheDocument();
      const images = screen.getAllByRole('img');
      const messageImage = images.find(img => img.getAttribute('src') === '/test-image.jpg');
      expect(messageImage).toBeInTheDocument();
    });

    test('renders video message correctly', () => {
      const videoContent: MediaContent = {
        url: '/test-video.mp4',
        type: 'video',
        filename: 'test.mp4',
        size: 5120,
        duration: 120
      };

      render(
        <ChatMessage
          {...baseProps}
          type="video"
          content={videoContent}
          testId="video-message"
        />
      );

      expect(screen.getByTestId('video-message')).toBeInTheDocument();
      const video = screen.getByTestId('video-message').querySelector('video');
      expect(video).toHaveAttribute('src', '/test-video.mp4');
    });

    test('renders audio message correctly', () => {
      const audioContent: MediaContent = {
        url: '/test-audio.mp3',
        type: 'audio',
        filename: 'test.mp3',
        duration: 60
      };

      render(
        <ChatMessage
          {...baseProps}
          type="audio"
          content={audioContent}
          testId="audio-message"
        />
      );

      expect(screen.getByTestId('audio-message')).toBeInTheDocument();
      const audio = screen.getByTestId('audio-message').querySelector('audio');
      expect(audio).toHaveAttribute('src', '/test-audio.mp3');
    });

    test('renders file message correctly', () => {
      const fileContent: MediaContent = {
        url: '/test-file.pdf',
        type: 'file',
        filename: 'document.pdf',
        size: 2048
      };

      render(
        <ChatMessage
          {...baseProps}
          type="file"
          content={fileContent}
          testId="file-message"
        />
      );

      expect(screen.getByTestId('file-message')).toBeInTheDocument();
      expect(screen.getByText('document.pdf')).toBeInTheDocument();
    });

    test('renders system message correctly', () => {
      const systemContent: SystemContent = {
        type: 'join',
        message: 'User joined the chat',
        metadata: { userId: '123' }
      };

      render(
        <ChatMessage
          {...baseProps}
          type="system"
          content={systemContent}
          testId="system-message"
        />
      );

      expect(screen.getByTestId('system-message')).toBeInTheDocument();
      expect(screen.getByText('User joined the chat')).toBeInTheDocument();
      expect(screen.getByTestId('system-message')).toHaveClass('chat-message-system');
    });
  });

  describe('Message States', () => {
    test('applies own message styling when isOwn is true', () => {
      render(<ChatMessage {...baseProps} isOwn testId="own-message" />);

      const message = screen.getByTestId('own-message');
      expect(message).toHaveClass('chat-message-own');
    });

    test('displays message status correctly', () => {
      const statuses: MessageStatus[] = ['sending', 'sent', 'delivered', 'read', 'failed'];

      statuses.forEach(status => {
        const { unmount } = render(
          <ChatMessage {...baseProps} status={status} testId={`message-${status}`} />
        );

        const message = screen.getByTestId(`message-${status}`);
        expect(message).toHaveAttribute('data-status', status);
        unmount();
      });
    });

    test('shows editing state correctly', () => {
      render(<ChatMessage {...baseProps} editing testId="editing-message" />);

      const message = screen.getByTestId('editing-message');
      expect(message).toHaveClass('chat-message-editing');
    });

    test('shows selected state correctly', () => {
      render(<ChatMessage {...baseProps} selected testId="selected-message" />);

      const message = screen.getByTestId('selected-message');
      expect(message).toHaveClass('chat-message-selected');
    });
  });

  describe('Message Features', () => {
    test('displays reactions correctly', () => {
      const reactions: Reaction[] = [
        { emoji: '👍', count: 3, users: ['1', '2', '3'], hasUserReacted: true },
        { emoji: '❤️', count: 1, users: ['2'], hasUserReacted: false }
      ];

      render(<ChatMessage {...baseProps} reactions={reactions} testId="message-reactions" />);

      expect(screen.getByText('👍')).toBeInTheDocument();
      expect(screen.getByText('3')).toBeInTheDocument();
      expect(screen.getByText('❤️')).toBeInTheDocument();
      expect(screen.getByText('1')).toBeInTheDocument();
    });

    test('displays reply message correctly', () => {
      const replyMessage: Partial<ChatMessageProps> = {
        content: 'Original message',
        sender: mockOtherUser
      };

      render(
        <ChatMessage {...baseProps} reply={replyMessage} testId="message-with-reply" />
      );

      expect(screen.getByText('Original message')).toBeInTheDocument();
      // Reply only shows content, not sender name
    });

    test('shows/hides timestamp based on showTimestamp prop', () => {
      const { rerender } = render(
        <ChatMessage {...baseProps} showTimestamp testId="message-timestamp" />
      );

      expect(screen.getByTestId('message-timestamp')).toContainElement(
        screen.getByText(/19:00/)  // UTC time shows as 19:00
      );

      rerender(
        <ChatMessage {...baseProps} showTimestamp={false} testId="message-timestamp" />
      );

      expect(screen.queryByText(/12:00/)).not.toBeInTheDocument();
    });

    test('shows/hides avatar based on showAvatar prop', () => {
      const { rerender } = render(
        <ChatMessage {...baseProps} showAvatar testId="message-avatar" />
      );

      expect(screen.getByRole('img', { name: /John Doe/i })).toBeInTheDocument();

      rerender(
        <ChatMessage {...baseProps} showAvatar={false} testId="message-avatar" />
      );

      expect(screen.queryByRole('img', { name: /John Doe/i })).not.toBeInTheDocument();
    });

    test('shows/hides sender name based on showSender prop', () => {
      const { rerender } = render(
        <ChatMessage {...baseProps} showSender testId="message-sender" />
      );

      expect(screen.getByText('John Doe')).toBeInTheDocument();

      rerender(
        <ChatMessage {...baseProps} showSender={false} testId="message-sender" />
      );

      expect(screen.queryByText('John Doe')).not.toBeInTheDocument();
    });
  });

  describe('Event Handlers', () => {
    test('calls onClick when message is clicked', () => {
      const onClick = vi.fn();
      render(<ChatMessage {...baseProps} onClick={onClick} testId="clickable-message" />);

      fireEvent.click(screen.getByTestId('clickable-message'));
      expect(onClick).toHaveBeenCalledTimes(1);
    });

    test('calls onReaction when reaction is clicked', () => {
      const onReaction = vi.fn();
      const reactions: Reaction[] = [
        { emoji: '👍', count: 1, users: ['1'], hasUserReacted: false }
      ];

      render(
        <ChatMessage
          {...baseProps}
          reactions={reactions}
          onReaction={onReaction}
          testId="reaction-message"
        />
      );

      fireEvent.click(screen.getByText('👍'));
      expect(onReaction).toHaveBeenCalledWith('👍');
    });

    test('calls onReply when reply button is clicked', () => {
      const onReply = vi.fn();
      render(<ChatMessage {...baseProps} onReply={onReply} testId="reply-message" />);

      const replyButton = screen.getByRole('button', { name: /reply/i });
      fireEvent.click(replyButton);
      expect(onReply).toHaveBeenCalledTimes(1);
    });

    test('calls onEdit when edit button is clicked', () => {
      const onEdit = vi.fn();
      render(<ChatMessage {...baseProps} onEdit={onEdit} isOwn testId="edit-message" />);

      const editButton = screen.getByRole('button', { name: /edit/i });
      fireEvent.click(editButton);
      expect(onEdit).toHaveBeenCalledTimes(1);
    });

    test('calls onDelete when delete button is clicked', () => {
      const onDelete = vi.fn();
      render(<ChatMessage {...baseProps} onDelete={onDelete} isOwn testId="delete-message" />);

      const deleteButton = screen.getByRole('button', { name: /delete/i });
      fireEvent.click(deleteButton);
      expect(onDelete).toHaveBeenCalledTimes(1);
    });
  });

  describe('Accessibility', () => {
    test('supports ARIA label', () => {
      render(
        <ChatMessage
          {...baseProps}
          aria-label="Message from John Doe"
          testId="aria-message"
        />
      );

      const message = screen.getByTestId('aria-message');
      expect(message).toHaveAttribute('aria-label', 'Message from John Doe');
    });

    test('supports custom role', () => {
      render(<ChatMessage {...baseProps} role="article" testId="role-message" />);

      const message = screen.getByTestId('role-message');
      expect(message).toHaveAttribute('role', 'article');
    });

    test('supports keyboard navigation', () => {
      render(<ChatMessage {...baseProps} tabIndex={0} testId="keyboard-message" />);

      const message = screen.getByTestId('keyboard-message');
      expect(message).toHaveAttribute('tabIndex', '0');
    });
  });
});

describe('TypingIndicator Contract Tests', () => {
  const typingUsers: User[] = [
    { id: '1', name: 'John Doe', status: 'online' },
    { id: '2', name: 'Jane Smith', status: 'online' }
  ];

  test('renders typing indicator with single user', () => {
    render(
      <TypingIndicator users={[typingUsers[0]]} testId="typing-single" />
    );

    expect(screen.getByTestId('typing-single')).toBeInTheDocument();
    expect(screen.getByText(/John Doe.*typing/i)).toBeInTheDocument();
  });

  test('renders typing indicator with multiple users', () => {
    render(
      <TypingIndicator users={typingUsers} testId="typing-multiple" />
    );

    expect(screen.getByTestId('typing-multiple')).toBeInTheDocument();
    expect(screen.getByText(/John Doe.*Jane Smith.*typing/i)).toBeInTheDocument();
  });

  test('limits displayed users based on maxUsers', () => {
    const manyUsers = [
      ...typingUsers,
      { id: '3', name: 'Bob Wilson', status: 'online' as const },
      { id: '4', name: 'Alice Brown', status: 'online' as const }
    ];

    render(
      <TypingIndicator users={manyUsers} maxUsers={2} testId="typing-limited" />
    );

    const indicator = screen.getByTestId('typing-limited');
    // With 4 users and maxUsers=2, shows "John Doe, Jane Smith and 2 others are typing"
    expect(indicator).toContainElement(screen.getByText(/John Doe, Jane Smith and 2 others are typing/i));
  });
});

describe('MessageGroup Contract Tests', () => {
  const groupMessages: ChatMessageProps[] = [
    {
      type: 'text',
      content: 'First message',
      timestamp: new Date('2023-01-01T12:00:00Z'),
      sender: mockUser
    },
    {
      type: 'text',
      content: 'Second message',
      timestamp: new Date('2023-01-01T12:01:00Z'),
      sender: mockUser
    }
  ];

  test('renders message group correctly', () => {
    render(
      <MessageGroup
        messages={groupMessages}
        sender={mockUser}
        testId="message-group-test"
      />
    );

    expect(screen.getByTestId('message-group-test')).toBeInTheDocument();
    expect(screen.getByText('First message')).toBeInTheDocument();
    expect(screen.getByText('Second message')).toBeInTheDocument();
  });

  test('applies own group styling when isOwn is true', () => {
    render(
      <MessageGroup
        messages={groupMessages}
        sender={mockUser}
        isOwn
        testId="own-group"
      />
    );

    const group = screen.getByTestId('own-group');
    expect(group).toHaveClass('message-group-own');
  });

  test('shows avatar when showAvatar is true', () => {
    render(
      <MessageGroup
        messages={groupMessages}
        sender={mockUser}
        showAvatar
        testId="group-avatar"
      />
    );

    expect(screen.getByRole('img', { name: /John Doe/i })).toBeInTheDocument();
  });

  test('shows timestamp when showTimestamp is true', () => {
    render(
      <MessageGroup
        messages={groupMessages}
        sender={mockUser}
        showTimestamp
        testId="group-timestamp"
      />
    );

    expect(screen.getByText(/19:01/)).toBeInTheDocument();  // UTC time shows as 19:01 for second message
  });
});

describe('MessageBubble Contract Tests', () => {
  test('renders bubble with default variant', () => {
    render(
      <MessageBubble testId="bubble-default">
        Bubble content
      </MessageBubble>
    );

    const bubble = screen.getByTestId('bubble-default');
    expect(bubble).toHaveClass('message-bubble-default');
    expect(screen.getByText('Bubble content')).toBeInTheDocument();
  });

  test('applies variant classes correctly', () => {
    const variants: MessageBubbleProps['variant'][] = ['default', 'own', 'system'];

    variants.forEach(variant => {
      const { unmount } = render(
        <MessageBubble variant={variant} testId={`bubble-${variant}`}>
          Content
        </MessageBubble>
      );

      const bubble = screen.getByTestId(`bubble-${variant}`);
      expect(bubble).toHaveClass(`message-bubble-${variant}`);
      unmount();
    });
  });

  test('applies position classes when grouped', () => {
    const positions: MessageBubbleProps['position'][] = ['first', 'middle', 'last', 'single'];

    positions.forEach(position => {
      const { unmount } = render(
        <MessageBubble grouped position={position} testId={`bubble-${position}`}>
          Content
        </MessageBubble>
      );

      const bubble = screen.getByTestId(`bubble-${position}`);
      expect(bubble).toHaveClass(`message-bubble-${position}`);
      unmount();
    });
  });

  test('shows tail when tail prop is true', () => {
    render(
      <MessageBubble tail testId="bubble-tail">
        Content with tail
      </MessageBubble>
    );

    const bubble = screen.getByTestId('bubble-tail');
    expect(bubble).toHaveClass('message-bubble-tail');
  });
});