package handlers

import (
	"tchat.dev/streaming/repository"
	"tchat.dev/streaming/services"
)

// Handlers groups all HTTP handlers
type Handlers struct {
	CreateStream            *CreateStreamHandler
	ListStreams             *ListStreamsHandler
	GetStream               *GetStreamHandler
	UpdateStream            *UpdateStreamHandler
	StartStream             *StartStreamHandler
	EndStream               *EndStreamHandler
	SendChat                *SendChatHandler
	GetChat                 *GetChatHandler
	DeleteChat              *DeleteChatHandler
	SendReaction            *SendReactionHandler
	FeatureProduct          *FeatureProductHandler
	ListProducts            *ListProductsHandler
	GetAnalytics            *GetAnalyticsHandler
	NotificationPreferences *NotificationPreferencesHandler
	SignalingService        *services.SignalingService
}

// NewHandlers initializes all handlers with dependencies
func NewHandlers(
	liveStreamRepo repository.LiveStreamRepositoryInterface,
	chatRepo repository.ChatMessageRepository,
	reactionRepo repository.StreamReactionRepository,
	productRepo repository.StreamProductRepository,
	analyticsRepo repository.StreamAnalyticsRepository,
	prefRepo repository.NotificationPreferenceRepository,
	webrtcService *services.WebRTCService,
	signalingService *services.SignalingService,
	recordingService *services.RecordingService,
	kycService *services.KYCService,
) *Handlers {
	return &Handlers{
		CreateStream:            NewCreateStreamHandler(liveStreamRepo, kycService, webrtcService),
		ListStreams:             NewListStreamsHandler(liveStreamRepo),
		GetStream:               NewGetStreamHandler(liveStreamRepo, productRepo),
		UpdateStream:            NewUpdateStreamHandler(liveStreamRepo),
		StartStream:             NewStartStreamHandler(liveStreamRepo, webrtcService, signalingService, recordingService, kycService),
		EndStream:               NewEndStreamHandler(liveStreamRepo, recordingService),
		SendChat:                NewSendChatHandler(liveStreamRepo, chatRepo, signalingService),
		GetChat:                 NewGetChatHandler(liveStreamRepo, chatRepo),
		DeleteChat:              NewDeleteChatHandler(liveStreamRepo, chatRepo),
		SendReaction:            NewSendReactionHandler(liveStreamRepo, reactionRepo, signalingService),
		FeatureProduct:          NewFeatureProductHandler(liveStreamRepo, productRepo),
		ListProducts:            NewListProductsHandler(liveStreamRepo, productRepo),
		GetAnalytics:            NewGetAnalyticsHandler(liveStreamRepo, analyticsRepo),
		NotificationPreferences: NewNotificationPreferencesHandler(prefRepo),
		SignalingService:        signalingService,
	}
}