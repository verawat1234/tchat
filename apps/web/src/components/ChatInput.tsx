import React, { useRef, useState } from 'react';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { Card } from './ui/card';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from './ui/dropdown-menu';
import { 
  Send, Paperclip, Mic, Smile, File, Image, Camera, MapPin, CreditCard, Gift,
  X, Pause, Play
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

interface ChatInputProps {
  messageInput: string;
  setMessageInput: (value: string) => void;
  onSendMessage: () => void;
  replyToMessage: Message | null;
  setReplyToMessage: (message: Message | null) => void;
  editingMessage: Message | null;
  setEditingMessage: (message: Message | null) => void;
  isRecording: boolean;
  recordingDuration: number;
  onStartRecording: () => void;
  onStopRecording: () => void;
  onCancelRecording: () => void;
}

export function ChatInput({
  messageInput,
  setMessageInput,
  onSendMessage,
  replyToMessage,
  setReplyToMessage,
  editingMessage,
  setEditingMessage,
  isRecording,
  recordingDuration,
  onStartRecording,
  onStopRecording,
  onCancelRecording
}: ChatInputProps) {
  const [showAttachmentMenu, setShowAttachmentMenu] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const imageInputRef = useRef<HTMLInputElement>(null);

  const formatRecordingTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const handleFileUpload = () => {
    fileInputRef.current?.click();
    setShowAttachmentMenu(false);
  };

  const handleImageUpload = () => {
    imageInputRef.current?.click();
    setShowAttachmentMenu(false);
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      toast.success(`Uploading ${file.name}...`);
    }
  };

  const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      toast.success(`Uploading image...`);
    }
  };

  const handleLocationShare = () => {
    toast.success('Sharing current location...');
    setShowAttachmentMenu(false);
  };

  const handlePaymentSend = () => {
    toast.success('Payment dialog would open');
    setShowAttachmentMenu(false);
  };

  const handleGiftSend = () => {
    toast.success('Gift selection would open');
    setShowAttachmentMenu(false);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      onSendMessage();
    }
  };

  const cancelReply = () => {
    setReplyToMessage(null);
  };

  const cancelEdit = () => {
    setEditingMessage(null);
    setMessageInput('');
  };

  return (
    <div className="border-t border-border bg-card">
      {/* Reply/Edit Preview */}
      {(replyToMessage || editingMessage) && (
        <div className="px-4 py-2 border-b border-border bg-muted/50">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <p className="text-xs text-muted-foreground mb-1">
                {editingMessage ? 'Editing message' : `Replying to ${replyToMessage?.senderName}`}
              </p>
              <p className="text-sm truncate">
                {editingMessage?.content || replyToMessage?.content}
              </p>
            </div>
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={editingMessage ? cancelEdit : cancelReply}
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
        </div>
      )}

      {/* Voice Recording UI */}
      {isRecording && (
        <div className="px-4 py-3 bg-red-50 border-b border-border">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-3 h-3 bg-red-500 rounded-full animate-pulse"></div>
              <span className="text-sm font-medium text-red-700">
                Recording... {formatRecordingTime(recordingDuration)}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={onCancelRecording}>
                Cancel
              </Button>
              <Button size="sm" onClick={onStopRecording}>
                <Pause className="w-4 h-4 mr-1" />
                Send
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Main Input Area */}
      <div className="p-4">
        <div className="flex items-end gap-2">
          {/* Attachment Menu */}
          <DropdownMenu open={showAttachmentMenu} onOpenChange={setShowAttachmentMenu}>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="flex-shrink-0">
                <Paperclip className="w-5 h-5" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-64">
              <DropdownMenuLabel>Send Attachment</DropdownMenuLabel>
              <DropdownMenuSeparator />
              
              <DropdownMenuItem onClick={handleImageUpload}>
                <Image className="w-4 h-4 mr-2" />
                Photo & Video
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={handleFileUpload}>
                <File className="w-4 h-4 mr-2" />
                Document
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={handleLocationShare}>
                <MapPin className="w-4 h-4 mr-2" />
                Location
              </DropdownMenuItem>
              
              <DropdownMenuSeparator />
              
              <DropdownMenuItem onClick={handlePaymentSend}>
                <CreditCard className="w-4 h-4 mr-2" />
                Send Payment
              </DropdownMenuItem>
              
              <DropdownMenuItem onClick={handleGiftSend}>
                <Gift className="w-4 h-4 mr-2" />
                Send Gift
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* Message Input */}
          <div className="flex-1 relative">
            <Input
              placeholder={
                editingMessage ? "Edit message..." : 
                replyToMessage ? "Reply..." : 
                "Type a message..."
              }
              value={messageInput}
              onChange={(e) => setMessageInput(e.target.value)}
              onKeyPress={handleKeyPress}
              className="pr-10"
              disabled={isRecording}
            />
            
            {/* Emoji Button */}
            <Button 
              variant="ghost" 
              size="icon" 
              className="absolute right-1 top-1/2 -translate-y-1/2 h-8 w-8"
              onClick={() => toast.success('Emoji picker would open')}
            >
              <Smile className="w-4 h-4" />
            </Button>
          </div>

          {/* Voice/Send Button */}
          {messageInput.trim() ? (
            <Button onClick={onSendMessage} size="icon" className="flex-shrink-0">
              <Send className="w-5 h-5" />
            </Button>
          ) : (
            <Button 
              variant={isRecording ? "destructive" : "ghost"}
              size="icon" 
              className="flex-shrink-0"
              onClick={isRecording ? onStopRecording : onStartRecording}
              onMouseDown={!isRecording ? onStartRecording : undefined}
            >
              <Mic className="w-5 h-5" />
            </Button>
          )}
        </div>

        {/* Quick Actions Row */}
        <div className="flex items-center gap-2 mt-2">
          <Badge variant="outline" className="text-xs cursor-pointer hover:bg-accent" onClick={() => setMessageInput(messageInput + "👍")}>
            👍
          </Badge>
          <Badge variant="outline" className="text-xs cursor-pointer hover:bg-accent" onClick={() => setMessageInput(messageInput + "❤️")}>
            ❤️
          </Badge>
          <Badge variant="outline" className="text-xs cursor-pointer hover:bg-accent" onClick={() => setMessageInput(messageInput + "😊")}>
            😊
          </Badge>
          <Badge variant="outline" className="text-xs cursor-pointer hover:bg-accent" onClick={() => setMessageInput(messageInput + "🙏")}>
            🙏
          </Badge>
        </div>
      </div>

      {/* Hidden File Inputs */}
      <input
        ref={fileInputRef}
        type="file"
        hidden
        onChange={handleFileChange}
        accept=".pdf,.doc,.docx,.txt,.xls,.xlsx,.ppt,.pptx"
      />
      <input
        ref={imageInputRef}
        type="file"
        hidden
        onChange={handleImageChange}
        accept="image/*,video/*"
        multiple
      />
    </div>
  );
}