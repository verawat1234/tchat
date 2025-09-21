import React from 'react';
import { Button } from './ui/button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from './ui/dropdown-menu';
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuLabel, ContextMenuSeparator, ContextMenuTrigger } from './ui/context-menu';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Badge } from './ui/badge';
import { 
  Reply, Forward, Copy, Edit3, Trash2, Pin, Archive, VolumeX, Volume2, UserX, 
  RotateCcw, Smile, MapPin, Gift, CreditCard, Bookmark, MoreVertical, 
  File, Download, Play, Bot, Phone, Video, Star, Lock
} from 'lucide-react';
import { toast } from "sonner";

interface Message {
  id: string;
  senderId: string;
  senderName: string;
  content: string;
  timestamp: string;
  type: 'text' | 'payment' | 'system' | 'voice' | 'file' | 'image';
  isOwn: boolean;
  fileUrl?: string;
  fileName?: string;
  fileSize?: string;
  duration?: string;
}

interface ChatActionsProps {
  messages: Message[];
  selectedMessages: string[];
  isSelectionMode: boolean;
  onReplyToMessage: (message: Message) => void;
  onEditMessage: (message: Message) => void;
  onDeleteMessage: (messageId: string) => void;
  onForwardMessage: (messageId: string) => void;
  onCopyMessage: (content: string) => void;
  onPinMessage: (messageId: string) => void;
  onSelectMessage: (messageId: string) => void;
}

export function ChatMessages({ 
  messages, 
  selectedMessages, 
  isSelectionMode, 
  onReplyToMessage,
  onEditMessage,
  onDeleteMessage,
  onForwardMessage,
  onCopyMessage,
  onPinMessage,
  onSelectMessage
}: ChatActionsProps) {

  const renderMessage = (message: Message) => (
    <ContextMenu key={message.id}>
      <ContextMenuTrigger>
        <div
          className={`flex gap-3 ${message.isOwn ? 'flex-row-reverse' : ''} ${
            isSelectionMode ? 'cursor-pointer' : ''
          } ${
            selectedMessages.includes(message.id) ? 'bg-accent/50 rounded-lg p-2 -m-2' : ''
          }`}
          onClick={() => onSelectMessage(message.id)}
        >
          {!message.isOwn && (
            <Avatar className="w-8 h-8 flex-shrink-0">
              <AvatarFallback className="text-xs">
                {message.type === 'system' ? <Bot className="w-4 h-4" /> : message.senderName.charAt(0)}
              </AvatarFallback>
            </Avatar>
          )}
          
          <div className={`flex flex-col max-w-md ${message.isOwn ? 'items-end' : 'items-start'}`}>
            {!message.isOwn && (
              <span className="text-xs text-muted-foreground mb-1 flex items-center gap-1">
                {message.type === 'system' && <Bot className="w-3 h-3" />}
                {message.senderName}
              </span>
            )}
            
            <div
              className={`rounded-lg px-3 py-2 ${
                message.type === 'system'
                  ? 'bg-chart-1/10 border border-chart-1/20'
                  : message.isOwn
                  ? 'bg-primary text-primary-foreground'
                  : 'bg-muted'
              }`}
            >
              {message.type === 'voice' ? (
                <div className="flex items-center gap-2">
                  <button className="p-1 rounded-full bg-background/20">
                    <Play className="w-3 h-3" />
                  </button>
                  <div className="flex-1 h-1 bg-background/20 rounded-full">
                    <div className="w-1/3 h-full bg-background/40 rounded-full"></div>
                  </div>
                  <span className="text-xs opacity-70">{message.duration}</span>
                </div>
              ) : message.type === 'file' ? (
                <div className="flex items-center gap-2">
                  <File className="w-4 h-4" />
                  <div>
                    <p className="text-sm font-medium">{message.fileName}</p>
                    <p className="text-xs opacity-70">{message.fileSize}</p>
                  </div>
                  <Button size="sm" variant="ghost" className="h-6 w-6 p-0">
                    <Download className="w-3 h-3" />
                  </Button>
                </div>
              ) : (
                <p className="text-sm whitespace-pre-line">{message.content}</p>
              )}
            </div>
            
            <span className="text-xs text-muted-foreground mt-1">
              {message.timestamp}
            </span>
          </div>
        </div>
      </ContextMenuTrigger>
      
      <ContextMenuContent>
        <ContextMenuItem onClick={() => onReplyToMessage(message)}>
          <Reply className="w-4 h-4 mr-2" />
          Reply
        </ContextMenuItem>
        <ContextMenuItem onClick={() => onCopyMessage(message.content)}>
          <Copy className="w-4 h-4 mr-2" />
          Copy
        </ContextMenuItem>
        <ContextMenuItem onClick={() => onForwardMessage(message.id)}>
          <Forward className="w-4 h-4 mr-2" />
          Forward
        </ContextMenuItem>
        {message.isOwn && (
          <ContextMenuItem onClick={() => onEditMessage(message)}>
            <Edit3 className="w-4 h-4 mr-2" />
            Edit
          </ContextMenuItem>
        )}
        <ContextMenuItem onClick={() => onPinMessage(message.id)}>
          <Pin className="w-4 h-4 mr-2" />
          Pin
        </ContextMenuItem>
        <ContextMenuSeparator />
        <ContextMenuItem 
          onClick={() => onDeleteMessage(message.id)}
          className="text-destructive"
        >
          <Trash2 className="w-4 h-4 mr-2" />
          Delete
        </ContextMenuItem>
      </ContextMenuContent>
    </ContextMenu>
  );

  return (
    <div className="space-y-4">
      {messages.map(renderMessage)}
    </div>
  );
}

interface ChatHeaderActionsProps {
  selectedDialog: any;
  chatActions: {
    isMuted: boolean;
    isPinned: boolean;
    isArchived: boolean;
    isBlocked: boolean;
  };
  onVideoCall: () => void;
  onVoiceCall: () => void;
  onMuteChat: () => void;
  onPinChat: () => void;
  onArchiveChat: () => void;
  onBlockContact: () => void;
  onClearHistory: () => void;
  toggleSelectionMode: () => void;
  isSelectionMode: boolean;
}

export function ChatHeaderActions({
  selectedDialog,
  chatActions,
  onVideoCall,
  onVoiceCall,
  onMuteChat,
  onPinChat,
  onArchiveChat,
  onBlockContact,
  onClearHistory,
  toggleSelectionMode,
  isSelectionMode
}: ChatHeaderActionsProps) {
  return (
    <div className="flex items-center gap-2">
      <Button variant="ghost" size="icon" onClick={onVoiceCall}>
        <Phone className="w-5 h-5" />
      </Button>
      <Button variant="ghost" size="icon" onClick={onVideoCall}>
        <Video className="w-5 h-5" />
      </Button>
      <Button 
        variant={isSelectionMode ? "default" : "ghost"} 
        size="icon" 
        onClick={toggleSelectionMode}
      >
        <Bookmark className="w-5 h-5" />
      </Button>
      
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon">
            <MoreVertical className="w-5 h-5" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>Chat Actions</DropdownMenuLabel>
          <DropdownMenuSeparator />
          
          <DropdownMenuItem onClick={onPinChat}>
            <Pin className="w-4 h-4 mr-2" />
            {chatActions.isPinned ? 'Unpin Chat' : 'Pin Chat'}
          </DropdownMenuItem>
          
          <DropdownMenuItem onClick={onMuteChat}>
            {chatActions.isMuted ? (
              <>
                <Volume2 className="w-4 h-4 mr-2" />
                Unmute Chat
              </>
            ) : (
              <>
                <VolumeX className="w-4 h-4 mr-2" />
                Mute Chat
              </>
            )}
          </DropdownMenuItem>
          
          <DropdownMenuItem onClick={onArchiveChat}>
            <Archive className="w-4 h-4 mr-2" />
            Archive Chat
          </DropdownMenuItem>
          
          <DropdownMenuSeparator />
          
          <DropdownMenuItem onClick={onClearHistory}>
            <RotateCcw className="w-4 h-4 mr-2" />
            Clear History
          </DropdownMenuItem>
          
          <DropdownMenuItem 
            onClick={onBlockContact}
            className="text-destructive"
          >
            <UserX className="w-4 h-4 mr-2" />
            {chatActions.isBlocked ? 'Unblock Contact' : 'Block Contact'}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}