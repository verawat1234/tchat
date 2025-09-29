package providers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"tchat.dev/notification/models"
)

// SMSProvider implements ChannelProvider for SMS notifications
type SMSProvider struct {
	twilioAccountSID string
	twilioAuthToken  string
	twilioFromNumber string
}

// SMSConfig holds SMS provider configuration
type SMSConfig struct {
	TwilioAccountSID string `json:"twilio_account_sid"`
	TwilioAuthToken  string `json:"twilio_auth_token"`
	TwilioFromNumber string `json:"twilio_from_number"`
}

// NewSMSProvider creates a new SMS provider
func NewSMSProvider(config SMSConfig) *SMSProvider {
	return &SMSProvider{
		twilioAccountSID: config.TwilioAccountSID,
		twilioAuthToken:  config.TwilioAuthToken,
		twilioFromNumber: config.TwilioFromNumber,
	}
}

// Send sends an SMS notification
func (s *SMSProvider) Send(ctx context.Context, notification *models.Notification) error {
	if notification.RecipientPhone == "" {
		return fmt.Errorf("recipient phone number is required for SMS notifications")
	}

	// Build SMS content
	message := s.buildSMSMessage(notification)

	// Send SMS via Twilio
	return s.sendTwilioSMS(ctx, notification.RecipientPhone, message)
}

// SendBatch sends multiple SMS notifications
func (s *SMSProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	// Twilio doesn't have native batch SMS, so we send individually
	for _, notification := range notifications {
		if err := s.Send(ctx, notification); err != nil {
			log.Printf("Failed to send SMS to %s: %v", notification.RecipientPhone, err)
			// Continue with other messages
		}
	}
	return nil
}

// ValidateConfig validates the SMS provider configuration
func (s *SMSProvider) ValidateConfig() error {
	if s.twilioAccountSID == "" {
		return fmt.Errorf("Twilio Account SID is required")
	}
	if s.twilioAuthToken == "" {
		return fmt.Errorf("Twilio Auth Token is required")
	}
	if s.twilioFromNumber == "" {
		return fmt.Errorf("Twilio from number is required")
	}
	return nil
}

// GetProviderName returns the provider name
func (s *SMSProvider) GetProviderName() string {
	return "sms"
}

// SupportsChannels returns the channels this provider supports
func (s *SMSProvider) SupportsChannels() []models.NotificationChannel {
	return []models.NotificationChannel{models.NotificationChannelSMS}
}

// buildSMSMessage builds the SMS message content
func (s *SMSProvider) buildSMSMessage(notification *models.Notification) string {
	var message strings.Builder

	// Add title if present
	if notification.Title != "" {
		message.WriteString(notification.Title)
		if notification.Body != "" {
			message.WriteString(": ")
		}
	}

	// Add body
	if notification.Body != "" {
		message.WriteString(notification.Body)
	}

	// Add action URL if present
	if notification.ActionURL != "" {
		message.WriteString(fmt.Sprintf("\n\nLink: %s", notification.ActionURL))
	}

	// Truncate if too long (SMS limit is usually 160 characters)
	content := message.String()
	if len(content) > 160 {
		content = content[:157] + "..."
	}

	return content
}

// sendTwilioSMS sends an SMS via Twilio API
func (s *SMSProvider) sendTwilioSMS(ctx context.Context, to, message string) error {
	// This would use the Twilio Go SDK to send SMS
	// For now, it's a placeholder implementation
	log.Printf("Twilio SMS: Sending to %s: %s", to, message)

	// In a real implementation, this would make an HTTP request to Twilio
	// Example Twilio API call structure:
	/*
		client := twilio.NewRestClient(s.twilioAccountSID, s.twilioAuthToken)
		params := &openapi.CreateMessageParams{}
		params.SetTo(to)
		params.SetFrom(s.twilioFromNumber)
		params.SetBody(message)

		_, err := client.Api.CreateMessage(params)
		return err
	*/

	return nil
}

// AWSSSMSProvider implements SMS provider using AWS SNS
type AWSSMSProvider struct {
	region    string
	accessKey string
	secretKey string
}

// NewAWSSMSProvider creates a new AWS SMS provider
func NewAWSSMSProvider(region, accessKey, secretKey string) *AWSSMSProvider {
	return &AWSSMSProvider{
		region:    region,
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// Send sends an SMS via AWS SNS
func (a *AWSSMSProvider) Send(ctx context.Context, notification *models.Notification) error {
	if notification.RecipientPhone == "" {
		return fmt.Errorf("recipient phone number is required for SMS notifications")
	}

	message := a.buildSMSMessage(notification)
	log.Printf("AWS SNS SMS: Sending to %s: %s", notification.RecipientPhone, message)

	// In a real implementation, this would use AWS SDK
	/*
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(a.region),
		}))
		svc := sns.New(sess)

		input := &sns.PublishInput{
			Message:     aws.String(message),
			PhoneNumber: aws.String(notification.RecipientPhone),
		}

		_, err := svc.Publish(input)
		return err
	*/

	return nil
}

// SendBatch sends multiple SMS notifications via AWS SNS
func (a *AWSSMSProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	for _, notification := range notifications {
		if err := a.Send(ctx, notification); err != nil {
			log.Printf("Failed to send AWS SMS to %s: %v", notification.RecipientPhone, err)
		}
	}
	return nil
}

// ValidateConfig validates AWS SMS configuration
func (a *AWSSMSProvider) ValidateConfig() error {
	if a.region == "" {
		return fmt.Errorf("AWS region is required")
	}
	if a.accessKey == "" {
		return fmt.Errorf("AWS access key is required")
	}
	if a.secretKey == "" {
		return fmt.Errorf("AWS secret key is required")
	}
	return nil
}

// GetProviderName returns the provider name
func (a *AWSSMSProvider) GetProviderName() string {
	return "aws-sms"
}

// SupportsChannels returns supported channels
func (a *AWSSMSProvider) SupportsChannels() []models.NotificationChannel {
	return []models.NotificationChannel{models.NotificationChannelSMS}
}

// buildSMSMessage builds SMS message for AWS
func (a *AWSSMSProvider) buildSMSMessage(notification *models.Notification) string {
	var message strings.Builder

	if notification.Title != "" {
		message.WriteString(notification.Title)
		if notification.Body != "" {
			message.WriteString(": ")
		}
	}

	if notification.Body != "" {
		message.WriteString(notification.Body)
	}

	if notification.ActionURL != "" {
		message.WriteString(fmt.Sprintf("\n\nLink: %s", notification.ActionURL))
	}

	content := message.String()
	if len(content) > 160 {
		content = content[:157] + "..."
	}

	return content
}

// LineNotifyProvider implements notifications using LINE Notify (popular in Southeast Asia)
type LineNotifyProvider struct {
	accessToken string
}

// NewLineNotifyProvider creates a new LINE Notify provider
func NewLineNotifyProvider(accessToken string) *LineNotifyProvider {
	return &LineNotifyProvider{
		accessToken: accessToken,
	}
}

// Send sends a notification via LINE Notify
func (l *LineNotifyProvider) Send(ctx context.Context, notification *models.Notification) error {
	message := l.buildLineMessage(notification)
	log.Printf("LINE Notify: Sending message: %s", message)

	// In a real implementation, this would make an HTTP request to LINE Notify API
	/*
		data := url.Values{}
		data.Set("message", message)

		req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", strings.NewReader(data.Encode()))
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", "Bearer "+l.accessToken)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		return nil
	*/

	return nil
}

// SendBatch sends multiple notifications via LINE Notify
func (l *LineNotifyProvider) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	for _, notification := range notifications {
		if err := l.Send(ctx, notification); err != nil {
			log.Printf("Failed to send LINE notification: %v", err)
		}
	}
	return nil
}

// ValidateConfig validates LINE Notify configuration
func (l *LineNotifyProvider) ValidateConfig() error {
	if l.accessToken == "" {
		return fmt.Errorf("LINE Notify access token is required")
	}
	return nil
}

// GetProviderName returns the provider name
func (l *LineNotifyProvider) GetProviderName() string {
	return "line-notify"
}

// SupportsChannels returns supported channels
func (l *LineNotifyProvider) SupportsChannels() []models.NotificationChannel {
	return []models.NotificationChannel{models.NotificationChannelSMS} // LINE Notify acts as messaging
}

// buildLineMessage builds LINE Notify message
func (l *LineNotifyProvider) buildLineMessage(notification *models.Notification) string {
	var message strings.Builder

	if notification.Title != "" {
		message.WriteString(fmt.Sprintf("📱 %s\n", notification.Title))
	}

	if notification.Body != "" {
		message.WriteString(notification.Body)
	}

	if notification.ActionURL != "" {
		message.WriteString(fmt.Sprintf("\n\n🔗 %s", notification.ActionURL))
	}

	return message.String()
}