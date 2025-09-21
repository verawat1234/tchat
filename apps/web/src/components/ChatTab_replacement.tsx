  return (
    <div className="h-full flex">
      <Tabs defaultValue="chat" className="flex flex-col w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="chat" className="flex items-center gap-2">
            <MessageCircle className="w-4 h-4" />
            Chat
          </TabsTrigger>
          <TabsTrigger value="work" className="flex items-center gap-2">
            <Briefcase className="w-4 h-4" />
            Work
          </TabsTrigger>
        </TabsList>

        {/* Personal Chat Tab */}
        <TabsContent value="chat" className="flex-1 flex">
          {/* Chat List */}
          <div className="w-80 border-r border-border flex flex-col">
            <div className="p-4 border-b border-border">
              <div className="flex items-center gap-2 mb-3">
                <div className="flex-1">
                  <Input
                    placeholder="Search chats..."
                    className="w-full"
                    startIcon={<Search className="w-4 h-4" />}
                  />
                </div>
                <Button size="icon" variant="ghost" onClick={onNewChat}>
                  <Plus className="w-4 h-4" />
                </Button>
              </div>

              {/* Tag Filter */}
              <div className="flex items-center gap-2">
                <Button
                  variant={showTagFilter ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setShowTagFilter(!showTagFilter)}
                >
                  <Filter className="w-4 h-4 mr-1" />
                  Filter
                </Button>
                
                {selectedTags.length > 0 && (
                  <Badge variant="secondary" className="text-xs">
                    {selectedTags.length} selected
                  </Badge>
                )}
              </div>

              {showTagFilter && (
                <div className="mt-3 p-3 bg-muted rounded-lg">
                  <div className="flex flex-wrap gap-1">
                    {availableTags.map(tag => (
                      <button
                        key={tag}
                        onClick={() => {
                          setSelectedTags(prev =>
                            prev.includes(tag)
                              ? prev.filter(t => t !== tag)
                              : [...prev, tag]
                          );
                        }}
                        className={`px-2 py-1 rounded-full text-xs transition-colors ${
                          selectedTags.includes(tag)
                            ? 'bg-primary text-primary-foreground'
                            : 'bg-background hover:bg-accent'
                        }`}
                      >
                        <Hash className="w-2 h-2 mr-1 inline" />
                        {tag}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>

            <ScrollArea className="flex-1">
              {renderDialogList(personalDialogs)}
            </ScrollArea>
          </div>

          {/* Chat Content */}
          <div className="flex-1 flex flex-col">
            {selectedDialog ? (
              <>
                {/* Chat Header */}
                <div className="border-b border-border p-4 flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Avatar className="w-10 h-10">
                      <AvatarImage src={personalDialogs.find(d => d.id === selectedDialog)?.avatar} />
                      <AvatarFallback>
                        {personalDialogs.find(d => d.id === selectedDialog)?.name.charAt(0)}
                      </AvatarFallback>
                    </Avatar>
                    <div>
                      <h3 className="font-medium">
                        {personalDialogs.find(d => d.id === selectedDialog)?.name}
                      </h3>
                      <p className="text-sm text-muted-foreground">
                        {personalDialogs.find(d => d.id === selectedDialog)?.isOnline ? 'Online' : 'Last seen recently'}
                      </p>
                    </div>
                  </div>

                  <ChatHeaderActions
                    selectedDialog={personalDialogs.find(d => d.id === selectedDialog)}
                    chatActions={chatActions}
                    onVideoCall={() => onVideoCall?.(personalDialogs.find(d => d.id === selectedDialog))}
                    onVoiceCall={() => onVoiceCall?.(personalDialogs.find(d => d.id === selectedDialog))}
                    onMuteChat={handleMuteChat}
                    onPinChat={handlePinChat}
                    onArchiveChat={handleArchiveChat}
                    onBlockContact={handleBlockContact}
                    onClearHistory={handleClearHistory}
                    toggleSelectionMode={toggleSelectionMode}
                    isSelectionMode={isSelectionMode}
                  />
                </div>

                {/* Messages */}
                <ScrollArea className="flex-1 p-4">
                  <ChatMessages
                    messages={selectedDialog === 'family' ? familyMessages : []}
                    selectedMessages={selectedMessages}
                    isSelectionMode={isSelectionMode}
                    onReplyToMessage={handleReplyToMessage}
                    onEditMessage={handleEditMessage}
                    onDeleteMessage={handleDeleteMessage}
                    onForwardMessage={handleForwardMessage}
                    onCopyMessage={handleCopyMessage}
                    onPinMessage={handlePinMessage}
                    onSelectMessage={handleSelectMessage}
                  />
                </ScrollArea>

                {/* Message Input */}
                <ChatInput
                  messageInput={messageInput}
                  setMessageInput={setMessageInput}
                  onSendMessage={handleSendMessage}
                  replyToMessage={replyToMessage}
                  setReplyToMessage={setReplyToMessage}
                  editingMessage={editingMessage}
                  setEditingMessage={setEditingMessage}
                  isRecording={isRecording}
                  recordingDuration={recordingDuration}
                  onStartRecording={startVoiceRecording}
                  onStopRecording={stopVoiceRecording}
                  onCancelRecording={cancelVoiceRecording}
                />
              </>
            ) : (
              <div className="flex-1 flex items-center justify-center">
                <div className="text-center">
                  <MessageCircle className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                  <p className="text-muted-foreground">Select a chat to start messaging</p>
                </div>
              </div>
            )}
          </div>
        </TabsContent>

        {/* Work Tab */}
        <TabsContent value="work" className="flex-1 flex flex-col">
          {/* Prominent Workspace Switcher */}
          <WorkspaceSwitcher
            selectedWorkspace={selectedWorkspace || ''}
            userWorkspaces={userWorkspaces}
            onWorkspaceChange={onWorkspaceChange || (() => {})}
            variant="prominent"
            showMetrics={true}
          />

          <div className="flex flex-1">
            {/* Business Dialogs List */}
            <div className="w-80 border-r border-border flex flex-col">
              <div className="p-4 border-b border-border">
                <div className="flex items-center gap-2 mb-3">
                  <div className="flex-1">
                    <Input
                      placeholder="Search customers..."
                      className="w-full"
                      startIcon={<Search className="w-4 h-4" />}
                    />
                  </div>
                  <Button size="icon" variant="ghost">
                    <UserPlus className="w-4 h-4" />
                  </Button>
                </div>

                {/* Work-specific filters */}
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm">
                    <Target className="w-4 h-4 mr-1" />
                    Hot Leads
                  </Button>
                  <Button variant="outline" size="sm">
                    <Crown className="w-4 h-4 mr-1" />
                    VIP
                  </Button>
                </div>
              </div>

              <ScrollArea className="flex-1">
                {businessDialogs.length > 0 ? (
                  renderDialogList(businessDialogs, true)
                ) : (
                  <div className="p-4 text-center">
                    <Building2 className="w-8 h-8 text-muted-foreground mx-auto mb-2" />
                    <p className="text-sm text-muted-foreground">No customers yet</p>
                    <p className="text-xs text-muted-foreground">Start building relationships!</p>
                  </div>
                )}
              </ScrollArea>
            </div>

            {/* Work Content */}
            <div className="flex-1 flex flex-col">
              {selectedDialog && businessDialogs.find(d => d.id === selectedDialog) ? (
                <>
                  {/* Business Chat Header */}
                  <div className="border-b border-border p-4 flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Avatar className="w-10 h-10">
                        <AvatarImage src={businessDialogs.find(d => d.id === selectedDialog)?.avatar} />
                        <AvatarFallback>
                          {businessDialogs.find(d => d.id === selectedDialog)?.name.charAt(0)}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <div className="flex items-center gap-2">
                          <h3 className="font-medium">
                            {businessDialogs.find(d => d.id === selectedDialog)?.name}
                          </h3>
                          {businessDialogs.find(d => d.id === selectedDialog)?.leadScore && (
                            <div className={`flex items-center gap-1 px-2 py-1 rounded-full text-xs ${getLeadScoreColor(businessDialogs.find(d => d.id === selectedDialog)?.leadScore || 0)}`}>
                              {getLeadScoreIcon(businessDialogs.find(d => d.id === selectedDialog)?.leadScore || 0)}
                              <span>{businessDialogs.find(d => d.id === selectedDialog)?.leadScore}</span>
                            </div>
                          )}
                        </div>
                        <p className="text-sm text-muted-foreground">
                          {businessDialogs.find(d => d.id === selectedDialog)?.lastActivity}
                        </p>
                      </div>
                    </div>

                    <ChatHeaderActions
                      selectedDialog={businessDialogs.find(d => d.id === selectedDialog)}
                      chatActions={chatActions}
                      onVideoCall={() => onVideoCall?.(businessDialogs.find(d => d.id === selectedDialog))}
                      onVoiceCall={() => onVoiceCall?.(businessDialogs.find(d => d.id === selectedDialog))}
                      onMuteChat={handleMuteChat}
                      onPinChat={handlePinChat}
                      onArchiveChat={handleArchiveChat}
                      onBlockContact={handleBlockContact}
                      onClearHistory={handleClearHistory}
                      toggleSelectionMode={toggleSelectionMode}
                      isSelectionMode={isSelectionMode}
                    />
                  </div>

                  {/* Business Messages */}
                  <ScrollArea className="flex-1 p-4">
                    <ChatMessages
                      messages={businessMessages}
                      selectedMessages={selectedMessages}
                      isSelectionMode={isSelectionMode}
                      onReplyToMessage={handleReplyToMessage}
                      onEditMessage={handleEditMessage}
                      onDeleteMessage={handleDeleteMessage}
                      onForwardMessage={handleForwardMessage}
                      onCopyMessage={handleCopyMessage}
                      onPinMessage={handlePinMessage}
                      onSelectMessage={handleSelectMessage}
                    />
                  </ScrollArea>

                  {/* Business Message Input */}
                  <ChatInput
                    messageInput={messageInput}
                    setMessageInput={setMessageInput}
                    onSendMessage={handleSendMessage}
                    replyToMessage={replyToMessage}
                    setReplyToMessage={setReplyToMessage}
                    editingMessage={editingMessage}
                    setEditingMessage={setEditingMessage}
                    isRecording={isRecording}
                    recordingDuration={recordingDuration}
                    onStartRecording={startVoiceRecording}
                    onStopRecording={stopVoiceRecording}
                    onCancelRecording={cancelVoiceRecording}
                  />
                </>
              ) : (
                renderWorkDashboard()
              )}
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}